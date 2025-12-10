package main

import (
	"context"
	"crypto/sha1"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	mathr "math/rand"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xtaci/kcp-go"
	"github.com/xtaci/smux"
	"golang.org/x/crypto/pbkdf2"
	"gopkg.in/yaml.v3"
)

var (
	configs map[string]interface{}
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func readConfig() {
	f, err := os.ReadFile("Config/config.yml")

	if err != nil {
		panic(err)
	}

	var data map[string]interface{}
	err = yaml.Unmarshal(f, &data)

	if err != nil {
		panic(err)
	}
	configs = data
}

func copyStream(clientConn *websocket.Conn, st *state) {
	first := true
	for {
		messageType, payload, err := clientConn.ReadMessage()

		if err != nil {
			return
		}

		var clientID clientID = clientID(binary.BigEndian.Uint32(payload[:8]))
		payload = payload[8:]
		if first {
			st.clients.update(clientID, clientConn)
			first = false
		}
		if messageType == websocket.BinaryMessage {
			st.QueueIncoming(payload, clientID)
		}
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request, st *state) {
	conn_client, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	go copyStream(conn_client, st)
}

// TURBO TUNNEL
type clientID uint32

func newClientID() clientID {
	return clientID(mathr.Uint32())
}

func (addr clientID) Network() string {
	return "clientid"
}

func (addr clientID) String() string {
	return fmt.Sprintf("%08x", uint32(addr))
}

type state struct {
	closed    chan struct{}
	recvQueue chan *taggetPacket
	localAddr net.Addr
	clients   *clientMap
}

type taggetPacket struct {
	data []byte
	addr net.Addr
}

func newQueuePacketConn() *state {
	id := newClientID()
	return &state{
		localAddr: id,
		recvQueue: make(chan *taggetPacket, 100),
		closed:    make(chan struct{}),
		clients:   newClientMap(),
	}
}

func (st *state) QueueIncoming(p []byte, id net.Addr) {
	st.recvQueue <- &taggetPacket{data: p, addr: id}
}

func (st *state) ReadFrom(p []byte) (int, net.Addr, error) {
	select {
	case packet := <-st.recvQueue:
		return copy(p, packet.data), packet.addr, nil
	case <-st.closed:
		return 0, nil, &net.OpError{Op: "read", Net: st.LocalAddr().Network(), Source: st.LocalAddr(), Err: errors.New("closed conn")}
	}
}

func (st *state) WriteTo(p []byte, addr net.Addr) (int, error) {
	select {
	case <-st.closed:
		return 0, &net.OpError{Op: "write", Net: addr.Network(), Source: st.LocalAddr(), Addr: addr, Err: errors.New("closed conn")}
	default:
	}

	conn := st.clients.get(addr)
	if conn == nil {
		return 0, &net.OpError{Op: "write", Net: addr.Network(), Addr: addr, Source: st.LocalAddr(), Err: fmt.Errorf("no mapped net.Conn")}
	}

	conn.WriteMessage(websocket.BinaryMessage, p)
	return len(p), nil
}

func (st *state) Close() error {
	select {
	case <-st.closed:
		return &net.OpError{Op: "close", Net: st.LocalAddr().Network(), Addr: st.LocalAddr(), Err: errors.New("closed conn")}
	default:
		close(st.closed)
		return nil
	}
}

func (st *state) LocalAddr() net.Addr {
	return st.localAddr
}

func (st *state) SetDeadline(t time.Time) error {
	return errors.New("not implemented")
}

func (st *state) SetReadDeadline(t time.Time) error {
	return errors.New("not implemented")
}

func (st *state) SetWriteDeadline(t time.Time) error {
	return errors.New("not implemented")
}

func handleSOCKS5(stream *smux.Stream) *net.TCPConn {
	buf := make([]byte, 2)
	if _, err := io.ReadFull(stream, buf); err != nil {
		fmt.Println("Failed to read version identifier/method selection:", err)
		return nil
	}

	if buf[0] != 0x05 {
		fmt.Println("Unsupported SOCKS version:", buf[0])
		return nil
	}

	methods := make([]byte, buf[1])
	if _, err := io.ReadFull(stream, methods); err != nil {
		fmt.Println("Failed to read authentication methods:", err)
		return nil
	}

	if _, err := stream.Write([]byte{0x05, 0x00}); err != nil {
		fmt.Println("Failed to write method selection response:", err)
		return nil
	}

	buf = make([]byte, 4)
	if _, err := io.ReadFull(stream, buf); err != nil {
		fmt.Println("Failed to read SOCKS5 request:", err)
		return nil
	}

	if buf[1] != 0x01 {
		fmt.Println("Unsupported SOCKS command:", buf[1])
		return nil
	}

	var targetAddr string
	switch buf[3] {
	case 0x01: // IPv4
		addr := make([]byte, 4)
		if _, err := io.ReadFull(stream, addr); err != nil {
			fmt.Println("Failed to read IPv4 address:", err)
			return nil
		}
		targetAddr = net.IP(addr).String()
	case 0x03: // Domain name
		addrLen := make([]byte, 1)
		if _, err := io.ReadFull(stream, addrLen); err != nil {
			fmt.Println("Failed to read domain name length:", err)
			return nil
		}
		addr := make([]byte, addrLen[0])
		if _, err := io.ReadFull(stream, addr); err != nil {
			fmt.Println("Failed to read domain name:", err)
			return nil
		}
		targetAddr = string(addr)
	case 0x04: // IPv6
		addr := make([]byte, 16)
		if _, err := io.ReadFull(stream, addr); err != nil {
			fmt.Println("Failed to read IPv6 address:", err)
			return nil
		}
		targetAddr = net.IP(addr).String()
	default:
		fmt.Println("Unsupported address type:", buf[3])
		return nil
	}

	port := make([]byte, 2)
	if _, err := io.ReadFull(stream, port); err != nil {
		fmt.Println("Failed to read port:", err)
		return nil
	}
	targetAddr = fmt.Sprintf("%s:%d", targetAddr, int(port[0])<<8|int(port[1]))

	dialer := &net.Dialer{}

	targetConn, err := dialer.DialContext(context.Background(), "tcp", targetAddr)
	if err != nil {
		fmt.Println("Failed to connect to target address:", err)
		stream.Write([]byte{0x05, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		return nil
	}

	stream.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	return targetConn.(*net.TCPConn)
}

func handleStream(stream *smux.Stream) error {
	conn := handleSOCKS5(stream)

	if conn == nil {
		return errors.New("error in socks5")
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		io.Copy(conn, stream)
		stream.Close()
		conn.Close()
	}()
	go func() {
		defer wg.Done()
		io.Copy(stream, conn)
		stream.Close()
		conn.Close()
	}()
	wg.Wait()
	return nil
}

func acceptStreams(sess *smux.Session) error {
	for {
		stream, err := sess.AcceptStream()
		if err != nil {
			if err, ok := err.(*net.OpError); ok && err.Temporary() {
				log.Printf("temporary error in sess.AcceptStream: %v", err)
				continue
			}
			return err
		}

		go func() {
			defer stream.Close()
			log.Printf("begin stream %v (session %v)", stream.ID(), stream.RemoteAddr())
			err := handleStream(stream)
			if err != nil {
				log.Printf("error in handleStream: %v", err)
			}
			log.Printf("end stream %v (session %v)", stream.ID(), stream.RemoteAddr())
		}()
	}
}

func acceptSessions(ln *kcp.Listener) error {
	for {
		conn, err := ln.AcceptKCP()
		if err != nil {
			if err, ok := err.(*net.OpError); ok && err.Temporary() {
				log.Printf("temporary error in ln.Accept: %v", err)
				continue
			}
			return err
		}

		conn.SetNoDelay(1, 10, 2, 1)

		go func() {
			defer conn.Close()

			sess, err := smux.Server(conn, &smux.Config{
				Version:           1,
				KeepAliveInterval: 10 * time.Second,
				KeepAliveTimeout:  100 * time.Second,
				MaxFrameSize:      32768,
				MaxReceiveBuffer:  4194304,
				MaxStreamBuffer:   65536})
			if err != nil {
				log.Printf("error in smux.Server: %v", err)
				return
			}
			defer sess.Close()

			log.Printf("begin session %v", sess.RemoteAddr())
			err = acceptStreams(sess)
			if err != nil {
				log.Printf("error in acceptStreams: %v", err)
			}
			log.Printf("end session %v", sess.RemoteAddr())
		}()
	}
}

func run(conn *state) error {
	defer conn.Close()

	key := pbkdf2.Key([]byte("demo pass"), []byte("demo salt"), 1024, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)

	ln, err := kcp.ServeConn(block, 0, 0, conn)
	if err != nil {
		return err
	}
	defer ln.Close()

	return acceptSessions(ln)
}

func main() {
	readConfig()

	conn := newQueuePacketConn()
	go run(conn)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleConnections(w, r, conn)
	})

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS13,
	}

	server := &http.Server{
		Addr:      fmt.Sprintf("0.0.0.0:%s", configs["localport"]),
		TLSConfig: tlsConfig,
	}

	fmt.Println("Ready")

	err := server.ListenAndServeTLS("Config/server.crt", "Config/server.key")
	fmt.Println(err)
}

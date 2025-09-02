package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/sha1"
	"crypto/tls"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"math"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/fogleman/gg"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media"
	"github.com/xtaci/kcp-go"
	"github.com/xtaci/smux"
	"golang.org/x/crypto/pbkdf2"
	"gopkg.in/yaml.v2"
)

const (
	CLIENT               = 2
	SIZE_OF_HEADER       = 11
	PICTURE_ID_MAX_VALUE = 32767
	PINGDEADLINE         = time.Second * 10

	WEB   = 3
	PROXY = 1

	width     = 1280
	height    = 720
	fps       = 30
	frameSize = width * height * 3 / 2
)

type state struct {
	connSignaling  *websocket.Conn
	peerConnection *webrtc.PeerConnection
	conn           *rOConn

	pendingCandidate []*webrtc.ICECandidateInit

	doneSingalingBool bool
	doneSignaling     chan int
	signalingConnMux  sync.Mutex

	iceConnectedCtx       context.Context
	iceConnectedCtxCancel context.CancelFunc

	sequenceNumer     uint32
	pictureIDTwoBytes uint64

	fragmentedPackets map[uint32]*FragmentedPacket

	closed chan int
}

type FragmentedPacket struct {
	data      map[uint8][]byte
	lastChunk uint8
}

var (
	configs   map[string]interface{}
	proxyType int
)

var first = true

func (s *state) encapsulateWeb(remaing, frame []byte, chunkNumber uint8, conn *rOConn) ([]byte, []byte, uint8) {
	var lenFrame int = len(frame)
	var result []byte = make([]byte, 0)
	var bypassBytes = 10
	if first {
		frame[0] = frame[0] & 0b11111110
		first = false
	}

	lenFrame = lenFrame - 10
	result = append(result, frame[:10]...)

	if (conn.lenOutgoingQueue() == 0 && len(remaing) == 0) || lenFrame <= SIZE_OF_HEADER {
		frame[bypassBytes] = 0
		return frame, nil, 0
	}

	var data []byte = remaing
	var reamaingArray []byte = nil

	for lenFrame > SIZE_OF_HEADER {
		if len(data) == 0 {
			data = conn.OutgoingQueue()
			if len(data) == 0 {
				break
			}
		}

		//FLAG PACKET HAS CONTENT
		result = append(result, byte(1))

		//SEQUENCE NUMBER OF THE PACKET
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, uint32(s.sequenceNumer))
		result = append(result, bytes...)

		//CHUNK OF THE PACKET
		result = append(result, byte(chunkNumber))
		chunkNumber++

		//LEN OF THE DATA IN THE PACKET
		lenData := uint32(min(len(data), max(lenFrame-SIZE_OF_HEADER, 0)))
		bytes = make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, lenData)
		result = append(result, bytes...)

		//FLAG LAST CHUNK
		finalChunk := 0
		if len(data) == int(lenData) {
			finalChunk = 1
		}
		result = append(result, byte(finalChunk))

		//DATA
		result = append(result, data[:lenData]...)

		data = data[lenData:]
		lenFrame = lenFrame - int(lenData) - SIZE_OF_HEADER

		if len(data) == 0 {
			s.sequenceNumer = s.sequenceNumer + 1
			chunkNumber = 0
			data = nil
		}
	}

	if len(result) < len(frame) {
		aux := make([]byte, len(frame)-len(result))
		result = append(result, aux...)
	}

	if len(data) > 0 {
		reamaingArray = data
	}

	s.pictureIDTwoBytes++
	return result, reamaingArray, chunkNumber
}

func (s *state) desencapsulateWeb(frame []byte, conn *rOConn) {
	var lenData uint32 = 0
	var sequenceNumber uint32 = 0
	var chunk uint8 = 0
	var finalChunk uint8 = 0
	var data []byte = make([]byte, 0)

	for i := 4; i < len(frame) && frame[i] == 1; i += (int(lenData) + SIZE_OF_HEADER) {
		sequenceNumber = binary.BigEndian.Uint32([]byte{frame[i+1], frame[i+2], frame[i+3], frame[i+4]})
		chunk = frame[i+5]
		lenData = binary.BigEndian.Uint32([]byte{frame[i+6], frame[i+7], frame[i+8], frame[i+9]})
		finalChunk = frame[i+10]
		data = frame[i+SIZE_OF_HEADER : i+SIZE_OF_HEADER+int(lenData)]

		packet := s.reconstructPacket(sequenceNumber, chunk, finalChunk, data)
		if len(packet) != 0 {
			conn.QueueIncoming(packet)
		}
	}
}

func (s *state) encapsulate(remaing, frame []byte, chunkNumber uint8, conn *rOConn) ([]byte, []byte, uint8) {
	var value_inc int
	if s.pictureIDTwoBytes > PICTURE_ID_MAX_VALUE {
		s.pictureIDTwoBytes = 0
	}

	//FRAMES ARE DIVIDED INTO CHUNCKS OF DIFFERENT SIZES DEPNDING ON THE PICTURE ID (A FLAG ON THE FRAME). EACH CHUNK WILL BE SENT IN AN INDIVIDUAL RTP PACKET.
	if s.pictureIDTwoBytes == 0 {
		value_inc = 1187
	} else if s.pictureIDTwoBytes > 0 && s.pictureIDTwoBytes < 128 {
		value_inc = 1185
	} else {
		value_inc = 1184
	}

	//IF FRAME IS TO SMALL OR NO CONTENT NEEDS TO BE SENT, SEND FRAME AS IS BUT WITH FLAG PACKET HAS CONTENT FALSE
	if (conn.lenOutgoingQueue() == 0 && len(remaing) == 0) || len(frame) <= SIZE_OF_HEADER {

		for i := 0; i < len(frame); i += value_inc {
			frame[i] = 0
		}

		s.pictureIDTwoBytes++
		return frame, nil, 0
	}

	//LOOP TO ADD CONTENT TO FRAMES
	var data []byte = remaing
	var reamaingArray []byte = nil
	var result []byte = make([]byte, 0)
	var remaingSize int = value_inc
	var lenFrame int = len(frame)
	for lenFrame > SIZE_OF_HEADER {
		if len(data) == 0 {
			data = conn.OutgoingQueue()
			if len(data) == 0 {
				break
			}
		}

		//FLAG PACKET HAS CONTENT
		result = append(result, byte(1))

		//SEQUENCE NUMBER OF THE PACKET
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, uint32(s.sequenceNumer))
		result = append(result, bytes...)

		//CHUNK OF THE PACKET
		result = append(result, byte(chunkNumber))
		chunkNumber++

		//LEN OF THE DATA IN THE PACKET
		lenData := uint32(min(len(data), max(remaingSize-SIZE_OF_HEADER, 0), max(lenFrame-SIZE_OF_HEADER, 0)))
		bytes = make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, lenData)
		result = append(result, bytes...)

		//FLAG LAST CHUNK
		finalChunk := 0
		if len(data) == int(lenData) {
			finalChunk = 1
		}
		result = append(result, byte(finalChunk))

		//DATA
		result = append(result, data[:lenData]...)

		//UPDATES FOR NEXT ITERATION
		data = data[lenData:]
		lenFrame = lenFrame - int(lenData) - SIZE_OF_HEADER
		remaingSize = remaingSize - int(lenData) - SIZE_OF_HEADER

		if remaingSize <= SIZE_OF_HEADER {
			filler := make([]byte, remaingSize)
			result = append(result, filler...)
			remaingSize = value_inc
			lenFrame = lenFrame - SIZE_OF_HEADER
		}

		if len(data) == 0 {
			s.sequenceNumer = s.sequenceNumer + 1
			chunkNumber = 0
			data = nil
		}
	}

	if len(result) < len(frame) {
		aux := make([]byte, len(frame)-len(result))
		result = append(result, aux...)
	}

	if len(data) > 0 {
		reamaingArray = data
	}

	s.pictureIDTwoBytes++
	return result, reamaingArray, chunkNumber
}

func (s *state) decapsulate(frame []byte, conn *rOConn) {
	var headerSizeBytes = 0

	//VP8 FRAME HEADER, THE VALUE TOTAL VALUE OF THE PAYLOAD DEPENDS ON SOME FLAGS
	if frame[0]&0b10010000 == 0b00010000 {
		headerSizeBytes++
	}

	if frame[0]&0b10000000 == 0b10000000 {
		headerSizeBytes = headerSizeBytes + 2

		if frame[1]&0b10000000 == 0b10000000 {
			headerSizeBytes++

			if frame[2] >= 128 {
				headerSizeBytes++
			}
		}

		if frame[1]&0b01000000 == 0b01000000 {
			headerSizeBytes++
		}

		if frame[1]&0b00100000 == 0b00100000 || frame[1]&0b00010000 == 0b00010000 {
			headerSizeBytes++
		}
	}

	var lenData uint32 = 0
	var sequenceNumber uint32 = 0
	var chunk uint8 = 0
	var finalChunk uint8 = 0
	var data []byte = make([]byte, 0)

	for i := headerSizeBytes; i < len(frame) && frame[i] == 1; i += (int(lenData) + SIZE_OF_HEADER) {
		sequenceNumber = binary.BigEndian.Uint32([]byte{frame[i+1], frame[i+2], frame[i+3], frame[i+4]})
		chunk = frame[i+5]
		lenData = binary.BigEndian.Uint32([]byte{frame[i+6], frame[i+7], frame[i+8], frame[i+9]})
		finalChunk = frame[i+10]
		data = frame[i+SIZE_OF_HEADER : i+SIZE_OF_HEADER+int(lenData)]

		packet := s.reconstructPacket(sequenceNumber, chunk, finalChunk, data)
		if len(packet) != 0 {
			conn.QueueIncoming(packet)
		}
	}
}

func (s *state) reconstructPacket(sequenceNumber uint32, chunk, finalChunk uint8, data []byte) []byte {
	if chunk == 0 && finalChunk == 1 {
		return data
	}

	result := make([]byte, 0)
	packet, exists := s.fragmentedPackets[sequenceNumber]
	if !exists {
		fragmentedP := FragmentedPacket{
			data:      make(map[uint8][]byte),
			lastChunk: 0,
		}
		fragmentedP.data[chunk] = data
		s.fragmentedPackets[sequenceNumber] = &fragmentedP

		//Keep cleaning up stuff
		maxValueSequenceNumber := uint32(uint64(1<<32) - 1)
		getSymmetricPosition := maxValueSequenceNumber - sequenceNumber
		delete(s.fragmentedPackets, getSymmetricPosition)

		result = nil
	} else {
		_, exists = packet.data[chunk]
		if !exists {
			packet.data[chunk] = data

			if finalChunk == 1 {
				packet.lastChunk = chunk
			}

			if packet.lastChunk != 0 {
				for i := 0; i <= int(packet.lastChunk); i++ {
					fragment, has := packet.data[uint8(i)]
					if has {
						result = append(result, fragment...)
					} else {
						result = nil
						break
					}
				}
			}
		}
	}

	if len(result) != 0 {
		delete(s.fragmentedPackets, sequenceNumber)
	}

	return result
}

func rgbaToYCbCr(img image.Image) *image.YCbCr {
	bounds := img.Bounds()
	ycbcr := image.NewYCbCr(bounds, image.YCbCrSubsampleRatio420)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			yy, cb, cr := color.RGBToYCbCr(r8, g8, b8)
			ycbcr.Y[ycbcr.YOffset(x, y)] = yy
			ycbcr.Cb[ycbcr.COffset(x, y)] = cb
			ycbcr.Cr[ycbcr.COffset(x, y)] = cr
		}
	}
	return ycbcr
}

func generateFrame(t float64) *image.YCbCr {
	dc := gg.NewContext(width, height)
	dc.SetRGB(0, 0, 0)
	dc.Clear()

	x := width/2 + int(200*math.Sin(t*2))
	y := height/2 + int(100*math.Cos(t*3))
	dc.SetColor(color.RGBA{255, 50, 50, 255})
	dc.DrawCircle(float64(x), float64(y), 40)
	dc.Fill()

	rgba := dc.Image()
	return rgbaToYCbCr(rgba)
}

func writeRawYUV(w io.Writer, img *image.YCbCr) error {
	if img.SubsampleRatio != image.YCbCrSubsampleRatio420 {
		return fmt.Errorf("invalid format")
	}
	_, err := w.Write(img.Y)
	if err != nil {
		return err
	}
	_, err = w.Write(img.Cb)
	if err != nil {
		return err
	}
	_, err = w.Write(img.Cr)
	return err
}

func (c *state) startVideoStream() {
	videoTrack, err := webrtc.NewTrackLocalStaticSample(
		webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8}, "video", "pion",
	)
	if err != nil {
		log.Fatalf("Failed to create video track: %v", err)
	}

	rtpSender, err := c.peerConnection.AddTrack(videoTrack)
	if err != nil {
		log.Fatalf("Failed to add video track: %v", err)
	}

	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				log.Printf("RTCP read error: %v", rtcpErr)
				return
			}
		}
	}()

	cmd := exec.Command("ffmpeg",
		"-f", "rawvideo",
		"-pix_fmt", "yuv420p",
		"-s", fmt.Sprintf("%dx%d", width, height),
		"-framerate", fmt.Sprint(fps),
		"-i", "pipe:0",
		"-c:v", "libvpx",
		"-b:v", "2M",
		"-deadline", "realtime",
		"-cpu-used", "8",
		"-threads", "4",
		"-error-resilient", "1",
		"-f", "ivf",
		"pipe:1",
	)

	ffmpegIn, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("Failed to get ffmpeg stdin: %v", err)
	}
	ffmpegOut, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to get ffmpeg stdout: %v", err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start ffmpeg: %v", err)
	}

	go func() {
		<-c.iceConnectedCtx.Done()
		t0 := time.Now()
		ticker := time.NewTicker(time.Second / fps)
		defer ticker.Stop()

		for range ticker.C {
			t := time.Since(t0).Seconds()
			yuv := generateFrame(t)
			if err := writeRawYUV(ffmpegIn, yuv); err != nil {
				log.Printf("Failed to write YUV to ffmpeg: %v", err)
				return
			}
		}
	}()

	go func() {
		<-c.iceConnectedCtx.Done()

		const ivfHeaderSize = 32
		const frameHeaderSize = 12

		var remainingData []byte
		var chunkNumber uint8
		var newFrame []byte

		bufReader := bufio.NewReader(ffmpegOut)

		_, err := io.CopyN(io.Discard, bufReader, ivfHeaderSize)
		if err != nil {
			log.Printf("Failed to skip IVF header: %v", err)
			return
		}

		for {
			header := make([]byte, frameHeaderSize)
			_, err := io.ReadFull(bufReader, header)
			if err != nil {
				log.Printf("Error reading IVF frame header: %v", err)
				return
			}

			frameSize := binary.LittleEndian.Uint32(header[0:4])

			frame := make([]byte, frameSize)
			_, err = io.ReadFull(bufReader, frame)
			if err != nil {
				log.Printf("Error reading IVF frame: %v", err)
				return
			}

			if proxyType == WEB {
				newFrame, remainingData, chunkNumber = c.encapsulateWeb(remainingData, frame, chunkNumber, c.conn)
			} else {
				newFrame, remainingData, chunkNumber = c.encapsulate(remainingData, frame, chunkNumber, c.conn)
			}

			if len(newFrame) > 0 {
				err = videoTrack.WriteSample(media.Sample{
					Data:     newFrame,
					Duration: time.Second / fps,
				})
				if err != nil {
					log.Printf("Error writing sample: %v", err)
				}
			}
		}
	}()
}

func (c *state) writePing() {
	ticker := time.NewTicker(PINGDEADLINE)
	for {
		select {
		case <-c.doneSignaling:
			c.connSignaling.Close()
			return
		case <-ticker.C:
			c.signalingConnMux.Lock()
			c.connSignaling.WriteMessage(websocket.BinaryMessage, []byte("ping"))
			c.signalingConnMux.Unlock()
		}
	}
}

func (c *state) setupConnection() {
	url := fmt.Sprintf("%s", configs["brokerAddr"])
	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			MinVersion:         tls.VersionTLS13,
			InsecureSkipVerify: true,
		},
	}

	conn_signaling, _, err := dialer.Dial(url, nil)
	if err != nil {
		panic(err)
	}

	data := []byte{CLIENT}
	data = append(data, []byte(fmt.Sprintf("%s", configs["bridgeaddr"]))...)

	c.signalingConnMux.Lock()
	conn_signaling.WriteMessage(websocket.BinaryMessage, data)
	c.signalingConnMux.Unlock()

	c.connSignaling = conn_signaling

	go c.writePing()

	// m := &webrtc.MediaEngine{}
	// //videoRTCPFeedback := []webrtc.RTCPFeedback{{Type: "goog-remb", Parameter: ""}, {Type: "ccm", Parameter: "fir"}, {Type: "nack", Parameter: ""}, {Type: "nack", Parameter: "pli"}}
	// // Setup the codecs you want to use.
	// // We'll use a VP8 and Opus but you can also define your own
	// if err := m.RegisterCodec(webrtc.RTPCodecParameters{
	// 	RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8, ClockRate: 90000, Channels: 0, SDPFmtpLine: "", RTCPFeedback: nil},
	// 	PayloadType:        96,
	// }, webrtc.RTPCodecTypeVideo); err != nil {
	// 	panic(err)
	// }
	// if err := m.RegisterCodec(webrtc.RTPCodecParameters{
	// 	RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus, ClockRate: 48000, Channels: 2, SDPFmtpLine: "minptime=10;useinbandfec=1", RTCPFeedback: nil},
	// 	PayloadType:        111,
	// }, webrtc.RTPCodecTypeAudio); err != nil {
	// 	panic(err)
	// }

	// // Create a InterceptorRegistry. This is the user configurable RTP/RTCP Pipeline.
	// // This provides NACKs, RTCP Reports and other features. If you use `webrtc.NewPeerConnection`
	// // this is enabled by default. If you are manually managing You MUST create a InterceptorRegistry
	// // for each PeerConnection.
	// i := &interceptor.Registry{}

	// // Use the default set of Interceptors
	// if err = webrtc.RegisterDefaultInterceptors(m, i); err != nil {
	// 	panic(err)
	// }

	// // Create the API object with the MediaEngine
	// api := webrtc.NewAPI(webrtc.WithMediaEngine(m), webrtc.WithInterceptorRegistry(i))

	// i := &interceptor.Registry{}
	// m := &webrtc.MediaEngine{}
	// if err := m.RegisterDefaultCodecs(); err != nil {
	// 	panic(err)
	// }

	// Create a Congestion Controller. This analyzes inbound and outbound data and provides
	// suggestions on how much we should be sending.
	//
	// Passing `nil` means we use the default Estimation Algorithm which is Google Congestion Control.
	// You can use the other ones that Pion provides, or write your own!
	// congestionController, err := cc.NewInterceptor(func() (cc.BandwidthEstimator, error) {
	// 	return gcc.NewSendSideBWE(gcc.SendSideBWEInitialBitrate(500_000))
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// estimatorChan := make(chan cc.BandwidthEstimator, 1)
	// congestionController.OnNewPeerConnection(func(id string, estimator cc.BandwidthEstimator) { //nolint: revive
	// 	estimatorChan <- estimator
	// })

	// i.Add(congestionController)
	// if err = webrtc.ConfigureTWCCHeaderExtensionSender(m, i); err != nil {
	// 	panic(err)
	// }

	// if err = webrtc.RegisterDefaultInterceptors(m, i); err != nil {
	// 	panic(err)
	// }

	// api := webrtc.NewAPI(webrtc.WithMediaEngine(m), webrtc.WithInterceptorRegistry(i))

	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				//URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	// estimator := <-estimatorChan

	// go func() {
	// 	for {
	// 		targetBitrate := estimator.GetTargetBitrate()
	// 		fmt.Println(targetBitrate)
	// 		time.Sleep(10 * time.Second)
	// 	}

	// }()

	c.peerConnection = peerConnection

	messageType, payload, err := c.connSignaling.ReadMessage()
	if err != nil {
		return
	}

	if messageType == websocket.BinaryMessage || messageType == websocket.TextMessage {
		if payload[0] == WEB {
			proxyType = WEB
		} else if payload[0] == PROXY {
			proxyType = PROXY
		} else {
			panic("Invalid type of proxy.")
		}
	}

	c.peerConnection.OnICECandidate(func(cand *webrtc.ICECandidate) {
		if cand == nil {
			return
		}

		candBytes, err := json.Marshal(cand.ToJSON())
		if err != nil {
			panic(err)
		}

		c.signalingConnMux.Lock()
		c.connSignaling.WriteMessage(websocket.BinaryMessage, candBytes)
		c.signalingConnMux.Unlock()
	})

	c.peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		if s == webrtc.PeerConnectionStateConnected {

			c.iceConnectedCtxCancel()

			if !c.doneSingalingBool {
				close(c.doneSignaling)
			}

			c.doneSingalingBool = true
			fmt.Println("connected")

			fileName := "/tmp/signal_file"
			_, err := os.Create(fileName)
			if err != nil {
				fmt.Println("Error creating signal file:", err)
				return
			}
		}

		if s == webrtc.PeerConnectionStateFailed {
			fmt.Println("Peer Connection has gone to failed exiting")
			close(c.closed)
		}

		if s == webrtc.PeerConnectionStateClosed {
			fmt.Println("Peer Connection has gone to closed exiting")
		}
	})

	c.peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		codec := track.Codec()
		if strings.EqualFold(codec.MimeType, webrtc.MimeTypeVP8) {
			for {
				packet, _, err := track.ReadRTP()
				if err != nil {
					return
				}

				if proxyType == WEB {
					c.desencapsulateWeb(packet.Payload, c.conn)
				} else {
					c.decapsulate(packet.Payload, c.conn)
				}
			}
		} else {
			for {
				_, _, err := track.ReadRTP()
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	})
}

// Acabar esta função quando acaba a cena
func (c *state) handleWSMessages() {
	defer func() {
		c.connSignaling.Close()
	}()

	for {
		messageType, payload, err := c.connSignaling.ReadMessage()
		if err != nil {
			return
		}

		if messageType == websocket.BinaryMessage || messageType == websocket.TextMessage {
			var (
				candidate webrtc.ICECandidateInit
				offer     webrtc.SessionDescription
			)

			switch {
			case json.Unmarshal(payload, &offer) == nil && offer.SDP != "":
				if sdpErr := c.peerConnection.SetRemoteDescription(offer); sdpErr != nil {
					panic(sdpErr)
				}

				if len(c.pendingCandidate) > 0 {
					for _, cand := range c.pendingCandidate {
						if candidateErr := c.peerConnection.AddICECandidate(*cand); candidateErr != nil {
							panic(candidateErr)
						}
					}
				}

				answer, err := c.peerConnection.CreateAnswer(nil)
				if err != nil {
					panic(err)
				}

				payload, err := json.Marshal(answer)
				if err != nil {
					panic(err)
				}

				c.signalingConnMux.Lock()
				c.connSignaling.WriteMessage(websocket.BinaryMessage, payload)
				c.signalingConnMux.Unlock()

				err = c.peerConnection.SetLocalDescription(answer)
				if err != nil {
					panic(err)
				}

			case json.Unmarshal(payload, &candidate) == nil && candidate.Candidate != "":
				if c.peerConnection.RemoteDescription() == nil {
					c.pendingCandidate = append(c.pendingCandidate, &candidate)
				} else {
					if candidateErr := c.peerConnection.AddICECandidate(candidate); candidateErr != nil {
						panic(candidateErr)
					}
				}
			default:
				panic("Unknown message")
			}
		}
	}
}

// Config file
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

// TURBO TUNNEL
type clientID []byte

func newClientID() clientID {
	buf := make([]byte, 8)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return buf
}

func (addr clientID) Network() string {
	return "clientid"
}

func (addr clientID) String() string {
	return hex.EncodeToString(addr)
}

type rOConn struct {
	idClient   []byte
	remoteAddr net.Addr
	localAddr  net.Addr
	closed     chan struct{}
	recvQueue  chan []byte
	sendQueue  chan []byte
}

func newrOConn() *rOConn {
	id := newClientID()
	conn := &rOConn{
		idClient:   id,
		remoteAddr: id,
		localAddr:  id,
		closed:     make(chan struct{}),
		recvQueue:  make(chan []byte, 100),
		sendQueue:  make(chan []byte, 100),
	}

	go conn.dialAndExchange()

	return conn
}

func (conn *rOConn) dialAndExchange() {
	for {
		iceConnectedCtx, iceConnectedCtxCancel := context.WithCancel(context.Background())
		state := state{
			connSignaling:         nil,
			peerConnection:        nil,
			conn:                  conn,
			sequenceNumer:         0,
			pictureIDTwoBytes:     0,
			fragmentedPackets:     make(map[uint32]*FragmentedPacket),
			doneSignaling:         make(chan int),
			doneSingalingBool:     false,
			iceConnectedCtx:       iceConnectedCtx,
			iceConnectedCtxCancel: iceConnectedCtxCancel,
			pendingCandidate:      make([]*webrtc.ICECandidateInit, 0),

			closed: make(chan int),
		}

		state.setupConnection()
		state.startVideoStream()

		go state.handleWSMessages()

		<-state.closed
	}
}

func (conn *rOConn) QueueIncoming(p []byte) {
	conn.recvQueue <- p
}

func (conn *rOConn) OutgoingQueue() []byte {
	select {
	case packet := <-conn.sendQueue:
		return packet
	default:
		return nil
	}
}

func (conn *rOConn) lenOutgoingQueue() int {
	return len(conn.sendQueue)
}

func (conn *rOConn) ReadFrom(p []byte) (int, net.Addr, error) {
	select {
	case data := <-conn.recvQueue:
		return copy(p, data), conn.RemoteAddr(), nil
	case <-conn.closed:
		return 0, nil, &net.OpError{Op: "read", Net: conn.RemoteAddr().Network(), Source: conn.LocalAddr(), Addr: conn.RemoteAddr(), Err: errors.New("closed conn")}
	}
}

func (conn *rOConn) WriteTo(p []byte, addr net.Addr) (int, error) {
	c := make([]byte, len(p)+8)

	copy(c[:8], conn.idClient)
	copy(c[8:], p)

	select {
	case <-conn.closed:
		return 0, &net.OpError{Op: "write", Net: conn.RemoteAddr().Network(), Source: conn.LocalAddr(), Addr: conn.RemoteAddr(), Err: errors.New("closed conn")}
	case conn.sendQueue <- c:
	default:
	}

	return len(c), nil
}

func (conn *rOConn) Close() error {
	select {

	case <-conn.closed:
		return &net.OpError{Op: "close", Net: conn.LocalAddr().Network(), Addr: conn.LocalAddr(), Err: errors.New("closed conn")}
	default:
		close(conn.closed)
		return nil
	}
}

func (conn *rOConn) LocalAddr() net.Addr {
	return conn.localAddr
}

func (conn *rOConn) RemoteAddr() net.Addr {
	return conn.remoteAddr
}

func (conn *rOConn) SetDeadline(t time.Time) error {
	return errors.New("not implemented")
}

func (conn *rOConn) SetReadDeadline(t time.Time) error {
	return errors.New("not implemented")
}

func (conn *rOConn) SetWriteDeadline(t time.Time) error {
	return errors.New("not implemented")
}

func handleLocalConn(conn *net.TCPConn, sess *smux.Session) error {
	stream, err := sess.OpenStream()
	if err != nil {
		return err
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

func acceptLocalConns(ln *net.TCPListener, sess *smux.Session) error {
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			if err, ok := err.(*net.OpError); ok && err.Temporary() {
				log.Printf("temporary error in ln.Accept: %v", err)
				continue
			}
			return err
		}

		go func() {
			defer conn.Close()
			err := handleLocalConn(conn, sess)
			if err != nil {
				log.Printf("error in handleLocalConn: %v", err)
			}
		}()
	}
}

func run() error {
	key := pbkdf2.Key([]byte("demo pass"), []byte("demo salt"), 1024, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)

	listenAddr := fmt.Sprintf("%s", configs["kcpListener"])
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()

	conn := newrOConn()
	defer conn.Close()

	kcpConn, err := kcp.NewConn2(nil, block, 0, 0, conn)
	if err != nil {
		return err
	}
	defer kcpConn.Close()

	kcpConn.SetNoDelay(1, 10, 2, 1)
	sess, err := smux.Client(kcpConn, &smux.Config{
		Version:           1,
		KeepAliveInterval: 10 * time.Second,
		KeepAliveTimeout:  100 * time.Second,
		MaxFrameSize:      32768,
		MaxReceiveBuffer:  4194304,
		MaxStreamBuffer:   65536})

	if err != nil {
		return err
	}
	defer sess.Close()

	return acceptLocalConns(ln.(*net.TCPListener), sess)
}

func main() {
	proxyType = 0
	readConfig()
	run()
}

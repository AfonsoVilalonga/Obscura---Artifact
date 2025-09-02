from netfilterqueue import NetfilterQueue
from scapy.layers.inet import IP, UDP
from random import random

FRAME_DROP_RATE = 0.25  # Probability to drop video frame
TARGET_BYTES = b'WebRTC'
ALT_TARGET_BYTES = b'webrtc'
packet_info = {}

def should_drop_frame():
    return random() < FRAME_DROP_RATE

def process_packet(packet):
    raw = packet.get_payload()
    pkt = IP(raw)

    if not pkt.haslayer(UDP):
        packet.accept()
        return

    udp = pkt[UDP]
    payload = bytes(udp.payload)

    if not payload:
        packet.accept()
        return

    key = f"{pkt.dst}:{udp.dport}"

    if payload[0] == 0x16:  # Likely DTLS packet
        if TARGET_BYTES in payload or ALT_TARGET_BYTES in payload:
            if key not in packet_info:
                packet_info[key] = {
                    "target_ip": pkt.dst,
                    "target_port": udp.dport,
                }
    elif key in packet_info and len(payload) > 1:
        marker_bit = (payload[1] & 0x80) >> 7
        payload_type = payload[1] & 0x7F

        if payload_type > 113 and marker_bit == 1:
            if should_drop_frame():
                packet.drop()
                return

    packet.accept()

def main():
    nfqueue = NetfilterQueue()
    try:
        nfqueue.bind(0, process_packet)
        print("NFQUEUE packet handler running...")
        nfqueue.run()
    except KeyboardInterrupt:
        print("Stopping...")
    finally:
        nfqueue.unbind()

if __name__ == "__main__":
    main()

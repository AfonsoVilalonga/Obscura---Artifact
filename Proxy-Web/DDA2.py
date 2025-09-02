from netfilterqueue import NetfilterQueue
from scapy.layers.inet import IP, UDP
from random import random

FRAME_DROP_RATE = 0.25  # Probability to drop a video frame
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

    # Use SSRC + ports as key (stream identifier)
    key = (pkt.src, pkt.dst, udp.sport, udp.dport)
    print(key)
    # Extract RTP marker bit and payload type (assume RTP over UDP)
    if len(payload) > 2:
        marker_bit = (payload[1] & 0x80) >> 7
        payload_type = payload[1] & 0x7F

        # Drop only video payloads (e.g., VP8/VP9)
        if payload_type > 113:  # Adjust based on actual codec
            if key not in packet_info:
                # Initialize stream tracking
                packet_info[key] = {
                    "drop_frame": False
                }

            stream = packet_info[key]

            # If it's the first packet of a frame, decide if we're dropping it
            if marker_bit == 0 and not stream["drop_frame"]:
                stream["drop_frame"] = should_drop_frame()

            # If marker bit is 1, it's the last packet in the frame
            if marker_bit == 1:
                drop_now = stream["drop_frame"]
                stream["drop_frame"] = False  # Reset for next frame
                if drop_now:
                    packet.drop()
                    return
            else:
                if stream["drop_frame"]:
                    packet.drop()
                    return

    # Default: accept
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

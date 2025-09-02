#!/bin/bash
cleanup() {
    echo "Cleaning up..."

    sudo iptables -F

    exit 0
}

trap cleanup SIGINT


sudo python3 DDA.py &
sleep 5
sudo iptables -I FORWARD -j NFQUEUE --queue-num 0







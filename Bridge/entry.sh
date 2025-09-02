#!/bin/bash
NETWORK_P=enp0s8
NETWORK_H=enp0s9

cleanup() {
    echo "Cleaning up..."

    sudo tc qdisc del dev ${NETWORK_H} root 2>/dev/null
    sudo tc qdisc del dev ${NETWORK_P} root 2>/dev/null

    exit 0
}

trap cleanup SIGINT

sudo tc qdisc add dev ${NETWORK_P} root netem delay 25ms 
sudo tc qdisc add dev ${NETWORK_H} root netem delay 7ms

go build .
./Bridge

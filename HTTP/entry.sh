#!/bin/bash

NETWORK_B=enp0s8

cleanup() {
    echo "Cleaning up..."

    sudo tc qdisc del dev ${NETWORK_B} root 2>/dev/null

    exit 0
}

trap cleanup SIGINT

sudo tc qdisc add dev ${NETWORK_B} root netem delay 8ms


go build .
./http

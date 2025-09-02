#!/bin/bash

TEST_N="2"
NETWORK_P=enp0s9

cleanup() {
    echo "Cleaning up..."

    sudo tc qdisc del dev ${NETWORK_P} root 2>/dev/null

    if [ -n "$CLIENT_PID" ]; then
        kill "$CLIENT_PID" 2>/dev/null
    fi

    sudo iptables -F

    exit 0
}

trap cleanup SIGINT

if [ "$TEST_N" = "1" ]; then
    sudo tc qdisc add dev ${NETWORK_P} root netem delay 7ms
elif [ "$TEST_N" = "2" ]; then
    sudo tc qdisc add dev ${NETWORK_P} root netem delay 25ms 
elif [ "$TEST_N" = "3" ]; then
    sudo tc qdisc add dev ${NETWORK_P} root netem delay 50ms
elif [ "$TEST_N" = "4" ]; then
    sudo tc qdisc add dev ${NETWORK_P} root netem delay 25ms loss 2%
elif [ "$TEST_N" = "5" ]; then
    sudo tc qdisc add dev ${NETWORK_P} root netem delay 25ms loss 5%
elif [ "$TEST_N" = "6" ]; then
    sudo tc qdisc add dev ${NETWORK_P} root netem delay 25ms loss 10%
elif [ "$TEST_N" = "7" ]; then
    sudo tc qdisc add dev ${NETWORK_P} root handle 1: htb default 12
    sudo tc class add dev ${NETWORK_P} parent 1: classid 1:12 htb rate 250kbit ceil 250kbit
    sudo tc qdisc add dev ${NETWORK_P} parent 1:12 netem delay 25ms
elif [ "$TEST_N" = "8" ]; then
    sudo tc qdisc add dev ${NETWORK_P} root handle 1: htb default 12
    sudo tc class add dev ${NETWORK_P} parent 1: classid 1:12 htb rate 750kbit ceil 750kbit
    sudo tc qdisc add dev ${NETWORK_P} parent 1:12 netem delay 25ms
elif [ "$TEST_N" = "9" ]; then
    sudo tc qdisc add dev ${NETWORK_P} root handle 1: htb default 12
    sudo tc class add dev ${NETWORK_P} parent 1: classid 1:12 htb rate 1500kbit ceil 1500kbit
    sudo tc qdisc add dev ${NETWORK_P} parent 1:12 netem delay 25ms
fi

go build .
./Client &
CLIENT_PID=$!

while [ ! -f /tmp/signal_file ]; do
  sleep 1  
done

./t.sh > t.txt
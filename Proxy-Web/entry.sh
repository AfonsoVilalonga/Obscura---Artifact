#!/bin/bash

TEST_N="2"
NETWORK_C=enp0s10
NETWORK_B=enp0s9
BROWSER_N=${BROWSER_NAME:-f}
DDA="no"

sudo ip link set enp0s10 up
sudo dhclient -v enp0s10

cleanup() {
    echo "Cleaning up..."

    sudo tc qdisc del dev ${NETWORK_C} root 2>/dev/null
    sudo tc qdisc del dev ${NETWORK_B} root 2>/dev/null

    if [ -n "$NODE_PID" ]; then
        kill "$NODE_PID" 2>/dev/null
    fi

    pkill -f firefoxP.py
    pkill -f chromeP.py

    sudo iptables -F

    exit 0
}


trap cleanup SIGINT


if [ "$TEST_N" = "1" ]; then
    sudo tc qdisc add dev ${NETWORK_C} root netem delay 8ms

    sudo tc qdisc add dev ${NETWORK_B} root netem delay 25ms
elif [ "$TEST_N" = "2" ]; then
    sudo tc qdisc add dev ${NETWORK_C} root netem delay 25ms

    sudo tc qdisc add dev ${NETWORK_B} root netem delay 25ms
elif [ "$TEST_N" = "3" ]; then
    sudo tc qdisc add dev ${NETWORK_C} root netem delay 50ms

    sudo tc qdisc add dev ${NETWORK_B} root netem delay 25ms
elif [ "$TEST_N" = "4" ]; then
    sudo tc qdisc add dev ${NETWORK_C} root netem delay 25ms loss 2%

    sudo tc qdisc add dev ${NETWORK_B} root netem delay 25ms
elif [ "$TEST_N" = "5" ]; then
    sudo tc qdisc add dev ${NETWORK_C} root netem delay 25ms loss 5%

    sudo tc qdisc add dev ${NETWORK_B} root netem delay 25ms
elif [ "$TEST_N" = "6" ]; then
    sudo tc qdisc add dev ${NETWORK_C} root netem delay 25ms loss 10%

    sudo tc qdisc add dev ${NETWORK_B} root netem delay 25ms
elif [ "$TEST_N" = "7" ]; then
    sudo tc qdisc add dev ${NETWORK_C} root handle 1: htb default 12
    sudo tc class add dev ${NETWORK_C} parent 1: classid 1:12 htb rate 250kbit ceil 250kbit
    sudo tc qdisc add dev ${NETWORK_C} parent 1:12 netem delay 25ms

    sudo tc qdisc add dev ${NETWORK_B} root netem delay 25ms
elif [ "$TEST_N" = "8" ]; then
    sudo tc qdisc add dev ${NETWORK_C} root handle 1: htb default 12
    sudo tc class add dev ${NETWORK_C} parent 1: classid 1:12 htb rate 750kbit ceil 750kbit
    sudo tc qdisc add dev ${NETWORK_C} parent 1:12 netem delay 25ms

    sudo tc qdisc add dev ${NETWORK_B} root netem delay 25ms
elif [ "$TEST_N" = "9" ]; then
    sudo tc qdisc add dev ${NETWORK_C} root handle 1: htb default 12
    sudo tc class add dev ${NETWORK_C} parent 1: classid 1:12 htb rate 1500kbit ceil 1500kbit
    sudo tc qdisc add dev ${NETWORK_C} parent 1:12 netem delay 25ms

    sudo tc qdisc add dev ${NETWORK_B} root netem delay 25ms
fi

if [ "$DDA" = "case_b" ]; then
    sudo python3 DDA.py &
    sleep 5
    sudo iptables -I OUTPUT -j NFQUEUE --queue-num 0
fi

if [ "$DDA" = "case_w" ]; then
    sudo python3 DDA2.py &
    sleep 5
    sudo iptables -I OUTPUT -j NFQUEUE --queue-num 0
fi

node index.js &
NODE_PID=$!
sleep 5

cd Selenium

if [ "$BROWSER_N" = "f" ]; then
    python3 firefoxP.py
else
    python3 chromeP.py
fi



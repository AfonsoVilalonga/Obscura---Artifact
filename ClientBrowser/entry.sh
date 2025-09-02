#!/bin/bash

TEST_N="2"
NETWORK_N=enp0s9
BROWSER_N=${BROWSER_NAME:-c}
DDA="T"

sudo ip link set enp0s9 up
sudo dhclient -v enp0s9

cleanup() {
    echo "Cleaning up..."

    sudo tc qdisc del dev ${NETWORK_N} root 2>/dev/null

    if [ -n "$NODE_PID" ]; then
        kill "$NODE_PID" 2>/dev/null
    fi

    if [ -n "$CLIENT_PID" ]; then
        kill "$CLIENT_PID" 2>/dev/null
    fi

    pkill -f firefoxC.py
    pkill -f chromeC.py

    sudo iptables -F

    exit 0
}

trap cleanup SIGINT


if [ "$TEST_N" = "1" ]; then
    sudo tc qdisc add dev ${NETWORK_N} root netem delay 7ms
elif [ "$TEST_N" = "2" ]; then
    sudo tc qdisc add dev ${NETWORK_N} root netem delay 25ms
elif [ "$TEST_N" = "3" ]; then
    sudo tc qdisc add dev ${NETWORK_N} root netem delay 50ms
elif [ "$TEST_N" = "4" ]; then
    sudo tc qdisc add dev ${NETWORK_N} root netem delay 25ms loss 2%
elif [ "$TEST_N" = "5" ]; then
    sudo tc qdisc add dev ${NETWORK_N} root netem delay 25ms loss 5%
elif [ "$TEST_N" = "6" ]; then
    sudo tc qdisc add dev ${NETWORK_N} root netem delay 25ms loss 10%
elif [ "$TEST_N" = "7" ]; then
    sudo tc qdisc add dev ${NETWORK_N} root handle 1: htb default 12
    sudo tc class add dev ${NETWORK_N} parent 1: classid 1:12 htb rate 250kbit ceil 250kbit
    sudo tc qdisc add dev ${NETWORK_N} parent 1:12 netem delay 25ms
elif [ "$TEST_N" = "8" ]; then
    sudo tc qdisc add dev ${NETWORK_N} root handle 1: htb default 12
    sudo tc class add dev ${NETWORK_N} parent 1: classid 1:12 htb rate 750kbit ceil 750kbit
    sudo tc qdisc add dev ${NETWORK_N} parent 1:12 netem delay 25ms
elif [ "$TEST_N" = "9" ]; then
    sudo tc qdisc add dev ${NETWORK_N} root handle 1: htb default 12
    sudo tc class add dev ${NETWORK_N} parent 1: classid 1:12 htb rate 1500kbit ceil 1500kbit
    sudo tc qdisc add dev ${NETWORK_N} parent 1:12 netem delay 25ms
fi

if [ "$DDA" = "T" ]; then
    sudo python3 DDA.py &
    sleep 5
    sudo iptables -I OUTPUT -j NFQUEUE --queue-num 0
fi

go build .
./clientbrowser &
CLIENT_PID=$!
sleep 5

cd Node-Server
node index.js &
NODE_PID=$!
sleep 5

cd ..
cd Selenium

if [ "$BROWSER_N" = "f" ]; then
    python3 firefoxC.py &
else
    python3 chromeC.py &
fi

while [ ! -f /tmp/signal_file ]; do
  sleep 1  
done

cd ..
./t.sh > t.txt



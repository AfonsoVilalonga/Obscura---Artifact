#!/bin/bash

# Example of continuous monitoring (commented out)
# while true; do
#     curl --socks5 127.0.0.1:10005 -w "Time to first byte: %{time_starttransfer} s\nTime until transfer began: %{time_pretransfer} s\nTotal time: %{time_total} s\nDownload speed: %{speed_download} bytes/sec\n" -o /dev/null http://192.168.30.2:8080/download
#     sleep 5
# done

# Loop 5 times and measure download speed
for i in $(seq 1 5); do
    curl --socks5 localhost:10005 -w "%{speed_download}\n" -o /dev/null -s http://192.168.50.50:8080/download
    sleep 2
done

# Loop 101 times to download files (commented out)
# for i in $(seq 1 101); do
#     curl --socks5 localhost:10005 http://192.168.50.50:8080/download -O
#     sleep 2
# done

echo "DONE"

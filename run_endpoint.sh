#!/bin/bash

# Set up the routing needed for the simulation.
/setup.sh

if [ "$ROLE" == "client" ]; then
    # Wait for the simulator to start up.
    /wait-for-it.sh sim:57832 -s -t 30
    echo "Starting QUIC client..."
    ./HTTP3-Discovery/quic-go-http3/client/main -h https://193.167.100.100:4433
elif [ "$ROLE" == "server" ]; then
    echo "Running QUIC server on  193.167.100.100:4433"
    tcpdump -i eth0 port 4433 -w packets.pcap &
    ./HTTP3-Discovery/quic-go-http3/server/main -h  193.167.100.100:4433 -c ./HTTP3-Discovery/certificates/cert.crt -k ./HTTP3-Discovery/certificates/priv.key -log ./logs/keylogs.txt -qlog ./logs
fi

echo "ending image"

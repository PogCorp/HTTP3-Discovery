# HTTP3-Discovery
The following shows how to use the [network simulator](https://github.com/quic-interop/quic-network-simulator) with projects involving the QUIC protocol.

# Setup
1. Clone the network simulator
```
git clone https://github.com/quic-interop/quic-network-simulator.git
```

2. Create a new folder inside the cloned repo
```
cd quic-network-simulator
mkdir discovery-http-quic-go
```

3. Copy the Dockerfile and bash script included in this branch
```
cp Dockerfile run_endpoint.sh /path/to/network/simulator/discovery-http-quic-go/
```

4. Inside the network simulator repo, build the docker compose project
```
CLIENT=discovery-http-quic-go  \
SERVER=discovery-http-quic-go \
docker-compose build
```

5. Run the docker compose selecting a network configuration as described in [scenarios](https://github.com/quic-interop/quic-network-simulator/tree/master/sim/scenarios/simple-p2p) and [tcp scenarios](https://github.com/quic-interop/quic-network-simulator/tree/master/sim/scenarios/tcp-cross-traffic)
```
CLIENT=discovery-http-quic-go \
SERVER=discovery-http-quic-go \
SCENARIO="simple-p2p --delay=15ms --bandwidth=10Mbps --queue=25" \
docker-compose up
```

# Debugging
To debug the runtime results of the simulation you will need to attach to the server docker container and follow this next steps:

1. Attach to the container
```
docker exec -it server /bin/bash
```

2. Stop the running _tcpdump_ process
```
pkill -SIGINT tcpdump
```

3. Open another terminal and find the id of the running container
```
docker ps
```
4. Open another terminal and transfer the log files from the simulation
```
docker cp <container-id>:packets.pcap .
docker cp <container-id>:/logs .
```

As a result you should have _qlog_ and _keylog_ files to analyse the simulation, as well as a pcap file that can be used in Wireshark (in conjunction to the _keylog_) to track all the packets transmitted during the test.

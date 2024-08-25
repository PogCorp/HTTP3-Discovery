FROM martenseemann/quic-network-simulator-endpoint:latest

RUN apt-get update
RUN apt-get install -y wget tar git vim python3

RUN wget https://dl.google.com/go/go1.22.5.linux-amd64.tar.gz && \
  tar xfz go1.22.5.linux-amd64.tar.gz && \
  rm go1.22.5.linux-amd64.tar.gz

ENV PATH="/go/bin:${PATH}"

# download and build your QUIC implementation
RUN git clone https://github.com/PogCorp/HTTP3-Discovery.git
WORKDIR /HTTP3-Discovery
RUN git checkout jet/debug-with-keylogs
RUN ls
RUN openssl req -newkey rsa:4096 \
            -x509 \
            -sha256 \
            -days 3650 \
            -nodes \
            -out cert.crt \
            -keyout priv.key \
            -subj "/C=BR/ST=São Paulo/L=São Paulo/O=PogCorp/OU=DEV/CN=pogcorp@gmail.com"
RUN mkdir certificates
RUN mv cert.crt priv.key ./certificates
WORKDIR /HTTP3-Discovery/quic-go-http3/server/
RUN go mod tidy
RUN go build ./cmd/main.go
WORKDIR /HTTP3-Discovery/quic-go-http3/client/
RUN go mod tidy
RUN go build ./cmd/main.go
WORKDIR /

# copy run script and run it
COPY run_endpoint.sh .
RUN chmod +x run_endpoint.sh
ENTRYPOINT [ "./run_endpoint.sh" ]

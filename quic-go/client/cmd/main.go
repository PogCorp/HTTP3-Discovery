package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/qlog"
)

func client(addr string) {
	log.Println("starting client")

	var message string

	for {
		_, err := fmt.Scanf("%s", &message)
		if err != nil {
			log.Println("failed to extract message")
			continue
		}
		if message == "stop" {
			log.Println("closing the program")
			return
		}

		log.Println(message)

		tlsConf := &tls.Config{
			InsecureSkipVerify: true,
			NextProtos:         []string{"echo"},
		}

		connUdp, err := net.ListenPacket("udp6", ":0")
		if err != nil {
			log.Println(err)
			return
		}

		// NOTE: example from https://github.com/golang/go/commit/645d4726f0f36c3aec9c864f47411a74c20ebc70
		udpAddr, err := net.ResolveUDPAddr("udp6", addr)
		if err != nil {
			log.Println(err)
			return
		}
		conn, err := quic.Dial(context.Background(), connUdp, udpAddr, tlsConf, &quic.Config{Tracer: qlog.DefaultTracer})
		//conn, err := quic.DialAddr(context.Background(), addr, tlsConf, &quic.Config{Tracer: qlog.DefaultTracer})
		if err != nil {
			log.Println("Big retard", err)
			return
		}

		stream, err := conn.OpenStreamSync(context.Background())
		if err != nil {
			log.Println(err)
		}

		fmt.Printf("Client: Sending '%s'\n", message)
		_, err = stream.Write([]byte(message))
		if err != nil {
			log.Println(err)
		}

		buf := make([]byte, len(message))
		_, err = stream.Read(buf)

		if err != nil {
			log.Println("failed to read message from server", err)
			continue
		}
		fmt.Printf("Client: Got '%s'\n", buf)

		stream.Close()
		err = conn.CloseWithError(0, "")
		if err != nil {
			log.Println("error while closing conection", err)
		}
	}
}

func main() {
	log.Println("starting echo client")
	addr := flag.String("h", "127.0.0.1:8080", "host:port")
	flag.Parse()

	if *addr == "" {
		log.Fatalln("certificate not provided")
	}

	client(*addr)
}

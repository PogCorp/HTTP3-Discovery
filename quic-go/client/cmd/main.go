package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"

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

		tlsConf := &tls.Config{
			InsecureSkipVerify: true,
		}
		conn, err := quic.DialAddr(context.Background(), addr, tlsConf, &quic.Config{Tracer: qlog.DefaultTracer})
		if err != nil {
			log.Println(err)
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

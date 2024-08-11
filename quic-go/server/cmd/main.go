package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/qlog"
)

func server(addr, cert, key, keylogPath string) {
	config, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		log.Println("bad tls configuration")
		return
	}

	w, err := os.OpenFile(keylogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Println("logWritter could not be opened")
		return
	}

	listener, err := quic.ListenAddr(addr, &tls.Config{Certificates: []tls.Certificate{config}, KeyLogWriter: w}, &quic.Config{Tracer: qlog.DefaultTracer})
	if err != nil {
		log.Println("got error while creating listener", err)
		return
	}
	defer listener.Close()

	log.Println("initializing server")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	go func() {
		for sig := range stop {
			log.Printf("Stoping QUIC Server now, got signal: %s\n", sig.String())
			cancelFunc()
		}
	}()

	for cancelCtx.Err() == nil {
		conn, err := listener.Accept(cancelCtx)
		if err != nil {
			log.Println("Accept failed", err)
			continue
		}

		go func() {
			for cancelCtx.Err() == nil {
				stream, err := conn.AcceptStream(cancelCtx)
				if err != nil {
					log.Println("AcceptStream received error", err)
					return
				}
				buffer := make([]byte, 256)
				_, err = stream.Read(buffer)
				if err != nil {
					log.Println("unable to read from stream")
				}

				message := string(buffer)
				id := stream.StreamID()
				log.Printf("[stream:%d]: received %s\n", int64(id), message)

				_, err = stream.Write([]byte(strings.ToUpper(message)))
				if err != nil {
					log.Println("server was unable to send response", err)
				}
				stream.Close()
			}
		}()
	}
}

func main() {

	log.Println("starting echo server")

	addr := flag.String("h", "127.0.0.1:8080", "host:port")
	cert := flag.String("c", "", "/path/to/certificate.crt")
	key := flag.String("k", "", "/path/to/private.key")
	logPath := flag.String("log", "", "path/to/logs")

	flag.Parse()

	if *cert == "" {
		log.Fatalln("certificate not provided")
	}

	if *key == "" {
		log.Fatalln("key not provided")
	}

	if *logPath == "" {
		log.Fatalln("keylog path not provided")
	}

	server(*addr, *cert, *key, *logPath)
}

package main

import (
	"crypto/tls"
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/quic-go/qlog"
)

func client(addr string) {

	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},    // set a TLS client config, if desired
		QUICConfig:      &quic.Config{Tracer: qlog.DefaultTracer}, // QUIC connection options
	}
	defer roundTripper.Close()
	client := &http.Client{
		Transport: roundTripper,
	}

	res, err := client.Get(addr)
	if err != nil {
		log.Println("client request failed", err)
		return
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("failed to read body")
		return
	}
	log.Printf("status: %d, respose: %s", res.StatusCode, string(body))
}

func main() {
	log.Println("starting http client")
	addr := flag.String("h", "127.0.0.1:8080", "host:port")
	flag.Parse()

	if *addr == "" {
		log.Fatalln("certificate not provided")
	}

	client(*addr)
}

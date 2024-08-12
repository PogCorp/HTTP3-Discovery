package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/quic-go/qlog"
)

func server(addr, cert, key, keylogPath string) {
	config, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		log.Println("unable to load certificate", err)
		return
	}
	w, err := os.OpenFile(keylogPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Println("logWritter could not be opened")
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tm := time.Now().Format(time.RFC1123)
		_, err := w.Write([]byte("The time is: " + tm))
		if err != nil {
			log.Println("unble to write response")
		}
		w.WriteHeader(200)
	})
	server := http3.Server{
		Handler:    mux,
		Addr:       addr,
		TLSConfig:  http3.ConfigureTLSConfig(&tls.Config{Certificates: []tls.Certificate{config}, KeyLogWriter: w}),
		QUICConfig: &quic.Config{Tracer: qlog.DefaultTracer},
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Println("unable to ListenAndServe", err)
	}
}

func main() {
	log.Println("starting http/3 server")

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

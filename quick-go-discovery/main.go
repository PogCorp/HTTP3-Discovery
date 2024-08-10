package main

import (
	"quick-go-discovery/client"
	"quick-go-discovery/server"
)

func main() {
	go func() {
		server.EchoServer()
	}()

	client.ClientMain()
}

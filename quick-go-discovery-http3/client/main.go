package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

type Person struct {
	Name string `json:"name"`
}

func main() {
	transport := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		QUICConfig: &quic.Config{
			KeepAlivePeriod: time.Minute * 30,
			MaxIdleTimeout:  time.Minute * 30,
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Minute * 30,
	}

	payload, err := json.Marshal(Person{
		Name: "Rafa",
	})
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	resp, err := client.Post("https://localhost:6121/demo/echo", "application/json", strings.NewReader(string(payload)))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		os.Exit(1)
	}

	fmt.Println("Response:", string(body))
}

package client

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/logging"
	"github.com/quic-go/quic-go/qlog"
)

const addr = "localhost:4242"

// func main() {
// 	if err := clientMain(); err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 	}
// }

func ClientMain() error {
	qlogFilename := fmt.Sprintf("client_%s.qlog", time.Now().Format("20060102_150405"))
	qlogFile, err := os.Create(filepath.Join(".", qlogFilename))
	if err != nil {
		return fmt.Errorf("failed to create qlog file: %v", err)
	}
	defer qlogFile.Close()

	config := &quic.Config{
		EnableDatagrams: true,
		Tracer: func(ctx context.Context, p logging.Perspective, ci quic.ConnectionID) *logging.ConnectionTracer {
			return qlog.NewConnectionTracer(qlogFile, p, ci)
		},
	}

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}

	conn, err := quic.DialAddr(context.Background(), addr, tlsConf, config)
	if err != nil {
		return err
	}
	defer conn.CloseWithError(0, "")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter message to send: ")
		message, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		message = message[:len(message)-1]

		if message == "END" {
			break
		}

		stream, err := conn.OpenStreamSync(context.Background())
		if err != nil {
			return err
		}

		fmt.Printf("Client: Sending '%s'\n", message)
		_, err = stream.Write([]byte(message))
		if err != nil {
			return err
		}

		buf := make([]byte, len(message))
		_, err = io.ReadFull(stream, buf)
		if err != nil {
			return err
		}
		fmt.Printf("Client: Got '%s'\n", buf)

		stream.Close()
	}

	return nil
}

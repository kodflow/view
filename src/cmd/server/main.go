package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/kodflow/view/internal/capture"
	"github.com/kodflow/view/internal/protocol"
	"github.com/kodflow/view/internal/stealth"
)

func main() {
	tcpPort := flag.Int("port", protocol.DefaultTCPPort, "TCP port for stealth listener")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	srv := &stealth.Server{
		OnClient: handleClient,
	}

	addr := fmt.Sprintf("0.0.0.0:%d", *tcpPort)
	log.Printf("listening on %s (stealth SSH)", addr)
	if err := srv.Listen(ctx, addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	log.Printf("client connected: %s", conn.RemoteAddr())

	for {
		msgType, _, err := protocol.ReadMessage(conn)
		if err != nil {
			log.Printf("client %s disconnected: %v", conn.RemoteAddr(), err)
			return
		}

		switch msgType {
		case protocol.MsgCapture:
			data, err := capture.CaptureScreen(context.Background())
			if err != nil {
				log.Printf("capture error: %v", err)
				protocol.WriteMessage(conn, protocol.MsgError, []byte(err.Error()))
				continue
			}
			if err := protocol.WriteMessage(conn, protocol.MsgFrame, data); err != nil {
				log.Printf("send error: %v", err)
				return
			}
			log.Printf("sent frame: %d bytes", len(data))

		default:
			log.Printf("unknown message type: 0x%02x", msgType)
		}
	}
}

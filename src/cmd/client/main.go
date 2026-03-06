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
	"time"

	"github.com/kodflow/view/internal/protocol"
	"github.com/kodflow/view/internal/web"
)

func main() {
	serverAddr := flag.String("server", "localhost", "server address (ip or ip:port)")
	serverPort := flag.Int("port", protocol.DefaultTCPPort, "server port (used if --server has no port)")
	interval := flag.Duration("interval", 500*time.Millisecond, "capture interval")
	noBrowser := flag.Bool("no-browser", false, "don't open browser automatically")
	flag.Parse()

	// Add port if not present
	if _, _, err := net.SplitHostPort(*serverAddr); err != nil {
		*serverAddr = fmt.Sprintf("%s:%d", *serverAddr, *serverPort)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Start web UI
	webServer := web.New()
	port, err := webServer.Start(ctx)
	if err != nil {
		log.Fatalf("web server: %v", err)
	}
	url := fmt.Sprintf("http://localhost:%d", port)
	log.Printf("web UI at %s", url)

	if !*noBrowser {
		web.OpenBrowser(url)
	}

	// Connect and stream, retry forever
	for ctx.Err() == nil {
		if err := streamFromServer(ctx, *serverAddr, webServer, *interval); err != nil {
			log.Printf("connection error: %v — reconnecting in 3s...", err)
			select {
			case <-time.After(3 * time.Second):
			case <-ctx.Done():
				return
			}
		}
	}
}

func streamFromServer(ctx context.Context, addr string, webServer *web.Server, interval time.Duration) error {
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	defer conn.Close()

	// Send magic bytes
	if _, err := conn.Write(protocol.MagicBytes[:]); err != nil {
		return fmt.Errorf("send magic: %w", err)
	}
	log.Printf("connected to %s", addr)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := protocol.WriteMessage(conn, protocol.MsgCapture, nil); err != nil {
				return fmt.Errorf("send capture request: %w", err)
			}

			msgType, data, err := protocol.ReadMessage(conn)
			if err != nil {
				return fmt.Errorf("read frame: %w", err)
			}

			switch msgType {
			case protocol.MsgFrame:
				webServer.PushFrame(data)
			case protocol.MsgError:
				log.Printf("server error: %s", string(data))
			default:
				log.Printf("unexpected message type: 0x%02x", msgType)
			}
		}
	}
}

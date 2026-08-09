package stealth

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/kodflow/view/internal/protocol"
)

// Server is a TCP server that mimics SSH for non-authenticated connections
// and serves the custom protocol for clients sending magic bytes.
type Server struct {
	OnClient func(conn net.Conn)
}

// Listen starts the stealth TCP server. It blocks until ctx is cancelled.
func (s *Server) Listen(ctx context.Context, addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	defer ln.Close()

	go func() {
		<-ctx.Done()
		ln.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	// Read first 8 bytes with a timeout
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buf := make([]byte, len(protocol.MagicBytes))
	n, err := io.ReadFull(conn, buf)

	if err != nil || n < len(protocol.MagicBytes) {
		// No data or incomplete — scanner/nmap: send SSH banner
		s.sendSSHBanner(conn)
		return
	}

	if !protocol.IsMagicBytes(buf) {
		// Wrong magic — send SSH banner
		s.sendSSHBanner(conn)
		return
	}

	// Valid client — reset deadline and hand off
	conn.SetReadDeadline(time.Time{})
	log.Printf("client connected: %s", conn.RemoteAddr())
	if s.OnClient != nil {
		s.OnClient(conn)
	}
}

func (s *Server) sendSSHBanner(conn net.Conn) {
	conn.Write([]byte(protocol.SSHBanner))
	// Keep connection open briefly like a real SSH server
	time.Sleep(3 * time.Second)
}

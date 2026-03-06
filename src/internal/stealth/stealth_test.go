package stealth

import (
	"context"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/kodflow/view/internal/protocol"
)

func startTestServer(t *testing.T, onClient func(net.Conn)) (string, context.CancelFunc) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	s := &Server{OnClient: onClient}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	ln.Close()

	go s.Listen(ctx, addr)
	time.Sleep(50 * time.Millisecond) // let server start

	return addr, cancel
}

func TestMagicBytesAccepted(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	clientConnected := false
	addr, cancel := startTestServer(t, func(conn net.Conn) {
		clientConnected = true
		wg.Done()
	})
	defer cancel()

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	conn.Write(protocol.MagicBytes[:])
	wg.Wait()

	if !clientConnected {
		t.Error("OnClient was not called")
	}
}

func TestNoMagicBytesSendsSSHBanner(t *testing.T) {
	addr, cancel := startTestServer(t, nil)
	defer cancel()

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// Send wrong bytes
	conn.Write([]byte("GET / HTTP/1.1\r\n"))

	buf := make([]byte, 256)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	got := string(buf[:n])
	if got != protocol.SSHBanner {
		t.Errorf("got %q, want SSH banner", got)
	}
}

func TestNmapLikeNoData(t *testing.T) {
	addr, cancel := startTestServer(t, nil)
	defer cancel()

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// Send nothing, like nmap NULL probe
	buf := make([]byte, 256)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	got := string(buf[:n])
	if got != protocol.SSHBanner {
		t.Errorf("got %q, want SSH banner", got)
	}
}

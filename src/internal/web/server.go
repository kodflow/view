package web

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"sync"

	"nhooyr.io/websocket"
)

//go:embed static/*
var staticFiles embed.FS

// Server serves the web UI and pushes frames to connected browsers via WebSocket.
type Server struct {
	mu      sync.RWMutex
	clients map[*wsClient]struct{}
}

type wsClient struct {
	conn *websocket.Conn
	ctx  context.Context
}

func New() *Server {
	return &Server{
		clients: make(map[*wsClient]struct{}),
	}
}

// Start starts the HTTP server on a random port and returns the port number.
func (s *Server) Start(ctx context.Context) (int, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, fmt.Errorf("listen: %w", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port

	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return 0, fmt.Errorf("sub fs: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(staticFS)))
	mux.HandleFunc("/ws", s.handleWS)

	srv := &http.Server{Handler: mux}
	go func() {
		<-ctx.Done()
		srv.Close()
	}()
	go srv.Serve(ln)

	return port, nil
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Printf("websocket accept: %v", err)
		return
	}

	client := &wsClient{conn: conn, ctx: r.Context()}
	s.mu.Lock()
	s.clients[client] = struct{}{}
	s.mu.Unlock()

	// Block until connection closes
	for {
		_, _, err := conn.Read(r.Context())
		if err != nil {
			break
		}
	}

	s.mu.Lock()
	delete(s.clients, client)
	s.mu.Unlock()
}

// PushFrame sends JPEG frame data to all connected WebSocket clients.
func (s *Server) PushFrame(data []byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for client := range s.clients {
		err := client.conn.Write(client.ctx, websocket.MessageBinary, data)
		if err != nil {
			client.conn.Close(websocket.StatusGoingAway, "write error")
		}
	}
}

// OpenBrowser opens the default browser to the given URL.
func OpenBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	cmd.Start()
}

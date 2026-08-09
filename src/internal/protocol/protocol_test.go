package protocol

import (
	"net"
	"testing"
)

func TestIsMagicBytes(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{"valid", MagicBytes[:], true},
		{"too short", []byte{0x56, 0x49}, false},
		{"wrong", []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, false},
		{"empty", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMagicBytes(tt.data); got != tt.want {
				t.Errorf("IsMagicBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteReadMessage(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	payload := []byte("hello world")
	errCh := make(chan error, 1)

	go func() {
		errCh <- WriteMessage(server, MsgFrame, payload)
	}()

	msgType, data, err := ReadMessage(client)
	if err != nil {
		t.Fatalf("ReadMessage: %v", err)
	}
	if err := <-errCh; err != nil {
		t.Fatalf("WriteMessage: %v", err)
	}
	if msgType != MsgFrame {
		t.Errorf("msgType = %d, want %d", msgType, MsgFrame)
	}
	if string(data) != string(payload) {
		t.Errorf("payload = %q, want %q", data, payload)
	}
}

func TestWriteReadEmptyPayload(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	go func() {
		WriteMessage(server, MsgCapture, nil)
	}()

	msgType, data, err := ReadMessage(client)
	if err != nil {
		t.Fatalf("ReadMessage: %v", err)
	}
	if msgType != MsgCapture {
		t.Errorf("msgType = %d, want %d", msgType, MsgCapture)
	}
	if len(data) != 0 {
		t.Errorf("payload len = %d, want 0", len(data))
	}
}

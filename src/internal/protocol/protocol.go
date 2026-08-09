package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

var MagicBytes = [8]byte{0x56, 0x49, 0x45, 0x57, 0x50, 0x52, 0x4F, 0x54} // "VIEWPROT"

var DiscoveryRequest = [4]byte{0x56, 0x49, 0x45, 0x57} // "VIEW"

const (
	SSHBanner      = "SSH-2.0-OpenSSH_9.7\r\n"
	DefaultTCPPort = 22

	MsgCapture byte = 0x01
	MsgFrame   byte = 0x02
	MsgError   byte = 0xFF

	MaxPayloadSize = 50 * 1024 * 1024 // 50MB max frame
)

func IsMagicBytes(data []byte) bool {
	if len(data) < len(MagicBytes) {
		return false
	}
	for i, b := range MagicBytes {
		if data[i] != b {
			return false
		}
	}
	return true
}

// WriteMessage writes [1 byte type][4 bytes length big-endian][payload]
func WriteMessage(conn net.Conn, msgType byte, payload []byte) error {
	header := make([]byte, 5)
	header[0] = msgType
	binary.BigEndian.PutUint32(header[1:5], uint32(len(payload)))
	if _, err := conn.Write(header); err != nil {
		return fmt.Errorf("write header: %w", err)
	}
	if len(payload) > 0 {
		if _, err := conn.Write(payload); err != nil {
			return fmt.Errorf("write payload: %w", err)
		}
	}
	return nil
}

// ReadMessage reads a complete message from the connection.
func ReadMessage(conn net.Conn) (msgType byte, payload []byte, err error) {
	header := make([]byte, 5)
	if _, err := io.ReadFull(conn, header); err != nil {
		return 0, nil, fmt.Errorf("read header: %w", err)
	}
	msgType = header[0]
	length := binary.BigEndian.Uint32(header[1:5])
	if length > MaxPayloadSize {
		return 0, nil, fmt.Errorf("payload too large: %d bytes", length)
	}
	if length == 0 {
		return msgType, nil, nil
	}
	payload = make([]byte, length)
	if _, err := io.ReadFull(conn, payload); err != nil {
		return 0, nil, fmt.Errorf("read payload: %w", err)
	}
	return msgType, payload, nil
}

package discovery

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/kodflow/view/internal/protocol"
)

// ListenForDiscovery listens for UDP broadcast discovery requests and responds
// with the server's TCP port. Runs until ctx is cancelled.
func ListenForDiscovery(ctx context.Context, udpPort, tcpPort int) error {
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("0.0.0.0:%d", udpPort))
	if err != nil {
		return fmt.Errorf("resolve udp addr: %w", err)
	}
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return fmt.Errorf("listen udp: %w", err)
	}
	defer conn.Close()

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	buf := make([]byte, 64)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			continue
		}
		if n < len(protocol.DiscoveryRequest) {
			continue
		}
		if buf[0] != protocol.DiscoveryRequest[0] ||
			buf[1] != protocol.DiscoveryRequest[1] ||
			buf[2] != protocol.DiscoveryRequest[2] ||
			buf[3] != protocol.DiscoveryRequest[3] {
			continue
		}

		resp := make([]byte, 6)
		copy(resp[:4], protocol.DiscoveryRequest[:])
		binary.BigEndian.PutUint16(resp[4:6], uint16(tcpPort))
		conn.WriteToUDP(resp, remoteAddr)
	}
}

// DiscoverServer sends UDP broadcast to find the server on the LAN.
// Returns the server IP and TCP port.
func DiscoverServer(ctx context.Context, udpPort int) (string, int, error) {
	broadcastAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", udpPort))
	if err != nil {
		return "", 0, fmt.Errorf("resolve broadcast: %w", err)
	}

	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		return "", 0, fmt.Errorf("listen udp: %w", err)
	}
	defer conn.Close()

	for attempt := 0; attempt < 3; attempt++ {
		if ctx.Err() != nil {
			return "", 0, ctx.Err()
		}

		if _, err := conn.WriteToUDP(protocol.DiscoveryRequest[:], broadcastAddr); err != nil {
			return "", 0, fmt.Errorf("send discovery: %w", err)
		}

		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		buf := make([]byte, 64)
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue // timeout, retry
		}
		if n < 6 {
			continue
		}
		if buf[0] != protocol.DiscoveryRequest[0] ||
			buf[1] != protocol.DiscoveryRequest[1] ||
			buf[2] != protocol.DiscoveryRequest[2] ||
			buf[3] != protocol.DiscoveryRequest[3] {
			continue
		}

		port := int(binary.BigEndian.Uint16(buf[4:6]))
		return remoteAddr.IP.String(), port, nil
	}

	return "", 0, fmt.Errorf("no server found after 3 attempts")
}

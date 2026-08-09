//go:build darwin

package capture

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

// CaptureScreen captures the macOS screen using the built-in screencapture CLI.
// Returns raw JPEG bytes.
func CaptureScreen(ctx context.Context) ([]byte, error) {
	tmpFile, err := os.CreateTemp("", "view-*.jpg")
	if err != nil {
		return nil, fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	cmd := exec.CommandContext(ctx, "/usr/sbin/screencapture", "-x", "-C", "-t", "jpg", tmpPath)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("screencapture: %w", err)
	}

	data, err := os.ReadFile(tmpPath)
	if err != nil {
		return nil, fmt.Errorf("read capture: %w", err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("capture produced empty file")
	}
	return data, nil
}

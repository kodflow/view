//go:build !darwin

package capture

import (
	"context"
	"fmt"
)

// CaptureScreen is a stub for non-macOS platforms.
// The server can only run on macOS.
func CaptureScreen(_ context.Context) ([]byte, error) {
	return nil, fmt.Errorf("screen capture is only supported on macOS")
}

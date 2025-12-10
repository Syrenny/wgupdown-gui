package wgupdown

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	pkexecPath   = "/usr/bin/pkexec"
	wgupdownPath = "/usr/local/bin/wgupdown"
)

// Up brings the WireGuard interface up using `wgupdown` utility
func Up(ctx context.Context, ifaceName string) error {
	cmd := exec.CommandContext(ctx, pkexecPath, wgupdownPath, "up", ifaceName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("wgupdown up failed: %w", err)
	}
	return nil
}

// Down brings the WireGuard interface down using `wgupdown` utility
func Down(ctx context.Context, ifaceName string) error {
	cmd := exec.CommandContext(ctx, pkexecPath, wgupdownPath, "down", ifaceName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("wgupdown down failed: %w", err)
	}
	return nil
}

// IsUp checks if the interface is up using `wgupdown` utility
func IsUp(ctx context.Context, ifaceName string) (bool, error) {
	cmd := exec.CommandContext(ctx, pkexecPath, wgupdownPath, "status", ifaceName)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return false, fmt.Errorf("wgupdown status failed: %w, stderr: %s", err, stderr.String())
	}

	statusStr := strings.TrimSpace(out.String())

	switch statusStr {
	case "up":
		return true, nil
	case "down":
		return false, nil
	default:
		return false, fmt.Errorf("unexpected output from wgupdown: %q", statusStr)
	}
}

package wgupdown

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"github.com/Syrenny/wgupdown-gui/pkg/wireguard"
)

const (
	sudoPath       = "/usr/bin/sudo"
	wgupdownPath   = "/usr/local/bin/wgupdown"
	nonInteractive = "-n"
)

func runHelper(ctx context.Context, action string, ifaceName string) error {
	cmd := exec.CommandContext(ctx, sudoPath, nonInteractive, wgupdownPath, action, ifaceName)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return fmt.Errorf("wgupdown %s %s failed: %w: %s", action, ifaceName, err, stderr.String())
		}

		return fmt.Errorf("wgupdown %s %s failed: %w", action, ifaceName, err)
	}

	return nil
}

func Up(ctx context.Context, ifaceName string) error {
	return runHelper(ctx, "up", ifaceName)
}

func Down(ctx context.Context, ifaceName string) error {
	return runHelper(ctx, "down", ifaceName)
}

func IsUp(ctx context.Context, ifaceName string) (bool, error) {
	return wireguard.IsUp(ctx, ifaceName)
}

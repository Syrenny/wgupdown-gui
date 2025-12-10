package wireguard

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	pkexecPath  = "/usr/bin/pkexec"
	wgquickPath = "/usr/bin/wg-quick"
)

// Up brings the WireGuard interface up using `wg-quick up`
func Up(ctx context.Context, ifaceName string) error {
	cmd := exec.CommandContext(ctx, pkexecPath, wgquickPath, "up", ifaceName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("wg-quick up failed: %w", err)
	}
	return nil
}

// Down brings the WireGuard interface down using `wg-quick down`
func Down(ctx context.Context, ifaceName string) error {
	cmd := exec.CommandContext(ctx, pkexecPath, wgquickPath, "down", ifaceName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("wg-quick down failed: %w", err)
	}
	return nil
}

// IsUp checks if the interface exists in /sys/class/net
func IsUp(ifaceName string) (bool, error) {
	ifacePath := filepath.Join("/sys/class/net", ifaceName)
	_, err := os.Stat(ifacePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("cannot stat interface: %w", err)
}

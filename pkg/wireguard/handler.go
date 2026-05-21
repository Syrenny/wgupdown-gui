package wireguard

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

var ifaceNamePattern = regexp.MustCompile(`^[A-Za-z0-9_.+=-]+$`)

func validateInterfaceName(ifaceName string) error {
	if !ifaceNamePattern.MatchString(ifaceName) {
		return fmt.Errorf("invalid interface name: %q", ifaceName)
	}

	return nil
}

func runWGQuick(ctx context.Context, action string, ifaceName string) error {
	if err := validateInterfaceName(ifaceName); err != nil {
		return err
	}

	wgQuickPath, err := exec.LookPath("wg-quick")
	if err != nil {
		return fmt.Errorf("wg-quick not found: %w", err)
	}

	cmd := exec.CommandContext(ctx, wgQuickPath, action, ifaceName)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("wg-quick %s %s failed: %w", action, ifaceName, err)
	}

	return nil
}

// Up brings the WireGuard interface up using `wg-quick up`.
func Up(ctx context.Context, ifaceName string) error {
	return runWGQuick(ctx, "up", ifaceName)
}

// Down brings the WireGuard interface down using `wg-quick down`.
func Down(ctx context.Context, ifaceName string) error {
	return runWGQuick(ctx, "down", ifaceName)
}

// IsUp checks if the interface exists in /sys/class/net.
func IsUp(_ context.Context, ifaceName string) (bool, error) {
	if err := validateInterfaceName(ifaceName); err != nil {
		return false, err
	}

	ifacePath := filepath.Join("/sys/class/net", ifaceName)
	_, err := os.Stat(ifacePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, fmt.Errorf("cannot stat interface %q: %w", ifaceName, err)
}

package wgupdown

import (
	"fmt"
	"os/exec"
)

// Up brings the VPN up
func Up() error {
	cmd := exec.Command("systemctl", "start", "wgupdown.service")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("start failed: %v: %s", err, string(out))
	}
	return nil
}

// Down brings the VPN down
func Down() error {
	cmd := exec.Command("systemctl", "stop", "wgupdown.service")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("stop failed: %v: %s", err, string(out))
	}
	return nil
}

func IsUp() (bool, error) {
	cmd := exec.Command("systemctl", "is-active", "--quiet", "wgupdown.service")
	err := cmd.Run()
	if err == nil {
		return true, nil
	}
	if _, ok := err.(*exec.ExitError); ok {
		return false, nil
	}
	return false, err
}

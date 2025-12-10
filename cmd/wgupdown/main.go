package main

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/Syrenny/wgupdown-gui/pkg/wireguard"
)

func validateIface(iface string) error {
	matched, err := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, iface)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("invalid interface name: %s", iface)
	}
	return nil
}

func main() {
	ctx := context.Background()
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: %s <up|down|status> <iface>\n", os.Args[0])
		os.Exit(2)
	}

	action := os.Args[1]
	iface := os.Args[2]

	if err := validateIface(iface); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch action {
	case "up":
		if err := wireguard.Up(ctx, iface); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "down":
		if err := wireguard.Down(ctx, iface); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "status":
		up, err := wireguard.IsUp(iface)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if up {
			fmt.Println("up")
		} else {
			fmt.Println("down")
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown action: %s\n", action)
		os.Exit(2)
	}
}

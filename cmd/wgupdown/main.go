package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Syrenny/wgupdown-gui/pkg/wireguard"
)

func main() {
	ctx := context.Background()
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: %s <up|down|status> <iface>\n", os.Args[0])
		os.Exit(2)
	}

	action := os.Args[1]
	ifaceName := os.Args[2]

	switch action {
	case "up":
		if err := wireguard.Up(ctx, ifaceName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "down":
		if err := wireguard.Down(ctx, ifaceName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "status":
		up, err := wireguard.IsUp(ctx, ifaceName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if up {
			fmt.Println("up")
			return
		}

		fmt.Println("down")
	default:
		fmt.Fprintf(os.Stderr, "unknown action: %s\n", action)
		os.Exit(2)
	}
}

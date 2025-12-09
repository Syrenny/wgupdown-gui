package main

import (
	"fmt"
	"log"
	"time"

	wgupdowngui "github.com/Syrenny/wgupdown-gui"
	"github.com/Syrenny/wgupdown-gui/pkg/wgupdown"

	"github.com/Syrenny/wgupdown-gui/config"
	"github.com/getlantern/systray"
)

const configPath = "/etc/wgupdown-gui/config.yaml"
const envDir = "/etc/wgupdown-gui/wgupdown-gui.auto.env"

func main() {
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	err = config.CreateEnvFile(envDir, cfg)
	if err != nil {
		log.Fatalf("EnvFile error: %s", err)
	}

	// Run systray
	systray.Run(onReady, onExit)
}

func onReady() {
	// load icon
	iconData, err := wgupdowngui.Assets.ReadFile("assets/icon.png")
	if err != nil {
		log.Printf("failed to read icon: %v", err)
	} else {
		systray.SetIcon(iconData)
	}
	systray.SetTitle("wgupdown")
	systray.SetTooltip("wgupdown — toggle wg0 quickly")

	// menu
	toggle := systray.AddMenuItem("Toggle VPN", "Toggle wg0 up/down")
	quit := systray.AddMenuItem("Quit", "Quit wgupdown")

	// update menu text according to current state initially
	updateToggleText(toggle)

	// run goroutine to handle clicks
	go func() {
		for {
			select {
			case <-toggle.ClickedCh:
				handleToggle(toggle)
			case <-quit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()

	// keep updating status every 5s (in case interface changed outside)
	go func() {
		for range time.Tick(5 * time.Second) {
			updateToggleText(toggle)
		}
	}()
}

func onExit() {
	// cleanup if needed
}

// handleToggle performs up/down and updates menu label
func handleToggle(toggle *systray.MenuItem) {
	up, err := wgupdown.IsUp()
	if err != nil {
		// show temporary tooltip by changing title briefly
		systray.SetTitle("Error: " + err.Error())
		log.Println("failed to check status:", err)
		return
	}

	if up {
		if err := wgupdown.Down(); err != nil {
			showErr(err)
			return
		}
	} else {
		if err := wgupdown.Up(); err != nil {
			showErr(err)
			return
		}
	}
	// small delay to let system apply change
	time.Sleep(400 * time.Millisecond)
	updateToggleText(toggle)
}

func updateToggleText(toggle *systray.MenuItem) {
	up, err := wgupdown.IsUp()
	if err != nil {
		toggle.SetTitle("Toggle VPN (status unknown)")
		return
	}
	if up {
		toggle.SetTitle("Disable VPN (wg0 is up)")
	} else {
		toggle.SetTitle("Enable VPN (wg0 is down)")
	}
}

func showErr(err error) {
	// quick log + tooltip
	log.Println("wgupdown error:", err)
	systray.SetTooltip(fmt.Sprintf("Error: %v", err))
}

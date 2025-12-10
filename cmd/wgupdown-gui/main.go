package main

import (
	"log"

	"github.com/Syrenny/wgupdown-gui/config"
	"github.com/Syrenny/wgupdown-gui/internal/app"
)

const (
	configPath = "/etc/wgupdown-gui/config.yaml"
)

func main() {
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run app
	app.Run(*cfg)
}

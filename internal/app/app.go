package app

import (
	"context"

	"github.com/Syrenny/wgupdown-gui/config"
	"github.com/Syrenny/wgupdown-gui/internal/service"
	"github.com/getlantern/systray"
)

func Run(cfg config.Config) {
	ctx, cancel := context.WithCancel(context.Background())

	deps := service.ServicesDependencies{
		Ctx: ctx,
		Cfg: cfg,
	}
	services := service.NewServices(deps)

	// Run systray
	systrayApp := NewSystrayApp(ctx, cancel, cfg, services)

	systray.Run(systrayApp.OnReady, systrayApp.OnExit)
}

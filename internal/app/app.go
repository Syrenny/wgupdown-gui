package app

import (
	"context"

	"github.com/Syrenny/wgupdown-gui/config"
	"github.com/Syrenny/wgupdown-gui/internal/service"
	"github.com/getlantern/systray"
)

func Run(cfg config.Config) {
	ctx := context.Background()

	deps := service.ServicesDependencies{
		Ctx: ctx,
		Cfg: cfg,
	}
	services := service.NewServices(deps)

	// Run systray
	systray_app := NewSystrayApp(ctx, cfg, services)

	systray.Run(systray_app.OnReady, systray_app.OnExit)
}

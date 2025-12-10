package service

import (
	"context"
	"fmt"
	"log"

	"github.com/Syrenny/wgupdown-gui/config"
	"github.com/Syrenny/wgupdown-gui/pkg/wgupdown"
	"github.com/getlantern/systray"
)

type GuiService struct {
	ctx context.Context
	cfg config.Config
}

func NewGuiService(ctx context.Context, cfg config.Config) *GuiService {
	return &GuiService{
		ctx: ctx,
		cfg: cfg,
	}
}

// handleToggle performs up/down and updates menu label
func (s *GuiService) HandleToggle(toggle *systray.MenuItem) error {
	up, err := wgupdown.IsUp(s.ctx, s.cfg.Wireguard.Interface)
	if err != nil {
		log.Println("failed to check status:", err)
		return err
	}

	if up {
		if err := wgupdown.Down(s.ctx, s.cfg.Wireguard.Interface); err != nil {
			return err
		}
	} else {
		if err := wgupdown.Up(s.ctx, s.cfg.Wireguard.Interface); err != nil {
			return err
		}
	}
	return nil
}

func (s *GuiService) UpdateToggleText(toggle *systray.MenuItem) error {
	up, err := wgupdown.IsUp(s.ctx, s.cfg.Wireguard.Interface)
	if err != nil {
		return err
	}
	if up {
		toggle.SetTitle("Disable VPN")
	} else {
		toggle.SetTitle("Enable VPN")
	}
	return nil
}

func (s *GuiService) ShowErr(err error) {
	// quick log + tooltip
	log.Println("wgupdown error:", err)
	systray.SetTooltip(fmt.Sprintf("Error: %v", err))
}

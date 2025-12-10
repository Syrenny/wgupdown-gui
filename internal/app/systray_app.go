package app

import (
	"context"
	"log"

	wgupdowngui "github.com/Syrenny/wgupdown-gui"
	"github.com/Syrenny/wgupdown-gui/config"
	"github.com/Syrenny/wgupdown-gui/internal/service"
	"github.com/getlantern/systray"
)

type SystrayApp interface {
	OnReady()
	OnExit()
}

type SystrayAppImpl struct {
	ctx      context.Context
	cfg      config.Config
	services *service.Services
}

func NewSystrayApp(ctx context.Context, cfg config.Config, services *service.Services) *SystrayAppImpl {
	return &SystrayAppImpl{
		ctx:      ctx,
		cfg:      cfg,
		services: services,
	}
}

func (s *SystrayAppImpl) OnReady() {
	iconData, err := wgupdowngui.Assets.ReadFile("assets/icon.png")
	if err != nil {
		log.Printf("failed to read icon: %v", err)
	} else {
		systray.SetIcon(iconData)
	}

	// menu
	toggle := systray.AddMenuItem("Toggle VPN", "Toggle VPN up/down")

	// update menu text according to current state initially
	err = s.services.Gui.UpdateToggleText(toggle)
	if err != nil {
		s.services.Gui.ShowErr(err)
	}

	// run goroutine to handle clicks
	go func() {
		for range toggle.ClickedCh {
			err := s.services.Gui.UpdateToggleText(toggle)
			if err != nil {
				s.services.Gui.ShowErr(err)
			}

			err = s.services.Gui.HandleToggle(toggle)
			if err != nil {
				s.services.Gui.ShowErr(err)
			}

			err = s.services.Gui.UpdateToggleText(toggle)
			if err != nil {
				s.services.Gui.ShowErr(err)
			}
		}
	}()
}

func (s *SystrayAppImpl) OnExit() {
	// Implementation for tray exit event
}

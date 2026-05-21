package service

import (
	"context"

	"github.com/Syrenny/wgupdown-gui/config"
	"github.com/getlantern/systray"
)

type WgUpDown interface {
	Up(
		ctx context.Context,
		ifaceName string,
	) error
	Down(
		ctx context.Context,
		ifaceName string,
	) error
	IsUp(
		ctx context.Context,
		ifaceName string,
	) (bool, error)
}

type Gui interface {
	HandleToggle(toggle *systray.MenuItem) error
	UpdateToggleText(toggle *systray.MenuItem) error
	ShowErr(err error)
}

type Services struct {
	Gui Gui
}

type ServicesDependencies struct {
	Ctx context.Context
	Cfg config.Config
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Gui: NewGuiService(deps.Ctx, deps.Cfg),
	}
}

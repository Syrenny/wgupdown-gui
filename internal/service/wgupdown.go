package service

import (
	"context"

	"github.com/Syrenny/wgupdown-gui/pkg/wgupdown"
)

type WgUpDownService struct {
	ctx context.Context
}

func NewWgUpDownService(ctx context.Context) *WgUpDownService {
	return &WgUpDownService{
		ctx: ctx,
	}
}

func (s *WgUpDownService) Up(
	ctx context.Context,
	ifaceName string,
) error {
	return wgupdown.Up(ctx, ifaceName)
}

func (s *WgUpDownService) Down(
	ctx context.Context,
	ifaceName string,
) error {
	return wgupdown.Down(ctx, ifaceName)
}

func (s *WgUpDownService) IsUp(
	ctx context.Context,
	ifaceName string,
) (bool, error) {
	return wgupdown.IsUp(ctx, ifaceName)
}

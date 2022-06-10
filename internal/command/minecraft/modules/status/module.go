package status

import (
	"context"

	"github.com/portcullis/application"
)

type cfg struct {
}

type module struct {
	cfg *cfg
}

func New() application.Module {
	return &module{
		cfg: &cfg{},
	}
}

func (m *module) Start(ctx context.Context) error {
	// TODO: periodically poll the status and report as metrics
	return nil
}

func (m *module) Stop(cctx context.Context) error {
	return nil
}

package gamerules

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"github.com/portcullis/application"
	"github.com/renevo/mcutils/internal/command/minecraft/modules/ext"
	"github.com/renevo/mcutils/pkg/minecraft"
)

type cfg struct {
	Rules map[string]string `hcl:"game_rules,optional"`
}

type module struct {
	cfg *cfg
}

func New() application.Module {
	return &module{
		cfg: &cfg{
			Rules: map[string]string{},
		},
	}
}

func (m *module) Config() (interface{}, error) {
	return m.cfg, nil
}

func (m *module) Start(ctx context.Context) error {
	subscriber := ext.Subscriber(ctx)

	ch, err := subscriber.Subscribe(ctx, minecraft.StateOnline)
	if err != nil {
		return errors.Wrap(err, "failed to subscribe to server state")
	}

	ctx, cancel := context.WithCancel(ctx)
	go m.setGameRules(ctx, cancel, ch)

	return nil
}

func (m *module) Stop(ctx context.Context) error {
	return nil
}

func (m *module) setGameRules(ctx context.Context, cancelFn context.CancelFunc, ch <-chan *message.Message) {
	select {
	case <-ctx.Done():
		return
	case msg := <-ch:
		msg.Ack()

		log := ext.Logger(ctx)
		srv := ext.Minecraft(ctx)
		for rule, value := range m.cfg.Rules {
			if err := srv.ExecuteCommand(fmt.Sprintf("gamerule %s %s", rule, value)); err != nil {
				log.Errorf("failed to set game rule %q", rule)
			}
		}
		cancelFn()
	}
}

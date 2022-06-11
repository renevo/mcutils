package startupcommands

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"github.com/portcullis/application"
	"github.com/renevo/mcutils/internal/command/minecraft/modules/ext"
	"github.com/renevo/mcutils/pkg/minecraft"
)

type cfg struct {
	Commands []string `hcl:"startup_commands,optional"`
}

type module struct {
	cfg *cfg
}

func New() application.Module {
	return &module{
		cfg: &cfg{},
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
	go m.runStartupCommands(ctx, cancel, ch)

	return nil
}

func (m *module) Stop(ctx context.Context) error {
	return nil
}

func (m *module) runStartupCommands(ctx context.Context, cancelFn context.CancelFunc, ch <-chan *message.Message) {
	log := ext.Logger(ctx)
	srv := ext.Minecraft(ctx)
	defer cancelFn()

	select {
	case <-ctx.Done():
		return
	case msg := <-ch:
		msg.Ack()

		// insert a delay before we push commands
		t := time.NewTimer(time.Second * 5)
		defer t.Stop()

		select {
		case <-ctx.Done():
			return
		case <-t.C:

		}

		for _, cmd := range m.cfg.Commands {
			if err := srv.ExecuteCommand(cmd); err != nil {
				log.Errorf("Failed to execute command %q: %v", cmd, err)
			}
		}
	}
}

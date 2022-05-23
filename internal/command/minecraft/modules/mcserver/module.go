package mcserver

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/portcullis/application"
	"github.com/portcullis/logging"
	"github.com/renevo/mcutils/internal/command/minecraft/modules/ext"
	"github.com/renevo/mcutils/pkg/minecraft"
	"github.com/sirupsen/logrus"
)

type cfg struct {
	Minecraft *minecraft.Server `hcl:"minecraft,block"`
}

type module struct {
	cfg *cfg
	wg  sync.WaitGroup
}

func New() application.Module {
	m := &module{
		cfg: &cfg{
			Minecraft: minecraft.Default(),
		},
	}

	return m
}

func (m *module) Config() (interface{}, error) {
	return m.cfg, nil
}

func (m *module) Initialize(ctx context.Context) (context.Context, error) {
	// create our global logger
	log := logrus.WithFields(logrus.Fields{"version": m.cfg.Minecraft.Version, "snapshot": m.cfg.Minecraft.Snapshot, "name": m.cfg.Minecraft.Name})
	if m.cfg.Minecraft.FabricJar() != "" {
		log = log.WithFields(logrus.Fields{"flavor": "fabric", "fabric": m.cfg.Minecraft.FabricVersionLoader})
	} else {
		log = log.WithField("flavor", "vanilla")
	}

	// wire up the portcullis logger to the logrus one
	logging.DefaultLog = logging.New(logging.WithWriter(logging.WriterFunc(func(e logging.Entry) {
		switch e.Level {
		case logging.LevelError:
			log.Errorf(e.Message, e.Arguments...)
		case logging.LevelWarning:
			log.Warningf(e.Message, e.Arguments...)
		case logging.LevelInformational:
			log.Infof(e.Message, e.Arguments...)
		default:
			log.Debugf(e.Message, e.Arguments...)
		}
	})))

	// context injections
	ctx = ext.WithLogger(ctx, log)
	ctx = ext.WithMinecraft(ctx, m.cfg.Minecraft)

	return ctx, nil
}

func (m *module) Start(ctx context.Context) error {
	log := ext.Logger(ctx)
	srv := ext.Minecraft(ctx)

	_, err := srv.Install(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to install server")
	}

	m.wg.Add(1)
	go func() {
		if err := srv.Run(ctx, log); err != nil {
			application.FromContext(ctx).Exit(errors.Wrap(err, "server failed to run"))
		}

		m.wg.Done()
	}()

	return nil
}

func (m *module) Stop(ctx context.Context) error {
	log := ext.Logger(ctx)
	srv := ext.Minecraft(ctx)

	log.Info("Saving minecraft server")
	if err := srv.ExecuteCommand("save-all"); err != nil {
		log.Errorf("Failed to save: %v", err)
	}

	log.Info("Stopping minecraft server")
	if err := srv.ExecuteCommand("stop"); err != nil {
		log.Errorf("Failed to stop - server may be zombied: %v", err)
	}

	m.wg.Wait()
	return nil
}

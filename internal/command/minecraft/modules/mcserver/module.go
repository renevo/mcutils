package mcserver

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/portcullis/application"
	"github.com/portcullis/logging"
	"github.com/renevo/mcutils/internal/command/minecraft/modules/ext"
	"github.com/renevo/mcutils/internal/control"
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
	publisher := ext.Publisher(ctx)
	subscriber := ext.Subscriber(ctx)

	_, err := srv.Install(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to install server")
	}

	start := time.Now()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	subscription, err := subscriber.Subscribe(ctx, minecraft.StateOnline)
	if err != nil {
		return errors.Wrapf(err, "failed to subscribe to %q", minecraft.StateOnline)
	}

	m.wg.Add(1)
	go func() {
		if err := srv.Run(ctx, log, publisher); err != nil {
			application.FromContext(ctx).Exit(errors.Wrap(err, "server failed to run"))
		}

		m.wg.Done()
	}()

	select {
	case <-ctx.Done():
	case msg := <-subscription:
		msg.Ack()
		log.Infof("Server online after %v", time.Since(start))
	}

	return nil
}

func (m *module) Stop(ctx context.Context) error {
	log := ext.Logger(ctx)
	srv := ext.Minecraft(ctx)
	ctrl := control.MinecraftController{Server: srv, Subscriber: ext.Subscriber(ctx)}

	// at some point we might want to have a timeout, just not sure how we would deal with that, since we have to wait for the pid to close, unless we want to implement an srv.Kill()
	ctx = context.Background()

	// might move this to another module, and have it do a countdown
	_ = ctrl.Emote(ctx, "is shutting down")

	log.Info("Saving minecraft server")
	if err := ctrl.SaveGame(ctx); err != nil {
		log.Errorf("Failed to save: %v", err)
	}

	log.Info("Stopping minecraft server")
	if err := ctrl.Stop(ctx); err != nil {
		log.Errorf("Failed to stop - server may be zombied: %v", err)
	}

	m.wg.Wait()
	return nil
}

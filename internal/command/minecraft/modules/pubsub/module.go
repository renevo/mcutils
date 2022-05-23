package pubsub

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/portcullis/application"
	"github.com/renevo/mcutils/internal/command/minecraft/modules/ext"
)

type module struct {
	bus *gochannel.GoChannel
}

func New() application.Module {
	return &module{}
}

func (m *module) Initialize(ctx context.Context) (context.Context, error) {
	// at some point in the future we might hook this up to a fanout and include like nats or other external pubsub systems
	m.bus = gochannel.NewGoChannel(gochannel.Config{
		OutputChannelBuffer:            100,
		Persistent:                     false,
		BlockPublishUntilSubscriberAck: false,
	}, &logger{ctx})

	ctx = ext.WithPublisher(ctx, m.bus)
	ctx = ext.WithSubscriber(ctx, m.bus)

	return ctx, nil
}

func (m *module) Start(ctx context.Context) error {
	return nil
}

func (m *module) Stop(ctx context.Context) error {
	return nil
}

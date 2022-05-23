package ext

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
)

const (
	subscriberContext = contextKey("message.subscriber")
)

func WithSubscriber(ctx context.Context, subscriber message.Subscriber) context.Context {
	return context.WithValue(ctx, subscriberContext, subscriber)
}

func Subscriber(ctx context.Context) message.Subscriber {
	v := ctx.Value(subscriberContext)
	if v == nil {
		return nil
	}

	return v.(message.Subscriber)
}

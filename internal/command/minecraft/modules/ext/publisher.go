package ext

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
)

const (
	publisherContext = contextKey("message.publisher")
)

func WithPublisher(ctx context.Context, publisher message.Publisher) context.Context {
	return context.WithValue(ctx, publisherContext, publisher)
}

func Publisher(ctx context.Context) message.Publisher {
	v := ctx.Value(publisherContext)
	if v == nil {
		return nil
	}

	return v.(message.Publisher)
}

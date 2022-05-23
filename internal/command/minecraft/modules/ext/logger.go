package ext

import (
	"context"

	"github.com/sirupsen/logrus"
)

const (
	logContext = contextKey("logger")
)

func WithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, logContext, logger)
}

func Logger(ctx context.Context) *logrus.Entry {
	logger := ctx.Value(logContext)
	if logger == nil {
		return logrus.WithContext(ctx)
	}

	return logger.(*logrus.Entry)
}

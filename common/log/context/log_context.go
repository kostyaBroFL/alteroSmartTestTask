package context

import (
	"context"

	"github.com/sirupsen/logrus"
)

type loggerContextKey struct{}

var loggerKey loggerContextKey

// FromContext return the logger from this context.
func FromContext(ctx context.Context) *logrus.Entry {
	return ctx.Value(loggerKey).(*logrus.Entry)
}

// WithLogger adds logger to context and return the resulting context.
func WithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

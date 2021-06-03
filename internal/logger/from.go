package logger

import (
	"context"

	"go.uber.org/zap"
)

// From creates a logger from the current context
// adds contextual attributes if possible
func From(ctx context.Context, options ...Option) *zap.Logger {
	logger := zap.L()

	for _, option := range options {
		logger = option(logger)
	}

	fields := make([]zap.Field, 0, 2)

	return logger.With(fields...)
}

// Option defines a function to set modify the logger of a Logger
type Option func(*zap.Logger) *zap.Logger

// WithBase sets the base zap Logger for a Logger
func WithBase(base *zap.Logger) Option {
	return func(_ *zap.Logger) *zap.Logger {
		return base
	}
}

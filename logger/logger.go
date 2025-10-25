package logger

import (
	"context"
)

const loggerKey = "logger"

type Logger interface {
}

func WithLogger(ctx context.Context, log Logger) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}
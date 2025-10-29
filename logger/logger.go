package logger

import (
	"context"

	"go.uber.org/zap"
)

const (
	defaultLogLevel = "info"
	loggerRequestIDKey = "x-request-id"
	loggerTraceIDKey = "x-trace-id"
	loggerKey          = "logger"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Sync() error
}

type L struct {
	z zap.Logger
}

func WithLogger(ctx context.Context, log Logger) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}

func FromContext(ctx context.Context) Logger {
	logger, ok := ctx.Value(loggerKey).(Logger)
	if !ok {
		panic("logger not found in context")
	}
	return logger
}

func NewLogger(loglevel string) (Logger, error) {
	LoggerCfg := zap.NewProductionConfig()
	LoggerCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	if loglevel == "debug" {
		LoggerCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	logger, err := LoggerCfg.Build()
	if err != nil {
		return nil, err
	}

	lo := L{z: *logger}

	return &lo, nil
}

func (l *L) Info(ctx context.Context, msg string, fields ...zap.Field) {
	id := ctx.Value(loggerRequestIDKey).(string)
	if id != "" {
		fields = append(fields, zap.String(loggerRequestIDKey, id))
	}
	l.z.Info(msg, fields...)
}

func (l *L) Error(ctx context.Context, msg string, fields ...zap.Field) {
	id := ctx.Value(loggerRequestIDKey).(string)
	if id != "" {
		fields = append(fields, zap.String(loggerRequestIDKey, id))
	}
	l.z.Error(msg, fields...) 
}

func (l *L) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	id := ctx.Value(loggerRequestIDKey).(string)
	if id != "" {
		fields = append(fields, zap.String(loggerRequestIDKey, id))
	}
	l.z.Debug(msg, fields...)
}

func (l *L) Sync() error {
	return l.z.Sync()
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, loggerRequestIDKey, requestID)
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, loggerRequestIDKey, traceID)
}
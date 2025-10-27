package logger

import (
	"context"

	"go.uber.org/zap"
)

const (
	loggerRequestIDKey = "x-request-id"
	loggerTraceIDKey = "x-trace-id"
	loggerKey          = "logger"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Debug(ctx context.Context, msg string, fields ...zap.Field)
}

type L struct {
	z zap.Logger
}

func WithLogger(ctx context.Context, log Logger) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}

func NewLogger(loglevel string) (Logger, error) {
	LoggerCfg := zap.NewProductionConfig()
	LoggerCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	if loglevel == "dev" {
		LoggerCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	logger, err := LoggerCfg.Build()
	if err != nil {
		return nil, err
	}
	// defer logger.Sync()

	lo := L{*logger}

	return &lo, nil
}

func (l *L) Info(ctx context.Context, msg string, fields ...zap.Field) {
	id := ctx.Value(loggerRequestIDKey).(string)
	fields = append(fields, zap.String(loggerRequestIDKey, id))
	l.z.Info(msg, fields...)
}

func (l *L) Error(ctx context.Context, msg string, fields ...zap.Field) {
	id := ctx.Value(loggerRequestIDKey).(string)
	fields = append(fields, zap.String(loggerRequestIDKey, id))
	l.z.Error(msg, fields...) 
}

func (l *L) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	id := ctx.Value(loggerRequestIDKey).(string)
	fields = append(fields, zap.String(loggerRequestIDKey, id))
	l.z.Debug(msg, fields...)
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, loggerRequestIDKey, requestID)
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, loggerRequestIDKey, traceID)
}
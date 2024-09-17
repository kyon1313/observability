package apw_logging

import (
	"context"

	"go.uber.org/zap"
)

// OtelLogging defines methods for logging operations.
type OtelLogging interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	Panic(args ...interface{})
	Panicf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Logf(template string, args ...interface{})
	WithContext(ctx context.Context) OtelLogging
}

type otelLog struct {
	logger *zap.SugaredLogger
}

// NewOtelLogging creates a new instance of OtelLogging.
func NewOtelLogging() OtelLogging {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	return &otelLog{
		logger: sugar,
	}
}

func (l *otelLog) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *otelLog) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(template, args...)
}

func (l *otelLog) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *otelLog) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
}

func (l *otelLog) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *otelLog) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args...)
}

func (l *otelLog) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *otelLog) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args...)
}

func (l *otelLog) DPanic(args ...interface{}) {
	l.logger.DPanic(args...)
}

func (l *otelLog) DPanicf(template string, args ...interface{}) {
	l.logger.DPanicf(template, args...)
}

func (l *otelLog) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *otelLog) Panicf(template string, args ...interface{}) {
	l.logger.Panicf(template, args...)
}

func (l *otelLog) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *otelLog) Fatalf(template string, args ...interface{}) {
	l.logger.Fatalf(template, args...)
}

func (l *otelLog) Logf(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
}

func (l *otelLog) WithContext(ctx context.Context) OtelLogging {
	return &otelLog{
		logger: l.logger.With(ctx),
	}
}

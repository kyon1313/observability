package otelBuilder

import (
	apw_logging "github.com/kyon1313/observability/logs"
	apw_tracing "github.com/kyon1313/observability/tracing"
)

type Otel struct {
	Tracing apw_tracing.OtelTracing
	Logs    apw_logging.OtelLogging
}

func NewOtel(tracing apw_tracing.OtelTracing, logs apw_logging.OtelLogging) *Otel {
	return &Otel{
		Tracing: tracing,
		Logs:    logs,
	}
}

func (o *Otel) GetTracing() apw_tracing.OtelTracing {

	return o.Tracing
}

func (o *Otel) GetLogs() apw_logging.OtelLogging {
	return o.Logs
}

package otelBuilder

import (
	"context"
	"fmt"
	apw_logging "otel-test/logs"
	apw_tracing "otel-test/tracing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

const (
	tracerName = "default-tracer"
)

type Header map[string]string

type OtelTracingBuilder struct {
	serviceName        string
	traceOpts          []trace.BatchSpanProcessorOption
	traceExporterOpts  []otlptracehttp.Option
	useConsoleExporter bool
}

func NewOtelTracingBuilder() *OtelTracingBuilder {
	return &OtelTracingBuilder{}
}

func (b *OtelTracingBuilder) WithInsecure(insecure bool) *OtelTracingBuilder {
	if insecure {
		b.traceExporterOpts = append(b.traceExporterOpts, otlptracehttp.WithInsecure())
	}
	return b
}

func (b *OtelTracingBuilder) WithEndpoint(otlpEndpoint string) *OtelTracingBuilder {
	if otlpEndpoint != "" {
		b.traceExporterOpts = append(b.traceExporterOpts, otlptracehttp.WithEndpoint(otlpEndpoint))
	}
	return b
}
func (b *OtelTracingBuilder) WithHeaders(headers Header) *OtelTracingBuilder {
	headerMap := make(map[string]string)
	for key, value := range headers {
		headerMap[key] = value
	}
	b.traceExporterOpts = append(b.traceExporterOpts, otlptracehttp.WithHeaders(headerMap))
	return b
}

func (b *OtelTracingBuilder) WithAuthHeader(token string) *OtelTracingBuilder {
	return b.WithHeaders(Header{
		"Authorization": "Bearer " + token,
	})
}

func (b *OtelTracingBuilder) WithServiceName(serviceName string) *OtelTracingBuilder {
	if serviceName != "" {
		b.serviceName = serviceName
	}
	return b
}

func (b *OtelTracingBuilder) WithTraceBatchSpanProcessorOption(opts ...trace.BatchSpanProcessorOption) *OtelTracingBuilder {
	b.traceOpts = opts
	return b
}

func (b *OtelTracingBuilder) WithConsoleExporter() *OtelTracingBuilder {
	b.useConsoleExporter = true
	return b
}

func (b *OtelTracingBuilder) Build(ctx context.Context, l apw_logging.OtelLogging) (apw_tracing.OtelTracing, error) {
	var traceExporter trace.SpanExporter
	var err error

	if b.useConsoleExporter {
		traceExporter, err = newConsoleTraceExporter()
		if err != nil {
			return nil, fmt.Errorf("failed to create console trace exporter: %w", err)
		}
	} else {
		traceExporter, err = otlptracehttp.New(ctx, b.traceExporterOpts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP trace exporter with options %v: %w", b.traceExporterOpts, err)
		}
	}

	resourceOpts := resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceName(b.serviceName))
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter, b.traceOpts...),
		trace.WithResource(resourceOpts),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tracerProvider)

	// Return tracing instance
	return apw_tracing.NewTracing(tracerProvider.Tracer(b.serviceName), l), nil
}

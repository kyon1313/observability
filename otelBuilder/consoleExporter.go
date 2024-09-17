package otelBuilder

import (
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	oteltrace "go.opentelemetry.io/otel/sdk/trace"
)

// Console Exporter, only for testing
func newConsoleTraceExporter() (oteltrace.SpanExporter, error) {
	return stdouttrace.New(stdouttrace.WithPrettyPrint())
}

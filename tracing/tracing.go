package _tracing

import (
	"context"
	"fmt"
	"net/http"
	_logging "otel-test/logs"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type OtelTracing interface {
	StartSpan(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
	EndSpan(span trace.Span)
	SetStatus(span trace.Span, code codes.Code, description string)
	AddAttribute(span trace.Span, key string, value any)
	AddEvent(ctx context.Context, span trace.Span, eventName string, opts ...trace.EventOption)
	SetOKStatus(span trace.Span, description string, attrs ...attribute.KeyValue)
	SetNoContentStatus(span trace.Span, description string, attrs ...attribute.KeyValue)
	GetTracer() trace.Tracer

	AddAttributes(span trace.Span, err error, attrs ...attribute.KeyValue)
	RecordError(span trace.Span, err error, source string)
}

type tracing struct {
	tracer trace.Tracer
	l      _logging.OtelLogging
}

// NewTracing initializes a new OtelTracing instance with the given Tracer.
func NewTracing(tracer trace.Tracer, l _logging.OtelLogging) OtelTracing {

	return &tracing{tracer: tracer, l: l}
}

// StartSpan creates a new span with the given name and options.
func (t *tracing) StartSpan(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, spanName, opts...)
}

// EndSpan ends the given span.
func (t *tracing) EndSpan(span trace.Span) {
	span.End()
}

// SetStatus sets the status code and description for the given span.
func (t *tracing) SetStatus(span trace.Span, code codes.Code, description string) {
	span.SetStatus(code, description)
}

// AddAttributes adds multiple attributes to the given span.
func (t *tracing) AddAttribute(span trace.Span, key string, value any) {
	switch v := value.(type) {
	case string:
		span.SetAttributes(attribute.String(key, v))
	case int:
		span.SetAttributes(attribute.Int(key, v))
	case float64:
		span.SetAttributes(attribute.Float64(key, v))
	case bool:
		span.SetAttributes(attribute.Bool(key, v))
	default:
		span.SetAttributes(attribute.String(key, fmt.Sprintf("%v", v)))
	}
}

// AddEvent records an event with a name and optional attributes in the given span.
func (t *tracing) AddEvent(ctx context.Context, span trace.Span, eventName string, opts ...trace.EventOption) {
	span.AddEvent(eventName, opts...)
}

// SetOKStatus sets the status to OK with an optional description and attributes.
func (t *tracing) SetOKStatus(span trace.Span, description string, attrs ...attribute.KeyValue) {
	span.SetStatus(codes.Ok, description)
	span.SetAttributes(attrs...)
}

// SetNoContentStatus sets the status to indicate no content with an optional description and attributes.
func (t *tracing) SetNoContentStatus(span trace.Span, description string, attrs ...attribute.KeyValue) {
	span.SetStatus(codes.Code(http.StatusNoContent), description)
	span.SetAttributes(attrs...)
}

// GetTracer returns the underlying Tracer instance.
func (t *tracing) GetTracer() trace.Tracer {
	return t.tracer
}

func (s *tracing) AddAttributes(span trace.Span, err error, attrs ...attribute.KeyValue) {
	if err != nil {
		span.SetAttributes(attrs...)
		return
	}

	attrs = append(attrs, attribute.Int("status.Code", 1))
	span.SetAttributes(attrs...)
}

func (t *tracing) RecordError(span trace.Span, err error, source string) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
	span.SetAttributes(attribute.String("error.source", source))
}

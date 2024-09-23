package _tracing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	_logging "github.com/kyon1313/observability/logs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type OtelTracing interface {
	StartSpan(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
	EndSpan(span trace.Span, err error)
	SetStatus(span trace.Span, code codes.Code, description string)
	AddAttribute(span trace.Span, key string, value any)
	AddEvent(span trace.Span, eventName string, attrs ...attribute.KeyValue)
	SetOKStatus(span trace.Span, description string, attrs ...attribute.KeyValue)
	SetNoContentStatus(span trace.Span, description string, attrs ...attribute.KeyValue)
	GetTracer() trace.Tracer

	AddAttributes(span trace.Span, err error, attrs ...attribute.KeyValue)
	RecordError(span trace.Span, err error, source string)

	ExtractSpanContext(ctx context.Context, r *http.Request) context.Context
	InjectSpanContext(ctx context.Context, r *http.Request)
	SpanFromContext(ctx context.Context) trace.Span
	StartSpanFromContext(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span)

	AddBaggage(ctx context.Context, key, value string) context.Context
	GetBaggage(ctx context.Context, key string) string
	LogTrace(span trace.Span, err *error, layer string, response any) func()
}

type tracing struct {
	tracer trace.Tracer
	l      _logging.OtelLogging
}

// NewTracing initializes a new OtelTracing instance with the given Tracer.
func NewTracing(tracer trace.Tracer, l _logging.OtelLogging) OtelTracing {
	return &tracing{tracer: tracer, l: l}
}

func (t *tracing) LogTrace(span trace.Span, err *error, layer string, response any) func() {
	return func() {
		if *err != nil && span.IsRecording() {
			t.RecordError(span, *err, layer)
			t.SetStatus(span, codes.Error, (*err).Error())
		} else if *err == nil {
			t.SetOKStatus(span, "Operation completed successfully")
		}

		// Marshal the response to JSON
		var jsonResponse string
		if response != nil {
			if jsonBytes, err := json.Marshal(response); err == nil {
				jsonResponse = string(jsonBytes)
			} else {
				jsonResponse = fmt.Sprintf("Failed to marshal response: %v", err)
			}
		} else {
			jsonResponse = "null"
		}

		// Record the response as an event
		if span.IsRecording() {
			t.AddEvent(span, layer, attribute.String("response", jsonResponse))
		}

		t.EndSpan(span, *err)
	}
}

// StartSpan creates a new span with the given name and options.
func (t *tracing) StartSpan(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, spanName, opts...)
}

// EndSpan ends the given span.
func (t *tracing) EndSpan(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "Success")
	}
	span.End()
}

// SetStatus sets the status code and description for the given span.
func (t *tracing) SetStatus(span trace.Span, code codes.Code, description string) {
	span.SetStatus(code, description)
}

// AddAttribute adds an attribute to the given span.
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
func (t *tracing) AddEvent(span trace.Span, eventName string, attrs ...attribute.KeyValue) {
	span.AddEvent(eventName, trace.WithAttributes(attrs...))
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

// AddAttributes adds multiple attributes to the given span.
func (t *tracing) AddAttributes(span trace.Span, err error, attrs ...attribute.KeyValue) {
	if err != nil {
		span.SetAttributes(attrs...)
		return
	}

	attrs = append(attrs, attribute.Int("status.Code", 1))
	span.SetAttributes(attrs...)
}

// RecordError records an error in the given span.
func (t *tracing) RecordError(span trace.Span, err error, source string) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
	span.SetAttributes(attribute.String("error.source", source))
}

// ExtractSpanContext extracts the span context from the incoming request headers.
func (t *tracing) ExtractSpanContext(ctx context.Context, r *http.Request) context.Context {
	propagator := otel.GetTextMapPropagator()
	return propagator.Extract(ctx, propagation.HeaderCarrier(r.Header))
}

// InjectSpanContext injects the span context into the outgoing request headers.
func (t *tracing) InjectSpanContext(ctx context.Context, r *http.Request) {
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, propagation.HeaderCarrier(r.Header))
}

// SpanFromContext retrieves the span from the context.
func (t *tracing) SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// StartSpanFromContext starts a new span with the given name and options, using the context.
func (t *tracing) StartSpanFromContext(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, spanName, opts...)
}

// GetTracer returns the underlying Tracer instance.
func (t *tracing) GetTracer() trace.Tracer {
	return t.tracer
}

// AddBaggage adds a key-value pair to the baggage in the context.
func (t *tracing) AddBaggage(ctx context.Context, key, value string) context.Context {
	b, _ := baggage.NewMember(key, value)
	bag, _ := baggage.New(b)
	return baggage.ContextWithBaggage(ctx, bag)
}

// GetBaggage retrieves the value of a key from the baggage in the context.
func (t *tracing) GetBaggage(ctx context.Context, key string) string {
	bag := baggage.FromContext(ctx)
	member := bag.Member(key)
	return member.Value()
}

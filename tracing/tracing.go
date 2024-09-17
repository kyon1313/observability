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

// func (s *tracing) RecordError(ctx context.Context, span trace.Span, err *errs.ErrorService) {
// 	s.SetStatus(span, codes.Error, err.Error())
// 	s.AddAttributes(
// 		span,
// 		err,
// 		semconv.ExceptionMessageKey.String(err.Error()),
// 		semconv.ExceptionTypeKey.String(err.StatusText),
// 		semconv.ExceptionStacktraceKey.String(err.GetStackTrace()),
// 		attribute.Int("status.Code", 2),
// 	)
// 	span.AddEvent("ExceptionOccurred", trace.WithAttributes(
// 		semconv.ExceptionTypeKey.String(err.StatusText),
// 		semconv.ExceptionMessageKey.String(err.Error()),
// 		attribute.String("exception.code", err.ErrorCode),
// 		attribute.Int("status.Code", 2),
// 	))
// 	if span.SpanContext().HasSpanID() {
// 		s.l.WithContext(ctx).Errorf("%s", err.Error())
// 	}
// 	span.RecordError(err)
// }

// func (s *tracing) AddErrorAttributes(ctx context.Context, span trace.Span, err *errs.ErrorService) {
// 	if err.StatusCode >= 400 {
// 		s.AddAttributes(
// 			span,
// 			err,
// 			semconv.ExceptionMessageKey.String(err.Error()),
// 			semconv.ExceptionTypeKey.String(err.StatusText),
// 			semconv.ExceptionStacktraceKey.String(err.GetStackTrace()),
// 		)
// 		span.AddEvent("ExceptionOccurred", trace.WithAttributes(
// 			semconv.ExceptionTypeKey.String(err.StatusText),
// 			semconv.ExceptionMessageKey.String(err.Error()),
// 			attribute.String("exception.code", err.ErrorCode),
// 		))
// 		if span.SpanContext().HasSpanID() {
// 			s.l.WithContext(ctx).Errorf("[\n\nAEA\n\n]%s", err.Error())
// 		}
// 	} else {
// 		s.AddAttributes(
// 			span,
// 			err,
// 			attribute.String("error.type", err.StatusText),
// 			attribute.String("error.message", err.Error()),
// 			attribute.String("error.stacktrace", err.GetStackTrace()),
// 		)
// 		span.AddEvent("ErrorHandled", trace.WithAttributes(
// 			attribute.String("error.type", err.StatusText),
// 			attribute.String("error.message", err.Error()),
// 			attribute.String("error.code", err.ErrorCode),
// 		))
// 		if span.SpanContext().HasSpanID() {
// 			s.l.WithContext(ctx).Warnf("%s", err.Error())
// 		}
// 	}

// }

func (s *tracing) AddAttributes(span trace.Span, err error, attrs ...attribute.KeyValue) {
	if err != nil {
		span.SetAttributes(attrs...)
		return
	}

	attrs = append(attrs, attribute.Int("status.Code", 1))
	span.SetAttributes(attrs...)
}

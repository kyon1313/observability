package otelBuilder

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	apw_logging "otel-test/logs"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func TracingMiddleware(l apw_logging.OtelLogging, tracer trace.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		if ShouldIgnoreRequest(c) {
			c.Next()
			return
		}

		ctx, span := tracer.Start(c.Request.Context(), fmt.Sprintf("HTTP %s %s", c.Request.Method, c.Request.URL.Path))
		defer span.End()

		spanContext := span.SpanContext()
		traceID := spanContext.TraceID().String()

		c.Set("X-Request-Id", traceID)
		c.Header("X-Request-Id", traceID)

		currentTime := time.Now()
		l.Debug("Request received",
			zap.String("date", currentTime.Format("2006/01/02 - 15:04:05")),
			zap.String("request_id", traceID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)

		body, err := c.GetRawData()
		if err == nil && len(body) > 0 {
			l.Debug("Request body",
				zap.String("request_id", traceID),
				zap.String("body", string(body)),
			)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		if len(c.Request.URL.RawQuery) > 0 {
			l.Debug("Request parameters",
				zap.String("request_id", traceID),
				zap.String("query_params", c.Request.URL.RawQuery),
			)
		}

		c.Request = c.Request.WithContext(ctx)

		w := &responseBodyWriter{
			body:           &bytes.Buffer{},
			ResponseWriter: c.Writer,
			statusCode:     http.StatusOK,
		}
		c.Writer = w

		c.Next()

		span.SetAttributes(
			attribute.Int("http.status_code", w.statusCode),
		)

		l.WithContext(c.Request.Context()).Debugf("Response sent",
			zap.Int("status_code", w.statusCode),
			zap.String("data", w.body.String()),
		)
	}
}

func ShouldIgnoreRequest(c *gin.Context) bool {
	pathsToIgnore := []string{
		"/health",
		"/healthcheck",
		"/metrics",
		"/swagger/",
	}

	for _, path := range pathsToIgnore {
		if c.Request.URL.Path == path || strings.HasPrefix(c.Request.URL.Path, path) {
			return true
		}
	}

	return false
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (r *responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r *responseBodyWriter) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

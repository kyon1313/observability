package otelBuilder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	apw_logging "github.com/kyon1313/observability/logs"

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

		traceID := span.SpanContext().TraceID().String()
		c.Set("X-Request-Id", traceID)
		c.Header("X-Request-Id", traceID)

		logRequestDetails(l, c, traceID)

		body, requestBody := readRequestBody(l, c, traceID)
		if requestBody != nil {
			setSpanAttributes(span, "request.body", requestBody)
		}

		c.Request = c.Request.WithContext(ctx)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer, statusCode: http.StatusOK}
		c.Writer = w

		c.Next()

		responseBody := w.body.String()
		logResponseBody(l, traceID, responseBody)

		if w.statusCode >= 400 {
			span.SetAttributes(attribute.Int("status.Code", 2))
		} else {
			setSpanAttributes(span, "response.body", parseJSON(responseBody))
			span.SetAttributes(attribute.Int("status.Code", 1))
		}

		span.SetAttributes(attribute.Int("http.status_code", w.statusCode))
	}
}

func logRequestDetails(l apw_logging.OtelLogging, c *gin.Context, traceID string) {
	currentTime := time.Now()
	l.Debug("Request received",
		zap.String("date", currentTime.Format("2006/01/02 - 15:04:05")),
		zap.String("request_id", traceID),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	)

	if len(c.Request.URL.RawQuery) > 0 {
		l.Debug("Request parameters", zap.String("request_id", traceID), zap.String("query_params", c.Request.URL.RawQuery))
	}
}

func readRequestBody(l apw_logging.OtelLogging, c *gin.Context, traceID string) ([]byte, map[string]interface{}) {
	body, err := c.GetRawData()
	if err != nil || len(body) == 0 {
		l.Debug("Failed to read request body or body is empty", zap.String("request_id", traceID), zap.Error(err))
		return nil, nil
	}

	var requestBody map[string]interface{}
	if err := json.Unmarshal(body, &requestBody); err != nil {
		l.Debug("Failed to unmarshal request body", zap.String("request_id", traceID), zap.Error(err))
		return body, nil
	}

	l.Debug("Request body read successfully", zap.String("request_id", traceID), zap.Any("body", requestBody))
	return body, requestBody
}

func logResponseBody(l apw_logging.OtelLogging, traceID, responseBody string) {
	l.Debug("Response body", zap.String("request_id", traceID), zap.String("body", responseBody))
}

func setSpanAttributes(span trace.Span, prefix string, data map[string]interface{}) {
	for key, value := range data {
		span.SetAttributes(attribute.String(fmt.Sprintf("%s.%s", prefix, key), fmt.Sprintf("%v", value)))
	}
}

func parseJSON(data string) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal([]byte(data), &result)
	return result
}

func ShouldIgnoreRequest(c *gin.Context) bool {
	pathsToIgnore := []string{"/health", "/healthcheck", "/metrics", "/swagger/"}
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

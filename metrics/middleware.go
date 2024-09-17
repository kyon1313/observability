package metrics

import (
	"time"

	"github.com/gin-gonic/gin"
)

type MetricsMiddlewareDecorator struct {
	metrics *Metrics
}

func NewMetricsMiddlewareDecorator(metrics *Metrics) *MetricsMiddlewareDecorator {
	return &MetricsMiddlewareDecorator{metrics: metrics}
}

func (m *MetricsMiddlewareDecorator) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()

		rw := &responseWriter{ResponseWriter: c.Writer}
		c.Writer = rw

		// Simulate active sessions
		m.metrics.Gauges["active_sessions"].WithLabelValues(path).Inc()
		defer m.metrics.Gauges["active_sessions"].WithLabelValues(path).Dec()

		c.Next()

		duration := time.Since(start).Seconds()

		m.metrics.Counters["http_requests_total"].WithLabelValues(path).Inc()
		m.metrics.Histograms["http_request_duration_seconds"].WithLabelValues(path).Observe(duration)

		if c.Writer.Status() >= 400 {
			m.metrics.Counters["http_errors_total"].WithLabelValues(path).Inc()
		}
	}
}

type responseWriter struct {
	gin.ResponseWriter
	size int
}

func (w *responseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}

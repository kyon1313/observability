package main

import (
	"context"
	apw_logging "otel-test/logs"
	"otel-test/metrics"
	"otel-test/otelBuilder"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

const (
	METRICURLENDPOINT = "http://localhost:8081/metrics"
	JAEGERENDPOINT    = "jaeger:4318"
)

var otelConfig = initOtel()

func initOtel() *otelBuilder.Otel {
	ctx := context.Background()
	l := apw_logging.NewOtelLogging()

	batchOpts := []trace.BatchSpanProcessorOption{
		trace.WithBatchTimeout(time.Second * 10),
	}

	tracing, err := otelBuilder.NewOtelTracingBuilder().
		WithEndpoint(JAEGERENDPOINT).
		WithInsecure(true).
		WithServiceName("testing-api").
		WithTraceBatchSpanProcessorOption(batchOpts...).
		Build(ctx, l)

	if err != nil {
		l.Error("Failed to initialize OpenTelemetry", err)
		return nil
	}

	return otelBuilder.NewOtel(tracing, l)
}

func main() {
	tracerProvider := otel.GetTracerProvider()
	tracer := tracerProvider.Tracer("apw-test")

	repo := NewRepository(tracer)
	svc := NewService(repo, tracer)
	h := NewHandler(svc, initOtel().Logs)

	metricBuilder := metrics.NewMetricsBuilder().
		AddCounter("http_requests_total", "Total number of HTTP requests", []string{"path"}).
		AddHistogram("http_request_duration_seconds", "Duration of HTTP requests in seconds", prometheus.DefBuckets, []string{"path"}).
		AddCounter("http_errors_total", "Total number of HTTP errors", []string{"path"}).
		AddGauge("active_sessions", "The current number of active sessions", []string{"path"}).
		AddGauge("queue_size", "The current size of the queue", []string{"path"}).
		Build()

	r := gin.Default()

	metricsMiddleware := metrics.NewMetricsMiddlewareDecorator(metricBuilder)
	r.Use(metricsMiddleware.Middleware(), otelBuilder.TracingMiddleware(otelConfig.Logs, tracer))

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.GET("/test2", h.PossibleErrorRequest)

	r.GET("/test3", h.SlowRequest)

	r.Run(":8080")
}

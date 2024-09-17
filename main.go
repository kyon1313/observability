package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	apw_logging "otel-test/logs"
	"otel-test/metrics"
	"otel-test/otelBuilder"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/trace"
)

const (
	metricEndpointURl = "http://localhost:8081/metrics"
	jaegerEndpointUrl = "jaeger:4318"
)

var otelConfig = initOtel()

func initOtel() *otelBuilder.Otel {
	ctx := context.Background()
	l := apw_logging.NewOtelLogging()

	batchOpts := []trace.BatchSpanProcessorOption{
		trace.WithBatchTimeout(time.Second * 10),
	}

	tracing, err := otelBuilder.NewOtelTracingBuilder().
		WithEndpoint(jaegerEndpointUrl).
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

	metricBuilder := metrics.NewMetricsBuilder().
		AddCounter("http_requests_total", "Total number of HTTP requests", []string{"path"}).
		AddHistogram("http_request_duration_seconds", "Duration of HTTP requests in seconds", prometheus.DefBuckets, []string{"path"}).
		AddCounter("http_errors_total", "Total number of HTTP errors", []string{"path"}).
		AddGauge("active_sessions", "The current number of active sessions", []string{"path"}).
		AddGauge("queue_size", "The current size of the queue", []string{"path"}).
		Build()

	r := gin.Default()

	metricsMiddleware := metrics.NewMetricsMiddlewareDecorator(metricBuilder)
	r.Use(metricsMiddleware.Middleware())

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Normal Request"})
	})

	r.GET("/test2", posibleErrorRequest)

	r.GET("/test3", slowRequest)

	r.Run(":8080")
}

func slowRequest(ctx *gin.Context) {
	c := ctx.Request.Context()

	c, span := otelConfig.Tracing.GetTracer().Start(c, "slow request")
	defer span.End()

	random := Rand(c) // Generate a random number between 1 and 5
	sleepDuration := time.Duration(random) * time.Second
	time.Sleep(sleepDuration)

	span.SetAttributes(attribute.String("request_took", fmt.Sprintf("%ds", random)))
	ctx.JSON(200, gin.H{"message": "Slow Request!"})
}

func posibleErrorRequest(ctx *gin.Context) {
	c := ctx.Request.Context()

	c, span := otelConfig.Tracing.GetTracer().Start(c, "posible error request")
	defer span.End()

	random := Rand(c)
	if random <= 3 {
		ctx.JSON(200, gin.H{"message": "Good Request!"})
		return
	}

	err := errors.New("error request")
	span.RecordError(err)
	span.SetStatus(codes.Error, "error encountered!")
	span.SetAttributes(attribute.String("error.message", err.Error()))
}

func Rand(ctx context.Context) int {

	_, span := otelConfig.Tracing.GetTracer().Start(ctx, "random number")
	defer span.End()

	return rand.Intn(5) + 1
}

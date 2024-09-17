package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	Counters   map[string]*prometheus.CounterVec
	Histograms map[string]*prometheus.HistogramVec
	Gauges     map[string]*prometheus.GaugeVec
}

type MetricsBuilder struct {
	metrics *Metrics
}

func NewMetricsBuilder() *MetricsBuilder {
	return &MetricsBuilder{
		metrics: &Metrics{
			Counters:   make(map[string]*prometheus.CounterVec),
			Histograms: make(map[string]*prometheus.HistogramVec),
			Gauges:     make(map[string]*prometheus.GaugeVec),
		},
	}
}

func (b *MetricsBuilder) AddCounter(name, help string, labels []string) *MetricsBuilder {
	b.metrics.Counters[name] = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: name,
		Help: help,
	}, labels)
	return b
}

func (b *MetricsBuilder) AddHistogram(name, help string, buckets []float64, labels []string) *MetricsBuilder {
	b.metrics.Histograms[name] = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    name,
		Help:    help,
		Buckets: buckets,
	}, labels)
	return b
}

func (b *MetricsBuilder) AddGauge(name, help string, labels []string) *MetricsBuilder {
	b.metrics.Gauges[name] = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: name,
		Help: help,
	}, labels)
	return b
}

func (b *MetricsBuilder) Build() *Metrics {
	return b.metrics
}

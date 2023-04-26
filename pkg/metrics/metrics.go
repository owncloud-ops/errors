package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	total    prometheus.Counter
	duration prometheus.Histogram
}

// NewMetrics creates new Metrics collector.
func NewMetrics() Metrics {
	const namespace, subsystem = "http", "requests"

	return Metrics{
		total: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "count_total",
			Help:      "counter of http requests made",
		}),
		duration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "duration_seconds",
			Help:      "histogram of the time (in seconds) each request took",
			Buckets:   append([]float64{.001, .003}, prometheus.DefBuckets...),
		}),
	}
}

// IncrementTotalRequests increments total requests counter.
func (w *Metrics) IncrementTotalRequests() { w.total.Inc() }

// ObserveRequestDuration observer requests duration histogram.
func (w *Metrics) ObserveRequestDuration(t time.Duration) { w.duration.Observe(t.Seconds()) }

// Register metrics with registerer.
func (w *Metrics) Register(reg prometheus.Registerer) error {
	if err := reg.Register(w.total); err != nil {
		return err
	}

	return reg.Register(w.duration)
}

package log

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	DefaultMetrics = NewMetrics()
)

type Metrics struct {
	logCounter *prometheus.CounterVec
}

func (m *Metrics) ReportLog(level logrus.Level) {
	m.logCounter.WithLabelValues(level.String()).Inc()
}

func NewMetrics() *Metrics {
	return &Metrics{
		logCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "opencoder_logs_total",
			Help: "Total number of logs produced by level",
		}, []string{"level"}),
	}
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by this Collector to the provided channel and returns once
// the last descriptor has been sent.
func (m *Metrics) Describe(ch chan<- *prometheus.Desc) {
	m.logCounter.Describe(ch)
}

// Collect is called by the Prometheus registry when collecting
// metrics. The implementation sends each collected metric via the
// provided channel and returns once the last metric has been sent.
func (m *Metrics) Collect(ch chan<- prometheus.Metric) {
	m.logCounter.Collect(ch)
}

func NewLogHook(metrics *Metrics) *Hook {
	return &Hook{metrics: metrics}
}

type Hook struct {
	metrics *Metrics
}

func (h *Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *Hook) Fire(entry *logrus.Entry) error {
	h.metrics.ReportLog(entry.Level)
	return nil
}

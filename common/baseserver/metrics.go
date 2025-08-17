package baseserver

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	serverVersionGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "opencoder",
		Subsystem: "server",
		Name:      "version",
		Help:      "Gauge of the current version of a opencoder server",
	}, []string{"version"})
)

func registerMetrics(reg *prometheus.Registry) error {
	metrics := []prometheus.Collector{
		serverVersionGauge,
	}
	for _, metric := range metrics {
		err := reg.Register(metric)
		if err != nil {
			return fmt.Errorf("failed to register metric: %w", err)
		}
	}

	return nil
}

func reportServerVersion(version string) {
	serverVersionGauge.WithLabelValues(version).Set(1)
}

package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	CheckDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "uptime_check_duration_seconds",
		Help:    "HTTP check latency in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"target", "url"})

	CheckUp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "uptime_check_up",
		Help: "1 if target is up, 0 if down",
	}, []string{"target", "url"})

	ChecksTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "uptime_checks_total",
		Help: "Total number of checks performed",
	}, []string{"target", "url", "status"})
)

func Register() {
	prometheus.MustRegister(CheckDuration, CheckUp, ChecksTotal)
}

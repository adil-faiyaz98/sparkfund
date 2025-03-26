package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	HttpRequestsTotal *prometheus.CounterVec
}

func NewMetrics() *Metrics {
	m := &Metrics{
		HttpRequestsTotal: prometheus.NewCounterVec(
			prometheus.Opts{
				Name: "http_requests_total",
				Help: "Number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
	}
	prometheus.Register(m.HttpRequestsTotal)
	return m
}

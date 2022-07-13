package okmonitor

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ApiCallMonitorEnabled bool = false

	ApiCallRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "okweb_api_call_response_duration_seconds",
			Help:    "Histogram of latencies for api call HTTP requests.",
			Buckets: []float64{.05, 0.1, .25, .5, .75, 1, 2, 5, 20, 60},
		},
		[]string{"host", "api", "method", "statusCode"},
	)
)

// init : auto run before main
func EnableApiCallMetrics() {
	if !ApiCallMonitorEnabled {
		ApiCallMonitorEnabled = true
		prometheus.MustRegister(ApiCallRequestDuration)
	}
}

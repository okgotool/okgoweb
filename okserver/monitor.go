package okserver

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	WebMonitor = &Monitor{}

	ApiAccessRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "okweb_api_response_duration_seconds",
			Help:    "Histogram of latencies for api HTTP requests.",
			Buckets: []float64{.05, 0.1, .25, .5, .75, 1, 2, 5, 20, 60},
		},
		[]string{"api", "method", "statusCode"},
	)

	ApiCallRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "okweb_api_call_response_duration_seconds",
			Help:    "Histogram of latencies for api call HTTP requests.",
			Buckets: []float64{.05, 0.1, .25, .5, .75, 1, 2, 5, 20, 60},
		},
		[]string{"host", "api", "method", "statusCode"},
	)
)

type (
	Monitor struct {
	}
)

func (w *Monitor) AddApis(router *gin.Engine) {
	router.GET("/metrics", w.MetricsHandler())
}

// MetricsHandler ：
// Swagger doc refer: https://github.com/swaggo/swag
// @Summary Metrics API
// @Description Metrics API
// @Router /metrics [get]
func (w *Monitor) MetricsHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// apiAccessMetricsMiddleware :
func (w *Monitor) ApiAccessMetricsMiddleware(c *gin.Context) {
	// metrics :
	path := c.FullPath()
	if len(path) < 1 {
		path = c.Request.URL.Path
	}

	startTime := time.Now()

	c.Next()

	// metrics : api cost time
	latency := time.Since(startTime).Seconds()
	statusCode := c.Writer.Status()
	ApiAccessRequestDuration.With(prometheus.Labels{
		"api":        path,
		"method":     c.Request.Method,
		"statusCode": fmt.Sprintf("%d", statusCode),
	}).Observe(latency)
}

// init : auto run before main
func (w *Monitor) InitApiMetrics() {
	prometheus.MustRegister(ApiAccessRequestDuration)
	prometheus.MustRegister(ApiCallRequestDuration)
}

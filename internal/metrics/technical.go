package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	RestRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hits_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"path"},
	)

	RestResponseDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "time_hits",
			Help:    "Duration of HTTP requests.",
			Buckets: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 15, 20, 25, 30, 35, 40, 50, 60, 70, 80, 90, 100},
		},
		[]string{"path", "method"},
	)

	RestEndpointsResponsesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hits_statuses",
			Help: "Statuses for HTTP responses.",
		},
		[]string{"path", "status"},
	)
)

func IncRestRequestsTotal(path string) {
	RestRequestsTotal.WithLabelValues(path).Inc()
}

func IncRestResponsesDuration(path, method string, timeServe time.Duration) {
	RestResponseDuration.WithLabelValues(path, method).Observe(float64(timeServe.Milliseconds()))
}

func IncRestResponsesStatusesTotal(path string, status int) {
	RestEndpointsResponsesTotal.WithLabelValues(path, http.StatusText(status)).Inc()
}

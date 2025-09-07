package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total requests",
		},
		[]string{"method", "endpoint", "status"},
	)
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status"},
	)
	EventCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_events_total",
			Help: "Counts events by name and status",
		},
		[]string{"event_name", "status"},
	)
)

const (
	EventPusblishExpenseApprove = "publish_expense_approve"
)

func Init() {
	prometheus.MustRegister(
		RequestsTotal,
		RequestDuration,
		EventCounter,
	)
}

func IncrementEvent(eventName, status string) {
	EventCounter.WithLabelValues(eventName, status).Inc()
}

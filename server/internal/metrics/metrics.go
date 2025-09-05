package metrics

import (
	"time"

	"github.com/gin-gonic/gin"
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
)

func Init() {
	prometheus.MustRegister(
		RequestsTotal,
		RequestDuration,
	)
}

func Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.FullPath() == "/metrics" {
			ctx.Next()
			return
		}

		start := time.Now()
		ctx.Next()

		duration := time.Since(start).Seconds()

		var status string
		if ctx.Writer.Status() >= 500 {
			status = "fail"
		} else {
			status = "ok"
		}

		RequestsTotal.WithLabelValues(ctx.Request.Method, ctx.FullPath(), status).Inc()
		RequestDuration.WithLabelValues(ctx.Request.Method, ctx.FullPath(), status).Observe(duration)
	}
}

package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "pattern", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "pattern"},
	)

	RedirectsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shortly_redirects_total",
			Help: "Total number of redirect requests by status",
		},
		[]string{"status"},
	)

	CacheHitsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "shortly_cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	CacheMissesTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "shortly_cache_misses_total",
			Help: "Total number of cache misses",
		},
	)

	RabbitMQConsumedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "shortly_rabbitmq_consumed_total",
			Help: "Total number of messages consumed from RabbitMQ",
		},
	)
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start).Seconds()
		pattern := r.Pattern
		if pattern == "" {
			pattern = r.URL.Path
		}

		HTTPRequestsTotal.WithLabelValues(r.Method, pattern, strconv.Itoa(wrapped.statusCode)).Inc()
		HTTPRequestDuration.WithLabelValues(r.Method, pattern).Observe(duration)
	})
}

func Handler() http.Handler {
	return promhttp.Handler()
}

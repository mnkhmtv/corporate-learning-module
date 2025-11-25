package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP метрики
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets, // 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10
		},
		[]string{"method", "endpoint"},
	)

	// Database метрики
	DbConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_active",
			Help: "Number of active database connections",
		},
	)

	DbConnectionsIdle = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_idle",
			Help: "Number of idle database connections",
		},
	)

	DbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
		},
		[]string{"operation"},
	)

	DbQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_queries_total",
			Help: "Total number of database queries",
		},
		[]string{"operation", "status"},
	)

	// Business метрики
	TrainingRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "training_requests_total",
			Help: "Total number of training requests",
		},
		[]string{"status"},
	)

	MentorsWorkload = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mentors_workload",
			Help: "Current workload of mentors (0-5 scale)",
		},
		[]string{"mentor_id", "mentor_name"},
	)

	LearningProcessesActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "learning_processes_active",
			Help: "Number of active learning processes",
		},
	)

	LearningProcessesCompleted = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "learning_processes_completed_total",
			Help: "Total number of completed learning processes",
		},
	)

	FeedbackRatingSum = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "feedback_rating_sum",
			Help: "Sum of all feedback ratings",
		},
	)

	FeedbackRatingCount = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "feedback_rating_count",
			Help: "Total number of feedback ratings",
		},
	)
)

// RecordHttpRequest records HTTP request metrics
func RecordHttpRequest(method, endpoint string, status int, duration time.Duration) {
	HttpRequestsTotal.WithLabelValues(method, endpoint, strconv.Itoa(status)).Inc()
	HttpRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

// RecordDbQuery records database query metrics
func RecordDbQuery(operation string, duration time.Duration, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}
	DbQueriesTotal.WithLabelValues(operation, status).Inc()
	DbQueryDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

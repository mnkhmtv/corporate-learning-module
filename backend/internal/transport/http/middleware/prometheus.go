package middleware

import (
	"time"

	"github.com/mnkhmtv/corporate-learning-module/backend/internal/pkg/metrics"

	"github.com/gin-gonic/gin"
)

// PrometheusMiddleware collects HTTP metrics for each request
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Record metrics
		duration := time.Since(start)
		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = "not_found"
		}

		metrics.RecordHttpRequest(
			c.Request.Method,
			endpoint,
			c.Writer.Status(),
			duration,
		)
	}
}

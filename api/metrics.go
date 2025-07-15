package api

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Metrics struct to store API metrics
type Metrics struct {
	TotalRequests   int64
	SuccessRequests int64
	ErrorRequests   int64
	TotalLatency    time.Duration
}

var apiMetrics = &Metrics{}

// MetricsMiddleware tracks API performance metrics
func MetricsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		// Process request
		ctx.Next()

		// Calculate metrics
		latency := time.Since(start)
		apiMetrics.TotalRequests++
		apiMetrics.TotalLatency += latency

		if ctx.Writer.Status() >= 200 && ctx.Writer.Status() < 400 {
			apiMetrics.SuccessRequests++
		} else {
			apiMetrics.ErrorRequests++
		}
	}
}

// GetMetrics returns current API metrics
func GetMetrics() *Metrics {
	return apiMetrics
}

// GetSuccessRate returns the success rate percentage
func GetSuccessRate() float64 {
	if apiMetrics.TotalRequests == 0 {
		return 0
	}
	return float64(apiMetrics.SuccessRequests) / float64(apiMetrics.TotalRequests) * 100
}

// GetAverageLatency returns average response time
func GetAverageLatency() time.Duration {
	if apiMetrics.TotalRequests == 0 {
		return 0
	}
	return apiMetrics.TotalLatency / time.Duration(apiMetrics.TotalRequests)
}

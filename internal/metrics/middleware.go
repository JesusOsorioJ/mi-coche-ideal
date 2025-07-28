package metrics

import (
	"time"

	"github.com/gin-gonic/gin"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		route := c.FullPath()
		if route == "" {
			route = "undefined"
		}

		RequestsTotal.WithLabelValues(c.Request.Method, route).Inc()
		RequestDuration.WithLabelValues(c.Request.Method, route).Observe(duration.Seconds())
	}
}

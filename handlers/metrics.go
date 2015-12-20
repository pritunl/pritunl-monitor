package handlers

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/gin-gonic/gin"
)

func metricsHandler() gin.HandlerFunc {
	prmHandler := prometheus.Handler()

	return func(c *gin.Context) {
		prmHandler.ServeHTTP(c.Writer, c.Request)
	}
}

package main

import (
	"bom_exporter/bom"
	"bom_exporter/bom/connection/ftp"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var requestDurations prometheus.Histogram

func metricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var h http.Handler

		timer := prometheus.NewTimer(requestDurations)
		defer timer.ObserveDuration()

		id := c.Query("id")
		if id == "" {
			h = promhttp.Handler()
		} else {
			registry := prometheus.NewPedanticRegistry()

			m := bom.New(ftp.New(id))

			err := m.RetrieveAndParse()
			if err != nil {
				log.Warnf("Failed to process: %s", err)
				c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("'%s' not found.", id)})
				return
			}
			registry.MustRegister(m)

			h = promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
		}

		h.ServeHTTP(c.Writer, c.Request)
	}
}

func init() {
	gin.SetMode(gin.ReleaseMode)
	requestDurations = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "bom",
		Subsystem: "metrics",
		Name:      "request_duration_seconds",
		Help:      "Histogram of request durations in seconds.",
		Buckets:   prometheus.DefBuckets})
	prometheus.MustRegister(requestDurations)
}

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.GET("/metrics", metricsHandler())

	r.Run()
}

package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pmoncadaisla/go-journey/pkg/domain"
	"github.com/pmoncadaisla/go-journey/pkg/journey"
	"github.com/pmoncadaisla/go-journey/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

func webserver() {

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	r.GET("/metrics", gin.WrapH(prometheus.Handler()))

	r.POST("/journey/receive", func(c *gin.Context) {
		var j domain.Journey
		if err := c.BindJSON(&j); err == nil {
			if metricsService.GetCounterValue(metrics.JOURNEYS_RECEIVED.String())+1 != j.ID {
				metricsService.CounterInc(metrics.HTTP_400_COUNT.String())
				c.JSON(400, gin.H{
					"status":              "bad_request",
					"highest_received_id": metricsService.GetCounterValue(metrics.JOURNEYS_RECEIVED.String()),
					"received_id":         j.ID,
					"message":             "received_id must be higher than highest_received_id",
				})
			} else {
				journey.Receive(j.ID, j.Time*time.Millisecond, finished)
				metricsService.CounterInc(metrics.HTTP_400_COUNT.String())
				c.JSON(200, gin.H{
					"status":       "received",
					"journey_id":   j.ID,
					"journey_time": j.Time,
				})
			}

		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	r.GET("/journey/next", func(c *gin.Context) {
		metricsService.CounterInc(metrics.HTTP_200_COUNT.String())
		c.JSON(200, gin.H{
			"id": metricsService.GetCounterValue(metrics.JOURNEYS_RECEIVED.String()) + 1,
		})
	})
	r.Run()

}

package router

import (
	"ginkvtest/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type result struct {
	value    string
	err      error
	duration float64
}

var (
	auroraResult result
	dynamoResult result
	redisResult  result
)

func SetupRoutes(r *gin.Engine, aurora service.DatabaseService, dynamo service.DatabaseService, redis service.DatabaseService) {
	r.GET("/benchmark/:key", func(c *gin.Context) {
		key := c.Param("key")
		g, _ := errgroup.WithContext(c.Request.Context()) // Use Gin's request context for cancellation

		// Aurora Performance
		g.Go(func() error {
			start := time.Now()
			value, err := aurora.GetValueByKey(key)
			auroraResult = result{
				value:    value,
				err:      err,
				duration: float64(time.Since(start).Microseconds()), // Convert to Microseconds
			}
			return err
		})

		// DynamoDB Performance
		g.Go(func() error {
			start := time.Now()
			value, err := dynamo.GetValueByKey(key)
			dynamoResult = result{
				value:    value,
				err:      err,
				duration: float64(time.Since(start).Microseconds()), // Convert to Microseconds
			}
			return err
		})

		// Redis Performance
		g.Go(func() error {
			start := time.Now()
			value, err := redis.GetValueByKey(key)
			redisResult = result{
				value:    value,
				err:      err,
				duration: float64(time.Since(start).Microseconds()), // Convert to Microseconds
			}
			return err
		})

		// Wait for all goroutines to complete
		if err := g.Wait(); err != nil {
			// If any error occurs, return a 500 status
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return the results
		c.JSON(http.StatusOK, gin.H{
			"aurora": gin.H{
				"value":    auroraResult.value,
				"error":    auroraResult.err,
				"duration": auroraResult.duration,
			},
			"dynamo": gin.H{
				"value":    dynamoResult.value,
				"error":    dynamoResult.err,
				"duration": dynamoResult.duration,
			},
			"redis": gin.H{
				"value":    redisResult.value,
				"error":    redisResult.err,
				"duration": redisResult.duration,
			},
		})
	})
}

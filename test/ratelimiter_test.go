package test

import (
	middleware "first-api/internal/middlewares"
	"first-api/internal/ratelimiter"
	"first-api/pkg/cache"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestTokenBucketMiddleware is a test function to verify the behavior of the rate limiter middleware.
func TestTokenBucketMiddleware(t *testing.T) {
	rateLimiterCache := cache.GetRateLimiterCache()

	// Create a new TokenBucket and set the initial token count in the cache

	tb := ratelimiter.NewTokenBucket(100, 1, time.Now())

	// Create a new Gin router and apply the middleware
	router := gin.New()
	router.Use(func(ctx *gin.Context) {
		middleware.RateLimitMiddleware(ctx, rateLimiterCache, tb)
	})

	// Define a test route
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Success!")
	})

	// Perform test requests
	for i := 1; i <= 5; i++ {
		req, _ := http.NewRequest("GET", "/test", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		// Check if the response code matches the expected status code
		if resp.Code != http.StatusOK {
			assert.Equal(t, http.StatusTooManyRequests, resp.Code)
		} else {
			assert.Equal(t, http.StatusOK, resp.Code)
		}
	}
}

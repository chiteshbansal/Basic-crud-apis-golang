package test

import (
	middleware "first-api/internal/middlewares"
	"first-api/internal/ratelimiter"
	"first-api/pkg/cache"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	// Import your middleware package
)

func TestTokenBucketMiddleware(t *testing.T) {
	// Create a test Redis client
	rateLimiterCache :=
		cache.NewRedisCache("localhost:6379", 1, 1000)

	ratelimiter.NewTokenBucket(100, 1, time.Now())

	// Create a new Gin router and apply the middleware
	router := gin.New()
	router.Use(func(ctx *gin.Context) {
		fmt.Println("using the ratelimiter middleware")
		middleware.RateLimitMiddleware(ctx, rateLimiterCache, ratelimiter.GetTokenBucket())
	})
	// Define a test route
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Success!")
	})

	// Perform test requests
	for i := 1; i <= 1000; i++ {
		req, _ := http.NewRequest("GET", "/test", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		fmt.Println("seding request no :", i)
		if resp.Code != http.StatusOK {
			fmt.Println("Expected status code 200, but got ", resp.Code)
			assert.Equal(t, http.StatusTooManyRequests, resp.Code)

		}
	}
}

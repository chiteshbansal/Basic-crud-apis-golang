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
	"github.com/stretchr/testify/mock"
)

// RateLimiterCache is a mock implementation of the cache.UserCache interface.
// It is used for testing the rate limiter middleware.
type RateLimiterCache struct {
	cache.UserCache
	mock.Mock
}

// Set is a mock implementation of the cache.UserCache Set method.
func (m *RateLimiterCache) Set(key string, value interface{}, exp *time.Duration) error {
	args := m.Called(key, value, exp)
	return args.Error(0)
}

// Get is a mock implementation of the cache.UserCache Get method.
func (m *RateLimiterCache) Get(key string) (interface{}, error) {
	args := m.Called(key)
	return args.Get(0), args.Error(1)
}

// TestTokenBucketMiddleware is a test function to verify the behavior of the rate limiter middleware.
func TestTokenBucketMiddleware(t *testing.T) {
	// Create a test Redis client
	rateLimiterCache := &RateLimiterCache{}

	// Create a new TokenBucket and set the initial token count in the cache
	ratelimiter.NewTokenBucket(100, 1, time.Now())
	tb := ratelimiter.GetTokenBucket()

	// Set up the mock expectations for the cache calls
	rateLimiterCache.On("Set", "token_bucket", 99, mock.AnythingOfType("*time.Duration")).Return(nil)
	rateLimiterCache.On("Get", "token_bucket").Return(tb.Limit+1, nil)

	// Create a new Gin router and apply the middleware
	router := gin.New()
	router.Use(func(ctx *gin.Context) {
		fmt.Println("Using the ratelimiter middleware")
		middleware.RateLimitMiddleware(ctx, rateLimiterCache, ratelimiter.GetTokenBucket())
	})

	// Define a test route
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Success!")
	})

	// Perform test requests
	for i := 1; i <= 500; i++ {
		req, _ := http.NewRequest("GET", "/test", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		fmt.Println("Sending request no:", i)

		// Check if the response code matches the expected status code
		if resp.Code != http.StatusOK {
			assert.Equal(t, http.StatusTooManyRequests, resp.Code)
		} else {
			assert.Equal(t, http.StatusOK, resp.Code)
		}
	}
}

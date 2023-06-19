package middleware

import (
	"first-api/internal/ratelimiter"
	"first-api/pkg/cache"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware(c *gin.Context, redisClient cache.UserCache, tb *ratelimiter.TokenBucket) {
	// Check if the Token Bucket has enough tokens
	if tb.Take(c, redisClient) {
		// Perform the request
		c.Next()
	} else {
		// Rate limit reached, reject the request
		c.AbortWithStatus(http.StatusTooManyRequests)
	}
}

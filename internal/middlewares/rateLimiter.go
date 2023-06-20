package middleware

import (
	"first-api/internal/ratelimiter"
	"first-api/pkg/cache"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware(c *gin.Context, redisClient cache.UserCache, tb *ratelimiter.TokenBucket) {
	// Check if the Token Bucket has enough tokens
	if tb.Take(c, redisClient) {
		// Perform the request
		// check for per client basis request
		clientID :=GetClientIdentity(c)
		clientTB := ratelimiter.GetClientBucket(clientID, "user",redisClient)
		fmt.Println("clientID:", clientID)
		if !clientTB.IsRequestAllowed(1) {
			fmt.Println("too many request by ",clientID)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "Too many requests , try again after sometimes!",
			})
			return
		}
		c.Next()
	} else {
		// Rate limit reached, reject the request
		c.AbortWithStatus(http.StatusTooManyRequests)
	}
}

func GetClientIdentity(c *gin.Context) string {
	ip := c.ClientIP()
	url := c.Request.URL.Path
	data := fmt.Sprintf("%s-%s", ip, url)
	return data
}

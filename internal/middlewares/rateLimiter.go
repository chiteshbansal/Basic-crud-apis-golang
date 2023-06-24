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
	if tb.IsTotalRequestAllowed(redisClient) {
		// Perform the request
		// check for per client basis request
		clientID := GetClientIdentity(c)
		clientTB := ratelimiter.GetClientBucket(clientID, "user", redisClient)
		if clientTB == nil {
			fmt.Println("failed to retreive client bucket from redis server")
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"message": "Something went wrong try again later ",
			})
			return
		}
		if !clientTB.IsClientRequestAllowed(1) {
			fmt.Println("too many request by ", clientID)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "Too many requests , try again after sometimes!",
			})
			return
		}
		err := ratelimiter.UpdateClientBucketInRedis(redisClient)
		if err != nil {
			fmt.Println("Error:", err)
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"message": err,
			})
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
	method :=c.Request.Method
	data := fmt.Sprintf("%s-%s-%s", ip, url,method)
	return data
}


package ratelimiter

import (
	"first-api/pkg/cache"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	TB *TokenBucket
)

type TokenBucket struct {
	limit    int
	tokens   int
	rate     int
	mu       sync.Mutex
	duration time.Time
}

// NewTokenBucket creates a new TokenBucket with the specified limit, rate, and duration.
func NewTokenBucket(limit, rate int, duration time.Time) {
	TB = &TokenBucket{
		limit:    limit,
		tokens:   limit,
		rate:     rate,
		duration: duration,
	}

}

func GetTokenBucket() *TokenBucket {
	return TB
}

// Take attempts to take a token from the TokenBucket.
// Returns true if a token is taken, false otherwise.
func (tb *TokenBucket) Take(c *gin.Context, redisClient cache.UserCache) bool {
	TB.mu.Lock()
	defer TB.mu.Unlock()

	// Key to identify the token bucket in Redis
	key := "token_bucket"

	now := time.Now()

	// Calculate the number of tokens to add since the last request
	fmt.Println(now, TB.duration)
	elapsed := now.Sub(TB.duration)
	fmt.Println(elapsed.Seconds(), elapsed, time.Second, (float64(elapsed.Nanoseconds()) / float64(time.Second)))
	tokensToAdd := int(float64(elapsed.Nanoseconds()) / float64(time.Second) * float64(TB.rate))

	// Refill tokens up to the limit
	fmt.Println("tokens add", tokensToAdd)
	TB.tokens += tokensToAdd
	if TB.tokens > TB.limit {
		TB.tokens = TB.limit
	}

	// Check if there re enough tokens to take
	if TB.tokens > 0 {
		fmt.Println("under limit Passing the req")
		TB.tokens--
		TB.duration = now
		return true

		// Check the current token count in Redis
		resultInterface, err := redisClient.Get(key)
		if err != nil || resultInterface == nil {
			fmt.Println("Failed to retrieve token count from Redis:", err)
			return false
		}
		result := resultInterface.(int)
		// Update the token count in Redis
		result -= 1
		if result >= 0 {
			err := redisClient.Set(key, result, nil)
			if err != nil {
				fmt.Println("Failed to update token count in Redis:", err)
			}
			return true
		}
	} else {
		fmt.Println("Request denied")
		return false
	}

	return false

}

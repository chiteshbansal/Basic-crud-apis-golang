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
	Limit    int
	tokens   int
	rate     int
	mu       sync.Mutex
	duration time.Time
}

// NewTokenBucket creates a new TokenBucket with the specified limit, rate, and duration.
func NewTokenBucket(limit, rate int, duration time.Time) {
	TB = &TokenBucket{
		Limit:    limit,
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
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// Key to identify the token bucket in Redis
	key := "token_bucket"

	now := time.Now()

	// Calculate the number of tokens to add since the last request
	elapsed := now.Sub(tb.duration)
	tokensToAdd := int(float64(elapsed.Nanoseconds()) / float64(time.Second) * float64(tb.rate))

	// Fetch the current token count from Redis
	resultInterface, err := redisClient.Get(key)
	if err != nil {
		fmt.Println("Failed to retrieve token count from Redis:", err)
		return false
	}

	var result int
	result, ok := resultInterface.(int)
	if !ok {
		// Handle the case when the underlying value is not an int
		// Convert the float64 value to int before using it
		result = int(resultInterface.(float64))
	}

	// Add the new tokens to the fetched token count
	result += tokensToAdd

	if result > tb.Limit {
		result = tb.Limit
	}

	// Check if there are enough tokens to take
	if result > 0 {
		// Decrement the token count in Redis
		result--
		err := redisClient.Set(key, result, nil)
		if err != nil {
			fmt.Println("Failed to update token count in Redis:", err)
		}

		// Update the local token count and duration
		tb.tokens = result
		tb.duration = now

		return true
	} else {
		fmt.Println("Request denied")
		return false
	}
}

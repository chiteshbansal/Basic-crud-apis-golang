package ratelimiter

import (
	"first-api/pkg/cache"
	"fmt"
	"math"
	"time"
)

var (
	TB *TokenBucket
)

type TokenBucket struct {
	Limit    int
	Tokens   int
	Rate     int
	Duration time.Time
}

// NewTokenBucket creates a new TokenBucket with the specified limit, Rate, and duration.
func NewTokenBucket(limit, rate int, duration time.Time) *TokenBucket {
	return &TokenBucket{
		Limit:    limit,
		Tokens:   limit,
		Rate:     rate,
		Duration: duration,
	}
}

// Take attempts to take a token from the TokenBucket.
// Returns true if a token is taken, false otherwise.
func (tb *TokenBucket) IsTotalRequestAllowed(redisClient cache.UserCache) bool {

	key := "token_bucket"
	mutex := redisClient.GetMutex(key)
	mutex.Lock()
	defer mutex.Unlock()

	// Key to identify the token bucket in Redis

	now := time.Now()

	// Calculate the number of tokens to add since the last request
	elapsed := now.Sub(tb.Duration)
	tokensToAdd := int(float64(elapsed.Nanoseconds()) / float64(time.Second) * float64(tb.Rate))

	// Fetch the current token count from Redis
	resultInterface, err := redisClient.Get(key)
	if err != nil {
		fmt.Println("Failed to retrieve token count from Redis:", err)
		err := redisClient.Set(key, tb.Limit, nil)
		if err != nil {
			fmt.Println("Failed to update token count in Redis:", err)
		}

		// Update the local token count and Duration
		tb.Tokens = tb.Limit - 1
		tb.Duration = now

		return true
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

		// Update the local token count and Duration
		tb.Tokens = result
		tb.Duration = now

		return true
	} else {
		fmt.Println("Request denied")
		return false
	}
}

func (tb *TokenBucket) refill() {
	now := time.Now()
	end := time.Since(tb.Duration)
	tokensToBeAdded := (end.Nanoseconds() * int64(tb.Rate)) / 1000000000
	tb.Tokens = int(math.Min(float64(tb.Tokens+int(tokensToBeAdded)), float64(tb.Limit)))
	tb.Duration = now
}

func (tb *TokenBucket) IsClientRequestAllowed(tokens int) bool {

	tb.refill()
	fmt.Println("client tokens", tb.Tokens)
	if tb.Tokens >= tokens {
		tb.Tokens = tb.Tokens - tokens
		return true
	}
	return false
}

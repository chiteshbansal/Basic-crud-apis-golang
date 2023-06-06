// Package cache provides functionalities to interact with a Redis cache.
package cache

import (
	"context"
	"encoding/json"
	model "first-api/api/Models"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// redisCache is a struct that holds the Redis host, database, and expiration details.
type redisCache struct {
	host    string
	db      int
	expires time.Duration
}

// NewRedisCache returns a new instance of redisCache.
func NewRedisCache(host string, db int, exp time.Duration) UserCache {
	return &redisCache{
		host:    host,
		db:      db,
		expires: exp,
	}
}

// getClient establishes a new connection with the Redis client.
func (cache *redisCache) getClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})

	ctx := context.Background()
	pong, err := client.Ping(ctx).Result()

	// If an error occurred while pinging, log it and proceed.
	if err != nil {
		fmt.Println("Error pinging Redis server:", err)
	}
	// Print the response from the Redis server.
	fmt.Println(pong)

	return client
}

// Set stores the user object in the Redis cache against the provided key.
func (cache *redisCache) Set(key string, user *model.User) {
	fmt.Println("Caching the query: ", key, user)

	// Obtain the Redis client.
	client := cache.getClient()

	// Marshal user object into JSON.
	json, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	
	// Set the key-value pair in the Redis cache with the expiry duration.
	client.Set(ctx, key, json, cache.expires*time.Second)
}

// Get retrieves the user object from the Redis cache for the given key.
func (cache *redisCache) Get(key string) *model.User {
	// Obtain the Redis client.
	client := cache.getClient()
	
	ctx := context.Background()

	// Get the value for the given key.
	val, err := client.Get(ctx, key).Result()
	if err != nil || val == "" {
		// Return nil if an error occurred or if the key does not exist.
		return nil
	}

	post := model.User{}

	// Unmarshal the JSON value into the user struct.
	err = json.Unmarshal([]byte(val), &post)
	if err != nil {
		panic(err)
	}

	// Return the pointer to the user object.
	return &post
}

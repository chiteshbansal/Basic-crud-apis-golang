// Package cache provides functionalities to interact with a Redis cache.
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redsync/redsync/v4"
	goredislib "github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

// redisCache is a struct that holds the Redis host, database, and expiration details.
type redisCache struct {
	host    string
	db      int
	expires time.Duration
}

// NewRedisCache returns a new instance of redisCache.
func NewRedisCache(host string, db int, exp time.Duration) *redisCache {
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

// Set stores the object in the Redis cache against the provided key.
func (cache *redisCache) Set(key string, value interface{}, exp *time.Duration) error {
	fmt.Println("Caching the query: ", key, value)

	// Obtain the Redis client.
	client := cache.getClient()

	// Marshal object into JSON.
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// Set the key-value pair in the Redis cache with the expiry duration.
	if exp == nil {
		err = client.Set(ctx, key, json, cache.expires*time.Second).Err()
	} else {
		err = client.Set(ctx, key, json, *exp).Err()
	}

	if err != nil {
		return err
	}

	return nil
}

// Get retrieves the object from the Redis cache for the given key.
func (cache *redisCache) Get(key string) (value interface{}, err error) {
	// Obtain the Redis client.
	client := cache.getClient()

	ctx := context.Background()

	// Get the value for the given key.
	val, err := client.Get(ctx, key).Result()
	if err != nil || val == "" {
		// Return nil if an error occurred or if the key does not exist.
		return nil, err
	}

	// Unmarshal the JSON value into the interface{}.
	err = json.Unmarshal([]byte(val), &value)
	if err != nil {
		return nil, err
	}

	// Return the interface{}.
	return value, nil
}

func (cache *redisCache) GetMutex(mutexname string) *redsync.Mutex {
	redisClient := cache.getClient()
	// Create a Redis pool
	pool := goredislib.NewPool(redisClient)

	// Create a Redsync instance using the Redis pool
	rs := redsync.New(pool)
	mutex := rs.NewMutex(mutexname)

	return mutex
}

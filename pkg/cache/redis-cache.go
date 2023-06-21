package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

// redisCache is a struct that holds the Redis host, database, and expiration details.
type redisCache struct {
	pool    *redis.Pool
	expires time.Duration
}

// NewRedisCache returns a new instance of redisCache.
func NewRedisCache(host string, db int, exp time.Duration) *redisCache {
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", host)
			if err != nil {
				return nil, err
			}
			if db != 0 {
				_, err := conn.Do("SELECT", db)
				if err != nil {
					conn.Close()
					return nil, err
				}
			}
			return conn, nil
		},
	}

	return &redisCache{
		pool:    pool,
		expires: exp,
	}
}

// Set stores the object in the Redis cache against the provided key.
func (cache *redisCache) Set(key string, value interface{}, exp *time.Duration) error {
	fmt.Println("Caching the query:", key, value)

	// Obtain a connection from the Redis pool.
	conn := cache.pool.Get()
	defer conn.Close()

	// Marshal object into JSON.
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Set the key-value pair in the Redis cache with the expiry duration.
	if exp == nil {
		fmt.Println(cache.expires*time.Second, cache.expires, (cache.expires.Seconds()))
		expiration := cache.expires * time.Second
		_, err = conn.Do("SETEX", key, expiration.Seconds(), jsonData)
	} else {
		fmt.Println((exp.Seconds()))
		_, err = conn.Do("SETEX", key, int(exp.Seconds()), jsonData)
	}

	if err != nil {
		return err
	}

	return nil
}

// Get retrieves the object from the Redis cache for the given key.
func (cache *redisCache) Get(key string) (value interface{}, err error) {
	// Obtain a connection from the Redis pool.
	conn := cache.pool.Get()
	defer conn.Close()

	// Get the value for the given key.
	val, err := redis.Bytes(conn.Do("GET", key))
	if err != nil || val == nil {
		// Return nil if an error occurred or if the key does not exist.
		return nil, err
	}

	// Unmarshal the JSON value into the interface{}.
	err = json.Unmarshal(val, &value)
	if err != nil {
		return nil, err
	}

	// Return the interface{}.
	return value, nil
}

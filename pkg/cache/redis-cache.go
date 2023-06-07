package cache

import (
	"context"
	"encoding/json"
	model "first-api/internal/models"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	host    string
	db      int
	expires time.Duration
}

func NewRedisCache(host string, db int, exp time.Duration) UserCache {
	return &redisCache{
		host:    host,
		db:      db,
		expires: exp,
	}
}

func (cache *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}

func (cache *redisCache) Set(key string, user *model.User) {
	fmt.Println("caching the query ", key, user)
	client := cache.getClient()
	json, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	// client.Set(key, json, cache.expires*time.Second)
	client.Set(ctx, key, json, cache.expires*time.Second)
}

func (cache *redisCache) Get(key string) *model.User {
	client := cache.getClient()
	ctx := context.Background()
	fmt.Println("client :", client)
	// return &model.User{}
	val, err := client.Get(ctx, key).Result()
	fmt.Println("value:", val)
	if err != nil {
		return nil
	}
	if val == "" {
		return nil
	}

	post := model.User{}
	err = json.Unmarshal([]byte(val), &post)
	if err != nil {
		panic(err)
	}

	return &post
}

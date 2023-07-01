package ratelimiter

import (

	"encoding/json"
	"errors"
	"first-api/pkg/cache"
	"fmt"
	"time"
)

var clientBucketMap = make(map[string]*TokenBucket)

type Rule struct {
	MaxTokens int
	Rate      int
}

func GetClientBucket(identifier string, usertype string, redisClient cache.UserCache) *TokenBucket {

	key := "client_bucket"
	// Fetch the current client bucket map from Redis
	mutex := redisClient.GetMutex(key)
	mutex.Lock()
	defer mutex.Unlock()
	result, err := redisClient.Get(key)
	if err != nil {
		fmt.Println("Failed to retrieve client bucket data from Redis:", err)
		// redis don't have anything for client_bucket
		clientBucketMap = make(map[string]*TokenBucket)
	} else {
		if jsonString, ok := result.(string); ok {
			err := json.Unmarshal([]byte(jsonString), &clientBucketMap)
			if err != nil {
				fmt.Println("Error unmarshaling JSON:", err)
				return nil
			}
		} else {
			fmt.Println("Invalid JSON format:", result)
			return nil
		}
	}

	if clientBucketMap[identifier] == nil {
		clientBucketMap[identifier] = NewTokenBucket(rulesMap[usertype].MaxTokens, rulesMap[usertype].Rate, time.Now())
	}

	err = UpdateClientBucketInRedis(redisClient)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return clientBucketMap[identifier]
}

func UpdateClientBucketInRedis(redisClient cache.UserCache) error {
	// // writing to redis again to update
	jsonData, err := json.Marshal(clientBucketMap)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return errors.New("Error marshaling JSON:")
	}

	// // Store JSON bytes in Redis
	err = redisClient.Set("client_bucket", string(jsonData), nil)
	if err != nil {
		fmt.Println("Error storing data in Redis:", err)
		return errors.New("Error storing data in Redis")

	}
	return nil
}

var rulesMap = map[string]Rule{
	"user":  {MaxTokens: 1, Rate: 1},
	"admin": {MaxTokens: 20, Rate: 1},
}

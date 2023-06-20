package ratelimiter

import (
	// "encoding/json"
	//"first-api/internal/utils"
	"first-api/pkg/cache"
	// "fmt"
	"time"
)

var clientBucketMap = make(map[string]*TokenBucket)

type Rule struct {
	MaxTokens int
	Rate      int
}

func GetClientBucket(identifier string, usertype string, redisClient cache.UserCache) *TokenBucket {

	// Fetch the current client bucket map from Redis

	// result, err := redisClient.Get("client_bucket")
	// if err != nil {
	// 	fmt.Println("Failed to retrieve client bucket data from Redis:", err)
	// 	// redis don't have anything for client_bucket
	// 	clientBucketMap = make(map[string]*TokenBucket)
	// } else {
	// 	fmt.Println("result", result)
	// 	resultBytes, err := utils.InterfaceToBytes(result)
	// 	if err != nil {
	// 		fmt.Println("Error converting data to bytes:", err)
	// 		return nil
	// 	}
	// 	err = json.Unmarshal(resultBytes, &clientBucketMap)
	// 	if err != nil {
	// 		fmt.Println("Error unmarshaling JSON:", err)
	// 		return nil
	// 	}
	// }

	if clientBucketMap[identifier] == nil {
		clientBucketMap[identifier] = NewTokenBucket(rulesMap[usertype].MaxTokens, rulesMap[usertype].Rate, time.Now())
	}

	// // writing to redis again to update
	// jsonData, err := json.Marshal(clientBucketMap)
	// if err != nil {
	// 	fmt.Println("Error marshaling JSON:", err)
	// 	return nil
	// }

	// // Store JSON bytes in Redis
	// err = redisClient.Set("client_bucket", string(jsonData), nil)
	// if err != nil {
	// 	fmt.Println("Error storing data in Redis:", err)
	// 	return nil
	// }
	return clientBucketMap[identifier]
}

var rulesMap = map[string]Rule{
	"user":  {MaxTokens: 1, Rate: 1},
	"admin": {MaxTokens: 20, Rate: 1},
}

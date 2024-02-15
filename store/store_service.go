package store

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

// Define struct wrapper around raw Redis client

type StorageService struct {
	redisClient *redis.Client
}

// Top level declarations for storeService and redis context
var (
	storeService = &StorageService{}
	ctx          = context.Background()
)

// Set expiration for each  key, value pair in redis
const storeExpiration = 12 * time.Hour

func InitializeStore() *StorageService {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	// Check if we can ping the Redis service
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Println(fmt.Sprintf("Error init Redis: %v\n", err))
	}
	fmt.Printf("Redis started successfully: pong msg: = %s", pong)
	storeService.redisClient = redisClient
	return storeService
}

func SaveUrlMapping(shortUrl string, originalUrl string) {
	err := storeService.redisClient.Set(ctx, shortUrl, originalUrl, storeExpiration).Err()
	if err != nil {
		log.Println(fmt.Sprintf("Failed saving key url. Error: %v - shortUrl: %s - originalUrl: %s\n", err, shortUrl, originalUrl))
	}
}

func GetUrlMapping(shortUrl string) string {
	res, err := storeService.redisClient.Get(ctx, shortUrl).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to get url. Error: %v - shortUrl: %s", err, shortUrl))
	}
	return res
}

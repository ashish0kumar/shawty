package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient() *redis.Client {
	fmt.Println("Connecting to redis server on: ", os.Getenv("REDIS_HOST"))

	// redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	return rdb
}

func SetKey(ctx *context.Context, rdb *redis.Client, key string, value string, ttl int) {

	// sets the key value pair in redis
	// uses the context defined in main by reference and a TTL of 0 (no expiration)

	fmt.Println("Setting key: ", key, "to", value)
	rdb.Set(*ctx, key, value, 0)
	fmt.Println("The key", key, "set to", value, "successfully")
}

func GetLongURL(ctx *context.Context, rdb *redis.Client, shortURL string) (string, error) {
	longURL, err := rdb.Get(*ctx, shortURL).Result()

	if err == redis.Nil {
		return "", fmt.Errorf("shortened URL not found")
	} else if err != nil {
		return "", fmt.Errorf("error fetching long URL from Redis: %v", err)
	}

	return longURL, nil
}

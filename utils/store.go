package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

// Prefix keys for different types of data
const (
	ShortURLPrefix = "short:" // For short URL to long URL mapping
	LongURLPrefix  = "long:"  // For long URL to short URL mapping
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

// SetKey stores both mappings: short->long and long->short
func SetKey(ctx *context.Context, rdb *redis.Client, shortURL string, longURL string, ttl int) {
	// Set short URL -> long URL mapping
	fmt.Println("Setting key: ", ShortURLPrefix+shortURL, "to", longURL)
	rdb.Set(*ctx, ShortURLPrefix+shortURL, longURL, 0)

	// Set long URL -> short URL mapping for future lookups
	fmt.Println("Setting reverse key: ", LongURLPrefix+longURL, "to", shortURL)
	rdb.Set(*ctx, LongURLPrefix+longURL, shortURL, 0)

	fmt.Println("URL mappings set successfully")
}

// GetLongURL retrieves the long URL for a given short code
func GetLongURL(ctx *context.Context, rdb *redis.Client, shortURL string) (string, error) {
	longURL, err := rdb.Get(*ctx, ShortURLPrefix+shortURL).Result()

	if err == redis.Nil {
		return "", fmt.Errorf("shortened URL not found")
	} else if err != nil {
		return "", fmt.Errorf("error fetching long URL from Redis: %v", err)
	}

	return longURL, nil
}

// GetExistingShortURL checks if a URL has already been shortened
func GetExistingShortURL(ctx *context.Context, rdb *redis.Client, longURL string) (string, error) {
	shortURL, err := rdb.Get(*ctx, LongURLPrefix+longURL).Result()

	if err == redis.Nil {
		return "", nil // No error, just not found
	} else if err != nil {
		return "", fmt.Errorf("error checking existing URL: %v", err)
	}

	return shortURL, nil
}

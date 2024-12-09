package service

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
)

type RedisService struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisService initializes the Redis client and ensures the default key exists.
func NewRedisService() (*RedisService, error) {
	// Set up the Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: "", // No password for local setup
		DB:       0,  // Use default DB
	})

	ctx := context.Background()
	service := &RedisService{client: client, ctx: ctx}

	// Ensure the default key exists
	service.ensureDefaultKeyExists("test_key", "test_value")

	return service, nil
}

// ensureDefaultKeyExists checks if a key exists, and inserts it with a default value if needed.
func (s *RedisService) ensureDefaultKeyExists(key string, defaultValue string) {
	val, err := s.client.Get(s.ctx, key).Result()

	if err == redis.Nil {
		// Key does not exist, set it
		fmt.Printf("Key '%s' not found, inserting default value...\n", key)
		err = s.client.Set(s.ctx, key, defaultValue, 0).Err()
		if err != nil {
			fmt.Println("Failed to insert default key-value:", err)
			return
		}
		fmt.Printf("Key '%s' initialized with value: '%s'\n", key, defaultValue)
	} else if err != nil {
		// Some other error occurred
		fmt.Println("Error checking key:", err)
	} else {
		// Key already exists
		fmt.Printf("Key '%s' already exists with value: '%s'\n", key, val)
	}
}

// GetValueByKey retrieves the value for a given key.
func (s *RedisService) GetValueByKey(key string) (string, error) {
	val, err := s.client.Get(s.ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key not found")
	} else if err != nil {
		return "", err
	}
	return val, nil
}

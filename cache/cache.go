package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type GoMovieCache interface {
	Get(key string, data interface{}) (interface{}, error)
	Set(key string, data interface{}) error
}

type goMovieCache struct {
	redisClient *redis.Client
}

func NewGoMovieCache(redisClient *redis.Client) GoMovieCache {
	return &goMovieCache{redisClient: redisClient}
}

func (c *goMovieCache) Get(key string, data interface{}) (interface{}, error) {
	val, err := c.redisClient.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(val, data)
	if err != nil {
		log.Printf("error unmarshalling data: %v", err)
		return nil, fmt.Errorf("error unmarshalling data: %v", err)
	}
	return data, nil
}

func (c *goMovieCache) Set(key string, data interface{}) error {
	marshal, err := json.Marshal(data)
	if err != nil {
		log.Printf("error marshalling data for redis: %v", err)
		return fmt.Errorf("error marshalling data for redis: %v", err)
	}
	return c.redisClient.Set(context.Background(), key, string(marshal), 20*time.Second).Err()
}

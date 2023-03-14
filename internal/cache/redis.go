package cache

import (
	"context"
	"fmt"
	"github.com/iBoBoTi/go-movie-api/internal/config"
	"github.com/redis/go-redis/v9"
	"log"
)

func NewRedisClient(config *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", config.RedisHost, config.RedisPort),
		Password: fmt.Sprintf("%v", config.RedisPassword),
		DB:       0, // use default DB
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(pong)
	return client
}

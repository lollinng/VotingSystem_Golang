package redis_db

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var Clients *redis.Client

func init() {
	Clients = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})
}

func Ping() error {
	pong, err := Clients.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	log.Printf("Connected to Redis! Ping response: %s", pong)
	return nil
}

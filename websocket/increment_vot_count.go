package websocket

import (
	"context"
	"sync"

	"govoting/redis_db"
)

var mu sync.RWMutex

func incrementVoteCount(option string) error {
	mu.Lock()
	defer mu.Unlock()

	// Increment vote count in Redis
	_, err := redis_db.Clients.HIncrBy(context.Background(), "votes", option, 1).Result()
	if err != nil {
		return err
	}

	return nil
}

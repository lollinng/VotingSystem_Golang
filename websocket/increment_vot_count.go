package websocket

import (
	"context"
	"sync"

	"govoting/redis_db"

	"github.com/gorilla/websocket"
)

var mu sync.RWMutex

func incrementVoteCount(option string, conn *websocket.Conn) error {
	mu.Lock()

	defer mu.Unlock()
	// if the client key exist which means he has voted

	// Increment vote count in Redis
	_, err := redis_db.Clients.HIncrBy(context.Background(), "votes", option, 1).Result()
	if err != nil {
		return err
	}

	return nil
}

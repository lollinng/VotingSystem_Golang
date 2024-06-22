package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	mu          sync.RWMutex // Mutex for locking votes map
	clients     = make(map[*websocket.Conn]bool)
)

func init() {
	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})
}

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Vote struct {
	Option string `json:"option"`
}

type VoteCount struct {
	Option string `json:"option"`
	Count  int    `json:"count"`
}

func main() {

	// Ping the Redis server to check connectivity
	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Printf("Connected to Redis! Ping response: %s", pong)
	// Set up HTTP routes

	http.HandleFunc("/ws", wsHandler)
	log.Fatal(http.ListenAndServe(":8040", nil))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		delete(clients, conn)
		conn.Close()
	}()

	clients[conn] = true

	handleClient(conn)
}

func handleClient(conn *websocket.Conn) {
	defer func() {
		delete(clients, conn)
		conn.Close()
	}()

	for {
		var vote Vote
		err := conn.ReadJSON(&vote)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading JSON: %v", err)
			}
			break
		}

		// Update vote count in Redis
		err = incrementVoteCount(vote.Option)
		if err != nil {
			log.Println("Error incrementing vote count in Redis:", err)
			break
		}

		// Broadcast updated vote count to all clients
		broadcastVoteCount()
	}
}

func incrementVoteCount(option string) error {
	mu.Lock()
	defer mu.Unlock()

	// Increment vote count in Redis
	_, err := redisClient.HIncrBy(context.Background(), "votes", option, 1).Result()
	if err != nil {
		return err
	}

	return nil
}

func broadcastVoteCount() {
	mu.RLock()
	defer mu.RUnlock()

	// Retrieve current vote counts from Redis
	votesMap, err := redisClient.HGetAll(context.Background(), "votes").Result()
	if err != nil {
		log.Println("Error retrieving vote counts from Redis:", err)
		return
	}
	fmt.Printf("votesMap: %+v\n", votesMap)

	for client := range clients {
		for option, countStr := range votesMap {
			count, _ := strconv.Atoi(countStr)
			voteCount := VoteCount{Option: option, Count: count}
			err := client.WriteJSON(voteCount)
			if err != nil {
				log.Println("Error writing JSON to client:", err)
				delete(clients, client)
				client.Close()
			}
		}
	}
}

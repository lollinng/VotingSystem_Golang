package websocket

import (
	"context"
	"log"
	"strconv"

	"govoting/models"
	"govoting/redis_db"
)

func broadcastVoteCount() {
	votesMap, err := redis_db.Clients.HGetAll(context.Background(), "votes").Result()
	if err != nil {
		log.Println("Error retrieving vote counts from Redis:", err)
		return
	}
	log.Printf("votesMap: %+v\n", votesMap)

	for client := range clients {
		for option, countStr := range votesMap {
			count, _ := strconv.Atoi(countStr)
			voteCount := models.VoteCount{Option: option, Count: count}
			err := client.WriteJSON(voteCount)
			if err != nil {
				log.Println("Error writing JSON to client:", err)
				delete(clients, client)
				client.Close()
			}
		}
	}
}

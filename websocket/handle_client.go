package websocket

import (
	"log"

	"govoting/models"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)

func handleClient(conn *websocket.Conn) {
	defer func() {
		delete(clients, conn)
		conn.Close()
	}()

	for {
		var vote models.Vote
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

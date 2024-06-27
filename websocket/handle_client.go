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

	var voted = false

	for {
		var vote models.Vote
		err := conn.ReadJSON(&vote)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading JSON: %v", err)
			}
			break
		}

		if voted == false {
			// Update vote count in Redis
			err = incrementVoteCount(vote.Option, conn)
			if err != nil {
				log.Println("Error incrementing vote count in Redis:", err)
				break
			}
			voted = true
		} else {
			log.Println("User has already voted, ignoring subsequent vote")
		}

		// Broadcast updated vote count to all clients
		broadcastVoteCount()
	}
}

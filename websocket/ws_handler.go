package websocket

import (
	"log"
	"net/http"
)

func WsHandler(w http.ResponseWriter, r *http.Request) {
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

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

	clients[conn] = true

	go handleClient(conn)
}

package main

import (
	"log"
	"net/http"

	"govoting/websocket"
)

func main() {
	// Set up HTTP routes
	http.HandleFunc("/ws", websocket.WsHandler)
	log.Fatal(http.ListenAndServe(":8040", nil))
}

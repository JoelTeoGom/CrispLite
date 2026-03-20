package main

import (
	"fmt"
	"net/http"

	"crisplite/internal/adapters/inbound/ws"
	database "crisplite/internal/adapters/outbound/stdout"
)

func main() {
	database := database.NewPostgresAdapter()

	wsHandler := ws.NewHandler(database)
	http.HandleFunc("/ws/v1/chat", wsHandler.ChatHandler)

	fmt.Println("WebSocket server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

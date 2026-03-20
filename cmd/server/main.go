package main

import (
	"fmt"
	"net/http"
	"time"

	"crisplite/internal/adapters/inbound/ws"
	"crisplite/internal/adapters/outbound/stdout"
	"crisplite/internal/application"
	"crisplite/internal/domain"
)

func main() {
	processor := stdout.NewBatchProcessor()
	messages := make(chan domain.Message, 10)

	go application.Batcher(messages, 10, 500*time.Millisecond, processor)

	wsHandler := ws.NewHandler(messages)
	http.HandleFunc("/ws/v1/chat", wsHandler.ChatHandler)

	fmt.Println("WebSocket server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

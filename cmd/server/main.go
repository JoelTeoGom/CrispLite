package main

import (
	"context"
	"crisplite/internal/adapter/outbound/postgres"
	"crisplite/internal/app"
	"crisplite/internal/config"
	"crisplite/internal/domain"
	"fmt"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	pool, err := postgres.InitPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer pool.Close()

	msgChannel := make(chan domain.Message, cfg.BatchSize)
	defer close(msgChannel)

	pg := postgres.NewMessageRepo(pool)

	batcher := app.NewBatcher(msgChannel, cfg.BatchSize, cfg.BatchInterval, pg)
	batcher.Start()
	//http.HandleFunc("/ws/v1/chat", wsHandler.ChatHandler)

	fmt.Printf("WebSocket server started on :%s\n", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, nil); err != nil {
		log.Fatalf("server: %v", err)
	}
}

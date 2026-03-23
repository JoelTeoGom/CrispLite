package main

import (
	"context"
	"crisplite/internal/adapter/outbound/env"
	"crisplite/internal/adapter/outbound/postgres"
	"crisplite/internal/app"
	"crisplite/internal/domain"
	"fmt"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	loader, err := env.NewConfigLoader(".env")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	cfg, err := loader.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	pool, err := postgres.NewPool(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer pool.Close()

	msgChannel := make(chan domain.Message, cfg.Batcher.Size)
	defer close(msgChannel)

	pg := postgres.NewMessageRepo(pool)

	batcher := app.NewBatcher(msgChannel, cfg.Batcher.Size, cfg.Batcher.Interval, pg)
	batcher.Start()
	//http.HandleFunc("/ws/v1/chat", wsHandler.ChatHandler)

	fmt.Printf("WebSocket server started on :%s\n", cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, nil); err != nil {
		log.Fatalf("server: %v", err)
	}
}

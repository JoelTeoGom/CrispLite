package main

import (
	"context"
	"crisplite/internal/adapter/inbound/rest"
	"crisplite/internal/adapter/outbound/config"
	locallogger "crisplite/internal/adapter/outbound/local_logger"
	"crisplite/internal/adapter/outbound/postgres"
	"crisplite/internal/app"
	"crisplite/internal/domain"
	"crisplite/internal/port/outbound"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	//CONFIGs
	loader, err := config.NewConfigLoader(".env")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	cfg, err := loader.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	log.Printf("Starting server in %s environment", cfg.Env)

	//LOGGING
	var loggerAdapter outbound.Logger
	if cfg.Env == domain.EnvLocal {
		loggerAdapter = locallogger.NewLocalLogger()
	} else {
		loggerAdapter = locallogger.NewLocalLogger() //TODO: replace with real logger
	}

	//DB
	pool, err := postgres.NewPool(ctx, cfg.Database, loggerAdapter)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer pool.Close()

	//REPOS
	messageRepo := postgres.NewMessageRepo(pool)
	userRepo := postgres.NewUserRepo(pool)

	//APP
	msgChannel := make(chan domain.Message, cfg.Batcher.Size)
	defer close(msgChannel)
	batcher := app.NewBatcher(msgChannel, cfg.Batcher.Size, cfg.Batcher.Interval, messageRepo)
	batcher.Start()
	defer batcher.Stop()

	router := http.NewServeMux()

	//SERVICES
	userService := app.NewUserService(userRepo)
	chatService := app.NewChatService(messageRepo, *batcher)

	//HANDLERS
	userHandler := rest.NewUserHandler(userService)
	chatHandler := rest.NewChatHandler(chatService)

	rest.RegisterRoutes(router, userHandler, chatHandler)

	if err := http.ListenAndServe(":"+cfg.Server.Port, router); err != nil {
		log.Fatalf("server: %v", err)
	}
}

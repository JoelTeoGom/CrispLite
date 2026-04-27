package main

import (
	"context"
	_ "crisplite/docs"
	"crisplite/internal/adapter/inbound/rest"
	"crisplite/internal/adapter/inbound/ws"
	"crisplite/internal/adapter/outbound/auth"
	"crisplite/internal/adapter/outbound/config"
	locallogger "crisplite/internal/adapter/outbound/local_logger"
	"crisplite/internal/adapter/outbound/postgres"
	"crisplite/internal/app"
	"crisplite/internal/domain"
	"crisplite/internal/port/outbound"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

const (
	_shutdownPeriod      = 15 * time.Second
	_shutdownHardPeriod  = 3 * time.Second
	_readinessDrainDelay = 5 * time.Second
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

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

	//REDIS

	//REPOS
	messageRepo := postgres.NewMessageRepo(pool, loggerAdapter)
	userRepo := postgres.NewUserRepo(pool, loggerAdapter)
	authRepo := postgres.NewAuthRepo(pool, loggerAdapter)

	//APP
	msgChannel := make(chan domain.Message, cfg.Batcher.Size)
	defer close(msgChannel)
	batcher := app.NewBatcher(msgChannel, cfg.Batcher.Size, cfg.Batcher.Interval, messageRepo, loggerAdapter)
	batcher.Start(ctx)
	defer batcher.Stop()

	//AUTH
	tokenService := auth.NewJWTService(cfg.Server.JWTSecret)
	router := http.NewServeMux()

	//SERVICES
	userService := app.NewUserService(userRepo, authRepo, tokenService, loggerAdapter)
	chatService := app.NewChatService(messageRepo, *batcher, loggerAdapter)
	hub := ws.NewHub(chatService)
	chatService.Hub = hub

	//HANDLERS
	authHandler := rest.NewAuthHandler(userService, loggerAdapter, cfg.Env)
	userHandler := rest.NewUserHandler(userService, tokenService, loggerAdapter)
	chatHandler := rest.NewChatHandler(chatService, loggerAdapter)

	handler := rest.RegisterRoutes(router, authHandler, userHandler, chatHandler, loggerAdapter, tokenService, cfg.Server.AllowedOrigin)

	if err := http.ListenAndServe(":"+cfg.Server.Port, handler); err != nil {
		log.Fatalf("server: %v", err)
	}
}

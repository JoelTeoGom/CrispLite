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
	"crisplite/internal/adapter/outbound/redis"
	"crisplite/internal/app"
	"crisplite/internal/domain"
	"crisplite/internal/port/outbound"
	"fmt"
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
	redisClient, err := redis.NewClient(ctx, loggerAdapter, cfg.Redis)
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	defer redisClient.Close()
	pubsub, err := redis.NewPubSub(ctx, loggerAdapter, cfg.Redis)
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	defer pubsub.Close()

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
	hub := ws.NewHub(chatService, pubsub)
	chatService.Hub = hub

	//HANDLERS
	authHandler := rest.NewAuthHandler(userService, loggerAdapter, cfg.Env)
	userHandler := rest.NewUserHandler(userService, tokenService, loggerAdapter)
	chatHandler := rest.NewChatHandler(chatService, loggerAdapter)

	handler := rest.RegisterRoutes(router, authHandler, userHandler, chatHandler, loggerAdapter, tokenService, cfg.Server.AllowedOrigin)

	serverErr := make(chan error, 1)
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: handler,
	}

	go func() {
		log.Printf("Server is running on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("Shutdown signal received, shutting down gracefully...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), _shutdownPeriod)
		log.Println()
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
		}
	case err := <-serverErr:
		log.Printf("Server error: %v", err)
	}

	fmt.Println("finally shudown process ended CLOSED FOREVER")
}

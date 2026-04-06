package main

import (
	"context"
	"crisplite/internal/adapter/inbound/rest"
	_ "crisplite/docs"
	"crisplite/internal/adapter/outbound/auth"
	"crisplite/internal/adapter/outbound/config"
	locallogger "crisplite/internal/adapter/outbound/local_logger"
	"crisplite/internal/adapter/outbound/postgres"
	"crisplite/internal/app"
	"crisplite/internal/domain"
	"crisplite/internal/port/outbound"
	"log"
	"net/http"
)

// @title           CrispLite API
// @version         1.0
// @description     CrispLite chat application API
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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
	messageRepo := postgres.NewMessageRepo(pool, loggerAdapter)
	userRepo := postgres.NewUserRepo(pool, loggerAdapter)
	authRepo := postgres.NewAuthRepo(pool, loggerAdapter)

	//APP
	msgChannel := make(chan domain.Message, cfg.Batcher.Size)
	defer close(msgChannel)
	batcher := app.NewBatcher(msgChannel, cfg.Batcher.Size, cfg.Batcher.Interval, messageRepo, loggerAdapter)
	batcher.Start()
	defer batcher.Stop()

	//AUTH
	tokenService := auth.NewJWTService(cfg.Server.JWTSecret)

	router := http.NewServeMux()

	//SERVICES
	userService := app.NewUserService(userRepo, authRepo, tokenService, loggerAdapter)
	chatService := app.NewChatService(messageRepo, *batcher, loggerAdapter)

	//HANDLERS
	authHandler := rest.NewAuthHandler(userService, loggerAdapter)
	userHandler := rest.NewUserHandler(userService, loggerAdapter)
	chatHandler := rest.NewChatHandler(chatService, loggerAdapter)

	handler := rest.RegisterRoutes(router, authHandler, userHandler, chatHandler, loggerAdapter, tokenService)

	if err := http.ListenAndServe(":"+cfg.Server.Port, handler); err != nil {
		log.Fatalf("server: %v", err)
	}
}

.PHONY: run build test migrate-up migrate-down migrate-create docker-up docker-down

# Config (loads DB_URL from .env)
include .env
export
MIGRATE_CMD = migrate -path migrations -database "$(DATABASE_URL)?sslmode=disable"

# App
run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

test:
	go test ./...

# Docker
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate-up:
	$(MIGRATE_CMD) up

migrate-down:
	$(MIGRATE_CMD) down 1

migrate-force:
	$(MIGRATE_CMD) force $(VERSION)

migrate-create:
	migrate create -ext sql -dir migrations -seq $(NAME)

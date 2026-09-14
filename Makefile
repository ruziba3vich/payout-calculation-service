# Makefile


ifneq (,$(wildcard .env))
    include .env
    export
endif

APP_NAME          ?= payout-calculation-service
DOCKER_COMPOSE    ?= docker compose
MIGRATE           ?= migrate
MIGRATIONS_PATH   ?= ./migrations
SQLC              ?= sqlc
SWAG              ?= $(shell go env GOPATH)/bin/swag
SQLC_CONFIG       ?= ./internal/infrastructure/persistence/postgres/sqlc.yaml

GREEN  := \033[0;32m
YELLOW := \033[0;33m
NC     := \033[0m


help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

run:
	@go run ./cmd/api

build:
	@go build -o bin/$(APP_NAME) ./cmd/api

test:
	@go test ./... -count=1 -race

tidy:
	@go mod tidy

sqlc:
	@echo "$(YELLOW)→ generating $(SQLC_CONFIG)$(NC)"
	@$(SQLC) generate -f $(SQLC_CONFIG)

swag:
	@echo "$(YELLOW)→ generating swagger$(NC)"
	@$(SWAG) init -g cmd/api/main.go -o docs --parseDependency --parseInternal -q

sqlc-vet:
	@echo "$(YELLOW)→ compiling $(SQLC_CONFIG)$(NC)"
	@$(SQLC) compile -f $(SQLC_CONFIG)

migrate-up:
	@if [ -z "$(DB_URL)" ]; then \
		echo "DB_URL is not set. Please define it in .env or export it"; \
		exit 1; \
	fi
	@if [ -z "$(VERSION)" ]; then \
		echo "$(YELLOW)→ migrating up to latest...$(NC)"; \
		$(MIGRATE) -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up; \
	else \
		echo "$(YELLOW)→ migrating up to version $(VERSION)...$(NC)"; \
		$(MIGRATE) -path $(MIGRATIONS_PATH) -database "$(DB_URL)" goto $(VERSION); \
	fi

migrate-down:
	@if [ -z "$(DB_URL)" ]; then \
		echo "DB_URL is not set. Please define it in .env or export it"; \
		exit 1; \
	fi
	@if [ -z "$(VERSION)" ]; then \
		echo "$(YELLOW)→ rolling back 1 migration...$(NC)"; \
		$(MIGRATE) -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1; \
	else \
		echo "$(YELLOW)→ rolling back to version $(VERSION)...$(NC)"; \
		$(MIGRATE) -path $(MIGRATIONS_PATH) -database "$(DB_URL)" goto $(VERSION); \
	fi

migrate-force:
	@if [ -z "$(DB_URL)" ]; then \
		echo "DB_URL is not set. Please define it in .env or export it"; \
		exit 1; \
	fi
	@if [ -z "$(VERSION)" ]; then \
		echo "VERSION is required. Example: make migrate-force VERSION=1"; \
		exit 1; \
	fi
	@$(MIGRATE) -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force $(VERSION)

migrate-create:
	@if [ -z "$(filter-out $@,$(MAKECMDGOALS))" ]; then \
		echo "Name is required. Example: make migrate-create add_payouts"; \
		exit 1; \
	fi
	@$(MIGRATE) create -ext sql -dir $(MIGRATIONS_PATH) -seq $(filter-out $@,$(MAKECMDGOALS))

%:
	@:

docker-up:
	@$(DOCKER_COMPOSE) up -d --build

docker-down:
	@$(DOCKER_COMPOSE) down -v

docker-logs:
	@$(DOCKER_COMPOSE) logs -f

clean:
	@rm -rf bin/
	@go clean

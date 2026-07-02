##change lint version because we need to use 1.24.0 ver of go for proto etc.
GOLANGCI_LINT_VERSION := v1.64.8

# Docker compose settings
DOCKER_COMPOSE_FILE := docker/docker-compose.yml
ENV_FILE := .env

.PHONY: lint
lint: ## run go linter
	docker run --rm -v `pwd`:/app -w /app golangci/golangci-lint:$(GOLANGCI_LINT_VERSION) golangci-lint run --timeout 2m

.PHONY: gen
gen: ## generate go code from proto
	mkdir -p pkg/gen
	protoc --proto_path=pkg/pb \
	--go_out=pkg/gen --go_opt=paths=source_relative \
	--go-grpc_out=pkg/gen --go-grpc_opt=paths=source_relative \
	pkg/pb/*.proto

# Docker commands
.PHONY: docker-build
docker-build: ## build docker images
	docker compose --env-file $(ENV_FILE) -f $(DOCKER_COMPOSE_FILE) build

.PHONY: docker-up
docker-up: ## start all services
	docker compose --env-file $(ENV_FILE) -f $(DOCKER_COMPOSE_FILE) up -d

.PHONY: docker-up-build
docker-up-build: ## build and start all services
	docker compose --env-file $(ENV_FILE) -f $(DOCKER_COMPOSE_FILE) up --build -d

.PHONY: docker-down
docker-down: ## stop and remove containers
	docker compose --env-file $(ENV_FILE) -f $(DOCKER_COMPOSE_FILE) down

.PHONY: docker-down-volumes
docker-down-volumes: ## stop and remove containers with volumes
	docker compose --env-file $(ENV_FILE) -f $(DOCKER_COMPOSE_FILE) down -v

.PHONY: docker-logs
docker-logs: ## show logs from all services
	docker compose --env-file $(ENV_FILE) -f $(DOCKER_COMPOSE_FILE) logs -f

.PHONY: docker-logs-app
docker-logs-app: ## show logs from application services only
	docker compose --env-file $(ENV_FILE) -f $(DOCKER_COMPOSE_FILE) logs -f issue-tracker issue-tracker-gateway

.PHONY: docker-restart
docker-restart: ## restart all services
	docker compose --env-file $(ENV_FILE) -f $(DOCKER_COMPOSE_FILE) restart

.PHONY: docker-ps
docker-ps: ## show running containers
	docker compose --env-file $(ENV_FILE) -f $(DOCKER_COMPOSE_FILE) ps

.PHONY: docker-clean
docker-clean: docker-down-volumes ## clean up everything (containers, volumes, images)
	docker compose --env-file $(ENV_FILE) -f $(DOCKER_COMPOSE_FILE) down --rmi all -v
	docker system prune -f

# Development workflow commands
.PHONY: dev-up
dev-up: docker-up-build ## full development setup (build and start)

.PHONY: dev-down
dev-down: docker-down-volumes ## clean development teardown

.PHONY: dev-restart
dev-restart: docker-down docker-up-build ## restart development environment

# Database specific commands
.PHONY: db-up
db-up: ## start only database
	docker compose --env-file $(ENV_FILE) -f $(DOCKER_COMPOSE_FILE) up postgres -d

.PHONY: db-logs
db-logs: ## show database logs
	docker compose --env-file $(ENV_FILE) -f $(DOCKER_COMPOSE_FILE) logs -f postgres

# Help command
.PHONY: help
help: ## show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
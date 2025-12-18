# OrgNote Backend - Deployment Commands

# Default recipe
default:
    @just --list

# === LOCAL DEVELOPMENT ===

# Run tests
test:
    go test ./...

# Run tests with verbose output
test-v:
    go test -v ./...

# Run tests with coverage
test-coverage:
    go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html

local:
    docker compose -f docker-compose.db.local.yaml -f docker-compose.local.yaml --env-file .env up -d orgnote-mongo-local orgnote-minio-local orgnote-minio-init-local
    docker compose -f docker-compose.db.local.yaml -f docker-compose.local.yaml --env-file .env up --build orgnote-backend-local

local-down:
    docker compose -f docker-compose.db.local.yaml -f docker-compose.local.yaml --env-file .env down

# === TRAEFIK (shared) ===

traefik:
    docker compose -f docker-compose.traefik.yaml --env-file .env.traefik up -d

traefik-logs:
    docker logs orgnote-traefik -f

# === PRODUCTION ===

deploy-prod:
    docker compose -f docker-compose.traefik.yaml --env-file .env.traefik up -d
    docker compose -f docker-compose.db.prod.yaml -f docker-compose.prod.yaml --env-file .env.prod up -d --build

prod-down:
    docker compose -f docker-compose.db.prod.yaml -f docker-compose.prod.yaml --env-file .env.prod down

prod-logs:
    docker compose -f docker-compose.db.prod.yaml -f docker-compose.prod.yaml --env-file .env.prod logs -f

prod-logs-backend:
    docker logs orgnote-backend-prod -f

# === DEVELOPMENT (server) ===

deploy-dev:
    docker compose -f docker-compose.traefik.yaml --env-file .env.traefik up -d
    docker compose -f docker-compose.db.dev.yaml -f docker-compose.dev.yaml --env-file .env.dev up -d --build

dev-down:
    docker compose -f docker-compose.db.dev.yaml -f docker-compose.dev.yaml --env-file .env.dev down

dev-logs:
    docker compose -f docker-compose.db.dev.yaml -f docker-compose.dev.yaml --env-file .env.dev logs -f

dev-logs-backend:
    docker logs orgnote-backend-dev -f

# === UTILITIES ===

prune:
    docker system prune -f

ps:
    docker ps

# Check all services status
status:
    @echo "=== All OrgNote containers ===" && docker ps --filter name=orgnote

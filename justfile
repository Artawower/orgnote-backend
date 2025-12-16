# OrgNote Backend - Deployment Commands

# Default recipe
default:
    @just --list

# === LOCAL DEVELOPMENT ===

local:
    docker compose -f docker-compose.db.local.yaml up -d
    docker compose -f docker-compose.local.yaml --env-file .env.local up --build

local-down:
    docker compose -f docker-compose.local.yaml down
    docker compose -f docker-compose.db.local.yaml down

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
    docker compose -f docker-compose.prod.yaml down
    docker compose -f docker-compose.db.prod.yaml down

prod-logs:
    docker compose -f docker-compose.prod.yaml logs -f

prod-logs-backend:
    docker logs orgnote-backend-prod -f

# === DEVELOPMENT (server) ===

deploy-dev:
    docker compose -f docker-compose.traefik.yaml --env-file .env.traefik up -d
    docker compose -f docker-compose.db.dev.yaml -f docker-compose.dev.yaml --env-file .env.dev up -d --build

dev-down:
    docker compose -f docker-compose.dev.yaml down
    docker compose -f docker-compose.db.dev.yaml down

dev-logs:
    docker compose -f docker-compose.dev.yaml logs -f

dev-logs-backend:
    docker logs orgnote-backend-dev -f

# === UTILITIES ===

prune:
    docker system prune -f

ps:
    docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

# Check all services status
status:
    @echo "=== Traefik ===" && docker ps --filter name=orgnote-traefik --format "{{.Names}}: {{.Status}}" || true
    @echo "=== Prod ===" && docker ps --filter name=-prod --format "{{.Names}}: {{.Status}}" || true
    @echo "=== Dev ===" && docker ps --filter name=-dev --format "{{.Names}}: {{.Status}}" || true

set dotenv-load

# Local development
local:
    docker compose -f docker-compose.db.local.yaml -f docker-compose.local.yaml up

local-build:
    docker compose -f docker-compose.db.local.yaml -f docker-compose.local.yaml up --build --force-recreate

local-down:
    docker compose -f docker-compose.db.local.yaml -f docker-compose.local.yaml down

# Dev deployment (dev.org-note.com)
dev:
    docker compose -f docker-compose.traefik.yaml --env-file .env.traefik up -d
    docker compose -f docker-compose.db.dev.yaml --env-file .env.dev up -d
    docker compose -f docker-compose.dev.yaml --env-file .env.dev up -d

dev-build:
    docker compose -f docker-compose.traefik.yaml --env-file .env.traefik up -d
    docker compose -f docker-compose.db.dev.yaml --env-file .env.dev up -d
    docker compose -f docker-compose.dev.yaml --env-file .env.dev up -d --build --force-recreate

dev-down:
    docker compose -f docker-compose.dev.yaml down
    docker compose -f docker-compose.db.dev.yaml down

# Production deployment (org-note.com)
prod:
    docker compose -f docker-compose.traefik.yaml --env-file .env.traefik up -d
    docker compose -f docker-compose.db.prod.yaml --env-file .env.prod up -d
    docker compose -f docker-compose.prod.yaml --env-file .env.prod up -d

prod-build:
    docker compose -f docker-compose.traefik.yaml --env-file .env.traefik up -d
    docker compose -f docker-compose.db.prod.yaml --env-file .env.prod up -d
    docker compose -f docker-compose.prod.yaml --env-file .env.prod up -d --build --force-recreate

prod-down:
    docker compose -f docker-compose.prod.yaml down
    docker compose -f docker-compose.db.prod.yaml down

# Traefik management
traefik:
    docker compose -f docker-compose.traefik.yaml --env-file .env.traefik up -d

traefik-down:
    docker compose -f docker-compose.traefik.yaml down

traefik-logs:
    docker logs -f orgnote-traefik

# Logs
logs-local:
    docker logs -f orgnote-backend-local

logs-dev:
    docker logs -f orgnote-backend-dev

logs-prod:
    docker logs -f orgnote-backend-prod

# Testing
test:
    go test -v ./...

# Cleanup
prune:
    docker system prune -f

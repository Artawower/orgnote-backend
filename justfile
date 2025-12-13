run-dev:
  docker compose -f docker-compose.db.yaml -f docker-compose.s3.yaml -f docker-compose.dev.yaml up

run-dev-force:
  docker compose -f docker-compose.db.yaml -f docker-compose.s3.yaml -f docker-compose.dev.yaml up --build --force-recreate


test:
    go test -v ./...

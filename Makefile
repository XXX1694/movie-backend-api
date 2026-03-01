.PHONY: up down build restart logs ps clean test

# Launch everything
up:
	docker compose up --build -d

# Stop everything
down:
	docker compose down

# Build without starting
build:
	docker compose build

# Restart all services
restart:
	docker compose down
	docker compose up --build -d

# Show logs
logs:
	docker compose logs -f

# Show running containers
ps:
	docker compose ps

# Stop and remove volumes (full clean)
clean:
	docker compose down -v

# Run Go tests
test:
	go test ./...

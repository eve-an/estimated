COMPOSE = docker compose -f docker/docker-compose.yml

.PHONY: build up down logs restart

build:
	$(COMPOSE) build

up:
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f

restart: down build up

restart-app:
	$(COMPOSE) down estimated && $(COMPOSE) build estimated && $(COMPOSE) up -d estimated

run:
	go run ./cmd/myapp


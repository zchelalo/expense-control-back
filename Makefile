ifneq ("$(wildcard .env)", "")
	include .env
	export $(shell sed 's/=.*//' .env)
endif

DOCKER_COMPOSE_FILE = ./.docker/compose.yml
DOCKER_NETWORK_PREFIX = docker_
DOCKER_NETWORK_NAME = expense-control-back-network
DOCKER_NETWORK = $(DOCKER_NETWORK_PREFIX)$(DOCKER_NETWORK_NAME)

URI_DB = postgresql://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATE = docker run -it -v $(shell pwd)/internal/db/migrations:/migrations --network $(DOCKER_NETWORK) migrate/migrate -path /migrations -database "$(URI_DB)" -verbose

setup:
	$(MAKE) create-envs
	$(MAKE) create-keys
	$(MAKE) compose-build-detached
	docker run --rm --network=$(DOCKER_NETWORK) \
		-v $(shell pwd)/scripts:/scripts alpine sh /scripts/wait_for_db.sh $(DB_HOST) $(DB_PORT)
	$(MAKE) migrate-up

migrate-up:
	$(MIGRATE) up

migrate-up-1:
	$(MIGRATE) up 1

migrate-down:
	$(MIGRATE) down

migrate-down-1:
	$(MIGRATE) down 1

compose:
	docker compose -f $(DOCKER_COMPOSE_FILE) up

compose-detached:
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d

compose-build:
	docker compose -f $(DOCKER_COMPOSE_FILE) up --build

compose-build-detached:
	docker compose -f $(DOCKER_COMPOSE_FILE) up --build -d

compose-down:
	docker compose -f $(DOCKER_COMPOSE_FILE) down

create-envs:
	cp .env.example .env

create-keys:
	./scripts/create_keys.sh

sqlc:
	sqlc generate

.PHONY: migrate-up migrate-up-1 migrate-down migrate-down-1 setup compose compose-detached compose-build compose-build-detached compose-down create-envs create-keys sqlc
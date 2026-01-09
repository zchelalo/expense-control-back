ifneq ("$(wildcard .env)", "")
	include .env
	export $(shell sed 's/=.*//' .env)
endif

DOCKER_COMPOSE_FILE = ./.docker/compose.yml

setup:
	$(MAKE) create-envs
	$(MAKE) compose-build-detached

compose:
	docker compose -f $(DOCKER_COMPOSE_FILE) up

compose-build:
	docker compose -f $(DOCKER_COMPOSE_FILE) up --build

compose-build-detached:
	docker compose -f $(DOCKER_COMPOSE_FILE) up --build -d

compose-down:
	docker compose -f $(DOCKER_COMPOSE_FILE) down

create-envs:
	cp .env.example .env

.PHONY: setup compose compose-build compose-build-detached compose-down create-envs
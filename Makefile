# sketch-api-go/

ifneq (,$(wildcard .env))
	include .env
	export
endif

# CONSTANS
IMAGE := my-api
CONTAINER := my-api-container
ENV_FILE=.env
ENV_EXAMPLE=.env.example
COMPOSE=docker compose
BASE=-f compose.yml

# make запускается на хостовой машине, поэтому localhost.
MIGRATION_DB_HOST ?= localhost
MIGRATIONS_DIR := db/migrations
MIGRATION_DB_PORT ?= 5433
DATABASE_URL = postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(MIGRATION_DB_HOST):$(MIGRATION_DB_PORT)/$(POSTGRES_DB)?sslmode=disable

# JWT
PRIVATE_KEY=./secrets/private.pem
PUBLIC_KEY=./secrets/public.pem

TEST_PRIVATE_KEY=./secrets/test-private.pem
TEST_PUBLIC_KEY=./secrets/test-public.pem

.DEFAULT_GOAL := help

# App
.PHONY: help run

# Generation
.PHONY: env-gen api-v1gen sql-gen jwt-gen

# Docker
.PHONY: docker-build docker-run docker-stop docker-down docker-clean

# Compose
.PHONY: compose-dev compose-clean

# Migrates
.PHONY: migrate-create migrate-up migrate-down migrate-version

help:
	@echo "sketch-api-go project"
	@echo "Usage:"
	@echo "  make help         - show this help"
	@echo "  make gen          - generate Go code from the OpenAPI specification"
	@echo "  make run          - generate code and run the API server locally"
	@echo "  make docker-build - generate code and build the Docker image"
	@echo "  make docker-run   - run the Docker container on port 8080"
	@echo "  make docker-stop  - stop the running Docker container"
	@echo "  make docker-down  - remove the Docker container"
	@echo "  make docker-clean - remove Docker resources and build cache"

env-gen:
	@if [ ! -f $(ENV_FILE) ]; then \
		echo "Creating .env from .env.example"; \
		cp $(ENV_EXAMPLE) $(ENV_FILE); \
	fi

# Генерирует JWT RSA ключи, если они отсутствуют
jwt-gen:
	@if [ -f $(PRIVATE_KEY) ] && [ -f $(PUBLIC_KEY) ]; then \
		echo "JWT keys already exist. Skipping generation."; \
	else \
		echo "Generating JWT RS256 keys..."; \
		mkdir -p secrets; \
		openssl genrsa -out $(PRIVATE_KEY) 2048; \
		openssl rsa -in $(PRIVATE_KEY) -pubout -out $(PUBLIC_KEY); \
		echo "JWT keys generated in ./secrets"; \
	fi

jwt-test-gen:
	@if [ -f $(TEST_PRIVATE_KEY) ] && [ -f $(TEST_PUBLIC_KEY) ]; then \
		echo "JWT keys already exist. Skipping generation."; \
	else \
		echo "Generating JWT RS256 keys..."; \
		mkdir -p secrets; \
		openssl genrsa -out $(TEST_PRIVATE_KEY) 2048; \
		openssl rsa -in $(TEST_PRIVATE_KEY) -pubout -out $(TEST_PUBLIC_KEY); \
		echo "JWT keys generated in ./secrets"; \
	fi


api-v1gen:
	@echo "===Generating Go code from OpenAPI==="
	go tool oapi-codegen -config ./api/v1/configs/server.yml ./api/v1/openapi.yml
	go tool oapi-codegen -config ./api/v1/configs/models.yml ./api/v1/openapi.yml
	go tool oapi-codegen -config ./api/v1/configs/spec.yml ./api/v1/openapi.yml
	@echo "===Generation complete!==="

run:
	@echo "===Run App==="
	go run ./cmd/server/

docker-build: gen
	docker build -t $(IMAGE) .

docker-run:
	docker run \
		--rm \
		--name $(CONTAINER) \
		-p 8080:8080 \
		$(IMAGE)

docker-stop:
	docker stop $(CONTAINER)

docker-down:
	-docker rm -f $(CONTAINER)

docker-clean: docker-down
	-docker image rm $(IMAGE)
	docker builder prune -af

compose-dev:
	$(COMPOSE) --env-file ./$(ENV_FILE) \
		$(BASE) \
		-f infra/compose/dev.yml \
		up --build

compose-clean:
	$(COMPOSE) -f compose.yml -f infra/compose/dev.yml down -v --remove-orphans
	docker builder prune -af


migrate-create:
	@test -n "$(name)" || \
		(echo "Usage: make migrate-create name=add_user_status" && exit 1)
	migrate create \
		-ext sql \
		-dir $(MIGRATIONS_DIR) \
		-seq \
		$(name)

migrate-up:
	migrate \
		-path $(MIGRATIONS_DIR) \
		-database "$(DATABASE_URL)" \
		up

migrate-down:
	migrate \
		-path $(MIGRATIONS_DIR) \
		-database "$(DATABASE_URL)" \
		down 1

migrate-version:
	migrate \
		-path $(MIGRATIONS_DIR) \
		-database "$(DATABASE_URL)" \
		version

sql-gen:
	sqlc generate

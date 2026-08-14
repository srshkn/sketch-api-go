# sketch-api-go/

ifneq (,$(wildcard .env))
	include .env
	export
endif

.DEFAULT_GOAL := help


# -------------------------------------------------------------------------
# CONSTANTS

# Environment
ENV_FILE=.env
ENV_EXAMPLE=.env.example

# Compose
COMPOSE=docker compose
BASE=-f compose.yml

# Migrations
MIGRATION_DB_HOST ?= localhost
MIGRATIONS_DIR := db/migrations
MIGRATION_DB_PORT ?= 5433
DATABASE_URL = postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(MIGRATION_DB_HOST):$(MIGRATION_DB_PORT)/$(POSTGRES_DB)?sslmode=disable

# JWT
PRIVATE_KEY=./secrets/private.pem
PUBLIC_KEY=./secrets/public.pem
TEST_PRIVATE_KEY=./secrets/test-private.pem
TEST_PUBLIC_KEY=./secrets/test-public.pem

# -------------------------------------------------------------------------
# PHONY

# Base
.PHONY: help

# Generation
.PHONY: env-gen sql-gen jwt-gen jwt-test-gen api-v1gen

# Compose
.PHONY: compose-dev compose-clean

# Migrations
.PHONY: migrate-create migrate-up migrate-down migrate-version

# -------------------------------------------------------------------------
# BASE

help:
	@echo "sketch-api-go"
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Generation:"
	@echo "  env-gen              Create .env from .env.example"
	@echo "  jwt-gen              Generate JWT RSA keys"
	@echo "  jwt-test-gen         Generate JWT RSA keys for tests"
	@echo "  api-v1gen            Generate Go code from OpenAPI specification"
	@echo "  sql-gen              Generate Go code from SQL queries"
	@echo ""
	@echo "Compose:"
	@echo "  compose-dev          Start development environment with Docker Compose"
	@echo "  compose-clean        Stop and remove development environment and volumes"
	@echo ""
	@echo "Migrations:"
	@echo "  migrate-create       Create a new migration (name=<migration_name>)"
	@echo "  migrate-up           Apply all pending migrations"
	@echo "  migrate-down         Roll back the last migration"
	@echo "  migrate-version      Show current migration version"

# -------------------------------------------------------------------------
# GENERATION

env-gen:
	@if [ ! -f $(ENV_FILE) ]; then \
		echo "Creating .env from .env.example"; \
		cp $(ENV_EXAMPLE) $(ENV_FILE); \
	fi

sql-gen:
	sqlc generate

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
	go tool oapi-codegen -config ./api/v1/configs/server.yml ./api/v1/openapi.yml
	go tool oapi-codegen -config ./api/v1/configs/models.yml ./api/v1/openapi.yml
	go tool oapi-codegen -config ./api/v1/configs/spec.yml ./api/v1/openapi.yml

# -------------------------------------------------------------------------
# COMPOSE

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

# -------------------------------------------------------------------------
# MIGRATIONS

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

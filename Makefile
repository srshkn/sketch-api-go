# sketch-api-go

.DEFAULT_GOAL := help

.PHONY: help run

help:
	@echo "sketch-api-go project"
	@echo "Usage:"
	@echo "  make run - ..."

run:
	go run ./cmd/server/main.go

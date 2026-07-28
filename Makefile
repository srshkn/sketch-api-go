# sketch-api-go

IMAGE := my-api
CONTAINER := my-api-container

.DEFAULT_GOAL := help

.PHONY: help run docker-build docker-run docker-stop

help:
	@echo "sketch-api-go project"
	@echo "Usage:"
	@echo "  make run - ..."

run:
	go run ./cmd/server/main.go

docker-build:
	docker build -t $(IMAGE) .

docker-run:
	docker run \
		--rm \
		--name $(CONTAINER) \
		-p 8080:8080 \
		$(IMAGE)

docker-stop:
	docker stop $(CONTAINER)

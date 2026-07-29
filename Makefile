# sketch-api-go/

IMAGE := my-api
CONTAINER := my-api-container

.DEFAULT_GOAL := help

.PHONY: help gen run docker-build docker-run docker-stop docker-down docker-clean

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

gen:
	@echo "===Generating Go code from OpenAPI==="
	go tool oapi-codegen -config ./api/configs/server.yml ./api/openapi.yml
	go tool oapi-codegen -config ./api/configs/models.yml ./api/openapi.yml
	go tool oapi-codegen -config ./api/configs/spec.yml ./api/openapi.yml
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

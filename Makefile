# sketch-api-go/

IMAGE := my-api
CONTAINER := my-api-container

.DEFAULT_GOAL := help

.PHONY: help gen run docker-build docker-run docker-stop

help:
	@echo "sketch-api-go project"
	@echo "Usage:"
	@echo "  make help - ..."
	@echo "  make gen - ..."
	@echo "  make run - ..."
	@echo "  make docker-build - ..."
	@echo "  make docker-run - ..."
	@echo "  make docker-stop - ..."

gen:
	@echo "===Generating Go code from OpenAPI==="
	go tool oapi-codegen -config ./api/configs/server.yml ./api/openapi.yml
	go tool oapi-codegen -config ./api/configs/models.yml ./api/openapi.yml
	go tool oapi-codegen -config ./api/configs/spec.yml ./api/openapi.yml
	@echo "===Generation complete!==="

run: gen
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

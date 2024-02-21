.DEFAULT_GOAL := build

fmt:
	@go fmt ./...
.PHONY:fmt

lint: fmt
	@golangci-lint run --issues-exit-code 0

vet: lint
	@go vet ./...
.PHONY:lint

test: vet
	@go test ./...
.PHONY:vet

build: generate-docs
	@go build -o video-hosting cmd/app/main.go
.PHONY:build

run: generate-docs
	@go run cmd/app/main.go
.PHONY:run

deploy:
	@docker compose -f deployments/docker-compose.yaml up -d
.PHONY:deploy

swagger: vet
	@swag init -g ./cmd/app/main.go -o docs
.PHONY:generate-docs

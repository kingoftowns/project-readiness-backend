# Makefile for GitLab Readiness API
# This file contains common commands for development and deployment

# Variables
BINARY_NAME=gitlab-readiness-api
MAIN_PATH=cmd/api/main.go
DOCKER_IMAGE=gitlab-readiness-api:latest

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Default target
.DEFAULT_GOAL := help

## help: Display this help message
.PHONY: help
help:
	@echo "GitLab Readiness API - Available Commands"
	@echo "========================================"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## build: Build the application binary
.PHONY: build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

## run: Run the application with default settings
.PHONY: run
run:
	DATABASE_TYPE=sqlite DATABASE_URL=gitlab_readiness.db $(GOCMD) run $(MAIN_PATH)

## run-dev: Run with hot reload (requires air)
.PHONY: run-dev
run-dev:
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	DATABASE_TYPE=sqlite DATABASE_URL=gitlab_readiness.db air

## test: Run all tests
.PHONY: test
test:
	DATABASE_TYPE=sqlite DATABASE_URL=:memory: $(GOTEST) -v ./...

## test-coverage: Run tests with coverage report
.PHONY: test-coverage
test-coverage:
	DATABASE_TYPE=sqlite DATABASE_URL=:memory: $(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

## clean: Clean build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

## deps: Download and tidy dependencies
.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) tidy

## lint: Run golangci-lint
.PHONY: lint
lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

## fmt: Format code
.PHONY: fmt
fmt:
	$(GOCMD) fmt ./...

## migrate-up: Run database migrations
.PHONY: migrate-up
migrate-up:
	@echo "Running migrations..."
	@$(GOCMD) run $(MAIN_PATH) migrate up

## docker-build: Build Docker image
.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE) .

## docker-run: Run application in Docker
.PHONY: docker-run
docker-run:
	docker run -p 8080:8080 -e DATABASE_TYPE=sqlite -e DATABASE_URL=/data/gitlab_readiness.db -v $(PWD)/data:/data $(DOCKER_IMAGE)

## postgres-start: Start PostgreSQL using Docker
.PHONY: postgres-start
postgres-start:
	docker run -d --name gitlab-postgres -p 5432:5432 \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DB=gitlab_readiness \
		postgres:16-alpine

## postgres-stop: Stop PostgreSQL container
.PHONY: postgres-stop
postgres-stop:
	docker stop gitlab-postgres
	docker rm gitlab-postgres

## api-docs: Generate API documentation
.PHONY: api-docs
api-docs:
	@echo "API documentation can be found in docs/API.md"

# Development shortcuts
.PHONY: dev
dev: run-dev

.PHONY: t
t: test

.PHONY: b
b: build
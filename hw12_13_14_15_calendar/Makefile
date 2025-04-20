BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"
CONFIG_PATH := "./configs/config.yaml"
MIGRATIONS_DIR := "./internal/storage/migrations"
DB_DSN := "postgres://user:password@localhost:5432/calendar?sslmode=disable"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config $(CONFIG_PATH)

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.63.4

lint: install-lint-deps
	golangci-lint run ./...

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_DSN) up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_DSN) down

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_DSN) status

migrate-create:
	goose -dir $(MIGRATIONS_DIR) create $(name) sql

.PHONY: build run build-img run-img version test lint migrate-up migrate-down migrate-status migrate-create

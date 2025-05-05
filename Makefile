# -----------------------------------------------------------------------------
# Algo‑Kite Makefile – build, test, lint, run in dev or Docker
# -----------------------------------------------------------------------------

APP_NAME := gokite
PKG := ./...

# Go settings
GO ?= go
GOFLAGS := -mod=vendor

# Docker settings
REGISTRY ?= local
IMAGE := $(REGISTRY)/$(APP_NAME):latest

.PHONY: help build test run clean docker lint fmt

help:
	@echo "Common targets:"
	@echo "  make build     – build the binary (./bin/$(APP_NAME))"
	@echo "  make test      – run unit tests with coverage"
	@echo "  make run       – start the app with live reload (air)"
	@echo "  make docker    – build Docker image $(IMAGE)"
	@echo "  make lint      – run golangci‑lint"
	@echo "  make fmt       – gofumpt all packages"

# -----------------------------------------------------------------------------
# Core commands
# -----------------------------------------------------------------------------

build:
	@mkdir -p bin
	$(GO) build $(GOFLAGS) -o bin/$(APP_NAME) cmd/server/*.go


test:
	$(GO) test $(GOFLAGS) -race -coverprofile=coverage.out $(PKG)

run: build
	@echo "\n[dev] starting $(APP_NAME)…"
	./bin/$(APP_NAME)

clean:
	rm -rf bin coverage.out

# -----------------------------------------------------------------------------
# Docker – expects Dockerfile in project root
# -----------------------------------------------------------------------------

docker:
	docker build -t $(IMAGE) .

# -----------------------------------------------------------------------------
# Dev quality tools (optional – add to go.mod/tools)
# -----------------------------------------------------------------------------

lint:
	@command -v golangci-lint >/dev/null 2>&1 || { echo >&2 "Install golangci-lint first"; exit 1; }
	golangci-lint run ./...

fmt:
	@command -v gofumpt >/dev/null 2>&1 || { echo >&2 "Install gofumpt first"; exit 1; }
	gofumpt -l -w $(shell $(GO) list -f '{{.Dir}}' $(PKG))

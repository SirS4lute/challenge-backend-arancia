SHELL := /bin/bash

APP_NAME := todo-api
BIN_DIR := $(CURDIR)/bin
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint

.PHONY: help run test lint fmt tidy tools clean

help:
	@echo "Targets:"
	@echo "  run   - run the API locally"
	@echo "  test  - run unit/integration tests"
	@echo "  lint  - run golangci-lint (installs locally if needed)"
	@echo "  fmt   - format code (gofmt)"
	@echo "  tidy  - go mod tidy"
	@echo "  clean - remove local build artifacts"

run:
	@go run ./cmd/api

test:
	@go test ./...

fmt:
	@gofmt -w .

tidy:
	@go mod tidy

tools: $(GOLANGCI_LINT)

$(GOLANGCI_LINT):
	@mkdir -p $(BIN_DIR)
	@GOBIN=$(BIN_DIR) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

lint: tools
	@$(GOLANGCI_LINT) run ./...

clean:
	@rm -rf $(BIN_DIR)


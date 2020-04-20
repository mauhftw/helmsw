## Set up script constants
BIN=helmsw
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get

mod: ## Run mod tidy
	@echo "Installing dependencies..."
	@$(GOCMD) mod tidy

build: ## Build a version
	@echo "Building $(BIN) application into dist folder..."
	@$(GOBUILD) -o dist/$(BIN)

clean: ## Remove temporary files
	@echo "Cleaning files..."
	@$(GOCLEAN) && rm dist/$(BIN)

help:
	@echo "Command list: "
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: mod build clean help

## Sets default command
.DEFAULT_GOAL := help

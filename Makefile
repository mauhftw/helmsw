# Set up script constants
BIN=helmsw
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get

# Install all the build and lint dependencies
.PHONY: setup
setup:
	@echo "Installing dependencies..."
	@$(GOGET) -u github.com/alecthomas/gometalinter
	gometalinter --install --update
	@$(MAKE) dep

# Run dep ensure and prune
.PHONY: dep
dep:
	@echo "Vendoring dependencies..."
	dep ensure

# Run goimports on all go files
.PHONY: fmt
fmt:
	@echo "Formating all .go files..."
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do goimports -w "$$file"; done

# Run all the linters
.PHONY: lint
lint:
	@echo "Running linter on all .go files..."
	gometalinter --vendor --disable-all \
		--enable=deadcode \
		--enable=ineffassign \
		--enable=gosimple \
		--enable=staticcheck \
		--enable=gofmt \
		--enable=goimports \
		--enable=misspell \
		--enable=errcheck \
		--enable=vet \
		--enable=vetshadow \
		--deadline=10m \
		./...

# Build a version
.PHONY: build
build:
	@echo "Building $(BIN) application into dist folder..."
	@$(GOBUILD) -o dist/$(BIN)

# Remove temporary files
.PHONY: clean
clean:
	@echo "Cleaning files..."
	@$(GOCLEAN) && rm dist/$(BIN)

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@echo "Command list: "
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# Sets default command
.DEFAULT_GOAL := build

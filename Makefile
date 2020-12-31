NAME := tmux-project
VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
COMMANDS := tmux-project-create tmux-project-list
LDFLAGS := -X 'main.version=$(VERSION)' \
           -X 'main.revision=$(REVISION)'
GOIMPORTS ?= goimports
GOCILINT ?= golangci-lint
GO ?= GO111MODULE=on go
.DEFAULT_GOAL := help

.PHONY: fmt
fmt: ## Formatting source codes.
	@$(GOIMPORTS) -w *.go config cmd


.PHONY: lint
lint: ## Run golint and go vet.
	@$(GOCILINT) run --no-config --disable-all --enable=goimports --enable=misspell *.go
	@$(GOCILINT) run --no-config --disable-all --enable=goimports --enable=misspell ./config/*.go
	@$(GOCILINT) run --no-config --disable-all --enable=goimports --enable=misspell ./cmd/*/*.go

.PHONY: test
test:  ## Run the tests.
	@$(GO) test ./...

.PHONY: build
build: main.go  ## Build a binary.
	$(foreach cmd,$(COMMANDS), $(GO) build -ldflags "$(LDFLAGS)" ./cmd/$(cmd);)

.PHONY: cross
cross: main.go  ## Build binaries for cross platform.
	mkdir -p pkg
	@# darwin
	@for arch in "amd64" "386"; do \
		GOOS=darwin GOARCH=$${arch} make build; \
		zip pkg/tmux-project_$(VERSION)_darwin_$${arch}.zip $(COMMANDS); \
	done;
	@# linux
	@for arch in "amd64" "386" "arm64" "arm"; do \
		GOOS=linux GOARCH=$${arch} make build; \
		zip pkg/tmux-project_$(VERSION)_linux_$${arch}.zip $(COMMANDS); \
	done;

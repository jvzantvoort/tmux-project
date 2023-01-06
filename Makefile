NAME := tmux-project
VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
COMMANDS := tmux-project pinfo
LDFLAGS := -X 'main.version=$(VERSION)' \
           -X 'main.revision=$(REVISION)' -w -s
GOIMPORTS ?= goimports
GOCILINT ?= golangci-lint
GO ?= GO111MODULE=on go
.DEFAULT_GOAL := help

.PHONY: fmt
fmt: ## Formatting source codes.
	find . -type f -name '*.go' -not -path '*/vendor/*' -exec $(GOIMPORTS) -w "{}" \;

.PHONY: clean
clean:
	@rm -f $(COMMANDS) || true; \
	rm -rf pkg || true; \
	rm -rf tags || true

.PHONY: refresh
refresh: tags
	@go-bindata -pkg tmuxproject messages; \
	pushd projecttype; \
	go-bindata -pkg projecttype templates; \
	popd

.PHONY: tags
tags:
	@find "$${PWD}" -type f -name '*.go' -not -path '*/vendor/*'| sed "s,$${PWD}/,," | xargs gotags >tags

.PHONY: pretty
pretty:
	@find "$${PWD}" -type f -name '*.go' -not -path '*/vendor/*' -exec goimports -w "{}" \;; \
	find "$${PWD}" -type f -name '*.go' -not -path '*/vendor/*' -exec gofmt -w "{}" \;

.PHONY: lint
lint: ## Run golint and go vet.
	@$(GOCILINT) run --no-config --disable-all --enable=goimports --enable=misspell *.go
	@$(GOCILINT) run --no-config --disable-all --enable=goimports --enable=misspell ./config/*.go
	@$(GOCILINT) run --no-config --disable-all --enable=goimports --enable=misspell ./cmd/*/*.go

.PHONY: test
test:
	@$(GO) test ./...

.PHONY: update
update:
	@test -e go.mod || $(GO) mod init
	@$(GO) mod tidy
	@$(GO) mod vendor

.PHONY: build
build: update
	$(foreach cmd,$(COMMANDS), $(GO) build -ldflags "$(LDFLAGS)" ./cmd/$(cmd);)

.PHONY: install
install:
	$(GO) install ./...

.PHONY: cross
cross: main.go  ## Build binaries for cross platform.
	mkdir -p pkg
	@# darwin
	@for arch in "amd64"; do \
		GOOS=darwin GOARCH=$${arch} make build; \
		zip pkg/tmux-project_$(VERSION)_darwin_$${arch}.zip $(COMMANDS); \
	done;
	@# linux
	@for arch in "amd64" "386" "arm64" "arm"; do \
		GOOS=linux GOARCH=$${arch} make build; \
		zip pkg/tmux-project_$(VERSION)_linux_$${arch}.zip $(COMMANDS); \
	done;

# renovate: datasource=github-releases depName=mvdan/gofumpt
GOFUMPT_PACKAGE_VERSION := v0.6.0
# renovate: datasource=github-releases depName=golangci/golangci-lint
GOLANGCI_LINT_PACKAGE_VERSION := v1.59.0

SHELL := bash
NAME := errors
IMPORT := github.com/owncloud-ops/$(NAME)
DIST := dist
DIST_DIRS := $(DIST)

GO ?= go
CWD ?= $(shell pwd)
PACKAGES ?= $(shell go list ./...)
SOURCES ?= $(shell find . -name "*.go" -type f)

GOFUMPT_PACKAGE ?= mvdan.cc/gofumpt@$(GOFUMPT_PACKAGE_VERSION)
GOLANGCI_LINT_PACKAGE ?= github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_PACKAGE_VERSION)
XGO_PACKAGE ?= src.techknowlogick.com/xgo@latest
GOTESTSUM_PACKAGE ?= gotest.tools/gotestsum@latest

XGO_VERSION := go-1.22.x
XGO_TARGETS ?= linux/amd64,linux/arm64

TAGS ?= netgo

ifndef VERSION
	ifneq ($(DRONE_TAG),)
		VERSION ?= $(subst v,,$(DRONE_TAG))
	else
		VERSION ?= $(shell git rev-parse --short HEAD)
	endif
endif

ifndef DATE
	DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%S%z")
endif

LDFLAGS += -s -w -X "$(IMPORT)/pkg/version.String=$(VERSION)" -X "$(IMPORT)/pkg/version.Date=$(DATE)"

.PHONY: all
all: build

.PHONY: clean
clean:
	$(GO) clean -i ./...
	rm -rf $(DIST_DIRS)

.PHONY: fmt
fmt:
	$(GO) run $(GOFUMPT_PACKAGE) -extra -w $(SOURCES)

.PHONY: golangci-lint
golangci-lint:
	$(GO) run $(GOLANGCI_LINT_PACKAGE) run

.PHONY: lint
lint: golangci-lint

.PHONY: test
test:
	$(GO) run $(GOTESTSUM_PACKAGE) -- -coverprofile=coverage.out $(PACKAGES)

.PHONY: build
build: $(DIST)/$(NAME)

$(DIST)/$(NAME): $(SOURCES)
	$(GO) build -v -tags '$(TAGS)' -ldflags '-extldflags "-static" $(LDFLAGS)' -o $@ ./cmd/$(NAME)

$(DIST_DIRS):
	mkdir -p $(DIST_DIRS)

.PHONY: xgo
xgo: | $(DIST_DIRS)
	$(GO) run $(XGO_PACKAGE) -go $(XGO_VERSION) -v -ldflags '-extldflags "-static" $(LDFLAGS)' -tags '$(TAGS)' -targets '$(XGO_TARGETS)' -out $(NAME) --pkg cmd/$(NAME) .
	cp /build/* $(CWD)/$(DIST)
	ls -l $(CWD)/$(DIST)

.PHONY: checksum
checksum:
	cd $(DIST); $(foreach file,$(wildcard $(DIST)/$(NAME)-*),sha256sum $(notdir $(file)) > $(notdir $(file)).sha256;)
	ls -l $(CWD)/$(DIST)

.PHONY: release
release: xgo checksum

.PHONY: deps
deps:
	$(GO) mod download
	$(GO) install $(GOFUMPT_PACKAGE)
	$(GO) install $(GOLANGCI_LINT_PACKAGE)
	$(GO) install $(XGO_PACKAGE)

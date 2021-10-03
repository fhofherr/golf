.DEFAULT_GOAL := test

GO ?= go

GO_FILES := $(shell find -iname '*.go')
GO_PACKAGES := $(shell $(GO) list ./... | paste -s -d ',')

.PHONY: help
help: ## Show this help.
	@awk -F ':|##' '/^[^\t].+?:.*?##/ {\
		printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF \
	}' $(MAKEFILE_LIST)

# ----------------------------------------------------------------------------
# Tools
# ----------------------------------------------------------------------------
TOOLS_GO := tools.go
TOOLS_DIR := .tools
TOOLS_BIN_DIR := $(TOOLS_DIR)/bin

.PHONY: tools
tools: $(TOOLS_BIN_DIR)/toolmgr ## Install all tools

$(TOOLS_BIN_DIR)/toolmgr:
	GOBIN=$(abspath ./$(TOOLS_BIN_DIR)) $(GO) install github.com/fhofherr/toolmgr

$(TOOLS_BIN_DIR)/%: $(TOOLS_GO) | $(TOOLS_BIN_DIR)/toolmgr
	$(TOOLS_BIN_DIR)/toolmgr -bin-dir $(TOOLS_BIN_DIR) -tools-go $<

.PRECIOUS: $(TOOLS_BIN_DIR)/%

# ----------------------------------------------------------------------------
# Tests
# ----------------------------------------------------------------------------
PRE_COMMIT ?= pre-commit

.PHONY: lint
lint: $(TOOLS_BIN_DIR)/revive ## Run linter on all files.
	$(PRE_COMMIT) run --all-files

.PHONY: test
test: .coverage.out ## Run all tests

.PHONY: coverage-html
coverage-html: .coverage.out ## Show HTML coverage report
	$(GO) tool cover -html=$<

.coverage.out: $(GO_FILES)
	@$(GO) test -race -coverpkg=$(GO_PACKAGES) -covermode=atomic -coverprofile=$@ ./...
	@$(GO) tool cover -func=$@ | tail -n1

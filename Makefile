# ted - Terminal EDitor Makefile
# Professional, user-friendly build system

# Project Information
PROJECT_NAME := ted
MODULE_NAME := github.com/AndrewDonelson/ted
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build Configuration
BINARY_NAME := ted
BINARY_DIR := bin
MAIN_PACKAGE := .

# Go Configuration
GO := go
GOFMT := gofmt
GOVET := go vet
GOLINT := golangci-lint
COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html

# Build Flags
LDFLAGS := -X $(MODULE_NAME)/internal/version.Version=$(VERSION) \
           -X $(MODULE_NAME)/internal/version.BuildTime=$(BUILD_TIME) \
           -X $(MODULE_NAME)/internal/version.GitCommit=$(GIT_COMMIT) \
           -s -w

# Colors for output
COLOR_RESET := \033[0m
COLOR_BOLD := \033[1m
COLOR_GREEN := \033[32m
COLOR_YELLOW := \033[33m
COLOR_BLUE := \033[34m
COLOR_CYAN := \033[36m

# Default target
.DEFAULT_GOAL := help

# Phony targets
.PHONY: help build install clean test test-coverage test-verbose test-race lint vet fmt fmt-check run dev deps deps-update deps-tidy cross-build all check

##@ General

help: ## Display this help message
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)ted - Terminal EDitor$(COLOR_RESET)"
	@echo "$(COLOR_BOLD)Available targets:$(COLOR_RESET)\n"
	@awk 'BEGIN {FS = ":.*##"; printf ""} /^[a-zA-Z_-]+:.*?##/ { printf "  $(COLOR_GREEN)%-20s$(COLOR_RESET) %s\n", $$1, $$2 } /^##@/ { printf "\n$(COLOR_BOLD)$(COLOR_BLUE)%s$(COLOR_RESET)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(COLOR_BOLD)Project:$(COLOR_RESET) $(PROJECT_NAME)"
	@echo "$(COLOR_BOLD)Version:$(COLOR_RESET) $(VERSION)"
	@echo "$(COLOR_BOLD)Module:$(COLOR_RESET) $(MODULE_NAME)"

##@ Building

build: ## Build the binary (default: bin/ted)
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)Building $(BINARY_NAME)...$(COLOR_RESET)"
	@mkdir -p $(BINARY_DIR)
	@$(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "$(COLOR_GREEN)✓ Build complete: $(BINARY_DIR)/$(BINARY_NAME)$(COLOR_RESET)"

install: ## Install the binary to $GOPATH/bin or $GOBIN
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)Installing $(BINARY_NAME)...$(COLOR_RESET)"
	@$(GO) install -ldflags "$(LDFLAGS)" $(MAIN_PACKAGE)
	@echo "$(COLOR_GREEN)✓ Installation complete$(COLOR_RESET)"

clean: ## Remove build artifacts
	@echo "$(COLOR_BOLD)$(COLOR_YELLOW)Cleaning...$(COLOR_RESET)"
	@rm -rf $(BINARY_DIR)
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	@$(GO) clean -cache -testcache -modcache
	@echo "$(COLOR_GREEN)✓ Clean complete$(COLOR_RESET)"

##@ Testing

test: ## Run all tests
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)Running tests...$(COLOR_RESET)"
	@$(GO) test ./... -v
	@echo "$(COLOR_GREEN)✓ Tests complete$(COLOR_RESET)"

test-coverage: ## Run tests with coverage report
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)Running tests with coverage...$(COLOR_RESET)"
	@$(GO) test ./... -coverprofile=$(COVERAGE_FILE) -covermode=atomic
	@$(GO) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "$(COLOR_GREEN)✓ Coverage report generated: $(COVERAGE_HTML)$(COLOR_RESET)"
	@$(GO) tool cover -func=$(COVERAGE_FILE) | tail -1

test-verbose: ## Run tests with verbose output
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)Running tests (verbose)...$(COLOR_RESET)"
	@$(GO) test ./... -v -cover

test-race: ## Run tests with race detector
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)Running tests with race detector...$(COLOR_RESET)"
	@$(GO) test ./... -race -v
	@echo "$(COLOR_GREEN)✓ Race detection complete$(COLOR_RESET)"

test-bench: ## Run benchmark tests
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)Running benchmarks...$(COLOR_RESET)"
	@$(GO) test ./... -bench=. -benchmem

##@ Code Quality

lint: ## Run all linters (requires golangci-lint)
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)Running linters...$(COLOR_RESET)"
	@if command -v $(GOLINT) > /dev/null; then \
		$(GOLINT) run ./...; \
	else \
		echo "$(COLOR_YELLOW)⚠ golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(COLOR_RESET)"; \
	fi

vet: ## Run go vet
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)Running go vet...$(COLOR_RESET)"
	@$(GOVET) ./...
	@echo "$(COLOR_GREEN)✓ go vet complete$(COLOR_RESET)"

fmt: ## Format code with gofmt
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)Formatting code...$(COLOR_RESET)"
	@$(GOFMT) -s -w .
	@echo "$(COLOR_GREEN)✓ Formatting complete$(COLOR_RESET)"

fmt-check: ## Check if code is formatted correctly
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)Checking code format...$(COLOR_RESET)"
	@if [ $$($(GOFMT) -l . | wc -l) -ne 0 ]; then \
		echo "$(COLOR_YELLOW)⚠ Code is not formatted. Run 'make fmt' to fix.$(COLOR_RESET)"; \
		$(GOFMT) -d .; \
		exit 1; \
	else \
		echo "$(COLOR_GREEN)✓ Code is properly formatted$(COLOR_RESET)"; \
	fi

check: fmt-check vet test ## Run all checks (format, vet, tests)
	@echo "$(COLOR_GREEN)✓ All checks passed$(COLOR_RESET)"

##@ Development

run: build ## Build and run the binary
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)Running $(BINARY_NAME)...$(COLOR_RESET)"
	@./$(BINARY_DIR)/$(BINARY_NAME) $(ARGS)

dev: ## Run in development mode (with auto-rebuild on file changes)
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)Starting development mode...$(COLOR_RESET)"
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "$(COLOR_YELLOW)⚠ air not found. Install with: go install github.com/cosmtrek/air@latest$(COLOR_RESET)"; \
		echo "$(COLOR_YELLOW)  Falling back to regular build...$(COLOR_RESET)"; \
		$(MAKE) run; \
	fi

##@ Dependencies

deps: ## Download dependencies
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)Downloading dependencies...$(COLOR_RESET)"
	@$(GO) mod download
	@echo "$(COLOR_GREEN)✓ Dependencies downloaded$(COLOR_RESET)"

deps-update: ## Update all dependencies
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)Updating dependencies...$(COLOR_RESET)"
	@$(GO) get -u ./...
	@$(GO) mod tidy
	@echo "$(COLOR_GREEN)✓ Dependencies updated$(COLOR_RESET)"

deps-tidy: ## Tidy dependencies
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)Tidying dependencies...$(COLOR_RESET)"
	@$(GO) mod tidy
	@echo "$(COLOR_GREEN)✓ Dependencies tidied$(COLOR_RESET)"

deps-vendor: ## Vendor dependencies
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)Vendoring dependencies...$(COLOR_RESET)"
	@$(GO) mod vendor
	@echo "$(COLOR_GREEN)✓ Dependencies vendored$(COLOR_RESET)"

##@ Cross-Platform Builds

cross-build: ## Build for all supported platforms
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)Building for all platforms...$(COLOR_RESET)"
	@mkdir -p $(BINARY_DIR)
	@echo "$(COLOR_CYAN)Building for Linux (amd64)...$(COLOR_RESET)"
	@GOOS=linux GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	@echo "$(COLOR_CYAN)Building for Linux (arm64)...$(COLOR_RESET)"
	@GOOS=linux GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PACKAGE)
	@echo "$(COLOR_CYAN)Building for macOS (amd64)...$(COLOR_RESET)"
	@GOOS=darwin GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	@echo "$(COLOR_CYAN)Building for macOS (arm64)...$(COLOR_RESET)"
	@GOOS=darwin GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)
	@echo "$(COLOR_CYAN)Building for Windows (amd64)...$(COLOR_RESET)"
	@GOOS=windows GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)
	@echo "$(COLOR_GREEN)✓ Cross-platform builds complete$(COLOR_RESET)"
	@ls -lh $(BINARY_DIR)/

##@ Release

release: clean cross-build ## Prepare release builds
	@echo "$(COLOR_BOLD)$(COLOR_GREEN)✓ Release builds ready in $(BINARY_DIR)/$(COLOR_RESET)"

##@ Information

info: ## Display project information
	@echo "$(COLOR_BOLD)Project Information:$(COLOR_RESET)"
	@echo "  Name:    $(PROJECT_NAME)"
	@echo "  Version: $(VERSION)"
	@echo "  Module:  $(MODULE_NAME)"
	@echo "  Commit:  $(GIT_COMMIT)"
	@echo "  Build:   $(BUILD_TIME)"
	@echo ""
	@echo "$(COLOR_BOLD)Go Information:$(COLOR_RESET)"
	@$(GO) version
	@echo ""
	@echo "$(COLOR_BOLD)Build Information:$(COLOR_RESET)"
	@echo "  Binary:  $(BINARY_DIR)/$(BINARY_NAME)"
	@echo "  Main:    $(MAIN_PACKAGE)"

version: ## Display version information
	@echo "$(VERSION)"

##@ All-in-One

all: clean deps-tidy fmt-check vet test build ## Run full build pipeline (clean, deps, format, vet, test, build)
	@echo "$(COLOR_BOLD)$(COLOR_GREEN)✓ Full build pipeline complete$(COLOR_RESET)"


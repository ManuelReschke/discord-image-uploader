# Discord Image Uploader - Makefile
# Go build automation for cross-platform compilation

# Variables
APP_NAME := discord-image-uploader
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GO_VERSION := $(shell go version | cut -d ' ' -f 3)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Directories
BIN_DIR := bin
CMD_DIR := cmd
DATA_DIR := data
CONFIG_DIR := config

# Build flags
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -s -w"

# Platform targets
PLATFORMS := \
	windows/amd64 \
	windows/386 \
	linux/amd64 \
	linux/386 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64

# Colors for output
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
MAGENTA := \033[35m
CYAN := \033[36m
WHITE := \033[37m
RESET := \033[0m

# Default target
.PHONY: all
all: clean build

# Print fancy banner
.PHONY: banner
banner:
	@echo "$(CYAN)"
	@echo "╔══════════════════════════════════════════════════════════════╗"
	@echo "║                 Discord Image Uploader                       ║"
	@echo "║                   Build System v2.0                         ║"
	@echo "╚══════════════════════════════════════════════════════════════╝"
	@echo "$(RESET)"
	@echo "$(WHITE)Version:    $(GREEN)$(VERSION)$(RESET)"
	@echo "$(WHITE)Go Version: $(GREEN)$(GO_VERSION)$(RESET)"
	@echo "$(WHITE)Git Commit: $(GREEN)$(GIT_COMMIT)$(RESET)"
	@echo "$(WHITE)Build Time: $(GREEN)$(BUILD_TIME)$(RESET)"
	@echo ""

# Help target
.PHONY: help
help:
	@echo "$(CYAN)Available targets:$(RESET)"
	@echo "  $(GREEN)build$(RESET)        - Build for current platform"
	@echo "  $(GREEN)build-all$(RESET)    - Build for all platforms"
	@echo "  $(GREEN)build-linux$(RESET)  - Build for Linux (amd64, 386, arm64)"
	@echo "  $(GREEN)build-windows$(RESET) - Build for Windows (amd64, 386)"
	@echo "  $(GREEN)build-mac$(RESET)    - Build for macOS (amd64, arm64)"
	@echo "  $(GREEN)clean$(RESET)        - Clean build artifacts"
	@echo "  $(GREEN)test$(RESET)         - Run tests"
	@echo "  $(GREEN)fmt$(RESET)          - Format code"
	@echo "  $(GREEN)lint$(RESET)         - Run linter"
	@echo "  $(GREEN)deps$(RESET)         - Download dependencies"
	@echo "  $(GREEN)tidy$(RESET)         - Clean up go.mod"
	@echo "  $(GREEN)install$(RESET)      - Install to GOPATH/bin"
	@echo "  $(GREEN)package$(RESET)      - Create release packages"
	@echo "  $(GREEN)dev$(RESET)          - Quick build for development"
	@echo "  $(GREEN)run$(RESET)          - Build and run with config"
	@echo "  $(GREEN)release$(RESET)      - Full release build"
	@echo "  $(GREEN)docker$(RESET)       - Build Docker image"
	@echo "  $(GREEN)help$(RESET)         - Show this help"

# Development build (fast, current platform only)
.PHONY: dev
dev: banner
	@echo "$(YELLOW)🔨 Building for development...$(RESET)"
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) ./$(CMD_DIR)
	@echo "$(GREEN)✅ Development build complete: $(BIN_DIR)/$(APP_NAME)$(RESET)"

# Build for current platform
.PHONY: build
build: banner deps
	@echo "$(YELLOW)🔨 Building for current platform...$(RESET)"
	@mkdir -p $(BIN_DIR)
	go build $(LDFLAGS) -o $(BIN_DIR)/$(APP_NAME) ./$(CMD_DIR)
	@echo "$(GREEN)✅ Build complete: $(BIN_DIR)/$(APP_NAME)$(RESET)"

# Build for all platforms
.PHONY: build-all
build-all: banner deps
	@echo "$(YELLOW)🌍 Building for all platforms...$(RESET)"
	@mkdir -p $(BIN_DIR)
	@$(foreach platform,$(PLATFORMS), \
		$(call build_platform,$(platform)))
	@echo "$(GREEN)✅ All builds complete!$(RESET)"
	@echo "$(CYAN)📦 Built binaries:$(RESET)"
	@ls -la $(BIN_DIR)/

# Build for Linux platforms
.PHONY: build-linux
build-linux: banner deps
	@echo "$(YELLOW)🐧 Building for Linux platforms...$(RESET)"
	@mkdir -p $(BIN_DIR)
	@$(call build_platform,linux/amd64)
	@$(call build_platform,linux/386)
	@$(call build_platform,linux/arm64)
	@echo "$(GREEN)✅ Linux builds complete!$(RESET)"

# Build for Windows platforms
.PHONY: build-windows
build-windows: banner deps
	@echo "$(YELLOW)🪟 Building for Windows platforms...$(RESET)"
	@mkdir -p $(BIN_DIR)
	@$(call build_platform,windows/amd64)
	@$(call build_platform,windows/386)
	@echo "$(GREEN)✅ Windows builds complete!$(RESET)"

# Build for macOS platforms
.PHONY: build-mac
build-mac: banner deps
	@echo "$(YELLOW)🍎 Building for macOS platforms...$(RESET)"
	@mkdir -p $(BIN_DIR)
	@$(call build_platform,darwin/amd64)
	@$(call build_platform,darwin/arm64)
	@echo "$(GREEN)✅ macOS builds complete!$(RESET)"

# Function to build for specific platform
define build_platform
	$(eval GOOS := $(word 1,$(subst /, ,$1)))
	$(eval GOARCH := $(word 2,$(subst /, ,$1)))
	$(eval EXT := $(if $(filter windows,$(GOOS)),.exe,))
	$(eval FILENAME := $(APP_NAME)-$(GOOS)-$(GOARCH)$(EXT))
	@echo "$(BLUE)  Building $(GOOS)/$(GOARCH)...$(RESET)"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(BIN_DIR)/$(FILENAME) ./$(CMD_DIR)
	@echo "$(GREEN)    ✓ $(FILENAME)$(RESET)"
endef

# Clean build artifacts
.PHONY: clean
clean:
	@echo "$(YELLOW)🧹 Cleaning build artifacts...$(RESET)"
	@rm -rf $(BIN_DIR)
	@rm -rf $(DATA_DIR)/*.json
	@rm -rf dist/
	@echo "$(GREEN)✅ Clean complete!$(RESET)"

# Download dependencies
.PHONY: deps
deps:
	@echo "$(YELLOW)📦 Downloading dependencies...$(RESET)"
	@go mod download
	@echo "$(GREEN)✅ Dependencies downloaded!$(RESET)"

# Tidy up go.mod
.PHONY: tidy
tidy:
	@echo "$(YELLOW)🧹 Tidying go.mod...$(RESET)"
	@go mod tidy
	@echo "$(GREEN)✅ go.mod tidied!$(RESET)"

# Format code
.PHONY: fmt
fmt:
	@echo "$(YELLOW)📝 Formatting code...$(RESET)"
	@go fmt ./...
	@echo "$(GREEN)✅ Code formatted!$(RESET)"

# Run linter
.PHONY: lint
lint:
	@echo "$(YELLOW)🔍 Running linter...$(RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(RED)⚠️  golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(RESET)"; \
		go vet ./...; \
	fi
	@echo "$(GREEN)✅ Linting complete!$(RESET)"

# Run tests
.PHONY: test
test:
	@echo "$(YELLOW)🧪 Running tests...$(RESET)"
	@go test -v ./...
	@echo "$(GREEN)✅ Tests complete!$(RESET)"

# Install to GOPATH/bin
.PHONY: install
install: deps
	@echo "$(YELLOW)📥 Installing to GOPATH/bin...$(RESET)"
	@go install $(LDFLAGS) ./$(CMD_DIR)
	@echo "$(GREEN)✅ Installed to $(shell go env GOPATH)/bin/$(APP_NAME)$(RESET)"

# Build and run with default config
.PHONY: run
run: build
	@echo "$(YELLOW)🚀 Running application...$(RESET)"
	@if [ ! -f $(CONFIG_DIR)/config.json ]; then \
		echo "$(RED)⚠️  Config file not found. Copying example...$(RESET)"; \
		cp $(CONFIG_DIR)/config.example.json $(CONFIG_DIR)/config.json; \
		echo "$(YELLOW)Please edit $(CONFIG_DIR)/config.json with your settings$(RESET)"; \
		exit 1; \
	fi
	@./$(BIN_DIR)/$(APP_NAME) -config $(CONFIG_DIR)/config.json

# Create release packages
.PHONY: package
package: clean build-all
	@echo "$(YELLOW)📦 Creating release packages...$(RESET)"
	@mkdir -p dist
	@for binary in $(BIN_DIR)/*; do \
		if [ -f "$$binary" ]; then \
			base=$$(basename "$$binary"); \
			echo "$(BLUE)  Packaging $$base...$(RESET)"; \
			mkdir -p "dist/$$base"; \
			cp "$$binary" "dist/$$base/"; \
			cp -r $(CONFIG_DIR) "dist/$$base/"; \
			cp README.md "dist/$$base/"; \
			cp LICENSE "dist/$$base/" 2>/dev/null || true; \
			cd dist && tar -czf "$$base.tar.gz" "$$base/" && cd ..; \
			echo "$(GREEN)    ✓ dist/$$base.tar.gz$(RESET)"; \
		fi \
	done
	@echo "$(GREEN)✅ Release packages complete!$(RESET)"

# Full release build
.PHONY: release
release: banner clean fmt lint test build-all package
	@echo "$(GREEN)"
	@echo "🎉 Release build complete!"
	@echo "   Version: $(VERSION)"
	@echo "   Binaries: $(BIN_DIR)/"
	@echo "   Packages: dist/"
	@echo "$(RESET)"

# Docker build
.PHONY: docker
docker:
	@echo "$(YELLOW)🐳 Building Docker image...$(RESET)"
	@docker build -t $(APP_NAME):$(VERSION) .
	@docker tag $(APP_NAME):$(VERSION) $(APP_NAME):latest
	@echo "$(GREEN)✅ Docker image built: $(APP_NAME):$(VERSION)$(RESET)"

# Show build info
.PHONY: info
info: banner
	@echo "$(CYAN)Build Information:$(RESET)"
	@echo "  App Name:     $(APP_NAME)"
	@echo "  Version:      $(VERSION)"
	@echo "  Git Commit:   $(GIT_COMMIT)"
	@echo "  Build Time:   $(BUILD_TIME)"
	@echo "  Go Version:   $(GO_VERSION)"
	@echo ""
	@echo "$(CYAN)Directories:$(RESET)"
	@echo "  Binary Dir:   $(BIN_DIR)/"
	@echo "  Source Dir:   $(CMD_DIR)/"
	@echo "  Config Dir:   $(CONFIG_DIR)/"
	@echo ""
	@echo "$(CYAN)Supported Platforms:$(RESET)"
	@$(foreach platform,$(PLATFORMS),echo "  $(platform)";)

# Development workflow shortcuts
.PHONY: quick
quick: clean dev run

.PHONY: check
check: fmt lint test

# Setup development environment
.PHONY: setup
setup:
	@echo "$(YELLOW)🔧 Setting up development environment...$(RESET)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go mod tidy
	@mkdir -p $(BIN_DIR) $(DATA_DIR)
	@if [ ! -f $(CONFIG_DIR)/config.json ]; then \
		cp $(CONFIG_DIR)/config.example.json $(CONFIG_DIR)/config.json; \
		echo "$(YELLOW)Please edit $(CONFIG_DIR)/config.json with your settings$(RESET)"; \
	fi
	@echo "$(GREEN)✅ Development environment ready!$(RESET)"

# Show git status and pending changes
.PHONY: status
status:
	@echo "$(CYAN)Git Status:$(RESET)"
	@git status --short
	@echo ""
	@echo "$(CYAN)Last 5 commits:$(RESET)"
	@git log --oneline -5
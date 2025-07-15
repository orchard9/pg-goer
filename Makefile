# pg-goer Makefile

BINARY_NAME=pg-goer
MAIN_PACKAGE=./cmd/pg-goer
GO=go
GOLANGCI_LINT=golangci-lint

# Version information
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse HEAD)
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)

# Build the binary
.PHONY: build
build:
	$(GO) build -o $(BINARY_NAME) $(MAIN_PACKAGE)

# Build the binary with version information
.PHONY: build-release
build-release:
	$(GO) build -ldflags="$(LDFLAGS)" -o $(BINARY_NAME) $(MAIN_PACKAGE)

# Run CI checks (tests + linting)
.PHONY: ci
ci: test lint

# Run tests
.PHONY: test
test:
	$(GO) test -v -race -coverprofile=coverage.out ./...

# Run linter
.PHONY: lint
lint:
	$(GOLANGCI_LINT) run

# Clean build artifacts
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)
	rm -f coverage.out

# Install dependencies
.PHONY: deps
deps:
	$(GO) mod download
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run the application
.PHONY: run
run:
	$(GO) run $(MAIN_PACKAGE)

# Format code
.PHONY: fmt
fmt:
	$(GO) fmt ./...

# Vet code
.PHONY: vet
vet:
	$(GO) vet ./...

# UAT: Start PostgreSQL database and build binary
.PHONY: uat-up
uat-up: build
	@echo "Starting UAT PostgreSQL database..."
	cd uat && docker-compose up -d
	cd uat && ./wait-for-postgres.sh

# UAT: Run tests against database
.PHONY: uat-run
uat-run:
	cd uat && ./run-tests.sh

# UAT: Stop database and clean up
.PHONY: uat-down
uat-down:
	@echo "Stopping UAT database and cleaning up..."
	cd uat && docker-compose down -v --remove-orphans 2>/dev/null || true
	rm -f uat/uat-test-output.md uat/schema-filtered-output.md

# UAT: Complete cycle (up + run + down)
.PHONY: uat
uat: uat-up uat-run uat-down

# Clean UAT artifacts (alias for uat-down)
.PHONY: uat-clean
uat-clean: uat-down

# Run integration tests
.PHONY: integration
integration:
	$(GO) test -v -tags=integration ./tests/integration/...

# Clean integration test artifacts
.PHONY: integration-clean
integration-clean:
	cd tests/integration && docker-compose down -v --remove-orphans 2>/dev/null || true

# Run all tests (unit + integration)
.PHONY: test-all
test-all: test integration

# Clean all test artifacts
.PHONY: clean-all
clean-all: clean uat-clean integration-clean

# Release targets
# ===============

# Cross-platform build targets
.PHONY: build-linux-amd64
build-linux-amd64:
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags="$(LDFLAGS)" -o $(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)

.PHONY: build-linux-arm64
build-linux-arm64:
	GOOS=linux GOARCH=arm64 $(GO) build -ldflags="$(LDFLAGS)" -o $(BINARY_NAME)-linux-arm64 $(MAIN_PACKAGE)

.PHONY: build-darwin-amd64
build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags="$(LDFLAGS)" -o $(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)

.PHONY: build-darwin-arm64
build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 $(GO) build -ldflags="$(LDFLAGS)" -o $(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)

.PHONY: build-windows-amd64
build-windows-amd64:
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags="$(LDFLAGS)" -o $(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)

# Build all platforms
.PHONY: build-all
build-all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64

# Create release archives
.PHONY: package
package: build-all
	@echo "Creating release packages..."
	@mkdir -p dist
	@# Linux amd64
	@tar -czf dist/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	@# Linux arm64
	@tar -czf dist/$(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64
	@# macOS amd64
	@tar -czf dist/$(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	@# macOS arm64
	@tar -czf dist/$(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64
	@# Windows amd64
	@zip dist/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	@echo "Release packages created in dist/"

# Full release preparation
.PHONY: release-prep
release-prep: clean ci build-all package
	@echo "Release preparation complete!"
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"
	@echo ""
	@echo "Release artifacts:"
	@ls -la dist/

# Clean release artifacts
.PHONY: clean-release
clean-release:
	rm -f $(BINARY_NAME)-*
	rm -rf dist/

# Create a new release tag
.PHONY: tag
tag:
	@if [ -z "$(TAG)" ]; then echo "Usage: make tag TAG=v1.0.0"; exit 1; fi
	@git tag -a $(TAG) -m "Release $(TAG)"
	@echo "Created tag $(TAG)"
	@echo "Push with: git push origin $(TAG)"

# Show version information
.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"

# Validate CI configuration locally
.PHONY: validate-ci
validate-ci:
	@./scripts/validate-ci.sh

# Complete release workflow
.PHONY: release
release: release-prep
	@echo ""
	@echo "🎉 Release $(VERSION) is ready!"
	@echo ""
	@echo "Next steps:"
	@echo "1. Review the release artifacts in dist/"
	@echo "2. Create a git tag: make tag TAG=v$(VERSION)"
	@echo "3. Push the tag: git push origin v$(VERSION)"
	@echo "4. GitHub Actions will automatically create the release"
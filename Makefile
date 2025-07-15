# pg-goer Makefile

BINARY_NAME=pg-goer
MAIN_PACKAGE=./cmd/pg-goer
GO=go
GOLANGCI_LINT=golangci-lint

# Build the binary
.PHONY: build
build:
	$(GO) build -o $(BINARY_NAME) $(MAIN_PACKAGE)

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
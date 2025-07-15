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
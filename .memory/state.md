# Project State

## Documentation Files
- README.md: Main project documentation and quick start guide
- CLAUDE.md: Development guide for AI-assisted coding with memory management
- LICENSE.md: MIT license text
- coding_guidelines.md: Go coding standards and best practices
- code_architecture.md: System design and component overview
- usage.md: Detailed usage instructions and examples
- why.md: Project rationale and benefits
- contributing.md: Contribution guidelines for developers
- example-output.md: Live example markdown documentation (5,459 bytes, 215 lines)
- example-output.json: Live example JSON documentation (11,864 bytes structured data)

## Build and Configuration
- Makefile: Build automation with ci, build, test, lint, and utility targets
- .golangci.yml: golangci-lint configuration for code quality
- go.mod: Go module definition with PostgreSQL driver dependency
- go.sum: Go module checksums
- coverage.out: Test coverage output file

## Source Code
- cmd/pg-goer/main.go: CLI entry point with connection string, output file, and schema filtering flags
- cmd/pg-goer/main_test.go: Unit tests for CLI functionality

## Core Models
- pkg/models/schema.go: Comprehensive data structures for Schema, Table, Column, ForeignKey, Index, Trigger, Extension, View, and Sequence

## Database Analysis
- internal/analyzer/connection.go: PostgreSQL connection management with proper pooling and timeouts
- internal/analyzer/connection_test.go: Tests for connection handling
- internal/analyzer/schema.go: Complete PostgreSQL schema analysis with generic query helpers for tables, columns, foreign keys, indexes, triggers, extensions, views, sequences, and row counts
- internal/analyzer/schema_test.go: Comprehensive tests for all schema analysis functionality

## Report Generation
- internal/reporter/markdown.go: Enhanced markdown report generator with TOC, database summary, professional formatting, and comprehensive PostgreSQL object documentation
- internal/reporter/markdown_test.go: Tests for markdown generation including all new formatting features
- internal/reporter/json.go: Complete JSON report generator with structured output for programmatic consumption
- internal/reporter/json_test.go: Tests for JSON generation and structure validation
- internal/generator/mermaid.go: Mermaid ER diagram generator creating visual database relationships with proper schema handling
- internal/generator/mermaid_test.go: Comprehensive tests for Mermaid diagram generation with 97% coverage

## Memory Files
- .memory/state.md: Current file listing with descriptions
- .memory/working-memory.md: Recent, current, and future work tracking
- .memory/semantic-memory.md: Project facts and characteristics
- .memory/vision.md: Long-term project vision and principles
- .memory/tasks.md: Comprehensive task list with descriptions and completion tracking

## User Acceptance Testing
- uat/docker-compose.yml: PostgreSQL 15 test database configuration with custom test data
- uat/init/01-setup-database.sql: Comprehensive e-commerce test schema with 5 tables, foreign keys, and sample data
- uat/run-tests.sh: Enhanced UAT script validating all functionality and generating example outputs
- uat/wait-for-postgres.sh: Database readiness validation script
- uat/test-uat.sh: Legacy UAT script (preserved for compatibility)
- uat/expected-output/sample-db-docs.md: Reference documentation for validation
- uat/README.md: UAT documentation and usage instructions

## Integration Testing
- tests/integration/docker-compose.yml: Integration test PostgreSQL configuration
- tests/integration/integration_test.go: Docker-based integration tests
- tests/integration/testdata/01-basic-schema.sql: Test schema for integration testing
- tests/integration/README.md: Integration testing documentation

## CI/CD and Scripts
- .github/workflows/ci.yml: GitHub Actions pipeline with testing, linting, security scanning, UAT, and multi-platform builds
- scripts/validate-ci.sh: CI validation script

## Directories
- internal/: Core application logic with analyzer, reporter, and generator packages
- pkg/: Public packages with data models
- cmd/: Command-line application entry point
- uat/: User Acceptance Testing with comprehensive Docker PostgreSQL validation
- tests/: Integration testing framework with Docker
- docs/: Additional documentation
- .memory/: Memory management files for development context
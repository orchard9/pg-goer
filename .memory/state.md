# Project State

## Documentation Files
- README.md: Main project documentation and quick start guide
- CLAUDE.md: Development guide for AI-assisted coding
- LICENSE.md: MIT license text
- coding_guidelines.md: Go coding standards and best practices
- code_architecture.md: System design and component overview
- usage.md: Detailed usage instructions and examples
- why.md: Project rationale and benefits
- contributing.md: Contribution guidelines for developers

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
- pkg/models/schema.go: Data structures for Schema, Table, Column, and ForeignKey

## Database Analysis
- internal/analyzer/connection.go: PostgreSQL connection management with proper pooling and timeouts
- internal/analyzer/connection_test.go: Tests for connection handling
- internal/analyzer/schema.go: Schema analysis logic with tables, columns, foreign keys, and row counts from information_schema
- internal/analyzer/schema_test.go: Tests for schema analysis functionality including foreign key and row count queries

## Report Generation
- internal/reporter/markdown.go: Enhanced markdown report generator with TOC, database summary, and professional formatting
- internal/reporter/markdown_test.go: Tests for markdown generation including new formatting features
- internal/generator/mermaid.go: Mermaid ER diagram generator creating visual database relationships
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
- uat/test-uat.sh: Automated UAT script validating CLI, database connectivity, schema analysis, and output generation
- uat/expected-output/sample-db-docs.md: Reference documentation for validation
- uat/README.md: UAT documentation and usage instructions
- uat/.gitignore: UAT-specific ignore patterns

## Directories
- internal/generator/: Mermaid ER diagram generation with comprehensive test coverage
- docs/: Additional documentation (empty)
- tests/: Test files (empty, integration tests planned)
- uat/: User Acceptance Testing with Docker PostgreSQL
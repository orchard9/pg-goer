# Project Tasks

## Pending Tasks

### 1. Create Makefile with make ci and make build targets
Set up the build system foundation. make ci runs tests and linting, make build creates the binary.
This is critical - we'll run make ci after every change to ensure quality.

### 2. Set up basic Go project linting and testing in make ci
Configure golangci-lint and go test in the ci target.
Essential for maintaining code quality throughout development.

### 3. Create pkg/models/schema.go with Table and Column structs
Define core data models that represent database schema.
These structs will be used throughout the application for type safety.

### 4. Create internal/analyzer/connection.go with Connect function
Implement PostgreSQL connection handling with proper error management.
Must handle connection strings, timeouts, and connection pooling.

### 5. Write test for connection string parsing
Test various PostgreSQL connection string formats before implementation.
Following TDD ensures we handle edge cases properly.

### 6. Implement connection string parsing
Parse PostgreSQL connection strings to extract host, port, database, credentials.
Support both URI and key-value formats.

### 7. Query information_schema for tables list
First actual database query to get all tables in the database.
This is the foundation of our schema analysis.

### 8. Query columns for each table
Retrieve column details: name, type, nullable, defaults for each table.
Core functionality for documenting table structures.

### 9. Create internal/reporter/markdown.go with basic output
Build markdown generation logic starting with simple table documentation.
This creates our first visible output.

### 10. Wire up CLI to analyzer and reporter
Connect the CLI flags to the analyzer and reporter components.
Makes the tool actually functional end-to-end.

### 11. Add foreign key constraint queries
Query PostgreSQL for foreign key relationships between tables.
Essential for generating accurate ER diagrams.

### 12. Create internal/generator/mermaid.go for diagrams
Implement Mermaid diagram syntax generation from schema data.
Provides visual representation of database relationships.

### 13. Generate basic Mermaid ER diagram
Convert foreign key data into Mermaid ER diagram syntax.
Makes relationships visually clear in documentation.

### 14. Add table row count queries
Efficiently query row counts for all tables (using pg_stat_user_tables).
Provides useful size metrics without full table scans.

### 15. Format complete markdown report
Combine all components into a polished markdown document.
Include TOC, sections for tables, relationships, and diagrams.

### 16. Add integration test with Docker PostgreSQL
Set up Docker-based testing against real PostgreSQL instances.
Ensures our tool works with actual databases, not just mocks.

### 17. Add JSON output format
Implement JSON report generation as alternative to markdown.
Enables programmatic consumption of the documentation.

### 18. Set up GitHub Actions CI
Configure automated testing and building on every push.
Maintains quality and provides pre-built binaries.

### 19. Create release process in Makefile
Add make release target for creating versioned binaries.
Automates the release process for consistency.

### 20. Build and test final binary
Final validation of the complete tool with all features.
Ensure single binary works as designed with no dependencies.

## Completed Tasks

<!-- Completed tasks will be moved here with implementation notes -->
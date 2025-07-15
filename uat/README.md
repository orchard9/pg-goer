# User Acceptance Testing (UAT) for pg-goer

This directory contains User Acceptance Tests that validate pg-goer against real PostgreSQL databases using Docker.

## Overview

The UAT suite creates a realistic test environment with:
- PostgreSQL 15 database in Docker
- Comprehensive e-commerce-like schema with 5 tables
- Foreign key relationships between tables
- Sample data (40+ records)
- Validation of all pg-goer features

## Quick Start

```bash
# Complete UAT cycle (recommended)
make uat

# Or run individual steps for development/debugging:
make uat-up    # Start database and build binary
make uat-run   # Run tests (can repeat multiple times)
make uat-down  # Stop database and clean up
```

## Split Commands for Reliability

The UAT is split into three focused commands following our core values:

### `make uat-up`
- **Reliable**: Builds binary and starts PostgreSQL with health checks
- **Elegant**: Single responsibility - environment setup
- **Efficient**: Only does setup, can be run once for multiple test iterations

### `make uat-run` 
- **Reliable**: Assumes database is ready, focused error messages
- **Elegant**: Pure testing logic without Docker complexity
- **Efficient**: Fast execution, no Docker overhead

### `make uat-down`
- **Reliable**: Always works, even if containers aren't running
- **Elegant**: Simple cleanup with clear status messages  
- **Efficient**: Quick shutdown and file cleanup

This design enables:
- **Development**: Start database once (`uat-up`), iterate on tests (`uat-run`)
- **CI/CD**: Each step can be cached, parallelized, and have clear failure points
- **Debugging**: Can inspect database state between steps

## Test Structure

```
uat/
├── docker-compose.yml          # PostgreSQL service configuration
├── init/
│   └── 01-setup-database.sql   # Test schema and sample data
├── test-uat.sh                 # Main UAT test script
├── expected-output/
│   └── sample-db-docs.md       # Reference output for validation
└── README.md                   # This file
```

## Test Database Schema

The test database includes:

- **users**: User accounts with email, names, status
- **orders**: Purchase orders linked to users
- **order_items**: Individual items within orders
- **categories**: Product categories with hierarchical structure
- **products**: Product catalog linked to categories

This schema tests:
- Foreign key relationships
- Various PostgreSQL data types
- Default values and constraints
- Primary and unique keys
- Generated columns
- Self-referencing relationships

## What the UAT Tests

1. **Binary Compilation**: Ensures `make build` works
2. **CLI Functionality**: Tests `--help` and `--version` flags
3. **Database Connection**: Validates PostgreSQL connectivity
4. **Schema Analysis**: Verifies table, column, and relationship discovery
5. **Documentation Generation**: Confirms markdown output creation
6. **Content Validation**: Checks for required sections and data
7. **Foreign Key Documentation**: Validates relationship documentation
8. **Mermaid Diagrams**: Ensures ER diagrams are generated
9. **Row Count Analysis**: Confirms table statistics inclusion
10. **Schema Filtering**: Tests `--schemas` flag functionality

## Manual Testing

You can also run manual tests:

```bash
# Start the test database
docker-compose up -d

# Wait for startup, then test manually
../pg-goer "postgresql://testuser:testpass@localhost:5432/testdb"

# View the generated documentation
cat database-docs.md

# Clean up
docker-compose down -v
```

## Expected Output

The UAT validates that the generated documentation includes:

- Complete table of contents with anchor links
- Database summary with table and row counts
- Mermaid ER diagram showing all relationships
- Detailed table documentation with columns, types, and constraints
- Foreign key relationships properly documented
- Row counts for statistical analysis

## Troubleshooting

### Docker Issues
- Ensure Docker is running
- Check port 5432 is available
- Try `docker-compose down -v` to clean state

### Connection Issues
- Verify PostgreSQL health check passes
- Check firewall settings
- Ensure connection string format is correct

### Test Failures
- Review test output for specific error messages
- Check generated documentation content
- Verify all expected tables and relationships exist

## Integration with CI

This UAT can be integrated into CI pipelines:

```bash
# Add to Makefile
make uat

# Or run directly in CI
cd uat && ./test-uat.sh
```

The script exits with status 0 on success, non-zero on failure, making it suitable for automated testing.
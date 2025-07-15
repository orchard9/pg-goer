# Integration Tests

This directory contains integration tests for pg-goer that run against a real PostgreSQL database using Docker.

## Overview

The integration tests validate the complete functionality of pg-goer by:

1. Starting a PostgreSQL database in Docker
2. Loading test schema and data
3. Testing database connection and schema analysis
4. Validating documentation generation
5. Verifying all features work together

## Running Tests

```bash
# Run integration tests
make integration

# Clean up after tests
make integration-clean
```

## Test Structure

- `docker-compose.yml`: PostgreSQL 15 Alpine configuration for fast startup
- `testdata/01-basic-schema.sql`: Test schema with comprehensive feature coverage
- `integration_test.go`: Go integration tests with build tag `integration`

## Test Database Schema

The test database includes:

- **users**: Primary table with various column types and constraints
- **posts**: Content table with foreign key to users
- **comments**: Comments with foreign keys to both users and posts
- **categories**: Hierarchical categories with self-referencing foreign key
- **post_categories**: Junction table for many-to-many relationships

This schema covers:
- Primary keys and unique constraints
- Foreign key relationships (one-to-many, many-to-many, self-referencing)
- Various PostgreSQL data types (VARCHAR, TEXT, INTEGER, BOOLEAN, TIMESTAMP, JSONB)
- Default values and CHECK constraints
- Composite primary keys

## Test Coverage

The integration tests validate:

1. **Database Connection**: Basic connectivity and authentication
2. **Schema Analysis**: Table, column, and constraint discovery
3. **Foreign Key Analysis**: Relationship detection and documentation
4. **Row Count Analysis**: Table statistics gathering
5. **Complete Workflow**: End-to-end documentation generation
6. **Content Validation**: Ensures all expected sections and relationships are documented

## CI/CD Integration

These tests are designed for CI/CD pipelines:

- Use Alpine PostgreSQL for fast startup
- Include health checks for reliable container startup
- Automatically clean up Docker resources
- Exit with proper status codes for CI integration
- Use build tags to separate from unit tests

## Performance

Integration tests are optimized for speed:

- Minimal test data (3-4 rows per table)
- Alpine PostgreSQL image
- In-memory storage (tmpfs)
- Fast health check intervals
- Parallel test execution where possible

## Troubleshooting

### Docker Issues
```bash
# Check if containers are running
docker ps | grep pg-goer-integration

# View logs
docker logs pg-goer-integration-test

# Force cleanup
make integration-clean
```

### Connection Issues
- Ensure port 5556 is available
- Check Docker daemon is running
- Verify PostgreSQL health check passes

### Test Failures
- Review test output for specific assertions
- Check database logs for connection/query issues
- Ensure schema loaded correctly with test data

## Development

To add new integration tests:

1. Add test functions to `integration_test.go`
2. Use the `+build integration` tag
3. Follow existing patterns for setup/teardown
4. Update test data in `testdata/` if needed

To modify the test schema:
1. Edit `testdata/01-basic-schema.sql`
2. Update expected values in test assertions
3. Run tests to verify changes work correctly
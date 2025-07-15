# Usage

## Installation
```bash
go install github.com/orchard9/pg-goer@latest
```

## Basic Usage
```bash
pg-goer "postgresql://user:password@localhost/dbname"
```

## Options
```bash
pg-goer [flags] <connection-string>

Flags:
  -o, --output string    Output file (default: README.md)
  -f, --format string    Output format: markdown, json (default: markdown)
  --no-diagram          Skip ER diagram generation
  --no-stats            Skip table statistics
  -h, --help            Show help
```

## Examples

### Generate report for local database
```bash
pg-goer "postgresql://localhost/myapp"
```

### Output to specific file
```bash
pg-goer -o database-docs.md "postgresql://localhost/myapp"
```

### JSON output
```bash
pg-goer -f json "postgresql://localhost/myapp"
```

## Environment Variables
- `PGCONNECT_TIMEOUT`: Connection timeout (default: 10s)
- `PGGOER_MAX_TABLES`: Maximum tables to analyze (default: 1000)
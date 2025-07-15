# PG Go ER

Automatically generate comprehensive PostgreSQL database documentation including ER diagrams, table statistics, and schema definitions.

*This is an AI-written project, built with modern software engineering practices.*

## Features
- Complete schema documentation
- Entity-relationship diagrams
- Table size and activity metrics
- Zero configuration
- Single binary distribution

## Quick Start
```bash
go install github.com/orchard9/pg-goer@latest
pg-goer "postgresql://localhost/mydb"
```

## Output Example
PG Go ER generates a markdown report containing:
- Database overview
- Table definitions with columns and types
- Foreign key relationships
- Table statistics (row counts, last activity)
- Mermaid ER diagram

📄 **[View Sample Output](example-output.md)** - See what the generated documentation looks like

## Installation

### From Source
```bash
git clone https://github.com/orchard9/pg-goer
cd pg-goer
go build -o pg-goer cmd/pg-goer/main.go
```

### Using Go
```bash
go install github.com/orchard9/pg-goer@latest
```

## Requirements
- PostgreSQL 12+
- Go 1.21+ (for building from source)

## Documentation
- [Usage Guide](usage.md)
- [Architecture](code_architecture.md)
- [Contributing](contributing.md)
- [Why PG Go ER?](why.md)

## License
MIT License. See [LICENSE.md](LICENSE.md) for details.

## Support
Report issues at https://github.com/orchard9/pg-goer/issues
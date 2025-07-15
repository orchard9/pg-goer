# Code Architecture

## Overview
PG Go ER follows a clean architecture pattern with clear separation of concerns.

## Core Components

### Database Analyzer
- Connects to PostgreSQL
- Extracts schema metadata
- Gathers table statistics
- Identifies relationships

### ER Generator
- Processes schema data
- Creates entity relationships
- Generates Mermaid diagrams

### Report Builder
- Formats analysis results
- Creates markdown output
- Structures documentation

## Data Flow
```
PostgreSQL → Analyzer → Data Models → Generator → Report
```

## Key Interfaces
```go
type Analyzer interface {
    Connect(dsn string) error
    AnalyzeSchema() (*Schema, error)
    Close() error
}

type Generator interface {
    Generate(schema *Schema) (*Diagram, error)
}

type Reporter interface {
    Build(schema *Schema, diagram *Diagram) (string, error)
}
```

## Design Principles
- Single responsibility
- Dependency injection
- Interface segregation
- Testability first
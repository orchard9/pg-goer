# Semantic Memory

## Project Facts
- PG GoER is a comprehensive PostgreSQL database documentation generator
- The project is written in Go 1.21+ following elegant efficiency principles
- It generates both markdown and JSON documentation formats
- The tool uses Mermaid syntax for ER diagram generation
- It follows clean architecture principles with strong separation of concerns
- The project is licensed under MIT and fully open source
- It is designed as a single binary command-line tool with zero dependencies
- The default output format is markdown with JSON as alternative
- It provides comprehensive PostgreSQL schema analysis including:
  - Tables, columns, constraints, indexes, triggers
  - Views, sequences, extensions, foreign key relationships
  - Row counts, data types, default values
- Connection timeout defaults to 10 seconds with proper connection pooling
- The tool includes comprehensive UAT testing with Docker PostgreSQL
- Example outputs are automatically generated during UAT in both formats
- All GitHub Actions CI/CD pipeline is functional with multi-platform builds
- The project is explicitly AI-written and transparent about it
- It prioritizes simplicity, reliability, and zero configuration
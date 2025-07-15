# Why PG Go ER?

## The Problem
Database documentation is often outdated, incomplete, or non-existent. Developers waste time reverse-engineering schemas and relationships. Existing tools require complex dependencies or manual processes.

## The Solution
PG Go ER automatically generates comprehensive database documentation by analyzing your PostgreSQL database structure in real-time. One command, zero dependencies, instant documentation.

## Key Benefits
- **Single Binary**: No Java, Python, or Node.js required
- **PostgreSQL Native**: Built specifically for PostgreSQL, not a generic tool
- **Always Current**: Documentation reflects actual database state
- **GitHub-Ready**: Markdown + Mermaid diagrams work perfectly in repos
- **CI/CD Friendly**: Automate documentation in your pipeline
- **Privacy First**: Runs locally, no cloud uploads

## Use Cases
- Onboarding new team members
- Database migration planning
- Architecture reviews
- Documentation compliance
- Development reference
- Automated documentation in CI/CD

## Why Not Alternatives?

### SchemaSpy/SchemaCrawler
- Requires Java runtime
- Complex configuration
- Generic multi-database approach
- Heavier resource usage

### dbdocs.io
- Online service (privacy concerns)
- Manual schema uploads
- Not automation-friendly

### pgModeler
- Primarily GUI tool
- Complex for simple documentation
- Overkill for most use cases

### pgAdmin
- GUI-based, not CLI-friendly
- Can't automate in pipelines
- No markdown output

## Our Philosophy
PG Go ER does one thing exceptionally well: generate PostgreSQL documentation. No bloat, no complex configurations, no heavy dependencies. Just a fast, reliable tool that fits seamlessly into modern development workflows.
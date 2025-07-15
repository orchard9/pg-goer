# Contributing

## Code of Conduct
Be respectful. We're here to build great software together.

## How to Contribute

### Reporting Issues
- Check existing issues first
- Provide clear reproduction steps
- Include PostgreSQL version

### Pull Requests
1. Fork the repository
2. Create a feature branch
3. Write tests for new functionality
4. Ensure all tests pass
5. Submit PR with clear description

### Development Setup
```bash
git clone https://github.com/yourusername/pg-goer
cd pg-goer
go mod download
go test ./...
```

### Testing
- Unit tests required for new features
- Integration tests use Docker PostgreSQL
- Run `make test` before submitting

### Commit Messages
- Use present tense
- Keep under 72 characters
- Reference issues when applicable

### Code Review
- All PRs require review
- Address feedback promptly
- Keep discussions focused

## Release Process
Maintainers handle releases following semantic versioning.
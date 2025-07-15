# Coding Guidelines

## Language
Go 1.21+

## Code Style
- Follow standard Go conventions and `gofmt`
- Use meaningful variable names
- Keep functions small and focused
- Prefer composition over inheritance
- Handle errors explicitly

## Project Structure
```
pg-goer/
├── cmd/           # CLI entry points
├── internal/      # Private application code
├── pkg/           # Public packages
├── docs/          # Documentation
└── tests/         # Test files
```

## Error Handling
- Return errors, don't panic
- Wrap errors with context
- Use custom error types for domain errors

## Testing
- Write table-driven tests
- Aim for 80%+ coverage
- Mock external dependencies

## Dependencies
- Minimize external dependencies
- Use standard library when possible
- Pin dependency versions

## Performance
- Profile before optimizing
- Use connection pooling
- Stream large results

## Security
- Never log connection strings
- Validate all inputs
- Use prepared statements
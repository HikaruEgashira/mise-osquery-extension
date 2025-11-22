# Contributing

Thank you for your interest in contributing to this project!

## Development Setup

1. Clone the repository:
```bash
git clone https://github.com/HikaruEgashira/version-manager-packages-osquery-extension.git
cd version-manager-packages-osquery-extension
```

2. Install dependencies:
```bash
go mod download
```

3. Run tests:
```bash
go test -v ./...
```

## Building

Build the extension:
```bash
go build -o version_manager_packages_extension
```

## Testing

Run all tests:
```bash
go test -v ./...
```

Run tests with coverage:
```bash
go test -v -cover ./...
```

## Testing with osquery

1. Build the extension:
```bash
go build -o version_manager_packages_extension
```

2. Run osqueryi with the extension:
```bash
osqueryi --extension ./version_manager_packages_extension
```

3. Query the table:
```sql
SELECT * FROM version_manager_packages;
```

## Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Run `go vet` to catch common mistakes
- Add tests for new features

## Pull Request Process

1. Fork the repository
2. Create a new branch for your feature
3. Make your changes
4. Run tests and ensure they pass
5. Submit a pull request

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

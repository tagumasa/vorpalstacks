# Contributing to Vorpalstacks

Thank you for your interest in contributing! This guide covers the basics.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/vorpalstacks.git`
3. Build: `make build`
4. Run: `make run`

## Development Setup

- Go 1.25+
- Docker (for Lambda functionality)

## Making Changes

1. Create a feature branch from `main`
2. Make your changes
3. Add godoc comments to all exported symbols
4. Run `make lint` to check for issues
5. Run `make test` to verify tests pass
6. Commit with a clear message
7. Open a pull request

## Code Style

- Follow existing conventions in each package
- Use `gofmt` for formatting
- Add godoc comments to all exported symbols (functions, types, constants, variables)

## Testing

### Unit Tests

```bash
make test
```

### SDK Integration Tests

```bash
# Start the server
SIGNATURE_VERIFICATION_ENABLED=false PORT=8080 DATA_PATH=./data TEST_MODE=true ./vorpalstacks &

# Run SDK tests
cd sdk-tests && ./sdk-tests-test -service all -v
```

### CLI Integration Tests

```bash
cd scripts/services && bash test_sqs.sh
cd scripts/services && bash test_dynamodb.sh
```

## Adding a New AWS Service

1. Create service handler in `internal/services/aws/<service>/`
2. Create store in `internal/store/aws/<service>/`
3. Register routing in `internal/server/http/`
4. Add godoc comments to all exported symbols
5. Add SDK tests in `sdk-tests/`
6. Update documentation in `docs/services.md`

## Pull Request Process

- Keep PRs focused on a single concern
- Include tests for new functionality
- Update documentation as needed
- Ensure `make lint` and `make test` pass

## License

By contributing, you agree that your contributions will be licensed under the FSL-1.1-MIT license (for the root project) and Apache 2.0 (for `pkg/`).

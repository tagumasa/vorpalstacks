# AGENTS.md

## Project Overview

**vorpalstacks** is a lightweight edge and on-premise cloud platform providing AWS-compatible services. It is a real implementation (not a mock) that supports 30 AWS services with full SDK compatibility.

- **Language**: Go 1.25+
- **License**: FSL-1.1-MIT (root), Apache 2.0 (`pkg/`)
- **SDK Tests**: 594 tests in `sdk-tests/` using AWS Go SDK v2

## Architecture

- `cmd/` — Entry points
- `internal/services/aws/` — Service implementations (one package per AWS service)
- `internal/store/aws/` — Persistent storage layer (PebbleDB per region)
- `internal/server/` — HTTP routing, request parsing, dispatching
- `internal/core/` — Cross-cutting concerns (resilience, auth, middleware)
- `pkg/` — Reusable libraries (SQL parser, JWT, VTL engine, filter patterns)
- `web_ui/` — Flutter admin console (gRPC-Web)
- `proto/` — Protocol buffer definitions

## Key Patterns

### Service Implementation

Each AWS service has a handler in `internal/services/aws/<service>/` that implements CRUD operations. Services receive parsed requests and return structured responses.

### Request Routing

AWS requests are routed by protocol (REST-XML, REST-JSON, AWS JSON, Query) via `internal/server/http/`. IAM actions are routed through `internal/server/actionregistry/registry.go`.

### Storage

Region-specific data is stored in PebbleDB under `DATA_PATH/<region>/pebble/`. Global services (IAM, Route53, CloudFront, STS) use `DATA_PATH/global/pebble/`.

### Protocol Quirks

- Query protocol list parameters use `Key.member.N` format (e.g., `ClientIDList.member.1=c1`)
- JSON round-trip through `interface{}` loses type information — handle both typed structs and `map[string]interface{}`
- AWS SDK field names sometimes differ from the XML/JSON field names

## Testing

```bash
# Build
go build -o vorpalstacks .

# Start server (development)
SIGNATURE_VERIFICATION_ENABLED=false PORT=8080 DATA_PATH=./data TEST_MODE=true ./vorpalstacks &

# SDK tests
cd sdk-tests && ./sdk-tests-test -service all -v

# CLI integration tests
cd scripts/services && bash test_sqs.sh
cd scripts/services && bash test_dynamodb.sh
cd scripts/services && bash test_s3.sh
```

## Contributing Guidelines

- Follow existing code conventions in each package
- Add godoc comments to all exported symbols
- Run `make test` before submitting changes
- Do not modify `pkg/sqlparser/dependency/` (vendored from vitess)
- Proto files are generated — do not edit directly

## staticcheck False Positives

The following are known false positives and should not be fixed:
- **U1000**: Unused functions (retained for future expansion)
- **SA4006**: Unused assignment values (intentional)
- **SA5011**: Nil deref in tests
- **SA1031**: Nil check before interface range
- `pkg/sqlparser/dependency/hack/` — External dependency, do not modify

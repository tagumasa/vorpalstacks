# AGENTS.md

## Project Overview

**vorpalstacks** is a lightweight edge and on-premise cloud platform providing AWS-compatible services. It is a real implementation (not a mock) that supports 32 AWS services with full SDK compatibility.

- **Language**: Go 1.25+
- **License**: FSL-1.1-MIT (root), Apache 2.0 (`pkg/`)
- **SDK Tests**: 2216 SDK tests + 29 integration + 17 WebSocket (2262 Go total), 631 Python, 2028 TypeScript, 2019 C#

## Architecture

```
cmd/                          — Entry points (server, proto_generator, pebble-loader)
main.go                       — Application entry point
internal/
  services/aws/<service>/     — Service implementations (one package per AWS service)
  store/aws/                  — Persistent storage layer (PebbleDB per region)
  store/api/                  — Storage interface definitions
  store/config/               — Storage configuration
  server/http/                — HTTP routing, request parsing, dispatching
  server/actionregistry/      — IAM action routing registry
  server/authorization/       — IAM policy evaluation
  server/grpcweb/             — Connect-RPC gRPC-Web admin server
  server/dispatcher/          — Request dispatching
  server/listener/            — Server lifecycle management
  server/apps/                — Application wiring
  core/resilience/            — Resilience patterns
  core/auth/                  — Authentication
  core/telemetry/             — Telemetry / observability
  core/storage/               — Storage helpers
  core/hooks/                 — Lifecycle hooks
  core/logs/                  — Logging utilities
  eventbus/                   — Event-driven inter-service communication
  smithy/                     — AWS Smithy model processing
  client/                     — Internal service clients
  pb/                         — Generated protobuf / Connect-RPC code
  utils/                      — Shared utilities
  tools/                      — Development tools
  common/                     — Common types and constants
pkg/
  sqlparser/                  — SQL parser (vendored from vitess)
  vsjwt/                      — JWT library
  vtl/                        — VTL (Velocity Template Language) engine
  filterpattern/              — CloudWatch Logs filter patterns
  cypherparser/               — openCypher query parser
  gremlinparser/              — Gremlin query parser
webconsole/                   — Web admin console (TypeScript/React, npm build, embedded via embed.FS)
proto/                        — Protocol buffer definitions
sdk-tests/                    — AWS SDK v2 integration tests
scripts/services/             — CLI integration test scripts per service
docs/                         — Documentation
```

## Implemented Services (32)

ACM, API Gateway, AppSync, Athena, CloudFront, CloudTrail, CloudWatch Logs, CloudWatch Metrics, Cognito IDP, Cognito Identity, DynamoDB, EC2, EventBridge, IAM, Kinesis, KMS, Lambda, Neptune, Neptune Data, Neptune Graph, Route53, S3, Scheduler, Secrets Manager, SESv2, SNS, SQS, SSM, STS, Step Functions, Timestream (Query + Write), WAFv2

## Key Patterns

### Service Implementation

Each AWS service has a handler in `internal/services/aws/<service>/` that implements CRUD operations. Services receive parsed requests and return structured responses.

### Request Routing

AWS requests are routed by protocol (REST-XML, REST-JSON, AWS JSON, Query) via `internal/server/http/`. IAM actions are routed through `internal/server/actionregistry/registry.go`.

### Authorization

IAM policy evaluation is handled in `internal/server/authorization/`. Enable with `AUTHORIZATION_ENABLED=true`. Actions are registered per service in the action registry.

### Storage

Region-specific data is stored in PebbleDB under `DATA_PATH/<region>/pebble/`. Global services (IAM, Route53, CloudFront, STS) use `DATA_PATH/global/pebble/`.

```
DATA_PATH/
├── us-east-1/
│   ├── pebble/       # Region-isolated PebbleDB
│   ├── code/         # Lambda function code
│   └── logs/         # CloudWatch Logs chunks
├── global/
│   └── pebble/       # Global services (IAM, Route53, CloudFront, STS)
└── uploads/          # S3 multipart uploads (temporary)
```

### Admin API

Connect-RPC gRPC-Web admin interface runs on a separate port (`GRPC_WEB_PORT`, default 50090). The web console at `/webconsole/` is served from this port. Proto definitions are in `proto/`, generated Go code in `internal/pb/`.

### Protocol Quirks

- Query protocol list parameters use `Key.member.N` format (e.g., `ClientIDList.member.1=c1`)
- JSON round-trip through `interface{}` loses type information — handle both typed structs and `map[string]interface{}`
- AWS SDK field names sometimes differ from the XML/JSON field names

## Build & Run

```bash
# Build (compiles webconsole frontend + Go binary)
make build

# Run in development mode
SIGNATURE_VERIFICATION_ENABLED=false DATA_PATH=./data ./vorpalstacks_server

# Run with Lambda support (requires Docker)
SIGNATURE_VERIFICATION_ENABLED=false DATA_PATH=./data DOCKER_HOST=unix:///var/run/docker.sock ./vorpalstacks_server

# Format and tidy
make fmt
make tidy
```

- AWS API endpoints: port 50080 (`PORT`)
- Admin console: `http://localhost:50090/webconsole/` (`GRPC_WEB_PORT`)

## Testing

```bash
# Unit tests
make test

# SDK tests (requires running server with audit enabled)
# Start server: VS_AUDIT_ENABLED=true SIGNATURE_VERIFICATION_ENABLED=false DATA_PATH=./data ./vorpalstacks_server
cd sdk-tests && go build -o sdk-tests-all .
./sdk-tests-all -service all -endpoint http://127.0.0.1:50080 -v

# CLI integration tests (all services)
make test-cli

# Individual CLI tests
cd scripts/services && bash test_sqs.sh
cd scripts/services && bash test_dynamodb.sh
cd scripts/services && bash test_s3.sh
cd scripts/services && bash test_iam.sh
cd scripts/services && bash test_lambda.sh
cd scripts/services && bash test_stepfunctions.sh
# ... see scripts/services/ for all available test scripts
```

## Key Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `50080` | HTTP server port |
| `DATA_PATH` | `./data` | Persistent data storage path |
| `AWS_REGION` | `us-east-1` | Default region |
| `AWS_ACCOUNT_ID` | `000000000000` | AWS account ID |
| `SIGNATURE_VERIFICATION_ENABLED` | `true` | AWS Signature V4 verification |
| `GRPC_WEB_PORT` | `50090` | gRPC-Web admin server port |
| `TLS_ENABLED` | `false` | Enable TLS |
| `AUTHORIZATION_ENABLED` | `false` | Enable IAM policy evaluation |
| `DOCKER_HOST` | `unix:///var/run/docker.sock` | Docker socket for Lambda |

See `docs/configuration.md` for the complete list.

## Contributing Guidelines

- Follow existing code conventions in each package
- Add godoc comments to all exported symbols
- Run `make test` before submitting changes
- Do not modify `pkg/sqlparser/dependency/` (vendored from vitess)
- Proto files are generated — do not edit generated code directly
- Use `make proto` to regenerate protobuf Go code from `proto/` definitions
- Use `make proto-generate` to generate `.proto` files from Smithy JSON models

## Linting

Configured via `.golangci.yml` with: errcheck, gosimple, govet, ineffassign, staticcheck, unused.

## staticcheck False Positives

The following are known false positives and should not be fixed:
- **U1000**: Unused functions (retained for future expansion)
- **SA4006**: Unused assignment values (intentional)
- **SA5011**: Nil deref in tests
- **SA1031**: Nil check before interface range
- `pkg/sqlparser/dependency/hack/` — External dependency, do not modify

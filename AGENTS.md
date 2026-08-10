# AGENTS.md

## Project Overview

**vorpalstacks** is a lightweight edge and on-premise cloud platform providing AWS-compatible services. It is a real implementation (not a mock) that supports 34 AWS service APIs with full SDK compatibility.

- **Language**: Go 1.25+
- **License**: FSL-1.1-MIT (root), Apache 2.0 (`pkg/`)
- **SDK Tests**: 2,702 SDK + 47 integration + 17 WebSocket = 2,766 Go tests. Python, TypeScript, and C# test suites exist in separate repositories.

## Architecture

```
main.go                       — Application entry point
cmd/
  vstacks/                    — Admin CLI (server control, IAM, config, service modes, backup)
  proto_generator/            — Generates .proto files from Smithy JSON models
  pebble-loader/              — Loads Smithy JSON models into Pebble storage
internal/
  services/aws/<service>/     — Service implementations (one package per AWS service)
  store/aws/                  — Persistent storage layer (PebbleDB per region)
  store/api/                  — Storage interface definitions
  store/config/               — Configuration storage (PebbleDB-backed app_config)
  server/http/                — HTTP routing, request parsing, dispatching
    chain/                    — Middleware chain with gateway support
    classifier/               — Request protocol/target classification
    router/                   — Service index for route lookup
  server/actionregistry/      — IAM action routing registry
  server/authorization/       — IAM policy evaluation
  server/grpcweb/             — Connect-RPC gRPC-Web admin server
  server/dispatcher/          — Request dispatching, audit, fallback
  server/listener/            — Server lifecycle management
  server/fqdnrouter/          — Host-based FQDN routing for per-resource hostnames
  server/portalloc/           — Per-resource port allocator (Individual mode)
  server/apps/                — Application wiring, conditional service registration
  config/                     — Bootstrap (env → struct) and runtime config (PebbleDB)
  serviceconfig/              — Static service registry (names, env vars, port keys, FQDN support)
  core/resilience/            — Resilience patterns (circuit breaker, retry, bulkhead, cache)
  core/auth/                  — Authentication
  core/telemetry/             — Telemetry / observability (OpenTelemetry)
  core/storage/               — Storage helpers (region manager, blob, multipart)
  core/hooks/                 — Lifecycle hooks
  core/logs/                  — CloudWatch Logs storage primitives
  client/mobyclient/          — Docker/Moby client for Lambda container lifecycle
  eventbus/                   — Event-driven inter-service communication (with transactional outbox)
  smithy/                     — AWS Smithy model processing
  common/                     — Shared types and constants (audit, pagination, protocol, serviceports, tags, iotutil, ...)
  pb/                         — Generated protobuf / Connect-RPC code (aws/, storage/)
  utils/                      — Shared utilities (aws, crypto, graphql, naming, ptrutil, timeutils, ...)
  tools/                      — Development tools (converter, proto_generator, vstackscli)
pkg/
  sqlparser/                  — SQL parser (vendored from vitess)
  vsjwt/                      — JWT library
  vtl/                        — VTL (Velocity Template Language) engine
  filterpattern/              — CloudWatch Logs filter patterns
  metricmath/                 — CloudWatch metric math expression engine
  cypherparser/               — openCypher query parser, planner, and executor
  gremlinparser/              — Gremlin query parser and executor
webconsole/                   — Web admin console (React 19 + Vite + TypeScript, npm build, embedded via embed.FS)
proto/                        — Protocol buffer definitions
sdk-tests/                    — AWS SDK v2 integration tests (independent Go module)
docs/                         — Documentation
```

## Implemented Services (34)

ACM, API Gateway, AppSync, Athena, CloudFront, CloudTrail, CloudWatch Logs, CloudWatch Metrics, Cognito IDP, Cognito Identity, DynamoDB, EventBridge, IAM, IoT, Kinesis, KMS, Lambda, Neptune, Neptune Data, Neptune Graph, RDS Data (vMySQL), Route53, S3, Scheduler, Secrets Manager, SESv2, SNS, SQS, SSM, STS, Step Functions, Timestream Query, Timestream Write, WAFv2

**Notes**:
- EC2 exists as a data-plane-only package (VPC/Subnet/Security Group helpers for other services) and is not counted as a standalone service.
- Neptune, Neptune Data, Neptune Graph, and RDS Data (vMySQL) share the composite `internal/services/aws/rds/` package, each with its own subdirectory and `service.go`.
- IoT, CloudTrail, and RDS MySQL default to **disabled** (`IOT_ENABLED=false`, `CLOUDTRAIL_ENABLED=false`, `RDS_MYSQL_ENABLED=false`). All other services default to enabled.

## Key Patterns

### Service Implementation

Each AWS service has a handler in `internal/services/aws/<service>/` that implements CRUD operations. Services receive parsed requests and return structured responses.

### Request Routing

AWS requests are routed by protocol (REST-XML, REST-JSON, AWS JSON, Query) via `internal/server/http/`. The classifier detects the protocol and extracts service/operation. IAM actions are routed through `internal/server/actionregistry/registry.go`.

### Service Modes

Every service has a `<NAME>_ENABLED` environment variable (e.g. `S3_ENABLED`, `IOT_ENABLED`). Set `ALL_SERVICES_ENABLED=true` to force-enable all services regardless of individual flags. Conditional wiring is in `internal/server/apps/optional.go`. The static service registry lives in `internal/serviceconfig/`.

→ See `docs/configuration.md` § Service Enablement for the full list.

### Routing Modes

Two port-routing strategies for services that expose per-resource hostnames (S3 Website, API Gateway, Cognito, CloudFront, Lambda URL, AppSync, Neptune):

- **FQDN mode** (default): all traffic routes through `:50080` via `Host` header matching. No extra listener.
- **Individual mode**: each resource gets its own port from the dynamic range (50200–50400).

Implemented in `internal/server/fqdnrouter/` and `internal/server/portalloc/`. Switchable at runtime via `vstacks config set ports.<service>.mode individual|fqdn`.

→ See `docs/architecture.md` § Port Architecture for details.

### Config Store

Runtime configuration is stored in PebbleDB (`app_config` bucket). Priority: **Pebble (persistent) > Environment variable > Default constant**. At startup, `internal/config/bootstrap.go` reads env vars into `BootstrapConfig`, then `store/config` seeds defaults and applies env overrides to Pebble. After startup, only Pebble is read.

→ See `docs/configuration.md` § Admin Config System and `docs/architecture.md` § Config Store Architecture.

### Authorization

IAM policy evaluation is handled in `internal/server/authorization/`. Enable with `AUTHORIZATION_ENABLED=true`. Actions are registered per service in the action registry. Additional env vars: `AUTHORIZATION_FAILURE_MODE` (strict/permissive), `AUTHORIZATION_CACHE_TTL_SECONDS`, `AUTHORIZATION_DEFAULT_ACCESS_KEY_ID`.

→ See `docs/configuration.md` § IAM Authorization.

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
SIGNATURE_VERIFICATION_ENABLED=false DATA_PATH=./data ./vorpalstacks

# Or use make run (wraps go run with signature verification disabled)
make run

# Run with Lambda support (requires Docker)
SIGNATURE_VERIFICATION_ENABLED=false DATA_PATH=./data DOCKER_HOST=unix:///var/run/docker.sock ./vorpalstacks

# Run with all services (including IoT, CloudTrail, RDS MySQL)
ALL_SERVICES_ENABLED=true SIGNATURE_VERIFICATION_ENABLED=false DATA_PATH=./data ./vorpalstacks

# Format and tidy
make fmt
make tidy
```

- AWS API endpoints: port 50080 (`PORT`)
- Admin console: `http://localhost:50090/webconsole/` (`GRPC_WEB_PORT`)

> **WARNING (Linux) — gopls OOM**: gopls (the Go language server) running concurrently with `go build` / `go vet` / `go test` can exhaust memory and crash the machine. **Always kill gopls before running any go command**: `pkill gopls || true`. Go telemetry should also be disabled: run `go telemetry off` once, or set `GOTELEMETRY=off` in the environment.

### vstacks CLI

The `vstacks` binary (`cmd/vstacks/`) provides server control, IAM management, configuration, service mode control, and backup operations. It communicates via gRPC-Web (server control) or reads PebbleDB directly (IAM, config, backup).

```
vstacks server status|stop
vstacks iam create-user -user <name> | list-users | ...
vstacks config get <key> | set <key> <value> | list | schema
vstacks service get <name> | enable <name> | disable <name> | set-mode <name> -mode <MODE>
vstacks backup create | restore <file> | list
```

→ See `docs/configuration.md` § vstacks CLI for the full command reference.

## Testing

```bash
# Unit tests
make test

# SDK tests (requires running server)
# Start server:
#   ALL_SERVICES_ENABLED=true VS_AUDIT_ENABLED=true \
#     SIGNATURE_VERIFICATION_ENABLED=false DATA_PATH=./data ./vorpalstacks
cd sdk-tests && go build -o sdk-tests-all .
./sdk-tests-all -service all -endpoint http://127.0.0.1:50080 -v

# Run specific service or test type
./sdk-tests-all -service sqs -endpoint http://127.0.0.1:50080 -v
./sdk-tests-all -service all -type sdk -v           # SDK tests only
./sdk-tests-all -service all -type integration -v    # Cross-service integration tests
./sdk-tests-all -service all -type ws -v             # WebSocket tests
./sdk-tests-all -service all -parallel 4 -v          # Parallel execution
```

`sdk-tests/` is an **independent Go module** (`vorpalstacks-sdk-tests`) with its own `go.mod`. Test services are registered declaratively in `sdk-tests/testutil/runner.go`.

## Key Environment Variables

Core settings most commonly used:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `50080` | HTTP server port |
| `GRPC_WEB_PORT` | `50090` | gRPC-Web admin server port |
| `DATA_PATH` | `./data` | Persistent data storage path |
| `AWS_REGION` | `us-east-1` | Default region |
| `AWS_ACCOUNT_ID` | `000000000000` | AWS account ID |
| `SIGNATURE_VERIFICATION_ENABLED` | `true` | AWS Signature V4 verification |
| `DOCKER_HOST` | `unix:///var/run/docker.sock` | Docker socket for Lambda |
| `ALL_SERVICES_ENABLED` | `false` | Force-enable all services |
| `TEST_MODE` | `false` | Allow unauthenticated access (testing only) |

Additional variable groups (see `docs/configuration.md` for the complete list):
- **Service enablement**: `<NAME>_ENABLED` per service (e.g. `IOT_ENABLED`, `S3_ENABLED`)
- **TLS**: `TLS_ENABLED`, `TLS_PORT`, `TLS_CERT_PATH`, `TLS_KEY_PATH`, `TLS_HOSTNAME`
- **Authorization**: `AUTHORIZATION_ENABLED`, `AUTHORIZATION_FAILURE_MODE`, `AUTHORIZATION_CACHE_TTL_SECONDS`
- **Networking**: `BIND_MODE`, `BIND_INTERFACE`, `TRUST_FORWARDED_HEADERS`
- **Telemetry**: `OTEL_TRACES_EXPORTER`, `OTEL_EXPORTER_OTLP_ENDPOINT`
- **Routing**: `ROUTE53_DNS_ENABLED`, `USE_CHAIN_GATEWAY`

## Git Rules

- **NEVER** modify `user.name`, `user.email`, or any identity-related git config (`git config user.email`, etc.) without explicit user instruction. Changing the author identity corrupts commit attribution and is strictly forbidden.
- When a git command fails, read the error message and fix the root cause. Do **not** run `git config` as a workaround.
- Do not force-push, rebase, or rewrite history unless the user explicitly requests it.

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

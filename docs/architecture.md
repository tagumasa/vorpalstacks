# Architecture

This document describes the architecture of Vorpalstacks.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                  HTTP Server (Chi)                              │
│                 :50080 (configurable)                           │
└────────────────────────────┬────────────────────────────────────┘
                             │
               ┌──────────────┴──────────────┐
               │                             │
               ▼                             ▼
┌───────────────────────────┐   ┌───────────────────────────┐
│    AWS Service Handler    │   │   API Gateway Runtime     │
│   (Standard endpoints)    │   │   (REST API execution)    │
└─────────────┬─────────────┘   └─────────────┬─────────────┘
              │                             │
              └──────────────┬──────────────┘
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Dispatcher                                  │
│   Action Registry → Signature Verification → IAM Authorization  │
│              Routes operations to service handlers              │
└────────────────────────────┬────────────────────────────────────┘
                             │
         ┌───────────────────┼───────────────────┐
         ▼                   ▼                   ▼
┌───────────────┐   ┌───────────────┐   ┌───────────────┐
│  Service A    │   │  Service B    │   │  Service C    │
│ (handler +    │   │ (handler +    │   │ (handler +    │
│  store)       │   │  store)       │   │  store)       │
└───────┬───────┘   └───────┬───────┘   └───────┬───────┘
        │                   │                   │
        └───────────────────┼───────────────────┘
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                  Store Layer (Pebble)                            │
│              Region-Isolated Key-Value Store                     │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│              gRPC-Web Admin Server (Connect-RPC)                │
│                 :50090 (configurable)                           │
│                                                                 │
│  Admin handlers for all 32 services (admin_handler.go)          │
│  Runtime config, service status, port mapping                   │
└─────────────────────────────────────────────────────────────────┘
```

## Port Architecture

All port constants are centralised in `internal/common/serviceports/ports.go`.

### Fixed Listeners

| Port | Purpose |
|------|---------|
| `50080` | HTTP server — all AWS API endpoints |
| `50090` | gRPC-Web admin API + web console |
| `50443` | HTTPS server (TLS) |
| `50088` | Route53 DNS server |
| `50089` | Route53 health check target |

### Service Endpoints (FQDN / Individual)

Service endpoints (S3 Website, API Gateway, Cognito, CloudFront, Lambda URL, AppSync, Neptune) default to **FQDN mode**: all traffic is routed through `:50080` via `Host` header matching. No additional listener is bound.

In **Individual mode**, each resource gets its own port from the dynamic range (`50200–50400`). Neptune uses `50107` for its first cluster when in Individual mode.

| Reserved Slot | Service |
|--------------|---------|
| `50101` | S3 Website |
| `50102` | API Gateway |
| `50103` | Cognito Hosted UI |
| `50104` | CloudFront |
| `50105` | Lambda Function URL |
| `50106` | AppSync Events |
| `50107` | Neptune (first cluster) |

### Port Mode Configuration

```bash
# View current mode
vstacks config get ports.apigateway.mode

# Switch to individual ports
vstacks config set ports.apigateway.mode individual

# Revert to FQDN mode
vstacks config set ports.apigateway.mode fqdn
```

## Two-Layer Architecture

### Service Layer (`internal/services/aws/`)

Handles business logic:

- Parameter extraction and validation
- Operation execution
- Response formatting
- Cross-service communication

Each service consists of:
- `service.go` — Service struct and handler registration
- `*_operations.go` — Operation implementations
- `admin_handler.go` — gRPC-Web Connect-RPC admin handlers
- `errors.go` — Service-specific error types
- `helpers.go` — Utility functions

### Store Layer (`internal/store/aws/`)

Handles data persistence:

- CRUD operations
- ARN generation
- Tag management
- Data serialisation (JSON)

### Common Layer (`internal/services/aws/common/`)

Shared utilities:
- `request/` — RequestContext, parameter extraction, StoreProvider
- `tags/` — Generic tag operations framework
- `iam/` — Role validation (RolePolicyProvider interface)
- `endpoint/` — URL generation with port mapping

## Request Flow

### HTTP API Path

1. **HTTP Request** arrives at Chi router (`:50080`)
2. **Classifier** (`internal/server/http/classifier/`) detects protocol and extracts service/operation
3. **Signature Verification** (if enabled) validates AWS SigV4
4. **IAM Authorization** (if enabled) evaluates policies
5. **Service Handler** processes the request
6. **Store Layer** persists/retrieves data via StorageManager
7. **Response** serialised and returned

### gRPC-Web Admin Path

1. **Connect-RPC Request** arrives at gRPC-Web server (`:50090`)
2. **Admin Handler** processes the request
3. **StorageManager** provides region-aware storage access
4. **Response** returned as Connect-RPC message

## Action Registry

`internal/server/actionregistry/registry.go` maps operation names to services for routing:

```
Operation Name → Service Name → Handler
CreateQueue    → sqs           → SQSService.CreateQueue
ListTopics     → sns           → SNSService.ListTopics
```

Required for Query protocol services (SQS, IAM). JSON protocol services use content-type-based routing.

## IAM Authorization

`internal/server/dispatcher/authorization.go` provides policy-based access control:

```
Request → Extract Access Key → Gather Policies → Evaluate → Allow/Deny
                                (inline + attached + group)
```

- `ResourceExtractor` maps operations to ARN patterns (14+ services)
- `ActionMapper` maps operations to IAM action names (23+ services)
- Policy evaluation cache with configurable TTL and size

## Protocol Support

| Protocol | Content-Type | Services |
|----------|-------------|----------|
| AWS JSON 1.1 | `application/x-amz-json-1.1` | IAM, API Gateway, SESv2, Lambda, Cognito, KMS, Athena, WAFv2, EventBridge, CloudTrail, Step Functions, CloudWatch Logs, Route53, SecretsManager, ACM, CloudWatch, Scheduler, SSM, STS, Timestream, Kinesis |
| AWS JSON 1.0 | `application/x-amz-json-1.0` | DynamoDB |
| REST-XML | XML over HTTP | S3, CloudFront |
| AWS Query | `application/x-www-form-urlencoded` | SQS |
| Connect-RPC | `application/connect+proto` | All 32 services (admin API on :50090) |

## Service Integration Patterns

### Event-Driven

```
EventBridge ──event──▶ Lambda
      │
      ├── SNS ──notification──▶ Lambda
      │        │
      │        └── SQS ◀──message── Lambda
      │
      ├── SQS (direct)
      ├── Step Functions
      └── CloudWatch Logs
```

### Scheduled

```
Scheduler ──schedule──▶ Lambda
      │
      ├── SQS
      └── SNS
```

### API Gateway Runtime

```
HTTP Request ──▶ API Gateway ──▶ Integration
                                       │
                     ┌─────────────────┼─────────────────┐
                     ▼                 ▼                 ▼
                  Lambda             SQS              SNS
```

### Cross-Cutting

```
KMS ──encryption──▶ SSM, DynamoDB, S3
IAM Authorization ──policy check──▶ All services
CloudTrail Audit ──logging──▶ All API operations
```

## Storage Architecture

### Region-Isolated Storage

```
DATA_PATH/
├── us-east-1/     → PebbleDB, Lambda code, CloudWatch Logs
├── us-west-2/     → PebbleDB, Lambda code, CloudWatch Logs
├── global/        → PebbleDB (IAM, Route53, CloudFront, STS)
└── graph/         → GraphDB (NeptuneGraph)
```

### StorageManager

`server.StorageManager()` provides region-aware storage:

```go
store := server.StorageManager().GetOrCreateStore(region)
```

Global services use `server.Storage()` which returns the global PebbleDB.

### Hybrid Blob Store (S3)

```
Small Objects (< threshold) → Pebble KV Store
Large Objects (>= threshold) → File System (dataDir/blobs/)
```

- Thread-safe with mutex protection
- Supports multipart uploads
- Metadata stored separately in Pebble

## Security Components

### JWT Token Management (`pkg/vsjwt`)

Cognito-compatible JWT token handling:
- Access tokens (authorization)
- ID tokens (identity information)
- Refresh tokens (opaque)
- RS256 signature verification

### IAM Role Validation (`internal/services/aws/common/iam`)

High-performance IAM role validation with caching:
- Role ARN parsing and validation
- Trust policy evaluation (service principals, conditions)
- Ristretto-based policy cache (16MB, 10min TTL)

Supported services: Lambda, Step Functions, EventBridge, Scheduler, CloudTrail

### KMS Backend (`internal/services/aws/kms/hsm`)

Key management with HSM-backed cryptographic operations:
- **HSM Interface**: Pluggable backend (`hsm.Backend` interface)
- **Persistent Backend**: AES-256-GCM encrypted keys persisted to disk
- **Memory Backend**: In-memory keys (testing only)
- Supported operations: Encrypt, Decrypt, Sign, Verify, GenerateMAC, VerifyMAC, GenerateDataKey, asymmetric key pairs (RSA, ECC)

## Config Store Architecture

### Single Source of Truth: Pebble

All runtime configuration is stored in PebbleDB (`app_config` bucket). The startup flow is:

1. `NewStore()` creates the store with defaults from `serviceports` constants
2. `Initialise()` → `seedDefaults()` writes missing keys to Pebble
3. `applyEnvOverrides()` overwrites Pebble values with any set ENV vars
4. From this point, **only Pebble is read** — no further ENV lookups

### Access Pattern

```go
// Always use the global singleton via appconfig
store := appconfig.GetStore()
port := store.GetInt("server.port")
```

**Never** call `storeconfig.NewStore()` directly — it creates a separate instance without `Initialise()`.

### Priority Order

**Pebble (persistent) > Environment variable > Default constant**

ENV overrides are applied once at startup. Changes made via the web console or CLI persist in Pebble but are overwritten on next startup if the ENV var is set.

## Scalability Considerations

- Single-process architecture (suitable for edge/on-premise)
- Pebble provides efficient key-value storage with WAL
- Docker provides Lambda runtime isolation
- No external database dependencies
- Region-isolated storage for multi-region support

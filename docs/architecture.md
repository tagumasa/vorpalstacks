# Architecture

This document describes the architecture of Vorpalstacks.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                  HTTP Server (Chi)                              │
│                 :8080 (configurable)                            │
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
│                 :9090 (configurable)                            │
│                                                                 │
│  Admin handlers for all 32 services (admin_handler.go)          │
│  Runtime config, service status, port mapping                   │
└─────────────────────────────────────────────────────────────────┘
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

1. **HTTP Request** arrives at Chi router (`:8080`)
2. **Classifier** (`internal/server/http/classifier/`) detects protocol and extracts service/operation
3. **Signature Verification** (if enabled) validates AWS SigV4
4. **IAM Authorization** (if enabled) evaluates policies
5. **Service Handler** processes the request
6. **Store Layer** persists/retrieves data via StorageManager
7. **Response** serialised and returned

### gRPC-Web Admin Path

1. **Connect-RPC Request** arrives at gRPC-Web server (`:9090`)
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
| Connect-RPC | `application/connect+proto` | All 32 services (admin API on :9090) |

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

## Scalability Considerations

- Single-process architecture (suitable for edge/on-premise)
- Pebble provides efficient key-value storage with WAL
- Docker provides Lambda runtime isolation
- No external database dependencies
- Region-isolated storage for multi-region support

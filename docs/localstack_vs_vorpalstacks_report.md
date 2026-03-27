# Technical Comparison Report: LocalStack vs Vorpalstacks

## Executive Summary

This report provides a technical comparison between LocalStack v4.14.0 (open source) and Vorpalstacks based on source code analysis. The two projects target different primary use cases: LocalStack is designed as a general-purpose AWS cloud emulator for development and testing, while Vorpalstacks is designed for edge and on-premises deployment as a single binary.

---

## 1. Architecture Overview

### LocalStack (open source)

| Aspect | Details |
|--------|---------|
| **Language** | Python (>=3.10, target 3.13) |
| **Version** | 4.14.0 |
| **Deployment** | Docker container (required) |
| **Core Dependencies** | Python 3.13, moto-ext 5.1.25, DynamoDB Local v2.x (Java), kinesis-mock 0.5.2 |
| **Service Implementation** | Hybrid: Custom providers + Moto fallback + External binaries |
| **License** | Apache License 2.0 |

**Implementation tiers:**
```
LocalStack (30 service APIs)
├── Pure Custom (14 services)
│   S3, SQS, SNS, Lambda, KMS, Step Functions,
│   CloudFormation, CloudWatch v2, EventBridge v2,
│   DynamoDB Streams, Elasticsearch, OpenSearch, Firehose
├── Custom + Moto Fallback (23 services)
│   IAM, STS, ACM, API Gateway, CloudWatch v1,
│   CloudWatch Logs, EC2, Route53, Route53 Resolver,
│   Secrets Manager, SSM, SES, S3 Control, Scheduler,
│   Config, Redshift, SWF, Resource Groups, Transcribe,
│   Support, Resource Groups Tagging API, EventBridge v1
├── Custom + External Binary (3 services)
│   DynamoDB (DynamoDB Local - Java),
│   Kinesis (kinesis-mock),
│   OpenSearch (native binary)
└── External Runtime Dependencies
    Lambda (lambda-runtime-init binary),
    Step Functions (JPype/JVM for JSONata),
    Transcribe (Vosk + FFmpeg)
```

### Vorpalstacks

| Aspect | Details |
|--------|---------|
| **Language** | Go 1.25.8 |
| **Deployment** | Single binary (~138 MB) |
| **Core Dependencies** | Pebble v2.1.4 (CockroachDB KV store), Moby client v0.3.0 (Docker) |
| **Service Implementation** | Custom implementation (all services) |
| **License** | FSL-1.1-MIT (root), Apache 2.0 (pkg/) |

**Architecture:**
```
Vorpalstacks (29 service APIs, 25 AWS services)
├── All Custom Implementation (1,081 operations)
│   ├── IAM (156), API Gateway (80), S3 (71),
│   ├── SESv2 incl. SESv1 (58), Cognito User Pools (51),
│   ├── Lambda (50), KMS (45), Kinesis (39),
│   ├── EventBridge (39), WAF v1 (40), WAF v2 (36),
│   ├── CloudFront (32), DynamoDB (57), Athena (37),
│   ├── CloudWatch Logs (29), Step Functions (28),
│   ├── CloudTrail (24), Secrets Manager (21),
│   ├── CloudWatch Metrics (17), Route53 (18),
│   ├── Timestream Write (19), SNS (27), SQS (20),
│   ├── ACM (16), Cognito Identity (13),
│   ├── Timestream Query (13), SSM (12), Scheduler (12),
│   └── STS (7)
├── Pebble Storage (persistent KV store)
│   ├── Region-specific (us-east-1, us-west-2, etc.)
│   └── Global (IAM, Route53, CloudFront, STS)
├── gRPC Internal Protocol
├── Docker Runtime (Lambda execution)
└── Service Integration
    EventBridge → Lambda, SNS → Lambda, SQS → Lambda, etc.
```

**Service API grouping note:**
- SESv2 also handles SESv1 (Query protocol)
- CloudWatch Metrics + CloudWatch Logs are separate service APIs
- WAF (v1) + WAF v2 are separate service APIs
- Cognito Identity + Cognito User Pools are separate service APIs
- Timestream Write + Timestream Query are separate service APIs

---

## 2. Service Implementation Analysis

### 2.1 SQS (Simple Queue Service)

#### LocalStack
- **Implementation**: Pure custom (no Moto fallback)
- **Storage**: In-memory Python dictionaries
- **Evidence**: `localstack-core/localstack/services/sqs/provider.py`

#### Vorpalstacks
- **Implementation**: Custom with Pebble storage
- **Storage**: Pebble KV store with JSON serialisation via `BaseStore`
- **Evidence**: `internal/store/aws/sqs/store.go`
- **Operations**: 20 registered Query handlers

**Observation**: Both implement SQS from scratch. LocalStack uses in-memory storage; Vorpalstacks uses persistent Pebble storage.

---

### 2.2 SNS (Simple Notification Service)

#### LocalStack
- **Implementation**: Pure custom (no Moto fallback)
- **Storage**: In-memory Python dictionaries
- **Evidence**: `localstack-core/localstack/services/sns/provider.py`

#### Vorpalstacks
- **Implementation**: Custom with Pebble storage
- **Storage**: Pebble KV store with JSON serialisation
- **Operations**: 27 registered Query handlers

**Observation**: Both implement SNS from scratch. LocalStack uses in-memory storage; Vorpalstacks uses persistent Pebble storage.

---

### 2.3 Lambda

#### LocalStack
- **Implementation**: Pure custom (no Moto fallback)
- **Execution**: Docker containers via lambda-runtime-init binary
- **Evidence**: `localstack-core/localstack/services/lambda_/provider.py`

#### Vorpalstacks
- **Implementation**: Custom with Moby (Docker) client
- **Execution**: Docker containers via Moby client API
- **Operations**: 50 registered gRPC handlers

**Observation**: Both execute Lambda functions in real Docker containers.

---

### 2.4 IAM (Identity and Access Management)

#### LocalStack
- **Implementation**: Custom provider with Moto fallback (`MotoFallbackDispatcher`)
- **Evidence**: `localstack-core/localstack/services/iam/provider.py`:
```python
from moto.iam.models import (
    IAMBackend,
    Role as MotoRole,
    User as MotoUser,
)
from moto.iam import iam_backends
from localstack.services.moto import call_moto

def create_user(self, context, req):
    result = call_moto(context)
```
- **Storage**: Moto's in-memory backend (with selective `call_moto()` calls)

#### Vorpalstacks
- **Implementation**: Custom with Pebble storage
- **Storage**: Persistent Pebble storage
- **Operations**: 156 registered Query handlers

**Observation**: LocalStack IAM delegates to the Moto mock library for many operations; Vorpalstacks implements IAM operations directly with persistent storage.

---

### 2.5 DynamoDB

#### LocalStack
- **Implementation**: Custom provider with DynamoDB Local (Java) as HTTP fallback
- **Evidence**: `localstack-core/localstack/services/dynamodb/server.py`
- **Storage**: Java process with file-based persistence
- **Fallback**: `HttpFallbackDispatcher` forwards unhandled requests to DynamoDB Local

#### Vorpalstacks
- **Implementation**: Custom with Pebble storage
- **Storage**: Pebble KV store
- **Operations**: 57 registered gRPC handlers
- **Features**: Table CRUD, Item operations, Query, Scan, GSI/LSI, Streams, PartiQL

**Observation**: LocalStack uses AWS's official DynamoDB Local (Java binary) as a fallback; Vorpalstacks implements DynamoDB-compatible operations in Go from scratch.

---

### 2.6 EventBridge

#### LocalStack
- **Implementation**: v1/legacy uses Moto fallback; v2 (default) is pure custom
- **Evidence**: `localstack-core/localstack/services/events/`

#### Vorpalstacks
- **Implementation**: Custom with Lambda invocation integration
- **Evidence** (`internal/services/aws/eventbridge/event_operations.go`):
```go
case "lambda":
    statusCode, result, err := s.lambdaInvoker.InvokeForGateway(ctx, functionName, payload)
```
- **Operations**: 39 registered Query handlers

**Observation**: Both implement EventBridge with actual Lambda invocation. LocalStack's default (v2) is pure custom; the legacy v1 falls back to Moto.

---

## 3. Service Coverage Comparison

### 3.1 Backend Technology

| Service | LocalStack | Vorpalstacks |
|---------|------------|--------------|
| **SQS** | Custom (in-memory) | Pebble (persistent) |
| **SNS** | Custom (in-memory) | Pebble (persistent) |
| **Lambda** | Custom + Docker (real) | Custom + Docker (real) |
| **IAM** | Custom + Moto fallback | Pebble (persistent) |
| **STS** | Custom + Moto fallback | Pebble (persistent) |
| **DynamoDB** | Custom + DynamoDB Local (Java) | Pebble (persistent) |
| **S3** | Custom (pure) | Pebble + file |
| **KMS** | Custom (pure) | Pebble (persistent) |
| **EventBridge** | v1: Moto, v2: Custom (pure) | Pebble (persistent) |
| **Step Functions** | Custom (pure) + JPype/JVM | Pebble (persistent) |
| **CloudFormation** | Custom (pure) | Not implemented |
| **Kinesis** | Custom + kinesis-mock | Pebble (persistent) |
| **ACM** | Custom + Moto fallback | Pebble (persistent) |
| **API Gateway** | Custom + Moto fallback | Pebble (persistent) |
| **Route53** | Custom + Moto fallback | Pebble (persistent) |
| **Secrets Manager** | Custom + Moto fallback | Pebble (persistent) |
| **SSM** | Custom + Moto fallback | Pebble (persistent) |
| **SES** | Custom + Moto fallback | Pebble (persistent) |
| **CloudWatch Logs** | Custom + Moto fallback | Pebble (persistent) |
| **CloudWatch Metrics** | v1: Moto, v2: Custom (pure) | Pebble (persistent) |
| **Scheduler** | Custom + Moto fallback | Pebble (persistent) |
| **CloudTrail** | Not implemented | Pebble (persistent) |
| **Athena** | Not implemented | Pebble (persistent) |
| **Timestream** | Not implemented | Pebble (persistent) |
| **Cognito** | Not implemented | Pebble (persistent) |
| **CloudFront** | Not implemented | Pebble (persistent) |
| **WAF** | Not implemented | Pebble (persistent) |
| **EC2** | Custom + Moto fallback | Not implemented |
| **OpenSearch / ES** | Custom + native binary | Not implemented |
| **Firehose** | Custom (pure) | Not implemented |
| **Redshift** | Custom + Moto fallback | Not implemented |
| **Config** | Custom + Moto fallback | Not implemented |
| **Transcribe** | Custom + Moto fallback + Vosk | Not implemented |

### 3.2 Services Unique to Each Project

| LocalStack only | Vorpalstacks only |
|-----------------|-------------------|
| EC2 | Athena |
| OpenSearch / Elasticsearch | CloudTrail |
| Firehose | Timestream (Write + Query) |
| CloudFormation | Cognito (Identity + User Pools) |
| Redshift | CloudFront |
| Config | WAF (v1 + v2) |
| S3 Control | |
| Route53 Resolver | |
| SWF | |
| Resource Groups | |
| Transcribe | |
| Support | |

---

## 4. Implementation Patterns

### LocalStack

**Moto Fallback Pattern** (`localstack-core/localstack/services/providers.py`):
```python
from localstack.services.moto import MotoFallbackDispatcher

def iam():
    from localstack.services.iam.provider import IamProvider
    provider = IamProvider()
    return Service.for_provider(provider, dispatch_table_factory=MotoFallbackDispatcher)
```

**Service categories by implementation:**
- **Pure Custom** (14): S3, SQS, SNS, Lambda, KMS, Step Functions, CloudFormation, CloudWatch v2, EventBridge v2, DynamoDB Streams, Elasticsearch, OpenSearch, Firehose
- **Custom + Moto Fallback** (23): IAM, STS, ACM, API Gateway, CloudWatch v1, CloudWatch Logs, EC2, Route53, Route53 Resolver, Secrets Manager, SSM, SES, S3 Control, Scheduler, Config, Redshift, SWF, Resource Groups, Transcribe, Support, Resource Groups Tagging API, EventBridge v1
- **Custom + External Binary** (3): DynamoDB (DynamoDB Local), Kinesis (kinesis-mock), OpenSearch (native binary)

### Vorpalstacks

**Store Pattern** (`internal/store/aws/common/base_store.go`):
```go
type BaseStore struct {
    storage storage.Storage
    prefix  string
}

func (s *BaseStore) Get(key string, dest interface{}) error {
    data, err := s.storage.Get(s.prefix + key)
    if err != nil {
        return err
    }
    return json.Unmarshal(data, dest)
}
```

**Service Integration Pattern** (main.go):
```go
lambdaService = lambda.NewLambdaServiceWithLogs(store, dockerClient, logsStore, ...)
snsService = sns.NewSNSServiceWithClients(..., lambdaService)
eventsService = events.NewEventsServiceWithClients(..., lambdaService)
```

All 29 service APIs use custom Go implementation with Pebble storage. Services communicate internally via gRPC, with dependency injection for cross-service integration (e.g., EventBridge invoking Lambda).

---

## 5. Key Findings

### LocalStack Implementation Breakdown

| Category | Count | Examples |
|----------|-------|---------|
| **Pure Custom** | 14 | S3, SQS, SNS, Lambda, KMS, Step Functions, CloudFormation |
| **Custom + Moto Fallback** | 23 | IAM, STS, ACM, EC2, Route53, Secrets Manager, SSM, SES |
| **External Binary** | 3 | DynamoDB (DynamoDB Local), Kinesis (kinesis-mock), OpenSearch |

### Vorpalstacks Implementation Summary

| Aspect | Detail |
|--------|--------|
| **Service APIs** | 29 (covering 25 AWS services) |
| **Implemented Operations** | 1,081 |
| **Storage** | Pebble v2.1.4 (persistent KV store, all services) |
| **Lambda Execution** | Docker containers via Moby client |
| **Service Integration** | EventBridge → Lambda, SNS → Lambda, SQS → Lambda, etc. |
| **SDK Test Coverage** | 594 tests across 25 services (100% pass rate) |
| **Implementation** | Custom Go code for all services |

---

## 6. Performance Characteristics

### Deployment Comparison

| Metric | LocalStack | Vorpalstacks |
|---------|-----------|--------------|
| **Binary size** | N/A (Docker image) | ~138 MB (single binary) |
| **Deployment** | Docker required | Single binary, Docker optional (Lambda only) |
| **Runtime Dependencies** | Python, Docker, Java (DynamoDB Local), kinesis-mock | Docker (Lambda only) |

### Appendix: Vorpalstacks Benchmark Results (Reference)

**Platform**: AMD Ryzen 7 5700U (16 cores), Linux, Go 1.25.8, Pebble v2.1.4
**Date**: 2026-03-27, 5 iterations per benchmark

> **Note**: These figures are environment-dependent and provided for reference only. Direct comparison with other systems is not meaningful without identical hardware, configuration, and workload.

| Service | Operation | Avg Latency | Memory/op | ops/sec |
|---------|-----------|-------------|-----------|---------|
| **DynamoDB** | PutItem | 1.02ms | 34.2KB | ~980 |
| | GetItem | 0.38ms | 36.4KB | ~2,630 |
| | UpdateItem | 0.71ms | 35.2KB | ~1,410 |
| | Query (100 items) | 6.81ms | 253KB | ~147 |
| | Scan (100 items) | 6.63ms | 252KB | ~151 |
| | BatchWriteItem (25 items) | 8.37ms | 59.6KB | ~119 batches/sec |
| | ParallelGetItem (10) | 0.11ms | 49.4KB | ~9,090 |
| **SQS** | SendMessage | 0.67ms | 34.1KB | ~1,490 |
| | ReceiveMessage | 16.7ms | 70.9KB | ~60 |
| | SendMessageBatch (10 msgs) | 3.73ms | 56.1KB | ~268 batches/sec |
| | ParallelSendReceive (10) | 0.35ms | 35.1KB | ~2,860 |
| **S3** | PutObject (1KB) | 1.49ms | 44.6KB | ~671 |
| | GetObject (1KB) | 0.27ms | 38.1KB | ~3,700 |
| | HeadObject (1KB) | 0.25ms | 35.4KB | ~4,000 |
| | ListObjects (100 items) | 10.3ms | 1.71MB | ~97 |
| | DeleteObject | 1.86ms | 34.8KB | ~538 |
| | ParallelGetObject (10) | 0.11ms | 41.5KB | ~9,260 |
| **Lambda** | Invoke (Node.js 20, cold start) | 200ms | 40.0KB | ~5 |
| | GetFunction | 0.29ms | 35.9KB | ~3,450 |
| | ListFunctions | 0.29ms | 35.1KB | ~3,450 |
| | ParallelInvoke (10) | 32.1ms | 39.3KB | ~31 |
| **Athena** | StartQueryExecution | 2.84ms | 32.6KB | ~352 |
| | ListWorkGroups | 0.27ms | 33.4KB | ~3,700 |
| | GetWorkGroup | 0.28ms | 33.6KB | ~3,570 |
| | ListQueryExecutions | 2.76ms | 45.5KB | ~362 |
| | ParallelStartQuery (10) | 2.78ms | 33.9KB | ~360 |
| **Timestream** | WriteRecords | 0.37ms | 36.7KB | ~2,700 |
| | ListDatabases | 0.32ms | 35.6KB | ~3,125 |
| | ListTables | 0.34ms | 36.5KB | ~2,940 |
| | Query | 0.23ms | 35.2KB | ~4,350 |
| | ParallelWriteRecords (10) | 0.11ms | 45.1KB | ~9,090 |

---

## 7. Summary

| Aspect | LocalStack | Vorpalstacks |
|--------|-----------|--------------|
| **Primary Use Case** | Development and testing | Edge and on-premises deployment |
| **Deployment Model** | Docker container (required) | Single binary |
| **Language** | Python | Go |
| **Service APIs** | 30 | 29 |
| **Implementation** | Hybrid (Custom + Moto + External) | All custom |
| **Data Persistence** | Mixed (in-memory + file) | Pebble (persistent, all services) |
| **Lambda Execution** | Docker (real) | Docker (real) |
| **SDK Test Coverage** | Not publicly documented | 594/594 (100%) |
| **License** | Apache 2.0 | FSL-1.1-MIT / Apache 2.0 (pkg/)¹ |
| **External Dependencies** | Docker + Java + kinesis-mock | Docker (Lambda only) |

Neither project is a superset of the other. LocalStack covers infrastructure services (EC2, CloudFormation, OpenSearch, Firehose) not present in Vorpalstacks, while Vorpalstacks covers observability and analytics services (CloudTrail, Athena, Timestream, CloudFront, WAF, Cognito) not present in LocalStack.

---

¹ The root licence is planned to be changed to MIT once the project reaches production stability after the initial bug-fix phase.

*Analysis Date: 2026-03-27*
*Information Source: Source code analysis of LocalStack v4.14.0 and Vorpalstacks (main branch)*

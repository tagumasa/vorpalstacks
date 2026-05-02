# Technical Comparison Report: LocalStack vs Vorpalstacks

*Analysis Date: 2026-04-20*

---

## Executive Summary

This report compares LocalStack (v2026.03) and Vorpalstacks as AWS service emulators for development, testing, and deployment.

As of March 2026, LocalStack transitioned from open-source to a commercial model. An auth token (account) is now required to run any edition. The free Hobby plan is restricted to non-commercial use and excludes many services. Paid plans (Base, Ultimate, Enterprise) unlock additional services, state persistence, and enterprise features.

Vorpalstacks is a self-contained single-binary emulator written in Go. It requires no account or auth token, uses persistent storage for all services, and is licensed under FSL-1.1-MIT (planned MIT after production stability). It targets 32 AWS service APIs with custom implementations.

The two projects have partially overlapping but distinct service coverage. Neither is a superset of the other. The choice between them depends on which AWS services are required, deployment constraints, licensing needs, and whether commercial use is intended.

---

## 1. Architecture Overview

### LocalStack (v2026.03)

| Aspect | Details |
|--------|---------|
| **Language** | Python (>=3.10) |
| **Version** | v2026.03 (calendar versioning; last semver was v4.14.0) |
| **Deployment** | Docker container (required) |
| **Auth** | Auth token required for all editions |
| **Core Dependencies** | Python, moto-ext, DynamoDB Local v2.x (Java), kinesis-mock |
| **Service Implementation** | Hybrid: Custom providers + Moto fallback + External binaries |
| **State Persistence** | In-memory by default; local persistence and Cloud Pods on Base+ plans only |
| **License** | Source-available (non-OSS since March 2026) |

**Implementation tiers:**

| Tier | Services | Count |
|------|----------|-------|
| Pure Custom | S3, SQS, SNS, Lambda, KMS, Step Functions, CloudFormation, CloudWatch v2, EventBridge v2, DynamoDB Streams, Elasticsearch, OpenSearch, Firehose | 14 |
| Custom + Moto Fallback | IAM, STS, ACM, API Gateway, CloudWatch v1, CloudWatch Logs, EC2, Route53, Route53 Resolver, Secrets Manager, SSM, SES, S3 Control, Scheduler, Config, Redshift, SWF, Resource Groups, Transcribe, Support, Resource Groups Tagging API, EventBridge v1 | 23 |
| Custom + External Binary | DynamoDB (DynamoDB Local — Java), Kinesis (kinesis-mock), OpenSearch (native binary) | 3 |
| External Runtime | Lambda (lambda-runtime-init), Step Functions (JPype/JVM), Transcribe (Vosk + FFmpeg) | — |

### Vorpalstacks

| Aspect | Details |
|--------|---------|
| **Language** | Go 1.25+ |
| **Deployment** | Single binary (110 MB stripped, 154 MB debug); Docker optional (Lambda execution only) |
| **Auth** | None required; fully self-contained |
| **Core Dependencies** | Pebble v2.1.4 (persistent KV store), Moby client (Docker, Lambda only) |
| **Service Implementation** | Custom Go implementation for all services |
| **State Persistence** | Pebble persistent KV store for all services (region-specific and global) |
| **License** | FSL-1.1-MIT (root), Apache 2.0 (pkg/); planned MIT after production stability |
| **Admin** | gRPC-Web admin API, `vstacks` CLI, PebbleDB-backed runtime configuration |

**Architecture:**

```
Vorpalstacks (32 AWS service APIs, 1,255 operations)
├── All Custom Implementation (Go)
│   ├── IAM (159), API Gateway (80), S3 (71),
│   ├── AppSync (75), Neptune (70), Cognito User Pools (68),
│   ├── SESv2 incl. SESv1 (63), DynamoDB (57),
│   ├── Lambda (50), KMS (45), NeptuneData (43),
│   ├── EventBridge (40), Step Functions (39), Kinesis (39),
│   ├── CloudFront (38), WAFv2 (37), Athena (37),
│   ├── NeptuneGraph (34), SNS (31),
│   ├── CloudWatch Logs (29), CloudTrail (24), Secrets Manager (24),
│   ├── SQS (23), Cognito Identity (23), CloudWatch (21),
│   ├── Timestream Write (20), Route53 (20),
│   ├── ACM (16), Timestream Query (15),
│   ├── SSM (12), Scheduler (12), STS (11)
│   └── EC2 helpers only (VPC, Security Groups, Subnets — no API handlers)
├── Pebble Storage (persistent KV, all services)
│   ├── Region-specific (us-east-1, us-west-2, etc.)
│   ├── Global (IAM, Route53, CloudFront, STS)
│   └── Graph (NeptuneGraph)
├── Internal Event Bus
│   EventBridge→Lambda, SNS→Lambda, SQS→Lambda,
│   Cognito triggers→Lambda, Secrets Manager rotation→Lambda,
│   CloudWatch alarm evaluation
├── Docker Runtime (Lambda execution)
├── IAM Authorization (policy-based access control)
├── gRPC-Web Admin API + vstacks CLI
└── Configurable CORS (multi-origin)
```

---

## 2. Service Coverage Comparison

### 2.1 Services in Both Projects (Hobby-tier comparison)

| Category | Service | LocalStack (Hobby) | Vorpalstacks |
|----------|---------|--------------------:|:------------:|
| **Analytics** | Athena | ❌ | ✅ |
| | Kinesis Streams | ✅ | ✅ |
| | Kinesis Firehose | ✅ | ❌ |
| | OpenSearch | ✅ | ❌ |
| | Redshift | ✅ | ❌ |
| **App Integration** | EventBridge | ✅ | ✅ |
| | EventBridge Scheduler | ✅ | ✅ |
| | SNS | ✅ | ✅ |
| | SQS | ✅ | ✅ |
| | Step Functions | ✅ | ✅ |
| | SWF | ✅ | ❌ |
| **Business** | SES (v1) | ✅ | ✅¹ |
| | SESv2 | ❌ (Base+) | ✅¹ |
| **Compute** | Lambda | ✅ | ✅ |
| | EC2 | ✅ | ❌³ |
| **Databases** | DynamoDB | ✅ | ✅ |
| | Neptune | ❌ (Base+) | ✅ |
| **Management & Governance** | CloudFormation | ✅ | ❌ |
| | CloudWatch Metrics | ✅ | ✅ |
| | CloudWatch Logs | ✅ | ✅ |
| | CloudTrail | ❌ (Base+) | ✅ |
| | Config | ✅ | ❌ |
| | Resource Groups | ✅ | ❌ |
| | SSM | ✅ | ✅ |
| **Networking** | API Gateway | ✅ (REST only) | ✅ |
| | CloudFront | ❌ (Base+) | ✅ |
| | Route53 | ✅ | ✅ |
| | Route53 Resolver | ✅ | ❌ |
| **Security** | ACM | ✅ | ✅ |
| | Cognito (Identity + User Pools) | ❌ (Base+) | ✅ |
| | IAM | ✅ | ✅² |
| | KMS | ✅ | ✅ |
| | Secrets Manager | ✅ | ✅ |
| | STS | ✅ | ✅ |
| | WAFv2 | ❌ (Base+) | ✅ |
| **Storage** | S3 | ✅ | ✅ |
| | S3 Control | ✅ | ❌ |
| | Timestream (Write + Query) | ❌ (Base+) | ✅ |

¹ Vorpalstacks implements SESv2, which also handles SESv1 Query protocol.
 ² Vorpalstacks implements IAM with policy-based access control enforcement.
³ Vorpalstacks includes EC2-related helper operations (VPC, Security Groups, Subnets) but does not register any EC2 API handlers.

**Totals:** LocalStack Hobby covers 27 of the services listed above; Vorpalstacks covers 27. Many of the services marked ❌ for LocalStack Hobby are available on paid plans (Base, Ultimate, Enterprise).

### 2.2 LocalStack Paid-Only Services (Base or higher required)

These services are not available on the free Hobby plan:

| Category | Services |
|----------|----------|
| Analytics | Athena, EMR, Glue, Lake Formation, MSK, Flink |
| App Integration | SESv2, Pinpoint, Amplify, AppSync, EventBridge Pipes |
| Business | — |
| Compute | ECS, EKS, ECR, Batch |
| Customer Enablement | — |
| Databases | Neptune, RDS, ElastiCache, DocumentDB, MemoryDB, Timestream |
| Machine Learning | Bedrock, SageMaker, Textract, Transcribe (paid) |
| Management & Governance | CloudTrail, Account Management, Organizations, CodeCommit/Build/Deploy/Pipeline |
| Media | MediaConvert |
| Networking | CloudFront, ELB, Cloud Map |
| Security | WAF, Cognito (both), Shield, Private CA, IAM Identity Store/Center |
| Storage | S3 Glacier, EFS, Backup, Transfer Family, DMS |
| Other | IoT (all), MQ, Managed Blockchain, Cost Explorer, Airflow, Fault Injection, AppConfig, Auto Scaling, RAM, Resiliency Testing, K8s Operator, IAM Policy Enforcement |

### 2.3 Services Unique to Each Project

#### LocalStack Only (across all plans)

| Service | Notes |
|---------|-------|
| CloudFormation | Infrastructure-as-code orchestration |
| ECS / EKS / ECR | Container orchestration and registry |
| OpenSearch / Elasticsearch | Search and analytics engine (native binary) |
| Firehose | Kinesis data delivery |
| Redshift | Data warehousing |
| RDS / ElastiCache / DocumentDB / MemoryDB | Managed databases |
| EMR / Glue / Lake Formation / MSK / Flink | Big data and analytics |
| IoT (all) | IoT device management |
| Amplify | Frontend hosting |
| Bedrock / SageMaker / Textract | ML and AI services |
| MediaConvert | Video transcoding |
| Transfer Family / DMS | Data transfer and migration |
| Cloud Map | Service discovery |
| CodeBuild/Deploy/Pipeline | CI/CD |
| X-Ray | Distributed tracing |
| ELB / Auto Scaling | Load balancing and scaling |
| S3 Glacier / EFS / Backup | Additional storage options |
| Cost Explorer | Cost analysis |
| Airflow | Workflow orchestration |
| MQ | Message broker |
| Managed Blockchain | Blockchain networks |
| SWF | Simple Workflow Service |
| Config | Resource configuration tracking |
| S3 Control | S3 account-level operations |
| Route53 Resolver | DNS resolver |
| Support API | AWS support |
| Resource Groups | Resource grouping |
| Cloud Pods | State snapshot sharing (Base+) |
| K8s Operator | Kubernetes integration |

#### Vorpalstacks Only

| Service | Notes |
|---------|-------|
| Athena | Interactive query engine |
| CloudTrail | API call logging and auditing |
| Timestream (Write + Query) | Time-series database |
| Cognito (Identity + User Pools) | Authentication and identity management |
| CloudFront | Content delivery network |
| WAFv2 | Web application firewall |
| Neptune (Management + Data + Graph) | Graph database |
| AppSync | Managed GraphQL API |
| IAM Authorization | Policy-based access control enforcement |

---

## 3. Licensing and Access

| Aspect | LocalStack | Vorpalstacks |
|--------|-----------|--------------|
| **License Model** | Commercial (source-available since March 2026) | FSL-1.1-MIT (root), Apache 2.0 (pkg/) |
| **Auth Required** | Yes — account/token for all editions | No |
| **Free Tier** | Hobby plan (non-commercial only) | Full functionality, no restrictions |
| **Commercial Use (Free)** | Not permitted on Hobby plan | Permitted |
| **Paid Plans** | Base ($/month), Ultimate ($/month), Enterprise | None — fully self-contained |
| **State Persistence (Free)** | Not available on Hobby | All services, persistent Pebble storage |
| **Telemetry** | Enforced on Hobby; Enterprise can disable | No telemetry |
| **Planned Licence Change** | N/A | MIT after production stability |

LocalStack's Hobby plan restricts usage to non-commercial purposes and lacks state persistence, many services, and enterprise features. Vorpalstacks imposes no usage restrictions and requires no account.

---

## 4. Deployment and Operations

| Aspect | LocalStack | Vorpalstacks |
|--------|-----------|--------------|
| **Binary Size** | N/A (Docker image) | 110 MB (stripped), 154 MB (debug) |
| **Docker Required** | Yes (mandatory) | Optional (Lambda execution only) |
| **Runtime Dependencies** | Python, Docker, Java (DynamoDB Local), kinesis-mock | Docker (Lambda only) |
| **Data Persistence** | In-memory default; local + Cloud Pods on Base+ | Pebble persistent KV (all services) |
| **Configuration** | Environment variables, LocalStack configuration files | gRPC-Web admin API, `vstacks` CLI, PebbleDB-backed config |
| **CORS** | Default behaviour | Configurable multi-origin |
| **Multi-Region** | Supported (per-container) | Supported (region-specific Pebble buckets) |
| **Startup Time** | Docker pull + container startup | Single binary execution |
| **Process Model** | Docker container | Native OS process |

---

## 5. Implementation Approach

### LocalStack: Hybrid Implementation

LocalStack uses a layered approach combining three implementation strategies:

1. **Custom providers**: Services implemented directly in Python (e.g., S3, SQS, SNS, Lambda)
2. **Moto fallback**: Operations delegated to the Moto mock library when the custom provider lacks coverage (e.g., IAM, STS, ACM, EC2, Route53)
3. **External binaries**: Real service binaries embedded or proxied (e.g., DynamoDB Local Java process, kinesis-mock, OpenSearch native binary)

This hybrid approach provides broad API coverage by falling back to Moto for edge cases, but introduces behavioural inconsistencies between the custom and Moto layers. External binaries add runtime dependency overhead (Java for DynamoDB Local).

### Vorpalstacks: All-Custom Implementation

All 32 service APIs are implemented in Go with no external service binaries or mock library fallbacks. Every operation is backed by Pebble persistent storage. Cross-service integration uses an internal event bus for service-to-service communication (e.g., EventBridge invoking Lambda, SNS triggering Lambda, CloudWatch alarm evaluation).

This approach provides consistent behaviour across all services and persistent state without external dependencies beyond Docker for Lambda execution. The trade-off is that each service must be implemented from scratch, which limits total service count compared to a hybrid approach.

---

## 6. Performance Characteristics

### Deployment Comparison

| Metric | LocalStack | Vorpalstacks |
|---------|-----------|--------------|
| **Binary size** | N/A (Docker image) | 110 MB (stripped), 154 MB (debug) |
| **Deployment** | Docker required | Single binary, Docker optional (Lambda only) |
| **Runtime Dependencies** | Python, Docker, Java (DynamoDB Local), kinesis-mock | Docker (Lambda only) |

### SDK Test Coverage

| Project | Test Count | Details |
|---------|-----------|---------|
| LocalStack | Not publicly documented | — |
| Vorpalstacks | 2,262 tests | 2,216 SDK + 29 integration + 17 WebSocket (100% pass rate) |

### Appendix: Vorpalstacks Benchmark Results (Reference)

**Platform**: AMD Ryzen 7 5700U (16 cores), Linux, Go 1.25.8, Pebble v2.1.4
**Date**: 2026-04-21, 1 iteration per benchmark, fresh server/data per service

> **Note**: These figures are environment-dependent and provided for reference only. Direct comparison with other systems is not meaningful without identical hardware, configuration, and workload.

| Service | Operation | Avg Latency | ops/sec |
|---------|-----------|-------------|---------|
| **DynamoDB** | PutItem | 1.08ms | ~926 |
| | GetItem | 0.87ms | ~1,150 |
| | UpdateItem | 1.22ms | ~820 |
| | Query (1k items) | 6.62ms | ~151 |
| | Scan (1k items) | 6.08ms | ~164 |
| | BatchWriteItem (25 items) | 8.71ms | ~115 batches/sec |
| | ParallelGetItem (10) | 1.05ms | ~952 |
| **SQS** | SendMessage | 0.82ms | ~1,220 |
| | ReceiveMessage | 12.04ms | ~83 |
| | SendMessageBatch (10 msgs) | 1.06ms | ~943 batches/sec |
| | ParallelSendReceive (10) | 0.95ms | ~1,053 |
| **S3** | PutObject | 1.27ms | ~787 |
| | GetObject | 0.96ms | ~1,042 |
| | HeadObject | 0.83ms | ~1,205 |
| | ListObjects (100 items) | 2.25ms | ~444 |
| | DeleteObject | 0.76ms | ~1,316 |
| | ParallelGetObject (10) | 0.82ms | ~1,220 |
| **Lambda** | Invoke (Node.js 20, cold start) | 317ms | ~3 |
| | GetFunction | 1.15ms | ~870 |
| | ListFunctions | 0.80ms | ~1,250 |
| | ParallelInvoke (10) | 70.4ms | ~14 |
| **Athena** | StartQueryExecution | 0.90ms | ~1,111 |
| | ListWorkGroups | 0.82ms | ~1,220 |
| | GetWorkGroup | 0.85ms | ~1,176 |
| | ListQueryExecutions | 0.75ms | ~1,333 |
| | ParallelStartQuery (10) | 0.85ms | ~1,176 |
| **Timestream** | WriteRecords | 14.1ms | ~71 |
| | ListDatabases | 1.32ms | ~758 |
| | ListTables | 0.90ms | ~1,111 |
| | Query | 11.8ms | ~85 |
| | ParallelWriteRecords (10) | 19.5ms | ~51 |
| **AppSync** | Introspection_Schema | 0.83ms | ~1,205 |
| | Query_NilResolver | 3.90ms | ~256 |
| | Mutation_Single (create+delete) | 1.12ms | ~893 |
| | WS Connect_Ack | 0.86ms | ~1,163 |
| | WS Publish_Subscribe_Delivery | 1.01ms | ~990 |
| **NeptuneGraph** | OpenCypher_Simple | 1.25ms | ~800 |
| | OpenCypher_MultiHop (100 nodes, 100 edges) | 3.13ms | ~319 |
| | GetGraphSummary (100 nodes, 100 edges) | 0.80ms | ~1,250 |
| **NeptuneData** | GetEngineStatus | 0.66ms | ~1,515 |
| | GremlinQuery_Simple | 1.13ms | ~885 |
| | GremlinQuery_Neighbor | 1.77ms | ~565 |
| | CypherQuery_Simple | 1.42ms | ~704 |
| | GetPropertygraphSummary (100 nodes, 100 edges) | 0.72ms | ~1,389 |

---

## 7. Summary

| Aspect | LocalStack | Vorpalstacks |
|--------|-----------|--------------|
| **Primary Target** | Development and testing | Edge and on-premises deployment |
| **Version** | v2026.03 (calendar versioning) | Continuous (Git main) |
| **Language** | Python | Go 1.25+ |
| **Deployment** | Docker container (required) | Single binary (110 MB stripped) |
| **Auth Required** | Yes (all editions) | No |
| **Free Tier** | Hobby (non-commercial, limited services) | Full functionality, no restrictions |
| **Commercial Use (Free)** | Not permitted | Permitted |
| **Service APIs** | ~50+ (across all plans) | 32 |
| **Implementation** | Hybrid (Custom + Moto + External binaries) | All custom Go |
| **Data Persistence** | In-memory default; paid plans add local + Cloud Pods | Pebble persistent KV (all services) |
| **Lambda Execution** | Docker (real) | Docker (real) |
| **IAM Enforcement** | Base+ plans only | Policy-based (all editions) |
| **SDK Test Coverage** | Not publicly documented | 2,262 tests (100% pass rate) |
| **Runtime Dependencies** | Python, Docker, Java, kinesis-mock | Docker (Lambda only) |
| **CORS** | Default | Configurable multi-origin |
| **Admin Interface** | LocalStack Web UI (paid) | gRPC-Web admin API + `vstacks` CLI |
| **Telemetry** | Enforced (Hobby); disableable (Enterprise) | None |
| **License** | Source-available (commercial) | FSL-1.1-MIT / Apache 2.0 (pkg/)¹ |

Neither project is a superset of the other. LocalStack covers a wider range of infrastructure services (EC2, CloudFormation, container orchestration, managed databases) across its paid plans, while Vorpalstacks covers analytics, observability, and security services (Athena, CloudTrail, Timestream, CloudFront, WAFv2, Cognito, Neptune, AppSync) that require paid LocalStack plans or are not available at all.

Key differentiators:
- **No auth required**: Vorpalstacks runs without an account or token; LocalStack requires authentication for all editions since March 2026.
- **Commercial use**: Vorpalstacks permits commercial use at no cost; LocalStack's Hobby plan is restricted to non-commercial use.
- **Deployment**: Vorpalstacks is a single binary with no mandatory Docker dependency; LocalStack requires Docker.
- **Persistence**: Vorpalstacks persists all service data by default; LocalStack Hobby has no persistence.
- **Service breadth**: LocalStack (paid plans) covers significantly more services; Vorpalstacks covers services not available in LocalStack Hobby.
- **Implementation consistency**: Vorpalstacks uses a single implementation approach (custom Go + Pebble); LocalStack mixes custom providers, Moto fallback, and external binaries, which can lead to behavioural inconsistencies.

---

¹ The root licence is planned to be changed to MIT once the project reaches production stability after the initial bug-fix phase.

*Information Sources: LocalStack official documentation (accessed 2026-04-20), Vorpalstacks source code (main branch).*

# Vorpalstacks

[日本語](README.ja.md) | [中文](README.zh.md)

> **Warning: This is a beta release.** Vorpalstacks is under active development. While 32 AWS services are implemented with 1962 passing SDK tests, 24 cross-service integration tests, and 17 WebSocket tests (2003 total, plus 631 Python, 629 TypeScript, 606 C#), not all edge cases and AWS behaviours are fully covered. Expect breaking changes. Bug reports and contributions are welcome.

A lightweight edge and on-premise cloud platform providing AWS-compatible services.

## Overview

Vorpalstacks enables running AWS-compatible services in environments where full AWS connectivity is not available:

- Edge computing scenarios
- On-premise deployments
- Development and testing environments
- Air-gapped networks

## Features

> **What this is**: A real implementation of AWS-compatible APIs, not a mock framework. Each service stores data in PebbleDB and supports multi-region isolation.

> **What this is not**: A fully faithful reproduction of every AWS behaviour. Some edge cases, undocumented behaviours, and advanced features may differ from AWS. See [docs/services.md](docs/services.md) for the current scope of each service.

- **AWS API Compatible**: Works with existing AWS SDKs and CLI
- **Thirty-two AWS Services**: S3, SQS, SNS, Lambda, DynamoDB, API Gateway, AppSync, Step Functions, WAFv2, Kinesis, KMS, Neptune, Neptune Graph, and more
- **IAM Authorization**: Full IAM policy evaluation with user/group/role-based access control
- **DynamoDB PartiQL**: SQL-like queries with WHERE functions (attribute_exists, begins_with, contains, size)
- **S3 SelectObjectContent**: SQL queries on CSV/JSON objects with event streaming
- **Multi-Region Support**: Region-isolated storage with dedicated PebbleDB per region
- **Global Services**: IAM, Route53, CloudFront, STS with shared global storage
- **gRPC-Web Admin API**: Connect-RPC admin interface on separate port for all services
- **Lightweight**: Single binary, minimal dependencies
- **Persistent Storage**: Pebble-based key-value store
- **Docker Integration**: Lambda functions run in containers
- **Service Integration**: Event-driven communication between services with 24 cross-service integration tests
- **TLS Support**: Optional HTTPS with auto-generated or custom certificates
- **LocalStack Comparison**: See [docs/localstack_vs_vorpalstacks_report.md](docs/localstack_vs_vorpalstacks_report.md) for a technical comparison with LocalStack

## Implemented Services

| Service | Coverage | Notes |
|---------|----------|-------|
| ACM | Full | |
| API Gateway | Broad | No client certificates, documentation, or SDK generation |
| AppSync | Full | 74 control-plane operations, GraphQL execution, WebSocket pub/sub |
| Athena | Broad | No capacity reservations or notebook sessions |
| CloudFront | Selective | No actual edge traffic distribution |
| CloudTrail | Broad | No event data stores or SQL queries |
| CloudWatch Logs | Selective | No Logs Insights queries or export |
| CloudWatch Metrics | Broad | No metric streams or anomaly detection |
| Cognito IDP | Selective | No external IdP, no hosted UI |
| Cognito Identity | Selective | Basic identity pool support |
| DynamoDB | Full | |
| EventBridge | Broad | No global endpoints or partner event sources |
| IAM | Broad | No policy simulator or organisations integration |
| Kinesis | Full | |
| KMS | Full | |
| Lambda | Broad | No durable functions or code signing |
| Neptune | Full | Property graph + RDF, openCypher/Gremlin, bulk loader |
| Neptune Graph | Full | 34 control-plane operations, openCypher, vector search |
| Route53 | Selective | DNS record management only |
| S3 | Broad | No analytics, inventory, or S3 Express |
| Scheduler | Full | |
| Secrets Manager | Full | |
| SESv2 | Broad | No deliverability testing or multi-tenancy |
| SNS | Broad | SMS sending not supported |
| SQS | Full | |
| SSM | Selective | Parameter Store only |
| STS | Full | |
| Step Functions | Full | |
| Timestream | Full | |
| WAFv2 | Broad | |

See [docs/services.md](docs/services.md) for detailed coverage tiers and service integration patterns.

## Quick Start

### Build

```bash
go build -o vorpalstacks .
```

### Run (Development Mode)

```bash
SIGNATURE_VERIFICATION_ENABLED=false DATA_PATH=./tmp/data ./vorpalstacks
```

### Run with Docker (for Lambda)

```bash
SIGNATURE_VERIFICATION_ENABLED=false DATA_PATH=./tmp/data ./vorpalstacks
```

### Use with AWS CLI

```bash
export AWS_ENDPOINT_URL=http://localhost:8080

aws --endpoint-url=http://localhost:8080 --no-sign-request sns list-topics
aws --endpoint-url=http://localhost:8080 --no-sign-request sqs list-queues
aws --endpoint-url=http://localhost:8080 --no-sign-request lambda list-functions
```

## Testing

### Unit Tests

```bash
make test
```

### SDK Tests (AWS Go SDK v2)

```bash
SIGNATURE_VERIFICATION_ENABLED=false TEST_MODE=true DATA_PATH=./data-test ./vorpalstacks &

cd sdk-tests && ./sdk-tests-test -service all -v
```

### CLI Integration Tests

```bash
cd scripts/services && bash test_sqs.sh
cd scripts/services && bash test_dynamodb.sh
cd scripts/services && bash test_s3.sh
cd scripts/services && bash test_iam.sh
```

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `DATA_PATH` | `./data` | Path for persistent data storage |
| `AWS_REGION` | `us-east-1` | Default region |
| `AWS_ACCOUNT_ID` | `000000000000` | AWS account ID |
| `SIGNATURE_VERIFICATION_ENABLED` | `true` | Enable AWS Signature V4 verification |
| `GRPC_WEB_PORT` | `9090` | gRPC-Web admin server port |
| `TLS_ENABLED` | `false` | Enable TLS |
| `AUTHORIZATION_ENABLED` | `false` | Enable IAM policy evaluation |

See [Configuration](docs/configuration.md) for the complete list.

### Example: Production Configuration

```bash
export AWS_ACCESS_KEY_ID=your-key-id
export AWS_SECRET_ACCESS_KEY=your-secret-key
export SIGNATURE_VERIFICATION_ENABLED=true
export AUTHORIZATION_ENABLED=true
export DATA_PATH=/var/lib/vorpalstacks
./vorpalstacks
```

## Data Storage Structure

```
DATA_PATH/
├── us-east-1/               # Region-specific storage
│   ├── pebble/              # PebbleDB for us-east-1
│   ├── code/                # Lambda function code
│   └── logs/                # CloudWatch Logs chunks
├── global/                  # Global services (IAM, Route53, CloudFront, STS)
│   └── pebble/
└── uploads/                 # S3 multipart uploads (temporary)
```

Each region has isolated storage including PebbleDB, Lambda code, and logs. Global services share dedicated storage.

## Docker Requirements

For Lambda functionality:

1. Docker must be installed and running
2. Lambda runtime images will be pulled automatically
3. Docker socket must be accessible at `DOCKER_HOST` (default: `unix:///var/run/docker.sock`)

## Documentation

- [Architecture](docs/architecture.md) - System architecture overview
- [Services](docs/services.md) - Implemented AWS services
- [Configuration](docs/configuration.md) - Environment variables and settings
- [Integration](docs/integration.md) - Service-to-service communication

## Known Limitations

- Not all AWS operations are implemented for every service — see [docs/services.md](docs/services.md) for details
- Some edge cases and undocumented AWS behaviours may differ
- No cross-account access control (single-account mode)
- CloudFront distributions do not serve actual edge traffic
- Cognito hosted UI domains are not supported (requires CloudFront edge)
- SQS FIFO queues have limited support
- DynamoDB Streams and Global Tables are not implemented

## Roadmap

- **Short-term**: Bug fixes, refactoring, and stability improvements
- **Service expansion**: AWS IoT Core
- **Terraform**: 28 services have passed basic conformance testing — see [vorpalstacks-conformance-tests](https://github.com/tagumasa/vorpalstacks-conformance-tests) for details and how to run

See [CHANGELOG.md](CHANGELOG.md) for release history.

## Requirements

- Go 1.25+
- Docker (for Lambda functionality)

## Performance

Vorpalstacks implements all 32 services as native Go binaries backed by PebbleDB, avoiding the overhead of interpreted languages or external process dependencies.

This architecture enables sub-millisecond latencies for core operations, making it practical to run extensive API tests (1962 SDK + 24 integration + 17 WebSocket Go tests, 631 Python, 629 TypeScript, 606 C# tests) directly within CI/CD pipelines without containerization overhead.

### Benchmark Results (Reference)

Platform: AMD Ryzen 7 5700U (16 cores), Linux, Go 1.25.8, Pebble v2.1.4

> **Note**: These figures are environment-dependent. Direct comparison with other systems is not meaningful without identical hardware, configuration, and workload.

| Service | Operation | Avg Latency | ops/sec |
|---------|-----------|-------------|---------|
| DynamoDB | GetItem | 0.38ms | ~2,630 |
| S3 | GetObject (1KB) | 0.27ms | ~3,700 |
| SQS | SendMessage | 0.67ms | ~1,490 |

## License

This project is licensed under the [Functional Source License, Version 1.1, MIT Future License (FSL-1.1-MIT)](LICENSE).

> **Note**: The root licence will change to MIT after the project reaches production stability.

The `pkg/` directory contains code licensed under Apache License 2.0 — see `pkg/sqlparser/LICENSE.md` and `pkg/vsjwt/LICENSE` for details.

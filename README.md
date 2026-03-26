# Vorpalstacks

> **Warning: This is a beta release.** Vorpalstacks is under active development. While 30 AWS services are implemented with 594 passing SDK tests, not all edge cases and AWS behaviours are fully covered. Expect breaking changes. Bug reports and contributions are welcome.

A lightweight edge and on-premise cloud platform providing AWS-compatible services.

## Overview

Vorpalstacks enables running AWS-compatible services in environments where full AWS connectivity is not available:

- Edge computing scenarios
- On-premise deployments
- Development and testing environments
- Air-gapped networks

## Features

> **What this is**: A real implementation of AWS-compatible APIs, not a mock framework. Each service stores data in PebbleDB and supports multi-region isolation.

> **What this is not**: A 100% faithful reproduction of every AWS behaviour. Some edge cases, undocumented behaviours, and advanced features may differ from AWS. See [docs/services.md](docs/services.md) for the current scope of each service.

- **AWS API Compatible**: Works with existing AWS SDKs and CLI
- **30 AWS Services**: S3, SQS, SNS, Lambda, DynamoDB, API Gateway, Step Functions, WAF, WAFv2, Kinesis, KMS, and more
- **IAM Authorization**: Full IAM policy evaluation with user/group/role-based access control
- **DynamoDB PartiQL**: SQL-like queries with WHERE functions (attribute_exists, begins_with, contains, size)
- **S3 SelectObjectContent**: SQL queries on CSV/JSON objects with event streaming
- **Multi-Region Support**: Region-isolated storage with dedicated PebbleDB per region
- **Global Services**: IAM, Route53, CloudFront, STS with shared global storage
- **gRPC-Web Admin API**: Connect-RPC admin interface on separate port for all services
- **Lightweight**: Single binary, minimal dependencies
- **Persistent Storage**: Pebble-based key-value store
- **Docker Integration**: Lambda functions run in containers
- **Service Integration**: Event-driven communication between services
- **TLS Support**: Optional HTTPS with auto-generated or custom certificates

## Implemented Services

| Service | Protocol | Operations | Scope |
|---------|----------|------------|-------|
| ACM | REST-JSON 1.1 | 16 | Regional |
| API Gateway | REST-JSON 1.1 | 80 | Regional |
| Athena | AWS JSON 1.1 | 37 | Regional |
| CloudFront | REST-XML | 32 | Global |
| CloudTrail | AWS JSON 1.1 | 24 | Regional |
| CloudWatch Logs | AWS JSON 1.1 | 29 | Regional |
| CloudWatch Metrics | AWS JSON 1.1 | 17 | Regional |
| Cognito IDP | AWS JSON 1.1 | 51 | Regional |
| Cognito Identity | AWS JSON 1.1 | 13 | Regional |
| DynamoDB | AWS JSON 1.0 | 57 | Regional |
| EventBridge | AWS JSON 1.1 | 39 | Regional |
| IAM | AWS JSON 1.1 | 155 | Global |
| Kinesis | AWS JSON 1.1 | 39 | Regional |
| KMS | AWS JSON 1.1 | 45 | Regional |
| Lambda | AWS JSON 1.1 | 50 | Regional |
| Route53 | AWS JSON 1.1 | 18 | Global |
| S3 | REST-XML | 68 | Regional |
| Scheduler | AWS JSON 1.1 | 12 | Regional |
| SecretsManager | AWS JSON 1.1 | 21 | Regional |
| SESv2 | AWS JSON 1.1 | 58 | Regional |
| SNS | AWS JSON 1.1 | 31 | Regional |
| SQS | AWS Query | 23 | Regional |
| SSM | AWS JSON 1.1 | 12 | Regional |
| STS | AWS JSON 1.1 | 11 | Global |
| Step Functions | AWS JSON 1.1 | 28 | Regional |
| Timestream Query | AWS JSON 1.1 | 13 | Regional |
| Timestream Write | AWS JSON 1.1 | 19 | Regional |
| WAF | AWS JSON 1.1 | 34 | Regional |
| WAFv2 | AWS JSON 1.1 | 36 | Regional |

**SDK Tests**: 594/594 (100%) pass rate

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
- [Services](docs/services.md) - Implemented AWS services and operations
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
- **Terraform**: Improving provider compatibility for existing services

See [CHANGELOG.md](CHANGELOG.md) for release history.

## Requirements

- Go 1.25+
- Docker (for Lambda functionality)

## License

This project is licensed under the [Functional Source License, Version 1.1, MIT Future License (FSL-1.1-MIT)](LICENSE).

The `pkg/` directory contains code licensed under Apache License 2.0 — see `pkg/sqlparser/LICENSE.md` and `pkg/vsjwt/LICENSE` for details.

# Vorpalstacks Documentation

Vorpalstacks is a lightweight edge and on-premise cloud platform providing AWS-compatible services.

## Overview

Vorpalstacks enables running AWS-compatible services in environments where full AWS connectivity is not available or desired:

- Edge computing scenarios
- On-premise deployments
- Development and testing environments
- Air-gapped networks

## Key Features

- **AWS API Compatible**: Works with existing AWS SDKs and CLI
- **30 AWS Services**: S3, SQS, SNS, Lambda, DynamoDB, Kinesis, KMS, and more
- **IAM Authorization**: Full policy-based access control
- **gRPC-Web Admin API**: Connect-RPC admin interface for all services
- **Lightweight**: Single binary, minimal dependencies
- **Persistent Storage**: Pebble-based key-value store
- **Docker Integration**: Lambda functions run in containers
- **TLS Support**: Optional HTTPS with auto-generated certificates
- **Multi-Region**: Region-isolated storage with global services

## Documentation

| Document | Description |
|----------|-------------|
| [Architecture](architecture.md) | System architecture, request flow, storage |
| [Services](services.md) | Implemented AWS services (29 APIs, 25 AWS services, 594/594 tests) |
| [Configuration](configuration.md) | Environment variables, TLS, gRPC-Web, IAM Auth |
| [Integration](integration.md) | Service-to-service communication patterns |

## Quick Start

```bash
go build -o vorpalstacks .
SIGNATURE_VERIFICATION_ENABLED=false DATA_PATH=./tmp/data ./vorpalstacks
```

## AWS CLI Usage

```bash
aws --endpoint-url=http://localhost:8080 sns list-topics
aws --endpoint-url=http://localhost:8080 sqs list-queues
aws --endpoint-url=http://localhost:8080 lambda list-functions
```

## For Developers

See [plans/](../plans/) for development guidelines, API gap analyses, and implementation guides.

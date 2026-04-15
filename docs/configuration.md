# Configuration

This document describes configuration options for Vorpalstacks.

## Environment Variables

### Core Settings

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `DATA_PATH` | `./data` | Path for persistent data storage |
| `AWS_REGION` | `us-east-1` | Default region |
| `AWS_ACCOUNT_ID` | `000000000000` | AWS account ID |
| `AWS_ACCESS_KEY_ID` | - | Default access key ID |
| `AWS_SECRET_ACCESS_KEY` | - | Default secret access key |
| `SIGNATURE_VERIFICATION_ENABLED` | `true` | Enable AWS Signature V4 verification |
| `USE_CHAIN_GATEWAY` | `false` | Enable chain gateway mode |

### TLS Settings

| Variable | Default | Description |
|----------|---------|-------------|
| `TLS_ENABLED` | `false` | Enable TLS for HTTP server |
| `TLS_PORT` | `8443` | TLS server port |
| `TLS_CERT_PATH` | `auto` | Path to TLS certificate (`auto` for self-signed) |
| `TLS_KEY_PATH` | `auto` | Path to TLS private key (`auto` for self-signed) |
| `TLS_HOSTNAME` | (empty) | TLS hostname for auto-generated certificates |

### gRPC-Web Admin API

| Variable | Default | Description |
|----------|---------|-------------|
| `GRPC_WEB_PORT` | `9090` | gRPC-Web admin server port |
| `GRPC_WEB_BIND_ADDR` | `127.0.0.1` | gRPC-Web bind address |

### IAM Authorization

| Variable | Default | Description |
|----------|---------|-------------|
| `AUTHORIZATION_ENABLED` | `false` | Enable IAM policy-based access control |
| `AUTHORIZATION_ROOT_ACCESS_KEYS` | (empty) | Comma-separated root access keys that bypass authorization |
| `AUTHORIZATION_DEFAULT_ACCESS_KEY_ID` | (empty) | Default access key when signature verification is disabled |
| `AUTHORIZATION_FAILURE_MODE` | `strict` | Authorization failure mode (`strict` = deny, `permissive` = allow) |
| `AUTHORIZATION_CACHE_TTL_SECONDS` | `300` | Policy evaluation cache TTL |
| `AUTHORIZATION_CACHE_MAX_SIZE` | `1000` | Maximum cached policy evaluations |
| `TEST_MODE` | `false` | Allow unauthenticated access in test mode |

### CloudTrail Audit Logging

| Variable | Default | Description |
|----------|---------|-------------|
| `VS_AUDIT_ENABLED` | `false` | Enable CloudTrail-compatible audit logging |

### Service Enablement

All services can be enabled/disabled individually via environment variables. Set to `false` to disable a service.

#### Required Services (default: `true`)

| Variable | Default | Service |
|----------|---------|---------|
| `ACM_ENABLED` | `true` | AWS Certificate Manager |
| `APIGATEWAY_ENABLED` | `true` | API Gateway |
| `CLOUDTRAIL_ENABLED` | `true` | CloudTrail |
| `CLOUDWATCH_ENABLED` | `true` | CloudWatch |
| `LOGS_ENABLED` | `true` | CloudWatch Logs |
| `COGNITO_ENABLED` | `true` | Cognito IDP |
| `COGNITO_IDENTITY_ENABLED` | `true` | Cognito Identity |
| `DYNAMODB_ENABLED` | `true` | DynamoDB |
| `EVENTS_ENABLED` | `true` | EventBridge |
| `KINESIS_ENABLED` | `true` | Kinesis |
| `KMS_ENABLED` | `true` | Key Management Service |
| `LAMBDA_ENABLED` | `true` | Lambda |
| `S3_ENABLED` | `true` | Simple Storage Service |
| `SCHEDULER_ENABLED` | `true` | EventBridge Scheduler |
| `SECRETSMANAGER_ENABLED` | `true` | Secrets Manager |
| `SESV2_ENABLED` | `true` | Simple Email Service v2 |
| `STEPFUNCTIONS_ENABLED` | `true` | Step Functions |
| `SNS_ENABLED` | `true` | Simple Notification Service |
| `SQS_ENABLED` | `true` | Simple Queue Service |
| `SSM_ENABLED` | `true` | Systems Manager Parameter Store |
| `STS_ENABLED` | `true` | Security Token Service |
| `IAM_ENABLED` | `true` | Identity and Access Management |

#### Optional Services (default: `true`)

| Variable | Default | Service |
|----------|---------|---------|
| `ATHENA_ENABLED` | `true` | Athena |
| `APPSYNC_ENABLED` | `true` | AppSync |
| `CLOUDFRONT_ENABLED` | `true` | CloudFront |
| `NEPTUNE_ENABLED` | `true` | Neptune |
| `NEPTUNE_DATA_ENABLED` | `true` | Neptune Data (Gremlin/SPARQL) |
| `NEPTUNE_GRAPH_ENABLED` | `true` | Neptune Graph |
| `ROUTE53_ENABLED` | `true` | Route53 |
| `TIMESTREAM_QUERY_ENABLED` | `true` | Timestream Query |
| `TIMESTREAM_WRITE_ENABLED` | `true` | Timestream Write |
| `WAFV2_ENABLED` | `true` | WAFv2 |

### Route53 DNS Server

| Variable | Default | Description |
|----------|---------|-------------|
| `ROUTE53_DNS_ENABLED` | `false` | Enable Route53 DNS server |
| `ROUTE53_DNS_BIND_ADDR` | `127.0.0.1` | DNS server bind address |

### Lambda Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `DOCKER_HOST` | `unix:///var/run/docker.sock` | Docker daemon endpoint |

## Starting the Server

### Development Mode

```bash
SIGNATURE_VERIFICATION_ENABLED=false \
DATA_PATH=./data \
./vorpalstacks
```

### Minimal (SQS + SNS only)

```bash
SNS_ENABLED=true \
SQS_ENABLED=true \
IAM_ENABLED=false \
STS_ENABLED=false \
KMS_ENABLED=false \
S3_ENABLED=false \
DYNAMODB_ENABLED=false \
SECRETSMANAGER_ENABLED=false \
ACM_ENABLED=false \
CLOUDWATCH_ENABLED=false \
CLOUDTRAIL_ENABLED=false \
LOGS_ENABLED=false \
# ... disable all other services ...
SIGNATURE_VERIFICATION_ENABLED=false \
DATA_PATH=./data \
./vorpalstacks
```

### Production (with TLS and Authorization)

```bash
TLS_ENABLED=true \
AUTHORIZATION_ENABLED=true \
AUTHORIZATION_ROOT_ACCESS_KEYS=AKIAIOSFODNN7EXAMPLE \
SIGNATURE_VERIFICATION_ENABLED=true \
DATA_PATH=/var/lib/vorpalstacks \
./vorpalstacks
```

### Test Mode

```bash
SIGNATURE_VERIFICATION_ENABLED=false \
TEST_MODE=true \
DATA_PATH=./data-test \
./vorpalstacks
```

## Data Storage Structure

```
DATA_PATH/
├── us-east-1/               # Region-specific storage
│   ├── pebble/              # PebbleDB for us-east-1
│   ├── code/                # Lambda function code
│   │   └── {function-name}/
│   │       └── code.zip
│   └── logs/                # CloudWatch Logs chunks
├── us-west-2/               # Region-specific storage
│   ├── pebble/
│   ├── code/
│   └── logs/
├── global/                  # Global services (IAM, Route53, CloudFront, STS)
│   └── pebble/
└── uploads/                 # S3 multipart uploads (temporary)
```

Each region has isolated storage including PebbleDB, Lambda code, and logs. Global services share dedicated storage.

## Docker Requirements

For Lambda functionality:

1. Docker must be installed and running
2. Lambda runtime images will be pulled automatically:
   - `public.ecr.aws/lambda/python:3.11`
   - `public.ecr.aws/lambda/node:18`
   - `public.ecr.aws/lambda/java:11`
   - etc.

3. Docker socket must be accessible at `DOCKER_HOST`

## AWS CLI Configuration

Configure AWS CLI to use Vorpalstacks:

```bash
export AWS_ENDPOINT_URL=http://localhost:8080
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_REGION=us-east-1
```

Or use per-command:

```bash
aws --endpoint-url=http://localhost:8080 \
    --region us-east-1 \
    sns list-topics
```

## Sample .env File

```env
# Core
PORT=8080
DATA_PATH=./data
AWS_REGION=us-east-1
AWS_ACCOUNT_ID=123456789012
SIGNATURE_VERIFICATION_ENABLED=false

# Required Services (all default to true, listed for reference)
# ACM_ENABLED=true
# APIGATEWAY_ENABLED=true
# CLOUDTRAIL_ENABLED=true
# CLOUDWATCH_ENABLED=true
# LOGS_ENABLED=true
# COGNITO_ENABLED=true
# COGNITO_IDENTITY_ENABLED=true
# DYNAMODB_ENABLED=true
# EVENTS_ENABLED=true
# KINESIS_ENABLED=true
# KMS_ENABLED=true
# LAMBDA_ENABLED=true
# S3_ENABLED=true
# SCHEDULER_ENABLED=true
# SECRETSMANAGER_ENABLED=true
# SESV2_ENABLED=true
# STEPFUNCTIONS_ENABLED=true
# SNS_ENABLED=true
# SQS_ENABLED=true
# SSM_ENABLED=true
# STS_ENABLED=true
# IAM_ENABLED=true

# Optional Services (all default to true, listed for reference)
# ATHENA_ENABLED=true
# APPSYNC_ENABLED=true
# CLOUDFRONT_ENABLED=true
# NEPTUNE_ENABLED=true
# NEPTUNE_DATA_ENABLED=true
# NEPTUNE_GRAPH_ENABLED=true
# ROUTE53_ENABLED=true
# TIMESTREAM_QUERY_ENABLED=true
# TIMESTREAM_WRITE_ENABLED=true
# WAFV2_ENABLED=true

# Route53 DNS (default: false)
# ROUTE53_DNS_ENABLED=false
# ROUTE53_DNS_BIND_ADDR=127.0.0.1

# gRPC-Web Admin
GRPC_WEB_PORT=9090

# Lambda
DOCKER_HOST=unix:///var/run/docker.sock
```

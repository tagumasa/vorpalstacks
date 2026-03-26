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
| `METADATA_PATH` | (empty) | Path to Smithy API metadata (loaded into Pebble DB at startup) |
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

The following services can be enabled/disabled individually. Set to `false` to disable:

| Variable | Default | Service |
|----------|---------|---------|
| `SNS_ENABLED` | `true` | Simple Notification Service |
| `SQS_ENABLED` | `true` | Simple Queue Service |
| `LAMBDA_ENABLED` | `true` | Lambda |
| `EVENTS_ENABLED` | `true` | EventBridge |
| `STEPFUNCTIONS_ENABLED` | `true` | Step Functions |
| `APIGATEWAY_ENABLED` | `true` | API Gateway |
| `SSM_ENABLED` | `true` | Systems Manager |
| `LOGS_ENABLED` | `true` | CloudWatch Logs |
| `COGNITO_ENABLED` | `true` | Cognito IDP |
| `COGNITO_IDENTITY_ENABLED` | `true` | Cognito Identity |
| `SCHEDULER_ENABLED` | `true` | EventBridge Scheduler |
| `KINESIS_ENABLED` | `true` | Kinesis |
| `CLOUDTRAIL_ENABLED` | `true` | CloudTrail |
| `SESV2_ENABLED` | `true` | Simple Email Service v2 |
| `TIMESTREAM_WRITE_ENABLED` | `true` | Timestream Write |
| `TIMESTREAM_QUERY_ENABLED` | `true` | Timestream Query |
| `ATHENA_ENABLED` | `true` | Athena |
| `ROUTE53_DNS_ENABLED` | `false` | Route53 DNS Server |

The following services are **always enabled** (no toggle): IAM, STS, KMS, S3, Route53, CloudFront, ACM, CloudWatch, DynamoDB, SecretsManager, WAF, WAFv2

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
DATA_PATH=./tmp/data \
./vorpalstacks
```

### Minimal (SQS + SNS)

```bash
SNS_ENABLED=true \
SQS_ENABLED=true \
SIGNATURE_VERIFICATION_ENABLED=false \
DATA_PATH=./tmp/data \
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

# Services
SNS_ENABLED=true
SQS_ENABLED=true
LAMBDA_ENABLED=true
EVENTS_ENABLED=true
STEPFUNCTIONS_ENABLED=true
APIGATEWAY_ENABLED=true
LOGS_ENABLED=true
COGNITO_ENABLED=true
COGNITO_IDENTITY_ENABLED=true
SCHEDULER_ENABLED=true
KINESIS_ENABLED=true
CLOUDTRAIL_ENABLED=true
SESV2_ENABLED=true
TIMESTREAM_WRITE_ENABLED=true
TIMESTREAM_QUERY_ENABLED=true
ATHENA_ENABLED=true
SSM_ENABLED=true

# gRPC-Web Admin
GRPC_WEB_PORT=9090

# Lambda
DOCKER_HOST=unix:///var/run/docker.sock
```

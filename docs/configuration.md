# Configuration

This document describes configuration options for Vorpalstacks.

## Environment Variables

### Core Settings

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `50080` | HTTP server port |
| `DATA_PATH` | `./data` | Path for persistent data storage |
| `AWS_REGION` | `us-east-1` | Default region |
| `AWS_ACCOUNT_ID` | `000000000000` | AWS account ID |
| `AWS_ACCESS_KEY_ID` | - | Default access key ID |
| `AWS_SECRET_ACCESS_KEY` | - | Default secret access key |
| `SIGNATURE_VERIFICATION_ENABLED` | `true` | Enable AWS Signature V4 verification |
| `USE_CHAIN_GATEWAY` | `false` | Enable chain gateway mode |
| `BIND_MODE` | `all` | Bind mode: `all` (0.0.0.0), `localhost` (127.0.0.1), or `interface` (custom IP) |
| `BIND_INTERFACE` | (empty) | Custom bind IP address (used only when `BIND_MODE=interface`) |

### TLS Settings

| Variable | Default | Description |
|----------|---------|-------------|
| `TLS_ENABLED` | `false` | Enable TLS for HTTP server |
| `TLS_PORT` | `50443` | TLS server port |
| `TLS_CERT_PATH` | `auto` | Path to TLS certificate (`auto` for self-signed) |
| `TLS_KEY_PATH` | `auto` | Path to TLS private key (`auto` for self-signed) |
| `TLS_HOSTNAME` | (empty) | TLS hostname for auto-generated certificates |

### gRPC-Web Admin API

| Variable | Default | Description |
|----------|---------|-------------|
| `GRPC_WEB_PORT` | `50090` | gRPC-Web admin server port (also serves web console) |

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

### Lambda Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `DOCKER_HOST` | `unix:///var/run/docker.sock` | Docker daemon endpoint |

## Port Architecture

All port defaults are defined in `internal/common/serviceports/ports.go`. Changing a constant migrates that port everywhere (server, CLI, SDK tests).

### Fixed Listener Ports

| Port | Env Var | Description |
|------|---------|-------------|
| `50080` | `PORT` | Main HTTP server (all AWS API endpoints) |
| `50090` | `GRPC_WEB_PORT` | gRPC-Web admin API + web console |
| `50443` | `TLS_PORT` | HTTPS server |
| `50088` | - | Route53 DNS server |
| `50089` | - | Route53 health check target |

### Service Endpoint Ports

Service endpoints (S3 Website, API Gateway, Cognito, CloudFront, Lambda URL, AppSync, Neptune) default to **FQDN mode**: all traffic is routed through the main HTTP port (`50080`) using `Host` header matching. No individual listener is bound.

When switched to **Individual mode** (via `vstacks config set ports.<service>.mode individual`), each resource gets its own port allocated from the dynamic range.

| Default Port | Config Key | Service |
|-------------|------------|---------|
| `50101` | `ports.s3_website` | S3 Website |
| `50102` | `ports.apigateway` | API Gateway |
| `50103` | `ports.cognito_hosted` | Cognito Hosted UI |
| `50104` | `ports.cloudfront` | CloudFront |
| `50105` | `ports.lambda_url` | Lambda Function URL |
| `50106` | `ports.appsync_events` | AppSync Events |
| `50107` | `ports.neptune` | Neptune DB (first cluster) |

### Dynamic Port Range

When services run in Individual mode, per-resource ports are allocated from:

| Config Key | Default |
|------------|---------|
| `ports.dynamic_range_start` | `50200` |
| `ports.dynamic_range_end` | `50400` |

## Starting the Server

### Development Mode

```bash
SIGNATURE_VERIFICATION_ENABLED=false \
DATA_PATH=./data \
./vorpalstacks
```

The admin console is available at `http://localhost:50090/webconsole/`. AWS API endpoints are on port `50080`.

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
export AWS_ENDPOINT_URL=http://localhost:50080
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_REGION=us-east-1
```

Or use per-command:

```bash
aws --endpoint-url=http://localhost:50080 \
    --region us-east-1 \
    sns list-topics
```

## Admin Config System

In addition to environment variables, runtime configuration can be managed through the **admin config system** backed by PebbleDB (`app_config` bucket). Settings are readable via the gRPC-Web Admin API and the `vstacks` CLI.

Priority order: **Store (persistent) > Environment variable > Default**

### app_config Keys

#### Server

| Key | Type | Default | Env Var | Description |
|-----|------|---------|---------|-------------|
| `server.port` | PORT | `50080` | `PORT` | Main HTTP server port |
| `server.grpc_web_port` | PORT | `50090` | `GRPC_WEB_PORT` | gRPC-Web admin port |
| `server.bind_mode` | STRING | `all` | `BIND_MODE` | Bind mode: `all`, `localhost`, or `interface` |
| `server.bind_interface` | STRING | (empty) | `BIND_INTERFACE` | Custom bind IP (bind_mode=interface only) |
| `server.tls_enabled` | BOOL | `false` | `TLS_ENABLED` | Enable TLS for HTTPS server |
| `server.tls_port` | PORT | `50443` | `TLS_PORT` | HTTPS server port |
| `server.tls_cert_path` | STRING | `auto` | `TLS_CERT_PATH` | TLS certificate path |
| `server.tls_key_path` | STRING | `auto` | `TLS_KEY_PATH` | TLS private key path |

#### AWS Identity (read-only)

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `aws.account_id` | STRING | (empty) | AWS account ID |
| `aws.region` | STRING | `us-east-1` | Default AWS region |

#### Storage

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `storage.data_path` | STRING | `./data` | Data storage path |
| `storage.metadata_path` | STRING | (empty) | Metadata path (defaults to data_path) |

#### Features

| Key | Type | Default | Env Var | Description |
|-----|------|---------|---------|-------------|
| `features.audit_enabled` | BOOL | `false` | `VS_AUDIT_ENABLED` | CloudTrail audit logging |
| `features.signature_verification` | BOOL | `true` | `SIGNATURE_VERIFICATION_ENABLED` | AWS signature verification |
| `features.route53_dns` | BOOL | `false` | `ROUTE53_DNS_ENABLED` | Route53 DNS server |

#### Endpoints

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `endpoints.base_url` | URL | `http://localhost:50080` | Base URL for generated endpoints |
| `endpoints.s3_website_suffix` | STRING | `s3-website.{region}.amazonaws.com` | S3 website domain suffix |
| `endpoints.cognito_suffix` | STRING | `auth.{region}.amazoncognito.com` | Cognito hosted UI suffix |
| `endpoints.apigateway_suffix` | STRING | `execute-api.{region}.amazonaws.com` | API Gateway suffix |

#### Ports

| Key | Type | Default | Env Var | Description |
|-----|------|---------|---------|-------------|
| `ports.s3_website` | PORT | `50101` | `VS_PORT_S3_WEBSITE` | S3 website default port |
| `ports.apigateway` | PORT | `50102` | `VS_PORT_APIGW` | API Gateway invoke URL port |
| `ports.cognito_hosted` | PORT | `50103` | `VS_PORT_COGNITO_HOSTED` | Cognito Hosted UI port |
| `ports.cloudfront` | PORT | `50104` | `VS_PORT_CLOUDFRONT` | CloudFront distribution port |
| `ports.lambda_url` | PORT | `50105` | `VS_PORT_LAMBDA_URL` | Lambda Function URL port |
| `ports.appsync_events` | PORT | `50106` | `VS_PORT_APPSYNC_EVENTS` | AppSync Events port |
| `ports.neptune` | PORT | `50107` | `VS_PORT_NEPTUNE` | Neptune DB cluster default port |
| `ports.route53_dns` | PORT | `50088` | `VS_PORT_ROUTE53_DNS` | Route53 DNS server port |
| `ports.route53_healthcheck` | PORT | `50089` | `VS_PORT_ROUTE53_HEALTHCHECK` | Route53 health check endpoint port |
| `ports.dynamic_range_start` | INT | `50200` | `VS_PORT_DYNAMIC_START` | Dynamic port range start |
| `ports.dynamic_range_end` | INT | `50400` | `VS_PORT_DYNAMIC_END` | Dynamic port range end |

#### Port Modes

Each service port supports a mode setting:

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `ports.s3_website.mode` | STRING | `fqdn` | S3 Website port mode (`fqdn` or `individual`) |
| `ports.apigateway.mode` | STRING | `fqdn` | API Gateway port mode |
| `ports.cognito_hosted.mode` | STRING | `fqdn` | Cognito Hosted UI port mode |
| `ports.cloudfront.mode` | STRING | `fqdn` | CloudFront port mode |
| `ports.lambda_url.mode` | STRING | `fqdn` | Lambda URL port mode |
| `ports.appsync_events.mode` | STRING | `fqdn` | AppSync Events port mode |
| `ports.neptune.mode` | STRING | `fqdn` | Neptune port mode |

- **`fqdn`**: All resources share the main HTTP port, routed by `Host` header. No additional listener.
- **`individual`**: Each resource gets its own port from the dynamic range.

#### HTTP / CORS

| Key | Type | Default | Env Var | Description |
|-----|------|---------|---------|-------------|
| `http.cors_allowed_origins` | STRING | `*` | `VS_CORS_ALLOWED_ORIGINS` | Comma-separated allowed origins (`*` = all) |
| `http.cors_allowed_methods` | STRING | `GET, POST, PUT, DELETE, OPTIONS, HEAD` | `VS_CORS_ALLOWED_METHODS` | Allowed CORS methods |
| `http.cors_allowed_headers` | STRING | `Authorization, Content-Type, X-Amz-Target, X-Amz-Date, X-Amz-Content-Sha256` | `VS_CORS_ALLOWED_HEADERS` | Allowed CORS request headers |
| `http.cors_expose_headers` | STRING | `x-amzn-RequestId, x-amzn-ErrorType, x-amzn-ErrorMessage` | `VS_CORS_EXPOSE_HEADERS` | CORS headers exposed to the client |

When `http.cors_allowed_origins` contains a comma-separated list (e.g. `http://localhost:3000,https://app.example.com`), the server matches the request `Origin` header against the list and responds with only the matching single origin. If no match is found, the `Access-Control-Allow-Origin` header is omitted entirely. The `Vary: Origin` header is added automatically for non-wildcard configurations.

## vstacks CLI

The `vstacks` command-line tool provides server control, IAM management, configuration, service mode control, and backup operations. It communicates with the server via gRPC-Web (server control) or reads PebbleDB directly (IAM, config, service, backup).

```
vstacks [options] <group> <command> [args]
```

### Options

| Option | Default | Description |
|--------|---------|-------------|
| `-data <path>` | `./data` | Path to data directory |
| `-account-id <id>` | `123456789012` | AWS Account ID |
| `-endpoint <url>` | `http://127.0.0.1:50090` | gRPC-Web admin endpoint |
| `-http-endpoint <url>` | `http://127.0.0.1:50080` | HTTP server endpoint |

### server — Server Control

| Command | Description |
|---------|-------------|
| `vstacks server status` | Check if the server is running (health endpoint) |
| `vstacks server stop` | Trigger graceful shutdown via gRPC-Web |

### iam — IAM User Management

| Command | Description |
|---------|-------------|
| `vstacks iam create-user -user <name> [-path <path>]` | Create an IAM user |
| `vstacks iam delete-user -user <name>` | Delete an IAM user |
| `vstacks iam list-users [-path-prefix <prefix>]` | List IAM users |
| `vstacks iam get-user -user <name>` | Get IAM user details |
| `vstacks iam create-access-key -user <name>` | Create an access key |
| `vstacks iam list-access-keys -user <name>` | List access keys |
| `vstacks iam delete-access-key -access-key-id <id>` | Delete an access key |
| `vstacks iam create-login-profile -user <name> -password <pw>` | Create a login profile |
| `vstacks iam delete-login-profile -user <name>` | Delete a login profile |

### config — Application Configuration (app_config)

| Command | Description |
|---------|-------------|
| `vstacks config get <key>` | Get a config value |
| `vstacks config set <key> <value>` | Set a config value |
| `vstacks config reset <key>` | Reset a config value to its default |
| `vstacks config list [-category <cat>]` | List all config entries |
| `vstacks config schema` | Show the full config schema |

### service — Service Mode Configuration (service_config)

| Command | Description |
|---------|-------------|
| `vstacks service get <name>` | Get a service's configuration |
| `vstacks service set-mode <name> -mode <MODE>` | Set service mode (`IMPLEMENTED`, `FALLBACK`, `ERROR_INJECTION`) |
| `vstacks service enable <name>` | Enable a service |
| `vstacks service disable <name>` | Disable a service |
| `vstacks service list` | List all services |

### backup — Data Backup and Restore

| Command | Description |
|---------|-------------|
| `vstacks backup create [-output <path>]` | Create a zip backup of the data directory |
| `vstacks backup list` | List existing backups |
| `vstacks backup restore <file>` | Restore from a backup (current data archived) |
| `vstacks backup info <file>` | Show backup contents and size |

### Examples

```bash
# Check server status
vstacks server status

# View current CORS configuration
vstacks config get http.cors_allowed_origins

# Restrict CORS to specific origins
vstacks config set http.cors_allowed_origins "http://localhost:3000,https://app.example.com"

# Reset CORS to default (allow all)
vstacks config reset http.cors_allowed_origins

# Switch API Gateway to individual port mode
vstacks config set ports.apigateway.mode individual

# Set a service to fallback mode
vstacks service set-mode dynamodb -mode FALLBACK

# Create a backup
vstacks backup create -output /tmp/my-backup.zip

# Graceful shutdown
vstacks server stop
```

## Sample .env File

```env
# Core
PORT=50080
DATA_PATH=./data
AWS_REGION=us-east-1
AWS_ACCOUNT_ID=123456789012
SIGNATURE_VERIFICATION_ENABLED=false

# Bind mode (default: all = 0.0.0.0)
# BIND_MODE=all
# BIND_INTERFACE=

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

# gRPC-Web Admin
GRPC_WEB_PORT=50090

# TLS
# TLS_ENABLED=false
# TLS_PORT=50443

# CORS (override defaults; omitted = allow all origins)
# VS_CORS_ALLOWED_ORIGINS=*
# VS_CORS_ALLOWED_METHODS=GET, POST, PUT, DELETE, OPTIONS, HEAD
# VS_CORS_ALLOWED_HEADERS=Authorization, Content-Type, X-Amz-Target, X-Amz-Date, X-Amz-Content-Sha256
# VS_CORS_EXPOSE_HEADERS=x-amzn-RequestId, x-amzn-ErrorType, x-amzn-ErrorMessage

# Lambda
DOCKER_HOST=unix:///var/run/docker.sock
```

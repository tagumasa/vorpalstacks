# Lambda Runtime Notes

This document covers runtime behaviour and configuration specific to Lambda function execution in vorpalstacks.

## Resource-Based Policy Principal Validation

`AddPermission` validates the `Principal` parameter against a known service-principal allowlist (`validServicePrincipals` in `validators.go`). This is stricter than real AWS, which accepts any `*.amazonaws.com` suffix. Rationale: on an edge/on-premises platform without the full IAM service-principal registry, a bare suffix check would accept typos (e.g. `lamda.amazonaws.com`) and spoofed principals (e.g. `evil.amazonaws.com`).

Accepted principal formats:

| Format | Example | Mapping |
|--------|---------|---------|
| Wildcard | `*` | `"Principal": "*"` |
| IAM ARN | `arn:aws:iam::123456789012:role/MyRole` | `"Principal": {"AWS": "arn:..."}` |
| Known service | `lambda.amazonaws.com` | `"Principal": {"Service": "lambda.amazonaws.com"}` |
| Unknown `*.amazonaws.com` | `evil.amazonaws.com` | **rejected** |

To add a new service principal, append it to `validServicePrincipals` in `internal/services/aws/lambda/validators.go`.

## Container Execution

Lambda functions run in Docker containers. The runtime lifecycle is:

1. On first invocation, vorpalstacks pulls the runtime image (e.g. `public.ecr.aws/lambda/nodejs:22`) and creates a container.
2. The container persists across subsequent invocations of the same function version until explicitly deleted or the function is updated.
3. Container names follow the pattern `lambda-{function-name}-{version}`.

### Supported Runtimes

The following runtimes are supported for new function creation. EOL runtimes (e.g. `nodejs20.x`, `go1.x`, `python3.9`) are rejected at validation time:

| Language | Identifiers |
|----------|------------|
| Node.js | `nodejs24.x`, `nodejs22.x` |
| Python | `python3.14`, `python3.13`, `python3.12`, `python3.11`, `python3.10` |
| Java | `java25`, `java21`, `java17`, `java11`, `java8.al2` |
| .NET | `dotnet10`, `dotnet9`*, `dotnet8` |
| Ruby | `ruby4.0`, `ruby3.4`, `ruby3.3` |
| OS-only | `provided.al2023`, `provided.al2` |

\* `dotnet9` is container-image only (no `--runtime` CLI support); it can be set via `UpdateFunctionConfiguration` on Image-package functions.

Deprecated runtimes that still have Docker base images (e.g. `nodejs20.x`, `go1.x`) will execute existing functions but cannot be selected for new functions. Go is supported via the `provided.al2023` custom runtime.

## Host Endpoint Injection

Containers run on a Docker bridge network and cannot access `localhost` on the host. Vorpalstacks automatically injects two things into every Lambda container:

### ExtraHosts

```
host.docker.internal → host-gateway
```

This resolves to the Docker bridge gateway IP (typically `172.17.0.1`), allowing containers to reach services running on the host.

### AWS_ENDPOINT_URL Environment Variable

```
AWS_ENDPOINT_URL=http://host.docker.internal:{PORT}
```

This is injected automatically unless the function's environment already defines `AWS_ENDPOINT_URL`. The port reflects the `PORT` environment variable used to start vorpalstacks.

Most AWS SDKs respect `AWS_ENDPOINT_URL` and will route all service calls to vorpalstacks instead of real AWS endpoints. This means your Lambda code can use standard SDK clients without any endpoint configuration:

```javascript
const { DynamoDBClient } = require("@aws-sdk/client-dynamodb");
const client = new DynamoDBClient({}); // automatically uses AWS_ENDPOINT_URL
```

```python
import boto3
dynamodb = boto3.client("dynamodb")  # automatically uses AWS_ENDPOINT_URL
```

### Overriding the Endpoint

To point a specific service at a different endpoint, set the service-specific variable:

```bash
aws lambda update-function-configuration \
    --function-name my-func \
    --environment Variables="{DYNAMODB_ENDPOINT=http://custom:8000}"
```

To disable automatic injection entirely, set `AWS_ENDPOINT_URL` yourself:

```bash
aws lambda update-function-configuration \
    --function-name my-func \
    --environment Variables="{AWS_ENDPOINT_URL=https://real.aws.amazon.com}"
```

## API Gateway Integration

API Gateway invoke requests are served on a dedicated secondary listener (default port 50102) to avoid routing conflicts with S3 on the main server port. This means:

- The invoke URL returned by `CreateStage` includes the configured API Gateway port.
- SDK clients must use the endpoint URL from the stage response, not the default main server port.
- The secondary port is configurable via `VS_PORT_APIGW` environment variable.

### Example: Full API Gateway → Lambda → DynamoDB Flow

```bash
# 1. Create a Lambda function
aws lambda create-function \
    --function-name api-handler \
    --runtime nodejs22.x \
    --handler index.handler \
    --role arn:aws:iam::000000000000:role/lambda-role \
    --zip-file fileb://function.zip

# 2. Create an API Gateway REST API with Lambda proxy integration
API_ID=$(aws apigateway create-rest-api --name my-api --query 'id' --output text)
ROOT=$(aws apigateway get-resources --rest-api-id $API_ID --query 'items[0].id' --output text)
ITEMS=$(aws apigateway create-resource --rest-api-id $API_ID --parent-id $ROOT --path-part items --query 'id' --output text)

aws apigateway put-method --rest-api-id $API_ID --resource-id $ITEMS --http-method GET --authorization-type NONE
aws apigateway put-method --rest-api-id $API_ID --resource-id $ITEMS --http-method PUT --authorization-type NONE

LAMBDA_URI="arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:000000000000:function:api-handler/invocations"
aws apigateway put-integration --rest-api-id $API_ID --resource-id $ITEMS --http-method GET --type AWS_PROXY --integration-http-method POST --uri $LAMBDA_URI
aws apigateway put-integration --rest-api-id $API_ID --resource-id $ITEMS --http-method PUT --type AWS_PROXY --integration-http-method POST --uri $LAMBDA_URI

DEPLOY=$(aws apigateway create-deployment --rest-api-id $API_ID --query 'id' --output text)
aws apigateway create-stage --rest-api-id $API_ID --stage-name prod --deployment-id $DEPLOY

# 3. Invoke via the API Gateway port (50102 by default)
curl http://localhost:50102/restapis/$API_ID/prod/_user_request_/items
curl -X PUT http://localhost:50102/restapis/$API_ID/prod/_user_request_/items \
    -H "Content-Type: application/json" \
    -d '{"id":"item1","data":"hello"}'
```

## Container Lifecycle

| Operation | Behaviour |
|-----------|-----------|
| First invocation | Pull image, create container, start, invoke |
| Subsequent invocations | Reuse existing running container |
| `UpdateFunctionCode` | Stop and remove old container, new container created on next invocation |
| `DeleteFunction` | Remove container if running |
| Server restart | Containers are not restarted; new containers are created on next invocation |

## Secondary Listener Ports

Vorpalstacks uses dedicated ports for invoke-style endpoints:

| Service | Default Port | Env Variable | Purpose |
|---------|-------------|-------------|---------|
| Main server | 50080 | `PORT` | Management API, all CRUD operations |
| API Gateway invoke | 50102 | `VS_PORT_APIGW` | Runtime invoke (`_user_request_` paths) |
| Lambda Function URL | 50105 | `VS_PORT_LAMBDA_URL` | Direct function URL invocation |
| S3 Website | 50101 | `VS_PORT_S3_WEBSITE` | Static website hosting |
| CloudFront | 50104 | `VS_PORT_CLOUDFRONT` | Distribution edge |
| Cognito Hosted UI | 50103 | `VS_PORT_COGNITO_HOSTED` | Authentication UI |
| AppSync Events | 50106 | `VS_PORT_APPSYNC_EVENTS` | WebSocket + HTTP publish |

If a secondary port is set to the same value as the main server port, the listener is skipped and the main router handles requests instead.

## Troubleshooting

### Container name conflict

```
Conflict. The container name "/lambda-my-func-LATEST" is already in use
```

Remove the stale container:

```bash
docker rm -f lambda-my-func-LATEST
```

### Lambda cannot reach other services

Verify `AWS_ENDPOINT_URL` is set inside the container:

```bash
docker exec lambda-my-func-LATEST env | grep AWS_ENDPOINT_URL
```

If missing, the function's environment may be overriding it. Check with `aws lambda get-function-configuration --function-name my-func`.

### Image pull failures

Ensure Docker can reach `public.ecr.aws`. In air-gapped environments, pre-pull and tag images:

```bash
docker pull public.ecr.aws/lambda/nodejs:22
```

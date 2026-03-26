# Service Integration

This document describes how services communicate and integrate with each other.

## Overview

Vorpalstacks services can trigger and invoke each other, mirroring AWS service integration patterns. IAM Authorization provides cross-cutting policy enforcement for all service integrations.

## Integration Matrix

| Source | Target | Mechanism |
|--------|--------|-----------|
| EventBridge | Lambda | Direct invocation |
| EventBridge | SQS | SendMessage |
| EventBridge | SNS | Publish |
| EventBridge | Step Functions | StartExecution |
| EventBridge | CloudWatch Logs | PutLogEvents |
| Scheduler | Lambda | Direct invocation |
| Scheduler | SQS | SendMessage |
| Scheduler | SNS | Publish |
| SNS | Lambda | Direct invocation |
| SNS | SQS | SendMessage |
| SNS | HTTP | HTTP POST |
| API Gateway | Lambda | Direct invocation |
| API Gateway | SQS | SendMessage |
| API Gateway | SNS | Publish |
| Step Functions | Lambda | Direct invocation |
| Step Functions | SQS | SendMessage |
| Step Functions | SNS | Publish |
| Lambda | CloudWatch Logs | Log streaming |
| KMS | SSM | Default encryption key |
| KMS | DynamoDB | Encryption at rest |
| KMS | S3 | Server-side encryption |

## EventBridge → Lambda

When an event matches a rule pattern, EventBridge invokes the target Lambda function.

### Configuration

```bash
aws events put-rule \
    --name my-rule \
    --event-pattern '{"source": ["my-app"]}'

aws events put-targets \
    --rule my-rule \
    --targets '{"Id": "1", "Arn": "arn:aws:lambda:us-east-1:123456789012:function:my-func"}'
```

### Input Transformation

```json
{
  "Id": "1",
  "Arn": "arn:aws:lambda:...",
  "Input": "{\"custom\": \"payload\"}",
  "InputPath": "$.detail"
}
```

## SNS → Lambda

SNS subscriptions can trigger Lambda functions.

### Configuration

```bash
aws sns subscribe \
    --topic-arn arn:aws:sns:us-east-1:123456789012:my-topic \
    --protocol lambda \
    --notification-endpoint arn:aws:lambda:us-east-1:123456789012:function:my-func
```

### Event Format

```json
{
  "Records": [{
    "EventSource": "aws:sns",
    "Sns": {
      "TopicArn": "arn:aws:sns:...",
      "Message": "message content",
      "MessageAttributes": {...}
    }
  }]
}
```

## API Gateway → Lambda

API Gateway can invoke Lambda for REST API backends.

### Integration Types

| Type | Description |
|------|-------------|
| AWS_PROXY | Direct Lambda invocation, event format |
| AWS | Lambda invocation with VTL templates |
| HTTP | Proxy to HTTP endpoint |
| HTTP_PROXY | Proxy to HTTP endpoint with request/response passthrough |
| MOCK | Fixed response for testing |

### Configuration

```bash
aws apigateway put-integration \
    --rest-api-id abc123 \
    --resource-id xyz789 \
    --http-method POST \
    --type AWS_PROXY \
    --uri arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:123456789012:function:my-func/invocations
```

## Step Functions → Lambda

Task states can invoke Lambda functions.

```json
{
  "Type": "Task",
  "Resource": "arn:aws:lambda:us-east-1:123456789012:function:my-func",
  "Next": "NextState"
}
```

## Scheduler → Lambda/SQS/SNS

EventBridge Scheduler supports timed invocation of Lambda, SQS, and SNS targets.

```bash
aws scheduler create-schedule \
    --name my-schedule \
    --schedule-expression "rate(5 minutes)" \
    --target '{"Arn": "arn:aws:lambda:...:function:my-func", "RoleArn": "arn:aws:iam::...:role/..."}'
```

## Lambda → CloudWatch Logs

Lambda automatically streams logs to CloudWatch Logs.

- Log group: `/aws/lambda/{function-name}` (auto-created)
- Log stream: `YYYY/MM/DD/[$LATEST]{uuid}` (one per invocation)

## SQS → Lambda (Event Source Mapping)

Lambda can poll SQS queues via event source mappings.

```bash
aws lambda create-event-source-mapping \
    --function-name my-func \
    --batch-size 10 \
    --event-source-arn arn:aws:sqs:us-east-1:123456789012:my-queue
```

## Cross-Cutting Features

### IAM Authorization

When `AUTHORIZATION_ENABLED=true`, all service-to-service integrations are subject to IAM policy evaluation. The dispatcher extracts the access key from the request, gathers attached policies, and evaluates them before routing to the service handler.

### CloudTrail Audit Logging

When `VS_AUDIT_ENABLED=true`, all API operations (including cross-service calls) are recorded as CloudTrail events with Loki-style indexing for efficient queries.

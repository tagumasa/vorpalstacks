# Implemented Services

**Last Updated**: 2026-03-26
**Total Services**: 30 (28 AWS services + 2 internal admin services)
**SDK Tests**: 594/594 (100%)

---

## AWS Services

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

### Service Scope

| Scope | Services |
|-------|----------|
| Global | IAM, STS, Route53, CloudFront |
| Regional | All others |

### Protocol

| Protocol | Services |
|----------|----------|
| AWS JSON 1.1 | IAM, API Gateway, SESv2, DynamoDB, Lambda, Cognito IDP, KMS, EventBridge, Kinesis, Athena, WAF, CloudTrail, Step Functions, SNS, CloudWatch Logs, Route53, Cognito Identity, Timestream Query, SecretsManager, ACM, CloudWatch Metrics, Scheduler, SSM, STS, Timestream Write, WAFv2 |
| REST-XML | S3, CloudFront |
| AWS Query | SQS |
| AWS JSON 1.0 | DynamoDB |

### Service Integration

| Source | Target | Description |
|--------|--------|-------------|
| EventBridge | Lambda, SQS, SNS, Step Functions, CloudWatch Logs | Event-driven invocation |
| Scheduler | Lambda, SQS, SNS | Scheduled invocation |
| SNS | Lambda, SQS | Pub/Sub fanout |
| Step Functions | Lambda, SQS, SNS | Workflow orchestration |
| API Gateway | Lambda | HTTP-to-Lambda proxy |
| Lambda | CloudWatch Logs | Automatic log streaming |
| KMS | SSM, DynamoDB, S3 | Encryption key provider |

### Cross-Cutting Features

- **IAM Authorization**: Policy-based access control (env: `AUTHORIZATION_ENABLED`)
- **CloudTrail Audit Logging**: API operation recording (env: `VS_AUDIT_ENABLED`)
- **gRPC-Web Admin API**: Connect-RPC admin interface on port 9090 (env: `GRPC_WEB_PORT`)

---

**Source**: Handler registration counts from service.go files

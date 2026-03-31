# Implemented Services

**Last Updated**: 2026-03-31
**Total**: 30 service APIs (26 AWS services)
**SDK Tests**: 890/890 (100%)

---

## Coverage Tiers

| Tier | Description |
|------|-------------|
| **Full** | All practical operations supported |
| **Broad** | Core workflows supported; advanced features remaining |
| **Selective** | Key features supported; see notes for limitations |

## AWS Services

| Service | Coverage | Notes |
|---------|----------|-------|
| ACM | Full | |
| API Gateway | Broad | No client certificates, documentation, or SDK generation |
| Athena | Broad | No capacity reservations or notebook sessions |
| CloudFront | Selective | No actual edge traffic distribution |
| CloudTrail | Broad | No event data stores or SQL queries |
| CloudWatch Logs | Selective | No Logs Insights queries or export |
| CloudWatch Metrics | Broad | No metric streams or anomaly detection |
| Cognito IDP | Selective | No external IdP, no hosted UI |
| Cognito Identity | Selective | Basic identity pool support |
| DynamoDB | Full | |
| EventBridge | Broad | No global endpoints or partner event sources |
| IAM | Broad | No organisations integration |
| Kinesis | Full | |
| KMS | Full | |
| Lambda | Broad | No durable functions or code signing |
| Route53 | Selective | DNS record management only |
| Neptune | Full | Property graph + RDF, openCypher/Gremlin, bulk loader, management API |
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
| WAF | Selective | No managed rule groups or logging configuration |
| WAFv2 | Broad | |

### Service Scope

| Scope | Services |
|-------|----------|
| Global | IAM, STS, Route53, CloudFront |
| Regional | All others |

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

**Source**: Handler registration counts from service.go files. For detailed API gap analysis, see [plans/](../plans/).

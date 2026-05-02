# Implemented Services

**Last Updated**: 2026-04-21
**Total**: 32 AWS services
**SDK Tests**: 2262 (2216 SDK + 29 integration + 17 WebSocket)

---

## Coverage Tiers

| Tier | Description |
|------|-------------|
| **Full** | All practical operations supported |
| **Broad** | Core workflows supported; advanced features remaining |
| **Selective** | Key features supported; see notes for limitations |

## AWS Services

### Required Services (22, default: enabled)

| Service | Coverage | Notes |
|---------|----------|-------|
| ACM | Full | |
| API Gateway | Broad | No client certificates, documentation, or SDK generation |
| CloudTrail | Broad | No event data stores or SQL queries |
| CloudWatch Metrics | Broad | No metric streams or anomaly detection |
| CloudWatch Logs | Selective | No Logs Insights queries or export |
| Cognito IDP | Selective | No external IdP |
| Cognito Identity | Selective | Basic identity pool support |
| DynamoDB | Full | |
| EventBridge | Broad | No global endpoints or partner event sources |
| IAM | Broad | No organisations integration |
| Kinesis | Full | |
| KMS | Full | |
| Lambda | Broad | No durable functions or code signing |
| S3 | Broad | No analytics, inventory, or S3 Express |
| Scheduler | Full | |
| Secrets Manager | Full | |
| SESv2 | Broad | No deliverability testing or multi-tenancy |
| SFN (Step Functions) | Full | |
| SNS | Broad | SMS sending not supported |
| SQS | Full | |
| SSM | Selective | Parameter Store only |
| STS | Full | |

### Optional Services (10, default: enabled)

| Service | Coverage | Notes |
|---------|----------|-------|
| Athena | Broad | No capacity reservations or notebook sessions |
| AppSync | Broad | GraphQL API with VTL resolvers, real-time subscriptions |
| CloudFront | Selective | No actual edge traffic distribution |
| Neptune | Full | Property graph + RDF, openCypher/Gremlin, bulk loader, management API |
| NeptuneData | Broad | Gremlin/SPARQL query endpoint |
| NeptuneGraph | Broad | Graph engine with graph/SPARQL/neptune-analytics APIs |
| Route53 | Selective | DNS record management only |
| Timestream Query | Broad | SQL query engine |
| Timestream Write | Broad | Time-series data ingestion |
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
| API Gateway | Lambda, SQS, SNS | HTTP-to-service proxy |
| Lambda | CloudWatch Logs | Automatic log streaming |
| Lambda | SQS | Event source mapping (polling) |
| KMS | SSM, DynamoDB, S3 | Encryption key provider |

### Cross-Cutting Features

- **IAM Authorization**: Policy-based access control (env: `AUTHORIZATION_ENABLED`)
- **CloudTrail Audit Logging**: API operation recording (env: `VS_AUDIT_ENABLED`)
- **gRPC-Web Admin API**: Connect-RPC admin interface on port 9090 (env: `GRPC_WEB_PORT`)

---

**Source**: Handler registration counts from service.go files. For detailed API gap analysis, see [plans/](../plans/).

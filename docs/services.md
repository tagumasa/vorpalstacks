# Implemented Services

**Last Updated**: 2026-09-05
**Total**: 35 AWS services — single source of truth for the supported-service count, per the AWS SDK service classification (Timestream Write and Timestream Query are separate SDK services)
**SDK Tests**: 3,478 passed, 0 failed (3,396 SDK + 65 integration + 17 WebSocket)

---

## Coverage Tiers

| Tier | Description |
|------|-------------|
| **Full** | All practical operations supported |
| **Broad** | Core workflows supported; advanced features remaining |
| **Selective** | Key features supported; see notes for limitations |

## AWS Services

### Required Services (21, default: enabled)

| Service | Coverage | Notes |
|---------|----------|-------|
| ACM | Broad | No ACME protocol (19 ops) |
| API Gateway | Broad | No client certificates, documentation, or SDK generation |
| CloudWatch Metrics | Broad | No metric streams or anomaly detection |
| CloudWatch Logs | Selective | No Logs Insights queries or export |
| Cognito IDP | Selective | No external IdP |
| Cognito Identity | Selective | Basic identity pool support |
| DynamoDB | Broad | ION import/export not supported. Streams and Global Tables implemented with multi-active replication. |
| EventBridge | Broad | No global endpoints or partner event sources |
| IAM | Broad | No custom-policy simulator (`SimulatePrincipalPolicy` and `ListPoliciesGrantingServiceAccess` implemented) or organisations integration. GetHumanReadableSummary excluded (external LLM dependence). Delegation request APIs not implemented. |
| Kinesis | Full | |
| KMS | Full | |
| Lambda | Broad | No durable functions, code signing, capacity providers, recursive loop detection, function scaling, or managed runtime updates. |
| S3 | Broad | Bucket inventory/metrics configurations not implemented. No analytics/intelligent-tiering configurations (no storage-class tiering), object annotations, bucket ABAC, S3 Express, S3 Metadata tables, GetObjectTorrent, or WriteGetObjectResponse. Object Lock, CORS, lifecycle, SSE encryption fully enforced. |
| Scheduler | Full | Templated targets limited to platform-implemented services (Lambda, SQS, SNS, Kinesis, Step Functions, EventBridge). ECS and Firehose targets are accepted but fail at delivery until those services exist. SageMaker, CodeBuild, CodePipeline, and Inspector targets are permanently out of scope (services not implemented on this platform). |
| Secrets Manager | Full | ListTagsForResource is implemented beyond the 2017-10-17 API model (the operation does not exist in the model, so AWS SDKs never generate a client method for it; the platform operation serves raw-HTTP/console consumers). Managed external secret rotation members are configuration storage and echo only — the partner integration itself is external. |
| SESv2 | Broad | No deliverability testing, dedicated IP address management, import/export jobs, multi-region endpoints, tenant management, custom verification email templates, reputation management, or account pricing plans |
| SFN (Step Functions) | Full | |
| SNS | Broad | SMS sending, email/email-json delivery, and mobile push (application protocol) not supported. Platform application/endpoint CRUD is available but no push delivery. Subscription FilterPolicy and RawMessageDelivery are supported. |
| SQS | Broad | SSE-KMS encryption not applied (attributes accepted but messages are not encrypted). FIFO advanced attributes (DeduplicationScope, FifoThroughputLimit, RedriveAllowPolicy) accepted but not enforced. Per-account request-rate quotas are not enforced (AWS-account-tied rate limiting is out of scope; the RequestThrottled error shape exists for wire-contract completeness). |
| SSM | Selective | Parameter Store only |
| STS | Full | |

### Optional Services (14)

| Service | Coverage | Default | Notes |
|---------|----------|---------|-------|
| Athena | Broad | enabled | No capacity reservations or notebook sessions |
| AppSync | Broad | enabled | GraphQL API with VTL resolvers, real-time subscriptions |
| CloudFront | Broad | enabled | Origin proxy with cache behaviours, TTL edge cache, invalidation, CNAME aliases, continuous deployment policies, viewer TLS serving, and ViewerProtocolPolicy enforcement |
| CloudTrail | Broad | **disabled** | Audit logging. No event data stores or SQL queries |
| EC2 | Selective | enabled | Basic instance management |
| IoT Core | Broad | enabled | 272 operations: things, certificates, policies, rules engine (11 action types), jobs, shadows, device management |
| Neptune | Full | enabled | Property graph + RDF, openCypher/Gremlin, bulk loader, management API |
| NeptuneData | Broad | enabled | Gremlin/SPARQL query endpoint |
| NeptuneGraph | Broad | enabled | Graph engine with graph/SPARQL/neptune-analytics APIs |
| RDS Data | Full | **disabled** | MySQL-compatible SQL via vmysql engine (requires `RDS_MYSQL_ENABLED=true` or `ALL_SERVICES_ENABLED=true`) |
| Route53 | Selective | enabled | DNS record management only |
| Timestream Query | Broad | enabled | SQL query engine |
| Timestream Write | Broad | enabled | Time-series data ingestion |
| WAFv2 | Broad | enabled | Signature managed rule groups evaluate locally; the six data-dependent groups (IP reputation, anonymous IP, Bot Control, ATP, ACFP, Anti-DDoS) and the Known Bad Inputs ReactJS RCE rule never match locally (their inputs exist only inside AWS); Monetize payment settlement is not verified |

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
| S3 | Lambda, SQS, SNS, EventBridge | S3 event notifications |
| IoT Core | Lambda, SQS, SNS, Kinesis, DynamoDB, S3, CloudWatch, CloudWatch Logs, Step Functions, Timestream, Republish | IoT rule engine actions (11 types) |
| KMS | SSM, DynamoDB, S3 | Encryption key provider |

### Platform Behaviour Notes

Platform decisions and restrictions on behaviour AWS leaves unspecified. The service tables above carry feature availability only.

- **Athena — TEST_MODE**: query execution history is purged at startup.
- **CloudFront — viewer TLS serving**: SNI per distribution, from the attached ACM/IAM certificate.
- **Kinesis — SubscribeToShard heartbeat interval**: 15 s (provisional; AWS does not document the exact value).
- **Lambda — AddPermission Principal**: restricted to a known service-principal allowlist (see `validServicePrincipals` in `validators.go`); unrecognised `*.amazonaws.com` principals are rejected.
- **Secrets Manager — BatchGetSecretValue `MaxResults`**: AWS documents the requirement ("To use this parameter, you must also use the Filters parameter") but not the behaviour when it is violated; requests pairing `MaxResults` with `SecretIdList` are rejected with `InvalidParameterException` (400).
- **WAFv2 — GeoMatch/AsnMatch resolution**: embedded table derived from RIR delegated-extended allocation data and a RouteViews routing-table snapshot (regenerated by internal/tools/wafgeogen); country results follow registry allocations rather than a commercial geolocation database and can differ from AWS WAF for reassigned addresses.
- **WAFv2 — CAPTCHA/Challenge/Monetize token**: the aws-waf-token cookie is HMAC-signed locally; text/html clients receive a JavaScript interstitial whose proof-of-work challenge is exchanged at the reserved POST /awswaf/token endpoint, with the documented 405/202 interrupts and immunity times; Monetize validates its configuration and interrupts with the 402 price manifest, but payment settlement requires blockchain network access and is not verified.
- **WAFv2 — managed rule groups**: served from a local catalog of the fifteen documented groups; the nine signature groups (63 rules) evaluate against local statements derived from the published rule descriptions — AWS publishes no exact match patterns, so each local statement covers the documented examples plus the canonical public signatures of the same threat class; the six data-dependent groups (73 rules) and the Known Bad Inputs ReactJS RCE rule never match locally because their inputs (threat-intelligence feeds, device fingerprints, ML models, unpublished advisory patterns) exist only inside AWS.
- **WAFv2 — HeaderOrder component**: the wire order is preserved on HTTP/1.1 connections; HTTP/2 connections fall back to the header map's order.
- **AWS IoT — managed job templates**: the platform ships no AWS-provided managed-job-template catalogue (the catalogue content is AWS's copyrighted material), so `DescribeManagedJobTemplate` resolves every template name to `ResourceNotFoundException` and `ListManagedJobTemplates` returns an empty list.

### Cross-Cutting Features

- **IAM Authorization**: Policy-based access control (env: `AUTHORIZATION_ENABLED`)
- **CloudTrail Audit Logging**: API operation recording (env: `CLOUDTRAIL_ENABLED`)
- **gRPC-Web Admin API**: Connect-RPC admin interface on port 50090 (env: `GRPC_WEB_PORT`)

---

**Source**: Handler registration counts from service.go files.

# Implemented Services

**Last Updated**: 2026-09-06
**Total**: 35 AWS services — single source of truth for the supported-service count, per the AWS SDK service classification (Timestream Write and Timestream Query are separate SDK services)
**SDK Tests**: over 3,000 passing (Go SDK, cross-service integration, and WebSocket suites; exact counts live in `sdk-tests/README.md`)

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
| ACM | Broad | No ACME protocol |
| API Gateway | Broad | No client certificates, documentation, SDK generation, or VpcLink |
| CloudWatch Metrics | Broad | No metric streams or anomaly detection |
| CloudWatch Logs | Selective | No Logs Insights queries or export |
| Cognito IDP | Selective | No external IdP |
| Cognito Identity | Selective | Basic identity pool support |
| DynamoDB | Broad | No ION import/export; Streams and Global Tables implemented |
| EventBridge | Broad | No global endpoints or partner event sources |
| IAM | Broad | No policy-simulator family beyond `SimulatePrincipalPolicy` and `ListPoliciesGrantingServiceAccess`; no organisations integration, GetHumanReadableSummary, or delegation request APIs |
| Kinesis | Full | |
| KMS | Full | |
| Lambda | Broad | No durable functions, code signing, capacity providers, recursive loop detection, function scaling, managed runtime updates, or resource-policy public-access rejection |
| S3 | Broad | Inventory and metrics configurations implemented; no analytics/intelligent-tiering configurations, object annotations, bucket ABAC, S3 Express, S3 Metadata tables, GetObjectTorrent, or WriteGetObjectResponse; Object Lock, CORS, lifecycle, and SSE enforced |
| Scheduler | Full | Templated targets limited to platform-implemented services; no SageMaker, CodeBuild, CodePipeline, or Inspector targets |
| Secrets Manager | Full | No managed external rotation execution |
| SESv2 | Broad | No deliverability testing, dedicated IP address management, import/export jobs, multi-region endpoints, tenant management, custom verification email templates, reputation management, or account pricing plans |
| SFN (Step Functions) | Full | |
| SNS | Broad | No SMS, email/email-json, or mobile push (application protocol) delivery; platform application/endpoint CRUD, FilterPolicy, and RawMessageDelivery implemented |
| SQS | Broad | No SSE-KMS message encryption, FIFO advanced-attribute enforcement (DeduplicationScope, FifoThroughputLimit, RedriveAllowPolicy), or per-account request-rate quotas |
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
| IoT Core | Broad | enabled | Things, certificates, policies, rules engine, jobs, shadows, and device management |
| Neptune | Full | enabled | Property graph + RDF, openCypher/Gremlin, bulk loader, management API |
| NeptuneData | Broad | enabled | Gremlin/SPARQL query endpoint |
| NeptuneGraph | Broad | enabled | Graph engine with graph/SPARQL/neptune-analytics APIs |
| RDS Data | Full | **disabled** | MySQL-compatible SQL via vmysql engine |
| Route53 | Selective | enabled | DNS record management only |
| Timestream Query | Broad | enabled | SQL query engine |
| Timestream Write | Broad | enabled | Time-series data ingestion |
| WAFv2 | Broad | enabled | Signature managed rule groups implemented; no data-dependent groups (IP reputation, anonymous IP, Bot Control, ATP, ACFP, Anti-DDoS) or Known Bad Inputs ReactJS RCE rule; Monetize payment settlement unverified |

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

Platform behaviour detail and restrictions, including where AWS leaves behaviour unspecified. The service tables above carry feature availability only.

- **API Gateway — VPC_LINK connection type**: rejects at both integration create and the /connectionType replace path; a VPC_LINK integration would route through a VpcLink to a Network Load Balancer, which this platform does not provide.
- **Athena — TEST_MODE**: query execution history is purged at startup.
- **CloudFront — viewer TLS serving**: SNI per distribution, from the attached ACM/IAM certificate.
- **Cognito IDP — user-pool domains**: the four domain operations are implemented; domain entries resolve to the platform endpoint suffix (`<domain>.auth.<cognito_suffix>` with the region substituted) rather than AWS-hosted CloudFront domains, which cannot exist in an edge/on-premises deployment.
- **DynamoDB — Streams and Global Tables**: implemented with multi-active replication.
- **IAM — GetHumanReadableSummary**: excluded; the summary requires external LLM generation.
- **Kinesis — SubscribeToShard heartbeat interval**: 15 s (provisional; AWS does not document the exact value).
- **Lambda — AddPermission Principal**: restricted to a known service-principal allowlist (see `validServicePrincipals` in `validators.go`); unrecognised `*.amazonaws.com` principals are rejected.
- **Lambda — resource-policy operations**: PutResourcePolicy/GetResourcePolicy/DeleteResourcePolicy share the AddPermission policy substrate — one policy per function, so a Put replaces permissions added through AddPermission and both GetPolicy and GetResourcePolicy serve the same statements, with Put-path statements reproduced from their verbatim JSON. Mutations carry a policy revision (RevisionId precondition answered with PreconditionFailedException); the policy document is bounded by the 20,480-character quota, and DeleteResourcePolicy succeeds on a resource whose policy is already gone.
- **Lambda — resource-policy public-access rejection**: excluded; the platform accepts resource-policy documents that grant wildcard-principal access instead of answering PublicPolicyException. The Smithy model declares the error on AddPermission, PutResourcePolicy, and RemovePermission, and its message names the trigger loosely — a policy that "would grant public access to it, which isn't allowed" — but defines no boundary for public access; the Lambda documentation's own valid examples include Condition-scoped wildcard-principal Allow statements, and current AWS enforcement is tied to the account-level block-public-access setting, which the vendored models do not carry.
- **Lambda — RuntimeManagementConfig**: excluded; the configuration API drives a managed runtime-update mechanism (a runtime-version registry that pins and advances each function's runtime patch level). This platform has no update engine — `RuntimeVersionConfig` is configuration echo only — so the API surface would be inert.
- **Lambda — RecursionConfig**: excluded; its only enforced semantics is the recursion loop detector's behaviour (Terminate/AllowOnFail). This platform has no loop-detection engine, so the default Terminate would report terminating behaviour that never occurs.
- **Lambda — ScalingConfig**: excluded; the model scopes the API to Managed-Instances functions running on Lambda capacity providers, which this platform does not provide.
- **S3 — inventory report delivery**: reports deliver on daily/weekly UTC boundaries to the S3 destination as CSV (gzip), Parquet (snappy), and ORC (ZLIB), with manifest.json, manifest.checksum, and the Hive symlink; report files honour the configuration's SSE-S3/SSE-KMS encryption choice. The report columns IntelligentTieringAccessTier, ChecksumAlgorithm, and LifecycleExpirationDate are emitted empty (no single-tier substrate). Both configuration families are bounded by the 1,000-configuration limit with 100-item pagination.
- **S3 — metrics configurations**: per-filter CloudWatch request metrics in the AWS/S3 namespace; requests on both the object and bucket planes count into AllRequests plus their per-operation metric, and each minute window publishes a CloudWatch statistic set (sample count, sum, min, max, so Average carries the documented error rate and bytes-per-request semantics). A filter carrying an access-point ARN generates no datapoints (no access-point substrate).
- **S3 — analytics/intelligent-tiering configurations**: excluded; this single-tier platform has no storage-class transition substrate.
- **RDS Data — enablement**: requires `RDS_MYSQL_ENABLED=true` or `ALL_SERVICES_ENABLED=true`.
- **Scheduler — templated targets**: platform-implemented targets are Lambda, SQS, SNS, Kinesis, Step Functions, and EventBridge; ECS and Firehose targets are accepted by rule templates but delivery fails until those services exist on the platform; SageMaker, CodeBuild, CodePipeline, and Inspector targets are permanently out of scope (those services are not implemented on this platform).
- **Secrets Manager — ListTagsForResource and managed rotation members**: the operation does not exist in the 2017-10-17 model, so AWS SDKs never generate a client method for it; the platform operation serves raw-HTTP/console consumers. Managed external rotation members are configuration storage and echo only — the partner integration itself is external.
- **SQS — SSE-KMS and request throttling**: SSE-KMS attributes are accepted but messages are stored unencrypted; per-account request-rate quotas are not enforced, and the RequestThrottled error shape exists for wire-contract completeness only.
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

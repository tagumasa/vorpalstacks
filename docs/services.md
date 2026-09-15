# Implemented Services

**Last Updated**: 2026-09-12
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
| API Gateway | Broad | No client certificates, documentation parts, GetSdk, VpcLink, or domain name access associations |
| CloudWatch Metrics | Broad | No metric streams or anomaly detection |
| CloudWatch Logs | Selective | No Logs Insights queries or export |
| Cognito IDP | Selective | No external IdP; no Firehose log-delivery export |
| Cognito Identity | Selective | Identity pools only; no external-IdP role-mapping claims |
| DynamoDB | Broad | No ION import/export |
| EventBridge | Broad | No global endpoints or partner event sources |
| IAM | Broad | No policy-simulator family beyond `SimulatePrincipalPolicy` and `ListPoliciesGrantingServiceAccess`; no organisations integration, GetHumanReadableSummary, or delegation request APIs |
| Kinesis | Full | |
| KMS | Full | |
| Lambda | Broad | No durable functions, code signing, capacity providers, recursive loop detection, function scaling, managed runtime updates, or resource-policy public-access rejection |
| S3 | Broad | No analytics/intelligent-tiering configurations, object annotations, bucket ABAC, S3 Express, S3 Metadata tables, GetObjectTorrent, WriteGetObjectResponse, S3 Select (SelectObjectContent), or hardware-bound storage classes |
| Scheduler | Full | No SageMaker, CodeBuild, CodePipeline, or Inspector targets |
| Secrets Manager | Full | No managed external rotation execution |
| SESv2 | Broad | No deliverability testing, dedicated IP address management, import/export jobs, multi-region endpoints, tenant management, custom verification email templates, reputation management, or account pricing plans |
| SFN (Step Functions) | Full | No integrations with unimplemented services (batch, ECS, Glue, and the rest of the AWS SDK namespace) |
| SNS | Broad | No SMS, email/email-json, or mobile push (application protocol) delivery |
| SQS | Broad | No SSE-KMS message encryption, FIFO advanced-attribute enforcement (DeduplicationScope, FifoThroughputLimit, RedriveAllowPolicy), or per-account request-rate quotas |
| SSM | Selective | Parameter Store only |
| STS | Full | |

### Optional Services (14)

| Service | Coverage | Default | Notes |
|---------|----------|---------|-------|
| Athena | Broad | enabled | No capacity reservations or notebook sessions |
| AppSync | Broad | enabled | GraphQL APIs, VTL resolvers, real-time subscriptions; AMAZON_BEDROCK_RUNTIME data source type excluded; Events API OpenID Connect authorisation excluded; appsyncDomainName is synthesised, not provisioned |
| CloudFront | Broad | enabled | Origin proxy, cache behaviours, TTL edge cache, invalidation, CNAME aliases, continuous deployment policies, viewer TLS, ViewerProtocolPolicy |
| CloudTrail | Broad | **disabled** | No event data stores or SQL queries |
| EC2 | Selective | enabled | Instance management only |
| IoT Core | Broad | enabled | Things, certificates, policies, rules engine, jobs, shadows, and device management |
| Neptune | Full | enabled | Property graph + RDF, openCypher/Gremlin, bulk loader, management API |
| NeptuneData | Broad | enabled | Gremlin/SPARQL query endpoint |
| NeptuneGraph | Broad | enabled | Graph, SPARQL, and neptune-analytics APIs |
| RDS Data | Full | **disabled** | MySQL-compatible SQL |
| Route53 | Selective | enabled | DNS record management only |
| Timestream Query | Broad | enabled | SQL query engine |
| Timestream Write | Broad | enabled | Time-series data ingestion |
| WAFv2 | Broad | enabled | No data-dependent groups (IP reputation, anonymous IP, Bot Control, ATP, ACFP, Anti-DDoS) or Known Bad Inputs ReactJS RCE rule |

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
| Step Functions | Lambda, SQS, SNS, EventBridge, DynamoDB, Step Functions | Workflow orchestration |
| API Gateway | Lambda, SQS, SNS | HTTP-to-service proxy |
| Lambda | CloudWatch Logs | Automatic log streaming |
| Lambda | SQS | Event source mapping (polling) |
| S3 | Lambda, SQS, SNS, EventBridge | S3 event notifications |
| IoT Core | Lambda, SQS, SNS, Kinesis, DynamoDB, S3, CloudWatch, CloudWatch Logs, Step Functions, Timestream, Republish | IoT rule engine actions (11 types) |
| KMS | SSM, DynamoDB, S3 | Encryption key provider |

### Platform Behaviour Notes

Platform behaviour detail and restrictions, including where AWS leaves behaviour unspecified. The service tables above carry feature availability only.

- **API Gateway — VPC_LINK connection type**: rejects at both integration create and the /connectionType replace path; a VPC_LINK integration would route through a VpcLink to a Network Load Balancer, which this platform does not provide.
- **API Gateway — domain name access associations**: the four operations are absent because the feature exists only for VPC-endpoint-backed access — the access association source type admits VPCE as its only valid value and applies to PRIVATE custom domain names, and the platform provides no VPC-endpoint substrate (the same basis as the VpcLink exclusion).
- **API Gateway — GetSdk**: excluded with the documentation-parts and client-certificate families; SDK package generation has no consumer on this platform, while the static SDK-type metadata operations (GetSdkType, GetSdkTypes) are implemented.
- **AppSync — AMAZON_BEDROCK_RUNTIME data source type**: rejected at creation; the type names the AWS-side Bedrock LLM runtime, which this platform does not provide — admitting it would leave every resolver on such a data source failing at dispatch.
- **AppSync — Events API OpenID Connect authorisation**: excluded; verifying an OIDC token requires reaching the external identity provider that issued it, and external IdP integration is outside this platform's scope. Connections and operations presenting OPENID_CONNECT credentials are denied (fail-closed); the other Events authorisation modes — API key, IAM, Cognito user pools, Lambda authoriser — are verified.
- **AppSync — appsyncDomainName**: synthesised deterministically as `<domainName>.appsync-api.<region>.amazonaws.com`; AWS provisions a CloudFront distribution and returns its domain, which cannot exist in an edge/on-premises deployment — no distribution is created and no DNS record is served, so the value is metadata for client-side CNAME setup (the same basis as the Cognito user-pool domain note).
- **Athena — TEST_MODE**: query execution history is purged at startup.
- **CloudFront — viewer TLS serving**: SNI per distribution, from the attached ACM/IAM certificate.
- **Cognito IDP — user-pool domains**: the four domain operations are implemented; domain entries resolve to the platform endpoint suffix (`<domain>.auth.<cognito_suffix>` with the region substituted) rather than AWS-hosted CloudFront domains, which cannot exist in an edge/on-premises deployment.
- **Cognito IDP — imported users and SRP sign-in**: CSV-imported users carry an imported password hash and no SRP verifier; USER_SRP_AUTH rejects them with NotAuthorizedException until native credentials exist, while USER_PASSWORD_AUTH verifies the imported hash and transparently migrates the account to native bcrypt+SRP credentials. RESET_REQUIRED imports without a hash complete the migration through the NEW_PASSWORD_REQUIRED challenge or the forgot-password flow. AdminResetUserPassword deactivates imported credentials as well — every post-reset sign-in answers PasswordResetRequiredException until the forgot-password flow completes.
- **Cognito IDP — log delivery destinations**: SetLogDeliveryConfiguration accepts all three destination members (CloudWatch Logs, S3, Firehose) and both event sources (userAuthEvents, userNotification) with the AWS contract enforced — the log level is bound to the event source (userNotification at ERROR, userAuthEvents at INFO), each configuration names exactly one destination, and userNotification goes to CloudWatch Logs alone. userAuthEvents records export to the CloudWatch Logs destination and to the S3 destination as per-event JSON objects under `AWSLogs/<account>/cognito-idp/<region>/<poolID>/`; userNotification records (message-delivery notifications on code delivery) export to the CloudWatch Logs destination alone. The Firehose destination is accepted and stored with delivery pending the platform Firehose service (the service is unimplemented — a missing-substrate carryover; delivery lands together with that service).
- **Cognito IDP — machine-to-machine authorisation configuration**: the `ALLOW_CLIENT_TOKEN_AUTH` explicit authentication flow (required by GetClientToken) is accepted on the AWS API plane, where it must be the client's only flow; the admin console's generated proto predates the value and does not project it, so an M2M app client is configured through the AWS API plane alone.
- **Cognito Identity — role-mapping claim source**: role mappings read claims from linked platform user-pool ID tokens — Logins keys of the issuer form `cognito-idp.<region>.amazonaws.com/<poolID>`, validated by signature, issuer and token_use; the ID token derives its `cognito:roles` and `cognito:preferred_role` claims from the IAM roles of the user's groups (highest precedence nominates the preferred role). A login under any other provider name belongs to an external identity provider, which the platform does not implement: it carries no verifiable token, contributes no claims, and the mapping's AmbiguousRoleResolution governs (fail-closed under Deny). Principal-tag attribute maps feed session tags on the issued credentials — custom mappings take their values from the mapped ID-token claims, and UseDefaults applies the aud and sub defaults (app client ID and user ID) — and do not feed role resolution.
- **DynamoDB — Streams and Global Tables**: implemented with multi-active replication.
- **IAM — GetHumanReadableSummary**: excluded; the summary requires external LLM generation.
- **Kinesis — SubscribeToShard heartbeat interval**: 15 s (provisional; AWS does not document the exact value).
- **Lambda — AddPermission Principal**: restricted to a known service-principal allowlist (see `validServicePrincipals` in `validators.go`); unrecognised `*.amazonaws.com` principals are rejected.
- **Lambda — resource-policy operations**: PutResourcePolicy/GetResourcePolicy/DeleteResourcePolicy share the AddPermission policy substrate — one policy per function, so a Put replaces permissions added through AddPermission and both GetPolicy and GetResourcePolicy serve the same statements, with Put-path statements reproduced from their verbatim JSON. Mutations carry a policy revision (RevisionId precondition answered with PreconditionFailedException); the policy document is bounded by the 20,480-character quota, and DeleteResourcePolicy succeeds on a resource whose policy is already gone.
- **Lambda — resource-policy public-access rejection**: excluded; the platform accepts resource-policy documents that grant wildcard-principal access instead of answering PublicPolicyException. The Smithy model declares the error on AddPermission, PutResourcePolicy, and RemovePermission, and its message names the trigger loosely — a policy that "would grant public access to it, which isn't allowed" — but defines no boundary for public access; the Lambda documentation's own valid examples include Condition-scoped wildcard-principal Allow statements, and current AWS enforcement is tied to the account-level block-public-access setting, which the vendored models do not carry.
- **Lambda — RuntimeManagementConfig**: excluded; the configuration API drives a managed runtime-update mechanism (a runtime-version registry that pins and advances each function's runtime patch level). This platform has no update engine — `RuntimeVersionConfig` is configuration echo only — so the API surface would be inert.
- **Lambda — RecursionConfig**: excluded; its only enforced semantics is the recursion loop detector's behaviour (Terminate/AllowOnFail). This platform has no loop-detection engine, so the default Terminate would report terminating behaviour that never occurs.
- **Lambda — ScalingConfig**: excluded; the model scopes the API to Managed-Instances functions running on Lambda capacity providers, which this platform does not provide.
- **S3 — inventory report delivery**: reports deliver on daily/weekly UTC boundaries to the S3 destination as CSV (gzip), Parquet (snappy), and ORC (ZLIB), with manifest.json, manifest.checksum, and the Hive symlink; report files honour the configuration's SSE-S3/SSE-KMS encryption choice. The report columns IntelligentTieringAccessTier and ChecksumAlgorithm are emitted empty (no access-tier or checksum substrate); the LifecycleExpirationDate column carries the ISO-8601 projection of an applicable expiration rule and stays empty when none applies or the object's replication has not succeeded. Both configuration families are bounded by the 1,000-configuration limit with 100-item pagination.
- **S3 — metrics configurations**: per-filter CloudWatch request metrics in the AWS/S3 namespace; requests on both the object and bucket planes count into AllRequests plus their per-operation metric, and each minute window publishes a CloudWatch statistic set (sample count, sum, min, max, so Average carries the documented error rate and bytes-per-request semantics). A filter carrying an access-point ARN generates no datapoints (no access-point substrate).
- **S3 — analytics/intelligent-tiering configurations**: excluded; this platform has no access-tier machinery (its storage classes are software classifications) and no access-pattern analysis substrate to configure. Lifecycle-driven storage-class transitions are implemented — see the lifecycle enforcement note.
- **S3 — storage classes**: the eight software classes (STANDARD, REDUCED_REDUNDANCY, STANDARD_IA, ONEZONE_IA, INTELLIGENT_TIERING, GLACIER, GLACIER_IR, DEEP_ARCHIVE) persist and round-trip on puts, copies, replication destinations, and lifecycle transitions; the acceptance set is defined once beside the store constants and every validation site derives from it. The AWS enum's hardware- and backup-bound classes (OUTPOSTS, SNOW, EXPRESS_ONEZONE, FSX_OPENZFS, FSX_ONTAP, AWS_BACKUP_WARM, AWS_BACKUP_LOW_COST_WARM) are rejected at validation: they name physical media and backup tiers this platform does not have.
- **S3 — lifecycle enforcement**: Transition and NoncurrentVersionTransition are executed by the lifecycle sweep. A Days window counts from the object's LastModified for current versions and from the successor version's creation for noncurrent versions, and rounds up to the next midnight UTC — enforcement and the x-amz-expiration projection share the one calculation; a Date entry applies to every matching object once the date has passed. Entries replay in rule order under the supported-transitions waterfall, so a move the waterfall forbids is skipped rather than regressing a class. The default 128 KiB transition minimum applies unless the rule's filter carries an explicit size bound; objects with a Pending or Failed replication status are not transitioned; NewerNoncurrentVersions gates both noncurrent actions on an exceeded newer-noncurrent count and requires a rule filter. Get and Head responses carry x-amz-expiration (expiry-date at HTTP-date, URL-encoded rule-id) when an enabled expiration rule projects onto the current version, and expiration and transition actions are prevented on objects whose replication has not succeeded. The `x-amz-transition-default-minimum-object-size` header is not implemented: it selects a pre-September-2024 AWS compatibility behaviour for configuration bases that predate the current default, and this platform has none.
- **S3 — Select**: excluded; AWS no longer offers S3 Select to new customers, and the platform does not reimplement it — `POST /{bucket}/{key}?select` is refused with 501 NotImplemented. The operation remains in the AWS API surface for existing AWS customers, so workloads that require it are not supported here.
- **RDS Data — enablement**: queries run on the vmysql engine; the service requires `RDS_MYSQL_ENABLED=true` or `ALL_SERVICES_ENABLED=true`.
- **SFN — task integration surface**: the optimised and AWS SDK integration forms run for the services the platform carries — Lambda (invoke), SQS (sendMessage), SNS (publish), EventBridge (putEvents), DynamoDB (getItem, putItem, updateItem, deleteItem), and the Step Functions self-integration (startExecution across the plain, .sync, .sync:2 and waitForTaskToken patterns, with the callback pattern available wherever AWS documents it). Definition validation enforces the per-service pattern and action combinations AWS documents; resource forms naming services the platform does not carry (batch, ECS, Glue, and the rest of the AWS SDK namespace) remain valid definitions and fail at run time with a cause naming the unavailable integration, the same way a task referencing any other nonexistent resource fails.
- **Scheduler — templated targets**: platform-implemented targets are Lambda, SQS, SNS, Kinesis, Step Functions, and EventBridge; ECS and Firehose targets are accepted by rule templates but delivery fails until those services exist on the platform; SageMaker, CodeBuild, CodePipeline, and Inspector targets are permanently out of scope (those services are not implemented on this platform).
- **Secrets Manager — ListTagsForResource and managed rotation members**: the operation does not exist in the 2017-10-17 model, so AWS SDKs never generate a client method for it; the platform operation serves raw-HTTP/console consumers. Managed external rotation members are configuration storage and echo only — the partner integration itself is external.
- **SQS — SSE-KMS and request throttling**: SSE-KMS attributes are accepted but messages are stored unencrypted; per-account request-rate quotas are not enforced, and no RequestThrottled error is ever returned (the shape lives only in the generated wire types).
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

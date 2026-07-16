# Changelog

All notable changes to Vorpalstacks will be documented in this file.


## [0.1.0] - Upcoming

Vorpalstacks will switch to **MIT + Apache 2.0** dual licensing,
replacing the current FSL-1.1-MIT (root) / Apache 2.0 (pkg/) split.

### Added

- **API Gateway** — MethodSettings and CanarySettings patch on UpdateStage (metrics, logging, throttling, caching, canary traffic, stage variable overrides; `remove` resets to defaults); stage throttling enforcement (token-bucket rate limiter, wildcard settings merge `*/*` → exact `path/method`, live rate/burst updates without restart); HTTP and DynamoDB non-proxy integration pipeline (request templates, passthrough behaviour NEVER/WHEN_NO_TEMPLATES/WHEN_NO_MATCH, response selection by status code, response templates); binary media types (Lambda proxy base64 encoding/decoding; `type/*` and `*/*` wildcards); integration content handling (CONVERT_TO_BINARY on request, CONVERT_TO_TEXT on response); stage variable substitution in integration URIs (`${stageVariables.varName}`); per-integration timeout (`TimeoutInMillis`; 29 s default via `context.WithTimeout`).
- **AppSync** — PER_RESOLVER_CACHING (resolver result cache with SHA-256 key from `$context.arguments.*`/`$context.source.*` paths, TTL-based expiry, FlushApiCache cache+schema invalidation); API key expiry validation (1–365 days per AWS spec); API key count limit enforcement (5000 per GraphQL API); CachingConfig TTL validation (1–3600 seconds).
- **CloudTrail** — CloudTrail Lake query API (StartQuery with SQL-like SELECT/FROM/WHERE parser supporting eventName, eventSource, username, accessKeyId, readOnly, and eventTime range filters; GetQueryResults with column-order-preserving `[][]map[string]string` storage; DescribeQuery, CancelQuery, ListQueries); channel management for external data ingestion (Create/Get/List/Update/DeleteChannel); EDS ingestion control (Start/StopEventDataStoreIngestion); EDS soft-delete with RestoreEventDataStore (PENDING_DELETION status with DeletedTimestamp); EDS federation (Enable/DisableFederation with FederationRoleArn); event configuration (Get/PutEventConfiguration); delegated admin registration (Register/DeregisterOrganizationDelegatedAdmin); EDS name validation (3–128 chars, `^[a-zA-Z0-9._\-]+$`); RetentionPeriod validation (7–3653 days); tag validation (max 50, key 1–128 chars, value 0–256 chars, `aws:` prefix reserved) on CreateTrail, CreateEventDataStore, and TagResource; LookupEvents MaxResults capped at 50.
- **CloudWatch** — Extended statistics computation engine (percentile p{n} with linear interpolation, trimmed mean tm{n}, trimmed count tc{n}, trimmed sum ts{n}, winsorized mean wm{n}, interquartile mean IQM); ExtendedStatistic support on PutMetricAlarm (mutually exclusive with Statistic, format-validated), alarm evaluator, GetMetricStatistics (ExtendedStatistics parameter), GetMetricData (automatic extended-stat routing), and dashboard widgets; anomaly detection comparison operators (LessThanLowerOrGreaterThanUpperThreshold, LessThanLowerThreshold, GreaterThanUpperThreshold); PutMetricAlarm validation (Statistic whitelist, Statistic/ExtendedStatistic mutual exclusivity, TreatMissingData enum, Period 10/20/30/multiples of 60, Period×EvaluationPeriods ≤ 7 days, DatapointsToAlarm ≤ EvaluationPeriods, Dimensions ≤ 30, Actions ≤ 5 per type, AlarmName ≤ 255 chars, tag validation max 50/key 1–128/value 0–256/`aws:` prefix reserved); PutCompositeAlarm validation (AlarmName ≤ 255 chars, Actions ≤ 5 per type, tag validation); PutMetricData validation (≤ 1000 entries, ≤ 30 Dimensions per datum, ≤ 150 Values, Value/Values/StatisticValues mutual exclusivity via presence flag, Values/Counts length match); Unit and ExtendedStatistic fields in alarm response; metric math expression engine (`pkg/metricmath`: tokenizer, recursive-descent parser with operator precedence +− < ×÷ < ^ right-associative < unary, AST evaluator with timestamp-aligned binary ops, scalar broadcasting, and 12 functions — FILL, AVG, SUM, MIN, MAX, ABS, CEIL, FLOOR, LOG10, LN, SQRT, EXP); GetMetricData metric math (3-phase pipeline: MetricStat resolution → expression dependency-ordered evaluation → `ReturnData`-filtered response; period-aligned grid expansion with NaN for missing periods enabling FILL); composite alarm evaluation (AlarmRule boolean parser supporting `ALARM("name")`, `AND`, `OR`, `NOT`, `TRUE`, `FALSE`; level-by-level topological sort via Kahn's algorithm with cyclic-dependency detection; ActionsSuppressor with WaitPeriod and ExtensionPeriod suppression windows); metric math alarm evaluation (Metrics array with ThresholdMetricId, dependency-ordered expression resolution); AlarmRule syntax validation at PutCompositeAlarm creation time; ActionsSuppressedBy/Reason, ActionsSuppressor, ThresholdMetricId, and Metrics fields in alarm response; anomaly detection models (PutAnomalyDetector for single-metric keyed on namespace/metric/dimensions/stat and metric-math keyed on query hash; DeleteAnomalyDetector; DescribeAnomalyDetectors with namespace/metric/dimensions/type/ID filtering); ANOMALY_DETECTION_BAND in GetMetricData and alarm evaluation (EWMA + N×stddev upper/lower band, NaN filtering, two-pass metric resolution); DescribeAlarmContributors with band-breach data points; Contributor Insights rules (Put/Delete/Describe/Enable/DisableInsightRules, GetInsightRuleReport with period-bucketed metric datapoints, PutManagedInsightRules with template-based rule generation, ListManagedInsightRules with ResourceARN filter); alarm mute rules with per-region evaluator suppression (Put/Delete/Get/ListAlarmMuteRules, SCHEDULED/ACTIVE/EXPIRED status computation, ACTIVE rule check suppresses alarm action dispatch); PutLogAlarm (log query result alarm with scheduled query configuration); PutMetricAlarm EvaluationCriteria with PromQLCriteria (Query, PendingPeriod, RecoveryPeriod); OTel enrichment control (Get/Start/StopOTelEnrichment with mutex-protected state); dataset KMS key association (GetDataset, Associate/DisassociateDatasetKmsKey with mutex-protected map).
- **Cognito Identity** — STS-backed credential issuance for GetCredentialsForIdentity (CredentialIssuer interface injected via wiring, resolving authenticated/unauthenticated roles from identity pool config with CustomRoleArn override, issuing real AssumeRoleWithWebIdentity sessions); RS256-signed JWT OpenID tokens for GetOpenIdToken and GetOpenIdTokenForDeveloperIdentity (lazy RSA key pair via vsjwt, iss/sub/aud/iat/exp/jti claims, kid header); GetId identity reuse for authenticated logins (FindIdentityByLogins prevents duplicate identity creation); LookupDeveloperIdentity NextToken pagination; GetOpenIdTokenForDeveloperIdentity TokenDuration parameter (15–3600 seconds); MergeDeveloperIdentities logins consolidation (source identity provider links transferred to destination before deletion).
- **Cognito IDP** — Per-client token validity (AccessTokenValidity, IDTokenValidity in minutes, RefreshTokenValidity in days from UserPoolClient config, replacing hardcoded 60/60/30 defaults); ID token persistence (stored alongside access and refresh tokens for validation and revocation); SignUp enforcement of AdminCreateUserConfig.AllowAdminCreateUserOnly (rejects self-registration when admin-only creation is enabled); GenerateSecret=false suppression on CreateUserPoolClient (JSON bool and Query string protocols); AllowedOAuthFlowsUserPoolClient parse and always-on response; PreventUserExistenceErrors as string on UserPoolClient (was bool on PasswordPolicy); PreSignUp Lambda trigger now receives ValidationData and ClientMetadata from request; identity provider AttributeMapping and IdpIdentifiers on Create/Update/Describe; DeleteUserPool cascade deletion of resource servers, identity providers, and domains.
- **DynamoDB** — Kinesis Data Stream destinations (item change records dispatched to configured streams on PutItem/DeleteItem/UpdateItem/BatchWriteItem/TransactWriteItems); Global Tables multi-active replication (atomic propagation of puts/deletes with index entries, item count, and table size to all replica regions); TTL expiration worker (5-minute scan, atomic deletion with stream records tagged `{type:"Service", principalId:"dynamodb.amazonaws.com"}`); GSI backfill on UpdateTable (existing items auto-populated when new GSI created); auto-scaling settings round-trip for DescribeTableReplicaAutoScaling/UpdateTableReplicaAutoScaling.
- **EventBridge** — Scheduled rule evaluation (rate() and cron() expressions with UTC minute-boundary aligned ticks; cron supports wildcards, ranges, lists, steps, and DOW names); target retry policy enforcement (configurable MaximumRetryAttempts/MaximumEventAgeInSeconds, exponential backoff with jitter, DLQ routing after exhaustion); Kinesis target PartitionKeyPath (JSON path extraction from payload); SQS target MessageGroupId (FIFO queue delivery); archive retention worker (hourly purge of expired events with counter decrement); PutTargets 5-targets-per-rule quota enforcement.
- **IAM** — Pagination (Marker/MaxItems) and PathPrefix filtering on `ListAttachedPolicies`; pagination on `ListInlinePolicies`; tag validation (key 1–128 chars, value 0–256 chars, max 50 per resource) on all Create and TagResource operations; policy document structure validation (requires `Statement` with `Effect` of `Allow` or `Deny`).
- **IoT** — Fleet indexing query engine (SearchIndex, GetStatistics, GetCardinality, GetPercentiles, GetBucketsAggregation; recursive-descent query parser supporting thingName/thingType/attributes.*/isConnected, comparison operators >, >=, <, <=, <>, boolean AND/OR/NOT, parenthesised grouping, and * wildcard matching; aggregation pipelines for count/sum/min/max/average, cardinality, percentile linear interpolation, and term buckets); MQTT broker mTLS enforcement (tls.RequireAndVerifyClientCert with CA-root-signed server certificate regenerated on each start; client identity from TLS certificate fingerprint exclusively — MQTT Username ignored per AWS spec); connection lifecycle tracking (OnSessionEstablished/OnDisconnect hooks recording cert ID → Unix timestamp via sync.Map); GetThingConnectivityData live status from broker tracker; IoT policy condition operators (StringEqualsIgnoreCase, StringNotEqualsIgnoreCase, Bool case-insensitive, ARN operators, IpAddress/NotIpAddress CIDR matching, Numeric comparison operators, Date comparison operators with RFC 3339 and epoch ≥10-digit parsing); IoT policy version management (5-version limit, default version document sync to Policy record for broker, default version deletion protection, SetDefaultPolicyVersion not-found error); IoT rules substitution templates (RepublishTopic, BucketName, QueueURL, TopicARN, StreamName, DynamoDB hashKeyValue/rangeKeyValue/operation/payloadField, CloudWatch metricName/metricNamespace/metricUnit/metricValue/metricTimestamp); IoT rules DynamoDB action operations (INSERT/UPDATE/DELETE via operation field); IoT rules SQL three-valued logic (AND/OR/NOT with UNKNOWN semantics from IS NULL); IoT rules republish payload copy-on-write isolation; ListMetricValues fleet metric query execution; paho.golang E2E MQTT test suite (mTLS connect, pub/sub round-trip, policy enforcement); crypto utilities consolidation (ca/certid.go → internal/utils/crypto: FingerprintPEM, FingerprintX509, SHA256Hash, key PEM encode/decode, ECDSA generation, certificate creation/verification, serial number generation).
- **Kinesis** — ExplicitHashKey shard routing on PutRecord/PutRecords (direct hash-key-to-shard mapping bypassing MD5 partition key hashing; invalid keys rejected with InvalidArgumentException); GetRecords retention period filtering (records older than the stream's RetentionPeriodHours excluded); Stream WarmThroughput persistence on UpdateStreamWarmThroughput; StartStreamEncryption validation (EncryptionType=KMS required, KeyId required); UpdateShardCount ScalingType validation (UNIFORM_SCALING only); ListStreams limit expanded to 10000.
- **KMS** — EncryptionAlgorithm parameter resolution in Encrypt/Decrypt/ReEncrypt (client-specified algorithm or key-spec default: SYMMETRIC_DEFAULT for symmetric, RSAES_OAEP_SHA_256 for RSA; validated against key's supported list); tag quota enforcement on CreateKey and TagResource (max 50 tags, key 1–128 chars, value 0–256 chars, `aws:` prefix reserved, merged cumulative count check on TagResource); EnableKeyRotation restricted to SYMMETRIC_DEFAULT keys with AWS_KMS origin.
- **Lambda** — Asynchronous invocation retry with EventInvokeConfig (configurable MaximumRetryAttempts, MaximumEventAgeInSeconds, exponential backoff capped at 60 s); async destination delivery (OnSuccess/OnFailure to SQS/SNS/Lambda/EventBridge); event source mapping filter criteria for SQS, Kinesis, and DynamoDB Streams (prefix, suffix, numeric range, exists, anything-but, equals-ignore-case; JSON body descent for SQS); alias weighted routing (RoutingConfig additional version weights with random selection); layer version policy APIs (AddLayerVersionPermission, RemoveLayerVersionPermission, GetLayerVersionPolicy, GetLayerVersionByArn); reserved concurrency enforcement (atomic inflight counter with CAS); LogType=Tail (base64-encoded last 4 KB of execution logs); provisioned concurrency pre-warm (container created on config put).
- **S3** — Cross-region replication (Put/Get/DeleteBucketReplication; async object copy with prefix and tag filter evaluation; delete-marker propagation when DeleteMarkerReplication is Enabled; multi-region destination bucket resolution via `FindBucket`); CORS preflight handling (OPTIONS) and actual-request `Access-Control-*` header injection from bucket CORS rules; background lifecycle worker (expiration by Days/Date, noncurrent-version expiration, abort-incomplete-multipart-upload; 5-minute scan interval with graceful shutdown); Object Lock enforcement on DeleteObject and DeleteObjects (Legal Hold blocks unconditionally; COMPLIANCE retention blocks unconditionally; GOVERNANCE retention bypassable with `x-amz-bypass-governance-retention: true`); PutObject `x-amz-tagging` header parsing and persistence.
- **Secrets Manager** — ClientRequestToken support in CreateSecret and PutSecretValue (version ID override with 32–64 char validation and idempotency: same token + same value returns existing version, same token + different value rejected); RotateImmediately parameter in RotateSecret (defaults true; when false, only rotation config is updated); UpdateSecretVersionStage three-mode operation (add label via MoveToVersionId, remove label via RemoveFromVersionId, move label with both); PutSecretValue non-AWSCURRENT staging label support (versions created without demoting existing current version); metadata-only UpdateSecret (no new version when SecretString/SecretBinary absent); tag quota enforcement on CreateSecret and TagResource (max 50 tags, key 1–128 chars, value 0–256 chars, merged total count check on TagResource); EventBus logger injection; CloudWatch alarm evaluator immediate first evaluation on startup; EventBus outbox pending retention (stale OutboxPending entries purged after 30 minutes, preventing accumulation from crash recovery).
- **Step Functions** — Choice state operator expansion (StringLessThanEquals, StringGreaterThanEquals, StringMatches glob with `*` wildcard, NumericLessThanEquals, NumericGreaterThanEquals, TimestampLessThanEquals, TimestampGreaterThanEquals, IsNull, IsBoolean, IsString, IsNumeric, IsTimestamp, all 16 `*Path` dynamic-value comparisons); Wait state Timestamp and TimestampPath; Map and Parallel state ResultPath; TaskStateEntered/TaskStateExited history events; Retry MaxDelaySeconds cap; SQS SendMessage integration passthrough (DelaySeconds, MessageGroupId, MessageDeduplicationId, MessageAttributes); SNS Publish MessageAttributes passthrough.

### Fixed

- **API Gateway**: response parameter mapping key/value corrected (key is destination `method.response.header.{name}`, value is source — was inverted; static single-quoted values now supported); response template content-type no longer hardcoded to `application/json` (now checks backend `Content-Type` then client `Accept`).
- **AppSync**: `$util.appendError` incorrectly halted resolver execution — was treated identically to `$util.error`/`$util.unauthorized`, returning nil data instead of partial data with collected errors (now distinguished via `FatalError` flag on `AppSyncError`; non-fatal errors allow continued dispatch and return data alongside errors in both unit and pipeline resolvers).
- **CloudTrail**: UpdateTrail now uses parameter existence checks so explicitly-provided empty strings clear fields (SnsTopicName, S3KeyPrefix, CloudWatchLogs ARNs, KmsKeyId — was ignoring empty strings); CreateTrail now requires S3BucketName (was silently accepting empty); DeleteTrail now resolves trail ARN in addition to name (was name-only); DeleteEventDataStore now soft-deletes to PENDING_DELETION status instead of hard-deleting (AWS spec compliance).
- **CloudWatch**: ComparisonOperator name corrected from `LessThanOrEqualToLowerOrGreaterThanUpperThreshold` to `LessThanLowerOrGreaterThanUpperThreshold` (AWS spec); GetMetricData, dashboard widget, and metric data result statistic matching now case-insensitive (was exact-match, rejecting lowercase stat names from some SDKs); widget unknown statistics now resolve from ExtendedStats instead of silently defaulting to Average.
- **CloudWatch**: integration alarm tests now use Period 60 with period-aligned timestamps (was Period 1, rejected by validation added in previous batch; missing Timestamp caused empty evaluation windows).
- **CloudWatch**: alarm mute rule status recomputed at read time in GetAlarmMuteRule and ListAlarmMuteRules so SCHEDULED→ACTIVE→EXPIRED transitions are always accurate without a background reaper; ErrAlarmMuteRuleNotFound sentinel error for errors.Is discrimination; ListAlarmMuteRules summary response reduced to 5 fields; anomaly band computation extracted to shared prepareAnomalyBand; anomaly detector response emits SingleMetricAnomalyDetector nested object and StateValue (was flat State — AnomalyDetectorType removed as non-existent in AWS API spec); PutLogAlarm action arrays now accept lowercase parameter aliases.
- **Cognito Identity**: GetCredentialsForIdentity was returning fabricated UUID-based credentials instead of real STS sessions (now issues AssumeRoleWithWebIdentity via injected CredentialIssuer); GetOpenIdToken and GetOpenIdTokenForDeveloperIdentity were returning raw UUID strings instead of signed JWT tokens; GetId was always creating a new identity even when an existing identity with matching logins existed (now reuses existing identity for authenticated requests); AllowClassicFlow and ServerSideTokenCheck always present in responses (were omitted when false, inconsistent with AWS API).
- **Cognito IDP**: AdminDeleteUser now deletes all user tokens (access, ID, refresh — was orphaning valid tokens after user deletion); PreSignUp Lambda trigger was passing nil for ValidationData and ClientMetadata instead of request values (now parsed from JSON and Query protocol parameters).
- **DynamoDB**: condition expressions (`attribute_exists`/`attribute_not_exists`/`size`/comparisons) now resolve nested document paths (`a.b.c`, `a[0].b`, `a[0][1]`); `importDynamoDBJSONData` now updates GSI index entries, item count, and table size atomically (was raw `Items().Put` that skipped indexes); `BatchWriteItem` now dispatches Kinesis destination records for successful operations.
- **EventBridge**: PutEvents now respects entry Time field (was always `time.Now()`); `matchOperator` numeric patterns now fail on malformed input (was `continue`-skipping and matching everything); `exists:false` content filter now matches absent fields in matchEventPattern and TestEventPattern; cascade-delete on DeleteEventBus now paginates all rules and targets (was truncating at 1000); replay goroutine panic now logged and replay marked FAILED (was silently swallowed); `DeleteExpiredArchiveEvents` SizeBytes decrement now uses JSON byte length (was map key count, causing counter drift).
- **IAM**: `DeleteGroup` and `DeleteUser` now return `DeleteConflictException` when dependent resources exist (was silently cascade-deleting); `RenameUser` restructured to two-phase migration with rollback (all resources migrated before user record swap; completed steps reversed on failure); `RenameUser` now migrates signing certificates, SSH public keys, and service-specific credentials (were left under old name); `DeleteUser` and `cascadeDeleteUser` now clean up service-specific credentials (was orphaned).
- **IoT**: MQTT broker always installed allowAllHook because authProvider was nil at Start() time, permitting all connections without authentication (replaced with fail-closed AuthHook that reads provider lazily via getAuthProvider); policy lookups in OnConnectAuthenticate and OnACLCheck passed raw SHA-256 cert ID instead of full principal ARN to getPoliciesForPrincipal, silently failing all MQTT policy enforcement (added CertificatePrincipal to AuthProvider for ARN construction); iot:Connect resource matching used MQTT wildcards (+/#) instead of AWS policy wildcards (*/?), preventing client ARN patterns like client/* from matching (engine now selects wildcardMatch for iot:Connect, topicMatch for Publish/Subscribe/Receive); connection tracker hooks resolved cert ID from MQTT Username instead of TLS certificate fingerprint, allowing connectivity-status spoofing (all hooks now use extractCertificateID exclusively).
- **Kinesis**: GetRecords MillisBehindLatest now calculated from last record's ApproximateArrivalTimestamp (was hardcoded 0); PutRecord and PutRecords now validate data size against base64-decoded payload length using the stream's MaxRecordSizeInKiB (was checking base64-encoded length against hardcoded 1 MiB); PutRecord no longer rejects empty Data (AWS allows empty-data records); UpdateMaxRecordSize upper bound raised to 10240 KiB (was 1024).
- **KMS**: ListKeyRotations no longer fabricates a rotation entry from CreationDate (returns empty list when no rotations tracked); EnableKeyRotation/DisableKeyRotation now reject Disabled, PendingDeletion, and PendingImport key states (DisabledException/KMSInvalidStateException per AWS key-state table); KeyStore.Enable rejects PendingImport (was enabling key without material); KeyStore.Disable rejects PendingDeletion (use CancelKeyDeletion); ScheduleDeletion rejects double-scheduling and saves/restores KeyRotationEnabled on deletion/cancel cycle.
- **Lambda**: ESM filtered-out records now acknowledged on all event sources — SQS messages deleted, Kinesis/DynamoDB Streams checkpoints advanced (was infinite re-polling after visibility timeout); async retry loop now uses detached `context.Background()` (was inheriting HTTP request context, killing retries and destination delivery after 202 response); EventBridge destination delivery now extracts event bus name via `ExtractEventBusNameFromARN` (was passing raw resource string `event-bus/default`); filter `exists:false` now matches absent keys regardless of other criteria in the same array (OR semantics was incorrectly restricted to single-criterion arrays).
- **S3**: versioning cache stale reads after status change (ristretto `Wait()` added to `VersioningCache.Set` for immediate visibility).
- **Secrets Manager**: DeleteSecret now validates RecoveryWindowInDays (7–30 inclusive) and rejects ForceDeleteWithoutRecovery + RecoveryWindowInDays combination and deletion of replicated primary secrets; PutSecretValue now requires at least one of SecretString/SecretBinary and rejects both simultaneously; DeletedDate in DescribeSecret now set to scheduled deletion date per AWS spec (was request timestamp); nil pointer panics prevented in rotation engine (LambdaInvoker nil check) and Lambda async destination delivery (SQS/SNS/Lambda/EventBridge invoker nil checks); failed Lambda rotation no longer persists incorrect LastRotatedDate/NextRotationDate (save/restore on error); NextRotationDate now calculated after LastRotatedDate update; AWSCURRENT staging label cannot be removed without specifying a target version; UpdateSecretVersionStage now supports remove-only mode (MoveToVersionId omitted); error mapping unified (fmt.Errorf → mapStoreError across rotation and PutSecretValue paths).
- **SNS**: filter policy now enforced before each delivery; failed SQS/HTTP/Lambda deliveries routed to the subscription's DLQ via RedrivePolicy; FilterPolicy/FilterPolicyScope/RedrivePolicy validated at Subscribe and SetSubscriptionAttributes; subscription attributes consolidated behind `Attributes`-map accessors (was reading a dead `RawMessageDelivery` field); `anything-but` now requires attribute presence to match AWS; invalid filter policy JSON fails closed.
- **SNS**: signing certificate serial number hardcoded to 1 (now random 128-bit via `vcrypto.GenerateSerialNumber`).
- **SQS**: `maxReceiveCount` now honored correctly when moving to DLQ (was `>=`, delivering one fewer receive than configured); `ReceiveMessage` now filters `MessageAttributes` and system `Attributes` per `MessageAttributeNames` / `AttributeNames` / `MessageSystemAttributeNames` (was returning all unconditionally); specifying both `AttributeNames` and `MessageSystemAttributeNames` now returns `InvalidParameterCombination`; background cleanup goroutine deletes messages exceeding `MessageRetentionPeriod` (storage hygiene; receive-time filter was already correct); store `Close()` registered as shutdown hook.
- **Step Functions**: `States.ALL` retry/catch now correctly excludes `States.DataLimitExceeded` and `States.RuntimeExceeded` (was matching all errors unconditionally, preventing proper handling of uncatchable execution-limit failures).

## [0.0.12] - 2026-07-12

### Added

- **IoT Core** — Full MQTT broker with SHA-256 certificate-fingerprint authentication and IoT policy enforcement (Connect/Publish/Subscribe ACLs); thing registry, device shadows, certificate/CA/policy management, job execution fleet management, provisioning templates, topic rules engine (SQL with IN/NOT IN operators; 8 action types including DynamoDB/Kinesis/Lambda/SNS/SQS/Republish/IoTEvents/HTTP), per-region task engine with panic recovery, and CloudTrail audit stream integration.

- **IoT Events** — Detector model state machine (onInput/onEnter/onExit lifecycle, EventBus dispatch), alarm models, and data-plane operations (BatchPutMessage, Batch*Detector, Batch*Alarm) with IoTEventsData signing classifier.

- **DynamoDB Streams** — Stream ARN management, shard iterators, GetRecords with sequence-number pagination, atomic stream capture within item mutation transactions (PutItem, UpdateItem, DeleteItem, BatchWriteItem, TransactWriteItems, PartiQL), and DescribeStream with StreamViewType-aware image serialization.

- **RDS: embedded MySQL engine and Data API** — Isolated MySQL instances via go-mysql-server backed by Pebble KV storage (`RDS_MYSQL_ENABLED=true`); RDS Data API (ExecuteStatement, BatchExecute, Begin/Commit/Rollback Transaction) with named parameter substitution and 5-minute transaction TTL.

- **Neptune: TinkerPop Gremlin Server WebSocket** — Native WebSocket endpoint at `/gremlin` with GraphSON v3 serialization; per-cluster isolated graph engine with dedicated HTTP listener (REST + WebSocket); engines restored on restart.

- **S3: streaming chunked encryption** — `EncryptStream`/`DecryptChunked` process 64 MB chunks through AES-GCM without buffering entire body; per-chunk `PartEncryptionInfos` in SSE metadata.

- **Web console enhancements** — 3-panel S3 object browser (breadcrumb navigation, folder traversal, upload/download, batch delete, text/image preview); RDS instance management with engine lifecycle; PebbleDB key inspector; DynamoDB nested table→items navigation with structured attribute editor; CloudWatch alarm create dialog.

- **`ALL_SERVICES_ENABLED` environment variable** — Set to `true` to force-enable every service at once, overriding individual `*_ENABLED` flags.

- **API Gateway: validation and pagination** — List operations support `limit`/`position` pagination; `UpdateRestApi`, `CreateAuthorizer`, `CreateStage` validate input.

### Changed

- **CloudTrail audit config unified to `CLOUDTRAIL_ENABLED`** — Replaces `VS_AUDIT_ENABLED` and `features.audit_enabled`; now an Optional service (default `false`), propagated at startup.

- **Protobuf definitions fully regenerated** — All 34 service `.proto` files regenerated; webconsole TypeScript generation migrated from `protoc-gen-js` to `protobuf-es`; legacy `webconsole/src/gen/` directory removed (~270K lines deleted).

- **Web console: all service pages migrated to unified 3-panel inspector layout** — DataTable + Splitter + inspector components with global text filter, JSON/CSV export, and accessibility improvements (focus trap, `aria-modal`).

- **Neptune service packages relocated under `rds/`** — Neptune, NeptuneData, NeptuneGraph implementations and stores moved to `internal/{services,store}/aws/rds/`.

- **Dispatcher: removed resilience wrapper** — HTTP request dispatcher no longer wraps handler invocations in the resilience layer.

- **Internal refactoring across 20+ services** — Error mapping unified via `errors.Is` sentinel checks; store-level pagination (`common.List`, `ListMatching`, `ListMatchingProto`); response serialization consolidation; deterministic sorting for parameter lists and tags; dead code removal.

### Fixed

- **IAM**: `aws:*` condition key resolution (was comparing resolved value as key); explicit Deny precedence within a policy regardless of statement order; trust policy matching for Federated/Service/Account principals; `ResetPassword` lost-update race (added key lock); `GetAccountAuthorizationDetails` pagination beyond 1000.

- **S3**: PebbleDB key delimiter changed to null byte (`\x00`) eliminating collision risk with keys containing `#`; KeyLocker mutex race on persistent lock keys; multipart part sorting before assembly; versioned Range requests; `UploadPartCopy` Range propagation; `IsLatest` flag consistency on version delete; `CreateMultipartUpload` cleanup on failure; object lock operations I/O optimization; malformed XML/invalid part-number-marker error handling; `_latest` pointer lazy-write race eliminated.

- **DynamoDB**: Key attribute (partition/sort key) update rejection in `UpdateItem`/`TransactWriteItems`; item size calculation (AWS formula for Number/Map/List); stream OLD image capture for `OLD_IMAGE`/`NEW_AND_OLD_IMAGES` view types; atomic stream capture within transactions (was separate post-commit write with data-loss risk); backup scan error propagation; streams cleanup on `DeleteTableCascade` (now inside the transaction); pagination offset panic guard.

- **IoT**: Pagination offset out-of-range panic guard; task engine panic recovery (deferred closure fix); republish iteration counter type-assertion fix (was infinite loop); rule action map defensive copy (was overwriting persistent templates); shadow/`DeleteThing` unified lock domain (was concurrent, causing data loss); provisioning template body JSON validation on Create/Update/Proto conversion; billing group `ScanPrefix` error propagation; `ListShadowNames` pagination defaults and marker-not-found handling; `ioteventsdata` SigV4 signing classifier mapping.

- **Auth**: SigV4 bypass restricted to STS service only (was exploitable via form-encoded body); SigV4 RFC 3986 percent-encoding correction; `X-Forwarded-For` header ignored by default (security hardening).

- **EventBus**: Outbox writes now `pebble.Sync` for crash-safe at-least-once delivery; `UpdateStatus` serialized behind per-event mutex; `Publish` fail-fast on missing `EventRegistry`.

- **Misc**: Blob hybrid I/O error propagation (was silently truncating); `ProvisioningTemplate` proto conversion error propagation (`ProtoStoreConfig` interface updated to return errors); CloudWatch `ActionsEnabled` proto field promoted to `optional` for gRPC-Web parity; Cognito `SignUp` `AutoConfirmUser` inversion; Step Functions retry `MaxAttempts` off-by-one; Step Functions activity queue goroutine leak; Timestream query row double-wrap; zip-slip protection in `vstackscli` backup.

### Removed

- **Legacy webconsole TypeScript protobuf generation** — 37 files (~270K lines) from `webconsole/src/gen/` replaced by protobuf-es output.

## [0.0.11]

### ⚠️ Breaking Changes

- **Management console rewritten from Flutter/Dart to TypeScript/React**: Removed `web_ui/` (Flutter) and replaced with `webconsole/` (TypeScript, React, Vite, Tailwind CSS). The web console is now embedded into the Go binary via `//go:embed`. Build requires Node.js 22; `make build` runs proto generation and webconsole build automatically.
- **Configuration system overhauled**: Removed `ConfigSource` tracking (ENV/STORE/DEFAULT). Replaced `server.bind_addr` with `server.bind_mode` (all/localhost/interface) + `server.bind_interface`. Port defaults centralized in new `internal/common/serviceports/` package; environment variable reads removed from defaults layer.
- **Port allocation restructured**: Service endpoint ports defined as constants in new `serviceports` package. Added dynamic port allocator (`internal/server/portalloc/`) with FQDN/Individual port mode per service. Dynamic port range configurable via `VS_PORT_DYNAMIC_START`/`VS_PORT_DYNAMIC_END`.

### Added

- **TypeScript/React web console**: gRPC-Web admin interface with i18n (EN/JA/ZH), Tailwind CSS, JSON viewer, sidebar navigation, and settings management
- **Per-service enable/disable configuration**: 32 services individually controllable via `services.<name>.enabled` config keys and gRPC RPCs (EnableService/DisableService)
- **Root user authentication**: LoginRoot, InitialSetup, IsRootInitialized gRPC methods for initial setup flow
- **Server metrics RPC**: GetServerMetrics in AdminConfigService for runtime performance data
- **Protobuf service definitions**: Added proto files for ACM, AdminAuth, AdminConfig, APIGateway, AppSync, Athena, CloudFront, CloudTrail, and more in `proto/aws/`
- **Dynamic port allocation**: New `internal/server/portalloc/` package for resolving and binding service ports dynamically
- **Storage manager**: Centralized global storage initialization in `main.go`

### Changed

- **ACM overhauled**: Account configuration persisted to PebbleDB (was in-memory, lost on restart); certificate listing uses `common.List` with `ListByStatus` for server-side pagination replacing handler-level manual scan; tag operations migrated to `tagutil` framework removing `GetTags`/`AddTags`/`RemoveTags` store methods; `certificateToDetailResponse` omits zero-value `NotBefore`/`NotAfter` instead of fabricating timestamps; tag error mapping returns `ResourceNotFoundException`/`ValidationException` instead of raw errors; admin handler uses `GenerateCertificateId` replacing `time.Now().UnixNano()`; removed `GetByDomain`, `DeleteAccountConfiguration`, and wrapper functions
- **All AWS service admin handlers**: Unified constructor signature with StorageManager and AccountID parameters
- **API Gateway overhauled**: Patch operation parsing consolidated into `parsePatchOperations` helper, eliminating duplicated type-assertion boilerplate across all Update operations; deep-copy helpers unified via generic `deepCopyMap[T]`; unused `*ListResult` types removed; admin store uses `GetOrCreateStoreE` singleton pattern; fixed nil map serialization in `Update` (`ensureRestApiMaps`), binary media type remove with JSON Patch `-` (append marker), and missing `rootResourceId` in REST API responses; `UpdateApiKey` now invalidates runtime cache via `RemoveApiKey`
- **AppSync overhauled**: Response serialization consolidated into `response_mapper.go`, parse helpers into `helpers.go`; `extractTypeName` unified via `graphql.ExtractTypeName` shared utility; `GetApiById`/`GetGraphqlApiById` use `#id:` secondary index replacing full-table scans; fixed `UpdateApiById`/`UpdateGraphqlApiById` Put/Delete ordering (data loss on rename), `ServiceQuotaExceededException` HTTP 402→429; `DeleteDomainName`/`DeleteGraphqlApiById` cascade to domain associations; error mapping unified via `mapStoreError` across all operations; removed `jsonUnmarshal` alias, `sharedHTTPClientOnce`, `UpdateChannelNamespaceTags`
- **Athena overhauled**: HAVING/DDL result construction deduplicated (`evaluateWhere` reuse, `emptyDDLResult` helper, `buildResultSetFromStoredTable` delegation); `DropDatabase` cascades to tables and table data; `DeleteWorkGroup` cascades to named queries, prepared statements, query executions, and results; query execution store adds `#wg:` secondary index for workgroup-scoped operations, replacing full-table scans; cleanup changed from `sync.Once` to 24h periodic with cascade result deletion; removed `AthenaStoresInterface`/`AthenaStore` aggregate, `ListQueryExecutionsWithPagination`, `GetDefaultWorkGroup`
- **CloudFront overhauled**: `UpdateDistribution` fixed WAF WebACL change detection — old config was overwritten before comparison so WAF association/disassociation never fired; `ListInvalidations` removed duplicate top-level fields outside `InvalidationList`; tag parsing deduplicated into `parseXMLTags` shared by `CreateDistributionWithTags` and `TagResource` (replacing 2 inline copies); `newCloudFrontError` wrapper removed in favour of direct `awserrors.NewAWSError` (10 sites); `generateETag` service method removed (already handled inside `UpdateWithLastModified`); `TrustedSigners` no-op assignment removed; `UpdateCachePolicy` now parses `ParametersInCacheKeyAndForwardedToOrigin` matching `CreateCachePolicy`; store `Create` methods on 4 policy stores add mutex; `AdminHandler` caches `DistributionStore` instead of recreating per call; `DistributionServer.getDistributionStore` returns cached store on error instead of nil store
- **CloudTrail overhauled**: Trail resolution consolidated into `ResolveTrail` replacing scattered `GetTrail`/`GetTrailByARN` conditionals; `resolveBool`/`applyTrailUpdates` unify bool parsing and field assignment for `CreateTrail`/`UpdateTrail`; `buildIndexKeys`/`applyIndexKeys` deduplicate index operations across add/remove/txn paths; `eventMatchesQuery` delegates to `protoMatchesQuery`; `listAllTrails` replaces `MaxItems:10000` with paginated iteration; `ListTrails` uses store-level Marker pagination; `GetEventByID` unified to always read proto; `EventCategory` added to Event proto message
- **CloudWatch alarm operations refactored**: Create-or-update logic deduplicated into `upsertAlarm` shared by `PutMetricAlarm`/`PutCompositeAlarm`; tag parsing extracted to `parseAlarmTags`; alarm type resolution and action history logging extracted to `resolveAlarmType`/`addAlarmActionHistory`; `alarmToResponse` uses `alarm.CreatedAt` as fallback replacing non-deterministic `time.Now()`; removed dead CWL types from cloudwatch store package
- **CloudWatch Logs overhauled**: subscription and metric filter evaluation deduplicated via `evaluateMetricFilters`/`deliverSubscriptionEvents` replacing duplicated HTTP and bus handler paths; destination persistence migrated from JSON to protobuf with dedicated proto message; `DeleteLogGroup` uses paginated stream deletion (was limited to 1000); `LogGroupClass` added to proto round-trip (was lost on persistence); admin `CreateLogGroup` uses `NewLogGroup` constructor setting `Region`/`AccountID`; `PutLogEvents` rejects empty events; `log_writer` tolerates `ErrLogStreamAlreadyExists` on auto-create; `ErrLimitExceeded` store error mapping added
- **Cognito Identity overhauled**: Login/provider map parsing consolidated into generic `parseMapParam`, removing `parseLogins` and `parseSupportedLoginProviders`; map-key extraction deduplicated via `formatLoginKeys`; fixed `GetOpenIdToken` variable shadow causing redundant store fetch; `DeleteIdentityPool` now cascades to developer identities and principal tag attribute maps; `NewIdentityPool` timestamps deferred to store layer; removed unused `generateIdentityID`
- **Cognito IDP overhauled**: Authentication flows (`USER_PASSWORD_AUTH`, `ADMIN_NO_SRP_AUTH`, refresh token) deduplicated into `authenticateUser`/`refreshAuthToken`/`respondToNewPasswordChallenge` shared methods; `findTokenByValue[T]` generic helper replaces three identical `GetTokenByValue` implementations; `tokenKey` unifies `refreshTokenKey`/`idTokenKey`/`accessTokenKey`; `applyUserPoolUpdates` consolidates `CreateUserPool`/`UpdateUserPool` field assignment; `parseIntParam` eliminates `getIntParam`/`getIntParamOK` duplication; `parseUserPoolAddOns`/`parseAccountRecoverySetting`/`parseUsernameConfiguration`/`parseDeviceConfiguration` parsers added; `DeleteUserPool` now cascades to users, groups, clients, tokens, and challenge sessions; pagination added to `ListGroups`, `ListUsers`, `ListUsersInGroup`, `AdminListGroupsForUser`, `ListUserPoolClients`; `SignUp` generates and persists confirmation code; `ConfirmSignUp` validates code with constant-time comparison and expiry check; hosted UI hardened against XSS via `html.EscapeString`; `handleTokenEndpoint` returns `expired_grant` for expired refresh tokens
- **DynamoDB overhauled**: GSI/LSI query logic deduplicated via `queryByIndex`; key condition extraction consolidated into `extractKeyConditionFromSchema`; sort logic unified with `sortItemsBySortKeyWithIndexDirection`; `DeleteTableCascade` now removes backups, exports, imports, and global table entries; transaction validation uses `NewTxn` interface method instead of type assertion; `ItemStore.Put` no longer mutates input attributes map; fixed `BatchWriteItem` table-not-found error handling, `CreateTable` TableClass persistence, backup duplicate name prevention, `DeleteAllForTable` error handling, and unparseable key handling in `BatchGetItem`; deterministic ordering for attribute definitions and GSI updates; scan operations now surface unmarshal errors
- **EC2 overhauled**: 4 Authorize/Revoke Ingress/Egress operations deduplicated into `modifySecurityGroupRules` helper with closure pattern; `ipRuleEquals` extends duplicate/revoke matching from protocol+port-only to full IpRanges/Ipv6Ranges/UserIdGroupPairs comparison; `parseIPRules` parses GroupId/GroupName/UserId/VpcId/Ipv6Ranges/Description (was GroupId+CidrIp only); `DescribeVpcs`/`DescribeSubnets`/`DescribeSecurityGroups` gain EC2 Filter.N support (`parseFilters` + typed match helpers); `DeleteVpc` checks for dependent subnets and security groups (`DependencyViolation`); `CreateSecurityGroup` rejects duplicate GroupName within VPC (`InvalidGroup.Duplicate`); `storeForRegion` extracts store lookup removing 4 duplicated paths; `ec2store.Tag` replaced with shared `types.Tag`; `getStringParam`→`request.GetStringParam`, `ParseTags`→`tags.ParseTagsWithPrefix`, `parseInt64`→`strconv.ParseInt`
- **EventBridge overhauled**: `eventBusToMap`/`ruleToMap`/`archiveToMap` consolidate response serialization; `mapStoreError` unifies ~20 store error checks; `resolveTaggableResource` supports 5 taggable resource types; `GetStoreForRegion` extracts store caching; Connection/API Destination split into dedicated files; `DeleteEventBus` cascades to rules→targets→tags→archives→events; `DeleteArchive` cascades to events; `ListRuleNamesByTarget` rewritten with full-scan + client-side pagination; fixed `RemoveTargets` missing rule check, `ListApiDestinations` duplicate fields; added duplicate-name guards on Create and in-use guard on `DeleteConnection`; `DescribeEventBus`/`DescribeRule` now return Tags; `strings.HasPrefix` replaces manual prefix matching; removed `getEventBusByArn`/`getRuleByArn`, admin `storageManager`
- **IAM service refactored**: Global IAM store singleton (`GetOrCreateGlobalStore`) prevents redundant policy seeding; attached and inline policy operations deduplicated via ops-struct pattern; store hoisting in instance profile responses; fixed `GetAccountAuthorizationDetails` pagination (was truncated at 1000) and `GetAccountSummary` ServerCertificates count (was hardcoded 0); credential report access keys now sorted deterministically
- **Kinesis overhauled**: Stream name/ARN resolution consolidated into `resolveStreamName` (7 operations); `PutRecord` shard selection and write atomic via `PutRecordWithShardSelection` (TOCTOU fix); `GetRecords` returns actual EncryptionType instead of hardcoded "NONE"; `ListStreams` uses store-level pagination replacing full-scan; `initShardIDCounter` prevents shard ID collisions on restart; `cleanExpiredIterators` throttled to 5-minute intervals; consumer/format helpers extracted; error mapping unified via `mapStoreError`
- **KMS overhauled**: UUID v4 key ID generation on `CreateKey` (was empty string); `DeleteAlias` authorization check; `Sign`/`Verify` algorithm validation against key's supported algorithms; key pair generation moved to HSM layer with `GenerateKeyPair` Backend method; `resolveKeyByKeyID` delegates to `resolveKey`; `GetPublicKey` reads algorithm lists from persisted key struct; grant response building deduplicated; fixed `ListRetirableGrants` nil-entry slice; `KeyUsageEncryptDecrypt` excludes ECC key specs (AWS spec compliance); `buildKeyMetadata` uses persisted `MultiRegionConfiguration`
- **Lambda overhauled**: Function and version configuration responses deduplicated via `configFields`/`buildConfigMap`; CORS config parsing consolidated with `applyCorsFields`; deprecated `UpdateAlias` and `AddPolicy` store methods removed in favor of atomic variants; fixed async invocations (`Invoke`/`InvokeAsync`) not passing resolved version, `SetProvisionedConcurrency` using hardcoded account ID, and host endpoint port formatting; `AddPermission` now rejects duplicate StatementId; `DeadLetterConfig` validates TargetArn, `TracingConfig` validates Mode
- **Neptune ecosystem overhauled**: Neptune events migrated to `ProtoStore[Event]`; `clusterToResponseMap` replaces JSON round-trip; parameter group handlers deduplicated; `OpenClusterEngine`/`clusterBucket` accept explicit region parameter. NeptuneData Gremlin/OpenCypher query status/list/cancel deduplicated into shared methods; admin gRPC handlers consolidated; loader cancellation refactored to per-job channels enabling concurrent independent cancellation; `nodeIDMap` instance field removed for per-load local maps; `traversalToSteps` shared between explain and nested traversal; ntriples loader checks cancel in scan loop. NeptuneGraph bulk import batching consolidated into `importBatcher` struct; periodic expired-task cleanup via `regionCleanups` sync.Map with 24h production interval (5m/30m in test mode); dead `isCreate` parameter, `parquetTypeToPb` branch, `sortedPropKeys` removed; `endpointToResponse` includes `graphId`; `resolveGraphIdentifier` deduplicates header lookup; `s.region`→`reqCtx.GetRegion()`; proto3 zero-value fields mapped to nil
- **Route 53 overhauled**: Health check and delegation set response construction deduplicated into `healthCheckToResponse` (4 sites) and `buildDelegationSetResponse` (3 sites); tag parameter parsing extracted to `parseResourceParams` removing duplicated validation; admin `CreateHostedZone` auto-creates NS/SOA records, converts proto VPC region enum via `protoVPCRegionToAWS`, sets `Region`, and always initializes `Config`; `ListResourceRecordSets` defaults `maxItems` to 300; `normalizeZoneNameForCreate` migrated to exported `NormalizeZoneName`; fixed `ChangeTagsForResource` double-tagging bug; removed unused store types and helpers
- **S3 service overhauled**: Store acquisition hoisted to handler layer; EventBridge notification support; SNS/SQS `MessageAttributes`; copy operations preserve source StorageClass; `CompleteMultipartUpload` ETag validation; fixed `ListObjectVersions` pagination, `ListMultipartUploads` key parsing, website routing rules, and chunked encoding bounds
- **Scheduler overhauled**: Target parsing deduplicated via `coerceToMap`/`parseTargetFromMap` replacing duplicated string/map branches; case-insensitive field access unified with variadic-key helpers (`getMapField`/`getFloatField`/`getBoolField`/`getSliceField`); AWS cron field conversion consolidated into `awsCronReplacements` map loop; empty-input default unified across all target types into `scheduleInput`; `getScheduleNameAndGroup`/`getListGroupName` extract duplicated parameter parsing; `getStoreForSchedule` extracts engine store lookup with cache-first optimization; store prefix checks simplified with `strings.HasPrefix`; `handleBusDelivery` adds engine nil guard
- **SESv2 overhauled**: 13 Put/Create/Update operations deduplicated via `updateConfigSet`/`updateEmailIdentity`/`putEmailIdentityPolicy` helpers with modifier closures; `generateDkimToken` switched from deterministic `time.Now().UnixNano()` to `crypto/rand`; `parseStoredTime` fallback changed from non-deterministic `time.Now()` to `0`; `ListEmailIdentities` derives `VerificationStatus` from DKIM/VerifiedForSending (was hardcoded `"SUCCESS"`); `PutEmailIdentityDkimSigningAttributes` now reads `SigningAttributesOrigin` (was no-op); `GetContactList` response includes `DisplayName`/`Description` on topics; duplicate `SuppressionReason` field removed; store sentinel errors propagated directly replacing `awserrors` wraps; dead `ContactListStore`/`ContactStore` types removed; AdminHandler caches store per region via `sync.Map`
- **Secrets Manager overhauled**: ARN→name secondary index (`#arn:` prefix) replaces full-table scans; `UpdateSecretVersionStage` stage rotation atomic via store-level `MoveStage` (race condition fix); `DeleteSecret`/`GetSecretVersionByStage` use `secret.VersionIDs` replacing full-table scans; password generation uses Fisher-Yates shuffle (was deterministic fixed positions); fixed `UpdateSecret` SecretString/SecretBinary overwrite, operator precedence in replication condition; `applySecretFilters` deduplicated via `filterSecrets`/`matchesAny`; `List`/`ListProto` skip `#`-prefixed index keys
- **SNS overhauled**: Message attribute parsing consolidated into `parseMessageAttributes` with `firstString` case-insensitive key lookup and Binary attribute base64 decode (was dead `[]byte` cast); `messageAttributeValue` helper ensures Binary values are base64-encoded in SQS payload, HTTP notification, and Lambda event envelope; SNS→SQS delivery uses typed `SQSMessageAttribute` (`DataType`/`BinaryValue`/`StringValue`) preserving attribute types through the event bus; `Subject` and `MessageAttributes` added to `SNSDeliveryEvent` across Publish, PublishBatch, and PublishToTopic; subscription list construction deduplicated into `buildSubscriptionList`; permission injection extracted to `injectPermissionsIntoPolicy`; `parseAttributes` case-insensitive key search consolidated into loop; endpoint attribute merge extracted to `mergeEndpointAttrs`; `buildNotificationPayload` unifies Lambda/HTTP delivery; `ErrTopicNotFound` package-level variable replaces `NewTopicNotFoundException`; unused `TopicListResult`/`SubscriptionListResult` removed; `GetTopicAttributes` returns `KmsMasterKeyId`, `EffectiveDeliveryPolicy`, `CreatedDate`, `LastModifiedTime`; `SetTopicAttributes` syncs attributes to struct; fixed `ListEndpointsByPlatformApplication` marker inclusion and signature string composition
- **SQS overhauled**: FIFO ReceiveMessage with per-MessageGroup in-flight isolation; `SQSInvoker.SendMessage` signature changed to `SQSSendOptions` struct; `SetQueueAttributes` expanded to `FifoQueue`, `ContentBasedDeduplication`, `RedrivePolicy`; `ChangeMessageVisibility` and DLQ move now use atomic transactions
- **SSM overhauled**: List parameter parsing consolidated into `parseStringList` helper (5 call sites); `GetParametersByPath` path normalized with trailing `/` suffix (fixes `/foobar` false-positive prefix match); admin `PutParameter` supports `KeyID`, `Tier`, and `Tags`; `DescribeParameters` Type/Tier hardcoded defaults removed in favor of existing switch mapping; `DeleteParameter` history deletion paginated (MaxItems 1000 loop replacing single 10000 batch); `DescribeParameters` gains `DataType` filter; `GetParametersByPath` non-recursive path matching simplified; fixed `PutParameter` tag guard dropping tags on overwrite
- **Step Functions overhauled**: `DeleteStateMachine` cascades to executions→history events, versions, and aliases via `ForEach` full-scan with collect-then-delete pattern; `DeleteActivity` cascades to pending tasks with channel close; activity queue drain now closes channel to prevent goroutine leaks; `deleteExecutionHistory` extracted as helper; fixed `StartSyncExecution` removing incorrect EXPRESS rejection (was backwards—EXPRESS is the type that supports sync); `ListMapRuns` omits empty `nextToken`; `extractStateMachineName` replaced with `arnutil.ExtractStateMachineNameFromARN` for proper ARN parsing; `getStringOrEmpty` deduplicated to `getStr`; `isErrNotFound` uses `errors.Is` sentinel errors
- **STS and admin_auth overhauled**: `resolveRoleForAssume` deduplicates role resolution across AssumeRole/AssumeRoleWithSAML/AssumeRoleWithWebIdentity; ARN format fixes in `GetCallerIdentity` (SAML, FederatedUser, WebIdentity); SAML assertion and WebIdentity token validation; `GetDelegatedAccessToken` uses delegated token store with single-use redemption; admin_auth JWT generation uses `vsjwt.JWTUser` directly, removing type-switch duplication
- **Timestream Query and Write overhauled**: Scheduled query response formatting deduplicated into `formatScheduledQueryBaseResponse` with list/description overrides replacing two independent builders; `ScheduledQueryState` type consolidated into `ScheduledQueryStatus` (same values); `UpdateScheduledQuery` store switch matches only valid states (was silently treating invalid as DISABLED); `GetQueryStatus`/`ListQueryResults` dead methods and `CachedRows` redundant field removed; `ErrInvalidParameter` merged into `ErrValidationException` (13 call sites); `mapStoreError` uses `errors.Is` sentinel errors replacing string comparison; `ListTables` validates empty database name; store `writeChunkBatch` returns `[]error` for partial version conflict reporting; `BatchLoadTaskStore` unused `TagStore` embed, `TimestreamStoresInterface`/`TimestreamStores` aggregate, and `NewRecordStore` non-indexed variant removed
- **WAFv2 overhauled**: Resource summary builders deduplicated into `buildWebACLSummary`/`buildRuleGroupSummary`/`buildIPSetSummary`/`buildRegexPatternSetSummary` (10 sites); IP address and regex list parsing extracted to `parseAddressList`/`parseRegularExpressionList`; `convertAction` handles lowercase JSON-deserialized keys; `convertActionToResponse` supports Captcha and Challenge actions; `UpdateRuleGroup` supports partial updates; `PutLoggingConfiguration` preserves `RedactedFields` as pass-through; `ListLoggingConfigurations` uses pagination helpers; `LoggingStore.Create` adds mutex; `Get` replaced by `GetByResourceArn`; removed unused store types and fields
- **CI/CD**: Added Node.js 22 setup, webconsole dependency install, and proto generation steps
- **Cypher parser refactored**: Edge target resolution deduplicated into `resolveEdgeTarget` (replacing 8 inline blocks across 6 files); result modifier logic (ORDER BY/SKIP/LIMIT/DISTINCT) extracted to `applyResultModifiers` (5 call sites); property resolution consolidated from 5 near-identical functions into `resolvePropsMap`/`resolvePropValue` with thin wrappers; `crossProduct` renamed to `crossProductBindings`; fixed MATCH WHERE silently ignored when WITH clause present (`matchWhere`/`matchVarLengthBindings` removed incorrect `q.With` guards); dead `evalRowExpr` wrapper and `OpNeq` unreachable branch in `evalComparison` removed; edge property check moved before row allocation in optional match
- **Gremlin parser refactored**: Edge traversal step execution deduplicated into `execEdgeTraversal` (replaces identical `execOutE`/`execInE`/`execBothE` bodies); `execFrom`/`execTo` merged into `execFromTo` with `tagKey` parameter; `execGroup`/`execGroupCount` share `groupMods` and `groupBy` generic grouping helper; `execProperty` Node/Edge update paths consolidated via `elementIDForUpdate`/`resolvePropertyValue`; `execElementMap` map construction extracted to `buildElementMap`; `isSlice` simplified from `reflect` to `[]any` type assertion; `resolveArgValue` non-string map key handling changed from returning error value to `continue` (was unreachable error path); fixed `property()` step `set` cardinality writing redundant updates when value already present (idempotency restored via `valuesEqual` guard)
- **Dependencies**: Added gopsutil (indirect: purego, go-ole, plan9stats, perfstat, go-sysconf, numcpus, wmi)
- **Documentation**: Updated README (EN/JA/ZH) with GUI screenshot and Docker requirements
- **Makefile**: Build now generates proto TypeScript before webconsole build; binary name updated

### Removed

- **Flutter/Dart web UI** (`web_ui/` directory deleted)
- `internal/common/defaults/` package (superseded by `serviceports`)
- `ConfigSource` type and `Source` field from config entries
- `getEnvString`, `getEnvInt`, `getEnvBool` helper functions from config defaults

## [0.0.10] - 2026-04-27

### Added
- **Temporary session credentials support** (`internal/server/http/auth_middleware.go`): AWS Signature middleware verifies both static and temporary (STS) credentials via `SessionResolver`
- **S3Invoker for EventBus** (`internal/eventbus/service_invokers.go`): cross-service S3 object operations enabling DynamoDB export/import to S3
- **DynamoDB Step Functions integration**: `GetItem`, `PutItem`, `DeleteItem`, `UpdateItem` task support in Step Functions service
- **CloudFront ListKeyGroups operation** and Location header handling in responses
- **IAM MFADevice name requirement** and server certificate rename support
- **Resource tracking in audit events**: `ResourceTypes` field and `buildResources` function for CloudTrail event enrichment
- **IAMPrincipalResolver**: access key ID to username resolution for audit logging
- **Defaults package** (`internal/common/defaults/`): centralized `DefaultRegion` constant replacing scattered references
- **SDK tests massively expanded** (290 files, +66,464/-44,016 lines): full lifecycle coverage for ACM, API Gateway (14 test files), AppSync (13 files), Athena (8 files), CloudWatch/Logs, Cognito, DynamoDB, EventBridge, KMS, Neptune, S3, SESv2, SQS, WAFv2 (5 files), and more — now 2262+ tests total

### Changed
- **awscli full-service validation pass**: comprehensive awscli-driven testing across all services uncovering and fixing integration issues
- **DynamoDB import/export refactored** (`import_export_operations.go`, 260 lines): S3-backed export/import with improved error handling and validation
- **KMS key operations refactored**: `CreateKey` handles external-origin keys (PendingImport state), `ImportKeyMaterial` validates tokens and decrypts wrapped keys
- **SESv2 configuration handling enhanced**: VDM attributes, contact list topic handling, configuration set event destination management
- **API Gateway service hardened**: conflict detection in `CreateResource`, authorization type validation in `PutMethod`, `CreateRestApi` Policy parameter, cascading `UpdateResource`, streamlined `TestInvokeMethod`
- **Cognito user pool configuration expanded**: email/SMS configuration parsing, new user pool attributes and settings
- **AppSync introspection and associations fixed**: `ListTypesByAssociation` validation, `ListSourceApiAssociations` merged API references, `StartDataSourceIntrospection` consistent failure response
- **Athena service hardened**: resource-not-found errors for catalog/workgroup/prepared statement queries, pagination in `ListDataCatalogs`/`ListPreparedStatements`, SELECT without FROM support
- **CloudWatch metric response formatting** and log group operations improved
- **State machine definition JSON validation** added, duplicate creation error handling
- **Audit event builder enhanced**: S3 event determination logic, event source mappings, event filtering by `EventID`/`ReadOnly`/`EventSource`/`AccessKeyID`
- **SDK test utilities refactored**: monolithic test files split into focused modules (e.g., `apigateway.go` → 14 files, `appsync.go` → 13 files, `athena.go` → 8 files)

### Fixed
- S3 encryption handling
- DynamoDB condition expression handling
- API Gateway error responses with specific not-found exceptions
- SQS policy generation in `GetQueueAttributes`
- SSM parameter label management (`UnlabelParameterVersion`), validation error handling, `PutParameter` tags silently dropped on overwrite, `GetParametersByPath` prefix false-positive match
- WAF duplicate resource error messages
 - OIDC provider ARN construction for URL prefixes
- SESv2 response structure inconsistencies
- Region handling edge cases across services

## [0.0.9] - 2026-04-26

### Added
- **EC2 resource stubs** (new): VPC, Subnet, Security Group data models and store layer for cross-service resource references (`internal/services/aws/ec2/`, `internal/store/aws/ec2/`) — not a functional compute service
- **vstacks CLI** (`cmd/vstacks/`): management tool for server control, IAM, config, service, and backup operations
- **PebbleDB batch operations**: `CommitSync` method for atomic writes; batch interface (`internal/core/storage/batch.go`, `pebble_batch.go`)
- **Goroutine pool** (`internal/core/resilience/goroutine.go`): bounded goroutine execution for resilience
- Unit tests: AppSync schema operations & WebSocket server, DynamoDB item condition/expression/key condition/table updater, Neptune pagination & global/event operations, store tag tests (362 lines)

### Changed
- **Neptune/NeptuneGraph enhanced**: graph engine moved from `pkg/graphengine/` to `internal/core/storage/graphengine/` with expanded KV-backend store (845+ lines tests); service handler split (`service.go` → export/import/queries/snapshots/tags/responses); admin operations refactored (cluster, subnet group, error handling); Cypher parser DDL support (CREATE/DROP index) and modularization into ~10 files; Gremlin parser modularized into 5 step files + helpers; N-Triples serialization added (`internal/utils/ntriples/`); EventBus NeptuneGraph invoker; NeptuneGraph parser and error handling improved
- **Request parsing refactored** (`internal/common/request/`): monolithic parser simplified via registry pattern; per-service parsers added; `interfaces.go` removed
- **Tag handling refactored** (`internal/common/tags/`): operations extracted into dedicated handler with `UntagResource` by keys; removed `internal/common/types/tag.go`
- **Service handler file splits**: API Gateway integration (`aws.go` → `aws_dynamodb.go` + `aws_sns.go` + `mapping.go`), S3 (`handler_bucket.go` → put/get/delete/response)
- **Lambda ESM poller and tag operations refactored**: improved stream polling and streamlined tag management
- **AppSync refactored**: evaluation operations expanded (+507 lines), datasource simplified (−300 lines), tag operations enhanced (+180 lines)
- **EventBus concurrency handling enhanced** (`bus.go`, `invoker.go`)
- **PebbleDB improvements**: TTL handling, transaction management, iterator, and lock mechanisms
- **Blob storage refactored** (`blob_hybrid.go`)
- **Docker (Moby) client refactored** (`internal/client/mobyclient/`): container, network, and volume operations with dedicated error types (`errors.go`)
- **Server wiring refactored** (`internal/server/apps/`): adapters, optional services, and service initialization streamlined
- **Resource extractor refactored** (`internal/server/authorization/resource_extractor.go`)
- Pagination helpers cleaned up, XML encoding improved, region handling simplified
- **SDK tests enhanced**: expanded coverage for SNS, SQS, Secrets Manager, Step Functions (+1661/-527 lines)

### Fixed
- Adaptive timeout redundant return statement (`internal/core/resilience/adaptive_timeout.go`)
- Error handling and unused code cleanup across multiple services (AppSync, NeptuneGraph, KMS, S3, SESv2, Route53, WAF, etc.)
- Lambda redundant response header handling

### Removed
- **Utility package consolidation** (~4270 lines): deleted `internal/utils/archiver/`, `internal/utils/naming/`, `internal/utils/netutils/`, `internal/utils/timeutils/`, `internal/utils/aws/conditions/`, `internal/utils/aws/helpers/`, `internal/utils/aws/parsers/`, `internal/utils/crypto/` (consolidated into core packages)
- Deleted `internal/common/types/tag.go`, `internal/common/request/interfaces.go`, `internal/store/aws/common/types.go` (absorbed into other packages)
- Deleted service-specific operation files absorbed by refactoring: Cognito user operations, Lambda function operations, Secrets Manager secret operations

### Documentation
- Updated README (English, Japanese, Chinese) with latency metrics and standardized `DATA_PATH`
- Updated configuration reference (`docs/configuration.md`)
- Updated architecture and services documentation
- Updated LocalStack comparison report

## [0.0.8] - 2026-04-09

### Added
- Neptune Graph (Neptune Analytics) service with graph store, RDF/SPARQL support, vector embeddings, topK search, and query procedures (`neptunegraph/`, 2265+ lines)
- Neptune Graph SDK tests (895 tests) with host prefix middleware workaround
- Graph engine vector embedding storage with in-memory cache, cosine/Euclidean/inner-product distance functions, and brute-force topK search (`pkg/graphengine/vector.go`, 771 lines)
- Lambda AWS Event Stream binary encoding support (`lambda/eventstream.go`) for streaming invoke responses
- Test helper scripts: `run_tests.sh`, `setup_test_credentials.sh`, `cleanup_test_resources.sh`, `test_authorization.sh`
- NeptuneGraph request parser (`internal/common/request/neptunegraph_parser.go`)
- App wiring layer (`internal/server/apps/`) for modular service initialization and dependency injection
- Authorization module extracted from dispatcher (`internal/server/authorization/`) with enhanced resource extraction for NeptuneGraph
- `internal/common/request/context.go` as shared request context (replaces per-service copy)
- Unit tests for endpoint builder, error factories, handler registrar, tags operations, resilience (adaptive timeout, bulkhead, circuit breaker, retry, cache, health, metrics), PebbleDB, auth credentials, eventstream encoder, and validators

### Changed
- Common packages moved from `internal/services/aws/common/` to `internal/common/` (auth, endpoint, errors, iam, kms, lambda, mock, pagination, protocol, region, request, response, tags, types, audit)
- Event bus moved from `internal/server/eventbus/` to `internal/eventbus/`
- All 26+ service implementations updated to import from new `internal/common/` paths
- `main.go` restructured: service initialization delegated to apps wiring layer (675-line refactor)
- All services now receive dependencies via setter injection (removing direct storage manager coupling)
- DynamoDB error constants expanded with comprehensive godoc documentation; removed unused expression files (`item_condition.go`, `item_expression.go`, `item_sort.go`, `partiql_expression.go`, `partiql_value_parser.go`, `input_output.go`)
- Neptune service expanded: cluster, instance, snapshot, parameter group, and subnet group operations with pagination support
- Neptune Data service enhanced with query status tracking, statistics management, and improved error handling
- S3 SSE-S3 encryption support, enhanced access control, and updated chunked upload handling
- Lambda invoke operations enhanced with response streaming and event stream support
- Cypher parser extended with CALL statement support and pipeline execution
- Gremlin parser enhanced with additional filter and source step operations
- AppSync GraphQL datasource simplified, WebSocket server improved
- CloudFront distribution server and policy operations enhanced
- CloudWatch alarm evaluator multi-region support with configurable tick intervals
- Scheduler engine enhanced with Step Functions target and configurable intervals
- Step Functions enhanced with JSONata evaluation, Map state processor, and redrive operations
- EventBridge refactored to multi-region store model with improved archive and replay operations
- PebbleDB transaction handling, iterator, TTL, and bucket operations improved
- Graph engine store and traversal implementations optimized
- Filter pattern evaluator and parser improved
- SDK tests updated for CloudWatch, Cognito, Lambda, NeptuneData, SESv2, SNS, SQS, StepFunctions, STS, Timestream, WAFv2; added NeptuneGraph test registration
- Updated documentation: architecture, configuration, services, Terraform guide, LocalStack comparison, README files
- Dependabot: bumped AWS SDK v2 dependencies (lambda, cloudwatchlogs, kinesis, s3, eventstream, OpenTelemetry)

### Fixed
- CloudTrail recorder import path corrected after audit package relocation
- SNS publish ARN resolution for event bus delivery
- Various store interface and type consistency fixes across services

## [0.0.7] - 2026-04-05

### Added
- AppSync service (GraphQL API v1, Event API v2) with full CRUD for APIs, data sources, resolvers, functions, types, schemas, API keys, caches, domain names, and merged APIs
- AppSync GraphQL engine with VTL resolver execution, introspection, pipeline resolvers, and `$ctx` context variables (`$ctx.args`, `$ctx.source`, `$ctx.identity`, `$ctx.info`, `$ctx.stash`, etc.)
- AppSync WebSocket server for real-time subscriptions and mutations
- AppSync request parser supporting both v1 (GraphQL) and v2 (Event API) protocols
- VTL engine AppSync context extensions (`pkg/vtl/appsync_context.go`, `appsync_util.go`)
- WAF gRPC-Web proto definitions and generated code
- Gremlin parser enhancements: comprehensive AST, lexer improvements, new filter steps (`where`, `has`, `select`, `order`, `group`, `dedup`, `path`), and source steps (`V`, `E`, `addV`, `addE`)
- Event bus resource extractors for AppSync and Neptune (dispatcher audit logging)
- Integration test suite (1748 lines) covering cross-service event flows: EventBridge→Lambda/SQS/SNS/Kinesis/StepFunctions, ESM SQS/Kinesis→Lambda, CloudWatch Alarm→SNS/Lambda/StepFunctions, Scheduler→Lambda/SQS/SNS/StepFunctions, SFN Task→Lambda/SQS/SNS, S3 Notification→Lambda
- SDK test service registry with category support (SDK, WebSocket, Integration)
- SDK tests for AppSync (2135 lines) and AppSync WebSocket (740 lines)
- Additional SDK test coverage for DynamoDB (backup), KMS, Lambda, and other services

### Changed
- Proto packages renamed to match AWS service naming conventions (`rds→neptune`, `email→sesv2`, `states→sfn`, `monitoring→cloudwatch`, `events→cloudwatchevents`, `logs→cloudwatchlogs`, etc.)
- Neptune admin handler expanded with full DescribeDBClusters/Instances and gRPC-Web management operations; proto migrated from RDS package
- Step Functions now subscribes to event bus for cross-service start execution events (EventBridge, Scheduler, CloudWatch Alarms); added Map state `ItemProcessor` and enhanced parallel execution
- EventBridge refactored to multi-region store support (`sync.Map` per-region stores); added Kinesis stream target delivery
- Lambda ESM poller enhanced with Kinesis stream source support (shard iterator, checkpoint tracking)
- CloudWatch alarm evaluator extended to multi-region evaluation with TEST_MODE 1-second tick interval
- EventBridge Scheduler engine enhanced with Step Functions target support and TEST_MODE 1-second ticker
- SNS/SQS admin handlers refactored to delegate to shared service instances instead of owning storage
- SNS publish operations now resolve SNS topic ARN components for event bus delivery
- AppSync ARN builder added to `internal/utils/aws/arn/builder_services.go`

### Documentation
- Updated README files (English, Japanese, Chinese) with AppSync service and updated test counts
- Updated SDK test README with integration test documentation
- Updated ACM proto with expanded field coverage

## [0.0.6] - 2026-04-04

### Added
- Event bus system (`internal/server/eventbus/`) with Pebble outbox store, IAM role resolver, policy evaluator, and subscription management for cross-service event routing
- HTTP request classifier (`internal/server/http/classifier/`) replacing legacy ServiceRouter — classifies requests by protocol (REST-XML, REST-JSON, AWS JSON, Query, CBOR) and service
- CloudWatch alarm evaluator engine extracted from admin handler (`alarm_evaluator.go`, 629 lines)
- Lambda Event Source Mapping (ESM) poller for DynamoDB/Kinesis streams (`esm_poller.go`)
- Secrets Manager rotation engine with configurable rotation strategies (`rotation_engine.go`)
- Cognito User Pool trigger pipeline (Pre/Post Sign-up, Sign-in, MFA, Token, Auth challenges — `triggers.go`)
- S3 event notifications (ObjectCreated, ObjectRemoved, ObjectRestore) published via event bus (`notifications.go`)
- CloudWatch Logs subscription filter delivery via event bus (`bus_handlers.go`, `log_writer.go`)
- Neptune Data API: Gremlin, OpenCypher, loader, statistics, and explain handlers refactored into dedicated files
- PebbleDB iterator extracted into `db_iter.go`, lazy iterator (`db_lazy_iterator.go`), and TTL support (`db_ttl.go`)
- Blob storage split into `blob_multipart.go`, `blob_reader.go`, `blob_versioning.go`
- AWS SigV4 credential parser utility (`internal/utils/aws/authutil/credential.go`)
- Common request helpers: `GetBoolParam`, `GetIntParam` in `internal/services/aws/common/request/`

### Changed
- Dispatcher refactored: extracted `executeHandler` pipeline, added `DispatchClassified` for classifier-based routing
- Admin handlers across 20+ services refactored: gRPC stubs removed from Neptune Data, Kinesis, SNS, SQS, EventBridge; large monolithic handlers trimmed
- gRPC-Web error constructors renamed from `ErrXxx()` to `NewXxxError()` for consistency
- NeptuneData error type simplified — removed wrapper struct, returns `*awserrors.AWSError` directly
- Service index builder extracted (`router/service_index.go`)
- Removed legacy files: `request_extraction.go`, `router/service_router.go`, `router/path_patterns.go`, `core/events/bus.go`
- Removed docs: `docs/integration.md`, `docs/new-service-guide.md`
- IAM role resolver wired into event bus at startup for trust policy evaluation

### Documentation
- Updated architecture doc (`docs/architecture.md`)
- Updated LocalStack comparison report (`docs/localstack_vs_vorpalstacks_report.md`)

## [0.0.5] - 2026-04-01

### Added
- Neptune graph database service (property graph + RDF, openCypher/Gremlin, bulk loader)
- Neptune Data API (Neptunedata) with routing and request handling

### Changed
- Updated gRPC dependency from v1.79.1 to v1.79.3
- Service count updated to 30, Go SDK tests to 890 passing

### Documentation
- Added Japanese and Chinese READMEs
- Added performance benchmark section
- Added Terraform conformance tests and LocalStack comparison links
- Updated roadmap: Neptune promoted to implemented

## [0.0.4] - 2026-03-30

### Added
- Step Functions JSONata query language support (full AWS specification compliance)
- Workflow variables with scope management (Assign field on all state types)
- Built-in variables: `$states.input`, `$states.result`, `$states.context`, `$states.errorOutput`
- JSONata Output, Arguments, Items, Condition, Seconds expressions
- TestState Variables support and VariableReferences on DescribeStateMachine
- EvaluationFailed history events for JSONata query errors
- JSONata-only / JSONPath-only field validation on Create/UpdateStateMachine
- Custom AWS JSONata functions: `$uuid`, `$hash`, `$random`, `$parse`, `$partition`, `$range`
- Context object (`$states.context`) matching `$$.` intrinsic function structure
- HeartbeatSeconds / TimeoutSeconds expression evaluation in JSONata mode
- MapItemValue / MapItemIndex for Map state ItemSelector
- Catch.Assign and Catch.Output for JSONata error handling
- Lambda Function URL server, S3 website hosting, CloudFront distribution server
- Cognito hosted UI server, listener manager, gRPC-Web admin_auth registration
- Lambda and new service guide documentation

### Fixed
- InputOutput unmarshal for string values in JSONata definitions
- Choice Condition unwrapping for JSONata expressions
- Variable names incorrectly stored with `$` prefix
- VariableReferences regex, builtin filtering, and describe-time recomputation
- Double-wrapping of `states` key in variable map construction
- BuiltinFunction value leaking into JSONata expression results
- Various service log/request refactors

## [0.0.3] - 2026-03-29

### Added
- Route53 path-based routing (create/get/delete by zone name)
- Route53 NS and SOA record auto-creation on hosted zone creation
- WAF ListRules and ListRateBasedRules APIs
- XML namespace encoder for REST-XML protocols
- Integration test scripts for API Gateway → Lambda → DynamoDB and S3

### Fixed
- CloudWatch millisecond timestamps converted to epoch seconds in JSON responses
- Kinesis SubscribeToShard initial-response event and nested JSON parameter parsing
- S3 URL-decoded copy-source header
- Route53/CloudFront XML list encoding (nil-slice → empty elements)
- SSM DescribeParameters filter by type/path
- CloudFront response headers policy, invalidation, origin access control, and cache policy list operations
- Go SDK tests 594/594, Python 605/605, TypeScript 603/603, C# 622/622 — all passing

### Documentation
- Updated API coverage stats (SecretsManager 100%, StepFunctions 100%, KMS 90%)

## [0.0.2] - 2026-03-28

### Added
- IAM Service Last Accessed Details with CloudTrail integration
- IAM policy condition support (StringEquals, StringNotEquals, Bool, ArnEquals, ArnNotEquals, IpAddress, NotIpAddress, DateLessThan, DateGreaterThan, Null)
- IAM SimulatePrincipalPolicy API
- IAM EvaluateTrustPolicy with Federated principal matching
- KMS authoriseOperation() wired into all KMS handlers (~30 call sites)
- STS EvaluateTrustPolicy wired into AssumeRole, AssumeRoleWithSAML, AssumeRoleWithWebIdentity
- Secrets Manager replication APIs (ReplicateSecretToRegions, UpdateReplicationInfo, RemoveRegionsFromReplication)
- Step Functions RedriveExecution and TestState APIs
- CloudFront ListDistributionsByWebACLId API
- WAFv2 GetWebACLForResource API
- CloudFront-WAF association synchronisation
- BootstrapConfig for centralised environment variable loading

### Changed
- Replaced 14 multi-constructor variants with setter injection across 7 services (SNS, EventBridge, Step Functions, Scheduler, CloudWatch Logs, Lambda, API Gateway Runtime)
- Extracted gRPC-Web admin handler registration to `grpcweb.RegisterAllAdminHandlers()`
- Moved Cognito JWKS handler from main.go closure into `CognitoService.JWKSHandler()`
- Split large files for SRP: `s3/handler.go` → 4 files, `dynamodb/convert.go` → 4 files, `kinesis/store.go` → 5 files
- Replaced service operation counts with Full/Broad/Selective coverage tiers in documentation
- main.go reduced from 502 to 345 lines (31% reduction)

### Fixed
- IAM trust policy parsing for single-string StringList values
- IAM SLA nil dereference, XML key names, audit principal resolution
- API Gateway Runtime setter injection order (both stores now mutate same factory)
- WAF and WAFv2 now use regional storage instead of global

### Documentation
- Added LocalStack comparison report
- Refreshed services.md with coverage tiers

## [0.0.1] - 2026-03-27

### Added
- Initial public beta release
- 30 AWS services implemented (S3, SQS, SNS, Lambda, DynamoDB, IAM, KMS, Kinesis, API Gateway, Step Functions, EventBridge, CloudFront, WAF, WAFv2, Cognito, STS, Route53, CloudWatch, CloudTrail, Secrets Manager, SSM, Scheduler, SESv2, Athena, Timestream, ACM)
- 594 SDK integration tests (100% pass rate)
- gRPC-Web admin console (Flutter)
- Multi-region support with PebbleDB storage
- IAM policy evaluation
- Docker-based Lambda execution
- Terraform and OpenTofu compatibility guide
- GitHub Actions CI

### Documentation
- Architecture overview
- Service reference with operation counts
- Configuration reference
- Terraform & OpenTofu guide with 18 verified services
- Contributing guidelines

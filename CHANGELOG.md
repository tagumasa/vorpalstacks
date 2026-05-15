# Changelog

All notable changes to Vorpalstacks will be documented in this file.

## [Unreleased]

### Added

- **Neptune: TinkerPop Gremlin Server WebSocket protocol** — Neptune clusters now expose a native WebSocket endpoint at `/gremlin` supporting the TinkerPop Gremlin Server binary protocol with GraphSON v3 serialization. Clients can connect directly using Gremlin language drivers (e.g. `gremlinpython`, `Gremlin-JavaScript`) for real-time traversal evaluation against per-cluster graph engines.

- **Neptune: Per-cluster data plane listeners** — Each Neptune DB cluster automatically opens an isolated graph engine with a dedicated HTTP listener on a dynamically allocated port. The data plane supports both REST (OpenCypher/Gremlin HTTP) and WebSocket protocols. Engines are restored automatically on server restart.

- **Neptune: GraphSON v3 encoder/decoder** — Full GraphSON v3 type system support for TinkerPop-compliant serialization of vertices, edges, paths, maps, lists, and scalar types over WebSocket.

- **Integration tests: Neptune direct protocol tests** — New test suite exercising raw HTTP POST and WebSocket paths bypassing the AWS SDK, covering OpenCypher queries, Gremlin HTTP queries, Gremlin WebSocket sessions, and GraphSON v3 response structure validation.

- **RDS: embedded MySQL engine via go-mysql-server** — New `rdbengine` package (`internal/core/storage/rdbengine/`) provides an embedded relational database engine backed by Pebble KV storage with row-level CRUD, secondary indexes, unique constraints, catalog management, and type-aware column encoding for MySQL compatibility. The `vmysql` service adapter exposes each RDS MySQL instance as an isolated `go-mysql-server` engine with a dynamically allocated TCP listener, enabling direct MySQL protocol connectivity from any MySQL client. Enabled via `RDS_MYSQL_ENABLED=true`.

- **RDS: shared store layer for RDS-compatible services** — New `internal/store/aws/rds/` package with protobuf-backed generic stores for `DBCluster`, `DBInstance`, `DBClusterSnapshot`, parameter groups, subnet groups, global clusters, event subscriptions, events, and tags. Neptune store refactored to embed `RDSStore` and export type aliases, enabling reuse by future RDS engines (MySQL, PostgreSQL). `rds.Engine` and `rds.GetPorter` interfaces decouple engine lifecycle management from specific service implementations.

- **RDS Data API: SQL execution via go-mysql-server** — New `rdsdata` service package (`internal/services/aws/rds/rdsdata/`) implements 6 RDS Data API operations (`ExecuteStatement`, `BatchExecuteStatement`, `ExecuteSql`, `BeginTransaction`, `CommitTransaction`, `RollbackTransaction`) that execute SQL against vmysql instances through the go-mysql-server `sqle.Engine`. Supports named parameter substitution, column metadata, `formatRecordsAs=JSON`, transaction management with 5-minute TTL and background stale-transaction reaper. Routed via `rds-data` signing service name and root-path operation mapping (`/Execute` → `ExecuteStatement`, etc.). Enabled when `RDS_MYSQL_ENABLED=true` (shares the RDS MySQL gate). SDK tests: 16 cases covering CRUD, transactions, batch, multi-statement, and error paths.

### Changed

- **Neptune service packages relocated under `rds/`** — Neptune, NeptuneData, and NeptuneGraph service implementations moved from `internal/services/aws/{neptune,neptunedata,neptunegraph}/` to `internal/services/aws/rds/{neptune,neptunedata,neptunegraph}/`. Neptune and NeptuneGraph stores moved from `internal/store/aws/{neptune,neptunegraph}/` to `internal/store/aws/rds/{neptune,neptunegraph}/`. All import paths updated accordingly.

- **Dispatcher: removed resilience wrapper** — The HTTP request dispatcher no longer wraps handler invocations in the resilience (circuit breaker, bulkhead, adaptive timeout) layer. The `resilience.ServiceResilienceConfig` parameter has been removed from `NewDispatcher` and `NewServer`.

### Fixed

- **S3: error handling for malformed XML, invalid part-number-marker, and unparseable redirect/routing codes** — Previously silently ignored parse errors; now returns proper `MalformedXML`/`InvalidArgument` (400) responses.

- **Neptune: operational fixes for cluster modification, role removal, tag cleanup, and parameter ordering** — `ModifyDBCluster` now syncs `Endpoint.Port` on port change, creates new cluster before reparenting resources on rename (with rollback on failure). `RemoveRoleFromDBCluster` supports `FeatureName` filtering. `RemoveFromGlobalCluster` matches by cluster ID extracted from ARN. `CreateDBInstance` applies tags on creation. Delete operations for instances, snapshots, parameter groups, and subnet groups now cascade tag cleanup. Parameter lists and tags are sorted deterministically. NeptuneGraph `Close`, `DeleteGraph`, `StopGraph`, and `ResetGraph` fixed to release mutex before calling `wg.Wait`/`db.Close`, preventing deadlock when held during blocking operations. NeptuneDataService `Close` renamed to `Shutdown` and refactored to collect engine entries under lock then close outside it.

- **API Gateway: deadlock fix in cascading deletes and missing lock coverage** — `DeleteDomainName` and `DeleteUsagePlan` previously called back into `DeleteBasePathMapping`/`DeleteUsagePlanKey` (which acquire the same mutex), causing deadlock. Introduced `deleteBasePathMappingLocked`/`deleteUsagePlanKeyLocked` internal methods to avoid re-acquisition. Added missing mutex protection to `UpdateBasePathMapping`, `DeleteBasePathMapping`, `UpdateUsagePlan`, `DeleteUsagePlanKey`, and several `RestApiStore` mutation methods.

- **AppSync: pipeline resolver before-template result propagation and rename collision** — Pipeline resolvers now parse and forward the before-step `RequestMappingTemplate` output to downstream functions via `$ctx.prev.result` (previously discarded). `UpdateApiById` and `UpdateGraphqlApiById` now reject name changes that would collide with an existing API. `dispatchRDS` implemented: RELATIONAL_DATABASE data sources now execute SQL through the EventBus `RDSDataInvoker` instead of returning "not yet implemented". VTL (`statements` array) and JS resolver (`sql` field) payload formats both supported.

- **Athena: ListWorkGroups pagination** — Previously fetched up to 10,000 workgroups into memory and performed manual offset-based slicing, silently truncating results beyond the cap. Now delegates to the store's `Marker`/`MaxItems`/`NextMarker` pagination with the correct default page size of 50.

- **Lambda: SQS event source batching, multi-region container naming, and Python runtime** — SQS polling now fetches in a loop to honor `BatchSize` up to 10,000 (previously capped at a single `ReceiveMessage` call of max 10). Container names include region to prevent collisions in multi-region deployments. `copyCodeToContainer` selects `index.py` for Python runtimes.

- **EventBridge: wildcard suffix matching and InputTransformer value serialization** — Fixed `matchWildcardPattern` to use `strings.HasSuffix` when no `*` is present in the pattern. `applyInputTransformer` now JSON-serializes non-string placeholder values instead of using `fmt.Sprintf("%v")`.

- **Lambda ESM poller: remove redundant accountID field** — `esmPoller.accountID` was a stale copy of `lambdaSvc.accountID`. All references now use the canonical source.

- **Store: replace hardcoded pagination limits with generic utility functions** — Added `ListMatching`, `ListMatchingProto`, `ForEachAll`, `ForEachAllProto` to `common/base_store.go`. 17 call sites migrated from ad-hoc `MaxItems: 1000/10000` to the new APIs, eliminating arbitrary result caps on internal lookups and cleanup operations. Secrets Manager `ListSecrets` now accepts a store-level filter callback so that `IncludePlannedDeletion` and `Filter` parameters are evaluated inside the Pebble iterator instead of in a post-fetch pass; the unsorted path delegates to store-level `Marker`/`MaxItems` pagination for bounded memory usage. Cognito IDP `ListUsers`, `ListGroups`, `ListUserPoolClients`, `ListUserPools`, `ListResourceServers`, and `ListIdentityProviders` migrated from in-memory full-scan + handler-level pagination to store-level `Marker`/`MaxItems` with the `common.List` generic API; `ListUsers` passes its `Filter` parameter as a store callback so filtered items are skipped before counting toward the page limit. `ListUserPools` MaxResults upper bound corrected from 50 to 60 (AWS spec compliance). CloudWatch `DescribeAlarms`, `DescribeAlarmsForMetric`, `DescribeAlarmHistory`, and `ListDashboards` migrated from full-scan + handler-level filtering to store-level `common.List[T]` with `Marker`/`MaxItems` pagination and filter callbacks evaluated inside the Pebble iterator. CloudWatch `ListMetrics` migrated from unbounded in-memory scan to marker-based pagination (custom implementation for filesystem-backed chunk storage). `ListMetrics` default page size set to 500.

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

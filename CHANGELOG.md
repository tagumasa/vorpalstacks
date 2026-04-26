# Changelog

All notable changes to Vorpalstacks will be documented in this file.

## [Unreleased]

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

# Changelog

All notable changes to Vorpalstacks will be documented in this file.

## [Unreleased]

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

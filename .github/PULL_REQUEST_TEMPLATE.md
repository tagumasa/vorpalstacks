## Summary

Brief description of the changes.

## Affected AWS Service(s)

- [ ] None (infrastructure/tooling)
- [ ] S3
- [ ] SQS
- [ ] SNS
- [ ] Lambda
- [ ] DynamoDB
- [ ] IAM
- [ ] KMS
- [ ] Kinesis
- [ ] API Gateway
- [ ] Step Functions
- [ ] EventBridge
- [ ] CloudFront
- [ ] WAF / WAFv2
- [ ] Cognito
- [ ] Other: _______________

## Type of Change

- [ ] Bug fix (non-breaking)
- [ ] New AWS operation(s) implemented
- [ ] New AWS service added
- [ ] Refactoring / internal cleanup
- [ ] Documentation update
- [ ] Breaking change (describe below)

## Testing

- [ ] `make lint` passes (build + vet + fmt)
- [ ] `make test` passes (unit tests)
- [ ] SDK tests verified: `./sdk-tests-test -service <service> -v`
- [ ] 594/594 SDK tests still pass (no regressions)
- [ ] Manual testing with AWS CLI / SDK

## Checklist

- [ ] Godoc comments added to all new exported symbols
- [ ] No changes to `internal/core/resilience/errors.go` (`IsRetryable` must remain unchanged)
- [ ] No hardcoded account IDs in server code
- [ ] Proto files not directly edited (use `proto_generator` if needed)
- [ ] `pkg/sqlparser/dependency/` not modified
- [ ] New store operations use PebbleDB correctly (region isolation)
- [ ] AWS API response format matches AWS SDK expectations

## Breaking Changes

If this PR introduces breaking changes, describe what users need to change:

_(none)_

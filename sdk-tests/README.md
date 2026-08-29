# SDK-Based Tests for VorpalStacks

## Overview

This directory contains comprehensive SDK-based tests for verifying AWS service functionality on the VorpalStacks emulator. The tests use the official AWS SDK for Go v2 to ensure compatibility with real AWS APIs.

## Features

- **Independent Go Module**: Uses its own `go.mod` file, not inherited from parent project
- **AWS SDK v2**: Official AWS Go SDK v2 for production-grade testing
- **Comprehensive Coverage**: Tests covering every supported AWS service (list in [docs/services.md](../docs/services.md)) with 3,441 test cases (3,374 SDK + 50 cross-service integration + 17 WebSocket)
- **Easy to Run**: Simple CLI for running tests per service or all at once

## Supported Services

Rows follow the runner's service registrations: the Timestream row is one
combined service running both Timestream Write and Timestream Query tests,
which are separate services in the SDK classification used by
[docs/services.md](../docs/services.md).

| Service | Tests | Pass Rate | Status |
|---------|--------|-----------|--------|
| ACM | 54 | 100% | ✅ Perfect |
| API Gateway | 135 | 100% | ✅ Perfect |
| AppSync | 181 | 100% | ✅ Perfect |
| Athena | 68 | 100% | ✅ Perfect |
| CloudFront | 92 | 100% | ✅ Perfect |
| CloudTrail | 113 | 100% | ✅ Perfect |
| CloudWatch | 25 | 100% | ✅ Perfect |
| CloudWatch Logs | 62 | 100% | ✅ Perfect |
| Cognito | 94 | 100% | ✅ Perfect |
| Cognito Identity | 46 | 100% | ✅ Perfect |
| DynamoDB | 214 | 100% | ✅ Perfect |
| EC2 | 37 | 100% | ✅ Perfect |
| EventBridge | 67 | 100% | ✅ Perfect |
| IAM | 223 | 100% | ✅ Perfect |
| IoT | 416 | 100% | ✅ Perfect |
| Kinesis | 54 | 100% | ✅ Perfect |
| KMS | 102 | 100% | ✅ Perfect |
| Lambda | 119 | 100% | ✅ Perfect |
| Neptune | 103 | 100% | ✅ Perfect |
| NeptuneData | 170 | 100% | ✅ Perfect |
| NeptuneGraph | 49 | 100% | ✅ Perfect |
| RDS Data | 17 | 100% | ✅ Perfect |
| Route53 | 48 | 100% | ✅ Perfect |
| S3 | 143 | 100% | ✅ Perfect |
| Scheduler | 70 | 100% | ✅ Perfect |
| SecretsManager | 63 | 100% | ✅ Perfect |
| SESv2 | 94 | 100% | ✅ Perfect |
| SNS | 83 | 100% | ✅ Perfect |
| SQS | 87 | 100% | ✅ Perfect |
| SSM | 49 | 100% | ✅ Perfect |
| STS | 57 | 100% | ✅ Perfect |
| StepFunctions | 116 | 100% | ✅ Perfect |
| Timestream (Write+Query) | 51 | 100% | ✅ Perfect |
| WAFv2 | 72 | 100% | ✅ Perfect |

**Overall: 3,441/3,441 tests passing (100%) — 3,374 SDK + 50 integration + 17 WebSocket** (confirmed 2026-08-30 on main; per-session deltas live in git history)

*CloudTrail audit tests require `CLOUDTRAIL_ENABLED=true` (or `ALL_SERVICES_ENABLED=true`).*

## Prerequisites

1. **Go 1.25+** installed
2. **VorpalStacks server** running on `http://localhost:50080`
3. **AWS credentials** set (can be dummy values for testing)

### TEST_MODE and Seeded Tokens

When `TEST_MODE=true`, the server seeds delegated tokens (e.g. `dummy-trade-in-token-verify`) for STS `GetDelegatedAccessToken` tests. These tokens are **re-used across test runs** — `RedeemDelegationToken` re-seeds them immediately after use in TEST_MODE so that repeated test runs pass. In production, tokens are single-use as expected.

## Installation

```bash
cd sdk-tests
go mod tidy
go build -o sdk-tests-all .
```

## Usage

### Start VorpalStacks Server

```bash
# From project root
pkill -9 vorpalstacks 2>/dev/null; sleep 1
ALL_SERVICES_ENABLED=true SIGNATURE_VERIFICATION_ENABLED=false PORT=50080 DATA_PATH=./data TEST_MODE=true tmp/vorpalstacks > tmp/server.log 2>&1 &
```

### Start with CloudTrail Audit Enabled

To run CloudTrail audit integration tests, set `CLOUDTRAIL_ENABLED=true` (or use `ALL_SERVICES_ENABLED=true` to enable everything):

```bash
pkill -9 vorpalstacks 2>/dev/null; sleep 1
ALL_SERVICES_ENABLED=true SIGNATURE_VERIFICATION_ENABLED=false PORT=50080 DATA_PATH=./data TEST_MODE=true tmp/vorpalstacks > tmp/server.log 2>&1 &
```

Then run tests with the env var:

```bash
ALL_SERVICES_ENABLED=true ./sdk-tests/sdk-tests-all -service all -v
```

### Run Tests

```bash
# Test specific service
./sdk-tests-all -service cognito -v

# Test multiple services
./sdk-tests-all -service cognito,kinesis,acm -v

# Test all services
ALL_SERVICES_ENABLED=true ./sdk-tests-all -service all -v
```

### Available Services

```
dynamodb
sqs
sns
s3
lambda
iam
kms
events
stepfunctions
apigateway
logs
cloudwatch
ssm
secretsmanager
sts
scheduler
cognito
cognito-identity
sesv2
kinesis
acm
athena
timestream
route53
cloudfront
cloudtrail
wafv2
iot
neptune
neptunedata
neptunegraph
appsync
rdsdata
ec2
```

### Options

```
-endpoint string
    VorpalStacks endpoint (default "http://localhost:50080")
-region string
    AWS region (default "us-east-1")
-format string
    Output format: table, json (default "table")
-v
    Verbose output
```

## Examples

### Example 1: Test Lambda Service

```bash
cd sdk-tests
./sdk-tests-all -service lambda -v
```

Output:
```
Running: CreateFunction...
✓ CreateFunction (0.03s)
Running: ListFunctions...
✓ ListFunctions (0.01s)
...

=== LAMBDA ===
✓ CreateFunction (0.03s)
✓ ListFunctions (0.01s)
...

Summary: 22 passed, 1 failed
```

### Example 2: Test with Custom Endpoint

```bash
./sdk-tests-all -endpoint http://localhost:9000 -service dynamodb -v
```

### Example 3: Test Multiple Services

```bash
./sdk-tests-all -service sqs,sns,kms -v
```

## Test Structure

Each service test file follows this pattern:

```go
package testutil

func (r *TestRunner) RunServiceTests() []TestResult {
    // Load AWS SDK configuration
    cfg, err := config.LoadDefaultAWSConfig(...)
    
    // Create client
    client := service.NewFromConfig(cfg)
    
    // Run tests
    results = append(results, r.RunTest("service", "TestName", func() error {
        // Test implementation
        _, err := client.SomeOperation(ctx, &input)
        return err
    }))
    
    return results
}
```

## Cross-Service Integration Tests

In addition to per-service SDK tests, cross-service integration tests verify end-to-end delivery between services and direct protocol access to Neptune graph engines. These tests create real resources via SDK, trigger cross-service connections, and verify data arrives at the destination.

### Running Integration Tests

```bash
# From project root — server must be running
./sdk-tests/sdk-tests-all -service integration -endpoint http://127.0.0.1:50080
```

### Test Matrix

| Source | → Lambda | → SQS | → SNS | → Kinesis | → Step Functions |
|--------|----------|-------|-------|-----------|------------------|
| EventBridge | ✓ | ✓ | ✓ | ✓ | ✓ |
| CloudWatch Alarm | ✓ | | | | ✓ |
| CloudWatch Logs | ✓ | | | ✓ | |
| Scheduler | ✓ | ✓ | ✓ | | ✓ |
| Step Functions Task | ✓ | ✓ | ✓ | | |
| S3 Notifications | ✓ | ✓ | ✓ | | |
| SNS | ✓ | ✓ | | | |
| ESM (SQS) | ✓ | | | | |
| ESM (Kinesis) | ✓ | | | | |

### Verification Methods

| Method | Used By | What It Checks |
|--------|---------|----------------|
| CW Logs `/aws/lambda/<fn>` | EB/S3/CWAlarm/CWLogs/SNS→Lambda | Lambda was invoked by the source service (Docker execution writes logs) |
| `ReceiveMessage` | EB/Scheduler/SFN/S3/SNS→SQS, EB/Scheduler→SNS→SQS | Message arrived in destination queue |
| `GetRecords` | EB/CWLogs→Kinesis | Records written to Kinesis stream |
| `ListExecutions` | EB/CWAlarm/Scheduler→SFN | Step Functions execution was triggered |
| `DescribeAlarms` State | CWAlarm→Lambda/SNS/SFN | Alarm transitioned to ALARM state |
| ESM message consumption | ESM→SQS | Messages deleted from SQS after Lambda invocation |

### Verification Methods

3 tests verify that CloudTrail captures audit events from cross-service operations. These require `CLOUDTRAIL_ENABLED=true` (or `ALL_SERVICES_ENABLED=true`) at server startup; without it they are automatically skipped.

| Test | What It Verifies |
|------|-----------------|
| `CloudTrailAudit_CreateTrail_VerifyEvent` | CreateTrail generates a CloudTrail event findable by `EventName=CreateTrail` |
| `CloudTrailAudit_S3_PutObject` | S3 PutObject generates a CloudTrail event findable by `EventSource=s3.amazonaws.com` |
| `CloudTrailAudit_CrossService_EventSource` | Events from both `cloudtrail.amazonaws.com` and `s3.amazonaws.com` coexist in LookupEvents |

### Running

```bash
# Start server with audit enabled
ALL_SERVICES_ENABLED=true SIGNATURE_VERIFICATION_ENABLED=false PORT=50080 DATA_PATH=./data TEST_MODE=true tmp/vorpalstacks > tmp/server.log 2>&1 &

# Run audit tests
ALL_SERVICES_ENABLED=true ./sdk-tests/sdk-tests-all -service cloudtrail-audit -v
```

## Adding New Tests

### Dynamic Account ID and Region

**Never hardcode account IDs (`000000000000`, `123456789012`) or region (`us-east-1`) in test files.** The `TestRunner` obtains these dynamically via STS `GetCallerIdentity` at startup:

```go
r.AccountID()  // e.g. "test" in TEST_MODE, real account ID in production
r.Region       // "us-east-1" by default, configurable via -region flag
```

Test contexts store these in `region` and `accountID` fields:

```go
type myServiceTestCtx struct {
    client    *myservice.Client
    ctx       context.Context
    runner    *TestRunner
    region    string
    accountID string
}

func newMyServiceTestContext(r *TestRunner) *myServiceTestCtx {
    return &myServiceTestCtx{
        region:    r.Region,
        accountID: r.AccountID(),
        // ...
    }
}
```

When constructing ARNs, always use the dynamic values:

```go
// Correct
arn := fmt.Sprintf("arn:aws:iam::%s:role/%s", tc.accountID, roleName)

// Wrong — hardcoded account ID
arn := fmt.Sprintf("arn:aws:iam::000000000000:role/%s", roleName)
```

For IAM role ARNs used across integration tests, use the `intRoleARN` helper:

```go
func intRoleARN(roleName, accountID string) string {
    return fmt.Sprintf("arn:aws:iam::%s:role/%s", accountID, roleName)
}
```

**Note**: When using `fmt.Sprintf` inside backtick raw strings (e.g. JSON trust policies), wrap the entire raw string in `fmt.Sprintf` — placing `fmt.Sprintf(...)` inside a raw string literal does not evaluate it.

1. Create a new test file in `testutil/`:
   ```bash
   sdk-tests/testutil/myservice.go
   ```

2. Implement the test runner:
   ```go
   func (r *TestRunner) RunMyServiceTests() []TestResult {
       // Implementation
   }
   ```

3. Remove stub from `stubs.go`:
   ```go
   // Remove: func (r *TestRunner) RunMyServiceTests() []TestResult { return []TestResult{} }
   ```

4. Add service mapping in `main.go` if needed

5. Install SDK package:
   ```bash
   cd sdk-tests
   go get github.com/aws/aws-sdk-go-v2/service/myservice
   ```

6. Build and test:
   ```bash
   go build -o sdk-tests-all .
   ./sdk-tests-all -service myservice -v
   ```

## Output Formats

### Table Format (Default)
```
=== SERVICE ===
✓ TestName (0.00s)
✗ TestName (0.00s) - error message

Summary: 10 passed, 1 failed
```

### JSON Format
```json
{
  "results": [
    {
      "service": "lambda",
      "testName": "CreateFunction",
      "status": "PASS",
      "duration": "0.03s",
      "error": ""
    }
  ]
}
```

## Troubleshooting

### Server Not Running
```
Error: connection refused
Solution: Start VorpalStacks server on localhost:50080
```

### Build Errors
```
Error: package not found
Solution: Run `go mod tidy` to update dependencies
```

### Test Failures
Tests may fail due to:
1. Server-side implementation gaps
2. API validation errors (missing required fields)
3. Timeout issues (for async operations)

Review the error message to identify the cause.

## Reference Materials

- **VorpalStacks Docs**: `../README.md`

## Contributing

When adding new tests:
1. Follow existing patterns
2. Include CRUD operations where applicable
3. Test both success and failure cases
4. Clean up created resources
5. Handle async operations properly

## License

Same as parent VorpalStacks project.

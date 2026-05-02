# SDK-Based Tests for VorpalStacks

## Overview

This directory contains comprehensive SDK-based tests for verifying AWS service functionality on the VorpalStacks emulator. The tests use the official AWS SDK for Go v2 to ensure compatibility with real AWS APIs.

## Features

- **Independent Go Module**: Uses its own `go.mod` file, not inherited from parent project
- **AWS SDK v2**: Official AWS Go SDK v2 for production-grade testing
- **Comprehensive Coverage**: Tests for 32 AWS services with 2260 test cases (2216 SDK + 27 cross-service integration + 17 WebSocket)
- **Easy to Run**: Simple CLI for running tests per service or all at once

## Supported Services

| Service | Tests | Pass Rate | Status |
|---------|--------|-----------|--------|
| ACM | 43 | 100% | ✅ Perfect |
| API Gateway | 129 | 100% | ✅ Perfect |
| AppSync | 170 | 100% | ✅ Perfect |
| Athena | 65 | 100% | ✅ Perfect |
| CloudFront | 60 | 100% | ✅ Perfect |
| CloudTrail | 65 | 100% | ✅ Perfect |
| CloudWatch | 24 | 100% | ✅ Perfect |
| CloudWatch Logs | 43 | 100% | ✅ Perfect |
| Cognito | 67 | 100% | ✅ Perfect |
| Cognito Identity | 43 | 100% | ✅ Perfect |
| DynamoDB | 111 | 100% | ✅ Perfect |
| EventBridge | 59 | 100% | ✅ Perfect |
| IAM | 152 | 100% | ✅ Perfect |
| Kinesis | 51 | 100% | ✅ Perfect |
| KMS | 97 | 100% | ✅ Perfect |
| Lambda | 68 | 100% | ✅ Perfect |
| Neptune | 97 | 100% | ✅ Perfect |
| NeptuneData | 168 | 100% | ✅ Perfect |
| NeptuneGraph | 47 | 100% | ✅ Perfect |
| Route53 | 44 | 100% | ✅ Perfect |
| S3 | 90 | 100% | ✅ Perfect |
| Scheduler | 38 | 100% | ✅ Perfect |
| SecretsManager | 41 | 100% | ✅ Perfect |
| SESv2 | 79 | 100% | ✅ Perfect |
| SNS | 63 | 100% | ✅ Perfect |
| SQS | 49 | 100% | ✅ Perfect |
| SSM | 44 | 100% | ✅ Perfect |
| STS | 41 | 100% | ✅ Perfect |
| StepFunctions | 53 | 100% | ✅ Perfect |
| Timestream | 50 | 100% | ✅ Perfect |
| WAF | Removed | No longer a supported service |
| WAFv2 | 61 | 100% | ✅ Perfect |

**Overall: 2262/2262 tests passing (100%) — 2216 SDK + 29 integration + 17 WebSocket**

*CloudTrail audit tests require `VS_AUDIT_ENABLED=true`.*

## Prerequisites

1. **Go 1.25+** installed
2. **VorpalStacks server** running on `http://localhost:8080`
3. **AWS credentials** set (can be dummy values for testing)

## Installation

```bash
cd sdk-tests
go mod tidy
go build -o sdk-tests-test .
```

## Usage

### Start VorpalStacks Server

```bash
# From project root
pkill -9 vorpalstacks 2>/dev/null; sleep 1
SIGNATURE_VERIFICATION_ENABLED=false PORT=8080 DATA_PATH=./data TEST_MODE=true tmp/vorpalstacks > tmp/server.log 2>&1 &
```

### Start with CloudTrail Audit Enabled

To run CloudTrail audit integration tests, set `VS_AUDIT_ENABLED=true`:

```bash
pkill -9 vorpalstacks 2>/dev/null; sleep 1
rm -rf data/us-east-1 data/global
SIGNATURE_VERIFICATION_ENABLED=false VS_AUDIT_ENABLED=true PORT=8080 DATA_PATH=./data TEST_MODE=true tmp/vorpalstacks > tmp/server.log 2>&1 &
```

Then run tests with the env var:

```bash
VS_AUDIT_ENABLED=true ./sdk-tests/tmp/sdk-tests-all -service all -v
```

### Run Tests

```bash
# Test specific service
./sdk-tests-test -service cognito -v

# Test multiple services
./sdk-tests-test -service cognito,kinesis,acm -v

# Test all services
./sdk-tests-test -service all -v
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
sesv2
kinesis
acm
athena
timestream
route53
cloudfront
cloudtrail
waf
neptune
neptunedata
neptunegraph
```

### Options

```
-endpoint string
    VorpalStacks endpoint (default "http://localhost:8080")
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
./sdk-tests-test -service lambda -v
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
./sdk-tests-test -endpoint http://localhost:9000 -service dynamodb -v
```

### Example 3: Test Multiple Services

```bash
./sdk-tests-test -service sqs,sns,kms -v
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

In addition to per-service SDK tests, 27 cross-service integration tests verify end-to-end delivery between services. These tests create real resources via SDK, trigger cross-service connections, and verify data arrives at the destination.

### Running Integration Tests

```bash
# From project root — server must be running
./sdk-tests/tmp/sdk-tests-all -service integration -endpoint http://127.0.0.1:8080
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

3 tests verify that CloudTrail captures audit events from cross-service operations. These require `VS_AUDIT_ENABLED=true` at server startup; without it they are automatically skipped.

| Test | What It Verifies |
|------|-----------------|
| `CloudTrailAudit_CreateTrail_VerifyEvent` | CreateTrail generates a CloudTrail event findable by `EventName=CreateTrail` |
| `CloudTrailAudit_S3_PutObject` | S3 PutObject generates a CloudTrail event findable by `EventSource=s3.amazonaws.com` |
| `CloudTrailAudit_CrossService_EventSource` | Events from both `cloudtrail.amazonaws.com` and `s3.amazonaws.com` coexist in LookupEvents |

### Running

```bash
# Start server with audit enabled
SIGNATURE_VERIFICATION_ENABLED=false VS_AUDIT_ENABLED=true PORT=8080 DATA_PATH=./data TEST_MODE=true tmp/vorpalstacks > tmp/server.log 2>&1 &

# Run audit tests
VS_AUDIT_ENABLED=true ./sdk-tests/tmp/sdk-tests-all -service cloudtrail-audit -v
```

## Adding New Tests

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
   go build -o sdk-tests-test .
   ./sdk-tests-test -service myservice -v
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
Solution: Start VorpalStacks server on localhost:8080
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

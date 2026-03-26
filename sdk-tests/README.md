# SDK-Based Tests for Vorpalstacks

## Overview

Comprehensive SDK-based tests using the official AWS SDK for Go v2 to verify AWS service compatibility.

- **594 tests** across **25 services**
- **100% pass rate** (594/594)
- Independent Go module

## Supported Services

| Service | Tests | Status |
|---------|--------|--------|
| DynamoDB | 57 | Passing |
| SQS | 23 | Passing |
| SNS | 31 | Passing |
| S3 | 68 | Passing |
| Lambda | 50 | Passing |
| IAM | 155 | Passing |
| KMS | 45 | Passing |
| EventBridge | 39 | Passing |
| Step Functions | 28 | Passing |
| API Gateway | 80 | Passing |
| CloudWatch Logs | 29 | Passing |
| CloudWatch Metrics | 17 | Passing |
| SSM | 12 | Passing |
| SecretsManager | 21 | Passing |
| STS | 11 | Passing |
| Scheduler | 12 | Passing |
| Cognito IDP | 51 | Passing |
| SESv2 | 58 | Passing |
| Kinesis | 39 | Passing |
| ACM | 16 | Passing |
| Athena | 37 | Passing |
| Timestream | 32 | Passing |
| Route53 | 18 | Passing |
| CloudFront | 32 | Passing |
| CloudTrail | 24 | Passing |
| WAF | 34 | Passing |

## Prerequisites

- Go 1.25+
- Vorpalstacks server running on `http://localhost:8080`

## Building

```bash
cd sdk-tests
go build -o sdk-tests-test .
```

## Running

```bash
# Start the server
SIGNATURE_VERIFICATION_ENABLED=false PORT=8080 DATA_PATH=./data TEST_MODE=true ./vorpalstacks &

# Test all services
./sdk-tests-test -service all -v

# Test a specific service
./sdk-tests-test -service lambda -v

# Test multiple services
./sdk-tests-test -service cognito,kinesis,acm -v
```

## Options

```
-endpoint string
    Vorpalstacks endpoint (default "http://localhost:8080")
-region string
    AWS region (default "us-east-1")
-format string
    Output format: table, json (default "table")
-v
    Verbose output
```

## Adding New Tests

1. Create a new test file in `testutil/` following existing patterns
2. Implement the test runner in `testutil/<service>.go`
3. Remove stub from `stubs.go`
4. Add service mapping in `main.go`
5. Install the SDK package: `go get github.com/aws/aws-sdk-go-v2/service/<service>`
6. Build and test

## License

Same as parent Vorpalstacks project (FSL-1.1-MIT).

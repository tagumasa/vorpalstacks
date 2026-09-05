package lambda

// This file carries the per-invocation context record: one request id, one
// log stream and one deadline per execution. The runtime wrappers derive the
// handler's second argument from it, the execution logs reuse the same
// request id so context.awsRequestId matches the START/END/REPORT lines, and
// the custom-runtime Runtime API emulation serves the same values as
// Lambda-Runtime-* headers.

import (
	"fmt"
	"time"

	lambdastore "vorpalstacks/internal/store/aws/lambda"

	"github.com/google/uuid"
)

// invocationRecord is the server-side source of truth for the AWS Lambda
// context object of a single execution. Every member the handler's context
// exposes is derived from this record.
type invocationRecord struct {
	RequestID         string
	LogGroupName      string
	LogStreamName     string
	InvokedARN        string
	FunctionName      string
	Version           string
	MemorySize        int32
	TimeoutSeconds    int32
	Deadline          time.Time
	ClientContextJSON string
}

// lambdaLogGroupName is the CloudWatch log group of a function's
// invocations — the documented AWS_LAMBDA_LOG_GROUP_NAME value and the
// group the platform writes to, defined once.
func lambdaLogGroupName(functionName string) string {
	return "/aws/lambda/" + functionName
}

// lambdaLogStreamName is the CloudWatch Logs stream convention the
// platform writes: YYYY/MM/DD/[version]id[:8]. The id is the request id
// for an invocation record and the environment id for a sandbox.
func lambdaLogStreamName(now time.Time, version, id string) string {
	return fmt.Sprintf("%d/%02d/%02d/[%s]%s", now.Year(), now.Month(), now.Day(), version, id[:8])
}

// newInvocationRecord builds the record for one execution. The log stream
// name follows the CloudWatch Logs convention the platform writes:
// YYYY/MM/DD/[version]requestID[:8].
func newInvocationRecord(functionName, version, invokedARN string, memorySize, timeoutSeconds int32, clientContextJSON string) invocationRecord {
	// A stored configuration can carry a zero timeout (for example an
	// image-package function); the effective deadline must still exist, so
	// the record normalises to the create-side default.
	if timeoutSeconds <= 0 {
		timeoutSeconds = lambdastore.DefaultFunctionTimeoutSeconds
	}
	now := time.Now().UTC()
	requestID := uuid.New().String()
	return invocationRecord{
		RequestID:         requestID,
		LogGroupName:      lambdaLogGroupName(functionName),
		LogStreamName:     lambdaLogStreamName(now, version, requestID),
		InvokedARN:        invokedARN,
		FunctionName:      functionName,
		Version:           version,
		MemorySize:        memorySize,
		TimeoutSeconds:    timeoutSeconds,
		Deadline:          now.Add(time.Duration(timeoutSeconds) * time.Second),
		ClientContextJSON: clientContextJSON,
	}
}

// DeadlineUnixMS returns the deadline as Unix epoch milliseconds — the unit
// the Node.js/Python getRemainingTimeInMillis closures and the Runtime API's
// Lambda-Runtime-Deadline-Ms header both consume.
func (r invocationRecord) DeadlineUnixMS() int64 {
	return r.Deadline.UnixNano() / int64(time.Millisecond)
}

// qualifiedInvokeARN returns the ARN the invoker used: the plain function
// ARN without a qualifier, ARN:alias or ARN:version with one. The $LATEST
// qualifier leaves the plain ARN, matching the AWS context contract that
// invokedFunctionArn "indicates if the invoker specified a version number
// or alias".
func qualifiedInvokeARN(functionArn, qualifier string) string {
	if qualifier == "" || qualifier == "$LATEST" {
		return functionArn
	}
	return functionArn + ":" + qualifier
}

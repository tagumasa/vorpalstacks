package lambda

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
	"time"

	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

func TestNewInvocationRecord(t *testing.T) {
	now := time.Now().UTC()
	rec := newInvocationRecord("fn-name", "$LATEST", "arn:aws:lambda:us-east-1:123456789012:function:fn-name", 256, 10, "")
	if len(rec.RequestID) != 36 || strings.Count(rec.RequestID, "-") != 4 {
		t.Fatalf("RequestID must be a UUID, got %q", rec.RequestID)
	}
	if rec.LogGroupName != "/aws/lambda/fn-name" {
		t.Fatalf("LogGroupName = %q, want /aws/lambda/fn-name", rec.LogGroupName)
	}
	wantPrefix := fmt.Sprintf("%04d/%02d/%02d/[$LATEST]", now.Year(), now.Month(), now.Day())
	if !strings.HasPrefix(rec.LogStreamName, wantPrefix) {
		t.Fatalf("LogStreamName = %q, want prefix %q", rec.LogStreamName, wantPrefix)
	}
	if !strings.HasSuffix(rec.LogStreamName, rec.RequestID[:8]) {
		t.Fatalf("LogStreamName = %q, want it to end with the request id head %q", rec.LogStreamName, rec.RequestID[:8])
	}
	remaining := rec.DeadlineUnixMS() - now.UnixMilli()
	if remaining <= 0 || remaining > 10_000 {
		t.Fatalf("deadline must be ~10s away, got %dms", remaining)
	}
}

func TestQualifiedInvokeARN(t *testing.T) {
	const arn = "arn:aws:lambda:us-east-1:123456789012:function:fn"
	cases := []struct {
		qualifier string
		want      string
	}{
		{"", arn},
		{"$LATEST", arn},
		{"prod", arn + ":prod"},
		{"3", arn + ":3"},
	}
	for _, tc := range cases {
		if got := qualifiedInvokeARN(arn, tc.qualifier); got != tc.want {
			t.Fatalf("qualifiedInvokeARN(%q) = %q, want %q", tc.qualifier, got, tc.want)
		}
	}
}

func TestBuildInvokeCommandContext(t *testing.T) {
	s := &LambdaService{}
	const ccJSON = `{"custom":{"k":"v"}}`
	rec := newInvocationRecord("ctxfn", "7", "arn:aws:lambda:eu-west-1:123456789012:function:ctxfn:7", 512, 30, ccJSON)

	nodeCmd, _ := s.buildInvokeCommand("nodejs22.x", "index", "handler", `{"a":1}`, rec)
	nodeScript := nodeCmd[len(nodeCmd)-1]
	nodeNeedles := []string{
		"functionName:'ctxfn'",
		"functionVersion:'7'",
		"invokedFunctionArn:'arn:aws:lambda:eu-west-1:123456789012:function:ctxfn:7'",
		"memoryLimitInMB:512",
		"awsRequestId:'" + rec.RequestID + "'",
		"logGroupName:'/aws/lambda/ctxfn'",
		fmt.Sprintf("getRemainingTimeInMillis:()=>%d-Date.now()", rec.DeadlineUnixMS()),
		base64.StdEncoding.EncodeToString([]byte(ccJSON)),
		base64.StdEncoding.EncodeToString([]byte(`{"a":1}`)),
	}
	for _, needle := range nodeNeedles {
		if !strings.Contains(nodeScript, needle) {
			t.Fatalf("node wrapper must contain %q", needle)
		}
	}

	pyCmd, _ := s.buildInvokeCommand("python3.13", "index", "handler", `{}`, rec)
	pyScript := pyCmd[len(pyCmd)-1]
	pyNeedles := []string{
		"function_name='ctxfn'",
		"function_version='7'",
		"invoked_function_arn='arn:aws:lambda:eu-west-1:123456789012:function:ctxfn:7'",
		"memory_limit_in_mb=512",
		"aws_request_id='" + rec.RequestID + "'",
		"log_group_name='/aws/lambda/ctxfn'",
		fmt.Sprintf("get_remaining_time_in_millis=lambda:%d-int(time.time()*1000)", rec.DeadlineUnixMS()),
		"cc=json.loads(base64.b64decode(",
	}
	for _, needle := range pyNeedles {
		if !strings.Contains(pyScript, needle) {
			t.Fatalf("python wrapper must contain %q", needle)
		}
	}

	// Without a client context the wrappers embed the language null literal.
	emptyRec := newInvocationRecord("fn2", "$LATEST", "arn:aws:lambda:us-east-1:1:function:fn2", 128, 3, "")
	nodeCmd2, _ := s.buildInvokeCommand("nodejs22.x", "index", "handler", `{}`, emptyRec)
	if !strings.Contains(nodeCmd2[len(nodeCmd2)-1], "cc=null") {
		t.Fatalf("node wrapper without client context must embed null")
	}
	pyCmd2, _ := s.buildInvokeCommand("python3.13", "index", "handler", `{}`, emptyRec)
	if !strings.Contains(pyCmd2[len(pyCmd2)-1], "cc=None") {
		t.Fatalf("python wrapper without client context must embed None")
	}
}

func TestTimeoutEnvelope(t *testing.T) {
	if got := timeoutEnvelope(1); got != `{"errorMessage":"Task timed out after 1.00 seconds"}` {
		t.Fatalf("timeoutEnvelope(1) = %q", got)
	}
	if got := timeoutEnvelope(3); got != `{"errorMessage":"Task timed out after 3.00 seconds"}` {
		t.Fatalf("timeoutEnvelope(3) = %q", got)
	}
}

func TestNewInvocationRecordNormalisesTimeout(t *testing.T) {
	rec := newInvocationRecord("fn", "$LATEST", "arn", 128, 0, "")
	if rec.TimeoutSeconds != lambdastore.DefaultFunctionTimeoutSeconds {
		t.Fatalf("a zero timeout must normalise to the default, got %d", rec.TimeoutSeconds)
	}
	if rec.DeadlineUnixMS() <= time.Now().UnixMilli() {
		t.Fatalf("the normalised deadline must lie in the future")
	}
}

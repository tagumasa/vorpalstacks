package lambda

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveFunctionRef(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		wantName string
		wantQual string
	}{
		{"bare name", "my-function", "my-function", ""},
		{"name with alias suffix", "my-function:prod", "my-function", "prod"},
		{"name with version suffix", "my-function:2", "my-function", "2"},
		{"full ARN without qualifier", "arn:aws:lambda:us-west-2:123456789012:function:my-function", "my-function", ""},
		{"full ARN with alias", "arn:aws:lambda:us-west-2:123456789012:function:my-function:prod", "my-function", "prod"},
		{"full ARN with version", "arn:aws:lambda:us-west-2:123456789012:function:my-function:3", "my-function", "3"},
		{"partial ARN without qualifier", "123456789012:function:my-function", "my-function", ""},
		{"partial ARN with alias", "123456789012:function:my-function:prod", "my-function", "prod"},
		{"empty reference", "", "", ""},
		{"non-function ARN passes through", "arn:aws:sqs:us-east-1:123456789012:queue", "arn:aws:sqs:us-east-1:123456789012:queue", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			name, qualifier := resolveFunctionRef(tc.input)
			assert.Equal(t, tc.wantName, name)
			assert.Equal(t, tc.wantQual, qualifier)
		})
	}

	// The FunctionName @length(1,140) bound applies to the raw input: a
	// reference longer than 140 characters must resolve to an empty name
	// so validation rejects it, whatever form it takes.
	t.Run("reference exceeding the raw length bound", func(t *testing.T) {
		long := "arn:aws:lambda:us-west-2:123456789012:function:" + strings.Repeat("a", 141)
		name, _ := resolveFunctionRef(long)
		assert.Equal(t, "", name)
	})
}

func TestMergeQualifier(t *testing.T) {
	assert.Equal(t, "prod", mergeQualifier("prod", "5"), "explicit parameter wins")
	assert.Equal(t, "5", mergeQualifier("", "5"), "embedded qualifier used when parameter absent")
	assert.Equal(t, "", mergeQualifier("", ""), "no qualifier when both absent")
}

func TestFinalJSONDocument(t *testing.T) {
	cases := []struct {
		name   string
		stdout string
		want   string
		ok     bool
	}{
		{"single document", `{"errorMessage":"boom"}`, `{"errorMessage":"boom"}`, true},
		{"log line before document", "console output\n{\"errorMessage\":\"boom\"}", `{"errorMessage":"boom"}`, true},
		{"several log lines", "one\ntwo\nthree\n{}", "{}", true},
		{"scalar document", `"still json"`, `"still json"`, true},
		{"log output only", "console output", "", false},
		{"empty output", "", "", false},
		{"log line and non-JSON return value", "console output\nplain text", "", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			doc, ok := finalJSONDocument([]byte(tc.stdout))
			assert.Equal(t, tc.ok, ok)
			if tc.ok {
				assert.JSONEq(t, tc.want, string(doc))
			}
		})
	}
}

func TestClassifyFunctionError(t *testing.T) {
	// The classification reads the execution outcome only — the payload
	// never decides it, because on AWS a handler returning an
	// error-shaped document is a plain success.
	cases := []struct {
		name     string
		exitCode int
		want     string
	}{
		{"zero exit is success", 0, ""},
		{"handled error exit is Handled", execExitHandledError, "Handled"},
		{"uncaught error exit is Unhandled", 1, "Unhandled"},
		{"timeout exit is Unhandled", execExitTimedOut, "Unhandled"},
		{"killed process exit is Unhandled", 137, "Unhandled"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, classifyFunctionError(tc.exitCode))
		})
	}
}

func TestExtractWindowStateFinalDocument(t *testing.T) {
	// A windowed handler that logs before returning its state must still be
	// read: the state document is the final JSON document of the payload.
	state, err := extractWindowState([]byte(`{"a":1}` + "\n" + `{"state":{"count":2}}`))
	if err != nil {
		t.Fatalf("extractWindowState with log lines: %v", err)
	}
	var decoded struct {
		Count int `json:"count"`
	}
	if err := json.Unmarshal(state, &decoded); err != nil {
		t.Fatalf("unmarshal state: %v", err)
	}
	assert.Equal(t, 2, decoded.Count)

	if _, err := extractWindowState([]byte("log line")); err == nil {
		t.Fatalf("expected error for payload without a state document")
	}
}

func TestSplitResultPayload(t *testing.T) {
	const marker = "__VORPALSTACKS_RESULT_test__"
	cases := []struct {
		name    string
		stdout  string
		payload string
		logs    string
	}{
		{
			name:    "console output precedes the framed object",
			stdout:  "log one\nlog two\n" + marker + `{"ok":true}` + marker,
			payload: `{"ok":true}`,
			logs:    "log one\nlog two\n",
		},
		{
			name:    "no console output",
			stdout:  marker + `{"ok":true}` + marker,
			payload: `{"ok":true}`,
			logs:    "",
		},
		{
			name:    "multi-line string return is preserved exactly",
			stdout:  "log\n" + marker + "line one\nline two" + marker + "\n",
			payload: "line one\nline two",
			logs:    "log\n",
		},
		{
			name:    "empty return value region",
			stdout:  "log\n" + marker + marker,
			payload: "",
			logs:    "log\n",
		},
		{
			name:    "missing closing marker means no return value",
			stdout:  "log\n" + marker + `{"partial":`,
			payload: "",
			logs:    "log\n" + marker + `{"partial":`,
		},
		{
			name:    "missing opening marker",
			stdout:  "only logs",
			payload: "",
			logs:    "only logs",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			payload, logs := splitResultPayload(tc.stdout, marker)
			if payload != tc.payload {
				t.Fatalf("payload = %q, want %q", payload, tc.payload)
			}
			if logs != tc.logs {
				t.Fatalf("logs = %q, want %q", logs, tc.logs)
			}
		})
	}
}

func TestBuildInvokeCommandFraming(t *testing.T) {
	s := &LambdaService{}
	rec := newInvocationRecord("fn", "$LATEST", "arn:aws:lambda:us-east-1:123456789012:function:fn", 128, 3, "")
	nodeCmd, nodeMarker := s.buildInvokeCommand("nodejs22.x", "index", "handler", "{}", rec)
	if nodeCmd[0] != "timeout" || nodeCmd[1] != "-k" || nodeCmd[2] != "2" || nodeCmd[3] != "3" {
		t.Fatalf("the wrapper command must be wrapped in the function timeout with a kill escalation, got %v", nodeCmd[:4])
	}
	if nodeMarker == "" {
		t.Fatalf("nodejs invocations must frame the return value with a marker")
	}
	if got := strings.Count(nodeCmd[len(nodeCmd)-1], nodeMarker); got != 6 {
		t.Fatalf("the node wrapper must frame the promise return, the callback response, and the error document (two markers each), found %d marker occurrences", got)
	}
	if nodeMarker == strings.TrimSuffix(strings.TrimPrefix(nodeCmd[len(nodeCmd)-1], nodeMarker), nodeMarker) {
		t.Fatalf("the marker must not dominate the wrapper script")
	}

	pyCmd, pyMarker := s.buildInvokeCommand("python3.13", "index", "handler", "{}", rec)
	if pyCmd[0] != "timeout" || pyCmd[1] != "-k" {
		t.Fatalf("python invocations must also be wrapped in the timeout with a kill escalation, got %v", pyCmd[:2])
	}
	if pyMarker == "" {
		t.Fatalf("python invocations must frame the return value with a marker")
	}
	if got := strings.Count(pyCmd[len(pyCmd)-1], pyMarker); got != 4 {
		t.Fatalf("the python wrapper must frame the serialised value and the error document, found %d marker occurrences", got)
	}
	if pyMarker == nodeMarker {
		t.Fatalf("markers must be generated per invocation")
	}

	bootCmd, bootMarker := s.buildInvokeCommand("provided.al2023", "index", "handler", "{}", rec)
	if bootMarker != "" {
		t.Fatalf("custom runtimes write the response as their whole stdout and carry no marker")
	}
	// Bootstrap runtimes must not be handed guest-side utility
	// requirements: their deadline is enforced on the host, so the command
	// is the image's own entry contract without a timeout prefix.
	if len(bootCmd) != 1 || bootCmd[0] != "/var/runtime/bootstrap" {
		t.Fatalf("custom runtimes run their bootstrap command without a timeout prefix, got %v", bootCmd)
	}
}

func TestBuildInvokeCommandEscapesHandlerTokens(t *testing.T) {
	s := &LambdaService{}
	rec := newInvocationRecord("fn", "$LATEST", "arn:aws:lambda:us-east-1:123456789012:function:fn", 128, 3, "")
	// The handler validation constrains the format but not the character
	// set, so a quote in a module or function name must not break out of
	// the wrapper's single-quoted literals.
	nodeCmd, _ := s.buildInvokeCommand("nodejs22.x", "in'dex", "h'andler", "{}", rec)
	nodeScript := nodeCmd[len(nodeCmd)-1]
	if !strings.Contains(nodeScript, `in\'dex`) || !strings.Contains(nodeScript, `h\'andler`) {
		t.Fatalf("the node wrapper must escape quotes in the handler tokens")
	}
	if strings.Contains(nodeScript, `'in'dex'`) || strings.Contains(nodeScript, `'h'andler'`) {
		t.Fatalf("unescaped handler tokens must not reach the node wrapper script")
	}

	pyCmd, _ := s.buildInvokeCommand("python3.13", "in'dex", "h'andler", "{}", rec)
	pyScript := pyCmd[len(pyCmd)-1]
	if !strings.Contains(pyScript, `in\'dex`) || !strings.Contains(pyScript, `h\'andler`) {
		t.Fatalf("the python wrapper must escape quotes in the handler tokens")
	}
	if strings.Contains(pyScript, `'in'dex'`) || strings.Contains(pyScript, `'h'andler'`) {
		t.Fatalf("unescaped handler tokens must not reach the python wrapper script")
	}
}

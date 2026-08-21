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
	envelope := `{"errorMessage":"handled-boom","errorType":"Error"}`
	cases := []struct {
		name     string
		exitCode int
		stdout   string
		want     string
	}{
		{"non-zero exit is Unhandled", 1, "", "Unhandled"},
		{"clean error envelope is Handled", 0, envelope, "Handled"},
		{"log lines before the envelope keep it Handled", 0, "INFO start\nINFO end\n" + envelope, "Handled"},
		{"log output without an envelope is success", 0, "INFO only", ""},
		{"non-JSON payload is success", 0, "plain text", ""},
		{"payload without errorMessage is success", 0, `{"result":"ok"}`, ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, classifyFunctionError(tc.exitCode, tc.stdout))
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
	nodeCmd, nodeMarker := s.buildInvokeCommand("nodejs22.x", "index", "handler", "{}")
	if nodeMarker == "" {
		t.Fatalf("nodejs invocations must frame the return value with a marker")
	}
	if got := strings.Count(nodeCmd[2], nodeMarker); got != 4 {
		t.Fatalf("the node wrapper must frame the serialised value twice (two markers per write), found %d marker occurrences", got)
	}
	if nodeMarker == strings.TrimSuffix(strings.TrimPrefix(nodeCmd[2], nodeMarker), nodeMarker) {
		t.Fatalf("the marker must not dominate the wrapper script")
	}

	pyCmd, pyMarker := s.buildInvokeCommand("python3.13", "index", "handler", "{}")
	if pyMarker == "" {
		t.Fatalf("python invocations must frame the return value with a marker")
	}
	if got := strings.Count(pyCmd[2], pyMarker); got != 2 {
		t.Fatalf("the python wrapper must frame the serialised value, found %d marker occurrences", got)
	}
	if pyMarker == nodeMarker {
		t.Fatalf("markers must be generated per invocation")
	}

	bootCmd, bootMarker := s.buildInvokeCommand("provided.al2023", "index", "handler", "{}")
	if bootMarker != "" {
		t.Fatalf("custom runtimes write the response as their whole stdout and carry no marker")
	}
	if len(bootCmd) == 0 {
		t.Fatalf("custom runtimes keep their bootstrap command")
	}
}

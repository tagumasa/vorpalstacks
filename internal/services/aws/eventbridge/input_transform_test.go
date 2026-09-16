package eventbridge

import (
	"encoding/json"
	"testing"
	"time"

	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// docExampleEvent is the EC2 instance state-change event the input
// transformation documentation builds its example table on.
func docExampleEvent() *eventsstore.Event {
	return &eventsstore.Event{
		ID:         "c12ad73a-34ac-4d33-a8a0-2bb51b6bc0a6",
		Version:    "0",
		DetailType: "EC2 Instance State-change Notification",
		Source:     "aws.ec2",
		Account:    "123456789012",
		Time:       time.Date(2026, 9, 15, 10, 20, 30, 0, time.UTC),
		Region:     "us-east-1",
		Resources:  []string{"arn:aws:ec2:us-east-1:123456789012:instance/i-abcd1111"},
		Detail: map[string]interface{}{
			"instance": "i-0123456789",
			"state":    "RUNNING",
		},
		EventBusName: "default",
	}
}

const (
	docRuleARN   = "arn:aws:events:us-east-1:123456789012:rule/example"
	docRuleName  = "example"
	docInstance  = "$.detail.instance"
	docStatePath = "$.detail.state"
)

func docPaths(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func transformPayload(t *testing.T, event *eventsstore.Event, transformer *eventsstore.InputTransformer) []byte {
	t.Helper()
	svc := &EventsService{}
	return svc.applyInputTransform(docRuleARN, docRuleName, event, transformer)
}

func decodeJSON(t *testing.T, payload []byte) interface{} {
	t.Helper()
	var value interface{}
	if err := json.Unmarshal(payload, &value); err != nil {
		t.Fatalf("payload is not valid JSON: %v\n%s", err, payload)
	}
	return value
}

func asString(t *testing.T, value interface{}) string {
	t.Helper()
	str, ok := value.(string)
	if !ok {
		t.Fatalf("expected a JSON string, got %T: %v", value, value)
	}
	return str
}

func asMap(t *testing.T, value interface{}) map[string]interface{} {
	t.Helper()
	m, ok := value.(map[string]interface{})
	if !ok {
		t.Fatalf("expected a JSON object, got %T: %v", value, value)
	}
	return m
}

// The transformation example table in the input transformation
// documentation: string templates deliver raw string values as a JSON
// string, JSON templates auto-quote string values in value positions and
// insert them raw inside string literals.
func TestInputTransformDocExampleTable(t *testing.T) {
	cases := []struct {
		name       string
		paths      map[string]string
		template   string
		assert     func(*testing.T, interface{})
		exactBytes string // set when the delivered JSON text is asserted verbatim
	}{
		{
			name:     "simple string",
			paths:    docPaths("instance", docInstance, "state", docStatePath),
			template: `instance <instance> is in <state>`,
			assert: func(t *testing.T, v interface{}) {
				if got := asString(t, v); got != "instance i-0123456789 is in RUNNING" {
					t.Fatalf("string template output %q", got)
				}
			},
		},
		{
			name:     "string with escaped quotes",
			paths:    docPaths("instance", docInstance, "state", docStatePath),
			template: `instance \"<instance>\" is in <state>`,
			assert: func(t *testing.T, v interface{}) {
				if got := asString(t, v); got != `instance "i-0123456789" is in RUNNING` {
					t.Fatalf("escaped-quote template output %q", got)
				}
			},
		},
		{
			name:     "simple JSON",
			paths:    docPaths("instance", docInstance, "state", docStatePath),
			template: `{"instance" : <instance>, "state": <state>}`,
			assert: func(t *testing.T, v interface{}) {
				m := asMap(t, v)
				if m["instance"] != "i-0123456789" || m["state"] != "RUNNING" {
					t.Fatalf("JSON template output %v", m)
				}
			},
			exactBytes: `{"instance" : "i-0123456789", "state": "RUNNING"}`,
		},
		{
			name:     "JSON with strings and variables",
			paths:    docPaths("instance", docInstance, "state", docStatePath),
			template: `{"instance" : <instance>, "state" : "<state>", "instanceStatus" : "instance \"<instance>\" is in <state>"}`,
			assert: func(t *testing.T, v interface{}) {
				m := asMap(t, v)
				if m["state"] != "RUNNING" {
					t.Fatalf("quoted variable must not double-quote: %v", m)
				}
				if got := asString(t, m["instanceStatus"]); got != `instance "i-0123456789" is in RUNNING` {
					t.Fatalf("in-string escaped quotes: %q", got)
				}
			},
		},
		{
			name:     "JSON with mixed types",
			paths:    docPaths("instance", docInstance, "state", docStatePath),
			template: `{"instance" : <instance>, "state" : [ 9, <state>, true ], "Transformed" : "Yes"}`,
			assert: func(t *testing.T, v interface{}) {
				m := asMap(t, v)
				arr, ok := m["state"].([]interface{})
				if !ok || len(arr) != 3 || arr[0] != float64(9) || arr[1] != "RUNNING" || arr[2] != true {
					t.Fatalf("mixed-type array: %v", m["state"])
				}
			},
		},
		{
			name:     "object variable in JSON value",
			paths:    docPaths("detail", "$.detail"),
			template: `{"detail" : <detail>, "kept" : true}`,
			assert: func(t *testing.T, v interface{}) {
				m := asMap(t, v)
				detail := asMap(t, m["detail"])
				if detail["instance"] != "i-0123456789" {
					t.Fatalf("object variable: %v", m)
				}
			},
		},
		{
			name:     "object variable inside a string loses its quotes",
			paths:    docPaths("detail", "$.detail"),
			template: `{"note" : "Detail is <detail>"}`,
			assert: func(t *testing.T, v interface{}) {
				m := asMap(t, v)
				if got := asString(t, m["note"]); got != `Detail is {instance:i-0123456789,state:RUNNING}` {
					t.Fatalf("in-string object must drop internal quotes: %q", got)
				}
			},
		},
	}

	event := docExampleEvent()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			payload := transformPayload(t, event, &eventsstore.InputTransformer{
				InputPathsMap: tc.paths,
				InputTemplate: tc.template,
			})
			if tc.exactBytes != "" && string(payload) != tc.exactBytes {
				t.Fatalf("delivered bytes\n got: %s\nwant: %s", payload, tc.exactBytes)
			}
			tc.assert(t, decodeJSON(t, payload))
		})
	}
}

// The reserved variables the platform supplies: the rule context, the
// ingest timestamp, and the event with and without its detail member.
func TestInputTransformReservedVariables(t *testing.T) {
	event := docExampleEvent()

	// Value position: strings are quoted, the event object is raw JSON.
	payload := transformPayload(t, event, &eventsstore.InputTransformer{
		InputTemplate: `{"rule" : <aws.events.rule-arn>, "name" : <aws.events.rule-name>, "ingested" : <aws.events.event.ingestion-time>, "event" : <aws.events.event>, "full" : <aws.events.event.json>}`,
	})
	m := asMap(t, decodeJSON(t, payload))
	if m["rule"] != docRuleARN || m["name"] != docRuleName {
		t.Fatalf("rule context variables: %v", m)
	}
	if m["ingested"] != "2026-09-15T10:20:30Z" {
		t.Fatalf("ingestion-time variable: %v", m["ingested"])
	}
	inner := asMap(t, m["event"])
	if _, has := inner["detail"]; has {
		t.Fatalf("aws.events.event must omit the detail member: %v", inner)
	}
	if inner["id"] != event.ID || inner["source"] != "aws.ec2" {
		t.Fatalf("aws.events.event members: %v", inner)
	}
	full := asMap(t, m["full"])
	if asMap(t, full["detail"])["state"] != "RUNNING" {
		t.Fatalf("aws.events.event.json must carry detail: %v", full)
	}

	// The receipt stamp takes precedence over the event Time: the
	// variable is "generated by EventBridge and can't be overwritten",
	// so a stamped event renders the stamp even when the publisher
	// supplied a different Time. The unstamped fixture above pins the
	// fallback: a record predating the stamp renders its Time.
	stamped := docExampleEvent()
	stamped.IngestionTime = time.Date(2026, 9, 16, 8, 0, 0, 0, time.UTC)
	stampedPayload := transformPayload(t, stamped, &eventsstore.InputTransformer{
		InputTemplate: `{"ingested" : <aws.events.event.ingestion-time>}`,
	})
	if got := asMap(t, decodeJSON(t, stampedPayload))["ingested"]; got != "2026-09-16T08:00:00Z" {
		t.Fatalf("stamped ingestion-time variable: %v", got)
	}

	// String context: the reserved strings substitute raw.
	payload = transformPayload(t, event, &eventsstore.InputTransformer{
		InputTemplate: `"<aws.events.rule-name> triggered"`,
	})
	if got := asString(t, decodeJSON(t, payload)); got != "example triggered" {
		t.Fatalf("reserved variable in a string template: %q", got)
	}
}

// Paths resolve through the documented subset: array indices, dashed
// keys, wildcards over arrays and objects. A variable whose path does not
// resolve is elided from the output.
func TestInputTransformPathResolution(t *testing.T) {
	event := docExampleEvent()
	event.Resources = []string{"arn:one", "arn:two"}
	event.Detail = map[string]interface{}{
		"instance-id": "i-0123456789",
		"nested-deep": map[string]interface{}{"state": "RUNNING"},
		"counts":      []interface{}{float64(3), float64(5)},
	}

	cases := []struct {
		path  string
		want  interface{}
		found bool
	}{
		{path: "$.resources[0]", want: "arn:one", found: true},
		{path: "$.resources[1]", want: "arn:two", found: true},
		{path: "$.resources[2]", found: false},
		{path: "$.detail.nested-deep.state", want: "RUNNING", found: true},
		{path: "$.detail.instance-id", want: "i-0123456789", found: true},
		{path: "$.detail.counts[*]", want: []interface{}{float64(3), float64(5)}, found: true},
		{path: "$.detail.*", want: []interface{}{[]interface{}{float64(3), float64(5)}, "i-0123456789", map[string]interface{}{"state": "RUNNING"}}, found: true},
		{path: "$.detail.missing", found: false},
		{path: "$.detail.counts[9]", found: false},
		{path: "$.detail[0]", found: false},
		{path: "$['detail']", found: false},
	}

	envelope := eventEnvelope(event)
	for _, tc := range cases {
		got, found := resolveEventPathString(envelope, tc.path)
		if found != tc.found {
			t.Fatalf("path %q: found=%v want %v", tc.path, found, tc.found)
		}
		if !found {
			continue
		}
		wantJSON, _ := json.Marshal(tc.want)
		gotJSON, _ := json.Marshal(got)
		if string(gotJSON) != string(wantJSON) {
			t.Fatalf("path %q: got %s want %s", tc.path, gotJSON, wantJSON)
		}
	}

	// Elision: the unresolved placeholder does not appear in the output.
	// The elided template is no longer valid JSON, so the string-template
	// fallback delivers it as a JSON string.
	payload := transformPayload(t, event, &eventsstore.InputTransformer{
		InputPathsMap: docPaths("instance", "$.detail.missing"),
		InputTemplate: `{"instance": <instance>}`,
	})
	if got := asString(t, decodeJSON(t, payload)); got != `{"instance": }` {
		t.Fatalf("unresolved variable must be elided, got %s", got)
	}
}

// A JSON-string constant Input is delivered verbatim and an InputPath
// extraction may yield a scalar; an unresolvable InputPath passes the
// whole event.
func TestBuildTargetPayloadInputArms(t *testing.T) {
	svc := &EventsService{}
	event := docExampleEvent()

	payload := svc.buildTargetPayload(docRuleARN, docRuleName, event, eventsstore.Target{
		Input: `"HelloWorld!"`,
	})
	if string(payload) != `"HelloWorld!"` {
		t.Fatalf("constant Input must pass through verbatim, got %s", payload)
	}

	payload = svc.buildTargetPayload(docRuleARN, docRuleName, event, eventsstore.Target{
		InputPath: "$.detail.state",
	})
	if string(payload) != `"RUNNING"` {
		t.Fatalf("scalar InputPath extraction, got %s", payload)
	}

	payload = svc.buildTargetPayload(docRuleARN, docRuleName, event, eventsstore.Target{
		InputPath: "$.resources[0]",
	})
	if string(payload) != `"arn:aws:ec2:us-east-1:123456789012:instance/i-abcd1111"` {
		t.Fatalf("array-index InputPath extraction, got %s", payload)
	}

	payload = svc.buildTargetPayload(docRuleARN, docRuleName, event, eventsstore.Target{
		InputPath: "$.detail.absent",
	})
	if asMap(t, decodeJSON(t, payload))["id"] != event.ID {
		t.Fatalf("unresolvable InputPath must pass the whole event, got %s", payload)
	}

	payload = svc.buildTargetPayload(docRuleARN, docRuleName, event, eventsstore.Target{})
	if asMap(t, decodeJSON(t, payload))["detail-type"] != event.DetailType {
		t.Fatalf("no input configuration must pass the envelope, got %s", payload)
	}
}

// A template that stays text after substitution is delivered as a JSON
// string, including the documented multi-line shape where every line of
// the template is wrapped in quotes.
func TestInputTransformStringTemplateShape(t *testing.T) {
	event := docExampleEvent()
	payload := transformPayload(t, event, &eventsstore.InputTransformer{
		InputPathsMap: docPaths("state", docStatePath, "instance", docInstance),
		InputTemplate: "\"<instance> found\"\n\"state <state>\"",
	})
	lines := splitLines(asString(t, decodeJSON(t, payload)))
	if len(lines) != 2 || lines[0] != `"i-0123456789 found"` || lines[1] != `"state RUNNING"` {
		t.Fatalf("multi-line string template: %q", lines)
	}

	// An unknown placeholder keeps its literal text: template placeholders
	// carry no validation.
	payload = transformPayload(t, event, &eventsstore.InputTransformer{
		InputTemplate: `value <not-defined> stays`,
	})
	if got := asString(t, decodeJSON(t, payload)); got != "value <not-defined> stays" {
		t.Fatalf("unknown placeholder: %q", got)
	}
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	return append(lines, s[start:])
}

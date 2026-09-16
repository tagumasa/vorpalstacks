package eventbridge

import (
	"testing"
	"time"

	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// vocabEvent returns the AWS documentation's EC2 state-change event as the
// base for vocabulary cases, with optional mutations per case.
func vocabEvent(mutate func(*eventsstore.Event)) *eventsstore.Event {
	event := &eventsstore.Event{
		ID:           "7bf73129-1428-4cd3-a780-95db273d1602",
		Version:      "0",
		DetailType:   "EC2 Instance State-change Notification",
		Source:       "aws.ec2",
		Account:      "123456789012",
		Time:         time.Date(2015, 11, 11, 21, 29, 54, 0, time.UTC),
		Region:       "us-east-1",
		Resources:    []string{"arn:aws:ec2:us-east-1:123456789012:instance/i-abcd1111"},
		EventBusName: "default",
		Detail: map[string]interface{}{
			"instance-id":     "i-abcd1111",
			"state":           "running",
			"c-count":         float64(3),
			"d-count":         float64(7),
			"x-limit":         float64(301.8),
			"service":         "EventBridge",
			"FileName":        "dir/icon.png",
			"FilePath":        "/usr/local/bin/tool",
			"sourceIPAddress": "10.0.0.5",
			"LastName":        "",
		},
	}
	if mutate != nil {
		mutate(event)
	}
	return event
}

// TestPatternVocabularyMatrix is the matcher family's class-closure test:
// every documented operator × operand form × value type in one table,
// evaluated through the delivery plane's entry (matchEventPattern over the
// rendered envelope). Rows cite the documentation form they pin.
func TestPatternVocabularyMatrix(t *testing.T) {
	cases := []struct {
		name    string
		pattern string
		event   *eventsstore.Event
		want    bool
	}{
		// Equals on scalars, and the multi-value OR list.
		{"equals string", `{"source": ["aws.ec2"]}`, vocabEvent(nil), true},
		{"equals string miss", `{"source": ["aws.s3"]}`, vocabEvent(nil), false},
		{"equals list any", `{"source": ["aws.ec2", "aws.s3"]}`, vocabEvent(nil), true},
		{"equals empty string", `{"detail": {"LastName": [""]}}`, vocabEvent(nil), true},
		{"equals nested level", `{"detail": {"instance-id": ["i-abcd1111"]}}`, vocabEvent(nil), true},

		// Equals stays within the value's JSON type: pattern values follow
		// JSON rules (strings, numbers, true/false/null), so a string event
		// value never matches a numeric or boolean pattern value and vice
		// versa, whatever their text renders to.
		{"equals number", `{"detail": {"c-count": [3]}}`, vocabEvent(nil), true},
		{"equals string event vs number pattern", `{"detail": {"state": [1]}}`,
			vocabEvent(func(e *eventsstore.Event) { e.Detail["state"] = "1" }), false},
		{"equals string event vs bool pattern", `{"detail": {"state": [true]}}`,
			vocabEvent(func(e *eventsstore.Event) { e.Detail["state"] = "true" }), false},
		{"equals number event vs string pattern", `{"detail": {"c-count": ["3"]}}`, vocabEvent(nil), false},
		{"equals bool", `{"detail": {"flag": [true]}}`,
			vocabEvent(func(e *eventsstore.Event) { e.Detail["flag"] = true }), true},
		{"equals bool event vs number pattern", `{"detail": {"flag": [1]}}`,
			vocabEvent(func(e *eventsstore.Event) { e.Detail["flag"] = true }), false},
		{"equals null", `{"detail": {"absent": [null]}}`,
			vocabEvent(func(e *eventsstore.Event) { e.Detail["absent"] = nil }), true},
		{"equals null pattern vs string event", `{"detail": {"state": [null]}}`, vocabEvent(nil), false},

		// Prefix: plain and the equals-ignore-case conjunction.
		{"prefix on time", `{"time": [{"prefix": "2015-11-11"}]}`, vocabEvent(nil), true},
		{"prefix miss", `{"time": [{"prefix": "2016"}]}`, vocabEvent(nil), false},
		{"prefix ignore case", `{"detail": {"service": [{"prefix": {"equals-ignore-case": "eventb"}}]}}`, vocabEvent(nil), true},
		{"prefix ignore case miss", `{"detail": {"service": [{"prefix": {"equals-ignore-case": "bridge"}}]}}`, vocabEvent(nil), false},

		// Suffix: plain and the equals-ignore-case conjunction.
		{"suffix", `{"detail": {"FileName": [{"suffix": ".png"}]}}`, vocabEvent(nil), true},
		{"suffix case sensitive", `{"detail": {"FileName": [{"suffix": ".PNG"}]}}`, vocabEvent(nil), false},
		{"suffix ignore case", `{"detail": {"FileName": [{"suffix": {"equals-ignore-case": ".PNG"}}]}}`, vocabEvent(nil), true},

		// Equals-ignore-case standalone.
		{"equals-ignore-case", `{"detail-type": [{"equals-ignore-case": "ec2 instance state-change notification"}]}`, vocabEvent(nil), true},

		// Numeric: equality, range, scientific notation.
		{"numeric equals", `{"detail": {"x-limit": [{"numeric": ["=", 3.018e2]}]}}`, vocabEvent(nil), true},
		{"numeric range", `{"detail": {"c-count": [{"numeric": [">", 0, "<=", 5]}]}}`, vocabEvent(nil), true},
		{"numeric range miss", `{"detail": {"d-count": [{"numeric": ["<", 10, ">", 100]}]}}`, vocabEvent(nil), false},

		// Anything-but: scalars, lists, and every documented conjunction.
		{"anything-but string", `{"detail": {"state": [{"anything-but": "initializing"}]}}`, vocabEvent(nil), true},
		{"anything-but string excluded", `{"detail": {"state": [{"anything-but": "running"}]}}`, vocabEvent(nil), false},
		{"anything-but number", `{"detail": {"x-limit": [{"anything-but": 123}]}}`, vocabEvent(nil), true},
		{"anything-but list miss", `{"detail": {"state": [{"anything-but": ["stopped", "running"]}]}}`, vocabEvent(nil), false},
		{"anything-but number list", `{"detail": {"x-limit": [{"anything-but": [100, 200, 301.8]}]}}`, vocabEvent(nil), false},
		{"anything-but prefix", `{"detail": {"state": [{"anything-but": {"prefix": "init"}}]}}`, vocabEvent(nil), true},
		{
			"anything-but prefix list excludes",
			`{"detail": {"state": [{"anything-but": {"prefix": ["init", "run"]}}]}}`,
			vocabEvent(func(e *eventsstore.Event) { e.Detail["state"] = "running" }), false,
		},
		{
			"anything-but prefix list passes",
			`{"detail": {"state": [{"anything-but": {"prefix": ["start", "stop"]}}]}}`,
			vocabEvent(func(e *eventsstore.Event) { e.Detail["state"] = "initializing" }), true,
		},
		{
			"anything-but suffix",
			`{"detail": {"FileName": [{"anything-but": {"suffix": ".txt"}}]}}`,
			vocabEvent(func(e *eventsstore.Event) { e.Detail["FileName"] = "notes.txt" }), false,
		},
		{
			"anything-but suffix list",
			`{"detail": {"FileName": [{"anything-but": {"suffix": [".txt", ".rtf"]}}]}}`,
			vocabEvent(func(e *eventsstore.Event) { e.Detail["FileName"] = "document.rtf" }), false,
		},
		{
			"anything-but equals-ignore-case excludes casings",
			`{"detail": {"state": [{"anything-but": {"equals-ignore-case": "initializing"}}]}}`,
			vocabEvent(func(e *eventsstore.Event) { e.Detail["state"] = "INITIALIZING" }), false,
		},
		{
			"anything-but equals-ignore-case list",
			`{"detail": {"state": [{"anything-but": {"equals-ignore-case": ["initializing", "stopped"]}}]}}`,
			vocabEvent(func(e *eventsstore.Event) { e.Detail["state"] = "Stopped" }), false,
		},
		{
			"anything-but wildcard",
			`{"detail": {"FilePath": [{"anything-but": {"wildcard": "*/lib/*"}}]}}`,
			vocabEvent(func(e *eventsstore.Event) { e.Detail["FilePath"] = "/usr/lib/tool" }), false,
		},
		{
			"anything-but wildcard list",
			`{"detail": {"FilePath": [{"anything-but": {"wildcard": ["*/lib/*", "*/bin/*"]}}]}}`,
			vocabEvent(func(e *eventsstore.Event) { e.Detail["FilePath"] = "/usr/bin/tool" }), false,
		},
		{
			"anything-but wildcard list passes",
			`{"detail": {"FilePath": [{"anything-but": {"wildcard": ["*/lib/*", "*/bin/*"]}}]}}`,
			vocabEvent(func(e *eventsstore.Event) { e.Detail["FilePath"] = "/etc/tool" }), true,
		},

		// Exists: presence and absence, at the top level and inside detail.
		{"exists true", `{"detail": {"state": [{"exists": true}]}}`, vocabEvent(nil), true},
		{"exists true absent", `{"detail": {"missing": [{"exists": true}]}}`, vocabEvent(nil), false},
		{"exists false absent", `{"detail": {"missing": [{"exists": false}]}}`, vocabEvent(nil), true},
		{"exists false present", `{"detail": {"state": [{"exists": false}]}}`, vocabEvent(nil), false},
		{"exists false top level", `{"trace-header": [{"exists": false}]}`, vocabEvent(nil), true},

		// Wildcards: anchoring and the backslash escapes.
		{"wildcard", `{"detail": {"FileName": [{"wildcard": "dir/*.png"}]}}`, vocabEvent(nil), true},
		{"wildcard anchored", `{"detail": {"FileName": [{"wildcard": "dir/*"}]}}`,
			vocabEvent(func(e *eventsstore.Event) { e.Detail["FileName"] = "xdir/icon.png" }), false},
		{"wildcard escaped star is literal", `{"detail": {"FileName": [{"wildcard": "a\\*b"}]}}`,
			vocabEvent(func(e *eventsstore.Event) { e.Detail["FileName"] = "a*b" }), true},
		{"wildcard escaped star rejects wildcard use", `{"detail": {"FileName": [{"wildcard": "a\\*b"}]}}`,
			vocabEvent(func(e *eventsstore.Event) { e.Detail["FileName"] = "aXb" }), false},
		{"wildcard escaped backslash is literal", `{"detail": {"FilePath": [{"wildcard": "/usr/local\\\\bin/*"}]}}`,
			vocabEvent(func(e *eventsstore.Event) { e.Detail["FilePath"] = "/usr/local\\bin/tool" }), true},
		{"wildcard escaped backslash rejects the unescaped", `{"detail": {"FilePath": [{"wildcard": "/usr/local\\\\bin/*"}]}}`,
			vocabEvent(nil), false},

		// CIDR.
		{"cidr in range", `{"detail": {"sourceIPAddress": [{"cidr": "10.0.0.0/24"}]}}`, vocabEvent(nil), true},
		{"cidr out of range", `{"detail": {"sourceIPAddress": [{"cidr": "192.168.0.0/24"}]}}`, vocabEvent(nil), false},

		// Array-valued event members: the intersection rule.
		{"event array equals", `{"resources": ["arn:aws:ec2:us-east-1:123456789012:instance/i-abcd1111"]}`, vocabEvent(nil), true},
		{"event array equals miss", `{"resources": ["arn:aws:ec2:us-east-1:123456789012:instance/i-other"]}`, vocabEvent(nil), false},
		{
			"pattern and event arrays intersect",
			`{"resources": ["arn:a", "arn:aws:ec2:us-east-1:123456789012:instance/i-abcd1111"]}`,
			vocabEvent(func(e *eventsstore.Event) {
				e.Resources = []string{"arn:aws:ec2:us-east-1:123456789012:instance/i-abcd1111", "arn:b"}
			}), true,
		},
		{
			"operator over event array elements",
			`{"detail": {"items": [{"prefix": "x-"}]}}`,
			vocabEvent(func(e *eventsstore.Event) { e.Detail["items"] = []interface{}{"a-1", "x-2"} }), true,
		},
		{
			// Operators on array members generalise the intersection rule
			// per element (the documented core case is equality): the
			// member matches when any element satisfies the pattern.
			"anything-but over event array elements",
			`{"detail": {"tags": [{"anything-but": "prod"}]}}`,
			vocabEvent(func(e *eventsstore.Event) { e.Detail["tags"] = []interface{}{"west", "prod"} }), true,
		},
		{
			"anything-but over event array all elements excluded",
			`{"detail": {"tags": [{"anything-but": "prod"}]}}`,
			vocabEvent(func(e *eventsstore.Event) { e.Detail["tags"] = []interface{}{"prod"} }), false,
		},

		// $or: top level, inside detail, and combined with siblings (AND).
		{
			"$or top level first disjunct",
			`{"$or": [{"source": ["aws.ec2"]}, {"source": ["aws.s3"]}]}`,
			vocabEvent(nil), true,
		},
		{
			"$or top level second disjunct",
			`{"$or": [{"source": ["aws.s3"]}, {"source": ["aws.ec2"]}]}`,
			vocabEvent(nil), true,
		},
		{
			"$or no disjunct matches",
			`{"$or": [{"source": ["aws.s3"]}, {"source": ["aws.lambda"]}]}`,
			vocabEvent(nil), false,
		},
		{
			"$or combines with siblings by AND",
			`{"region": ["us-east-1"], "$or": [{"source": ["aws.s3"]}, {"detail-type": ["EC2 Instance State-change Notification"]}]}`,
			vocabEvent(nil), true,
		},
		{
			"$or sibling failure fails the pattern",
			`{"region": ["eu-west-1"], "$or": [{"source": ["aws.ec2"]}]}`,
			vocabEvent(nil), false,
		},
		{
			"$or inside detail (documentation example)",
			`{"detail": {"$or": [{"c-count": [{"numeric": [">", 0, "<=", 5]}]}, {"d-count": [{"numeric": ["<", 10]}]}]}}`,
			vocabEvent(nil), true,
		},
		{
			"$or inside detail no disjunct",
			`{"detail": {"$or": [{"c-count": [{"numeric": [">", 100]}]}, {"d-count": [{"numeric": ["<", 1]}]}]}}`,
			vocabEvent(nil), false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := (&EventsService{}).matchEventPattern(tc.event, tc.pattern); got != tc.want {
				t.Fatalf("matchEventPattern = %v, want %v", got, tc.want)
			}
		})
	}
}

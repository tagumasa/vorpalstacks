package cloudwatchlogs

import "testing"

// resolveCase is one row of an identifier resolver's table.
type resolveCase struct {
	name       string
	identifier string
	want       string
}

// runResolveTable runs one identifier resolver's table: every resolver in
// the family passes bare values through, resolves its own family's ARN
// form, and passes foreign ARNs through unchanged for the not-found report.
func runResolveTable(t *testing.T, resolve func(string) string, tests []resolveCase) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolve(tt.identifier); got != tt.want {
				t.Errorf("resolve(%q) = %q, want %q", tt.identifier, got, tt.want)
			}
		})
	}
}

func TestResolveLogGroupIdentifier(t *testing.T) {
	runResolveTable(t, resolveLogGroupIdentifier, []resolveCase{
		{"bare name passes through", "/aws/lambda/my-function", "/aws/lambda/my-function"},
		{"log group ARN resolves", "arn:aws:logs:us-east-1:123456789012:log-group:/aws/lambda/my-function", "/aws/lambda/my-function"},
		{"object ARN tolerates the :* suffix", "arn:aws:logs:us-east-1:123456789012:log-group:/aws/lambda/my-function:*", "/aws/lambda/my-function"},
		{"log stream ARN yields its group", "arn:aws:logs:us-east-1:123456789012:log-group:my-group:log-stream:my-stream", "my-group"},
		{"unresolvable ARN passes through for the not-found report", "arn:aws:logs:us-east-1:123456789012:destination:my-dest", "arn:aws:logs:us-east-1:123456789012:destination:my-dest"},
		{"empty stays empty", "", ""},
	})
}

func TestResolveScheduledQueryIdentifier(t *testing.T) {
	runResolveTable(t, resolveScheduledQueryIdentifier, []resolveCase{
		{"bare id passes through", "sq-1234567890", "sq-1234567890"},
		{"scheduled query ARN resolves", "arn:aws:logs:us-east-1:123456789012:scheduled-query:sq-1234567890", "sq-1234567890"},
		{"unresolvable ARN passes through", "arn:aws:logs:us-east-1:123456789012:log-group:my-group", "arn:aws:logs:us-east-1:123456789012:log-group:my-group"},
	})
}

func TestResolveLookupTableIdentifier(t *testing.T) {
	runResolveTable(t, resolveLookupTableIdentifier, []resolveCase{
		{"bare name passes through", "my_table", "my_table"},
		{"lookup table ARN resolves", "arn:aws:logs:us-east-1:123456789012:lookup-table:my_table", "my_table"},
		{"unresolvable ARN passes through", "arn:aws:logs:us-east-1:123456789012:destination:my-dest", "arn:aws:logs:us-east-1:123456789012:destination:my-dest"},
	})
}

func TestLogGroupNameOrIdentifier(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]interface{}
		want   string
	}{
		{"identifier ARN resolves", map[string]interface{}{"logGroupIdentifier": "arn:aws:logs:us-east-1:123456789012:log-group:g"}, "g"},
		{"identifier name passes through", map[string]interface{}{"logGroupIdentifier": "g"}, "g"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := logGroupNameOrIdentifier(tt.params)
			if err != nil {
				t.Fatalf("logGroupNameOrIdentifier() error: %v", err)
			}
			if got != tt.want {
				t.Errorf("logGroupNameOrIdentifier() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLogGroupTarget(t *testing.T) {
	if got := logGroupTarget("named", "arn:aws:logs:us-east-1:123456789012:log-group:other"); got != "named" {
		t.Errorf("explicit name must win, got %q", got)
	}
	if got := logGroupTarget("", "arn:aws:logs:us-east-1:123456789012:log-group:arn-resolved"); got != "arn-resolved" {
		t.Errorf("identifier must resolve when the name is absent, got %q", got)
	}
}

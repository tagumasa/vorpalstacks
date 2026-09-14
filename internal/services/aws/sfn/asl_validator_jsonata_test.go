package sfn

import (
	"strings"
	"testing"
)

// TestJSONataExpressionDelimiters pins the creation-time delimitation
// contract: "the string must start with `{%` with no leading spaces, and
// must end with `%}` with no trailing spaces. Improperly opening or
// closing the expression will result in a validation error."
func TestJSONataExpressionDelimiters(t *testing.T) {
	cases := []struct {
		name    string
		output  string
		wantErr bool
	}{
		{"well-formed expression", "{% $states.input.x %}", false},
		{"plain literal", "just text", false},
		{"unterminated expression", "{% $states.input.x", true},
		{"leading space before opener", " {% $states.input.x %}", true},
		{"trailing space after closer", "{% $states.input.x %} ", true},
		{"closer without opener", "$states.input.x %}", true},
		{"mid-string delimiters", "prefix {% $x %} suffix", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			def := `{"QueryLanguage": "JSONata", "StartAt": "A", "States": {"A": {"Type": "Pass", "Output": {"v": ` + jsonString(tc.output) + `}, "End": true}}}`
			err := validateDefinitionStructure(def, "STANDARD")
			if tc.wantErr && err == nil {
				t.Fatalf("definition with output %q was accepted", tc.output)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("definition with output %q rejected: %v", tc.output, err)
			}
		})
	}
}

// TestJSONataStatesNodeAccessibility pins the creation-time accessibility
// contract: "Attempting to access $states.result or $states.errorOutput
// in fields and states where they are not accessible will be caught at
// creation, update, or validation of the state machine."
func TestJSONataStatesNodeAccessibility(t *testing.T) {
	cases := []struct {
		name    string
		def     string
		wantErr bool
	}{
		{"result in Pass Output", `{"QueryLanguage": "JSONata", "StartAt": "A", "States": {"A": {"Type": "Pass", "Output": "{% $states.result %}", "End": true}}}`, true},
		{"result in Choice condition", `{"QueryLanguage": "JSONata", "StartAt": "A", "States": {"A": {"Type": "Choice", "Choices": [{"Condition": "{% $states.result > 1 %}", "Next": "B"}], "Default": "B"}, "B": {"Type": "Succeed"}}}`, true},
		{"result in Wait Seconds", `{"QueryLanguage": "JSONata", "StartAt": "A", "States": {"A": {"Type": "Wait", "Seconds": "{% $states.result %}", "End": true}}}`, true},
		{"result in Task Output", `{"QueryLanguage": "JSONata", "StartAt": "A", "States": {"A": {"Type": "Task", "Resource": "arn:aws:states:::lambda:invoke", "Output": "{% $states.result %}", "End": true}}}`, false},
		{"result in Map Output", `{"QueryLanguage": "JSONata", "StartAt": "A", "States": {"A": {"Type": "Map", "Output": "{% $states.result %}", "ItemProcessor": {"StartAt": "W", "States": {"W": {"Type": "Pass", "End": true}}}, "End": true}}}`, false},
		{"similar node name is not result", `{"QueryLanguage": "JSONata", "StartAt": "A", "States": {"A": {"Type": "Pass", "Output": "{% $states.input.resultX %}", "End": true}}}`, false},
		{"errorOutput in state Output", `{"QueryLanguage": "JSONata", "StartAt": "A", "States": {"A": {"Type": "Task", "Resource": "arn:aws:states:::lambda:invoke", "Output": "{% $states.errorOutput %}", "End": true}}}`, true},
		{"errorOutput in Catch Assign", `{"QueryLanguage": "JSONata", "StartAt": "A", "States": {"A": {"Type": "Task", "Resource": "arn:aws:states:::lambda:invoke", "Catch": [{"ErrorEquals": ["States.ALL"], "Assign": {"e": "{% $states.errorOutput %}"}, "Next": "B"}], "End": true}, "B": {"Type": "Succeed"}}}`, false},
		{"errorOutput in Catch Output", `{"QueryLanguage": "JSONata", "StartAt": "A", "States": {"A": {"Type": "Task", "Resource": "arn:aws:states:::lambda:invoke", "Catch": [{"ErrorEquals": ["States.ALL"], "Output": "{% $states.errorOutput %}", "Next": "B"}], "End": true}, "B": {"Type": "Succeed"}}}`, false},
		{"errorOutput in state Assign", `{"QueryLanguage": "JSONata", "StartAt": "A", "States": {"A": {"Type": "Map", "Assign": {"e": "{% $states.errorOutput %}"}, "ItemProcessor": {"StartAt": "W", "States": {"W": {"Type": "Pass", "End": true}}}, "End": true}}}`, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateDefinitionStructure(tc.def, "STANDARD")
			if tc.wantErr && err == nil {
				t.Fatal("definition was accepted")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("definition rejected: %v", err)
			}
		})
	}
}

// jsonString embeds a Go string as a JSON string literal.
func jsonString(s string) string {
	var sb []byte
	sb = append(sb, '"')
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"':
			sb = append(sb, '\\', '"')
		case '\\':
			sb = append(sb, '\\', '\\')
		default:
			sb = append(sb, s[i])
		}
	}
	sb = append(sb, '"')
	return string(sb)
}

// TestTruncateForDiagnosticCutsOnRuneBoundary pins the diagnostic
// truncation: the cut lands on a rune boundary so a multi-byte value
// never prints as mojibake.
func TestTruncateForDiagnosticCutsOnRuneBoundary(t *testing.T) {
	cjk := strings.Repeat("あ", 60)
	got := truncateForDiagnostic(cjk)
	want := strings.Repeat("あ", 40) + "…"
	if got != want {
		t.Fatalf("truncate = %q, want %q", got, want)
	}
	if short := truncateForDiagnostic("short"); short != "short" {
		t.Fatalf("short string altered: %q", short)
	}
}

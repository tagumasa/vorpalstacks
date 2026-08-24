package sfn

import (
	"strings"
	"testing"
)

// codesOf extracts the diagnostic codes of the given severity.
func codesOf(diags []aslDiagnostic, severity string) []string {
	var codes []string
	for _, d := range diags {
		if d.Severity == severity {
			codes = append(codes, d.Code)
		}
	}
	return codes
}

func hasCode(diags []aslDiagnostic, severity, code string) bool {
	for _, c := range codesOf(diags, severity) {
		if c == code {
			return true
		}
	}
	return false
}

// TestASLValidatorDocumentedCodes pins the documented diagnostic code set:
// each minimal definition triggers exactly the code it exists for, with
// the documented "/States/<Name>/<Field>" location shape where applicable.
func TestASLValidatorDocumentedCodes(t *testing.T) {
	tests := []struct {
		name        string
		definition  string
		severity    string
		code        string
		wantLoc     string
		notWantCode string
	}{
		{
			name:       "empty definition",
			definition: "",
			severity:   "ERROR", code: "MISSING_DESCRIPTION",
		},
		{
			name:       "syntax error",
			definition: `{"StartAt":`,
			severity:   "ERROR", code: "INVALID_JSON_DESCRIPTION",
		},
		{
			name:       "unknown next",
			definition: `{"StartAt":"A","States":{"A":{"Type":"Pass","Next":"B"}}}`,
			severity:   "ERROR", code: "MISSING_TRANSITION_TARGET",
			wantLoc: "/States/A/Next",
		},
		{
			name:       "choice default target missing",
			definition: `{"StartAt":"C","States":{"C":{"Type":"Choice","Choices":[{"Variable":"$.x","StringEquals":"a","Next":"A"}],"Default":"Z"},"A":{"Type":"Succeed"}}}`,
			severity:   "ERROR", code: "MISSING_TRANSITION_TARGET",
			wantLoc: "/States/C/Default",
		},
		{
			name:       "no terminal state",
			definition: `{"StartAt":"A","States":{"A":{"Type":"Pass","Next":"B"},"B":{"Type":"Pass","Next":"A"}}}`,
			severity:   "ERROR", code: "MISSING_END_STATE",
		},
		{
			name:       "duplicate state key in raw json",
			definition: `{"StartAt":"A","States":{"A":{"Type":"Succeed"},"A":{"Type":"Fail"}}}`,
			severity:   "ERROR", code: "DUPLICATE_STATE_NAME",
			wantLoc: "/States/A",
		},
		{
			name:       "duplicate state name across scopes",
			definition: `{"StartAt":"P","States":{"P":{"Type":"Parallel","End":true,"Branches":[{"StartAt":"A","States":{"A":{"Type":"Succeed"}}},{"StartAt":"A","States":{"A":{"Type":"Succeed"}}}]}}}`,
			severity:   "ERROR", code: "DUPLICATE_STATE_NAME",
		},
		{
			name:       "state name too long",
			definition: `{"StartAt":"` + strings.Repeat("n", 81) + `","States":{"` + strings.Repeat("n", 81) + `":{"Type":"Succeed"}}}`,
			severity:   "ERROR", code: "INVALID_STATE_NAME",
		},
		{
			name:       "invalid task resource",
			definition: `{"StartAt":"T","States":{"T":{"Type":"Task","Resource":"not-an-arn","End":true}}}`,
			severity:   "ERROR", code: "INVALID_RESOURCE",
			wantLoc: "/States/T/Resource",
		},
		{
			name:       "invalid label characters",
			definition: `{"StartAt":"M","States":{"M":{"Type":"Map","Label":"bad label","ItemProcessor":{"StartAt":"S","States":{"S":{"Type":"Succeed"}}},"End":true}}}`,
			severity:   "ERROR", code: "INVALID_LABEL_NAME",
			wantLoc: "/States/M/Label",
		},
		{
			name:       "duplicate labels",
			definition: `{"StartAt":"M","States":{"M":{"Type":"Map","Label":"dup","ItemProcessor":{"StartAt":"S","States":{"S":{"Type":"Succeed"}}},"Next":"N"},"N":{"Type":"Map","Label":"dup","ItemProcessor":{"StartAt":"S2","States":{"S2":{"Type":"Succeed"}}},"End":true}}}`,
			severity:   "ERROR", code: "DUPLICATE_LABEL_NAME",
		},
		{
			name:       "missing startat",
			definition: `{"States":{"A":{"Type":"Succeed"}}}`,
			severity:   "ERROR", code: "SCHEMA_VALIDATION_FAILED",
			wantLoc: "/StartAt",
		},
		{
			name:       "missing states",
			definition: `{"StartAt":"A"}`,
			severity:   "ERROR", code: "SCHEMA_VALIDATION_FAILED",
			wantLoc: "/States",
		},
		{
			name:       "startat not a state",
			definition: `{"StartAt":"Z","States":{"A":{"Type":"Succeed"}}}`,
			severity:   "ERROR", code: "SCHEMA_VALIDATION_FAILED",
		},
		{
			name:       "missing type",
			definition: `{"StartAt":"A","States":{"A":{"End":true}}}`,
			severity:   "ERROR", code: "SCHEMA_VALIDATION_FAILED",
		},
		{
			name:       "invalid type",
			definition: `{"StartAt":"A","States":{"A":{"Type":"Queue","End":true}}}`,
			severity:   "ERROR", code: "SCHEMA_VALIDATION_FAILED",
		},
		{
			name:       "task without resource",
			definition: `{"StartAt":"T","States":{"T":{"Type":"Task","End":true}}}`,
			severity:   "ERROR", code: "SCHEMA_VALIDATION_FAILED",
		},
		{
			name:       "next and end together",
			definition: `{"StartAt":"A","States":{"A":{"Type":"Pass","Next":"B","End":true},"B":{"Type":"Succeed"}}}`,
			severity:   "ERROR", code: "SCHEMA_VALIDATION_FAILED",
		},
		{
			name:       "choice without choices",
			definition: `{"StartAt":"C","States":{"C":{"Type":"Choice","Default":"A"},"A":{"Type":"Succeed"}}}`,
			severity:   "ERROR", code: "SCHEMA_VALIDATION_FAILED",
		},
		{
			name:       "map without iterator",
			definition: `{"StartAt":"M","States":{"M":{"Type":"Map","End":true}}}`,
			severity:   "ERROR", code: "SCHEMA_VALIDATION_FAILED",
		},
		{
			name:       "map with both iterator forms",
			definition: `{"StartAt":"M","States":{"M":{"Type":"Map","Iterator":{"StartAt":"S","States":{"S":{"Type":"Succeed"}}},"ItemProcessor":{"StartAt":"S2","States":{"S2":{"Type":"Succeed"}}},"End":true}}}`,
			severity:   "ERROR", code: "SCHEMA_VALIDATION_FAILED",
		},
		{
			name:       "parallel without branches",
			definition: `{"StartAt":"P","States":{"P":{"Type":"Parallel","End":true}}}`,
			severity:   "ERROR", code: "SCHEMA_VALIDATION_FAILED",
		},
		{
			name:       "itemreader unsupported resource",
			definition: `{"StartAt":"M","States":{"M":{"Type":"Map","ItemReader":{"Resource":"arn:aws:states:::dynamodb:getItem","Parameters":{"Bucket":"b","Key":"k"}},"ItemProcessor":{"StartAt":"S","States":{"S":{"Type":"Succeed"}}},"End":true}}}`,
			severity:   "ERROR", code: "SCHEMA_VALIDATION_FAILED",
			wantLoc: "/States/M/ItemReader/Resource",
		},
		{
			name:       "itemreader maxitems and path together",
			definition: `{"StartAt":"M","States":{"M":{"Type":"Map","ItemReader":{"Resource":"arn:aws:states:::s3:getObject","Parameters":{"Bucket":"b","Key":"k"},"ReaderConfig":{"MaxItems":5,"MaxItemsPath":"$.m"}},"ItemProcessor":{"StartAt":"S","States":{"S":{"Type":"Succeed"}}},"End":true}}}`,
			severity:   "ERROR", code: "SCHEMA_VALIDATION_FAILED",
		},
		{
			name:       "resultwriter without parameters",
			definition: `{"StartAt":"M","States":{"M":{"Type":"Map","ResultWriter":{"Resource":"arn:aws:states:::s3:putObject"},"ItemProcessor":{"StartAt":"S","States":{"S":{"Type":"Succeed"}}},"End":true}}}`,
			severity:   "ERROR", code: "SCHEMA_VALIDATION_FAILED",
		},
		{
			name:       "tolerated percentage out of range",
			definition: `{"StartAt":"M","States":{"M":{"Type":"Map","ToleratedFailurePercentage":150,"ItemProcessor":{"StartAt":"S","States":{"S":{"Type":"Succeed"}}},"End":true}}}`,
			severity:   "ERROR", code: "SCHEMA_VALIDATION_FAILED",
		},
		{
			name:       "retry backoff below one",
			definition: `{"StartAt":"T","States":{"T":{"Type":"Task","Resource":"arn:aws:lambda:us-east-1:000000000000:function:f","Retry":[{"ErrorEquals":["States.ALL"],"BackoffRate":0.5}],"End":true}}}`,
			severity:   "ERROR", code: "SCHEMA_VALIDATION_FAILED",
		},
		{
			name:       "distributed mode in express machine",
			definition: `{"StartAt":"M","States":{"M":{"Type":"Map","ItemProcessor":{"ProcessorConfig":{"Mode":"DISTRIBUTED","ExecutionType":"STANDARD"},"StartAt":"S","States":{"S":{"Type":"Succeed"}}},"End":true}}}`,
			severity:   "ERROR", code: "SCHEMA_VALIDATION_FAILED",
		},
		{
			name:        "valid minimal definition produces no errors",
			definition:  `{"StartAt":"A","States":{"A":{"Type":"Pass","Result":{"x":1},"End":true}}}`,
			notWantCode: "SCHEMA_VALIDATION_FAILED",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diags := validateASLStructure(tc.definition, "EXPRESS")
			if tc.notWantCode != "" {
				for _, d := range diags {
					if d.Code == tc.notWantCode && d.Severity == "ERROR" {
						t.Fatalf("unexpected %s: %+v", tc.notWantCode, d)
					}
				}
				return
			}
			if !hasCode(diags, tc.severity, tc.code) {
				t.Fatalf("%s diagnostic missing (severity %s), got %+v", tc.code, tc.severity, diags)
			}
			if tc.wantLoc != "" {
				found := false
				for _, d := range diags {
					if d.Code == tc.code && d.Location == tc.wantLoc {
						found = true
					}
				}
				if !found {
					t.Fatalf("%s at location %s missing, got %+v", tc.code, tc.wantLoc, diags)
				}
			}
		})
	}
}

// TestASLValidatorDistributedModeAllowedForStandard pins that the same
// Distributed definition passes for a STANDARD machine.
func TestASLValidatorDistributedModeAllowedForStandard(t *testing.T) {
	def := `{"StartAt":"M","States":{"M":{"Type":"Map","ItemProcessor":{"ProcessorConfig":{"Mode":"DISTRIBUTED","ExecutionType":"STANDARD"},"StartAt":"S","States":{"S":{"Type":"Succeed"}}},"End":true}}}`
	if diags := validateASLStructure(def, "STANDARD"); hasCode(diags, "ERROR", "SCHEMA_VALIDATION_FAILED") {
		t.Fatalf("standard machine rejected: %+v", diags)
	}
}

// TestASLValidatorWarnings pins the documented warning codes: NO_PATH for
// path-looking values under field names without the Path suffix, NO_DOLLAR
// for intrinsic-looking values, PASS_RESULT_IS_STATIC for a path-looking
// Pass Result. Warnings never fail the definition.
func TestASLValidatorWarnings(t *testing.T) {
	def := `{"StartAt":"T","States":{"T":{"Type":"Task","Resource":"arn:aws:lambda:us-east-1:000000000000:function:f","Parameters":{"plain":"$.input","intrinsic":"States.Format('{}')"},"End":true}}}`
	diags := validateASLStructure(def, "STANDARD")
	if !hasCode(diags, "WARNING", "NO_PATH") {
		t.Fatalf("NO_PATH warning missing: %+v", diags)
	}
	if !hasCode(diags, "WARNING", "NO_DOLLAR") {
		t.Fatalf("NO_DOLLAR warning missing: %+v", diags)
	}
	for _, d := range diags {
		if d.Severity == "ERROR" {
			t.Fatalf("warning-only definition produced an error: %+v", d)
		}
	}

	passDef := `{"StartAt":"P","States":{"P":{"Type":"Pass","Result":"$.static","End":true}}}`
	diags = validateASLStructure(passDef, "STANDARD")
	if !hasCode(diags, "WARNING", "PASS_RESULT_IS_STATIC") {
		t.Fatalf("PASS_RESULT_IS_STATIC warning missing: %+v", diags)
	}
}

// TestValidateDefinitionStructureCreateRejection pins the creation-path
// behaviour: ERROR diagnostics reject with the InvalidDefinition shape and
// warning-only definitions are accepted.
func TestValidateDefinitionStructureCreateRejection(t *testing.T) {
	if err := validateDefinitionStructure(`{"StartAt":"A","States":{"A":{"Type":"Pass","Next":"X"}}}`, "STANDARD"); err == nil {
		t.Fatal("broken definition accepted at creation time")
	}
	if err := validateDefinitionStructure(`{"StartAt":"A","States":{"A":{"Type":"Pass","Result":"$.w","End":true}}}`, "STANDARD"); err != nil {
		t.Fatalf("warning-only definition rejected at creation time: %v", err)
	}
}

// TestASLValidatorItemBatcher pins the ItemBatcher structural contract:
// the sizing pairs are mutually exclusive, at least one sizing value is
// required, the byte cap stays within the 256 KiB child-execution input
// bound and the fixed BatchInput is an object.
func TestASLValidatorItemBatcher(t *testing.T) {
	valid := `{"StartAt":"M","States":{"M":{"Type":"Map","ItemProcessor":{"StartAt":"S","States":{"S":{"Type":"Succeed"}}},"ItemBatcher":{"MaxItemsPerBatch":500,"BatchInput":{"k":"v"}},"End":true}}}`
	if diags := validateASLStructure(valid, "STANDARD"); len(codesOf(diags, "ERROR")) != 0 {
		t.Fatalf("valid ItemBatcher produced errors: %+v", diags)
	}

	cases := []struct {
		name       string
		batcher    string
		wantLoc    string
		wantSubstr string
	}{
		{
			name:       "both max item forms",
			batcher:    `{"MaxItemsPerBatch":5,"MaxItemsPerBatchPath":"$.m"}`,
			wantLoc:    "/States/M/ItemBatcher/MaxItemsPerBatch",
			wantSubstr: "both MaxItemsPerBatch and MaxItemsPerBatchPath",
		},
		{
			name:       "both byte cap forms",
			batcher:    `{"MaxInputBytesPerBatch":100,"MaxInputBytesPerBatchPath":"$.b"}`,
			wantLoc:    "/States/M/ItemBatcher/MaxInputBytesPerBatch",
			wantSubstr: "both MaxInputBytesPerBatch and MaxInputBytesPerBatchPath",
		},
		{
			name:       "no sizing value",
			batcher:    `{"BatchInput":{"k":"v"}}`,
			wantLoc:    "/States/M/ItemBatcher",
			wantSubstr: "must specify MaxItemsPerBatch, MaxInputBytesPerBatch or both",
		},
		{
			name:       "byte cap over the 256 kib bound",
			batcher:    `{"MaxInputBytesPerBatch":262145}`,
			wantLoc:    "/States/M/ItemBatcher/MaxInputBytesPerBatch",
			wantSubstr: "262144",
		},
		{
			name:       "byte cap not a positive integer",
			batcher:    `{"MaxInputBytesPerBatch":0}`,
			wantLoc:    "/States/M/ItemBatcher/MaxInputBytesPerBatch",
			wantSubstr: "must be an integer from 1 to 262144",
		},
		{
			name:       "item count not positive",
			batcher:    `{"MaxItemsPerBatch":0}`,
			wantLoc:    "/States/M/ItemBatcher/MaxItemsPerBatch",
			wantSubstr: "must be a positive integer",
		},
		{
			name:       "batch input not an object",
			batcher:    `{"MaxItemsPerBatch":5,"BatchInput":"scalar"}`,
			wantLoc:    "/States/M/ItemBatcher/BatchInput",
			wantSubstr: "must be a JSON object",
		},
		{
			name:       "batch input both forms",
			batcher:    `{"MaxItemsPerBatch":5,"BatchInput":{"k":"v"},"BatchInputPath":"$.f"}`,
			wantLoc:    "/States/M/ItemBatcher/BatchInput",
			wantSubstr: "both BatchInput and BatchInputPath",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			def := `{"StartAt":"M","States":{"M":{"Type":"Map","ItemProcessor":{"StartAt":"S","States":{"S":{"Type":"Succeed"}}},"ItemBatcher":` + tc.batcher + `,"End":true}}}`
			diags := validateASLStructure(def, "STANDARD")
			found := false
			for _, d := range diags {
				if d.Severity == "ERROR" && d.Code == "SCHEMA_VALIDATION_FAILED" &&
					d.Location == tc.wantLoc && strings.Contains(d.Message, tc.wantSubstr) {
					found = true
				}
			}
			if !found {
				t.Fatalf("expected SCHEMA_VALIDATION_FAILED at %s mentioning %q, got %+v", tc.wantLoc, tc.wantSubstr, diags)
			}
		})
	}
}

// TestASLValidatorUnknownStateFields pins the closed-world state-field
// census: members the state's type does not declare are
// SCHEMA_VALIDATION_FAILED errors and reject at creation time, while the
// documented direct members (including the reference-path forms) stay
// legal.
func TestASLValidatorUnknownStateFields(t *testing.T) {
	unknown := `{"StartAt":"P","States":{"P":{"Type":"Pass","ResultPat":"$.x","End":true}}}`
	diags := validateASLStructure(unknown, "STANDARD")
	found := false
	for _, d := range diags {
		if d.Severity == "ERROR" && d.Code == "SCHEMA_VALIDATION_FAILED" &&
			d.Location == "/States/P/ResultPat" && strings.Contains(d.Message, "unknown field") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected an unknown-field diagnostic at /States/P/ResultPat, got %+v", diags)
	}
	if err := validateDefinitionStructure(unknown, "STANDARD"); err == nil {
		t.Fatal("definition with an unknown state field accepted at creation time")
	}

	failInputPath := `{"StartAt":"F","States":{"F":{"Type":"Fail","InputPath":"$.x"}}}`
	diags = validateASLStructure(failInputPath, "STANDARD")
	found = false
	for _, d := range diags {
		if d.Severity == "ERROR" && d.Location == "/States/F/InputPath" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected Fail InputPath rejection, got %+v", diags)
	}

	failBoth := `{"StartAt":"F","States":{"F":{"Type":"Fail","Error":"E","ErrorPath":"$.e","Cause":"C","CausePath":"$.c"}}}`
	diags = validateASLStructure(failBoth, "STANDARD")
	errorOK, causeOK := false, false
	for _, d := range diags {
		if d.Location == "/States/F/ErrorPath" {
			errorOK = true
		}
		if d.Location == "/States/F/CausePath" {
			causeOK = true
		}
	}
	if !errorOK || !causeOK {
		t.Fatalf("expected Error and Cause literal/path exclusivity diagnostics, got %+v", diags)
	}

	taskBoth := `{"StartAt":"T","States":{"T":{"Type":"Task","Resource":"arn:aws:states:::lambda:invoke","TimeoutSeconds":10,"TimeoutSecondsPath":"$.t","End":true}}}`
	diags = validateASLStructure(taskBoth, "STANDARD")
	found = false
	for _, d := range diags {
		if d.Location == "/States/T/TimeoutSecondsPath" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected TimeoutSeconds/TimeoutSecondsPath exclusivity diagnostic, got %+v", diags)
	}

	legal := `{"StartAt":"M","States":{"M":{"Type":"Map","ItemProcessor":{"StartAt":"S","States":{"S":{"Type":"Succeed"}}},"MaxConcurrencyPath":"$.mc","End":true}}}`
	if diags := validateASLStructure(legal, "STANDARD"); len(codesOf(diags, "ERROR")) != 0 {
		t.Fatalf("MaxConcurrencyPath must stay legal: %+v", diags)
	}
	legalFail := `{"StartAt":"F","States":{"F":{"Type":"Fail","ErrorPath":"$.e","CausePath":"$.c"}}}`
	if diags := validateASLStructure(legalFail, "STANDARD"); len(codesOf(diags, "ERROR")) != 0 {
		t.Fatalf("Fail reference paths must stay legal: %+v", diags)
	}
	legalTask := `{"StartAt":"T","States":{"T":{"Type":"Task","Resource":"arn:aws:states:::lambda:invoke","TimeoutSecondsPath":"$.t","HeartbeatSecondsPath":"$.h","End":true}}}`
	if diags := validateASLStructure(legalTask, "STANDARD"); len(codesOf(diags, "ERROR")) != 0 {
		t.Fatalf("Task timeout reference paths must stay legal: %+v", diags)
	}
}

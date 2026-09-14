package sfn

import (
	"fmt"
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
			// "Express Workflows only support Request Response
			// integrations" — the .sync suffix names a pattern beyond
			// Request Response and is rejected at validation.
			name:       "express machine with sync pattern resource",
			definition: `{"StartAt":"T","States":{"T":{"Type":"Task","Resource":"arn:aws:states:::states:startExecution.sync","Parameters":{"StateMachineArn":"arn:aws:states:us-east-1:000000000000:stateMachine:x"},"End":true}}}`,
			severity:   "ERROR", code: "SCHEMA_VALIDATION_FAILED",
			wantLoc: "/States/T/Resource",
		},
		{
			// A state-level Input member is not part of the documented
			// common field set — the census rejects it on both dialects.
			name:       "state-level input member",
			definition: `{"StartAt":"A","States":{"A":{"Type":"Pass","Input":{"y":2},"End":true}}}`,
			severity:   "ERROR", code: "SCHEMA_VALIDATION_FAILED",
			wantLoc: "/States/A/Input",
		},
		{
			// A payload-template ".$" key carries a path string; any other
			// JSON type is an invalid template.
			name:       "payload template non-string path key",
			definition: `{"StartAt":"T","States":{"T":{"Type":"Task","Resource":"arn:aws:lambda:us-east-1:000000000000:function:f","Parameters":{"num.$":5},"End":true}}}`,
			severity:   "ERROR", code: "SCHEMA_VALIDATION_FAILED",
			wantLoc: "/States/T/Parameters/num.$",
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

// TestASLValidatorPatternAllowedForStandard pins the counterpart of the
// Express integration-pattern gate: a STANDARD machine accepts the
// .sync pattern the gate rejects on Express workflows.
func TestASLValidatorPatternAllowedForStandard(t *testing.T) {
	def := `{"StartAt":"T","States":{"T":{"Type":"Task","Resource":"arn:aws:states:::states:startExecution.sync","Parameters":{"StateMachineArn":"arn:aws:states:us-east-1:000000000000:stateMachine:x"},"End":true}}}`
	if diags := validateASLStructure(def, "STANDARD"); hasCode(diags, "ERROR", "SCHEMA_VALIDATION_FAILED") {
		t.Fatalf("standard machine rejected a .sync resource: %+v", diags)
	}
}

// TestASLValidatorCatcherDialectSeparation pins the Catch member
// separation: a JSONata Catcher folds its error through Assign/Output and
// must not carry the JSONPath ResultPath, while a JSONPath Catcher must
// not carry the JSONata Assign/Output.
func TestASLValidatorCatcherDialectSeparation(t *testing.T) {
	jsonataResultPath := `{"QueryLanguage":"JSONata","StartAt":"T","States":{"T":{"Type":"Task","Resource":"arn:aws:lambda:us-east-1:000000000000:function:f","Catch":[{"ErrorEquals":["States.ALL"],"ResultPath":"$.err","Next":"A"}],"Next":"A"},"A":{"Type":"Succeed"}}}`
	if diags := validateASLStructure(jsonataResultPath, "STANDARD"); !hasCode(diags, "ERROR", "SCHEMA_VALIDATION_FAILED") {
		t.Fatalf("JSONata Catcher with ResultPath accepted: %+v", diags)
	}

	jsonPathOutput := `{"StartAt":"T","States":{"T":{"Type":"Task","Resource":"arn:aws:lambda:us-east-1:000000000000:function:f","Catch":[{"ErrorEquals":["States.ALL"],"Output":{"x":1},"Next":"A"}],"Next":"A"},"A":{"Type":"Succeed"}}}`
	if diags := validateASLStructure(jsonPathOutput, "STANDARD"); !hasCode(diags, "ERROR", "SCHEMA_VALIDATION_FAILED") {
		t.Fatalf("JSONPath Catcher with Output accepted: %+v", diags)
	}

	jsonPathValid := `{"StartAt":"T","States":{"T":{"Type":"Task","Resource":"arn:aws:lambda:us-east-1:000000000000:function:f","Catch":[{"ErrorEquals":["States.ALL"],"ResultPath":"$.err","Next":"A"}],"Next":"A"},"A":{"Type":"Succeed"}}}`
	if diags := validateASLStructure(jsonPathValid, "STANDARD"); hasCode(diags, "ERROR", "SCHEMA_VALIDATION_FAILED") {
		t.Fatalf("JSONPath Catcher with ResultPath rejected: %+v", diags)
	}
}

// TestASLValidatorJSONataNestedChoiceRules pins that nested Choice rules
// (inside And/Or/Not combinators) receive the same JSONata validation as
// top-level ones: an improperly delimited Condition and an inaccessible
// $states.result reference are both caught at creation no matter the
// nesting depth.
func TestASLValidatorJSONataNestedChoiceRules(t *testing.T) {
	badDelimiter := `{"QueryLanguage":"JSONata","StartAt":"C","States":{"C":{"Type":"Choice","Choices":[{"And":[{"Condition":"{% $x > 1 %}"},{"Condition":"broken {% $y"}],"Default":"A"}],"Default":"A"},"A":{"Type":"Succeed"}}}`
	if diags := validateASLStructure(badDelimiter, "STANDARD"); !hasCode(diags, "ERROR", "SCHEMA_VALIDATION_FAILED") {
		t.Fatalf("nested rule with an improperly delimited Condition accepted: %+v", diags)
	}

	inaccessibleResult := `{"QueryLanguage":"JSONata","StartAt":"C","States":{"C":{"Type":"Choice","Choices":[{"Not":{"Condition":"{% $states.result > 1 %}"}}],"Default":"A"},"A":{"Type":"Succeed"}}}`
	if diags := validateASLStructure(inaccessibleResult, "STANDARD"); !hasCode(diags, "ERROR", "SCHEMA_VALIDATION_FAILED") {
		t.Fatalf("nested rule referencing $states.result in a Choice accepted: %+v", diags)
	}

	valid := `{"QueryLanguage":"JSONata","StartAt":"C","States":{"C":{"Type":"Choice","Choices":[{"And":[{"Condition":"{% $x > 1 %}"},{"Condition":"{% $y < 2 %}"}],"Next":"A"}],"Default":"A"},"A":{"Type":"Succeed"}}}`
	if diags := validateASLStructure(valid, "STANDARD"); hasCode(diags, "ERROR", "SCHEMA_VALIDATION_FAILED") {
		t.Fatalf("well-formed nested rules rejected: %+v", diags)
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

// TestPayloadTemplateCensusCoversItemReaderAndResultWriter pins that the
// creation-time ".$"-key census reaches the Distributed Map's nested
// payload templates: a non-string value under ItemReader.Parameters or
// ResultWriter.Parameters is rejected at creation, not deferred to the
// run-time States.ParameterPathFailure path.
func TestPayloadTemplateCensusCoversItemReaderAndResultWriter(t *testing.T) {
	readerBad := `{"StartAt":"M","States":{"M":{"Type":"Map","ItemProcessor":{"StartAt":"S","States":{"S":{"Type":"Succeed"}}},"ItemReader":{"Resource":"arn:aws:states:::s3:getObject","Parameters":{"Bucket.$":5}},"End":true}}}`
	if diags := validateASLStructure(readerBad, "STANDARD"); !hasCode(diags, "ERROR", "SCHEMA_VALIDATION_FAILED") {
		t.Fatalf("ItemReader.Parameters non-string path value accepted: %+v", diags)
	}

	writerBad := `{"StartAt":"M","States":{"M":{"Type":"Map","ItemProcessor":{"StartAt":"S","States":{"S":{"Type":"Succeed"}}},"ResultWriter":{"Resource":"arn:aws:states:::s3:putObject","Parameters":{"Prefix.$":7}},"End":true}}}`
	if diags := validateASLStructure(writerBad, "STANDARD"); !hasCode(diags, "ERROR", "SCHEMA_VALIDATION_FAILED") {
		t.Fatalf("ResultWriter.Parameters non-string path value accepted: %+v", diags)
	}

	readerGood := `{"StartAt":"M","States":{"M":{"Type":"Map","ItemProcessor":{"StartAt":"S","States":{"S":{"Type":"Succeed"}}},"ItemReader":{"Resource":"arn:aws:states:::s3:getObject","Parameters":{"Bucket.$":"$.bucket"}},"End":true}}}`
	if diags := validateASLStructure(readerGood, "STANDARD"); hasCode(diags, "ERROR", "SCHEMA_VALIDATION_FAILED") {
		t.Fatalf("well-formed ItemReader.Parameters path string rejected: %+v", diags)
	}
}

// TestSingleDiagnosisPerLocation pins the one-diagnosis-per-defect rule: a
// path-looking Pass Result carries PASS_RESULT_IS_STATIC alone (the scan
// walk does not stack NO_PATH at the same location), and a Fail state's
// Output member is reported once by the unknown-field census (not again by
// a Fail-specific member list).
func TestSingleDiagnosisPerLocation(t *testing.T) {
	passDef := `{"StartAt":"P","States":{"P":{"Type":"Pass","Result":"$.static","End":true}}}`
	diags := validateASLStructure(passDef, "STANDARD")
	passResultDiags := 0
	for _, d := range diags {
		if d.Location == "/States/P/Result" {
			passResultDiags++
		}
	}
	if passResultDiags != 1 {
		t.Fatalf("Pass Result produced %d diagnostics at /States/P/Result, want exactly the PASS_RESULT_IS_STATIC one: %+v", passResultDiags, diags)
	}

	failDef := `{"StartAt":"F","States":{"F":{"Type":"Fail","Error":"E","Cause":"c","Output":"{}"}}}`
	diags = validateASLStructure(failDef, "STANDARD")
	failOutputDiags := 0
	for _, d := range diags {
		if d.Location == "/States/F/Output" {
			failOutputDiags++
		}
	}
	if failOutputDiags != 1 {
		t.Fatalf("Fail Output produced %d diagnostics at /States/F/Output, want exactly the unknown-field rejection: %+v", failOutputDiags, diags)
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

// TestASLValidatorCatchEdgeSatisfiesTerminalReachability pins that the
// terminal-reachability walk traverses Catch edges: a Catcher's Next is an
// ordinary transition, so a workflow whose terminal state is reachable
// only through a catch target is legal, while a genuinely terminal-less
// workflow is still rejected.
func TestASLValidatorCatchEdgeSatisfiesTerminalReachability(t *testing.T) {
	catchTerminal := `{"StartAt":"T","States":{"T":{"Type":"Task","Resource":"arn:aws:lambda:us-east-1:000000000000:function:f","Next":"T","Catch":[{"ErrorEquals":["States.ALL"],"Next":"Handle"}]},"Handle":{"Type":"Succeed"}}}`
	if diags := validateASLStructure(catchTerminal, "STANDARD"); hasCode(diags, "ERROR", "MISSING_END_STATE") {
		t.Fatalf("catch-terminal workflow rejected as terminal-less: %+v", diags)
	}

	terminalLess := `{"StartAt":"A","States":{"A":{"Type":"Pass","Next":"A"}}}`
	if diags := validateASLStructure(terminalLess, "STANDARD"); !hasCode(diags, "ERROR", "MISSING_END_STATE") {
		t.Fatalf("a workflow with no terminal state at all must stay rejected: %+v", diags)
	}
}

// TestASLValidatorChoiceRuleMixesRejected pins the three Choice-rule
// validation gaps: a combinator combined with a comparator, a non-object
// And/Or entry, and a JSONata rule combining a combinator with the
// JSONPath-only Variable field are all schema failures at creation.
func TestASLValidatorChoiceRuleMixesRejected(t *testing.T) {
	mix := `{"StartAt":"C","States":{"C":{"Type":"Choice","Choices":[{"And":[{"Variable":"$.a","IsString":true}],"StringEquals":"x","Next":"A"}],"Default":"A"},"A":{"Type":"Succeed"}}}`
	if diags := validateASLStructure(mix, "STANDARD"); !hasCode(diags, "ERROR", "SCHEMA_VALIDATION_FAILED") {
		t.Fatalf("combinator+comparator rule accepted: %+v", diags)
	}

	nonObjectEntry := `{"StartAt":"C","States":{"C":{"Type":"Choice","Choices":[{"And":["foo"],"Next":"A"}],"Default":"A"},"A":{"Type":"Succeed"}}}`
	if diags := validateASLStructure(nonObjectEntry, "STANDARD"); !hasCode(diags, "ERROR", "SCHEMA_VALIDATION_FAILED") {
		t.Fatalf("non-object And entry accepted: %+v", diags)
	}

	jsonataVariableCombinator := `{"QueryLanguage":"JSONata","StartAt":"C","States":{"C":{"Type":"Choice","Choices":[{"And":[{"Condition":"{% $x > 1 %}"}],"Variable":"$.a","Next":"A"}],"Default":"A"},"A":{"Type":"Succeed"}}}`
	if diags := validateASLStructure(jsonataVariableCombinator, "STANDARD"); !hasCode(diags, "ERROR", "SCHEMA_VALIDATION_FAILED") {
		t.Fatalf("JSONata combinator with the JSONPath-only Variable accepted: %+v", diags)
	}
}

// TestASLValidatorJSONataChoiceRejectsOutputPath pins the dialect rule on
// Choice: OutputPath is a JSONPath-only member ("it is invalid to use
// path-based fields when using JSONata"), so a JSONata Choice carrying it
// fails at creation instead of being silently inert, while the JSONPath
// form keeps it.
func TestASLValidatorJSONataChoiceRejectsOutputPath(t *testing.T) {
	jsonata := `{"QueryLanguage":"JSONata","StartAt":"C","States":{"C":{"Type":"Choice","Choices":[{"Condition":"{% $x > 1 %}","Next":"A"}],"OutputPath":"$.o","Default":"A"},"A":{"Type":"Succeed"}}}`
	if diags := validateASLStructure(jsonata, "STANDARD"); !hasCode(diags, "ERROR", "SCHEMA_VALIDATION_FAILED") {
		t.Fatalf("JSONata Choice with OutputPath accepted: %+v", diags)
	}

	jsonPath := `{"StartAt":"C","States":{"C":{"Type":"Choice","Choices":[{"Variable":"$.x","NumericEquals":1,"Next":"A"}],"OutputPath":"$.o","Default":"A"},"A":{"Type":"Succeed"}}}`
	if diags := validateASLStructure(jsonPath, "STANDARD"); hasCode(diags, "ERROR", "SCHEMA_VALIDATION_FAILED") {
		t.Fatalf("JSONPath Choice with OutputPath rejected: %+v", diags)
	}
}

// TestASLValidatorStateLevelConditionRejected pins that a state-level
// Condition on Choice is unknown on both dialects — Condition is
// documented only inside Choice rules — so the census rejects it instead
// of admitting an inert member (JSONata) or citing a nonexistent
// JSONata-only contract (JSONPath).
func TestASLValidatorStateLevelConditionRejected(t *testing.T) {
	jsonata := `{"QueryLanguage":"JSONata","StartAt":"C","States":{"C":{"Type":"Choice","Condition":"{% $x > 1 %}","Choices":[{"Condition":"{% $y > 2 %}","Next":"A"}],"Default":"A"},"A":{"Type":"Succeed"}}}`
	if diags := validateASLStructure(jsonata, "STANDARD"); !hasCode(diags, "ERROR", "SCHEMA_VALIDATION_FAILED") {
		t.Fatalf("JSONata state-level Condition accepted: %+v", diags)
	}

	jsonPath := `{"StartAt":"C","States":{"C":{"Type":"Choice","Condition":"$.x","Choices":[{"Variable":"$.x","NumericEquals":1,"Next":"A"}],"Default":"A"},"A":{"Type":"Succeed"}}}`
	if diags := validateASLStructure(jsonPath, "STANDARD"); !hasCode(diags, "ERROR", "SCHEMA_VALIDATION_FAILED") {
		t.Fatalf("JSONPath state-level Condition accepted: %+v", diags)
	}
}

// TestASLValidatorNotEqualsOperatorsRejected pins the comparator
// vocabulary: the documented supported-operator list (specification and
// Choice page agree) has no NotEquals variants, so a rule using one is a
// schema failure at creation instead of a rule that silently never
// matches.
func TestASLValidatorNotEqualsOperatorsRejected(t *testing.T) {
	for _, operator := range []string{"StringNotEquals", "StringNotEqualsPath", "NumericNotEquals", "NumericNotEqualsPath"} {
		def := fmt.Sprintf(`{"StartAt":"C","States":{"C":{"Type":"Choice","Choices":[`+
			`{"Variable":"$.x","%s":%s,"Next":"A"}],"Default":"A"},"A":{"Type":"Succeed"}}}`,
			operator, map[bool]string{true: `"a"`, false: `"$.a"`}[strings.HasSuffix(operator, "Path")])
		if diags := validateASLStructure(def, "STANDARD"); !hasCode(diags, "ERROR", "SCHEMA_VALIDATION_FAILED") {
			t.Fatalf("%s accepted: %+v", operator, diags)
		}
	}
}

// TestASLValidatorChoiceComparatorValues pins the value typing of the
// Choice comparators at definition validation: a literal whose JSON type
// contradicts its operator, or a timestamp literal outside the RFC3339
// profile with the uppercase-T and uppercase-Z restrictions, is a schema
// failure at creation instead of a rule that silently never matches at
// run time. The documented well-typed forms — including fractional
// seconds and a numeric offset — stay legal, and the typing applies
// inside combinators through the recursive rule walk.
func TestASLValidatorChoiceComparatorValues(t *testing.T) {
	cases := []struct {
		name    string
		rule    string
		wantRej bool
	}{
		{name: "string operator with number", rule: `{"Variable":"$.v","StringEquals":123,"Next":"A"}`, wantRej: true},
		{name: "string operator with boolean", rule: `{"Variable":"$.v","StringGreaterThan":true,"Next":"A"}`, wantRej: true},
		{name: "StringMatches with number", rule: `{"Variable":"$.v","StringMatches":5,"Next":"A"}`, wantRej: true},
		{name: "numeric operator with string", rule: `{"Variable":"$.v","NumericEquals":"1","Next":"A"}`, wantRej: true},
		{name: "boolean operator with string", rule: `{"Variable":"$.v","BooleanEquals":"yes","Next":"A"}`, wantRej: true},
		{name: "IsPresent with string", rule: `{"Variable":"$.v","IsPresent":"true","Next":"A"}`, wantRej: true},
		{name: "timestamp operator with number", rule: `{"Variable":"$.v","TimestampEquals":123,"Next":"A"}`, wantRej: true},
		{name: "timestamp literal not a timestamp", rule: `{"Variable":"$.v","TimestampEquals":"not-a-timestamp","Next":"A"}`, wantRej: true},
		{name: "timestamp literal lowercase T", rule: `{"Variable":"$.v","TimestampEquals":"2001-01-01t12:00:00Z","Next":"A"}`, wantRej: true},
		{name: "timestamp literal lowercase Z", rule: `{"Variable":"$.v","TimestampEquals":"2001-01-01T12:00:00z","Next":"A"}`, wantRej: true},
		{name: "timestamp literal missing offset", rule: `{"Variable":"$.v","TimestampEquals":"2001-01-01T12:00:00","Next":"A"}`, wantRej: true},
		{name: "path operator with number", rule: `{"Variable":"$.v","StringEqualsPath":7,"Next":"A"}`, wantRej: true},
		{name: "Variable not a string", rule: `{"Variable":9,"StringEquals":"x","Next":"A"}`, wantRej: true},
		{name: "mistyped literal inside Not", rule: `{"Not":{"Variable":"$.v","TimestampEquals":"bad"},"Next":"A"}`, wantRej: true},
		{name: "documented string example", rule: `{"Variable":"$.foo","StringEquals":"MyString","Next":"A"}`},
		{name: "documented numeric example", rule: `{"Variable":"$.foo","NumericEquals":1,"Next":"A"}`},
		{name: "documented boolean example", rule: `{"Variable":"$.possiblyNull","IsNull":true,"Next":"A"}`},
		{name: "timestamp with uppercase Z", rule: `{"Variable":"$.foo","TimestampEquals":"2001-01-01T12:00:00Z","Next":"A"}`},
		{name: "timestamp with fractional seconds", rule: `{"Variable":"$.foo","TimestampEquals":"2001-01-01T12:00:00.500Z","Next":"A"}`},
		{name: "timestamp with numeric offset", rule: `{"Variable":"$.foo","TimestampEquals":"2001-01-01T12:00:00+09:00","Next":"A"}`},
		{name: "documented path example", rule: `{"Variable":"$.foo","StringEqualsPath":"$.bar","Next":"A"}`},
		{name: "variable reference path form", rule: `{"Variable":"$limit","StringEquals":"x","Next":"A"}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			def := `{"StartAt":"C","States":{"C":{"Type":"Choice","Choices":[` + tc.rule + `],"Default":"A"},"A":{"Type":"Succeed"}}}`
			diags := validateASLStructure(def, "STANDARD")
			rejected := hasCode(diags, "ERROR", "SCHEMA_VALIDATION_FAILED")
			if tc.wantRej && !rejected {
				t.Fatalf("rule accepted: %s — diagnostics %+v", tc.rule, diags)
			}
			if !tc.wantRej && rejected {
				t.Fatalf("legal rule rejected: %s — diagnostics %+v", tc.rule, diags)
			}
		})
	}
}

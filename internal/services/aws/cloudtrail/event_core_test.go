package cloudtrail

import (
	"testing"
	"time"

	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// TestEpochFromRawKeepsMillisecondPrecision pins the wire contract of the
// JSON-protocol time bounds: the SDKs serialise epochs with millisecond
// decimal precision, and the parse must preserve it — a bound floored to
// its whole second excludes events recorded in that same partial second.
func TestEpochFromRawKeepsMillisecondPrecision(t *testing.T) {
	tm, ok := epochFromRaw(1768687142.971)
	if !ok {
		t.Fatal("expected float64 epoch to parse")
	}
	want := time.Unix(1768687142, 971*int64(time.Millisecond)).UTC()
	if !tm.Equal(want) {
		t.Fatalf("epochFromRaw = %v, want %v", tm, want)
	}

	if _, ok := epochFromRaw("not-a-number"); ok {
		t.Fatal("expected non-float64 value to be rejected")
	}
	if _, ok := epochFromRaw(nil); ok {
		t.Fatal("expected nil value to be rejected")
	}
}

// TestParseAdvancedEventSelectorsWireForms pins the single advanced-selector
// parser's accepted wire forms: the JSON-protocol list and the JSON string
// holding the same list must both parse, and both selector-bearing resource
// families (trails, event data stores) share the parser, so the forms cannot
// drift between them.
func TestParseAdvancedEventSelectorsWireForms(t *testing.T) {
	listForm := []interface{}{
		map[string]interface{}{
			"Name": "management",
			"FieldSelectors": []interface{}{
				map[string]interface{}{
					"Field":  "eventCategory",
					"Equals": []interface{}{"Management"},
				},
			},
		},
	}
	stringForm := `[{"Name":"management","FieldSelectors":[{"Field":"eventCategory","Equals":["Management"]}]}]`

	fromList, listErr := parseAdvancedEventSelectors(listForm)
	if listErr != nil {
		t.Fatalf("list form parse failed: %v", listErr)
	}
	fromString, strErr := parseAdvancedEventSelectors(stringForm)
	if strErr != nil {
		t.Fatalf("string form parse failed: %v", strErr)
	}

	if len(fromList) != 1 || len(fromString) != 1 {
		t.Fatalf("expected one selector from each form, got list=%d string=%d", len(fromList), len(fromString))
	}
	for _, got := range [][]cloudtrailstore.AdvancedEventSelector{fromList, fromString} {
		sel := got[0]
		if sel.Name != "management" {
			t.Errorf("selector Name = %q, want %q", sel.Name, "management")
		}
		if len(sel.FieldSelectors) != 1 || sel.FieldSelectors[0].Field != "eventCategory" {
			t.Fatalf("unexpected field selectors: %+v", sel.FieldSelectors)
		}
		if eq := sel.FieldSelectors[0].Equals; len(eq) != 1 || eq[0] != "Management" {
			t.Errorf("FieldSelector Equals = %v, want [Management]", eq)
		}
	}

	if _, err := parseAdvancedEventSelectors(42); err == nil {
		t.Error("expected an error for an unsupported wire type")
	}
	if _, err := parseAdvancedEventSelectors("not json"); err == nil {
		t.Error("expected an error for a malformed JSON string")
	}
	// FieldSelectors is model-required on every selector: a selector
	// without any field selector is rejected, as is a field selector
	// without a Field.
	if _, err := parseAdvancedEventSelectors([]interface{}{map[string]interface{}{"Name": "empty"}}); err == nil {
		t.Error("expected an error for a selector with no FieldSelectors")
	}
	if _, err := parseAdvancedEventSelectors(`[{"FieldSelectors":[{"Equals":["Management"]}]}]`); err == nil {
		t.Error("expected an error for a field selector without a Field")
	}
}

package sfn

import "testing"

// TestBranchStatesInheritMachineQueryLanguage pins the spec default: "the
// default query language for each state inside the Map's ItemProcessor or
// the Parallel's Branches fields is the state machine's query language,
// and is independent of the Map or Parallel State's query language."
func TestBranchStatesInheritMachineQueryLanguage(t *testing.T) {
	jsonataParallel := `{
		"QueryLanguage": "JSONata",
		"StartAt": "P",
		"States": {
			"P": {
				"Type": "Parallel",
				"Branches": [{
					"StartAt": "A",
					"States": {
						"A": {"Type": "Pass", "Output": "{% $states.input %}", "End": true}
					}
				}],
				"End": true
			}
		}
	}`
	if err := validateDefinitionStructure(jsonataParallel, "STANDARD"); err != nil {
		t.Fatalf("a JSONata machine's branch state with JSONata fields was rejected: %v", err)
	}

	jsonataMap := `{
		"QueryLanguage": "JSONata",
		"StartAt": "M",
		"States": {
			"M": {
				"Type": "Map",
				"Items": [1, 2],
				"ItemProcessor": {
					"StartAt": "A",
					"States": {
						"A": {"Type": "Pass", "Output": "{% $states.input %}", "End": true}
					}
				},
				"End": true
			}
		}
	}`
	if err := validateDefinitionStructure(jsonataMap, "STANDARD"); err != nil {
		t.Fatalf("a JSONata machine's ItemProcessor state with JSONata fields was rejected: %v", err)
	}

	// JSONPath-only fields inside a JSONata machine's branches are
	// rejected — the accept side of the same default.
	jsonataBranchJSONPathField := `{
		"QueryLanguage": "JSONata",
		"StartAt": "P",
		"States": {
			"P": {
				"Type": "Parallel",
				"Branches": [{
					"StartAt": "A",
					"States": {
						"A": {"Type": "Pass", "InputPath": "$", "End": true}
					}
				}],
				"End": true
			}
		}
	}`
	if err := validateDefinitionStructure(jsonataBranchJSONPathField, "STANDARD"); err == nil {
		t.Fatal("a JSONPath-only field inside a JSONata machine's branch was accepted")
	}

	// A JSONPath state override inside a JSONata machine is illegal:
	// "You cannot revert a top-level JSONata-based state machine to a mix
	// of JSONata and JSONPath states."
	branchStateOverride := `{
		"QueryLanguage": "JSONata",
		"StartAt": "P",
		"States": {
			"P": {
				"Type": "Parallel",
				"Branches": [{
					"StartAt": "A",
					"States": {
						"A": {"Type": "Pass", "QueryLanguage": "JSONPath", "Parameters": {"x.$": "$.x"}, "End": true}
					}
				}],
				"End": true
			}
		}
	}`
	if err := validateDefinitionStructure(branchStateOverride, "STANDARD"); err == nil {
		t.Fatal("a JSONPath state override inside a JSONata machine was accepted")
	}
}

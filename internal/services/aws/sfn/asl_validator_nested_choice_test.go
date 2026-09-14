package sfn

import (
	"strings"
	"testing"
)

// TestNestedChoiceRulesAccepted pins the documented nested-rule contract:
// "The values of the And and Or operators must be non-empty arrays of
// Choice Rules that must not themselves contain Next fields. Likewise,
// the value of a Not operator must be a single Choice Rule that must not
// contain Next fields." and "the Next field can appear only in a top-level
// Choice Rule."
func TestNestedChoiceRulesAccepted(t *testing.T) {
	// The documented form: nested rules carry no Next, the outer rule
	// carries it.
	def := `{
		"StartAt": "Pick",
		"States": {
			"Pick": {
				"Type": "Choice",
				"Choices": [{
					"And": [
						{"Variable": "$.keyThatMightNotExist", "IsPresent": true},
						{"Variable": "$.keyThatMightNotExist", "StringEquals": "foo"}
					],
					"Next": "Yes"
				}],
				"Default": "No"
			},
			"Yes": {"Type": "Succeed"},
			"No": {"Type": "Succeed"}
		}
	}`
	if err := validateDefinitionStructure(def, "STANDARD"); err != nil {
		t.Fatalf("the documented nested And rule was rejected: %v", err)
	}

	nestedNot := `{
		"StartAt": "Pick",
		"States": {
			"Pick": {
				"Type": "Choice",
				"Choices": [{
					"Not": {"Variable": "$.a", "IsPresent": true},
					"Next": "Yes"
				}],
				"Default": "No"
			},
			"Yes": {"Type": "Succeed"},
			"No": {"Type": "Succeed"}
		}
	}`
	if err := validateDefinitionStructure(nestedNot, "STANDARD"); err != nil {
		t.Fatalf("the nested Not rule was rejected: %v", err)
	}

	// Next inside a nested rule is illegal.
	nestedNext := `{
		"StartAt": "Pick",
		"States": {
			"Pick": {
				"Type": "Choice",
				"Choices": [{
					"And": [
						{"Variable": "$.a", "IsPresent": true, "Next": "Yes"}
					],
					"Next": "Yes"
				}],
				"Default": "No"
			},
			"Yes": {"Type": "Succeed"},
			"No": {"Type": "Succeed"}
		}
	}`
	if err := validateDefinitionStructure(nestedNext, "STANDARD"); err == nil {
		t.Fatal("a Next inside a nested rule was accepted")
	}

	// And/Or must be non-empty arrays of rules.
	emptyAnd := `{
		"StartAt": "Pick",
		"States": {
			"Pick": {
				"Type": "Choice",
				"Choices": [{"And": [], "Next": "Yes"}],
				"Default": "No"
			},
			"Yes": {"Type": "Succeed"},
			"No": {"Type": "Succeed"}
		}
	}`
	if err := validateDefinitionStructure(emptyAnd, "STANDARD"); err == nil {
		t.Fatal("an empty And array was accepted")
	}

	// A top-level rule still requires Next.
	missingNext := `{
		"StartAt": "Pick",
		"States": {
			"Pick": {
				"Type": "Choice",
				"Choices": [{"Variable": "$.a", "IsPresent": true}],
				"Default": "No"
			},
			"No": {"Type": "Succeed"}
		}
	}`
	if err := validateDefinitionStructure(missingNext, "STANDARD"); err == nil {
		t.Fatal("a top-level rule without Next was accepted")
	}
}

// TestChoiceRuleDialectAndRetrierCensus pins the dialect separation on
// choice rules ("Variables and comparison fields are only available for
// JSONPath. Condition is only available for JSONata."), the Variable
// requirement for Is* rules, and the States.ALL placement plus
// retrier/catcher member census.
func TestChoiceRuleDialectAndRetrierCensus(t *testing.T) {
	// A JSONata rule with a JSONPath comparator is rejected.
	jsonataJSONPathRule := `{
		"QueryLanguage": "JSONata", "StartAt": "Pick",
		"States": {
			"Pick": {"Type": "Choice", "Choices": [{"Variable": "$.x", "StringEquals": "a", "Next": "Y"}], "Default": "N"},
			"Y": {"Type": "Succeed"}, "N": {"Type": "Succeed"}
		}
	}`
	if err := validateDefinitionStructure(jsonataJSONPathRule, "STANDARD"); err == nil {
		t.Fatal("a JSONPath comparator rule in a JSONata Choice was accepted")
	}

	// A JSONata rule without Condition is rejected.
	jsonataNoCondition := `{
		"QueryLanguage": "JSONata", "StartAt": "Pick",
		"States": {
			"Pick": {"Type": "Choice", "Choices": [{"Next": "Y"}], "Default": "N"},
			"Y": {"Type": "Succeed"}, "N": {"Type": "Succeed"}
		}
	}`
	if err := validateDefinitionStructure(jsonataNoCondition, "STANDARD"); err == nil {
		t.Fatal("a JSONata choice rule without Condition was accepted")
	}

	// A JSONPath Is* rule without Variable is rejected.
	isWithoutVariable := `{
		"StartAt": "Pick",
		"States": {
			"Pick": {"Type": "Choice", "Choices": [{"IsNull": true, "Next": "Y"}], "Default": "N"},
			"Y": {"Type": "Succeed"}, "N": {"Type": "Succeed"}
		}
	}`
	if err := validateDefinitionStructure(isWithoutVariable, "STANDARD"); err == nil {
		t.Fatal("an Is* rule without Variable was accepted")
	}
	isWithVariable := strings.Replace(isWithoutVariable,
		`{"IsNull": true, "Next": "Y"}`, `{"Variable": "$.x", "IsNull": true, "Next": "Y"}`, 1)
	if err := validateDefinitionStructure(isWithVariable, "STANDARD"); err != nil {
		t.Fatalf("an Is* rule with Variable was rejected: %v", err)
	}

	// States.ALL must appear alone and in the last retrier/catcher.
	allNotAlone := `{
		"StartAt": "T",
		"States": {
			"T": {"Type": "Task", "Resource": "arn:aws:states:::lambda:invoke",
				"Retry": [{"ErrorEquals": ["States.ALL", "States.Timeout"]}], "End": true}
		}
	}`
	if err := validateDefinitionStructure(allNotAlone, "STANDARD"); err == nil {
		t.Fatal("States.ALL alongside another error name was accepted")
	}
	allNotLast := `{
		"StartAt": "T",
		"States": {
			"T": {"Type": "Task", "Resource": "arn:aws:states:::lambda:invoke",
				"Retry": [{"ErrorEquals": ["States.ALL"]}, {"ErrorEquals": ["States.Timeout"]}], "End": true}
		}
	}`
	if err := validateDefinitionStructure(allNotLast, "STANDARD"); err == nil {
		t.Fatal("States.ALL in a non-final retrier was accepted")
	}
	allLastValid := `{
		"StartAt": "T",
		"States": {
			"T": {"Type": "Task", "Resource": "arn:aws:states:::lambda:invoke",
				"Retry": [{"ErrorEquals": ["States.Timeout"]}, {"ErrorEquals": ["States.ALL"]}], "End": true}
		}
	}`
	if err := validateDefinitionStructure(allLastValid, "STANDARD"); err != nil {
		t.Fatalf("States.ALL in the final retrier was rejected: %v", err)
	}

	// Retrier and catcher members are a closed census.
	typoMember := `{
		"StartAt": "T",
		"States": {
			"T": {"Type": "Task", "Resource": "arn:aws:states:::lambda:invoke",
				"Retry": [{"ErrorEquals": ["States.Timeout"], "intervalSeconds": 5}], "End": true}
		}
	}`
	if err := validateDefinitionStructure(typoMember, "STANDARD"); err == nil {
		t.Fatal("an unknown retrier member was accepted")
	}
}

package sfn

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func validateDefinitionJSONataFields(definition string) error {
	var def map[string]interface{}
	if err := json.Unmarshal([]byte(definition), &def); err != nil {
		return nil
	}

	topQL, _ := def["QueryLanguage"].(string)
	if topQL == "" {
		topQL = "JSONPath"
	}

	states, ok := def["States"].(map[string]interface{})
	if !ok {
		return nil
	}

	for stateName, stateData := range states {
		stateMap, ok := stateData.(map[string]interface{})
		if !ok {
			continue
		}

		stateType, _ := stateMap["Type"].(string)
		stateQL, _ := stateMap["QueryLanguage"].(string)
		if stateQL == "" {
			stateQL = topQL
		}

		if stateQL == "JSONata" {
			jsonPathOnlyFields := getJSONPathOnlyFields(stateType, stateMap)
			if len(jsonPathOnlyFields) > 0 {
				return NewInvalidDefinitionException(fmt.Sprintf(
					"State '%s' uses JSONata QueryLanguage but contains JSONPath-only field(s): %s",
					stateName, strings.Join(jsonPathOnlyFields, ", ")))
			}
		} else {
			jsonataOnlyFields := getJSONataOnlyFields(stateType, stateMap)
			if len(jsonataOnlyFields) > 0 {
				return NewInvalidDefinitionException(fmt.Sprintf(
					"State '%s' uses JSONPath QueryLanguage but contains JSONata-only field(s): %s",
					stateName, strings.Join(jsonataOnlyFields, ", ")))
			}
		}
	}

	return nil
}

func getJSONPathOnlyFields(stateType string, stateMap map[string]interface{}) []string {
	var forbidden []string

	switch stateType {
	case "Pass", "Task", "Parallel", "Map", "Succeed":
		if _, ok := stateMap["InputPath"]; ok {
			forbidden = append(forbidden, "InputPath")
		}
		if _, ok := stateMap["OutputPath"]; ok {
			forbidden = append(forbidden, "OutputPath")
		}
		if _, ok := stateMap["Parameters"]; ok {
			forbidden = append(forbidden, "Parameters")
		}
	case "Choice":
		if _, ok := stateMap["InputPath"]; ok {
			forbidden = append(forbidden, "InputPath")
		}
	case "Wait":
		if _, ok := stateMap["InputPath"]; ok {
			forbidden = append(forbidden, "InputPath")
		}
		if _, ok := stateMap["OutputPath"]; ok {
			forbidden = append(forbidden, "OutputPath")
		}
		if _, ok := stateMap["SecondsPath"]; ok {
			forbidden = append(forbidden, "SecondsPath")
		}
		if _, ok := stateMap["TimestampPath"]; ok {
			forbidden = append(forbidden, "TimestampPath")
		}
	}

	switch stateType {
	case "Task":
		if _, ok := stateMap["ResultPath"]; ok {
			forbidden = append(forbidden, "ResultPath")
		}
		if _, ok := stateMap["ResultSelector"]; ok {
			forbidden = append(forbidden, "ResultSelector")
		}
		if _, ok := stateMap["TimeoutSecondsPath"]; ok {
			forbidden = append(forbidden, "TimeoutSecondsPath")
		}
		if _, ok := stateMap["HeartbeatSecondsPath"]; ok {
			forbidden = append(forbidden, "HeartbeatSecondsPath")
		}
	case "Pass":
		if _, ok := stateMap["ResultPath"]; ok {
			forbidden = append(forbidden, "ResultPath")
		}
		if _, ok := stateMap["ResultSelector"]; ok {
			forbidden = append(forbidden, "ResultSelector")
		}
	case "Map":
		if _, ok := stateMap["ItemsPath"]; ok {
			forbidden = append(forbidden, "ItemsPath")
		}
		if _, ok := stateMap["ResultPath"]; ok {
			forbidden = append(forbidden, "ResultPath")
		}
		if _, ok := stateMap["ResultSelector"]; ok {
			forbidden = append(forbidden, "ResultSelector")
		}
	case "Parallel":
		if _, ok := stateMap["ResultPath"]; ok {
			forbidden = append(forbidden, "ResultPath")
		}
		if _, ok := stateMap["ResultSelector"]; ok {
			forbidden = append(forbidden, "ResultSelector")
		}
	case "Fail":
		if _, ok := stateMap["CausePath"]; ok {
			forbidden = append(forbidden, "CausePath")
		}
		if _, ok := stateMap["ErrorPath"]; ok {
			forbidden = append(forbidden, "ErrorPath")
		}
	}

	return forbidden
}

func getJSONataOnlyFields(stateType string, stateMap map[string]interface{}) []string {
	var forbidden []string

	if _, ok := stateMap["Arguments"]; ok {
		if stateType == "Task" || stateType == "Parallel" {
			forbidden = append(forbidden, "Arguments")
		}
	}

	if _, ok := stateMap["Items"]; ok {
		if stateType == "Map" {
			forbidden = append(forbidden, "Items")
		}
	}

	if _, ok := stateMap["Condition"]; ok {
		if stateType == "Choice" {
			forbidden = append(forbidden, "Condition")
		}
	}

	if stateType == "Choice" {
		if choices, ok := stateMap["Choices"].([]interface{}); ok {
			for _, choice := range choices {
				if choiceMap, ok := choice.(map[string]interface{}); ok {
					if _, ok := choiceMap["Condition"]; ok {
						forbidden = append(forbidden, "Condition")
						break
					}
				}
			}
		}
	}

	return forbidden
}

var variableRefRegex = regexp.MustCompile(`\$([a-zA-Z_][a-zA-Z0-9_]*)`)

func extractVariableReferences(definition string) map[string][]string {
	var def map[string]interface{}
	if err := json.Unmarshal([]byte(definition), &def); err != nil {
		return nil
	}

	states, ok := def["States"].(map[string]interface{})
	if !ok {
		return nil
	}

	assignedVars := make(map[string]bool)
	for _, stateData := range states {
		stateMap, ok := stateData.(map[string]interface{})
		if !ok {
			continue
		}
		if assign, ok := stateMap["Assign"].(map[string]interface{}); ok {
			for name := range assign {
				clean := strings.TrimPrefix(name, "$")
				assignedVars[clean] = true
			}
		}
		if choices, ok := stateMap["Choices"].([]interface{}); ok {
			for _, choice := range choices {
				if choiceMap, ok := choice.(map[string]interface{}); ok {
					if assign, ok := choiceMap["Assign"].(map[string]interface{}); ok {
						for name := range assign {
							clean := strings.TrimPrefix(name, "$")
							assignedVars[clean] = true
						}
					}
				}
			}
		}
	}

	result := make(map[string][]string)

	for stateName, stateData := range states {
		stateMap, ok := stateData.(map[string]interface{})
		if !ok {
			continue
		}

		refs := collectVariableRefsFromState(stateMap, assignedVars)
		if len(refs) > 0 {
			result[stateName] = refs
		}
	}

	return result
}

func collectVariableRefsFromState(stateMap map[string]interface{}, assignedVars map[string]bool) []string {
	seen := make(map[string]bool)
	var refs []string

	for _, field := range []string{"Assign", "Output", "Arguments", "Items", "Condition"} {
		if field == "Condition" {
			if choices, ok := stateMap["Choices"].([]interface{}); ok {
				for _, choice := range choices {
					if choiceMap, ok := choice.(map[string]interface{}); ok {
						scanValueForVariableRefs(choiceMap["Condition"], seen, &refs, assignedVars)
					}
				}
			}
			continue
		}
		if val, ok := stateMap[field]; ok {
			scanValueForVariableRefs(val, seen, &refs, assignedVars)
		}
	}

	return refs
}

var jsonataBuiltins = map[string]bool{
	"states": true, "context": true,
	"abs": true, "count": true, "sum": true, "max": true, "min": true, "average": true,
	"string": true, "substring": true, "length": true, "uppercase": true, "lowercase": true,
	"trim": true, "pad": true, "contains": true, "split": true, "join": true,
	"match": true, "replace": true, "base64encode": true, "base64decode": true,
	"number": true, "round": true, "floor": true, "ceil": true, "sqrt": true, "power": true,
	"random": true, "boolean": true, "not": true, "exists": true, "type": true,
	"each": true, "filter": true, "flatten": true, "keys": true, "lookup": true,
	"map": true, "merge": true, "reverse": true, "sort": true, "spread": true,
	"sift": true, "distinct": true, "single": true, "tail": true, "append": true,
	"errors": true, "fromMillis": true, "toMillis": true, "millis": true,
	"now": true, "uuid": true, "parse": true, "hash": true, "partition": true, "range": true,
	"description": true, "url": true, "encodeUrlComponent": true,
	"assert": true, "error": true, "order": true,
}

func scanValueForVariableRefs(v interface{}, seen map[string]bool, refs *[]string, assignedVars map[string]bool) {
	switch val := v.(type) {
	case string:
		if strings.HasPrefix(val, "{%") && strings.HasSuffix(val, "%}") {
			expr := strings.TrimPrefix(strings.TrimSuffix(val, "%}"), "{%")
			for _, match := range variableRefRegex.FindAllStringSubmatch(expr, -1) {
				name := match[1]
				if !seen[name] && (!jsonataBuiltins[name] || assignedVars[name]) {
					seen[name] = true
					*refs = append(*refs, name)
				}
			}
		}
	case map[string]interface{}:
		for _, child := range val {
			scanValueForVariableRefs(child, seen, refs, assignedVars)
		}
	case []interface{}:
		for _, child := range val {
			scanValueForVariableRefs(child, seen, refs, assignedVars)
		}
	}
}

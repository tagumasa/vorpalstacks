// Package lambda provides AWS Lambda service operations for vorpalstacks.
package lambda

import (
	"fmt"

	"vorpalstacks/internal/common/request"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

// parseFunctionResponseTypes validates the requested response-type list.
// "A list of current response type enums applied to the event source
// mapping" — the provided list is the mapping's complete list, so callers
// replace any stored value with the parsed result.
func parseFunctionResponseTypes(raw []interface{}) ([]string, error) {
	parsed := make([]string, 0, len(raw))
	for _, frt := range raw {
		s, ok := frt.(string)
		if !ok {
			continue
		}
		if s != "ReportBatchItemFailures" {
			return nil, NewInvalidParameter("FunctionResponseTypes",
				fmt.Sprintf("FunctionResponseTypes must be ReportBatchItemFailures; got %q", s))
		}
		parsed = append(parsed, s)
	}
	return parsed, nil
}

func parseDestinationConfig(destMap map[string]interface{}) *lambdastore.DestinationConfig {
	config := &lambdastore.DestinationConfig{}
	if onSuccess, ok := destMap["OnSuccess"].(map[string]interface{}); ok {
		config.OnSuccess = &lambdastore.OnSuccess{
			Destination: request.GetStringParam(onSuccess, "Destination"),
		}
	}
	if onFailure, ok := destMap["OnFailure"].(map[string]interface{}); ok {
		config.OnFailure = &lambdastore.OnFailure{
			Destination: request.GetStringParam(onFailure, "Destination"),
		}
	}
	return config
}

func toDestinationConfig(d *lambdastore.DestinationConfig) map[string]interface{} {
	result := map[string]interface{}{}
	if d.OnSuccess != nil {
		result["OnSuccess"] = map[string]interface{}{
			"Destination": d.OnSuccess.Destination,
		}
	}
	if d.OnFailure != nil {
		result["OnFailure"] = map[string]interface{}{
			"Destination": d.OnFailure.Destination,
		}
	}
	return result
}

func parseFilterCriteria(filterMap map[string]interface{}) *lambdastore.FilterCriteria {
	criteria := &lambdastore.FilterCriteria{}
	if filters, ok := filterMap["Filters"].([]interface{}); ok {
		for _, f := range filters {
			if filterObj, ok := f.(map[string]interface{}); ok {
				filter := lambdastore.Filter{}
				if pattern, ok := filterObj["Pattern"].(string); ok {
					filter.Pattern = pattern
				}
				criteria.Filters = append(criteria.Filters, filter)
			}
		}
	}
	return criteria
}

func toFilterCriteria(f *lambdastore.FilterCriteria) map[string]interface{} {
	filters := make([]map[string]interface{}, 0, len(f.Filters))
	for _, filter := range f.Filters {
		filters = append(filters, map[string]interface{}{
			"Pattern": filter.Pattern,
		})
	}
	return map[string]interface{}{
		"Filters": filters,
	}
}

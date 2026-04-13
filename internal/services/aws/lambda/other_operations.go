// Package lambda provides AWS Lambda service operations for vorpalstacks.
package lambda

import (
	"vorpalstacks/internal/common/request"
	lambdastore "vorpalstacks/internal/store/aws/lambda"
)

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

package sfn

import (
	"encoding/json"
	"strings"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

func (e *Executor) applyInputPath(input string, inputPath string) string {
	if inputPath == "" || inputPath == "$" {
		return input
	}

	var inputData interface{}
	if err := json.Unmarshal([]byte(input), &inputData); err != nil {
		return input
	}

	dataMap, ok := inputData.(map[string]interface{})
	if !ok {
		return input
	}

	filtered, exists := getJSONPathValueRaw(dataMap, inputPath)
	if !exists {
		return "{}"
	}

	result, err := json.Marshal(filtered)
	if err != nil {
		return input
	}
	return string(result)
}

func (e *Executor) applyOutputPath(output string, outputPath string) string {
	return e.applyInputPath(output, outputPath)
}

func (e *Executor) applyParameters(input string, params *sfnstore.Parameters) string {
	if params == nil || params.Values == nil {
		return input
	}

	var inputData interface{}
	if err := json.Unmarshal([]byte(input), &inputData); err != nil {
		return input
	}

	dataMap, ok := inputData.(map[string]interface{})
	if !ok {
		dataMap = make(map[string]interface{})
	}

	result := make(map[string]interface{})
	for key, value := range params.Values {
		if strings.HasSuffix(key, ".$") {
			cleanKey := strings.TrimSuffix(key, ".$")
			if jsonPath, ok := value.(string); ok {
				if strings.HasPrefix(jsonPath, "$$.") {
					result[cleanKey] = e.getContextValue(jsonPath)
				} else if resolved, exists := getJSONPathValueRaw(dataMap, jsonPath); exists {
					result[cleanKey] = resolved
				}
			}
		} else {
			processedValue := e.processParameterValue(value, dataMap)
			result[key] = processedValue
		}
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return input
	}
	return string(resultJSON)
}

func (e *Executor) processParameterValue(value interface{}, inputData map[string]interface{}) interface{} {
	switch v := value.(type) {
	case string:
		if strings.HasSuffix(v, ".$") {
			jsonPath := strings.TrimSuffix(v, ".$")
			if strings.HasPrefix(jsonPath, "$$.") {
				return e.getContextValue(jsonPath)
			}
			if val, exists := getJSONPathValueRaw(inputData, jsonPath); exists {
				return val
			}
			return nil
		}
		return v
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, val := range v {
			if strings.HasSuffix(key, ".$") {
				jsonPath, ok := val.(string)
				if !ok {
					continue
				}
				if strings.HasPrefix(jsonPath, "$$.") {
					result[strings.TrimSuffix(key, ".$")] = e.getContextValue(jsonPath)
				} else if resolved, exists := getJSONPathValueRaw(inputData, jsonPath); exists {
					result[strings.TrimSuffix(key, ".$")] = resolved
				}
			} else {
				result[key] = e.processParameterValue(val, inputData)
			}
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = e.processParameterValue(item, inputData)
		}
		return result
	default:
		return value
	}
}

func (e *Executor) applyResultPath(input, output, resultPath string) string {
	if resultPath == "" || resultPath == "$" {
		return output
	}

	var inputData map[string]interface{}
	if err := json.Unmarshal([]byte(input), &inputData); err != nil {
		return output
	}

	var outputData interface{}
	if err := json.Unmarshal([]byte(output), &outputData); err != nil {
		return output
	}

	setNestedPath(inputData, resultPath, outputData)
	mergedJSON, err := json.Marshal(inputData)
	if err != nil {
		return output
	}
	return string(mergedJSON)
}

func setNestedPath(data map[string]interface{}, path string, value interface{}) {
	path = strings.TrimPrefix(path, "$.")
	parts := strings.Split(path, ".")
	current := data
	for i, part := range parts {
		if i == len(parts)-1 {
			current[part] = value
			return
		}
		if next, ok := current[part].(map[string]interface{}); ok {
			current = next
		} else {
			next := make(map[string]interface{})
			current[part] = next
			current = next
		}
	}
}

func (e *Executor) applyResultSelector(result string, selector *sfnstore.ResultSelector) string {
	if selector == nil || selector.Fields == nil {
		return result
	}

	var resultData interface{}
	if err := json.Unmarshal([]byte(result), &resultData); err != nil {
		return result
	}

	dataMap, ok := resultData.(map[string]interface{})
	if !ok {
		return result
	}

	output := make(map[string]interface{})
	for key, value := range selector.Fields {
		if strings.HasSuffix(key, ".$") {
			cleanKey := strings.TrimSuffix(key, ".$")
			if jsonPath, ok := value.(string); ok {
				if strings.HasPrefix(jsonPath, "$$.") {
					output[cleanKey] = e.getContextValue(jsonPath)
				} else if resolved, exists := getJSONPathValueRaw(dataMap, jsonPath); exists {
					output[cleanKey] = resolved
				}
			}
		} else {
			processedValue := e.processParameterValue(value, dataMap)
			output[key] = processedValue
		}
	}

	outputJSON, err := json.Marshal(output)
	if err != nil {
		return result
	}
	return string(outputJSON)
}

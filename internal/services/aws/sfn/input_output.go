package sfn

import (
	"context"
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

func buildVarsMap(statesVar interface{}, scope *VariableScope) map[string]interface{} {
	vars, ok := statesVar.(map[string]interface{})
	if !ok {
		vars = map[string]interface{}{}
		if statesVar != nil {
			vars["states"] = statesVar
		}
	}
	if scope != nil {
		allVars := scope.GetAll()
		for k, v := range allVars {
			if _, exists := vars[k]; !exists {
				vars[k] = v
			}
		}
	}
	return vars
}

func evaluateAssign(ctx context.Context, assign map[string]interface{}, statesVar interface{}, scope *VariableScope) (map[string]interface{}, error) {
	if len(assign) == 0 {
		return nil, nil
	}

	vars := buildVarsMap(statesVar, scope)

	evaluated := make(map[string]interface{}, len(assign))
	for name, value := range assign {
		resolved, err := ResolveTemplate(ctx, value, nil, vars)
		if err != nil {
			return nil, err
		}
		evaluated[strings.TrimPrefix(name, "$")] = resolved
	}

	return evaluated, nil
}

func (e *Executor) applyAssign(ctx context.Context, assign map[string]interface{}, statesVar interface{}, scope *VariableScope) error {
	if len(assign) == 0 {
		return nil
	}

	vars := buildVarsMap(statesVar, scope)

	evaluated := make(map[string]interface{}, len(assign))
	for name, value := range assign {
		resolved, err := ResolveTemplate(ctx, value, nil, vars)
		if err != nil {
			return err
		}
		evaluated[strings.TrimPrefix(name, "$")] = resolved
	}

	if scope != nil {
		return scope.SetAll(evaluated)
	}
	return nil
}

func (e *Executor) applyJSONataOutput(ctx context.Context, output interface{}, statesVar interface{}, scope *VariableScope) (interface{}, error) {
	if output == nil {
		return nil, nil
	}

	vars := buildVarsMap(statesVar, scope)

	return ResolveTemplate(ctx, output, nil, vars)
}

func (e *Executor) applyJSONataArguments(ctx context.Context, arguments interface{}, statesVar interface{}, scope *VariableScope) (string, error) {
	if arguments == nil {
		return "{}", nil
	}

	vars := buildVarsMap(statesVar, scope)

	result, err := ResolveTemplate(ctx, arguments, nil, vars)
	if err != nil {
		return "{}", err
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "{}", err
	}
	return string(resultJSON), nil
}

func (e *Executor) applyItemSelector(ctx context.Context, execCtx *ExecutionContext, selector interface{}, itemValue interface{}, isJSONata bool) (interface{}, error) {
	if selector == nil {
		return itemValue, nil
	}

	if isJSONata {
		statesVar := e.buildStatesVarWithContext(execCtx, itemValue, nil, nil)
		vars := buildVarsMap(statesVar, execCtx.VariableScope)
		return ResolveTemplate(ctx, selector, nil, vars)
	}

	return e.applyItemSelectorJSONPath(selector, itemValue)
}

func (e *Executor) applyItemSelectorJSONPath(selector interface{}, itemValue interface{}) (interface{}, error) {
	selectorMap, ok := selector.(map[string]interface{})
	if !ok {
		return itemValue, nil
	}

	itemMap, ok := itemValue.(map[string]interface{})
	if !ok {
		return itemValue, nil
	}

	output := make(map[string]interface{})
	for key, value := range selectorMap {
		if strings.HasSuffix(key, ".$") {
			cleanKey := strings.TrimSuffix(key, ".$")
			if jsonPath, ok := value.(string); ok {
				if strings.HasPrefix(jsonPath, "$$.") {
					output[cleanKey] = e.getContextValue(jsonPath)
				} else if resolved, exists := getJSONPathValueRaw(itemMap, jsonPath); exists {
					output[cleanKey] = resolved
				}
			}
		} else {
			output[key] = value
		}
	}
	return output, nil
}

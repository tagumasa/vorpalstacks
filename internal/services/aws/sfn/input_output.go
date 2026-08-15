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

// applyParameters applies a Parameters block to the state input. taskToken
// is the token minted for the current activity-task attempt so
// $$.Task.Token resolves to the exact token the worker must return; callers
// outside task Parameters pass an empty string, and a $$.Task.Token
// reference there fails the evaluation. Failures are classified here as
// States.Runtime so every state type surfaces the same error code.
func (e *Executor) applyParameters(taskToken string, input string, params *sfnstore.Parameters) (string, *ExecutionError) {
	if params == nil || params.Values == nil {
		return input, nil
	}

	var inputData interface{}
	if err := json.Unmarshal([]byte(input), &inputData); err != nil {
		return input, nil
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
					ctxVal, ctxErr := e.getContextValue(taskToken, jsonPath)
					if ctxErr != nil {
						return "", newJSONPathEvalError("Parameters", ctxErr)
					}
					result[cleanKey] = ctxVal
				} else if resolved, exists := getJSONPathValueRaw(dataMap, jsonPath); exists {
					result[cleanKey] = resolved
				}
			}
		} else {
			processedValue, procErr := e.processParameterValue(taskToken, value, dataMap)
			if procErr != nil {
				return "", newJSONPathEvalError("Parameters", procErr)
			}
			result[key] = processedValue
		}
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return input, nil
	}
	return string(resultJSON), nil
}

func (e *Executor) processParameterValue(taskToken string, value interface{}, inputData map[string]interface{}) (interface{}, error) {
	switch v := value.(type) {
	case string:
		if strings.HasSuffix(v, ".$") {
			jsonPath := strings.TrimSuffix(v, ".$")
			if strings.HasPrefix(jsonPath, "$$.") {
				return e.getContextValue(taskToken, jsonPath)
			}
			if val, exists := getJSONPathValueRaw(inputData, jsonPath); exists {
				return val, nil
			}
			return nil, nil
		}
		return v, nil
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, val := range v {
			if strings.HasSuffix(key, ".$") {
				jsonPath, ok := val.(string)
				if !ok {
					continue
				}
				if strings.HasPrefix(jsonPath, "$$.") {
					ctxVal, ctxErr := e.getContextValue(taskToken, jsonPath)
					if ctxErr != nil {
						return nil, ctxErr
					}
					result[strings.TrimSuffix(key, ".$")] = ctxVal
				} else if resolved, exists := getJSONPathValueRaw(inputData, jsonPath); exists {
					result[strings.TrimSuffix(key, ".$")] = resolved
				}
			} else {
				processed, procErr := e.processParameterValue(taskToken, val, inputData)
				if procErr != nil {
					return nil, procErr
				}
				result[key] = processed
			}
		}
		return result, nil
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			processed, procErr := e.processParameterValue(taskToken, item, inputData)
			if procErr != nil {
				return nil, procErr
			}
			result[i] = processed
		}
		return result, nil
	default:
		return value, nil
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

// applyResultSelector applies a ResultSelector block to a state's result.
// taskToken is the token of the task attempt that produced the result
// (activity tasks only): the context object exposes $$.Task.Token in
// ResultSelector, and for an activity task it resolves to the token the
// attempt actually ran under. States without a token pass an empty string;
// a $$.Task.Token reference there fails the evaluation, classified here as
// States.Runtime so every state type surfaces the same error code.
func (e *Executor) applyResultSelector(result string, selector *sfnstore.ResultSelector, taskToken string) (string, *ExecutionError) {
	if selector == nil || selector.Fields == nil {
		return result, nil
	}

	var resultData interface{}
	if err := json.Unmarshal([]byte(result), &resultData); err != nil {
		return result, nil
	}

	dataMap, ok := resultData.(map[string]interface{})
	if !ok {
		return result, nil
	}

	output := make(map[string]interface{})
	for key, value := range selector.Fields {
		if strings.HasSuffix(key, ".$") {
			cleanKey := strings.TrimSuffix(key, ".$")
			if jsonPath, ok := value.(string); ok {
				if strings.HasPrefix(jsonPath, "$$.") {
					ctxVal, ctxErr := e.getContextValue(taskToken, jsonPath)
					if ctxErr != nil {
						return "", newJSONPathEvalError("ResultSelector", ctxErr)
					}
					output[cleanKey] = ctxVal
				} else if resolved, exists := getJSONPathValueRaw(dataMap, jsonPath); exists {
					output[cleanKey] = resolved
				}
			}
		} else {
			processedValue, procErr := e.processParameterValue(taskToken, value, dataMap)
			if procErr != nil {
				return "", newJSONPathEvalError("ResultSelector", procErr)
			}
			output[key] = processedValue
		}
	}

	outputJSON, err := json.Marshal(output)
	if err != nil {
		return result, nil
	}
	return string(outputJSON), nil
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

// applyItemSelector evaluates a JSONata ItemSelector template against a
// single map item. JSONata failures stay classified as
// States.QueryEvaluationError by the caller; the JSONPath dialect is
// handled by applyItemSelectorJSONPath, which classifies as States.Runtime.
func (e *Executor) applyItemSelector(ctx context.Context, execCtx *ExecutionContext, selector interface{}, itemValue interface{}) (interface{}, error) {
	if selector == nil {
		return itemValue, nil
	}

	statesVar := e.buildStatesVarWithContext(execCtx, itemValue, nil, nil)
	vars := buildVarsMap(statesVar, execCtx.VariableScope)
	return ResolveTemplate(ctx, selector, nil, vars)
}

// applyItemSelectorJSONPath applies a JSONPath ItemSelector to a single map
// item. A Map ItemSelector evaluates outside any task, so no attempt token
// exists; a $$.Task.Token reference fails the evaluation, classified here
// as States.Runtime so every state type surfaces the same error code.
func (e *Executor) applyItemSelectorJSONPath(selector interface{}, itemValue interface{}) (interface{}, *ExecutionError) {
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
					ctxVal, ctxErr := e.getContextValue("", jsonPath)
					if ctxErr != nil {
						return nil, newJSONPathEvalError("ItemSelector", ctxErr)
					}
					output[cleanKey] = ctxVal
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

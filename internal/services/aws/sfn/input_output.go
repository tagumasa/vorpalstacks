package sfn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// parseVariableReference splits "$name.rest" into the variable name and the
// remaining subpath. ok is false for the context object ("$$…"), plain
// input paths ("$.…" or "$") and names whose first character is not an
// ID_Start character — the '$' is the reference sigil, never part of the
// name.
func parseVariableReference(path string) (name, rest string, ok bool) {
	if len(path) < 2 || path[0] != '$' || path[1] == '$' || path[1] == '.' {
		return "", "", false
	}
	s := path[1:]
	end := strings.IndexByte(s, '.')
	head := s
	if end >= 0 {
		head = s[:end]
	}
	if head == "" {
		return "", "", false
	}
	for i, r := range head {
		if i == 0 {
			if !isIDStart(r) {
				return "", "", false
			}
		} else if !isIDContinue(r) {
			return "", "", false
		}
	}
	if end >= 0 {
		return head, s[end:], true
	}
	return head, "", true
}

// resolveVariableRef reads a workflow variable (with an optional subpath)
// from the execution's variable scope. An undefined variable is an error,
// not a missing value: "If $x had not been previously assigned, the
// example would fail because $x would be undefined."
func (e *Executor) resolveVariableRef(execCtx *ExecutionContext, name, rest string) (interface{}, bool, error) {
	if execCtx == nil || execCtx.VariableScope == nil {
		return nil, false, fmt.Errorf("variable $%s is not defined", name)
	}
	value, defined := execCtx.VariableScope.Get(name)
	if !defined {
		return nil, false, fmt.Errorf("variable $%s is not defined", name)
	}
	if rest == "" {
		return value, true, nil
	}
	resolved, found := getJSONPathValueAny(value, "$"+rest)
	if !found {
		return nil, false, nil
	}
	return resolved, true, nil
}

// applyInputPath selects the state input. The path may address the input
// ($.), the context object ($$. — InputPath is one of the documented $$.-
// fields) or a workflow variable ($name). An unresolvable path or a path
// that cannot apply to the input shape fails with States.Runtime —
// "errors at runtime, such as attempting to apply InputPath or OutputPath
// on a null JSON payload" — never a silent substitution.
func (e *Executor) applyInputPath(execCtx *ExecutionContext, input string, inputPath string) (string, *ExecutionError) {
	if inputPath == "" || inputPath == "$" {
		return input, nil
	}
	runtimeFailure := func(cause string) *ExecutionError {
		return &ExecutionError{ErrorCode: "States.Runtime", Cause: "InputPath: " + cause}
	}
	if isContextPath(inputPath) {
		v, err := e.getContextValue(execCtx, "", inputPath)
		if err != nil {
			return "", runtimeFailure(err.Error())
		}
		marshalled, mErr := json.Marshal(v)
		if mErr != nil {
			return "", runtimeFailure(mErr.Error())
		}
		return string(marshalled), nil
	}
	if name, rest, isVar := parseVariableReference(inputPath); isVar {
		v, found, err := e.resolveVariableRef(execCtx, name, rest)
		if err != nil {
			return "", runtimeFailure(err.Error())
		}
		if !found {
			return "", runtimeFailure(fmt.Sprintf("path %s selected no value", inputPath))
		}
		marshalled, mErr := json.Marshal(v)
		if mErr != nil {
			return "", runtimeFailure(mErr.Error())
		}
		return string(marshalled), nil
	}

	var inputData interface{}
	if err := json.Unmarshal([]byte(input), &inputData); err != nil {
		return "", runtimeFailure("state input is not valid JSON")
	}

	filtered, exists := getJSONPathValueAny(inputData, inputPath)
	if !exists {
		return "", runtimeFailure(fmt.Sprintf("path %s selected no value", inputPath))
	}

	result, err := json.Marshal(filtered)
	if err != nil {
		return "", runtimeFailure(err.Error())
	}
	return string(result), nil
}

func (e *Executor) applyOutputPath(execCtx *ExecutionContext, output string, outputPath string) (string, *ExecutionError) {
	selected, err := e.applyInputPath(execCtx, output, outputPath)
	if err != nil {
		err.Cause = strings.Replace(err.Cause, "InputPath:", "OutputPath:", 1)
	}
	return selected, err
}

// resolvePayloadTemplateValue resolves the value of a ".$"-suffixed payload
// template key: an intrinsic invocation evaluates against the state input,
// the value may address the context object ($$.) or a workflow variable
// ($name), and otherwise it is an input JSONPath. found is false only when
// a plain path selects no value, which every payload builder fails with
// States.ParameterPathFailure under its field's name; evaluation failures
// classify as States.Runtime.
func (e *Executor) resolvePayloadTemplateValue(execCtx *ExecutionContext, taskToken, value string, data interface{}) (interface{}, bool, error) {
	if looksLikeIntrinsic(value) {
		resolved, err := e.evaluateIntrinsic(execCtx, taskToken, value, data, 1)
		if err != nil {
			return nil, false, err
		}
		return resolved, true, nil
	}
	if isContextPath(value) {
		ctxVal, ctxErr := e.getContextValue(execCtx, taskToken, value)
		if ctxErr != nil {
			return nil, false, ctxErr
		}
		return ctxVal, true, nil
	}
	if name, rest, isVar := parseVariableReference(value); isVar {
		return e.resolveVariableRef(execCtx, name, rest)
	}
	if resolved, exists := getJSONPathValueAny(data, value); exists {
		return resolved, true, nil
	}
	return nil, false, nil
}

// applyParameters applies a Parameters block to the state input. taskToken
// is the token minted for the current activity-task attempt so
// $$.Task.Token resolves to the exact token the worker must return; callers
// outside task Parameters pass an empty string, and a $$.Task.Token
// reference there fails the evaluation. A malformed template keeps its
// States.ParameterPathFailure identity; evaluation failures classify as
// States.Runtime so every state type surfaces the same error code.
func (e *Executor) applyParameters(execCtx *ExecutionContext, taskToken string, input string, params *sfnstore.Parameters) (string, *ExecutionError) {
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
			// A ".$"-suffixed key carries a path string by the
			// payload-template convention; any other JSON type is an
			// invalid template, not data to pass through — failing beats
			// silently dropping the entry from the payload.
			jsonPath, isString := value.(string)
			if !isString {
				return "", &ExecutionError{ErrorCode: "States.ParameterPathFailure", Cause: fmt.Sprintf("Parameters: the value of the parameter %q must be a path string", key)}
			}
			resolved, found, evalErr := e.resolvePayloadTemplateValue(execCtx, taskToken, jsonPath, dataMap)
			if evalErr != nil {
				return "", classifyPayloadTemplateError("Parameters", evalErr)
			}
			if !found {
				return "", parameterPathFailure("Parameters", jsonPath)
			}
			result[cleanKey] = resolved
		} else {
			processedValue, procErr := e.processParameterValue(execCtx, taskToken, value, dataMap)
			if procErr != nil {
				return "", classifyPayloadTemplateError("Parameters", procErr)
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

func (e *Executor) processParameterValue(execCtx *ExecutionContext, taskToken string, value interface{}, inputData interface{}) (interface{}, error) {
	switch v := value.(type) {
	case string:
		// The ".$" convention is key-side only: a literal string value that
		// happens to end in ".$" is data, not a path reference.
		return v, nil
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, val := range v {
			if strings.HasSuffix(key, ".$") {
				// Same payload-template convention one level deeper: a
				// ".$"-suffixed key with a non-string value is an invalid
				// template, not an entry to drop.
				jsonPath, isString := val.(string)
				if !isString {
					return nil, &ExecutionError{ErrorCode: "States.ParameterPathFailure", Cause: fmt.Sprintf("the value of the parameter %q must be a path string", key)}
				}
				resolved, found, evalErr := e.resolvePayloadTemplateValue(execCtx, taskToken, jsonPath, inputData)
				if evalErr != nil {
					return nil, evalErr
				}
				if !found {
					return nil, parameterPathFailure("", jsonPath)
				}
				result[strings.TrimSuffix(key, ".$")] = resolved
			} else {
				processed, procErr := e.processParameterValue(execCtx, taskToken, val, inputData)
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
			processed, procErr := e.processParameterValue(execCtx, taskToken, item, inputData)
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

// applyResultPath folds the state result into the state input at the
// reference path. A configured ResultPath that cannot apply to the input
// the state received — a non-object input has nowhere to inject — fails
// with States.ResultPathMatchFailure: "Suppose a state's input is the
// string \"foo\", and its \"ResultPath\" field has the value \"$.x\" … Then
// ResultPath cannot apply and the interpreter fails the machine".
func (e *Executor) applyResultPath(input, output, resultPath string) (string, *ExecutionError) {
	if resultPath == "" || resultPath == "$" {
		return output, nil
	}

	var inputData map[string]interface{}
	if err := json.Unmarshal([]byte(input), &inputData); err != nil {
		return "", &ExecutionError{ErrorCode: "States.ResultPathMatchFailure", Cause: fmt.Sprintf("ResultPath %s cannot apply to a non-object state input", resultPath)}
	}

	var outputData interface{}
	if err := json.Unmarshal([]byte(output), &outputData); err != nil {
		return output, nil
	}

	setNestedPath(inputData, resultPath, outputData)
	mergedJSON, err := json.Marshal(inputData)
	if err != nil {
		return output, nil
	}
	return string(mergedJSON), nil
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
// a $$.Task.Token reference there fails the evaluation. A malformed
// template keeps its States.ParameterPathFailure identity; evaluation
// failures classify as States.Runtime so every state type surfaces the
// same error code.
func (e *Executor) applyResultSelector(execCtx *ExecutionContext, result string, selector *sfnstore.ResultSelector, taskToken string) (string, *ExecutionError) {
	if selector == nil || selector.Fields == nil {
		return result, nil
	}

	var resultData interface{}
	if err := json.Unmarshal([]byte(result), &resultData); err != nil {
		return result, nil
	}

	output := make(map[string]interface{})
	for key, value := range selector.Fields {
		if strings.HasSuffix(key, ".$") {
			cleanKey := strings.TrimSuffix(key, ".$")
			// A ".$"-suffixed key carries a path string by the
			// payload-template convention; any other JSON type is an
			// invalid template, not an entry to drop silently.
			jsonPath, isString := value.(string)
			if !isString {
				return "", &ExecutionError{ErrorCode: "States.ParameterPathFailure", Cause: fmt.Sprintf("ResultSelector: the value of the parameter %q must be a path string", key)}
			}
			resolved, found, evalErr := e.resolvePayloadTemplateValue(execCtx, taskToken, jsonPath, resultData)
			if evalErr != nil {
				return "", classifyPayloadTemplateError("ResultSelector", evalErr)
			}
			if !found {
				return "", parameterPathFailure("ResultSelector", jsonPath)
			}
			output[cleanKey] = resolved
		} else {
			processedValue, procErr := e.processParameterValue(execCtx, taskToken, value, resultData)
			if procErr != nil {
				return "", classifyPayloadTemplateError("ResultSelector", procErr)
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
		// JSONata Assign keys reference the variable with its sigil
		// ("$x"); the stored name is the bare identifier, which the naming
		// rules then validate.
		bare := strings.TrimPrefix(name, "$")
		if err := ValidateVariableName(bare); err != nil {
			return nil, err
		}
		evaluated[bare] = resolved
	}

	return evaluated, nil
}

// isContextOrVariablePath reports whether a path addresses the context
// object or a workflow variable rather than the state input.
func isContextOrVariablePath(path string) bool {
	if isContextPath(path) {
		return true
	}
	_, _, ok := parseVariableReference(path)
	return ok
}

// parameterPathFailure is the payload-template failure the ASL spec names:
// "If the path is legal but cannot be applied successfully, the interpreter
// fails the machine execution" with States.ParameterPathFailure — a path
// that selects no value is never a silent key omission. field names the
// payload block the path belonged to when known; the nested template
// walker prefixes it at the classification boundary.
func parameterPathFailure(field, path string) *ExecutionError {
	cause := fmt.Sprintf("the path %q selected no value", path)
	if field != "" {
		cause = field + ": " + cause
	}
	return &ExecutionError{ErrorCode: "States.ParameterPathFailure", Cause: cause}
}

// classifyPayloadTemplateError splits a payload-template failure into its
// wire classification: a typed States.ParameterPathFailure keeps its
// identity (its cause relayed under the field's name), and every other
// failure is the JSONPath evaluation-error class States.Runtime.
func classifyPayloadTemplateError(field string, err error) *ExecutionError {
	var paramErr *ExecutionError
	if errors.As(err, &paramErr) && paramErr.ErrorCode == "States.ParameterPathFailure" {
		return &ExecutionError{ErrorCode: paramErr.ErrorCode, Cause: field + ": " + paramErr.Cause}
	}
	return newJSONPathEvalError(field, err)
}

// applyJSONPathCatchAssign evaluates a JSONPath Catch handler's Assign
// against the error output and stages it for the executor's post-state
// application; the values become visible in the next state.
func (e *Executor) applyJSONPathCatchAssign(execCtx *ExecutionContext, assign map[string]interface{}, catchOutput string) *ExecutionError {
	var root interface{}
	if err := json.Unmarshal([]byte(catchOutput), &root); err != nil {
		root = nil
	}
	evaluated, err := e.evaluateJSONPathAssign(execCtx, assign, root)
	if err != nil {
		return newJSONPathEvalError("Assign", err)
	}
	execCtx.PendingAssign = evaluated
	return nil
}

// evaluateJSONPathAssign evaluates a JSONPath Assign block as a payload
// template over the given root: keys name variables, ".$"-suffixed keys
// carry path expressions (input paths, context nodes, variables, intrinsic
// invocations), and plain keys are literal values that may nest further
// templates. The result feeds PendingAssign; new values become visible in
// the next state.
func (e *Executor) evaluateJSONPathAssign(execCtx *ExecutionContext, assign map[string]interface{}, root interface{}) (map[string]interface{}, error) {
	evaluated := make(map[string]interface{}, len(assign))
	for key, value := range assign {
		name := key
		if strings.HasSuffix(key, ".$") {
			name = strings.TrimSuffix(key, ".$")
			path, isPath := value.(string)
			if !isPath {
				return nil, fmt.Errorf("Assign %q must carry a JSONPath expression string", key)
			}
			resolved, found, err := e.resolvePayloadTemplateValue(execCtx, "", path, root)
			if err != nil {
				return nil, err
			}
			if !found {
				return nil, fmt.Errorf("Assign path %q selected no value", path)
			}
			if err := ValidateVariableName(name); err != nil {
				return nil, err
			}
			evaluated[name] = resolved
			continue
		}
		if err := ValidateVariableName(name); err != nil {
			return nil, err
		}
		processed, err := e.processParameterValue(execCtx, "", value, root)
		if err != nil {
			return nil, err
		}
		evaluated[name] = processed
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
// exists; a $$.Task.Token reference fails the evaluation. A malformed
// template keeps its States.ParameterPathFailure identity; evaluation
// failures classify as States.Runtime so every state type surfaces the
// same error code.
func (e *Executor) applyItemSelectorJSONPath(execCtx *ExecutionContext, selector interface{}, itemValue interface{}) (interface{}, *ExecutionError) {
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
			// A ".$"-suffixed key carries a path string by the
			// payload-template convention; any other JSON type is an
			// invalid template, not an entry to drop silently.
			jsonPath, isString := value.(string)
			if !isString {
				return nil, &ExecutionError{ErrorCode: "States.ParameterPathFailure", Cause: fmt.Sprintf("ItemSelector: the value of the parameter %q must be a path string", key)}
			}
			resolved, found, evalErr := e.resolvePayloadTemplateValue(execCtx, "", jsonPath, itemMap)
			if evalErr != nil {
				return nil, classifyPayloadTemplateError("ItemSelector", evalErr)
			}
			if !found {
				return nil, parameterPathFailure("ItemSelector", jsonPath)
			}
			output[cleanKey] = resolved
		} else {
			// Non-path keys still resolve nested payload templates, exactly
			// as Parameters does.
			processed, procErr := e.processParameterValue(execCtx, "", value, itemMap)
			if procErr != nil {
				return nil, classifyPayloadTemplateError("ItemSelector", procErr)
			}
			output[key] = processed
		}
	}
	return output, nil
}

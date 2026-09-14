package sfn

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

func (e *Executor) executeChoice(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.ChoiceState) (string, string, error) {
	isJSONata := IsJSONataState(state, execCtx.QueryLanguage)

	eventId := execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:             execCtx.Execution.ExecutionArn,
		EventId:                  eventId,
		PreviousEventId:          eventId - 1,
		Type:                     "ChoiceStateEntered",
		Timestamp:                time.Now().UTC(),
		StateEnteredEventDetails: stateEnteredDetails(execCtx, execCtx.Input),
	})

	if isJSONata {
		return e.executeChoiceJSONata(ctx, execCtx, state)
	}

	// InputPath selects the effective input the rules evaluate against;
	// OutputPath filters the value handed to the transitioned state.
	processedInput, ipErr := e.applyInputPath(execCtx, execCtx.Input, state.GetInputPath())
	if ipErr != nil {
		return "", "", ipErr
	}

	var inputData map[string]interface{}
	if err := json.Unmarshal([]byte(processedInput), &inputData); err != nil {
		return "", "", fmt.Errorf("failed to parse input JSON: %w", err)
	}

	for _, rule := range state.Choices {
		matched, ruleErr := e.evaluateChoiceRule(execCtx, rule, inputData)
		if ruleErr != nil {
			return "", "", ruleErr
		}
		if matched {
			if len(rule.Assign) > 0 {
				// JSONPath Choice-rule Assign evaluates against the state
				// input; the values become visible in the next state.
				evaluated, err := e.evaluateJSONPathAssign(execCtx, rule.Assign, inputData)
				if err != nil {
					return "", "", newJSONPathEvalError("Assign", err)
				}
				execCtx.PendingAssign = evaluated
			}
			output, opErr := e.applyOutputPath(execCtx, processedInput, state.GetOutputPath())
			if opErr != nil {
				return "", "", opErr
			}
			eventId = execCtx.nextEventId()
			e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
				ExecutionArn:            execCtx.Execution.ExecutionArn,
				EventId:                 eventId,
				PreviousEventId:         eventId - 1,
				Type:                    "ChoiceStateExited",
				Timestamp:               time.Now().UTC(),
				StateExitedEventDetails: choiceExitedDetails(execCtx, output, rule.Next),
			})
			return output, rule.Next, nil
		}
	}

	nextState := state.Default
	if nextState == "" {
		return "", "", &ExecutionError{ErrorCode: "States.NoChoiceMatched", Cause: "no choice rule matched and no default state specified"}
	}

	// "Step Functions evaluates the top-level Assign only when no Choice
	// Rule matches and the workflow transitions to the Default state" —
	// a matched rule's Assign takes over instead.
	if len(state.Assign) > 0 {
		evaluated, err := e.evaluateJSONPathAssign(execCtx, state.Assign, inputData)
		if err != nil {
			return "", "", newJSONPathEvalError("Assign", err)
		}
		execCtx.PendingAssign = evaluated
	}

	output, opErr := e.applyOutputPath(execCtx, processedInput, state.GetOutputPath())
	if opErr != nil {
		return "", "", opErr
	}

	eventId = execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:            execCtx.Execution.ExecutionArn,
		EventId:                 eventId,
		PreviousEventId:         eventId - 1,
		Type:                    "ChoiceStateExited",
		Timestamp:               time.Now().UTC(),
		StateExitedEventDetails: choiceExitedDetails(execCtx, output, nextState),
	})

	return output, nextState, nil
}

// choiceExitedDetails builds the generic exit details for a Choice state,
// recording the transition the choice took so the resume machinery can
// recover it from history; the response serialiser never emits the
// internal NextState member.
func choiceExitedDetails(execCtx *ExecutionContext, output, nextState string) *sfnstore.StateExitedEventDetails {
	details := stateExitedDetails(execCtx, output)
	details.NextState = nextState
	return details
}

func (e *Executor) executeChoiceJSONata(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.ChoiceState) (string, string, error) {
	var inputData interface{}
	if execCtx.Input != "" {
		if err := json.Unmarshal([]byte(execCtx.Input), &inputData); err != nil {
			return "", "", e.newQueryEvalError(ctx, execCtx, "Input", "failed to parse input JSON")
		}
	}
	statesVar := e.buildStatesVarWithContext(execCtx, inputData, nil, nil)

	for _, rule := range state.Choices {
		if rule.Condition != "" {
			vars := buildVarsMap(statesVar, execCtx.VariableScope)
			result, err := EvaluateJSONata(ctx, UnwrapExpression(rule.Condition), nil, vars)
			if err != nil {
				return "", "", e.newQueryEvalError(ctx, execCtx, "Condition", err.Error())
			}
			matched, ok := result.(bool)
			if !ok {
				return "", "", e.newQueryEvalError(ctx, execCtx, "Condition", fmt.Sprintf("the condition expression must evaluate to a boolean, got %v", result))
			}
			if !matched {
				continue
			}

			if len(rule.Assign) > 0 {
				evaluated, err := evaluateAssign(ctx, rule.Assign, statesVar, execCtx.VariableScope)
				if err != nil {
					return "", "", e.newQueryEvalError(ctx, execCtx, "Assign", err.Error())
				}
				execCtx.PendingAssign = evaluated
			}

			eventId := execCtx.nextEventId()
			e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
				ExecutionArn:            execCtx.Execution.ExecutionArn,
				EventId:                 eventId,
				PreviousEventId:         eventId - 1,
				Type:                    "ChoiceStateExited",
				Timestamp:               time.Now().UTC(),
				StateExitedEventDetails: choiceExitedDetails(execCtx, execCtx.Input, rule.Next),
			})
			return execCtx.Input, rule.Next, nil
		}
	}

	nextState := state.Default
	if nextState == "" {
		return "", "", &ExecutionError{ErrorCode: "States.NoChoiceMatched", Cause: "no choice rule matched and no default state specified"}
	}

	if len(state.Assign) > 0 {
		evaluated, err := evaluateAssign(ctx, state.Assign, statesVar, execCtx.VariableScope)
		if err != nil {
			return "", "", e.newQueryEvalError(ctx, execCtx, "Assign", err.Error())
		}
		execCtx.PendingAssign = evaluated
	}

	eventId := execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:            execCtx.Execution.ExecutionArn,
		EventId:                 eventId,
		PreviousEventId:         eventId - 1,
		Type:                    "ChoiceStateExited",
		Timestamp:               time.Now().UTC(),
		StateExitedEventDetails: choiceExitedDetails(execCtx, execCtx.Input, nextState),
	})

	return execCtx.Input, nextState, nil
}

// choiceRuleValue resolves a Choice rule operand across the three
// addressable sources: the state input, the context object ("JSONPath
// states can refer to the context ($$.) from … Variable (in Choice
// states)") and the workflow variables ("You can reference a variable in
// any field that accepts a JSONpath expression"). found distinguishes an
// absent plain path (the rule simply does not match) from a resolution
// error — an undefined variable reference or an unresolvable context node
// fails the execution, never a silent non-match.
func (e *Executor) choiceRuleValue(execCtx *ExecutionContext, inputData map[string]interface{}, path string) (interface{}, bool, error) {
	if isContextPath(path) {
		v, err := e.getContextValue(execCtx, "", path)
		if err != nil {
			return nil, false, err
		}
		return v, true, nil
	}
	if name, rest, isVar := parseVariableReference(path); isVar {
		v, found, err := e.resolveVariableRef(execCtx, name, rest)
		if err != nil {
			return nil, false, err
		}
		return v, found, nil
	}
	v, found := getJSONPathValueRaw(inputData, path)
	return v, found, nil
}

func (e *Executor) evaluateChoiceRule(execCtx *ExecutionContext, rule *sfnstore.ChoiceRule, inputData map[string]interface{}) (bool, *ExecutionError) {
	if len(rule.And) > 0 {
		for _, r := range rule.And {
			matched, err := e.evaluateChoiceRule(execCtx, r, inputData)
			if err != nil {
				return false, err
			}
			if !matched {
				return false, nil
			}
		}
		return true, nil
	}

	if len(rule.Or) > 0 {
		for _, r := range rule.Or {
			matched, err := e.evaluateChoiceRule(execCtx, r, inputData)
			if err != nil {
				return false, err
			}
			if matched {
				return true, nil
			}
		}
		return false, nil
	}

	if rule.Not != nil {
		matched, err := e.evaluateChoiceRule(execCtx, rule.Not, inputData)
		if err != nil {
			return false, err
		}
		return !matched, nil
	}

	// IsPresent observes presence itself: a reference that fails to
	// resolve is reported as absent rather than failing the execution —
	// the operator exists precisely to ask whether a node is there.
	if rule.IsPresent != nil {
		_, exists, perr := e.choiceRuleValue(execCtx, inputData, rule.Variable)
		return *rule.IsPresent == (perr == nil && exists), nil
	}

	varValue, exists, verr := e.choiceRuleValue(execCtx, inputData, rule.Variable)
	if verr != nil {
		return false, &ExecutionError{ErrorCode: "States.Runtime", Cause: "Choice Variable: " + verr.Error()}
	}
	if !exists {
		return false, nil
	}

	// --- String comparisons ---
	// "For each of these operators, the corresponding value must be of the
	// appropriate type: string, number, Boolean, or timestamp. Step
	// Functions doesn't attempt to match a numeric field to a string
	// value." Only a JSON string operand can match a string operator — a
	// numeric or Boolean variable never stringifies — while "because
	// timestamp fields are logically strings, it's possible that a field
	// considered to be a timestamp can be matched by a StringEquals
	// comparator".
	strVal, strValOK := varValue.(string)
	if strValOK {
		if rule.StringEquals != nil {
			return strVal == *rule.StringEquals, nil
		}

		if rule.StringLessThan != nil {
			return strVal < *rule.StringLessThan, nil
		}

		if rule.StringGreaterThan != nil {
			return strVal > *rule.StringGreaterThan, nil
		}

		if rule.StringLessThanEquals != nil {
			return strVal <= *rule.StringLessThanEquals, nil
		}

		if rule.StringGreaterThanEquals != nil {
			return strVal >= *rule.StringGreaterThanEquals, nil
		}

		if rule.StringMatches != nil {
			return globMatch(strVal, *rule.StringMatches), nil
		}
	} else {
		anyStringOperator := rule.StringEquals != nil || rule.StringLessThan != nil ||
			rule.StringGreaterThan != nil || rule.StringLessThanEquals != nil ||
			rule.StringGreaterThanEquals != nil || rule.StringMatches != nil
		if anyStringOperator {
			return false, nil
		}
	}

	// --- Numeric comparisons ---

	if rule.NumericEquals != nil {
		if numVal, ok := toFloat64(varValue); ok {
			return numVal == *rule.NumericEquals, nil
		}
	}

	if rule.NumericLessThan != nil {
		if numVal, ok := toFloat64(varValue); ok {
			return numVal < *rule.NumericLessThan, nil
		}
	}

	if rule.NumericGreaterThan != nil {
		if numVal, ok := toFloat64(varValue); ok {
			return numVal > *rule.NumericGreaterThan, nil
		}
	}

	if rule.NumericLessThanEquals != nil {
		if numVal, ok := toFloat64(varValue); ok {
			return numVal <= *rule.NumericLessThanEquals, nil
		}
	}

	if rule.NumericGreaterThanEquals != nil {
		if numVal, ok := toFloat64(varValue); ok {
			return numVal >= *rule.NumericGreaterThanEquals, nil
		}
	}

	// --- Boolean comparisons ---

	if rule.BooleanEquals != nil {
		if boolVal, ok := varValue.(bool); ok {
			return boolVal == *rule.BooleanEquals, nil
		}
	}

	// --- Timestamp comparisons ---

	if rule.TimestampEquals != nil {
		return e.compareTimestamp(varValue, *rule.TimestampEquals, "equals"), nil
	}

	if rule.TimestampLessThan != nil {
		return e.compareTimestamp(varValue, *rule.TimestampLessThan, "less"), nil
	}

	if rule.TimestampGreaterThan != nil {
		return e.compareTimestamp(varValue, *rule.TimestampGreaterThan, "greater"), nil
	}

	if rule.TimestampLessThanEquals != nil {
		return e.compareTimestamp(varValue, *rule.TimestampLessThanEquals, "lessEquals"), nil
	}

	if rule.TimestampGreaterThanEquals != nil {
		return e.compareTimestamp(varValue, *rule.TimestampGreaterThanEquals, "greaterEquals"), nil
	}

	// --- Type-test comparisons ---

	if rule.IsNull != nil {
		return *rule.IsNull == (varValue == nil), nil
	}

	if rule.IsBoolean != nil {
		_, isBool := varValue.(bool)
		return *rule.IsBoolean == isBool, nil
	}

	if rule.IsString != nil {
		_, isStr := varValue.(string)
		return *rule.IsString == isStr, nil
	}

	if rule.IsNumeric != nil {
		_, isNum := toFloat64(varValue)
		return *rule.IsNumeric == isNum, nil
	}

	if rule.IsTimestamp != nil {
		if s, ok := varValue.(string); ok {
			_, err := time.Parse(time.RFC3339, s)
			return *rule.IsTimestamp == (err == nil), nil
		}
		return !*rule.IsTimestamp, nil
	}

	// --- *Path comparisons (compare against a dynamic value from input) ---

	if rule.StringEqualsPath != "" {
		return e.comparePathValue(execCtx, inputData, rule.StringEqualsPath, varValue, "stringEquals")
	}
	if rule.StringLessThanPath != "" {
		return e.comparePathValue(execCtx, inputData, rule.StringLessThanPath, varValue, "stringLessThan")
	}
	if rule.StringGreaterThanPath != "" {
		return e.comparePathValue(execCtx, inputData, rule.StringGreaterThanPath, varValue, "stringGreaterThan")
	}
	if rule.StringLessThanEqualsPath != "" {
		return e.comparePathValue(execCtx, inputData, rule.StringLessThanEqualsPath, varValue, "stringLessThanEquals")
	}
	if rule.StringGreaterThanEqualsPath != "" {
		return e.comparePathValue(execCtx, inputData, rule.StringGreaterThanEqualsPath, varValue, "stringGreaterThanEquals")
	}
	if rule.NumericEqualsPath != "" {
		return e.comparePathValue(execCtx, inputData, rule.NumericEqualsPath, varValue, "numericEquals")
	}
	if rule.NumericLessThanPath != "" {
		return e.comparePathValue(execCtx, inputData, rule.NumericLessThanPath, varValue, "numericLessThan")
	}
	if rule.NumericGreaterThanPath != "" {
		return e.comparePathValue(execCtx, inputData, rule.NumericGreaterThanPath, varValue, "numericGreaterThan")
	}
	if rule.NumericLessThanEqualsPath != "" {
		return e.comparePathValue(execCtx, inputData, rule.NumericLessThanEqualsPath, varValue, "numericLessThanEquals")
	}
	if rule.NumericGreaterThanEqualsPath != "" {
		return e.comparePathValue(execCtx, inputData, rule.NumericGreaterThanEqualsPath, varValue, "numericGreaterThanEquals")
	}
	if rule.BooleanEqualsPath != "" {
		return e.comparePathValue(execCtx, inputData, rule.BooleanEqualsPath, varValue, "booleanEquals")
	}
	if rule.TimestampEqualsPath != "" {
		return e.comparePathValue(execCtx, inputData, rule.TimestampEqualsPath, varValue, "timestampEquals")
	}
	if rule.TimestampLessThanPath != "" {
		return e.comparePathValue(execCtx, inputData, rule.TimestampLessThanPath, varValue, "timestampLessThan")
	}
	if rule.TimestampGreaterThanPath != "" {
		return e.comparePathValue(execCtx, inputData, rule.TimestampGreaterThanPath, varValue, "timestampGreaterThan")
	}
	if rule.TimestampLessThanEqualsPath != "" {
		return e.comparePathValue(execCtx, inputData, rule.TimestampLessThanEqualsPath, varValue, "timestampLessThanEquals")
	}
	if rule.TimestampGreaterThanEqualsPath != "" {
		return e.comparePathValue(execCtx, inputData, rule.TimestampGreaterThanEqualsPath, varValue, "timestampGreaterThanEquals")
	}

	return false, nil
}

func (e *Executor) compareTimestamp(value interface{}, expected string, op string) bool {
	var valTime time.Time
	var err error

	switch v := value.(type) {
	case string:
		valTime, err = time.Parse(time.RFC3339, v)
		if err != nil {
			return false
		}
	default:
		return false
	}

	expectedTime, err := time.Parse(time.RFC3339, expected)
	if err != nil {
		return false
	}

	switch op {
	case "equals":
		return valTime.Equal(expectedTime)
	case "less":
		return valTime.Before(expectedTime)
	case "greater":
		return valTime.After(expectedTime)
	case "lessEquals":
		return valTime.Before(expectedTime) || valTime.Equal(expectedTime)
	case "greaterEquals":
		return valTime.After(expectedTime) || valTime.Equal(expectedTime)
	}
	return false
}

// comparePathValue reads the comparison value for a *Path operator — an
// input path, a context-object node or a workflow variable — and compares
// it against varValue according to op. String operators require both
// operands to be JSON strings (the same type discipline as the constant
// operators); a resolution error on the comparator's own path fails the
// execution.
func (e *Executor) comparePathValue(execCtx *ExecutionContext, inputData map[string]interface{}, path string, varValue interface{}, op string) (bool, *ExecutionError) {
	compareVal, exists, rerr := e.choiceRuleValue(execCtx, inputData, path)
	if rerr != nil {
		return false, &ExecutionError{ErrorCode: "States.Runtime", Cause: "Choice comparison path: " + rerr.Error()}
	}
	if !exists {
		return false, nil
	}
	if op == "stringEquals" || op == "stringLessThan" || op == "stringGreaterThan" ||
		op == "stringLessThanEquals" || op == "stringGreaterThanEquals" {
		varStr, ok1 := varValue.(string)
		cmpStr, ok2 := compareVal.(string)
		if !ok1 || !ok2 {
			return false, nil
		}
		switch op {
		case "stringEquals":
			return varStr == cmpStr, nil
		case "stringLessThan":
			return varStr < cmpStr, nil
		case "stringGreaterThan":
			return varStr > cmpStr, nil
		case "stringLessThanEquals":
			return varStr <= cmpStr, nil
		case "stringGreaterThanEquals":
			return varStr >= cmpStr, nil
		}
		return false, nil
	}
	switch op {
	case "numericEquals":
		a, ok1 := toFloat64(varValue)
		b, ok2 := toFloat64(compareVal)
		return ok1 && ok2 && a == b, nil
	case "numericLessThan":
		a, ok1 := toFloat64(varValue)
		b, ok2 := toFloat64(compareVal)
		return ok1 && ok2 && a < b, nil
	case "numericGreaterThan":
		a, ok1 := toFloat64(varValue)
		b, ok2 := toFloat64(compareVal)
		return ok1 && ok2 && a > b, nil
	case "numericLessThanEquals":
		a, ok1 := toFloat64(varValue)
		b, ok2 := toFloat64(compareVal)
		return ok1 && ok2 && a <= b, nil
	case "numericGreaterThanEquals":
		a, ok1 := toFloat64(varValue)
		b, ok2 := toFloat64(compareVal)
		return ok1 && ok2 && a >= b, nil
	case "booleanEquals":
		va, ok1 := varValue.(bool)
		vb, ok2 := compareVal.(bool)
		return ok1 && ok2 && va == vb, nil
	case "timestampEquals":
		return e.compareTimestamp(varValue, fmt.Sprintf("%v", compareVal), "equals"), nil
	case "timestampLessThan":
		return e.compareTimestamp(varValue, fmt.Sprintf("%v", compareVal), "less"), nil
	case "timestampGreaterThan":
		return e.compareTimestamp(varValue, fmt.Sprintf("%v", compareVal), "greater"), nil
	case "timestampLessThanEquals":
		return e.compareTimestamp(varValue, fmt.Sprintf("%v", compareVal), "lessEquals"), nil
	case "timestampGreaterThanEquals":
		return e.compareTimestamp(varValue, fmt.Sprintf("%v", compareVal), "greaterEquals"), nil
	}
	return false, nil
}

// globMatch matches a value against an AWS ASL StringMatches pattern.
// AWS supports `*` (zero or more characters) and `\*` / `\?` for literals.
func globMatch(value, pattern string) bool {
	// Convert the AWS glob pattern to a Go regexp.
	var sb strings.Builder
	sb.WriteString("^")
	for i := 0; i < len(pattern); i++ {
		c := pattern[i]
		if c == '\\' && i+1 < len(pattern) {
			next := pattern[i+1]
			sb.WriteString(regexp.QuoteMeta(string(next)))
			i++
		} else if c == '*' {
			sb.WriteString(".*")
		} else {
			sb.WriteString(regexp.QuoteMeta(string(c)))
		}
	}
	sb.WriteString("$")
	re, err := regexp.Compile(sb.String())
	if err != nil {
		return false
	}
	return re.MatchString(value)
}

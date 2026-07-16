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

func (e *Executor) executeChoice(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.ChoiceState) (string, error) {
	isJSONata := IsJSONataState(state, execCtx.QueryLanguage)

	eventId := execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "ChoiceStateEntered",
		Timestamp:       time.Now().UTC(),
		ChoiceStateEnteredEventDetails: &sfnstore.ChoiceStateEnteredEventDetails{
			Input: execCtx.Input,
			Name:  execCtx.CurrentState,
		},
	})

	if isJSONata {
		return e.executeChoiceJSONata(ctx, execCtx, state)
	}

	var inputData map[string]interface{}
	if err := json.Unmarshal([]byte(execCtx.Input), &inputData); err != nil {
		return "", fmt.Errorf("failed to parse input JSON: %w", err)
	}

	for _, rule := range state.Choices {
		if e.evaluateChoiceRule(rule, inputData) {
			eventId = execCtx.nextEventId()
			e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
				ExecutionArn:    execCtx.Execution.ExecutionArn,
				EventId:         eventId,
				PreviousEventId: eventId - 1,
				Type:            "ChoiceStateExited",
				Timestamp:       time.Now().UTC(),
				ChoiceStateExitedEventDetails: &sfnstore.ChoiceStateExitedEventDetails{
					Output:    execCtx.Input,
					Name:      execCtx.CurrentState,
					NextState: rule.Next,
				},
			})
			return rule.Next, nil
		}
	}

	nextState := state.Default
	if nextState == "" {
		return "", fmt.Errorf("no choice rule matched and no default state specified")
	}

	eventId = execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "ChoiceStateExited",
		Timestamp:       time.Now().UTC(),
		ChoiceStateExitedEventDetails: &sfnstore.ChoiceStateExitedEventDetails{
			Output:    execCtx.Input,
			Name:      execCtx.CurrentState,
			NextState: nextState,
		},
	})

	return nextState, nil
}

func (e *Executor) executeChoiceJSONata(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.ChoiceState) (string, error) {
	var inputData interface{}
	if execCtx.Input != "" {
		if err := json.Unmarshal([]byte(execCtx.Input), &inputData); err != nil {
			return "", e.newQueryEvalError(ctx, execCtx, "Input", "failed to parse input JSON")
		}
	}
	statesVar := e.buildStatesVarWithContext(execCtx, inputData, nil, nil)

	for _, rule := range state.Choices {
		if rule.Condition != "" {
			vars := buildVarsMap(statesVar, execCtx.VariableScope)
			result, err := EvaluateJSONata(ctx, UnwrapExpression(rule.Condition), nil, vars)
			if err != nil {
				return "", e.newQueryEvalError(ctx, execCtx, "Condition", err.Error())
			}
			matched, ok := result.(bool)
			if !ok || !matched {
				continue
			}

			if len(rule.Assign) > 0 {
				evaluated, err := evaluateAssign(ctx, rule.Assign, statesVar, execCtx.VariableScope)
				if err != nil {
					return "", e.newQueryEvalError(ctx, execCtx, "Assign", err.Error())
				}
				execCtx.PendingAssign = evaluated
			}

			eventId := execCtx.nextEventId()
			e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
				ExecutionArn:    execCtx.Execution.ExecutionArn,
				EventId:         eventId,
				PreviousEventId: eventId - 1,
				Type:            "ChoiceStateExited",
				Timestamp:       time.Now().UTC(),
				ChoiceStateExitedEventDetails: &sfnstore.ChoiceStateExitedEventDetails{
					Output:    execCtx.Input,
					Name:      execCtx.CurrentState,
					NextState: rule.Next,
				},
			})
			return rule.Next, nil
		}
	}

	nextState := state.Default
	if nextState == "" {
		return "", fmt.Errorf("no choice rule matched and no default state specified")
	}

	eventId := execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "ChoiceStateExited",
		Timestamp:       time.Now().UTC(),
		ChoiceStateExitedEventDetails: &sfnstore.ChoiceStateExitedEventDetails{
			Output:    execCtx.Input,
			Name:      execCtx.CurrentState,
			NextState: nextState,
		},
	})

	return nextState, nil
}

func (e *Executor) evaluateChoiceRule(rule *sfnstore.ChoiceRule, inputData map[string]interface{}) bool {
	if len(rule.And) > 0 {
		for _, r := range rule.And {
			if !e.evaluateChoiceRule(r, inputData) {
				return false
			}
		}
		return true
	}

	if len(rule.Or) > 0 {
		for _, r := range rule.Or {
			if e.evaluateChoiceRule(r, inputData) {
				return true
			}
		}
		return false
	}

	if rule.Not != nil {
		return !e.evaluateChoiceRule(rule.Not, inputData)
	}

	if rule.IsPresent != nil {
		_, exists := getJSONPathValueRaw(inputData, rule.Variable)
		return *rule.IsPresent == exists
	}

	varValue, exists := getJSONPathValueRaw(inputData, rule.Variable)
	if !exists {
		return false
	}

	// --- String comparisons ---

	if rule.StringEquals != nil {
		if expected, ok := rule.StringEquals[rule.Variable]; ok {
			return fmt.Sprintf("%v", varValue) == expected
		}
	}

	if rule.StringLessThan != nil {
		if expected, ok := rule.StringLessThan[rule.Variable]; ok {
			return fmt.Sprintf("%v", varValue) < expected
		}
	}

	if rule.StringGreaterThan != nil {
		if expected, ok := rule.StringGreaterThan[rule.Variable]; ok {
			return fmt.Sprintf("%v", varValue) > expected
		}
	}

	if rule.StringLessThanEquals != "" {
		return fmt.Sprintf("%v", varValue) <= rule.StringLessThanEquals
	}

	if rule.StringGreaterThanEquals != "" {
		return fmt.Sprintf("%v", varValue) >= rule.StringGreaterThanEquals
	}

	if rule.StringMatches != "" {
		return globMatch(fmt.Sprintf("%v", varValue), rule.StringMatches)
	}

	// --- Numeric comparisons ---

	if rule.NumericEquals != nil {
		if expected, ok := rule.NumericEquals[rule.Variable]; ok {
			if numVal, ok := toFloat64(varValue); ok {
				return numVal == expected
			}
		}
	}

	if rule.NumericLessThan != nil {
		if expected, ok := rule.NumericLessThan[rule.Variable]; ok {
			if numVal, ok := toFloat64(varValue); ok {
				return numVal < expected
			}
		}
	}

	if rule.NumericGreaterThan != nil {
		if expected, ok := rule.NumericGreaterThan[rule.Variable]; ok {
			if numVal, ok := toFloat64(varValue); ok {
				return numVal > expected
			}
		}
	}

	if rule.NumericLessThanEquals != nil {
		if numVal, ok := toFloat64(varValue); ok {
			return numVal <= *rule.NumericLessThanEquals
		}
	}

	if rule.NumericGreaterThanEquals != nil {
		if numVal, ok := toFloat64(varValue); ok {
			return numVal >= *rule.NumericGreaterThanEquals
		}
	}

	// --- Boolean comparisons ---

	if rule.BooleanEquals != nil {
		if expected, ok := rule.BooleanEquals[rule.Variable]; ok {
			if boolVal, ok := varValue.(bool); ok {
				return boolVal == expected
			}
		}
	}

	// --- Timestamp comparisons ---

	if rule.TimestampEquals != nil {
		if expected, ok := rule.TimestampEquals[rule.Variable]; ok {
			return e.compareTimestamp(varValue, expected, "equals")
		}
	}

	if rule.TimestampLessThan != nil {
		if expected, ok := rule.TimestampLessThan[rule.Variable]; ok {
			return e.compareTimestamp(varValue, expected, "less")
		}
	}

	if rule.TimestampGreaterThan != nil {
		if expected, ok := rule.TimestampGreaterThan[rule.Variable]; ok {
			return e.compareTimestamp(varValue, expected, "greater")
		}
	}

	if rule.TimestampLessThanEquals != "" {
		return e.compareTimestamp(varValue, rule.TimestampLessThanEquals, "lessEquals")
	}

	if rule.TimestampGreaterThanEquals != "" {
		return e.compareTimestamp(varValue, rule.TimestampGreaterThanEquals, "greaterEquals")
	}

	// --- Type-test comparisons ---

	if rule.IsNull != nil {
		return *rule.IsNull == (varValue == nil)
	}

	if rule.IsBoolean != nil {
		_, isBool := varValue.(bool)
		return *rule.IsBoolean == isBool
	}

	if rule.IsString != nil {
		_, isStr := varValue.(string)
		return *rule.IsString == isStr
	}

	if rule.IsNumeric != nil {
		_, isNum := toFloat64(varValue)
		return *rule.IsNumeric == isNum
	}

	if rule.IsTimestamp != nil {
		if s, ok := varValue.(string); ok {
			_, err := time.Parse(time.RFC3339, s)
			return *rule.IsTimestamp == (err == nil)
		}
		return !*rule.IsTimestamp
	}

	// --- *Path comparisons (compare against a dynamic value from input) ---

	if rule.StringEqualsPath != "" {
		return e.comparePathValue(inputData, rule.StringEqualsPath, varValue, "stringEquals")
	}
	if rule.StringLessThanPath != "" {
		return e.comparePathValue(inputData, rule.StringLessThanPath, varValue, "stringLessThan")
	}
	if rule.StringGreaterThanPath != "" {
		return e.comparePathValue(inputData, rule.StringGreaterThanPath, varValue, "stringGreaterThan")
	}
	if rule.StringLessThanEqualsPath != "" {
		return e.comparePathValue(inputData, rule.StringLessThanEqualsPath, varValue, "stringLessThanEquals")
	}
	if rule.StringGreaterThanEqualsPath != "" {
		return e.comparePathValue(inputData, rule.StringGreaterThanEqualsPath, varValue, "stringGreaterThanEquals")
	}
	if rule.NumericEqualsPath != "" {
		return e.comparePathValue(inputData, rule.NumericEqualsPath, varValue, "numericEquals")
	}
	if rule.NumericLessThanPath != "" {
		return e.comparePathValue(inputData, rule.NumericLessThanPath, varValue, "numericLessThan")
	}
	if rule.NumericGreaterThanPath != "" {
		return e.comparePathValue(inputData, rule.NumericGreaterThanPath, varValue, "numericGreaterThan")
	}
	if rule.NumericLessThanEqualsPath != "" {
		return e.comparePathValue(inputData, rule.NumericLessThanEqualsPath, varValue, "numericLessThanEquals")
	}
	if rule.NumericGreaterThanEqualsPath != "" {
		return e.comparePathValue(inputData, rule.NumericGreaterThanEqualsPath, varValue, "numericGreaterThanEquals")
	}
	if rule.BooleanEqualsPath != "" {
		return e.comparePathValue(inputData, rule.BooleanEqualsPath, varValue, "booleanEquals")
	}
	if rule.TimestampEqualsPath != "" {
		return e.comparePathValue(inputData, rule.TimestampEqualsPath, varValue, "timestampEquals")
	}
	if rule.TimestampLessThanPath != "" {
		return e.comparePathValue(inputData, rule.TimestampLessThanPath, varValue, "timestampLessThan")
	}
	if rule.TimestampGreaterThanPath != "" {
		return e.comparePathValue(inputData, rule.TimestampGreaterThanPath, varValue, "timestampGreaterThan")
	}
	if rule.TimestampLessThanEqualsPath != "" {
		return e.comparePathValue(inputData, rule.TimestampLessThanEqualsPath, varValue, "timestampLessThanEquals")
	}
	if rule.TimestampGreaterThanEqualsPath != "" {
		return e.comparePathValue(inputData, rule.TimestampGreaterThanEqualsPath, varValue, "timestampGreaterThanEquals")
	}

	return false
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

// comparePathValue reads a comparison value from inputData at the given
// JSONPath and compares it against varValue according to op.
func (e *Executor) comparePathValue(inputData map[string]interface{}, path string, varValue interface{}, op string) bool {
	compareVal, exists := getJSONPathValueRaw(inputData, path)
	if !exists {
		return false
	}
	switch op {
	case "stringEquals":
		return fmt.Sprintf("%v", varValue) == fmt.Sprintf("%v", compareVal)
	case "stringLessThan":
		return fmt.Sprintf("%v", varValue) < fmt.Sprintf("%v", compareVal)
	case "stringGreaterThan":
		return fmt.Sprintf("%v", varValue) > fmt.Sprintf("%v", compareVal)
	case "stringLessThanEquals":
		return fmt.Sprintf("%v", varValue) <= fmt.Sprintf("%v", compareVal)
	case "stringGreaterThanEquals":
		return fmt.Sprintf("%v", varValue) >= fmt.Sprintf("%v", compareVal)
	case "numericEquals":
		a, ok1 := toFloat64(varValue)
		b, ok2 := toFloat64(compareVal)
		return ok1 && ok2 && a == b
	case "numericLessThan":
		a, ok1 := toFloat64(varValue)
		b, ok2 := toFloat64(compareVal)
		return ok1 && ok2 && a < b
	case "numericGreaterThan":
		a, ok1 := toFloat64(varValue)
		b, ok2 := toFloat64(compareVal)
		return ok1 && ok2 && a > b
	case "numericLessThanEquals":
		a, ok1 := toFloat64(varValue)
		b, ok2 := toFloat64(compareVal)
		return ok1 && ok2 && a <= b
	case "numericGreaterThanEquals":
		a, ok1 := toFloat64(varValue)
		b, ok2 := toFloat64(compareVal)
		return ok1 && ok2 && a >= b
	case "booleanEquals":
		va, ok1 := varValue.(bool)
		vb, ok2 := compareVal.(bool)
		return ok1 && ok2 && va == vb
	case "timestampEquals":
		return e.compareTimestamp(varValue, fmt.Sprintf("%v", compareVal), "equals")
	case "timestampLessThan":
		return e.compareTimestamp(varValue, fmt.Sprintf("%v", compareVal), "less")
	case "timestampGreaterThan":
		return e.compareTimestamp(varValue, fmt.Sprintf("%v", compareVal), "greater")
	case "timestampLessThanEquals":
		return e.compareTimestamp(varValue, fmt.Sprintf("%v", compareVal), "lessEquals")
	case "timestampGreaterThanEquals":
		return e.compareTimestamp(varValue, fmt.Sprintf("%v", compareVal), "greaterEquals")
	}
	return false
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

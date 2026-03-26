package sfn

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

func (e *Executor) executeChoice(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.ChoiceState) (string, error) {
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

	if rule.StringEquals != nil {
		if expected, ok := rule.StringEquals[rule.Variable]; ok {
			return fmt.Sprintf("%v", varValue) == expected
		}
		for varPath, expected := range rule.StringEquals {
			if val, ok := getJSONPathValueRaw(inputData, varPath); ok {
				return fmt.Sprintf("%v", val) == expected
			}
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

	if rule.BooleanEquals != nil {
		if expected, ok := rule.BooleanEquals[rule.Variable]; ok {
			if boolVal, ok := varValue.(bool); ok {
				return boolVal == expected
			}
		}
	}

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
	}
	return false
}

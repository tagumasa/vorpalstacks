package sfn

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

func (e *Executor) executePass(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.PassState) (string, string, error) {
	isJSONata := IsJSONataState(state, execCtx.QueryLanguage)

	if isJSONata {
		return e.executePassJSONata(ctx, execCtx, state)
	}

	processedInput := e.applyInputPath(execCtx.Input, state.GetInputPath())
	processedInput = e.applyParameters(processedInput, state.Parameters)

	eventId := execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "PassStateEntered",
		Timestamp:       time.Now().UTC(),
		PassStateEnteredEventDetails: &sfnstore.PassStateEnteredEventDetails{
			Input: processedInput,
			Name:  execCtx.CurrentState,
		},
	})

	output := processedInput
	if state.Result != nil {
		resultJSON, err := json.Marshal(state.Result)
		if err != nil {
			return "", "", fmt.Errorf("failed to marshal result: %w", err)
		}
		output = string(resultJSON)
	}

	output = e.applyResultSelector(output, state.GetResultSelector())
	output = e.applyResultPath(processedInput, output, state.ResultPath)
	output = e.applyOutputPath(output, state.GetOutputPath())

	eventId = execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "PassStateExited",
		Timestamp:       time.Now().UTC(),
		PassStateExitedEventDetails: &sfnstore.PassStateExitedEventDetails{
			Output: output,
			Name:   execCtx.CurrentState,
		},
	})

	return output, state.Next, nil
}

func (e *Executor) executePassJSONata(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.PassState) (string, string, error) {
	var inputData interface{}
	if execCtx.Input != "" {
		if err := json.Unmarshal([]byte(execCtx.Input), &inputData); err != nil {
			return "", "", e.newQueryEvalError(ctx, execCtx, "Input", "failed to parse input JSON")
		}
	}

	if state.JSONataOutput == nil && len(state.OutputRaw) > 0 {
		var err error
		state.JSONataOutput, err = resolveJSONataOutput(state)
		if err != nil {
			return "", "", e.newQueryEvalError(ctx, execCtx, "Output", err.Error())
		}
	}

	statesVar := e.buildStatesVarWithContext(execCtx, inputData, nil, nil)

	eventId := execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "PassStateEntered",
		Timestamp:       time.Now().UTC(),
		PassStateEnteredEventDetails: &sfnstore.PassStateEnteredEventDetails{
			Input: execCtx.Input,
			Name:  execCtx.CurrentState,
		},
	})

	if len(state.Assign) > 0 {
		evaluated, err := evaluateAssign(ctx, state.Assign, statesVar, execCtx.VariableScope)
		if err != nil {
			return "", "", e.newQueryEvalError(ctx, execCtx, "Assign", err.Error())
		}
		execCtx.PendingAssign = evaluated
	}

	output := inputData
	if state.JSONataOutput != nil {
		resolved, err := e.applyJSONataOutput(ctx, state.JSONataOutput, statesVar, execCtx.VariableScope)
		if err != nil {
			return "", "", e.newQueryEvalError(ctx, execCtx, "Output", err.Error())
		}
		output = resolved
	} else if state.Result != nil {
		output = state.Result
	}

	outputJSON, err := json.Marshal(output)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal output: %w", err)
	}
	outputStr := string(outputJSON)

	eventId = execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "PassStateExited",
		Timestamp:       time.Now().UTC(),
		PassStateExitedEventDetails: &sfnstore.PassStateExitedEventDetails{
			Output: outputStr,
			Name:   execCtx.CurrentState,
		},
	})

	return outputStr, state.Next, nil
}

func (e *Executor) executeWait(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.WaitState) (string, string, error) {
	isJSONata := IsJSONataState(state, execCtx.QueryLanguage)

	if isJSONata {
		return e.executeWaitJSONata(ctx, execCtx, state)
	}

	processedInput := e.applyInputPath(execCtx.Input, state.GetInputPath())

	var waitSeconds int32 = state.GetSeconds()

	if state.SecondsPath != "" {
		var inputData map[string]interface{}
		if err := json.Unmarshal([]byte(processedInput), &inputData); err == nil {
			if val, ok := getJSONPathValueRaw(inputData, state.SecondsPath); ok {
				if numVal, ok := toFloat64(val); ok {
					waitSeconds = int32(numVal)
				}
			}
		}
	}

	eventId := execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "WaitStateEntered",
		Timestamp:       time.Now().UTC(),
		WaitStateEnteredEventDetails: &sfnstore.WaitStateEnteredEventDetails{
			Input:   processedInput,
			Name:    execCtx.CurrentState,
			Seconds: int64(waitSeconds),
		},
	})

	if waitSeconds > 0 {
		timer := time.NewTimer(time.Duration(waitSeconds) * time.Second)
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return "", "", ctx.Err()
		}
	}

	output := e.applyOutputPath(processedInput, state.GetOutputPath())

	eventId = execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "WaitStateExited",
		Timestamp:       time.Now().UTC(),
		WaitStateExitedEventDetails: &sfnstore.WaitStateExitedEventDetails{
			Output: output,
			Name:   execCtx.CurrentState,
		},
	})

	return output, state.Next, nil
}

func (e *Executor) executeWaitJSONata(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.WaitState) (string, string, error) {
	var inputData interface{}
	if execCtx.Input != "" {
		if err := json.Unmarshal([]byte(execCtx.Input), &inputData); err != nil {
			return "", "", e.newQueryEvalError(ctx, execCtx, "Input", "failed to parse input JSON")
		}
	}

	statesVar := e.buildStatesVarWithContext(execCtx, inputData, nil, nil)

	eventId := execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "WaitStateEntered",
		Timestamp:       time.Now().UTC(),
		WaitStateEnteredEventDetails: &sfnstore.WaitStateEnteredEventDetails{
			Input: execCtx.Input,
			Name:  execCtx.CurrentState,
		},
	})

	if len(state.Assign) > 0 {
		evaluated, err := evaluateAssign(ctx, state.Assign, statesVar, execCtx.VariableScope)
		if err != nil {
			return "", "", e.newQueryEvalError(ctx, execCtx, "Assign", err.Error())
		}
		execCtx.PendingAssign = evaluated
	}

	var waitSeconds int32 = state.GetSeconds()
	if s, ok := state.Seconds.(string); ok && IsExpression(s) {
		vars := buildVarsMap(statesVar, execCtx.VariableScope)
		result, err := EvaluateJSONata(ctx, UnwrapExpression(s), nil, vars)
		if err != nil {
			return "", "", e.newQueryEvalError(ctx, execCtx, "Seconds", err.Error())
		}
		if f, ok := toFloat64(result); ok {
			waitSeconds = int32(f)
		}
	}

	if waitSeconds > 0 {
		timer := time.NewTimer(time.Duration(waitSeconds) * time.Second)
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return "", "", ctx.Err()
		}
	}

	outputJSON, err := json.Marshal(inputData)
	if err != nil {
		outputJSON = []byte("{}")
	}

	eventId = execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "WaitStateExited",
		Timestamp:       time.Now().UTC(),
		WaitStateExitedEventDetails: &sfnstore.WaitStateExitedEventDetails{
			Output: string(outputJSON),
			Name:   execCtx.CurrentState,
		},
	})

	return string(outputJSON), state.Next, nil
}

func (e *Executor) executeFail(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.FailState) error {
	cause := state.GetCause()
	if IsJSONataState(state, execCtx.QueryLanguage) && state.Cause != nil {
		if s, ok := state.Cause.(string); ok && IsExpression(s) {
			var inputData interface{}
			if execCtx.Input != "" {
				if err := json.Unmarshal([]byte(execCtx.Input), &inputData); err != nil {
					return e.newQueryEvalError(ctx, execCtx, "Input", "failed to parse input JSON")
				}
			}
			statesVar := e.buildStatesVarWithContext(execCtx, inputData, nil, nil)
			vars := buildVarsMap(statesVar, execCtx.VariableScope)
			result, err := EvaluateJSONata(ctx, UnwrapExpression(s), nil, vars)
			if err != nil {
				_ = e.newQueryEvalError(ctx, execCtx, "Cause", err.Error())
				cause = fmt.Sprintf("QueryEvaluationError: %s", err.Error())
			} else {
				switch v := result.(type) {
				case string:
					cause = v
				default:
					b, _ := json.Marshal(result)
					cause = string(b)
				}
			}
		}
	}

	eventId := execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "FailStateEntered",
		Timestamp:       time.Now().UTC(),
		FailStateEnteredEventDetails: &sfnstore.FailStateEnteredEventDetails{
			Input: execCtx.Input,
			Name:  execCtx.CurrentState,
			Error: state.Error,
			Cause: cause,
		},
	})

	return fmt.Errorf("state machine failed: %s - %s", state.Error, cause)
}

func (e *Executor) executeSucceed(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.SucceedState) (string, string, error) {
	isJSONata := IsJSONataState(state, execCtx.QueryLanguage)

	if isJSONata {
		if state.JSONataOutput == nil && len(state.OutputRaw) > 0 {
			var err error
			state.JSONataOutput, err = resolveJSONataOutput(state)
			if err != nil {
				return "", "", e.newQueryEvalError(ctx, execCtx, "Output", err.Error())
			}
		}
		if state.JSONataOutput != nil {
			var inputData interface{}
			if execCtx.Input != "" {
				if err := json.Unmarshal([]byte(execCtx.Input), &inputData); err != nil {
					return "", "", e.newQueryEvalError(ctx, execCtx, "Input", "failed to parse input JSON")
				}
			}
			statesVar := e.buildStatesVarWithContext(execCtx, inputData, nil, nil)

			eventId := execCtx.nextEventId()
			e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
				ExecutionArn:    execCtx.Execution.ExecutionArn,
				EventId:         eventId,
				PreviousEventId: eventId - 1,
				Type:            "SucceedStateEntered",
				Timestamp:       time.Now().UTC(),
				SucceedStateEnteredEventDetails: &sfnstore.SucceedStateEnteredEventDetails{
					Input: execCtx.Input,
					Name:  execCtx.CurrentState,
				},
			})

			resolved, err := e.applyJSONataOutput(ctx, state.JSONataOutput, statesVar, execCtx.VariableScope)
			if err != nil {
				return "", "", e.newQueryEvalError(ctx, execCtx, "Output", err.Error())
			}
			outputJSON, err := json.Marshal(resolved)
			if err != nil {
				return "", "", fmt.Errorf("failed to marshal output: %w", err)
			}
			return string(outputJSON), "", nil
		}
	}

	processedInput := e.applyInputPath(execCtx.Input, state.GetInputPath())

	eventId := execCtx.nextEventId()
	e.logHistoryEvent(ctx, execCtx.Execution, &sfnstore.ExecutionHistoryEvent{
		ExecutionArn:    execCtx.Execution.ExecutionArn,
		EventId:         eventId,
		PreviousEventId: eventId - 1,
		Type:            "SucceedStateEntered",
		Timestamp:       time.Now().UTC(),
		SucceedStateEnteredEventDetails: &sfnstore.SucceedStateEnteredEventDetails{
			Input: execCtx.Input,
			Name:  execCtx.CurrentState,
		},
	})

	output := e.applyOutputPath(processedInput, state.GetOutputPath())

	return output, "", nil
}

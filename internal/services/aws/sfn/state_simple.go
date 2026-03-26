package sfn

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

func (e *Executor) executePass(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.PassState) (string, string, error) {
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

func (e *Executor) executeWait(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.WaitState) (string, string, error) {
	processedInput := e.applyInputPath(execCtx.Input, state.GetInputPath())

	var waitSeconds int32 = state.Seconds

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

func (e *Executor) executeFail(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.FailState) error {
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
			Cause: state.Cause,
		},
	})

	return fmt.Errorf("state machine failed: %s - %s", state.Error, state.Cause)
}

func (e *Executor) executeSucceed(ctx context.Context, execCtx *ExecutionContext, state *sfnstore.SucceedState) (string, string, error) {
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

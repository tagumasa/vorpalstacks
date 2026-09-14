package sfn

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"time"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

func (e *Executor) findMatchingRetryPolicy(policies []*sfnstore.RetryPolicy, errorCode string) *sfnstore.RetryPolicy {
	for _, policy := range policies {
		for _, errPattern := range policy.ErrorEquals {
			if e.errorMatchesPattern(errorCode, errPattern) {
				return policy
			}
		}
	}
	return nil
}

func (e *Executor) findMatchingCatchPolicy(policies []*sfnstore.CatchPolicy, errorCode string) *sfnstore.CatchPolicy {
	for _, policy := range policies {
		for _, errPattern := range policy.ErrorEquals {
			if e.errorMatchesPattern(errorCode, errPattern) {
				return policy
			}
		}
	}
	return nil
}

func (e *Executor) errorMatchesPattern(errorCode, pattern string) bool {
	if pattern == "States.ALL" {
		// States.ALL matches all errors except States.DataLimitExceeded
		// and States.Runtime (RuntimeExceeded).
		excludedByAll := map[string]bool{
			"States.DataLimitExceeded": true,
			"States.RuntimeExceeded":   true,
			"States.Runtime":           true,
		}
		return !excludedByAll[errorCode]
	}
	if pattern == "States.TaskFailed" {
		// "The name States.TaskFailed also acts a wildcard and matches any
		// error except for States.Timeout" (error handling documentation) —
		// a prefix match would both miss foreign names (Lambda.TooManyRequestsException)
		// and wrongly admit States.Timeout.
		return errorCode != "States.Timeout"
	}
	if pattern == errorCode {
		return true
	}
	return false
}

func (e *Executor) calculateBackoffInterval(policy *sfnstore.RetryPolicy, attempt int32) time.Duration {
	interval := float64(policy.IntervalSeconds)
	if interval == 0 {
		interval = 1
	}
	backoffRate := policy.BackoffRate
	if backoffRate == 0 {
		backoffRate = 2.0
	}
	interval = interval * math.Pow(backoffRate, float64(attempt-1))
	if policy.MaxDelaySeconds > 0 && interval > float64(policy.MaxDelaySeconds) {
		interval = float64(policy.MaxDelaySeconds)
	}
	if policy.JitterStrategy == "FULL" {
		// "If you set JitterStrategy as FULL, the first retry interval is
		// randomized between 0 and 2 seconds, the second retry interval is
		// randomized between 0 and 4 seconds" — with IntervalSeconds 2 the
		// window is the attempt's deterministic interval itself
		// (MaxDelaySeconds cap included), not a multiple of it.
		interval = rand.Float64() * interval
	}
	return time.Duration(interval) * time.Second
}

// sleepForRetry waits the retrier's backoff interval for the given attempt
// number and reports true when the execution was interrupted instead of
// the interval elapsing.
func (e *Executor) sleepForRetry(ctx context.Context, policy *sfnstore.RetryPolicy, attempt int32) bool {
	timer := time.NewTimer(e.calculateBackoffInterval(policy, attempt))
	select {
	case <-timer.C:
		return false
	case <-ctx.Done():
		timer.Stop()
		return true
	}
}

// buildCatchOutput renders a Catch handler's error output, folding it into
// the state input at the handler's ResultPath. A configured ResultPath that
// cannot apply to the input fails with States.ResultPathMatchFailure, as on
// the state's own ResultPath.
func (e *Executor) buildCatchOutput(input, errorCode, cause, resultPath string) (string, *ExecutionError) {
	errorInfo := map[string]interface{}{
		"Error": errorCode,
		"Cause": cause,
	}
	errorJSON, err := json.Marshal(errorInfo)
	if err != nil {
		// The fallback must still produce valid JSON: each member is
		// marshalled on its own, so control characters and quotes inside
		// the error code or cause stay escaped instead of breaking the
		// document a hand-assembled Sprintf would produce.
		codeJSON, codeErr := json.Marshal(errorCode)
		causeJSON, causeErr := json.Marshal(cause)
		if codeErr != nil || causeErr != nil {
			return `{"Error":"","Cause":""}`, nil
		}
		return fmt.Sprintf(`{"Error":%s,"Cause":%s}`, codeJSON, causeJSON), nil
	}

	if resultPath == "" {
		return string(errorJSON), nil
	}

	var inputData map[string]interface{}
	if err := json.Unmarshal([]byte(input), &inputData); err != nil {
		return "", &ExecutionError{ErrorCode: "States.ResultPathMatchFailure", Cause: fmt.Sprintf("ResultPath %s cannot apply to a non-object state input", resultPath)}
	}

	setNestedPath(inputData, resultPath, errorInfo)
	mergedJSON, err := json.Marshal(inputData)
	if err != nil {
		return string(errorJSON), nil
	}
	return string(mergedJSON), nil
}

package sfn

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
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
		return true
	}
	if pattern == "States.Timeout" && strings.HasPrefix(errorCode, "States.Timeout") {
		return true
	}
	if pattern == "States.TaskFailed" && strings.HasPrefix(errorCode, "States.TaskFailed") {
		return true
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
	return time.Duration(interval) * time.Second
}

func (e *Executor) buildCatchOutput(input, errorCode, cause, resultPath string) string {
	errorInfo := map[string]interface{}{
		"Error": errorCode,
		"Cause": cause,
	}
	errorJSON, err := json.Marshal(errorInfo)
	if err != nil {
		return fmt.Sprintf(`{"Error":"%s","Cause":"%s"}`, errorCode, cause)
	}

	if resultPath == "" {
		return string(errorJSON)
	}

	var inputData map[string]interface{}
	if err := json.Unmarshal([]byte(input), &inputData); err != nil {
		return string(errorJSON)
	}

	setNestedPath(inputData, resultPath, errorInfo)
	mergedJSON, err := json.Marshal(inputData)
	if err != nil {
		return string(errorJSON)
	}
	return string(mergedJSON)
}

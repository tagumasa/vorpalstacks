// Package schedules provides AWS EventBridge schedule parsing utilities for vorpalstacks.
package schedules

import (
	"strings"
)

func parseScheduleExpression(expression string) (scheduleType, string, error) {
	expression = strings.TrimSpace(expression)

	if strings.HasPrefix(expression, "rate(") {
		content := strings.TrimSuffix(strings.TrimPrefix(expression, "rate("), ")")
		return scheduleTypeRate, strings.TrimSpace(content), nil
	}

	if strings.HasPrefix(expression, "cron(") {
		content := strings.TrimSuffix(strings.TrimPrefix(expression, "cron("), ")")
		return scheduleTypeCron, strings.TrimSpace(content), nil
	}

	return "", "", nil
}

type scheduleType string

const (
	scheduleTypeRate scheduleType = "rate"
	scheduleTypeCron scheduleType = "cron"
)

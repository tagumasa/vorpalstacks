// Package validators provides AWS-specific utility functions for vorpalstacks.
package validators

import (
	"vorpalstacks/internal/utils/aws/schedules"
)

// IsValidScheduleExpression returns true if the given expression is a valid
// AWS EventBridge schedule expression (rate or cron format).
//
// Parameters:
//   - expression: The schedule expression to validate
//
// Returns:
//   - bool: True if the expression is valid
func IsValidScheduleExpression(expression string) bool {
	scheduleType, err := schedules.GetScheduleType(expression)
	if err != nil {
		return false
	}

	if scheduleType == schedules.ScheduleTypeRate {
		_, err := schedules.GetRateSchedule(expression)
		return err == nil
	}

	if scheduleType == schedules.ScheduleTypeCron {
		_, err := schedules.GetCronSchedule(expression)
		return err == nil
	}

	return false
}

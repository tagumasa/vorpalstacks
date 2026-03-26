// Package schedules provides AWS-specific utility functions for vorpalstacks.
package schedules

import (
	"fmt"
)

// ScheduleType represents the type of a schedule expression.
type ScheduleType string

// Schedule type constants for rate and cron expressions.
const (
	ScheduleTypeRate ScheduleType = "rate"
	ScheduleTypeCron ScheduleType = "cron"
)

// GetScheduleType determines the type of a schedule expression (rate or cron).
//
// Parameters:
//   - expression: The schedule expression to check
//
// Returns:
//   - ScheduleType: The type of schedule (ScheduleTypeRate or ScheduleTypeCron)
//   - error: An error if the expression is invalid
//
// Example:
//
//	scheduleType, err := GetScheduleType("rate(1 hour)")
func GetScheduleType(expression string) (ScheduleType, error) {
	scheduleType, _, err := parseScheduleExpression(expression)
	if err != nil {
		return "", err
	}
	if scheduleType == "" {
		return "", fmt.Errorf("invalid schedule expression: must start with 'rate(' or 'cron('")
	}
	return ScheduleType(scheduleType), nil
}

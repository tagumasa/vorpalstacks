// Package schedules provides AWS-specific utility functions for vorpalstacks.
package schedules

import (
	"fmt"
	"strings"
)

// CronSchedule represents a cron-based schedule expression.
type CronSchedule struct {
	// Minutes is the minutes field (0-59)
	Minutes string
	// Hours is the hours field (0-23)
	Hours string
	// DayOfMonth is the day of month field (1-31)
	DayOfMonth string
	// Month is the month field (1-12 or JAN-DEC)
	Month string
	// DayOfWeek is the day of week field (0-6 or SUN-SAT)
	DayOfWeek string
	// Year is the year field (optional)
	Year string
}

// GetCronSchedule parses an AWS cron-based schedule expression.
// Cron expressions have 6 fields: minutes, hours, day of month, month, day of week, and year.
//
// Parameters:
//   - expression: The cron expression (e.g., "cron(0 12 * * ? *)")
//
// Returns:
//   - *CronSchedule: The parsed schedule
//   - error: An error if parsing fails
//
// Example:
//
//	schedule, err := GetCronSchedule("cron(0 12 * * ? *)")
func GetCronSchedule(expression string) (*CronSchedule, error) {
	scheduleType, content, err := parseScheduleExpression(expression)
	if err != nil {
		return nil, err
	}
	if scheduleType != scheduleTypeCron {
		return nil, fmt.Errorf("not a cron expression: %s", expression)
	}

	fields := strings.Fields(content)
	if len(fields) != 6 {
		return nil, fmt.Errorf("invalid cron expression: must have 6 fields, got %d", len(fields))
	}

	return &CronSchedule{
		Minutes:    fields[0],
		Hours:      fields[1],
		DayOfMonth: fields[2],
		Month:      fields[3],
		DayOfWeek:  fields[4],
		Year:       fields[5],
	}, nil
}

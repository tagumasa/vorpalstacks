// Package schedules provides AWS-specific utility functions for vorpalstacks.
package schedules

import (
	"fmt"
	"strconv"
	"strings"
)

// RateSchedule represents a rate-based schedule expression.
type RateSchedule struct {
	// Value is the numeric value for the rate (e.g., 5)
	Value int
	// Unit is the time unit (minute, minutes, hour, hours, day, days)
	Unit string
}

// GetRateSchedule parses an AWS rate-based schedule expression.
// Rate expressions start with "rate(" and specify how often an event occurs.
//
// Parameters:
//   - expression: The rate expression (e.g., "rate(5 minutes)")
//
// Returns:
//   - *RateSchedule: The parsed schedule
//   - error: An error if parsing fails
//
// Example:
//
//	schedule, err := GetRateSchedule("rate(1 hour)")
func GetRateSchedule(expression string) (*RateSchedule, error) {
	scheduleType, content, err := parseScheduleExpression(expression)
	if err != nil {
		return nil, err
	}
	if scheduleType != scheduleTypeRate {
		return nil, fmt.Errorf("not a rate expression: %s", expression)
	}

	parts := strings.Fields(content)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid rate expression: %s", expression)
	}

	value, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid rate value: %s", parts[0])
	}

	if value < 1 {
		return nil, fmt.Errorf("rate value must be at least 1: %d", value)
	}

	unit := strings.ToLower(parts[1])

	validUnits := map[string]bool{
		"minute": true, "minutes": true,
		"hour": true, "hours": true,
		"day": true, "days": true,
	}

	if !validUnits[unit] {
		return nil, fmt.Errorf("invalid rate unit: %s (must be minute/minutes, hour/hours, or day/days)", parts[1])
	}

	return &RateSchedule{
		Value: value,
		Unit:  unit,
	}, nil
}

package eventbridge

import (
	"encoding/json"
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
)

// testEventPatternCore validates the EventPattern and Event members and
// reports whether the event matches the pattern: both must be non-empty
// JSON objects (the pattern additionally capped at 4096 characters), and
// every top-level pattern key must match the event value — a key absent
// from the event only satisfies an explicit exists-false pattern.
func (s *EventsService) testEventPatternCore(patternStr, eventStr string) (bool, error) {
	if patternStr == "" {
		return false, awserrors.NewValidationException("Parameter EventPattern is required")
	}
	if !validateEventPatternLength(patternStr) {
		return false, awserrors.NewValidationException("EventPattern must be at most 4096 characters")
	}
	if eventStr == "" {
		return false, awserrors.NewValidationException("Parameter Event is required")
	}

	var patternMap, eventMap map[string]interface{}
	if err := json.Unmarshal([]byte(patternStr), &patternMap); err != nil {
		return false, awserrors.NewInvalidEventPatternException(fmt.Sprintf("EventPattern is not valid JSON: %s", err))
	}
	if err := json.Unmarshal([]byte(eventStr), &eventMap); err != nil {
		return false, awserrors.NewValidationException(fmt.Sprintf("Event is not valid JSON: %s", err))
	}

	result := true
	for key, patternValue := range patternMap {
		eventValue, exists := eventMap[key]
		if !exists {
			if isExistsFalsePattern(patternValue) {
				continue
			}
			result = false
			break
		}
		if !s.matchValue(eventValue, patternValue) {
			result = false
			break
		}
	}
	return result, nil
}

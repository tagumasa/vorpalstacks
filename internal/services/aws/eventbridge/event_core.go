package eventbridge

import (
	"encoding/json"
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
)

// mandatoryEventMembers lists the members the model's TestEventPattern
// Event documentation makes mandatory for the event under test: "The event,
// in JSON format, to test against the event pattern. The JSON must follow
// the format specified in AWS Events, and the following fields are
// mandatory: id, account, source, time, region, resources, detail-type."
var mandatoryEventMembers = []string{
	"id", "account", "source", "time", "region", "resources", "detail-type",
}

// testEventPatternCore validates the EventPattern and Event members and
// reports whether the event matches the pattern: the pattern must pass the
// structural validation shared with every acceptance site, the event must
// be a JSON object carrying the mandatory members, and evaluation runs
// through the same level matcher the delivery plane uses.
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

	if err := validateEventPatternStructure(patternStr); err != nil {
		return false, err
	}

	var patternMap map[string]interface{}
	if err := json.Unmarshal([]byte(patternStr), &patternMap); err != nil {
		return false, awserrors.NewInvalidEventPatternException(fmt.Sprintf("EventPattern is not valid JSON: %s", err))
	}
	var eventMap map[string]interface{}
	if err := json.Unmarshal([]byte(eventStr), &eventMap); err != nil {
		return false, awserrors.NewValidationException(fmt.Sprintf("Event is not valid JSON: %s", err))
	}
	for _, member := range mandatoryEventMembers {
		if _, present := eventMap[member]; !present {
			return false, awserrors.NewValidationException("Event is missing mandatory member: " + member)
		}
	}

	return s.matchPatternLevel(eventMap, patternMap), nil
}

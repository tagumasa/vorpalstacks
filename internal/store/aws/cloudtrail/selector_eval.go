package cloudtrail

import "strings"

// Advanced event selector evaluation — the one evaluator for both trail
// log-file delivery and event data store ingestion, the two surfaces an
// advanced selector governs. The operator combination semantics are the
// documented ones (CloudTrail User Guide, "How CloudTrail evaluates
// multiple conditions for a field"):
//
//   - DESELECT operators (NotEquals, NotStartsWith, NotEndsWith) are
//     AND'd together: if any DESELECT operator condition is met, the
//     event is not delivered.
//   - SELECT operators (Equals, StartsWith, EndsWith) are OR'd together:
//     at least one SELECT operator condition must be met.
//   - Combinations of SELECT and DESELECT operators follow the above
//     rules and both groups are AND'd together.
//
// Two fail-closed rulings predate this evaluator and are preserved: a
// field whose value the event does not carry never holds, and a field
// selector carrying no operator values at all never holds — both exclude
// the event rather than deliver it unfiltered.

// advancedFieldVocabulary is the accepted advanced-selector field set with
// each field's documented operator restriction: eventCategory, readOnly,
// resources.type and errorCode accept Equals only,
// sessionCredentialFromConsole accepts Equals and NotEquals only, and a
// nil entry means the field permits all six operator lists
// (AdvancedFieldSelector API reference and the user guide's data-event
// field table). eventSource's Equals-only restriction is conditional on
// the selector's event category (network activity selectors) and lives in
// the selector-level validator, not this per-field table. username is the
// one addition beyond the documented set — this platform's event records
// carry it as a first-class member, and both ingestion surfaces resolve
// it.
var advancedFieldVocabulary = map[string]map[string]bool{
	"eventCategory":                {"Equals": true},
	"readOnly":                     {"Equals": true},
	"resources.type":               {"Equals": true},
	"errorCode":                    {"Equals": true},
	"sessionCredentialFromConsole": {"Equals": true, "NotEquals": true},
	"eventSource":                  nil,
	"eventName":                    nil,
	"eventType":                    nil,
	"userIdentity.arn":             nil,
	"resources.ARN":                nil,
	"vpcEndpointId":                nil,
	"username":                     nil,
}

// AdvancedSelectorSelectsNetworkActivity reports whether the advanced
// selector's eventCategory pins network activity events — the category
// whose eventSource field the documentation both requires and restricts
// to Equals ("For network activity events, this is a required field that
// only uses the Equals operator").
func AdvancedSelectorSelectsNetworkActivity(sel AdvancedEventSelector) bool {
	for _, fs := range sel.FieldSelectors {
		if !strings.EqualFold(fs.Field, "eventCategory") {
			continue
		}
		for _, v := range fs.Equals {
			if strings.EqualFold(v, "NetworkActivity") {
				return true
			}
		}
	}
	return false
}

// AdvancedSelectorFieldRules reports whether field belongs to the
// advanced event selector vocabulary, and the operator lists the field is
// restricted to when the documentation restricts it (nil means all six
// lists are permitted). Field comparison folds case, matching the
// evaluator.
func AdvancedSelectorFieldRules(field string) (restrictedTo map[string]bool, known bool) {
	for canonical, restricted := range advancedFieldVocabulary {
		if strings.EqualFold(field, canonical) {
			return restricted, true
		}
	}
	return nil, false
}

// AdvancedSelectorMatches reports whether one advanced event selector
// holds for the event: every field selector must hold.
func AdvancedSelectorMatches(sel AdvancedEventSelector, e *Event) bool {
	for _, fs := range sel.FieldSelectors {
		if !advancedFieldSelectorHolds(fs, e) {
			return false
		}
	}
	return true
}

// advancedFieldSelectorHolds evaluates one field selector. The event's
// value set for the field must be non-empty, the selector must carry at
// least one operator value, and one of the values must satisfy both
// operator groups. resources.type and resources.ARN resolve one value per
// referenced resource; the per-value evaluation is the documented
// single-value semantics generalised to the field's value set — a value
// that satisfies the selector qualifies the event.
func advancedFieldSelectorHolds(fs AdvancedFieldSelector, e *Event) bool {
	if !advancedFieldSelectorHasValues(fs) {
		return false
	}
	values, ok := advancedFieldValues(fs.Field, e)
	if !ok {
		return false
	}
	for _, value := range values {
		if selectorOperatorGroupsHold(fs, value) {
			return true
		}
	}
	return false
}

// advancedFieldSelectorHasValues reports whether any operator list of the
// field selector carries at least one value.
func advancedFieldSelectorHasValues(fs AdvancedFieldSelector) bool {
	return len(fs.Equals) > 0 || len(fs.NotEquals) > 0 ||
		len(fs.StartsWith) > 0 || len(fs.NotStartsWith) > 0 ||
		len(fs.EndsWith) > 0 || len(fs.NotEndsWith) > 0
}

// selectorOperatorGroupsHold applies the documented group semantics to one
// field value: at least one SELECT operator (Equals, StartsWith, EndsWith)
// condition must be met — an absent SELECT group holds vacuously, the
// deselect-only form the documentation shows — and no DESELECT operator
// (NotEquals, NotStartsWith, NotEndsWith) condition may be met.
func selectorOperatorGroupsHold(fs AdvancedFieldSelector, value string) bool {
	selectMet := len(fs.Equals) == 0 && len(fs.StartsWith) == 0 && len(fs.EndsWith) == 0
	for _, w := range fs.Equals {
		if value == w {
			selectMet = true
			break
		}
	}
	if !selectMet {
		for _, w := range fs.StartsWith {
			if strings.HasPrefix(value, w) {
				selectMet = true
				break
			}
		}
	}
	if !selectMet {
		for _, w := range fs.EndsWith {
			if strings.HasSuffix(value, w) {
				selectMet = true
				break
			}
		}
	}
	if !selectMet {
		return false
	}
	for _, w := range fs.NotEquals {
		if value == w {
			return false
		}
	}
	for _, w := range fs.NotStartsWith {
		if strings.HasPrefix(value, w) {
			return false
		}
	}
	for _, w := range fs.NotEndsWith {
		if strings.HasSuffix(value, w) {
			return false
		}
	}
	return true
}

// advancedFieldValues resolves an advanced-selector field name to the
// event's value(s) for it. sessionCredentialFromConsole resolves "false"
// for every event: the record member is "not shown unless the value is
// true" (CloudTrail record contents) and this platform's event records
// never mark a console session, so no event is console-originated.
// vpcEndpointId resolves no value: the event records this platform writes
// carry no VPC endpoint identity and absence has no default, so a
// selector on it excludes every event — fail-closed, like any field the
// event does not hold.
func advancedFieldValues(field string, e *Event) ([]string, bool) {
	switch {
	case strings.EqualFold(field, "eventCategory"):
		return singletonFieldValue(e.EventCategory)
	case strings.EqualFold(field, "eventSource"):
		return singletonFieldValue(e.EventSource)
	case strings.EqualFold(field, "eventName"):
		return singletonFieldValue(e.EventName)
	case strings.EqualFold(field, "eventType"):
		return singletonFieldValue(e.EventType)
	case strings.EqualFold(field, "username"):
		if e.UserIdentity == nil {
			return nil, false
		}
		return singletonFieldValue(e.UserIdentity.UserName)
	case strings.EqualFold(field, "userIdentity.arn"):
		if e.UserIdentity == nil {
			return nil, false
		}
		return singletonFieldValue(e.UserIdentity.ARN)
	case strings.EqualFold(field, "readOnly"):
		return singletonFieldValue(e.ReadOnly)
	case strings.EqualFold(field, "resources.type"):
		values := make([]string, 0, len(e.Resources))
		for _, r := range e.Resources {
			values = append(values, r.ResourceType)
		}
		return values, len(values) > 0
	case strings.EqualFold(field, "resources.ARN"):
		values := make([]string, 0, len(e.Resources))
		for _, r := range e.Resources {
			values = append(values, r.ResourceName)
		}
		return values, len(values) > 0
	case strings.EqualFold(field, "errorCode"):
		return singletonFieldValue(e.ErrorCode)
	case strings.EqualFold(field, "sessionCredentialFromConsole"):
		// The member is absent unless true (CloudTrail record contents),
		// and no event this platform writes carries a console-session
		// marker, so every event's value is false.
		return []string{"false"}, true
	}
	return nil, false
}

// singletonFieldValue wraps a single-valued field, reporting ok=false for
// the empty value.
func singletonFieldValue(v string) ([]string, bool) {
	if v == "" {
		return nil, false
	}
	return []string{v}, true
}

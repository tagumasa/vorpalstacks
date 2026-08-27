package iot

import (
	"regexp"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// Scheduled Audit Core (scheduled audit configuration). A scheduled audit
// names the Device Defender checks that run periodically; records live
// under the generic-KV "scheduledAudit/" prefix.
// ---------------------------------------------------------------------------

// ScheduledAuditInput carries the mutable fields of a scheduled audit. The
// *Provided flags distinguish explicitly supplied members from omitted ones
// so a partial update leaves the omitted members unchanged.
type ScheduledAuditInput struct {
	Name                     string
	Frequency                string
	FrequencyProvided        bool
	DayOfMonth               string
	DayOfMonthProvided       bool
	DayOfWeek                string
	DayOfWeekProvided        bool
	TargetCheckNames         []string
	TargetCheckNamesProvided bool
	Tags                     map[string]string
}

// auditFrequencies is the AuditFrequency enum member set.
var auditFrequencies = map[string]bool{
	"DAILY": true, "WEEKLY": true, "BIWEEKLY": true, "MONTHLY": true,
}

// dayOfWeekValues is the DayOfWeek enum member set.
var dayOfWeekValues = map[string]bool{
	"SUN": true, "MON": true, "TUE": true, "WED": true, "THU": true, "FRI": true, "SAT": true,
}

// dayOfMonthPattern is the model's DayOfMonth pattern: "1" through "31"
// or the literal "LAST".
var dayOfMonthPattern = regexp.MustCompile(`^([1-9]|[12][0-9]|3[01])$|^LAST$`)

// validateScheduledAuditFrequency enforces the AuditFrequency enum and the
// conditional members the model documents: dayOfMonth is required for
// MONTHLY and dayOfWeek for WEEKLY/BIWEEKLY.
func validateScheduledAuditFrequency(frequency, dayOfMonth, dayOfWeek string) error {
	if !auditFrequencies[frequency] {
		return iotstore.ErrInvalidRequest
	}
	if dayOfWeek != "" && !dayOfWeekValues[dayOfWeek] {
		return iotstore.ErrInvalidRequest
	}
	if dayOfMonth != "" && !dayOfMonthPattern.MatchString(dayOfMonth) {
		return iotstore.ErrInvalidRequest
	}
	switch frequency {
	case "MONTHLY":
		if dayOfMonth == "" {
			return iotstore.ErrInvalidRequest
		}
	case "WEEKLY", "BIWEEKLY":
		if dayOfWeek == "" {
			return iotstore.ErrInvalidRequest
		}
	}
	return nil
}

// createScheduledAuditCore validates and persists a scheduled audit and
// returns its ARN.
func (s *IoTService) createScheduledAuditCore(store iotstore.IotStoreInterface, in ScheduledAuditInput) (string, error) {
	if in.Frequency == "" {
		return "", iotstore.ErrMissingParam
	}
	if len(in.TargetCheckNames) == 0 {
		return "", iotstore.ErrMissingParam
	}
	if err := validateScheduledAuditFrequency(in.Frequency, in.DayOfMonth, in.DayOfWeek); err != nil {
		return "", err
	}
	rec, err := s.bulkCreateCore(store, "scheduledAudit", in.Name, map[string]interface{}{
		"frequency":        in.Frequency,
		"dayOfMonth":       in.DayOfMonth,
		"dayOfWeek":        in.DayOfWeek,
		"targetCheckNames": in.TargetCheckNames,
	})
	if err != nil {
		return "", err
	}
	arn := iotstore.BuildScheduledAuditARN(store.GetAccountID(), store.GetRegion(), bulkName(rec))
	if len(in.Tags) > 0 {
		if err := store.TagResource(arn, in.Tags); err != nil {
			return "", err
		}
	}
	return arn, nil
}

// deleteScheduledAuditCore removes a scheduled audit record and its tags.
func (s *IoTService) deleteScheduledAuditCore(store iotstore.IotStoreInterface, name string) error {
	arn := iotstore.BuildScheduledAuditARN(store.GetAccountID(), store.GetRegion(), name)
	_ = store.DeleteAllTags(arn)
	return s.bulkDeleteCore(store, "scheduledAudit", name)
}

// ScheduledAuditRecord is the persisted scheduled audit record plus its ARN.
type ScheduledAuditRecord struct {
	Rec map[string]interface{}
	Arn string
}

// describeScheduledAuditCore loads a scheduled audit record and computes its
// ARN. An unknown name yields ErrScheduledAuditNotFound.
func (s *IoTService) describeScheduledAuditCore(store iotstore.IotStoreInterface, name string) (*ScheduledAuditRecord, error) {
	rec, exists, err := s.bulkGetCore(store, "scheduledAudit", name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrScheduledAuditNotFound
	}
	return &ScheduledAuditRecord{
		Rec: rec,
		Arn: iotstore.BuildScheduledAuditARN(store.GetAccountID(), store.GetRegion(), name),
	}, nil
}

// ScheduledAuditListItem is one ListScheduledAudits entry.
type ScheduledAuditListItem struct {
	Name      string
	Arn       string
	Frequency interface{}
}

// listScheduledAuditsCore lists every scheduled audit with its ARN.
func (s *IoTService) listScheduledAuditsCore(store iotstore.IotStoreInterface) ([]ScheduledAuditListItem, error) {
	items, err := s.bulkListCore(store, "scheduledAudit")
	if err != nil {
		return nil, err
	}
	out := make([]ScheduledAuditListItem, 0, len(items))
	for _, item := range items {
		name, _ := item["name"].(string)
		out = append(out, ScheduledAuditListItem{
			Name:      name,
			Arn:       iotstore.BuildScheduledAuditARN(store.GetAccountID(), store.GetRegion(), name),
			Frequency: item["frequency"],
		})
	}
	return out, nil
}

// updateScheduledAuditCore merges the supplied fields into an existing
// scheduled audit and returns its ARN. Only explicitly supplied members are
// applied; the conditional-member rules are validated against the merged
// state, so changing the frequency to MONTHLY keeps a stored dayOfMonth and
// changing it to WEEKLY without a dayOfWeek in the request or the record is
// rejected. An unknown name yields ErrScheduledAuditNotFound.
func (s *IoTService) updateScheduledAuditCore(store iotstore.IotStoreInterface, in ScheduledAuditInput) (string, error) {
	rec, exists, err := s.bulkGetCore(store, "scheduledAudit", in.Name)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", iotstore.ErrScheduledAuditNotFound
	}
	merge := map[string]interface{}{}
	if in.FrequencyProvided {
		merge["frequency"] = in.Frequency
	}
	if in.DayOfMonthProvided {
		merge["dayOfMonth"] = in.DayOfMonth
	}
	if in.DayOfWeekProvided {
		merge["dayOfWeek"] = in.DayOfWeek
	}
	if in.TargetCheckNamesProvided {
		merge["targetCheckNames"] = in.TargetCheckNames
	}
	if in.FrequencyProvided || in.DayOfMonthProvided || in.DayOfWeekProvided {
		frequency, _ := rec["frequency"].(string)
		if v, ok := merge["frequency"].(string); ok {
			frequency = v
		}
		dayOfMonth, _ := rec["dayOfMonth"].(string)
		if v, ok := merge["dayOfMonth"].(string); ok {
			dayOfMonth = v
		}
		dayOfWeek, _ := rec["dayOfWeek"].(string)
		if v, ok := merge["dayOfWeek"].(string); ok {
			dayOfWeek = v
		}
		if err := validateScheduledAuditFrequency(frequency, dayOfMonth, dayOfWeek); err != nil {
			return "", err
		}
	}
	updated, _, err := s.bulkUpdateCore(store, "scheduledAudit", in.Name, merge)
	if err != nil {
		return "", err
	}
	return iotstore.BuildScheduledAuditARN(store.GetAccountID(), store.GetRegion(), bulkName(updated)), nil
}

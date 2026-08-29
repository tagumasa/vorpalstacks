package eventbridge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
	awstypes "vorpalstacks/internal/common/tags"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// ---------------------------------------------------------------------------
// Input / Result structs (transport-agnostic)
// ---------------------------------------------------------------------------

// BusLogConfigInput carries the wire LogConfig members with per-member
// presence flags so an explicitly provided empty string — an out-of-enum
// value — stays distinguishable from an omitted member through the
// transport boundary.
type BusLogConfigInput struct {
	IncludeDetailSet bool
	IncludeDetail    string
	LevelSet         bool
	Level            string
}

// validateLogConfigInput rejects out-of-enum values for every provided
// member: the IncludeDetail and Level enums admit only NONE/FULL and
// OFF/ERROR/INFO/TRACE respectively, so an explicitly provided empty
// string is invalid exactly like any other non-member value.
func validateLogConfigInput(lc *BusLogConfigInput) error {
	if lc == nil {
		return nil
	}
	if lc.IncludeDetailSet && !isValidLogIncludeDetail(lc.IncludeDetail) {
		return awserrors.NewValidationException("LogConfig.IncludeDetail must be one of: NONE, FULL")
	}
	if lc.LevelSet && !isValidLogLevel(lc.Level) {
		return awserrors.NewValidationException("LogConfig.Level must be one of: OFF, ERROR, INFO, TRACE")
	}
	return nil
}

// storeLogConfig projects the input members onto the stored configuration;
// members absent from the request stay empty.
func storeLogConfig(lc *BusLogConfigInput) *eventsstore.BusLogConfig {
	if lc == nil {
		return nil
	}
	return &eventsstore.BusLogConfig{
		IncludeDetail: lc.IncludeDetail,
		Level:         lc.Level,
	}
}

// CreateEventBusInput carries the parameters for CreateEventBus in a
// transport-agnostic form shared by the HTTP API and the admin handler.
type CreateEventBusInput struct {
	Name             string
	Description      string
	Policy           string
	KmsKeyIdentifier string
	DeadLetterConfig *eventsstore.DeadLetterConfig
	LogConfig        *BusLogConfigInput
	Tags             []awstypes.Tag
}

// CreateEventBusResult holds the outcome of CreateEventBusCore.
type CreateEventBusResult struct {
	EventBus *eventsstore.EventBus
}

// DeleteEventBusInput carries the parameters for DeleteEventBus.
type DeleteEventBusInput struct {
	Name string
}

// ListEventBusesInput carries the parameters for ListEventBuses.
type ListEventBusesInput struct {
	NamePrefix string
	Limit      int32
	NextToken  string
}

// ListRulesInput carries the parameters for ListRules.
type ListRulesInput struct {
	EventBusName         string
	EventBusNameProvided bool
	NamePrefix           string
	Limit                int32
	NextToken            string
}

// resolveEventBusNameCore applies the EventBusName presence semantics: an
// absent member addresses the default event bus while an explicitly
// provided empty value is rejected per the Smithy @length(min=1) trait.
func resolveEventBusNameCore(name string, provided bool) (string, error) {
	if name != "" {
		return name, nil
	}
	if !provided {
		return "default", nil
	}
	return "", awserrors.NewValidationException("EventBusName must not be empty")
}

// ---------------------------------------------------------------------------
// Core functions
// ---------------------------------------------------------------------------

// createEventBusCore validates input, creates the event bus in the store and
// applies tags.  Shared by the HTTP API handler and the admin handler.
func (s *EventsService) createEventBusCore(ctx context.Context, store *eventsstore.EventsStore, input CreateEventBusInput) (*CreateEventBusResult, error) {
	if input.Name == "" {
		return nil, awserrors.NewValidationException("Event bus name is required")
	}
	if input.Name == "default" {
		return nil, awserrors.NewValidationException("Cannot create event bus named 'default'")
	}
	if !validateEventBusName(input.Name) {
		return nil, awserrors.NewValidationException("Event bus name must match the pattern and be 1-256 characters")
	}
	if input.KmsKeyIdentifier != "" && !validateKmsKeyIdentifier(input.KmsKeyIdentifier) {
		return nil, awserrors.NewValidationException("KmsKeyIdentifier must be a valid KMS ARN")
	}
	if err := validateLogConfigInput(input.LogConfig); err != nil {
		return nil, err
	}

	eventBus := &eventsstore.EventBus{
		Name:             input.Name,
		Description:      input.Description,
		Policy:           input.Policy,
		KmsKeyIdentifier: input.KmsKeyIdentifier,
		DeadLetterConfig: input.DeadLetterConfig,
		LogConfig:        storeLogConfig(input.LogConfig),
	}

	if err := store.CreateEventBus(ctx, eventBus); err != nil {
		return nil, mapStoreError(err, input.Name)
	}

	if len(input.Tags) > 0 {
		if err := store.TagStore.TagFromSlice(eventBus.ARN, input.Tags); err != nil {
			return nil, err
		}
	}

	return &CreateEventBusResult{EventBus: eventBus}, nil
}

// deleteEventBusCore validates input and performs a strict cascade-delete
// (rules → targets → archives → bus).  If any cascade step fails the bus
// is left in place and an InternalException is returned, mirroring the
// HTTP API contract.
func (s *EventsService) deleteEventBusCore(ctx context.Context, store *eventsstore.EventsStore, input DeleteEventBusInput) error {
	if input.Name == "" {
		return awserrors.NewValidationException("Event bus name is required")
	}
	if input.Name == "default" {
		return awserrors.NewValidationException("Cannot delete event bus 'default'")
	}

	if _, err := store.GetEventBus(ctx, input.Name); err != nil {
		return mapStoreError(err, input.Name)
	}

	var cascadeErr error

	ruleToken := ""
	for cascadeErr == nil {
		rulesResult, err := store.ListRules(ctx, input.Name, "", 1000, ruleToken)
		if err != nil {
			cascadeErr = fmt.Errorf("DeleteEventBus: list rules: %w", err)
			break
		}
		for _, rule := range rulesResult.Rules {
			targetToken := ""
			for cascadeErr == nil {
				targetsResult, tErr := store.ListTargetsByRule(ctx, input.Name, rule.Name, 1000, targetToken)
				if tErr != nil {
					cascadeErr = fmt.Errorf("DeleteEventBus: list targets for rule %s: %w", rule.Name, tErr)
					break
				}
				for _, t := range targetsResult.Targets {
					if err := store.DeleteTarget(ctx, input.Name, rule.Name, t.ID); err != nil {
						cascadeErr = fmt.Errorf("DeleteEventBus: delete target %s: %w", t.ID, err)
						break
					}
				}
				if cascadeErr != nil {
					break
				}
				if targetsResult.NextToken == "" {
					break
				}
				targetToken = targetsResult.NextToken
			}
			if cascadeErr != nil {
				break
			}
			if err := store.DeleteRule(ctx, input.Name, rule.Name); err != nil {
				cascadeErr = fmt.Errorf("DeleteEventBus: delete rule %s: %w", rule.Name, err)
				break
			}
			lastFireTimes.Delete(rule.ARN)
			_ = store.TagStore.Delete(rule.ARN)
		}
		if cascadeErr != nil {
			break
		}
		if rulesResult.NextToken == "" {
			break
		}
		ruleToken = rulesResult.NextToken
	}

	if cascadeErr == nil {
		archives, err := store.ListArchivesForEventBus(ctx, input.Name)
		if err != nil {
			cascadeErr = fmt.Errorf("DeleteEventBus: list archives: %w", err)
		} else {
			for _, a := range archives {
				if err := store.DeleteArchiveEvents(ctx, a.Name); err != nil {
					cascadeErr = fmt.Errorf("DeleteEventBus: delete archive events %s: %w", a.Name, err)
					break
				}
				if err := store.DeleteArchive(ctx, a.Name); err != nil {
					cascadeErr = fmt.Errorf("DeleteEventBus: delete archive %s: %w", a.Name, err)
					break
				}
			}
		}
	}

	if cascadeErr != nil {
		return awserrors.NewAWSError(
			"InternalException",
			cascadeErr.Error(),
			http.StatusInternalServerError,
		)
	}

	return mapStoreError(store.DeleteEventBus(ctx, input.Name), input.Name)
}

// listEventBusesCore validates input and delegates to the store.
func (s *EventsService) listEventBusesCore(ctx context.Context, store *eventsstore.EventsStore, input ListEventBusesInput) (*eventsstore.EventBusListResult, error) {
	limit := input.Limit
	if limit == 0 {
		limit = 100
	}
	if limit < 1 || limit > 100 {
		return nil, awserrors.NewValidationException("Limit must be between 1 and 100")
	}
	return store.ListEventBuses(ctx, input.NamePrefix, limit, input.NextToken)
}

// listRulesCore validates input and delegates to the store. The event bus
// name is resolved first: this operation carries no other required member,
// so the empty-bus-name rejection takes precedence over the limit window.
func (s *EventsService) listRulesCore(ctx context.Context, store *eventsstore.EventsStore, input ListRulesInput) (*eventsstore.RuleListResult, error) {
	eventBusName, err := resolveEventBusNameCore(input.EventBusName, input.EventBusNameProvided)
	if err != nil {
		return nil, err
	}
	limit := input.Limit
	if limit == 0 {
		limit = 100
	}
	if limit < 1 || limit > 100 {
		return nil, awserrors.NewValidationException("Limit must be between 1 and 100")
	}
	return store.ListRules(ctx, eventBusName, input.NamePrefix, limit, input.NextToken)
}

// DescribeEventBusResult holds the outcome of describeEventBusCore: the
// event bus record plus its tags.
type DescribeEventBusResult struct {
	EventBus *eventsstore.EventBus
	Tags     []awstypes.Tag
}

// UpdateEventBusInput carries the parameters for UpdateEventBus. The *Set
// flags distinguish an omitted member from an explicitly provided empty one
// so the merge semantics survive the transport boundary.
type UpdateEventBusInput struct {
	Name                string
	DescriptionSet      bool
	Description         string
	PolicySet           bool
	Policy              string
	KmsKeyIdentifierSet bool
	KmsKeyIdentifier    string
	DeadLetterConfigSet bool
	DeadLetterConfig    *eventsstore.DeadLetterConfig
	LogConfigSet        bool
	LogConfig           *BusLogConfigInput
}

// PutPermissionInput carries the parameters for PutPermission. The Policy
// members select the full-policy-document mode; the individual members
// select the statement mode.
type PutPermissionInput struct {
	BusName         string
	BusNameProvided bool
	PolicySet       bool
	Policy          string
	Principal       string
	StatementId     string
	Action          string
	Condition       string
}

// RemovePermissionInput carries the parameters for RemovePermission.
type RemovePermissionInput struct {
	BusName         string
	BusNameProvided bool
	StatementId     string
	RemoveAll       bool
}

// describeEventBusCore validates input and fetches the event bus with its
// tags.
func (s *EventsService) describeEventBusCore(ctx context.Context, store *eventsstore.EventsStore, name string) (*DescribeEventBusResult, error) {
	eventBus, err := store.GetEventBus(ctx, name)
	if err != nil {
		return nil, mapStoreError(err, name)
	}

	result := &DescribeEventBusResult{EventBus: eventBus}
	if tagSlice, err := store.TagStore.ListAsSlice(eventBus.ARN); err == nil && len(tagSlice) > 0 {
		result.Tags = tagSlice
	}
	return result, nil
}

// updateEventBusCore validates input, merges the provided members onto the
// stored event bus and persists the update.
func (s *EventsService) updateEventBusCore(ctx context.Context, store *eventsstore.EventsStore, input UpdateEventBusInput) (*eventsstore.EventBus, error) {
	if input.Name == "" {
		return nil, awserrors.NewValidationException("Event bus name is required")
	}

	eventBus, err := store.GetEventBus(ctx, input.Name)
	if err != nil {
		return nil, mapStoreError(err, input.Name)
	}

	if input.DescriptionSet {
		if !validateDescription(input.Description) {
			return nil, errDescriptionTooLong()
		}
		eventBus.Description = input.Description
	}

	if input.PolicySet {
		eventBus.Policy = input.Policy
	}

	if input.KmsKeyIdentifierSet {
		eventBus.KmsKeyIdentifier = input.KmsKeyIdentifier
	}
	if input.DeadLetterConfigSet {
		eventBus.DeadLetterConfig = input.DeadLetterConfig
	}
	if input.LogConfigSet {
		if err := validateLogConfigInput(input.LogConfig); err != nil {
			return nil, err
		}
		eventBus.LogConfig = storeLogConfig(input.LogConfig)
	}

	if err := store.UpdateEventBus(ctx, eventBus); err != nil {
		return nil, err
	}
	return eventBus, nil
}

// putPermissionCore validates input and adds a resource policy statement to
// the event bus, granting the given principal permission to put events.
// Supports two modes: (1) individual parameters (Principal, StatementId,
// Action, Condition) and (2) a complete policy document via the Policy
// member.
func (s *EventsService) putPermissionCore(ctx context.Context, store *eventsstore.EventsStore, input PutPermissionInput) error {
	busName, err := resolveEventBusNameCore(input.BusName, input.BusNameProvided)
	if err != nil {
		return err
	}
	eventBus, err := store.GetEventBus(ctx, busName)
	if err != nil {
		return mapStoreError(err, busName)
	}

	// Mode 1: Full policy document provided via the Policy member.
	if input.PolicySet && input.Policy != "" {
		if err := validateEventBusPolicySize(input.Policy); err != nil {
			return err
		}
		var policyDoc map[string]interface{}
		if err := json.Unmarshal([]byte(input.Policy), &policyDoc); err != nil {
			return awserrors.NewValidationException("Invalid policy document")
		}
		if _, ok := policyDoc["Version"]; !ok {
			policyDoc["Version"] = "2012-10-17"
		}
		policyBytes, err := json.Marshal(policyDoc)
		if err != nil {
			return fmt.Errorf("failed to marshal policy: %w", err)
		}
		eventBus.Policy = string(policyBytes)
		return store.UpdateEventBus(ctx, eventBus)
	}

	// Mode 2: Individual parameters (Principal, StatementId, Action, Condition).
	principal := input.Principal
	statementID := input.StatementId
	action := input.Action
	if action == "" {
		action = "events:PutEvents"
	}

	if principal == "" || statementID == "" {
		return awserrors.NewValidationException("Principal and StatementId are required")
	}

	var policyDoc map[string]interface{}
	if eventBus.Policy != "" {
		if err := json.Unmarshal([]byte(eventBus.Policy), &policyDoc); err != nil {
			policyDoc = make(map[string]interface{})
		}
	}
	if _, ok := policyDoc["Version"]; !ok {
		policyDoc["Version"] = "2012-10-17"
	}

	statement := map[string]interface{}{
		"Sid":       statementID,
		"Effect":    "Allow",
		"Principal": map[string]interface{}{"AWS": principal},
		"Action":    action,
		"Resource":  eventBus.ARN,
	}
	if input.Condition != "" {
		var cond map[string]interface{}
		if err := json.Unmarshal([]byte(input.Condition), &cond); err != nil {
			return awserrors.NewValidationException("Condition must be valid JSON: " + err.Error())
		}
		statement["Condition"] = cond
	}

	statements, _ := policyDoc["Statement"].([]interface{})
	replaced := false
	for i, s := range statements {
		if stmt, ok := s.(map[string]interface{}); ok {
			if sid, _ := stmt["Sid"].(string); sid == statementID {
				statements[i] = statement
				replaced = true
				break
			}
		}
	}
	if !replaced {
		statements = append(statements, statement)
	}
	policyDoc["Statement"] = statements

	policyBytes, err := json.Marshal(policyDoc)
	if err != nil {
		return fmt.Errorf("failed to marshal policy: %w", err)
	}
	eventBus.Policy = string(policyBytes)

	return store.UpdateEventBus(ctx, eventBus)
}

// removePermissionCore validates input and removes a resource policy
// statement from the event bus identified by its StatementId, or the whole
// policy when RemoveAll is set.
func (s *EventsService) removePermissionCore(ctx context.Context, store *eventsstore.EventsStore, input RemovePermissionInput) error {
	busName, err := resolveEventBusNameCore(input.BusName, input.BusNameProvided)
	if err != nil {
		return err
	}

	// AWS requires either StatementId or RemoveAllPermissions=true.
	if !input.RemoveAll && input.StatementId == "" {
		return awserrors.NewValidationException("StatementId is required when RemoveAllPermissions is not true")
	}

	eventBus, err := store.GetEventBus(ctx, busName)
	if err != nil {
		return mapStoreError(err, busName)
	}

	if eventBus.Policy == "" {
		return nil
	}

	// RemoveAllPermissions clears the policy entirely.
	if input.RemoveAll {
		eventBus.Policy = ""
		return store.UpdateEventBus(ctx, eventBus)
	}

	var policyDoc map[string]interface{}
	if err := json.Unmarshal([]byte(eventBus.Policy), &policyDoc); err != nil {
		return nil
	}

	statements, ok := policyDoc["Statement"].([]interface{})
	if !ok {
		return nil
	}

	filtered := make([]interface{}, 0, len(statements))
	for _, s := range statements {
		if stmt, ok := s.(map[string]interface{}); ok {
			if sid, _ := stmt["Sid"].(string); sid != input.StatementId {
				filtered = append(filtered, s)
			}
		}
	}

	if len(filtered) == 0 {
		delete(policyDoc, "Statement")
	} else {
		policyDoc["Statement"] = filtered
	}

	policyBytes, err := json.Marshal(policyDoc)
	if err != nil {
		return fmt.Errorf("failed to marshal policy: %w", err)
	}
	eventBus.Policy = string(policyBytes)

	return store.UpdateEventBus(ctx, eventBus)
}

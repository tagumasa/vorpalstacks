package eventbridge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
// Policy is not a member: neither create request shape models it, and the
// bus policy's write surface is PutPermission/RemovePermission alone.
type CreateEventBusInput struct {
	Name             string
	Description      string
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
// absent member addresses the default event bus while an explicitly provided
// empty value is rejected per the Smithy @length(min=1) trait. The member
// type is EventBusNameOrArn — "The name or ARN of the event bus" (model
// member documentation) — so the ARN form is normalised to the canonical bus
// name through the same extractor the PutEvents ingress plane uses, keeping
// rule/archive keys name-addressed on every plane.
func resolveEventBusNameCore(name string, provided bool) (string, error) {
	if name != "" {
		return resolveEventBusName(name), nil
	}
	if !provided {
		return "default", nil
	}
	return "", awserrors.NewValidationException("EventBusName must not be empty")
}

// resolveNonPartnerEventBusNameCore resolves the permission family's bus
// reference: the member targets NonPartnerEventBusName (pattern
// ^[\.\-_A-Za-z0-9]+$, 1-256), whose character class excludes both the ARN
// form and the slash-bearing partner event-source names, so any other form
// is a pattern violation before the shared resolver's default-bus fallback
// applies.
func resolveNonPartnerEventBusNameCore(name string, provided bool) (string, error) {
	if name != "" && !validateNonPartnerEventBusName(name) {
		return "", awserrors.NewValidationException("Event bus name must match the pattern and be 1-256 characters")
	}
	return resolveEventBusNameCore(name, provided)
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
	if !validateDescription(input.Description) {
		return nil, errDescriptionTooLong()
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
	// Name targets the plain EventBusName shape, so the ARN form is a
	// pattern violation here (unlike the EventBusNameOrArn rule family).
	if !validateEventBusName(input.Name) {
		return awserrors.NewValidationException("Event bus name must match the pattern and be 1-256 characters")
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
			s.fireDedup.deleteLastFire(rule.ARN)
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
	if input.NamePrefix != "" {
		if err := validateListNamePrefix(input.NamePrefix, "event-bus"); err != nil {
			return nil, err
		}
	}
	limit, err := normaliseListLimit(input.Limit)
	if err != nil {
		return nil, err
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
	if input.NamePrefix != "" {
		if err := validateListNamePrefix(input.NamePrefix, "rule"); err != nil {
			return nil, err
		}
	}
	limit, err := normaliseListLimit(input.Limit)
	if err != nil {
		return nil, err
	}
	return store.ListRules(ctx, eventBusName, input.NamePrefix, limit, input.NextToken)
}

// DescribeEventBusResult holds the outcome of describeEventBusCore: the
// event bus record. Tags are not a DescribeEventBusResponse member — the
// tag read surface is ListTagsForResource.
type DescribeEventBusResult struct {
	EventBus *eventsstore.EventBus
}

// UpdateEventBusInput carries the parameters for UpdateEventBus. The *Set
// flags distinguish an omitted member from an explicitly provided empty one
// so the merge semantics survive the transport boundary; NameProvided does
// the same for the optional Name member. Policy is not a member: the update
// request shape does not model it, and the bus policy's write surface is
// PutPermission/RemovePermission alone.
type UpdateEventBusInput struct {
	Name                string
	NameProvided        bool
	DescriptionSet      bool
	Description         string
	KmsKeyIdentifierSet bool
	KmsKeyIdentifier    string
	DeadLetterConfigSet bool
	DeadLetterConfig    *eventsstore.DeadLetterConfig
	LogConfigSet        bool
	LogConfig           *BusLogConfigInput
}

// PutPermissionCondition carries the typed Condition member of PutPermission:
// the modelled structure with Key, Type and Value ("The Condition is a JSON
// string which must contain Type, Key, and Value fields" — the members are
// documented as the single supported pair Type=StringEquals,
// Key=aws:PrincipalOrgID).
type PutPermissionCondition struct {
	Key   string
	Type  string
	Value string
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
	Condition       *PutPermissionCondition
}

// RemovePermissionInput carries the parameters for RemovePermission.
type RemovePermissionInput struct {
	BusName         string
	BusNameProvided bool
	StatementId     string
	RemoveAll       bool
}

// describeEventBusCore resolves the bus name ("The name or ARN of the
// event bus to show details for. If you omit this, the default event bus
// is displayed" — model member documentation) and fetches the event bus
// with its tags.
func (s *EventsService) describeEventBusCore(ctx context.Context, store *eventsstore.EventsStore, name string, provided bool) (*DescribeEventBusResult, error) {
	eventBusName, err := resolveEventBusNameCore(name, provided)
	if err != nil {
		return nil, err
	}
	eventBus, err := store.GetEventBus(ctx, eventBusName)
	if err != nil {
		return nil, mapStoreError(err, eventBusName)
	}

	result := &DescribeEventBusResult{EventBus: eventBus}
	return result, nil
}

// updateEventBusCore validates input, merges the provided members onto the
// stored event bus and persists the update through the atomic record
// mutation.
func (s *EventsService) updateEventBusCore(ctx context.Context, store *eventsstore.EventsStore, input UpdateEventBusInput) (*eventsstore.EventBus, error) {
	// Name is optional on the update request (no required trait on the
	// model member; API reference "Required: No") and follows the family
	// semantics of every optional bus-name member: an omitted name
	// addresses the default event bus, an explicitly provided empty value
	// violates the EventBusName length floor. The member targets the plain
	// EventBusName form, so no ARN canonicalisation applies and the ARN
	// form is a pattern violation.
	eventBusName := input.Name
	if eventBusName == "" {
		if !input.NameProvided {
			eventBusName = "default"
		} else {
			return nil, awserrors.NewValidationException("Event bus name must not be empty")
		}
	} else if !validateEventBusName(eventBusName) {
		return nil, awserrors.NewValidationException("Event bus name must match the pattern and be 1-256 characters")
	}

	var updated *eventsstore.EventBus
	if err := store.MutateEventBus(ctx, eventBusName, func(eventBus *eventsstore.EventBus) error {
		if input.DescriptionSet {
			if !validateDescription(input.Description) {
				return errDescriptionTooLong()
			}
			eventBus.Description = input.Description
		}

		if input.KmsKeyIdentifierSet {
			if input.KmsKeyIdentifier != "" && !validateKmsKeyIdentifier(input.KmsKeyIdentifier) {
				return awserrors.NewValidationException("KmsKeyIdentifier must be a valid KMS ARN")
			}
			eventBus.KmsKeyIdentifier = input.KmsKeyIdentifier
		}
		if input.DeadLetterConfigSet {
			eventBus.DeadLetterConfig = input.DeadLetterConfig
		}
		if input.LogConfigSet {
			if err := validateLogConfigInput(input.LogConfig); err != nil {
				return err
			}
			eventBus.LogConfig = storeLogConfig(input.LogConfig)
		}

		eventBus.LastModifiedAt = time.Now().UTC()
		updated = eventBus
		return nil
	}); err != nil {
		return nil, mapStoreError(err, eventBusName)
	}
	return updated, nil
}

// putPermissionCore validates input and adds a resource policy statement to
// the event bus, granting the given principal permission to put events.
// Supports two modes: (1) individual parameters (Principal, StatementId,
// Action, Condition) and (2) a complete policy document via the Policy
// member.
func (s *EventsService) putPermissionCore(ctx context.Context, store *eventsstore.EventsStore, input PutPermissionInput) error {
	busName, err := resolveNonPartnerEventBusNameCore(input.BusName, input.BusNameProvided)
	if err != nil {
		return err
	}
	if _, err := store.GetEventBus(ctx, busName); err != nil {
		return mapStoreError(err, busName)
	}

	// Mode 1: Full policy document provided via the Policy member. The
	// member's presence selects this mode — "You can include a Policy
	// parameter in the request instead of using the StatementId, Action,
	// Principal, or Condition parameters" — so an explicitly empty Policy
	// is an invalid policy document rather than a fall-through to the
	// statement mode. The write runs through the atomic record mutation so
	// it cannot interleave with a concurrent statement-mode merge.
	if input.PolicySet {
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
		return store.MutateEventBus(ctx, busName, func(eventBus *eventsstore.EventBus) error {
			eventBus.Policy = string(policyBytes)
			eventBus.LastModifiedAt = time.Now().UTC()
			return nil
		})
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
	if !putPermissionPrincipalPattern.MatchString(principal) {
		return awserrors.NewValidationException("Principal must be a 12-digit account ID or \"*\"")
	}
	if len(statementID) > maxPutPermissionStatementIDLength || !putPermissionStatementIDPattern.MatchString(statementID) {
		return awserrors.NewValidationException("StatementId must match [a-zA-Z0-9-_] and be 1-64 characters")
	}
	if len(action) > maxPutPermissionActionLength || !putPermissionActionPattern.MatchString(action) {
		return awserrors.NewValidationException("Action must match events:[a-zA-Z] and be 1-64 characters")
	}
	if input.Condition != nil {
		c := input.Condition
		if c.Type == "" || c.Key == "" || c.Value == "" {
			return awserrors.NewValidationException("Condition must contain Type, Key, and Value")
		}
		if c.Type != "StringEquals" {
			return awserrors.NewValidationException("Condition.Type supports only StringEquals")
		}
		if c.Key != "aws:PrincipalOrgID" {
			return awserrors.NewValidationException("Condition.Key supports only aws:PrincipalOrgID")
		}
	}

	// The Sid merge reads and writes the current policy inside the atomic
	// record mutation: two concurrent statement-mode merges cannot lose
	// each other's statements.
	return store.MutateEventBus(ctx, busName, func(eventBus *eventsstore.EventBus) error {
		// The map is allocated up front: a bus whose Policy member is
		// still empty must take its first statement without a nil-map
		// assignment panic.
		policyDoc := make(map[string]interface{})
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
		if input.Condition != nil {
			// The resource-policy Condition element nests the operator
			// under the condition type and the organisation ID under the
			// key: {"StringEquals": {"aws:PrincipalOrgID": "o-..."}}.
			statement["Condition"] = map[string]interface{}{
				input.Condition.Type: map[string]interface{}{
					input.Condition.Key: input.Condition.Value,
				},
			}
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
		// The 10 KB ceiling governs the merged document — repeated
		// statement-mode writes must not grow the policy past the size a
		// full-document mode write could ever set. Returning the error
		// from inside the mutation aborts the merge atomically.
		if err := validateEventBusPolicySize(string(policyBytes)); err != nil {
			return err
		}
		eventBus.Policy = string(policyBytes)
		eventBus.LastModifiedAt = time.Now().UTC()
		return nil
	})
}

// removePermissionCore validates input and removes a resource policy
// statement from the event bus identified by its StatementId, or the whole
// policy when RemoveAll is set.
func (s *EventsService) removePermissionCore(ctx context.Context, store *eventsstore.EventsStore, input RemovePermissionInput) error {
	busName, err := resolveNonPartnerEventBusNameCore(input.BusName, input.BusNameProvided)
	if err != nil {
		return err
	}

	// AWS requires either StatementId or RemoveAllPermissions=true.
	if !input.RemoveAll && input.StatementId == "" {
		return awserrors.NewValidationException("StatementId is required when RemoveAllPermissions is not true")
	}

	if _, err := store.GetEventBus(ctx, busName); err != nil {
		return mapStoreError(err, busName)
	}

	// The statement removal reads and writes the current policy inside the
	// atomic record mutation, mirroring the statement-mode merge.
	return store.MutateEventBus(ctx, busName, func(eventBus *eventsstore.EventBus) error {
		if eventBus.Policy == "" {
			return nil
		}

		// RemoveAllPermissions clears the policy entirely.
		if input.RemoveAll {
			eventBus.Policy = ""
			eventBus.LastModifiedAt = time.Now().UTC()
			return nil
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
		eventBus.LastModifiedAt = time.Now().UTC()
		return nil
	})
}

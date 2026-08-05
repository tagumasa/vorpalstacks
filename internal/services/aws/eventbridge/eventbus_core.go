package eventbridge

import (
	"context"
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	awstypes "vorpalstacks/internal/utils/aws/types"
)

// ---------------------------------------------------------------------------
// Input / Result structs (transport-agnostic)
// ---------------------------------------------------------------------------

// CreateEventBusInput carries the parameters for CreateEventBus in a
// transport-agnostic form shared by the HTTP API and the admin handler.
type CreateEventBusInput struct {
	Name             string
	Description      string
	Policy           string
	KmsKeyIdentifier string
	DeadLetterConfig *eventsstore.DeadLetterConfig
	LogConfig        *eventsstore.BusLogConfig
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
	EventBusName string
	NamePrefix   string
	Limit        int32
	NextToken    string
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
	if input.LogConfig != nil {
		if input.LogConfig.IncludeDetail != "" && !isValidLogIncludeDetail(input.LogConfig.IncludeDetail) {
			return nil, awserrors.NewValidationException("LogConfig.IncludeDetail must be one of: NONE, FULL")
		}
		if input.LogConfig.Level != "" && !isValidLogLevel(input.LogConfig.Level) {
			return nil, awserrors.NewValidationException("LogConfig.Level must be one of: OFF, ERROR, INFO, TRACE")
		}
	}

	eventBus := &eventsstore.EventBus{
		Name:             input.Name,
		Description:      input.Description,
		Policy:           input.Policy,
		KmsKeyIdentifier: input.KmsKeyIdentifier,
		DeadLetterConfig: input.DeadLetterConfig,
		LogConfig:        input.LogConfig,
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

// listRulesCore validates input and delegates to the store.
func (s *EventsService) listRulesCore(ctx context.Context, store *eventsstore.EventsStore, input ListRulesInput) (*eventsstore.RuleListResult, error) {
	eventBusName := input.EventBusName
	if eventBusName == "" {
		eventBusName = "default"
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

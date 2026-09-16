package eventbridge

import (
	"context"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ---------------------------------------------------------------------------
// Input / Result structs (transport-agnostic)
// ---------------------------------------------------------------------------

// CreateApiDestinationInput carries the parameters for CreateApiDestination.
// The *Set flags distinguish an omitted member from an explicitly provided
// one so the merge semantics survive the transport boundary.
type CreateApiDestinationInput struct {
	Name                   string
	ConnectionArn          string
	HttpMethod             string
	InvocationEndpoint     string
	DescriptionSet         bool
	Description            string
	InvocationRateLimitSet bool
	InvocationRateLimit    int32
}

// UpdateApiDestinationInput carries the parameters for UpdateApiDestination.
type UpdateApiDestinationInput struct {
	Name                   string
	DescriptionSet         bool
	Description            string
	HttpMethodSet          bool
	HttpMethod             string
	InvocationEndpointSet  bool
	InvocationEndpoint     string
	ConnectionArnSet       bool
	ConnectionArn          string
	InvocationRateLimitSet bool
	InvocationRateLimit    int32
}

// ListApiDestinationsInput carries the parameters for ListApiDestinations.
type ListApiDestinationsInput struct {
	NamePrefix    string
	ConnectionArn string
	Limit         int32
	NextToken     string
}

// ---------------------------------------------------------------------------
// Core functions
// ---------------------------------------------------------------------------

// createApiDestinationCore validates input and creates the API destination.
func (s *EventsService) createApiDestinationCore(ctx context.Context, store *eventsstore.EventsStore, input CreateApiDestinationInput) (*eventsstore.ApiDestination, error) {
	if input.Name == "" {
		return nil, awserrors.NewValidationException("Api destination name is required")
	}
	if !validateResourceName(input.Name, "api-destination") {
		return nil, awserrors.NewValidationException("Api destination name must match the pattern and be 1-64 characters")
	}

	if err := s.validateConnectionArn(ctx, store, input.ConnectionArn); err != nil {
		return nil, err
	}

	// HttpMethod is a required member of CreateApiDestinationRequest (no
	// documented default — an omitted method is a ValidationException).
	if input.HttpMethod == "" {
		return nil, awserrors.NewValidationException("HttpMethod is required")
	}
	if !validHttpMethods[input.HttpMethod] {
		return nil, awserrors.NewValidationException("HttpMethod must be one of: GET, POST, PUT, DELETE, HEAD, OPTIONS, PATCH")
	}

	if input.InvocationEndpoint == "" {
		return nil, awserrors.NewValidationException("InvocationEndpoint is required")
	}
	if !validateInvocationEndpoint(input.InvocationEndpoint) {
		return nil, awserrors.NewValidationException("InvocationEndpoint must be an https URL matching the HttpsEndpoint pattern and at most 2048 characters")
	}

	apiDest := &eventsstore.ApiDestination{
		Name:               input.Name,
		ConnectionARN:      input.ConnectionArn,
		HttpMethod:         input.HttpMethod,
		InvocationEndpoint: input.InvocationEndpoint,
	}

	if input.DescriptionSet {
		if !validateDescription(input.Description) {
			return nil, errDescriptionTooLong()
		}
		apiDest.Description = input.Description
	}

	if input.InvocationRateLimitSet {
		if !validateInvocationRateLimit(input.InvocationRateLimit) {
			return nil, awserrors.NewValidationException("InvocationRateLimitPerSecond must be at least 1")
		}
		apiDest.InvocationRateLimitPerSecond = input.InvocationRateLimit
	}

	if err := store.CreateApiDestination(ctx, apiDest); err != nil {
		return nil, mapStoreError(err, input.Name)
	}
	return apiDest, nil
}

// deleteApiDestinationCore validates input and deletes the API destination.
func (s *EventsService) deleteApiDestinationCore(ctx context.Context, store *eventsstore.EventsStore, name string) error {
	if name == "" {
		return awserrors.NewValidationException("Api destination name is required")
	}
	return mapStoreError(store.DeleteApiDestination(ctx, name), name)
}

// getApiDestinationCore validates input and fetches the API destination.
func (s *EventsService) getApiDestinationCore(ctx context.Context, store *eventsstore.EventsStore, name string) (*eventsstore.ApiDestination, error) {
	if name == "" {
		return nil, awserrors.NewValidationException("Api destination name is required")
	}
	apiDest, err := store.GetApiDestination(ctx, name)
	if err != nil {
		return nil, mapStoreError(err, name)
	}
	return apiDest, nil
}

// updateApiDestinationCore validates input, merges the provided members onto
// the stored record and persists the update through the atomic record
// mutation.
func (s *EventsService) updateApiDestinationCore(ctx context.Context, store *eventsstore.EventsStore, input UpdateApiDestinationInput) (*eventsstore.ApiDestination, error) {
	if input.Name == "" {
		return nil, awserrors.NewValidationException("Api destination name is required")
	}
	if input.ConnectionArnSet && input.ConnectionArn != "" {
		if err := s.validateConnectionArn(ctx, store, input.ConnectionArn); err != nil {
			return nil, err
		}
	}

	var updated *eventsstore.ApiDestination
	if err := store.MutateApiDestination(ctx, input.Name, func(apiDest *eventsstore.ApiDestination) error {
		if input.DescriptionSet {
			if !validateDescription(input.Description) {
				return errDescriptionTooLong()
			}
			apiDest.Description = input.Description
		}
		if input.HttpMethodSet && input.HttpMethod != "" {
			if !validHttpMethods[input.HttpMethod] {
				return awserrors.NewValidationException("HttpMethod must be one of: GET, POST, PUT, DELETE, HEAD, OPTIONS, PATCH")
			}
			apiDest.HttpMethod = input.HttpMethod
		}
		if input.InvocationEndpointSet && input.InvocationEndpoint != "" {
			if !validateInvocationEndpoint(input.InvocationEndpoint) {
				return awserrors.NewValidationException("InvocationEndpoint must be an https URL matching the HttpsEndpoint pattern and at most 2048 characters")
			}
			apiDest.InvocationEndpoint = input.InvocationEndpoint
		}
		if input.ConnectionArnSet && input.ConnectionArn != "" {
			apiDest.ConnectionARN = input.ConnectionArn
		}
		if input.InvocationRateLimitSet {
			if !validateInvocationRateLimit(input.InvocationRateLimit) {
				return awserrors.NewValidationException("InvocationRateLimitPerSecond must be at least 1")
			}
			apiDest.InvocationRateLimitPerSecond = input.InvocationRateLimit
		}

		apiDest.LastModifiedAt = time.Now().UTC()
		updated = apiDest
		return nil
	}); err != nil {
		return nil, mapStoreError(err, input.Name)
	}
	return updated, nil
}

// validateConnectionArn enforces the ConnectionArn contract on
// CreateApiDestination and UpdateApiDestination: the value must match the
// Smithy @pattern (arn:...:events:...:connection/<name>/<id>) and must name
// an existing connection — the destination's invocations draw their
// credentials from that connection.
func (s *EventsService) validateConnectionArn(ctx context.Context, store *eventsstore.EventsStore, arn string) error {
	if arn == "" {
		return awserrors.NewValidationException("ConnectionArn is required")
	}
	if !connectionArnPattern.MatchString(arn) {
		return awserrors.NewValidationException("ConnectionArn " + arn + " does not match the required pattern arn:<partition>:events:<region>:<account>:connection/<name>/<id>")
	}
	connectionName := svcarn.ExtractConnectionNameFromARN(arn)
	if _, err := store.GetConnection(ctx, connectionName); err != nil {
		if err == eventsstore.ErrConnectionNotFound {
			return NewResourceNotFoundException("Connection '" + connectionName + "' does not exist")
		}
		return err
	}
	return nil
}

// listApiDestinationsCore validates the query (limit window and NamePrefix
// charset, both per the model) and lists the API destinations.
func (s *EventsService) listApiDestinationsCore(ctx context.Context, store *eventsstore.EventsStore, input ListApiDestinationsInput) (*eventsstore.ApiDestinationListResult, error) {
	if input.NamePrefix != "" {
		if err := validateListNamePrefix(input.NamePrefix, "api-destination"); err != nil {
			return nil, err
		}
	}
	limit, err := normaliseListLimit(input.Limit)
	if err != nil {
		return nil, err
	}
	return store.ListApiDestinations(ctx, input.NamePrefix, input.ConnectionArn, limit, input.NextToken)
}

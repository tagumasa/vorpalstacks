package eventbridge

import (
	"context"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// ---------------------------------------------------------------------------
// Input / Result structs (transport-agnostic)
// ---------------------------------------------------------------------------

// CreateApiDestinationInput carries the parameters for CreateApiDestination.
// The *Set flags distinguish an omitted member from an explicitly provided
// empty one so the merge semantics survive the transport boundary.
type CreateApiDestinationInput struct {
	Name                string
	ConnectionArn       string
	HttpMethod          string
	InvocationEndpoint  string
	DescriptionSet      bool
	Description         string
	InvocationRateLimit int32
}

// UpdateApiDestinationInput carries the parameters for UpdateApiDestination.
type UpdateApiDestinationInput struct {
	Name                  string
	DescriptionSet        bool
	Description           string
	HttpMethodSet         bool
	HttpMethod            string
	InvocationEndpointSet bool
	InvocationEndpoint    string
	ConnectionArnSet      bool
	ConnectionArn         string
	InvocationRateLimit   int32
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

	if input.ConnectionArn == "" {
		return nil, awserrors.NewValidationException("ConnectionArn is required")
	}

	httpMethod := input.HttpMethod
	if httpMethod == "" {
		httpMethod = "POST"
	}
	if !validHttpMethods[httpMethod] {
		return nil, awserrors.NewValidationException("HttpMethod must be one of: GET, POST, PUT, DELETE, HEAD, OPTIONS, PATCH")
	}

	if input.InvocationEndpoint == "" {
		return nil, awserrors.NewValidationException("InvocationEndpoint is required")
	}

	apiDest := &eventsstore.ApiDestination{
		Name:               input.Name,
		ConnectionARN:      input.ConnectionArn,
		HttpMethod:         httpMethod,
		InvocationEndpoint: input.InvocationEndpoint,
	}

	if input.DescriptionSet {
		if !validateDescription(input.Description) {
			return nil, errDescriptionTooLong()
		}
		apiDest.Description = input.Description
	}

	if input.InvocationRateLimit > 0 {
		if !validateInvocationRateLimit(input.InvocationRateLimit) {
			return nil, awserrors.NewValidationException("InvocationRateLimitPerSecond must be between 1 and 300")
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
// the stored record and persists the update.
func (s *EventsService) updateApiDestinationCore(ctx context.Context, store *eventsstore.EventsStore, input UpdateApiDestinationInput) (*eventsstore.ApiDestination, error) {
	if input.Name == "" {
		return nil, awserrors.NewValidationException("Api destination name is required")
	}

	apiDest, err := store.GetApiDestination(ctx, input.Name)
	if err != nil {
		return nil, mapStoreError(err, input.Name)
	}

	if input.DescriptionSet {
		if !validateDescription(input.Description) {
			return nil, errDescriptionTooLong()
		}
		apiDest.Description = input.Description
	}
	if input.HttpMethodSet && input.HttpMethod != "" {
		if !validHttpMethods[input.HttpMethod] {
			return nil, awserrors.NewValidationException("HttpMethod must be one of: GET, POST, PUT, DELETE, HEAD, OPTIONS, PATCH")
		}
		apiDest.HttpMethod = input.HttpMethod
	}
	if input.InvocationEndpointSet && input.InvocationEndpoint != "" {
		apiDest.InvocationEndpoint = input.InvocationEndpoint
	}
	if input.ConnectionArnSet && input.ConnectionArn != "" {
		apiDest.ConnectionARN = input.ConnectionArn
	}
	if input.InvocationRateLimit > 0 {
		if !validateInvocationRateLimit(input.InvocationRateLimit) {
			return nil, awserrors.NewValidationException("InvocationRateLimitPerSecond must be between 1 and 300")
		}
		apiDest.InvocationRateLimitPerSecond = input.InvocationRateLimit
	}

	apiDest.LastModifiedAt = time.Now().UTC()

	if err := store.UpdateApiDestination(ctx, apiDest); err != nil {
		return nil, err
	}
	return apiDest, nil
}

// listApiDestinationsCore applies the documented limit window and lists the
// API destinations.
func (s *EventsService) listApiDestinationsCore(ctx context.Context, store *eventsstore.EventsStore, input ListApiDestinationsInput) (*eventsstore.ApiDestinationListResult, error) {
	limit := input.Limit
	if limit < 0 || limit > 100 {
		return nil, awserrors.NewValidationException("Limit must be between 0 and 100")
	}
	if limit == 0 {
		limit = 50
	}
	return store.ListApiDestinations(ctx, input.NamePrefix, input.ConnectionArn, limit, input.NextToken)
}

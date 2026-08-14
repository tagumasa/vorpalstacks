package eventbridge

import (
	"context"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

func apiDestinationToMap(d *eventsstore.ApiDestination) map[string]interface{} {
	result := map[string]interface{}{
		"ApiDestinationArn":   d.ARN,
		"ApiDestinationState": string(d.State),
		"ConnectionArn":       d.ConnectionARN,
		"InvocationEndpoint":  d.InvocationEndpoint,
		"HttpMethod":          d.HttpMethod,
		"CreationTime":        d.CreatedAt.Unix(),
	}
	if d.Name != "" {
		result["Name"] = d.Name
	}
	if d.InvocationRateLimitPerSecond > 0 {
		result["InvocationRateLimitPerSecond"] = d.InvocationRateLimitPerSecond
	}
	if d.Description != "" {
		result["Description"] = d.Description
	}
	if !d.LastModifiedAt.IsZero() {
		result["LastModifiedTime"] = d.LastModifiedAt.Unix()
	}
	return result
}

// CreateApiDestination creates an API destination for EventBridge.
func (s *EventsService) CreateApiDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewValidationException("Api destination name is required")
	}
	if !validateResourceName(name, "api-destination") {
		return nil, awserrors.NewValidationException("Api destination name must match the pattern and be 1-64 characters")
	}

	connectionArn := request.GetParamLowerFirst(req.Parameters, "ConnectionArn")
	if connectionArn == "" {
		return nil, awserrors.NewValidationException("ConnectionArn is required")
	}

	httpMethod := request.GetParamLowerFirst(req.Parameters, "HttpMethod")
	if httpMethod == "" {
		httpMethod = "POST"
	}
	if !validHttpMethods[httpMethod] {
		return nil, awserrors.NewValidationException("HttpMethod must be one of: GET, POST, PUT, DELETE, HEAD, OPTIONS, PATCH")
	}

	invocationEndpoint := request.GetParamLowerFirst(req.Parameters, "InvocationEndpoint")
	if invocationEndpoint == "" {
		return nil, awserrors.NewValidationException("InvocationEndpoint is required")
	}

	apiDest := &eventsstore.ApiDestination{
		Name:               name,
		ConnectionARN:      connectionArn,
		HttpMethod:         httpMethod,
		InvocationEndpoint: invocationEndpoint,
	}

	if desc, ok := req.Parameters["Description"].(string); ok {
		if !validateDescription(desc) {
			return nil, awserrors.NewValidationException("Description must be at most 512 characters")
		}
		apiDest.Description = desc
	}

	if rateLimit := int32(request.GetIntParam(req.Parameters, "InvocationRateLimitPerSecond")); rateLimit > 0 {
		if !validateInvocationRateLimit(rateLimit) {
			return nil, awserrors.NewValidationException("InvocationRateLimitPerSecond must be between 1 and 300")
		}
		apiDest.InvocationRateLimitPerSecond = rateLimit
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.CreateApiDestination(ctx, apiDest); err != nil {
		return nil, mapStoreError(err, name)
	}

	return map[string]interface{}{
		"ApiDestinationArn":   apiDest.ARN,
		"CreationTime":        apiDest.CreatedAt.Unix(),
		"ApiDestinationState": string(apiDest.State),
	}, nil
}

// DeleteApiDestination deletes an API destination.
func (s *EventsService) DeleteApiDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewValidationException("Api destination name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeleteApiDestination(ctx, name); err != nil {
		return nil, mapStoreError(err, name)
	}

	return response.EmptyResponse(), nil
}

// DescribeApiDestination returns information about an API destination.
func (s *EventsService) DescribeApiDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewValidationException("Api destination name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	apiDest, err := store.GetApiDestination(ctx, name)
	if err != nil {
		return nil, mapStoreError(err, name)
	}

	result := apiDestinationToMap(apiDest)
	result["Name"] = apiDest.Name

	return result, nil
}

// UpdateApiDestination updates an existing EventBridge API destination.
func (s *EventsService) UpdateApiDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewValidationException("Api destination name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	apiDest, err := store.GetApiDestination(ctx, name)
	if err != nil {
		return nil, mapStoreError(err, name)
	}

	if desc, ok := req.Parameters["Description"].(string); ok {
		if !validateDescription(desc) {
			return nil, awserrors.NewValidationException("Description must be at most 512 characters")
		}
		apiDest.Description = desc
	}
	if httpMethod, ok := req.Parameters["HttpMethod"].(string); ok && httpMethod != "" {
		if !validHttpMethods[httpMethod] {
			return nil, awserrors.NewValidationException("HttpMethod must be one of: GET, POST, PUT, DELETE, HEAD, OPTIONS, PATCH")
		}
		apiDest.HttpMethod = httpMethod
	}
	if endpoint, ok := req.Parameters["InvocationEndpoint"].(string); ok && endpoint != "" {
		apiDest.InvocationEndpoint = endpoint
	}
	if connArn, ok := req.Parameters["ConnectionArn"].(string); ok && connArn != "" {
		apiDest.ConnectionARN = connArn
	}
	if rateLimit := int32(request.GetIntParam(req.Parameters, "InvocationRateLimitPerSecond")); rateLimit > 0 {
		if !validateInvocationRateLimit(rateLimit) {
			return nil, awserrors.NewValidationException("InvocationRateLimitPerSecond must be between 1 and 300")
		}
		apiDest.InvocationRateLimitPerSecond = rateLimit
	}

	apiDest.LastModifiedAt = time.Now().UTC()

	if err := store.UpdateApiDestination(ctx, apiDest); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ApiDestinationArn":   apiDest.ARN,
		"ApiDestinationState": string(apiDest.State),
		"CreationTime":        apiDest.CreatedAt.Unix(),
		"LastModifiedTime":    apiDest.LastModifiedAt.Unix(),
	}, nil
}

// ListApiDestinations lists API destinations with optional filtering.
func (s *EventsService) ListApiDestinations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	namePrefix := request.GetParamLowerFirst(req.Parameters, "NamePrefix")
	connectionArn := request.GetStringParam(req.Parameters, "ConnectionArn")

	limit := int32(request.GetIntParam(req.Parameters, "Limit"))
	if limit < 0 || limit > 100 {
		return nil, awserrors.NewValidationException("Limit must be between 0 and 100")
	}
	if limit == 0 {
		limit = 50
	}

	nextToken := pagination.GetMarker(req.Parameters, "NextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := store.ListApiDestinations(ctx, namePrefix, connectionArn, limit, nextToken)
	if err != nil {
		return nil, err
	}

	destinations := make([]map[string]interface{}, 0, len(result.ApiDestinations))
	for _, dest := range result.ApiDestinations {
		destinations = append(destinations, apiDestinationToMap(dest))
	}

	resp := map[string]interface{}{
		"ApiDestinations": destinations,
	}
	pagination.SetNextToken(resp, "NextToken", result.NextToken)

	return resp, nil
}

package eventbridge

import (
	"context"

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

// parseCreateApiDestinationInput reads the CreateApiDestination wire request
// into the transport-agnostic Core input.
func parseCreateApiDestinationInput(req *request.ParsedRequest) CreateApiDestinationInput {
	input := CreateApiDestinationInput{
		Name:               request.GetParamLowerFirst(req.Parameters, "Name"),
		ConnectionArn:      request.GetParamLowerFirst(req.Parameters, "ConnectionArn"),
		HttpMethod:         request.GetParamLowerFirst(req.Parameters, "HttpMethod"),
		InvocationEndpoint: request.GetParamLowerFirst(req.Parameters, "InvocationEndpoint"),
	}
	if desc, ok := req.Parameters["Description"].(string); ok {
		input.DescriptionSet = true
		input.Description = desc
	}
	input.InvocationRateLimit = int32(request.GetIntParam(req.Parameters, "InvocationRateLimitPerSecond"))
	return input
}

// parseUpdateApiDestinationInput reads the UpdateApiDestination wire request
// into the transport-agnostic Core input.
func parseUpdateApiDestinationInput(req *request.ParsedRequest) UpdateApiDestinationInput {
	input := UpdateApiDestinationInput{
		Name: request.GetParamLowerFirst(req.Parameters, "Name"),
	}
	if desc, ok := req.Parameters["Description"].(string); ok {
		input.DescriptionSet = true
		input.Description = desc
	}
	if httpMethod, ok := req.Parameters["HttpMethod"].(string); ok {
		input.HttpMethodSet = true
		input.HttpMethod = httpMethod
	}
	if endpoint, ok := req.Parameters["InvocationEndpoint"].(string); ok {
		input.InvocationEndpointSet = true
		input.InvocationEndpoint = endpoint
	}
	if connArn, ok := req.Parameters["ConnectionArn"].(string); ok {
		input.ConnectionArnSet = true
		input.ConnectionArn = connArn
	}
	input.InvocationRateLimit = int32(request.GetIntParam(req.Parameters, "InvocationRateLimitPerSecond"))
	return input
}

// CreateApiDestination creates an API destination for EventBridge.
func (s *EventsService) CreateApiDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := parseCreateApiDestinationInput(req)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	apiDest, err := s.createApiDestinationCore(ctx, store, input)
	if err != nil {
		return nil, err
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

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteApiDestinationCore(ctx, store, name); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DescribeApiDestination returns information about an API destination.
func (s *EventsService) DescribeApiDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	apiDest, err := s.getApiDestinationCore(ctx, store, name)
	if err != nil {
		return nil, err
	}

	result := apiDestinationToMap(apiDest)
	result["Name"] = apiDest.Name

	return result, nil
}

// UpdateApiDestination updates an existing EventBridge API destination.
func (s *EventsService) UpdateApiDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := parseUpdateApiDestinationInput(req)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	apiDest, err := s.updateApiDestinationCore(ctx, store, input)
	if err != nil {
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
	input := ListApiDestinationsInput{
		NamePrefix:    request.GetParamLowerFirst(req.Parameters, "NamePrefix"),
		ConnectionArn: request.GetStringParam(req.Parameters, "ConnectionArn"),
		Limit:         int32(request.GetIntParam(req.Parameters, "Limit")),
		NextToken:     pagination.GetMarker(req.Parameters, "NextToken"),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.listApiDestinationsCore(ctx, store, input)
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

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

func connectionToMap(c *eventsstore.Connection) map[string]interface{} {
	result := map[string]interface{}{
		"ConnectionArn":     c.ARN,
		"Name":              c.Name,
		"ConnectionState":   string(c.State),
		"AuthorizationType": c.AuthorizationType,
		"CreationTime":      c.CreatedAt.Unix(),
	}
	if c.Description != "" {
		result["Description"] = c.Description
	}
	if c.StateReason != "" {
		result["StateReason"] = c.StateReason
	}
	if !c.LastModifiedAt.IsZero() {
		result["LastModifiedTime"] = c.LastModifiedAt.Unix()
	}
	if !c.LastAuthorizedAt.IsZero() {
		result["LastAuthorizedTime"] = c.LastAuthorizedAt.Unix()
	}
	return result
}

// CreateConnection creates a new EventBridge connection.
func (s *EventsService) CreateConnection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewValidationException("Connection name is required")
	}

	authType := request.GetParamLowerFirst(req.Parameters, "AuthorizationType")
	if authType == "" {
		return nil, awserrors.NewValidationException("AuthorizationType is required")
	}
	if !validAuthTypes[authType] {
		return nil, awserrors.NewValidationException("AuthorizationType must be one of: API_KEY, BASIC, OAUTH_CLIENT_CREDENTIALS")
	}

	connection := &eventsstore.Connection{
		Name:              name,
		AuthorizationType: authType,
	}

	if desc, ok := req.Parameters["Description"].(string); ok {
		connection.Description = desc
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.CreateConnection(ctx, connection); err != nil {
		return nil, mapStoreError(err, name)
	}

	return map[string]interface{}{
		"ConnectionArn":   connection.ARN,
		"CreationTime":    connection.CreatedAt.Unix(),
		"ConnectionState": string(connection.State),
	}, nil
}

// DeleteConnection deletes an EventBridge connection.
func (s *EventsService) DeleteConnection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewValidationException("Connection name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	connection, err := store.GetConnection(ctx, name)
	if err != nil {
		if err == eventsstore.ErrConnectionNotFound {
			return nil, NewResourceNotFoundException("Connection '" + name + "' does not exist")
		}
		return nil, err
	}

	allDests, err := store.ListApiDestinations(ctx, "", 1000, "")
	if err == nil {
		for _, d := range allDests.ApiDestinations {
			if d.ConnectionARN == connection.ARN {
				return nil, awserrors.NewValidationException("Connection '" + name + "' is in use by API destination '" + d.Name + "'")
			}
		}
	}

	if err := store.DeleteConnection(ctx, name); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DescribeConnection returns information about a connection.
func (s *EventsService) DescribeConnection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewValidationException("Connection name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	connection, err := store.GetConnection(ctx, name)
	if err != nil {
		return nil, mapStoreError(err, name)
	}

	return connectionToMap(connection), nil
}

// UpdateConnection updates an existing EventBridge connection.
func (s *EventsService) UpdateConnection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewValidationException("Connection name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	connection, err := store.GetConnection(ctx, name)
	if err != nil {
		return nil, mapStoreError(err, name)
	}

	if desc, ok := req.Parameters["Description"].(string); ok {
		connection.Description = desc
	}
	if authType, ok := req.Parameters["AuthorizationType"].(string); ok && authType != "" {
		if !validAuthTypes[authType] {
			return nil, awserrors.NewValidationException("AuthorizationType must be one of: API_KEY, BASIC, OAUTH_CLIENT_CREDENTIALS")
		}
		connection.AuthorizationType = authType
	}

	connection.LastModifiedAt = time.Now().UTC()
	connection.State = eventsstore.ConnectionStateAuthorized

	if err := store.UpdateConnection(ctx, connection); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ConnectionArn":    connection.ARN,
		"ConnectionState":  string(connection.State),
		"CreationTime":     connection.CreatedAt.Unix(),
		"LastModifiedTime": connection.LastModifiedAt.Unix(),
	}, nil
}

// DeauthorizeConnection deauthorises an EventBridge connection, revoking its authorisation.
func (s *EventsService) DeauthorizeConnection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewValidationException("Connection name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeauthorizeConnection(ctx, name); err != nil {
		return nil, mapStoreError(err, name)
	}

	connection, _ := store.GetConnection(ctx, name)
	connArn := ""
	if connection != nil {
		connArn = connection.ARN
	}

	return map[string]interface{}{
		"ConnectionArn":   connArn,
		"ConnectionState": string(eventsstore.ConnectionStateDeauthorized),
	}, nil
}

// ListConnections lists connections with optional filtering.
func (s *EventsService) ListConnections(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	namePrefix := request.GetParamLowerFirst(req.Parameters, "NamePrefix")
	stateStr := request.GetParamLowerFirst(req.Parameters, "State")

	limit := int32(request.GetIntParam(req.Parameters, "Limit"))
	if limit == 0 {
		limit = 50
	}
	if limit > 100 {
		return nil, awserrors.NewValidationException("Limit must be between 1 and 100")
	}

	nextToken := pagination.GetMarker(req.Parameters, "NextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := store.ListConnections(ctx, namePrefix, stateStr, limit, nextToken)
	if err != nil {
		return nil, err
	}

	connections := make([]map[string]interface{}, 0, len(result.Connections))
	for _, conn := range result.Connections {
		connections = append(connections, connectionToMap(conn))
	}

	resp := map[string]interface{}{
		"Connections": connections,
	}
	pagination.SetNextToken(resp, "NextToken", result.NextToken)

	return resp, nil
}

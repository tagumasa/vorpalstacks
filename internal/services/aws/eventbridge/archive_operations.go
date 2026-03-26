package eventbridge

import (
	"context"
	"time"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	"vorpalstacks/internal/utils/aws/arn"
)

var validAuthTypes = map[string]bool{
	"API_KEY":                  true,
	"BASIC":                    true,
	"OAUTH_CLIENT_CREDENTIALS": true,
}

var validHttpMethods = map[string]bool{
	"GET":     true,
	"POST":    true,
	"PUT":     true,
	"DELETE":  true,
	"HEAD":    true,
	"OPTIONS": true,
	"PATCH":   true,
}

// CreateArchive creates an archive of events.
func (s *EventsService) CreateArchive(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "ArchiveName")
	if name == "" {
		return nil, NewValidationException("Archive name is required")
	}

	eventSourceArn := request.GetParamLowerFirst(req.Parameters, "EventSourceArn")
	if eventSourceArn == "" {
		return nil, NewValidationException("EventSourceArn is required")
	}

	eventBusName := arn.ExtractEventBusNameFromARN(eventSourceArn)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Check if event bus exists
	if _, err := store.GetEventBus(ctx, eventBusName); err != nil {
		if err == eventsstore.ErrEventBusNotFound {
			return nil, NewResourceNotFoundException("Event bus '" + eventBusName + "' does not exist")
		}
		return nil, err
	}

	archive := &eventsstore.Archive{
		Name:           name,
		EventBusName:   eventBusName,
		EventSourceARN: eventSourceArn,
	}

	if desc, ok := req.Parameters["Description"].(string); ok {
		archive.Description = desc
	}

	if pattern, ok := req.Parameters["EventPattern"].(string); ok {
		if !isValidEventPattern(pattern) {
			return nil, NewValidationException("EventPattern must be a valid JSON object")
		}
		archive.EventPattern = pattern
	}

	if retention := int32(request.GetIntParam(req.Parameters, "RetentionDays")); retention > 0 {
		archive.RetentionDays = retention
	}

	if err := store.CreateArchive(ctx, archive); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ArchiveArn":   archive.ARN,
		"CreationTime": archive.CreatedAt.Unix(),
		"State":        string(archive.State),
	}, nil
}

// DeleteArchive deletes an archive.
func (s *EventsService) DeleteArchive(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "ArchiveName")
	if name == "" {
		return nil, NewValidationException("Archive name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeleteArchive(ctx, name); err != nil {
		if err == eventsstore.ErrArchiveNotFound {
			return nil, NewResourceNotFoundException("Archive '" + name + "' does not exist")
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DescribeArchive returns information about an archive.
func (s *EventsService) DescribeArchive(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "ArchiveName")
	if name == "" {
		return nil, NewValidationException("Archive name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	archive, err := store.GetArchive(ctx, name)
	if err != nil {
		if err == eventsstore.ErrArchiveNotFound {
			return nil, NewResourceNotFoundException("Archive '" + name + "' does not exist")
		}
		return nil, err
	}

	result := map[string]interface{}{
		"ArchiveArn":     archive.ARN,
		"ArchiveName":    archive.Name,
		"EventSourceArn": archive.EventSourceARN,
		"State":          string(archive.State),
		"CreationTime":   archive.CreatedAt.Unix(),
		"EventCount":     archive.EventCount,
		"SizeBytes":      archive.SizeBytes,
	}

	if archive.Description != "" {
		result["Description"] = archive.Description
	}
	if archive.EventPattern != "" {
		result["EventPattern"] = archive.EventPattern
	}
	if archive.RetentionDays > 0 {
		result["RetentionDays"] = archive.RetentionDays
	}

	return result, nil
}

// ListArchives lists archives with optional filtering.
func (s *EventsService) ListArchives(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	namePrefix := request.GetParamLowerFirst(req.Parameters, "NamePrefix")
	stateStr := request.GetParamLowerFirst(req.Parameters, "State")

	limit := int32(request.GetIntParam(req.Parameters, "Limit"))
	if limit == 0 {
		limit = 50
	}
	if limit > 100 {
		return nil, NewValidationException("Limit must be between 1 and 100")
	}

	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := store.ListArchives(ctx, namePrefix, stateStr, limit, nextToken)
	if err != nil {
		return nil, err
	}

	archives := make([]map[string]interface{}, 0, len(result.Archives))
	for _, archive := range result.Archives {
		a := map[string]interface{}{
			"ArchiveArn":     archive.ARN,
			"ArchiveName":    archive.Name,
			"EventSourceArn": archive.EventSourceARN,
			"State":          string(archive.State),
			"CreationTime":   archive.CreatedAt.Unix(),
			"EventCount":     archive.EventCount,
			"SizeBytes":      archive.SizeBytes,
		}
		if archive.Description != "" {
			a["Description"] = archive.Description
		}
		if archive.EventPattern != "" {
			a["EventPattern"] = archive.EventPattern
		}
		if archive.RetentionDays > 0 {
			a["RetentionDays"] = archive.RetentionDays
		}
		archives = append(archives, a)
	}

	response := map[string]interface{}{
		"Archives": archives,
	}
	if result.NextToken != "" {
		response["NextToken"] = result.NextToken
	}

	return response, nil
}

// UpdateArchive updates an existing EventBridge archive.
func (s *EventsService) UpdateArchive(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "ArchiveName")
	if name == "" {
		return nil, NewValidationException("Archive name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	archive, err := store.GetArchive(ctx, name)
	if err != nil {
		if err == eventsstore.ErrArchiveNotFound {
			return nil, NewResourceNotFoundException("Archive '" + name + "' does not exist")
		}
		return nil, err
	}

	if desc, ok := req.Parameters["Description"].(string); ok {
		archive.Description = desc
	}
	if pattern, ok := req.Parameters["EventPattern"].(string); ok {
		if pattern != "" && !isValidEventPattern(pattern) {
			return nil, NewValidationException("EventPattern must be a valid JSON object")
		}
		archive.EventPattern = pattern
	}
	if retention := int32(request.GetIntParam(req.Parameters, "RetentionDays")); retention > 0 {
		archive.RetentionDays = retention
	}

	if err := store.UpdateArchive(ctx, archive); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ArchiveArn":     archive.ARN,
		"ArchiveName":    archive.Name,
		"EventSourceArn": archive.EventSourceARN,
		"State":          string(archive.State),
		"CreationTime":   archive.CreatedAt.Unix(),
		"EventCount":     archive.EventCount,
		"SizeBytes":      archive.SizeBytes,
	}, nil
}

// UpdateConnection updates an existing EventBridge connection.
func (s *EventsService) UpdateConnection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, NewValidationException("Connection name is required")
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

	if desc, ok := req.Parameters["Description"].(string); ok {
		connection.Description = desc
	}
	if authType, ok := req.Parameters["AuthorizationType"].(string); ok && authType != "" {
		if !validAuthTypes[authType] {
			return nil, NewValidationException("AuthorizationType must be one of: API_KEY, BASIC, OAUTH_CLIENT_CREDENTIALS")
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
		return nil, NewValidationException("Connection name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeauthorizeConnection(ctx, name); err != nil {
		if err == eventsstore.ErrConnectionNotFound {
			return nil, NewResourceNotFoundException("Connection '" + name + "' does not exist")
		}
		return nil, err
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

// UpdateApiDestination updates an existing EventBridge API destination.
func (s *EventsService) UpdateApiDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, NewValidationException("Api destination name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	apiDest, err := store.GetApiDestination(ctx, name)
	if err != nil {
		if err == eventsstore.ErrApiDestinationNotFound {
			return nil, NewResourceNotFoundException("Api destination '" + name + "' does not exist")
		}
		return nil, err
	}

	if desc, ok := req.Parameters["Description"].(string); ok {
		apiDest.Description = desc
	}
	if httpMethod, ok := req.Parameters["HttpMethod"].(string); ok && httpMethod != "" {
		if !validHttpMethods[httpMethod] {
			return nil, NewValidationException("HttpMethod must be one of: GET, POST, PUT, DELETE, HEAD, OPTIONS, PATCH")
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

// CreateConnection creates a new EventBridge connection.
func (s *EventsService) CreateConnection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, NewValidationException("Connection name is required")
	}

	authType := request.GetParamLowerFirst(req.Parameters, "AuthorizationType")
	if authType == "" {
		return nil, NewValidationException("AuthorizationType is required")
	}
	if !validAuthTypes[authType] {
		return nil, NewValidationException("AuthorizationType must be one of: API_KEY, BASIC, OAUTH_CLIENT_CREDENTIALS")
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
		return nil, err
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
		return nil, NewValidationException("Connection name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeleteConnection(ctx, name); err != nil {
		if err == eventsstore.ErrConnectionNotFound {
			return nil, NewResourceNotFoundException("Connection '" + name + "' does not exist")
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DescribeConnection returns information about a connection.
func (s *EventsService) DescribeConnection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, NewValidationException("Connection name is required")
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

	result := map[string]interface{}{
		"ConnectionArn":     connection.ARN,
		"Name":              connection.Name,
		"AuthorizationType": connection.AuthorizationType,
		"ConnectionState":   string(connection.State),
		"CreationTime":      connection.CreatedAt.Unix(),
	}

	if connection.Description != "" {
		result["Description"] = connection.Description
	}
	if connection.StateReason != "" {
		result["StateReason"] = connection.StateReason
	}
	if !connection.LastModifiedAt.IsZero() {
		result["LastModifiedTime"] = connection.LastModifiedAt.Unix()
	}
	if !connection.LastAuthorizedAt.IsZero() {
		result["LastAuthorizedTime"] = connection.LastAuthorizedAt.Unix()
	}

	return result, nil
}

// CreateApiDestination creates an API destination for EventBridge.
func (s *EventsService) CreateApiDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, NewValidationException("Api destination name is required")
	}

	connectionArn := request.GetParamLowerFirst(req.Parameters, "ConnectionArn")
	if connectionArn == "" {
		return nil, NewValidationException("ConnectionArn is required")
	}

	httpMethod := request.GetParamLowerFirst(req.Parameters, "HttpMethod")
	if httpMethod == "" {
		httpMethod = "POST"
	}
	if !validHttpMethods[httpMethod] {
		return nil, NewValidationException("HttpMethod must be one of: GET, POST, PUT, DELETE, HEAD, OPTIONS, PATCH")
	}

	invocationEndpoint := request.GetParamLowerFirst(req.Parameters, "InvocationEndpoint")
	if invocationEndpoint == "" {
		return nil, NewValidationException("InvocationEndpoint is required")
	}

	apiDest := &eventsstore.ApiDestination{
		Name:               name,
		ConnectionARN:      connectionArn,
		HttpMethod:         httpMethod,
		InvocationEndpoint: invocationEndpoint,
	}

	if desc, ok := req.Parameters["Description"].(string); ok {
		apiDest.Description = desc
	}

	if rateLimit := int32(request.GetIntParam(req.Parameters, "InvocationRateLimitPerSecond")); rateLimit > 0 {
		apiDest.InvocationRateLimitPerSecond = rateLimit
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.CreateApiDestination(ctx, apiDest); err != nil {
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
	if name == "" {
		return nil, NewValidationException("Api destination name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeleteApiDestination(ctx, name); err != nil {
		if err == eventsstore.ErrApiDestinationNotFound {
			return nil, NewResourceNotFoundException("Api destination '" + name + "' does not exist")
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DescribeApiDestination returns information about an API destination.
func (s *EventsService) DescribeApiDestination(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	if name == "" {
		return nil, NewValidationException("Api destination name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	apiDest, err := store.GetApiDestination(ctx, name)
	if err != nil {
		if err == eventsstore.ErrApiDestinationNotFound {
			return nil, NewResourceNotFoundException("Api destination '" + name + "' does not exist")
		}
		return nil, err
	}

	result := map[string]interface{}{
		"ApiDestinationArn":   apiDest.ARN,
		"Name":                apiDest.Name,
		"ConnectionArn":       apiDest.ConnectionARN,
		"HttpMethod":          apiDest.HttpMethod,
		"InvocationEndpoint":  apiDest.InvocationEndpoint,
		"ApiDestinationState": string(apiDest.State),
		"CreationTime":        apiDest.CreatedAt.Unix(),
	}

	if apiDest.Description != "" {
		result["Description"] = apiDest.Description
	}
	if apiDest.InvocationRateLimitPerSecond > 0 {
		result["InvocationRateLimitPerSecond"] = apiDest.InvocationRateLimitPerSecond
	}
	if !apiDest.LastModifiedAt.IsZero() {
		result["LastModifiedTime"] = apiDest.LastModifiedAt.Unix()
	}

	return result, nil
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
		return nil, NewValidationException("Limit must be between 1 and 100")
	}

	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")

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
		c := map[string]interface{}{
			"ConnectionArn":     conn.ARN,
			"ConnectionName":    conn.Name,
			"ConnectionState":   string(conn.State),
			"AuthorizationType": conn.AuthorizationType,
			"CreationTime":      conn.CreatedAt.Unix(),
		}
		if conn.Description != "" {
			c["Description"] = conn.Description
		}
		if !conn.LastModifiedAt.IsZero() {
			c["LastModifiedTime"] = conn.LastModifiedAt.Unix()
		}
		if !conn.LastAuthorizedAt.IsZero() {
			c["LastAuthorizedTime"] = conn.LastAuthorizedAt.Unix()
		}
		connections = append(connections, c)
	}

	response := map[string]interface{}{
		"Connections": connections,
	}
	if result.NextToken != "" {
		response["NextToken"] = result.NextToken
	}

	return response, nil
}

// ListApiDestinations lists API destinations with optional filtering.
func (s *EventsService) ListApiDestinations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	namePrefix := request.GetParamLowerFirst(req.Parameters, "NamePrefix")

	limit := int32(request.GetIntParam(req.Parameters, "Limit"))
	if limit == 0 {
		limit = 50
	}
	if limit > 100 {
		return nil, NewValidationException("Limit must be between 1 and 100")
	}

	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := store.ListApiDestinations(ctx, namePrefix, limit, nextToken)
	if err != nil {
		return nil, err
	}

	destinations := make([]map[string]interface{}, 0, len(result.ApiDestinations))
	for _, dest := range result.ApiDestinations {
		d := map[string]interface{}{
			"ApiDestinationArn":            dest.ARN,
			"ApiDestinationName":           dest.Name,
			"ApiDestinationState":          string(dest.State),
			"ConnectionArn":                dest.ConnectionARN,
			"InvocationEndpoint":           dest.InvocationEndpoint,
			"InvocationRateLimitPerSecond": dest.InvocationRateLimitPerSecond,
			"HttpMethod":                   dest.HttpMethod,
			"CreationTime":                 dest.CreatedAt.Unix(),
		}
		if dest.Description != "" {
			d["Description"] = dest.Description
		}
		if !dest.LastModifiedAt.IsZero() {
			d["LastModifiedTime"] = dest.LastModifiedAt.Unix()
		}
		destinations = append(destinations, d)
	}

	response := map[string]interface{}{
		"ApiDestinations": destinations,
	}
	if result.NextToken != "" {
		response["NextToken"] = result.NextToken
	}

	return response, nil
}

package eventbridge

import (
	"context"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
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

func archiveToMap(a *eventsstore.Archive) map[string]interface{} {
	result := map[string]interface{}{
		"ArchiveArn":     a.ARN,
		"ArchiveName":    a.Name,
		"EventSourceArn": a.EventSourceARN,
		"State":          string(a.State),
		"CreationTime":   a.CreatedAt.Unix(),
		"EventCount":     a.EventCount,
		"SizeBytes":      a.SizeBytes,
	}
	if a.Description != "" {
		result["Description"] = a.Description
	}
	if a.EventPattern != "" {
		result["EventPattern"] = a.EventPattern
	}
	if a.RetentionDays > 0 {
		result["RetentionDays"] = a.RetentionDays
	}
	return result
}

// CreateArchive creates an archive of events.
func (s *EventsService) CreateArchive(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "ArchiveName")
	if name == "" {
		return nil, awserrors.NewValidationException("Archive name is required")
	}

	eventSourceArn := request.GetParamLowerFirst(req.Parameters, "EventSourceArn")
	if eventSourceArn == "" {
		return nil, awserrors.NewValidationException("EventSourceArn is required")
	}

	eventBusName := arn.ExtractEventBusNameFromARN(eventSourceArn)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Check if event bus exists
	if _, err := store.GetEventBus(ctx, eventBusName); err != nil {
		return nil, mapStoreError(err, eventBusName)
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
			return nil, awserrors.NewValidationException("EventPattern must be a valid JSON object")
		}
		archive.EventPattern = pattern
	}

	if _, ok := req.Parameters["RetentionDays"]; ok {
		archive.RetentionDays = int32(request.GetIntParam(req.Parameters, "RetentionDays"))
	}

	if err := store.CreateArchive(ctx, archive); err != nil {
		return nil, mapStoreError(err, name)
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
		return nil, awserrors.NewValidationException("Archive name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeleteArchive(ctx, name); err != nil {
		return nil, mapStoreError(err, name)
	}

	_ = store.DeleteArchiveEvents(ctx, name)

	return response.EmptyResponse(), nil
}

// DescribeArchive returns information about an archive.
func (s *EventsService) DescribeArchive(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "ArchiveName")
	if name == "" {
		return nil, awserrors.NewValidationException("Archive name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	archive, err := store.GetArchive(ctx, name)
	if err != nil {
		return nil, mapStoreError(err, name)
	}

	return archiveToMap(archive), nil
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
		return nil, awserrors.NewValidationException("Limit must be between 1 and 100")
	}

	nextToken := pagination.GetMarker(req.Parameters, "NextToken")

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
		archives = append(archives, archiveToMap(archive))
	}

	response := map[string]interface{}{
		"Archives": archives,
	}
	pagination.SetNextToken(response, "NextToken", result.NextToken)

	return response, nil
}

// UpdateArchive updates an existing EventBridge archive.
func (s *EventsService) UpdateArchive(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "ArchiveName")
	if name == "" {
		return nil, awserrors.NewValidationException("Archive name is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	archive, err := store.GetArchive(ctx, name)
	if err != nil {
		return nil, mapStoreError(err, name)
	}

	if desc, ok := req.Parameters["Description"].(string); ok {
		archive.Description = desc
	}
	if pattern, ok := req.Parameters["EventPattern"].(string); ok {
		if pattern != "" && !isValidEventPattern(pattern) {
			return nil, awserrors.NewValidationException("EventPattern must be a valid JSON object")
		}
		archive.EventPattern = pattern
	}
	if _, ok := req.Parameters["RetentionDays"]; ok {
		archive.RetentionDays = int32(request.GetIntParam(req.Parameters, "RetentionDays"))
	}

	if err := store.UpdateArchive(ctx, archive); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ArchiveArn":   archive.ARN,
		"State":        string(archive.State),
		"CreationTime": archive.CreatedAt.Unix(),
	}, nil
}

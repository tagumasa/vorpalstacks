package eventbridge

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
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

// archiveToListItem serialises an Archive for the ListArchives response.
// Per Smithy Archive shape: ArchiveName, EventSourceArn, State,
// StateReason, RetentionDays, SizeBytes, EventCount, CreationTime.
func archiveToListItem(a *eventsstore.Archive) map[string]interface{} {
	result := map[string]interface{}{
		"ArchiveName":    a.Name,
		"EventSourceArn": a.EventSourceARN,
		"State":          string(a.State),
		"CreationTime":   a.CreatedAt.Unix(),
		"EventCount":     a.EventCount,
		"SizeBytes":      a.SizeBytes,
	}
	if a.StateReason != "" {
		result["StateReason"] = a.StateReason
	}
	if a.RetentionDays > 0 {
		result["RetentionDays"] = a.RetentionDays
	}
	return result
}

// archiveToDescribeMap serialises an Archive for the
// DescribeArchiveResponse shape. Includes the Describe-only fields:
// ArchiveArn, Description, EventPattern, KmsKeyIdentifier.
func archiveToDescribeMap(a *eventsstore.Archive) map[string]interface{} {
	result := archiveToListItem(a)
	result["ArchiveArn"] = a.ARN
	if a.Description != "" {
		result["Description"] = a.Description
	}
	if a.EventPattern != "" {
		result["EventPattern"] = a.EventPattern
	}
	if a.KmsKeyIdentifier != "" {
		result["KmsKeyIdentifier"] = a.KmsKeyIdentifier
	}
	return result
}

// parseArchiveMergeInput reads the create/update archive merge members from
// the wire request into the presence-flag carrier.
func parseArchiveMergeInput(req *request.ParsedRequest) ArchiveMergeMembers {
	var m ArchiveMergeMembers
	if desc, ok := req.Parameters["Description"].(string); ok {
		m.DescriptionSet = true
		m.Description = desc
	}
	if pattern, ok := req.Parameters["EventPattern"].(string); ok {
		m.EventPatternSet = true
		m.EventPattern = pattern
	}
	if _, ok := req.Parameters["RetentionDays"]; ok {
		m.RetentionDaysSet = true
		m.RetentionDays = int32(request.GetIntParam(req.Parameters, "RetentionDays"))
	}
	if kms, ok := req.Parameters["KmsKeyIdentifier"].(string); ok {
		m.KmsKeyIdentifierSet = true
		m.KmsKeyIdentifier = kms
	}
	return m
}

// CreateArchive creates an archive of events.
func (s *EventsService) CreateArchive(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := CreateArchiveInput{
		ArchiveName:         request.GetParamLowerFirst(req.Parameters, "ArchiveName"),
		EventSourceArn:      request.GetParamLowerFirst(req.Parameters, "EventSourceArn"),
		ArchiveMergeMembers: parseArchiveMergeInput(req),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	archive, err := s.createArchiveCore(ctx, store, input)
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"ArchiveArn":   archive.ARN,
		"CreationTime": archive.CreatedAt.Unix(),
		"State":        string(archive.State),
	}
	if archive.StateReason != "" {
		resp["StateReason"] = archive.StateReason
	}
	return resp, nil
}

// DeleteArchive deletes an archive.
func (s *EventsService) DeleteArchive(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "ArchiveName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteArchiveCore(ctx, store, name); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DescribeArchive returns information about an archive.
func (s *EventsService) DescribeArchive(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "ArchiveName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	archive, err := s.getArchiveCore(ctx, store, name)
	if err != nil {
		return nil, err
	}

	return archiveToDescribeMap(archive), nil
}

// ListArchives lists archives with optional filtering.
func (s *EventsService) ListArchives(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := ListArchivesInput{
		NamePrefix:     request.GetParamLowerFirst(req.Parameters, "NamePrefix"),
		State:          request.GetParamLowerFirst(req.Parameters, "State"),
		EventSourceArn: request.GetStringParam(req.Parameters, "EventSourceArn"),
		Limit:          int32(request.GetIntParam(req.Parameters, "Limit")),
		NextToken:      pagination.GetMarker(req.Parameters, "NextToken"),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.listArchivesCore(ctx, store, input)
	if err != nil {
		return nil, err
	}

	archives := make([]map[string]interface{}, 0, len(result.Archives))
	for _, archive := range result.Archives {
		archives = append(archives, archiveToListItem(archive))
	}

	response := map[string]interface{}{
		"Archives": archives,
	}
	pagination.SetNextToken(response, "NextToken", result.NextToken)

	return response, nil
}

// UpdateArchive updates an existing EventBridge archive.
func (s *EventsService) UpdateArchive(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := UpdateArchiveInput{
		ArchiveName:         request.GetParamLowerFirst(req.Parameters, "ArchiveName"),
		ArchiveMergeMembers: parseArchiveMergeInput(req),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	archive, err := s.updateArchiveCore(ctx, store, input)
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"ArchiveArn":   archive.ARN,
		"State":        string(archive.State),
		"CreationTime": archive.CreatedAt.Unix(),
	}
	if archive.StateReason != "" {
		resp["StateReason"] = archive.StateReason
	}
	return resp, nil
}

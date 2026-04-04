package appsync

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"

	"vorpalstacks/internal/services/aws/common/request"
)

// AssociateSourceGraphqlApi links a source GraphQL API to a merged API.
func (s *AppSyncService) AssociateSourceGraphqlApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	mergedApiId := request.GetStringParam(req.Parameters, "mergedApiIdentifier")
	if mergedApiId == "" {
		return nil, NewBadRequestException("mergedApiIdentifier is required")
	}

	sourceApiId := request.GetStringParam(req.Parameters, "sourceApiIdentifier")
	if sourceApiId == "" {
		return nil, NewBadRequestException("sourceApiIdentifier is required")
	}

	if _, err := store.GetGraphqlApiById(sourceApiId); err != nil {
		return mapStoreError(err)
	}

	description := request.GetStringParam(req.Parameters, "description")

	assoc := &appsyncstore.SourceApiAssociation{
		AssociationId:              uuid.New().String(),
		MergedApiId:                mergedApiId,
		SourceApiId:                sourceApiId,
		MergedApiArn:               fmt.Sprintf("arn:aws:appsync:%s:%s:apis/%s", store.GetRegion(), store.GetAccountID(), mergedApiId),
		SourceApiArn:               fmt.Sprintf("arn:aws:appsync:%s:%s:apis/%s", store.GetRegion(), store.GetAccountID(), sourceApiId),
		AssociationArn:             fmt.Sprintf("arn:aws:appsync:%s:%s:mergedApis/%s/sourceApiAssociations/%s", store.GetRegion(), store.GetAccountID(), mergedApiId, uuid.New().String()),
		SourceApiAssociationStatus: "MERGE_SCHEDULED",
		Description:                description,
	}

	if err := store.CreateMergedApiAssociation(assoc); err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{"sourceApiAssociation": mergedApiAssociationToMap(assoc)}, nil
}

// GetSourceApiAssociation retrieves a source API association.
func (s *AppSyncService) GetSourceApiAssociation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	mergedApiId := request.GetStringParam(req.Parameters, "mergedApiIdentifier")
	if mergedApiId == "" {
		return nil, NewBadRequestException("mergedApiIdentifier is required")
	}

	associationId := request.GetStringParam(req.Parameters, "associationId")
	if associationId == "" {
		return nil, NewBadRequestException("associationId is required")
	}

	assoc, err := store.GetMergedApiAssociation(mergedApiId, associationId)
	if err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{"sourceApiAssociation": mergedApiAssociationToMap(assoc)}, nil
}

// UpdateSourceApiAssociation updates a source API association description and config.
func (s *AppSyncService) UpdateSourceApiAssociation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	mergedApiId := request.GetStringParam(req.Parameters, "mergedApiIdentifier")
	if mergedApiId == "" {
		return nil, NewBadRequestException("mergedApiIdentifier is required")
	}

	associationId := request.GetStringParam(req.Parameters, "associationId")
	if associationId == "" {
		return nil, NewBadRequestException("associationId is required")
	}

	assoc, err := store.GetMergedApiAssociation(mergedApiId, associationId)
	if err != nil {
		return mapStoreError(err)
	}

	description := request.GetStringParam(req.Parameters, "description")
	if description != "" {
		assoc.Description = description
	}

	if err := store.UpdateMergedApiAssociation(assoc); err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{"sourceApiAssociation": mergedApiAssociationToMap(assoc)}, nil
}

// DisassociateSourceGraphqlApi removes a source API association from a merged API.
func (s *AppSyncService) DisassociateSourceGraphqlApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	mergedApiId := request.GetStringParam(req.Parameters, "mergedApiIdentifier")
	if mergedApiId == "" {
		return nil, NewBadRequestException("mergedApiIdentifier is required")
	}

	associationId := request.GetStringParam(req.Parameters, "associationId")
	if associationId == "" {
		return nil, NewBadRequestException("associationId is required")
	}

	assoc, err := store.GetMergedApiAssociation(mergedApiId, associationId)
	if err != nil {
		return mapStoreError(err)
	}

	if err := store.DeleteMergedApiAssociation(mergedApiId, associationId); err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{"sourceApiAssociationStatus": assoc.SourceApiAssociationStatus}, nil
}

// StartSchemaMerge initiates a schema merge for a source API association.
func (s *AppSyncService) StartSchemaMerge(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	mergedApiId := request.GetStringParam(req.Parameters, "mergedApiIdentifier")
	if mergedApiId == "" {
		return nil, NewBadRequestException("mergedApiIdentifier is required")
	}

	associationId := request.GetStringParam(req.Parameters, "associationId")
	if associationId == "" {
		return nil, NewBadRequestException("associationId is required")
	}

	assoc, err := store.GetMergedApiAssociation(mergedApiId, associationId)
	if err != nil {
		return mapStoreError(err)
	}

	assoc.SourceApiAssociationStatus = "MERGE_IN_PROGRESS"
	if err := store.UpdateMergedApiAssociation(assoc); err != nil {
		return mapStoreError(err)
	}

	assoc.SourceApiAssociationStatus = "MERGE_SUCCESS"
	_ = store.UpdateMergedApiAssociation(assoc)

	return map[string]interface{}{"sourceApiAssociationStatus": "MERGE_SUCCESS"}, nil
}

// AssociateMergedGraphqlApi links a merged API from the source API side.
func (s *AppSyncService) AssociateMergedGraphqlApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	sourceApiId := request.GetStringParam(req.Parameters, "sourceApiIdentifier")
	if sourceApiId == "" {
		return nil, NewBadRequestException("sourceApiIdentifier is required")
	}

	mergedApiId := request.GetStringParam(req.Parameters, "mergedApiIdentifier")
	if mergedApiId == "" {
		return nil, NewBadRequestException("mergedApiIdentifier is required")
	}

	if _, err := store.GetGraphqlApiById(sourceApiId); err != nil {
		return mapStoreError(err)
	}

	description := request.GetStringParam(req.Parameters, "description")

	assoc := &appsyncstore.SourceApiAssociation{
		AssociationId:              uuid.New().String(),
		MergedApiId:                mergedApiId,
		SourceApiId:                sourceApiId,
		MergedApiArn:               fmt.Sprintf("arn:aws:appsync:%s:%s:apis/%s", store.GetRegion(), store.GetAccountID(), mergedApiId),
		SourceApiArn:               fmt.Sprintf("arn:aws:appsync:%s:%s:apis/%s", store.GetRegion(), store.GetAccountID(), sourceApiId),
		AssociationArn:             fmt.Sprintf("arn:aws:appsync:%s:%s:sourceApis/%s/mergedApiAssociations/%s", store.GetRegion(), store.GetAccountID(), sourceApiId, uuid.New().String()),
		SourceApiAssociationStatus: "MERGE_SCHEDULED",
		Description:                description,
	}

	if err := store.CreateMergedApiAssociation(assoc); err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{"sourceApiAssociation": mergedApiAssociationToMap(assoc)}, nil
}

// DisassociateMergedGraphqlApi removes a merged API association from the source API side.
// The source-side path only provides sourceApiIdentifier and associationId, so we
// look up the association by UUID to recover the mergedApiId for the composite key.
func (s *AppSyncService) DisassociateMergedGraphqlApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	associationId := request.GetStringParam(req.Parameters, "associationId")
	if associationId == "" {
		return nil, NewBadRequestException("associationId is required")
	}

	assoc, err := store.GetMergedApiAssociationById(associationId)
	if err != nil {
		return mapStoreError(err)
	}

	if err := store.DeleteMergedApiAssociation(assoc.MergedApiId, associationId); err != nil {
		return mapStoreError(err)
	}

	return map[string]interface{}{"sourceApiAssociationStatus": assoc.SourceApiAssociationStatus}, nil
}

// ListSourceApiAssociations lists source API associations for a GraphQL API.
// Uses filter-based listing by sourceApiId since the store key is composite (mergedApiId/associationId).
func (s *AppSyncService) ListSourceApiAssociations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	apiId := request.GetStringParam(req.Parameters, "apiId")
	if apiId == "" {
		return nil, NewBadRequestException("apiId is required")
	}

	if _, err := store.GetGraphqlApiById(apiId); err != nil {
		return mapStoreError(err)
	}

	opts := parsePaginationOptions(req)
	assocs, nextToken, err := store.ListAssociationsBySourceApi(apiId, opts)
	if err != nil {
		return mapStoreError(err)
	}

	items := make([]map[string]interface{}, 0, len(assocs))
	for _, a := range assocs {
		items = append(items, mergedApiAssociationSummaryToMap(a))
	}

	response := map[string]interface{}{"sourceApiAssociationSummaries": items}
	if nextToken != "" {
		response["nextToken"] = nextToken
	}
	return response, nil
}

func mergedApiAssociationToMap(a *appsyncstore.SourceApiAssociation) map[string]interface{} {
	result := map[string]interface{}{
		"associationId":              a.AssociationId,
		"mergedApiId":                a.MergedApiId,
		"sourceApiId":                a.SourceApiId,
		"sourceApiAssociationStatus": a.SourceApiAssociationStatus,
	}
	if a.AssociationArn != "" {
		result["associationArn"] = a.AssociationArn
	}
	if a.MergedApiArn != "" {
		result["mergedApiArn"] = a.MergedApiArn
	}
	if a.SourceApiArn != "" {
		result["sourceApiArn"] = a.SourceApiArn
	}
	if a.Description != "" {
		result["description"] = a.Description
	}
	return result
}

func mergedApiAssociationSummaryToMap(a *appsyncstore.SourceApiAssociation) map[string]interface{} {
	result := map[string]interface{}{
		"associationId": a.AssociationId,
		"mergedApiId":   a.MergedApiId,
		"sourceApiId":   a.SourceApiId,
	}
	if a.AssociationArn != "" {
		result["associationArn"] = a.AssociationArn
	}
	if a.MergedApiArn != "" {
		result["mergedApiArn"] = a.MergedApiArn
	}
	if a.SourceApiArn != "" {
		result["sourceApiArn"] = a.SourceApiArn
	}
	if a.Description != "" {
		result["description"] = a.Description
	}
	return result
}

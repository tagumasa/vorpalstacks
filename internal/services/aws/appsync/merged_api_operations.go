package appsync

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// AssociateSourceGraphqlApi creates a source API association addressed from
// the merged API side.
// POST /v1/mergedApis/{mergedApiIdentifier}/sourceApiAssociations
func (s *AppSyncService) AssociateSourceGraphqlApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	assocConfig, err := parseSourceApiAssociationConfig(req.Parameters)
	if err != nil {
		return nil, err
	}

	in := associateSourceApiInput{
		MergedApiId: request.GetStringParam(req.Parameters, "mergedApiIdentifier"),
		SourceApiId: request.GetStringParam(req.Parameters, "sourceApiIdentifier"),
		Description: request.GetStringParam(req.Parameters, "description"),
		AssocConfig: assocConfig,
	}

	assoc, err := s.associateSourceGraphqlApiCore(store, in)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"sourceApiAssociation": mergedApiAssociationToMap(assoc)}, nil
}

// GetSourceApiAssociation retrieves one source API association of a merged
// API.
// GET /v1/mergedApis/{mergedApiIdentifier}/sourceApiAssociations/{associationId}
func (s *AppSyncService) GetSourceApiAssociation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	assoc, err := s.getSourceApiAssociationCore(store,
		request.GetStringParam(req.Parameters, "mergedApiIdentifier"),
		request.GetStringParam(req.Parameters, "associationId"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"sourceApiAssociation": mergedApiAssociationToMap(assoc)}, nil
}

// UpdateSourceApiAssociation updates an existing source API association.
// POST /v1/mergedApis/{mergedApiIdentifier}/sourceApiAssociations/{associationId}
func (s *AppSyncService) UpdateSourceApiAssociation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	assocConfig, err := parseSourceApiAssociationConfig(req.Parameters)
	if err != nil {
		return nil, err
	}

	assoc, err := s.updateSourceApiAssociationCore(store,
		request.GetStringParam(req.Parameters, "mergedApiIdentifier"),
		request.GetStringParam(req.Parameters, "associationId"),
		request.GetStringParam(req.Parameters, "description"),
		assocConfig)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"sourceApiAssociation": mergedApiAssociationToMap(assoc)}, nil
}

// DisassociateSourceGraphqlApi schedules the deletion of a source API
// association addressed from the merged API side.
// DELETE /v1/mergedApis/{mergedApiIdentifier}/sourceApiAssociations/{associationId}
func (s *AppSyncService) DisassociateSourceGraphqlApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	status, err := s.disassociateSourceGraphqlApiCore(store,
		request.GetStringParam(req.Parameters, "mergedApiIdentifier"),
		request.GetStringParam(req.Parameters, "associationId"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"sourceApiAssociationStatus": status}, nil
}

// StartSchemaMerge triggers a schema merge for a source API association.
// POST /v1/mergedApis/{mergedApiIdentifier}/sourceApiAssociations/{associationId}/merge
func (s *AppSyncService) StartSchemaMerge(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	status, err := s.startSchemaMergeCore(store,
		request.GetStringParam(req.Parameters, "mergedApiIdentifier"),
		request.GetStringParam(req.Parameters, "associationId"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"sourceApiAssociationStatus": status}, nil
}

// AssociateMergedGraphqlApi creates a source API association addressed from
// the source API side.
// POST /v1/sourceApis/{sourceApiIdentifier}/mergedApiAssociations
func (s *AppSyncService) AssociateMergedGraphqlApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	assocConfig, err := parseSourceApiAssociationConfig(req.Parameters)
	if err != nil {
		return nil, err
	}

	in := associateSourceApiInput{
		MergedApiId: request.GetStringParam(req.Parameters, "mergedApiIdentifier"),
		SourceApiId: request.GetStringParam(req.Parameters, "sourceApiIdentifier"),
		Description: request.GetStringParam(req.Parameters, "description"),
		AssocConfig: assocConfig,
	}

	assoc, err := s.associateMergedGraphqlApiCore(store, in)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"sourceApiAssociation": mergedApiAssociationToMap(assoc)}, nil
}

// DisassociateMergedGraphqlApi schedules the deletion of a source API
// association addressed from the source API side.
// DELETE /v1/sourceApis/{sourceApiIdentifier}/mergedApiAssociations/{associationId}
func (s *AppSyncService) DisassociateMergedGraphqlApi(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	status, err := s.disassociateMergedGraphqlApiCore(store,
		request.GetStringParam(req.Parameters, "sourceApiIdentifier"),
		request.GetStringParam(req.Parameters, "associationId"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{"sourceApiAssociationStatus": status}, nil
}

// ListSourceApiAssociations lists the source API associations of a GraphQL
// API acting as a merged API.
// GET /v1/apis/{apiId}/sourceApiAssociations
func (s *AppSyncService) ListSourceApiAssociations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return mapStoreError(err)
	}

	assocs, nextToken, err := s.listSourceApiAssociationsCore(store,
		request.GetStringParam(req.Parameters, "apiId"),
		request.GetIntParam(req.Parameters, "maxResults"),
		request.GetStringParam(req.Parameters, "nextToken"))
	if err != nil {
		return nil, err
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

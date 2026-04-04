package appsync

import (
	"context"

	"vorpalstacks/internal/services/aws/common/request"
)

// StartDataSourceIntrospection initiates an introspection of a relational data source.
// Returns HTTP 501 — external IdP/data-source introspection is not supported locally.
// POST /v1/datasources/introspections
func (s *AppSyncService) StartDataSourceIntrospection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return nil, ErrInternalFailureException
}

// GetDataSourceIntrospection retrieves the result of a data source introspection.
// Returns HTTP 501 — external IdP/data-source introspection is not supported locally.
// GET /v1/datasources/introspections/{introspectionId}
func (s *AppSyncService) GetDataSourceIntrospection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return nil, ErrInternalFailureException
}

// ListTypesByAssociation returns types associated with a merged API source association.
// Returns HTTP 501 — merged API is not yet implemented.
// GET /v1/mergedApis/{mergedApiIdentifier}/sourceApiAssociations/{associationId}/types
func (s *AppSyncService) ListTypesByAssociation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return nil, ErrInternalFailureException
}

// EvaluateCode evaluates JavaScript code against an AppSync runtime.
// Returns HTTP 501 — APPSYNC_JS runtime is not yet implemented.
// POST /v1/dataplane-evaluatecode
func (s *AppSyncService) EvaluateCode(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return nil, ErrInternalFailureException
}

// EvaluateMappingTemplate evaluates a VTL mapping template.
// Returns HTTP 501 — template evaluation endpoint is not yet implemented.
// POST /v1/dataplane-evaluatetemplate
func (s *AppSyncService) EvaluateMappingTemplate(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return nil, ErrInternalFailureException
}

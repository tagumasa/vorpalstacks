package athena

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	athenastore "vorpalstacks/internal/store/aws/athena"
)

// CreateNamedQuery creates a new named query in the Athena workgroup.
func (s *AthenaService) CreateNamedQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := CreateNamedQueryInput{
		Name:        request.GetParamCaseInsensitive(req.Parameters, "Name"),
		Description: request.GetParamCaseInsensitive(req.Parameters, "Description"),
		Database:    request.GetParamCaseInsensitive(req.Parameters, "Database"),
		QueryString: request.GetParamCaseInsensitive(req.Parameters, "QueryString"),
		WorkGroup:   request.GetParamCaseInsensitive(req.Parameters, "WorkGroup"),
	}

	namedQueryId, err := s.createNamedQueryCore(reqCtx, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"NamedQueryId": namedQueryId,
	}, nil
}

// GetNamedQuery retrieves the specified named query.
func (s *AthenaService) GetNamedQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	namedQueryId := request.GetParamCaseInsensitive(req.Parameters, "NamedQueryId")

	namedQuery, err := s.getNamedQueryCore(reqCtx, namedQueryId)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"NamedQuery": s.namedQueryToResponse(namedQuery),
	}, nil
}

// DeleteNamedQuery deletes the specified named query.
func (s *AthenaService) DeleteNamedQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	namedQueryId := request.GetParamCaseInsensitive(req.Parameters, "NamedQueryId")

	if err := s.deleteNamedQueryCore(reqCtx, namedQueryId); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListNamedQueries retrieves a list of named queries in the specified workgroup.
func (s *AthenaService) ListNamedQueries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	maxResults, hasMaxResults := request.GetIntParamCaseInsensitive(req.Parameters, "MaxResults")
	input := ListNamedQueriesInput{
		WorkGroup:     request.GetParamCaseInsensitive(req.Parameters, "WorkGroup"),
		MaxResults:    maxResults,
		HasMaxResults: hasMaxResults,
		NextToken:     pagination.GetMarker(req.Parameters, "NextToken"),
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	ids, nextMarker, err := listNamedQueriesCore(stores, input)
	if err != nil {
		return nil, err
	}

	return pagination.BuildListResponse("NamedQueryIds", ids, nextMarker), nil
}

// UpdateNamedQuery updates the specified named query.
// Per the Smithy model, NamedQueryId, Name, and QueryString are all REQUIRED.
func (s *AthenaService) UpdateNamedQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := UpdateNamedQueryInput{
		NamedQueryId: request.GetParamCaseInsensitive(req.Parameters, "NamedQueryId"),
		Name:         request.GetParamCaseInsensitive(req.Parameters, "Name"),
		Description:  request.GetParamCaseInsensitive(req.Parameters, "Description"),
		QueryString:  request.GetParamCaseInsensitive(req.Parameters, "QueryString"),
	}

	if err := s.updateNamedQueryCore(reqCtx, input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// BatchGetNamedQuery retrieves multiple named queries in a single request.
func (s *AthenaService) BatchGetNamedQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := BatchGetNamedQueryInput{
		NamedQueryIds: request.GetArrayParam(req.Parameters, "NamedQueryIds"),
	}

	namedQueries, unprocessedIds, err := s.batchGetNamedQueryCore(reqCtx, input)
	if err != nil {
		return nil, err
	}

	var queryResponses []map[string]interface{}
	for _, nq := range namedQueries {
		queryResponses = append(queryResponses, s.namedQueryToResponse(nq))
	}

	var unprocessedResponses []map[string]interface{}
	for _, id := range unprocessedIds {
		unprocessedResponses = append(unprocessedResponses, map[string]interface{}{
			"NamedQueryId": id,
		})
	}

	return map[string]interface{}{
		"NamedQueries":             queryResponses,
		"UnprocessedNamedQueryIds": unprocessedResponses,
	}, nil
}

func (s *AthenaService) namedQueryToResponse(nq *athenastore.NamedQuery) map[string]interface{} {
	return map[string]interface{}{
		"Name":         nq.Name,
		"Description":  nq.Description,
		"Database":     nq.Database,
		"QueryString":  nq.QueryString,
		"NamedQueryId": nq.NamedQueryId,
		"WorkGroup":    nq.WorkGroup,
	}
}

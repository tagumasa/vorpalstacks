package athena

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	athenastore "vorpalstacks/internal/store/aws/athena"
)

// CreatePreparedStatement creates a new prepared statement in the Athena workgroup.
func (s *AthenaService) CreatePreparedStatement(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := CreatePreparedStatementInput{
		StatementName:  request.GetParamCaseInsensitive(req.Parameters, "StatementName"),
		WorkGroup:      request.GetParamCaseInsensitive(req.Parameters, "WorkGroup"),
		Description:    request.GetParamCaseInsensitive(req.Parameters, "Description"),
		QueryStatement: request.GetParamCaseInsensitive(req.Parameters, "QueryStatement"),
	}

	if err := s.createPreparedStatementCore(reqCtx, input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetPreparedStatement retrieves the specified prepared statement.
func (s *AthenaService) GetPreparedStatement(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := GetPreparedStatementInput{
		StatementName: request.GetParamCaseInsensitive(req.Parameters, "StatementName"),
		WorkGroup:     request.GetParamCaseInsensitive(req.Parameters, "WorkGroup"),
	}

	preparedStatement, err := s.getPreparedStatementCore(reqCtx, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"PreparedStatement": s.preparedStatementToResponse(preparedStatement),
	}, nil
}

// DeletePreparedStatement deletes the specified prepared statement.
func (s *AthenaService) DeletePreparedStatement(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := GetPreparedStatementInput{
		StatementName: request.GetParamCaseInsensitive(req.Parameters, "StatementName"),
		WorkGroup:     request.GetParamCaseInsensitive(req.Parameters, "WorkGroup"),
	}

	if err := s.deletePreparedStatementCore(reqCtx, input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListPreparedStatements retrieves a list of prepared statements in the specified workgroup.
func (s *AthenaService) ListPreparedStatements(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	maxResults, hasMaxResults := request.GetIntParamCaseInsensitive(req.Parameters, "MaxResults")
	input := ListPreparedStatementsInput{
		WorkGroup:     request.GetParamCaseInsensitive(req.Parameters, "WorkGroup"),
		MaxResults:    maxResults,
		HasMaxResults: hasMaxResults,
		NextToken:     pagination.GetMarker(req.Parameters, "NextToken"),
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	preparedStatements, nextMarker, err := listPreparedStatementsCore(stores, input)
	if err != nil {
		return nil, err
	}

	summaries := make([]map[string]interface{}, 0, len(preparedStatements))
	for _, ps := range preparedStatements {
		summaries = append(summaries, map[string]interface{}{
			"StatementName":    ps.StatementName,
			"LastModifiedTime": float64(ps.LastModifiedTime.UnixNano()) / 1e9,
		})
	}

	return pagination.BuildListResponse("PreparedStatements", summaries, nextMarker), nil
}

// UpdatePreparedStatement updates the specified prepared statement.
// Per the Smithy model, StatementName, WorkGroup, and QueryStatement are all REQUIRED.
func (s *AthenaService) UpdatePreparedStatement(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := UpdatePreparedStatementInput{
		StatementName:  request.GetParamCaseInsensitive(req.Parameters, "StatementName"),
		WorkGroup:      request.GetParamCaseInsensitive(req.Parameters, "WorkGroup"),
		QueryStatement: request.GetParamCaseInsensitive(req.Parameters, "QueryStatement"),
		Description:    request.GetParamCaseInsensitive(req.Parameters, "Description"),
	}

	if err := s.updatePreparedStatementCore(reqCtx, input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// BatchGetPreparedStatement retrieves multiple prepared statements in a single request.
func (s *AthenaService) BatchGetPreparedStatement(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := BatchGetPreparedStatementInput{
		WorkGroup:              request.GetParamCaseInsensitive(req.Parameters, "WorkGroup"),
		PreparedStatementNames: request.GetArrayParam(req.Parameters, "PreparedStatementNames"),
	}

	preparedStatements, unprocessedNames, err := s.batchGetPreparedStatementCore(reqCtx, input)
	if err != nil {
		return nil, err
	}

	var statementResponses []map[string]interface{}
	for _, ps := range preparedStatements {
		statementResponses = append(statementResponses, s.preparedStatementToResponse(ps))
	}

	var unprocessedResponses []map[string]interface{}
	for _, name := range unprocessedNames {
		unprocessedResponses = append(unprocessedResponses, map[string]interface{}{
			"StatementName": name,
		})
	}

	return map[string]interface{}{
		"PreparedStatements":                statementResponses,
		"UnprocessedPreparedStatementNames": unprocessedResponses,
	}, nil
}

func (s *AthenaService) preparedStatementToResponse(ps *athenastore.PreparedStatement) map[string]interface{} {
	return map[string]interface{}{
		"StatementName":    ps.StatementName,
		"QueryStatement":   ps.QueryStatement,
		"WorkGroupName":    ps.WorkGroupName,
		"Description":      ps.Description,
		"LastModifiedTime": float64(ps.LastModifiedTime.UnixNano()) / 1e9,
	}
}

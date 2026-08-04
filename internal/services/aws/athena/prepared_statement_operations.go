package athena

import (
	"context"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	athenastore "vorpalstacks/internal/store/aws/athena"
)

// CreatePreparedStatement creates a new prepared statement in the Athena workgroup.
func (s *AthenaService) CreatePreparedStatement(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	statementName := request.GetParamCaseInsensitive(req.Parameters, "StatementName")
	if statementName == "" {
		return nil, ErrInvalidRequestException
	}

	if err := validateStatementName(statementName); err != nil {
		return nil, err
	}

	workGroup := request.GetParamCaseInsensitive(req.Parameters, "WorkGroup")
	if workGroup == "" {
		workGroup = "primary"
	}

	description := request.GetParamCaseInsensitive(req.Parameters, "Description")
	if description != "" {
		if err := validateDescriptionString(description); err != nil {
			return nil, err
		}
	}
	queryStatement := request.GetParamCaseInsensitive(req.Parameters, "QueryStatement")

	if len(queryStatement) > maxQueryStringSize {
		return nil, ErrInvalidRequestException
	}

	if queryStatement == "" {
		return nil, ErrInvalidRequestException
	}

	preparedStatement := &athenastore.PreparedStatement{
		StatementName:  statementName,
		WorkGroupName:  workGroup,
		QueryStatement: queryStatement,
		Description:    description,
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := stores.preparedStatementStore.CreatePreparedStatement(preparedStatement); err != nil {
		if err == athenastore.ErrPreparedStatementAlreadyExists {
			return nil, ErrResourceAlreadyExistsException
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetPreparedStatement retrieves the specified prepared statement.
func (s *AthenaService) GetPreparedStatement(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	statementName := request.GetParamCaseInsensitive(req.Parameters, "StatementName")
	if statementName == "" {
		return nil, ErrInvalidRequestException
	}

	workGroup := request.GetParamCaseInsensitive(req.Parameters, "WorkGroup")
	if workGroup == "" {
		workGroup = "primary"
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	preparedStatement, err := stores.preparedStatementStore.GetPreparedStatement(workGroup, statementName)
	if err != nil {
		if err == athenastore.ErrPreparedStatementNotFound {
			return nil, preparedStatementNotFound(statementName)
		}
		return nil, err
	}

	return map[string]interface{}{
		"PreparedStatement": s.preparedStatementToResponse(preparedStatement),
	}, nil
}

// DeletePreparedStatement deletes the specified prepared statement.
func (s *AthenaService) DeletePreparedStatement(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	statementName := request.GetParamCaseInsensitive(req.Parameters, "StatementName")
	if statementName == "" {
		return nil, ErrInvalidRequestException
	}

	workGroup := request.GetParamCaseInsensitive(req.Parameters, "WorkGroup")
	if workGroup == "" {
		workGroup = "primary"
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := stores.preparedStatementStore.DeletePreparedStatement(workGroup, statementName); err != nil {
		if err == athenastore.ErrPreparedStatementNotFound {
			return nil, preparedStatementNotFound(statementName)
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListPreparedStatements retrieves a list of prepared statements in the specified workgroup.
func (s *AthenaService) ListPreparedStatements(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	workGroup := request.GetParamCaseInsensitive(req.Parameters, "WorkGroup")
	if workGroup == "" {
		workGroup = "primary"
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	preparedStatements, err := stores.preparedStatementStore.ListPreparedStatements(workGroup)
	if err != nil {
		return nil, err
	}

	var summaries []map[string]interface{}
	for _, ps := range preparedStatements {
		summaries = append(summaries, map[string]interface{}{
			"StatementName":    ps.StatementName,
			"LastModifiedTime": float64(ps.LastModifiedTime.UnixNano()) / 1e9,
		})
	}

	maxResults, err := validateMaxResults(req.Parameters, 50, 1, 50)
	if err != nil {
		return nil, err
	}
	marker := pagination.GetMarker(req.Parameters, "NextToken")
	pageResult := pagination.PaginateSlice(summaries, marker, maxResults, func(item map[string]interface{}) string {
		return item["StatementName"].(string)
	})

	return pagination.BuildListResponse("PreparedStatements", pageResult.Items, pageResult.NextMarker), nil
}

// UpdatePreparedStatement updates the specified prepared statement.
// Per the Smithy model, StatementName, WorkGroup, and QueryStatement are all REQUIRED.
func (s *AthenaService) UpdatePreparedStatement(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	statementName := request.GetParamCaseInsensitive(req.Parameters, "StatementName")
	if statementName == "" {
		return nil, ErrInvalidRequestException
	}

	workGroup := request.GetParamCaseInsensitive(req.Parameters, "WorkGroup")
	if workGroup == "" {
		workGroup = "primary"
	}

	queryStatement := request.GetParamCaseInsensitive(req.Parameters, "QueryStatement")
	if queryStatement == "" {
		return nil, awserrors.NewInvalidParameterException("QueryStatement is required for UpdatePreparedStatement")
	}

	if len(queryStatement) > maxQueryStringSize {
		return nil, ErrInvalidRequestException
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	preparedStatement, err := stores.preparedStatementStore.GetPreparedStatement(workGroup, statementName)
	if err != nil {
		if err == athenastore.ErrPreparedStatementNotFound {
			return nil, preparedStatementNotFound(statementName)
		}
		return nil, err
	}

	preparedStatement.QueryStatement = queryStatement

	description := request.GetParamCaseInsensitive(req.Parameters, "Description")
	if description != "" {
		preparedStatement.Description = description
	}

	if err := stores.preparedStatementStore.UpdatePreparedStatement(preparedStatement); err != nil {
		if err == athenastore.ErrPreparedStatementNotFound {
			return nil, preparedStatementNotFound(statementName)
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// BatchGetPreparedStatement retrieves multiple prepared statements in a single request.
func (s *AthenaService) BatchGetPreparedStatement(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	preparedStatementsRaw := request.GetArrayParam(req.Parameters, "PreparedStatementNames")
	if len(preparedStatementsRaw) == 0 {
		return nil, ErrInvalidRequestException
	}

	if len(preparedStatementsRaw) > 50 {
		return nil, ErrInvalidRequestException
	}

	workGroup := request.GetParamCaseInsensitive(req.Parameters, "WorkGroup")
	if workGroup == "" {
		workGroup = "primary"
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var preparedStatements []map[string]interface{}
	var unprocessedNames []map[string]interface{}

	for _, nameRaw := range preparedStatementsRaw {
		name, ok := nameRaw.(string)
		if !ok {
			continue
		}

		ps, err := stores.preparedStatementStore.GetPreparedStatement(workGroup, name)
		if err != nil {
			unprocessedNames = append(unprocessedNames, map[string]interface{}{
				"StatementName": name,
			})
			continue
		}
		preparedStatements = append(preparedStatements, s.preparedStatementToResponse(ps))
	}

	return map[string]interface{}{
		"PreparedStatements":                preparedStatements,
		"UnprocessedPreparedStatementNames": unprocessedNames,
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

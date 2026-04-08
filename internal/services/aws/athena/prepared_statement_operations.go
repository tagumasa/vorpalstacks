package athena

import (
	"context"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	athenastore "vorpalstacks/internal/store/aws/athena"
)

// CreatePreparedStatement creates a new prepared statement in the Athena workgroup.
func (s *Service) CreatePreparedStatement(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	statementName := request.GetParamCaseInsensitive(req.Parameters, "StatementName")
	if statementName == "" {
		return nil, ErrInvalidRequestException
	}

	if len(statementName) > 256 {
		return nil, ErrInvalidParameterException
	}

	workGroup := request.GetParamCaseInsensitive(req.Parameters, "WorkGroup")
	if workGroup == "" {
		workGroup = "primary"
	}

	description := request.GetParamCaseInsensitive(req.Parameters, "Description")
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
			return nil, ErrInvalidRequestException
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetPreparedStatement retrieves the specified prepared statement.
func (s *Service) GetPreparedStatement(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
			return nil, ErrInvalidRequestException
		}
		return nil, err
	}

	return map[string]interface{}{
		"PreparedStatement": s.preparedStatementToResponse(preparedStatement),
	}, nil
}

// DeletePreparedStatement deletes the specified prepared statement.
func (s *Service) DeletePreparedStatement(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
			return nil, ErrInvalidRequestException
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListPreparedStatements retrieves a list of prepared statements in the specified workgroup.
func (s *Service) ListPreparedStatements(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	return map[string]interface{}{
		"PreparedStatements": summaries,
		"NextToken":          "",
	}, nil
}

// UpdatePreparedStatement updates the specified prepared statement.
func (s *Service) UpdatePreparedStatement(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
			return nil, ErrResourceNotFoundException
		}
		return nil, err
	}

	description := request.GetParamCaseInsensitive(req.Parameters, "Description")
	if description != "" {
		preparedStatement.Description = description
	}

	queryStatement := request.GetParamCaseInsensitive(req.Parameters, "QueryStatement")
	if queryStatement != "" {
		preparedStatement.QueryStatement = queryStatement
	}

	if err := stores.preparedStatementStore.UpdatePreparedStatement(preparedStatement); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// BatchGetPreparedStatement retrieves multiple prepared statements in a single request.
func (s *Service) BatchGetPreparedStatement(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

func (s *Service) preparedStatementToResponse(ps *athenastore.PreparedStatement) map[string]interface{} {
	return map[string]interface{}{
		"StatementName":    ps.StatementName,
		"QueryStatement":   ps.QueryStatement,
		"WorkGroupName":    ps.WorkGroupName,
		"Description":      ps.Description,
		"LastModifiedTime": float64(ps.LastModifiedTime.UnixNano()) / 1e9,
	}
}

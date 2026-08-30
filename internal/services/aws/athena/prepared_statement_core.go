package athena

import (
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	athenastore "vorpalstacks/internal/store/aws/athena"
)

// --- DTOs ---

// CreatePreparedStatementInput carries the parsed wire members of a
// CreatePreparedStatement request.
type CreatePreparedStatementInput struct {
	StatementName  string
	WorkGroup      string
	Description    string
	QueryStatement string
}

// GetPreparedStatementInput carries the statement identity; the workgroup
// defaults to "primary" inside the Core exactly as the handler did.
type GetPreparedStatementInput struct {
	StatementName string
	WorkGroup     string
}

// ListPreparedStatementsInput carries the workgroup filter plus the raw
// MaxResults window (presence-flagged) and pagination marker.
type ListPreparedStatementsInput struct {
	WorkGroup     string
	MaxResults    int
	HasMaxResults bool
	NextToken     string
}

// UpdatePreparedStatementInput carries the parsed wire members of an
// UpdatePreparedStatement request.
type UpdatePreparedStatementInput struct {
	StatementName  string
	WorkGroup      string
	QueryStatement string
	Description    string
}

// BatchGetPreparedStatementInput carries the workgroup plus the raw
// PreparedStatementNames wire array; the Core applies the entry-count gates
// and string filtering.
type BatchGetPreparedStatementInput struct {
	WorkGroup              string
	PreparedStatementNames []interface{}
}

// --- Core functions ---

// createPreparedStatementCore validates the create request and persists the
// prepared statement. The store is acquired only after the request has been
// validated, the order the original handler applied.
func (s *AthenaService) createPreparedStatementCore(reqCtx *request.RequestContext, input CreatePreparedStatementInput) error {
	if input.StatementName == "" {
		return ErrInvalidRequestException
	}

	if err := validateStatementName(input.StatementName); err != nil {
		return err
	}

	workGroup := input.WorkGroup
	if workGroup == "" {
		workGroup = "primary"
	}

	if input.Description != "" {
		if err := validateDescriptionString(input.Description); err != nil {
			return err
		}
	}

	if err := validateQueryStringSize(input.QueryStatement); err != nil {
		return err
	}

	if input.QueryStatement == "" {
		return ErrInvalidRequestException
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	preparedStatement := &athenastore.PreparedStatement{
		StatementName:  input.StatementName,
		WorkGroupName:  workGroup,
		QueryStatement: input.QueryStatement,
		Description:    input.Description,
	}

	if err := stores.preparedStatementStore.CreatePreparedStatement(preparedStatement); err != nil {
		if err == athenastore.ErrPreparedStatementAlreadyExists {
			return alreadyExistsInvalidRequest("PreparedStatement", input.StatementName)
		}
		return err
	}

	return nil
}

// getPreparedStatementCore fetches a prepared statement scoped to its
// workgroup, mapping the store not-found sentinel onto the API error.
func (s *AthenaService) getPreparedStatementCore(reqCtx *request.RequestContext, input GetPreparedStatementInput) (*athenastore.PreparedStatement, error) {
	if input.StatementName == "" {
		return nil, ErrInvalidRequestException
	}

	workGroup := input.WorkGroup
	if workGroup == "" {
		workGroup = "primary"
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	preparedStatement, err := stores.preparedStatementStore.GetPreparedStatement(workGroup, input.StatementName)
	if err != nil {
		if err == athenastore.ErrPreparedStatementNotFound {
			return nil, preparedStatementNotFound(input.StatementName)
		}
		return nil, err
	}
	return preparedStatement, nil
}

// deletePreparedStatementCore deletes a prepared statement scoped to its
// workgroup, mapping the store not-found sentinel onto the API error.
func (s *AthenaService) deletePreparedStatementCore(reqCtx *request.RequestContext, input GetPreparedStatementInput) error {
	if input.StatementName == "" {
		return ErrInvalidRequestException
	}

	workGroup := input.WorkGroup
	if workGroup == "" {
		workGroup = "primary"
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	if err := stores.preparedStatementStore.DeletePreparedStatement(workGroup, input.StatementName); err != nil {
		if err == athenastore.ErrPreparedStatementNotFound {
			return preparedStatementNotFound(input.StatementName)
		}
		return err
	}
	return nil
}

// listPreparedStatementsCore lists the workgroup's prepared statements and
// pages them by statement name with the documented window semantics (default
// 50, range 1-50) applied after the list walk, matching the original
// validation position.
func listPreparedStatementsCore(stores *athenaStores, input ListPreparedStatementsInput) ([]*athenastore.PreparedStatement, string, error) {
	workGroup := input.WorkGroup
	if workGroup == "" {
		workGroup = "primary"
	}

	preparedStatements, err := stores.preparedStatementStore.ListPreparedStatements(workGroup)
	if err != nil {
		return nil, "", err
	}

	maxResults, err := resolveMaxResults(input.MaxResults, input.HasMaxResults, 50, 1, 50)
	if err != nil {
		return nil, "", err
	}

	pageResult := pagination.PaginateSlice(preparedStatements, input.NextToken, maxResults, func(item *athenastore.PreparedStatement) string {
		return item.StatementName
	})

	return pageResult.Items, pageResult.NextMarker, nil
}

// updatePreparedStatementCore validates the update request (per the Smithy
// model, StatementName, WorkGroup and QueryStatement are all REQUIRED),
// applies it to the stored record and persists it.
func (s *AthenaService) updatePreparedStatementCore(reqCtx *request.RequestContext, input UpdatePreparedStatementInput) error {
	if input.StatementName == "" {
		return ErrInvalidRequestException
	}

	workGroup := input.WorkGroup
	if workGroup == "" {
		workGroup = "primary"
	}

	if input.QueryStatement == "" {
		return invalidRequestParameter("QueryStatement is required for UpdatePreparedStatement")
	}

	if err := validateQueryStringSize(input.QueryStatement); err != nil {
		return err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	preparedStatement, err := stores.preparedStatementStore.GetPreparedStatement(workGroup, input.StatementName)
	if err != nil {
		if err == athenastore.ErrPreparedStatementNotFound {
			return preparedStatementNotFound(input.StatementName)
		}
		return err
	}

	preparedStatement.QueryStatement = input.QueryStatement

	if input.Description != "" {
		preparedStatement.Description = input.Description
	}

	if err := stores.preparedStatementStore.UpdatePreparedStatement(preparedStatement); err != nil {
		if err == athenastore.ErrPreparedStatementNotFound {
			return preparedStatementNotFound(input.StatementName)
		}
		return err
	}

	return nil
}

// batchGetPreparedStatementCore applies the entry-count gates (1-50 raw
// names), fetches each prepared statement and partitions found/unprocessed.
func (s *AthenaService) batchGetPreparedStatementCore(reqCtx *request.RequestContext, input BatchGetPreparedStatementInput) ([]*athenastore.PreparedStatement, []string, error) {
	if len(input.PreparedStatementNames) == 0 {
		return nil, nil, ErrInvalidRequestException
	}

	if len(input.PreparedStatementNames) > 50 {
		return nil, nil, ErrInvalidRequestException
	}

	workGroup := input.WorkGroup
	if workGroup == "" {
		workGroup = "primary"
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, nil, err
	}

	var found []*athenastore.PreparedStatement
	var unprocessed []string

	for _, nameRaw := range input.PreparedStatementNames {
		name, ok := nameRaw.(string)
		if !ok {
			continue
		}

		ps, err := stores.preparedStatementStore.GetPreparedStatement(workGroup, name)
		if err != nil {
			unprocessed = append(unprocessed, name)
			continue
		}
		found = append(found, ps)
	}

	return found, unprocessed, nil
}

package athena

import (
	"sort"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	athenastore "vorpalstacks/internal/store/aws/athena"
)

// --- DTOs ---

// CreateNamedQueryInput carries the parsed wire members of a
// CreateNamedQuery request.
type CreateNamedQueryInput struct {
	Name        string
	Description string
	Database    string
	QueryString string
	WorkGroup   string
}

// ListNamedQueriesInput carries the workgroup filter plus the raw MaxResults
// window (presence-flagged) and pagination marker.
type ListNamedQueriesInput struct {
	WorkGroup     string
	MaxResults    int
	HasMaxResults bool
	NextToken     string
}

// UpdateNamedQueryInput carries the parsed wire members of an
// UpdateNamedQuery request.
type UpdateNamedQueryInput struct {
	NamedQueryId string
	Name         string
	Description  string
	QueryString  string
}

// BatchGetNamedQueryInput carries the raw NamedQueryIds wire array; the Core
// applies the entry-count gates and string filtering.
type BatchGetNamedQueryInput struct {
	NamedQueryIds []interface{}
}

// --- Core functions ---

// createNamedQueryCore validates the create request and persists the named
// query, returning its minted id. The store is acquired only after the
// request has been validated, the order the original handler applied.
func (s *AthenaService) createNamedQueryCore(reqCtx *request.RequestContext, input CreateNamedQueryInput) (string, error) {
	if input.Name == "" {
		return "", ErrInvalidRequestException
	}

	if err := validateNameString(input.Name); err != nil {
		return "", err
	}

	if input.Description != "" {
		if err := validateNamedQueryDescriptionString(input.Description); err != nil {
			return "", err
		}
	}

	if err := validateQueryStringSize(input.QueryString); err != nil {
		return "", err
	}

	workGroup := input.WorkGroup
	if workGroup == "" {
		workGroup = "primary"
	}

	if input.Database == "" || input.QueryString == "" {
		return "", ErrInvalidRequestException
	}

	if err := validateDatabaseString(input.Database); err != nil {
		return "", err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return "", err
	}

	namedQuery := &athenastore.NamedQuery{
		Name:        input.Name,
		Description: input.Description,
		Database:    input.Database,
		QueryString: input.QueryString,
		WorkGroup:   workGroup,
	}

	if err := stores.namedQueryStore.CreateNamedQuery(namedQuery); err != nil {
		if err == athenastore.ErrNamedQueryAlreadyExists {
			return "", alreadyExistsInvalidRequest("NamedQuery", input.Name)
		}
		return "", err
	}

	return namedQuery.NamedQueryId, nil
}

// getNamedQueryCore fetches a named query, mapping the store not-found
// sentinel onto the API error. The store is acquired after the identifier
// validation, the order the original handler applied.
func (s *AthenaService) getNamedQueryCore(reqCtx *request.RequestContext, namedQueryId string) (*athenastore.NamedQuery, error) {
	if namedQueryId == "" {
		return nil, ErrInvalidRequestException
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	namedQuery, err := stores.namedQueryStore.GetNamedQuery(namedQueryId)
	if err != nil {
		if err == athenastore.ErrNamedQueryNotFound {
			return nil, namedQueryNotFound(namedQueryId)
		}
		return nil, err
	}
	return namedQuery, nil
}

// deleteNamedQueryCore deletes a named query, mapping the store not-found
// sentinel onto the API error. The store is acquired after the identifier
// validation, the order the original handler applied.
func (s *AthenaService) deleteNamedQueryCore(reqCtx *request.RequestContext, namedQueryId string) error {
	if namedQueryId == "" {
		return ErrInvalidRequestException
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	if err := stores.namedQueryStore.DeleteNamedQuery(namedQueryId); err != nil {
		if err == athenastore.ErrNamedQueryNotFound {
			return namedQueryNotFound(namedQueryId)
		}
		return err
	}
	return nil
}

// listNamedQueriesCore lists the workgroup's named-query ids (sorted) and
// pages them with the documented window semantics (default 50, range 0-50)
// applied after the list walk, matching the original validation position.
func listNamedQueriesCore(stores *athenaStores, input ListNamedQueriesInput) ([]string, string, error) {
	workGroup := input.WorkGroup
	if workGroup == "" {
		workGroup = "primary"
	}

	namedQueries, err := stores.namedQueryStore.ListNamedQueries(workGroup)
	if err != nil {
		return nil, "", err
	}

	ids := make([]string, 0, len(namedQueries))
	for _, nq := range namedQueries {
		ids = append(ids, nq.NamedQueryId)
	}
	sort.Strings(ids)

	maxResults, err := resolveMaxResults(input.MaxResults, input.HasMaxResults, 50, 0, 50)
	if err != nil {
		return nil, "", err
	}

	pageResult := pagination.PaginateSlice(ids, input.NextToken, maxResults, func(id string) string {
		return id
	})

	return pageResult.Items, pageResult.NextMarker, nil
}

// updateNamedQueryCore validates the update request (per the Smithy model,
// NamedQueryId, Name and QueryString are all REQUIRED), applies it to the
// stored record and persists it.
func (s *AthenaService) updateNamedQueryCore(reqCtx *request.RequestContext, input UpdateNamedQueryInput) error {
	if input.NamedQueryId == "" {
		return ErrInvalidRequestException
	}

	if input.Name == "" {
		return invalidRequestParameter("Name is required for UpdateNamedQuery")
	}
	if err := validateNameString(input.Name); err != nil {
		return err
	}

	if input.QueryString == "" {
		return invalidRequestParameter("QueryString is required for UpdateNamedQuery")
	}
	if err := validateQueryStringSize(input.QueryString); err != nil {
		return err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	namedQuery, err := stores.namedQueryStore.GetNamedQuery(input.NamedQueryId)
	if err != nil {
		if err == athenastore.ErrNamedQueryNotFound {
			return namedQueryNotFound(input.NamedQueryId)
		}
		return err
	}

	namedQuery.Name = input.Name
	namedQuery.QueryString = input.QueryString

	if input.Description != "" {
		if err := validateNamedQueryDescriptionString(input.Description); err != nil {
			return err
		}
		namedQuery.Description = input.Description
	}

	if err := stores.namedQueryStore.UpdateNamedQuery(input.NamedQueryId, namedQuery); err != nil {
		if err == athenastore.ErrNamedQueryNotFound {
			return namedQueryNotFound(input.NamedQueryId)
		}
		return err
	}

	return nil
}

// batchGetNamedQueryCore applies the entry-count gates (1-50 raw ids),
// fetches each named query and partitions found/unprocessed. The store is
// acquired only after the gates, the order the original handler applied.
func (s *AthenaService) batchGetNamedQueryCore(reqCtx *request.RequestContext, input BatchGetNamedQueryInput) ([]*athenastore.NamedQuery, []string, error) {
	if len(input.NamedQueryIds) == 0 {
		return nil, nil, ErrInvalidRequestException
	}

	if len(input.NamedQueryIds) > 50 {
		return nil, nil, ErrInvalidRequestException
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, nil, err
	}

	var namedQueryIds []string
	for _, id := range input.NamedQueryIds {
		if idStr, ok := id.(string); ok {
			namedQueryIds = append(namedQueryIds, idStr)
		}
	}

	var found []*athenastore.NamedQuery
	var unprocessed []string

	for _, id := range namedQueryIds {
		namedQuery, err := stores.namedQueryStore.GetNamedQuery(id)
		if err != nil {
			unprocessed = append(unprocessed, id)
			continue
		}
		found = append(found, namedQuery)
	}

	return found, unprocessed, nil
}

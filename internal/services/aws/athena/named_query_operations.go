package athena

import (
	"context"
	"strconv"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	athenastore "vorpalstacks/internal/store/aws/athena"
)

// CreateNamedQuery creates a new named query in the Athena workgroup.
func (s *Service) CreateNamedQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "Name")
	if name == "" {
		return nil, ErrInvalidRequestException
	}

	if len(name) > 128 {
		return nil, ErrInvalidParameterException
	}

	description := request.GetParamCaseInsensitive(req.Parameters, "Description")
	database := request.GetParamCaseInsensitive(req.Parameters, "Database")
	queryString := request.GetParamCaseInsensitive(req.Parameters, "QueryString")

	if len(queryString) > maxQueryStringSize {
		return nil, ErrInvalidRequestException
	}

	workGroup := request.GetParamCaseInsensitive(req.Parameters, "WorkGroup")

	if workGroup == "" {
		workGroup = "primary"
	}

	if database == "" || queryString == "" {
		return nil, ErrInvalidRequestException
	}

	namedQuery := &athenastore.NamedQuery{
		Name:        name,
		Description: description,
		Database:    database,
		QueryString: queryString,
		WorkGroup:   workGroup,
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := stores.namedQueryStore.CreateNamedQuery(namedQuery); err != nil {
		if err == athenastore.ErrNamedQueryAlreadyExists {
			return nil, ErrResourceAlreadyExistsException
		}
		return nil, err
	}

	return map[string]interface{}{
		"NamedQueryId": namedQuery.NamedQueryId,
	}, nil
}

// GetNamedQuery retrieves the specified named query.
func (s *Service) GetNamedQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	namedQueryId := request.GetParamCaseInsensitive(req.Parameters, "NamedQueryId")
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
			return nil, ErrResourceNotFoundException
		}
		return nil, err
	}

	return map[string]interface{}{
		"NamedQuery": s.namedQueryToResponse(namedQuery),
	}, nil
}

// DeleteNamedQuery deletes the specified named query.
func (s *Service) DeleteNamedQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	namedQueryId := request.GetParamCaseInsensitive(req.Parameters, "NamedQueryId")
	if namedQueryId == "" {
		return nil, ErrInvalidRequestException
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := stores.namedQueryStore.DeleteNamedQuery(namedQueryId); err != nil {
		if err == athenastore.ErrNamedQueryNotFound {
			return nil, ErrResourceNotFoundException
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListNamedQueries retrieves a list of named queries in the specified workgroup.
func (s *Service) ListNamedQueries(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	workGroup := request.GetParamCaseInsensitive(req.Parameters, "WorkGroup")
	if workGroup == "" {
		workGroup = "primary"
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	namedQueries, err := stores.namedQueryStore.ListNamedQueries(workGroup)
	if err != nil {
		return nil, err
	}

	maxResults := 50
	if maxStr := request.GetParamCaseInsensitive(req.Parameters, "MaxResults"); maxStr != "" {
		if val, err := strconv.Atoi(maxStr); err == nil && val > 0 {
			maxResults = val
		}
	}

	offset := 0
	if nextToken := request.GetParamCaseInsensitive(req.Parameters, "NextToken"); nextToken != "" {
		if val, err := strconv.Atoi(nextToken); err == nil && val >= 0 {
			offset = val
		}
	}

	var ids []string
	for i, nq := range namedQueries {
		if i < offset {
			continue
		}
		if len(ids) >= maxResults {
			break
		}
		ids = append(ids, nq.NamedQueryId)
	}

	result := map[string]interface{}{
		"NamedQueryIds": ids,
	}

	if offset+maxResults < len(namedQueries) {
		result["NextToken"] = strconv.Itoa(offset + maxResults)
	}

	return result, nil
}

// UpdateNamedQuery updates the specified named query.
func (s *Service) UpdateNamedQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	namedQueryId := request.GetParamCaseInsensitive(req.Parameters, "NamedQueryId")
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
			return nil, ErrResourceNotFoundException
		}
		return nil, err
	}

	name := request.GetParamCaseInsensitive(req.Parameters, "Name")
	if name != "" {
		namedQuery.Name = name
	}

	description := request.GetParamCaseInsensitive(req.Parameters, "Description")
	if description != "" {
		namedQuery.Description = description
	}

	queryString := request.GetParamCaseInsensitive(req.Parameters, "QueryString")
	if queryString != "" {
		namedQuery.QueryString = queryString
	}

	if err := stores.namedQueryStore.UpdateNamedQuery(namedQueryId, namedQuery); err != nil {
		if err == athenastore.ErrNamedQueryNotFound {
			return nil, ErrResourceNotFoundException
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// BatchGetNamedQuery retrieves multiple named queries in a single request.
func (s *Service) BatchGetNamedQuery(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	namedQueryIdsRaw := request.GetArrayParam(req.Parameters, "NamedQueryIds")
	if len(namedQueryIdsRaw) == 0 {
		return nil, ErrInvalidRequestException
	}

	if len(namedQueryIdsRaw) > 50 {
		return nil, ErrInvalidRequestException
	}

	var namedQueryIds []string
	for _, id := range namedQueryIdsRaw {
		if idStr, ok := id.(string); ok {
			namedQueryIds = append(namedQueryIds, idStr)
		}
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	var namedQueries []map[string]interface{}
	var unprocessedIds []map[string]interface{}

	for _, id := range namedQueryIds {
		namedQuery, err := stores.namedQueryStore.GetNamedQuery(id)
		if err != nil {
			unprocessedIds = append(unprocessedIds, map[string]interface{}{
				"NamedQueryId": id,
			})
			continue
		}
		namedQueries = append(namedQueries, s.namedQueryToResponse(namedQuery))
	}

	return map[string]interface{}{
		"NamedQueries":             namedQueries,
		"UnprocessedNamedQueryIds": unprocessedIds,
	}, nil
}

func (s *Service) namedQueryToResponse(nq *athenastore.NamedQuery) map[string]interface{} {
	return map[string]interface{}{
		"Name":         nq.Name,
		"Description":  nq.Description,
		"Database":     nq.Database,
		"QueryString":  nq.QueryString,
		"NamedQueryId": nq.NamedQueryId,
		"WorkGroup":    nq.WorkGroup,
	}
}

package cloudwatchlogs

import (
	"context"
	"fmt"
	"time"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// PutQueryDefinition creates or updates a saved query definition.
func (s *LogsService) PutQueryDefinition(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	queryString := request.GetParamLowerFirst(req.Parameters, "QueryString")
	queryLanguage := request.GetParamLowerFirst(req.Parameters, "QueryLanguage")

	if err := validateQueryDefinitionName(name); err != nil {
		return nil, err
	}
	if err := validateQueryString(queryString); err != nil {
		return nil, err
	}

	queryDefinitionId := request.GetParamLowerFirst(req.Parameters, "QueryDefinitionId")

	var parameters map[string]interface{}
	if p, ok := req.Parameters["parameters"]; ok {
		if m, ok := p.(map[string]interface{}); ok {
			parameters = m
		}
	}

	var logGroupNames []string
	if names, ok := req.Parameters["logGroupNames"]; ok {
		if arr, ok := names.([]interface{}); ok {
			for _, item := range arr {
				if s, ok := item.(string); ok {
					logGroupNames = append(logGroupNames, s)
				}
			}
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if queryDefinitionId == "" {
		queryDefinitionId = generateQueryDefinitionId()
	}

	qd := &logsstore.QueryDefinition{
		QueryDefinitionId: queryDefinitionId,
		Name:              name,
		QueryString:       queryString,
		LogGroupNames:     logGroupNames,
		QueryLanguage:     queryLanguage,
		Parameters:        parameters,
		LastModified:      time.Now().UTC().UnixMilli(),
	}

	if err := store.PutQueryDefinitionEntry(qd); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"queryDefinitionId": queryDefinitionId,
	}, nil
}

// DeleteQueryDefinition deletes a saved query definition.
func (s *LogsService) DeleteQueryDefinition(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryDefinitionId := request.GetParamLowerFirst(req.Parameters, "QueryDefinitionId")
	if queryDefinitionId == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteQueryDefinitionEntry(queryDefinitionId); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"success": true,
	}, nil
}

// DescribeQueryDefinitions lists saved query definitions.
func (s *LogsService) DescribeQueryDefinitions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	namePrefix := request.GetParamLowerFirst(req.Parameters, "QueryDefinitionNamePrefix")
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")
	maxResults, err := validateListLimit(int32(request.GetIntParam(req.Parameters, "MaxResults")), 1000, 1000)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	allDefs, err := store.ListQueryDefinitions(namePrefix)
	if err != nil {
		return nil, mapStoreError(err)
	}

	result := pagination.PaginateSlice(allDefs, nextToken, int(maxResults), func(qd *logsstore.QueryDefinition) string {
		return qd.QueryDefinitionId
	})

	defs := make([]map[string]interface{}, len(result.Items))
	for i, qd := range result.Items {
		entry := map[string]interface{}{
			"queryDefinitionId": qd.QueryDefinitionId,
			"name":              qd.Name,
			"queryString":       qd.QueryString,
			"lastModified":      qd.LastModified,
		}
		if len(qd.LogGroupNames) > 0 {
			entry["logGroupNames"] = qd.LogGroupNames
		}
		if qd.QueryLanguage != "" {
			entry["queryLanguage"] = qd.QueryLanguage
		}
		if qd.Parameters != nil {
			entry["parameters"] = qd.Parameters
		}
		defs[i] = entry
	}

	resp := map[string]interface{}{
		"queryDefinitions": defs,
	}
	if result.NextMarker != "" {
		resp["nextToken"] = result.NextMarker
	}

	return resp, nil
}

func generateQueryDefinitionId() string {
	return fmt.Sprintf("qd-%d", time.Now().UnixNano())
}

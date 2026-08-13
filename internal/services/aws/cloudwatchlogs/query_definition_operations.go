package cloudwatchlogs

import (
	"context"
	"fmt"
	"time"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// --- Core methods ---

func (s *LogsService) putQueryDefinitionCore(name, queryString, queryLanguage, queryDefinitionId string, logGroupNames []string, parameters map[string]interface{}, region string) (string, error) {
	if err := validateQueryDefinitionName(name); err != nil {
		return "", err
	}
	if err := validateQueryString(queryString); err != nil {
		return "", err
	}
	if !validateQueryLanguage(queryLanguage) {
		return "", NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid queryLanguage: %s. Allowed values: CWLI, SQL, PPL", queryLanguage), 400)
	}

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return "", err
	}

	if queryDefinitionId == "" {
		queryDefinitionId = fmt.Sprintf("qd-%d", time.Now().UnixNano())
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
		return "", mapStoreError(err)
	}
	return queryDefinitionId, nil
}

func (s *LogsService) deleteQueryDefinitionCore(queryDefinitionId, region string) error {
	if queryDefinitionId == "" {
		return ErrMissingParameter
	}

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return err
	}

	if err := store.DeleteQueryDefinitionEntry(queryDefinitionId); err != nil {
		return mapStoreError(err)
	}
	return nil
}

func (s *LogsService) describeQueryDefinitionsCore(namePrefix, nextToken, region string, maxResults int32) ([]*logsstore.QueryDefinition, string, error) {
	limit, err := validateListLimit(maxResults, 1000, 1000)
	if err != nil {
		return nil, "", err
	}

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return nil, "", err
	}

	allDefs, err := store.ListQueryDefinitions(namePrefix)
	if err != nil {
		return nil, "", mapStoreError(err)
	}

	result := pagination.PaginateSlice(allDefs, nextToken, int(limit), func(qd *logsstore.QueryDefinition) string {
		return qd.QueryDefinitionId
	})

	return result.Items, result.NextMarker, nil
}

// --- HTTP handlers ---

func (s *LogsService) PutQueryDefinition(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	queryString := request.GetParamLowerFirst(req.Parameters, "QueryString")
	queryLanguage := request.GetParamLowerFirst(req.Parameters, "QueryLanguage")
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

	id, err := s.putQueryDefinitionCore(name, queryString, queryLanguage, queryDefinitionId, logGroupNames, parameters, reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"queryDefinitionId": id,
	}, nil
}

func (s *LogsService) DeleteQueryDefinition(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryDefinitionId := request.GetParamLowerFirst(req.Parameters, "QueryDefinitionId")

	if err := s.deleteQueryDefinitionCore(queryDefinitionId, reqCtx.GetRegion()); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
	}, nil
}

func (s *LogsService) DescribeQueryDefinitions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	namePrefix := request.GetParamLowerFirst(req.Parameters, "QueryDefinitionNamePrefix")
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")
	maxResults := int32(request.GetIntParam(req.Parameters, "MaxResults"))

	defs, nextMarker, err := s.describeQueryDefinitionsCore(namePrefix, nextToken, reqCtx.GetRegion(), maxResults)
	if err != nil {
		return nil, err
	}

	formatted := make([]map[string]interface{}, len(defs))
	for i, qd := range defs {
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
		formatted[i] = entry
	}

	resp := map[string]interface{}{
		"queryDefinitions": formatted,
	}
	if nextMarker != "" {
		resp["nextToken"] = nextMarker
	}

	return resp, nil
}

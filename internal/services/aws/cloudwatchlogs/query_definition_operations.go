package cloudwatchlogs

import (
	"context"
	"fmt"
	"time"
	"unicode/utf8"

	"vorpalstacks/internal/common/request"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// --- Core methods ---

func (s *LogsService) putQueryDefinitionCore(name, queryString, queryLanguage, queryDefinitionId string, logGroupNames []string, parameters []logsstore.QueryParameter, region string) (string, error) {
	if err := validateQueryDefinitionName(name); err != nil {
		return "", err
	}
	if err := validateQueryString(queryString, 1); err != nil {
		return "", err
	}
	if !validateQueryLanguage(queryLanguage) {
		return "", NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid queryLanguage: %s. Allowed values: CWLI, SQL, PPL", queryLanguage), 400)
	}
	if err := validateQueryParameters(parameters); err != nil {
		return "", err
	}

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return "", err
	}

	if queryDefinitionId == "" {
		queryDefinitionId = fmt.Sprintf("qd-%d", time.Now().UnixNano())
	} else {
		if utf8.RuneCountInString(queryDefinitionId) > logsstore.MaxQueryDefinitionIdLength {
			return "", NewLogsError("InvalidParameterException",
				fmt.Sprintf("queryDefinitionId must be between 1 and %d characters", logsstore.MaxQueryDefinitionIdLength), 400)
		}
		// The update form names the definition it updates ("To update a
		// query definition, specify its queryDefinitionId in your
		// request"): an unknown id is the operation's declared
		// ResourceNotFoundException, never a silent create.
		if _, err := store.GetQueryDefinitionEntry(queryDefinitionId); err != nil {
			return "", mapStoreError(err)
		}
	}

	qd := &logsstore.QueryDefinition{
		QueryDefinitionId: queryDefinitionId,
		Name:              name,
		QueryString:       queryString,
		LogGroupNames:     logGroupNames,
		QueryLanguage:     queryLanguage,
		Parameters:        parameters,
	}

	if err := store.PutQueryDefinitionEntry(qd); err != nil {
		return "", mapStoreError(err)
	}
	return queryDefinitionId, nil
}

func (s *LogsService) deleteQueryDefinitionCore(queryDefinitionId, region string) error {
	if queryDefinitionId == "" {
		return errRequiredMember("queryDefinitionId")
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

// describeQueryDefinitionsCore lists saved queries under the optional
// name prefix and query-language filter. An absent stored language is the
// CWLI default, so a CWLI filter still matches definitions saved without
// an explicit language — the same normalisation the DescribeQueries
// filter applies.
func (s *LogsService) describeQueryDefinitionsCore(namePrefix, queryLanguage, nextToken, region string, maxResults int32) ([]*logsstore.QueryDefinition, string, error) {
	limit, err := validateListLimit(maxResults, logsstore.MaxListMaxResults, logsstore.MaxListMaxResults)
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

	filtered := make([]*logsstore.QueryDefinition, 0, len(allDefs))
	for _, qd := range allDefs {
		if queryLanguage != "" {
			stored := qd.QueryLanguage
			if stored == "" {
				stored = "CWLI"
			}
			if stored != queryLanguage {
				continue
			}
		}
		filtered = append(filtered, qd)
	}

	// The nextToken is the scoped listing vocabulary's form ("The token
	// expires after 24 hours"): the raw key stays inside the token, so a
	// token minted by another filter shape or past its expiry rejects.
	scope := listingScope("query-definitions", namePrefix, queryLanguage, region)
	result, err := paginateScopedListing(scope, nextToken, filtered, int(limit), func(qd *logsstore.QueryDefinition) string {
		return qd.QueryDefinitionId
	})
	if err != nil {
		return nil, "", err
	}
	return result.Items, result.NextMarker, nil
}

// --- HTTP handlers ---

func (s *LogsService) PutQueryDefinition(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamLowerFirst(req.Parameters, "Name")
	queryString := request.GetParamLowerFirst(req.Parameters, "QueryString")
	queryLanguage := request.GetParamLowerFirst(req.Parameters, "QueryLanguage")
	queryDefinitionId := request.GetParamLowerFirst(req.Parameters, "QueryDefinitionId")

	parameters := parseQueryParameterList(req.Parameters)

	logGroupNames := request.GetStringList(req.Parameters, "LogGroupNames")

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
	queryLanguage := request.GetParamLowerFirst(req.Parameters, "QueryLanguage")
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")
	maxResults := int32(request.GetIntParam(req.Parameters, "MaxResults"))

	defs, nextMarker, err := s.describeQueryDefinitionsCore(namePrefix, queryLanguage, nextToken, reqCtx.GetRegion(), maxResults)
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
		if len(qd.Parameters) > 0 {
			entry["parameters"] = formatQueryParameters(qd.Parameters)
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

// parseQueryParameterList reads the wire form of a saved query's
// parameters — a list of name/defaultValue/description members — into the
// store type; validation of the parsed values lives in Core.
func parseQueryParameterList(params map[string]interface{}) []logsstore.QueryParameter {
	raw := request.GetArrayParam(params, "parameters")
	if raw == nil {
		return nil
	}
	out := make([]logsstore.QueryParameter, 0, len(raw))
	for _, item := range raw {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		qp := logsstore.QueryParameter{
			Name:         request.GetStringParam(m, "name"),
			DefaultValue: request.GetStringParam(m, "defaultValue"),
			Description:  request.GetStringParam(m, "description"),
		}
		out = append(out, qp)
	}
	return out
}

// formatQueryParameters renders the parameter list back to the wire form:
// every member present on a parameter appears, absent ones stay absent.
func formatQueryParameters(params []logsstore.QueryParameter) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(params))
	for _, p := range params {
		entry := map[string]interface{}{"name": p.Name}
		if p.DefaultValue != "" {
			entry["defaultValue"] = p.DefaultValue
		}
		if p.Description != "" {
			entry["description"] = p.Description
		}
		out = append(out, entry)
	}
	return out
}

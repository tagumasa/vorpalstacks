package cloudwatchlogs

import (
	"context"
	"encoding/csv"
	"fmt"
	"regexp"
	"strings"

	"vorpalstacks/internal/common/request"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// lookupTableNameRe validates the documented name pattern of lookup tables
// (alphanumeric characters and underscores).
var lookupTableNameRe = regexp.MustCompile(logsstore.LookupTableNamePattern)

// lookupTableArn builds the ARN of a lookup table.
func lookupTableArn(region, accountID, name string) string {
	return fmt.Sprintf("arn:aws:logs:%s:%s:lookup-table:%s", region, accountID, name)
}

// lookupTableNameFromArn extracts the table name from an ARN; a bare name
// is accepted as-is.
func lookupTableNameFromArn(identifier string) string {
	if strings.HasPrefix(identifier, "arn:") {
		if idx := strings.LastIndex(identifier, ":"); idx >= 0 {
			return identifier[idx+1:]
		}
	}
	return identifier
}

// stringParamPresent extracts a string parameter with presence semantics
// across the case variants, distinguishing "absent" from "empty".
func stringParamPresent(params map[string]interface{}, key string) (string, bool) {
	for _, k := range []string{key, request.LowerFirst(key), strings.ToLower(key)} {
		if v, ok := params[k]; ok {
			if str, ok := v.(string); ok {
				return str, true
			}
		}
	}
	return "", false
}

// parseLookupCSV parses a lookup table body: the first record is the header
// row and must be non-empty; ragged data rows are padded with empty values
// when shorter and rejected when longer than the header.
func parseLookupCSV(body string) (fields []string, records [][]string, err error) {
	reader := csv.NewReader(strings.NewReader(body))
	reader.FieldsPerRecord = -1
	all, readErr := reader.ReadAll()
	if readErr != nil {
		return nil, nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid tableBody CSV: %v", readErr), 400)
	}
	if len(all) == 0 || len(all[0]) == 0 {
		return nil, nil, NewLogsError("InvalidParameterException",
			"tableBody CSV must include a header row with column names", 400)
	}
	fields = all[0]
	for _, rec := range all[1:] {
		if len(rec) > len(fields) {
			return nil, nil, NewLogsError("InvalidParameterException",
				fmt.Sprintf("tableBody CSV row has %d values but the header has %d columns", len(rec), len(fields)), 400)
		}
		if len(rec) < len(fields) {
			padded := make([]string, len(fields))
			copy(padded, rec)
			rec = padded
		}
		records = append(records, rec)
	}
	return fields, records, nil
}

// CreateLookupTable creates a lookup table from CSV data or query results.
func (s *LogsService) CreateLookupTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &LookupTableInput{
		Name:        request.GetParamLowerFirst(req.Parameters, "LookupTableName"),
		Description: request.GetParamLowerFirst(req.Parameters, "Description"),
		TableBody:   request.GetParamLowerFirst(req.Parameters, "TableBody"),
		QueryId:     request.GetParamLowerFirst(req.Parameters, "QueryId"),
		KmsKeyId:    request.GetParamLowerFirst(req.Parameters, "KmsKeyId"),
		Tags:        parseTagsFromParams(req.Parameters),
	}
	name, createdAt, err := s.createLookupTableCore(store, in, reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"lookupTableArn": lookupTableArn(reqCtx.GetRegion(), s.accountID, name),
		"createdAt":      createdAt,
	}, nil
}

// GetLookupTable retrieves the full content of a lookup table.
func (s *LogsService) GetLookupTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.getLookupTableCore(store, &GetLookupTableInput{
		Identifier: request.GetParamLowerFirst(req.Parameters, "LookupTableArn"),
		Region:     reqCtx.GetRegion(),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"lookupTableArn":  result.Arn,
		"lookupTableName": result.Name,
		"description":     result.Description,
		"tableBody":       result.TableBody,
		"sizeBytes":       result.SizeBytes,
		"lastUpdatedTime": result.LastUpdatedTime,
		"kmsKeyId":        result.KmsKeyId,
	}, nil
}

// UpdateLookupTable replaces the whole content of a lookup table.
func (s *LogsService) UpdateLookupTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	in := &UpdateLookupTableInput{
		Identifier: request.GetParamLowerFirst(req.Parameters, "LookupTableArn"),
		Region:     reqCtx.GetRegion(),
	}
	if v, present := stringParamPresent(req.Parameters, "Description"); present {
		in.Description, in.DescriptionSet = v, true
	}
	if v, present := stringParamPresent(req.Parameters, "TableBody"); present {
		in.TableBody, in.TableBodySet = v, true
	}
	if v, present := stringParamPresent(req.Parameters, "QueryId"); present {
		in.QueryId, in.QueryIdSet = v, true
	}
	if v, present := stringParamPresent(req.Parameters, "KmsKeyId"); present {
		in.KmsKeyId, in.KmsKeyIdSet = v, true
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.updateLookupTableCore(store, in)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"lookupTableArn":  result.Arn,
		"lastUpdatedTime": result.LastUpdatedTime,
	}, nil
}

// DeleteLookupTable deletes a lookup table permanently.
func (s *LogsService) DeleteLookupTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteLookupTableCore(store, request.GetParamLowerFirst(req.Parameters, "LookupTableArn")); err != nil {
		return nil, err
	}
	return map[string]interface{}{}, nil
}

// DescribeLookupTables lists lookup tables, optionally filtered by name
// prefix, sorted by name in ascending order.
func (s *LogsService) DescribeLookupTables(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.describeLookupTablesCore(store, &DescribeLookupTablesInput{
		Prefix:     request.GetParamLowerFirst(req.Parameters, "LookupTableNamePrefix"),
		MaxResults: int32(request.GetIntParam(req.Parameters, "MaxResults")),
		NextToken:  request.GetParamLowerFirst(req.Parameters, "NextToken"),
		Region:     reqCtx.GetRegion(),
	})
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(result.LookupTables))
	for _, lt := range result.LookupTables {
		items = append(items, map[string]interface{}{
			"lookupTableArn":  lt.Arn,
			"lookupTableName": lt.Name,
			"description":     lt.Description,
			"tableFields":     lt.TableFields,
			"recordsCount":    lt.RecordsCount,
			"sizeBytes":       lt.SizeBytes,
			"lastUpdatedTime": lt.LastUpdatedTime,
			"kmsKeyId":        lt.KmsKeyId,
		})
	}
	resp := map[string]interface{}{
		"lookupTables": items,
	}
	if result.NextToken != "" {
		resp["nextToken"] = result.NextToken
	}
	return resp, nil
}

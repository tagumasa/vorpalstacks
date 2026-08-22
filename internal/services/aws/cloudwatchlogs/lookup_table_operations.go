package cloudwatchlogs

import (
	"context"
	"encoding/csv"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

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

// LookupTableInput holds the parameters of CreateLookupTable and
// UpdateLookupTable.
type LookupTableInput struct {
	Name        string
	Description string
	TableBody   string
	QueryId     string
	KmsKeyId    string
	Tags        map[string]string
}

// resolveLookupTableBody returns the CSV content from the direct tableBody
// parameter or by materialising the results of the referenced query.
func (s *LogsService) resolveLookupTableBody(in *LookupTableInput) (string, error) {
	if in.TableBody != "" && in.QueryId != "" {
		return "", NewLogsError("ValidationException",
			"Specify either tableBody or queryId, but not both", 400)
	}
	if in.TableBody != "" {
		return in.TableBody, nil
	}
	if in.QueryId != "" {
		val, ok := s.queries.Load(in.QueryId)
		if !ok {
			return "", NewLogsError("ResourceNotFoundException",
				fmt.Sprintf("Query %s not found", in.QueryId), 404)
		}
		qs := val.(*queryState)
		if qs.status != "Complete" && qs.status != "Cancelled" {
			return "", NewLogsError("InvalidParameterException",
				fmt.Sprintf("Query %s has status %s; only completed or cancelled queries can populate a lookup table", in.QueryId, qs.status), 400)
		}
		return rowsToCSV(qs.results), nil
	}
	return "", NewLogsError("ValidationException",
		"Specify either tableBody or queryId", 400)
}

// rowsToCSV renders query result rows as CSV: the first row's column order
// provides the header, later rows contribute any additional columns.
func rowsToCSV(rows []queryResultRow) string {
	if len(rows) == 0 {
		return ""
	}
	var header []string
	seen := map[string]bool{}
	for i := range rows {
		for _, k := range rows[i].ordered() {
			if !seen[k] {
				seen[k] = true
				header = append(header, k)
			}
		}
	}
	var b strings.Builder
	w := csv.NewWriter(&b)
	_ = w.Write(header)
	for i := range rows {
		rec := make([]string, len(header))
		for j, k := range header {
			rec[j] = rows[i].fields[k]
		}
		_ = w.Write(rec)
	}
	w.Flush()
	return b.String()
}

// validateLookupTableSpec checks the documented name and description
// constraints.
func validateLookupTableSpec(name, description string) error {
	if name == "" || len(name) > logsstore.MaxLookupTableNameLength || !lookupTableNameRe.MatchString(name) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid lookupTableName %q: 1-%d characters, alphanumeric and underscores only",
				name, logsstore.MaxLookupTableNameLength), 400)
	}
	// LookupTableDescription @length(0-1024) counts Unicode characters.
	if utf8.RuneCountInString(description) > logsstore.MaxLookupTableDescriptionLength {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("description exceeds %d characters", logsstore.MaxLookupTableDescriptionLength), 400)
	}
	return nil
}

// applyLookupTableBody validates the CSV body, records the metadata derived
// from it, and stores it either as plaintext or envelope-encrypted under the
// table's KMS key.
func (s *LogsService) applyLookupTableBody(lt *logsstore.LookupTable, body, region string) error {
	if len(body) > logsstore.MaxLookupTableBodyBytes {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("tableBody exceeds %d bytes", logsstore.MaxLookupTableBodyBytes), 400)
	}
	fields, records, err := parseLookupCSV(body)
	if err != nil {
		return err
	}
	lt.TableBody = body
	lt.TableFields = fields
	lt.RecordsCount = int64(len(records))
	lt.SizeBytes = int64(len(body))
	lt.EncryptedBody = nil
	lt.EncryptedDataKey = nil
	lt.ContentNonce = nil
	if lt.KmsKeyId != "" {
		tableArn := lookupTableArn(region, s.accountID, lt.Name)
		encrypted, dataKey, nonce, err := s.encryptLookupTableBody([]byte(body), lt.KmsKeyId, tableArn)
		if err != nil {
			return err
		}
		lt.TableBody = ""
		lt.EncryptedBody = encrypted
		lt.EncryptedDataKey = dataKey
		lt.ContentNonce = nonce
	}
	return nil
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

func (s *LogsService) createLookupTableCore(store *logsstore.Store, in *LookupTableInput, region string) (string, int64, error) {
	if err := validateLookupTableSpec(in.Name, in.Description); err != nil {
		return "", 0, err
	}
	if len(in.Tags) > logsstore.MaxLookupTableTags {
		return "", 0, NewLogsError("InvalidParameterException",
			fmt.Sprintf("a maximum of %d tags can be attached to a lookup table", logsstore.MaxLookupTableTags), 400)
	}
	if err := s.validateLookupTableKmsKey(in.KmsKeyId); err != nil {
		return "", 0, err
	}
	if _, err := store.GetLookupTable(in.Name); err == nil {
		return "", 0, NewLogsError("ResourceAlreadyExistsException",
			fmt.Sprintf("Lookup table %s already exists", in.Name), 400)
	}
	body, err := s.resolveLookupTableBody(in)
	if err != nil {
		return "", 0, err
	}
	count, err := store.CountLookupTables()
	if err != nil {
		return "", 0, err
	}
	if count >= logsstore.MaxLookupTables {
		return "", 0, NewLogsError("LimitExceededException",
			fmt.Sprintf("A maximum of %d lookup tables can exist per account per Region", logsstore.MaxLookupTables), 400)
	}
	lt := &logsstore.LookupTable{
		Name:        in.Name,
		Description: in.Description,
		KmsKeyId:    in.KmsKeyId,
		Tags:        in.Tags,
	}
	if err := s.applyLookupTableBody(lt, body, region); err != nil {
		return "", 0, err
	}
	lt.CreationTime = time.Now().UTC().UnixMilli()
	if err := store.PutLookupTable(lt); err != nil {
		return "", 0, err
	}
	return lt.Name, lt.CreationTime, nil
}

// GetLookupTable retrieves the full content of a lookup table.
func (s *LogsService) GetLookupTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	identifier := request.GetParamLowerFirst(req.Parameters, "LookupTableArn")
	if identifier == "" {
		return nil, ErrMissingParameter
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	lt, err := store.GetLookupTable(lookupTableNameFromArn(identifier))
	if err != nil {
		return nil, mapLookupTableStoreError(err)
	}
	body, err := s.lookupTablePlainBody(lt, reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}
	arn := lookupTableArn(reqCtx.GetRegion(), s.accountID, lt.Name)
	return map[string]interface{}{
		"lookupTableArn":  arn,
		"lookupTableName": lt.Name,
		"description":     lt.Description,
		"tableBody":       body,
		"sizeBytes":       lt.SizeBytes,
		"lastUpdatedTime": lt.LastUpdatedTime,
		"kmsKeyId":        lt.KmsKeyId,
	}, nil
}

// mapLookupTableStoreError maps store sentinels to API errors.
func mapLookupTableStoreError(err error) error {
	if err == logsstore.ErrResourceNotFound {
		return NewLogsError("ResourceNotFoundException", "Lookup table not found", 400)
	}
	return err
}

// UpdateLookupTable replaces the whole content of a lookup table.
func (s *LogsService) UpdateLookupTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	identifier := request.GetParamLowerFirst(req.Parameters, "LookupTableArn")
	if identifier == "" {
		return nil, ErrMissingParameter
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	lt, err := store.GetLookupTable(lookupTableNameFromArn(identifier))
	if err != nil {
		return nil, mapLookupTableStoreError(err)
	}

	in := &LookupTableInput{
		Name:     lt.Name,
		KmsKeyId: lt.KmsKeyId,
	}
	if v, present := stringParamPresent(req.Parameters, "Description"); present {
		in.Description = v
	} else {
		in.Description = lt.Description
	}
	if v, present := stringParamPresent(req.Parameters, "TableBody"); present {
		in.TableBody = v
	}
	if v, present := stringParamPresent(req.Parameters, "QueryId"); present {
		in.QueryId = v
	}
	if v, present := stringParamPresent(req.Parameters, "KmsKeyId"); present {
		in.KmsKeyId = v
	}
	if err := s.validateLookupTableKmsKey(in.KmsKeyId); err != nil {
		return nil, err
	}

	body, err := s.resolveLookupTableBody(in)
	if err != nil {
		return nil, err
	}
	lt.Description = in.Description
	lt.KmsKeyId = in.KmsKeyId
	if err := s.applyLookupTableBody(lt, body, reqCtx.GetRegion()); err != nil {
		return nil, err
	}
	if err := store.PutLookupTable(lt); err != nil {
		return nil, err
	}
	arn := lookupTableArn(reqCtx.GetRegion(), s.accountID, lt.Name)
	return map[string]interface{}{
		"lookupTableArn":  arn,
		"lastUpdatedTime": lt.LastUpdatedTime,
	}, nil
}

// DeleteLookupTable deletes a lookup table permanently.
func (s *LogsService) DeleteLookupTable(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	identifier := request.GetParamLowerFirst(req.Parameters, "LookupTableArn")
	if identifier == "" {
		return nil, ErrMissingParameter
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeleteLookupTable(lookupTableNameFromArn(identifier)); err != nil {
		return nil, mapLookupTableStoreError(err)
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
	prefix := request.GetParamLowerFirst(req.Parameters, "LookupTableNamePrefix")
	maxResults := int32(request.GetIntParam(req.Parameters, "MaxResults"))
	if maxResults < 0 || maxResults > 100 {
		return nil, NewLogsError("InvalidParameterException",
			"maxResults must be between 1 and 100", 400)
	}
	if maxResults == 0 {
		maxResults = 100
	}
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")

	tables, err := store.ListLookupTables(prefix)
	if err != nil {
		return nil, err
	}

	offset := 0
	if nextToken != "" {
		n, err := parseInt(nextToken)
		if err != nil {
			return nil, NewLogsError("InvalidParameterException", "Invalid nextToken", 400)
		}
		offset = n
	}
	if offset > len(tables) {
		offset = len(tables)
	}
	end := offset + int(maxResults)
	if end > len(tables) {
		end = len(tables)
	}

	items := make([]map[string]interface{}, 0, end-offset)
	for _, lt := range tables[offset:end] {
		items = append(items, map[string]interface{}{
			"lookupTableArn":  lookupTableArn(reqCtx.GetRegion(), s.accountID, lt.Name),
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
	if end < len(tables) {
		resp["nextToken"] = fmt.Sprintf("%d", end)
	}
	return resp, nil
}

package cloudwatchlogs

import (
	"encoding/csv"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

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

// UpdateLookupTableInput carries the presence-parsed members of
// UpdateLookupTable: a member is applied only when its Set flag is true,
// otherwise the stored value is preserved.
type UpdateLookupTableInput struct {
	Identifier     string
	Description    string
	DescriptionSet bool
	TableBody      string
	TableBodySet   bool
	QueryId        string
	QueryIdSet     bool
	KmsKeyId       string
	KmsKeyIdSet    bool
	Region         string
}

// UpdateLookupTableResult is the response payload of UpdateLookupTable.
type UpdateLookupTableResult struct {
	Arn             string
	LastUpdatedTime int64
}

// GetLookupTableInput identifies a lookup table by name or ARN.
type GetLookupTableInput struct {
	Identifier string
	Region     string
}

// GetLookupTableResult is the full content projection of a lookup table.
type GetLookupTableResult struct {
	Arn             string
	Name            string
	Description     string
	TableBody       string
	SizeBytes       int64
	LastUpdatedTime int64
	KmsKeyId        string
}

// DescribeLookupTablesInput holds the parsed list parameters. MaxResults is
// the raw client value: 0 means the member was omitted and the documented
// default of 50 applies, values above the documented maximum of 100 are
// rejected, and negative values are rejected as invalid input (the model
// documents no minimum).
type DescribeLookupTablesInput struct {
	Prefix     string
	MaxResults int32
	NextToken  string
	Region     string
}

// LookupTableSummary is one entry of the DescribeLookupTables response.
type LookupTableSummary struct {
	Arn             string
	Name            string
	Description     string
	TableFields     []string
	RecordsCount    int64
	SizeBytes       int64
	LastUpdatedTime int64
	KmsKeyId        string
}

// DescribeLookupTablesResult is the paginated DescribeLookupTables payload.
type DescribeLookupTablesResult struct {
	LookupTables []LookupTableSummary
	NextToken    string
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
				fmt.Sprintf("Query %s not found", in.QueryId), 400)
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

// mapLookupTableStoreError maps store sentinels to API errors.
func mapLookupTableStoreError(err error) error {
	if err == logsstore.ErrResourceNotFound {
		return NewLogsError("ResourceNotFoundException", "Lookup table not found", 400)
	}
	return err
}

// createLookupTableCore validates input and creates a lookup table.
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

// getLookupTableCore resolves the identifier (name or ARN) and returns the
// decrypted full content of the lookup table.
func (s *LogsService) getLookupTableCore(store *logsstore.Store, in *GetLookupTableInput) (*GetLookupTableResult, error) {
	if in.Identifier == "" {
		return nil, ErrMissingParameter
	}
	lt, err := store.GetLookupTable(lookupTableNameFromArn(in.Identifier))
	if err != nil {
		return nil, mapLookupTableStoreError(err)
	}
	body, err := s.lookupTablePlainBody(lt, in.Region)
	if err != nil {
		return nil, err
	}
	return &GetLookupTableResult{
		Arn:             lookupTableArn(in.Region, s.accountID, lt.Name),
		Name:            lt.Name,
		Description:     lt.Description,
		TableBody:       body,
		SizeBytes:       lt.SizeBytes,
		LastUpdatedTime: lt.LastUpdatedTime,
		KmsKeyId:        lt.KmsKeyId,
	}, nil
}

// updateLookupTableCore merges the provided members onto the stored table
// (omitted members keep their stored values), revalidates the body source and
// the KMS key, and persists the replacement content.
func (s *LogsService) updateLookupTableCore(store *logsstore.Store, in *UpdateLookupTableInput) (*UpdateLookupTableResult, error) {
	if in.Identifier == "" {
		return nil, ErrMissingParameter
	}
	lt, err := store.GetLookupTable(lookupTableNameFromArn(in.Identifier))
	if err != nil {
		return nil, mapLookupTableStoreError(err)
	}

	update := &LookupTableInput{
		Name:     lt.Name,
		KmsKeyId: lt.KmsKeyId,
	}
	if in.DescriptionSet {
		update.Description = in.Description
	} else {
		update.Description = lt.Description
	}
	if in.TableBodySet {
		update.TableBody = in.TableBody
	}
	if in.QueryIdSet {
		update.QueryId = in.QueryId
	}
	if in.KmsKeyIdSet {
		update.KmsKeyId = in.KmsKeyId
	}
	if err := s.validateLookupTableKmsKey(update.KmsKeyId); err != nil {
		return nil, err
	}

	body, err := s.resolveLookupTableBody(update)
	if err != nil {
		return nil, err
	}
	lt.Description = update.Description
	lt.KmsKeyId = update.KmsKeyId
	if err := s.applyLookupTableBody(lt, body, in.Region); err != nil {
		return nil, err
	}
	if err := store.PutLookupTable(lt); err != nil {
		return nil, err
	}
	return &UpdateLookupTableResult{
		Arn:             lookupTableArn(in.Region, s.accountID, lt.Name),
		LastUpdatedTime: lt.LastUpdatedTime,
	}, nil
}

// deleteLookupTableCore validates input and deletes a lookup table
// permanently.
func (s *LogsService) deleteLookupTableCore(store *logsstore.Store, identifier string) error {
	if identifier == "" {
		return ErrMissingParameter
	}
	if err := store.DeleteLookupTable(lookupTableNameFromArn(identifier)); err != nil {
		return mapLookupTableStoreError(err)
	}
	return nil
}

// describeLookupTablesCore lists lookup tables filtered by name prefix,
// paginated by numeric offset, sorted by name in ascending order.
func (s *LogsService) describeLookupTablesCore(store *logsstore.Store, in *DescribeLookupTablesInput) (*DescribeLookupTablesResult, error) {
	if in.MaxResults < 0 || in.MaxResults > logsstore.MaxDescribeLookupTablesResults {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("maxResults must be between 1 and %d", logsstore.MaxDescribeLookupTablesResults), 400)
	}
	maxResults := in.MaxResults
	if maxResults == 0 {
		maxResults = logsstore.DefaultDescribeLookupTablesResults
	}

	tables, err := store.ListLookupTables(in.Prefix)
	if err != nil {
		return nil, err
	}

	offset := 0
	if in.NextToken != "" {
		n, err := parseInt(in.NextToken)
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

	result := &DescribeLookupTablesResult{
		LookupTables: make([]LookupTableSummary, 0, end-offset),
	}
	for _, lt := range tables[offset:end] {
		result.LookupTables = append(result.LookupTables, LookupTableSummary{
			Arn:             lookupTableArn(in.Region, s.accountID, lt.Name),
			Name:            lt.Name,
			Description:     lt.Description,
			TableFields:     lt.TableFields,
			RecordsCount:    lt.RecordsCount,
			SizeBytes:       lt.SizeBytes,
			LastUpdatedTime: lt.LastUpdatedTime,
			KmsKeyId:        lt.KmsKeyId,
		})
	}
	if end < len(tables) {
		result.NextToken = fmt.Sprintf("%d", end)
	}
	return result, nil
}

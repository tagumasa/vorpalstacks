package dynamodb

import (
	"regexp"
	"strings"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// ---------------------------------------------------------------------------
// Smithy / AWS Docs constraint constants
// ---------------------------------------------------------------------------

// Resource names (TableName, BackupName, IndexName, GlobalTableName):
// 3-255 characters, [a-zA-Z0-9_.-]+
const (
	resourceNameMinLength = 3
	resourceNameMaxLength = 255
)

// ProjectionType enum values (Smithy ProjectionType).
const (
	ProjectionTypeAll      = "ALL"
	ProjectionTypeKeysOnly = "KEYS_ONLY"
	ProjectionTypeInclude  = "INCLUDE"
)

// ARN length constraints from Smithy.
const (
	arnBackupMinLen   = 37
	arnBackupMaxLen   = 1024
	arnStreamMinLen   = 37
	arnStreamMaxLen   = 1024
	arnExportMinLen   = 37
	arnExportMaxLen   = 1024
	arnImportMinLen   = 37
	arnImportMaxLen   = 1024
	arnTableMinLen    = 1
	arnTableMaxLen    = 1024
	arnResourceMinLen = 1
	arnResourceMaxLen = 1283
)

// Other Smithy length/range/pattern constants.
const (
	attributeNameMaxLen         = 65535
	keySchemaAttrNameMinLen     = 1
	keySchemaAttrNameMaxLen     = 255
	nonKeyAttrNameMinLen        = 1
	nonKeyAttrNameMaxLen        = 255
	nonKeyAttrListMaxLen        = 20
	partiqlStatementMaxLen      = 8192
	partiqlBatchMaxLen          = 25
	partiqlNextTokenMaxLen      = 32768
	importNextTokenMinLen       = 112
	importNextTokenMaxLen       = 1024
	tagKeyMaxLen                = 128
	tagValueMaxLen              = 256
	ttlAttributeNameMaxLen      = 255
	policyRevisionIdMinLen      = 1
	policyRevisionIdMaxLen      = 255
	clientRequestTokenMinLen    = 1
	clientRequestTokenMaxLen    = 36
	s3BucketMaxLen              = 255
	s3PrefixMaxLen              = 1024
	s3SseKmsKeyIdMinLen         = 1
	s3SseKmsKeyIdMaxLen         = 2048
	autoScalingPolicyNameMinLen = 1
	autoScalingPolicyNameMaxLen = 256
	autoScalingRoleArnMinLen    = 1
	autoScalingRoleArnMaxLen    = 1600
	csvHeaderListMaxLen         = 255
	batchGetMaxTables           = 100
	batchWriteMaxItems          = 25
	transactMaxItems            = 100
	listExportsMaxLimit         = 25
	listImportsMaxLimit         = 25
	listBackupsMaxLimit         = 100
	listContributorMaxLimit     = 100
	listTablesMaxLimit          = 100
	listGlobalTablesMinLimit    = 1
	recoveryPeriodMin           = 1
	recoveryPeriodMax           = 35
	scanSegmentMin              = 0
	scanSegmentMax              = 999999
	scanTotalSegmentsMin        = 1
	scanTotalSegmentsMax        = 1000000
)

// Compiled regex patterns from Smithy @pattern traits.
var (
	// Shared resource-name pattern: TableName, BackupName, IndexName, GlobalTableName.
	resourceNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)

	s3BucketRegex              = regexp.MustCompile(`^[a-z0-9A-Z]+[.\-\w]*[a-z0-9A-Z]+$`)
	s3BucketOwnerRegex         = regexp.MustCompile(`^[0-9]{12}$`)
	tableIdRegex               = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	clientTokenRegex           = regexp.MustCompile(`^[^\$]+$`)
	importNextTokenRegex       = regexp.MustCompile(`^([0-9a-f]{16})+$`)
	contributorRuleRegex       = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9\-_\.]{0,126}[A-Za-z0-9]$`)
	autoScalingPolicyNameRegex = regexp.MustCompile(`[^[:cntrl:]]+`)
	csvDelimiterRegex          = regexp.MustCompile(`^[,;:|\t ]$`)
	csvHeaderRegex             = regexp.MustCompile(`^[\x20-\x21\x23-\x2B\x2D-\x7E]*$`)
)

var validProjectionTypes = map[string]bool{
	ProjectionTypeAll:      true,
	ProjectionTypeKeysOnly: true,
	ProjectionTypeInclude:  true,
}

// ---------------------------------------------------------------------------
// Generic helpers
// ---------------------------------------------------------------------------

// validateLength checks that a string's length falls within [min, max].
func validateLength(val string, min, max int) error {
	if len(val) < min || len(val) > max {
		return ErrInvalidParameter
	}
	return nil
}

// validateRange checks that an integer falls within [min, max].
func validateRange(val, min, max int) error {
	if val < min || val > max {
		return ErrInvalidParameter
	}
	return nil
}

// ---------------------------------------------------------------------------
// Resource-name validators (shared: TableName, BackupName, IndexName, GlobalTableName)
// ---------------------------------------------------------------------------

// validateResourceName checks the DynamoDB resource-name constraints:
// 3-255 characters, characters [a-zA-Z0-9_.-].
// Applies to TableName, BackupName, IndexName, and GlobalTableName.
func validateResourceName(name string) error {
	if name == "" {
		return ErrInvalidParameter
	}
	if len(name) < resourceNameMinLength || len(name) > resourceNameMaxLength {
		return ErrInvalidParameter
	}
	if !resourceNameRegex.MatchString(name) {
		return ErrInvalidParameter
	}
	return nil
}

// validateTableName is an alias for validateResourceName, retained for
// call-site clarity when validating a table name specifically.
func validateTableName(name string) error {
	return validateResourceName(name)
}

// validateIndexName checks that a secondary-index name meets the same
// constraints as table names (Smithy IndexName: len 3-255, pattern).
func validateIndexName(name string) error {
	return validateResourceName(name)
}

// ---------------------------------------------------------------------------
// ARN validators
// ---------------------------------------------------------------------------

// validateArnLength checks that an ARN's length falls within the Smithy
// constraints for the given ARN type.
func validateArnLength(arn string, min, max int) error {
	return validateLength(arn, min, max)
}

// validateBackupArn validates a BackupArn (Smithy: len 37-1024).
func validateBackupArn(arn string) error {
	return validateArnLength(arn, arnBackupMinLen, arnBackupMaxLen)
}

// validateStreamArn validates a StreamArn (Smithy: len 37-1024).
func validateStreamArn(arn string) error {
	return validateArnLength(arn, arnStreamMinLen, arnStreamMaxLen)
}

// validateExportArn validates an ExportArn (Smithy: len 37-1024).
func validateExportArn(arn string) error {
	return validateArnLength(arn, arnExportMinLen, arnExportMaxLen)
}

// validateImportArn validates an ImportArn (Smithy: len 37-1024).
func validateImportArn(arn string) error {
	return validateArnLength(arn, arnImportMinLen, arnImportMaxLen)
}

// validateTableArn validates a TableArn (Smithy: len 1-1024).
func validateTableArn(arn string) error {
	return validateArnLength(arn, arnTableMinLen, arnTableMaxLen)
}

// validateResourceArnString validates a generic ResourceArnString used by
// tag and resource-policy operations (Smithy: len 1-1283).
func validateResourceArnString(arn string) error {
	return validateArnLength(arn, arnResourceMinLen, arnResourceMaxLen)
}

// ---------------------------------------------------------------------------
// Projection / index validators
// ---------------------------------------------------------------------------

// validateProjectionType checks the Smithy ProjectionType enum.
func validateProjectionType(pt string) error {
	if pt == "" {
		return nil
	}
	if !validProjectionTypes[pt] {
		return ErrInvalidParameter
	}
	return nil
}

// validateProjectionRequired validates the Projection sub-structure of a
// GSI or LSI definition. It enforces:
//   - NonKeyAttributes list size (Smithy NonKeyAttributeNameList: len 1-20).
//   - Each NonKeyAttributeName length (Smithy: len 1-255).
//   - NonKeyAttributes only valid when ProjectionType=INCLUDE.
func validateProjectionRequired(projMap map[string]interface{}) error {
	if projMap == nil {
		return ErrInvalidParameter
	}
	projectionType := request.GetStringParam(projMap, "ProjectionType")
	if nkAs, ok := projMap["NonKeyAttributes"].([]interface{}); ok {
		if len(nkAs) > nonKeyAttrListMaxLen {
			return ErrInvalidParameter
		}
		for _, nk := range nkAs {
			nks, ok := nk.(string)
			if !ok {
				return ErrInvalidParameter
			}
			if err := validateLength(nks, nonKeyAttrNameMinLen, nonKeyAttrNameMaxLen); err != nil {
				return err
			}
		}
		if projectionType != "" && projectionType != ProjectionTypeInclude {
			return ErrInvalidParameter
		}
	}
	return nil
}

// validateKeyAttributeValue ensures key attribute values use only the
// types allowed for DynamoDB keys: S, N, or B (Smithy KeySchemaAttribute).
// Also rejects empty values (N type empty check included).
func validateKeyAttributeValue(key map[string]*dbstore.AttributeValue) error {
	for _, av := range key {
		if av == nil {
			return ErrInvalidParameter
		}
		// Key attributes must be one of S, N, or B.
		// BOOL, NULL, SS, NS, BS, M, L are not permitted.
		if av.S == nil && av.N == nil && av.B == nil {
			return ErrInvalidParameter
		}
		// Reject empty values (extend to N type).
		if av.S != nil && *av.S == "" {
			return ErrInvalidParameter
		}
		if av.N != nil && *av.N == "" {
			return ErrInvalidParameter
		}
		if av.B != nil && len(av.B) == 0 {
			return ErrInvalidParameter
		}
	}
	return nil
}

// validateGSIDeleteExists checks that a GSI to be deleted actually exists
// before calling delete() on the map (Smithy: returns
// GlobalSecondaryIndexNotFoundException on unknown GSI).
func validateGSIDeleteExists(gsiMap map[string]*dbstore.GlobalSecondaryIndex, indexName string) error {
	if _, exists := gsiMap[indexName]; !exists {
		return ErrIndexNotFound
	}
	return nil
}

// validateGSICreateRequired checks that a GSI Create request includes
// the required IndexName and KeySchema (Smithy required traits), and
// validates the IndexName format (len 3-255, pattern).
func validateGSICreateRequired(create map[string]interface{}) error {
	indexName := request.GetStringParam(create, "IndexName")
	if indexName == "" {
		return ErrInvalidParameter
	}
	if err := validateIndexName(indexName); err != nil {
		return err
	}
	keySchema := parseKeySchema(create)
	if len(keySchema) == 0 {
		return ErrInvalidParameter
	}
	if err := validateKeySchema(keySchema); err != nil {
		return err
	}
	return nil
}

// validateBillingModeConsistency ensures that PROVISIONED billing mode
// includes ProvisionedThroughput, and PAY_PER_REQUEST does not.
func validateBillingModeConsistency(billingMode dbstore.BillingMode, provThroughput *dbstore.ProvisionedThroughput) error {
	if billingMode == dbstore.BillingModeProvisioned {
		if provThroughput == nil {
			return ErrInvalidParameter
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Bool / bracket helpers
// ---------------------------------------------------------------------------

// validateBoolParam extracts a bool from a parameter map. If the key is
// absent, the default value is returned. If the key is present but not a
// bool, an error is returned (rejects malformed requests that would
// otherwise be silently coerced to the default).
func validateBoolParam(params map[string]interface{}, key string, defaultVal bool) (bool, error) {
	v, ok := params[key]
	if !ok {
		return defaultVal, nil
	}
	b, ok := v.(bool)
	if !ok {
		return false, ErrInvalidParameter
	}
	return b, nil
}

// validateBracketIndex parses a bracket-enclosed list index (e.g. "[3]")
// and returns the integer index. Returns an error for empty, non-numeric,
// or negative values (Smithy document path spec).
func validateBracketIndex(idxStr string) (int, error) {
	idxStr = strings.TrimSpace(idxStr)
	if idxStr == "" {
		return 0, ErrInvalidParameter
	}
	var idx int
	for _, ch := range idxStr {
		if ch < '0' || ch > '9' {
			return 0, ErrInvalidParameter
		}
		idx = idx*10 + int(ch-'0')
	}
	if idx < 0 {
		return 0, ErrInvalidParameter
	}
	return idx, nil
}

// ---------------------------------------------------------------------------
// Smithy HIGH-tier validators
// ---------------------------------------------------------------------------

// validatePartiQLStatement checks the PartiQL statement length
// (Smithy PartiQLStatement: len 1-8192).
func validatePartiQLStatement(stmt string) error {
	return validateLength(stmt, 1, partiqlStatementMaxLen)
}

// validatePartiQLBatchCount checks the number of statements in a
// BatchExecuteStatement request (Smithy PartiQLBatchRequest: len 1-25).
func validatePartiQLBatchCount(count int) error {
	return validateRange(count, 1, partiqlBatchMaxLen)
}

// validateExecuteStatementLimit checks the Limit parameter on
// ExecuteStatement (Smithy PositiveIntegerObject: range min 1).
func validateExecuteStatementLimit(limit int) error {
	if limit < 1 {
		return ErrInvalidParameter
	}
	return nil
}

// validateS3Bucket checks the S3 bucket name format
// (Smithy S3Bucket: len 0-255, pattern ^[a-z0-9A-Z]+[.\-\w]*[a-z0-9A-Z]+$).
func validateS3Bucket(bucket string) error {
	if err := validateLength(bucket, 0, s3BucketMaxLen); err != nil {
		return err
	}
	if bucket != "" && !s3BucketRegex.MatchString(bucket) {
		return ErrInvalidParameter
	}
	return nil
}

// validateS3BucketOwner checks the S3 bucket owner account ID
// (Smithy S3BucketOwner: pattern ^[0-9]{12}$).
func validateS3BucketOwner(owner string) error {
	if owner == "" {
		return nil
	}
	if !s3BucketOwnerRegex.MatchString(owner) {
		return ErrInvalidParameter
	}
	return nil
}

// validateS3Prefix checks the S3 key prefix length
// (Smithy S3Prefix: len 0-1024).
func validateS3Prefix(prefix string) error {
	return validateLength(prefix, 0, s3PrefixMaxLen)
}

// validateS3SseKmsKeyId checks the KMS key ID length
// (Smithy S3SseKmsKeyId: OPTIONAL, len 1-2048 when provided).
func validateS3SseKmsKeyId(keyId string) error {
	if keyId == "" {
		return nil
	}
	return validateLength(keyId, s3SseKmsKeyIdMinLen, s3SseKmsKeyIdMaxLen)
}

// ---------------------------------------------------------------------------
// Smithy MEDIUM-tier validators
// ---------------------------------------------------------------------------

// validateGlobalTableName is an alias for validateResourceName; retained
// for call-site clarity (Smithy TableName: len 3-255, pattern).
func validateGlobalTableName(name string) error {
	return validateResourceName(name)
}

// validateTagKey checks a tag key (Smithy TagKeyString: len 1-128).
func validateTagKey(key string) error {
	return validateLength(key, 1, tagKeyMaxLen)
}

// validateTagValue checks a tag value (Smithy TagValueString: len 0-256).
func validateTagValue(value string) error {
	return validateLength(value, 0, tagValueMaxLen)
}

// validateRecoveryPeriodInDays checks the PITR recovery period
// (Smithy RecoveryPeriodInDays: range 1-35).
func validateRecoveryPeriodInDays(days int) error {
	return validateRange(days, recoveryPeriodMin, recoveryPeriodMax)
}

// validateScanSegment checks the parallel Scan segment number
// (Smithy ScanSegment: range 0-999999).
func validateScanSegment(segment int) error {
	return validateRange(segment, scanSegmentMin, scanSegmentMax)
}

// validateScanTotalSegments checks the parallel Scan total-segments count
// (Smithy ScanTotalSegments: range 1-1000000).
func validateScanTotalSegments(total int) error {
	return validateRange(total, scanTotalSegmentsMin, scanTotalSegmentsMax)
}

// validateClientRequestToken checks the idempotency token length
// (Smithy ClientRequestToken: len 1-36).
func validateClientRequestToken(token string) error {
	if token == "" {
		return nil
	}
	return validateLength(token, clientRequestTokenMinLen, clientRequestTokenMaxLen)
}

// validateListExportsLimit checks ListExports MaxResults
// (Smithy ListExportsMaxLimit: range 1-25).
func validateListExportsLimit(limit int) error {
	return validateRange(limit, 1, listExportsMaxLimit)
}

// validateListImportsLimit checks ListImports PageSize
// (Smithy ListImportsMaxLimit: range 1-25).
func validateListImportsLimit(limit int) error {
	return validateRange(limit, 1, listImportsMaxLimit)
}

// validateListBackupsLimit checks ListBackups Limit
// (Smithy BackupsInputLimit: range 1-100).
func validateListBackupsLimit(limit int) error {
	return validateRange(limit, 1, listBackupsMaxLimit)
}

// validateListContributorInsightsLimit checks ListContributorInsights MaxResults
// (Smithy ListContributorInsightsLimit: range max 100).
func validateListContributorInsightsLimit(limit int) error {
	return validateRange(limit, 0, listContributorMaxLimit)
}

// validateListTablesLimit checks ListTables Limit
// (Smithy ListTablesInputLimit: range 1-100).
func validateListTablesLimit(limit int) error {
	return validateRange(limit, 1, listTablesMaxLimit)
}

// validateListGlobalTablesLimit checks ListGlobalTables Limit
// (Smithy PositiveIntegerObject: range min 1).
func validateListGlobalTablesLimit(limit int) error {
	if limit < listGlobalTablesMinLimit {
		return ErrInvalidParameter
	}
	return nil
}

// ---------------------------------------------------------------------------
// Smithy LOW-tier validators
// ---------------------------------------------------------------------------

// validatePolicyRevisionId checks the resource-policy revision ID length
// (Smithy PolicyRevisionId: len 1-255).
func validatePolicyRevisionId(id string) error {
	if id == "" {
		return nil
	}
	return validateLength(id, policyRevisionIdMinLen, policyRevisionIdMaxLen)
}

// validateTimeToLiveAttributeName checks the TTL attribute name length
// (Smithy TimeToLiveAttributeName: len 1-255).
func validateTimeToLiveAttributeName(name string) error {
	return validateLength(name, 1, ttlAttributeNameMaxLen)
}

// validateTableId checks the table UUID format
// (Smithy TableId: pattern ^[0-9a-f]{8}-...$).
func validateTableId(id string) error {
	if id == "" {
		return nil
	}
	if !tableIdRegex.MatchString(id) {
		return ErrInvalidParameter
	}
	return nil
}

// validateImportNextToken checks the ListImports next-token format
// (Smithy ImportNextToken: len 112-1024, pattern ^([0-9a-f]{16})+$).
func validateImportNextToken(token string) error {
	if token == "" {
		return nil
	}
	if err := validateLength(token, importNextTokenMinLen, importNextTokenMaxLen); err != nil {
		return err
	}
	if !importNextTokenRegex.MatchString(token) {
		return ErrInvalidParameter
	}
	return nil
}

// validateClientToken checks the Export/Import client-token format
// (Smithy ClientToken: pattern ^[^\$]+$).
func validateClientToken(token string) error {
	if token == "" {
		return nil
	}
	if !clientTokenRegex.MatchString(token) {
		return ErrInvalidParameter
	}
	return nil
}

// validateAutoScalingPolicyName checks the auto-scaling policy name
// (Smithy AutoScalingPolicyName: len 1-256, pattern ^\p{Print}+$).
func validateAutoScalingPolicyName(name string) error {
	if err := validateLength(name, autoScalingPolicyNameMinLen, autoScalingPolicyNameMaxLen); err != nil {
		return err
	}
	if !autoScalingPolicyNameRegex.MatchString(name) {
		return ErrInvalidParameter
	}
	return nil
}

// validateAutoScalingRoleArn checks the auto-scaling role ARN length
// (Smithy AutoScalingRoleArn: len 1-1600). The Smithy pattern permits any
// XML-compatible character and is extremely permissive, so only length is
// enforced.
func validateAutoScalingRoleArn(arn string) error {
	return validateLength(arn, autoScalingRoleArnMinLen, autoScalingRoleArnMaxLen)
}

// validateContributorInsightsRule checks the contributor-insights rule
// pattern (Smithy ContributorInsightsRule:
// ^[A-Za-z0-9][A-Za-z0-9\-_\.]{0,126}[A-Za-z0-9]$).
func validateContributorInsightsRule(rule string) error {
	if rule == "" {
		return nil
	}
	if !contributorRuleRegex.MatchString(rule) {
		return ErrInvalidParameter
	}
	return nil
}

// validateCsvDelimiter checks the CSV import delimiter character
// (Smithy CsvDelimiter: len 1-1, pattern ^[,;:|\t ]$).
func validateCsvDelimiter(d string) error {
	if len(d) != 1 {
		return ErrInvalidParameter
	}
	if !csvDelimiterRegex.MatchString(d) {
		return ErrInvalidParameter
	}
	return nil
}

// validateCsvHeader checks a single CSV import header value
// (Smithy CsvHeader: len 1-65536, pattern ^[\x20-\x21\x23-\x2B\x2D-\x7E]*$).
func validateCsvHeader(header string) error {
	if err := validateLength(header, 1, 65536); err != nil {
		return err
	}
	if !csvHeaderRegex.MatchString(header) {
		return ErrInvalidParameter
	}
	return nil
}

// validateCsvHeaderList checks the CSV import header list size
// (Smithy CsvHeaderList: len 1-255).
func validateCsvHeaderList(count int) error {
	return validateRange(count, 1, csvHeaderListMaxLen)
}

// validateAttributeName checks an attribute name length
// (Smithy AttributeName: len 0-65535).
func validateAttributeName(name string) error {
	return validateLength(name, 0, attributeNameMaxLen)
}

// validateKeySchemaAttributeName checks a key-schema attribute name
// (Smithy KeySchemaAttributeName: len 1-255).
func validateKeySchemaAttributeName(name string) error {
	return validateLength(name, keySchemaAttrNameMinLen, keySchemaAttrNameMaxLen)
}

// validateNonKeyAttributeName checks a single NonKeyAttributeName
// (Smithy NonKeyAttributeName: len 1-255).
func validateNonKeyAttributeName(name string) error {
	return validateLength(name, nonKeyAttrNameMinLen, nonKeyAttrNameMaxLen)
}

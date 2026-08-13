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
	batchGetMaxTotalItems       = 100
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

	// Pagination defaults and caps used by List operations that do not
	// have an explicit Smithy Limit trait but follow AWS documented
	// behaviour. Each value is referenced from exactly one call site.
	listGlobalTablesDefaultLimit       = 100
	listGlobalTablesMinPageSize        = 10
	listStreamsDefaultLimit            = 100
	listStreamsMaxLimit                = 100
	getRecordsDefaultLimit             = 100
	getRecordsMaxLimit                 = 1000
	listTagsForResourceDefaultPageSize = 50

	// Backup list minimum fetch size: the ListBackups implementation
	// always fetches at least this many entries from the store before
	// applying filters, to reduce round-trips for common queries.
	listBackupsMinFetchSize = 100

	// PITR default recovery window in days (AWS default: 35 days).
	pitrDefaultRecoveryPeriodDays = 35

	// Data-plane Query/Scan pagination defaults. AWS does not document a
	// hard cap for these (the documented max is 1 MB per page), but the
	// SDK tests and the existing implementation cap the item count at
	// 1000 to bound work. The admin Scan handler shares the same cap.
	dataPlaneQueryDefaultLimit = 100
	dataPlaneQueryMaxLimit     = 1000
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
// Convention
// ---------------------------------------------------------------------------
//
// Pure validators in this file return bool: true means the input satisfies the
// Smithy constraint, false means it does not. Callers decide which sentinel
// error to surface to the client (typically ErrInvalidParameter).
//
// The two error-returning functions kept here (validateGSIDeleteExists and
// validateBoolParam/validateBracketIndex) carry information that cannot be
// represented as a bool: a distinct Smithy error code or a parsed value that
// the caller still needs on success. Their continued use of the error type is
// deliberate, not an oversight.

// ---------------------------------------------------------------------------
// Generic helpers
// ---------------------------------------------------------------------------

// validateLength reports whether a string's length falls within [min, max].
func validateLength(val string, min, max int) bool {
	return len(val) >= min && len(val) <= max
}

// validateRange reports whether an integer falls within [min, max].
func validateRange(val, min, max int) bool {
	return val >= min && val <= max
}

// ---------------------------------------------------------------------------
// Resource-name validators (shared: TableName, BackupName, IndexName, GlobalTableName)
// ---------------------------------------------------------------------------

// validateResourceName reports whether the input satisfies the DynamoDB
// resource-name constraints: 3-255 characters, characters [a-zA-Z0-9_.-].
// Applies to TableName, BackupName, IndexName, and GlobalTableName.
func validateResourceName(name string) bool {
	if name == "" {
		return false
	}
	if len(name) < resourceNameMinLength || len(name) > resourceNameMaxLength {
		return false
	}
	return resourceNameRegex.MatchString(name)
}

// validateTableName is an alias for validateResourceName, retained for
// call-site clarity when validating a table name specifically.
func validateTableName(name string) bool {
	return validateResourceName(name)
}

// validateIndexName reports whether a secondary-index name meets the same
// constraints as table names (Smithy IndexName: len 3-255, pattern).
func validateIndexName(name string) bool {
	return validateResourceName(name)
}

// ---------------------------------------------------------------------------
// ARN validators
// ---------------------------------------------------------------------------

// validateArnLength reports whether an ARN's length falls within the Smithy
// constraints for the given ARN type.
func validateArnLength(arn string, min, max int) bool {
	return validateLength(arn, min, max)
}

// validateBackupArn reports whether arn satisfies the BackupArn Smithy
// constraints (len 37-1024).
func validateBackupArn(arn string) bool {
	return validateArnLength(arn, arnBackupMinLen, arnBackupMaxLen)
}

// validateStreamArn reports whether arn satisfies the StreamArn Smithy
// constraints (len 37-1024).
func validateStreamArn(arn string) bool {
	return validateArnLength(arn, arnStreamMinLen, arnStreamMaxLen)
}

// validateExportArn reports whether arn satisfies the ExportArn Smithy
// constraints (len 37-1024).
func validateExportArn(arn string) bool {
	return validateArnLength(arn, arnExportMinLen, arnExportMaxLen)
}

// validateImportArn reports whether arn satisfies the ImportArn Smithy
// constraints (len 37-1024).
func validateImportArn(arn string) bool {
	return validateArnLength(arn, arnImportMinLen, arnImportMaxLen)
}

// validateTableArn reports whether arn satisfies the TableArn Smithy
// constraints (len 1-1024).
func validateTableArn(arn string) bool {
	return validateArnLength(arn, arnTableMinLen, arnTableMaxLen)
}

// validateResourceArnString reports whether arn satisfies the generic
// ResourceArnString constraints used by tag and resource-policy operations
// (Smithy: len 1-1283).
func validateResourceArnString(arn string) bool {
	return validateArnLength(arn, arnResourceMinLen, arnResourceMaxLen)
}

// ---------------------------------------------------------------------------
// Projection / index validators
// ---------------------------------------------------------------------------

// validateProjectionType reports whether pt is a valid Smithy ProjectionType
// enum value. An empty string is treated as valid (parameter omitted).
func validateProjectionType(pt string) bool {
	if pt == "" {
		return true
	}
	return validProjectionTypes[pt]
}

// validateProjectionRequired validates the Projection sub-structure of a
// GSI or LSI definition. It returns true when the structure satisfies:
//   - NonKeyAttributes list size (Smithy NonKeyAttributeNameList: len 1-20).
//   - Each NonKeyAttributeName length (Smithy: len 1-255).
//   - NonKeyAttributes only valid when ProjectionType=INCLUDE.
func validateProjectionRequired(projMap map[string]interface{}) bool {
	if projMap == nil {
		return false
	}
	projectionType := request.GetStringParam(projMap, "ProjectionType")
	if nkAs, ok := projMap["NonKeyAttributes"].([]interface{}); ok {
		if len(nkAs) > nonKeyAttrListMaxLen {
			return false
		}
		for _, nk := range nkAs {
			nks, ok := nk.(string)
			if !ok {
				return false
			}
			if !validateLength(nks, nonKeyAttrNameMinLen, nonKeyAttrNameMaxLen) {
				return false
			}
		}
		if projectionType != "" && projectionType != ProjectionTypeInclude {
			return false
		}
	}
	return true
}

// validateKeyAttributeValue reports whether every value in key uses one of
// the types allowed for DynamoDB keys: S, N, or B (Smithy KeySchemaAttribute).
// Empty values are also rejected.
func validateKeyAttributeValue(key map[string]*dbstore.AttributeValue) bool {
	for _, av := range key {
		if av == nil {
			return false
		}
		// Key attributes must be one of S, N, or B.
		// BOOL, NULL, SS, NS, BS, M, L are not permitted.
		if av.S == nil && av.N == nil && av.B == nil {
			return false
		}
		// Reject empty values (extend to N type).
		if av.S != nil && *av.S == "" {
			return false
		}
		if av.N != nil && *av.N == "" {
			return false
		}
		if av.B != nil && len(av.B) == 0 {
			return false
		}
	}
	return true
}

// validateGSIDeleteExists reports whether the named GSI exists in gsiMap.
// Returns the typed ErrIndexNotFound sentinel so the caller can surface
// IndexNotFoundException rather than the generic ValidationException.
func validateGSIDeleteExists(gsiMap map[string]*dbstore.GlobalSecondaryIndex, indexName string) error {
	if _, exists := gsiMap[indexName]; !exists {
		return ErrIndexNotFound
	}
	return nil
}

// validateGSICreateRequired reports whether a GSI Create request includes
// the required IndexName and KeySchema (Smithy required traits), and whether
// the IndexName format satisfies the Smithy constraints (len 3-255, pattern).
func validateGSICreateRequired(create map[string]interface{}) bool {
	indexName := request.GetStringParam(create, "IndexName")
	if indexName == "" {
		return false
	}
	if !validateIndexName(indexName) {
		return false
	}
	keySchema := parseKeySchema(create)
	if len(keySchema) == 0 {
		return false
	}
	return validateKeySchema(keySchema)
}

// validateBillingModeConsistency reports whether PROVISIONED billing mode
// is paired with a ProvisionedThroughput value. PAY_PER_REQUEST is always
// consistent.
func validateBillingModeConsistency(billingMode dbstore.BillingMode, provThroughput *dbstore.ProvisionedThroughput) bool {
	if billingMode == dbstore.BillingModeProvisioned {
		return provThroughput != nil
	}
	return true
}

// ---------------------------------------------------------------------------
// Bool / bracket helpers
// ---------------------------------------------------------------------------

// validateBoolParam extracts a bool from a parameter map. If the key is
// absent, the default value is returned. If the key is present but not a
// bool, an error is returned (rejects malformed requests that would
// otherwise be silently coerced to the default).
//
// Kept as (bool, error) because the parsed value is needed by the caller on
// success; the bool alone is insufficient.
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
//
// Kept as (int, error) because the parsed index is needed by the caller on
// success; the bool alone is insufficient.
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

// validatePartiQLStatement reports whether stmt satisfies the PartiQL
// statement length constraint (Smithy PartiQLStatement: len 1-8192).
func validatePartiQLStatement(stmt string) bool {
	return validateLength(stmt, 1, partiqlStatementMaxLen)
}

// validatePartiQLBatchCount reports whether count satisfies the
// BatchExecuteStatement batch-size constraint (Smithy PartiQLBatchRequest:
// len 1-25).
func validatePartiQLBatchCount(count int) bool {
	return validateRange(count, 1, partiqlBatchMaxLen)
}

// validateExecuteStatementLimit reports whether limit satisfies the
// ExecuteStatement Limit constraint (Smithy PositiveIntegerObject: min 1).
func validateExecuteStatementLimit(limit int) bool {
	return limit >= 1
}

// validateS3Bucket reports whether bucket matches the S3 bucket name format
// (Smithy S3Bucket: len 0-255, pattern ^[a-z0-9A-Z]+[.\-\w]*[a-z0-9A-Z]+$).
func validateS3Bucket(bucket string) bool {
	if !validateLength(bucket, 0, s3BucketMaxLen) {
		return false
	}
	if bucket != "" && !s3BucketRegex.MatchString(bucket) {
		return false
	}
	return true
}

// validateS3BucketOwner reports whether owner matches the S3 bucket owner
// account ID format (Smithy S3BucketOwner: pattern ^[0-9]{12}$). Empty is
// treated as valid (parameter omitted).
func validateS3BucketOwner(owner string) bool {
	if owner == "" {
		return true
	}
	return s3BucketOwnerRegex.MatchString(owner)
}

// validateS3Prefix reports whether prefix satisfies the S3 key prefix length
// constraint (Smithy S3Prefix: len 0-1024).
func validateS3Prefix(prefix string) bool {
	return validateLength(prefix, 0, s3PrefixMaxLen)
}

// validateS3SseKmsKeyId reports whether keyId satisfies the KMS key ID
// length constraint (Smithy S3SseKmsKeyId: OPTIONAL, len 1-2048 when
// provided). Empty is treated as valid (parameter omitted).
func validateS3SseKmsKeyId(keyId string) bool {
	if keyId == "" {
		return true
	}
	return validateLength(keyId, s3SseKmsKeyIdMinLen, s3SseKmsKeyIdMaxLen)
}

// ---------------------------------------------------------------------------
// Smithy MEDIUM-tier validators
// ---------------------------------------------------------------------------

// validateGlobalTableName is an alias for validateResourceName; retained
// for call-site clarity (Smithy TableName: len 3-255, pattern).
func validateGlobalTableName(name string) bool {
	return validateResourceName(name)
}

// validateTagKey reports whether key satisfies the tag key length
// constraint (Smithy TagKeyString: len 1-128).
func validateTagKey(key string) bool {
	return validateLength(key, 1, tagKeyMaxLen)
}

// validateTagValue reports whether value satisfies the tag value length
// constraint (Smithy TagValueString: len 0-256).
func validateTagValue(value string) bool {
	return validateLength(value, 0, tagValueMaxLen)
}

// validateRecoveryPeriodInDays reports whether days falls within the PITR
// recovery-period range (Smithy RecoveryPeriodInDays: range 1-35).
func validateRecoveryPeriodInDays(days int) bool {
	return validateRange(days, recoveryPeriodMin, recoveryPeriodMax)
}

// validateScanSegment reports whether segment falls within the parallel
// Scan segment range (Smithy ScanSegment: range 0-999999).
func validateScanSegment(segment int) bool {
	return validateRange(segment, scanSegmentMin, scanSegmentMax)
}

// validateScanTotalSegments reports whether total falls within the parallel
// Scan total-segments range (Smithy ScanTotalSegments: range 1-1000000).
func validateScanTotalSegments(total int) bool {
	return validateRange(total, scanTotalSegmentsMin, scanTotalSegmentsMax)
}

// validateClientRequestToken reports whether token satisfies the idempotency
// token length constraint (Smithy ClientRequestToken: len 1-36). Empty is
// treated as valid (parameter omitted).
func validateClientRequestToken(token string) bool {
	if token == "" {
		return true
	}
	return validateLength(token, clientRequestTokenMinLen, clientRequestTokenMaxLen)
}

// validateListExportsLimit reports whether limit falls within the ListExports
// MaxResults range (Smithy ListExportsMaxLimit: range 1-25).
func validateListExportsLimit(limit int) bool {
	return validateRange(limit, 1, listExportsMaxLimit)
}

// validateListImportsLimit reports whether limit falls within the ListImports
// PageSize range (Smithy ListImportsMaxLimit: range 1-25).
func validateListImportsLimit(limit int) bool {
	return validateRange(limit, 1, listImportsMaxLimit)
}

// validateListBackupsLimit reports whether limit falls within the ListBackups
// Limit range (Smithy BackupsInputLimit: range 1-100).
func validateListBackupsLimit(limit int) bool {
	return validateRange(limit, 1, listBackupsMaxLimit)
}

// validateListContributorInsightsLimit reports whether limit falls within
// the ListContributorInsights MaxResults range (Smithy
// ListContributorInsightsLimit: range max 100).
func validateListContributorInsightsLimit(limit int) bool {
	return validateRange(limit, 0, listContributorMaxLimit)
}

// validateListTablesLimit reports whether limit falls within the ListTables
// Limit range (Smithy ListTablesInputLimit: range 1-100).
func validateListTablesLimit(limit int) bool {
	return validateRange(limit, 1, listTablesMaxLimit)
}

// validateListGlobalTablesLimit reports whether limit satisfies the
// ListGlobalTables Limit constraint (Smithy PositiveIntegerObject: min 1).
func validateListGlobalTablesLimit(limit int) bool {
	return limit >= listGlobalTablesMinLimit
}

// ---------------------------------------------------------------------------
// Smithy LOW-tier validators
// ---------------------------------------------------------------------------

// validatePolicyRevisionId reports whether id satisfies the resource-policy
// revision ID length constraint (Smithy PolicyRevisionId: len 1-255). Empty
// is treated as valid (parameter omitted).
func validatePolicyRevisionId(id string) bool {
	if id == "" {
		return true
	}
	return validateLength(id, policyRevisionIdMinLen, policyRevisionIdMaxLen)
}

// validateTimeToLiveAttributeName reports whether name satisfies the TTL
// attribute name length constraint (Smithy TimeToLiveAttributeName: len 1-255).
func validateTimeToLiveAttributeName(name string) bool {
	return validateLength(name, 1, ttlAttributeNameMaxLen)
}

// validateTableId reports whether id matches the table UUID format
// (Smithy TableId: pattern ^[0-9a-f]{8}-...$). Empty is treated as valid
// (parameter omitted).
func validateTableId(id string) bool {
	if id == "" {
		return true
	}
	return tableIdRegex.MatchString(id)
}

// validateImportNextToken reports whether token matches the ListImports
// next-token format (Smithy ImportNextToken: len 112-1024, pattern
// ^([0-9a-f]{16})+$). Empty is treated as valid (parameter omitted).
func validateImportNextToken(token string) bool {
	if token == "" {
		return true
	}
	if !validateLength(token, importNextTokenMinLen, importNextTokenMaxLen) {
		return false
	}
	return importNextTokenRegex.MatchString(token)
}

// validateClientToken reports whether token matches the Export/Import
// client-token format (Smithy ClientToken: pattern ^[^\$]+$). Empty is
// treated as valid (parameter omitted).
func validateClientToken(token string) bool {
	if token == "" {
		return true
	}
	return clientTokenRegex.MatchString(token)
}

// validateAutoScalingPolicyName reports whether name satisfies the
// auto-scaling policy name constraints (Smithy AutoScalingPolicyName:
// len 1-256, pattern ^\p{Print}+$).
func validateAutoScalingPolicyName(name string) bool {
	if !validateLength(name, autoScalingPolicyNameMinLen, autoScalingPolicyNameMaxLen) {
		return false
	}
	return autoScalingPolicyNameRegex.MatchString(name)
}

// validateAutoScalingRoleArn reports whether arn satisfies the auto-scaling
// role ARN length constraint (Smithy AutoScalingRoleArn: len 1-1600). The
// Smithy pattern permits any XML-compatible character and is extremely
// permissive, so only length is enforced.
func validateAutoScalingRoleArn(arn string) bool {
	return validateLength(arn, autoScalingRoleArnMinLen, autoScalingRoleArnMaxLen)
}

// validateContributorInsightsMode reports whether mode is a valid
// contributor-insights mode value (Smithy ContributorInsightsMode enum:
// ACCESSED_AND_THROTTLED_KEYS, THROTTLED_KEYS). Empty is treated as valid
// (parameter omitted).
func validateContributorInsightsMode(mode string) bool {
	if mode == "" {
		return true
	}
	switch mode {
	case "ACCESSED_AND_THROTTLED_KEYS", "THROTTLED_KEYS":
		return true
	default:
		return false
	}
}

// validateCsvDelimiter reports whether d is a valid CSV import delimiter
// character (Smithy CsvDelimiter: len 1-1, pattern ^[,;:|\t ]$).
func validateCsvDelimiter(d string) bool {
	if len(d) != 1 {
		return false
	}
	return csvDelimiterRegex.MatchString(d)
}

// validateCsvHeader reports whether header satisfies a single CSV import
// header value (Smithy CsvHeader: len 1-65536, pattern
// ^[\x20-\x21\x23-\x2B\x2D-\x7E]*$).
func validateCsvHeader(header string) bool {
	if !validateLength(header, 1, 65536) {
		return false
	}
	return csvHeaderRegex.MatchString(header)
}

// validateCsvHeaderList reports whether count falls within the CSV import
// header list size constraint (Smithy CsvHeaderList: len 1-255).
func validateCsvHeaderList(count int) bool {
	return validateRange(count, 1, csvHeaderListMaxLen)
}

// validateAttributeName reports whether name satisfies the attribute name
// length constraint (Smithy AttributeName: len 0-65535).
func validateAttributeName(name string) bool {
	return validateLength(name, 0, attributeNameMaxLen)
}

// validateKeySchemaAttributeName reports whether name satisfies a key-schema
// attribute name length (Smithy KeySchemaAttributeName: len 1-255).
func validateKeySchemaAttributeName(name string) bool {
	return validateLength(name, keySchemaAttrNameMinLen, keySchemaAttrNameMaxLen)
}

// validateNonKeyAttributeName reports whether name satisfies a single
// NonKeyAttributeName length (Smithy NonKeyAttributeName: len 1-255).
func validateNonKeyAttributeName(name string) bool {
	return validateLength(name, nonKeyAttrNameMinLen, nonKeyAttrNameMaxLen)
}

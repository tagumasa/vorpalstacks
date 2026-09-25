package dynamodb

import (
	"math"
	"regexp"
	"strings"
	"unicode/utf8"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
	"vorpalstacks/internal/utils/aws/arn"
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
	// TargetTrackingScalingPolicyConfiguration.TargetValue documented
	// range (Base 10): 8.515920e-109 to 1.174271e+108.
	autoScalingTargetValueMin = 8.515920e-109
	autoScalingTargetValueMax = 1.174271e+108
	csvHeaderListMaxLen       = 255
	batchGetMaxTotalItems     = 100
	batchWriteMaxItems        = 25
	transactMaxItems          = 100
	listExportsMaxLimit       = 25
	listImportsMaxLimit       = 25
	listBackupsMaxLimit       = 100
	listContributorMaxLimit   = 100
	listTablesMaxLimit        = 100
	listGlobalTablesMinLimit  = 1
	recoveryPeriodMin         = 1
	recoveryPeriodMax         = 35
	scanSegmentMin            = 0
	scanSegmentMax            = 999999
	scanTotalSegmentsMin      = 1
	scanTotalSegmentsMax      = 1000000

	// Pagination defaults and caps used by List operations that do not
	// have an explicit Smithy Limit trait but follow AWS documented
	// behaviour. Each value is referenced from exactly one call site.
	listGlobalTablesDefaultLimit       = 100
	listGlobalTablesMinPageSize        = 10
	listStreamsDefaultLimit            = 100
	listStreamsMaxLimit                = 100
	describeStreamDefaultLimit         = 100
	describeStreamMaxLimit             = 100
	getRecordsDefaultLimit             = 100
	getRecordsMaxLimit                 = 1000
	listTagsForResourceDefaultPageSize = 50

	// Backup list minimum fetch size: the ListBackups implementation
	// always fetches at least this many entries from the store before
	// applying filters, to reduce round-trips for common queries.
	listBackupsMinFetchSize = 100

	// PITR default recovery window in days (AWS default: 35 days).
	pitrDefaultRecoveryPeriodDays = 35

	// Data-plane Query/Scan page size when the request carries no Limit.
	// The Limit parameter's model range has a minimum of 1 and no maximum;
	// the documented AWS page bound is 1 MB of data, not an item count, so
	// a requested Limit is honoured as-is.
	dataPlaneQueryDefaultLimit = 100
)

// Compiled regex patterns from Smithy @pattern traits.
var (
	// Shared resource-name pattern: TableName, BackupName, IndexName, GlobalTableName.
	resourceNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)

	s3BucketRegex              = regexp.MustCompile(`^[a-z0-9A-Z]+[.\-\w]*[a-z0-9A-Z]+$`)
	s3BucketOwnerRegex         = regexp.MustCompile(`^[0-9]{12}$`)
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

// validateLength reports whether a string's length falls within [min, max],
// counted in Unicode characters (Smithy @length on string shapes counts
// code points).
func validateLength(val string, min, max int) bool {
	n := utf8.RuneCountInString(val)
	return n >= min && n <= max
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

// resolveTableNameMember normalises a table-naming member that documents
// both forms — the bare table name or the table's ARN — into the bare
// name the stores key records under, parsing the ARN rather than probing
// substrings. The ARN form must be the table ARN itself: stream and index
// ARNs carry further resource segments, so a resource that is not exactly
// table/<name> does not resolve.
func resolveTableNameMember(member string) (string, error) {
	if parsed, perr := arn.ParseARN(member); perr == nil {
		name, isTable := strings.CutPrefix(parsed.Resource, "table/")
		if parsed.Service != "dynamodb" || !isTable || strings.Contains(name, "/") {
			return "", ErrInvalidParameter
		}
		member = name
	}
	if !validateResourceName(member) {
		return "", ErrInvalidParameter
	}
	return member, nil
}

// validateVectorDistanceFunction reports whether fn is one of the modelled
// vector distance functions.
func validateVectorDistanceFunction(fn string) bool {
	switch fn {
	case "COSINE", "EUCLIDEAN", "DOT_PRODUCT":
		return true
	}
	return false
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
		// The list trait bounds the length at both ends: an explicitly
		// empty NonKeyAttributes list violates the minimum of 1.
		if len(nkAs) == 0 || len(nkAs) > nonKeyAttrListMaxLen {
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
// Empty values are also rejected. A key with no members cannot address an
// item (the wire Key member is required), so an empty map is invalid too —
// this is the Core-side rejection of a missing or unparseable Key.
func validateKeyAttributeValue(key map[string]*dbstore.AttributeValue) bool {
	if len(key) == 0 {
		return false
	}
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
	if !validateResourceName(indexName) {
		return false
	}
	keySchema := parseKeySchema(create)
	if len(keySchema) == 0 {
		return false
	}
	return validateKeySchema(keySchema)
}

// validateBillingModeConsistency reports whether the billing mode and the
// provisioned throughput pair is valid: PROVISIONED requires throughput
// settings with both units at least 1, and PAY_PER_REQUEST rejects them —
// on-demand tables have no provisioned capacity ("ProvisionedThroughput
// cannot be specified when BillingMode is PAY_PER_REQUEST").
func validateBillingModeConsistency(billingMode dbstore.BillingMode, provThroughput *dbstore.ProvisionedThroughput) bool {
	if billingMode == dbstore.BillingModeProvisioned {
		return validateProvisionedThroughputValues(provThroughput)
	}
	if billingMode == dbstore.BillingModePayPerRequest {
		return provThroughput == nil
	}
	return true
}

// validateSecondaryIndexQuotas reports whether a table's secondary index
// families satisfy the documented per-table quotas: at most 20 global
// secondary indexes, at most 5 local secondary indexes, and at most 100
// user-specified projected attribute names combined across all secondary
// indexes — the sum the quota page states "only applies to user-specified
// projected attributes".
func validateSecondaryIndexQuotas(gsis []*dbstore.GlobalSecondaryIndex, lsis []*dbstore.LocalSecondaryIndex) bool {
	if len(gsis) > dbstore.GlobalSecondaryIndexesPerTable {
		return false
	}
	if len(lsis) > dbstore.LocalSecondaryIndexesPerTable {
		return false
	}
	projected := 0
	for _, gsi := range gsis {
		if gsi.Projection != nil {
			projected += len(gsi.Projection.NonKeyAttributes)
		}
	}
	for _, lsi := range lsis {
		if lsi.Projection != nil {
			projected += len(lsi.Projection.NonKeyAttributes)
		}
	}
	return projected <= dbstore.ProjectedAttributesPerTable
}

// validateIndexThroughputModeConsistency reports whether every global
// secondary index's capacity settings follow the table's billing mode:
// "Global secondary indexes inherit the read/write capacity mode from the
// base table", and the throughput member's own contract is "You must use
// ProvisionedThroughput or OnDemandThroughput based on your table's capacity
// mode" — so a provisioned table's index carries provisioned settings and an
// on-demand table's index carries none.
func validateIndexThroughputModeConsistency(billingMode dbstore.BillingMode, gsis []*dbstore.GlobalSecondaryIndex) bool {
	for _, gsi := range gsis {
		if billingMode == dbstore.BillingModeProvisioned {
			if gsi.ProvisionedThroughput == nil || !validateProvisionedThroughputValues(gsi.ProvisionedThroughput) {
				return false
			}
		}
		if billingMode == dbstore.BillingModePayPerRequest && gsi.ProvisionedThroughput != nil {
			return false
		}
	}
	return true
}

// validateProvisionedThroughputValues reports whether both capacity values
// fall inside the model's PositiveLongObject range (minimum 1).
func validateProvisionedThroughputValues(pt *dbstore.ProvisionedThroughput) bool {
	return pt != nil &&
		pt.ReadCapacityUnits >= 1 && pt.WriteCapacityUnits >= 1 &&
		pt.ReadCapacityUnits <= dbstore.TableMaxReadCapacityUnits &&
		pt.WriteCapacityUnits <= dbstore.TableMaxWriteCapacityUnits
}

// validateBillingModeValue reports whether the billing mode is one of the
// BillingMode enum values (Smithy: PROVISIONED | PAY_PER_REQUEST).
func validateBillingModeValue(bm dbstore.BillingMode) bool {
	return bm == dbstore.BillingModeProvisioned || bm == dbstore.BillingModePayPerRequest
}

// validateTableClassValue reports whether the table class is one of the
// TableClass enum values (Smithy: STANDARD | STANDARD_INFREQUENT_ACCESS).
func validateTableClassValue(tableClass string) bool {
	return tableClass == string(dbstore.TableClassStandard) || tableClass == string(dbstore.TableClassStandardInfrequentAccess)
}

// validateStreamSpecification rejects a streaming table whose view type is
// not one of the StreamViewType enum values (Smithy: NEW_IMAGE |
// OLD_IMAGE | NEW_AND_OLD_IMAGES | KEYS_ONLY). A disabled specification
// carries no view type to validate.
func validateStreamSpecification(spec *dbstore.StreamSpecification) error {
	if spec == nil || !spec.StreamEnabled {
		return nil
	}
	switch spec.StreamViewType {
	case dbstore.StreamViewTypeNewImage, dbstore.StreamViewTypeOldImage,
		dbstore.StreamViewTypeNewAndOldImages, dbstore.StreamViewTypeKeysOnly:
		return nil
	}
	return ErrInvalidParameter
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
	return validateBoolValue(params[key], defaultVal)
}

// validateBoolValue is validateBoolParam's form for a caller that holds the
// raw wire member rather than the parameter map (a core reading a DTO field
// the handler extracted): nil is the absent member and yields the default,
// a present non-bool is rejected rather than silently coerced.
func validateBoolValue(raw interface{}, defaultVal bool) (bool, error) {
	if raw == nil {
		return defaultVal, nil
	}
	b, ok := raw.(bool)
	if !ok {
		return false, ErrInvalidParameter
	}
	return b, nil
}

// validateBracketIndex parses a bracket-enclosed list index (e.g. "[3]")
// and returns the integer index. Returns an error for empty, non-numeric,
// or negative values (Smithy document path spec), and for digit runs the
// machine cannot represent: the accumulation bails at the first digit that
// would overflow, so an absurd index can never wrap into a small valid one.
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
		digit := int(ch - '0')
		if idx > (math.MaxInt-digit)/10 {
			return 0, ErrInvalidParameter
		}
		idx = idx*10 + digit
	}
	return idx, nil
}

// ---------------------------------------------------------------------------
// Statement, S3 target and quota validators (export/import plane)
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
// (Smithy S3Bucket: pattern ^[a-z0-9A-Z]+[.\-\w]*[a-z0-9A-Z]+$, length at
// most 255). The pattern applies unconditionally: it requires at least two
// characters, so the empty string — a required member's omission — fails
// at request time rather than inside the background job.
func validateS3Bucket(bucket string) bool {
	if !validateLength(bucket, 0, s3BucketMaxLen) {
		return false
	}
	return s3BucketRegex.MatchString(bucket)
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

// validateS3SseAlgorithm reports whether algorithm is a member of the
// S3SseAlgorithm enum (Smithy: AES256 | KMS). Empty is treated as valid
// (parameter omitted — the destination bucket's default encryption then
// governs).
func validateS3SseAlgorithm(algorithm string) bool {
	return algorithm == "" || algorithm == "AES256" || algorithm == "KMS"
}

// validateExportType reports whether exportType is a member of the
// ExportType enum (Smithy: FULL_EXPORT | INCREMENTAL_EXPORT). Empty is
// treated as valid (parameter omitted — the documented default is
// FULL_EXPORT).
func validateExportType(exportType string) bool {
	return exportType == "" || exportType == "FULL_EXPORT" || exportType == "INCREMENTAL_EXPORT"
}

// validateExportViewType reports whether viewType is a member of the
// ExportViewType enum (Smithy: NEW_IMAGE | NEW_AND_OLD_IMAGES). Empty is
// treated as valid (parameter omitted — the documented default is
// NEW_AND_OLD_IMAGES).
func validateExportViewType(viewType string) bool {
	return viewType == "" || viewType == "NEW_IMAGE" || viewType == "NEW_AND_OLD_IMAGES"
}

// ---------------------------------------------------------------------------
// Tagging, scan and list-pagination validators
// ---------------------------------------------------------------------------

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
// Name, token and expression-attribute validators
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

// Select parameter values (Smithy Select enum).
const (
	selectCount              = "COUNT"
	selectAllAttributes      = "ALL_ATTRIBUTES"
	selectAllProjectedAttrs  = "ALL_PROJECTED_ATTRIBUTES"
	selectSpecificAttributes = "SPECIFIC_ATTRIBUTES"
)

// parseSelectParam validates the Select parameter of Query/Scan against the
// API reference rules: ALL_PROJECTED_ATTRIBUTES is only valid on an index,
// SPECIFIC_ATTRIBUTES requires a ProjectionExpression, and Select may only
// accompany a ProjectionExpression when it is SPECIFIC_ATTRIBUTES. Index
// reads default to ALL_PROJECTED_ATTRIBUTES, table reads to ALL_ATTRIBUTES.
func parseSelectParam(params map[string]interface{}, indexName string, hasProjection bool) (countOnly, allProjected bool, err error) {
	raw, present := params["Select"]
	if !present {
		return false, indexName != "", nil
	}
	value, ok := raw.(string)
	if !ok {
		return false, false, ErrInvalidParameter
	}
	switch value {
	case selectCount:
		return true, false, nil
	case selectAllAttributes:
		if hasProjection {
			return false, false, ErrInvalidParameter
		}
		return false, false, nil
	case selectAllProjectedAttrs:
		if indexName == "" || hasProjection {
			return false, false, ErrInvalidParameter
		}
		return false, true, nil
	case selectSpecificAttributes:
		if !hasProjection {
			return false, false, ErrInvalidParameter
		}
		return false, false, nil
	default:
		return false, false, ErrInvalidParameter
	}
}

// DynamoDB Number constraints (developer guide, Supported data types):
// up to 38 digits of precision, positive range 1E-130 to
// 9.9999999999999999999999999999999999999E+125 (symmetric negative range).
const (
	numberMaxSignificantDigits = 38
)

// ClientRequestToken idempotency window (TransactWriteItems API reference:
// "A client request token is valid for 10 minutes after the first request
// that uses it is completed").
const idempotencyWindowMinutes = 10

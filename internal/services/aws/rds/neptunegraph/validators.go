package neptunegraph

import (
	"fmt"
	"regexp"
	"strings"

	ngstore "vorpalstacks/internal/store/aws/rds/neptunegraph"
)

// ---------------------------------------------------------------------------
// Smithy-derived constants
// ---------------------------------------------------------------------------

const (
	minProvisionedMemory = 16
	maxProvisionedMemory = 24576
	maxReplicaCount      = 2
	minVectorSearchDim   = 1
	maxVectorSearchDim   = 65536
	maxTags              = 50
	maxTagKeyLen         = 128
	maxTagValueLen       = 256
	maxKmsKeyArnLen      = 1024
	maxDestinationLen    = 1024
	planCacheTTLSeconds  = 300
	planCacheCapacity    = 1000
)

// ---------------------------------------------------------------------------
// Smithy-derived patterns
// ---------------------------------------------------------------------------

var (
	graphNameRegex    = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)
	snapshotNameRegex = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)
	tagKeyRegex       = regexp.MustCompile(`^[a-zA-Z+\-=._:/]+$`)
	kmsKeyArnPattern  = regexp.MustCompile(`^arn:aws[^:]*:kms:[a-zA-Z0-9-]*:[0-9]{12}:key/[a-zA-Z0-9-]{36}$`)
	roleArnPattern    = regexp.MustCompile(`^arn:aws[^:]*:iam::\d{12}:(role|role/service-role)(/[\w+=,.@-]+)+$`)
)

// ---------------------------------------------------------------------------
// Smithy-derived enum sets
// ---------------------------------------------------------------------------

var validQueryLanguage = map[string]bool{
	"OPEN_CYPHER": true,
}

var validPlanCacheType = map[string]bool{
	"ENABLED":  true,
	"DISABLED": true,
	"AUTO":     true,
}

var validExplainMode = map[string]bool{
	"STATIC":  true,
	"DETAILS": true,
}

var validExportFormat = map[string]bool{
	"CSV":     true,
	"PARQUET": true,
}

var validImportFormat = map[string]bool{
	"CSV":         true,
	"OPEN_CYPHER": true,
	"PARQUET":     true,
	"NTRIPLES":    true,
}

var validParquetType = map[string]bool{
	"COLUMNAR": true,
}

var validBlankNodeHandling = map[string]bool{
	"convertToIri": true,
}

var validGraphSummaryMode = map[string]bool{
	"BASIC":    true,
	"DETAILED": true,
}

var validQueryStateInput = map[string]bool{
	"ALL":        true,
	"RUNNING":    true,
	"WAITING":    true,
	"CANCELLING": true,
}

// ---------------------------------------------------------------------------
// Validators — each returns *httpError (nil on success)
// ---------------------------------------------------------------------------

// validateGraphName validates a graph name per the Smithy GraphName shape:
// @length(min:1, max:63) @pattern(^(?!g-)[a-z][a-z0-9]*(-[a-z0-9]+)*$).
func validateGraphName(name string) error {
	if name == "" || strings.HasPrefix(name, "g-") || !graphNameRegex.MatchString(name) || len(name) > 63 {
		return newValidationException("ILLEGAL_ARGUMENT", "graphName must be 1-63 chars, lowercase letters, digits, hyphens, and must not start with 'g-'")
	}
	return nil
}

// validateSnapshotName validates a snapshot name per the Smithy SnapshotName shape:
// @length(min:1, max:63) @pattern(^(?!gs-)[a-z][a-z0-9]*(-[a-z0-9]+)*$).
func validateSnapshotName(name string) error {
	if name == "" || strings.HasPrefix(name, "gs-") || !snapshotNameRegex.MatchString(name) || len(name) > 63 {
		return newValidationException("ILLEGAL_ARGUMENT", "snapshotName must be 1-63 chars, lowercase letters, digits, hyphens, and must not start with 'gs-'")
	}
	return nil
}

// validateProvisionedMemory validates provisionedMemory per the Smithy
// ProvisionedMemory shape: @range(min:16, max:24576).
// When required is true (CreateGraph), an unset value is rejected.
func validateProvisionedMemory(mem int, required bool) error {
	if required && mem == 0 {
		return newValidationException("ILLEGAL_ARGUMENT", "provisionedMemory is required")
	}
	if mem != 0 && (mem < minProvisionedMemory || mem > maxProvisionedMemory) {
		return newValidationException("CONSTRAINT_VIOLATION", fmt.Sprintf("provisionedMemory must be between %d and %d", minProvisionedMemory, maxProvisionedMemory))
	}
	return nil
}

// validateReplicaCount validates replicaCount per the Smithy ReplicaCount
// shape: @range(min:0, max:2).
func validateReplicaCount(rc int) error {
	if rc < 0 || rc > maxReplicaCount {
		return newValidationException("CONSTRAINT_VIOLATION", fmt.Sprintf("replicaCount must be between 0 and %d", maxReplicaCount))
	}
	return nil
}

// validateVectorSearchDimension validates the vector search dimension per the
// Smithy VectorSearchDimension shape: @range(min:1, max:65536).
func validateVectorSearchDimension(dim int) error {
	if dim < minVectorSearchDim || dim > maxVectorSearchDim {
		return newValidationException("CONSTRAINT_VIOLATION", fmt.Sprintf("vectorSearchConfiguration.dimension must be between %d and %d", minVectorSearchDim, maxVectorSearchDim))
	}
	return nil
}

// validateQueryLanguage validates the query language per the Smithy
// QueryLanguage enum: OPEN_CYPHER only. The field is required.
func validateQueryLanguage(lang string) error {
	if lang == "" {
		return newValidationException("ILLEGAL_ARGUMENT", "language is required")
	}
	if !validQueryLanguage[lang] {
		return newValidationException("ILLEGAL_ARGUMENT", fmt.Sprintf("unsupported language: %s (only OPEN_CYPHER is supported)", lang))
	}
	return nil
}

// validatePlanCache validates the planCache parameter per the Smithy
// PlanCacheType enum: ENABLED, DISABLED, AUTO. The field is optional.
func validatePlanCache(pc string) error {
	if pc == "" {
		return nil
	}
	if !validPlanCacheType[pc] {
		return newValidationException("ILLEGAL_ARGUMENT", fmt.Sprintf("invalid planCache: %s (must be ENABLED, DISABLED, or AUTO)", pc))
	}
	return nil
}

// validateExplainMode validates the explain mode per the Smithy ExplainMode
// enum: STATIC, DETAILS. The field is optional.
func validateExplainMode(mode string) error {
	if mode == "" {
		return nil
	}
	if !validExplainMode[mode] {
		return newValidationException("ILLEGAL_ARGUMENT", fmt.Sprintf("invalid explain mode: %s (must be STATIC or DETAILS)", mode))
	}
	return nil
}

// validateExportFormat validates the export format per the Smithy ExportFormat
// enum: PARQUET, CSV.
func validateExportFormat(format string) error {
	if format == "" {
		return newValidationException("ILLEGAL_ARGUMENT", "format is required")
	}
	if !validExportFormat[format] {
		return newValidationException("ILLEGAL_ARGUMENT", fmt.Sprintf("invalid format: %s (must be CSV or PARQUET)", format))
	}
	return nil
}

// validateImportFormat validates the import format per the Smithy Format enum:
// CSV, OPEN_CYPHER, PARQUET, NTRIPLES. The field is optional.
func validateImportFormat(format string) error {
	if format == "" {
		return nil
	}
	if !validImportFormat[format] {
		return newValidationException("ILLEGAL_ARGUMENT", fmt.Sprintf("invalid format: %s (must be CSV, OPEN_CYPHER, PARQUET, or NTRIPLES)", format))
	}
	return nil
}

// validateParquetType validates the parquet type per the Smithy ParquetType
// enum: COLUMNAR. The field is optional.
func validateParquetType(pt string) error {
	if pt == "" {
		return nil
	}
	if !validParquetType[pt] {
		return newValidationException("ILLEGAL_ARGUMENT", fmt.Sprintf("invalid parquetType: %s (must be COLUMNAR)", pt))
	}
	return nil
}

// validateBlankNodeHandling validates the blankNodeHandling per the Smithy
// BlankNodeHandling enum: convertToIri. The field is optional.
func validateBlankNodeHandling(bnh string) error {
	if bnh == "" {
		return nil
	}
	if !validBlankNodeHandling[bnh] {
		return newValidationException("ILLEGAL_ARGUMENT", fmt.Sprintf("invalid blankNodeHandling: %s (must be convertToIri)", bnh))
	}
	return nil
}

// validateGraphSummaryMode validates the summary mode per the Smithy
// GraphSummaryMode enum: BASIC, DETAILED. The field is optional.
func validateGraphSummaryMode(mode string) error {
	if mode == "" {
		return nil
	}
	if !validGraphSummaryMode[mode] {
		return newValidationException("ILLEGAL_ARGUMENT", fmt.Sprintf("invalid mode: %s (must be BASIC or DETAILED)", mode))
	}
	return nil
}

// validateQueryStateInput validates the query state filter per the Smithy
// QueryStateInput enum: ALL, RUNNING, WAITING, CANCELLING. The field is optional.
func validateQueryStateInput(state string) error {
	if state == "" {
		return nil
	}
	if !validQueryStateInput[state] {
		return newValidationException("ILLEGAL_ARGUMENT", fmt.Sprintf("invalid state: %s (must be ALL, RUNNING, WAITING, or CANCELLING)", state))
	}
	return nil
}

// validateKmsKeyArn validates a KMS key ARN per the Smithy KmsKeyArn shape:
// @length(min:1, max:1024) @pattern(^arn:aws[^:]*:kms:...).
// The field is optional unless required is true.
func validateKmsKeyArn(arn string, required bool) error {
	if arn == "" {
		if required {
			return newValidationException("ILLEGAL_ARGUMENT", "kmsKeyIdentifier is required")
		}
		return nil
	}
	if len(arn) > maxKmsKeyArnLen {
		return newValidationException("ILLEGAL_ARGUMENT", fmt.Sprintf("kmsKeyIdentifier must be at most %d characters", maxKmsKeyArnLen))
	}
	if !kmsKeyArnPattern.MatchString(arn) {
		return newValidationException("ILLEGAL_ARGUMENT", "kmsKeyIdentifier must be a valid KMS key ARN")
	}
	return nil
}

// validateRoleArn validates a role ARN per the Smithy RoleArn shape:
// @pattern(^arn:aws[^:]*:iam::\d{12}:(role|role/service-role)(/[\w+=,.@-]+)+$).
func validateRoleArn(arn string) error {
	if arn == "" {
		return newValidationException("ILLEGAL_ARGUMENT", "roleArn is required")
	}
	if !roleArnPattern.MatchString(arn) {
		return newValidationException("ILLEGAL_ARGUMENT", "roleArn must be a valid IAM role ARN")
	}
	return nil
}

// validateDestination validates the export destination per the Smithy
// shape: @length(min:1, max:1024).
func validateDestination(dest string) error {
	if dest == "" {
		return newValidationException("ILLEGAL_ARGUMENT", "destination is required")
	}
	if len(dest) > maxDestinationLen {
		return newValidationException("ILLEGAL_ARGUMENT", fmt.Sprintf("destination must be at most %d characters", maxDestinationLen))
	}
	return nil
}

// validateImportOptions validates NeptuneImportOptions required fields when
// the neptune sub-object is present.
func validateImportOptions(opts *ngstore.ImportOptions) error {
	if opts == nil || opts.Neptune == nil {
		return nil
	}
	n := opts.Neptune
	if n.S3ExportPath == "" {
		return newValidationException("ILLEGAL_ARGUMENT", "importOptions.neptune.s3ExportPath is required")
	}
	if len(n.S3ExportPath) > maxKmsKeyArnLen {
		return newValidationException("ILLEGAL_ARGUMENT", "importOptions.neptune.s3ExportPath must be 1-1024 characters")
	}
	if n.S3ExportKmsKeyId == "" {
		return newValidationException("ILLEGAL_ARGUMENT", "importOptions.neptune.s3ExportKmsKeyId is required")
	}
	if len(n.S3ExportKmsKeyId) > maxKmsKeyArnLen {
		return newValidationException("ILLEGAL_ARGUMENT", "importOptions.neptune.s3ExportKmsKeyId must be 1-1024 characters")
	}
	return nil
}

// validateTags validates a tag map per Smithy TagMap constraints:
// max 50 entries, TagKey @length(1,128) @pattern, TagValue @length(0,256).
func validateTags(tags map[string]string) error {
	if len(tags) > maxTags {
		return newValidationException("ILLEGAL_ARGUMENT", fmt.Sprintf("too many tags: %d (max %d)", len(tags), maxTags))
	}
	for k, v := range tags {
		if len(k) > maxTagKeyLen || !tagKeyRegex.MatchString(k) {
			return newValidationException("ILLEGAL_ARGUMENT", fmt.Sprintf("invalid tag key: %q", k))
		}
		if len(v) > maxTagValueLen {
			return newValidationException("ILLEGAL_ARGUMENT", fmt.Sprintf("tag value too long for key %q (max %d)", k, maxTagValueLen))
		}
	}
	return nil
}

// clampMaxResults clamps a maxResults value to the Smithy MaxResults range
// @range(min:1, max:100). Unset (0) defaults to 100.
func clampMaxResults(v int) int {
	if v < 1 {
		return 100
	}
	if v > 100 {
		return 100
	}
	return v
}

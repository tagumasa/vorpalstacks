package neptunedata

import "regexp"

// ---------------------------------------------------------------------------
// Smithy-derived constants
// ---------------------------------------------------------------------------

const (
	engineVersion = "1.3.3.0"
)

// ---------------------------------------------------------------------------
// Smithy-derived patterns
// ---------------------------------------------------------------------------

var (
	iamRoleArnPattern  = regexp.MustCompile(`^arn:aws:iam::\d*:role/.+`)
	s3SourceURIPattern = regexp.MustCompile(`^(s3://|https://s3\.)`)
)

// ---------------------------------------------------------------------------
// Smithy-derived enum sets
// ---------------------------------------------------------------------------

// Format enum (Smithy: CSV, NQUADS, NTRIPLES, OPENCYPHER, RDFXML, TURTLE).
// The AWS SDK sends lowercase values in the JSON body, matching the
// Neptune Data API wire convention.
var validLoaderFormats = map[string]bool{
	"csv":        true,
	"opencypher": true,
	"ntriples":   true,
	"nquads":     true,
	"rdfxml":     true,
	"turtle":     true,
}

// Mode enum (Smithy: AUTO, NEW, RESUME).
var validLoaderModes = map[string]bool{
	"RESUME": true,
	"NEW":    true,
	"AUTO":   true,
}

// Parallelism enum (Smithy: HIGH, LOW, MEDIUM, OVERSUBSCRIBE).
var validLoaderParallelism = map[string]bool{
	"LOW":           true,
	"MEDIUM":        true,
	"HIGH":          true,
	"OVERSUBSCRIBE": true,
}

// OpenCypherExplainMode enum (Smithy: DETAILS, DYNAMIC, STATIC).
// The AWS SDK sends lowercase values in the JSON body.
var validExplainModes = map[string]bool{
	"static":  true,
	"details": true,
	"dynamic": true,
}

// IteratorType enum (Smithy: AFTER_SEQUENCE_NUMBER, AT_SEQUENCE_NUMBER,
// LATEST, TRIM_HORIZON).
var validIteratorTypes = map[string]bool{
	"AT_SEQUENCE_NUMBER":    true,
	"AFTER_SEQUENCE_NUMBER": true,
	"LATEST":                true,
	"TRIM_HORIZON":          true,
}

// StatisticsAutoGenerationMode values used in the JSON body. The AWS SDK
// sends camelCase strings matching the Neptune Data API wire convention.
var validManageStatsModes = map[string]bool{
	"disableAutoCompute": true,
	"enableAutoCompute":  true,
	"refresh":            true,
}

// GraphSummaryType values used in the URL path parameter. The Neptune Data
// API accepts lowercase mode values.
var validGraphSummaryModes = map[string]bool{
	"basic":    true,
	"detailed": true,
}

// ---------------------------------------------------------------------------
// Validator functions
// ---------------------------------------------------------------------------

func validateLoaderFormat(format string) bool {
	return validLoaderFormats[format]
}

func validateLoaderMode(mode string) bool {
	return mode == "" || validLoaderModes[mode]
}

func validateLoaderParallelism(p string) bool {
	return p == "" || validLoaderParallelism[p]
}

func validateExplainMode(mode string) bool {
	return mode == "" || validExplainModes[mode]
}

func validateIteratorType(t string) bool {
	return t == "" || validIteratorTypes[t]
}

func validateManageStatsMode(mode string) bool {
	return validManageStatsModes[mode]
}

func validateGraphSummaryMode(mode string) bool {
	return mode == "" || validGraphSummaryModes[mode]
}

func validateIamRoleArn(arn string) bool {
	return iamRoleArnPattern.MatchString(arn)
}

func validateS3SourceURI(uri string) bool {
	return s3SourceURIPattern.MatchString(uri)
}

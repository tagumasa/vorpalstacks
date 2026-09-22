package cloudwatchlogs

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	tagutil "vorpalstacks/internal/common/tags"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The listing validators: the enum tables the filters and statuses
// read, the page-limit resolution, and the ListLogGroups and identifier
// families' member bounds.

// --- Enum maps ---

var (
	validLogGroupClasses = map[string]bool{
		"STANDARD":          true,
		"INFREQUENT_ACCESS": true,
		"DELIVERY":          true,
	}

	validPolicyTypes = map[string]bool{
		"DATA_PROTECTION_POLICY":     true,
		"SUBSCRIPTION_FILTER_POLICY": true,
		"FIELD_INDEX_POLICY":         true,
		"TRANSFORMER_POLICY":         true,
		"METRIC_EXTRACTION_POLICY":   true,
	}

	validDistributions = map[string]bool{
		"Random":      true,
		"ByLogStream": true,
	}

	validScheduledQueryStates = map[string]bool{
		"ENABLED":  true,
		"DISABLED": true,
	}

	// validQueryStatusFilter is the model's QueryStatus enum — the status
	// member of DescribeQueries.
	validQueryStatusFilter = map[string]bool{
		"Scheduled": true, "Running": true, "Complete": true, "Failed": true,
		"Cancelled": true, "Timeout": true, "Unknown": true,
	}

	// validExportStatusCode is the model's ExportTaskStatusCode enum — the
	// statusCode member of DescribeExportTasks.
	validExportStatusCode = map[string]bool{
		"CANCELLED": true, "COMPLETED": true, "FAILED": true,
		"PENDING": true, "PENDING_CANCEL": true, "RUNNING": true,
	}

	// validImportStatus is the model's ImportStatus enum — the importStatus
	// member of DescribeImportTasks.
	validImportStatus = map[string]bool{
		"IN_PROGRESS": true, "CANCELLED": true, "COMPLETED": true, "FAILED": true,
	}

	// validExecutionStatus is the model's ExecutionStatus enum — the
	// executionStatuses member of GetScheduledQueryHistory.
	validExecutionStatus = map[string]bool{
		"Running": true, "InvalidQuery": true, "Complete": true, "Failed": true, "Timeout": true,
	}
)

// --- Enum validators ---

func validateLogGroupClass(class string) bool {
	return validLogGroupClasses[class]
}

func validatePolicyType(t string) bool {
	return validPolicyTypes[t]
}

func validateDistribution(d string) bool {
	return validDistributions[d]
}

func validateScheduledQueryState(state string) bool {
	return validScheduledQueryStates[state]
}

// --- Range / limit validators ---

// resolveListLimit applies a default when limit is unset (0) and
// enforces the range trait's 1..max floor and ceiling, rejecting with the
// identity the calling operation family's declared error list carries. A
// negative value is a supplied value outside the range, not an absent
// member — it rejects rather than silently becoming the default page.
func resolveListLimit(limit, defaultVal, maxVal int32, rejectionCode string) (int32, error) {
	if limit == 0 {
		return defaultVal, nil
	}
	if limit < 0 || limit > maxVal {
		return 0, NewLogsError(rejectionCode,
			fmt.Sprintf("Limit must be between 1 and %d", maxVal), 400)
	}
	return limit, nil
}

// validateListLimit applies a default when limit is unset (0) and
// enforces the range trait's bounds, rejecting a negative value as a
// supplied value outside the range rather than defaulting it. Returns
// the resolved limit or an InvalidParameterException when the
// caller-supplied value is out of range — the parameter identity the
// legacy operation families declare.
func validateListLimit(limit, defaultVal, maxVal int32) (int32, error) {
	return resolveListLimit(limit, defaultVal, maxVal, "InvalidParameterException")
}

// validateListLimitValidation is the same bound check for the operation
// families whose declared error lists carry ValidationException and not
// InvalidParameterException — the vended-delivery family and the
// scheduled-query lists.
func validateListLimitValidation(limit, defaultVal, maxVal int32) (int32, error) {
	return resolveListLimit(limit, defaultVal, maxVal, "ValidationException")
}

// --- Log group identifier count validator (50 max) ---

// maxLogGroupIdentifiers is the documented list bound of the
// logGroupIdentifiers families that take up to fifty entries — the
// scheduled query's member and DescribeLogGroups' member ("You can
// specify as many as 50 log groups in the array").
const maxLogGroupIdentifiers = 50

// maxDescribeAccountIdentifiers is DescribeLogGroups' accountIdentifiers
// bound: "You can specify as many as 20 account IDs in the array." The
// same value serves ListLogGroups through the store's exported constant.
const maxDescribeAccountIdentifiers = logsstore.MaxListAccountIdentifiers

// validateLogGroupIdentifierCount enforces the AWS-documented limit of 1-50
// log group identifiers per scheduled query.
func validateLogGroupIdentifierCount(ids []string) error {
	if len(ids) > maxLogGroupIdentifiers {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Number of log group identifiers must not exceed %d", maxLogGroupIdentifiers), 400)
	}
	return nil
}

// logGroupIdentifierElementPattern is the Smithy pattern trait on the
// shared LogGroupIdentifier shape.
var logGroupIdentifierElementPattern = regexp.MustCompile(`^[\w#+=/:,.@-]*$`)

// validateLogGroupIdentifierElements enforces the per-element traits of
// the shared LogGroupIdentifier shape (1-2048 characters of the
// identifier alphabet); the count bound lives in
// validateLogGroupIdentifierCount.
func validateLogGroupIdentifierElements(ids []string) error {
	for _, id := range ids {
		if id == "" || len(id) > logsstore.MaxLogGroupIdentifierLength || !logGroupIdentifierElementPattern.MatchString(id) {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("Each log group identifier must be 1 to %d characters of the allowed alphabet", logsstore.MaxLogGroupIdentifierLength), 400)
		}
	}
	return nil
}

// --- ListLogGroups filter-member validators ---

// tagFilterKeyPattern is the Smithy pattern trait on TagFilterKey.
var tagFilterKeyPattern = regexp.MustCompile(`^([\p{L}\p{Z}\p{N}_.:/=+\-@]+)$`)

// tagFilterValuePattern is the Smithy pattern trait on the TagFilterValues
// element: an optional negation prefix and optional leading and trailing
// wildcards over the restricted alphabet.
var tagFilterValuePattern = regexp.MustCompile(`^!?\*?([\p{L}\p{Z}\p{N}_.:/=+\-@]*)\*?$`)

// validateLogGroupTagFilters enforces the logGroupTags member's bounds and
// per-element traits: "Array Members: Minimum number of 1 item. Maximum
// number of 5 items", each key 1-128 over its alphabet, each values list
// 0-5 entries of 0-259 characters.
func validateLogGroupTagFilters(filters []TagFilterInput) error {
	if len(filters) > logsstore.MaxListLogGroupTagFilters {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("logGroupTags accepts as many as %d tag filters; the request carries %d",
				logsstore.MaxListLogGroupTagFilters, len(filters)), 400)
	}
	for _, tf := range filters {
		if tf.Key == "" || utf8.RuneCountInString(tf.Key) > tagutil.MaxTagKeyLength || !tagFilterKeyPattern.MatchString(tf.Key) {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("Tag filter key %q must be 1 to %d characters of the allowed alphabet", tf.Key, tagutil.MaxTagKeyLength), 400)
		}
		if len(tf.Values) > logsstore.MaxTagFilterValues {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("A tag filter accepts as many as %d values; the filter for key %s carries %d",
					logsstore.MaxTagFilterValues, tf.Key, len(tf.Values)), 400)
		}
		for _, v := range tf.Values {
			if utf8.RuneCountInString(v) > logsstore.MaxTagFilterValueLength || !tagFilterValuePattern.MatchString(v) {
				return NewLogsError("InvalidParameterException",
					fmt.Sprintf("Tag filter value %q for key %s must be at most %d characters of the allowed alphabet", v, tf.Key, logsstore.MaxTagFilterValueLength), 400)
			}
		}
	}
	return nil
}

// validateDataSourceFilters enforces the dataSources member's bounds and
// its elements' required name: "Array Members: Minimum number of 1 item.
// Maximum number of 5 items"; the DataSourceFilter name member carries
// the required trait.
func validateDataSourceFilters(filters []DataSourceFilterInput) error {
	if len(filters) > logsstore.MaxListDataSourceFilters {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("dataSources accepts as many as %d filters; the request carries %d",
				logsstore.MaxListDataSourceFilters, len(filters)), 400)
	}
	for _, f := range filters {
		if f.Name == "" {
			return NewLogsError("InvalidParameterException",
				"The name member of every data source filter is required", 400)
		}
	}
	return nil
}

// validateFieldIndexNameFilters enforces the fieldIndexNames member's
// bounds and element traits: "You can specify 1 to 20 field index names,
// each with 1 to 512 characters" over the member's pattern alphabet.
func validateFieldIndexNameFilters(names []string) error {
	if len(names) > logsstore.MaxListFieldIndexNames {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("fieldIndexNames accepts as many as %d names; the request carries %d",
				logsstore.MaxListFieldIndexNames, len(names)), 400)
	}
	for _, name := range names {
		if name == "" || len(name) > logsstore.MaxFieldIndexNameFilterLength || !logGroupNamePattern.MatchString(name) {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("Field index name %q must be 1 to %d characters of the allowed alphabet", name, logsstore.MaxFieldIndexNameFilterLength), 400)
		}
	}
	return nil
}

// validateAccountIdentifierList enforces the accountIdentifiers element
// trait: "Length Constraints: Fixed length of 12" over digits. An empty
// element is the documented null form (all accounts) and stays legal.
func validateAccountIdentifierList(ids []string) error {
	for _, id := range ids {
		if id == "" {
			continue
		}
		if len(id) != 12 || !ownerAccountIDPattern.MatchString(id) {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("Account identifier %q must be exactly 12 digits", id), 400)
		}
	}
	return nil
}

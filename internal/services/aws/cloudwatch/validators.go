package cloudwatch

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	tagutil "vorpalstacks/internal/common/tags"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
	"vorpalstacks/internal/utils/aws/types"
)

// CloudWatch alarm limits per AWS spec.
const (
	maxAlarmTags        = 50
	maxAlarmTagKeyLen   = 128
	maxAlarmTagValueLen = 256
	maxAlarmNameLen     = 255
	maxActionsPerType   = 5
	maxDimensions       = 30
	maxEvaluationWindow = 604800 // 7 days in seconds (Period * EvaluationPeriods)
)

// CloudWatch PutMetricData limits per AWS spec.
const (
	maxMetricDataPerRequest = 1000
	maxValuesPerDatum       = 150
	maxMetricDimensions     = 30
)

// String length limits derived from Smithy traits.
const (
	maxNamespaceLen          = 255
	maxMetricNameLen         = 255
	maxInsightRuleNameLen    = 128
	maxInsightRuleDefLen     = 8192
	maxAlarmDescLen          = 1024
	maxStateReasonLen        = 1023
	minKmsKeyArnLen          = 20
	maxKmsKeyArnLen          = 2048
	maxMetricIdLen           = 255
	maxDashboardNameLen      = 255
	maxDashboardNames        = 100
	maxListMetricsDimensions = 10
	maxDescribeAlarmsRecords = 100
)

// validStatistics is the set of statistic values accepted by
// PutMetricAlarm and GetMetricStatistics (Smithy enum Statistic).
var validStatistics = map[string]bool{
	"SampleCount": true,
	"Average":     true,
	"Sum":         true,
	"Minimum":     true,
	"Maximum":     true,
}

// validTreatMissingData is the set of TreatMissingData values accepted by AWS.
var validTreatMissingData = map[string]bool{
	"breaching":    true,
	"notBreaching": true,
	"ignore":       true,
	"missing":      true,
}

// validComparisonOperators includes all AWS-supported comparison operators,
// including those for anomaly detection alarms.
var validComparisonOperators = map[string]bool{
	"GreaterThanOrEqualToThreshold":            true,
	"GreaterThanThreshold":                     true,
	"LessThanThreshold":                        true,
	"LessThanOrEqualToThreshold":               true,
	"LessThanLowerOrGreaterThanUpperThreshold": true,
	"LessThanLowerThreshold":                   true,
	"GreaterThanUpperThreshold":                true,
}

// validRuleStates is the set of RuleState values accepted by PutInsightRule
// (Smithy InsightRuleState).
var validRuleStates = map[string]bool{
	"ENABLED":  true,
	"DISABLED": true,
}

// validOutputFormats is the set of OutputFormat values accepted by
// GetMetricWidgetImage.
var validOutputFormats = map[string]bool{
	"png":  true,
	"json": true,
}

// validStandardUnits is the set of Unit values accepted by PutMetricAlarm
// and PutMetricData (Smithy enum StandardUnit, 27 values).
var validStandardUnits = map[string]bool{
	"Seconds":          true,
	"Microseconds":     true,
	"Milliseconds":     true,
	"Bytes":            true,
	"Kilobytes":        true,
	"Megabytes":        true,
	"Gigabytes":        true,
	"Terabytes":        true,
	"Bits":             true,
	"Kilobits":         true,
	"Megabits":         true,
	"Gigabits":         true,
	"Terabits":         true,
	"Percent":          true,
	"Count":            true,
	"Bytes/Second":     true,
	"Kilobytes/Second": true,
	"Megabytes/Second": true,
	"Gigabytes/Second": true,
	"Terabytes/Second": true,
	"Bits/Second":      true,
	"Kilobits/Second":  true,
	"Megabits/Second":  true,
	"Gigabits/Second":  true,
	"Terabits/Second":  true,
	"Count/Second":     true,
	"None":             true,
}

// validHistoryItemTypes is the set of HistoryItemType values accepted by
// DescribeAlarmHistory (Smithy enum HistoryItemType).
var validHistoryItemTypes = map[string]bool{
	"Action":                      true,
	"AlarmContributorAction":      true,
	"AlarmContributorStateUpdate": true,
	"ConfigurationUpdate":         true,
	"StateUpdate":                 true,
}

// validStateValues is the set of StateValue values accepted by
// DescribeAlarms and SetAlarmState (Smithy enum StateValue).
var validStateValues = map[string]bool{
	"OK":                true,
	"ALARM":             true,
	"INSUFFICIENT_DATA": true,
}

// validAlarmTypes is the set of AlarmType values accepted by
// DescribeAlarms (Smithy enum AlarmType).
var validAlarmTypes = map[string]bool{
	"CompositeAlarm": true,
	"LogAlarm":       true,
	"MetricAlarm":    true,
}

// Compiled regex patterns for Smithy trait validation.
var (
	insightRuleNamePattern = regexp.MustCompile(`^[\x20-\x7E]+$`)
	insightRuleDefPattern  = regexp.MustCompile(`^[\x00-\x7F]+$`)
	kmsKeyArnPattern       = regexp.MustCompile(`^arn:[a-zA-Z0-9-]+:kms:[a-zA-Z0-9-]+:\d{12}:key/[a-f0-9-]+$`)
	dashboardNamePattern   = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
)

// validateAlarmTagList validates tags against AWS CloudWatch tag limits
// (max 50 tags, key 1-128 chars, value 0-256 chars, no "aws:" prefix).
func validateAlarmTagList(tags []types.Tag) error {
	if len(tags) > maxAlarmTags {
		return awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Number of tags must not exceed %d", maxAlarmTags))
	}
	for _, t := range tags {
		if len(t.Key) == 0 || len(t.Key) > maxAlarmTagKeyLen {
			return awserrors.NewInvalidParameterValueException(
				"Tag key length must be between 1 and 128 characters")
		}
		if len(t.Value) > maxAlarmTagValueLen {
			return awserrors.NewInvalidParameterValueException(
				"Tag value length must not exceed 256 characters")
		}
		if strings.HasPrefix(strings.ToLower(t.Key), "aws:") {
			return awserrors.NewInvalidParameterValueException(
				"Tag keys cannot start with 'aws:'")
		}
	}
	return nil
}

// parseAndValidateAlarmTags parses tags from request parameters and validates
// them against AWS CloudWatch tag limits.
func parseAndValidateAlarmTags(params map[string]interface{}) (map[string]string, error) {
	tagList := tagutil.ParseTags(params, "Tags")
	if len(tagList) == 0 {
		tagList = tagutil.ParseTags(params, "tags")
	}
	if err := validateAlarmTagList(tagList); err != nil {
		return nil, err
	}
	return tagutil.ToMap(tagList), nil
}

// validateAlarmActions validates that each action list does not exceed
// the AWS limit of 5 actions per type.
func validateAlarmActions(alarmActions, okActions, insufficientDataActions []string) error {
	if len(alarmActions) > maxActionsPerType {
		return awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Number of AlarmActions must not exceed %d", maxActionsPerType))
	}
	if len(okActions) > maxActionsPerType {
		return awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Number of OKActions must not exceed %d", maxActionsPerType))
	}
	if len(insufficientDataActions) > maxActionsPerType {
		return awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Number of InsufficientDataActions must not exceed %d", maxActionsPerType))
	}
	return nil
}

// validatePeriod checks that the period is 10, 20, 30, or a multiple of 60.
func validatePeriod(period int32) error {
	if period == 0 {
		return nil
	}
	if period == 10 || period == 20 || period == 30 {
		return nil
	}
	if period >= 60 && period%60 == 0 {
		return nil
	}
	return awserrors.NewInvalidParameterValueException(
		"Period must be 10, 20, 30, or a multiple of 60")
}

// validateExtendedStatistic checks that an extended statistic string matches
// one of the patterns supported by CloudWatch: p{n}, tm{n}, wm{n} (0 ≤ n,
// n < 50 for tm/wm), tc{n}, ts{n} (n ≥ 0), or IQM.
func validateExtendedStatistic(stat string) error {
	if stat == "IQM" {
		return nil
	}
	lower := strings.ToLower(stat)
	switch {
	case strings.HasPrefix(lower, "p"):
		n, err := strconv.ParseFloat(stat[1:], 64)
		if err != nil || n < 0 || n > 100 {
			return awserrors.NewInvalidParameterValueException(
				fmt.Sprintf("Invalid ExtendedStatistic: %s. Percentile must be 0-100", stat))
		}
	case strings.HasPrefix(lower, "tm") || strings.HasPrefix(lower, "wm"):
		n, err := strconv.ParseFloat(stat[2:], 64)
		if err != nil || n < 0 || n >= 50 {
			return awserrors.NewInvalidParameterValueException(
				fmt.Sprintf("Invalid ExtendedStatistic: %s. Value must be 0 to less than 50", stat))
		}
	case strings.HasPrefix(lower, "tc") || strings.HasPrefix(lower, "ts"):
		n, err := strconv.Atoi(stat[2:])
		if err != nil || n < 0 {
			return awserrors.NewInvalidParameterValueException(
				fmt.Sprintf("Invalid ExtendedStatistic: %s. Value must be a non-negative integer", stat))
		}
	default:
		return awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Invalid ExtendedStatistic: %s. Must be p{n}, tm{n}, tc{n}, ts{n}, wm{n}, or IQM", stat))
	}
	return nil
}

// validateRuleState checks that the RuleState is one of the AWS-accepted
// values: ENABLED or DISABLED (Smithy InsightRuleState).
func validateRuleState(state string) error {
	if !validRuleStates[state] {
		return awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Invalid RuleState: %s. Must be ENABLED or DISABLED", state))
	}
	return nil
}

// validateOutputFormat checks that the OutputFormat is one of the
// AWS-accepted values for GetMetricWidgetImage: png or json.
func validateOutputFormat(format string) error {
	if !validOutputFormats[format] {
		return awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Invalid OutputFormat: %s. Must be png or json", format))
	}
	return nil
}

// validateNamespace validates a CloudWatch Namespace
// (Smithy: length 1-255, pattern ^[^:]).
func validateNamespace(ns string) error {
	if len(ns) == 0 || len(ns) > maxNamespaceLen {
		return awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Namespace length must be between 1 and %d characters", maxNamespaceLen))
	}
	if ns[0] == ':' {
		return awserrors.NewInvalidParameterValueException(
			"Namespace must not start with ':'")
	}
	return nil
}

// validateMetricName validates a CloudWatch MetricName
// (Smithy: length 1-255).
func validateMetricName(name string) error {
	if len(name) == 0 || len(name) > maxMetricNameLen {
		return awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("MetricName length must be between 1 and %d characters", maxMetricNameLen))
	}
	return nil
}

// validateInsightRuleName validates an InsightRuleName
// (Smithy: length 1-128, pattern ^[\x20-\x7E]+$).
func validateInsightRuleName(name string) error {
	if len(name) == 0 || len(name) > maxInsightRuleNameLen {
		return awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("RuleName length must be between 1 and %d characters", maxInsightRuleNameLen))
	}
	if !insightRuleNamePattern.MatchString(name) {
		return awserrors.NewInvalidParameterValueException(
			"RuleName must contain only printable ASCII characters (0x20-0x7E)")
	}
	return nil
}

// validateInsightRuleDefinition validates an InsightRuleDefinition
// (Smithy: length 1-8192, pattern ^[\x00-\x7F]+$).
func validateInsightRuleDefinition(def string) error {
	if len(def) == 0 || len(def) > maxInsightRuleDefLen {
		return awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("RuleDefinition length must be between 1 and %d characters", maxInsightRuleDefLen))
	}
	if !insightRuleDefPattern.MatchString(def) {
		return awserrors.NewInvalidParameterValueException(
			"RuleDefinition must contain only ASCII characters (0x00-0x7F)")
	}
	return nil
}

// validateAlarmDescription validates an AlarmDescription
// (Smithy: length 0-1024).
func validateAlarmDescription(desc string) error {
	if len(desc) > maxAlarmDescLen {
		return awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("AlarmDescription must not exceed %d characters", maxAlarmDescLen))
	}
	return nil
}

// validateStateReason validates a StateReason
// (Smithy: length 0-1023).
func validateStateReason(reason string) error {
	if len(reason) > maxStateReasonLen {
		return awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("StateReason must not exceed %d characters", maxStateReasonLen))
	}
	return nil
}

// validateKmsKeyArn validates a KmsKeyArn
// (Smithy: length 20-2048, pattern ^arn:...kms:...key/...).
func validateKmsKeyArn(arn string) error {
	if len(arn) < minKmsKeyArnLen || len(arn) > maxKmsKeyArnLen {
		return awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("KmsKeyArn length must be between %d and %d characters", minKmsKeyArnLen, maxKmsKeyArnLen))
	}
	if !kmsKeyArnPattern.MatchString(arn) {
		return awserrors.NewInvalidParameterValueException(
			"KmsKeyArn must be a valid KMS key ARN (arn:aws:kms:region:accountId:key/keyId)")
	}
	return nil
}

// validateThresholdMetricId validates a ThresholdMetricId
// (Smithy MetricId: length 1-255).
func validateThresholdMetricId(id string) error {
	if id == "" {
		return nil
	}
	if len(id) > maxMetricIdLen {
		return awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("ThresholdMetricId must not exceed %d characters", maxMetricIdLen))
	}
	return nil
}

// validateDashboardName validates a DashboardName
// (AWS docs: max 255, valid characters A-Z a-z 0-9 - _).
func validateDashboardName(name string) error {
	if len(name) == 0 || len(name) > maxDashboardNameLen {
		return awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("DashboardName length must be between 1 and %d characters", maxDashboardNameLen))
	}
	if !dashboardNamePattern.MatchString(name) {
		return awserrors.NewInvalidParameterValueException(
			"DashboardName must contain only alphanumeric characters, hyphens, and underscores")
	}
	return nil
}

// validateUnit validates a Unit value against the Smithy StandardUnit enum.
func validateUnit(unit string) error {
	if unit == "" {
		return nil
	}
	if !validStandardUnits[unit] {
		return awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Invalid Unit: %s", unit))
	}
	return nil
}

// validateHistoryItemType validates a HistoryItemType value against the
// Smithy HistoryItemType enum.
func validateHistoryItemType(t string) error {
	if t == "" {
		return nil
	}
	if !validHistoryItemTypes[t] {
		return awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Invalid HistoryItemType: %s", t))
	}
	return nil
}

// validateStateValueFilter validates a StateValue filter against the
// Smithy StateValue enum.
func validateStateValueFilter(v string) error {
	if v == "" {
		return nil
	}
	if !validStateValues[v] {
		return awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Invalid StateValue: %s. Must be OK, ALARM, or INSUFFICIENT_DATA", v))
	}
	return nil
}

// validateAlarmTypeFilters validates AlarmType filter values against the
// Smithy AlarmType enum.
func validateAlarmTypeFilters(types []string) error {
	for _, t := range types {
		if !validAlarmTypes[t] {
			return awserrors.NewInvalidParameterValueException(
				fmt.Sprintf("Invalid AlarmType: %s. Must be CompositeAlarm, LogAlarm, or MetricAlarm", t))
		}
	}
	return nil
}

// validateDatapointsToAlarm validates a DatapointsToAlarm value
// (Smithy: range min 1).
func validateDatapointsToAlarm(v int32) error {
	if v < 1 {
		return awserrors.NewInvalidParameterValueException(
			"DatapointsToAlarm must be at least 1")
	}
	return nil
}

// validateStorageResolution validates a StorageResolution value
// (Smithy: range min 1).
func validateStorageResolution(v int32) error {
	if v < 1 {
		return awserrors.NewInvalidParameterValueException(
			"StorageResolution must be at least 1")
	}
	return nil
}

// capMaxRecords caps a MaxRecords value to the Smithy-specified maximum.
// Values of 0 are replaced with the default. Values exceeding maxBound are
// capped to maxBound.
func capMaxRecords(v, defaultVal, maxBound int) int {
	if v == 0 {
		return defaultVal
	}
	if v > maxBound {
		return maxBound
	}
	return v
}

// validateMetricDatum checks that a single MetricDatum adheres to AWS
// constraints: at most one of Value, Values+Counts, or StatisticValues
// is provided; Values and Counts must have matching lengths; Values
// must not exceed 150 entries; and Dimensions must not exceed 30.
func validateMetricDatum(datum cwstore.MetricDatum) error {
	if len(datum.Dimensions) > maxMetricDimensions {
		return awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("A MetricDatum can have at most %d Dimensions", maxMetricDimensions))
	}

	hasValue := datum.HasValue
	hasValues := len(datum.Values) > 0
	hasStatSet := datum.StatisticValues != nil

	modeCount := 0
	if hasValue {
		modeCount++
	}
	if hasValues {
		modeCount++
	}
	if hasStatSet {
		modeCount++
	}
	if modeCount > 1 {
		return awserrors.NewInvalidParameterValueException(
			"A MetricDatum must not specify more than one of Value, Values, or StatisticValues")
	}

	if len(datum.Values) > maxValuesPerDatum {
		return awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("A MetricDatum Values array must not exceed %d entries", maxValuesPerDatum))
	}

	if hasValues && len(datum.Counts) > 0 && len(datum.Values) != len(datum.Counts) {
		return awserrors.NewInvalidParameterValueException(
			"Values and Counts arrays must have the same length")
	}

	return nil
}

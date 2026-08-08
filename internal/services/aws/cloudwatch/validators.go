package cloudwatch

import (
	"fmt"
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

// validStatistics is the set of statistic values accepted by PutMetricAlarm.
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

// validRuleStates is the set of RuleState values accepted by PutInsightRule.
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
// values: ENABLED or DISABLED.
func validateRuleState(state string) bool {
	return validRuleStates[state]
}

// validateOutputFormat checks that the OutputFormat is one of the
// AWS-accepted values for GetMetricWidgetImage.
func validateOutputFormat(format string) bool {
	return validOutputFormats[format]
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

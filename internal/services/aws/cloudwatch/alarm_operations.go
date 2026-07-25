package cloudwatch

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
	"vorpalstacks/internal/store/aws/common"
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

func getAlarmStringParam(params map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if v, ok := params[key]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	return ""
}

func getAlarmFloatParam(params map[string]interface{}, keys ...string) float64 {
	for _, key := range keys {
		if v, ok := params[key]; ok {
			switch val := v.(type) {
			case float64:
				return val
			case float32:
				return float64(val)
			case int:
				return float64(val)
			case int64:
				return float64(val)
			}
		}
	}
	return 0
}

func getAlarmIntParam(params map[string]interface{}, keys ...string) int {
	for _, key := range keys {
		if v, ok := params[key]; ok {
			switch val := v.(type) {
			case int:
				return val
			case int32:
				return int(val)
			case int64:
				return int(val)
			case uint64:
				return int(val)
			case float64:
				return int(val)
			case string:
				n, _ := strconv.Atoi(val)
				return n
			}
		}
	}
	return 0
}

func getAlarmBoolParam(params map[string]interface{}, keys []string, defaultVal bool) bool {
	for _, key := range keys {
		if v, ok := params[key]; ok {
			switch val := v.(type) {
			case bool:
				return val
			case string:
				if val == "false" {
					return false
				}
			}
		}
	}
	return defaultVal
}

func getAlarmStringListParam(params map[string]interface{}, keys ...string) []string {
	for _, key := range keys {
		if v, ok := params[key]; ok {
			switch val := v.(type) {
			case []interface{}:
				result := make([]string, 0, len(val))
				for _, item := range val {
					if s, ok := item.(string); ok {
						result = append(result, s)
					}
				}
				return result
			case []string:
				return val
			}
		}
	}
	return nil
}

func parseAlarmTags(params map[string]interface{}) map[string]string {
	tagList := tagutil.ParseTags(params, "Tags")
	if len(tagList) == 0 {
		tagList = tagutil.ParseTags(params, "tags")
	}
	return tagutil.ToMap(tagList)
}

// parseAndValidateAlarmTags parses tags from request parameters and validates
// them against AWS CloudWatch tag limits (max 50 tags, key 1-128 chars,
// value 0-256 chars, no "aws:" prefix on keys).
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

// parseEvaluationCriteria extracts the EvaluationCriteria parameter
// from the request, supporting the PromQLCriteria member.
func parseEvaluationCriteria(params map[string]interface{}) *cwstore.EvaluationCriteria {
	var ecRaw interface{}
	if v, ok := params["EvaluationCriteria"]; ok {
		ecRaw = v
	} else if v, ok := params["evaluationCriteria"]; ok {
		ecRaw = v
	}
	if ecRaw == nil {
		return nil
	}

	ecMap, ok := ecRaw.(map[string]interface{})
	if !ok {
		return nil
	}

	result := &cwstore.EvaluationCriteria{}

	if promql, ok := ecMap["PromQLCriteria"]; ok {
		result.PromQLCriteria = parsePromQLCriteria(promql)
	} else if promql, ok := ecMap["promqlCriteria"]; ok {
		result.PromQLCriteria = parsePromQLCriteria(promql)
	}

	if result.PromQLCriteria == nil {
		return nil
	}
	return result
}

// parsePromQLCriteria extracts PromQL criteria from a raw map.
func parsePromQLCriteria(raw interface{}) *cwstore.AlarmPromQLCriteria {
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil
	}
	criteria := &cwstore.AlarmPromQLCriteria{}
	if q := getAlarmStringParam(m, "Query", "query"); q != "" {
		criteria.Query = q
	}
	if pp := getAlarmIntParam(m, "PendingPeriod", "pendingPeriod"); pp > 0 {
		criteria.PendingPeriod = int32(pp)
	}
	if rp := getAlarmIntParam(m, "RecoveryPeriod", "recoveryPeriod"); rp > 0 {
		criteria.RecoveryPeriod = int32(rp)
	}
	return criteria
}

// evaluationCriteriaToResponse serialises EvaluationCriteria for the
// AWS API response format.
func evaluationCriteriaToResponse(ec *cwstore.EvaluationCriteria) map[string]interface{} {
	if ec == nil {
		return nil
	}
	result := map[string]interface{}{}
	if ec.PromQLCriteria != nil {
		promql := map[string]interface{}{
			"Query": ec.PromQLCriteria.Query,
		}
		if ec.PromQLCriteria.PendingPeriod > 0 {
			promql["PendingPeriod"] = ec.PromQLCriteria.PendingPeriod
		}
		if ec.PromQLCriteria.RecoveryPeriod > 0 {
			promql["RecoveryPeriod"] = ec.PromQLCriteria.RecoveryPeriod
		}
		result["PromQLCriteria"] = promql
	}
	return result
}

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

func alarmToResponse(alarm *cwstore.Alarm) map[string]interface{} {
	stateUpdatedTs := alarm.StateUpdatedTimestamp
	if stateUpdatedTs.IsZero() {
		stateUpdatedTs = alarm.CreatedAt
	}

	result := map[string]interface{}{
		"AlarmName":             alarm.Name,
		"AlarmArn":              alarm.ARN,
		"Namespace":             alarm.Namespace,
		"MetricName":            alarm.MetricName,
		"ComparisonOperator":    alarm.ComparisonOperator,
		"Threshold":             alarm.Threshold,
		"EvaluationPeriods":     alarm.EvaluationPeriods,
		"DatapointsToAlarm":     alarm.DatapointsToAlarm,
		"Period":                alarm.Period,
		"Statistic":             alarm.Statistic,
		"StateValue":            alarm.State,
		"StateReason":           alarm.StateReason,
		"StateUpdatedTimestamp": stateUpdatedTs.UTC().UnixMilli(),
		"TreatMissingData":      alarm.TreatMissingData,
		"ActionsEnabled":        alarm.ActionsEnabled,
	}
	if alarm.StateReasonData != "" {
		result["StateReasonData"] = alarm.StateReasonData
	}
	if alarm.Unit != "" {
		result["Unit"] = string(alarm.Unit)
	}
	if alarm.AlarmType != "" {
		result["AlarmType"] = alarm.AlarmType
	}
	if alarm.AlarmRule != "" {
		result["AlarmRule"] = alarm.AlarmRule
	}
	if alarm.ActionsSuppressedBy != "" {
		result["ActionsSuppressedBy"] = alarm.ActionsSuppressedBy
	}
	if alarm.ActionsSuppressedReason != "" {
		result["ActionsSuppressedReason"] = alarm.ActionsSuppressedReason
	}
	if alarm.ActionsSuppressor != "" {
		result["ActionsSuppressor"] = alarm.ActionsSuppressor
	}
	if alarm.ActionsSuppressorWaitPeriod != 0 {
		result["ActionsSuppressorWaitPeriod"] = alarm.ActionsSuppressorWaitPeriod
	}
	if alarm.ActionsSuppressorExtPeriod != 0 {
		result["ActionsSuppressorExtensionPeriod"] = alarm.ActionsSuppressorExtPeriod
	}
	if alarm.ThresholdMetricID != "" {
		result["ThresholdMetricId"] = alarm.ThresholdMetricID
	}
	if len(alarm.Metrics) > 0 {
		result["Metrics"] = metricDataQueriesToResponse(alarm.Metrics)
	}
	if ec := evaluationCriteriaToResponse(alarm.EvaluationCriteria); ec != nil {
		result["EvaluationCriteria"] = ec
	}
	if alarm.AlarmDescription != "" {
		result["AlarmDescription"] = alarm.AlarmDescription
	}
	if alarm.ExtendedStatistic != "" {
		result["ExtendedStatistic"] = alarm.ExtendedStatistic
	}
	if len(alarm.AlarmActions) > 0 {
		result["AlarmActions"] = alarm.AlarmActions
	}
	if len(alarm.OKActions) > 0 {
		result["OKActions"] = alarm.OKActions
	}
	if len(alarm.InsufficientDataActions) > 0 {
		result["InsufficientDataActions"] = alarm.InsufficientDataActions
	}
	if len(alarm.Dimensions) > 0 {
		dims := make([]map[string]interface{}, len(alarm.Dimensions))
		for j, d := range alarm.Dimensions {
			dims[j] = map[string]interface{}{
				"Name":  d.Name,
				"Value": d.Value,
			}
		}
		result["Dimensions"] = dims
	}
	return result
}

func parseAlarmDimensions(params map[string]interface{}) []cwstore.Dimension {
	var dimensions []cwstore.Dimension
	if dimsRaw, ok := params["Dimensions"]; ok {
		if dimsList, ok := dimsRaw.([]interface{}); ok {
			for _, d := range dimsList {
				if dimMap, ok := d.(map[string]interface{}); ok {
					dimensions = append(dimensions, cwstore.Dimension{
						Name:  getAlarmStringParam(dimMap, "Name", "name"),
						Value: getAlarmStringParam(dimMap, "Value", "value"),
					})
				}
			}
		}
	} else if dimsRaw, ok := params["dimensions"]; ok {
		if dimsList, ok := dimsRaw.([]interface{}); ok {
			for _, d := range dimsList {
				if dimMap, ok := d.(map[string]interface{}); ok {
					dimensions = append(dimensions, cwstore.Dimension{
						Name:  getAlarmStringParam(dimMap, "Name", "name"),
						Value: getAlarmStringParam(dimMap, "Value", "value"),
					})
				}
			}
		}
	}
	return dimensions
}

// PutMetricAlarm creates or updates a metric alarm.
func (s *CloudWatchService) upsertAlarm(store *cwstore.AlarmStore, alarm *cwstore.Alarm, alarmType string) (*cwstore.Alarm, error) {
	existing, err := store.GetAlarm(alarm.Name)
	if err == nil && existing != nil {
		alarm.ARN = existing.ARN
		alarm.CreatedAt = existing.CreatedAt
		alarm.State = existing.State
		alarm.StateUpdatedTimestamp = existing.StateUpdatedTimestamp
		if err := store.UpdateAlarm(alarm); err != nil {
			return nil, err
		}
	} else {
		created, err := store.CreateAlarm(alarm)
		if err != nil {
			return nil, err
		}
		alarm = created
	}

	if len(alarm.Tags) > 0 {
		if err := store.Tag(alarm.ARN, alarm.Tags); err != nil {
			return nil, err
		}
	}

	if err := store.AddAlarmHistory(&cwstore.AlarmHistoryEntry{
		AlarmName:       alarm.Name,
		AlarmType:       alarmType,
		Timestamp:       time.Now().UTC().UnixMilli(),
		HistoryItemType: cwstore.HistoryItemTypeConfigurationUpdate,
		HistorySummary:  "Alarm was created or updated",
	}); err != nil {
		logs.Warn("failed to add alarm history", logs.String("alarm", alarm.Name), logs.Err(err))
	}

	return alarm, nil
}

// PutMetricAlarm creates or updates a CloudWatch alarm.
func (s *CloudWatchService) PutMetricAlarm(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	alarmName := getAlarmStringParam(req.Parameters, "AlarmName", "alarmName")
	if alarmName == "" {
		return nil, ErrInvalidParameter
	}
	if len(alarmName) > maxAlarmNameLen {
		return nil, awserrors.NewInvalidParameterValueException(
			"AlarmName must not exceed 255 characters")
	}

	namespace := getAlarmStringParam(req.Parameters, "Namespace", "namespace")
	metricName := getAlarmStringParam(req.Parameters, "MetricName", "metricName")

	alarm := cwstore.NewAlarm(alarmName, namespace, metricName)
	alarm.Dimensions = parseAlarmDimensions(req.Parameters)
	alarm.ComparisonOperator = getAlarmStringParam(req.Parameters, "ComparisonOperator", "comparisonOperator")
	alarm.Threshold = getAlarmFloatParam(req.Parameters, "Threshold", "threshold")
	alarm.EvaluationPeriods = int32(getAlarmIntParam(req.Parameters, "EvaluationPeriods", "evaluationPeriods"))
	alarm.Period = int32(getAlarmIntParam(req.Parameters, "Period", "period"))
	alarm.Statistic = getAlarmStringParam(req.Parameters, "Statistic", "statistic")
	alarm.ExtendedStatistic = getAlarmStringParam(req.Parameters, "ExtendedStatistic", "extendedStatistic")
	alarm.TreatMissingData = getAlarmStringParam(req.Parameters, "TreatMissingData", "treatMissingData")
	alarm.AlarmDescription = getAlarmStringParam(req.Parameters, "AlarmDescription", "alarmDescription")
	alarm.ActionsEnabled = getAlarmBoolParam(req.Parameters, []string{"ActionsEnabled", "actionsEnabled"}, true)
	alarm.DatapointsToAlarm = int32(getAlarmIntParam(req.Parameters, "DatapointsToAlarm", "datapointsToAlarm"))
	if alarm.DatapointsToAlarm == 0 {
		alarm.DatapointsToAlarm = alarm.EvaluationPeriods
	}
	alarm.AlarmActions = getAlarmStringListParam(req.Parameters, "AlarmActions", "alarmActions")
	alarm.OKActions = getAlarmStringListParam(req.Parameters, "OKActions", "okActions")
	alarm.InsufficientDataActions = getAlarmStringListParam(req.Parameters, "InsufficientDataActions", "insufficientDataActions")
	alarm.Unit = cwstore.StandardUnit(getAlarmStringParam(req.Parameters, "Unit", "unit"))
	alarm.ThresholdMetricID = getAlarmStringParam(req.Parameters, "ThresholdMetricId", "thresholdMetricId")
	if metricsRaw, ok := req.Parameters["Metrics"]; ok {
		alarm.Metrics = parseMetricDataQueries(metricsRaw)
	} else if metricsRaw, ok := req.Parameters["metrics"]; ok {
		alarm.Metrics = parseMetricDataQueries(metricsRaw)
	}

	if alarm.ComparisonOperator == "" {
		alarm.ComparisonOperator = "GreaterThanOrEqualToThreshold"
	}
	if !validComparisonOperators[alarm.ComparisonOperator] {
		return nil, awserrors.NewInvalidParameterValueException(fmt.Sprintf("Invalid ComparisonOperator: %s", alarm.ComparisonOperator))
	}
	if alarm.EvaluationPeriods == 0 {
		alarm.EvaluationPeriods = 1
	}
	if alarm.Period == 0 {
		alarm.Period = 60
	}
	if err := validatePeriod(alarm.Period); err != nil {
		return nil, err
	}
	if int64(alarm.Period)*int64(alarm.EvaluationPeriods) > maxEvaluationWindow {
		return nil, awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Period * EvaluationPeriods must not exceed %d seconds (7 days)", maxEvaluationWindow))
	}
	if alarm.DatapointsToAlarm > alarm.EvaluationPeriods {
		return nil, awserrors.NewInvalidParameterValueException(
			"DatapointsToAlarm must be less than or equal to EvaluationPeriods")
	}
	if alarm.Statistic == "" && alarm.ExtendedStatistic == "" && len(alarm.Metrics) == 0 {
		alarm.Statistic = "Average"
	}
	if alarm.Statistic != "" && !validStatistics[alarm.Statistic] {
		return nil, awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Invalid Statistic: %s. Must be one of SampleCount, Average, Sum, Minimum, Maximum", alarm.Statistic))
	}
	if alarm.Statistic != "" && alarm.ExtendedStatistic != "" {
		return nil, awserrors.NewInvalidParameterValueException(
			"Statistic and ExtendedStatistic are mutually exclusive")
	}
	if alarm.ExtendedStatistic != "" {
		if err := validateExtendedStatistic(alarm.ExtendedStatistic); err != nil {
			return nil, err
		}
	}
	if alarm.TreatMissingData == "" {
		alarm.TreatMissingData = "missing"
	}
	if !validTreatMissingData[alarm.TreatMissingData] {
		return nil, awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Invalid TreatMissingData: %s. Must be one of breaching, notBreaching, ignore, missing", alarm.TreatMissingData))
	}
	if len(alarm.Dimensions) > maxDimensions {
		return nil, awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Number of Dimensions must not exceed %d", maxDimensions))
	}
	if err := validateAlarmActions(alarm.AlarmActions, alarm.OKActions, alarm.InsufficientDataActions); err != nil {
		return nil, err
	}

	// Parse EvaluationCriteria (PromQL support).
	alarm.EvaluationCriteria = parseEvaluationCriteria(req.Parameters)

	var tagErr error
	alarm.Tags, tagErr = parseAndValidateAlarmTags(req.Parameters)
	if tagErr != nil {
		return nil, tagErr
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	alarm, err = s.upsertAlarm(store.alarms, alarm, cwstore.AlarmTypeMetricAlarm)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"AlarmArn": alarm.ARN,
	}, nil
}

// DescribeAlarms returns a list of alarms.
func (s *CloudWatchService) DescribeAlarms(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	alarmNamePrefix := getAlarmStringParam(req.Parameters, "AlarmNamePrefix", "alarmNamePrefix")
	childrenOfAlarmName := getAlarmStringParam(req.Parameters, "ChildrenOfAlarmName", "childrenOfAlarmName")
	parentsOfAlarmName := getAlarmStringParam(req.Parameters, "ParentsOfAlarmName", "parentsOfAlarmName")
	stateValueFilter := getAlarmStringParam(req.Parameters, "StateValue", "stateValue")
	actionPrefix := getAlarmStringParam(req.Parameters, "ActionPrefix", "actionPrefix")

	var alarmNames []string
	if namesRaw, ok := req.Parameters["AlarmNames"]; ok {
		if namesList, ok := namesRaw.([]interface{}); ok {
			for _, n := range namesList {
				if ns, ok := n.(string); ok {
					alarmNames = append(alarmNames, ns)
				}
			}
		}
	} else if namesRaw, ok := req.Parameters["alarmNames"]; ok {
		if namesList, ok := namesRaw.([]interface{}); ok {
			for _, n := range namesList {
				if ns, ok := n.(string); ok {
					alarmNames = append(alarmNames, ns)
				}
			}
		}
	}

	var alarmTypes []string
	if typesRaw, ok := req.Parameters["AlarmTypes"]; ok {
		if typesList, ok := typesRaw.([]interface{}); ok {
			for _, t := range typesList {
				if ts, ok := t.(string); ok {
					alarmTypes = append(alarmTypes, ts)
				}
			}
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Build the combined filter function for StateValue, ActionPrefix
	// and AlarmTypes. This must be constructed before the
	// ChildrenOfAlarmName and ParentsOfAlarmName branches so that all
	// filters are consistently applied.
	var typeFilter func(*cwstore.Alarm) bool
	if len(alarmTypes) > 0 {
		typeSet := make(map[string]bool)
		for _, t := range alarmTypes {
			typeSet[t] = true
		}
		typeFilter = func(a *cwstore.Alarm) bool {
			at := a.AlarmType
			if at == "" {
				at = cwstore.AlarmTypeMetricAlarm
			}
			return typeSet[at]
		}
	}
	if stateValueFilter != "" {
		svf := typeFilter
		typeFilter = func(a *cwstore.Alarm) bool {
			if svf != nil && !svf(a) {
				return false
			}
			return a.State == stateValueFilter
		}
	}
	if actionPrefix != "" {
		apf := typeFilter
		typeFilter = func(a *cwstore.Alarm) bool {
			if apf != nil && !apf(a) {
				return false
			}
			return alarmMatchesActionPrefix(a, actionPrefix)
		}
	}

	// ChildrenOfAlarmName: resolve the composite alarm's rule to find
	// child alarm names, then filter to only those children.
	if childrenOfAlarmName != "" {
		parent, err := store.alarms.GetAlarm(childrenOfAlarmName)
		if err != nil {
			return s.buildDescribeAlarmsResponse(nil, alarmTypes, ""), nil
		}
		rule, err := parseAlarmRule(parent.AlarmRule)
		if err != nil {
			return s.buildDescribeAlarmsResponse(nil, alarmTypes, ""), nil
		}
		childNames := rule.childAlarmNames()
		var children []*cwstore.Alarm
		for _, name := range childNames {
			if a, err := store.alarms.GetAlarm(name); err == nil {
				if typeFilter == nil || typeFilter(a) {
					children = append(children, a)
				}
			}
		}
		return s.buildDescribeAlarmsResponse(children, alarmTypes, ""), nil
	}

	// ParentsOfAlarmName: find all composite alarms whose AlarmRule
	// references the specified alarm name.
	if parentsOfAlarmName != "" {
		allAlarms, err := store.alarms.ListAlarms("")
		if err != nil {
			return nil, err
		}
		var parents []*cwstore.Alarm
		for _, a := range allAlarms {
			if a.AlarmType != cwstore.AlarmTypeCompositeAlarm || a.AlarmRule == "" {
				continue
			}
			rule, err := parseAlarmRule(a.AlarmRule)
			if err != nil {
				continue
			}
			for _, childName := range rule.childAlarmNames() {
				if childName == parentsOfAlarmName {
					if typeFilter == nil || typeFilter(a) {
						parents = append(parents, a)
					}
					break
				}
			}
		}
		return s.buildDescribeAlarmsResponse(parents, alarmTypes, ""), nil
	}

	var alarms []*cwstore.Alarm
	if len(alarmNames) > 0 {
		for _, name := range alarmNames {
			alarm, err := store.alarms.GetAlarm(name)
			if err == nil {
				if typeFilter == nil || typeFilter(alarm) {
					alarms = append(alarms, alarm)
				}
			}
		}
		return s.buildDescribeAlarmsResponse(alarms, alarmTypes, ""), nil
	}

	marker := pagination.GetMarker(req.Parameters, "NextToken")
	maxRecords := pagination.GetMaxItems(req.Parameters, 100, "MaxRecords")

	opts := common.ListOptions{Marker: marker, MaxItems: maxRecords}
	result, err := store.alarms.ListAlarmsPaginated(alarmNamePrefix, opts, typeFilter)
	if err != nil {
		return nil, err
	}

	return s.buildDescribeAlarmsResponse(result.Items, alarmTypes, result.NextMarker), nil
}

// scheduledQueryConfigToResponse serialises a ScheduledQueryConfig to
// the Smithy ScheduledQueryConfiguration response format.
func scheduledQueryConfigToResponse(sqc *cwstore.ScheduledQueryConfig) map[string]interface{} {
	result := map[string]interface{}{}
	if sqc.QueryString != "" {
		result["QueryString"] = sqc.QueryString
	}
	if len(sqc.LogGroupIdentifiers) > 0 {
		result["LogGroupIdentifiers"] = sqc.LogGroupIdentifiers
	}
	if sqc.QueryARN != "" {
		result["QueryARN"] = sqc.QueryARN
	}
	if sqc.ScheduledQueryRoleARN != "" {
		result["ScheduledQueryRoleARN"] = sqc.ScheduledQueryRoleARN
	}
	if sqc.ScheduleConfiguration != nil {
		sc := map[string]interface{}{}
		if sqc.ScheduleConfiguration.ScheduleExpression != "" {
			sc["ScheduleExpression"] = sqc.ScheduleConfiguration.ScheduleExpression
		}
		if sqc.ScheduleConfiguration.StartTimeOffset != 0 {
			sc["StartTimeOffset"] = sqc.ScheduleConfiguration.StartTimeOffset
		}
		if sqc.ScheduleConfiguration.EndTimeOffset != 0 {
			sc["EndTimeOffset"] = sqc.ScheduleConfiguration.EndTimeOffset
		}
		result["ScheduleConfiguration"] = sc
	}
	if sqc.AggregationExpression != "" {
		result["AggregationExpression"] = sqc.AggregationExpression
	}
	if len(sqc.Tags) > 0 {
		tags := make([]map[string]interface{}, 0, len(sqc.Tags))
		for k, v := range sqc.Tags {
			tags = append(tags, map[string]interface{}{
				"Key":   k,
				"Value": v,
			})
		}
		result["Tags"] = tags
	}
	return result
}

// logAlarmToResponse builds a response map for a log alarm, including
// only the members defined in the Smithy LogAlarm shape.
func logAlarmToResponse(alarm *cwstore.Alarm) map[string]interface{} {
	stateUpdatedTs := alarm.StateUpdatedTimestamp
	if stateUpdatedTs.IsZero() {
		stateUpdatedTs = alarm.CreatedAt
	}
	result := map[string]interface{}{
		"AlarmName":             alarm.Name,
		"AlarmArn":              alarm.ARN,
		"ActionsEnabled":        alarm.ActionsEnabled,
		"StateValue":            alarm.State,
		"StateReason":           alarm.StateReason,
		"StateUpdatedTimestamp": stateUpdatedTs.UTC().UnixMilli(),
		"Threshold":             alarm.Threshold,
		"ComparisonOperator":    alarm.ComparisonOperator,
		"TreatMissingData":      alarm.TreatMissingData,
	}
	if alarm.AlarmDescription != "" {
		result["AlarmDescription"] = alarm.AlarmDescription
	}
	if alarm.StateReasonData != "" {
		result["StateReasonData"] = alarm.StateReasonData
	}
	if len(alarm.AlarmActions) > 0 {
		result["AlarmActions"] = alarm.AlarmActions
	}
	if len(alarm.OKActions) > 0 {
		result["OKActions"] = alarm.OKActions
	}
	if len(alarm.InsufficientDataActions) > 0 {
		result["InsufficientDataActions"] = alarm.InsufficientDataActions
	}
	// ScheduledQueryConfiguration (log alarm).
	if alarm.ScheduledQueryConfiguration != nil {
		result["ScheduledQueryConfiguration"] = scheduledQueryConfigToResponse(alarm.ScheduledQueryConfiguration)
	}
	if alarm.QueryResultsToEvaluate > 0 {
		result["QueryResultsToEvaluate"] = alarm.QueryResultsToEvaluate
	}
	if alarm.QueryResultsToAlarm > 0 {
		result["QueryResultsToAlarm"] = alarm.QueryResultsToAlarm
	}
	if alarm.ActionLogLineCount > 0 {
		result["ActionLogLineCount"] = alarm.ActionLogLineCount
	}
	if alarm.ActionLogLineRoleArn != "" {
		result["ActionLogLineRoleArn"] = alarm.ActionLogLineRoleArn
	}
	return result
}

func (s *CloudWatchService) buildDescribeAlarmsResponse(alarms []*cwstore.Alarm, alarmTypes []string, nextToken string) map[string]interface{} {
	metricAlarms := make([]map[string]interface{}, 0)
	compositeAlarms := make([]map[string]interface{}, 0)
	logAlarms := make([]map[string]interface{}, 0)
	for _, alarm := range alarms {
		switch alarm.AlarmType {
		case cwstore.AlarmTypeCompositeAlarm:
			compositeAlarms = append(compositeAlarms, alarmToResponse(alarm))
		case cwstore.AlarmTypeLogAlarm:
			logAlarms = append(logAlarms, logAlarmToResponse(alarm))
		default:
			metricAlarms = append(metricAlarms, alarmToResponse(alarm))
		}
	}

	response := map[string]interface{}{
		"MetricAlarms":    metricAlarms,
		"CompositeAlarms": compositeAlarms,
		"LogAlarms":       logAlarms,
	}
	if nextToken != "" {
		response["NextToken"] = nextToken
	}
	return response
}

// DescribeAlarmsForMetric returns alarms for a specific metric.
func (s *CloudWatchService) DescribeAlarmsForMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	namespace := getAlarmStringParam(req.Parameters, "Namespace", "namespace")
	metricName := getAlarmStringParam(req.Parameters, "MetricName", "metricName")

	if namespace == "" || metricName == "" {
		return nil, ErrInvalidParameter
	}

	dimensions := parseAlarmDimensions(req.Parameters)
	statistic := getAlarmStringParam(req.Parameters, "Statistic", "statistic")
	extendedStatistic := getAlarmStringParam(req.Parameters, "ExtendedStatistic", "extendedStatistic")
	period := int32(getAlarmIntParam(req.Parameters, "Period", "period"))
	unit := getAlarmStringParam(req.Parameters, "Unit", "unit")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	marker := pagination.GetMarker(req.Parameters, "NextToken")
	maxRecords := pagination.GetMaxItems(req.Parameters, 100, "MaxRecords")

	filter := func(a *cwstore.Alarm) bool {
		if a.Namespace != namespace || a.MetricName != metricName {
			return false
		}
		if len(dimensions) > 0 && len(a.Dimensions) > 0 {
			if !dimensionsMatch(a.Dimensions, dimensions) {
				return false
			}
		}
		if statistic != "" && a.Statistic != statistic {
			return false
		}
		if extendedStatistic != "" && a.ExtendedStatistic != extendedStatistic {
			return false
		}
		if period > 0 && a.Period != period {
			return false
		}
		if unit != "" && string(a.Unit) != unit {
			return false
		}
		return true
	}

	opts := common.ListOptions{Marker: marker, MaxItems: maxRecords}
	result, err := store.alarms.ListAlarmsPaginated("", opts, filter)
	if err != nil {
		return nil, err
	}

	metricAlarms := make([]map[string]interface{}, len(result.Items))
	for i, alarm := range result.Items {
		metricAlarms[i] = alarmToResponse(alarm)
	}

	response := map[string]interface{}{
		"MetricAlarms": metricAlarms,
	}
	if result.NextMarker != "" {
		response["NextToken"] = result.NextMarker
	}
	return response, nil
}

// alarmMatchesActionPrefix checks whether any of the alarm's action
// lists contain an action whose ARN starts with the given prefix.
func alarmMatchesActionPrefix(a *cwstore.Alarm, prefix string) bool {
	check := func(actions []string) bool {
		for _, act := range actions {
			if strings.HasPrefix(act, prefix) {
				return true
			}
		}
		return false
	}
	return check(a.AlarmActions) || check(a.OKActions) || check(a.InsufficientDataActions)
}

func dimensionsMatch(a, b []cwstore.Dimension) bool {
	if len(a) != len(b) {
		return false
	}
	aMap := make(map[string]string)
	for _, d := range a {
		aMap[d.Name] = d.Value
	}
	for _, d := range b {
		if aMap[d.Name] != d.Value {
			return false
		}
	}
	return true
}

// DeleteAlarms deletes one or more alarms.
func (s *CloudWatchService) DeleteAlarms(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var alarmNames []string
	if namesRaw, ok := req.Parameters["AlarmNames"]; ok {
		if namesList, ok := namesRaw.([]interface{}); ok {
			for _, n := range namesList {
				if ns, ok := n.(string); ok {
					alarmNames = append(alarmNames, ns)
				}
			}
		}
	} else if namesRaw, ok := req.Parameters["alarmNames"]; ok {
		if namesList, ok := namesRaw.([]interface{}); ok {
			for _, n := range namesList {
				if ns, ok := n.(string); ok {
					alarmNames = append(alarmNames, ns)
				}
			}
		}
	}

	if len(alarmNames) == 0 {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Check for ResourceConflict: any composite alarm referencing the
	// alarms being deleted would break its AlarmRule.
	allAlarms, err := store.alarms.ListAlarms("")
	if err != nil {
		return nil, err
	}
	deleteSet := make(map[string]bool, len(alarmNames))
	for _, n := range alarmNames {
		deleteSet[n] = true
	}
	for _, a := range allAlarms {
		if a.AlarmType != cwstore.AlarmTypeCompositeAlarm || a.AlarmRule == "" {
			continue
		}
		rule, err := parseAlarmRule(a.AlarmRule)
		if err != nil {
			continue
		}
		for _, childName := range rule.childAlarmNames() {
			if deleteSet[childName] {
				return nil, awserrors.NewAWSError("ResourceConflict",
					fmt.Sprintf("Alarm %s is referenced by composite alarm %s", childName, a.Name),
					http.StatusConflict)
			}
		}
	}

	for _, name := range alarmNames {
		if err := store.alarms.DeleteAlarm(name); err != nil {
			if errors.Is(err, cwstore.ErrAlarmNotFound) {
				return nil, ErrAlarmNotFound
			}
			return nil, err
		}
	}

	return response.EmptyResponse(), nil
}

// SetAlarmState sets the state of an alarm.
func (s *CloudWatchService) SetAlarmState(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	alarmName := getAlarmStringParam(req.Parameters, "AlarmName", "alarmName")
	stateValue := getAlarmStringParam(req.Parameters, "StateValue", "stateValue")
	stateReason := getAlarmStringParam(req.Parameters, "StateReason", "stateReason")
	stateReasonData := getAlarmStringParam(req.Parameters, "StateReasonData", "stateReasonData")

	if alarmName == "" || stateValue == "" {
		return nil, ErrInvalidParameter
	}

	validStates := map[string]bool{
		"OK":                true,
		"ALARM":             true,
		"INSUFFICIENT_DATA": true,
	}
	if !validStates[stateValue] {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.alarms.SetAlarmState(alarmName, stateValue, stateReason, stateReasonData); err != nil {
		if errors.Is(err, cwstore.ErrAlarmNotFound) {
			return nil, ErrAlarmNotFound
		}
		return nil, err
	}

	alarm, _ := store.alarms.GetAlarm(alarmName)
	if err := store.alarms.AddAlarmHistory(&cwstore.AlarmHistoryEntry{
		AlarmName:       alarmName,
		AlarmType:       resolveAlarmType(alarm),
		Timestamp:       time.Now().UTC().UnixMilli(),
		HistoryItemType: cwstore.HistoryItemTypeStateUpdate,
		HistorySummary:  stateReason,
	}); err != nil {
		logs.Warn("failed to add alarm state history", logs.String("alarm", alarmName), logs.Err(err))
	}

	return response.EmptyResponse(), nil
}

func resolveAlarmType(alarm *cwstore.Alarm) string {
	if alarm != nil && alarm.AlarmType == cwstore.AlarmTypeCompositeAlarm {
		return cwstore.AlarmTypeCompositeAlarm
	}
	return cwstore.AlarmTypeMetricAlarm
}

func addAlarmActionHistory(store *cwstore.AlarmStore, alarmName, summary string) {
	alarm, _ := store.GetAlarm(alarmName)
	alarmType := resolveAlarmType(alarm)
	if err := store.AddAlarmHistory(&cwstore.AlarmHistoryEntry{
		AlarmName:       alarmName,
		AlarmType:       alarmType,
		Timestamp:       time.Now().UTC().UnixMilli(),
		HistoryItemType: cwstore.HistoryItemTypeAction,
		HistorySummary:  summary,
	}); err != nil {
		logs.Warn("failed to add alarm history", logs.String("alarm", alarmName), logs.Err(err))
	}
}

func cwAlarmTagConfig(s *CloudWatchService, reqCtx *request.RequestContext) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.StandardARNConfig,
		TagFunc: func(_ context.Context, resourceKey string, tags []types.Tag) error {
			store, err := s.store(reqCtx)
			if err != nil {
				return err
			}
			return store.alarms.TagFromSlice(resourceKey, tags)
		},
		UntagFunc: func(ctx context.Context, resourceKey string, tagKeys []string) error {
			store, err := s.store(reqCtx)
			if err != nil {
				return err
			}
			return store.alarms.Untag(resourceKey, tagKeys)
		},
		ListFunc: func(ctx context.Context, resourceKey string) ([]types.Tag, error) {
			store, err := s.store(reqCtx)
			if err != nil {
				return nil, err
			}
			return store.alarms.ListAsSlice(resourceKey)
		},
		EmptyResponse: func() (interface{}, error) { return response.EmptyResponse(), nil },
	}
}

// TagResource adds tags to a CloudWatch alarm.
func (s *CloudWatchService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagutil.HandleTag(ctx, req, cwAlarmTagConfig(s, reqCtx))
}

// UntagResource removes tags from a CloudWatch alarm.
func (s *CloudWatchService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagutil.HandleUntag(ctx, req, cwAlarmTagConfig(s, reqCtx))
}

// ListTagsForResource lists tags for a CloudWatch alarm.
func (s *CloudWatchService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagutil.HandleList(ctx, req, cwAlarmTagConfig(s, reqCtx))
}

// EnableAlarmActions enables actions for the specified alarms.
func (s *CloudWatchService) EnableAlarmActions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	alarmNames := getAlarmStringListParam(req.Parameters, "AlarmNames", "alarmNames")
	if len(alarmNames) == 0 {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	for _, name := range alarmNames {
		if err := store.alarms.SetAlarmActionsEnabled(name, true); err != nil {
			if errors.Is(err, cwstore.ErrAlarmNotFound) {
				return nil, ErrAlarmNotFound
			}
			return nil, err
		}
		addAlarmActionHistory(store.alarms, name, "Alarm actions enabled")
	}

	return response.EmptyResponse(), nil
}

// DisableAlarmActions disables actions for the specified alarms.
func (s *CloudWatchService) DisableAlarmActions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	alarmNames := getAlarmStringListParam(req.Parameters, "AlarmNames", "alarmNames")
	if len(alarmNames) == 0 {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	for _, name := range alarmNames {
		if err := store.alarms.SetAlarmActionsEnabled(name, false); err != nil {
			if errors.Is(err, cwstore.ErrAlarmNotFound) {
				return nil, ErrAlarmNotFound
			}
			return nil, err
		}
		addAlarmActionHistory(store.alarms, name, "Alarm actions disabled")
	}

	return response.EmptyResponse(), nil
}

// DescribeAlarmHistory retrieves the history for the specified alarm.
func (s *CloudWatchService) DescribeAlarmHistory(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	alarmName := getAlarmStringParam(req.Parameters, "AlarmName", "alarmName")
	historyItemType := getAlarmStringParam(req.Parameters, "HistoryItemType", "historyItemType")
	startDate := parseTimestampFromMap(req.Parameters, "StartDate")
	endDate := parseTimestampFromMap(req.Parameters, "EndDate")
	scanBy := getAlarmStringParam(req.Parameters, "ScanBy", "scanBy")

	var alarmTypeFilter map[string]bool
	if typesRaw, ok := req.Parameters["AlarmTypes"]; ok {
		if typesList, ok := typesRaw.([]interface{}); ok {
			alarmTypeFilter = make(map[string]bool)
			for _, t := range typesList {
				if ts, ok := t.(string); ok {
					alarmTypeFilter[ts] = true
				}
			}
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	marker := pagination.GetMarker(req.Parameters, "NextToken")
	maxRecords := pagination.GetMaxItems(req.Parameters, 100, "MaxRecords")

	var startMs, endMs int64
	if !startDate.IsZero() {
		startMs = startDate.UnixMilli()
	}
	if !endDate.IsZero() {
		endMs = endDate.UnixMilli()
	}

	opts := cwstore.ListAlarmHistoryOpts{
		AlarmName:       alarmName,
		HistoryItemType: historyItemType,
		AlarmTypes:      alarmTypeFilter,
		StartDate:       startMs,
		EndDate:         endMs,
		ScanBy:          scanBy,
		ListOpts:        common.ListOptions{Marker: marker, MaxItems: maxRecords},
	}
	result, err := store.alarms.ListAlarmHistoryPaginated(opts)
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(result.Items))
	for _, entry := range result.Items {
		item := map[string]interface{}{
			"AlarmName":       entry.AlarmName,
			"AlarmType":       entry.AlarmType,
			"Timestamp":       entry.Timestamp,
			"HistoryItemType": entry.HistoryItemType,
			"HistorySummary":  entry.HistorySummary,
		}
		if entry.HistoryData != "" {
			item["HistoryData"] = entry.HistoryData
		}
		items = append(items, item)
	}

	response := map[string]interface{}{
		"AlarmHistoryItems": items,
	}
	if result.NextMarker != "" {
		response["NextToken"] = result.NextMarker
	}
	return response, nil
}

// PutCompositeAlarm creates or updates a composite alarm.
func (s *CloudWatchService) PutCompositeAlarm(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	alarmName := getAlarmStringParam(req.Parameters, "AlarmName", "alarmName")
	alarmRule := getAlarmStringParam(req.Parameters, "AlarmRule", "alarmRule")

	if alarmName == "" || alarmRule == "" {
		return nil, ErrInvalidParameter
	}
	if len(alarmName) > maxAlarmNameLen {
		return nil, awserrors.NewInvalidParameterValueException(
			"AlarmName must not exceed 255 characters")
	}

	alarm := cwstore.NewAlarm(alarmName, "", "")
	alarm.AlarmRule = alarmRule
	alarm.AlarmType = cwstore.AlarmTypeCompositeAlarm
	alarm.AlarmDescription = getAlarmStringParam(req.Parameters, "AlarmDescription", "alarmDescription")
	alarm.ActionsEnabled = getAlarmBoolParam(req.Parameters, []string{"ActionsEnabled", "actionsEnabled"}, true)
	alarm.AlarmActions = getAlarmStringListParam(req.Parameters, "AlarmActions", "alarmActions")
	alarm.OKActions = getAlarmStringListParam(req.Parameters, "OKActions", "okActions")
	alarm.InsufficientDataActions = getAlarmStringListParam(req.Parameters, "InsufficientDataActions", "insufficientDataActions")
	alarm.ActionsSuppressor = getAlarmStringParam(req.Parameters, "ActionsSuppressor", "actionsSuppressor")
	alarm.ActionsSuppressorWaitPeriod = int32(getAlarmIntParam(req.Parameters, "ActionsSuppressorWaitPeriod", "actionsSuppressorWaitPeriod"))
	alarm.ActionsSuppressorExtPeriod = int32(getAlarmIntParam(req.Parameters, "ActionsSuppressorExtensionPeriod", "actionsSuppressorExtPeriod"))

	// Validate the AlarmRule expression at creation time so that syntax
	// errors are reported immediately rather than at evaluation time.
	if _, err := parseAlarmRule(alarmRule); err != nil {
		return nil, awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Invalid AlarmRule: %v", err))
	}

	if err := validateAlarmActions(alarm.AlarmActions, alarm.OKActions, alarm.InsufficientDataActions); err != nil {
		return nil, err
	}

	var tagErr error
	alarm.Tags, tagErr = parseAndValidateAlarmTags(req.Parameters)
	if tagErr != nil {
		return nil, tagErr
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	alarm, err = s.upsertAlarm(store.alarms, alarm, cwstore.AlarmTypeCompositeAlarm)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"AlarmArn": alarm.ARN,
	}, nil
}

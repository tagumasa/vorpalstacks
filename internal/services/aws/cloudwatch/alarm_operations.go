package cloudwatch

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
	"vorpalstacks/internal/utils/aws/types"
)

func getAlarmStringParam(params map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if v, ok := params[key]; ok {
			if s, ok := v.(string); ok {
				return s
			}
			logs.Warn("parameter type mismatch: expected string", logs.String("key", key), logs.String("actualType", fmt.Sprintf("%T", v)))
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
			default:
				logs.Warn("parameter type mismatch: expected numeric", logs.String("key", key), logs.String("actualType", fmt.Sprintf("%T", v)))
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
				n, err := strconv.Atoi(val)
				if err != nil {
					logs.Warn("parameter parse error: invalid integer string", logs.String("key", key), logs.String("value", val))
				}
				return n
			default:
				logs.Warn("parameter type mismatch: expected numeric", logs.String("key", key), logs.String("actualType", fmt.Sprintf("%T", v)))
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
				if val == "true" {
					return true
				}
				logs.Warn("parameter parse error: invalid bool string", logs.String("key", key), logs.String("value", val))
			default:
				logs.Warn("parameter type mismatch: expected bool", logs.String("key", key), logs.String("actualType", fmt.Sprintf("%T", v)))
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
					} else {
						logs.Warn("list item type mismatch: expected string", logs.String("key", key), logs.String("actualType", fmt.Sprintf("%T", item)))
					}
				}
				return result
			case []string:
				return val
			default:
				logs.Warn("parameter type mismatch: expected list", logs.String("key", key), logs.String("actualType", fmt.Sprintf("%T", v)))
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
		// AWS: PutMetricAlarm Tags are only applied on creation, not update.
		if err := store.UpdateAlarm(alarm); err != nil {
			return nil, err
		}
	} else {
		created, err := store.CreateAlarm(alarm)
		if err != nil {
			return nil, err
		}
		alarm = created

		if len(alarm.Tags) > 0 {
			if err := store.Tag(alarm.ARN, alarm.Tags); err != nil {
				return nil, err
			}
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
	input := &PutMetricAlarmInput{
		AlarmName:               getAlarmStringParam(req.Parameters, "AlarmName", "alarmName"),
		Namespace:               getAlarmStringParam(req.Parameters, "Namespace", "namespace"),
		MetricName:              getAlarmStringParam(req.Parameters, "MetricName", "metricName"),
		Dimensions:              parseAlarmDimensions(req.Parameters),
		ComparisonOperator:      getAlarmStringParam(req.Parameters, "ComparisonOperator", "comparisonOperator"),
		Threshold:               getAlarmFloatParam(req.Parameters, "Threshold", "threshold"),
		EvaluationPeriods:       int32(getAlarmIntParam(req.Parameters, "EvaluationPeriods", "evaluationPeriods")),
		Period:                  int32(getAlarmIntParam(req.Parameters, "Period", "period")),
		Statistic:               getAlarmStringParam(req.Parameters, "Statistic", "statistic"),
		ExtendedStatistic:       getAlarmStringParam(req.Parameters, "ExtendedStatistic", "extendedStatistic"),
		TreatMissingData:        getAlarmStringParam(req.Parameters, "TreatMissingData", "treatMissingData"),
		AlarmDescription:        getAlarmStringParam(req.Parameters, "AlarmDescription", "alarmDescription"),
		ActionsEnabled:          getAlarmBoolParam(req.Parameters, []string{"ActionsEnabled", "actionsEnabled"}, true),
		DatapointsToAlarm:       int32(getAlarmIntParam(req.Parameters, "DatapointsToAlarm", "datapointsToAlarm")),
		AlarmActions:            getAlarmStringListParam(req.Parameters, "AlarmActions", "alarmActions"),
		OKActions:               getAlarmStringListParam(req.Parameters, "OKActions", "okActions"),
		InsufficientDataActions: getAlarmStringListParam(req.Parameters, "InsufficientDataActions", "insufficientDataActions"),
		Unit:                    getAlarmStringParam(req.Parameters, "Unit", "unit"),
		ThresholdMetricID:       getAlarmStringParam(req.Parameters, "ThresholdMetricId", "thresholdMetricId"),
		EvaluationCriteria:      parseEvaluationCriteria(req.Parameters),
	}

	if metricsRaw, ok := req.Parameters["Metrics"]; ok {
		input.Metrics = parseMetricDataQueries(metricsRaw)
	} else if metricsRaw, ok := req.Parameters["metrics"]; ok {
		input.Metrics = parseMetricDataQueries(metricsRaw)
	}

	var tagErr error
	input.Tags, tagErr = parseAndValidateAlarmTags(req.Parameters)
	if tagErr != nil {
		return nil, tagErr
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	arn, err := s.putMetricAlarmCore(store, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"AlarmArn": arn,
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

	input := &DescribeAlarmsInput{
		AlarmNamePrefix:     alarmNamePrefix,
		AlarmNames:          alarmNames,
		StateValueFilter:    stateValueFilter,
		ActionPrefix:        actionPrefix,
		AlarmTypes:          alarmTypes,
		ChildrenOfAlarmName: childrenOfAlarmName,
		ParentsOfAlarmName:  parentsOfAlarmName,
		NextToken:           pagination.GetMarker(req.Parameters, "NextToken"),
		MaxRecords:          pagination.GetMaxItems(req.Parameters, 100, "MaxRecords"),
	}

	alarms, nextToken, err := s.describeAlarmsCore(store, input)
	if err != nil {
		return nil, err
	}

	return s.buildDescribeAlarmsResponse(alarms, alarmTypes, nextToken), nil
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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	alarms, nextMarker, err := s.describeAlarmsForMetricCore(store, &DescribeAlarmsForMetricInput{
		Namespace:         getAlarmStringParam(req.Parameters, "Namespace", "namespace"),
		MetricName:        getAlarmStringParam(req.Parameters, "MetricName", "metricName"),
		Dimensions:        parseAlarmDimensions(req.Parameters),
		Statistic:         getAlarmStringParam(req.Parameters, "Statistic", "statistic"),
		ExtendedStatistic: getAlarmStringParam(req.Parameters, "ExtendedStatistic", "extendedStatistic"),
		Period:            int32(getAlarmIntParam(req.Parameters, "Period", "period")),
		Unit:              getAlarmStringParam(req.Parameters, "Unit", "unit"),
		NextToken:         pagination.GetMarker(req.Parameters, "NextToken"),
		MaxRecords:        pagination.GetMaxItems(req.Parameters, 100, "MaxRecords"),
	})
	if err != nil {
		return nil, err
	}

	metricAlarms := make([]map[string]interface{}, len(alarms))
	for i, alarm := range alarms {
		metricAlarms[i] = alarmToResponse(alarm)
	}

	resp := map[string]interface{}{
		"MetricAlarms": metricAlarms,
	}
	if nextMarker != "" {
		resp["NextToken"] = nextMarker
	}
	return resp, nil
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

// DeleteAlarms deletes one or more CloudWatch alarms.
func (s *CloudWatchService) DeleteAlarms(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	alarmNames := getAlarmStringListParam(req.Parameters, "AlarmNames", "alarmNames")
	if len(alarmNames) == 0 {
		return nil, ErrInvalidParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteAlarmsCore(store, &DeleteAlarmsInput{AlarmNames: alarmNames}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// SetAlarmState sets the state of an alarm.
func (s *CloudWatchService) SetAlarmState(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	alarmName := getAlarmStringParam(req.Parameters, "AlarmName", "alarmName")
	stateValue := getAlarmStringParam(req.Parameters, "StateValue", "stateValue")
	stateReason := getAlarmStringParam(req.Parameters, "StateReason", "stateReason")
	stateReasonData := getAlarmStringParam(req.Parameters, "StateReasonData", "stateReasonData")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := s.setAlarmStateCore(store, &SetAlarmStateInput{
		AlarmName:       alarmName,
		StateValue:      stateValue,
		StateReason:     stateReason,
		StateReasonData: stateReasonData,
	}); err != nil {
		return nil, err
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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.enableAlarmActionsCore(store, &EnableAlarmActionsInput{
		AlarmNames: getAlarmStringListParam(req.Parameters, "AlarmNames", "alarmNames"),
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DisableAlarmActions disables actions for the specified alarms.
func (s *CloudWatchService) DisableAlarmActions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.disableAlarmActionsCore(store, &EnableAlarmActionsInput{
		AlarmNames: getAlarmStringListParam(req.Parameters, "AlarmNames", "alarmNames"),
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DescribeAlarmHistory retrieves the history for the specified alarm.
func (s *CloudWatchService) DescribeAlarmHistory(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
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

	items, nextMarker, err := s.describeAlarmHistoryCore(store, &DescribeAlarmHistoryInput{
		AlarmName:       getAlarmStringParam(req.Parameters, "AlarmName", "alarmName"),
		AlarmTypes:      alarmTypes,
		HistoryItemType: getAlarmStringParam(req.Parameters, "HistoryItemType", "historyItemType"),
		StartDate:       parseTimestampFromMap(req.Parameters, "StartDate"),
		EndDate:         parseTimestampFromMap(req.Parameters, "EndDate"),
		ScanBy:          getAlarmStringParam(req.Parameters, "ScanBy", "scanBy"),
		NextToken:       pagination.GetMarker(req.Parameters, "NextToken"),
		MaxRecords:      pagination.GetMaxItems(req.Parameters, 100, "MaxRecords"),
	})
	if err != nil {
		return nil, err
	}

	historyItems := make([]map[string]interface{}, 0, len(items))
	for _, entry := range items {
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
		historyItems = append(historyItems, item)
	}

	resp := map[string]interface{}{
		"AlarmHistoryItems": historyItems,
	}
	if nextMarker != "" {
		resp["NextToken"] = nextMarker
	}
	return resp, nil
}

// PutCompositeAlarm creates or updates a composite alarm.
func (s *CloudWatchService) PutCompositeAlarm(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	tags, tagErr := parseAndValidateAlarmTags(req.Parameters)
	if tagErr != nil {
		return nil, tagErr
	}

	arn, err := s.putCompositeAlarmCore(store, &PutCompositeAlarmInput{
		AlarmName:                   getAlarmStringParam(req.Parameters, "AlarmName", "alarmName"),
		AlarmRule:                   getAlarmStringParam(req.Parameters, "AlarmRule", "alarmRule"),
		AlarmDescription:            getAlarmStringParam(req.Parameters, "AlarmDescription", "alarmDescription"),
		ActionsEnabled:              getAlarmBoolParam(req.Parameters, []string{"ActionsEnabled", "actionsEnabled"}, true),
		AlarmActions:                getAlarmStringListParam(req.Parameters, "AlarmActions", "alarmActions"),
		OKActions:                   getAlarmStringListParam(req.Parameters, "OKActions", "okActions"),
		InsufficientDataActions:     getAlarmStringListParam(req.Parameters, "InsufficientDataActions", "insufficientDataActions"),
		ActionsSuppressor:           getAlarmStringParam(req.Parameters, "ActionsSuppressor", "actionsSuppressor"),
		ActionsSuppressorWaitPeriod: int32(getAlarmIntParam(req.Parameters, "ActionsSuppressorWaitPeriod", "actionsSuppressorWaitPeriod")),
		ActionsSuppressorExtPeriod:  int32(getAlarmIntParam(req.Parameters, "ActionsSuppressorExtensionPeriod", "actionsSuppressorExtPeriod")),
		Tags:                        tags,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"AlarmArn": arn,
	}, nil
}

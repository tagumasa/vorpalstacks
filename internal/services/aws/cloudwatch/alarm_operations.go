package cloudwatch

import (
	"context"
	"errors"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
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
			case int64:
				return int(val)
			case float64:
				return int(val)
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

func alarmToResponse(alarm *cwstore.Alarm, stateUpdatedTs time.Time) map[string]interface{} {
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
		"StateUpdatedTimestamp": stateUpdatedTs.UTC().UnixMilli(),
		"TreatMissingData":      alarm.TreatMissingData,
		"ActionsEnabled":        alarm.ActionsEnabled,
	}
	if alarm.AlarmType != "" {
		result["AlarmType"] = alarm.AlarmType
	}
	if alarm.AlarmRule != "" {
		result["AlarmRule"] = alarm.AlarmRule
	}
	if alarm.AlarmDescription != "" {
		result["AlarmDescription"] = alarm.AlarmDescription
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
func (s *CloudWatchService) PutMetricAlarm(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	alarmName := getAlarmStringParam(req.Parameters, "AlarmName", "alarmName")
	if alarmName == "" {
		return nil, ErrInvalidParameter
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

	if alarm.ComparisonOperator == "" {
		alarm.ComparisonOperator = "GreaterThanOrEqualToThreshold"
	}
	if alarm.EvaluationPeriods == 0 {
		alarm.EvaluationPeriods = 1
	}
	if alarm.Period == 0 {
		alarm.Period = 60
	}
	if alarm.Statistic == "" {
		alarm.Statistic = "Average"
	}
	if alarm.TreatMissingData == "" {
		alarm.TreatMissingData = "missing"
	}

	tagList := tagutil.ParseTags(req.Parameters, "Tags")
	if len(tagList) == 0 {
		tagList = tagutil.ParseTags(req.Parameters, "tags")
	}
	alarm.Tags = tagutil.ToMap(tagList)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, err := store.alarms.GetAlarm(alarmName)
	if err == nil && existing != nil {
		alarm.ARN = existing.ARN
		alarm.CreatedAt = existing.CreatedAt
		alarm.State = existing.State
		alarm.StateUpdatedTimestamp = existing.StateUpdatedTimestamp
		if err := store.alarms.UpdateAlarm(alarm); err != nil {
			return nil, err
		}
	} else {
		created, err := store.alarms.CreateAlarm(alarm)
		if err != nil {
			return nil, err
		}
		alarm = created
	}

	if len(alarm.Tags) > 0 {
		if err := store.alarms.TagResource(alarm.ARN, alarm.Tags); err != nil {
			return nil, err
		}
	}

	store.alarms.AddAlarmHistory(&cwstore.AlarmHistoryEntry{
		AlarmName:       alarmName,
		AlarmType:       cwstore.AlarmTypeMetricAlarm,
		Timestamp:       time.Now().UTC().UnixMilli(),
		HistoryItemType: cwstore.HistoryItemTypeConfigurationUpdate,
		HistorySummary:  "Alarm was created or updated",
	})

	return map[string]interface{}{
		"AlarmArn": alarm.ARN,
	}, nil
}

// DescribeAlarms returns a list of alarms.
func (s *CloudWatchService) DescribeAlarms(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	alarmNamePrefix := getAlarmStringParam(req.Parameters, "AlarmNamePrefix", "alarmNamePrefix")

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

	var alarms []*cwstore.Alarm
	if len(alarmNames) > 0 {
		for _, name := range alarmNames {
			alarm, err := store.alarms.GetAlarm(name)
			if err == nil {
				alarms = append(alarms, alarm)
			}
		}
	} else {
		alarms, err = store.alarms.ListAlarms(alarmNamePrefix)
		if err != nil {
			return nil, err
		}
	}

	if len(alarmTypes) > 0 {
		typeSet := make(map[string]bool)
		for _, t := range alarmTypes {
			typeSet[t] = true
		}
		filtered := make([]*cwstore.Alarm, 0, len(alarms))
		for _, alarm := range alarms {
			alarmType := alarm.AlarmType
			if alarmType == "" {
				alarmType = cwstore.AlarmTypeMetricAlarm
			}
			if typeSet[alarmType] {
				filtered = append(filtered, alarm)
			}
		}
		alarms = filtered
	}

	metricAlarms := make([]map[string]interface{}, 0)
	compositeAlarms := make([]map[string]interface{}, 0)
	for _, alarm := range alarms {
		stateUpdatedTs := alarm.StateUpdatedTimestamp
		if stateUpdatedTs.IsZero() {
			stateUpdatedTs = time.Now()
		}
		resp := alarmToResponse(alarm, stateUpdatedTs)
		if alarm.AlarmType == cwstore.AlarmTypeCompositeAlarm {
			compositeAlarms = append(compositeAlarms, resp)
		} else {
			metricAlarms = append(metricAlarms, resp)
		}
	}
	return map[string]interface{}{
		"MetricAlarms":    metricAlarms,
		"CompositeAlarms": compositeAlarms,
	}, nil
}

// DescribeAlarmsForMetric returns alarms for a specific metric.
func (s *CloudWatchService) DescribeAlarmsForMetric(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	namespace := getAlarmStringParam(req.Parameters, "Namespace", "namespace")
	metricName := getAlarmStringParam(req.Parameters, "MetricName", "metricName")

	if namespace == "" || metricName == "" {
		return nil, ErrInvalidParameter
	}

	dimensions := parseAlarmDimensions(req.Parameters)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	alarms, err := store.alarms.ListAlarms("")
	if err != nil {
		return nil, err
	}

	var matchedAlarms []*cwstore.Alarm
	for _, alarm := range alarms {
		if alarm.Namespace == namespace && alarm.MetricName == metricName {
			if len(dimensions) > 0 && len(alarm.Dimensions) > 0 {
				if dimensionsMatch(alarm.Dimensions, dimensions) {
					matchedAlarms = append(matchedAlarms, alarm)
				}
			} else {
				matchedAlarms = append(matchedAlarms, alarm)
			}
		}
	}

	metricAlarms := make([]map[string]interface{}, len(matchedAlarms))
	for i, alarm := range matchedAlarms {
		stateUpdatedTs := alarm.StateUpdatedTimestamp
		if stateUpdatedTs.IsZero() {
			stateUpdatedTs = time.Now()
		}
		metricAlarms[i] = alarmToResponse(alarm, stateUpdatedTs)
	}

	return map[string]interface{}{
		"MetricAlarms": metricAlarms,
	}, nil
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

	for _, name := range alarmNames {
		if err := store.alarms.DeleteAlarm(name); err != nil {
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

	if err := store.alarms.SetAlarmState(alarmName, stateValue, stateReason); err != nil {
		if errors.Is(err, cwstore.ErrAlarmNotFound) {
			return nil, ErrAlarmNotFound
		}
		return nil, err
	}

	alarm, _ := store.alarms.GetAlarm(alarmName)
	alarmType := cwstore.AlarmTypeMetricAlarm
	if alarm != nil && alarm.AlarmType == cwstore.AlarmTypeCompositeAlarm {
		alarmType = cwstore.AlarmTypeCompositeAlarm
	}
	store.alarms.AddAlarmHistory(&cwstore.AlarmHistoryEntry{
		AlarmName:       alarmName,
		AlarmType:       alarmType,
		Timestamp:       time.Now().UTC().UnixMilli(),
		HistoryItemType: cwstore.HistoryItemTypeStateUpdate,
		HistorySummary:  stateReason,
	})

	return response.EmptyResponse(), nil
}

// TagResource adds tags to a CloudWatch alarm.
func (s *CloudWatchService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := request.GetStringParam(req.Parameters, "ResourceARN")
	if resourceARN == "" {
		return nil, ErrInvalidParameter
	}
	tags := tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags")
	if len(tags) == 0 {
		return nil, ErrInvalidParameter
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return response.EmptyResponse(), store.alarms.TagResource(resourceARN, tagutil.ToMap(tags))
}

// UntagResource removes tags from a CloudWatch alarm.
func (s *CloudWatchService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := request.GetStringParam(req.Parameters, "ResourceARN")
	if resourceARN == "" {
		return nil, ErrInvalidParameter
	}
	tagKeys := request.GetStringList(req.Parameters, "TagKeys")
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return response.EmptyResponse(), store.alarms.UntagResource(resourceARN, tagKeys)
}

// ListTagsForResource lists tags for a CloudWatch alarm.
func (s *CloudWatchService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceARN := request.GetStringParam(req.Parameters, "ResourceARN")
	if resourceARN == "" {
		return nil, ErrInvalidParameter
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	tags, _ := store.alarms.ListTags(resourceARN)
	return map[string]interface{}{
		"Tags": tagutil.MapToResponse(tags),
	}, nil
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
		alarm, _ := store.alarms.GetAlarm(name)
		alarmType := cwstore.AlarmTypeMetricAlarm
		if alarm != nil && alarm.AlarmType == cwstore.AlarmTypeCompositeAlarm {
			alarmType = cwstore.AlarmTypeCompositeAlarm
		}
		store.alarms.AddAlarmHistory(&cwstore.AlarmHistoryEntry{
			AlarmName:       name,
			AlarmType:       alarmType,
			Timestamp:       time.Now().UTC().UnixMilli(),
			HistoryItemType: cwstore.HistoryItemTypeAction,
			HistorySummary:  "Alarm actions enabled",
		})
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
		alarm, _ := store.alarms.GetAlarm(name)
		alarmType := cwstore.AlarmTypeMetricAlarm
		if alarm != nil && alarm.AlarmType == cwstore.AlarmTypeCompositeAlarm {
			alarmType = cwstore.AlarmTypeCompositeAlarm
		}
		store.alarms.AddAlarmHistory(&cwstore.AlarmHistoryEntry{
			AlarmName:       name,
			AlarmType:       alarmType,
			Timestamp:       time.Now().UTC().UnixMilli(),
			HistoryItemType: cwstore.HistoryItemTypeAction,
			HistorySummary:  "Alarm actions disabled",
		})
	}

	return response.EmptyResponse(), nil
}

// DescribeAlarmHistory retrieves the history for the specified alarm.
func (s *CloudWatchService) DescribeAlarmHistory(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	alarmName := getAlarmStringParam(req.Parameters, "AlarmName", "alarmName")
	historyItemType := getAlarmStringParam(req.Parameters, "HistoryItemType", "historyItemType")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	entries, err := store.alarms.ListAlarmHistory(alarmName, historyItemType)
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(entries))
	for _, entry := range entries {
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

	return map[string]interface{}{
		"AlarmHistoryItems": items,
	}, nil
}

// PutCompositeAlarm creates or updates a composite alarm.
func (s *CloudWatchService) PutCompositeAlarm(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	alarmName := getAlarmStringParam(req.Parameters, "AlarmName", "alarmName")
	alarmRule := getAlarmStringParam(req.Parameters, "AlarmRule", "alarmRule")

	if alarmName == "" || alarmRule == "" {
		return nil, ErrInvalidParameter
	}

	alarm := cwstore.NewAlarm(alarmName, "", "")
	alarm.AlarmRule = alarmRule
	alarm.AlarmType = cwstore.AlarmTypeCompositeAlarm
	alarm.Namespace = ""
	alarm.MetricName = ""

	alarm.AlarmDescription = getAlarmStringParam(req.Parameters, "AlarmDescription", "alarmDescription")
	alarm.ActionsEnabled = getAlarmBoolParam(req.Parameters, []string{"ActionsEnabled", "actionsEnabled"}, true)
	alarm.AlarmActions = getAlarmStringListParam(req.Parameters, "AlarmActions", "alarmActions")
	alarm.OKActions = getAlarmStringListParam(req.Parameters, "OKActions", "okActions")
	alarm.InsufficientDataActions = getAlarmStringListParam(req.Parameters, "InsufficientDataActions", "insufficientDataActions")

	tagList := tagutil.ParseTags(req.Parameters, "Tags")
	if len(tagList) == 0 {
		tagList = tagutil.ParseTags(req.Parameters, "tags")
	}
	alarm.Tags = tagutil.ToMap(tagList)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, err := store.alarms.GetAlarm(alarmName)
	if err == nil && existing != nil {
		alarm.ARN = existing.ARN
		alarm.CreatedAt = existing.CreatedAt
		alarm.State = existing.State
		alarm.StateUpdatedTimestamp = existing.StateUpdatedTimestamp
		if err := store.alarms.UpdateAlarm(alarm); err != nil {
			return nil, err
		}
	} else {
		created, err := store.alarms.CreateAlarm(alarm)
		if err != nil {
			return nil, err
		}
		alarm = created
	}

	if len(alarm.Tags) > 0 {
		if err := store.alarms.TagResource(alarm.ARN, alarm.Tags); err != nil {
			return nil, err
		}
	}

	store.alarms.AddAlarmHistory(&cwstore.AlarmHistoryEntry{
		AlarmName:       alarmName,
		AlarmType:       cwstore.AlarmTypeCompositeAlarm,
		Timestamp:       time.Now().UTC().UnixMilli(),
		HistoryItemType: cwstore.HistoryItemTypeConfigurationUpdate,
		HistorySummary:  "Alarm was created or updated",
	})

	return map[string]interface{}{
		"AlarmArn": alarm.ARN,
	}, nil
}

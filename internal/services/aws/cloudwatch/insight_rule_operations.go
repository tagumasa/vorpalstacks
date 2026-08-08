package cloudwatch

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/eventbus"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
	"vorpalstacks/internal/store/aws/common"
)

// PutInsightRule creates or updates a CloudWatch Contributor Insights
// rule. The rule analyses log data to find the top contributors (most
// active sources) for a given metric.
//
// AWS docs: https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/API_PutInsightRule.html
func (s *CloudWatchService) PutInsightRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	ruleName := getAlarmStringParam(req.Parameters, "RuleName", "ruleName")
	if ruleName == "" {
		return nil, awserrors.NewMissingParameter("RuleName is required")
	}

	ruleState := getAlarmStringParam(req.Parameters, "RuleState", "ruleState")
	if ruleState == "" {
		ruleState = "ENABLED"
	}
	if !validateRuleState(ruleState) {
		return nil, awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Invalid RuleState: %s. Must be ENABLED or DISABLED", ruleState))
	}

	ruleDefinition := getAlarmStringParam(req.Parameters, "RuleDefinition", "ruleDefinition")

	applyOnTransformed := false
	if v, ok := req.Parameters["ApplyOnTransformedLogs"]; ok {
		if b, ok := v.(bool); ok {
			applyOnTransformed = b
		}
	} else if v, ok := req.Parameters["applyOnTransformedLogs"]; ok {
		if b, ok := v.(bool); ok {
			applyOnTransformed = b
		}
	}

	tags, tagErr := parseAndValidateAlarmTags(req.Parameters)
	if tagErr != nil {
		return nil, tagErr
	}

	rule := &cwstore.InsightRule{
		Name:                   ruleName,
		State:                  ruleState,
		Definition:             ruleDefinition,
		ApplyOnTransformedLogs: applyOnTransformed,
		Tags:                   tags,
	}

	result, err := store.insightRules.PutInsightRule(rule)
	if err != nil {
		return nil, fmt.Errorf("failed to put insight rule: %w", err)
	}

	_ = result
	return map[string]interface{}{}, nil
}

// DeleteInsightRules deletes one or more Contributor Insights rules.
//
// AWS docs: https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/API_DeleteInsightRules.html
func (s *CloudWatchService) DeleteInsightRules(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	ruleNames := parseStringArrayParam(req.Parameters, "RuleNames", "ruleNames")
	if len(ruleNames) == 0 {
		return nil, awserrors.NewMissingParameter("RuleNames is required")
	}

	_, notFound, err := store.insightRules.DeleteInsightRules(ruleNames)
	if err != nil {
		return nil, fmt.Errorf("failed to delete insight rules: %w", err)
	}

	failures := make([]map[string]interface{}, 0, len(notFound))
	for _, name := range notFound {
		failures = append(failures, map[string]interface{}{
			"FailureName":        name,
			"ExceptionName":      "ResourceNotFoundException",
			"FailureDescription": fmt.Sprintf("Rule %s does not exist", name),
		})
	}

	return map[string]interface{}{
		"Failures": failures,
	}, nil
}

// DescribeInsightRules lists the Contributor Insights rules in the
// account.
//
// AWS docs: https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/API_DescribeInsightRules.html
func (s *CloudWatchService) DescribeInsightRules(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	marker := pagination.GetMarker(req.Parameters, "NextToken")
	maxResults := pagination.GetMaxItems(req.Parameters, 100, "MaxResults")

	opts := common.ListOptions{Marker: marker, MaxItems: maxResults}
	result, err := store.insightRules.ListInsightRulesPaginated(false, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list insight rules: %w", err)
	}

	results := make([]map[string]interface{}, 0, len(result.Items))
	for _, r := range result.Items {
		results = append(results, insightRuleToResponse(r))
	}

	resp := map[string]interface{}{
		"InsightRules": results,
	}
	if result.NextMarker != "" {
		resp["NextToken"] = result.NextMarker
	}
	return resp, nil
}

// EnableInsightRules enables the specified Contributor Insights rules.
//
// AWS docs: https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/API_EnableInsightRules.html
func (s *CloudWatchService) EnableInsightRules(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	ruleNames := parseStringArrayParam(req.Parameters, "RuleNames", "ruleNames")
	if len(ruleNames) == 0 {
		return nil, awserrors.NewMissingParameter("RuleNames is required")
	}

	_, notFound, err := store.insightRules.SetRuleState(ruleNames, "ENABLED")
	if err != nil {
		return nil, fmt.Errorf("failed to enable insight rules: %w", err)
	}

	failures := make([]map[string]interface{}, 0, len(notFound))
	for _, name := range notFound {
		failures = append(failures, map[string]interface{}{
			"FailureName":        name,
			"ExceptionName":      "ResourceNotFoundException",
			"FailureDescription": fmt.Sprintf("Rule %s does not exist", name),
		})
	}

	return map[string]interface{}{
		"Failures": failures,
	}, nil
}

// DisableInsightRules disables the specified Contributor Insights rules.
//
// AWS docs: https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/API_DisableInsightRules.html
func (s *CloudWatchService) DisableInsightRules(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	ruleNames := parseStringArrayParam(req.Parameters, "RuleNames", "ruleNames")
	if len(ruleNames) == 0 {
		return nil, awserrors.NewMissingParameter("RuleNames is required")
	}

	_, notFound, err := store.insightRules.SetRuleState(ruleNames, "DISABLED")
	if err != nil {
		return nil, fmt.Errorf("failed to disable insight rules: %w", err)
	}

	failures := make([]map[string]interface{}, 0, len(notFound))
	for _, name := range notFound {
		failures = append(failures, map[string]interface{}{
			"FailureName":        name,
			"ExceptionName":      "ResourceNotFoundException",
			"FailureDescription": fmt.Sprintf("Rule %s does not exist", name),
		})
	}

	return map[string]interface{}{
		"Failures": failures,
	}, nil
}

// GetInsightRuleReport returns a report for a Contributor Insights rule.
// The report includes aggregated contributor data and metric data points
// for the specified time range. Data is sourced from CloudTrail events
// queried via the event bus CloudTrailInvoker.
//
// AWS docs: https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/API_GetInsightRuleReport.html
func (s *CloudWatchService) GetInsightRuleReport(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	ruleName := getAlarmStringParam(req.Parameters, "RuleName", "ruleName")
	if ruleName == "" {
		return nil, awserrors.NewMissingParameter("RuleName is required")
	}

	rule, err := store.insightRules.GetInsightRule(ruleName)
	if err != nil {
		return nil, awserrors.NewResourceNotFoundException("InsightRule", ruleName)
	}

	startTime := parseTimestampFromMap(req.Parameters, "StartTime")
	endTime := parseTimestampFromMap(req.Parameters, "EndTime")
	period := getAlarmIntParam(req.Parameters, "Period", "period")
	maxContributorCount := getAlarmIntParam(req.Parameters, "MaxContributorCount", "maxContributorCount")
	if maxContributorCount <= 0 {
		maxContributorCount = 10
	}
	if maxContributorCount > 100 {
		maxContributorCount = 100
	}

	def := parseInsightRuleDefinition(rule.Definition)

	contributorKeyField := "Username"
	if strings.Contains(strings.ToLower(def.ContributorValue), "eventname") {
		contributorKeyField = "EventName"
	} else if strings.Contains(strings.ToLower(def.ContributorValue), "eventsource") {
		contributorKeyField = "EventSource"
	}

	aggregationStatistic := "Sum"
	if strings.EqualFold(def.AggregateOn, "Count") {
		aggregationStatistic = "Count"
	}

	keyLabels := []string{contributorKeyField}

	var events []eventbus.CloudTrailEventInfo
	if s.bus != nil && s.bus.CloudTrailInvoker() != nil {
		events = fetchInsightEvents(ctx, s.bus.CloudTrailInvoker(), reqCtx.Region, startTime, endTime)
	}

	type contributorAgg struct {
		key        string
		count      int
		datapoints map[time.Time]int
	}
	contributorMap := make(map[string]*contributorAgg)
	periodDuration := time.Duration(period) * time.Second

	for _, e := range events {
		var key string
		switch contributorKeyField {
		case "EventName":
			key = e.EventName
		case "EventSource":
			key = e.EventSource
		default:
			key = e.Username
		}
		if key == "" {
			continue
		}

		bucket := e.EventTime.Truncate(periodDuration)
		agg, exists := contributorMap[key]
		if !exists {
			agg = &contributorAgg{key: key, datapoints: make(map[time.Time]int)}
			contributorMap[key] = agg
		}
		agg.count++
		agg.datapoints[bucket]++
	}

	allContributors := make([]*contributorAgg, 0, len(contributorMap))
	for _, agg := range contributorMap {
		allContributors = append(allContributors, agg)
	}
	sort.Slice(allContributors, func(i, j int) bool { return allContributors[i].count > allContributors[j].count })

	if len(allContributors) > maxContributorCount {
		allContributors = allContributors[:maxContributorCount]
	}

	contributorList := make([]interface{}, 0, len(allContributors))
	for _, c := range allContributors {
		contributorList = append(contributorList, map[string]interface{}{
			"Keys":           []string{c.key},
			"AggregateValue": float64(c.count),
			"Datapoints":     buildContributorDatapoints(c.datapoints, startTime, endTime, periodDuration),
		})
	}

	aggregateValue := 0.0
	for _, c := range allContributors {
		aggregateValue += float64(c.count)
	}

	return map[string]interface{}{
		"KeyLabels":              keyLabels,
		"AggregationStatistic":   aggregationStatistic,
		"AggregateValue":         aggregateValue,
		"ApproximateUniqueCount": len(contributorMap),
		"Contributors":           contributorList,
		"MetricDatapoints":       buildInsightMetricDatapointsFromEvents(events, startTime, endTime, period, periodDuration, contributorKeyField),
	}, nil
}

// PutManagedInsightRules creates or updates managed insight
// rules for a specified resource.
//
// AWS docs: https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/API_PutManagedInsightRules.html
func (s *CloudWatchService) PutManagedInsightRules(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var managedRules []interface{}
	if v, ok := req.Parameters["ManagedRules"]; ok {
		managedRules, _ = v.([]interface{})
	} else if v, ok := req.Parameters["managedRules"]; ok {
		managedRules, _ = v.([]interface{})
	}

	if len(managedRules) == 0 {
		return nil, awserrors.NewMissingParameter("ManagedRules is required")
	}

	var failures []map[string]interface{}

	for _, raw := range managedRules {
		m, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		templateName := getAlarmStringParam(m, "TemplateName", "templateName")
		resourceARN := getAlarmStringParam(m, "ResourceARN", "resourceArn")
		tags := parseAlarmTags(m)

		if templateName == "" || resourceARN == "" {
			failures = append(failures, map[string]interface{}{
				"FailureName":        templateName,
				"ExceptionName":      "InvalidParameterValueException",
				"FailureDescription": "TemplateName and ResourceARN are required",
			})
			continue
		}

		// Generate a deterministic rule name from template and resource.
		ruleName := fmt.Sprintf("ManagedRule:%s:%s", templateName, resourceARN)

		rule := &cwstore.InsightRule{
			Name:         ruleName,
			State:        "ENABLED",
			ManagedRule:  true,
			TemplateName: templateName,
			ResourceARN:  resourceARN,
			Tags:         tags,
		}

		_, err := store.insightRules.PutInsightRule(rule)
		if err != nil {
			failures = append(failures, map[string]interface{}{
				"FailureName":        templateName,
				"ExceptionName":      "InternalServiceFault",
				"FailureDescription": err.Error(),
			})
		}
	}
	return map[string]interface{}{
		"Failures": failures,
	}, nil
}

// ListManagedInsightRules lists the managed Contributor Insights rules
// for a specified resource.
//
// AWS docs: https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/API_ListManagedInsightRules.html
func (s *CloudWatchService) ListManagedInsightRules(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	resourceARN := getAlarmStringParam(req.Parameters, "ResourceARN", "resourceArn")

	marker := pagination.GetMarker(req.Parameters, "NextToken")
	maxResults := pagination.GetMaxItems(req.Parameters, 100, "MaxResults")

	opts := common.ListOptions{Marker: marker, MaxItems: maxResults}
	result, err := store.insightRules.ListInsightRulesPaginated(true, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list managed insight rules: %w", err)
	}

	results := make([]map[string]interface{}, 0, len(result.Items))
	for _, r := range result.Items {
		if resourceARN != "" && r.ResourceARN != resourceARN {
			continue
		}
		results = append(results, managedInsightRuleToResponse(r))
	}

	resp := map[string]interface{}{
		"ManagedRules": results,
	}
	if result.NextMarker != "" {
		resp["NextToken"] = result.NextMarker
	}
	return resp, nil
}

// insightRuleToResponse serialises an InsightRule into the AWS API
// response format for DescribeInsightRules.
func insightRuleToResponse(r *cwstore.InsightRule) map[string]interface{} {
	resp := map[string]interface{}{
		"Name":        r.Name,
		"State":       r.State,
		"ManagedRule": r.ManagedRule,
	}
	if r.Schema != "" {
		resp["Schema"] = r.Schema
	}
	if r.Definition != "" {
		resp["Definition"] = r.Definition
	}
	if r.ApplyOnTransformedLogs {
		resp["ApplyOnTransformedLogs"] = true
	}
	return resp
}

// managedInsightRuleToResponse serialises an InsightRule into the AWS
// API response format for ListManagedInsightRules.
func managedInsightRuleToResponse(r *cwstore.InsightRule) map[string]interface{} {
	resp := map[string]interface{}{
		"TemplateName": r.TemplateName,
	}
	if r.ResourceARN != "" {
		resp["ResourceARN"] = r.ResourceARN
	}
	resp["RuleState"] = map[string]interface{}{
		"RuleName": r.Name,
		"State":    r.State,
	}
	return resp
}

// insightRuleDefinition holds the parsed fields from the JSON Definition
// string stored in an InsightRule. The Definition follows the AWS
// CloudWatchLogRule schema.
type insightRuleDefinition struct {
	ContributorValue string
	AggregateOn      string
}

// parseInsightRuleDefinition extracts the contributor key path and
// aggregation type from the rule definition JSON string.
func parseInsightRuleDefinition(definition string) insightRuleDefinition {
	var raw struct {
		AggregateOn string `json:"AggregateOn"`
		Contributor struct {
			Value string `json:"Value"`
		} `json:"Contributor"`
	}
	if definition == "" {
		return insightRuleDefinition{}
	}
	if err := json.Unmarshal([]byte(definition), &raw); err != nil {
		return insightRuleDefinition{}
	}
	return insightRuleDefinition{
		ContributorValue: raw.Contributor.Value,
		AggregateOn:      raw.AggregateOn,
	}
}

const insightEventMaxPages = 20

func fetchInsightEvents(ctx context.Context, invoker eventbus.CloudTrailInvoker, region string, startTime, endTime time.Time) []eventbus.CloudTrailEventInfo {
	if startTime.IsZero() || endTime.IsZero() {
		return nil
	}
	var all []eventbus.CloudTrailEventInfo
	nextToken := ""
	for page := 0; page < insightEventMaxPages; page++ {
		events, nt, err := invoker.LookupEvents(ctx, region, "", nextToken, startTime, endTime, 50)
		if err != nil {
			break
		}
		all = append(all, events...)
		if nt == "" {
			break
		}
		nextToken = nt
	}
	return all
}

// buildInsightMetricDatapointsFromEvents aggregates CloudTrail events
// into time buckets and computes per-bucket statistics for the
// GetInsightRuleReport response.
func buildInsightMetricDatapointsFromEvents(events []eventbus.CloudTrailEventInfo, startTime, endTime time.Time, period int, periodDuration time.Duration, contributorKeyField string) []map[string]interface{} {
	if startTime.IsZero() || endTime.IsZero() || period <= 0 {
		return []map[string]interface{}{}
	}

	type bucketAgg struct {
		contributorCounts map[string]int
		sampleCount       int
		maxValue          float64
		minValue          float64
	}
	buckets := make(map[time.Time]*bucketAgg)

	for _, e := range events {
		bucket := e.EventTime.Truncate(periodDuration)
		agg, ok := buckets[bucket]
		if !ok {
			agg = &bucketAgg{contributorCounts: make(map[string]int), minValue: -1}
			buckets[bucket] = agg
		}

		var key string
		switch contributorKeyField {
		case "EventName":
			key = e.EventName
		case "EventSource":
			key = e.EventSource
		default:
			key = e.Username
		}
		if key == "" {
			continue
		}
		agg.contributorCounts[key]++
		agg.sampleCount++
		v := float64(agg.contributorCounts[key])
		if v > agg.maxValue {
			agg.maxValue = v
		}
		if agg.minValue < 0 || v < agg.minValue {
			agg.minValue = v
		}
	}

	var datapoints []map[string]interface{}
	for t := startTime.Truncate(periodDuration); !t.After(endTime); t = t.Add(periodDuration) {
		agg, ok := buckets[t]
		if !ok {
			datapoints = append(datapoints, map[string]interface{}{
				"Timestamp":           t.Format("2006-01-02T15:04:05Z"),
				"UniqueContributors":  0,
				"MaxContributorValue": 0,
				"SampleCount":         0,
				"Average":             0.0,
				"Sum":                 0.0,
				"Minimum":             0.0,
				"Maximum":             0.0,
			})
			continue
		}

		totalContributors := len(agg.contributorCounts)
		sum := float64(agg.sampleCount)
		avg := 0.0
		if totalContributors > 0 {
			avg = sum / float64(totalContributors)
		}
		minVal := agg.minValue
		if minVal < 0 {
			minVal = 0
		}

		datapoints = append(datapoints, map[string]interface{}{
			"Timestamp":           t.Format("2006-01-02T15:04:05Z"),
			"UniqueContributors":  totalContributors,
			"MaxContributorValue": agg.maxValue,
			"SampleCount":         agg.sampleCount,
			"Average":             avg,
			"Sum":                 sum,
			"Minimum":             minVal,
			"Maximum":             agg.maxValue,
		})
	}

	return datapoints
}

// buildContributorDatapoints builds the per-contributor time series
// for a single contributor's datapoints in the Contributors array.
func buildContributorDatapoints(dpMap map[time.Time]int, startTime, endTime time.Time, periodDuration time.Duration) []map[string]interface{} {
	var datapoints []map[string]interface{}
	for t := startTime.Truncate(periodDuration); !t.After(endTime); t = t.Add(periodDuration) {
		count := dpMap[t]
		datapoints = append(datapoints, map[string]interface{}{
			"Timestamp": t.Format("2006-01-02T15:04:05Z"),
			"Value":     float64(count),
		})
	}
	return datapoints
}

// parseStringArrayParam extracts a string array parameter from the
// request parameters, trying both PascalCase and camelCase keys.
func parseStringArrayParam(params map[string]interface{}, keys ...string) []string {
	for _, key := range keys {
		if v, ok := params[key]; ok {
			if arr, ok := v.([]interface{}); ok {
				result := make([]string, 0, len(arr))
				for _, item := range arr {
					if s, ok := item.(string); ok {
						result = append(result, s)
					}
				}
				return result
			}
		}
	}
	return nil
}

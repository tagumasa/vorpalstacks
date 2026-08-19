package cloudwatch

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/eventbus"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
)

// PutInsightRule creates or updates a CloudWatch Contributor Insights
// rule.
func (s *CloudWatchService) PutInsightRule(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

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

	err = s.putInsightRuleCore(store, &PutInsightRuleInput{
		RuleName:               getAlarmStringParam(req.Parameters, "RuleName", "ruleName"),
		RuleState:              getAlarmStringParam(req.Parameters, "RuleState", "ruleState"),
		RuleDefinition:         getAlarmStringParam(req.Parameters, "RuleDefinition", "ruleDefinition"),
		ApplyOnTransformedLogs: applyOnTransformed,
		Tags:                   tags,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// DeleteInsightRules deletes one or more Contributor Insights rules.
func (s *CloudWatchService) DeleteInsightRules(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	notFound, err := s.deleteInsightRulesCore(store, &DeleteInsightRulesInput{
		RuleNames: parseStringArrayParam(req.Parameters, "RuleNames", "ruleNames"),
	})
	if err != nil {
		return nil, err
	}

	return buildInsightRuleFailures(notFound), nil
}

// DescribeInsightRules lists the Contributor Insights rules.
func (s *CloudWatchService) DescribeInsightRules(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	items, nextMarker, err := s.describeInsightRulesCore(store, &DescribeInsightRulesInput{
		NextToken:  pagination.GetMarker(req.Parameters, "NextToken"),
		MaxResults: pagination.GetMaxItems(req.Parameters, 100, "MaxResults"),
	})
	if err != nil {
		return nil, err
	}

	results := make([]map[string]interface{}, 0, len(items))
	for _, r := range items {
		results = append(results, insightRuleToResponse(r))
	}

	// DynamoDB-derived contributor rules are merged into the final page of
	// the listing; they live on the DynamoDB side and are derived from
	// insights-enabled tables.
	if nextMarker == "" && s.bus != nil && s.bus.DynamoDBInvoker() != nil {
		if dynamoRules, err := s.bus.DynamoDBInvoker().ContributorRules(ctx, reqCtx.Region); err == nil {
			for _, r := range dynamoRules {
				results = append(results, map[string]interface{}{
					"Name":        r.Name,
					"State":       "ENABLED",
					"ManagedRule": true,
				})
			}
		}
	}

	resp := map[string]interface{}{
		"InsightRules": results,
	}
	if nextMarker != "" {
		resp["NextToken"] = nextMarker
	}
	return resp, nil
}

// EnableInsightRules enables the specified Contributor Insights rules.
func (s *CloudWatchService) EnableInsightRules(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	notFound, err := s.setInsightRuleStateCore(store, &SetInsightRuleStateInput{
		RuleNames: parseStringArrayParam(req.Parameters, "RuleNames", "ruleNames"),
		State:     "ENABLED",
	})
	if err != nil {
		return nil, err
	}

	return buildInsightRuleFailures(notFound), nil
}

// DisableInsightRules disables the specified Contributor Insights rules.
func (s *CloudWatchService) DisableInsightRules(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	notFound, err := s.setInsightRuleStateCore(store, &SetInsightRuleStateInput{
		RuleNames: parseStringArrayParam(req.Parameters, "RuleNames", "ruleNames"),
		State:     "DISABLED",
	})
	if err != nil {
		return nil, err
	}

	return buildInsightRuleFailures(notFound), nil
}

// GetInsightRuleReport returns a report for a Contributor Insights rule.
func (s *CloudWatchService) GetInsightRuleReport(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	ruleName := getAlarmStringParam(req.Parameters, "RuleName", "ruleName")

	// DynamoDB contributor rules are derived from insights-enabled tables;
	// their reports come from the DynamoDB access aggregation through the
	// bus instead of the CloudTrail event path.
	if tableName, layout, isDynamo := dynamoDBContributorRuleParts(ruleName); isDynamo {
		return s.getDynamoDBContributorReport(ctx, reqCtx, ruleName, tableName, layout, req.Parameters)
	}

	rule, err := s.getInsightRuleCore(store, ruleName)
	if err != nil {
		return nil, err
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

// dynamoDBContributorRulePrefix is the rule-name prefix of the DynamoDB
// contributor insights rules surfaced through the DynamoDB API.
const dynamoDBContributorRulePrefix = "DynamoDBContributorInsights-"

// dynamoDBContributorRuleParts reports whether a rule name is a
// DynamoDB-derived contributor insights rule and extracts the tracked
// table and key layout from it. Table names may contain hyphens, so the
// layout is the first segment, the creation timestamp the last, and
// everything between them the table name.
func dynamoDBContributorRuleParts(ruleName string) (tableName, layout string, ok bool) {
	rest, found := strings.CutPrefix(ruleName, dynamoDBContributorRulePrefix)
	if !found {
		return "", "", false
	}
	parts := strings.Split(rest, "-")
	if len(parts) < 3 {
		return "", "", false
	}
	switch parts[0] {
	case "PKC", "SKC", "PKT", "SKT":
	default:
		return "", "", false
	}
	if _, err := strconv.ParseInt(parts[len(parts)-1], 10, 64); err != nil {
		return "", "", false
	}
	return strings.Join(parts[1:len(parts)-1], "-"), parts[0], true
}

// getDynamoDBContributorReport serves GetInsightRuleReport for a
// DynamoDB-derived rule from the DynamoDB key access aggregation. Only
// MaxContributorValue and Maximum carry meaningful values, matching the
// documented behaviour for DynamoDB-managed rules.
func (s *CloudWatchService) getDynamoDBContributorReport(ctx context.Context, reqCtx *request.RequestContext, ruleName, tableName, layout string, params map[string]interface{}) (interface{}, error) {
	startTime := parseTimestampFromMap(params, "StartTime")
	endTime := parseTimestampFromMap(params, "EndTime")
	period := getAlarmIntParam(params, "Period", "period")
	maxContributorCount := getAlarmIntParam(params, "MaxContributorCount", "maxContributorCount")
	if maxContributorCount <= 0 {
		maxContributorCount = 10
	}
	// DynamoDB-managed rules expose at most twenty-five contributors.
	if maxContributorCount > 25 {
		maxContributorCount = 25
	}
	if period <= 0 {
		period = 60
	}

	if s.bus == nil || s.bus.DynamoDBInvoker() == nil {
		return nil, awserrors.NewResourceNotFoundException("InsightRule", ruleName)
	}
	stats, err := s.bus.DynamoDBInvoker().ContributorStats(ctx, reqCtx.Region, tableName, layout, startTime, endTime, maxContributorCount)
	if err != nil {
		return nil, err
	}

	keyLabels := []string{"PartitionKey"}
	if layout == "SKC" || layout == "SKT" {
		keyLabels = []string{"PartitionKey", "SortKey"}
	}

	contributors := make([]interface{}, 0, len(stats))
	aggregateValue := 0.0
	maxContributorValue := 0.0
	for _, stat := range stats {
		aggregateValue += stat.Units
		if stat.Units > maxContributorValue {
			maxContributorValue = stat.Units
		}
		contributors = append(contributors, map[string]interface{}{
			"Keys": stat.Keys,
			// Each contributor carries its own approximate contribution,
			// not the metric-level statistics.
			"ApproximateAggregateValue": stat.Units,
			"Datapoints": []map[string]interface{}{{
				"Timestamp":        endTime.Unix(),
				"ApproximateValue": stat.Units,
			}},
		})
	}

	return map[string]interface{}{
		"KeyLabels":              keyLabels,
		"AggregationStatistic":   "Sum",
		"AggregationPeriod":      period,
		"AggregateValue":         aggregateValue,
		"ApproximateUniqueCount": len(stats),
		"Contributors":           contributors,
		// For DynamoDB-managed rules only MaxContributorValue and Maximum
		// carry useful statistics in the metric-level data points.
		"MetricDatapoints": []map[string]interface{}{{
			"Timestamp":           endTime.Unix(),
			"MaxContributorValue": maxContributorValue,
			"Maximum":             aggregateValue,
		}},
	}, nil
}

// PutManagedInsightRules creates or updates managed insight rules.
func (s *CloudWatchService) PutManagedInsightRules(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var rawRules []interface{}
	if v, ok := req.Parameters["ManagedRules"]; ok {
		rawRules, _ = v.([]interface{})
	} else if v, ok := req.Parameters["managedRules"]; ok {
		rawRules, _ = v.([]interface{})
	}

	if len(rawRules) == 0 {
		return nil, awserrors.NewMissingParameter("ManagedRules is required")
	}

	managedRules := make([]PutManagedInsightRuleItem, 0, len(rawRules))
	for _, raw := range rawRules {
		m, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		managedRules = append(managedRules, PutManagedInsightRuleItem{
			TemplateName: getAlarmStringParam(m, "TemplateName", "templateName"),
			ResourceARN:  getAlarmStringParam(m, "ResourceARN", "resourceArn"),
			Tags:         parseAlarmTags(m),
		})
	}

	failures := s.putManagedInsightRulesCore(store, &PutManagedInsightRulesInput{ManagedRules: managedRules})
	return map[string]interface{}{"Failures": failures}, nil
}

// ListManagedInsightRules lists the managed Contributor Insights rules.
func (s *CloudWatchService) ListManagedInsightRules(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	items, nextMarker, err := s.listManagedInsightRulesCore(store, &ListManagedInsightRulesInput{
		ResourceARN: getAlarmStringParam(req.Parameters, "ResourceARN", "resourceArn"),
		NextToken:   pagination.GetMarker(req.Parameters, "NextToken"),
		MaxResults:  pagination.GetMaxItems(req.Parameters, 100, "MaxResults"),
	})
	if err != nil {
		return nil, err
	}

	results := make([]map[string]interface{}, 0, len(items))
	for _, r := range items {
		results = append(results, managedInsightRuleToResponse(r))
	}

	resp := map[string]interface{}{
		"ManagedRules": results,
	}
	if nextMarker != "" {
		resp["NextToken"] = nextMarker
	}
	return resp, nil
}

// buildInsightRuleFailures converts a list of not-found rule names into
// the AWS API Failures response format.
func buildInsightRuleFailures(notFound []string) map[string]interface{} {
	failures := make([]map[string]interface{}, 0, len(notFound))
	for _, name := range notFound {
		failures = append(failures, map[string]interface{}{
			"FailureName":        name,
			"ExceptionName":      "ResourceNotFoundException",
			"FailureDescription": fmt.Sprintf("Rule %s does not exist", name),
		})
	}
	return map[string]interface{}{"Failures": failures}
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
// string stored in an InsightRule.
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
// into time buckets for the GetInsightRuleReport response.
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

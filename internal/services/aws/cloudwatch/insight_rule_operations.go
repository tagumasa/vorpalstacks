package cloudwatch

import (
	"context"
	"fmt"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
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

	tags := parseAlarmTags(req.Parameters)

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

	rules, err := store.insightRules.ListInsightRules(false)
	if err != nil {
		return nil, fmt.Errorf("failed to list insight rules: %w", err)
	}

	results := make([]map[string]interface{}, 0, len(rules))
	for _, r := range rules {
		results = append(results, insightRuleToResponse(r))
	}

	return map[string]interface{}{
		"InsightRules": results,
		"NextToken":    nil,
	}, nil
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
// for the specified time range.
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

	// Build a basic report structure. The full contributor analysis
	// would require CloudTrail log processing; this returns the report
	// structure with empty contributors and aggregated metric datapoints.
	keyLabels := []string{}
	aggregationStatistic := "Sum"

	return map[string]interface{}{
		"KeyLabels":              keyLabels,
		"AggregationStatistic":   aggregationStatistic,
		"AggregateValue":         0.0,
		"ApproximateUniqueCount": 0,
		"Contributors":           []interface{}{},
		"MetricDatapoints":       buildInsightMetricDatapoints(startTime, endTime, period, rule),
	}, nil
}

// PutManagedInsightRules creates or updates managed Contributor Insights
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

	rules, err := store.insightRules.ListInsightRules(true)
	if err != nil {
		return nil, fmt.Errorf("failed to list managed insight rules: %w", err)
	}

	// Filter by resource ARN if provided.
	results := make([]map[string]interface{}, 0, len(rules))
	for _, r := range rules {
		if resourceARN != "" && r.ResourceARN != resourceARN {
			continue
		}
		results = append(results, managedInsightRuleToResponse(r))
	}

	return map[string]interface{}{
		"ManagedRules": results,
		"NextToken":    nil,
	}, nil
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

// buildInsightMetricDatapoints builds a list of metric data point
// structures for the GetInsightRuleReport response. Each data point
// represents a period bucket with aggregated contributor statistics.
func buildInsightMetricDatapoints(startTime, endTime time.Time, period int, rule *cwstore.InsightRule) []map[string]interface{} {
	if startTime.IsZero() || endTime.IsZero() || period <= 0 {
		return []map[string]interface{}{}
	}

	periodDuration := time.Duration(period) * time.Second
	var datapoints []map[string]interface{}

	for t := startTime.Truncate(periodDuration); !t.After(endTime); t = t.Add(periodDuration) {
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

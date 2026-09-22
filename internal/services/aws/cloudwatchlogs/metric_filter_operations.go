// Copyright 2026 Vorpalstacks Authors
// SPDX-License-Identifier: Apache-2.0

package cloudwatchlogs

import (
	"context"
	"fmt"
	"strconv"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	"vorpalstacks/pkg/filterpattern"
)

// PutMetricFilter creates or updates a metric filter for the specified CloudWatch Logs log group.
func (s *LogsService) PutMetricFilter(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	filterName := request.GetParamLowerFirst(req.Parameters, "FilterName")
	filterPattern := request.GetParamLowerFirst(req.Parameters, "FilterPattern")

	if err := validateLogGroupName(logGroupName); err != nil {
		return nil, err
	}
	if err := validateFilterName(filterName); err != nil {
		return nil, err
	}

	var transformations []logsstore.MetricTransformation
	for i := 1; ; i++ {
		metricName := request.GetParamLowerFirst(req.Parameters, "MetricTransformations."+strconv.Itoa(i)+".MetricName")
		metricNamespace := request.GetParamLowerFirst(req.Parameters, "MetricTransformations."+strconv.Itoa(i)+".MetricNamespace")
		if metricName == "" && metricNamespace == "" {
			break
		}
		nameKey := "MetricTransformations." + strconv.Itoa(i) + ".MetricName"
		namespaceKey := "MetricTransformations." + strconv.Itoa(i) + ".MetricNamespace"
		valueKey := "MetricTransformations." + strconv.Itoa(i) + ".MetricValue"
		t := logsstore.MetricTransformation{
			MetricName:         metricName,
			MetricNamespace:    metricNamespace,
			MetricValue:        request.GetParamLowerFirst(req.Parameters, valueKey),
			Unit:               request.GetParamLowerFirst(req.Parameters, "MetricTransformations."+strconv.Itoa(i)+".Unit"),
			MetricNameSet:      request.HasParamLowerFirst(req.Parameters, nameKey),
			MetricNamespaceSet: request.HasParamLowerFirst(req.Parameters, namespaceKey),
			MetricValueSet:     request.HasParamLowerFirst(req.Parameters, valueKey),
		}
		dvKey := "MetricTransformations." + strconv.Itoa(i) + ".DefaultValue"
		if dvStr := request.GetParamLowerFirst(req.Parameters, dvKey); dvStr != "" {
			t.DefaultValue = request.GetFloatParam(req.Parameters, dvKey)
			t.DefaultValueSet = true
		}
		transformations = append(transformations, t)
	}

	if len(transformations) == 0 {
		transformations = parseMetricTransformationsFromMap(req)
	}

	fieldSelectionCriteria := request.GetParamLowerFirst(req.Parameters, "FieldSelectionCriteria")

	if err := s.putMetricFilterCore(logGroupName, filterName, filterPattern,
		request.HasParamLowerFirst(req.Parameters, "FilterPattern"), transformations,
		request.GetBoolParam(req.Parameters, "ApplyOnTransformedLogs"),
		fieldSelectionCriteria,
		request.GetStringList(req.Parameters, "EmitSystemFieldDimensions"),
		reqCtx.GetRegion()); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

func (s *LogsService) putMetricFilterCore(logGroupName, filterName, filterPattern string, filterPatternSet bool, transformations []logsstore.MetricTransformation, applyOnTransformedLogs bool, fieldSelectionCriteria string, emitSystemFieldDimensions []string, region string) error {
	// filterPattern is a required member whose shape allows the empty
	// string, so presence — not value — is the requiredness test; a
	// present empty pattern matches every event and stays legal.
	if !filterPatternSet {
		return errRequiredMember("filterPattern")
	}
	// The pattern validates in the Core — the syntax gate and the regex
	// quota below are the plane-independent contract, not the HTTP
	// handler's.
	if err := validateFilterPattern(filterPattern); err != nil {
		return err
	}
	// The MetricTransformations target's length trait is exactly 1..1 —
	// one transformation per metric filter.
	if len(transformations) == 0 {
		return NewLogsError("InvalidParameterException",
			"metricTransformations must contain exactly one transformation", 400)
	}
	if len(transformations) > 1 {
		return NewLogsError("InvalidParameterException",
			"metricTransformations must contain exactly one transformation; multiple transformations are not supported", 400)
	}
	for _, t := range transformations {
		// metricName, metricNamespace and metricValue are required members
		// whose shapes allow the empty string — the same presence test as
		// the pattern's.
		if !t.MetricNameSet {
			return errRequiredMember("metricTransformations.member.metricName")
		}
		if !t.MetricNamespaceSet {
			return errRequiredMember("metricTransformations.member.metricNamespace")
		}
		if !t.MetricValueSet {
			return errRequiredMember("metricTransformations.member.metricValue")
		}
		if err := validateMetricName(t.MetricName); err != nil {
			return err
		}
		if err := validateMetricNamespace(t.MetricNamespace); err != nil {
			return err
		}
		if err := validateMetricValue(t.MetricValue); err != nil {
			return err
		}
		// "The unit to assign to the metric. If you omit this, the unit
		// is set as None" — the empty string is the omitted form, any
		// other value must be on the enum.
		if t.Unit != "" && !validStandardUnits[t.Unit] {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("Invalid unit: %s. Valid values are the CloudWatch StandardUnit values (Seconds, Count, Percent, ...)", t.Unit), 400)
		}
		if err := validateTransformDimensions(t, emitSystemFieldDimensions); err != nil {
			return err
		}
	}
	if err := validateEmitSystemFieldDimensions(emitSystemFieldDimensions); err != nil {
		return err
	}
	if err := validateFieldSelectionCriteria(fieldSelectionCriteria); err != nil {
		return err
	}
	if err := validateFieldSelectionCriteriaSyntax(fieldSelectionCriteria); err != nil {
		return err
	}

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return err
	}

	if _, err = store.GetLogGroup(logGroupName); err != nil {
		return mapStoreError(err)
	}

	// The documented per-group regex quota spans both filter families,
	// so the census runs here — above either family's store. The filter
	// being replaced keeps its slot.
	if filterpattern.PatternContainsRegex(filterPattern) {
		if err := refuseRegexPatternQuotaExceeded(store, logGroupName, filterName, ""); err != nil {
			return err
		}
	}

	filter := logsstore.NewMetricFilter(logGroupName, filterName, filterPattern, transformations)
	filter.ApplyOnTransformedLogs = applyOnTransformedLogs
	filter.FieldSelectionCriteria = fieldSelectionCriteria
	filter.EmitSystemFieldDimensions = emitSystemFieldDimensions

	if err := store.PutMetricFilter(filter); err != nil {
		return mapStoreError(err)
	}
	return nil
}

// refuseRegexPatternQuotaExceeded enforces the documented per-group
// regex quota — "There is a maximum of 5 filter patterns containing
// regex for each log group when creating metric filters or subscription
// filters" (Filter and pattern syntax) — across both filter families,
// so the census lives above either family's store. The filter being
// created or updated is exempt in its own family: replacing a regex
// pattern keeps its slot.
func refuseRegexPatternQuotaExceeded(store *logsstore.Store, group, replacingMetricFilter, replacingSubscriptionFilter string) error {
	count := 0
	marker := ""
	for {
		filters, next, err := store.ListMetricFilters(group, "", marker, logsstore.ListingPageSize)
		if err != nil {
			return err
		}
		for _, mf := range filters {
			if mf.Name == replacingMetricFilter {
				continue
			}
			if filterpattern.PatternContainsRegex(mf.FilterPattern) {
				count++
			}
		}
		if next == "" {
			break
		}
		marker = next
	}
	subs, err := store.ListSubscriptionFilters(group, "")
	if err != nil {
		return err
	}
	for _, sf := range subs {
		if sf.FilterName == replacingSubscriptionFilter {
			continue
		}
		if filterpattern.PatternContainsRegex(sf.FilterPattern) {
			count++
		}
	}
	if count >= logsstore.MaxRegexFilterPatternsPerGroup {
		return NewLogsError("LimitExceededException",
			fmt.Sprintf("A log group can have a maximum of %d filter patterns containing regex", logsstore.MaxRegexFilterPatternsPerGroup), 400)
	}
	return nil
}

func parseMetricTransformationsFromMap(req *request.ParsedRequest) []logsstore.MetricTransformation {
	var transformations []logsstore.MetricTransformation

	var mt interface{}
	if v, ok := req.Parameters["metricTransformations"]; ok {
		mt = v
	} else if v, ok := req.Parameters["MetricTransformations"]; ok {
		mt = v
	}

	if mt != nil {
		if mtList, ok := mt.([]interface{}); ok {
			for _, item := range mtList {
				if m, ok := item.(map[string]interface{}); ok {
					t := logsstore.MetricTransformation{}
					if v, ok := m["metricName"].(string); ok {
						t.MetricName = v
						t.MetricNameSet = true
					} else if v, ok := m["MetricName"].(string); ok {
						t.MetricName = v
						t.MetricNameSet = true
					}
					if v, ok := m["metricNamespace"].(string); ok {
						t.MetricNamespace = v
						t.MetricNamespaceSet = true
					} else if v, ok := m["MetricNamespace"].(string); ok {
						t.MetricNamespace = v
						t.MetricNamespaceSet = true
					}
					if v, ok := m["metricValue"].(string); ok {
						t.MetricValue = v
						t.MetricValueSet = true
					} else if v, ok := m["MetricValue"].(string); ok {
						t.MetricValue = v
						t.MetricValueSet = true
					}
					if v, ok := m["defaultValue"].(float64); ok {
						t.DefaultValue = v
						t.DefaultValueSet = true
					} else if v, ok := m["DefaultValue"].(float64); ok {
						t.DefaultValue = v
						t.DefaultValueSet = true
					}
					if v, ok := m["unit"].(string); ok {
						t.Unit = v
					} else if v, ok := m["Unit"].(string); ok {
						t.Unit = v
					}
					if raw, ok := m["dimensions"]; ok {
						t.Dimensions = stringMapMember(raw)
					} else if raw, ok := m["Dimensions"]; ok {
						t.Dimensions = stringMapMember(raw)
					}
					transformations = append(transformations, t)
				}
			}
		}
	}

	return transformations
}

// stringMapMember coerces one map-valued request member (the dimensions
// member's JSON form) to its string map, dropping non-string entries.
func stringMapMember(raw interface{}) map[string]string {
	list, ok := raw.(map[string]interface{})
	if !ok {
		return nil
	}
	out := make(map[string]string, len(list))
	for k, v := range list {
		if s, ok := v.(string); ok {
			out[k] = s
		}
	}
	return out
}

// DeleteMetricFilter deletes the specified metric filter from the CloudWatch Logs log group.
func (s *LogsService) DeleteMetricFilter(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	filterName := request.GetParamLowerFirst(req.Parameters, "FilterName")

	if err := s.deleteMetricFilterCore(logGroupName, filterName, reqCtx.GetRegion()); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

func (s *LogsService) deleteMetricFilterCore(logGroupName, filterName, region string) error {
	if err := validateLogGroupName(logGroupName); err != nil {
		return err
	}
	if err := validateFilterName(filterName); err != nil {
		return err
	}

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return err
	}

	if err := store.DeleteMetricFilter(logGroupName, filterName); err != nil {
		return mapStoreError(err)
	}
	return nil
}

// DescribeMetricFiltersInput is the transport-agnostic input for the
// metric-filter listing.
type DescribeMetricFiltersInput struct {
	LogGroupName       string
	FilterNamePrefix   string
	MetricName         string
	MetricNameSet      bool
	MetricNamespace    string
	MetricNamespaceSet bool
	NextToken          string
	Limit              int32
	Region             string
}

// describeMetricFiltersCore lists metric filters under the store's
// prefix walk; the metricName/metricNamespace members ("Filters results
// to include only those with the specified metric name... you must also
// include the metricNamespace parameter", and vice versa) require each
// other and filter by the transformations' published metric.
func (s *LogsService) describeMetricFiltersCore(input DescribeMetricFiltersInput) ([]*logsstore.MetricFilter, string, error) {
	l, err := validateListLimit(input.Limit, logsstore.DefaultDescribeLimit, logsstore.MaxDescribeLimit)
	if err != nil {
		return nil, "", err
	}
	// The mutual-requirement rule both members document; the shapes allow
	// the empty string, so presence — not value — is the test.
	if input.MetricNameSet != input.MetricNamespaceSet {
		missing := "metricNamespace"
		if input.MetricNamespaceSet {
			missing = "metricName"
		}
		return nil, "", NewLogsError("InvalidParameterException",
			fmt.Sprintf("If you specify one of metricName and metricNamespace, you must specify both; %s is missing", missing), 400)
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return nil, "", err
	}

	// The operation declares ResourceNotFoundException: a listing against
	// a group that does not exist fails rather than serving an empty
	// page (the sibling stream listing carries the same check).
	if _, err := store.GetLogGroup(input.LogGroupName); err != nil {
		return nil, "", mapStoreError(err)
	}

	// The listing pages through the scoped token vocabulary: the group,
	// prefix and metric members digest into the scope, and the store's raw
	// key cursor never surfaces bare — a group-A token replayed against
	// group B rejects instead of silently skipping that group's keys.
	scope := listingScope("metricfilters", input.LogGroupName, input.FilterNamePrefix,
		input.MetricName+"\x1f"+input.MetricNamespace)

	if input.MetricNameSet {
		// The metric members filter by transformation content, which the
		// store's prefix walk knows nothing about: the core materialises
		// the group's matching filters (the per-group filter quota bounds
		// the walk) and pages the matched set by name — the name-ordered
		// walk makes the marker exact, and the page marker comes from the
		// pagination result like every other listing.
		var matched []*logsstore.MetricFilter
		scan := ""
		for {
			page, next, err := store.ListMetricFilters(input.LogGroupName, input.FilterNamePrefix, scan, logsstore.ListingPageSize)
			if err != nil {
				return nil, "", mapStoreError(err)
			}
			for _, f := range page {
				if metricFilterPublishesMetric(f, input.MetricName, input.MetricNamespace) {
					matched = append(matched, f)
				}
			}
			if next == "" {
				break
			}
			scan = next
		}
		result, err := paginateScopedListing(scope, input.NextToken, matched, int(l), func(f *logsstore.MetricFilter) string {
			return f.Name
		})
		if err != nil {
			return nil, "", err
		}
		return result.Items, result.NextMarker, nil
	}

	cursor := ""
	if input.NextToken != "" {
		var err error
		cursor, err = decodeListingToken(input.NextToken, scope)
		if err != nil {
			return nil, "", err
		}
	}
	filters, nextMarker, err := store.ListMetricFilters(input.LogGroupName, input.FilterNamePrefix, cursor, int(l))
	if err != nil {
		return nil, "", mapStoreError(err)
	}
	nextMarker, err = wrapScopedListingMarker(scope, nextMarker)
	if err != nil {
		return nil, "", err
	}
	return filters, nextMarker, nil
}

// metricFilterPublishesMetric reports whether one filter's
// transformations publish the named metric in the named namespace.
func metricFilterPublishesMetric(f *logsstore.MetricFilter, metricName, metricNamespace string) bool {
	for _, t := range f.MetricTransformations {
		if t.MetricName == metricName && t.MetricNamespace == metricNamespace {
			return true
		}
	}
	return false
}

// DescribeMetricFilters returns a list of metric filters for the specified CloudWatch Logs log group.
func (s *LogsService) DescribeMetricFilters(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	filterNamePrefix := request.GetParamLowerFirst(req.Parameters, "FilterNamePrefix")
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")
	limit := int32(request.GetIntParam(req.Parameters, "Limit"))

	filters, nextMarker, err := s.describeMetricFiltersCore(DescribeMetricFiltersInput{
		LogGroupName:       logGroupName,
		FilterNamePrefix:   filterNamePrefix,
		MetricName:         request.GetParamLowerFirst(req.Parameters, "MetricName"),
		MetricNameSet:      request.HasParamLowerFirst(req.Parameters, "MetricName"),
		MetricNamespace:    request.GetParamLowerFirst(req.Parameters, "MetricNamespace"),
		MetricNamespaceSet: request.HasParamLowerFirst(req.Parameters, "MetricNamespace"),
		NextToken:          nextToken,
		Limit:              limit,
		Region:             reqCtx.GetRegion(),
	})
	if err != nil {
		return nil, err
	}

	var metricFilters []map[string]interface{}
	for _, f := range filters {
		var transformations []map[string]interface{}
		for _, t := range f.MetricTransformations {
			transformation := map[string]interface{}{
				"metricName":      t.MetricName,
				"metricNamespace": t.MetricNamespace,
				"metricValue":     t.MetricValue,
			}
			if t.DefaultValueSet {
				transformation["defaultValue"] = t.DefaultValue
			}
			if len(t.Dimensions) > 0 {
				transformation["dimensions"] = t.Dimensions
			}
			if t.Unit != "" {
				transformation["unit"] = t.Unit
			}
			transformations = append(transformations, transformation)
		}

		entry := map[string]interface{}{
			"filterName":            f.Name,
			"logGroupName":          f.LogGroupName,
			"filterPattern":         f.FilterPattern,
			"metricTransformations": transformations,
			"creationTime":          f.CreatedAt.UnixMilli(),
		}
		if f.ApplyOnTransformedLogs {
			entry["applyOnTransformedLogs"] = f.ApplyOnTransformedLogs
		}
		if f.FieldSelectionCriteria != "" {
			entry["fieldSelectionCriteria"] = f.FieldSelectionCriteria
		}
		if len(f.EmitSystemFieldDimensions) > 0 {
			entry["emitSystemFieldDimensions"] = f.EmitSystemFieldDimensions
		}
		metricFilters = append(metricFilters, entry)
	}

	response := map[string]interface{}{
		"metricFilters": metricFilters,
	}
	if nextMarker != "" {
		response["nextToken"] = nextMarker
	}

	return response, nil
}

// TestMetricFilter tests a filter pattern against a set of log event messages.
func (s *LogsService) TestMetricFilter(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	result, err := s.testMetricFilterCore(testMetricFilterInput{
		FilterPattern:    request.GetParamLowerFirst(req.Parameters, "FilterPattern"),
		FilterPatternSet: request.HasParamLowerFirst(req.Parameters, "FilterPattern"),
		LogEventMessages: request.GetStringList(req.Parameters, "LogEventMessages"),
	})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"matches": result,
	}, nil
}

// testMetricFilterInput is the transport-agnostic input for TestMetricFilter.
type testMetricFilterInput struct {
	FilterPattern    string
	FilterPatternSet bool
	LogEventMessages []string
}

// testMetricFilterCore validates the members — the pattern's requiredness
// is a presence test (the shape allows the empty string, which matches
// every event, the same requiredness putMetricFilterCore implements) and
// the message list carries 1-50 entries (the logEventMessages member is
// required over the LogEventMessages target) — and matches every message
// against the pattern, reporting each match's eventNumber as the matched
// event's position in the input list.
func (s *LogsService) testMetricFilterCore(input testMetricFilterInput) ([]map[string]interface{}, error) {
	if !input.FilterPatternSet {
		return nil, errRequiredMember("filterPattern")
	}
	if err := validateFilterPattern(input.FilterPattern); err != nil {
		return nil, err
	}
	if len(input.LogEventMessages) == 0 {
		return nil, NewLogsError("InvalidParameterException",
			"logEventMessages must contain between 1 and 50 messages", 400)
	}
	if len(input.LogEventMessages) > logsstore.MaxTestMetricFilterLogEventMessages {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("logEventMessages must contain between 1 and %d messages", logsstore.MaxTestMetricFilterLogEventMessages), 400)
	}

	matcher := filterpattern.NewMatcher()
	var matches []map[string]interface{}
	for i, msg := range input.LogEventMessages {
		matched, extracted := matcher.ExtractMatches(input.FilterPattern, msg)
		if matched {
			evMap := make(map[string]interface{}, len(extracted))
			for k, v := range extracted {
				evMap[k] = v
			}
			matches = append(matches, map[string]interface{}{
				"eventMessage":    msg,
				"eventNumber":     int64(i + 1),
				"extractedValues": evMap,
			})
		}
	}
	return matches, nil
}

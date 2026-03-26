// Copyright 2026 Vorpalstacks Authors
// SPDX-License-Identifier: Apache-2.0

package cloudwatchlogs

import (
	"context"
	"strconv"

	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/services/aws/common/response"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	"vorpalstacks/pkg/filterpattern"
)

// PutMetricFilter creates or updates a metric filter for the specified CloudWatch Logs log group.
func (s *LogsService) PutMetricFilter(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	filterName := request.GetParamLowerFirst(req.Parameters, "FilterName")
	filterPattern := request.GetParamLowerFirst(req.Parameters, "FilterPattern")

	if logGroupName == "" || filterName == "" {
		return nil, ErrMissingParameter
	}

	var transformations []logsstore.MetricTransformation
	for i := 1; ; i++ {
		metricName := request.GetParamLowerFirst(req.Parameters, "MetricTransformations."+strconv.Itoa(i)+".MetricName")
		metricNamespace := request.GetParamLowerFirst(req.Parameters, "MetricTransformations."+strconv.Itoa(i)+".MetricNamespace")
		if metricName == "" && metricNamespace == "" {
			break
		}
		transformations = append(transformations, logsstore.MetricTransformation{
			MetricName:      metricName,
			MetricNamespace: metricNamespace,
			MetricValue:     request.GetParamLowerFirst(req.Parameters, "MetricTransformations."+strconv.Itoa(i)+".MetricValue"),
			DefaultValue:    request.GetFloatParam(req.Parameters, "MetricTransformations."+strconv.Itoa(i)+".DefaultValue"),
		})
	}

	if len(transformations) == 0 {
		transformations = parseMetricTransformationsFromMap(req)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	_, err = store.GetLogGroup(logGroupName)
	if err != nil {
		return nil, mapStoreError(err)
	}

	filter := logsstore.NewMetricFilter(logGroupName, filterName, filterPattern, transformations)
	if err := store.PutMetricFilter(filter); err != nil {
		return nil, mapStoreError(err)
	}

	return response.EmptyResponse(), nil
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
					} else if v, ok := m["MetricName"].(string); ok {
						t.MetricName = v
					}
					if v, ok := m["metricNamespace"].(string); ok {
						t.MetricNamespace = v
					} else if v, ok := m["MetricNamespace"].(string); ok {
						t.MetricNamespace = v
					}
					if v, ok := m["metricValue"].(string); ok {
						t.MetricValue = v
					} else if v, ok := m["MetricValue"].(string); ok {
						t.MetricValue = v
					}
					if v, ok := m["defaultValue"].(float64); ok {
						t.DefaultValue = v
						t.DefaultValueSet = true
					} else if v, ok := m["DefaultValue"].(float64); ok {
						t.DefaultValue = v
						t.DefaultValueSet = true
					}
					transformations = append(transformations, t)
				}
			}
		}
	}

	return transformations
}

// DeleteMetricFilter deletes the specified metric filter from the CloudWatch Logs log group.
func (s *LogsService) DeleteMetricFilter(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	filterName := request.GetParamLowerFirst(req.Parameters, "FilterName")

	if logGroupName == "" || filterName == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteMetricFilter(logGroupName, filterName); err != nil {
		return nil, mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// DescribeMetricFilters returns a list of metric filters for the specified CloudWatch Logs log group.
func (s *LogsService) DescribeMetricFilters(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	filterNamePrefix := request.GetParamLowerFirst(req.Parameters, "FilterNamePrefix")
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")
	limit := int32(request.GetIntParam(req.Parameters, "Limit"))
	if limit <= 0 {
		limit = 50
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	filters, nextMarker, err := store.ListMetricFilters(logGroupName, filterNamePrefix, nextToken, int(limit))
	if err != nil {
		return nil, mapStoreError(err)
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
			transformations = append(transformations, transformation)
		}

		metricFilters = append(metricFilters, map[string]interface{}{
			"filterName":            f.Name,
			"logGroupName":          f.LogGroupName,
			"filterPattern":         f.FilterPattern,
			"metricTransformations": transformations,
			"creationTime":          f.CreatedAt.UnixMilli(),
		})
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
	filterPattern := request.GetParamLowerFirst(req.Parameters, "FilterPattern")
	if filterPattern == "" {
		return nil, ErrInvalidParameter
	}

	var logEventMessages []string
	if v, ok := req.Parameters["logEventMessages"]; ok {
		if arr, ok := v.([]interface{}); ok {
			for _, item := range arr {
				if s, ok := item.(string); ok {
					logEventMessages = append(logEventMessages, s)
				}
			}
		}
	}

	if len(logEventMessages) > 10 {
		return nil, ErrInvalidParameter
	}

	matcher := filterpattern.NewMatcher()
	var matches []map[string]interface{}
	for _, msg := range logEventMessages {
		if matcher.Matches(filterPattern, msg) {
			matches = append(matches, map[string]interface{}{
				"message":  msg,
				"matched":  true,
				"metadata": map[string]interface{}{},
			})
		}
	}

	return map[string]interface{}{
		"matches": matches,
	}, nil
}

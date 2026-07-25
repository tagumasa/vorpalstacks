package cloudwatch

import (
	"context"
	"fmt"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
	"vorpalstacks/internal/store/aws/common"
)

// PutAnomalyDetector creates or updates a CloudWatch anomaly detection
// model for a metric. The model can be a single-metric detector
// (identified by Namespace, MetricName, Dimensions, Stat) or a
// metric-math detector (identified by MetricDataQueries).
//
// AWS docs: https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/API_PutAnomalyDetector.html
func (s *CloudWatchService) PutAnomalyDetector(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Check for SingleMetricAnomalyDetector parameter (newer API form).
	if smad, ok := req.Parameters["SingleMetricAnomalyDetector"]; ok {
		return s.putSingleMetricAnomalyDetector(store, smad, req.Parameters)
	}
	if smad, ok := req.Parameters["singleMetricAnomalyDetector"]; ok {
		return s.putSingleMetricAnomalyDetector(store, smad, req.Parameters)
	}

	// Check for MetricMathAnomalyDetector parameter.
	if mmad, ok := req.Parameters["MetricMathAnomalyDetector"]; ok {
		return s.putMetricMathAnomalyDetector(store, mmad)
	}
	if mmad, ok := req.Parameters["metricMathAnomalyDetector"]; ok {
		return s.putMetricMathAnomalyDetector(store, mmad)
	}

	// Legacy form: flat Namespace/MetricName/Dimensions/Stat parameters.
	namespace := getAlarmStringParam(req.Parameters, "Namespace", "namespace")
	metricName := getAlarmStringParam(req.Parameters, "MetricName", "metricName")
	stat := getAlarmStringParam(req.Parameters, "Stat", "stat")
	dimensions := parseAlarmDimensions(req.Parameters)

	if namespace == "" || metricName == "" || stat == "" {
		return nil, awserrors.NewMissingParameter(
			"Namespace, MetricName, and Stat are required for single-metric anomaly detectors")
	}

	detector := &cwstore.AnomalyDetector{
		Namespace:  namespace,
		MetricName: metricName,
		Stat:       stat,
		Dimensions: dimensions,
	}

	// Parse Configuration if present (top-level for legacy flat form).
	if cfg, ok := req.Parameters["Configuration"]; ok {
		detector.AnomalyDetectorConfiguration = parseAnomalyDetectorConfiguration(cfg)
	} else if cfg, ok := req.Parameters["configuration"]; ok {
		detector.AnomalyDetectorConfiguration = parseAnomalyDetectorConfiguration(cfg)
	}

	// Parse MetricCharacteristics if present (top-level for legacy flat form).
	if mc, ok := req.Parameters["MetricCharacteristics"]; ok {
		detector.MetricCharacteristics = parseMetricCharacteristics(mc)
	} else if mc, ok := req.Parameters["metricCharacteristics"]; ok {
		detector.MetricCharacteristics = parseMetricCharacteristics(mc)
	}

	saved, err := store.anomalyDetectors.PutAnomalyDetector(detector)
	if err != nil {
		return nil, fmt.Errorf("failed to put anomaly detector: %w", err)
	}

	return map[string]interface{}{
		"AnomalyDetectorId": saved.ID,
	}, nil
}

func (s *CloudWatchService) putSingleMetricAnomalyDetector(store *cloudwatchStores, raw interface{}, params map[string]interface{}) (interface{}, error) {
	smad, ok := raw.(map[string]interface{})
	if !ok {
		return nil, awserrors.NewInvalidParameterValueException(
			"SingleMetricAnomalyDetector must be an object")
	}

	namespace := getAlarmStringParam(smad, "Namespace", "namespace")
	metricName := getAlarmStringParam(smad, "MetricName", "metricName")
	stat := getAlarmStringParam(smad, "Stat", "stat")
	accountID := getAlarmStringParam(smad, "AccountId", "accountId")
	dimensions := parseAlarmDimensions(smad)

	if namespace == "" || metricName == "" || stat == "" {
		return nil, awserrors.NewMissingParameter(
			"Namespace, MetricName, and Stat are required in SingleMetricAnomalyDetector")
	}

	detector := &cwstore.AnomalyDetector{
		Namespace:  namespace,
		MetricName: metricName,
		Stat:       stat,
		Dimensions: dimensions,
		AccountID:  accountID,
	}

	// Parse Configuration if present.
	if cfg, ok := smad["Configuration"]; ok {
		detector.AnomalyDetectorConfiguration = parseAnomalyDetectorConfiguration(cfg)
	} else if cfg, ok := smad["configuration"]; ok {
		detector.AnomalyDetectorConfiguration = parseAnomalyDetectorConfiguration(cfg)
	}

	// Parse MetricCharacteristics if present (from top-level params).
	if mc, ok := params["MetricCharacteristics"]; ok {
		detector.MetricCharacteristics = parseMetricCharacteristics(mc)
	} else if mc, ok := params["metricCharacteristics"]; ok {
		detector.MetricCharacteristics = parseMetricCharacteristics(mc)
	}

	saved, err := store.anomalyDetectors.PutAnomalyDetector(detector)
	if err != nil {
		return nil, fmt.Errorf("failed to put anomaly detector: %w", err)
	}

	return map[string]interface{}{
		"AnomalyDetectorId": saved.ID,
	}, nil
}

func (s *CloudWatchService) putMetricMathAnomalyDetector(store *cloudwatchStores, raw interface{}) (interface{}, error) {
	mmad, ok := raw.(map[string]interface{})
	if !ok {
		return nil, awserrors.NewInvalidParameterValueException(
			"MetricMathAnomalyDetector must be an object")
	}

	queries := parseMetricDataQueries(mmad["MetricDataQueries"])
	if len(queries) == 0 {
		queries = parseMetricDataQueries(mmad["metricDataQueries"])
	}
	if len(queries) == 0 {
		return nil, awserrors.NewMissingParameter(
			"MetricDataQueries is required in MetricMathAnomalyDetector")
	}

	detector := &cwstore.AnomalyDetector{
		MetricDataQueries: queries,
	}

	if cfg, ok := mmad["Configuration"]; ok {
		detector.AnomalyDetectorConfiguration = parseAnomalyDetectorConfiguration(cfg)
	} else if cfg, ok := mmad["configuration"]; ok {
		detector.AnomalyDetectorConfiguration = parseAnomalyDetectorConfiguration(cfg)
	}

	saved, err := store.anomalyDetectors.PutMetricMathAnomalyDetector(detector)
	if err != nil {
		return nil, fmt.Errorf("failed to put metric math anomaly detector: %w", err)
	}

	return map[string]interface{}{
		"AnomalyDetectorId": saved.ID,
	}, nil
}

// DeleteAnomalyDetector deletes the specified anomaly detection model.
//
// AWS docs: https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/API_DeleteAnomalyDetector.html
func (s *CloudWatchService) DeleteAnomalyDetector(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Try AnomalyDetectorId first (for metric-math detectors).
	detectorID := getAlarmStringParam(req.Parameters, "AnomalyDetectorId", "anomalyDetectorId")

	// Legacy form: flat Namespace/MetricName/Dimensions/Stat.
	namespace := getAlarmStringParam(req.Parameters, "Namespace", "namespace")
	metricName := getAlarmStringParam(req.Parameters, "MetricName", "metricName")
	stat := getAlarmStringParam(req.Parameters, "Stat", "stat")
	dimensions := parseAlarmDimensions(req.Parameters)

	// Check SingleMetricAnomalyDetector parameter.
	if smad, ok := req.Parameters["SingleMetricAnomalyDetector"]; ok {
		if m, ok := smad.(map[string]interface{}); ok {
			if namespace == "" {
				namespace = getAlarmStringParam(m, "Namespace", "namespace")
			}
			if metricName == "" {
				metricName = getAlarmStringParam(m, "MetricName", "metricName")
			}
			if stat == "" {
				stat = getAlarmStringParam(m, "Stat", "stat")
			}
			if len(dimensions) == 0 {
				dimensions = parseAlarmDimensions(m)
			}
		}
	}

	if detectorID == "" && namespace == "" && metricName == "" && stat == "" {
		return nil, awserrors.NewMissingParameter(
			"Either AnomalyDetectorId or (Namespace, MetricName, Stat) must be specified")
	}

	if err := store.anomalyDetectors.DeleteAnomalyDetector(namespace, metricName, dimensions, stat, detectorID); err != nil {
		return nil, awserrors.NewResourceNotFoundException("AnomalyDetector", err.Error())
	}

	return map[string]interface{}{}, nil
}

// DescribeAnomalyDetectors lists the anomaly detection models in the
// account. Results can be filtered by namespace, metric name, dimensions,
// anomaly detector type, or anomaly detector IDs.
//
// AWS docs: https://docs.aws.amazon.com/AmazonCloudWatch/latest/APIReference/API_DescribeAnomalyDetectors.html
func (s *CloudWatchService) DescribeAnomalyDetectors(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	namespace := getAlarmStringParam(req.Parameters, "Namespace", "namespace")
	metricName := getAlarmStringParam(req.Parameters, "MetricName", "metricName")

	// Smithy input is AnomalyDetectorTypes (plural list), not
	// AnomalyDetectorType (singular string). The SDK sends
	// AnomalyDetectorTypes.1=…, AnomalyDetectorTypes.2=… etc.
	var detectorTypes []string
	if typesRaw, ok := req.Parameters["AnomalyDetectorTypes"]; ok {
		detectorTypes = toStringSlice(typesRaw)
	} else if typesRaw, ok := req.Parameters["anomalyDetectorTypes"]; ok {
		detectorTypes = toStringSlice(typesRaw)
	}
	// Validate detector types if provided.
	for _, dt := range detectorTypes {
		if dt != cwstore.AnomalyDetectorTypeSingleMetric && dt != cwstore.AnomalyDetectorTypeMetricMath {
			return nil, awserrors.NewInvalidParameterValueException(
				fmt.Sprintf("AnomalyDetectorType must be %s or %s",
					cwstore.AnomalyDetectorTypeSingleMetric, cwstore.AnomalyDetectorTypeMetricMath))
		}
	}

	dimensions := parseAlarmDimensions(req.Parameters)

	var detectorIDs []string
	if ids, ok := req.Parameters["AnomalyDetectorIds"]; ok {
		detectorIDs = toStringSlice(ids)
	} else if ids, ok := req.Parameters["anomalyDetectorIds"]; ok {
		detectorIDs = toStringSlice(ids)
	}

	marker := pagination.GetMarker(req.Parameters, "NextToken")
	maxResults := pagination.GetMaxItems(req.Parameters, 100, "MaxResults")

	opts := cwstore.DescribeAnomalyDetectorsOpts{
		AnomalyDetectorIDs:   detectorIDs,
		AnomalyDetectorTypes: detectorTypes,
		Namespace:            namespace,
		MetricName:           metricName,
		Dimensions:           dimensions,
		ListOpts:             common.ListOptions{Marker: marker, MaxItems: maxResults},
	}

	detectors, nextToken, err := store.anomalyDetectors.DescribeAnomalyDetectors(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to describe anomaly detectors: %w", err)
	}

	results := make([]map[string]interface{}, 0, len(detectors))
	for _, d := range detectors {
		results = append(results, anomalyDetectorToResponse(d))
	}

	resp := map[string]interface{}{
		"AnomalyDetectors": results,
	}
	if nextToken != "" {
		resp["NextToken"] = nextToken
	}
	return resp, nil
}

// anomalyDetectorToResponse serialises an AnomalyDetector into the
// AWS API response format.  Both the deprecated flat fields and the
// current nested objects (SingleMetricAnomalyDetector /
// MetricMathAnomalyDetector) are emitted for backwards compatibility.
func anomalyDetectorToResponse(d *cwstore.AnomalyDetector) map[string]interface{} {
	resp := map[string]interface{}{
		"StateValue": d.State,
	}

	if d.ID != "" {
		resp["AnomalyDetectorId"] = d.ID
	}

	// Build the SingleMetricAnomalyDetector nested object for
	// single-metric detectors.  The flat fields below are deprecated
	// in the Smithy model but still emitted for older SDK clients.
	if d.AnomalyDetectorType != cwstore.AnomalyDetectorTypeMetricMath {
		single := map[string]interface{}{}
		if d.Namespace != "" {
			single["Namespace"] = d.Namespace
			resp["Namespace"] = d.Namespace
		}
		if d.MetricName != "" {
			single["MetricName"] = d.MetricName
			resp["MetricName"] = d.MetricName
		}
		if d.Stat != "" {
			single["Stat"] = d.Stat
			resp["Stat"] = d.Stat
		}
		if len(d.Dimensions) > 0 {
			dims := dimensionsToResponse(d.Dimensions)
			single["Dimensions"] = dims
			resp["Dimensions"] = dims
		}
		if d.AccountID != "" {
			single["AccountId"] = d.AccountID
			resp["AccountId"] = d.AccountID
		}
		resp["SingleMetricAnomalyDetector"] = single
	}

	if d.AnomalyDetectorConfiguration != nil {
		cfg := map[string]interface{}{}
		if len(d.AnomalyDetectorConfiguration.ExcludedTimeRanges) > 0 {
			ranges := make([]map[string]interface{}, len(d.AnomalyDetectorConfiguration.ExcludedTimeRanges))
			for i, r := range d.AnomalyDetectorConfiguration.ExcludedTimeRanges {
				ranges[i] = map[string]interface{}{
					"StartTime": r.StartTime.Format("2006-01-02T15:04:05Z"),
					"EndTime":   r.EndTime.Format("2006-01-02T15:04:05Z"),
				}
			}
			cfg["ExcludedTimeRanges"] = ranges
		}
		if d.AnomalyDetectorConfiguration.MetricTimezone != "" {
			cfg["MetricTimezone"] = d.AnomalyDetectorConfiguration.MetricTimezone
		}
		resp["Configuration"] = cfg
	}
	if d.MetricCharacteristics != nil {
		resp["MetricCharacteristics"] = map[string]interface{}{
			"PeriodicSpikes": d.MetricCharacteristics.PeriodicSpikes,
		}
	}
	if len(d.MetricDataQueries) > 0 {
		resp["MetricMathAnomalyDetector"] = map[string]interface{}{
			"MetricDataQueries": metricDataQueriesToResponse(d.MetricDataQueries),
		}
	}

	return resp
}

// parseAnomalyDetectorConfiguration extracts the configuration from a
// raw map.
func parseAnomalyDetectorConfiguration(raw interface{}) *cwstore.AnomalyDetectorConfiguration {
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil
	}

	cfg := &cwstore.AnomalyDetectorConfiguration{}

	if etr, ok := m["ExcludedTimeRanges"]; ok {
		if ranges, ok := etr.([]interface{}); ok {
			for _, r := range ranges {
				rm, ok := r.(map[string]interface{})
				if !ok {
					continue
				}
				timeRange := cwstore.Range{}
				if start, ok := rm["StartTime"]; ok {
					timeRange.StartTime = parseTime(start)
				}
				if end, ok := rm["EndTime"]; ok {
					timeRange.EndTime = parseTime(end)
				}
				cfg.ExcludedTimeRanges = append(cfg.ExcludedTimeRanges, timeRange)
			}
		}
	}

	if tz, ok := m["MetricTimezone"]; ok {
		if tzStr, ok := tz.(string); ok {
			cfg.MetricTimezone = tzStr
		}
	}

	return cfg
}

// parseMetricCharacteristics extracts metric characteristics from a raw map.
func parseMetricCharacteristics(raw interface{}) *cwstore.MetricCharacteristics {
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil
	}
	mc := &cwstore.MetricCharacteristics{}
	if ps, ok := m["PeriodicSpikes"]; ok {
		if b, ok := ps.(bool); ok {
			mc.PeriodicSpikes = b
		}
	}
	return mc
}

// dimensionsToResponse converts a slice of Dimension structs to the
// AWS API response format.
func dimensionsToResponse(dims []cwstore.Dimension) []map[string]interface{} {
	result := make([]map[string]interface{}, len(dims))
	for i, d := range dims {
		result[i] = map[string]interface{}{
			"Name":  d.Name,
			"Value": d.Value,
		}
	}
	return result
}

// toStringSlice converts an interface{} (expected to be []interface{})
// to a []string.
func toStringSlice(raw interface{}) []string {
	arr, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	result := make([]string, 0, len(arr))
	for _, v := range arr {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

// parseTime parses a time value from an RFC3339 string.
func parseTime(raw interface{}) time.Time {
	s, ok := raw.(string)
	if !ok {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

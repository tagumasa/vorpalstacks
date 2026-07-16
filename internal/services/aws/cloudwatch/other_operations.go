package cloudwatch

import (
	"context"
	"fmt"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
)

// GetOTelEnrichment returns the current status of OTel enrichment.
func (s *CloudWatchService) GetOTelEnrichment(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	s.globalMu.Lock()
	status := s.otelStatus
	s.globalMu.Unlock()

	return map[string]interface{}{
		"Status": status,
	}, nil
}

// StartOTelEnrichment starts the OTel enrichment process.
func (s *CloudWatchService) StartOTelEnrichment(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	s.globalMu.Lock()
	s.otelStatus = "STARTED"
	s.globalMu.Unlock()

	return map[string]interface{}{}, nil
}

// StopOTelEnrichment stops the OTel enrichment process.
func (s *CloudWatchService) StopOTelEnrichment(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	s.globalMu.Lock()
	s.otelStatus = "STOPPED"
	s.globalMu.Unlock()

	return map[string]interface{}{}, nil
}

// --- Datasets ---

// GetDataset returns information about a dataset, including any
// associated KMS key.
func (s *CloudWatchService) GetDataset(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	identifier := getAlarmStringParam(req.Parameters, "DatasetIdentifier", "datasetIdentifier")
	if identifier == "" {
		identifier = "default"
	}

	s.globalMu.Lock()
	kmsKey := s.datasetKMS[identifier]
	s.globalMu.Unlock()

	return map[string]interface{}{
		"DatasetId": identifier,
		"Arn":       fmt.Sprintf("arn:aws:cloudwatch:::dataset/%s", identifier),
		"KmsKeyArn": kmsKey,
	}, nil
}

// AssociateDatasetKmsKey associates a KMS key with a dataset.
func (s *CloudWatchService) AssociateDatasetKmsKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	identifier := getAlarmStringParam(req.Parameters, "DatasetIdentifier", "datasetIdentifier")
	kmsKeyArn := getAlarmStringParam(req.Parameters, "KmsKeyArn", "kmsKeyArn")

	if identifier == "" || kmsKeyArn == "" {
		return nil, awserrors.NewMissingParameter("DatasetIdentifier and KmsKeyArn are required")
	}

	s.globalMu.Lock()
	s.datasetKMS[identifier] = kmsKeyArn
	s.globalMu.Unlock()

	return map[string]interface{}{}, nil
}

// DisassociateDatasetKmsKey removes the KMS key association from a
// dataset.
func (s *CloudWatchService) DisassociateDatasetKmsKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	identifier := getAlarmStringParam(req.Parameters, "DatasetIdentifier", "datasetIdentifier")
	if identifier == "" {
		return nil, awserrors.NewMissingParameter("DatasetIdentifier is required")
	}

	s.globalMu.Lock()
	delete(s.datasetKMS, identifier)
	s.globalMu.Unlock()

	return map[string]interface{}{}, nil
}

// --- Other operations ---

// alarmContributor represents a single contributor to an anomaly
// detection alarm evaluation.
type alarmContributor struct {
	Timestamp   string  `json:"timestamp"`
	MetricValue float64 `json:"metricValue"`
	BandUpper   float64 `json:"bandUpper"`
	BandLower   float64 `json:"bandLower"`
}

// DescribeAlarmContributors returns the top contributors for an alarm
// that uses anomaly detection. It resolves the alarm's
// ANOMALY_DETECTION_BAND expression, computes the EWMA band from the
// referenced metric's historical data, and returns the data points
// that fall outside the band as contributors.
func (s *CloudWatchService) DescribeAlarmContributors(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	alarmName := getAlarmStringParam(req.Parameters, "AlarmName", "alarmName")
	if alarmName == "" {
		return nil, awserrors.NewMissingParameter("AlarmName is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	alarm, err := store.alarms.GetAlarm(alarmName)
	if err != nil || alarm == nil {
		return map[string]interface{}{
			"AlarmName":        alarmName,
			"Contributors":     []interface{}{},
			"RoughResultCount": 0,
		}, nil
	}

	// Only anomaly detection alarms have contributors.
	if len(alarm.Metrics) == 0 || !hasAnomalyDetectionBand(alarm.Metrics) {
		return map[string]interface{}{
			"AlarmName":        alarmName,
			"Contributors":     []interface{}{},
			"RoughResultCount": 0,
		}, nil
	}

	contributors, err := computeAlarmContributors(alarm, store.metrics)
	if err != nil {
		return map[string]interface{}{
			"AlarmName":        alarmName,
			"Contributors":     []interface{}{},
			"RoughResultCount": 0,
		}, nil
	}

	return map[string]interface{}{
		"AlarmName":        alarmName,
		"Contributors":     contributors,
		"RoughResultCount": len(contributors),
	}, nil
}

// PutLogAlarm creates or updates a log-based alarm. Log alarms evaluate
// CloudWatch Logs query results against a threshold.
func (s *CloudWatchService) PutLogAlarm(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	alarmName := getAlarmStringParam(req.Parameters, "AlarmName", "alarmName")
	if alarmName == "" {
		return nil, awserrors.NewMissingParameter("AlarmName is required")
	}

	description := getAlarmStringParam(req.Parameters, "AlarmDescription", "alarmDescription")
	comparisonOperator := getAlarmStringParam(req.Parameters, "ComparisonOperator", "comparisonOperator")
	threshold := getAlarmFloatParam(req.Parameters, "Threshold", "threshold")
	treatMissingData := getAlarmStringParam(req.Parameters, "TreatMissingData", "treatMissingData")
	actionsEnabled := getAlarmBoolParam(req.Parameters, []string{"ActionsEnabled", "actionsEnabled"}, true)

	// Parse scheduled query configuration.
	sqc := ""
	if v, ok := req.Parameters["ScheduledQueryConfiguration"]; ok {
		if m, ok := v.(map[string]interface{}); ok {
			if sq, ok := m["ScheduledQueryName"]; ok {
				if sqStr, ok := sq.(string); ok {
					sqc = sqStr
				}
			}
		}
	}

	// Store as a metric alarm variant with the log alarm metadata.
	alarm := cwstore.NewAlarm(alarmName, "AWS/Logs", "LogQueryResult")
	alarm.AlarmDescription = description
	alarm.ComparisonOperator = comparisonOperator
	alarm.Threshold = threshold
	alarm.TreatMissingData = treatMissingData
	alarm.ActionsEnabled = actionsEnabled
	alarm.AlarmActions = parseStringArrayParam(req.Parameters, "AlarmActions", "alarmActions")
	alarm.OKActions = parseStringArrayParam(req.Parameters, "OKActions", "okActions")
	alarm.InsufficientDataActions = parseStringArrayParam(req.Parameters, "InsufficientDataActions", "insufficientDataActions")

	// Store the scheduled query name in the AlarmRule field.
	if sqc != "" {
		alarm.AlarmRule = sqc
	}

	tags, tagErr := parseAndValidateAlarmTags(req.Parameters)
	if tagErr != nil {
		return nil, tagErr
	}
	alarm.Tags = tags

	result, err := s.upsertAlarm(store.alarms, alarm, cwstore.AlarmTypeMetricAlarm)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"AlarmArn": result.ARN,
	}, nil
}

// computeAlarmContributors computes the anomaly band for the given
// alarm and returns data points that fall outside the band as
// contributors.
func computeAlarmContributors(alarm *cwstore.Alarm, metricStore *cwstore.MetricChunkStore) ([]alarmContributor, error) {
	ctx := prepareAnomalyBand(alarm, metricStore)
	if ctx == nil {
		return nil, fmt.Errorf("no anomaly detection band or metric data available")
	}

	var contributors []alarmContributor
	for i, s := range ctx.stats {
		if s.Timestamp.Before(ctx.startTime) || s.Timestamp.After(ctx.endTime) {
			continue
		}
		val := statisticValue(s, ctx.statLower)
		if val < ctx.lower[i] || val > ctx.upper[i] {
			contributors = append(contributors, alarmContributor{
				Timestamp:   s.Timestamp.Format(time.RFC3339),
				MetricValue: val,
				BandUpper:   ctx.upper[i],
				BandLower:   ctx.lower[i],
			})
		}
	}

	return contributors, nil
}

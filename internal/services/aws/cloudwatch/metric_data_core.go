package cloudwatch

import (
	"encoding/json"
	"fmt"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
)

// GetMetricWidgetImageInput holds parameters for GetMetricWidgetImage.
type GetMetricWidgetImageInput struct {
	MetricWidget string
	OutputFormat string
}

// getMetricWidgetImageCore validates input and retrieves the metric
// widget image. Returns the parsed widget definition, validated format,
// and any validation error.
func (s *CloudWatchService) getMetricWidgetImageCore(input *GetMetricWidgetImageInput) (*widgetDef, string, error) {
	if input.MetricWidget == "" {
		return nil, "", ErrInvalidParameter
	}
	if !json.Valid([]byte(input.MetricWidget)) {
		return nil, "", ErrInvalidParameter
	}

	format := input.OutputFormat
	if format == "" {
		format = "png"
	}
	if err := validateOutputFormat(format); err != nil {
		return nil, "", err
	}

	def, err := parseWidgetDefinition(input.MetricWidget)
	if err != nil {
		return nil, "", ErrInvalidParameter
	}

	return def, format, nil
}

// GetMetricDataInput holds validated parameters for GetMetricData.
type GetMetricDataInput struct {
	StartTime         time.Time
	EndTime           time.Time
	MetricDataQueries []cwstore.MetricDataQuery
	ScanBy            string
}

// getMetricDataCore validates GetMetricData input and returns the metric
// store for the complex multi-query evaluation that remains in the handler.
func (s *CloudWatchService) getMetricDataCore(stores *cloudwatchStores, input *GetMetricDataInput) (*cwstore.MetricChunkStore, error) {
	if input.StartTime.IsZero() || input.EndTime.IsZero() {
		return nil, awserrors.NewMissingParameter("StartTime and EndTime are required")
	}
	if !input.StartTime.Before(input.EndTime) {
		return nil, awserrors.NewInvalidParameterValueException("StartTime must be before EndTime")
	}
	if len(input.MetricDataQueries) == 0 {
		return nil, awserrors.NewMissingParameter("MetricDataQueries is required")
	}
	if input.ScanBy != "" && input.ScanBy != "TimestampAscending" && input.ScanBy != "TimestampDescending" {
		return nil, awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Invalid ScanBy: %s", input.ScanBy))
	}
	return stores.metrics, nil
}

// describeAlarmContributorsCore validates input and retrieves the alarm
// for contributor analysis. Returns the alarm and metric store.
func (s *CloudWatchService) describeAlarmContributorsCore(stores *cloudwatchStores, alarmName string) (*cwstore.Alarm, *cwstore.MetricChunkStore, error) {
	if alarmName == "" {
		return nil, nil, awserrors.NewMissingParameter("AlarmName is required")
	}

	alarm, err := stores.alarms.GetAlarm(alarmName)
	if err != nil || alarm == nil {
		return nil, nil, awserrors.NewResourceNotFoundException("Alarm", alarmName)
	}

	if len(alarm.Metrics) == 0 || !hasAnomalyDetectionBand(alarm.Metrics) {
		return nil, nil, awserrors.NewInvalidParameterValueException(
			"The specified alarm does not use anomaly detection")
	}

	return alarm, stores.metrics, nil
}

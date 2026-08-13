package cloudwatch

import (
	"fmt"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
)

// PutMetricDataInput holds parameters for PutMetricData.
type PutMetricDataInput struct {
	Namespace  string
	MetricData []cwstore.MetricDatum
}

// GetMetricStatisticsInput holds parameters for GetMetricStatistics.
type GetMetricStatisticsInput struct {
	Namespace          string
	MetricName         string
	StartTime          time.Time
	EndTime            time.Time
	Period             int32
	Statistics         []string
	ExtendedStatistics []string
	Dimensions         []cwstore.Dimension
}

// putMetricDataCore validates input and stores metric data points.
func (s *CloudWatchService) putMetricDataCore(stores *cloudwatchStores, input *PutMetricDataInput) error {
	if err := validateNamespace(input.Namespace); err != nil {
		return err
	}
	if len(input.MetricData) == 0 {
		return ErrInvalidParameter
	}
	if len(input.MetricData) > maxMetricDataPerRequest {
		return awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Number of MetricData entries must not exceed %d", maxMetricDataPerRequest))
	}
	for _, datum := range input.MetricData {
		if err := validateMetricDatum(datum); err != nil {
			return err
		}
		if err := validateMetricName(datum.MetricName); err != nil {
			return err
		}
		if datum.StorageResolution != 0 {
			if err := validateStorageResolution(datum.StorageResolution); err != nil {
				return err
			}
		}
		if !datum.Timestamp.IsZero() {
			now := time.Now()
			if datum.Timestamp.Before(now.Add(-14 * 24 * time.Hour)) {
				return awserrors.NewInvalidParameterValueException(
					"A MetricData timestamp must not be more than 14 days in the past")
			}
			if datum.Timestamp.After(now.Add(2 * time.Hour)) {
				return awserrors.NewInvalidParameterValueException(
					"A MetricData timestamp must not be more than 2 hours in the future")
			}
		}
	}

	return stores.metrics.PutMetricData(input.Namespace, input.MetricData)
}

// getMetricStatisticsCore validates input and retrieves metric statistics.
func (s *CloudWatchService) getMetricStatisticsCore(stores *cloudwatchStores, input *GetMetricStatisticsInput) ([]*cwstore.MetricStatistics, error) {
	if err := validateNamespace(input.Namespace); err != nil {
		return nil, err
	}
	if err := validateMetricName(input.MetricName); err != nil {
		return nil, err
	}
	if input.StartTime.IsZero() || input.EndTime.IsZero() {
		return nil, ErrInvalidParameter
	}
	if !input.StartTime.Before(input.EndTime) {
		return nil, awserrors.NewInvalidParameterValueException(
			"StartTime must be before EndTime")
	}
	if input.Period <= 0 {
		return nil, awserrors.NewInvalidParameterValueException(
			"Period must be a positive integer")
	}

	for _, stat := range input.Statistics {
		if !validStatistics[stat] {
			return nil, awserrors.NewInvalidParameterValueException(
				fmt.Sprintf("Invalid Statistics value: %s. Must be one of SampleCount, Average, Sum, Minimum, Maximum", stat))
		}
	}
	for _, est := range input.ExtendedStatistics {
		if err := validateExtendedStatistic(est); err != nil {
			return nil, err
		}
	}

	if len(input.Statistics) == 0 && len(input.ExtendedStatistics) == 0 {
		return nil, awserrors.NewMissingParameter(
			"Must specify either Statistics or ExtendedStatistics")
	}

	if len(input.Dimensions) > maxDimensions {
		return nil, awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("Number of Dimensions must not exceed %d", maxDimensions))
	}

	query := cwstore.MetricQuery{
		Namespace:          input.Namespace,
		MetricName:         input.MetricName,
		Dimensions:         input.Dimensions,
		StartTime:          input.StartTime,
		EndTime:            input.EndTime,
		Period:             input.Period,
		Statistics:         input.Statistics,
		ExtendedStatistics: input.ExtendedStatistics,
	}

	return stores.metrics.GetMetricStatistics(query)
}

package cloudwatch

import (
	"errors"
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
	"vorpalstacks/internal/store/aws/common"
)

// PutAnomalyDetectorInput holds parameters for PutAnomalyDetector.
type PutAnomalyDetectorInput struct {
	Namespace  string
	MetricName string
	Stat       string
	Dimensions []cwstore.Dimension
	AccountID  string
	Detector   *cwstore.AnomalyDetector
}

// DeleteAnomalyDetectorInput holds parameters for DeleteAnomalyDetector.
type DeleteAnomalyDetectorInput struct {
	DetectorID string
	Namespace  string
	MetricName string
	Stat       string
	Dimensions []cwstore.Dimension
}

// DescribeAnomalyDetectorsInput holds parameters for DescribeAnomalyDetectors.
type DescribeAnomalyDetectorsInput struct {
	Namespace     string
	MetricName    string
	Dimensions    []cwstore.Dimension
	DetectorTypes []string
	DetectorIDs   []string
	NextToken     string
	MaxResults    int
}

// validateAnomalyDetectorStat validates the Stat field of an anomaly
// detector against the Smithy Statistic/ExtendedStatistic patterns.
func validateAnomalyDetectorStat(stat string) error {
	if stat == "" {
		return awserrors.NewMissingParameter("Stat is required")
	}
	if validStatistics[stat] {
		return nil
	}
	return validateExtendedStatistic(stat)
}

// putAnomalyDetectorCore validates input and stores an anomaly detector.
func (s *CloudWatchService) putAnomalyDetectorCore(stores *cloudwatchStores, input *PutAnomalyDetectorInput) (*cwstore.AnomalyDetector, error) {
	if err := validateNamespace(input.Namespace); err != nil {
		return nil, err
	}
	if err := validateMetricName(input.MetricName); err != nil {
		return nil, err
	}
	if err := validateAnomalyDetectorStat(input.Stat); err != nil {
		return nil, err
	}

	detector := input.Detector
	if detector == nil {
		detector = &cwstore.AnomalyDetector{
			Namespace:  input.Namespace,
			MetricName: input.MetricName,
			Stat:       input.Stat,
			Dimensions: input.Dimensions,
			AccountID:  input.AccountID,
		}
	}

	saved, err := stores.anomalyDetectors.PutAnomalyDetector(detector)
	if err != nil {
		return nil, fmt.Errorf("failed to put anomaly detector: %w", err)
	}
	return saved, nil
}

// putMetricMathAnomalyDetectorCore validates the MetricMathAnomalyDetector
// form of PutAnomalyDetector and stores the detector. The raw member map is
// carried as-is so member parsing keeps its original order. Returns the
// saved detector with its assigned ID.
func (s *CloudWatchService) putMetricMathAnomalyDetectorCore(stores *cloudwatchStores, raw interface{}) (*cwstore.AnomalyDetector, error) {
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

	saved, err := stores.anomalyDetectors.PutMetricMathAnomalyDetector(detector)
	if err != nil {
		return nil, fmt.Errorf("failed to put metric math anomaly detector: %w", err)
	}
	return saved, nil
}

// deleteAnomalyDetectorCore validates input and deletes an anomaly
// detector. Returns a descriptive identifier for error reporting.
func (s *CloudWatchService) deleteAnomalyDetectorCore(stores *cloudwatchStores, input *DeleteAnomalyDetectorInput) error {
	if input.DetectorID == "" && input.Namespace == "" && input.MetricName == "" && input.Stat == "" {
		return awserrors.NewMissingParameter(
			"Either AnomalyDetectorId or (Namespace, MetricName, Stat) must be specified")
	}

	if err := stores.anomalyDetectors.DeleteAnomalyDetector(
		input.Namespace, input.MetricName, input.Dimensions, input.Stat, input.DetectorID); err != nil {

		if errors.Is(err, cwstore.ErrResourceNotFound) {
			identifier := input.DetectorID
			if identifier == "" {
				identifier = fmt.Sprintf("%s/%s", input.Namespace, input.MetricName)
			}
			return awserrors.NewResourceNotFoundException("AnomalyDetector", identifier)
		}
		return fmt.Errorf("failed to delete anomaly detector: %w", err)
	}
	return nil
}

// describeAnomalyDetectorsCore validates input and lists anomaly detectors.
func (s *CloudWatchService) describeAnomalyDetectorsCore(stores *cloudwatchStores, input *DescribeAnomalyDetectorsInput) ([]*cwstore.AnomalyDetector, string, error) {
	for _, dt := range input.DetectorTypes {
		if dt != cwstore.AnomalyDetectorTypeSingleMetric && dt != cwstore.AnomalyDetectorTypeMetricMath {
			return nil, "", awserrors.NewInvalidParameterValueException(
				fmt.Sprintf("AnomalyDetectorType must be %s or %s",
					cwstore.AnomalyDetectorTypeSingleMetric, cwstore.AnomalyDetectorTypeMetricMath))
		}
	}

	opts := cwstore.DescribeAnomalyDetectorsOpts{
		AnomalyDetectorIDs:   input.DetectorIDs,
		AnomalyDetectorTypes: input.DetectorTypes,
		Namespace:            input.Namespace,
		MetricName:           input.MetricName,
		Dimensions:           input.Dimensions,
		ListOpts:             common.ListOptions{Marker: input.NextToken, MaxItems: input.MaxResults},
	}

	detectors, nextToken, err := stores.anomalyDetectors.DescribeAnomalyDetectors(opts)
	if err != nil {
		return nil, "", fmt.Errorf("failed to describe anomaly detectors: %w", err)
	}
	return detectors, nextToken, nil
}

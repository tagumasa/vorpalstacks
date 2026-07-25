package cloudwatch

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

func anomalyDetectorBucketName(region string) string {
	return "cw_anomaly_detectors-" + region
}

// AnomalyDetectorStore provides storage operations for CloudWatch
// anomaly detectors.
type AnomalyDetectorStore struct {
	*common.BaseStore
	arnBuilder *svcarn.ARNBuilder
	accountID  string
	region     string
	mu         sync.Mutex
}

// NewAnomalyDetectorStore creates a new AnomalyDetectorStore instance.
func NewAnomalyDetectorStore(store storage.BasicStorage, accountID, region string) *AnomalyDetectorStore {
	return &AnomalyDetectorStore{
		BaseStore:  common.NewBaseStore(store.Bucket(anomalyDetectorBucketName(region)), "cloudwatch-anomaly-detectors"),
		arnBuilder: svcarn.NewARNBuilder(accountID, region),
		accountID:  accountID,
		region:     region,
	}
}

func (s *AnomalyDetectorStore) buildArn(id string) string {
	return s.arnBuilder.Build("cloudwatch", "anomaly-detector:"+id)
}

// PutAnomalyDetector creates or updates a single-metric anomaly detector.
// The detector is keyed on the combination of Namespace, MetricName,
// sorted Dimensions, and Stat.
func (s *AnomalyDetectorStore) PutAnomalyDetector(detector *AnomalyDetector) (*AnomalyDetector, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate a deterministic ID from the metric identity.
	key, err := singleMetricDetectorKey(detector.Namespace, detector.MetricName, detector.Dimensions, detector.Stat)
	if err != nil {
		return nil, err
	}

	existing := &AnomalyDetector{}
	if err := s.BaseStore.Get(key, existing); err == nil {
		existing.MetricCharacteristics = detector.MetricCharacteristics
		existing.AnomalyDetectorConfiguration = detector.AnomalyDetectorConfiguration
		existing.UpdatedAt = time.Now().UTC()
		if err := s.BaseStore.Put(key, existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	now := time.Now().UTC()
	detector.ID = key
	detector.ARN = s.buildArn(key)
	detector.State = "PENDING_TRAINING"
	detector.CreatedAt = now
	detector.UpdatedAt = now
	if detector.AnomalyDetectorType == "" {
		detector.AnomalyDetectorType = AnomalyDetectorTypeSingleMetric
	}
	if err := s.BaseStore.Put(key, detector); err != nil {
		return nil, err
	}
	return detector, nil
}

// PutMetricMathAnomalyDetector creates or updates a metric-math anomaly
// detector.
func (s *AnomalyDetectorStore) PutMetricMathAnomalyDetector(detector *AnomalyDetector) (*AnomalyDetector, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key, err := metricMathDetectorKey(detector.MetricDataQueries)
	if err != nil {
		return nil, err
	}

	existing := &AnomalyDetector{}
	if err := s.BaseStore.Get(key, existing); err == nil {
		existing.MetricDataQueries = detector.MetricDataQueries
		existing.AnomalyDetectorConfiguration = detector.AnomalyDetectorConfiguration
		existing.UpdatedAt = time.Now().UTC()
		if err := s.BaseStore.Put(key, existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	now := time.Now().UTC()
	detector.ID = key
	detector.ARN = s.buildArn(key)
	detector.State = "PENDING_TRAINING"
	detector.CreatedAt = now
	detector.UpdatedAt = now
	detector.AnomalyDetectorType = AnomalyDetectorTypeMetricMath
	if err := s.BaseStore.Put(key, detector); err != nil {
		return nil, err
	}
	return detector, nil
}

// DeleteAnomalyDetector deletes a single-metric anomaly detector by its
// metric identity (namespace, metric name, dimensions, stat).
func (s *AnomalyDetectorStore) DeleteAnomalyDetector(namespace, metricName string, dimensions []Dimension, stat, detectorID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if detectorID != "" {
		if !s.BaseStore.Exists(detectorID) {
			return fmt.Errorf("ResourceNotFound: anomaly detector %s not found", detectorID)
		}
		return s.BaseStore.Delete(detectorID)
	}

	key, err := singleMetricDetectorKey(namespace, metricName, dimensions, stat)
	if err != nil {
		return err
	}
	if !s.BaseStore.Exists(key) {
		return fmt.Errorf("ResourceNotFound: anomaly detector not found for %s/%s/%s", namespace, metricName, stat)
	}
	return s.BaseStore.Delete(key)
}

// DescribeAnomalyDetectors returns anomaly detectors matching the given
// filters. When anomalyDetectorIDs is non-empty, only those IDs are
// returned. Otherwise, filters by namespace, metricName, dimensions,
// and anomalyDetectorType are applied. Supports pagination via ListOpts.
func (s *AnomalyDetectorStore) DescribeAnomalyDetectors(opts DescribeAnomalyDetectorsOpts) ([]*AnomalyDetector, string, error) {
	if len(opts.AnomalyDetectorIDs) > 0 {
		var detectors []*AnomalyDetector
		for _, id := range opts.AnomalyDetectorIDs {
			d := &AnomalyDetector{}
			if err := s.BaseStore.Get(id, d); err != nil {
				continue
			}
			detectors = append(detectors, d)
		}
		return detectors, "", nil
	}

	prefix := "detector:"
	if opts.Namespace != "" {
		prefix = singleMetricKeyPrefix(opts.Namespace, opts.MetricName)
	}

	listOpts := opts.ListOpts
	listOpts.Prefix = prefix

	var filter func(*AnomalyDetector) bool
	if len(opts.AnomalyDetectorTypes) > 0 || opts.Namespace != "" || opts.MetricName != "" || len(opts.Dimensions) > 0 {
		filter = func(d *AnomalyDetector) bool {
			if len(opts.AnomalyDetectorTypes) > 0 {
				matched := false
				for _, t := range opts.AnomalyDetectorTypes {
					if d.AnomalyDetectorType == t {
						matched = true
						break
					}
				}
				if !matched {
					return false
				}
			}
			if opts.Namespace != "" && d.Namespace != opts.Namespace {
				return false
			}
			if opts.MetricName != "" && d.MetricName != opts.MetricName {
				return false
			}
			if len(opts.Dimensions) > 0 && !dimensionsMatch(d.Dimensions, opts.Dimensions) {
				return false
			}
			return true
		}
	}

	result, err := common.List[AnomalyDetector](s.BaseStore, listOpts, filter)
	if err != nil {
		return nil, "", err
	}
	return result.Items, result.NextMarker, nil
}

// DescribeAnomalyDetectorsOpts holds filter options for listing
// anomaly detectors.
type DescribeAnomalyDetectorsOpts struct {
	AnomalyDetectorIDs   []string
	AnomalyDetectorTypes []string
	Namespace            string
	MetricName           string
	Dimensions           []Dimension
	ListOpts             common.ListOptions
}

// singleMetricDetectorKey builds a deterministic storage key for a
// single-metric anomaly detector from its metric identity.
func singleMetricDetectorKey(namespace, metricName string, dimensions []Dimension, stat string) (string, error) {
	if namespace == "" || metricName == "" || stat == "" {
		return "", fmt.Errorf("namespace, metricName, and stat are required for single-metric anomaly detectors")
	}
	return singleMetricKeyPrefix(namespace, metricName) + ":" + hashDimensions(dimensions) + ":" + strings.ToLower(stat), nil
}

func singleMetricKeyPrefix(namespace, metricName string) string {
	if metricName == "" {
		return "detector:" + namespace
	}
	return "detector:" + namespace + ":" + metricName
}

// metricMathDetectorKey builds a deterministic key for a metric-math
// anomaly detector from its metric data queries.
func metricMathDetectorKey(queries []MetricDataQuery) (string, error) {
	if len(queries) == 0 {
		return "", fmt.Errorf("metric data queries are required for metric-math anomaly detectors")
	}
	data, err := json.Marshal(queries)
	if err != nil {
		return "", fmt.Errorf("failed to marshal metric data queries: %w", err)
	}
	hash := sha256.Sum256(data)
	return "detector:mm:" + hex.EncodeToString(hash[:16]), nil
}

func hashDimensions(dimensions []Dimension) string {
	if len(dimensions) == 0 {
		return "nodims"
	}
	dims := make([]string, len(dimensions))
	for i, d := range dimensions {
		dims[i] = d.Name + "=" + d.Value
	}
	sort.Strings(dims)
	combined := strings.Join(dims, ",")
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:8])
}

func dimensionsMatch(a, b []Dimension) bool {
	if len(a) != len(b) {
		return false
	}
	am := make(map[string]string, len(a))
	for _, d := range a {
		am[d.Name] = d.Value
	}
	for _, d := range b {
		if am[d.Name] != d.Value {
			return false
		}
	}
	return true
}

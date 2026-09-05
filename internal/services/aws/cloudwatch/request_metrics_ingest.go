package cloudwatch

import (
	"context"
	"fmt"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
)

// handleS3RequestMetrics ingests the aggregated S3 request-metrics
// datapoints the event bus delivers into the event region's metric store.
func (s *CloudWatchService) handleS3RequestMetrics(ctx context.Context, evt *eventbus.S3RequestMetricsEvent) eventbus.HandlerResult {
	if len(evt.Datapoints) == 0 {
		return eventbus.HandlerResult{}
	}
	metrics, err := s.regionMetricStore(evt.EventRegion())
	if err != nil {
		logs.Warn("cloudwatch: s3 request metrics dropped, no metric store",
			logs.String("region", evt.EventRegion()), logs.Err(err))
		return eventbus.HandlerResult{}
	}

	datums := make([]cwstore.MetricDatum, 0, len(evt.Datapoints))
	for _, dp := range evt.Datapoints {
		dimensions := make([]cwstore.Dimension, 0, len(dp.Dimensions))
		for _, d := range dp.Dimensions {
			dimensions = append(dimensions, cwstore.Dimension{Name: d.Name, Value: d.Value})
		}
		datum := cwstore.MetricDatum{
			MetricName: dp.MetricName,
			Value:      dp.Value,
			HasValue:   true,
			Timestamp:  dp.Timestamp,
			Unit:       cwstore.StandardUnit(dp.Unit),
			Dimensions: dimensions,
		}
		// The S3 aggregator publishes one statistic set per window, so
		// Average keeps its documented per-request meaning (the error rate,
		// bytes per request) instead of averaging the window sums.
		if dp.SampleCount > 0 {
			datum.StatisticValues = &cwstore.StatisticSet{
				SampleCount: dp.SampleCount,
				Sum:         dp.Value,
				Minimum:     dp.Minimum,
				Maximum:     dp.Maximum,
			}
		}
		datums = append(datums, datum)
	}
	if err := metrics.PutMetricData(requestMetricsNamespace, datums); err != nil {
		return eventbus.HandlerResult{Error: err}
	}
	return eventbus.HandlerResult{}
}

// regionMetricStore resolves one region's metric store without a request
// context, through the same storage manager the request plane uses.
func (s *CloudWatchService) regionMetricStore(region string) (*cwstore.MetricChunkStore, error) {
	if cached, ok := s.stores.Load(region); ok {
		if typed, ok := cached.(*cloudwatchStores); ok {
			return typed.metrics, nil
		}
	}
	if s.storageManager == nil {
		return nil, fmt.Errorf("cloudwatch: no storage manager for region %q", region)
	}
	storage, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, fmt.Errorf("cloudwatch: storage for region %q: %w", region, err)
	}
	metrics, err := cwstore.NewMetricChunkStoreWithIndex(storage, region, s.dataPath)
	if err != nil {
		return nil, fmt.Errorf("cloudwatch: metric store for region %q: %w", region, err)
	}
	stores := &cloudwatchStores{metrics: metrics}
	if actual, loaded := s.stores.LoadOrStore(region, stores); loaded {
		metrics.Close()
		if typed, ok := actual.(*cloudwatchStores); ok {
			return typed.metrics, nil
		}
	}
	return metrics, nil
}

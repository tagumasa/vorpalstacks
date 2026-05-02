package testutil

import (
	"fmt"
	"math"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

func (tc *cloudwatchTestCtx) metricTests() []TestResult {
	var results []TestResult

	namespace := tc.uniquePrefix("TestNS")
	metricName := tc.uniquePrefix("TestMetric")

	results = append(results, tc.runner.RunTest("cloudwatch", "PutMetricData_SingleValue", func() error {
		_, err := tc.client.PutMetricData(tc.ctx, &cloudwatch.PutMetricDataInput{
			Namespace: aws.String(namespace),
			MetricData: []types.MetricDatum{
				{
					MetricName: aws.String(metricName),
					Value:      aws.Float64(42.0),
					Unit:       types.StandardUnitCount,
					Timestamp:  aws.Time(time.Now()),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put metric data: %v", err)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudwatch", "ListMetrics_FilterByNamespace", func() error {
		resp, err := tc.client.ListMetrics(tc.ctx, &cloudwatch.ListMetricsInput{
			Namespace:  aws.String(namespace),
			MetricName: aws.String(metricName),
		})
		if err != nil {
			return fmt.Errorf("list metrics: %v", err)
		}
		if len(resp.Metrics) == 0 {
			return fmt.Errorf("expected at least 1 metric for ns=%s name=%s", namespace, metricName)
		}
		found := false
		for _, m := range resp.Metrics {
			if m.Namespace != nil && *m.Namespace == namespace && m.MetricName != nil && *m.MetricName == metricName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("metric %s/%s not found in ListMetrics result", namespace, metricName)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudwatch", "GetMetricStatistics_ReturnsValue", func() error {
		now := time.Now()
		resp, err := tc.client.GetMetricStatistics(tc.ctx, &cloudwatch.GetMetricStatisticsInput{
			Namespace:  aws.String(namespace),
			MetricName: aws.String(metricName),
			StartTime:  aws.Time(now.Add(-1 * time.Hour)),
			EndTime:    aws.Time(now),
			Period:     aws.Int32(300),
			Statistics: []types.Statistic{types.StatisticAverage},
		})
		if err != nil {
			return fmt.Errorf("get metric statistics: %v", err)
		}
		if len(resp.Datapoints) == 0 {
			return fmt.Errorf("expected at least 1 datapoint")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudwatch", "PutMetricData_GetMetricStatistics_Roundtrip", func() error {
		testNS := tc.uniquePrefix("RTNS")
		testMetric := tc.uniquePrefix("RTMetric")
		ts := time.Now().Add(-5 * time.Minute).Truncate(time.Minute)

		_, err := tc.client.PutMetricData(tc.ctx, &cloudwatch.PutMetricDataInput{
			Namespace: aws.String(testNS),
			MetricData: []types.MetricDatum{
				{
					MetricName: aws.String(testMetric),
					Value:      aws.Float64(42.0),
					Unit:       types.StandardUnitNone,
					Timestamp:  aws.Time(ts),
				},
				{
					MetricName: aws.String(testMetric),
					Value:      aws.Float64(58.0),
					Unit:       types.StandardUnitNone,
					Timestamp:  aws.Time(ts.Add(3 * time.Minute)),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		listResp, err := tc.client.ListMetrics(tc.ctx, &cloudwatch.ListMetricsInput{
			Namespace:  aws.String(testNS),
			MetricName: aws.String(testMetric),
		})
		if err != nil {
			return fmt.Errorf("list: %v", err)
		}
		if len(listResp.Metrics) == 0 {
			return fmt.Errorf("metric not found in ListMetrics")
		}
		if listResp.Metrics[0].Namespace == nil || *listResp.Metrics[0].Namespace != testNS {
			return fmt.Errorf("namespace mismatch: got %v", listResp.Metrics[0].Namespace)
		}
		if listResp.Metrics[0].MetricName == nil || *listResp.Metrics[0].MetricName != testMetric {
			return fmt.Errorf("metric name mismatch: got %v", listResp.Metrics[0].MetricName)
		}

		statsResp, err := tc.client.GetMetricStatistics(tc.ctx, &cloudwatch.GetMetricStatisticsInput{
			Namespace:  aws.String(testNS),
			MetricName: aws.String(testMetric),
			StartTime:  aws.Time(ts.Add(-1 * time.Minute)),
			EndTime:    aws.Time(ts.Add(10 * time.Minute)),
			Period:     aws.Int32(300),
			Statistics: []types.Statistic{types.StatisticSum},
		})
		if err != nil {
			return fmt.Errorf("get stats: %v", err)
		}
		if len(statsResp.Datapoints) == 0 {
			return fmt.Errorf("no datapoints returned")
		}
		var totalSum float64
		for _, dp := range statsResp.Datapoints {
			if dp.Sum != nil {
				totalSum += *dp.Sum
			}
		}
		if math.Abs(totalSum-100.0) > 0.01 {
			return fmt.Errorf("expected sum~100.0 (42+58), got %.2f", totalSum)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudwatch", "GetMetricData_Basic", func() error {
		testNS := tc.uniquePrefix("MDNS")
		testMetric := tc.uniquePrefix("MDMetric")
		now := time.Now()

		_, err := tc.client.PutMetricData(tc.ctx, &cloudwatch.PutMetricDataInput{
			Namespace: aws.String(testNS),
			MetricData: []types.MetricDatum{
				{
					MetricName: aws.String(testMetric),
					Value:      aws.Float64(10.0),
					Timestamp:  aws.Time(now.Add(-3 * time.Minute)),
				},
				{
					MetricName: aws.String(testMetric),
					Value:      aws.Float64(20.0),
					Timestamp:  aws.Time(now.Add(-1 * time.Minute)),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put metric data: %v", err)
		}

		resp, err := tc.client.GetMetricData(tc.ctx, &cloudwatch.GetMetricDataInput{
			StartTime: aws.Time(now.Add(-10 * time.Minute)),
			EndTime:   aws.Time(now.Add(1 * time.Minute)),
			MetricDataQueries: []types.MetricDataQuery{
				{
					Id: aws.String("m1"),
					MetricStat: &types.MetricStat{
						Metric: &types.Metric{
							Namespace:  aws.String(testNS),
							MetricName: aws.String(testMetric),
						},
						Period: aws.Int32(60),
						Stat:   aws.String("Sum"),
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("get metric data: %v", err)
		}
		if len(resp.MetricDataResults) == 0 {
			return fmt.Errorf("expected at least 1 MetricDataResult")
		}
		r := resp.MetricDataResults[0]
		if r.Id == nil || *r.Id != "m1" {
			return fmt.Errorf("expected id=m1, got %v", r.Id)
		}
		var total float64
		for _, v := range r.Values {
			total += v
		}
		if math.Abs(total-30.0) > 0.01 {
			return fmt.Errorf("expected sum=30.0 (10+20), got %.2f", total)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("cloudwatch", "PutMetricData_WithDimensions", func() error {
		dimNS := tc.uniquePrefix("DimNS")
		dimMetric := tc.uniquePrefix("DimMetric")
		_, err := tc.client.PutMetricData(tc.ctx, &cloudwatch.PutMetricDataInput{
			Namespace: aws.String(dimNS),
			MetricData: []types.MetricDatum{
				{
					MetricName: aws.String(dimMetric),
					Value:      aws.Float64(1.0),
					Dimensions: []types.Dimension{
						{Name: aws.String("Host"), Value: aws.String("server-1")},
						{Name: aws.String("Region"), Value: aws.String("us-east-1")},
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("put with dimensions: %v", err)
		}

		listResp, err := tc.client.ListMetrics(tc.ctx, &cloudwatch.ListMetricsInput{
			Namespace:  aws.String(dimNS),
			MetricName: aws.String(dimMetric),
		})
		if err != nil {
			return fmt.Errorf("list: %v", err)
		}
		if len(listResp.Metrics) == 0 {
			return fmt.Errorf("metric not found")
		}
		dims := listResp.Metrics[0].Dimensions
		if len(dims) != 2 {
			return fmt.Errorf("expected 2 dimensions, got %d", len(dims))
		}
		hasHost := false
		hasRegion := false
		for _, d := range dims {
			if d.Name != nil && *d.Name == "Host" && d.Value != nil && *d.Value == "server-1" {
				hasHost = true
			}
			if d.Name != nil && *d.Name == "Region" && d.Value != nil && *d.Value == "us-east-1" {
				hasRegion = true
			}
		}
		if !hasHost || !hasRegion {
			return fmt.Errorf("expected Host=server-1 and Region=us-east-1 dimensions")
		}
		return nil
	}))

	return results
}

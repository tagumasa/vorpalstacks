package cloudwatch

import (
	"strings"
	"time"

	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
	"vorpalstacks/pkg/metricmath"
)

// evaluateMetricMathAlarm evaluates an alarm that uses a Metrics array
// (metric math). It executes all MetricDataQuery entries, resolves
// Expression queries via the metricmath engine, and compares the
// ThresholdMetricId result against the alarm threshold.
func evaluateMetricMathAlarm(alarm *cwstore.Alarm, metricStore *cwstore.MetricChunkStore) *alarmEvalResult {
	now := time.Now().UTC()
	endTime := now.Truncate(time.Duration(alarm.Period) * time.Second)
	if endTime.After(now) {
		endTime = now
	}
	startTime := endTime.Add(-time.Duration(alarm.Period*alarm.EvaluationPeriods) * time.Second)

	queryResults := make(map[string][]metricmath.DataPoint)
	exprPending := make(map[string]*cwstore.MetricDataQuery)

	for i := range alarm.Metrics {
		q := &alarm.Metrics[i]
		if q.MetricStat != nil {
			statLower := strings.ToLower(q.MetricStat.Stat)
			isExtended := !isBasicStatistic(statLower)

			mq := cwstore.MetricQuery{
				Namespace:  q.MetricStat.Metric.Namespace,
				MetricName: q.MetricStat.Metric.MetricName,
				Dimensions: q.MetricStat.Metric.Dimensions,
				StartTime:  startTime,
				EndTime:    endTime,
				Period:     q.MetricStat.Period,
			}
			if isExtended {
				mq.ExtendedStatistics = []string{q.MetricStat.Stat}
			} else {
				mq.Statistics = []string{q.MetricStat.Stat}
			}
			stats, err := metricStore.GetMetricStatistics(mq)
			if err != nil {
				continue
			}
			dataPoints := make([]metricmath.DataPoint, 0, len(stats))
			for _, s := range stats {
				val := extractStatValue(s, statLower, q.MetricStat.Stat, isExtended)
				dataPoints = append(dataPoints, metricmath.DataPoint{Timestamp: s.Timestamp, Value: val})
			}
			queryResults[q.Id] = dataPoints
		} else if q.Expression != "" {
			exprPending[q.Id] = q
		}
	}

	for len(exprPending) > 0 {
		progress := false
		for id, q := range exprPending {
			ast, err := metricmath.Parse(q.Expression)
			if err != nil {
				delete(exprPending, id)
				continue
			}
			refs := ast.References()
			ready := true
			for _, ref := range refs {
				if _, ok := queryResults[ref]; !ok {
					if _, isExpr := exprPending[ref]; isExpr {
						ready = false
						break
					}
				}
			}
			if !ready {
				continue
			}
			result, err := ast.Eval(queryResults)
			if err == nil {
				queryResults[id] = result
			}
			delete(exprPending, id)
			progress = true
		}
		if !progress {
			break
		}
	}

	thresholdData, ok := queryResults[alarm.ThresholdMetricID]
	if !ok || len(thresholdData) == 0 {
		return determineStateTransition(alarm, 0, 0)
	}

	breachCount := 0
	for _, dp := range thresholdData {
		if dp.Timestamp.Before(startTime) || dp.Timestamp.After(endTime) {
			continue
		}
		if isBreaching(dp.Value, alarm.Threshold, alarm.ComparisonOperator) {
			breachCount++
		}
	}
	return determineStateTransition(alarm, breachCount, len(thresholdData))
}

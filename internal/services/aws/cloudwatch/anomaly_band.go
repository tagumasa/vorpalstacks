package cloudwatch

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
	"vorpalstacks/pkg/metricmath"
)

// hasAnomalyDetectionBand returns true if any of the given metric data
// queries contains an ANOMALY_DETECTION_BAND expression.
func hasAnomalyDetectionBand(metrics []cwstore.MetricDataQuery) bool {
	for _, q := range metrics {
		if q.Expression != "" && strings.Contains(
			strings.ToUpper(q.Expression), "ANOMALY_DETECTION_BAND") {
			return true
		}
	}
	return false
}

// resolveAnomalyDetectionBand processes an ANOMALY_DETECTION_BAND
// expression query by computing the EWMA band from the referenced
// metric's data points. The result is stored as two series in
// queryResults: {id} for the upper bound and {id}.low for the lower
// bound.
func resolveAnomalyDetectionBand(q *cwstore.MetricDataQuery, queryResults map[string][]metricmath.DataPoint) {
	stdDevMultiplier := extractStdDevMultiplier(q.Expression)

	// Find the referenced metric series. The expression format is
	// ANOMALY_DETECTION_BAND(metricId, N) — extract the metricId.
	metricID := extractBandMetricID(q.Expression)
	if metricID == "" {
		return
	}

	points, ok := queryResults[metricID]
	if !ok || len(points) == 0 {
		return
	}

	// Filter NaN values for band computation.
	valid := make([]metricmath.DataPoint, 0, len(points))
	for _, p := range points {
		if !math.IsNaN(p.Value) {
			valid = append(valid, p)
		}
	}
	if len(valid) == 0 {
		return
	}

	// Sort by timestamp.
	sort.Slice(valid, func(i, j int) bool {
		return valid[i].Timestamp.Before(valid[j].Timestamp)
	})

	upper, lower := computeBandFromDataPoints(valid, stdDevMultiplier)

	queryResults[q.Id] = upper
	queryResults[q.Id+".low"] = lower
}

// extractBandInnerArgs extracts the comma-separated arguments inside an
// ANOMALY_DETECTION_BAND(...) expression. Returns empty string if the
// expression does not contain a well-formed call.
func extractBandInnerArgs(expr string) string {
	upper := strings.ToUpper(expr)
	idx := strings.Index(upper, "ANOMALY_DETECTION_BAND(")
	if idx < 0 {
		return ""
	}
	rest := expr[idx+len("ANOMALY_DETECTION_BAND("):]
	depth := 0
	for i, c := range rest {
		if c == '(' {
			depth++
		} else if c == ')' {
			if depth == 0 {
				return rest[:i]
			}
			depth--
		}
	}
	return rest
}

// extractBandMetricID extracts the metric ID from an
// ANOMALY_DETECTION_BAND(metricId, N) expression.
func extractBandMetricID(expr string) string {
	rest := extractBandInnerArgs(expr)
	if rest == "" {
		return ""
	}
	commaIdx := strings.Index(rest, ",")
	if commaIdx < 0 {
		return strings.TrimSpace(rest)
	}
	return strings.TrimSpace(rest[:commaIdx])
}

// computeBandFromDataPoints computes the anomaly detection band from a
// series of DataPoints using EWMA + N*stddev. Returns upper and lower
// bound series.
func computeBandFromDataPoints(points []metricmath.DataPoint, stdDevMultiplier float64) (upper, lower []metricmath.DataPoint) {
	n := len(points)
	if n == 0 {
		return nil, nil
	}

	values := make([]float64, n)
	for i, p := range points {
		values[i] = p.Value
	}

	// Compute EWMA.
	ewma := make([]float64, n)
	ewma[0] = values[0]
	for i := 1; i < n; i++ {
		ewma[i] = ewmaAlpha*values[i] + (1-ewmaAlpha)*ewma[i-1]
	}

	// Compute residual stddev.
	residuals := make([]float64, n)
	for i := 0; i < n; i++ {
		residuals[i] = values[i] - ewma[i]
	}

	meanRes := mean(residuals)
	variance := 0.0
	for _, r := range residuals {
		variance += (r - meanRes) * (r - meanRes)
	}
	variance /= float64(n)
	stdDev := math.Sqrt(variance)
	bandWidth := stdDevMultiplier * stdDev

	upper = make([]metricmath.DataPoint, n)
	lower = make([]metricmath.DataPoint, n)
	for i := 0; i < n; i++ {
		upper[i] = metricmath.DataPoint{Timestamp: points[i].Timestamp, Value: ewma[i] + bandWidth}
		lower[i] = metricmath.DataPoint{Timestamp: points[i].Timestamp, Value: ewma[i] - bandWidth}
	}

	return upper, lower
}

// ewmaAlpha is the smoothing factor for the exponentially weighted
// moving average used in anomaly detection. A value of 0.2 gives
// moderate smoothing, similar to CloudWatch's default behaviour.
const ewmaAlpha = 0.2

// computeAnomalyBand computes an anomaly detection band from a series
// of metric statistics. It uses an exponentially weighted moving average
// (EWMA) and the standard deviation of the residuals to build upper and
// lower bounds.
//
// Parameters:
//   - stats: the historical metric statistics for the evaluation window
//   - stdDevMultiplier: the number of standard deviations for the band (e.g. 2)
//
// Returns two slices of float64 representing the upper and lower band
// values, aligned by index to the input stats slice.
func computeAnomalyBand(stats []*cwstore.MetricStatistics, stdDevMultiplier float64) (upper, lower []float64) {
	n := len(stats)
	if n == 0 {
		return nil, nil
	}

	values := make([]float64, n)
	for i, s := range stats {
		values[i] = s.Average
	}

	// Compute EWMA.
	ewma := make([]float64, n)
	ewma[0] = values[0]
	for i := 1; i < n; i++ {
		ewma[i] = ewmaAlpha*values[i] + (1-ewmaAlpha)*ewma[i-1]
	}

	// Compute residuals and their standard deviation.
	residuals := make([]float64, n)
	for i := 0; i < n; i++ {
		residuals[i] = values[i] - ewma[i]
	}

	meanResidual := mean(residuals)
	variance := 0.0
	for _, r := range residuals {
		variance += (r - meanResidual) * (r - meanResidual)
	}
	variance /= float64(n)
	stdDev := math.Sqrt(variance)

	bandWidth := stdDevMultiplier * stdDev

	upper = make([]float64, n)
	lower = make([]float64, n)
	for i := 0; i < n; i++ {
		upper[i] = ewma[i] + bandWidth
		lower[i] = ewma[i] - bandWidth
	}

	return upper, lower
}

// anomalyBandContext holds the precomputed band data shared between
// evaluateAnomalyBandAlarm and computeAlarmContributors.
type anomalyBandContext struct {
	stats     []*cwstore.MetricStatistics
	upper     []float64
	lower     []float64
	statLower string
	startTime time.Time
	endTime   time.Time
}

// prepareAnomalyBand performs the two-pass metric lookup and EWMA band
// computation shared by evaluateAnomalyBandAlarm and
// computeAlarmContributors.  Returns nil if the alarm does not contain
// an anomaly detection band or no metric data is available.
func prepareAnomalyBand(alarm *cwstore.Alarm, metricStore *cwstore.MetricChunkStore) *anomalyBandContext {
	// Two-pass approach: AWS standard Metrics array ordering places
	// the MetricStat entry before the ANOMALY_DETECTION_BAND expression,
	// so a single pass would never find the referenced metric.
	//
	// Pass 1: find the ANOMALY_DETECTION_BAND expression entry.
	var bandQuery *cwstore.MetricDataQuery
	var stdDevMultiplier float64 = 2 // default

	for i := range alarm.Metrics {
		q := &alarm.Metrics[i]
		if q.Expression != "" && containsAnomalyDetectionBand(q.Expression) {
			bandQuery = q
			stdDevMultiplier = extractStdDevMultiplier(q.Expression)
			break
		}
	}

	if bandQuery == nil {
		return nil
	}

	// Pass 2: find the metric query referenced by the band expression.
	referencedID := extractBandMetricID(bandQuery.Expression)
	var metricQuery *cwstore.MetricDataQuery

	for i := range alarm.Metrics {
		q := &alarm.Metrics[i]
		if q.Id == referencedID && q.MetricStat != nil {
			metricQuery = q
			break
		}
	}

	if metricQuery == nil || metricQuery.MetricStat == nil {
		return nil
	}

	// Fetch metric statistics for the evaluation window.
	now := time.Now().UTC()
	period := metricQuery.MetricStat.Period
	if period == 0 {
		period = alarm.Period
	}
	periodDuration := time.Duration(period) * time.Second
	endTime := now.Truncate(periodDuration)
	startTime := endTime.Add(-time.Duration(period*alarm.EvaluationPeriods) * time.Second)

	statLower := ""
	if metricQuery.MetricStat.Stat != "" {
		statLower = metricQuery.MetricStat.Stat
	} else {
		statLower = alarm.Statistic
	}

	mq := cwstore.MetricQuery{
		Namespace:  metricQuery.MetricStat.Metric.Namespace,
		MetricName: metricQuery.MetricStat.Metric.MetricName,
		Dimensions: metricQuery.MetricStat.Metric.Dimensions,
		StartTime:  startTime,
		EndTime:    endTime,
		Period:     period,
		Statistics: []string{statLower},
	}

	stats, err := metricStore.GetMetricStatistics(mq)
	if err != nil || len(stats) == 0 {
		return nil
	}

	// Sort stats by timestamp.
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Timestamp.Before(stats[j].Timestamp)
	})

	upper, lower := computeAnomalyBand(stats, stdDevMultiplier)

	return &anomalyBandContext{
		stats:     stats,
		upper:     upper,
		lower:     lower,
		statLower: statLower,
		startTime: startTime,
		endTime:   endTime,
	}
}

// evaluateAnomalyBandAlarm evaluates an alarm that uses anomaly detection.
// It finds the ANOMALY_DETECTION_BAND expression in the alarm's Metrics
// array, computes the band from the referenced metric's historical data,
// and counts how many data points fall outside the band.
//
// Returns (breachCount, totalDataPoints).
func evaluateAnomalyBandAlarm(alarm *cwstore.Alarm, metricStore *cwstore.MetricChunkStore) (int, int) {
	ctx := prepareAnomalyBand(alarm, metricStore)
	if ctx == nil {
		return 0, 0
	}

	breachCount := 0
	for i, s := range ctx.stats {
		if s.Timestamp.Before(ctx.startTime) || s.Timestamp.After(ctx.endTime) {
			continue
		}
		val := statisticValue(s, ctx.statLower)
		switch alarm.ComparisonOperator {
		case "LessThanLowerOrGreaterThanUpperThreshold":
			if val < ctx.lower[i] || val > ctx.upper[i] {
				breachCount++
			}
		case "LessThanLowerThreshold":
			if val < ctx.lower[i] {
				breachCount++
			}
		case "GreaterThanUpperThreshold":
			if val > ctx.upper[i] {
				breachCount++
			}
		}
	}

	return breachCount, len(ctx.stats)
}

// containsAnomalyDetectionBand checks if an expression string contains
// the ANOMALY_DETECTION_BAND function call.
func containsAnomalyDetectionBand(expr string) bool {
	return strings.Contains(strings.ToUpper(expr), "ANOMALY_DETECTION_BAND")
}

// extractStdDevMultiplier extracts the standard deviation multiplier
// from an ANOMALY_DETECTION_BAND expression. Returns 2.0 as default.
func extractStdDevMultiplier(expr string) float64 {
	rest := strings.TrimSpace(extractBandInnerArgs(expr))
	if rest == "" {
		return 2.0
	}
	commaIdx := strings.LastIndex(rest, ",")
	if commaIdx < 0 {
		return 2.0
	}
	numStr := strings.TrimSpace(rest[commaIdx+1:])
	var val float64
	if _, err := fmt.Sscanf(numStr, "%f", &val); err != nil || val <= 0 {
		return 2.0
	}
	return val
}

// mean computes the arithmetic mean of a float64 slice.
func mean(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals))
}

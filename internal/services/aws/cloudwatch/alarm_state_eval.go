package cloudwatch

import (
	"fmt"
	"strings"
	"time"

	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
)

// evaluateAlarm performs a single alarm evaluation by querying metric
// statistics for the alarm's configured namespace, metric name, dimensions,
// period, and statistic. It then compares the returned data points against
// the alarm's threshold and comparison operator.
//
// When the alarm has a Metrics array (metric math alarm), the queries are
// executed in dependency order and the ThresholdMetricId result is used.
//
// Returns nil when no state transition is needed.
func evaluateAlarm(alarm *cwstore.Alarm, metricStore *cwstore.MetricChunkStore) *alarmEvalResult {
	// Anomaly detection alarms: when the Metrics array contains an
	// ANOMALY_DETECTION_BAND expression and the ComparisonOperator is
	// one of the anomaly operators, use the band evaluation path.
	if len(alarm.Metrics) > 0 && hasAnomalyDetectionBand(alarm.Metrics) {
		breachCount, totalDataPoints := evaluateAnomalyBandAlarm(alarm, metricStore)
		return determineStateTransition(alarm, breachCount, totalDataPoints)
	}

	if len(alarm.Metrics) > 0 && alarm.ThresholdMetricID != "" {
		return evaluateMetricMathAlarm(alarm, metricStore)
	}

	now := time.Now().UTC()
	endTime := now.Truncate(time.Duration(alarm.Period) * time.Second)
	if endTime.After(now) {
		endTime = now
	}
	return evaluateAlarmWindow(alarm, metricStore, endTime)
}

// evaluateAlarmWindow evaluates one alarm over the EvaluationPeriods
// period buckets ending at endTime. The window is half-open
// [startTime, endTime): the bucket timestamped endTime belongs to the next
// evaluation, so consecutive evaluations never count a bucket twice.
func evaluateAlarmWindow(alarm *cwstore.Alarm, metricStore *cwstore.MetricChunkStore, endTime time.Time) *alarmEvalResult {
	startTime := endTime.Add(-time.Duration(alarm.Period*alarm.EvaluationPeriods) * time.Second)

	query := cwstore.MetricQuery{
		Namespace:  alarm.Namespace,
		MetricName: alarm.MetricName,
		Dimensions: alarm.Dimensions,
		StartTime:  startTime,
		EndTime:    endTime,
		Period:     alarm.Period,
		Statistics: []string{alarm.Statistic},
	}

	// When ExtendedStatistic is set (e.g. p90, p99), request it instead.
	if alarm.ExtendedStatistic != "" {
		query.Statistics = nil
		query.ExtendedStatistics = []string{alarm.ExtendedStatistic}
	}

	stats, err := metricStore.GetMetricStatistics(query)
	if err != nil {
		return nil
	}

	breachCount, inWindow := countBreaches(alarm, stats, startTime, endTime)

	return determineStateTransition(alarm, breachCount, inWindow)
}

// countBreaches iterates over the aggregated statistics returned by
// GetMetricStatistics and counts how many period buckets breach the
// alarm's threshold. Only buckets in the half-open window
// [startTime, endTime) are considered; each period bucket is expected to
// contain at most one aggregated data point. The second return value is
// the number of in-window buckets, used for missing-data accounting.
func countBreaches(alarm *cwstore.Alarm, stats []*cwstore.MetricStatistics, startTime, endTime time.Time) (breaches, inWindow int) {
	for _, s := range stats {
		if s.Timestamp.Before(startTime) || !s.Timestamp.Before(endTime) {
			continue
		}
		inWindow++

		var value float64
		if alarm.ExtendedStatistic != "" {
			value = statisticValue(s, alarm.ExtendedStatistic)
		} else {
			value = statisticValue(s, alarm.Statistic)
		}
		if isBreaching(value, alarm.Threshold, alarm.ComparisonOperator) {
			breaches++
		}
	}

	return breaches, inWindow
}

// statisticValue extracts the requested statistic (e.g. "Average", "Sum")
// from a MetricStatistics struct. For extended statistics (percentiles),
// the value is looked up from ExtendedStats.
func statisticValue(stats *cwstore.MetricStatistics, statistic string) float64 {
	if stats.ExtendedStats != nil {
		if v, ok := stats.ExtendedStats[statistic]; ok {
			return v
		}
	}
	switch strings.ToLower(statistic) {
	case "sum":
		return stats.Sum
	case "average":
		return stats.Average
	case "minimum":
		return stats.Minimum
	case "maximum":
		return stats.Maximum
	case "samplecount":
		return stats.SampleCount
	}
	return stats.Average
}

// isBreaching returns true when the given metric value satisfies the
// alarm's comparison operator against the threshold.
func isBreaching(value, threshold float64, operator string) bool {
	switch operator {
	case "GreaterThanOrEqualToThreshold":
		return value >= threshold
	case "GreaterThanThreshold":
		return value > threshold
	case "LessThanOrEqualToThreshold":
		return value <= threshold
	case "LessThanThreshold":
		return value < threshold
	case "LessThanLowerOrGreaterThanUpperThreshold":
		// Anomaly detection: value is breaching if outside the band.
		// Threshold represents the band boundary; without a full anomaly
		// model this falls back to a simple less-than check.
		return value < threshold
	case "LessThanLowerThreshold":
		return value < threshold
	case "GreaterThanUpperThreshold":
		return value > threshold
	default:
		return value >= threshold
	}
}

// determineStateTransition computes the new alarm state based on the
// number of breaching data points, the total number of evaluation periods,
// and the alarm's TreatMissingData configuration. When fewer data points
// than EvaluationPeriods are returned, the missing periods are handled
// according to TreatMissingData: "breaching" treats them as breaching;
// "notBreaching" and "missing" treat them as not breaching; "ignore"
// excludes them from evaluation entirely.
// Returns nil when the state has not changed.
func determineStateTransition(alarm *cwstore.Alarm, breachCount int, totalDataPoints int) *alarmEvalResult {
	oldState := alarm.State
	totalPeriods := int(alarm.EvaluationPeriods)
	datapointsToAlarm := int(alarm.DatapointsToAlarm)
	if datapointsToAlarm == 0 {
		datapointsToAlarm = totalPeriods
	}

	missingPeriods := totalPeriods - totalDataPoints
	if missingPeriods < 0 {
		missingPeriods = 0
	}

	// Per AWS docs, TreatMissingData controls how missing data points are
	// filled in when fewer real data points than EvaluationPeriods exist:
	//   "breaching"    — missing periods count as breaching
	//   "notBreaching" — missing periods count as not breaching
	//   "ignore"       — missing periods don't affect evaluation; maintain state
	//   "missing"      — missing periods don't count as breaching; if ALL data
	//                    points are missing, transition to INSUFFICIENT_DATA
	// When the alarm does not specify a value, "missing" is the AWS default
	// behaviour — an unset value must not silently behave like "breaching".
	treatMissing := alarm.TreatMissingData
	if treatMissing == "" {
		treatMissing = "missing"
	}

	var effectiveBreaches int
	switch treatMissing {
	case "breaching":
		effectiveBreaches = breachCount + missingPeriods
	default:
		// "missing", "notBreaching", "ignore" and any other stored value
		// keep the real breach count.
		effectiveBreaches = breachCount
	}

	var newState string
	var reason string

	switch {
	case effectiveBreaches >= datapointsToAlarm:
		newState = "ALARM"
		reason = fmt.Sprintf("Threshold Crossed: %d datapoints [%s] %s %v (threshold %v).",
			effectiveBreaches, alarm.Statistic, operatorPhrase(alarm.ComparisonOperator), alarm.Threshold, alarm.Threshold)

	case totalDataPoints == 0 && treatMissing == "missing":
		newState = "INSUFFICIENT_DATA"
		reason = fmt.Sprintf("Insufficient Metrics for %d datapoints. TreatMissingData=%s transitions to INSUFFICIENT_DATA.",
			missingPeriods, treatMissing)

	case treatMissing == "ignore" && totalDataPoints == 0:
		// For "ignore", retain the current state when ALL data points
		// are missing. If any real data points exist, evaluate normally.
		return nil

	default:
		newState = "OK"
		reason = fmt.Sprintf("Threshold Crossed: %d out of %d datapoints were not breaching.",
			totalPeriods-effectiveBreaches, totalPeriods)
	}

	if newState == oldState {
		return nil
	}

	var actionsToFire []string
	switch newState {
	case "ALARM":
		actionsToFire = alarm.AlarmActions
	case "OK":
		actionsToFire = alarm.OKActions
	case "INSUFFICIENT_DATA":
		actionsToFire = alarm.InsufficientDataActions
	}

	return &alarmEvalResult{
		alarm:         alarm,
		oldState:      oldState,
		newState:      newState,
		breachCount:   effectiveBreaches,
		reason:        reason,
		actionsToFire: actionsToFire,
	}
}

// isAlarmMuted checks if any ACTIVE alarm mute rule targets the given
// alarm name using the per-region mute rule store. Returns false if the
// mute rule store is unavailable.
func isAlarmMuted(alarmName string, muteRuleStore *cwstore.AlarmMuteRuleStore) bool {
	if muteRuleStore == nil {
		return false
	}
	return muteRuleStore.IsAlarmMuted(alarmName)
}

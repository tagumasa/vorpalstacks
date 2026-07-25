package cloudwatch

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/eventbus"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
	svcarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/pkg/metricmath"
)

const (
	// defaultEvalInterval is the default tick interval for the alarm
	// evaluator loop. One minute matches the minimum configurable alarm
	// Period, ensuring no evaluation window is missed.
	defaultEvalInterval = 60 * time.Second

	// testEvalInterval is the tick interval used when TEST_MODE is enabled.
	// A short interval allows integration tests to verify alarm evaluation
	// without waiting for the full default period.
	testEvalInterval = 1 * time.Second

	// defaultEvalWorkers is the number of concurrent goroutines used to
	// evaluate alarms in parallel.
	defaultEvalWorkers = 4
)

// alarmEvalResult captures the outcome of a single alarm evaluation.
// It is returned from evaluateAlarm and used by the dispatcher to decide
// whether a state transition has occurred.
type alarmEvalResult struct {
	alarm         *cwstore.Alarm
	oldState      string
	newState      string
	breachCount   int
	reason        string
	actionsToFire []string
}

// alarmEvaluator periodically evaluates all metric alarms against their
// configured metrics. When a state transition is detected it updates the
// alarm state in the store, publishes a CloudWatchAlarmStateEvent via the
// event bus, and dispatches alarm actions to SNS topics and Lambda
// functions.
type alarmEvaluator struct {
	mu       sync.Mutex
	running  bool
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	interval time.Duration
	workers  int
	logger   logs.Logger
}

// newAlarmEvaluator creates a new alarm evaluator with the given tick
// interval and worker count. If interval is zero, defaultEvalInterval is
// used; if workers is zero, defaultEvalWorkers is used.
func newAlarmEvaluator(interval time.Duration, workers int, logger logs.Logger) *alarmEvaluator {
	if interval <= 0 {
		interval = defaultEvalInterval
		if os.Getenv("TEST_MODE") == "true" {
			interval = testEvalInterval
		}
	}
	if workers <= 0 {
		workers = defaultEvalWorkers
	}
	return &alarmEvaluator{
		interval: interval,
		workers:  workers,
		logger:   logger,
	}
}

// Start launches the background evaluation loop. It is safe to call Start
// multiple times; subsequent calls are no-ops until Stop has been called.
func (e *alarmEvaluator) Start(ctx context.Context, s *CloudWatchService) {
	e.mu.Lock()
	if e.running {
		e.mu.Unlock()
		return
	}
	e.running = true
	ctx, e.cancel = context.WithCancel(ctx)
	e.mu.Unlock()

	e.wg.Add(1)
	go e.evalLoop(ctx, s)
}

// Stop gracefully shuts down the evaluation loop, waiting for any
// in-flight evaluations to complete.
func (e *alarmEvaluator) Stop() {
	e.mu.Lock()
	if !e.running {
		e.mu.Unlock()
		return
	}
	e.running = false
	if e.cancel != nil {
		e.cancel()
	}
	e.mu.Unlock()
	e.wg.Wait()
}

// evalLoop ticks at the configured interval, lists all metric alarms, and
// evaluates each one. Errors during individual alarm evaluation are logged
// but do not halt the loop.
func (e *alarmEvaluator) evalLoop(ctx context.Context, s *CloudWatchService) {
	defer e.wg.Done()
	defer func() { resilience.RecoverAndRestart("alarm evalLoop", &e.wg, func() { e.evalLoop(ctx, s) }) }()
	ticker := time.NewTicker(e.interval)
	defer ticker.Stop()

	// Perform an immediate first evaluation so that alarms configured before
	// the server started are evaluated without waiting for the first tick.
	e.evaluateAll(ctx, s)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			e.evaluateAll(ctx, s)
		}
	}
}

// evaluateAll fetches all metric alarms from all regions and dispatches
// them to a pool of worker goroutines for parallel evaluation.
func (e *alarmEvaluator) evaluateAll(ctx context.Context, s *CloudWatchService) {
	regions := s.getEvaluatorRegions()
	for _, region := range regions {
		e.evaluateAllForRegion(ctx, s, region)
	}
}

func (e *alarmEvaluator) evaluateAllForRegion(ctx context.Context, s *CloudWatchService, region string) {
	alarmStore, metricStore, muteRuleStore := s.evaluatorStoresForRegion(region)
	if alarmStore == nil || metricStore == nil {
		return
	}

	alarms, err := alarmStore.ListAlarms("")
	if err != nil {
		e.log("failed to list alarms for evaluation", "error", err)
		return
	}

	metricAlarms := make([]*cwstore.Alarm, 0, len(alarms))
	compositeAlarms := make([]*cwstore.Alarm, 0)
	for _, a := range alarms {
		if a.AlarmType == cwstore.AlarmTypeCompositeAlarm {
			compositeAlarms = append(compositeAlarms, a)
		} else {
			metricAlarms = append(metricAlarms, a)
		}
	}

	if len(metricAlarms) > 0 {
		e.evaluateMetricAlarms(ctx, s, metricAlarms, metricStore, alarmStore, muteRuleStore)
	}

	if len(compositeAlarms) > 0 {
		e.evaluateCompositeAlarms(ctx, s, compositeAlarms, alarmStore, muteRuleStore)
	}
}

func (e *alarmEvaluator) evaluateMetricAlarms(ctx context.Context, s *CloudWatchService, alarms []*cwstore.Alarm, metricStore *cwstore.MetricChunkStore, alarmStore *cwstore.AlarmStore, muteRuleStore *cwstore.AlarmMuteRuleStore) {
	type evalJob struct {
		alarm *cwstore.Alarm
	}
	jobs := make(chan evalJob, len(alarms))
	for _, a := range alarms {
		jobs <- evalJob{alarm: a}
	}
	close(jobs)

	var wg sync.WaitGroup
	for i := 0; i < e.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { resilience.RecoverPanic("cloudwatch alarm evaluator worker") }()
			for job := range jobs {
				result := evaluateAlarm(job.alarm, metricStore)
				if result == nil {
					continue
				}
				s.handleAlarmStateTransition(ctx, result, alarmStore, muteRuleStore)
			}
		}()
	}
	wg.Wait()
}

func (e *alarmEvaluator) evaluateCompositeAlarms(ctx context.Context, s *CloudWatchService, composites []*cwstore.Alarm, alarmStore *cwstore.AlarmStore, muteRuleStore *cwstore.AlarmMuteRuleStore) {
	// Parse all AlarmRules and build the dependency graph: for each
	// composite alarm, determine which other alarms (metric or composite)
	// it references. Only references to *composite* alarms create edges
	// in the topological sort; metric alarm references are always
	// already evaluated before this function is called.
	//
	// Two-pass approach: first register ALL composites in the lookup
	// map, then build dependency edges. This ensures forward references
	// (e.g. composite A references composite B, but B appears later in
	// the slice) are correctly detected.
	nameToComposite := make(map[string]*cwstore.Alarm, len(composites))
	parsedRules := make(map[string]alarmRuleNode, len(composites))

	// Pass 1: register all composites and parse their rules.
	for _, c := range composites {
		nameToComposite[c.Name] = c
		node, err := parseAlarmRule(c.AlarmRule)
		if err != nil {
			e.log("failed to parse alarm rule", "alarm", c.Name, "rule", c.AlarmRule, "error", err)
			continue
		}
		parsedRules[c.Name] = node
	}

	// Pass 2: build dependency edges using the fully-populated map.
	dependencies := make(map[string][]string, len(composites))
	for _, c := range composites {
		node := parsedRules[c.Name]
		if node == nil {
			dependencies[c.Name] = nil
			continue
		}
		var compositeDeps []string
		for _, childName := range node.childAlarmNames() {
			if _, isComposite := nameToComposite[childName]; isComposite {
				compositeDeps = append(compositeDeps, childName)
			}
		}
		dependencies[c.Name] = compositeDeps
	}

	// Topological sort using Kahn's algorithm. Process in levels so that
	// all alarms at the same level (no inter-dependencies) can be
	// evaluated in a single pass.
	ordered, err := topologicalSortLevels(dependencies)
	if err != nil {
		e.log("failed to topologically sort composite alarms", "error", err)
		return
	}

	// Fetch all alarm states once and build a lookup map. This avoids
	// calling ListAlarms inside evaluateCompositeAlarm for each composite.
	allAlarms, err := alarmStore.ListAlarms("")
	if err != nil {
		e.log("failed to list alarms for composite evaluation", "error", err)
		return
	}
	alarmStateMap := make(map[string]string, len(allAlarms))
	for _, a := range allAlarms {
		alarmStateMap[a.Name] = a.State
	}

	for _, level := range ordered {
		for _, name := range level {
			alarm := nameToComposite[name]
			if alarm == nil {
				continue
			}
			rule := parsedRules[name]
			if rule == nil {
				continue
			}
			result := evaluateCompositeAlarm(alarm, rule, alarmStateMap)
			if result == nil {
				continue
			}
			s.handleAlarmStateTransition(ctx, result, alarmStore, muteRuleStore)
		}
	}
}

// evaluateCompositeAlarm evaluates a composite alarm by parsing its
// AlarmRule expression, looking up the current state of all referenced
// child alarms, and evaluating the boolean expression. Returns nil when
// no state transition is needed.
func evaluateCompositeAlarm(alarm *cwstore.Alarm, rule alarmRuleNode, alarmStateMap map[string]string) *alarmEvalResult {
	childNames := rule.childAlarmNames()

	// Build a map of alarm name -> is-in-ALARM-state using the
	// pre-fetched state map.
	alarmStates := make(map[string]bool, len(childNames))
	for _, name := range childNames {
		state := alarmStateMap[name]
		alarmStates[name] = state == "ALARM"
	}

	isBreaching := rule.evaluate(alarmStates)

	oldState := alarm.State
	var newState string
	var reason string
	if isBreaching {
		newState = "ALARM"
		reason = fmt.Sprintf("Composite alarm rule evaluated to ALARM (was %s).", oldState)
	} else {
		newState = "OK"
		reason = fmt.Sprintf("Composite alarm rule evaluated to OK (was %s).", oldState)
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
		breachCount:   0,
		reason:        reason,
		actionsToFire: actionsToFire,
	}
}

// topologicalSortLevels performs a level-by-level topological sort using
// Kahn's algorithm. The input is a map of node → nodes-it-depends-on
// (i.e., nodes that must be evaluated before it). The output is a slice
// of levels, where each level contains nodes that can be evaluated in
// parallel. Returns an error if the graph contains a cycle.
func topologicalSortLevels(dependencies map[string][]string) ([][]string, error) {
	inDegree := make(map[string]int)
	dependents := make(map[string][]string)

	for node := range dependencies {
		if _, exists := inDegree[node]; !exists {
			inDegree[node] = 0
		}
		for _, dep := range dependencies[node] {
			if _, exists := inDegree[dep]; !exists {
				inDegree[dep] = 0
			}
			inDegree[node]++
			dependents[dep] = append(dependents[dep], node)
		}
	}

	var levels [][]string
	for {
		var level []string
		for node, deg := range inDegree {
			if deg == 0 {
				level = append(level, node)
			}
		}
		if len(level) == 0 {
			break
		}
		for _, node := range level {
			delete(inDegree, node)
		}
		for _, node := range level {
			for _, dependent := range dependents[node] {
				if _, exists := inDegree[dependent]; exists {
					inDegree[dependent]--
				}
			}
		}
		levels = append(levels, level)
	}

	if len(inDegree) > 0 {
		var cyclic []string
		for node := range inDegree {
			cyclic = append(cyclic, node)
		}
		return nil, fmt.Errorf("cyclic dependency detected among composite alarms: %v", cyclic)
	}

	return levels, nil
}

// statistics for the alarm's configured namespace, metric name, dimensions,
// period, and statistic. It then compares the returned data points against
// the alarm's threshold and comparison operator.
//
// Returns nil when no state transition is needed.
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

// checkActionsSuppressor evaluates whether a composite alarm's actions
// should be suppressed based on its ActionsSuppressor configuration. The
// suppressor alarm must be in ALARM state for at least WaitPeriod seconds
// before suppression begins. After the suppressor transitions out of ALARM,
// suppression continues for ExtensionPeriod seconds.
//
// Returns (true, reason) if actions should be suppressed.
func checkActionsSuppressor(alarm *cwstore.Alarm, alarmStore *cwstore.AlarmStore) (bool, string) {
	suppressorName := extractAlarmNameFromARN(alarm.ActionsSuppressor)
	if suppressorName == "" {
		suppressorName = alarm.ActionsSuppressor
	}

	suppressor, err := alarmStore.GetAlarm(suppressorName)
	if err != nil || suppressor == nil {
		return false, ""
	}

	now := time.Now().UTC()
	if suppressor.State == "ALARM" {
		waitDuration := time.Duration(alarm.ActionsSuppressorWaitPeriod) * time.Second
		if now.Sub(suppressor.StateUpdatedTimestamp) >= waitDuration {
			return true, "WaitPeriod"
		}
	} else {
		extDuration := time.Duration(alarm.ActionsSuppressorExtPeriod) * time.Second
		if extDuration > 0 && now.Sub(suppressor.StateUpdatedTimestamp) < extDuration {
			return true, "ExtensionPeriod"
		}
	}
	return false, ""
}

// extractAlarmNameFromARN extracts the alarm name from a CloudWatch alarm
// ARN. Returns empty string if the ARN format is not recognised.
func extractAlarmNameFromARN(arn string) string {
	parts := strings.Split(arn, ":")
	if len(parts) < 6 {
		return ""
	}
	resource := parts[5]
	if strings.HasPrefix(resource, "alarm:") {
		return strings.TrimPrefix(resource, "alarm:")
	}
	return ""
}

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

	breachCount := countBreaches(alarm, stats, startTime, endTime)

	return determineStateTransition(alarm, breachCount, len(stats))
}

// countBreaches iterates over the aggregated statistics returned by
// GetMetricStatistics and counts how many period buckets breach the
// alarm's threshold. Each period bucket is expected to contain at most
// one aggregated data point.
func countBreaches(alarm *cwstore.Alarm, stats []*cwstore.MetricStatistics, startTime, endTime time.Time) int {
	if len(stats) == 0 {
		return 0
	}

	breaches := 0

	for _, s := range stats {
		if s.Timestamp.Before(startTime) || s.Timestamp.After(endTime) {
			continue
		}

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

	return breaches
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
// according to TreatMissingData: "missing" and "breaching" treat them as
// breaching; "ignore" and "notBreaching" treat them as not breaching.
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
	//   "breaching"  — missing periods count as breaching
	//   "notBreaching" — missing periods count as not breaching
	//   "ignore"     — missing periods don't affect evaluation; maintain state
	//   "missing"    — missing periods don't count as breaching; if ALL data
	//                  points are missing, transition to INSUFFICIENT_DATA
	var effectiveBreaches int
	switch alarm.TreatMissingData {
	case "ignore", "notBreaching", "missing":
		effectiveBreaches = breachCount
	case "breaching":
		effectiveBreaches = breachCount + missingPeriods
	default:
		effectiveBreaches = breachCount + missingPeriods
	}

	var newState string
	var reason string

	switch {
	case effectiveBreaches >= datapointsToAlarm:
		newState = "ALARM"
		reason = fmt.Sprintf("Threshold Crossed: %d datapoints [%s] %s %v (threshold %v).",
			effectiveBreaches, alarm.Statistic, operatorPhrase(alarm.ComparisonOperator), alarm.Threshold, alarm.Threshold)

	case totalDataPoints == 0 && alarm.TreatMissingData == "missing":
		newState = "INSUFFICIENT_DATA"
		reason = fmt.Sprintf("Insufficient Metrics for %d datapoints. TreatMissingData=%s transitions to INSUFFICIENT_DATA.",
			missingPeriods, alarm.TreatMissingData)

	case alarm.TreatMissingData == "ignore" && totalDataPoints == 0:
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

// handleAlarmStateTransition is called when the evaluator detects an alarm
// state change. It updates the alarm state in the store, records alarm
// history, publishes a CloudWatchAlarmStateEvent via the event bus, and
// dispatches alarm actions to SNS topics and Lambda functions.
func (s *CloudWatchService) handleAlarmStateTransition(ctx context.Context, result *alarmEvalResult, alarmStore *cwstore.AlarmStore, muteRuleStore *cwstore.AlarmMuteRuleStore) {
	alarm := result.alarm

	if err := alarmStore.SetAlarmState(alarm.Name, result.newState, result.reason, ""); err != nil {
		s.log("failed to set alarm state", "alarm", alarm.Name, "new_state", result.newState, "error", err)
		return
	}

	historyType := cwstore.AlarmTypeMetricAlarm
	if alarm.AlarmType == cwstore.AlarmTypeCompositeAlarm {
		historyType = cwstore.AlarmTypeCompositeAlarm
	}

	if err := alarmStore.AddAlarmHistory(&cwstore.AlarmHistoryEntry{
		AlarmName:       alarm.Name,
		AlarmType:       historyType,
		Timestamp:       time.Now().UTC().UnixMilli(),
		HistoryItemType: cwstore.HistoryItemTypeStateUpdate,
		HistorySummary:  result.reason,
	}); err != nil {
		s.log("failed to add alarm history", "alarm", alarm.Name, "error", err)
	}

	if !alarm.ActionsEnabled {
		return
	}

	// Alarm Mute Rules: if any ACTIVE mute rule targets this alarm,
	// suppress all alarm actions.
	if isAlarmMuted(result.alarm.Name, muteRuleStore) {
		return
	}

	// ActionsSuppressor: when the suppressor alarm is in ALARM state and
	// has been so for at least ActionsSuppressorWaitPeriod seconds, the
	// composite alarm's actions are suppressed. When the suppressor
	// transitions out of ALARM, suppression continues for
	// ActionsSuppressorExtensionPeriod seconds.
	if alarm.ActionsSuppressor != "" {
		if suppressed, reason := checkActionsSuppressor(alarm, alarmStore); suppressed {
			if err := alarmStore.SetAlarmActionsSuppressed(alarm.Name, reason); err != nil {
				s.log("failed to set actions suppressed", "alarm", alarm.Name, "error", err)
			}
			result.actionsToFire = nil
		}
	}

	s.publishAlarmStateEvent(ctx, result)
	s.dispatchAlarmActions(ctx, result)
}

// publishAlarmStateEvent publishes a CloudWatchAlarmStateEvent to the
// event bus. The event carries the alarm ARN, previous state, new state,
// and the reason for the transition.
func (s *CloudWatchService) publishAlarmStateEvent(ctx context.Context, result *alarmEvalResult) {
	if s.bus == nil {
		return
	}

	_, _, alarmRegion, _, _ := svcarn.SplitARN(result.alarm.ARN)
	if alarmRegion == "" {
		alarmRegion = s.region
	}

	evt := &eventbus.CloudWatchAlarmStateEvent{
		EventBase: eventbus.EventBase{
			Timestamp: time.Now().UTC(),
			Source:    "aws:cloudwatch",
			Region:    alarmRegion,
			AccountID: s.accountID,
			Caller: eventbus.CallerContext{
				ServicePrincipal: "cloudwatch.amazonaws.com",
				AccountID:        s.accountID,
			},
		},
		AlarmName:     result.alarm.Name,
		AlarmARN:      result.alarm.ARN,
		PreviousState: result.oldState,
		NewState:      result.newState,
		Reason:        result.reason,
	}

	if err := s.bus.Publish(ctx, evt); err != nil {
		logs.Warn("failed to publish alarm state change event", logs.String("alarmName", result.alarm.Name), logs.Err(err))
	}
}

// dispatchAlarmActions iterates over the action ARNs for the new state and
// dispatches notifications to SNS topics (via the event bus) and Lambda
// functions (via the direct invoker). Region and account ID are extracted
// from each action ARN to support cross-region alarm actions.
func (s *CloudWatchService) dispatchAlarmActions(ctx context.Context, result *alarmEvalResult) {
	for _, actionArn := range result.actionsToFire {
		switch svcarn.GetServiceFromARN(actionArn) {
		case "sns":
			s.dispatchAlarmToSNS(ctx, actionArn, result)
		case "lambda":
			s.dispatchAlarmToLambda(ctx, actionArn, result)
		case "states":
			s.dispatchAlarmToStepFunctions(ctx, actionArn, result)
		}
	}
}

// dispatchAlarmToSNS publishes the alarm state change notification to an
// SNS topic via the event bus. Region and account ID are extracted from the
// topic ARN.
func (s *CloudWatchService) dispatchAlarmToSNS(ctx context.Context, topicArn string, result *alarmEvalResult) {
	if s.bus == nil {
		return
	}

	allowed, evalErr := s.bus.EvaluateTargetPolicy(ctx, topicArn, "sns", "cloudwatch.amazonaws.com", "sns:Publish", topicArn)
	if evalErr != nil {
		s.log("resource policy evaluation failed for alarm SNS delivery, dropping notification", "topicArn", topicArn, "error", evalErr)
		return
	}
	if !allowed {
		return
	}

	_, _, region, accountID, _ := svcarn.SplitARN(topicArn)
	messageID := fmt.Sprintf("%d", time.Now().UnixNano())

	payload := buildAlarmNotificationPayload(result)
	messageBytes, _ := json.Marshal(payload)

	snsEvt := &eventbus.SNSDeliveryEvent{
		EventBase: eventbus.EventBase{
			Timestamp: time.Now().UTC(),
			Source:    "aws:cloudwatch",
			Region:    region,
			AccountID: accountID,
			Caller: eventbus.CallerContext{
				ServicePrincipal: "cloudwatch.amazonaws.com",
				AccountID:        accountID,
			},
		},
		TopicARN:  topicArn,
		MessageID: messageID,
		Message:   string(messageBytes),
		Subject:   fmt.Sprintf("ALARM: \"%s\" in %s", result.alarm.Name, result.newState),
	}
	snsEvt.Region = region

	if err := s.bus.Publish(ctx, snsEvt); err != nil {
		logs.Warn("failed to publish alarm SNS notification", logs.String("alarmName", result.alarm.Name), logs.Err(err))
	}
}

// dispatchAlarmToLambda invokes a Lambda function with the alarm state
// change notification payload. The function name is extracted from the
// function ARN.
func (s *CloudWatchService) dispatchAlarmToLambda(ctx context.Context, functionArn string, result *alarmEvalResult) {
	if s.bus == nil {
		return
	}

	allowed, evalErr := s.bus.EvaluateTargetPolicy(ctx, functionArn, "lambda", "cloudwatch.amazonaws.com", "lambda:InvokeFunction", functionArn)
	if evalErr != nil {
		s.log("resource policy evaluation failed for alarm Lambda delivery, dropping notification", "functionArn", functionArn, "error", evalErr)
		return
	}
	if !allowed {
		return
	}

	fnName := svcarn.ExtractFunctionNameFromARN(functionArn)
	payload := buildAlarmNotificationPayload(result)
	payloadBytes, _ := json.Marshal(payload)

	_, _, err := s.bus.LambdaInvoker().InvokeForGateway(ctx, functionArn, payloadBytes)
	if err != nil {
		s.log("lambda dispatch failed for alarm action", "alarm", result.alarm.Name, "function", fnName, "error", err)
	}
}

func (s *CloudWatchService) dispatchAlarmToStepFunctions(ctx context.Context, stateMachineArn string, result *alarmEvalResult) {
	if s.bus == nil {
		return
	}

	_, _, smRegion, _, _ := svcarn.SplitARN(stateMachineArn)
	if smRegion == "" {
		smRegion = s.region
	}

	payload := buildAlarmNotificationPayload(result)
	payloadBytes, _ := json.Marshal(payload)

	evt := &eventbus.StepFunctionsStartExecutionEvent{
		StateMachineArn: stateMachineArn,
		Input:           string(payloadBytes),
	}
	evt.Region = smRegion
	evt.AccountID = s.accountID

	if err := s.bus.Publish(ctx, evt); err != nil {
		logs.Warn("failed to publish alarm Step Function event", logs.String("alarmName", result.alarm.Name), logs.Err(err))
	}
}

// operatorPhrase returns a human-readable phrase describing the comparison
// direction, suitable for inclusion in alarm state change reason strings.
func operatorPhrase(operator string) string {
	switch operator {
	case "GreaterThanOrEqualToThreshold":
		return "were at or above"
	case "GreaterThanThreshold":
		return "were above"
	case "LessThanOrEqualToThreshold":
		return "were at or below"
	case "LessThanThreshold":
		return "were below"
	case "LessThanLowerOrGreaterThanUpperThreshold":
		return "were outside the anomaly band"
	case "LessThanLowerThreshold":
		return "were below the anomaly band lower bound"
	case "GreaterThanUpperThreshold":
		return "were above the anomaly band upper bound"
	default:
		return "crossed"
	}
}

// buildAlarmNotificationPayload constructs the CloudWatch alarm
// notification payload matching the format AWS sends to SNS topics and
// Lambda functions. This includes the alarm description, metric details,
// and state transition information.
func buildAlarmNotificationPayload(result *alarmEvalResult) map[string]interface{} {
	alarm := result.alarm
	now := time.Now().UTC()

	_, _, alarmRegion, _, _ := svcarn.SplitARN(alarm.ARN)

	return map[string]interface{}{
		"AlarmName":          alarm.Name,
		"AlarmArn":           alarm.ARN,
		"AlarmDescription":   alarm.AlarmDescription,
		"AlarmConfiguration": buildAlarmConfiguration(alarm),
		"PreviousState": map[string]interface{}{
			"StateValue":      result.oldState,
			"StateReason":     "",
			"StateReasonData": "",
		},
		"NewState": map[string]interface{}{
			"StateValue":      result.newState,
			"StateReason":     result.reason,
			"StateReasonData": "",
			"TriggeredTime":   now.Format(time.RFC3339),
		},
		"NewStateReason":     result.reason,
		"StateChangeTime":    now.Format(time.RFC3339),
		"Region":             alarmRegion,
		"MetricName":         alarm.MetricName,
		"Namespace":          alarm.Namespace,
		"Statistic":          alarm.Statistic,
		"Period":             alarm.Period,
		"EvaluationPeriods":  alarm.EvaluationPeriods,
		"Threshold":          alarm.Threshold,
		"ComparisonOperator": alarm.ComparisonOperator,
		"TreatMissingData":   alarm.TreatMissingData,
	}
}

// buildAlarmConfiguration serialises the alarm's key configuration fields
// into a nested map for inclusion in the notification payload.
func buildAlarmConfiguration(alarm *cwstore.Alarm) map[string]interface{} {
	config := map[string]interface{}{
		"AlarmName":          alarm.Name,
		"AlarmArn":           alarm.ARN,
		"AlarmType":          alarm.AlarmType,
		"MetricName":         alarm.MetricName,
		"Namespace":          alarm.Namespace,
		"Statistic":          alarm.Statistic,
		"Period":             alarm.Period,
		"EvaluationPeriods":  alarm.EvaluationPeriods,
		"Threshold":          alarm.Threshold,
		"ComparisonOperator": alarm.ComparisonOperator,
		"TreatMissingData":   alarm.TreatMissingData,
		"ActionsEnabled":     alarm.ActionsEnabled,
	}

	if alarm.DatapointsToAlarm > 0 {
		config["DatapointsToAlarm"] = alarm.DatapointsToAlarm
	}
	if len(alarm.Dimensions) > 0 {
		dims := make([]map[string]string, len(alarm.Dimensions))
		for i, d := range alarm.Dimensions {
			dims[i] = map[string]string{"name": d.Name, "value": d.Value}
		}
		config["Dimensions"] = dims
	}
	if len(alarm.AlarmActions) > 0 {
		config["AlarmActions"] = alarm.AlarmActions
	}
	if len(alarm.OKActions) > 0 {
		config["OKActions"] = alarm.OKActions
	}
	if len(alarm.InsufficientDataActions) > 0 {
		config["InsufficientDataActions"] = alarm.InsufficientDataActions
	}

	return config
}

// evaluatorStoresForRegion returns the alarm store, metric store, and
// alarm mute rule store for the specified region.
func (s *CloudWatchService) evaluatorStoresForRegion(region string) (*cwstore.AlarmStore, *cwstore.MetricChunkStore, *cwstore.AlarmMuteRuleStore) {
	if region == "" {
		return nil, nil, nil
	}

	if cached, ok := s.stores.Load(region); ok {
		if typed, ok := cached.(*cloudwatchStores); ok {
			return typed.alarms, typed.metrics, typed.alarmMuteRules
		}
	}

	var storage storage.BasicStorage
	if s.storageManager != nil {
		var err error
		storage, err = s.storageManager.GetStorage(region)
		if err != nil {
			return nil, nil, nil
		}
	}
	if storage == nil {
		return nil, nil, nil
	}

	alarmStore := cwstore.NewAlarmStore(storage, s.accountID, region)
	metricStore, err := cwstore.NewMetricChunkStoreWithIndex(storage, region, s.dataPath)
	if err != nil {
		return alarmStore, nil, nil
	}

	stores := &cloudwatchStores{
		metrics:          metricStore,
		alarms:           alarmStore,
		dashboards:       cwstore.NewDashboardStore(storage, s.accountID, region),
		anomalyDetectors: cwstore.NewAnomalyDetectorStore(storage, s.accountID, region),
		insightRules:     cwstore.NewInsightRuleStore(storage, region),
		alarmMuteRules:   cwstore.NewAlarmMuteRuleStore(storage, s.accountID, region),
	}
	if actual, loaded := s.stores.LoadOrStore(region, stores); loaded {
		metricStore.Close()
		if typed, ok := actual.(*cloudwatchStores); ok {
			return typed.alarms, typed.metrics, typed.alarmMuteRules
		}
	}

	return alarmStore, metricStore, stores.alarmMuteRules
}

// getEvaluatorRegions returns the list of regions to evaluate alarms for.
func (s *CloudWatchService) getEvaluatorRegions() []string {
	if s.storageManager != nil {
		return s.storageManager.GetActiveRegions()
	}
	if s.region != "" {
		return []string{s.region}
	}
	return nil
}

// log emits a structured log message if a logger is configured on the
// service. Used by the alarm evaluator and action dispatch methods.
func (s *CloudWatchService) log(msg string, keyvals ...interface{}) {
	if s.logger == nil {
		return
	}
	fields := make([]logs.Field, 0, len(keyvals)/2)
	for i := 0; i+1 < len(keyvals); i += 2 {
		fields = append(fields, logs.Field{Key: fmt.Sprint(keyvals[i]), Value: keyvals[i+1]})
	}
	s.logger.Info(msg, fields...)
}

// log is a convenience method on the alarmEvaluator for structured logging.
func (e *alarmEvaluator) log(msg string, keyvals ...interface{}) {
	if e.logger == nil {
		return
	}
	fields := make([]logs.Field, 0, len(keyvals)/2)
	for i := 0; i+1 < len(keyvals); i += 2 {
		fields = append(fields, logs.Field{Key: fmt.Sprint(keyvals[i]), Value: keyvals[i+1]})
	}
	e.logger.Info(msg, fields...)
}

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
	alarmStore, metricStore := s.evaluatorStoresForRegion(region)
	if alarmStore == nil || metricStore == nil {
		return
	}

	alarms, err := alarmStore.ListAlarms("")
	if err != nil {
		e.log("failed to list alarms for evaluation", "error", err)
		return
	}

	metricAlarms := make([]*cwstore.Alarm, 0, len(alarms))
	for _, a := range alarms {
		if a.AlarmType == "" || a.AlarmType == cwstore.AlarmTypeMetricAlarm {
			metricAlarms = append(metricAlarms, a)
		}
	}

	if len(metricAlarms) == 0 {
		return
	}

	type evalJob struct {
		alarm *cwstore.Alarm
	}
	jobs := make(chan evalJob, len(metricAlarms))
	for _, a := range metricAlarms {
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
				s.handleAlarmStateTransition(ctx, result, alarmStore)
			}
		}()
	}
	wg.Wait()
}

// evaluateAlarm performs a single alarm evaluation by querying metric
// statistics for the alarm's configured namespace, metric name, dimensions,
// period, and statistic. It then compares the returned data points against
// the alarm's threshold and comparison operator.
//
// Returns nil when no state transition is needed.
func evaluateAlarm(alarm *cwstore.Alarm, metricStore *cwstore.MetricChunkStore) *alarmEvalResult {
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

		value := statisticValue(s, alarm.Statistic)
		if isBreaching(value, alarm.Threshold, alarm.ComparisonOperator) {
			breaches++
		}
	}

	return breaches
}

// statisticValue extracts the requested statistic (e.g. "Average", "Sum")
// from a MetricStatistics struct.
func statisticValue(stats *cwstore.MetricStatistics, statistic string) float64 {
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
	case "LessThanOrEqualToLowerOrGreaterThanUpperThreshold":
		return value <= threshold
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

	var effectiveBreaches int
	switch alarm.TreatMissingData {
	case "ignore", "notBreaching":
		effectiveBreaches = breachCount
	case "breaching", "missing":
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

	case missingPeriods >= datapointsToAlarm && alarm.TreatMissingData != "ignore" && alarm.TreatMissingData != "notBreaching":
		newState = "INSUFFICIENT_DATA"
		reason = fmt.Sprintf("Threshold Crossed: %d datapoints were missing. TreatMissingData=%s transitions to INSUFFICIENT_DATA.",
			missingPeriods, alarm.TreatMissingData)

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

// handleAlarmStateTransition is called when the evaluator detects an alarm
// state change. It updates the alarm state in the store, records alarm
// history, publishes a CloudWatchAlarmStateEvent via the event bus, and
// dispatches alarm actions to SNS topics and Lambda functions.
func (s *CloudWatchService) handleAlarmStateTransition(ctx context.Context, result *alarmEvalResult, alarmStore *cwstore.AlarmStore) {
	alarm := result.alarm

	if err := alarmStore.SetAlarmState(alarm.Name, result.newState, result.reason); err != nil {
		s.log("failed to set alarm state", "alarm", alarm.Name, "new_state", result.newState, "error", err)
		return
	}

	if err := alarmStore.AddAlarmHistory(&cwstore.AlarmHistoryEntry{
		AlarmName:       alarm.Name,
		AlarmType:       cwstore.AlarmTypeMetricAlarm,
		Timestamp:       time.Now().UTC().UnixMilli(),
		HistoryItemType: cwstore.HistoryItemTypeStateUpdate,
		HistorySummary:  result.reason,
	}); err != nil {
		s.log("failed to add alarm history", "alarm", alarm.Name, "error", err)
	}

	if !alarm.ActionsEnabled {
		return
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
	case "LessThanOrEqualToLowerOrGreaterThanUpperThreshold":
		return "were outside"
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

// evaluatorStoresForRegion returns the alarm store and metric store for
// the specified region.
func (s *CloudWatchService) evaluatorStoresForRegion(region string) (*cwstore.AlarmStore, *cwstore.MetricChunkStore) {
	if region == "" {
		return nil, nil
	}

	if cached, ok := s.stores.Load(region); ok {
		if typed, ok := cached.(*cloudwatchStores); ok {
			return typed.alarms, typed.metrics
		}
	}

	var storage storage.BasicStorage
	if s.storageManager != nil {
		var err error
		storage, err = s.storageManager.GetStorage(region)
		if err != nil {
			return nil, nil
		}
	}
	if storage == nil {
		return nil, nil
	}

	alarmStore := cwstore.NewAlarmStore(storage, s.accountID, region)
	metricStore, err := cwstore.NewMetricChunkStoreWithIndex(storage, region, s.dataPath)
	if err != nil {
		return alarmStore, nil
	}

	stores := &cloudwatchStores{
		metrics:    metricStore,
		alarms:     alarmStore,
		dashboards: cwstore.NewDashboardStore(storage, s.accountID, region),
	}
	if actual, loaded := s.stores.LoadOrStore(region, stores); loaded {
		metricStore.Close()
		if typed, ok := actual.(*cloudwatchStores); ok {
			return typed.alarms, typed.metrics
		}
	}

	return alarmStore, metricStore
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

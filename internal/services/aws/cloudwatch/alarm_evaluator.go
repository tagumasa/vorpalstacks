package cloudwatch

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"
	"vorpalstacks/internal/core/storage"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
)

const (
	// defaultEvalInterval is the default tick interval for the alarm
	// evaluator loop. Alarms with a Period shorter than the tick are
	// evaluated once per completed period boundary (see
	// evaluateAlarmJob), so sub-minute periods do not miss windows.
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

	// subMinuteMu guards lastEvaluated, the per-alarm watermark of the
	// last completed period boundary already evaluated. It is only used
	// for alarms whose Period is shorter than the tick interval, where a
	// single window per tick would skip periods.
	subMinuteMu   sync.Mutex
	lastEvaluated map[string]time.Time
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
		interval:      interval,
		workers:       workers,
		logger:        logger,
		lastEvaluated: make(map[string]time.Time),
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
	e.pruneWatermarks(region, alarms)

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
		e.evaluateMetricAlarms(ctx, s, region, metricAlarms, metricStore, alarmStore, muteRuleStore)
	}

	if len(compositeAlarms) > 0 {
		e.evaluateCompositeAlarms(ctx, s, compositeAlarms, alarmStore, muteRuleStore)
	}
}

func (e *alarmEvaluator) evaluateMetricAlarms(ctx context.Context, s *CloudWatchService, region string, alarms []*cwstore.Alarm, metricStore *cwstore.MetricChunkStore, alarmStore *cwstore.AlarmStore, muteRuleStore *cwstore.AlarmMuteRuleStore) {
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
				result := e.evaluateAlarmJob(region, job.alarm, metricStore)
				if result == nil {
					continue
				}
				s.handleAlarmStateTransition(ctx, result, alarmStore, muteRuleStore)
			}
		}()
	}
	wg.Wait()
}

// watermarkKey scopes a sub-minute evaluation watermark to the alarm's
// region: alarm names are unique per region only, so a name-only key would
// let same-named alarms in different regions starve each other.
func watermarkKey(region, alarmName string) string {
	return region + "/" + alarmName
}

// pruneWatermarks drops this region's sub-minute evaluation watermarks for
// alarms that no longer exist, so the map stays bounded and a recreated
// alarm does not resume from a stale boundary.
func (e *alarmEvaluator) pruneWatermarks(region string, alarms []*cwstore.Alarm) {
	live := make(map[string]bool, len(alarms))
	for _, a := range alarms {
		live[watermarkKey(region, a.Name)] = true
	}
	prefix := region + "/"
	e.subMinuteMu.Lock()
	for k := range e.lastEvaluated {
		if strings.HasPrefix(k, prefix) && !live[k] {
			delete(e.lastEvaluated, k)
		}
	}
	e.subMinuteMu.Unlock()
}

// evaluateAlarmJob evaluates one alarm per tick. Alarms whose Period is at
// least the tick interval evaluate the most recent completed window. Alarms
// with a shorter Period (10, 20 or 30 seconds against the 60-second default
// tick) get one evaluation per completed period boundary since their last
// evaluation, mirroring AWS's per-period evaluation, so no period window is
// skipped. It returns on the first state transition; the watermark advances
// to the boundary that produced it.
func (e *alarmEvaluator) evaluateAlarmJob(region string, alarm *cwstore.Alarm, metricStore *cwstore.MetricChunkStore) *alarmEvalResult {
	if alarm.Period <= 0 || time.Duration(alarm.Period)*time.Second >= e.interval {
		return evaluateAlarm(alarm, metricStore)
	}

	period := time.Duration(alarm.Period) * time.Second
	now := time.Now().UTC()
	endTime := now.Truncate(period)
	key := watermarkKey(region, alarm.Name)

	e.subMinuteMu.Lock()
	last, tracked := e.lastEvaluated[key]
	e.subMinuteMu.Unlock()

	if tracked && !last.Before(endTime) {
		// No new period boundary completed since the last tick.
		return nil
	}
	if !tracked {
		// First evaluation after startup: cover only the most recent
		// window; older boundaries predate this process's watch.
		result := evaluateAlarmWindow(alarm, metricStore, endTime)
		e.setWatermark(key, endTime)
		return result
	}

	for b := last.Add(period); !b.After(endTime); b = b.Add(period) {
		result := evaluateAlarmWindow(alarm, metricStore, b)
		if result != nil {
			e.setWatermark(key, b)
			return result
		}
	}
	e.setWatermark(key, endTime)
	return nil
}

func (e *alarmEvaluator) setWatermark(alarmName string, boundary time.Time) {
	e.subMinuteMu.Lock()
	e.lastEvaluated[alarmName] = boundary
	e.subMinuteMu.Unlock()
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
	// evaluated in a single pass. Nodes on (or downstream of) a cycle are
	// returned separately: creation-time validation should keep cycles
	// out, but a stray cycle must only skip its own alarms, not abort the
	// whole region's composite evaluation.
	ordered, cyclic := topologicalSortLevels(dependencies)
	if len(cyclic) > 0 {
		sort.Strings(cyclic)
		e.log("skipping composite alarms in a dependency cycle", "alarms", cyclic)
	}
	cyclicSet := make(map[string]bool, len(cyclic))
	for _, name := range cyclic {
		cyclicSet[name] = true
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
			if cyclicSet[name] {
				continue
			}
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
			// Refresh the shared state map so composites on deeper
			// levels evaluate against this level's fresh state instead
			// of the pre-tick snapshot; otherwise a multi-level
			// escalation always lags one tick behind.
			alarmStateMap[name] = result.newState
		}
	}
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

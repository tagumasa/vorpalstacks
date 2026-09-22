// Copyright 2026 Vorpalstacks Authors
// SPDX-License-Identifier: Apache-2.0

package cloudwatchlogs

import (
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/worker"
	"vorpalstacks/internal/core/logs"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	"vorpalstacks/pkg/filterpattern"
)

// The metric-filter emission module: how an ingested batch becomes
// CloudWatch metric data. The documented semantics ("CloudWatch aggregates
// and reports metric values every minute"; "If your metric filter doesn't
// find matches in either records during the second minute, the default
// value for that minute is 0") decompose into three rules the emitter
// follows:
//   - every matched event publishes the transformation's metric value,
//     resolved as a number or as a value reference into the event;
//   - a matchless evaluation publishes nothing by itself: the minute may
//     still observe a later match, so the default is decided only once
//     the minute has closed without one — the next evaluation carries
//     the close in-band and the minute-close worker carries it when the
//     group goes quiet;
//   - a minute that observed a match publishes no default at all, so the
//     default replaces the value only for entirely matchless minutes.

// evaluateMetricFilters evaluates every metric filter of the log group
// against the ingested events and publishes the matching metric data.
// Listing rides the fetch-all layer so the group's complete filter set
// evaluates; a listing failure surfaces through the server log instead of
// silently disabling metric emission.
func (s *LogsService) evaluateMetricFilters(store *logsstore.Store, region, logGroupName, logStream string, events []logsstore.LogEntry, transformed map[string]string) {
	if s.metricInvoker() == nil {
		logs.Warn("cwMetricInvoker is nil, skipping metric filter evaluation",
			logs.String("logGroup", logGroupName))
		return
	}

	filters, err := fetchAllMetricFilters(store, logGroupName, "")
	if err != nil {
		logs.Error("Failed to list metric filters for evaluation",
			logs.String("logGroup", logGroupName), logs.Err(err))
		return
	}
	if len(filters) == 0 {
		return
	}

	matcher := filterpattern.NewMatcher()
	now := time.Now()

	for _, filter := range filters {
		// The criteria member gates which events the filter processes
		// ("specifies which log events should be processed by this metric
		// filter based on system fields such as source account and source
		// region") — a filter whose criteria the batch's system fields
		// fail never evaluates.
		if !evalFieldSelectionCriteria(filter.FieldSelectionCriteria, region, s.accountID) {
			continue
		}
		anyMatch := false
		for _, event := range events {
			// A filter with applyOnTransformedLogs set evaluates the
			// transformed form when one exists (falling back to the
			// original for events the recipe left untransformed).
			message := transformedForFilters(event.Message,
				logsstore.TransformedMessageDigest(event.Timestamp, logStream, event.Message), transformed, filter.ApplyOnTransformedLogs)
			matched, extracted := matcher.ExtractMatches(filter.FilterPattern, message)
			if !matched {
				continue
			}
			anyMatch = true
			ts := time.UnixMilli(event.Timestamp)
			if ts.IsZero() || ts.After(now) {
				ts = now
			}
			for _, transform := range filter.MetricTransformations {
				value, resolvable := resolveMetricValue(transform.MetricValue, extracted)
				if !resolvable {
					continue
				}
				dimensions, ok := s.resolveTransformDimensions(transform, filter.EmitSystemFieldDimensions, region, extracted)
				if !ok {
					continue
				}
				if dimensions == nil && transform.Unit == "" {
					if err := s.metricInvoker().PutMetricData(region, transform.MetricNamespace, transform.MetricName, value, ts); err != nil {
						logs.Error("Failed to put metric data", logs.Err(err))
					}
					continue
				}
				if err := s.metricInvoker().PutMetricDataWithDimensionsAndUnit(region, transform.MetricNamespace, transform.MetricName, dimensions, transform.Unit, value, ts); err != nil {
					logs.Error("Failed to put metric data", logs.Err(err))
				}
			}
		}
		// The default belongs to the closed minute, never to the running
		// one: a matchless batch at :05 followed by a matched batch at
		// :40 publishes the minute's values alone. The close publishes
		// the previous minute's owed default first, then this batch
		// marks the current minute matched or merely observed.
		s.closeElapsedMetricDefaultMinute(region, logGroupName, filter, now)
		if anyMatch {
			s.markMetricMatched(region, logGroupName, filter.Name, now)
			continue
		}
		s.recordMetricMinutePresence(region, logGroupName, filter.Name, now)
	}
}

// resolveMetricValue resolves a transformation's metricValue against one
// matched event's extracted values. The documented forms are a number and
// a value reference into the event ("You can set metric values to
// numbers, named identifiers, or numeric identifiers"; the worked example
// runs the JSON pattern { $.latency = * } with metricValue $.latency and
// publishes the event's numeric 50). A reference that the event does not
// carry, or whose value is not numeric, resolves to nothing and the event
// contributes no data point. An absent value keeps the historical unit
// increment.
func resolveMetricValue(metricValue string, extracted map[string]string) (float64, bool) {
	v := strings.TrimSpace(metricValue)
	if v == "" {
		return 1, true
	}
	if strings.HasPrefix(v, "$") {
		for _, key := range []string{v, strings.TrimPrefix(v, "$."), strings.TrimPrefix(v, "$")} {
			if key == "" {
				continue
			}
			if raw, ok := extracted[key]; ok {
				if f, err := strconv.ParseFloat(strings.TrimSpace(raw), 64); err == nil {
					return f, true
				}
			}
		}
		return 0, false
	}
	if f, err := strconv.ParseFloat(v, 64); err == nil {
		return f, true
	}
	return 0, false
}

// resolveTransformDimensions resolves one transformation's dimensions for
// a matched event: the values are value references into the event (the
// operation documentation's worked example runs "dimensions":
// {"Request": "$request"}), literals otherwise; the emitSystemFieldDimensions
// members contribute "@aws.account" and "@aws.region" ("A list of system
// fields to emit as additional dimensions in the generated metrics"). A
// reference the event does not carry resolves to nothing and the event
// contributes no data point — the same treatment an unresolvable
// metricValue reference receives. A transformation with no dimensions at
// either source returns a nil map so the undimensioned datum publishes
// through the plain call.
func (s *LogsService) resolveTransformDimensions(transform logsstore.MetricTransformation, emitSystemFieldDimensions []string, region string, extracted map[string]string) (map[string]string, bool) {
	if len(transform.Dimensions) == 0 && len(emitSystemFieldDimensions) == 0 {
		return nil, true
	}
	dims := make(map[string]string, len(transform.Dimensions)+len(emitSystemFieldDimensions))
	for k, v := range transform.Dimensions {
		if strings.HasPrefix(v, "$") {
			resolved, ok := extractStringValue(extracted, v)
			if !ok {
				return nil, false
			}
			dims[k] = resolved
			continue
		}
		dims[k] = v
	}
	for _, field := range emitSystemFieldDimensions {
		switch field {
		case "@aws.account":
			dims[field] = s.accountID
		case "@aws.region":
			dims[field] = region
		}
	}
	return dims, true
}

// extractStringValue resolves one $-prefixed value reference against a
// matched event's extracted values, over the same key ladder
// resolveMetricValue walks.
func extractStringValue(extracted map[string]string, ref string) (string, bool) {
	for _, key := range []string{ref, strings.TrimPrefix(ref, "$."), strings.TrimPrefix(ref, "$")} {
		if key == "" {
			continue
		}
		if v, ok := extracted[key]; ok {
			return v, true
		}
	}
	return "", false
}

// emitTransformationMetrics publishes the documented log transformer
// metrics for one ingested batch: "CloudWatch Logs publishes
// transformation metrics to CloudWatch. These metrics include
// TransformedLogEvents, TransformedBytes, and TransformationErrors" in
// the AWS/Logs namespace — TransformedLogEvents counting the events the
// recipe transformed, TransformedBytes their uncompressed output
// volume, TransformationErrors the events that failed. The dimensions
// follow the level the effective transformer came from: LogGroupname
// "is used only for log-group-level transformers" and PolicyLevel
// "is used only for account-level transformers. Currently the only
// valid value for this dimension is AccountPolicy".
func (s *LogsService) emitTransformationMetrics(region, logGroupName, level string, events, bytes, failures int) {
	if events == 0 && bytes == 0 && failures == 0 {
		return
	}
	invoker := s.metricInvoker()
	if invoker == nil {
		logs.Warn("cwMetricInvoker is nil, skipping transformation metrics",
			logs.String("logGroup", logGroupName))
		return
	}
	dimensions := map[string]string{}
	switch level {
	case transformerLevelGroup:
		dimensions["LogGroupname"] = logGroupName
	case transformerLevelAccount:
		dimensions["PolicyLevel"] = "AccountPolicy"
	default:
		return
	}
	now := time.Now()
	put := func(metricName string, value float64) {
		if err := invoker.PutMetricDataWithDimensions(region, "AWS/Logs", metricName, dimensions, value, now); err != nil {
			logs.Error("Failed to put transformation metric data",
				logs.String("metric", metricName), logs.Err(err))
		}
	}
	if events > 0 {
		put("TransformedLogEvents", float64(events))
	}
	if bytes > 0 {
		put("TransformedBytes", float64(bytes))
	}
	if failures > 0 {
		put("TransformationErrors", float64(failures))
	}
}

// metricMinuteState describes one (log group, filter) pair's current
// wall-clock minute for default-value emission.
type metricMinuteState struct {
	minute         int64
	matched        bool
	defaultEmitted bool
}

// metricDefaultKey is the state map's identity for one filter of one
// group in one region.
func metricDefaultKey(region, logGroupName, filterName string) string {
	return region + "\x00" + logGroupName + "\x00" + filterName
}

// markMetricMatched records that the filter matched during the given
// minute; a minute with a match never reports the default value.
func (s *LogsService) markMetricMatched(region, logGroupName, filterName string, now time.Time) {
	minute := now.Unix() / 60
	s.metricDefaultMu.Lock()
	defer s.metricDefaultMu.Unlock()
	if s.metricDefaultMinutes == nil {
		s.metricDefaultMinutes = make(map[string]metricMinuteState)
	}
	key := metricDefaultKey(region, logGroupName, filterName)
	st := s.metricDefaultMinutes[key]
	if st.minute != minute {
		st = metricMinuteState{minute: minute}
	}
	st.matched = true
	s.metricDefaultMinutes[key] = st
}

// recordMetricMinutePresence registers a matchless evaluation of the
// filter in the given minute without emitting anything: the minute may
// still observe a later match, so its default — if owed at all — is
// decided when the minute closes, never while it is still running.
func (s *LogsService) recordMetricMinutePresence(region, logGroupName, filterName string, now time.Time) {
	minute := now.Unix() / 60
	s.metricDefaultMu.Lock()
	defer s.metricDefaultMu.Unlock()
	if s.metricDefaultMinutes == nil {
		s.metricDefaultMinutes = make(map[string]metricMinuteState)
	}
	key := metricDefaultKey(region, logGroupName, filterName)
	if st, ok := s.metricDefaultMinutes[key]; ok && st.minute == minute {
		return
	}
	s.metricDefaultMinutes[key] = metricMinuteState{minute: minute}
}

// elapsedMetricDefault is one closed matchless minute that owes its
// default emission.
type elapsedMetricDefault struct {
	region     string
	logGroup   string
	filterName string
	minute     int64
}

// claimElapsedMetricDefaultMinutes claims every closed matchless minute
// of the state map — minutes that observed evaluations, no match and no
// emission — marking each claimed so no second claimant can publish it
// again. The minute-close worker and the in-band close share this claim
// seam.
func (s *LogsService) claimElapsedMetricDefaultMinutes(now time.Time) []elapsedMetricDefault {
	minute := now.Unix() / 60
	s.metricDefaultMu.Lock()
	defer s.metricDefaultMu.Unlock()
	var owed []elapsedMetricDefault
	for key, st := range s.metricDefaultMinutes {
		if st.minute >= minute || st.matched || st.defaultEmitted || st.minute == 0 {
			continue
		}
		parts := strings.SplitN(key, "\x00", 3)
		if len(parts) != 3 {
			continue
		}
		st.defaultEmitted = true
		s.metricDefaultMinutes[key] = st
		owed = append(owed, elapsedMetricDefault{
			region: parts[0], logGroup: parts[1], filterName: parts[2], minute: st.minute,
		})
	}
	return owed
}

// claimElapsedMetricDefaultMinute claims one filter's just-ended
// matchless minute — minutes that observed evaluations, no match and no
// emission — marking it claimed so no second claimant (the in-band close
// or the minute-close worker racing it) can publish the minute twice.
func (s *LogsService) claimElapsedMetricDefaultMinute(region, logGroupName, filterName string, now time.Time) (int64, bool) {
	minute := now.Unix() / 60
	s.metricDefaultMu.Lock()
	defer s.metricDefaultMu.Unlock()
	key := metricDefaultKey(region, logGroupName, filterName)
	st := s.metricDefaultMinutes[key]
	if st.minute == 0 || st.minute >= minute || st.matched || st.defaultEmitted {
		return 0, false
	}
	st.defaultEmitted = true
	s.metricDefaultMinutes[key] = st
	return st.minute, true
}

// closeElapsedMetricDefaultMinute publishes the default value a just
// ended matchless minute of this filter owes — the next evaluation
// carries the close in-band when the group keeps receiving batches.
func (s *LogsService) closeElapsedMetricDefaultMinute(region, logGroupName string, filter *logsstore.MetricFilter, now time.Time) {
	minute, owed := s.claimElapsedMetricDefaultMinute(region, logGroupName, filter.Name, now)
	if !owed {
		return
	}
	s.publishMetricDefaultForMinute(region, filter, minute)
}

// publishMetricDefaultForMinute publishes one filter's transformations'
// default values timestamped at the given minute's start. Dimensions
// never ride the default path — the conflict rule ("If you assign
// dimensions to a metric created by a metric filter, you can't assign a
// default value for that metric") leaves the default-carrying
// transformation dimensionless — but the unit still names the metric's
// unit. A failed put is logged and the minute's default is lost: the
// matched-value path treats its own put failures the same way, and the
// minute has closed — the state has moved on, so no later evaluation
// can retry it.
func (s *LogsService) publishMetricDefaultForMinute(region string, filter *logsstore.MetricFilter, minute int64) {
	ts := time.Unix(minute*60, 0)
	for _, transform := range filter.MetricTransformations {
		if !transform.DefaultValueSet {
			continue
		}
		if transform.Unit != "" {
			if err := s.metricInvoker().PutMetricDataWithDimensionsAndUnit(region, transform.MetricNamespace, transform.MetricName, nil, transform.Unit, transform.DefaultValue, ts); err != nil {
				logs.Error("Failed to put metric data", logs.Err(err))
			}
			continue
		}
		if err := s.metricInvoker().PutMetricData(region, transform.MetricNamespace, transform.MetricName, transform.DefaultValue, ts); err != nil {
			logs.Error("Failed to put metric data", logs.Err(err))
		}
	}
}

// closeQuietMetricDefaultMinutes closes finished matchless minutes whose
// groups went quiet: a minute that observed evaluations but no match
// owes its default value even when no later batch arrives to close it
// in-band. Each owed group's filters are fetched once and the defaults
// publish at their minutes' own timestamps.
func (s *LogsService) closeQuietMetricDefaultMinutes() {
	if s.metricInvoker() == nil {
		return
	}
	owed := s.claimElapsedMetricDefaultMinutes(time.Now())
	byGroup := make(map[string][]elapsedMetricDefault)
	for _, o := range owed {
		gk := o.region + "\x00" + o.logGroup
		byGroup[gk] = append(byGroup[gk], o)
	}
	for gk, items := range byGroup {
		parts := strings.SplitN(gk, "\x00", 2)
		if len(parts) != 2 {
			continue
		}
		region, group := parts[0], parts[1]
		store, err := s.getLogsStoreByRegion(region)
		if err != nil {
			logs.Error("Failed to resolve store for metric default close",
				logs.String("region", region), logs.Err(err))
			continue
		}
		filters, err := fetchAllMetricFilters(store, group, "")
		if err != nil {
			logs.Error("Failed to list filters for metric default close",
				logs.String("logGroup", group), logs.Err(err))
			continue
		}
		byName := make(map[string]*logsstore.MetricFilter, len(filters))
		for _, f := range filters {
			byName[f.Name] = f
		}
		for _, item := range items {
			// A filter deleted before its minute closed owes nothing.
			if f, ok := byName[item.filterName]; ok {
				s.publishMetricDefaultForMinute(region, f, item.minute)
			}
		}
	}
}

// metricDefaultMinuteCloseTick bounds how long a closed matchless minute
// waits for the worker when its group receives no further batch.
// TEST_MODE shortens it so SDK-level default pins observe the close in
// seconds.
var metricDefaultMinuteCloseTick = worker.Cadence(30*time.Second, time.Second)

// startMetricDefaultMinuteCloser runs the minute-close worker under the
// service's panic-respawn policy.
func (s *LogsService) startMetricDefaultMinuteCloser() {
	s.startRespawningWorker("metric default minute closer",
		worker.TickerLoop(s.ctx.Done(), metricDefaultMinuteCloseTick, s.closeQuietMetricDefaultMinutes))
}

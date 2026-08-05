package eventbridge

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// schedulerTickInterval is the granularity at which the scheduler evaluates
// rules. AWS EventBridge uses 1-minute minimum granularity for rate/cron.
const schedulerTickInterval = 1 * time.Minute

// lastFireTimes tracks the last time each scheduled rule was fired.
// Keyed by rule ARN. Accessed only by the single scheduler goroutine.
var lastFireTimes sync.Map

// startScheduler launches a background goroutine that ticks every minute and
// fires ENABLED rules whose ScheduleExpression matches the current time.
// It also starts a retention worker that purges expired archive events hourly.
func (s *EventsService) startScheduler() {
	ctx, cancel := context.WithCancel(context.Background())
	s.schedCancel = cancel
	s.schedWg.Add(2)
	go func() {
		defer s.schedWg.Done()
		s.runScheduler(ctx)
	}()
	go func() {
		defer s.schedWg.Done()
		s.runRetentionWorker(ctx)
	}()
}

func (s *EventsService) runScheduler(ctx context.Context) {
	// Align the first tick to the next UTC minute boundary so that
	// rate/cron expressions are evaluated on the same boundary as AWS
	// EventBridge, regardless of when the server started.
	now := time.Now().UTC()
	elapsed := time.Duration(now.Second())*time.Second + time.Duration(now.Nanosecond())
	wait := schedulerTickInterval - elapsed
	if wait > 0 {
		select {
		case <-time.After(wait):
		case <-ctx.Done():
			return
		}
	}

	ticker := time.NewTicker(schedulerTickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			s.tickScheduledRules(ctx, now.UTC())
		}
	}
}

func (s *EventsService) tickScheduledRules(ctx context.Context, now time.Time) {
	s.eventsStores.Range(func(key, value any) bool {
		region, _ := key.(string)
		store, ok := value.(*eventsstore.EventsStore)
		if !ok {
			return true
		}
		s.fireScheduledRulesForRegion(ctx, region, store, now)
		return true
	})
}

func (s *EventsService) fireScheduledRulesForRegion(ctx context.Context, region string, store *eventsstore.EventsStore, now time.Time) {
	nextToken := ""
	for {
		result, err := store.ListRules(ctx, "", "", 1000, nextToken)
		if err != nil {
			logs.Warn("eventbridge scheduler: failed to list rules",
				logs.String("region", region),
				logs.Err(err))
			return
		}
		for _, rule := range result.Rules {
			if rule.State != eventsstore.RuleStateEnabled && rule.State != eventsstore.RuleStateEnabledWithAllCloudtrailManagementEvents {
				continue
			}
			if rule.ScheduleExpression == "" {
				continue
			}
			if !shouldFireSchedule(rule.ARN, rule.ScheduleExpression, now) {
				continue
			}
			s.fireScheduledRule(ctx, region, store, rule, now)
		}
		if result.NextToken == "" {
			break
		}
		nextToken = result.NextToken
	}
}

func (s *EventsService) fireScheduledRule(ctx context.Context, region string, store *eventsstore.EventsStore, rule *eventsstore.Rule, now time.Time) {
	logs.Debug("eventbridge scheduler: firing scheduled rule",
		logs.String("rule", rule.Name),
		logs.String("region", region),
		logs.String("schedule", rule.ScheduleExpression))

	event := &eventsstore.Event{
		ID:           generateEventID(),
		Version:      "0",
		DetailType:   "Scheduled Event",
		Source:       "aws.events",
		Account:      s.accountID,
		Time:         now,
		Region:       region,
		Detail:       map[string]interface{}{},
		EventBusName: rule.EventBusName,
	}

	if err := s.deliverEventWithStore(ctx, region, event, rule.EventBusName, store); err != nil {
		logs.Warn("eventbridge scheduler: failed to deliver scheduled event",
			logs.String("rule", rule.Name),
			logs.String("region", region),
			logs.Err(err))
	}
}

// shouldFireSchedule determines whether a schedule expression should fire at
// the given time. It uses lastFireTimes to ensure each rule fires at most once
// per evaluation.
func shouldFireSchedule(ruleARN, expr string, now time.Time) bool {
	if strings.HasPrefix(expr, "rate(") {
		return shouldFireRate(ruleARN, expr, now)
	}
	if strings.HasPrefix(expr, "cron(") {
		return shouldFireCron(ruleARN, expr, now)
	}
	return false
}

// shouldFireRate parses rate(value unit) expressions and fires when the
// elapsed time since the last fire is >= the rate duration.
func shouldFireRate(ruleARN, expr string, now time.Time) bool {
	inner := strings.TrimPrefix(expr, "rate(")
	inner = strings.TrimSuffix(inner, ")")
	inner = strings.TrimSpace(inner)

	parts := strings.Fields(inner)
	if len(parts) != 2 {
		return false
	}
	value, err := strconv.Atoi(parts[0])
	if err != nil || value <= 0 {
		return false
	}

	var duration time.Duration
	switch unit := strings.ToLower(parts[1]); unit {
	case "minute", "minutes":
		duration = time.Duration(value) * time.Minute
	case "hour", "hours":
		duration = time.Duration(value) * time.Hour
	case "day", "days":
		duration = time.Duration(value) * 24 * time.Hour
	case "week", "weeks":
		duration = time.Duration(value) * 7 * 24 * time.Hour
	default:
		return false
	}

	last, ok := getLastFire(ruleARN)
	if !ok {
		setLastFire(ruleARN, now)
		return true
	}
	if now.Sub(last) >= duration {
		setLastFire(ruleARN, now)
		return true
	}
	return false
}

// shouldFireCron evaluates a cron(...) expression against the current time.
// AWS cron format: cron(minutes hours day-of-month month day-of-week year)
func shouldFireCron(ruleARN, expr string, now time.Time) bool {
	inner := strings.TrimPrefix(expr, "cron(")
	inner = strings.TrimSuffix(inner, ")")
	inner = strings.TrimSpace(inner)

	fields := strings.Fields(inner)
	if len(fields) < 5 || len(fields) > 6 {
		return false
	}

	if !cronFieldMatches(fields[0], now.Minute(), 0, 59) {
		return false
	}
	if !cronFieldMatches(fields[1], now.Hour(), 0, 23) {
		return false
	}
	if !cronFieldMatches(fields[2], now.Day(), 1, 31) {
		return false
	}
	if !cronMonthMatches(fields[3], now.Month()) {
		return false
	}
	if !cronDOWMatches(fields[4], int(now.Weekday())) {
		return false
	}

	// Year field (optional, defaults to *)
	if len(fields) == 6 && !cronFieldMatches(fields[5], now.Year(), 1970, 2199) {
		return false
	}

	// Ensure we only fire once per minute.
	last, ok := getLastFire(ruleARN)
	if ok && last.Truncate(time.Minute).Equal(now.Truncate(time.Minute)) {
		return false
	}
	setLastFire(ruleARN, now)
	return true
}

// cronFieldMatches checks whether a cron field matches the given value.
// Supports: * (wildcard), exact values, ranges (a-b), lists (a,b,c),
// and step values (*/n or a-b/n). The ? character is treated as *.
func cronFieldMatches(field string, value, min, max int) bool {
	field = strings.TrimSpace(field)
	if field == "*" || field == "?" {
		return true
	}

	for _, part := range strings.Split(field, ",") {
		part = strings.TrimSpace(part)
		if matchesCronPart(part, value, min, max) {
			return true
		}
	}
	return false
}

func matchesCronPart(part string, value, min, max int) bool {
	// Step value: a/n (a-max/n), a-b/n, or * /n
	if idx := strings.Index(part, "/"); idx != -1 {
		rangePart := part[:idx]
		stepStr := part[idx+1:]
		step, err := strconv.Atoi(stepStr)
		if err != nil || step <= 0 {
			return false
		}
		start, end := min, max
		if rangePart != "*" && rangePart != "?" {
			// Plain numeric a/n: equivalent to a-max/n.
			if !strings.Contains(rangePart, "-") {
				n, err := strconv.Atoi(rangePart)
				if err != nil || n < min || n > max {
					return false
				}
				start = n
			} else {
				s, e, ok := parseRange(rangePart, min, max)
				if !ok {
					return false
				}
				start, end = s, e
			}
		}
		for i := start; i <= end; i += step {
			if i == value {
				return true
			}
		}
		return false
	}

	// Range: a-b
	if strings.Contains(part, "-") {
		s, e, ok := parseRange(part, min, max)
		if !ok {
			return false
		}
		return value >= s && value <= e
	}

	// Exact value
	n, err := strconv.Atoi(part)
	if err != nil {
		return false
	}
	return n == value
}

func parseRange(s string, min, max int) (int, int, bool) {
	parts := strings.SplitN(s, "-", 2)
	if len(parts) != 2 {
		return 0, 0, false
	}
	start, err1 := strconv.Atoi(parts[0])
	end, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return 0, 0, false
	}
	if start < min {
		start = min
	}
	if end > max {
		end = max
	}
	return start, end, true
}

var monthMap = map[string]int{
	"JAN": 1, "FEB": 2, "MAR": 3, "APR": 4, "MAY": 5, "JUN": 6,
	"JUL": 7, "AUG": 8, "SEP": 9, "OCT": 10, "NOV": 11, "DEC": 12,
}

// expandCronNames returns a copy of field in which every occurrence of a
// month name (JAN..DEC) is replaced with its numeric value. The replacement
// is applied to every comma-separated, range, list, and step sub-expression
// so that downstream range/list/step validators can treat the field purely
// as integers.
func expandCronNames(field string, names map[string]int) string {
	upper := strings.ToUpper(field)
	for name, num := range names {
		upper = strings.ReplaceAll(upper, name, strconv.Itoa(num))
	}
	return upper
}

func cronMonthMatches(field string, month time.Month) bool {
	expanded := expandCronNames(field, monthMap)
	return cronFieldMatches(expanded, int(month), 1, 12)
}

var dowNameToAWS = map[string]int{
	"SUN": 1, "MON": 2, "TUE": 3, "WED": 4, "THU": 5, "FRI": 6, "SAT": 7,
}

// cronDOWMatches checks the day-of-week field. AWS cron uses 1=SUN..7=SAT.
// Go's time.Weekday uses 0=SUNDAY..6=SATURDAY. We convert Go's value to
// AWS convention (goDOW+1) before matching.
//
// Both cronMonthMatches and cronDOWMatches expand the supplied field so that
// names inside ranges (e.g. MON-FRI), lists (MON,WED,FRI), and steps
// (MON/2) are translated to numeric form before cronFieldMatches parses
// them. This closes the long-standing asymmetry where only the bare single
// name was being recognised.
func cronDOWMatches(field string, goDOW int) bool {
	awsDOW := goDOW + 1
	if awsDOW > 7 {
		awsDOW = 7
	}

	expanded := expandCronNames(field, dowNameToAWS)
	if expanded == "?" || expanded == "*" {
		return true
	}
	return cronFieldMatches(expanded, awsDOW, 1, 7)
}

func getLastFire(ruleARN string) (time.Time, bool) {
	v, ok := lastFireTimes.Load(ruleARN)
	if !ok {
		return time.Time{}, false
	}
	t, ok := v.(time.Time)
	return t, ok
}

func setLastFire(ruleARN string, t time.Time) {
	lastFireTimes.Store(ruleARN, t)
}

// retentionTickInterval controls how often the retention worker runs.
const retentionTickInterval = 1 * time.Hour

func (s *EventsService) runRetentionWorker(ctx context.Context) {
	ticker := time.NewTicker(retentionTickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			s.purgeExpiredArchiveEvents(ctx, now.UTC())
		}
	}
}

func (s *EventsService) purgeExpiredArchiveEvents(ctx context.Context, now time.Time) {
	s.eventsStores.Range(func(key, value any) bool {
		store, ok := value.(*eventsstore.EventsStore)
		if !ok {
			return true
		}
		token := ""
		for {
			result, err := store.ListArchives(ctx, "", "", "", 1000, token)
			if err != nil {
				return true
			}
			for _, archive := range result.Archives {
				if archive.RetentionDays <= 0 {
					continue
				}
				cutoff := now.AddDate(0, 0, -int(archive.RetentionDays))
				if err := store.DeleteExpiredArchiveEvents(ctx, archive.Name, cutoff); err != nil {
					logs.Warn("eventbridge retention: failed to purge expired events",
						logs.String("archive", archive.Name),
						logs.Err(err))
				}
			}
			if result.NextToken == "" {
				break
			}
			token = result.NextToken
		}
		return true
	})
}

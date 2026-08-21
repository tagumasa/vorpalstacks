// Package scheduleexpr provides AWS schedule expression parsing for
// cron(), rate(), and at() expressions. It is shared by the EventBridge
// Scheduler engine and the Timestream Query scheduled-query engine.
package scheduleexpr

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/utils/timeutils"

	"github.com/robfig/cron/v3"
)

var (
	atExprPattern   = regexp.MustCompile(`^at\((\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})\)$`)
	rateExprPattern = regexp.MustCompile(`^rate\((\d+)\s+(minute|minutes|hour|hours|day|days)\)$`)
	cronExprPattern = regexp.MustCompile(`^cron\((.+)\)$`)
)

// Rate expression value bounds enforced by validateRateFormat: the value
// may not exceed one year expressed in the chosen unit.
const (
	// MaxRateMinutes is the largest accepted rate() value in minutes.
	MaxRateMinutes = 525600
	// MaxRateHours is the largest accepted rate() value in hours.
	MaxRateHours = 8760
	// MaxRateDays is the largest accepted rate() value in days.
	MaxRateDays = 365
)

// ValidateExpression checks whether an AWS schedule expression is
// syntactically valid. It accepts cron(), rate(), and at() expressions
// following the AWS EventBridge / Scheduler format rules:
//   - Overall length must not exceed 256 characters.
//   - at(): ISO-8601 timestamp with a valid calendar date.
//   - rate(): positive value with a unit that agrees with the value
//     (singular for 1, plural for values greater than 1) and at most one
//     year expressed in the chosen unit (MaxRateMinutes, MaxRateHours or
//     MaxRateDays).
//   - cron(): exactly 6 whitespace-delimited fields inside the parentheses.
//
// This is a format check only; it does not guarantee that the resulting
// schedule will ever fire. Use NextExecutionTime to compute actual fire
// times.
func ValidateExpression(expr string) bool {
	if len(expr) > 256 {
		return false
	}

	if matches := atExprPattern.FindStringSubmatch(expr); len(matches) == 2 {
		if _, err := time.Parse(timeutils.ISO8601NoZFormat, matches[1]); err != nil {
			return false
		}
		return true
	}

	if validateRateFormat(expr) {
		return true
	}

	if matches := cronExprPattern.FindStringSubmatch(expr); len(matches) == 2 {
		fields := strings.Fields(matches[1])
		if len(fields) != 6 {
			return false
		}
		return true
	}

	return false
}

// validateRateFormat checks a rate() expression: the value must be a
// positive number (>= 1), the unit must agree with the value — singular
// for 1, plural for values greater than 1 — and the value must not
// exceed one year expressed in the chosen unit.
func validateRateFormat(expr string) bool {
	matches := rateExprPattern.FindStringSubmatch(expr)
	if len(matches) != 3 {
		return false
	}
	if !rateValueAgrees(matches[1], matches[2]) {
		return false
	}
	value, _ := strconv.Atoi(matches[1])
	// Agreement has already been checked, so a singular unit implies
	// value == 1, which is within every bound.
	switch matches[2] {
	case "minute", "minutes":
		return value <= MaxRateMinutes
	case "hour", "hours":
		return value <= MaxRateHours
	case "day", "days":
		return value <= MaxRateDays
	}
	return false
}

// NextExecutionTime calculates the next execution time for an AWS schedule
// expression (cron(...), rate(...), or at(...)) at or after the given
// reference time. creationTime is the baseline for rate() calculations.
// startDate, if non-nil, overrides creationTime as the rate() baseline.
func NextExecutionTime(expr string, now time.Time, creationTime time.Time, startDate *time.Time) (time.Time, error) {
	if strings.HasPrefix(expr, "at(") {
		matches := atExprPattern.FindStringSubmatch(expr)
		if len(matches) == 2 {
			t, err := time.Parse(timeutils.ISO8601NoZFormat, matches[1])
			if err != nil {
				return time.Time{}, err
			}
			return t, nil
		}
	}

	if strings.HasPrefix(expr, "rate(") {
		matches := rateExprPattern.FindStringSubmatch(expr)
		if len(matches) == 3 {
			value, _ := strconv.Atoi(matches[1])
			unit := matches[2]

			var duration time.Duration
			switch strings.TrimSuffix(unit, "s") {
			case "minute":
				duration = time.Duration(value) * time.Minute
			case "hour":
				duration = time.Duration(value) * time.Hour
			case "day":
				duration = time.Duration(value) * 24 * time.Hour
			}

			base := creationTime
			if startDate != nil {
				base = *startDate
			}
			elapsed := now.Sub(base)
			periods := int(elapsed / duration)
			nextTime := base.Add(time.Duration(periods) * duration)
			return nextTime, nil
		}
	}

	if strings.HasPrefix(expr, "cron(") {
		return parseCronNextTime(expr, now)
	}

	return time.Time{}, fmt.Errorf("unsupported schedule expression: %s", expr)
}

func parseCronNextTime(expr string, now time.Time) (time.Time, error) {
	matches := cronExprPattern.FindStringSubmatch(expr)
	if len(matches) != 2 {
		return time.Time{}, fmt.Errorf("invalid cron expression: %s", expr)
	}

	fields := strings.Fields(matches[1])
	if len(fields) != 6 {
		return time.Time{}, fmt.Errorf("AWS cron expression must have 6 fields, got %d", len(fields))
	}

	needsManualEval := false
	for _, f := range fields {
		if strings.Contains(f, "L") || strings.Contains(f, "W") || strings.Contains(f, "#") {
			needsManualEval = true
			break
		}
	}

	if needsManualEval {
		return awsCronNextTime(fields, now)
	}

	standardCron := fmt.Sprintf("%s %s %s %s %s",
		convertAWSCronField(fields[0]),
		convertAWSCronField(fields[1]),
		convertAWSCronField(fields[2]),
		convertAWSCronField(fields[3]),
		convertAWSDowField(fields[4]),
	)

	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(standardCron)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse cron expression: %w", err)
	}

	nextTime := schedule.Next(now.Add(-time.Minute))

	yearField := convertAWSCronField(fields[5])
	if yearField != "*" && yearField != "" {
		for y := nextTime.Year(); y <= now.Year()+100; y++ {
			if isYearAllowed(yearField, y) {
				if y == nextTime.Year() {
					return nextTime, nil
				}
				baseline := time.Date(y, 1, 1, 0, 0, 0, 0, now.Location())
				return schedule.Next(baseline.Add(-time.Minute)), nil
			}
		}
		return time.Time{}, fmt.Errorf("cron year field %q has no valid year in range", fields[5])
	}

	return nextTime, nil
}

func awsCronNextTime(fields []string, now time.Time) (time.Time, error) {
	minuteMatcher, err := buildCronFieldMatcher(fields[0], 0, 59)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid cron minutes field %q: %w", fields[0], err)
	}
	hourMatcher, err := buildCronFieldMatcher(fields[1], 0, 23)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid cron hours field %q: %w", fields[1], err)
	}
	domMatcher, err := buildDayOfMonthMatcher(fields[2])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid cron day-of-month field %q: %w", fields[2], err)
	}
	monthMatcher, err := buildCronFieldMatcher(convertAWSCronField(fields[3]), 1, 12)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid cron month field %q: %w", fields[3], err)
	}
	dowMatcher, err := buildDayOfWeekMatcher(fields[4])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid cron day-of-week field %q: %w", fields[4], err)
	}
	yearMatcher, err := buildCronFieldMatcher(convertAWSCronField(fields[5]), 1970, 2199)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid cron year field %q: %w", fields[5], err)
	}

	domIsAny := fields[2] == "*" || fields[2] == "?"
	dowIsAny := fields[4] == "*" || fields[4] == "?"

	candidate := now.Truncate(time.Minute)
	limit := candidate.AddDate(1, 0, 0)

	for candidate.Before(limit) {
		if !minuteMatcher(candidate.Minute()) {
			candidate = candidate.Add(time.Minute)
			continue
		}
		if !hourMatcher(candidate.Hour()) {
			candidate = time.Date(candidate.Year(), candidate.Month(), candidate.Day(), candidate.Hour()+1, 0, 0, 0, candidate.Location())
			continue
		}
		if !monthMatcher(int(candidate.Month())) {
			candidate = time.Date(candidate.Year(), candidate.Month()+1, 1, 0, 0, 0, 0, candidate.Location())
			continue
		}
		if !yearMatcher(candidate.Year()) {
			candidate = time.Date(candidate.Year()+1, 1, 1, 0, 0, 0, 0, candidate.Location())
			continue
		}

		dayMatches := false
		if domIsAny && dowIsAny {
			dayMatches = true
		} else if domIsAny {
			dayMatches = dowMatcher(candidate)
		} else if dowIsAny {
			dayMatches = domMatcher(candidate)
		} else {
			dayMatches = domMatcher(candidate) || dowMatcher(candidate)
		}

		if !dayMatches {
			candidate = time.Date(candidate.Year(), candidate.Month(), candidate.Day()+1, 0, 0, 0, 0, candidate.Location())
			continue
		}

		return candidate, nil
	}

	return time.Time{}, fmt.Errorf("no matching time found within one year for cron expression")
}

func buildCronFieldMatcher(field string, min, max int) (func(int) bool, error) {
	if field == "*" {
		return func(v int) bool { return v >= min && v <= max }, nil
	}

	allowed := make(map[int]bool)
	for _, part := range strings.Split(field, ",") {
		part = strings.TrimSpace(part)
		if part == "*" {
			for v := min; v <= max; v++ {
				allowed[v] = true
			}
			continue
		}
		if strings.Contains(part, "/") {
			slashParts := strings.SplitN(part, "/", 2)
			base := slashParts[0]
			step, err := strconv.Atoi(slashParts[1])
			if err != nil || step <= 0 {
				return nil, fmt.Errorf("invalid step: %s", slashParts[1])
			}
			start := min
			end := max
			if base != "*" {
				if strings.Contains(base, "-") {
					rangeParts := strings.SplitN(base, "-", 2)
					start, _ = strconv.Atoi(rangeParts[0])
					end, _ = strconv.Atoi(rangeParts[1])
				} else {
					start, _ = strconv.Atoi(base)
				}
			}
			for v := start; v <= end; v += step {
				if v >= min && v <= max {
					allowed[v] = true
				}
			}
			continue
		}
		if strings.Contains(part, "-") {
			rangeParts := strings.SplitN(part, "-", 2)
			start, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				return nil, err
			}
			end, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				return nil, err
			}
			for v := start; v <= end; v++ {
				if v >= min && v <= max {
					allowed[v] = true
				}
			}
			continue
		}
		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		if val >= min && val <= max {
			allowed[val] = true
		}
	}

	return func(v int) bool { return allowed[v] }, nil
}

func buildDayOfMonthMatcher(field string) (func(time.Time) bool, error) {
	converted := convertAWSCronField(field)
	if converted == "*" {
		return func(t time.Time) bool { return true }, nil
	}

	hasL := strings.Contains(field, "L")
	hasW := strings.Contains(field, "W")

	if !hasL && !hasW {
		matcher, err := buildCronFieldMatcher(converted, 1, 31)
		if err != nil {
			return nil, err
		}
		return func(t time.Time) bool { return matcher(t.Day()) }, nil
	}

	return func(t time.Time) bool {
		for _, part := range strings.Split(field, ",") {
			part = strings.TrimSpace(part)
			if part == "*" || part == "?" {
				return true
			}
			if part == "L" {
				lastDay := lastDayOfMonth(t.Year(), t.Month())
				if t.Day() == lastDay {
					return true
				}
				continue
			}
			if strings.HasSuffix(part, "L") {
				offsetStr := strings.TrimSuffix(part, "L")
				offset := 0
				if offsetStr != "" {
					offset, _ = strconv.Atoi(offsetStr)
				}
				lastDay := lastDayOfMonth(t.Year(), t.Month())
				if t.Day() == lastDay+offset {
					return true
				}
				continue
			}
			if strings.HasSuffix(part, "W") {
				dayStr := strings.TrimSuffix(part, "W")
				targetDay, err := strconv.Atoi(dayStr)
				if err != nil {
					continue
				}
				if isNearestWeekday(t, targetDay) {
					return true
				}
				continue
			}
			val, err := strconv.Atoi(convertAWSCronField(part))
			if err == nil && val == t.Day() {
				return true
			}
		}
		return false
	}, nil
}

func buildDayOfWeekMatcher(field string) (func(time.Time) bool, error) {
	if field == "*" || field == "?" {
		return func(t time.Time) bool { return true }, nil
	}

	hasL := strings.Contains(field, "L")
	hasHash := strings.Contains(field, "#")

	if !hasL && !hasHash {
		allowed := make(map[int]bool)
		for _, part := range strings.Split(field, ",") {
			part = strings.TrimSpace(part)
			if part == "*" || part == "?" {
				for i := 0; i < 7; i++ {
					allowed[i] = true
				}
				continue
			}
			step := 1
			base := part
			if idx := strings.Index(part, "/"); idx != -1 {
				base = part[:idx]
				step, _ = strconv.Atoi(part[idx+1:])
				if step <= 0 {
					step = 1
				}
			}
			if idx := strings.Index(base, "-"); idx != -1 {
				lo, err := parseAWSDoWValue(base[:idx])
				if err != nil {
					return nil, err
				}
				hi, err := parseAWSDoWValue(base[idx+1:])
				if err != nil {
					return nil, err
				}
				for d := lo; d <= hi; d += step {
					allowed[awsDowToGo(d)] = true
				}
			} else if base == "*" {
				for d := 1; d <= 7; d += step {
					allowed[awsDowToGo(d)] = true
				}
			} else {
				d, err := parseAWSDoWValue(base)
				if err != nil {
					return nil, err
				}
				allowed[awsDowToGo(d)] = true
			}
		}
		return func(t time.Time) bool { return allowed[int(t.Weekday())] }, nil
	}

	return func(t time.Time) bool {
		for _, part := range strings.Split(field, ",") {
			part = strings.TrimSpace(part)
			if part == "*" || part == "?" {
				return true
			}
			if strings.HasSuffix(part, "L") {
				dayStr := strings.TrimSuffix(part, "L")
				// A bare L in the day-of-week field is the last day of
				// the week — Saturday, day 7 in AWS numbering — and
				// fires every week, unlike nL which is the last given
				// weekday of the month.
				if dayStr == "" {
					if int(t.Weekday()) == awsDowToGo(7) {
						return true
					}
					continue
				}
				dayNum, err := parseAWSDoWValue(dayStr)
				if err != nil {
					continue
				}
				goDay := awsDowToGo(dayNum)
				if isLastWeekdayOfMonth(t, goDay) {
					return true
				}
				continue
			}
			if strings.Contains(part, "#") {
				hashParts := strings.SplitN(part, "#", 2)
				dayNum, err := parseAWSDoWValue(hashParts[0])
				if err != nil {
					continue
				}
				weekNum, err := strconv.Atoi(hashParts[1])
				if err != nil || weekNum < 1 || weekNum > 5 {
					continue
				}
				goDay := awsDowToGo(dayNum)
				if isNthWeekdayOfMonth(t, goDay, weekNum) {
					return true
				}
				continue
			}
			val, err := parseAWSDoWValue(part)
			if err != nil {
				continue
			}
			goDay := awsDowToGo(val)
			if int(t.Weekday()) == goDay {
				return true
			}
		}
		return false
	}, nil
}

func parseAWSDoWValue(s string) (int, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	switch s {
	case "SUN":
		return 1, nil
	case "MON":
		return 2, nil
	case "TUE":
		return 3, nil
	case "WED":
		return 4, nil
	case "THU":
		return 5, nil
	case "FRI":
		return 6, nil
	case "SAT":
		return 7, nil
	case "":
		return 0, fmt.Errorf("empty day-of-week value")
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid day-of-week value: %q", s)
	}
	if n < 1 || n > 7 {
		return 0, fmt.Errorf("day-of-week value must be 1-7, got %d", n)
	}
	return n, nil
}

func awsDowToGo(awsDay int) int {
	if awsDay >= 1 && awsDay <= 7 {
		return awsDay - 1
	}
	return awsDay
}

func lastDayOfMonth(year int, month time.Month) int {
	firstOfNext := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	return firstOfNext.AddDate(0, 0, -1).Day()
}

func isNearestWeekday(t time.Time, targetDay int) bool {
	lastDay := lastDayOfMonth(t.Year(), t.Month())
	if targetDay < 1 || targetDay > lastDay {
		targetDay = lastDay
	}

	target := time.Date(t.Year(), t.Month(), targetDay, 0, 0, 0, 0, t.Location())
	wd := target.Weekday()

	var nearestDay int
	switch wd {
	case time.Saturday:
		if targetDay-1 >= 1 {
			nearestDay = targetDay - 1
		} else {
			nearestDay = targetDay + 2
		}
	case time.Sunday:
		if targetDay+1 <= lastDay {
			nearestDay = targetDay + 1
		} else {
			nearestDay = targetDay - 2
		}
	default:
		nearestDay = targetDay
	}

	return t.Day() == nearestDay
}

func isLastWeekdayOfMonth(t time.Time, goDay int) bool {
	if int(t.Weekday()) != goDay {
		return false
	}
	nextWeek := t.AddDate(0, 0, 7)
	return nextWeek.Month() != t.Month()
}

func isNthWeekdayOfMonth(t time.Time, goDay, weekNum int) bool {
	if int(t.Weekday()) != goDay {
		return false
	}
	weekOfMonth := (t.Day()-1)/7 + 1
	return weekOfMonth == weekNum
}

var awsCronReplacements = []struct{ old, new string }{
	{"?", "*"},
	{"SUN", "0"}, {"MON", "1"}, {"TUE", "2"}, {"WED", "3"}, {"THU", "4"}, {"FRI", "5"}, {"SAT", "6"},
	{"JAN", "1"}, {"FEB", "2"}, {"MAR", "3"}, {"APR", "4"}, {"MAY", "5"}, {"JUN", "6"},
	{"JUL", "7"}, {"AUG", "8"}, {"SEP", "9"}, {"OCT", "10"}, {"NOV", "11"}, {"DEC", "12"},
}

func convertAWSCronField(field string) string {
	for _, r := range awsCronReplacements {
		field = strings.ReplaceAll(field, r.old, r.new)
	}
	return field
}

// convertAWSDowField translates the AWS day-of-week field into the
// numbering the standard cron parser expects: AWS numbers the days 1-7
// (1=Sunday) while the parser uses 0-6 (0=Sunday). Numeric tokens —
// plain values, range endpoints and step bases — are offset by one;
// SUN-SAT names are handled by convertAWSCronField afterwards and need
// no offset. Expressions containing L, W or # never reach this
// converter: they are evaluated by awsCronNextTime instead.
func convertAWSDowField(field string) string {
	offset := make([]string, 0, 4)
	for _, part := range strings.Split(field, ",") {
		offset = append(offset, offsetDowPart(part))
	}
	return convertAWSCronField(strings.Join(offset, ","))
}

// offsetDowPart offsets every numeric token of a single comma element
// (value, range endpoints, step base) from AWS to standard numbering.
// Non-numeric tokens (wildcards, day names) are returned unchanged.
func offsetDowPart(part string) string {
	if idx := strings.Index(part, "/"); idx != -1 {
		return offsetDowPart(part[:idx]) + part[idx:]
	}
	if idx := strings.Index(part, "-"); idx != -1 {
		return offsetDowNum(part[:idx]) + "-" + offsetDowNum(part[idx+1:])
	}
	return offsetDowNum(part)
}

// offsetDowNum maps a single numeric day-of-week token from the AWS
// 1=Sunday..7=Saturday numbering to the standard 0=Sunday..6=Saturday
// numbering. Tokens outside 1-7 are left as they are so malformed input
// fails in the parser rather than being silently rewritten.
func offsetDowNum(token string) string {
	token = strings.TrimSpace(token)
	n, err := strconv.Atoi(token)
	if err != nil || n < 1 || n > 7 {
		return token
	}
	return strconv.Itoa(n - 1)
}

func isYearAllowed(field string, year int) bool {
	if field == "*" || field == "" {
		return true
	}
	now := time.Now().Year()
	for _, part := range strings.Split(field, ",") {
		part = strings.TrimSpace(part)
		if part == "*" {
			return true
		}
		if strings.Contains(part, "/") {
			slashParts := strings.SplitN(part, "/", 2)
			base := slashParts[0]
			step, _ := strconv.Atoi(slashParts[1])
			if step <= 0 {
				step = 1
			}
			startYear := now
			if base != "*" {
				startYear, _ = strconv.Atoi(base)
			}
			if year >= startYear && (year-startYear)%step == 0 {
				return true
			}
			continue
		}
		if strings.Contains(part, "-") {
			rangeParts := strings.SplitN(part, "-", 2)
			startYear, _ := strconv.Atoi(rangeParts[0])
			endYear, _ := strconv.Atoi(rangeParts[1])
			if year >= startYear && year <= endYear {
				return true
			}
			continue
		}
		y, err := strconv.Atoi(part)
		if err == nil && y == year {
			return true
		}
	}
	return false
}

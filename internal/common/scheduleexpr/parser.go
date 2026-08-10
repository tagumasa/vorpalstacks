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
		convertAWSCronField(fields[4]),
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
	if converted == "*" || converted == "?" {
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

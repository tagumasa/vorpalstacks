package scheduleexpr

import (
	"regexp"
	"strconv"
	"strings"
)

// ruleCronExprPattern matches a cron(...) expression whose inner body is
// non-empty (field-level validation follows).
var ruleCronExprPattern = regexp.MustCompile(`^cron\(.+\)$`)

// ruleRateExprPattern matches a rate(value unit) expression. minute,
// hour and day are the only accepted units.
var ruleRateExprPattern = regexp.MustCompile(`^rate\((\d+)\s+(minute|minutes|hour|hours|day|days)\)$`)

// ValidateRuleExpression checks a scheduled-rule expression against the
// Amazon EventBridge PutRule contract — the rule-validation profile of
// this package. It differs from ValidateExpression (the EventBridge
// Scheduler / Timestream profile) in two ways: at() expressions are not
// rules and are rejected, and rate() values carry no upper bound (the
// EventBridge PutRule model specifies only the overall length trait,
// @length(0, 256), which is enforced here). An empty expression is
// valid: the parameter is optional.
func ValidateRuleExpression(expr string) bool {
	if expr == "" {
		return true
	}
	if len(expr) > 256 {
		return false
	}
	if ruleCronExprPattern.MatchString(expr) {
		return validateRuleCronFields(expr)
	}
	if matches := ruleRateExprPattern.FindStringSubmatch(expr); len(matches) == 3 {
		return rateValueAgrees(matches[1], matches[2])
	}
	return false
}

// rateValueAgrees enforces the AWS rate() contract: the value is a
// positive number and the unit agrees in number with the value
// (singular for 1, plural for values greater than 1).
func rateValueAgrees(valueStr, unit string) bool {
	value, err := strconv.Atoi(valueStr)
	if err != nil || value < 1 {
		return false
	}
	if value == 1 {
		switch unit {
		case "minute", "hour", "day":
			return true
		}
		return false
	}
	switch unit {
	case "minutes", "hours", "days":
		return true
	}
	return false
}

// validateRuleCronFields checks the inner contents of cron(...) against
// the AWS six-field layout: minutes hours day-of-month month day-of-week
// year. All six fields are required. Per-field ranges follow the AWS
// EventBridge cron reference. The day-of-month and day-of-week fields are
// mutually exclusive — one of them must be "?" — and "?" is a wildcard in
// those two fields only. Invalid fields (e.g. minute=99) are rejected
// here instead of failing silently at runtime.
func validateRuleCronFields(expr string) bool {
	inner := strings.TrimPrefix(expr, "cron(")
	inner = strings.TrimSuffix(inner, ")")
	inner = strings.TrimSpace(inner)

	fields := strings.Fields(inner)
	if len(fields) != 6 {
		return false
	}

	// Day-of-month and day-of-week cannot both carry values or "*": if a
	// value or "*" is specified in one, the other must be "?".
	if fields[2] != "?" && fields[4] != "?" {
		return false
	}

	// "?" is only a wildcard in the day-of-month and day-of-week fields.
	for i, f := range fields {
		if i != 2 && i != 4 && strings.Contains(f, "?") {
			return false
		}
	}

	normalised := make([]string, len(fields))
	for i, f := range fields {
		normalised[i] = normaliseRuleCronField(i, f)
	}

	if !validateRuleCronField(normalised[0], 0, 59) {
		return false
	}
	if !validateRuleCronField(normalised[1], 0, 23) {
		return false
	}
	if !validateRuleDayOfMonth(normalised[2]) {
		return false
	}
	if !validateRuleCronField(normalised[3], 1, 12) {
		return false
	}
	if !validateRuleDayOfWeek(normalised[4]) {
		return false
	}
	return validateRuleCronField(normalised[5], 1970, 2199)
}

// normaliseRuleCronField replaces month/DOW names with numeric literals
// so the numeric validators can be expressed uniformly. fieldIndex 3 is
// month (JAN..DEC), fieldIndex 4 is day-of-week (SUN..SAT).
func normaliseRuleCronField(fieldIndex int, field string) string {
	upper := strings.ToUpper(field)
	var table map[string]string
	if fieldIndex == 3 {
		table = map[string]string{
			"JAN": "1", "FEB": "2", "MAR": "3", "APR": "4", "MAY": "5", "JUN": "6",
			"JUL": "7", "AUG": "8", "SEP": "9", "OCT": "10", "NOV": "11", "DEC": "12",
		}
	} else if fieldIndex == 4 {
		table = map[string]string{
			"SUN": "1", "MON": "2", "TUE": "3", "WED": "4",
			"THU": "5", "FRI": "6", "SAT": "7",
		}
	} else {
		return field
	}
	for name, num := range table {
		upper = strings.ReplaceAll(upper, name, num)
	}
	return upper
}

// validateRuleCronField verifies that an already name-normalised cron
// field contains only legal values and operators within [min, max]:
// *, ?, exact, list (a,b,c), range (a-b) and step (a/n, a-b/n, * /n,
// ?/n). a/n is equivalent to a-max/n. The day-of-month and day-of-week
// wildcards L, W and # are handled by their dedicated validators below.
func validateRuleCronField(field string, min, max int) bool {
	field = strings.TrimSpace(field)
	if field == "" || field == "*" || field == "?" {
		return true
	}
	for _, part := range strings.Split(field, ",") {
		if !validateRuleCronPart(strings.TrimSpace(part), min, max) {
			return false
		}
	}
	return true
}

// validateRuleDayOfMonth validates the day-of-month field: numeric
// values 1-31 with lists, ranges and steps, plus the AWS wildcards L
// (the last day of the month) and nW (the weekday nearest day n).
func validateRuleDayOfMonth(field string) bool {
	if field == "*" || field == "?" {
		return true
	}
	for _, part := range strings.Split(field, ",") {
		part = strings.TrimSpace(part)
		if part == "L" {
			continue
		}
		if strings.HasSuffix(part, "W") {
			n, err := strconv.Atoi(strings.TrimSuffix(part, "W"))
			if err != nil || n < 1 || n > 31 {
				return false
			}
			continue
		}
		if !validateRuleCronPart(part, 1, 31) {
			return false
		}
	}
	return true
}

// validateRuleDayOfWeek validates the day-of-week field: numeric values
// 1-7 (1=Sunday) with lists, ranges and steps, plus the AWS wildcards L
// (the last day of the week), nL (the last weekday n of the month, as in
// 6L for the last Friday) and n#m (the m-th weekday n of the month, as
// in 3#2 for the second Tuesday). When # is used, the field must be a
// single expression — lists are not valid with it.
func validateRuleDayOfWeek(field string) bool {
	if field == "*" || field == "?" {
		return true
	}
	if strings.Contains(field, "#") {
		parts := strings.Split(field, "#")
		if len(parts) != 2 {
			return false
		}
		day, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil || day < 1 || day > 7 {
			return false
		}
		nth, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		return err == nil && nth >= 1 && nth <= 5
	}
	for _, part := range strings.Split(field, ",") {
		part = strings.TrimSpace(part)
		if part == "L" {
			continue
		}
		if strings.HasSuffix(part, "L") {
			n, err := strconv.Atoi(strings.TrimSuffix(part, "L"))
			if err != nil || n < 1 || n > 7 {
				return false
			}
			continue
		}
		if !validateRuleCronPart(part, 1, 7) {
			return false
		}
	}
	return true
}

// validateRuleCronPart validates a single comma-separated element of a
// cron field: an exact value, a range (a-b) or a step (a/n, a-b/n,
// * /n, ?/n) with all numeric operands within [min, max].
func validateRuleCronPart(part string, min, max int) bool {
	if part == "" {
		return false
	}
	if idx := strings.Index(part, "/"); idx != -1 {
		rangePart := part[:idx]
		stepStr := part[idx+1:]
		step, err := strconv.Atoi(stepStr)
		if err != nil || step <= 0 {
			return false
		}
		if rangePart == "*" || rangePart == "?" {
			return true
		}
		// Plain numeric a/n: equivalent to a-max/n.
		if !strings.Contains(rangePart, "-") {
			n, err := strconv.Atoi(rangePart)
			if err != nil || n < min || n > max {
				return false
			}
			return true
		}
		return validateRuleCronRange(rangePart, min, max)
	}
	if strings.Contains(part, "-") {
		return validateRuleCronRange(part, min, max)
	}
	n, err := strconv.Atoi(part)
	if err != nil {
		return false
	}
	return n >= min && n <= max
}

// validateRuleCronRange validates a-b segments within [min, max]. Both
// endpoints must be integers and start <= end.
func validateRuleCronRange(s string, min, max int) bool {
	parts := strings.SplitN(s, "-", 2)
	if len(parts) != 2 {
		return false
	}
	start, err1 := strconv.Atoi(parts[0])
	end, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return false
	}
	if start < min || start > max {
		return false
	}
	if end < min || end > max {
		return false
	}
	if start > end {
		return false
	}
	return true
}

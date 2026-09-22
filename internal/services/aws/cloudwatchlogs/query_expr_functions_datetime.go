package cloudwatchlogs

import (
	"strconv"
	"strings"
	"time"
)

// The datetime function family of the documented Logs Insights QL
// function library and the period/pattern helpers only its cases use.
// "Maps and lists are treated as null for string, number, and datetime
// functions" — the structure-null coercion runs in the library entry
// before these cases see their arguments.

// evalDatetimeFunctions holds the datetime family's dispatch: false
// means the name is not this family's.
func evalDatetimeFunctions(name string, f funcArgs) (interface{}, bool) {
	switch name {
	case "datefloor":
		return dateRound(f.str(0), f.str(1), false), true
	case "dateceil":
		return dateRound(f.str(0), f.str(1), true), true
	case "frommillis":
		if v, ok := f.num(0); ok {
			return timestampValue(int64(v)), true
		}
		return nil, true
	case "tomillis":
		if v, ok := f.num(0); ok {
			return v, true
		}
		// Timestamps stored as result-formatted or ISO strings are parsed
		// back to milliseconds.
		if ms, ok := parseResultTimestamp(f.str(0)); ok {
			return float64(ms), true
		}
		if t, err := time.Parse(time.RFC3339Nano, f.str(0)); err == nil {
			return float64(t.UnixMilli()), true
		}
		return nil, true
	case "now":
		// "Returns the time that the query processing was started, in
		// epoch seconds" — the start stamp when the context carries one,
		// so the value stays fixed across a long-running pipeline
		// instead of drifting with each evaluation.
		if f.ctx != nil {
			if f.ctx.startedMs > 0 {
				return float64(f.ctx.startedMs / 1000), true
			}
			return float64(f.ctx.now() / 1000), true
		}
		return float64(time.Now().UnixMilli() / 1000), true
	case "parsedate":
		layout := javaTimeLayout(f.str(1))
		tz := ""
		if len(f.vals) > 2 {
			tz = f.str(2)
		}
		loc := time.UTC
		if tz != "" {
			if l, err := time.LoadLocation(tz); err == nil {
				loc = l
			}
		}
		if t, err := time.ParseInLocation(layout, f.str(0), loc); err == nil {
			return float64(t.UnixMilli()), true
		}
		return nil, true
	case "formatdate", "strftime":
		layout := strftimeLayout(f.str(1))
		tz := ""
		if len(f.vals) > 2 {
			tz = f.str(2)
		}
		loc := time.UTC
		if tz != "" {
			if l, err := time.LoadLocation(tz); err == nil {
				loc = l
			}
		}
		ms, ok := f.num(0)
		if !ok {
			if parsed, ok2 := parseResultTimestamp(f.str(0)); ok2 {
				return time.UnixMilli(parsed).In(loc).Format(layout), true
			}
			if t, err := time.Parse(time.RFC3339Nano, f.str(0)); err == nil {
				return t.In(loc).Format(layout), true
			}
			return nil, true
		}
		return time.UnixMilli(int64(ms)).In(loc).Format(layout), true
	}
	return nil, false
}

// dateFloorCeil rounds a timestamp down (floor) or up (ceil) to the period
// boundary and truncates, yielding a timestamp-typed value.
func dateRound(ts, period string, ceil bool) interface{} {
	ms, ok := asNumber(ts)
	if !ok {
		if parsed, ok2 := parseResultTimestamp(ts); ok2 {
			ms = float64(parsed)
		} else if t, err := time.Parse(time.RFC3339Nano, ts); err == nil {
			ms = float64(t.UnixMilli())
		} else {
			return ""
		}
	}
	dur, ok := parsePeriodMillis(period)
	if !ok {
		return ""
	}
	v := int64(ms)
	if ceil {
		v = ((v + dur - 1) / dur) * dur
	} else {
		v = (v / dur) * dur
	}
	return timestampValue(v)
}

// parsePeriodMillis parses a period literal such as 5m, 10s, 1h, 2d, 1mo
// into milliseconds. Month, quarter, and year use calendar approximations.
// Full-word units and their documented abbreviations (with optional plural
// s) are accepted.
func parsePeriodMillis(p string) (int64, bool) {
	s := strings.TrimSpace(strings.ToLower(p))
	if s == "" {
		return 0, false
	}
	// Longest unit names first so that "milliseconds" is not read as "ms"
	// followed by stray characters.
	units := []struct {
		name string
		ms   float64
	}{
		{"milliseconds", 1}, {"millisecond", 1}, {"msec", 1}, {"msecs", 1}, {"ms", 1},
		{"seconds", 1000}, {"second", 1000}, {"secs", 1000}, {"sec", 1000}, {"s", 1000},
		{"minutes", 60 * 1000}, {"minute", 60 * 1000}, {"mins", 60 * 1000}, {"min", 60 * 1000}, {"m", 60 * 1000},
		{"hours", 3600 * 1000}, {"hrs", 3600 * 1000}, {"hr", 3600 * 1000}, {"hour", 3600 * 1000}, {"h", 3600 * 1000},
		{"days", 24 * 3600 * 1000}, {"day", 24 * 3600 * 1000}, {"d", 24 * 3600 * 1000},
		{"weeks", 7 * 24 * 3600 * 1000}, {"week", 7 * 24 * 3600 * 1000}, {"w", 7 * 24 * 3600 * 1000},
		{"months", 30 * 24 * 3600 * 1000}, {"month", 30 * 24 * 3600 * 1000}, {"mons", 30 * 24 * 3600 * 1000}, {"mon", 30 * 24 * 3600 * 1000}, {"mo", 30 * 24 * 3600 * 1000},
		{"quarters", 91 * 24 * 3600 * 1000}, {"quarter", 91 * 24 * 3600 * 1000}, {"qtrs", 91 * 24 * 3600 * 1000}, {"qtr", 91 * 24 * 3600 * 1000}, {"q", 91 * 24 * 3600 * 1000},
		{"years", 365 * 24 * 3600 * 1000}, {"year", 365 * 24 * 3600 * 1000}, {"yrs", 365 * 24 * 3600 * 1000}, {"yr", 365 * 24 * 3600 * 1000}, {"y", 365 * 24 * 3600 * 1000},
	}
	for _, u := range units {
		if strings.HasSuffix(s, u.name) {
			numStr := strings.TrimSpace(strings.TrimSuffix(s, u.name))
			if numStr == "" {
				numStr = "1"
			}
			n, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return 0, false
			}
			return int64(n * u.ms), true
		}
	}
	// Bare numbers are treated as milliseconds.
	if n, err := strconv.ParseFloat(s, 64); err == nil {
		return int64(n), true
	}
	return 0, false
}

// capPeriodMillis applies the documented caps for time units: 1000 for
// milliseconds, 60 for seconds and minutes, 24 for hours. The cap depends on
// the unit the period was written in, not on the total duration.
func capPeriodMillis(ms, unitMs int64) int64 {
	var cap int64
	switch unitMs {
	case 1:
		cap = 1000
	case 1000:
		cap = 60 * 1000
	case 60 * 1000:
		cap = 60 * 60 * 1000
	case 3600 * 1000:
		cap = 24 * 3600 * 1000
	default:
		return ms
	}
	if ms > cap {
		return cap
	}
	return ms
}

// javaTimeLayout converts the documented subset of Java DateTimeFormatter
// patterns to Go layouts for parseDate.
func javaTimeLayout(pattern string) string {
	repl := []struct{ java, goLayout string }{
		{"yyyy", "2006"}, {"yy", "06"}, {"MM", "01"}, {"dd", "02"},
		{"HH", "15"}, {"mm", "04"}, {"ss", "05"}, {"SSS", "000"},
		{"a", "PM"}, {"EEEE", "Monday"}, {"EEE", "Mon"}, {"Z", "Z07:00"}, {"X", "Z07:00"},
	}
	out := pattern
	for _, r := range repl {
		out = strings.ReplaceAll(out, r.java, r.goLayout)
	}
	return out
}

// strftimeLayout converts the documented strftime-style specifiers to Go
// layouts for formatDate.
func strftimeLayout(pattern string) string {
	repl := []struct{ spec, goLayout string }{
		{"%Y", "2006"}, {"%y", "06"}, {"%m", "01"}, {"%d", "02"},
		{"%H", "15"}, {"%M", "04"}, {"%S", "05"}, {"%j", "002"},
		{"%p", "PM"}, {"%B", "January"}, {"%b", "Jan"}, {"%Z", "MST"},
	}
	out := pattern
	for _, r := range repl {
		out = strings.ReplaceAll(out, r.spec, r.goLayout)
	}
	return out
}

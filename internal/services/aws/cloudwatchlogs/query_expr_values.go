package cloudwatchlogs

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Value semantics for query expressions: coercion, comparison, and rendering.

// --- value helpers ---

// timestampValue is a timestamp-typed result value in epoch milliseconds.
// The documented result type of fromMillis, datefloor, dateceil and bin is
// Timestamp: asNumber keeps the numeric interpretation working inside
// expressions while asString renders the timestamp form results present.
type timestampValue int64

func truthy(v interface{}) bool {
	b, ok := v.(bool)
	return ok && b
}

// boolToNum renders a boolean as the documented 1/0 number result.
func boolToNum(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

func asNumber(v interface{}) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case int:
		return float64(x), true
	case int64:
		return float64(x), true
	case timestampValue:
		return float64(x), true
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(x), 64)
		if err != nil {
			return 0, false
		}
		return f, true
	case bool:
		return 0, false
	}
	return 0, false
}

// formatNumber renders a numeric value the way query results present it.
func formatNumber(f float64) string {
	if f == math.Trunc(f) && math.Abs(f) < 1e15 {
		return strconv.FormatInt(int64(f), 10)
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// asString converts any value to its result-string representation. Maps and
// lists render as canonical JSON, matching how structure values round-trip
// through the string-typed result rows.
func asString(v interface{}) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case bool:
		if x {
			return "true"
		}
		return "false"
	case float64:
		return formatNumber(x)
	case int:
		return strconv.Itoa(x)
	case int64:
		return strconv.FormatInt(x, 10)
	case timestampValue:
		return time.UnixMilli(int64(x)).UTC().Format(resultTimestampLayout)
	case []interface{}, map[string]interface{}:
		b, err := json.Marshal(x)
		if err != nil {
			return ""
		}
		return string(b)
	}
	return fmt.Sprintf("%v", v)
}

// storeValue converts an evaluated value to the string stored in a row,
// keeping numbers numeric-looking so later numeric comparisons still work.
func storeValue(v interface{}) string {
	return asString(v)
}

func valuesEqual(l, r interface{}) bool {
	if l == nil || r == nil {
		return false
	}
	if ln, lok := asNumber(l); lok {
		if rn, rok := asNumber(r); rok {
			return ln == rn
		}
	}
	if lb, lok := l.(bool); lok {
		if rb, rok := r.(bool); rok {
			return lb == rb
		}
		return false
	}
	if _, rok := r.(bool); rok {
		return false
	}
	ls, lok := l.(string)
	rs, rok2 := r.(string)
	if lok && rok2 {
		return ls == rs
	}
	return asString(l) == asString(r)
}

func compareValues(l, r interface{}, op string) interface{} {
	if l == nil || r == nil {
		return false
	}
	ln, lok := asNumber(l)
	rn, rok := asNumber(r)
	var cmp int
	if lok && rok {
		switch {
		case ln < rn:
			cmp = -1
		case ln > rn:
			cmp = 1
		}
	} else {
		ls, rs := asString(l), asString(r)
		cmp = strings.Compare(ls, rs)
	}
	switch op {
	case "<":
		return cmp < 0
	case "<=":
		return cmp <= 0
	case ">":
		return cmp > 0
	case ">=":
		return cmp >= 0
	}
	return false
}

// globMatch implements the wildcard matching used by like with quoted
// patterns: * matches any run of characters.
func globMatch(pattern, text string) bool {
	if !strings.Contains(pattern, "*") {
		return pattern == text
	}
	re := regexp.QuoteMeta(pattern)
	re = strings.ReplaceAll(re, `\*`, ".*")
	matched, _ := regexp.MatchString("^"+re+"$", text)
	return matched
}

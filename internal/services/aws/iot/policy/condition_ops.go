package policy

import (
	"net"
	"strconv"
	"strings"
	"time"
)

// ConditionOperator evaluates a condition: actual is the runtime value
// (e.g. the connecting client's source IP) and expected is the value from
// the policy document.
type ConditionOperator func(actual, expected string) bool

var conditionOperators = map[string]ConditionOperator{
	// --- String operators ---
	"StringEquals":              func(a, b string) bool { return a == b },
	"StringNotEquals":           func(a, b string) bool { return a != b },
	"StringEqualsIgnoreCase":    func(a, b string) bool { return strings.EqualFold(a, b) },
	"StringNotEqualsIgnoreCase": func(a, b string) bool { return !strings.EqualFold(a, b) },
	"StringLike":                func(a, b string) bool { return wildcardMatch(a, b) },
	"StringNotLike":             func(a, b string) bool { return !wildcardMatch(a, b) },

	// --- Bool operator ---
	// AWS treats Bool values as case-insensitive "true"/"false".
	"Bool": func(a, b string) bool { return strings.EqualFold(a, b) },

	// --- ARN operators ---
	// ArnEquals is semantically identical to StringEquals for ARNs.
	"ArnEquals":    func(a, b string) bool { return a == b },
	"ArnNotEquals": func(a, b string) bool { return a != b },
	// ArnLike uses the same wildcard matching as StringLike.
	"ArnLike":    func(a, b string) bool { return wildcardMatch(a, b) },
	"ArnNotLike": func(a, b string) bool { return !wildcardMatch(a, b) },

	// --- IP address operators ---
	"IpAddress":    ipMatch,
	"NotIpAddress": func(a, b string) bool { return !ipMatch(a, b) },

	// --- Numeric operators ---
	// All numeric operators return false when either value cannot be
	// parsed as a number (fail-closed).
	"NumericEquals":            func(a, b string) bool { return numericOp(a, b, func(x, y float64) bool { return x == y }) },
	"NumericNotEquals":         func(a, b string) bool { return numericOp(a, b, func(x, y float64) bool { return x != y }) },
	"NumericLessThan":          func(a, b string) bool { return numericOp(a, b, func(x, y float64) bool { return x < y }) },
	"NumericLessThanEquals":    func(a, b string) bool { return numericOp(a, b, func(x, y float64) bool { return x <= y }) },
	"NumericGreaterThan":       func(a, b string) bool { return numericOp(a, b, func(x, y float64) bool { return x > y }) },
	"NumericGreaterThanEquals": func(a, b string) bool { return numericOp(a, b, func(x, y float64) bool { return x >= y }) },

	// --- Date operators ---
	// Expects RFC 3339 or epoch seconds. All return false on parse
	// failure (fail-closed).
	"DateLessThan": func(a, b string) bool { return dateOp(a, b, func(x, y time.Time) bool { return x.Before(y) }) },
	"DateLessThanEquals": func(a, b string) bool {
		return dateOp(a, b, func(x, y time.Time) bool { return x.Before(y) || x.Equal(y) })
	},
	"DateGreaterThan": func(a, b string) bool { return dateOp(a, b, func(x, y time.Time) bool { return x.After(y) }) },
	"DateGreaterThanEquals": func(a, b string) bool {
		return dateOp(a, b, func(x, y time.Time) bool { return x.After(y) || x.Equal(y) })
	},
}

// ipMatch checks whether actual IP falls within the expected CIDR or
// exact IP address. Both IPv4 and IPv6 are supported.
func ipMatch(actual, expected string) bool {
	expected = strings.TrimSpace(expected)
	// If expected has no /, treat as /32 (IPv4) or /128 (IPv6).
	if !strings.Contains(expected, "/") {
		ip := net.ParseIP(expected)
		if ip == nil {
			return false
		}
		if ip.To4() != nil {
			expected += "/32"
		} else {
			expected += "/128"
		}
	}
	_, cidr, err := net.ParseCIDR(expected)
	if err != nil {
		return false
	}
	ip := net.ParseIP(actual)
	if ip == nil {
		return false
	}
	return cidr.Contains(ip)
}

// numericOp parses two numeric strings and applies the comparison
// function. Returns false when either value fails to parse
// (fail-closed: the condition simply does not match).
func numericOp(a, b string, cmp func(af, bf float64) bool) bool {
	af, errA := strconv.ParseFloat(a, 64)
	bf, errB := strconv.ParseFloat(b, 64)
	if errA != nil || errB != nil {
		return false
	}
	return cmp(af, bf)
}

// dateOp parses two date strings and applies the comparison function.
// Returns false when either value fails to parse (fail-closed).
func dateOp(a, b string, cmp func(at, bt time.Time) bool) bool {
	at := parseDate(a)
	bt := parseDate(b)
	if at.IsZero() || bt.IsZero() {
		return false
	}
	return cmp(at, bt)
}

func parseDate(s string) time.Time {
	s = strings.TrimSpace(s)
	// Try epoch seconds first. Require at least 10 digits to avoid
	// misinterpreting short numbers (e.g. "2024" as epoch 1970+2024s).
	if epoch, err := strconv.ParseInt(s, 10, 64); err == nil && len(s) >= 10 {
		return time.Unix(epoch, 0).UTC()
	}
	// Try common RFC 3339 / ISO 8601 layouts.
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC()
		}
	}
	return time.Time{}
}

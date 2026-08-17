package cloudwatchlogs

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"math"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// The documented Logs Insights QL function library and its helpers.

// --- function library ---
// The functions below implement the documented Logs Insights QL function
// set: general, numeric, string, IP address, JSON, and datetime functions.

func callQueryFunction(name string, args []exprNode, row *queryResultRow, ctx *execContext) interface{} {
	vals := make([]interface{}, len(args))
	for i, a := range args {
		vals[i] = a.eval(row, ctx)
	}
	arg := func(i int) interface{} {
		if i < len(vals) {
			return vals[i]
		}
		return nil
	}
	num := func(i int) (float64, bool) { return asNumber(arg(i)) }
	str := func(i int) string { return asString(arg(i)) }

	switch name {
	// General functions.
	case "ispresent":
		return arg(0) != nil
	case "ispresentornull":
		return arg(0) != nil
	case "coalesce":
		for _, v := range vals {
			if v != nil && asString(v) != "" {
				return v
			}
		}
		return nil
	case "case":
		// case(cond1, val1, ..., [default]): pairs then optional default.
		i := 0
		for i+1 < len(vals) {
			if truthy(vals[i]) {
				return vals[i+1]
			}
			i += 2
		}
		if len(vals)%2 == 1 {
			return vals[len(vals)-1]
		}
		return nil
	case "if":
		// The condition is validated as a three-argument call at parse time;
		// arg bounds defensively so a nil context cannot panic here.
		if truthy(arg(0)) {
			return arg(1)
		}
		return arg(2)
	case "isnumeric":
		_, ok := num(0)
		return ok
	case "messagesize":
		return float64(len(str(0)))
	case "querystarttime":
		if ctx != nil {
			return float64(ctx.startTime)
		}
		return nil
	case "queryendtime":
		if ctx != nil {
			return float64(ctx.endTime)
		}
		return nil
	case "querytimerange":
		if ctx != nil {
			return float64(ctx.endTime - ctx.startTime)
		}
		return nil

	// Numeric functions.
	case "abs":
		if f, ok := num(0); ok {
			return math.Abs(f)
		}
		return nil
	case "ceil":
		if f, ok := num(0); ok {
			return math.Ceil(f)
		}
		return nil
	case "floor":
		if f, ok := num(0); ok {
			return math.Floor(f)
		}
		return nil
	case "greatest":
		best, ok := num(0)
		if !ok {
			return nil
		}
		for i := 1; i < len(vals); i++ {
			if f, ok := num(i); ok && f > best {
				best = f
			}
		}
		return best
	case "least":
		best, ok := num(0)
		if !ok {
			return nil
		}
		for i := 1; i < len(vals); i++ {
			if f, ok := num(i); ok && f < best {
				best = f
			}
		}
		return best
	case "log":
		if f, ok := num(0); ok && f > 0 {
			return math.Log(f)
		}
		return nil
	case "round":
		if f, ok := num(0); ok {
			d := 0.0
			if len(vals) > 1 {
				if d2, ok := num(1); ok {
					d = d2
				}
			}
			pow := math.Pow(10, d)
			return math.Round(f*pow) / pow
		}
		return nil
	case "sqrt":
		if f, ok := num(0); ok && f >= 0 {
			return math.Sqrt(f)
		}
		return nil
	case "haversine":
		lat1, ok1 := num(0)
		lon1, ok2 := num(1)
		lat2, ok3 := num(2)
		lon2, ok4 := num(3)
		if !ok1 || !ok2 || !ok3 || !ok4 {
			return nil
		}
		const toRad = math.Pi / 180
		const earthKm = 6371.0
		dLat := (lat2 - lat1) * toRad
		dLon := (lon2 - lon1) * toRad
		a := math.Sin(dLat/2)*math.Sin(dLat/2) +
			math.Cos(lat1*toRad)*math.Cos(lat2*toRad)*math.Sin(dLon/2)*math.Sin(dLon/2)
		c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
		return earthKm * c
	case "tonumber":
		if f, ok := num(0); ok {
			return f
		}
		return nil
	case "toint":
		if f, ok := num(0); ok {
			return float64(int32(f))
		}
		return nil
	case "tolong":
		if f, ok := num(0); ok {
			return float64(int64(f))
		}
		return nil
	case "todouble":
		if f, ok := num(0); ok {
			return f
		}
		return nil

	// String functions documented to return Number yield 1/0.
	case "isempty":
		return boolToNum(arg(0) == nil || str(0) == "")
	case "isblank":
		return boolToNum(arg(0) == nil || strings.TrimSpace(str(0)) == "")
	case "concat":
		var b strings.Builder
		for i := range vals {
			b.WriteString(str(i))
		}
		return b.String()
	case "ltrim":
		if len(vals) > 1 {
			return strings.TrimLeft(str(0), str(1))
		}
		return strings.TrimLeft(str(0), " \t\r\n")
	case "rtrim":
		if len(vals) > 1 {
			return strings.TrimRight(str(0), str(1))
		}
		return strings.TrimRight(str(0), " \t\r\n")
	case "trim":
		if len(vals) > 1 {
			return strings.Trim(str(0), str(1))
		}
		return strings.Trim(str(0), " \t\r\n")
	case "strlen":
		return float64(len([]rune(str(0))))
	case "toupper":
		return strings.ToUpper(str(0))
	case "tolower":
		return strings.ToLower(str(0))
	case "substr":
		s := []rune(str(0))
		start, ok := num(1)
		if !ok {
			return nil
		}
		si := int(start)
		if si < 0 {
			si = len(s) + si
		}
		if si < 0 || si > len(s) {
			return ""
		}
		if len(vals) > 2 {
			l, ok := num(2)
			if !ok {
				return nil
			}
			ei := si + int(l)
			if ei > len(s) {
				ei = len(s)
			}
			if ei < si {
				return ""
			}
			return string(s[si:ei])
		}
		return string(s[si:])
	case "replace":
		return strings.ReplaceAll(str(0), str(1), str(2))
	case "regexreplace":
		re, err := regexp.Compile(str(1))
		if err != nil {
			return nil
		}
		return re.ReplaceAllString(str(0), str(2))
	case "strcontains":
		if len(vals) > 2 && truthy(vals[2]) {
			return boolToNum(strings.Contains(strings.ToLower(str(0)), strings.ToLower(str(1))))
		}
		return boolToNum(strings.Contains(str(0), str(1)))
	case "startswith":
		return boolToNum(strings.HasPrefix(str(0), str(1)))
	case "endswith":
		return boolToNum(strings.HasSuffix(str(0), str(1)))
	case "urlencode":
		return url.QueryEscape(str(0))
	case "urldecode":
		if dec, err := url.QueryUnescape(str(0)); err == nil {
			return dec
		}
		return nil
	case "base64encode":
		return base64.StdEncoding.EncodeToString([]byte(str(0)))
	case "base64decode":
		if dec, err := base64.StdEncoding.DecodeString(str(0)); err == nil {
			return string(dec)
		}
		return nil
	case "split":
		var parts []interface{}
		for _, p := range strings.Split(str(0), str(1)) {
			parts = append(parts, p)
		}
		return parts
	case "hextoascii":
		if b, err := hexDecodeString(str(0)); err == nil {
			return string(b)
		}
		return nil
	case "hextodec":
		ui, err := strconv.ParseUint(strings.TrimPrefix(strings.ToLower(str(0)), "0x"), 16, 64)
		if err != nil {
			return nil
		}
		return float64(ui)
	case "dectohex":
		f, ok := num(0)
		if !ok {
			return nil
		}
		n := int64(f)
		prefix := "0x"
		if n < 0 {
			prefix = "-0x"
			n = -n
		}
		return prefix + strconv.FormatUint(uint64(n), 16)

	// IP address functions.
	case "isvalidip":
		return net.ParseIP(str(0)) != nil
	case "isvalidipv4":
		ip := net.ParseIP(str(0))
		return ip != nil && ip.To4() != nil && strings.Count(str(0), ".") == 3
	case "isvalidipv6":
		ip := net.ParseIP(str(0))
		return ip != nil && strings.Contains(str(0), ":")
	case "isipinsubnet":
		return ipInSubnet(str(0), str(1), false)
	case "isipv4insubnet":
		return ipInSubnet(str(0), str(1), true)
	case "isipv6insubnet":
		ip := net.ParseIP(str(0))
		_, cidr, err := net.ParseCIDR(str(1))
		return err == nil && ip != nil && strings.Contains(str(1), ":") && cidr.Contains(ip)
	case "ipv4tonumber":
		ip := net.ParseIP(str(0))
		if ip == nil || ip.To4() == nil {
			return nil
		}
		v := uint32(0)
		for _, b := range ip.To4() {
			v = v<<8 | uint32(b)
		}
		return float64(v)
	case "isprivateip":
		ip := net.ParseIP(str(0))
		if ip == nil || ip.To4() == nil {
			return false
		}
		for _, cidr := range []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"} {
			_, nw, _ := net.ParseCIDR(cidr)
			if nw.Contains(ip) {
				return true
			}
		}
		return false
	case "ispublicip":
		ip := net.ParseIP(str(0))
		if ip == nil {
			return false
		}
		private := []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "127.0.0.0/8", "169.254.0.0/16"}
		reserved := []string{"0.0.0.0/8", "192.0.2.0/24", "198.51.100.0/24", "203.0.113.0/24", "224.0.0.0/4", "240.0.0.0/4"}
		for _, cidr := range append(private, reserved...) {
			_, nw, _ := net.ParseCIDR(cidr)
			if nw.Contains(ip) {
				return false
			}
		}
		return true
	case "isreservedip":
		ip := net.ParseIP(str(0))
		if ip == nil {
			return false
		}
		reserved := []string{"0.0.0.0/8", "10.0.0.0/8", "127.0.0.0/8", "169.254.0.0/16", "172.16.0.0/12",
			"192.0.2.0/24", "192.168.0.0/16", "198.51.100.0/24", "203.0.113.0/24", "224.0.0.0/4", "240.0.0.0/4"}
		for _, cidr := range reserved {
			_, nw, _ := net.ParseCIDR(cidr)
			if nw.Contains(ip) {
				return true
			}
		}
		return false

	// JSON functions.
	case "jsonparse":
		s := strings.TrimSpace(str(0))
		if strings.HasPrefix(s, "{") || strings.HasPrefix(s, "[") {
			var decoded interface{}
			if err := json.Unmarshal([]byte(s), &decoded); err == nil {
				return decoded
			}
		}
		return nil
	case "jsonstringify":
		return asString(arg(0))
	case "jsonarraysize":
		var decoded interface{}
		if err := json.Unmarshal([]byte(strings.TrimSpace(str(0))), &decoded); err != nil {
			return float64(0)
		}
		if list, ok := decoded.([]interface{}); ok {
			return float64(len(list))
		}
		return float64(0)
	case "jsonarraycontains":
		var decoded interface{}
		if err := json.Unmarshal([]byte(strings.TrimSpace(str(0))), &decoded); err != nil {
			return false
		}
		list, ok := decoded.([]interface{})
		if !ok {
			return false
		}
		for _, item := range list {
			if valuesEqual(item, arg(1)) {
				return true
			}
		}
		return false

	// Datetime functions.
	case "datefloor":
		return dateRound(str(0), str(1), false)
	case "dateceil":
		return dateRound(str(0), str(1), true)
	case "frommillis":
		if f, ok := num(0); ok {
			return timestampValue(int64(f))
		}
		return nil
	case "tomillis":
		if f, ok := num(0); ok {
			return f
		}
		// Timestamps stored as result-formatted or ISO strings are parsed
		// back to milliseconds.
		if ms, ok := parseResultTimestamp(str(0)); ok {
			return float64(ms)
		}
		if t, err := time.Parse(time.RFC3339Nano, str(0)); err == nil {
			return float64(t.UnixMilli())
		}
		return nil
	case "now":
		if ctx != nil {
			return float64(ctx.now() / 1000)
		}
		return float64(time.Now().UnixMilli() / 1000)
	case "parsedate":
		layout := javaTimeLayout(str(1))
		tz := ""
		if len(vals) > 2 {
			tz = str(2)
		}
		loc := time.UTC
		if tz != "" {
			if l, err := time.LoadLocation(tz); err == nil {
				loc = l
			}
		}
		if t, err := time.ParseInLocation(layout, str(0), loc); err == nil {
			return float64(t.UnixMilli())
		}
		return nil
	case "formatdate", "strftime":
		layout := strftimeLayout(str(1))
		tz := ""
		if len(vals) > 2 {
			tz = str(2)
		}
		loc := time.UTC
		if tz != "" {
			if l, err := time.LoadLocation(tz); err == nil {
				loc = l
			}
		}
		ms, ok := num(0)
		if !ok {
			if parsed, ok2 := parseResultTimestamp(str(0)); ok2 {
				return time.UnixMilli(parsed).In(loc).Format(layout)
			}
			if t, err := time.Parse(time.RFC3339Nano, str(0)); err == nil {
				return t.In(loc).Format(layout)
			}
			return nil
		}
		return time.UnixMilli(int64(ms)).In(loc).Format(layout)

	// Hashing functions.
	case "md5":
		sum := md5.Sum([]byte(str(0)))
		return hex.EncodeToString(sum[:])
	case "sha256":
		sum := sha256.Sum256([]byte(str(0)))
		return hex.EncodeToString(sum[:])
	}
	return nil
}

func hexDecodeString(s string) ([]byte, error) {
	s = strings.TrimPrefix(strings.ToLower(s), "0x")
	return hex.DecodeString(s)
}

// ipInSubnet reports whether the address is within the subnet; when v4Only
// is set the address must be IPv4 and the subnet IPv4 as well.
func ipInSubnet(addr, subnet string, v4Only bool) bool {
	ip := net.ParseIP(addr)
	_, cidr, err := net.ParseCIDR(subnet)
	if err != nil || ip == nil {
		return false
	}
	if v4Only && (ip.To4() == nil || strings.Contains(subnet, ":")) {
		return false
	}
	return cidr.Contains(ip)
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
		{"a", "PM"}, {"EEE", "Mon"}, {"EEEE", "Monday"}, {"Z", "Z07:00"}, {"X", "Z07:00"},
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

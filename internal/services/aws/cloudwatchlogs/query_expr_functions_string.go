package cloudwatchlogs

import (
	"encoding/base64"
	"encoding/hex"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// The string function family of the documented Logs Insights QL function
// library, with the encoding helpers only its cases use. "Maps and lists
// are treated as null for string, number, and datetime functions" — the
// structure-null coercion runs in the library entry before these cases
// see their arguments.

// evalStringFunctions holds the string family's dispatch: false means
// the name is not this family's. The functions documented to return
// Number yield 1/0.
func evalStringFunctions(name string, f funcArgs) (interface{}, bool) {
	switch name {
	case "isempty":
		return boolToNum(f.arg(0) == nil || f.str(0) == ""), true
	case "isblank":
		return boolToNum(f.arg(0) == nil || strings.TrimSpace(f.str(0)) == ""), true
	case "concat":
		var b strings.Builder
		for i := range f.vals {
			b.WriteString(f.str(i))
		}
		return b.String(), true
	case "ltrim":
		if len(f.vals) > 1 {
			return strings.TrimLeft(f.str(0), f.str(1)), true
		}
		return strings.TrimLeft(f.str(0), " \t\r\n"), true
	case "rtrim":
		if len(f.vals) > 1 {
			return strings.TrimRight(f.str(0), f.str(1)), true
		}
		return strings.TrimRight(f.str(0), " \t\r\n"), true
	case "trim":
		if len(f.vals) > 1 {
			return strings.Trim(f.str(0), f.str(1)), true
		}
		return strings.Trim(f.str(0), " \t\r\n"), true
	case "strlen":
		return float64(len([]rune(f.str(0)))), true
	case "toupper":
		return strings.ToUpper(f.str(0)), true
	case "tolower":
		return strings.ToLower(f.str(0)), true
	case "substr":
		s := []rune(f.str(0))
		start, ok := f.num(1)
		if !ok {
			return nil, true
		}
		si := int(start)
		if si < 0 {
			si = len(s) + si
		}
		if si < 0 || si > len(s) {
			return "", true
		}
		if len(f.vals) > 2 {
			l, ok := f.num(2)
			if !ok {
				return nil, true
			}
			ei := si + int(l)
			if ei > len(s) {
				ei = len(s)
			}
			if ei < si {
				return "", true
			}
			return string(s[si:ei]), true
		}
		return string(s[si:]), true
	case "replace":
		return strings.ReplaceAll(f.str(0), f.str(1), f.str(2)), true
	case "regexreplace":
		re, err := regexp.Compile(f.str(1))
		if err != nil {
			return nil, true
		}
		return re.ReplaceAllString(f.str(0), f.str(2)), true
	case "strcontains":
		if len(f.vals) > 2 && truthy(f.vals[2]) {
			return boolToNum(strings.Contains(strings.ToLower(f.str(0)), strings.ToLower(f.str(1)))), true
		}
		return boolToNum(strings.Contains(f.str(0), f.str(1))), true
	case "startswith":
		return boolToNum(strings.HasPrefix(f.str(0), f.str(1))), true
	case "endswith":
		return boolToNum(strings.HasSuffix(f.str(0), f.str(1))), true
	case "urlencode":
		return url.QueryEscape(f.str(0)), true
	case "urldecode":
		if dec, err := url.QueryUnescape(f.str(0)); err == nil {
			return dec, true
		}
		return nil, true
	case "base64encode":
		return base64.StdEncoding.EncodeToString([]byte(f.str(0))), true
	case "base64decode":
		if dec, err := base64.StdEncoding.DecodeString(f.str(0)); err == nil {
			return string(dec), true
		}
		return nil, true
	case "split":
		var parts []interface{}
		for _, p := range strings.Split(f.str(0), f.str(1)) {
			parts = append(parts, p)
		}
		return parts, true
	case "hextoascii":
		if b, err := hexDecodeString(f.str(0)); err == nil {
			return string(b), true
		}
		return nil, true
	case "hextodec":
		ui, err := strconv.ParseUint(strings.TrimPrefix(strings.ToLower(f.str(0)), "0x"), 16, 64)
		if err != nil {
			return nil, true
		}
		return float64(ui), true
	case "dectohex":
		v, ok := f.num(0)
		if !ok {
			return nil, true
		}
		n := int64(v)
		prefix := "0x"
		if n < 0 {
			prefix = "-0x"
			n = -n
		}
		return prefix + strconv.FormatUint(uint64(n), 16), true
	}
	return nil, false
}

func hexDecodeString(s string) ([]byte, error) {
	s = strings.TrimPrefix(strings.ToLower(s), "0x")
	return hex.DecodeString(s)
}

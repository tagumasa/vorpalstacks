package cloudwatchlogs

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// The grok processor's pattern library and compiler. Every pattern the
// documentation's "Supported grok patterns" table names resolves to the
// regular expression its description and examples define; the three
// common-log composites (APACHE_ACCESS_LOG, NGINX_ACCESS_LOG,
// SYSLOG5424) expand to their documented typed field sets.

// grokDotPlaceholder stands in for a field name's dots (Go's named
// groups admit word characters only); dotted field names are JSON paths
// in the output.
const grokDotPlaceholder = "__GROKDOT__"

// grokPatterns maps every documented pattern name to its regular
// expression. Composite patterns whose output carries extra fields name
// those groups inline (HOSTPORT's PORT, URIHOST's port, SYSLOGPROG's
// program and pid, SYSLOGFACILITY's facility and priority).
var grokPatterns = map[string]string{
	"USERNAME":             `[A-Za-z0-9._-]+`,
	"USER":                 `[A-Za-z0-9._-]+`,
	"INT":                  `[-+]?[0-9]+`,
	"NUMBER":               `[-+]?[0-9]+(?:\.[0-9]+)?`,
	"BASE10NUM":            `[-+]?(?:[0-9]+(?:\.[0-9]+)?|\.[0-9]+)`,
	"BASE16NUM":            `[-+]?(?:0[xX])?[0-9A-Fa-f]+`,
	"POSINT":               `[1-9][0-9]*`,
	"NONNEGINT":            `[0-9]+`,
	"WORD":                 `\w+`,
	"NOTSPACE":             `\S+`,
	"SPACE":                `\s*`,
	"DATA":                 `.*?`,
	"GREEDYDATA":           `.*`,
	"GREEDYDATA_MULTILINE": `(?s).*`,
	"QUOTEDSTRING":         `"(?:[^"\\]|\\.)*"|'(?:[^'\\]|\\.)*'`,
	"UUID":                 `[0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{4}-[0-9A-Fa-f]{12}`,
	"URN":                  `urn:[A-Za-z0-9][A-Za-z0-9-]{0,31}:[^\s]+`,
	// "It will not match ARNs that are missing information between
	// colons" — the region and account segments are never empty.
	"ARN": `arn:(?:aws|aws-cn|aws-us-gov):[A-Za-z0-9_-]+:[A-Za-z0-9-]+:[0-9]+:[^\s]+`,

	"CISCOMAC":   `[0-9A-Fa-f]{4}\.[0-9A-Fa-f]{4}\.[0-9A-Fa-f]{4}`,
	"WINDOWSMAC": `[0-9A-Fa-f]{2}-[0-9A-Fa-f]{2}-[0-9A-Fa-f]{2}-[0-9A-Fa-f]{2}-[0-9A-Fa-f]{2}-[0-9A-Fa-f]{2}`,
	"COMMONMAC":  `[0-9A-Fa-f]{2}:[0-9A-Fa-f]{2}:[0-9A-Fa-f]{2}:[0-9A-Fa-f]{2}:[0-9A-Fa-f]{2}:[0-9A-Fa-f]{2}`,
	"IPV4":       `(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])(?:\.(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])){3}`,
	"IPV6":       `(?:[0-9A-Fa-f]{1,4}:){7}[0-9A-Fa-f]{1,4}|(?:[0-9A-Fa-f]{1,4}:)*:(?:[0-9A-Fa-f]{1,4}:)*[0-9A-Fa-f]{1,4}|::(?:[Ff]{4}(?:0:0){0,2}:)?(?:[0-9]{1,3}\.){3}[0-9]{1,3}`,
	"HOSTNAME":   `[A-Za-z0-9](?:[A-Za-z0-9._-]*)`,
	"HOST":       `[A-Za-z0-9](?:[A-Za-z0-9._-]*)`,

	// UNIXPATH "Matches URL paths, potentially including query
	// parameters" (documented example: /search?q=regex).
	"UNIXPATH": `/[^\s]*`,
	"WINPATH":  `[A-Za-z]:\\[^\s]+`,
	"TTY":      `/dev/(?:pts/)?tty[0-9]+`,
	"URIPROTO": `[A-Za-z]+(?:\+[A-Za-z]+)*`,
	"URIPATH":  `/[^\s?]+`,
	"URIPARAM": `\?[^\s]+`,

	"MONTH":     `(?:Jan(?:uary)?|Feb(?:ruary)?|Mar(?:ch)?|Apr(?:il)?|May|Jun(?:e)?|Jul(?:y)?|Aug(?:ust)?|Sep(?:tember)?|Sept|Oct(?:ober)?|Nov(?:ember)?|Dec(?:ember)?)`,
	"MONTHNUM":  `(?:0?[1-9]|1[0-2])`,
	"MONTHNUM2": `(?:0[1-9]|1[0-2])`,
	"MONTHDAY":  `(?:0?[1-9]|[12][0-9]|3[01])`,
	"YEAR":      `(?:[0-9]{2})?[0-9]{2}`,
	"DAY":       `(?:Mon(?:day)?|Tue(?:sday)?|Wed(?:nesday)?|Thu(?:rsday)?|Fri(?:day)?|Sat(?:urday)?|Sun(?:day)?)`,
	"HOUR":      `(?:0?[0-9]|1[0-9]|2[0-3])`,
	"MINUTE":    `[0-5][0-9]`,
	"SECOND":    `(?:0?[0-9]|[0-5][0-9]|60)(?:[.:][0-9]+)?`,
	"TIME":      `(?:0?[0-9]|1[0-9]|2[0-3]):[0-5][0-9](?::(?:0?[0-9]|[0-5][0-9]|60)(?:[.:][0-9]+)?)?`,
	"DATE_US":   `(?:0?[1-9]|1[0-2])[/-](?:0?[1-9]|[12][0-9]|3[01])[/-](?:[0-9]{2})?[0-9]{2}`,
	"DATE_EU":   `(?:0?[1-9]|[12][0-9]|3[01])[/.-](?:0?[1-9]|1[0-2])[/.-](?:[0-9]{2})?[0-9]{2}`,
	// The offset hours are "(H)H" — one or two digits (documented
	// examples: -530 and +5:30 alongside +05:30).
	"ISO8601_TIMEZONE":  `(?:Z|[+-][0-9]{1,2}(?::?[0-9]{2})?)`,
	"ISO8601_SECOND":    `(?:0?[0-9]|[0-5][0-9]|60)(?:[.:][0-9]+)?`,
	"TIMESTAMP_ISO8601": `(?:[0-9]{2})?[0-9]{2}-(?:0?[1-9]|1[0-2])-(?:0?[1-9]|[12][0-9]|3[01])T(?:0?[0-9]|1[0-9]|2[0-3]):[0-5][0-9](?::(?:0?[0-9]|[0-5][0-9])(?:[.:][0-9]+)?)?(?:%{ISO8601_TIMEZONE})?`,
	// DATE "Matches either a date in the US format ... or in the EU
	// format"; DATESTAMP "Matches %{DATE} followed by %{TIME} pattern,
	// separated by space or hyphen" — both date orders qualify.
	"DATE":               `(?:0?[1-9]|1[0-2])[/-](?:0?[1-9]|[12][0-9]|3[01])[/-](?:[0-9]{2})?[0-9]{2}|(?:0?[1-9]|[12][0-9]|3[01])[/.-](?:0?[1-9]|1[0-2])[/.-](?:[0-9]{2})?[0-9]{2}`,
	"DATESTAMP":          `(?:%{DATE})[ -](?:%{TIME})`,
	"TZ":                 `(?:PST|PDT|MST|MDT|CST|CDT|EST|EDT|UTC)`,
	"DATESTAMP_RFC822":   `(?:Mon(?:day)?|Tue(?:sday)?|Wed(?:nesday)?|Thu(?:rsday)?|Fri(?:day)?|Sat(?:urday)?|Sun(?:day)?) (?:Jan(?:uary)?|Feb(?:ruary)?|Mar(?:ch)?|Apr(?:il)?|May|Jun(?:e)?|Jul(?:y)?|Aug(?:ust)?|Sep(?:tember)?|Sept|Oct(?:ober)?|Nov(?:ember)?|Dec(?:ember)?) (?:0?[1-9]|[12][0-9]|3[01]) (?:[0-9]{2})?[0-9]{2} (?:0?[0-9]|1[0-9]|2[0-3]):[0-5][0-9](?::(?:0?[0-9]|[0-5][0-9])(?:[.:][0-9]+)?)? (?:PST|PDT|MST|MDT|CST|CDT|EST|EDT|UTC)`,
	"DATESTAMP_RFC2822":  `(?:Mon(?:day)?|Tue(?:sday)?|Wed(?:nesday)?|Thu(?:rsday)?|Fri(?:day)?|Sat(?:urday)?|Sun(?:day)?), (?:0?[1-9]|[12][0-9]|3[01]) (?:Jan(?:uary)?|Feb(?:ruary)?|Mar(?:ch)?|Apr(?:il)?|May|Jun(?:e)?|Jul(?:y)?|Aug(?:ust)?|Sep(?:tember)?|Sept|Oct(?:ober)?|Nov(?:ember)?|Dec(?:ember)?) (?:[0-9]{2})?[0-9]{2} (?:0?[0-9]|1[0-9]|2[0-3]):[0-5][0-9](?::(?:0?[0-9]|[0-5][0-9])(?:[.:][0-9]+)?)? (?:Z|[+-][0-9]{2}(?::?[0-9]{2})?)`,
	"DATESTAMP_OTHER":    `(?:Mon(?:day)?|Tue(?:sday)?|Wed(?:nesday)?|Thu(?:rsday)?|Fri(?:day)?|Sat(?:urday)?|Sun(?:day)?) (?:Jan(?:uary)?|Feb(?:ruary)?|Mar(?:ch)?|Apr(?:il)?|May|Jun(?:e)?|Jul(?:y)?|Aug(?:ust)?|Sep(?:tember)?|Sept|Oct(?:ober)?|Nov(?:ember)?|Dec(?:ember)?) (?:0?[1-9]|[12][0-9]|3[01]) (?:0?[0-9]|1[0-9]|2[0-3]):[0-5][0-9](?::(?:0?[0-9]|[0-5][0-9])(?:[.:][0-9]+)?)? (?:PST|PDT|MST|MDT|CST|CDT|EST|EDT|UTC) (?:[0-9]{2})?[0-9]{2}`,
	"DATESTAMP_EVENTLOG": `(?:[0-9]{2})?[0-9]{2}(?:0[1-9]|1[0-2])(?:0?[1-9]|[12][0-9]|3[01])(?:0?[0-9]|1[0-9]|2[0-3])[0-5][0-9](?:0?[0-9]|[0-5][0-9])`,

	"LOGLEVEL": `(?:[Aa]lert|ALERT|[Tt]race|TRACE|[Dd]ebug|DEBUG|[Nn]otice|NOTICE|[Ii]nfo|INFO|[Ww]arn(?:ing)?|WARN|WARNING|[Ee]rr(?:or)?|ERR|ERROR|[Cc]rit(?:ical)?|CRIT|CRITICAL|[Ff]atal|FATAL|[Ss]evere|SEVERE|[Ee]merg(?:ency)?|EMERG|EMERGENCY)`,
	// HTTPDATE's and SYSLOGTIMESTAMP's MonthName "Matches full or
	// abbreviated english month names (Example: \"Jan\" or \"January\")"
	// — the MONTH pattern's vocabulary.
	"HTTPDATE":        `(?:0?[1-9]|[12][0-9]|3[01])/(?:%{MONTH})/(?:[0-9]{2})?[0-9]{2}:(?:0?[0-9]|1[0-9]|2[0-3]):[0-5][0-9](?::(?:0?[0-9]|[0-5][0-9])(?:[.:][0-9]+)?)? [-+]?[0-9]+`,
	"SYSLOGTIMESTAMP": `(?:%{MONTH}) (?: ?[0-9]|[12][0-9]|3[01]) (?:0?[0-9]|1[0-9]|2[0-3]):[0-5][0-9](?::(?:0?[0-9]|[0-5][0-9])(?:[.:][0-9]+)?)?`,
	"PROG":            `[\w./%+-]+`,
	"SYSLOGFACILITY":  `<(?P<facility>[-+]?[0-9]+)\.(?P<priority>[-+]?[0-9]+)>`,
}

// grokAliases resolves the union patterns the table defines over other
// patterns; they expand after the primary table.
var grokAliases = map[string]string{
	"IP":           `(%{IPV6}|%{IPV4})`,
	"IPORHOST":     `(%{HOSTNAME}|%{IPV6}|%{IPV4})`,
	"MAC":          `(%{CISCOMAC}|%{WINDOWSMAC}|%{COMMONMAC})`,
	"PATH":         `(%{UNIXPATH}|%{WINPATH})`,
	"URIPATHPARAM": `(%{URIPATH}(?:%{URIPARAM})?)`,
	"URI":          `(%{URIPROTO}://(?:[^\s/@]+@)?%{IPORHOST}(?::[0-9]+)?(?:%{URIPATHPARAM})?)`,
	"HOSTPORT":     `(%{IPORHOST}:(?P<PORT>[0-9]+))`,
	// URIHOST's port is "a colon and a port number" (the documented
	// worked example: %{URIHOST:host} on example.com:443 captures
	// host=example.com:443 with port=443) — the colon cannot be part of
	// the host alternation, which excludes it.
	"URIHOST":    `(%{IPORHOST}(?::(?P<port>[0-9]+))?)`,
	"SYSLOGPROG": `(%{PROG}(?:\[(?P<pid>[-+]?[0-9]+)\])?)`,
	"SYSLOGHOST": `(%{HOSTNAME}|%{IPV6}|%{IPV4})`,
}

// grokCommonLogPatterns carries the three common-log composites with
// their documented typed outputs. They must be the first pattern in a
// match and accept only one trailing DATA/GREEDYDATA pattern; the
// structural rule is enforced by validateGrokMatch.
var grokCommonLogPatterns = map[string]string{
	"APACHE_ACCESS_LOG": `(?P<remote_host>%{IPORHOST}) \S+ \S+ \[(?P<__TS_timestamp>%{HTTPDATE})\] "(?P<http_method>[A-Z]+) (?P<request>[^ ]*) HTTP/(?P<http_version>[0-9.]+)" (?P<__NUM_status_code>[0-9]+) (?P<__NUM_response_size>[0-9]+|-)`,
	"NGINX_ACCESS_LOG":  `(?P<remote_host>%{IPORHOST}) \S+ (?P<auth_user>[^\s\[]+) \[(?P<__TS_timestamp>%{HTTPDATE})\] "(?P<http_method>[A-Z]+) (?P<request>[^ ]*) HTTP/(?P<http_version>[0-9.]+)" (?P<__NUM_status_code>[0-9]+) (?P<__NUM_response_size>[0-9]+|-)(?: "(?P<referrer>[^"]*)")?(?: "(?P<agent>[^"]*)")?`,
	"SYSLOG5424":        `<(?P<__NUM_pri>[0-9]+)>(?P<__NUM_version>[0-9]+) (?P<timestamp>%{TIMESTAMP_ISO8601}) (?P<hostname>\S+) (?P<app>\S+) (?P<__NIL_procid>\S+) (?P<msg_id>\S+) (?P<structured_data>\[[^\]]*\]|\-)(?P<message>(?:\[[^\]]*\])*.*)`,
}

// grokNumericPrefix marks the common-log composite groups whose
// documented outputs carry numbers rather than strings (a plain grok
// capture of the same field name stays a string).
const grokNumericPrefix = "__NUM_"

// grokTimestampPrefix marks the common-log composites' captured
// httpdate: the documented output renames the field to timestamp and
// converts the value to ISO-8601 UTC ("timestamp":
// "2023-08-03T12:34:56Z" for the documented [03/Aug/2023:12:34:56
// +0000]). A plain grok capture named timestamp keeps the raw text (the
// documented %{HTTPDATE:timestamp} example emits the httpdate verbatim).
const grokTimestampPrefix = "__TS_"

// grokNilPrefix marks a composite group whose protocol's nil token
// means the field is absent: RFC 5424 spells PROCID "-" when nil, and
// the documented SYSLOG5424 output omits procid for the sample whose
// PROCID is the nil token — a present PROCID is captured under the
// stripped name.
const grokNilPrefix = "__NIL_"

// httpdateToISO renders a captured httpdate as the ISO-8601
// UTC-normalised timestamp the common-log composites' documented outputs
// carry; an unparseable value falls back to the raw text. The layout's
// unpadded day, hour and second references accept both spellings — the
// HTTPDATE grammar admits one-digit day, hour and second (`0?[1-9]`
// classes), and the zero-padded references would reject exactly the
// values the accepting pattern matches.
func httpdateToISO(httpdate string) string {
	if t, err := time.Parse("2/Jan/2006:15:04:5 -0700", httpdate); err == nil {
		return t.UTC().Format("2006-01-02T15:04:05Z07:00")
	}
	return httpdate
}

var grokTokenRe = regexp.MustCompile(`%\{([A-Z0-9_]+)(?::([^}]*))?\}`)

// grokFieldNamePattern bounds a grok capture's FIELD_NAME: dot-separated
// segments of ASCII word characters. The dot is the documented JSON-path
// notation; the word-character vocabulary is what Go's named capture
// groups admit (dots are mangled to a placeholder before compilation).
// A name outside the vocabulary cannot compile, and compileGrok's
// failure path silently disables the processor — an accepted transformer
// that extracts nothing — so the match expression's validator rejects it
// up front.
var grokFieldNamePattern = regexp.MustCompile(`^[A-Za-z0-9_]+(?:\.[A-Za-z0-9_]+)*$`)

// grokFiveLimitPatterns are the patterns whose documented table row
// caps usage at five ("Maximum pattern limit: 5").
var grokFiveLimitPatterns = map[string]bool{
	"NOTSPACE": true, "SPACE": true, "DATA": true, "GREEDYDATA": true,
	"IPV6": true, "IP": true, "HOSTNAME": true, "HOST": true,
	"IPORHOST": true, "HOSTPORT": true, "URIHOST": true, "ARN": true,
	"WINPATH": true, "PATH": true, "URIPARAM": true, "URIPATHPARAM": true,
	"URI": true, "SYSLOGHOST": true,
}

// grokCombinationPatterns share one budget: "Any combination of the
// following patterns can be used as many as five times".
var grokCombinationPatterns = map[string]bool{
	"URI": true, "URIPARAM": true, "URIPATHPARAM": true, "SPACE": true,
	"DATA": true, "GREEDYDATA": true, "GREEDYDATA_MULTILINE": true,
}

// grokTrailingAfterComposite are the only patterns a common-log
// composite may be followed by, and only one of them ("you can follow
// them only with exactly one DATA. GREEDYDATA or GREEDYDATA_MULTILINE
// pattern"; the documented invalid examples reject every other
// follower and anything preceding the composite).
var grokTrailingAfterComposite = map[string]bool{
	"DATA": true, "GREEDYDATA": true, "GREEDYDATA_MULTILINE": true,
}

// grokMaxPatternsInMatch is the documented combination ceiling
// ("As many as 20 grok patterns can be combined to match a log event").
const grokMaxPatternsInMatch = 20

// validateGrokMatch enforces the documented match-expression structure:
// the total pattern count, the per-pattern and shared usage limits, and
// the common-log composites' first-and-alone rule.
func validateGrokMatch(match string) error {
	tokens := grokTokenRe.FindAllStringSubmatch(match, -1)
	if len(tokens) > grokMaxPatternsInMatch {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("A grok match can combine at most %d grok patterns", grokMaxPatternsInMatch), 400)
	}
	uses := map[string]int{}
	combinationUses := 0
	compositeAt := -1
	for i, token := range tokens {
		name := token[1]
		if _, known := grokExpansion(name); !known {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("Unknown grok pattern %q: only predefined grok patterns are supported", name), 400)
		}
		if field := token[2]; field != "" && !grokFieldNamePattern.MatchString(field) {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("Grok field name %q must be a dot-separated path whose segments use letters, digits and underscores only", field), 400)
		}
		uses[name]++
		if grokCombinationPatterns[name] {
			combinationUses++
		}
		if _, composite := grokCommonLogPatterns[name]; composite {
			if compositeAt >= 0 {
				return NewLogsError("InvalidParameterException",
					"A grok match can include only one common log format pattern", 400)
			}
			if i > 0 {
				return NewLogsError("InvalidParameterException",
					"A common log format pattern must be the first pattern in the match", 400)
			}
			compositeAt = i
		}
	}
	if combinationUses > 5 {
		return NewLogsError("InvalidParameterException",
			"The URI, URIPARAM, URIPATHPARAM, SPACE, DATA, GREEDYDATA and GREEDYDATA_MULTILINE patterns can be used at most five times in combination", 400)
	}
	for name, count := range uses {
		if name == "GREEDYDATA_MULTILINE" && count > 1 {
			return NewLogsError("InvalidParameterException",
				"The GREEDYDATA_MULTILINE pattern can be used at most once in a grok match", 400)
		}
		if grokFiveLimitPatterns[name] && count > 5 {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("The %s grok pattern can be used at most five times in a grok match", name), 400)
		}
	}
	if compositeAt >= 0 {
		trailing := tokens[compositeAt+1:]
		if len(trailing) > 1 || (len(trailing) == 1 && !grokTrailingAfterComposite[trailing[0][1]]) {
			return NewLogsError("InvalidParameterException",
				"A common log format pattern can be followed only by one DATA, GREEDYDATA or GREEDYDATA_MULTILINE pattern", 400)
		}
	}
	return nil
}

// compiledGrok carries the compiled match expression and the field
// names whose pattern was QUOTEDSTRING — those captures drop their
// surrounding quotes ("Output: {\"msg\": \"Hello, world!\"}" — the
// documented example's value is unquoted).
type compiledGrok struct {
	re           *regexp.Regexp
	quotedFields map[string]bool
}

// compileGrok translates one grok match expression into a compiled
// regular expression whose named groups carry the field names.
func compileGrok(match string) (*compiledGrok, bool) {
	quoted := map[string]bool{}
	expr := grokTokenRe.ReplaceAllStringFunc(match, func(token string) string {
		groups := grokTokenRe.FindStringSubmatch(token)
		patternName, field := groups[1], groups[2]
		expansion, ok := grokExpansion(patternName)
		if !ok {
			return regexp.QuoteMeta(token)
		}
		if field != "" {
			name := strings.ReplaceAll(field, ".", grokDotPlaceholder)
			if patternName == "QUOTEDSTRING" {
				quoted[name] = true
			}
			return "(?P<" + name + ">" + expansion + ")"
		}
		return "(?:" + expansion + ")"
	})
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, false
	}
	return &compiledGrok{re: re, quotedFields: quoted}, true
}

// grokCompileCache memoises compiled match expressions keyed by the
// match string: the ingestion hot path applies a recipe to every event
// of every batch, and recompiling an immutable expression per event
// (token walk, alias expansion, regexp.Compile) is pure repeated work.
// A failed compile caches as a nil entry — the string is immutable, so
// its outcome cannot change. The cache resets when full: distinct match
// strings are bounded by the configured recipes, and a reset only costs
// recompiles, never correctness.
const grokCompileCacheMax = 128

var (
	grokCompileMu    sync.Mutex
	grokCompileCache = map[string]*compiledGrok{}
)

// compiledGrokFor resolves a match expression through the compile cache.
func compiledGrokFor(match string) (*compiledGrok, bool) {
	grokCompileMu.Lock()
	cached, ok := grokCompileCache[match]
	grokCompileMu.Unlock()
	if ok {
		return cached, cached != nil
	}
	compiled, ok := compileGrok(match)
	grokCompileMu.Lock()
	if len(grokCompileCache) >= grokCompileCacheMax {
		grokCompileCache = map[string]*compiledGrok{}
	}
	if ok {
		grokCompileCache[match] = compiled
	} else {
		grokCompileCache[match] = nil
	}
	grokCompileMu.Unlock()
	return compiled, ok
}

// grokExpansion resolves one pattern name to its regular expression,
// expanding nested pattern references (aliases reference other
// patterns) and protecting against reference cycles.
func grokExpansion(name string, seen ...string) (string, bool) {
	for _, s := range seen {
		if s == name {
			return "", false
		}
	}
	if pattern, ok := grokCommonLogPatterns[name]; ok {
		return expandGrokRefs(pattern, append(seen, name)...), true
	}
	if pattern, ok := grokAliases[name]; ok {
		return expandGrokRefs(pattern, append(seen, name)...), true
	}
	if pattern, ok := grokPatterns[name]; ok {
		return expandGrokRefs(pattern, append(seen, name)...), true
	}
	return "", false
}

// expandGrokRefs recursively expands %{PATTERN} references inside an
// already-known expansion body.
func expandGrokRefs(body string, seen ...string) string {
	return grokTokenRe.ReplaceAllStringFunc(body, func(token string) string {
		groups := grokTokenRe.FindStringSubmatch(token)
		name := groups[1]
		if expansion, ok := grokExpansion(name, seen...); ok {
			return "(?:" + expansion + ")"
		}
		return regexp.QuoteMeta(token)
	})
}

// grokCaptureValue types one captured value: the composite numeric
// fields render as JSON numbers, everything else as the matched string.
func grokCaptureValue(name, value string) interface{} {
	if name == "structured_data" {
		// The documented SYSLOG5424 output carries the structured-data
		// content without its brackets.
		return strings.TrimSuffix(strings.TrimPrefix(value, "["), "]")
	}
	return value
}

// grokNumericCapture types a composite numeric group's value.
func grokNumericCapture(value string) interface{} {
	if n, err := strconv.ParseInt(value, 10, 64); err == nil {
		return n
	}
	return value
}

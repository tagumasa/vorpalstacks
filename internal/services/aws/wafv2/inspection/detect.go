package inspection

import (
	"regexp"
	"strings"
)

// byteMatchConstraint reports whether the transformed candidate
// satisfies the positional constraint against the search string. The
// five constraints follow the ByteMatchStatement documentation
// (PositionalConstraint enum: EXACTLY, STARTS_WITH, ENDS_WITH, CONTAINS,
// CONTAINS_WORD).
func byteMatchConstraint(candidate, search []byte, constraint string) bool {
	switch constraint {
	case "EXACTLY":
		return string(candidate) == string(search)
	case "STARTS_WITH":
		return len(candidate) >= len(search) && string(candidate[:len(search)]) == string(search)
	case "ENDS_WITH":
		return len(candidate) >= len(search) && string(candidate[len(candidate)-len(search):]) == string(search)
	case "CONTAINS":
		return strings.Contains(string(candidate), string(search))
	case "CONTAINS_WORD":
		return containsWord(candidate, search)
	}
	return false
}

// containsWord reports whether the search string occurs in the
// candidate delimited on both sides by word boundaries. A word
// boundary is any character other than a-z, A-Z, 0-9 and underscore.
func containsWord(candidate, search []byte) bool {
	if len(search) == 0 {
		return false
	}
	s, sub := string(candidate), string(search)
	for start := 0; start < len(s); {
		idx := strings.Index(s[start:], sub)
		if idx < 0 {
			return false
		}
		idx += start
		beforeOK := idx == 0 || !isWordByte(s[idx-1])
		after := idx + len(sub)
		afterOK := after == len(s) || !isWordByte(s[after])
		if beforeOK && afterOK {
			return true
		}
		start = idx + 1
	}
	return false
}

func isWordByte(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' || c == '_'
}

// SQL injection and cross-site scripting detection. AWS does not
// publish the token sets its SQLi and XSS match statements use; the
// patterns below are a conservative local approximation of the
// documented detection classes (SQL injection character sequences and
// keywords; script tags, event handlers and script URIs). Sensitivity
// level LOW applies only the high-confidence core patterns; HIGH adds
// the extended set, matching the documented trade-off between
// detection coverage and false positives.
var (
	sqliCorePatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(--|#|/\*)\s*$`),
		regexp.MustCompile(`(?i)\bunion\b\s+\bselect\b`),
		regexp.MustCompile(`(?i)\bselect\b\s+.*\bfrom\b`),
		regexp.MustCompile(`(?i)\binsert\b\s+\binto\b`),
		regexp.MustCompile(`(?i)\bdrop\b\s+\btable\b`),
		regexp.MustCompile(`(?i)\bdelete\b\s+\bfrom\b`),
		regexp.MustCompile(`(?i)\bupdate\b\s+\w+\s+\bset\b`),
		regexp.MustCompile(`(?i)['"]\s*or\s*['"]?\d+['"]?\s*=\s*['"]?\d+`),
		regexp.MustCompile(`(?i)\bwaitfor\b\s+\bdelay\b`),
		regexp.MustCompile(`(?i)\bsleep\s*\(`),
		regexp.MustCompile(`(?i)\bbenchmark\s*\(`),
		regexp.MustCompile(`(?i)\bexec(?:ute)?\b\s*\(`),
	}
	sqliExtendedPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)\bload_file\s*\(`),
		regexp.MustCompile(`(?i)\binformation_schema\b`),
		regexp.MustCompile(`(?i)\bpg_sleep\b`),
		regexp.MustCompile(`(?i)\bxp_cmdshell\b`),
		regexp.MustCompile(`(?i)\bchar\s*\(\s*\d+\s*\)`),
		regexp.MustCompile(`(?i)0x[0-9a-f]{6,}`),
		regexp.MustCompile(`(?i)'\s*=\s*'`),
		regexp.MustCompile(`(?i)\bbenchmark\b`),
		regexp.MustCompile(`(?i)\bhaving\b\s+\d+\s*=\s*\d+`),
	}
	xssPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)<\s*script\b`),
		regexp.MustCompile(`(?i)<\s*/\s*script\s*>`),
		regexp.MustCompile(`(?i)javascript\s*:`),
		regexp.MustCompile(`(?i)vbscript\s*:`),
		regexp.MustCompile(`(?i)\bon(?:error|load|click|mouseover|focus|blur|submit|input|change)\s*=`),
		regexp.MustCompile(`(?i)<\s*(?:iframe|embed|object|svg|img)\b[^>]*\bon\w+\s*=`),
		regexp.MustCompile(`(?i)\bexpression\s*\(`),
		regexp.MustCompile(`(?i)<\s*(?:iframe|embed|object)\b`),
		regexp.MustCompile(`(?i)\balert\s*\(`),
		regexp.MustCompile(`(?i)\bdocument\.cookie\b`),
		regexp.MustCompile(`(?i)\beval\s*\(`),
	}
)

// sqliDetected reports whether the candidate carries SQL injection
// patterns at the configured sensitivity level. The API model
// documents LOW as the default, so any level other than HIGH
// (including an unset value) uses the core patterns only.
func sqliDetected(candidate []byte, sensitivityLevel string) bool {
	if matchAnyPattern(candidate, sqliCorePatterns) {
		return true
	}
	if sensitivityLevel != "HIGH" {
		return false
	}
	return matchAnyPattern(candidate, sqliExtendedPatterns)
}

// xssDetected reports whether the candidate carries cross-site
// scripting patterns.
func xssDetected(candidate []byte) bool {
	return matchAnyPattern(candidate, xssPatterns)
}

func matchAnyPattern(candidate []byte, patterns []*regexp.Regexp) bool {
	for _, p := range patterns {
		if p.Match(candidate) {
			return true
		}
	}
	return false
}

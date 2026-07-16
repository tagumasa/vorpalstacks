// Package ntriples provides parsers for the N-Triples RDF serialisation format
// (https://www.w3.org/TR/n-triples/). Used by Neptune loader and NeptuneGraph
// import tasks to convert .nt source files into graph nodes and edges.
package ntriples

import "strings"

// ParseLine parses a single N-Triples line into its subject, predicate, and
// object components. Returns ok=false for blank lines or malformed input.
func ParseLine(line string) (subject, predicate, object string, ok bool) {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return "", "", "", false
	}
	if line[len(line)-1] == '.' {
		line = strings.TrimSpace(line[:len(line)-1])
	}

	var pos int

	subject, pos, ok = parseTerm(line, pos)
	if !ok {
		return "", "", "", false
	}
	for pos < len(line) && line[pos] == ' ' {
		pos++
	}

	predicate, pos, ok = parseTerm(line, pos)
	if !ok {
		return "", "", "", false
	}
	for pos < len(line) && line[pos] == ' ' {
		pos++
	}

	object, _, ok = parseTerm(line, pos)
	if !ok {
		return "", "", "", false
	}

	return subject, predicate, object, true
}

// ExtractLocalName strips angle brackets from a URI and returns the last
// segment after '/' or '#'.
func ExtractLocalName(uri string) string {
	uri = strings.Trim(uri, "<>")
	if idx := strings.LastIndex(uri, "/"); idx >= 0 {
		return uri[idx+1:]
	}
	if idx := strings.LastIndex(uri, "#"); idx >= 0 {
		return uri[idx+1:]
	}
	return uri
}

// IsAlphaNumHyphen reports whether c is an ASCII letter, digit, or hyphen.
func IsAlphaNumHyphen(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-'
}

func parseTerm(s string, pos int) (string, int, bool) {
	if pos >= len(s) {
		return "", pos, false
	}

	if s[pos] == '<' {
		end := strings.IndexByte(s[pos+1:], '>')
		if end < 0 {
			return "", pos, false
		}
		return s[pos : pos+end+2], pos + end + 2, true
	}

	if s[pos] == '"' {
		i := pos + 1
		for i < len(s) {
			if s[i] == '\\' && i+1 < len(s) {
				i += 2
				continue
			}
			if s[i] == '"' {
				break
			}
			i++
		}
		if i >= len(s) {
			return "", pos, false
		}
		end := i + 1
		if end < len(s) && s[end] == '@' {
			j := end + 1
			for j < len(s) && IsAlphaNumHyphen(s[j]) {
				j++
			}
			return s[pos:j], j, true
		}
		if end+1 < len(s) && s[end] == '^' && end+2 < len(s) && s[end+1] == '^' && end+2 < len(s) && s[end+2] == '<' {
			close := strings.IndexByte(s[end+2:], '>')
			if close >= 0 {
				return s[pos : end+2+close+1], end + 2 + close + 1, true
			}
		}
		return s[pos:end], end, true
	}

	if strings.HasPrefix(s[pos:], "_:") {
		i := pos + 2
		for i < len(s) && s[i] != ' ' && s[i] != '\t' && s[i] != '.' {
			i++
		}
		return s[pos:i], i, true
	}

	return "", pos, false
}

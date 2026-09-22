/*
 * Copyright 2026 Vorpalstacks
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package filterpattern

import (
	"errors"
	"regexp"
	"strings"
)

// ErrInvalidPattern reports a filter pattern the documented grammar
// rejects; the service layer maps it to InvalidParameterException so an
// invalid pattern never stores as a match-nothing filter.
var ErrInvalidPattern = errors.New("invalid filter pattern")

// The documented regex dialect of filter patterns: "Filter patterns
// with regex can only include the following:" alphanumeric characters;
// the symbol characters ': _ # = @ / ; , -'; and the operators
// '^ $ ? [ ] { } | \ * + .'. "The ( and ) operators are not supported.
// You cannot use parentheses to define a subpattern." "Multi-byte
// characters are not supported." — "%something!% would be rejected
// since '!' is not supported" (Filter and pattern syntax, CloudWatch
// Logs User Guide). Whitespace inside a span is tolerated as a
// separator artifact of the surrounding grammar, not a documented
// rejection.
const regexDialectSymbols = ":_#=@/;,-^$?[]{}|\\*+."

// MaxRegexPerPattern is the documented ceiling on regex spans inside
// one delimited or JSON filter pattern: "There is a limit of 2 regex
// for each filter pattern when creating a delimited or JSON filter
// pattern for metric filters and subscription filters or when filtering
// log events or Live Tail."
const MaxRegexPerPattern = 2

// The documented wildcard-selector quotas: "You can only use up to one
// wildcard selector in a property selector" and "You can use up to
// three wildcard selectors in a filter pattern with compound
// expressions".
const (
	MaxWildcardsPerSelector = 1
	MaxWildcardsPerPattern  = 3
)

// countSelectorWildcards counts the wildcard selectors one JSON property
// selector carries: the [*] array form and the .* object form.
func countSelectorWildcards(selector string) int {
	count := strings.Count(selector, "[*]")
	rest := strings.TrimPrefix(selector, "$.")
	for _, seg := range strings.Split(rest, ".") {
		name := seg
		if i := strings.Index(name, "["); i >= 0 {
			name = name[:i]
		}
		if name == "*" {
			count++
		}
	}
	return count
}

// extractJSONSelectors returns every JSON property selector ($.-prefixed)
// the pattern carries, walking quote-aware so quoted values never yield
// selectors.
func extractJSONSelectors(pattern string) []string {
	var selectors []string
	inSingle, inDouble := false, false
	for i := 0; i < len(pattern); i++ {
		c := pattern[i]
		switch {
		case c == '\'' && !inDouble:
			inSingle = !inSingle
		case c == '"':
			inDouble = !inDouble
		case c == '$' && !inSingle && !inDouble && i+1 < len(pattern) && pattern[i+1] == '.':
			j := i + 1
			for j < len(pattern) {
				r := pattern[j]
				if r == ' ' || r == '\t' || r == '=' || r == '!' || r == '<' || r == '>' ||
					r == ')' || r == '}' || r == '"' || r == '\'' || r == ',' {
					break
				}
				j++
			}
			if j > i+1 {
				selectors = append(selectors, pattern[i:j])
				i = j - 1
			}
		}
	}
	return selectors
}

func validRegexRune(r rune) bool {
	switch {
	case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
		return true
	case r == ' ' || r == '\t':
		return true
	case r < 0x80:
		return strings.ContainsRune(regexDialectSymbols, r)
	}
	return false
}

// scanRegexSpans walks the pattern with the evaluator's quote nesting
// — single quotes nest inside double quotes — and returns every
// unquoted %...% span. A percent sign with no closing partner before a
// quote or the end stays a literal character ("50% off"). Unbalanced
// quotes at the end report ErrInvalidPattern.
func scanRegexSpans(pattern string) ([]string, error) {
	var spans []string
	inSingle, inDouble := false, false
	for i := 0; i < len(pattern); i++ {
		c := pattern[i]
		switch {
		case c == '\'' && !inDouble:
			inSingle = !inSingle
		case c == '"':
			inDouble = !inDouble
		case c == '%' && !inSingle && !inDouble:
			j := i + 1
			for j < len(pattern) && pattern[j] != '%' && pattern[j] != '"' && pattern[j] != '\'' {
				j++
			}
			if j < len(pattern) && pattern[j] == '%' {
				spans = append(spans, pattern[i+1:j])
				i = j
			}
		}
	}
	if inSingle || inDouble {
		return nil, ErrInvalidPattern
	}
	return spans, nil
}

// ValidatePattern reports whether a filter pattern is syntactically
// valid per the documented grammar: balanced quotes, the delimited
// form's closing bracket, the JSON form's closing brace, every %regex%
// span inside the documented dialect and compilable, and at most two
// regex spans in a delimited or JSON pattern. An empty pattern is
// valid — the reads accept an absent pattern, and the requiredness of
// the member is each operation's own rule.
func ValidatePattern(pattern string) error {
	if pattern == "" {
		return nil
	}
	spans, err := scanRegexSpans(pattern)
	if err != nil {
		return err
	}
	trimmed := strings.TrimSpace(pattern)
	twoRegexCeiling := false
	switch {
	case strings.HasPrefix(trimmed, "["):
		twoRegexCeiling = true
		if !strings.Contains(trimmed, "]") {
			return ErrInvalidPattern
		}
	case strings.HasPrefix(trimmed, "{"):
		twoRegexCeiling = true
		if !strings.HasSuffix(trimmed, "}") {
			return ErrInvalidPattern
		}
	}
	if twoRegexCeiling && len(spans) > MaxRegexPerPattern {
		return ErrInvalidPattern
	}
	// The wildcard-selector quotas: "You can only use up to one wildcard
	// selector in a property selector" and "You can use up to three
	// wildcard selectors in a filter pattern with compound expressions".
	if strings.HasPrefix(trimmed, "{") {
		total := 0
		for _, selector := range extractJSONSelectors(pattern) {
			n := countSelectorWildcards(selector)
			if n > MaxWildcardsPerSelector {
				return ErrInvalidPattern
			}
			total += n
		}
		if total > MaxWildcardsPerPattern {
			return ErrInvalidPattern
		}
	}
	for _, span := range spans {
		for _, r := range span {
			if !validRegexRune(r) {
				return ErrInvalidPattern
			}
		}
		if _, err := regexp.Compile(span); err != nil {
			return ErrInvalidPattern
		}
	}
	return nil
}

// PatternContainsRegex reports whether the pattern carries at least one
// %regex% span — the census basis of the documented per-group quota:
// "There is a maximum of 5 filter patterns containing regex for each
// log group when creating metric filters or subscription filters."
func PatternContainsRegex(pattern string) bool {
	spans, err := scanRegexSpans(pattern)
	return err == nil && len(spans) > 0
}

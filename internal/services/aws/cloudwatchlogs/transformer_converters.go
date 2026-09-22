package cloudwatchlogs

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// The datatype converters: typeConverter and dateTimeConverter with
// their timezone, layout and localisation tables.

// --- datatype converters ---

func applyTypeConverter(record map[string]interface{}, cfg map[string]interface{}) {
	entries, _ := cfg["entries"].([]interface{})
	for _, raw := range entries {
		entry, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		key, _ := entry["key"].(string)
		targetType, _ := entry["type"].(string)
		if key == "" || targetType == "" {
			continue
		}
		value, exists := getPath(record, key)
		if !exists {
			continue
		}
		converted, ok := convertType(value, targetType)
		if !ok {
			continue
		}
		setPath(record, key, converted, true)
	}
}

func convertType(value interface{}, targetType string) (interface{}, bool) {
	switch targetType {
	case "integer":
		switch v := value.(type) {
		case int64:
			return v, true
		case float64:
			return int64(v), true
		case string:
			if n, err := strconv.ParseInt(v, 10, 64); err == nil {
				return n, true
			}
		case bool:
			if v {
				return int64(1), true
			}
			return int64(0), true
		}
	case "double":
		switch v := value.(type) {
		case int64:
			return float64(v), true
		case float64:
			return v, true
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f, true
			}
		}
	case "string":
		return fmt.Sprintf("%v", value), true
	case "boolean":
		switch v := value.(type) {
		case int64:
			return v != 0, true
		case bool:
			return v, true
		case string:
			if b, err := strconv.ParseBool(v); err == nil {
				return b, true
			}
		case float64:
			return v != 0, true
		}
	}
	return nil, false
}

// applyDateTimeConverter parses a datetime string under the configured
// match patterns and renders the target format.
func applyDateTimeConverter(record map[string]interface{}, cfg map[string]interface{}) {
	source, _ := cfg["source"].(string)
	target, _ := cfg["target"].(string)
	if source == "" || target == "" {
		return
	}
	var patterns []string
	if raw, ok := cfg["matchPatterns"].([]interface{}); ok {
		for _, p := range raw {
			if s, ok := p.(string); ok {
				patterns = append(patterns, s)
			}
		}
	}
	if len(patterns) == 0 {
		return
	}
	targetFormat := stringMember(cfg, "targetFormat", "yyyy-MM-dd'T'HH:mm:ss.SSS'Z'")
	sourceTZ, _ := cfg["sourceTimezone"].(string)
	targetTZ, _ := cfg["targetTimezone"].(string)
	locale, _ := cfg["locale"].(string)

	value, exists := getPath(record, source)
	if !exists {
		return
	}
	text, ok := value.(string)
	if !ok {
		return
	}
	srcLoc := loadTZ(sourceTZ)
	tgtLoc := loadTZ(targetTZ)
	parseText := localizeDatetimeText(text, locale)
	for _, pattern := range patterns {
		goPattern := javaToGoLayout(pattern)
		if t, err := time.ParseInLocation(goPattern, parseText, srcLoc); err == nil {
			rendered := localiseZoneName(t.In(tgtLoc).Format(javaToGoLayout(targetFormat)), locale)
			setPath(record, target, rendered, true)
			return
		}
	}
}

func loadTZ(name string) *time.Location {
	if name == "" {
		return time.UTC
	}
	if loc, err := time.LoadLocation(name); err == nil {
		return loc
	}
	return time.UTC
}

// javaToGoLayout translates the documented Java datetime letters the
// match patterns use into Go reference layouts, for the letters the
// documented patterns exercise. Single-quoted spans are literals in the
// Java syntax: their content passes through verbatim with the quotes
// removed (a doubled quote escapes one), so a quoted 'T' or 'Z' never
// re-enters the letter translation — the documented default
// targetFormat ends in the literal 'Z', and translating that Z as the
// zone letter corrupted every default render.
func javaToGoLayout(pattern string) string {
	replacements := []struct{ java, goLayout string }{
		{"yyyy", "2006"}, {"yy", "06"},
		{"MMMM", "January"}, {"MMM", "Jan"}, {"MM", "01"},
		{"EEEE", "Monday"}, {"EEE", "Mon"},
		{"dd", "02"},
		{"HH", "15"}, {"mm", "04"}, {"ss", "05"},
		// SSS is the bare fraction in the Java syntax — the separator
		// before it is the pattern's own literal, so the documented
		// "ss.SSS" translates to Go's ".000" fraction without injecting
		// a second dot.
		{"SSS", "000"},
		{"Z", "Z07:00"}, {"XXX", "Z07:00"},
		// z is Java's general zone-name letter: Go's reference
		// abbreviation (MST) renders it, localised afterwards when the
		// configured locale names the zone differently.
		{"z", "MST"},
	}
	var b strings.Builder
	inQuote := false
	for i := 0; i < len(pattern); i++ {
		c := pattern[i]
		if c == '\'' {
			if inQuote && i+1 < len(pattern) && pattern[i+1] == '\'' {
				b.WriteByte('\'')
				i++
				continue
			}
			inQuote = !inQuote
			continue
		}
		if inQuote {
			b.WriteByte(c)
			continue
		}
		matched := false
		for _, r := range replacements {
			if strings.HasPrefix(pattern[i:], r.java) {
				b.WriteString(r.goLayout)
				i += len(r.java) - 1
				matched = true
				break
			}
		}
		if !matched {
			b.WriteByte(c)
		}
	}
	return b.String()
}

// germanZoneNames carries the German zone abbreviations Go's reference
// render cannot produce: Java's z letter is locale-aware, and the
// documented German example targetFormat "yyyy-MM-dd'T'HH:mm:ss z"
// renders Europe/Berlin as "MEZ". Every zone and locale the documented
// examples do not exercise keeps Go's own abbreviation (the English
// short name Java itself uses).
var germanZoneNames = [][2]string{
	{"CEST", "MESZ"},
	{"CET", "MEZ"},
}

// localiseZoneName translates a rendered zone abbreviation into the
// configured locale's zone name.
func localiseZoneName(rendered, locale string) string {
	if strings.ToLower(locale) != "de" {
		return rendered
	}
	for _, pair := range germanZoneNames {
		rendered = replaceWholeWord(rendered, pair[0], pair[1])
	}
	return rendered
}

// stringMember reads a string member with a default.
func stringMember(m map[string]interface{}, key, def string) string {
	if v, ok := m[key].(string); ok && v != "" {
		return v
	}
	return def
}

// localizedDatetimeNames carries the month and day names of the locales
// the documented datetimeConverter examples exercise; a localized name
// parses after translation to Go's English reference names.
var localizedDatetimeNames = map[string][][2]string{
	"de": {
		{"Januar", "January"}, {"Februar", "February"}, {"März", "March"}, {"Mai", "May"},
		{"Juni", "June"}, {"Juli", "July"}, {"Oktober", "October"}, {"Dezember", "December"},
		{"Samstag", "Saturday"}, {"Sonntag", "Sunday"}, {"Montag", "Monday"}, {"Dienstag", "Tuesday"},
		{"Mittwoch", "Wednesday"}, {"Donnerstag", "Thursday"}, {"Freitag", "Friday"},
	},
	"fr": {
		{"janvier", "January"}, {"février", "February"}, {"mars", "March"}, {"avril", "April"},
		{"mai", "May"}, {"juin", "June"}, {"juillet", "July"}, {"août", "August"},
		{"septembre", "September"}, {"octobre", "October"}, {"novembre", "November"}, {"décembre", "December"},
	},
	"es": {
		{"enero", "January"}, {"febrero", "February"}, {"marzo", "March"}, {"abril", "April"},
		{"mayo", "May"}, {"junio", "June"}, {"julio", "July"}, {"agosto", "August"},
		{"septiembre", "September"}, {"octubre", "October"}, {"noviembre", "November"}, {"diciembre", "December"},
	},
}

// localizeDatetimeText translates a text's localized month and day
// names to their English parse names (whole-word, case-sensitive).
func localizeDatetimeText(text, locale string) string {
	replacements, ok := localizedDatetimeNames[strings.ToLower(locale)]
	if !ok {
		return text
	}
	result := text
	for _, pair := range replacements {
		result = replaceWholeWord(result, pair[0], pair[1])
	}
	return result
}

// replaceWholeWord replaces case-sensitive whole words: a word is a
// maximal run of Unicode letters, scanned by rune — the locale tables
// carry month names whose letters reach outside ASCII (März, février,
// août, décembre), and a byte-class boundary would end those words at
// the first non-ASCII letter, leaving the entry unreachable.
func replaceWholeWord(text, from, to string) string {
	var b strings.Builder
	wordStart := -1
	for i, r := range text {
		if unicode.IsLetter(r) {
			if wordStart < 0 {
				wordStart = i
			}
			continue
		}
		if wordStart >= 0 {
			if text[wordStart:i] == from {
				b.WriteString(to)
			} else {
				b.WriteString(text[wordStart:i])
			}
			wordStart = -1
		}
		b.WriteRune(r)
	}
	if wordStart >= 0 {
		if text[wordStart:] == from {
			b.WriteString(to)
		} else {
			b.WriteString(text[wordStart:])
		}
	}
	return b.String()
}

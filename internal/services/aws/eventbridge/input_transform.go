package eventbridge

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"

	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// The input transformer follows the EventBridge input transformation
// contract. Paths use the documented JSONPath subset — dot notation with
// keys of alphanumerics, dashes, underscores and forward slashes, array
// indices and wildcards ("Supported syntax includes: dot notation, dashes,
// underscores, alphanumeric characters, array indices, wildcards (*),
// forward slashes"). Placeholders are <name>; a variable whose path
// resolves to nothing is elided from the output, string values are quoted
// when substituted in a JSON value position ("EventBridge automatically
// adds quotes to string variable values during transformation... does not
// add quotes to variables that represent JSON objects or arrays"), and a
// template that does not become valid JSON is delivered as a JSON string.

// Reserved template variables supplied by the platform (the "Predefined
// variables" section): the rule context, the ingest timestamp, and the
// event with and without its detail member. Names are matched verbatim.
const (
	reservedVarRuleARN       = "aws.events.rule-arn"
	reservedVarRuleName      = "aws.events.rule-name"
	reservedVarEvent         = "aws.events.event"
	reservedVarEventJSON     = "aws.events.event.json"
	reservedVarIngestionTime = "aws.events.event.ingestion-time"
)

// pathStepKind discriminates the steps of a parsed event path.
type pathStepKind int

const (
	pathStepKey pathStepKind = iota
	pathStepIndex
	pathStepWildcard
)

type pathStep struct {
	kind  pathStepKind
	key   string
	index int
}

// isEventPathKeyChar reports whether c may appear in a path key segment:
// alphanumeric, dash, underscore or forward slash. Dots separate segments
// and brackets carry indices, so they cannot be part of a key.
func isEventPathKeyChar(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' ||
		c == '_' || c == '-' || c == '/'
}

// parseEventPath parses the documented JSONPath subset. The leading $ is
// optional; every following segment is a dot-prefixed key, an array index
// bracket, or a wildcard. A path outside the subset does not parse, and an
// unparsable path never resolves (there is no validation when creating
// JSON paths, so delivery tolerates them).
func parseEventPath(path string) ([]pathStep, bool) {
	rest := strings.TrimPrefix(path, "$")
	var steps []pathStep
	for len(rest) > 0 {
		if rest[0] == '.' {
			rest = rest[1:]
			if len(rest) == 0 {
				return nil, false
			}
		}
		switch {
		case rest[0] == '*':
			steps = append(steps, pathStep{kind: pathStepWildcard})
			rest = rest[1:]
		case rest[0] == '[':
			closing := strings.IndexByte(rest, ']')
			if closing < 0 {
				return nil, false
			}
			inner := rest[1:closing]
			rest = rest[closing+1:]
			if inner == "*" {
				steps = append(steps, pathStep{kind: pathStepWildcard})
				continue
			}
			idx, err := strconv.Atoi(inner)
			if err != nil || idx < 0 || inner != strconv.Itoa(idx) {
				return nil, false
			}
			steps = append(steps, pathStep{kind: pathStepIndex, index: idx})
		default:
			i := 0
			for i < len(rest) && isEventPathKeyChar(rest[i]) {
				i++
			}
			if i == 0 {
				return nil, false
			}
			steps = append(steps, pathStep{kind: pathStepKey, key: rest[:i]})
			rest = rest[i:]
			if len(rest) > 0 && rest[0] != '.' && rest[0] != '[' {
				return nil, false
			}
		}
	}
	return steps, true
}

// resolveEventPath walks the parsed steps. A wildcard over an array yields
// its elements; over an object it yields the member values in key order —
// the collected array keeps the result a single JSON value. Any step that
// does not apply to the current value leaves the path unresolved.
func resolveEventPath(root interface{}, steps []pathStep) (interface{}, bool) {
	current := root
	for _, step := range steps {
		switch step.kind {
		case pathStepKey:
			m, ok := current.(map[string]interface{})
			if !ok {
				return nil, false
			}
			value, exists := m[step.key]
			if !exists {
				return nil, false
			}
			current = value
		case pathStepIndex:
			switch v := current.(type) {
			case []interface{}:
				if step.index >= len(v) {
					return nil, false
				}
				current = v[step.index]
			case []string:
				// In-process envelopes carry resources as []string;
				// envelopes that crossed a JSON hop carry []interface{}.
				if step.index >= len(v) {
					return nil, false
				}
				current = v[step.index]
			default:
				return nil, false
			}
		case pathStepWildcard:
			switch v := current.(type) {
			case []interface{}:
				current = append([]interface{}(nil), v...)
			case []string:
				widened := make([]interface{}, len(v))
				for i, s := range v {
					widened[i] = s
				}
				current = widened
			case map[string]interface{}:
				keys := make([]string, 0, len(v))
				for k := range v {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				values := make([]interface{}, 0, len(v))
				for _, k := range keys {
					values = append(values, v[k])
				}
				current = values
			default:
				return nil, false
			}
		}
	}
	return current, true
}

// resolveEventPathString resolves a path expression against an event
// envelope in one step.
func resolveEventPathString(envelope map[string]interface{}, path string) (interface{}, bool) {
	steps, ok := parseEventPath(path)
	if !ok {
		return nil, false
	}
	return resolveEventPath(envelope, steps)
}

// transformVariable is one substitutable <name>: its resolved value and
// whether the path resolved at all. An unresolved variable is elided from
// the output ("that variable isn't created and won't appear in the
// output").
type transformVariable struct {
	value interface{}
	found bool
}

// applyInputTransform renders the payload a transformer target receives.
// The template is first substituted in JSON mode, where string values are
// quoted in value positions and left raw inside string literals; if the
// result is valid JSON it is the payload. Otherwise the template is a
// string template: it is unescaped as a JSON string literal, its values
// substituted raw, and the resulting text delivered as a JSON string.
func (s *EventsService) applyInputTransform(ruleARN, ruleName string, event *eventsstore.Event, transformer *eventsstore.InputTransformer) []byte {
	if transformer == nil || transformer.InputTemplate == "" {
		return marshalEventEnvelope(event)
	}

	envelope := eventEnvelope(event)
	vars := make(map[string]*transformVariable, len(transformer.InputPathsMap)+5)
	for name, path := range transformer.InputPathsMap {
		value, found := resolveEventPathString(envelope, path)
		vars[name] = &transformVariable{value: value, found: found}
	}

	eventWithoutDetail := make(map[string]interface{}, len(envelope))
	for k, v := range envelope {
		eventWithoutDetail[k] = v
	}
	delete(eventWithoutDetail, "detail")

	// aws.events.event.ingestion-time is "The time at which the event
	// was received by EventBridge. This is an ISO 8601 timestamp. This
	// variable is generated by EventBridge and can't be overwritten"
	// (input transformation page, predefined variables): the platform's
	// receipt stamp, never the publisher-suppliable event Time. Records
	// that predate the stamp (pre-index archive replays, synthetic
	// events) fall back to Time, their closest surviving receipt.
	ingestionTime := event.IngestionTime
	if ingestionTime.IsZero() {
		ingestionTime = event.Time
	}
	vars[reservedVarRuleARN] = &transformVariable{value: ruleARN, found: true}
	vars[reservedVarRuleName] = &transformVariable{value: ruleName, found: true}
	vars[reservedVarEvent] = &transformVariable{value: eventWithoutDetail, found: true}
	vars[reservedVarEventJSON] = &transformVariable{value: envelope, found: true}
	vars[reservedVarIngestionTime] = &transformVariable{value: formatEventTime(ingestionTime), found: true}

	template := transformer.InputTemplate
	if substituted := substituteTemplate(template, vars, true); json.Valid([]byte(substituted)) {
		return []byte(substituted)
	}
	unescaped := unescapeJSONStringLiteral(template)
	payload, err := json.Marshal(substituteTemplate(unescaped, vars, false))
	if err != nil {
		return marshalEventEnvelope(event)
	}
	return payload
}

// substituteTemplate replaces every <name> placeholder whose name is a
// known variable. jsonMode selects JSON-template semantics: values in
// value positions are JSON-encoded (strings quoted), while values inside
// string literals are inserted raw with object and array quotes stripped
// ("EventBridge removes any internal quotes to ensure a valid string").
// Outside jsonMode every value takes the raw form. An unknown name keeps
// its literal text — template placeholders carry no acceptance
// validation.
func substituteTemplate(template string, vars map[string]*transformVariable, jsonMode bool) string {
	var b strings.Builder
	b.Grow(len(template))
	inString := false
	i := 0
	for i < len(template) {
		c := template[i]
		if c == '<' {
			end := strings.IndexByte(template[i+1:], '>')
			if end >= 0 {
				name := template[i+1 : i+1+end]
				if variable, known := vars[name]; known {
					b.WriteString(renderVariable(variable, jsonMode && !inString))
					i += end + 2
					continue
				}
			}
			b.WriteByte(c)
			i++
			continue
		}
		if inString {
			if c == '\\' && i+1 < len(template) {
				b.WriteByte(c)
				b.WriteByte(template[i+1])
				i += 2
				continue
			}
			if c == '"' {
				inString = false
			}
			b.WriteByte(c)
			i++
			continue
		}
		if c == '"' {
			inString = true
		}
		b.WriteByte(c)
		i++
	}
	return b.String()
}

// renderVariable renders one substitution. In a JSON value position every
// value is JSON-encoded, which supplies the automatic quoting for strings;
// in a string context the raw text is inserted and object and array values
// have their quotes stripped.
func renderVariable(variable *transformVariable, jsonValuePosition bool) string {
	if !variable.found {
		return ""
	}
	if jsonValuePosition {
		b, err := json.Marshal(variable.value)
		if err != nil {
			return ""
		}
		return string(b)
	}
	if str, ok := variable.value.(string); ok {
		return str
	}
	b, err := json.Marshal(variable.value)
	if err != nil {
		return ""
	}
	return strings.ReplaceAll(string(b), `"`, "")
}

// unescapeJSONStringLiteral resolves the backslash escapes a string
// template carries ("To have the InputTemplate include quote marks within
// a JSON string, escape each quote marks with a slash"). Unrecognised
// escapes are kept verbatim.
func unescapeJSONStringLiteral(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if s[i] != '\\' || i+1 >= len(s) {
			b.WriteByte(s[i])
			continue
		}
		i++
		switch s[i] {
		case '"':
			b.WriteByte('"')
		case '\\':
			b.WriteByte('\\')
		case '/':
			b.WriteByte('/')
		case 'b':
			b.WriteByte('\b')
		case 'f':
			b.WriteByte('\f')
		case 'n':
			b.WriteByte('\n')
		case 'r':
			b.WriteByte('\r')
		case 't':
			b.WriteByte('\t')
		case 'u':
			if i+4 < len(s) {
				if code, err := strconv.ParseUint(s[i+1:i+5], 16, 32); err == nil {
					b.WriteRune(rune(code))
					i += 4
					continue
				}
			}
			b.WriteByte('\\')
			b.WriteByte(s[i])
		default:
			b.WriteByte('\\')
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

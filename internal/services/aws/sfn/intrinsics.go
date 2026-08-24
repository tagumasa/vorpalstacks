package sfn

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"

	sfnstore "vorpalstacks/internal/store/aws/sfn"
)

// Intrinsic-function evaluation for JSONPath states. Intrinsics appear as the
// value of a ".$"-suffixed payload-template key and, in the Fail state, as the
// ErrorPath/CausePath members. Every rule below follows the
// intrinsic-functions documentation: argument validation, the 10-level
// nesting bound, the reserved-character escapes and the per-function limits.

// errIntrinsicEscape reports a backslash that is not part of one of the four
// documented escape sequences (\' \{ \} \\) — the documentation calls this an
// open escape backslash and requires a runtime error.
var errIntrinsicEscape = errors.New("open escape backslash in intrinsic invocation string")

// formatPlaceholder is a private-use rune standing in for one "{}"
// substitution slot of a States.Format template while escapes are decoded.
const formatPlaceholder = '\uE000'

// intrinsicFunctions is the documented intrinsic set (the same list the
// J2119 statelint grammar accepts inside payload-template values).
var intrinsicFunctions = map[string]bool{
	"Array": true, "ArrayPartition": true, "ArrayContains": true,
	"ArrayRange": true, "ArrayGetItem": true, "ArrayLength": true,
	"ArrayUnique": true, "Base64Encode": true, "Base64Decode": true,
	"Hash": true, "JsonMerge": true, "JsonToString": true,
	"StringToJson": true, "MathRandom": true, "MathAdd": true,
	"StringSplit": true, "UUID": true, "Format": true,
}

type intrinsicArgKind int

const (
	argIntrinsic     intrinsicArgKind = iota // nested States.Name(...)
	argPath                                  // $.x or $$.Context.X path
	argStringLiteral                         // '...' with escapes
	argJSONLiteral                           // bare number, boolean or null
)

type intrinsicArg struct {
	kind  intrinsicArgKind
	text  string      // invocation or path text, or the raw literal content
	value interface{} // decoded value for JSON literals
}

// formatTemplate is a decoded States.Format first argument: the template with
// each live "{}" pair replaced by formatPlaceholder and the number of those
// pairs. Escaped braces (\{ \}) survive as literal braces and never act as
// placeholders.
type formatTemplate struct {
	marked       string
	placeholders int
}

// evaluateIntrinsic evaluates one intrinsic invocation against the state
// input. depth counts the invocation itself as level one so that a chain of
// ten nested functions is the documented limit.
func (e *Executor) evaluateIntrinsic(taskToken string, invocation string, data interface{}, depth int) (interface{}, error) {
	if depth > sfnstore.MaxIntrinsicNesting {
		return nil, fmt.Errorf("intrinsic nesting exceeds %d levels", sfnstore.MaxIntrinsicNesting)
	}
	name, args, err := parseIntrinsicInvocation(invocation)
	if err != nil {
		return nil, err
	}
	if !intrinsicFunctions[name] {
		return nil, fmt.Errorf("States.%s is not a supported intrinsic function", name)
	}
	vals := make([]interface{}, len(args))
	for idx, arg := range args {
		switch arg.kind {
		case argIntrinsic:
			v, err := e.evaluateIntrinsic(taskToken, arg.text, data, depth+1)
			if err != nil {
				return nil, err
			}
			vals[idx] = v
		case argPath:
			v, found, err := e.resolveIntrinsicPathArg(taskToken, arg.text, data)
			if err != nil {
				return nil, err
			}
			if !found {
				return nil, fmt.Errorf("path %s in intrinsic argument selected no value", arg.text)
			}
			vals[idx] = v
		case argStringLiteral:
			if name == "Format" && idx == 0 {
				tmpl, err := decodeIntrinsicTemplateLiteral(arg.text)
				if err != nil {
					return nil, err
				}
				vals[idx] = tmpl
			} else {
				decoded, err := decodeIntrinsicString(arg.text)
				if err != nil {
					return nil, err
				}
				vals[idx] = decoded
			}
		default:
			vals[idx] = arg.value
		}
	}
	return applyIntrinsic(name, vals)
}

// resolveIntrinsicPathArg resolves a $. or $$. argument path against the
// state input and the context object. Array roots accept bare and bracketed
// numeric segments.
func (e *Executor) resolveIntrinsicPathArg(taskToken, path string, data interface{}) (interface{}, bool, error) {
	if strings.HasPrefix(path, "$$.") {
		v, err := e.getContextValue(taskToken, path)
		if err != nil {
			return nil, false, err
		}
		return v, true, nil
	}
	if v, exists := getJSONPathValueAny(data, path); exists {
		return v, true, nil
	}
	return nil, false, nil
}

// parseIntrinsicInvocation splits a States.Name(args...) invocation into its
// name and classified arguments, validating the escape grammar as it goes.
func parseIntrinsicInvocation(invocation string) (string, []intrinsicArg, error) {
	const prefix = "States."
	if !strings.HasPrefix(invocation, prefix) {
		return "", nil, fmt.Errorf("not an intrinsic invocation: %s", invocation)
	}
	open := strings.IndexByte(invocation, '(')
	if open < 0 {
		return "", nil, fmt.Errorf("malformed intrinsic invocation: %s", invocation)
	}
	name := strings.TrimSpace(invocation[len(prefix):open])
	if name == "" || strings.ContainsAny(name, " \t\r\n") {
		return "", nil, fmt.Errorf("malformed intrinsic function name in: %s", invocation)
	}
	args := []intrinsicArg{}
	i := open + 1
	for {
		i = skipIntrinsicSpace(invocation, i)
		if i >= len(invocation) {
			return "", nil, fmt.Errorf("unterminated intrinsic invocation: %s", invocation)
		}
		if invocation[i] == ')' {
			break
		}
		arg, next, err := parseIntrinsicArg(invocation, i)
		if err != nil {
			return "", nil, err
		}
		args = append(args, arg)
		i = skipIntrinsicSpace(invocation, next)
		if i >= len(invocation) {
			return "", nil, fmt.Errorf("unterminated intrinsic invocation: %s", invocation)
		}
		switch invocation[i] {
		case ',':
			i++
		case ')':
		default:
			return "", nil, fmt.Errorf("unexpected character %q in intrinsic invocation: %s", invocation[i], invocation)
		}
	}
	i++
	if strings.TrimSpace(invocation[i:]) != "" {
		return "", nil, fmt.Errorf("unexpected text after intrinsic invocation: %s", invocation)
	}
	return name, args, nil
}

// parseIntrinsicArg parses one argument starting at position i and returns
// the classified argument plus the position after it.
func parseIntrinsicArg(invocation string, i int) (intrinsicArg, int, error) {
	switch {
	case invocation[i] == '\'':
		return parseIntrinsicStringArg(invocation, i)
	case strings.HasPrefix(invocation[i:], "States."):
		return parseIntrinsicNestedArg(invocation, i)
	case invocation[i] == '$':
		return parseIntrinsicPathArg(invocation, i)
	default:
		return parseIntrinsicJSONArg(invocation, i)
	}
}

// parseIntrinsicStringArg parses a single-quoted literal. The returned text
// holds the raw content between the quotes; escapes are decoded later so
// States.Format can distinguish escaped braces from placeholders.
func parseIntrinsicStringArg(invocation string, i int) (intrinsicArg, int, error) {
	j := i + 1
	for j < len(invocation) {
		switch invocation[j] {
		case '\\':
			if j+1 >= len(invocation) {
				return intrinsicArg{}, 0, errIntrinsicEscape
			}
			j += 2
		case '\'':
			return intrinsicArg{kind: argStringLiteral, text: invocation[i+1 : j]}, j + 1, nil
		default:
			j++
		}
	}
	return intrinsicArg{}, 0, fmt.Errorf("unterminated string literal in intrinsic invocation: %s", invocation)
}

// parseIntrinsicNestedArg parses a nested invocation by scanning to its
// matching close paren, skipping quoted strings so parens and commas inside
// literals do not terminate the scan.
func parseIntrinsicNestedArg(invocation string, i int) (intrinsicArg, int, error) {
	open := strings.IndexByte(invocation[i:], '(')
	if open < 0 {
		return intrinsicArg{}, 0, fmt.Errorf("malformed nested intrinsic invocation: %s", invocation[i:])
	}
	depth := 0
	k := i + open
	for k < len(invocation) {
		c := invocation[k]
		if c == '\'' {
			k++
			for k < len(invocation) {
				if invocation[k] == '\\' {
					k += 2
					continue
				}
				if invocation[k] == '\'' {
					break
				}
				k++
			}
			if k >= len(invocation) {
				return intrinsicArg{}, 0, fmt.Errorf("unterminated string literal in intrinsic invocation: %s", invocation)
			}
			k++
			continue
		}
		if c == '(' {
			depth++
		} else if c == ')' {
			depth--
			if depth == 0 {
				return intrinsicArg{kind: argIntrinsic, text: invocation[i : k+1]}, k + 1, nil
			}
		}
		k++
	}
	return intrinsicArg{}, 0, fmt.Errorf("unbalanced parentheses in intrinsic invocation: %s", invocation[i:])
}

// parseIntrinsicPathArg parses a JSONPath argument. Quoted segments (square
// bracket notation) may contain commas, parens and quotes, so the scan skips
// over them.
func parseIntrinsicPathArg(invocation string, i int) (intrinsicArg, int, error) {
	j := i
	for j < len(invocation) && invocation[j] != ',' && invocation[j] != ')' {
		if invocation[j] == '\'' {
			j++
			for j < len(invocation) {
				if invocation[j] == '\\' {
					j += 2
					continue
				}
				if invocation[j] == '\'' {
					break
				}
				j++
			}
		}
		j++
	}
	return intrinsicArg{kind: argPath, text: invocation[i:j]}, j, nil
}

// parseIntrinsicJSONArg parses a bare number, boolean or null literal.
func parseIntrinsicJSONArg(invocation string, i int) (intrinsicArg, int, error) {
	j := i
	for j < len(invocation) && invocation[j] != ',' && invocation[j] != ')' {
		j++
	}
	text := strings.TrimSpace(invocation[i:j])
	var value interface{}
	if err := json.Unmarshal([]byte(text), &value); err != nil {
		return intrinsicArg{}, 0, fmt.Errorf("invalid literal argument %q in intrinsic invocation", text)
	}
	switch value.(type) {
	case float64, bool, nil:
	default:
		return intrinsicArg{}, 0, fmt.Errorf("argument %q must be a quoted string, a path, or a number, boolean or null literal", text)
	}
	return intrinsicArg{kind: argJSONLiteral, value: value}, j, nil
}

func skipIntrinsicSpace(s string, i int) int {
	for i < len(s) && (s[i] == ' ' || s[i] == '\t' || s[i] == '\r' || s[i] == '\n') {
		i++
	}
	return i
}

// decodeIntrinsicString decodes the four documented escape sequences in a
// string literal argument.
func decodeIntrinsicString(raw string) (string, error) {
	var b strings.Builder
	for i := 0; i < len(raw); i++ {
		c := raw[i]
		if c != '\\' {
			b.WriteByte(c)
			continue
		}
		i++
		if i >= len(raw) {
			return "", errIntrinsicEscape
		}
		switch raw[i] {
		case '\'', '{', '}', '\\':
			b.WriteByte(raw[i])
		default:
			return "", errIntrinsicEscape
		}
	}
	return b.String(), nil
}

// decodeIntrinsicTemplateLiteral decodes a States.Format template literal.
// Escaped braces become literal braces; each unescaped "{}" pair becomes a
// placeholder marker.
func decodeIntrinsicTemplateLiteral(raw string) (formatTemplate, error) {
	var b strings.Builder
	placeholders := 0
	pendingOpen := false
	flushOpen := func() {
		if pendingOpen {
			b.WriteByte('{')
			pendingOpen = false
		}
	}
	for i := 0; i < len(raw); i++ {
		c := raw[i]
		if c == '\\' {
			i++
			if i >= len(raw) {
				return formatTemplate{}, errIntrinsicEscape
			}
			switch raw[i] {
			case '\'', '{', '}', '\\':
				flushOpen()
				b.WriteByte(raw[i])
			default:
				return formatTemplate{}, errIntrinsicEscape
			}
			continue
		}
		switch {
		case c == '{':
			flushOpen()
			pendingOpen = true
		case c == '}' && pendingOpen:
			b.WriteRune(formatPlaceholder)
			placeholders++
			pendingOpen = false
		default:
			flushOpen()
			b.WriteByte(c)
		}
	}
	flushOpen()
	return formatTemplate{marked: b.String(), placeholders: placeholders}, nil
}

// markFormatTemplate converts a template that arrived as a resolved value (a
// path, a nested intrinsic or any non-literal source) — such templates have
// no escape layer, so every "{}" pair is a placeholder.
func markFormatTemplate(s string) formatTemplate {
	return formatTemplate{
		marked:       strings.ReplaceAll(s, "{}", string(formatPlaceholder)),
		placeholders: strings.Count(s, "{}"),
	}
}

// applyIntrinsic validates arity and arguments and computes the function
// result.
func applyIntrinsic(name string, vals []interface{}) (interface{}, error) {
	switch name {
	case "Array":
		return append([]interface{}{}, vals...), nil

	case "ArrayPartition":
		if err := intrinsicArity(name, len(vals), 2, 2); err != nil {
			return nil, err
		}
		arr, err := intrinsicArgArray(name, vals[0], 1)
		if err != nil {
			return nil, err
		}
		if err := intrinsicArrayWithinPayloadLimit(name, arr); err != nil {
			return nil, err
		}
		chunk, err := intrinsicRoundedInt(name, vals[1], 2)
		if err != nil {
			return nil, err
		}
		if chunk <= 0 {
			return nil, fmt.Errorf("States.ArrayPartition chunk size must be a non-zero positive integer")
		}
		out := make([]interface{}, 0, (len(arr)+int(chunk)-1)/int(chunk))
		for start := 0; start < len(arr); start += int(chunk) {
			end := start + int(chunk)
			if end > len(arr) {
				end = len(arr)
			}
			out = append(out, append([]interface{}{}, arr[start:end]...))
		}
		return out, nil

	case "ArrayContains":
		if err := intrinsicArity(name, len(vals), 2, 2); err != nil {
			return nil, err
		}
		arr, err := intrinsicArgArray(name, vals[0], 1)
		if err != nil {
			return nil, err
		}
		if err := intrinsicArrayWithinPayloadLimit(name, arr); err != nil {
			return nil, err
		}
		for _, item := range arr {
			if jsonValuesEqual(item, vals[1]) {
				return true, nil
			}
		}
		return false, nil

	case "ArrayRange":
		if err := intrinsicArity(name, len(vals), 3, 3); err != nil {
			return nil, err
		}
		first, err := intrinsicRoundedInt(name, vals[0], 1)
		if err != nil {
			return nil, err
		}
		last, err := intrinsicRoundedInt(name, vals[1], 2)
		if err != nil {
			return nil, err
		}
		step, err := intrinsicRoundedInt(name, vals[2], 3)
		if err != nil {
			return nil, err
		}
		if step == 0 {
			return nil, fmt.Errorf("States.ArrayRange step must be non-zero")
		}
		out := []interface{}{}
		for v := first; (step > 0 && v <= last) || (step < 0 && v >= last); v += step {
			out = append(out, float64(v))
			if len(out) > sfnstore.MaxArrayRangeElements {
				return nil, fmt.Errorf("States.ArrayRange result exceeds %d elements", sfnstore.MaxArrayRangeElements)
			}
		}
		return out, nil

	case "ArrayGetItem":
		if err := intrinsicArity(name, len(vals), 2, 2); err != nil {
			return nil, err
		}
		arr, err := intrinsicArgArray(name, vals[0], 1)
		if err != nil {
			return nil, err
		}
		index, err := intrinsicRoundedInt(name, vals[1], 2)
		if err != nil {
			return nil, err
		}
		if index < 0 || index >= int64(len(arr)) {
			return nil, fmt.Errorf("States.ArrayGetItem index %d is out of range for an array of %d elements", index, len(arr))
		}
		return arr[index], nil

	case "ArrayLength":
		if err := intrinsicArity(name, len(vals), 1, 1); err != nil {
			return nil, err
		}
		arr, err := intrinsicArgArray(name, vals[0], 1)
		if err != nil {
			return nil, err
		}
		return float64(len(arr)), nil

	case "ArrayUnique":
		if err := intrinsicArity(name, len(vals), 1, 1); err != nil {
			return nil, err
		}
		arr, err := intrinsicArgArray(name, vals[0], 1)
		if err != nil {
			return nil, err
		}
		out := make([]interface{}, 0, len(arr))
		for _, item := range arr {
			duplicate := false
			for _, kept := range out {
				if jsonValuesEqual(item, kept) {
					duplicate = true
					break
				}
			}
			if !duplicate {
				out = append(out, item)
			}
		}
		return out, nil

	case "Base64Encode":
		if err := intrinsicArity(name, len(vals), 1, 1); err != nil {
			return nil, err
		}
		data, err := intrinsicArgString(name, vals[0], 1)
		if err != nil {
			return nil, err
		}
		if err := intrinsicStringWithinLimit(name, data); err != nil {
			return nil, err
		}
		return base64.StdEncoding.EncodeToString([]byte(data)), nil

	case "Base64Decode":
		if err := intrinsicArity(name, len(vals), 1, 1); err != nil {
			return nil, err
		}
		data, err := intrinsicArgString(name, vals[0], 1)
		if err != nil {
			return nil, err
		}
		if err := intrinsicStringWithinLimit(name, data); err != nil {
			return nil, err
		}
		decoded, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			return nil, fmt.Errorf("States.Base64Decode argument is not a valid Base64 string")
		}
		return string(decoded), nil

	case "Hash":
		if err := intrinsicArity(name, len(vals), 2, 2); err != nil {
			return nil, err
		}
		data, err := intrinsicArgString(name, vals[0], 1)
		if err != nil {
			return nil, err
		}
		if err := intrinsicStringWithinLimit(name, data); err != nil {
			return nil, err
		}
		algorithm, err := intrinsicArgString(name, vals[1], 2)
		if err != nil {
			return nil, err
		}
		var sum []byte
		switch algorithm {
		case "MD5":
			s := md5.Sum([]byte(data))
			sum = s[:]
		case "SHA-1":
			s := sha1.Sum([]byte(data))
			sum = s[:]
		case "SHA-256":
			s := sha256.Sum256([]byte(data))
			sum = s[:]
		case "SHA-384":
			s := sha512.Sum384([]byte(data))
			sum = s[:]
		case "SHA-512":
			s := sha512.Sum512([]byte(data))
			sum = s[:]
		default:
			return nil, fmt.Errorf("States.Hash algorithm must be one of MD5, SHA-1, SHA-256, SHA-384 or SHA-512")
		}
		return hex.EncodeToString(sum), nil

	case "JsonMerge":
		if err := intrinsicArity(name, len(vals), 3, 3); err != nil {
			return nil, err
		}
		base, err := intrinsicArgObject(name, vals[0], 1)
		if err != nil {
			return nil, err
		}
		overlay, err := intrinsicArgObject(name, vals[1], 2)
		if err != nil {
			return nil, err
		}
		deep, ok := vals[2].(bool)
		if !ok || deep {
			return nil, fmt.Errorf("States.JsonMerge third argument must be the boolean false; only shallow merging is supported")
		}
		merged := make(map[string]interface{}, len(base)+len(overlay))
		for k, v := range base {
			merged[k] = v
		}
		for k, v := range overlay {
			merged[k] = v
		}
		return merged, nil

	case "StringToJson":
		if err := intrinsicArity(name, len(vals), 1, 1); err != nil {
			return nil, err
		}
		data, err := intrinsicArgString(name, vals[0], 1)
		if err != nil {
			return nil, err
		}
		var parsed interface{}
		if err := json.Unmarshal([]byte(data), &parsed); err != nil {
			return nil, fmt.Errorf("States.StringToJson argument is not valid JSON")
		}
		return parsed, nil

	case "JsonToString":
		if err := intrinsicArity(name, len(vals), 1, 1); err != nil {
			return nil, err
		}
		encoded, err := json.Marshal(vals[0])
		if err != nil {
			return nil, fmt.Errorf("States.JsonToString argument could not be serialised")
		}
		return string(encoded), nil

	case "MathRandom":
		if err := intrinsicArity(name, len(vals), 2, 3); err != nil {
			return nil, err
		}
		start, err := intrinsicRoundedInt(name, vals[0], 1)
		if err != nil {
			return nil, err
		}
		end, err := intrinsicRoundedInt(name, vals[1], 2)
		if err != nil {
			return nil, err
		}
		if end <= start {
			return nil, fmt.Errorf("States.MathRandom end must be greater than start")
		}
		var draw int64
		if len(vals) == 3 {
			seed, err := intrinsicRoundedInt(name, vals[2], 3)
			if err != nil {
				return nil, err
			}
			draw = rand.New(rand.NewSource(seed)).Int63n(end - start)
		} else {
			draw = rand.Int63n(end - start)
		}
		return float64(start + draw), nil

	case "MathAdd":
		if err := intrinsicArity(name, len(vals), 2, 2); err != nil {
			return nil, err
		}
		augend, err := intrinsicRoundedInt(name, vals[0], 1)
		if err != nil {
			return nil, err
		}
		addend, err := intrinsicRoundedInt(name, vals[1], 2)
		if err != nil {
			return nil, err
		}
		sum := augend + addend
		if sum < math.MinInt32 || sum > math.MaxInt32 {
			return nil, fmt.Errorf("States.MathAdd result must be in the range -2147483648 to 2147483647")
		}
		return float64(sum), nil

	case "StringSplit":
		if err := intrinsicArity(name, len(vals), 2, 2); err != nil {
			return nil, err
		}
		data, err := intrinsicArgString(name, vals[0], 1)
		if err != nil {
			return nil, err
		}
		delimiters, err := intrinsicArgString(name, vals[1], 2)
		if err != nil {
			return nil, err
		}
		if delimiters == "" {
			return []interface{}{data}, nil
		}
		// Every character of the delimiter argument is a delimiting
		// character; adjacent delimiters delimit empty segments (split
		// semantics rather than field extraction).
		out := []interface{}{}
		var current strings.Builder
		for _, r := range data {
			if strings.ContainsRune(delimiters, r) {
				out = append(out, current.String())
				current.Reset()
			} else {
				current.WriteRune(r)
			}
		}
		out = append(out, current.String())
		return out, nil

	case "UUID":
		if err := intrinsicArity(name, len(vals), 0, 0); err != nil {
			return nil, err
		}
		return uuid.New().String(), nil

	case "Format":
		if err := intrinsicArity(name, len(vals), 1, -1); err != nil {
			return nil, err
		}
		var tmpl formatTemplate
		switch v := vals[0].(type) {
		case formatTemplate:
			tmpl = v
		case string:
			tmpl = markFormatTemplate(v)
		default:
			return nil, fmt.Errorf("States.Format first argument must be a string")
		}
		if tmpl.placeholders != len(vals)-1 {
			return nil, fmt.Errorf("States.Format has %d {} placeholders but %d substitution arguments", tmpl.placeholders, len(vals)-1)
		}
		parts := strings.Split(tmpl.marked, string(formatPlaceholder))
		var b strings.Builder
		b.WriteString(parts[0])
		for i := 1; i < len(parts); i++ {
			b.WriteString(formatIntrinsicArgument(vals[i]))
			b.WriteString(parts[i])
		}
		return b.String(), nil

	default:
		return nil, fmt.Errorf("States.%s is not a supported intrinsic function", name)
	}
}

// formatIntrinsicArgument renders one States.Format substitution argument:
// strings insert verbatim, every other value inserts its JSON serialisation.
func formatIntrinsicArgument(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	encoded, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(encoded)
}

// intrinsicArity checks the argument count; max < 0 means unbounded.
func intrinsicArity(name string, got, min, max int) error {
	if got < min || (max >= 0 && got > max) {
		if max < 0 {
			return fmt.Errorf("States.%s takes at least %d argument(s), got %d", name, min, got)
		}
		if min == max {
			return fmt.Errorf("States.%s takes exactly %d argument(s), got %d", name, min, got)
		}
		return fmt.Errorf("States.%s takes between %d and %d arguments, got %d", name, min, max, got)
	}
	return nil
}

func intrinsicArgString(name string, v interface{}, position int) (string, error) {
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("States.%s argument %d must be a string", name, position)
	}
	return s, nil
}

func intrinsicArgArray(name string, v interface{}, position int) ([]interface{}, error) {
	arr, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf("States.%s argument %d must be an array", name, position)
	}
	return arr, nil
}

func intrinsicArgObject(name string, v interface{}, position int) (map[string]interface{}, error) {
	obj, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("States.%s argument %d must be an object", name, position)
	}
	return obj, nil
}

// intrinsicRoundedInt coerces a numeric argument to an integer, rounding
// non-integer values to the nearest integer as the documentation specifies.
func intrinsicRoundedInt(name string, v interface{}, position int) (int64, error) {
	f, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("States.%s argument %d must be a number", name, position)
	}
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return 0, fmt.Errorf("States.%s argument %d is not a finite number", name, position)
	}
	r := math.Round(f)
	if r > math.MaxInt64 || r < math.MinInt64 {
		return 0, fmt.Errorf("States.%s argument %d is out of integer range", name, position)
	}
	return int64(r), nil
}

// intrinsicStringWithinLimit enforces the documented 10,000 character bound
// on the data strings of States.Base64Encode, States.Base64Decode and
// States.Hash.
func intrinsicStringWithinLimit(name string, data string) error {
	if utf8.RuneCountInString(data) > sfnstore.MaxIntrinsicStringChars {
		return fmt.Errorf("States.%s data string exceeds %d characters", name, sfnstore.MaxIntrinsicStringChars)
	}
	return nil
}

// intrinsicArrayWithinPayloadLimit enforces the documented payload size limit
// on array arguments to States.ArrayPartition and States.ArrayContains.
func intrinsicArrayWithinPayloadLimit(name string, arr []interface{}) error {
	encoded, err := json.Marshal(arr)
	if err != nil {
		return fmt.Errorf("States.%s input array could not be measured against the payload limit", name)
	}
	if len(encoded) > sfnstore.MaxExecutionDataBytes {
		return fmt.Errorf("States.%s input array exceeds the %d byte payload limit", name, sfnstore.MaxExecutionDataBytes)
	}
	return nil
}

// jsonValuesEqual compares two JSON values structurally; numbers arrive as
// float64 from every decode path this package uses.
func jsonValuesEqual(a, b interface{}) bool {
	if an, ok := a.(float64); ok {
		bn, ok := b.(float64)
		return ok && an == bn
	}
	return jsonDeepEqual(a, b)
}

func jsonDeepEqual(a, b interface{}) bool {
	switch av := a.(type) {
	case map[string]interface{}:
		bv, ok := b.(map[string]interface{})
		if !ok || len(av) != len(bv) {
			return false
		}
		for k, v := range av {
			other, exists := bv[k]
			if !exists || !jsonDeepEqual(v, other) {
				return false
			}
		}
		return true
	case []interface{}:
		bv, ok := b.([]interface{})
		if !ok || len(av) != len(bv) {
			return false
		}
		for i, v := range av {
			if !jsonDeepEqual(v, bv[i]) {
				return false
			}
		}
		return true
	case string:
		bs, ok := b.(string)
		return ok && av == bs
	case bool:
		bb, ok := b.(bool)
		return ok && av == bb
	case nil:
		return b == nil
	default:
		return false
	}
}

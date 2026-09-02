package inspection

import (
	"bytes"
	"encoding/json"
	"strings"

	wafstore "vorpalstacks/internal/store/aws/waf"
)

// fieldCandidate is one extracted request-component value. Body-based
// candidates carry the oversize flag so that match statements can
// honour the configured OversizeHandling when the component exceeded
// the inspection limit.
type fieldCandidate struct {
	value     []byte
	oversized bool
}

// extractFieldToMatch returns the candidate values for the configured
// request component. Statements match when any transformed candidate
// satisfies the statement's condition.
func extractFieldToMatch(req *Request, ftm *wafstore.FieldToMatch) []fieldCandidate {
	if ftm == nil {
		return nil
	}
	switch {
	case ftm.Method != nil:
		return singleCandidate(req.Method)
	case ftm.UriPath != nil:
		return singleCandidate(req.URIPath)
	case ftm.QueryString != nil:
		return singleCandidate(req.RawQuery)
	case ftm.AllQueryArguments != nil:
		return queryArgumentCandidates(req.RawQuery, "")
	case ftm.SingleQueryArgument != nil:
		return queryArgumentCandidates(req.RawQuery, ftm.SingleQueryArgument.Name)
	case ftm.SingleHeader != nil:
		var out []fieldCandidate
		for _, v := range req.headerValues(ftm.SingleHeader.Name) {
			out = append(out, fieldCandidate{value: []byte(v)})
		}
		return out
	case ftm.HeaderOrder != nil:
		// The Developer Guide specifies colon separation with no added
		// spaces, for example host:user-agent:accept.
		var b strings.Builder
		for i, h := range req.Headers {
			if i > 0 {
				b.WriteByte(':')
			}
			b.WriteString(strings.ToLower(h.Name))
		}
		return []fieldCandidate{{value: []byte(b.String())}}
	case ftm.Headers != nil:
		return headersCandidates(req, ftm.Headers)
	case ftm.Cookies != nil:
		return cookiesCandidates(req, ftm.Cookies)
	case ftm.Body != nil:
		return bodyCandidates(req, ftm.Body)
	case ftm.JsonBody != nil:
		return jsonBodyCandidates(req, ftm.JsonBody)
	case ftm.UriFragment != nil:
		// Fragments are stripped by clients and intermediaries before
		// the request reaches the protected resource, so WAF applies
		// the statement's FallbackBehaviour instead of component data.
		if ftm.UriFragment.FallbackBehavior == "MATCH" {
			return matchSentinel()
		}
		return nil
	case ftm.JA3Fingerprint != nil:
		// The platform's HTTP listeners do not expose TLS ClientHello
		// fingerprints, so the configured fallback applies.
		if ftm.JA3Fingerprint.FallbackBehavior == "MATCH" {
			return matchSentinel()
		}
		return nil
	case ftm.JA4Fingerprint != nil:
		if ftm.JA4Fingerprint.FallbackBehavior == "MATCH" {
			return matchSentinel()
		}
		return nil
	}
	return nil
}

// matchSentinel produces a candidate that satisfies every condition,
// implementing MATCH fallbacks.
func matchSentinel() []fieldCandidate {
	return []fieldCandidate{{value: nil, oversized: true}}
}

func singleCandidate(v string) []fieldCandidate {
	return []fieldCandidate{{value: []byte(v)}}
}

func queryArgumentCandidates(rawQuery, name string) []fieldCandidate {
	var out []fieldCandidate
	for _, pair := range strings.Split(rawQuery, "&") {
		if pair == "" {
			continue
		}
		key, value, _ := strings.Cut(pair, "=")
		if name != "" && !stringsEqualFold(key, name) {
			continue
		}
		out = append(out, fieldCandidate{value: []byte(value)})
	}
	return out
}

// headersCandidates implements the Headers component: the MatchPattern
// selects which headers to inspect, MatchScope selects whether the
// whole header (ALL), the name (KEY) or the value (VALUE) is inspected,
// and OversizeHandling applies when the header section exceeds the
// inspection limits (first 8 KB / 200 headers).
func headersCandidates(req *Request, cfg *wafstore.Headers) []fieldCandidate {
	if headerSectionOversized(req) {
		switch cfg.OversizeHandling {
		case "MATCH":
			return matchSentinel()
		case "NO_MATCH":
			return nil
		}
		// CONTINUE: inspect the (already bounded) headers as-is.
	}
	var out []fieldCandidate
	for _, h := range req.Headers {
		if !headerSelected(h.Name, cfg.MatchPattern) {
			continue
		}
		switch cfg.MatchScope {
		case "KEY":
			out = append(out, fieldCandidate{value: []byte(h.Name)})
		case "VALUE":
			out = append(out, fieldCandidate{value: []byte(h.Value)})
		default: // ALL
			// All requires a match in the keys or the values or both,
			// so the name and the value are separate candidates. A
			// joined name:value candidate would additionally let a
			// search string straddle the boundary, which the match
			// definition excludes.
			out = append(out, fieldCandidate{value: []byte(h.Name)})
			out = append(out, fieldCandidate{value: []byte(h.Value)})
		}
	}
	return out
}

// headerSelected implements the Developer Guide's match rule for header
// keys: the string match is not case sensitive and is performed after
// trimming leading and trailing spaces from both the request header
// name and the rule's string.
func headerSelected(name string, pattern wafstore.HeaderMatchPattern) bool {
	if pattern.All != nil {
		return true
	}
	trimmed := trimASCIISpace(name)
	if len(pattern.IncludedHeaders) > 0 {
		for _, included := range pattern.IncludedHeaders {
			if stringsEqualFold(trimmed, trimASCIISpace(included)) {
				return true
			}
		}
		return false
	}
	if len(pattern.ExcludedHeaders) > 0 {
		for _, excluded := range pattern.ExcludedHeaders {
			if stringsEqualFold(trimmed, trimASCIISpace(excluded)) {
				return false
			}
		}
	}
	return true
}

func cookiesCandidates(req *Request, cfg *wafstore.Cookies) []fieldCandidate {
	totalLen := 0
	for _, c := range req.Cookies {
		totalLen += len(c.Name) + len(c.Value)
	}
	if totalLen > wafstore.MaxInspectionCookieBytes || len(req.Cookies) > wafstore.MaxInspectionCookieCount {
		switch cfg.OversizeHandling {
		case "MATCH":
			return matchSentinel()
		case "NO_MATCH":
			return nil
		}
	}
	var out []fieldCandidate
	for _, c := range req.Cookies {
		if !cookieSelected(c.Name, cfg.MatchPattern) {
			continue
		}
		switch cfg.MatchScope {
		case "KEY":
			out = append(out, fieldCandidate{value: []byte(c.Name)})
		case "VALUE":
			out = append(out, fieldCandidate{value: []byte(c.Value)})
		default: // ALL
			// As with headers, All inspects the keys or the values or
			// both, so each is its own candidate.
			out = append(out, fieldCandidate{value: []byte(c.Name)})
			out = append(out, fieldCandidate{value: []byte(c.Value)})
		}
	}
	return out
}

// cookieSelected implements the Developer Guide's match rule for cookie
// keys: the string match for a key is case sensitive and must be exact.
func cookieSelected(name string, pattern wafstore.CookieMatchPattern) bool {
	if pattern.All != nil {
		return true
	}
	if len(pattern.IncludedCookies) > 0 {
		for _, included := range pattern.IncludedCookies {
			if name == included {
				return true
			}
		}
		return false
	}
	if len(pattern.ExcludedCookies) > 0 {
		for _, excluded := range pattern.ExcludedCookies {
			if name == excluded {
				return false
			}
		}
	}
	return true
}

func bodyCandidates(req *Request, cfg *wafstore.Body) []fieldCandidate {
	if req.BodyTruncated {
		switch cfg.OversizeHandling {
		case "MATCH":
			return matchSentinel()
		case "NO_MATCH":
			return nil
		}
		// CONTINUE: inspect the bounded prefix.
	}
	return []fieldCandidate{{value: req.Body, oversized: req.BodyTruncated}}
}

// jsonBodyCandidates implements the JSON body component. The body is
// parsed as JSON and the IncludedPaths (RFC 6901 JSON Pointer: slash-
// separated reference tokens with numeric array indices and ~0/~1
// escaping) select the nodes to inspect. A selected container node
// contributes the keys and values of its subtree, excluding the key the
// path itself addressed; a selected scalar node contributes its value.
// MatchScope chooses the key, the value, or both. A body that is not
// valid JSON follows InvalidFallbackBehavior; the default (omitted)
// behaviour evaluates the content only up to the point where parsing
// encountered the error, per the Developer Guide.
func jsonBodyCandidates(req *Request, cfg *wafstore.JsonBody) []fieldCandidate {
	if req.BodyTruncated {
		switch cfg.OversizeHandling {
		case "MATCH":
			return matchSentinel()
		case "NO_MATCH":
			return nil
		}
		// CONTINUE: parse the bounded prefix; a body truncated mid-JSON
		// follows InvalidFallbackBehaviour.
	}
	var root interface{}
	if err := json.Unmarshal(req.Body, &root); err != nil || !jsonRootAllowed(root) {
		switch cfg.InvalidFallbackBehavior {
		case "MATCH":
			return matchSentinel()
		case "EVALUATE_AS_STRING":
			return []fieldCandidate{{value: req.Body}}
		case "NO_MATCH":
			return nil
		default:
			// None (default behaviour): evaluate the valid prefix.
			root = partialJSONDocument(req.Body)
			if !jsonRootAllowed(root) {
				return nil
			}
		}
	}
	var out []fieldCandidate
	for _, node := range selectJSONNodes(root, cfg.MatchPattern) {
		appendJSONElementCandidates(&out, node, cfg.MatchScope)
	}
	return out
}

// jsonRootAllowed reports whether a parsed document root is usable for
// element extraction. The Developer Guide lists content whose root
// node isn't an object or an array among the JSON parsing error
// states, so such content follows InvalidFallbackBehavior.
func jsonRootAllowed(root interface{}) bool {
	switch root.(type) {
	case map[string]interface{}, []interface{}:
		return true
	}
	return false
}

// partialJSONDocument best-effort parses the valid prefix of a JSON
// document: complete elements before the first parsing error are kept
// and everything from the error position on is dropped. A nil return
// means nothing usable was parsed.
func partialJSONDocument(body []byte) interface{} {
	dec := json.NewDecoder(bytes.NewReader(body))
	value, _ := decodePartialValue(dec)
	return value
}

// decodePartialValue decodes one value from the decoder, keeping the
// complete members of an incomplete container. The boolean reports
// whether the value parsed completely.
func decodePartialValue(dec *json.Decoder) (interface{}, bool) {
	tok, err := dec.Token()
	if err != nil {
		return nil, false
	}
	delim, isDelim := tok.(json.Delim)
	if !isDelim {
		return tok, true
	}
	switch delim {
	case '{':
		obj := make(map[string]interface{})
		for dec.More() {
			keyTok, err := dec.Token()
			if err != nil {
				return nonEmptyJSON(obj), false
			}
			key, ok := keyTok.(string)
			if !ok {
				return nonEmptyJSON(obj), false
			}
			val, complete := decodePartialValue(dec)
			if val != nil {
				obj[key] = val
			}
			if !complete {
				return nonEmptyJSON(obj), false
			}
		}
		if _, err := dec.Token(); err != nil {
			return nonEmptyJSON(obj), false
		}
		return obj, true
	case '[':
		arr := make([]interface{}, 0, 4)
		for dec.More() {
			val, complete := decodePartialValue(dec)
			if val != nil {
				arr = append(arr, val)
			}
			if !complete {
				return nonEmptyJSON(arr), false
			}
		}
		if _, err := dec.Token(); err != nil {
			return nonEmptyJSON(arr), false
		}
		return arr, true
	}
	return nil, false
}

// nonEmptyJSON returns the container when it holds at least one
// element, and nil otherwise, so a container with no complete members
// contributes no candidates.
func nonEmptyJSON(container interface{}) interface{} {
	switch typed := container.(type) {
	case map[string]interface{}:
		if len(typed) > 0 {
			return typed
		}
	case []interface{}:
		if len(typed) > 0 {
			return typed
		}
	}
	return nil
}

// appendJSONElementCandidates adds the inspectable string elements of
// one path-selected node. A container node contributes every key and
// string value in its subtree (the subtree walk starts one level below
// the selected node, so the key the included path itself addressed is
// never evaluated); a string node contributes its own value.
func appendJSONElementCandidates(out *[]fieldCandidate, node interface{}, scope string) {
	if s, ok := node.(string); ok {
		if scope != "KEY" {
			*out = append(*out, fieldCandidate{value: []byte(s)})
		}
		return
	}
	var nodes []jsonNode
	collectAllJSON(node, &nodes)
	for _, n := range nodes {
		switch scope {
		case "KEY":
			if n.key != "" {
				*out = append(*out, fieldCandidate{value: []byte(n.key)})
			}
		case "VALUE":
			if s, ok := n.value.(string); ok {
				*out = append(*out, fieldCandidate{value: []byte(s)})
			}
		default: // ALL
			if n.key != "" {
				*out = append(*out, fieldCandidate{value: []byte(n.key)})
			}
			if s, ok := n.value.(string); ok {
				*out = append(*out, fieldCandidate{value: []byte(s)})
			}
		}
	}
}

type jsonNode struct {
	key   string
	value interface{}
}

// selectJSONNodes returns the nodes the pattern selects: the root for
// All, otherwise the nodes the JSON Pointer paths address.
func selectJSONNodes(root interface{}, pattern wafstore.JsonMatchPattern) []interface{} {
	if pattern.All != nil {
		return []interface{}{root}
	}
	var nodes []interface{}
	for _, path := range pattern.IncludedPaths {
		nodes = append(nodes, walkJSONPath(root, splitJSONPath(path))...)
	}
	return nodes
}

func collectAllJSON(v interface{}, nodes *[]jsonNode) {
	switch typed := v.(type) {
	case map[string]interface{}:
		for k, child := range typed {
			*nodes = append(*nodes, jsonNode{key: k, value: child})
			collectAllJSON(child, nodes)
		}
	case []interface{}:
		for _, child := range typed {
			collectAllJSON(child, nodes)
		}
	}
}

// splitJSONPath splits an RFC 6901 JSON Pointer into its reference
// tokens, unescaping ~1 to / and ~0 to ~. A pointer without the
// leading slash addresses nothing.
func splitJSONPath(path string) []string {
	if !strings.HasPrefix(path, "/") {
		return nil
	}
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return nil
	}
	segments := strings.Split(trimmed, "/")
	for i, seg := range segments {
		seg = strings.ReplaceAll(seg, "~1", "/")
		segments[i] = strings.ReplaceAll(seg, "~0", "~")
	}
	return segments
}

func walkJSONPath(v interface{}, segments []string) []interface{} {
	if len(segments) == 0 {
		return []interface{}{v}
	}
	seg := segments[0]
	var out []interface{}
	switch typed := v.(type) {
	case map[string]interface{}:
		if child, ok := typed[seg]; ok {
			out = append(out, walkJSONPath(child, segments[1:])...)
		}
	case []interface{}:
		if idx, ok := arrayIndex(seg); ok && idx < len(typed) {
			out = append(out, walkJSONPath(typed[idx], segments[1:])...)
		}
	}
	return out
}

func arrayIndex(seg string) (int, bool) {
	if len(seg) == 0 || (len(seg) > 1 && seg[0] == '0') {
		return 0, false
	}
	n := 0
	for i := 0; i < len(seg); i++ {
		if seg[i] < '0' || seg[i] > '9' {
			return 0, false
		}
		n = n*10 + int(seg[i]-'0')
	}
	return n, true
}

// headerSectionOversized reports whether the request's header section
// exceeds the AWS WAF inspection limits for headers: at most the first
// 8 KB of the header section and at most the first 200 headers.
func headerSectionOversized(req *Request) bool {
	if len(req.Headers) > wafstore.MaxInspectionHeaderCount {
		return true
	}
	total := 0
	for _, h := range req.Headers {
		total += len(h.Name) + len(h.Value) + 2
	}
	return total > wafstore.MaxInspectionHeaderBytes
}

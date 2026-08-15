package iot

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	storecommon "vorpalstacks/internal/store/aws/common"
	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// Fleet Indexing Query Engine
//
// Implements the AWS IoT fleet indexing query syntax for SearchIndex and
// related aggregation operations. Supports:
//   - thingName:<pattern>    (with * wildcard)
//   - thingType:<value>
//   - attributes.<name>:<value>
//   - attributes.<name>:>|>=|<|<=|<> <number>
//   - isConnected:true|false
//   - Boolean operators AND, OR, NOT
//   - Parenthesised grouping
//   - "*" to match all things
//
// Reference: https://docs.aws.amazon.com/iot/latest/developerguide/query-syntax.html
// ---------------------------------------------------------------------------

// queryNode is a node in the parsed query AST.
type queryNode interface {
	match(thing *iotstore.Thing, conn connectedSet) bool
}

// connectedSet reports whether a thing name is currently connected.
type connectedSet map[string]bool

// matchAllNode matches every thing.
type matchAllNode struct{}

func (matchAllNode) match(_ *iotstore.Thing, _ connectedSet) bool { return true }

// fieldNode matches a single field comparison.
type fieldNode struct {
	field   string
	op      string // "", ">", ">=", "<", "<=", "<>"
	value   string
	numeric float64
	isNum   bool
}

func (n *fieldNode) match(thing *iotstore.Thing, conn connectedSet) bool {
	var actual string
	switch {
	case strings.HasPrefix(n.field, "attributes."):
		attrName := strings.TrimPrefix(n.field, "attributes.")
		actual = thing.Attributes[attrName]
		if actual == "" {
			return false
		}
	case n.field == "thingName":
		actual = thing.ThingName
	case n.field == "thingType":
		actual = thing.ThingTypeName
	case n.field == "isConnected":
		connected := conn[thing.ThingName]
		switch strings.ToLower(n.value) {
		case "true":
			return connected
		case "false":
			return !connected
		default:
			return false
		}
	default:
		return false
	}

	switch n.op {
	case "":
		// String match with wildcard support.
		if strings.Contains(n.value, "*") {
			return wildcardMatchString(actual, n.value)
		}
		return actual == n.value
	case ">":
		return numericCompare(actual) > n.numeric
	case ">=":
		return numericCompare(actual) >= n.numeric
	case "<":
		return numericCompare(actual) < n.numeric
	case "<=":
		return numericCompare(actual) <= n.numeric
	case "<>":
		if n.isNum {
			return numericCompare(actual) != n.numeric
		}
		return actual != n.value
	}
	return false
}

func numericCompare(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return math.NaN()
	}
	return v
}

func wildcardMatchString(s, pattern string) bool {
	// Simple wildcard: * matches any characters.
	parts := strings.Split(pattern, "*")
	if len(parts) == 1 {
		return s == pattern
	}
	idx := 0
	for i, p := range parts {
		if p == "" {
			continue
		}
		if i == 0 {
			if !strings.HasPrefix(s[idx:], p) {
				return false
			}
		} else if i == len(parts)-1 {
			if !strings.HasSuffix(s, p) {
				return false
			}
		}
		pos := strings.Index(s[idx:], p)
		if pos < 0 {
			return false
		}
		idx += pos + len(p)
	}
	return true
}

// andNode matches when all children match.
type andNode struct{ children []queryNode }

func (n *andNode) match(t *iotstore.Thing, conn connectedSet) bool {
	for _, c := range n.children {
		if !c.match(t, conn) {
			return false
		}
	}
	return true
}

// orNode matches when any child matches.
type orNode struct{ children []queryNode }

func (n *orNode) match(t *iotstore.Thing, conn connectedSet) bool {
	for _, c := range n.children {
		if c.match(t, conn) {
			return true
		}
	}
	return false
}

// notNode matches when the child does not match.
type notNode struct{ child queryNode }

func (n *notNode) match(t *iotstore.Thing, conn connectedSet) bool { return !n.child.match(t, conn) }

// parseQuery parses an AWS IoT fleet indexing query string into a queryNode.
func parseQuery(query string) (queryNode, error) {
	query = strings.TrimSpace(query)
	if query == "" || query == "*" {
		return matchAllNode{}, nil
	}
	p := &queryParser{tokens: tokenizeQuery(query)}
	return p.parseOr()
}

type queryParser struct {
	tokens []string
	pos    int
}

func (p *queryParser) peek() string {
	if p.pos >= len(p.tokens) {
		return ""
	}
	return p.tokens[p.pos]
}

func (p *queryParser) next() string {
	if p.pos >= len(p.tokens) {
		return ""
	}
	t := p.tokens[p.pos]
	p.pos++
	return t
}

func (p *queryParser) parseOr() (queryNode, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	children := []queryNode{left}
	for {
		upper := strings.ToUpper(p.peek())
		if upper != "OR" {
			break
		}
		p.next()
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		children = append(children, right)
	}
	if len(children) == 1 {
		return children[0], nil
	}
	return &orNode{children: children}, nil
}

func (p *queryParser) parseAnd() (queryNode, error) {
	left, err := p.parseNot()
	if err != nil {
		return nil, err
	}
	children := []queryNode{left}
	for {
		upper := strings.ToUpper(p.peek())
		if upper != "AND" {
			break
		}
		p.next()
		right, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		children = append(children, right)
	}
	if len(children) == 1 {
		return children[0], nil
	}
	return &andNode{children: children}, nil
}

func (p *queryParser) parseNot() (queryNode, error) {
	if strings.ToUpper(p.peek()) == "NOT" {
		p.next()
		child, err := p.parseAtom()
		if err != nil {
			return nil, err
		}
		return &notNode{child: child}, nil
	}
	return p.parseAtom()
}

func (p *queryParser) parseAtom() (queryNode, error) {
	tok := p.peek()
	if tok == "(" {
		p.next()
		node, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		if p.peek() != ")" {
			return nil, fmt.Errorf("expected closing parenthesis")
		}
		p.next()
		return node, nil
	}
	// Parse field:value or field:op value
	return p.parseTerm()
}

func (p *queryParser) parseTerm() (queryNode, error) {
	term := p.next()
	if term == "" {
		return nil, fmt.Errorf("unexpected end of query")
	}
	// Split on the first colon.
	colonIdx := strings.Index(term, ":")
	if colonIdx < 0 {
		// Bare word: treat as thingName match.
		return &fieldNode{field: "thingName", value: term}, nil
	}
	field := term[:colonIdx]
	rest := term[colonIdx+1:]

	// Check for comparison operators.
	for _, op := range []string{">=", "<=", "<>", ">", "<"} {
		if strings.HasPrefix(rest, op) {
			val := strings.TrimSpace(strings.TrimPrefix(rest, op))
			if val == "" {
				// A lone operator (e.g. "field:>" when the tokenizer
				// split "field:> 20" on the space) must not silently
				// match nothing.
				return nil, fmt.Errorf("missing comparison value after %q in %q", op, term)
			}
			node := &fieldNode{field: field, op: op, value: val}
			if num, err := strconv.ParseFloat(val, 64); err == nil {
				node.numeric = num
				node.isNum = true
			}
			return node, nil
		}
	}
	return &fieldNode{field: field, value: rest}, nil
}

// tokenizeQuery splits a query string into tokens. Parentheses are separate
// tokens; everything else is split on whitespace, except within a field:value
// term where the colon binds the two sides.
func tokenizeQuery(query string) []string {
	var tokens []string
	i := 0
	for i < len(query) {
		ch := query[i]
		switch ch {
		case ' ', '\t', '\n':
			i++
		case '(', ')':
			tokens = append(tokens, string(ch))
			i++
		default:
			// Read until whitespace or parenthesis.
			start := i
			for i < len(query) && query[i] != ' ' && query[i] != '\t' && query[i] != '\n' && query[i] != '(' && query[i] != ')' {
				i++
			}
			tokens = append(tokens, query[start:i])
		}
	}
	return tokens
}

// ---------------------------------------------------------------------------
// Fleet indexing search helpers
// ---------------------------------------------------------------------------

// searchThings loads all things from the store and filters them by the
// parsed query. Returns the full matching slice; pagination is applied by
// the caller.
func (s *IoTService) searchThings(reqCtx *request.RequestContext, queryString string) ([]*iotstore.Thing, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	qNode, err := parseQuery(queryString)
	if err != nil {
		return nil, iotstore.ErrInvalidRequest
	}

	// Pre-compute the set of connected thing names from all brokers in
	// the current region.  This is used by the isConnected query filter.
	conn := s.buildConnectedSet(reqCtx, store)

	// Paginate through ALL things, applying the query filter.
	var matched []*iotstore.Thing
	var opts storecommon.ListOptions
	for {
		result, err := store.ListThings(opts, "", "")
		if err != nil {
			return nil, err
		}
		for i := range result.Items {
			if qNode.match(result.Items[i], conn) {
				matched = append(matched, result.Items[i])
			}
		}
		if result.NextMarker == "" {
			break
		}
		opts.Marker = result.NextMarker
	}
	return matched, nil
}

// buildConnectedSet returns a set of thing names that have at least one
// principal certificate currently connected to any broker. The account and
// region are required to construct the full certificate ARN that
// ListThingsForPrincipal uses as a PebbleDB key prefix.
func (s *IoTService) buildConnectedSet(reqCtx *request.RequestContext, store iotstore.IotStoreInterface) connectedSet {
	conn := connectedSet{}
	accountID := ""
	region := ""
	if reqCtx != nil {
		accountID = reqCtx.GetAccountID()
		region = reqCtx.GetRegion()
	}
	for _, b := range s.brokers {
		for _, certID := range b.ConnectedCertIDs() {
			// ListThingsForPrincipal keys on the full certificate ARN,
			// not the raw SHA-256 hash.
			principal := iotstore.BuildCertificateARN(accountID, region, certID)
			thingNames, err := store.ListThingsForPrincipal(principal)
			if err != nil {
				continue
			}
			for _, tn := range thingNames {
				conn[tn] = true
			}
		}
	}
	return conn
}

// ---------------------------------------------------------------------------
// SearchIndex handler — replaces the previous stub.
// ---------------------------------------------------------------------------

func (s *IoTService) searchIndexImpl(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryString := request.GetParamCaseInsensitive(req.Parameters, "queryString")
	if queryString == "" {
		return nil, iotstore.ErrMissingParam
	}

	matched, err := s.searchThings(reqCtx, queryString)
	if err != nil {
		return nil, err
	}

	// Apply pagination.
	start, end, err := applyPagination(len(matched), req.Parameters)
	if err != nil {
		return nil, err
	}
	page := matched[start:min(end, len(matched))]

	things := make([]map[string]interface{}, 0, len(page))
	for _, t := range page {
		id := t.ThingID
		if id == "" {
			id = t.ThingName
		}
		things = append(things, map[string]interface{}{
			"thingName":     t.ThingName,
			"thingId":       id,
			"thingTypeName": t.ThingTypeName,
			"attributes":    t.Attributes,
			"version":       t.Version,
		})
	}

	resp := map[string]interface{}{
		"things":    things,
		"nextToken": nextPageToken(len(matched), end, req.Parameters),
	}
	return resp, nil
}

// ---------------------------------------------------------------------------
// GetStatistics handler
// ---------------------------------------------------------------------------

func (s *IoTService) getStatisticsImpl(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryString := request.GetParamCaseInsensitive(req.Parameters, "queryString")
	if queryString == "" {
		return nil, iotstore.ErrMissingParam
	}

	matched, err := s.searchThings(reqCtx, queryString)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"count": int64(len(matched)),
	}

	// If aggregationField is specified, compute numeric statistics.
	if aggField := request.GetParamCaseInsensitive(req.Parameters, "aggregationField"); aggField != "" {
		var values []float64
		for _, t := range matched {
			if v := getNumericAttribute(t, aggField); !math.IsNaN(v) {
				values = append(values, v)
			}
		}
		if len(values) > 0 {
			sum := 0.0
			minV := values[0]
			maxV := values[0]
			for _, v := range values {
				sum += v
				if v < minV {
					minV = v
				}
				if v > maxV {
					maxV = v
				}
			}
			stats["sum"] = sum
			stats["minimum"] = minV
			stats["maximum"] = maxV
			stats["average"] = sum / float64(len(values))
			stats["valueCount"] = int64(len(values))
		}
	}

	return map[string]interface{}{"statistics": stats}, nil
}

// ---------------------------------------------------------------------------
// GetCardinality handler
// ---------------------------------------------------------------------------

func (s *IoTService) getCardinalityImpl(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryString := request.GetParamCaseInsensitive(req.Parameters, "queryString")
	if queryString == "" {
		return nil, iotstore.ErrMissingParam
	}
	aggField := request.GetParamCaseInsensitive(req.Parameters, "aggregationField")

	matched, err := s.searchThings(reqCtx, queryString)
	if err != nil {
		return nil, err
	}

	// Without an aggregationField, AWS GetCardinality returns the total
	// number of things matching the query.
	if aggField == "" {
		return map[string]interface{}{"cardinality": int64(len(matched))}, nil
	}

	seen := make(map[string]bool)
	for _, t := range matched {
		val := getStringAttribute(t, aggField)
		if val != "" {
			seen[val] = true
		}
	}

	return map[string]interface{}{"cardinality": int64(len(seen))}, nil
}

// ---------------------------------------------------------------------------
// GetPercentiles handler
// ---------------------------------------------------------------------------

func (s *IoTService) getPercentilesImpl(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryString := request.GetParamCaseInsensitive(req.Parameters, "queryString")
	if queryString == "" {
		return nil, iotstore.ErrMissingParam
	}
	aggField := request.GetParamCaseInsensitive(req.Parameters, "aggregationField")

	// Without an aggregationField, return empty percentiles.
	if aggField == "" {
		return map[string]interface{}{"percentiles": []map[string]interface{}{}}, nil
	}

	matched, err := s.searchThings(reqCtx, queryString)
	if err != nil {
		return nil, err
	}

	var values []float64
	for _, t := range matched {
		if v := getNumericAttribute(t, aggField); !math.IsNaN(v) {
			values = append(values, v)
		}
	}

	pct := request.GetIntParam(req.Parameters, "percent")

	var percentiles []map[string]interface{}
	if len(values) > 0 {
		sort.Float64s(values)
		if pct > 0 && pct <= 100 {
			percentiles = append(percentiles, map[string]interface{}{
				"percent": int64(pct),
				"value":   percentileValue(values, float64(pct)),
			})
		} else {
			for _, p := range []int{50, 90, 95, 99} {
				percentiles = append(percentiles, map[string]interface{}{
					"percent": int64(p),
					"value":   percentileValue(values, float64(p)),
				})
			}
		}
	}

	return map[string]interface{}{"percentiles": percentiles}, nil
}

// ---------------------------------------------------------------------------
// GetBucketsAggregation handler
// ---------------------------------------------------------------------------

func (s *IoTService) getBucketsAggregationImpl(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queryString := request.GetParamCaseInsensitive(req.Parameters, "queryString")
	if queryString == "" {
		return nil, iotstore.ErrMissingParam
	}
	aggField := request.GetParamCaseInsensitive(req.Parameters, "aggregationField")
	if aggField == "" {
		return nil, iotstore.ErrMissingParam
	}

	matched, err := s.searchThings(reqCtx, queryString)
	if err != nil {
		return nil, err
	}

	bucketCounts := make(map[string]int64)
	for _, t := range matched {
		val := getStringAttribute(t, aggField)
		bucketCounts[val]++
	}

	buckets := make([]map[string]interface{}, 0, len(bucketCounts))
	for val, count := range bucketCounts {
		buckets = append(buckets, map[string]interface{}{
			"keyValue": val,
			"count":    count,
		})
	}

	return map[string]interface{}{
		"totalCount":  int64(len(matched)),
		"bucketCount": int64(len(buckets)),
		"buckets":     buckets,
	}, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func getNumericAttribute(t *iotstore.Thing, field string) float64 {
	val := getStringAttribute(t, field)
	if val == "" {
		return math.NaN()
	}
	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return math.NaN()
	}
	return v
}

func getStringAttribute(t *iotstore.Thing, field string) string {
	switch {
	case strings.HasPrefix(field, "attributes."):
		return t.Attributes[strings.TrimPrefix(field, "attributes.")]
	case field == "thingName":
		return t.ThingName
	case field == "thingType":
		return t.ThingTypeName
	}
	return ""
}

func percentileValue(sorted []float64, pct float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	idx := (pct / 100) * float64(len(sorted)-1)
	lower := int(math.Floor(idx))
	upper := int(math.Ceil(idx))
	if lower == upper || upper >= len(sorted) {
		return sorted[lower]
	}
	frac := idx - float64(lower)
	return sorted[lower]*(1-frac) + sorted[upper]*frac
}

// defaultSearchIndexMaxResults is the page size applied when the request
// omits maxResults. The AWS API documentation does not state a default, so
// the local convention of fifty applies; an explicitly supplied maxResults
// is always validated against the documented range instead.
const defaultSearchIndexMaxResults = 50

// applyPagination slices the matched set for one response page. A supplied
// maxResults must be numeric and within the documented SearchIndex range
// (SearchIndexMaxResultsMin..SearchIndexMaxResultsCap): a non-numeric or
// out-of-range value is a malformed request and is rejected with
// InvalidRequestException — silently paging at the default or accepting an
// unbounded page would return results the caller never asked for.
func applyPagination(total int, params map[string]interface{}) (int, int, error) {
	limit := 0
	if raw, ok := params["maxResults"]; ok && raw != nil {
		switch v := raw.(type) {
		case float64:
			limit = int(v)
		case int:
			limit = v
		case string:
			n, err := strconv.Atoi(v)
			if err != nil {
				return 0, 0, iotstore.ErrInvalidRequest
			}
			limit = n
		default:
			return 0, 0, iotstore.ErrInvalidRequest
		}
		if limit < iotstore.SearchIndexMaxResultsMin || limit > iotstore.SearchIndexMaxResultsCap {
			return 0, 0, iotstore.ErrInvalidRequest
		}
	}
	if limit == 0 {
		limit = defaultSearchIndexMaxResults
	}
	start := 0
	if token, ok := params["nextToken"].(string); ok && token != "" {
		n, err := strconv.Atoi(token)
		if err != nil {
			return 0, 0, iotstore.ErrInvalidRequest
		}
		start = n
	}
	if start >= total {
		return total, total, nil
	}
	end := start + limit
	if end > total {
		end = total
	}
	return start, end, nil
}

func nextPageToken(total, end int, params map[string]interface{}) string {
	if end >= total {
		return ""
	}
	return strconv.Itoa(end)
}

// ---- Fleet Indexing / Metrics --------------------------------------

func (s *IoTService) DescribeIndex(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetParamCaseInsensitive(req.Parameters, "indexName")
	return map[string]interface{}{"indexName": name, "indexStatus": "ACTIVE"}, nil
}
func (s *IoTService) ListIndices(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// Smithy: ListIndicesResponse.indexNames is a list of IndexName (string).
	return map[string]interface{}{"indexNames": []string{"AWS_Things"}}, nil
}
func (s *IoTService) SearchIndex(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.searchIndexImpl(ctx, reqCtx, req)
}
func (s *IoTService) GetBucketsAggregation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.getBucketsAggregationImpl(ctx, reqCtx, req)
}
func (s *IoTService) GetCardinality(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.getCardinalityImpl(ctx, reqCtx, req)
}
func (s *IoTService) GetPercentiles(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.getPercentilesImpl(ctx, reqCtx, req)
}
func (s *IoTService) GetStatistics(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return s.getStatisticsImpl(ctx, reqCtx, req)
}
func (s *IoTService) ListMetricValues(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	metricName := request.GetParamCaseInsensitive(req.Parameters, "metricName")
	if metricName == "" {
		return nil, iotstore.ErrMissingParam
	}
	rec, _, exists, err := s.bulkGet(reqCtx, "fleetMetric", req, "metricName")
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrFleetMetricNotFound
	}

	// Execute the fleet metric's query to compute the current value.
	queryString, _ := rec["queryString"].(string)
	aggField, _ := rec["aggregationField"].(string)
	aggType, _ := rec["aggregationType"].(string)

	matched, err := s.searchThings(reqCtx, queryString)
	if err != nil {
		return nil, err
	}

	var value float64
	var count int64
	if aggField != "" {
		var values []float64
		for _, t := range matched {
			if v := getNumericAttribute(t, aggField); !math.IsNaN(v) {
				values = append(values, v)
			}
		}
		count = int64(len(values))
		if count > 0 {
			switch strings.ToUpper(aggType) {
			case "AVERAGE":
				sum := 0.0
				for _, v := range values {
					sum += v
				}
				value = sum / float64(count)
			case "SUM":
				for _, v := range values {
					value += v
				}
			case "MINIMUM":
				value = values[0]
				for _, v := range values {
					if v < value {
						value = v
					}
				}
			case "MAXIMUM":
				value = values[0]
				for _, v := range values {
					if v > value {
						value = v
					}
				}
			default: // COUNT or unspecified
				value = float64(count)
			}
		}
	} else {
		value = float64(len(matched))
		count = int64(len(matched))
	}

	now := time.Now().UTC().Unix()
	values := []map[string]interface{}{
		{
			"timestamp": now,
			"value":     value,
		},
	}
	return paginatedMaps("metricValues", values, req.Parameters), nil
}

package dynamodb

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// searchVectorsCoreResult carries the SearchVectors response in a
// transport-independent form: the projected items, their similarity scores
// in the same order, and the vector search request bytes consumed (the
// serialised byte length of the request's SearchVector list — AWS documents
// the member but not its formula).
type searchVectorsCoreResult struct {
	Items                    []map[string]interface{}
	Scores                   []float64
	VectorSearchRequestBytes int
}

// searchVectorsCore is the single entry point for SearchVectors: it
// validates the request against the table and vector index, evaluates the
// search-condition expression subset, runs the brute-force top-k search, and
// projects the results to the index projection. A missing table, a missing
// index, or an index that is not ACTIVE is ResourceNotFoundException, per
// the API reference error list.
func (s *DynamoDBService) searchVectorsCore(ctx context.Context, store dbstore.DynamoDBStoreInterface, params map[string]interface{}) (*searchVectorsCoreResult, error) {
	if _, err := getReturnConsumedCapacity(params); err != nil {
		return nil, err
	}
	tableName, err := resolveVectorSearchTable(request.GetStringParam(params, "TableName"))
	if err != nil {
		return nil, err
	}
	indexName := request.GetStringParam(params, "IndexName")
	if indexName == "" || !validateResourceName(indexName) {
		return nil, ErrInvalidParameter
	}

	table, err := store.Tables().Get(tableName)
	if err != nil {
		if dbstore.IsTableNotFound(err) {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}
	if table.Status != dbstore.TableStatusActive {
		return nil, ErrResourceNotFound
	}

	var vi *dbstore.VectorIndex
	for _, cand := range table.VectorIndexes {
		if cand.IndexName == indexName {
			vi = cand
			break
		}
	}
	if vi == nil || vi.IndexStatus != dbstore.IndexStatusActive {
		return nil, ErrResourceNotFound
	}

	rawVec, ok := params["SearchVector"].([]interface{})
	if !ok || len(rawVec) == 0 {
		return nil, ErrInvalidParameter
	}
	queryVec := make([]float64, len(rawVec))
	for i, el := range rawVec {
		av, avErr := parseAttributeValue(el)
		if avErr != nil || av.N == nil {
			return nil, ErrInvalidParameter
		}
		// SearchVector elements are 32-bit IEEE-754 values in DynamoDB list
		// format (SearchVectors API reference): bitSize 32 rounds each parsed
		// element to its float32 value, returned as the exact float64 of that
		// value. The distance arithmetic stays double-precision — the
		// documented contract pins the element width only — and operates on
		// exactly-representable float32 inputs.
		f, perr := strconv.ParseFloat(*av.N, 32)
		if perr != nil {
			return nil, ErrInvalidParameter
		}
		queryVec[i] = f
	}
	if int64(len(queryVec)) != vi.Dimensions {
		return nil, ErrInvalidParameter
	}

	topK := request.GetIntParam(params, "TopK")
	if topK < dbstore.VectorSearchTopKMin || topK > dbstore.VectorSearchTopKMax {
		return nil, ErrInvalidParameter
	}

	names, err := parseExpressionAttributeNames(params)
	if err != nil {
		return nil, err
	}
	values, err := parseExpressionAttributeValues(params)
	if err != nil {
		return nil, err
	}

	condition, err := parseSearchCondition(request.GetStringParam(params, "SearchConditionExpression"), names, values, vi)
	if err != nil {
		return nil, err
	}

	projection, err := parseProjectionExpression(params)
	if err != nil {
		return nil, err
	}
	allowed, err := vectorProjectionSet(vi, table, projection)
	if err != nil {
		return nil, err
	}

	var hits []dbstore.VectorSearchHit
	err = store.View(ctx, func(txn *dbstore.DynamoDBTxn) error {
		var terr error
		hits, terr = txn.VectorTopK(tableName, indexName, queryVec, topK, condition.matches)
		return terr
	})
	if err != nil {
		return nil, err
	}

	items := make([]map[string]interface{}, 0, len(hits))
	scores := make([]float64, 0, len(hits))
	for _, h := range hits {
		merged := make(map[string]*dbstore.AttributeValue, len(h.Item.Key)+len(h.Item.Attributes))
		for k, v := range h.Item.Key {
			merged[k] = v
		}
		for k, v := range h.Item.Attributes {
			merged[k] = v
		}
		if allowed != nil {
			for k := range merged {
				if !allowed[k] {
					delete(merged, k)
				}
			}
		}
		items = append(items, buildItemResponse(merged))
		scores = append(scores, h.Score)
	}

	wireVec, _ := json.Marshal(rawVec)
	return &searchVectorsCoreResult{
		Items:                    items,
		Scores:                   scores,
		VectorSearchRequestBytes: len(wireVec),
	}, nil
}

// resolveVectorSearchTable accepts the SearchVectors TableName member in
// both documented forms — "The name or Amazon Resource Name (ARN) of the
// table containing the vector index" — and returns the table name. It is
// the one DynamoDB table parameter modelled as TableArn rather than
// TableName; a plain table name keeps the shared resource-name validation.
// The ARN form must be the table ARN itself: stream and index ARNs carry
// further resource segments, so a resource that is not exactly table/<name>
// does not resolve (the family's ParseTableARN accepts stream ARNs and
// would widen the member past its documentation).
func resolveVectorSearchTable(member string) (string, error) {
	return resolveTableNameMember(member)
}

// ---------------------------------------------------------------------------
// Search-condition expression
// ---------------------------------------------------------------------------

// searchConditionClause is one AND-joined comparison of the search-condition
// expression: a search-schema attribute against a substituted value. The
// published contract supports the equality operator only — for HASH and
// INLINE_FILTER attributes alike (the vendored model snapshot's
// range-operator wording predates that revision of the documentation).
type searchConditionClause struct {
	attrName string
	value    *dbstore.AttributeValue
}

func (c *searchConditionClause) matches(item *dbstore.Item) bool {
	var av *dbstore.AttributeValue
	if item.Attributes != nil {
		av = item.Attributes[c.attrName]
	}
	if av == nil && item.Key != nil {
		av = item.Key[c.attrName]
	}
	return av != nil && attributeValuesEqual(av, c.value)
}

// searchConditionConjunction is the parsed SearchConditionExpression: every
// clause must match (the grammar is a flat AND of comparisons over top-level
// search-schema attributes).
type searchConditionConjunction struct {
	clauses []*searchConditionClause
}

func (c *searchConditionConjunction) matches(item *dbstore.Item) bool {
	for _, clause := range c.clauses {
		if !clause.matches(item) {
			return false
		}
	}
	return true
}

// hasClause reports whether the conjunction carries a comparison on the
// named attribute — the requirement check for HASH search-schema elements.
func (c *searchConditionConjunction) hasClause(attrName string) bool {
	for _, clause := range c.clauses {
		if clause.attrName == attrName {
			return true
		}
	}
	return false
}

// parseSearchCondition parses the SearchConditionExpression contract: a
// flat AND of equality comparisons whose attribute side must be a top-level
// attribute declared in the index search schema. Equality is the only
// supported operator for HASH and INLINE_FILTER attributes alike. A HASH
// element partitions the index, so its attribute must appear as an equality
// clause — a search without the partition-key value is rejected. With no
// HASH element in the schema, an empty expression matches everything.
// Every token position the walk advances past — the value after "=" and
// the attribute after a trailing AND included — is bounds-checked, so a
// truncated expression is the invalid-parameter error, never a read past
// the token slice or a silently dropped tail.
func parseSearchCondition(expr string, names map[string]string, values map[string]*dbstore.AttributeValue, vi *dbstore.VectorIndex) (*searchConditionConjunction, error) {
	schemaRole := make(map[string]dbstore.SearchSchemaElementType)
	var hashNames []string
	for _, e := range vi.SearchSchema {
		schemaRole[e.AttributeName] = e.SearchSchemaElementType
		if e.SearchSchemaElementType == dbstore.SearchSchemaElementTypeHash {
			hashNames = append(hashNames, e.AttributeName)
		}
	}

	resolveAttr := func(token string) (string, error) {
		name := token
		if resolved, ok := names[token]; ok {
			name = resolved
		}
		if strings.ContainsAny(name, ".[") {
			return "", ErrInvalidParameter // only top-level attributes
		}
		if _, inSchema := schemaRole[name]; !inSchema {
			return "", ErrInvalidParameter
		}
		return name, nil
	}
	resolveValue := func(token string) (*dbstore.AttributeValue, error) {
		v, ok := values[token]
		if !ok {
			return nil, ErrInvalidParameter
		}
		return v, nil
	}

	conj := &searchConditionConjunction{}
	if expr != "" {
		tokens := strings.Fields(expr)
		i := 0
		for i < len(tokens) {
			attrName, err := resolveAttr(tokens[i])
			if err != nil {
				return nil, err
			}
			i++
			if i >= len(tokens) || tokens[i] != "=" {
				return nil, ErrInvalidParameter
			}
			i++
			if i >= len(tokens) {
				return nil, ErrInvalidParameter
			}
			clause := &searchConditionClause{attrName: attrName}
			if clause.value, err = resolveValue(tokens[i]); err != nil {
				return nil, err
			}
			i++
			conj.clauses = append(conj.clauses, clause)

			if i < len(tokens) {
				if tokens[i] != "AND" {
					return nil, ErrInvalidParameter
				}
				i++
				if i >= len(tokens) {
					return nil, ErrInvalidParameter
				}
			}
		}
		if len(conj.clauses) == 0 {
			return nil, ErrInvalidParameter
		}
	}
	for _, hashName := range hashNames {
		if !conj.hasClause(hashName) {
			return nil, ErrInvalidParameter
		}
	}
	return conj, nil
}

// ---------------------------------------------------------------------------
// Projection
// ---------------------------------------------------------------------------

// vectorProjectionSet resolves the retrievable attribute set. The index
// projection defines the universe (ALL = everything; KEYS_ONLY = key
// attributes plus the vector attribute; INCLUDE = keys, the vector
// attribute, the non-key projection list, and the search-schema elements —
// inline filters are projected by definition), and an explicit
// ProjectionExpression narrows it further: requesting an attribute outside
// the projection is rejected, per "Only attributes projected into the
// vector index can be retrieved".
func vectorProjectionSet(vi *dbstore.VectorIndex, table *dbstore.Table, projection [][]docPathPart) (map[string]bool, error) {
	if vi.Projection == nil {
		return nil, fmt.Errorf("vector index %s has no projection", vi.IndexName)
	}

	var projected map[string]bool
	switch vi.Projection.ProjectionType {
	case "", "ALL":
		projected = nil
	default:
		projected = make(map[string]bool)
		for _, ks := range table.KeySchema {
			projected[ks.AttributeName] = true
		}
		projected[vi.VectorAttributeName] = true
		if vi.Projection.ProjectionType == "INCLUDE" {
			for _, a := range vi.Projection.NonKeyAttributes {
				projected[a] = true
			}
		}
		for _, e := range vi.SearchSchema {
			projected[e.AttributeName] = true
		}
	}

	if len(projection) == 0 {
		return projected, nil
	}
	// Key attributes are part of every projection, so an explicit
	// ProjectionExpression narrows the non-key set only. The set is named
	// by top-level attribute — the segment a document path starts from.
	narrowed := make(map[string]bool, len(projection))
	for _, ks := range table.KeySchema {
		narrowed[ks.AttributeName] = true
	}
	for _, path := range projection {
		attr := ""
		if len(path) > 0 {
			attr = path[0].name
		}
		if projected != nil && !projected[attr] {
			return nil, ErrInvalidParameter
		}
		narrowed[attr] = true
	}
	return narrowed, nil
}

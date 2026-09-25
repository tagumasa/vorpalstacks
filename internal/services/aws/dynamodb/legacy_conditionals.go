package dynamodb

import (
	"fmt"
	"net/http"
	"sort"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// The legacy conditional parameter family — Expected, KeyConditions,
// QueryFilter, ScanFilter, ConditionalOperator and AttributesToGet — is
// deprecated in favour of the expression parameters that replaced them,
// but the members remain on the wire contract (developer guide,
// LegacyConditionalParameters). The Expected and filter translators below
// build the typed condition tree directly from already-validated wire
// members — no expression string is generated or reparsed. The
// KeyConditions translator still renders an expression string: its
// consumer is the key-condition engine, which has its own decomposition;
// that string round-trip closes when the key plane adopts the same typed
// core.
//
// A legacy parameter never combines with the expression parameter that
// replaced it, nor with the expression substitution maps: such a request
// carries two contradictory specifications of one condition, and the
// documented precedent inside Expected itself — Value/Exists alongside
// AttributeValueList/ComparisonOperator — is a ValidationException, never
// a silent resolution.

// legacyValidationException is the rejection every legacy translator and
// exclusivity check answers with.
func legacyValidationException(message string) error {
	return NewAPIError("com.amazon.coral.validate#ValidationException", message, http.StatusBadRequest)
}

// legacyJoiner resolves ConditionalOperator: AND by default, OR when asked.
// Any other value violates the ConditionalOperator enum.
func legacyJoiner(params map[string]interface{}, member string) (string, error) {
	switch op := request.GetStringParam(params, "ConditionalOperator"); op {
	case "":
		return " AND ", nil
	case "AND":
		return " AND ", nil
	case "OR":
		return " OR ", nil
	default:
		return "", legacyValidationException(fmt.Sprintf("Invalid ConditionalOperator value %q: must be one of AND, OR (with %s)", op, member))
	}
}

// legacyComparisonOperators maps every ComparisonOperator value to the
// expression form the key-condition engine consumes. The arity is the
// number of AttributeValueList entries the operator consumes: 1 for the
// comparators, 0 for the existence tests, 2 for BETWEEN, one-or-more for
// IN. The condition-plane translators build typed nodes directly instead
// (legacyOperatorNode); this string rendering serves the KeyConditions
// translation alone.
var legacyComparisonOperators = map[string]struct {
	expression func(name string, values []string) string
	arity      int
}{
	"EQ":           {func(n string, v []string) string { return n + " = " + v[0] }, 1},
	"NE":           {func(n string, v []string) string { return n + " <> " + v[0] }, 1},
	"LE":           {func(n string, v []string) string { return n + " <= " + v[0] }, 1},
	"LT":           {func(n string, v []string) string { return n + " < " + v[0] }, 1},
	"GE":           {func(n string, v []string) string { return n + " >= " + v[0] }, 1},
	"GT":           {func(n string, v []string) string { return n + " > " + v[0] }, 1},
	"NOT_NULL":     {func(n string, v []string) string { return "attribute_exists(" + n + ")" }, 0},
	"NULL":         {func(n string, v []string) string { return "attribute_not_exists(" + n + ")" }, 0},
	"CONTAINS":     {func(n string, v []string) string { return "contains(" + n + ", " + v[0] + ")" }, 1},
	"NOT_CONTAINS": {func(n string, v []string) string { return "NOT contains(" + n + ", " + v[0] + ")" }, 1},
	"BEGINS_WITH":  {func(n string, v []string) string { return "begins_with(" + n + ", " + v[0] + ")" }, 1},
	"IN":           {func(n string, v []string) string { return n + " IN (" + joinComma(v) + ")" }, -1},
	"BETWEEN":      {func(n string, v []string) string { return n + " BETWEEN " + v[0] + " AND " + v[1] }, 2},
}

// legacyComparisonOperatorArity is the arity table the condition-plane
// node builder validates against — the same contract the string map above
// carries for the key plane.
var legacyComparisonOperatorArity = map[string]int{
	"EQ": 1, "NE": 1, "LE": 1, "LT": 1, "GE": 1, "GT": 1,
	"NOT_NULL": 0, "NULL": 0,
	"CONTAINS": 1, "NOT_CONTAINS": 1, "BEGINS_WITH": 1,
	"IN": -1, "BETWEEN": 2,
}

// legacyAttrPath is the single-segment path of a legacy attribute name:
// the legacy members name top-level attributes, so a name containing a
// separator character addresses the attribute so named, never a document
// path through it.
func legacyAttrPath(attrName string) []docPathPart {
	return []docPathPart{{name: attrName}}
}

// legacyOperatorNode builds one {ComparisonOperator, AttributeValueList}
// condition entry as its typed condition-tree node, validating the
// operator against the full legacy enum and the list length against the
// operator's arity.
func legacyOperatorNode(attrName, op string, rawList interface{}, member string) (condNode, error) {
	arity, known := legacyComparisonOperatorArity[op]
	if !known {
		return nil, legacyValidationException(fmt.Sprintf("Invalid %s ComparisonOperator %q", member, op))
	}
	if arity == 0 {
		if rawList != nil {
			if list, ok := rawList.([]interface{}); ok && len(list) > 0 {
				return nil, legacyValidationException(fmt.Sprintf("%s ComparisonOperator %s takes no AttributeValueList entries", member, op))
			}
		}
		if op == "NOT_NULL" {
			return &condFuncNode{kind: condFuncExists, path: legacyAttrPath(attrName)}, nil
		}
		return &condFuncNode{kind: condFuncNotExists, path: legacyAttrPath(attrName)}, nil
	}
	values, err := parseLegacyValueList(rawList, member)
	if err != nil {
		return nil, err
	}
	if arity > 0 && len(values) != arity {
		return nil, legacyValidationException(fmt.Sprintf("%s ComparisonOperator %s requires exactly %d AttributeValueList entr%s, got %d",
			member, op, arity, pluralY(len(values)), len(values)))
	}
	if arity < 0 && len(values) == 0 {
		return nil, legacyValidationException(member + " ComparisonOperator IN requires at least one AttributeValueList entry")
	}
	path := legacyAttrPath(attrName)
	lhs := &operandPath{parts: path}
	switch op {
	case "EQ":
		return &condCompareNode{op: "=", lhs: lhs, rhs: &operandValue{v: values[0]}}, nil
	case "NE":
		return &condCompareNode{op: "<>", lhs: lhs, rhs: &operandValue{v: values[0]}}, nil
	case "LE":
		return &condCompareNode{op: "<=", lhs: lhs, rhs: &operandValue{v: values[0]}}, nil
	case "LT":
		return &condCompareNode{op: "<", lhs: lhs, rhs: &operandValue{v: values[0]}}, nil
	case "GE":
		return &condCompareNode{op: ">=", lhs: lhs, rhs: &operandValue{v: values[0]}}, nil
	case "GT":
		return &condCompareNode{op: ">", lhs: lhs, rhs: &operandValue{v: values[0]}}, nil
	case "CONTAINS":
		return &condFuncNode{kind: condFuncContains, path: path, arg: &operandValue{v: values[0]}}, nil
	case "NOT_CONTAINS":
		return &condNotNode{operand: &condFuncNode{kind: condFuncContains, path: path, arg: &operandValue{v: values[0]}}}, nil
	case "BEGINS_WITH":
		return &condFuncNode{kind: condFuncBegins, path: path, arg: &operandValue{v: values[0]}}, nil
	case "IN":
		candidates := make([]condOperand, len(values))
		for i, v := range values {
			candidates[i] = &operandValue{v: v}
		}
		return &condLegacyInNode{operand: lhs, candidates: candidates}, nil
	case "BETWEEN":
		return &condBetweenNode{operand: lhs, low: &operandValue{v: values[0]}, high: &operandValue{v: values[1]}}, nil
	}
	return nil, legacyValidationException(fmt.Sprintf("Invalid %s ComparisonOperator %q", member, op))
}

// legacyJoinCondNodes combines translated clauses under the resolved
// ConditionalOperator; nil when no clause was translated.
func legacyJoinCondNodes(joiner string, clauses []condNode) condNode {
	if len(clauses) == 0 {
		return nil
	}
	if len(clauses) == 1 {
		return clauses[0]
	}
	if joiner == " OR " {
		return &condOrNode{operands: clauses}
	}
	return &condAndNode{operands: clauses}
}

func joinComma(values []string) string {
	out := ""
	for i, v := range values {
		if i > 0 {
			out += ", "
		}
		out += v
	}
	return out
}

// legacyTokenGen mints unique substitution tokens for one translated
// request: names as #<prefix>n, values as :<prefix>n.
type legacyTokenGen struct {
	namePrefix  string
	valuePrefix string
	names       map[string]string
	values      map[string]*dbstore.AttributeValue
}

func newLegacyTokenGen(namePrefix, valuePrefix string) *legacyTokenGen {
	return &legacyTokenGen{
		namePrefix:  namePrefix,
		valuePrefix: valuePrefix,
		names:       make(map[string]string),
		values:      make(map[string]*dbstore.AttributeValue),
	}
}

func (g *legacyTokenGen) name(attr string) string {
	token := fmt.Sprintf("#%s%d", g.namePrefix, len(g.names))
	g.names[token] = attr
	return token
}

func (g *legacyTokenGen) value(v *dbstore.AttributeValue) string {
	token := fmt.Sprintf(":%s%d", g.valuePrefix, len(g.values))
	g.values[token] = v
	return token
}

// parseLegacyValueList parses an AttributeValueList member into typed
// values, reporting the list position of any malformed entry.
func parseLegacyValueList(raw interface{}, member string) ([]*dbstore.AttributeValue, error) {
	list, ok := raw.([]interface{})
	if !ok {
		return nil, legacyValidationException(member + " AttributeValueList must be a list")
	}
	values := make([]*dbstore.AttributeValue, 0, len(list))
	for i, entry := range list {
		parsed, err := parseAttributeValue(entry)
		if err != nil {
			return nil, legacyValidationException(fmt.Sprintf("%s AttributeValueList entry %d: %s", member, i, err.Error()))
		}
		values = append(values, parsed)
	}
	return values, nil
}

// legacyOperatorClause renders one {ComparisonOperator, AttributeValueList}
// condition entry as the expression clause the key-condition engine
// consumes, validating the operator against the full legacy enum and the
// list length against the operator's arity. Key-plane rendering only: the
// condition-plane translators build typed nodes (legacyOperatorNode).
func legacyOperatorClause(gen *legacyTokenGen, attrName, op string, rawList interface{}, member string) (string, error) {
	spec, ok := legacyComparisonOperators[op]
	if !ok {
		return "", legacyValidationException(fmt.Sprintf("Invalid %s ComparisonOperator %q", member, op))
	}
	if spec.arity == 0 {
		if rawList != nil {
			if list, ok := rawList.([]interface{}); ok && len(list) > 0 {
				return "", legacyValidationException(fmt.Sprintf("%s ComparisonOperator %s takes no AttributeValueList entries", member, op))
			}
		}
		return spec.expression(gen.name(attrName), nil), nil
	}
	values, err := parseLegacyValueList(rawList, member)
	if err != nil {
		return "", err
	}
	if spec.arity > 0 && len(values) != spec.arity {
		return "", legacyValidationException(fmt.Sprintf("%s ComparisonOperator %s requires exactly %d AttributeValueList entr%s, got %d",
			member, op, spec.arity, pluralY(len(values)), len(values)))
	}
	if spec.arity < 0 && len(values) == 0 {
		return "", legacyValidationException(member + " ComparisonOperator IN requires at least one AttributeValueList entry")
	}
	refs := make([]string, len(values))
	for i, v := range values {
		refs[i] = gen.value(v)
	}
	return spec.expression(gen.name(attrName), refs), nil
}

func pluralY(n int) string {
	if n == 1 {
		return "y"
	}
	return "ies"
}

// resolveLegacyWriteCondition applies the Expected exclusivity rules and
// translation for the three conditional writes (PutItem, DeleteItem,
// UpdateItem). It returns nil when the request carries no Expected member.
func resolveLegacyWriteCondition(params map[string]interface{}, conditionExpr string, names map[string]string, values map[string]*dbstore.AttributeValue) (*conditionSpec, error) {
	raw := params["Expected"]
	if raw == nil {
		return nil, nil
	}
	if _, ok := raw.(map[string]interface{}); !ok {
		return nil, legacyValidationException("Expected must be a map of attribute conditions")
	}
	if conditionExpr != "" {
		return nil, legacyValidationException("ConditionExpression and Expected cannot be used together: use ConditionExpression")
	}
	if len(names) > 0 || len(values) > 0 {
		return nil, legacyValidationException("Expected cannot be used with ExpressionAttributeNames or ExpressionAttributeValues")
	}
	spec, err := translateExpected(params)
	if err != nil {
		return nil, err
	}
	if spec.tree == nil {
		return nil, nil
	}
	return &spec, nil
}

// resolveLegacyFilterParam applies the QueryFilter/ScanFilter exclusivity
// rules and translation. It returns nil when the request carries no filter
// member of that name.
func resolveLegacyFilterParam(params map[string]interface{}, member, filterExpr string) (*conditionSpec, error) {
	raw := params[member]
	if raw == nil {
		return nil, nil
	}
	if _, ok := raw.(map[string]interface{}); !ok {
		return nil, legacyValidationException(member + " must be a map of attribute conditions")
	}
	if filterExpr != "" {
		return nil, legacyValidationException("FilterExpression and " + member + " cannot be used together: use FilterExpression")
	}
	// The legacy filter carries its values inline (AttributeValueList), so
	// the expression attribute maps have nothing to substitute — the same
	// exclusivity the Expected path applies: supplying them alongside the
	// legacy member is rejected, never silently ignored.
	if params["ExpressionAttributeNames"] != nil || params["ExpressionAttributeValues"] != nil {
		return nil, legacyValidationException(member + " cannot be used with ExpressionAttributeNames or ExpressionAttributeValues")
	}
	spec, err := translateLegacyFilter(raw, params, member)
	if err != nil {
		return nil, err
	}
	if spec.tree == nil {
		return nil, nil
	}
	return &spec, nil
}

// resolveFilterCondition resolves a read request's filter from whichever
// family carried it — FilterExpression with its substitution maps, or the
// legacy QueryFilter/ScanFilter member — and returns the compiled
// condition. Exclusivity and parse both run here, before the read
// executes: a malformed filter is a ValidationException, never a silently
// empty page. The substitution maps are read only when a native filter
// consumes them; a legacy filter never combines with either map — legacy
// conditional parameters and expression parameters do not mix in a single
// call — so the pairing is rejected inside resolveLegacyFilterParam
// whether or not the map would have been consumed.
func resolveFilterCondition(params map[string]interface{}, legacyMember, filterExpr string) (*compiledCondition, error) {
	legacyFilter, err := resolveLegacyFilterParam(params, legacyMember, filterExpr)
	if err != nil {
		return nil, err
	}
	if legacyFilter != nil {
		return legacyFilter.compile()
	}
	if filterExpr == "" {
		return nil, nil
	}
	names, namesErr := parseExpressionAttributeNames(params)
	if namesErr != nil {
		return nil, namesErr
	}
	values, valuesErr := parseExpressionAttributeValues(params)
	if valuesErr != nil {
		return nil, valuesErr
	}
	return compileConditionExpression(filterExpr, names, values)
}

// legacyMemberPresent reports whether a legacy map member is present and
// non-nil on the wire.
func legacyMemberPresent(params map[string]interface{}, member string) bool {
	return params[member] != nil
}

// translateExpected builds the typed condition tree of the Expected map
// of PutItem, DeleteItem and UpdateItem. Each entry is either the
// {Value, Exists} form or the {ComparisonOperator, AttributeValueList}
// form — carrying both in one entry is the documented
// ValidationException — and multiple entries combine with the
// ConditionalOperator (AND by default).
func translateExpected(params map[string]interface{}) (conditionSpec, error) {
	raw, ok := params["Expected"].(map[string]interface{})
	if !ok || len(raw) == 0 {
		return conditionSpec{}, nil
	}
	joiner, err := legacyJoiner(params, "Expected")
	if err != nil {
		return conditionSpec{}, err
	}
	attrs := make([]string, 0, len(raw))
	for attr := range raw {
		attrs = append(attrs, attr)
	}
	sort.Strings(attrs)

	var clauses []condNode
	for _, attr := range attrs {
		entry, ok := raw[attr].(map[string]interface{})
		if !ok {
			return conditionSpec{}, legacyValidationException("Expected entry " + attr + " must be a condition map")
		}
		_, hasOp := entry["ComparisonOperator"]
		var hasList bool
		if list, ok := entry["AttributeValueList"].([]interface{}); ok && len(list) > 0 {
			hasList = true
		} else if _, present := entry["AttributeValueList"]; present && !ok {
			hasList = true
		}
		_, hasValue := entry["Value"]
		_, hasExists := entry["Exists"]
		if (hasOp || hasList) && (hasValue || hasExists) {
			return conditionSpec{}, legacyValidationException(
				"Expected entry " + attr + " cannot carry Value or Exists together with AttributeValueList and ComparisonOperator")
		}
		if hasOp || hasList {
			op, _ := entry["ComparisonOperator"].(string)
			node, nodeErr := legacyOperatorNode(attr, op, entry["AttributeValueList"], "Expected")
			if nodeErr != nil {
				return conditionSpec{}, nodeErr
			}
			clauses = append(clauses, node)
			continue
		}
		// {Value, Exists} form: Value compares by equality; Exists tests
		// existence (true by default). Exists=false together with Value is
		// the literal conjunction of the two documented checks.
		namePath := legacyAttrPath(attr)
		if hasValue {
			value, valueErr := parseAttributeValue(entry["Value"])
			if valueErr != nil {
				return conditionSpec{}, legacyValidationException("Expected entry " + attr + " Value: " + valueErr.Error())
			}
			clauses = append(clauses, &condCompareNode{op: "=", lhs: &operandPath{parts: namePath}, rhs: &operandValue{v: value}})
		}
		// Exists is a boolean member: present-but-mistyped is a request
		// error, never a silently ignored value.
		if rawExists, present := entry["Exists"]; present {
			exists, isBool := rawExists.(bool)
			if !isBool {
				return conditionSpec{}, legacyValidationException("Expected entry " + attr + " Exists must be a boolean")
			}
			if exists {
				if !hasValue {
					clauses = append(clauses, &condFuncNode{kind: condFuncExists, path: namePath})
				}
			} else {
				clauses = append(clauses, &condFuncNode{kind: condFuncNotExists, path: namePath})
			}
		} else if !hasValue {
			return conditionSpec{}, legacyValidationException("Expected entry " + attr + " carries no condition")
		}
	}

	return conditionSpec{tree: legacyJoinCondNodes(joiner, clauses)}, nil
}

// translateLegacyFilter builds the typed condition tree of a QueryFilter
// or ScanFilter map. Filter conditions are the {ComparisonOperator,
// AttributeValueList} form alone — the Condition shape carries no Value or
// Exists — and entries combine with the ConditionalOperator (AND by
// default).
func translateLegacyFilter(raw interface{}, params map[string]interface{}, member string) (conditionSpec, error) {
	filterMap, ok := raw.(map[string]interface{})
	if !ok || len(filterMap) == 0 {
		return conditionSpec{}, nil
	}
	joiner, err := legacyJoiner(params, member)
	if err != nil {
		return conditionSpec{}, err
	}
	attrs := make([]string, 0, len(filterMap))
	for attr := range filterMap {
		attrs = append(attrs, attr)
	}
	sort.Strings(attrs)

	var clauses []condNode
	for _, attr := range attrs {
		entry, ok := filterMap[attr].(map[string]interface{})
		if !ok {
			return conditionSpec{}, legacyValidationException(member + " entry " + attr + " must be a condition map")
		}
		if _, hasValue := entry["Value"]; hasValue {
			return conditionSpec{}, legacyValidationException(member + " entry " + attr + " cannot carry Value: filter conditions use ComparisonOperator and AttributeValueList")
		}
		if _, hasExists := entry["Exists"]; hasExists {
			return conditionSpec{}, legacyValidationException(member + " entry " + attr + " cannot carry Exists: filter conditions use ComparisonOperator and AttributeValueList")
		}
		op, _ := entry["ComparisonOperator"].(string)
		node, nodeErr := legacyOperatorNode(attr, op, entry["AttributeValueList"], member)
		if nodeErr != nil {
			return conditionSpec{}, nodeErr
		}
		clauses = append(clauses, node)
	}

	return conditionSpec{tree: legacyJoinCondNodes(joiner, clauses)}, nil
}

// translateKeyConditions renders the Query KeyConditions map as the
// KeyConditionExpression that replaces it. Conditions may name only the
// key attributes of the table or index being queried: the partition key
// as an EQ condition is required, and one optional sort-key condition may
// use EQ, LE, LT, GE, GT, BEGINS_WITH or BETWEEN.
func translateKeyConditions(table *dbstore.Table, indexName string, raw interface{}) (conditionSpec, error) {
	keyCondMap, ok := raw.(map[string]interface{})
	if !ok || len(keyCondMap) == 0 {
		return conditionSpec{}, nil
	}
	hashName, rangeName := tableKeyAttributeNames(table, indexName)

	gen := newLegacyTokenGen("lk", "lkv")
	var pkClause, skClause string
	for attr, entryRaw := range keyCondMap {
		entry, ok := entryRaw.(map[string]interface{})
		if !ok {
			return conditionSpec{}, legacyValidationException("KeyConditions entry " + attr + " must be a condition map")
		}
		op, _ := entry["ComparisonOperator"].(string)
		switch {
		case attr == hashName:
			if op != "EQ" {
				return conditionSpec{}, legacyValidationException("KeyConditions partition key condition must use EQ: " + attr)
			}
			clause, err := legacyOperatorClause(gen, attr, op, entry["AttributeValueList"], "KeyConditions")
			if err != nil {
				return conditionSpec{}, err
			}
			pkClause = clause
		case rangeName != "" && attr == rangeName:
			clause, err := legacyKeySortClause(gen, attr, op, entry["AttributeValueList"])
			if err != nil {
				return conditionSpec{}, err
			}
			skClause = clause
		default:
			return conditionSpec{}, legacyValidationException("KeyConditions can only reference key schema attributes, not: " + attr)
		}
	}
	if pkClause == "" {
		return conditionSpec{}, NewAPIError("com.amazon.coral.validate#ValidationException",
			"Query condition missed key schema element: "+hashName, http.StatusBadRequest)
	}
	expr := pkClause
	if skClause != "" {
		expr += " AND " + skClause
	}
	return conditionSpec{Expr: expr, Names: gen.names, Values: gen.values}, nil
}

// legacyKeySortClause renders the optional sort-key condition; the key
// condition family admits only the seven sort-compatible operators.
func legacyKeySortClause(gen *legacyTokenGen, attr, op string, rawList interface{}) (string, error) {
	switch op {
	case "EQ", "LE", "LT", "GE", "GT", "BEGINS_WITH", "BETWEEN":
		return legacyOperatorClause(gen, attr, op, rawList, "KeyConditions")
	default:
		return "", legacyValidationException("KeyConditions sort key condition does not support ComparisonOperator " + op)
	}
}

// tableKeyAttributeNames returns the partition and sort key attribute
// names of the table, or of the index when one is being queried.
func tableKeyAttributeNames(table *dbstore.Table, indexName string) (string, string) {
	if indexName != "" {
		if hashAttr, rangeAttr, found := indexKeyAttributeNames(table, indexName); found || hashAttr != "" {
			return hashAttr, rangeAttr
		}
	}
	hashAttr, rangeAttr := "", ""
	for _, ks := range table.KeySchema {
		switch ks.KeyType {
		case dbstore.KeyTypeHash:
			hashAttr = ks.AttributeName
		case dbstore.KeyTypeRange:
			rangeAttr = ks.AttributeName
		}
	}
	return hashAttr, rangeAttr
}

// parseAttributesToGet reads the AttributesToGet member: an array of one
// or more attribute names, whose absence from an item simply omits the
// attribute from the result. The names are plain attribute names — the
// member predates expression attribute names and cannot address elements
// inside a List or Map.
func parseAttributesToGet(params map[string]interface{}) ([]string, error) {
	raw, present := params["AttributesToGet"]
	if !present || raw == nil {
		return nil, nil
	}
	list, ok := raw.([]interface{})
	if !ok || len(list) == 0 {
		return nil, legacyValidationException("AttributesToGet must be an array of one or more attribute names")
	}
	names := make([]string, 0, len(list))
	for i, entry := range list {
		name, ok := entry.(string)
		if !ok || name == "" {
			return nil, legacyValidationException(fmt.Sprintf("AttributesToGet entry %d must be an attribute name", i))
		}
		names = append(names, name)
	}
	return names, nil
}

// resolveProjectionMembers resolves the projection of a read operation
// from whichever family the request used: ProjectionExpression, or the
// legacy AttributesToGet array. The two never combine. An AttributesToGet
// entry is one plain attribute name — a single path segment, never a
// document path — so a name containing '.' addresses the top-level
// attribute so named.
func resolveProjectionMembers(params map[string]interface{}) ([][]docPathPart, error) {
	projExpr := request.GetStringParam(params, "ProjectionExpression")
	legacy, legacyErr := parseAttributesToGet(params)
	if legacyErr != nil {
		return nil, legacyErr
	}
	if projExpr != "" && legacy != nil {
		return nil, legacyValidationException("ProjectionExpression and AttributesToGet cannot be used together: use ProjectionExpression")
	}
	if projExpr != "" {
		return parseProjectionExpression(params)
	}
	// The absent-member case must stay nil: the read planes gate
	// projection application on projection != nil.
	if len(legacy) == 0 {
		return nil, nil
	}
	projection := make([][]docPathPart, len(legacy))
	for i, name := range legacy {
		projection[i] = []docPathPart{{name: name}}
	}
	return projection, nil
}

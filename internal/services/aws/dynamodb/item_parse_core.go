package dynamodb

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"regexp"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func parseItem(v interface{}) (map[string]*dbstore.AttributeValue, error) {
	if v == nil {
		return nil, nil
	}

	m, ok := v.(map[string]interface{})
	if !ok {
		return nil, nil
	}

	return parseAttributeValueMap(m)
}

// parseKey parses a DynamoDB key map. Key attributes must be one of
// S, N, or B — other attribute types (BOOL, NULL, SS, NS, BS, M, L)
// are rejected to match the DynamoDB key schema constraint. Empty values
// are also rejected (S/N/B must be non-empty).
func parseKey(v interface{}) (map[string]*dbstore.AttributeValue, error) {
	parsed, err := parseItem(v)
	if err != nil {
		return nil, err
	}
	if parsed == nil {
		return nil, nil
	}
	if !validateKeyAttributeValue(parsed) {
		return nil, ErrInvalidParameter
	}
	return parsed, nil
}

func parseAttributeValueMap(m map[string]interface{}) (map[string]*dbstore.AttributeValue, error) {
	result := make(map[string]*dbstore.AttributeValue)
	for k, v := range m {
		parsed, err := parseAttributeValue(v)
		if err != nil {
			return nil, err
		}
		result[k] = parsed
	}
	return result, nil
}

// attributeValueMembers is the AttributeValue union membership: a wire
// attribute value carries exactly one of these members.
var attributeValueMembers = []string{"B", "BOOL", "BS", "L", "M", "N", "NS", "NULL", "S", "SS"}

// wireStringSet normalises the SS member from either wire shape into the
// plain string slice the set validation works on. A non-string element is
// the invalid-set error; a shape the member does not carry at all returns
// nothing, leaving the caller to try the next member.
func wireStringSet(raw interface{}) ([]string, error) {
	switch list := raw.(type) {
	case []string:
		return list, nil
	case []interface{}:
		strs := make([]string, 0, len(list))
		for _, s := range list {
			str, ok := s.(string)
			if !ok {
				return nil, fmt.Errorf("Supplied AttributeValue is not a valid string set")
			}
			strs = append(strs, str)
		}
		return strs, nil
	}
	return nil, nil
}

// wireNumberSet normalises the NS member from either wire shape, validating
// each element as a DynamoDB number on the way through.
func wireNumberSet(raw interface{}) ([]string, error) {
	switch list := raw.(type) {
	case []string:
		for _, str := range list {
			if !isValidDynamoDBNumber(str) {
				return nil, fmt.Errorf("Supplied AttributeValue is not a valid number: %s", str)
			}
		}
		return list, nil
	case []interface{}:
		nums := make([]string, 0, len(list))
		for _, n := range list {
			str, ok := n.(string)
			if !ok {
				return nil, fmt.Errorf("Supplied AttributeValue is not a valid number set")
			}
			if !isValidDynamoDBNumber(str) {
				return nil, fmt.Errorf("Supplied AttributeValue is not a valid number: %s", str)
			}
			nums = append(nums, str)
		}
		return nums, nil
	}
	return nil, nil
}

// wireBinarySet normalises the BS member from either wire shape: the
// []interface{} form carries raw byte slices or base64 strings, the
// []string form base64 strings alone.
func wireBinarySet(raw interface{}) ([][]byte, error) {
	switch list := raw.(type) {
	case []string:
		binaries := make([][]byte, 0, len(list))
		for _, str := range list {
			decoded, err := base64.StdEncoding.DecodeString(str)
			if err != nil {
				return nil, fmt.Errorf("Supplied AttributeValue is not valid base64")
			}
			binaries = append(binaries, decoded)
		}
		return binaries, nil
	case []interface{}:
		binaries := make([][]byte, 0, len(list))
		for _, b := range list {
			if bytes, ok := b.([]byte); ok {
				binaries = append(binaries, bytes)
			} else if str, ok := b.(string); ok {
				decoded, err := base64.StdEncoding.DecodeString(str)
				if err != nil {
					return nil, fmt.Errorf("Supplied AttributeValue is not valid base64")
				}
				binaries = append(binaries, decoded)
			} else {
				return nil, fmt.Errorf("Supplied AttributeValue is not a valid binary set")
			}
		}
		return binaries, nil
	}
	return nil, nil
}

// parseAttributeValue decodes one wire attribute value. The union contract
// is enforced structurally before any member decodes: more than one member
// present, or a NULL member whose value is not true, is rejected instead of
// resolved by member order.
func parseAttributeValue(v interface{}) (*dbstore.AttributeValue, error) {
	if v == nil {
		return dbstore.NullValue(), nil
	}

	m, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Supplied AttributeValue is not a map")
	}

	present := 0
	for _, member := range attributeValueMembers {
		if _, ok := m[member]; ok {
			present++
		}
	}
	if present > 1 {
		return nil, fmt.Errorf("Supplied AttributeValue has more than one datatypes set")
	}
	if nv, ok := m["NULL"]; ok {
		if isTrue, isBool := nv.(bool); !isBool || !isTrue {
			return nil, fmt.Errorf("Null value must be true")
		}
		return dbstore.NullValue(), nil
	}

	if s, ok := m["S"].(string); ok {
		return dbstore.StringValue(s), nil
	}
	if n, ok := m["N"].(string); ok {
		if !isValidDynamoDBNumber(n) {
			return nil, fmt.Errorf("Supplied AttributeValue is not a valid number: %s", n)
		}
		return dbstore.NumberValue(n), nil
	}
	if b, ok := m["B"].([]byte); ok {
		return dbstore.BinaryValue(b), nil
	}
	if b, ok := m["B"].(string); ok {
		decoded, err := base64.StdEncoding.DecodeString(b)
		if err != nil {
			return nil, fmt.Errorf("Supplied AttributeValue is not valid base64")
		}
		return dbstore.BinaryValue(decoded), nil
	}
	if b, ok := m["BOOL"].(bool); ok {
		return dbstore.BoolValue(b), nil
	}
	// The set members arrive in either wire shape — the JSON decoder's
	// []interface{} or the request codec's []string (a rendered key
	// round-trips through internal request parameters) — so each set
	// normalises first and validates once.
	if strs, err := wireStringSet(m["SS"]); strs != nil || err != nil {
		if err != nil {
			return nil, err
		}
		if len(strs) == 0 {
			return nil, fmt.Errorf("Supplied AttributeValue is an empty string set")
		}
		if hasDuplicateString(strs) {
			return nil, fmt.Errorf("Supplied AttributeValue string set contains duplicates")
		}
		return dbstore.StringSet(strs), nil
	}
	if nums, err := wireNumberSet(m["NS"]); nums != nil || err != nil {
		if err != nil {
			return nil, err
		}
		if len(nums) == 0 {
			return nil, fmt.Errorf("Supplied AttributeValue is an empty number set")
		}
		if hasDuplicateNumber(nums) {
			return nil, fmt.Errorf("Supplied AttributeValue number set contains duplicates")
		}
		return dbstore.NumberSet(nums), nil
	}
	if binaries, err := wireBinarySet(m["BS"]); binaries != nil || err != nil {
		if err != nil {
			return nil, err
		}
		if len(binaries) == 0 {
			return nil, fmt.Errorf("Supplied AttributeValue is an empty binary set")
		}
		if hasDuplicateBinary(binaries) {
			return nil, fmt.Errorf("Supplied AttributeValue binary set contains duplicates")
		}
		return dbstore.BinarySet(binaries), nil
	}
	if mm, ok := m["M"].(map[string]interface{}); ok {
		parsed, err := parseAttributeValueMap(mm)
		if err != nil {
			return nil, err
		}
		return dbstore.MapValue(parsed), nil
	}
	if l, ok := m["L"].([]interface{}); ok {
		list := make([]*dbstore.AttributeValue, 0, len(l))
		for _, item := range l {
			parsed, err := parseAttributeValue(item)
			if err != nil {
				return nil, err
			}
			list = append(list, parsed)
		}
		return dbstore.ListValue(list), nil
	}

	return nil, fmt.Errorf("Supplied AttributeValue is empty: exactly one member must be set")
}

// buildItemResponse renders a typed attribute map in the wire shape; the
// codec itself lives in the store package next to the AttributeValue type.
func buildItemResponse(attrs map[string]*dbstore.AttributeValue) map[string]interface{} {
	return dbstore.BuildItemWire(attrs)
}

// buildAttributeValueResponse renders one typed attribute value in the wire
// shape, delegating to the store package's single wire codec.
func buildAttributeValueResponse(av *dbstore.AttributeValue) map[string]interface{} {
	return dbstore.BuildAttributeValueWire(av)
}

func buildItemsResponse(items []*dbstore.Item) []map[string]interface{} {
	result := make([]map[string]interface{}, len(items))
	for i, item := range items {
		result[i] = buildItemResponse(item.Attributes)
	}
	return result
}

func buildLastEvaluatedKeyWithIndex(item *dbstore.Item, table *dbstore.Table, indexName string) map[string]interface{} {
	if item == nil || item.Key == nil {
		return nil
	}

	key := buildItemResponse(item.Key)
	if indexName == "" || table == nil || item.Attributes == nil {
		return key
	}

	var indexKeySchema []*dbstore.KeySchemaElement
	for _, gsi := range table.GlobalSecondaryIndexes {
		if gsi.IndexName == indexName {
			indexKeySchema = gsi.KeySchema
			break
		}
	}
	if indexKeySchema == nil {
		for _, lsi := range table.LocalSecondaryIndexes {
			if lsi.IndexName == indexName {
				indexKeySchema = lsi.KeySchema
				break
			}
		}
	}

	for _, ks := range indexKeySchema {
		if _, exists := key[ks.AttributeName]; !exists {
			if attrVal, ok := item.Attributes[ks.AttributeName]; ok {
				key[ks.AttributeName] = buildAttributeValueResponse(attrVal)
			}
		}
	}

	return key
}

func parseExpressionAttributeNames(params map[string]interface{}) (map[string]string, error) {
	names := make(map[string]string)
	ean, ok := params["ExpressionAttributeNames"].(map[string]interface{})
	if !ok {
		return names, nil
	}
	for k, v := range ean {
		vs, ok := v.(string)
		if !ok {
			return nil, ErrInvalidParameter
		}
		names[k] = vs
	}
	return names, nil
}

func parseExpressionAttributeValues(params map[string]interface{}) (map[string]*dbstore.AttributeValue, error) {
	values := make(map[string]*dbstore.AttributeValue)
	if eav, ok := params["ExpressionAttributeValues"].(map[string]interface{}); ok {
		for k, v := range eav {
			parsed, err := parseAttributeValue(v)
			if err != nil {
				return nil, fmt.Errorf("invalid value for expression attribute %q: %w", k, err)
			}
			values[k] = parsed
		}
	}
	return values, nil
}

// parseProjectionExpression parses the ProjectionExpression member into
// resolved document paths. Each comma-separated token goes through the
// document-path grammar FIRST and its expression attribute names are then
// resolved per segment, so an alias standing for a name that itself
// contains '.' or '[' stays a single top-level segment — a resolved value
// never re-enters the path grammar. An alias the names map does not
// define is ErrInvalidParameter, the rejection every expression plane
// applies; a token the path grammar rejects is the same error.
func parseProjectionExpression(params map[string]interface{}) ([][]docPathPart, error) {
	projExpr := request.GetStringParam(params, "ProjectionExpression")
	if projExpr == "" {
		return nil, nil
	}

	names, err := parseExpressionAttributeNames(params)
	if err != nil {
		return nil, err
	}
	var projection [][]docPathPart

	attrs := splitAndTrim(projExpr, ",")
	for _, attr := range attrs {
		parts, resolveErr := resolveDocPathParts(attr, names)
		if resolveErr != nil {
			return nil, ErrInvalidParameter
		}
		projection = append(projection, parts)
	}

	return projection, nil
}

// applyProjection narrows an item's attributes to the requested document
// paths. The paths arrive already parsed and resolved (see
// parseProjectionExpression); a whole path's leaf value is returned under
// its own name, while a multi-segment path's leaf is rebuilt inside the
// container shape the path descends through, merged across paths sharing
// a top-level name.
func applyProjection(attrs map[string]*dbstore.AttributeValue, projection [][]docPathPart) map[string]*dbstore.AttributeValue {
	result := make(map[string]*dbstore.AttributeValue)
	for _, parts := range projection {
		val := getDocPathValue(attrs, parts)
		if val == nil {
			continue
		}
		if len(parts) <= 1 {
			if len(parts) == 1 {
				result[parts[0].name] = val
			}
			continue
		}
		topKey := parts[0].name
		nested := buildNestedValue(parts[1:], val)
		if existing, ok := result[topKey]; ok {
			mergeNestedValues(existing, nested, parts[1:])
		} else {
			result[topKey] = nested
		}
	}
	return result
}

func buildNestedValue(parts []docPathPart, leaf *dbstore.AttributeValue) *dbstore.AttributeValue {
	if len(parts) == 0 {
		return leaf
	}
	part := parts[0]
	if part.isIndex {
		child := buildNestedValue(parts[1:], leaf)
		return &dbstore.AttributeValue{L: []*dbstore.AttributeValue{child}}
	}
	m := map[string]*dbstore.AttributeValue{
		part.name: buildNestedValue(parts[1:], leaf),
	}
	return &dbstore.AttributeValue{M: m}
}

func mergeNestedValues(existing *dbstore.AttributeValue, nested *dbstore.AttributeValue, parts []docPathPart) {
	if len(parts) == 0 || existing == nil {
		return
	}
	part := parts[0]
	if part.isIndex {
		if existing.L == nil {
			return
		}
		if len(nested.L) == 0 {
			return
		}
		existing.L = append(existing.L, nested.L[0])
		if len(parts) > 1 && existing.L[len(existing.L)-1] != nil {
			mergeNestedValues(existing.L[len(existing.L)-1], nested.L[0], parts[1:])
		}
	} else {
		if existing.M == nil || nested.M == nil {
			return
		}
		for k, v := range nested.M {
			if existing.M[k] == nil {
				existing.M[k] = v
			} else if len(parts) > 1 {
				mergeNestedValues(existing.M[k], v, parts[1:])
			}
		}
	}
}

func parseExclusiveStartKey(params map[string]interface{}) (map[string]*dbstore.AttributeValue, error) {
	raw, present := params["ExclusiveStartKey"]
	if !present || raw == nil {
		return nil, nil
	}
	esk, ok := raw.(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}
	if len(esk) == 0 {
		return nil, ErrInvalidParameter
	}
	return parseAttributeValueMap(esk)
}

func getExpressionAttributes(params map[string]interface{}) (map[string]string, map[string]*dbstore.AttributeValue, error) {
	names, err := parseExpressionAttributeNames(params)
	if err != nil {
		return nil, nil, err
	}
	values, err := parseExpressionAttributeValues(params)
	if err != nil {
		return nil, nil, err
	}
	return names, values, nil
}

// dynamoDBNumberPattern matches the decimal grammar DynamoDB accepts for
// Number values: an optional sign, digits with an optional fraction, and an
// optional decimal exponent. Fraction forms such as "1/3" that big.Rat
// would accept are not part of the grammar.
var dynamoDBNumberPattern = regexp.MustCompile(`^[+-]?(\d+(\.\d*)?|\.\d+)([eE][+-]?\d+)?$`)

// Number magnitude boundaries as exact rationals: values must satisfy
// 1E-130 <= |v| < 1E+126 (the documented range tops out at
// 9.9999999999999999999999999999999999999E+125, which every 38-digit
// mantissa stays below).
var (
	numberMaxMagnitude = pow10Rat(126)
	numberMinMagnitude = pow10Rat(-130)
)

func pow10Rat(exp int) *big.Rat {
	ten := big.NewInt(10)
	if exp >= 0 {
		num := new(big.Int).Exp(ten, big.NewInt(int64(exp)), nil)
		return new(big.Rat).SetInt(num)
	}
	den := new(big.Int).Exp(ten, big.NewInt(int64(-exp)), nil)
	return new(big.Rat).SetFrac(big.NewInt(1), den)
}

func isValidDynamoDBNumber(n string) bool {
	if n == "" {
		return false
	}
	if !dynamoDBNumberPattern.MatchString(n) {
		return false
	}
	value, ok := new(big.Rat).SetString(n)
	if !ok {
		return false
	}
	if value.Sign() == 0 {
		return true
	}
	abs := new(big.Rat).Abs(value)
	if abs.Cmp(numberMaxMagnitude) >= 0 || abs.Cmp(numberMinMagnitude) < 0 {
		return false
	}
	return dbstore.CountSignificantDigits(n) <= numberMaxSignificantDigits
}

func buildUpdatedAttributesResponse(attrs map[string]*dbstore.AttributeValue, updatedAttrNames []string) map[string]interface{} {
	result := make(map[string]interface{})
	for _, attrName := range updatedAttrNames {
		if v, ok := attrs[attrName]; ok {
			result[attrName] = buildAttributeValueResponse(v)
		}
	}
	return result
}

// hasDuplicateString reports whether the string set contains a repeated
// element. Set elements must be unique; empty string and binary elements are
// allowed.
func hasDuplicateString(items []string) bool {
	seen := make(map[string]bool, len(items))
	for _, item := range items {
		if seen[item] {
			return true
		}
		seen[item] = true
	}
	return false
}

// hasDuplicateNumber reports whether the number set contains a repeated
// value, comparing numerically so that "1" and "1.0" count as duplicates.
func hasDuplicateNumber(items []string) bool {
	seen := make(map[string]bool, len(items))
	for _, item := range items {
		normalized := normalizeNumberString(item)
		if seen[normalized] {
			return true
		}
		seen[normalized] = true
	}
	return false
}

// hasDuplicateBinary reports whether the binary set contains a repeated
// element.
func hasDuplicateBinary(items [][]byte) bool {
	seen := make(map[string]bool, len(items))
	for _, item := range items {
		key := string(item)
		if seen[key] {
			return true
		}
		seen[key] = true
	}
	return false
}

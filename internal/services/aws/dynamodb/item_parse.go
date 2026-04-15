package dynamodb

import (
	"encoding/base64"
	"math/big"
	"strings"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func parseItem(v interface{}) map[string]*dbstore.AttributeValue {
	if v == nil {
		return nil
	}

	m, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	return parseAttributeValueMap(m)
}

func parseKey(v interface{}) map[string]*dbstore.AttributeValue {
	return parseItem(v)
}

func parseAttributeValueMap(m map[string]interface{}) map[string]*dbstore.AttributeValue {
	result := make(map[string]*dbstore.AttributeValue)
	for k, v := range m {
		result[k] = parseAttributeValue(v)
	}
	return result
}

func parseAttributeValue(v interface{}) *dbstore.AttributeValue {
	if v == nil {
		return dbstore.NullValue()
	}

	m, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	if s, ok := m["S"].(string); ok {
		return dbstore.StringValue(s)
	}
	if n, ok := m["N"].(string); ok {
		if !isValidDynamoDBNumber(n) {
			return nil
		}
		return dbstore.NumberValue(n)
	}
	if b, ok := m["B"].([]byte); ok {
		return dbstore.BinaryValue(b)
	}
	if b, ok := m["B"].(string); ok {
		decoded, err := base64.StdEncoding.DecodeString(b)
		if err != nil {
			return nil
		}
		return dbstore.BinaryValue(decoded)
	}
	if b, ok := m["BOOL"].(bool); ok {
		return dbstore.BoolValue(b)
	}
	if _, ok := m["NULL"]; ok {
		return dbstore.NullValue()
	}
	if ss, ok := m["SS"].([]interface{}); ok {
		var strs []string
		for _, s := range ss {
			if str, ok := s.(string); ok {
				strs = append(strs, str)
			}
		}
		if len(strs) == 0 {
			return nil
		}
		return dbstore.StringSet(strs)
	}
	if ns, ok := m["NS"].([]interface{}); ok {
		var nums []string
		for _, n := range ns {
			if str, ok := n.(string); ok {
				if isValidDynamoDBNumber(str) {
					nums = append(nums, str)
				}
			}
		}
		if len(nums) == 0 {
			return nil
		}
		return dbstore.NumberSet(nums)
	}
	if bs, ok := m["BS"].([]interface{}); ok {
		var binaries [][]byte
		for _, b := range bs {
			if bytes, ok := b.([]byte); ok {
				binaries = append(binaries, bytes)
			} else if str, ok := b.(string); ok {
				decoded, err := base64.StdEncoding.DecodeString(str)
				if err != nil {
					return nil
				}
				binaries = append(binaries, decoded)
			}
		}
		if len(binaries) == 0 {
			return nil
		}
		return dbstore.BinarySet(binaries)
	}
	if mm, ok := m["M"].(map[string]interface{}); ok {
		return dbstore.MapValue(parseAttributeValueMap(mm))
	}
	if l, ok := m["L"].([]interface{}); ok {
		var list []*dbstore.AttributeValue
		for _, item := range l {
			if parsed := parseAttributeValue(item); parsed != nil {
				list = append(list, parsed)
			}
		}
		return dbstore.ListValue(list)
	}

	return nil
}

func buildItemResponse(attrs map[string]*dbstore.AttributeValue) map[string]interface{} {
	if attrs == nil {
		return nil
	}

	result := make(map[string]interface{})
	for k, v := range attrs {
		result[k] = buildAttributeValueResponse(v)
	}
	return result
}

func buildAttributeValueResponse(av *dbstore.AttributeValue) map[string]interface{} {
	if av == nil {
		return nil
	}

	result := make(map[string]interface{})

	if av.S != nil {
		result["S"] = *av.S
	}
	if av.N != nil {
		result["N"] = *av.N
	}
	if av.B != nil {
		result["B"] = base64.StdEncoding.EncodeToString(av.B)
	}
	if av.BOOL != nil {
		result["BOOL"] = *av.BOOL
	}
	if av.NULL != nil && *av.NULL {
		result["NULL"] = true
	}
	if av.SS != nil {
		result["SS"] = av.SS
	}
	if av.NS != nil {
		result["NS"] = av.NS
	}
	if av.BS != nil {
		var encodedBS []string
		for _, b := range av.BS {
			encodedBS = append(encodedBS, base64.StdEncoding.EncodeToString(b))
		}
		result["BS"] = encodedBS
	}
	if av.M != nil {
		m := make(map[string]interface{})
		for k, v := range av.M {
			m[k] = buildAttributeValueResponse(v)
		}
		result["M"] = m
	}
	if av.L != nil {
		var l []interface{}
		for _, item := range av.L {
			l = append(l, buildAttributeValueResponse(item))
		}
		result["L"] = l
	}

	return result
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

func parseExpressionAttributeNames(params map[string]interface{}) map[string]string {
	names := make(map[string]string)
	if ean, ok := params["ExpressionAttributeNames"].(map[string]interface{}); ok {
		for k, v := range ean {
			if vs, ok := v.(string); ok {
				names[k] = vs
			}
		}
	}
	return names
}

func parseExpressionAttributeValues(params map[string]interface{}) map[string]*dbstore.AttributeValue {
	values := make(map[string]*dbstore.AttributeValue)
	if eav, ok := params["ExpressionAttributeValues"].(map[string]interface{}); ok {
		for k, v := range eav {
			values[k] = parseAttributeValue(v)
		}
	}
	return values
}

func parseProjectionExpression(params map[string]interface{}) []string {
	projExpr := request.GetStringParam(params, "ProjectionExpression")
	if projExpr == "" {
		return nil
	}

	names := parseExpressionAttributeNames(params)
	var projection []string

	attrs := splitAndTrim(projExpr, ",")
	for _, attr := range attrs {
		resolved := resolvePathTokens(attr, names)
		projection = append(projection, resolved)
	}

	return projection
}

func resolvePathTokens(path string, names map[string]string) string {
	if names == nil {
		return path
	}

	var result strings.Builder
	current := ""
	i := 0

	for i < len(path) {
		c := path[i]
		if c == '.' || c == '[' {
			if current != "" {
				if resolved, ok := names[current]; ok {
					result.WriteString(resolved)
				} else {
					result.WriteString(current)
				}
				current = ""
			}
			result.WriteByte(c)
			i++
		} else if c == ']' {
			if current != "" {
				result.WriteString(current)
				current = ""
			}
			result.WriteByte(c)
			i++
		} else {
			current += string(c)
			i++
		}
	}

	if current != "" {
		if resolved, ok := names[current]; ok {
			result.WriteString(resolved)
		} else {
			result.WriteString(current)
		}
	}

	return result.String()
}

func applyProjection(attrs map[string]*dbstore.AttributeValue, projection []string) map[string]*dbstore.AttributeValue {
	result := make(map[string]*dbstore.AttributeValue)
	for _, path := range projection {
		val := getNestedAttrValueForProjection(attrs, path)
		if val == nil {
			continue
		}
		parts := parseProjPathParts(path)
		if len(parts) <= 1 {
			key := path
			if len(parts) == 1 {
				key = parts[0].name
			}
			result[key] = val
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

func buildNestedValue(parts []projPathPart, leaf *dbstore.AttributeValue) *dbstore.AttributeValue {
	if len(parts) == 0 {
		return leaf
	}
	part := parts[0]
	if part.isIndex {
		l := make([]*dbstore.AttributeValue, part.index+1)
		l[part.index] = buildNestedValue(parts[1:], leaf)
		return &dbstore.AttributeValue{L: l}
	}
	m := map[string]*dbstore.AttributeValue{
		part.name: buildNestedValue(parts[1:], leaf),
	}
	return &dbstore.AttributeValue{M: m}
}

func mergeNestedValues(existing *dbstore.AttributeValue, nested *dbstore.AttributeValue, parts []projPathPart) {
	if len(parts) == 0 || existing == nil {
		return
	}
	part := parts[0]
	if part.isIndex {
		if existing.L == nil || part.index >= len(existing.L) {
			return
		}
		if existing.L[part.index] == nil {
			existing.L[part.index] = nested.L[part.index]
		} else if len(parts) > 1 {
			mergeNestedValues(existing.L[part.index], nested.L[part.index], parts[1:])
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

func getNestedAttrValueForProjection(attrs map[string]*dbstore.AttributeValue, path string) *dbstore.AttributeValue {
	parts := parseProjPathParts(path)
	if len(parts) == 0 {
		if v, ok := attrs[path]; ok {
			return v
		}
		return nil
	}

	if len(parts) == 1 {
		if v, ok := attrs[parts[0].name]; ok {
			return v
		}
		return nil
	}

	topAttr := parts[0].name
	topVal, ok := attrs[topAttr]
	if !ok {
		return nil
	}

	current := topVal
	for i := 1; i < len(parts); i++ {
		part := parts[i]
		if part.isIndex {
			if current.L == nil || part.index >= len(current.L) {
				return nil
			}
			current = current.L[part.index]
		} else {
			if current.M == nil {
				return nil
			}
			if v, exists := current.M[part.name]; exists {
				current = v
			} else {
				return nil
			}
		}
	}
	return current
}

type projPathPart struct {
	name    string
	index   int
	isIndex bool
}

func parseProjPathParts(path string) []projPathPart {
	var parts []projPathPart
	current := ""
	i := 0

	for i < len(path) {
		c := path[i]
		switch c {
		case '.':
			if current != "" {
				parts = append(parts, projPathPart{name: current})
				current = ""
			}
			i++
		case '[':
			if current != "" {
				parts = append(parts, projPathPart{name: current})
				current = ""
			}
			j := i + 1
			for j < len(path) && path[j] != ']' {
				j++
			}
			idxStr := path[i+1 : j]
			idx := 0
			for _, ch := range idxStr {
				if ch >= '0' && ch <= '9' {
					idx = idx*10 + int(ch-'0')
				}
			}
			parts = append(parts, projPathPart{index: idx, isIndex: true})
			i = j + 1
		default:
			current += string(c)
			i++
		}
	}

	if current != "" {
		parts = append(parts, projPathPart{name: current})
	}

	return parts
}

func parseExclusiveStartKey(params map[string]interface{}) map[string]*dbstore.AttributeValue {
	if esk, ok := params["ExclusiveStartKey"].(map[string]interface{}); ok {
		return parseAttributeValueMap(esk)
	}
	return nil
}

func getExpressionAttributes(params map[string]interface{}) (map[string]string, map[string]*dbstore.AttributeValue) {
	names := make(map[string]string)
	values := make(map[string]*dbstore.AttributeValue)

	if namesMap, ok := params["ExpressionAttributeNames"].(map[string]interface{}); ok {
		for k, v := range namesMap {
			if vs, ok := v.(string); ok {
				names[k] = vs
			}
		}
	}

	if valuesMap, ok := params["ExpressionAttributeValues"].(map[string]interface{}); ok {
		for k, v := range valuesMap {
			if attrVal := parseAttributeValue(v); attrVal != nil {
				values[k] = attrVal
			}
		}
	}

	return names, values
}

func isValidDynamoDBNumber(n string) bool {
	if n == "" {
		return false
	}

	_, ok := new(big.Rat).SetString(n)
	if !ok {
		return false
	}

	digitCount := 0
	for _, c := range n {
		if c >= '0' && c <= '9' {
			digitCount++
		}
	}
	return digitCount <= 38
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

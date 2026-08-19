package dynamodb

import (
	"fmt"
	"net/http"
	"strings"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

const maxItemSizeBytes = 400 * 1024

func isKeyAttribute(table *dbstore.Table, attrName string) bool {
	for _, ks := range table.KeySchema {
		if ks.AttributeName == attrName {
			return true
		}
	}
	return false
}

func validateNotKeyAttributes(table *dbstore.Table, paths []string) error {
	for _, p := range paths {
		if isKeyAttribute(table, p) {
			return NewAPIError(
				"com.amazon.coral.validate#ValidationException",
				fmt.Sprintf("One or more parameter values were invalid: Cannot update attribute %s. This attribute is part of the key schema", p),
				400,
			)
		}
	}
	return nil
}

func getReturnConsumedCapacity(params map[string]interface{}) string {
	return request.GetStringParam(params, "ReturnConsumedCapacity")
}

func buildConsumedCapacityResponse(tableName string, capacityUnits float64) map[string]interface{} {
	return map[string]interface{}{
		"TableName":     tableName,
		"CapacityUnits": capacityUnits,
	}
}

func buildConsumedCapacityResponseWithIndex(tableName string, indexName string, capacityUnits float64, isLSI bool) map[string]interface{} {
	resp := map[string]interface{}{
		"TableName":     tableName,
		"CapacityUnits": capacityUnits,
	}
	if indexName != "" {
		resp["Table"] = map[string]interface{}{"CapacityUnits": capacityUnits}
		if isLSI {
			resp["LocalSecondaryIndexes"] = map[string]interface{}{
				indexName: map[string]interface{}{"CapacityUnits": capacityUnits},
			}
		} else {
			resp["GlobalSecondaryIndexes"] = map[string]interface{}{
				indexName: map[string]interface{}{"CapacityUnits": capacityUnits},
			}
		}
	}
	return resp
}

func (s *DynamoDBService) extractKeyFromItem(table *dbstore.Table, item map[string]*dbstore.AttributeValue) map[string]*dbstore.AttributeValue {
	key := make(map[string]*dbstore.AttributeValue)

	for _, ks := range table.KeySchema {
		attr, ok := item[ks.AttributeName]
		if !ok {
			return nil
		}
		key[ks.AttributeName] = attr
	}

	return key
}

func calculateItemSize(item map[string]*dbstore.AttributeValue) int64 {
	var size int64
	for attrName, av := range item {
		size += int64(len(attrName))
		size += calculateAttributeValueSize(av)
	}
	return size
}

func calculateAttributeValueSize(av *dbstore.AttributeValue) int64 {
	if av == nil {
		return 0
	}

	if av.S != nil {
		return int64(len(*av.S))
	}
	if av.N != nil {
		return calculateNumberSize(*av.N)
	}
	if av.B != nil {
		return int64(len(av.B))
	}
	if av.BOOL != nil {
		return 1
	}
	if av.NULL != nil {
		return 1
	}
	if av.SS != nil {
		var size int64
		for _, s := range av.SS {
			size += int64(len(s))
		}
		return size
	}
	if av.NS != nil {
		var size int64
		for _, n := range av.NS {
			size += calculateNumberSize(n)
		}
		return size
	}
	if av.BS != nil {
		var size int64
		for _, b := range av.BS {
			size += int64(len(b))
		}
		return size
	}
	if av.M != nil {
		var size int64 = 3
		for k, v := range av.M {
			size += int64(len(k))
			size += calculateAttributeValueSize(v)
		}
		return size
	}
	if av.L != nil {
		var size int64 = 3
		for _, v := range av.L {
			size += calculateAttributeValueSize(v)
		}
		return size
	}
	return 0
}

// calculateNumberSize returns the size in bytes for a DynamoDB Number value.
// AWS counts each pair of significant digits as 1 byte, minimum 1 byte.
func calculateNumberSize(numStr string) int64 {
	if numStr == "" {
		return 1
	}
	significantDigits := 0
	for _, c := range numStr {
		if c >= '0' && c <= '9' {
			significantDigits++
		}
	}
	if significantDigits == 0 {
		return 1
	}
	return int64((significantDigits + 1) / 2)
}

func validateKeyValueNotEmpty(key map[string]*dbstore.AttributeValue) bool {
	for _, av := range key {
		if av == nil {
			return false
		}
		if av.S != nil && *av.S == "" {
			return false
		}
		if av.N != nil && *av.N == "" {
			return false
		}
		if av.B != nil && len(av.B) == 0 {
			return false
		}
	}
	return true
}

func tokenizeExpression(expr string) []string {
	var tokens []string
	var current strings.Builder
	i := 0

	for i < len(expr) {
		ch := expr[i]

		if ch == ' ' || ch == '\t' || ch == '\n' {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			i++
			continue
		}

		if ch == '\'' {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			j := i + 1
			for j < len(expr) && expr[j] != '\'' {
				j++
			}
			if j < len(expr) {
				tokens = append(tokens, expr[i:j+1])
				i = j + 1
			} else {
				tokens = append(tokens, expr[i:])
				i = j
			}
			continue
		}

		if ch == '=' || ch == '<' || ch == '>' || ch == ',' {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}

			if ch == '<' && i+1 < len(expr) {
				next := expr[i+1]
				if next == '=' || next == '>' {
					tokens = append(tokens, string(ch)+string(next))
					i += 2
					continue
				}
			}

			if ch == '>' && i+1 < len(expr) && expr[i+1] == '=' {
				tokens = append(tokens, ">=")
				i += 2
				continue
			}

			tokens = append(tokens, string(ch))
			i++
			continue
		}

		if ch == '(' {
			if current.Len() > 0 && isConditionFunctionName(current.String()) {
				current.WriteByte(ch)
				depth := 1
				i++
				for i < len(expr) && depth > 0 {
					c := expr[i]
					current.WriteByte(c)
					if c == '(' {
						depth++
					} else if c == ')' {
						depth--
					}
					i++
				}
				tokens = append(tokens, current.String())
				current.Reset()
				continue
			}
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, "(")
			i++
			continue
		}

		if ch == ')' {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, ")")
			i++
			continue
		}

		current.WriteByte(ch)
		i++
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

func splitAndTrim(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || string(s[i]) == sep {
			part := strings.TrimSpace(s[start:i])
			if part != "" {
				result = append(result, part)
			}
			start = i + 1
		}
	}
	return result
}

func isIdentRune(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}
	return true
}

// isConditionFunctionName returns true for DynamoDB expression function
// names that should absorb their parenthetical arguments as a single token.
// This prevents keywords like IN from being treated as function calls
// when written as IN(:v1, :v2) without a space before the parenthesis.
func isConditionFunctionName(s string) bool {
	switch s {
	case "attribute_exists", "attribute_not_exists", "attribute_type",
		"begins_with", "contains", "size":
		return true
	}
	return false
}

// attributeValueType reports the scalar type descriptor ("S", "N", or "B")
// carried by an attribute value; set, document, boolean, and null values
// carry no scalar type and cannot serve as key attributes.
func attributeValueType(av *dbstore.AttributeValue) string {
	switch {
	case av.S != nil:
		return string(dbstore.ScalarAttributeTypeS)
	case av.N != nil:
		return string(dbstore.ScalarAttributeTypeN)
	case av.B != nil:
		return string(dbstore.ScalarAttributeTypeB)
	default:
		return ""
	}
}

// attributeTypeDefinitions indexes the table's attribute definitions by
// attribute name.
func attributeTypeDefinitions(table *dbstore.Table) map[string]dbstore.ScalarAttributeType {
	defs := make(map[string]dbstore.ScalarAttributeType, len(table.AttributeDefinitions))
	for _, def := range table.AttributeDefinitions {
		defs[def.AttributeName] = def.AttributeType
	}
	return defs
}

// keyTypeMismatchError builds the ValidationException DynamoDB answers a
// wrong-typed key attribute with.
func keyTypeMismatchError(attrName string, expected, actual dbstore.ScalarAttributeType) error {
	return NewAPIError("com.amazon.coral.validate#ValidationException",
		fmt.Sprintf("Type mismatch for key %s expected: %s actual: %s", attrName, expected, actual),
		http.StatusBadRequest)
}

// validateKeySchemaAttrTypes checks the attributes named by a key schema
// against the table's attribute definitions. Values of non-scalar types are
// reported as mismatches against the expected scalar type.
func validateKeySchemaAttrTypes(table *dbstore.Table, keySchema []*dbstore.KeySchemaElement, attrs map[string]*dbstore.AttributeValue) error {
	defs := attributeTypeDefinitions(table)
	for _, ks := range keySchema {
		attr, ok := attrs[ks.AttributeName]
		if !ok || attr == nil {
			continue
		}
		expected, hasDef := defs[ks.AttributeName]
		if !hasDef {
			continue
		}
		actual := attributeValueType(attr)
		if actual == "" {
			return keyTypeMismatchError(ks.AttributeName, expected, expected)
		}
		if actual != string(expected) {
			return keyTypeMismatchError(ks.AttributeName, expected, dbstore.ScalarAttributeType(actual))
		}
	}
	return nil
}

// validateKeyTypes checks a supplied primary key against the table's key
// schema and attribute definitions.
func validateKeyTypes(table *dbstore.Table, key map[string]*dbstore.AttributeValue) error {
	return validateKeySchemaAttrTypes(table, table.KeySchema, key)
}

// validateItemKeyTypes validates the primary key carried by an item along
// with any secondary index key attributes the item carries, matching the
// BatchWriteItem contract that index key attribute types must agree with the
// schema definitions.
func validateItemKeyTypes(table *dbstore.Table, item map[string]*dbstore.AttributeValue) error {
	if err := validateKeyTypes(table, item); err != nil {
		return err
	}
	for _, idx := range table.GlobalSecondaryIndexes {
		if err := validateKeySchemaAttrTypes(table, idx.KeySchema, item); err != nil {
			return err
		}
	}
	for _, idx := range table.LocalSecondaryIndexes {
		if err := validateKeySchemaAttrTypes(table, idx.KeySchema, item); err != nil {
			return err
		}
	}
	return nil
}

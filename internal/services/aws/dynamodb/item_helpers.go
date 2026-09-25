package dynamodb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// errScanSufficient is returned from a scan callback to signal that
// enough items have been collected and the scan can stop early.
var errScanSufficient = errors.New("scan sufficient items collected")

// stringSetCanonical is the canonical form of string-set and binary-set
// membership: the member itself.
func stringSetCanonical(s string) string { return s }

// addStringSetMembers appends every incoming member whose canonical form
// the stored set does not already carry, preserving the stored order and
// appending the new members in statement order. The canonical form is both
// the identity key and the stored representation of an appended member —
// number sets pass their canonical number form so that two spellings of
// one numeric value ("1" and "01") are one set member, while string sets
// pass the member itself. Previously stored members keep their stored
// spelling untouched.
func addStringSetMembers(existing, incoming []string, canonical func(string) string) []string {
	present := make(map[string]bool, len(existing))
	for _, m := range existing {
		present[canonical(m)] = true
	}
	for _, m := range incoming {
		key := canonical(m)
		if !present[key] {
			present[key] = true
			existing = append(existing, key)
		}
	}
	return existing
}

// removeStringSetMembers filters the stored set down to the members whose
// canonical form the removal set does not carry. Survivors keep their
// stored spelling: the canonical form serves as the identity key alone.
func removeStringSetMembers(existing, removals []string, canonical func(string) string) []string {
	toDelete := make(map[string]bool, len(removals))
	for _, m := range removals {
		toDelete[canonical(m)] = true
	}
	remaining := make([]string, 0, len(existing))
	for _, m := range existing {
		if !toDelete[canonical(m)] {
			remaining = append(remaining, m)
		}
	}
	return remaining
}

// addBinarySetMembers appends every incoming member the stored set does
// not already carry, preserving stored order. A binary member is its own
// identity.
func addBinarySetMembers(existing, incoming [][]byte) [][]byte {
	present := make(map[string]bool, len(existing))
	for _, b := range existing {
		present[string(b)] = true
	}
	for _, b := range incoming {
		if !present[string(b)] {
			present[string(b)] = true
			existing = append(existing, b)
		}
	}
	return existing
}

// removeBinarySetMembers filters the stored set down to the members the
// removal set does not carry. A binary member is its own identity.
func removeBinarySetMembers(existing, removals [][]byte) [][]byte {
	toDelete := make(map[string]bool, len(removals))
	for _, b := range removals {
		toDelete[string(b)] = true
	}
	remaining := make([][]byte, 0, len(existing))
	for _, b := range existing {
		if !toDelete[string(b)] {
			remaining = append(remaining, b)
		}
	}
	return remaining
}

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

// getReturnConsumedCapacity validates and reads the ReturnConsumedCapacity
// member: the Smithy enum trait admits exactly INDEXES, TOTAL and NONE, an
// unset member behaving as NONE. An unknown value is a ValidationException
// the caller must answer before executing anything.
func getReturnConsumedCapacity(params map[string]interface{}) (string, error) {
	v := request.GetStringParam(params, "ReturnConsumedCapacity")
	return v, validateReturnConsumedCapacityValue(v)
}

// validateReturnConsumedCapacityValue is getReturnConsumedCapacity's form
// for a core holding the extracted member: it admits the enum's values
// plus the absent empty string and rejects everything else.
func validateReturnConsumedCapacityValue(v string) error {
	switch v {
	case "", "NONE", "TOTAL", "INDEXES":
		return nil
	default:
		return NewAPIError("com.amazon.coral.validate#ValidationException",
			fmt.Sprintf("Invalid ReturnConsumedCapacity value %q: must be one of INDEXES, TOTAL, NONE", v), http.StatusBadRequest)
	}
}

// itemReadUnits returns the read capacity one item read consumes: the
// item's size rounded up to 4 KB multiples, halved for an eventually
// consistent read ("One read request unit represents one strongly
// consistent read operation per second, or two eventually consistent read
// operations per second, for an item up to 4 KB in size"). The size basis
// is what the read evaluated — the full stored item, before any
// projection narrows the returned attributes.
func itemReadUnits(itemSizeBytes int64, consistentRead bool) float64 {
	units := (itemSizeBytes + dbstore.ReadCapacityUnitBytes - 1) / dbstore.ReadCapacityUnitBytes
	if units < 1 {
		units = 1
	}
	if consistentRead {
		return float64(units)
	}
	return float64(units) / 2
}

// itemWriteUnits returns the write capacity one item write consumes: the
// written item's size rounded up to 1 KB multiples ("One write request
// unit represents one write operation per second, for an item up to 1 KB
// in size"). A write of an item that does not exist (a delete of a missing
// key) still consumes the one-unit minimum.
func itemWriteUnits(itemSizeBytes int64) float64 {
	units := (itemSizeBytes + dbstore.WriteCapacityUnitBytes - 1) / dbstore.WriteCapacityUnitBytes
	if units < 1 {
		units = 1
	}
	return float64(units)
}

// getItemCollectionMetricsSetting validates and reads the
// ReturnItemCollectionMetrics member: the Smithy enum trait admits exactly
// SIZE and NONE, an unset member behaving as NONE.
func getItemCollectionMetricsSetting(params map[string]interface{}) (string, error) {
	v := request.GetStringParam(params, "ReturnItemCollectionMetrics")
	return v, validateItemCollectionMetricsValue(v)
}

// validateItemCollectionMetricsValue is getItemCollectionMetricsSetting's
// form for a core holding the extracted member.
func validateItemCollectionMetricsValue(v string) error {
	switch v {
	case "", "NONE", "SIZE":
		return nil
	default:
		return NewAPIError("com.amazon.coral.validate#ValidationException",
			fmt.Sprintf("Invalid ReturnItemCollectionMetrics value %q: must be one of SIZE, NONE", v), http.StatusBadRequest)
	}
}

func buildConsumedCapacityResponse(tableName string, capacityUnits float64) map[string]interface{} {
	return map[string]interface{}{
		"TableName":     tableName,
		"CapacityUnits": capacityUnits,
	}
}

// buildReadConsumedCapacityResponse renders a read operation's
// ConsumedCapacity entry: the ConsumedCapacity shape defines both the
// aggregate member and the read-specific member, and read operations
// populate them with the same figure. Under ReturnConsumedCapacity=INDEXES
// the per-table detail joins the aggregate.
func buildReadConsumedCapacityResponse(tableName string, capacityUnits float64, includeTableDetail bool) map[string]interface{} {
	resp := map[string]interface{}{
		"TableName":         tableName,
		"CapacityUnits":     capacityUnits,
		"ReadCapacityUnits": capacityUnits,
	}
	if includeTableDetail {
		resp["Table"] = map[string]interface{}{
			"CapacityUnits":     capacityUnits,
			"ReadCapacityUnits": capacityUnits,
		}
	}
	return resp
}

// buildConsumedCapacityResponseWithVector adds the per-index vector write
// bytes to a ConsumedCapacity response when the write indexed any vectors.
func buildConsumedCapacityResponseWithVector(tableName string, capacityUnits float64, vectorIndexes map[string]interface{}) map[string]interface{} {
	resp := buildConsumedCapacityResponse(tableName, capacityUnits)
	if vectorIndexes != nil {
		resp["VectorIndexes"] = vectorIndexes
	}
	return resp
}

// vectorWriteCapacityForItems builds the ConsumedCapacity.VectorIndexes map
// for a write that modifies attributes indexed by a vector index: per index,
// the serialised byte length of the vector attribute values the write puts
// into the index (AWS documents the member but not its formula). Items whose
// vector attribute is absent, malformed, or dimension-mismatched are not
// indexed and report nothing. Returns nil when the table has no vector
// indexes or none of the items carries an indexable vector.
func vectorWriteCapacityForItems(table *dbstore.Table, items ...*dbstore.Item) map[string]interface{} {
	if table == nil || len(table.VectorIndexes) == 0 {
		return nil
	}
	var indexes map[string]interface{}
	for _, vi := range table.VectorIndexes {
		var bytes int
		for _, item := range items {
			if item == nil {
				continue
			}
			var attr *dbstore.AttributeValue
			if item.Attributes != nil {
				attr = item.Attributes[vi.VectorAttributeName]
			}
			if attr == nil && item.Key != nil {
				attr = item.Key[vi.VectorAttributeName]
			}
			if attr == nil || attr.L == nil || int64(len(attr.L)) != vi.Dimensions {
				continue
			}
			if encoded, err := json.Marshal(buildItemResponse(map[string]*dbstore.AttributeValue{"v": attr})); err == nil {
				bytes += len(encoded)
			}
		}
		if bytes > 0 {
			if indexes == nil {
				indexes = make(map[string]interface{})
			}
			indexes[vi.IndexName] = map[string]interface{}{
				"VectorWriteRequestBytes": float64(bytes),
			}
		}
	}
	return indexes
}

func buildConsumedCapacityResponseWithIndex(tableName string, indexName string, capacityUnits float64, isLSI bool) map[string]interface{} {
	// The per-table detail belongs to every INDEXES response — a base-table
	// read reports the table alone; an index read adds the index entry.
	resp := map[string]interface{}{
		"TableName":     tableName,
		"CapacityUnits": capacityUnits,
		"Table":         map[string]interface{}{"CapacityUnits": capacityUnits},
	}
	if indexName != "" {
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

// isConditionFunctionName returns true for DynamoDB expression function
// names that should absorb their parenthetical arguments as a single token.
// This prevents keywords like IN from being treated as function calls
// when written as IN(:v1, :v2) without a space before the parenthesis.
// Function names are case-sensitive (documented on the comparison
// operators and functions page), so the match is exact — an uppercase
// form does not glue and the downstream grammar rejects it.
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

// validateKeySchemaMembership enforces that an addressing Key carries
// exactly the table's primary-key attributes: a Key naming an attribute
// outside the key schema — or missing one — is answered with the
// ValidationException "The provided key element does not match the schema".
// The rule binds only request Key members; an ExclusiveStartKey for an index
// read legitimately carries the index key attributes alongside the primary
// key, and a PutItem Item naturally holds non-key attributes, so neither
// path may call this.
func validateKeySchemaMembership(table *dbstore.Table, key map[string]*dbstore.AttributeValue) error {
	if len(key) != len(table.KeySchema) {
		return keySchemaMismatchError()
	}
	for name := range key {
		if !isKeyAttribute(table, name) {
			return keySchemaMismatchError()
		}
	}
	return nil
}

// keySchemaMismatchError is the ValidationException DynamoDB answers a Key
// whose member set differs from the table's key schema with.
func keySchemaMismatchError() error {
	return NewAPIError("com.amazon.coral.validate#ValidationException",
		"The provided key element does not match the schema", http.StatusBadRequest)
}

// validateKeyTypes checks a supplied primary key against the table's key
// schema and attribute definitions, then enforces the documented length
// limits on string-typed key values (partition 2048 bytes, sort 1024 —
// the constraint is on the attribute value itself, so every plane that
// addresses an item by primary key applies it).
func validateKeyTypes(table *dbstore.Table, key map[string]*dbstore.AttributeValue) error {
	if err := validateKeySchemaAttrTypes(table, table.KeySchema, key); err != nil {
		return err
	}
	for _, ks := range table.KeySchema {
		attr, ok := key[ks.AttributeName]
		if !ok || attr == nil || attr.S == nil {
			continue
		}
		if ks.KeyType == dbstore.KeyTypeHash && len(*attr.S) > dbstore.MaxPartitionKeyBytes {
			// The missing space before the byte count is the service's own
			// message form.
			return NewAPIError("com.amazon.coral.validate#ValidationException",
				fmt.Sprintf("One or more parameter values were invalid: Size of hashkey has exceeded the maximum size limit of%d bytes", dbstore.MaxPartitionKeyBytes),
				http.StatusBadRequest)
		}
		if ks.KeyType == dbstore.KeyTypeRange && len(*attr.S) > dbstore.MaxSortKeyBytes {
			return NewAPIError("com.amazon.coral.validate#ValidationException",
				fmt.Sprintf("One or more parameter values were invalid: Aggregated size of all range keys has exceeded the size limit of %d bytes", dbstore.MaxSortKeyBytes),
				http.StatusBadRequest)
		}
	}
	return nil
}

// validateIndexStartKeyTypes checks the index-key attributes an index
// read's ExclusiveStartKey carries against the queried index's schema:
// the service's own LastEvaluatedKey emits the index keys alongside the
// primary key, and a wrong-typed echoed value would encode to a marker
// position its own type dictates rather than the named item's. Absent
// attributes are the marker composer's own presence contract; only
// present values are type-checked here.
func validateIndexStartKeyTypes(table *dbstore.Table, indexName string, esk map[string]*dbstore.AttributeValue) error {
	hashName, sortName, _ := indexKeyAttributeNames(table, indexName)
	defs := attributeTypeDefinitions(table)
	for _, name := range []string{hashName, sortName} {
		if name == "" {
			continue
		}
		val, ok := esk[name]
		if !ok || val == nil {
			continue
		}
		expected, hasDef := defs[name]
		if !hasDef {
			continue
		}
		actual := attributeValueType(val)
		if actual == "" {
			return keyTypeMismatchError(name, expected, expected)
		}
		if actual != string(expected) {
			return keyTypeMismatchError(name, expected, dbstore.ScalarAttributeType(actual))
		}
	}
	return nil
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

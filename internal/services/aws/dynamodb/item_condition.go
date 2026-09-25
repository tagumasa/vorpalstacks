package dynamodb

import (
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// The condition grammar itself — parse, validation and evaluation — lives
// in condition_ast.go; this file holds the document-path resolution the
// update plane shares with it, and the attribute predicates those engines
// and the response builders use.

// getNestedAttributeValue resolves a document path (e.g. "a.b.c" or
// "a[0].b") against an item's attributes, navigating Map and List
// AttributeValue types, with expression attribute names resolved per path
// segment — the grammar defines one alias per element (#pr.#1star), so the
// raw token is parsed first and each "#alias" segment is substituted on
// the parsed parts. An alias the names map does not define is a request
// error, the rejection every expression plane applies; a segment whose
// attribute is missing resolves to nil — a missing attribute, unlike an
// undefined alias, is a legal read. A malformed path — an invalid bracket
// index, or a path that does not begin with an attribute name — is a
// request error, not a missing attribute.
func getNestedAttributeValue(attrs map[string]*dbstore.AttributeValue, path string, names map[string]string) (*dbstore.AttributeValue, error) {
	if attrs == nil || path == "" {
		return nil, nil
	}
	parts, err := resolveDocPathParts(path, names)
	if err != nil {
		return nil, err
	}
	return getDocPathValue(attrs, parts), nil
}

func getHashKeyName(table *dbstore.Table) string {
	for _, ks := range table.KeySchema {
		if ks.KeyType == dbstore.KeyTypeHash {
			return ks.AttributeName
		}
	}
	return ""
}

func getAttributeTypeName(attr *dbstore.AttributeValue) string {
	switch {
	case attr.S != nil:
		return "S"
	case attr.SS != nil:
		return "SS"
	case attr.N != nil:
		return "N"
	case attr.NS != nil:
		return "NS"
	case attr.B != nil:
		return "B"
	case attr.BS != nil:
		return "BS"
	case attr.BOOL != nil:
		return "BOOL"
	case attr.NULL != nil:
		return "NULL"
	case attr.L != nil:
		return "L"
	case attr.M != nil:
		return "M"
	}
	return ""
}

func computeAttributeSize(attr *dbstore.AttributeValue) (int, bool) {
	switch {
	case attr.S != nil:
		return len(*attr.S), true
	case attr.B != nil:
		return len(attr.B), true
	case attr.SS != nil:
		return len(attr.SS), true
	case attr.NS != nil:
		return len(attr.NS), true
	case attr.BS != nil:
		return len(attr.BS), true
	case attr.L != nil:
		return len(attr.L), true
	case attr.M != nil:
		return len(attr.M), true
	}
	return 0, false
}

func listContainsValue(list []*dbstore.AttributeValue, target *dbstore.AttributeValue) bool {
	for _, elem := range list {
		if attributeValuesEqual(elem, target) {
			return true
		}
	}
	return false
}

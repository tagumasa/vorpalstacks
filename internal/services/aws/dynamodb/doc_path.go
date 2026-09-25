package dynamodb

import (
	"net/http"
	"strings"
	"unicode/utf8"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// docPathPart is one segment of a parsed document path: either an attribute
// name within a Map or an index within a List. A document path always begins
// with a name part; index parts only ever follow one.
type docPathPart struct {
	name    string
	index   int
	isIndex bool
}

// docPathPartsEqual answers whether two parsed document paths address the
// same value.
func docPathPartsEqual(a, b []docPathPart) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// parseDocPath is the single document-path grammar shared by condition
// evaluation, projection assembly and update mutation. It splits "a.b[0].c"
// into name and bracket-index segments, validating every bracket index
// through validateBracketIndex and rejecting a path that does not begin
// with an attribute name.
func parseDocPath(path string) ([]docPathPart, error) {
	var parts []docPathPart
	current := ""
	i := 0

	for i < len(path) {
		c := path[i]
		switch c {
		case '.':
			if current != "" {
				parts = append(parts, docPathPart{name: current})
				current = ""
			}
			i++
		case '[':
			if current != "" {
				parts = append(parts, docPathPart{name: current})
				current = ""
			}
			j := i + 1
			for j < len(path) && path[j] != ']' {
				j++
			}
			// The list dereference operator is the closed pair "[n]": a
			// bracket left unterminated until the end of the path is not a
			// document path the expression grammar produces.
			if j == len(path) {
				return nil, ErrInvalidParameter
			}
			idx, idxErr := validateBracketIndex(path[i+1 : j])
			if idxErr != nil {
				return nil, idxErr
			}
			parts = append(parts, docPathPart{index: idx, isIndex: true})
			i = j + 1
		default:
			// A segment name carries whole runes: taking a byte at a time
			// would re-encode each byte's value as a rune and mangle a
			// multi-byte name.
			_, size := utf8.DecodeRuneInString(path[i:])
			current += path[i : i+size]
			i += size
		}
	}

	if current != "" {
		parts = append(parts, docPathPart{name: current})
	}

	// A document path addresses a top-level attribute first; a leading
	// bracket index has no attribute to index into.
	if len(parts) > 0 && parts[0].isIndex {
		return nil, ErrInvalidParameter
	}

	return parts, nil
}

// resolveDocPathParts parses a raw document-path token and resolves every
// name segment that is an expression attribute name. The documented
// grammar defines one alias per path element (#pr.#1star), so a compound
// token resolves segment by segment after the path grammar has split it.
// Resolution happens on the parsed segments, never on a re-serialised
// string: an alias standing for a name that itself contains '.' or '['
// stays a single segment. An alias absent from the names map is a
// validation error.
func resolveDocPathParts(path string, names map[string]string) ([]docPathPart, error) {
	return resolveDocPathSegments(path, names, true)
}

// resolveDocPathPartsLenient is the analysis-pass variant: an alias absent
// from the names map stays as the literal segment name, so an analysis
// that runs before an expression's application — the key-attribute
// guard's top-level segment extraction — never pre-empts the apply pass's
// own rejection of the undefined alias.
func resolveDocPathPartsLenient(path string, names map[string]string) ([]docPathPart, error) {
	return resolveDocPathSegments(path, names, false)
}

func resolveDocPathSegments(path string, names map[string]string, strict bool) ([]docPathPart, error) {
	parts, err := parseDocPath(path)
	if err != nil {
		return nil, err
	}
	for i := range parts {
		if parts[i].isIndex || !strings.HasPrefix(parts[i].name, "#") {
			continue
		}
		resolved, ok := names[parts[i].name]
		if !ok {
			if !strict {
				continue
			}
			return nil, ErrInvalidParameter
		}
		parts[i].name = resolved
	}
	return parts, nil
}

// getDocPathValue resolves parsed document-path parts against an item's
// attributes, navigating Map and List AttributeValue types. Returns nil when
// any segment is missing or a container's type does not match the part kind.
func getDocPathValue(attrs map[string]*dbstore.AttributeValue, parts []docPathPart) *dbstore.AttributeValue {
	if attrs == nil || len(parts) == 0 {
		return nil
	}
	current, ok := attrs[parts[0].name]
	if !ok || current == nil {
		return nil
	}
	for _, part := range parts[1:] {
		if part.isIndex {
			if current.L == nil || part.index >= len(current.L) {
				return nil
			}
			current = current.L[part.index]
		} else {
			if current.M == nil {
				return nil
			}
			v, exists := current.M[part.name]
			if !exists || v == nil {
				return nil
			}
			current = v
		}
	}
	return current
}

// errInvalidUpdateDocumentPath is the ValidationException the service
// answers a SET whose document path cannot be traversed with. The guide:
// "You cannot update nested map attributes if the parent map does not
// exist. If you attempt to update a nested attribute ... when the parent
// map ... does not exist, DynamoDB returns a ValidationException with the
// message 'The document path provided in the update expression is invalid
// for update.'" The same path-validity rule covers a parent that exists
// but is not the container type the path descends through: items meant to
// carry nested updates later initialise their parent maps empty.
func errInvalidUpdateDocumentPath() error {
	return NewAPIError("com.amazon.coral.validate#ValidationException",
		"The document path provided in the update expression is invalid for update", http.StatusBadRequest)
}

func setNestedValue(attrs map[string]*dbstore.AttributeValue, parts []docPathPart, value *dbstore.AttributeValue) error {
	if len(parts) == 0 {
		return nil
	}

	if len(parts) == 1 {
		attrs[parts[0].name] = value
		return nil
	}

	// A nested SET traverses existing containers only — no implicit parent
	// creation at any depth.
	current, exists := attrs[parts[0].name]
	if !exists || current == nil {
		return errInvalidUpdateDocumentPath()
	}

	return setNestedValueRecursive(current, parts[1:], value)
}

func setNestedValueRecursive(current *dbstore.AttributeValue, parts []docPathPart, value *dbstore.AttributeValue) error {
	if len(parts) == 0 {
		return nil
	}

	part := parts[0]
	isLast := len(parts) == 1

	if part.isIndex {
		if current.L == nil {
			return errInvalidUpdateDocumentPath()
		}
		if isLast {
			// The documented SET rule for lists: an element the path names
			// that does not already exist appends at the end of the list —
			// an out-of-range index never positions the write and never
			// pads the gap.
			if part.index >= len(current.L) {
				current.L = append(current.L, value)
				return nil
			}
			current.L[part.index] = value
			return nil
		}
		if part.index >= len(current.L) {
			return errInvalidUpdateDocumentPath()
		}
		next := current.L[part.index]
		if next == nil {
			return errInvalidUpdateDocumentPath()
		}
		return setNestedValueRecursive(next, parts[1:], value)
	}

	if current.M == nil {
		return errInvalidUpdateDocumentPath()
	}
	if isLast {
		current.M[part.name] = value
		return nil
	}
	next, ok := current.M[part.name]
	if !ok || next == nil {
		return errInvalidUpdateDocumentPath()
	}
	return setNestedValueRecursive(next, parts[1:], value)
}

func removeNestedValue(attrs map[string]*dbstore.AttributeValue, parts []docPathPart) error {
	if len(parts) == 0 {
		return nil
	}

	if len(parts) == 1 {
		delete(attrs, parts[0].name)
		return nil
	}

	current, exists := attrs[parts[0].name]
	if !exists {
		return nil
	}

	removeNestedValueRecursive(current, parts[1:])
	return nil
}

func removeNestedValueRecursive(current *dbstore.AttributeValue, parts []docPathPart) {
	if len(parts) == 0 {
		return
	}

	part := parts[0]
	isLast := len(parts) == 1

	if part.isIndex {
		if current.L == nil || part.index >= len(current.L) {
			return
		}
		if isLast {
			newList := make([]*dbstore.AttributeValue, 0, len(current.L)-1)
			for i, v := range current.L {
				if i != part.index {
					newList = append(newList, v)
				}
			}
			current.L = newList
		} else {
			removeNestedValueRecursive(current.L[part.index], parts[1:])
		}
	} else {
		if current.M == nil {
			return
		}
		if isLast {
			delete(current.M, part.name)
		} else {
			if next, ok := current.M[part.name]; ok {
				removeNestedValueRecursive(next, parts[1:])
			}
		}
	}
}

package dynamodb

import (
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
			current += string(c)
			i++
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

func setNestedValue(attrs map[string]*dbstore.AttributeValue, path string, value *dbstore.AttributeValue) error {
	parts, err := parseDocPath(path)
	if err != nil {
		return err
	}
	if len(parts) == 0 {
		return nil
	}

	if len(parts) == 1 {
		attrs[parts[0].name] = value
		return nil
	}

	current, exists := attrs[parts[0].name]
	nextPartIsIndex := parts[1].isIndex
	if !exists {
		if nextPartIsIndex {
			return ErrInvalidParameter
		}
		current = dbstore.MapValue(make(map[string]*dbstore.AttributeValue))
		attrs[parts[0].name] = current
	} else if nextPartIsIndex && current.L == nil {
		return ErrInvalidParameter
	} else if !nextPartIsIndex && current.M == nil {
		current = dbstore.MapValue(make(map[string]*dbstore.AttributeValue))
		attrs[parts[0].name] = current
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
		if current.L == nil || part.index >= len(current.L) {
			return ErrInvalidParameter
		}
		if isLast {
			current.L[part.index] = value
			return nil
		} else {
			return setNestedValueRecursive(current.L[part.index], parts[1:], value)
		}
	} else {
		if current.M == nil {
			if current.S != nil || current.N != nil || current.B != nil ||
				current.BOOL != nil || current.NULL != nil || current.L != nil ||
				current.SS != nil || current.NS != nil || current.BS != nil {
				return ErrInvalidParameter
			}
			current.M = make(map[string]*dbstore.AttributeValue)
		}
		if isLast {
			current.M[part.name] = value
			return nil
		} else {
			if _, ok := current.M[part.name]; !ok {
				current.M[part.name] = dbstore.MapValue(make(map[string]*dbstore.AttributeValue))
			}
			return setNestedValueRecursive(current.M[part.name], parts[1:], value)
		}
	}
}

func removeNestedValue(attrs map[string]*dbstore.AttributeValue, path string) error {
	parts, err := parseDocPath(path)
	if err != nil {
		return err
	}
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

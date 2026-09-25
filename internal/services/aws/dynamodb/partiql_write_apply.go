package dynamodb

import (
	"strings"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
	"vorpalstacks/pkg/sqlparser"
)

// This file carries the write-clause appliers of the PartiQL write plane:
// the SET/ADD/DELETE/REMOVE assignment application over a target item, the
// operand model that materialises SET operands (literals, paths, and the
// documented functions), and the clause-target naming the engines and the
// validation share. The write engines that drive these appliers live in
// partiql_operations_core.go.

// applyAddAssignments applies the ADD clauses of a PartiQL UPDATE: each
// assignment's operand materialises through the shared value materialiser
// and folds into its target path through the ADD action semantics, the
// stored base read from the pre-statement image preUpdate — the
// documented multi-action evaluation basis. The update plane answers a
// type-incompatible pair with ErrTypeMismatch — the shared appliers'
// contract, identical on every face.
func applyAddAssignments(attrs map[string]*dbstore.AttributeValue, preUpdate map[string]*dbstore.AttributeValue, assignments []setAssignment, params *partiQLParams) error {
	for _, asgn := range assignments {
		parts, pErr := parseDocPath(asgn.attrName)
		if pErr != nil {
			return pErr
		}
		addValue, err := exprToAttributeValueWithParams(asgn.value, params)
		if err != nil {
			return err
		}
		if err := applyActionAtPath(attrs, preUpdate, parts, true, addValue); err != nil {
			return err
		}
	}
	return nil
}

// applyDeleteAssignments applies the DELETE clauses of a PartiQL UPDATE
// through the DELETE action semantics at each target path, the stored
// base read from preUpdate exactly as applyAddAssignments reads it; the
// mismatch error contract follows applyAddAssignments.
func applyDeleteAssignments(attrs map[string]*dbstore.AttributeValue, preUpdate map[string]*dbstore.AttributeValue, assignments []setAssignment, params *partiQLParams) error {
	for _, asgn := range assignments {
		parts, pErr := parseDocPath(asgn.attrName)
		if pErr != nil {
			return pErr
		}
		delValue, err := exprToAttributeValueWithParams(asgn.value, params)
		if err != nil {
			return err
		}
		if err := applyActionAtPath(attrs, preUpdate, parts, false, delValue); err != nil {
			return err
		}
	}
	return nil
}

// applySetAssignments writes SET clause assignments into attrs through
// the shared update-operand model: the SET functions and every path
// operand evaluate against the pre-statement image preUpdate — the
// documented multi-action basis, so a later assignment never observes an
// earlier one's write — and literals or bound parameters materialise at
// compile time, an unresolvable placeholder failing the statement. Each
// assignment's left-hand side parses once here into its document path
// (both parser paths hand the rendered path over), so a nested target
// writes the member it names under the no-implicit-parent-creation rule.
func applySetAssignments(attrs map[string]*dbstore.AttributeValue, preUpdate map[string]*dbstore.AttributeValue, assignments []setAssignment, params *partiQLParams) error {
	for _, asgn := range assignments {
		parts, pErr := parseDocPath(asgn.attrName)
		if pErr != nil {
			return pErr
		}
		operand, err := partiqlUpdateOperand(asgn.value, params)
		if err != nil {
			return err
		}
		if fn, ok := operand.(updateSetFunction); ok {
			// The only documented form of the SET functions addresses the
			// assignment's own target: set_add(path, value) with the path
			// naming the target the fold writes. A first argument naming
			// any other path fits the documented form nowhere, so the
			// statement fails instead of folding into the target while
			// the named path is never read.
			if fn.path != asgn.attrName {
				return ErrInvalidParameter
			}
			value, vErr := fn.operand.eval(preUpdate)
			if vErr != nil {
				return vErr
			}
			if value == nil {
				return ErrInvalidParameter
			}
			if err := applyActionAtPath(attrs, preUpdate, parts, fn.add, value); err != nil {
				return err
			}
			continue
		}
		attrValue, evalErr := operand.eval(preUpdate)
		if evalErr != nil {
			return evalErr
		}
		if attrValue == nil {
			// Every SET clause operand must resolve to a value: a document
			// path naming an attribute the item does not hold is a
			// validation error.
			return ErrInvalidParameter
		}
		if err := setNestedValue(attrs, parts, attrValue); err != nil {
			return err
		}
	}
	return nil
}

// applyActionAtPath folds a set action's operand into the value the target
// path holds, through the ADD and DELETE action semantics: a top-level
// path uses the action appliers themselves, a nested path the same
// combination over the nested read and write walks. The stored base reads
// the pre-statement image preUpdate while the write lands on attrs. The
// absent-target rules are the actions' own — an absent target creates
// from the operand for the add and is a no-op for the delete — and an
// absent parent is the invalid document path the nested walks answer with.
func applyActionAtPath(attrs map[string]*dbstore.AttributeValue, preUpdate map[string]*dbstore.AttributeValue, parts []docPathPart, add bool, value *dbstore.AttributeValue) error {
	if len(parts) == 1 {
		if add {
			_, err := applyAddAction(attrs, preUpdate, parts[0].name, value)
			return err
		}
		_, err := applyDeleteAction(attrs, preUpdate, parts[0].name, value)
		return err
	}
	existing := getDocPathValue(preUpdate, parts)
	if add {
		if existing == nil {
			return setNestedValue(attrs, parts, value)
		}
		combined, cErr := combineAddValues(existing, value)
		if cErr != nil {
			return cErr
		}
		// A nil combination with no error is the empty-set no-change.
		if combined == nil {
			return nil
		}
		return setNestedValue(attrs, parts, combined)
	}
	if existing == nil {
		return nil
	}
	remaining, emptied, cErr := combineDeleteValues(existing, value)
	if cErr != nil {
		return cErr
	}
	if emptied {
		return removeNestedValue(attrs, parts)
	}
	if remaining == nil {
		return nil
	}
	return setNestedValue(attrs, parts, remaining)
}

// partiqlUpdateOperand compiles one PartiQL SET-clause value into the
// shared update-operand model: the two SET functions map onto their nodes,
// a column reference reads as the document path it names (a qualified
// reference addresses a nested path), and every other expression
// materialises through the shared value materialiser as a literal or a
// bound parameter.
func partiqlUpdateOperand(expr sqlparser.Expr, params *partiQLParams) (updateOperand, error) {
	if funcExpr, ok := expr.(*sqlparser.FuncExpr); ok {
		switch {
		case strings.EqualFold(funcExpr.Name.String(), "if_not_exists") && len(funcExpr.Exprs) >= 2:
			path, isPath := partiqlFuncPathArg(funcExpr.Exprs[0])
			if !isPath {
				// The first argument does not name a path, so the guard
				// cannot run and the fallback decides.
				fallbackExpr, argOk := funcExprArg(funcExpr.Exprs[1])
				if !argOk {
					return nil, ErrInvalidParameter
				}
				return partiqlUpdateOperand(fallbackExpr, params)
			}
			fallbackExpr, argOk := funcExprArg(funcExpr.Exprs[1])
			if !argOk {
				return nil, ErrInvalidParameter
			}
			fallback, err := partiqlUpdateOperand(fallbackExpr, params)
			if err != nil {
				return nil, err
			}
			return updateIfNotExists{path: updatePath{path: path}, fallback: fallback}, nil
		case strings.EqualFold(funcExpr.Name.String(), "list_append") && len(funcExpr.Exprs) == 2:
			firstExpr, firstOk := funcExprArg(funcExpr.Exprs[0])
			secondExpr, secondOk := funcExprArg(funcExpr.Exprs[1])
			if !firstOk || !secondOk {
				return nil, ErrInvalidParameter
			}
			first, err := partiqlUpdateOperand(firstExpr, params)
			if err != nil {
				return nil, err
			}
			second, err := partiqlUpdateOperand(secondExpr, params)
			if err != nil {
				return nil, err
			}
			return updateListAppend{first: first, second: second}, nil
		case strings.EqualFold(funcExpr.Name.String(), "set_add") && len(funcExpr.Exprs) == 2:
			// The documented form is set_add(path, value) with the path
			// addressing the assignment's own target; the SET applier
			// enforces the correspondence and folds the value through the
			// ADD action.
			path, isPath := partiqlFuncPathArg(funcExpr.Exprs[0])
			if !isPath {
				return nil, ErrInvalidParameter
			}
			valueExpr, argOk := funcExprArg(funcExpr.Exprs[1])
			if !argOk {
				return nil, ErrInvalidParameter
			}
			inner, err := partiqlUpdateOperand(valueExpr, params)
			if err != nil {
				return nil, err
			}
			return updateSetFunction{add: true, path: path, operand: inner}, nil
		case strings.EqualFold(funcExpr.Name.String(), "set_delete") && len(funcExpr.Exprs) == 2:
			path, isPath := partiqlFuncPathArg(funcExpr.Exprs[0])
			if !isPath {
				return nil, ErrInvalidParameter
			}
			valueExpr, argOk := funcExprArg(funcExpr.Exprs[1])
			if !argOk {
				return nil, ErrInvalidParameter
			}
			inner, err := partiqlUpdateOperand(valueExpr, params)
			if err != nil {
				return nil, err
			}
			return updateSetFunction{add: false, path: path, operand: inner}, nil
		}
	}
	if col, ok := expr.(*sqlparser.ColName); ok {
		return updatePath{path: colNameDocumentPath(col)}, nil
	}
	v, err := exprToAttributeValueWithParams(expr, params)
	if err != nil {
		return nil, err
	}
	return updateValue{value: v}, nil
}

// funcExprArg unwraps one function argument's AliasedExpr shell; the
// function-argument positions the operand compiler reads carry no other
// SelectExpr shape, and a different shape fails the statement.
func funcExprArg(arg sqlparser.SelectExpr) (sqlparser.Expr, bool) {
	if aliased, ok := arg.(*sqlparser.AliasedExpr); ok {
		return aliased.Expr, true
	}
	return nil, false
}

// partiqlFuncPathArg renders a function argument as the document path it
// names, reporting whether the argument is a column reference at all.
func partiqlFuncPathArg(arg sqlparser.SelectExpr) (string, bool) {
	expr, ok := funcExprArg(arg)
	if !ok {
		return "", false
	}
	if col, ok := expr.(*sqlparser.ColName); ok {
		return colNameDocumentPath(col), true
	}
	return "", false
}

// colNameDocumentPath renders a column reference as the document path it
// names: a bare name addresses a top-level attribute, a qualified name
// the nested path underneath its qualifiers. The grammar nests one
// qualifier level — a three-segment reference is a single ColName — so
// every segment renders: dropping the outermost qualifier would name the
// reference's two-segment shadow instead of the path it addresses. Every
// segment renders its raw value, never the identifier printer's escaped
// form — clause targets, operands, and the set-function correspondence all
// render through this walk, so an escapable name (a keyword, non-ASCII)
// stays one and the same path on every face.
func colNameDocumentPath(col *sqlparser.ColName) string {
	segments := make([]string, 0, 3)
	if lead := col.Qualifier.Qualifier.String(); lead != "" {
		segments = append(segments, lead)
	}
	if mid := col.Qualifier.Name.String(); mid != "" {
		segments = append(segments, mid)
	}
	segments = append(segments, col.Name.String())
	return strings.Join(segments, ".")
}

// applyRemoveAttrs removes the attributes the REMOVE clause names, each
// through its document path — a nested or indexed target removes the
// member it names, and a missing path is the removal no-op the walk
// defines.
func applyRemoveAttrs(attrs map[string]*dbstore.AttributeValue, removeAttrs []string) error {
	for _, name := range removeAttrs {
		parts, pErr := parseDocPath(name)
		if pErr != nil {
			return pErr
		}
		if err := removeNestedValue(attrs, parts); err != nil {
			return err
		}
	}
	return nil
}

// clauseTargetRoot renders the top-level attribute one clause target
// addresses: everything before the first path separator, so a nested
// write guards and reports on the attribute it roots at.
func clauseTargetRoot(name string) string {
	if i := strings.IndexAny(name, ".["); i >= 0 {
		return name[:i]
	}
	return name
}

// updateClauseTargetNames collects every top-level attribute name an UPDATE
// statement writes across its SET, REMOVE, ADD and DELETE clauses. The key
// attributes identify the item being updated, so any clause touching one of
// them is rejected before the statement applies.
func updateClauseTargetNames(clauses updateClauses) []string {
	names := make([]string, 0, len(clauses.setAssignments)+len(clauses.removeAttrs)+
		len(clauses.addAssignments)+len(clauses.deleteAssignments))
	for _, asgn := range clauses.setAssignments {
		names = append(names, clauseTargetRoot(asgn.attrName))
	}
	for _, name := range clauses.removeAttrs {
		names = append(names, clauseTargetRoot(name))
	}
	for _, asgn := range clauses.addAssignments {
		names = append(names, clauseTargetRoot(asgn.attrName))
	}
	for _, asgn := range clauses.deleteAssignments {
		names = append(names, clauseTargetRoot(asgn.attrName))
	}
	return names
}

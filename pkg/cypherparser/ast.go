// AST types for the openCypher parser.
//
// Based on goraphdb/cypher_ast.go core types (CypherQuery, CypherWrite, CypherMerge,
// Pattern, NodePattern, RelPattern, Expression), extended with:
//   - RelPattern.Props for edge properties (missing in goraphdb)
//   - WITH/UNWIND/DISTINCT/REMOVE clause types (not in goraphdb)
//   - Expression: IN, EXISTS, STARTS WITH, ENDS WITH, CONTAINS, CASE/WHEN,
//     IS NULL/NOT NULL, aggregation functions, list literals, arithmetic
//   - CypherPipeline for multi-segment queries (WITH ... MATCH ... RETURN)

package cypherparser

import "vorpalstacks/pkg/graphengine"

// ExplainMode controls whether the query is explained, profiled, or executed normally.
// Based on goraphdb ExplainMode.
type ExplainMode int

const (
	ExplainNone ExplainMode = iota
	ExplainOnly
	ExplainProfile
)

// CypherQuery is the top-level AST node for a read query.
// Based on goraphdb CypherQuery, extended with Distinct, multiple MATCH clauses,
// REMOVE clause, and SKIP support.
//
//	[EXPLAIN|PROFILE] [WITH items]
//	  [MATCH Pattern [WHERE Expr]]*
//	  [OPTIONAL MATCH Pattern [WHERE Expr]]
//	  [UNWIND expr AS var]
//	  [SET items]
//	  [DELETE vars]
//	  [REMOVE items]
//	  RETURN [DISTINCT] Items
//	  [ORDER BY Items] [SKIP n] [LIMIT n]
type CypherQuery struct {
	Explain       ExplainMode
	With          *WithClause
	Matches       []MatchClause
	Match         MatchClause
	Where         *Expression
	OptionalMatch *MatchClause
	OptionalWhere *Expression
	Unwind        *UnwindClause
	Set           []SetItem
	Delete        []string
	DeleteDetach  bool
	Remove        []RemoveItem
	Create        *CypherWrite
	Return        ReturnClause
	OrderBy       []OrderItem
	Skip          *int
	Limit         *int
	LimitParam    string
	SkipParam     string
}

// WithClause represents a WITH clause that projects intermediate results.
// Our addition — not in goraphdb.
//
//	WITH [DISTINCT] expr [AS alias], ...
type WithClause struct {
	Distinct bool
	Items    []WithItem
	OrderBy  []OrderItem
	Skip     *int
	Limit    *int
	Where    *Expression
}

// WithItem is a single expression in a WITH clause.
type WithItem struct {
	Expr  Expression
	Alias string
}

// UnwindClause represents an UNWIND clause that expands a list into rows.
// Our addition — not in goraphdb.
//
//	UNWIND expr AS var
type UnwindClause struct {
	Expr Expression
	Var  string
}

// MatchClause holds the pattern that follows a MATCH keyword.
// Based on goraphdb MatchClause.
type MatchClause struct {
	Pattern Pattern
}

// RemoveItem represents a single REMOVE operation.
// Our addition — not in goraphdb.
//
//	REMOVE n:Label          → {Variable: "n", Kind: RemoveLabel, Label: "Label"}
//	REMOVE n.prop           → {Variable: "n", Kind: RemoveProperty, Property: "prop"}
type RemoveItem struct {
	Variable string
	Kind     RemoveKind
	Label    string
	Property string
}

// RemoveKind identifies the kind of REMOVE operation.
type RemoveKind int

const (
	// RemoveLabel indicates a label removal (REMOVE n:Label).
	RemoveLabel RemoveKind = iota
	// RemoveProperty indicates a property removal (REMOVE n.prop).
	RemoveProperty
)

// Pattern is a sequence of node–rel–node–rel–…–node.
// Based on goraphdb Pattern.
type Pattern struct {
	Nodes []NodePattern
	Rels  []RelPattern
}

// NodePattern represents a single node in a MATCH or CREATE pattern.
// Based on goraphdb NodePattern.
//
//	(n)                → Variable="n", Labels=nil, Props=nil
//	(n:Person:Actor)   → Variable="n", Labels=["Person","Actor"], Props=nil
//	(n {name:"Alice"}) → Variable="n", Labels=nil, Props={"name":"Alice"}
type NodePattern struct {
	Variable string
	Labels   []string
	Props    map[string]any
}

// RelPattern represents a relationship (edge) in a pattern.
// Based on goraphdb RelPattern, extended with Props field for edge properties
// (goraphdb did not support inline edge properties).
//
//	-[:FOLLOWS]->       → Label="FOLLOWS", Dir=Outgoing
//	-[r:FOLLOWS {since: 2020}]->  → Variable="r", Props={"since": 2020}
//	-[:FOLLOWS*1..3]->  → VarLength=true, MinHops=1, MaxHops=3
//	<-[:KNOWS]-         → Dir=Incoming
//	-[:KNOWS]-          → Dir=Both
type RelPattern struct {
	Variable  string
	Label     string
	Dir       graphengine.Direction
	VarLength bool
	MinHops   int
	MaxHops   int
	Props     map[string]any
}

// SetItem represents a single property assignment in a SET clause.
// Based on goraphdb SetItem.
//
//	SET n.name = "Alice"   → {Variable:"n", Property:"name", Value:litExpr("Alice")}
//	SET n += {map}         → {Variable:"n", Value:mapExpr(...), Merge:true}
//	SET n = {map}          → {Variable:"n", Value:mapExpr(...), ReplaceAll:true}
//	SET n:Label            → {Variable:"n", Label:"Label", SetLabel:true}
type SetItem struct {
	Variable   string
	Property   string
	Value      Expression
	Merge      bool
	ReplaceAll bool
	SetLabel   bool
	Label      string
}

// ReturnClause holds the items after RETURN.
// Based on goraphdb ReturnClause, extended with Distinct flag.
type ReturnClause struct {
	Distinct bool
	Items    []ReturnItem
}

// ReturnItem is a single expression in the RETURN clause, optionally aliased.
// Based on goraphdb ReturnItem.
type ReturnItem struct {
	Expr  Expression
	Alias string
}

// OrderItem is a single expression in ORDER BY.
// Based on goraphdb OrderItem.
type OrderItem struct {
	Expr Expression
	Desc bool
}

// ---------------------------------------------------------------------------
// Expressions
// ---------------------------------------------------------------------------
// Based on goraphdb Expression types (Literal, VarRef, PropAccess, FuncCall,
// Comparison, And, Or, Not, Param), extended with:
//   - In, Exists, StartsWith, EndsWith, Contains, IsNull, IsNotNull
//   - CaseSwitch (CASE/WHEN/THEN/ELSE)
//   - ListLiteral, MapLiteral
//   - Aggregation (Count, Sum, Avg, Min, Max, Collect)
//   - ArithAdd, ArithSub, ArithMul, ArithDiv, ArithMod
// ---------------------------------------------------------------------------

// ExprKind identifies the type of a Cypher expression.
type ExprKind int

// Expression kind constants. The first group mirrors goraphdb; the remaining
// groups are our additions for extended openCypher support.
const (
	ExprLiteral ExprKind = iota
	ExprVarRef
	ExprPropAccess
	ExprFuncCall
	ExprComparison
	ExprAnd
	ExprOr
	ExprNot
	ExprParam

	// String predicates (our addition)
	ExprIn
	ExprExists
	ExprStartsWith
	ExprEndsWith
	ExprContains
	ExprIsNull
	ExprIsNotNull
	ExprRegexMatch

	// CASE expression (our addition)
	ExprCase

	// Literals (our addition)
	ExprListLiteral
	ExprMapLiteral

	// Aggregation functions (our addition)
	ExprAggregation

	// Arithmetic (our addition)
	ExprArithAdd
	ExprArithSub
	ExprArithMul
	ExprArithDiv
	ExprArithMod

	// Label/property projection (our addition)
	ExprLabels
	ExprProperties
)

// CompOp identifies a comparison operator.
type CompOp int

// Comparison operator constants.
const (
	OpEq CompOp = iota
	OpNeq
	OpLt
	OpGt
	OpLte
	OpGte
)

// AggFunc identifies an aggregation function.
type AggFunc int

// Aggregation function constants.
const (
	AggCount AggFunc = iota
	AggSum
	AggAvg
	AggMin
	AggMax
	AggCollect
)

// Expression is a polymorphic AST node for all expression types.
// Based on goraphdb Expression struct, extended with fields for each new ExprKind.
type Expression struct {
	Kind ExprKind

	// ExprLiteral
	LitValue any

	// ExprVarRef
	Variable string

	// ExprPropAccess
	Object           string
	Property         string
	NestedPropObject *Expression

	// ExprFuncCall
	FuncName string
	Args     []Expression

	// ExprComparison
	Left  *Expression
	Op    CompOp
	Right *Expression

	// ExprAnd / ExprOr
	Operands []Expression

	// ExprNot
	Inner *Expression

	// ExprParam
	ParamName string

	// ExprIn: left IN right (where right is a list or subquery)
	InLeft  *Expression
	InRight *Expression

	// ExprExists: EXISTS { pattern }
	ExistsPattern *Pattern

	// ExprStartsWith / ExprEndsWith / ExprContains
	StrExpr  *Expression
	StrValue *Expression

	// ExprIsNull / ExprIsNotNull
	IsNullExpr *Expression

	// ExprRegexMatch: expr =~ "pattern"
	RegexExpr    *Expression
	RegexPattern *Expression

	// ExprCase: CASE [expr] WHEN ... THEN ... ELSE ... END
	CaseOperand *Expression // non-nil for simple CASE (has operand)
	CaseWhens   []CaseWhen
	CaseElse    *Expression

	// ExprListLiteral
	ListItems []Expression

	// ExprMapLiteral
	MapPairs []MapPair

	// ExprAggregation
	AggFunc     AggFunc
	AggArg      *Expression // nil for COUNT(*)
	AggDistinct bool

	// ExprArith*
	ArithLeft  *Expression
	ArithRight *Expression

	// ExprLabels / ExprProperties
	LabelsExpr     *Expression
	PropertiesExpr *Expression
}

// CaseWhen represents a WHEN ... THEN ... branch in a CASE expression.
type CaseWhen struct {
	Condition *Expression
	Result    *Expression
}

// MapPair is a key-value pair in a map literal.
type MapPair struct {
	Key   string
	Value Expression
}

// paramRef is a sentinel type used in property maps to represent a $param reference.
// Based on goraphdb paramRef.
type paramRef string

// ---------------------------------------------------------------------------
// Write queries
// ---------------------------------------------------------------------------

// CypherWrite is the top-level AST node for a write query (CREATE).
// Based on goraphdb CypherWrite, extended with DISTINCT on RETURN.
//
//	CREATE pattern (',' pattern)* [RETURN [DISTINCT] items]
type CypherWrite struct {
	Creates []CreatePattern
	Return  *ReturnClause
	OrderBy []OrderItem
	Skip    *int
	Limit   *int
}

// CreatePattern is a chain of alternating nodes and relationships to create.
// Based on goraphdb CreatePattern.
type CreatePattern struct {
	Nodes []NodePattern
	Rels  []RelPattern
}

// CypherMerge is the top-level AST node for a MERGE query.
// Based on goraphdb CypherMerge.
//
//	MERGE (n:Label {key: value})
//	  [ON CREATE SET n.prop = val, ...]
//	  [ON MATCH SET n.prop = val, ...]
//	  [RETURN ...]
type CypherMerge struct {
	Pattern     MergePattern
	OnCreateSet []SetItem
	OnMatchSet  []SetItem
	Return      *ReturnClause
	OrderBy     []OrderItem
	Skip        *int
	Limit       *int
}

// MergePattern describes a single node to be merged (matched or created).
// Based on goraphdb MergePattern.
type MergePattern struct {
	Variable string
	Labels   []string
	Props    map[string]any
}

// CypherPipeline represents a multi-segment query with WITH clauses.
// Our addition — not in goraphdb.
//
//	MATCH (n) WITH n.name AS name WHERE name STARTS WITH "A" RETURN name
type CypherPipeline struct {
	Segments []PipelineSegment
}

// PipelineSegment is one segment of a pipeline query.
type PipelineSegment struct {
	With  *WithClause
	Query *CypherQuery
}

// ---------------------------------------------------------------------------
// Convenience constructors
// ---------------------------------------------------------------------------
// Based on goraphdb convenience constructors (litExpr, varRefExpr, etc.),
// extended with new expression types.

func litExpr(v any) Expression {
	return Expression{Kind: ExprLiteral, LitValue: v}
}

func varRefExpr(name string) Expression {
	return Expression{Kind: ExprVarRef, Variable: name}
}

func propExpr(obj, prop string) Expression {
	return Expression{Kind: ExprPropAccess, Object: obj, Property: prop}
}

func propExprForExpr(inner Expression, prop string) Expression {
	return Expression{Kind: ExprPropAccess, NestedPropObject: &inner, Property: prop}
}

func funcCallExpr(name string, args ...Expression) Expression {
	return Expression{Kind: ExprFuncCall, FuncName: name, Args: args}
}

func compExpr(left Expression, op CompOp, right Expression) Expression {
	return Expression{Kind: ExprComparison, Left: &left, Op: op, Right: &right}
}

func andExpr(operands ...Expression) Expression {
	return Expression{Kind: ExprAnd, Operands: operands}
}

func orExpr(operands ...Expression) Expression {
	return Expression{Kind: ExprOr, Operands: operands}
}

func notExpr(inner Expression) Expression {
	return Expression{Kind: ExprNot, Inner: &inner}
}

func inExpr(left, right Expression) Expression {
	return Expression{Kind: ExprIn, InLeft: &left, InRight: &right}
}

func existsExpr(pattern *Pattern) Expression {
	return Expression{Kind: ExprExists, ExistsPattern: pattern}
}

func startsWithExpr(str, prefix Expression) Expression {
	return Expression{Kind: ExprStartsWith, StrExpr: &str, StrValue: &prefix}
}

func endsWithExpr(str, suffix Expression) Expression {
	return Expression{Kind: ExprEndsWith, StrExpr: &str, StrValue: &suffix}
}

func containsExpr(str, sub Expression) Expression {
	return Expression{Kind: ExprContains, StrExpr: &str, StrValue: &sub}
}

func isNullExpr(expr Expression) Expression {
	return Expression{Kind: ExprIsNull, IsNullExpr: &expr}
}

func isNotNullExpr(expr Expression) Expression {
	return Expression{Kind: ExprIsNotNull, IsNullExpr: &expr}
}

func regexMatchExpr(expr, pattern Expression) Expression {
	return Expression{Kind: ExprRegexMatch, RegexExpr: &expr, RegexPattern: &pattern}
}

func caseExpr(operand *Expression, whens []CaseWhen, elseExpr *Expression) Expression {
	return Expression{Kind: ExprCase, CaseOperand: operand, CaseWhens: whens, CaseElse: elseExpr}
}

func listExpr(items ...Expression) Expression {
	return Expression{Kind: ExprListLiteral, ListItems: items}
}

func mapExpr(pairs ...MapPair) Expression {
	return Expression{Kind: ExprMapLiteral, MapPairs: pairs}
}

func aggExpr(fn AggFunc, arg *Expression) Expression {
	return Expression{Kind: ExprAggregation, AggFunc: fn, AggArg: arg}
}

func addExpr(left, right Expression) Expression {
	return Expression{Kind: ExprArithAdd, ArithLeft: &left, ArithRight: &right}
}

func subExpr(left, right Expression) Expression {
	return Expression{Kind: ExprArithSub, ArithLeft: &left, ArithRight: &right}
}

func mulExpr(left, right Expression) Expression {
	return Expression{Kind: ExprArithMul, ArithLeft: &left, ArithRight: &right}
}

func divExpr(left, right Expression) Expression {
	return Expression{Kind: ExprArithDiv, ArithLeft: &left, ArithRight: &right}
}

func modExpr(left, right Expression) Expression {
	return Expression{Kind: ExprArithMod, ArithLeft: &left, ArithRight: &right}
}

func labelsExpr(expr Expression) Expression {
	return Expression{Kind: ExprLabels, LabelsExpr: &expr}
}

func propertiesExpr(expr Expression) Expression {
	return Expression{Kind: ExprProperties, PropertiesExpr: &expr}
}

func intPtr(n int) *int {
	return &n
}

func derefInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

// DDLStatement represents a data definition language statement (CREATE INDEX, CREATE CONSTRAINT, SHOW).
type DDLStatement struct {
	Kind       string
	Label      string
	Property   string
	Constraint string
}

// CypherCall represents a CALL statement for invoking a stored procedure.
// CALL procName(args...) YIELD field1, field2, ... RETURN ...
type CypherCall struct {
	Name       string
	Args       []Expression
	YieldItems []string
}

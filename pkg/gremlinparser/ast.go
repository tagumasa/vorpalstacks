// AST types for the Gremlin traversal parser.

package gremlinparser

// Query is the top-level AST node, containing one or more semicolon-separated statements.
//
//	g.V().has('name','Alice').out('knows').toList(); g.V().count()
type Query struct {
	Statements []Statement
}

// Statement holds a single Gremlin traversal with an optional terminal step.
//
//	g.V().out()           → Traversal with nil Terminal
//	g.V().count().toList() → Traversal with Terminal{Name: "toList"}
type Statement struct {
	Traversal *Traversal
}

// Traversal represents a chained sequence of steps, optionally with a terminal step.
// Source is "__" for nested traversals (__.out()), "" for top-level traversals.
//
//	g.V().out('knows').values('name') → Steps=[V, out, values], Terminal=nil
type Traversal struct {
	Source   string
	Steps    []Step
	Terminal *TerminalStep
}

// TerminalStep is a terminating method such as toList(), next(), iterate(), etc.
type TerminalStep struct {
	Name string
	Args []Argument
}

// Step is a single step in a traversal, with its name, arguments, and modulators.
// Modulators are steps like by(), as(), emit() that modify the preceding step's behaviour.
//
//	has('name', 'Alice').as('a').by('name')
//	→ Step{Name:"has", Args:[...], Modulators:[Step{Name:"as", Args:["a"]}, Step{Name:"by", Args:["name"]}]}
type Step struct {
	Name       string
	Args       []Argument
	Modulators []Step
}

// ArgumentKind identifies the type of an argument in a step.
type ArgumentKind int

const (
	ArgNull            ArgumentKind = iota
	ArgString                       // 'hello'
	ArgInt                          // 42
	ArgFloat                        // 3.14
	ArgBool                         // true / false
	ArgList                         // [1, 2, 3]
	ArgMap                          // [key: val, ...]
	ArgPredicate                    // P.eq(42), TextP.containing('foo')
	ArgNestedTraversal              // __.out('knows')
	ArgEnum                         // T.label, Order.asc, Cardinality.single
	ArgVarRef                       // reserved for future use
)

// Argument is a polymorphic value in a step's argument list.
// Only fields corresponding to the Kind are populated.
type Argument struct {
	Kind  ArgumentKind
	Str   string
	Int   int64
	Float float64
	Bool  bool
	List  []Argument
	Map   []MapEntry
	Pred  *Predicate
	Trav  *Traversal
	Enum  *EnumRef
}

// MapEntry is a single key-value pair in a map literal argument.
type MapEntry struct {
	Key   Argument
	Value Argument
}

// Predicate represents a Gremlin predicate expression (P or TextP).
// Supports chained .and(), .or(), and .negate() combinators.
//
//	P.eq(42)                        → Type:"P", Method:"eq", Args:[42]
//	P.gt(10).and(P.lt(50))          → Type:"P", Method:"and", Args:[P.eq(10), P.lt(50)]
//	TextP.containing('foo').negate() → Type:"TextP", Method:"containing", Negate:true
type Predicate struct {
	Type   string // "P" or "TextP"
	Method string // "eq", "neq", "gt", "gte", "lt", "lte", "between", "inside", "outside", "within", "without", "containing", "notContaining", "startingWith", "notStartingWith", "endingWith", "notEndingWith", "regex", "notRegex", "and", "or"
	Args   []Argument
	And    *Predicate
	Or     *Predicate
	Negate bool
}

// EnumRef represents an enum value such as T.label, Order.asc, Cardinality.single.
// Category is the enum type name (e.g. "T", "Order", "Cardinality"), Value is the member.
type EnumRef struct {
	Category string
	Value    string
}

package gremlinparser

type Query struct {
	Statements []Statement
}

type Statement struct {
	Traversal *Traversal
}

type Traversal struct {
	Source   string
	Steps    []Step
	Terminal *TerminalStep
}

type TerminalStep struct {
	Name string
	Args []Argument
}

type Step struct {
	Name       string
	Args       []Argument
	Modulators []Step
}

type ArgumentKind int

const (
	ArgNull ArgumentKind = iota
	ArgString
	ArgInt
	ArgFloat
	ArgBool
	ArgList
	ArgMap
	ArgPredicate
	ArgNestedTraversal
	ArgEnum
	ArgVarRef
)

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

type MapEntry struct {
	Key   Argument
	Value Argument
}

type Predicate struct {
	Type   string
	Method string
	Args   []Argument
	And    *Predicate
	Or     *Predicate
	Negate bool
}

type EnumRef struct {
	Category string
	Value    string
}

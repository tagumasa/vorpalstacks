package metricmath

import "fmt"

// Expr is the interface implemented by all AST nodes in a parsed metric
// math expression. Each node knows how to evaluate itself given a set
// of variable bindings (time-series data) and can report which variables
// it references.
type Expr interface {
	// Eval evaluates this node and returns the resulting time series.
	Eval(data map[string][]DataPoint) ([]DataPoint, error)

	// References returns the set of variable names this node depends on.
	References() []string
}

// binaryExpr represents a binary arithmetic operation (+, -, *, /, ^).
type binaryExpr struct {
	op    Token
	left  Expr
	right Expr
}

func (e *binaryExpr) Eval(data map[string][]DataPoint) ([]DataPoint, error) {
	left, err := e.left.Eval(data)
	if err != nil {
		return nil, err
	}
	right, err := e.right.Eval(data)
	if err != nil {
		return nil, err
	}
	return alignAndCompute(left, right, e.op.Type), nil
}

func (e *binaryExpr) References() []string {
	return append(e.left.References(), e.right.References()...)
}

// unaryExpr represents a unary minus applied to a sub-expression.
type unaryExpr struct {
	operand Expr
}

func (e *unaryExpr) Eval(data map[string][]DataPoint) ([]DataPoint, error) {
	pts, err := e.operand.Eval(data)
	if err != nil {
		return nil, err
	}
	result := make([]DataPoint, len(pts))
	for i, p := range pts {
		result[i] = DataPoint{Timestamp: p.Timestamp, Value: -p.Value}
	}
	return result, nil
}

func (e *unaryExpr) References() []string {
	return e.operand.References()
}

// numberExpr represents a numeric literal. When evaluated, it produces
// a single data point with value n at a zero timestamp. When combined
// with a time series via a binary operator, alignAndCompute replicates
// it at each timestamp of the other operand.
type numberExpr struct {
	value float64
}

func (e *numberExpr) Eval(_ map[string][]DataPoint) ([]DataPoint, error) {
	return []DataPoint{{Value: e.value}}, nil
}

func (e *numberExpr) References() []string { return nil }

// varRefExpr represents a reference to another metric data query by ID
// (e.g., m1, e1). The variable name is looked up in the data map.
type varRefExpr struct {
	name string
}

func (e *varRefExpr) Eval(data map[string][]DataPoint) ([]DataPoint, error) {
	pts, ok := data[e.name]
	if !ok {
		return nil, fmt.Errorf("metricmath: undefined variable %q", e.name)
	}
	result := make([]DataPoint, len(pts))
	copy(result, pts)
	return result, nil
}

func (e *varRefExpr) References() []string { return []string{e.name} }

// funcCallExpr represents a function call such as FILL(m1, 0) or AVG(m1).
type funcCallExpr struct {
	name string
	args []Expr
}

func (e *funcCallExpr) Eval(data map[string][]DataPoint) ([]DataPoint, error) {
	args := make([][]DataPoint, len(e.args))
	for i, arg := range e.args {
		val, err := arg.Eval(data)
		if err != nil {
			return nil, err
		}
		args[i] = val
	}
	return evalFunction(e.name, args)
}

func (e *funcCallExpr) References() []string {
	var refs []string
	for _, arg := range e.args {
		refs = append(refs, arg.References()...)
	}
	return refs
}

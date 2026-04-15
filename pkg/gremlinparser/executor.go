// Gremlin traversal executor.
//
// Executes a parsed Gremlin Query AST against graphengine.GraphReader/GraphWriter.
//
// Architecture:
//   - Traverser wraps each element as it flows through the step pipeline
//   - StepExecutor is the function signature for each registered step
//   - Steps are registered via RegisterStep() in init() functions across files
//   - The source step (V, E, addV, etc.) produces initial traversers
//   - Subsequent steps consume traversers from the previous step
//   - Modulators (by, as, emit, etc.) are applied after each step
//   - Terminal steps (toList, next, iterate, etc.) produce the final result

package gremlinparser

import (
	"context"
	"fmt"

	"vorpalstacks/pkg/graphengine"
)

// Traverser carries an element through the traversal pipeline along with
// metadata: path history, loop counters, bulk count, tags, and sack value.
type Traverser struct {
	Element any
	Path    []any
	Loops   map[string]int
	Bulk    int
	Tags    map[string]any
	Sack    any
}

// newTraverser creates a Traverser with a single-element path, empty loops/tags, and bulk=1.
func newTraverser(element any) *Traverser {
	return &Traverser{
		Element: element,
		Path:    []any{element},
		Loops:   make(map[string]int),
		Tags:    make(map[string]any),
		Bulk:    1,
	}
}

// clone creates a deep copy of the traverser, safe to mutate independently.
func (t *Traverser) clone() *Traverser {
	path := make([]any, len(t.Path))
	copy(path, t.Path)
	loops := make(map[string]int, len(t.Loops))
	for k, v := range t.Loops {
		loops[k] = v
	}
	tags := make(map[string]any, len(t.Tags))
	for k, v := range t.Tags {
		tags[k] = v
	}
	return &Traverser{
		Element: t.Element,
		Path:    path,
		Loops:   loops,
		Bulk:    t.Bulk,
		Tags:    tags,
		Sack:    t.Sack,
	}
}

// ExecContext carries the execution state: context, graph reader/writer, parameters, and side-effects.
type ExecContext struct {
	Ctx         context.Context
	Reader      graphengine.GraphReader
	Writer      graphengine.GraphWriter
	Params      map[string]any
	SideEffects map[string]any
}

// StepExecutor is the function signature for a step implementation.
// Receives the execution context, incoming traversers, and the step definition;
// returns the outgoing traversers.
type StepExecutor func(*ExecContext, []*Traverser, Step) ([]*Traverser, error)

// stepRegistry maps step names to their executor functions.
var stepRegistry = map[string]StepExecutor{}

// RegisterStep adds a step executor to the registry. Called from init() in step files.
func RegisterStep(name string, executor StepExecutor) {
	stepRegistry[name] = executor
}

// ExecuteQuery runs all statements in a query against the graph, returning the combined result.
// For single-statement queries, returns the result directly; for multiple statements,
// flattens all results into a single slice.
func ExecuteQuery(ctx context.Context, reader graphengine.GraphReader, writer graphengine.GraphWriter, query *Query, params map[string]any) (any, error) {
	if len(query.Statements) == 0 {
		return nil, nil
	}

	ec := &ExecContext{
		Ctx:         ctx,
		Reader:      reader,
		Writer:      writer,
		Params:      params,
		SideEffects: make(map[string]any),
	}

	var results []any

	for _, stmt := range query.Statements {
		result, err := executeStatement(ec, &stmt)
		if err != nil {
			return nil, err
		}
		results = append(results, result...)
	}

	if len(results) == 1 {
		return results[0], nil
	}
	if len(results) == 0 {
		return []any{}, nil
	}
	flattened := make([]any, 0)
	for _, r := range results {
		if arr, ok := r.([]any); ok {
			flattened = append(flattened, arr...)
		} else {
			flattened = append(flattened, r)
		}
	}
	return flattened, nil
}

// executeStatement runs a single traversal statement through the step pipeline.
// The first step is treated as a source step (receives nil traversers).
// Modulators are applied after each step.
func executeStatement(ec *ExecContext, stmt *Statement) ([]any, error) {
	trav := stmt.Traversal
	if len(trav.Steps) == 0 {
		return nil, nil
	}

	var traversers []*Traverser

	sourceStep := trav.Steps[0]
	sourceExec, ok := stepRegistry[sourceStep.Name]
	if !ok {
		return nil, fmt.Errorf("gremlin: unknown step %q", sourceStep.Name)
	}

	traversers, err := sourceExec(ec, nil, sourceStep)
	if err != nil {
		return nil, err
	}
	traversers, err = applyModulators(ec, traversers, sourceStep)
	if err != nil {
		return nil, err
	}

	for i := 1; i < len(trav.Steps); i++ {
		step := trav.Steps[i]
		exec, ok := stepRegistry[step.Name]
		if !ok {
			return nil, fmt.Errorf("gremlin: unknown step %q", step.Name)
		}

		if step.Name == "addE" && len(step.Modulators) > 0 {
			traversers, err = applyModulators(ec, traversers, step)
			if err != nil {
				return nil, err
			}
			traversers, err = exec(ec, traversers, step)
			if err != nil {
				return nil, err
			}
			continue
		}

		traversers, err = exec(ec, traversers, step)
		if err != nil {
			return nil, err
		}
		traversers, err = applyModulators(ec, traversers, step)
		if err != nil {
			return nil, err
		}
	}

	if trav.Terminal != nil {
		return executeTerminal(ec, traversers, trav.Terminal)
	}

	result := make([]any, len(traversers))
	for i, t := range traversers {
		result[i] = t.Element
	}
	return result, nil
}

// applyModulators runs all modulator steps attached to a step in order.
func applyModulators(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	for _, mod := range step.Modulators {
		exec, ok := stepRegistry[mod.Name]
		if !ok {
			continue
		}
		var err error
		traversers, err = exec(ec, traversers, mod)
		if err != nil {
			return nil, err
		}
	}
	return traversers, nil
}

// executeTerminal handles terminal steps that produce the final result.
// toList/toBulkSet returns all elements; toSet deduplicates; next returns n elements;
// tryNext returns at most one; iterate discards all; hasNext returns boolean.
func executeTerminal(ec *ExecContext, traversers []*Traverser, term *TerminalStep) ([]any, error) {
	switch term.Name {
	case "toList":
		result := make([]any, len(traversers))
		for i, t := range traversers {
			result[i] = t.Element
		}
		return result, nil
	case "toSet", "toBulkSet":
		seen := make(map[string]bool)
		var result []any
		for _, t := range traversers {
			key := dedupKey(t.Element)
			if !seen[key] {
				seen[key] = true
				result = append(result, t.Element)
			}
		}
		return result, nil
	case "next":
		n := 1
		if len(term.Args) > 0 && term.Args[0].Kind == ArgInt {
			n = int(term.Args[0].Int)
		}
		result := make([]any, 0, n)
		for i := 0; i < n && i < len(traversers); i++ {
			result = append(result, traversers[i].Element)
		}
		return result, nil
	case "tryNext":
		if len(traversers) > 0 {
			return []any{traversers[0].Element}, nil
		}
		return []any{nil}, nil
	case "iterate":
		return nil, nil
	case "hasNext":
		if len(traversers) > 0 {
			return []any{true}, nil
		}
		return []any{false}, nil
	default:
		result := make([]any, len(traversers))
		for i, t := range traversers {
			result[i] = t.Element
		}
		return result, nil
	}
}

// asNode extracts a *graphengine.Node from a traverser's element, returning false if not a node.
func asNode(t *Traverser) (*graphengine.Node, bool) {
	n, ok := t.Element.(*graphengine.Node)
	return n, ok
}

// asEdge extracts a *graphengine.Edge from a traverser's element, returning false if not an edge.
func asEdge(t *Traverser) (*graphengine.Edge, bool) {
	e, ok := t.Element.(*graphengine.Edge)
	return e, ok
}

// argString returns the string value of an argument.
func argString(arg Argument) string {
	return arg.Str
}

// argInt returns the integer value of an argument.
func argInt(arg Argument) int64 {
	return arg.Int
}

// argStrings extracts a slice of strings from a list of arguments.
func argStrings(args []Argument) []string {
	result := make([]string, len(args))
	for i, a := range args {
		result[i] = a.Str
	}
	return result
}

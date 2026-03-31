package gremlinparser

import (
	"context"
	"fmt"

	"vorpalstacks/pkg/graphengine"
)

type Traverser struct {
	Element any
	Path    []any
	Loops   map[string]int
	Bulk    int
	Tags    map[string]any
	Sack    any
}

func newTraverser(element any) *Traverser {
	return &Traverser{
		Element: element,
		Path:    []any{element},
		Loops:   make(map[string]int),
		Tags:    make(map[string]any),
		Bulk:    1,
	}
}

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

type ExecContext struct {
	Ctx         context.Context
	Reader      graphengine.GraphReader
	Writer      graphengine.GraphWriter
	Params      map[string]any
	SideEffects map[string]any
}

type StepExecutor func(*ExecContext, []*Traverser, Step) ([]*Traverser, error)

var stepRegistry = map[string]StepExecutor{}

func RegisterStep(name string, executor StepExecutor) {
	stepRegistry[name] = executor
}

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

func asNode(t *Traverser) (*graphengine.Node, bool) {
	n, ok := t.Element.(*graphengine.Node)
	return n, ok
}

func asEdge(t *Traverser) (*graphengine.Edge, bool) {
	e, ok := t.Element.(*graphengine.Edge)
	return e, ok
}

func argString(arg Argument) string {
	return arg.Str
}

func argInt(arg Argument) int64 {
	return arg.Int
}

func argStrings(args []Argument) []string {
	result := make([]string, len(args))
	for i, a := range args {
		result[i] = a.Str
	}
	return result
}

func resolveEnum(arg Argument) string {
	if arg.Kind == ArgEnum && arg.Enum != nil {
		return arg.Enum.Value
	}
	return arg.Str
}

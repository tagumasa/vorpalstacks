package neptunedata

import (
	"fmt"

	"vorpalstacks/pkg/gremlinparser"
)

// explainGremlinQuery parses a Gremlin query and produces a step-by-step
// explanation of the traversal plan.
func explainGremlinQuery(query string) map[string]interface{} {
	parsed, err := gremlinparser.Parse(query)
	if err != nil {
		return map[string]interface{}{
			"error": fmt.Sprintf("parse error: %v", err),
		}
	}

	steps := []map[string]interface{}{}
	if len(parsed.Statements) > 0 && parsed.Statements[0].Traversal != nil {
		for _, step := range parsed.Statements[0].Traversal.Steps {
			stepInfo := map[string]interface{}{
				"name": step.Name,
			}
			if len(step.Args) > 0 {
				args := make([]interface{}, len(step.Args))
				for i, arg := range step.Args {
					args[i] = describeArg(arg)
				}
				stepInfo["args"] = args
			}
			steps = append(steps, stepInfo)
		}
	}

	return map[string]interface{}{
		"steps": steps,
	}
}

// profileGremlinQuery returns an explain plan augmented with profiling metrics.
func profileGremlinQuery(query string) map[string]interface{} {
	explain := explainGremlinQuery(query)
	explain["profile"] = map[string]interface{}{
		"metrics": map[string]interface{}{
			"dur":     0,
			"count":   1,
			"size":    0,
			"time":    0,
			"incTime": 0,
		},
		"indices": map[string]interface{}{},
	}
	return explain
}

// describeArg converts a Gremlin argument to a serialisable representation for
// explain/profile output.
func describeArg(arg gremlinparser.Argument) interface{} {
	switch arg.Kind {
	case gremlinparser.ArgString:
		return arg.Str
	case gremlinparser.ArgInt:
		return arg.Int
	case gremlinparser.ArgFloat:
		return arg.Float
	case gremlinparser.ArgBool:
		return arg.Bool
	case gremlinparser.ArgNull:
		return nil
	case gremlinparser.ArgEnum:
		if arg.Enum != nil {
			return arg.Enum.Value
		}
		return nil
	case gremlinparser.ArgList:
		items := make([]interface{}, len(arg.List))
		for i, a := range arg.List {
			items[i] = describeArg(a)
		}
		return items
	case gremlinparser.ArgMap:
		m := make(map[string]interface{})
		for _, entry := range arg.Map {
			if key, ok := describeArg(entry.Key).(string); ok {
				m[key] = describeArg(entry.Value)
			}
		}
		return m
	case gremlinparser.ArgPredicate:
		if arg.Pred != nil {
			return map[string]interface{}{
				"type":   arg.Pred.Type,
				"method": arg.Pred.Method,
				"args":   describeArgs(arg.Pred.Args),
			}
		}
		return nil
	case gremlinparser.ArgNestedTraversal:
		if arg.Trav != nil {
			steps := []map[string]interface{}{}
			for _, step := range arg.Trav.Steps {
				stepInfo := map[string]interface{}{
					"name": step.Name,
				}
				if len(step.Args) > 0 {
					args := make([]interface{}, len(step.Args))
					for i, a := range step.Args {
						args[i] = describeArg(a)
					}
					stepInfo["args"] = args
				}
				steps = append(steps, stepInfo)
			}
			return steps
		}
		return nil
	default:
		return arg.Str
	}
}

// describeArgs converts a slice of Gremlin arguments to serialisable representations.
func describeArgs(args []gremlinparser.Argument) []interface{} {
	result := make([]interface{}, len(args))
	for i, a := range args {
		result[i] = describeArg(a)
	}
	return result
}

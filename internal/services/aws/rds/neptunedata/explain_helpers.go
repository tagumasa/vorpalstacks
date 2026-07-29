package neptunedata

import (
	"fmt"
	"strings"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/pkg/gremlinparser"
)

// explainGremlinQuery parses a Gremlin query and produces a step-by-step
// explanation of the traversal plan. Returns an error if parsing fails.
func explainGremlinQuery(query string) (map[string]interface{}, error) {
	parsed, err := gremlinparser.Parse(query)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	var steps []map[string]interface{}
	if len(parsed.Statements) > 0 && parsed.Statements[0].Traversal != nil {
		steps = traversalToSteps(parsed.Statements[0].Traversal)
	}

	return map[string]interface{}{
		"steps": steps,
	}, nil
}

// traversalToSteps converts a Gremlin traversal into a step-by-step
// explanation. Shared between top-level explain and nested traversal arguments.
func traversalToSteps(trav *gremlinparser.Traversal) []map[string]interface{} {
	steps := []map[string]interface{}{}
	for _, step := range trav.Steps {
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
	return steps
}

// profileOptions holds the parameters from the Smithy jsonName-prefixed
// profile query fields. These control the level of detail in profiling output.
type profileOptions struct {
	results    *bool  // profile.results: include step results
	chop       *int   // profile.chop: truncate result strings to N chars
	serializer string // profile.serializer: output format for results
	indexOps   *bool  // profile.indexOps: include index operation details
}

// formatExplainOutput converts the explain plan map to a text representation
// matching the AWS Neptune explain query response format.
func formatExplainOutput(plan map[string]interface{}) string {
	steps, ok := plan["steps"].([]map[string]interface{})
	if !ok {
		return ""
	}
	var b strings.Builder
	b.WriteString("***************************************************\n")
	b.WriteString("*                 Gremlin Explain                 *\n")
	b.WriteString("***************************************************\n")
	for i, step := range steps {
		name, _ := step["name"].(string)
		fmt.Fprintf(&b, "%d\t%s", i+1, name)
		if args, ok := step["args"].([]interface{}); ok && len(args) > 0 {
			b.WriteString("(")
			for j, arg := range args {
				if j > 0 {
					b.WriteString(",")
				}
				fmt.Fprintf(&b, "%v", arg)
			}
			b.WriteString(")")
		}
		b.WriteString("\n")
	}
	return b.String()
}

// formatProfileOutput converts the profile plan map to a text representation
// matching the AWS Neptune profile query response format. Includes profiling
// metrics when profile options are set.
func formatProfileOutput(plan map[string]interface{}, opts profileOptions) string {
	var b strings.Builder
	b.WriteString("***************************************************\n")
	b.WriteString("*                 Gremlin Profile                 *\n")
	b.WriteString("***************************************************\n")

	if steps, ok := plan["steps"].([]map[string]interface{}); ok {
		for i, step := range steps {
			name, _ := step["name"].(string)
			fmt.Fprintf(&b, "%d\t%s", i+1, name)
			if args, ok := step["args"].([]interface{}); ok && len(args) > 0 {
				b.WriteString("(")
				for j, arg := range args {
					if j > 0 {
						b.WriteString(",")
					}
					fmt.Fprintf(&b, "%v", arg)
				}
				b.WriteString(")")
			}
			b.WriteString("\n")

			// Include profiling metrics when available
			if prof, ok := step["profile"].(map[string]interface{}); ok {
				fmt.Fprintf(&b, "\t\tmetrics: ")
				metrics, _ := prof["metrics"].(map[string]interface{})
				first := true
				for k, v := range metrics {
					if !first {
						b.WriteString(", ")
					}
					fmt.Fprintf(&b, "%s=%v", k, v)
					first = false
				}
				b.WriteString("\n")
			}
		}
	}

	// Include index operations if requested
	if opts.indexOps != nil && *opts.indexOps {
		if prof, ok := plan["profile"].(map[string]interface{}); ok {
			if indices, ok := prof["indices"].(map[string]interface{}); ok && len(indices) > 0 {
				b.WriteString("\nIndex Operations:\n")
				for name, data := range indices {
					fmt.Fprintf(&b, "\t%s: %v\n", name, data)
				}
			}
		}
	}

	return b.String()
}

// profileGremlinQueryEx returns an explain plan augmented with profiling
// metrics, honouring the profile.* options.
func profileGremlinQueryEx(query string, opts profileOptions) (map[string]interface{}, error) {
	plan, err := explainGremlinQuery(query)
	if err != nil {
		return nil, err
	}

	chopVal := 250
	if opts.chop != nil {
		chopVal = *opts.chop
	}

	profMetrics := map[string]interface{}{
		"metrics": map[string]interface{}{
			"dur":        0,
			"count":      1,
			"size":       0,
			"time":       0,
			"incTime":    0,
			"traversers": map[string]interface{}{"count": 0},
		},
		"indices": map[string]interface{}{},
		"results": map[string]interface{}{
			"chop": chopVal,
		},
	}

	// Attach profiling data to each step
	if steps, ok := plan["steps"].([]map[string]interface{}); ok {
		for i := range steps {
			steps[i]["profile"] = profMetrics
		}
	}

	plan["profile"] = profMetrics
	return plan, nil
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
			return traversalToSteps(arg.Trav)
		}
		return nil
	default:
		logs.Warn("describeArg: unhandled arg kind", logs.Int("kind", int(arg.Kind)))
		return fmt.Sprintf("<unknown:%d>", arg.Kind)
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

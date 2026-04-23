// Source step executors: V, E, addV, addE graph element access and creation.
//
// This file also contains the single init() registration point for all steps
// defined across steps_source.go, steps_traversal.go, steps_property.go,
// steps_match_math.go, steps_filter.go, and helpers.go.

package gremlinparser

import (
	"fmt"

	"vorpalstacks/pkg/graphengine"
)

func init() {
	RegisterStep("V", execV)
	RegisterStep("E", execE)
	RegisterStep("addV", execAddV)
	RegisterStep("addE", execAddE)
	RegisterStep("out", execOut)
	RegisterStep("in", execIn)
	RegisterStep("both", execBoth)
	RegisterStep("outE", execOutE)
	RegisterStep("inE", execInE)
	RegisterStep("bothE", execBothE)
	RegisterStep("outV", execOutV)
	RegisterStep("inV", execInV)
	RegisterStep("otherV", execOtherV)
	RegisterStep("has", execHas)
	RegisterStep("hasLabel", execHasLabel)
	RegisterStep("hasId", execHasId)
	RegisterStep("hasNot", execHasNot)
	RegisterStep("values", execValues)
	RegisterStep("valueMap", execValueMap)
	RegisterStep("properties", execProperties)
	RegisterStep("propertyMap", execPropertyMap)
	RegisterStep("id", execId)
	RegisterStep("label", execLabel)
	RegisterStep("property", execProperty)
	RegisterStep("drop", execDrop)
	RegisterStep("count", execCount)
	RegisterStep("identity", execIdentity)
	RegisterStep("keys", execKeys)
	RegisterStep("sample", execSample)
	RegisterStep("sack", execSack)
	RegisterStep("match", execMatch)
	RegisterStep("math", execMath)
}

// execV is the vertex source step. With no arguments, iterates all vertices;
// with arguments, looks up vertices by ID.
func execV(ec *ExecContext, _ []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) == 0 {
		var traversers []*Traverser
		err := ec.Reader.ForEachNode(func(n *graphengine.Node) error {
			traversers = append(traversers, newTraverser(n))
			return nil
		})
		return traversers, err
	}

	var traversers []*Traverser
	for _, arg := range step.Args {
		id, err := toNodeID(arg)
		if err != nil {
			continue
		}
		n, err := ec.Reader.GetNode(id)
		if err != nil {
			continue
		}
		traversers = append(traversers, newTraverser(n))
	}
	return traversers, nil
}

// execE is the edge source step. With no arguments, iterates all edges;
// with arguments, looks up edges by ID.
func execE(ec *ExecContext, _ []*Traverser, step Step) ([]*Traverser, error) {
	if len(step.Args) == 0 {
		var traversers []*Traverser
		err := ec.Reader.ForEachEdge(func(e *graphengine.Edge) error {
			traversers = append(traversers, newTraverser(e))
			return nil
		})
		return traversers, err
	}

	var traversers []*Traverser
	for _, arg := range step.Args {
		id, err := toEdgeID(arg)
		if err != nil {
			continue
		}
		e, err := ec.Reader.GetEdge(id)
		if err != nil {
			continue
		}
		traversers = append(traversers, newTraverser(e))
	}
	return traversers, nil
}

// execAddV creates a new vertex with the given label(s).
func execAddV(ec *ExecContext, _ []*Traverser, step Step) ([]*Traverser, error) {
	if ec.Writer == nil {
		return nil, fmt.Errorf("gremlin: addV requires a GraphWriter")
	}

	var labels []string
	if len(step.Args) > 0 {
		labels = []string{argString(step.Args[0])}
	}

	id, err := ec.Writer.AddNode(labels, nil)
	if err != nil {
		return nil, err
	}
	n, err := ec.Reader.GetNode(id)
	if err != nil {
		return nil, err
	}
	return []*Traverser{newTraverser(n)}, nil
}

// execAddE creates a new edge with the given label. The source vertex comes from
// the incoming traverser (or a from() modulator), and the target from a to() modulator.
func execAddE(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	if ec.Writer == nil {
		return nil, fmt.Errorf("gremlin: addE requires a GraphWriter")
	}

	if len(step.Args) == 0 {
		return nil, fmt.Errorf("gremlin: addE requires a label argument")
	}
	label := argString(step.Args[0])

	if len(traversers) == 0 {
		return nil, fmt.Errorf("gremlin: addE requires a source traverser")
	}

	tags := traversers[0].Tags
	fromID, hasFrom := tags["__from"]
	toID, hasTo := tags["__to"]

	if !hasFrom {
		sourceNode, ok := asNode(traversers[0])
		if ok {
			for _, t := range traversers {
				t.Tags["__from"] = sourceNode.ID
			}
			fromID = sourceNode.ID
			hasFrom = true
		}
	}

	if !hasTo {
		if len(traversers) > 1 {
			targetNode, ok := asNode(traversers[1])
			if ok {
				for _, t := range traversers {
					t.Tags["__to"] = targetNode.ID
				}
				toID = targetNode.ID
				hasTo = true
			}
		} else {
			return nil, fmt.Errorf("gremlin: addE requires a to() step when there is only one traverser")
		}
	}

	if !hasFrom || !hasTo {
		return nil, fmt.Errorf("gremlin: addE requires both from() and to() steps")
	}

	fromNodeID, ok := fromID.(graphengine.NodeID)
	if !ok {
		return nil, fmt.Errorf("gremlin: addE from ID is not a valid NodeID")
	}
	toNodeIDVal, ok := toID.(graphengine.NodeID)
	if !ok {
		return nil, fmt.Errorf("gremlin: addE to ID is not a valid NodeID")
	}

	var result []*Traverser
	for _, t := range traversers {
		edgeID, err := ec.Writer.AddEdge(fromNodeID, toNodeIDVal, label, nil)
		if err != nil {
			return nil, err
		}

		e, err := ec.Reader.GetEdge(edgeID)
		if err != nil {
			return nil, err
		}

		nt := t.clone()
		nt.Element = e
		nt.Path = append(nt.Path, e)
		result = append(result, nt)
	}

	return result, nil
}

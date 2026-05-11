// Traversal step executors: out, in, both, outE, inE, bothE, outV, inV, otherV.

package gremlinparser

import (
	"fmt"

	"vorpalstacks/internal/core/storage/graphengine"
)

// execOut traverses outgoing edges to adjacent vertices, filtered by optional edge labels.
func execOut(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	labels := argStrings(step.Args)
	var result []*Traverser

	for _, t := range traversers {
		n, ok := asNode(t)
		if !ok {
			continue
		}
		edges, err := ec.Reader.GetEdges(n.ID, graphengine.Outgoing, singleLabelFilter(labels))
		if err != nil {
			continue
		}
		edges = filterEdgesByLabels(edges, labels)
		for _, e := range edges {
			target, err := ec.Reader.GetNode(e.To)
			if err != nil {
				continue
			}
			nt := t.clone()
			nt.Element = target
			nt.Path = append(nt.Path, target)
			result = append(result, nt)
		}
	}
	return result, nil
}

// execIn traverses incoming edges to adjacent vertices, filtered by optional edge labels.
func execIn(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	labels := argStrings(step.Args)
	var result []*Traverser

	for _, t := range traversers {
		n, ok := asNode(t)
		if !ok {
			continue
		}
		edges, err := ec.Reader.GetEdges(n.ID, graphengine.Incoming, singleLabelFilter(labels))
		if err != nil {
			continue
		}
		edges = filterEdgesByLabels(edges, labels)
		for _, e := range edges {
			target, err := ec.Reader.GetNode(e.From)
			if err != nil {
				continue
			}
			nt := t.clone()
			nt.Element = target
			nt.Path = append(nt.Path, target)
			result = append(result, nt)
		}
	}
	return result, nil
}

// execBoth traverses both incoming and outgoing edges to adjacent vertices.
func execBoth(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	labels := argStrings(step.Args)
	var result []*Traverser

	for _, t := range traversers {
		n, ok := asNode(t)
		if !ok {
			continue
		}
		edges, err := ec.Reader.GetEdges(n.ID, graphengine.Both, singleLabelFilter(labels))
		if err != nil {
			continue
		}
		edges = filterEdgesByLabels(edges, labels)
		for _, e := range edges {
			var targetID graphengine.NodeID
			if e.To == n.ID {
				targetID = e.From
			} else {
				targetID = e.To
			}
			target, err := ec.Reader.GetNode(targetID)
			if err != nil {
				continue
			}
			nt := t.clone()
			nt.Element = target
			nt.Path = append(nt.Path, target)
			result = append(result, nt)
		}
	}
	return result, nil
}

func execOutE(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	return execEdgeTraversal(ec, traversers, step, graphengine.Outgoing, "outE")
}

func execInE(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	return execEdgeTraversal(ec, traversers, step, graphengine.Incoming, "inE")
}

func execBothE(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	return execEdgeTraversal(ec, traversers, step, graphengine.Both, "bothE")
}

func execEdgeTraversal(ec *ExecContext, traversers []*Traverser, step Step, dir graphengine.Direction, stepName string) ([]*Traverser, error) {
	labels := argStrings(step.Args)
	var result []*Traverser

	for _, t := range traversers {
		n, ok := asNode(t)
		if !ok {
			return nil, fmt.Errorf("gremlin: %s() requires vertex input, got %T", stepName, t.Element)
		}
		edges, err := ec.Reader.GetEdges(n.ID, dir, singleLabelFilter(labels))
		if err != nil {
			continue
		}
		edges = filterEdgesByLabels(edges, labels)
		for _, e := range edges {
			nt := t.clone()
			nt.Element = e
			nt.Path = append(nt.Path, e)
			result = append(result, nt)
		}
	}
	return result, nil
}

// execOutV returns the source vertex of each edge.
func execOutV(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	var result []*Traverser
	for _, t := range traversers {
		e, ok := asEdge(t)
		if !ok {
			return nil, fmt.Errorf("gremlin: outV() requires edge input, got %T", t.Element)
		}
		n, err := ec.Reader.GetNode(e.From)
		if err != nil {
			continue
		}
		nt := t.clone()
		nt.Element = n
		nt.Path = append(nt.Path, n)
		result = append(result, nt)
	}
	return result, nil
}

// execInV returns the target vertex of each edge.
func execInV(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	var result []*Traverser
	for _, t := range traversers {
		e, ok := asEdge(t)
		if !ok {
			return nil, fmt.Errorf("gremlin: inV() requires edge input, got %T", t.Element)
		}
		n, err := ec.Reader.GetNode(e.To)
		if err != nil {
			continue
		}
		nt := t.clone()
		nt.Element = n
		nt.Path = append(nt.Path, n)
		result = append(result, nt)
	}
	return result, nil
}

// execOtherV returns the vertex at the opposite end of each edge from the traverser's path.
func execOtherV(ec *ExecContext, traversers []*Traverser, step Step) ([]*Traverser, error) {
	var result []*Traverser
	for _, t := range traversers {
		e, ok := asEdge(t)
		if !ok {
			continue
		}
		var otherID graphengine.NodeID
		if len(t.Path) > 1 {
			if prevNode, ok := t.Path[len(t.Path)-2].(*graphengine.Node); ok {
				if e.From == prevNode.ID {
					otherID = e.To
				} else {
					otherID = e.From
				}
			}
		}
		if otherID == 0 {
			continue
		}
		n, err := ec.Reader.GetNode(otherID)
		if err != nil {
			continue
		}
		nt := t.clone()
		nt.Element = n
		nt.Path = append(nt.Path, n)
		result = append(result, nt)
	}
	return result, nil
}

// singleLabelFilter returns the single label if exactly one is provided, otherwise ""
// (the graph engine returns all edges when label filter is empty).
func singleLabelFilter(labels []string) string {
	if len(labels) == 1 {
		return labels[0]
	}
	return ""
}

// filterEdgesByLabels filters edges by label when more than one label is specified.
// For a single label, the graph engine already handles filtering efficiently.
func filterEdgesByLabels(edges []*graphengine.Edge, labels []string) []*graphengine.Edge {
	if len(labels) <= 1 {
		return edges
	}
	labelSet := make(map[string]bool, len(labels))
	for _, l := range labels {
		labelSet[l] = true
	}
	filtered := make([]*graphengine.Edge, 0, len(edges))
	for _, e := range edges {
		if labelSet[e.Label] {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

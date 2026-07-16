package arn

// NeptuneGraphBuilder provides methods for constructing NeptuneGraph ARNs.
type NeptuneGraphBuilder struct{ *ARNBuilder }

// NeptuneGraph returns a NeptuneGraphBuilder for constructing NeptuneGraph ARNs.
func (b *ARNBuilder) NeptuneGraph() *NeptuneGraphBuilder { return &NeptuneGraphBuilder{b} }

// Graph constructs an ARN for a NeptuneGraph graph.
func (b *NeptuneGraphBuilder) Graph(id string) string {
	return b.Build("neptune-graph", "graph/"+id)
}

// Snapshot constructs an ARN for a NeptuneGraph snapshot.
func (b *NeptuneGraphBuilder) Snapshot(id string) string {
	return b.Build("neptune-graph", "snapshot/"+id)
}

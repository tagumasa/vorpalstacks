// Cypher query plan — tree of operators for EXPLAIN and PROFILE.
//
// Based on goraphdb/query_plan.go (PlanNode, QueryPlan, buildPlan,
// buildNodeMatchPlan, buildSingleHopPlan, buildVarLengthPlan, formatting),
// adapted for GraphReader interface and our AST types.
//
// Architecture:
//   - EXPLAIN: builds and returns a plan tree without executing the query
//   - PROFILE: executes the query normally, returning both the plan tree
//     and the result (timing instrumentation is a future enhancement)
//
// Plan operators:
//   - AllNodesScan: full node scan (no index)
//   - NodeByLabelScan: label index scan
//   - NodeIndexSeek: property index scan
//   - Filter: inline property or WHERE clause filtering
//   - Expand(All): single-hop relationship traversal
//   - VarLengthExpand: variable-length path traversal
//   - Projection: WITH or RETURN column projection
//   - Sort: ORDER BY
//   - Limit: SKIP + LIMIT
//   - ProduceResults: final result output

package cypherparser

import (
	"fmt"
	"strings"

	"vorpalstacks/pkg/graphengine"
)

// PlanOperator identifies a type of operator in a query execution plan.
type PlanOperator string

// Plan operator constants.
const (
	OpAllNodesScan  PlanOperator = "AllNodesScan"
	OpLabelScan     PlanOperator = "NodeByLabelScan"
	OpIndexSeek     PlanOperator = "NodeIndexSeek"
	OpFilter        PlanOperator = "Filter"
	OpExpand        PlanOperator = "Expand(All)"
	OpExpandVarLen  PlanOperator = "VarLengthExpand"
	OpProjection    PlanOperator = "Projection"
	OpSort          PlanOperator = "Sort"
	OpLimit         PlanOperator = "Limit"
	OpProduceResult PlanOperator = "ProduceResults"
	OpUnwind        PlanOperator = "Unwind"
	OpWith          PlanOperator = "With"
)

// PlanNode represents a single node in a query execution plan tree.
type PlanNode struct {
	Operator   PlanOperator   `json:"operator"`
	Details    string         `json:"details,omitempty"`
	EstRows    int            `json:"estimatedRows"`
	ActualRows int            `json:"actualRows,omitempty"`
	Args       map[string]any `json:"args,omitempty"`
	Children   []*PlanNode    `json:"children,omitempty"`
}

// QueryPlan holds the root of a query execution plan and, for PROFILE queries, the query result.
type QueryPlan struct {
	Root    *PlanNode     `json:"root"`
	Profile bool          `json:"profile"`
	Result  *CypherResult `json:"result,omitempty"`
}

// String returns a human-readable tree representation of the query plan.
func (qp *QueryPlan) String() string {
	var sb strings.Builder
	if qp.Profile {
		sb.WriteString("PROFILE:\n")
	} else {
		sb.WriteString("EXPLAIN:\n")
	}
	qp.Root.format(&sb, "", true)
	return sb.String()
}

func (n *PlanNode) format(sb *strings.Builder, prefix string, isLast bool) {
	connector := "├── "
	if isLast {
		connector = "└── "
	}
	if prefix == "" {
		connector = ""
	}

	sb.WriteString(prefix)
	sb.WriteString(connector)
	sb.WriteString(string(n.Operator))
	if n.Details != "" {
		sb.WriteString(" (")
		sb.WriteString(n.Details)
		sb.WriteString(")")
	}

	if n.ActualRows > 0 {
		sb.WriteString(fmt.Sprintf(" [rows=%d]", n.ActualRows))
	} else if n.EstRows > 0 {
		sb.WriteString(fmt.Sprintf(" [est. rows=%d]", n.EstRows))
	}

	sb.WriteString("\n")

	childPrefix := prefix
	if prefix != "" {
		if isLast {
			childPrefix += "    "
		} else {
			childPrefix += "│   "
		}
	}

	for i, child := range n.Children {
		child.format(sb, childPrefix, i == len(n.Children)-1)
	}
}

func buildPlan(q *CypherQuery, reader graphengine.GraphReader) *PlanNode {
	pat := q.Match.Pattern

	var scanNode *PlanNode

	switch {
	case len(pat.Nodes) == 1 && len(pat.Rels) == 0:
		scanNode = buildNodeMatchPlan(q, reader)

	case len(pat.Nodes) == 2 && len(pat.Rels) == 1:
		rel := pat.Rels[0]
		if rel.VarLength {
			scanNode = buildVarLengthPlan(q)
		} else {
			scanNode = buildSingleHopPlan(q, reader)
		}

	default:
		scanNode = &PlanNode{Operator: OpAllNodesScan, Details: "unsupported pattern"}
	}

	top := scanNode

	if q.Unwind != nil {
		top = &PlanNode{
			Operator: OpUnwind,
			Details:  fmt.Sprintf("%s AS %s", formatExprBrief(q.Unwind.Expr), q.Unwind.Var),
			Children: []*PlanNode{top},
		}
	}

	if q.With != nil {
		top = &PlanNode{
			Operator: OpWith,
			Details:  formatWithItems(q.With),
			Children: []*PlanNode{top},
		}
	}

	if len(q.OrderBy) > 0 {
		top = &PlanNode{
			Operator: OpSort,
			Details:  formatOrderBy(q.OrderBy),
			Children: []*PlanNode{top},
		}
	}

	if (q.Skip != nil && *q.Skip > 0) || (q.Limit != nil && *q.Limit > 0) {
		detail := ""
		if q.Skip != nil && *q.Skip > 0 {
			detail = fmt.Sprintf("SKIP %d", *q.Skip)
		}
		if q.Limit != nil && *q.Limit > 0 {
			if detail != "" {
				detail += " "
			}
			detail += fmt.Sprintf("LIMIT %d", *q.Limit)
		}
		top = &PlanNode{
			Operator: OpLimit,
			Details:  detail,
			Children: []*PlanNode{top},
		}
	}

	top = &PlanNode{
		Operator: OpProduceResult,
		Details:  formatReturn(q.Return),
		Children: []*PlanNode{top},
	}

	return top
}

func buildNodeMatchPlan(q *CypherQuery, reader graphengine.GraphReader) *PlanNode {
	np := q.Match.Pattern.Nodes[0]
	varName := np.Variable
	if varName == "" {
		varName = "_n"
	}

	var scan *PlanNode

	if len(np.Labels) > 0 {
		scan = &PlanNode{
			Operator: OpLabelScan,
			Details:  fmt.Sprintf("%s:%s", varName, strings.Join(np.Labels, ":")),
		}
	} else if len(np.Props) > 0 {
		scan = &PlanNode{Operator: OpAllNodesScan, Details: varName}
	} else if q.Where != nil {
		if prop, _, ok := extractWhereEquality(q.Where, varName); ok {
			scan = &PlanNode{
				Operator: OpIndexSeek,
				Details:  fmt.Sprintf("%s.%s (from WHERE)", varName, prop),
			}
		} else {
			scan = &PlanNode{Operator: OpAllNodesScan, Details: varName}
		}
	} else {
		scan = &PlanNode{Operator: OpAllNodesScan, Details: varName}
	}

	if len(np.Props) > 0 || q.Where != nil {
		filterDetail := ""
		if len(np.Props) > 0 {
			filterDetail = fmt.Sprintf("inline props on %s", varName)
		}
		if q.Where != nil {
			if filterDetail != "" {
				filterDetail += " AND WHERE clause"
			} else {
				filterDetail = "WHERE clause"
			}
		}
		scan = &PlanNode{Operator: OpFilter, Details: filterDetail, Children: []*PlanNode{scan}}
	}

	return scan
}

func buildSingleHopPlan(q *CypherQuery, reader graphengine.GraphReader) *PlanNode {
	np := q.Match.Pattern.Nodes[0]
	rel := q.Match.Pattern.Rels[0]

	varA := np.Variable
	if varA == "" {
		varA = "_a"
	}

	source := buildNodeMatchPlan(&CypherQuery{
		Match: MatchClause{Pattern: Pattern{Nodes: []NodePattern{np}}},
	}, reader)

	expandDetail := fmt.Sprintf("(%s)-[:%s]->", varA, rel.Label)
	if rel.Dir == graphengine.Incoming {
		expandDetail = fmt.Sprintf("(%s)<-[:%s]-", varA, rel.Label)
	} else if rel.Dir == graphengine.Both {
		expandDetail = fmt.Sprintf("(%s)-[:%s]-", varA, rel.Label)
	}

	expand := &PlanNode{
		Operator: OpExpand,
		Details:  expandDetail,
		Children: []*PlanNode{source},
	}

	if q.Where != nil {
		expand = &PlanNode{Operator: OpFilter, Details: "WHERE clause", Children: []*PlanNode{expand}}
	}

	return expand
}

func buildVarLengthPlan(q *CypherQuery) *PlanNode {
	rel := q.Match.Pattern.Rels[0]
	maxHops := rel.MaxHops
	if maxHops < 0 {
		maxHops = defaultMaxDepth
	}
	detail := fmt.Sprintf("[:%s*%d..%d]", rel.Label, rel.MinHops, maxHops)
	return &PlanNode{
		Operator: OpExpandVarLen,
		Details:  detail,
		Children: []*PlanNode{
			{Operator: OpAllNodesScan, Details: "source"},
		},
	}
}

func formatOrderBy(items []OrderItem) string {
	var parts []string
	for _, item := range items {
		s := formatExprBrief(item.Expr)
		if item.Desc {
			s += " DESC"
		}
		parts = append(parts, s)
	}
	return strings.Join(parts, ", ")
}

func formatReturn(rc ReturnClause) string {
	var parts []string
	for _, item := range rc.Items {
		s := formatExprBrief(item.Expr)
		if item.Alias != "" {
			s += " AS " + item.Alias
		}
		parts = append(parts, s)
	}
	result := strings.Join(parts, ", ")
	if rc.Distinct {
		result = "DISTINCT " + result
	}
	return result
}

func formatWithItems(wc *WithClause) string {
	var parts []string
	for _, item := range wc.Items {
		s := formatExprBrief(item.Expr)
		if item.Alias != "" {
			s += " AS " + item.Alias
		}
		parts = append(parts, s)
	}
	result := strings.Join(parts, ", ")
	if wc.Distinct {
		result = "DISTINCT " + result
	}
	return result
}

func formatExprBrief(e Expression) string {
	switch e.Kind {
	case ExprVarRef:
		return e.Variable
	case ExprPropAccess:
		return e.Object + "." + e.Property
	case ExprFuncCall:
		return e.FuncName + "(...)"
	case ExprLiteral:
		return fmt.Sprintf("%v", e.LitValue)
	case ExprParam:
		return "$" + e.ParamName
	case ExprAggregation:
		if e.AggArg == nil {
			return aggFuncName(e.AggFunc) + "(*)"
		}
		return aggFuncName(e.AggFunc) + "(" + formatExprBrief(*e.AggArg) + ")"
	default:
		return "expr"
	}
}

// BuildExplainPlan builds a query plan without executing the query (EXPLAIN mode).
func BuildExplainPlan(q *CypherQuery, reader graphengine.GraphReader) *QueryPlan {
	return &QueryPlan{
		Root:    buildPlan(q, reader),
		Profile: false,
	}
}

// BuildProfilePlan builds a query plan with the query result attached (PROFILE mode).
func BuildProfilePlan(q *CypherQuery, reader graphengine.GraphReader, result *CypherResult) *QueryPlan {
	return &QueryPlan{
		Root:    buildPlan(q, reader),
		Profile: true,
		Result:  result,
	}
}

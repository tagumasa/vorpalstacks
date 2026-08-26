package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptunedata"
)

func (r *TestRunner) runNeptunedataGremlinBasicTests(tc *neptunedataContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("neptunedata", "Gremlin_ResetGraph", func() error {
		return tc.fastReset()
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteGremlinQuery_AddVertex", func() error {
		return tc.gremlin("g.addV('person').property('name','marko').property('age',29)")
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteGremlinQuery_AddMoreVertices", func() error {
		queries := []string{
			"g.addV('person').property('name','vadas').property('age',27)",
			"g.addV('person').property('name','josh').property('age',32)",
			"g.addV('software').property('name','lop').property('lang','java')",
		}
		for _, q := range queries {
			if err := tc.gremlin(q); err != nil {
				return err
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteGremlinQuery_AddEdges", func() error {
		queries := []string{
			"g.V().has('name','marko').addE('knows').to(g.V().has('name','vadas')).property('weight',0.5)",
			"g.V().has('name','marko').addE('knows').to(g.V().has('name','josh')).property('weight',1.0)",
			"g.V().has('name','marko').addE('created').to(g.V().has('name','lop')).property('weight',0.4)",
		}
		for _, q := range queries {
			if err := tc.gremlin(q); err != nil {
				return err
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteGremlinQuery_Count", func() error {
		return tc.gremlinContains("g.V().count()", "4")
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteGremlinQuery_Traversal", func() error {
		s, err := tc.gremlinResult("g.V().has('name','marko').out('knows').values('name')")
		if err != nil {
			return err
		}
		for _, name := range []string{"vadas", "josh"} {
			if !strings.Contains(s, name) {
				return fmt.Errorf("expected '%s' in gremlin traversal results, got %s", name, s)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteGremlinQuery_ValueMap", func() error {
		s, err := tc.gremlinResult("g.V().has('name','marko').valueMap()")
		if err != nil {
			return err
		}
		for _, key := range []string{"name", "age"} {
			if !strings.Contains(s, key) {
				return fmt.Errorf("expected key '%s' in valueMap result, got %s", key, s)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteGremlinQuery_HasLabel", func() error {
		return tc.gremlinContains("g.V().hasLabel('software').count()", "1")
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteGremlinQuery_EdgeCount", func() error {
		return tc.gremlinContains("g.E().count()", "3")
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteGremlinQuery_DropVertex", func() error {
		return tc.gremlin("g.V().has('name','lop').drop()")
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteGremlinQuery_VerifyDrop", func() error {
		return tc.gremlinContains("g.V().has('name','lop').count()", "0")
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteGremlinExplainQuery", func() error {
		resp, err := tc.client.ExecuteGremlinExplainQuery(tc.ctx, &neptunedata.ExecuteGremlinExplainQueryInput{
			GremlinQuery: aws.String("g.V().count()"),
		})
		if err != nil {
			return err
		}
		if resp.Output == nil {
			return fmt.Errorf("expected non-nil explain output")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteGremlinProfileQuery", func() error {
		resp, err := tc.client.ExecuteGremlinProfileQuery(tc.ctx, &neptunedata.ExecuteGremlinProfileQueryInput{
			GremlinQuery: aws.String("g.V().count()"),
		})
		if err != nil {
			return err
		}
		if resp.Output == nil {
			return fmt.Errorf("expected non-nil profile output")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "GetGremlinQueryStatus", func() error {
		resp, err := tc.client.ExecuteGremlinQuery(tc.ctx, &neptunedata.ExecuteGremlinQueryInput{
			GremlinQuery: aws.String("g.V().count()"),
		})
		if err != nil {
			return err
		}
		if resp.RequestId == nil {
			return fmt.Errorf("expected requestId from gremlin query")
		}

		statusResp, err := tc.client.GetGremlinQueryStatus(tc.ctx, &neptunedata.GetGremlinQueryStatusInput{
			QueryId: resp.RequestId,
		})
		if err != nil {
			return err
		}
		if statusResp.QueryId == nil || *statusResp.QueryId != *resp.RequestId {
			return fmt.Errorf("queryId mismatch: expected %s, got %v", *resp.RequestId, statusResp.QueryId)
		}
		return nil
	}))

	return results
}

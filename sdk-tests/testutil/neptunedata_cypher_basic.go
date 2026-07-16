package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptunedata"
	"github.com/aws/aws-sdk-go-v2/service/neptunedata/types"
)

func (r *TestRunner) runNeptunedataCypherBasicTests(tc *neptunedataContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("neptunedata", "ExecuteOpenCypherQuery_CreateNode", func() error {
		resp, err := tc.client.ExecuteOpenCypherQuery(tc.ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("CREATE (n:Person {name: 'marko', age: 29})"),
		})
		if err != nil {
			return err
		}
		if resp.Results == nil {
			return fmt.Errorf("expected non-nil results")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteOpenCypherQuery_CreateMoreNodes", func() error {
		resp, err := tc.client.ExecuteOpenCypherQuery(tc.ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("CREATE (n:Person {name: 'vadas', age: 27})"),
		})
		if err != nil {
			return err
		}
		if resp.Results == nil {
			return fmt.Errorf("expected non-nil results for first create")
		}
		resp, err = tc.client.ExecuteOpenCypherQuery(tc.ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("CREATE (n:Person {name: 'josh', age: 32})"),
		})
		if err != nil {
			return err
		}
		if resp.Results == nil {
			return fmt.Errorf("expected non-nil results for second create")
		}
		resp, err = tc.client.ExecuteOpenCypherQuery(tc.ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("CREATE (n:Software {name: 'lop', lang: 'java'})"),
		})
		if err != nil {
			return err
		}
		if resp.Results == nil {
			return fmt.Errorf("expected non-nil results for third create")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteOpenCypherQuery_CreateRelationships", func() error {
		resp, err := tc.client.ExecuteOpenCypherQuery(tc.ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("MATCH (a:Person {name: 'marko'}), (b:Person {name: 'vadas'}) CREATE (a)-[:KNOWS {weight: 0.5}]->(b)"),
		})
		if err != nil {
			return err
		}
		if resp.Results == nil {
			return fmt.Errorf("expected non-nil results for first edge")
		}
		resp, err = tc.client.ExecuteOpenCypherQuery(tc.ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("MATCH (a:Person {name: 'marko'}), (b:Person {name: 'josh'}) CREATE (a)-[:KNOWS {weight: 1.0}]->(b)"),
		})
		if err != nil {
			return err
		}
		if resp.Results == nil {
			return fmt.Errorf("expected non-nil results for second edge")
		}
		resp, err = tc.client.ExecuteOpenCypherQuery(tc.ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("MATCH (a:Person {name: 'marko'}), (b:Software {name: 'lop'}) CREATE (a)-[:CREATED {weight: 0.4}]->(b)"),
		})
		if err != nil {
			return err
		}
		if resp.Results == nil {
			return fmt.Errorf("expected non-nil results for third edge")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteOpenCypherQuery_MatchAllNodes", func() error {
		resp, err := tc.client.ExecuteOpenCypherQuery(tc.ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("MATCH (n) RETURN n.name ORDER BY n.name"),
		})
		if err != nil {
			return err
		}
		if resp.Results == nil {
			return fmt.Errorf("expected non-nil results")
		}
		s := marshalDoc(resp.Results)
		for _, name := range []string{`"marko"`, `"vadas"`, `"josh"`, `"lop"`} {
			if !strings.Contains(s, name) {
				return fmt.Errorf("expected node name %s in results, got %s", name, s)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteOpenCypherQuery_MatchByProperty", func() error {
		resp, err := tc.client.ExecuteOpenCypherQuery(tc.ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("MATCH (n:Person {name: 'marko'}) RETURN n.age"),
		})
		if err != nil {
			return err
		}
		if resp.Results == nil {
			return fmt.Errorf("expected non-nil results")
		}
		s := marshalDoc(resp.Results)
		if !strings.Contains(s, "29") {
			return fmt.Errorf("expected age 29 in results, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteOpenCypherQuery_Traversal", func() error {
		resp, err := tc.client.ExecuteOpenCypherQuery(tc.ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("MATCH (a:Person {name: 'marko'})-[:KNOWS]->(friend) RETURN friend.name"),
		})
		if err != nil {
			return err
		}
		if resp.Results == nil {
			return fmt.Errorf("expected non-nil results")
		}
		s := marshalDoc(resp.Results)
		for _, name := range []string{`"vadas"`, `"josh"`} {
			if !strings.Contains(s, name) {
				return fmt.Errorf("expected friend %s in traversal results, got %s", name, s)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteOpenCypherQuery_Aggregation", func() error {
		resp, err := tc.client.ExecuteOpenCypherQuery(tc.ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("MATCH (n:Person) RETURN count(n) AS cnt"),
		})
		if err != nil {
			return err
		}
		if resp.Results == nil {
			return fmt.Errorf("expected non-nil results")
		}
		s := marshalDoc(resp.Results)
		if !strings.Contains(s, "3") {
			return fmt.Errorf("expected count 3 in results, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteOpenCypherQuery_Parameters", func() error {
		resp, err := tc.client.ExecuteOpenCypherQuery(tc.ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("MATCH (n:Person {name: $name}) RETURN n.age"),
			Parameters:      aws.String(`{"name": "marko"}`),
		})
		if err != nil {
			return err
		}
		if resp.Results == nil {
			return fmt.Errorf("expected non-nil results")
		}
		s := marshalDoc(resp.Results)
		if !strings.Contains(s, "29") {
			return fmt.Errorf("expected age 29 from parameterised query, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteOpenCypherQuery_Delete", func() error {
		resp, err := tc.client.ExecuteOpenCypherQuery(tc.ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("MATCH (n:Software {name: 'lop'}) DETACH DELETE n"),
		})
		if err != nil {
			return err
		}
		if resp.Results == nil {
			return fmt.Errorf("expected non-nil results for delete")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteOpenCypherQuery_VerifyDelete", func() error {
		resp, err := tc.client.ExecuteOpenCypherQuery(tc.ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("MATCH (n:Software) RETURN count(n) AS cnt"),
		})
		if err != nil {
			return err
		}
		if resp.Results == nil {
			return fmt.Errorf("expected non-nil results")
		}
		s := marshalDoc(resp.Results)
		if !strings.Contains(s, "0") {
			return fmt.Errorf("expected count 0 after delete, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteOpenCypherExplainQuery", func() error {
		resp, err := tc.client.ExecuteOpenCypherExplainQuery(tc.ctx, &neptunedata.ExecuteOpenCypherExplainQueryInput{
			OpenCypherQuery: aws.String("MATCH (n) RETURN n LIMIT 1"),
			ExplainMode:     types.OpenCypherExplainModeStatic,
		})
		if err != nil {
			return err
		}
		if resp.Results == nil {
			return fmt.Errorf("expected non-nil explain results")
		}
		return nil
	}))

	return results
}

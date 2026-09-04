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
		return tc.cypher("CREATE (n:Person {name: 'marko', age: 29})")
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteOpenCypherQuery_CreateMoreNodes", func() error {
		if err := tc.cypher("CREATE (n:Person {name: 'vadas', age: 27})"); err != nil {
			return err
		}
		if err := tc.cypher("CREATE (n:Person {name: 'josh', age: 32})"); err != nil {
			return err
		}
		return tc.cypher("CREATE (n:Software {name: 'lop', lang: 'java'})")
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteOpenCypherQuery_CreateRelationships", func() error {
		if err := tc.cypher("MATCH (a:Person {name: 'marko'}), (b:Person {name: 'vadas'}) CREATE (a)-[:KNOWS {weight: 0.5}]->(b)"); err != nil {
			return err
		}
		if err := tc.cypher("MATCH (a:Person {name: 'marko'}), (b:Person {name: 'josh'}) CREATE (a)-[:KNOWS {weight: 1.0}]->(b)"); err != nil {
			return err
		}
		return tc.cypher("MATCH (a:Person {name: 'marko'}), (b:Software {name: 'lop'}) CREATE (a)-[:CREATED {weight: 0.4}]->(b)")
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteOpenCypherQuery_MatchAllNodes", func() error {
		return tc.cypherContains("MATCH (n) RETURN n.name ORDER BY n.name", `"marko"`, `"vadas"`, `"josh"`, `"lop"`)
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteOpenCypherQuery_MatchByProperty", func() error {
		return tc.cypherContains("MATCH (n:Person {name: 'marko'}) RETURN n.age", "29")
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteOpenCypherQuery_Traversal", func() error {
		return tc.cypherContains("MATCH (a:Person {name: 'marko'})-[:KNOWS]->(friend) RETURN friend.name", `"vadas"`, `"josh"`)
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteOpenCypherQuery_Aggregation", func() error {
		return tc.cypherContains("MATCH (n:Person) RETURN count(n) AS cnt", "3")
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteOpenCypherQuery_Parameters", func() error {
		resp, err := tc.client.ExecuteOpenCypherQuery(tc.ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("MATCH (n:Person {name: $name}) RETURN n.age"),
			Parameters:      aws.String(`{"name": "marko"}`),
		})
		if err != nil {
			return err
		}
		s := marshalDoc(resp.Results)
		if !strings.Contains(s, "29") {
			return fmt.Errorf("expected age 29 from parameterised query, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteOpenCypherQuery_Delete", func() error {
		return tc.cypher("MATCH (n:Software {name: 'lop'}) DETACH DELETE n")
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteOpenCypherQuery_VerifyDelete", func() error {
		return tc.cypherContains("MATCH (n:Software) RETURN count(n) AS cnt", "0")
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

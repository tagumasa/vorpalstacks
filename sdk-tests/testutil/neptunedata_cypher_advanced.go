package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptunedata"
	"github.com/aws/aws-sdk-go-v2/service/neptunedata/types"
)

func (r *TestRunner) runNeptunedataCypherAdvancedTests(tc *neptunedataContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("neptunedata", "Cypher_Advanced_Reset", func() error {
		return tc.fastReset()
	}))

	tc.cypher("CREATE (a:Person {name:'alice',age:25}), (b:Person {name:'bob',age:30}), (c:Person {name:'charlie',age:35}), (d:Person {name:'dave',age:40})")
	tc.cypher("MATCH (a:Person {name:'alice'}), (b:Person {name:'bob'}) CREATE (a)-[:KNOWS {weight:0.5}]->(b)")
	tc.cypher("MATCH (a:Person {name:'alice'}), (c:Person {name:'charlie'}) CREATE (a)-[:KNOWS {weight:1.0}]->(c)")
	tc.cypher("MATCH (b:Person {name:'bob'}), (d:Person {name:'dave'}) CREATE (b)-[:KNOWS {weight:0.7}]->(d)")
	tc.cypher("MATCH (c:Person {name:'charlie'}), (d:Person {name:'dave'}) CREATE (c)-[:KNOWS {weight:0.3}]->(d)")
	tc.cypher("MATCH (a:Person {name:'alice'}), (b:Person {name:'bob'}) CREATE (a)-[:WORKS_WITH]->(b)")
	tc.cypher("MATCH (b:Person {name:'bob'}), (c:Person {name:'charlie'}) CREATE (b)-[:WORKS_WITH]->(c)")

	results = append(results, r.RunTest("neptunedata", "Cypher_WhereClause", func() error {
		return tc.cypherContains("MATCH (n:Person) WHERE n.age > 30 RETURN n.name", `"charlie"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_WhereAND", func() error {
		return tc.cypherContains("MATCH (n:Person) WHERE n.age > 25 AND n.age < 40 RETURN n.name", `"bob"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_WhereOR", func() error {
		s, err := tc.cypherResult("MATCH (n:Person) WHERE n.name = 'alice' OR n.name = 'dave' RETURN n.name ORDER BY n.name")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"alice"`) || !strings.Contains(s, `"dave"`) {
			return fmt.Errorf("expected alice and dave, got %s", s)
		}
		return nil
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_WhereNOT", func() error {
		s, err := tc.cypherResult("MATCH (n:Person) WHERE NOT n.name = 'alice' AND NOT n.name = 'dave' RETURN n.name ORDER BY n.name")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"bob"`) || !strings.Contains(s, `"charlie"`) {
			return fmt.Errorf("expected bob and charlie, got %s", s)
		}
		return nil
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_WhereContains", func() error {
		return tc.cypherContains("MATCH (n:Person) WHERE n.name CONTAINS 'li' RETURN n.name", `"alice"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_WhereStartsWith", func() error {
		return tc.cypherContains("MATCH (n:Person) WHERE n.name STARTS WITH 'al' RETURN n.name", `"alice"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_WhereEndsWith", func() error {
		return tc.cypherContains("MATCH (n:Person) WHERE n.name ENDS WITH 'ie' RETURN n.name", `"charlie"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_WhereIN", func() error {
		s, err := tc.cypherResult("MATCH (n:Person) WHERE n.name IN ['alice', 'dave'] RETURN n.name ORDER BY n.name")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"alice"`) || !strings.Contains(s, `"dave"`) {
			return fmt.Errorf("expected alice and dave, got %s", s)
		}
		return nil
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_WhereIsNull", func() error {
		tc.cypher("CREATE (n:Person {name:'eve'})")
		return tc.cypherContains("MATCH (n:Person) WHERE n.age IS NULL RETURN n.name", `"eve"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_WhereIsNotNull", func() error {
		s, err := tc.cypherResult("MATCH (n:Person) WHERE n.age IS NOT NULL RETURN count(n) AS cnt")
		if err != nil {
			return err
		}
		if !strings.Contains(s, "4") {
			return fmt.Errorf("expected 4 persons with age, got %s", s)
		}
		return nil
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_ReturnDistinct", func() error {
		s, err := tc.cypherResult("MATCH (n:Person)-[:KNOWS]->(m) RETURN DISTINCT m.name ORDER BY m.name")
		if err != nil {
			return err
		}
		for _, name := range []string{`"bob"`, `"charlie"`, `"dave"`} {
			if !strings.Contains(s, name) {
				return fmt.Errorf("expected %s in DISTINCT results, got %s", name, s)
			}
		}
		return nil
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_OrderByDesc", func() error {
		s, err := tc.cypherResult("MATCH (n:Person) WHERE n.age IS NOT NULL RETURN n.name, n.age ORDER BY n.age DESC")
		if err != nil {
			return err
		}
		idx := strings.Index(s, `"dave"`)
		idx2 := strings.Index(s, `"alice"`)
		if idx > idx2 {
			return fmt.Errorf("expected dave before alice in DESC order, got %s", s)
		}
		return nil
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_Skip", func() error {
		return tc.cypherContains("MATCH (n:Person) WHERE n.age IS NOT NULL RETURN n.name ORDER BY n.age SKIP 2", `"charlie"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_Limit", func() error {
		return tc.cypherContains("MATCH (n:Person) RETURN n.name LIMIT 1", `"alice"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_SkipAndLimit", func() error {
		s, err := tc.cypherResult("MATCH (n:Person) WHERE n.age IS NOT NULL RETURN n.name ORDER BY n.age SKIP 1 LIMIT 1")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"bob"`) {
			return fmt.Errorf("expected bob at skip 1 limit 1, got %s", s)
		}
		return nil
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_AggregationSum", func() error {
		return tc.cypherContains("MATCH (n:Person) WHERE n.age IS NOT NULL RETURN sum(n.age) AS total", `130`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_AggregationAvg", func() error {
		return tc.cypherContains("MATCH (n:Person) WHERE n.age IS NOT NULL RETURN avg(n.age) AS avg", `32.5`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_AggregationMin", func() error {
		return tc.cypherContains("MATCH (n:Person) WHERE n.age IS NOT NULL RETURN min(n.age) AS mn", `25`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_AggregationMax", func() error {
		return tc.cypherContains("MATCH (n:Person) WHERE n.age IS NOT NULL RETURN max(n.age) AS mx", `40`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_AggregationCollect", func() error {
		s, err := tc.cypherResult("MATCH (n:Person) WHERE n.age IS NOT NULL RETURN collect(n.name) AS names")
		if err != nil {
			return err
		}
		for _, name := range []string{`"alice"`, `"bob"`, `"charlie"`, `"dave"`} {
			if !strings.Contains(s, name) {
				return fmt.Errorf("expected %s in collect, got %s", name, s)
			}
		}
		return nil
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_CountDistinct", func() error {
		return tc.cypherContains("MATCH (n:Person)-[:KNOWS]->(m) RETURN count(DISTINCT m.name) AS cnt", `3`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_SetProperty", func() error {
		if err := tc.cypher("MATCH (n:Person {name:'alice'}) SET n.age = 26"); err != nil {
			return err
		}
		return tc.cypherContains("MATCH (n:Person {name:'alice'}) RETURN n.age", `26`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_SetMergeProperties", func() error {
		if err := tc.cypher("MATCH (n:Person {name:'alice'}) SET n += {city:'NYC',active:true}"); err != nil {
			return err
		}
		return tc.cypherContains("MATCH (n:Person {name:'alice'}) RETURN n.city", `"NYC"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_SetLabel", func() error {
		if err := tc.cypher("MATCH (n:Person {name:'alice'}) SET n:Employee"); err != nil {
			return err
		}
		return tc.cypherContains("MATCH (n:Employee) RETURN n.name", `"alice"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_RemoveLabel", func() error {
		if err := tc.cypher("MATCH (n:Person {name:'alice'}) REMOVE n:Employee"); err != nil {
			return err
		}
		s, err := tc.cypherResult("MATCH (n:Employee) RETURN count(n) AS cnt")
		if err != nil {
			return err
		}
		if !strings.Contains(s, "0") {
			return fmt.Errorf("expected 0 employees after REMOVE, got %s", s)
		}
		return nil
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_RemoveProperty", func() error {
		if err := tc.cypher("MATCH (n:Person {name:'alice'}) REMOVE n.city"); err != nil {
			return err
		}
		s, err := tc.cypherResult("MATCH (n:Person {name:'alice'}) RETURN n.city")
		if err != nil {
			return err
		}
		if strings.Contains(s, `"NYC"`) {
			return fmt.Errorf("expected city removed, got %s", s)
		}
		return nil
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_DeleteNonDetach", func() error {
		tc.cypher("CREATE (n:Temp {name:'temp_node'})")
		if err := tc.cypher("MATCH (n:Temp) DELETE n"); err != nil {
			return err
		}
		return tc.cypherContains("MATCH (n:Temp) RETURN count(n) AS cnt", `0`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_OptionalMatch", func() error {
		s, err := tc.cypherResult("MATCH (a:Person {name:'alice'}) OPTIONAL MATCH (a)-[:WORKS_WITH]->(c) RETURN a.name, c.name ORDER BY c.name")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"alice"`) {
			return fmt.Errorf("expected alice, got %s", s)
		}
		if !strings.Contains(s, `"bob"`) {
			return fmt.Errorf("expected bob (WORKS_WITH target), got %s", s)
		}
		return nil
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_OptionalMatchNull", func() error {
		s, err := tc.cypherResult("MATCH (a:Person {name:'dave'}) OPTIONAL MATCH (a)-[:WORKS_WITH]->(c) RETURN a.name, c.name")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"dave"`) {
			return fmt.Errorf("expected dave, got %s", s)
		}
		return nil
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_WithClause", func() error {
		return tc.cypherContains("MATCH (n:Person) WHERE n.age IS NOT NULL WITH n.name AS name ORDER BY name LIMIT 2 RETURN name", `"alice"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_WithAggregation", func() error {
		return tc.cypherContains("MATCH (n:Person)-[:KNOWS]->(m) WITH m.name AS friend, count(*) AS cnt RETURN friend, cnt ORDER BY cnt DESC", `"bob"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_Unwind", func() error {
		if err := tc.cypher("UNWIND [10,20,30] AS x CREATE (n:Number {val:x})"); err != nil {
			return err
		}
		return tc.cypherContains("MATCH (n:Number) RETURN sum(n.val) AS total", `60`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_MultiHop", func() error {
		return tc.cypherContains("MATCH (a:Person {name:'alice'})-[:KNOWS]->(b)-[:KNOWS]->(c) RETURN c.name", `"dave"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_VarLengthPath", func() error {
		return tc.cypherContains("MATCH (a:Person {name:'alice'})-[:KNOWS*1..2]->(b) RETURN DISTINCT b.name ORDER BY b.name", `"dave"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_CommaPatterns", func() error {
		s, err := tc.cypherResult("MATCH (a:Person {name:'alice'}), (b:Person {name:'bob'}) RETURN a.name, b.name")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"alice"`) || !strings.Contains(s, `"bob"`) {
			return fmt.Errorf("expected both alice and bob, got %s", s)
		}
		return nil
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_MultiMatch", func() error {
		return tc.cypherContains("MATCH (a:Person {name:'alice'})-[:KNOWS]->(b) MATCH (b)-[:KNOWS]->(c) RETURN c.name", `"dave"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_MergeCreate", func() error {
		if err := tc.cypher("MERGE (n:Animal {name:'rex'})"); err != nil {
			return err
		}
		return tc.cypherContains("MATCH (n:Animal) RETURN n.name", `"rex"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_MergeIdempotent", func() error {
		tc.cypher("MATCH (n:Animal {name:'rex'}) DELETE n")
		if err := tc.cypher("MERGE (n:Animal {name:'rex'}) ON CREATE SET n.species = 'dog'"); err != nil {
			return err
		}
		return tc.cypherContains("MATCH (n:Animal {name:'rex'}) RETURN n.species", `"dog"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_MergeOnMatchSet", func() error {
		if err := tc.cypher("MERGE (n:Animal {name:'rex'}) ON MATCH SET n.visits = 1"); err != nil {
			return err
		}
		return tc.cypherContains("MATCH (n:Animal {name:'rex'}) RETURN n.visits", `1`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_TypeFunction", func() error {
		return tc.cypherContains("MATCH (a:Person {name:'alice'})-[r:KNOWS]->(b) RETURN type(r)", `"KNOWS"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_IdFunction", func() error {
		s, err := tc.cypherResult("MATCH (n:Person {name:'alice'}) RETURN id(n) AS nid")
		if err != nil {
			return err
		}
		if s == "" {
			return fmt.Errorf("expected non-empty id result")
		}
		return nil
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_CoalesceFunction", func() error {
		return tc.cypherContains("MATCH (n:Person {name:'eve'}) RETURN coalesce(n.age, 0) AS age", `0`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_CaseExpression", func() error {
		return tc.cypherContains("MATCH (n:Person {name:'alice'}) RETURN CASE WHEN n.age > 30 THEN 'senior' ELSE 'junior' END AS category", `"junior"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_Profile", func() error {
		resp, err := tc.client.ExecuteOpenCypherExplainQuery(tc.ctx, &neptunedata.ExecuteOpenCypherExplainQueryInput{
			OpenCypherQuery: aws.String("PROFILE MATCH (n) RETURN n LIMIT 1"),
			ExplainMode:     types.OpenCypherExplainModeDetails,
		})
		if err != nil {
			return err
		}
		if resp.Results == nil {
			return fmt.Errorf("expected non-nil profile results")
		}
		return nil
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_MatchEdgeReturn", func() error {
		return tc.cypherContains("MATCH (a:Person {name:'alice'})-[r:KNOWS]->(b) RETURN r.weight", `0.5`)
	}))
	results = append(results, r.RunTest("neptunedata", "Cypher_IncomingDirection", func() error {
		return tc.cypherContains("MATCH (b:Person {name:'bob'})<-[:KNOWS]-(a) RETURN a.name", `"alice"`)
	}))

	return results
}

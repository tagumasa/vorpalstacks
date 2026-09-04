package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptunedata"
)

func (r *TestRunner) runNeptunedataGremlinAdvancedTests(tc *neptunedataContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("neptunedata", "Gremlin_Advanced_Reset", func() error {
		return tc.fastReset()
	}))

	tc.gremlin("g.addV('person').property('name','alice').property('age',25)")
	tc.gremlin("g.addV('person').property('name','bob').property('age',30)")
	tc.gremlin("g.addV('person').property('name','charlie').property('age',35)")
	tc.gremlin("g.addV('person').property('name','dave').property('age',40)")
	tc.gremlin("g.V().has('name','alice').addE('knows').to(g.V().has('name','bob')).property('weight',0.5)")
	tc.gremlin("g.V().has('name','alice').addE('knows').to(g.V().has('name','charlie')).property('weight',1.0)")
	tc.gremlin("g.V().has('name','bob').addE('knows').to(g.V().has('name','dave')).property('weight',0.7)")
	tc.gremlin("g.V().has('name','charlie').addE('knows').to(g.V().has('name','dave')).property('weight',0.3)")

	results = append(results, r.RunTest("neptunedata", "Gremlin_In", func() error {
		return tc.gremlinContains("g.V().has('name','bob').in('knows').values('name')", `"alice"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Both", func() error {
		return tc.gremlinContains("g.V().has('name','bob').both('knows').values('name').order()", `"alice"`, `"dave"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_OutE", func() error {
		return tc.gremlinContains("g.V().has('name','alice').outE('knows').count()", "2")
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_InE", func() error {
		return tc.gremlinContains("g.V().has('name','bob').inE('knows').count()", "1")
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_BothE", func() error {
		return tc.gremlinContains("g.V().has('name','bob').bothE().count()", "2")
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_OutV", func() error {
		return tc.gremlinContains("g.E().hasLabel('knows').outV().has('name','alice').count()", `2`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_InV", func() error {
		return tc.gremlinContains("g.E().hasLabel('knows').inV().has('name','bob').count()", `1`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_OtherV", func() error {
		return tc.gremlinContains("g.V().has('name','alice').outE('knows').otherV().has('name','bob').count()", `1`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_HasPredicate", func() error {
		return tc.gremlinContains("g.V().has('age',P.gt(30)).values('name')", `"charlie"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_HasGte", func() error {
		return tc.gremlinContains("g.V().has('age',P.gte(30)).count()", "3")
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_HasLt", func() error {
		return tc.gremlinContains("g.V().has('age',P.lt(30)).values('name')", `"alice"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_HasNeq", func() error {
		return tc.gremlinContains("g.V().has('name',P.neq('alice')).count()", "3")
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_HasBetween", func() error {
		return tc.gremlinContains("g.V().has('age',P.between(25,35)).count()", "3")
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_HasWithin", func() error {
		return tc.gremlinContains("g.V().has('name',P.within('alice','dave')).values('name').order()", `"alice"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_HasNot", func() error {
		_, _ = tc.client.ExecuteOpenCypherQuery(tc.ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("CREATE (n:Person {name:'frank'})"),
		})
		return tc.gremlinContains("g.V().hasNot('age').values('name')", `"frank"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_StartingWith", func() error {
		return tc.gremlinContains("g.V().has('name',P.startingWith('al')).values('name')", `"alice"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Containing", func() error {
		return tc.gremlinContains("g.V().has('name',P.containing('li')).values('name')", `"alice"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_EndingWith", func() error {
		return tc.gremlinContains("g.V().has('name',P.endingWith('ie')).values('name')", `"charlie"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Id", func() error {
		s, err := tc.gremlinResult("g.V().has('name','alice').id()")
		if err != nil {
			return err
		}
		if s == "" {
			return fmt.Errorf("expected non-empty id result")
		}
		return nil
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Label", func() error {
		return tc.gremlinContains("g.V().has('name','alice').label()", `"person"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Keys", func() error {
		return tc.gremlinContains("g.V().has('name','alice').keys()", `"name"`, `"age"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Properties", func() error {
		return tc.gremlinContains("g.V().has('name','alice').properties('name')", `"name"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_PropertyMap", func() error {
		return tc.gremlinContains("g.V().has('name','alice').propertyMap('name','age')", `"alice"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_PropertyMutation", func() error {
		tc.gremlin("g.V().has('name','alice').property('city','NYC')")
		return tc.gremlinContains("g.V().has('name','alice').values('city')", `"NYC"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_AsSelect", func() error {
		return tc.gremlinContains("g.V().has('name','alice').as('a').out('knows').as('b').select('a','b').by('name')", `"alice"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_WhereTraversal", func() error {
		return tc.gremlinContains("g.V().where(out().count().is(gt(0))).values('name').dedup()", `"alice"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Filter", func() error {
		return tc.gremlinContains("g.V().filter(out('knows').count().is(gt(1))).values('name')", `"alice"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Limit", func() error {
		return tc.gremlinContains("g.V().hasLabel('person').limit(2).count()", `2`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Range", func() error {
		return tc.gremlinContains("g.V().hasLabel('person').range(1,3).count()", `2`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Skip", func() error {
		return tc.gremlinContains("g.V().hasLabel('person').order().by('name').skip(2).values('name')", `"charlie"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Dedup", func() error {
		tc.gremlin("g.addV('dup').property('name','dup1')")
		tc.gremlin("g.addV('dup').property('name','dup1')")
		return tc.gremlinContains("g.V().hasLabel('dup').dedup().by('name').count()", "1")
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_OrderByAsc", func() error {
		return tc.gremlinContains("g.V().hasLabel('person').order().by('name',asc).values('name').limit(1)", `"alice"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_OrderByDesc", func() error {
		return tc.gremlinContains("g.V().hasLabel('person').order().by('name',desc).values('name').limit(1)", `"dave"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_GroupCount", func() error {
		return tc.gremlinContains("g.V().hasLabel('person').groupCount().by('name')", `"alice"`, `"bob"`, `"charlie"`, `"dave"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_GroupBy", func() error {
		return tc.gremlinContains("g.V().hasLabel('person').group().by('name').by('age')", `"alice"`, `"bob"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Path", func() error {
		return tc.gremlinContains("g.V().has('name','alice').out('knows').limit(1).path()", `"alice"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Union", func() error {
		return tc.gremlinContains("g.V().has('name','bob').union(in('knows'),out('knows')).values('name').dedup().order()", `"alice"`, `"dave"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Coalesce", func() error {
		return tc.gremlinContains("g.V().has('name','bob').coalesce(out('nonexistent'),in('knows')).values('name')", `"alice"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Choose", func() error {
		return tc.gremlinContains("g.V().has('name','alice').choose(has('age'),values('age'),values('name'))", `25`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Optional", func() error {
		s, err := tc.gremlinResult("g.V().has('name','dave').optional(out('knows')).values('name')")
		if err != nil {
			return err
		}
		if strings.Contains(s, `"charlie"`) {
			return fmt.Errorf("optional should return original when no outgoing knows, got %s", s)
		}
		return nil
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_RepeatTimes", func() error {
		return tc.gremlinContains("g.V().has('name','alice').repeat(out('knows')).times(2).values('name')", `"dave"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_RepeatEmit", func() error {
		return tc.gremlinContains("g.V().has('name','alice').repeat(out('knows')).emit().values('name').dedup().order()", `"bob"`, `"charlie"`, `"dave"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Fold", func() error {
		return tc.gremlinContains("g.V().hasLabel('person').values('name').fold()", `"alice"`, `"bob"`, `"charlie"`, `"dave"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Unfold", func() error {
		return tc.gremlinContains("g.V().hasLabel('person').values('name').fold().unfold().count()", `4`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Constant", func() error {
		return tc.gremlinContains("g.V().has('name','alice').constant('hello').limit(1)", `"hello"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Is", func() error {
		return tc.gremlinContains("g.V().hasLabel('person').values('age').is(P.gt(30))", `35`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Not", func() error {
		return tc.gremlinContains("g.V().hasLabel('person').not(out('knows')).values('name')", `"dave"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_And", func() error {
		return tc.gremlinContains("g.V().and(out('knows'),in('knows')).values('name')", `"bob"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Or", func() error {
		return tc.gremlinContains("g.V().or(has('name','alice'),has('name','dave')).values('name').order()", `"alice"`, `"dave"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Mean", func() error {
		return tc.gremlinContains("g.V().hasLabel('person').values('age').mean()", `32.5`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Sum", func() error {
		return tc.gremlinContains("g.V().hasLabel('person').values('age').sum()", `130`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Min", func() error {
		return tc.gremlinContains("g.V().hasLabel('person').values('age').min()", `25`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Max", func() error {
		return tc.gremlinContains("g.V().hasLabel('person').values('age').max()", `40`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Tail", func() error {
		s, err := tc.gremlinResult("g.V().hasLabel('person').order().by('name').tail(1).values('name')")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"dave"`) {
			return fmt.Errorf("expected dave as tail, got %s", s)
		}
		return nil
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Inject", func() error {
		return tc.gremlinContains("g.inject('a','b').count()", `2`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_MergeV", func() error {
		tc.gremlin("g.mergeV([~label:'item',name:'widget',price:10])")
		return tc.gremlinContains("g.V().has('name','widget').values('price')", `10`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_MergeVIdempotent", func() error {
		tc.gremlin("g.mergeV([~label:'item',name:'widget',price:10])")
		return tc.gremlinContains("g.V().has('name','widget').count()", "1")
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_ElementMap", func() error {
		return tc.gremlinContains("g.V().has('name','alice').elementMap()", `"alice"`, `"person"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_Project", func() error {
		return tc.gremlinContains("g.V().as('a').out('knows').as('b').project('from','to').by('name').limit(1)", `"from"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_CountLocal", func() error {
		return tc.gremlinContains("g.V().hasLabel('person').values('name').fold().count(local)", `4`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_MultiHop", func() error {
		return tc.gremlinContains("g.V().has('name','alice').out().out().values('name')", `"dave"`)
	}))
	results = append(results, r.RunTest("neptunedata", "Gremlin_SimplePath", func() error {
		return tc.gremlinContains("g.V().has('name','alice').repeat(__.out('knows')).times(2).simplePath().values('name').dedup()", `"dave"`)
	}))

	return results
}

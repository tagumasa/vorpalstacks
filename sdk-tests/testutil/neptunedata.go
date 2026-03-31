package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptunedata"
	"github.com/aws/aws-sdk-go-v2/service/neptunedata/types"
	"vorpalstacks-sdk-tests/config"
)

type smithyUnmarshaler interface {
	UnmarshalSmithyDocument(interface{}) error
}

func marshalDoc(v interface{}) string {
	if u, ok := v.(smithyUnmarshaler); ok {
		var target interface{}
		if err := u.UnmarshalSmithyDocument(&target); err == nil {
			v = target
		}
	}
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}

func containsDoc(v interface{}, substr string) bool {
	return strings.Contains(marshalDoc(v), substr)
}

func (r *TestRunner) RunNeptunedataTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "neptunedata",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := neptunedata.NewFromConfig(cfg)
	ctx := context.Background()

	results = append(results, r.RunTest("neptunedata", "GetEngineStatus", func() error {
		resp, err := client.GetEngineStatus(ctx, &neptunedata.GetEngineStatusInput{})
		if err != nil {
			return err
		}
		if resp.Status == nil || *resp.Status != "healthy" {
			return fmt.Errorf("expected status=healthy, got %v", resp.Status)
		}
		if resp.StartTime == nil || *resp.StartTime == "" {
			return fmt.Errorf("expected non-empty startTime")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "GetEngineStatus_GremlinVersion", func() error {
		resp, err := client.GetEngineStatus(ctx, &neptunedata.GetEngineStatusInput{})
		if err != nil {
			return err
		}
		if resp.Gremlin == nil || resp.Gremlin.Version == nil || *resp.Gremlin.Version == "" {
			return fmt.Errorf("expected non-empty gremlin version")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "GetEngineStatus_OpenCypherVersion", func() error {
		resp, err := client.GetEngineStatus(ctx, &neptunedata.GetEngineStatusInput{})
		if err != nil {
			return err
		}
		if resp.Opencypher == nil || resp.Opencypher.Version == nil || *resp.Opencypher.Version == "" {
			return fmt.Errorf("expected non-empty opencypher version")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "GetEngineStatus_Role", func() error {
		resp, err := client.GetEngineStatus(ctx, &neptunedata.GetEngineStatusInput{})
		if err != nil {
			return err
		}
		if resp.Role == nil || (*resp.Role != "writer" && *resp.Role != "reader") {
			return fmt.Errorf("expected role=writer|reader, got %v", resp.Role)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteFastReset_Initiate", func() error {
		resp, err := client.ExecuteFastReset(ctx, &neptunedata.ExecuteFastResetInput{
			Action: types.ActionInitializeReset,
		})
		if err != nil {
			return err
		}
		if resp.Payload == nil || resp.Payload.Token == nil || *resp.Payload.Token == "" {
			return fmt.Errorf("expected non-empty fast reset token")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteFastReset_Perform", func() error {
		initResp, err := client.ExecuteFastReset(ctx, &neptunedata.ExecuteFastResetInput{
			Action: types.ActionInitializeReset,
		})
		if err != nil {
			return err
		}
		if initResp.Payload == nil || initResp.Payload.Token == nil {
			return fmt.Errorf("expected non-empty token from initiateDatabaseReset")
		}

		performResp, err := client.ExecuteFastReset(ctx, &neptunedata.ExecuteFastResetInput{
			Action: types.ActionPerformReset,
			Token:  initResp.Payload.Token,
		})
		if err != nil {
			return err
		}
		if performResp.Status == nil || *performResp.Status == "" {
			return fmt.Errorf("expected non-empty status from performDatabaseReset")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteOpenCypherQuery_CreateNode", func() error {
		_, err := client.ExecuteOpenCypherQuery(ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("CREATE (n:Person {name: 'marko', age: 29})"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteOpenCypherQuery_CreateMoreNodes", func() error {
		_, err := client.ExecuteOpenCypherQuery(ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("CREATE (n:Person {name: 'vadas', age: 27})"),
		})
		if err != nil {
			return err
		}
		_, err = client.ExecuteOpenCypherQuery(ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("CREATE (n:Person {name: 'josh', age: 32})"),
		})
		if err != nil {
			return err
		}
		_, err = client.ExecuteOpenCypherQuery(ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("CREATE (n:Software {name: 'lop', lang: 'java'})"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteOpenCypherQuery_CreateRelationships", func() error {
		_, err := client.ExecuteOpenCypherQuery(ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("MATCH (a:Person {name: 'marko'}), (b:Person {name: 'vadas'}) CREATE (a)-[:KNOWS {weight: 0.5}]->(b)"),
		})
		if err != nil {
			return err
		}
		_, err = client.ExecuteOpenCypherQuery(ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("MATCH (a:Person {name: 'marko'}), (b:Person {name: 'josh'}) CREATE (a)-[:KNOWS {weight: 1.0}]->(b)"),
		})
		if err != nil {
			return err
		}
		_, err = client.ExecuteOpenCypherQuery(ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("MATCH (a:Person {name: 'marko'}), (b:Software {name: 'lop'}) CREATE (a)-[:CREATED {weight: 0.4}]->(b)"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteOpenCypherQuery_MatchAllNodes", func() error {
		resp, err := client.ExecuteOpenCypherQuery(ctx, &neptunedata.ExecuteOpenCypherQueryInput{
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
		resp, err := client.ExecuteOpenCypherQuery(ctx, &neptunedata.ExecuteOpenCypherQueryInput{
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
		resp, err := client.ExecuteOpenCypherQuery(ctx, &neptunedata.ExecuteOpenCypherQueryInput{
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
		resp, err := client.ExecuteOpenCypherQuery(ctx, &neptunedata.ExecuteOpenCypherQueryInput{
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
		resp, err := client.ExecuteOpenCypherQuery(ctx, &neptunedata.ExecuteOpenCypherQueryInput{
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
		_, err := client.ExecuteOpenCypherQuery(ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("MATCH (n:Software {name: 'lop'}) DETACH DELETE n"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteOpenCypherQuery_VerifyDelete", func() error {
		resp, err := client.ExecuteOpenCypherQuery(ctx, &neptunedata.ExecuteOpenCypherQueryInput{
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
		resp, err := client.ExecuteOpenCypherExplainQuery(ctx, &neptunedata.ExecuteOpenCypherExplainQueryInput{
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

	results = append(results, r.RunTest("neptunedata", "Gremlin_ResetGraph", func() error {
		initResp, err := client.ExecuteFastReset(ctx, &neptunedata.ExecuteFastResetInput{
			Action: types.ActionInitializeReset,
		})
		if err != nil {
			return err
		}
		if initResp.Payload == nil || initResp.Payload.Token == nil {
			return fmt.Errorf("expected token from initiateDatabaseReset")
		}
		_, err = client.ExecuteFastReset(ctx, &neptunedata.ExecuteFastResetInput{
			Action: types.ActionPerformReset,
			Token:  initResp.Payload.Token,
		})
		return err
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteGremlinQuery_AddVertex", func() error {
		resp, err := client.ExecuteGremlinQuery(ctx, &neptunedata.ExecuteGremlinQueryInput{
			GremlinQuery: aws.String("g.addV('person').property('name','marko').property('age',29)"),
		})
		if err != nil {
			return err
		}
		if resp.RequestId == nil || *resp.RequestId == "" {
			return fmt.Errorf("expected non-empty requestId")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteGremlinQuery_AddMoreVertices", func() error {
		queries := []string{
			"g.addV('person').property('name','vadas').property('age',27)",
			"g.addV('person').property('name','josh').property('age',32)",
			"g.addV('software').property('name','lop').property('lang','java')",
		}
		for _, q := range queries {
			_, err := client.ExecuteGremlinQuery(ctx, &neptunedata.ExecuteGremlinQueryInput{
				GremlinQuery: aws.String(q),
			})
			if err != nil {
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
			_, err := client.ExecuteGremlinQuery(ctx, &neptunedata.ExecuteGremlinQueryInput{
				GremlinQuery: aws.String(q),
			})
			if err != nil {
				return err
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteGremlinQuery_Count", func() error {
		resp, err := client.ExecuteGremlinQuery(ctx, &neptunedata.ExecuteGremlinQueryInput{
			GremlinQuery: aws.String("g.V().count()"),
		})
		if err != nil {
			return err
		}
		s := marshalDoc(resp.Result)
		if !strings.Contains(s, "4") {
			return fmt.Errorf("expected count 4 in gremlin result, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteGremlinQuery_Traversal", func() error {
		resp, err := client.ExecuteGremlinQuery(ctx, &neptunedata.ExecuteGremlinQueryInput{
			GremlinQuery: aws.String("g.V().has('name','marko').out('knows').values('name')"),
		})
		if err != nil {
			return err
		}
		s := marshalDoc(resp.Result)
		for _, name := range []string{"vadas", "josh"} {
			if !strings.Contains(s, name) {
				return fmt.Errorf("expected '%s' in gremlin traversal results, got %s", name, s)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteGremlinQuery_ValueMap", func() error {
		resp, err := client.ExecuteGremlinQuery(ctx, &neptunedata.ExecuteGremlinQueryInput{
			GremlinQuery: aws.String("g.V().has('name','marko').valueMap()"),
		})
		if err != nil {
			return err
		}
		s := marshalDoc(resp.Result)
		for _, key := range []string{"name", "age"} {
			if !strings.Contains(s, key) {
				return fmt.Errorf("expected key '%s' in valueMap result, got %s", key, s)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteGremlinQuery_HasLabel", func() error {
		resp, err := client.ExecuteGremlinQuery(ctx, &neptunedata.ExecuteGremlinQueryInput{
			GremlinQuery: aws.String("g.V().hasLabel('software').count()"),
		})
		if err != nil {
			return err
		}
		s := marshalDoc(resp.Result)
		if !strings.Contains(s, "1") {
			return fmt.Errorf("expected count 1 for software label, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteGremlinQuery_EdgeCount", func() error {
		resp, err := client.ExecuteGremlinQuery(ctx, &neptunedata.ExecuteGremlinQueryInput{
			GremlinQuery: aws.String("g.E().count()"),
		})
		if err != nil {
			return err
		}
		s := marshalDoc(resp.Result)
		if !strings.Contains(s, "3") {
			return fmt.Errorf("expected count 3 for edges, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteGremlinQuery_DropVertex", func() error {
		_, err := client.ExecuteGremlinQuery(ctx, &neptunedata.ExecuteGremlinQueryInput{
			GremlinQuery: aws.String("g.V().has('name','lop').drop()"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteGremlinQuery_VerifyDrop", func() error {
		resp, err := client.ExecuteGremlinQuery(ctx, &neptunedata.ExecuteGremlinQueryInput{
			GremlinQuery: aws.String("g.V().has('name','lop').count()"),
		})
		if err != nil {
			return err
		}
		s := marshalDoc(resp.Result)
		if !strings.Contains(s, "0") {
			return fmt.Errorf("expected count 0 after drop, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ExecuteGremlinExplainQuery", func() error {
		resp, err := client.ExecuteGremlinExplainQuery(ctx, &neptunedata.ExecuteGremlinExplainQueryInput{
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
		resp, err := client.ExecuteGremlinProfileQuery(ctx, &neptunedata.ExecuteGremlinProfileQueryInput{
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
		resp, err := client.ExecuteGremlinQuery(ctx, &neptunedata.ExecuteGremlinQueryInput{
			GremlinQuery: aws.String("g.V().count()"),
		})
		if err != nil {
			return err
		}
		if resp.RequestId == nil {
			return fmt.Errorf("expected requestId from gremlin query")
		}

		statusResp, err := client.GetGremlinQueryStatus(ctx, &neptunedata.GetGremlinQueryStatusInput{
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

	results = append(results, r.RunTest("neptunedata", "GetOpenCypherQueryStatus", func() error {
		_, err := client.ExecuteOpenCypherQuery(ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("MATCH (n) RETURN count(n)"),
		})
		return err
	}))

	results = append(results, r.RunTest("neptunedata", "GetPropertygraphStatistics", func() error {
		resp, err := client.GetPropertygraphStatistics(ctx, &neptunedata.GetPropertygraphStatisticsInput{})
		if err != nil {
			return err
		}
		if resp.Status == nil {
			return fmt.Errorf("expected non-nil status")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "GetPropertygraphSummary", func() error {
		resp, err := client.GetPropertygraphSummary(ctx, &neptunedata.GetPropertygraphSummaryInput{
			Mode: types.GraphSummaryTypeBasic,
		})
		if err != nil {
			return err
		}
		if resp.StatusCode == nil {
			return fmt.Errorf("expected non-nil statusCode")
		}
		return nil
	}))

	var loaderJobID string

	results = append(results, r.RunTest("neptunedata", "StartLoaderJob", func() error {
		resp, err := client.StartLoaderJob(ctx, &neptunedata.StartLoaderJobInput{
			Source:         aws.String("s3://test-bucket/data"),
			Format:         types.FormatCsv,
			IamRoleArn:     aws.String("arn:aws:iam::000000000000:role/NeptuneLoadRole"),
			S3BucketRegion: types.S3BucketRegionUsEast1,
		})
		if err != nil {
			return err
		}
		if resp.Payload == nil {
			return fmt.Errorf("expected non-nil loader job payload")
		}
		data := marshalDoc(resp.Payload)
		var payloadMap map[string]interface{}
		if err := json.Unmarshal([]byte(data), &payloadMap); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}
		if lid, ok := payloadMap["loadId"].(string); ok && lid != "" {
			loaderJobID = lid
		} else {
			return fmt.Errorf("expected loadId in payload, got %s", data)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "GetLoaderJobStatus", func() error {
		if loaderJobID == "" {
			return fmt.Errorf("no loader job ID from StartLoaderJob")
		}
		resp, err := client.GetLoaderJobStatus(ctx, &neptunedata.GetLoaderJobStatusInput{
			LoadId: aws.String(loaderJobID),
		})
		if err != nil {
			return err
		}
		if resp.Payload == nil {
			return fmt.Errorf("expected non-nil loader job status payload")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ListLoaderJobs", func() error {
		resp, err := client.ListLoaderJobs(ctx, &neptunedata.ListLoaderJobsInput{})
		if err != nil {
			return err
		}
		if resp.Payload == nil {
			return fmt.Errorf("expected non-nil list loader jobs payload")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "CancelLoaderJob", func() error {
		if loaderJobID == "" {
			return fmt.Errorf("no loader job ID from StartLoaderJob")
		}
		_, err := client.CancelLoaderJob(ctx, &neptunedata.CancelLoaderJobInput{
			LoadId: aws.String(loaderJobID),
		})
		return err
	}))

	results = append(results, r.RunTest("neptunedata", "GetSparqlStatistics_Unsupported", func() error {
		_, err := client.GetSparqlStatistics(ctx, &neptunedata.GetSparqlStatisticsInput{})
		if err == nil {
			return fmt.Errorf("expected error for unsupported SPARQL statistics")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "GetRDFGraphSummary_Unsupported", func() error {
		_, err := client.GetRDFGraphSummary(ctx, &neptunedata.GetRDFGraphSummaryInput{})
		if err == nil {
			return fmt.Errorf("expected error for unsupported RDF graph summary")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "StartMLDataProcessingJob_Unsupported", func() error {
		_, err := client.StartMLDataProcessingJob(ctx, &neptunedata.StartMLDataProcessingJobInput{
			InputDataS3Location:     aws.String("s3://test/ml-input"),
			ProcessedDataS3Location: aws.String("s3://test/ml-output"),
		})
		if err == nil {
			return fmt.Errorf("expected error for unsupported ML data processing job")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ListGremlinQueries", func() error {
		_, err := client.ListGremlinQueries(ctx, &neptunedata.ListGremlinQueriesInput{})
		return err
	}))

	results = append(results, r.RunTest("neptunedata", "ListOpenCypherQueries", func() error {
		_, err := client.ListOpenCypherQueries(ctx, &neptunedata.ListOpenCypherQueriesInput{})
		return err
	}))

	return results
}

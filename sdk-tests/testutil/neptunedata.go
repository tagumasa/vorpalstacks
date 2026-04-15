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
		if err := AssertErrorContains(err, "UnsupportedOperationException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "GetRDFGraphSummary_Unsupported", func() error {
		_, err := client.GetRDFGraphSummary(ctx, &neptunedata.GetRDFGraphSummaryInput{})
		if err := AssertErrorContains(err, "UnsupportedOperationException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "StartMLDataProcessingJob_Unsupported", func() error {
		_, err := client.StartMLDataProcessingJob(ctx, &neptunedata.StartMLDataProcessingJobInput{
			InputDataS3Location:     aws.String("s3://test/ml-input"),
			ProcessedDataS3Location: aws.String("s3://test/ml-output"),
		})
		if err := AssertErrorContains(err, "UnsupportedOperationException"); err != nil {
			return err
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

	// --- Cypher advanced tests ---
	results = append(results, r.RunTest("neptunedata", "Cypher_Advanced_Reset", func() error {
		initResp, err := client.ExecuteFastReset(ctx, &neptunedata.ExecuteFastResetInput{
			Action: types.ActionInitializeReset,
		})
		if err != nil {
			return err
		}
		if initResp.Payload == nil || initResp.Payload.Token == nil {
			return fmt.Errorf("expected token")
		}
		_, err = client.ExecuteFastReset(ctx, &neptunedata.ExecuteFastResetInput{
			Action: types.ActionPerformReset,
			Token:  initResp.Payload.Token,
		})
		return err
	}))

	cypher := func(q string) error {
		_, err := client.ExecuteOpenCypherQuery(ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String(q),
		})
		return err
	}
	cypherResult := func(q string) (string, error) {
		resp, err := client.ExecuteOpenCypherQuery(ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String(q),
		})
		if err != nil {
			return "", err
		}
		return marshalDoc(resp.Results), nil
	}
	cypherContains := func(q, substr string) error {
		s, err := cypherResult(q)
		if err != nil {
			return err
		}
		if !strings.Contains(s, substr) {
			return fmt.Errorf("expected %q in result, got %s", substr, s)
		}
		return nil
	}

	cypher("CREATE (a:Person {name:'alice',age:25}), (b:Person {name:'bob',age:30}), (c:Person {name:'charlie',age:35}), (d:Person {name:'dave',age:40})")
	cypher("MATCH (a:Person {name:'alice'}), (b:Person {name:'bob'}) CREATE (a)-[:KNOWS {weight:0.5}]->(b)")
	cypher("MATCH (a:Person {name:'alice'}), (c:Person {name:'charlie'}) CREATE (a)-[:KNOWS {weight:1.0}]->(c)")
	cypher("MATCH (b:Person {name:'bob'}), (d:Person {name:'dave'}) CREATE (b)-[:KNOWS {weight:0.7}]->(d)")
	cypher("MATCH (c:Person {name:'charlie'}), (d:Person {name:'dave'}) CREATE (c)-[:KNOWS {weight:0.3}]->(d)")
	cypher("MATCH (a:Person {name:'alice'}), (b:Person {name:'bob'}) CREATE (a)-[:WORKS_WITH]->(b)")
	cypher("MATCH (b:Person {name:'bob'}), (c:Person {name:'charlie'}) CREATE (b)-[:WORKS_WITH]->(c)")

	results = append(results, r.RunTest("neptunedata", "Cypher_WhereClause", func() error {
		return cypherContains("MATCH (n:Person) WHERE n.age > 30 RETURN n.name", `"charlie"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_WhereAND", func() error {
		return cypherContains("MATCH (n:Person) WHERE n.age > 25 AND n.age < 40 RETURN n.name", `"bob"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_WhereOR", func() error {
		s, err := cypherResult("MATCH (n:Person) WHERE n.name = 'alice' OR n.name = 'dave' RETURN n.name ORDER BY n.name")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"alice"`) || !strings.Contains(s, `"dave"`) {
			return fmt.Errorf("expected alice and dave, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_WhereNOT", func() error {
		s, err := cypherResult("MATCH (n:Person) WHERE NOT n.name = 'alice' AND NOT n.name = 'dave' RETURN n.name ORDER BY n.name")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"bob"`) || !strings.Contains(s, `"charlie"`) {
			return fmt.Errorf("expected bob and charlie, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_WhereContains", func() error {
		return cypherContains("MATCH (n:Person) WHERE n.name CONTAINS 'li' RETURN n.name", `"alice"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_WhereStartsWith", func() error {
		return cypherContains("MATCH (n:Person) WHERE n.name STARTS WITH 'al' RETURN n.name", `"alice"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_WhereEndsWith", func() error {
		return cypherContains("MATCH (n:Person) WHERE n.name ENDS WITH 'ie' RETURN n.name", `"charlie"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_WhereIN", func() error {
		s, err := cypherResult("MATCH (n:Person) WHERE n.name IN ['alice', 'dave'] RETURN n.name ORDER BY n.name")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"alice"`) || !strings.Contains(s, `"dave"`) {
			return fmt.Errorf("expected alice and dave, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_WhereIsNull", func() error {
		cypher("CREATE (n:Person {name:'eve'})")
		return cypherContains("MATCH (n:Person) WHERE n.age IS NULL RETURN n.name", `"eve"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_WhereIsNotNull", func() error {
		s, err := cypherResult("MATCH (n:Person) WHERE n.age IS NOT NULL RETURN count(n) AS cnt")
		if err != nil {
			return err
		}
		if !strings.Contains(s, "4") {
			return fmt.Errorf("expected 4 persons with age, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_ReturnDistinct", func() error {
		s, err := cypherResult("MATCH (n:Person)-[:KNOWS]->(m) RETURN DISTINCT m.name ORDER BY m.name")
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
		s, err := cypherResult("MATCH (n:Person) WHERE n.age IS NOT NULL RETURN n.name, n.age ORDER BY n.age DESC")
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
		return cypherContains("MATCH (n:Person) WHERE n.age IS NOT NULL RETURN n.name ORDER BY n.age SKIP 2", `"charlie"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_Limit", func() error {
		return cypherContains("MATCH (n:Person) RETURN n.name LIMIT 1", `"alice"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_SkipAndLimit", func() error {
		s, err := cypherResult("MATCH (n:Person) WHERE n.age IS NOT NULL RETURN n.name ORDER BY n.age SKIP 1 LIMIT 1")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"bob"`) {
			return fmt.Errorf("expected bob at skip 1 limit 1, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_AggregationSum", func() error {
		return cypherContains("MATCH (n:Person) WHERE n.age IS NOT NULL RETURN sum(n.age) AS total", `130`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_AggregationAvg", func() error {
		return cypherContains("MATCH (n:Person) WHERE n.age IS NOT NULL RETURN avg(n.age) AS avg", `32.5`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_AggregationMin", func() error {
		return cypherContains("MATCH (n:Person) WHERE n.age IS NOT NULL RETURN min(n.age) AS mn", `25`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_AggregationMax", func() error {
		return cypherContains("MATCH (n:Person) WHERE n.age IS NOT NULL RETURN max(n.age) AS mx", `40`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_AggregationCollect", func() error {
		s, err := cypherResult("MATCH (n:Person) WHERE n.age IS NOT NULL RETURN collect(n.name) AS names")
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
		return cypherContains("MATCH (n:Person)-[:KNOWS]->(m) RETURN count(DISTINCT m.name) AS cnt", `3`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_SetProperty", func() error {
		if err := cypher("MATCH (n:Person {name:'alice'}) SET n.age = 26"); err != nil {
			return err
		}
		return cypherContains("MATCH (n:Person {name:'alice'}) RETURN n.age", `26`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_SetMergeProperties", func() error {
		if err := cypher("MATCH (n:Person {name:'alice'}) SET n += {city:'NYC',active:true}"); err != nil {
			return err
		}
		return cypherContains("MATCH (n:Person {name:'alice'}) RETURN n.city", `"NYC"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_SetLabel", func() error {
		if err := cypher("MATCH (n:Person {name:'alice'}) SET n:Employee"); err != nil {
			return err
		}
		return cypherContains("MATCH (n:Employee) RETURN n.name", `"alice"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_RemoveLabel", func() error {
		if err := cypher("MATCH (n:Person {name:'alice'}) REMOVE n:Employee"); err != nil {
			return err
		}
		s, err := cypherResult("MATCH (n:Employee) RETURN count(n) AS cnt")
		if err != nil {
			return err
		}
		if !strings.Contains(s, "0") {
			return fmt.Errorf("expected 0 employees after REMOVE, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_RemoveProperty", func() error {
		if err := cypher("MATCH (n:Person {name:'alice'}) REMOVE n.city"); err != nil {
			return err
		}
		s, err := cypherResult("MATCH (n:Person {name:'alice'}) RETURN n.city")
		if err != nil {
			return err
		}
		if strings.Contains(s, `"NYC"`) {
			return fmt.Errorf("expected city removed, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_DeleteNonDetach", func() error {
		cypher("CREATE (n:Temp {name:'temp_node'})")
		if err := cypher("MATCH (n:Temp) DELETE n"); err != nil {
			return err
		}
		return cypherContains("MATCH (n:Temp) RETURN count(n) AS cnt", `0`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_OptionalMatch", func() error {
		s, err := cypherResult("MATCH (a:Person {name:'alice'}) OPTIONAL MATCH (a)-[:WORKS_WITH]->(c) RETURN a.name, c.name ORDER BY c.name")
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
		s, err := cypherResult("MATCH (a:Person {name:'dave'}) OPTIONAL MATCH (a)-[:WORKS_WITH]->(c) RETURN a.name, c.name")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"dave"`) {
			return fmt.Errorf("expected dave, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_WithClause", func() error {
		return cypherContains("MATCH (n:Person) WHERE n.age IS NOT NULL WITH n.name AS name ORDER BY name LIMIT 2 RETURN name", `"alice"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_WithAggregation", func() error {
		return cypherContains("MATCH (n:Person)-[:KNOWS]->(m) WITH m.name AS friend, count(*) AS cnt RETURN friend, cnt ORDER BY cnt DESC", `"bob"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_Unwind", func() error {
		if err := cypher("UNWIND [10,20,30] AS x CREATE (n:Number {val:x})"); err != nil {
			return err
		}
		return cypherContains("MATCH (n:Number) RETURN sum(n.val) AS total", `60`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_MultiHop", func() error {
		return cypherContains("MATCH (a:Person {name:'alice'})-[:KNOWS]->(b)-[:KNOWS]->(c) RETURN c.name", `"dave"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_VarLengthPath", func() error {
		return cypherContains("MATCH (a:Person {name:'alice'})-[:KNOWS*1..2]->(b) RETURN DISTINCT b.name ORDER BY b.name", `"dave"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_CommaPatterns", func() error {
		s, err := cypherResult("MATCH (a:Person {name:'alice'}), (b:Person {name:'bob'}) RETURN a.name, b.name")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"alice"`) || !strings.Contains(s, `"bob"`) {
			return fmt.Errorf("expected both alice and bob, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_MultiMatch", func() error {
		return cypherContains("MATCH (a:Person {name:'alice'})-[:KNOWS]->(b) MATCH (b)-[:KNOWS]->(c) RETURN c.name", `"dave"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_MergeCreate", func() error {
		if err := cypher("MERGE (n:Animal {name:'rex'})"); err != nil {
			return err
		}
		return cypherContains("MATCH (n:Animal) RETURN n.name", `"rex"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_MergeIdempotent", func() error {
		cypher("MATCH (n:Animal {name:'rex'}) DELETE n")
		if err := cypher("MERGE (n:Animal {name:'rex'}) ON CREATE SET n.species = 'dog'"); err != nil {
			return err
		}
		return cypherContains("MATCH (n:Animal {name:'rex'}) RETURN n.species", `"dog"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_MergeOnMatchSet", func() error {
		if err := cypher("MERGE (n:Animal {name:'rex'}) ON MATCH SET n.visits = 1"); err != nil {
			return err
		}
		return cypherContains("MATCH (n:Animal {name:'rex'}) RETURN n.visits", `1`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_TypeFunction", func() error {
		return cypherContains("MATCH (a:Person {name:'alice'})-[r:KNOWS]->(b) RETURN type(r)", `"KNOWS"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_IdFunction", func() error {
		_, err := cypherResult("MATCH (n:Person {name:'alice'}) RETURN id(n) AS nid")
		return err
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_CoalesceFunction", func() error {
		return cypherContains("MATCH (n:Person {name:'eve'}) RETURN coalesce(n.age, 0) AS age", `0`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_CaseExpression", func() error {
		return cypherContains("MATCH (n:Person {name:'alice'}) RETURN CASE WHEN n.age > 30 THEN 'senior' ELSE 'junior' END AS category", `"junior"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_Profile", func() error {
		_, err := client.ExecuteOpenCypherExplainQuery(ctx, &neptunedata.ExecuteOpenCypherExplainQueryInput{
			OpenCypherQuery: aws.String("PROFILE MATCH (n) RETURN n LIMIT 1"),
			ExplainMode:     types.OpenCypherExplainModeDetails,
		})
		return err
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_MatchEdgeReturn", func() error {
		return cypherContains("MATCH (a:Person {name:'alice'})-[r:KNOWS]->(b) RETURN r.weight", `0.5`)
	}))

	results = append(results, r.RunTest("neptunedata", "Cypher_IncomingDirection", func() error {
		return cypherContains("MATCH (b:Person {name:'bob'})<-[:KNOWS]-(a) RETURN a.name", `"alice"`)
	}))

	// --- Gremlin advanced tests ---
	results = append(results, r.RunTest("neptunedata", "Gremlin_Advanced_Reset", func() error {
		initResp, err := client.ExecuteFastReset(ctx, &neptunedata.ExecuteFastResetInput{
			Action: types.ActionInitializeReset,
		})
		if err != nil {
			return err
		}
		if initResp.Payload == nil || initResp.Payload.Token == nil {
			return fmt.Errorf("expected token")
		}
		_, err = client.ExecuteFastReset(ctx, &neptunedata.ExecuteFastResetInput{
			Action: types.ActionPerformReset,
			Token:  initResp.Payload.Token,
		})
		return err
	}))

	gremlin := func(q string) error {
		_, err := client.ExecuteGremlinQuery(ctx, &neptunedata.ExecuteGremlinQueryInput{
			GremlinQuery: aws.String(q),
		})
		return err
	}
	gremlinResult := func(q string) (string, error) {
		resp, err := client.ExecuteGremlinQuery(ctx, &neptunedata.ExecuteGremlinQueryInput{
			GremlinQuery: aws.String(q),
		})
		if err != nil {
			return "", err
		}
		return marshalDoc(resp.Result), nil
	}
	gremlinContains := func(q, substr string) error {
		s, err := gremlinResult(q)
		if err != nil {
			return err
		}
		if !strings.Contains(s, substr) {
			return fmt.Errorf("expected %q in result, got %s", substr, s)
		}
		return nil
	}

	gremlin("g.addV('person').property('name','alice').property('age',25)")
	gremlin("g.addV('person').property('name','bob').property('age',30)")
	gremlin("g.addV('person').property('name','charlie').property('age',35)")
	gremlin("g.addV('person').property('name','dave').property('age',40)")
	gremlin("g.V().has('name','alice').addE('knows').to(g.V().has('name','bob')).property('weight',0.5)")
	gremlin("g.V().has('name','alice').addE('knows').to(g.V().has('name','charlie')).property('weight',1.0)")
	gremlin("g.V().has('name','bob').addE('knows').to(g.V().has('name','dave')).property('weight',0.7)")
	gremlin("g.V().has('name','charlie').addE('knows').to(g.V().has('name','dave')).property('weight',0.3)")

	results = append(results, r.RunTest("neptunedata", "Gremlin_In", func() error {
		return gremlinContains("g.V().has('name','bob').in('knows').values('name')", `"alice"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Both", func() error {
		s, err := gremlinResult("g.V().has('name','bob').both('knows').values('name').order()")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"alice"`) || !strings.Contains(s, `"dave"`) {
			return fmt.Errorf("expected alice and dave, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_OutE", func() error {
		s, err := gremlinResult("g.V().has('name','alice').outE('knows').count()")
		if err != nil {
			return err
		}
		if !strings.Contains(s, "2") {
			return fmt.Errorf("expected 2 outgoing edges from alice, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_InE", func() error {
		s, err := gremlinResult("g.V().has('name','bob').inE('knows').count()")
		if err != nil {
			return err
		}
		if !strings.Contains(s, "1") {
			return fmt.Errorf("expected 1 incoming edge to bob, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_BothE", func() error {
		s, err := gremlinResult("g.V().has('name','bob').bothE().count()")
		if err != nil {
			return err
		}
		if !strings.Contains(s, "2") {
			return fmt.Errorf("expected 2 edges for bob (1 in + 1 out), got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_OutV", func() error {
		return gremlinContains("g.E().hasLabel('knows').outV().has('name','alice').count()", `2`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_InV", func() error {
		return gremlinContains("g.E().hasLabel('knows').inV().has('name','bob').count()", `1`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_OtherV", func() error {
		return gremlinContains("g.V().has('name','alice').outE('knows').otherV().has('name','bob').count()", `1`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_HasPredicate", func() error {
		return gremlinContains("g.V().has('age',P.gt(30)).values('name')", `"charlie"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_HasGte", func() error {
		s, err := gremlinResult("g.V().has('age',P.gte(30)).count()")
		if err != nil {
			return err
		}
		if !strings.Contains(s, "3") {
			return fmt.Errorf("expected 3 vertices with age>=30, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_HasLt", func() error {
		return gremlinContains("g.V().has('age',P.lt(30)).values('name')", `"alice"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_HasNeq", func() error {
		s, err := gremlinResult("g.V().has('name',P.neq('alice')).count()")
		if err != nil {
			return err
		}
		if !strings.Contains(s, "3") {
			return fmt.Errorf("expected 3 non-alice vertices, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_HasBetween", func() error {
		s, err := gremlinResult("g.V().has('age',P.between(25,35)).count()")
		if err != nil {
			return err
		}
		if !strings.Contains(s, "3") {
			return fmt.Errorf("expected 3 vertices with age between 25-35, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_HasWithin", func() error {
		return gremlinContains("g.V().has('name',P.within('alice','dave')).values('name').order()", `"alice"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_HasNot", func() error {
		cypher("CREATE (n:Person {name:'frank'})")
		return gremlinContains("g.V().hasNot('age').values('name')", `"frank"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_StartingWith", func() error {
		return gremlinContains("g.V().has('name',P.startingWith('al')).values('name')", `"alice"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Containing", func() error {
		return gremlinContains("g.V().has('name',P.containing('li')).values('name')", `"alice"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_EndingWith", func() error {
		return gremlinContains("g.V().has('name',P.endingWith('ie')).values('name')", `"charlie"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Id", func() error {
		_, err := gremlinResult("g.V().has('name','alice').id()")
		return err
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Label", func() error {
		return gremlinContains("g.V().has('name','alice').label()", `"person"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Keys", func() error {
		s, err := gremlinResult("g.V().has('name','alice').keys()")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"name"`) || !strings.Contains(s, `"age"`) {
			return fmt.Errorf("expected name and age keys, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Properties", func() error {
		return gremlinContains("g.V().has('name','alice').properties('name')", `"name"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_PropertyMap", func() error {
		return gremlinContains("g.V().has('name','alice').propertyMap('name','age')", `"alice"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_PropertyMutation", func() error {
		gremlin("g.V().has('name','alice').property('city','NYC')")
		return gremlinContains("g.V().has('name','alice').values('city')", `"NYC"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_AsSelect", func() error {
		return gremlinContains("g.V().has('name','alice').as('a').out('knows').as('b').select('a','b').by('name')", `"alice"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_WhereTraversal", func() error {
		return gremlinContains("g.V().where(out().count().is(gt(0))).values('name').dedup()", `"alice"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Filter", func() error {
		return gremlinContains("g.V().filter(out('knows').count().is(gt(1))).values('name')", `"alice"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Limit", func() error {
		return gremlinContains("g.V().hasLabel('person').limit(2).count()", `2`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Range", func() error {
		return gremlinContains("g.V().hasLabel('person').range(1,3).count()", `2`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Skip", func() error {
		s, err := gremlinResult("g.V().hasLabel('person').order().by('name').skip(2).values('name')")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"charlie"`) {
			return fmt.Errorf("expected charlie at skip 2, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Dedup", func() error {
		gremlin("g.addV('dup').property('name','dup1')")
		gremlin("g.addV('dup').property('name','dup1')")
		s, err := gremlinResult("g.V().hasLabel('dup').dedup().by('name').count()")
		if err != nil {
			return err
		}
		if !strings.Contains(s, "1") {
			return fmt.Errorf("expected 1 after dedup by name, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_OrderByAsc", func() error {
		s, err := gremlinResult("g.V().hasLabel('person').order().by('name',asc).values('name').limit(1)")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"alice"`) {
			return fmt.Errorf("expected alice first in ASC order, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_OrderByDesc", func() error {
		s, err := gremlinResult("g.V().hasLabel('person').order().by('name',desc).values('name').limit(1)")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"dave"`) {
			return fmt.Errorf("expected dave first in DESC order, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_GroupCount", func() error {
		s, err := gremlinResult("g.V().hasLabel('person').groupCount().by('name')")
		if err != nil {
			return err
		}
		for _, name := range []string{`"alice"`, `"bob"`, `"charlie"`, `"dave"`} {
			if !strings.Contains(s, name) {
				return fmt.Errorf("expected %s in groupCount, got %s", name, s)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_GroupBy", func() error {
		s, err := gremlinResult("g.V().hasLabel('person').group().by('name').by('age')")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"alice"`) || !strings.Contains(s, `"bob"`) {
			return fmt.Errorf("expected grouped names, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Path", func() error {
		s, err := gremlinResult("g.V().has('name','alice').out('knows').limit(1).path()")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"alice"`) {
			return fmt.Errorf("expected alice in path, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Union", func() error {
		s, err := gremlinResult("g.V().has('name','bob').union(in('knows'),out('knows')).values('name').dedup().order()")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"alice"`) || !strings.Contains(s, `"dave"`) {
			return fmt.Errorf("expected alice and dave from union, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Coalesce", func() error {
		return gremlinContains("g.V().has('name','bob').coalesce(out('nonexistent'),in('knows')).values('name')", `"alice"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Choose", func() error {
		return gremlinContains("g.V().has('name','alice').choose(has('age'),values('age'),values('name'))", `25`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Optional", func() error {
		s, err := gremlinResult("g.V().has('name','dave').optional(out('knows')).values('name')")
		if err != nil {
			return err
		}
		if strings.Contains(s, `"charlie"`) {
			return fmt.Errorf("optional should return original when no outgoing knows, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_RepeatTimes", func() error {
		return gremlinContains("g.V().has('name','alice').repeat(out('knows')).times(2).values('name')", `"dave"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_RepeatEmit", func() error {
		s, err := gremlinResult("g.V().has('name','alice').repeat(out('knows')).emit().values('name').dedup().order()")
		if err != nil {
			return err
		}
		for _, name := range []string{`"bob"`, `"charlie"`, `"dave"`} {
			if !strings.Contains(s, name) {
				return fmt.Errorf("expected %s in repeat-emit, got %s", name, s)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Fold", func() error {
		s, err := gremlinResult("g.V().hasLabel('person').values('name').fold()")
		if err != nil {
			return err
		}
		for _, name := range []string{`"alice"`, `"bob"`, `"charlie"`, `"dave"`} {
			if !strings.Contains(s, name) {
				return fmt.Errorf("expected %s in fold, got %s", name, s)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Unfold", func() error {
		return gremlinContains("g.V().hasLabel('person').values('name').fold().unfold().count()", `4`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Constant", func() error {
		return gremlinContains("g.V().has('name','alice').constant('hello').limit(1)", `"hello"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Is", func() error {
		return gremlinContains("g.V().hasLabel('person').values('age').is(P.gt(30))", `35`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Not", func() error {
		return gremlinContains("g.V().hasLabel('person').not(out('knows')).values('name')", `"dave"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_And", func() error {
		return gremlinContains("g.V().and(out('knows'),in('knows')).values('name')", `"bob"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Or", func() error {
		s, err := gremlinResult("g.V().or(has('name','alice'),has('name','dave')).values('name').order()")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"alice"`) || !strings.Contains(s, `"dave"`) {
			return fmt.Errorf("expected alice or dave, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Mean", func() error {
		return gremlinContains("g.V().hasLabel('person').values('age').mean()", `32.5`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Sum", func() error {
		return gremlinContains("g.V().hasLabel('person').values('age').sum()", `130`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Min", func() error {
		return gremlinContains("g.V().hasLabel('person').values('age').min()", `25`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Max", func() error {
		return gremlinContains("g.V().hasLabel('person').values('age').max()", `40`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Tail", func() error {
		s, err := gremlinResult("g.V().hasLabel('person').order().by('name').tail(1).values('name')")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"dave"`) {
			return fmt.Errorf("expected dave as tail, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Inject", func() error {
		return gremlinContains("g.inject('a','b').count()", `2`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_MergeV", func() error {
		gremlin("g.mergeV([~label:'item',name:'widget',price:10])")
		return gremlinContains("g.V().has('name','widget').values('price')", `10`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_MergeVIdempotent", func() error {
		gremlin("g.mergeV([~label:'item',name:'widget',price:10])")
		s, err := gremlinResult("g.V().has('name','widget').count()")
		if err != nil {
			return err
		}
		if !strings.Contains(s, "1") {
			return fmt.Errorf("expected 1 widget after idempotent mergeV, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_ElementMap", func() error {
		s, err := gremlinResult("g.V().has('name','alice').elementMap()")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"alice"`) || !strings.Contains(s, `"person"`) {
			return fmt.Errorf("expected alice and person in elementMap, got %s", s)
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_Project", func() error {
		return gremlinContains("g.V().as('a').out('knows').as('b').project('from','to').by('name').limit(1)", `"from"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_CountLocal", func() error {
		return gremlinContains("g.V().hasLabel('person').values('name').fold().count(local)", `4`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_MultiHop", func() error {
		return gremlinContains("g.V().has('name','alice').out().out().values('name')", `"dave"`)
	}))

	results = append(results, r.RunTest("neptunedata", "Gremlin_SimplePath", func() error {
		s, err := gremlinResult("g.V().has('name','alice').repeat(__.out('knows')).times(2).simplePath().values('name').dedup()")
		if err != nil {
			return err
		}
		if !strings.Contains(s, `"dave"`) {
			return fmt.Errorf("expected dave in simplePath result, got %s", s)
		}
		return nil
	}))

	// --- Server API tests ---
	results = append(results, r.RunTest("neptunedata", "CancelGremlinQuery", func() error {
		_, err := client.CancelGremlinQuery(ctx, &neptunedata.CancelGremlinQueryInput{
			QueryId: aws.String("nonexistent-query-id"),
		})
		if err := AssertErrorContains(err, "BadRequestException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "CancelOpenCypherQuery", func() error {
		_, err := client.ExecuteOpenCypherQuery(ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("MATCH (n) RETURN count(n)"),
		})
		if err != nil {
			return err
		}
		_, err = client.ListOpenCypherQueries(ctx, &neptunedata.ListOpenCypherQueriesInput{})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ManagePropertygraphStatistics_Disable", func() error {
		_, err := client.ManagePropertygraphStatistics(ctx, &neptunedata.ManagePropertygraphStatisticsInput{
			Mode: types.StatisticsAutoGenerationModeDisableAutocompute,
		})
		return err
	}))

	results = append(results, r.RunTest("neptunedata", "ManagePropertygraphStatistics_Enable", func() error {
		_, err := client.ManagePropertygraphStatistics(ctx, &neptunedata.ManagePropertygraphStatisticsInput{
			Mode: types.StatisticsAutoGenerationModeEnableAutocompute,
		})
		return err
	}))

	results = append(results, r.RunTest("neptunedata", "ManagePropertygraphStatistics_Refresh", func() error {
		resp, err := client.ManagePropertygraphStatistics(ctx, &neptunedata.ManagePropertygraphStatisticsInput{
			Mode: types.StatisticsAutoGenerationModeRefresh,
		})
		if err != nil {
			return err
		}
		if resp.Status == nil {
			return fmt.Errorf("expected non-nil status from refresh")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "DeletePropertygraphStatistics", func() error {
		resp, err := client.DeletePropertygraphStatistics(ctx, &neptunedata.DeletePropertygraphStatisticsInput{})
		if err != nil {
			return err
		}
		if resp.Status == nil {
			return fmt.Errorf("expected non-nil status from delete statistics")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "GetPropertygraphStream", func() error {
		resp, err := client.GetPropertygraphStream(ctx, &neptunedata.GetPropertygraphStreamInput{})
		if err != nil {
			return err
		}
		if resp.Format == nil {
			return fmt.Errorf("expected non-nil format from stream")
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "GetPropertygraphSummary_Detailed", func() error {
		_, err := client.GetPropertygraphSummary(ctx, &neptunedata.GetPropertygraphSummaryInput{
			Mode: types.GraphSummaryTypeDetailed,
		})
		return err
	}))

	results = append(results, r.RunTest("neptunedata", "ErrorCase_InvalidCypherSyntax", func() error {
		_, err := client.ExecuteOpenCypherQuery(ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String("INVALID CYPHER QUERY"),
		})
		if err := AssertErrorContains(err, "MalformedQueryException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ErrorCase_InvalidGremlinSyntax", func() error {
		_, err := client.ExecuteGremlinQuery(ctx, &neptunedata.ExecuteGremlinQueryInput{
			GremlinQuery: aws.String("g.INVALID_STEP()"),
		})
		if err := AssertErrorContains(err, "MalformedQueryException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ErrorCase_FastResetInvalidToken", func() error {
		_, err := client.ExecuteFastReset(ctx, &neptunedata.ExecuteFastResetInput{
			Action: types.ActionPerformReset,
			Token:  aws.String("invalid-token-12345"),
		})
		if err := AssertErrorContains(err, "PreconditionsFailedException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ErrorCase_NonExistentLoaderJob", func() error {
		_, err := client.GetLoaderJobStatus(ctx, &neptunedata.GetLoaderJobStatusInput{
			LoadId: aws.String("nonexistent-load-id"),
		})
		if err := AssertErrorContains(err, "BulkLoadIdNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ErrorCase_CancelNonExistentLoaderJob", func() error {
		_, err := client.CancelLoaderJob(ctx, &neptunedata.CancelLoaderJobInput{
			LoadId: aws.String("nonexistent-load-id"),
		})
		if err := AssertErrorContains(err, "BulkLoadIdNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ErrorCase_EmptyCypherQuery", func() error {
		_, err := client.ExecuteOpenCypherQuery(ctx, &neptunedata.ExecuteOpenCypherQueryInput{
			OpenCypherQuery: aws.String(""),
		})
		if err := AssertErrorContains(err, "MissingParameterException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunedata", "ErrorCase_EmptyGremlinQuery", func() error {
		_, err := client.ExecuteGremlinQuery(ctx, &neptunedata.ExecuteGremlinQueryInput{
			GremlinQuery: aws.String(""),
		})
		if err := AssertErrorContains(err, "MissingParameterException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}

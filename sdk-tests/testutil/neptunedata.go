package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
	"github.com/aws/aws-sdk-go-v2/service/neptunedata"
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

type neptunedataContext struct {
	client    *neptunedata.Client
	ctx       context.Context
	clusterID string
	region    string
	accountID string
}

func (r *TestRunner) runNeptunedataClusterTests(tc *neptunedataContext) []TestResult {
	var results []TestResult
	results = append(results, r.runNeptunedataCypherBasicTests(tc)...)
	results = append(results, r.runNeptunedataCypherAdvancedTests(tc)...)
	results = append(results, r.runNeptunedataGremlinBasicTests(tc)...)
	results = append(results, r.runNeptunedataGremlinAdvancedTests(tc)...)
	results = append(results, r.runNeptunedataStatisticsTests(tc)...)
	results = append(results, r.runNeptunedataLoaderTests(tc)...)
	results = append(results, r.runNeptunedataEdgeTests(tc)...)
	return results
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

	tc := &neptunedataContext{
		client:    neptunedata.NewFromConfig(cfg),
		ctx:       context.Background(),
		region:    r.region,
		accountID: r.accountID,
	}

	results = append(results, r.runNeptunedataEngineTests(tc)...)
	results = append(results, r.runNeptunedataResetTests(tc)...)

	clusterID := fmt.Sprintf("sdk-test-%d", time.Now().UnixNano())
	clusterPort, cleanup := r.createNeptuneClusterForDataTests(clusterID)
	defer cleanup()
	if clusterPort > 0 {
		clusterCfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: fmt.Sprintf("http://127.0.0.1:%d", clusterPort),
			Region:   r.region,
		})
		if err != nil {
			log.Printf("neptunedata: cluster config load failed: %v", err)
			results = append(results, SetupFailResult("neptunedata", "cluster config load failed: %v", err))
		} else {
			tcCluster := &neptunedataContext{
				client:    neptunedata.NewFromConfig(clusterCfg),
				ctx:       context.Background(),
				clusterID: clusterID,
			}
			results = append(results, r.runNeptunedataClusterTests(tcCluster)...)
		}
	} else {
		results = append(results, SetupFailResult("neptunedata", "cluster creation returned port 0"))
	}

	results = append(results, r.runNeptunedataServerAPITests(tc)...)
	results = append(results, r.runNeptunedataUnsupportedTests(tc)...)

	return results
}

func (r *TestRunner) createNeptuneClusterForDataTests(clusterID string) (int, func()) {
	noop := func() {}

	neptuneCfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return 0, noop
	}

	neptuneClient := neptune.NewFromConfig(neptuneCfg)

	_, err = neptuneClient.CreateDBCluster(context.Background(), &neptune.CreateDBClusterInput{
		DBClusterIdentifier: aws.String(clusterID),
		Engine:              aws.String("neptune"),
	})
	if err != nil {
		return 0, noop
	}

	cleanup := func() {
		_, delErr := neptuneClient.DeleteDBCluster(context.Background(), &neptune.DeleteDBClusterInput{
			DBClusterIdentifier: aws.String(clusterID),
			SkipFinalSnapshot:   aws.Bool(true),
		})
		if delErr != nil {
			log.Printf("neptunedata: failed to delete cluster %s: %v", clusterID, delErr)
		}
	}

	resp, err := neptuneClient.DescribeDBClusters(context.Background(), &neptune.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(clusterID),
	})
	if err != nil || len(resp.DBClusters) == 0 {
		cleanup()
		return 0, noop
	}

	endpoint := resp.DBClusters[0].Endpoint
	if endpoint == nil || *endpoint == "" {
		cleanup()
		return 0, noop
	}

	var port int
	u := fmt.Sprintf("http://%s", *endpoint)
	if parsed, parseErr := url.Parse(u); parseErr == nil && parsed.Port() != "" {
		if _, scanErr := fmt.Sscanf(parsed.Port(), "%d", &port); scanErr != nil {
			log.Printf("neptunedata: failed to parse port from %q: %v", *endpoint, scanErr)
		}
	}
	if port == 0 {
		log.Printf("neptunedata: no port in endpoint %q, falling back to default", *endpoint)
	}

	return port, cleanup
}

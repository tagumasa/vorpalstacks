package testutil

import (
	"context"
	"encoding/json"
	"fmt"

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
	client *neptunedata.Client
	ctx    context.Context
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
		client: neptunedata.NewFromConfig(cfg),
		ctx:    context.Background(),
	}

	results = append(results, r.runNeptunedataEngineTests(tc)...)
	results = append(results, r.runNeptunedataResetTests(tc)...)
	results = append(results, r.runNeptunedataCypherBasicTests(tc)...)
	results = append(results, r.runNeptunedataCypherAdvancedTests(tc)...)
	results = append(results, r.runNeptunedataGremlinBasicTests(tc)...)
	results = append(results, r.runNeptunedataGremlinAdvancedTests(tc)...)
	results = append(results, r.runNeptunedataStatisticsTests(tc)...)
	results = append(results, r.runNeptunedataLoaderTests(tc)...)
	results = append(results, r.runNeptunedataServerAPITests(tc)...)
	results = append(results, r.runNeptunedataUnsupportedTests(tc)...)
	results = append(results, r.runNeptunedataEdgeTests(tc)...)

	return results
}

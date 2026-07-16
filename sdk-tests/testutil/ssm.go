package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"vorpalstacks-sdk-tests/config"
)

type ssmTestContext struct {
	client *ssm.Client
	ctx    context.Context
	ts     string
	r      *TestRunner
}

func (tc *ssmTestContext) uniqueName(prefix string) string {
	return fmt.Sprintf("%s/param-%s-%d", prefix, tc.ts, time.Now().UnixNano())
}

func (tc *ssmTestContext) putParam(name, value string, paramType types.ParameterType, opts ...func(*ssm.PutParameterInput)) (*ssm.PutParameterOutput, error) {
	input := &ssm.PutParameterInput{
		Name:  aws.String(name),
		Value: aws.String(value),
		Type:  paramType,
	}
	for _, opt := range opts {
		opt(input)
	}
	return tc.client.PutParameter(tc.ctx, input)
}

func (tc *ssmTestContext) deleteParam(name string) {
	tc.client.DeleteParameter(tc.ctx, &ssm.DeleteParameterInput{Name: aws.String(name)})
}

func (r *TestRunner) RunSSMTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "ssm",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := ssm.NewFromConfig(cfg)
	ctx := context.Background()
	ts := fmt.Sprintf("%d", time.Now().UnixNano())

	tc := &ssmTestContext{
		client: client,
		ctx:    ctx,
		ts:     ts,
		r:      r,
	}

	results = append(results, r.runSSMParameterCRUD(tc)...)
	results = append(results, r.runSSMParameterList(tc)...)
	results = append(results, r.runSSMTag(tc)...)
	results = append(results, r.runSSMVersion(tc)...)
	results = append(results, r.runSSMBatch(tc)...)
	results = append(results, r.runSSMEdge(tc)...)

	return results
}

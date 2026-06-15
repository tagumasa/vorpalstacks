package testutil

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/iot"
	"vorpalstacks-sdk-tests/config"
)

type iotTestContext struct {
	client *iot.Client
	ctx    context.Context
}

func (r *TestRunner) newIoTTestContext() (*iotTestContext, error) {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &iotTestContext{
		client: iot.NewFromConfig(cfg),
		ctx:    context.Background(),
	}, nil
}

func (r *TestRunner) RunIoTTests() []TestResult {
	tc, err := r.newIoTTestContext()
	if err != nil {
		return []TestResult{{
			Service:  "iot",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    err.Error(),
		}}
	}

	var results []TestResult
	results = append(results, r.runIoTThingTests(tc)...)
	results = append(results, r.runIoTCertTests(tc)...)
	results = append(results, r.runIoTPolicyTests(tc)...)
	results = append(results, r.runIoTParamFixTests(tc)...)
	return results
}

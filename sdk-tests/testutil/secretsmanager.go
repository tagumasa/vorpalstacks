package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"vorpalstacks-sdk-tests/config"
)

type secretsManagerTestContext struct {
	client *secretsmanager.Client
	ctx    context.Context
	region string
}

func (r *TestRunner) initSecretsManager() (*secretsManagerTestContext, error) {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &secretsManagerTestContext{
		client: secretsmanager.NewFromConfig(cfg),
		ctx:    context.Background(),
		region: r.region,
	}, nil
}

func (tc *secretsManagerTestContext) uniqueName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

func (tc *secretsManagerTestContext) createAndDeleteSecret(name, value string) (string, error) {
	resp, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
		Name:         &name,
		SecretString: &value,
	})
	if err != nil {
		return "", err
	}
	return *resp.ARN, nil
}

func (tc *secretsManagerTestContext) forceDeleteSecret(name string) {
	tc.client.DeleteSecret(tc.ctx, &secretsmanager.DeleteSecretInput{
		SecretId:                   &name,
		ForceDeleteWithoutRecovery: boolPtr(true),
	})
}

func boolPtr(b bool) *bool {
	return &b
}

func int64Ptr(v int64) *int64 {
	return &v
}

func (r *TestRunner) RunSecretsManagerTests() []TestResult {
	tc, err := r.initSecretsManager()
	if err != nil {
		return []TestResult{{
			Service:  "secretsmanager",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    err.Error(),
		}}
	}

	var results []TestResult
	results = append(results, r.runSecretsManagerSecretTests(tc)...)
	results = append(results, r.runSecretsManagerValueTests(tc)...)
	results = append(results, r.runSecretsManagerPolicyTests(tc)...)
	results = append(results, r.runSecretsManagerTagTests(tc)...)
	results = append(results, r.runSecretsManagerRotationTests(tc)...)
	results = append(results, r.runSecretsManagerPasswordTests(tc)...)
	results = append(results, r.runSecretsManagerBatchTests(tc)...)
	results = append(results, r.runSecretsManagerReplicationTests(tc)...)
	results = append(results, r.runSecretsManagerEdgeTests(tc)...)
	return results
}

package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"vorpalstacks-sdk-tests/config"
)

type secretsManagerTestContext struct {
	client    *secretsmanager.Client
	ctx       context.Context
	region    string
	accountID string
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
		client:    secretsmanager.NewFromConfig(cfg),
		ctx:       context.Background(),
		region:    r.region,
		accountID: r.accountID,
	}, nil
}

func (tc *secretsManagerTestContext) uniqueName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

// createSecret creates a secret carrying only Name and SecretString. Tests
// that pass additional input members (Description, SecretBinary, Tags) or
// exercise CreateSecret itself keep their literal CreateSecretInput.
func (tc *secretsManagerTestContext) createSecret(name, value string) (*secretsmanager.CreateSecretOutput, error) {
	resp, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
		Name:         &name,
		SecretString: &value,
	})
	if err != nil {
		return nil, fmt.Errorf("create %s: %w", name, err)
	}
	return resp, nil
}

// describeSecret reads a secret by plain SecretId; decorated Describe calls
// and the DescribeSecret operation test keep their literals.
func (tc *secretsManagerTestContext) describeSecret(name string) (*secretsmanager.DescribeSecretOutput, error) {
	resp, err := tc.client.DescribeSecret(tc.ctx, &secretsmanager.DescribeSecretInput{
		SecretId: &name,
	})
	if err != nil {
		return nil, fmt.Errorf("describe %s: %w", name, err)
	}
	return resp, nil
}

// getSecretValue reads the current secret value by plain SecretId; the
// GetSecretValue operation test keeps its literal.
func (tc *secretsManagerTestContext) getSecretValue(name string) (*secretsmanager.GetSecretValueOutput, error) {
	resp, err := tc.client.GetSecretValue(tc.ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &name,
	})
	if err != nil {
		return nil, fmt.Errorf("get %s: %w", name, err)
	}
	return resp, nil
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

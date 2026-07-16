package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func (r *TestRunner) runSecretsManagerPolicyTests(tc *secretsManagerTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("secretsmanager", "PutResourcePolicy_Basic", func() error {
		name := tc.uniqueName("PolicyTest")
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String("policy-test"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(name)

		policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"secretsmanager:GetSecretValue","Resource":"*"}]}`
		putResp, err := tc.client.PutResourcePolicy(tc.ctx, &secretsmanager.PutResourcePolicyInput{
			SecretId:       aws.String(name),
			ResourcePolicy: aws.String(policy),
		})
		if err != nil {
			return fmt.Errorf("put policy: %v", err)
		}
		if putResp.ARN == nil {
			return fmt.Errorf("ARN is nil")
		}
		if putResp.Name == nil || *putResp.Name != name {
			return fmt.Errorf("name mismatch")
		}

		getResp, err := tc.client.GetResourcePolicy(tc.ctx, &secretsmanager.GetResourcePolicyInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("get policy: %v", err)
		}
		if getResp.ResourcePolicy == nil || *getResp.ResourcePolicy != policy {
			return fmt.Errorf("policy mismatch: got %q", aws.ToString(getResp.ResourcePolicy))
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "DeleteResourcePolicy_Basic", func() error {
		name := tc.uniqueName("DelPolicy")
		_, err := tc.client.CreateSecret(tc.ctx, &secretsmanager.CreateSecretInput{
			Name:         aws.String(name),
			SecretString: aws.String("del-policy"),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.forceDeleteSecret(name)

		policy := `{"Version":"2012-10-17","Statement":[]}`
		_, err = tc.client.PutResourcePolicy(tc.ctx, &secretsmanager.PutResourcePolicyInput{
			SecretId:       aws.String(name),
			ResourcePolicy: aws.String(policy),
		})
		if err != nil {
			return fmt.Errorf("put policy: %v", err)
		}

		delResp, err := tc.client.DeleteResourcePolicy(tc.ctx, &secretsmanager.DeleteResourcePolicyInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("delete policy: %v", err)
		}
		if delResp.ARN == nil {
			return fmt.Errorf("ARN is nil")
		}

		getResp, err := tc.client.GetResourcePolicy(tc.ctx, &secretsmanager.GetResourcePolicyInput{
			SecretId: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("get policy: %v", err)
		}
		if getResp.ResourcePolicy != nil {
			return fmt.Errorf("policy should be nil after deletion")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "ValidateResourcePolicy_Valid", func() error {
		policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"*","Resource":"*"}]}`
		resp, err := tc.client.ValidateResourcePolicy(tc.ctx, &secretsmanager.ValidateResourcePolicyInput{
			ResourcePolicy: aws.String(policy),
		})
		if err != nil {
			return err
		}
		if !resp.PolicyValidationPassed {
			return fmt.Errorf("expected policy validation to pass")
		}
		if len(resp.ValidationErrors) > 0 {
			return fmt.Errorf("expected no validation errors, got %d", len(resp.ValidationErrors))
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "ValidateResourcePolicy_Invalid", func() error {
		resp, err := tc.client.ValidateResourcePolicy(tc.ctx, &secretsmanager.ValidateResourcePolicyInput{
			ResourcePolicy: aws.String("not valid json {"),
		})
		if err != nil {
			return err
		}
		if resp.PolicyValidationPassed {
			return fmt.Errorf("expected policy validation to fail for invalid JSON")
		}
		if len(resp.ValidationErrors) == 0 {
			return fmt.Errorf("expected validation errors for invalid JSON")
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "GetResourcePolicy_NonExistent", func() error {
		_, err := tc.client.GetResourcePolicy(tc.ctx, &secretsmanager.GetResourcePolicyInput{
			SecretId: aws.String("nonexistent-policy-secret"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "PutResourcePolicy_NonExistent", func() error {
		_, err := tc.client.PutResourcePolicy(tc.ctx, &secretsmanager.PutResourcePolicyInput{
			SecretId:       aws.String("nonexistent-policy-secret"),
			ResourcePolicy: aws.String(`{"Version":"2012-10-17"}`),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}

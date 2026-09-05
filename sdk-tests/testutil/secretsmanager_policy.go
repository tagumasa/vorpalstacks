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
		_, err := tc.createSecret(name, "policy-test")
		if err != nil {
			return err
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
		_, err := tc.createSecret(name, "del-policy")
		if err != nil {
			return err
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

	results = append(results, r.RunTest("secretsmanager", "ValidateResourcePolicy_Verdict", func() error {
		// The validation verdict must track the policy's validity: a valid
		// document passes with no errors, invalid JSON fails with errors.
		rows := []struct {
			name       string
			policy     string
			wantPassed bool
			wantErrors bool
		}{
			{"Valid", `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"*","Resource":"*"}]}`, true, false},
			{"Invalid", "not valid json {", false, true},
		}
		for _, row := range rows {
			resp, err := tc.client.ValidateResourcePolicy(tc.ctx, &secretsmanager.ValidateResourcePolicyInput{
				ResourcePolicy: aws.String(row.policy),
			})
			if err != nil {
				return fmt.Errorf("%s: %w", row.name, err)
			}
			if resp.PolicyValidationPassed != row.wantPassed {
				return fmt.Errorf("%s: PolicyValidationPassed = %v, want %v", row.name, resp.PolicyValidationPassed, row.wantPassed)
			}
			if (len(resp.ValidationErrors) > 0) != row.wantErrors {
				return fmt.Errorf("%s: got %d validation errors, wantErrors %v", row.name, len(resp.ValidationErrors), row.wantErrors)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("secretsmanager", "ValidateResourcePolicy_TooLongRejected", func() error {
		// NonEmptyResourcePolicyType caps the document at 20480 bytes; the
		// same constraint binds ValidateResourcePolicy as PutResourcePolicy.
		pad := make([]byte, 21000)
		for i := range pad {
			pad[i] = 'a'
		}
		policy := `{"Version":"2012-10-17","Statement":[{"Sid":"` + string(pad) + `","Effect":"Allow","Principal":"*","Action":"secretsmanager:GetSecretValue","Resource":"*"}]}`
		_, err := tc.client.ValidateResourcePolicy(tc.ctx, &secretsmanager.ValidateResourcePolicyInput{
			ResourcePolicy: aws.String(policy),
		})
		if err == nil {
			return fmt.Errorf("oversized ResourcePolicy should be rejected")
		}
		return expectAWSErrorCode(err, "InvalidParameterException")
	}))

	return results
}

package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
)

const testIoTPolicyDocument = `{
  "Version": "2012-10-17",
  "Statement": [
    {"Effect": "Allow", "Action": "iot:Connect", "Resource": "*"}
  ]
}`

// runIoTPolicyTests covers the Policy lifecycle: Create/Get/List/Delete,
// non-default policy versions (Create/SetDefault/List/Get/Delete), and
// Attach/Detach to a certificate principal with ListAttachedPolicies
// verification, plus Duplicate/NotFound negative paths.
func (r *TestRunner) runIoTPolicyTests(tc *iotTestContext) []TestResult {
	var results []TestResult
	policyName := uniqueName("policy")

	// Create a throwaway certificate to use as an attach principal.
	cert, certCleanup, certErr := tc.createCertificate(true)
	if certErr != nil {
		return []TestResult{{Service: "iot", TestName: "Policy_Setup", Status: "FAIL", Error: certErr.Error()}}
	}
	defer certCleanup()
	certARN := cert.ARN

	defer func() {
		if policyName != "" {
			tc.client.DeletePolicy(tc.ctx, &iot.DeletePolicyInput{PolicyName: aws.String(policyName)})
		}
	}()

	results = append(results, r.RunTest("iot", "Policy_CreatePolicy", func() error {
		out, err := tc.client.CreatePolicy(tc.ctx, &iot.CreatePolicyInput{
			PolicyName:     aws.String(policyName),
			PolicyDocument: aws.String(testIoTPolicyDocument),
		})
		if err != nil {
			return fmt.Errorf("CreatePolicy failed: %w", err)
		}
		if out.PolicyName == nil || *out.PolicyName != policyName {
			return fmt.Errorf("expected policyName=%s, got %v", policyName, out.PolicyName)
		}
		if out.PolicyArn == nil || *out.PolicyArn == "" {
			return fmt.Errorf("expected non-empty policyArn")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Policy_CreatePolicy_Duplicate_Conflict", func() error {
		_, err := tc.client.CreatePolicy(tc.ctx, &iot.CreatePolicyInput{
			PolicyName:     aws.String(policyName),
			PolicyDocument: aws.String(testIoTPolicyDocument),
		})
		return expectConflict(err)
	}))

	results = append(results, r.RunTest("iot", "Policy_GetPolicy", func() error {
		out, err := tc.client.GetPolicy(tc.ctx, &iot.GetPolicyInput{PolicyName: aws.String(policyName)})
		if err != nil {
			return fmt.Errorf("GetPolicy failed: %w", err)
		}
		if out.PolicyName == nil || *out.PolicyName != policyName {
			return fmt.Errorf("expected policyName=%s, got %v", policyName, out.PolicyName)
		}
		if out.PolicyDocument == nil || *out.PolicyDocument == "" {
			return fmt.Errorf("expected non-empty policyDocument")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Policy_GetPolicy_NotFound", func() error {
		_, err := tc.client.GetPolicy(tc.ctx, &iot.GetPolicyInput{PolicyName: aws.String(uniqueName("nope-policy"))})
		return expectNotFound(err)
	}))

	results = append(results, r.RunTest("iot", "Policy_ListPolicies_IncludesCreated", func() error {
		out, err := tc.client.ListPolicies(tc.ctx, &iot.ListPoliciesInput{})
		if err != nil {
			return fmt.Errorf("ListPolicies failed: %w", err)
		}
		for _, p := range out.Policies {
			if p.PolicyName != nil && *p.PolicyName == policyName {
				return nil
			}
		}
		return fmt.Errorf("%s not found in list of %d policies", policyName, len(out.Policies))
	}))

	// ── Non-default policy versions ──
	var createdVersionID string
	results = append(results, r.RunTest("iot", "PolicyVersion_CreatePolicyVersion", func() error {
		out, err := tc.client.CreatePolicyVersion(tc.ctx, &iot.CreatePolicyVersionInput{
			PolicyName:     aws.String(policyName),
			PolicyDocument: aws.String(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"iot:*","Resource":"*"}]}`),
			SetAsDefault:   false,
		})
		if err != nil {
			return fmt.Errorf("CreatePolicyVersion failed: %w", err)
		}
		if out.PolicyVersionId == nil || *out.PolicyVersionId == "" {
			return fmt.Errorf("expected non-empty policyVersionId")
		}
		createdVersionID = *out.PolicyVersionId
		return nil
	}))

	results = append(results, r.RunTest("iot", "PolicyVersion_GetPolicyVersion", func() error {
		if createdVersionID == "" {
			return fmt.Errorf("no version id captured from CreatePolicyVersion")
		}
		out, err := tc.client.GetPolicyVersion(tc.ctx, &iot.GetPolicyVersionInput{
			PolicyName:      aws.String(policyName),
			PolicyVersionId: aws.String(createdVersionID),
		})
		if err != nil {
			return fmt.Errorf("GetPolicyVersion failed: %w", err)
		}
		if aws.ToString(out.PolicyVersionId) != createdVersionID {
			return fmt.Errorf("expected versionId=%s", createdVersionID)
		}
		if out.PolicyDocument == nil || *out.PolicyDocument == "" {
			return fmt.Errorf("expected non-empty policyDocument")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "PolicyVersion_ListPolicyVersions", func() error {
		out, err := tc.client.ListPolicyVersions(tc.ctx, &iot.ListPolicyVersionsInput{PolicyName: aws.String(policyName)})
		if err != nil {
			return fmt.Errorf("ListPolicyVersions failed: %w", err)
		}
		// Default version (1) + the non-default we just created.
		if len(out.PolicyVersions) < 2 {
			return fmt.Errorf("expected at least 2 versions, got %d", len(out.PolicyVersions))
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "PolicyVersion_SetDefaultPolicyVersion", func() error {
		listOut, err := tc.client.ListPolicyVersions(tc.ctx, &iot.ListPolicyVersionsInput{PolicyName: aws.String(policyName)})
		if err != nil {
			return fmt.Errorf("ListPolicyVersions failed: %w", err)
		}
		var nonDefault string
		for _, v := range listOut.PolicyVersions {
			if v.VersionId != nil && !v.IsDefaultVersion {
				nonDefault = *v.VersionId
				break
			}
		}
		if nonDefault == "" {
			return fmt.Errorf("expected a non-default version to promote")
		}
		if _, err := tc.client.SetDefaultPolicyVersion(tc.ctx, &iot.SetDefaultPolicyVersionInput{
			PolicyName:      aws.String(policyName),
			PolicyVersionId: aws.String(nonDefault),
		}); err != nil {
			return fmt.Errorf("SetDefaultPolicyVersion failed: %w", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "PolicyVersion_DeletePolicyVersion", func() error {
		listOut, err := tc.client.ListPolicyVersions(tc.ctx, &iot.ListPolicyVersionsInput{PolicyName: aws.String(policyName)})
		if err != nil {
			return fmt.Errorf("ListPolicyVersions failed: %w", err)
		}
		var nonDefault string
		for _, v := range listOut.PolicyVersions {
			if v.VersionId != nil && !v.IsDefaultVersion {
				nonDefault = *v.VersionId
				break
			}
		}
		if nonDefault == "" {
			return nil // nothing to delete
		}
		_, err = tc.client.DeletePolicyVersion(tc.ctx, &iot.DeletePolicyVersionInput{
			PolicyName:      aws.String(policyName),
			PolicyVersionId: aws.String(nonDefault),
		})
		return err
	}))

	// ── Attach / Detach to the certificate principal ──
	results = append(results, r.RunTest("iot", "Policy_AttachPolicy", func() error {
		if _, err := tc.client.AttachPolicy(tc.ctx, &iot.AttachPolicyInput{
			PolicyName: aws.String(policyName),
			Target:     aws.String(certARN),
		}); err != nil {
			return fmt.Errorf("AttachPolicy failed: %w", err)
		}
		attached, err := tc.client.ListAttachedPolicies(tc.ctx, &iot.ListAttachedPoliciesInput{Target: aws.String(certARN)})
		if err != nil {
			return fmt.Errorf("ListAttachedPolicies failed: %w", err)
		}
		for _, p := range attached.Policies {
			if p.PolicyName != nil && *p.PolicyName == policyName {
				return nil
			}
		}
		return fmt.Errorf("policy %s not found in attached policies", policyName)
	}))

	results = append(results, r.RunTest("iot", "Policy_DetachPolicy", func() error {
		if _, err := tc.client.DetachPolicy(tc.ctx, &iot.DetachPolicyInput{
			PolicyName: aws.String(policyName),
			Target:     aws.String(certARN),
		}); err != nil {
			return fmt.Errorf("DetachPolicy failed: %w", err)
		}
		attached, err := tc.client.ListAttachedPolicies(tc.ctx, &iot.ListAttachedPoliciesInput{Target: aws.String(certARN)})
		if err != nil {
			return fmt.Errorf("ListAttachedPolicies failed: %w", err)
		}
		for _, p := range attached.Policies {
			if p.PolicyName != nil && *p.PolicyName == policyName {
				return fmt.Errorf("policy %s should no longer be attached", policyName)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Policy_DeletePolicy", func() error {
		_, err := tc.client.DeletePolicy(tc.ctx, &iot.DeletePolicyInput{PolicyName: aws.String(policyName)})
		if err != nil {
			return fmt.Errorf("DeletePolicy failed: %w", err)
		}
		policyName = "" // suppress deferred delete
		return nil
	}))

	results = append(results, r.RunTest("iot", "Policy_DeletePolicy_NotFound", func() error {
		_, err := tc.client.DeletePolicy(tc.ctx, &iot.DeletePolicyInput{PolicyName: aws.String(uniqueName("nope-policy"))})
		return expectNotFound(err)
	}))

	return results
}

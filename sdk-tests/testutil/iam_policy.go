package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func (r *TestRunner) iamPolicyTests(tc *iamTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("iam", "CreatePolicy", func() error {
		resp, err := tc.client.CreatePolicy(tc.ctx, &iam.CreatePolicyInput{
			PolicyName:     aws.String(tc.policy),
			PolicyDocument: aws.String(s3FullAccessPolicy),
		})
		if err != nil {
			return err
		}
		if resp.Policy == nil {
			return fmt.Errorf("policy is nil")
		}
		if aws.ToString(resp.Policy.Arn) == "" {
			return fmt.Errorf("policy arn is empty")
		}
		if aws.ToString(resp.Policy.PolicyName) != tc.policy {
			return fmt.Errorf("policy name mismatch: got %s, want %s", aws.ToString(resp.Policy.PolicyName), tc.policy)
		}
		tc.policyArn = *resp.Policy.Arn
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetPolicy", func() error {
		resp, err := tc.client.GetPolicy(tc.ctx, &iam.GetPolicyInput{
			PolicyArn: aws.String(tc.policyArn),
		})
		if err != nil {
			return err
		}
		if resp.Policy == nil {
			return fmt.Errorf("policy is nil")
		}
		if aws.ToString(resp.Policy.PolicyName) != tc.policy {
			return fmt.Errorf("policy name mismatch: got %v, want %s", resp.Policy.PolicyName, tc.policy)
		}
		if aws.ToString(resp.Policy.Arn) != tc.policyArn {
			return fmt.Errorf("policy arn mismatch")
		}
		if aws.ToString(resp.Policy.DefaultVersionId) == "" {
			return fmt.Errorf("default version id is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListPolicies", func() error {
		var found bool
		var marker *string
		for {
			resp, err := tc.client.ListPolicies(tc.ctx, &iam.ListPoliciesInput{
				Scope:  types.PolicyScopeTypeLocal,
				Marker: marker,
			})
			if err != nil {
				return err
			}
			for _, p := range resp.Policies {
				if aws.ToString(p.PolicyName) == tc.policy {
					found = true
					if aws.ToString(p.Arn) != tc.policyArn {
						return fmt.Errorf("policy arn mismatch in list")
					}
					break
				}
			}
			if found || !resp.IsTruncated || resp.Marker == nil {
				break
			}
			marker = resp.Marker
		}
		if !found {
			return fmt.Errorf("policy %s not found in ListPolicies", tc.policy)
		}
		return nil
	}))

	// Policy tags
	results = append(results, r.RunTest("iam", "TagPolicy", func() error {
		_, err := tc.client.TagPolicy(tc.ctx, &iam.TagPolicyInput{
			PolicyArn: aws.String(tc.policyArn),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
			},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListPolicyTags(tc.ctx, &iam.ListPolicyTagsInput{
			PolicyArn: aws.String(tc.policyArn),
		})
		if err != nil {
			return fmt.Errorf("ListPolicyTags after tag: %w", err)
		}
		if !iamTagPresent(resp.Tags, "Environment", "test") {
			return fmt.Errorf("Environment=test tag not found after TagPolicy")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListPolicyTags", func() error {
		resp, err := tc.client.ListPolicyTags(tc.ctx, &iam.ListPolicyTagsInput{
			PolicyArn: aws.String(tc.policyArn),
		})
		if err != nil {
			return err
		}
		if !iamTagPresent(resp.Tags, "Environment", "test") {
			return fmt.Errorf("Environment=test tag not found")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UntagPolicy", func() error {
		_, err := tc.client.UntagPolicy(tc.ctx, &iam.UntagPolicyInput{
			PolicyArn: aws.String(tc.policyArn),
			TagKeys:   []string{"Environment"},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListPolicyTags(tc.ctx, &iam.ListPolicyTagsInput{
			PolicyArn: aws.String(tc.policyArn),
		})
		if err != nil {
			return err
		}
		if iamTagPresent(resp.Tags, "Environment", "test") {
			return fmt.Errorf("Environment tag should be removed")
		}
		return nil
	}))

	// Attached policies — User
	results = append(results, r.RunTest("iam", "AttachUserPolicy", func() error {
		_, err := tc.client.AttachUserPolicy(tc.ctx, &iam.AttachUserPolicyInput{
			UserName:  aws.String(tc.user),
			PolicyArn: aws.String(tc.policyArn),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListAttachedUserPolicies(tc.ctx, &iam.ListAttachedUserPoliciesInput{
			UserName: aws.String(tc.user),
		})
		if err != nil {
			return fmt.Errorf("ListAttachedUserPolicies after attach: %w", err)
		}
		if !iamFindAttachedPolicy(resp.AttachedPolicies, tc.policyArn) {
			return fmt.Errorf("policy %s not found after AttachUserPolicy", tc.policyArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListAttachedUserPolicies", func() error {
		resp, err := tc.client.ListAttachedUserPolicies(tc.ctx, &iam.ListAttachedUserPoliciesInput{
			UserName: aws.String(tc.user),
		})
		if err != nil {
			return err
		}
		if !iamFindAttachedPolicy(resp.AttachedPolicies, tc.policyArn) {
			return fmt.Errorf("policy %s not found in ListAttachedUserPolicies", tc.policyArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DetachUserPolicy", func() error {
		_, err := tc.client.DetachUserPolicy(tc.ctx, &iam.DetachUserPolicyInput{
			UserName:  aws.String(tc.user),
			PolicyArn: aws.String(tc.policyArn),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListAttachedUserPolicies(tc.ctx, &iam.ListAttachedUserPoliciesInput{
			UserName: aws.String(tc.user),
		})
		if err != nil {
			return err
		}
		if iamFindAttachedPolicy(resp.AttachedPolicies, tc.policyArn) {
			return fmt.Errorf("policy should be detached from user")
		}
		return nil
	}))

	// Attached policies — Group
	results = append(results, r.RunTest("iam", "AttachGroupPolicy", func() error {
		_, err := tc.client.AttachGroupPolicy(tc.ctx, &iam.AttachGroupPolicyInput{
			GroupName: aws.String(tc.group),
			PolicyArn: aws.String(tc.policyArn),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListAttachedGroupPolicies(tc.ctx, &iam.ListAttachedGroupPoliciesInput{
			GroupName: aws.String(tc.group),
		})
		if err != nil {
			return fmt.Errorf("ListAttachedGroupPolicies after attach: %w", err)
		}
		if !iamFindAttachedPolicy(resp.AttachedPolicies, tc.policyArn) {
			return fmt.Errorf("policy %s not found after AttachGroupPolicy", tc.policyArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListAttachedGroupPolicies", func() error {
		resp, err := tc.client.ListAttachedGroupPolicies(tc.ctx, &iam.ListAttachedGroupPoliciesInput{
			GroupName: aws.String(tc.group),
		})
		if err != nil {
			return err
		}
		if !iamFindAttachedPolicy(resp.AttachedPolicies, tc.policyArn) {
			return fmt.Errorf("policy %s not found in ListAttachedGroupPolicies", tc.policyArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DetachGroupPolicy", func() error {
		_, err := tc.client.DetachGroupPolicy(tc.ctx, &iam.DetachGroupPolicyInput{
			GroupName: aws.String(tc.group),
			PolicyArn: aws.String(tc.policyArn),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListAttachedGroupPolicies(tc.ctx, &iam.ListAttachedGroupPoliciesInput{
			GroupName: aws.String(tc.group),
		})
		if err != nil {
			return fmt.Errorf("ListAttachedGroupPolicies after detach: %w", err)
		}
		if iamFindAttachedPolicy(resp.AttachedPolicies, tc.policyArn) {
			return fmt.Errorf("policy should be detached from group")
		}
		return nil
	}))

	// Attached policies — Role + ListEntitiesForPolicy
	results = append(results, r.RunTest("iam", "ListEntitiesForPolicy_Role", func() error {
		_, err := tc.client.AttachRolePolicy(tc.ctx, &iam.AttachRolePolicyInput{
			RoleName:  aws.String(tc.role),
			PolicyArn: aws.String(tc.policyArn),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListEntitiesForPolicy(tc.ctx, &iam.ListEntitiesForPolicyInput{
			PolicyArn: aws.String(tc.policyArn),
		})
		if err != nil {
			return err
		}
		if resp.PolicyRoles == nil {
			return fmt.Errorf("policy roles list is nil")
		}
		roleFound := false
		for _, pr := range resp.PolicyRoles {
			if aws.ToString(pr.RoleName) == tc.role {
				roleFound = true
				break
			}
		}
		if !roleFound {
			return fmt.Errorf("role %s not found in ListEntitiesForPolicy", tc.role)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListAttachedRolePolicies", func() error {
		resp, err := tc.client.ListAttachedRolePolicies(tc.ctx, &iam.ListAttachedRolePoliciesInput{
			RoleName: aws.String(tc.role),
		})
		if err != nil {
			return err
		}
		if !iamFindAttachedPolicy(resp.AttachedPolicies, tc.policyArn) {
			return fmt.Errorf("policy %s not found in ListAttachedRolePolicies", tc.policyArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DetachRolePolicy", func() error {
		_, err := tc.client.DetachRolePolicy(tc.ctx, &iam.DetachRolePolicyInput{
			RoleName:  aws.String(tc.role),
			PolicyArn: aws.String(tc.policyArn),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListAttachedRolePolicies(tc.ctx, &iam.ListAttachedRolePoliciesInput{
			RoleName: aws.String(tc.role),
		})
		if err != nil {
			return fmt.Errorf("ListAttachedRolePolicies after detach: %w", err)
		}
		if iamFindAttachedPolicy(resp.AttachedPolicies, tc.policyArn) {
			return fmt.Errorf("policy should be detached from role")
		}
		return nil
	}))

	// Policy versioning
	results = append(results, r.RunTest("iam", "CreatePolicyVersion", func() error {
		v2Document := `{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Action": "ec2:*",
				"Resource": "*"
			}]
		}`
		resp, err := tc.client.CreatePolicyVersion(tc.ctx, &iam.CreatePolicyVersionInput{
			PolicyArn:      aws.String(tc.policyArn),
			PolicyDocument: aws.String(v2Document),
			SetAsDefault:   false,
		})
		if err != nil {
			return err
		}
		if resp.PolicyVersion == nil {
			return fmt.Errorf("policy version is nil")
		}
		if aws.ToString(resp.PolicyVersion.VersionId) == "" {
			return fmt.Errorf("version id is empty")
		}
		if resp.PolicyVersion.IsDefaultVersion {
			return fmt.Errorf("expected non-default version")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListPolicyVersions", func() error {
		resp, err := tc.client.ListPolicyVersions(tc.ctx, &iam.ListPolicyVersionsInput{
			PolicyArn: aws.String(tc.policyArn),
		})
		if err != nil {
			return err
		}
		if len(resp.Versions) < 2 {
			return fmt.Errorf("expected at least 2 policy versions, got %d", len(resp.Versions))
		}
		defaultCount := 0
		for _, v := range resp.Versions {
			if v.IsDefaultVersion {
				defaultCount++
			}
		}
		if defaultCount != 1 {
			return fmt.Errorf("expected exactly 1 default version, got %d", defaultCount)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetPolicyVersion", func() error {
		resp, err := tc.client.ListPolicyVersions(tc.ctx, &iam.ListPolicyVersionsInput{
			PolicyArn: aws.String(tc.policyArn),
		})
		if err != nil {
			return err
		}
		var defaultVid string
		for _, v := range resp.Versions {
			if v.IsDefaultVersion && v.VersionId != nil {
				defaultVid = *v.VersionId
				break
			}
		}
		if defaultVid == "" {
			return fmt.Errorf("no default version found")
		}
		getResp, err := tc.client.GetPolicyVersion(tc.ctx, &iam.GetPolicyVersionInput{
			PolicyArn: aws.String(tc.policyArn),
			VersionId: aws.String(defaultVid),
		})
		if err != nil {
			return err
		}
		if getResp.PolicyVersion == nil {
			return fmt.Errorf("policy version is nil")
		}
		if !getResp.PolicyVersion.IsDefaultVersion {
			return fmt.Errorf("expected default version")
		}
		if aws.ToString(getResp.PolicyVersion.Document) == "" {
			return fmt.Errorf("policy document is empty in version")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "SetDefaultPolicyVersion", func() error {
		resp, err := tc.client.ListPolicyVersions(tc.ctx, &iam.ListPolicyVersionsInput{
			PolicyArn: aws.String(tc.policyArn),
		})
		if err != nil {
			return err
		}
		var nonDefaultVid string
		for _, v := range resp.Versions {
			if v.VersionId != nil && !v.IsDefaultVersion {
				nonDefaultVid = *v.VersionId
				break
			}
		}
		if nonDefaultVid == "" {
			return fmt.Errorf("no non-default version found")
		}
		_, err = tc.client.SetDefaultPolicyVersion(tc.ctx, &iam.SetDefaultPolicyVersionInput{
			PolicyArn: aws.String(tc.policyArn),
			VersionId: aws.String(nonDefaultVid),
		})
		if err != nil {
			return err
		}
		verifyResp, err := tc.client.GetPolicyVersion(tc.ctx, &iam.GetPolicyVersionInput{
			PolicyArn: aws.String(tc.policyArn),
			VersionId: aws.String(nonDefaultVid),
		})
		if err != nil {
			return err
		}
		if !verifyResp.PolicyVersion.IsDefaultVersion {
			return fmt.Errorf("version should now be default")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeletePolicyVersion", func() error {
		resp, err := tc.client.ListPolicyVersions(tc.ctx, &iam.ListPolicyVersionsInput{
			PolicyArn: aws.String(tc.policyArn),
		})
		if err != nil {
			return err
		}
		var nonDefaultVid string
		for _, v := range resp.Versions {
			if v.VersionId != nil && !v.IsDefaultVersion {
				nonDefaultVid = *v.VersionId
				break
			}
		}
		if nonDefaultVid == "" {
			return fmt.Errorf("no non-default version found to delete")
		}
		_, err = tc.client.DeletePolicyVersion(tc.ctx, &iam.DeletePolicyVersionInput{
			PolicyArn: aws.String(tc.policyArn),
			VersionId: aws.String(nonDefaultVid),
		})
		return err
	}))

	return results
}

package testutil

import (
	"fmt"
	"strings"

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
		policies, err := iamPaginate(func(marker *string) ([]types.Policy, *string, error) {
			resp, err := tc.client.ListPolicies(tc.ctx, &iam.ListPoliciesInput{
				Scope:  types.PolicyScopeTypeLocal,
				Marker: marker,
			})
			if err != nil {
				return nil, nil, err
			}
			return resp.Policies, resp.Marker, nil
		})
		if err != nil {
			return err
		}
		var matched *types.Policy
		for i := range policies {
			if aws.ToString(policies[i].PolicyName) == tc.policy {
				matched = &policies[i]
				break
			}
		}
		if matched == nil {
			return fmt.Errorf("policy %s not found in ListPolicies", tc.policy)
		}
		if aws.ToString(matched.Arn) != tc.policyArn {
			return fmt.Errorf("policy arn mismatch in list")
		}

		// Scope=Local lists customer managed policies only: the AWS managed
		// catalogue is seeded at startup, so any leak is visible here. The
		// count is not asserted — only the arn prefix membership.
		const awsManagedPrefix = "arn:aws:iam::aws:policy/"
		for i := range policies {
			if strings.HasPrefix(aws.ToString(policies[i].Arn), awsManagedPrefix) {
				return fmt.Errorf("Local listing returned AWS managed policy %s", aws.ToString(policies[i].Arn))
			}
		}

		// An omitted Scope defaults to All: the listing contains both the
		// customer managed policies and the AWS managed catalogue seeded at
		// startup.
		defaultScoped, err := iamPaginate(func(marker *string) ([]types.Policy, *string, error) {
			resp, err := tc.client.ListPolicies(tc.ctx, &iam.ListPoliciesInput{
				Marker: marker,
			})
			if err != nil {
				return nil, nil, err
			}
			return resp.Policies, resp.Marker, nil
		})
		if err != nil {
			return err
		}
		hasAWSManaged, hasCustomer := false, false
		for i := range defaultScoped {
			arn := aws.ToString(defaultScoped[i].Arn)
			if strings.HasPrefix(arn, awsManagedPrefix) {
				hasAWSManaged = true
			}
			if arn == aws.ToString(matched.Arn) {
				hasCustomer = true
			}
		}
		if !hasAWSManaged {
			return fmt.Errorf("default-scope listing returned no AWS managed policies")
		}
		if !hasCustomer {
			return fmt.Errorf("default-scope listing omitted the customer managed policy %s", aws.ToString(matched.Arn))
		}

		// Scope=AWS lists only the AWS managed catalogue.
		awsScoped, err := iamPaginate(func(marker *string) ([]types.Policy, *string, error) {
			resp, err := tc.client.ListPolicies(tc.ctx, &iam.ListPoliciesInput{
				Scope:  types.PolicyScopeTypeAws,
				Marker: marker,
			})
			if err != nil {
				return nil, nil, err
			}
			return resp.Policies, resp.Marker, nil
		})
		if err != nil {
			return err
		}
		if len(awsScoped) == 0 {
			return fmt.Errorf("AWS listing returned no policies")
		}
		for i := range awsScoped {
			if !strings.HasPrefix(aws.ToString(awsScoped[i].Arn), awsManagedPrefix) {
				return fmt.Errorf("AWS listing returned non-AWS managed policy %s", aws.ToString(awsScoped[i].Arn))
			}
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

	// AWS-managed policy attach/detach
	results = append(results, r.RunTest("iam", "AttachUserPolicy_AWSManaged", func() error {
		awsManagedArn := "arn:aws:iam::aws:policy/ReadOnlyAccess"
		_, err := tc.client.AttachUserPolicy(tc.ctx, &iam.AttachUserPolicyInput{
			UserName:  aws.String(tc.user),
			PolicyArn: aws.String(awsManagedArn),
		})
		if err != nil {
			return fmt.Errorf("AttachUserPolicy with AWS-managed ARN failed: %w", err)
		}
		resp, err := tc.client.ListAttachedUserPolicies(tc.ctx, &iam.ListAttachedUserPoliciesInput{
			UserName: aws.String(tc.user),
		})
		if err != nil {
			return err
		}
		if !iamFindAttachedPolicy(resp.AttachedPolicies, awsManagedArn) {
			return fmt.Errorf("AWS-managed policy %s not found after attach", awsManagedArn)
		}
		// Detach to clean up
		_, _ = tc.client.DetachUserPolicy(tc.ctx, &iam.DetachUserPolicyInput{
			UserName:  aws.String(tc.user),
			PolicyArn: aws.String(awsManagedArn),
		})
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

	// SetAsDefault=true exercises the 2-step PutVersion + SetDefaultVersion
	// path.  The invariant "exactly one default version per policy" must
	// hold after the swap (regression guard).
	results = append(results, r.RunTest("iam", "CreatePolicyVersion_SetAsDefault", func() error {
		v3Document := `{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Action": "s3:*",
				"Resource": "*"
			}]
		}`
		resp, err := tc.client.CreatePolicyVersion(tc.ctx, &iam.CreatePolicyVersionInput{
			PolicyArn:      aws.String(tc.policyArn),
			PolicyDocument: aws.String(v3Document),
			SetAsDefault:   true,
		})
		if err != nil {
			return err
		}
		if resp.PolicyVersion == nil {
			return fmt.Errorf("policy version is nil")
		}
		if !resp.PolicyVersion.IsDefaultVersion {
			return fmt.Errorf("expected new version to be the default")
		}
		newVersionId := aws.ToString(resp.PolicyVersion.VersionId)

		// Re-list and confirm exactly one default remains, and it is the
		// newly-created version.
		listResp, err := tc.client.ListPolicyVersions(tc.ctx, &iam.ListPolicyVersionsInput{
			PolicyArn: aws.String(tc.policyArn),
		})
		if err != nil {
			return err
		}
		defaultCount := 0
		defaultVersionId := ""
		for _, v := range listResp.Versions {
			if v.IsDefaultVersion {
				defaultCount++
				defaultVersionId = aws.ToString(v.VersionId)
			}
		}
		if defaultCount != 1 {
			return fmt.Errorf("expected exactly 1 default version after SetAsDefault, got %d", defaultCount)
		}
		if defaultVersionId != newVersionId {
			return fmt.Errorf("expected default version %s, got %s", newVersionId, defaultVersionId)
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

	results = append(results, r.RunTest("iam", "AWSManagedPolicy_Protection", func() error {
		arn := fmt.Sprintf("arn:aws:iam::aws:policy/ReadOnlyAccess")

		resp, err := tc.client.CreatePolicyVersion(tc.ctx, &iam.CreatePolicyVersionInput{
			PolicyArn:      aws.String(arn),
			PolicyDocument: aws.String(iamAllowPolicy("s3:GetObject")),
			SetAsDefault:   false,
		})
		if err == nil {
			// Remove the accidental version so the seeded policy stays
			// intact even when the guard is missing.
			if resp.PolicyVersion != nil && resp.PolicyVersion.VersionId != nil {
				_, _ = tc.client.DeletePolicyVersion(tc.ctx, &iam.DeletePolicyVersionInput{
					PolicyArn: aws.String(arn),
					VersionId: resp.PolicyVersion.VersionId,
				})
			}
			return fmt.Errorf("CreatePolicyVersion on an AWS managed policy must be rejected")
		}
		if !isInvalidInputError(err) {
			return fmt.Errorf("CreatePolicyVersion on an AWS managed policy: got %v, want InvalidInput", err)
		}

		versions, err := tc.client.ListPolicyVersions(tc.ctx, &iam.ListPolicyVersionsInput{PolicyArn: aws.String(arn)})
		if err != nil {
			return fmt.Errorf("ListPolicyVersions on the AWS managed policy: %w", err)
		}
		defaultVersion := ""
		for _, v := range versions.Versions {
			if v.IsDefaultVersion && v.VersionId != nil {
				defaultVersion = *v.VersionId
				break
			}
		}
		if defaultVersion == "" {
			return fmt.Errorf("no default version found on the AWS managed policy")
		}

		if _, err := tc.client.SetDefaultPolicyVersion(tc.ctx, &iam.SetDefaultPolicyVersionInput{
			PolicyArn: aws.String(arn),
			VersionId: aws.String(defaultVersion),
		}); err == nil || !isInvalidInputError(err) {
			return fmt.Errorf("SetDefaultPolicyVersion on an AWS managed policy: got %v, want InvalidInput", err)
		}

		if _, err := tc.client.DeletePolicyVersion(tc.ctx, &iam.DeletePolicyVersionInput{
			PolicyArn: aws.String(arn),
			VersionId: aws.String(defaultVersion),
		}); err == nil || !isInvalidInputError(err) {
			return fmt.Errorf("DeletePolicyVersion on an AWS managed policy: got %v, want InvalidInput", err)
		}

		if _, err := tc.client.DeletePolicy(tc.ctx, &iam.DeletePolicyInput{
			PolicyArn: aws.String(arn),
		}); err == nil || !isInvalidInputError(err) {
			return fmt.Errorf("DeletePolicy on an AWS managed policy: got %v, want InvalidInput", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "AWSManagedPolicy_MissingArn_NoSuchEntity", func() error {
		// An AWS-managed ARN that does not exist is a missing resource,
		// not an unmodifiable one: existence is reported first.
		arn := fmt.Sprintf("arn:aws:iam::aws:policy/NeverSeeded-%s", tc.ts)

		if _, err := tc.client.DeletePolicy(tc.ctx, &iam.DeletePolicyInput{
			PolicyArn: aws.String(arn),
		}); err == nil || !containsErrorCode(err, "NoSuchEntity") {
			return fmt.Errorf("DeletePolicy on an unseeded AWS-managed ARN: got %v, want NoSuchEntity", err)
		}

		if _, err := tc.client.CreatePolicyVersion(tc.ctx, &iam.CreatePolicyVersionInput{
			PolicyArn:      aws.String(arn),
			PolicyDocument: aws.String(iamAllowPolicy("s3:GetObject")),
			SetAsDefault:   false,
		}); err == nil || !containsErrorCode(err, "NoSuchEntity") {
			return fmt.Errorf("CreatePolicyVersion on an unseeded AWS-managed ARN: got %v, want NoSuchEntity", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "CreatePolicy_DescriptionTooLong", func() error {
		_, err := tc.client.CreatePolicy(tc.ctx, &iam.CreatePolicyInput{
			PolicyName:     aws.String(fmt.Sprintf("LongDesc-%s", tc.ts)),
			PolicyDocument: aws.String(iamAllowPolicy("s3:GetObject")),
			Description:    aws.String(strings.Repeat("d", 1001)),
		})
		if err == nil {
			return fmt.Errorf("a 1001-character description must be rejected")
		}
		if !isInvalidInputError(err) {
			return fmt.Errorf("long description: got %v, want InvalidInput", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "SimulatePrincipalPolicy", func() error {
		user := fmt.Sprintf("SimUser-%s", tc.ts)
		cleanupUser, err := tc.createUser(user)
		if err != nil {
			return err
		}
		defer cleanupUser()

		simArn, cleanupSim, err := tc.createPolicy(fmt.Sprintf("SimPolicy-%s", tc.ts), iamAllowPolicy("s3:ListBucket"))
		if err != nil {
			return err
		}
		defer cleanupSim()
		if _, err := tc.client.AttachUserPolicy(tc.ctx, &iam.AttachUserPolicyInput{
			UserName:  aws.String(user),
			PolicyArn: aws.String(simArn),
		}); err != nil {
			return err
		}
		defer tc.client.DetachUserPolicy(tc.ctx, &iam.DetachUserPolicyInput{
			UserName:  aws.String(user),
			PolicyArn: aws.String(simArn),
		})

		userArn := fmt.Sprintf("arn:aws:iam::%s:user/%s", tc.accountID, user)
		resp, err := tc.client.SimulatePrincipalPolicy(tc.ctx, &iam.SimulatePrincipalPolicyInput{
			PolicySourceArn: aws.String(userArn),
			ActionNames:     []string{"s3:ListBucket", "dynamodb:ListTables"},
			ResourceArns:    []string{"*"},
		})
		if err != nil {
			return err
		}
		if len(resp.EvaluationResults) != 2 {
			return fmt.Errorf("evaluation results: got %d, want 2", len(resp.EvaluationResults))
		}
		decisions := map[string]string{}
		for _, ev := range resp.EvaluationResults {
			decisions[aws.ToString(ev.EvalActionName)] = string(ev.EvalDecision)
		}
		if decisions["s3:ListBucket"] != "allowed" {
			return fmt.Errorf("s3:ListBucket decision: got %s, want allowed", decisions["s3:ListBucket"])
		}
		if decisions["dynamodb:ListTables"] == "allowed" {
			return fmt.Errorf("dynamodb:ListTables must not be allowed")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "SimulatePrincipalPolicy_ResourceSpecificResults", func() error {
		user := fmt.Sprintf("SimRsr-%s", tc.ts)
		cleanupUser, err := tc.createUser(user)
		if err != nil {
			return err
		}
		defer cleanupUser()

		userArn := fmt.Sprintf("arn:aws:iam::%s:user/%s", tc.accountID, user)
		allowDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:ListBucket"],"Resource":["arn:aws:s3:::bucket-a"]}]}`
		resp, err := tc.client.SimulatePrincipalPolicy(tc.ctx, &iam.SimulatePrincipalPolicyInput{
			PolicySourceArn: aws.String(userArn),
			ActionNames:     []string{"s3:ListBucket"},
			ResourceArns:    []string{"arn:aws:s3:::bucket-a", "arn:aws:s3:::bucket-b"},
			PolicyInputList: []string{allowDoc},
		})
		if err != nil {
			return err
		}
		if len(resp.EvaluationResults) != 1 {
			return fmt.Errorf("evaluation results: got %d, want 1 (one per action)", len(resp.EvaluationResults))
		}
		ev := resp.EvaluationResults[0]
		if string(ev.EvalDecision) != "implicitDeny" {
			return fmt.Errorf("aggregate decision: got %s, want implicitDeny (bucket-b denies)", ev.EvalDecision)
		}
		if aws.ToString(ev.EvalResourceName) != "*" {
			return fmt.Errorf("top-level resource name: got %s, want * (the ARN-template slot)", aws.ToString(ev.EvalResourceName))
		}
		if len(ev.ResourceSpecificResults) != 2 {
			return fmt.Errorf("resource-specific results: got %d, want 2", len(ev.ResourceSpecificResults))
		}
		rsr := ev.ResourceSpecificResults[0]
		if aws.ToString(rsr.EvalResourceName) != "arn:aws:s3:::bucket-a" || string(rsr.EvalResourceDecision) != "allowed" {
			return fmt.Errorf("bucket-a resource result: got %s/%s", aws.ToString(rsr.EvalResourceName), rsr.EvalResourceDecision)
		}
		if len(rsr.MatchedStatements) != 1 || aws.ToString(rsr.MatchedStatements[0].SourcePolicyId) != "PolicyInputList.1" {
			return fmt.Errorf("bucket-a matched statements: %+v", rsr.MatchedStatements)
		}
		if rsr2 := ev.ResourceSpecificResults[1]; string(rsr2.EvalResourceDecision) != "implicitDeny" {
			return fmt.Errorf("bucket-b resource decision: got %s, want implicitDeny", rsr2.EvalResourceDecision)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "SimulatePrincipalPolicy_Pagination", func() error {
		user := fmt.Sprintf("SimPage-%s", tc.ts)
		cleanupUser, err := tc.createUser(user)
		if err != nil {
			return err
		}
		defer cleanupUser()

		userArn := fmt.Sprintf("arn:aws:iam::%s:user/%s", tc.accountID, user)
		page1, err := tc.client.SimulatePrincipalPolicy(tc.ctx, &iam.SimulatePrincipalPolicyInput{
			PolicySourceArn: aws.String(userArn),
			ActionNames:     []string{"s3:ListBucket", "sqs:ListQueues"},
			MaxItems:        aws.Int32(1),
		})
		if err != nil {
			return err
		}
		if len(page1.EvaluationResults) != 1 || !page1.IsTruncated || aws.ToString(page1.Marker) == "" {
			return fmt.Errorf("first page: got %d results, truncated=%v, marker=%q", len(page1.EvaluationResults), page1.IsTruncated, aws.ToString(page1.Marker))
		}
		if got := aws.ToString(page1.EvaluationResults[0].EvalActionName); got != "s3:ListBucket" {
			return fmt.Errorf("first page action: got %s, want s3:ListBucket", got)
		}
		page2, err := tc.client.SimulatePrincipalPolicy(tc.ctx, &iam.SimulatePrincipalPolicyInput{
			PolicySourceArn: aws.String(userArn),
			ActionNames:     []string{"s3:ListBucket", "sqs:ListQueues"},
			MaxItems:        aws.Int32(1),
			Marker:          page1.Marker,
		})
		if err != nil {
			return err
		}
		if len(page2.EvaluationResults) != 1 || page2.IsTruncated {
			return fmt.Errorf("second page: got %d results, truncated=%v", len(page2.EvaluationResults), page2.IsTruncated)
		}
		if got := aws.ToString(page2.EvaluationResults[0].EvalActionName); got != "sqs:ListQueues" {
			return fmt.Errorf("second page action: got %s, want sqs:ListQueues", got)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "SimulatePrincipalPolicy_ResourcePolicyAndExclusions", func() error {
		user := fmt.Sprintf("SimRp-%s", tc.ts)
		cleanupUser, err := tc.createUser(user)
		if err != nil {
			return err
		}
		defer cleanupUser()

		userArn := fmt.Sprintf("arn:aws:iam::%s:user/%s", tc.accountID, user)
		rpDoc := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":"%s"},"Action":["s3:PutObject"],"Resource":["arn:aws:s3:::mary/*"]}]}`, userArn)
		allowAll := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:*"],"Resource":["*"]}]}`
		identityArn, cleanupIdentity, err := tc.createPolicy(fmt.Sprintf("SimRpPol-%s", tc.ts), allowAll)
		if err != nil {
			return err
		}
		defer cleanupIdentity()
		if _, err := tc.client.AttachUserPolicy(tc.ctx, &iam.AttachUserPolicyInput{
			UserName:  aws.String(user),
			PolicyArn: aws.String(identityArn),
		}); err != nil {
			return err
		}
		defer tc.client.DetachUserPolicy(tc.ctx, &iam.DetachUserPolicyInput{
			UserName:  aws.String(user),
			PolicyArn: aws.String(identityArn),
		})

		// Same account: the resource policy alone allows.
		resp, err := tc.client.SimulatePrincipalPolicy(tc.ctx, &iam.SimulatePrincipalPolicyInput{
			PolicySourceArn: aws.String(userArn),
			ActionNames:     []string{"s3:PutObject"},
			ResourceArns:    []string{"arn:aws:s3:::mary/Test"},
			ResourcePolicy:  aws.String(rpDoc),
		})
		if err != nil {
			return err
		}
		if string(resp.EvaluationResults[0].EvalDecision) != "allowed" {
			return fmt.Errorf("same-account resource policy decision: got %s, want allowed", resp.EvaluationResults[0].EvalDecision)
		}
		// Per-resource matched statements carry every matched statement:
		// the attached identity policy's allow and the resource policy's.
		rsr := resp.EvaluationResults[0].ResourceSpecificResults[0]
		byType := map[string]int{}
		for _, st := range rsr.MatchedStatements {
			byType[string(st.SourcePolicyType)]++
		}
		if byType["user-managed"] != 1 || byType["resource"] != 1 {
			return fmt.Errorf("same-account matched statements: %+v", rsr.MatchedStatements)
		}

		// Cross account with an identity allow: intersection allows; the
		// decision details name every policy type.
		cross, err := tc.client.SimulatePrincipalPolicy(tc.ctx, &iam.SimulatePrincipalPolicyInput{
			PolicySourceArn: aws.String(userArn),
			ActionNames:     []string{"s3:PutObject"},
			ResourceArns:    []string{"arn:aws:s3:::mary/Test"},
			ResourcePolicy:  aws.String(rpDoc),
			ResourceOwner:   aws.String("999999999999"),
		})
		if err != nil {
			return err
		}
		ev := cross.EvaluationResults[0]
		if string(ev.EvalDecision) != "allowed" {
			return fmt.Errorf("cross-account decision: got %s, want allowed", ev.EvalDecision)
		}
		if v, ok := ev.EvalDecisionDetails["Resource Policy"]; !ok || string(v) != "allowed" {
			return fmt.Errorf("cross-account decision details: %v", ev.EvalDecisionDetails)
		}

		// Excluding the attached identity policy (the exclusion list
		// addresses the principal's gathered policies, never the
		// caller-supplied input documents) leaves the cross-account
		// intersection denied.
		excluded, err := tc.client.SimulatePrincipalPolicy(tc.ctx, &iam.SimulatePrincipalPolicyInput{
			PolicySourceArn: aws.String(userArn),
			ActionNames:     []string{"s3:PutObject"},
			ResourceArns:    []string{"arn:aws:s3:::mary/Test"},
			ResourcePolicy:  aws.String(rpDoc),
			ResourceOwner:   aws.String("999999999999"),
			PolicyExclusionList: []types.PolicyIdentifier{&types.PolicyIdentifierMemberPolicyType{
				Value: types.PolicyIdentifierPolicyTypeUserManaged,
			}},
		})
		if err != nil {
			return err
		}
		if string(excluded.EvaluationResults[0].EvalDecision) != "implicitDeny" {
			return fmt.Errorf("excluded identity decision: got %s, want implicitDeny", excluded.EvaluationResults[0].EvalDecision)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "SimulateCustomPolicy", func() error {
		allowDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:ListBucket"],"Resource":["*"]}]}`
		resp, err := tc.client.SimulateCustomPolicy(tc.ctx, &iam.SimulateCustomPolicyInput{
			PolicyInputList: []string{allowDoc},
			ActionNames:     []string{"s3:ListBucket", "sqs:ListQueues"},
		})
		if err != nil {
			return err
		}
		if len(resp.EvaluationResults) != 2 {
			return fmt.Errorf("evaluation results: got %d, want 2", len(resp.EvaluationResults))
		}
		decisions := map[string]string{}
		for _, ev := range resp.EvaluationResults {
			decisions[aws.ToString(ev.EvalActionName)] = string(ev.EvalDecision)
		}
		if decisions["s3:ListBucket"] != "allowed" || decisions["sqs:ListQueues"] != "implicitDeny" {
			return fmt.Errorf("custom decisions: got %v", decisions)
		}

		// An empty PolicyInputList is rejected client-side; a malformed
		// document is rejected server-side.
		if _, err := tc.client.SimulateCustomPolicy(tc.ctx, &iam.SimulateCustomPolicyInput{
			PolicyInputList: []string{"{not a policy document"},
			ActionNames:     []string{"s3:ListBucket"},
		}); err == nil {
			return fmt.Errorf("a malformed custom policy document must be rejected")
		}

		// An ordered organisation policy list is a hierarchy of levels, each
		// carrying SCP documents. A level that omits the action caps the
		// decision and is reported through the organisation decision detail;
		// cross-account, the SCP deny also suppresses the decision details.
		scpUser := fmt.Sprintf("SimCustScp-%s", tc.ts)
		cleanupScpUser, err := tc.createUser(scpUser)
		if err != nil {
			return err
		}
		defer cleanupScpUser()
		scpUserArn := fmt.Sprintf("arn:aws:iam::%s:user/%s", tc.accountID, scpUser)
		allowBoth := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["s3:ListBucket","sqs:ListQueues"],"Resource":["*"]}]}`
		allowSQS := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["sqs:*"],"Resource":["*"]}]}`
		rpDoc := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":"%s"},"Action":["s3:ListBucket"],"Resource":["arn:aws:s3:::bucket"]}]}`, scpUserArn)
		cross, err := tc.client.SimulateCustomPolicy(tc.ctx, &iam.SimulateCustomPolicyInput{
			PolicyInputList: []string{allowBoth},
			OrderedOrganizationPolicyInputList: []types.OrderedOrganizationPolicyType{
				{ServiceControlPolicyInputList: []string{allowSQS}},
			},
			ActionNames:    []string{"s3:ListBucket", "sqs:ListQueues"},
			CallerArn:      aws.String(scpUserArn),
			ResourceArns:   []string{"arn:aws:s3:::bucket"},
			ResourcePolicy: aws.String(rpDoc),
			ResourceOwner:  aws.String("999999999999"),
		})
		if err != nil {
			return err
		}
		byAction := map[string]types.EvaluationResult{}
		for _, ev := range cross.EvaluationResults {
			byAction[aws.ToString(ev.EvalActionName)] = ev
		}
		capped := byAction["s3:ListBucket"]
		if string(capped.EvalDecision) != "implicitDeny" {
			return fmt.Errorf("SCP-capped decision: got %s, want implicitDeny", capped.EvalDecision)
		}
		if capped.OrganizationsDecisionDetail == nil || capped.OrganizationsDecisionDetail.AllowedByOrganizations {
			return fmt.Errorf("SCP-capped org detail: %+v, want AllowedByOrganizations=false", capped.OrganizationsDecisionDetail)
		}
		if len(capped.EvalDecisionDetails) != 0 {
			return fmt.Errorf("SCP-denied decision details must be suppressed, got %v", capped.EvalDecisionDetails)
		}
		inList := byAction["sqs:ListQueues"]
		if inList.OrganizationsDecisionDetail == nil || !inList.OrganizationsDecisionDetail.AllowedByOrganizations {
			return fmt.Errorf("in-list org detail: %+v, want AllowedByOrganizations=true", inList.OrganizationsDecisionDetail)
		}
		if len(inList.EvalDecisionDetails) == 0 {
			return fmt.Errorf("cross-account decision details without an SCP deny must be reported")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListPoliciesGrantingServiceAccess", func() error {
		group := fmt.Sprintf("LPGSA-grp-%s", tc.ts)
		user := fmt.Sprintf("LPGSA-usr-%s", tc.ts)
		cleanupGroup, err := tc.createGroup(group)
		if err != nil {
			return err
		}
		defer cleanupGroup()
		cleanupUser, err := tc.createUser(user)
		if err != nil {
			return err
		}
		defer cleanupUser()
		if _, err := tc.client.AddUserToGroup(tc.ctx, &iam.AddUserToGroupInput{GroupName: aws.String(group), UserName: aws.String(user)}); err != nil {
			return err
		}
		defer tc.client.RemoveUserFromGroup(tc.ctx, &iam.RemoveUserFromGroupInput{GroupName: aws.String(group), UserName: aws.String(user)})

		usrPolArn, cleanupUsrPol, err := tc.createPolicy(fmt.Sprintf("LPGSA-usr-%s", tc.ts), iamAllowPolicy("s3:ListBucket"))
		if err != nil {
			return err
		}
		defer cleanupUsrPol()
		if _, err := tc.client.AttachUserPolicy(tc.ctx, &iam.AttachUserPolicyInput{
			UserName:  aws.String(user),
			PolicyArn: aws.String(usrPolArn),
		}); err != nil {
			return err
		}
		defer tc.client.DetachUserPolicy(tc.ctx, &iam.DetachUserPolicyInput{
			UserName:  aws.String(user),
			PolicyArn: aws.String(usrPolArn),
		})

		grpPolArn, cleanupGrpPol, err := tc.createPolicy(fmt.Sprintf("LPGSA-grp-%s", tc.ts), iamAllowPolicy("ec2:*"))
		if err != nil {
			return err
		}
		defer cleanupGrpPol()
		if _, err := tc.client.AttachGroupPolicy(tc.ctx, &iam.AttachGroupPolicyInput{
			GroupName: aws.String(group),
			PolicyArn: aws.String(grpPolArn),
		}); err != nil {
			return err
		}
		defer tc.client.DetachGroupPolicy(tc.ctx, &iam.DetachGroupPolicyInput{
			GroupName: aws.String(group),
			PolicyArn: aws.String(grpPolArn),
		})

		if _, err := tc.client.PutGroupPolicy(tc.ctx, &iam.PutGroupPolicyInput{
			GroupName:      aws.String(group),
			PolicyName:     aws.String("lpgsa-group-inline"),
			PolicyDocument: aws.String(iamAllowPolicy("logs:CreateLogGroup")),
		}); err != nil {
			return err
		}
		defer tc.client.DeleteGroupPolicy(tc.ctx, &iam.DeleteGroupPolicyInput{
			GroupName:  aws.String(group),
			PolicyName: aws.String("lpgsa-group-inline"),
		})

		if _, err := tc.client.PutUserPolicy(tc.ctx, &iam.PutUserPolicyInput{
			UserName:       aws.String(user),
			PolicyName:     aws.String("lpgsa-user-inline"),
			PolicyDocument: aws.String(iamAllowPolicy("iam:ListUsers")),
		}); err != nil {
			return err
		}
		defer tc.client.DeleteUserPolicy(tc.ctx, &iam.DeleteUserPolicyInput{
			UserName:   aws.String(user),
			PolicyName: aws.String("lpgsa-user-inline"),
		})

		userArn := fmt.Sprintf("arn:aws:iam::%s:user/%s", tc.accountID, user)
		namespaces := []string{"s3", "ec2", "logs", "iam", "sns"}
		resp, err := tc.client.ListPoliciesGrantingServiceAccess(tc.ctx, &iam.ListPoliciesGrantingServiceAccessInput{
			Arn:               aws.String(userArn),
			ServiceNamespaces: namespaces,
		})
		if err != nil {
			return err
		}
		if len(resp.PoliciesGrantingServiceAccess) != len(namespaces) {
			return fmt.Errorf("expected one entry per requested namespace, got %d", len(resp.PoliciesGrantingServiceAccess))
		}
		byNS := map[string][]types.PolicyGrantingServiceAccess{}
		for i, e := range resp.PoliciesGrantingServiceAccess {
			if aws.ToString(e.ServiceNamespace) != namespaces[i] {
				return fmt.Errorf("entry %d namespace: got %s, want %s (request order)", i, aws.ToString(e.ServiceNamespace), namespaces[i])
			}
			byNS[aws.ToString(e.ServiceNamespace)] = e.Policies
		}

		s3p := byNS["s3"]
		if len(s3p) != 1 || aws.ToString(s3p[0].PolicyArn) != usrPolArn ||
			string(s3p[0].PolicyType) != "MANAGED" || aws.ToString(s3p[0].PolicyName) != fmt.Sprintf("LPGSA-usr-%s", tc.ts) {
			return fmt.Errorf("s3 namespace entries: %+v", s3p)
		}
		ec2p := byNS["ec2"]
		if len(ec2p) != 1 || aws.ToString(ec2p[0].PolicyArn) != grpPolArn || string(ec2p[0].PolicyType) != "MANAGED" {
			return fmt.Errorf("ec2 namespace entries: %+v", ec2p)
		}
		logsp := byNS["logs"]
		if len(logsp) != 1 || string(logsp[0].PolicyType) != "INLINE" || aws.ToString(logsp[0].PolicyName) != "lpgsa-group-inline" ||
			aws.ToString(logsp[0].EntityName) != group || string(logsp[0].EntityType) != "GROUP" {
			return fmt.Errorf("logs namespace entries: %+v", logsp)
		}
		iamp := byNS["iam"]
		if len(iamp) != 1 || string(iamp[0].PolicyType) != "INLINE" || aws.ToString(iamp[0].PolicyName) != "lpgsa-user-inline" ||
			aws.ToString(iamp[0].EntityName) != user || string(iamp[0].EntityType) != "USER" {
			return fmt.Errorf("iam namespace entries: %+v", iamp)
		}
		if len(byNS["sns"]) != 0 {
			return fmt.Errorf("sns namespace must have no granting policies: %+v", byNS["sns"])
		}
		if resp.IsTruncated {
			return fmt.Errorf("a single page must hold every entry")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListPoliciesGrantingServiceAccess_RoleAndBoundary", func() error {
		role := fmt.Sprintf("LPGSA-role-%s", tc.ts)
		cleanupRole, err := tc.createRole(role)
		if err != nil {
			return err
		}
		defer cleanupRole()

		rolePolArn, cleanupRolePol, err := tc.createPolicy(fmt.Sprintf("LPGSA-role-%s", tc.ts), iamAllowPolicy("sns:Publish"))
		if err != nil {
			return err
		}
		defer cleanupRolePol()
		if _, err := tc.client.AttachRolePolicy(tc.ctx, &iam.AttachRolePolicyInput{
			RoleName:  aws.String(role),
			PolicyArn: aws.String(rolePolArn),
		}); err != nil {
			return err
		}
		defer tc.client.DetachRolePolicy(tc.ctx, &iam.DetachRolePolicyInput{
			RoleName:  aws.String(role),
			PolicyArn: aws.String(rolePolArn),
		})

		resp, err := tc.client.ListPoliciesGrantingServiceAccess(tc.ctx, &iam.ListPoliciesGrantingServiceAccessInput{
			Arn:               aws.String(fmt.Sprintf("arn:aws:iam::%s:role/%s", tc.accountID, role)),
			ServiceNamespaces: []string{"sns"},
		})
		if err != nil {
			return err
		}
		if len(resp.PoliciesGrantingServiceAccess) != 1 || len(resp.PoliciesGrantingServiceAccess[0].Policies) != 1 {
			return fmt.Errorf("role sns entry shape unexpected: %+v", resp.PoliciesGrantingServiceAccess)
		}
		roleEntry := resp.PoliciesGrantingServiceAccess[0].Policies[0]
		if aws.ToString(roleEntry.PolicyArn) != rolePolArn || string(roleEntry.PolicyType) != "MANAGED" {
			return fmt.Errorf("role sns entry: %+v", roleEntry)
		}

		// A managed policy attached only as a permissions boundary is not
		// reported as granting service access.
		user := fmt.Sprintf("LPGSA-bnd-%s", tc.ts)
		cleanupUser, err := tc.createUser(user)
		if err != nil {
			return err
		}
		defer cleanupUser()
		bndArn, cleanupBnd, err := tc.createPolicy(fmt.Sprintf("LPGSA-bnd-%s", tc.ts), iamAllowPolicy("sqs:SendMessage"))
		if err != nil {
			return err
		}
		defer cleanupBnd()
		if _, err := tc.client.PutUserPermissionsBoundary(tc.ctx, &iam.PutUserPermissionsBoundaryInput{
			UserName:            aws.String(user),
			PermissionsBoundary: aws.String(bndArn),
		}); err != nil {
			return err
		}
		defer tc.client.DeleteUserPermissionsBoundary(tc.ctx, &iam.DeleteUserPermissionsBoundaryInput{UserName: aws.String(user)})

		userArn := fmt.Sprintf("arn:aws:iam::%s:user/%s", tc.accountID, user)
		bndResp, err := tc.client.ListPoliciesGrantingServiceAccess(tc.ctx, &iam.ListPoliciesGrantingServiceAccessInput{
			Arn:               aws.String(userArn),
			ServiceNamespaces: []string{"sqs"},
		})
		if err != nil {
			return err
		}
		if len(bndResp.PoliciesGrantingServiceAccess) != 1 || len(bndResp.PoliciesGrantingServiceAccess[0].Policies) != 0 {
			return fmt.Errorf("permissions boundaries must not be reported: %+v", bndResp.PoliciesGrantingServiceAccess)
		}

		// A missing identity is rejected. The missing-ServiceNamespaces
		// path cannot be exercised through the SDK: the client validates
		// the required member before the request is sent; the server-side
		// rejection is pinned by a unit test instead.
		if _, err := tc.client.ListPoliciesGrantingServiceAccess(tc.ctx, &iam.ListPoliciesGrantingServiceAccessInput{
			Arn:               aws.String(fmt.Sprintf("arn:aws:iam::%s:user/LPGSA-nobody-%s", tc.accountID, tc.ts)),
			ServiceNamespaces: []string{"s3"},
		}); err == nil || !containsErrorCode(err, "NoSuchEntity") {
			return fmt.Errorf("unknown identity ARN: got %v, want NoSuchEntity", err)
		}
		return nil
	}))

	// The policyDescriptionType bound is a Smithy @length trait counted in
	// Unicode characters; an 800-character CJK description is 2400 bytes
	// but must be accepted.
	results = append(results, r.RunTest("iam", "CreatePolicy_DescriptionMultibyteAccepted", func() error {
		name := tc.policy + "-mb-desc"
		desc := strings.Repeat("\u65e5", 800)
		resp, err := tc.client.CreatePolicy(tc.ctx, &iam.CreatePolicyInput{
			PolicyName:     aws.String(name),
			PolicyDocument: aws.String(s3FullAccessPolicy),
			Description:    aws.String(desc),
		})
		if err != nil {
			return fmt.Errorf("create with multibyte description: %v", err)
		}
		arn := aws.ToString(resp.Policy.Arn)
		defer tc.client.DeletePolicy(tc.ctx, &iam.DeletePolicyInput{PolicyArn: aws.String(arn)})

		getResp, err := tc.client.GetPolicy(tc.ctx, &iam.GetPolicyInput{PolicyArn: aws.String(arn)})
		if err != nil {
			return err
		}
		if aws.ToString(getResp.Policy.Description) != desc {
			return fmt.Errorf("multibyte description not persisted faithfully")
		}
		return nil
	}))

	return results
}

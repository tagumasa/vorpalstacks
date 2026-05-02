package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func (r *TestRunner) iamRoleTests(tc *iamTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("iam", "CreateRole", func() error {
		resp, err := tc.client.CreateRole(tc.ctx, &iam.CreateRoleInput{
			RoleName:                 aws.String(tc.role),
			AssumeRolePolicyDocument: aws.String(assumeRolePolicy),
			Description:              aws.String("Test role for SDK tests"),
			MaxSessionDuration:       aws.Int32(43200),
		})
		if err != nil {
			return err
		}
		if resp.Role == nil {
			return fmt.Errorf("role is nil")
		}
		if aws.ToString(resp.Role.RoleName) != tc.role {
			return fmt.Errorf("role name mismatch: got %s, want %s", aws.ToString(resp.Role.RoleName), tc.role)
		}
		if aws.ToString(resp.Role.Arn) == "" {
			return fmt.Errorf("role arn is empty")
		}
		if aws.ToString(resp.Role.Description) != "Test role for SDK tests" {
			return fmt.Errorf("description mismatch: got %s", aws.ToString(resp.Role.Description))
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "CreateRole_InvalidName", func() error {
		_, err := tc.client.CreateRole(tc.ctx, &iam.CreateRoleInput{
			RoleName:                 aws.String("invalid:role-name"),
			AssumeRolePolicyDocument: aws.String(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":["lambda.amazonaws.com"]},"Action":["sts:AssumeRole"}]}`),
		})
		if err == nil {
			return fmt.Errorf("expected error for invalid role name with colon")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetRole", func() error {
		resp, err := tc.client.GetRole(tc.ctx, &iam.GetRoleInput{
			RoleName: aws.String(tc.role),
		})
		if err != nil {
			return err
		}
		if resp.Role == nil {
			return fmt.Errorf("role is nil")
		}
		if aws.ToString(resp.Role.RoleName) != tc.role {
			return fmt.Errorf("role name mismatch")
		}
		if aws.ToString(resp.Role.AssumeRolePolicyDocument) == "" {
			return fmt.Errorf("assume role policy document is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListRoles", func() error {
		var found bool
		var marker *string
		for {
			resp, err := tc.client.ListRoles(tc.ctx, &iam.ListRolesInput{Marker: marker})
			if err != nil {
				return err
			}
			for _, r := range resp.Roles {
				if aws.ToString(r.RoleName) == tc.role {
					found = true
					if aws.ToString(r.Arn) == "" {
						return fmt.Errorf("role arn is empty in list")
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
			return fmt.Errorf("role %s not found in ListRoles", tc.role)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UpdateRoleDescription", func() error {
		_, err := tc.client.UpdateRoleDescription(tc.ctx, &iam.UpdateRoleDescriptionInput{
			RoleName:    aws.String(tc.role),
			Description: aws.String("Updated role description"),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetRole(tc.ctx, &iam.GetRoleInput{
			RoleName: aws.String(tc.role),
		})
		if err != nil {
			return err
		}
		if aws.ToString(resp.Role.Description) != "Updated role description" {
			return fmt.Errorf("description not updated: got %s", aws.ToString(resp.Role.Description))
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UpdateRole", func() error {
		_, err := tc.client.UpdateRole(tc.ctx, &iam.UpdateRoleInput{
			RoleName:           aws.String(tc.role),
			MaxSessionDuration: aws.Int32(3600),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetRole(tc.ctx, &iam.GetRoleInput{
			RoleName: aws.String(tc.role),
		})
		if err != nil {
			return err
		}
		if aws.ToInt32(resp.Role.MaxSessionDuration) != 3600 {
			return fmt.Errorf("max session duration not updated: got %d", aws.ToInt32(resp.Role.MaxSessionDuration))
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UpdateAssumeRolePolicy", func() error {
		newTrustPolicy := `{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {"Service": "ec2.amazonaws.com"},
				"Action": "sts:AssumeRole"
			}]
		}`
		_, err := tc.client.UpdateAssumeRolePolicy(tc.ctx, &iam.UpdateAssumeRolePolicyInput{
			RoleName:       aws.String(tc.role),
			PolicyDocument: aws.String(newTrustPolicy),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetRole(tc.ctx, &iam.GetRoleInput{
			RoleName: aws.String(tc.role),
		})
		if err != nil {
			return err
		}
		if aws.ToString(resp.Role.AssumeRolePolicyDocument) == "" {
			return fmt.Errorf("assume role policy document is empty after update")
		}
		return nil
	}))

	// Role tags
	results = append(results, r.RunTest("iam", "TagRole", func() error {
		_, err := tc.client.TagRole(tc.ctx, &iam.TagRoleInput{
			RoleName: aws.String(tc.role),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "ListRoleTags", func() error {
		resp, err := tc.client.ListRoleTags(tc.ctx, &iam.ListRoleTagsInput{
			RoleName: aws.String(tc.role),
		})
		if err != nil {
			return err
		}
		if !iamTagPresent(resp.Tags, "Environment", "test") {
			return fmt.Errorf("Environment=test tag not found")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UntagRole", func() error {
		_, err := tc.client.UntagRole(tc.ctx, &iam.UntagRoleInput{
			RoleName: aws.String(tc.role),
			TagKeys:  []string{"Environment"},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListRoleTags(tc.ctx, &iam.ListRoleTagsInput{
			RoleName: aws.String(tc.role),
		})
		if err != nil {
			return err
		}
		if iamTagPresent(resp.Tags, "Environment", "test") {
			return fmt.Errorf("Environment tag should be removed")
		}
		return nil
	}))

	// Permissions boundary: deferred to iamPermissionsBoundaryTests (needs policyArn)

	// Inline policies
	results = append(results, r.RunTest("iam", "PutRolePolicy", func() error {
		_, err := tc.client.PutRolePolicy(tc.ctx, &iam.PutRolePolicyInput{
			RoleName:       aws.String(tc.role),
			PolicyName:     aws.String(tc.roleInlinePolicy),
			PolicyDocument: aws.String(logsFullAccessPolicy),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "GetRolePolicy", func() error {
		resp, err := tc.client.GetRolePolicy(tc.ctx, &iam.GetRolePolicyInput{
			RoleName:   aws.String(tc.role),
			PolicyName: aws.String(tc.roleInlinePolicy),
		})
		if err != nil {
			return err
		}
		if resp.PolicyDocument == nil || *resp.PolicyDocument == "" {
			return fmt.Errorf("policy document is empty")
		}
		if aws.ToString(resp.RoleName) != tc.role {
			return fmt.Errorf("role name mismatch in GetRolePolicy")
		}
		if aws.ToString(resp.PolicyName) != tc.roleInlinePolicy {
			return fmt.Errorf("policy name mismatch in GetRolePolicy")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListRolePolicies", func() error {
		resp, err := tc.client.ListRolePolicies(tc.ctx, &iam.ListRolePoliciesInput{
			RoleName: aws.String(tc.role),
		})
		if err != nil {
			return err
		}
		found := false
		for _, name := range resp.PolicyNames {
			if name == tc.roleInlinePolicy {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("inline policy %s not found in ListRolePolicies", tc.roleInlinePolicy)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteRolePolicy", func() error {
		_, err := tc.client.DeleteRolePolicy(tc.ctx, &iam.DeleteRolePolicyInput{
			RoleName:   aws.String(tc.role),
			PolicyName: aws.String(tc.roleInlinePolicy),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListRolePolicies(tc.ctx, &iam.ListRolePoliciesInput{
			RoleName: aws.String(tc.role),
		})
		if err != nil {
			return err
		}
		for _, name := range resp.PolicyNames {
			if name == tc.roleInlinePolicy {
				return fmt.Errorf("inline policy should be deleted")
			}
		}
		return nil
	}))

	// Service-linked role
	results = append(results, r.RunTest("iam", "CreateServiceLinkedRole", func() error {
		resp, err := tc.client.CreateServiceLinkedRole(tc.ctx, &iam.CreateServiceLinkedRoleInput{
			AWSServiceName: aws.String("lambda.amazonaws.com"),
			Description:    aws.String("Test service-linked role"),
		})
		if err != nil {
			return err
		}
		if resp.Role == nil {
			return fmt.Errorf("role is nil")
		}
		if resp.Role.RoleName == nil {
			return fmt.Errorf("role name is nil")
		}
		tc.svcLinkedRoleName = *resp.Role.RoleName
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteServiceLinkedRole", func() error {
		resp, err := tc.client.DeleteServiceLinkedRole(tc.ctx, &iam.DeleteServiceLinkedRoleInput{
			RoleName: aws.String(tc.svcLinkedRoleName),
		})
		if err != nil {
			return err
		}
		if resp.DeletionTaskId == nil {
			return fmt.Errorf("deletion task id is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetServiceLinkedRoleDeletionStatus", func() error {
		_, err := tc.client.GetServiceLinkedRoleDeletionStatus(tc.ctx, &iam.GetServiceLinkedRoleDeletionStatusInput{
			DeletionTaskId: aws.String("test-task-id"),
		})
		return err
	}))

	return results
}

package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func (r *TestRunner) iamRoleTests(tc *iamTestContext) []TestResult {
	var results []TestResult

	// exampleRoleTemplateArn addresses the Example role template — the one
	// catalogue entry whose full content AWS documents (API Reference
	// examples of GetRoleTemplateVersion and AcquireRole).
	const exampleRoleTemplateArn = "arn:aws:iam::aws:role-template/awsserviceprincipal/Example:1"

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
		roles, err := iamPaginate(func(marker *string) ([]types.Role, *string, error) {
			resp, err := tc.client.ListRoles(tc.ctx, &iam.ListRolesInput{Marker: marker})
			if err != nil {
				return nil, nil, err
			}
			return resp.Roles, resp.Marker, nil
		})
		if err != nil {
			return err
		}
		var matched *types.Role
		for i := range roles {
			if aws.ToString(roles[i].RoleName) == tc.role {
				matched = &roles[i]
				break
			}
		}
		if matched == nil {
			return fmt.Errorf("role %s not found in ListRoles", tc.role)
		}
		if aws.ToString(matched.Arn) == "" {
			return fmt.Errorf("role arn is empty in list")
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
		_, err := tc.client.UpdateAssumeRolePolicy(tc.ctx, &iam.UpdateAssumeRolePolicyInput{
			RoleName:       aws.String(tc.role),
			PolicyDocument: aws.String(ec2AssumeRolePolicy),
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
		if err != nil {
			return err
		}
		resp, err := tc.client.ListRoleTags(tc.ctx, &iam.ListRoleTagsInput{
			RoleName: aws.String(tc.role),
		})
		if err != nil {
			return fmt.Errorf("ListRoleTags after tag: %w", err)
		}
		if !iamTagPresent(resp.Tags, "Environment", "test") {
			return fmt.Errorf("Environment=test tag not found after TagRole")
		}
		return nil
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
		if err != nil {
			return err
		}
		resp, err := tc.client.GetRolePolicy(tc.ctx, &iam.GetRolePolicyInput{
			RoleName:   aws.String(tc.role),
			PolicyName: aws.String(tc.roleInlinePolicy),
		})
		if err != nil {
			return fmt.Errorf("GetRolePolicy after PutRolePolicy: %w", err)
		}
		if resp.PolicyDocument == nil || *resp.PolicyDocument == "" {
			return fmt.Errorf("policy document is empty after PutRolePolicy")
		}
		return nil
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
		tc.deletionTaskId = *resp.DeletionTaskId
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetServiceLinkedRoleDeletionStatus", func() error {
		// A task id longer than the documented 1000-character maximum is
		// malformed input, not a missing task.
		if _, err := tc.client.GetServiceLinkedRoleDeletionStatus(tc.ctx, &iam.GetServiceLinkedRoleDeletionStatusInput{
			DeletionTaskId: aws.String(strings.Repeat("t", 1001)),
		}); err == nil {
			return fmt.Errorf("an over-length task id must be rejected")
		} else if !containsErrorCode(err, "InvalidInput") {
			return fmt.Errorf("over-length task id: got %v, want InvalidInput", err)
		}

		resp, err := tc.client.GetServiceLinkedRoleDeletionStatus(tc.ctx, &iam.GetServiceLinkedRoleDeletionStatusInput{
			DeletionTaskId: aws.String(tc.deletionTaskId),
		})
		if err != nil {
			return err
		}
		if resp.Status == "" {
			return fmt.Errorf("deletion status is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "CreateServiceLinkedRole_Naming", func() error {
		resp, err := tc.client.CreateServiceLinkedRole(tc.ctx, &iam.CreateServiceLinkedRoleInput{
			AWSServiceName: aws.String("ecs.amazonaws.com"),
		})
		if err != nil {
			return err
		}
		if resp.Role == nil {
			return fmt.Errorf("role is nil")
		}
		if aws.ToString(resp.Role.RoleName) != "AWSServiceRoleForECS" {
			return fmt.Errorf("service-linked role name: got %s, want AWSServiceRoleForECS", aws.ToString(resp.Role.RoleName))
		}
		if aws.ToString(resp.Role.Path) != "/aws-service-role/ecs.amazonaws.com/" {
			return fmt.Errorf("service-linked role path: got %s", aws.ToString(resp.Role.Path))
		}
		_, err = tc.client.DeleteServiceLinkedRole(tc.ctx, &iam.DeleteServiceLinkedRoleInput{
			RoleName: aws.String("AWSServiceRoleForECS"),
		})
		if err != nil {
			return fmt.Errorf("cleanup DeleteServiceLinkedRole: %w", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetRoleTemplateVersion", func() error {
		resp, err := tc.client.GetRoleTemplateVersion(tc.ctx, &iam.GetRoleTemplateVersionInput{
			TemplateArn: aws.String(exampleRoleTemplateArn),
		})
		if err != nil {
			return err
		}
		template := resp.RoleTemplateVersion
		if template == nil {
			return fmt.Errorf("role template version is nil")
		}
		if aws.ToString(template.TemplateName) != "Example" {
			return fmt.Errorf("template name: got %s, want Example", aws.ToString(template.TemplateName))
		}
		if aws.ToInt32(template.MajorVersion) != 1 || aws.ToInt32(template.MinorVersion) != 1 || aws.ToInt32(template.DefaultMinorVersion) != 1 {
			return fmt.Errorf("template versions: got major %d minor %d default %d, want 1/1/1",
				aws.ToInt32(template.MajorVersion), aws.ToInt32(template.MinorVersion), aws.ToInt32(template.DefaultMinorVersion))
		}
		if !template.Enabled {
			return fmt.Errorf("template Enabled: got false, want true")
		}
		if aws.ToString(template.RoleNamePattern) != "Example-@{Department}" {
			return fmt.Errorf("role name pattern: got %s", aws.ToString(template.RoleNamePattern))
		}
		if aws.ToString(template.RolePathPattern) != "/awsserviceprincipal/" {
			return fmt.Errorf("role path pattern: got %s", aws.ToString(template.RolePathPattern))
		}
		if !strings.Contains(aws.ToString(template.AssumeRolePolicyDocumentTemplate), "ec2.amazonaws.com") {
			return fmt.Errorf("trust policy template does not grant ec2.amazonaws.com: %s", aws.ToString(template.AssumeRolePolicyDocumentTemplate))
		}
		if len(template.ParametersDefinition) != 1 || aws.ToString(template.ParametersDefinition[0].Name) != "Department" {
			return fmt.Errorf("parameters definition: got %v, want one Department parameter", template.ParametersDefinition)
		}

		// An unknown template ARN is rejected with NoSuchEntity.
		_, err = tc.client.GetRoleTemplateVersion(tc.ctx, &iam.GetRoleTemplateVersionInput{
			TemplateArn: aws.String("arn:aws:iam::aws:role-template/awsserviceprincipal/Absent:1"),
		})
		if err == nil {
			return fmt.Errorf("an unknown template ARN must be rejected")
		}
		if !containsErrorCode(err, "NoSuchEntity") {
			return fmt.Errorf("unknown template ARN: got %v, want NoSuchEntity", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "AcquireRole", func() error {
		resp, err := tc.client.AcquireRole(tc.ctx, &iam.AcquireRoleInput{
			TemplateArn: aws.String(exampleRoleTemplateArn),
			ReplacementValues: map[string]types.ReplacementValueEntry{
				"Department": {Values: []string{"Engineering"}},
			},
		})
		if err != nil {
			return err
		}
		// The acquired role must not outlive the test: every assertion
		// failure below returns early, and a leaked Example-Engineering
		// makes the retry collide before the assertions run.
		defer tc.client.DeleteRole(tc.ctx, &iam.DeleteRoleInput{
			RoleName: aws.String("Example-Engineering"),
		})
		role := resp.Role
		if role == nil {
			return fmt.Errorf("role is nil")
		}
		if aws.ToString(role.RoleName) != "Example-Engineering" {
			return fmt.Errorf("acquired role name: got %s, want Example-Engineering", aws.ToString(role.RoleName))
		}
		if aws.ToString(role.Path) != "/awsserviceprincipal/" {
			return fmt.Errorf("acquired role path: got %s, want /awsserviceprincipal/", aws.ToString(role.Path))
		}
		if !strings.Contains(aws.ToString(role.AssumeRolePolicyDocument), "ec2.amazonaws.com") {
			return fmt.Errorf("acquired trust policy does not grant ec2.amazonaws.com: %s", aws.ToString(role.AssumeRolePolicyDocument))
		}
		if role.SourceRoleTemplate == nil {
			return fmt.Errorf("the acquired role must record its source template")
		}
		if aws.ToString(role.SourceRoleTemplate.TemplateArn) != exampleRoleTemplateArn {
			return fmt.Errorf("source template arn: got %s", aws.ToString(role.SourceRoleTemplate.TemplateArn))
		}
		if aws.ToInt32(role.SourceRoleTemplate.TemplateMinorVersion) != 1 {
			return fmt.Errorf("source template minor version: got %d, want 1 (the resolved default)", aws.ToInt32(role.SourceRoleTemplate.TemplateMinorVersion))
		}
		got, err := tc.client.GetRole(tc.ctx, &iam.GetRoleInput{RoleName: aws.String("Example-Engineering")})
		if err != nil {
			return err
		}
		if got.Role == nil || got.Role.SourceRoleTemplate == nil {
			return fmt.Errorf("GetRole must return the recorded source template")
		}
		if aws.ToString(got.Role.SourceRoleTemplate.TemplateArn) != exampleRoleTemplateArn {
			return fmt.Errorf("GetRole source template arn: got %s", aws.ToString(got.Role.SourceRoleTemplate.TemplateArn))
		}

		// A replacement value list beyond the documented bound of 20 is
		// rejected.
		many := make([]string, 21)
		for i := range many {
			many[i] = "v"
		}
		_, err = tc.client.AcquireRole(tc.ctx, &iam.AcquireRoleInput{
			TemplateArn: aws.String(exampleRoleTemplateArn),
			ReplacementValues: map[string]types.ReplacementValueEntry{
				"Department": {Values: []string{"Overflow"}},
				"Other":      {Values: many},
			},
		})
		if err == nil {
			return fmt.Errorf("a 21-value replacement list must be rejected")
		}
		if !isInvalidInputError(err) {
			return fmt.Errorf("over-length replacement list: got %v, want InvalidInput", err)
		}

		// Acquiring the same template with the same parameter again
		// collides with the existing role name.
		_, err = tc.client.AcquireRole(tc.ctx, &iam.AcquireRoleInput{
			TemplateArn: aws.String(exampleRoleTemplateArn),
			ReplacementValues: map[string]types.ReplacementValueEntry{
				"Department": {Values: []string{"Engineering"}},
			},
		})
		if err == nil {
			return fmt.Errorf("a role name collision must be rejected")
		}
		if !containsErrorCode(err, "NameConflict") {
			return fmt.Errorf("role name collision: got %v, want NameConflict", err)
		}

		// A required template parameter must be supplied.
		_, err = tc.client.AcquireRole(tc.ctx, &iam.AcquireRoleInput{
			TemplateArn: aws.String(exampleRoleTemplateArn),
		})
		if err == nil {
			return fmt.Errorf("a missing required template parameter must be rejected")
		}
		if !isInvalidInputError(err) {
			return fmt.Errorf("missing template parameter: got %v, want InvalidInput", err)
		}
		return nil
	}))

	return results
}

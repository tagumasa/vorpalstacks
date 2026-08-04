package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func (r *TestRunner) iamEdgeTests(tc *iamTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("iam", "Error_DeleteNonExistentUser", func() error {
		_, err := tc.client.DeleteUser(tc.ctx, &iam.DeleteUserInput{
			UserName: aws.String("NonExistentUser-" + tc.ts),
		})
		if err := AssertErrorContains(err, "NoSuchEntity"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "Error_GetNonExistentRole", func() error {
		_, err := tc.client.GetRole(tc.ctx, &iam.GetRoleInput{
			RoleName: aws.String("NonExistentRole-" + tc.ts),
		})
		if err := AssertErrorContains(err, "NoSuchEntity"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "Error_AttachPolicyToNonExistentUser", func() error {
		_, err := tc.client.AttachUserPolicy(tc.ctx, &iam.AttachUserPolicyInput{
			UserName:  aws.String("NonExistentUser-" + tc.ts),
			PolicyArn: aws.String(tc.policyArn),
		})
		if err := AssertErrorContains(err, "NoSuchEntity"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "Error_DeleteDefaultPolicyVersion", func() error {
		resp, err := tc.client.ListPolicyVersions(tc.ctx, &iam.ListPolicyVersionsInput{
			PolicyArn: aws.String(tc.policyArn),
		})
		if err != nil {
			return fmt.Errorf("failed to list policy versions: %v", err)
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
		_, err = tc.client.DeletePolicyVersion(tc.ctx, &iam.DeletePolicyVersionInput{
			PolicyArn: aws.String(tc.policyArn),
			VersionId: aws.String(defaultVid),
		})
		if err == nil {
			return fmt.Errorf("expected error when deleting default policy version")
		}
		if !strings.Contains(err.Error(), "InvalidInput") && !strings.Contains(err.Error(), "DeleteConflict") {
			return fmt.Errorf("expected InvalidInput or DeleteConflict error, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "Error_CreateDuplicateUser", func() error {
		_, err := tc.client.CreateUser(tc.ctx, &iam.CreateUserInput{
			UserName: aws.String(tc.user),
		})
		if err := AssertErrorContains(err, "EntityAlreadyExists"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "Error_CreateDuplicatePolicy", func() error {
		_, err := tc.client.CreatePolicy(tc.ctx, &iam.CreatePolicyInput{
			PolicyName:     aws.String(tc.policy),
			PolicyDocument: aws.String(s3FullAccessPolicy),
		})
		if err := AssertErrorContains(err, "EntityAlreadyExists"); err != nil {
			return err
		}
		return nil
	}))

	return results
}

func (r *TestRunner) iamPaginationTests(tc *iamTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("iam", "ListUsers_Pagination", func() error {
		pgTs := tc.ts
		var pgUsers []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("PagUser-%s-%d", pgTs, i)
			_, err := tc.client.CreateUser(tc.ctx, &iam.CreateUserInput{UserName: aws.String(name)})
			if err != nil {
				return fmt.Errorf("create user %s: %v", name, err)
			}
			pgUsers = append(pgUsers, name)
		}

		var allUsers []types.User
		var marker *string
		for {
			resp, err := tc.client.ListUsers(tc.ctx, &iam.ListUsersInput{
				PathPrefix: aws.String("/"),
				Marker:     marker,
				MaxItems:   aws.Int32(2),
			})
			if err != nil {
				for _, name := range pgUsers {
					tc.client.DeleteUser(tc.ctx, &iam.DeleteUserInput{UserName: aws.String(name)})
				}
				return fmt.Errorf("list users page: %v", err)
			}
			for _, u := range resp.Users {
				if strings.HasPrefix(aws.ToString(u.UserName), "PagUser-"+pgTs) {
					allUsers = append(allUsers, u)
				}
			}
			if resp.IsTruncated && resp.Marker != nil {
				marker = resp.Marker
			} else {
				break
			}
		}

		for _, name := range pgUsers {
			tc.client.DeleteUser(tc.ctx, &iam.DeleteUserInput{UserName: aws.String(name)})
		}

		if len(allUsers) != 5 {
			return fmt.Errorf("expected 5 paginated users, got %d", len(allUsers))
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListRoles_Pagination", func() error {
		pgTs := tc.ts
		var pgRoles []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("PagRole-%s-%d", pgTs, i)
			_, err := tc.client.CreateRole(tc.ctx, &iam.CreateRoleInput{
				RoleName:                 aws.String(name),
				AssumeRolePolicyDocument: aws.String(assumeRolePolicy),
			})
			if err != nil {
				return fmt.Errorf("create role %s: %v", name, err)
			}
			pgRoles = append(pgRoles, name)
		}

		var allRoles []types.Role
		var marker *string
		for {
			resp, err := tc.client.ListRoles(tc.ctx, &iam.ListRolesInput{
				PathPrefix: aws.String("/"),
				Marker:     marker,
				MaxItems:   aws.Int32(2),
			})
			if err != nil {
				for _, name := range pgRoles {
					tc.client.DeleteRole(tc.ctx, &iam.DeleteRoleInput{RoleName: aws.String(name)})
				}
				return fmt.Errorf("list roles page: %v", err)
			}
			for _, r := range resp.Roles {
				if strings.HasPrefix(aws.ToString(r.RoleName), "PagRole-"+pgTs) {
					allRoles = append(allRoles, r)
				}
			}
			if resp.IsTruncated && resp.Marker != nil {
				marker = resp.Marker
			} else {
				break
			}
		}

		for _, name := range pgRoles {
			tc.client.DeleteRole(tc.ctx, &iam.DeleteRoleInput{RoleName: aws.String(name)})
		}

		if len(allRoles) != 5 {
			return fmt.Errorf("expected 5 paginated roles, got %d", len(allRoles))
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListPolicies_Pagination", func() error {
		pgTs := tc.ts
		var pgPolicyArns []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("PagPolicy-%s-%d", pgTs, i)
			resp, err := tc.client.CreatePolicy(tc.ctx, &iam.CreatePolicyInput{
				PolicyName:     aws.String(name),
				PolicyDocument: aws.String(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"*","Resource":"*"}]}`),
			})
			if err != nil {
				return fmt.Errorf("create policy %s: %v", name, err)
			}
			pgPolicyArns = append(pgPolicyArns, aws.ToString(resp.Policy.Arn))
		}

		var allPolicies []types.Policy
		var marker *string
		for {
			resp, err := tc.client.ListPolicies(tc.ctx, &iam.ListPoliciesInput{
				Scope:    types.PolicyScopeTypeLocal,
				Marker:   marker,
				MaxItems: aws.Int32(2),
			})
			if err != nil {
				for _, arn := range pgPolicyArns {
					tc.client.DeletePolicy(tc.ctx, &iam.DeletePolicyInput{PolicyArn: aws.String(arn)})
				}
				return fmt.Errorf("list policies page: %v", err)
			}
			for _, p := range resp.Policies {
				if strings.HasPrefix(aws.ToString(p.PolicyName), "PagPolicy-"+pgTs) {
					allPolicies = append(allPolicies, p)
				}
			}
			if resp.IsTruncated && resp.Marker != nil {
				marker = resp.Marker
			} else {
				break
			}
		}

		for _, arn := range pgPolicyArns {
			tc.client.DeletePolicy(tc.ctx, &iam.DeletePolicyInput{PolicyArn: aws.String(arn)})
		}

		if len(allPolicies) != 5 {
			return fmt.Errorf("expected 5 paginated policies, got %d", len(allPolicies))
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "Error_DeleteInlinePolicyNonExistentUser", func() error {
		_, err := tc.client.DeleteUserPolicy(tc.ctx, &iam.DeleteUserPolicyInput{
			UserName:   aws.String("NonExistentUser-" + tc.ts),
			PolicyName: aws.String("SomePolicy"),
		})
		if err := AssertErrorContains(err, "NoSuchEntity"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "Error_CreateUserInvalidName", func() error {
		_, err := tc.client.CreateUser(tc.ctx, &iam.CreateUserInput{
			UserName: aws.String("invalid user!@#$"),
		})
		if err := AssertErrorContains(err, "InvalidInput"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "Error_CreateGroupInvalidName", func() error {
		_, err := tc.client.CreateGroup(tc.ctx, &iam.CreateGroupInput{
			GroupName: aws.String("invalid group!@#$"),
		})
		if err := AssertErrorContains(err, "InvalidInput"); err != nil {
			return err
		}
		return nil
	}))

	// CreateUser with invalid path must be rejected.
	results = append(results, r.RunTest("iam", "Error_CreateUserInvalidPath", func() error {
		_, err := tc.client.CreateUser(tc.ctx, &iam.CreateUserInput{
			UserName: aws.String("edge-user-" + tc.ts),
			Path:     aws.String("invalid-no-slashes"),
		})
		if err := AssertErrorContains(err, "InvalidInput"); err != nil {
			return err
		}
		return nil
	}))

	// CreateRole with invalid path must be rejected.
	results = append(results, r.RunTest("iam", "Error_CreateRoleInvalidPath", func() error {
		_, err := tc.client.CreateRole(tc.ctx, &iam.CreateRoleInput{
			RoleName:                 aws.String("edge-role-" + tc.ts),
			AssumeRolePolicyDocument: aws.String(assumeRolePolicy),
			Path:                     aws.String("invalid-no-slashes"),
		})
		if err := AssertErrorContains(err, "InvalidInput"); err != nil {
			return err
		}
		return nil
	}))

	// CreatePolicy with empty Action string must be rejected.
	results = append(results, r.RunTest("iam", "Error_CreatePolicyEmptyAction", func() error {
		_, err := tc.client.CreatePolicy(tc.ctx, &iam.CreatePolicyInput{
			PolicyName:     aws.String("edge-policy-" + tc.ts),
			PolicyDocument: aws.String(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"","Resource":"*"}]}`),
		})
		if err == nil {
			return fmt.Errorf("expected error for empty Action in policy")
		}
		return nil
	}))

	// CreatePolicy with Action and NotAction both present must be rejected.
	results = append(results, r.RunTest("iam", "Error_CreatePolicyActionNotActionExclusive", func() error {
		_, err := tc.client.CreatePolicy(tc.ctx, &iam.CreatePolicyInput{
			PolicyName:     aws.String("edge-policy-ex-" + tc.ts),
			PolicyDocument: aws.String(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"s3:*","NotAction":"ec2:*","Resource":"*"}]}`),
		})
		if err == nil {
			return fmt.Errorf("expected error for Action + NotAction in policy")
		}
		return nil
	}))

	// DeleteServiceLinkedRole on a non-service-linked role must return DeleteConflict.
	results = append(results, r.RunTest("iam", "Error_DeleteServiceLinkedRoleNonSLR", func() error {
		roleName := "EdgeSLRTest-" + tc.ts
		_, err := tc.client.CreateRole(tc.ctx, &iam.CreateRoleInput{
			RoleName:                 aws.String(roleName),
			AssumeRolePolicyDocument: aws.String(assumeRolePolicy),
		})
		if err != nil {
			return fmt.Errorf("setup CreateRole failed: %w", err)
		}
		_, err = tc.client.DeleteServiceLinkedRole(tc.ctx, &iam.DeleteServiceLinkedRoleInput{
			RoleName: aws.String(roleName),
		})
		if err == nil {
			tc.client.DeleteRole(tc.ctx, &iam.DeleteRoleInput{RoleName: aws.String(roleName)})
			return fmt.Errorf("expected error when deleting non-service-linked role via DeleteServiceLinkedRole")
		}
		if !strings.Contains(err.Error(), "DeleteConflict") && !strings.Contains(err.Error(), "ConflictException") {
			tc.client.DeleteRole(tc.ctx, &iam.DeleteRoleInput{RoleName: aws.String(roleName)})
			return fmt.Errorf("expected DeleteConflict error, got: %v", err)
		}
		tc.client.DeleteRole(tc.ctx, &iam.DeleteRoleInput{RoleName: aws.String(roleName)})
		return nil
	}))

	// DeleteAccountAlias with wrong alias name must return NoSuchEntity.
	results = append(results, r.RunTest("iam", "Error_DeleteAccountAliasMismatch", func() error {
		_, err := tc.client.DeleteAccountAlias(tc.ctx, &iam.DeleteAccountAliasInput{
			AccountAlias: aws.String("nonexistent-alias-" + tc.ts),
		})
		if err := AssertErrorContains(err, "NoSuchEntity"); err != nil {
			return err
		}
		return nil
	}))

	// Deleting a user that has a permissions boundary attached must
	// decrement the referenced policy's PermissionsBoundaryUsageCount.
	results = append(results, r.RunTest("iam", "DeleteUser_DecrementsPermissionsBoundaryUsageCount", func() error {
		policyName := "EdgePBDriftUser-" + tc.ts
		createOut, err := tc.client.CreatePolicy(tc.ctx, &iam.CreatePolicyInput{
			PolicyName:     aws.String(policyName),
			PolicyDocument: aws.String(s3FullAccessPolicy),
		})
		if err != nil {
			return fmt.Errorf("setup CreatePolicy failed: %w", err)
		}
		policyArn := aws.ToString(createOut.Policy.Arn)
		defer tc.client.DeletePolicy(tc.ctx, &iam.DeletePolicyInput{PolicyArn: aws.String(policyArn)})

		userName := "EdgePBDriftUser-" + tc.ts
		if _, err := tc.client.CreateUser(tc.ctx, &iam.CreateUserInput{
			UserName: aws.String(userName),
		}); err != nil {
			return fmt.Errorf("setup CreateUser failed: %w", err)
		}

		if _, err := tc.client.PutUserPermissionsBoundary(tc.ctx, &iam.PutUserPermissionsBoundaryInput{
			UserName:            aws.String(userName),
			PermissionsBoundary: aws.String(policyArn),
		}); err != nil {
			tc.client.DeleteUser(tc.ctx, &iam.DeleteUserInput{UserName: aws.String(userName)})
			return fmt.Errorf("PutUserPermissionsBoundary failed: %w", err)
		}

		// Verify the counter is incremented before deletion.
		getBefore, err := tc.client.GetPolicy(tc.ctx, &iam.GetPolicyInput{PolicyArn: aws.String(policyArn)})
		if err != nil {
			return fmt.Errorf("GetPolicy (before) failed: %w", err)
		}
		if aws.ToInt32(getBefore.Policy.PermissionsBoundaryUsageCount) != 1 {
			return fmt.Errorf("expected PermissionsBoundaryUsageCount=1 before delete, got %d", aws.ToInt32(getBefore.Policy.PermissionsBoundaryUsageCount))
		}

		if _, err := tc.client.DeleteUser(tc.ctx, &iam.DeleteUserInput{UserName: aws.String(userName)}); err != nil {
			return fmt.Errorf("DeleteUser failed: %w", err)
		}

		// After deletion the counter must return to 0.
		getAfter, err := tc.client.GetPolicy(tc.ctx, &iam.GetPolicyInput{PolicyArn: aws.String(policyArn)})
		if err != nil {
			return fmt.Errorf("GetPolicy (after) failed: %w", err)
		}
		if aws.ToInt32(getAfter.Policy.PermissionsBoundaryUsageCount) != 0 {
			return fmt.Errorf("expected PermissionsBoundaryUsageCount=0 after delete, got %d", aws.ToInt32(getAfter.Policy.PermissionsBoundaryUsageCount))
		}
		return nil
	}))

	// Deleting a role that has a permissions boundary attached must
	// decrement the referenced policy's PermissionsBoundaryUsageCount.
	results = append(results, r.RunTest("iam", "DeleteRole_DecrementsPermissionsBoundaryUsageCount", func() error {
		policyName := "EdgePBDriftRole-" + tc.ts
		createOut, err := tc.client.CreatePolicy(tc.ctx, &iam.CreatePolicyInput{
			PolicyName:     aws.String(policyName),
			PolicyDocument: aws.String(s3FullAccessPolicy),
		})
		if err != nil {
			return fmt.Errorf("setup CreatePolicy failed: %w", err)
		}
		policyArn := aws.ToString(createOut.Policy.Arn)
		defer tc.client.DeletePolicy(tc.ctx, &iam.DeletePolicyInput{PolicyArn: aws.String(policyArn)})

		roleName := "EdgePBDriftRole-" + tc.ts
		if _, err := tc.client.CreateRole(tc.ctx, &iam.CreateRoleInput{
			RoleName:                 aws.String(roleName),
			AssumeRolePolicyDocument: aws.String(assumeRolePolicy),
		}); err != nil {
			return fmt.Errorf("setup CreateRole failed: %w", err)
		}

		if _, err := tc.client.PutRolePermissionsBoundary(tc.ctx, &iam.PutRolePermissionsBoundaryInput{
			RoleName:            aws.String(roleName),
			PermissionsBoundary: aws.String(policyArn),
		}); err != nil {
			tc.client.DeleteRole(tc.ctx, &iam.DeleteRoleInput{RoleName: aws.String(roleName)})
			return fmt.Errorf("PutRolePermissionsBoundary failed: %w", err)
		}

		// Verify the counter is incremented before deletion.
		getBefore, err := tc.client.GetPolicy(tc.ctx, &iam.GetPolicyInput{PolicyArn: aws.String(policyArn)})
		if err != nil {
			return fmt.Errorf("GetPolicy (before) failed: %w", err)
		}
		if aws.ToInt32(getBefore.Policy.PermissionsBoundaryUsageCount) != 1 {
			return fmt.Errorf("expected PermissionsBoundaryUsageCount=1 before delete, got %d", aws.ToInt32(getBefore.Policy.PermissionsBoundaryUsageCount))
		}

		if _, err := tc.client.DeleteRole(tc.ctx, &iam.DeleteRoleInput{RoleName: aws.String(roleName)}); err != nil {
			return fmt.Errorf("DeleteRole failed: %w", err)
		}

		// After deletion the counter must return to 0.
		getAfter, err := tc.client.GetPolicy(tc.ctx, &iam.GetPolicyInput{PolicyArn: aws.String(policyArn)})
		if err != nil {
			return fmt.Errorf("GetPolicy (after) failed: %w", err)
		}
		if aws.ToInt32(getAfter.Policy.PermissionsBoundaryUsageCount) != 0 {
			return fmt.Errorf("expected PermissionsBoundaryUsageCount=0 after delete, got %d", aws.ToInt32(getAfter.Policy.PermissionsBoundaryUsageCount))
		}
		return nil
	}))

	// CreateRole with Description exceeding Smithy roleDescriptionType
	// length(0-1000) must be rejected.
	results = append(results, r.RunTest("iam", "Error_CreateRoleDescriptionTooLong", func() error {
		longDesc := strings.Repeat("a", 1001)
		_, err := tc.client.CreateRole(tc.ctx, &iam.CreateRoleInput{
			RoleName:                 aws.String("EdgeRoleLongDesc-" + tc.ts),
			AssumeRolePolicyDocument: aws.String(assumeRolePolicy),
			Description:              aws.String(longDesc),
		})
		if err == nil {
			tc.client.DeleteRole(tc.ctx, &iam.DeleteRoleInput{RoleName: aws.String("EdgeRoleLongDesc-" + tc.ts)})
			return fmt.Errorf("expected error for Description length > 1000")
		}
		return nil
	}))

	// ListPolicies with an invalid Scope must be rejected
	results = append(results, r.RunTest("iam", "Error_ListPoliciesInvalidScope", func() error {
		_, err := tc.client.ListPolicies(tc.ctx, &iam.ListPoliciesInput{
			Scope: types.PolicyScopeType("Bogus"),
		})
		if err == nil {
			return fmt.Errorf("expected error for invalid Scope value")
		}
		return nil
	}))

	// CreateServiceLinkedRole with CustomSuffix exceeding Smithy
	// customSuffixType length(1-64) must be rejected.
	results = append(results, r.RunTest("iam", "Error_CreateServiceLinkedRoleCustomSuffixTooLong", func() error {
		longSuffix := strings.Repeat("a", 65)
		_, err := tc.client.CreateServiceLinkedRole(tc.ctx, &iam.CreateServiceLinkedRoleInput{
			AWSServiceName: aws.String("testservice.amazonaws.com"),
			CustomSuffix:   aws.String(longSuffix),
		})
		if err == nil {
			return fmt.Errorf("expected error for CustomSuffix length > 64")
		}
		return nil
	}))

	// CreateUser with PermissionsBoundary at creation time must set the
	// boundary in a single API call.
	results = append(results, r.RunTest("iam", "CreateUser_WithPermissionsBoundary", func() error {
		policyName := "EdgeCreateUserPB-" + tc.ts
		createOut, err := tc.client.CreatePolicy(tc.ctx, &iam.CreatePolicyInput{
			PolicyName:     aws.String(policyName),
			PolicyDocument: aws.String(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"s3:Get*","Resource":"*"}]}`),
		})
		if err != nil {
			return fmt.Errorf("CreatePolicy for PB test: %w", err)
		}
		pbArn := aws.ToString(createOut.Policy.Arn)

		userName := "edge-user-pb-create-" + tc.ts
		createUserResp, err := tc.client.CreateUser(tc.ctx, &iam.CreateUserInput{
			UserName:            aws.String(userName),
			PermissionsBoundary: aws.String(pbArn),
		})
		if err != nil {
			return fmt.Errorf("CreateUser with PermissionsBoundary failed: %w", err)
		}
		if createUserResp.User.PermissionsBoundary == nil {
			return fmt.Errorf("CreateUser response missing PermissionsBoundary")
		}
		if aws.ToString(createUserResp.User.PermissionsBoundary.PermissionsBoundaryArn) != pbArn {
			return fmt.Errorf("PermissionsBoundary ARN mismatch: got %s, want %s",
				aws.ToString(createUserResp.User.PermissionsBoundary.PermissionsBoundaryArn), pbArn)
		}

		// Verify PB usage count was incremented
		getPolicy, err := tc.client.GetPolicy(tc.ctx, &iam.GetPolicyInput{
			PolicyArn: aws.String(pbArn),
		})
		if err != nil {
			return fmt.Errorf("GetPolicy for PB count check: %w", err)
		}
		if aws.ToInt32(getPolicy.Policy.PermissionsBoundaryUsageCount) != 1 {
			return fmt.Errorf("PermissionsBoundaryUsageCount should be 1, got %d",
				aws.ToInt32(getPolicy.Policy.PermissionsBoundaryUsageCount))
		}

		// Clean up
		_, _ = tc.client.DeleteUser(tc.ctx, &iam.DeleteUserInput{UserName: aws.String(userName)})
		_, _ = tc.client.DeletePolicy(tc.ctx, &iam.DeletePolicyInput{PolicyArn: aws.String(pbArn)})
		return nil
	}))

	// CreateRole with PermissionsBoundary at creation time
	results = append(results, r.RunTest("iam", "CreateRole_WithPermissionsBoundary", func() error {
		policyName := "EdgeCreateRolePB-" + tc.ts
		createOut, err := tc.client.CreatePolicy(tc.ctx, &iam.CreatePolicyInput{
			PolicyName:     aws.String(policyName),
			PolicyDocument: aws.String(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"s3:Get*","Resource":"*"}]}`),
		})
		if err != nil {
			return fmt.Errorf("CreatePolicy for role PB test: %w", err)
		}
		pbArn := aws.ToString(createOut.Policy.Arn)

		roleName := "edge-role-pb-create-" + tc.ts
		createRoleResp, err := tc.client.CreateRole(tc.ctx, &iam.CreateRoleInput{
			RoleName:                 aws.String(roleName),
			AssumeRolePolicyDocument: aws.String(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"ec2.amazonaws.com"},"Action":"sts:AssumeRole"}]}`),
			PermissionsBoundary:      aws.String(pbArn),
		})
		if err != nil {
			return fmt.Errorf("CreateRole with PermissionsBoundary failed: %w", err)
		}
		if createRoleResp.Role.PermissionsBoundary == nil {
			return fmt.Errorf("CreateRole response missing PermissionsBoundary")
		}
		if aws.ToString(createRoleResp.Role.PermissionsBoundary.PermissionsBoundaryArn) != pbArn {
			return fmt.Errorf("PermissionsBoundary ARN mismatch: got %s, want %s",
				aws.ToString(createRoleResp.Role.PermissionsBoundary.PermissionsBoundaryArn), pbArn)
		}

		// Clean up
		_, _ = tc.client.DeleteRole(tc.ctx, &iam.DeleteRoleInput{RoleName: aws.String(roleName)})
		_, _ = tc.client.DeletePolicy(tc.ctx, &iam.DeletePolicyInput{PolicyArn: aws.String(pbArn)})
		return nil
	}))

	return results
}

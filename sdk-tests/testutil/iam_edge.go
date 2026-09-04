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

	// ListEntities_Pagination pins marker-based traversal across the
	// three core entity lists. Five entities of each kind are created
	// and then walked two per page; during a full regression other
	// services create resources concurrently, so a single page is never
	// guaranteed to hold everything.
	results = append(results, r.RunTest("iam", "ListEntities_Pagination", func() error {
		pgTs := tc.ts
		wildcardDoc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"*","Resource":"*"}]}`

		cases := []struct {
			label    string
			prefix   string
			create   func(name string) (key string, err error)
			remove   func(key string)
			listPage func(marker *string) ([]string, *string, error)
		}{
			{
				label:  "user",
				prefix: "PagUser",
				create: func(name string) (string, error) {
					if _, err := tc.client.CreateUser(tc.ctx, &iam.CreateUserInput{UserName: aws.String(name)}); err != nil {
						return "", err
					}
					return name, nil
				},
				remove: func(name string) {
					tc.client.DeleteUser(tc.ctx, &iam.DeleteUserInput{UserName: aws.String(name)})
				},
				listPage: func(marker *string) ([]string, *string, error) {
					resp, err := tc.client.ListUsers(tc.ctx, &iam.ListUsersInput{
						PathPrefix: aws.String("/"),
						Marker:     marker,
						MaxItems:   aws.Int32(2),
					})
					if err != nil {
						return nil, nil, err
					}
					names := make([]string, 0, len(resp.Users))
					for _, u := range resp.Users {
						names = append(names, aws.ToString(u.UserName))
					}
					return names, resp.Marker, nil
				},
			},
			{
				label:  "role",
				prefix: "PagRole",
				create: func(name string) (string, error) {
					if _, err := tc.client.CreateRole(tc.ctx, &iam.CreateRoleInput{
						RoleName:                 aws.String(name),
						AssumeRolePolicyDocument: aws.String(assumeRolePolicy),
					}); err != nil {
						return "", err
					}
					return name, nil
				},
				remove: func(name string) {
					tc.client.DeleteRole(tc.ctx, &iam.DeleteRoleInput{RoleName: aws.String(name)})
				},
				listPage: func(marker *string) ([]string, *string, error) {
					resp, err := tc.client.ListRoles(tc.ctx, &iam.ListRolesInput{
						PathPrefix: aws.String("/"),
						Marker:     marker,
						MaxItems:   aws.Int32(2),
					})
					if err != nil {
						return nil, nil, err
					}
					names := make([]string, 0, len(resp.Roles))
					for _, ro := range resp.Roles {
						names = append(names, aws.ToString(ro.RoleName))
					}
					return names, resp.Marker, nil
				},
			},
			{
				label:  "policy",
				prefix: "PagPolicy",
				create: func(name string) (string, error) {
					resp, err := tc.client.CreatePolicy(tc.ctx, &iam.CreatePolicyInput{
						PolicyName:     aws.String(name),
						PolicyDocument: aws.String(wildcardDoc),
					})
					if err != nil {
						return "", err
					}
					return aws.ToString(resp.Policy.Arn), nil
				},
				remove: func(arn string) {
					tc.client.DeletePolicy(tc.ctx, &iam.DeletePolicyInput{PolicyArn: aws.String(arn)})
				},
				listPage: func(marker *string) ([]string, *string, error) {
					resp, err := tc.client.ListPolicies(tc.ctx, &iam.ListPoliciesInput{
						Scope:    types.PolicyScopeTypeLocal,
						Marker:   marker,
						MaxItems: aws.Int32(2),
					})
					if err != nil {
						return nil, nil, err
					}
					names := make([]string, 0, len(resp.Policies))
					for _, p := range resp.Policies {
						names = append(names, aws.ToString(p.PolicyName))
					}
					return names, resp.Marker, nil
				},
			},
		}

		for _, pc := range cases {
			var created []string
			for i := 0; i < 5; i++ {
				name := fmt.Sprintf("%s-%s-%d", pc.prefix, pgTs, i)
				key, err := pc.create(name)
				if err != nil {
					for _, k := range created {
						pc.remove(k)
					}
					return fmt.Errorf("create %s %s: %v", pc.label, name, err)
				}
				created = append(created, key)
			}

			collected, err := iamPaginate(pc.listPage)

			for _, k := range created {
				pc.remove(k)
			}
			if err != nil {
				return fmt.Errorf("list %s pages: %v", pc.label, err)
			}

			count := 0
			for _, name := range collected {
				if strings.HasPrefix(name, pc.prefix+"-"+pgTs) {
					count++
				}
			}
			if count != 5 {
				return fmt.Errorf("expected 5 paginated %ss, got %d", pc.label, count)
			}
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
		cleanupRole, err := tc.createRole(roleName)
		if err != nil {
			return fmt.Errorf("setup CreateRole failed: %w", err)
		}
		defer cleanupRole()
		_, err = tc.client.DeleteServiceLinkedRole(tc.ctx, &iam.DeleteServiceLinkedRoleInput{
			RoleName: aws.String(roleName),
		})
		if err == nil {
			return fmt.Errorf("expected error when deleting non-service-linked role via DeleteServiceLinkedRole")
		}
		if !strings.Contains(err.Error(), "DeleteConflict") && !strings.Contains(err.Error(), "ConflictException") {
			return fmt.Errorf("expected DeleteConflict error, got: %v", err)
		}
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
		pbArn, cleanupPolicy, err := tc.createPolicy("EdgePBDriftUser-"+tc.ts, s3FullAccessPolicy)
		if err != nil {
			return fmt.Errorf("setup CreatePolicy failed: %w", err)
		}
		defer cleanupPolicy()

		userName := "EdgePBDriftUser-" + tc.ts
		cleanupUser, err := tc.createUser(userName)
		if err != nil {
			return fmt.Errorf("setup CreateUser failed: %w", err)
		}

		return iamAssertBoundaryUsageDrift(tc, pbArn,
			func() error {
				if _, err := tc.client.PutUserPermissionsBoundary(tc.ctx, &iam.PutUserPermissionsBoundaryInput{
					UserName:            aws.String(userName),
					PermissionsBoundary: aws.String(pbArn),
				}); err != nil {
					cleanupUser()
					return fmt.Errorf("PutUserPermissionsBoundary failed: %w", err)
				}
				return nil
			},
			func() error {
				if _, err := tc.client.DeleteUser(tc.ctx, &iam.DeleteUserInput{UserName: aws.String(userName)}); err != nil {
					return fmt.Errorf("DeleteUser failed: %w", err)
				}
				return nil
			},
		)
	}))

	// Deleting a role that has a permissions boundary attached must
	// decrement the referenced policy's PermissionsBoundaryUsageCount.
	results = append(results, r.RunTest("iam", "DeleteRole_DecrementsPermissionsBoundaryUsageCount", func() error {
		pbArn, cleanupPolicy, err := tc.createPolicy("EdgePBDriftRole-"+tc.ts, s3FullAccessPolicy)
		if err != nil {
			return fmt.Errorf("setup CreatePolicy failed: %w", err)
		}
		defer cleanupPolicy()

		roleName := "EdgePBDriftRole-" + tc.ts
		cleanupRole, err := tc.createRole(roleName)
		if err != nil {
			return fmt.Errorf("setup CreateRole failed: %w", err)
		}

		return iamAssertBoundaryUsageDrift(tc, pbArn,
			func() error {
				if _, err := tc.client.PutRolePermissionsBoundary(tc.ctx, &iam.PutRolePermissionsBoundaryInput{
					RoleName:            aws.String(roleName),
					PermissionsBoundary: aws.String(pbArn),
				}); err != nil {
					cleanupRole()
					return fmt.Errorf("PutRolePermissionsBoundary failed: %w", err)
				}
				return nil
			},
			func() error {
				if _, err := tc.client.DeleteRole(tc.ctx, &iam.DeleteRoleInput{RoleName: aws.String(roleName)}); err != nil {
					return fmt.Errorf("DeleteRole failed: %w", err)
				}
				return nil
			},
		)
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
		pbArn, cleanupPolicy, err := tc.createPolicy("EdgeCreateUserPB-"+tc.ts, iamAllowPolicy("s3:Get*"))
		if err != nil {
			return fmt.Errorf("CreatePolicy for PB test: %w", err)
		}
		defer cleanupPolicy()

		userName := "edge-user-pb-create-" + tc.ts
		createUserResp, err := tc.client.CreateUser(tc.ctx, &iam.CreateUserInput{
			UserName:            aws.String(userName),
			PermissionsBoundary: aws.String(pbArn),
		})
		if err != nil {
			return fmt.Errorf("CreateUser with PermissionsBoundary failed: %w", err)
		}
		defer tc.client.DeleteUser(tc.ctx, &iam.DeleteUserInput{UserName: aws.String(userName)})

		if createUserResp.User.PermissionsBoundary == nil {
			return fmt.Errorf("CreateUser response missing PermissionsBoundary")
		}
		if got := aws.ToString(createUserResp.User.PermissionsBoundary.PermissionsBoundaryArn); got != pbArn {
			return fmt.Errorf("PermissionsBoundary ARN mismatch: got %s, want %s", got, pbArn)
		}

		// Verify PB usage count was incremented; the deferred cleanups
		// remove the user before the policy it bounds.
		return iamAssertBoundaryUsageCount(tc, "after create-with-boundary", pbArn, 1)
	}))

	// CreateRole with PermissionsBoundary at creation time
	results = append(results, r.RunTest("iam", "CreateRole_WithPermissionsBoundary", func() error {
		pbArn, cleanupPolicy, err := tc.createPolicy("EdgeCreateRolePB-"+tc.ts, iamAllowPolicy("s3:Get*"))
		if err != nil {
			return fmt.Errorf("CreatePolicy for role PB test: %w", err)
		}
		defer cleanupPolicy()

		roleName := "edge-role-pb-create-" + tc.ts
		createRoleResp, err := tc.client.CreateRole(tc.ctx, &iam.CreateRoleInput{
			RoleName:                 aws.String(roleName),
			AssumeRolePolicyDocument: aws.String(ec2AssumeRolePolicy),
			PermissionsBoundary:      aws.String(pbArn),
		})
		if err != nil {
			return fmt.Errorf("CreateRole with PermissionsBoundary failed: %w", err)
		}
		defer tc.client.DeleteRole(tc.ctx, &iam.DeleteRoleInput{RoleName: aws.String(roleName)})

		if createRoleResp.Role.PermissionsBoundary == nil {
			return fmt.Errorf("CreateRole response missing PermissionsBoundary")
		}
		if got := aws.ToString(createRoleResp.Role.PermissionsBoundary.PermissionsBoundaryArn); got != pbArn {
			return fmt.Errorf("PermissionsBoundary ARN mismatch: got %s, want %s", got, pbArn)
		}

		return iamAssertBoundaryUsageCount(tc, "after create-with-boundary", pbArn, 1)
	}))

	return results
}

// iamAssertBoundaryUsageCount verifies that the policy's
// PermissionsBoundaryUsageCount equals want, where what locates the
// check in failure messages.
func iamAssertBoundaryUsageCount(tc *iamTestContext, what, pbArn string, want int32) error {
	getPolicy, err := tc.client.GetPolicy(tc.ctx, &iam.GetPolicyInput{PolicyArn: aws.String(pbArn)})
	if err != nil {
		return fmt.Errorf("GetPolicy (%s) failed: %w", what, err)
	}
	if got := aws.ToInt32(getPolicy.Policy.PermissionsBoundaryUsageCount); got != want {
		return fmt.Errorf("expected PermissionsBoundaryUsageCount=%d %s, got %d", want, what, got)
	}
	return nil
}

// iamAssertBoundaryUsageDrift attaches pbArn to an entity via attach,
// asserts the referenced policy's usage count is 1, deletes the
// boundary-bearing entity via removeEntity, and asserts the count
// returns to 0.
func iamAssertBoundaryUsageDrift(tc *iamTestContext, pbArn string, attach, removeEntity func() error) error {
	if err := attach(); err != nil {
		return err
	}
	if err := iamAssertBoundaryUsageCount(tc, "before delete", pbArn, 1); err != nil {
		return err
	}
	if err := removeEntity(); err != nil {
		return err
	}
	return iamAssertBoundaryUsageCount(tc, "after delete", pbArn, 0)
}

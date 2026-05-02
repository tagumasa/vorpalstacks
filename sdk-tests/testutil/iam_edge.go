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

	return results
}

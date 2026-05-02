package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

func (r *TestRunner) iamPermissionsBoundaryTests(tc *iamTestContext) []TestResult {
	var results []TestResult

	// User permissions boundary
	results = append(results, r.RunTest("iam", "PutUserPermissionsBoundary", func() error {
		_, err := tc.client.PutUserPermissionsBoundary(tc.ctx, &iam.PutUserPermissionsBoundaryInput{
			UserName:            aws.String(tc.user),
			PermissionsBoundary: aws.String(tc.policyArn),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "GetUser_PermissionsBoundary", func() error {
		resp, err := tc.client.GetUser(tc.ctx, &iam.GetUserInput{
			UserName: aws.String(tc.user),
		})
		if err != nil {
			return err
		}
		if resp.User.PermissionsBoundary == nil {
			return fmt.Errorf("permissions boundary is nil")
		}
		if aws.ToString(resp.User.PermissionsBoundary.PermissionsBoundaryArn) != tc.policyArn {
			return fmt.Errorf("permissions boundary arn mismatch: got %s, want %s",
				aws.ToString(resp.User.PermissionsBoundary.PermissionsBoundaryArn), tc.policyArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteUserPermissionsBoundary", func() error {
		_, err := tc.client.DeleteUserPermissionsBoundary(tc.ctx, &iam.DeleteUserPermissionsBoundaryInput{
			UserName: aws.String(tc.user),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetUser(tc.ctx, &iam.GetUserInput{
			UserName: aws.String(tc.user),
		})
		if err != nil {
			return err
		}
		if resp.User.PermissionsBoundary != nil {
			return fmt.Errorf("permissions boundary should be nil after deletion")
		}
		return nil
	}))

	// Role permissions boundary
	results = append(results, r.RunTest("iam", "PutRolePermissionsBoundary", func() error {
		_, err := tc.client.PutRolePermissionsBoundary(tc.ctx, &iam.PutRolePermissionsBoundaryInput{
			RoleName:            aws.String(tc.role),
			PermissionsBoundary: aws.String(tc.policyArn),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "GetRole_PermissionsBoundary", func() error {
		resp, err := tc.client.GetRole(tc.ctx, &iam.GetRoleInput{
			RoleName: aws.String(tc.role),
		})
		if err != nil {
			return err
		}
		if resp.Role.PermissionsBoundary == nil {
			return fmt.Errorf("permissions boundary is nil")
		}
		if aws.ToString(resp.Role.PermissionsBoundary.PermissionsBoundaryArn) != tc.policyArn {
			return fmt.Errorf("permissions boundary arn mismatch: got %s, want %s",
				aws.ToString(resp.Role.PermissionsBoundary.PermissionsBoundaryArn), tc.policyArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteRolePermissionsBoundary", func() error {
		_, err := tc.client.DeleteRolePermissionsBoundary(tc.ctx, &iam.DeleteRolePermissionsBoundaryInput{
			RoleName: aws.String(tc.role),
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
		if resp.Role.PermissionsBoundary != nil {
			return fmt.Errorf("permissions boundary should be nil after deletion")
		}
		return nil
	}))

	return results
}

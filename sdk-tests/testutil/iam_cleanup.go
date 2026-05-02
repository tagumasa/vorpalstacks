package testutil

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

func (r *TestRunner) iamCleanup(tc *iamTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("iam", "_Cleanup_DeleteLoginProfile", func() error {
		_, err := tc.client.DeleteLoginProfile(tc.ctx, &iam.DeleteLoginProfileInput{
			UserName: aws.String(tc.user),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "_Cleanup_DeleteAccessKey", func() error {
		_, err := tc.client.DeleteAccessKey(tc.ctx, &iam.DeleteAccessKeyInput{
			UserName:    aws.String(tc.user),
			AccessKeyId: aws.String(tc.accessKeyId),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "_Cleanup_DeleteInstanceProfile", func() error {
		_, err := tc.client.DeleteInstanceProfile(tc.ctx, &iam.DeleteInstanceProfileInput{
			InstanceProfileName: aws.String(tc.profile),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "_Cleanup_DeleteUser", func() error {
		_, err := tc.client.DeleteUser(tc.ctx, &iam.DeleteUserInput{
			UserName: aws.String(tc.user),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "_Cleanup_DeleteGroup", func() error {
		_, err := tc.client.DeleteGroup(tc.ctx, &iam.DeleteGroupInput{
			GroupName: aws.String(tc.group),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "_Cleanup_DeleteRole", func() error {
		_, err := tc.client.DeleteRole(tc.ctx, &iam.DeleteRoleInput{
			RoleName: aws.String(tc.role),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "_Cleanup_DeletePolicy", func() error {
		_, err := tc.client.DeletePolicy(tc.ctx, &iam.DeletePolicyInput{
			PolicyArn: aws.String(tc.policyArn),
		})
		return err
	}))

	return results
}

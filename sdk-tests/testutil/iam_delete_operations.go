package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

// iamDeleteOperationTests verifies the Delete* operations beyond the
// silent _Cleanup_ path: each test creates a resource, deletes it through
// the API, and confirms the resource is gone.
func (r *TestRunner) iamDeleteOperationTests(tc *iamTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("iam", "DeleteUser_Verified", func() error {
		user := fmt.Sprintf("DelUser-%s", tc.ts)
		// The deletion below is the operation under test, so the helper's
		// cleanup is discarded.
		_, err := tc.createUser(user)
		if err != nil {
			return err
		}
		if _, err := tc.client.DeleteUser(tc.ctx, &iam.DeleteUserInput{UserName: aws.String(user)}); err != nil {
			return err
		}
		_, err = tc.client.GetUser(tc.ctx, &iam.GetUserInput{UserName: aws.String(user)})
		return iamAssertNoSuchEntity(err, "GetUser after delete")
	}))

	results = append(results, r.RunTest("iam", "DeleteGroup_Verified", func() error {
		group := fmt.Sprintf("DelGroup-%s", tc.ts)
		// The deletion below is the operation under test, so the helper's
		// cleanup is discarded.
		_, err := tc.createGroup(group)
		if err != nil {
			return err
		}
		if _, err := tc.client.DeleteGroup(tc.ctx, &iam.DeleteGroupInput{GroupName: aws.String(group)}); err != nil {
			return err
		}
		_, err = tc.client.GetGroup(tc.ctx, &iam.GetGroupInput{GroupName: aws.String(group)})
		return iamAssertNoSuchEntity(err, "GetGroup after delete")
	}))

	results = append(results, r.RunTest("iam", "DeleteRole_Verified", func() error {
		role := fmt.Sprintf("DelRole-%s", tc.ts)
		cleanupRole, err := tc.createRole(role)
		if err != nil {
			return err
		}
		cleanupRole()
		_, err = tc.client.GetRole(tc.ctx, &iam.GetRoleInput{RoleName: aws.String(role)})
		return iamAssertNoSuchEntity(err, "GetRole after delete")
	}))

	results = append(results, r.RunTest("iam", "DeletePolicy_Verified", func() error {
		policy := fmt.Sprintf("DelPolicy-%s", tc.ts)
		// The deletion below is the operation under test, so the helper's
		// cleanup is discarded.
		arn, _, err := tc.createPolicy(policy, iamAllowPolicy("s3:GetObject"))
		if err != nil {
			return err
		}
		if _, err := tc.client.DeletePolicy(tc.ctx, &iam.DeletePolicyInput{PolicyArn: aws.String(arn)}); err != nil {
			return err
		}
		_, err = tc.client.GetPolicy(tc.ctx, &iam.GetPolicyInput{PolicyArn: aws.String(arn)})
		return iamAssertNoSuchEntity(err, "GetPolicy after delete")
	}))

	results = append(results, r.RunTest("iam", "DeleteAccessKey_Verified", func() error {
		user := fmt.Sprintf("DelKey-%s", tc.ts)
		cleanupUser, err := tc.createUser(user)
		if err != nil {
			return err
		}
		defer cleanupUser()

		key, err := tc.client.CreateAccessKey(tc.ctx, &iam.CreateAccessKeyInput{UserName: aws.String(user)})
		if err != nil {
			return err
		}
		if _, err := tc.client.DeleteAccessKey(tc.ctx, &iam.DeleteAccessKeyInput{
			UserName:    aws.String(user),
			AccessKeyId: key.AccessKey.AccessKeyId,
		}); err != nil {
			return err
		}
		list, err := tc.client.ListAccessKeys(tc.ctx, &iam.ListAccessKeysInput{UserName: aws.String(user)})
		if err != nil {
			return err
		}
		for _, k := range list.AccessKeyMetadata {
			if aws.ToString(k.AccessKeyId) == aws.ToString(key.AccessKey.AccessKeyId) {
				return fmt.Errorf("access key still listed after delete")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteInstanceProfile_Verified", func() error {
		profile := fmt.Sprintf("DelProf-%s", tc.ts)
		// The deletion below is the operation under test, so the helper's
		// cleanup is discarded.
		_, err := tc.createInstanceProfile(profile)
		if err != nil {
			return err
		}
		if _, err := tc.client.DeleteInstanceProfile(tc.ctx, &iam.DeleteInstanceProfileInput{
			InstanceProfileName: aws.String(profile),
		}); err != nil {
			return err
		}
		_, err = tc.client.GetInstanceProfile(tc.ctx, &iam.GetInstanceProfileInput{
			InstanceProfileName: aws.String(profile),
		})
		return iamAssertNoSuchEntity(err, "GetInstanceProfile after delete")
	}))

	results = append(results, r.RunTest("iam", "DeleteLoginProfile_Verified", func() error {
		user := fmt.Sprintf("DelLogin-%s", tc.ts)
		cleanupUser, err := tc.createUser(user)
		if err != nil {
			return err
		}
		defer cleanupUser()

		if _, err := tc.client.CreateLoginProfile(tc.ctx, &iam.CreateLoginProfileInput{
			UserName: aws.String(user),
			Password: aws.String("Valid!Pw1-Del-2026"),
		}); err != nil {
			return err
		}
		if _, err := tc.client.DeleteLoginProfile(tc.ctx, &iam.DeleteLoginProfileInput{
			UserName: aws.String(user),
		}); err != nil {
			return err
		}
		_, err = tc.client.GetLoginProfile(tc.ctx, &iam.GetLoginProfileInput{UserName: aws.String(user)})
		return iamAssertNoSuchEntity(err, "GetLoginProfile after delete")
	}))

	return results
}

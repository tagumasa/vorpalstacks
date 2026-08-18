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
		if _, err := tc.client.CreateUser(tc.ctx, &iam.CreateUserInput{UserName: aws.String(user)}); err != nil {
			return err
		}
		if _, err := tc.client.DeleteUser(tc.ctx, &iam.DeleteUserInput{UserName: aws.String(user)}); err != nil {
			return err
		}
		_, err := tc.client.GetUser(tc.ctx, &iam.GetUserInput{UserName: aws.String(user)})
		if err == nil || !containsErrorCode(err, "NoSuchEntity") {
			return fmt.Errorf("GetUser after delete: got %v, want NoSuchEntity", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteGroup_Verified", func() error {
		group := fmt.Sprintf("DelGroup-%s", tc.ts)
		if _, err := tc.client.CreateGroup(tc.ctx, &iam.CreateGroupInput{GroupName: aws.String(group)}); err != nil {
			return err
		}
		if _, err := tc.client.DeleteGroup(tc.ctx, &iam.DeleteGroupInput{GroupName: aws.String(group)}); err != nil {
			return err
		}
		_, err := tc.client.GetGroup(tc.ctx, &iam.GetGroupInput{GroupName: aws.String(group)})
		if err == nil || !containsErrorCode(err, "NoSuchEntity") {
			return fmt.Errorf("GetGroup after delete: got %v, want NoSuchEntity", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteRole_Verified", func() error {
		role := fmt.Sprintf("DelRole-%s", tc.ts)
		if _, err := tc.client.CreateRole(tc.ctx, &iam.CreateRoleInput{
			RoleName:                 aws.String(role),
			AssumeRolePolicyDocument: aws.String(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"ec2.amazonaws.com"},"Action":"sts:AssumeRole"}]}`),
		}); err != nil {
			return err
		}
		if _, err := tc.client.DeleteRole(tc.ctx, &iam.DeleteRoleInput{RoleName: aws.String(role)}); err != nil {
			return err
		}
		_, err := tc.client.GetRole(tc.ctx, &iam.GetRoleInput{RoleName: aws.String(role)})
		if err == nil || !containsErrorCode(err, "NoSuchEntity") {
			return fmt.Errorf("GetRole after delete: got %v, want NoSuchEntity", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeletePolicy_Verified", func() error {
		policy := fmt.Sprintf("DelPolicy-%s", tc.ts)
		created, err := tc.client.CreatePolicy(tc.ctx, &iam.CreatePolicyInput{
			PolicyName:     aws.String(policy),
			PolicyDocument: aws.String(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"s3:GetObject","Resource":"*"}]}`),
		})
		if err != nil {
			return err
		}
		if _, err := tc.client.DeletePolicy(tc.ctx, &iam.DeletePolicyInput{PolicyArn: created.Policy.Arn}); err != nil {
			return err
		}
		_, err = tc.client.GetPolicy(tc.ctx, &iam.GetPolicyInput{PolicyArn: created.Policy.Arn})
		if err == nil || !containsErrorCode(err, "NoSuchEntity") {
			return fmt.Errorf("GetPolicy after delete: got %v, want NoSuchEntity", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteAccessKey_Verified", func() error {
		user := fmt.Sprintf("DelKey-%s", tc.ts)
		if _, err := tc.client.CreateUser(tc.ctx, &iam.CreateUserInput{UserName: aws.String(user)}); err != nil {
			return err
		}
		defer tc.client.DeleteUser(tc.ctx, &iam.DeleteUserInput{UserName: aws.String(user)})

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
		if _, err := tc.client.CreateInstanceProfile(tc.ctx, &iam.CreateInstanceProfileInput{
			InstanceProfileName: aws.String(profile),
		}); err != nil {
			return err
		}
		if _, err := tc.client.DeleteInstanceProfile(tc.ctx, &iam.DeleteInstanceProfileInput{
			InstanceProfileName: aws.String(profile),
		}); err != nil {
			return err
		}
		_, err := tc.client.GetInstanceProfile(tc.ctx, &iam.GetInstanceProfileInput{
			InstanceProfileName: aws.String(profile),
		})
		if err == nil || !containsErrorCode(err, "NoSuchEntity") {
			return fmt.Errorf("GetInstanceProfile after delete: got %v, want NoSuchEntity", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteLoginProfile_Verified", func() error {
		user := fmt.Sprintf("DelLogin-%s", tc.ts)
		if _, err := tc.client.CreateUser(tc.ctx, &iam.CreateUserInput{UserName: aws.String(user)}); err != nil {
			return err
		}
		defer tc.client.DeleteUser(tc.ctx, &iam.DeleteUserInput{UserName: aws.String(user)})

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
		_, err := tc.client.GetLoginProfile(tc.ctx, &iam.GetLoginProfileInput{UserName: aws.String(user)})
		if err == nil || !containsErrorCode(err, "NoSuchEntity") {
			return fmt.Errorf("GetLoginProfile after delete: got %v, want NoSuchEntity", err)
		}
		return nil
	}))

	return results
}

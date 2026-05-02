package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func (r *TestRunner) iamUserTests(tc *iamTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("iam", "CreateUser", func() error {
		resp, err := tc.client.CreateUser(tc.ctx, &iam.CreateUserInput{
			UserName: aws.String(tc.user),
			Tags: []types.Tag{
				{Key: aws.String("CreatedBy"), Value: aws.String("sdk-test")},
			},
		})
		if err != nil {
			return err
		}
		if resp.User == nil {
			return fmt.Errorf("user is nil")
		}
		if aws.ToString(resp.User.UserName) != tc.user {
			return fmt.Errorf("username mismatch: got %s, want %s", aws.ToString(resp.User.UserName), tc.user)
		}
		if aws.ToString(resp.User.Arn) == "" {
			return fmt.Errorf("arn is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetUser", func() error {
		resp, err := tc.client.GetUser(tc.ctx, &iam.GetUserInput{
			UserName: aws.String(tc.user),
		})
		if err != nil {
			return err
		}
		if resp.User == nil {
			return fmt.Errorf("user is nil")
		}
		if aws.ToString(resp.User.UserName) != tc.user {
			return fmt.Errorf("username mismatch: got %s, want %s", aws.ToString(resp.User.UserName), tc.user)
		}
		if aws.ToString(resp.User.UserId) == "" {
			return fmt.Errorf("user id is empty")
		}
		if aws.ToString(resp.User.Arn) == "" {
			return fmt.Errorf("arn is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListUsers", func() error {
		var allUsers []types.User
		var marker *string
		for {
			resp, err := tc.client.ListUsers(tc.ctx, &iam.ListUsersInput{
				Marker: marker,
			})
			if err != nil {
				return err
			}
			for _, u := range resp.Users {
				if aws.ToString(u.UserName) == tc.user {
					allUsers = append(allUsers, u)
				}
			}
			if resp.IsTruncated && resp.Marker != nil {
				marker = resp.Marker
			} else {
				break
			}
		}
		if len(allUsers) == 0 {
			return fmt.Errorf("user %s not found in ListUsers", tc.user)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UpdateUser", func() error {
		newName := fmt.Sprintf("UpdatedUser-%s", tc.ts)
		_, err := tc.client.UpdateUser(tc.ctx, &iam.UpdateUserInput{
			UserName:    aws.String(tc.user),
			NewUserName: aws.String(newName),
		})
		if err != nil {
			return err
		}
		tc.user = newName
		return nil
	}))

	// Access keys
	results = append(results, r.RunTest("iam", "CreateAccessKey", func() error {
		resp, err := tc.client.CreateAccessKey(tc.ctx, &iam.CreateAccessKeyInput{
			UserName: aws.String(tc.user),
		})
		if err != nil {
			return err
		}
		if resp.AccessKey == nil {
			return fmt.Errorf("access key is nil")
		}
		if aws.ToString(resp.AccessKey.AccessKeyId) == "" {
			return fmt.Errorf("access key id is empty")
		}
		if aws.ToString(resp.AccessKey.SecretAccessKey) == "" {
			return fmt.Errorf("secret access key is empty")
		}
		if resp.AccessKey.Status != types.StatusTypeActive {
			return fmt.Errorf("expected Active status, got %s", resp.AccessKey.Status)
		}
		if aws.ToString(resp.AccessKey.UserName) != tc.user {
			return fmt.Errorf("username mismatch on access key")
		}
		tc.accessKeyId = *resp.AccessKey.AccessKeyId
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListAccessKeys", func() error {
		resp, err := tc.client.ListAccessKeys(tc.ctx, &iam.ListAccessKeysInput{
			UserName: aws.String(tc.user),
		})
		if err != nil {
			return err
		}
		found := false
		for _, k := range resp.AccessKeyMetadata {
			if aws.ToString(k.AccessKeyId) == tc.accessKeyId {
				found = true
				if aws.ToString(k.UserName) != tc.user {
					return fmt.Errorf("username mismatch in key metadata")
				}
				if k.Status != types.StatusTypeActive {
					return fmt.Errorf("expected Active status, got %s", k.Status)
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("access key %s not found", tc.accessKeyId)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UpdateAccessKey", func() error {
		_, err := tc.client.UpdateAccessKey(tc.ctx, &iam.UpdateAccessKeyInput{
			AccessKeyId: aws.String(tc.accessKeyId),
			Status:      types.StatusTypeInactive,
			UserName:    aws.String(tc.user),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListAccessKeys(tc.ctx, &iam.ListAccessKeysInput{
			UserName: aws.String(tc.user),
		})
		if err != nil {
			return err
		}
		for _, k := range resp.AccessKeyMetadata {
			if aws.ToString(k.AccessKeyId) == tc.accessKeyId {
				if k.Status != types.StatusTypeInactive {
					return fmt.Errorf("expected Inactive status after update, got %s", k.Status)
				}
				return nil
			}
		}
		return fmt.Errorf("access key not found after update")
	}))

	results = append(results, r.RunTest("iam", "GetAccessKeyLastUsed", func() error {
		resp, err := tc.client.GetAccessKeyLastUsed(tc.ctx, &iam.GetAccessKeyLastUsedInput{
			AccessKeyId: aws.String(tc.accessKeyId),
		})
		if err != nil {
			return err
		}
		if aws.ToString(resp.UserName) != tc.user {
			return fmt.Errorf("username mismatch: got %s, want %s", aws.ToString(resp.UserName), tc.user)
		}
		return nil
	}))

	// Login profile
	results = append(results, r.RunTest("iam", "CreateLoginProfile", func() error {
		resp, err := tc.client.CreateLoginProfile(tc.ctx, &iam.CreateLoginProfileInput{
			UserName:              aws.String(tc.user),
			Password:              aws.String("TempPassword123!"),
			PasswordResetRequired: true,
		})
		if err != nil {
			return err
		}
		if resp.LoginProfile == nil {
			return fmt.Errorf("login profile is nil")
		}
		if aws.ToString(resp.LoginProfile.UserName) != tc.user {
			return fmt.Errorf("login profile username mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetLoginProfile", func() error {
		resp, err := tc.client.GetLoginProfile(tc.ctx, &iam.GetLoginProfileInput{
			UserName: aws.String(tc.user),
		})
		if err != nil {
			return err
		}
		if resp.LoginProfile == nil {
			return fmt.Errorf("login profile is nil")
		}
		if aws.ToString(resp.LoginProfile.UserName) != tc.user {
			return fmt.Errorf("login profile username mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UpdateLoginProfile", func() error {
		_, err := tc.client.UpdateLoginProfile(tc.ctx, &iam.UpdateLoginProfileInput{
			UserName:              aws.String(tc.user),
			Password:              aws.String("NewPassword456!"),
			PasswordResetRequired: aws.Bool(true),
		})
		return err
	}))

	// User tags
	results = append(results, r.RunTest("iam", "TagUser", func() error {
		_, err := tc.client.TagUser(tc.ctx, &iam.TagUserInput{
			UserName: aws.String(tc.user),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "ListUserTags", func() error {
		resp, err := tc.client.ListUserTags(tc.ctx, &iam.ListUserTagsInput{
			UserName: aws.String(tc.user),
		})
		if err != nil {
			return err
		}
		if !iamTagPresent(resp.Tags, "Environment", "test") {
			return fmt.Errorf("Environment=test tag not found")
		}
		if !iamTagPresent(resp.Tags, "CreatedBy", "sdk-test") {
			return fmt.Errorf("CreatedBy=sdk-test tag not found (from CreateUser)")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UntagUser", func() error {
		_, err := tc.client.UntagUser(tc.ctx, &iam.UntagUserInput{
			UserName: aws.String(tc.user),
			TagKeys:  []string{"Environment"},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListUserTags(tc.ctx, &iam.ListUserTagsInput{
			UserName: aws.String(tc.user),
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

	// User inline policies
	results = append(results, r.RunTest("iam", "PutUserPolicy", func() error {
		_, err := tc.client.PutUserPolicy(tc.ctx, &iam.PutUserPolicyInput{
			UserName:       aws.String(tc.user),
			PolicyName:     aws.String(tc.userInlinePolicy),
			PolicyDocument: aws.String(s3FullAccessPolicy),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "GetUserPolicy", func() error {
		resp, err := tc.client.GetUserPolicy(tc.ctx, &iam.GetUserPolicyInput{
			UserName:   aws.String(tc.user),
			PolicyName: aws.String(tc.userInlinePolicy),
		})
		if err != nil {
			return err
		}
		if resp.PolicyDocument == nil || *resp.PolicyDocument == "" {
			return fmt.Errorf("policy document is empty")
		}
		if aws.ToString(resp.UserName) != tc.user {
			return fmt.Errorf("username mismatch in GetUserPolicy")
		}
		if aws.ToString(resp.PolicyName) != tc.userInlinePolicy {
			return fmt.Errorf("policy name mismatch in GetUserPolicy")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListUserPolicies", func() error {
		resp, err := tc.client.ListUserPolicies(tc.ctx, &iam.ListUserPoliciesInput{
			UserName: aws.String(tc.user),
		})
		if err != nil {
			return err
		}
		found := false
		for _, name := range resp.PolicyNames {
			if name == tc.userInlinePolicy {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("inline policy %s not found in ListUserPolicies", tc.userInlinePolicy)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteUserPolicy", func() error {
		_, err := tc.client.DeleteUserPolicy(tc.ctx, &iam.DeleteUserPolicyInput{
			UserName:   aws.String(tc.user),
			PolicyName: aws.String(tc.userInlinePolicy),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListUserPolicies(tc.ctx, &iam.ListUserPoliciesInput{
			UserName: aws.String(tc.user),
		})
		if err != nil {
			return err
		}
		for _, name := range resp.PolicyNames {
			if name == tc.userInlinePolicy {
				return fmt.Errorf("inline policy should be deleted")
			}
		}
		return nil
	}))

	return results
}

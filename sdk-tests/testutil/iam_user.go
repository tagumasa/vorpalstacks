package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"vorpalstacks-sdk-tests/config"
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
		users, err := iamPaginate(func(marker *string) ([]types.User, *string, error) {
			resp, err := tc.client.ListUsers(tc.ctx, &iam.ListUsersInput{Marker: marker})
			if err != nil {
				return nil, nil, err
			}
			return resp.Users, resp.Marker, nil
		})
		if err != nil {
			return err
		}
		for _, u := range users {
			if aws.ToString(u.UserName) == tc.user {
				return nil
			}
		}
		return fmt.Errorf("user %s not found in ListUsers", tc.user)
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
		resp, err := tc.client.GetUser(tc.ctx, &iam.GetUserInput{
			UserName: aws.String(newName),
		})
		if err != nil {
			return fmt.Errorf("GetUser with new name after UpdateUser: %w", err)
		}
		if aws.ToString(resp.User.UserName) != newName {
			return fmt.Errorf("username not updated: got %s, want %s", aws.ToString(resp.User.UserName), newName)
		}

		// A new name outside the entity-name pattern must be rejected, and
		// renaming onto an existing user must fail with EntityAlreadyExists.
		if _, err := tc.client.UpdateUser(tc.ctx, &iam.UpdateUserInput{
			UserName:    aws.String(newName),
			NewUserName: aws.String("invalid user name!"),
		}); err == nil || !isInvalidInputError(err) {
			return fmt.Errorf("invalid NewUserName: got %v, want InvalidInput", err)
		}
		other := fmt.Sprintf("UpdateOther-%s", tc.ts)
		cleanupOther, err := tc.createUser(other)
		if err != nil {
			return err
		}
		defer cleanupOther()
		if _, err := tc.client.UpdateUser(tc.ctx, &iam.UpdateUserInput{
			UserName:    aws.String(newName),
			NewUserName: aws.String(other),
		}); err == nil || !containsErrorCode(err, "EntityAlreadyExists") {
			return fmt.Errorf("rename onto existing user: got %v, want EntityAlreadyExists", err)
		}
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
		if err != nil {
			return err
		}
		resp, err := tc.client.GetLoginProfile(tc.ctx, &iam.GetLoginProfileInput{
			UserName: aws.String(tc.user),
		})
		if err != nil {
			return fmt.Errorf("GetLoginProfile after update: %w", err)
		}
		if resp.LoginProfile == nil {
			return fmt.Errorf("login profile is nil after UpdateLoginProfile")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ChangePassword", func() error {
		user := fmt.Sprintf("ChangePw-%s", tc.ts)
		cleanupUser, err := tc.createUser(user)
		if err != nil {
			return fmt.Errorf("CreateUser for ChangePassword: %w", err)
		}
		defer cleanupUser()

		const oldPass = "Valid!Old1-Pw-2026"
		const newPass = "Valid!New1-Pw-2026"
		if _, err := tc.client.CreateLoginProfile(tc.ctx, &iam.CreateLoginProfileInput{
			UserName:              aws.String(user),
			Password:              aws.String(oldPass),
			PasswordResetRequired: false,
		}); err != nil {
			return fmt.Errorf("CreateLoginProfile for ChangePassword: %w", err)
		}
		defer tc.client.DeleteLoginProfile(tc.ctx, &iam.DeleteLoginProfileInput{UserName: aws.String(user)})

		key, err := tc.client.CreateAccessKey(tc.ctx, &iam.CreateAccessKeyInput{UserName: aws.String(user)})
		if err != nil {
			return fmt.Errorf("CreateAccessKey for ChangePassword: %w", err)
		}
		defer tc.client.DeleteAccessKey(tc.ctx, &iam.DeleteAccessKeyInput{
			UserName:    aws.String(user),
			AccessKeyId: key.AccessKey.AccessKeyId,
		})

		// ChangePassword carries only OldPassword/NewPassword on the wire;
		// the target user is the authenticated caller itself.
		cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
			Endpoint: r.endpoint,
			Region:   r.region,
		})
		if err != nil {
			return fmt.Errorf("load config for user client: %w", err)
		}
		cfg.Credentials = credentials.NewStaticCredentialsProvider(
			*key.AccessKey.AccessKeyId,
			*key.AccessKey.SecretAccessKey,
			"",
		)
		userClient := iam.NewFromConfig(cfg)

		if _, err := userClient.ChangePassword(tc.ctx, &iam.ChangePasswordInput{
			OldPassword: aws.String(oldPass),
			NewPassword: aws.String(newPass),
		}); err != nil {
			return fmt.Errorf("ChangePassword as the owning user failed: %w", err)
		}

		// The superseded password must not be accepted as OldPassword.
		if _, err := userClient.ChangePassword(tc.ctx, &iam.ChangePasswordInput{
			OldPassword: aws.String(oldPass),
			NewPassword: aws.String("Valid!New2-Pw-2026"),
		}); err == nil {
			return fmt.Errorf("ChangePassword accepted the superseded old password")
		}
		return nil
	}))

	// User tags
	results = append(results, r.RunTest("iam", "TagUser", func() error {
		_, err := tc.client.TagUser(tc.ctx, &iam.TagUserInput{
			UserName: aws.String(tc.user),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
			},
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListUserTags(tc.ctx, &iam.ListUserTagsInput{
			UserName: aws.String(tc.user),
		})
		if err != nil {
			return fmt.Errorf("ListUserTags after tag: %w", err)
		}
		if !iamTagPresent(resp.Tags, "Environment", "test") {
			return fmt.Errorf("Environment=test tag not found after TagUser")
		}
		return nil
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

	// Tag keys and values are limited to 128/256 Unicode characters, not
	// bytes, per the Smithy tagKeyType/tagValueType @length traits.
	results = append(results, r.RunTest("iam", "TagUser_MultibyteTagAccepted", func() error {
		key := strings.Repeat("\u65e5", 100)
		value := strings.Repeat("\u672c", 200)
		_, err := tc.client.TagUser(tc.ctx, &iam.TagUserInput{
			UserName: aws.String(tc.user),
			Tags:     []types.Tag{{Key: aws.String(key), Value: aws.String(value)}},
		})
		if err != nil {
			return fmt.Errorf("multibyte tag within the Unicode length limits rejected: %w", err)
		}
		resp, err := tc.client.ListUserTags(tc.ctx, &iam.ListUserTagsInput{
			UserName: aws.String(tc.user),
		})
		if err != nil {
			return err
		}
		if !iamTagPresent(resp.Tags, key, value) {
			return fmt.Errorf("multibyte tag not found after TagUser")
		}
		if _, err := tc.client.UntagUser(tc.ctx, &iam.UntagUserInput{
			UserName: aws.String(tc.user),
			TagKeys:  []string{key},
		}); err != nil {
			return fmt.Errorf("cleanup UntagUser: %w", err)
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
		if err != nil {
			return err
		}
		resp, err := tc.client.GetUserPolicy(tc.ctx, &iam.GetUserPolicyInput{
			UserName:   aws.String(tc.user),
			PolicyName: aws.String(tc.userInlinePolicy),
		})
		if err != nil {
			return fmt.Errorf("GetUserPolicy after PutUserPolicy: %w", err)
		}
		if resp.PolicyDocument == nil || *resp.PolicyDocument == "" {
			return fmt.Errorf("policy document is empty after PutUserPolicy")
		}
		return nil
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

	results = append(results, r.RunTest("iam", "CreateLoginProfile_PasswordTooLong", func() error {
		user := fmt.Sprintf("LongPw-%s", tc.ts)
		cleanupUser, err := tc.createUser(user)
		if err != nil {
			return err
		}
		defer cleanupUser()

		_, err = tc.client.CreateLoginProfile(tc.ctx, &iam.CreateLoginProfileInput{
			UserName: aws.String(user),
			Password: aws.String(strings.Repeat("A1!a", 33)), // 132 characters
		})
		if err == nil {
			return fmt.Errorf("a password longer than 128 characters must be rejected")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListUserTags_Pagination", func() error {
		// The user carries CreatedBy from CreateUser; add two more tags so
		// pagination has three entries to traverse.
		if _, err := tc.client.TagUser(tc.ctx, &iam.TagUserInput{
			UserName: aws.String(tc.user),
			Tags: []types.Tag{
				{Key: aws.String("Page1"), Value: aws.String("a")},
				{Key: aws.String("Page2"), Value: aws.String("b")},
			},
		}); err != nil {
			return err
		}
		defer tc.client.UntagUser(tc.ctx, &iam.UntagUserInput{
			UserName: aws.String(tc.user),
			TagKeys:  []string{"Page1", "Page2"},
		})

		first, err := tc.client.ListUserTags(tc.ctx, &iam.ListUserTagsInput{
			UserName: aws.String(tc.user),
			MaxItems: aws.Int32(2),
		})
		if err != nil {
			return err
		}
		if len(first.Tags) != 2 {
			return fmt.Errorf("first page size: got %d, want 2", len(first.Tags))
		}
		if !first.IsTruncated || first.Marker == nil || *first.Marker == "" {
			return fmt.Errorf("first page must be truncated with a marker")
		}

		second, err := tc.client.ListUserTags(tc.ctx, &iam.ListUserTagsInput{
			UserName: aws.String(tc.user),
			Marker:   first.Marker,
			MaxItems: aws.Int32(2),
		})
		if err != nil {
			return err
		}
		if len(second.Tags) != 1 {
			return fmt.Errorf("second page size: got %d, want 1", len(second.Tags))
		}
		if second.IsTruncated {
			return fmt.Errorf("second page must not be truncated")
		}
		return nil
	}))

	return results
}

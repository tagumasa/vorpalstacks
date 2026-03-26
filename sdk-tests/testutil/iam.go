package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunIAMTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "iam",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := iam.NewFromConfig(cfg)
	ctx := context.Background()

	userName := fmt.Sprintf("TestUser-%d", time.Now().UnixNano())
	groupName := fmt.Sprintf("TestGroup-%d", time.Now().UnixNano())
	roleName := fmt.Sprintf("TestRole-%d", time.Now().UnixNano())
	policyName := fmt.Sprintf("TestPolicy-%d", time.Now().UnixNano())
	userInlinePolicyName := fmt.Sprintf("UserPolicy-%d", time.Now().UnixNano())
	roleInlinePolicyName := fmt.Sprintf("RolePolicy-%d", time.Now().UnixNano())
	var accessKeyId string

	results = append(results, r.RunTest("iam", "CreateUser", func() error {
		_, err := client.CreateUser(ctx, &iam.CreateUserInput{
			UserName: aws.String(userName),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "GetUser", func() error {
		resp, err := client.GetUser(ctx, &iam.GetUserInput{
			UserName: aws.String(userName),
		})
		if err != nil {
			return err
		}
		if resp.User == nil {
			return fmt.Errorf("user is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListUsers", func() error {
		resp, err := client.ListUsers(ctx, &iam.ListUsersInput{})
		if err != nil {
			return err
		}
		if resp.Users == nil {
			return fmt.Errorf("users list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "CreateAccessKey", func() error {
		resp, err := client.CreateAccessKey(ctx, &iam.CreateAccessKeyInput{
			UserName: aws.String(userName),
		})
		if err != nil {
			return err
		}
		accessKeyId = *resp.AccessKey.AccessKeyId
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListAccessKeys", func() error {
		resp, err := client.ListAccessKeys(ctx, &iam.ListAccessKeysInput{
			UserName: aws.String(userName),
		})
		if err != nil {
			return err
		}
		if resp.AccessKeyMetadata == nil {
			return fmt.Errorf("access keys list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "CreateLoginProfile", func() error {
		_, err := client.CreateLoginProfile(ctx, &iam.CreateLoginProfileInput{
			UserName: aws.String(userName),
			Password: aws.String("TempPassword123!"),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "GetLoginProfile", func() error {
		resp, err := client.GetLoginProfile(ctx, &iam.GetLoginProfileInput{
			UserName: aws.String(userName),
		})
		if err != nil {
			return err
		}
		if resp.LoginProfile == nil {
			return fmt.Errorf("login profile is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UpdateUser", func() error {
		newUserName := fmt.Sprintf("UpdatedUser-%d", time.Now().UnixNano())
		_, err := client.UpdateUser(ctx, &iam.UpdateUserInput{
			UserName:    aws.String(userName),
			NewUserName: aws.String(newUserName),
		})
		if err == nil {
			userName = newUserName
		}
		return err
	}))

	results = append(results, r.RunTest("iam", "TagUser", func() error {
		_, err := client.TagUser(ctx, &iam.TagUserInput{
			UserName: aws.String(userName),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "ListUserTags", func() error {
		resp, err := client.ListUserTags(ctx, &iam.ListUserTagsInput{
			UserName: aws.String(userName),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil {
			return fmt.Errorf("tags is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UntagUser", func() error {
		_, err := client.UntagUser(ctx, &iam.UntagUserInput{
			UserName: aws.String(userName),
			TagKeys:  []string{"Environment"},
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "CreateGroup", func() error {
		_, err := client.CreateGroup(ctx, &iam.CreateGroupInput{
			GroupName: aws.String(groupName),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "GetGroup", func() error {
		resp, err := client.GetGroup(ctx, &iam.GetGroupInput{
			GroupName: aws.String(groupName),
		})
		if err != nil {
			return err
		}
		if resp.Group == nil {
			return fmt.Errorf("group is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListGroups", func() error {
		resp, err := client.ListGroups(ctx, &iam.ListGroupsInput{})
		if err != nil {
			return err
		}
		if resp.Groups == nil {
			return fmt.Errorf("groups list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "AddUserToGroup", func() error {
		_, err := client.AddUserToGroup(ctx, &iam.AddUserToGroupInput{
			GroupName: aws.String(groupName),
			UserName:  aws.String(userName),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "RemoveUserFromGroup", func() error {
		_, err := client.RemoveUserFromGroup(ctx, &iam.RemoveUserFromGroupInput{
			GroupName: aws.String(groupName),
			UserName:  aws.String(userName),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "CreateRole", func() error {
		assumeRolePolicy := `{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {"Service": "lambda.amazonaws.com"},
				"Action": "sts:AssumeRole"
			}]
		}`
		_, err := client.CreateRole(ctx, &iam.CreateRoleInput{
			RoleName:                 aws.String(roleName),
			AssumeRolePolicyDocument: aws.String(assumeRolePolicy),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "CreateRole_InvalidName", func() error {
		_, err := client.CreateRole(ctx, &iam.CreateRoleInput{
			RoleName:                 aws.String("invalid:role-name"),
			AssumeRolePolicyDocument: aws.String(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":["lambda.amazonaws.com"]},"Action":["sts:AssumeRole"}]}`),
		})
		if err == nil {
			return fmt.Errorf("expected error for invalid role name with colon")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetRole", func() error {
		resp, err := client.GetRole(ctx, &iam.GetRoleInput{
			RoleName: aws.String(roleName),
		})
		if err != nil {
			return err
		}
		if resp.Role == nil {
			return fmt.Errorf("role is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListRoles", func() error {
		resp, err := client.ListRoles(ctx, &iam.ListRolesInput{})
		if err != nil {
			return err
		}
		if resp.Roles == nil {
			return fmt.Errorf("roles list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UpdateRoleDescription", func() error {
		_, err := client.UpdateRoleDescription(ctx, &iam.UpdateRoleDescriptionInput{
			RoleName:    aws.String(roleName),
			Description: aws.String("Updated role description"),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "TagRole", func() error {
		_, err := client.TagRole(ctx, &iam.TagRoleInput{
			RoleName: aws.String(roleName),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "ListRoleTags", func() error {
		resp, err := client.ListRoleTags(ctx, &iam.ListRoleTagsInput{
			RoleName: aws.String(roleName),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil {
			return fmt.Errorf("tags is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UntagRole", func() error {
		_, err := client.UntagRole(ctx, &iam.UntagRoleInput{
			RoleName: aws.String(roleName),
			TagKeys:  []string{"Environment"},
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "CreatePolicy", func() error {
		policyDocument := `{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Action": "s3:*",
				"Resource": "*"
			}]
		}`
		_, err := client.CreatePolicy(ctx, &iam.CreatePolicyInput{
			PolicyName:     aws.String(policyName),
			PolicyDocument: aws.String(policyDocument),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "ListPolicies", func() error {
		resp, err := client.ListPolicies(ctx, &iam.ListPoliciesInput{})
		if err != nil {
			return err
		}
		if resp.Policies == nil {
			return fmt.Errorf("policies list is nil")
		}
		return nil
	}))

	profileName := fmt.Sprintf("TestProfile-%d", time.Now().UnixNano())
	results = append(results, r.RunTest("iam", "CreateInstanceProfile", func() error {
		_, err := client.CreateInstanceProfile(ctx, &iam.CreateInstanceProfileInput{
			InstanceProfileName: aws.String(profileName),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "GetInstanceProfile", func() error {
		resp, err := client.GetInstanceProfile(ctx, &iam.GetInstanceProfileInput{
			InstanceProfileName: aws.String(profileName),
		})
		if err != nil {
			return err
		}
		if resp.InstanceProfile == nil {
			return fmt.Errorf("instance profile is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListInstanceProfiles", func() error {
		resp, err := client.ListInstanceProfiles(ctx, &iam.ListInstanceProfilesInput{})
		if err != nil {
			return err
		}
		if resp.InstanceProfiles == nil {
			return fmt.Errorf("instance profiles list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "AddRoleToInstanceProfile", func() error {
		_, err := client.AddRoleToInstanceProfile(ctx, &iam.AddRoleToInstanceProfileInput{
			InstanceProfileName: aws.String(profileName),
			RoleName:            aws.String(roleName),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "RemoveRoleFromInstanceProfile", func() error {
		_, err := client.RemoveRoleFromInstanceProfile(ctx, &iam.RemoveRoleFromInstanceProfileInput{
			InstanceProfileName: aws.String(profileName),
			RoleName:            aws.String(roleName),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "DeleteInstanceProfile", func() error {
		_, err := client.DeleteInstanceProfile(ctx, &iam.DeleteInstanceProfileInput{
			InstanceProfileName: aws.String(profileName),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "PutUserPolicy", func() error {
		policyDocument := `{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Action": "s3:GetObject",
				"Resource": "*"
			}]
		}`
		_, err := client.PutUserPolicy(ctx, &iam.PutUserPolicyInput{
			UserName:       aws.String(userName),
			PolicyName:     aws.String(userInlinePolicyName),
			PolicyDocument: aws.String(policyDocument),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "GetUserPolicy", func() error {
		resp, err := client.GetUserPolicy(ctx, &iam.GetUserPolicyInput{
			UserName:   aws.String(userName),
			PolicyName: aws.String(userInlinePolicyName),
		})
		if err != nil {
			return err
		}
		if resp.PolicyDocument == nil || *resp.PolicyDocument == "" {
			return fmt.Errorf("policy document is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListUserPolicies", func() error {
		resp, err := client.ListUserPolicies(ctx, &iam.ListUserPoliciesInput{
			UserName: aws.String(userName),
		})
		if err != nil {
			return err
		}
		if resp.PolicyNames == nil {
			return fmt.Errorf("policy names list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "PutRolePolicy", func() error {
		policyDocument := `{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Action": "logs:*",
				"Resource": "*"
			}]
		}`
		_, err := client.PutRolePolicy(ctx, &iam.PutRolePolicyInput{
			RoleName:       aws.String(roleName),
			PolicyName:     aws.String(roleInlinePolicyName),
			PolicyDocument: aws.String(policyDocument),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "GetRolePolicy", func() error {
		resp, err := client.GetRolePolicy(ctx, &iam.GetRolePolicyInput{
			RoleName:   aws.String(roleName),
			PolicyName: aws.String(roleInlinePolicyName),
		})
		if err != nil {
			return err
		}
		if resp.PolicyDocument == nil || *resp.PolicyDocument == "" {
			return fmt.Errorf("policy document is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListRolePolicies", func() error {
		resp, err := client.ListRolePolicies(ctx, &iam.ListRolePoliciesInput{
			RoleName: aws.String(roleName),
		})
		if err != nil {
			return err
		}
		if resp.PolicyNames == nil {
			return fmt.Errorf("policy names list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetAccountSummary", func() error {
		resp, err := client.GetAccountSummary(ctx, &iam.GetAccountSummaryInput{})
		if err != nil {
			return err
		}
		if resp.SummaryMap == nil {
			return fmt.Errorf("summary map is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteUserPolicy", func() error {
		_, err := client.DeleteUserPolicy(ctx, &iam.DeleteUserPolicyInput{
			UserName:   aws.String(userName),
			PolicyName: aws.String(userInlinePolicyName),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "DeleteRolePolicy", func() error {
		_, err := client.DeleteRolePolicy(ctx, &iam.DeleteRolePolicyInput{
			RoleName:   aws.String(roleName),
			PolicyName: aws.String(roleInlinePolicyName),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "DeleteLoginProfile", func() error {
		_, err := client.DeleteLoginProfile(ctx, &iam.DeleteLoginProfileInput{
			UserName: aws.String(userName),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "DeleteAccessKey", func() error {
		_, err := client.DeleteAccessKey(ctx, &iam.DeleteAccessKeyInput{
			UserName:    aws.String(userName),
			AccessKeyId: aws.String(accessKeyId),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "DeleteUser", func() error {
		_, err := client.DeleteUser(ctx, &iam.DeleteUserInput{
			UserName: aws.String(userName),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "DeleteGroup", func() error {
		_, err := client.DeleteGroup(ctx, &iam.DeleteGroupInput{
			GroupName: aws.String(groupName),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "DeleteRole", func() error {
		_, err := client.DeleteRole(ctx, &iam.DeleteRoleInput{
			RoleName: aws.String(roleName),
		})
		return err
	}))

	return results
}

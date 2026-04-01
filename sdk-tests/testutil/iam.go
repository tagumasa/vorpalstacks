package testutil

import (
	"context"
	"fmt"
	"strings"
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

	ts := fmt.Sprintf("%d", time.Now().UnixNano())
	userName := fmt.Sprintf("TestUser-%s", ts)
	groupName := fmt.Sprintf("TestGroup-%s", ts)
	roleName := fmt.Sprintf("TestRole-%s", ts)
	policyName := fmt.Sprintf("TestPolicy-%s", ts)
	profileName := fmt.Sprintf("TestProfile-%s", ts)
	userInlinePolicyName := fmt.Sprintf("UserPolicy-%s", ts)
	roleInlinePolicyName := fmt.Sprintf("RolePolicy-%s", ts)
	groupInlinePolicyName := fmt.Sprintf("GroupPolicy-%s", ts)
	var accessKeyId string
	var policyArn string

	// ========== CREATE PHASE ==========

	results = append(results, r.RunTest("iam", "CreateUser", func() error {
		resp, err := client.CreateUser(ctx, &iam.CreateUserInput{
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
		if resp.AccessKey == nil {
			return fmt.Errorf("access key is nil")
		}
		if resp.AccessKey.AccessKeyId == nil {
			return fmt.Errorf("access key id is nil")
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
		resp, err := client.CreateLoginProfile(ctx, &iam.CreateLoginProfileInput{
			UserName: aws.String(userName),
			Password: aws.String("TempPassword123!"),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
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
		newUserName := fmt.Sprintf("UpdatedUser-%s", ts)
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
		resp, err := client.CreateGroup(ctx, &iam.CreateGroupInput{
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

	results = append(results, r.RunTest("iam", "UpdateGroup", func() error {
		_, err := client.UpdateGroup(ctx, &iam.UpdateGroupInput{
			GroupName:    aws.String(groupName),
			NewGroupName: aws.String(groupName + "-renamed"),
		})
		if err == nil {
			groupName = groupName + "-renamed"
		}
		return err
	}))

	results = append(results, r.RunTest("iam", "AddUserToGroup", func() error {
		_, err := client.AddUserToGroup(ctx, &iam.AddUserToGroupInput{
			GroupName: aws.String(groupName),
			UserName:  aws.String(userName),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "ListGroupsForUser", func() error {
		resp, err := client.ListGroupsForUser(ctx, &iam.ListGroupsForUserInput{
			UserName: aws.String(userName),
		})
		if err != nil {
			return err
		}
		if resp.Groups == nil {
			return fmt.Errorf("groups list is nil")
		}
		found := false
		for _, g := range resp.Groups {
			if g.GroupName != nil && *g.GroupName == groupName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("group %s not found in ListGroupsForUser response", groupName)
		}
		return nil
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
		resp, err := client.CreateRole(ctx, &iam.CreateRoleInput{
			RoleName:                 aws.String(roleName),
			AssumeRolePolicyDocument: aws.String(assumeRolePolicy),
		})
		if err != nil {
			return err
		}
		if resp.Role == nil {
			return fmt.Errorf("role is nil")
		}
		if resp.Role.Arn == nil {
			return fmt.Errorf("role arn is nil")
		}
		return nil
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

	results = append(results, r.RunTest("iam", "UpdateRole", func() error {
		_, err := client.UpdateRole(ctx, &iam.UpdateRoleInput{
			RoleName:           aws.String(roleName),
			MaxSessionDuration: aws.Int32(3600),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "UpdateAssumeRolePolicy", func() error {
		newTrustPolicy := `{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Principal": {"Service": "ec2.amazonaws.com"},
				"Action": "sts:AssumeRole"
			}]
		}`
		_, err := client.UpdateAssumeRolePolicy(ctx, &iam.UpdateAssumeRolePolicyInput{
			RoleName:       aws.String(roleName),
			PolicyDocument: aws.String(newTrustPolicy),
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
		resp, err := client.CreatePolicy(ctx, &iam.CreatePolicyInput{
			PolicyName:     aws.String(policyName),
			PolicyDocument: aws.String(policyDocument),
		})
		if err != nil {
			return err
		}
		if resp.Policy == nil {
			return fmt.Errorf("policy is nil")
		}
		if resp.Policy.Arn == nil {
			return fmt.Errorf("policy arn is nil")
		}
		policyArn = *resp.Policy.Arn
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetPolicy", func() error {
		resp, err := client.GetPolicy(ctx, &iam.GetPolicyInput{
			PolicyArn: aws.String(policyArn),
		})
		if err != nil {
			return err
		}
		if resp.Policy == nil {
			return fmt.Errorf("policy is nil")
		}
		if resp.Policy.PolicyName == nil || *resp.Policy.PolicyName != policyName {
			return fmt.Errorf("policy name mismatch: got %v, want %s", resp.Policy.PolicyName, policyName)
		}
		return nil
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

	results = append(results, r.RunTest("iam", "TagPolicy", func() error {
		_, err := client.TagPolicy(ctx, &iam.TagPolicyInput{
			PolicyArn: aws.String(policyArn),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "ListPolicyTags", func() error {
		resp, err := client.ListPolicyTags(ctx, &iam.ListPolicyTagsInput{
			PolicyArn: aws.String(policyArn),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil {
			return fmt.Errorf("tags is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UntagPolicy", func() error {
		_, err := client.UntagPolicy(ctx, &iam.UntagPolicyInput{
			PolicyArn: aws.String(policyArn),
			TagKeys:   []string{"Environment"},
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "CreateInstanceProfile", func() error {
		resp, err := client.CreateInstanceProfile(ctx, &iam.CreateInstanceProfileInput{
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

	results = append(results, r.RunTest("iam", "ListInstanceProfilesForRole", func() error {
		resp, err := client.ListInstanceProfilesForRole(ctx, &iam.ListInstanceProfilesForRoleInput{
			RoleName: aws.String(roleName),
		})
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

	results = append(results, r.RunTest("iam", "TagInstanceProfile", func() error {
		_, err := client.TagInstanceProfile(ctx, &iam.TagInstanceProfileInput{
			InstanceProfileName: aws.String(profileName),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "ListInstanceProfileTags", func() error {
		resp, err := client.ListInstanceProfileTags(ctx, &iam.ListInstanceProfileTagsInput{
			InstanceProfileName: aws.String(profileName),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil {
			return fmt.Errorf("tags is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UntagInstanceProfile", func() error {
		_, err := client.UntagInstanceProfile(ctx, &iam.UntagInstanceProfileInput{
			InstanceProfileName: aws.String(profileName),
			TagKeys:             []string{"Environment"},
		})
		return err
	}))

	// ========== INTERACTION PHASE ==========

	results = append(results, r.RunTest("iam", "AttachUserPolicy", func() error {
		_, err := client.AttachUserPolicy(ctx, &iam.AttachUserPolicyInput{
			UserName:  aws.String(userName),
			PolicyArn: aws.String(policyArn),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "ListAttachedUserPolicies", func() error {
		resp, err := client.ListAttachedUserPolicies(ctx, &iam.ListAttachedUserPoliciesInput{
			UserName: aws.String(userName),
		})
		if err != nil {
			return err
		}
		if resp.AttachedPolicies == nil {
			return fmt.Errorf("attached policies list is nil")
		}
		found := false
		for _, p := range resp.AttachedPolicies {
			if p.PolicyArn != nil && *p.PolicyArn == policyArn {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("policy %s not found in ListAttachedUserPolicies", policyArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DetachUserPolicy", func() error {
		_, err := client.DetachUserPolicy(ctx, &iam.DetachUserPolicyInput{
			UserName:  aws.String(userName),
			PolicyArn: aws.String(policyArn),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "PutUserPermissionsBoundary", func() error {
		_, err := client.PutUserPermissionsBoundary(ctx, &iam.PutUserPermissionsBoundaryInput{
			UserName:            aws.String(userName),
			PermissionsBoundary: aws.String(policyArn),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "GetUser_PermissionsBoundary", func() error {
		resp, err := client.GetUser(ctx, &iam.GetUserInput{
			UserName: aws.String(userName),
		})
		if err != nil {
			return err
		}
		if resp.User == nil {
			return fmt.Errorf("user is nil")
		}
		if resp.User.PermissionsBoundary == nil {
			return fmt.Errorf("permissions boundary is nil")
		}
		if resp.User.PermissionsBoundary.PermissionsBoundaryArn == nil || *resp.User.PermissionsBoundary.PermissionsBoundaryArn != policyArn {
			return fmt.Errorf("permissions boundary arn mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteUserPermissionsBoundary", func() error {
		_, err := client.DeleteUserPermissionsBoundary(ctx, &iam.DeleteUserPermissionsBoundaryInput{
			UserName: aws.String(userName),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "AttachGroupPolicy", func() error {
		_, err := client.AttachGroupPolicy(ctx, &iam.AttachGroupPolicyInput{
			GroupName: aws.String(groupName),
			PolicyArn: aws.String(policyArn),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "ListAttachedGroupPolicies", func() error {
		resp, err := client.ListAttachedGroupPolicies(ctx, &iam.ListAttachedGroupPoliciesInput{
			GroupName: aws.String(groupName),
		})
		if err != nil {
			return err
		}
		if resp.AttachedPolicies == nil {
			return fmt.Errorf("attached policies list is nil")
		}
		found := false
		for _, p := range resp.AttachedPolicies {
			if p.PolicyArn != nil && *p.PolicyArn == policyArn {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("policy %s not found in ListAttachedGroupPolicies", policyArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListEntitiesForPolicy_Role", func() error {
		_, err := client.AttachRolePolicy(ctx, &iam.AttachRolePolicyInput{
			RoleName:  aws.String(roleName),
			PolicyArn: aws.String(policyArn),
		})
		if err != nil {
			return err
		}
		resp, err := client.ListEntitiesForPolicy(ctx, &iam.ListEntitiesForPolicyInput{
			PolicyArn: aws.String(policyArn),
		})
		if err != nil {
			return err
		}
		if resp.PolicyRoles == nil {
			return fmt.Errorf("policy roles list is nil")
		}
		if resp.PolicyGroups == nil {
			return fmt.Errorf("policy groups list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DetachGroupPolicy", func() error {
		_, err := client.DetachGroupPolicy(ctx, &iam.DetachGroupPolicyInput{
			GroupName: aws.String(groupName),
			PolicyArn: aws.String(policyArn),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "PutGroupPolicy", func() error {
		policyDocument := `{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Action": "logs:*",
				"Resource": "*"
			}]
		}`
		_, err := client.PutGroupPolicy(ctx, &iam.PutGroupPolicyInput{
			GroupName:      aws.String(groupName),
			PolicyName:     aws.String(groupInlinePolicyName),
			PolicyDocument: aws.String(policyDocument),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "GetGroupPolicy", func() error {
		resp, err := client.GetGroupPolicy(ctx, &iam.GetGroupPolicyInput{
			GroupName:  aws.String(groupName),
			PolicyName: aws.String(groupInlinePolicyName),
		})
		if err != nil {
			return err
		}
		if resp.PolicyDocument == nil || *resp.PolicyDocument == "" {
			return fmt.Errorf("policy document is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListGroupPolicies", func() error {
		resp, err := client.ListGroupPolicies(ctx, &iam.ListGroupPoliciesInput{
			GroupName: aws.String(groupName),
		})
		if err != nil {
			return err
		}
		if resp.PolicyNames == nil {
			return fmt.Errorf("policy names list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteGroupPolicy", func() error {
		_, err := client.DeleteGroupPolicy(ctx, &iam.DeleteGroupPolicyInput{
			GroupName:  aws.String(groupName),
			PolicyName: aws.String(groupInlinePolicyName),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "AttachRolePolicy_FullCycle", func() error {
		_, err := client.AttachRolePolicy(ctx, &iam.AttachRolePolicyInput{
			RoleName:  aws.String(roleName),
			PolicyArn: aws.String(policyArn),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "ListAttachedRolePolicies", func() error {
		resp, err := client.ListAttachedRolePolicies(ctx, &iam.ListAttachedRolePoliciesInput{
			RoleName: aws.String(roleName),
		})
		if err != nil {
			return err
		}
		if resp.AttachedPolicies == nil {
			return fmt.Errorf("attached policies list is nil")
		}
		found := false
		for _, p := range resp.AttachedPolicies {
			if p.PolicyArn != nil && *p.PolicyArn == policyArn {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("policy %s not found in ListAttachedRolePolicies", policyArn)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DetachRolePolicy", func() error {
		_, err := client.DetachRolePolicy(ctx, &iam.DetachRolePolicyInput{
			RoleName:  aws.String(roleName),
			PolicyArn: aws.String(policyArn),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "PutRolePermissionsBoundary", func() error {
		_, err := client.PutRolePermissionsBoundary(ctx, &iam.PutRolePermissionsBoundaryInput{
			RoleName:            aws.String(roleName),
			PermissionsBoundary: aws.String(policyArn),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "GetRole_PermissionsBoundary", func() error {
		resp, err := client.GetRole(ctx, &iam.GetRoleInput{
			RoleName: aws.String(roleName),
		})
		if err != nil {
			return err
		}
		if resp.Role == nil {
			return fmt.Errorf("role is nil")
		}
		if resp.Role.PermissionsBoundary == nil {
			return fmt.Errorf("permissions boundary is nil")
		}
		if resp.Role.PermissionsBoundary.PermissionsBoundaryArn == nil || *resp.Role.PermissionsBoundary.PermissionsBoundaryArn != policyArn {
			return fmt.Errorf("permissions boundary arn mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteRolePermissionsBoundary", func() error {
		_, err := client.DeleteRolePermissionsBoundary(ctx, &iam.DeleteRolePermissionsBoundaryInput{
			RoleName: aws.String(roleName),
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

	results = append(results, r.RunTest("iam", "DeleteUserPolicy", func() error {
		_, err := client.DeleteUserPolicy(ctx, &iam.DeleteUserPolicyInput{
			UserName:   aws.String(userName),
			PolicyName: aws.String(userInlinePolicyName),
		})
		return err
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

	results = append(results, r.RunTest("iam", "DeleteRolePolicy", func() error {
		_, err := client.DeleteRolePolicy(ctx, &iam.DeleteRolePolicyInput{
			RoleName:   aws.String(roleName),
			PolicyName: aws.String(roleInlinePolicyName),
		})
		return err
	}))

	// ========== POLICY VERSIONING ==========

	results = append(results, r.RunTest("iam", "CreatePolicyVersion", func() error {
		v2Document := `{
			"Version": "2012-10-17",
			"Statement": [{
				"Effect": "Allow",
				"Action": "ec2:*",
				"Resource": "*"
			}]
		}`
		resp, err := client.CreatePolicyVersion(ctx, &iam.CreatePolicyVersionInput{
			PolicyArn:      aws.String(policyArn),
			PolicyDocument: aws.String(v2Document),
			SetAsDefault:   false,
		})
		if err != nil {
			return err
		}
		if resp.PolicyVersion == nil {
			return fmt.Errorf("policy version is nil")
		}
		if resp.PolicyVersion.VersionId == nil {
			return fmt.Errorf("version id is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListPolicyVersions", func() error {
		resp, err := client.ListPolicyVersions(ctx, &iam.ListPolicyVersionsInput{
			PolicyArn: aws.String(policyArn),
		})
		if err != nil {
			return err
		}
		if resp.Versions == nil {
			return fmt.Errorf("versions list is nil")
		}
		if len(resp.Versions) < 2 {
			return fmt.Errorf("expected at least 2 policy versions, got %d", len(resp.Versions))
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "SetDefaultPolicyVersion", func() error {
		resp, err := client.ListPolicyVersions(ctx, &iam.ListPolicyVersionsInput{
			PolicyArn: aws.String(policyArn),
		})
		if err != nil {
			return err
		}
		var nonDefaultVersionId string
		for _, v := range resp.Versions {
			if v.VersionId != nil && !v.IsDefaultVersion {
				nonDefaultVersionId = *v.VersionId
				break
			}
		}
		if nonDefaultVersionId == "" {
			return fmt.Errorf("no non-default version found")
		}
		_, err = client.SetDefaultPolicyVersion(ctx, &iam.SetDefaultPolicyVersionInput{
			PolicyArn: aws.String(policyArn),
			VersionId: aws.String(nonDefaultVersionId),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "GetPolicyVersion", func() error {
		resp, err := client.ListPolicyVersions(ctx, &iam.ListPolicyVersionsInput{
			PolicyArn: aws.String(policyArn),
		})
		if err != nil {
			return err
		}
		var defaultVersionId string
		for _, v := range resp.Versions {
			if v.IsDefaultVersion && v.VersionId != nil {
				defaultVersionId = *v.VersionId
				break
			}
		}
		if defaultVersionId == "" {
			return fmt.Errorf("no default version found")
		}
		getResp, err := client.GetPolicyVersion(ctx, &iam.GetPolicyVersionInput{
			PolicyArn: aws.String(policyArn),
			VersionId: aws.String(defaultVersionId),
		})
		if err != nil {
			return err
		}
		if getResp.PolicyVersion == nil {
			return fmt.Errorf("policy version is nil")
		}
		if !getResp.PolicyVersion.IsDefaultVersion {
			return fmt.Errorf("expected default version")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeletePolicyVersion", func() error {
		resp, err := client.ListPolicyVersions(ctx, &iam.ListPolicyVersionsInput{
			PolicyArn: aws.String(policyArn),
		})
		if err != nil {
			return err
		}
		var nonDefaultVersionId string
		for _, v := range resp.Versions {
			if v.VersionId != nil && !v.IsDefaultVersion {
				nonDefaultVersionId = *v.VersionId
				break
			}
		}
		if nonDefaultVersionId == "" {
			return fmt.Errorf("no non-default version found to delete")
		}
		_, err = client.DeletePolicyVersion(ctx, &iam.DeletePolicyVersionInput{
			PolicyArn: aws.String(policyArn),
			VersionId: aws.String(nonDefaultVersionId),
		})
		return err
	}))

	// ========== ACCESS KEY OPERATIONS ==========

	results = append(results, r.RunTest("iam", "UpdateAccessKey", func() error {
		_, err := client.UpdateAccessKey(ctx, &iam.UpdateAccessKeyInput{
			AccessKeyId: aws.String(accessKeyId),
			Status:      types.StatusTypeInactive,
			UserName:    aws.String(userName),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "GetAccessKeyLastUsed", func() error {
		resp, err := client.GetAccessKeyLastUsed(ctx, &iam.GetAccessKeyLastUsedInput{
			AccessKeyId: aws.String(accessKeyId),
		})
		if err != nil {
			return err
		}
		if resp.UserName == nil {
			return fmt.Errorf("username is nil")
		}
		return nil
	}))

	// ========== LOGIN PROFILE OPERATIONS ==========

	results = append(results, r.RunTest("iam", "UpdateLoginProfile", func() error {
		_, err := client.UpdateLoginProfile(ctx, &iam.UpdateLoginProfileInput{
			UserName:              aws.String(userName),
			Password:              aws.String("NewPassword456!"),
			PasswordResetRequired: aws.Bool(true),
		})
		return err
	}))

	// ========== ACCOUNT OPERATIONS ==========

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

	results = append(results, r.RunTest("iam", "GetAccountAuthorizationDetails", func() error {
		resp, err := client.GetAccountAuthorizationDetails(ctx, &iam.GetAccountAuthorizationDetailsInput{
			Filter: []types.EntityType{types.EntityTypeUser},
		})
		if err != nil {
			return err
		}
		if resp.UserDetailList == nil {
			return fmt.Errorf("user detail list is nil")
		}
		return nil
	}))

	accountAlias := "test-alias-" + ts

	results = append(results, r.RunTest("iam", "CreateAccountAlias", func() error {
		_, err := client.CreateAccountAlias(ctx, &iam.CreateAccountAliasInput{
			AccountAlias: aws.String(accountAlias),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "ListAccountAliases", func() error {
		resp, err := client.ListAccountAliases(ctx, &iam.ListAccountAliasesInput{})
		if err != nil {
			return err
		}
		if resp.AccountAliases == nil {
			return fmt.Errorf("account aliases list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteAccountAlias", func() error {
		_, err := client.DeleteAccountAlias(ctx, &iam.DeleteAccountAliasInput{
			AccountAlias: aws.String(accountAlias),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "UpdateAccountPasswordPolicy", func() error {
		_, err := client.UpdateAccountPasswordPolicy(ctx, &iam.UpdateAccountPasswordPolicyInput{
			MinimumPasswordLength:      aws.Int32(12),
			RequireUppercaseCharacters: true,
			RequireLowercaseCharacters: true,
			RequireNumbers:             true,
			RequireSymbols:             true,
			AllowUsersToChangePassword: true,
			MaxPasswordAge:             aws.Int32(90),
			PasswordReusePrevention:    aws.Int32(5),
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "GetAccountPasswordPolicy", func() error {
		resp, err := client.GetAccountPasswordPolicy(ctx, &iam.GetAccountPasswordPolicyInput{})
		if err != nil {
			return err
		}
		if resp.PasswordPolicy == nil {
			return fmt.Errorf("password policy is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteAccountPasswordPolicy", func() error {
		_, err := client.DeleteAccountPasswordPolicy(ctx, &iam.DeleteAccountPasswordPolicyInput{})
		return err
	}))

	// ========== SERVICE-LINKED ROLE ==========

	var serviceLinkedRoleName string

	results = append(results, r.RunTest("iam", "CreateServiceLinkedRole", func() error {
		resp, err := client.CreateServiceLinkedRole(ctx, &iam.CreateServiceLinkedRoleInput{
			AWSServiceName: aws.String("lambda.amazonaws.com"),
			Description:    aws.String("Test service-linked role"),
		})
		if err != nil {
			return err
		}
		if resp.Role == nil {
			return fmt.Errorf("role is nil")
		}
		if resp.Role.RoleName == nil {
			return fmt.Errorf("role name is nil")
		}
		serviceLinkedRoleName = *resp.Role.RoleName
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteServiceLinkedRole", func() error {
		resp, err := client.DeleteServiceLinkedRole(ctx, &iam.DeleteServiceLinkedRoleInput{
			RoleName: aws.String(serviceLinkedRoleName),
		})
		if err != nil {
			return err
		}
		if resp.DeletionTaskId == nil {
			return fmt.Errorf("deletion task id is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetServiceLinkedRoleDeletionStatus", func() error {
		_, err := client.GetServiceLinkedRoleDeletionStatus(ctx, &iam.GetServiceLinkedRoleDeletionStatusInput{
			DeletionTaskId: aws.String("test-task-id"),
		})
		return err
	}))

	// ========== SAML PROVIDER ==========

	samlProviderName := fmt.Sprintf("TestSAML-%s", ts)
	samlMetadata := `<?xml version="1.0" encoding="UTF-8"?>
<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata" entityID="https://idp.example.com">
  <md:IDPSSODescriptor protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
    <md:SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="https://idp.example.com/sso"/>
  </md:IDPSSODescriptor>
</md:EntityDescriptor>`
	var samlProviderArn string

	results = append(results, r.RunTest("iam", "CreateSAMLProvider", func() error {
		resp, err := client.CreateSAMLProvider(ctx, &iam.CreateSAMLProviderInput{
			Name:                 aws.String(samlProviderName),
			SAMLMetadataDocument: aws.String(samlMetadata),
			Tags: []types.Tag{
				{Key: aws.String("Source"), Value: aws.String("test")},
			},
		})
		if err != nil {
			return err
		}
		if resp.SAMLProviderArn == nil {
			return fmt.Errorf("saml provider arn is nil")
		}
		samlProviderArn = *resp.SAMLProviderArn
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetSAMLProvider", func() error {
		resp, err := client.GetSAMLProvider(ctx, &iam.GetSAMLProviderInput{
			SAMLProviderArn: aws.String(samlProviderArn),
		})
		if err != nil {
			return err
		}
		if resp.SAMLMetadataDocument == nil || *resp.SAMLMetadataDocument == "" {
			return fmt.Errorf("saml metadata document is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListSAMLProviders", func() error {
		resp, err := client.ListSAMLProviders(ctx, &iam.ListSAMLProvidersInput{})
		if err != nil {
			return err
		}
		if resp.SAMLProviderList == nil {
			return fmt.Errorf("saml provider list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UpdateSAMLProvider", func() error {
		resp, err := client.UpdateSAMLProvider(ctx, &iam.UpdateSAMLProviderInput{
			SAMLProviderArn:      aws.String(samlProviderArn),
			SAMLMetadataDocument: aws.String(samlMetadata),
		})
		if err != nil {
			return err
		}
		if resp.SAMLProviderArn == nil {
			return fmt.Errorf("saml provider arn is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "TagSAMLProvider", func() error {
		_, err := client.TagSAMLProvider(ctx, &iam.TagSAMLProviderInput{
			SAMLProviderArn: aws.String(samlProviderArn),
			Tags: []types.Tag{
				{Key: aws.String("Environment"), Value: aws.String("test")},
			},
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "ListSAMLProviderTags", func() error {
		resp, err := client.ListSAMLProviderTags(ctx, &iam.ListSAMLProviderTagsInput{
			SAMLProviderArn: aws.String(samlProviderArn),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil {
			return fmt.Errorf("tags is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UntagSAMLProvider", func() error {
		_, err := client.UntagSAMLProvider(ctx, &iam.UntagSAMLProviderInput{
			SAMLProviderArn: aws.String(samlProviderArn),
			TagKeys:         []string{"Environment"},
		})
		return err
	}))

	results = append(results, r.RunTest("iam", "DeleteSAMLProvider", func() error {
		_, err := client.DeleteSAMLProvider(ctx, &iam.DeleteSAMLProviderInput{
			SAMLProviderArn: aws.String(samlProviderArn),
		})
		return err
	}))

	// ========== VIRTUAL MFA DEVICE ==========

	virtualMFADeviceName := fmt.Sprintf("TestMFA-%s", ts)
	var virtualMFASerial string

	results = append(results, r.RunTest("iam", "CreateVirtualMFADevice", func() error {
		resp, err := client.CreateVirtualMFADevice(ctx, &iam.CreateVirtualMFADeviceInput{
			VirtualMFADeviceName: aws.String(virtualMFADeviceName),
			Tags: []types.Tag{
				{Key: aws.String("Purpose"), Value: aws.String("test")},
			},
		})
		if err != nil {
			return err
		}
		if resp.VirtualMFADevice == nil {
			return fmt.Errorf("virtual mfa device is nil")
		}
		if resp.VirtualMFADevice.SerialNumber == nil {
			return fmt.Errorf("serial number is nil")
		}
		virtualMFASerial = *resp.VirtualMFADevice.SerialNumber
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListVirtualMFADevices", func() error {
		resp, err := client.ListVirtualMFADevices(ctx, &iam.ListVirtualMFADevicesInput{})
		if err != nil {
			return err
		}
		if resp.VirtualMFADevices == nil {
			return fmt.Errorf("virtual mfa devices list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteVirtualMFADevice", func() error {
		_, err := client.DeleteVirtualMFADevice(ctx, &iam.DeleteVirtualMFADeviceInput{
			SerialNumber: aws.String(virtualMFASerial),
		})
		return err
	}))

	// ========== ERROR PATH TESTS ==========

	results = append(results, r.RunTest("iam", "Error_DeleteNonExistentUser", func() error {
		_, err := client.DeleteUser(ctx, &iam.DeleteUserInput{
			UserName: aws.String("NonExistentUser-" + ts),
		})
		if err == nil {
			return fmt.Errorf("expected error for deleting non-existent user")
		}
		if !strings.Contains(err.Error(), "NoSuchEntity") {
			return fmt.Errorf("expected NoSuchEntity error, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "Error_GetNonExistentRole", func() error {
		_, err := client.GetRole(ctx, &iam.GetRoleInput{
			RoleName: aws.String("NonExistentRole-" + ts),
		})
		if err == nil {
			return fmt.Errorf("expected error for getting non-existent role")
		}
		if !strings.Contains(err.Error(), "NoSuchEntity") {
			return fmt.Errorf("expected NoSuchEntity error, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "Error_AttachPolicyToNonExistentUser", func() error {
		_, err := client.AttachUserPolicy(ctx, &iam.AttachUserPolicyInput{
			UserName:  aws.String("NonExistentUser-" + ts),
			PolicyArn: aws.String(policyArn),
		})
		if err == nil {
			return fmt.Errorf("expected error for attaching policy to non-existent user")
		}
		if !strings.Contains(err.Error(), "NoSuchEntity") {
			return fmt.Errorf("expected NoSuchEntity error, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "Error_DeleteDefaultPolicyVersion", func() error {
		resp, err := client.ListPolicyVersions(ctx, &iam.ListPolicyVersionsInput{
			PolicyArn: aws.String(policyArn),
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
		_, err = client.DeletePolicyVersion(ctx, &iam.DeletePolicyVersionInput{
			PolicyArn: aws.String(policyArn),
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
		_, err := client.CreateUser(ctx, &iam.CreateUserInput{
			UserName: aws.String(userName),
		})
		if err == nil {
			return fmt.Errorf("expected error for creating duplicate user")
		}
		if !strings.Contains(err.Error(), "EntityAlreadyExists") {
			return fmt.Errorf("expected EntityAlreadyExists error, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "Error_CreateDuplicatePolicy", func() error {
		policyDocument := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"s3:*","Resource":"*"}]}`
		_, err := client.CreatePolicy(ctx, &iam.CreatePolicyInput{
			PolicyName:     aws.String(policyName),
			PolicyDocument: aws.String(policyDocument),
		})
		if err == nil {
			return fmt.Errorf("expected error for creating duplicate policy")
		}
		if !strings.Contains(err.Error(), "EntityAlreadyExists") {
			return fmt.Errorf("expected EntityAlreadyExists error, got: %v", err)
		}
		return nil
	}))

	// ========== DELETE PHASE ==========

	results = append(results, r.RunTest("iam", "DeleteInstanceProfile", func() error {
		_, err := client.DeleteInstanceProfile(ctx, &iam.DeleteInstanceProfileInput{
			InstanceProfileName: aws.String(profileName),
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

	results = append(results, r.RunTest("iam", "DeletePolicy", func() error {
		_, err := client.DeletePolicy(ctx, &iam.DeletePolicyInput{
			PolicyArn: aws.String(policyArn),
		})
		return err
	}))

	return results
}

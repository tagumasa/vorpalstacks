package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func (r *TestRunner) iamGroupTests(tc *iamTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("iam", "CreateGroup", func() error {
		resp, err := tc.client.CreateGroup(tc.ctx, &iam.CreateGroupInput{
			GroupName: aws.String(tc.group),
		})
		if err != nil {
			return err
		}
		if resp.Group == nil {
			return fmt.Errorf("group is nil")
		}
		if aws.ToString(resp.Group.GroupName) != tc.group {
			return fmt.Errorf("group name mismatch: got %s, want %s", aws.ToString(resp.Group.GroupName), tc.group)
		}
		if aws.ToString(resp.Group.Arn) == "" {
			return fmt.Errorf("group arn is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetGroup", func() error {
		resp, err := tc.client.GetGroup(tc.ctx, &iam.GetGroupInput{
			GroupName: aws.String(tc.group),
		})
		if err != nil {
			return err
		}
		if resp.Group == nil {
			return fmt.Errorf("group is nil")
		}
		if aws.ToString(resp.Group.GroupName) != tc.group {
			return fmt.Errorf("group name mismatch")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListGroups", func() error {
		groups, err := iamPaginate(func(marker *string) ([]types.Group, *string, error) {
			resp, err := tc.client.ListGroups(tc.ctx, &iam.ListGroupsInput{Marker: marker})
			if err != nil {
				return nil, nil, err
			}
			return resp.Groups, resp.Marker, nil
		})
		if err != nil {
			return err
		}
		var matched *types.Group
		for i := range groups {
			if aws.ToString(groups[i].GroupName) == tc.group {
				matched = &groups[i]
				break
			}
		}
		if matched == nil {
			return fmt.Errorf("group %s not found in ListGroups", tc.group)
		}
		if aws.ToString(matched.Arn) == "" {
			return fmt.Errorf("group arn is empty in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "UpdateGroup", func() error {
		newName := tc.group + "-renamed"
		_, err := tc.client.UpdateGroup(tc.ctx, &iam.UpdateGroupInput{
			GroupName:    aws.String(tc.group),
			NewGroupName: aws.String(newName),
		})
		if err != nil {
			return err
		}
		tc.group = newName
		resp, err := tc.client.GetGroup(tc.ctx, &iam.GetGroupInput{
			GroupName: aws.String(tc.group),
		})
		if err != nil {
			return err
		}
		if aws.ToString(resp.Group.GroupName) != tc.group {
			return fmt.Errorf("group name not updated")
		}

		// A new name outside the group-name pattern must be rejected, and
		// renaming onto an existing group must fail with EntityAlreadyExists.
		if _, err := tc.client.UpdateGroup(tc.ctx, &iam.UpdateGroupInput{
			GroupName:    aws.String(tc.group),
			NewGroupName: aws.String("invalid group name!"),
		}); err == nil || !isInvalidInputError(err) {
			return fmt.Errorf("invalid NewGroupName: got %v, want InvalidInput", err)
		}
		other := fmt.Sprintf("UpdateGroupOther-%s", tc.ts)
		cleanupOther, err := tc.createGroup(other)
		if err != nil {
			return err
		}
		defer cleanupOther()
		if _, err := tc.client.UpdateGroup(tc.ctx, &iam.UpdateGroupInput{
			GroupName:    aws.String(tc.group),
			NewGroupName: aws.String(other),
		}); err == nil || !containsErrorCode(err, "EntityAlreadyExists") {
			return fmt.Errorf("rename onto existing group: got %v, want EntityAlreadyExists", err)
		}
		return nil
	}))

	// User-Group membership
	results = append(results, r.RunTest("iam", "AddUserToGroup", func() error {
		_, err := tc.client.AddUserToGroup(tc.ctx, &iam.AddUserToGroupInput{
			GroupName: aws.String(tc.group),
			UserName:  aws.String(tc.user),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetGroup(tc.ctx, &iam.GetGroupInput{
			GroupName: aws.String(tc.group),
		})
		if err != nil {
			return fmt.Errorf("GetGroup after AddUserToGroup: %w", err)
		}
		found := false
		for _, u := range resp.Users {
			if aws.ToString(u.UserName) == tc.user {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("user %s not found in group after AddUserToGroup", tc.user)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListGroupsForUser", func() error {
		resp, err := tc.client.ListGroupsForUser(tc.ctx, &iam.ListGroupsForUserInput{
			UserName: aws.String(tc.user),
		})
		if err != nil {
			return err
		}
		found := false
		for _, g := range resp.Groups {
			if aws.ToString(g.GroupName) == tc.group {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("group %s not found in ListGroupsForUser", tc.group)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "AddUserToGroup_Idempotent", func() error {
		// Re-adding an existing membership succeeds without change; the
		// documented error surface has no duplicate-membership error.
		if _, err := tc.client.AddUserToGroup(tc.ctx, &iam.AddUserToGroupInput{
			GroupName: aws.String(tc.group),
			UserName:  aws.String(tc.user),
		}); err != nil {
			return fmt.Errorf("second AddUserToGroup for an existing membership failed: %w", err)
		}

		resp, err := tc.client.GetGroup(tc.ctx, &iam.GetGroupInput{
			GroupName: aws.String(tc.group),
		})
		if err != nil {
			return fmt.Errorf("GetGroup after repeated AddUserToGroup: %w", err)
		}
		count := 0
		for _, u := range resp.Users {
			if aws.ToString(u.UserName) == tc.user {
				count++
			}
		}
		if count != 1 {
			return fmt.Errorf("user %s appears %d times in group, want exactly 1", tc.user, count)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "RemoveUserFromGroup", func() error {
		_, err := tc.client.RemoveUserFromGroup(tc.ctx, &iam.RemoveUserFromGroupInput{
			GroupName: aws.String(tc.group),
			UserName:  aws.String(tc.user),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListGroupsForUser(tc.ctx, &iam.ListGroupsForUserInput{
			UserName: aws.String(tc.user),
		})
		if err != nil {
			return err
		}
		for _, g := range resp.Groups {
			if aws.ToString(g.GroupName) == tc.group {
				return fmt.Errorf("user should be removed from group")
			}
		}
		return nil
	}))

	// Group inline policies
	results = append(results, r.RunTest("iam", "PutGroupPolicy", func() error {
		_, err := tc.client.PutGroupPolicy(tc.ctx, &iam.PutGroupPolicyInput{
			GroupName:      aws.String(tc.group),
			PolicyName:     aws.String(tc.groupInlinePolicy),
			PolicyDocument: aws.String(logsFullAccessPolicy),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.GetGroupPolicy(tc.ctx, &iam.GetGroupPolicyInput{
			GroupName:  aws.String(tc.group),
			PolicyName: aws.String(tc.groupInlinePolicy),
		})
		if err != nil {
			return fmt.Errorf("GetGroupPolicy after PutGroupPolicy: %w", err)
		}
		if resp.PolicyDocument == nil || *resp.PolicyDocument == "" {
			return fmt.Errorf("policy document is empty after PutGroupPolicy")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "GetGroupPolicy", func() error {
		resp, err := tc.client.GetGroupPolicy(tc.ctx, &iam.GetGroupPolicyInput{
			GroupName:  aws.String(tc.group),
			PolicyName: aws.String(tc.groupInlinePolicy),
		})
		if err != nil {
			return err
		}
		if resp.PolicyDocument == nil || *resp.PolicyDocument == "" {
			return fmt.Errorf("policy document is empty")
		}
		if aws.ToString(resp.GroupName) != tc.group {
			return fmt.Errorf("group name mismatch in GetGroupPolicy")
		}
		if aws.ToString(resp.PolicyName) != tc.groupInlinePolicy {
			return fmt.Errorf("policy name mismatch in GetGroupPolicy")
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "ListGroupPolicies", func() error {
		resp, err := tc.client.ListGroupPolicies(tc.ctx, &iam.ListGroupPoliciesInput{
			GroupName: aws.String(tc.group),
		})
		if err != nil {
			return err
		}
		found := false
		for _, name := range resp.PolicyNames {
			if name == tc.groupInlinePolicy {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("inline policy %s not found in ListGroupPolicies", tc.groupInlinePolicy)
		}
		return nil
	}))

	results = append(results, r.RunTest("iam", "DeleteGroupPolicy", func() error {
		_, err := tc.client.DeleteGroupPolicy(tc.ctx, &iam.DeleteGroupPolicyInput{
			GroupName:  aws.String(tc.group),
			PolicyName: aws.String(tc.groupInlinePolicy),
		})
		if err != nil {
			return err
		}
		resp, err := tc.client.ListGroupPolicies(tc.ctx, &iam.ListGroupPoliciesInput{
			GroupName: aws.String(tc.group),
		})
		if err != nil {
			return err
		}
		for _, name := range resp.PolicyNames {
			if name == tc.groupInlinePolicy {
				return fmt.Errorf("inline policy should be deleted")
			}
		}
		return nil
	}))

	return results
}

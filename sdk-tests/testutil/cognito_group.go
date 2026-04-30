package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func (r *TestRunner) cognitoGroupTests(ctx context.Context, client *cognitoidentityprovider.Client, userPoolID string) []TestResult {
	var results []TestResult

	groupName := fmt.Sprintf("group-%d", time.Now().UnixNano())
	results = append(results, r.RunTest("cognito", "CreateGroup", func() error {
		resp, err := client.CreateGroup(ctx, &cognitoidentityprovider.CreateGroupInput{
			GroupName:  aws.String(groupName),
			UserPoolId: aws.String(userPoolID),
		})
		if err != nil {
			return err
		}
		if resp.Group == nil {
			return fmt.Errorf("group is nil")
		}
		if resp.Group.GroupName == nil || *resp.Group.GroupName != groupName {
			return fmt.Errorf("GroupName mismatch: got %v, want %s", resp.Group.GroupName, groupName)
		}
		if resp.Group.UserPoolId == nil || *resp.Group.UserPoolId != userPoolID {
			return fmt.Errorf("UserPoolId mismatch: got %v, want %s", resp.Group.UserPoolId, userPoolID)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "ListGroups", func() error {
		resp, err := client.ListGroups(ctx, &cognitoidentityprovider.ListGroupsInput{
			UserPoolId: aws.String(userPoolID),
		})
		if err != nil {
			return err
		}
		if len(resp.Groups) == 0 {
			return fmt.Errorf("expected at least one group")
		}
		found := false
		for _, g := range resp.Groups {
			if g.GroupName != nil && *g.GroupName == groupName {
				found = true
				if g.UserPoolId == nil || *g.UserPoolId != userPoolID {
					return fmt.Errorf("UserPoolId mismatch in ListGroups")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created group %s not found in ListGroups", groupName)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "GetGroup", func() error {
		newGroupName := fmt.Sprintf("get-group-%d", time.Now().UnixNano())
		_, err := client.CreateGroup(ctx, &cognitoidentityprovider.CreateGroupInput{
			GroupName:  aws.String(newGroupName),
			UserPoolId: aws.String(userPoolID),
		})
		if err != nil {
			return fmt.Errorf("create group: %v", err)
		}
		defer client.DeleteGroup(ctx, &cognitoidentityprovider.DeleteGroupInput{
			GroupName:  aws.String(newGroupName),
			UserPoolId: aws.String(userPoolID),
		})
		resp, err := client.GetGroup(ctx, &cognitoidentityprovider.GetGroupInput{
			GroupName:  aws.String(newGroupName),
			UserPoolId: aws.String(userPoolID),
		})
		if err != nil {
			return fmt.Errorf("GetGroup failed: %v", err)
		}
		if resp.Group == nil {
			return fmt.Errorf("group is nil")
		}
		if resp.Group.GroupName == nil || *resp.Group.GroupName != newGroupName {
			return fmt.Errorf("GroupName mismatch: got %v", resp.Group.GroupName)
		}
		if resp.Group.UserPoolId == nil || *resp.Group.UserPoolId != userPoolID {
			return fmt.Errorf("UserPoolId mismatch: got %v", resp.Group.UserPoolId)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "UpdateGroup", func() error {
		ugGroupName := fmt.Sprintf("ug-group-%d", time.Now().UnixNano())
		_, err := client.CreateGroup(ctx, &cognitoidentityprovider.CreateGroupInput{
			GroupName:   aws.String(ugGroupName),
			UserPoolId:  aws.String(userPoolID),
			Description: aws.String("Original description"),
		})
		if err != nil {
			return fmt.Errorf("create group: %v", err)
		}
		defer client.DeleteGroup(ctx, &cognitoidentityprovider.DeleteGroupInput{
			GroupName:  aws.String(ugGroupName),
			UserPoolId: aws.String(userPoolID),
		})
		_, err = client.UpdateGroup(ctx, &cognitoidentityprovider.UpdateGroupInput{
			GroupName:   aws.String(ugGroupName),
			UserPoolId:  aws.String(userPoolID),
			Description: aws.String("Updated description"),
			Precedence:  aws.Int32(10),
		})
		if err != nil {
			return fmt.Errorf("UpdateGroup failed: %v", err)
		}
		resp, err := client.GetGroup(ctx, &cognitoidentityprovider.GetGroupInput{
			GroupName:  aws.String(ugGroupName),
			UserPoolId: aws.String(userPoolID),
		})
		if err != nil {
			return fmt.Errorf("GetGroup after update failed: %v", err)
		}
		if resp.Group.Description == nil || *resp.Group.Description != "Updated description" {
			return fmt.Errorf("description not updated: got %v", resp.Group.Description)
		}
		if resp.Group.Precedence == nil || *resp.Group.Precedence != 10 {
			return fmt.Errorf("Precedence not set: got %v, want 10", resp.Group.Precedence)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "AdminAddUserToGroup", func() error {
		ugUser := fmt.Sprintf("ug-user-%d", time.Now().UnixNano())
		ugGroup := fmt.Sprintf("ug-group-%d", time.Now().UnixNano())
		_, err := client.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
			UserPoolId:        aws.String(userPoolID),
			Username:          aws.String(ugUser),
			TemporaryPassword: aws.String("TempPass123!"),
			MessageAction:     types.MessageActionTypeSuppress,
		})
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		defer client.AdminDeleteUser(ctx, &cognitoidentityprovider.AdminDeleteUserInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(ugUser),
		})
		_, err = client.CreateGroup(ctx, &cognitoidentityprovider.CreateGroupInput{
			GroupName:  aws.String(ugGroup),
			UserPoolId: aws.String(userPoolID),
		})
		if err != nil {
			return fmt.Errorf("create group: %v", err)
		}
		defer client.DeleteGroup(ctx, &cognitoidentityprovider.DeleteGroupInput{
			GroupName:  aws.String(ugGroup),
			UserPoolId: aws.String(userPoolID),
		})
		_, err = client.AdminAddUserToGroup(ctx, &cognitoidentityprovider.AdminAddUserToGroupInput{
			UserPoolId: aws.String(userPoolID),
			GroupName:  aws.String(ugGroup),
			Username:   aws.String(ugUser),
		})
		if err != nil {
			return fmt.Errorf("AdminAddUserToGroup failed: %v", err)
		}
		listResp, err := client.ListUsersInGroup(ctx, &cognitoidentityprovider.ListUsersInGroupInput{
			UserPoolId: aws.String(userPoolID),
			GroupName:  aws.String(ugGroup),
		})
		if err != nil {
			return fmt.Errorf("ListUsersInGroup: %v", err)
		}
		found := false
		for _, u := range listResp.Users {
			if u.Username != nil && *u.Username == ugUser {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("user not found in group after AdminAddUserToGroup")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "ListUsersInGroup", func() error {
		ugUser2 := fmt.Sprintf("ug2-user-%d", time.Now().UnixNano())
		ugGroup2 := fmt.Sprintf("ug2-group-%d", time.Now().UnixNano())
		_, err := client.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
			UserPoolId:        aws.String(userPoolID),
			Username:          aws.String(ugUser2),
			TemporaryPassword: aws.String("TempPass123!"),
			MessageAction:     types.MessageActionTypeSuppress,
		})
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		defer client.AdminDeleteUser(ctx, &cognitoidentityprovider.AdminDeleteUserInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(ugUser2),
		})
		_, err = client.CreateGroup(ctx, &cognitoidentityprovider.CreateGroupInput{
			GroupName:  aws.String(ugGroup2),
			UserPoolId: aws.String(userPoolID),
		})
		if err != nil {
			return fmt.Errorf("create group: %v", err)
		}
		defer client.DeleteGroup(ctx, &cognitoidentityprovider.DeleteGroupInput{
			GroupName:  aws.String(ugGroup2),
			UserPoolId: aws.String(userPoolID),
		})
		_, err = client.AdminAddUserToGroup(ctx, &cognitoidentityprovider.AdminAddUserToGroupInput{
			UserPoolId: aws.String(userPoolID),
			GroupName:  aws.String(ugGroup2),
			Username:   aws.String(ugUser2),
		})
		if err != nil {
			return fmt.Errorf("AdminAddUserToGroup failed: %v", err)
		}
		listResp, err := client.ListUsersInGroup(ctx, &cognitoidentityprovider.ListUsersInGroupInput{
			UserPoolId: aws.String(userPoolID),
			GroupName:  aws.String(ugGroup2),
		})
		if err != nil {
			return fmt.Errorf("ListUsersInGroup failed: %v", err)
		}
		found := false
		for _, u := range listResp.Users {
			if u.Username != nil && *u.Username == ugUser2 {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("user not found in ListUsersInGroup")
		}
		_, err = client.AdminRemoveUserFromGroup(ctx, &cognitoidentityprovider.AdminRemoveUserFromGroupInput{
			UserPoolId: aws.String(userPoolID),
			GroupName:  aws.String(ugGroup2),
			Username:   aws.String(ugUser2),
		})
		if err != nil {
			return fmt.Errorf("AdminRemoveUserFromGroup failed: %v", err)
		}
		listResp2, err := client.ListUsersInGroup(ctx, &cognitoidentityprovider.ListUsersInGroupInput{
			UserPoolId: aws.String(userPoolID),
			GroupName:  aws.String(ugGroup2),
		})
		if err != nil {
			return fmt.Errorf("ListUsersInGroup after remove failed: %v", err)
		}
		for _, u := range listResp2.Users {
			if u.Username != nil && *u.Username == ugUser2 {
				return fmt.Errorf("user still in group after AdminRemoveUserFromGroup")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "AdminListGroupsForUser", func() error {
		lgUser := fmt.Sprintf("lg-user-%d", time.Now().UnixNano())
		lgGroup1 := fmt.Sprintf("lg-group1-%d", time.Now().UnixNano())
		lgGroup2 := fmt.Sprintf("lg-group2-%d", time.Now().UnixNano())
		_, err := client.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
			UserPoolId:        aws.String(userPoolID),
			Username:          aws.String(lgUser),
			TemporaryPassword: aws.String("TempPass123!"),
			MessageAction:     types.MessageActionTypeSuppress,
		})
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		defer client.AdminDeleteUser(ctx, &cognitoidentityprovider.AdminDeleteUserInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(lgUser),
		})
		_, err = client.CreateGroup(ctx, &cognitoidentityprovider.CreateGroupInput{
			GroupName:  aws.String(lgGroup1),
			UserPoolId: aws.String(userPoolID),
		})
		if err != nil {
			return fmt.Errorf("create group1: %v", err)
		}
		defer client.DeleteGroup(ctx, &cognitoidentityprovider.DeleteGroupInput{
			GroupName:  aws.String(lgGroup1),
			UserPoolId: aws.String(userPoolID),
		})
		_, err = client.CreateGroup(ctx, &cognitoidentityprovider.CreateGroupInput{
			GroupName:  aws.String(lgGroup2),
			UserPoolId: aws.String(userPoolID),
		})
		if err != nil {
			return fmt.Errorf("create group2: %v", err)
		}
		defer client.DeleteGroup(ctx, &cognitoidentityprovider.DeleteGroupInput{
			GroupName:  aws.String(lgGroup2),
			UserPoolId: aws.String(userPoolID),
		})
		_, err = client.AdminAddUserToGroup(ctx, &cognitoidentityprovider.AdminAddUserToGroupInput{
			UserPoolId: aws.String(userPoolID),
			GroupName:  aws.String(lgGroup1),
			Username:   aws.String(lgUser),
		})
		if err != nil {
			return fmt.Errorf("add to group1: %v", err)
		}
		_, err = client.AdminAddUserToGroup(ctx, &cognitoidentityprovider.AdminAddUserToGroupInput{
			UserPoolId: aws.String(userPoolID),
			GroupName:  aws.String(lgGroup2),
			Username:   aws.String(lgUser),
		})
		if err != nil {
			return fmt.Errorf("add to group2: %v", err)
		}
		resp, err := client.AdminListGroupsForUser(ctx, &cognitoidentityprovider.AdminListGroupsForUserInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(lgUser),
		})
		if err != nil {
			return fmt.Errorf("AdminListGroupsForUser failed: %v", err)
		}
		if len(resp.Groups) < 2 {
			return fmt.Errorf("expected at least 2 groups, got %d", len(resp.Groups))
		}
		foundG1, foundG2 := false, false
		for _, g := range resp.Groups {
			if g.GroupName != nil {
				if *g.GroupName == lgGroup1 {
					foundG1 = true
				}
				if *g.GroupName == lgGroup2 {
					foundG2 = true
				}
			}
		}
		if !foundG1 || !foundG2 {
			return fmt.Errorf("expected both groups in AdminListGroupsForUser, found g1=%v g2=%v", foundG1, foundG2)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "DeleteGroup", func() error {
		_, err := client.DeleteGroup(ctx, &cognitoidentityprovider.DeleteGroupInput{
			GroupName:  aws.String(groupName),
			UserPoolId: aws.String(userPoolID),
		})
		if err != nil {
			return err
		}
		_, err = client.GetGroup(ctx, &cognitoidentityprovider.GetGroupInput{
			GroupName:  aws.String(groupName),
			UserPoolId: aws.String(userPoolID),
		})
		if err == nil {
			return fmt.Errorf("expected error getting deleted group")
		}
		return nil
	}))

	return results
}

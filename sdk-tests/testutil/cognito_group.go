package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

func (r *TestRunner) cognitoGroupTests(tc *cognitoIDPContext) []TestResult {
	var results []TestResult

	groupName := tc.unique("group")
	results = append(results, r.RunTest("cognito", "CreateGroup", func() error {
		resp, err := tc.client.CreateGroup(tc.ctx, &cognitoidentityprovider.CreateGroupInput{
			GroupName:  aws.String(groupName),
			UserPoolId: aws.String(tc.userPoolID),
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
		if resp.Group.UserPoolId == nil || *resp.Group.UserPoolId != tc.userPoolID {
			return fmt.Errorf("UserPoolId mismatch: got %v, want %s", resp.Group.UserPoolId, tc.userPoolID)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "ListGroups", func() error {
		resp, err := tc.client.ListGroups(tc.ctx, &cognitoidentityprovider.ListGroupsInput{
			UserPoolId: aws.String(tc.userPoolID),
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
				if g.UserPoolId == nil || *g.UserPoolId != tc.userPoolID {
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
		newGroupName := tc.unique("get-group")
		cleanupGroup, err := tc.createGroup(newGroupName)
		if err != nil {
			return err
		}
		defer cleanupGroup()
		resp, err := tc.client.GetGroup(tc.ctx, &cognitoidentityprovider.GetGroupInput{
			GroupName:  aws.String(newGroupName),
			UserPoolId: aws.String(tc.userPoolID),
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
		if resp.Group.UserPoolId == nil || *resp.Group.UserPoolId != tc.userPoolID {
			return fmt.Errorf("UserPoolId mismatch: got %v", resp.Group.UserPoolId)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "UpdateGroup", func() error {
		ugGroupName := tc.unique("ug-group")
		cleanupUgGroup, err := tc.createGroup(ugGroupName,
			func(input *cognitoidentityprovider.CreateGroupInput) {
				input.Description = aws.String("Original description")
			})
		if err != nil {
			return err
		}
		defer cleanupUgGroup()
		_, err = tc.client.UpdateGroup(tc.ctx, &cognitoidentityprovider.UpdateGroupInput{
			GroupName:   aws.String(ugGroupName),
			UserPoolId:  aws.String(tc.userPoolID),
			Description: aws.String("Updated description"),
			Precedence:  aws.Int32(10),
		})
		if err != nil {
			return fmt.Errorf("UpdateGroup failed: %v", err)
		}
		resp, err := tc.client.GetGroup(tc.ctx, &cognitoidentityprovider.GetGroupInput{
			GroupName:  aws.String(ugGroupName),
			UserPoolId: aws.String(tc.userPoolID),
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
		ugUser := tc.unique("ug-user")
		ugGroup := tc.unique("aug-group")
		cleanupUgUser, err := tc.adminCreateUser(ugUser)
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		defer cleanupUgUser()
		cleanupUgGroup, err := tc.createGroup(ugGroup)
		if err != nil {
			return err
		}
		defer cleanupUgGroup()
		_, err = tc.client.AdminAddUserToGroup(tc.ctx, &cognitoidentityprovider.AdminAddUserToGroupInput{
			UserPoolId: aws.String(tc.userPoolID),
			GroupName:  aws.String(ugGroup),
			Username:   aws.String(ugUser),
		})
		if err != nil {
			return fmt.Errorf("AdminAddUserToGroup failed: %v", err)
		}
		listResp, err := tc.client.ListUsersInGroup(tc.ctx, &cognitoidentityprovider.ListUsersInGroupInput{
			UserPoolId: aws.String(tc.userPoolID),
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
		ugUser2 := tc.unique("ug2-user")
		ugGroup2 := tc.unique("ug2-group")
		cleanupUgUser2, err := tc.adminCreateUser(ugUser2)
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		defer cleanupUgUser2()
		cleanupUgGroup2, err := tc.createGroup(ugGroup2)
		if err != nil {
			return err
		}
		defer cleanupUgGroup2()
		_, err = tc.client.AdminAddUserToGroup(tc.ctx, &cognitoidentityprovider.AdminAddUserToGroupInput{
			UserPoolId: aws.String(tc.userPoolID),
			GroupName:  aws.String(ugGroup2),
			Username:   aws.String(ugUser2),
		})
		if err != nil {
			return fmt.Errorf("AdminAddUserToGroup failed: %v", err)
		}
		listResp, err := tc.client.ListUsersInGroup(tc.ctx, &cognitoidentityprovider.ListUsersInGroupInput{
			UserPoolId: aws.String(tc.userPoolID),
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
		_, err = tc.client.AdminRemoveUserFromGroup(tc.ctx, &cognitoidentityprovider.AdminRemoveUserFromGroupInput{
			UserPoolId: aws.String(tc.userPoolID),
			GroupName:  aws.String(ugGroup2),
			Username:   aws.String(ugUser2),
		})
		if err != nil {
			return fmt.Errorf("AdminRemoveUserFromGroup failed: %v", err)
		}
		listResp2, err := tc.client.ListUsersInGroup(tc.ctx, &cognitoidentityprovider.ListUsersInGroupInput{
			UserPoolId: aws.String(tc.userPoolID),
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
		lgUser := tc.unique("lg-user")
		lgGroup1 := tc.unique("lg-group1")
		lgGroup2 := tc.unique("lg-group2")
		cleanupLgUser, err := tc.adminCreateUser(lgUser)
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		defer cleanupLgUser()
		cleanupLgGroup1, err := tc.createGroup(lgGroup1)
		if err != nil {
			return err
		}
		defer cleanupLgGroup1()
		cleanupLgGroup2, err := tc.createGroup(lgGroup2)
		if err != nil {
			return err
		}
		defer cleanupLgGroup2()
		_, err = tc.client.AdminAddUserToGroup(tc.ctx, &cognitoidentityprovider.AdminAddUserToGroupInput{
			UserPoolId: aws.String(tc.userPoolID),
			GroupName:  aws.String(lgGroup1),
			Username:   aws.String(lgUser),
		})
		if err != nil {
			return fmt.Errorf("add to group1: %v", err)
		}
		_, err = tc.client.AdminAddUserToGroup(tc.ctx, &cognitoidentityprovider.AdminAddUserToGroupInput{
			UserPoolId: aws.String(tc.userPoolID),
			GroupName:  aws.String(lgGroup2),
			Username:   aws.String(lgUser),
		})
		if err != nil {
			return fmt.Errorf("add to group2: %v", err)
		}
		resp, err := tc.client.AdminListGroupsForUser(tc.ctx, &cognitoidentityprovider.AdminListGroupsForUserInput{
			UserPoolId: aws.String(tc.userPoolID),
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
		_, err := tc.client.DeleteGroup(tc.ctx, &cognitoidentityprovider.DeleteGroupInput{
			GroupName:  aws.String(groupName),
			UserPoolId: aws.String(tc.userPoolID),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.GetGroup(tc.ctx, &cognitoidentityprovider.GetGroupInput{
			GroupName:  aws.String(groupName),
			UserPoolId: aws.String(tc.userPoolID),
		})
		if err == nil {
			return fmt.Errorf("expected error getting deleted group")
		}
		return nil
	}))

	return results
}

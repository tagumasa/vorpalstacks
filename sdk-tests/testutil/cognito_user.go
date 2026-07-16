package testutil

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func (r *TestRunner) cognitoUserTests(ctx context.Context, client *cognitoidentityprovider.Client, userPoolID string) []TestResult {
	var results []TestResult

	username := fmt.Sprintf("user-%d", time.Now().UnixNano())
	results = append(results, r.RunTest("cognito", "AdminCreateUser", func() error {
		resp, err := client.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
			UserPoolId:        aws.String(userPoolID),
			Username:          aws.String(username),
			TemporaryPassword: aws.String("TempPass123!"),
			MessageAction:     types.MessageActionTypeSuppress,
		})
		if err != nil {
			return err
		}
		if resp.User == nil {
			return fmt.Errorf("user is nil")
		}
		if resp.User.Username == nil || *resp.User.Username != username {
			return fmt.Errorf("username mismatch: got %v, want %s", resp.User.Username, username)
		}
		if resp.User.UserStatus != types.UserStatusTypeForceChangePassword {
			return fmt.Errorf("expected UserStatus FORCE_CHANGE_PASSWORD, got %v", resp.User.UserStatus)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "AdminGetUser", func() error {
		resp, err := client.AdminGetUser(ctx, &cognitoidentityprovider.AdminGetUserInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(username),
		})
		if err != nil {
			return err
		}
		if resp.Username == nil || *resp.Username != username {
			return fmt.Errorf("username mismatch: got %v, want %s", resp.Username, username)
		}
		if !resp.Enabled {
			return fmt.Errorf("expected user to be enabled")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "ListUsers", func() error {
		resp, err := client.ListUsers(ctx, &cognitoidentityprovider.ListUsersInput{
			UserPoolId: aws.String(userPoolID),
		})
		if err != nil {
			return err
		}
		if len(resp.Users) == 0 {
			return fmt.Errorf("expected at least one user")
		}
		found := false
		for _, u := range resp.Users {
			if u.Username != nil && *u.Username == username {
				found = true
				if u.UserStatus == "" {
					return fmt.Errorf("UserStatus is empty in listing")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created user %s not found in ListUsers", username)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "AdminDisableUser", func() error {
		_, err := client.AdminDisableUser(ctx, &cognitoidentityprovider.AdminDisableUserInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(username),
		})
		if err != nil {
			return err
		}
		getResp, err := client.AdminGetUser(ctx, &cognitoidentityprovider.AdminGetUserInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(username),
		})
		if err != nil {
			return fmt.Errorf("AdminGetUser after disable: %v", err)
		}
		if getResp.Enabled {
			return fmt.Errorf("expected user to be disabled after AdminDisableUser")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "AdminEnableUser", func() error {
		_, err := client.AdminEnableUser(ctx, &cognitoidentityprovider.AdminEnableUserInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(username),
		})
		if err != nil {
			return err
		}
		getResp, err := client.AdminGetUser(ctx, &cognitoidentityprovider.AdminGetUserInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(username),
		})
		if err != nil {
			return fmt.Errorf("AdminGetUser after enable: %v", err)
		}
		if !getResp.Enabled {
			return fmt.Errorf("expected user to be enabled after AdminEnableUser")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "AdminUpdateUserAttributes", func() error {
		attrUser2 := fmt.Sprintf("attr2-user-%d", time.Now().UnixNano())
		_, err := client.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
			UserPoolId:        aws.String(userPoolID),
			Username:          aws.String(attrUser2),
			TemporaryPassword: aws.String("TempPass123!"),
			MessageAction:     types.MessageActionTypeSuppress,
		})
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		defer client.AdminDeleteUser(ctx, &cognitoidentityprovider.AdminDeleteUserInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(attrUser2),
		})
		_, err = client.AdminUpdateUserAttributes(ctx, &cognitoidentityprovider.AdminUpdateUserAttributesInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(attrUser2),
			UserAttributes: []types.AttributeType{
				{Name: aws.String("email"), Value: aws.String("updated@example.com")},
				{Name: aws.String("phone_number"), Value: aws.String("+441234567890")},
			},
		})
		if err != nil {
			return fmt.Errorf("AdminUpdateUserAttributes failed: %v", err)
		}
		getResp, err := client.AdminGetUser(ctx, &cognitoidentityprovider.AdminGetUserInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(attrUser2),
		})
		if err != nil {
			return fmt.Errorf("AdminGetUser failed: %v", err)
		}
		found := false
		for _, attr := range getResp.UserAttributes {
			if attr.Name != nil && *attr.Name == "email" && attr.Value != nil && *attr.Value == "updated@example.com" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("updated email attribute not found")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "AdminDeleteUserAttributes", func() error {
		daUser := fmt.Sprintf("da-user-%d", time.Now().UnixNano())
		_, err := client.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
			UserPoolId:        aws.String(userPoolID),
			Username:          aws.String(daUser),
			TemporaryPassword: aws.String("TempPass123!"),
			MessageAction:     types.MessageActionTypeSuppress,
			UserAttributes: []types.AttributeType{
				{Name: aws.String("email"), Value: aws.String("da@example.com")},
				{Name: aws.String("name"), Value: aws.String("DA User")},
			},
		})
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		defer client.AdminDeleteUser(ctx, &cognitoidentityprovider.AdminDeleteUserInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(daUser),
		})
		_, err = client.AdminDeleteUserAttributes(ctx, &cognitoidentityprovider.AdminDeleteUserAttributesInput{
			UserPoolId:         aws.String(userPoolID),
			Username:           aws.String(daUser),
			UserAttributeNames: []string{"name"},
		})
		if err != nil {
			return fmt.Errorf("AdminDeleteUserAttributes failed: %v", err)
		}
		getResp, err := client.AdminGetUser(ctx, &cognitoidentityprovider.AdminGetUserInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(daUser),
		})
		if err != nil {
			return fmt.Errorf("AdminGetUser failed: %v", err)
		}
		for _, attr := range getResp.UserAttributes {
			if attr.Name != nil && *attr.Name == "name" {
				return fmt.Errorf("attribute 'name' should have been deleted")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "AdminResetUserPassword", func() error {
		rpUser := fmt.Sprintf("rp-user-%d", time.Now().UnixNano())
		_, err := client.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
			UserPoolId:        aws.String(userPoolID),
			Username:          aws.String(rpUser),
			TemporaryPassword: aws.String("TempPass123!"),
			MessageAction:     types.MessageActionTypeSuppress,
		})
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		defer client.AdminDeleteUser(ctx, &cognitoidentityprovider.AdminDeleteUserInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(rpUser),
		})
		_, err = client.AdminSetUserPassword(ctx, &cognitoidentityprovider.AdminSetUserPasswordInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(rpUser),
			Password:   aws.String("PermPass123!"),
			Permanent:  true,
		})
		if err != nil {
			return fmt.Errorf("AdminSetUserPassword failed: %v", err)
		}
		_, err = client.AdminResetUserPassword(ctx, &cognitoidentityprovider.AdminResetUserPasswordInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(rpUser),
		})
		if err != nil {
			return fmt.Errorf("AdminResetUserPassword failed: %v", err)
		}
		getResp, err := client.AdminGetUser(ctx, &cognitoidentityprovider.AdminGetUserInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(rpUser),
		})
		if err != nil {
			return fmt.Errorf("AdminGetUser failed: %v", err)
		}
		if getResp.UserStatus != types.UserStatusTypeForceChangePassword && getResp.UserStatus != types.UserStatusTypeResetRequired {
			return fmt.Errorf("expected FORCE_CHANGE_PASSWORD or RESET_REQUIRED after reset, got %v", getResp.UserStatus)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "AdminSetUserPassword", func() error {
		spUser := fmt.Sprintf("sp-user-%d", time.Now().UnixNano())
		_, err := client.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
			UserPoolId:        aws.String(userPoolID),
			Username:          aws.String(spUser),
			TemporaryPassword: aws.String("TempPass123!"),
			MessageAction:     types.MessageActionTypeSuppress,
		})
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		defer client.AdminDeleteUser(ctx, &cognitoidentityprovider.AdminDeleteUserInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(spUser),
		})
		_, err = client.AdminSetUserPassword(ctx, &cognitoidentityprovider.AdminSetUserPasswordInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(spUser),
			Password:   aws.String("NewPermPass123!"),
			Permanent:  true,
		})
		if err != nil {
			return fmt.Errorf("AdminSetUserPassword failed: %v", err)
		}
		getResp, err := client.AdminGetUser(ctx, &cognitoidentityprovider.AdminGetUserInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(spUser),
		})
		if err != nil {
			return fmt.Errorf("AdminGetUser failed: %v", err)
		}
		if getResp.UserStatus != types.UserStatusTypeConfirmed {
			return fmt.Errorf("expected CONFIRMED after permanent password, got %v", getResp.UserStatus)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "SignUp", func() error {
		signUpClientName := fmt.Sprintf("signup-client-%d", time.Now().UnixNano())
		signUpClientResp, err := client.CreateUserPoolClient(ctx, &cognitoidentityprovider.CreateUserPoolClientInput{
			UserPoolId: aws.String(userPoolID),
			ClientName: aws.String(signUpClientName),
		})
		if err != nil {
			return fmt.Errorf("create client: %v", err)
		}
		signUpClientID := *signUpClientResp.UserPoolClient.ClientId
		defer client.DeleteUserPoolClient(ctx, &cognitoidentityprovider.DeleteUserPoolClientInput{
			ClientId:   aws.String(signUpClientID),
			UserPoolId: aws.String(userPoolID),
		})
		signUpUser := fmt.Sprintf("signup-user-%d", time.Now().UnixNano())
		resp, err := client.SignUp(ctx, &cognitoidentityprovider.SignUpInput{
			ClientId: aws.String(signUpClientID),
			Username: aws.String(signUpUser),
			Password: aws.String("SignUpPass123!"),
			UserAttributes: []types.AttributeType{
				{Name: aws.String("email"), Value: aws.String("signup@example.com")},
			},
		})
		if err != nil {
			return fmt.Errorf("SignUp failed: %v", err)
		}
		if resp.UserSub == nil || *resp.UserSub == "" {
			return fmt.Errorf("UserSub is nil or empty")
		}
		if resp.UserConfirmed {
			return fmt.Errorf("expected UserConfirmed=false after SignUp")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "ConfirmSignUp", func() error {
		confirmClientName := fmt.Sprintf("confirm-client-%d", time.Now().UnixNano())
		confirmClientResp, err := client.CreateUserPoolClient(ctx, &cognitoidentityprovider.CreateUserPoolClientInput{
			UserPoolId: aws.String(userPoolID),
			ClientName: aws.String(confirmClientName),
		})
		if err != nil {
			return fmt.Errorf("create client: %v", err)
		}
		confirmClientID := *confirmClientResp.UserPoolClient.ClientId
		defer client.DeleteUserPoolClient(ctx, &cognitoidentityprovider.DeleteUserPoolClientInput{
			ClientId:   aws.String(confirmClientID),
			UserPoolId: aws.String(userPoolID),
		})
		confirmUser := fmt.Sprintf("confirm-user-%d", time.Now().UnixNano())
		_, err = client.SignUp(ctx, &cognitoidentityprovider.SignUpInput{
			ClientId: aws.String(confirmClientID),
			Username: aws.String(confirmUser),
			Password: aws.String("ConfirmPass123!"),
		})
		if err != nil {
			return fmt.Errorf("SignUp failed: %v", err)
		}
		_, err = client.ConfirmSignUp(ctx, &cognitoidentityprovider.ConfirmSignUpInput{
			ClientId:         aws.String(confirmClientID),
			Username:         aws.String(confirmUser),
			ConfirmationCode: aws.String("123456"),
		})
		if err != nil {
			return fmt.Errorf("ConfirmSignUp failed: %v", err)
		}
		getResp, err := client.AdminGetUser(ctx, &cognitoidentityprovider.AdminGetUserInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(confirmUser),
		})
		if err != nil {
			return fmt.Errorf("AdminGetUser after confirm: %v", err)
		}
		if getResp.UserStatus != types.UserStatusTypeConfirmed {
			return fmt.Errorf("expected CONFIRMED after ConfirmSignUp, got %v", getResp.UserStatus)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "AdminInitiateAuth", func() error {
		authUser := fmt.Sprintf("auth-user-%d", time.Now().UnixNano())
		_, err := client.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
			UserPoolId:        aws.String(userPoolID),
			Username:          aws.String(authUser),
			TemporaryPassword: aws.String("TempPass123!"),
			MessageAction:     types.MessageActionTypeSuppress,
		})
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		defer client.AdminDeleteUser(ctx, &cognitoidentityprovider.AdminDeleteUserInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(authUser),
		})
		_, err = client.AdminSetUserPassword(ctx, &cognitoidentityprovider.AdminSetUserPasswordInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(authUser),
			Password:   aws.String("AuthPass123!"),
			Permanent:  true,
		})
		if err != nil {
			return fmt.Errorf("set password: %v", err)
		}
		authClientName := fmt.Sprintf("auth-client-%d", time.Now().UnixNano())
		authClientResp, err := client.CreateUserPoolClient(ctx, &cognitoidentityprovider.CreateUserPoolClientInput{
			UserPoolId: aws.String(userPoolID),
			ClientName: aws.String(authClientName),
		})
		if err != nil {
			return fmt.Errorf("create auth client: %v", err)
		}
		authClientID := *authClientResp.UserPoolClient.ClientId
		defer client.DeleteUserPoolClient(ctx, &cognitoidentityprovider.DeleteUserPoolClientInput{
			ClientId:   aws.String(authClientID),
			UserPoolId: aws.String(userPoolID),
		})
		authResp, err := client.AdminInitiateAuth(ctx, &cognitoidentityprovider.AdminInitiateAuthInput{
			UserPoolId: aws.String(userPoolID),
			ClientId:   aws.String(authClientID),
			AuthFlow:   types.AuthFlowTypeAdminNoSrpAuth,
			AuthParameters: map[string]string{
				"USERNAME": authUser,
				"PASSWORD": "AuthPass123!",
			},
		})
		if err != nil {
			return fmt.Errorf("AdminInitiateAuth failed: %v", err)
		}
		if authResp.AuthenticationResult == nil {
			return fmt.Errorf("AuthenticationResult is nil")
		}
		if authResp.AuthenticationResult.AccessToken == nil || *authResp.AuthenticationResult.AccessToken == "" {
			return fmt.Errorf("AccessToken is nil or empty")
		}
		if authResp.AuthenticationResult.IdToken == nil || *authResp.AuthenticationResult.IdToken == "" {
			return fmt.Errorf("IdToken is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "AdminUserGlobalSignOut", func() error {
		gsoUser := fmt.Sprintf("gso-user-%d", time.Now().UnixNano())
		_, err := client.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
			UserPoolId:        aws.String(userPoolID),
			Username:          aws.String(gsoUser),
			TemporaryPassword: aws.String("TempPass123!"),
			MessageAction:     types.MessageActionTypeSuppress,
		})
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		defer client.AdminDeleteUser(ctx, &cognitoidentityprovider.AdminDeleteUserInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(gsoUser),
		})
		_, err = client.AdminSetUserPassword(ctx, &cognitoidentityprovider.AdminSetUserPasswordInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(gsoUser),
			Password:   aws.String("GSOPass123!"),
			Permanent:  true,
		})
		if err != nil {
			return fmt.Errorf("set password: %v", err)
		}
		gsoClientName := fmt.Sprintf("gso-client-%d", time.Now().UnixNano())
		gsoClientResp, err := client.CreateUserPoolClient(ctx, &cognitoidentityprovider.CreateUserPoolClientInput{
			UserPoolId: aws.String(userPoolID),
			ClientName: aws.String(gsoClientName),
		})
		if err != nil {
			return fmt.Errorf("create client: %v", err)
		}
		gsoClientID := *gsoClientResp.UserPoolClient.ClientId
		defer client.DeleteUserPoolClient(ctx, &cognitoidentityprovider.DeleteUserPoolClientInput{
			ClientId:   aws.String(gsoClientID),
			UserPoolId: aws.String(userPoolID),
		})
		authResp, err := client.AdminInitiateAuth(ctx, &cognitoidentityprovider.AdminInitiateAuthInput{
			UserPoolId: aws.String(userPoolID),
			ClientId:   aws.String(gsoClientID),
			AuthFlow:   types.AuthFlowTypeAdminNoSrpAuth,
			AuthParameters: map[string]string{
				"USERNAME": gsoUser,
				"PASSWORD": "GSOPass123!",
			},
		})
		if err != nil {
			return fmt.Errorf("AdminInitiateAuth: %v", err)
		}
		if authResp.AuthenticationResult == nil || authResp.AuthenticationResult.AccessToken == nil {
			return fmt.Errorf("AccessToken is nil before sign-out")
		}
		accessToken := *authResp.AuthenticationResult.AccessToken

		_, err = client.AdminUserGlobalSignOut(ctx, &cognitoidentityprovider.AdminUserGlobalSignOutInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(gsoUser),
		})
		if err != nil {
			return fmt.Errorf("AdminUserGlobalSignOut failed: %v", err)
		}
		_, err = client.GetUser(ctx, &cognitoidentityprovider.GetUserInput{
			AccessToken: aws.String(accessToken),
		})
		if err == nil {
			return fmt.Errorf("expected error using access token after global sign-out")
		}
		authResp2, err := client.AdminInitiateAuth(ctx, &cognitoidentityprovider.AdminInitiateAuthInput{
			UserPoolId: aws.String(userPoolID),
			ClientId:   aws.String(gsoClientID),
			AuthFlow:   types.AuthFlowTypeAdminNoSrpAuth,
			AuthParameters: map[string]string{
				"USERNAME": gsoUser,
				"PASSWORD": "GSOPass123!",
			},
		})
		if err != nil {
			return fmt.Errorf("re-auth after global sign-out failed: %v", err)
		}
		if authResp2.AuthenticationResult == nil || authResp2.AuthenticationResult.AccessToken == nil {
			return fmt.Errorf("AccessToken is nil after re-auth")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "GlobalSignOut", func() error {
		_, err := client.GlobalSignOut(ctx, &cognitoidentityprovider.GlobalSignOutInput{
			AccessToken: aws.String("dummy-token"),
		})
		if err == nil {
			return fmt.Errorf("expected error for dummy access token")
		}
		var notAuthEx *types.NotAuthorizedException
		if !errors.As(err, &notAuthEx) {
			return fmt.Errorf("expected NotAuthorizedException for invalid token, got: %v", err)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "AdminDeleteUser", func() error {
		_, err := client.AdminDeleteUser(ctx, &cognitoidentityprovider.AdminDeleteUserInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(username),
		})
		if err != nil {
			return err
		}
		_, err = client.AdminGetUser(ctx, &cognitoidentityprovider.AdminGetUserInput{
			UserPoolId: aws.String(userPoolID),
			Username:   aws.String(username),
		})
		if err == nil {
			return fmt.Errorf("expected error getting deleted user")
		}
		var notFoundEx *types.UserNotFoundException
		if !errors.As(err, &notFoundEx) {
			return fmt.Errorf("expected UserNotFoundException, got: %v", err)
		}
		return nil
	}))

	return results
}

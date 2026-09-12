package testutil

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func (r *TestRunner) cognitoUserTests(tc *cognitoIDPContext) []TestResult {
	var results []TestResult

	username := tc.unique("user")
	results = append(results, r.RunTest("cognito", "AdminCreateUser", func() error {
		resp, err := tc.client.AdminCreateUser(tc.ctx, &cognitoidentityprovider.AdminCreateUserInput{
			UserPoolId:        aws.String(tc.userPoolID),
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
		resp, err := tc.client.AdminGetUser(tc.ctx, &cognitoidentityprovider.AdminGetUserInput{
			UserPoolId: aws.String(tc.userPoolID),
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
		sub := ""
		for _, attr := range resp.UserAttributes {
			if attr.Name != nil && *attr.Name == "sub" {
				sub = aws.ToString(attr.Value)
			}
		}
		if sub == "" {
			return fmt.Errorf("AdminGetUser response carries no sub attribute")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "ListUsers", func() error {
		resp, err := tc.client.ListUsers(tc.ctx, &cognitoidentityprovider.ListUsersInput{
			UserPoolId: aws.String(tc.userPoolID),
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

	// The sub attribute is a filterable attribute of every user: the filter
	// resolves the created user by its subject identifier alone.
	results = append(results, r.RunTest("cognito", "ListUsers_FilterSub", func() error {
		getResp, err := tc.client.AdminGetUser(tc.ctx, &cognitoidentityprovider.AdminGetUserInput{
			UserPoolId: aws.String(tc.userPoolID),
			Username:   aws.String(username),
		})
		if err != nil {
			return fmt.Errorf("AdminGetUser: %v", err)
		}
		sub := ""
		for _, attr := range getResp.UserAttributes {
			if attr.Name != nil && *attr.Name == "sub" {
				sub = aws.ToString(attr.Value)
			}
		}
		if sub == "" {
			return fmt.Errorf("AdminGetUser carries no sub to filter on")
		}
		resp, err := tc.client.ListUsers(tc.ctx, &cognitoidentityprovider.ListUsersInput{
			UserPoolId: aws.String(tc.userPoolID),
			Filter:     aws.String(`sub = "` + sub + `"`),
		})
		if err != nil {
			return err
		}
		found := false
		for _, u := range resp.Users {
			if u.Username != nil && *u.Username == username {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("sub filter did not resolve user %s (%d users returned)", username, len(resp.Users))
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "AdminDisableUser", func() error {
		_, err := tc.client.AdminDisableUser(tc.ctx, &cognitoidentityprovider.AdminDisableUserInput{
			UserPoolId: aws.String(tc.userPoolID),
			Username:   aws.String(username),
		})
		if err != nil {
			return err
		}
		getResp, err := tc.client.AdminGetUser(tc.ctx, &cognitoidentityprovider.AdminGetUserInput{
			UserPoolId: aws.String(tc.userPoolID),
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
		_, err := tc.client.AdminEnableUser(tc.ctx, &cognitoidentityprovider.AdminEnableUserInput{
			UserPoolId: aws.String(tc.userPoolID),
			Username:   aws.String(username),
		})
		if err != nil {
			return err
		}
		getResp, err := tc.client.AdminGetUser(tc.ctx, &cognitoidentityprovider.AdminGetUserInput{
			UserPoolId: aws.String(tc.userPoolID),
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
		attrUser2 := tc.unique("attr2-user")
		cleanupAttrUser2, err := tc.adminCreateUser(attrUser2)
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		defer cleanupAttrUser2()
		_, err = tc.client.AdminUpdateUserAttributes(tc.ctx, &cognitoidentityprovider.AdminUpdateUserAttributesInput{
			UserPoolId: aws.String(tc.userPoolID),
			Username:   aws.String(attrUser2),
			UserAttributes: []types.AttributeType{
				{Name: aws.String("email"), Value: aws.String("updated@example.com")},
				{Name: aws.String("phone_number"), Value: aws.String("+441234567890")},
			},
		})
		if err != nil {
			return fmt.Errorf("AdminUpdateUserAttributes failed: %v", err)
		}
		getResp, err := tc.client.AdminGetUser(tc.ctx, &cognitoidentityprovider.AdminGetUserInput{
			UserPoolId: aws.String(tc.userPoolID),
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

	// Attribute writes are validated against the pool schema: sub is
	// immutable, and an attribute the pool schema does not define cannot be
	// written at all.
	results = append(results, r.RunTest("cognito", "AdminUpdateUserAttributes_SchemaRejected", func() error {
		schemaUser := tc.unique("schema-user")
		cleanupSchemaUser, err := tc.adminCreateUser(schemaUser)
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		defer cleanupSchemaUser()
		for _, c := range []struct {
			name  string
			value string
		}{
			{"sub", "forged-subject"},
			{"custom:undefined", "x"},
		} {
			_, err := tc.client.AdminUpdateUserAttributes(tc.ctx, &cognitoidentityprovider.AdminUpdateUserAttributesInput{
				UserPoolId: aws.String(tc.userPoolID),
				Username:   aws.String(schemaUser),
				UserAttributes: []types.AttributeType{
					{Name: aws.String(c.name), Value: aws.String(c.value)},
				},
			})
			if err := expectAWSErrorCode(err, "InvalidParameterException"); err != nil {
				return fmt.Errorf("updating %s: %v", c.name, err)
			}
		}
		getResp, err := tc.client.AdminGetUser(tc.ctx, &cognitoidentityprovider.AdminGetUserInput{
			UserPoolId: aws.String(tc.userPoolID),
			Username:   aws.String(schemaUser),
		})
		if err != nil {
			return fmt.Errorf("AdminGetUser after rejections: %v", err)
		}
		for _, attr := range getResp.UserAttributes {
			if attr.Name != nil && *attr.Name == "sub" && aws.ToString(attr.Value) == "forged-subject" {
				return fmt.Errorf("rejected sub update was persisted")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "AdminDeleteUserAttributes", func() error {
		daUser := tc.unique("da-user")
		_, err := tc.client.AdminCreateUser(tc.ctx, &cognitoidentityprovider.AdminCreateUserInput{
			UserPoolId:        aws.String(tc.userPoolID),
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
		defer tc.client.AdminDeleteUser(tc.ctx, &cognitoidentityprovider.AdminDeleteUserInput{
			UserPoolId: aws.String(tc.userPoolID),
			Username:   aws.String(daUser),
		})
		_, err = tc.client.AdminDeleteUserAttributes(tc.ctx, &cognitoidentityprovider.AdminDeleteUserAttributesInput{
			UserPoolId:         aws.String(tc.userPoolID),
			Username:           aws.String(daUser),
			UserAttributeNames: []string{"name"},
		})
		if err != nil {
			return fmt.Errorf("AdminDeleteUserAttributes failed: %v", err)
		}
		getResp, err := tc.client.AdminGetUser(tc.ctx, &cognitoidentityprovider.AdminGetUserInput{
			UserPoolId: aws.String(tc.userPoolID),
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
		rpUser := tc.unique("rp-user")
		cleanupRpUser, err := tc.createConfirmedUser(rpUser, "PermPass123!")
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		defer cleanupRpUser()
		_, err = tc.client.AdminResetUserPassword(tc.ctx, &cognitoidentityprovider.AdminResetUserPasswordInput{
			UserPoolId: aws.String(tc.userPoolID),
			Username:   aws.String(rpUser),
		})
		if err != nil {
			return fmt.Errorf("AdminResetUserPassword failed: %v", err)
		}
		getResp, err := tc.client.AdminGetUser(tc.ctx, &cognitoidentityprovider.AdminGetUserInput{
			UserPoolId: aws.String(tc.userPoolID),
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
		spUser := tc.unique("sp-user")
		cleanupSpUser, err := tc.adminCreateUser(spUser)
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		defer cleanupSpUser()
		_, err = tc.client.AdminSetUserPassword(tc.ctx, &cognitoidentityprovider.AdminSetUserPasswordInput{
			UserPoolId: aws.String(tc.userPoolID),
			Username:   aws.String(spUser),
			Password:   aws.String("NewPermPass123!"),
			Permanent:  true,
		})
		if err != nil {
			return fmt.Errorf("AdminSetUserPassword failed: %v", err)
		}
		getResp, err := tc.client.AdminGetUser(tc.ctx, &cognitoidentityprovider.AdminGetUserInput{
			UserPoolId: aws.String(tc.userPoolID),
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
		signUpClientID, cleanupSignUpClientID, err := tc.createPoolClient(tc.userPoolID, tc.unique("signup-client"))
		if err != nil {
			return fmt.Errorf("create client: %v", err)
		}
		defer cleanupSignUpClientID()
		signUpUser := tc.unique("signup-user")
		resp, err := tc.client.SignUp(tc.ctx, &cognitoidentityprovider.SignUpInput{
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
		confirmClientID, cleanupConfirmClientID, err := tc.createPoolClient(tc.userPoolID, tc.unique("confirm-client"))
		if err != nil {
			return fmt.Errorf("create client: %v", err)
		}
		defer cleanupConfirmClientID()
		confirmUser := tc.unique("confirm-user")
		_, err = tc.client.SignUp(tc.ctx, &cognitoidentityprovider.SignUpInput{
			ClientId: aws.String(confirmClientID),
			Username: aws.String(confirmUser),
			Password: aws.String("ConfirmPass123!"),
		})
		if err != nil {
			return fmt.Errorf("SignUp failed: %v", err)
		}
		_, err = tc.client.ConfirmSignUp(tc.ctx, &cognitoidentityprovider.ConfirmSignUpInput{
			ClientId:         aws.String(confirmClientID),
			Username:         aws.String(confirmUser),
			ConfirmationCode: aws.String("123456"),
		})
		if err != nil {
			return fmt.Errorf("ConfirmSignUp failed: %v", err)
		}
		getResp, err := tc.client.AdminGetUser(tc.ctx, &cognitoidentityprovider.AdminGetUserInput{
			UserPoolId: aws.String(tc.userPoolID),
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
		authUser := tc.unique("auth-user")
		cleanupAuthUser, err := tc.createConfirmedUser(authUser, "AuthPass123!")
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		defer cleanupAuthUser()
		authClientID, cleanupAuthClientID, err := tc.createPoolClient(tc.userPoolID, tc.unique("auth-client"))
		if err != nil {
			return fmt.Errorf("create auth client: %v", err)
		}
		defer cleanupAuthClientID()
		authResp, err := tc.client.AdminInitiateAuth(tc.ctx, &cognitoidentityprovider.AdminInitiateAuthInput{
			UserPoolId: aws.String(tc.userPoolID),
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

	// GetUser returns the authenticated caller's profile: the model's
	// response members Username and UserAttributes, with the required sub
	// attribute present.
	results = append(results, r.RunTest("cognito", "GetUser", func() error {
		getUser := tc.unique("getuser")
		cleanupGetUser, err := tc.createConfirmedUser(getUser, "GetUserPass123!")
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		defer cleanupGetUser()
		getUserClientID, cleanupGetUserClient, err := tc.createPoolClient(tc.userPoolID, tc.unique("getuser-client"))
		if err != nil {
			return fmt.Errorf("create client: %v", err)
		}
		defer cleanupGetUserClient()
		authResp, err := tc.client.AdminInitiateAuth(tc.ctx, &cognitoidentityprovider.AdminInitiateAuthInput{
			UserPoolId: aws.String(tc.userPoolID),
			ClientId:   aws.String(getUserClientID),
			AuthFlow:   types.AuthFlowTypeAdminNoSrpAuth,
			AuthParameters: map[string]string{
				"USERNAME": getUser,
				"PASSWORD": "GetUserPass123!",
			},
		})
		if err != nil {
			return fmt.Errorf("AdminInitiateAuth: %v", err)
		}
		if authResp.AuthenticationResult == nil || authResp.AuthenticationResult.AccessToken == nil {
			return fmt.Errorf("AccessToken is nil")
		}
		resp, err := tc.client.GetUser(tc.ctx, &cognitoidentityprovider.GetUserInput{
			AccessToken: authResp.AuthenticationResult.AccessToken,
		})
		if err != nil {
			return err
		}
		if resp.Username == nil || *resp.Username != getUser {
			return fmt.Errorf("username mismatch: got %v, want %s", resp.Username, getUser)
		}
		sub := ""
		for _, attr := range resp.UserAttributes {
			if attr.Name != nil && *attr.Name == "sub" {
				sub = aws.ToString(attr.Value)
			}
		}
		if sub == "" {
			return fmt.Errorf("GetUser response carries no sub attribute")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "AdminUserGlobalSignOut", func() error {
		gsoUser := tc.unique("gso-user")
		cleanupGsoUser, err := tc.createConfirmedUser(gsoUser, "GSOPass123!")
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		defer cleanupGsoUser()
		gsoClientID, cleanupGsoClientID, err := tc.createPoolClient(tc.userPoolID, tc.unique("gso-client"))
		if err != nil {
			return fmt.Errorf("create client: %v", err)
		}
		defer cleanupGsoClientID()
		authResp, err := tc.client.AdminInitiateAuth(tc.ctx, &cognitoidentityprovider.AdminInitiateAuthInput{
			UserPoolId: aws.String(tc.userPoolID),
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

		_, err = tc.client.AdminUserGlobalSignOut(tc.ctx, &cognitoidentityprovider.AdminUserGlobalSignOutInput{
			UserPoolId: aws.String(tc.userPoolID),
			Username:   aws.String(gsoUser),
		})
		if err != nil {
			return fmt.Errorf("AdminUserGlobalSignOut failed: %v", err)
		}
		_, err = tc.client.GetUser(tc.ctx, &cognitoidentityprovider.GetUserInput{
			AccessToken: aws.String(accessToken),
		})
		if err == nil {
			return fmt.Errorf("expected error using access token after global sign-out")
		}
		authResp2, err := tc.client.AdminInitiateAuth(tc.ctx, &cognitoidentityprovider.AdminInitiateAuthInput{
			UserPoolId: aws.String(tc.userPoolID),
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

	results = append(results, r.RunTest("cognito", "AdminDisableUser_RevokesTokens", func() error {
		disUser := tc.unique("disable-user")
		cleanupDisUser, err := tc.createConfirmedUser(disUser, "DisablePass123!")
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		defer cleanupDisUser()
		disClientID, cleanupDisClientID, err := tc.createPoolClient(tc.userPoolID, tc.unique("disable-client"))
		if err != nil {
			return fmt.Errorf("create client: %v", err)
		}
		defer cleanupDisClientID()
		authResp, err := tc.client.AdminInitiateAuth(tc.ctx, &cognitoidentityprovider.AdminInitiateAuthInput{
			UserPoolId: aws.String(tc.userPoolID),
			ClientId:   aws.String(disClientID),
			AuthFlow:   types.AuthFlowTypeAdminNoSrpAuth,
			AuthParameters: map[string]string{
				"USERNAME": disUser,
				"PASSWORD": "DisablePass123!",
			},
		})
		if err != nil {
			return fmt.Errorf("AdminInitiateAuth: %v", err)
		}
		if authResp.AuthenticationResult == nil || authResp.AuthenticationResult.AccessToken == nil || authResp.AuthenticationResult.RefreshToken == nil {
			return fmt.Errorf("tokens are nil before disable")
		}
		accessToken := *authResp.AuthenticationResult.AccessToken
		refreshToken := *authResp.AuthenticationResult.RefreshToken

		if _, err := tc.client.AdminDisableUser(tc.ctx, &cognitoidentityprovider.AdminDisableUserInput{
			UserPoolId: aws.String(tc.userPoolID),
			Username:   aws.String(disUser),
		}); err != nil {
			return fmt.Errorf("AdminDisableUser: %v", err)
		}

		_, err = tc.client.GetUser(tc.ctx, &cognitoidentityprovider.GetUserInput{
			AccessToken: aws.String(accessToken),
		})
		if err == nil {
			return fmt.Errorf("expected error using access token after disable")
		}
		var notAuthEx *types.NotAuthorizedException
		if !errors.As(err, &notAuthEx) {
			return fmt.Errorf("expected NotAuthorizedException for the pre-disable access token, got: %v", err)
		}

		_, err = tc.client.AdminInitiateAuth(tc.ctx, &cognitoidentityprovider.AdminInitiateAuthInput{
			UserPoolId: aws.String(tc.userPoolID),
			ClientId:   aws.String(disClientID),
			AuthFlow:   types.AuthFlowTypeRefreshTokenAuth,
			AuthParameters: map[string]string{
				"REFRESH_TOKEN": refreshToken,
			},
		})
		if err == nil {
			return fmt.Errorf("expected error refreshing with the pre-disable refresh token")
		}
		if !errors.As(err, &notAuthEx) {
			return fmt.Errorf("expected NotAuthorizedException for the pre-disable refresh token, got: %v", err)
		}

		// Re-enabling restores authentication: the block was the disabled
		// profile, not a destroyed one.
		if _, err := tc.client.AdminEnableUser(tc.ctx, &cognitoidentityprovider.AdminEnableUserInput{
			UserPoolId: aws.String(tc.userPoolID),
			Username:   aws.String(disUser),
		}); err != nil {
			return fmt.Errorf("AdminEnableUser: %v", err)
		}
		reAuthResp, err := tc.client.AdminInitiateAuth(tc.ctx, &cognitoidentityprovider.AdminInitiateAuthInput{
			UserPoolId: aws.String(tc.userPoolID),
			ClientId:   aws.String(disClientID),
			AuthFlow:   types.AuthFlowTypeAdminNoSrpAuth,
			AuthParameters: map[string]string{
				"USERNAME": disUser,
				"PASSWORD": "DisablePass123!",
			},
		})
		if err != nil {
			return fmt.Errorf("sign-in after re-enable failed: %v", err)
		}
		if reAuthResp.AuthenticationResult == nil || reAuthResp.AuthenticationResult.AccessToken == nil {
			return fmt.Errorf("AccessToken is nil after re-enable")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "GlobalSignOut", func() error {
		_, err := tc.client.GlobalSignOut(tc.ctx, &cognitoidentityprovider.GlobalSignOutInput{
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
		_, err := tc.client.AdminDeleteUser(tc.ctx, &cognitoidentityprovider.AdminDeleteUserInput{
			UserPoolId: aws.String(tc.userPoolID),
			Username:   aws.String(username),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.AdminGetUser(tc.ctx, &cognitoidentityprovider.AdminGetUserInput{
			UserPoolId: aws.String(tc.userPoolID),
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

	// A pool with email as an alias attribute implements the staged claim
	// model: duplicate values coexist while unverified, the second
	// confirmation is rejected with AliasExistsException, alias values sign
	// in, and AdminCreateUser's ForceAliasCreation migrates a verified claim
	// to the new user.
	results = append(results, r.RunTest("cognito", "AliasAttributes_ClaimAndSignIn", func() error {
		aliasPoolID, cleanupAliasPool, err := tc.createUserPool(tc.unique("alias-pool"), func(in *cognitoidentityprovider.CreateUserPoolInput) {
			in.AliasAttributes = []types.AliasAttributeType{types.AliasAttributeTypeEmail}
			in.AutoVerifiedAttributes = []types.VerifiedAttributeType{types.VerifiedAttributeTypeEmail}
			in.Policies = &types.UserPoolPolicyType{
				PasswordPolicy: &types.PasswordPolicyType{
					MinimumLength:    aws.Int32(8),
					RequireUppercase: true,
					RequireLowercase: true,
					RequireNumbers:   true,
				},
			}
		})
		if err != nil {
			return fmt.Errorf("create alias pool: %v", err)
		}
		defer cleanupAliasPool()
		aliasClientID, cleanupAliasClient, err := tc.createPoolClient(aliasPoolID, tc.unique("alias-client"))
		if err != nil {
			return fmt.Errorf("create client: %v", err)
		}
		defer cleanupAliasClient()

		shared := tc.unique("alias-shared") + "@example.com"
		holder := tc.unique("alias-holder")
		duplicate := tc.unique("alias-dup")
		signUpWithSharedEmail := func(username string) error {
			_, err := tc.client.SignUp(tc.ctx, &cognitoidentityprovider.SignUpInput{
				ClientId: aws.String(aliasClientID),
				Username: aws.String(username),
				Password: aws.String("AliasPass123!"),
				UserAttributes: []types.AttributeType{
					{Name: aws.String("email"), Value: aws.String(shared)},
				},
			})
			return err
		}
		if err := signUpWithSharedEmail(holder); err != nil {
			return fmt.Errorf("sign-up holder: %v", err)
		}
		if err := signUpWithSharedEmail(duplicate); err != nil {
			return fmt.Errorf("sign-up with a duplicate alias value must succeed: %v", err)
		}

		if _, err := tc.client.ConfirmSignUp(tc.ctx, &cognitoidentityprovider.ConfirmSignUpInput{
			ClientId:         aws.String(aliasClientID),
			Username:         aws.String(holder),
			ConfirmationCode: aws.String("123456"),
		}); err != nil {
			return fmt.Errorf("confirm holder: %v", err)
		}
		_, err = tc.client.ConfirmSignUp(tc.ctx, &cognitoidentityprovider.ConfirmSignUpInput{
			ClientId:         aws.String(aliasClientID),
			Username:         aws.String(duplicate),
			ConfirmationCode: aws.String("123456"),
		})
		if codeErr := expectAWSErrorCode(err, "AliasExistsException"); codeErr != nil {
			return fmt.Errorf("confirming a duplicate alias value: %v", codeErr)
		}

		// Signing in with the alias value reaches the holder's account.
		authResp, err := tc.client.AdminInitiateAuth(tc.ctx, &cognitoidentityprovider.AdminInitiateAuthInput{
			UserPoolId: aws.String(aliasPoolID),
			ClientId:   aws.String(aliasClientID),
			AuthFlow:   types.AuthFlowTypeAdminNoSrpAuth,
			AuthParameters: map[string]string{
				"USERNAME": shared,
				"PASSWORD": "AliasPass123!",
			},
		})
		if err != nil {
			return fmt.Errorf("sign-in with the alias value: %v", err)
		}
		if authResp.AuthenticationResult == nil || authResp.AuthenticationResult.AccessToken == nil {
			return fmt.Errorf("alias sign-in returned no tokens: %+v", authResp)
		}
		profile, err := tc.client.GetUser(tc.ctx, &cognitoidentityprovider.GetUserInput{
			AccessToken: authResp.AuthenticationResult.AccessToken,
		})
		if err != nil {
			return fmt.Errorf("GetUser after alias sign-in: %v", err)
		}
		if aws.ToString(profile.Username) != holder {
			return fmt.Errorf("alias sign-in resolved username %q, want %q", aws.ToString(profile.Username), holder)
		}

		verifiedEmail := []types.AttributeType{
			{Name: aws.String("email"), Value: aws.String(shared)},
			{Name: aws.String("email_verified"), Value: aws.String("true")},
		}
		_, err = tc.client.AdminCreateUser(tc.ctx, &cognitoidentityprovider.AdminCreateUserInput{
			UserPoolId:     aws.String(aliasPoolID),
			Username:       aws.String(tc.unique("alias-challenger")),
			MessageAction:  types.MessageActionTypeSuppress,
			UserAttributes: verifiedEmail,
		})
		if codeErr := expectAWSErrorCode(err, "AliasExistsException"); codeErr != nil {
			return fmt.Errorf("AdminCreateUser claiming a verified alias: %v", codeErr)
		}

		taker := tc.unique("alias-taker")
		if _, err := tc.client.AdminCreateUser(tc.ctx, &cognitoidentityprovider.AdminCreateUserInput{
			UserPoolId:         aws.String(aliasPoolID),
			Username:           aws.String(taker),
			MessageAction:      types.MessageActionTypeSuppress,
			TemporaryPassword:  aws.String("TakerTemp123!"),
			ForceAliasCreation: true,
			UserAttributes:     verifiedEmail,
		}); err != nil {
			return fmt.Errorf("AdminCreateUser with ForceAliasCreation: %v", err)
		}
		holderResp, err := tc.client.AdminGetUser(tc.ctx, &cognitoidentityprovider.AdminGetUserInput{
			UserPoolId: aws.String(aliasPoolID),
			Username:   aws.String(holder),
		})
		if err != nil {
			return fmt.Errorf("AdminGetUser holder after migration: %v", err)
		}
		emailVerified := ""
		for _, attr := range holderResp.UserAttributes {
			if aws.ToString(attr.Name) == "email_verified" {
				emailVerified = aws.ToString(attr.Value)
			}
		}
		if emailVerified != "false" {
			return fmt.Errorf("previous holder email_verified = %q after alias migration, want false", emailVerified)
		}

		// After the migration the alias resolves to the new holder, whose
		// admin-created account answers with the new-password challenge.
		migratedResp, err := tc.client.AdminInitiateAuth(tc.ctx, &cognitoidentityprovider.AdminInitiateAuthInput{
			UserPoolId: aws.String(aliasPoolID),
			ClientId:   aws.String(aliasClientID),
			AuthFlow:   types.AuthFlowTypeAdminNoSrpAuth,
			AuthParameters: map[string]string{
				"USERNAME": shared,
				"PASSWORD": "TakerTemp123!",
			},
		})
		if err != nil {
			return fmt.Errorf("sign-in with the migrated alias: %v", err)
		}
		if migratedResp.ChallengeName != types.ChallengeNameTypeNewPasswordRequired {
			return fmt.Errorf("migrated alias sign-in challenge = %v, want NEW_PASSWORD_REQUIRED", migratedResp.ChallengeName)
		}
		return nil
	}))

	// The event history records the account lifecycle, not sign-in alone: a
	// self-service sign-up leaves a SignUp event with response Pass, and the
	// sign-in that follows leaves a SignIn event on the same user.
	eventsPoolID, eventsPoolCleanup, err := tc.createUserPool(tc.unique("events-pool"))
	if err != nil {
		results = append(results, r.RunTest("cognito", "AdminListUserAuthEvents", func() error {
			return fmt.Errorf("create events pool: %v", err)
		}))
		return results
	}
	defer eventsPoolCleanup()
	eventsClientID, eventsClientCleanup, err := tc.createPoolClient(eventsPoolID, tc.unique("events-client"))
	if err != nil {
		results = append(results, r.RunTest("cognito", "AdminListUserAuthEvents", func() error {
			return fmt.Errorf("create events app client: %v", err)
		}))
		return results
	}
	defer eventsClientCleanup()
	eventsUser := tc.unique("events-user")
	results = append(results, r.RunTest("cognito", "AdminListUserAuthEvents", func() error {
		if _, err := tc.client.SignUp(tc.ctx, &cognitoidentityprovider.SignUpInput{
			ClientId: aws.String(eventsClientID),
			Username: aws.String(eventsUser),
			Password: aws.String("EventsPass1!"),
		}); err != nil {
			return fmt.Errorf("self-service sign-up: %v", err)
		}

		events := map[string]string{}
		nextToken := ""
		for {
			input := &cognitoidentityprovider.AdminListUserAuthEventsInput{
				UserPoolId: aws.String(eventsPoolID),
				Username:   aws.String(eventsUser),
				MaxResults: aws.Int32(20),
			}
			if nextToken != "" {
				input.NextToken = aws.String(nextToken)
			}
			resp, err := tc.client.AdminListUserAuthEvents(tc.ctx, input)
			if err != nil {
				return fmt.Errorf("admin list user auth events: %v", err)
			}
			for _, e := range resp.AuthEvents {
				if e.EventType != "" {
					events[string(e.EventType)] = string(e.EventResponse)
				}
			}
			if resp.NextToken == nil {
				break
			}
			nextToken = *resp.NextToken
		}
		if response, ok := events["SignUp"]; !ok || response != "Pass" {
			return fmt.Errorf("SignUp event = (%q, %v), want Pass", response, ok)
		}
		if _, ok := events["SignIn"]; ok {
			return fmt.Errorf("unexpected SignIn event on a user that never signed in: %v", events)
		}
		return nil
	}))

	return results
}

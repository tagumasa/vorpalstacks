package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity/types"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	idptypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"vorpalstacks-sdk-tests/config"
)

// mintUserPoolIDToken provisions a throwaway user pool with a confirmed user
// and signs in through the admin password flow, returning the issuer-form
// Logins key and the ID token the user pool minted — the token role mappings
// read claims from. When roleGroup is true the user also joins a
// precedence-1 group whose role ARN is a cognito-trusted IAM role created
// here, so the token carries cognito:roles and cognito:preferred_role.
func (tc *cognitoIdentityContext) mintUserPoolIDToken(username string, roleGroup bool) (loginsKey, idToken string, cleanup func(), err error) {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: tc.r.endpoint,
		Region:   tc.r.region,
	})
	if err != nil {
		return "", "", nil, err
	}
	idp := cognitoidentityprovider.NewFromConfig(cfg)
	ctx := tc.ctx

	var cleanups []func()
	runCleanups := func() {
		for i := len(cleanups) - 1; i >= 0; i-- {
			cleanups[i]()
		}
	}
	// The deferred failure path runs the collected cleanups through the
	// local variable: an error return assigns nil to the named cleanup
	// result, and calling that would panic.
	defer func() {
		if err != nil {
			runCleanups()
		}
	}()

	poolResp, err := idp.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
		PoolName: aws.String(tc.unique("rolemap-pool")),
	})
	if err != nil {
		return "", "", nil, fmt.Errorf("create user pool: %w", err)
	}
	poolID := *poolResp.UserPool.Id
	cleanups = append(cleanups, func() {
		_, _ = idp.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{UserPoolId: aws.String(poolID)})
	})

	clientResp, err := idp.CreateUserPoolClient(ctx, &cognitoidentityprovider.CreateUserPoolClientInput{
		UserPoolId: aws.String(poolID),
		ClientName: aws.String(tc.unique("rolemap-client")),
	})
	if err != nil {
		return "", "", nil, fmt.Errorf("create app client: %w", err)
	}
	clientID := *clientResp.UserPoolClient.ClientId

	groupRoleArn := ""
	if roleGroup {
		iamClient := iam.NewFromConfig(cfg)
		roleName := tc.unique("cognito-rolemap")
		// The group role must be assumable by the Cognito service
		// principal; the pool's group creation validates the trust.
		const cognitoTrust = `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"cognito-identity.amazonaws.com"},"Action":"sts:AssumeRole"}]}`
		if _, err = iamClient.CreateRole(ctx, &iam.CreateRoleInput{
			RoleName:                 aws.String(roleName),
			AssumeRolePolicyDocument: aws.String(cognitoTrust),
		}); err != nil {
			return "", "", nil, fmt.Errorf("create IAM role: %w", err)
		}
		cleanups = append(cleanups, func() {
			_, _ = iamClient.DeleteRole(ctx, &iam.DeleteRoleInput{RoleName: aws.String(roleName)})
		})
		groupRoleArn = fmt.Sprintf("arn:aws:iam::%s:role/%s", tc.r.accountID, roleName)
		if _, err = idp.CreateGroup(ctx, &cognitoidentityprovider.CreateGroupInput{
			UserPoolId: aws.String(poolID),
			GroupName:  aws.String("rolemap-group"),
			RoleArn:    aws.String(groupRoleArn),
			Precedence: aws.Int32(1),
		}); err != nil {
			return "", "", nil, fmt.Errorf("create role group: %w", err)
		}
	}

	const password = "RoleMapPass123!"
	if _, err = idp.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId: aws.String(poolID),
		Username:   aws.String(username),
		// The name attribute is the claim the Rules mapping matches; the
		// ID token carries it as a claim.
		UserAttributes: []idptypes.AttributeType{
			{Name: aws.String("name"), Value: aws.String("RoleMap Tester")},
		},
		MessageAction: idptypes.MessageActionTypeSuppress,
	}); err != nil {
		return "", "", nil, fmt.Errorf("admin create user: %w", err)
	}
	if _, err = idp.AdminSetUserPassword(ctx, &cognitoidentityprovider.AdminSetUserPasswordInput{
		UserPoolId: aws.String(poolID),
		Username:   aws.String(username),
		Password:   aws.String(password),
		Permanent:  true,
	}); err != nil {
		return "", "", nil, fmt.Errorf("set permanent password: %w", err)
	}
	if roleGroup {
		if _, err = idp.AdminAddUserToGroup(ctx, &cognitoidentityprovider.AdminAddUserToGroupInput{
			UserPoolId: aws.String(poolID),
			Username:   aws.String(username),
			GroupName:  aws.String("rolemap-group"),
		}); err != nil {
			return "", "", nil, fmt.Errorf("add user to role group: %w", err)
		}
	}

	authResp, err := idp.AdminInitiateAuth(ctx, &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow:   "ADMIN_USER_PASSWORD_AUTH",
		UserPoolId: aws.String(poolID),
		ClientId:   aws.String(clientID),
		AuthParameters: map[string]string{
			"USERNAME": username,
			"PASSWORD": password,
		},
	})
	if err != nil {
		return "", "", nil, fmt.Errorf("admin initiate auth: %w", err)
	}
	if authResp.AuthenticationResult == nil || authResp.AuthenticationResult.IdToken == nil || *authResp.AuthenticationResult.IdToken == "" {
		return "", "", nil, fmt.Errorf("sign-in returned no ID token: %+v", authResp.AuthenticationResult)
	}
	loginsKey = fmt.Sprintf("cognito-idp.%s.amazonaws.com/%s", tc.r.region, poolID)
	return loginsKey, *authResp.AuthenticationResult.IdToken, runCleanups, nil
}

func (r *TestRunner) cognitoIdentityCredentialsTests(tc *cognitoIdentityContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cognito-identity", "GetCredentialsForIdentity", func() error {
		resp, err := tc.client.GetCredentialsForIdentity(tc.ctx, &cognitoidentity.GetCredentialsForIdentityInput{
			IdentityId: aws.String(tc.identityID),
		})
		if err != nil {
			return err
		}
		if *resp.IdentityId != tc.identityID {
			return fmt.Errorf("expected identity ID %s, got %s", tc.identityID, *resp.IdentityId)
		}
		if resp.Credentials == nil {
			return fmt.Errorf("credentials is nil")
		}
		if resp.Credentials.AccessKeyId == nil || *resp.Credentials.AccessKeyId == "" {
			return fmt.Errorf("AccessKeyId is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetOpenIdToken", func() error {
		resp, err := tc.client.GetOpenIdToken(tc.ctx, &cognitoidentity.GetOpenIdTokenInput{
			IdentityId: aws.String(tc.identityID),
		})
		if err != nil {
			return err
		}
		if *resp.IdentityId != tc.identityID {
			return fmt.Errorf("expected identity ID %s, got %s", tc.identityID, *resp.IdentityId)
		}
		if resp.Token == nil || *resp.Token == "" {
			return fmt.Errorf("token is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "GetOpenIdToken_WithLogins", func() error {
		resp, err := tc.client.GetOpenIdToken(tc.ctx, &cognitoidentity.GetOpenIdTokenInput{
			IdentityId: aws.String(tc.identityID),
			Logins: map[string]string{
				"graph.facebook.com": "new-token",
			},
		})
		if err != nil {
			return err
		}
		if resp.Token == nil || *resp.Token == "" {
			return fmt.Errorf("token is nil or empty")
		}
		return nil
	}))

	// The claim set role mappings evaluate comes from the linked platform
	// user-pool ID token: a Rules mapping matches a claim the token carries.
	results = append(results, r.RunTest("cognito-identity", "GetCredentialsForIdentity_RoleMappingRules", func() error {
		loginsKey, idToken, cleanupMint, err := tc.mintUserPoolIDToken(tc.unique("rolemap-rules"), false)
		if err != nil {
			return err
		}
		defer cleanupMint()

		poolID, cleanupPool, err := tc.createIdPool(tc.unique("test-rolemapping"))
		if err != nil {
			return err
		}
		defer cleanupPool()

		mappedRole := fmt.Sprintf("arn:aws:iam::%s:role/mapped-role", r.accountID)
		_, err = tc.client.SetIdentityPoolRoles(tc.ctx, &cognitoidentity.SetIdentityPoolRolesInput{
			IdentityPoolId: aws.String(poolID),
			Roles: map[string]string{
				"authenticated":   fmt.Sprintf("arn:aws:iam::%s:role/auth-role", r.accountID),
				"unauthenticated": fmt.Sprintf("arn:aws:iam::%s:role/unauth-role", r.accountID),
			},
			RoleMappings: map[string]types.RoleMapping{
				loginsKey: {
					Type:                    types.RoleMappingTypeRules,
					AmbiguousRoleResolution: types.AmbiguousRoleResolutionTypeAuthenticatedRole,
					RulesConfiguration: &types.RulesConfigurationType{
						Rules: []types.MappingRule{{
							Claim:     aws.String("name"),
							MatchType: types.MappingRuleMatchTypeEquals,
							Value:     aws.String("RoleMap Tester"),
							RoleARN:   aws.String(mappedRole),
						}},
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("set pool roles with mappings: %v", err)
		}

		idResp, err := tc.client.GetId(tc.ctx, &cognitoidentity.GetIdInput{
			IdentityPoolId: aws.String(poolID),
			Logins:         map[string]string{loginsKey: idToken},
		})
		if err != nil {
			return err
		}
		credResp, err := tc.client.GetCredentialsForIdentity(tc.ctx, &cognitoidentity.GetCredentialsForIdentityInput{
			IdentityId: idResp.IdentityId,
			Logins:     map[string]string{loginsKey: idToken},
		})
		if err != nil {
			return fmt.Errorf("credentials through a matching rule: %v", err)
		}
		if credResp.Credentials == nil || credResp.Credentials.AccessKeyId == nil || *credResp.Credentials.AccessKeyId == "" {
			return fmt.Errorf("credentials is nil or empty through a matching rule")
		}

		// A CustomRoleArn alongside a single granted role is ignored — the
		// parameter selects only among multiple received roles — and the
		// mapped role is issued either way.
		ignoredResp, err := tc.client.GetCredentialsForIdentity(tc.ctx, &cognitoidentity.GetCredentialsForIdentityInput{
			IdentityId:    idResp.IdentityId,
			CustomRoleArn: aws.String(fmt.Sprintf("arn:aws:iam::%s:role/other-role", r.accountID)),
			Logins:        map[string]string{loginsKey: idToken},
		})
		if err != nil {
			return fmt.Errorf("CustomRoleArn alongside a single granted role: %v", err)
		}
		if ignoredResp.Credentials == nil || ignoredResp.Credentials.AccessKeyId == nil || *ignoredResp.Credentials.AccessKeyId == "" {
			return fmt.Errorf("credentials is nil or empty alongside an ignored CustomRoleArn")
		}

		// The mapped role itself is an authorised CustomRoleArn.
		_, err = tc.client.GetCredentialsForIdentity(tc.ctx, &cognitoidentity.GetCredentialsForIdentityInput{
			IdentityId:    idResp.IdentityId,
			CustomRoleArn: aws.String(mappedRole),
			Logins:        map[string]string{loginsKey: idToken},
		})
		return err
	}))

	// A Token mapping with a Deny resolution reads the cognito:preferred_role
	// claim the ID token derives from the user's role-bearing group; a token
	// without role claims resolves ambiguous and is denied.
	results = append(results, r.RunTest("cognito-identity", "GetCredentialsForIdentity_TokenMappingFromIDToken", func() error {
		groupLoginsKey, groupToken, cleanupGroup, err := tc.mintUserPoolIDToken(tc.unique("rolemap-token"), true)
		if err != nil {
			return err
		}
		defer cleanupGroup()
		plainLoginsKey, plainToken, cleanupPlain, err := tc.mintUserPoolIDToken(tc.unique("rolemap-notoken"), false)
		if err != nil {
			return err
		}
		defer cleanupPlain()

		poolID, cleanupPool, err := tc.createIdPool(tc.unique("token-rolemapping"))
		if err != nil {
			return err
		}
		defer cleanupPool()

		_, err = tc.client.SetIdentityPoolRoles(tc.ctx, &cognitoidentity.SetIdentityPoolRolesInput{
			IdentityPoolId: aws.String(poolID),
			Roles: map[string]string{
				"authenticated":   fmt.Sprintf("arn:aws:iam::%s:role/auth-role", r.accountID),
				"unauthenticated": fmt.Sprintf("arn:aws:iam::%s:role/unauth-role", r.accountID),
			},
			RoleMappings: map[string]types.RoleMapping{
				groupLoginsKey: {Type: types.RoleMappingTypeToken, AmbiguousRoleResolution: types.AmbiguousRoleResolutionTypeDeny},
				plainLoginsKey: {Type: types.RoleMappingTypeToken, AmbiguousRoleResolution: types.AmbiguousRoleResolutionTypeDeny},
			},
		})
		if err != nil {
			return fmt.Errorf("set pool roles with token mappings: %v", err)
		}

		groupID, err := tc.client.GetId(tc.ctx, &cognitoidentity.GetIdInput{
			IdentityPoolId: aws.String(poolID),
			Logins:         map[string]string{groupLoginsKey: groupToken},
		})
		if err != nil {
			return err
		}
		credResp, err := tc.client.GetCredentialsForIdentity(tc.ctx, &cognitoidentity.GetCredentialsForIdentityInput{
			IdentityId: groupID.IdentityId,
			Logins:     map[string]string{groupLoginsKey: groupToken},
		})
		if err != nil {
			return fmt.Errorf("credentials through the group-derived preferred role: %v", err)
		}
		if credResp.Credentials == nil || credResp.Credentials.AccessKeyId == nil || *credResp.Credentials.AccessKeyId == "" {
			return fmt.Errorf("credentials is nil or empty through the preferred role")
		}

		// The token without role claims resolves ambiguous under Deny and
		// the request is refused.
		plainID, err := tc.client.GetId(tc.ctx, &cognitoidentity.GetIdInput{
			IdentityPoolId: aws.String(poolID),
			Logins:         map[string]string{plainLoginsKey: plainToken},
		})
		if err != nil {
			return err
		}
		_, err = tc.client.GetCredentialsForIdentity(tc.ctx, &cognitoidentity.GetCredentialsForIdentityInput{
			IdentityId: plainID.IdentityId,
			Logins:     map[string]string{plainLoginsKey: plainToken},
		})
		if err == nil {
			return fmt.Errorf("expected NotAuthorizedException for a token without role claims under a Deny resolution")
		}
		return AssertErrorContains(err, "NotAuthorizedException")
	}))

	// A presented platform token that fails validation — here a tampered
	// signature — fails the credential request closed. GetId accepts the
	// login link; the validation strikes when credentials are requested.
	results = append(results, r.RunTest("cognito-identity", "GetCredentialsForIdentity_ForgedTokenRefused", func() error {
		loginsKey, idToken, cleanupMint, err := tc.mintUserPoolIDToken(tc.unique("rolemap-forged"), false)
		if err != nil {
			return err
		}
		defer cleanupMint()

		poolID, cleanupPool, err := tc.createIdPool(tc.unique("forged-rolemapping"))
		if err != nil {
			return err
		}
		defer cleanupPool()

		_, err = tc.client.SetIdentityPoolRoles(tc.ctx, &cognitoidentity.SetIdentityPoolRolesInput{
			IdentityPoolId: aws.String(poolID),
			Roles: map[string]string{
				"authenticated": fmt.Sprintf("arn:aws:iam::%s:role/auth-role", r.accountID),
			},
			RoleMappings: map[string]types.RoleMapping{
				loginsKey: {Type: types.RoleMappingTypeToken, AmbiguousRoleResolution: types.AmbiguousRoleResolutionTypeAuthenticatedRole},
			},
		})
		if err != nil {
			return fmt.Errorf("set pool roles: %v", err)
		}

		idResp, err := tc.client.GetId(tc.ctx, &cognitoidentity.GetIdInput{
			IdentityPoolId: aws.String(poolID),
			Logins:         map[string]string{loginsKey: "forged." + idToken},
		})
		if err != nil {
			return err
		}
		_, err = tc.client.GetCredentialsForIdentity(tc.ctx, &cognitoidentity.GetCredentialsForIdentityInput{
			IdentityId: idResp.IdentityId,
			Logins:     map[string]string{loginsKey: "forged." + idToken},
		})
		if err == nil {
			return fmt.Errorf("expected NotAuthorizedException for a forged platform token")
		}
		return AssertErrorContains(err, "NotAuthorizedException")
	}))

	return results
}

package testutil

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunCognitoTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "cognito",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := cognitoidentityprovider.NewFromConfig(cfg)
	ctx := context.Background()

	userPoolName := fmt.Sprintf("test-pool-%d", time.Now().UnixNano())

	var userPoolID string
	results = append(results, r.RunTest("cognito", "CreateUserPool", func() error {
		resp, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(userPoolName),
			Policies: &types.UserPoolPolicyType{
				PasswordPolicy: &types.PasswordPolicyType{
					MinimumLength:    aws.Int32(8),
					RequireUppercase: true,
					RequireLowercase: true,
					RequireNumbers:   true,
					RequireSymbols:   false,
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.UserPool == nil {
			return fmt.Errorf("UserPool is nil")
		}
		if resp.UserPool.Id == nil {
			return fmt.Errorf("UserPool.Id is nil")
		}
		if resp.UserPool.Name == nil || *resp.UserPool.Name != userPoolName {
			return fmt.Errorf("UserPool.Name mismatch: got %v, want %s", resp.UserPool.Name, userPoolName)
		}
		if resp.UserPool.Arn == nil {
			return fmt.Errorf("UserPool.Arn is nil")
		}
		userPoolID = *resp.UserPool.Id
		return nil
	}))

	if userPoolID != "" {
		results = append(results, r.RunTest("cognito", "DescribeUserPool", func() error {
			resp, err := client.DescribeUserPool(ctx, &cognitoidentityprovider.DescribeUserPoolInput{
				UserPoolId: aws.String(userPoolID),
			})
			if err != nil {
				return err
			}
			if resp.UserPool == nil {
				return fmt.Errorf("UserPool is nil")
			}
			if resp.UserPool.Id == nil {
				return fmt.Errorf("UserPool.Id is nil")
			}
			if resp.UserPool.Name == nil || *resp.UserPool.Name != userPoolName {
				return fmt.Errorf("UserPool.Name mismatch: got %v, want %s", resp.UserPool.Name, userPoolName)
			}
			if resp.UserPool.Arn == nil {
				return fmt.Errorf("UserPool.Arn is nil")
			}
			return nil
		}))

		clientName := fmt.Sprintf("test-client-%d", time.Now().UnixNano())
		var clientID string
		results = append(results, r.RunTest("cognito", "CreateUserPoolClient", func() error {
			resp, err := client.CreateUserPoolClient(ctx, &cognitoidentityprovider.CreateUserPoolClientInput{
				UserPoolId: aws.String(userPoolID),
				ClientName: aws.String(clientName),
			})
			if err != nil {
				return err
			}
			if resp.UserPoolClient == nil {
				return fmt.Errorf("UserPoolClient is nil")
			}
			if resp.UserPoolClient.ClientId == nil {
				return fmt.Errorf("UserPoolClient.ClientId is nil")
			}
			if resp.UserPoolClient.ClientName == nil || *resp.UserPoolClient.ClientName != clientName {
				return fmt.Errorf("ClientName mismatch: got %v, want %s", resp.UserPoolClient.ClientName, clientName)
			}
			if resp.UserPoolClient.ClientSecret != nil && *resp.UserPoolClient.ClientSecret == "" {
				return fmt.Errorf("ClientSecret should not be empty string if set")
			}
			clientID = *resp.UserPoolClient.ClientId
			return nil
		}))

		if clientID != "" {
			results = append(results, r.RunTest("cognito", "DescribeUserPoolClient", func() error {
				resp, err := client.DescribeUserPoolClient(ctx, &cognitoidentityprovider.DescribeUserPoolClientInput{
					ClientId:   aws.String(clientID),
					UserPoolId: aws.String(userPoolID),
				})
				if err != nil {
					return err
				}
				if resp.UserPoolClient == nil {
					return fmt.Errorf("UserPoolClient is nil")
				}
				if resp.UserPoolClient.ClientId == nil || *resp.UserPoolClient.ClientId != clientID {
					return fmt.Errorf("ClientId mismatch: got %v, want %s", resp.UserPoolClient.ClientId, clientID)
				}
				if resp.UserPoolClient.ClientName == nil || *resp.UserPoolClient.ClientName != clientName {
					return fmt.Errorf("ClientName mismatch: got %v, want %s", resp.UserPoolClient.ClientName, clientName)
				}
				if resp.UserPoolClient.UserPoolId == nil || *resp.UserPoolClient.UserPoolId != userPoolID {
					return fmt.Errorf("UserPoolId mismatch: got %v, want %s", resp.UserPoolClient.UserPoolId, userPoolID)
				}
				return nil
			}))

			results = append(results, r.RunTest("cognito", "UpdateUserPoolClient", func() error {
				resp, err := client.UpdateUserPoolClient(ctx, &cognitoidentityprovider.UpdateUserPoolClientInput{
					ClientId:   aws.String(clientID),
					UserPoolId: aws.String(userPoolID),
					ClientName: aws.String("updated-client"),
				})
				if err != nil {
					return err
				}
				if resp.UserPoolClient == nil {
					return fmt.Errorf("UserPoolClient is nil")
				}
				if resp.UserPoolClient.ClientName == nil || *resp.UserPoolClient.ClientName != "updated-client" {
					return fmt.Errorf("ClientName not updated: got %v, want updated-client", resp.UserPoolClient.ClientName)
				}
				return nil
			}))
		}

		domainName := fmt.Sprintf("test-domain-%d", time.Now().UnixNano())
		results = append(results, r.RunTest("cognito", "CreateUserPoolDomain", func() error {
			resp, err := client.CreateUserPoolDomain(ctx, &cognitoidentityprovider.CreateUserPoolDomainInput{
				Domain:     aws.String(domainName),
				UserPoolId: aws.String(userPoolID),
			})
			if err != nil {
				return err
			}
			if resp.CloudFrontDomain == nil || *resp.CloudFrontDomain == "" {
				return fmt.Errorf("CloudFrontDomain is nil or empty")
			}
			return nil
		}))

		results = append(results, r.RunTest("cognito", "DescribeUserPoolDomain", func() error {
			resp, err := client.DescribeUserPoolDomain(ctx, &cognitoidentityprovider.DescribeUserPoolDomainInput{
				Domain: aws.String(domainName),
			})
			if err != nil {
				return err
			}
			if resp.DomainDescription == nil {
				return fmt.Errorf("DomainDescription is nil")
			}
			if resp.DomainDescription.UserPoolId == nil || *resp.DomainDescription.UserPoolId != userPoolID {
				return fmt.Errorf("UserPoolId mismatch: got %v, want %s", resp.DomainDescription.UserPoolId, userPoolID)
			}
			if resp.DomainDescription.Domain == nil || *resp.DomainDescription.Domain != domainName {
				return fmt.Errorf("Domain mismatch: got %v, want %s", resp.DomainDescription.Domain, domainName)
			}
			return nil
		}))

		results = append(results, r.RunTest("cognito", "ListUserPoolClients", func() error {
			resp, err := client.ListUserPoolClients(ctx, &cognitoidentityprovider.ListUserPoolClientsInput{
				UserPoolId: aws.String(userPoolID),
				MaxResults: aws.Int32(10),
			})
			if err != nil {
				return err
			}
			if len(resp.UserPoolClients) == 0 {
				return fmt.Errorf("expected at least one user pool client")
			}
			return nil
		}))

		results = append(results, r.RunTest("cognito", "ListUserPools", func() error {
			resp, err := client.ListUserPools(ctx, &cognitoidentityprovider.ListUserPoolsInput{
				MaxResults: aws.Int32(10),
			})
			if err != nil {
				return err
			}
			if len(resp.UserPools) == 0 {
				return fmt.Errorf("expected at least one user pool")
			}
			return nil
		}))

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
				return fmt.Errorf("Group is nil")
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
			return nil
		}))

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
				return fmt.Errorf("User is nil")
			}
			if resp.User.Username == nil || *resp.User.Username != username {
				return fmt.Errorf("Username mismatch: got %v, want %s", resp.User.Username, username)
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
				return fmt.Errorf("Username mismatch: got %v, want %s", resp.Username, username)
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
			return nil
		}))

		identifier := fmt.Sprintf("resource-%d", time.Now().UnixNano())
		results = append(results, r.RunTest("cognito", "CreateResourceServer", func() error {
			resp, err := client.CreateResourceServer(ctx, &cognitoidentityprovider.CreateResourceServerInput{
				UserPoolId: aws.String(userPoolID),
				Identifier: aws.String(identifier),
				Name:       aws.String("Test Resource Server"),
			})
			if err != nil {
				return err
			}
			if resp.ResourceServer == nil {
				return fmt.Errorf("ResourceServer is nil")
			}
			if resp.ResourceServer.Identifier == nil || *resp.ResourceServer.Identifier != identifier {
				return fmt.Errorf("Identifier mismatch: got %v, want %s", resp.ResourceServer.Identifier, identifier)
			}
			if resp.ResourceServer.Name == nil || *resp.ResourceServer.Name != "Test Resource Server" {
				return fmt.Errorf("Name mismatch: got %v, want Test Resource Server", resp.ResourceServer.Name)
			}
			if len(resp.ResourceServer.Scopes) != 0 {
				return fmt.Errorf("Scopes should be empty when no scopes specified")
			}
			return nil
		}))

		results = append(results, r.RunTest("cognito", "ListResourceServers", func() error {
			resp, err := client.ListResourceServers(ctx, &cognitoidentityprovider.ListResourceServersInput{
				UserPoolId: aws.String(userPoolID),
			})
			if err != nil {
				return err
			}
			if len(resp.ResourceServers) == 0 {
				return fmt.Errorf("expected at least one resource server")
			}
			return nil
		}))

		results = append(results, r.RunTest("cognito", "UpdateUserPool", func() error {
			_, err := client.UpdateUserPool(ctx, &cognitoidentityprovider.UpdateUserPoolInput{
				UserPoolId: aws.String(userPoolID),
				Policies: &types.UserPoolPolicyType{
					PasswordPolicy: &types.PasswordPolicyType{
						MinimumLength:    aws.Int32(10),
						RequireUppercase: true,
						RequireLowercase: true,
						RequireNumbers:   true,
						RequireSymbols:   true,
					},
				},
			})
			return err
		}))

		results = append(results, r.RunTest("cognito", "CreateIdentityProvider", func() error {
			resp, err := client.CreateIdentityProvider(ctx, &cognitoidentityprovider.CreateIdentityProviderInput{
				UserPoolId:   aws.String(userPoolID),
				ProviderName: aws.String("TestProvider"),
				ProviderType: types.IdentityProviderTypeType("Facebook"),
				ProviderDetails: map[string]string{
					"client_id":        "test-client-id",
					"client_secret":    "test-client-secret",
					"authorize_scopes": "public_profile,email",
				},
			})
			if err != nil {
				return err
			}
			if resp.IdentityProvider == nil {
				return fmt.Errorf("IdentityProvider is nil")
			}
			if resp.IdentityProvider.ProviderName == nil || *resp.IdentityProvider.ProviderName != "TestProvider" {
				return fmt.Errorf("ProviderName mismatch: got %v, want TestProvider", resp.IdentityProvider.ProviderName)
			}
			if resp.IdentityProvider.ProviderType != types.IdentityProviderTypeTypeFacebook {
				return fmt.Errorf("ProviderType mismatch: got %v, want Facebook", resp.IdentityProvider.ProviderType)
			}
			if resp.IdentityProvider.UserPoolId == nil || *resp.IdentityProvider.UserPoolId != userPoolID {
				return fmt.Errorf("UserPoolId mismatch: got %v, want %s", resp.IdentityProvider.UserPoolId, userPoolID)
			}
			return nil
		}))

		results = append(results, r.RunTest("cognito", "ListIdentityProviders", func() error {
			resp, err := client.ListIdentityProviders(ctx, &cognitoidentityprovider.ListIdentityProvidersInput{
				UserPoolId: aws.String(userPoolID),
			})
			if err != nil {
				return err
			}
			if len(resp.Providers) == 0 {
				return fmt.Errorf("expected at least one identity provider")
			}
			return nil
		}))

		results = append(results, r.RunTest("cognito", "SetUserPoolMfaConfig", func() error {
			_, err := client.SetUserPoolMfaConfig(ctx, &cognitoidentityprovider.SetUserPoolMfaConfigInput{
				UserPoolId:       aws.String(userPoolID),
				MfaConfiguration: types.UserPoolMfaTypeOn,
			})
			if err != nil {
				return fmt.Errorf("SetUserPoolMfaConfig failed: %v", err)
			}
			mfaResp, err := client.GetUserPoolMfaConfig(ctx, &cognitoidentityprovider.GetUserPoolMfaConfigInput{
				UserPoolId: aws.String(userPoolID),
			})
			if err != nil {
				return fmt.Errorf("GetUserPoolMfaConfig failed: %v", err)
			}
			if mfaResp.MfaConfiguration != types.UserPoolMfaTypeOn {
				return fmt.Errorf("expected MfaConfiguration ON, got %v", mfaResp.MfaConfiguration)
			}
			_, err = client.SetUserPoolMfaConfig(ctx, &cognitoidentityprovider.SetUserPoolMfaConfigInput{
				UserPoolId:       aws.String(userPoolID),
				MfaConfiguration: types.UserPoolMfaTypeOff,
			})
			if err != nil {
				return fmt.Errorf("SetUserPoolMfaConfig (OFF) failed: %v", err)
			}
			return nil
		}))

		results = append(results, r.RunTest("cognito", "GetUserPoolMfaConfig", func() error {
			resp, err := client.GetUserPoolMfaConfig(ctx, &cognitoidentityprovider.GetUserPoolMfaConfigInput{
				UserPoolId: aws.String(userPoolID),
			})
			if err != nil {
				return err
			}
			if resp.MfaConfiguration == "" && resp.SoftwareTokenMfaConfiguration == nil && resp.SmsMfaConfiguration == nil && resp.EmailMfaConfiguration == nil {
				return fmt.Errorf("expected at least one MFA config field to be set")
			}
			return nil
		}))

		results = append(results, r.RunTest("cognito", "AdminDisableUser", func() error {
			_, err := client.AdminDisableUser(ctx, &cognitoidentityprovider.AdminDisableUserInput{
				UserPoolId: aws.String(userPoolID),
				Username:   aws.String(username),
			})
			return err
		}))

		results = append(results, r.RunTest("cognito", "AdminEnableUser", func() error {
			_, err := client.AdminEnableUser(ctx, &cognitoidentityprovider.AdminEnableUserInput{
				UserPoolId: aws.String(userPoolID),
				Username:   aws.String(username),
			})
			return err
		}))

		results = append(results, r.RunTest("cognito", "AdminDeleteUser", func() error {
			_, err := client.AdminDeleteUser(ctx, &cognitoidentityprovider.AdminDeleteUserInput{
				UserPoolId: aws.String(userPoolID),
				Username:   aws.String(username),
			})
			return err
		}))

		results = append(results, r.RunTest("cognito", "DeleteUserPoolDomain", func() error {
			_, err := client.DeleteUserPoolDomain(ctx, &cognitoidentityprovider.DeleteUserPoolDomainInput{
				Domain:     aws.String(domainName),
				UserPoolId: aws.String(userPoolID),
			})
			return err
		}))

		if clientID != "" {
			results = append(results, r.RunTest("cognito", "DeleteUserPoolClient", func() error {
				_, err := client.DeleteUserPoolClient(ctx, &cognitoidentityprovider.DeleteUserPoolClientInput{
					ClientId:   aws.String(clientID),
					UserPoolId: aws.String(userPoolID),
				})
				return err
			}))
		}

		results = append(results, r.RunTest("cognito", "DeleteGroup", func() error {
			_, err := client.DeleteGroup(ctx, &cognitoidentityprovider.DeleteGroupInput{
				GroupName:  aws.String(groupName),
				UserPoolId: aws.String(userPoolID),
			})
			return err
		}))

		results = append(results, r.RunTest("cognito", "GetCSVHeader", func() error {
			resp, err := client.GetCSVHeader(ctx, &cognitoidentityprovider.GetCSVHeaderInput{
				UserPoolId: aws.String(userPoolID),
			})
			if err != nil {
				return err
			}
			if len(resp.CSVHeader) == 0 {
				return fmt.Errorf("expected non-empty CSV header")
			}
			return nil
		}))

		results = append(results, r.RunTest("cognito", "DescribeRiskConfiguration", func() error {
			resp, err := client.DescribeRiskConfiguration(ctx, &cognitoidentityprovider.DescribeRiskConfigurationInput{
				UserPoolId: aws.String(userPoolID),
			})
			if err != nil {
				return err
			}
			if resp.RiskConfiguration == nil {
				return fmt.Errorf("RiskConfiguration is nil")
			}
			return nil
		}))

		// === NEW BUG FIX VERIFICATION TESTS ===

		results = append(results, r.RunTest("cognito", "DescribeIdentityProvider", func() error {
			resp, err := client.DescribeIdentityProvider(ctx, &cognitoidentityprovider.DescribeIdentityProviderInput{
				UserPoolId:   aws.String(userPoolID),
				ProviderName: aws.String("TestProvider"),
			})
			if err != nil {
				return fmt.Errorf("DescribeIdentityProvider failed: %v", err)
			}
			if resp.IdentityProvider == nil {
				return fmt.Errorf("IdentityProvider is nil")
			}
			if resp.IdentityProvider.ProviderName == nil || *resp.IdentityProvider.ProviderName != "TestProvider" {
				return fmt.Errorf("ProviderName mismatch: got %v", resp.IdentityProvider.ProviderName)
			}
			if resp.IdentityProvider.ProviderType != types.IdentityProviderTypeTypeFacebook {
				return fmt.Errorf("ProviderType mismatch: got %v", resp.IdentityProvider.ProviderType)
			}
			return nil
		}))

		results = append(results, r.RunTest("cognito", "UpdateIdentityProvider", func() error {
			_, err := client.UpdateIdentityProvider(ctx, &cognitoidentityprovider.UpdateIdentityProviderInput{
				UserPoolId:      aws.String(userPoolID),
				ProviderName:    aws.String("TestProvider"),
				ProviderDetails: map[string]string{"updated_key": "updated_value"},
			})
			if err != nil {
				return fmt.Errorf("UpdateIdentityProvider failed: %v", err)
			}
			descResp, err := client.DescribeIdentityProvider(ctx, &cognitoidentityprovider.DescribeIdentityProviderInput{
				UserPoolId:   aws.String(userPoolID),
				ProviderName: aws.String("TestProvider"),
			})
			if err != nil {
				return fmt.Errorf("DescribeIdentityProvider after update failed: %v", err)
			}
			if descResp.IdentityProvider.ProviderDetails == nil {
				return fmt.Errorf("ProviderDetails is nil after update")
			}
			if descResp.IdentityProvider.ProviderDetails["updated_key"] != "updated_value" {
				return fmt.Errorf("ProviderDetails not updated")
			}
			return nil
		}))

		results = append(results, r.RunTest("cognito", "DeleteIdentityProvider", func() error {
			delProvider := fmt.Sprintf("del-provider-%d", time.Now().UnixNano())
			_, err := client.CreateIdentityProvider(ctx, &cognitoidentityprovider.CreateIdentityProviderInput{
				UserPoolId:      aws.String(userPoolID),
				ProviderName:    aws.String(delProvider),
				ProviderType:    types.IdentityProviderTypeTypeGoogle,
				ProviderDetails: map[string]string{"client_id": "test"},
			})
			if err != nil {
				return fmt.Errorf("create provider: %v", err)
			}
			_, err = client.DeleteIdentityProvider(ctx, &cognitoidentityprovider.DeleteIdentityProviderInput{
				UserPoolId:   aws.String(userPoolID),
				ProviderName: aws.String(delProvider),
			})
			if err != nil {
				return fmt.Errorf("DeleteIdentityProvider failed: %v", err)
			}
			_, err = client.DescribeIdentityProvider(ctx, &cognitoidentityprovider.DescribeIdentityProviderInput{
				UserPoolId:   aws.String(userPoolID),
				ProviderName: aws.String(delProvider),
			})
			if err == nil {
				return fmt.Errorf("expected error after deleting identity provider")
			}
			return nil
		}))

		results = append(results, r.RunTest("cognito", "DescribeResourceServer", func() error {
			resp, err := client.DescribeResourceServer(ctx, &cognitoidentityprovider.DescribeResourceServerInput{
				UserPoolId: aws.String(userPoolID),
				Identifier: aws.String(identifier),
			})
			if err != nil {
				return fmt.Errorf("DescribeResourceServer failed: %v", err)
			}
			if resp.ResourceServer == nil {
				return fmt.Errorf("ResourceServer is nil")
			}
			if resp.ResourceServer.Identifier == nil || *resp.ResourceServer.Identifier != identifier {
				return fmt.Errorf("Identifier mismatch: got %v", resp.ResourceServer.Identifier)
			}
			if resp.ResourceServer.Name == nil || *resp.ResourceServer.Name != "Test Resource Server" {
				return fmt.Errorf("Name mismatch: got %v", resp.ResourceServer.Name)
			}
			return nil
		}))

		results = append(results, r.RunTest("cognito", "UpdateResourceServer", func() error {
			resp, err := client.UpdateResourceServer(ctx, &cognitoidentityprovider.UpdateResourceServerInput{
				UserPoolId: aws.String(userPoolID),
				Identifier: aws.String(identifier),
				Name:       aws.String("Updated Resource Server"),
				Scopes: []types.ResourceServerScopeType{
					{ScopeName: aws.String("read"), ScopeDescription: aws.String("Read access")},
				},
			})
			if err != nil {
				return fmt.Errorf("UpdateResourceServer failed: %v", err)
			}
			if resp.ResourceServer == nil {
				return fmt.Errorf("ResourceServer is nil")
			}
			if resp.ResourceServer.Name == nil || *resp.ResourceServer.Name != "Updated Resource Server" {
				return fmt.Errorf("Name not updated: got %v", resp.ResourceServer.Name)
			}
			return nil
		}))

		results = append(results, r.RunTest("cognito", "DeleteResourceServer", func() error {
			delRS := fmt.Sprintf("del-rs-%d", time.Now().UnixNano())
			_, err := client.CreateResourceServer(ctx, &cognitoidentityprovider.CreateResourceServerInput{
				UserPoolId: aws.String(userPoolID),
				Identifier: aws.String(delRS),
				Name:       aws.String("Deletable RS"),
			})
			if err != nil {
				return fmt.Errorf("create resource server: %v", err)
			}
			_, err = client.DeleteResourceServer(ctx, &cognitoidentityprovider.DeleteResourceServerInput{
				UserPoolId: aws.String(userPoolID),
				Identifier: aws.String(delRS),
			})
			if err != nil {
				return fmt.Errorf("DeleteResourceServer failed: %v", err)
			}
			_, err = client.DescribeResourceServer(ctx, &cognitoidentityprovider.DescribeResourceServerInput{
				UserPoolId: aws.String(userPoolID),
				Identifier: aws.String(delRS),
			})
			if err == nil {
				return fmt.Errorf("expected error after deleting resource server")
			}
			return nil
		}))

		results = append(results, r.RunTest("cognito", "UpdateUserPoolDomain", func() error {
			udDomain := fmt.Sprintf("ud-domain-%d", time.Now().UnixNano())
			_, err := client.CreateUserPoolDomain(ctx, &cognitoidentityprovider.CreateUserPoolDomainInput{
				Domain:     aws.String(udDomain),
				UserPoolId: aws.String(userPoolID),
			})
			if err != nil {
				return fmt.Errorf("create domain: %v", err)
			}
			defer client.DeleteUserPoolDomain(ctx, &cognitoidentityprovider.DeleteUserPoolDomainInput{
				Domain:     aws.String(udDomain),
				UserPoolId: aws.String(userPoolID),
			})
			resp, err := client.UpdateUserPoolDomain(ctx, &cognitoidentityprovider.UpdateUserPoolDomainInput{
				Domain:     aws.String(udDomain),
				UserPoolId: aws.String(userPoolID),
			})
			if err != nil {
				return fmt.Errorf("UpdateUserPoolDomain failed: %v", err)
			}
			if resp.CloudFrontDomain == nil || *resp.CloudFrontDomain == "" {
				return fmt.Errorf("CloudFrontDomain is nil or empty")
			}
			return nil
		}))

		// === UNTESTED HANDLER TESTS ===

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
				return fmt.Errorf("Group is nil")
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
				return fmt.Errorf("Description not updated: got %v", resp.Group.Description)
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
			authResp, err := client.AdminInitiateAuth(ctx, &cognitoidentityprovider.AdminInitiateAuthInput{
				UserPoolId: aws.String(userPoolID),
				ClientId:   aws.String(clientID),
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
			_, err = client.AdminUserGlobalSignOut(ctx, &cognitoidentityprovider.AdminUserGlobalSignOutInput{
				UserPoolId: aws.String(userPoolID),
				Username:   aws.String(gsoUser),
			})
			if err != nil {
				return fmt.Errorf("AdminUserGlobalSignOut failed: %v", err)
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

		// === DELETE USER POOL ===

		results = append(results, r.RunTest("cognito", "DeleteUserPool", func() error {
			_, err := client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: aws.String(userPoolID),
			})
			return err
		}))

		// === TAG TESTS (own pools) ===

		results = append(results, r.RunTest("cognito", "TagResource", func() error {
			newPoolName := fmt.Sprintf("test-pool-tags-%d", time.Now().UnixNano())
			newPool, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
				PoolName: aws.String(newPoolName),
			})
			if err != nil {
				return err
			}
			if newPool.UserPool == nil || newPool.UserPool.Arn == nil {
				return fmt.Errorf("new pool Arn is nil")
			}
			_, err = client.TagResource(ctx, &cognitoidentityprovider.TagResourceInput{
				ResourceArn: newPool.UserPool.Arn,
				Tags: map[string]string{
					"Environment": "test",
					"Owner":       "test-user",
				},
			})
			if err != nil {
				return err
			}
			listResp, listErr := client.ListTagsForResource(ctx, &cognitoidentityprovider.ListTagsForResourceInput{
				ResourceArn: newPool.UserPool.Arn,
			})
			if listErr != nil {
				client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
					UserPoolId: newPool.UserPool.Id,
				})
				return listErr
			}
			if listResp.Tags == nil {
				client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
					UserPoolId: newPool.UserPool.Id,
				})
				return fmt.Errorf("Tags is nil after tagging")
			}
			if listResp.Tags["Environment"] != "test" {
				client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
					UserPoolId: newPool.UserPool.Id,
				})
				return fmt.Errorf("tag Environment not found after TagResource")
			}
			client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: newPool.UserPool.Id,
			})
			return nil
		}))

		results = append(results, r.RunTest("cognito", "ListTagsForResource", func() error {
			newPoolName := fmt.Sprintf("test-pool-listtags-%d", time.Now().UnixNano())
			newPool, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
				PoolName: aws.String(newPoolName),
			})
			if err != nil {
				return err
			}
			if newPool.UserPool == nil || newPool.UserPool.Arn == nil {
				return fmt.Errorf("new pool Arn is nil")
			}
			_, err = client.TagResource(ctx, &cognitoidentityprovider.TagResourceInput{
				ResourceArn: newPool.UserPool.Arn,
				Tags: map[string]string{
					"Test": "value",
				},
			})
			if err != nil {
				return err
			}
			resp, err := client.ListTagsForResource(ctx, &cognitoidentityprovider.ListTagsForResourceInput{
				ResourceArn: newPool.UserPool.Arn,
			})
			if err != nil {
				client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
					UserPoolId: newPool.UserPool.Id,
				})
				return err
			}
			if resp.Tags == nil {
				client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
					UserPoolId: newPool.UserPool.Id,
				})
				return fmt.Errorf("Tags is nil")
			}
			if resp.Tags["Test"] != "value" {
				client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
					UserPoolId: newPool.UserPool.Id,
				})
				return fmt.Errorf("expected tag Test=value, got %v", resp.Tags["Test"])
			}
			client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: newPool.UserPool.Id,
			})
			return nil
		}))

		results = append(results, r.RunTest("cognito", "UntagResource", func() error {
			newPoolName := fmt.Sprintf("test-pool-untag-%d", time.Now().UnixNano())
			newPool, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
				PoolName: aws.String(newPoolName),
			})
			if err != nil {
				return err
			}
			if newPool.UserPool == nil || newPool.UserPool.Arn == nil {
				return fmt.Errorf("new pool Arn is nil")
			}
			_, err = client.TagResource(ctx, &cognitoidentityprovider.TagResourceInput{
				ResourceArn: newPool.UserPool.Arn,
				Tags: map[string]string{
					"Test": "value",
				},
			})
			if err != nil {
				return err
			}
			_, err = client.UntagResource(ctx, &cognitoidentityprovider.UntagResourceInput{
				ResourceArn: newPool.UserPool.Arn,
				TagKeys:     []string{"Test"},
			})
			if err != nil {
				return err
			}
			listResp, listErr := client.ListTagsForResource(ctx, &cognitoidentityprovider.ListTagsForResourceInput{
				ResourceArn: newPool.UserPool.Arn,
			})
			if listErr != nil {
				client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
					UserPoolId: newPool.UserPool.Id,
				})
				return listErr
			}
			if _, exists := listResp.Tags["Test"]; exists {
				client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
					UserPoolId: newPool.UserPool.Id,
				})
				return fmt.Errorf("tag Test should have been removed after UntagResource")
			}
			client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: newPool.UserPool.Id,
			})
			return nil
		}))
	}

	// === ERROR / EDGE CASE TESTS ===

	results = append(results, r.RunTest("cognito", "DescribeUserPool_NonExistent", func() error {
		_, err := client.DescribeUserPool(ctx, &cognitoidentityprovider.DescribeUserPoolInput{
			UserPoolId: aws.String("us-east-1_nonexistentpool"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent user pool")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "DeleteUserPool_NonExistent", func() error {
		_, err := client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: aws.String("us-east-1_nonexistentpool"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent user pool")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "AdminGetUser_NonExistent", func() error {
		errPoolName := fmt.Sprintf("err-pool-%d", time.Now().UnixNano())
		createResp, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(errPoolName),
		})
		if err != nil {
			return fmt.Errorf("create pool: %v", err)
		}
		defer client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: createResp.UserPool.Id,
		})

		_, err = client.AdminGetUser(ctx, &cognitoidentityprovider.AdminGetUserInput{
			UserPoolId: createResp.UserPool.Id,
			Username:   aws.String("nonexistent-user-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent user")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "CreateUserPool_DuplicateName", func() error {
		dupPoolName := fmt.Sprintf("dup-pool-%d", time.Now().UnixNano())
		resp1, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(dupPoolName),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: resp1.UserPool.Id,
		})

		resp2, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(dupPoolName),
		})
		if err != nil {
			return fmt.Errorf("duplicate name should be allowed (unique IDs), got: %v", err)
		}
		if resp2.UserPool.Id == nil || *resp2.UserPool.Id == *resp1.UserPool.Id {
			return fmt.Errorf("duplicate pool should have different ID")
		}
		defer client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: resp2.UserPool.Id,
		})
		return nil
	}))

	results = append(results, r.RunTest("cognito", "AdminCreateUser_VerifyAttributes", func() error {
		errPoolName := fmt.Sprintf("attr-pool-%d", time.Now().UnixNano())
		createResp, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(errPoolName),
		})
		if err != nil {
			return fmt.Errorf("create pool: %v", err)
		}
		defer client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: createResp.UserPool.Id,
		})

		attrUser := fmt.Sprintf("attr-user-%d", time.Now().UnixNano())
		createUserResp, err := client.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
			UserPoolId:        createResp.UserPool.Id,
			Username:          aws.String(attrUser),
			TemporaryPassword: aws.String("TempPass123!"),
			MessageAction:     types.MessageActionTypeSuppress,
			UserAttributes: []types.AttributeType{
				{Name: aws.String("email"), Value: aws.String("test@example.com")},
				{Name: aws.String("name"), Value: aws.String("Test User")},
			},
		})
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		if createUserResp.User == nil {
			return fmt.Errorf("user is nil")
		}
		if createUserResp.User.Username == nil || *createUserResp.User.Username != attrUser {
			return fmt.Errorf("username mismatch")
		}
		if !createUserResp.User.Enabled {
			return fmt.Errorf("user should be enabled")
		}
		if createUserResp.User.UserStatus != types.UserStatusTypeForceChangePassword {
			return fmt.Errorf("expected FORCE_CHANGE_PASSWORD status, got %v", createUserResp.User.UserStatus)
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "ListUsers_ContainsCreated", func() error {
		errPoolName := fmt.Sprintf("list-pool-%d", time.Now().UnixNano())
		createResp, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(errPoolName),
		})
		if err != nil {
			return fmt.Errorf("create pool: %v", err)
		}
		defer client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: createResp.UserPool.Id,
		})

		listUser := fmt.Sprintf("list-user-%d", time.Now().UnixNano())
		_, err = client.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
			UserPoolId:        createResp.UserPool.Id,
			Username:          aws.String(listUser),
			TemporaryPassword: aws.String("TempPass123!"),
			MessageAction:     types.MessageActionTypeSuppress,
		})
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}

		resp, err := client.ListUsers(ctx, &cognitoidentityprovider.ListUsersInput{
			UserPoolId: createResp.UserPool.Id,
		})
		if err != nil {
			return err
		}
		found := false
		for _, u := range resp.Users {
			if u.Username != nil && *u.Username == listUser {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created user not found in ListUsers")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "ListGroups_ContainsCreated", func() error {
		errPoolName := fmt.Sprintf("grp-pool-%d", time.Now().UnixNano())
		createResp, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(errPoolName),
		})
		if err != nil {
			return fmt.Errorf("create pool: %v", err)
		}
		defer client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: createResp.UserPool.Id,
		})

		testGroup := fmt.Sprintf("test-grp-%d", time.Now().UnixNano())
		_, err = client.CreateGroup(ctx, &cognitoidentityprovider.CreateGroupInput{
			GroupName:   aws.String(testGroup),
			UserPoolId:  createResp.UserPool.Id,
			Description: aws.String("Test group description"),
		})
		if err != nil {
			return fmt.Errorf("create group: %v", err)
		}

		resp, err := client.ListGroups(ctx, &cognitoidentityprovider.ListGroupsInput{
			UserPoolId: createResp.UserPool.Id,
		})
		if err != nil {
			return err
		}
		found := false
		for _, g := range resp.Groups {
			if g.GroupName != nil && *g.GroupName == testGroup {
				found = true
				if g.Description == nil || *g.Description != "Test group description" {
					return fmt.Errorf("group description mismatch")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created group not found in ListGroups")
		}
		return nil
	}))

	// === ADDITIONAL ERROR / EDGE CASE TESTS ===

	results = append(results, r.RunTest("cognito", "GetGroup_NonExistent", func() error {
		nePoolName := fmt.Sprintf("ge-pool-%d", time.Now().UnixNano())
		createResp, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(nePoolName),
		})
		if err != nil {
			return fmt.Errorf("create pool: %v", err)
		}
		defer client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: createResp.UserPool.Id,
		})
		_, err = client.GetGroup(ctx, &cognitoidentityprovider.GetGroupInput{
			GroupName:  aws.String("nonexistent-group-xyz"),
			UserPoolId: createResp.UserPool.Id,
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent group")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "DescribeIdentityProvider_NonExistent", func() error {
		nePoolName := fmt.Sprintf("dip-pool-%d", time.Now().UnixNano())
		createResp, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(nePoolName),
		})
		if err != nil {
			return fmt.Errorf("create pool: %v", err)
		}
		defer client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: createResp.UserPool.Id,
		})
		_, err = client.DescribeIdentityProvider(ctx, &cognitoidentityprovider.DescribeIdentityProviderInput{
			UserPoolId:   createResp.UserPool.Id,
			ProviderName: aws.String("nonexistent-idp-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent identity provider")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "DescribeResourceServer_NonExistent", func() error {
		nePoolName := fmt.Sprintf("drs-pool-%d", time.Now().UnixNano())
		createResp, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(nePoolName),
		})
		if err != nil {
			return fmt.Errorf("create pool: %v", err)
		}
		defer client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: createResp.UserPool.Id,
		})
		_, err = client.DescribeResourceServer(ctx, &cognitoidentityprovider.DescribeResourceServerInput{
			UserPoolId: createResp.UserPool.Id,
			Identifier: aws.String("nonexistent-rs-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent resource server")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "DeleteIdentityProvider_NonExistent", func() error {
		nePoolName := fmt.Sprintf("dlip-pool-%d", time.Now().UnixNano())
		createResp, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(nePoolName),
		})
		if err != nil {
			return fmt.Errorf("create pool: %v", err)
		}
		defer client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: createResp.UserPool.Id,
		})
		_, err = client.DeleteIdentityProvider(ctx, &cognitoidentityprovider.DeleteIdentityProviderInput{
			UserPoolId:   createResp.UserPool.Id,
			ProviderName: aws.String("nonexistent-idp-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent identity provider")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "DeleteResourceServer_NonExistent", func() error {
		nePoolName := fmt.Sprintf("dlrs-pool-%d", time.Now().UnixNano())
		createResp, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(nePoolName),
		})
		if err != nil {
			return fmt.Errorf("create pool: %v", err)
		}
		defer client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: createResp.UserPool.Id,
		})
		_, err = client.DeleteResourceServer(ctx, &cognitoidentityprovider.DeleteResourceServerInput{
			UserPoolId: createResp.UserPool.Id,
			Identifier: aws.String("nonexistent-rs-xyz"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent resource server")
		}
		return nil
	}))

	return results
}

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

	// Test 1: Create User Pool
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
		// Test 2: Describe User Pool
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

		// Test 3: Create User Pool Client
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

		// Test 4: Describe User Pool Client
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

			// Test 5: Update User Pool Client
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

		// Test 6: Create User Pool Domain
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

		// Test 7: Describe User Pool Domain
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

		// Test 8: List User Pool Clients
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

		// Test 9: List User Pools
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

		// Test 10: Create Group
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

		// Test 11: List Groups
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

		// Test 12: Admin Create User
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

		// Test 14: Admin Get User
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

		// Test 15: List Users
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

		// Test 16: Create Resource Server
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

		// Test 17: List Resource Servers
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

		// Test 18: Update User Pool
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

		// Test 19: Create Identity Provider
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

		// Test 20: List Identity Providers
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

		// Test 21: Set User Pool MFA Config
		results = append(results, r.RunTest("cognito", "SetUserPoolMfaConfig", func() error {
			_, err := client.SetUserPoolMfaConfig(ctx, &cognitoidentityprovider.SetUserPoolMfaConfigInput{
				UserPoolId: aws.String(userPoolID),
				SmsMfaConfiguration: &types.SmsMfaConfigType{
					SmsConfiguration: &types.SmsConfigurationType{
						SnsCallerArn: aws.String("arn:aws:sns:us-east-1:123456789012:sms-topic"),
						ExternalId:   aws.String("external-id"),
					},
				},
			})
			return err
		}))

		// Test 22: Get User Pool MFA Config
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

		// Test 23: Admin Disable User
		results = append(results, r.RunTest("cognito", "AdminDisableUser", func() error {
			_, err := client.AdminDisableUser(ctx, &cognitoidentityprovider.AdminDisableUserInput{
				UserPoolId: aws.String(userPoolID),
				Username:   aws.String(username),
			})
			return err
		}))

		// Test 24: Admin Enable User
		results = append(results, r.RunTest("cognito", "AdminEnableUser", func() error {
			_, err := client.AdminEnableUser(ctx, &cognitoidentityprovider.AdminEnableUserInput{
				UserPoolId: aws.String(userPoolID),
				Username:   aws.String(username),
			})
			return err
		}))

		// Test 25: Admin Delete User
		results = append(results, r.RunTest("cognito", "AdminDeleteUser", func() error {
			_, err := client.AdminDeleteUser(ctx, &cognitoidentityprovider.AdminDeleteUserInput{
				UserPoolId: aws.String(userPoolID),
				Username:   aws.String(username),
			})
			return err
		}))

		// Test 26: Delete User Pool Domain
		results = append(results, r.RunTest("cognito", "DeleteUserPoolDomain", func() error {
			_, err := client.DeleteUserPoolDomain(ctx, &cognitoidentityprovider.DeleteUserPoolDomainInput{
				Domain:     aws.String(domainName),
				UserPoolId: aws.String(userPoolID),
			})
			return err
		}))

		// Test 27: Delete User Pool Client
		if clientID != "" {
			results = append(results, r.RunTest("cognito", "DeleteUserPoolClient", func() error {
				_, err := client.DeleteUserPoolClient(ctx, &cognitoidentityprovider.DeleteUserPoolClientInput{
					ClientId:   aws.String(clientID),
					UserPoolId: aws.String(userPoolID),
				})
				return err
			}))
		}

		// Test 28: Delete Group
		results = append(results, r.RunTest("cognito", "DeleteGroup", func() error {
			_, err := client.DeleteGroup(ctx, &cognitoidentityprovider.DeleteGroupInput{
				GroupName:  aws.String(groupName),
				UserPoolId: aws.String(userPoolID),
			})
			return err
		}))

		// Test 29: Get CSV Header (before deleting user pool)
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

		// Test 30: Describe Risk Configuration (before deleting user pool)
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

		// Test 31: Delete User Pool
		results = append(results, r.RunTest("cognito", "DeleteUserPool", func() error {
			_, err := client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: aws.String(userPoolID),
			})
			return err
		}))

		// Test 32: Tag Resource (with new pool)
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

		// Test 33: List Tags for Resource
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

		// Test 34: Untag Resource
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

		// Test 35: Global Sign Out
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

	return results
}

package testutil

import (
	"context"
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
	results = append(results, r.RunTest("cognito", "CreateUserPool", func() error {
		_, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
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
		return err
	}))

	// Get the User Pool ID for subsequent tests
	var userPoolID string
	if len(results) > 0 && results[0].Error == "" {
		listResp, _ := client.ListUserPools(ctx, &cognitoidentityprovider.ListUserPoolsInput{
			MaxResults: aws.Int32(10),
		})
		if len(listResp.UserPools) > 0 {
			userPoolID = aws.ToString(listResp.UserPools[0].Id)
		}
	}

	if userPoolID != "" {
		// Test 2: Describe User Pool
		results = append(results, r.RunTest("cognito", "DescribeUserPool", func() error {
			_, err := client.DescribeUserPool(ctx, &cognitoidentityprovider.DescribeUserPoolInput{
				UserPoolId: aws.String(userPoolID),
			})
			return err
		}))

		// Test 3: Create User Pool Client
		clientName := fmt.Sprintf("test-client-%d", time.Now().UnixNano())
		var clientID string
		results = append(results, r.RunTest("cognito", "CreateUserPoolClient", func() error {
			resp, err := client.CreateUserPoolClient(ctx, &cognitoidentityprovider.CreateUserPoolClientInput{
				UserPoolId: aws.String(userPoolID),
				ClientName: aws.String(clientName),
			})
			if err == nil {
				clientID = aws.ToString(resp.UserPoolClient.ClientId)
			}
			return err
		}))

		// Test 4: Describe User Pool Client
		if clientID != "" {
			results = append(results, r.RunTest("cognito", "DescribeUserPoolClient", func() error {
				_, err := client.DescribeUserPoolClient(ctx, &cognitoidentityprovider.DescribeUserPoolClientInput{
					ClientId:   aws.String(clientID),
					UserPoolId: aws.String(userPoolID),
				})
				return err
			}))

			// Test 5: Update User Pool Client
			results = append(results, r.RunTest("cognito", "UpdateUserPoolClient", func() error {
				_, err := client.UpdateUserPoolClient(ctx, &cognitoidentityprovider.UpdateUserPoolClientInput{
					ClientId:   aws.String(clientID),
					UserPoolId: aws.String(userPoolID),
					ClientName: aws.String("updated-client"),
				})
				return err
			}))
		}

		// Test 6: Create User Pool Domain
		domainName := fmt.Sprintf("test-domain-%d", time.Now().UnixNano())
		results = append(results, r.RunTest("cognito", "CreateUserPoolDomain", func() error {
			_, err := client.CreateUserPoolDomain(ctx, &cognitoidentityprovider.CreateUserPoolDomainInput{
				Domain:     aws.String(domainName),
				UserPoolId: aws.String(userPoolID),
			})
			return err
		}))

		// Test 7: Describe User Pool Domain
		results = append(results, r.RunTest("cognito", "DescribeUserPoolDomain", func() error {
			_, err := client.DescribeUserPoolDomain(ctx, &cognitoidentityprovider.DescribeUserPoolDomainInput{
				Domain: aws.String(domainName),
			})
			return err
		}))

		// Test 8: List User Pool Clients
		results = append(results, r.RunTest("cognito", "ListUserPoolClients", func() error {
			_, err := client.ListUserPoolClients(ctx, &cognitoidentityprovider.ListUserPoolClientsInput{
				UserPoolId: aws.String(userPoolID),
				MaxResults: aws.Int32(10),
			})
			return err
		}))

		// Test 9: List User Pools
		results = append(results, r.RunTest("cognito", "ListUserPools", func() error {
			_, err := client.ListUserPools(ctx, &cognitoidentityprovider.ListUserPoolsInput{
				MaxResults: aws.Int32(10),
			})
			return err
		}))

		// Test 10: Create Group
		groupName := fmt.Sprintf("group-%d", time.Now().UnixNano())
		results = append(results, r.RunTest("cognito", "CreateGroup", func() error {
			_, err := client.CreateGroup(ctx, &cognitoidentityprovider.CreateGroupInput{
				GroupName:  aws.String(groupName),
				UserPoolId: aws.String(userPoolID),
			})
			return err
		}))

		// Test 11: List Groups
		results = append(results, r.RunTest("cognito", "ListGroups", func() error {
			_, err := client.ListGroups(ctx, &cognitoidentityprovider.ListGroupsInput{
				UserPoolId: aws.String(userPoolID),
			})
			return err
		}))

		// Test 12: Admin Create User
		username := fmt.Sprintf("user-%d", time.Now().UnixNano())
		results = append(results, r.RunTest("cognito", "AdminCreateUser", func() error {
			_, err := client.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
				UserPoolId:        aws.String(userPoolID),
				Username:          aws.String(username),
				TemporaryPassword: aws.String("TempPass123!"),
				MessageAction:     types.MessageActionTypeSuppress,
			})
			return err
		}))

		// Test 14: Admin Get User
		results = append(results, r.RunTest("cognito", "AdminGetUser", func() error {
			_, err := client.AdminGetUser(ctx, &cognitoidentityprovider.AdminGetUserInput{
				UserPoolId: aws.String(userPoolID),
				Username:   aws.String(username),
			})
			return err
		}))

		// Test 15: List Users
		results = append(results, r.RunTest("cognito", "ListUsers", func() error {
			_, err := client.ListUsers(ctx, &cognitoidentityprovider.ListUsersInput{
				UserPoolId: aws.String(userPoolID),
			})
			return err
		}))

		// Test 16: Create Resource Server
		identifier := fmt.Sprintf("resource-%d", time.Now().UnixNano())
		results = append(results, r.RunTest("cognito", "CreateResourceServer", func() error {
			_, err := client.CreateResourceServer(ctx, &cognitoidentityprovider.CreateResourceServerInput{
				UserPoolId: aws.String(userPoolID),
				Identifier: aws.String(identifier),
				Name:       aws.String("Test Resource Server"),
			})
			return err
		}))

		// Test 17: List Resource Servers
		results = append(results, r.RunTest("cognito", "ListResourceServers", func() error {
			_, err := client.ListResourceServers(ctx, &cognitoidentityprovider.ListResourceServersInput{
				UserPoolId: aws.String(userPoolID),
			})
			return err
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
			_, err := client.CreateIdentityProvider(ctx, &cognitoidentityprovider.CreateIdentityProviderInput{
				UserPoolId:   aws.String(userPoolID),
				ProviderName: aws.String("TestProvider"),
				ProviderType: types.IdentityProviderTypeType("Facebook"),
				ProviderDetails: map[string]string{
					"client_id":        "test-client-id",
					"client_secret":    "test-client-secret",
					"authorize_scopes": "public_profile,email",
				},
			})
			return err
		}))

		// Test 20: List Identity Providers
		results = append(results, r.RunTest("cognito", "ListIdentityProviders", func() error {
			_, err := client.ListIdentityProviders(ctx, &cognitoidentityprovider.ListIdentityProvidersInput{
				UserPoolId: aws.String(userPoolID),
			})
			return err
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
			_, err := client.GetUserPoolMfaConfig(ctx, &cognitoidentityprovider.GetUserPoolMfaConfigInput{
				UserPoolId: aws.String(userPoolID),
			})
			return err
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
			_, err := client.GetCSVHeader(ctx, &cognitoidentityprovider.GetCSVHeaderInput{
				UserPoolId: aws.String(userPoolID),
			})
			return err
		}))

		// Test 30: Describe Risk Configuration (before deleting user pool)
		results = append(results, r.RunTest("cognito", "DescribeRiskConfiguration", func() error {
			_, err := client.DescribeRiskConfiguration(ctx, &cognitoidentityprovider.DescribeRiskConfigurationInput{
				UserPoolId: aws.String(userPoolID),
			})
			return err
		}))

		// Test 31: Delete User Pool
		results = append(results, r.RunTest("cognito", "DeleteUserPool", func() error {
			_, err := client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: aws.String(userPoolID),
			})
			return err
		}))

		// Test 30: Tag Resource (with new pool)
		results = append(results, r.RunTest("cognito", "TagResource", func() error {
			newPoolName := fmt.Sprintf("test-pool-tags-%d", time.Now().UnixNano())
			newPool, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
				PoolName: aws.String(newPoolName),
			})
			if err != nil {
				return err
			}
			_, err = client.TagResource(ctx, &cognitoidentityprovider.TagResourceInput{
				ResourceArn: newPool.UserPool.Arn,
				Tags: map[string]string{
					"Environment": "test",
					"Owner":       "test-user",
				},
			})
			// Cleanup
			client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: newPool.UserPool.Id,
			})
			return err
		}))

		// Test 31: List Tags for Resource
		results = append(results, r.RunTest("cognito", "ListTagsForResource", func() error {
			newPoolName := fmt.Sprintf("test-pool-listtags-%d", time.Now().UnixNano())
			newPool, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
				PoolName: aws.String(newPoolName),
			})
			if err != nil {
				return err
			}
			client.TagResource(ctx, &cognitoidentityprovider.TagResourceInput{
				ResourceArn: newPool.UserPool.Arn,
				Tags: map[string]string{
					"Test": "value",
				},
			})
			_, err = client.ListTagsForResource(ctx, &cognitoidentityprovider.ListTagsForResourceInput{
				ResourceArn: newPool.UserPool.Arn,
			})
			// Cleanup
			client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: newPool.UserPool.Id,
			})
			return err
		}))

		// Test 32: Untag Resource
		results = append(results, r.RunTest("cognito", "UntagResource", func() error {
			newPoolName := fmt.Sprintf("test-pool-untag-%d", time.Now().UnixNano())
			newPool, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
				PoolName: aws.String(newPoolName),
			})
			if err != nil {
				return err
			}
			client.TagResource(ctx, &cognitoidentityprovider.TagResourceInput{
				ResourceArn: newPool.UserPool.Arn,
				Tags: map[string]string{
					"Test": "value",
				},
			})
			_, err = client.UntagResource(ctx, &cognitoidentityprovider.UntagResourceInput{
				ResourceArn: newPool.UserPool.Arn,
				TagKeys:     []string{"Test"},
			})
			// Cleanup
			client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: newPool.UserPool.Id,
			})
			return err
		}))

		// Test 35: Global Sign Out (will fail but should handle gracefully)
		results = append(results, r.RunTest("cognito", "GlobalSignOut", func() error {
			// This requires an access token, which we don't have
			// Just test that the API is callable
			_, err := client.GlobalSignOut(ctx, &cognitoidentityprovider.GlobalSignOutInput{
				AccessToken: aws.String("dummy-token"),
			})
			// Expected to fail, but API should be accessible
			return err
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

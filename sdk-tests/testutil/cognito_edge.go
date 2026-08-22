package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func (r *TestRunner) cognitoEdgeCaseTests(ctx context.Context, client *cognitoidentityprovider.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cognito", "DescribeUserPool_NonExistent", func() error {
		_, err := client.DescribeUserPool(ctx, &cognitoidentityprovider.DescribeUserPoolInput{
			UserPoolId: aws.String(fmt.Sprintf("%s_nonexistentpool", r.region)),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "DeleteUserPool_NonExistent", func() error {
		_, err := client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: aws.String(fmt.Sprintf("%s_nonexistentpool", r.region)),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
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
		if err := AssertErrorContains(err, "UserNotFoundException"); err != nil {
			return err
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
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
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
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
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
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
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
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
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
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "ListUserPools_Pagination", func() error {
		pgTs := fmt.Sprintf("%d", time.Now().UnixNano())
		var pgPoolIds []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("PagPool-%s-%d", pgTs, i)
			createResp, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
				PoolName: aws.String(name),
			})
			if err != nil {
				for _, pid := range pgPoolIds {
					client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{UserPoolId: aws.String(pid)})
				}
				return fmt.Errorf("create user pool %s: %v", name, err)
			}
			if createResp.UserPool == nil || createResp.UserPool.Id == nil {
				for _, pid := range pgPoolIds {
					client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{UserPoolId: aws.String(pid)})
				}
				return fmt.Errorf("create user pool %s: UserPool.Id is nil", name)
			}
			pgPoolIds = append(pgPoolIds, *createResp.UserPool.Id)
		}
		pageCount := 0
		var nextToken *string
		for {
			resp, err := client.ListUserPools(ctx, &cognitoidentityprovider.ListUserPoolsInput{
				MaxResults: aws.Int32(2),
				NextToken:  nextToken,
			})
			if err != nil {
				for _, pid := range pgPoolIds {
					client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{UserPoolId: aws.String(pid)})
				}
				return fmt.Errorf("list user pools page: %v", err)
			}
			pageCount++
			if resp.NextToken != nil && *resp.NextToken != "" {
				nextToken = resp.NextToken
			} else {
				break
			}
		}
		for _, pid := range pgPoolIds {
			client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{UserPoolId: aws.String(pid)})
		}
		if pageCount < 2 {
			return fmt.Errorf("expected at least 2 pages with MaxResults=2, got %d", pageCount)
		}
		return nil
	}))

	// MaxResults above the Smithy PoolQueryLimitType maximum of 60 must be
	// rejected with InvalidParameterException instead of being silently
	// clamped to the default.
	results = append(results, r.RunTest("cognitoidp", "ListUserPools_MaxResultsOverLimit_Rejected", func() error {
		_, err := client.ListUserPools(ctx, &cognitoidentityprovider.ListUserPoolsInput{
			MaxResults: aws.Int32(61),
		})
		return AssertErrorContains(err, "InvalidParameterException")
	}))

	// ListUsers.Limit is QueryLimitType (range 0-60); 61 must be rejected.
	results = append(results, r.RunTest("cognitoidp", "ListUsers_LimitOverLimit_Rejected", func() error {
		pools, err := client.ListUserPools(ctx, &cognitoidentityprovider.ListUserPoolsInput{MaxResults: aws.Int32(1)})
		if err != nil {
			return fmt.Errorf("list pools: %v", err)
		}
		poolID := "us-east-1_NONEXISTENT"
		if len(pools.UserPools) > 0 && pools.UserPools[0].Id != nil {
			poolID = *pools.UserPools[0].Id
		}
		_, err = client.ListUsers(ctx, &cognitoidentityprovider.ListUsersInput{
			UserPoolId: aws.String(poolID),
			Limit:      aws.Int32(61),
		})
		return AssertErrorContains(err, "InvalidParameterException")
	}))

	return results
}

// cognitoPoolValidationNegativeTests pins the model-derived range and
// pattern validation on the pool and client creation paths.
func (r *TestRunner) cognitoPoolValidationNegativeTests(ctx context.Context, client *cognitoidentityprovider.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cognito", "CreateUserPool_InvalidPasswordPolicyMinimumLength", func() error {
		_, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(fmt.Sprintf("invalid-policy-%d", time.Now().UnixNano())),
			Policies: &types.UserPoolPolicyType{
				PasswordPolicy: &types.PasswordPolicyType{
					MinimumLength: aws.Int32(5),
				},
			},
		})
		if err := AssertErrorContains(err, "InvalidParameterException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "CreateUserPool_InvalidTemporaryPasswordValidityDays", func() error {
		_, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(fmt.Sprintf("invalid-days-%d", time.Now().UnixNano())),
			Policies: &types.UserPoolPolicyType{
				PasswordPolicy: &types.PasswordPolicyType{
					MinimumLength:                 aws.Int32(8),
					TemporaryPasswordValidityDays: 366,
				},
			},
		})
		if err := AssertErrorContains(err, "InvalidParameterException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "CreateUserPool_InvalidLambdaConfigArn", func() error {
		_, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(fmt.Sprintf("invalid-lambda-%d", time.Now().UnixNano())),
			LambdaConfig: &types.LambdaConfigType{
				PreSignUp: aws.String("not-an-arn"),
			},
		})
		if err := AssertErrorContains(err, "InvalidParameterException"); err != nil {
			return err
		}
		return nil
	}))

	// Standard attribute names are part of the schema name domain: AWS
	// documents modifying standard attribute properties through the
	// CreateUserPool Schema parameter, including the 21-character
	// phone_number_verified. Only custom attribute names carry the
	// 1-20 character constraint.
	results = append(results, r.RunTest("cognito", "CreateUserPool_StandardSchemaAttributes", func() error {
		schemaPoolName := fmt.Sprintf("standard-schema-%d", time.Now().UnixNano())
		createResp, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(schemaPoolName),
			Schema: []types.SchemaAttributeType{
				{
					Name:              aws.String("email"),
					AttributeDataType: types.AttributeDataTypeString,
					Required:          aws.Bool(true),
				},
				{
					Name:              aws.String("phone_number_verified"),
					AttributeDataType: types.AttributeDataTypeBoolean,
					Mutable:           aws.Bool(true),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("standard schema attributes rejected: %v", err)
		}
		defer client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: createResp.UserPool.Id,
		})

		_, err = client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(schemaPoolName + "-custom"),
			Schema: []types.SchemaAttributeType{
				{
					Name:              aws.String("customnameexceedinglimit"),
					AttributeDataType: types.AttributeDataTypeString,
				},
			},
		})
		return AssertErrorContains(err, "InvalidParameterException")
	}))

	poolName := fmt.Sprintf("client-validation-%d", time.Now().UnixNano())
	var poolID string
	results = append(results, r.RunTest("cognito", "CreateUserPoolClient_InvalidClientName", func() error {
		createResp, err := client.CreateUserPool(ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(poolName),
		})
		if err != nil {
			return fmt.Errorf("create pool: %v", err)
		}
		poolID = *createResp.UserPool.Id

		_, err = client.CreateUserPoolClient(ctx, &cognitoidentityprovider.CreateUserPoolClientInput{
			UserPoolId: aws.String(poolID),
			ClientName: aws.String("bad:name"),
		})
		if err := AssertErrorContains(err, "InvalidParameterException"); err != nil {
			return err
		}
		return nil
	}))
	if poolID != "" {
		results = append(results, r.RunTest("cognito", "UpdateUserPool_InvalidTemporaryPasswordValidityDays", func() error {
			_, err := client.UpdateUserPool(ctx, &cognitoidentityprovider.UpdateUserPoolInput{
				UserPoolId: aws.String(poolID),
				Policies: &types.UserPoolPolicyType{
					PasswordPolicy: &types.PasswordPolicyType{
						MinimumLength:                 aws.Int32(8),
						TemporaryPasswordValidityDays: 366,
					},
				},
			})
			if err := AssertErrorContains(err, "InvalidParameterException"); err != nil {
				return err
			}
			_, _ = client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: aws.String(poolID),
			})
			return nil
		}))
	}

	return results
}

package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func (r *TestRunner) cognitoEdgeCaseTests(tc *cognitoIDPContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cognito", "DescribeUserPool_NonExistent", func() error {
		_, err := tc.client.DescribeUserPool(tc.ctx, &cognitoidentityprovider.DescribeUserPoolInput{
			UserPoolId: aws.String(fmt.Sprintf("%s_nonexistentpool", r.region)),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "DeleteUserPool_NonExistent", func() error {
		_, err := tc.client.DeleteUserPool(tc.ctx, &cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: aws.String(fmt.Sprintf("%s_nonexistentpool", r.region)),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "AdminGetUser_NonExistent", func() error {
		poolID, cleanupPool, err := tc.createUserPool(tc.unique("err-pool"))
		if err != nil {
			return fmt.Errorf("create pool: %v", err)
		}
		defer cleanupPool()
		_, err = tc.client.AdminGetUser(tc.ctx, &cognitoidentityprovider.AdminGetUserInput{
			UserPoolId: aws.String(poolID),
			Username:   aws.String("nonexistent-user-xyz"),
		})
		if err := AssertErrorContains(err, "UserNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "CreateUserPool_DuplicateName", func() error {
		dupPoolName := tc.unique("dup-pool")
		pool1, cleanup1, err := tc.createUserPool(dupPoolName)
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer cleanup1()
		pool2, cleanup2, err := tc.createUserPool(dupPoolName)
		if err != nil {
			return fmt.Errorf("duplicate name should be allowed (unique IDs), got: %v", err)
		}
		if pool2 == "" || pool2 == pool1 {
			return fmt.Errorf("duplicate pool should have different ID")
		}
		defer cleanup2()
		return nil
	}))

	results = append(results, r.RunTest("cognito", "AdminCreateUser_VerifyAttributes", func() error {
		poolID, cleanupPool, err := tc.createUserPool(tc.unique("attr-pool"))
		if err != nil {
			return fmt.Errorf("create pool: %v", err)
		}
		defer cleanupPool()
		attrUser := tc.unique("attr-user")
		createUserResp, err := tc.client.AdminCreateUser(tc.ctx, &cognitoidentityprovider.AdminCreateUserInput{
			UserPoolId:        aws.String(poolID),
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
		poolID, cleanupPool, err := tc.createUserPool(tc.unique("list-pool"))
		if err != nil {
			return fmt.Errorf("create pool: %v", err)
		}
		defer cleanupPool()
		listUser := tc.unique("list-user")
		_, err = tc.client.AdminCreateUser(tc.ctx, &cognitoidentityprovider.AdminCreateUserInput{
			UserPoolId:        aws.String(poolID),
			Username:          aws.String(listUser),
			TemporaryPassword: aws.String("TempPass123!"),
			MessageAction:     types.MessageActionTypeSuppress,
		})
		if err != nil {
			return fmt.Errorf("create user: %v", err)
		}
		resp, err := tc.client.ListUsers(tc.ctx, &cognitoidentityprovider.ListUsersInput{
			UserPoolId: aws.String(poolID),
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
		poolID, cleanupPool, err := tc.createUserPool(tc.unique("grp-pool"))
		if err != nil {
			return fmt.Errorf("create pool: %v", err)
		}
		defer cleanupPool()
		testGroup := tc.unique("test-grp")
		_, err = tc.client.CreateGroup(tc.ctx, &cognitoidentityprovider.CreateGroupInput{
			GroupName:   aws.String(testGroup),
			UserPoolId:  aws.String(poolID),
			Description: aws.String("Test group description"),
		})
		if err != nil {
			return fmt.Errorf("create group: %v", err)
		}
		resp, err := tc.client.ListGroups(tc.ctx, &cognitoidentityprovider.ListGroupsInput{
			UserPoolId: aws.String(poolID),
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
		poolID, cleanupPool, err := tc.createUserPool(tc.unique("ge-pool"))
		if err != nil {
			return fmt.Errorf("create pool: %v", err)
		}
		defer cleanupPool()
		_, err = tc.client.GetGroup(tc.ctx, &cognitoidentityprovider.GetGroupInput{
			GroupName:  aws.String("nonexistent-group-xyz"),
			UserPoolId: aws.String(poolID),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "DescribeIdentityProvider_NonExistent", func() error {
		poolID, cleanupPool, err := tc.createUserPool(tc.unique("dip-pool"))
		if err != nil {
			return fmt.Errorf("create pool: %v", err)
		}
		defer cleanupPool()
		_, err = tc.client.DescribeIdentityProvider(tc.ctx, &cognitoidentityprovider.DescribeIdentityProviderInput{
			UserPoolId:   aws.String(poolID),
			ProviderName: aws.String("nonexistent-idp-xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "DescribeResourceServer_NonExistent", func() error {
		poolID, cleanupPool, err := tc.createUserPool(tc.unique("drs-pool"))
		if err != nil {
			return fmt.Errorf("create pool: %v", err)
		}
		defer cleanupPool()
		_, err = tc.client.DescribeResourceServer(tc.ctx, &cognitoidentityprovider.DescribeResourceServerInput{
			UserPoolId: aws.String(poolID),
			Identifier: aws.String("nonexistent-rs-xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "DeleteIdentityProvider_NonExistent", func() error {
		poolID, cleanupPool, err := tc.createUserPool(tc.unique("dlip-pool"))
		if err != nil {
			return fmt.Errorf("create pool: %v", err)
		}
		defer cleanupPool()
		_, err = tc.client.DeleteIdentityProvider(tc.ctx, &cognitoidentityprovider.DeleteIdentityProviderInput{
			UserPoolId:   aws.String(poolID),
			ProviderName: aws.String("nonexistent-idp-xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "DeleteResourceServer_NonExistent", func() error {
		poolID, cleanupPool, err := tc.createUserPool(tc.unique("dlrs-pool"))
		if err != nil {
			return fmt.Errorf("create pool: %v", err)
		}
		defer cleanupPool()
		_, err = tc.client.DeleteResourceServer(tc.ctx, &cognitoidentityprovider.DeleteResourceServerInput{
			UserPoolId: aws.String(poolID),
			Identifier: aws.String("nonexistent-rs-xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito", "ListUserPools_Pagination", func() error {
		var cleanupPools []func()
		defer func() {
			for _, cleanup := range cleanupPools {
				cleanup()
			}
		}()
		for i := 0; i < 5; i++ {
			_, cleanupPool, err := tc.createUserPool(fmt.Sprintf("PagPool-%s-%d", tc.ts, i))
			if err != nil {
				return err
			}
			cleanupPools = append(cleanupPools, cleanupPool)
		}

		pageCount := 0
		pools, err := paginate(func(next *string) ([]types.UserPoolDescriptionType, *string, error) {
			pageCount++
			resp, err := tc.client.ListUserPools(tc.ctx, &cognitoidentityprovider.ListUserPoolsInput{
				MaxResults: aws.Int32(2),
				NextToken:  next,
			})
			if err != nil {
				return nil, nil, err
			}
			return resp.UserPools, resp.NextToken, nil
		})
		if err != nil {
			return fmt.Errorf("list user pools page: %v", err)
		}
		if pageCount < 2 {
			return fmt.Errorf("expected at least 2 pages with MaxResults=2, got %d", pageCount)
		}
		sawOwn := 0
		for _, pool := range pools {
			if pool.Name != nil && strings.Contains(*pool.Name, "PagPool-"+tc.ts) {
				sawOwn++
			}
		}
		if sawOwn != 5 {
			return fmt.Errorf("expected all 5 paginated pools across pages, got %d", sawOwn)
		}
		return nil
	}))

	// MaxResults above the Smithy PoolQueryLimitType maximum of 60 must be
	// rejected with InvalidParameterException instead of being silently
	// clamped to the default.
	results = append(results, r.RunTest("cognito", "ListUserPools_MaxResultsOverLimit_Rejected", func() error {
		_, err := tc.client.ListUserPools(tc.ctx, &cognitoidentityprovider.ListUserPoolsInput{
			MaxResults: aws.Int32(61),
		})
		return AssertErrorContains(err, "InvalidParameterException")
	}))

	// ListUsers.Limit is QueryLimitType (range 0-60); 61 must be rejected.
	results = append(results, r.RunTest("cognito", "ListUsers_LimitOverLimit_Rejected", func() error {
		pools, err := tc.client.ListUserPools(tc.ctx, &cognitoidentityprovider.ListUserPoolsInput{MaxResults: aws.Int32(1)})
		if err != nil {
			return fmt.Errorf("list pools: %v", err)
		}
		poolID := "us-east-1_NONEXISTENT"
		if len(pools.UserPools) > 0 && pools.UserPools[0].Id != nil {
			poolID = *pools.UserPools[0].Id
		}
		_, err = tc.client.ListUsers(tc.ctx, &cognitoidentityprovider.ListUsersInput{
			UserPoolId: aws.String(poolID),
			Limit:      aws.Int32(61),
		})
		return AssertErrorContains(err, "InvalidParameterException")
	}))

	return results
}

// cognitoPoolValidationNegativeTests pins the model-derived range and
// pattern validation on the pool and client creation paths.
func (r *TestRunner) cognitoPoolValidationNegativeTests(tc *cognitoIDPContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("cognito", "CreateUserPool_InvalidPasswordPolicyMinimumLength", func() error {
		_, err := tc.client.CreateUserPool(tc.ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(tc.unique("invalid-policy")),
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
		_, err := tc.client.CreateUserPool(tc.ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(tc.unique("invalid-days")),
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
		_, err := tc.client.CreateUserPool(tc.ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(tc.unique("invalid-lambda")),
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
		schemaPoolName := tc.unique("standard-schema")
		_, cleanupPool, err := tc.createUserPool(schemaPoolName,
			func(input *cognitoidentityprovider.CreateUserPoolInput) {
				input.Schema = []types.SchemaAttributeType{
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
				}
			})
		if err != nil {
			return fmt.Errorf("standard schema attributes rejected: %v", err)
		}
		defer cleanupPool()

		_, err = tc.client.CreateUserPool(tc.ctx, &cognitoidentityprovider.CreateUserPoolInput{
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

	// The schema projection pins the AWS wire contract: SchemaAttributes is
	// SDK-visible, carries the full standard attribute set with the pool's
	// settings, and prefixes custom attribute names.
	results = append(results, r.RunTest("cognito", "CreateUserPool_SchemaAttributesPrefixedOutput", func() error {
		createResp, err := tc.client.CreateUserPool(tc.ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(tc.unique("schema-output")),
			Schema: []types.SchemaAttributeType{
				{Name: aws.String("rank"), AttributeDataType: types.AttributeDataTypeString, Mutable: aws.Bool(true)},
				{Name: aws.String("email"), AttributeDataType: types.AttributeDataTypeString, Required: aws.Bool(true)},
			},
		})
		if err != nil {
			return fmt.Errorf("create pool with schema: %v", err)
		}
		if createResp.UserPool == nil || createResp.UserPool.Id == nil {
			return fmt.Errorf("UserPool missing from create response")
		}
		poolID := *createResp.UserPool.Id
		defer tc.client.DeleteUserPool(tc.ctx, &cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: aws.String(poolID),
		})

		assertSchema := func(attrs []types.SchemaAttributeType) error {
			byName := make(map[string]types.SchemaAttributeType, len(attrs))
			for _, a := range attrs {
				if a.Name == nil {
					return fmt.Errorf("schema attribute without a name")
				}
				byName[*a.Name] = a
			}
			for _, want := range []string{"sub", "email", "phone_number_verified", "updated_at", "identities", "custom:rank"} {
				if _, ok := byName[want]; !ok {
					return fmt.Errorf("SchemaAttributes missing %q (have %v)", want, byName)
				}
			}
			if sub := byName["sub"]; sub.Required == nil || !*sub.Required {
				return fmt.Errorf("sub must be a required attribute")
			}
			if email := byName["email"]; email.Required == nil || !*email.Required {
				return fmt.Errorf("email must keep the supplied Required=true")
			}
			if _, ok := byName["rank"]; ok {
				return fmt.Errorf("custom attribute returned without the custom: prefix")
			}
			return nil
		}
		if err := assertSchema(createResp.UserPool.SchemaAttributes); err != nil {
			return fmt.Errorf("create response: %v", err)
		}

		descResp, err := tc.client.DescribeUserPool(tc.ctx, &cognitoidentityprovider.DescribeUserPoolInput{
			UserPoolId: aws.String(poolID),
		})
		if err != nil {
			return fmt.Errorf("describe pool: %v", err)
		}
		if descResp.UserPool == nil {
			return fmt.Errorf("UserPool missing from describe response")
		}
		return assertSchema(descResp.UserPool.SchemaAttributes)
	}))

	// The CSV header carries the documented base columns (including
	// cognito:mfa_enabled) plus the pool's custom attribute columns with
	// the custom: prefix, without duplicating standard attribute names.
	results = append(results, r.RunTest("cognito", "GetCSVHeader_CustomAttributeColumns", func() error {
		poolID, cleanupPool, err := tc.createUserPool(tc.unique("csv-header"),
			func(input *cognitoidentityprovider.CreateUserPoolInput) {
				input.Schema = []types.SchemaAttributeType{
					{Name: aws.String("rank"), AttributeDataType: types.AttributeDataTypeString, Mutable: aws.Bool(true)},
					{Name: aws.String("email"), AttributeDataType: types.AttributeDataTypeString, Required: aws.Bool(true)},
				}
			})
		if err != nil {
			return fmt.Errorf("create pool with schema: %v", err)
		}
		defer cleanupPool()

		resp, err := tc.client.GetCSVHeader(tc.ctx, &cognitoidentityprovider.GetCSVHeaderInput{
			UserPoolId: aws.String(poolID),
		})
		if err != nil {
			return err
		}
		emailColumns, hasMfa, hasUsername, hasRank := 0, false, false, false
		for _, column := range resp.CSVHeader {
			switch column {
			case "email":
				emailColumns++
			case "cognito:mfa_enabled":
				hasMfa = true
			case "cognito:username":
				hasUsername = true
			case "custom:rank":
				hasRank = true
			}
		}
		if !hasMfa || !hasUsername {
			return fmt.Errorf("base columns missing (mfa_enabled=%v username=%v): %v", hasMfa, hasUsername, resp.CSVHeader)
		}
		if !hasRank {
			return fmt.Errorf("custom attribute column custom:rank missing: %v", resp.CSVHeader)
		}
		if emailColumns != 1 {
			return fmt.Errorf("standard attribute column duplicated: email appears %d times: %v", emailColumns, resp.CSVHeader)
		}
		return nil
	}))

	// UpdateUserPool carries no Schema member in the model; an update must
	// never drop custom attributes added through AddCustomAttributes.
	results = append(results, r.RunTest("cognito", "UpdateUserPool_PreservesSchemaAttributes", func() error {
		poolID, cleanupPool, err := tc.createUserPool(tc.unique("schema-preserve"))
		if err != nil {
			return fmt.Errorf("create pool: %v", err)
		}
		defer cleanupPool()

		if _, err := tc.client.AddCustomAttributes(tc.ctx, &cognitoidentityprovider.AddCustomAttributesInput{
			UserPoolId: aws.String(poolID),
			CustomAttributes: []types.SchemaAttributeType{
				{Name: aws.String("loyalty"), AttributeDataType: types.AttributeDataTypeString, Mutable: aws.Bool(true)},
			},
		}); err != nil {
			return fmt.Errorf("add custom attributes: %v", err)
		}
		if _, err := tc.client.UpdateUserPool(tc.ctx, &cognitoidentityprovider.UpdateUserPoolInput{
			UserPoolId: aws.String(poolID),
			Policies: &types.UserPoolPolicyType{
				PasswordPolicy: &types.PasswordPolicyType{
					MinimumLength: aws.Int32(10),
				},
			},
		}); err != nil {
			return fmt.Errorf("update pool: %v", err)
		}

		descResp, err := tc.client.DescribeUserPool(tc.ctx, &cognitoidentityprovider.DescribeUserPoolInput{
			UserPoolId: aws.String(poolID),
		})
		if err != nil {
			return fmt.Errorf("describe pool: %v", err)
		}
		if descResp.UserPool == nil {
			return fmt.Errorf("UserPool missing from describe response")
		}
		for _, a := range descResp.UserPool.SchemaAttributes {
			if a.Name != nil && *a.Name == "custom:loyalty" {
				return nil
			}
		}
		return fmt.Errorf("custom:loyalty missing from SchemaAttributes after update")
	}))

	poolName := tc.unique("client-validation")
	var poolID string
	results = append(results, r.RunTest("cognito", "CreateUserPoolClient_InvalidClientName", func() error {
		createResp, err := tc.client.CreateUserPool(tc.ctx, &cognitoidentityprovider.CreateUserPoolInput{
			PoolName: aws.String(poolName),
		})
		if err != nil {
			return fmt.Errorf("create pool: %v", err)
		}
		poolID = *createResp.UserPool.Id

		_, err = tc.client.CreateUserPoolClient(tc.ctx, &cognitoidentityprovider.CreateUserPoolClientInput{
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
			_, err := tc.client.UpdateUserPool(tc.ctx, &cognitoidentityprovider.UpdateUserPoolInput{
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
			_, _ = tc.client.DeleteUserPool(tc.ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: aws.String(poolID),
			})
			return nil
		}))
	}

	return results
}

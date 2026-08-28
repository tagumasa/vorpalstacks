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

	tc := &cognitoIDPContext{
		r:      r,
		ctx:    ctx,
		client: client,
		ts:     fmt.Sprintf("%d", time.Now().UnixNano()),
	}

	userPoolName := tc.unique("test-pool")
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
		tc.userPoolID = *resp.UserPool.Id
		return nil
	}))

	if tc.userPoolID != "" {
		results = append(results, r.cognitoPoolCoreTests(tc)...)
		results = append(results, r.cognitoClientTests(tc)...)
		results = append(results, r.cognitoGroupTests(tc)...)
		results = append(results, r.cognitoUserTests(tc)...)
		results = append(results, r.cognitoWebAuthnTests(tc)...)
		results = append(results, r.cognitoIDPTests(tc)...)
		results = append(results, r.cognitoResourceServerTests(tc)...)
		results = append(results, r.cognitoSRPTests(tc)...)
		results = append(results, r.cognitoReplicaTermsBrandingTests(tc)...)

		userPoolID := tc.userPoolID
		results = append(results, r.RunTest("cognito", "DeleteUserPool", func() error {
			_, err := client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: aws.String(userPoolID),
			})
			return err
		}))
	}

	results = append(results, r.cognitoEdgeCaseTests(tc)...)
	results = append(results, r.cognitoPoolValidationNegativeTests(tc)...)
	results = append(results, r.cognitoImportJobTests(tc)...)

	return results
}

// cognitoIDPContext carries the user-pool suite state: the idp client,
// the lifecycle pool created by the CreateUserPool test (empty in the
// edge/negative contexts that own their pools), and a unique per-suite
// suffix for naming ephemeral resources.
type cognitoIDPContext struct {
	r          *TestRunner
	ctx        context.Context
	client     *cognitoidentityprovider.Client
	userPoolID string
	ts         string
}

// unique returns a per-suite unique resource name.
func (tc *cognitoIDPContext) unique(prefix string) string {
	return prefix + "-" + tc.ts
}

// createUserPool creates a throwaway user pool and returns its ID plus a
// cleanup closure deleting it. opts mutate the create input for tests
// needing Schema, AutoVerifiedAttributes, or Policies beyond PoolName.
func (tc *cognitoIDPContext) createUserPool(name string, opts ...func(*cognitoidentityprovider.CreateUserPoolInput)) (string, func(), error) {
	input := &cognitoidentityprovider.CreateUserPoolInput{PoolName: aws.String(name)}
	for _, opt := range opts {
		opt(input)
	}
	resp, err := tc.client.CreateUserPool(tc.ctx, input)
	if err != nil {
		return "", func() {}, fmt.Errorf("create user pool %s: %w", name, err)
	}
	if resp.UserPool == nil || resp.UserPool.Id == nil {
		return "", func() {}, fmt.Errorf("create user pool %s: UserPool.Id is nil", name)
	}
	poolID := *resp.UserPool.Id
	return poolID, func() {
		_, _ = tc.client.DeleteUserPool(tc.ctx, &cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: aws.String(poolID),
		})
	}, nil
}

// createPoolClient creates a throwaway app client on the given pool and
// returns its ID plus a cleanup closure deleting it.
func (tc *cognitoIDPContext) createPoolClient(poolID, name string, opts ...func(*cognitoidentityprovider.CreateUserPoolClientInput)) (string, func(), error) {
	input := &cognitoidentityprovider.CreateUserPoolClientInput{
		UserPoolId: aws.String(poolID),
		ClientName: aws.String(name),
	}
	for _, opt := range opts {
		opt(input)
	}
	resp, err := tc.client.CreateUserPoolClient(tc.ctx, input)
	if err != nil {
		return "", func() {}, fmt.Errorf("create app client %s: %w", name, err)
	}
	if resp.UserPoolClient == nil || resp.UserPoolClient.ClientId == nil {
		return "", func() {}, fmt.Errorf("create app client %s: ClientId is nil", name)
	}
	clientID := *resp.UserPoolClient.ClientId
	return clientID, func() {
		_, _ = tc.client.DeleteUserPoolClient(tc.ctx, &cognitoidentityprovider.DeleteUserPoolClientInput{
			ClientId:   aws.String(clientID),
			UserPoolId: aws.String(poolID),
		})
	}, nil
}

// adminCreateUser creates a suppressed-invitation throwaway user with the
// standard temporary password on the suite pool and returns a cleanup
// closure deleting it. opts may override MessageAction (empty string
// restores the unsuppressed invitation behaviour).
func (tc *cognitoIDPContext) adminCreateUser(username string, opts ...func(*cognitoidentityprovider.AdminCreateUserInput)) (func(), error) {
	input := &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId:        aws.String(tc.userPoolID),
		Username:          aws.String(username),
		TemporaryPassword: aws.String("TempPass123!"),
		MessageAction:     types.MessageActionTypeSuppress,
	}
	for _, opt := range opts {
		opt(input)
	}
	if _, err := tc.client.AdminCreateUser(tc.ctx, input); err != nil {
		return func() {}, err
	}
	return func() {
		_, _ = tc.client.AdminDeleteUser(tc.ctx, &cognitoidentityprovider.AdminDeleteUserInput{
			UserPoolId: aws.String(tc.userPoolID),
			Username:   aws.String(username),
		})
	}, nil
}

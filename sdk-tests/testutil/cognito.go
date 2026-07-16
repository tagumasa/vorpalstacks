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
		results = append(results, r.cognitoPoolCoreTests(ctx, client, userPoolID)...)
		results = append(results, r.cognitoClientTests(ctx, client, userPoolID)...)
		results = append(results, r.cognitoGroupTests(ctx, client, userPoolID)...)
		results = append(results, r.cognitoUserTests(ctx, client, userPoolID)...)
		results = append(results, r.cognitoIDPTests(ctx, client, userPoolID)...)
		results = append(results, r.cognitoResourceServerTests(ctx, client, userPoolID)...)

		results = append(results, r.RunTest("cognito", "DeleteUserPool", func() error {
			_, err := client.DeleteUserPool(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: aws.String(userPoolID),
			})
			return err
		}))
	}

	results = append(results, r.cognitoEdgeCaseTests(ctx, client)...)

	return results
}

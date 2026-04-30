package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunCognitoIdentityTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "cognito-identity",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := cognitoidentity.NewFromConfig(cfg)
	ctx := context.Background()

	poolName := fmt.Sprintf("test-idpool-%d", time.Now().UnixNano())

	var poolID string
	results = append(results, r.RunTest("cognito-identity", "CreateIdentityPool", func() error {
		resp, err := client.CreateIdentityPool(ctx, &cognitoidentity.CreateIdentityPoolInput{
			IdentityPoolName:               aws.String(poolName),
			AllowUnauthenticatedIdentities: true,
		})
		if err != nil {
			return err
		}
		if resp.IdentityPoolId == nil || *resp.IdentityPoolId == "" {
			return fmt.Errorf("IdentityPoolId is nil or empty")
		}
		if resp.IdentityPoolName == nil || *resp.IdentityPoolName != poolName {
			return fmt.Errorf("IdentityPoolName mismatch: got %v, want %s", resp.IdentityPoolName, poolName)
		}
		if !resp.AllowUnauthenticatedIdentities {
			return fmt.Errorf("expected AllowUnauthenticatedIdentities true")
		}
		poolID = *resp.IdentityPoolId
		return nil
	}))

	if poolID != "" {
		results = append(results, r.cognitoIdentityPoolTests(ctx, client, poolID)...)
		results = append(results, r.cognitoIdentityRolesTests(ctx, client, poolID)...)

		idResults, identityID := r.cognitoIdentityIdTests(ctx, client, poolID)
		results = append(results, idResults...)

		results = append(results, r.cognitoIdentityCredentialsTests(ctx, client, poolID, identityID)...)
		results = append(results, r.cognitoIdentityDeveloperTests(ctx, client, poolID, identityID)...)
		results = append(results, r.cognitoIdentityTagsTests(ctx, client, poolID)...)

		results = append(results, r.RunTest("cognito-identity", "DeleteIdentityPool", func() error {
			_, err := client.DeleteIdentityPool(ctx, &cognitoidentity.DeleteIdentityPoolInput{
				IdentityPoolId: aws.String(poolID),
			})
			if err != nil {
				return err
			}
			_, err = client.DescribeIdentityPool(ctx, &cognitoidentity.DescribeIdentityPoolInput{
				IdentityPoolId: aws.String(poolID),
			})
			if err == nil {
				return fmt.Errorf("expected error for deleted pool")
			}
			return nil
		}))
	}

	results = append(results, r.cognitoIdentityEdgeTests(ctx, client)...)

	return results
}

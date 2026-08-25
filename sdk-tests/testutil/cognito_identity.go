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

	tc := &cognitoIdentityContext{
		r:      r,
		ctx:    ctx,
		client: client,
		ts:     fmt.Sprintf("%d", time.Now().UnixNano()),
	}

	poolName := tc.unique("test-idpool")

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
		tc.poolID = *resp.IdentityPoolId
		return nil
	}))

	if tc.poolID != "" {
		results = append(results, r.cognitoIdentityPoolTests(tc)...)
		results = append(results, r.cognitoIdentityRolesTests(tc)...)

		idResults := r.cognitoIdentityIdTests(tc)
		results = append(results, idResults...)

		results = append(results, r.cognitoIdentityCredentialsTests(tc)...)
		results = append(results, r.cognitoIdentityDeveloperTests(tc)...)
		results = append(results, r.cognitoIdentityTagsTests(tc)...)

		poolID := tc.poolID
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

	results = append(results, r.cognitoIdentityEdgeTests(tc)...)

	return results
}

// cognitoIdentityContext carries the identity-pool suite state: the
// identity client, the lifecycle pool created by the CreateIdentityPool
// test (empty in the edge context that owns its pools), the identity ID
// produced by the GetId stage, and a unique per-suite suffix.
type cognitoIdentityContext struct {
	r          *TestRunner
	ctx        context.Context
	client     *cognitoidentity.Client
	poolID     string
	identityID string
	ts         string
}

// unique returns a per-suite unique resource name.
func (tc *cognitoIdentityContext) unique(prefix string) string {
	return prefix + "-" + tc.ts
}

// createIdPool creates a throwaway identity pool with unauthenticated
// identities allowed by default and returns its ID plus a cleanup
// closure deleting it.
func (tc *cognitoIdentityContext) createIdPool(name string, opts ...func(*cognitoidentity.CreateIdentityPoolInput)) (string, func(), error) {
	input := &cognitoidentity.CreateIdentityPoolInput{
		IdentityPoolName:               aws.String(name),
		AllowUnauthenticatedIdentities: true,
	}
	for _, opt := range opts {
		opt(input)
	}
	resp, err := tc.client.CreateIdentityPool(tc.ctx, input)
	if err != nil {
		return "", func() {}, fmt.Errorf("create identity pool %s: %w", name, err)
	}
	poolID := aws.ToString(resp.IdentityPoolId)
	return poolID, func() {
		_, _ = tc.client.DeleteIdentityPool(tc.ctx, &cognitoidentity.DeleteIdentityPoolInput{
			IdentityPoolId: aws.String(poolID),
		})
	}, nil
}

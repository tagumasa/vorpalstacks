package testutil

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
)

func (r *TestRunner) cognitoIdentityIdTests(ctx context.Context, client *cognitoidentity.Client, poolID string) ([]TestResult, string) {
	var results []TestResult

	results = append(results, r.RunTest("cognito-identity", "GetId", func() error {
		resp, err := client.GetId(ctx, &cognitoidentity.GetIdInput{
			IdentityPoolId: aws.String(poolID),
		})
		if err != nil {
			return err
		}
		if resp.IdentityId == nil || *resp.IdentityId == "" {
			return fmt.Errorf("IdentityId is nil or empty")
		}
		return nil
	}))

	var identityID string
	results = append(results, r.RunTest("cognito-identity", "GetId_WithLogins", func() error {
		resp, err := client.GetId(ctx, &cognitoidentity.GetIdInput{
			IdentityPoolId: aws.String(poolID),
			Logins: map[string]string{
				"graph.facebook.com": "test-token",
			},
		})
		if err != nil {
			return err
		}
		if resp.IdentityId == nil || *resp.IdentityId == "" {
			return fmt.Errorf("IdentityId is nil or empty")
		}
		identityID = *resp.IdentityId
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "DescribeIdentity", func() error {
		resp, err := client.DescribeIdentity(ctx, &cognitoidentity.DescribeIdentityInput{
			IdentityId: aws.String(identityID),
		})
		if err != nil {
			return err
		}
		if resp.IdentityId == nil || *resp.IdentityId != identityID {
			return fmt.Errorf("expected identity ID %s, got %v", identityID, resp.IdentityId)
		}
		found := false
		for _, l := range resp.Logins {
			if l == "graph.facebook.com" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("expected graph.facebook.com in Logins")
		}
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "ListIdentities", func() error {
		resp, err := client.ListIdentities(ctx, &cognitoidentity.ListIdentitiesInput{
			IdentityPoolId: aws.String(poolID),
			MaxResults:     aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if *resp.IdentityPoolId != poolID {
			return fmt.Errorf("expected pool ID %s, got %s", poolID, *resp.IdentityPoolId)
		}
		if len(resp.Identities) < 1 {
			return fmt.Errorf("expected at least 1 identity, got %d", len(resp.Identities))
		}
		return nil
	}))

	var secondIdentityID string
	results = append(results, r.RunTest("cognito-identity", "GetId_SecondIdentity", func() error {
		resp, err := client.GetId(ctx, &cognitoidentity.GetIdInput{
			IdentityPoolId: aws.String(poolID),
			Logins: map[string]string{
				"accounts.google.com": "google-token",
			},
		})
		if err != nil {
			return err
		}
		if resp.IdentityId == nil || *resp.IdentityId == "" {
			return fmt.Errorf("IdentityId is nil or empty")
		}
		secondIdentityID = *resp.IdentityId
		return nil
	}))

	results = append(results, r.RunTest("cognito-identity", "DeleteIdentities", func() error {
		resp, err := client.DeleteIdentities(ctx, &cognitoidentity.DeleteIdentitiesInput{
			IdentityIdsToDelete: []string{secondIdentityID},
		})
		if err != nil {
			return err
		}
		if len(resp.UnprocessedIdentityIds) > 0 {
			return fmt.Errorf("unexpected unprocessed identities: %v", resp.UnprocessedIdentityIds)
		}
		_, err = client.DescribeIdentity(ctx, &cognitoidentity.DescribeIdentityInput{
			IdentityId: aws.String(secondIdentityID),
		})
		if err == nil {
			return fmt.Errorf("expected error for deleted identity")
		}
		return nil
	}))

	return results, identityID
}

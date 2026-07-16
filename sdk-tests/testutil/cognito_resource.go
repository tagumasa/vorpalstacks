package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func (r *TestRunner) cognitoResourceServerTests(ctx context.Context, client *cognitoidentityprovider.Client, userPoolID string) []TestResult {
	var results []TestResult

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
			return fmt.Errorf("identifier mismatch: got %v, want %s", resp.ResourceServer.Identifier, identifier)
		}
		if resp.ResourceServer.Name == nil || *resp.ResourceServer.Name != "Test Resource Server" {
			return fmt.Errorf("name mismatch: got %v, want Test Resource Server", resp.ResourceServer.Name)
		}
		if len(resp.ResourceServer.Scopes) != 0 {
			return fmt.Errorf("scopes should be empty when no scopes specified")
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
		found := false
		for _, rs := range resp.ResourceServers {
			if rs.Identifier != nil && *rs.Identifier == identifier {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("resource server %s not found in ListResourceServers", identifier)
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
			return fmt.Errorf("identifier mismatch: got %v", resp.ResourceServer.Identifier)
		}
		if resp.ResourceServer.Name == nil || *resp.ResourceServer.Name != "Test Resource Server" {
			return fmt.Errorf("name mismatch: got %v", resp.ResourceServer.Name)
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
			return fmt.Errorf("name not updated: got %v", resp.ResourceServer.Name)
		}
		if len(resp.ResourceServer.Scopes) != 1 {
			return fmt.Errorf("expected 1 scope, got %d", len(resp.ResourceServer.Scopes))
		}
		if resp.ResourceServer.Scopes[0].ScopeName == nil || *resp.ResourceServer.Scopes[0].ScopeName != "read" {
			return fmt.Errorf("scope name mismatch: got %v", resp.ResourceServer.Scopes[0].ScopeName)
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
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}

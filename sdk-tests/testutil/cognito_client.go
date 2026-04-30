package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

func (r *TestRunner) cognitoClientTests(ctx context.Context, client *cognitoidentityprovider.Client, userPoolID string) []TestResult {
	var results []TestResult

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
		found := false
		for _, c := range resp.UserPoolClients {
			if c.ClientId != nil && *c.ClientId == clientID {
				found = true
				if c.ClientName == nil || *c.ClientName == "" {
					return fmt.Errorf("ClientName is nil or empty in listing")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("created client %s not found in ListUserPoolClients", clientID)
		}
		return nil
	}))

	if clientID != "" {
		results = append(results, r.RunTest("cognito", "DeleteUserPoolClient", func() error {
			_, err := client.DeleteUserPoolClient(ctx, &cognitoidentityprovider.DeleteUserPoolClientInput{
				ClientId:   aws.String(clientID),
				UserPoolId: aws.String(userPoolID),
			})
			if err != nil {
				return err
			}
			_, err = client.DescribeUserPoolClient(ctx, &cognitoidentityprovider.DescribeUserPoolClientInput{
				ClientId:   aws.String(clientID),
				UserPoolId: aws.String(userPoolID),
			})
			if err == nil {
				return fmt.Errorf("expected error describing deleted client")
			}
			return nil
		}))
	}

	return results
}

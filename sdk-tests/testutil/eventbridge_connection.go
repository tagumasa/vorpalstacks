package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

func (r *TestRunner) runEventBridgeConnectionTests(ctx context.Context, client *eventbridge.Client) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("events", "CreateConnection", func() error {
		ccName := fmt.Sprintf("CcConn-%d", time.Now().UnixNano())
		resp, err := client.CreateConnection(ctx, &eventbridge.CreateConnectionInput{
			Name:              aws.String(ccName),
			AuthorizationType: types.ConnectionAuthorizationTypeBasic,
			AuthParameters: &types.CreateConnectionAuthRequestParameters{
				BasicAuthParameters: &types.CreateConnectionBasicAuthRequestParameters{
					Username: aws.String("testuser"),
					Password: aws.String("testpass"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create connection: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.ConnectionArn == nil || *resp.ConnectionArn == "" {
			return fmt.Errorf("connection ARN is nil or empty")
		}
		defer client.DeleteConnection(ctx, &eventbridge.DeleteConnectionInput{Name: aws.String(ccName)})
		return nil
	}))

	results = append(results, r.RunTest("events", "DescribeConnection", func() error {
		dcName := fmt.Sprintf("DcConn-%d", time.Now().UnixNano())
		_, err := client.CreateConnection(ctx, &eventbridge.CreateConnectionInput{
			Name:              aws.String(dcName),
			AuthorizationType: types.ConnectionAuthorizationTypeBasic,
			AuthParameters: &types.CreateConnectionAuthRequestParameters{
				BasicAuthParameters: &types.CreateConnectionBasicAuthRequestParameters{
					Username: aws.String("testuser"),
					Password: aws.String("testpass"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create connection: %v", err)
		}
		defer client.DeleteConnection(ctx, &eventbridge.DeleteConnectionInput{Name: aws.String(dcName)})

		resp, err := client.DescribeConnection(ctx, &eventbridge.DescribeConnectionInput{
			Name: aws.String(dcName),
		})
		if err != nil {
			return fmt.Errorf("describe connection: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.Name == nil || *resp.Name != dcName {
			return fmt.Errorf("connection name mismatch, got %v", resp.Name)
		}
		if resp.ConnectionArn == nil || *resp.ConnectionArn == "" {
			return fmt.Errorf("connection ARN is nil or empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "DeleteConnection", func() error {
		dlcName := fmt.Sprintf("DlcConn-%d", time.Now().UnixNano())
		_, err := client.CreateConnection(ctx, &eventbridge.CreateConnectionInput{
			Name:              aws.String(dlcName),
			AuthorizationType: types.ConnectionAuthorizationTypeBasic,
			AuthParameters: &types.CreateConnectionAuthRequestParameters{
				BasicAuthParameters: &types.CreateConnectionBasicAuthRequestParameters{
					Username: aws.String("u"),
					Password: aws.String("p"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create connection: %v", err)
		}

		resp, err := client.DeleteConnection(ctx, &eventbridge.DeleteConnectionInput{
			Name: aws.String(dlcName),
		})
		if err != nil {
			return fmt.Errorf("delete connection: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}

		_, err = client.DescribeConnection(ctx, &eventbridge.DescribeConnectionInput{
			Name: aws.String(dlcName),
		})
		if err == nil {
			return fmt.Errorf("expected error describing deleted connection")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "ListConnections", func() error {
		lcName := fmt.Sprintf("LcConn-%d", time.Now().UnixNano())
		_, err := client.CreateConnection(ctx, &eventbridge.CreateConnectionInput{
			Name:              aws.String(lcName),
			AuthorizationType: types.ConnectionAuthorizationTypeBasic,
			AuthParameters: &types.CreateConnectionAuthRequestParameters{
				BasicAuthParameters: &types.CreateConnectionBasicAuthRequestParameters{
					Username: aws.String("u"),
					Password: aws.String("p"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create connection: %v", err)
		}
		defer client.DeleteConnection(ctx, &eventbridge.DeleteConnectionInput{Name: aws.String(lcName)})

		resp, err := client.ListConnections(ctx, &eventbridge.ListConnectionsInput{})
		if err != nil {
			return fmt.Errorf("list connections: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.Connections == nil {
			return fmt.Errorf("connections list is nil")
		}
		found := false
		for _, c := range resp.Connections {
			if c.Name != nil && *c.Name == lcName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("expected connection %s in list", lcName)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "UpdateConnection", func() error {
		ucName := fmt.Sprintf("UcConn-%d", time.Now().UnixNano())
		_, err := client.CreateConnection(ctx, &eventbridge.CreateConnectionInput{
			Name:              aws.String(ucName),
			AuthorizationType: types.ConnectionAuthorizationTypeApiKey,
			AuthParameters: &types.CreateConnectionAuthRequestParameters{
				ApiKeyAuthParameters: &types.CreateConnectionApiKeyAuthRequestParameters{
					ApiKeyName:  aws.String("key"),
					ApiKeyValue: aws.String("value"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create connection: %v", err)
		}
		defer client.DeleteConnection(ctx, &eventbridge.DeleteConnectionInput{Name: aws.String(ucName)})

		resp, err := client.UpdateConnection(ctx, &eventbridge.UpdateConnectionInput{
			Name:        aws.String(ucName),
			Description: aws.String("updated connection description"),
		})
		if err != nil {
			return fmt.Errorf("update connection: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.ConnectionArn == nil || *resp.ConnectionArn == "" {
			return fmt.Errorf("connection ARN is nil or empty")
		}
		desc, err := client.DescribeConnection(ctx, &eventbridge.DescribeConnectionInput{
			Name: aws.String(ucName),
		})
		if err != nil {
			return fmt.Errorf("describe after update: %v", err)
		}
		if desc.Description == nil || *desc.Description != "updated connection description" {
			return fmt.Errorf("description not updated, got %v", desc.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "CreateApiDestination", func() error {
		cadName := fmt.Sprintf("CadDest-%d", time.Now().UnixNano())
		connName := fmt.Sprintf("CadConn-%d", time.Now().UnixNano())
		_, err := client.CreateConnection(ctx, &eventbridge.CreateConnectionInput{
			Name:              aws.String(connName),
			AuthorizationType: types.ConnectionAuthorizationTypeBasic,
			AuthParameters: &types.CreateConnectionAuthRequestParameters{
				BasicAuthParameters: &types.CreateConnectionBasicAuthRequestParameters{
					Username: aws.String("u"),
					Password: aws.String("p"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create connection: %v", err)
		}
		defer client.DeleteConnection(ctx, &eventbridge.DeleteConnectionInput{Name: aws.String(connName)})

		resp, err := client.CreateApiDestination(ctx, &eventbridge.CreateApiDestinationInput{
			Name:               aws.String(cadName),
			ConnectionArn:      aws.String(fmt.Sprintf("arn:aws:events:%s:000000000000:connection/%s", r.region, connName)),
			HttpMethod:         types.ApiDestinationHttpMethodPost,
			InvocationEndpoint: aws.String("https://example.com/webhook"),
			Description:        aws.String("test api destination"),
		})
		if err != nil {
			return fmt.Errorf("create api destination: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.ApiDestinationArn == nil || *resp.ApiDestinationArn == "" {
			return fmt.Errorf("api destination ARN is nil or empty")
		}
		defer client.DeleteApiDestination(ctx, &eventbridge.DeleteApiDestinationInput{Name: aws.String(cadName)})
		return nil
	}))

	results = append(results, r.RunTest("events", "DescribeApiDestination", func() error {
		dadName := fmt.Sprintf("DadDest-%d", time.Now().UnixNano())
		dadConn := fmt.Sprintf("DadConn-%d", time.Now().UnixNano())
		_, err := client.CreateConnection(ctx, &eventbridge.CreateConnectionInput{
			Name:              aws.String(dadConn),
			AuthorizationType: types.ConnectionAuthorizationTypeBasic,
			AuthParameters: &types.CreateConnectionAuthRequestParameters{
				BasicAuthParameters: &types.CreateConnectionBasicAuthRequestParameters{
					Username: aws.String("u"),
					Password: aws.String("p"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create connection: %v", err)
		}
		defer client.DeleteConnection(ctx, &eventbridge.DeleteConnectionInput{Name: aws.String(dadConn)})

		connARN := fmt.Sprintf("arn:aws:events:%s:000000000000:connection/%s", r.region, dadConn)
		_, err = client.CreateApiDestination(ctx, &eventbridge.CreateApiDestinationInput{
			Name:               aws.String(dadName),
			ConnectionArn:      aws.String(connARN),
			HttpMethod:         types.ApiDestinationHttpMethodPost,
			InvocationEndpoint: aws.String("https://example.com/webhook"),
			Description:        aws.String("test api destination for describe"),
		})
		if err != nil {
			return fmt.Errorf("create api destination: %v", err)
		}
		defer client.DeleteApiDestination(ctx, &eventbridge.DeleteApiDestinationInput{Name: aws.String(dadName)})

		resp, err := client.DescribeApiDestination(ctx, &eventbridge.DescribeApiDestinationInput{
			Name: aws.String(dadName),
		})
		if err != nil {
			return fmt.Errorf("describe api destination: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.Name == nil || *resp.Name != dadName {
			return fmt.Errorf("name mismatch, got %v", resp.Name)
		}
		if resp.ConnectionArn == nil || *resp.ConnectionArn != connARN {
			return fmt.Errorf("connection ARN mismatch, got %v", resp.ConnectionArn)
		}
		if resp.HttpMethod != types.ApiDestinationHttpMethodPost {
			return fmt.Errorf("http method mismatch, got %v", resp.HttpMethod)
		}
		if resp.InvocationEndpoint == nil || *resp.InvocationEndpoint != "https://example.com/webhook" {
			return fmt.Errorf("invocation endpoint mismatch, got %v", resp.InvocationEndpoint)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "DeleteApiDestination", func() error {
		dladName := fmt.Sprintf("DladDest-%d", time.Now().UnixNano())
		dladConn := fmt.Sprintf("DladConn-%d", time.Now().UnixNano())
		_, err := client.CreateConnection(ctx, &eventbridge.CreateConnectionInput{
			Name:              aws.String(dladConn),
			AuthorizationType: types.ConnectionAuthorizationTypeBasic,
			AuthParameters: &types.CreateConnectionAuthRequestParameters{
				BasicAuthParameters: &types.CreateConnectionBasicAuthRequestParameters{
					Username: aws.String("u"),
					Password: aws.String("p"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create connection: %v", err)
		}
		defer client.DeleteConnection(ctx, &eventbridge.DeleteConnectionInput{Name: aws.String(dladConn)})

		connARN := fmt.Sprintf("arn:aws:events:%s:000000000000:connection/%s", r.region, dladConn)
		_, err = client.CreateApiDestination(ctx, &eventbridge.CreateApiDestinationInput{
			Name:               aws.String(dladName),
			ConnectionArn:      aws.String(connARN),
			HttpMethod:         types.ApiDestinationHttpMethodPost,
			InvocationEndpoint: aws.String("https://example.com/webhook"),
		})
		if err != nil {
			return fmt.Errorf("create api destination: %v", err)
		}

		resp, err := client.DeleteApiDestination(ctx, &eventbridge.DeleteApiDestinationInput{
			Name: aws.String(dladName),
		})
		if err != nil {
			return fmt.Errorf("delete api destination: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}

		_, err = client.DescribeApiDestination(ctx, &eventbridge.DescribeApiDestinationInput{
			Name: aws.String(dladName),
		})
		if err == nil {
			return fmt.Errorf("expected error describing deleted api destination")
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "ListApiDestinations", func() error {
		ladName := fmt.Sprintf("LadDest-%d", time.Now().UnixNano())
		ladConn := fmt.Sprintf("LadConn-%d", time.Now().UnixNano())
		_, err := client.CreateConnection(ctx, &eventbridge.CreateConnectionInput{
			Name:              aws.String(ladConn),
			AuthorizationType: types.ConnectionAuthorizationTypeBasic,
			AuthParameters: &types.CreateConnectionAuthRequestParameters{
				BasicAuthParameters: &types.CreateConnectionBasicAuthRequestParameters{
					Username: aws.String("u"),
					Password: aws.String("p"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create connection: %v", err)
		}
		defer client.DeleteConnection(ctx, &eventbridge.DeleteConnectionInput{Name: aws.String(ladConn)})

		connARN := fmt.Sprintf("arn:aws:events:%s:000000000000:connection/%s", r.region, ladConn)
		_, err = client.CreateApiDestination(ctx, &eventbridge.CreateApiDestinationInput{
			Name:               aws.String(ladName),
			ConnectionArn:      aws.String(connARN),
			HttpMethod:         types.ApiDestinationHttpMethodPost,
			InvocationEndpoint: aws.String("https://example.com/webhook"),
		})
		if err != nil {
			return fmt.Errorf("create api destination: %v", err)
		}
		defer client.DeleteApiDestination(ctx, &eventbridge.DeleteApiDestinationInput{Name: aws.String(ladName)})

		resp, err := client.ListApiDestinations(ctx, &eventbridge.ListApiDestinationsInput{})
		if err != nil {
			return fmt.Errorf("list api destinations: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.ApiDestinations == nil {
			return fmt.Errorf("api destinations list is nil")
		}
		found := false
		for _, d := range resp.ApiDestinations {
			if d.Name != nil && *d.Name == ladName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("expected api destination %s in list", ladName)
		}
		return nil
	}))

	results = append(results, r.RunTest("events", "UpdateApiDestination", func() error {
		uadName := fmt.Sprintf("UadDest-%d", time.Now().UnixNano())
		uadConn := fmt.Sprintf("UadConn-%d", time.Now().UnixNano())
		_, err := client.CreateConnection(ctx, &eventbridge.CreateConnectionInput{
			Name:              aws.String(uadConn),
			AuthorizationType: types.ConnectionAuthorizationTypeBasic,
			AuthParameters: &types.CreateConnectionAuthRequestParameters{
				BasicAuthParameters: &types.CreateConnectionBasicAuthRequestParameters{
					Username: aws.String("u"),
					Password: aws.String("p"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create connection: %v", err)
		}
		defer client.DeleteConnection(ctx, &eventbridge.DeleteConnectionInput{Name: aws.String(uadConn)})

		connARN := fmt.Sprintf("arn:aws:events:%s:000000000000:connection/%s", r.region, uadConn)
		_, err = client.CreateApiDestination(ctx, &eventbridge.CreateApiDestinationInput{
			Name:               aws.String(uadName),
			ConnectionArn:      aws.String(connARN),
			HttpMethod:         types.ApiDestinationHttpMethodPost,
			InvocationEndpoint: aws.String("https://example.com/original"),
		})
		if err != nil {
			return fmt.Errorf("create api destination: %v", err)
		}
		defer client.DeleteApiDestination(ctx, &eventbridge.DeleteApiDestinationInput{Name: aws.String(uadName)})

		resp, err := client.UpdateApiDestination(ctx, &eventbridge.UpdateApiDestinationInput{
			Name:               aws.String(uadName),
			Description:        aws.String("updated description"),
			InvocationEndpoint: aws.String("https://example.com/updated"),
		})
		if err != nil {
			return fmt.Errorf("update api destination: %v", err)
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		if resp.ApiDestinationArn == nil || *resp.ApiDestinationArn == "" {
			return fmt.Errorf("api destination ARN is nil or empty")
		}
		desc, err := client.DescribeApiDestination(ctx, &eventbridge.DescribeApiDestinationInput{
			Name: aws.String(uadName),
		})
		if err != nil {
			return fmt.Errorf("describe after update: %v", err)
		}
		if desc.InvocationEndpoint == nil || *desc.InvocationEndpoint != "https://example.com/updated" {
			return fmt.Errorf("invocation endpoint not updated, got %v", desc.InvocationEndpoint)
		}
		return nil
	}))

	return results
}

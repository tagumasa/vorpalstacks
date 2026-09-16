package testutil

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
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
		_, cleanupConn, err := createEventBridgeTestConnection(ctx, client, dcName, func(input *eventbridge.CreateConnectionInput) {
			input.AuthParameters.BasicAuthParameters.Username = aws.String("testuser")
			input.AuthParameters.BasicAuthParameters.Password = aws.String("testpass")
		})
		if err != nil {
			return err
		}
		defer cleanupConn()

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
		// The deletion below is the operation under test, so the helper's
		// cleanup is discarded.
		_, _, err := createEventBridgeTestConnection(ctx, client, dlcName)
		if err != nil {
			return err
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
		_, cleanupConn, err := createEventBridgeTestConnection(ctx, client, lcName)
		if err != nil {
			return err
		}
		defer cleanupConn()

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
		_, cleanupConn, err := createEventBridgeTestConnection(ctx, client, ucName, func(input *eventbridge.CreateConnectionInput) {
			input.AuthorizationType = types.ConnectionAuthorizationTypeApiKey
			input.AuthParameters = &types.CreateConnectionAuthRequestParameters{
				ApiKeyAuthParameters: &types.CreateConnectionApiKeyAuthRequestParameters{
					ApiKeyName:  aws.String("key"),
					ApiKeyValue: aws.String("value"),
				},
			}
		})
		if err != nil {
			return err
		}
		defer cleanupConn()

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
		connARN, cleanupConn, err := createEventBridgeTestConnection(ctx, client, connName)
		if err != nil {
			return err
		}
		defer cleanupConn()

		resp, err := client.CreateApiDestination(ctx, &eventbridge.CreateApiDestinationInput{
			Name:               aws.String(cadName),
			ConnectionArn:      aws.String(connARN),
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
		connARN, cleanupConn, err := createEventBridgeTestConnection(ctx, client, dadConn)
		if err != nil {
			return err
		}
		defer cleanupConn()

		cleanupDest, err := createEventBridgeTestApiDestination(ctx, client, dadName, connARN, func(input *eventbridge.CreateApiDestinationInput) {
			input.Description = aws.String("test api destination for describe")
		})
		if err != nil {
			return err
		}
		defer cleanupDest()

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
		connARN, cleanupConn, err := createEventBridgeTestConnection(ctx, client, dladConn)
		if err != nil {
			return err
		}
		defer cleanupConn()

		// The deletion below is the operation under test, so the helper's
		// cleanup is discarded.
		_, err = createEventBridgeTestApiDestination(ctx, client, dladName, connARN)
		if err != nil {
			return err
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
		connARN, cleanupConn, err := createEventBridgeTestConnection(ctx, client, ladConn)
		if err != nil {
			return err
		}
		defer cleanupConn()

		cleanupDest, err := createEventBridgeTestApiDestination(ctx, client, ladName, connARN)
		if err != nil {
			return err
		}
		defer cleanupDest()

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
		connARN, cleanupConn, err := createEventBridgeTestConnection(ctx, client, uadConn)
		if err != nil {
			return err
		}
		defer cleanupConn()

		cleanupDest, err := createEventBridgeTestApiDestination(ctx, client, uadName, connARN, func(input *eventbridge.CreateApiDestinationInput) {
			input.InvocationEndpoint = aws.String("https://example.com/original")
		})
		if err != nil {
			return err
		}
		defer cleanupDest()

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

	// The connection wire-shape closure test: the full AuthParameters tree
	// (header/query/body parameter lists of {IsValueSecret, Key, Value})
	// survives the SDK round trip, the credential members the response
	// shapes omit stay absent, and the credential secret is surfaced as
	// SecretArn.
	results = append(results, r.RunTest("events", "DescribeConnection_AuthParametersRoundTrip", func() error {
		rtName := fmt.Sprintf("RtConn-%d", time.Now().UnixNano())
		connARN, cleanupConn, err := createEventBridgeTestConnection(ctx, client, rtName, func(input *eventbridge.CreateConnectionInput) {
			input.AuthParameters.InvocationHttpParameters = &types.ConnectionHttpParameters{
				HeaderParameters: []types.ConnectionHeaderParameter{
					{Key: aws.String("X-Plain"), Value: aws.String("plain-value"), IsValueSecret: false},
					{Key: aws.String("X-Marked"), Value: aws.String("stored-blanked"), IsValueSecret: true},
				},
				QueryStringParameters: []types.ConnectionQueryStringParameter{
					{Key: aws.String("qs"), Value: aws.String("qv"), IsValueSecret: false},
				},
				BodyParameters: []types.ConnectionBodyParameter{
					{Key: aws.String("bk"), Value: aws.String("bv"), IsValueSecret: false},
				},
			}
		})
		if err != nil {
			return err
		}
		defer cleanupConn()

		resp, err := client.DescribeConnection(ctx, &eventbridge.DescribeConnectionInput{Name: aws.String(rtName)})
		if err != nil {
			return fmt.Errorf("describe connection: %v", err)
		}
		if resp.SecretArn == nil || *resp.SecretArn == "" {
			return fmt.Errorf("SecretArn must be populated (the secret created from the authorization parameters)")
		}
		if resp.AuthParameters == nil || resp.AuthParameters.BasicAuthParameters == nil {
			return fmt.Errorf("basic auth parameters must be described")
		}
		if username := resp.AuthParameters.BasicAuthParameters.Username; username == nil || *username != "u" {
			return fmt.Errorf("Username must round-trip, got %v", username)
		}
		// Password has no field on the response struct at all — the
		// DescribeConnection response shape structurally omits it.
		hp := resp.AuthParameters.InvocationHttpParameters
		if hp == nil || len(hp.HeaderParameters) != 2 {
			return fmt.Errorf("both header parameters must be described, got %+v", hp)
		}
		plain := hp.HeaderParameters[0]
		if aws.ToString(plain.Key) != "X-Plain" || aws.ToString(plain.Value) != "plain-value" || plain.IsValueSecret {
			return fmt.Errorf("non-secret header parameter must round-trip exactly, got %+v", plain)
		}
		marked := hp.HeaderParameters[1]
		if aws.ToString(marked.Key) != "X-Marked" || !marked.IsValueSecret {
			return fmt.Errorf("secret-marked header parameter keeps its key and flag, got %+v", marked)
		}
		if aws.ToString(marked.Value) == "stored-blanked" {
			return fmt.Errorf("a secret-marked value must not come back in plaintext")
		}
		if hp.QueryStringParameters == nil || len(hp.QueryStringParameters) != 1 || aws.ToString(hp.QueryStringParameters[0].Key) != "qs" {
			return fmt.Errorf("query string parameters must round-trip, got %+v", hp.QueryStringParameters)
		}
		if hp.BodyParameters == nil || len(hp.BodyParameters) != 1 || aws.ToString(hp.BodyParameters[0].Key) != "bk" {
			return fmt.Errorf("body parameters must round-trip, got %+v", hp.BodyParameters)
		}
		_ = connARN
		return nil
	}))

	// The API destination delivery closure test: a rule target pointing at
	// an API destination invokes the endpoint over HTTPS with the event
	// payload, the connection's Basic credentials, the target's HTTP
	// parameters, and the fixed header set.
	results = append(results, r.RunTest("events", "ApiDestinationTargetDelivery", func() error {
		var mu sync.Mutex
		var requests []struct {
			method        string
			path          string
			query         string
			authorization string
			userAgent     string
			contentType   string
			custom        string
			body          string
		}
		endpoint := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			mu.Lock()
			requests = append(requests, struct {
				method        string
				path          string
				query         string
				authorization string
				userAgent     string
				contentType   string
				custom        string
				body          string
			}{r.Method, r.URL.Path, r.URL.RawQuery, r.Header.Get("Authorization"),
				r.Header.Get("User-Agent"), r.Header.Get("Content-Type"),
				r.Header.Get("X-SDK-Test"), string(body)})
			mu.Unlock()
			w.WriteHeader(http.StatusOK)
		}))
		defer endpoint.Close()

		busName := fmt.Sprintf("api-dest-bus-%d", time.Now().UnixNano())
		cleanupBus, err := createEventBridgeTestBus(ctx, client, busName)
		if err != nil {
			return err
		}
		defer cleanupBus()

		connName := fmt.Sprintf("ApiDestConn-%d", time.Now().UnixNano())
		connARN, cleanupConn, err := createEventBridgeTestConnection(ctx, client, connName)
		if err != nil {
			return err
		}

		destName := fmt.Sprintf("DeliverDest-%d", time.Now().UnixNano())
		dest, err := client.CreateApiDestination(ctx, &eventbridge.CreateApiDestinationInput{
			Name:               aws.String(destName),
			ConnectionArn:      aws.String(connARN),
			HttpMethod:         types.ApiDestinationHttpMethodPost,
			InvocationEndpoint: aws.String(endpoint.URL + "/hook"),
		})
		if err != nil {
			return fmt.Errorf("create api destination: %v", err)
		}
		cleanupDest := func() {
			_, _ = client.DeleteApiDestination(ctx, &eventbridge.DeleteApiDestinationInput{Name: aws.String(destName)})
		}

		ruleName := fmt.Sprintf("api-dest-rule-%d", time.Now().UnixNano())
		_, err = client.PutRule(ctx, &eventbridge.PutRuleInput{
			Name:         aws.String(ruleName),
			EventBusName: aws.String(busName),
			EventPattern: aws.String(`{"source": ["sdktests.api-destination"]}`),
		})
		if err != nil {
			cleanupDest()
			return fmt.Errorf("put rule: %v", err)
		}

		_, err = client.PutTargets(ctx, &eventbridge.PutTargetsInput{
			Rule:         aws.String(ruleName),
			EventBusName: aws.String(busName),
			Targets: []types.Target{{
				Id:  aws.String("api-dest-1"),
				Arn: dest.ApiDestinationArn,
				HttpParameters: &types.HttpParameters{
					HeaderParameters:      map[string]string{"X-SDK-Test": "yes"},
					QueryStringParameters: map[string]string{"tag": "api-dest"},
				},
			}},
		})
		if err != nil {
			cleanupDest()
			return fmt.Errorf("put targets: %v", err)
		}
		cleanupTarget := func() {
			_, _ = client.RemoveTargets(ctx, &eventbridge.RemoveTargetsInput{
				Rule: aws.String(ruleName), EventBusName: aws.String(busName), Ids: []string{"api-dest-1"},
			})
			_, _ = client.DeleteRule(ctx, &eventbridge.DeleteRuleInput{Name: aws.String(ruleName), EventBusName: aws.String(busName)})
			cleanupDest()
		}
		defer cleanupTarget()

		if _, err := client.PutEvents(ctx, &eventbridge.PutEventsInput{
			Entries: []types.PutEventsRequestEntry{{
				EventBusName: aws.String(busName),
				Source:       aws.String("sdktests.api-destination"),
				DetailType:   aws.String("Delivery"),
				Detail:       aws.String(`{"kind":"hook"}`),
			}},
		}); err != nil {
			return fmt.Errorf("put events: %v", err)
		}

		if err := waitFor(200*time.Millisecond, 10*time.Second, func() bool {
			mu.Lock()
			defer mu.Unlock()
			return len(requests) > 0
		}); err != nil {
			cleanupConn()
			return fmt.Errorf("api destination target never invoked the endpoint")
		}
		cleanupConn()

		mu.Lock()
		req := requests[0]
		mu.Unlock()
		if req.method != "POST" || req.path != "/hook" {
			return fmt.Errorf("invocation method/path = %s %s, want POST /hook", req.method, req.path)
		}
		expectedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("u:p"))
		if req.authorization != expectedAuth {
			return fmt.Errorf("Authorization = %q, want %q", req.authorization, expectedAuth)
		}
		if req.userAgent != "Amazon/EventBridge/ApiDestinations" {
			return fmt.Errorf("User-Agent = %q", req.userAgent)
		}
		if !strings.HasPrefix(req.contentType, "application/json") {
			return fmt.Errorf("Content-Type = %q", req.contentType)
		}
		if req.custom != "yes" {
			return fmt.Errorf("target header X-SDK-Test = %q", req.custom)
		}
		if !strings.Contains(req.query, "tag=api-dest") {
			return fmt.Errorf("target query parameter missing, got %q", req.query)
		}
		if !strings.Contains(req.body, `"kind":"hook"`) && !strings.Contains(req.body, `"kind": "hook"`) {
			return fmt.Errorf("event payload must reach the endpoint, got %s", req.body)
		}
		return nil
	}))

	return results
}

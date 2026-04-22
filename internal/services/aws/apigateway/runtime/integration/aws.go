package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"vorpalstacks/internal/eventbus"
	"vorpalstacks/pkg/vtl"
)

var lambdaFunctionRegex = regexp.MustCompile(`functions/(.+?)/invocations`)

// jsonMarshal is a convenience wrapper for JSON marshalling.
func jsonMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// AWSExecutor executes AWS integration requests for API Gateway.
type AWSExecutor struct {
	accountID string
	region    string
	bus       eventbus.Bus
}

// NewAWSExecutor creates a new AWSExecutor using the event bus for
// cross-service invocations.
func NewAWSExecutor(bus eventbus.Bus, accountID, region string) *AWSExecutor {
	return &AWSExecutor{
		accountID: accountID,
		region:    region,
		bus:       bus,
	}
}

// Execute executes an AWS integration request.
func (e *AWSExecutor) Execute(ctx context.Context, req *IntegrationRequest) (*IntegrationResponse, error) {
	if isLambdaURI(req.URI) {
		return e.executeLambda(ctx, req)
	}

	if isSQSURI(req.URI) {
		return e.executeSQS(ctx, req)
	}

	if isSNSURI(req.URI) {
		return e.executeSNS(ctx, req)
	}

	if isDynamoDBURI(req.URI) {
		return e.executeDynamoDB(ctx, req)
	}

	if isKinesisURI(req.URI) {
		return e.executeKinesis(ctx, req)
	}

	if isSFNURI(req.URI) {
		return e.executeStepFunctions(ctx, req)
	}

	return nil, &IntegrationError{
		Message:  "AWS integration not yet implemented for: " + req.URI,
		Type:     "InternalServerError",
		HTTPCode: http.StatusNotImplemented,
	}
}

func isLambdaURI(uri string) bool {
	return strings.Contains(uri, "lambda:")
}

// LambdaProxyResponse represents the response from a Lambda proxy integration.
type LambdaProxyResponse struct {
	StatusCode      int               `json:"statusCode"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
}

// LambdaProxyEvent represents the event structure passed to a Lambda function
// when using API Gateway Lambda proxy integration.
type LambdaProxyEvent struct {
	Resource                        string                 `json:"resource"`
	Path                            string                 `json:"path"`
	HTTPMethod                      string                 `json:"httpMethod"`
	Headers                         map[string]string      `json:"headers"`
	MultiValueHeaders               map[string][]string    `json:"multiValueHeaders"`
	QueryStringParameters           map[string]string      `json:"queryStringParameters"`
	MultiValueQueryStringParameters map[string][]string    `json:"multiValueQueryStringParameters"`
	PathParameters                  map[string]string      `json:"pathParameters"`
	StageVariables                  map[string]string      `json:"stageVariables"`
	RequestContext                  map[string]interface{} `json:"requestContext"`
	Body                            string                 `json:"body"`
	IsBase64Encoded                 bool                   `json:"isBase64Encoded"`
}

func (e *AWSExecutor) executeLambda(ctx context.Context, req *IntegrationRequest) (*IntegrationResponse, error) {
	if e.bus == nil || e.bus.LambdaInvoker() == nil {
		return nil, &IntegrationError{
			Message:  "Lambda client not configured",
			Type:     "InternalServerError",
			HTTPCode: 500,
		}
	}

	functionRef, err := extractFunctionRefFromURI(req.URI)
	if err != nil {
		return nil, &IntegrationError{
			Message:  fmt.Sprintf("Invalid Lambda URI: %v", err),
			Type:     "BadRequestException",
			HTTPCode: 400,
		}
	}

	req = applyRequestParameterMapping(req)

	isProxy := req.IntegrationType == "AWS_PROXY"

	var eventJSON []byte
	if isProxy {
		event := e.buildLambdaProxyEvent(req)
		eventJSON, err = json.Marshal(event)
		if err != nil {
			return nil, &IntegrationError{
				Message:  fmt.Sprintf("Failed to marshal Lambda event: %v", err),
				Type:     "InternalServerError",
				HTTPCode: 500,
			}
		}
	} else {
		contentType := req.Headers["Content-Type"]
		tmpl := req.RequestTemplates[contentType]
		if tmpl == "" {
			tmpl = req.RequestTemplates["application/json"]
		}
		if tmpl != "" {
			eventJSON, err = applyRequestTemplate(tmpl, req)
			if err != nil {
				return nil, &IntegrationError{
					Message:  fmt.Sprintf("Failed to apply request template: %v", err),
					Type:     "InternalServerError",
					HTTPCode: 500,
				}
			}
		} else if len(req.Body) > 0 {
			eventJSON = req.Body
		} else {
			eventJSON = []byte("{}")
		}
	}

	statusCode, payload, err := e.bus.LambdaInvoker().InvokeForGateway(ctx, functionRef, eventJSON)
	if err != nil {
		return nil, &IntegrationError{
			Message:  fmt.Sprintf("Lambda invocation failed: %v", err),
			Type:     "IntegrationFailure",
			HTTPCode: 502,
		}
	}

	if isProxy {
		lambdaResp, err := parseLambdaResponse(payload)
		if err != nil {
			return nil, &IntegrationError{
				Message:  fmt.Sprintf("Failed to parse Lambda response: %v", err),
				Type:     "InternalServerError",
				HTTPCode: 500,
			}
		}
		return lambdaResp, nil
	}

	resp := &IntegrationResponse{
		StatusCode:      int(statusCode),
		Headers:         make(map[string]string),
		Body:            payload,
		IsBase64Encoded: false,
	}

	if req.IntegrationResponses != nil {
		respConfig := matchIntegrationResponse(req.IntegrationResponses, string(payload), int(statusCode))
		if respConfig != nil {
			if respConfig.ResponseHeaders != nil {
				for k, v := range respConfig.ResponseHeaders {
					resp.Headers[k] = v
				}
			}

			if respConfig.ResponseTemplates != nil {
				contentType := "application/json"
				if tmpl, ok := respConfig.ResponseTemplates[contentType]; ok && tmpl != "" {
					transformed, err := applyResponseTemplate(tmpl, payload, req)
					if err != nil {
						return nil, &IntegrationError{
							Message:  fmt.Sprintf("Failed to apply response template: %v", err),
							Type:     "InternalServerError",
							HTTPCode: 500,
						}
					}
					resp.Body = transformed
				}
			}

			if respConfig.StatusCode != "" {
				_, _ = fmt.Sscanf(respConfig.StatusCode, "%d", &resp.StatusCode)
			}

			resp.Headers = applyResponseParameterMapping(respConfig.ResponseHeaders, resp.Headers, string(payload))
		}
	}

	return resp, nil
}

func (e *AWSExecutor) buildLambdaProxyEvent(req *IntegrationRequest) *LambdaProxyEvent {
	now := time.Now()
	requestID := fmt.Sprintf("%x", now.UnixNano())

	return &LambdaProxyEvent{
		Resource:                        req.Path,
		Path:                            req.Path,
		HTTPMethod:                      req.Method,
		Headers:                         req.Headers,
		MultiValueHeaders:               req.MultiValueHeaders,
		QueryStringParameters:           req.QueryParams,
		MultiValueQueryStringParameters: req.MultiValueQueryStringParams,
		PathParameters:                  req.PathParams,
		StageVariables:                  req.StageVariables,
		RequestContext: map[string]interface{}{
			"accountID":    e.accountID,
			"apiId":        req.RestApiId,
			"requestId":    requestID,
			"resourcePath": req.Path,
			"stage":        req.StageName,
			"region":       e.region,
			"http": map[string]interface{}{
				"method":    req.Method,
				"path":      req.Path,
				"protocol":  "HTTP/1.1",
				"sourceIp":  req.SourceIP,
				"userAgent": req.Headers["User-Agent"],
			},
			"requestTime":      now.Format("02/Jan/2006:15:04:05 -0700"),
			"requestTimeEpoch": now.UnixMilli(),
		},
		Body:            string(req.Body),
		IsBase64Encoded: false,
	}
}

// applyRequestTemplate applies a VTL request template with $input, $context,
// and $stageVariables substitutions.
func applyRequestTemplate(tmpl string, req *IntegrationRequest) ([]byte, error) {
	engine := vtl.NewEngine()

	engine.GetContext().Body = string(req.Body)
	if len(req.Body) > 0 {
		var bodyObj interface{}
		if err := json.Unmarshal(req.Body, &bodyObj); err == nil {
			engine.GetContext().JSONBody = bodyObj
		}
	}

	params := make(map[string]string)
	for k, v := range req.PathParams {
		params[k] = v
	}
	for k, v := range req.QueryParams {
		params[k] = v
	}
	for k, v := range req.Headers {
		params[k] = v
	}
	engine.GetContext().Params = params

	engine.GetContext().Context = map[string]interface{}{
		"stage":     req.StageName,
		"apiId":     req.RestApiId,
		"requestId": fmt.Sprintf("%x", time.Now().UnixNano()),
	}

	engine.GetContext().StageVars = req.StageVariables

	result, err := engine.Transform(tmpl)
	if err != nil {
		return nil, err
	}

	return []byte(result), nil
}

func matchIntegrationResponse(responses map[string]*IntegrationResponseConfig, responseBody string, statusCode int) *IntegrationResponseConfig {
	var defaultResp *IntegrationResponseConfig

	for _, resp := range responses {
		if resp.SelectionPattern == "" {
			if defaultResp == nil {
				defaultResp = resp
			}
			continue
		}

		matched, err := regexp.MatchString(resp.SelectionPattern, responseBody)
		if err == nil && matched {
			return resp
		}
	}

	if defaultResp != nil {
		return defaultResp
	}

	statusStr := fmt.Sprintf("%d", statusCode)
	return responses[statusStr]
}

// applyResponseTemplate applies a VTL response template with $input, $context,
// and $stageVariables substitutions.
func applyResponseTemplate(tmpl string, responseBody []byte, req *IntegrationRequest) ([]byte, error) {
	engine := vtl.NewEngine()

	engine.GetContext().Body = string(responseBody)
	if len(responseBody) > 0 {
		var bodyObj interface{}
		if err := json.Unmarshal(responseBody, &bodyObj); err == nil {
			engine.GetContext().JSONBody = bodyObj
		}
	}

	params := make(map[string]string)
	for k, v := range req.PathParams {
		params[k] = v
	}
	for k, v := range req.QueryParams {
		params[k] = v
	}
	for k, v := range req.Headers {
		params[k] = v
	}
	engine.GetContext().Params = params

	engine.GetContext().Context = map[string]interface{}{
		"stage":     req.StageName,
		"apiId":     req.RestApiId,
		"requestId": fmt.Sprintf("%x", time.Now().UnixNano()),
	}

	engine.GetContext().StageVars = req.StageVariables

	result, err := engine.Transform(tmpl)
	if err != nil {
		return nil, err
	}

	return []byte(result), nil
}

func extractFunctionRefFromURI(uri string) (string, error) {
	matches := lambdaFunctionRegex.FindStringSubmatch(uri)
	if len(matches) < 2 {
		return "", fmt.Errorf("invalid Lambda URI format")
	}
	return matches[1], nil
}

func parseLambdaResponse(body []byte) (*IntegrationResponse, error) {
	var lambdaResp LambdaProxyResponse
	if err := json.Unmarshal(body, &lambdaResp); err != nil {
		return nil, fmt.Errorf("failed to parse Lambda response: %w", err)
	}

	return &IntegrationResponse{
		StatusCode:      lambdaResp.StatusCode,
		Headers:         lambdaResp.Headers,
		Body:            []byte(lambdaResp.Body),
		IsBase64Encoded: lambdaResp.IsBase64Encoded,
	}, nil
}

package testutil

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	apigwtypes "github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

// runAPIGatewayTestInvokeToSQS verifies that TestInvokeMethod against an
// AWS-type SQS integration reaches the real backend through the wired event
// bus — the local executor switch previously built the AWS executor without
// a bus, so every AWS-integration test invocation failed with a
// "not configured" style error instead of invoking the service.
func (r *TestRunner) runAPIGatewayTestInvokeToSQS(ic *integClients, ts string) TestResult {
	const name = "APIGateway_TestInvoke_SQS"
	queueName := fmt.Sprintf("integ-apigw-ti-q-%s", ts)
	apiName := fmt.Sprintf("integ-apigw-ti-api-%s", ts)

	queueURL, err := ic.createQueue(queueName)
	if err != nil {
		return r.RunTest(integSvc, name, func() error { return fmt.Errorf("create queue: %w", err) })
	}
	defer ic.deleteQueue(queueURL)

	api, err := ic.apigateway.CreateRestApi(ic.ctx, &apigateway.CreateRestApiInput{
		Name: aws.String(apiName),
	})
	if err != nil {
		return r.RunTest(integSvc, name, func() error { return fmt.Errorf("create rest api: %w", err) })
	}
	apiID := aws.ToString(api.Id)
	defer ic.apigateway.DeleteRestApi(ic.ctx, &apigateway.DeleteRestApiInput{RestApiId: aws.String(apiID)})

	resource, err := ic.apigateway.CreateResource(ic.ctx, &apigateway.CreateResourceInput{
		RestApiId: api.Id, ParentId: api.RootResourceId, PathPart: aws.String("send"),
	})
	if err != nil {
		return r.RunTest(integSvc, name, func() error { return fmt.Errorf("create resource: %w", err) })
	}

	for _, call := range []struct {
		name string
		fn   func() error
	}{
		{"put method", func() error {
			_, err := ic.apigateway.PutMethod(ic.ctx, &apigateway.PutMethodInput{
				RestApiId: api.Id, ResourceId: resource.Id,
				HttpMethod: aws.String("POST"), AuthorizationType: aws.String("NONE"),
			})
			return err
		}},
		{"put integration", func() error {
			_, err := ic.apigateway.PutIntegration(ic.ctx, &apigateway.PutIntegrationInput{
				RestApiId: api.Id, ResourceId: resource.Id, HttpMethod: aws.String("POST"),
				Type:                  apigwtypes.IntegrationTypeAws,
				IntegrationHttpMethod: aws.String("POST"),
				Uri:                   aws.String(fmt.Sprintf("arn:aws:apigateway:%s:sqs:path/000000000000/%s", ic.region, queueName)),
			})
			return err
		}},
	} {
		if err := call.fn(); err != nil {
			return r.RunTest(integSvc, name, func() error { return fmt.Errorf("%s: %w", call.name, err) })
		}
	}

	body := `{"test":"apigw-testinvoke-sqs"}`
	resp, err := ic.apigateway.TestInvokeMethod(ic.ctx, &apigateway.TestInvokeMethodInput{
		RestApiId: api.Id, ResourceId: resource.Id, HttpMethod: aws.String("POST"),
		PathWithQueryString: aws.String("/send"),
		Body:                aws.String(body),
	})
	if err != nil {
		return r.RunTest(integSvc, name, func() error { return fmt.Errorf("test invoke: %w", err) })
	}
	if resp.Status >= 400 {
		return r.RunTest(integSvc, name, func() error {
			return fmt.Errorf("test invoke returned status %d, body %q, log %q", resp.Status, aws.ToString(resp.Body), aws.ToString(resp.Log))
		})
	}

	return r.pollVerify(name, defaultPollTimeout, func() error {
		return ic.verifyMessageContains(queueURL, "apigw-testinvoke-sqs")
	})
}

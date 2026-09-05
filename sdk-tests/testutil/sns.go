package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"vorpalstacks-sdk-tests/config"
)

type snsTestContext struct {
	client    *sns.Client
	ctx       context.Context
	region    string
	accountID string
}

func (r *TestRunner) initSNS() (*snsTestContext, error) {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &snsTestContext{
		client:    sns.NewFromConfig(cfg),
		ctx:       context.Background(),
		region:    r.region,
		accountID: r.accountID,
	}, nil
}

func (tc *snsTestContext) uniqueName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

func (tc *snsTestContext) createTopic(name string) (string, error) {
	resp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
		Name: ptrString(name),
	})
	if err != nil {
		return "", fmt.Errorf("create topic %q: %w", name, err)
	}
	return *resp.TopicArn, nil
}

func (tc *snsTestContext) deleteTopic(arn string) {
	tc.client.DeleteTopic(tc.ctx, &sns.DeleteTopicInput{
		TopicArn: ptrString(arn),
	})
}

// getTopicAttributes fetches the attribute map of one topic, in plain form:
// only the topic ARN, no optional members.
func (tc *snsTestContext) getTopicAttributes(topicArn string) (*sns.GetTopicAttributesOutput, error) {
	return tc.client.GetTopicAttributes(tc.ctx, &sns.GetTopicAttributesInput{
		TopicArn: aws.String(topicArn),
	})
}

// getSubscriptionAttributes fetches the attribute map of one subscription,
// in plain form: only the subscription ARN, no optional members.
func (tc *snsTestContext) getSubscriptionAttributes(subArn string) (*sns.GetSubscriptionAttributesOutput, error) {
	return tc.client.GetSubscriptionAttributes(tc.ctx, &sns.GetSubscriptionAttributesInput{
		SubscriptionArn: aws.String(subArn),
	})
}

// getPlatformApplicationAttributes fetches the attribute map of one platform
// application, in plain form: only the ARN, no optional members.
func (tc *snsTestContext) getPlatformApplicationAttributes(appArn string) (*sns.GetPlatformApplicationAttributesOutput, error) {
	return tc.client.GetPlatformApplicationAttributes(tc.ctx, &sns.GetPlatformApplicationAttributesInput{
		PlatformApplicationArn: aws.String(appArn),
	})
}

// getEndpointAttributes fetches the attribute map of one platform endpoint,
// in plain form: only the endpoint ARN, no optional members.
func (tc *snsTestContext) getEndpointAttributes(endpointArn string) (*sns.GetEndpointAttributesOutput, error) {
	return tc.client.GetEndpointAttributes(tc.ctx, &sns.GetEndpointAttributesInput{
		EndpointArn: aws.String(endpointArn),
	})
}

// listTags returns the tags attached to a topic ARN.
func (tc *snsTestContext) listTags(resourceArn string) ([]types.Tag, error) {
	resp, err := tc.client.ListTagsForResource(tc.ctx, &sns.ListTagsForResourceInput{
		ResourceArn: aws.String(resourceArn),
	})
	if err != nil {
		return nil, err
	}
	return resp.Tags, nil
}

// snsTagValue returns the value of the named tag and whether it is present.
func snsTagValue(tags []types.Tag, key string) (string, bool) {
	for _, t := range tags {
		if t.Key != nil && *t.Key == key {
			if t.Value == nil {
				return "", true
			}
			return *t.Value, true
		}
	}
	return "", false
}

// subscribeSQS subscribes a platform SQS queue ARN to a topic in plain form:
// protocol sqs, queue ARN built from the test account, no optional members.
// It returns the subscription ARN.
func (tc *snsTestContext) subscribeSQS(topicArn, queueName string) (string, error) {
	resp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
		TopicArn: aws.String(topicArn),
		Protocol: aws.String("sqs"),
		Endpoint: aws.String(fmt.Sprintf("arn:aws:sqs:%s:%s:%s", tc.region, tc.accountID, queueName)),
	})
	if err != nil {
		return "", fmt.Errorf("subscribe: %w", err)
	}
	return *resp.SubscriptionArn, nil
}

// expectNotFound asserts the typed NotFound contract of a failed call,
// tagging the failure with the operation name. SNS reports its
// NotFoundException type with the wire error code "NotFound".
func (tc *snsTestContext) expectNotFound(op string, err error) error {
	if aerr := expectAWSErrorCode(err, "NotFound"); aerr != nil {
		return fmt.Errorf("%s: %w", op, aerr)
	}
	return nil
}

// expectResourceNotFound asserts the typed ResourceNotFound contract of a
// failed tag operation, tagging the failure with the operation name. The
// SNS model maps the ResourceNotFoundException shape to the awsQueryError
// wire code "ResourceNotFound".
func (tc *snsTestContext) expectResourceNotFound(op string, err error) error {
	if aerr := expectAWSErrorCode(err, "ResourceNotFound"); aerr != nil {
		return fmt.Errorf("%s: %w", op, aerr)
	}
	return nil
}

// startTokenCaptureServer runs a local HTTP server that captures the
// confirmation token from the first subscription-notification body. AWS
// sends the token to the endpoint out-of-band; it is not exposed through
// GetSubscriptionAttributes.
func (tc *snsTestContext) startTokenCaptureServer() (*httptest.Server, <-chan string) {
	tokenCh := make(chan string, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var payload map[string]interface{}
		if json.Unmarshal(body, &payload) == nil {
			if t, ok := payload["Token"].(string); ok {
				select {
				case tokenCh <- t:
				default:
				}
			}
		}
		w.WriteHeader(http.StatusOK)
	}))
	return srv, tokenCh
}

// listAllTopics walks every ListTopics page and returns all topic ARNs.
func (tc *snsTestContext) listAllTopics() ([]string, error) {
	return paginate(func(next *string) ([]string, *string, error) {
		resp, err := tc.client.ListTopics(tc.ctx, &sns.ListTopicsInput{NextToken: next})
		if err != nil {
			return nil, nil, err
		}
		arns := make([]string, 0, len(resp.Topics))
		for _, t := range resp.Topics {
			arns = append(arns, aws.ToString(t.TopicArn))
		}
		return arns, resp.NextToken, nil
	})
}

// listAllSubscriptions walks every ListSubscriptions page and returns all
// subscriptions.
func (tc *snsTestContext) listAllSubscriptions() ([]types.Subscription, error) {
	return paginate(func(next *string) ([]types.Subscription, *string, error) {
		resp, err := tc.client.ListSubscriptions(tc.ctx, &sns.ListSubscriptionsInput{NextToken: next})
		if err != nil {
			return nil, nil, err
		}
		return resp.Subscriptions, resp.NextToken, nil
	})
}

// listAllSubscriptionsByTopic walks every ListSubscriptionsByTopic page and
// returns all subscriptions of one topic.
func (tc *snsTestContext) listAllSubscriptionsByTopic(topicArn string) ([]types.Subscription, error) {
	return paginate(func(next *string) ([]types.Subscription, *string, error) {
		resp, err := tc.client.ListSubscriptionsByTopic(tc.ctx, &sns.ListSubscriptionsByTopicInput{
			TopicArn:  aws.String(topicArn),
			NextToken: next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Subscriptions, resp.NextToken, nil
	})
}

// listAllPlatformApplications walks every ListPlatformApplications page and
// returns all platform application ARNs.
func (tc *snsTestContext) listAllPlatformApplications() ([]string, error) {
	return paginate(func(next *string) ([]string, *string, error) {
		resp, err := tc.client.ListPlatformApplications(tc.ctx, &sns.ListPlatformApplicationsInput{NextToken: next})
		if err != nil {
			return nil, nil, err
		}
		arns := make([]string, 0, len(resp.PlatformApplications))
		for _, a := range resp.PlatformApplications {
			arns = append(arns, aws.ToString(a.PlatformApplicationArn))
		}
		return arns, resp.NextToken, nil
	})
}

// listAllEndpointsByPlatformApplication walks every
// ListEndpointsByPlatformApplication page and returns all endpoint ARNs of
// one platform application.
func (tc *snsTestContext) listAllEndpointsByPlatformApplication(appArn string) ([]string, error) {
	return paginate(func(next *string) ([]string, *string, error) {
		resp, err := tc.client.ListEndpointsByPlatformApplication(tc.ctx, &sns.ListEndpointsByPlatformApplicationInput{
			PlatformApplicationArn: aws.String(appArn),
			NextToken:              next,
		})
		if err != nil {
			return nil, nil, err
		}
		arns := make([]string, 0, len(resp.Endpoints))
		for _, e := range resp.Endpoints {
			arns = append(arns, aws.ToString(e.EndpointArn))
		}
		return arns, resp.NextToken, nil
	})
}

// httpSubscription captures a confirmed http-protocol subscription created
// against a local test server.
type httpSubscription struct {
	TopicArn        string
	SubscriptionArn string
}

// createHTTPSubscription creates a topic with an http subscription against a
// local server, waits for the confirmation token, and confirms the
// subscription. The returned struct carries both ARNs for cleanup and
// assertions.
func (tc *snsTestContext) createHTTPSubscription(prefix string) (*httpSubscription, error) {
	tResp, err := tc.client.CreateTopic(tc.ctx, &sns.CreateTopicInput{
		Name: ptrString(tc.uniqueName(prefix)),
	})
	if err != nil {
		return nil, fmt.Errorf("create topic: %v", err)
	}

	srv, tokenCh := tc.startTokenCaptureServer()
	defer srv.Close()

	sResp, err := tc.client.Subscribe(tc.ctx, &sns.SubscribeInput{
		TopicArn:              tResp.TopicArn,
		Protocol:              aws.String("http"),
		Endpoint:              aws.String(srv.URL),
		ReturnSubscriptionArn: true,
	})
	if err != nil {
		tc.deleteTopic(*tResp.TopicArn)
		return nil, fmt.Errorf("subscribe: %v", err)
	}

	var token string
	select {
	case token = <-tokenCh:
	case <-time.After(10 * time.Second):
		tc.deleteTopic(*tResp.TopicArn)
		return nil, fmt.Errorf("timed out waiting for subscription confirmation")
	}

	if _, err := tc.client.ConfirmSubscription(tc.ctx, &sns.ConfirmSubscriptionInput{
		TopicArn: tResp.TopicArn,
		Token:    aws.String(token),
	}); err != nil {
		tc.deleteTopic(*tResp.TopicArn)
		return nil, fmt.Errorf("confirm: %v", err)
	}

	return &httpSubscription{
		TopicArn:        *tResp.TopicArn,
		SubscriptionArn: *sResp.SubscriptionArn,
	}, nil
}

func (r *TestRunner) RunSNSTests() []TestResult {
	tc, err := r.initSNS()
	if err != nil {
		return []TestResult{{
			Service:  "sns",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		}}
	}

	var results []TestResult

	results = append(results, r.runSNSTopicTests(tc)...)
	results = append(results, r.runSNSPublishTests(tc)...)
	results = append(results, r.runSNSSubscriptionTests(tc)...)
	results = append(results, r.runSNSPlatformTests(tc)...)
	results = append(results, r.runSNSTagTests(tc)...)
	results = append(results, r.runSNSPolicyTests(tc)...)
	results = append(results, r.runSNSEdgeTests(tc)...)

	return results
}

func ptrString(s string) *string { return &s }

func nanoTime() int64 { return time.Now().UnixNano() }

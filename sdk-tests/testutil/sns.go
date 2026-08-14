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

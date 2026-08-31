package sns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	snsstore "vorpalstacks/internal/store/aws/sns"

	"github.com/google/uuid"
)

// Subscribe creates a subscription to an SNS topic.
// https://docs.aws.amazon.com/sns/latest/api/API_Subscribe.html
func (s *SNSService) Subscribe(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.subscribeCore(store, reqCtx, SubscribeInput{
		TopicArn:              request.GetParamLowerFirst(req.Parameters, "TopicArn"),
		Protocol:              request.GetParamLowerFirst(req.Parameters, "Protocol"),
		Endpoint:              request.GetParamLowerFirst(req.Parameters, "Endpoint"),
		ReturnSubscriptionArn: request.GetParamLowerFirst(req.Parameters, "ReturnSubscriptionArn"),
		Attributes:            parseAttributes(req.Parameters),
	})
}

// Unsubscribe deletes a subscription.
// https://docs.aws.amazon.com/sns/latest/api/API_Unsubscribe.html
func (s *SNSService) Unsubscribe(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.unsubscribeCore(store, reqCtx, UnsubscribeInput{
		SubscriptionArn: request.GetParamLowerFirst(req.Parameters, "SubscriptionArn"),
	})
}

// ConfirmSubscription confirms a subscription request.
// https://docs.aws.amazon.com/sns/latest/api/API_ConfirmSubscription.html
func (s *SNSService) ConfirmSubscription(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.confirmSubscriptionCore(store, ConfirmSubscriptionInput{
		TopicArn:                  request.GetParamLowerFirst(req.Parameters, "TopicArn"),
		Token:                     request.GetParamLowerFirst(req.Parameters, "Token"),
		AuthenticateOnUnsubscribe: request.GetParamLowerFirst(req.Parameters, "AuthenticateOnUnsubscribe"),
	})
}

// GetSubscriptionAttributes returns the attributes of a subscription.
// https://docs.aws.amazon.com/sns/latest/api/API_GetSubscriptionAttributes.html
func (s *SNSService) GetSubscriptionAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.getSubscriptionAttributesCore(store, GetSubscriptionAttributesInput{
		SubscriptionArn: request.GetParamLowerFirst(req.Parameters, "SubscriptionArn"),
	})
}

// SetSubscriptionAttributes sets the attributes of a subscription.
// https://docs.aws.amazon.com/sns/latest/api/API_SetSubscriptionAttributes.html
func (s *SNSService) SetSubscriptionAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.setSubscriptionAttributesCore(store, SetSubscriptionAttributesInput{
		SubscriptionArn: request.GetParamLowerFirst(req.Parameters, "SubscriptionArn"),
		AttributeName:   request.GetParamLowerFirst(req.Parameters, "AttributeName"),
		AttributeValue:  request.GetParamLowerFirst(req.Parameters, "AttributeValue"),
	})
}

// ListSubscriptions lists the subscriptions.
// https://docs.aws.amazon.com/sns/latest/api/API_ListSubscriptions.html
func (s *SNSService) ListSubscriptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.listSubscriptionsCore(store, ListSubscriptionsInput{
		NextToken: pagination.GetMarker(req.Parameters, "NextToken"),
	})
}

// ListSubscriptionsByTopic lists the subscriptions by topic.
// https://docs.aws.amazon.com/sns/latest/api/API_ListSubscriptionsByTopic.html
func (s *SNSService) ListSubscriptionsByTopic(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.listSubscriptionsByTopicCore(store, ListSubscriptionsByTopicInput{
		TopicArn:  request.GetParamLowerFirst(req.Parameters, "TopicArn"),
		NextToken: pagination.GetMarker(req.Parameters, "NextToken"),
	})
}

// sendSubscriptionConfirmation sends a SubscriptionConfirmation message
// to an HTTP/HTTPS endpoint. This is best-effort: if the endpoint is
// unreachable the subscription simply stays in pending state until the
// subscriber retries.
func (s *SNSService) sendSubscriptionConfirmation(sub *snsstore.Subscription, region string) {
	if region == "" {
		region = "us-east-1"
	}

	payload := map[string]interface{}{
		"Type":         "SubscriptionConfirmation",
		"MessageId":    uuid.New().String(),
		"Token":        sub.ConfirmationToken,
		"TopicArn":     sub.TopicArn,
		"Message":      fmt.Sprintf("You have chosen to subscribe to the topic %s.\nTo confirm the subscription, visit the SubscribeURL included in this message.", sub.TopicArn),
		"SubscribeURL": fmt.Sprintf("https://sns.%s.amazonaws.com/?Action=ConfirmSubscription&TopicArn=%s&Token=%s", region, sub.TopicArn, sub.ConfirmationToken),
		"Timestamp":    time.Now().UTC().Format(time.RFC3339),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		logs.Warn("SNS: failed to marshal subscription confirmation",
			logs.String("subscriptionArn", sub.SubscriptionArn),
			logs.Err(err))
		return
	}

	req, err := http.NewRequest("POST", sub.Endpoint, bytes.NewReader(body))
	if err != nil {
		logs.Warn("SNS: failed to create confirmation request",
			logs.String("endpoint", sub.Endpoint),
			logs.Err(err))
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-amz-sns-message-type", "SubscriptionConfirmation")
	req.Header.Set("x-amz-sns-message-id", payload["MessageId"].(string))
	req.Header.Set("x-amz-sns-topic-arn", sub.TopicArn)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logs.Warn("SNS: failed to send subscription confirmation",
			logs.String("endpoint", sub.Endpoint),
			logs.String("subscriptionArn", sub.SubscriptionArn),
			logs.Err(err))
		return
	}
	defer resp.Body.Close()

	logs.Debug("SNS: subscription confirmation sent",
		logs.String("endpoint", sub.Endpoint),
		logs.Int("status", resp.StatusCode))
}

func buildSubscriptionList(items []*snsstore.Subscription) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(items))
	for _, sub := range items {
		subArn := sub.SubscriptionArn
		if sub.PendingConfirmation {
			subArn = "pending confirmation"
		}
		result = append(result, map[string]interface{}{
			"SubscriptionArn": subArn,
			"TopicArn":        sub.TopicArn,
			"Protocol":        sub.Protocol,
			"Endpoint":        sub.Endpoint,
			"Owner":           sub.Owner,
		})
	}
	return result
}

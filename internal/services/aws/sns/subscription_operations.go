package sns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/store/aws/common"
	snsstore "vorpalstacks/internal/store/aws/sns"

	"github.com/google/uuid"
)

// autoConfirmedProtocols lists protocols that vorpalstacks auto-confirms at
// Subscribe time, matching AWS behaviour. These subscriptions become
// immediately active without a confirmation round-trip.
var autoConfirmedProtocols = map[string]bool{
	"sqs":         true,
	"lambda":      true,
	"application": true,
}

// Subscribe creates a subscription to an SNS topic.
// https://docs.aws.amazon.com/sns/latest/api/API_Subscribe.html
func (s *SNSService) Subscribe(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	topicArn := request.GetParamLowerFirst(req.Parameters, "TopicArn")
	protocol := request.GetParamLowerFirst(req.Parameters, "Protocol")
	endpoint := request.GetParamLowerFirst(req.Parameters, "Endpoint")

	if topicArn == "" {
		return nil, awserrors.NewInvalidParameterException("TopicArn is required")
	}
	// Validate protocol value against the nine AWS-supported protocols.
	if err := validateProtocol(protocol); err != nil {
		return nil, err
	}
	if endpoint == "" {
		return nil, awserrors.NewInvalidParameterException("Endpoint is required")
	}
	// Validate endpoint format per protocol to catch grossly invalid
	// endpoints at Subscribe time rather than silently failing at delivery.
	if err := validateEndpointForProtocol(protocol, endpoint); err != nil {
		return nil, err
	}

	subscription := &snsstore.Subscription{
		TopicArn: topicArn,
		Protocol: protocol,
		Endpoint: endpoint,
		Owner:    reqCtx.GetAccountID(),
	}

	subscription.Attributes = parseAttributes(req.Parameters)

	for attrName, attrValue := range subscription.Attributes {
		if err := validateSubscriptionAttribute(attrName, attrValue); err != nil {
			return nil, err
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := store.CreateSubscription(subscription)
	if err != nil {
		if err == snsstore.ErrTopicNotFound {
			return nil, ErrTopicNotFound
		}
		return nil, err
	}

	needsConfirmation := !autoConfirmedProtocols[protocol]

	returnSubscriptionArn := request.GetParamLowerFirst(req.Parameters, "ReturnSubscriptionArn")
	if !needsConfirmation {
		if err := store.AutoConfirmSubscription(created); err != nil {
			return nil, err
		}
	} else if protocol == "http" || protocol == "https" {
		s.deliveryWg.Add(1)
		go func() {
			defer s.deliveryWg.Done()
			defer func() {
				if r := recover(); r != nil {
					logs.Error("SNS subscription confirmation panicked",
						logs.String("endpoint", created.Endpoint),
						logs.Any("panic", r))
				}
			}()
			s.sendSubscriptionConfirmation(created, reqCtx.GetRegion())
		}()
	}

	subArn := created.SubscriptionArn
	if needsConfirmation && strings.ToLower(returnSubscriptionArn) != "true" {
		subArn = "pending confirmation"
	}

	return map[string]interface{}{
		"SubscriptionArn": subArn,
	}, nil
}

// Unsubscribe deletes a subscription.
// https://docs.aws.amazon.com/sns/latest/api/API_Unsubscribe.html
func (s *SNSService) Unsubscribe(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	subscriptionArn := request.GetParamLowerFirst(req.Parameters, "SubscriptionArn")
	if subscriptionArn == "" {
		return nil, awserrors.NewInvalidParameterException("SubscriptionArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	subscription, err := store.GetSubscription(subscriptionArn)
	if err != nil {
		if err == snsstore.ErrSubscriptionNotFound {
			return nil, NewNotFoundException("Subscription does not exist")
		}
		return nil, err
	}

	// When the subscription was confirmed with AuthenticateOnUnsubscribe,
	// only the topic owner and the subscription owner may unsubscribe the
	// endpoint, and the request must be AWS-authenticated.
	if strings.EqualFold(subscription.Attributes["AuthenticateOnUnsubscribe"], "true") {
		if reqCtx.PrincipalType == request.PrincipalTypeAnonymous {
			return nil, ErrAuthorizationError
		}
		principalAccount := reqCtx.GetAccountID()
		if principalAccount != subscription.Owner {
			topic, err := store.GetTopic(subscription.TopicArn)
			if err != nil {
				if err == snsstore.ErrTopicNotFound {
					return nil, ErrTopicNotFound
				}
				return nil, err
			}
			if principalAccount != topic.Owner {
				return nil, ErrAuthorizationError
			}
		}
	}

	if err := store.DeleteSubscription(subscriptionArn); err != nil {
		if err == snsstore.ErrSubscriptionNotFound {
			return nil, NewNotFoundException("Subscription does not exist")
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ConfirmSubscription confirms a subscription request.
// https://docs.aws.amazon.com/sns/latest/api/API_ConfirmSubscription.html
func (s *SNSService) ConfirmSubscription(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	topicArn := request.GetParamLowerFirst(req.Parameters, "TopicArn")
	token := request.GetParamLowerFirst(req.Parameters, "Token")

	if topicArn == "" {
		return nil, awserrors.NewInvalidParameterException("TopicArn is required")
	}
	if token == "" {
		return nil, awserrors.NewInvalidParameterException("Token is required")
	}

	// AuthenticateOnUnsubscribe disallows unauthenticated unsubscribes of
	// this subscription. The parameter accepts only the boolean literals
	// "true" and "false"; nil means the parameter was not sent. A non-nil
	// value is persisted as the AuthenticateOnUnsubscribe subscription
	// attribute which Unsubscribe enforces.
	var authenticateOnUnsubscribe *bool
	if raw := request.GetParamLowerFirst(req.Parameters, "AuthenticateOnUnsubscribe"); raw != "" {
		switch strings.ToLower(raw) {
		case "true":
			val := true
			authenticateOnUnsubscribe = &val
		case "false":
			val := false
			authenticateOnUnsubscribe = &val
		default:
			return nil, awserrors.NewInvalidParameterException(
				fmt.Sprintf("Invalid AuthenticateOnUnsubscribe value %q: must be \"true\" or \"false\"", raw))
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	sub, err := store.FindSubscriptionByToken(topicArn, token)
	if err != nil {
		return nil, awserrors.NewInvalidParameterException("Subscription not found for token")
	}

	confirmed, err := store.ConfirmSubscription(sub.SubscriptionArn, token, authenticateOnUnsubscribe)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"SubscriptionArn": confirmed.SubscriptionArn,
	}, nil
}

// GetSubscriptionAttributes returns the attributes of a subscription.
// https://docs.aws.amazon.com/sns/latest/api/API_GetSubscriptionAttributes.html
func (s *SNSService) GetSubscriptionAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	subscriptionArn := request.GetParamLowerFirst(req.Parameters, "SubscriptionArn")
	if subscriptionArn == "" {
		return nil, awserrors.NewInvalidParameterException("SubscriptionArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	attrs, err := store.GetSubscriptionAttributes(subscriptionArn)
	if err != nil {
		if err == snsstore.ErrSubscriptionNotFound {
			return nil, NewNotFoundException("Subscription does not exist")
		}
		return nil, err
	}

	return map[string]interface{}{
		"Attributes": attrs,
	}, nil
}

// SetSubscriptionAttributes sets the attributes of a subscription.
// https://docs.aws.amazon.com/sns/latest/api/API_SetSubscriptionAttributes.html
func (s *SNSService) SetSubscriptionAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	subscriptionArn := request.GetParamLowerFirst(req.Parameters, "SubscriptionArn")
	attributeName := request.GetParamLowerFirst(req.Parameters, "AttributeName")
	attributeValue := request.GetParamLowerFirst(req.Parameters, "AttributeValue")

	if subscriptionArn == "" {
		return nil, awserrors.NewInvalidParameterException("SubscriptionArn is required")
	}
	if attributeName == "" {
		return nil, awserrors.NewInvalidParameterException("AttributeName is required")
	}

	if err := validateSubscriptionAttribute(attributeName, attributeValue); err != nil {
		return nil, err
	}

	attrs := map[string]string{attributeName: attributeValue}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.SetSubscriptionAttributes(subscriptionArn, attrs); err != nil {
		if err == snsstore.ErrSubscriptionNotFound {
			return nil, NewNotFoundException("Subscription does not exist")
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListSubscriptions lists the subscriptions.
// https://docs.aws.amazon.com/sns/latest/api/API_ListSubscriptions.html
func (s *SNSService) ListSubscriptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	nextToken := pagination.GetMarker(req.Parameters, "NextToken")
	result, err := store.ListSubscriptions(common.ListOptions{Marker: nextToken})
	if err != nil {
		return nil, err
	}

	subscriptions := buildSubscriptionList(result.Items)

	nextToken = ""
	if result.IsTruncated && result.NextMarker != "" {
		nextToken = result.NextMarker
	}
	return pagination.BuildListResponse("Subscriptions", subscriptions, nextToken), nil
}

// ListSubscriptionsByTopic lists the subscriptions by topic.
// https://docs.aws.amazon.com/sns/latest/api/API_ListSubscriptionsByTopic.html
func (s *SNSService) ListSubscriptionsByTopic(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	topicArn := request.GetParamLowerFirst(req.Parameters, "TopicArn")
	if topicArn == "" {
		return nil, awserrors.NewInvalidParameterException("TopicArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if _, err := store.GetTopic(topicArn); err != nil {
		return nil, awserrors.NewNotFoundException("Topic not found: " + topicArn)
	}

	nextToken := pagination.GetMarker(req.Parameters, "NextToken")
	result, err := store.ListSubscriptionsByTopic(topicArn, common.ListOptions{Marker: nextToken})
	if err != nil {
		return nil, err
	}

	subs := buildSubscriptionList(result.Items)

	nextToken = ""
	if result.IsTruncated && result.NextMarker != "" {
		nextToken = result.NextMarker
	}
	return pagination.BuildListResponse("Subscriptions", subs, nextToken), nil
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

package sns

import (
	"context"
	"crypto/subtle"
	"strconv"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/store/aws/common"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

// Subscribe creates a subscription to an SNS topic.
// https://docs.aws.amazon.com/sns/latest/api/API_Subscribe.html
func (s *SNSService) Subscribe(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	topicArn := request.GetParamLowerFirst(req.Parameters, "TopicArn")
	protocol := request.GetParamLowerFirst(req.Parameters, "Protocol")
	endpoint := request.GetParamLowerFirst(req.Parameters, "Endpoint")

	if topicArn == "" {
		return nil, awserrors.NewInvalidParameterException("TopicArn is required")
	}
	if protocol == "" {
		return nil, awserrors.NewInvalidParameterException("Protocol is required")
	}
	if endpoint == "" {
		return nil, awserrors.NewInvalidParameterException("Endpoint is required")
	}

	subscription := &snsstore.Subscription{
		TopicArn: topicArn,
		Protocol: protocol,
		Endpoint: endpoint,
		Owner:    reqCtx.GetAccountID(),
	}

	subscription.Attributes = parseSubscriptionAttributes(req.Parameters)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	created, err := store.CreateSubscription(subscription)
	if err != nil {
		if err == snsstore.ErrTopicNotFound {
			return nil, NewTopicNotFoundException()
		}
		return nil, err
	}

	needsConfirmation := protocol == "email" || protocol == "email-json" || protocol == "http" || protocol == "https" || protocol == "sms"

	returnSubscriptionArn := request.GetParamLowerFirst(req.Parameters, "ReturnSubscriptionArn")
	if !needsConfirmation {
		if err := store.AutoConfirmSubscription(created); err != nil {
			return nil, err
		}
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

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := store.ListSubscriptionsByTopic(topicArn, common.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, sub := range result.Items {
		if subtle.ConstantTimeCompare([]byte(sub.ConfirmationToken), []byte(token)) == 1 {
			confirmed, err := store.ConfirmSubscription(sub.SubscriptionArn, token)
			if err != nil {
				return nil, err
			}
			return map[string]interface{}{
				"SubscriptionArn": confirmed.SubscriptionArn,
			}, nil
		}
	}

	return nil, awserrors.NewInvalidParameterException("Subscription not found for token")
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

	subscriptions := make([]map[string]interface{}, 0, len(result.Items))
	for _, sub := range result.Items {
		subArn := sub.SubscriptionArn
		if sub.PendingConfirmation {
			subArn = "pending confirmation"
		}
		subscriptions = append(subscriptions, map[string]interface{}{
			"SubscriptionArn": subArn,
			"TopicArn":        sub.TopicArn,
			"Protocol":        sub.Protocol,
			"Endpoint":        sub.Endpoint,
			"Owner":           sub.Owner,
		})
	}

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

	nextToken := pagination.GetMarker(req.Parameters, "NextToken")
	result, err := store.ListSubscriptionsByTopic(topicArn, common.ListOptions{Marker: nextToken})
	if err != nil {
		return nil, err
	}

	subs := make([]map[string]interface{}, 0, len(result.Items))
	for _, sub := range result.Items {
		subArn := sub.SubscriptionArn
		if sub.PendingConfirmation {
			subArn = "pending confirmation"
		}
		subs = append(subs, map[string]interface{}{
			"SubscriptionArn": subArn,
			"TopicArn":        sub.TopicArn,
			"Protocol":        sub.Protocol,
			"Endpoint":        sub.Endpoint,
			"Owner":           sub.Owner,
		})
	}

	response := map[string]interface{}{
		"Subscriptions": subs,
	}
	if result.IsTruncated && result.NextMarker != "" {
		response["NextToken"] = result.NextMarker
	}
	return response, nil
}

func parseSubscriptionAttributes(params map[string]interface{}) map[string]string {
	result := make(map[string]string)

	if attrs, ok := params["Attributes"].(map[string]interface{}); ok {
		for k, v := range attrs {
			if vs, ok := v.(string); ok {
				result[k] = vs
			}
		}
	}
	if attrs, ok := params["attributes"].(map[string]interface{}); ok {
		for k, v := range attrs {
			if vs, ok := v.(string); ok {
				result[k] = vs
			}
		}
	}

	for i := 1; ; i++ {
		keyKey := "Attributes.entry." + strconv.Itoa(i) + ".key"
		valueKey := "Attributes.entry." + strconv.Itoa(i) + ".value"
		key := request.GetStringParam(params, keyKey)
		if key == "" {
			break
		}
		value := request.GetStringParam(params, valueKey)
		result[key] = value
	}

	return result
}

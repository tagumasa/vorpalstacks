package sns

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	snsstore "vorpalstacks/internal/store/aws/sns"
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

	return s.setSubscriptionAttributesCore(store, reqCtx, SetSubscriptionAttributesInput{
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

func buildSubscriptionList(items []*snsstore.Subscription) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(items))
	for _, sub := range items {
		subArn := sub.SubscriptionArn
		if sub.PendingConfirmation {
			subArn = pendingConfirmationARN
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

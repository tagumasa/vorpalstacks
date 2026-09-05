package sns

import (
	"fmt"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/store/aws/common"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

// autoConfirmedProtocols lists protocols that vorpalstacks auto-confirms at
// Subscribe time, matching AWS behaviour. These subscriptions become
// immediately active without a confirmation round-trip.
var autoConfirmedProtocols = map[string]bool{
	"sqs":         true,
	"lambda":      true,
	"application": true,
}

// SubscribeInput carries the parsed parameters for creating a subscription.
// Attributes holds the two-format attribute map parsed from the wire.
type SubscribeInput struct {
	TopicArn              string
	Protocol              string
	Endpoint              string
	ReturnSubscriptionArn string
	Attributes            map[string]string
}

// UnsubscribeInput carries the subscription ARN to delete.
type UnsubscribeInput struct {
	SubscriptionArn string
}

// ConfirmSubscriptionInput carries the parsed parameters for confirming a
// subscription. AuthenticateOnUnsubscribe holds the raw wire value; its
// boolean-literal validation runs inside the Core after the required-member
// checks, preserving the original error precedence.
type ConfirmSubscriptionInput struct {
	TopicArn                  string
	Token                     string
	AuthenticateOnUnsubscribe string
}

// GetSubscriptionAttributesInput carries the subscription ARN whose
// attributes are read.
type GetSubscriptionAttributesInput struct {
	SubscriptionArn string
}

// SetSubscriptionAttributesInput carries the attribute update for a
// subscription.
type SetSubscriptionAttributesInput struct {
	SubscriptionArn string
	AttributeName   string
	AttributeValue  string
}

// ListSubscriptionsInput carries the pagination token for listing
// subscriptions.
type ListSubscriptionsInput struct {
	NextToken string
}

// ListSubscriptionsByTopicInput carries the topic ARN and pagination token
// for listing a topic's subscriptions.
type ListSubscriptionsByTopicInput struct {
	TopicArn  string
	NextToken string
}

// subscribeCore is the single validation and persistence path for Subscribe.
// It needs the request context for the subscription owner account and the
// region used by the confirmation delivery.
func (s *SNSService) subscribeCore(store snsstore.SNSStoreInterface, reqCtx *request.RequestContext, in SubscribeInput) (interface{}, error) {
	if in.TopicArn == "" {
		return nil, NewInvalidParameter("TopicArn is required")
	}
	// Validate protocol value against the nine AWS-supported protocols.
	if err := validateProtocol(in.Protocol); err != nil {
		return nil, err
	}
	if in.Endpoint == "" {
		return nil, NewInvalidParameter("Endpoint is required")
	}
	// Validate endpoint format per protocol to catch grossly invalid
	// endpoints at Subscribe time rather than silently failing at delivery.
	if err := validateEndpointForProtocol(in.Protocol, in.Endpoint); err != nil {
		return nil, err
	}

	subscription := &snsstore.Subscription{
		TopicArn: in.TopicArn,
		Protocol: in.Protocol,
		Endpoint: in.Endpoint,
		Owner:    reqCtx.GetAccountID(),
	}

	subscription.Attributes = in.Attributes

	for attrName, attrValue := range subscription.Attributes {
		if err := validateSubscriptionAttribute(attrName, attrValue); err != nil {
			return nil, err
		}
	}

	created, err := store.CreateSubscription(subscription)
	if err != nil {
		if err == snsstore.ErrTopicNotFound {
			return nil, ErrTopicNotFound
		}
		return nil, err
	}

	needsConfirmation := !autoConfirmedProtocols[in.Protocol]

	if !needsConfirmation {
		if err := store.AutoConfirmSubscription(created); err != nil {
			return nil, err
		}
	} else if in.Protocol == "http" || in.Protocol == "https" {
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
	if needsConfirmation && strings.ToLower(in.ReturnSubscriptionArn) != "true" {
		subArn = "pending confirmation"
	}

	return map[string]interface{}{
		"SubscriptionArn": subArn,
	}, nil
}

// unsubscribeCore is the single validation and persistence path for
// Unsubscribe. It needs the request context for the owner-based
// authorisation checks on AuthenticateOnUnsubscribe subscriptions.
func (s *SNSService) unsubscribeCore(store snsstore.SNSStoreInterface, reqCtx *request.RequestContext, in UnsubscribeInput) (interface{}, error) {
	if in.SubscriptionArn == "" {
		return nil, NewInvalidParameter("SubscriptionArn is required")
	}

	subscription, err := store.GetSubscription(in.SubscriptionArn)
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

	if err := store.DeleteSubscription(in.SubscriptionArn); err != nil {
		if err == snsstore.ErrSubscriptionNotFound {
			return nil, NewNotFoundException("Subscription does not exist")
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// confirmSubscriptionCore is the single validation and persistence path for
// ConfirmSubscription.
func (s *SNSService) confirmSubscriptionCore(store snsstore.SNSStoreInterface, in ConfirmSubscriptionInput) (interface{}, error) {
	if in.TopicArn == "" {
		return nil, NewInvalidParameter("TopicArn is required")
	}
	if in.Token == "" {
		return nil, NewInvalidParameter("Token is required")
	}

	// AuthenticateOnUnsubscribe disallows unauthenticated unsubscribes of
	// this subscription. The parameter accepts only the boolean literals
	// "true" and "false"; nil means the parameter was not sent. A non-nil
	// value is persisted as the AuthenticateOnUnsubscribe subscription
	// attribute which Unsubscribe enforces.
	var authenticateOnUnsubscribe *bool
	if raw := in.AuthenticateOnUnsubscribe; raw != "" {
		switch strings.ToLower(raw) {
		case "true":
			val := true
			authenticateOnUnsubscribe = &val
		case "false":
			val := false
			authenticateOnUnsubscribe = &val
		default:
			return nil, NewInvalidParameter(
				fmt.Sprintf("Invalid AuthenticateOnUnsubscribe value %q: must be \"true\" or \"false\"", raw))
		}
	}

	sub, err := store.FindSubscriptionByToken(in.TopicArn, in.Token)
	if err != nil {
		return nil, NewInvalidParameter("Subscription not found for token")
	}

	confirmed, err := store.ConfirmSubscription(sub.SubscriptionArn, in.Token, authenticateOnUnsubscribe)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"SubscriptionArn": confirmed.SubscriptionArn,
	}, nil
}

// getSubscriptionAttributesCore is the single validation and persistence path
// for GetSubscriptionAttributes.
func (s *SNSService) getSubscriptionAttributesCore(store snsstore.SNSStoreInterface, in GetSubscriptionAttributesInput) (interface{}, error) {
	if in.SubscriptionArn == "" {
		return nil, NewInvalidParameter("SubscriptionArn is required")
	}

	attrs, err := store.GetSubscriptionAttributes(in.SubscriptionArn)
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

// setSubscriptionAttributesCore is the single validation and persistence path
// for SetSubscriptionAttributes.
func (s *SNSService) setSubscriptionAttributesCore(store snsstore.SNSStoreInterface, in SetSubscriptionAttributesInput) (interface{}, error) {
	if in.SubscriptionArn == "" {
		return nil, NewInvalidParameter("SubscriptionArn is required")
	}
	if in.AttributeName == "" {
		return nil, NewInvalidParameter("AttributeName is required")
	}

	if err := validateSubscriptionAttribute(in.AttributeName, in.AttributeValue); err != nil {
		return nil, err
	}

	attrs := map[string]string{in.AttributeName: in.AttributeValue}

	if err := store.SetSubscriptionAttributes(in.SubscriptionArn, attrs); err != nil {
		if err == snsstore.ErrSubscriptionNotFound {
			return nil, NewNotFoundException("Subscription does not exist")
		}
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// listSubscriptionsCore is the single validation and persistence path for
// ListSubscriptions.
func (s *SNSService) listSubscriptionsCore(store snsstore.SNSStoreInterface, in ListSubscriptionsInput) (interface{}, error) {
	nextToken := in.NextToken
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

// listSubscriptionsByTopicCore is the single validation and persistence path
// for ListSubscriptionsByTopic.
func (s *SNSService) listSubscriptionsByTopicCore(store snsstore.SNSStoreInterface, in ListSubscriptionsByTopicInput) (interface{}, error) {
	if in.TopicArn == "" {
		return nil, NewInvalidParameter("TopicArn is required")
	}

	if _, err := store.GetTopic(in.TopicArn); err != nil {
		return nil, awserrors.NewNotFoundException("Topic not found: " + in.TopicArn)
	}

	result, err := store.ListSubscriptionsByTopic(in.TopicArn, common.ListOptions{Marker: in.NextToken})
	if err != nil {
		return nil, err
	}

	subs := buildSubscriptionList(result.Items)

	nextToken := ""
	if result.IsTruncated && result.NextMarker != "" {
		nextToken = result.NextMarker
	}
	return pagination.BuildListResponse("Subscriptions", subs, nextToken), nil
}

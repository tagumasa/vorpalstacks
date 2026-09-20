package sns

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/store/aws/common"
	snsstore "vorpalstacks/internal/store/aws/sns"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

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
// region used by the confirmation delivery. Protocol acceptance,
// auto-confirm behaviour and endpoint validation are decided by the
// protocol registry (protocols.go) — the single protocol table shared
// with the delivery engine.
func (s *SNSService) subscribeCore(store snsstore.SNSStoreInterface, reqCtx *request.RequestContext, in SubscribeInput) (interface{}, error) {
	if in.TopicArn == "" {
		return nil, NewInvalidParameter("TopicArn is required")
	}
	// Validate protocol value against the registry's accepted protocols.
	if err := validateProtocol(in.Protocol); err != nil {
		return nil, err
	}
	if in.Endpoint == "" {
		return nil, NewInvalidParameter("Endpoint is required")
	}
	// Validate endpoint format per protocol to catch grossly invalid
	// endpoints at Subscribe time rather than silently failing at delivery.
	if err := validateProtocolEndpoint(in.Protocol, in.Endpoint); err != nil {
		return nil, err
	}

	topic, err := store.GetTopic(in.TopicArn)
	if err != nil {
		return nil, mapStoreError(err)
	}
	if err := rejectSubscriptionProtocol(in.Protocol, topic.IsFifoTopic()); err != nil {
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
		if err := validateSubscriptionAttribute(attrName, attrValue, in.Protocol); err != nil {
			return nil, err
		}
	}
	// The policy and its scope arrive as one attribute set; the pair check
	// validates the document under the scope this request establishes.
	if err := validateFilterPolicyScopePair(
		subscription.Attributes[snsstore.AttrFilterPolicy],
		subscription.Attributes[snsstore.AttrFilterPolicyScope]); err != nil {
		return nil, err
	}
	if raw, ok := subscription.Attributes[snsstore.AttrRedrivePolicy]; ok {
		if err := s.validateRedrivePolicyTarget(raw, reqCtx.GetRegion()); err != nil {
			return nil, err
		}
	}

	// The store is idempotent on the subscription natural key: a repeated
	// Subscribe returns the live subscription (created=false) instead of
	// accumulating a duplicate, and only a genuinely new record consumes
	// the topic's subscription and filter-policy quotas.
	created, isNew, err := store.CreateSubscription(subscription)
	if err != nil {
		return nil, mapStoreError(err)
	}

	needsConfirmation := !protocolAutoConfirms(in.Protocol)

	if isNew && !needsConfirmation {
		if err := store.AutoConfirmSubscription(created); err != nil {
			return nil, mapStoreError(err)
		}
	} else if needsConfirmation && protocolSendsConfirmationPost(in.Protocol) && (isNew || created.PendingConfirmation) {
		// The confirmation envelope carries the topic's SignatureVersion;
		// resolve it before the goroutine so the confirmation and later
		// notifications sign with the same version. A pending subscription
		// re-subscribed gets its confirmation request re-sent — the
		// endpoint never confirmed, so the retry gives it another chance —
		// while a CONFIRMED subscription does not: the repeated Subscribe
		// returns the live subscription, and a confirmation POST to an
		// endpoint that already confirmed is a request it has no reason to
		// act on again.
		dlv := s.resolveDeliveryContext(in.TopicArn, reqCtx.GetRegion())
		endpoint := created.Endpoint
		region := reqCtx.GetRegion()
		s.runTracked("subscription confirmation to "+endpoint, func() {
			s.sendSubscriptionConfirmation(created, region, dlv.signatureVersion)
		})
	}

	subArn := created.SubscriptionArn
	if needsConfirmation && strings.ToLower(in.ReturnSubscriptionArn) != "true" {
		subArn = pendingConfirmationARN
	}

	return map[string]interface{}{
		"SubscriptionArn": subArn,
	}, nil
}

// pendingConfirmationARN is the marker Subscribe and the list operations
// report in place of a confirmable ARN while a subscription awaits its
// confirmation round-trip.
const pendingConfirmationARN = "pending confirmation"

// rejectSubscriptionProtocol refuses protocols Subscribe cannot honestly
// accept: platform-refused protocols everywhere (the registry's refusal
// reason), and customer-managed endpoints on FIFO topics — "SNS FIFO
// topics can't deliver messages to customer managed endpoints ... Attempts
// to subscribe customer managed endpoints to SNS FIFO topics result in
// errors" (FIFO message delivery): email addresses, mobile apps, phone
// numbers for SMS, and HTTP(S) endpoints.
func rejectSubscriptionProtocol(protocol string, fifoTopic bool) error {
	info, ok := protocolRegistry[protocol]
	if !ok {
		return nil // validateProtocol already rejected unknown protocols.
	}
	if info.subscribeRefusal != "" {
		return NewInvalidParameter(fmt.Sprintf(
			"Invalid parameter: %s: %s", protocol, info.subscribeRefusal))
	}
	if fifoTopic && !info.fifoCapable {
		return NewInvalidParameter(fmt.Sprintf(
			"Invalid parameter: Protocol %s is a customer managed endpoint; SNS FIFO topics can't deliver messages to customer managed endpoints", protocol))
	}
	return nil
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
		return nil, mapStoreError(err)
	}

	// When the subscription was confirmed with AuthenticateOnUnsubscribe,
	// only the topic owner and the subscription owner may unsubscribe the
	// endpoint, and the request must be AWS-authenticated.
	if strings.EqualFold(subscription.Attributes[snsstore.AttrAuthenticateOnUnsubscribe], "true") {
		if reqCtx.PrincipalType == request.PrincipalTypeAnonymous {
			return nil, ErrAuthorizationError
		}
		principalAccount := reqCtx.GetAccountID()
		if principalAccount != subscription.Owner {
			topic, err := store.GetTopic(subscription.TopicArn)
			if err != nil {
				return nil, mapStoreError(err)
			}
			if principalAccount != topic.Owner {
				return nil, ErrAuthorizationError
			}
		}
	}

	if err := store.DeleteSubscription(in.SubscriptionArn); err != nil {
		return nil, mapStoreError(err)
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
		// The token miss answers the operation's InvalidParameter; a
		// storage failure keeps its own identity through the seam and
		// surfaces on the internal error plane instead.
		if errors.Is(err, snsstore.ErrSubscriptionNotFound) {
			return nil, NewInvalidParameter("Subscription not found for token")
		}
		return nil, mapStoreError(err)
	}

	confirmed, err := store.ConfirmSubscription(sub.SubscriptionArn, in.Token, authenticateOnUnsubscribe)
	if err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"SubscriptionArn": confirmed.SubscriptionArn,
	}, nil
}

// getSubscriptionAttributesCore is the single validation and persistence path
// for GetSubscriptionAttributes. The store materialises the stored map plus
// the documented synthesised keys; the EffectiveDeliveryPolicy member is
// resolved here against the subscription's own policy, the topic's default
// policy and the system defaults.
func (s *SNSService) getSubscriptionAttributesCore(store snsstore.SNSStoreInterface, in GetSubscriptionAttributesInput) (interface{}, error) {
	if in.SubscriptionArn == "" {
		return nil, NewInvalidParameter("SubscriptionArn is required")
	}

	subscription, err := store.GetSubscription(in.SubscriptionArn)
	if err != nil {
		return nil, mapStoreError(err)
	}

	attrs, err := store.GetSubscriptionAttributes(in.SubscriptionArn)
	if err != nil {
		return nil, mapStoreError(err)
	}

	subPolicy, err := parseSubscriptionDeliveryPolicy(subscription.Attributes[snsstore.AttrDeliveryPolicy])
	if err != nil {
		subPolicy = nil
	}
	var topicPolicy *deliveryPolicy
	if topic, err := store.GetTopic(subscription.TopicArn); err == nil {
		if topicPolicy, err = parseTopicDeliveryPolicy(topic.Attributes[snsstore.AttrDeliveryPolicy]); err != nil {
			topicPolicy = nil
		}
	}
	effective := resolveEffectiveSubscriptionPolicy(subPolicy, topicPolicy)
	if rendered, renderErr := effective.effectiveSubscriptionJSON(); renderErr == nil {
		attrs[snsstore.AttrEffectiveDelivery] = rendered
	}

	return map[string]interface{}{
		"Attributes": attrs,
	}, nil
}

// setSubscriptionAttributesCore is the single validation and persistence path
// for SetSubscriptionAttributes. The subscription is read first: the
// RawMessageDelivery protocol gate validates against the subscription's
// protocol, and a RedrivePolicy names a dead-letter queue that must exist.
func (s *SNSService) setSubscriptionAttributesCore(store snsstore.SNSStoreInterface, reqCtx *request.RequestContext, in SetSubscriptionAttributesInput) (interface{}, error) {
	if in.SubscriptionArn == "" {
		return nil, NewInvalidParameter("SubscriptionArn is required")
	}
	if in.AttributeName == "" {
		return nil, NewInvalidParameter("AttributeName is required")
	}

	subscription, err := store.GetSubscription(in.SubscriptionArn)
	if err != nil {
		return nil, mapStoreError(err)
	}

	if err := validateSubscriptionAttribute(in.AttributeName, in.AttributeValue, subscription.Protocol); err != nil {
		return nil, err
	}
	// FilterPolicy and FilterPolicyScope are coupled through the effective
	// pair this write leaves behind: a FilterPolicy write validates under
	// the scope that will apply (the request's value, else the stored one),
	// and a FilterPolicyScope write re-validates the stored policy — a
	// nested policy cannot be stranded under attribute-based matching.
	effectivePolicy, effectiveScope := subscription.Attributes[snsstore.AttrFilterPolicy], subscription.Attributes[snsstore.AttrFilterPolicyScope]
	if in.AttributeName == snsstore.AttrFilterPolicy {
		effectivePolicy = in.AttributeValue
	}
	if in.AttributeName == snsstore.AttrFilterPolicyScope {
		effectiveScope = in.AttributeValue
	}
	if err := validateFilterPolicyScopePair(effectivePolicy, effectiveScope); err != nil {
		return nil, err
	}
	if in.AttributeName == snsstore.AttrRedrivePolicy {
		if err := s.validateRedrivePolicyTarget(in.AttributeValue, reqCtx.GetRegion()); err != nil {
			return nil, err
		}
	}

	attrs := map[string]string{in.AttributeName: in.AttributeValue}

	if err := store.SetSubscriptionAttributes(in.SubscriptionArn, attrs); err != nil {
		return nil, mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// validateRedrivePolicyTarget enforces the documented dead-letter-queue
// constraints at attribute-set time: "The ARN must point to an Amazon SQS
// queue in the same AWS account and Region as your Amazon SNS
// subscription" (dead-letter queues documentation), and the queue must
// exist — the platform's SQS queue-attribute plane enforces the same
// contract for its own redrive policies, and a typo'd deadLetterTargetArn
// is otherwise discovered only when the failed message is already lost.
// Existence resolves through the SQS invoker seam because cross-service
// store imports are prohibited; when no invoker is wired the check fails
// closed rather than accepting an unverified target.
func (s *SNSService) validateRedrivePolicyTarget(value, region string) error {
	if strings.TrimSpace(value) == "" {
		// The empty value clears the association, the documented unset form.
		return nil
	}
	var rp snsstore.RedrivePolicy
	if err := json.Unmarshal([]byte(value), &rp); err != nil {
		return NewInvalidParameter(fmt.Sprintf("Invalid redrive policy: %s", err.Error()))
	}
	arn := rp.DeadLetterTargetArn
	if strings.TrimSpace(arn) == "" {
		return NewInvalidParameter("Invalid redrive policy: deadLetterTargetArn is required")
	}

	_, service, arnRegion, arnAccount, resource := svcarn.SplitARN(arn)
	if service != "sqs" || resource == "" {
		return NewInvalidParameter(fmt.Sprintf(
			"Invalid RedrivePolicy: deadLetterTargetArn must be an Amazon SQS queue ARN: %s", arn))
	}
	if arnRegion != region || arnAccount != s.accountID {
		return NewInvalidParameter(fmt.Sprintf(
			"Invalid RedrivePolicy: the dead-letter queue must be in the same account and Region as the subscription: %s", arn))
	}

	if s.bus == nil {
		return NewInvalidParameter("Invalid RedrivePolicy: the SQS queue cannot be verified because no service bus is wired")
	}
	sqsInvoker := s.bus.SQSInvoker()
	if sqsInvoker == nil {
		return NewInvalidParameter("Invalid RedrivePolicy: the SQS queue cannot be verified because the SQS invoker is unavailable")
	}
	queueName := svcarn.ExtractQueueNameFromARN(arn)
	if _, err := sqsInvoker.GetQueueByName(context.Background(), arnRegion, queueName); err != nil {
		return NewInvalidParameter(fmt.Sprintf(
			"Invalid RedrivePolicy: deadLetterTargetArn does not resolve to an existing Amazon SQS queue: %s", arn))
	}
	return nil
}

// listSubscriptionsCore is the single validation and persistence path for
// ListSubscriptions.
func (s *SNSService) listSubscriptionsCore(store snsstore.SNSStoreInterface, in ListSubscriptionsInput) (interface{}, error) {
	nextToken := in.NextToken
	result, err := store.ListSubscriptions(common.ListOptions{Marker: nextToken})
	if err != nil {
		return nil, mapStoreError(err)
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
		return nil, mapStoreError(err)
	}

	result, err := store.ListSubscriptionsByTopic(in.TopicArn, common.ListOptions{Marker: in.NextToken})
	if err != nil {
		return nil, mapStoreError(err)
	}

	subs := buildSubscriptionList(result.Items)

	nextToken := ""
	if result.IsTruncated && result.NextMarker != "" {
		nextToken = result.NextMarker
	}
	return pagination.BuildListResponse("Subscriptions", subs, nextToken), nil
}

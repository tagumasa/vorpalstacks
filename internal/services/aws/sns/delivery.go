package sns

// delivery.go is the fan-out engine: the seam a validated publish passes
// through on its way to the topic's subscriptions, the per-protocol
// delivery paths (SQS, HTTP/S, Lambda, dead-letter redrive) and the retry,
// throttle and shutdown machinery those paths share. The wire wrappers and
// validation live in publish_operations.go / publish_core.go; the envelope
// construction and signing the non-raw deliveries consume live in
// envelope.go.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	snsstore "vorpalstacks/internal/store/aws/sns"
	arnutil "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/aws/queueurl"
)

// dispatchPublish is the single dispatch seam between a validated message
// and its fan-out: through the event bus when wired (the subscribed
// handleBusDelivery re-lists the topic's subscriptions and delivers), or
// through a tracked delivery goroutine when no bus is configured. FIFO
// topics never take this seam — their per-group ordering requires the
// synchronous publishFifoOrdered path, which calls deliverToSubscriptions
// inside the store's group lock. The message and subscription records must
// already be private copies — the bus path serialises them for transport
// and the direct path hands them to a goroutine that outlives this call.
// The error policy is surface: this function never swallows a dispatch
// failure; its callers decide how the error reaches their own caller.
func (s *SNSService) dispatchPublish(msg *snsstore.Message, subscriptions []*snsstore.Subscription, region string) error {
	if s.bus != nil {
		// Serialise message attributes to raw JSON for transport through
		// the event bus (which must not depend on store-layer types).
		var msgAttrs map[string]json.RawMessage
		if len(msg.MessageAttributes) > 0 {
			msgAttrs = make(map[string]json.RawMessage, len(msg.MessageAttributes))
			for k, v := range msg.MessageAttributes {
				if raw, err := json.Marshal(v); err == nil {
					msgAttrs[k] = raw
				} else {
					// The attribute rides no further, so the delivered
					// message degrades — say so instead of dropping it
					// silently.
					logs.Warn("SNS: dropping message attribute the bus envelope cannot marshal",
						logs.String("attribute", k),
						logs.Err(err))
				}
			}
		}
		evt := &eventbus.SNSDeliveryEvent{
			TopicARN:               msg.TopicArn,
			MessageID:              msg.MessageId,
			Message:                msg.Message,
			Subject:                msg.Subject,
			MessageStructure:       msg.MessageStructure,
			MessageGroupId:         msg.MessageGroupId,
			MessageDeduplicationID: msg.MessageDeduplicationId,
			MessageAttributes:      msgAttrs,
		}
		evt.Region = region
		return s.bus.Publish(context.Background(), evt)
	}

	msgCopy := *msg
	s.deliverAsync(&msgCopy, subscriptions, region)
	return nil
}

// deliverToSubscriptions delivers one message to every confirmed,
// filter-matching subscription. Standard topics reach it through the async
// dispatch paths; FIFO topics call it synchronously inside the store's
// per-(topic, message group) lock, where one goroutine's sequential loop
// makes the delivery order per endpoint the sequence-allocation order.
func (s *SNSService) deliverToSubscriptions(msg *snsstore.Message, subscriptions []*snsstore.Subscription, region string) {
	dlv := s.resolveDeliveryContext(msg.TopicArn, region)

	for _, sub := range subscriptions {
		if sub.PendingConfirmation {
			continue
		}

		if !matchFilterPolicy(sub.GetFilterPolicy(), sub.GetFilterPolicyScope(), msg.MessageAttributes, msg.Message) {
			continue
		}

		var deliveryErr error
		if !protocolHasDeliveryHandler(sub.Protocol) {
			// Return a delivery error so the message routes to the DLQ when
			// a RedrivePolicy is configured, rather than being silently lost
			// by an earlier `logs.Warn + continue` path.
			deliveryErr = fmt.Errorf("unsupported protocol %q: no delivery handler available", sub.Protocol)
		} else {
			switch sub.Protocol {
			case "sqs":
				deliveryErr = s.deliverToSQS(msg, sub, region, dlv)
			case "http", "https":
				deliveryErr = s.deliverToHTTP(msg, sub, region, dlv)
			case "lambda":
				deliveryErr = s.deliverToLambda(msg, sub, region, dlv)
			default:
				// The registry and the engine drifted apart: the table
				// declares a handler this switch does not carry. Fail the
				// delivery loudly instead of succeeding silently.
				deliveryErr = fmt.Errorf("protocol %q: the registry declares a delivery handler the engine does not implement", sub.Protocol)
			}
		}

		if deliveryErr != nil {
			s.handleDeliveryFailure(msg, sub, region, deliveryErr, dlv)
		}
	}
}

// deliveryContext carries the topic-level delivery facts one fan-out pass
// resolves once: the signature version the envelopes carry (and the signing
// hash it selects) and the topic's default HTTP/S delivery policy every
// subscription's own policy resolves against.
type deliveryContext struct {
	signatureVersion string
	topicPolicy      *deliveryPolicy
}

// resolveDeliveryContext reads the topic's delivery-relevant attributes.
// When the topic cannot be read (deleted mid-flight, storage hiccup) the
// documented defaults apply — an envelope signed with the default version
// and the system default retry policy still delivers, which beats dropping
// the message on a degraded read.
func (s *SNSService) resolveDeliveryContext(topicArn, region string) deliveryContext {
	dlv := deliveryContext{
		// "By default, SignatureVersion is set to 1" (GetTopicAttributes);
		// absence means 1.
		signatureVersion: "1",
	}
	store, err := s.getSNSStoreByRegion(region)
	if err != nil {
		return dlv
	}
	topic, err := store.GetTopic(topicArn)
	if err != nil {
		logs.Warn("SNS delivery: topic unavailable while resolving the delivery context, using defaults",
			logs.String("topicArn", topicArn),
			logs.Err(err))
		return dlv
	}
	if v := topic.Attributes[snsstore.AttrSignatureVersion]; v == "1" || v == "2" {
		dlv.signatureVersion = v
	}
	if topicPolicy, err := parseTopicDeliveryPolicy(topic.Attributes[snsstore.AttrDeliveryPolicy]); err == nil {
		dlv.topicPolicy = topicPolicy
	}
	return dlv
}

// handleDeliveryFailure routes a failed delivery to the subscription's
// dead-letter queue when a RedrivePolicy is configured.
func (s *SNSService) handleDeliveryFailure(msg *snsstore.Message, sub *snsstore.Subscription, region string, deliveryErr error, dlv deliveryContext) {
	rp, err := sub.GetRedrivePolicy()
	if err != nil {
		logs.Warn("Failed to parse subscription RedrivePolicy",
			logs.String("subscriptionArn", sub.SubscriptionArn),
			logs.Err(err))
		return
	}
	if rp == nil || rp.DeadLetterTargetArn == "" {
		// No DLQ configured — the message is dropped. Log so that silent
		// loss (e.g. unsupported protocol without RedrivePolicy) leaves
		// a trace for operators.
		logs.Warn("SNS delivery failed with no DLQ configured — message dropped",
			logs.String("subscriptionArn", sub.SubscriptionArn),
			logs.String("topicArn", msg.TopicArn),
			logs.String("protocol", sub.Protocol),
			logs.String("messageId", msg.MessageId),
			logs.Err(deliveryErr))
		return
	}

	logs.Warn("SNS delivery failed, routing to DLQ",
		logs.String("subscriptionArn", sub.SubscriptionArn),
		logs.String("dlqArn", rp.DeadLetterTargetArn),
		logs.String("topicArn", msg.TopicArn),
		logs.Err(deliveryErr))

	if dlqErr := s.deliverToDLQ(msg, sub, rp.DeadLetterTargetArn, region, dlv); dlqErr != nil {
		logs.Error("DLQ delivery also failed — message permanently lost",
			logs.String("dlqArn", rp.DeadLetterTargetArn),
			logs.String("messageId", msg.MessageId),
			logs.Err(dlqErr))
	}
}

// deliverToSQS delivers the notification to the subscription's queue,
// honouring the subscription's raw-message setting and the queue's
// resource-policy authorisation.
func (s *SNSService) deliverToSQS(msg *snsstore.Message, sub *snsstore.Subscription, region string, dlv deliveryContext) error {
	return s.deliverToSQSEndpoint(msg, sub, sub.Endpoint, sub.IsRawMessageDelivery(), true, region, dlv)
}

// deliverToDLQ sends a failed message to the dead-letter SQS queue — SQS
// delivery with the redrive endpoint: the full signed envelope (a redrive
// receives the notification, not the bare payload) and no
// subscription-policy authorisation (the redrive target's policy is not
// the subscription's delivery contract).
func (s *SNSService) deliverToDLQ(msg *snsstore.Message, sub *snsstore.Subscription, dlqArn, region string, dlv deliveryContext) error {
	return s.deliverToSQSEndpoint(msg, sub, dlqArn, false, false, region, dlv)
}

// resolveQueueURL resolves an SQS endpoint to its queue URL and region: a
// queue ARN names its own region and resolves to the URL through the SQS
// invoker; a queue URL is used as-is and names its region when it embeds
// one (AWS-form URLs), falling back to the topic's region when it embeds
// none — the subscription lives in the topic's regional store, so a
// region-less platform URL addresses the topic's own region, never the
// global default.
func resolveQueueURL(sqsInvoker invokers.SQSInvoker, endpoint, fallbackRegion string) (queueURL, region string) {
	queueURL = endpoint
	region = fallbackRegion
	if strings.HasPrefix(endpoint, "arn:") {
		queueName := arnutil.ExtractQueueNameFromARN(endpoint)
		_, _, region, _, _ = arnutil.SplitARN(endpoint)
		if queueName != "" {
			resolvedURL, err := sqsInvoker.GetQueueByName(context.Background(), region, queueName)
			if err == nil {
				queueURL = resolvedURL
			} else {
				// The send below will fail against the raw ARN and route to
				// the DLQ path — the resolution failure is the diagnosis,
				// so it must not vanish.
				logs.Warn("SNS: resolving the queue ARN to its URL failed; sending will fail against the ARN",
					logs.String("queueName", queueName),
					logs.String("region", region),
					logs.Err(err))
			}
		}
	} else if fromURL := queueurl.RegionFromQueueURL(endpoint); fromURL != "" {
		region = fromURL
	}
	return queueURL, region
}

// sqsSendOptions builds the typed SQS send options from the message: the
// typed attribute conversion and the FIFO identifiers are the same for
// subscription delivery and DLQ redrive. MessageGroupId is forwarded to
// every SQS destination (standard queues accept it — fair queues), while
// MessageDeduplicationId is forwarded to FIFO queues only: a standard
// queue rejects the parameter, and a FIFO topic's standard-queue
// subscriptions (the documented best-effort channel) must still deliver.
func sqsSendOptions(msg *snsstore.Message, fifoDestination bool) invokers.SQSSendOptions {
	typedAttrs := make(map[string]invokers.SQSMessageAttribute, len(msg.MessageAttributes))
	for k, v := range msg.MessageAttributes {
		ta := invokers.SQSMessageAttribute{DataType: v.Type}
		if len(v.BinaryValue) > 0 {
			ta.BinaryValue = v.BinaryValue
		} else {
			ta.StringValue = v.StringValue
		}
		typedAttrs[k] = ta
	}
	opts := invokers.SQSSendOptions{
		TypedMessageAttributes: typedAttrs,
		MessageGroupID:         msg.MessageGroupId,
	}
	if fifoDestination {
		opts.MessageDeduplicationID = msg.MessageDeduplicationId
	}
	return opts
}

// queueNameFromDestination extracts the queue name from an SQS endpoint:
// queue ARNs carry it as the resource; queue URLs carry it as the final
// path segment.
func queueNameFromDestination(endpoint, queueURL string) string {
	if strings.HasPrefix(endpoint, "arn:") {
		return arnutil.ExtractQueueNameFromARN(endpoint)
	}
	if idx := strings.LastIndex(queueURL, "/"); idx >= 0 && idx+1 < len(queueURL) {
		return queueURL[idx+1:]
	}
	return ""
}

// isFIFOQueueDestination reports whether an SQS destination is a FIFO
// queue: FIFO queue names carry the mandatory ".fifo" suffix.
func isFIFOQueueDestination(endpoint, queueURL string) bool {
	return strings.HasSuffix(queueNameFromDestination(endpoint, queueURL), ".fifo")
}

// deliverToSQSEndpoint delivers the message to one SQS destination — a
// subscription's queue or a subscription's dead-letter queue. raw selects
// the bare protocol message over the signed notification envelope;
// authorise runs the queue resource-policy evaluation the subscription
// delivery contract requires. The authorisation is fail-closed: a queue
// whose ARN cannot be resolved, or whose evaluation errors, fails the
// delivery (routing to the DLQ) rather than delivering unverified — the
// same failure mode the Lambda path has always had.
func (s *SNSService) deliverToSQSEndpoint(msg *snsstore.Message, sub *snsstore.Subscription, endpoint string, raw, authorise bool, region string, dlv deliveryContext) error {
	if s.bus == nil {
		return fmt.Errorf("event bus not available for SQS delivery")
	}
	sqsInvoker := s.bus.SQSInvoker()
	if sqsInvoker == nil {
		return fmt.Errorf("SQS invoker not available for SQS delivery")
	}

	queueURL, sqsRegion := resolveQueueURL(sqsInvoker, endpoint, region)

	if authorise {
		queueARN, qErr := sqsInvoker.GetQueueARN(context.Background(), sqsRegion, queueURL)
		if qErr != nil {
			return fmt.Errorf("resource policy evaluation skipped: queue ARN resolution for %s failed (fail-closed): %w", queueURL, qErr)
		}
		if queueARN == "" {
			return fmt.Errorf("resource policy evaluation skipped: queue ARN resolution for %s returned no ARN (fail-closed)", queueURL)
		}
		allowed, evalErr := s.bus.EvaluateTargetPolicy(context.Background(), queueARN, "sqs", "sns.amazonaws.com", "sqs:SendMessage", queueARN)
		if evalErr != nil {
			logs.Warn("resource policy evaluation failed for SQS delivery, dropping message",
				logs.String("queueArn", queueARN),
				logs.String("topicArn", msg.TopicArn),
				logs.Err(evalErr))
			return fmt.Errorf("resource policy evaluation failed: %w", evalErr)
		}
		if !allowed {
			return fmt.Errorf("resource policy denied delivery to queue %s", queueARN)
		}
	}

	protocolMessage, err := extractProtocolMessage(msg, "sqs")
	if err != nil {
		return fmt.Errorf("extract protocol message for SQS: %w", err)
	}

	var body string
	if raw {
		body = protocolMessage
	} else {
		payload := s.buildNotificationEnvelope(msg, sub, region, protocolMessage, dlv.signatureVersion)
		s.signNotificationEnvelope(payload, region, dlv.signatureVersion)

		jsonData, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("marshal SQS notification: %w", err)
		}
		body = string(jsonData)
	}

	if _, _, err := sqsInvoker.SendMessage(context.Background(), sqsRegion, queueURL, body, sqsSendOptions(msg, isFIFOQueueDestination(endpoint, queueURL))); err != nil {
		return fmt.Errorf("send to queue %s: %w", queueURL, err)
	}

	return nil
}

// deliverToHTTP delivers the notification to an HTTP/S endpoint under the
// subscription's effective delivery policy: "Amazon SNS considers all 5XX
// errors and 429 (too many requests sent) errors as retryable. These
// errors are subject to the delivery policy. All other errors are
// considered as permanent failures and retries will not be attempted"
// (message delivery retries); a transport failure — the system hosting
// the endpoint never answering — is the defining server-side error the policy
// exists for and is retryable too. The retry schedule is the policy's
// four-phase ladder; when it is exhausted (or a permanent failure occurs)
// the error routes to handleDeliveryFailure like every delivery failure.
func (s *SNSService) deliverToHTTP(msg *snsstore.Message, sub *snsstore.Subscription, region string, dlv deliveryContext) error {
	subPolicy, _ := parseSubscriptionDeliveryPolicy(sub.Attributes[snsstore.AttrDeliveryPolicy])
	policy := resolveEffectiveSubscriptionPolicy(subPolicy, dlv.topicPolicy)

	protocolMessage, err := extractProtocolMessage(msg, sub.Protocol)
	if err != nil {
		return fmt.Errorf("extract protocol message for HTTP: %w", err)
	}

	var jsonData []byte
	if sub.IsRawMessageDelivery() {
		jsonData = []byte(protocolMessage)
	} else {
		payload := s.buildNotificationEnvelope(msg, sub, region, protocolMessage, dlv.signatureVersion)
		s.signNotificationEnvelope(payload, region, dlv.signatureVersion)

		jsonData, err = json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("marshal HTTP notification: %w", err)
		}
	}

	contentType := policy.headerContentType
	if contentType == "" {
		// "By default, Amazon SNS sends all the notifications to HTTP/S
		// endpoints with content type set to text/plain; charset=UTF-8."
		contentType = "text/plain; charset=UTF-8"
	}

	schedule := policy.retrySchedule()
	var lastErr error
	for attempt := 0; ; attempt++ {
		if policy.maxReceivesPerSecond > 0 {
			s.throttleDelivery(sub.SubscriptionArn, policy.maxReceivesPerSecond)
		}
		retryable, err := s.attemptHTTPDelivery(sub, jsonData, contentType, msg)
		if err == nil {
			return nil
		}
		lastErr = err
		if !retryable {
			return err
		}
		if attempt >= len(schedule) {
			return fmt.Errorf("HTTP delivery to %s exhausted the delivery policy's %d retries: %w", sub.Endpoint, len(schedule), lastErr)
		}
		if !s.sleepDelivery(schedule[attempt]) {
			return fmt.Errorf("HTTP delivery to %s abandoned at shutdown after attempt %d: %w", sub.Endpoint, attempt+1, lastErr)
		}
	}
}

// attemptHTTPDelivery performs one HTTP/S delivery attempt. The returned
// retryable flag applies the documented error classification: 5XX and 429
// responses and transport errors retry under the delivery policy, every
// other status is a permanent failure.
func (s *SNSService) attemptHTTPDelivery(sub *snsstore.Subscription, body []byte, contentType string, msg *snsstore.Message) (retryable bool, err error) {
	req, err := http.NewRequest("POST", sub.Endpoint, bytes.NewBuffer(body))
	if err != nil {
		return false, fmt.Errorf("create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-amz-sns-message-type", "Notification")
	req.Header.Set("x-amz-sns-message-id", msg.MessageId)
	req.Header.Set("x-amz-sns-topic-arn", msg.TopicArn)
	req.Header.Set("x-amz-sns-subscription-arn", sub.SubscriptionArn)

	resp, err := s.httpDeliveryClient().Do(req)
	if err != nil {
		// The endpoint never answered — the hosted system is unavailable,
		// the server-side error class the retry policy governs.
		return true, fmt.Errorf("HTTP delivery to %s: %w", sub.Endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 || resp.StatusCode >= 500 {
		return true, fmt.Errorf("HTTP delivery to %s returned retryable status %d", sub.Endpoint, resp.StatusCode)
	}
	if resp.StatusCode >= 400 {
		return false, fmt.Errorf("HTTP delivery to %s returned status %d", sub.Endpoint, resp.StatusCode)
	}

	logs.Debug("HTTP notification delivered",
		logs.String("endpoint", sub.Endpoint),
		logs.Int("status", resp.StatusCode))

	return false, nil
}

// throttleDelivery paces one delivery attempt to the subscription's
// maxReceivesPerSecond — "The maximum average number of message deliveries
// per second, per subscription". Waiting happens outside the limiter's
// critical section so concurrent deliveries to other subscriptions never
// queue behind this one's sleep.
func (s *SNSService) throttleDelivery(subscriptionArn string, ratePerSecond int) {
	interval := time.Second / time.Duration(ratePerSecond)
	s.throttleMu.Lock()
	if s.throttleNext == nil {
		s.throttleNext = make(map[string]time.Time)
	}
	now := time.Now()
	next := s.throttleNext[subscriptionArn]
	if next.Before(now) {
		next = now
	}
	s.throttleNext[subscriptionArn] = next.Add(interval)
	wait := next.Sub(now)
	s.throttleMu.Unlock()

	if wait > 0 {
		s.sleepDelivery(wait)
	}
}

// sleepDelivery waits for d, reporting false when the service's delivery
// lifecycle ended first — retry ladders then abandon their remaining
// attempts instead of holding delivery goroutines through shutdown.
func (s *SNSService) sleepDelivery(d time.Duration) bool {
	if d <= 0 {
		return s.deliveryCtx.Err() == nil
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-timer.C:
		return true
	case <-s.deliveryCtx.Done():
		return false
	}
}

// deliverToLambda invokes the subscription's Lambda function with the
// notification nested in the aws:sns event envelope. A missing bus or
// invoker is a delivery error like the SQS path's — the failure routes to
// handleDeliveryFailure (and the DLQ when configured) instead of
// returning success with no delivery.
func (s *SNSService) deliverToLambda(msg *snsstore.Message, sub *snsstore.Subscription, region string, dlv deliveryContext) error {
	if s.bus == nil {
		return fmt.Errorf("event bus not available for Lambda delivery")
	}

	lambdaInvoker := s.bus.LambdaInvoker()
	if lambdaInvoker == nil {
		return fmt.Errorf("Lambda invoker not available for Lambda delivery")
	}

	if region == "" {
		region = s.defaultRegion
	}

	functionARN := sub.Endpoint
	if !strings.HasPrefix(functionARN, "arn:") {
		functionARN = arnutil.NewARNBuilder(s.accountID, region).Lambda().Function(sub.Endpoint)
	}
	allowed, evalErr := s.bus.EvaluateTargetPolicy(context.Background(), functionARN, "lambda", "sns.amazonaws.com", "lambda:InvokeFunction", functionARN)
	if evalErr != nil {
		return fmt.Errorf("resource policy evaluation failed: %w", evalErr)
	}
	if !allowed {
		return fmt.Errorf("resource policy denied invocation of %s", functionARN)
	}

	protocolMessage, err := extractProtocolMessage(msg, "lambda")
	if err != nil {
		return fmt.Errorf("extract protocol message for Lambda: %w", err)
	}

	// Lambda receives the event envelope with the notification nested under
	// Sns — raw message delivery is documented for "Amazon SQS or HTTP/S
	// endpoints" only, so the Lambda plane has no raw form.
	snsPayload := s.buildNotificationEnvelope(msg, sub, region, protocolMessage, dlv.signatureVersion)
	s.signNotificationEnvelope(snsPayload, region, dlv.signatureVersion)

	record := map[string]interface{}{
		"EventSource":          "aws:sns",
		"EventVersion":         "1.0",
		"EventSubscriptionArn": sub.SubscriptionArn,
		"Sns":                  snsPayload,
	}

	eventEnvelope := map[string]interface{}{
		"Records": []interface{}{record},
	}

	jsonData, err := json.Marshal(eventEnvelope)
	if err != nil {
		return fmt.Errorf("marshal Lambda event envelope: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	functionName := sub.Endpoint
	if _, _, err := lambdaInvoker.InvokeForGateway(ctx, functionName, jsonData); err != nil {
		return fmt.Errorf("invoke Lambda %s (subscription %s, messageId %s): %w", functionName, sub.SubscriptionArn, msg.MessageId, err)
	}

	return nil
}

// sendSubscriptionConfirmation sends a SubscriptionConfirmation message
// to an HTTP/HTTPS endpoint. This is best-effort: if the endpoint is
// unreachable the subscription simply stays in pending state until the
// subscriber retries. signatureVersion is the topic's SignatureVersion
// attribute (default "1") — the confirmation envelope carries it as a
// member and the signing step selects the hash it names.
func (s *SNSService) sendSubscriptionConfirmation(sub *snsstore.Subscription, region string, signatureVersion string) {
	envelope := s.buildConfirmationEnvelope(sub, region, signatureVersion)
	s.signNotificationEnvelope(envelope, region, signatureVersion)

	body, err := json.Marshal(envelope)
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
	req.Header.Set("x-amz-sns-message-id", envelope["MessageId"].(string))
	req.Header.Set("x-amz-sns-topic-arn", sub.TopicArn)
	req.Header.Set("x-amz-sns-subscription-arn", sub.SubscriptionArn)

	resp, err := s.httpDeliveryClient().Do(req)
	if err != nil {
		logs.Warn("SNS: failed to send subscription confirmation",
			logs.String("endpoint", sub.Endpoint),
			logs.String("subscriptionArn", sub.SubscriptionArn),
			logs.Err(err))
		return
	}
	defer resp.Body.Close()

	// The POST is best-effort (a failed confirmation leaves the
	// subscription pending), but an error status is a warning, not a
	// success trace — the subscriber sees why the subscription never
	// activated.
	if resp.StatusCode >= 400 {
		logs.Warn("SNS: subscription confirmation endpoint answered with an error status",
			logs.String("endpoint", sub.Endpoint),
			logs.String("subscriptionArn", sub.SubscriptionArn),
			logs.Int("status", resp.StatusCode))
		return
	}

	logs.Debug("SNS: subscription confirmation sent",
		logs.String("endpoint", sub.Endpoint),
		logs.Int("status", resp.StatusCode))
}

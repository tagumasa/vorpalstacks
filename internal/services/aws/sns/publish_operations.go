package sns

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"vorpalstacks/internal/common/defaults"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	snsstore "vorpalstacks/internal/store/aws/sns"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// Publish publishes a message to an SNS topic.
// https://docs.aws.amazon.com/sns/latest/api/API_Publish.html
//
// All validation, deduplication, and delivery fan-out live inside the Core.
func (s *SNSService) Publish(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.publishCore(store, reqCtx, PublishInput{
		TopicArn:               request.GetParamLowerFirst(req.Parameters, "TopicArn"),
		TargetArn:              request.GetParamLowerFirst(req.Parameters, "TargetArn"),
		PhoneNumber:            request.GetParamLowerFirst(req.Parameters, "PhoneNumber"),
		Message:                request.GetParamLowerFirst(req.Parameters, "Message"),
		Subject:                request.GetParamLowerFirst(req.Parameters, "Subject"),
		MessageStructure:       request.GetParamLowerFirst(req.Parameters, "MessageStructure"),
		MessageGroupId:         request.GetParamLowerFirst(req.Parameters, "MessageGroupId"),
		MessageDeduplicationId: request.GetParamLowerFirst(req.Parameters, "MessageDeduplicationId"),
		Parameters:             req.Parameters,
	})
}

func (s *SNSService) deliverToSubscriptions(msg *snsstore.Message, subscriptions []*snsstore.Subscription, region string) {
	for _, sub := range subscriptions {
		if sub.PendingConfirmation {
			continue
		}

		if !matchFilterPolicy(sub.GetFilterPolicy(), sub.GetFilterPolicyScope(), msg.MessageAttributes, msg.Message) {
			continue
		}

		var deliveryErr error
		switch sub.Protocol {
		case "sqs":
			deliveryErr = s.deliverToSQS(msg, sub)
		case "http", "https":
			deliveryErr = s.deliverToHTTP(msg, sub, region)
		case "lambda":
			deliveryErr = s.deliverToLambda(msg, sub, region)
		default:
			// Return a delivery error so the message routes to the DLQ when
			// a RedrivePolicy is configured, rather than being silently lost
			// by an earlier `logs.Warn + continue` path.
			deliveryErr = fmt.Errorf("unsupported protocol %q: no delivery handler available", sub.Protocol)
		}

		if deliveryErr != nil {
			s.handleDeliveryFailure(msg, sub, region, deliveryErr)
		}
	}
}

// handleDeliveryFailure routes a failed delivery to the subscription's
// dead-letter queue when a RedrivePolicy is configured.
func (s *SNSService) handleDeliveryFailure(msg *snsstore.Message, sub *snsstore.Subscription, region string, deliveryErr error) {
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

	if dlqErr := s.deliverToDLQ(msg, sub, rp.DeadLetterTargetArn); dlqErr != nil {
		logs.Error("DLQ delivery also failed — message permanently lost",
			logs.String("dlqArn", rp.DeadLetterTargetArn),
			logs.String("messageId", msg.MessageId),
			logs.Err(dlqErr))
	}
}

// deliverToDLQ sends a failed message to the dead-letter SQS queue.
func (s *SNSService) deliverToDLQ(msg *snsstore.Message, sub *snsstore.Subscription, dlqArn string) error {
	if s.bus == nil {
		return fmt.Errorf("event bus not available for DLQ delivery")
	}

	sqsInvoker := s.bus.SQSInvoker()
	if sqsInvoker == nil {
		return fmt.Errorf("SQS invoker not available for DLQ delivery")
	}

	queueURL := dlqArn
	queueName := arnutil.ExtractQueueNameFromARN(dlqArn)
	_, _, dlqRegion, _, _ := arnutil.SplitARN(dlqArn)
	if queueName != "" {
		if resolvedURL, err := sqsInvoker.GetQueueByName(context.Background(), dlqRegion, queueName); err == nil {
			queueURL = resolvedURL
		}
	}

	protocolMessage, err := extractProtocolMessage(msg, "sqs")
	if err != nil {
		return fmt.Errorf("extract protocol message for DLQ: %w", err)
	}

	payload := map[string]interface{}{
		"Type":      "Notification",
		"MessageId": msg.MessageId,
		"TopicArn":  msg.TopicArn,
		"Subject":   msg.Subject,
		"Message":   protocolMessage,
		"Timestamp": msg.PublishedTimestamp.Format(time.RFC3339Nano),
	}

	if len(msg.MessageAttributes) > 0 {
		attrs := make(map[string]interface{})
		for k, v := range msg.MessageAttributes {
			attrs[k] = map[string]interface{}{
				"Type":  v.Type,
				"Value": messageAttributeValue(v),
			}
		}
		payload["MessageAttributes"] = attrs
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal DLQ message: %w", err)
	}

	typedAttrs := make(map[string]eventbus.SQSMessageAttribute, len(msg.MessageAttributes))
	for k, v := range msg.MessageAttributes {
		ta := eventbus.SQSMessageAttribute{DataType: v.Type}
		if len(v.BinaryValue) > 0 {
			ta.BinaryValue = v.BinaryValue
		} else {
			ta.StringValue = v.StringValue
		}
		typedAttrs[k] = ta
	}

	opts := eventbus.SQSSendOptions{
		TypedMessageAttributes: typedAttrs,
		MessageGroupID:         msg.MessageGroupId,
		MessageDeduplicationID: msg.MessageDeduplicationId,
	}

	if _, _, err := sqsInvoker.SendMessage(context.Background(), dlqRegion, queueURL, string(body), opts); err != nil {
		return fmt.Errorf("send to DLQ %s: %w", dlqArn, err)
	}

	return nil
}

func (s *SNSService) deliverToSQS(msg *snsstore.Message, sub *snsstore.Subscription) error {
	if s.bus == nil {
		return nil
	}

	sqsInvoker := s.bus.SQSInvoker()
	if sqsInvoker == nil {
		return nil
	}

	queueURL := sub.Endpoint
	sqsRegion := ""
	if strings.HasPrefix(queueURL, "arn:") {
		queueName := arnutil.ExtractQueueNameFromARN(queueURL)
		_, _, sqsRegion, _, _ = arnutil.SplitARN(queueURL)
		if queueName != "" {
			if resolvedURL, err := sqsInvoker.GetQueueByName(context.Background(), sqsRegion, queueName); err == nil {
				queueURL = resolvedURL
			}
		}
	}

	queueARN, qErr := sqsInvoker.GetQueueARN(context.Background(), sqsRegion, queueURL)
	if qErr == nil && queueARN != "" {
		allowed, evalErr := s.bus.EvaluateTargetPolicy(context.Background(), queueARN, "sqs", "sns.amazonaws.com", "sqs:SendMessage", queueARN)
		if evalErr != nil {
			logs.Warn("resource policy evaluation failed for SQS delivery, dropping message",
				logs.String("queueArn", queueARN),
				logs.String("topicArn", msg.TopicArn),
				logs.String("error", evalErr.Error()))
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
	if sub.IsRawMessageDelivery() {
		body = protocolMessage
	} else {
		payload := map[string]interface{}{
			"Type":      "Notification",
			"MessageId": msg.MessageId,
			"TopicArn":  msg.TopicArn,
			"Subject":   msg.Subject,
			"Message":   protocolMessage,
			"Timestamp": msg.PublishedTimestamp.Format(time.RFC3339Nano),
		}

		if len(msg.MessageAttributes) > 0 {
			attrs := make(map[string]interface{})
			for k, v := range msg.MessageAttributes {
				attrs[k] = map[string]interface{}{
					"Type":  v.Type,
					"Value": messageAttributeValue(v),
				}
			}
			payload["MessageAttributes"] = attrs
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			logs.Warn("Failed to marshal SQS notification", logs.String("error", err.Error()))
			return fmt.Errorf("marshal SQS notification: %w", err)
		}
		body = string(jsonData)
	}

	typedAttrs := make(map[string]eventbus.SQSMessageAttribute, len(msg.MessageAttributes))
	for k, v := range msg.MessageAttributes {
		ta := eventbus.SQSMessageAttribute{DataType: v.Type}
		if len(v.BinaryValue) > 0 {
			ta.BinaryValue = v.BinaryValue
		} else {
			ta.StringValue = v.StringValue
		}
		typedAttrs[k] = ta
	}

	opts := eventbus.SQSSendOptions{
		TypedMessageAttributes: typedAttrs,
		MessageGroupID:         msg.MessageGroupId,
		MessageDeduplicationID: msg.MessageDeduplicationId,
	}

	if _, _, err := sqsInvoker.SendMessage(context.Background(), sqsRegion, queueURL, body, opts); err != nil {
		return fmt.Errorf("send to queue %s: %w", queueURL, err)
	}

	return nil
}

func (s *SNSService) deliverToHTTP(msg *snsstore.Message, sub *snsstore.Subscription, region string) error {
	protocolMessage, err := extractProtocolMessage(msg, sub.Protocol)
	if err != nil {
		return fmt.Errorf("extract protocol message for HTTP: %w", err)
	}

	var jsonData []byte
	if sub.IsRawMessageDelivery() {
		jsonData = []byte(protocolMessage)
	} else {
		payload := s.buildNotificationPayload(msg, sub, region, protocolMessage)

		if signature, certURL := s.signPayload(payload, region); signature != "" {
			payload["Signature"] = signature
			payload["SigningCertURL"] = certURL
		}

		jsonData, err = json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("marshal HTTP notification: %w", err)
		}
	}

	req, err := http.NewRequest("POST", sub.Endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-amz-sns-message-type", "Notification")
	req.Header.Set("x-amz-sns-message-id", msg.MessageId)
	req.Header.Set("x-amz-sns-topic-arn", msg.TopicArn)

	client := s.httpClient
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP delivery to %s: %w", sub.Endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP delivery to %s returned status %d", sub.Endpoint, resp.StatusCode)
	}

	logs.Debug("HTTP notification delivered",
		logs.String("endpoint", sub.Endpoint),
		logs.Int("status", resp.StatusCode))

	return nil
}

// buildNotificationPayload constructs the base SNS notification payload common to
// all delivery protocols. The caller may add protocol-specific fields (e.g. signature).
func (s *SNSService) buildNotificationPayload(msg *snsstore.Message, sub *snsstore.Subscription, region string, message string) map[string]interface{} {
	if region == "" {
		region = defaults.DefaultRegion
	}
	payload := map[string]interface{}{
		"Type":             "Notification",
		"MessageId":        msg.MessageId,
		"TopicArn":         msg.TopicArn,
		"Message":          message,
		"Timestamp":        msg.PublishedTimestamp.Format(time.RFC3339),
		"SignatureVersion": "2",
		"UnsubscribeURL":   fmt.Sprintf("https://sns.%s.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=%s", region, sub.SubscriptionArn),
	}

	if msg.Subject != "" {
		payload["Subject"] = msg.Subject
	}

	if len(msg.MessageAttributes) > 0 {
		attrs := make(map[string]interface{}, len(msg.MessageAttributes))
		for k, v := range msg.MessageAttributes {
			attrs[k] = map[string]interface{}{
				"Type":  v.Type,
				"Value": messageAttributeValue(v),
			}
		}
		payload["MessageAttributes"] = attrs
	}

	return payload
}

func (s *SNSService) signPayload(payload map[string]interface{}, region string) (string, string) {
	s.initSigningKey()
	if s.signingKey == nil {
		logs.Warn("Failed to get cached RSA key for signing")
		return "", ""
	}

	certURL := fmt.Sprintf("https://sns.%s.amazonaws.com/SimpleNotificationService-%x.pem", region, sha256.Sum256(s.signingCertPEM))

	var strToSign string
	if m, ok := payload["Message"].(string); ok {
		strToSign += "Message\n" + m + "\n"
	}
	if mid, ok := payload["MessageId"].(string); ok {
		strToSign += "MessageId\n" + mid + "\n"
	}
	if sub, ok := payload["Subject"].(string); ok && sub != "" {
		strToSign += "Subject\n" + sub + "\n"
	}
	if ts, ok := payload["Timestamp"].(string); ok {
		strToSign += "Timestamp\n" + ts + "\n"
	}
	if ta, ok := payload["TopicArn"].(string); ok {
		strToSign += "TopicArn\n" + ta + "\n"
	}
	if t, ok := payload["Type"].(string); ok {
		strToSign += "Type\n" + t
	}

	hashed := sha256.Sum256([]byte(strToSign))
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.signingKey, crypto.SHA256, hashed[:])
	if err != nil {
		logs.Warn("Failed to sign payload", logs.String("error", err.Error()))
		return "", ""
	}

	return base64.StdEncoding.EncodeToString(signature), certURL
}

func (s *SNSService) deliverToLambda(msg *snsstore.Message, sub *snsstore.Subscription, region string) error {
	if s.bus == nil {
		return nil
	}

	lambdaInvoker := s.bus.LambdaInvoker()
	if lambdaInvoker == nil {
		return nil
	}

	if region == "" {
		region = defaults.DefaultRegion
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

	var jsonData []byte
	if sub.IsRawMessageDelivery() {
		jsonData = []byte(protocolMessage)
	} else {
		snsPayload := s.buildNotificationPayload(msg, sub, region, protocolMessage)

		record := map[string]interface{}{
			"EventSource":          "aws:sns",
			"EventVersion":         "1.0",
			"EventSubscriptionArn": sub.SubscriptionArn,
			"Sns":                  snsPayload,
		}

		eventEnvelope := map[string]interface{}{
			"Records": []interface{}{record},
		}

		jsonData, err = json.Marshal(eventEnvelope)
		if err != nil {
			return fmt.Errorf("marshal Lambda event envelope: %w", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	functionName := sub.Endpoint
	if _, _, err := lambdaInvoker.InvokeForGateway(ctx, functionName, jsonData); err != nil {
		return fmt.Errorf("invoke Lambda %s (subscription %s, messageId %s): %w", functionName, sub.SubscriptionArn, msg.MessageId, err)
	}

	return nil
}

// PublishBatch publishes multiple messages to an SNS topic in a single request.
// https://docs.aws.amazon.com/sns/latest/api/API_PublishBatch.html
//
// Two-pass design (inside the Core) — all entries are validated before any
// entry is delivered, and the batch-level size check rejects the entire batch
// before any delivery or dedup recording.
func (s *SNSService) PublishBatch(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Wire parsing only; entry validation and execution live in the Core.
	entries := request.GetListParam(req.Parameters, "PublishBatchRequestEntries")

	return s.publishBatchCore(store, reqCtx, PublishBatchInput{
		TopicArn: request.GetParamLowerFirst(req.Parameters, "TopicArn"),
		Entries:  entries,
	})
}

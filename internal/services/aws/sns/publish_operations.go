package sns

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/common/defaults"
	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	"vorpalstacks/internal/store/aws/common"
	snsstore "vorpalstacks/internal/store/aws/sns"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// Publish publishes a message to an SNS topic.
// https://docs.aws.amazon.com/sns/latest/api/API_Publish.html
func (s *SNSService) Publish(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	topicArn := request.GetParamLowerFirst(req.Parameters, "TopicArn")
	message := request.GetParamLowerFirst(req.Parameters, "Message")
	subject := request.GetParamLowerFirst(req.Parameters, "Subject")
	messageStructure := request.GetParamLowerFirst(req.Parameters, "MessageStructure")

	messageGroupId := request.GetParamLowerFirst(req.Parameters, "MessageGroupId")
	messageDeduplicationId := request.GetParamLowerFirst(req.Parameters, "MessageDeduplicationId")

	// M14: TargetArn is an AWS-supported alternative to TopicArn. L4:
	// PhoneNumber is silently accepted but SMS is out-of-scope — reject
	// explicitly so callers get a clear error instead of silent success.
	targetArn := request.GetParamLowerFirst(req.Parameters, "TargetArn")
	phoneNumber := request.GetParamLowerFirst(req.Parameters, "PhoneNumber")

	if phoneNumber != "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "PhoneNumber is not supported (SMS sending is not available)", 400)
	}

	if topicArn == "" && targetArn == "" {
		return nil, awserrors.NewInvalidParameterException("TopicArn (or TargetArn) is required")
	}
	if topicArn == "" {
		topicArn = targetArn
	}
	if message == "" {
		return nil, awserrors.NewInvalidParameterException("Message is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	topic, err := store.GetTopic(topicArn)
	if err != nil {
		if err == snsstore.ErrTopicNotFound {
			return nil, ErrTopicNotFound
		}
		return nil, err
	}

	if err := validatePublishParams(topic.IsFifoTopic(), topic.IsContentBasedDeduplication(), message, subject, messageStructure, messageGroupId, messageDeduplicationId); err != nil {
		return nil, err
	}

	if topic.IsFifoTopic() && messageDeduplicationId == "" {
		messageDeduplicationId = generateContentBasedDeduplicationId(message)
	}

	messageId := uuid.New().String()

	msg := &snsstore.Message{
		MessageId:              messageId,
		TopicArn:               topic.Arn,
		Subject:                subject,
		Message:                message,
		MessageStructure:       messageStructure,
		MessageGroupId:         messageGroupId,
		MessageDeduplicationId: messageDeduplicationId,
	}

	if err := parseMessageAttributes(req.Parameters, msg); err != nil {
		return nil, err
	}

	// Atomically check for duplicates and record the dedup ID. This runs
	// after all validation to prevent cache leaks when validation fails,
	// and is atomic to eliminate the TOCTOU race between separate
	// check (RLock) and record (Lock) operations.
	if topic.IsFifoTopic() && messageDeduplicationId != "" {
		if existingMsgID, isDuplicate := store.CheckAndRecordDeduplication(topicArn, messageDeduplicationId, messageId); isDuplicate {
			return map[string]interface{}{
				"MessageId": existingMsgID,
			}, nil
		}
	}

	msg.PublishedTimestamp = time.Now().UTC()
	msg.ReceivedTimestamp = time.Now().UTC()

	subscriptions, err := store.ListSubscriptionsByTopic(topicArn, common.ListOptions{})
	if err == nil && len(subscriptions.Items) > 0 {
		msgCopy := *msg
		subsCopy := make([]*snsstore.Subscription, len(subscriptions.Items))
		for i, sub := range subscriptions.Items {
			subCopy := *sub
			subsCopy[i] = &subCopy
		}
		region := reqCtx.GetRegion()

		if s.bus != nil {
			// Serialise message attributes to raw JSON for transport through
			// the event bus (which must not depend on store-layer types).
			var msgAttrs map[string]json.RawMessage
			if len(msg.MessageAttributes) > 0 {
				msgAttrs = make(map[string]json.RawMessage, len(msg.MessageAttributes))
				for k, v := range msg.MessageAttributes {
					raw, err := json.Marshal(v)
					if err == nil {
						msgAttrs[k] = raw
					}
				}
			}
			snsEvt := &eventbus.SNSDeliveryEvent{
				TopicARN:          topic.Arn,
				MessageID:         msg.MessageId,
				Message:           message,
				Subject:           subject,
				MessageStructure:  messageStructure,
				MessageGroupId:    messageGroupId,
				MessageAttributes: msgAttrs,
			}
			snsEvt.Region = region
			if err := s.bus.Publish(context.Background(), snsEvt); err != nil {
				logs.Warn("Failed to publish SNS delivery event to event bus; message is stored but subscribers may not be notified",
					logs.String("topicArn", topicArn),
					logs.String("messageId", messageId),
					logs.Err(err))
			}
		} else {
			s.deliverAsync(&msgCopy, subsCopy, region)
		}
	}

	result := map[string]interface{}{
		"MessageId": messageId,
	}
	if topic.IsFifoTopic() {
		result["SequenceNumber"] = store.GetNextSequenceNumber(topicArn, messageGroupId)
	}
	return result, nil
}

func generateContentBasedDeduplicationId(message string) string {
	hash := sha256.Sum256([]byte(message))
	return hex.EncodeToString(hash[:32])
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
			// H3: previously this was `logs.Warn + continue` (silent drop).
			// Return a delivery error so the message routes to the DLQ when
			// a RedrivePolicy is configured, rather than being silently lost.
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

// batchValidatedEntry holds a single PublishBatch entry that has passed all
// validation in the first pass. The second pass uses these entries for
// delivery without re-parsing. This separation ensures that a batch-level
// error (e.g. BatchRequestTooLong) rejects the entire batch before any entry
// is delivered, preserving the atomicity contract of the top-level error.
type batchValidatedEntry struct {
	id                     string
	message                string
	subject                string
	messageStructure       string
	messageGroupId         string
	messageDeduplicationId string
	msgAttrs               map[string]*snsstore.MessageAttribute
	isDuplicate            bool
	existingMsgID          string
}

// PublishBatch publishes multiple messages to an SNS topic in a single request.
// https://docs.aws.amazon.com/sns/latest/api/API_PublishBatch.html
//
// Two-pass design:
//   - Pass 1 validates every entry (params, attributes, dedup read-check,
//     size accumulation). Per-entry failures go into the Failed list; valid
//     entries are collected for Pass 2.
//   - After Pass 1 the total batch size is checked. If it exceeds the limit
//     the entire batch is rejected (BatchRequestTooLong) — no entry has been
//     delivered or dedup-recorded.
//   - Pass 2 delivers each validated entry, records dedup IDs, and collects
//     results.
func (s *SNSService) PublishBatch(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	topicArn := request.GetParamLowerFirst(req.Parameters, "TopicArn")
	if topicArn == "" {
		return nil, awserrors.NewInvalidParameterException("TopicArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	topic, err := store.GetTopic(topicArn)
	if err != nil {
		if err == snsstore.ErrTopicNotFound {
			return nil, ErrTopicNotFound
		}
		return nil, err
	}

	entryMaps := request.GetListParam(req.Parameters, "PublishBatchRequestEntries")
	if len(entryMaps) == 0 {
		return nil, awserrors.NewAWSError("EmptyBatchRequest", "Batch request does not contain any entries", 400)
	}
	if len(entryMaps) > maxBatchEntries {
		return nil, ErrTooManyEntriesInBatch
	}

	successful := make([]map[string]interface{}, 0)
	failed := make([]map[string]interface{}, 0)
	seenIds := make(map[string]bool, len(entryMaps))
	validated := make([]batchValidatedEntry, 0, len(entryMaps))
	batchTotalSize := 0

	subscriptions, err := store.ListSubscriptionsByTopic(topicArn, common.ListOptions{})
	if err != nil {
		return nil, err
	}
	region := reqCtx.GetRegion()

	// --- Pass 1: validate all entries (no side effects) ---
	for _, entryMap := range entryMaps {

		id, _ := entryMap["Id"].(string)
		if id == "" {
			failed = append(failed, map[string]interface{}{
				"Id":          "",
				"Code":        "InvalidBatchEntryId",
				"Message":     "A batch entry Id is required",
				"SenderFault": true,
			})
			continue
		}

		if seenIds[id] {
			return nil, ErrBatchEntryIdsNotDistinct
		}
		seenIds[id] = true

		message, _ := entryMap["Message"].(string)
		if message == "" {
			failed = append(failed, map[string]interface{}{
				"Id":          id,
				"Code":        "InvalidParameter",
				"Message":     "Message is required",
				"SenderFault": true,
			})
			continue
		}

		subject, _ := entryMap["Subject"].(string)
		messageGroupId, _ := entryMap["MessageGroupId"].(string)
		messageDeduplicationId, _ := entryMap["MessageDeduplicationId"].(string)
		messageStructure, _ := entryMap["MessageStructure"].(string)

		if err := validatePublishParams(topic.IsFifoTopic(), topic.IsContentBasedDeduplication(), message, subject, messageStructure, messageGroupId, messageDeduplicationId); err != nil {
			failed = append(failed, map[string]interface{}{
				"Id":          id,
				"Code":        "InvalidParameter",
				"Message":     err.Error(),
				"SenderFault": true,
			})
			continue
		}

		entry := batchValidatedEntry{
			id:               id,
			message:          message,
			subject:          subject,
			messageStructure: messageStructure,
			messageGroupId:   messageGroupId,
		}

		if topic.IsFifoTopic() {
			if messageDeduplicationId == "" {
				messageDeduplicationId = generateContentBasedDeduplicationId(message)
			}
			entry.messageDeduplicationId = messageDeduplicationId
			if existingMsgID, isDuplicate := store.CheckDeduplication(topicArn, messageDeduplicationId); isDuplicate {
				entry.isDuplicate = true
				entry.existingMsgID = existingMsgID
			}
		}

		if !entry.isDuplicate {
			msg := &snsstore.Message{}
			if err := parseMessageAttributes(entryMap, msg); err != nil {
				failed = append(failed, map[string]interface{}{
					"Id":          id,
					"Code":        "InvalidParameter",
					"Message":     err.Error(),
					"SenderFault": true,
				})
				continue
			}
			entry.msgAttrs = msg.MessageAttributes
			batchTotalSize += messageEntrySize(message, subject, entry.msgAttrs)
		}

		validated = append(validated, entry)
	}

	// Batch-level size check: reject the entire batch before any delivery.
	if batchTotalSize > maxBatchTotalSize {
		return nil, awserrors.NewAWSError("BatchRequestTooLong", fmt.Sprintf("Total batch request size %d exceeds maximum %d", batchTotalSize, maxBatchTotalSize), 400)
	}

	// --- Pass 2: deliver validated entries ---
	for _, entry := range validated {
		if entry.isDuplicate {
			successful = append(successful, map[string]interface{}{
				"Id":        entry.id,
				"MessageId": entry.existingMsgID,
			})
			continue
		}

		messageId := uuid.New().String()

		if topic.IsFifoTopic() && entry.messageDeduplicationId != "" {
			existingMsgID, isDuplicate := store.CheckAndRecordDeduplication(topicArn, entry.messageDeduplicationId, messageId)
			if isDuplicate {
				// A concurrent publish with the same dedup ID won the
				// race between Pass 1 and Pass 2. Return the existing
				// message ID without delivering.
				successful = append(successful, map[string]interface{}{
					"Id":        entry.id,
					"MessageId": existingMsgID,
				})
				continue
			}
		}

		msg := &snsstore.Message{
			MessageId:              messageId,
			TopicArn:               topicArn,
			Subject:                entry.subject,
			Message:                entry.message,
			MessageStructure:       entry.messageStructure,
			MessageGroupId:         entry.messageGroupId,
			MessageDeduplicationId: entry.messageDeduplicationId,
			MessageAttributes:      entry.msgAttrs,
		}
		msg.PublishedTimestamp = time.Now().UTC()
		msg.ReceivedTimestamp = time.Now().UTC()

		if len(subscriptions.Items) > 0 {
			msgCopy := *msg
			subsCopy := make([]*snsstore.Subscription, len(subscriptions.Items))
			for j, sub := range subscriptions.Items {
				subCopy := *sub
				subsCopy[j] = &subCopy
			}

			if s.bus != nil {
				var msgAttrs map[string]json.RawMessage
				if len(msg.MessageAttributes) > 0 {
					msgAttrs = make(map[string]json.RawMessage, len(msg.MessageAttributes))
					for k, v := range msg.MessageAttributes {
						raw, err := json.Marshal(v)
						if err == nil {
							msgAttrs[k] = raw
						}
					}
				}
				snsEvt := &eventbus.SNSDeliveryEvent{
					TopicARN:          topicArn,
					MessageID:         messageId,
					Message:           entry.message,
					Subject:           entry.subject,
					MessageStructure:  entry.messageStructure,
					MessageGroupId:    entry.messageGroupId,
					MessageAttributes: msgAttrs,
				}
				snsEvt.Region = region
				if err := s.bus.Publish(context.Background(), snsEvt); err != nil {
					logs.Warn("Failed to publish SNS event", logs.Err(err))
				}
			} else {
				s.deliverAsync(&msgCopy, subsCopy, region)
			}
		}

		result := map[string]interface{}{
			"Id":        entry.id,
			"MessageId": messageId,
		}

		if topic.IsFifoTopic() {
			result["SequenceNumber"] = store.GetNextSequenceNumber(topicArn, entry.messageGroupId)
		}

		successful = append(successful, result)
	}

	return map[string]interface{}{
		"Successful": successful,
		"Failed":     failed,
	}, nil
}

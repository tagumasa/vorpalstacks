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

	if topicArn == "" {
		return nil, awserrors.NewInvalidParameterException("TopicArn is required")
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

	if topic.FifoTopic {
		if messageGroupId == "" {
			return nil, awserrors.NewInvalidParameterException("MessageGroupId is required for FIFO topics")
		}

		if messageDeduplicationId != "" {
			if existingMsgID, isDuplicate := store.CheckDeduplication(topicArn, messageDeduplicationId); isDuplicate {
				return map[string]interface{}{
					"MessageId": existingMsgID,
				}, nil
			}
		} else if !topic.ContentBasedDeduplication {
			return nil, awserrors.NewInvalidParameterException("MessageDeduplicationId is required when ContentBasedDeduplication is false")
		} else {
			messageDeduplicationId = generateContentBasedDeduplicationId(message)
			if existingMsgID, isDuplicate := store.CheckDeduplication(topicArn, messageDeduplicationId); isDuplicate {
				return map[string]interface{}{
					"MessageId": existingMsgID,
				}, nil
			}
		}
	}

	messageId := uuid.New().String()

	if topic.FifoTopic && messageDeduplicationId != "" {
		store.RecordDeduplication(topicArn, messageDeduplicationId, messageId)
	}

	msg := &snsstore.Message{
		MessageId:              messageId,
		TopicArn:               topic.Arn,
		Subject:                subject,
		Message:                message,
		MessageStructure:       messageStructure,
		MessageGroupId:         messageGroupId,
		MessageDeduplicationId: messageDeduplicationId,
	}

	parseMessageAttributes(req.Parameters, msg)

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
				TopicARN:          topicArn,
				MessageID:         messageId,
				Message:           message,
				Subject:           subject,
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
	if topic.FifoTopic {
		result["SequenceNumber"] = store.GetNextSequenceNumber(topicArn, messageGroupId)
	}
	return result, nil
}

// parseMessageAttributes extracts SNS message attributes from a request params
// map (either top-level req.Parameters or a PublishBatch entry map) and
// populates the Message's MessageAttributes field.
func parseMessageAttributes(params map[string]interface{}, msg *snsstore.Message) {
	var attrs map[string]interface{}
	for _, key := range []string{"MessageAttributes", "messageAttributes"} {
		if m, ok := params[key].(map[string]interface{}); ok {
			attrs = m
			break
		}
	}
	if attrs == nil {
		return
	}

	msg.MessageAttributes = make(map[string]*snsstore.MessageAttribute, len(attrs))
	for k, v := range attrs {
		attrMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		attr := &snsstore.MessageAttribute{
			Type:        firstString(attrMap, "DataType", "dataType"),
			StringValue: firstString(attrMap, "StringValue", "stringValue"),
		}
		if raw := firstString(attrMap, "BinaryValue", "binaryValue"); raw != "" {
			if decoded, err := base64.StdEncoding.DecodeString(raw); err == nil {
				attr.BinaryValue = decoded
			}
		}
		msg.MessageAttributes[k] = attr
	}
}

// firstString returns the first non-empty string value found for any of the
// given keys in the map.
func firstString(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k].(string); ok && v != "" {
			return v
		}
	}
	return ""
}

// messageAttributeValue returns the serialisable value for an SNS message
// attribute. String/Number/String.Array types return the string value;
// Binary types return the base64-encoded representation.
func messageAttributeValue(attr *snsstore.MessageAttribute) string {
	if len(attr.BinaryValue) > 0 {
		return base64.StdEncoding.EncodeToString(attr.BinaryValue)
	}
	return attr.StringValue
}

func generateContentBasedDeduplicationId(message string) string {
	hash := sha256.Sum256([]byte(message))
	return hex.EncodeToString(hash[:32])
}

func extractProtocolMessage(msg *snsstore.Message, protocol string) string {
	if msg.MessageStructure != "json" {
		return msg.Message
	}

	var msgMap map[string]string
	if err := json.Unmarshal([]byte(msg.Message), &msgMap); err != nil {
		return msg.Message
	}

	if protocolMsg, ok := msgMap[protocol]; ok {
		return protocolMsg
	}
	if defaultMsg, ok := msgMap["default"]; ok {
		return defaultMsg
	}

	return msg.Message
}

func (s *SNSService) deliverToSubscriptions(msg *snsstore.Message, subscriptions []*snsstore.Subscription, region string) {
	for _, sub := range subscriptions {
		if sub.PendingConfirmation {
			continue
		}

		if !matchFilterPolicy(sub.GetFilterPolicy(), msg.MessageAttributes) {
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
	if queueName != "" {
		if resolvedURL, err := sqsInvoker.GetQueueByName(context.Background(), queueName); err == nil {
			queueURL = resolvedURL
		}
	}

	protocolMessage := extractProtocolMessage(msg, "sqs")

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

	if _, _, err := sqsInvoker.SendMessage(context.Background(), queueURL, string(body), opts); err != nil {
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
	if strings.HasPrefix(queueURL, "arn:") {
		queueName := arnutil.ExtractQueueNameFromARN(queueURL)
		if queueName != "" {
			if resolvedURL, err := sqsInvoker.GetQueueByName(context.Background(), queueName); err == nil {
				queueURL = resolvedURL
			}
		}
	}

	queueARN, qErr := sqsInvoker.GetQueueARN(context.Background(), queueURL)
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

	protocolMessage := extractProtocolMessage(msg, "sqs")

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

	if _, _, err := sqsInvoker.SendMessage(context.Background(), queueURL, body, opts); err != nil {
		return fmt.Errorf("send to queue %s: %w", queueURL, err)
	}

	return nil
}

func (s *SNSService) deliverToHTTP(msg *snsstore.Message, sub *snsstore.Subscription, region string) error {
	protocolMessage := extractProtocolMessage(msg, sub.Protocol)
	payload := s.buildNotificationPayload(msg, sub, region, protocolMessage)

	if signature, certURL := s.signPayload(payload, region); signature != "" {
		payload["Signature"] = signature
		payload["SigningCertURL"] = certURL
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal HTTP notification: %w", err)
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

	protocolMessage := extractProtocolMessage(msg, "lambda")
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

	jsonData, err := json.Marshal(eventEnvelope)
	if err != nil {
		return fmt.Errorf("marshal Lambda event envelope: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	functionName := sub.Endpoint
	if _, _, err := lambdaInvoker.InvokeForGateway(ctx, functionName, jsonData); err != nil {
		return fmt.Errorf("invoke Lambda %s: %w", functionName, err)
	}

	return nil
}

// PublishBatch publishes multiple messages to an SNS topic in a single request.
// https://docs.aws.amazon.com/sns/latest/api/API_PublishBatch.html
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
		return nil, awserrors.NewInvalidParameterException("PublishBatchRequestEntries is required")
	}
	if len(entryMaps) > 10 {
		return nil, awserrors.NewInvalidParameterException("PublishBatchRequestEntries cannot exceed 10 entries")
	}

	successful := make([]map[string]interface{}, 0)
	failed := make([]map[string]interface{}, 0)

	subscriptions, err := store.ListSubscriptionsByTopic(topicArn, common.ListOptions{})
	if err != nil {
		return nil, err
	}
	region := reqCtx.GetRegion()

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

		if topic.FifoTopic {
			if messageGroupId == "" {
				failed = append(failed, map[string]interface{}{
					"Id":          id,
					"Code":        "InvalidParameter",
					"Message":     "MessageGroupId is required for FIFO topics",
					"SenderFault": true,
				})
				continue
			}

			if messageDeduplicationId == "" {
				if topic.ContentBasedDeduplication {
					messageDeduplicationId = generateContentBasedDeduplicationId(message)
				} else {
					failed = append(failed, map[string]interface{}{
						"Id":          id,
						"Code":        "InvalidParameter",
						"Message":     "MessageDeduplicationId is required when ContentBasedDeduplication is false",
						"SenderFault": true,
					})
					continue
				}
			}

			if existingMsgID, isDuplicate := store.CheckDeduplication(topicArn, messageDeduplicationId); isDuplicate {
				successful = append(successful, map[string]interface{}{
					"Id":        id,
					"MessageId": existingMsgID,
				})
				continue
			}
		}

		messageId := uuid.New().String()

		if topic.FifoTopic && messageDeduplicationId != "" {
			store.RecordDeduplication(topicArn, messageDeduplicationId, messageId)
		}

		msg := &snsstore.Message{
			MessageId:              messageId,
			TopicArn:               topicArn,
			Subject:                subject,
			Message:                message,
			MessageStructure:       messageStructure,
			MessageGroupId:         messageGroupId,
			MessageDeduplicationId: messageDeduplicationId,
		}

		parseMessageAttributes(entryMap, msg)

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
					Message:           message,
					Subject:           subject,
					MessageGroupId:    messageGroupId,
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
			"Id":        id,
			"MessageId": messageId,
		}

		if topic.FifoTopic {
			sequenceNumber := store.GetNextSequenceNumber(topicArn, messageGroupId)
			result["SequenceNumber"] = sequenceNumber
		}

		successful = append(successful, result)
	}

	return map[string]interface{}{
		"Successful": successful,
		"Failed":     failed,
	}, nil
}

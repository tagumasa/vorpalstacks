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
			return nil, NewTopicNotFoundException()
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

	attrs, ok := req.Parameters["messageAttributes"].(map[string]interface{})
	if !ok {
		attrs, ok = req.Parameters["MessageAttributes"].(map[string]interface{})
	}
	if ok {
		msg.MessageAttributes = make(map[string]*snsstore.MessageAttribute)
		for k, v := range attrs {
			if attrMap, ok := v.(map[string]interface{}); ok {
				attr := &snsstore.MessageAttribute{}
				if dataType, ok := attrMap["dataType"].(string); ok {
					attr.Type = dataType
				}
				if dataType, ok := attrMap["DataType"].(string); ok {
					attr.Type = dataType
				}
				if stringValue, ok := attrMap["stringValue"].(string); ok {
					attr.StringValue = stringValue
				}
				if stringValue, ok := attrMap["StringValue"].(string); ok {
					attr.StringValue = stringValue
				}
				if binaryValue, ok := attrMap["binaryValue"].([]byte); ok {
					attr.BinaryValue = binaryValue
				}
				if binaryValue, ok := attrMap["BinaryValue"].([]byte); ok {
					attr.BinaryValue = binaryValue
				}
				msg.MessageAttributes[k] = attr
			}
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
		switch sub.Protocol {
		case "sqs":
			s.deliverToSQS(msg, sub)
		case "http", "https":
			s.deliverToHTTP(msg, sub, region)
		case "lambda":
			s.deliverToLambda(msg, sub, region)
		}
	}
}

func (s *SNSService) deliverToSQS(msg *snsstore.Message, sub *snsstore.Subscription) {
	if s.bus == nil {
		return
	}

	sqsInvoker := s.bus.SQSInvoker()
	if sqsInvoker == nil {
		return
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
			return
		}
		if !allowed {
			return
		}
	}

	protocolMessage := extractProtocolMessage(msg, "sqs")

	var body string
	if sub.RawMessageDelivery {
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
					"Value": v.StringValue,
				}
			}
			payload["MessageAttributes"] = attrs
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			logs.Warn("Failed to marshal SQS notification", logs.String("error", err.Error()))
			return
		}
		body = string(jsonData)
	}

	msgAttrs := make(map[string]string)
	for k, v := range msg.MessageAttributes {
		msgAttrs[k] = v.StringValue
	}

	if msg.MessageGroupId != "" || msg.MessageDeduplicationId != "" {
		logs.Warn("SNS to SQS FIFO delivery: MessageGroupId/MessageDeduplicationId not supported by SQSInvoker; attributes will be dropped",
			logs.String("topicArn", msg.TopicArn),
			logs.String("messageId", msg.MessageId))
	}

	if _, _, err := sqsInvoker.SendMessage(context.Background(), queueURL, body, 0, msgAttrs); err != nil {
		logs.Error("Failed to deliver SNS message to SQS queue — message may be permanently lost",
			logs.String("queueURL", queueURL),
			logs.String("messageId", msg.MessageId),
			logs.String("topicArn", msg.TopicArn),
			logs.Err(err))
	}
}

func (s *SNSService) deliverToHTTP(msg *snsstore.Message, sub *snsstore.Subscription, region string) {
	protocolMessage := extractProtocolMessage(msg, sub.Protocol)
	payload := s.buildNotificationPayloadWithMessage(msg, sub, region, protocolMessage)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		logs.Debug("Failed to marshal HTTP notification", logs.String("error", err.Error()))
		return
	}

	req, err := http.NewRequest("POST", sub.Endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		logs.Debug("Failed to create HTTP request", logs.String("error", err.Error()))
		return
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
		logs.Warn("HTTP delivery failed", logs.String("endpoint", sub.Endpoint), logs.String("error", err.Error()))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		logs.Warn("HTTP delivery returned error status", logs.String("endpoint", sub.Endpoint), logs.Int("status", resp.StatusCode))
		return
	}

	logs.Debug("HTTP notification delivered",
		logs.String("endpoint", sub.Endpoint),
		logs.Int("status", resp.StatusCode))
}

func (s *SNSService) buildNotificationPayloadWithMessage(msg *snsstore.Message, sub *snsstore.Subscription, region string, message string) map[string]interface{} {
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
		attrs := make(map[string]interface{})
		for k, v := range msg.MessageAttributes {
			attr := map[string]interface{}{
				"Type": v.Type,
			}
			if v.StringValue != "" {
				attr["Value"] = v.StringValue
			}
			attrs[k] = attr
		}
		payload["MessageAttributes"] = attrs
	}

	signature, certURL := s.signPayload(payload, region)
	if signature != "" {
		payload["Signature"] = signature
		payload["SigningCertURL"] = certURL
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
		strToSign += "Type\n" + t + "\n"
	}
	if unsubURL, ok := payload["UnsubscribeURL"].(string); ok {
		strToSign += "UnsubscribeURL\n" + unsubURL + "\n"
	}

	hashed := sha256.Sum256([]byte(strToSign))
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.signingKey, crypto.SHA256, hashed[:])
	if err != nil {
		logs.Warn("Failed to sign payload", logs.String("error", err.Error()))
		return "", ""
	}

	return base64.StdEncoding.EncodeToString(signature), certURL
}

func (s *SNSService) deliverToLambda(msg *snsstore.Message, sub *snsstore.Subscription, region string) {
	if s.bus == nil {
		return
	}

	lambdaInvoker := s.bus.LambdaInvoker()
	if lambdaInvoker == nil {
		return
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
		logs.Warn("resource policy evaluation failed for Lambda delivery, dropping message",
			logs.String("functionArn", functionARN),
			logs.String("topicArn", msg.TopicArn),
			logs.String("error", evalErr.Error()))
		return
	}
	if !allowed {
		return
	}

	protocolMessage := extractProtocolMessage(msg, "lambda")

	certURL := fmt.Sprintf("https://sns.%s.amazonaws.com/SimpleNotificationService-%x.pem", region, sha256.Sum256(s.signingCertPEM))
	s.initSigningKey()

	snsPayload := map[string]interface{}{
		"Type":             "Notification",
		"MessageId":        msg.MessageId,
		"TopicArn":         msg.TopicArn,
		"Subject":          msg.Subject,
		"Message":          protocolMessage,
		"Timestamp":        msg.PublishedTimestamp.Format(time.RFC3339Nano),
		"SignatureVersion": "2",
		"SigningCertURL":   certURL,
		"UnsubscribeURL":   fmt.Sprintf("https://sns.%s.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=%s", region, sub.SubscriptionArn),
	}

	if len(msg.MessageAttributes) > 0 {
		attrs := make(map[string]interface{})
		for k, v := range msg.MessageAttributes {
			attrs[k] = map[string]interface{}{
				"Type":  v.Type,
				"Value": v.StringValue,
			}
		}
		snsPayload["MessageAttributes"] = attrs
	}

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
		logs.Warn("Failed to marshal Lambda event envelope", logs.String("error", err.Error()))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	functionName := sub.Endpoint
	_, _, err = lambdaInvoker.InvokeForGateway(ctx, functionName, jsonData)
	if err != nil {
		logs.Warn("Lambda invocation failed", logs.String("function", functionName), logs.String("error", err.Error()))
	}
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
			return nil, NewTopicNotFoundException()
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

		attrs, ok := entryMap["MessageAttributes"].(map[string]interface{})
		if !ok {
			attrs, ok = entryMap["messageAttributes"].(map[string]interface{})
		}
		if ok {
			msg.MessageAttributes = make(map[string]*snsstore.MessageAttribute)
			for k, v := range attrs {
				if attrMap, ok := v.(map[string]interface{}); ok {
					attr := &snsstore.MessageAttribute{}
					if dataType, ok := attrMap["dataType"].(string); ok {
						attr.Type = dataType
					}
					if dataType, ok := attrMap["DataType"].(string); ok {
						attr.Type = dataType
					}
					if stringValue, ok := attrMap["stringValue"].(string); ok {
						attr.StringValue = stringValue
					}
					if stringValue, ok := attrMap["StringValue"].(string); ok {
						attr.StringValue = stringValue
					}
					if binaryValue, ok := attrMap["binaryValue"].([]byte); ok {
						attr.BinaryValue = binaryValue
					}
					if binaryValue, ok := attrMap["BinaryValue"].([]byte); ok {
						attr.BinaryValue = binaryValue
					}
					msg.MessageAttributes[k] = attr
				}
			}
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
				snsEvt := &eventbus.SNSDeliveryEvent{
					TopicARN:       topicArn,
					MessageID:      messageId,
					Message:        message,
					Subject:        subject,
					MessageGroupId: messageGroupId,
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

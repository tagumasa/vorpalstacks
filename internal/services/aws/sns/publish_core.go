package sns

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	"vorpalstacks/internal/store/aws/common"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

// PublishInput carries the parsed parameters for publishing a message.
// Parameters holds the raw request map because the message-attribute wire
// formats (JSON map and Query flat keys) are parsed inside the Core at their
// original validation position.
type PublishInput struct {
	TopicArn               string
	TargetArn              string
	PhoneNumber            string
	Message                string
	Subject                string
	MessageStructure       string
	MessageGroupId         string
	MessageDeduplicationId string
	Parameters             map[string]interface{}
}

// PublishBatchInput carries the topic ARN and the raw batch entries for a
// PublishBatch request. Entries are extracted on the wire side (infallible);
// every entry-level validation runs inside the Core in the original order.
type PublishBatchInput struct {
	TopicArn string
	Entries  []map[string]interface{}
}

// messageEntrySize estimates the serialised wire size of a single Publish or
// PublishBatch entry: message body + subject + all message attributes
// (name + DataType + StringValue or BinaryValue).  AWS counts the full
// serialised request toward the 256 KB batch limit, so excluding attributes
// under-estimates and allows oversized batches through.
func messageEntrySize(message, subject string, attrs map[string]*snsstore.MessageAttribute) int {
	size := len(message) + len(subject)
	for name, attr := range attrs {
		size += len(name)
		size += len(attr.Type)
		size += len(attr.StringValue)
		if attr.BinaryValue != nil {
			size += base64.StdEncoding.EncodedLen(len(attr.BinaryValue))
		}
	}
	return size
}

// parseMessageAttributes extracts SNS message attributes from a request params
// map and populates the Message's MessageAttributes field. Returns an error
// (fail-closed) when an attribute entry is malformed or has an invalid
// DataType instead of silently skipping it.
//
// The function enforces the AWS-documented message-attribute limits:
// maximum 10 attributes per message, String values up to 256 characters,
// Binary values up to 256 bytes, and attribute names matching
// [a-zA-Z0-9_.-]{1,256}.
func parseMessageAttributes(params map[string]interface{}, msg *snsstore.Message) error {
	var attrs map[string]interface{}
	for _, key := range []string{"MessageAttributes", "messageAttributes"} {
		if m, ok := params[key].(map[string]interface{}); ok {
			attrs = m
			break
		}
	}

	// Also handle the AWS Query API flat-key format. The SDK sends message
	// attributes as:
	//   MessageAttributes.entry.N.Name=<name>
	//   MessageAttributes.entry.N.Value.DataType=<type>
	//   MessageAttributes.entry.N.Value.StringValue=<value>
	// The HTTP query parser stores these as flat string keys, so we parse
	// them manually.
	if attrs == nil {
		attrs = make(map[string]interface{})
		for i := 1; i <= 30; i++ {
			name := request.GetStringParam(params, fmt.Sprintf("MessageAttributes.entry.%d.Name", i))
			if name == "" {
				break
			}
			attrs[name] = map[string]interface{}{
				"DataType":    request.GetStringParam(params, fmt.Sprintf("MessageAttributes.entry.%d.Value.DataType", i)),
				"StringValue": request.GetStringParam(params, fmt.Sprintf("MessageAttributes.entry.%d.Value.StringValue", i)),
				"BinaryValue": request.GetStringParam(params, fmt.Sprintf("MessageAttributes.entry.%d.Value.BinaryValue", i)),
			}
		}
		if len(attrs) == 0 {
			attrs = nil
		}
	}

	if attrs == nil {
		return nil
	}

	// Maximum 10 message attributes per AWS spec.
	if len(attrs) > maxMessageAttributes {
		return awserrors.NewInvalidParameterException(fmt.Sprintf("Too many message attributes: %d (maximum %d)", len(attrs), maxMessageAttributes))
	}

	msg.MessageAttributes = make(map[string]*snsstore.MessageAttribute, len(attrs))
	for k, v := range attrs {
		// Validate attribute name format per AWS spec.
		if err := validateMessageAttributeName(k); err != nil {
			return err
		}

		attrMap, ok := v.(map[string]interface{})
		if !ok {
			return awserrors.NewInvalidParameterException(fmt.Sprintf("Invalid message attribute %q: value must be a map", k))
		}

		dataType := firstString(attrMap, "DataType", "dataType")
		if dataType == "" {
			return awserrors.NewInvalidParameterException(fmt.Sprintf("Invalid message attribute %q: DataType is required", k))
		}
		if !validMessageAttributeDataTypes[dataType] {
			return awserrors.NewInvalidParameterException(fmt.Sprintf("Invalid message attribute %q: DataType %q is not valid (String, Number, String.Array, Binary)", k, dataType))
		}

		attr := &snsstore.MessageAttribute{
			Type:        dataType,
			StringValue: firstString(attrMap, "StringValue", "stringValue"),
		}
		if raw := firstString(attrMap, "BinaryValue", "binaryValue"); raw != "" {
			decoded, err := base64.StdEncoding.DecodeString(raw)
			if err != nil {
				return awserrors.NewInvalidParameterException(fmt.Sprintf("Invalid message attribute %q: BinaryValue is not valid base64", k))
			}
			attr.BinaryValue = decoded
		}

		// Validate attribute value sizes per AWS spec.
		if err := validateMessageAttributeLimits(k, attr.StringValue, attr.BinaryValue); err != nil {
			return err
		}

		msg.MessageAttributes[k] = attr
	}

	return nil
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

// generateContentBasedDeduplicationId derives the deduplication ID from the
// SHA-256 hash of the message body when a FIFO topic enables
// ContentBasedDeduplication.
func generateContentBasedDeduplicationId(message string) string {
	hash := sha256.Sum256([]byte(message))
	return hex.EncodeToString(hash[:32])
}

// publishCore is the single validation and persistence path for Publish. It
// needs the request context for the delivery region.
func (s *SNSService) publishCore(store snsstore.SNSStoreInterface, reqCtx *request.RequestContext, in PublishInput) (interface{}, error) {
	// TargetArn is an AWS-supported alternative to TopicArn. PhoneNumber
	// is silently accepted by AWS but SMS sending is out-of-scope here —
	// reject it explicitly so callers get a clear error instead of silent
	// success.
	if in.PhoneNumber != "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "PhoneNumber is not supported (SMS sending is not available)", 400)
	}

	if in.TopicArn == "" && in.TargetArn == "" {
		return nil, awserrors.NewInvalidParameterException("TopicArn (or TargetArn) is required")
	}
	if in.TopicArn == "" {
		in.TopicArn = in.TargetArn
	}
	if in.Message == "" {
		return nil, awserrors.NewInvalidParameterException("Message is required")
	}

	topic, err := store.GetTopic(in.TopicArn)
	if err != nil {
		if err == snsstore.ErrTopicNotFound {
			return nil, ErrTopicNotFound
		}
		return nil, err
	}

	if err := validatePublishParams(topic.IsFifoTopic(), topic.IsContentBasedDeduplication(), in.Message, in.Subject, in.MessageStructure, in.MessageGroupId, in.MessageDeduplicationId); err != nil {
		return nil, err
	}

	if topic.IsFifoTopic() && in.MessageDeduplicationId == "" {
		in.MessageDeduplicationId = generateContentBasedDeduplicationId(in.Message)
	}

	messageId := uuid.New().String()

	msg := &snsstore.Message{
		MessageId:              messageId,
		TopicArn:               topic.Arn,
		Subject:                in.Subject,
		Message:                in.Message,
		MessageStructure:       in.MessageStructure,
		MessageGroupId:         in.MessageGroupId,
		MessageDeduplicationId: in.MessageDeduplicationId,
	}

	if err := parseMessageAttributes(in.Parameters, msg); err != nil {
		return nil, err
	}

	// Atomically check for duplicates and record the dedup ID. This runs
	// after all validation to prevent cache leaks when validation fails,
	// and is atomic to eliminate the TOCTOU race between separate
	// check (RLock) and record (Lock) operations.
	if topic.IsFifoTopic() && in.MessageDeduplicationId != "" {
		if existingMsgID, isDuplicate := store.CheckAndRecordDeduplication(in.TopicArn, in.MessageDeduplicationId, messageId); isDuplicate {
			return map[string]interface{}{
				"MessageId": existingMsgID,
			}, nil
		}
	}

	msg.PublishedTimestamp = time.Now().UTC()
	msg.ReceivedTimestamp = time.Now().UTC()

	subscriptions, err := store.ListSubscriptionsByTopic(in.TopicArn, common.ListOptions{})
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
				Message:           in.Message,
				Subject:           in.Subject,
				MessageStructure:  in.MessageStructure,
				MessageGroupId:    in.MessageGroupId,
				MessageAttributes: msgAttrs,
			}
			snsEvt.Region = region
			if err := s.bus.Publish(context.Background(), snsEvt); err != nil {
				logs.Warn("Failed to publish SNS delivery event to event bus; message is stored but subscribers may not be notified",
					logs.String("topicArn", in.TopicArn),
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
		result["SequenceNumber"] = store.GetNextSequenceNumber(in.TopicArn, in.MessageGroupId)
	}
	return result, nil
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

// publishBatchCore is the single validation and persistence path for
// PublishBatch. It needs the request context for the delivery region.
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
func (s *SNSService) publishBatchCore(store snsstore.SNSStoreInterface, reqCtx *request.RequestContext, in PublishBatchInput) (interface{}, error) {
	if in.TopicArn == "" {
		return nil, awserrors.NewInvalidParameterException("TopicArn is required")
	}

	topic, err := store.GetTopic(in.TopicArn)
	if err != nil {
		if err == snsstore.ErrTopicNotFound {
			return nil, ErrTopicNotFound
		}
		return nil, err
	}

	entryMaps := in.Entries
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

	subscriptions, err := store.ListSubscriptionsByTopic(in.TopicArn, common.ListOptions{})
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
			if existingMsgID, isDuplicate := store.CheckDeduplication(in.TopicArn, messageDeduplicationId); isDuplicate {
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
			existingMsgID, isDuplicate := store.CheckAndRecordDeduplication(in.TopicArn, entry.messageDeduplicationId, messageId)
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
			TopicArn:               in.TopicArn,
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
					TopicARN:          in.TopicArn,
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
			result["SequenceNumber"] = store.GetNextSequenceNumber(in.TopicArn, entry.messageGroupId)
		}

		successful = append(successful, result)
	}

	return map[string]interface{}{
		"Successful": successful,
		"Failed":     failed,
	}, nil
}

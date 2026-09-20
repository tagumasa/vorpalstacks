package sns

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/common/request"
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

// attributeBytes measures the serialised wire size of a message's
// attributes: each entry contributes its name, DataType and StringValue or
// base64-encoded BinaryValue.
func attributeBytes(attrs map[string]*snsstore.MessageAttribute) int {
	size := 0
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

// messageEntrySize estimates the serialised wire size of a single Publish or
// PublishBatch entry: message body + subject + all message attributes
// (name + DataType + StringValue or BinaryValue).  AWS counts the full
// serialised request toward the 256 KB batch limit, so excluding attributes
// under-estimates and allows oversized batches through.
func messageEntrySize(message, subject string, attrs map[string]*snsstore.MessageAttribute) int {
	return len(message) + len(subject) + attributeBytes(attrs)
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
		// The Query flat-key entries are numbered contiguously from 1; the
		// walk ends at the first gap. There is no parse-side entry cap to
		// drift from the documented limit — the count check below is the
		// single enforcement of MaxMessageAttributes.
		for i := 1; ; i++ {
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
	if len(attrs) > snsstore.MaxMessageAttributes {
		return NewInvalidParameter(fmt.Sprintf("Too many message attributes: %d (maximum %d)", len(attrs), snsstore.MaxMessageAttributes))
	}

	msg.MessageAttributes = make(map[string]*snsstore.MessageAttribute, len(attrs))
	for k, v := range attrs {
		// Validate attribute name format per AWS spec.
		if err := validateMessageAttributeName(k); err != nil {
			return err
		}

		attrMap, ok := v.(map[string]interface{})
		if !ok {
			return NewInvalidParameter(fmt.Sprintf("Invalid message attribute %q: value must be a map", k))
		}

		dataType := firstString(attrMap, "DataType", "dataType")
		if dataType == "" {
			return NewInvalidParameter(fmt.Sprintf("Invalid message attribute %q: DataType is required", k))
		}
		if !validMessageAttributeDataTypes[dataType] {
			return NewInvalidParameter(fmt.Sprintf("Invalid message attribute %q: DataType %q is not valid (String, Number, String.Array, Binary)", k, dataType))
		}

		attr := &snsstore.MessageAttribute{
			Type:        dataType,
			StringValue: firstString(attrMap, "StringValue", "stringValue"),
		}
		if raw := firstString(attrMap, "BinaryValue", "binaryValue"); raw != "" {
			decoded, err := base64.StdEncoding.DecodeString(raw)
			if err != nil {
				return NewInvalidParameter(fmt.Sprintf("Invalid message attribute %q: BinaryValue is not valid base64", k))
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
	return hex.EncodeToString(hash[:])
}

// publishCore is the single validation and persistence path for Publish.
// region is the delivery region the request context carries.
func (s *SNSService) publishCore(store snsstore.SNSStoreInterface, region string, in PublishInput) (interface{}, error) {
	// TargetArn is an AWS-supported alternative to TopicArn. PhoneNumber
	// is silently accepted by AWS but SMS sending is out-of-scope here —
	// reject it explicitly so callers get a clear error instead of silent
	// success.
	if in.PhoneNumber != "" {
		return nil, NewInvalidParameter("PhoneNumber is not supported (SMS sending is not available)")
	}

	if in.TopicArn == "" && in.TargetArn == "" {
		return nil, NewInvalidParameter("TopicArn (or TargetArn) is required")
	}
	if in.TopicArn == "" {
		in.TopicArn = in.TargetArn
	}
	if in.Message == "" {
		return nil, NewInvalidParameter("Message is required")
	}

	topic, err := store.GetTopic(in.TopicArn)
	if err != nil {
		return nil, mapStoreError(err)
	}

	msg := &snsstore.Message{
		MessageId:              uuid.New().String(),
		TopicArn:               topic.Arn,
		Subject:                in.Subject,
		Message:                in.Message,
		MessageStructure:       in.MessageStructure,
		MessageGroupId:         in.MessageGroupId,
		MessageDeduplicationId: in.MessageDeduplicationId,
	}

	// Attributes parse ahead of validation: the size ceiling counts the
	// combined body and attributes (Publish), so the validator needs the
	// parsed set.
	if err := parseMessageAttributes(in.Parameters, msg); err != nil {
		return nil, err
	}

	if err := validatePublishParams(topic.IsFifoTopic(), topic.IsContentBasedDeduplication(), topic.MaximumMessageSizeBytes(), in.Message, in.Subject, in.MessageStructure, in.MessageGroupId, in.MessageDeduplicationId, msg.MessageAttributes); err != nil {
		return nil, err
	}

	if topic.IsFifoTopic() && in.MessageDeduplicationId == "" {
		in.MessageDeduplicationId = generateContentBasedDeduplicationId(in.Message)
		msg.MessageDeduplicationId = in.MessageDeduplicationId
	}

	msg.PublishedTimestamp = time.Now().UTC()

	subscriptions, err := store.ListAllSubscriptionsByTopic(in.TopicArn)
	if err != nil {
		return nil, mapStoreError(err)
	}

	if topic.IsFifoTopic() {
		result, err := s.publishFifoOrdered(store, msg, subscriptions, region, topic.PerGroupDeduplication())
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	if len(subscriptions) > 0 {
		if err := s.dispatchPublish(msg, subscriptions, region); err != nil {
			return nil, fmt.Errorf("dispatch publish: %w", err)
		}
	}

	return map[string]interface{}{
		"MessageId": msg.MessageId,
	}, nil
}

// publishFifoOrdered runs one FIFO publish: deduplication, sequence
// allocation and delivery execute inside the store's per-(topic, message
// group) critical section, so the sequence order the responses report is
// the delivery order every endpoint observes. Delivery is synchronous —
// the async fan-out seam cannot preserve a per-group order across its
// worker pool — and per-subscription failures route to the DLQ exactly as
// the async path routes them. A dedup hit is accepted but not delivered
// and answers with the original publish's MessageId and SequenceNumber
// (the sequence number is assigned to each message; a dedup hit IS the
// original message, which is also why no new number is allocated for it).
// perGroupDedup scopes the deduplication window to the message group when
// the topic's FifoThroughputScope is MessageGroup.
func (s *SNSService) publishFifoOrdered(store snsstore.SNSStoreInterface, msg *snsstore.Message, subscriptions []*snsstore.Subscription, region string, perGroupDedup bool) (map[string]interface{}, error) {
	release := store.AcquireFifoGroupLock(msg.TopicArn, msg.MessageGroupId)
	defer release()

	if messageID, sequenceNumber, hit, err := store.CheckFifoDeduplication(msg.TopicArn, msg.MessageGroupId, msg.MessageDeduplicationId, perGroupDedup); err != nil {
		return nil, mapStoreError(err)
	} else if hit {
		return map[string]interface{}{
			"MessageId":      messageID,
			"SequenceNumber": sequenceNumber,
		}, nil
	}

	sequenceNumber, err := store.AllocateFifoSequence(msg.TopicArn, msg.MessageGroupId)
	if err != nil {
		return nil, mapStoreError(err)
	}

	if messageID, existingSeq, hit, err := store.CheckAndRecordFifoDeduplication(msg.TopicArn, msg.MessageGroupId, msg.MessageDeduplicationId, msg.MessageId, sequenceNumber, perGroupDedup); err != nil {
		return nil, mapStoreError(err)
	} else if hit {
		// A publish recording the same dedup ID between the pre-check and
		// here — the window admits only one copy at the topic's dedup
		// scope, so this publish is accepted but not delivered and answers
		// with that publish's identifiers.
		return map[string]interface{}{
			"MessageId":      messageID,
			"SequenceNumber": existingSeq,
		}, nil
	}

	s.deliverWithRecover(msg, subscriptions, region)

	return map[string]interface{}{
		"MessageId":      msg.MessageId,
		"SequenceNumber": sequenceNumber,
	}, nil
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
	existingSeq            string
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
		return nil, NewInvalidParameter("TopicArn is required")
	}

	topic, err := store.GetTopic(in.TopicArn)
	if err != nil {
		return nil, mapStoreError(err)
	}

	entryMaps := in.Entries
	if len(entryMaps) == 0 {
		return nil, ErrEmptyBatchRequest
	}
	if len(entryMaps) > snsstore.MaxBatchEntries {
		return nil, ErrTooManyEntriesInBatch
	}

	successful := make([]map[string]interface{}, 0)
	failed := make([]map[string]interface{}, 0)
	seenIds := make(map[string]bool, len(entryMaps))
	validated := make([]batchValidatedEntry, 0, len(entryMaps))
	batchTotalSize := 0

	subscriptions, err := store.ListAllSubscriptionsByTopic(in.TopicArn)
	if err != nil {
		return nil, mapStoreError(err)
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

		// Attributes parse ahead of validation on every entry: the size
		// ceiling counts the combined body and attributes (Publish), and
		// the batch total counts every entry's wire size — duplicates
		// included, their bytes crossed the wire like any other's.
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

		if err := validatePublishParams(topic.IsFifoTopic(), topic.IsContentBasedDeduplication(), topic.MaximumMessageSizeBytes(), message, subject, messageStructure, messageGroupId, messageDeduplicationId, msg.MessageAttributes); err != nil {
			failed = append(failed, map[string]interface{}{
				"Id":          id,
				"Code":        "InvalidParameter",
				"Message":     err.Error(),
				"SenderFault": true,
			})
			continue
		}

		batchTotalSize += messageEntrySize(message, subject, msg.MessageAttributes)

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
			existingMsgID, existingSeq, isDuplicate, err := store.CheckFifoDeduplication(in.TopicArn, messageGroupId, messageDeduplicationId, topic.PerGroupDeduplication())
			if err != nil {
				return nil, mapStoreError(err)
			}
			if isDuplicate {
				entry.isDuplicate = true
				entry.existingMsgID = existingMsgID
				entry.existingSeq = existingSeq
			}
		}

		if !entry.isDuplicate {
			entry.msgAttrs = msg.MessageAttributes
		}

		validated = append(validated, entry)
	}

	// Batch-level size check: reject the entire batch before any delivery.
	if batchTotalSize > snsstore.MaxBatchTotalSize {
		return nil, ErrBatchRequestTooLong
	}

	// --- Pass 2: deliver validated entries ---
	for _, entry := range validated {
		if entry.isDuplicate {
			successful = append(successful, map[string]interface{}{
				"Id":             entry.id,
				"MessageId":      entry.existingMsgID,
				"SequenceNumber": entry.existingSeq,
			})
			continue
		}

		messageId := uuid.New().String()

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

		if topic.IsFifoTopic() {
			// The ordered path owns deduplication, sequence allocation and
			// synchronous delivery per entry; its result already carries
			// the batch's Id column.
			result, err := s.publishFifoOrdered(store, msg, subscriptions, region, topic.PerGroupDeduplication())
			if err != nil {
				failed = append(failed, map[string]interface{}{
					"Id":          entry.id,
					"Code":        "InternalError",
					"Message":     err.Error(),
					"SenderFault": false,
				})
				continue
			}
			result["Id"] = entry.id
			successful = append(successful, result)
			continue
		}

		if len(subscriptions) > 0 {
			if err := s.dispatchPublish(msg, subscriptions, region); err != nil {
				// A dispatch failure is surfaced per entry — the batch's
				// Failed list is the batch's own error channel; entries
				// already dispatched keep their Successful rows.
				failed = append(failed, map[string]interface{}{
					"Id":          entry.id,
					"Code":        "InternalError",
					"Message":     err.Error(),
					"SenderFault": false,
				})
				continue
			}
		}

		successful = append(successful, map[string]interface{}{
			"Id":        entry.id,
			"MessageId": messageId,
		})
	}

	return map[string]interface{}{
		"Successful": successful,
		"Failed":     failed,
	}, nil
}

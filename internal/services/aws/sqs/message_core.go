package sqs

import (
	"strconv"
	"strings"

	"vorpalstacks/internal/common/request"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

// ---------------------------------------------------------------------------
// Message-operation Core — single validation + persistence path for the
// message plane (send, receive, delete, visibility). The HTTP handlers in
// message_operations.go are thin adapters: wire parsing, DTO construction,
// the Core call and response serialisation. Parameters carries the raw
// request parameter map so the dual-format (JSON / flattened query) entry
// and attribute parsing keeps a single implementation here.
// ---------------------------------------------------------------------------

// SendMessageInput carries the parameters for sending a single message.
// Parameters holds the raw request map for the message-attribute and
// system-attribute members that arrive in two wire formats.
type SendMessageInput struct {
	QueueURL     string
	MessageBody  string
	DelaySeconds int32
	Parameters   map[string]interface{}
}

// SendMessageBatchInput carries the parameters for a batch send. Parameters
// holds the raw request map; batch entries are parsed and validated inside
// the Core so no entry is sent before every entry has passed validation.
type SendMessageBatchInput struct {
	QueueURL   string
	Parameters map[string]interface{}
}

// ReceiveMessageInput carries the parameters for receiving messages.
// Parameters holds the raw request map for the attribute-name lists and
// receive options that arrive in two wire formats.
type ReceiveMessageInput struct {
	QueueURL   string
	Parameters map[string]interface{}
}

// DeleteMessageInput carries the parameters for deleting a single message.
type DeleteMessageInput struct {
	QueueURL      string
	ReceiptHandle string
}

// DeleteMessageBatchInput carries the parameters for a batch delete.
// Parameters holds the raw request map so entry parsing happens inside the
// Core after the QueueUrl check, preserving the request-level error
// precedence (an empty QueueUrl wins over a malformed entry list).
type DeleteMessageBatchInput struct {
	QueueURL   string
	Parameters map[string]interface{}
}

// ChangeMessageVisibilityInput carries the parameters for changing the
// visibility timeout of a single message.
type ChangeMessageVisibilityInput struct {
	QueueURL          string
	ReceiptHandle     string
	VisibilityTimeout int32
}

// ChangeMessageVisibilityBatchInput carries the parameters for a batch
// visibility change. Parameters holds the raw request map so entry parsing
// happens inside the Core after the QueueUrl check, preserving the
// request-level error precedence.
type ChangeMessageVisibilityBatchInput struct {
	QueueURL   string
	Parameters map[string]interface{}
}

// sendMessageCore sends a message to an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_SendMessage.html
func (s *SQSService) sendMessageCore(store sqsstore.SQSStoreInterface, in SendMessageInput) (interface{}, error) {
	if in.QueueURL == "" {
		return nil, ErrMissingParameter
	}
	if in.MessageBody == "" {
		return nil, ErrMissingParameter
	}

	if err := rejectAttributeListValues(in.Parameters); err != nil {
		return nil, err
	}

	message := sqsstore.NewMessage(in.MessageBody)
	message.DelaySeconds = in.DelaySeconds
	message.MessageGroupID = request.GetParamCaseInsensitive(in.Parameters, "MessageGroupId")
	message.MessageDeduplicationID = request.GetParamCaseInsensitive(in.Parameters, "MessageDeduplicationId")

	messageAttributes := make(map[string]*sqsstore.MessageAttributeValue)

	if jsonAttrs, ok := in.Parameters["MessageAttributes"].(map[string]interface{}); ok && len(jsonAttrs) > 0 {
		for name, val := range jsonAttrs {
			attrMap, ok := val.(map[string]interface{})
			if !ok {
				continue
			}
			attrValue := &sqsstore.MessageAttributeValue{
				DataType: request.GetStringParam(attrMap, "DataType"),
			}
			if sv, ok := attrMap["StringValue"].(string); ok && sv != "" {
				attrValue.StringValue = &sv
			}
			if bv, ok := attrMap["BinaryValue"].(string); ok && bv != "" {
				decoded, dErr := sqsstore.DecodeBinaryValue(bv)
				if dErr != nil {
					return nil, ErrInvalidParameterValue
				}
				attrValue.BinaryValue = decoded
			}
			messageAttributes[name] = attrValue
		}
	} else {
		for i := 1; ; i++ {
			attrName := request.GetParamCaseInsensitive(in.Parameters, "MessageAttribute."+strconv.Itoa(i)+".Name")
			if attrName == "" {
				attrNameKey := "MessageAttribute." + strconv.Itoa(i) + ".Name"
				if val, ok := in.Parameters[attrNameKey].(string); ok {
					attrName = val
				}
			}
			if attrName == "" {
				break
			}

			dataType := request.GetParamCaseInsensitive(in.Parameters, "MessageAttribute."+strconv.Itoa(i)+".Value.DataType")
			if dataType == "" {
				dataTypeKey := "MessageAttribute." + strconv.Itoa(i) + ".Value.DataType"
				if val, ok := in.Parameters[dataTypeKey].(string); ok {
					dataType = val
				}
			}

			attrValue := &sqsstore.MessageAttributeValue{DataType: dataType}

			stringValue := request.GetParamCaseInsensitive(in.Parameters, "MessageAttribute."+strconv.Itoa(i)+".Value.StringValue")
			if stringValue == "" {
				svKey := "MessageAttribute." + strconv.Itoa(i) + ".Value.StringValue"
				if val, ok := in.Parameters[svKey].(string); ok {
					stringValue = val
				}
			}
			if stringValue != "" {
				attrValue.StringValue = &stringValue
			}

			binaryValue := request.GetParamCaseInsensitive(in.Parameters, "MessageAttribute."+strconv.Itoa(i)+".Value.BinaryValue")
			if binaryValue == "" {
				bvKey := "MessageAttribute." + strconv.Itoa(i) + ".Value.BinaryValue"
				if val, ok := in.Parameters[bvKey].(string); ok {
					binaryValue = val
				}
			}
			if binaryValue != "" {
				decoded, dErr := sqsstore.DecodeBinaryValue(binaryValue)
				if dErr != nil {
					return nil, ErrInvalidParameterValue
				}
				attrValue.BinaryValue = decoded
			}

			messageAttributes[attrName] = attrValue
		}
	}
	message.MessageAttributes = messageAttributes

	// Validate the complete message attribute set (count, names, data
	// types) through the single store-layer validator.
	if err := sqsstore.ValidateMessageAttributes(message.MessageAttributes); err != nil {
		return nil, convertStoreError(err)
	}

	// Parse MessageSystemAttributes (only AWSTraceHeader is valid for sends)
	systemAttrs, err := parseMessageSystemAttributes(in.Parameters)
	if err != nil {
		return nil, err
	}
	if traceHeader, ok := systemAttrs["AWSTraceHeader"]; ok && traceHeader.StringValue != nil {
		if message.Attributes == nil {
			message.Attributes = make(map[string]string)
		}
		message.Attributes["AWSTraceHeader"] = *traceHeader.StringValue
	}

	created, err := store.SendMessage(in.QueueURL, message)
	if err != nil {
		return nil, convertStoreError(err)
	}

	response := messageSendResponseFields(created)
	if len(systemAttrs) > 0 {
		response["MD5OfMessageSystemAttributes"] = sqsstore.CalculateMessageAttributesMD5(systemAttrs)
	}
	return response, nil
}

// sendMessageBatchCore sends multiple messages to an SQS queue in a single
// request.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_SendMessageBatch.html
//
// Two-pass design — all entries are parsed and validated before any message
// is sent. This prevents partial sends when a later entry fails validation.
func (s *SQSService) sendMessageBatchCore(store sqsstore.SQSStoreInterface, in SendMessageBatchInput) (interface{}, error) {
	if in.QueueURL == "" {
		return nil, ErrMissingParameter
	}

	if err := rejectAttributeListValues(in.Parameters); err != nil {
		return nil, err
	}

	// Pass 1: Parse and validate all entries (no sends).
	parsed, err := parseBatchSendEntries(in.Parameters)
	if err != nil {
		return nil, err
	}

	// QueueDoesNotExist is a request-level error: the whole request fails
	// before any entry is processed.
	if _, err := store.GetQueue(in.QueueURL); err != nil {
		return nil, convertStoreError(err)
	}

	// Pass 2: Send each validated message.
	successEntries := make([]map[string]interface{}, 0)
	failedEntries := make([]map[string]interface{}, 0)

	for _, e := range parsed {
		created, sendErr := store.SendMessage(in.QueueURL, e.message)
		if sendErr != nil {
			code, senderFault := mapStoreErrorToBatchCode(sendErr)
			failedEntries = append(failedEntries, map[string]interface{}{
				"Id":          e.id,
				"SenderFault": senderFault,
				"Code":        code,
				"Message":     sendErr.Error(),
			})
			continue
		}

		entry := map[string]interface{}{
			"Id": e.id,
		}
		for k, v := range messageSendResponseFields(created) {
			entry[k] = v
		}
		if len(e.systemAttrs) > 0 {
			entry["MD5OfMessageSystemAttributes"] = sqsstore.CalculateMessageAttributesMD5(e.systemAttrs)
		}
		successEntries = append(successEntries, entry)
	}

	return map[string]interface{}{
		"Successful": successEntries,
		"Failed":     failedEntries,
	}, nil
}

// receiveMessageCore receives one or more messages from an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_ReceiveMessage.html
func (s *SQSService) receiveMessageCore(store sqsstore.SQSStoreInterface, in ReceiveMessageInput) (interface{}, error) {
	if in.QueueURL == "" {
		return nil, ErrMissingParameter
	}

	// MaxNumberOfMessages carries the AWS contract "Valid values: 1 to 10.
	// Default: 1." — an explicitly provided value outside the range
	// (including 0) is rejected, a present value that is not an integer is
	// a wire-type violation, and only omission falls back to the default.
	maxNumberOfMessages := int32(1)
	if val, present, perr := request.GetIntParamStrictCaseInsensitive(in.Parameters, "MaxNumberOfMessages"); present {
		if perr != nil {
			return nil, ErrSerializationException
		}
		maxNumberOfMessages = int32(val)
		if err := validateMaxNumberOfMessages(maxNumberOfMessages); err != nil {
			return nil, err
		}
	}

	var visibilityTimeoutPtr *int32
	if vtVal, present, perr := request.GetIntParamStrictCaseInsensitive(in.Parameters, "VisibilityTimeout"); present {
		if perr != nil {
			return nil, ErrSerializationException
		}
		vt := int32(vtVal)
		visibilityTimeoutPtr = &vt
	}
	// WaitTimeSeconds: a negative sentinel tells the store the request
	// omitted the parameter so it applies the queue's
	// ReceiveMessageWaitTimeSeconds attribute; explicit values are validated
	// here and long-poll in the store.
	waitTimeSeconds := int32(-1)
	if val, present, perr := request.GetIntParamStrictCaseInsensitive(in.Parameters, "WaitTimeSeconds"); present {
		if perr != nil {
			return nil, ErrSerializationException
		}
		if val < 0 {
			return nil, ErrInvalidParameterValue
		}
		waitTimeSeconds = int32(val)
	}

	// Parse system attribute names: support both legacy AttributeNames
	// and newer MessageSystemAttributeNames. Specifying both is an error.
	legacyAttrNames := request.GetStringList(in.Parameters, "AttributeNames")
	newSysAttrNames := request.GetStringList(in.Parameters, "MessageSystemAttributeNames")

	if len(legacyAttrNames) > 0 && len(newSysAttrNames) > 0 {
		return nil, ErrInvalidParameterCombination
	}

	var sysAttrNames []string
	if len(legacyAttrNames) > 0 {
		sysAttrNames = legacyAttrNames
	} else {
		sysAttrNames = newSysAttrNames
	}

	// Parse message attribute names
	msgAttrNames := request.GetStringList(in.Parameters, "MessageAttributeNames")

	// Parse ReceiveRequestAttemptId for FIFO receive dedup
	receiveRequestAttemptId := request.GetParamCaseInsensitive(in.Parameters, "ReceiveRequestAttemptId")
	if err := validateReceiveRequestAttemptId(receiveRequestAttemptId); err != nil {
		return nil, err
	}

	messages, err := store.ReceiveMessage(in.QueueURL, maxNumberOfMessages, visibilityTimeoutPtr, waitTimeSeconds, receiveRequestAttemptId)
	if err != nil {
		return nil, convertStoreError(err)
	}

	sysAttrAll := shouldReturnAllAttributes(sysAttrNames)
	msgAttrAll := shouldReturnAllAttributes(msgAttrNames)

	messageList := make([]map[string]interface{}, 0, len(messages))
	for _, msg := range messages {
		msgMap := map[string]interface{}{
			"MessageId":     msg.ID,
			"ReceiptHandle": msg.ReceiptHandle,
			"MD5OfBody":     msg.MD5OfBody,
			"Body":          msg.Body,
		}
		if len(msg.MessageAttributes) > 0 {
			msgMap["MD5OfMessageAttributes"] = msg.MD5OfMessageAttributes
		}

		if msg.MessageGroupID != "" {
			msgMap["MessageGroupId"] = msg.MessageGroupID
		}
		if msg.MessageDeduplicationID != "" {
			msgMap["MessageDeduplicationId"] = msg.MessageDeduplicationID
		}

		// Filter message attributes
		if len(msg.MessageAttributes) > 0 {
			if msgAttrAll {
				attrs := make(map[string]interface{}, len(msg.MessageAttributes))
				for k, v := range msg.MessageAttributes {
					attrMap := map[string]interface{}{
						"DataType": v.DataType,
					}
					if v.StringValue != nil {
						attrMap["StringValue"] = *v.StringValue
					}
					if v.BinaryValue != nil {
						attrMap["BinaryValue"] = sqsstore.EncodeBinaryValue(v.BinaryValue)
					}
					attrs[k] = attrMap
				}
				msgMap["MessageAttributes"] = attrs
			} else {
				attrs := make(map[string]interface{})
				for k, v := range msg.MessageAttributes {
					if !isRequestedAttribute(k, msgAttrNames) {
						continue
					}
					attrMap := map[string]interface{}{
						"DataType": v.DataType,
					}
					if v.StringValue != nil {
						attrMap["StringValue"] = *v.StringValue
					}
					if v.BinaryValue != nil {
						attrMap["BinaryValue"] = sqsstore.EncodeBinaryValue(v.BinaryValue)
					}
					attrs[k] = attrMap
				}
				if len(attrs) > 0 {
					msgMap["MessageAttributes"] = attrs
				}
			}
		}

		// Filter system attributes
		if len(msg.Attributes) > 0 {
			if sysAttrAll {
				msgMap["Attributes"] = msg.Attributes
			} else {
				filtered := make(map[string]string)
				for k, v := range msg.Attributes {
					if isRequestedAttribute(k, sysAttrNames) {
						filtered[k] = v
					}
				}
				if len(filtered) > 0 {
					msgMap["Attributes"] = filtered
				}
			}
		}

		messageList = append(messageList, msgMap)
	}

	if len(messageList) == 0 {
		return map[string]interface{}{}, nil
	}

	return map[string]interface{}{
		"Messages": messageList,
	}, nil
}

// deleteMessageCore deletes a message from an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_DeleteMessage.html
func (s *SQSService) deleteMessageCore(store sqsstore.SQSStoreInterface, in DeleteMessageInput) error {
	if in.QueueURL == "" {
		return ErrMissingParameter
	}
	if in.ReceiptHandle == "" {
		return ErrMissingParameter
	}

	if err := store.DeleteMessage(in.QueueURL, in.ReceiptHandle); err != nil {
		return convertStoreError(err)
	}
	return nil
}

// deleteMessageBatchCore deletes multiple messages from an SQS queue in a
// single request.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_DeleteMessageBatch.html
//
// Two-pass design — all entries are validated before any deletion. This
// prevents partial state when a later entry fails validation.
func (s *SQSService) deleteMessageBatchCore(store sqsstore.SQSStoreInterface, in DeleteMessageBatchInput) (interface{}, error) {
	if in.QueueURL == "" {
		return nil, ErrMissingParameter
	}

	// Pass 1: Parse and validate all entries (no mutations).
	entries, err := parseDeleteBatchEntries(in.Parameters)
	if err != nil {
		return nil, err
	}

	seenIDs := make(map[string]bool, len(entries))
	for _, e := range entries {
		if err := sqsstore.ValidateBatchEntryId(e.id); err != nil {
			return nil, ErrInvalidBatchEntryId
		}
		if seenIDs[e.id] {
			return nil, ErrBatchEntryIdsNotDistinct
		}
		seenIDs[e.id] = true
	}
	if len(entries) > sqsstore.MaxBatchEntries {
		return nil, ErrTooManyEntriesInBatch
	}

	// QueueDoesNotExist is a request-level error: the whole request fails
	// before any entry is processed.
	if _, err := store.GetQueue(in.QueueURL); err != nil {
		return nil, convertStoreError(err)
	}

	// Pass 2: Execute deletions.
	successEntries := make([]map[string]interface{}, 0)
	failedEntries := make([]map[string]interface{}, 0)

	for _, e := range entries {
		if err := store.DeleteMessage(in.QueueURL, e.receiptHandle); err != nil {
			code, senderFault := mapStoreErrorToBatchCode(err)
			failedEntries = append(failedEntries, map[string]interface{}{
				"Id":          e.id,
				"SenderFault": senderFault,
				"Code":        code,
				"Message":     err.Error(),
			})
			continue
		}
		successEntries = append(successEntries, map[string]interface{}{
			"Id": e.id,
		})
	}

	return map[string]interface{}{
		"Successful": successEntries,
		"Failed":     failedEntries,
	}, nil
}

// changeMessageVisibilityCore changes the visibility timeout of a message.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_ChangeMessageVisibility.html
func (s *SQSService) changeMessageVisibilityCore(store sqsstore.SQSStoreInterface, in ChangeMessageVisibilityInput) error {
	if in.QueueURL == "" {
		return ErrMissingParameter
	}
	if in.ReceiptHandle == "" {
		return ErrMissingParameter
	}

	if err := store.ChangeMessageVisibility(in.QueueURL, in.ReceiptHandle, in.VisibilityTimeout); err != nil {
		return convertStoreError(err)
	}
	return nil
}

// changeMessageVisibilityBatchCore changes the visibility timeout for
// multiple messages in a single request.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_ChangeMessageVisibilityBatch.html
//
// Two-pass design — all entries are validated before any mutation.
func (s *SQSService) changeMessageVisibilityBatchCore(store sqsstore.SQSStoreInterface, in ChangeMessageVisibilityBatchInput) (interface{}, error) {
	if in.QueueURL == "" {
		return nil, ErrMissingParameter
	}

	// Pass 1: Parse and validate all entries (no mutations).
	entries, err := parseChangeVisibilityBatchEntries(in.Parameters)
	if err != nil {
		return nil, err
	}

	seenIDs := make(map[string]bool, len(entries))
	for _, e := range entries {
		if err := sqsstore.ValidateBatchEntryId(e.id); err != nil {
			return nil, ErrInvalidBatchEntryId
		}
		if seenIDs[e.id] {
			return nil, ErrBatchEntryIdsNotDistinct
		}
		seenIDs[e.id] = true
	}
	if len(entries) > sqsstore.MaxBatchEntries {
		return nil, ErrTooManyEntriesInBatch
	}

	// QueueDoesNotExist is a request-level error: the whole request fails
	// before any entry is processed.
	if _, err := store.GetQueue(in.QueueURL); err != nil {
		return nil, convertStoreError(err)
	}

	// Pass 2: Execute visibility changes.
	successEntries := make([]map[string]interface{}, 0)
	failedEntries := make([]map[string]interface{}, 0)

	for _, e := range entries {
		if err := store.ChangeMessageVisibility(in.QueueURL, e.receiptHandle, e.visibilityTimeout); err != nil {
			code, senderFault := mapStoreErrorToBatchCode(err)
			failedEntries = append(failedEntries, map[string]interface{}{
				"Id":          e.id,
				"SenderFault": senderFault,
				"Code":        code,
				"Message":     err.Error(),
			})
			continue
		}
		successEntries = append(successEntries, map[string]interface{}{
			"Id": e.id,
		})
	}

	return map[string]interface{}{
		"Successful": successEntries,
		"Failed":     failedEntries,
	}, nil
}

// rejectAttributeListValues rejects message-sending requests that carry
// message-attribute list values. The service model marks StringListValues
// and BinaryListValues as not implemented ("Not implemented. Reserved for
// future use."), and AWS rejects any request that includes them with
// UnsupportedOperation, so they must not be silently accepted and dropped.
func rejectAttributeListValues(params map[string]interface{}) error {
	// JSON protocol: SendMessage carries the attribute maps at the top level
	// and SendMessageBatch nests them inside each batch entry.
	attrMaps := make([]map[string]interface{}, 0, 8)
	for _, key := range []string{"MessageAttributes", "MessageSystemAttributes"} {
		if attrs, ok := params[key].(map[string]interface{}); ok {
			for _, val := range attrs {
				if attrMap, ok := val.(map[string]interface{}); ok {
					attrMaps = append(attrMaps, attrMap)
				}
			}
		}
	}
	if entries, ok := params["Entries"].([]interface{}); ok {
		for _, raw := range entries {
			entryMap, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}
			for _, key := range []string{"MessageAttributes", "MessageSystemAttributes"} {
				if attrs, ok := entryMap[key].(map[string]interface{}); ok {
					for _, val := range attrs {
						if attrMap, ok := val.(map[string]interface{}); ok {
							attrMaps = append(attrMaps, attrMap)
						}
					}
				}
			}
		}
	}
	for _, attrMap := range attrMaps {
		if list, ok := attrMap["StringListValues"].([]interface{}); ok && len(list) > 0 {
			return ErrUnsupportedOperation
		}
		if list, ok := attrMap["BinaryListValues"].([]interface{}); ok && len(list) > 0 {
			return ErrUnsupportedOperation
		}
	}
	// Query protocol: the flattened members arrive as indexed flat keys on
	// every message-sending operation (MessageAttribute.1.Value.StringListValue.1,
	// MessageSystemAttribute.1.Value.StringListValue.1, or the batch-entry
	// prefixed variants).
	for key := range params {
		lower := strings.ToLower(key)
		if strings.Contains(lower, ".value.stringlistvalue.") || strings.Contains(lower, ".value.binarylistvalue.") {
			return ErrUnsupportedOperation
		}
	}
	return nil
}

// messageEntrySize calculates the total byte size of a batch entry for the
// BatchRequestTooLong check. Includes message body and all attribute data.
// Delegates to the store helper so the size definition stays singular.
func messageEntrySize(body string, attrs map[string]*sqsstore.MessageAttributeValue) int {
	return sqsstore.MessageSize(body, attrs)
}

// batchSendEntry holds a parsed and validated SendMessageBatch entry awaiting
// dispatch. All validation is completed before any SendMessage call.
type batchSendEntry struct {
	id          string
	message     *sqsstore.Message
	systemAttrs map[string]*sqsstore.MessageAttributeValue
}

// parseBatchSendEntries extracts and validates all SendMessageBatch entries
// from request parameters. Supports both JSON (SDK v2) and flattened query
// (CLI / SDK v1) formats. Returns an error if any entry is malformed or
// fails validation — no messages are sent until this returns successfully.
func parseBatchSendEntries(params map[string]interface{}) ([]*batchSendEntry, error) {
	if jsonEntries, ok := params["Entries"].([]interface{}); ok && len(jsonEntries) > 0 {
		return parseBatchEntriesJSON(jsonEntries)
	}
	return parseBatchEntriesQuery(params)
}

func parseBatchEntriesJSON(jsonEntries []interface{}) ([]*batchSendEntry, error) {
	if len(jsonEntries) > sqsstore.MaxBatchEntries {
		return nil, ErrTooManyEntriesInBatch
	}

	seenIDs := make(map[string]bool)
	batchTotalSize := 0
	result := make([]*batchSendEntry, 0, len(jsonEntries))

	for _, raw := range jsonEntries {
		entryMap, ok := raw.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameterValue
		}

		id, _ := entryMap["Id"].(string)
		if id == "" {
			return nil, ErrInvalidParameterValue
		}

		if err := sqsstore.ValidateBatchEntryId(id); err != nil {
			return nil, ErrInvalidBatchEntryId
		}

		if seenIDs[id] {
			return nil, ErrBatchEntryIdsNotDistinct
		}
		seenIDs[id] = true

		messageBody, _ := entryMap["MessageBody"].(string)
		delaySeconds := int32(0)
		if ds, ok := entryMap["DelaySeconds"].(float64); ok {
			delaySeconds = int32(ds)
		}

		message := sqsstore.NewMessage(messageBody)
		message.DelaySeconds = delaySeconds
		if mgid, ok := entryMap["MessageGroupId"].(string); ok {
			message.MessageGroupID = mgid
		}
		if mdid, ok := entryMap["MessageDeduplicationId"].(string); ok {
			message.MessageDeduplicationID = mdid
		}

		if attrs, ok := entryMap["MessageAttributes"].(map[string]interface{}); ok {
			msgAttrs := make(map[string]*sqsstore.MessageAttributeValue)
			for attrName, attrVal := range attrs {
				if attrMap, ok := attrVal.(map[string]interface{}); ok {
					attr := &sqsstore.MessageAttributeValue{}
					if dt, ok := attrMap["DataType"].(string); ok {
						attr.DataType = dt
					}
					if sv, ok := attrMap["StringValue"].(string); ok {
						attr.StringValue = &sv
					}
					if bv, ok := attrMap["BinaryValue"].(string); ok {
						decoded, dErr := sqsstore.DecodeBinaryValue(bv)
						if dErr != nil {
							return nil, ErrInvalidParameterValue
						}
						attr.BinaryValue = decoded
					}
					msgAttrs[attrName] = attr
				}
			}
			message.MessageAttributes = msgAttrs
		}

		if err := sqsstore.ValidateMessageAttributes(message.MessageAttributes); err != nil {
			return nil, convertStoreError(err)
		}

		batchTotalSize += messageEntrySize(messageBody, message.MessageAttributes)
		if batchTotalSize > sqsstore.MaxMaximumMessageSize {
			return nil, ErrBatchRequestTooLong
		}

		sysAttrs, err := parseBatchEntrySystemAttributesJSON(entryMap)
		if err != nil {
			return nil, err
		}
		if len(sysAttrs) > 0 {
			if message.Attributes == nil {
				message.Attributes = make(map[string]string)
			}
			if th, ok := sysAttrs["AWSTraceHeader"]; ok && th.StringValue != nil {
				message.Attributes["AWSTraceHeader"] = *th.StringValue
			}
		}

		result = append(result, &batchSendEntry{
			id:          id,
			message:     message,
			systemAttrs: sysAttrs,
		})
	}

	if len(result) == 0 {
		return nil, ErrEmptyBatchRequest
	}
	return result, nil
}

func parseBatchEntriesQuery(params map[string]interface{}) ([]*batchSendEntry, error) {
	seenIDs := make(map[string]bool)
	entryCount := 0
	batchTotalSize := 0
	result := make([]*batchSendEntry, 0)

	for i := 1; ; i++ {
		prefix := "SendMessageBatchRequestEntry." + strconv.Itoa(i) + "."
		id := request.GetParamCaseInsensitive(params, prefix+"Id")
		if id == "" {
			// Query-protocol entries are contiguous; the loop is bounded by
			// the entry-count check below, which fails with
			// TooManyEntriesInBatchRequest instead of silently dropping
			// entries beyond a fixed limit.
			break
		}

		if err := sqsstore.ValidateBatchEntryId(id); err != nil {
			return nil, ErrInvalidBatchEntryId
		}

		if seenIDs[id] {
			return nil, ErrBatchEntryIdsNotDistinct
		}
		seenIDs[id] = true
		entryCount++

		if entryCount > sqsstore.MaxBatchEntries {
			return nil, ErrTooManyEntriesInBatch
		}

		messageBody := request.GetParamCaseInsensitive(params, prefix+"MessageBody")
		// DelaySeconds is an Integer member: a present value that is not an
		// integer is a wire-type violation, not an omitted member.
		delayVal, present, derr := request.GetIntParamStrictCaseInsensitive(params, prefix+"DelaySeconds")
		if derr != nil {
			return nil, ErrSerializationException
		}
		var delaySeconds int32
		if present {
			delaySeconds = int32(delayVal)
		}

		message := sqsstore.NewMessage(messageBody)
		message.DelaySeconds = delaySeconds
		message.MessageGroupID = request.GetParamCaseInsensitive(params, prefix+"MessageGroupId")
		message.MessageDeduplicationID = request.GetParamCaseInsensitive(params, prefix+"MessageDeduplicationId")

		msgAttrs := make(map[string]*sqsstore.MessageAttributeValue)
		for j := 1; ; j++ {
			attrPrefix := prefix + "MessageAttribute." + strconv.Itoa(j) + "."
			attrName := request.GetParamCaseInsensitive(params, attrPrefix+"Name")
			if attrName == "" {
				break
			}
			dataType := request.GetParamCaseInsensitive(params, attrPrefix+"Value.DataType")
			if dataType == "" {
				break
			}
			attr := &sqsstore.MessageAttributeValue{DataType: dataType}
			if sv := request.GetParamCaseInsensitive(params, attrPrefix+"Value.StringValue"); sv != "" {
				attr.StringValue = &sv
			}
			if bv := request.GetParamCaseInsensitive(params, attrPrefix+"Value.BinaryValue"); bv != "" {
				decoded, dErr := sqsstore.DecodeBinaryValue(bv)
				if dErr != nil {
					return nil, ErrInvalidParameterValue
				}
				attr.BinaryValue = decoded
			}
			msgAttrs[attrName] = attr
		}
		if len(msgAttrs) > 0 {
			message.MessageAttributes = msgAttrs
		}

		if err := sqsstore.ValidateMessageAttributes(msgAttrs); err != nil {
			return nil, convertStoreError(err)
		}

		batchTotalSize += messageEntrySize(messageBody, msgAttrs)
		if batchTotalSize > sqsstore.MaxMaximumMessageSize {
			return nil, ErrBatchRequestTooLong
		}

		sysAttrs, err := parseBatchEntrySystemAttributesQuery(params, i)
		if err != nil {
			return nil, err
		}
		if len(sysAttrs) > 0 {
			if message.Attributes == nil {
				message.Attributes = make(map[string]string)
			}
			if th, ok := sysAttrs["AWSTraceHeader"]; ok && th.StringValue != nil {
				message.Attributes["AWSTraceHeader"] = *th.StringValue
			}
		}

		result = append(result, &batchSendEntry{
			id:          id,
			message:     message,
			systemAttrs: sysAttrs,
		})
	}

	if len(result) == 0 {
		return nil, ErrEmptyBatchRequest
	}
	return result, nil
}

// parseBatchEntrySystemAttributesJSON extracts system attributes from a JSON
// batch entry map. Only AWSTraceHeader is valid for sends.
func parseBatchEntrySystemAttributesJSON(entryMap map[string]interface{}) (map[string]*sqsstore.MessageAttributeValue, error) {
	result := make(map[string]*sqsstore.MessageAttributeValue)
	sysAttrs, ok := entryMap["MessageSystemAttributes"].(map[string]interface{})
	if !ok {
		return result, nil
	}
	for name, val := range sysAttrs {
		if name != "AWSTraceHeader" {
			return nil, ErrInvalidParameterValue
		}
		attrMap, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		attr := &sqsstore.MessageAttributeValue{
			DataType: request.GetStringParam(attrMap, "DataType"),
		}
		if sv, ok := attrMap["StringValue"].(string); ok && sv != "" {
			attr.StringValue = &sv
		}
		if bv, ok := attrMap["BinaryValue"].(string); ok && bv != "" {
			decoded, dErr := sqsstore.DecodeBinaryValue(bv)
			if dErr != nil {
				return nil, ErrInvalidParameterValue
			}
			attr.BinaryValue = decoded
		}
		result[name] = attr
	}
	return result, nil
}

// parseBatchEntrySystemAttributesQuery extracts system attributes from
// flattened query parameters for batch entry at position entryIndex.
func parseBatchEntrySystemAttributesQuery(params map[string]interface{}, entryIndex int) (map[string]*sqsstore.MessageAttributeValue, error) {
	result := make(map[string]*sqsstore.MessageAttributeValue)
	for j := 1; ; j++ {
		prefix := "SendMessageBatchRequestEntry." + strconv.Itoa(entryIndex) + ".MessageSystemAttribute." + strconv.Itoa(j) + "."
		name := request.GetParamCaseInsensitive(params, prefix+"Name")
		if name == "" {
			break
		}
		if name != "AWSTraceHeader" {
			return nil, ErrInvalidParameterValue
		}
		dataType := request.GetParamCaseInsensitive(params, prefix+"Value.DataType")
		if dataType == "" {
			break
		}
		attr := &sqsstore.MessageAttributeValue{DataType: dataType}
		if sv := request.GetParamCaseInsensitive(params, prefix+"Value.StringValue"); sv != "" {
			attr.StringValue = &sv
		}
		result[name] = attr
	}
	return result, nil
}

// parseMessageSystemAttributes extracts MessageSystemAttributes from the
// request parameters. Supports both JSON map and flattened query formats.
// Only AWSTraceHeader is valid for sends. Unknown system attribute names are
// rejected with InvalidParameterValue.
func parseMessageSystemAttributes(params map[string]interface{}) (map[string]*sqsstore.MessageAttributeValue, error) {
	result := make(map[string]*sqsstore.MessageAttributeValue)
	if jsonAttrs, ok := params["MessageSystemAttributes"].(map[string]interface{}); ok {
		for name, val := range jsonAttrs {
			if name != "AWSTraceHeader" {
				return nil, ErrInvalidParameterValue
			}
			attrMap, ok := val.(map[string]interface{})
			if !ok {
				continue
			}
			attrValue := &sqsstore.MessageAttributeValue{
				DataType: request.GetStringParam(attrMap, "DataType"),
			}
			if sv, ok := attrMap["StringValue"].(string); ok && sv != "" {
				attrValue.StringValue = &sv
			}
			if bv, ok := attrMap["BinaryValue"].(string); ok && bv != "" {
				decoded, dErr := sqsstore.DecodeBinaryValue(bv)
				if dErr != nil {
					return nil, ErrInvalidParameterValue
				}
				attrValue.BinaryValue = decoded
			}
			result[name] = attrValue
		}
		return result, nil
	}

	for i := 1; ; i++ {
		name := request.GetParamCaseInsensitive(params, "MessageSystemAttribute."+strconv.Itoa(i)+".Name")
		if name == "" {
			nameKey := "MessageSystemAttribute." + strconv.Itoa(i) + ".Name"
			if val, ok := params[nameKey].(string); ok {
				name = val
			}
		}
		if name == "" {
			break
		}
		if name != "AWSTraceHeader" {
			return nil, ErrInvalidParameterValue
		}
		dataType := request.GetParamCaseInsensitive(params, "MessageSystemAttribute."+strconv.Itoa(i)+".Value.DataType")
		if dataType == "" {
			dataTypeKey := "MessageSystemAttribute." + strconv.Itoa(i) + ".Value.DataType"
			if val, ok := params[dataTypeKey].(string); ok {
				dataType = val
			}
		}
		attrValue := &sqsstore.MessageAttributeValue{DataType: dataType}
		if sv := request.GetParamCaseInsensitive(params, "MessageSystemAttribute."+strconv.Itoa(i)+".Value.StringValue"); sv != "" {
			attrValue.StringValue = &sv
		}
		result[name] = attrValue
	}

	return result, nil
}

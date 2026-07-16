package sqs

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

func convertStoreError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sqsstore.ErrQueueNotFound) {
		return ErrQueueDoesNotExist
	}
	if errors.Is(err, sqsstore.ErrInvalidReceiptHandle) {
		return ErrReceiptHandleIsInvalid
	}
	if errors.Is(err, sqsstore.ErrMessageTooLarge) {
		return ErrMessageTooLarge
	}
	if errors.Is(err, sqsstore.ErrMissingMessageGroupId) {
		return ErrMissingMessageGroupId
	}
	if errors.Is(err, sqsstore.ErrMissingDeduplicationId) {
		return ErrMissingDeduplicationId
	}
	if errors.Is(err, sqsstore.ErrInvalidParameterValue) {
		return ErrInvalidParameterValue
	}
	if errors.Is(err, sqsstore.ErrPurgeQueueInProgress) {
		return ErrPurgeQueueInProgress
	}
	if errors.Is(err, sqsstore.ErrTooManyTags) {
		return ErrTooManyTags
	}
	if errors.Is(err, sqsstore.ErrInvalidTagKey) {
		return ErrInvalidTagKey
	}
	if errors.Is(err, sqsstore.ErrInvalidTagValue) {
		return ErrInvalidTagValue
	}
	if errors.Is(err, sqsstore.ErrInvalidQueueName) {
		return ErrInvalidQueueName
	}
	return err
}

func mapStoreErrorToBatchCode(err error) (code string, senderFault bool) {
	switch {
	case errors.Is(err, sqsstore.ErrMissingMessageGroupId):
		return "MissingMessageGroupId", true
	case errors.Is(err, sqsstore.ErrMissingDeduplicationId):
		return "MissingDeduplicationId", true
	case errors.Is(err, sqsstore.ErrMessageTooLarge):
		return "MessageTooLarge", true
	case errors.Is(err, sqsstore.ErrInvalidReceiptHandle):
		return "ReceiptHandleIsInvalid", true
	case errors.Is(err, sqsstore.ErrInvalidParameterValue):
		return "InvalidParameterValue", true
	case errors.Is(err, sqsstore.ErrQueueNotFound):
		return "QueueDoesNotExist", true
	default:
		return "InternalError", false
	}
}

// shouldReturnAllAttributes returns true if the attribute name list is empty
// (default: return all) or contains a recognised wildcard token ("All" or ".*").
// Reference: AWS SQS ReceiveMessage API docs.
func shouldReturnAllAttributes(names []string) bool {
	if len(names) == 0 {
		return true
	}
	for _, n := range names {
		if n == "All" || n == ".*" {
			return true
		}
	}
	return false
}

// isRequestedAttribute checks if attrName matches any of the requested patterns.
// Supports exact match and ".*" prefix wildcard (e.g., "Prefix.*").
func isRequestedAttribute(attrName string, requested []string) bool {
	for _, r := range requested {
		if r == attrName {
			return true
		}
		if strings.HasSuffix(r, ".*") {
			prefix := r[:len(r)-2]
			if strings.HasPrefix(attrName, prefix) {
				return true
			}
		}
	}
	return false
}

// SendMessage sends a message to an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_SendMessage.html
func (s *SQSService) SendMessage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueURL := request.GetParamCaseInsensitive(req.Parameters, "QueueUrl")
	if queueURL == "" {
		return nil, ErrMissingParameter
	}

	messageBody := request.GetParamCaseInsensitive(req.Parameters, "MessageBody")
	if messageBody == "" {
		return nil, ErrMissingParameter
	}

	delaySeconds := int32(request.GetIntParam(req.Parameters, "DelaySeconds"))

	message := sqsstore.NewMessage(messageBody)
	message.DelaySeconds = delaySeconds
	message.MessageGroupID = request.GetParamCaseInsensitive(req.Parameters, "MessageGroupId")
	message.MessageDeduplicationID = request.GetParamCaseInsensitive(req.Parameters, "MessageDeduplicationId")

	messageAttributes := make(map[string]*sqsstore.MessageAttributeValue)

	if jsonAttrs, ok := req.Parameters["MessageAttributes"].(map[string]interface{}); ok && len(jsonAttrs) > 0 {
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
				attrValue.BinaryValue = sqsstore.DecodeBinaryValue(bv)
			}
			messageAttributes[name] = attrValue
		}
	} else {
		for i := 1; ; i++ {
			attrName := request.GetParamCaseInsensitive(req.Parameters, "MessageAttribute."+strconv.Itoa(i)+".Name")
			if attrName == "" {
				attrNameKey := "MessageAttribute." + strconv.Itoa(i) + ".Name"
				if val, ok := req.Parameters[attrNameKey].(string); ok {
					attrName = val
				}
			}
			if attrName == "" {
				break
			}

			dataType := request.GetParamCaseInsensitive(req.Parameters, "MessageAttribute."+strconv.Itoa(i)+".Value.DataType")
			if dataType == "" {
				dataTypeKey := "MessageAttribute." + strconv.Itoa(i) + ".Value.DataType"
				if val, ok := req.Parameters[dataTypeKey].(string); ok {
					dataType = val
				}
			}

			attrValue := &sqsstore.MessageAttributeValue{DataType: dataType}

			stringValue := request.GetParamCaseInsensitive(req.Parameters, "MessageAttribute."+strconv.Itoa(i)+".Value.StringValue")
			if stringValue == "" {
				svKey := "MessageAttribute." + strconv.Itoa(i) + ".Value.StringValue"
				if val, ok := req.Parameters[svKey].(string); ok {
					stringValue = val
				}
			}
			if stringValue != "" {
				attrValue.StringValue = &stringValue
			}

			binaryValue := request.GetParamCaseInsensitive(req.Parameters, "MessageAttribute."+strconv.Itoa(i)+".Value.BinaryValue")
			if binaryValue == "" {
				bvKey := "MessageAttribute." + strconv.Itoa(i) + ".Value.BinaryValue"
				if val, ok := req.Parameters[bvKey].(string); ok {
					binaryValue = val
				}
			}
			if binaryValue != "" {
				attrValue.BinaryValue = sqsstore.DecodeBinaryValue(binaryValue)
			}

			messageAttributes[attrName] = attrValue
		}
	}
	message.MessageAttributes = messageAttributes

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	created, err := store.SendMessage(queueURL, message)
	if err != nil {
		return nil, convertStoreError(err)
	}

	response := map[string]interface{}{
		"MessageId":              created.ID,
		"MD5OfMessageBody":       created.MD5OfBody,
		"MD5OfMessageAttributes": created.MD5OfMessageAttributes,
	}
	if created.SequenceNumber != "" {
		response["SequenceNumber"] = created.SequenceNumber
	}
	return response, nil
}

// SendMessageBatch sends multiple messages to an SQS queue in a single request.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_SendMessageBatch.html
func (s *SQSService) SendMessageBatch(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueURL := request.GetParamCaseInsensitive(req.Parameters, "QueueUrl")
	if queueURL == "" {
		if val, ok := req.Parameters["QueueUrl"].(string); ok {
			queueURL = val
		}
	}
	if queueURL == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	entries := make([]map[string]interface{}, 0)
	seenIDs := make(map[string]bool)
	entryCount := 0

	if jsonEntries, ok := req.Parameters["Entries"].([]interface{}); ok && len(jsonEntries) > 0 {
		if len(jsonEntries) > 10 {
			return nil, ErrTooManyEntriesInBatch
		}
		for _, entry := range jsonEntries {
			entryMap, ok := entry.(map[string]interface{})
			if !ok {
				continue
			}

			id, _ := entryMap["Id"].(string)
			if id == "" {
				continue
			}

			if err := sqsstore.ValidateBatchEntryId(id); err != nil {
				return nil, ErrInvalidBatchEntryId
			}

			if seenIDs[id] {
				return nil, ErrBatchEntryIdsNotDistinct
			}
			seenIDs[id] = true
			entryCount++

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
				messageAttributes := make(map[string]*sqsstore.MessageAttributeValue)
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
							attr.BinaryValue = sqsstore.DecodeBinaryValue(bv)
						}
						messageAttributes[attrName] = attr
					}
				}
				message.MessageAttributes = messageAttributes
			}

			created, err := store.SendMessage(queueURL, message)
			if err != nil {
				code, senderFault := mapStoreErrorToBatchCode(err)
				entries = append(entries, map[string]interface{}{
					"Id":          id,
					"SenderFault": senderFault,
					"Code":        code,
					"Message":     err.Error(),
				})
				continue
			}

			batchEntry := map[string]interface{}{
				"Id":                     id,
				"MessageId":              created.ID,
				"MD5OfMessageBody":       created.MD5OfBody,
				"MD5OfMessageAttributes": created.MD5OfMessageAttributes,
			}
			if created.SequenceNumber != "" {
				batchEntry["SequenceNumber"] = created.SequenceNumber
			}
			entries = append(entries, batchEntry)
		}
	} else {
		for i := 1; i <= 10; i++ {
			id := request.GetParamCaseInsensitive(req.Parameters, "SendMessageBatchRequestEntry."+strconv.Itoa(i)+".Id")
			if id == "" {
				continue
			}

			if err := sqsstore.ValidateBatchEntryId(id); err != nil {
				return nil, ErrInvalidBatchEntryId
			}

			if seenIDs[id] {
				return nil, ErrBatchEntryIdsNotDistinct
			}
			seenIDs[id] = true
			entryCount++

			if entryCount > 10 {
				return nil, ErrTooManyEntriesInBatch
			}

			messageBody := request.GetParamCaseInsensitive(req.Parameters, "SendMessageBatchRequestEntry."+strconv.Itoa(i)+".MessageBody")
			delaySeconds := int32(request.GetIntParam(req.Parameters, "SendMessageBatchRequestEntry."+strconv.Itoa(i)+".DelaySeconds"))

			message := sqsstore.NewMessage(messageBody)
			message.DelaySeconds = delaySeconds
			message.MessageGroupID = request.GetParamCaseInsensitive(req.Parameters, "SendMessageBatchRequestEntry."+strconv.Itoa(i)+".MessageGroupId")
			message.MessageDeduplicationID = request.GetParamCaseInsensitive(req.Parameters, "SendMessageBatchRequestEntry."+strconv.Itoa(i)+".MessageDeduplicationId")

			msgAttrs := make(map[string]*sqsstore.MessageAttributeValue)
			for j := 1; ; j++ {
				attrName := request.GetParamCaseInsensitive(req.Parameters, "SendMessageBatchRequestEntry."+strconv.Itoa(i)+".MessageAttribute."+strconv.Itoa(j)+".Name")
				if attrName == "" {
					break
				}
				dataType := request.GetParamCaseInsensitive(req.Parameters, "SendMessageBatchRequestEntry."+strconv.Itoa(i)+".MessageAttribute."+strconv.Itoa(j)+".Value.DataType")
				if dataType == "" {
					break
				}
				attr := &sqsstore.MessageAttributeValue{DataType: dataType}
				if sv := request.GetParamCaseInsensitive(req.Parameters, "SendMessageBatchRequestEntry."+strconv.Itoa(i)+".MessageAttribute."+strconv.Itoa(j)+".Value.StringValue"); sv != "" {
					attr.StringValue = &sv
				}
				if bv := request.GetParamCaseInsensitive(req.Parameters, "SendMessageBatchRequestEntry."+strconv.Itoa(i)+".MessageAttribute."+strconv.Itoa(j)+".Value.BinaryValue"); bv != "" {
					attr.BinaryValue = sqsstore.DecodeBinaryValue(bv)
				}
				msgAttrs[attrName] = attr
			}
			if len(msgAttrs) > 0 {
				message.MessageAttributes = msgAttrs
			}

			created, err := store.SendMessage(queueURL, message)
			if err != nil {
				code, senderFault := mapStoreErrorToBatchCode(err)
				entries = append(entries, map[string]interface{}{
					"Id":          id,
					"SenderFault": senderFault,
					"Code":        code,
					"Message":     err.Error(),
				})
				continue
			}

			batchEntry := map[string]interface{}{
				"Id":                     id,
				"MessageId":              created.ID,
				"MD5OfMessageBody":       created.MD5OfBody,
				"MD5OfMessageAttributes": created.MD5OfMessageAttributes,
			}
			if created.SequenceNumber != "" {
				batchEntry["SequenceNumber"] = created.SequenceNumber
			}
			entries = append(entries, batchEntry)
		}
	}

	if len(entries) == 0 {
		return nil, ErrEmptyBatchRequest
	}

	successEntries := make([]map[string]interface{}, 0)
	failedEntries := make([]map[string]interface{}, 0)
	for _, entry := range entries {
		if _, ok := entry["MessageId"]; ok {
			successEntries = append(successEntries, entry)
		} else {
			failedEntries = append(failedEntries, entry)
		}
	}

	return map[string]interface{}{
		"Successful": successEntries,
		"Failed":     failedEntries,
	}, nil
}

// ReceiveMessage receives one or more messages from an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_ReceiveMessage.html
func (s *SQSService) ReceiveMessage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueURL := request.GetParamCaseInsensitive(req.Parameters, "QueueUrl")
	if queueURL == "" {
		return nil, ErrMissingParameter
	}

	maxNumberOfMessages := int32(request.GetIntParam(req.Parameters, "MaxNumberOfMessages"))
	if maxNumberOfMessages <= 0 {
		maxNumberOfMessages = 1
	}

	var visibilityTimeoutPtr *int32
	if vtStr := request.GetParamCaseInsensitive(req.Parameters, "VisibilityTimeout"); vtStr != "" {
		vt := int32(request.GetIntParam(req.Parameters, "VisibilityTimeout"))
		visibilityTimeoutPtr = &vt
	}
	waitTimeSeconds := int32(request.GetIntParam(req.Parameters, "WaitTimeSeconds"))

	// Parse system attribute names: support both legacy AttributeNames
	// and newer MessageSystemAttributeNames. Specifying both is an error.
	legacyAttrNames := request.GetStringList(req.Parameters, "AttributeNames")
	newSysAttrNames := request.GetStringList(req.Parameters, "MessageSystemAttributeNames")

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
	msgAttrNames := request.GetStringList(req.Parameters, "MessageAttributeNames")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	messages, err := store.ReceiveMessage(queueURL, maxNumberOfMessages, visibilityTimeoutPtr, waitTimeSeconds)
	if err != nil {
		return nil, convertStoreError(err)
	}

	sysAttrAll := shouldReturnAllAttributes(sysAttrNames)
	msgAttrAll := shouldReturnAllAttributes(msgAttrNames)

	messageList := make([]map[string]interface{}, 0, len(messages))
	for _, msg := range messages {
		msgMap := map[string]interface{}{
			"MessageId":              msg.ID,
			"ReceiptHandle":          msg.ReceiptHandle,
			"MD5OfBody":              msg.MD5OfBody,
			"Body":                   msg.Body,
			"MD5OfMessageAttributes": msg.MD5OfMessageAttributes,
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

// DeleteMessage deletes a message from an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_DeleteMessage.html
func (s *SQSService) DeleteMessage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueURL := request.GetParamCaseInsensitive(req.Parameters, "QueueUrl")
	if queueURL == "" {
		return nil, ErrMissingParameter
	}

	receiptHandle := request.GetParamCaseInsensitive(req.Parameters, "ReceiptHandle")
	if receiptHandle == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteMessage(queueURL, receiptHandle); err != nil {
		return nil, convertStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// DeleteMessageBatch deletes multiple messages from an SQS queue in a single request.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_DeleteMessageBatch.html
func (s *SQSService) DeleteMessageBatch(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueURL := request.GetParamCaseInsensitive(req.Parameters, "QueueUrl")
	if queueURL == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	successEntries := make([]map[string]interface{}, 0)
	failedEntries := make([]map[string]interface{}, 0)
	seenIDs := make(map[string]bool)
	entryCount := 0

	if entries, ok := req.Parameters["Entries"].([]interface{}); ok && len(entries) > 0 {
		for _, entry := range entries {
			entryMap, ok := entry.(map[string]interface{})
			if !ok {
				continue
			}

			id, _ := entryMap["Id"].(string)
			receiptHandle, _ := entryMap["ReceiptHandle"].(string)

			if id == "" {
				continue
			}

			if err := sqsstore.ValidateBatchEntryId(id); err != nil {
				return nil, ErrInvalidBatchEntryId
			}

			if seenIDs[id] {
				return nil, ErrBatchEntryIdsNotDistinct
			}
			seenIDs[id] = true
			entryCount++

			if entryCount > 10 {
				return nil, ErrTooManyEntriesInBatch
			}

			if err := store.DeleteMessage(queueURL, receiptHandle); err != nil {
				failedEntries = append(failedEntries, map[string]interface{}{
					"Id":          id,
					"SenderFault": false,
					"Code":        "ReceiptHandleIsInvalid",
					"Message":     err.Error(),
				})
				continue
			}

			successEntries = append(successEntries, map[string]interface{}{
				"Id": id,
			})
		}
	} else {
		for i := 1; ; i++ {
			id := request.GetParamCaseInsensitive(req.Parameters, "DeleteMessageBatchRequestEntry."+strconv.Itoa(i)+".Id")
			if id == "" {
				idKey := "DeleteMessageBatchRequestEntry." + strconv.Itoa(i) + ".Id"
				if val, ok := req.Parameters[idKey].(string); ok {
					id = val
				}
			}
			if id == "" {
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

			if entryCount > 10 {
				return nil, ErrTooManyEntriesInBatch
			}

			receiptHandle := request.GetParamCaseInsensitive(req.Parameters, "DeleteMessageBatchRequestEntry."+strconv.Itoa(i)+".ReceiptHandle")
			if receiptHandle == "" {
				rhKey := "DeleteMessageBatchRequestEntry." + strconv.Itoa(i) + ".ReceiptHandle"
				if val, ok := req.Parameters[rhKey].(string); ok {
					receiptHandle = val
				}
			}

			if err := store.DeleteMessage(queueURL, receiptHandle); err != nil {
				failedEntries = append(failedEntries, map[string]interface{}{
					"Id":          id,
					"SenderFault": false,
					"Code":        "ReceiptHandleIsInvalid",
					"Message":     err.Error(),
				})
				continue
			}

			successEntries = append(successEntries, map[string]interface{}{
				"Id": id,
			})
		}
	}

	if entryCount == 0 {
		return nil, ErrEmptyBatchRequest
	}

	return map[string]interface{}{
		"Successful": successEntries,
		"Failed":     failedEntries,
	}, nil
}

// ChangeMessageVisibility changes the visibility timeout of a message.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_ChangeMessageVisibility.html
func (s *SQSService) ChangeMessageVisibility(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueURL := request.GetParamCaseInsensitive(req.Parameters, "QueueUrl")
	if queueURL == "" {
		return nil, ErrMissingParameter
	}

	receiptHandle := request.GetParamCaseInsensitive(req.Parameters, "ReceiptHandle")
	if receiptHandle == "" {
		return nil, ErrMissingParameter
	}

	visibilityTimeout := int32(request.GetIntParam(req.Parameters, "VisibilityTimeout"))

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.ChangeMessageVisibility(queueURL, receiptHandle, visibilityTimeout); err != nil {
		return nil, convertStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// ChangeMessageVisibilityBatch changes the visibility timeout of multiple messages in a single request.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_ChangeMessageVisibilityBatch.html
func (s *SQSService) ChangeMessageVisibilityBatch(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueURL := request.GetParamCaseInsensitive(req.Parameters, "QueueUrl")
	if queueURL == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	entries := make([]map[string]interface{}, 0)
	seenIDs := make(map[string]bool)
	entryCount := 0

	if jsonEntries, ok := req.Parameters["Entries"].([]interface{}); ok && len(jsonEntries) > 0 {
		for _, entry := range jsonEntries {
			entryMap, ok := entry.(map[string]interface{})
			if !ok {
				continue
			}

			id := request.GetStringParam(entryMap, "Id")
			receiptHandle := request.GetStringParam(entryMap, "ReceiptHandle")
			if id == "" {
				continue
			}

			if err := sqsstore.ValidateBatchEntryId(id); err != nil {
				return nil, ErrInvalidBatchEntryId
			}

			if seenIDs[id] {
				return nil, ErrBatchEntryIdsNotDistinct
			}
			seenIDs[id] = true
			entryCount++

			if entryCount > 10 {
				return nil, ErrTooManyEntriesInBatch
			}

			visibilityTimeout := int32(request.GetIntParam(entryMap, "VisibilityTimeout"))

			if err := store.ChangeMessageVisibility(queueURL, receiptHandle, visibilityTimeout); err != nil {
				code, senderFault := mapStoreErrorToBatchCode(err)
				entries = append(entries, map[string]interface{}{
					"Id":          id,
					"SenderFault": senderFault,
					"Code":        code,
					"Message":     err.Error(),
				})
				continue
			}

			entries = append(entries, map[string]interface{}{
				"Id": id,
			})
		}
	} else {
		for i := 1; ; i++ {
			id := request.GetParamCaseInsensitive(req.Parameters, "ChangeMessageVisibilityBatchRequestEntry."+strconv.Itoa(i)+".Id")
			if id == "" {
				idKey := "ChangeMessageVisibilityBatchRequestEntry." + strconv.Itoa(i) + ".Id"
				if val, ok := req.Parameters[idKey].(string); ok {
					id = val
				}
			}
			if id == "" {
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

			if entryCount > 10 {
				return nil, ErrTooManyEntriesInBatch
			}

			receiptHandle := request.GetParamCaseInsensitive(req.Parameters, "ChangeMessageVisibilityBatchRequestEntry."+strconv.Itoa(i)+".ReceiptHandle")
			if receiptHandle == "" {
				rhKey := "ChangeMessageVisibilityBatchRequestEntry." + strconv.Itoa(i) + ".ReceiptHandle"
				if val, ok := req.Parameters[rhKey].(string); ok {
					receiptHandle = val
				}
			}

			visibilityTimeout := int32(request.GetIntParam(req.Parameters, "ChangeMessageVisibilityBatchRequestEntry."+strconv.Itoa(i)+".VisibilityTimeout"))

			if err := store.ChangeMessageVisibility(queueURL, receiptHandle, visibilityTimeout); err != nil {
				code, senderFault := mapStoreErrorToBatchCode(err)
				entries = append(entries, map[string]interface{}{
					"Id":          id,
					"SenderFault": senderFault,
					"Code":        code,
					"Message":     err.Error(),
				})
				continue
			}

			entries = append(entries, map[string]interface{}{
				"Id": id,
			})
		}
	}

	if len(entries) == 0 {
		return nil, ErrEmptyBatchRequest
	}

	successEntries := make([]map[string]interface{}, 0)
	failedEntries := make([]map[string]interface{}, 0)
	for _, entry := range entries {
		if _, ok := entry["SenderFault"]; ok {
			failedEntries = append(failedEntries, entry)
		} else {
			successEntries = append(successEntries, entry)
		}
	}

	return map[string]interface{}{
		"Successful": successEntries,
		"Failed":     failedEntries,
	}, nil
}

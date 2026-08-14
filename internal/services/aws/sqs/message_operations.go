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

// messageEntrySize calculates the total byte size of a batch entry for the
// BatchRequestTooLong check. Includes message body and all attribute data.
// Delegates to the store helper so the size definition stays singular.
func messageEntrySize(body string, attrs map[string]*sqsstore.MessageAttributeValue) int {
	return sqsstore.MessageSize(body, attrs)
}

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
	if errors.Is(err, sqsstore.ErrOverLimit) {
		return ErrOverLimit
	}
	if errors.Is(err, sqsstore.ErrInvalidAttributeValue) {
		return ErrInvalidAttributeValue
	}
	if errors.Is(err, sqsstore.ErrBatchRequestTooLong) {
		return ErrBatchRequestTooLong
	}
	if errors.Is(err, sqsstore.ErrInvalidDataType) {
		return ErrInvalidParameterValue
	}
	if errors.Is(err, sqsstore.ErrTaskAlreadyTerminal) {
		return ErrInvalidParameterValue
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
	systemAttrs, err := parseMessageSystemAttributes(req.Parameters)
	if err != nil {
		return nil, err
	}
	if traceHeader, ok := systemAttrs["AWSTraceHeader"]; ok && traceHeader.StringValue != nil {
		if message.Attributes == nil {
			message.Attributes = make(map[string]string)
		}
		message.Attributes["AWSTraceHeader"] = *traceHeader.StringValue
	}

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
	if len(systemAttrs) > 0 {
		response["MD5OfMessageSystemAttributes"] = sqsstore.CalculateMessageAttributesMD5(systemAttrs)
	}
	if created.SequenceNumber != "" {
		response["SequenceNumber"] = created.SequenceNumber
	}
	return response, nil
}

// batchSendEntry holds a parsed and validated SendMessageBatch entry awaiting
// dispatch. All validation is completed before any SendMessage call.
type batchSendEntry struct {
	id          string
	message     *sqsstore.Message
	systemAttrs map[string]*sqsstore.MessageAttributeValue
}

// SendMessageBatch sends multiple messages to an SQS queue in a single request.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_SendMessageBatch.html
//
// Two-pass design — all entries are parsed and validated before any message
// is sent. This prevents partial sends when a later entry fails validation.
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

	// Pass 1: Parse and validate all entries (no sends).
	parsed, err := parseBatchSendEntries(req.Parameters)
	if err != nil {
		return nil, err
	}

	// Pass 2: Send each validated message.
	successEntries := make([]map[string]interface{}, 0)
	failedEntries := make([]map[string]interface{}, 0)

	for _, e := range parsed {
		created, sendErr := store.SendMessage(queueURL, e.message)
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
			"Id":                     e.id,
			"MessageId":              created.ID,
			"MD5OfMessageBody":       created.MD5OfBody,
			"MD5OfMessageAttributes": created.MD5OfMessageAttributes,
		}
		if len(e.systemAttrs) > 0 {
			entry["MD5OfMessageSystemAttributes"] = sqsstore.CalculateMessageAttributesMD5(e.systemAttrs)
		}
		if created.SequenceNumber != "" {
			entry["SequenceNumber"] = created.SequenceNumber
		}
		successEntries = append(successEntries, entry)
	}

	return map[string]interface{}{
		"Successful": successEntries,
		"Failed":     failedEntries,
	}, nil
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
	if len(jsonEntries) > 10 {
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

	for i := 1; i <= 10; i++ {
		prefix := "SendMessageBatchRequestEntry." + strconv.Itoa(i) + "."
		id := request.GetParamCaseInsensitive(params, prefix+"Id")
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

		messageBody := request.GetParamCaseInsensitive(params, prefix+"MessageBody")
		delaySeconds := int32(request.GetIntParam(params, prefix+"DelaySeconds"))

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

// ReceiveMessage receives one or more messages from an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_ReceiveMessage.html
func (s *SQSService) ReceiveMessage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueURL := request.GetParamCaseInsensitive(req.Parameters, "QueueUrl")
	if queueURL == "" {
		return nil, ErrMissingParameter
	}

	maxNumberOfMessages := int32(request.GetIntParam(req.Parameters, "MaxNumberOfMessages"))
	if err := validateMaxNumberOfMessages(maxNumberOfMessages); err != nil {
		// Treat unset (0) as default 1; reject only when explicitly out of range.
		if maxNumberOfMessages != 0 {
			return nil, err
		}
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

	// Parse ReceiveRequestAttemptId for FIFO receive dedup
	receiveRequestAttemptId := request.GetParamCaseInsensitive(req.Parameters, "ReceiveRequestAttemptId")
	if err := validateReceiveRequestAttemptId(receiveRequestAttemptId); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	messages, err := store.ReceiveMessage(queueURL, maxNumberOfMessages, visibilityTimeoutPtr, waitTimeSeconds, receiveRequestAttemptId)
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

// deleteBatchEntry holds a parsed DeleteMessageBatch entry awaiting execution.
type deleteBatchEntry struct {
	id            string
	receiptHandle string
}

// parseDeleteBatchEntries extracts entries from both JSON and query formats.
func parseDeleteBatchEntries(params map[string]interface{}) ([]deleteBatchEntry, error) {
	if entries, ok := params["Entries"].([]interface{}); ok && len(entries) > 0 {
		result := make([]deleteBatchEntry, 0, len(entries))
		for _, entry := range entries {
			entryMap, ok := entry.(map[string]interface{})
			if !ok {
				continue
			}
			id, _ := entryMap["Id"].(string)
			if id == "" {
				continue
			}
			receiptHandle, _ := entryMap["ReceiptHandle"].(string)
			result = append(result, deleteBatchEntry{id: id, receiptHandle: receiptHandle})
		}
		if len(result) == 0 {
			return nil, ErrEmptyBatchRequest
		}
		return result, nil
	}

	result := make([]deleteBatchEntry, 0)
	for i := 1; ; i++ {
		id := request.GetParamCaseInsensitive(params, "DeleteMessageBatchRequestEntry."+strconv.Itoa(i)+".Id")
		if id == "" {
			idKey := "DeleteMessageBatchRequestEntry." + strconv.Itoa(i) + ".Id"
			if val, ok := params[idKey].(string); ok {
				id = val
			}
		}
		if id == "" {
			break
		}
		receiptHandle := request.GetParamCaseInsensitive(params, "DeleteMessageBatchRequestEntry."+strconv.Itoa(i)+".ReceiptHandle")
		if receiptHandle == "" {
			rhKey := "DeleteMessageBatchRequestEntry." + strconv.Itoa(i) + ".ReceiptHandle"
			if val, ok := params[rhKey].(string); ok {
				receiptHandle = val
			}
		}
		result = append(result, deleteBatchEntry{id: id, receiptHandle: receiptHandle})
	}
	if len(result) == 0 {
		return nil, ErrEmptyBatchRequest
	}
	return result, nil
}

// DeleteMessageBatch deletes multiple messages from an SQS queue in a single request.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_DeleteMessageBatch.html
//
// Two-pass design — all entries are parsed and validated before any deletion.
// This prevents partial state when a later entry fails validation.
func (s *SQSService) DeleteMessageBatch(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueURL := request.GetParamCaseInsensitive(req.Parameters, "QueueUrl")
	if queueURL == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Pass 1: Parse and validate all entries (no mutations).
	entries, err := parseDeleteBatchEntries(req.Parameters)
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
	if len(entries) > 10 {
		return nil, ErrTooManyEntriesInBatch
	}

	// Pass 2: Execute deletions.
	successEntries := make([]map[string]interface{}, 0)
	failedEntries := make([]map[string]interface{}, 0)

	for _, e := range entries {
		if err := store.DeleteMessage(queueURL, e.receiptHandle); err != nil {
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

type changeVisibilityBatchEntry struct {
	id                string
	receiptHandle     string
	visibilityTimeout int32
}

// parseChangeVisibilityBatchEntries extracts entries from both JSON and query
// formats.
func parseChangeVisibilityBatchEntries(params map[string]interface{}) ([]changeVisibilityBatchEntry, error) {
	if jsonEntries, ok := params["Entries"].([]interface{}); ok && len(jsonEntries) > 0 {
		result := make([]changeVisibilityBatchEntry, 0, len(jsonEntries))
		for _, entry := range jsonEntries {
			entryMap, ok := entry.(map[string]interface{})
			if !ok {
				continue
			}
			id := request.GetStringParam(entryMap, "Id")
			if id == "" {
				continue
			}
			receiptHandle := request.GetStringParam(entryMap, "ReceiptHandle")
			visibilityTimeout := int32(request.GetIntParam(entryMap, "VisibilityTimeout"))
			result = append(result, changeVisibilityBatchEntry{
				id:                id,
				receiptHandle:     receiptHandle,
				visibilityTimeout: visibilityTimeout,
			})
		}
		if len(result) == 0 {
			return nil, ErrEmptyBatchRequest
		}
		return result, nil
	}

	result := make([]changeVisibilityBatchEntry, 0)
	for i := 1; ; i++ {
		id := request.GetParamCaseInsensitive(params, "ChangeMessageVisibilityBatchRequestEntry."+strconv.Itoa(i)+".Id")
		if id == "" {
			idKey := "ChangeMessageVisibilityBatchRequestEntry." + strconv.Itoa(i) + ".Id"
			if val, ok := params[idKey].(string); ok {
				id = val
			}
		}
		if id == "" {
			break
		}
		receiptHandle := request.GetParamCaseInsensitive(params, "ChangeMessageVisibilityBatchRequestEntry."+strconv.Itoa(i)+".ReceiptHandle")
		if receiptHandle == "" {
			rhKey := "ChangeMessageVisibilityBatchRequestEntry." + strconv.Itoa(i) + ".ReceiptHandle"
			if val, ok := params[rhKey].(string); ok {
				receiptHandle = val
			}
		}
		visibilityTimeout := int32(request.GetIntParam(params, "ChangeMessageVisibilityBatchRequestEntry."+strconv.Itoa(i)+".VisibilityTimeout"))
		result = append(result, changeVisibilityBatchEntry{
			id:                id,
			receiptHandle:     receiptHandle,
			visibilityTimeout: visibilityTimeout,
		})
	}
	if len(result) == 0 {
		return nil, ErrEmptyBatchRequest
	}
	return result, nil
}

// ChangeMessageVisibilityBatch changes the visibility timeout for multiple
// messages in a single request.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_ChangeMessageVisibilityBatch.html
//
// Two-pass design — all entries are parsed and validated before any mutation.
func (s *SQSService) ChangeMessageVisibilityBatch(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueURL := request.GetParamCaseInsensitive(req.Parameters, "QueueUrl")
	if queueURL == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Pass 1: Parse and validate all entries (no mutations).
	entries, err := parseChangeVisibilityBatchEntries(req.Parameters)
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
	if len(entries) > 10 {
		return nil, ErrTooManyEntriesInBatch
	}

	// Pass 2: Execute visibility changes.
	successEntries := make([]map[string]interface{}, 0)
	failedEntries := make([]map[string]interface{}, 0)

	for _, e := range entries {
		if err := store.ChangeMessageVisibility(queueURL, e.receiptHandle, e.visibilityTimeout); err != nil {
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

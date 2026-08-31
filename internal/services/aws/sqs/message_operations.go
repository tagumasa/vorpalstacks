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
	if errors.Is(err, sqsstore.ErrQueueDeletedRecently) {
		return ErrQueueDeletedRecently
	}
	if errors.Is(err, sqsstore.ErrInvalidReceiptHandle) {
		return ErrReceiptHandleIsInvalid
	}
	if errors.Is(err, sqsstore.ErrMessageNotInflight) {
		return ErrMessageNotInflight
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
	if errors.Is(err, sqsstore.ErrInvalidMessageContents) {
		return ErrInvalidMessageContents
	}
	if errors.Is(err, sqsstore.ErrTaskAlreadyTerminal) {
		return ErrInvalidParameterValue
	}
	if errors.Is(err, sqsstore.ErrTaskNotFound) {
		return ErrResourceNotFound
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

// messageSendResponseFields builds the response members shared by SendMessage
// and SendMessageBatch successful entries. MD5OfMessageAttributes is only
// present when the message carries attributes, matching the AWS response
// surface.
func messageSendResponseFields(msg *sqsstore.Message) map[string]interface{} {
	response := map[string]interface{}{
		"MessageId":        msg.ID,
		"MD5OfMessageBody": msg.MD5OfBody,
	}
	if len(msg.MessageAttributes) > 0 {
		response["MD5OfMessageAttributes"] = msg.MD5OfMessageAttributes
	}
	if msg.SequenceNumber != "" {
		response["SequenceNumber"] = msg.SequenceNumber
	}
	return response
}

// SendMessage sends a message to an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_SendMessage.html
func (s *SQSService) SendMessage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.sendMessageCore(store, SendMessageInput{
		QueueURL:     request.GetParamCaseInsensitive(req.Parameters, "QueueUrl"),
		MessageBody:  request.GetParamCaseInsensitive(req.Parameters, "MessageBody"),
		DelaySeconds: int32(request.GetIntParam(req.Parameters, "DelaySeconds")),
		Parameters:   req.Parameters,
	})
}

// SendMessageBatch sends multiple messages to an SQS queue in a single request.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_SendMessageBatch.html
//
// Two-pass design (inside the Core) — all entries are parsed and validated
// before any message is sent. This prevents partial sends when a later entry
// fails validation.
func (s *SQSService) SendMessageBatch(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueURL := request.GetParamCaseInsensitive(req.Parameters, "QueueUrl")
	if queueURL == "" {
		if val, ok := req.Parameters["QueueUrl"].(string); ok {
			queueURL = val
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.sendMessageBatchCore(store, SendMessageBatchInput{
		QueueURL:   queueURL,
		Parameters: req.Parameters,
	})
}

// ReceiveMessage receives one or more messages from an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_ReceiveMessage.html
func (s *SQSService) ReceiveMessage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.receiveMessageCore(store, ReceiveMessageInput{
		QueueURL:   request.GetParamCaseInsensitive(req.Parameters, "QueueUrl"),
		Parameters: req.Parameters,
	})
}

// DeleteMessage deletes a message from an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_DeleteMessage.html
func (s *SQSService) DeleteMessage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.deleteMessageCore(store, DeleteMessageInput{
		QueueURL:      request.GetParamCaseInsensitive(req.Parameters, "QueueUrl"),
		ReceiptHandle: request.GetParamCaseInsensitive(req.Parameters, "ReceiptHandle"),
	}); err != nil {
		return nil, err
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
// Two-pass design (inside the Core) — all entries are validated before any
// deletion. This prevents partial state when a later entry fails validation.
func (s *SQSService) DeleteMessageBatch(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.deleteMessageBatchCore(store, DeleteMessageBatchInput{
		QueueURL:   request.GetParamCaseInsensitive(req.Parameters, "QueueUrl"),
		Parameters: req.Parameters,
	})
}

// ChangeMessageVisibility changes the visibility timeout of a message.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_ChangeMessageVisibility.html
func (s *SQSService) ChangeMessageVisibility(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.changeMessageVisibilityCore(store, ChangeMessageVisibilityInput{
		QueueURL:          request.GetParamCaseInsensitive(req.Parameters, "QueueUrl"),
		ReceiptHandle:     request.GetParamCaseInsensitive(req.Parameters, "ReceiptHandle"),
		VisibilityTimeout: int32(request.GetIntParam(req.Parameters, "VisibilityTimeout")),
	}); err != nil {
		return nil, err
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
// Two-pass design (inside the Core) — all entries are validated before any
// mutation.
func (s *SQSService) ChangeMessageVisibilityBatch(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.changeMessageVisibilityBatchCore(store, ChangeMessageVisibilityBatchInput{
		QueueURL:   request.GetParamCaseInsensitive(req.Parameters, "QueueUrl"),
		Parameters: req.Parameters,
	})
}

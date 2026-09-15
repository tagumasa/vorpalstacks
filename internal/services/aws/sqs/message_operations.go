package sqs

import (
	"context"
	"strconv"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// SendMessage sends a message to an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_SendMessage.html
func (s *SQSService) SendMessage(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// DelaySeconds is an Integer member: a present value that is not an
	// integer is a wire-type violation, not an omitted member. An omitted
	// member stays 0, which applies the queue's DelaySeconds attribute
	// ("If you don't specify a value, the default value for the queue
	// applies").
	delaySeconds := int32(0)
	if val, present, perr := request.GetIntParamStrictCaseInsensitive(req.Parameters, "DelaySeconds"); present {
		if perr != nil {
			return nil, ErrSerializationException
		}
		delaySeconds = int32(val)
	}
	return s.sendMessageCore(store, SendMessageInput{
		QueueURL:     request.GetParamCaseInsensitive(req.Parameters, "QueueUrl"),
		MessageBody:  request.GetParamCaseInsensitive(req.Parameters, "MessageBody"),
		DelaySeconds: delaySeconds,
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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.sendMessageBatchCore(store, SendMessageBatchInput{
		QueueURL:   request.GetParamCaseInsensitive(req.Parameters, "QueueUrl"),
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

// idReceiptBatchEntry is the Id+ReceiptHandle wire layout shared by the
// DeleteMessageBatch and ChangeMessageVisibilityBatch entries;
// visibilityTimeout is the ChangeMessageVisibilityBatch extension member and
// stays zero for delete entries. visibilityTimeoutSet distinguishes an
// explicitly provided value (including 0) from an omitted one, which selects
// the queue's VisibilityTimeout attribute in the Core.
type idReceiptBatchEntry struct {
	id                   string
	receiptHandle        string
	visibilityTimeout    int32
	visibilityTimeoutSet bool
}

// parseIdReceiptBatchEntries parses Id+ReceiptHandle batch entries from both
// wire formats (the JSON Entries array and the flattened query
// queryPrefix+N+"." keys). The extension callbacks fill per-entry
// operation-specific members (the visibility batch's VisibilityTimeout) and
// may reject the request. Id is @required on every entry shape: an entry
// without it is rejected, never silently dropped from the batch.
func parseIdReceiptBatchEntries(
	params map[string]interface{},
	queryPrefix string,
	jsonExtension func(entryMap map[string]interface{}, entry *idReceiptBatchEntry) error,
	queryExtension func(params map[string]interface{}, entryPrefix string, entry *idReceiptBatchEntry) error,
) ([]idReceiptBatchEntry, error) {
	if entries, ok := params["Entries"].([]interface{}); ok && len(entries) > 0 {
		result := make([]idReceiptBatchEntry, 0, len(entries))
		for _, entry := range entries {
			entryMap, ok := entry.(map[string]interface{})
			if !ok {
				return nil, ErrInvalidParameterValue
			}
			id := request.GetStringParam(entryMap, "Id")
			if id == "" {
				return nil, ErrInvalidParameterValue
			}
			receiptHandle := request.GetStringParam(entryMap, "ReceiptHandle")
			// ReceiptHandle is @required on both entry shapes: an entry
			// without it rejects the request, never surfacing later as a
			// per-entry Failed ReceiptHandleIsInvalid result.
			if receiptHandle == "" {
				return nil, ErrMissingParameter
			}
			e := idReceiptBatchEntry{
				id:            id,
				receiptHandle: receiptHandle,
			}
			if jsonExtension != nil {
				if err := jsonExtension(entryMap, &e); err != nil {
					return nil, err
				}
			}
			result = append(result, e)
		}
		return result, nil
	}

	result := make([]idReceiptBatchEntry, 0)
	for i := 1; ; i++ {
		entryPrefix := queryPrefix + strconv.Itoa(i) + "."
		id := request.GetParamCaseInsensitive(params, entryPrefix+"Id")
		if id == "" {
			if hasQueryEntryMembers(params, entryPrefix) {
				// An entry-shaped key set without its required Id is a
				// malformed entry that rejects the request, not a list
				// terminator — the JSON arm rejects the same shape.
				return nil, ErrInvalidParameterValue
			}
			break
		}
		receiptHandle := request.GetParamCaseInsensitive(params, entryPrefix+"ReceiptHandle")
		// ReceiptHandle is @required on both entry shapes (see the JSON arm).
		if receiptHandle == "" {
			return nil, ErrMissingParameter
		}
		e := idReceiptBatchEntry{
			id:            id,
			receiptHandle: receiptHandle,
		}
		if queryExtension != nil {
			if err := queryExtension(params, entryPrefix, &e); err != nil {
				return nil, err
			}
		}
		result = append(result, e)
	}
	if len(result) == 0 {
		return nil, ErrEmptyBatchRequest
	}
	return result, nil
}

// parseDeleteBatchEntries extracts DeleteMessageBatch entries from both JSON
// and query formats.
func parseDeleteBatchEntries(params map[string]interface{}) ([]idReceiptBatchEntry, error) {
	return parseIdReceiptBatchEntries(params, "DeleteMessageBatchRequestEntry.", nil, nil)
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
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_ChangeMessageVisibility.html
func (s *SQSService) ChangeMessageVisibility(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// VisibilityTimeout is @required and an Integer member: an omitted value
	// is a missing required parameter (the Core rejects it), and a present
	// non-integer is a wire-type violation.
	visibilityTimeout := int32(0)
	visibilityTimeoutSet := false
	if val, present, perr := request.GetIntParamStrictCaseInsensitive(req.Parameters, "VisibilityTimeout"); present {
		if perr != nil {
			return nil, ErrSerializationException
		}
		visibilityTimeout = int32(val)
		visibilityTimeoutSet = true
	}
	if err := s.changeMessageVisibilityCore(store, ChangeMessageVisibilityInput{
		QueueURL:             request.GetParamCaseInsensitive(req.Parameters, "QueueUrl"),
		ReceiptHandle:        request.GetParamCaseInsensitive(req.Parameters, "ReceiptHandle"),
		VisibilityTimeout:    visibilityTimeout,
		VisibilityTimeoutSet: visibilityTimeoutSet,
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// parseChangeVisibilityBatchEntries extracts ChangeMessageVisibilityBatch
// entries (Id, ReceiptHandle, VisibilityTimeout) from both JSON and query
// formats. The entry-level VisibilityTimeout is Required: No (model and API
// reference): an omitted value is left unset for the Core to default to the
// queue's VisibilityTimeout attribute, and an explicitly provided value —
// including 0 — is honoured as-is.
func parseChangeVisibilityBatchEntries(params map[string]interface{}) ([]idReceiptBatchEntry, error) {
	return parseIdReceiptBatchEntries(params, "ChangeMessageVisibilityBatchRequestEntry.",
		func(entryMap map[string]interface{}, entry *idReceiptBatchEntry) error {
			val, present, perr := request.GetIntParamStrictCaseInsensitive(entryMap, "VisibilityTimeout")
			if perr != nil {
				return ErrSerializationException
			}
			entry.visibilityTimeout = int32(val)
			entry.visibilityTimeoutSet = present
			return nil
		},
		func(params map[string]interface{}, entryPrefix string, entry *idReceiptBatchEntry) error {
			val, present, perr := request.GetIntParamStrictCaseInsensitive(params, entryPrefix+"VisibilityTimeout")
			if perr != nil {
				return ErrSerializationException
			}
			entry.visibilityTimeout = int32(val)
			entry.visibilityTimeoutSet = present
			return nil
		},
	)
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

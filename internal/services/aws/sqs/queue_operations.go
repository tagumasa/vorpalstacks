package sqs

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
)

// CreateQueue creates a new SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_CreateQueue.html
func (s *SQSService) CreateQueue(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueName := request.GetParamCaseInsensitive(req.Parameters, "QueueName")

	attrs, attrErr := getQueryEntriesByStem(req.Parameters, "Attribute", "Name", "Value")
	if attrErr != nil {
		return nil, attrErr
	}
	if len(attrs) == 0 {
		attrs = request.ParseAttributes(req.Parameters, "Attributes")
	}

	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags"))
	if len(tags) == 0 {
		// The tags member carries xmlName "Tag" + xmlFlattened: the query
		// wire spells the entries Tag.N.Key / Tag.N.Value — the same stem
		// fallback TagQueue's tag parsing applies, under the flattened
		// stems' contiguity rule (a gap or a value-without-Key slot
		// rejects instead of truncating).
		strictTags, tagErr := getQueryEntriesByStem(req.Parameters, "Tag", "Key", "Value")
		if tagErr != nil {
			return nil, tagErr
		}
		tags = strictTags
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.createQueueCore(ctx, store, CreateQueueInput{
		QueueName: queueName,
		Region:    reqCtx.GetRegion(),
		Attrs:     attrs,
		Tags:      tags,
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"QueueUrl": result.QueueURL,
	}, nil
}

// DeleteQueue deletes an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_DeleteQueue.html
func (s *SQSService) DeleteQueue(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueURL := request.GetParamCaseInsensitive(req.Parameters, "QueueUrl")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteQueueCore(store, DeleteQueueInput{QueueURL: queueURL}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetQueueUrl returns the URL of an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_GetQueueUrl.html
func (s *SQSService) GetQueueUrl(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueName := request.GetParamCaseInsensitive(req.Parameters, "QueueName")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.getQueueUrlCore(store, GetQueueUrlInput{
		QueueName:              queueName,
		QueueOwnerAWSAccountID: request.GetParamCaseInsensitive(req.Parameters, "QueueOwnerAWSAccountId"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"QueueUrl": result.QueueURL,
	}, nil
}

// ListQueues lists the SQS queues.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_ListQueues.html
func (s *SQSService) ListQueues(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// MaxResults is documented as "Value range is 1 to 1000"; absent means
	// the default page size. NextToken is only returned when MaxResults was
	// set ("You must set MaxResults to receive a value for NextToken in the
	// response"). A present value that is not an integer is a wire-type
	// violation, not an omitted member.
	maxResults := 0
	maxResultsSet := false
	if val, present, perr := request.GetIntParamStrictCaseInsensitive(req.Parameters, "MaxResults"); present {
		if perr != nil {
			return nil, ErrSerializationException
		}
		maxResults = val
		maxResultsSet = true
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.listQueuesCore(store, ListQueuesInput{
		QueueNamePrefix: request.GetParamCaseInsensitive(req.Parameters, "QueueNamePrefix"),
		MaxResults:      maxResults,
		MaxResultsSet:   maxResultsSet,
		NextToken:       request.GetParamCaseInsensitive(req.Parameters, "NextToken"),
	})
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"QueueUrls": result.QueueURLs,
	}
	if result.NextToken != "" {
		resp["NextToken"] = result.NextToken
	}
	return resp, nil
}

// GetQueueAttributes returns the attributes of an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_GetQueueAttributes.html
func (s *SQSService) GetQueueAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var requestedAttrs []string
	// The AttributeNames member carries xmlName "AttributeName": the query
	// wire spells the indexed list AttributeName.N.
	requestedAttrs, stemErr := getQueryListByStem(req.Parameters, "AttributeName")
	if stemErr != nil {
		return nil, stemErr
	}

	if len(requestedAttrs) == 0 {
		requestedAttrs = request.GetStringList(req.Parameters, "AttributeNames")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.getQueueAttributesCore(store, GetQueueAttributesInput{
		QueueURL:       request.GetParamCaseInsensitive(req.Parameters, "QueueUrl"),
		AttributeNames: requestedAttrs,
	})
	if err != nil {
		return nil, err
	}
	return getQueueAttributesResponse(result), nil
}

// SetQueueAttributes sets the attributes of an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_SetQueueAttributes.html
func (s *SQSService) SetQueueAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	attrs, attrErr := getQueryEntriesByStem(req.Parameters, "Attribute", "Name", "Value")
	if attrErr != nil {
		return nil, attrErr
	}
	if len(attrs) == 0 {
		attrs = request.ParseAttributes(req.Parameters, "Attributes")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.setQueueAttributesCore(ctx, store, SetQueueAttributesInput{
		QueueURL: request.GetParamCaseInsensitive(req.Parameters, "QueueUrl"),
		Region:   reqCtx.GetRegion(),
		Attrs:    attrs,
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// PurgeQueue purges all messages from an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_PurgeQueue.html
func (s *SQSService) PurgeQueue(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := s.purgeQueueCore(store, PurgeQueueInput{
		QueueURL: request.GetParamCaseInsensitive(req.Parameters, "QueueUrl"),
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListDeadLetterSourceQueues lists the dead letter source queues for an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_ListDeadLetterSourceQueues.html
func (s *SQSService) ListDeadLetterSourceQueues(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// MaxResults is an Integer member: a present value that is not an
	// integer is a wire-type violation, not an omitted member. An omitted
	// member stays 0, which the Core replaces with the documented default
	// page of 1000; an explicitly supplied value below 1 is rejected (the
	// documented "Value range is 1 to 1000").
	maxResults := int32(0)
	maxResultsSet := false
	if val, present, perr := request.GetIntParamStrictCaseInsensitive(req.Parameters, "MaxResults"); present {
		if perr != nil {
			return nil, ErrSerializationException
		}
		maxResults = int32(val)
		maxResultsSet = true
	}
	result, err := s.listDeadLetterSourceQueuesCore(store, ListDeadLetterSourceQueuesInput{
		QueueURL:      request.GetParamCaseInsensitive(req.Parameters, "QueueUrl"),
		MaxResults:    maxResults,
		MaxResultsSet: maxResultsSet,
		NextToken:     request.GetParamCaseInsensitive(req.Parameters, "NextToken"),
	})
	if err != nil {
		return nil, err
	}
	return listDeadLetterSourceQueuesResponse(result), nil
}

// StartMessageMoveTask starts a message move task to move messages from one queue to another.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_StartMessageMoveTask.html
func (s *SQSService) StartMessageMoveTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// MaxNumberOfMessagesPerSecond is an Integer member: a present value
	// that is not an integer is a wire-type violation, not an omitted
	// member. An omitted member stays 0, which the Core treats as the
	// system-optimised variable rate ("If this field is left blank, the
	// system will optimize the rate").
	maxMessages := int32(0)
	if val, present, perr := request.GetIntParamStrictCaseInsensitive(req.Parameters, "MaxNumberOfMessagesPerSecond"); present {
		if perr != nil {
			return nil, ErrSerializationException
		}
		maxMessages = int32(val)
	}
	result, err := s.startMessageMoveTaskCore(store, StartMessageMoveTaskInput{
		SourceARN:      request.GetParamCaseInsensitive(req.Parameters, "SourceArn"),
		DestinationARN: request.GetParamCaseInsensitive(req.Parameters, "DestinationArn"),
		MaxMessages:    maxMessages,
	})
	if err != nil {
		return nil, err
	}
	return startMessageMoveTaskResponse(result), nil
}

// CancelMessageMoveTask cancels a message move task.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_CancelMessageMoveTask.html
func (s *SQSService) CancelMessageMoveTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.cancelMessageMoveTaskCore(store, CancelMessageMoveTaskInput{
		TaskHandle: request.GetParamCaseInsensitive(req.Parameters, "TaskHandle"),
	})
	if err != nil {
		return nil, err
	}
	return cancelMessageMoveTaskResponse(result), nil
}

// ListMessageMoveTasks lists the message move tasks for a source queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_ListMessageMoveTasks.html
func (s *SQSService) ListMessageMoveTasks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// MaxResults is an Integer member: a present value that is not an
	// integer is a wire-type violation, not an omitted member. An omitted
	// member stays 0, which the Core replaces with the documented default
	// of the most recent task.
	maxResults := int32(0)
	if val, present, perr := request.GetIntParamStrictCaseInsensitive(req.Parameters, "MaxResults"); present {
		if perr != nil {
			return nil, ErrSerializationException
		}
		maxResults = int32(val)
	}
	result, err := s.listMessageMoveTasksCore(store, ListMessageMoveTasksInput{
		SourceARN:  request.GetParamCaseInsensitive(req.Parameters, "SourceArn"),
		MaxResults: maxResults,
	})
	if err != nil {
		return nil, err
	}
	return listMessageMoveTasksResponse(result), nil
}

// ---------------------------------------------------------------------------
// DTO → wire-map serialisation. Each function emits exactly the response
// member keys its operation's output shape defines (member names × casing
// pinned against the service model by TestResponseVocabularyAgainstModel).
// ---------------------------------------------------------------------------

func getQueueAttributesResponse(result *GetQueueAttributesResult) map[string]interface{} {
	return map[string]interface{}{
		"Attributes": result.Attributes,
	}
}

func listDeadLetterSourceQueuesResponse(result *ListDeadLetterSourceQueuesResult) map[string]interface{} {
	// The output shape's member is queueUrls — the model's own casing,
	// asymmetrical with ListQueues' QueueUrls.
	resp := map[string]interface{}{
		"queueUrls": result.QueueURLs,
	}
	if result.NextToken != "" {
		resp["NextToken"] = result.NextToken
	}
	return resp
}

func startMessageMoveTaskResponse(result *StartMessageMoveTaskResult) map[string]interface{} {
	return map[string]interface{}{
		"TaskHandle": result.TaskHandle,
	}
}

func cancelMessageMoveTaskResponse(result *CancelMessageMoveTaskResult) map[string]interface{} {
	return map[string]interface{}{
		"ApproximateNumberOfMessagesMoved": result.ApproximateNumberOfMessagesMoved,
	}
}

func listMessageMoveTasksResponse(result *ListMessageMoveTasksResult) map[string]interface{} {
	entries := make([]map[string]interface{}, 0, len(result.Results))
	for _, t := range result.Results {
		entry := map[string]interface{}{
			"Status":                           t.Status,
			"SourceArn":                        t.SourceArn,
			"ApproximateNumberOfMessagesMoved": t.ApproximateNumberOfMessagesMoved,
			"StartedTimestamp":                 t.StartedTimestamp,
		}
		// The rate member is only reported when a fixed rate was requested;
		// an unset rate means the system-optimised variable rate.
		if t.MaxNumberOfMessagesPerSecond > 0 {
			entry["MaxNumberOfMessagesPerSecond"] = t.MaxNumberOfMessagesPerSecond
		}
		if t.TaskHandle != "" {
			entry["TaskHandle"] = t.TaskHandle
		}
		if t.DestinationArn != "" {
			entry["DestinationArn"] = t.DestinationArn
		}
		if t.ApproximateNumberOfMessagesToMove > 0 {
			entry["ApproximateNumberOfMessagesToMove"] = t.ApproximateNumberOfMessagesToMove
		}
		if t.FailureReason != "" {
			entry["FailureReason"] = t.FailureReason
		}
		entries = append(entries, entry)
	}
	return map[string]interface{}{
		"Results": entries,
	}
}

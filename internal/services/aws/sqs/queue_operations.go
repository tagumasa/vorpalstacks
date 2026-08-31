package sqs

import (
	"context"
	"strconv"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
)

// CreateQueue creates a new SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_CreateQueue.html
func (s *SQSService) CreateQueue(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueName := request.GetParamCaseInsensitive(req.Parameters, "QueueName")

	attrs := request.ParseQueryAttributes(req.Parameters, "Attribute")
	if len(attrs) == 0 {
		attrs = request.ParseAttributes(req.Parameters, "Attributes")
	}

	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags"))

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

	result, err := s.getQueueUrlCore(store, GetQueueUrlInput{QueueName: queueName})
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
	// response").
	maxResults := 0
	maxResultsSet := false
	if val, ok := request.GetIntParamCaseInsensitive(req.Parameters, "MaxResults"); ok {
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
	for i := 1; ; i++ {
		attrName := request.GetParamCaseInsensitive(req.Parameters, "AttributeName."+strconv.Itoa(i))
		if attrName == "" {
			break
		}
		requestedAttrs = append(requestedAttrs, attrName)
	}

	if len(requestedAttrs) == 0 {
		requestedAttrs = request.GetStringList(req.Parameters, "AttributeNames")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	return s.getQueueAttributesCore(store, GetQueueAttributesInput{
		QueueURL:       request.GetParamCaseInsensitive(req.Parameters, "QueueUrl"),
		AttributeNames: requestedAttrs,
	})
}

// SetQueueAttributes sets the attributes of an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_SetQueueAttributes.html
func (s *SQSService) SetQueueAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	attrs := make(map[string]string)

	for i := 1; ; i++ {
		attrName := request.GetParamCaseInsensitive(req.Parameters, "Attribute."+strconv.Itoa(i)+".Name")
		if attrName == "" {
			attrNameKey := "Attribute." + strconv.Itoa(i) + ".Name"
			if val, ok := req.Parameters[attrNameKey].(string); ok {
				attrName = val
			}
		}
		if attrName == "" {
			break
		}

		attrValue := request.GetParamCaseInsensitive(req.Parameters, "Attribute."+strconv.Itoa(i)+".Value")
		if attrValue == "" {
			attrValueKey := "Attribute." + strconv.Itoa(i) + ".Value"
			if val, ok := req.Parameters[attrValueKey].(string); ok {
				attrValue = val
			}
		}

		attrs[attrName] = attrValue
	}

	if len(attrs) == 0 {
		for k, v := range request.ParseAttributes(req.Parameters, "Attributes") {
			attrs[k] = v
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.setQueueAttributesCore(store, SetQueueAttributesInput{
		QueueURL: request.GetParamCaseInsensitive(req.Parameters, "QueueUrl"),
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
	return s.listDeadLetterSourceQueuesCore(store, ListDeadLetterSourceQueuesInput{
		QueueURL:   request.GetParamCaseInsensitive(req.Parameters, "QueueUrl"),
		MaxResults: int32(request.GetIntParam(req.Parameters, "MaxResults")),
		NextToken:  request.GetParamCaseInsensitive(req.Parameters, "NextToken"),
	})
}

// StartMessageMoveTask starts a message move task to move messages from one queue to another.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_StartMessageMoveTask.html
func (s *SQSService) StartMessageMoveTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.startMessageMoveTaskCore(store, StartMessageMoveTaskInput{
		SourceARN:      request.GetParamCaseInsensitive(req.Parameters, "SourceArn"),
		DestinationARN: request.GetParamCaseInsensitive(req.Parameters, "DestinationArn"),
		MaxMessages:    int32(request.GetIntParam(req.Parameters, "MaxNumberOfMessagesPerSecond")),
	})
}

// CancelMessageMoveTask cancels a message move task.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_CancelMessageMoveTask.html
func (s *SQSService) CancelMessageMoveTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.cancelMessageMoveTaskCore(store, CancelMessageMoveTaskInput{
		TaskHandle: request.GetParamCaseInsensitive(req.Parameters, "TaskHandle"),
	})
}

// ListMessageMoveTasks lists the message move tasks for a source queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_ListMessageMoveTasks.html
func (s *SQSService) ListMessageMoveTasks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return s.listMessageMoveTasksCore(store, ListMessageMoveTasksInput{
		SourceARN:  request.GetParamCaseInsensitive(req.Parameters, "SourceArn"),
		MaxResults: int32(request.GetIntParam(req.Parameters, "MaxResults")),
	})
}

package sqs

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/store/aws/common"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

func queuesHaveSameAttributes(q1, q2 *sqsstore.Queue) bool {
	if q1.VisibilityTimeout != q2.VisibilityTimeout {
		return false
	}
	if q1.MaximumMessageSize != q2.MaximumMessageSize {
		return false
	}
	if q1.MessageRetentionPeriod != q2.MessageRetentionPeriod {
		return false
	}
	if q1.DelaySeconds != q2.DelaySeconds {
		return false
	}
	if q1.ReceiveMessageWaitTimeSeconds != q2.ReceiveMessageWaitTimeSeconds {
		return false
	}
	if q1.FifoQueue != q2.FifoQueue {
		return false
	}
	if q1.ContentBasedDeduplication != q2.ContentBasedDeduplication {
		return false
	}
	if q1.Policy != q2.Policy {
		return false
	}
	return true
}

// CreateQueue creates a new SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_CreateQueue.html
func (s *SQSService) CreateQueue(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueName := request.GetParamCaseInsensitive(req.Parameters, "QueueName")
	if queueName == "" {
		return nil, ErrMissingParameter
	}

	queue := sqsstore.NewQueue(queueName, reqCtx.GetRegion(), reqCtx.GetAccountID())

	attrs := request.ParseQueryAttributes(req.Parameters, "Attribute")
	if len(attrs) == 0 {
		attrs = request.ParseAttributes(req.Parameters, "Attributes")
	}
	for attrName, attrValue := range attrs {
		if queue.Attributes == nil {
			queue.Attributes = make(map[string]string)
		}
		queue.Attributes[attrName] = attrValue

		switch attrName {
		case "VisibilityTimeout":
			if val, err := strconv.ParseInt(attrValue, 10, 32); err == nil {
				if val < 0 || val > 43200 {
					return nil, ErrInvalidParameterValue
				}
				queue.VisibilityTimeout = int32(val)
			} else {
				return nil, ErrInvalidParameterValue
			}
		case "MaximumMessageSize":
			if val, err := strconv.ParseInt(attrValue, 10, 32); err == nil {
				if val < 1024 || val > 262144 {
					return nil, ErrInvalidParameterValue
				}
				queue.MaximumMessageSize = int32(val)
			} else {
				return nil, ErrInvalidParameterValue
			}
		case "MessageRetentionPeriod":
			if val, err := strconv.ParseInt(attrValue, 10, 32); err == nil {
				if val < 60 || val > 1209600 {
					return nil, ErrInvalidParameterValue
				}
				queue.MessageRetentionPeriod = int32(val)
			} else {
				return nil, ErrInvalidParameterValue
			}
		case "DelaySeconds":
			if val, err := strconv.ParseInt(attrValue, 10, 32); err == nil {
				if val < 0 || val > 900 {
					return nil, ErrInvalidParameterValue
				}
				queue.DelaySeconds = int32(val)
			} else {
				return nil, ErrInvalidParameterValue
			}
		case "ReceiveMessageWaitTimeSeconds":
			if val, err := strconv.ParseInt(attrValue, 10, 32); err == nil {
				if val < 0 || val > 20 {
					return nil, ErrInvalidParameterValue
				}
				queue.ReceiveMessageWaitTimeSeconds = int32(val)
			} else {
				return nil, ErrInvalidParameterValue
			}
		case "FifoQueue":
			if val, err := strconv.ParseBool(attrValue); err == nil {
				queue.FifoQueue = val
			} else {
				return nil, ErrInvalidParameterValue
			}
		case "ContentBasedDeduplication":
			if val, err := strconv.ParseBool(attrValue); err == nil {
				queue.ContentBasedDeduplication = val
			} else {
				return nil, ErrInvalidParameterValue
			}
		case "Policy":
			queue.Policy = attrValue
		case "RedrivePolicy":
			var rdp sqsstore.RedrivePolicy
			if err := json.Unmarshal([]byte(attrValue), &rdp); err == nil {
				queue.RedrivePolicy = &rdp
			}
		}
	}

	queue.Tags = tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags"))

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	created, err := store.CreateQueue(queue)
	if err != nil {
		if errors.Is(err, sqsstore.ErrQueueAlreadyExists) {
			existingQueue, getErr := store.GetQueueByName(queueName)
			if getErr != nil {
				return nil, convertStoreError(getErr)
			}
			if !queuesHaveSameAttributes(queue, existingQueue) {
				return nil, ErrQueueNameExists
			}
			return map[string]interface{}{
				"QueueUrl": existingQueue.URL,
			}, nil
		}
		return nil, convertStoreError(err)
	}

	return map[string]interface{}{
		"QueueUrl": created.URL,
	}, nil
}

// DeleteQueue deletes an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_DeleteQueue.html
func (s *SQSService) DeleteQueue(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueURL := request.GetParamCaseInsensitive(req.Parameters, "QueueUrl")
	if queueURL == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.DeleteQueue(queueURL); err != nil {
		return nil, convertStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// GetQueueUrl returns the URL of an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_GetQueueUrl.html
func (s *SQSService) GetQueueUrl(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueName := request.GetParamCaseInsensitive(req.Parameters, "QueueName")
	if queueName == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	queue, err := store.GetQueueByName(queueName)
	if err != nil {
		return nil, convertStoreError(err)
	}

	return map[string]interface{}{
		"QueueUrl": queue.URL,
	}, nil
}

// ListQueues lists the SQS queues.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_ListQueues.html
func (s *SQSService) ListQueues(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := store.ListQueues(getListOptions(req))
	if err != nil {
		return nil, convertStoreError(err)
	}

	queueNamePrefix := request.GetParamCaseInsensitive(req.Parameters, "QueueNamePrefix")

	queueURLs := make([]string, 0, len(result.Items))
	for _, queue := range result.Items {
		if queueNamePrefix != "" {
			queueName := queue.Name
			if queueName == "" {
				parts := strings.Split(queue.URL, "/")
				queueName = parts[len(parts)-1]
			}
			if !strings.HasPrefix(queueName, queueNamePrefix) {
				continue
			}
		}
		queueURLs = append(queueURLs, queue.URL)
	}

	response := map[string]interface{}{
		"QueueUrls": queueURLs,
	}
	if result.NextMarker != "" {
		response["NextToken"] = result.NextMarker
	}
	return response, nil
}

// GetQueueAttributes returns the attributes of an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_GetQueueAttributes.html
func (s *SQSService) GetQueueAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueURL := request.GetParamCaseInsensitive(req.Parameters, "QueueUrl")
	if queueURL == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	queue, err := store.GetQueue(queueURL)
	if err != nil {
		return nil, convertStoreError(err)
	}

	visible, notVisible, delayed := store.GetMessageCounts(queueURL)

	allAttrs := make(map[string]string)
	allAttrs["QueueArn"] = queue.ARN
	allAttrs["ApproximateNumberOfMessages"] = strconv.FormatInt(int64(visible), 10)
	allAttrs["ApproximateNumberOfMessagesNotVisible"] = strconv.FormatInt(int64(notVisible), 10)
	allAttrs["ApproximateNumberOfMessagesDelayed"] = strconv.FormatInt(int64(delayed), 10)
	allAttrs["CreatedTimestamp"] = strconv.FormatInt(queue.CreatedTimestamp.Unix(), 10)
	allAttrs["LastModifiedTimestamp"] = strconv.FormatInt(queue.LastModifiedTimestamp.Unix(), 10)
	allAttrs["VisibilityTimeout"] = strconv.FormatInt(int64(queue.VisibilityTimeout), 10)
	allAttrs["MaximumMessageSize"] = strconv.FormatInt(int64(queue.MaximumMessageSize), 10)
	allAttrs["MessageRetentionPeriod"] = strconv.FormatInt(int64(queue.MessageRetentionPeriod), 10)
	allAttrs["DelaySeconds"] = strconv.FormatInt(int64(queue.DelaySeconds), 10)
	allAttrs["ReceiveMessageWaitTimeSeconds"] = strconv.FormatInt(int64(queue.ReceiveMessageWaitTimeSeconds), 10)

	for k, v := range queue.Attributes {
		allAttrs[k] = v
	}

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

	if len(requestedAttrs) == 0 {
		return map[string]interface{}{
			"Attributes": allAttrs,
		}, nil
	}

	attrs := make(map[string]string)
	for _, attrName := range requestedAttrs {
		if attrName == "All" {
			attrs = allAttrs
			break
		}
		if !sqsstore.IsValidAttributeName(attrName) {
			return nil, ErrInvalidAttributeName
		}
		if val, ok := allAttrs[attrName]; ok {
			attrs[attrName] = val
		}
	}

	return map[string]interface{}{
		"Attributes": attrs,
	}, nil
}

// SetQueueAttributes sets the attributes of an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_SetQueueAttributes.html
func (s *SQSService) SetQueueAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueURL := request.GetParamCaseInsensitive(req.Parameters, "QueueUrl")
	if queueURL == "" {
		return nil, ErrMissingParameter
	}

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

	if len(attrs) > 0 {
		store, err := s.store(reqCtx)
		if err != nil {
			return nil, err
		}
		if err := store.SetQueueAttributes(queueURL, attrs); err != nil {
			return nil, convertStoreError(err)
		}
	}

	return response.EmptyResponse(), nil
}

// PurgeQueue purges all messages from an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_PurgeQueue.html
func (s *SQSService) PurgeQueue(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueURL := request.GetParamCaseInsensitive(req.Parameters, "QueueUrl")
	if queueURL == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.PurgeQueue(queueURL); err != nil {
		return nil, convertStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// ListDeadLetterSourceQueues lists the dead letter source queues for an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_ListDeadLetterSourceQueues.html
func (s *SQSService) ListDeadLetterSourceQueues(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	queueURL := request.GetParamCaseInsensitive(req.Parameters, "QueueUrl")
	if queueURL == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	dlq, err := store.GetQueue(queueURL)
	if err != nil {
		return nil, convertStoreError(err)
	}

	queues, err := store.ListDeadLetterSourceQueues(dlq.ARN)
	if err != nil {
		return nil, convertStoreError(err)
	}

	queueURLs := make([]string, 0, len(queues))
	for _, q := range queues {
		queueURLs = append(queueURLs, q.URL)
	}

	return map[string]interface{}{
		"QueueUrls": queueURLs,
	}, nil
}

// StartMessageMoveTask starts a message move task to move messages from one queue to another.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_StartMessageMoveTask.html
func (s *SQSService) StartMessageMoveTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	sourceARN := request.GetParamCaseInsensitive(req.Parameters, "SourceArn")
	if sourceARN == "" {
		return nil, ErrMissingParameter
	}

	destARN := request.GetParamCaseInsensitive(req.Parameters, "DestinationArn")
	maxMessages := int32(request.GetIntParam(req.Parameters, "MaxNumberOfMessagesPerSecond"))
	if maxMessages == 0 {
		maxMessages = 1000
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	task, err := store.StartMessageMoveTask(sourceARN, destARN, maxMessages)
	if err != nil {
		return nil, convertStoreError(err)
	}

	return map[string]interface{}{
		"TaskHandle": task.TaskId,
	}, nil
}

// CancelMessageMoveTask cancels a message move task.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_CancelMessageMoveTask.html
func (s *SQSService) CancelMessageMoveTask(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskId := request.GetParamCaseInsensitive(req.Parameters, "TaskHandle")
	if taskId == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	task, err := store.CancelMessageMoveTask(taskId)
	if err != nil {
		return nil, convertStoreError(err)
	}

	return map[string]interface{}{
		"TaskId":                           task.TaskId,
		"ApproximateNumberOfMessagesMoved": task.MovedMessages,
	}, nil
}

// ListMessageMoveTasks lists the message move tasks for a source queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_ListMessageMoveTasks.html
func (s *SQSService) ListMessageMoveTasks(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	sourceARN := request.GetParamCaseInsensitive(req.Parameters, "SourceArn")
	if sourceARN == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	tasks, err := store.ListMessageMoveTasks(sourceARN)
	if err != nil {
		return nil, convertStoreError(err)
	}

	var results []interface{}
	for _, t := range tasks {
		results = append(results, map[string]interface{}{
			"TaskHandle":                   t.TaskId,
			"SourceArn":                    t.SourceQueueARN,
			"DestinationArn":               t.DestinationQueueARN,
			"Status":                       t.Status,
			"MaxNumberOfMessagesPerSecond": t.MaxNumberOfMessages,
		})
	}

	return map[string]interface{}{
		"Results": results,
	}, nil
}

// getListOptions extracts list options from the request.
func getListOptions(req *request.ParsedRequest) common.ListOptions {
	opts := common.ListOptions{MaxItems: 1000}
	if maxResults := request.GetParamCaseInsensitive(req.Parameters, "MaxResults"); maxResults != "" {
		if val, err := strconv.Atoi(maxResults); err == nil && val > 0 {
			opts.MaxItems = val
		}
	}
	opts.Marker = request.GetParamCaseInsensitive(req.Parameters, "NextToken")
	return opts
}

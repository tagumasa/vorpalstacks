package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/store/aws/common"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

var typedQueueAttributes = map[string]bool{
	"QueueArn": true, "ApproximateNumberOfMessages": true,
	"ApproximateNumberOfMessagesNotVisible": true,
	"ApproximateNumberOfMessagesDelayed":    true,
	"CreatedTimestamp":                      true, "LastModifiedTimestamp": true,
	"VisibilityTimeout": true, "MaximumMessageSize": true,
	"MessageRetentionPeriod": true, "DelaySeconds": true,
	"ReceiveMessageWaitTimeSeconds": true, "Policy": true,
	"FifoQueue": true, "ContentBasedDeduplication": true,
	"RedrivePolicy": true,
}

// requestAttrsMatchExisting compares only the attributes a CreateQueue
// request actually provides against the existing queue. AWS returns
// QueueNameExists "only if the request includes attributes whose values
// differ from those of the existing queue"; attributes the request omits are
// not compared and the existing queue URL is returned instead.
func requestAttrsMatchExisting(requestAttrs map[string]string, existing *sqsstore.Queue) bool {
	for name, value := range requestAttrs {
		if !queueAttrMatches(name, value, existing) {
			return false
		}
	}
	return true
}

// queueAttrMatches reports whether a single request attribute value matches
// the existing queue. Numeric and boolean attributes are compared by parsed
// value; attributes without a dedicated queue field fall back to the stored
// attribute map (request attributes are persisted verbatim on create).
func queueAttrMatches(name, value string, existing *sqsstore.Queue) bool {
	switch name {
	case "VisibilityTimeout":
		return int32AttrMatches(value, existing.VisibilityTimeout)
	case "MaximumMessageSize":
		return int32AttrMatches(value, existing.MaximumMessageSize)
	case "MessageRetentionPeriod":
		return int32AttrMatches(value, existing.MessageRetentionPeriod)
	case "DelaySeconds":
		return int32AttrMatches(value, existing.DelaySeconds)
	case "ReceiveMessageWaitTimeSeconds":
		return int32AttrMatches(value, existing.ReceiveMessageWaitTimeSeconds)
	case "FifoQueue":
		return boolAttrMatches(value, existing.FifoQueue)
	case "ContentBasedDeduplication":
		return boolAttrMatches(value, existing.ContentBasedDeduplication)
	case "Policy":
		return value == existing.Policy
	case "RedrivePolicy":
		rdp, err := sqsstore.ParseRedrivePolicy(value)
		if err != nil || existing.RedrivePolicy == nil {
			return false
		}
		return rdp.DeadLetterTargetARN == existing.RedrivePolicy.DeadLetterTargetARN &&
			rdp.MaxReceiveCount == existing.RedrivePolicy.MaxReceiveCount
	default:
		return value == existing.Attributes[name]
	}
}

func int32AttrMatches(value string, existing int32) bool {
	n, err := strconv.ParseInt(value, 10, 32)
	return err == nil && int32(n) == existing
}

func boolAttrMatches(value string, existing bool) bool {
	b, err := strconv.ParseBool(value)
	return err == nil && b == existing
}

// applyQueueAttributes validates and applies attribute key-value pairs to a Queue
// struct. Returns ErrInvalidParameterValue for any invalid attribute value.
func applyQueueAttributes(attrs map[string]string, queue *sqsstore.Queue) error {
	for attrName, attrValue := range attrs {
		if queue.Attributes == nil {
			queue.Attributes = make(map[string]string)
		}
		queue.Attributes[attrName] = attrValue

		switch attrName {
		case "VisibilityTimeout":
			if val, err := strconv.ParseInt(attrValue, 10, 32); err == nil {
				if val < 0 || val > int64(sqsstore.MaxVisibilityTimeout) {
					return ErrInvalidParameterValue
				}
				queue.VisibilityTimeout = int32(val)
			} else {
				return ErrInvalidParameterValue
			}
		case "MaximumMessageSize":
			if val, err := strconv.ParseInt(attrValue, 10, 32); err == nil {
				if val < int64(sqsstore.MinMaximumMessageSize) || val > int64(sqsstore.MaxMaximumMessageSize) {
					return ErrInvalidParameterValue
				}
				queue.MaximumMessageSize = int32(val)
			} else {
				return ErrInvalidParameterValue
			}
		case "MessageRetentionPeriod":
			if val, err := strconv.ParseInt(attrValue, 10, 32); err == nil {
				if val < int64(sqsstore.MinMessageRetentionPeriod) || val > int64(sqsstore.MaxMessageRetentionPeriod) {
					return ErrInvalidParameterValue
				}
				queue.MessageRetentionPeriod = int32(val)
			} else {
				return ErrInvalidParameterValue
			}
		case "DelaySeconds":
			if val, err := strconv.ParseInt(attrValue, 10, 32); err == nil {
				if val < int64(sqsstore.MinDelaySeconds) || val > int64(sqsstore.MaxDelaySeconds) {
					return ErrInvalidParameterValue
				}
				queue.DelaySeconds = int32(val)
			} else {
				return ErrInvalidParameterValue
			}
		case "ReceiveMessageWaitTimeSeconds":
			if val, err := strconv.ParseInt(attrValue, 10, 32); err == nil {
				if val < int64(sqsstore.MinReceiveMessageWaitTimeSeconds) || val > int64(sqsstore.MaxReceiveMessageWaitTimeSeconds) {
					return ErrInvalidParameterValue
				}
				queue.ReceiveMessageWaitTimeSeconds = int32(val)
			} else {
				return ErrInvalidParameterValue
			}
		case "FifoQueue":
			if val, err := strconv.ParseBool(attrValue); err == nil {
				queue.FifoQueue = val
			} else {
				return ErrInvalidParameterValue
			}
		case "ContentBasedDeduplication":
			if val, err := strconv.ParseBool(attrValue); err == nil {
				queue.ContentBasedDeduplication = val
			} else {
				return ErrInvalidParameterValue
			}
		case "Policy":
			if err := sqsstore.ValidatePolicyJSON(attrValue); err != nil {
				return convertStoreError(err)
			}
			queue.Policy = attrValue
		case "RedrivePolicy":
			rdp, err := sqsstore.ParseRedrivePolicy(attrValue)
			if err != nil {
				return ErrInvalidParameterValue
			}
			queue.RedrivePolicy = rdp
		case "KmsDataKeyReusePeriodSeconds":
			if val, err := strconv.ParseInt(attrValue, 10, 32); err == nil {
				if val < int64(sqsstore.MinKmsDataKeyReusePeriodSeconds) || val > int64(sqsstore.MaxKmsDataKeyReusePeriodSeconds) {
					return ErrInvalidParameterValue
				}
			} else {
				return ErrInvalidParameterValue
			}
		case "DeduplicationScope":
			if err := sqsstore.ValidateDeduplicationScope(attrValue); err != nil {
				return convertStoreError(err)
			}
		case "FifoThroughputLimit":
			if err := sqsstore.ValidateFifoThroughputLimit(attrValue); err != nil {
				return convertStoreError(err)
			}
		case "SqsManagedSseEnabled":
			if _, err := strconv.ParseBool(attrValue); err != nil {
				return ErrInvalidParameterValue
			}
		case "RedriveAllowPolicy":
			if attrValue != "" {
				if err := sqsstore.ValidateRedriveAllowPolicyJSON(attrValue); err != nil {
					return convertStoreError(err)
				}
			}
		}
	}
	return nil
}

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
	if val, ok := request.GetIntParamCaseInsensitive(req.Parameters, "MaxResults"); ok {
		if val < 1 || val > sqsstore.MaxListResults {
			return nil, ErrInvalidParameterValue
		}
		maxResults = val
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.listQueuesCore(store, ListQueuesInput{
		QueueNamePrefix: request.GetParamCaseInsensitive(req.Parameters, "QueueNamePrefix"),
		MaxResults:      maxResults,
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
	allAttrs["FifoQueue"] = strconv.FormatBool(queue.FifoQueue)
	allAttrs["ContentBasedDeduplication"] = strconv.FormatBool(queue.ContentBasedDeduplication)

	if queue.RedrivePolicy != nil {
		rdpJSON, _ := json.Marshal(map[string]interface{}{
			"deadLetterTargetArn": queue.RedrivePolicy.DeadLetterTargetARN,
			"maxReceiveCount":     queue.RedrivePolicy.MaxReceiveCount,
		})
		allAttrs["RedrivePolicy"] = string(rdpJSON)
	}

	if queue.Policy != "" {
		allAttrs["Policy"] = queue.Policy
	} else if len(queue.Permissions) > 0 {
		policyJSON := buildPolicyFromPermissions(queue.ARN, queue.Permissions)
		if policyJSON != "" {
			allAttrs["Policy"] = policyJSON
		}
	}

	for k, v := range queue.Attributes {
		if _, isKnown := typedQueueAttributes[k]; !isKnown {
			allAttrs[k] = v
		}
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
		// Validate at the HTTP boundary through the same attribute
		// validation path as CreateQueue (applyQueueAttributes). The store
		// re-validates on write as defence-in-depth for other callers.
		if err := applyQueueAttributes(attrs, &sqsstore.Queue{}); err != nil {
			return nil, err
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	// The store call resolves the queue even for an empty attribute map so a
	// nonexistent queue is rejected with QueueDoesNotExist.
	if err := store.SetQueueAttributes(queueURL, attrs); err != nil {
		return nil, convertStoreError(err)
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

	// MaxResults is 1-1000 per the AWS API reference; the default is 1000.
	maxResults := int32(request.GetIntParam(req.Parameters, "MaxResults"))
	if maxResults == 0 {
		maxResults = int32(sqsstore.MaxListResults)
	}
	if maxResults < 1 || maxResults > int32(sqsstore.MaxListResults) {
		return nil, ErrInvalidParameterValue
	}

	opts := common.ListOptions{
		MaxItems: int(maxResults),
		Marker:   request.GetParamCaseInsensitive(req.Parameters, "NextToken"),
	}

	result, err := store.ListDeadLetterSourceQueues(dlq.ARN, opts)
	if err != nil {
		return nil, convertStoreError(err)
	}

	queueURLs := make([]string, 0, len(result.Items))
	for _, q := range result.Items {
		queueURLs = append(queueURLs, q.URL)
	}

	resp := map[string]interface{}{
		"queueUrls": queueURLs,
	}
	if result.IsTruncated && result.NextMarker != "" {
		resp["NextToken"] = result.NextMarker
	}
	return resp, nil
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
	if err := validateMessageMoveRate(maxMessages); err != nil {
		return nil, err
	}
	// An unset rate (0) means the system-optimised variable rate and is
	// stored as-is; AWS documents a fixed-rate maximum of 500 messages per
	// second, so substituting a concrete rate would exceed it.

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

	maxResults := int32(request.GetIntParam(req.Parameters, "MaxResults"))
	if maxResults < 0 {
		return nil, ErrInvalidParameterValue
	}
	if maxResults == 0 {
		maxResults = 1
	}
	if maxResults > 10 {
		maxResults = 10
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	tasks, err := store.ListMessageMoveTasks(sourceARN, maxResults)
	if err != nil {
		return nil, convertStoreError(err)
	}

	var results []interface{}
	for _, t := range tasks {
		entry := map[string]interface{}{
			"Status":                           t.Status,
			"SourceArn":                        t.SourceQueueARN,
			"ApproximateNumberOfMessagesMoved": t.MovedMessages,
			"StartedTimestamp":                 t.StartTime.UnixMilli(),
		}
		// The rate member is only reported when a fixed rate was requested;
		// an unset rate means the system-optimised variable rate.
		if t.MaxNumberOfMessages > 0 {
			entry["MaxNumberOfMessagesPerSecond"] = t.MaxNumberOfMessages
		}
		if t.TaskId != "" {
			entry["TaskHandle"] = t.TaskId
		}
		if t.DestinationQueueARN != "" {
			entry["DestinationArn"] = t.DestinationQueueARN
		}
		if t.ApproximateNumberOfMessagesToMove > 0 {
			entry["ApproximateNumberOfMessagesToMove"] = t.ApproximateNumberOfMessagesToMove
		}
		if t.FailureReason != "" {
			entry["FailureReason"] = t.FailureReason
		}
		results = append(results, entry)
	}

	return map[string]interface{}{
		"Results": results,
	}, nil
}

func buildPrincipalARNs(accountIDs []string) []string {
	arns := make([]string, len(accountIDs))
	for i, id := range accountIDs {
		arns[i] = fmt.Sprintf("arn:aws:iam::%s:root", id)
	}
	return arns
}

func buildPolicyFromPermissions(queueARN string, permissions map[string]*sqsstore.Permission) string {
	type statement struct {
		Sid       string `json:"Sid"`
		Effect    string `json:"Effect"`
		Principal struct {
			AWS interface{} `json:"AWS"`
		} `json:"Principal"`
		Action   interface{} `json:"Action"`
		Resource string      `json:"Resource"`
	}
	type policy struct {
		Version   string      `json:"Version"`
		Id        string      `json:"Id,omitempty"`
		Statement []statement `json:"Statement"`
	}

	p := policy{
		Version:   "2012-10-17",
		Id:        queueARN + "/SQSDefaultPolicy",
		Statement: make([]statement, 0, len(permissions)),
	}

	for _, perm := range permissions {
		s := statement{
			Sid:      perm.Label,
			Effect:   "Allow",
			Resource: queueARN,
		}
		s.Principal.AWS = buildPrincipalARNs(perm.AWSAccountIDs)
		if len(perm.Actions) == 1 {
			s.Action = "sqs:" + perm.Actions[0]
		} else {
			actions := make([]string, len(perm.Actions))
			for i, a := range perm.Actions {
				actions[i] = "sqs:" + a
			}
			s.Action = actions
		}
		p.Statement = append(p.Statement, s)
	}

	b, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(b)
}

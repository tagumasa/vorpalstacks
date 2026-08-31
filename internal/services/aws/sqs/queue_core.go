package sqs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"vorpalstacks/internal/common/kmsutil"

	storecommon "vorpalstacks/internal/store/aws/common"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

// ---------------------------------------------------------------------------
// Transport-agnostic DTOs shared by HTTP handlers and the admin gRPC handler.
// ---------------------------------------------------------------------------

// CreateQueueInput contains the parameters for creating an SQS queue.
type CreateQueueInput struct {
	QueueName string
	Region    string
	Attrs     map[string]string
	Tags      map[string]string
}

// CreateQueueResult contains the output of a create-queue operation.
type CreateQueueResult struct {
	QueueURL string
}

// DeleteQueueInput contains the parameters for deleting an SQS queue.
type DeleteQueueInput struct {
	QueueURL string
}

// ListQueuesInput contains the parameters for listing SQS queues.
// MaxResultsSet distinguishes an explicitly supplied MaxResults from an
// omitted one (the wire value 0 alone is ambiguous with an unset member).
type ListQueuesInput struct {
	QueueNamePrefix string
	MaxResults      int
	MaxResultsSet   bool
	NextToken       string
}

// ListQueuesResult contains the output of a list-queues operation.
type ListQueuesResult struct {
	QueueURLs []string
	NextToken string
}

// GetQueueUrlInput contains the parameters for resolving a queue URL.
type GetQueueUrlInput struct {
	QueueName string
}

// GetQueueUrlResult contains the output of a get-queue-url operation.
type GetQueueUrlResult struct {
	QueueURL string
}

// GetQueueAttributesInput contains the parameters for reading queue
// attributes. AttributeNames carries the parsed requested-attribute list
// (empty means return all).
type GetQueueAttributesInput struct {
	QueueURL       string
	AttributeNames []string
}

// SetQueueAttributesInput contains the parameters for updating queue
// attributes.
type SetQueueAttributesInput struct {
	QueueURL string
	Attrs    map[string]string
}

// PurgeQueueInput contains the parameters for purging a queue.
type PurgeQueueInput struct {
	QueueURL string
}

// ListDeadLetterSourceQueuesInput contains the parameters for listing the
// dead-letter source queues of a queue.
type ListDeadLetterSourceQueuesInput struct {
	QueueURL   string
	MaxResults int32
	NextToken  string
}

// StartMessageMoveTaskInput contains the parameters for starting a message
// move task.
type StartMessageMoveTaskInput struct {
	SourceARN      string
	DestinationARN string
	MaxMessages    int32
}

// CancelMessageMoveTaskInput contains the parameters for cancelling a
// message move task.
type CancelMessageMoveTaskInput struct {
	TaskHandle string
}

// ListMessageMoveTasksInput contains the parameters for listing message
// move tasks.
type ListMessageMoveTasksInput struct {
	SourceARN  string
	MaxResults int32
}

// ---------------------------------------------------------------------------
// Core methods — shared business logic used by both the HTTP API handlers
// and the admin gRPC-Web handler. Each method performs validation, delegates
// to the store, and maps store errors to service-level AWS errors.
// ---------------------------------------------------------------------------

// createQueueCore creates a new queue or returns the existing URL when the same
// name and attributes are submitted (idempotent re-creation). It rejects when
// the name is taken with different attributes.
func (s *SQSService) createQueueCore(ctx context.Context, store sqsstore.SQSStoreInterface, in CreateQueueInput) (*CreateQueueResult, error) {
	if in.QueueName == "" {
		return nil, ErrMissingParameter
	}
	if !isValidQueueName(in.QueueName) {
		return nil, ErrInvalidQueueName
	}

	queue := sqsstore.NewQueue(in.QueueName, in.Region, store.GetAccountID())

	if err := applyQueueAttributes(in.Attrs, queue); err != nil {
		return nil, err
	}

	// Validate FIFO queue naming bidirectionally: FifoQueue=true requires
	// ".fifo" suffix and vice-versa. Standard queue names cannot contain dots.
	if queue.FifoQueue && !strings.HasSuffix(in.QueueName, ".fifo") {
		return nil, ErrInvalidParameterValue
	}
	if !queue.FifoQueue && strings.HasSuffix(in.QueueName, ".fifo") {
		return nil, ErrInvalidParameterValue
	}

	// Validate KMS key existence when a KmsMasterKeyId is set.
	if kmsKey, ok := in.Attrs["KmsMasterKeyId"]; ok && kmsKey != "" && s.kmsChecker != nil {
		if err := s.kmsChecker.CheckKey(ctx, in.Region, kmsKey); err != nil {
			return nil, mapKMSError(err)
		}
	}

	if len(in.Tags) > 0 {
		queue.Tags = in.Tags
	}

	created, err := store.CreateQueue(queue)
	if err != nil {
		if errors.Is(err, sqsstore.ErrQueueAlreadyExists) {
			existingQueue, getErr := store.GetQueueByName(in.QueueName)
			if getErr != nil {
				return nil, convertStoreError(getErr)
			}
			if !requestAttrsMatchExisting(in.Attrs, existingQueue) {
				return nil, ErrQueueNameExists
			}
			return &CreateQueueResult{QueueURL: existingQueue.URL}, nil
		}
		return nil, convertStoreError(err)
	}

	return &CreateQueueResult{QueueURL: created.URL}, nil
}

// deleteQueueCore deletes a queue by URL.
func (s *SQSService) deleteQueueCore(store sqsstore.SQSStoreInterface, in DeleteQueueInput) error {
	if in.QueueURL == "" {
		return ErrMissingParameter
	}
	if err := store.DeleteQueue(in.QueueURL); err != nil {
		return convertStoreError(err)
	}
	return nil
}

// listQueuesCore lists queues with optional prefix filtering and pagination.
// The prefix filter runs inside the store listing so MaxResults counts only
// matching queues. NextToken is only surfaced when MaxResults was explicitly
// set ("You must set MaxResults to receive a value for NextToken").
func (s *SQSService) listQueuesCore(store sqsstore.SQSStoreInterface, in ListQueuesInput) (*ListQueuesResult, error) {
	if in.MaxResults < 0 || in.MaxResults > sqsstore.MaxListResults {
		return nil, ErrInvalidParameterValue
	}
	// MaxResults is documented as "Value range is 1 to 1000"; an explicitly
	// supplied value below 1 is rejected, while an omitted value selects the
	// default page size.
	if in.MaxResultsSet && in.MaxResults < 1 {
		return nil, ErrInvalidParameterValue
	}
	maxResultsSet := in.MaxResults > 0
	maxItems := in.MaxResults
	if !maxResultsSet {
		maxItems = sqsstore.MaxListResults
	}
	result, err := store.ListQueues(storecommon.ListOptions{
		MaxItems: maxItems,
		Marker:   in.NextToken,
	}, in.QueueNamePrefix)
	if err != nil {
		return nil, convertStoreError(err)
	}

	queueURLs := make([]string, 0, len(result.Items))
	for _, queue := range result.Items {
		queueURLs = append(queueURLs, queue.URL)
	}

	nextToken := ""
	if maxResultsSet && result.IsTruncated {
		nextToken = result.NextMarker
	}
	return &ListQueuesResult{
		QueueURLs: queueURLs,
		NextToken: nextToken,
	}, nil
}

// getQueueUrlCore resolves a queue name to its URL.
func (s *SQSService) getQueueUrlCore(store sqsstore.SQSStoreInterface, in GetQueueUrlInput) (*GetQueueUrlResult, error) {
	if in.QueueName == "" {
		return nil, ErrMissingParameter
	}
	queue, err := store.GetQueueByName(in.QueueName)
	if err != nil {
		return nil, convertStoreError(err)
	}
	return &GetQueueUrlResult{QueueURL: queue.URL}, nil
}

// getQueueAttributesCore returns the attributes of an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_GetQueueAttributes.html
func (s *SQSService) getQueueAttributesCore(store sqsstore.SQSStoreInterface, in GetQueueAttributesInput) (interface{}, error) {
	if in.QueueURL == "" {
		return nil, ErrMissingParameter
	}

	queue, err := store.GetQueue(in.QueueURL)
	if err != nil {
		return nil, convertStoreError(err)
	}

	visible, notVisible, delayed := store.GetMessageCounts(in.QueueURL)

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

	requestedAttrs := in.AttributeNames

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

// setQueueAttributesCore sets the attributes of an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_SetQueueAttributes.html
func (s *SQSService) setQueueAttributesCore(store sqsstore.SQSStoreInterface, in SetQueueAttributesInput) error {
	if in.QueueURL == "" {
		return ErrMissingParameter
	}

	if len(in.Attrs) > 0 {
		// Validate through the same attribute validation path as CreateQueue
		// (applyQueueAttributes). The store re-validates on write as
		// defence-in-depth for other callers.
		if err := applyQueueAttributes(in.Attrs, &sqsstore.Queue{}); err != nil {
			return err
		}
	}

	// The store call resolves the queue even for an empty attribute map so a
	// nonexistent queue is rejected with QueueDoesNotExist.
	if err := store.SetQueueAttributes(in.QueueURL, in.Attrs); err != nil {
		return convertStoreError(err)
	}

	return nil
}

// purgeQueueCore purges all messages from an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_PurgeQueue.html
func (s *SQSService) purgeQueueCore(store sqsstore.SQSStoreInterface, in PurgeQueueInput) error {
	if in.QueueURL == "" {
		return ErrMissingParameter
	}

	if err := store.PurgeQueue(in.QueueURL); err != nil {
		return convertStoreError(err)
	}

	return nil
}

// listDeadLetterSourceQueuesCore lists the dead letter source queues for an
// SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_ListDeadLetterSourceQueues.html
func (s *SQSService) listDeadLetterSourceQueuesCore(store sqsstore.SQSStoreInterface, in ListDeadLetterSourceQueuesInput) (interface{}, error) {
	if in.QueueURL == "" {
		return nil, ErrMissingParameter
	}

	dlq, err := store.GetQueue(in.QueueURL)
	if err != nil {
		return nil, convertStoreError(err)
	}

	// MaxResults is 1-1000 per the AWS API reference; the default is 1000.
	maxResults := in.MaxResults
	if maxResults == 0 {
		maxResults = int32(sqsstore.MaxListResults)
	}
	if maxResults < 1 || maxResults > int32(sqsstore.MaxListResults) {
		return nil, ErrInvalidParameterValue
	}

	opts := storecommon.ListOptions{
		MaxItems: int(maxResults),
		Marker:   in.NextToken,
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

// startMessageMoveTaskCore starts a message move task to move messages from
// one queue to another.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_StartMessageMoveTask.html
func (s *SQSService) startMessageMoveTaskCore(store sqsstore.SQSStoreInterface, in StartMessageMoveTaskInput) (interface{}, error) {
	if in.SourceARN == "" {
		return nil, ErrMissingParameter
	}

	if err := validateMessageMoveRate(in.MaxMessages); err != nil {
		return nil, err
	}
	// An unset rate (0) means the system-optimised variable rate and is
	// stored as-is; AWS documents a fixed-rate maximum of 500 messages per
	// second, so substituting a concrete rate would exceed it.

	task, err := store.StartMessageMoveTask(in.SourceARN, in.DestinationARN, in.MaxMessages)
	if err != nil {
		return nil, convertStoreError(err)
	}

	return map[string]interface{}{
		"TaskHandle": task.TaskId,
	}, nil
}

// cancelMessageMoveTaskCore cancels a message move task.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_CancelMessageMoveTask.html
func (s *SQSService) cancelMessageMoveTaskCore(store sqsstore.SQSStoreInterface, in CancelMessageMoveTaskInput) (interface{}, error) {
	if in.TaskHandle == "" {
		return nil, ErrMissingParameter
	}

	task, err := store.CancelMessageMoveTask(in.TaskHandle)
	if err != nil {
		return nil, convertStoreError(err)
	}

	return map[string]interface{}{
		"ApproximateNumberOfMessagesMoved": task.MovedMessages,
	}, nil
}

// listMessageMoveTasksCore lists the message move tasks for a source queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_ListMessageMoveTasks.html
func (s *SQSService) listMessageMoveTasksCore(store sqsstore.SQSStoreInterface, in ListMessageMoveTasksInput) (interface{}, error) {
	if in.SourceARN == "" {
		return nil, ErrMissingParameter
	}

	maxResults := in.MaxResults
	if maxResults < 0 {
		return nil, ErrInvalidParameterValue
	}
	if maxResults == 0 {
		maxResults = 1
	}
	if maxResults > 10 {
		maxResults = 10
	}

	tasks, err := store.ListMessageMoveTasks(in.SourceARN, maxResults)
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

// mapKMSError converts a KMS key-checker error into the appropriate SQS Kms*
// error that matches the AWS SQS API error codes.
func mapKMSError(err error) error {
	switch {
	case errors.Is(err, kmsutil.ErrKeyNotFound):
		return ErrKmsNotFound
	case errors.Is(err, kmsutil.ErrKeyDisabled):
		return ErrKmsDisabled
	case errors.Is(err, kmsutil.ErrKeyInvalidState):
		return ErrKmsInvalidState
	case errors.Is(err, kmsutil.ErrKeyInvalidUsage):
		return ErrKmsInvalidKeyUsage
	default:
		return ErrInvalidParameterValue
	}
}

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

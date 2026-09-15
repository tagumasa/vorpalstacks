package sqs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"vorpalstacks/internal/common/kmsutil"
	"vorpalstacks/internal/utils/aws/arn"

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
	// QueueOwnerAWSAccountID names the account whose queue the caller is
	// resolving; empty means the caller's own account.
	QueueOwnerAWSAccountID string
}

// GetQueueUrlResult contains the output of a get-queue-url operation.
type GetQueueUrlResult struct {
	QueueURL string
}

// GetQueueAttributesInput contains the parameters for reading queue
// attributes. AttributeNames carries the parsed requested-attribute list;
// an empty list is the documented empty-result request, not a request for
// everything ("The AttributeNames parameter is optional, but if you don't
// specify values for this parameter, the request returns empty results."
// — model + API reference; "All" is the documented way to ask for every
// attribute).
type GetQueueAttributesInput struct {
	QueueURL       string
	AttributeNames []string
}

// GetQueueAttributesResult contains the output of reading queue attributes.
type GetQueueAttributesResult struct {
	Attributes map[string]string
}

// SetQueueAttributesInput contains the parameters for updating queue
// attributes. Region names the region the queue lives in, for the
// platform-extension KMS key existence check.
type SetQueueAttributesInput struct {
	QueueURL string
	Region   string
	Attrs    map[string]string
}

// PurgeQueueInput contains the parameters for purging a queue.
type PurgeQueueInput struct {
	QueueURL string
}

// ListDeadLetterSourceQueuesInput contains the parameters for listing the
// dead-letter source queues of a queue. MaxResultsSet distinguishes an
// explicitly supplied MaxResults from an omitted one, so an explicit 0 can
// be rejected against the documented "Value range is 1 to 1000" while an
// omitted member selects the default page.
type ListDeadLetterSourceQueuesInput struct {
	QueueURL      string
	MaxResults    int32
	MaxResultsSet bool
	NextToken     string
}

// ListDeadLetterSourceQueuesResult contains the output of listing the
// dead-letter source queues of a queue.
type ListDeadLetterSourceQueuesResult struct {
	QueueURLs []string
	NextToken string
}

// StartMessageMoveTaskInput contains the parameters for starting a message
// move task.
type StartMessageMoveTaskInput struct {
	SourceARN      string
	DestinationARN string
	MaxMessages    int32
}

// StartMessageMoveTaskResult contains the output of starting a message move
// task.
type StartMessageMoveTaskResult struct {
	TaskHandle string
}

// CancelMessageMoveTaskInput contains the parameters for cancelling a
// message move task.
type CancelMessageMoveTaskInput struct {
	TaskHandle string
}

// CancelMessageMoveTaskResult contains the output of cancelling a message
// move task.
type CancelMessageMoveTaskResult struct {
	ApproximateNumberOfMessagesMoved int64
}

// ListMessageMoveTasksInput contains the parameters for listing message
// move tasks.
type ListMessageMoveTasksInput struct {
	SourceARN  string
	MaxResults int32
}

// MessageMoveTaskDescription carries one entry of the ListMessageMoveTasks
// result, mirroring the model's result-entry shape. Members with a documented
// absent form (DestinationArn — NULL when the task named no destination;
// TaskHandle; MaxNumberOfMessagesPerSecond — NULL for the system-optimised
// rate; ApproximateNumberOfMessagesToMove; FailureReason) are omitted from
// the emitted entry when unset; the moved/failed counters and
// StartedTimestamp always emit their values, zero included.
type MessageMoveTaskDescription struct {
	TaskHandle                        string
	Status                            string
	SourceArn                         string
	DestinationArn                    string
	MaxNumberOfMessagesPerSecond      int32
	ApproximateNumberOfMessagesMoved  int64
	ApproximateNumberOfMessagesToMove int64
	FailureReason                     string
	StartedTimestamp                  int64
}

// ListMessageMoveTasksResult contains the output of listing message move
// tasks.
type ListMessageMoveTasksResult struct {
	Results []MessageMoveTaskDescription
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
	if err := sqsstore.ValidateQueueName(in.QueueName); err != nil {
		return nil, ErrInvalidQueueName
	}

	// The single validation path: every attribute name and value, shared
	// with the store's defence-in-depth re-check.
	if err := sqsstore.ValidateQueueAttributes(in.Attrs); err != nil {
		return nil, convertStoreError(err)
	}

	// FIFO queue naming is bidirectional: FifoQueue=true requires the
	// ".fifo" suffix and vice-versa. Standard queue names cannot contain
	// dots. The rule is the store's shared exported one, checked before any
	// store access on the validated attribute view.
	if err := sqsstore.ValidateFifoQueueName(in.QueueName, sqsstore.ParseBoolAttr(in.Attrs["FifoQueue"])); err != nil {
		return nil, convertStoreError(err)
	}

	queue := sqsstore.NewQueue(in.QueueName, in.Region, store.GetAccountID())

	if err := applyQueueAttributes(in.Attrs, queue); err != nil {
		return nil, err
	}

	// KMS key existence is a platform extension on every write path that
	// sets KmsMasterKeyId (the model carries no Kms-family errors on
	// CreateQueue or SetQueueAttributes).
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
	// QueueOwnerAWSAccountID selects the owning account. The platform runs
	// a single account, so a queue under any other owner does not exist in
	// this deployment: the queue's own account is read from its ARN and a
	// mismatch answers QueueDoesNotExist instead of silently returning the
	// local queue's URL.
	if in.QueueOwnerAWSAccountID != "" {
		if _, _, _, ownerAccount, _ := arn.SplitARN(queue.ARN); ownerAccount != in.QueueOwnerAWSAccountID {
			return nil, ErrQueueDoesNotExist
		}
	}
	return &GetQueueUrlResult{QueueURL: queue.URL}, nil
}

// getQueueAttributesCore returns the attributes of an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_GetQueueAttributes.html
func (s *SQSService) getQueueAttributesCore(store sqsstore.SQSStoreInterface, in GetQueueAttributesInput) (*GetQueueAttributesResult, error) {
	if in.QueueURL == "" {
		return nil, ErrMissingParameter
	}

	queue, err := store.GetQueue(in.QueueURL)
	if err != nil {
		return nil, convertStoreError(err)
	}

	// The attribute list is opt-in: an omitted AttributeNames returns empty
	// results (the queue must still resolve, so QueueDoesNotExist keeps its
	// precedence).
	if len(in.AttributeNames) == 0 {
		return &GetQueueAttributesResult{Attributes: map[string]string{}}, nil
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

	return &GetQueueAttributesResult{Attributes: attrs}, nil
}

// setQueueAttributesCore sets the attributes of an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_SetQueueAttributes.html
func (s *SQSService) setQueueAttributesCore(ctx context.Context, store sqsstore.SQSStoreInterface, in SetQueueAttributesInput) error {
	if in.QueueURL == "" {
		return ErrMissingParameter
	}
	// Attributes is @required (model): an omitted member is rejected rather
	// than becoming a validation-free no-op success.
	if len(in.Attrs) == 0 {
		return ErrMissingParameter
	}

	// The single validation path shared with CreateQueue: names and values.
	if err := sqsstore.ValidateQueueAttributes(in.Attrs); err != nil {
		return convertStoreError(err)
	}

	// Queue existence precedes the platform-extension KMS check: a request
	// naming a nonexistent queue answers QueueDoesNotExist regardless of the
	// key it carries, the same precedence the queue plane's other writes
	// have (the KMS check is an add-on, not the request's target).
	if _, err := store.GetQueue(in.QueueURL); err != nil {
		return convertStoreError(err)
	}

	// KMS key existence is checked on every write path that sets
	// KmsMasterKeyId (platform extension, symmetric with CreateQueue; the
	// model carries no Kms-family errors on either operation).
	if kmsKey, ok := in.Attrs["KmsMasterKeyId"]; ok && kmsKey != "" && s.kmsChecker != nil {
		if err := s.kmsChecker.CheckKey(ctx, in.Region, kmsKey); err != nil {
			return mapKMSError(err)
		}
	}

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
func (s *SQSService) listDeadLetterSourceQueuesCore(store sqsstore.SQSStoreInterface, in ListDeadLetterSourceQueuesInput) (*ListDeadLetterSourceQueuesResult, error) {
	if in.QueueURL == "" {
		return nil, ErrMissingParameter
	}

	// MaxResults is documented "Value range is 1 to 1000" — the same
	// sentence family ListQueues implements as reject-explicit-below-1;
	// an omitted member selects the default page.
	if in.MaxResultsSet && in.MaxResults < 1 {
		return nil, ErrInvalidParameterValue
	}
	maxResults := in.MaxResults
	if maxResults == 0 {
		maxResults = int32(sqsstore.MaxListResults)
	}
	if maxResults > int32(sqsstore.MaxListResults) {
		return nil, ErrInvalidParameterValue
	}

	dlq, err := store.GetQueue(in.QueueURL)
	if err != nil {
		return nil, convertStoreError(err)
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

	nextToken := ""
	// "You must set MaxResults to receive a value for NextToken in the
	// response" (model + API reference) — the same gate ListQueues applies:
	// a truncated default page surfaces no token, so the caller must opt into
	// pagination explicitly.
	if in.MaxResultsSet && result.IsTruncated && result.NextMarker != "" {
		nextToken = result.NextMarker
	}
	return &ListDeadLetterSourceQueuesResult{
		QueueURLs: queueURLs,
		NextToken: nextToken,
	}, nil
}

// startMessageMoveTaskCore starts a message move task to move messages from
// one queue to another.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_StartMessageMoveTask.html
func (s *SQSService) startMessageMoveTaskCore(store sqsstore.SQSStoreInterface, in StartMessageMoveTaskInput) (*StartMessageMoveTaskResult, error) {
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
		// StartMessageMoveTask documents exactly five errors — InvalidAddress,
		// InvalidSecurity, RequestThrottled, ResourceNotFoundException,
		// UnsupportedOperation (model error list, identical on the fetched API
		// reference page). Not-found failures project onto that vocabulary
		// instead of the queue-plane QueueDoesNotExist, which the operation
		// does not carry.
		if errors.Is(err, sqsstore.ErrQueueNotFound) {
			return nil, ErrResourceNotFound
		}
		return nil, convertStoreError(err)
	}

	return &StartMessageMoveTaskResult{
		TaskHandle: task.TaskId,
	}, nil
}

// cancelMessageMoveTaskCore cancels a message move task.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_CancelMessageMoveTask.html
func (s *SQSService) cancelMessageMoveTaskCore(store sqsstore.SQSStoreInterface, in CancelMessageMoveTaskInput) (*CancelMessageMoveTaskResult, error) {
	if in.TaskHandle == "" {
		return nil, ErrMissingParameter
	}

	task, err := store.CancelMessageMoveTask(in.TaskHandle)
	if err != nil {
		return nil, convertStoreError(err)
	}

	return &CancelMessageMoveTaskResult{
		ApproximateNumberOfMessagesMoved: int64(task.MovedMessages),
	}, nil
}

// listMessageMoveTasksCore lists the message move tasks for a source queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_ListMessageMoveTasks.html
func (s *SQSService) listMessageMoveTasksCore(store sqsstore.SQSStoreInterface, in ListMessageMoveTasksInput) (*ListMessageMoveTasksResult, error) {
	if in.SourceARN == "" {
		return nil, ErrMissingParameter
	}

	maxResults := in.MaxResults
	if maxResults < 0 {
		return nil, ErrInvalidParameterValue
	}
	if maxResults == 0 {
		maxResults = sqsstore.DefaultListMessageMoveTasksResults
	}
	if maxResults > sqsstore.MaxListMessageMoveTasksResults {
		maxResults = sqsstore.MaxListMessageMoveTasksResults
	}

	tasks, err := store.ListMessageMoveTasks(in.SourceARN, maxResults)
	if err != nil {
		return nil, convertStoreError(err)
	}

	results := make([]MessageMoveTaskDescription, 0, len(tasks))
	for _, t := range tasks {
		results = append(results, MessageMoveTaskDescription{
			TaskHandle:                        t.TaskId,
			Status:                            t.Status,
			SourceArn:                         t.SourceQueueARN,
			DestinationArn:                    t.DestinationQueueARN,
			MaxNumberOfMessagesPerSecond:      t.MaxNumberOfMessages,
			ApproximateNumberOfMessagesMoved:  int64(t.MovedMessages),
			ApproximateNumberOfMessagesToMove: int64(t.ApproximateNumberOfMessagesToMove),
			FailureReason:                     t.FailureReason,
			StartedTimestamp:                  t.StartTime.UnixMilli(),
		})
	}

	return &ListMessageMoveTasksResult{
		Results: results,
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

// typedQueueAttributes is not a name whitelist: the attribute-name
// vocabulary is the store's 22-name validAttributeNames set (consulted via
// IsValidAttributeName / ValidateQueueAttributes). This map records which
// attribute names getQueueAttributesCore materialises from dedicated queue
// fields and computed counts — everything not listed here is copied verbatim
// from the stored raw attribute map.
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
		if err != nil {
			return false
		}
		// The clear form matches a queue without a redrive policy.
		if rdp == nil {
			return existing.RedrivePolicy == nil
		}
		if existing.RedrivePolicy == nil {
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

// applyQueueAttributes applies validated attribute key-value pairs onto a
// Queue struct: every value lands in the raw attribute map and the typed
// fields are coerced onto the struct. Name and value validation happens
// once, in sqsstore.ValidateQueueAttributes, before this runs; the error
// return only guards a RedrivePolicy parse that validation has already
// admitted.
func applyQueueAttributes(attrs map[string]string, queue *sqsstore.Queue) error {
	if queue.Attributes == nil {
		queue.Attributes = make(map[string]string)
	}
	for attrName, attrValue := range attrs {
		queue.Attributes[attrName] = attrValue

		switch attrName {
		case "VisibilityTimeout":
			queue.VisibilityTimeout = sqsstore.ParseInt32Attr(attrValue)
		case "MaximumMessageSize":
			queue.MaximumMessageSize = sqsstore.ParseInt32Attr(attrValue)
		case "MessageRetentionPeriod":
			queue.MessageRetentionPeriod = sqsstore.ParseInt32Attr(attrValue)
		case "DelaySeconds":
			queue.DelaySeconds = sqsstore.ParseInt32Attr(attrValue)
		case "ReceiveMessageWaitTimeSeconds":
			queue.ReceiveMessageWaitTimeSeconds = sqsstore.ParseInt32Attr(attrValue)
		case "FifoQueue":
			queue.FifoQueue = sqsstore.ParseBoolAttr(attrValue)
		case "ContentBasedDeduplication":
			queue.ContentBasedDeduplication = sqsstore.ParseBoolAttr(attrValue)
		case "Policy":
			queue.Policy = attrValue
		case "RedrivePolicy":
			rdp, err := sqsstore.ParseRedrivePolicy(attrValue)
			if err != nil {
				return ErrInvalidParameterValue
			}
			queue.RedrivePolicy = rdp
			if rdp == nil {
				// The empty value clears the association: the raw attribute
				// is removed so the queue reports RedrivePolicy as unset.
				delete(queue.Attributes, attrName)
				continue
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

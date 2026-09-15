package sqs

import (
	"errors"
	"strconv"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

// ---------------------------------------------------------------------------
// Message-operation Core — single validation + persistence path for the
// message plane (send, receive, delete, visibility). The HTTP handlers in
// message_operations.go are thin adapters: wire parsing, DTO construction,
// the Core call and response serialisation. Parameters carries the raw
// request parameter map so the dual-format (JSON / flattened query) entry
// and attribute parsing keeps a single implementation here.
// ---------------------------------------------------------------------------

// SendMessageInput carries the parameters for sending a single message.
// Parameters holds the raw request map for the message-attribute and
// system-attribute members that arrive in two wire formats.
type SendMessageInput struct {
	QueueURL     string
	MessageBody  string
	DelaySeconds int32
	Parameters   map[string]interface{}
}

// SendMessageBatchInput carries the parameters for a batch send. Parameters
// holds the raw request map; batch entries are parsed and validated inside
// the Core so no entry is sent before every entry has passed validation.
type SendMessageBatchInput struct {
	QueueURL   string
	Parameters map[string]interface{}
}

// ReceiveMessageInput carries the parameters for receiving messages.
// Parameters holds the raw request map for the attribute-name lists and
// receive options that arrive in two wire formats.
type ReceiveMessageInput struct {
	QueueURL   string
	Parameters map[string]interface{}
}

// DeleteMessageInput carries the parameters for deleting a single message.
type DeleteMessageInput struct {
	QueueURL      string
	ReceiptHandle string
}

// DeleteMessageBatchInput carries the parameters for a batch delete.
// Parameters holds the raw request map so entry parsing happens inside the
// Core after the QueueUrl check, preserving the request-level error
// precedence (an empty QueueUrl wins over a malformed entry list).
type DeleteMessageBatchInput struct {
	QueueURL   string
	Parameters map[string]interface{}
}

// ChangeMessageVisibilityInput carries the parameters for changing the
// visibility timeout of a single message.
type ChangeMessageVisibilityInput struct {
	QueueURL      string
	ReceiptHandle string
	// VisibilityTimeoutSet distinguishes an explicit 0 (a valid value that
	// makes the message visible for immediate redelivery) from an omitted
	// member — VisibilityTimeout is @required on
	// ChangeMessageVisibilityRequest, so an omitted member is rejected,
	// never defaulted to 0.
	VisibilityTimeout    int32
	VisibilityTimeoutSet bool
}

// ChangeMessageVisibilityBatchInput carries the parameters for a batch
// visibility change. Parameters holds the raw request map so entry parsing
// happens inside the Core after the QueueUrl check, preserving the
// request-level error precedence.
type ChangeMessageVisibilityBatchInput struct {
	QueueURL   string
	Parameters map[string]interface{}
}

// storeErrorRow is one row of the single store→services error table. The
// single-op projection (convertStoreError) maps a store sentinel to the wire
// error a single-operation request returns; the batch projection
// (mapStoreErrorToBatchCode) renders the same sentinel as the Code and
// SenderFault of one Failed entry in a batch response. Both surfaces of an
// operation family read this one table, so the same store failure cannot map
// to different wire outcomes per surface.
type storeErrorRow struct {
	store  error // store-layer sentinel
	single error // services error for single-operation requests
	// batchCodeOverride carries a per-entry batch Code that intentionally
	// differs from the single error's JSON error code. Empty means the batch
	// Code is derived from the single error's code.
	batchCodeOverride string
}

// storeErrorTable is the complete store→services error vocabulary: every
// store sentinel a Core can observe through an error return is mapped once
// here (ErrQueueAlreadyExists is the one exception — createQueueCore
// handles it with errors.Is before any table projection). The batch Code of
// a row is the JSON error code of its single projection (a batch Failed
// entry reports the error code the same failure produces on the single
// operation) except where the override field documents a divergence.
var storeErrorTable = []storeErrorRow{
	{sqsstore.ErrQueueNotFound, ErrQueueDoesNotExist, ""},
	{sqsstore.ErrQueueDeletedRecently, ErrQueueDeletedRecently, ""},
	{sqsstore.ErrInvalidReceiptHandle, ErrReceiptHandleIsInvalid, ""},
	{sqsstore.ErrMessageNotInflight, ErrMessageNotInflight, ""},
	{sqsstore.ErrMessageTooLarge, ErrMessageTooLarge, ""},
	// The FIFO identifier failures predate the table and carry their own
	// vocabulary in batch Failed entries while the single-op wire code is
	// InvalidParameterValue. No AWS source enumerates the batch Code for
	// these failures, so the existing vocabulary is preserved.
	{sqsstore.ErrMissingMessageGroupId, ErrMissingMessageGroupId, "MissingMessageGroupId"},
	{sqsstore.ErrMissingDeduplicationId, ErrMissingDeduplicationId, "MissingDeduplicationId"},
	{sqsstore.ErrInvalidParameterValue, ErrInvalidParameterValue, ""},
	{sqsstore.ErrPurgeQueueInProgress, ErrPurgeQueueInProgress, ""},
	{sqsstore.ErrTooManyTags, ErrTooManyTags, ""},
	{sqsstore.ErrInvalidTagKey, ErrInvalidTagKey, ""},
	{sqsstore.ErrInvalidTagValue, ErrInvalidTagValue, ""},
	{sqsstore.ErrInvalidQueueName, ErrInvalidQueueName, ""},
	{sqsstore.ErrInvalidAttributeName, ErrInvalidAttributeName, ""},
	{sqsstore.ErrOverLimit, ErrOverLimit, ""},
	{sqsstore.ErrInvalidAttributeValue, ErrInvalidAttributeValue, ""},
	{sqsstore.ErrInvalidDataType, ErrInvalidParameterValue, ""},
	{sqsstore.ErrInvalidMessageContents, ErrInvalidMessageContents, ""},
	{sqsstore.ErrTaskAlreadyTerminal, ErrInvalidParameterValue, ""},
	{sqsstore.ErrTaskNotFound, ErrResourceNotFound, ""},
	{sqsstore.ErrInvalidBatchEntryId, ErrInvalidBatchEntryId, ""},
	// A move source that no queue redrives into violates the SourceArn
	// contract ("only ARNs of dead-letter queues ... are accepted", model).
	{sqsstore.ErrUnsupportedOperation, ErrUnsupportedOperation, ""},
	// A move ARN that does not parse as a queue ARN at all: the malformed-
	// identifier error of the move operations ("The specified ID is
	// invalid.", model — HTTP 404).
	{sqsstore.ErrInvalidAddress, ErrInvalidAddress, ""},
}

// lookupStoreErrorRow finds the table row whose store sentinel matches err.
func lookupStoreErrorRow(err error) *storeErrorRow {
	for i := range storeErrorTable {
		if errors.Is(err, storeErrorTable[i].store) {
			return &storeErrorTable[i]
		}
	}
	return nil
}

// convertStoreError projects a store error onto the single-operation wire
// through the shared table. Unrecognised errors pass through unchanged.
func convertStoreError(err error) error {
	if err == nil {
		return nil
	}
	if row := lookupStoreErrorRow(err); row != nil {
		return row.single
	}
	return err
}

// mapStoreErrorToBatchCode projects a store error onto a batch Failed entry
// (Code, SenderFault) through the shared table. Every mapped store sentinel
// describes a client-induced failure, so SenderFault is true for all of
// them; an unmapped store failure is reported as an internal error.
func mapStoreErrorToBatchCode(err error) (code string, senderFault bool) {
	row := lookupStoreErrorRow(err)
	if row == nil {
		return "InternalError", false
	}
	if row.batchCodeOverride != "" {
		return row.batchCodeOverride, true
	}
	if awsErr, ok := row.single.(*awserrors.AWSError); ok {
		return awsErr.GetCode(), true
	}
	return "InternalError", false
}

// requestsAllAttributes reports whether the attribute name list asks for
// every attribute through a recognised wildcard token ("All" or ".*").
// An empty list asks for nothing: "you can send a list of attribute names
// to receive, or you can return all of the attributes by specifying All
// or .* in your request" — the lists are opt-in, never a default-All.
// Reference: AWS SQS ReceiveMessage API docs.
func requestsAllAttributes(names []string) bool {
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

// applyTraceHeader merges a send's parsed system attributes into the
// message's attribute map: AWSTraceHeader is the only system attribute
// valid on sends, carried as an ordinary message attribute entry.
func applyTraceHeader(message *sqsstore.Message, systemAttrs map[string]*sqsstore.MessageAttributeValue) {
	if th, ok := systemAttrs["AWSTraceHeader"]; ok && th.StringValue != nil {
		if message.Attributes == nil {
			message.Attributes = make(map[string]string)
		}
		message.Attributes["AWSTraceHeader"] = *th.StringValue
	}
}

// sendMessageCore sends a message to an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_SendMessage.html
func (s *SQSService) sendMessageCore(store sqsstore.SQSStoreInterface, in SendMessageInput) (interface{}, error) {
	if in.QueueURL == "" {
		return nil, ErrMissingParameter
	}
	if in.MessageBody == "" {
		return nil, ErrMissingParameter
	}

	if err := rejectAttributeListValues(in.Parameters); err != nil {
		return nil, err
	}

	message := sqsstore.NewMessage(in.MessageBody)
	message.DelaySeconds = in.DelaySeconds
	message.MessageGroupID = request.GetParamCaseInsensitive(in.Parameters, "MessageGroupId")
	message.MessageDeduplicationID = request.GetParamCaseInsensitive(in.Parameters, "MessageDeduplicationId")

	messageAttributes, err := parseRequestMessageAttributes(in.Parameters)
	if err != nil {
		return nil, err
	}
	message.MessageAttributes = messageAttributes

	// Validate the complete message attribute set (count, names, data
	// types) through the single store-layer validator.
	if err := sqsstore.ValidateMessageAttributes(message.MessageAttributes); err != nil {
		return nil, convertStoreError(err)
	}

	// Parse MessageSystemAttributes (only AWSTraceHeader is valid for sends)
	systemAttrs, err := parseSystemAttributes(in.Parameters, "MessageSystemAttribute.")
	if err != nil {
		return nil, err
	}
	applyTraceHeader(message, systemAttrs)

	created, err := store.SendMessage(in.QueueURL, message)
	if err != nil {
		return nil, convertStoreError(err)
	}

	response := messageSendResponseFields(created)
	if len(systemAttrs) > 0 {
		response["MD5OfMessageSystemAttributes"] = sqsstore.CalculateMessageAttributesMD5(systemAttrs)
	}
	return response, nil
}

// sendMessageBatchCore sends multiple messages to an SQS queue in a single
// request.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_SendMessageBatch.html
//
// Two-pass design — all entries are parsed and validated before any message
// is sent. This prevents partial sends when a later entry fails validation.
func (s *SQSService) sendMessageBatchCore(store sqsstore.SQSStoreInterface, in SendMessageBatchInput) (interface{}, error) {
	if in.QueueURL == "" {
		return nil, ErrMissingParameter
	}

	if err := rejectAttributeListValues(in.Parameters); err != nil {
		return nil, err
	}

	// Pass 1: Parse and validate all entries (no sends).
	parsed, err := parseBatchSendEntries(in.Parameters)
	if err != nil {
		return nil, err
	}

	// QueueDoesNotExist is a request-level error: the whole request fails
	// before any entry is processed.
	if _, err := store.GetQueue(in.QueueURL); err != nil {
		return nil, convertStoreError(err)
	}

	// Pass 2: Send each validated message.
	successEntries := make([]map[string]interface{}, 0)
	failedEntries := make([]map[string]interface{}, 0)

	for _, e := range parsed {
		created, sendErr := store.SendMessage(in.QueueURL, e.message)
		if sendErr != nil {
			failedEntries = append(failedEntries, batchFailedEntry(e.id, sendErr))
			continue
		}

		entry := map[string]interface{}{
			"Id": e.id,
		}
		for k, v := range messageSendResponseFields(created) {
			entry[k] = v
		}
		if len(e.systemAttrs) > 0 {
			entry["MD5OfMessageSystemAttributes"] = sqsstore.CalculateMessageAttributesMD5(e.systemAttrs)
		}
		successEntries = append(successEntries, entry)
	}

	return map[string]interface{}{
		"Successful": successEntries,
		"Failed":     failedEntries,
	}, nil
}

// receiveMessageCore receives one or more messages from an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_ReceiveMessage.html
func (s *SQSService) receiveMessageCore(store sqsstore.SQSStoreInterface, in ReceiveMessageInput) (interface{}, error) {
	if in.QueueURL == "" {
		return nil, ErrMissingParameter
	}

	// MaxNumberOfMessages carries the AWS contract "Valid values: 1 to 10.
	// Default: 1." — an explicitly provided value outside the range
	// (including 0) is rejected, a present value that is not an integer is
	// a wire-type violation, and only omission falls back to the default.
	maxNumberOfMessages := int32(1)
	if val, present, perr := request.GetIntParamStrictCaseInsensitive(in.Parameters, "MaxNumberOfMessages"); present {
		if perr != nil {
			return nil, ErrSerializationException
		}
		maxNumberOfMessages = int32(val)
		if err := validateMaxNumberOfMessages(maxNumberOfMessages); err != nil {
			return nil, err
		}
	}

	var visibilityTimeoutPtr *int32
	if vtVal, present, perr := request.GetIntParamStrictCaseInsensitive(in.Parameters, "VisibilityTimeout"); present {
		if perr != nil {
			return nil, ErrSerializationException
		}
		vt := int32(vtVal)
		// The member shares the queue attribute's documented 0–43200 range
		// and is checked here for the same defence-in-depth its sibling
		// members carry, instead of deferring to the store alone.
		if vt < sqsstore.MinVisibilityTimeout || vt > sqsstore.MaxVisibilityTimeout {
			return nil, ErrInvalidParameterValue
		}
		visibilityTimeoutPtr = &vt
	}
	// WaitTimeSeconds: a negative sentinel tells the store the request
	// omitted the parameter so it applies the queue's
	// ReceiveMessageWaitTimeSeconds attribute; explicit values are validated
	// here and long-poll in the store.
	waitTimeSeconds := int32(-1)
	if val, present, perr := request.GetIntParamStrictCaseInsensitive(in.Parameters, "WaitTimeSeconds"); present {
		if perr != nil {
			return nil, ErrSerializationException
		}
		if val < 0 {
			return nil, ErrInvalidParameterValue
		}
		waitTimeSeconds = int32(val)
	}

	// Parse system attribute names: support both legacy AttributeNames
	// and newer MessageSystemAttributeNames. Specifying both is an error.
	legacyAttrNames := request.GetStringList(in.Parameters, "AttributeNames")
	newSysAttrNames := request.GetStringList(in.Parameters, "MessageSystemAttributeNames")
	if len(legacyAttrNames) == 0 && len(newSysAttrNames) == 0 {
		// Both members carry xmlName "AttributeName": on the query wire they
		// are literally the same keys, so a query request expresses exactly
		// one of the two members. The stem read must not feed both lists —
		// that would self-trigger the mutual rejection below on the very
		// form the query protocol makes indistinguishable.
		var stemErr error
		legacyAttrNames, stemErr = getQueryListByStem(in.Parameters, "AttributeName")
		if stemErr != nil {
			return nil, stemErr
		}
	}

	if len(legacyAttrNames) > 0 && len(newSysAttrNames) > 0 {
		return nil, ErrInvalidParameterCombination
	}

	// The two members' filter vocabularies are protocol-enforced on the AWS
	// side and rejected here rather than silently matching nothing. Both
	// target the message-system attribute set the ReceiveMessage reference
	// enumerates (All, SenderId, SentTimestamp, the receive counters, the
	// FIFO identifiers, AWSTraceHeader, DeadLetterQueueSourceArn); the
	// deprecated AttributeNames additionally admits SqsManagedSseEnabled
	// per its own documentation list (its Valid Values row reuses the
	// QueueAttributeName shape). MessageAttributeNames is a plain string
	// with the documented "All"/".*"/prefix wildcards and is not
	// enum-validated.
	for _, name := range legacyAttrNames {
		if !sqsstore.IsValidMessageSystemAttributeName(name) && name != "SqsManagedSseEnabled" {
			return nil, ErrInvalidAttributeName
		}
	}
	for _, name := range newSysAttrNames {
		if !sqsstore.IsValidMessageSystemAttributeName(name) {
			return nil, ErrInvalidAttributeName
		}
	}

	var sysAttrNames []string
	if len(legacyAttrNames) > 0 {
		sysAttrNames = legacyAttrNames
	} else {
		sysAttrNames = newSysAttrNames
	}

	// Parse message attribute names
	msgAttrNames := request.GetStringList(in.Parameters, "MessageAttributeNames")
	if len(msgAttrNames) == 0 {
		// The MessageAttributeNames member carries xmlName
		// "MessageAttributeName": the query wire spells the indexed list
		// MessageAttributeName.N.
		var stemErr error
		msgAttrNames, stemErr = getQueryListByStem(in.Parameters, "MessageAttributeName")
		if stemErr != nil {
			return nil, stemErr
		}
	}

	// Parse ReceiveRequestAttemptId for FIFO receive dedup
	receiveRequestAttemptId := request.GetParamCaseInsensitive(in.Parameters, "ReceiveRequestAttemptId")
	if err := validateReceiveRequestAttemptId(receiveRequestAttemptId); err != nil {
		return nil, err
	}

	messages, err := store.ReceiveMessage(in.QueueURL, maxNumberOfMessages, visibilityTimeoutPtr, waitTimeSeconds, receiveRequestAttemptId)
	if err != nil {
		return nil, convertStoreError(err)
	}

	sysAttrAll := requestsAllAttributes(sysAttrNames)
	msgAttrAll := requestsAllAttributes(msgAttrNames)

	messageList := make([]map[string]interface{}, 0, len(messages))
	for _, msg := range messages {
		// The Message shape carries exactly these four always-present
		// members plus Attributes, MD5OfMessageAttributes and
		// MessageAttributes — no top-level FIFO identifier members exist
		// in the model; MessageGroupId, MessageDeduplicationId and
		// SequenceNumber travel as system attributes (see below).
		msgMap := map[string]interface{}{
			"MessageId":     msg.ID,
			"ReceiptHandle": msg.ReceiptHandle,
			"MD5OfBody":     msg.MD5OfBody,
			"Body":          msg.Body,
		}

		// Message attributes are emitted only when the request asked for
		// them (the name list is opt-in). MD5OfMessageAttributes is the
		// digest of the attributes actually returned, so client-side
		// integrity checks that recompute from the response agree with it
		// under filtering; when the filter selects nothing, both keys are
		// omitted.
		if len(msgAttrNames) > 0 && len(msg.MessageAttributes) > 0 {
			emitted := make(map[string]*sqsstore.MessageAttributeValue, len(msg.MessageAttributes))
			for k, v := range msg.MessageAttributes {
				if msgAttrAll || isRequestedAttribute(k, msgAttrNames) {
					emitted[k] = v
				}
			}
			if len(emitted) > 0 {
				attrs := make(map[string]interface{}, len(emitted))
				for k, v := range emitted {
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
				msgMap["MD5OfMessageAttributes"] = sqsstore.CalculateMessageAttributesMD5(emitted)
			}
		}

		// Filter system attributes — same opt-in rule through their own
		// name list.
		if len(sysAttrNames) > 0 && len(msg.Attributes) > 0 {
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

// deleteMessageCore deletes a message from an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_DeleteMessage.html
func (s *SQSService) deleteMessageCore(store sqsstore.SQSStoreInterface, in DeleteMessageInput) error {
	if in.QueueURL == "" {
		return ErrMissingParameter
	}
	if in.ReceiptHandle == "" {
		return ErrMissingParameter
	}

	if err := store.DeleteMessage(in.QueueURL, in.ReceiptHandle); err != nil {
		return convertStoreError(err)
	}
	return nil
}

// deleteMessageBatchCore deletes multiple messages from an SQS queue in a
// single request.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_DeleteMessageBatch.html
//
// Two-pass design — all entries are validated before any deletion. This
// prevents partial state when a later entry fails validation.
func (s *SQSService) deleteMessageBatchCore(store sqsstore.SQSStoreInterface, in DeleteMessageBatchInput) (interface{}, error) {
	if in.QueueURL == "" {
		return nil, ErrMissingParameter
	}

	// Pass 1: Parse and validate all entries (no mutations).
	entries, err := parseDeleteBatchEntries(in.Parameters)
	if err != nil {
		return nil, err
	}

	ids := make([]string, len(entries))
	for i, e := range entries {
		ids[i] = e.id
	}
	if err := validateBatchEntryIDs(ids); err != nil {
		return nil, err
	}

	// QueueDoesNotExist is a request-level error: the whole request fails
	// before any entry is processed.
	if _, err := store.GetQueue(in.QueueURL); err != nil {
		return nil, convertStoreError(err)
	}

	// Pass 2: Execute deletions.
	successEntries := make([]map[string]interface{}, 0)
	failedEntries := make([]map[string]interface{}, 0)

	for _, e := range entries {
		if err := store.DeleteMessage(in.QueueURL, e.receiptHandle); err != nil {
			failedEntries = append(failedEntries, batchFailedEntry(e.id, err))
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

// changeMessageVisibilityCore changes the visibility timeout of a message.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_ChangeMessageVisibility.html
func (s *SQSService) changeMessageVisibilityCore(store sqsstore.SQSStoreInterface, in ChangeMessageVisibilityInput) error {
	if in.QueueURL == "" {
		return ErrMissingParameter
	}
	if in.ReceiptHandle == "" {
		return ErrMissingParameter
	}
	// VisibilityTimeout is @required ("Values range: 0 to 43200", model):
	// an omitted member is rejected rather than defaulted, since 0 carries
	// the destructive meaning "release for immediate redelivery".
	if !in.VisibilityTimeoutSet {
		return ErrMissingParameter
	}

	if err := store.ChangeMessageVisibility(in.QueueURL, in.ReceiptHandle, in.VisibilityTimeout); err != nil {
		return convertStoreError(err)
	}
	return nil
}

// changeMessageVisibilityBatchCore changes the visibility timeout for
// multiple messages in a single request.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/APIReference/API_ChangeMessageVisibilityBatch.html
//
// Two-pass design — all entries are validated before any mutation.
func (s *SQSService) changeMessageVisibilityBatchCore(store sqsstore.SQSStoreInterface, in ChangeMessageVisibilityBatchInput) (interface{}, error) {
	if in.QueueURL == "" {
		return nil, ErrMissingParameter
	}

	// Pass 1: Parse and validate all entries (no mutations).
	entries, err := parseChangeVisibilityBatchEntries(in.Parameters)
	if err != nil {
		return nil, err
	}

	ids := make([]string, len(entries))
	for i, e := range entries {
		ids[i] = e.id
	}
	if err := validateBatchEntryIDs(ids); err != nil {
		return nil, err
	}

	// QueueDoesNotExist is a request-level error: the whole request fails
	// before any entry is processed.
	queue, err := store.GetQueue(in.QueueURL)
	if err != nil {
		return nil, convertStoreError(err)
	}

	// An entry that omits VisibilityTimeout (Required: No on the entry
	// shape, model and API reference) selects the queue's VisibilityTimeout
	// attribute — the default the same member documents on ReceiveMessage —
	// never the destructive immediate-release 0.
	for i := range entries {
		if !entries[i].visibilityTimeoutSet {
			entries[i].visibilityTimeout = queue.VisibilityTimeout
		}
	}

	// Pass 2: Execute visibility changes.
	successEntries := make([]map[string]interface{}, 0)
	failedEntries := make([]map[string]interface{}, 0)

	for _, e := range entries {
		if err := store.ChangeMessageVisibility(in.QueueURL, e.receiptHandle, e.visibilityTimeout); err != nil {
			failedEntries = append(failedEntries, batchFailedEntry(e.id, err))
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

// rejectAttributeListValues rejects message-sending requests that carry
// message-attribute list values. The service model marks StringListValues
// and BinaryListValues as not implemented ("Not implemented. Reserved for
// future use."), and AWS rejects any request that includes them with
// UnsupportedOperation, so they must not be silently accepted and dropped.
func rejectAttributeListValues(params map[string]interface{}) error {
	// JSON protocol: SendMessage carries the attribute maps at the top level
	// and SendMessageBatch nests them inside each batch entry.
	attrMaps := make([]map[string]interface{}, 0, 8)
	for _, key := range []string{"MessageAttributes", "MessageSystemAttributes"} {
		if attrs, ok := params[key].(map[string]interface{}); ok {
			for _, val := range attrs {
				if attrMap, ok := val.(map[string]interface{}); ok {
					attrMaps = append(attrMaps, attrMap)
				}
			}
		}
	}
	if entries, ok := params["Entries"].([]interface{}); ok {
		for _, raw := range entries {
			entryMap, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}
			for _, key := range []string{"MessageAttributes", "MessageSystemAttributes"} {
				if attrs, ok := entryMap[key].(map[string]interface{}); ok {
					for _, val := range attrs {
						if attrMap, ok := val.(map[string]interface{}); ok {
							attrMaps = append(attrMaps, attrMap)
						}
					}
				}
			}
		}
	}
	for _, attrMap := range attrMaps {
		if list, ok := attrMap["StringListValues"].([]interface{}); ok && len(list) > 0 {
			return ErrUnsupportedOperation
		}
		if list, ok := attrMap["BinaryListValues"].([]interface{}); ok && len(list) > 0 {
			return ErrUnsupportedOperation
		}
	}
	// Query protocol: the flattened members arrive as indexed flat keys on
	// every message-sending operation (MessageAttribute.1.Value.StringListValue.1,
	// MessageSystemAttribute.1.Value.StringListValue.1, or the batch-entry
	// prefixed variants).
	for key := range params {
		lower := strings.ToLower(key)
		if strings.Contains(lower, ".value.stringlistvalue.") || strings.Contains(lower, ".value.binarylistvalue.") {
			return ErrUnsupportedOperation
		}
	}
	return nil
}

// messageEntrySize calculates the total byte size of a batch entry for the
// BatchRequestTooLong check. Includes message body and all attribute data.
// Delegates to the store helper so the size definition stays singular.
func messageEntrySize(body string, attrs map[string]*sqsstore.MessageAttributeValue) int {
	return sqsstore.MessageSize(body, attrs)
}

// validateBatchEntryIDs applies the request-level batch-entry rules every
// batch operation shares. Order of checks within this helper: the entry
// count first (TooManyEntriesInBatch), then each Id must be a valid
// batch-entry Id and distinct within the request. In the batch cores the
// helper runs AFTER the per-entry parse (a malformed entry rejects the
// request before the shared entry-ID rules see it) and before the queue
// resolves — the request-level precedence AWS's batch operations apply, as
// near as the sources pin it (the precise order among these rejections is
// not documented; the recorded order is the codebase's single convention).
func validateBatchEntryIDs(ids []string) error {
	if len(ids) > sqsstore.MaxBatchEntries {
		return ErrTooManyEntriesInBatch
	}
	seenIDs := make(map[string]bool, len(ids))
	for _, id := range ids {
		if err := sqsstore.ValidateBatchEntryId(id); err != nil {
			return ErrInvalidBatchEntryId
		}
		if seenIDs[id] {
			return ErrBatchEntryIdsNotDistinct
		}
		seenIDs[id] = true
	}
	return nil
}

// batchFailedEntry renders one Failed entry of a batch response through the
// shared store-error projection. The Message is the mapped single-operation
// error's text — the same failure's user-facing wording — rather than the
// store sentinel's internal phrasing; Code and SenderFault come from the
// same table row.
func batchFailedEntry(id string, err error) map[string]interface{} {
	code, senderFault := mapStoreErrorToBatchCode(err)
	message := err.Error()
	if row := lookupStoreErrorRow(err); row != nil {
		message = row.single.Error()
	}
	return map[string]interface{}{
		"Id":          id,
		"SenderFault": senderFault,
		"Code":        code,
		"Message":     message,
	}
}

// batchSendEntry holds a parsed and validated SendMessageBatch entry awaiting
// dispatch. All validation is completed before any SendMessage call.
type batchSendEntry struct {
	id          string
	message     *sqsstore.Message
	systemAttrs map[string]*sqsstore.MessageAttributeValue
}

// parseBatchSendEntries extracts and validates all SendMessageBatch entries
// from request parameters. Supports both JSON (SDK v2) and flattened query
// (CLI / SDK v1) formats. Returns an error if any entry is malformed or
// fails validation — no messages are sent until this returns successfully.
func parseBatchSendEntries(params map[string]interface{}) ([]*batchSendEntry, error) {
	var entries []*batchSendEntry
	var err error
	if jsonEntries, ok := params["Entries"].([]interface{}); ok && len(jsonEntries) > 0 {
		entries, err = parseBatchEntriesJSON(jsonEntries)
	} else {
		entries, err = parseBatchEntriesQuery(params)
	}
	if err != nil {
		return nil, err
	}

	ids := make([]string, len(entries))
	for i, e := range entries {
		ids[i] = e.id
	}
	if err := validateBatchEntryIDs(ids); err != nil {
		return nil, err
	}
	return entries, nil
}

func parseBatchEntriesJSON(jsonEntries []interface{}) ([]*batchSendEntry, error) {
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

		messageBody, _ := entryMap["MessageBody"].(string)
		// MessageBody is @required on SendMessageBatchRequestEntry: an
		// empty body is rejected exactly as the single-send path rejects
		// it, never silently sent.
		if messageBody == "" {
			return nil, ErrMissingParameter
		}
		// DelaySeconds is an Integer member: a present value that is not an
		// integer is a wire-type violation, not an omitted member. An omitted
		// member stays 0, which applies the queue's DelaySeconds attribute.
		delaySeconds := int32(0)
		if ds, present, derr := request.GetIntParamStrictCaseInsensitive(entryMap, "DelaySeconds"); present {
			if derr != nil {
				return nil, ErrSerializationException
			}
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
			msgAttrs, err := parseJSONMessageAttributes(attrs)
			if err != nil {
				return nil, err
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

		sysAttrs, err := parseSystemAttributes(entryMap, "")
		if err != nil {
			return nil, err
		}
		applyTraceHeader(message, sysAttrs)

		result = append(result, &batchSendEntry{
			id:          id,
			message:     message,
			systemAttrs: sysAttrs,
		})
	}

	return result, nil
}

func parseBatchEntriesQuery(params map[string]interface{}) ([]*batchSendEntry, error) {
	batchTotalSize := 0
	result := make([]*batchSendEntry, 0)

	for i := 1; ; i++ {
		prefix := "SendMessageBatchRequestEntry." + strconv.Itoa(i) + "."
		id := request.GetParamCaseInsensitive(params, prefix+"Id")
		if id == "" {
			if hasQueryEntryMembers(params, prefix) {
				// An entry-shaped key set without its required Id is a
				// malformed entry that rejects the request, not a list
				// terminator — the JSON arm rejects the same shape.
				return nil, ErrInvalidParameterValue
			}
			// Query-protocol entries are contiguous, so the first missing Id
			// ends the list; an over-count list is rejected by the shared
			// validateBatchEntryIDs pass in parseBatchSendEntries instead of
			// silently dropping entries beyond a fixed limit.
			break
		}

		messageBody := request.GetParamCaseInsensitive(params, prefix+"MessageBody")
		// MessageBody is @required on SendMessageBatchRequestEntry: an
		// empty body is rejected exactly as the single-send path rejects
		// it, never silently sent.
		if messageBody == "" {
			return nil, ErrMissingParameter
		}
		// DelaySeconds is an Integer member: a present value that is not an
		// integer is a wire-type violation, not an omitted member.
		delayVal, present, derr := request.GetIntParamStrictCaseInsensitive(params, prefix+"DelaySeconds")
		if derr != nil {
			return nil, ErrSerializationException
		}
		var delaySeconds int32
		if present {
			delaySeconds = int32(delayVal)
		}

		message := sqsstore.NewMessage(messageBody)
		message.DelaySeconds = delaySeconds
		message.MessageGroupID = request.GetParamCaseInsensitive(params, prefix+"MessageGroupId")
		message.MessageDeduplicationID = request.GetParamCaseInsensitive(params, prefix+"MessageDeduplicationId")

		msgAttrs, err := parseQueryMessageAttributes(params, prefix+"MessageAttribute.")
		if err != nil {
			return nil, err
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

		sysAttrs, err := parseSystemAttributes(params, prefix+"MessageSystemAttribute.")
		if err != nil {
			return nil, err
		}
		applyTraceHeader(message, sysAttrs)

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

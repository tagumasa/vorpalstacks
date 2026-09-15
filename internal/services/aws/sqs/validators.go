package sqs

import (
	"regexp"

	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

// ---------------------------------------------------------------------------
// Service-layer input validation (supplements store-layer validation to
// provide early rejection at the HTTP boundary — fail-closed defence-in-depth).
// Queue names and queue attributes validate through the store's shared
// exported rules (ValidateQueueName, ValidateQueueAttributes) at the Core.
// ---------------------------------------------------------------------------

const (
	maxReceiveAttemptIdLen = 128
)

var (
	// receiveAttemptIdRegex allows exactly the documented
	// ReceiveRequestAttemptId character set: "alphanumeric characters
	// (a-z, A-Z, 0-9) and punctuation
	// !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~" (AWS SQS API Reference) — every
	// printable non-space ASCII character.
	receiveAttemptIdRegex = regexp.MustCompile(`^[!-~]*$`)
)

// validateMaxNumberOfMessages rejects out-of-range values at the service layer.
// AWS SQS requires 1–10 messages per ReceiveMessage call.
func validateMaxNumberOfMessages(n int32) error {
	if n < sqsstore.MinMaxNumberOfMessages || n > sqsstore.MaxMaxNumberOfMessages {
		return ErrInvalidParameterValue
	}
	return nil
}

// validateReceiveRequestAttemptId checks length and character set for FIFO
// receive-request attempt IDs (max 128 chars, alphanumeric + punctuation).
func validateReceiveRequestAttemptId(s string) error {
	if len(s) > maxReceiveAttemptIdLen {
		return ErrInvalidParameterValue
	}
	if s != "" && !receiveAttemptIdRegex.MatchString(s) {
		return ErrInvalidParameterValue
	}
	return nil
}

// validatePermissionActionsCount enforces the AWS SQS limit of at most seven
// actions per AddPermission statement.
func validatePermissionActionsCount(actions []string) error {
	if len(actions) > sqsstore.MaxActionsPerStatement {
		return ErrOverLimit
	}
	return nil
}

// validateMessageMoveRate enforces the AWS SQS limit: MaxNumberOfMessagesPerSecond
// must be between 1 and 500. A value of 0 means "unset" and selects the
// system-optimised variable rate.
func validateMessageMoveRate(n int32) error {
	if n < 0 || n > sqsstore.MaxMessageMoveRate {
		return ErrInvalidParameterValue
	}
	return nil
}

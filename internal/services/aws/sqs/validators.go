package sqs

import (
	"regexp"
	"strings"

	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

// ---------------------------------------------------------------------------
// Service-layer input validation (supplements store-layer validation to
// provide early rejection at the HTTP boundary — fail-closed defence-in-depth).
// ---------------------------------------------------------------------------

const (
	maxReceiveAttemptIdLen = 128
	maxMessageMoveRate     = 500
)

var (
	// receiveAttemptIdRegex allows exactly the documented
	// ReceiveRequestAttemptId character set: "alphanumeric characters
	// (a-z, A-Z, 0-9) and punctuation
	// !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~" (AWS SQS API Reference) — every
	// printable non-space ASCII character.
	receiveAttemptIdRegex = regexp.MustCompile(`^[!-~]*$`)
)

// isValidQueueName checks whether a queue name conforms to the AWS SQS naming
// rules: 1–80 characters of [a-zA-Z0-9_-], with an optional .fifo suffix for
// FIFO queues.
func isValidQueueName(name string) bool {
	if len(name) == 0 || len(name) > sqsstore.MaxQueueNameLength {
		return false
	}
	if strings.HasSuffix(name, ".fifo") {
		prefix := name[:len(name)-5]
		if len(prefix) == 0 {
			return false
		}
		for _, c := range prefix {
			if !isAlphanumeric(c) && c != '-' && c != '_' {
				return false
			}
		}
		return true
	}
	for _, c := range name {
		if !isAlphanumeric(c) && c != '-' && c != '_' {
			return false
		}
	}
	return true
}

func isAlphanumeric(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

// validateMaxNumberOfMessages rejects out-of-range values at the service layer.
// AWS SQS requires 1–10 messages per ReceiveMessage call.
func validateMaxNumberOfMessages(n int32) error {
	if n < 1 || n > sqsstore.MaxMaxNumberOfMessages {
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
	if n < 0 || n > maxMessageMoveRate {
		return ErrInvalidParameterValue
	}
	return nil
}

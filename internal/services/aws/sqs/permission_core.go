package sqs

import (
	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

// ---------------------------------------------------------------------------
// Permission-operation Core — single validation + persistence path for
// queue sharing permissions. The HTTP handlers in permission_operations.go
// parse the account-ID and action lists (two wire formats) and delegate
// here.
// ---------------------------------------------------------------------------

// AddPermissionInput carries the parameters for adding a permission.
type AddPermissionInput struct {
	QueueURL      string
	Label         string
	AWSAccountIDs []string
	Actions       []string
}

// RemovePermissionInput carries the parameters for removing a permission.
type RemovePermissionInput struct {
	QueueURL string
	Label    string
}

// addPermissionCore adds a permission to an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_AddPermission.html
func (s *SQSService) addPermissionCore(store sqsstore.SQSStoreInterface, in AddPermissionInput) error {
	if in.QueueURL == "" {
		return ErrMissingParameter
	}
	if in.Label == "" {
		return ErrMissingParameter
	}
	// AWSAccountIds and Actions are @required on AddPermissionRequest: an
	// absent member rejects with the same MissingParameter family as the
	// shape's other required members, not the store's empty-list
	// InvalidParameterValue.
	if len(in.AWSAccountIDs) == 0 {
		return ErrMissingParameter
	}
	if len(in.Actions) == 0 {
		return ErrMissingParameter
	}

	if err := validatePermissionLabelFormat(in.Label); err != nil {
		return err
	}
	if err := validatePermissionActionsCount(in.Actions); err != nil {
		return err
	}

	if err := store.AddPermission(in.QueueURL, in.Label, in.AWSAccountIDs, in.Actions); err != nil {
		return convertStoreError(err)
	}

	return nil
}

// removePermissionCore removes a permission from an SQS queue.
// https://docs.aws.amazon.com/AWSSimpleQueueService/latest/API/API_RemovePermission.html
func (s *SQSService) removePermissionCore(store sqsstore.SQSStoreInterface, in RemovePermissionInput) error {
	if in.QueueURL == "" {
		return ErrMissingParameter
	}

	if in.Label == "" {
		return ErrMissingParameter
	}

	if err := store.RemovePermission(in.QueueURL, in.Label); err != nil {
		return convertStoreError(err)
	}

	return nil
}

// validatePermissionLabelFormat checks the label against the store's single
// permission-label rule (length and character set).
func validatePermissionLabelFormat(label string) error {
	return convertStoreError(sqsstore.ValidatePermissionLabel(label))
}

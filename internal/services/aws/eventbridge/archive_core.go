package eventbridge

import (
	"context"

	awserrors "vorpalstacks/internal/common/errors"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	"vorpalstacks/internal/utils/aws/arn"
)

// ---------------------------------------------------------------------------
// Input / Result structs (transport-agnostic)
// ---------------------------------------------------------------------------

// ArchiveMergeMembers carries the create/update archive merge members. The
// *Set flags distinguish an omitted member from an explicitly provided empty
// one so the merge semantics survive the transport boundary.
type ArchiveMergeMembers struct {
	DescriptionSet      bool
	Description         string
	EventPatternSet     bool
	EventPattern        string
	RetentionDaysSet    bool
	RetentionDays       int32
	KmsKeyIdentifierSet bool
	KmsKeyIdentifier    string
}

// CreateArchiveInput carries the parameters for CreateArchive.
type CreateArchiveInput struct {
	ArchiveName    string
	EventSourceArn string
	ArchiveMergeMembers
}

// UpdateArchiveInput carries the parameters for UpdateArchive.
type UpdateArchiveInput struct {
	ArchiveName string
	ArchiveMergeMembers
}

// ListArchivesInput carries the parameters for ListArchives.
type ListArchivesInput struct {
	NamePrefix     string
	State          string
	EventSourceArn string
	Limit          int32
	NextToken      string
}

// ---------------------------------------------------------------------------
// Core functions
// ---------------------------------------------------------------------------

// createArchiveCore validates input, checks the source event bus exists and
// creates the archive.
func (s *EventsService) createArchiveCore(ctx context.Context, store *eventsstore.EventsStore, input CreateArchiveInput) (*eventsstore.Archive, error) {
	if input.ArchiveName == "" {
		return nil, awserrors.NewValidationException("Archive name is required")
	}
	if !validateResourceName(input.ArchiveName, "archive") {
		return nil, awserrors.NewValidationException("Archive name must match the pattern and be 1-48 characters")
	}

	if input.EventSourceArn == "" {
		return nil, awserrors.NewValidationException("EventSourceArn is required")
	}

	eventBusName := arn.ExtractEventBusNameFromARN(input.EventSourceArn)

	// Check if event bus exists
	if _, err := store.GetEventBus(ctx, eventBusName); err != nil {
		return nil, mapStoreError(err, eventBusName)
	}

	archive := &eventsstore.Archive{
		Name:           input.ArchiveName,
		EventBusName:   eventBusName,
		EventSourceARN: input.EventSourceArn,
	}

	if input.DescriptionSet {
		if !validateDescription(input.Description) {
			return nil, errDescriptionTooLong()
		}
		archive.Description = input.Description
	}

	if input.EventPatternSet {
		if !validateEventPatternLength(input.EventPattern) {
			return nil, awserrors.NewValidationException("EventPattern must be at most 4096 characters")
		}
		if !isValidEventPattern(input.EventPattern) {
			return nil, awserrors.NewValidationException("EventPattern must be a valid JSON object")
		}
		archive.EventPattern = input.EventPattern
	}

	if input.RetentionDaysSet {
		archive.RetentionDays = input.RetentionDays
	}

	if input.KmsKeyIdentifierSet {
		if !validateKmsKeyIdentifier(input.KmsKeyIdentifier) {
			return nil, awserrors.NewValidationException("KmsKeyIdentifier must be a valid KMS ARN")
		}
		archive.KmsKeyIdentifier = input.KmsKeyIdentifier
	}

	if err := store.CreateArchive(ctx, archive); err != nil {
		return nil, mapStoreError(err, input.ArchiveName)
	}
	return archive, nil
}

// deleteArchiveCore validates input, deletes the archive and its stored
// events.
func (s *EventsService) deleteArchiveCore(ctx context.Context, store *eventsstore.EventsStore, name string) error {
	if name == "" {
		return awserrors.NewValidationException("Archive name is required")
	}

	if err := store.DeleteArchive(ctx, name); err != nil {
		return mapStoreError(err, name)
	}

	_ = store.DeleteArchiveEvents(ctx, name)

	return nil
}

// getArchiveCore validates input and fetches the archive.
func (s *EventsService) getArchiveCore(ctx context.Context, store *eventsstore.EventsStore, name string) (*eventsstore.Archive, error) {
	if name == "" {
		return nil, awserrors.NewValidationException("Archive name is required")
	}
	archive, err := store.GetArchive(ctx, name)
	if err != nil {
		return nil, mapStoreError(err, name)
	}
	return archive, nil
}

// updateArchiveCore validates input, merges the provided members onto the
// stored archive and persists the update.
func (s *EventsService) updateArchiveCore(ctx context.Context, store *eventsstore.EventsStore, input UpdateArchiveInput) (*eventsstore.Archive, error) {
	if input.ArchiveName == "" {
		return nil, awserrors.NewValidationException("Archive name is required")
	}

	archive, err := store.GetArchive(ctx, input.ArchiveName)
	if err != nil {
		return nil, mapStoreError(err, input.ArchiveName)
	}

	if input.DescriptionSet {
		if !validateDescription(input.Description) {
			return nil, errDescriptionTooLong()
		}
		archive.Description = input.Description
	}
	if input.EventPatternSet {
		if !validateEventPatternLength(input.EventPattern) {
			return nil, awserrors.NewValidationException("EventPattern must be at most 4096 characters")
		}
		if input.EventPattern != "" && !isValidEventPattern(input.EventPattern) {
			return nil, awserrors.NewValidationException("EventPattern must be a valid JSON object")
		}
		archive.EventPattern = input.EventPattern
	}
	if input.RetentionDaysSet {
		archive.RetentionDays = input.RetentionDays
	}
	if input.KmsKeyIdentifierSet {
		if !validateKmsKeyIdentifier(input.KmsKeyIdentifier) {
			return nil, awserrors.NewValidationException("KmsKeyIdentifier must be a valid KMS ARN")
		}
		archive.KmsKeyIdentifier = input.KmsKeyIdentifier
	}

	if err := store.UpdateArchive(ctx, archive); err != nil {
		return nil, err
	}
	return archive, nil
}

// listArchivesCore applies the documented limit window and lists the
// archives.
func (s *EventsService) listArchivesCore(ctx context.Context, store *eventsstore.EventsStore, input ListArchivesInput) (*eventsstore.ArchiveListResult, error) {
	limit := input.Limit
	if limit < 0 || limit > 100 {
		return nil, awserrors.NewValidationException("Limit must be between 0 and 100")
	}
	if limit == 0 {
		limit = 50
	}
	return store.ListArchives(ctx, input.NamePrefix, input.EventSourceArn, input.State, limit, input.NextToken)
}

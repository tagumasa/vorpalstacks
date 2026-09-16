package eventbridge

import (
	"context"
	"strings"

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
	// EventSourceArn targets the EventBusArn shape: a value that is not an
	// event-bus ARN is a pattern violation, not a missing-bus lookup.
	if !validateEventBusArn(input.EventSourceArn) {
		return nil, awserrors.NewValidationException("EventSourceArn must match the event bus ARN pattern")
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
		if err := validateEventPatternStructure(input.EventPattern); err != nil {
			return nil, err
		}
		archive.EventPattern = input.EventPattern
	}

	if input.RetentionDaysSet {
		if !validateRetentionDays(input.RetentionDays) {
			return nil, awserrors.NewValidationException("RetentionDays must be 0 or greater")
		}
		archive.RetentionDays = input.RetentionDays
	}

	if input.KmsKeyIdentifierSet {
		if !validateKmsKeyIdentifier(input.KmsKeyIdentifier) {
			return nil, awserrors.NewValidationException("KmsKeyIdentifier must be a key ARN, key ID, key alias or key alias ARN")
		}
		archive.KmsKeyIdentifier = input.KmsKeyIdentifier
	}

	if err := store.CreateArchive(ctx, archive); err != nil {
		return nil, mapStoreError(err, input.ArchiveName)
	}
	return archive, nil
}

// deleteArchiveCore validates input and deletes the archive. The stored
// events go with the record inside the store's locked delete — a caller-side
// sweep before the record delete could interleave with an ingress write and
// orphan the row under the dead archive name.
func (s *EventsService) deleteArchiveCore(ctx context.Context, store *eventsstore.EventsStore, name string) error {
	if name == "" {
		return awserrors.NewValidationException("Archive name is required")
	}

	return mapStoreError(store.DeleteArchive(ctx, name), name)
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
// stored archive and persists the update through the atomic record mutation
// (counter increments from the delivery path cannot be lost to the merge,
// and vice versa).
func (s *EventsService) updateArchiveCore(ctx context.Context, store *eventsstore.EventsStore, input UpdateArchiveInput) (*eventsstore.Archive, error) {
	if input.ArchiveName == "" {
		return nil, awserrors.NewValidationException("Archive name is required")
	}

	var updated *eventsstore.Archive
	if err := store.MutateArchive(ctx, input.ArchiveName, func(archive *eventsstore.Archive) error {
		if input.DescriptionSet {
			if !validateDescription(input.Description) {
				return errDescriptionTooLong()
			}
			archive.Description = input.Description
		}
		if input.EventPatternSet {
			if !validateEventPatternLength(input.EventPattern) {
				return awserrors.NewValidationException("EventPattern must be at most 4096 characters")
			}
			if err := validateEventPatternStructure(input.EventPattern); err != nil {
				return err
			}
			archive.EventPattern = input.EventPattern
		}
		if input.RetentionDaysSet {
			if !validateRetentionDays(input.RetentionDays) {
				return awserrors.NewValidationException("RetentionDays must be 0 or greater")
			}
			archive.RetentionDays = input.RetentionDays
		}
		if input.KmsKeyIdentifierSet {
			if !validateKmsKeyIdentifier(input.KmsKeyIdentifier) {
				return awserrors.NewValidationException("KmsKeyIdentifier must be a key ARN, key ID, key alias or key alias ARN")
			}
			archive.KmsKeyIdentifier = input.KmsKeyIdentifier
		}
		updated = archive
		return nil
	}); err != nil {
		return nil, mapStoreError(err, input.ArchiveName)
	}
	return updated, nil
}

// listArchivesCore validates the query and lists the archives: the limit
// window, the NamePrefix charset (both per the model), the State filter's
// enum membership and the EventSourceArn filter's EventBusArn form.
func (s *EventsService) listArchivesCore(ctx context.Context, store *eventsstore.EventsStore, input ListArchivesInput) (*eventsstore.ArchiveListResult, error) {
	if input.NamePrefix != "" {
		if err := validateListNamePrefix(input.NamePrefix, "archive"); err != nil {
			return nil, err
		}
	}
	if input.EventSourceArn != "" && !validateEventBusArn(input.EventSourceArn) {
		return nil, awserrors.NewValidationException("EventSourceArn must match the event bus ARN pattern")
	}
	if input.State != "" && !eventsstore.IsValidArchiveState(eventsstore.ArchiveState(input.State)) {
		return nil, awserrors.NewValidationException(
			"State must be a member of the ArchiveState enum: " + strings.Join(eventsstore.ArchiveStateVocabulary(), ", "))
	}
	limit, err := normaliseListLimit(input.Limit)
	if err != nil {
		return nil, err
	}
	return store.ListArchives(ctx, input.NamePrefix, input.EventSourceArn, input.State, limit, input.NextToken)
}

package appsync

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/core/resilience"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// The SourceApiAssociationStatus values this service writes. Every value
// must be a member of the SourceApiAssociationStatus enum of the AppSync
// API model; the membership is pinned by test.
const (
	assocStatusMergeScheduled    = "MERGE_SCHEDULED"
	assocStatusMergeInProgress   = "MERGE_IN_PROGRESS"
	assocStatusMergeSuccess      = "MERGE_SUCCESS"
	assocStatusDeletionScheduled = "DELETION_SCHEDULED"
	assocStatusDeletionFailed    = "DELETION_FAILED"
)

// associateSourceApiInput carries the parsed payload of the two association
// create operations (from the merged-API side and from the source-API side).
type associateSourceApiInput struct {
	MergedApiId string
	SourceApiId string
	Description string
	AssocConfig *appsyncstore.SourceApiAssociationConfig
}

// buildSourceApiAssociation validates both APIs and constructs the
// association record. fromMergedSide selects the addressing side of the two
// association-create operations; it only determines the association ARN's
// shape — merged-API-side ARNs nest under the merged API, source-API-side
// ARNs under the source API.
func (s *AppSyncService) buildSourceApiAssociation(store *appsyncstore.AppSyncStore, in associateSourceApiInput, fromMergedSide bool) (*appsyncstore.SourceApiAssociation, error) {
	if in.MergedApiId == "" {
		return nil, NewBadRequestException("mergedApiIdentifier is required")
	}
	if in.SourceApiId == "" {
		return nil, NewBadRequestException("sourceApiIdentifier is required")
	}

	if _, err := store.GetGraphqlApiById(in.SourceApiId); err != nil {
		return nil, mapStoreErrorE(err)
	}
	if _, err := store.GetGraphqlApiById(in.MergedApiId); err != nil {
		return nil, mapStoreErrorE(err)
	}

	assocID := uuid.New().String()
	arns := arnutil.NewARNBuilder(store.GetAccountID(), store.GetRegion())
	assoc := &appsyncstore.SourceApiAssociation{
		AssociationId:              assocID,
		MergedApiId:                in.MergedApiId,
		SourceApiId:                in.SourceApiId,
		MergedApiArn:               arns.AppSync().Api(in.MergedApiId),
		SourceApiArn:               arns.AppSync().Api(in.SourceApiId),
		SourceApiAssociationStatus: assocStatusMergeScheduled,
		Description:                in.Description,
		SourceApiAssociationConfig: in.AssocConfig,
	}
	if fromMergedSide {
		assoc.AssociationArn = arns.AppSync().SourceApiAssociation(in.MergedApiId, assocID)
	} else {
		assoc.AssociationArn = arns.AppSync().MergedApiAssociation(in.SourceApiId, assocID)
	}

	if err := store.CreateMergedApiAssociation(assoc); err != nil {
		return nil, mapStoreErrorE(err)
	}

	return assoc, nil
}

// associateSourceGraphqlApiCore validates the request and creates a source
// API association addressed from the merged API side.
func (s *AppSyncService) associateSourceGraphqlApiCore(store *appsyncstore.AppSyncStore, in associateSourceApiInput) (*appsyncstore.SourceApiAssociation, error) {
	return s.buildSourceApiAssociation(store, in, true)
}

// associateMergedGraphqlApiCore validates the request and creates a source
// API association addressed from the source API side.
func (s *AppSyncService) associateMergedGraphqlApiCore(store *appsyncstore.AppSyncStore, in associateSourceApiInput) (*appsyncstore.SourceApiAssociation, error) {
	return s.buildSourceApiAssociation(store, in, false)
}

// getSourceApiAssociationCore fetches one source API association of a merged
// API.
func (s *AppSyncService) getSourceApiAssociationCore(store *appsyncstore.AppSyncStore, mergedApiId, associationId string) (*appsyncstore.SourceApiAssociation, error) {
	if mergedApiId == "" {
		return nil, NewBadRequestException("mergedApiIdentifier is required")
	}
	if associationId == "" {
		return nil, NewBadRequestException("associationId is required")
	}

	assoc, err := store.GetMergedApiAssociation(mergedApiId, associationId)
	if err != nil {
		return nil, mapStoreErrorE(err)
	}

	return assoc, nil
}

// updateSourceApiAssociationCore applies an update to an existing source API
// association. Description and sourceApiAssociationConfig are presence-based:
// an omitted member keeps its stored value.
func (s *AppSyncService) updateSourceApiAssociationCore(store *appsyncstore.AppSyncStore, mergedApiId, associationId, description string, assocConfig *appsyncstore.SourceApiAssociationConfig) (*appsyncstore.SourceApiAssociation, error) {
	if mergedApiId == "" {
		return nil, NewBadRequestException("mergedApiIdentifier is required")
	}
	if associationId == "" {
		return nil, NewBadRequestException("associationId is required")
	}

	assoc, err := store.GetMergedApiAssociation(mergedApiId, associationId)
	if err != nil {
		return nil, mapStoreErrorE(err)
	}

	if description != "" {
		assoc.Description = description
	}
	if assocConfig != nil {
		assoc.SourceApiAssociationConfig = assocConfig
	}

	if err := store.UpdateMergedApiAssociation(assoc); err != nil {
		return nil, mapStoreErrorE(err)
	}

	return assoc, nil
}

// scheduleAssociationDeletion marks the association DELETION_SCHEDULED and
// performs the delayed deletion. On delete failure the status is persisted
// as DELETION_FAILED so the association does not silently remain scheduled
// forever. opLabel names the addressing side in log messages.
func (s *AppSyncService) scheduleAssociationDeletion(store *appsyncstore.AppSyncStore, assoc *appsyncstore.SourceApiAssociation, mergedApiId, associationId, opLabel string) (string, error) {
	assoc.SourceApiAssociationStatus = assocStatusDeletionScheduled
	if err := store.UpdateMergedApiAssociation(assoc); err != nil {
		return "", mapStoreErrorE(err)
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				resilience.LogPanic("appsync association async cleanup", r)
			}
		}()
		time.Sleep(5 * time.Second)
		if err := store.DeleteMergedApiAssociation(mergedApiId, associationId); err != nil {
			// The failure marker is written through a value copy so the
			// record handed to the caller is never mutated after the
			// response has been serialised.
			failed := *assoc
			failed.SourceApiAssociationStatus = assocStatusDeletionFailed
			if updateErr := store.UpdateMergedApiAssociation(&failed); updateErr != nil {
				logs.Warn("failed to persist the deletion-failed status marker",
					logs.String("mergedApiId", mergedApiId),
					logs.String("associationId", associationId),
					logs.Err(updateErr))
			}
			logs.Warn("async deletion of "+opLabel+" failed",
				logs.String("mergedApiId", mergedApiId),
				logs.String("associationId", associationId),
				logs.Err(err))
		}
	}()

	return assocStatusDeletionScheduled, nil
}

// disassociateSourceGraphqlApiCore schedules the deletion of a source API
// association addressed from the merged API side.
func (s *AppSyncService) disassociateSourceGraphqlApiCore(store *appsyncstore.AppSyncStore, mergedApiId, associationId string) (string, error) {
	if mergedApiId == "" {
		return "", NewBadRequestException("mergedApiIdentifier is required")
	}
	if associationId == "" {
		return "", NewBadRequestException("associationId is required")
	}

	assoc, err := store.GetMergedApiAssociation(mergedApiId, associationId)
	if err != nil {
		return "", mapStoreErrorE(err)
	}

	return s.scheduleAssociationDeletion(store, assoc, mergedApiId, associationId, "source API association")
}

// disassociateMergedGraphqlApiCore schedules the deletion of a source API
// association addressed from the source API side. The association is scoped
// by the required sourceApiIdentifier of the request.
func (s *AppSyncService) disassociateMergedGraphqlApiCore(store *appsyncstore.AppSyncStore, sourceApiId, associationId string) (string, error) {
	if sourceApiId == "" {
		return "", NewBadRequestException("sourceApiIdentifier is required")
	}
	if associationId == "" {
		return "", NewBadRequestException("associationId is required")
	}

	assoc, err := store.GetMergedApiAssociationById(associationId)
	if err != nil {
		return "", mapStoreErrorE(err)
	}
	if assoc.SourceApiId != sourceApiId {
		return "", NewNotFoundException(fmt.Sprintf("Source API association %s not found for source API %s", associationId, sourceApiId))
	}

	return s.scheduleAssociationDeletion(store, assoc, assoc.MergedApiId, associationId, "merged API association")
}

// startSchemaMergeCore transitions a source API association into
// MERGE_IN_PROGRESS and simulates the asynchronous merge to MERGE_SUCCESS.
func (s *AppSyncService) startSchemaMergeCore(store *appsyncstore.AppSyncStore, mergedApiId, associationId string) (string, error) {
	if mergedApiId == "" {
		return "", NewBadRequestException("mergedApiIdentifier is required")
	}
	if associationId == "" {
		return "", NewBadRequestException("associationId is required")
	}

	assoc, err := store.GetMergedApiAssociation(mergedApiId, associationId)
	if err != nil {
		return "", mapStoreErrorE(err)
	}

	assoc.SourceApiAssociationStatus = assocStatusMergeInProgress
	if err := store.UpdateMergedApiAssociation(assoc); err != nil {
		return "", mapStoreErrorE(err)
	}

	// Simulate async schema merge: transition MERGE_IN_PROGRESS →
	// MERGE_SUCCESS, writing through a value copy so the record handed to
	// the caller is never mutated after the response has been serialised.
	go func() {
		defer func() {
			if r := recover(); r != nil {
				resilience.LogPanic("appsync schema merge async", r)
			}
		}()
		time.Sleep(2 * time.Second)
		merged := *assoc
		merged.SourceApiAssociationStatus = assocStatusMergeSuccess
		now := time.Now().UTC()
		merged.LastSuccessfulMergeDate = &now
		if err := store.UpdateMergedApiAssociation(&merged); err != nil {
			logs.Warn("failed to persist merged API SUCCESS status",
				logs.String("mergedApiId", mergedApiId),
				logs.String("associationId", associationId),
				logs.Err(err))
		}
	}()

	return assocStatusMergeInProgress, nil
}

// listSourceApiAssociationsCore lists the source API associations of a
// GraphQL API acting as a merged API.
func (s *AppSyncService) listSourceApiAssociationsCore(store *appsyncstore.AppSyncStore, apiId string, maxResults int, nextToken string) ([]*appsyncstore.SourceApiAssociation, string, error) {
	if apiId == "" {
		return nil, "", NewBadRequestException("apiId is required")
	}

	if _, err := store.GetGraphqlApiById(apiId); err != nil {
		return nil, "", mapStoreErrorE(err)
	}

	opts, err := listOptionsFromParams(maxResults, nextToken)
	if err != nil {
		return nil, "", err
	}

	assocs, nextToken, err := store.ListMergedApiAssociations(apiId, opts)
	if err != nil {
		return nil, "", mapStoreErrorE(err)
	}

	return assocs, nextToken, nil
}

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

// associateSourceApiInput carries the parsed payload of the two association
// create operations (from the merged-API side and from the source-API side).
type associateSourceApiInput struct {
	MergedApiId string
	SourceApiId string
	Description string
	AssocConfig *appsyncstore.SourceApiAssociationConfig
}

// associateSourceGraphqlApiCore validates the request and creates a source
// API association addressed from the merged API side.
func (s *AppSyncService) associateSourceGraphqlApiCore(store *appsyncstore.AppSyncStore, in associateSourceApiInput) (*appsyncstore.SourceApiAssociation, error) {
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
	assoc := &appsyncstore.SourceApiAssociation{
		AssociationId:              assocID,
		MergedApiId:                in.MergedApiId,
		SourceApiId:                in.SourceApiId,
		MergedApiArn:               arnutil.NewARNBuilder(store.GetAccountID(), store.GetRegion()).AppSync().Api(in.MergedApiId),
		SourceApiArn:               arnutil.NewARNBuilder(store.GetAccountID(), store.GetRegion()).AppSync().Api(in.SourceApiId),
		AssociationArn:             arnutil.NewARNBuilder(store.GetAccountID(), store.GetRegion()).AppSync().SourceApiAssociation(in.MergedApiId, assocID),
		SourceApiAssociationStatus: "MERGE_SCHEDULED",
		Description:                in.Description,
		SourceApiAssociationConfig: in.AssocConfig,
	}

	if err := store.CreateMergedApiAssociation(assoc); err != nil {
		return nil, mapStoreErrorE(err)
	}

	return assoc, nil
}

// associateMergedGraphqlApiCore validates the request and creates a source
// API association addressed from the source API side.
func (s *AppSyncService) associateMergedGraphqlApiCore(store *appsyncstore.AppSyncStore, in associateSourceApiInput) (*appsyncstore.SourceApiAssociation, error) {
	if in.SourceApiId == "" {
		return nil, NewBadRequestException("sourceApiIdentifier is required")
	}
	if in.MergedApiId == "" {
		return nil, NewBadRequestException("mergedApiIdentifier is required")
	}

	if _, err := store.GetGraphqlApiById(in.SourceApiId); err != nil {
		return nil, mapStoreErrorE(err)
	}
	if _, err := store.GetGraphqlApiById(in.MergedApiId); err != nil {
		return nil, mapStoreErrorE(err)
	}

	assocID := uuid.New().String()
	assoc := &appsyncstore.SourceApiAssociation{
		AssociationId:              assocID,
		MergedApiId:                in.MergedApiId,
		SourceApiId:                in.SourceApiId,
		MergedApiArn:               arnutil.NewARNBuilder(store.GetAccountID(), store.GetRegion()).AppSync().Api(in.MergedApiId),
		SourceApiArn:               arnutil.NewARNBuilder(store.GetAccountID(), store.GetRegion()).AppSync().Api(in.SourceApiId),
		AssociationArn:             arnutil.NewARNBuilder(store.GetAccountID(), store.GetRegion()).AppSync().MergedApiAssociation(in.SourceApiId, assocID),
		SourceApiAssociationStatus: "MERGE_SCHEDULED",
		Description:                in.Description,
		SourceApiAssociationConfig: in.AssocConfig,
	}

	if err := store.CreateMergedApiAssociation(assoc); err != nil {
		return nil, mapStoreErrorE(err)
	}

	return assoc, nil
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

	assoc.SourceApiAssociationStatus = "DELETION_SCHEDULED"
	if err := store.UpdateMergedApiAssociation(assoc); err != nil {
		return "", mapStoreErrorE(err)
	}

	// Async deletion with persistent failure status.
	// On delete failure, the association status is updated to DELETION_FAILED
	// so it does not silently remain in DELETION_SCHEDULED forever.
	go func() {
		defer func() { resilience.RecoverPanic("appsync DisassociateSourceGraphqlApi async cleanup") }()
		time.Sleep(5 * time.Second)
		if err := store.DeleteMergedApiAssociation(mergedApiId, associationId); err != nil {
			assoc.SourceApiAssociationStatus = "DELETION_FAILED"
			if updateErr := store.UpdateMergedApiAssociation(assoc); updateErr != nil {
				logs.Warn("failed to persist DELETION_FAILED status",
					logs.String("mergedApiId", mergedApiId),
					logs.String("associationId", associationId),
					logs.Err(updateErr))
			}
			logs.Warn("async deletion of source API association failed",
				logs.String("mergedApiId", mergedApiId),
				logs.String("associationId", associationId),
				logs.Err(err))
		}
	}()

	return "DELETION_SCHEDULED", nil
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

	assoc.SourceApiAssociationStatus = "DELETION_SCHEDULED"
	if err := store.UpdateMergedApiAssociation(assoc); err != nil {
		return "", mapStoreErrorE(err)
	}

	mergedApiId := assoc.MergedApiId
	// Async deletion with persistent failure status.
	go func() {
		defer func() { resilience.RecoverPanic("appsync DisassociateMergedGraphqlApi async cleanup") }()
		time.Sleep(5 * time.Second)
		if err := store.DeleteMergedApiAssociation(mergedApiId, associationId); err != nil {
			assoc.SourceApiAssociationStatus = "DELETION_FAILED"
			if updateErr := store.UpdateMergedApiAssociation(assoc); updateErr != nil {
				logs.Warn("failed to persist DELETION_FAILED status",
					logs.String("mergedApiId", mergedApiId),
					logs.String("associationId", associationId),
					logs.Err(updateErr))
			}
			logs.Warn("async deletion of merged API association failed",
				logs.String("mergedApiId", mergedApiId),
				logs.String("associationId", associationId),
				logs.Err(err))
		}
	}()

	return "DELETION_SCHEDULED", nil
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

	assoc.SourceApiAssociationStatus = "MERGE_IN_PROGRESS"
	if err := store.UpdateMergedApiAssociation(assoc); err != nil {
		return "", mapStoreErrorE(err)
	}

	// Simulate async schema merge: transition MERGE_IN_PROGRESS → MERGE_SUCCESS.
	go func() {
		defer func() { resilience.RecoverPanic("appsync schema merge async") }()
		time.Sleep(2 * time.Second)
		assoc.SourceApiAssociationStatus = "MERGE_SUCCESS"
		now := time.Now().UTC()
		assoc.LastSuccessfulMergeDate = &now
		if err := store.UpdateMergedApiAssociation(assoc); err != nil {
			logs.Warn("failed to persist merged API SUCCESS status",
				logs.String("mergedApiId", mergedApiId),
				logs.String("associationId", associationId),
				logs.Err(err))
		}
	}()

	return "MERGE_IN_PROGRESS", nil
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

package cloudtrail

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/common/iam"
	tags "vorpalstacks/internal/common/tags"
	ctstore "vorpalstacks/internal/store/aws/cloudtrail"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input structs
// ---------------------------------------------------------------------------

// CreateEventDataStoreInput carries the create members for an event data
// store. Presence-tracked members distinguish "not provided" from an
// explicit value so the Core applies AWS defaults exactly.
type CreateEventDataStoreInput struct {
	Name                         string
	TerminationProtectionEnabled *bool
	MultiRegionEnabled           *bool
	OrganizationEnabled          *bool
	StartIngestion               *bool
	RetentionPeriodRaw           interface{}
	RetentionPeriodSet           bool
	KmsKeyId                     string
	BillingMode                  string
	AdvancedEventSelectorsRaw    interface{}
	AdvancedEventSelectorsSet    bool
	TagList                      []tags.Tag
	TagsRaw                      interface{}
	TagsSet                      bool
}

// UpdateEventDataStoreInput carries the update members for an event data
// store. Only provided members are applied; the rest keep their stored
// values.
type UpdateEventDataStoreInput struct {
	EventDataStore               string
	Name                         string
	TerminationProtectionEnabled *bool
	MultiRegionEnabled           *bool
	OrganizationEnabled          *bool
	RetentionPeriodRaw           interface{}
	RetentionPeriodSet           bool
	KmsKeyId                     string
	BillingMode                  string
	AdvancedEventSelectorsRaw    interface{}
	AdvancedEventSelectorsSet    bool
}

// EventDataStoreIDInput carries the event data store identifier.
type EventDataStoreIDInput struct {
	EventDataStore string
}

// ListEventDataStoresInput carries pagination parameters for listing event
// data stores.
type ListEventDataStoresInput struct {
	NextToken  string
	MaxResults int
}

// EnableFederationInput carries the members for EnableFederation.
type EnableFederationInput struct {
	EventDataStore    string
	FederationRoleArn string
	IAMValidator      *iam.IAMValidator
}

// DisableFederationInput carries the members for DisableFederation.
type DisableFederationInput struct {
	EventDataStore string
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// createEventDataStoreCore is the single entry point for CreateEventDataStore.
func (s *CloudTrailService) createEventDataStoreCore(store ctstore.CloudTrailStoreInterface, in CreateEventDataStoreInput) (map[string]interface{}, error) {
	name := in.Name
	if name == "" {
		return nil, newInvalidParameterException("Name is required")
	}
	if err := validateEventDataStoreName(name); err != nil {
		return nil, err
	}

	eds := ctstore.NewEventDataStore(name, store.GetAccountID(), store.GetRegion())

	if in.TerminationProtectionEnabled != nil {
		eds.TerminationProtectionEnabled = *in.TerminationProtectionEnabled
	}
	if in.MultiRegionEnabled != nil {
		eds.MultiRegionEnabled = *in.MultiRegionEnabled
	}
	if in.OrganizationEnabled != nil {
		eds.OrganizationEnabled = *in.OrganizationEnabled
	}
	if in.StartIngestion != nil {
		eds.IngestionEnabled = *in.StartIngestion
	}
	if rp, err := extractRetentionPeriod(in.RetentionPeriodRaw, in.RetentionPeriodSet); err != nil {
		return nil, err
	} else if rp > 0 {
		// FIXED_RETENTION_PRICING caps the retention period at 2557 days
		// (CreateEventDataStore RetentionPeriod).
		if in.BillingMode == "FIXED_RETENTION_PRICING" && rp > ctstore.MaxEventDataStoreRetentionDaysFixedPricing {
			return nil, newInvalidParameterException(
				fmt.Sprintf("RetentionPeriod must be between %d and %d days for FIXED_RETENTION_PRICING",
					ctstore.MinEventDataStoreRetentionDays, ctstore.MaxEventDataStoreRetentionDaysFixedPricing))
		}
		eds.RetentionPeriod = rp
	}
	if in.KmsKeyId != "" {
		if err := validateEventDataStoreKMSKeyID(in.KmsKeyId); err != nil {
			return nil, err
		}
		eds.KMSKeyID = in.KmsKeyId
	}
	if in.BillingMode != "" {
		if err := validateBillingMode(in.BillingMode); err != nil {
			return nil, err
		}
		eds.BillingMode = in.BillingMode
	}

	if in.AdvancedEventSelectorsSet {
		selectors, parseErr := parseAdvancedEventSelectors(in.AdvancedEventSelectorsRaw)
		if parseErr != nil {
			return nil, parseErr
		}
		// The shared validator carries the 500-value quota with the
		// selector count — the quota table binds trails and event data
		// stores alike.
		if err := validateAdvancedEventSelectors(selectors); err != nil {
			return nil, err
		}
		eds.AdvancedEventSelectors = selectors
	}

	// Validate and apply tags BEFORE creation to ensure atomicity.
	if in.TagsSet {
		if err := validateCloudTrailTags(in.TagList); err != nil {
			return nil, err
		}
		applyTags(&eds.Tags, in.TagsRaw)
	}

	// The event data store quota is enforced inside CreateEventDataStore
	// (the store counts every lifecycle stage under its mutex —
	// PENDING_DELETION included: "This includes event data stores in any
	// lifecycle stage", Quotas in AWS CloudTrail); the model declares
	// EventDataStoreMaxLimitExceededException on CreateEventDataStore.

	created, err := store.CreateEventDataStore(eds)
	if err != nil {
		if errors.Is(err, ctstore.ErrEventDataStoreQuotaExceeded) {
			return nil, newEventDataStoreMaxLimitExceededException(
				fmt.Sprintf("The maximum number of event data stores per Region (%d) has been reached", ctstore.MaxEventDataStoresPerRegion))
		}
		return nil, s.mapStoreError(err)
	}

	return formatEventDataStoreFor(created, edsProfileCreate), nil
}

// resolveEventDataStore resolves an event data store by ID or ARN. A
// selector that carries the ARN prefix but is not a well-formed event data
// store ARN answers the declared EventDataStoreARNInvalidException; store
// errors propagate for the caller's own mapping.
func (s *CloudTrailService) resolveEventDataStore(store ctstore.CloudTrailStoreInterface, idOrARN string) (*ctstore.EventDataStore, error) {
	if strings.HasPrefix(idOrARN, "arn:") && !isCloudTrailResourceARN(idOrARN, "eventdatastore/") {
		return nil, newEventDataStoreARNInvalidException(
			fmt.Sprintf("The specified event data store ARN is not valid: %s", idOrARN))
	}
	return store.GetEventDataStore(idOrARN)
}

// getEventDataStoreCore is the single entry point for GetEventDataStore.
func (s *CloudTrailService) getEventDataStoreCore(store ctstore.CloudTrailStoreInterface, in EventDataStoreIDInput) (map[string]interface{}, error) {
	if in.EventDataStore == "" {
		return nil, newInvalidParameterException("EventDataStore is required")
	}

	eds, err := s.resolveEventDataStore(store, in.EventDataStore)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return formatEventDataStoreFor(eds, edsProfileGet), nil
}

// listEventDataStoresCore is the single entry point for ListEventDataStores.
func (s *CloudTrailService) listEventDataStoresCore(store ctstore.CloudTrailStoreInterface, in ListEventDataStoresInput) (map[string]interface{}, error) {
	opts := storecommon.ListOptions{MaxItems: ctstore.DefaultListEventDataStoresResults}
	if in.NextToken != "" {
		opts.Marker = in.NextToken
	}
	if in.MaxResults > 0 {
		if in.MaxResults > ctstore.MaxListEventDataStoresResults {
			return nil, newInvalidMaxResultsException(
				fmt.Sprintf("MaxResults exceeds the maximum of %d", ctstore.MaxListEventDataStoresResults))
		}
		opts.MaxItems = in.MaxResults
	}

	result, err := store.ListEventDataStores(opts)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	items := make([]interface{}, 0, len(result.Items))
	for _, eds := range result.Items {
		items = append(items, formatEventDataStoreFor(eds, edsProfileList))
	}

	resp := map[string]interface{}{
		"EventDataStores": items,
	}
	if result.NextMarker != "" {
		resp["NextToken"] = result.NextMarker
	}

	return resp, nil
}

// updateEventDataStoreCore is the single entry point for UpdateEventDataStore.
// Input validation runs on the request values; the provided members are then
// applied to the stored record under the store mutex, so concurrent updates
// and lifecycle transitions cannot lose writes. The mutation closure also
// enforces the store's lifecycle preconditions: an inactive
// (PENDING_DELETION) store and a store with an import in progress are
// rejected, and the KMS key and the EXTENDABLE billing mode are immutable
// once set.
func (s *CloudTrailService) updateEventDataStoreCore(store ctstore.CloudTrailStoreInterface, in UpdateEventDataStoreInput) (map[string]interface{}, error) {
	if in.EventDataStore == "" {
		return nil, newInvalidParameterException("EventDataStore is required")
	}

	if in.Name != "" {
		if err := validateEventDataStoreName(in.Name); err != nil {
			return nil, err
		}
	}
	if in.BillingMode != "" {
		if err := validateBillingMode(in.BillingMode); err != nil {
			return nil, err
		}
	}
	if in.KmsKeyId != "" {
		if err := validateEventDataStoreKMSKeyID(in.KmsKeyId); err != nil {
			return nil, err
		}
	}
	retentionPeriod, err := extractRetentionPeriod(in.RetentionPeriodRaw, in.RetentionPeriodSet)
	if err != nil {
		return nil, err
	}

	// "Other parameters are optional, but at least one optional parameter
	// must be specified, or CloudTrail throws an error"
	// (UpdateEventDataStore).
	if in.Name == "" && in.TerminationProtectionEnabled == nil && in.MultiRegionEnabled == nil &&
		in.OrganizationEnabled == nil && !in.RetentionPeriodSet && in.KmsKeyId == "" &&
		in.BillingMode == "" && !in.AdvancedEventSelectorsSet {
		return nil, newInvalidParameterException(
			"At least one parameter to update is required")
	}

	var advancedSelectors []ctstore.AdvancedEventSelector
	if in.AdvancedEventSelectorsSet {
		parsed, parseErr := parseAdvancedEventSelectors(in.AdvancedEventSelectorsRaw)
		if parseErr != nil {
			return nil, parseErr
		}
		// The shared validator carries the 500-value quota with the
		// selector count — the quota table binds trails and event data
		// stores alike.
		if err := validateAdvancedEventSelectors(parsed); err != nil {
			return nil, err
		}
		advancedSelectors = parsed
	}

	eds, err := store.MutateEventDataStore(in.EventDataStore, func(eds *ctstore.EventDataStore) error {
		// The update preconditions run under the same mutex as the write so
		// a lifecycle transition cannot interleave between check and apply.
		if eds.Status == "PENDING_DELETION" {
			return newInactiveEventDataStoreException(
				"The event data store is inactive")
		}
		ongoing, listErr := storeHasOngoingImport(store, eds.EventDataStoreARN)
		if listErr != nil {
			return listErr
		}
		if ongoing {
			return newEventDataStoreHasOngoingImportException(
				"Cannot update an event data store with an import in progress")
		}
		// "After you associate an event data store with a KMS key, the KMS
		// key cannot be removed or changed" (CreateEventDataStore /
		// UpdateEventDataStore KmsKeyId). Re-asserting the stored value is
		// accepted as a no-op.
		if in.KmsKeyId != "" && eds.KMSKeyID != "" && eds.KMSKeyID != in.KmsKeyId {
			return newOperationNotPermittedException(
				"The KMS key of an event data store cannot be changed once associated")
		}
		// "You can't change the billing mode from
		// EXTENDABLE_RETENTION_PRICING to FIXED_RETENTION_PRICING"
		// (UpdateEventDataStore BillingMode).
		if in.BillingMode == "FIXED_RETENTION_PRICING" && eds.BillingMode != "" &&
			eds.BillingMode != in.BillingMode {
			return newOperationNotPermittedException(
				"The billing mode cannot change from EXTENDABLE_RETENTION_PRICING to FIXED_RETENTION_PRICING")
		}
		// The retention upper bound depends on the effective billing mode:
		// FIXED_RETENTION_PRICING caps the period at 2557 days.
		if retentionPeriod > 0 {
			effectiveBilling := eds.BillingMode
			if in.BillingMode != "" {
				effectiveBilling = in.BillingMode
			}
			maxRetention := int32(ctstore.MaxEventDataStoreRetentionDays)
			if effectiveBilling == "FIXED_RETENTION_PRICING" {
				maxRetention = ctstore.MaxEventDataStoreRetentionDaysFixedPricing
			}
			if retentionPeriod > maxRetention {
				return newInvalidParameterException(
					fmt.Sprintf("RetentionPeriod must be between %d and %d days for this billing mode",
						ctstore.MinEventDataStoreRetentionDays, maxRetention))
			}
		}

		if in.Name != "" {
			eds.Name = in.Name
		}
		if in.TerminationProtectionEnabled != nil {
			eds.TerminationProtectionEnabled = *in.TerminationProtectionEnabled
		}
		if in.MultiRegionEnabled != nil {
			eds.MultiRegionEnabled = *in.MultiRegionEnabled
		}
		if in.OrganizationEnabled != nil {
			eds.OrganizationEnabled = *in.OrganizationEnabled
		}
		if retentionPeriod > 0 {
			eds.RetentionPeriod = retentionPeriod
		}
		if in.KmsKeyId != "" {
			eds.KMSKeyID = in.KmsKeyId
		}
		if in.BillingMode != "" {
			eds.BillingMode = in.BillingMode
		}
		if in.AdvancedEventSelectorsSet {
			eds.AdvancedEventSelectors = advancedSelectors
		}
		return nil
	})
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return formatEventDataStoreFor(eds, edsProfileUpdate), nil
}

// deleteEventDataStoreCore is the single entry point for DeleteEventDataStore.
// It enforces the termination-protection, federation, ongoing-import and
// channel-association preconditions and the PENDING_DELETION flip as one
// atomic step under the store mutex: nothing can interleave between the
// checks and the status write. A precondition-check failure aborts the
// delete instead of silently skipping the guard.
func (s *CloudTrailService) deleteEventDataStoreCore(store ctstore.CloudTrailStoreInterface, in EventDataStoreIDInput) (map[string]interface{}, error) {
	if in.EventDataStore == "" {
		return nil, newInvalidParameterException("EventDataStore is required")
	}

	_, err := store.MutateEventDataStore(in.EventDataStore, func(eds *ctstore.EventDataStore) error {
		if eds.TerminationProtectionEnabled {
			return newEventDataStoreTerminationProtectedException(
				"The event data store cannot be deleted because termination protection is enabled")
		}

		if eds.FederationStatus == "ENABLED" {
			return newEventDataStoreFederationEnabledException(
				"Cannot delete event data store with federation enabled. Disable federation first.")
		}

		// Check for ongoing imports referencing this EDS.  The destinations
		// list stores the ARN values provided by the SDK, so we compare against
		// the EDS ARN (not the short ID).
		ongoing, listErr := storeHasOngoingImport(store, eds.EventDataStoreARN)
		if listErr != nil {
			return listErr
		}
		if ongoing {
			return newEventDataStoreHasOngoingImportException(
				"Cannot delete event data store with an ongoing import. Stop the import first.")
		}

		// Check for channels referencing this EDS as a destination. The
		// association is Destination.Location (the EDS ARN) on every
		// EVENT_DATA_STORE destination; every channel page is drained so
		// the guard is complete. A listing failure aborts the delete.
		channelMarker := ""
		for {
			channelResult, listErr := store.ListChannels(storecommon.ListOptions{
				MaxItems: ctstore.MaxListChannelsResults,
				Marker:   channelMarker,
			})
			if listErr != nil {
				return listErr
			}
			for _, ch := range channelResult.Items {
				for _, dest := range ch.Destinations {
					if dest.Type != ctstore.DestinationTypeEventDataStore {
						continue
					}
					if ctstore.ExtractEventDataStoreID(dest.Location) == eds.EventDataStoreID {
						return newChannelExistsForEDSException(
							"Cannot delete event data store because a channel is associated with it")
					}
				}
			}
			if channelResult.NextMarker == "" {
				break
			}
			channelMarker = channelResult.NextMarker
		}

		eds.Status = "PENDING_DELETION"
		now := time.Now().UTC()
		eds.DeletedTimestamp = &now
		return nil
	})
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{}, nil
}

// startEventDataStoreIngestionCore is the single entry point for
// StartEventDataStoreIngestion.
func (s *CloudTrailService) startEventDataStoreIngestionCore(store ctstore.CloudTrailStoreInterface, in EventDataStoreIDInput) (map[string]interface{}, error) {
	return s.setEventDataStoreIngestion(store, in, true)
}

// stopEventDataStoreIngestionCore is the single entry point for
// StopEventDataStoreIngestion.
func (s *CloudTrailService) stopEventDataStoreIngestionCore(store ctstore.CloudTrailStoreInterface, in EventDataStoreIDInput) (map[string]interface{}, error) {
	return s.setEventDataStoreIngestion(store, in, false)
}

// ingestionToggleableCategories is the category set both ingestion
// operations require: "the eventCategory must be Management, Data,
// NetworkActivity, or ConfigurationItem" (StartEventDataStoreIngestion and
// StopEventDataStoreIngestion).
var ingestionToggleableCategories = map[string]bool{
	"Management":        true,
	"Data":              true,
	"NetworkActivity":   true,
	"ConfigurationItem": true,
}

// edsIngestionCategories returns the eventCategory values the store's
// advanced selectors pin — the platform's record of an event data store's
// category. A store carrying no eventCategory selector ingests management
// events (the documented default the materialised selectors express).
func edsIngestionCategories(eds *ctstore.EventDataStore) []string {
	var categories []string
	seen := map[string]bool{}
	add := func(v string) {
		if v != "" && !seen[v] {
			seen[v] = true
			categories = append(categories, v)
		}
	}
	for _, sel := range eds.AdvancedEventSelectors {
		for _, fs := range sel.FieldSelectors {
			if strings.EqualFold(fs.Field, "eventCategory") {
				for _, v := range fs.Equals {
					add(v)
				}
			}
		}
	}
	if len(categories) == 0 {
		add("Management")
	}
	return categories
}

// setEventDataStoreIngestion applies the ingestion lifecycle shared by the
// start and stop operations: "To stop ingestion, the event data store
// Status must be ENABLED" and "To start ingestion, the event data store
// Status must be STOPPED_INGESTION" (StartEventDataStoreIngestion /
// StopEventDataStoreIngestion) — the stop lands the store in
// STOPPED_INGESTION, the start returns it to ENABLED. The category
// precondition, the status precondition and the transition run under the
// store mutex as one step.
func (s *CloudTrailService) setEventDataStoreIngestion(store ctstore.CloudTrailStoreInterface, in EventDataStoreIDInput, enabled bool) (map[string]interface{}, error) {
	if in.EventDataStore == "" {
		return nil, newInvalidParameterException("EventDataStore is required")
	}

	_, err := store.MutateEventDataStore(in.EventDataStore, func(eds *ctstore.EventDataStore) error {
		if err := validateEventDataStoreStatus(eds.Status); err != nil {
			return err
		}
		for _, category := range edsIngestionCategories(eds) {
			if !ingestionToggleableCategories[category] {
				return newInvalidEventDataStoreCategoryException(fmt.Sprintf(
					"Ingestion cannot be started or stopped on an event data store with eventCategory %s", category))
			}
		}
		want, next := "ENABLED", "STOPPED_INGESTION"
		verb := "stop"
		if enabled {
			want, next, verb = "STOPPED_INGESTION", "ENABLED", "start"
		}
		if eds.Status != want {
			return newInvalidEventDataStoreStatusException(fmt.Sprintf(
				"Event data store must be in %s state to %s ingestion", want, verb))
		}
		eds.Status = next
		eds.IngestionEnabled = enabled
		return nil
	})
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{}, nil
}

// restoreEventDataStoreCore is the single entry point for
// RestoreEventDataStore, which only restores a PENDING_DELETION store.
func (s *CloudTrailService) restoreEventDataStoreCore(store ctstore.CloudTrailStoreInterface, in EventDataStoreIDInput) (map[string]interface{}, error) {
	if in.EventDataStore == "" {
		return nil, newInvalidParameterException("EventDataStore is required")
	}

	eds, err := store.RestoreEventDataStore(in.EventDataStore)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return formatEventDataStoreFor(eds, edsProfileCreate), nil
}

// enableFederationCore is the single entry point for EnableFederation.
func (s *CloudTrailService) enableFederationCore(ctx context.Context, store ctstore.CloudTrailStoreInterface, in EnableFederationInput) (map[string]interface{}, error) {
	if in.EventDataStore == "" {
		return nil, newInvalidParameterException("EventDataStore is required")
	}

	if in.FederationRoleArn == "" {
		return nil, newInvalidParameterException("FederationRoleArn is required")
	}

	if in.IAMValidator != nil {
		if err := in.IAMValidator.ValidateRoleForService(ctx, in.FederationRoleArn, iam.ServicePrincipalCloudTrail); err != nil {
			return nil, err
		}
	}

	eds, err := store.MutateEventDataStore(in.EventDataStore, func(eds *ctstore.EventDataStore) error {
		// A PENDING_DELETION store is inactive — the same precondition the
		// update path enforces through the operation's declared
		// InactiveEventDataStoreException.
		if eds.Status == "PENDING_DELETION" {
			return newInactiveEventDataStoreException("The event data store is inactive")
		}
		eds.FederationStatus = "ENABLED"
		eds.FederationRoleARN = in.FederationRoleArn
		return nil
	})
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	// EnableFederationResponse carries the ARN, role and status alone.
	return map[string]interface{}{
		"EventDataStoreArn": eds.EventDataStoreARN,
		"FederationRoleArn": eds.FederationRoleARN,
		"FederationStatus":  eds.FederationStatus,
	}, nil
}

// disableFederationCore is the single entry point for DisableFederation.
func (s *CloudTrailService) disableFederationCore(store ctstore.CloudTrailStoreInterface, in DisableFederationInput) (map[string]interface{}, error) {
	if in.EventDataStore == "" {
		return nil, newInvalidParameterException("EventDataStore is required")
	}

	eds, err := store.MutateEventDataStore(in.EventDataStore, func(eds *ctstore.EventDataStore) error {
		// The inactive precondition mirrors the enable direction.
		if eds.Status == "PENDING_DELETION" {
			return newInactiveEventDataStoreException("The event data store is inactive")
		}
		eds.FederationStatus = "DISABLED"
		eds.FederationRoleARN = ""
		return nil
	})
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	// DisableFederationResponse carries the ARN and status alone.
	return map[string]interface{}{
		"EventDataStoreArn": eds.EventDataStoreARN,
		"FederationStatus":  eds.FederationStatus,
	}, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// formatEventDataStore builds the response map for an event data store.
// edsResponseProfile selects which model response shape a formatted event
// data store targets: the four shapes share a base member set but differ
// in the billing, tags, federation and partition-key members they carry.
type edsResponseProfile struct {
	billingMode   bool
	kmsKeyID      bool
	tagsList      bool
	federation    bool
	partitionKeys bool
}

// The per-operation profiles, each mirroring one model response shape:
// create and restore add billing, KMS key and tags to the base; get adds
// billing, KMS key, federation and partition keys; update drops the
// partition keys; the list item shape (EventDataStore) carries the base
// alone.
var (
	edsProfileCreate = edsResponseProfile{billingMode: true, kmsKeyID: true, tagsList: true}
	edsProfileGet    = edsResponseProfile{billingMode: true, kmsKeyID: true, federation: true, partitionKeys: true}
	edsProfileUpdate = edsResponseProfile{billingMode: true, kmsKeyID: true, federation: true}
	edsProfileList   = edsResponseProfile{}
)

// formatEventDataStoreFor renders an event data store for one response
// profile. Members the platform has no data for (PartitionKeys, until
// organization event data stores exist) stay absent, never fabricated.
func formatEventDataStoreFor(eds *ctstore.EventDataStore, profile edsResponseProfile) map[string]interface{} {
	resp := map[string]interface{}{
		"EventDataStoreArn":            eds.EventDataStoreARN,
		"Name":                         eds.Name,
		"TerminationProtectionEnabled": eds.TerminationProtectionEnabled,
		"Status":                       eds.Status,
		"MultiRegionEnabled":           eds.MultiRegionEnabled,
		"OrganizationEnabled":          eds.OrganizationEnabled,
		"RetentionPeriod":              eds.RetentionPeriod,
		"CreatedTimestamp":             eds.CreatedTimestamp,
		"UpdatedTimestamp":             eds.UpdatedTimestamp,
	}
	if profile.billingMode && eds.BillingMode != "" {
		resp["BillingMode"] = eds.BillingMode
	}
	if profile.kmsKeyID && eds.KMSKeyID != "" {
		resp["KmsKeyId"] = eds.KMSKeyID
	}
	if profile.federation && eds.FederationStatus != "" {
		resp["FederationStatus"] = eds.FederationStatus
	}
	if profile.federation && eds.FederationRoleARN != "" {
		resp["FederationRoleArn"] = eds.FederationRoleARN
	}
	if len(eds.AdvancedEventSelectors) > 0 {
		resp["AdvancedEventSelectors"] = formatAdvancedEventSelectors(eds.AdvancedEventSelectors)
	}
	if profile.tagsList && len(eds.Tags) > 0 {
		resp["TagsList"] = formatTagsList(eds.Tags)
	}
	return resp
}

// extractRetentionPeriod validates the RetentionPeriod wire value. Returns
// (0, nil) when the member was not provided.
func extractRetentionPeriod(raw interface{}, provided bool) (int32, error) {
	if !provided {
		return 0, nil
	}
	rp := int32(0)
	switch val := raw.(type) {
	case float64:
		rp = int32(val)
	case int:
		rp = int32(val)
	case int32:
		rp = val
	}
	if rp < ctstore.MinEventDataStoreRetentionDays || rp > ctstore.MaxEventDataStoreRetentionDays {
		return 0, newInvalidParameterException(
			fmt.Sprintf("RetentionPeriod must be between %d and %d days",
				ctstore.MinEventDataStoreRetentionDays, ctstore.MaxEventDataStoreRetentionDays))
	}
	return rp, nil
}

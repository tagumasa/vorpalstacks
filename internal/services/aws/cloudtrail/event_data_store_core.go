package cloudtrail

import (
	"context"
	"encoding/json"

	awserrors "vorpalstacks/internal/common/errors"
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
	IngestionEnabled             *bool
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
	IngestionEnabled             *bool
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
		return nil, awserrors.NewAWSError("InvalidEventDataStoreCategory", "Name is required", 400)
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
	if in.IngestionEnabled != nil {
		eds.IngestionEnabled = *in.IngestionEnabled
	}
	if in.StartIngestion != nil {
		eds.IngestionEnabled = *in.StartIngestion
	}
	if rp, err := extractRetentionPeriod(in.RetentionPeriodRaw, in.RetentionPeriodSet); err != nil {
		return nil, err
	} else if rp > 0 {
		eds.RetentionPeriod = rp
	}
	if in.KmsKeyId != "" {
		eds.KMSKeyID = in.KmsKeyId
	}
	if in.BillingMode != "" {
		if err := validateBillingMode(in.BillingMode); err != nil {
			return nil, err
		}
		eds.BillingMode = in.BillingMode
	}

	if in.AdvancedEventSelectorsSet {
		eds.AdvancedEventSelectors = parseAdvancedEventSelectors(in.AdvancedEventSelectorsRaw)
	}

	// Validate and apply tags BEFORE creation to ensure atomicity.
	if in.TagsSet {
		if err := validateCloudTrailTags(in.TagList); err != nil {
			return nil, err
		}
		applyEventDataStoreTags(eds, in.TagsRaw)
	}

	created, err := store.CreateEventDataStore(eds)
	if err != nil {
		if err == ctstore.ErrEventDataStoreAlreadyExists {
			return nil, awserrors.NewAWSError("EventDataStoreAlreadyExistsException", "Event data store already exists", 409)
		}
		return nil, err
	}

	return formatEventDataStore(created), nil
}

// getEventDataStoreCore is the single entry point for GetEventDataStore.
func (s *CloudTrailService) getEventDataStoreCore(store ctstore.CloudTrailStoreInterface, in EventDataStoreIDInput) (map[string]interface{}, error) {
	if in.EventDataStore == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "EventDataStore is required", 400)
	}

	eds, err := store.GetEventDataStore(in.EventDataStore)
	if err != nil {
		if err == ctstore.ErrEventDataStoreNotFound {
			return nil, awserrors.NewAWSError("EventDataStoreNotFoundException", "Event data store not found", 404)
		}
		return nil, err
	}

	return formatEventDataStore(eds), nil
}

// listEventDataStoresCore is the single entry point for ListEventDataStores.
func (s *CloudTrailService) listEventDataStoresCore(store ctstore.CloudTrailStoreInterface, in ListEventDataStoresInput) (map[string]interface{}, error) {
	opts := storecommon.ListOptions{MaxItems: 100}
	if in.NextToken != "" {
		opts.Marker = in.NextToken
	}
	if in.MaxResults > 0 {
		opts.MaxItems = in.MaxResults
	}

	result, err := store.ListEventDataStores(opts)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(result.Items))
	for _, eds := range result.Items {
		items = append(items, formatEventDataStore(eds))
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
func (s *CloudTrailService) updateEventDataStoreCore(store ctstore.CloudTrailStoreInterface, in UpdateEventDataStoreInput) (map[string]interface{}, error) {
	if in.EventDataStore == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "EventDataStore is required", 400)
	}

	eds, err := store.GetEventDataStore(in.EventDataStore)
	if err != nil {
		if err == ctstore.ErrEventDataStoreNotFound {
			return nil, awserrors.NewAWSError("EventDataStoreNotFoundException", "Event data store not found", 404)
		}
		return nil, err
	}

	if in.Name != "" {
		if err := validateEventDataStoreName(in.Name); err != nil {
			return nil, err
		}
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
	if in.IngestionEnabled != nil {
		eds.IngestionEnabled = *in.IngestionEnabled
	}
	if rp, err := extractRetentionPeriod(in.RetentionPeriodRaw, in.RetentionPeriodSet); err != nil {
		return nil, err
	} else if rp > 0 {
		eds.RetentionPeriod = rp
	}
	if in.KmsKeyId != "" {
		eds.KMSKeyID = in.KmsKeyId
	}
	if in.BillingMode != "" {
		if err := validateBillingMode(in.BillingMode); err != nil {
			return nil, err
		}
		eds.BillingMode = in.BillingMode
	}
	if in.AdvancedEventSelectorsSet {
		eds.AdvancedEventSelectors = parseAdvancedEventSelectors(in.AdvancedEventSelectorsRaw)
	}

	if err := store.UpdateEventDataStore(eds); err != nil {
		return nil, err
	}

	return formatEventDataStore(eds), nil
}

// deleteEventDataStoreCore is the single entry point for DeleteEventDataStore.
// It enforces the termination-protection, federation, ongoing-import and
// channel-association preconditions before deleting.
func (s *CloudTrailService) deleteEventDataStoreCore(store ctstore.CloudTrailStoreInterface, in EventDataStoreIDInput) (map[string]interface{}, error) {
	if in.EventDataStore == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "EventDataStore is required", 400)
	}

	eds, err := store.GetEventDataStore(in.EventDataStore)
	if err != nil {
		if err == ctstore.ErrEventDataStoreNotFound {
			return nil, awserrors.NewAWSError("EventDataStoreNotFoundException", "Event data store not found", 404)
		}
		return nil, err
	}

	if eds.TerminationProtectionEnabled {
		return nil, awserrors.NewAWSError("EventDataStoreTerminationProtectedException",
			"The event data store cannot be deleted because termination protection is enabled", 400)
	}

	if eds.FederationStatus == "ENABLED" {
		return nil, awserrors.NewAWSError("EventDataStoreFederationEnabledException",
			"Cannot delete event data store with federation enabled. Disable federation first.", 400)
	}

	// Check for ongoing imports referencing this EDS.  The destinations
	// list stores the ARN values provided by the SDK, so we compare against
	// the EDS ARN (not the short ID).
	imports, err := store.ListImports(storecommon.ListOptions{MaxItems: 1}, eds.EventDataStoreARN, "IN_PROGRESS")
	if err == nil && len(imports.Items) > 0 {
		return nil, awserrors.NewAWSError("EventDataStoreHasOngoingImportException",
			"Cannot delete event data store with an ongoing import. Stop the import first.", 400)
	}

	// Check for channels referencing this EDS as a destination.
	channelResult, err := store.ListChannels(storecommon.ListOptions{MaxItems: 1000})
	if err == nil {
		for _, ch := range channelResult.Items {
			for _, dest := range ch.Destinations {
				if dest.EDSARN == eds.EventDataStoreARN || dest.EDSARN == eds.EventDataStoreID {
					return nil, awserrors.NewAWSError("ChannelExistsForEDSException",
						"Cannot delete event data store because a channel is associated with it", 400)
				}
			}
		}
	}

	if err := store.DeleteEventDataStore(eds.EventDataStoreID); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// startEventDataStoreIngestionCore is the single entry point for
// StartEventDataStoreIngestion.
func (s *CloudTrailService) startEventDataStoreIngestionCore(store ctstore.CloudTrailStoreInterface, in EventDataStoreIDInput) (map[string]interface{}, error) {
	if in.EventDataStore == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "EventDataStore is required", 400)
	}

	eds, err := store.GetEventDataStore(in.EventDataStore)
	if err != nil {
		if err == ctstore.ErrEventDataStoreNotFound {
			return nil, awserrors.NewAWSError("EventDataStoreNotFoundException", "Event data store not found", 404)
		}
		return nil, err
	}

	if err := validateEventDataStoreStatus(eds.Status); err != nil {
		return nil, err
	}
	if eds.Status != "ENABLED" {
		return nil, awserrors.NewAWSError("InvalidEventDataStoreStatusException",
			"Event data store must be in ENABLED state to start ingestion", 400)
	}

	eds.IngestionEnabled = true
	if err := store.UpdateEventDataStore(eds); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// stopEventDataStoreIngestionCore is the single entry point for
// StopEventDataStoreIngestion.
func (s *CloudTrailService) stopEventDataStoreIngestionCore(store ctstore.CloudTrailStoreInterface, in EventDataStoreIDInput) (map[string]interface{}, error) {
	if in.EventDataStore == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "EventDataStore is required", 400)
	}

	eds, err := store.GetEventDataStore(in.EventDataStore)
	if err != nil {
		if err == ctstore.ErrEventDataStoreNotFound {
			return nil, awserrors.NewAWSError("EventDataStoreNotFoundException", "Event data store not found", 404)
		}
		return nil, err
	}

	if err := validateEventDataStoreStatus(eds.Status); err != nil {
		return nil, err
	}
	if eds.Status != "ENABLED" {
		return nil, awserrors.NewAWSError("InvalidEventDataStoreStatusException",
			"Event data store must be in ENABLED state to stop ingestion", 400)
	}

	eds.IngestionEnabled = false
	if err := store.UpdateEventDataStore(eds); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// restoreEventDataStoreCore is the single entry point for
// RestoreEventDataStore, which only restores a PENDING_DELETION store.
func (s *CloudTrailService) restoreEventDataStoreCore(store ctstore.CloudTrailStoreInterface, in EventDataStoreIDInput) (map[string]interface{}, error) {
	if in.EventDataStore == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "EventDataStore is required", 400)
	}

	eds, err := store.RestoreEventDataStore(in.EventDataStore)
	if err != nil {
		if err == ctstore.ErrEventDataStoreNotFound {
			return nil, awserrors.NewAWSError("EventDataStoreNotFoundException", "Event data store not found", 404)
		}
		if err == ctstore.ErrEventDataStoreNotPendingDeletion {
			return nil, awserrors.NewAWSError("OperationNotPermittedException",
				"Event data store is not in PENDING_DELETION state", 400)
		}
		return nil, err
	}

	return formatEventDataStore(eds), nil
}

// enableFederationCore is the single entry point for EnableFederation.
func (s *CloudTrailService) enableFederationCore(ctx context.Context, store ctstore.CloudTrailStoreInterface, in EnableFederationInput) (map[string]interface{}, error) {
	if in.EventDataStore == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "EventDataStore is required", 400)
	}

	if in.FederationRoleArn == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "FederationRoleArn is required", 400)
	}

	if in.IAMValidator != nil {
		if err := in.IAMValidator.ValidateRoleForService(ctx, in.FederationRoleArn, iam.ServicePrincipalCloudTrail); err != nil {
			return nil, err
		}
	}

	eds, err := store.GetEventDataStore(in.EventDataStore)
	if err != nil {
		if err == ctstore.ErrEventDataStoreNotFound {
			return nil, awserrors.NewAWSError("EventDataStoreNotFoundException", "Event data store not found", 404)
		}
		return nil, err
	}

	eds.FederationStatus = "ENABLED"
	eds.FederationRoleARN = in.FederationRoleArn
	if err := store.UpdateEventDataStore(eds); err != nil {
		return nil, err
	}

	return formatEventDataStore(eds), nil
}

// disableFederationCore is the single entry point for DisableFederation.
func (s *CloudTrailService) disableFederationCore(store ctstore.CloudTrailStoreInterface, in DisableFederationInput) (map[string]interface{}, error) {
	if in.EventDataStore == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "EventDataStore is required", 400)
	}

	eds, err := store.GetEventDataStore(in.EventDataStore)
	if err != nil {
		if err == ctstore.ErrEventDataStoreNotFound {
			return nil, awserrors.NewAWSError("EventDataStoreNotFoundException", "Event data store not found", 404)
		}
		return nil, err
	}

	eds.FederationStatus = "DISABLED"
	eds.FederationRoleARN = ""
	if err := store.UpdateEventDataStore(eds); err != nil {
		return nil, err
	}

	return formatEventDataStore(eds), nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// formatEventDataStore builds the response map for an event data store.
func formatEventDataStore(eds *ctstore.EventDataStore) map[string]interface{} {
	resp := map[string]interface{}{
		"EventDataStoreArn":            eds.EventDataStoreARN,
		"Name":                         eds.Name,
		"TerminationProtectionEnabled": eds.TerminationProtectionEnabled,
		"Status":                       eds.Status,
		"MultiRegionEnabled":           eds.MultiRegionEnabled,
		"OrganizationEnabled":          eds.OrganizationEnabled,
		"RetentionPeriod":              eds.RetentionPeriod,
		"IngestionEnabled":             eds.IngestionEnabled,
		"CreatedTimestamp":             eds.CreatedTimestamp,
		"UpdatedTimestamp":             eds.UpdatedTimestamp,
	}
	if eds.BillingMode != "" {
		resp["BillingMode"] = eds.BillingMode
	}
	if eds.KMSKeyID != "" {
		resp["KmsKeyId"] = eds.KMSKeyID
	}
	if eds.FederationStatus != "" {
		resp["FederationStatus"] = eds.FederationStatus
	}
	if eds.FederationRoleARN != "" {
		resp["FederationRoleArn"] = eds.FederationRoleARN
	}
	if len(eds.AdvancedEventSelectors) > 0 {
		selectors := make([]interface{}, 0, len(eds.AdvancedEventSelectors))
		for _, sel := range eds.AdvancedEventSelectors {
			selMap := map[string]interface{}{}
			if sel.Name != "" {
				selMap["Name"] = sel.Name
			}
			fields := make([]interface{}, 0, len(sel.FieldSelectors))
			for _, fs := range sel.FieldSelectors {
				fm := map[string]interface{}{"Field": fs.Field}
				if len(fs.Equals) > 0 {
					fm["Equals"] = fs.Equals
				}
				if len(fs.StartsWith) > 0 {
					fm["StartsWith"] = fs.StartsWith
				}
				if len(fs.EndsWith) > 0 {
					fm["EndsWith"] = fs.EndsWith
				}
				if len(fs.NotEquals) > 0 {
					fm["NotEquals"] = fs.NotEquals
				}
				if len(fs.NotStartsWith) > 0 {
					fm["NotStartsWith"] = fs.NotStartsWith
				}
				if len(fs.NotEndsWith) > 0 {
					fm["NotEndsWith"] = fs.NotEndsWith
				}
				fields = append(fields, fm)
			}
			selMap["FieldSelectors"] = fields
			selectors = append(selectors, selMap)
		}
		resp["AdvancedEventSelectors"] = selectors
	}
	if len(eds.Tags) > 0 {
		tagsList := make([]interface{}, 0, len(eds.Tags))
		for k, v := range eds.Tags {
			tagsList = append(tagsList, map[string]interface{}{
				"Key":   k,
				"Value": v,
			})
		}
		resp["TagsList"] = tagsList
	}
	return resp
}

// parseAdvancedEventSelectors parses the advanced event selectors from the
// raw wire value.
func parseAdvancedEventSelectors(raw interface{}) []ctstore.AdvancedEventSelector {
	var rawList []interface{}
	switch v := raw.(type) {
	case []interface{}:
		rawList = v
	case string:
		if err := json.Unmarshal([]byte(v), &rawList); err != nil {
			return nil
		}
	default:
		return nil
	}

	result := make([]ctstore.AdvancedEventSelector, 0, len(rawList))
	for _, item := range rawList {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		sel := ctstore.AdvancedEventSelector{}
		if name, ok := m["Name"].(string); ok {
			sel.Name = name
		}
		if fieldsRaw, ok := m["FieldSelectors"].([]interface{}); ok {
			for _, fRaw := range fieldsRaw {
				fm, ok := fRaw.(map[string]interface{})
				if !ok {
					continue
				}
				fs := ctstore.AdvancedFieldSelector{}
				if f, ok := fm["Field"].(string); ok {
					fs.Field = f
				}
				fs.Equals = toStringSlice(fm["Equals"])
				fs.StartsWith = toStringSlice(fm["StartsWith"])
				fs.EndsWith = toStringSlice(fm["EndsWith"])
				fs.NotEquals = toStringSlice(fm["NotEquals"])
				fs.NotStartsWith = toStringSlice(fm["NotStartsWith"])
				fs.NotEndsWith = toStringSlice(fm["NotEndsWith"])
				sel.FieldSelectors = append(sel.FieldSelectors, fs)
			}
		}
		result = append(result, sel)
	}
	return result
}

// toStringSlice converts an interface to a string slice.
func toStringSlice(v interface{}) []string {
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case []interface{}:
		result := make([]string, 0, len(val))
		for _, item := range val {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	case []string:
		return val
	default:
		return nil
	}
}

// applyEventDataStoreTags parses tags from the raw interface and applies them
// to the event data store.
func applyEventDataStoreTags(eds *ctstore.EventDataStore, raw interface{}) {
	if eds.Tags == nil {
		eds.Tags = make(map[string]string)
	}
	var tagsList []interface{}
	switch v := raw.(type) {
	case []interface{}:
		tagsList = v
	case string:
		if err := json.Unmarshal([]byte(v), &tagsList); err != nil {
			return
		}
	default:
		return
	}
	for _, item := range tagsList {
		if m, ok := item.(map[string]interface{}); ok {
			key, _ := m["Key"].(string)
			val, _ := m["Value"].(string)
			if key != "" {
				eds.Tags[key] = val
			}
		}
	}
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
	if rp < 7 || rp > 3653 {
		return 0, awserrors.NewAWSError("InvalidParameterException",
			"RetentionPeriod must be between 7 and 3653 days", 400)
	}
	return rp, nil
}

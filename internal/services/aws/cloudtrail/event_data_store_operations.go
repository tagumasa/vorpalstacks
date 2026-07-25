package cloudtrail

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/iam"
	"vorpalstacks/internal/common/request"
	tags "vorpalstacks/internal/common/tags"
	ctstore "vorpalstacks/internal/store/aws/cloudtrail"
	storecommon "vorpalstacks/internal/store/aws/common"
)

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

// CreateEventDataStore creates a new CloudTrail event data store.
func (s *CloudTrailService) CreateEventDataStore(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewAWSError("InvalidEventDataStoreCategory", "Name is required", 400)
	}
	if err := validateEventDataStoreName(name); err != nil {
		return nil, err
	}

	eds := ctstore.NewEventDataStore(name, store.GetAccountID(), store.GetRegion())

	if v, ok := req.Parameters["TerminationProtectionEnabled"]; ok {
		eds.TerminationProtectionEnabled = parseBool(v)
	}
	if v, ok := req.Parameters["MultiRegionEnabled"]; ok {
		eds.MultiRegionEnabled = parseBool(v)
	}
	if v, ok := req.Parameters["OrganizationEnabled"]; ok {
		eds.OrganizationEnabled = parseBool(v)
	}
	if v, ok := req.Parameters["IngestionEnabled"]; ok {
		eds.IngestionEnabled = parseBool(v)
	}
	if v, ok := req.Parameters["StartIngestion"]; ok {
		eds.IngestionEnabled = parseBool(v)
	}
	if rp, err := extractRetentionPeriod(req.Parameters); err != nil {
		return nil, err
	} else if rp > 0 {
		eds.RetentionPeriod = rp
	}
	if v := request.GetStringParam(req.Parameters, "KmsKeyId"); v != "" {
		eds.KMSKeyID = v
	}
	if v := request.GetStringParam(req.Parameters, "BillingMode"); v != "" {
		if v != "EXTENDABLE_RETENTION_PRICING" && v != "FIXED_RETENTION_PRICING" {
			return nil, awserrors.NewAWSError("InvalidParameterException",
				"BillingMode must be EXTENDABLE_RETENTION_PRICING or FIXED_RETENTION_PRICING", 400)
		}
		eds.BillingMode = v
	}

	if aesRaw, ok := req.Parameters["AdvancedEventSelectors"]; ok {
		eds.AdvancedEventSelectors = parseAdvancedEventSelectors(aesRaw)
	}

	// Validate and apply tags BEFORE creation to ensure atomicity.
	if tagsRaw, ok := req.Parameters["TagsList"]; ok {
		tagList := tags.ParseTags(req.Parameters, "TagsList")
		if err := validateCloudTrailTags(tagList); err != nil {
			return nil, err
		}
		applyEventDataStoreTags(eds, tagsRaw)
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

// GetEventDataStore retrieves details about the specified event data store.
func (s *CloudTrailService) GetEventDataStore(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	id := request.GetStringParam(req.Parameters, "EventDataStore")
	if id == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "EventDataStore is required", 400)
	}

	eds, err := store.GetEventDataStore(id)
	if err != nil {
		if err == ctstore.ErrEventDataStoreNotFound {
			return nil, awserrors.NewAWSError("EventDataStoreNotFoundException", "Event data store not found", 404)
		}
		return nil, err
	}

	return formatEventDataStore(eds), nil
}

// ListEventDataStores returns event data stores with pagination.
func (s *CloudTrailService) ListEventDataStores(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	opts := storecommon.ListOptions{MaxItems: 100}
	if nextToken := req.GetParam("NextToken"); nextToken != "" {
		opts.Marker = nextToken
	}
	if maxResults := request.GetIntParam(req.Parameters, "MaxResults"); maxResults > 0 {
		opts.MaxItems = maxResults
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

// UpdateEventDataStore updates the specified event data store.
func (s *CloudTrailService) UpdateEventDataStore(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	id := request.GetStringParam(req.Parameters, "EventDataStore")
	if id == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "EventDataStore is required", 400)
	}

	eds, err := store.GetEventDataStore(id)
	if err != nil {
		if err == ctstore.ErrEventDataStoreNotFound {
			return nil, awserrors.NewAWSError("EventDataStoreNotFoundException", "Event data store not found", 404)
		}
		return nil, err
	}

	if name := request.GetStringParam(req.Parameters, "Name"); name != "" {
		if err := validateEventDataStoreName(name); err != nil {
			return nil, err
		}
		eds.Name = name
	}
	if v, ok := req.Parameters["TerminationProtectionEnabled"]; ok {
		eds.TerminationProtectionEnabled = parseBool(v)
	}
	if v, ok := req.Parameters["MultiRegionEnabled"]; ok {
		eds.MultiRegionEnabled = parseBool(v)
	}
	if v, ok := req.Parameters["OrganizationEnabled"]; ok {
		eds.OrganizationEnabled = parseBool(v)
	}
	if v, ok := req.Parameters["IngestionEnabled"]; ok {
		eds.IngestionEnabled = parseBool(v)
	}
	if rp, err := extractRetentionPeriod(req.Parameters); err != nil {
		return nil, err
	} else if rp > 0 {
		eds.RetentionPeriod = rp
	}
	if v := request.GetStringParam(req.Parameters, "KmsKeyId"); v != "" {
		eds.KMSKeyID = v
	}
	if v := request.GetStringParam(req.Parameters, "BillingMode"); v != "" {
		if v != "EXTENDABLE_RETENTION_PRICING" && v != "FIXED_RETENTION_PRICING" {
			return nil, awserrors.NewAWSError("InvalidParameterException",
				"BillingMode must be EXTENDABLE_RETENTION_PRICING or FIXED_RETENTION_PRICING", 400)
		}
		eds.BillingMode = v
	}
	if aesRaw, ok := req.Parameters["AdvancedEventSelectors"]; ok {
		eds.AdvancedEventSelectors = parseAdvancedEventSelectors(aesRaw)
	}

	if err := store.UpdateEventDataStore(eds); err != nil {
		return nil, err
	}

	return formatEventDataStore(eds), nil
}

// DeleteEventDataStore deletes the specified event data store.
func (s *CloudTrailService) DeleteEventDataStore(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	id := request.GetStringParam(req.Parameters, "EventDataStore")
	if id == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "EventDataStore is required", 400)
	}

	eds, err := store.GetEventDataStore(id)
	if err != nil {
		if err == ctstore.ErrEventDataStoreNotFound {
			return nil, awserrors.NewAWSError("EventDataStoreNotFoundException", "Event data store not found", 404)
		}
		return nil, err
	}

	if eds.TerminationProtectionEnabled {
		return nil, awserrors.NewAWSError("EventDataStoreHasTerminationProtectionEnabled", "Cannot delete event data store with termination protection enabled", 400)
	}

	if eds.FederationStatus == "ENABLED" {
		return nil, awserrors.NewAWSError("OperationNotPermittedException",
			"Cannot delete event data store with federation enabled. Disable federation first.", 400)
	}

	if err := store.DeleteEventDataStore(eds.EventDataStoreID); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// parseBool handles bool values from request parameters which may be bool or
// json.Number.
func parseBool(v interface{}) bool {
	switch val := v.(type) {
	case bool:
		return val
	case string:
		return strings.ToLower(val) == "true"
	default:
		return false
	}
}

// parseAdvancedEventSelectors parses the advanced event selectors from the
// request parameters.
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

// edsNamePattern matches the Smithy model constraint for EventDataStoreName:
// ^[a-zA-Z0-9._\-]+$, length 3-128.
var edsNamePattern = regexp.MustCompile(`^[a-zA-Z0-9._\-]+$`)

// validateEventDataStoreName validates the EDS name against AWS spec
// (Smithy model: length 3-128, pattern ^[a-zA-Z0-9._\-]+$).
func validateEventDataStoreName(name string) error {
	if len(name) < 3 || len(name) > 128 {
		return awserrors.NewAWSError("InvalidParameterException",
			"Event data store name must be between 3 and 128 characters", 400)
	}
	if !edsNamePattern.MatchString(name) {
		return awserrors.NewAWSError("InvalidParameterException",
			"Event data store name contains invalid characters", 400)
	}
	return nil
}

// StartEventDataStoreIngestion enables ingestion on the specified event data
// store.
func (s *CloudTrailService) StartEventDataStoreIngestion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	id := request.GetStringParam(req.Parameters, "EventDataStore")
	if id == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "EventDataStore is required", 400)
	}

	eds, err := store.GetEventDataStore(id)
	if err != nil {
		if err == ctstore.ErrEventDataStoreNotFound {
			return nil, awserrors.NewAWSError("EventDataStoreNotFoundException", "Event data store not found", 404)
		}
		return nil, err
	}

	if eds.Status != "ENABLED" {
		return nil, awserrors.NewAWSError("OperationNotPermittedException",
			"Event data store must be in ENABLED state", 400)
	}

	eds.IngestionEnabled = true
	if err := store.UpdateEventDataStore(eds); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// StopEventDataStoreIngestion disables ingestion on the specified event data
// store.
func (s *CloudTrailService) StopEventDataStoreIngestion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	id := request.GetStringParam(req.Parameters, "EventDataStore")
	if id == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "EventDataStore is required", 400)
	}

	eds, err := store.GetEventDataStore(id)
	if err != nil {
		if err == ctstore.ErrEventDataStoreNotFound {
			return nil, awserrors.NewAWSError("EventDataStoreNotFoundException", "Event data store not found", 404)
		}
		return nil, err
	}

	if eds.Status != "ENABLED" {
		return nil, awserrors.NewAWSError("OperationNotPermittedException",
			"Event data store must be in ENABLED state", 400)
	}

	eds.IngestionEnabled = false
	if err := store.UpdateEventDataStore(eds); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// RestoreEventDataStore restores a PENDING_DELETION event data store.
func (s *CloudTrailService) RestoreEventDataStore(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	id := request.GetStringParam(req.Parameters, "EventDataStore")
	if id == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "EventDataStore is required", 400)
	}

	eds, err := store.RestoreEventDataStore(id)
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

// EnableFederation enables Lake query federation on the specified event data
// store.
func (s *CloudTrailService) EnableFederation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	id := request.GetStringParam(req.Parameters, "EventDataStore")
	if id == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "EventDataStore is required", 400)
	}

	federationRoleARN := request.GetStringParam(req.Parameters, "FederationRoleArn")
	if federationRoleARN == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "FederationRoleArn is required", 400)
	}

	validator := reqCtx.GetIAMValidator()
	if err := validator.ValidateRoleForService(ctx, federationRoleARN, iam.ServicePrincipalCloudTrail); err != nil {
		return nil, err
	}

	eds, err := store.GetEventDataStore(id)
	if err != nil {
		if err == ctstore.ErrEventDataStoreNotFound {
			return nil, awserrors.NewAWSError("EventDataStoreNotFoundException", "Event data store not found", 404)
		}
		return nil, err
	}

	eds.FederationStatus = "ENABLED"
	eds.FederationRoleARN = federationRoleARN
	if err := store.UpdateEventDataStore(eds); err != nil {
		return nil, err
	}

	return formatEventDataStore(eds), nil
}

// DisableFederation disables Lake query federation on the specified event data
// store.
func (s *CloudTrailService) DisableFederation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	id := request.GetStringParam(req.Parameters, "EventDataStore")
	if id == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "EventDataStore is required", 400)
	}

	eds, err := store.GetEventDataStore(id)
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

// extractRetentionPeriod extracts and validates the RetentionPeriod parameter
// from request parameters. Returns (0, nil) if not provided.
func extractRetentionPeriod(params map[string]interface{}) (int32, error) {
	v, ok := params["RetentionPeriod"]
	if !ok {
		return 0, nil
	}
	rp := int32(0)
	switch val := v.(type) {
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

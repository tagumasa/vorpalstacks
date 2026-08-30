package athena

import (
	"fmt"
	"sort"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
	athenastore "vorpalstacks/internal/store/aws/athena"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// --- Service-layer DTOs (no store-package types exposed) ---

// WorkGroupCreateInput is the transport-neutral input for creating a workgroup.
// Both the HTTP handler and the admin gRPC-Web handler populate this struct
// from their respective request formats and pass it to createWorkGroupCore.
type WorkGroupCreateInput struct {
	Name        string
	Description string
	Config      *WorkGroupConfigInput
	Tags        map[string]string
}

// WorkGroupConfigInput holds the subset of WorkGroup configuration fields that
// can be set at creation time. Mirrors athenastore.WorkGroupConfiguration but
// uses only primitive types so callers need not import the store package.
type WorkGroupConfigInput struct {
	OutputLocation                       string
	EnforceConfig                        bool
	PublishMetrics                       bool
	BytesScannedCutoff                   *int64
	RequesterPaysEnabled                 bool
	EngineVersionSelected                string
	EngineVersionEffective               string
	AdditionalConfiguration              string
	ExecutionRole                        string
	CustomerContentEncryptionKmsKey      string
	EnableMinimumEncryptionConfiguration bool
}

func (in *WorkGroupConfigInput) bytesScannedCutoffValue() int64 {
	if in.BytesScannedCutoff != nil {
		return *in.BytesScannedCutoff
	}
	return 0
}

// WorkGroupOut is the transport-neutral representation of a workgroup summary.
type WorkGroupOut struct {
	Name         string
	State        string
	Description  string
	CreationTime time.Time
}

// WorkGroupListResult holds a paginated list of workgroups.
type WorkGroupListResult struct {
	Items      []WorkGroupOut
	NextMarker string
}

// --- Core functions ---

// createWorkGroupCore validates the input, constructs the store object, and
// persists it. It is called by both the HTTP handler and the admin handler,
// ensuring consistent validation and error handling.
func createWorkGroupCore(stores *athenaStores, input WorkGroupCreateInput) error {
	if input.Name == "" {
		return awserrors.NewAWSError("InvalidRequestException", "Name is required", 400)
	}
	if err := validateWorkGroupName(input.Name); err != nil {
		return err
	}
	if input.Description != "" {
		if err := validateWorkGroupDescriptionString(input.Description); err != nil {
			return err
		}
	}
	if input.Config != nil {
		if input.Config.BytesScannedCutoff != nil {
			if err := validateBytesScannedCutoff(*input.Config.BytesScannedCutoff); err != nil {
				return err
			}
		}
		if input.Config.AdditionalConfiguration != "" {
			if err := validateAdditionalConfiguration(input.Config.AdditionalConfiguration); err != nil {
				return err
			}
		}
		if input.Config.ExecutionRole != "" {
			if err := validateExecutionRole(input.Config.ExecutionRole); err != nil {
				return err
			}
		}
	}
	if len(input.Tags) > 0 {
		if err := validateTags(input.Tags); err != nil {
			return err
		}
	}

	wg := &athenastore.WorkGroup{
		Name:        input.Name,
		Description: input.Description,
		State:       athenastore.WorkGroupStateEnabled,
	}
	if input.Config != nil {
		wg.Configuration = configInputToStore(input.Config)
	}

	if err := stores.workGroupStore.CreateWorkGroup(wg); err != nil {
		if err == athenastore.ErrWorkGroupAlreadyExists {
			return alreadyExistsInvalidRequest("WorkGroup", input.Name)
		}
		return err
	}

	if len(input.Tags) > 0 {
		arn := stores.workGroupStore.GetARN(input.Name)
		if err := stores.workGroupStore.Tag(arn, input.Tags); err != nil {
			return err
		}
	}
	return nil
}

// deleteWorkGroupCore deletes a workgroup with full cascade cleanup of
// dependent resources (named queries, prepared statements, query executions,
// and results). Both the HTTP handler and the admin handler must call this
// to avoid orphaned data.
func deleteWorkGroupCore(stores *athenaStores, name string) error {
	if name == "" {
		return awserrors.NewAWSError("InvalidRequestException", "WorkGroup is required", 400)
	}
	if name == "primary" {
		return awserrors.NewAWSError("InvalidRequestException",
			"The primary work group cannot be deleted.", 400)
	}

	stores.namedQueryStore.DeleteNamedQueriesByWorkGroup(name)
	stores.preparedStatementStore.DeletePreparedStatementsByWorkGroup(name)
	deletedQEIds, _ := stores.queryExecutionStore.DeleteQueryExecutionsByWorkGroup(name)
	stores.resultStore.DeleteResultsByIDs(deletedQEIds)

	if err := stores.workGroupStore.DeleteWorkGroup(name); err != nil {
		if err == athenastore.ErrWorkGroupNotFound {
			return workGroupNotFound(name)
		}
		return err
	}
	return nil
}

// listWorkGroupsCore lists workgroups with pagination.
func listWorkGroupsCore(stores *athenaStores, maxResults int, marker string) (*WorkGroupListResult, error) {
	result, err := stores.workGroupStore.ListWorkGroups(storecommon.ListOptions{
		Marker:   marker,
		MaxItems: maxResults,
	})
	if err != nil {
		return nil, err
	}

	out := &WorkGroupListResult{
		Items:      make([]WorkGroupOut, len(result.Items)),
		NextMarker: result.NextMarker,
	}
	for i, wg := range result.Items {
		out.Items[i] = WorkGroupOut{
			Name:         wg.Name,
			State:        string(wg.State),
			Description:  wg.Description,
			CreationTime: wg.CreatedTime,
		}
	}
	return out, nil
}

func configInputToStore(in *WorkGroupConfigInput) *athenastore.WorkGroupConfiguration {
	cfg := &athenastore.WorkGroupConfiguration{
		EnforceWorkGroupConfiguration:        in.EnforceConfig,
		PublishCloudWatchMetricsEnabled:      in.PublishMetrics,
		BytesScannedCutoffPerQuery:           in.bytesScannedCutoffValue(),
		RequesterPaysEnabled:                 in.RequesterPaysEnabled,
		AdditionalConfiguration:              in.AdditionalConfiguration,
		ExecutionRole:                        in.ExecutionRole,
		EnableMinimumEncryptionConfiguration: in.EnableMinimumEncryptionConfiguration,
	}
	if in.OutputLocation != "" {
		cfg.ResultConfiguration = &athenastore.ResultConfiguration{
			OutputLocation: in.OutputLocation,
		}
	}
	if in.CustomerContentEncryptionKmsKey != "" {
		cfg.CustomerContentEncryptionConfiguration = &athenastore.CustomerContentEncryptionConfiguration{
			KmsKey: in.CustomerContentEncryptionKmsKey,
		}
	}
	evSelected := in.EngineVersionSelected
	if evSelected == "" {
		evSelected = "AUTO"
	}
	evEffective := in.EngineVersionEffective
	if evEffective == "" {
		evEffective = "Athena engine version 3"
	}
	cfg.EngineVersion = &athenastore.EngineVersion{
		SelectedEngineVersion:  evSelected,
		EffectiveEngineVersion: evEffective,
	}
	return cfg
}

// getStoresForRegion returns the full athenaStores for the given region,
// creating a new store group if not already cached. This replaces the old
// GetWorkGroupStoreForRegion method and allows the admin handler to access
// all stores (needed for cascade cleanup in deleteWorkGroupCore).
func (s *AthenaService) getStoresForRegion(region string) (*athenaStores, error) {
	if v, ok := s.stores.Load(region); ok {
		return v.(*athenaStores), nil
	}
	if s.storageManager == nil {
		return nil, fmt.Errorf("athena storage manager not initialised")
	}
	st, err := s.storageManager.GetStorage(region)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage for region %s: %w", region, err)
	}
	stores := &athenaStores{
		workGroupStore:           athenastore.NewWorkGroupStore(st, s.accountID, region),
		namedQueryStore:          athenastore.NewNamedQueryStore(st, region),
		preparedStatementStore:   athenastore.NewPreparedStatementStore(st, region),
		queryExecutionStore:      athenastore.NewQueryExecutionStore(st, region),
		resultStore:              athenastore.NewResultStore(st, region),
		dataCatalogStore:         athenastore.NewDataCatalogStore(st, s.accountID, region),
		databaseStore:            athenastore.NewDatabaseStore(st, region),
		tableStore:               athenastore.NewTableStore(st, region),
		tableDataStore:           athenastore.NewTableDataStore(st, region),
		capacityReservationStore: athenastore.NewCapacityReservationStore(st, s.accountID, region),
	}
	actual, _ := s.stores.LoadOrStore(region, stores)
	return actual.(*athenaStores), nil
}

// getWorkGroupCore fetches a workgroup, mapping the store not-found sentinel
// onto the API error. The store is acquired after the name validation, the
// order the original handler applied.
func (s *AthenaService) getWorkGroupCore(reqCtx *request.RequestContext, name string) (*athenastore.WorkGroup, error) {
	if name == "" {
		return nil, ErrInvalidRequestException
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	workGroup, err := stores.workGroupStore.GetWorkGroup(name)
	if err != nil {
		if err == athenastore.ErrWorkGroupNotFound {
			return nil, workGroupNotFound(name)
		}
		return nil, err
	}
	return workGroup, nil
}

// UpdateWorkGroupInput carries the parsed wire members of an
// UpdateWorkGroup request; ConfigurationUpdates travels as the raw wire map
// so the Core applies the same update ladder the handler applied inline.
type UpdateWorkGroupInput struct {
	WorkGroup            string
	Description          string
	State                string
	ConfigurationUpdates map[string]interface{}
}

// updateWorkGroupCore validates the update request, applies the description,
// state and configuration updates presence-based onto the stored record and
// persists it. The store is acquired after the name validation, the order
// the original handler applied.
func (s *AthenaService) updateWorkGroupCore(reqCtx *request.RequestContext, input UpdateWorkGroupInput) error {
	if input.WorkGroup == "" {
		return ErrInvalidRequestException
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	workGroup, err := stores.workGroupStore.GetWorkGroup(input.WorkGroup)
	if err != nil {
		if err == athenastore.ErrWorkGroupNotFound {
			return workGroupNotFound(input.WorkGroup)
		}
		return err
	}

	if input.Description != "" {
		if err := validateWorkGroupDescriptionString(input.Description); err != nil {
			return err
		}
		workGroup.Description = input.Description
	}

	if input.State != "" {
		if err := validateWorkGroupState(input.State); err != nil {
			return err
		}
		workGroup.State = athenastore.WorkGroupState(input.State)
	}

	if input.ConfigurationUpdates != nil {
		if err := applyConfigurationUpdates(workGroup, input.ConfigurationUpdates); err != nil {
			return err
		}
	}

	if err := stores.workGroupStore.UpdateWorkGroup(workGroup); err != nil {
		return err
	}

	return nil
}

// validateResourceExists checks that the taggable resource named by a
// TagResource/UntagResource/ListTagsForResource ARN exists, returning the
// per-type not-found error.
func validateResourceExists(stores *athenaStores, resourceType, resourceName string) error {
	switch resourceType {
	case "workgroup":
		_, err := stores.workGroupStore.GetWorkGroup(resourceName)
		if err != nil {
			return workGroupNotFound(resourceName)
		}
	case "datacatalog":
		_, err := stores.dataCatalogStore.GetDataCatalog(resourceName)
		if err != nil {
			return dataCatalogNotFound(resourceName)
		}
	case "capacityreservation":
		_, err := stores.capacityReservationStore.GetCapacityReservation(resourceName)
		if err != nil {
			return awserrors.NewResourceNotFoundException("CapacityReservation", resourceName)
		}
	default:
		return ErrInvalidRequestException
	}
	return nil
}

// TagResourceInput carries the parsed wire members of a TagResource
// request.
type TagResourceInput struct {
	ResourceARN string
	Tags        map[string]string
}

// tagResourceCore validates the tags, checks the tagged resource exists and
// applies the tags through the per-type store. The store is acquired only
// after the ARN and tag validation, the order the original handler applied.
func (s *AthenaService) tagResourceCore(reqCtx *request.RequestContext, input TagResourceInput) error {
	if input.ResourceARN == "" {
		return ErrInvalidRequestException
	}

	matches := arnRegex.FindStringSubmatch(input.ResourceARN)
	if matches == nil {
		return ErrInvalidRequestException
	}

	resourceType := matches[1]
	resourceName := matches[2]

	if err := validateTags(input.Tags); err != nil {
		return err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	if err := validateResourceExists(stores, resourceType, resourceName); err != nil {
		return err
	}

	resourceArn := normalizeAthenaARN(input.ResourceARN, s.accountID)

	if len(input.Tags) > 0 {
		switch resourceType {
		case "workgroup":
			if err := stores.workGroupStore.Tag(resourceArn, input.Tags); err != nil {
				return err
			}
		case "datacatalog":
			if err := stores.dataCatalogStore.Tag(resourceArn, input.Tags); err != nil {
				return err
			}
		case "capacityreservation":
			if err := stores.capacityReservationStore.Tag(resourceArn, input.Tags); err != nil {
				return err
			}
		default:
			return ErrInvalidRequestException
		}
	}

	return nil
}

// UntagResourceInput carries the parsed wire members of an UntagResource
// request.
type UntagResourceInput struct {
	ResourceARN string
	TagKeys     []string
}

// untagResourceCore checks the tagged resource exists and removes the tag
// keys through the per-type store. The store is acquired only after the ARN
// validation, the order the original handler applied.
func (s *AthenaService) untagResourceCore(reqCtx *request.RequestContext, input UntagResourceInput) error {
	if input.ResourceARN == "" {
		return ErrInvalidRequestException
	}

	matches := arnRegex.FindStringSubmatch(input.ResourceARN)
	if matches == nil {
		return ErrInvalidRequestException
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	if err := validateResourceExists(stores, matches[1], matches[2]); err != nil {
		return err
	}

	resourceArn := normalizeAthenaARN(input.ResourceARN, s.accountID)

	if len(input.TagKeys) > 0 {
		switch matches[1] {
		case "workgroup":
			if err := stores.workGroupStore.Untag(resourceArn, input.TagKeys); err != nil {
				return err
			}
		case "datacatalog":
			if err := stores.dataCatalogStore.Untag(resourceArn, input.TagKeys); err != nil {
				return err
			}
		case "capacityreservation":
			if err := stores.capacityReservationStore.Untag(resourceArn, input.TagKeys); err != nil {
				return err
			}
		default:
			return ErrInvalidRequestException
		}
	}

	return nil
}

// ListTagsForResourceInput carries the parsed wire members of a
// ListTagsForResource request; the MaxResults window travels presence-
// flagged because its default is the dynamic full-list length.
type ListTagsForResourceInput struct {
	ResourceARN   string
	MaxResults    int
	HasMaxResults bool
	NextToken     string
}

// listTagsForResourceCore checks the tagged resource exists, lists its tags
// through the per-type store and pages the sorted tag list by key with the
// documented window semantics (default: every tag, minimum 75, no upper
// bound). The store is acquired only after the ARN validation, the order
// the original handler applied.
func (s *AthenaService) listTagsForResourceCore(reqCtx *request.RequestContext, input ListTagsForResourceInput) ([]map[string]interface{}, string, error) {
	if input.ResourceARN == "" {
		return nil, "", ErrInvalidRequestException
	}

	matches := arnRegex.FindStringSubmatch(input.ResourceARN)
	if matches == nil {
		return nil, "", ErrInvalidRequestException
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, "", err
	}

	if err := validateResourceExists(stores, matches[1], matches[2]); err != nil {
		return nil, "", err
	}

	resourceArn := normalizeAthenaARN(input.ResourceARN, s.accountID)

	var tags map[string]string
	switch matches[1] {
	case "workgroup":
		tags, err = stores.workGroupStore.List(resourceArn)
	case "datacatalog":
		tags, err = stores.dataCatalogStore.List(resourceArn)
	case "capacityreservation":
		tags, err = stores.capacityReservationStore.List(resourceArn)
	default:
		return nil, "", ErrInvalidRequestException
	}
	if err != nil {
		return nil, "", err
	}

	tagList := tagutil.MapToResponse(tags)
	sort.Slice(tagList, func(i, j int) bool {
		return tagList[i]["Key"].(string) < tagList[j]["Key"].(string)
	})

	// MaxTagsCount carries only a documented minimum of 75 — no upper
	// bound exists, so any value of at least 75 is accepted and the page
	// is served up to the platform's hard pagination cap.
	maxResults := len(tagList)
	if input.HasMaxResults {
		if input.MaxResults < 75 {
			return nil, "", invalidRequestParameter(
				fmt.Sprintf("MaxResults must be at least 75 (got %d)", input.MaxResults))
		}
		maxResults = input.MaxResults
	}

	pageResult := pagination.PaginateSlice(tagList, input.NextToken, maxResults, func(item map[string]interface{}) string {
		return item["Key"].(string)
	})

	return pageResult.Items, pageResult.NextMarker, nil
}

// parseEngineVersion converts the raw EngineVersion wire map into the store
// structure, applying the documented defaults.
func parseEngineVersion(engineVersion map[string]interface{}) *athenastore.EngineVersion {
	ev := &athenastore.EngineVersion{}
	if selected, ok := engineVersion["SelectedEngineVersion"].(string); ok {
		ev.SelectedEngineVersion = selected
	}
	if effective, ok := engineVersion["EffectiveEngineVersion"].(string); ok {
		ev.EffectiveEngineVersion = effective
	}
	if ev.SelectedEngineVersion == "" {
		ev.SelectedEngineVersion = "AUTO"
	}
	if ev.EffectiveEngineVersion == "" {
		ev.EffectiveEngineVersion = "Athena engine version 3"
	}
	return ev
}

// applyConfigurationUpdates applies a WorkGroupConfigurationUpdates wire map
// onto the stored workgroup configuration, validating every provided member
// and honouring the Remove* flags.
func applyConfigurationUpdates(workGroup *athenastore.WorkGroup, updates map[string]interface{}) error {
	if workGroup.Configuration == nil {
		workGroup.Configuration = &athenastore.WorkGroupConfiguration{}
	}

	if resultConfigUpdatesRaw, ok := updates["ResultConfigurationUpdates"]; ok {
		if resultConfigUpdates, ok := resultConfigUpdatesRaw.(map[string]interface{}); ok {
			if workGroup.Configuration.ResultConfiguration == nil {
				workGroup.Configuration.ResultConfiguration = &athenastore.ResultConfiguration{}
			}
			rc := workGroup.Configuration.ResultConfiguration

			if outputLocation, ok := resultConfigUpdates["OutputLocation"].(string); ok {
				rc.OutputLocation = outputLocation
			}
			if encConfigMap, ok := resultConfigUpdates["EncryptionConfiguration"].(map[string]interface{}); ok {
				rc.EncryptionConfiguration = &athenastore.EncryptionConfiguration{}
				if encOption, ok := encConfigMap["EncryptionOption"].(string); ok {
					rc.EncryptionConfiguration.EncryptionOption = encOption
				}
				if kmsKey, ok := encConfigMap["KmsKey"].(string); ok {
					rc.EncryptionConfiguration.KmsKey = kmsKey
				}
			}
			if expectedBucketOwner, ok := resultConfigUpdates["ExpectedBucketOwner"].(string); ok {
				rc.ExpectedBucketOwner = expectedBucketOwner
			}
			if aclConfigMap, ok := resultConfigUpdates["AclConfiguration"].(map[string]interface{}); ok {
				aclOption, _ := aclConfigMap["S3AclOption"].(string)
				if aclOption != "BUCKET_OWNER_FULL_CONTROL" {
					return invalidRequestParameter("AclConfiguration.S3AclOption must be BUCKET_OWNER_FULL_CONTROL")
				}
				rc.ACLConfiguration = &athenastore.ACLConfiguration{S3ACLOption: aclOption}
			}

			if remove, ok := resultConfigUpdates["RemoveOutputLocation"].(bool); ok && remove {
				rc.OutputLocation = ""
			}
			if remove, ok := resultConfigUpdates["RemoveEncryptionConfiguration"].(bool); ok && remove {
				rc.EncryptionConfiguration = nil
			}
			if remove, ok := resultConfigUpdates["RemoveExpectedBucketOwner"].(bool); ok && remove {
				rc.ExpectedBucketOwner = ""
			}
			if remove, ok := resultConfigUpdates["RemoveAclConfiguration"].(bool); ok && remove {
				rc.ACLConfiguration = nil
			}
		}
	}

	if enforce, ok := updates["EnforceWorkGroupConfiguration"].(bool); ok {
		workGroup.Configuration.EnforceWorkGroupConfiguration = enforce
	}

	if bytesScanned, ok := updates["BytesScannedCutoffPerQuery"].(float64); ok {
		workGroup.Configuration.BytesScannedCutoffPerQuery = int64(bytesScanned)
		if err := validateBytesScannedCutoff(workGroup.Configuration.BytesScannedCutoffPerQuery); err != nil {
			return err
		}
	}

	if remove, ok := updates["RemoveBytesScannedCutoffPerQuery"].(bool); ok && remove {
		workGroup.Configuration.BytesScannedCutoffPerQuery = 0
	}

	if requesterPays, ok := updates["RequesterPaysEnabled"].(bool); ok {
		workGroup.Configuration.RequesterPaysEnabled = requesterPays
	}

	if publish, ok := updates["PublishCloudWatchMetricsEnabled"].(bool); ok {
		workGroup.Configuration.PublishCloudWatchMetricsEnabled = publish
	}

	if engineVersionRaw, ok := updates["EngineVersion"]; ok {
		if engineVersion, ok := engineVersionRaw.(map[string]interface{}); ok {
			workGroup.Configuration.EngineVersion = parseEngineVersion(engineVersion)
		}
	}

	if additional, ok := updates["AdditionalConfiguration"].(string); ok {
		if err := validateAdditionalConfiguration(additional); err != nil {
			return err
		}
		workGroup.Configuration.AdditionalConfiguration = additional
	}

	if executionRole, ok := updates["ExecutionRole"].(string); ok {
		if err := validateExecutionRole(executionRole); err != nil {
			return err
		}
		workGroup.Configuration.ExecutionRole = executionRole
	}

	if custEncMap, ok := updates["CustomerContentEncryptionConfiguration"].(map[string]interface{}); ok {
		workGroup.Configuration.CustomerContentEncryptionConfiguration = &athenastore.CustomerContentEncryptionConfiguration{}
		if kmsKey, ok := custEncMap["KmsKey"].(string); ok {
			workGroup.Configuration.CustomerContentEncryptionConfiguration.KmsKey = kmsKey
		}
	}

	if remove, ok := updates["RemoveCustomerContentEncryptionConfiguration"].(bool); ok && remove {
		workGroup.Configuration.CustomerContentEncryptionConfiguration = nil
	}

	if enableMin, ok := updates["EnableMinimumEncryptionConfiguration"].(bool); ok {
		workGroup.Configuration.EnableMinimumEncryptionConfiguration = enableMin
	}

	return nil
}

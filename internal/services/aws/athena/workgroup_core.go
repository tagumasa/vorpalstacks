package athena

import (
	"fmt"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
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
	OutputLocation          string
	EnforceConfig           bool
	PublishMetrics          bool
	BytesScannedCutoff      int64
	RequesterPaysEnabled    bool
	EngineVersionSelected   string
	EngineVersionEffective  string
	AdditionalConfiguration string
	ExecutionRole           string
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
	if input.Config != nil && input.Config.BytesScannedCutoff != 0 {
		if err := validateBytesScannedCutoff(input.Config.BytesScannedCutoff); err != nil {
			return err
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
			return ErrResourceAlreadyExistsException
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

// --- Internal converters ---

// configInputToStore converts a service-layer WorkGroupConfigInput to the
// store-layer WorkGroupConfiguration type.
func configInputToStore(in *WorkGroupConfigInput) *athenastore.WorkGroupConfiguration {
	cfg := &athenastore.WorkGroupConfiguration{
		EnforceWorkGroupConfiguration:   in.EnforceConfig,
		PublishCloudWatchMetricsEnabled: in.PublishMetrics,
		BytesScannedCutoffPerQuery:      in.BytesScannedCutoff,
		RequesterPaysEnabled:            in.RequesterPaysEnabled,
		AdditionalConfiguration:         in.AdditionalConfiguration,
		ExecutionRole:                   in.ExecutionRole,
	}
	if in.OutputLocation != "" {
		cfg.ResultConfiguration = &athenastore.ResultConfiguration{
			OutputLocation: in.OutputLocation,
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

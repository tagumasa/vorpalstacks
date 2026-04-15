package neptune

import (
	"context"
	"fmt"
	"net/http"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	neptunestore "vorpalstacks/internal/store/aws/neptune"
)

// CreateDBInstance creates a new Neptune DB instance within a cluster.
func (s *NeptuneService) CreateDBInstance(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	id := request.GetStringParam(params, "DBInstanceIdentifier")
	if id == "" {
		return nil, fmt.Errorf("neptune: DBInstanceIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	clusterID := request.GetStringParam(params, "DBClusterIdentifier")
	if clusterID != "" {
		if _, err := store.GetCluster(clusterID); err != nil {
			return nil, fmt.Errorf("neptune: DBCluster %s not found", clusterID)
		}
	}

	now := time.Now()
	instance := &neptunestore.DBInstance{
		DBInstanceIdentifier:            id,
		DBClusterIdentifier:             clusterID,
		Engine:                          request.GetStringParam(params, "Engine"),
		EngineVersion:                   request.GetStringParam(params, "EngineVersion"),
		DBInstanceClass:                 request.GetStringParam(params, "DBInstanceClass"),
		Status:                          "available",
		AvailabilityZone:                request.GetStringParam(params, "AvailabilityZone"),
		DBParameterGroupName:            request.GetStringParam(params, "DBParameterGroupName"),
		DBSubnetGroupName:               request.GetStringParam(params, "DBSubnetGroupName"),
		EnableIAMDatabaseAuthentication: request.GetBoolParam(params, "EnableIAMDatabaseAuthentication"),
		PubliclyAccessible:              request.GetBoolParam(params, "PubliclyAccessible"),
		AutoMinorVersionUpgrade:         request.GetBoolParam(params, "AutoMinorVersionUpgrade"),
		InstanceCreateTime:              &now,
		AccountID:                       reqCtx.GetAccountID(),
		Region:                          reqCtx.GetRegion(),
		DBInstanceArn:                   fmt.Sprintf("arn:aws:rds:%s:%s:db:%s", reqCtx.GetRegion(), reqCtx.GetAccountID(), id),
	}

	if err := store.CreateInstance(instance); err != nil {
		return nil, translateStoreError(err)
	}

	recordEvent(store, "db-instance", id, instance.DBInstanceArn,
		fmt.Sprintf("DB instance %s created", id), []string{"creation"})

	return map[string]interface{}{
		"DBInstance": instance,
	}, nil
}

// DeleteDBInstance deletes the specified Neptune DB instance.
func (s *NeptuneService) DeleteDBInstance(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	id := request.GetStringParam(params, "DBInstanceIdentifier")
	if id == "" {
		return nil, fmt.Errorf("neptune: DBInstanceIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	instance, err := store.GetInstance(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	instance.Status = "deleting"

	if err := store.DeleteInstance(id); err != nil {
		return nil, translateStoreError(err)
	}

	recordEvent(store, "db-instance", id, instance.DBInstanceArn,
		fmt.Sprintf("DB instance %s deleted", id), []string{"deletion"})

	return map[string]interface{}{
		"DBInstance": instance,
	}, nil
}

// ModifyDBInstance updates the configuration of the specified DB instance.
func (s *NeptuneService) ModifyDBInstance(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	id := request.GetStringParam(params, "DBInstanceIdentifier")
	if id == "" {
		return nil, fmt.Errorf("neptune: DBInstanceIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	instance, err := store.GetInstance(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	if v := request.GetStringParam(params, "DBInstanceClass"); v != "" {
		instance.DBInstanceClass = v
	}
	if v := request.GetStringParam(params, "EngineVersion"); v != "" {
		instance.EngineVersion = v
	}
	if v := request.GetStringParam(params, "DBParameterGroupName"); v != "" {
		instance.DBParameterGroupName = v
	}
	if v := request.GetStringParam(params, "PreferredMaintenanceWindow"); v != "" {
		instance.PreferredMaintenanceWindow = v
	}
	if request.HasParam(params, "PubliclyAccessible") {
		instance.PubliclyAccessible = request.GetBoolParam(params, "PubliclyAccessible")
	}
	if request.HasParam(params, "AutoMinorVersionUpgrade") {
		instance.AutoMinorVersionUpgrade = request.GetBoolParam(params, "AutoMinorVersionUpgrade")
	}
	if request.HasParam(params, "EnableIAMDatabaseAuthentication") {
		instance.EnableIAMDatabaseAuthentication = request.GetBoolParam(params, "EnableIAMDatabaseAuthentication")
	}
	if request.HasParam(params, "CopyTagsToSnapshot") {
		instance.CopyTagsToSnapshot = request.GetBoolParam(params, "CopyTagsToSnapshot")
	}

	if err := store.UpdateInstance(instance); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBInstance": instance,
	}, nil
}

// DescribeDBInstances returns information about the specified DB instance or lists
// all instances when no identifier is provided.
func (s *NeptuneService) DescribeDBInstances(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	instanceID := request.GetStringParam(params, "DBInstanceIdentifier")
	if instanceID != "" {
		instance, err := store.GetInstance(instanceID)
		if err != nil {
			return nil, translateStoreError(err)
		}
		return map[string]interface{}{
			"DBInstances": protocol.XMLElements{ElementName: "DBInstance", Items: []interface{}{instance}},
		}, nil
	}

	instances, err := store.ListInstances()
	if err != nil {
		return nil, translateStoreError(err)
	}

	items := make([]interface{}, 0, len(instances))
	for _, i := range instances {
		items = append(items, i)
	}

	marker := request.GetStringParam(params, "Marker")
	maxRecords := request.GetIntParam(params, "MaxRecords")
	resultItems, nextMarker, isTruncated := paginateItems(items, marker, maxRecords, func(item interface{}) string {
		return item.(*neptunestore.DBInstance).DBInstanceIdentifier
	})

	result := map[string]interface{}{
		"DBInstances": protocol.XMLElements{ElementName: "DBInstance", Items: resultItems},
	}
	if isTruncated {
		result["Marker"] = nextMarker
	}
	return result, nil
}

// RebootDBInstance reboots the specified DB instance.
func (s *NeptuneService) RebootDBInstance(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	id := request.GetStringParam(params, "DBInstanceIdentifier")
	if id == "" {
		return nil, fmt.Errorf("neptune: DBInstanceIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	instance, err := store.GetInstance(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	if instance.Status != "available" {
		return nil, awserrors.NewAWSError("InvalidDBInstanceStateFault", fmt.Sprintf("instance %s is not in available state", id), http.StatusBadRequest)
	}

	instance.Status = "rebooting"
	if err := store.UpdateInstance(instance); err != nil {
		return nil, translateStoreError(err)
	}

	instance.Status = "available"
	if err := store.UpdateInstance(instance); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBInstance": instance,
	}, nil
}

package neptune

import (
	"context"
	"fmt"
	"time"

	"vorpalstacks/internal/services/aws/common/protocol"
	"vorpalstacks/internal/services/aws/common/request"
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

	if _, err := store.GetInstance(id); err == nil {
		return nil, fmt.Errorf("neptune: DBInstance %s already exists", id)
	}

	now := time.Now()
	instance := &neptunestore.DBInstance{
		DBInstanceIdentifier:            id,
		DBClusterIdentifier:             request.GetStringParam(params, "DBClusterIdentifier"),
		Engine:                          request.GetStringParam(params, "Engine"),
		EngineVersion:                   request.GetStringParam(params, "EngineVersion"),
		DBInstanceClass:                 request.GetStringParam(params, "DBInstanceClass"),
		Status:                          "creating",
		AvailabilityZone:                request.GetStringParam(params, "AvailabilityZone"),
		DBParameterGroupName:            request.GetStringParam(params, "DBParameterGroupName"),
		DBSubnetGroupName:               request.GetStringParam(params, "DBSubnetGroupName"),
		EnableIAMDatabaseAuthentication: request.GetBoolParam(params, "EnableIAMDatabaseAuthentication"),
		PubliclyAccessible:              request.GetBoolParam(params, "PubliclyAccessible"),
		AutoMinorVersionUpgrade:         request.GetBoolParam(params, "AutoMinorVersionUpgrade"),
		InstanceCreateTime:              &now,
		AccountID:                       reqCtx.GetAccountID(),
		Region:                          reqCtx.GetRegion(),
	}

	if err := store.CreateInstance(instance); err != nil {
		return nil, translateStoreError(err)
	}

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

	_ = store.DeleteInstance(id)
	instance.Status = "deleting"

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

	_ = store.UpdateInstance(instance)

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

	result := make([]interface{}, 0, len(instances))
	for _, i := range instances {
		result = append(result, i)
	}

	return map[string]interface{}{
		"DBInstances": protocol.XMLElements{ElementName: "DBInstance", Items: result},
	}, nil
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

	instance.Status = "rebooting"
	_ = store.UpdateInstance(instance)
	instance.Status = "available"
	_ = store.UpdateInstance(instance)

	return map[string]interface{}{
		"DBInstance": instance,
	}, nil
}

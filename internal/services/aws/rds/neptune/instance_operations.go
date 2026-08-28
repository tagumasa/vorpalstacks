package neptune

import (
	"context"

	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	rdssvc "vorpalstacks/internal/services/aws/rds"
)

// CreateDBInstance creates a new Neptune DB instance within a cluster.
func (s *NeptuneService) CreateDBInstance(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := &CreateDBInstanceInput{
		CreateInstanceParams: rdssvc.CreateInstanceParams{
			DBInstanceIdentifier:            request.GetStringParam(params, "DBInstanceIdentifier"),
			DBClusterIdentifier:             request.GetStringParam(params, "DBClusterIdentifier"),
			Engine:                          request.GetStringParam(params, "Engine"),
			EngineVersion:                   request.GetStringParam(params, "EngineVersion"),
			DBInstanceClass:                 request.GetStringParam(params, "DBInstanceClass"),
			AvailabilityZone:                request.GetStringParam(params, "AvailabilityZone"),
			DBParameterGroupName:            request.GetStringParam(params, "DBParameterGroupName"),
			DBSubnetGroupName:               request.GetStringParam(params, "DBSubnetGroupName"),
			PubliclyAccessible:              request.GetBoolParam(params, "PubliclyAccessible"),
			AutoMinorVersionUpgrade:         request.GetBoolParam(params, "AutoMinorVersionUpgrade"),
			Port:                            int32(request.GetIntParam(params, "Port")),
			StorageType:                     request.GetStringParam(params, "StorageType"),
			BackupRetentionPeriod:           int32(request.GetIntParam(params, "BackupRetentionPeriod")),
			AllocatedStorage:                int32(request.GetIntParam(params, "AllocatedStorage")),
			StorageEncrypted:                request.GetBoolParam(params, "StorageEncrypted"),
			KmsKeyId:                        request.GetStringParam(params, "KmsKeyId"),
			DeletionProtection:              request.GetBoolParam(params, "DeletionProtection"),
			MonitoringInterval:              int32(request.GetIntParam(params, "MonitoringInterval")),
			MonitoringRoleArn:               request.GetStringParam(params, "MonitoringRoleArn"),
			EnablePerformanceInsights:       request.GetBoolParam(params, "EnablePerformanceInsights"),
			PerformanceInsightsKMSKeyId:     request.GetStringParam(params, "PerformanceInsightsKMSKeyId"),
			EnableIAMDatabaseAuthentication: request.GetBoolParam(params, "EnableIAMDatabaseAuthentication"),
			MasterUsername:                  request.GetStringParam(params, "MasterUsername"),
			CopyTagsToSnapshot:              request.GetBoolParam(params, "CopyTagsToSnapshot"),
			AccountID:                       reqCtx.GetAccountID(),
			Region:                          reqCtx.GetRegion(),
		},
		MasterUserPassword: request.GetStringParam(params, "MasterUserPassword"),
		Tags:               parseNeptuneTags(params),
	}
	return s.createDBInstanceCore(ctx, store, in)
}

// DeleteDBInstance deletes the specified Neptune DB instance.
func (s *NeptuneService) DeleteDBInstance(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &DeleteDBInstanceInput{
		DeleteInstanceParams: rdssvc.DeleteInstanceParams{
			DBInstanceIdentifier: request.GetStringParam(params, "DBInstanceIdentifier"),
			AccountID:            reqCtx.GetAccountID(),
			Region:               reqCtx.GetRegion(),
		},
	}
	return s.deleteDBInstanceCore(ctx, store, in)
}

// ModifyDBInstance updates the configuration of the specified DB instance.
func (s *NeptuneService) ModifyDBInstance(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &ModifyDBInstanceInput{
		DBInstanceIdentifier:            request.GetStringParam(params, "DBInstanceIdentifier"),
		DBInstanceClass:                 request.GetStringParam(params, "DBInstanceClass"),
		EngineVersion:                   request.GetStringParam(params, "EngineVersion"),
		DBParameterGroupName:            request.GetStringParam(params, "DBParameterGroupName"),
		PreferredMaintenanceWindow:      request.GetStringParam(params, "PreferredMaintenanceWindow"),
		HasPubliclyAccessible:           request.HasParam(params, "PubliclyAccessible"),
		PubliclyAccessible:              request.GetBoolParam(params, "PubliclyAccessible"),
		HasAutoMinorVersionUpgrade:      request.HasParam(params, "AutoMinorVersionUpgrade"),
		AutoMinorVersionUpgrade:         request.GetBoolParam(params, "AutoMinorVersionUpgrade"),
		HasEnableIAMDatabaseAuth:        request.HasParam(params, "EnableIAMDatabaseAuthentication"),
		EnableIAMDatabaseAuthentication: request.GetBoolParam(params, "EnableIAMDatabaseAuthentication"),
		HasCopyTagsToSnapshot:           request.HasParam(params, "CopyTagsToSnapshot"),
		CopyTagsToSnapshot:              request.GetBoolParam(params, "CopyTagsToSnapshot"),
		NewDBInstanceIdentifier:         request.GetStringParam(params, "NewDBInstanceIdentifier"),
	}
	return s.modifyDBInstanceCore(ctx, store, in)
}

// DescribeDBInstances returns information about the specified DB instance or lists
// all instances when no identifier is provided.
func (s *NeptuneService) DescribeDBInstances(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	instances, nextMarker, err := rdssvc.QueryInstances(store, rdssvc.DescribeDBInstancesInput{
		DBInstanceIdentifier: request.GetStringParam(params, "DBInstanceIdentifier"),
		Filters:              nil,
		Marker:               request.GetStringParam(params, "Marker"),
		MaxRecords:           int32(request.GetIntParam(params, "MaxRecords")),
	})
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(instances))
	for _, i := range instances {
		items = append(items, i)
	}

	result := map[string]interface{}{
		"DBInstances": protocol.XMLElements{ElementName: "DBInstance", Items: items},
	}
	if nextMarker != "" {
		result["Marker"] = nextMarker
	}
	return result, nil
}

// RebootDBInstance reboots the specified DB instance.
func (s *NeptuneService) RebootDBInstance(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &RebootDBInstanceInput{
		DBInstanceIdentifier: request.GetStringParam(req.Parameters, "DBInstanceIdentifier"),
	}
	return s.rebootDBInstanceCore(ctx, store, in)
}

package neptune

import (
	"context"

	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	rdssvc "vorpalstacks/internal/services/aws/rds"
)

// parseRestoreDBClusterInput extracts the restore request members both
// restore operations consume from the wire parameters.
func parseRestoreDBClusterInput(params map[string]interface{}) *RestoreDBClusterInput {
	return &RestoreDBClusterInput{
		DBClusterIdentifier:         request.GetStringParam(params, "DBClusterIdentifier"),
		SnapshotIdentifier:          request.GetStringParam(params, "SnapshotIdentifier"),
		SourceDBClusterIdentifier:   request.GetStringParam(params, "SourceDBClusterIdentifier"),
		Engine:                      request.GetStringParam(params, "Engine"),
		EngineVersion:               request.GetStringParam(params, "EngineVersion"),
		BackupRetentionPeriod:       request.GetIntParam(params, "BackupRetentionPeriod"),
		Port:                        request.GetIntParam(params, "Port"),
		DBClusterParameterGroupName: request.GetStringParam(params, "DBClusterParameterGroupName"),
		DBSubnetGroupName:           request.GetStringParam(params, "DBSubnetGroupName"),
		StorageEncrypted:            request.GetBoolParam(params, "StorageEncrypted"),
		DeletionProtection:          request.GetBoolParam(params, "DeletionProtection"),
		PreferredBackupWindow:       request.GetStringParam(params, "PreferredBackupWindow"),
		PreferredMaintenanceWindow:  request.GetStringParam(params, "PreferredMaintenanceWindow"),
	}
}

// CreateDBClusterSnapshot creates a new snapshot of the specified DB cluster.
func (s *NeptuneService) CreateDBClusterSnapshot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &CreateDBClusterSnapshotInput{
		DBClusterSnapshotIdentifier: request.GetStringParam(params, "DBClusterSnapshotIdentifier"),
		DBClusterIdentifier:         request.GetStringParam(params, "DBClusterIdentifier"),
		AccountID:                   reqCtx.GetAccountID(),
		Region:                      reqCtx.GetRegion(),
	}
	return s.createDBClusterSnapshotCore(ctx, store, in)
}

// DeleteDBClusterSnapshot deletes the specified DB cluster snapshot.
func (s *NeptuneService) DeleteDBClusterSnapshot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &DeleteDBClusterSnapshotInput{
		DBClusterSnapshotIdentifier: request.GetStringParam(req.Parameters, "DBClusterSnapshotIdentifier"),
	}
	return s.deleteDBClusterSnapshotCore(ctx, store, in)
}

// DescribeDBClusterSnapshots returns information about the specified cluster
// snapshot or lists all snapshots when no identifier is provided.
func (s *NeptuneService) DescribeDBClusterSnapshots(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	snapshots, nextMarker, err := rdssvc.QueryClusterSnapshots(store, rdssvc.DescribeDBClusterSnapshotsInput{
		DBClusterSnapshotIdentifier: request.GetStringParam(params, "DBClusterSnapshotIdentifier"),
		DBClusterIdentifier:         request.GetStringParam(params, "DBClusterIdentifier"),
		SnapshotType:                request.GetStringParam(params, "SnapshotType"),
		Filters:                     nil,
		Marker:                      request.GetStringParam(params, "Marker"),
		MaxRecords:                  int32(request.GetIntParam(params, "MaxRecords")),
	})
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(snapshots))
	for _, snap := range snapshots {
		items = append(items, snap)
	}

	result := map[string]interface{}{
		"DBClusterSnapshots": protocol.XMLElements{ElementName: "DBClusterSnapshot", Items: items},
	}
	if nextMarker != "" {
		result["Marker"] = nextMarker
	}
	return result, nil
}

// CopyDBClusterSnapshot creates a copy of the specified DB cluster snapshot.
func (s *NeptuneService) CopyDBClusterSnapshot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &CopyDBClusterSnapshotInput{
		SourceDBClusterSnapshotIdentifier: request.GetStringParam(params, "SourceDBClusterSnapshotIdentifier"),
		TargetDBClusterSnapshotIdentifier: request.GetStringParam(params, "TargetDBClusterSnapshotIdentifier"),
		AccountID:                         reqCtx.GetAccountID(),
		Region:                            reqCtx.GetRegion(),
	}
	return s.copyDBClusterSnapshotCore(ctx, store, in)
}

// DescribeDBClusterSnapshotAttributes returns the attributes of the specified
// DB cluster snapshot.
func (s *NeptuneService) DescribeDBClusterSnapshotAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &DescribeDBClusterSnapshotAttributesInput{
		DBClusterSnapshotIdentifier: request.GetStringParam(req.Parameters, "DBClusterSnapshotIdentifier"),
	}
	snap, err := s.getDBClusterSnapshotCore(ctx, store, in)
	if err != nil {
		return nil, translateStoreError(err)
	}

	attrValues := snap.RestoreAttributeValues
	if attrValues == nil {
		attrValues = []string{}
	}
	attrItems := make([]interface{}, 0, len(attrValues))
	for _, v := range attrValues {
		attrItems = append(attrItems, v)
	}

	return map[string]interface{}{
		"DBClusterSnapshotAttributesResult": map[string]interface{}{
			"DBClusterSnapshotIdentifier": in.DBClusterSnapshotIdentifier,
			"DBClusterSnapshotAttributes": protocol.XMLElements{ElementName: "DBClusterSnapshotAttribute", Items: []interface{}{
				map[string]interface{}{"AttributeName": "restore", "AttributeValues": protocol.XMLElements{ElementName: "AttributeValue", Items: attrItems}},
			}},
		},
	}, nil
}

// ModifyDBClusterSnapshotAttribute modifies an attribute of the specified DB
// cluster snapshot.
func (s *NeptuneService) ModifyDBClusterSnapshotAttribute(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &ModifyDBClusterSnapshotAttributeInput{
		DBClusterSnapshotIdentifier: request.GetStringParam(params, "DBClusterSnapshotIdentifier"),
		AttributeName:               request.GetStringParam(params, "AttributeName"),
		ValuesToAdd:                 getNeptuneStringList(params, "ValuesToAdd", "AttributeValue", "member"),
		ValuesToRemove:              getNeptuneStringList(params, "ValuesToRemove", "AttributeValue", "member"),
	}
	return s.modifyDBClusterSnapshotAttributeCore(ctx, store, in)
}

// RestoreDBClusterFromSnapshot creates a new DB cluster from a DB cluster
// snapshot. Cluster snapshots are metadata-only — they do not capture
// row-level data.
func (s *NeptuneService) RestoreDBClusterFromSnapshot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := parseRestoreDBClusterInput(req.Parameters)
	in.AccountID = reqCtx.GetAccountID()
	in.Region = reqCtx.GetRegion()
	return s.restoreDBClusterFromSnapshotCore(ctx, store, in)
}

// RestoreDBClusterToPointInTime restores a DB cluster to a point in time from
// a source cluster.
func (s *NeptuneService) RestoreDBClusterToPointInTime(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := parseRestoreDBClusterInput(req.Parameters)
	in.AccountID = reqCtx.GetAccountID()
	in.Region = reqCtx.GetRegion()
	return s.restoreDBClusterToPointInTimeCore(ctx, store, in)
}

// PromoteReadReplicaDBCluster promotes a read replica cluster to a standalone
// writable cluster.
func (s *NeptuneService) PromoteReadReplicaDBCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &PromoteReadReplicaDBClusterInput{
		DBClusterIdentifier: request.GetStringParam(req.Parameters, "DBClusterIdentifier"),
	}
	return s.promoteReadReplicaDBClusterCore(ctx, store, in)
}

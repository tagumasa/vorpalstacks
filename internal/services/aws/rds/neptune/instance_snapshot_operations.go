package neptune

import (
	"context"
	"fmt"
	"net/http"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	rdssvc "vorpalstacks/internal/services/aws/rds"
	rdsstore "vorpalstacks/internal/store/aws/rds"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// CreateDBSnapshot creates a manual snapshot of a DB instance.
func (s *NeptuneService) CreateDBSnapshot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	snapshotID := request.GetStringParam(params, "DBSnapshotIdentifier")
	if snapshotID == "" {
		return nil, awserrors.NewMissingParameter("DBSnapshotIdentifier is required")
	}
	if err := rdssvc.ValidateDBSnapshotIdentifier(snapshotID); err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}
	instanceID := request.GetStringParam(params, "DBInstanceIdentifier")
	if instanceID == "" {
		return nil, awserrors.NewMissingParameter("DBInstanceIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	instance, err := store.GetInstance(instanceID)
	if err != nil {
		return nil, translateStoreError(err)
	}

	now := time.Now()
	snap := &rdsstore.DBInstanceSnapshot{
		DBSnapshotIdentifier:   snapshotID,
		DBInstanceIdentifier:   instanceID,
		SnapshotCreateTime:     &now,
		InstanceCreateTime:     instance.InstanceCreateTime,
		Engine:                 instance.Engine,
		EngineVersion:          instance.EngineVersion,
		SnapshotType:           "manual",
		Status:                 "available",
		AvailabilityZone:       instance.AvailabilityZone,
		DBSnapshotArn:          arnutil.NewARNBuilder(reqCtx.GetAccountID(), reqCtx.GetRegion()).RDS().Snapshot(snapshotID),
		IAMDatabaseAuthEnabled: instance.IAMDatabaseAuthenticationEnabled,
		AccountID:              reqCtx.GetAccountID(),
		Region:                 reqCtx.GetRegion(),
	}

	if err := store.CreateInstanceSnapshot(snap); err != nil {
		return nil, translateStoreError(err)
	}

	// Capture row-level data for MySQL instances so that
	// RestoreDBInstanceFromDBSnapshot can recover user tables.
	if s.snapOp != nil && instance.Engine == "mysql" {
		if err := s.snapOp.SnapshotData(instanceID, snapshotID); err != nil {
			logs.Warn("neptune: SnapshotData for manual snapshot failed",
				logs.String("instance", instanceID),
				logs.String("snapshot", snapshotID),
				logs.Err(err))
		}
	}

	return map[string]interface{}{
		"DBSnapshot": instanceSnapshotToMap(snap, reqCtx),
	}, nil
}

// DescribeDBSnapshots lists DB instance snapshots with optional filtering.
func (s *NeptuneService) DescribeDBSnapshots(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	snapshots, err := store.ListInstanceSnapshots()
	if err != nil {
		return nil, translateStoreError(err)
	}

	filterSnapID := request.GetStringParam(req.Parameters, "DBSnapshotIdentifier")
	filterInstID := request.GetStringParam(req.Parameters, "DBInstanceIdentifier")
	filterType := request.GetStringParam(req.Parameters, "SnapshotType")

	result := make([]interface{}, 0, len(snapshots))
	for _, snap := range snapshots {
		if filterSnapID != "" && snap.DBSnapshotIdentifier != filterSnapID {
			continue
		}
		if filterInstID != "" && snap.DBInstanceIdentifier != filterInstID {
			continue
		}
		if filterType != "" && snap.SnapshotType != filterType {
			continue
		}
		result = append(result, instanceSnapshotToMap(snap, reqCtx))
	}

	marker := request.GetStringParam(req.Parameters, "Marker")
	maxRecords := request.GetIntParam(req.Parameters, "MaxRecords")
	resultItems, nextMarker, isTruncated := paginateItems(result, marker, maxRecords, func(item interface{}) string {
		m := item.(map[string]interface{})
		if v, ok := m["DBSnapshotIdentifier"].(string); ok {
			return v
		}
		return ""
	})

	resp := map[string]interface{}{
		"DBSnapshots": protocol.XMLElements{ElementName: "DBSnapshot", Items: resultItems},
	}
	if isTruncated {
		resp["Marker"] = nextMarker
	}
	return resp, nil
}

// DeleteDBSnapshot deletes a DB instance snapshot.
func (s *NeptuneService) DeleteDBSnapshot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	snapshotID := request.GetStringParam(req.Parameters, "DBSnapshotIdentifier")
	if snapshotID == "" {
		return nil, awserrors.NewMissingParameter("DBSnapshotIdentifier is required")
	}
	if err := rdssvc.ValidateDBSnapshotIdentifier(snapshotID); err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	snapshot, err := store.GetInstanceSnapshot(snapshotID)
	if err != nil {
		return nil, translateStoreError(err)
	}

	if err := store.DeleteInstanceSnapshot(snapshotID); err != nil {
		return nil, translateStoreError(err)
	}

	// Clean up row-level data for MySQL snapshots so it does not leak.
	if s.snapOp != nil && snapshot.Engine == "mysql" {
		if err := s.snapOp.DeleteSnapshotData(snapshotID); err != nil {
			logs.Warn("neptune: DeleteSnapshotData failed",
				logs.String("snapshot", snapshotID),
				logs.Err(err))
		}
	}

	return map[string]interface{}{
		"DBSnapshot": instanceSnapshotToMap(snapshot, reqCtx),
	}, nil
}

// instanceSnapshotToMap converts a DBInstanceSnapshot to the AWS RDS API
// response map with PascalCase keys matching the Query protocol shape.
func instanceSnapshotToMap(s *rdsstore.DBInstanceSnapshot, reqCtx *request.RequestContext) map[string]interface{} {
	m := map[string]interface{}{
		"DBSnapshotIdentifier":             s.DBSnapshotIdentifier,
		"DBInstanceIdentifier":             s.DBInstanceIdentifier,
		"Engine":                           s.Engine,
		"EngineVersion":                    s.EngineVersion,
		"SnapshotType":                     s.SnapshotType,
		"Status":                           s.Status,
		"AllocatedStorage":                 s.AllocatedStorage,
		"StorageType":                      s.StorageType,
		"Port":                             s.Port,
		"AvailabilityZone":                 s.AvailabilityZone,
		"VpcId":                            s.VpcId,
		"MasterUsername":                   s.MasterUsername,
		"LicenseModel":                     s.LicenseModel,
		"StorageEncrypted":                 s.StorageEncrypted,
		"KmsKeyId":                         s.KmsKeyId,
		"DBSnapshotArn":                    s.DBSnapshotArn,
		"IAMDatabaseAuthenticationEnabled": s.IAMDatabaseAuthEnabled,
		"OptionGroupName":                  s.OptionGroupName,
	}
	if s.SnapshotCreateTime != nil {
		m["SnapshotCreateTime"] = s.SnapshotCreateTime.UTC().Format("2006-01-02T15:04:05Z")
	}
	if s.InstanceCreateTime != nil {
		m["InstanceCreateTime"] = s.InstanceCreateTime.UTC().Format("2006-01-02T15:04:05Z")
	}
	return m
}

// CopyDBSnapshot creates a copy of the specified DB instance snapshot.
func (s *NeptuneService) CopyDBSnapshot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	sourceID := request.GetStringParam(params, "SourceDBSnapshotIdentifier")
	if sourceID == "" {
		return nil, awserrors.NewMissingParameter("SourceDBSnapshotIdentifier is required")
	}
	targetID := request.GetStringParam(params, "TargetDBSnapshotIdentifier")
	if targetID == "" {
		return nil, awserrors.NewMissingParameter("TargetDBSnapshotIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	source, err := store.GetInstanceSnapshot(sourceID)
	if err != nil {
		return nil, translateStoreError(err)
	}

	now := time.Now()
	cp := *source
	cp.DBSnapshotIdentifier = targetID
	cp.SnapshotCreateTime = &now
	if source.InstanceCreateTime != nil {
		ct := *source.InstanceCreateTime
		cp.InstanceCreateTime = &ct
	}
	cp.DBSnapshotArn = arnutil.NewARNBuilder(reqCtx.GetAccountID(), reqCtx.GetRegion()).RDS().Snapshot(targetID)

	if err := store.CreateInstanceSnapshot(&cp); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBSnapshot": instanceSnapshotToMap(&cp, reqCtx),
	}, nil
}

// DescribeDBSnapshotAttributes returns the restore attributes of the
// specified DB instance snapshot.
func (s *NeptuneService) DescribeDBSnapshotAttributes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	snapshotID := request.GetStringParam(params, "DBSnapshotIdentifier")
	if snapshotID == "" {
		return nil, awserrors.NewMissingParameter("DBSnapshotIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	snap, err := store.GetInstanceSnapshot(snapshotID)
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
		"DBSnapshotAttributesResult": map[string]interface{}{
			"DBSnapshotIdentifier": snapshotID,
			"DBSnapshotAttributes": protocol.XMLElements{ElementName: "DBSnapshotAttribute", Items: []interface{}{
				map[string]interface{}{"AttributeName": "restore", "AttributeValues": protocol.XMLElements{ElementName: "AttributeValue", Items: attrItems}},
			}},
		},
	}, nil
}

// ModifyDBSnapshotAttribute modifies a restore attribute of the specified DB
// instance snapshot.
func (s *NeptuneService) ModifyDBSnapshotAttribute(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	snapshotID := request.GetStringParam(params, "DBSnapshotIdentifier")
	if snapshotID == "" {
		return nil, awserrors.NewMissingParameter("DBSnapshotIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	snap, err := store.GetInstanceSnapshot(snapshotID)
	if err != nil {
		return nil, translateStoreError(err)
	}

	attrName := request.GetStringParam(params, "AttributeName")
	valuesToAdd := getNeptuneStringList(params, "ValuesToAdd", "AttributeValue", "member")
	valuesToRemove := getNeptuneStringList(params, "ValuesToRemove", "AttributeValue", "member")

	if attrName == "" {
		attrName = "restore"
	}
	if attrName != "restore" {
		return nil, awserrors.NewValidationException(fmt.Sprintf("Unsupported attribute name '%s', only 'restore' is supported", attrName))
	}

	existing := make(map[string]bool, len(snap.RestoreAttributeValues))
	for _, v := range snap.RestoreAttributeValues {
		existing[v] = true
	}
	for _, v := range valuesToRemove {
		delete(existing, v)
	}
	for _, v := range valuesToAdd {
		existing[v] = true
	}
	result := make([]string, 0, len(existing))
	for v := range existing {
		result = append(result, v)
	}
	snap.RestoreAttributeValues = result
	if err := store.UpdateInstanceSnapshot(snap); err != nil {
		return nil, translateStoreError(err)
	}

	attrItems := make([]interface{}, 0, len(snap.RestoreAttributeValues))
	for _, v := range snap.RestoreAttributeValues {
		attrItems = append(attrItems, v)
	}

	return map[string]interface{}{
		"DBSnapshotAttributesResult": map[string]interface{}{
			"DBSnapshotIdentifier": snapshotID,
			"DBSnapshotAttributes": protocol.XMLElements{ElementName: "DBSnapshotAttribute", Items: []interface{}{
				map[string]interface{}{"AttributeName": attrName, "AttributeValues": protocol.XMLElements{ElementName: "AttributeValue", Items: attrItems}},
			}},
		},
	}, nil
}

// ModifyDBSnapshot updates the engine version of a DB snapshot. Neptune does
// not support in-place snapshot engine upgrades, so this validates the
// snapshot exists and returns it unchanged.
func (s *NeptuneService) ModifyDBSnapshot(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	snapshotID := request.GetStringParam(params, "DBSnapshotIdentifier")
	if snapshotID == "" {
		return nil, awserrors.NewMissingParameter("DBSnapshotIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	snap, err := store.GetInstanceSnapshot(snapshotID)
	if err != nil {
		return nil, translateStoreError(err)
	}

	if v := request.GetStringParam(params, "EngineVersion"); v != "" {
		snap.EngineVersion = v
		if err := store.UpdateInstanceSnapshot(snap); err != nil {
			return nil, translateStoreError(err)
		}
	}

	return map[string]interface{}{
		"DBSnapshot": instanceSnapshotToMap(snap, reqCtx),
	}, nil
}

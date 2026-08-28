package neptune

import (
	"context"
	"fmt"
	"net/http"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	rdssvc "vorpalstacks/internal/services/aws/rds"
	neptunestore "vorpalstacks/internal/store/aws/rds/neptune"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// CreateDBClusterInput carries the wire-parsed CreateDBCluster request. The
// embedded CreateClusterParams holds the members shared with the admin-plane
// cluster Core; the remaining members are Neptune HTTP-plane specifics.
type CreateDBClusterInput struct {
	rdssvc.CreateClusterParams
	MasterUserPassword          string
	PreferredBackupWindow       string
	PreferredMaintenanceWindow  string
	KmsKeyId                    string
	ReplicationSourceIdentifier string
	GlobalClusterIdentifier     string
	StorageType                 string
	VpcSecurityGroupIds         []string
	HasServerlessV2Scaling      bool
	ServerlessV2MinCapacity     float64
	ServerlessV2MaxCapacity     float64
	Tags                        []types.Tag
}

// DeleteDBClusterInput carries the wire-parsed DeleteDBCluster request.
type DeleteDBClusterInput struct {
	rdssvc.DeleteClusterParams
}

// ModifyDBClusterInput carries the wire-parsed ModifyDBCluster request. The
// Has* members preserve the wire presence of optional SCALAR members so an
// omitted boolean keeps its stored value instead of resetting it; list
// members are value-gated (an empty list cannot carry a clear on the query
// wire).
type ModifyDBClusterInput struct {
	DBClusterIdentifier             string
	EngineVersion                   string
	DBClusterParameterGroupName     string
	Port                            int
	BackupRetentionPeriod           int
	PreferredBackupWindow           string
	PreferredMaintenanceWindow      string
	StorageType                     string
	MasterUserPassword              string
	NetworkType                     string
	HasServerlessV2Scaling          bool
	ServerlessV2MinCapacity         float64
	ServerlessV2MaxCapacity         float64
	HasDeletionProtection           bool
	DeletionProtection              bool
	HasEnableIAMDatabaseAuth        bool
	EnableIAMDatabaseAuthentication bool
	VpcSecurityGroupIds             []string
	EnabledCloudwatchLogsExports    []string
	NewDBClusterIdentifier          string
	Region                          string
}

// StartDBClusterInput carries the wire-parsed StartDBCluster request.
type StartDBClusterInput struct {
	DBClusterIdentifier string
	Region              string
}

// StopDBClusterInput carries the wire-parsed StopDBCluster request.
type StopDBClusterInput struct {
	DBClusterIdentifier string
}

// FailoverDBClusterInput carries the wire-parsed FailoverDBCluster request.
type FailoverDBClusterInput struct {
	DBClusterIdentifier string
}

// AddRoleToDBClusterInput carries the wire-parsed AddRoleToDBCluster request.
type AddRoleToDBClusterInput struct {
	DBClusterIdentifier string
	RoleArn             string
	FeatureName         string
}

// RemoveRoleFromDBClusterInput carries the wire-parsed RemoveRoleFromDBCluster
// request.
type RemoveRoleFromDBClusterInput struct {
	DBClusterIdentifier string
	RoleArn             string
	FeatureName         string
}

// enrichClusterWithTags serialises a cluster record, attaching the TagList
// read from the tag store.
func enrichClusterWithTags(store neptunestore.NeptuneStoreInterface, cluster *neptunestore.DBCluster) map[string]interface{} {
	m := clusterToResponseMap(cluster)
	tags, err := store.GetTags(cluster.DBClusterArn)
	if err != nil {
		logs.Warn("failed to get tags for cluster", logs.String("cluster", cluster.DBClusterIdentifier), logs.Err(err))
	}
	if len(tags) > 0 {
		tagItems := make([]interface{}, 0, len(tags))
		for _, t := range tags {
			tagItems = append(tagItems, map[string]interface{}{"Key": t.Key, "Value": t.Value})
		}
		m["TagList"] = protocol.XMLElements{ElementName: "Tag", Items: tagItems}
	}
	if roles, ok := m["AssociatedRoles"].([]interface{}); ok && len(roles) > 0 {
		m["AssociatedRoles"] = protocol.XMLElements{ElementName: "DBClusterRole", Items: roles}
	}
	return m
}

// setClusterEndpoint persists the connection endpoint on the cluster.  It
// constructs the address from endpointAddress, stores it in cluster.Endpoint,
// and writes the update through to Pebble.
func (s *NeptuneService) setClusterEndpoint(store neptunestore.NeptuneStoreInterface, cluster *neptunestore.DBCluster, enginePort int) {
	addr := s.endpointAddressFor(cluster.DBClusterIdentifier, cluster.Engine)
	if addr == "" || enginePort <= 0 {
		return
	}
	cluster.Endpoint = &neptunestore.Endpoint{Address: addr, Port: enginePort}
	// ReaderEndpoint mirrors the cluster endpoint for Neptune's
	// single-writer topology. AWS Neptune surfaces both endpoints.
	if cluster.ReaderEndpoint == nil {
		cluster.ReaderEndpoint = &neptunestore.Endpoint{Address: addr, Port: enginePort}
	}
	if err := store.UpdateCluster(cluster); err != nil {
		logs.Warn("failed to persist cluster endpoint", logs.String("cluster", cluster.DBClusterIdentifier), logs.Err(err))
	}
}

// reparentClusterResources migrates tags and updates instance references when
// a cluster is renamed. Tag-copy, instance-list, and instance-update failures
// are propagated so the caller can roll back the rename. Old-tag removal is
// best-effort because the tags have already been copied to the new ARN.
func reparentClusterResources(store neptunestore.NeptuneStoreInterface, oldArn, newArn, oldID, newID string) error {
	tags, err := store.GetTags(oldArn)
	if err != nil {
		return fmt.Errorf("reparent: failed to get tags from %s: %w", oldArn, err)
	}
	if len(tags) > 0 {
		if err := store.AddTags(newArn, tags); err != nil {
			return fmt.Errorf("reparent: failed to copy tags to %s: %w", newArn, err)
		}
		keys := make([]string, len(tags))
		for i, t := range tags {
			keys[i] = t.Key
		}
		if err := store.RemoveTags(oldArn, keys); err != nil {
			logs.Warn("reparent: failed to remove old cluster tags after copy", logs.Err(err))
		}
	}

	instances, err := store.ListInstances()
	if err != nil {
		return fmt.Errorf("reparent: failed to list instances: %w", err)
	}
	for _, inst := range instances {
		if inst.DBClusterIdentifier == oldID {
			inst.DBClusterIdentifier = newID
			if err := store.UpdateInstance(inst); err != nil {
				return fmt.Errorf("reparent: failed to update instance %s cluster ref: %w", inst.DBInstanceIdentifier, err)
			}
		}
	}
	return nil
}

func cascadeDeleteClusterResources(store neptunestore.NeptuneStoreInterface, cluster *neptunestore.DBCluster) {
	clusterID := cluster.DBClusterIdentifier

	instances, err := store.ListInstances()
	if err != nil {
		logs.Warn("cascade: failed to list instances", logs.Err(err))
	} else {
		for _, inst := range instances {
			if inst.DBClusterIdentifier == clusterID {
				if delErr := store.DeleteInstance(inst.DBInstanceIdentifier); delErr != nil {
					logs.Warn("cascade: failed to delete instance", logs.String("instance", inst.DBInstanceIdentifier), logs.Err(delErr))
				} else {
					removeTagsForResource(store, inst.DBInstanceArn)
				}
			}
		}
	}

	endpoints, err := store.ListClusterEndpoints(clusterID)
	if err != nil {
		logs.Warn("cascade: failed to list cluster endpoints", logs.Err(err))
	} else {
		for _, ep := range endpoints {
			if delErr := store.DeleteClusterEndpoint(ep.DBClusterEndpointIdentifier); delErr != nil {
				logs.Warn("cascade: failed to delete cluster endpoint", logs.String("endpoint", ep.DBClusterEndpointIdentifier), logs.Err(delErr))
			}
		}
	}

	removeTagsForResource(store, cluster.DBClusterArn)
}

// removeClusterFromGlobal removes the cluster's membership entry from its
// parent global cluster, if any. Used by DeleteDBCluster and by the
// CreateDBCluster tag-failure rollback path.
func removeClusterFromGlobal(store neptunestore.NeptuneStoreInterface, cluster *neptunestore.DBCluster) {
	if cluster.GlobalClusterIdentifier == "" {
		return
	}
	gc, err := store.GetGlobalCluster(cluster.GlobalClusterIdentifier)
	if err != nil {
		return
	}
	filtered := make([]neptunestore.GlobalClusterMember, 0, len(gc.GlobalClusterMembers))
	for _, m := range gc.GlobalClusterMembers {
		if m.DBClusterArn != cluster.DBClusterArn {
			filtered = append(filtered, m)
		}
	}
	gc.GlobalClusterMembers = filtered
	if err := store.UpdateGlobalCluster(gc); err != nil {
		logs.Warn("failed to remove cluster from global cluster members",
			logs.String("cluster", cluster.DBClusterIdentifier), logs.Err(err))
	}
}

// createDBClusterCore validates and persists a new DB cluster, registers
// global-cluster membership, opens the engine, tags the cluster and records
// the creation event.
func (s *NeptuneService) createDBClusterCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *CreateDBClusterInput) (interface{}, error) {
	createParams := in.CreateClusterParams
	if err := rdssvc.ValidateCreateClusterParams(store, createParams); err != nil {
		return nil, neptuneTranslateError(err)
	}

	replicationSource := in.ReplicationSourceIdentifier
	if replicationSource != "" {
		if _, err := store.GetCluster(replicationSource); err != nil {
			return nil, awserrors.NewAWSError("InvalidParameterValue", fmt.Sprintf("replication source cluster %s not found", replicationSource), http.StatusBadRequest)
		}
	}

	backupRetention := createParams.BackupRetentionPeriod
	if backupRetention == 0 {
		backupRetention = 1
	}
	if err := validateBackupRetentionPeriod(backupRetention); err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}

	masterPasswordHash, err := hashMasterPassword(in.MasterUserPassword)
	if err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}

	cluster := rdssvc.BuildCluster(createParams)
	cluster.BackupRetentionPeriod = backupRetention
	cluster.PreferredBackupWindow = in.PreferredBackupWindow
	cluster.PreferredMaintenanceWindow = in.PreferredMaintenanceWindow
	cluster.KmsKeyId = in.KmsKeyId
	cluster.ReplicationSourceIdentifier = replicationSource
	cluster.GlobalClusterIdentifier = in.GlobalClusterIdentifier
	cluster.StorageType = in.StorageType
	cluster.MasterUserPasswordHash = masterPasswordHash
	cluster.DbClusterResourceId = fmt.Sprintf("cluster-%s", cluster.DBClusterIdentifier)
	cluster.NetworkType = "IPV4"
	cluster.EarliestRestorableTime = cluster.ClusterCreateTime
	cluster.LatestRestorableTime = cluster.ClusterCreateTime

	if len(in.VpcSecurityGroupIds) > 0 {
		if _, err := s.resolveSecurityGroups(ctx, in.Region, in.VpcSecurityGroupIds); err != nil {
			return nil, translateStoreError(err)
		}
		cluster.VpcSecurityGroupIds = in.VpcSecurityGroupIds
	}
	if in.HasServerlessV2Scaling {
		minCap := in.ServerlessV2MinCapacity
		maxCap := in.ServerlessV2MaxCapacity
		if minCap < 0.5 || minCap > 128 || maxCap < 1 || maxCap > 256 || minCap >= maxCap {
			return nil, awserrors.NewAWSError("InvalidParameterValue", "ServerlessV2ScalingConfiguration: MinCapacity must be 0.5-128, MaxCapacity 1-256, and MinCapacity < MaxCapacity", http.StatusBadRequest)
		}
		cluster.ServerlessV2ScalingConfiguration = &neptunestore.ServerlessV2ScalingConfiguration{
			MinCapacity: minCap,
			MaxCapacity: maxCap,
		}
	}

	if err := store.CreateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	if cluster.GlobalClusterIdentifier != "" {
		if gc, err := store.GetGlobalCluster(cluster.GlobalClusterIdentifier); err == nil {
			isWriter := len(gc.GlobalClusterMembers) == 0
			if !isWriter {
				hasWriter := false
				for _, m := range gc.GlobalClusterMembers {
					if m.IsWriter {
						hasWriter = true
						break
					}
				}
				if !hasWriter {
					isWriter = true
				}
			}
			gc.GlobalClusterMembers = append(gc.GlobalClusterMembers, neptunestore.GlobalClusterMember{
				DBClusterArn:            cluster.DBClusterArn,
				IsWriter:                isWriter,
				GlobalClusterIdentifier: gc.GlobalClusterIdentifier,
			})
			if err := store.UpdateGlobalCluster(gc); err != nil {
				logs.Warn("failed to register cluster as global cluster member", logs.String("cluster", cluster.DBClusterIdentifier), logs.Err(err))
			}
		}
	}

	var enginePort int
	if eng := s.engineFor(cluster.Engine); eng != nil {
		if port, err := eng.Open(in.Region, cluster.DBClusterIdentifier); err != nil {
			logs.Warn("failed to open cluster engine", logs.String("cluster", cluster.DBClusterIdentifier), logs.Err(err))
		} else {
			enginePort = port
		}
	}

	s.setClusterEndpoint(store, cluster, enginePort)

	if len(in.Tags) > 0 {
		if err := store.AddTags(cluster.DBClusterArn, in.Tags); err != nil {
			if eng := s.engineFor(cluster.Engine); eng != nil {
				eng.Close(cluster.DBClusterIdentifier)
			}
			removeClusterFromGlobal(store, cluster)
			store.DeleteCluster(cluster.DBClusterIdentifier)
			return nil, awserrors.NewAWSError("InvalidParameterValue", fmt.Sprintf("Failed to tag cluster: %v", err), http.StatusBadRequest)
		}
	}

	recordEvent(store, "db-cluster", cluster.DBClusterIdentifier, cluster.DBClusterArn,
		fmt.Sprintf("DB cluster %s created", cluster.DBClusterIdentifier), []string{"creation"})

	cluster.Status = "available"
	if err := store.UpdateCluster(cluster); err != nil {
		logs.Warn("failed to transition cluster to available", logs.String("cluster", cluster.DBClusterIdentifier), logs.Err(err))
	}
	s.scheduleTransition(in.Region, 500*time.Millisecond, func(st neptunestore.NeptuneStoreInterface) error {
		c, err := st.GetCluster(cluster.DBClusterIdentifier)
		if err != nil || c.Status != "creating" {
			return nil
		}
		c.Status = "available"
		return st.UpdateCluster(c)
	})

	return map[string]interface{}{
		"DBCluster": enrichClusterWithTags(store, cluster),
	}, nil
}

// deleteDBClusterCore validates and executes cluster deletion, including the
// final-snapshot rule, cascade cleanup and engine shutdown.
func (s *NeptuneService) deleteDBClusterCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *DeleteDBClusterInput) (interface{}, error) {
	cluster, err := rdssvc.ValidateDeleteClusterParams(store, in.DeleteClusterParams)
	if err != nil {
		return nil, neptuneTranslateError(err)
	}

	skipFinal := in.SkipFinalSnapshot
	finalSnapshotID := in.FinalDBSnapshotIdentifier
	if !skipFinal && finalSnapshotID == "" {
		return nil, awserrors.NewAWSError("InvalidParameterCombination", "SkipFinalSnapshot must be true or FinalDBSnapshotIdentifier must be specified", http.StatusBadRequest)
	}
	if skipFinal && finalSnapshotID != "" {
		return nil, awserrors.NewAWSError("InvalidParameterCombination", "Cannot specify both SkipFinalSnapshot and FinalDBSnapshotIdentifier", http.StatusBadRequest)
	}

	cluster.Status = "deleting"
	if err := store.UpdateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	if !skipFinal {
		snapshot := rdssvc.BuildFinalSnapshot(cluster, finalSnapshotID, in.AccountID, in.Region)
		if err := store.CreateSnapshot(snapshot); err != nil {
			cluster.Status = "available"
			store.UpdateCluster(cluster)
			return nil, translateStoreError(err)
		}
	}

	if err := store.DeleteCluster(cluster.DBClusterIdentifier); err != nil {
		cluster.Status = "available"
		if rbErr := store.UpdateCluster(cluster); rbErr != nil {
			logs.Warn("failed to roll back cluster status after delete failure", logs.String("cluster", cluster.DBClusterIdentifier), logs.Err(rbErr))
		}
		if !skipFinal {
			if delErr := store.DeleteSnapshot(finalSnapshotID); delErr != nil {
				logs.Warn("failed to clean up orphaned snapshot after delete failure", logs.String("snapshot", finalSnapshotID), logs.Err(delErr))
			}
		}
		return nil, translateStoreError(err)
	}

	cascadeDeleteClusterResources(store, cluster)
	removeClusterFromGlobal(store, cluster)

	if eng := s.engineFor(cluster.Engine); eng != nil {
		if err := eng.Close(cluster.DBClusterIdentifier); err != nil {
			logs.Warn("failed to close cluster engine", logs.String("cluster", cluster.DBClusterIdentifier), logs.Err(err))
		}
	}

	recordEvent(store, "db-cluster", cluster.DBClusterIdentifier, cluster.DBClusterArn,
		fmt.Sprintf("DB cluster %s deleted", cluster.DBClusterIdentifier), []string{"deletion"})

	return map[string]interface{}{
		"DBCluster": enrichClusterWithTags(store, cluster),
	}, nil
}

// modifyDBClusterCore applies presence-based member updates to a cluster and
// handles the rename path with resource reparenting.
func (s *NeptuneService) modifyDBClusterCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *ModifyDBClusterInput) (interface{}, error) {
	id := in.DBClusterIdentifier
	if id == "" {
		return nil, awserrors.NewMissingParameter("DBClusterIdentifier is required")
	}

	cluster, err := store.GetCluster(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	if v := in.EngineVersion; v != "" {
		cluster.EngineVersion = v
	}
	if v := in.DBClusterParameterGroupName; v != "" {
		// Validate referenced parameter group exists before assigning.
		if _, err := store.GetClusterParameterGroup(v); err != nil {
			return nil, awserrors.NewAWSError("DBClusterParameterGroupNotFoundFault", fmt.Sprintf("DB Cluster Parameter Group not found: %s", v), http.StatusNotFound)
		}
		cluster.DBClusterParameterGroupName = v
	}
	if v := in.Port; v > 0 {
		// Validate Port range on modify.
		if err := validatePort(v); err != nil {
			return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
		}
		cluster.Port = v
		if cluster.Endpoint != nil {
			cluster.Endpoint.Port = v
		}
	}
	if v := in.BackupRetentionPeriod; v > 0 {
		// Validate BackupRetentionPeriod range on modify.
		if err := validateBackupRetentionPeriod(v); err != nil {
			return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
		}
		cluster.BackupRetentionPeriod = v
	}
	if v := in.PreferredBackupWindow; v != "" {
		cluster.PreferredBackupWindow = v
	}
	if v := in.PreferredMaintenanceWindow; v != "" {
		cluster.PreferredMaintenanceWindow = v
	}
	if v := in.StorageType; v != "" {
		cluster.StorageType = v
	}
	// Accept MasterUserPassword on modify and store as bcrypt hash.
	if pwd := in.MasterUserPassword; pwd != "" {
		hash, hashErr := hashMasterPassword(pwd)
		if hashErr != nil {
			return nil, awserrors.NewAWSError("InvalidParameterValue", hashErr.Error(), http.StatusBadRequest)
		}
		cluster.MasterUserPasswordHash = hash
	}
	// Handle NetworkType on modify (Smithy ModifyDBClusterMessage member).
	if nt := in.NetworkType; nt != "" {
		cluster.NetworkType = nt
	}
	// Handle ServerlessV2ScalingConfiguration on modify.
	if in.HasServerlessV2Scaling {
		minCap := in.ServerlessV2MinCapacity
		maxCap := in.ServerlessV2MaxCapacity
		if minCap < 0.5 || minCap > 128 || maxCap < 1 || maxCap > 256 || minCap >= maxCap {
			return nil, awserrors.NewAWSError("InvalidParameterValue", "ServerlessV2ScalingConfiguration: MinCapacity must be 0.5-128, MaxCapacity 1-256, and MinCapacity < MaxCapacity", http.StatusBadRequest)
		}
		cluster.ServerlessV2ScalingConfiguration = &neptunestore.ServerlessV2ScalingConfiguration{
			MinCapacity: minCap,
			MaxCapacity: maxCap,
		}
	}
	if in.HasDeletionProtection {
		cluster.DeletionProtection = in.DeletionProtection
	}
	if in.HasEnableIAMDatabaseAuth {
		cluster.IAMDatabaseAuthenticationEnabled = in.EnableIAMDatabaseAuthentication
	}
	if len(in.VpcSecurityGroupIds) > 0 {
		if _, err := s.resolveSecurityGroups(ctx, in.Region, in.VpcSecurityGroupIds); err != nil {
			return nil, translateStoreError(err)
		}
		cluster.VpcSecurityGroupIds = in.VpcSecurityGroupIds
	}
	if len(in.EnabledCloudwatchLogsExports) > 0 {
		cluster.EnabledCloudwatchLogsExports = in.EnabledCloudwatchLogsExports
	}

	if err := store.UpdateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	if newPort := in.Port; newPort > 0 {
		s.setClusterEndpoint(store, cluster, newPort)
	}

	newID := in.NewDBClusterIdentifier
	if newID != "" && newID != id {
		oldArn := cluster.DBClusterArn
		oldID := cluster.DBClusterIdentifier
		cluster.DBClusterIdentifier = newID
		cluster.DBClusterArn = arnutil.NewARNBuilder(cluster.AccountID, cluster.Region).RDS().Cluster(newID)
		if err := store.CreateCluster(cluster); err != nil {
			cluster.DBClusterIdentifier = oldID
			cluster.DBClusterArn = oldArn
			return nil, translateStoreError(err)
		}
		if err := reparentClusterResources(store, oldArn, cluster.DBClusterArn, oldID, newID); err != nil {
			store.DeleteCluster(newID)
			cluster.DBClusterIdentifier = oldID
			cluster.DBClusterArn = oldArn
			return nil, awserrors.NewAWSError("InvalidDBClusterStateFault",
				fmt.Sprintf("cluster rename failed during resource reparenting: %v", err), http.StatusBadRequest)
		}
		if err := store.DeleteCluster(oldID); err != nil {
			logs.Error("failed to delete old cluster record after successful rename",
				logs.String("oldID", oldID), logs.String("newID", newID), logs.Err(err))
		}
	}

	return map[string]interface{}{
		"DBCluster": enrichClusterWithTags(store, cluster),
	}, nil
}

// startDBClusterCore transitions a stopped cluster back to available and
// reopens the engine.
func (s *NeptuneService) startDBClusterCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *StartDBClusterInput) (interface{}, error) {
	id := in.DBClusterIdentifier
	if id == "" {
		return nil, awserrors.NewMissingParameter("DBClusterIdentifier is required")
	}

	cluster, err := store.GetCluster(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	if cluster.Status != "stopped" {
		return nil, awserrors.NewAWSError("InvalidDBClusterStateFault", fmt.Sprintf("DBCluster %s is not in stopped state (current: %s)", id, cluster.Status), http.StatusBadRequest)
	}

	cluster.Status = "available"
	if err := store.UpdateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	var enginePort int
	if eng := s.engineFor(cluster.Engine); eng != nil {
		if port, err := eng.Open(in.Region, id); err != nil {
			logs.Warn("failed to open cluster engine on start", logs.String("cluster", id), logs.Err(err))
		} else {
			enginePort = port
		}
	}

	s.setClusterEndpoint(store, cluster, enginePort)

	return map[string]interface{}{
		"DBCluster": enrichClusterWithTags(store, cluster),
	}, nil
}

// stopDBClusterCore stops a running cluster and closes the engine.
func (s *NeptuneService) stopDBClusterCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *StopDBClusterInput) (interface{}, error) {
	id := in.DBClusterIdentifier
	if id == "" {
		return nil, awserrors.NewMissingParameter("DBClusterIdentifier is required")
	}

	cluster, err := store.GetCluster(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	if cluster.Status != "available" {
		return nil, awserrors.NewAWSError("InvalidDBClusterStateFault", fmt.Sprintf("DBCluster %s is not in available state (current: %s)", id, cluster.Status), http.StatusBadRequest)
	}

	if eng := s.engineFor(cluster.Engine); eng != nil {
		if err := eng.Close(id); err != nil {
			logs.Warn("failed to close cluster engine on stop", logs.String("cluster", id), logs.Err(err))
		}
	}

	cluster.Status = "stopped"
	if err := store.UpdateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBCluster": enrichClusterWithTags(store, cluster),
	}, nil
}

// failoverDBClusterCore re-asserts the writer role on a cluster; Neptune's
// single-writer topology makes this a state re-validation.
func (s *NeptuneService) failoverDBClusterCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *FailoverDBClusterInput) (interface{}, error) {
	id := in.DBClusterIdentifier
	if id == "" {
		return nil, awserrors.NewMissingParameter("DBClusterIdentifier is required")
	}

	cluster, err := store.GetCluster(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	// FailoverDBCluster requires the cluster to be in 'available' state.
	if cluster.Status != "available" {
		return nil, awserrors.NewAWSError("InvalidDBClusterStateFault", fmt.Sprintf("DBCluster %s is not in available state (current: %s)", id, cluster.Status), http.StatusBadRequest)
	}

	cluster.Status = "available"
	if err := store.UpdateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBCluster": enrichClusterWithTags(store, cluster),
	}, nil
}

// addRoleToDBClusterCore associates an IAM role with a cluster.
func (s *NeptuneService) addRoleToDBClusterCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *AddRoleToDBClusterInput) (interface{}, error) {
	id := in.DBClusterIdentifier
	if id == "" {
		return nil, awserrors.NewMissingParameter("DBClusterIdentifier is required")
	}
	roleArn := in.RoleArn
	if roleArn == "" {
		return nil, awserrors.NewMissingParameter("RoleArn is required")
	}
	// Validate IAM role ARN format.
	if !isValidIAMRoleArn(roleArn) {
		return nil, awserrors.NewAWSError("InvalidParameterValue", fmt.Sprintf("Invalid IAM role ARN: %s", roleArn), http.StatusBadRequest)
	}

	cluster, err := store.GetCluster(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	for _, r := range cluster.AssociatedRoles {
		if r.RoleArn == roleArn {
			return nil, awserrors.NewAWSError("DBClusterRoleAlreadyExistsFault", fmt.Sprintf("IAM role %s is already associated with cluster %s", roleArn, id), http.StatusBadRequest)
		}
	}
	cluster.AssociatedRoles = append(cluster.AssociatedRoles, neptunestore.DBClusterRole{
		RoleArn:     roleArn,
		FeatureName: in.FeatureName,
		Status:      "ACTIVE",
	})
	if err := store.UpdateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{}, nil
}

// removeRoleFromDBClusterCore disassociates an IAM role from a cluster.
func (s *NeptuneService) removeRoleFromDBClusterCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *RemoveRoleFromDBClusterInput) (interface{}, error) {
	id := in.DBClusterIdentifier
	if id == "" {
		return nil, awserrors.NewMissingParameter("DBClusterIdentifier is required")
	}
	roleArn := in.RoleArn
	if roleArn == "" {
		return nil, awserrors.NewMissingParameter("RoleArn is required")
	}

	cluster, err := store.GetCluster(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	found := false
	filtered := make([]neptunestore.DBClusterRole, 0, len(cluster.AssociatedRoles))
	for _, r := range cluster.AssociatedRoles {
		if r.RoleArn == roleArn && (in.FeatureName == "" || r.FeatureName == in.FeatureName) {
			found = true
			continue
		}
		filtered = append(filtered, r)
	}
	if !found {
		return nil, awserrors.NewAWSError("DBClusterRoleNotFoundFault", fmt.Sprintf("role %s is not associated with cluster %s", roleArn, id), http.StatusBadRequest)
	}
	cluster.AssociatedRoles = filtered
	if err := store.UpdateCluster(cluster); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{}, nil
}

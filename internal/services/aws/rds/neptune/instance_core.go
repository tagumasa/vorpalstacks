package neptune

import (
	"context"
	"fmt"
	"net/http"

	awserrors "vorpalstacks/internal/common/errors"
	types "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	rdssvc "vorpalstacks/internal/services/aws/rds"
	neptunestore "vorpalstacks/internal/store/aws/rds/neptune"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// CreateDBInstanceInput carries the wire-parsed CreateDBInstance request. The
// embedded CreateInstanceParams holds the members shared with the admin-plane
// instance Core.
type CreateDBInstanceInput struct {
	rdssvc.CreateInstanceParams
	MasterUserPassword string
	Tags               []types.Tag
}

// DeleteDBInstanceInput carries the wire-parsed DeleteDBInstance request.
type DeleteDBInstanceInput struct {
	rdssvc.DeleteInstanceParams
}

// ModifyDBInstanceInput carries the wire-parsed ModifyDBInstance request. The
// Has* members preserve the wire presence of optional members so an omitted
// boolean keeps its stored value instead of resetting it.
type ModifyDBInstanceInput struct {
	DBInstanceIdentifier            string
	DBInstanceClass                 string
	EngineVersion                   string
	DBParameterGroupName            string
	PreferredMaintenanceWindow      string
	HasPubliclyAccessible           bool
	PubliclyAccessible              bool
	HasAutoMinorVersionUpgrade      bool
	AutoMinorVersionUpgrade         bool
	HasEnableIAMDatabaseAuth        bool
	EnableIAMDatabaseAuthentication bool
	HasCopyTagsToSnapshot           bool
	CopyTagsToSnapshot              bool
	NewDBInstanceIdentifier         string
}

// RebootDBInstanceInput carries the wire-parsed RebootDBInstance request.
type RebootDBInstanceInput struct {
	DBInstanceIdentifier string
}

// createDBInstanceCore validates and persists a new DB instance, registers
// cluster membership, resolves the instance endpoint and records the
// creation event.
func (s *NeptuneService) createDBInstanceCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *CreateDBInstanceInput) (interface{}, error) {
	createParams := in.CreateInstanceParams
	clusterID := createParams.DBClusterIdentifier
	if clusterID != "" {
		if _, err := store.GetCluster(clusterID); err != nil {
			return nil, awserrors.NewAWSError("DBClusterNotFoundFault", fmt.Sprintf("DBCluster %s not found", clusterID), http.StatusNotFound)
		}
	}

	if err := rdssvc.ValidateCreateInstanceParams(store, createParams); err != nil {
		return nil, neptuneTranslateError(err)
	}

	masterPasswordHash, err := hashMasterPassword(in.MasterUserPassword)
	if err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}

	instance := rdssvc.BuildInstance(createParams)
	instance.DBInstanceStatus = "available"
	instance.MasterUserPasswordHash = masterPasswordHash

	if err := store.CreateInstance(instance); err != nil {
		return nil, translateStoreError(err)
	}

	if clusterID != "" {
		if cluster, err := store.GetCluster(clusterID); err == nil {
			isWriter := len(cluster.DBClusterMembers) == 0
			cluster.DBClusterMembers = append(cluster.DBClusterMembers, neptunestore.DBClusterMember{
				DBInstanceIdentifier: instance.DBInstanceIdentifier,
				IsClusterWriter:      isWriter,
				PromotionTier:        0,
			})
			if err := store.UpdateCluster(cluster); err != nil {
				logs.Warn("failed to add instance to cluster members", logs.String("instance", instance.DBInstanceIdentifier), logs.Err(err))
			}
		}
	}

	if clusterID != "" {
		if s.porter != nil {
			if port, err := s.porter.GetPort(clusterID); err == nil && port > 0 {
				instance.Endpoint = &neptunestore.Endpoint{
					Address: s.endpointAddressFor(instance.DBInstanceIdentifier, instance.Engine),
					Port:    port,
				}
				if err := store.UpdateInstance(instance); err != nil {
					logs.Warn("failed to persist instance endpoint", logs.String("instance", instance.DBInstanceIdentifier), logs.Err(err))
				}
			}
		}
	} else {
		engineType := instance.Engine
		if engineType == "" {
			engineType = "neptune"
		}
		if eng := s.engineFor(engineType); eng != nil {
			if port, err := eng.Open(createParams.Region, instance.DBInstanceIdentifier); err != nil {
				logs.Warn("failed to open instance engine", logs.String("instance", instance.DBInstanceIdentifier), logs.Err(err))
			} else {
				instance.Endpoint = &neptunestore.Endpoint{
					Address: s.endpointAddressFor(instance.DBInstanceIdentifier, engineType),
					Port:    port,
				}
				if err := store.UpdateInstance(instance); err != nil {
					logs.Warn("failed to persist instance endpoint", logs.String("instance", instance.DBInstanceIdentifier), logs.Err(err))
				}
			}
		}
	}

	recordEvent(store, "db-instance", instance.DBInstanceIdentifier, instance.DBInstanceArn,
		fmt.Sprintf("DB instance %s created", instance.DBInstanceIdentifier), []string{"creation"})

	if len(in.Tags) > 0 {
		if err := store.AddTags(instance.DBInstanceArn, in.Tags); err != nil {
			if instance.DBClusterIdentifier == "" {
				engineType := instance.Engine
				if engineType == "" {
					engineType = "neptune"
				}
				if eng := s.engineFor(engineType); eng != nil {
					eng.Close(instance.DBInstanceIdentifier)
				}
			}
			store.DeleteInstance(instance.DBInstanceIdentifier)
			return nil, awserrors.NewAWSError("InvalidParameterValue", fmt.Sprintf("Failed to tag instance: %v", err), http.StatusBadRequest)
		}
	}

	return map[string]interface{}{
		"DBInstance": instance,
	}, nil
}

// deleteDBInstanceCore validates and executes instance deletion, including
// cluster-membership cleanup, engine shutdown and tag removal.
func (s *NeptuneService) deleteDBInstanceCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *DeleteDBInstanceInput) (interface{}, error) {
	instance, err := rdssvc.ValidateDeleteInstanceParams(store, in.DeleteInstanceParams)
	if err != nil {
		return nil, neptuneTranslateError(err)
	}

	instance.DBInstanceStatus = "deleting"
	if err := store.UpdateInstance(instance); err != nil {
		return nil, translateStoreError(err)
	}

	if instance.DBClusterIdentifier == "" {
		engineType := instance.Engine
		if engineType == "" {
			engineType = "neptune"
		}
		if eng := s.engineFor(engineType); eng != nil {
			if err := eng.Close(instance.DBInstanceIdentifier); err != nil {
				logs.Warn("failed to close instance engine on delete", logs.String("instance", instance.DBInstanceIdentifier), logs.Err(err))
			}
		}
	}

	if err := store.DeleteInstance(instance.DBInstanceIdentifier); err != nil {
		return nil, translateStoreError(err)
	}

	if instance.DBClusterIdentifier != "" {
		if cluster, err := store.GetCluster(instance.DBClusterIdentifier); err == nil {
			filtered := make([]neptunestore.DBClusterMember, 0, len(cluster.DBClusterMembers))
			for _, mem := range cluster.DBClusterMembers {
				if mem.DBInstanceIdentifier != instance.DBInstanceIdentifier {
					filtered = append(filtered, mem)
				}
			}
			cluster.DBClusterMembers = filtered
			if err := store.UpdateCluster(cluster); err != nil {
				logs.Warn("failed to remove instance from cluster members", logs.String("instance", instance.DBInstanceIdentifier), logs.Err(err))
			}
		}
	}

	removeTagsForResource(store, instance.DBInstanceArn)

	recordEvent(store, "db-instance", instance.DBInstanceIdentifier, instance.DBInstanceArn,
		fmt.Sprintf("DB instance %s deleted", instance.DBInstanceIdentifier), []string{"deletion"})

	return map[string]interface{}{
		"DBInstance": instance,
	}, nil
}

// modifyDBInstanceCore applies member updates to an instance and handles the
// rename path with tag copy and cluster-membership repointing.
func (s *NeptuneService) modifyDBInstanceCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *ModifyDBInstanceInput) (interface{}, error) {
	id := in.DBInstanceIdentifier
	if id == "" {
		return nil, awserrors.NewMissingParameter("DBInstanceIdentifier is required")
	}

	instance, err := store.GetInstance(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	if v := in.DBInstanceClass; v != "" {
		instance.DBInstanceClass = v
	}
	if v := in.EngineVersion; v != "" {
		instance.EngineVersion = v
	}
	if v := in.DBParameterGroupName; v != "" {
		instance.DBParameterGroupName = v
	}
	if v := in.PreferredMaintenanceWindow; v != "" {
		instance.PreferredMaintenanceWindow = v
	}
	if in.HasPubliclyAccessible {
		instance.PubliclyAccessible = in.PubliclyAccessible
	}
	if in.HasAutoMinorVersionUpgrade {
		instance.AutoMinorVersionUpgrade = in.AutoMinorVersionUpgrade
	}
	if in.HasEnableIAMDatabaseAuth {
		instance.IAMDatabaseAuthenticationEnabled = in.EnableIAMDatabaseAuthentication
	}
	if in.HasCopyTagsToSnapshot {
		instance.CopyTagsToSnapshot = in.CopyTagsToSnapshot
	}

	if err := store.UpdateInstance(instance); err != nil {
		return nil, translateStoreError(err)
	}

	// Support NewDBInstanceIdentifier (instance rename).
	// Follows the same create-new + delete-old pattern as ModifyDBCluster.
	if newID := in.NewDBInstanceIdentifier; newID != "" && newID != id {
		oldArn := instance.DBInstanceArn
		instance.DBInstanceIdentifier = newID
		instance.DBInstanceArn = arnutil.NewARNBuilder(instance.AccountID, instance.Region).RDS().DBInstance(newID)
		if err := store.CreateInstance(instance); err != nil {
			instance.DBInstanceIdentifier = id
			instance.DBInstanceArn = oldArn
			return nil, translateStoreError(err)
		}
		// Copy tags from old to new ARN.
		if tags, err := store.GetTags(oldArn); err == nil && len(tags) > 0 {
			store.AddTags(instance.DBInstanceArn, tags)
		}
		if err := store.DeleteInstance(id); err != nil {
			logs.Warn("failed to delete old instance after rename", logs.String("oldID", id), logs.Err(err))
		}
		removeTagsForResource(store, oldArn)
		// Update DBClusterMembers if the instance belongs to a cluster.
		if instance.DBClusterIdentifier != "" {
			if cluster, err := store.GetCluster(instance.DBClusterIdentifier); err == nil {
				for i, mem := range cluster.DBClusterMembers {
					if mem.DBInstanceIdentifier == id {
						cluster.DBClusterMembers[i].DBInstanceIdentifier = newID
					}
				}
				store.UpdateCluster(cluster)
			}
		}
	}

	return map[string]interface{}{
		"DBInstance": instance,
	}, nil
}

// rebootDBInstanceCore reboots an available instance.
func (s *NeptuneService) rebootDBInstanceCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *RebootDBInstanceInput) (interface{}, error) {
	id := in.DBInstanceIdentifier
	if id == "" {
		return nil, awserrors.NewMissingParameter("DBInstanceIdentifier is required")
	}

	instance, err := store.GetInstance(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	if instance.DBInstanceStatus != "available" {
		return nil, awserrors.NewAWSError("InvalidDBInstanceStateFault", fmt.Sprintf("instance %s is not in available state", id), http.StatusBadRequest)
	}

	instance.DBInstanceStatus = "available"
	if err := store.UpdateInstance(instance); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"DBInstance": instance,
	}, nil
}

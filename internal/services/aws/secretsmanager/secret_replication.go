package secretsmanager

import (
	"context"
	"fmt"
	"net/http"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
)

var (
	// ErrReplicationNotConfigured indicates that the storage manager is not available for cross-region operations.
	ErrReplicationNotConfigured = awserrors.NewAWSError("InvalidRequestException", "Cross-region replication requires a storage manager", http.StatusBadRequest)
	// ErrInvalidReplicationRegion indicates the specified region is not valid.
	ErrInvalidReplicationRegion = awserrors.NewAWSError("InvalidParameterException", "Replication region is not valid", http.StatusBadRequest)
	// ErrSecretAlreadyReplicating indicates the secret is already being replicated to the specified region.
	ErrSecretAlreadyReplicating = awserrors.NewAWSError("InvalidRequestException", "Secret is already being replicated to the specified region", http.StatusBadRequest)
	// ErrNoReplicationConfigured indicates the secret has no replication configured.
	ErrNoReplicationConfigured = awserrors.NewAWSError("InvalidRequestException", "You can't perform this operation on the secret because it's not configured with replica regions", http.StatusBadRequest)
	// ErrReplicaNotFound indicates the specified replica region is not found.
	ErrReplicaNotFound = awserrors.NewAWSError("ResourceNotFoundException", "Replica region not found in replication configuration", http.StatusNotFound)
)

// replicateSecretToRegions is the shared replication engine used by both
// ReplicateSecretToRegions and CreateSecret (when AddReplicaRegions is
// provided inline). It replicates the secret and its versions to each
// target region, tracking per-replica success/failure accurately.
//
// Parameters:
//   - store: the primary-region store
//   - secret: the primary secret (mutated: ReplicationStatus updated)
//   - regions: the replica regions to create
//   - forceOverwrite: when true, overwrite existing replicas in target regions
//   - primaryRegion: the source region for replica-side metadata
func (s *SecretsManagerService) replicateSecretToRegions(
	store secretsmanagerstore.SecretStoreInterface,
	secret *secretsmanagerstore.Secret,
	regions []replicaRegion,
	forceOverwrite bool,
	primaryRegion string,
) {
	// Read the current version's value for replication. SecretString and
	// SecretBinary on the Secret struct are transient (json:"-") and are
	// not populated when the secret was deserialised from the store.
	var srcSecretString string
	var srcSecretBinary []byte
	if secret.CurrentVersion != "" {
		if currentVer, verErr := store.GetSecretVersion(secret.Name, secret.CurrentVersion); verErr == nil {
			srcSecretString = currentVer.SecretString
			srcSecretBinary = currentVer.SecretBinary
		}
	}

	for _, replicaRegion := range regions {
		// Check for existing replica. When forceOverwrite is true,
		// delete the old replica first; otherwise skip with a warning.
		alreadyExists := false
		for _, existing := range secret.ReplicationStatus {
			if existing.Region == replicaRegion.Region {
				alreadyExists = true
				break
			}
		}
		if alreadyExists && !forceOverwrite {
			// Already replicating — skip silently (caller handles errors).
			continue
		}

		regionStorage, err := s.storageManager.GetStorage(replicaRegion.Region)
		if err != nil {
			secret.ReplicationStatus = upsertReplicationStatus(secret.ReplicationStatus, secretsmanagerstore.ReplicationStatus{
				Region:        replicaRegion.Region,
				KmsKeyId:      replicaRegion.KmsKeyId,
				Status:        "Failed",
				StatusMessage: fmt.Sprintf("Region %s is not available: %v", replicaRegion.Region, err),
			})
			continue
		}

		replicaStore := secretsmanagerstore.NewSecretStore(regionStorage, s.accountID, replicaRegion.Region)

		// When forceOverwrite is true, delete any existing secret in
		// the target region before creating the replica.  This covers
		// both existing replicas (tracked in ReplicationStatus) and
		// unrelated secrets with the same name in the target region
		// CreateSecret inline replication where ReplicationStatus
		// is still empty, so alreadyExists is false).
		if forceOverwrite {
			_ = replicaStore.DeleteSecret(secret.Name)
		}

		replica := secretsmanagerstore.NewSecret(secret.Name)
		replica.Description = secret.Description
		replica.PrimaryRegion = primaryRegion
		replica.KmsKeyId = replicaRegion.KmsKeyId
		replica.SecretString = srcSecretString
		replica.SecretBinary = srcSecretBinary
		replica.Tags = make(map[string]string)
		for k, v := range secret.Tags {
			replica.Tags[k] = v
		}

		if _, err := replicaStore.CreateSecret(replica); err != nil {
			secret.ReplicationStatus = upsertReplicationStatus(secret.ReplicationStatus, secretsmanagerstore.ReplicationStatus{
				Region:        replicaRegion.Region,
				KmsKeyId:      replicaRegion.KmsKeyId,
				Status:        "Failed",
				StatusMessage: err.Error(),
			})
			continue
		}

		// Track version sync failures accurately. If any version
		// fails to sync, the replica status reflects "Failed" instead
		// of the misleading "InSync".
		syncFailures := 0
		if (srcSecretString != "" || len(srcSecretBinary) > 0) && replica.CurrentVersion != "" {
			srcVersions, err := store.ListSecretVersions(secret.Name)
			if err == nil && len(srcVersions) > 0 {
				for _, v := range srcVersions {
					versionKey := fmt.Sprintf("%s/%s/%s", s.accountID, secret.Name, v.VersionId)
					var srcVersion secretsmanagerstore.SecretVersion
					if srcErr := store.GetBaseStore().Get(versionKey, &srcVersion); srcErr == nil {
						replicaVersion := secretsmanagerstore.NewSecretVersion(srcVersion.VersionId)
						replicaVersion.SecretName = secret.Name
						replicaVersion.SecretString = srcVersion.SecretString
						replicaVersion.SecretBinary = srcVersion.SecretBinary
						replicaVersion.VersionStages = srcVersion.VersionStages
						if err := replicaStore.CreateVersionDirect(secret.Name, replicaVersion); err != nil {
							syncFailures++
							logs.Warn("Failed to create version on replica",
								logs.String("region", replicaRegion.Region),
								logs.String("version", v.VersionId),
								logs.Err(err))
						}
					}
				}
			}
		}

		replica.ReplicationStatus = []secretsmanagerstore.ReplicationStatus{
			{
				Region:           primaryRegion,
				KmsKeyId:         "",
				Status:           "InSync",
				LastAccessedDate: time.Now().UTC(),
			},
		}

		if err := replicaStore.UpdateSecretMetadata(replica); err != nil {
			logs.Error("Failed to set replication status on replica", logs.String("region", replicaRegion.Region), logs.Err(err))
		}

		finalStatus := "InSync"
		statusMessage := ""
		if syncFailures > 0 {
			finalStatus = "Failed"
			statusMessage = fmt.Sprintf("%d version(s) failed to sync", syncFailures)
		}
		secret.ReplicationStatus = upsertReplicationStatus(secret.ReplicationStatus, secretsmanagerstore.ReplicationStatus{
			Region:           replicaRegion.Region,
			KmsKeyId:         replicaRegion.KmsKeyId,
			Status:           finalStatus,
			StatusMessage:    statusMessage,
			LastAccessedDate: time.Now().UTC(),
		})
	}
}

// upsertReplicationStatus appends or replaces a replication status entry
// for the given region.
func upsertReplicationStatus(statuses []secretsmanagerstore.ReplicationStatus, rs secretsmanagerstore.ReplicationStatus) []secretsmanagerstore.ReplicationStatus {
	for i, existing := range statuses {
		if existing.Region == rs.Region {
			statuses[i] = rs
			return statuses
		}
	}
	return append(statuses, rs)
}

// ReplicateSecretToRegions replicates a secret to one or more regions.
func (s *SecretsManagerService) ReplicateSecretToRegions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if err := validateSecretId(secretId); err != nil {
		return nil, err
	}

	if s.storageManager == nil {
		return nil, ErrReplicationNotConfigured
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	secret, err := s.resolveSecret(reqCtx, secretId)
	if err != nil {
		return nil, err
	}

	addReplicaRegions, err := parseReplicaRegions(req.Parameters, "AddReplicaRegions")
	if err != nil {
		return nil, err
	}

	if len(addReplicaRegions) == 0 {
		return nil, awserrors.NewAWSError("InvalidParameterException", "AddReplicaRegions must not be empty", http.StatusBadRequest)
	}

	// ForceOverwriteReplicaSecret controls whether existing replicas
	// are overwritten in the target regions.
	forceOverwrite := request.GetBoolParam(req.Parameters, "ForceOverwriteReplicaSecret")

	// Check for duplicates when not force-overwriting.
	if !forceOverwrite {
		for _, existing := range secret.ReplicationStatus {
			for _, newRegion := range addReplicaRegions {
				if existing.Region == newRegion.Region {
					return nil, ErrSecretAlreadyReplicating
				}
			}
		}
	}

	s.replicateSecretToRegions(store, secret, addReplicaRegions, forceOverwrite, reqCtx.GetRegion())

	if err := store.UpdateSecretMetadata(secret); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"ARN":               secret.ARN,
		"Name":              secret.Name,
		"ReplicationStatus": buildReplicationStatusResponse(secret.ReplicationStatus),
	}, nil
}

// RemoveRegionsFromReplication removes replica regions from a secret.
func (s *SecretsManagerService) RemoveRegionsFromReplication(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if err := validateSecretId(secretId); err != nil {
		return nil, err
	}

	if s.storageManager == nil {
		return nil, ErrReplicationNotConfigured
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	secret, err := s.resolveSecret(reqCtx, secretId)
	if err != nil {
		return nil, err
	}

	removeRegions := request.GetStringList(req.Parameters, "RemoveReplicaRegions")
	if len(removeRegions) == 0 {
		return nil, awserrors.NewAWSError("InvalidParameterException", "RemoveReplicaRegions must not be empty", http.StatusBadRequest)
	}
	for _, r := range removeRegions {
		if err := validateRegion(r); err != nil {
			return nil, err
		}
	}

	remainingStatus := make([]secretsmanagerstore.ReplicationStatus, 0, len(secret.ReplicationStatus))
	removed := false

	for _, rs := range secret.ReplicationStatus {
		found := false
		for _, removeRegion := range removeRegions {
			if rs.Region == removeRegion {
				found = true
				removed = true
				regionStorage, storeErr := s.storageManager.GetStorage(rs.Region)
				if storeErr == nil {
					replicaStore := secretsmanagerstore.NewSecretStore(regionStorage, s.accountID, rs.Region)
					if delErr := replicaStore.DeleteSecret(secret.Name); delErr != nil {
						// Replica deletion failed. Keep the entry
						// with status "Failed" so the orphaned replica
						// is visible in DescribeSecret, rather than
						// silently disappearing from the replication
						// status.
						logs.Warn("Failed to delete replica secret; keeping status as Failed",
							logs.String("region", rs.Region), logs.Err(delErr))
						rs.Status = "Failed"
						rs.StatusMessage = fmt.Sprintf("Failed to delete replica: %v", delErr)
						remainingStatus = append(remainingStatus, rs)
					}
				} else {
					// Region storage unavailable — can't delete.
					logs.Warn("Region storage unavailable for replica deletion",
						logs.String("region", rs.Region), logs.Err(storeErr))
					rs.Status = "Failed"
					rs.StatusMessage = fmt.Sprintf("Region %s storage unavailable: %v", rs.Region, storeErr)
					remainingStatus = append(remainingStatus, rs)
				}
				break
			}
		}
		if !found {
			remainingStatus = append(remainingStatus, rs)
		}
	}

	if !removed {
		return nil, ErrReplicaNotFound
	}

	secret.ReplicationStatus = remainingStatus
	if err := store.UpdateSecretMetadata(secret); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"ARN":               secret.ARN,
		"Name":              secret.Name,
		"ReplicationStatus": buildReplicationStatusResponse(secret.ReplicationStatus),
	}, nil
}

// StopReplicationToReplica stops replication to a replica secret, promoting it to a standalone secret.
func (s *SecretsManagerService) StopReplicationToReplica(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if err := validateSecretId(secretId); err != nil {
		return nil, err
	}

	if s.storageManager == nil {
		return nil, ErrReplicationNotConfigured
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	secret, err := s.resolveSecret(reqCtx, secretId)
	if err != nil {
		return nil, err
	}

	if len(secret.ReplicationStatus) == 0 {
		return nil, ErrNoReplicationConfigured
	}
	// Note: AWS spec requires StopReplicationToReplica to be called from
	// the replica region only. On this edge platform all regions share a
	// single server instance, so the region distinction is not enforced.
	// PrimaryRegion is populated on replicas for future use.

	// Update the replica's LastAccessedDate to reflect the promotion
	// to a standalone secret, matching AWS behaviour.
	secret.ReplicationStatus = nil
	secret.LastAccessedDate = time.Now().UTC()
	if err := store.UpdateSecretMetadata(secret); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"ARN":  secret.ARN,
		"Name": secret.Name,
	}, nil
}

type replicaRegion struct {
	Region   string
	KmsKeyId string
}

func parseReplicaRegions(params map[string]interface{}, key string) ([]replicaRegion, error) {
	raw, ok := params[key]
	if !ok {
		return nil, nil
	}

	list, ok := raw.([]interface{})
	if !ok {
		return nil, ErrInvalidReplicationRegion
	}

	regions := make([]replicaRegion, 0, len(list))
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		region := ""
		if r, ok := m["Region"].(string); ok {
			region = r
		}
		kmsKey := ""
		if k, ok := m["KmsKeyId"].(string); ok {
			kmsKey = k
		}
		if region == "" {
			return nil, ErrInvalidReplicationRegion
		}
		if err := validateRegion(region); err != nil {
			return nil, err
		}
		if err := validateKmsKeyId(kmsKey); err != nil {
			return nil, err
		}
		regions = append(regions, replicaRegion{Region: region, KmsKeyId: kmsKey})
	}

	return regions, nil
}

func buildReplicationStatusResponse(statuses []secretsmanagerstore.ReplicationStatus) []interface{} {
	result := make([]interface{}, len(statuses))
	for i, rs := range statuses {
		entry := map[string]interface{}{
			"Region": rs.Region,
			"Status": rs.Status,
		}
		if rs.KmsKeyId != "" {
			entry["KmsKeyId"] = rs.KmsKeyId
		}
		if rs.StatusMessage != "" {
			entry["StatusMessage"] = rs.StatusMessage
		}
		if !rs.LastAccessedDate.IsZero() {
			entry["LastAccessedDate"] = rs.LastAccessedDate.Unix()
		}
		result[i] = entry
	}
	return result
}

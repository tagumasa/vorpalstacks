package secretsmanager

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	awserrors "vorpalstacks/internal/services/aws/common/errors"
	"vorpalstacks/internal/services/aws/common/request"
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

// ReplicateSecretToRegions replicates a secret to one or more regions.
// https://docs.aws.amazon.com/secretsmanager/latest/userguide/API_ReplicateSecretToRegions.html
func (s *SecretsManagerService) ReplicateSecretToRegions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if secretId == "" {
		return nil, awserrors.ErrMissingParameter
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

	for _, existing := range secret.ReplicationStatus {
		for _, newRegion := range addReplicaRegions {
			if existing.Region == newRegion.Region {
				return nil, ErrSecretAlreadyReplicating
			}
		}
	}

	for _, replicaRegion := range addReplicaRegions {
		regionStorage, err := s.storageManager.GetStorage(replicaRegion.Region)
		if err != nil {
			secret.ReplicationStatus = append(secret.ReplicationStatus, secretsmanagerstore.ReplicationStatus{
				Region:        replicaRegion.Region,
				KmsKeyId:      replicaRegion.KmsKeyId,
				Status:        "FAILED",
				StatusMessage: fmt.Sprintf("Region %s is not available: %v", replicaRegion.Region, err),
			})
			continue
		}

		replicaStore := secretsmanagerstore.NewSecretStore(regionStorage, s.accountID, replicaRegion.Region)

		replica := secretsmanagerstore.NewSecret(secret.Name)
		replica.Description = secret.Description
		replica.KmsKeyId = replicaRegion.KmsKeyId
		replica.SecretString = secret.SecretString
		replica.SecretBinary = secret.SecretBinary
		replica.Tags = make(map[string]string)
		for k, v := range secret.Tags {
			replica.Tags[k] = v
		}

		if _, err := replicaStore.CreateSecret(replica); err != nil {
			secret.ReplicationStatus = append(secret.ReplicationStatus, secretsmanagerstore.ReplicationStatus{
				Region:        replicaRegion.Region,
				KmsKeyId:      replicaRegion.KmsKeyId,
				Status:        "FAILED",
				StatusMessage: err.Error(),
			})
			continue
		}

		if len(secret.SecretString) > 0 || len(secret.SecretBinary) > 0 && replica.CurrentVersion != "" {
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
						replicaStore.CreateVersionDirect(secret.Name, replicaVersion)
					}
				}
			}
		}

		replica.ReplicationStatus = []secretsmanagerstore.ReplicationStatus{
			{
				Region:           reqCtx.GetRegion(),
				KmsKeyId:         "",
				Status:           "InSync",
				LastAccessedDate: time.Now().UTC(),
			},
		}

		if err := replicaStore.UpdateSecretMetadata(replica); err != nil {
			log.Printf("Failed to set replication status on replica in %s: %v", replicaRegion.Region, err)
		}

		secret.ReplicationStatus = append(secret.ReplicationStatus, secretsmanagerstore.ReplicationStatus{
			Region:           replicaRegion.Region,
			KmsKeyId:         replicaRegion.KmsKeyId,
			Status:           "InSync",
			LastAccessedDate: time.Now().UTC(),
		})
	}

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
// https://docs.aws.amazon.com/secretsmanager/latest/userguide/API_RemoveRegionsFromReplication.html
func (s *SecretsManagerService) RemoveRegionsFromReplication(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if secretId == "" {
		return nil, awserrors.ErrMissingParameter
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
					replicaStore.DeleteSecret(secret.Name)
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
// https://docs.aws.amazon.com/secretsmanager/latest/userguide/API_StopReplicationToReplica.html
func (s *SecretsManagerService) StopReplicationToReplica(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	secretId := request.GetStringParam(req.Parameters, "SecretId")
	if secretId == "" {
		return nil, awserrors.ErrMissingParameter
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

	secret.ReplicationStatus = nil
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
			continue
		}
		regions = append(regions, replicaRegion{Region: region, KmsKeyId: kmsKey})
	}

	return regions, nil
}

func buildReplicationStatusResponse(statuses []secretsmanagerstore.ReplicationStatus) []interface{} {
	result := make([]interface{}, len(statuses))
	for i, rs := range statuses {
		entry := map[string]interface{}{
			"Region":   rs.Region,
			"KmsKeyId": rs.KmsKeyId,
			"Status":   rs.Status,
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

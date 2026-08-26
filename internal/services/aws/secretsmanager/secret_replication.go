package secretsmanager

import (
	"context"

	"vorpalstacks/internal/common/request"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
)

// ReplicateSecretToRegions replicates a secret to one or more regions.
func (s *SecretsManagerService) ReplicateSecretToRegions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	addReplicaRegions, err := parseReplicaRegions(req.Parameters, "AddReplicaRegions")
	if err != nil {
		return nil, err
	}

	result, err := s.replicateSecretToRegionsOpCore(ctx, store, ReplicateSecretToRegionsInput{
		SecretId:                    request.GetStringParam(req.Parameters, "SecretId"),
		AddReplicaRegions:           addReplicaRegions,
		ForceOverwriteReplicaSecret: request.GetBoolParam(req.Parameters, "ForceOverwriteReplicaSecret"),
		Region:                      reqCtx.GetRegion(),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ARN":               result.ARN,
		"Name":              result.Name,
		"ReplicationStatus": buildReplicationStatusResponse(result.ReplicationStatus),
	}, nil
}

// RemoveRegionsFromReplication removes replica regions from a secret.
func (s *SecretsManagerService) RemoveRegionsFromReplication(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.removeRegionsFromReplicationCore(ctx, store, RemoveRegionsFromReplicationInput{
		SecretId:             request.GetStringParam(req.Parameters, "SecretId"),
		RemoveReplicaRegions: request.GetStringList(req.Parameters, "RemoveReplicaRegions"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ARN":               result.ARN,
		"Name":              result.Name,
		"ReplicationStatus": buildReplicationStatusResponse(result.ReplicationStatus),
	}, nil
}

// StopReplicationToReplica stops replication to a replica secret, promoting it to a standalone secret.
func (s *SecretsManagerService) StopReplicationToReplica(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.stopReplicationToReplicaCore(ctx, store, StopReplicationToReplicaInput{
		SecretId: request.GetStringParam(req.Parameters, "SecretId"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ARN":  result.ARN,
		"Name": result.Name,
	}, nil
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

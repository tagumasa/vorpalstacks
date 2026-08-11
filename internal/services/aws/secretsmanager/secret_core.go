package secretsmanager

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/store/aws/common"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
	"vorpalstacks/internal/utils/aws/types"
)

// ---------------------------------------------------------------------------
// Input structs — transport-agnostic DTOs shared by HTTP API and admin handler
// ---------------------------------------------------------------------------

// CreateSecretInput carries all fields needed for CreateSecret.
type CreateSecretInput struct {
	Name               string
	SecretString       string
	SecretBinary       []byte
	Description        string
	KmsKeyId           string
	Type               string
	ClientRequestToken string
	Tags               map[string]string
}

// DeleteSecretInput carries all fields needed for DeleteSecret.
type DeleteSecretInput struct {
	SecretId                   string
	ForceDeleteWithoutRecovery bool
	RecoveryWindowInDays       int
	HasRecoveryWindow          bool
}

// ListSecretsInput carries pagination and filter parameters for ListSecrets.
type ListSecretsInput struct {
	MaxResults             int
	NextToken              string
	Filters                []map[string]interface{}
	IncludePlannedDeletion bool
}

// CreateSecretResult holds the transport-agnostic result of CreateSecret.
type CreateSecretResult struct {
	ARN       string
	Name      string
	VersionID string
}

// DeleteSecretResult holds the transport-agnostic result of DeleteSecret.
type DeleteSecretResult struct {
	ARN          string
	Name         string
	DeletionDate time.Time
}

// ListSecretsResult holds the transport-agnostic result of ListSecrets.
type ListSecretsResult struct {
	Secrets   []*secretsmanagerstore.Secret
	NextToken string
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// createSecretCore is the single entry point for secret creation logic
// shared by the HTTP API and the admin gRPC handler. It performs all field
// validation, constructs the Secret struct, persists it, and returns the
// creation result.
func (s *SecretsManagerService) createSecretCore(ctx context.Context, store secretsmanagerstore.SecretStoreInterface, in CreateSecretInput) (*CreateSecretResult, error) {
	if err := validateSecretName(in.Name); err != nil {
		return nil, err
	}
	if in.SecretString != "" && len(in.SecretBinary) > 0 {
		return nil, awserrors.NewAWSError("InvalidParameterException",
			"You can't specify both SecretString and SecretBinary in the same request.", http.StatusBadRequest)
	}
	if err := validateSecretStringLength(in.SecretString); err != nil {
		return nil, err
	}
	if err := validateSecretBinaryLength(in.SecretBinary); err != nil {
		return nil, err
	}
	if err := validateDescription(in.Description); err != nil {
		return nil, err
	}
	if err := validateClientRequestToken(in.ClientRequestToken); err != nil {
		return nil, err
	}

	// Validate tags when provided as a map (admin handler path).
	if len(in.Tags) > 0 {
		tagList := make([]types.Tag, 0, len(in.Tags))
		for k, v := range in.Tags {
			tagList = append(tagList, types.Tag{Key: k, Value: v})
		}
		if err := validateSecretTags(tagList); err != nil {
			return nil, err
		}
	}

	// ClientRequestToken idempotency: if the token is provided and a
	// version with that ID already exists for a secret with this name,
	// check whether the values match (idempotent success) or differ
	// (error). This mirrors the AWS CreateSecret idempotency rules.
	if in.ClientRequestToken != "" {
		existing, metaErr := store.GetSecretForMetadata(in.Name)
		if metaErr == nil && len(existing.VersionIDs) > 0 {
			for _, vid := range existing.VersionIDs {
				if vid == in.ClientRequestToken {
					existingVer, verErr := store.GetSecretVersion(in.Name, in.ClientRequestToken)
					if verErr != nil {
						break
					}
					if existingVer.SecretString == in.SecretString &&
						bytesEqual(existingVer.SecretBinary, in.SecretBinary) {
						return &CreateSecretResult{
							ARN:       existing.ARN,
							Name:      existing.Name,
							VersionID: in.ClientRequestToken,
						}, nil
					}
					return nil, awserrors.NewAWSError("InvalidRequestException",
						fmt.Sprintf("You can't modify an existing secret version. The ClientRequestToken %s is already associated with a different version value.", in.ClientRequestToken),
						http.StatusBadRequest)
				}
			}
		}
	}

	secret := secretsmanagerstore.NewSecret(in.Name)
	secret.SecretString = in.SecretString
	secret.SecretBinary = in.SecretBinary
	secret.Description = in.Description
	secret.KmsKeyId = in.KmsKeyId
	secret.Type = in.Type
	secret.InitialVersionId = in.ClientRequestToken
	if len(in.Tags) > 0 {
		secret.Tags = in.Tags
	}

	created, err := store.CreateSecret(secret)
	if err != nil {
		return nil, mapStoreError(err)
	}

	return &CreateSecretResult{
		ARN:       created.ARN,
		Name:      created.Name,
		VersionID: created.CurrentVersion,
	}, nil
}

// deleteSecretCore is the single entry point for secret deletion logic
// shared by the HTTP API and the admin gRPC handler.
func (s *SecretsManagerService) deleteSecretCore(ctx context.Context, store secretsmanagerstore.SecretStoreInterface, in DeleteSecretInput) (*DeleteSecretResult, error) {
	if in.SecretId == "" {
		return nil, awserrors.ErrMissingParameter
	}
	if err := validateRecoveryWindow(in.RecoveryWindowInDays, in.HasRecoveryWindow, in.ForceDeleteWithoutRecovery); err != nil {
		return nil, err
	}

	// Resolve the secret by name or ARN, matching the HTTP handler's
	// resolveSecretForMetadata behaviour.
	var secret *secretsmanagerstore.Secret
	if strings.HasPrefix(in.SecretId, arnPrefix) {
		name, err := store.LookupNameByARN(in.SecretId)
		if err != nil {
			return nil, mapStoreError(err)
		}
		secret, err = store.GetSecretForMetadata(name)
		if err != nil {
			return nil, mapStoreError(err)
		}
	} else {
		var err error
		secret, err = store.GetSecretForMetadata(in.SecretId)
		if err != nil {
			return nil, mapStoreError(err)
		}
	}

	// You can't delete a primary secret that is replicated to other Regions.
	if len(secret.ReplicationStatus) > 0 {
		return nil, awserrors.NewAWSError("InvalidRequestException",
			"You can't delete a primary secret that is replicated to other Regions. Remove the replicas first.", http.StatusBadRequest)
	}

	var deletionDate time.Time
	if in.ForceDeleteWithoutRecovery {
		if err := store.DeleteSecret(secret.Name); err != nil {
			return nil, mapStoreError(err)
		}
		deletionDate = time.Now().UTC()
	} else {
		rw := in.RecoveryWindowInDays
		if rw == 0 {
			rw = 30
		}
		deletionDate = time.Now().UTC().AddDate(0, 0, rw)
		if err := store.ScheduleDeletion(secret.Name, deletionDate); err != nil {
			return nil, mapStoreError(err)
		}
	}

	return &DeleteSecretResult{
		ARN:          secret.ARN,
		Name:         secret.Name,
		DeletionDate: deletionDate,
	}, nil
}

// listSecretsCore is the single entry point for listing secrets,
// shared by the HTTP API and the admin gRPC handler.
func (s *SecretsManagerService) listSecretsCore(ctx context.Context, store secretsmanagerstore.SecretStoreInterface, in ListSecretsInput) (*ListSecretsResult, error) {
	if err := validateMaxFilters(len(in.Filters)); err != nil {
		return nil, err
	}

	secretFilter := buildSecretFilter(in.IncludePlannedDeletion, in.Filters)

	opts := common.ListOptions{
		Marker: in.NextToken,
	}
	if in.MaxResults < 0 {
		// Sentinel: caller wants all items for client-side sorting.
		opts.MaxItems = 0
	} else {
		maxResults := in.MaxResults
		if maxResults <= 0 {
			maxResults = 100
		}
		opts.MaxItems = maxResults
	}

	result, err := store.ListSecrets(opts, secretFilter)
	if err != nil {
		return nil, mapStoreError(err)
	}

	return &ListSecretsResult{
		Secrets:   result.Items,
		NextToken: result.NextMarker,
	}, nil
}

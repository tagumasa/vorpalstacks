package secretsmanager

import (
	"context"

	"vorpalstacks/internal/common/invokers"
)

// SecretsManagerBusInvoker returns an invokers.SecretsManagerInvoker backed
// by this service's Core functions, for cross-service consumers that own
// service secrets (e.g. EventBridge connection credentials).
func (s *SecretsManagerService) SecretsManagerBusInvoker() invokers.SecretsManagerInvoker {
	return &secretsManagerBusAdapter{s}
}

// secretsManagerBusAdapter adapts SecretsManagerService to the
// invokers.SecretsManagerInvoker contract.
type secretsManagerBusAdapter struct {
	s *SecretsManagerService
}

// CreateServiceSecret creates a secret holding secretString and returns its
// ARN. The description records the owning service so operators can tell
// service-owned secrets from user-created ones.
func (a *secretsManagerBusAdapter) CreateServiceSecret(ctx context.Context, region, name, secretString, description string) (string, error) {
	store, err := a.s.GetStoreForRegion(region)
	if err != nil {
		return "", err
	}
	created, err := a.s.createSecretCore(ctx, store, CreateSecretInput{
		Name:         name,
		SecretString: secretString,
		Description:  description,
		Region:       region,
	})
	if err != nil {
		return "", err
	}
	return created.ARN, nil
}

// GetServiceSecretString returns the current secret string of the secret
// addressed by ARN or name.
func (a *secretsManagerBusAdapter) GetServiceSecretString(ctx context.Context, region, secretId string) (string, error) {
	store, err := a.s.GetStoreForRegion(region)
	if err != nil {
		return "", err
	}
	result, err := a.s.getSecretValueCore(ctx, store, GetSecretValueInput{SecretId: secretId})
	if err != nil {
		return "", err
	}
	return result.Version.SecretString, nil
}

// UpdateServiceSecretString replaces the secret string of the secret
// addressed by ARN or name.
func (a *secretsManagerBusAdapter) UpdateServiceSecretString(ctx context.Context, region, secretId, secretString string) error {
	store, err := a.s.GetStoreForRegion(region)
	if err != nil {
		return err
	}
	_, err = a.s.updateSecretCore(ctx, store, UpdateSecretInput{
		SecretId:     secretId,
		SecretString: secretString,
	})
	return err
}

// DeleteServiceSecret immediately deletes the secret addressed by ARN or
// name. Service-owned secrets follow their owning resource's lifecycle, so
// no recovery window is scheduled. A missing secret is already in the
// caller's desired state and is not an error.
func (a *secretsManagerBusAdapter) DeleteServiceSecret(ctx context.Context, region, secretId string) error {
	store, err := a.s.GetStoreForRegion(region)
	if err != nil {
		return err
	}
	_, err = a.s.deleteSecretCore(ctx, store, DeleteSecretInput{
		SecretId:                   secretId,
		ForceDeleteWithoutRecovery: true,
	})
	if err == ErrSecretNotFound {
		return nil
	}
	return err
}

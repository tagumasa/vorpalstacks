// Package secretsmanager provides AWS Secrets Manager operations for vorpalstacks.
package secretsmanager

import (
	"vorpalstacks/internal/server/dispatcher"
	"vorpalstacks/internal/services/aws/common/request"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
)

// SecretsManagerService provides AWS Secrets Manager operations.
type SecretsManagerService struct {
	accountID string
}

// NewSecretsManagerService creates a new Secrets Manager service instance.
func NewSecretsManagerService(accountID string) *SecretsManagerService {
	return &SecretsManagerService{
		accountID: accountID,
	}
}

// store returns the Secrets Manager store for the given request context.
func (s *SecretsManagerService) store(reqCtx *request.RequestContext) (secretsmanagerstore.SecretStoreInterface, error) {
	store := reqCtx.GetSecretsManagerStore()
	if store != nil {
		return store, nil
	}
	storage, err := reqCtx.GetStorage()
	if err != nil {
		return nil, err
	}
	return secretsmanagerstore.NewSecretStore(storage, s.accountID, reqCtx.GetRegion()), nil
}

// RegisterHandlers registers the Secrets Manager handlers with the dispatcher.
func (s *SecretsManagerService) RegisterHandlers(d *dispatcher.Dispatcher) {
	d.RegisterHandlerForService("secretsmanager", "CreateSecret", s.CreateSecret)
	d.RegisterHandlerForService("secretsmanager", "GetSecretValue", s.GetSecretValue)
	d.RegisterHandlerForService("secretsmanager", "UpdateSecret", s.UpdateSecret)
	d.RegisterHandlerForService("secretsmanager", "DeleteSecret", s.DeleteSecret)
	d.RegisterHandlerForService("secretsmanager", "ListSecrets", s.ListSecrets)
	d.RegisterHandlerForService("secretsmanager", "DescribeSecret", s.DescribeSecret)
	d.RegisterHandlerForService("secretsmanager", "PutSecretValue", s.PutSecretValue)
	d.RegisterHandlerForService("secretsmanager", "TagResource", s.TagResource)
	d.RegisterHandlerForService("secretsmanager", "UntagResource", s.UntagResource)
	d.RegisterHandlerForService("secretsmanager", "ListTagsForResource", s.ListTagsForResource)
	d.RegisterHandlerForService("secretsmanager", "RestoreSecret", s.RestoreSecret)
	d.RegisterHandlerForService("secretsmanager", "RotateSecret", s.RotateSecret)
	d.RegisterHandlerForService("secretsmanager", "CancelRotateSecret", s.CancelRotateSecret)
	d.RegisterHandlerForService("secretsmanager", "GetResourcePolicy", s.GetResourcePolicy)
	d.RegisterHandlerForService("secretsmanager", "PutResourcePolicy", s.PutResourcePolicy)
	d.RegisterHandlerForService("secretsmanager", "DeleteResourcePolicy", s.DeleteResourcePolicy)
	d.RegisterHandlerForService("secretsmanager", "ValidateResourcePolicy", s.ValidateResourcePolicy)
	d.RegisterHandlerForService("secretsmanager", "GetRandomPassword", s.GetRandomPassword)
	d.RegisterHandlerForService("secretsmanager", "ListSecretVersionIds", s.ListSecretVersionIds)
	d.RegisterHandlerForService("secretsmanager", "UpdateSecretVersionStage", s.UpdateSecretVersionStage)
	d.RegisterHandlerForService("secretsmanager", "BatchGetSecretValue", s.BatchGetSecretValue)
}

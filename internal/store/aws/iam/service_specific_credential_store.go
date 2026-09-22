package iam

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const serviceSpecificCredentialBucketName = "iam_service_credentials"

// ServiceSpecificCredentialStore provides storage operations for IAM service-specific credentials.
type ServiceSpecificCredentialStore struct {
	uk         userKeyed[ServiceSpecificCredential]
	arnBuilder *ARNBuilder
}

// NewServiceSpecificCredentialStore creates a new ServiceSpecificCredentialStore instance.
func NewServiceSpecificCredentialStore(store storage.BasicStorage, accountID string) *ServiceSpecificCredentialStore {
	return &ServiceSpecificCredentialStore{
		uk: newUserKeyed[ServiceSpecificCredential](
			common.NewBaseStore(store.Bucket(serviceSpecificCredentialBucketName), "iam"),
			func(c *ServiceSpecificCredential) string { return c.ServiceSpecificCredentialId },
			func(c *ServiceSpecificCredential) string { return c.UserName },
		),
		arnBuilder: NewARNBuilder(accountID),
	}
}

// Get retrieves a service-specific credential by its ID.
func (s *ServiceSpecificCredentialStore) Get(credentialId string) (*ServiceSpecificCredential, error) {
	return getByKey[ServiceSpecificCredential](s.uk.BaseStore, credentialId, "get_service_specific_credential", ErrServiceSpecificCredentialNotFound)
}

// Put stores a service-specific credential, keyed by its ID.
func (s *ServiceSpecificCredentialStore) Put(cred *ServiceSpecificCredential) error {
	return s.uk.BaseStore.Put(cred.ServiceSpecificCredentialId, cred)
}

// Delete removes a service-specific credential by its ID.
func (s *ServiceSpecificCredentialStore) Delete(credentialId string) error {
	return s.uk.BaseStore.Delete(credentialId)
}

// Exists reports whether a service-specific credential exists with the given ID.
func (s *ServiceSpecificCredentialStore) Exists(credentialId string) bool {
	return s.uk.BaseStore.Exists(credentialId)
}

// Create generates a new service-specific credential for the given user and service.
func (s *ServiceSpecificCredentialStore) Create(userName, serviceName string, credentialAgeDays int) (*ServiceSpecificCredential, error) {
	id, err := GenerateServiceSpecificCredentialID()
	if err != nil {
		return nil, err
	}
	password, err := generateServicePassword()
	if err != nil {
		return nil, err
	}
	credName := fmt.Sprintf("%s-at-%d", userName, time.Now().Unix())
	now := time.Now().UTC()
	cred := &ServiceSpecificCredential{
		ServiceSpecificCredentialId:   id,
		ServiceSpecificCredentialName: credName,
		ServiceName:                   serviceName,
		UserName:                      userName,
		ServicePassword:               password,
		ServiceSpecificCredentialArn:  s.arnBuilder.UserARN("", userName),
		CreateDate:                    now,
		Status:                        "Active",
	}
	if credentialAgeDays > 0 {
		exp := now.AddDate(0, 0, credentialAgeDays)
		cred.ExpirationDate = &exp
	}
	if err := s.Put(cred); err != nil {
		return nil, err
	}
	return cred, nil
}

// ResetPassword generates a new password for the specified service-specific
// credential, keeping the same credential ID and all other metadata intact.
func (s *ServiceSpecificCredentialStore) ResetPassword(credentialId string) (*ServiceSpecificCredential, error) {
	var result *ServiceSpecificCredential
	err := s.uk.kl.WithLock(credentialId, func() error {
		cred, err := s.Get(credentialId)
		if err != nil {
			return err
		}
		password, err := generateServicePassword()
		if err != nil {
			return err
		}
		cred.ServicePassword = password
		if err := s.Put(cred); err != nil {
			return err
		}
		result = cred
		return nil
	})
	return result, err
}

// UpdateStatus changes the status of a service-specific credential (e.g. Active/Inactive).
func (s *ServiceSpecificCredentialStore) UpdateStatus(credentialId, status string) error {
	return s.uk.updateStatus(credentialId, status, s.Get, func(c *ServiceSpecificCredential, status string) { c.Status = status })
}

// ListByUserName returns all service-specific credentials for the given user.
func (s *ServiceSpecificCredentialStore) ListByUserName(userName string) ([]*ServiceSpecificCredential, error) {
	return s.uk.listByUserName(userName, "list_service_specific_credentials")
}

// FindByServiceAndSecret resolves a credential by its service name and
// secret — the lookup the CloudWatch Logs HTTP ingestion endpoints'
// bearer-token authentication needs (the token IS the credential's
// ServiceCredentialSecret). The comparison is constant-time.
func (s *ServiceSpecificCredentialStore) FindByServiceAndSecret(serviceName, secret string) (*ServiceSpecificCredential, error) {
	return common.FindFirst[ServiceSpecificCredential](s.uk.BaseStore, func(c *ServiceSpecificCredential) bool {
		return c.ServiceName == serviceName &&
			subtle.ConstantTimeCompare([]byte(c.ServicePassword), []byte(secret)) == 1
	})
}

// DeleteAllForUser removes all service-specific credentials belonging to the given user.
func (s *ServiceSpecificCredentialStore) DeleteAllForUser(userName string) error {
	return s.uk.deleteAllForUser(userName, "delete_user_service_credentials")
}

// MigrateUser updates the UserName field on all service-specific credentials
// from oldUserName to newUserName. Called during IAM user rename operations.
func (s *ServiceSpecificCredentialStore) MigrateUser(oldUserName, newUserName string) error {
	return s.uk.migrateUser(oldUserName, newUserName, "migrate_service_credentials", func(cred *ServiceSpecificCredential, newName string) {
		cred.UserName = newName
		cred.ServiceSpecificCredentialArn = s.arnBuilder.UserARN("", newName)
	})
}

// Count returns the total number of service-specific credentials.
func (s *ServiceSpecificCredentialStore) Count() int {
	return s.uk.BaseStore.Count()
}

func generateServicePassword() (string, error) {
	bytes := make([]byte, 30)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(bytes)
	trimmed := strings.ReplaceAll(encoded, "=", "")
	trimmed = strings.ReplaceAll(trimmed, "+", "")
	trimmed = strings.ReplaceAll(trimmed, "/", "")
	if len(trimmed) < 28 {
		return "", fmt.Errorf("generated password too short: %d < 28", len(trimmed))
	}
	return trimmed[:28], nil
}

package iam

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
)

const serviceSpecificCredentialBucketName = "iam_service_credentials"

// ServiceSpecificCredentialStore provides storage operations for IAM service-specific credentials.
type ServiceSpecificCredentialStore struct {
	bucket    storage.Bucket
	accountID string
}

// NewServiceSpecificCredentialStore creates a new ServiceSpecificCredentialStore instance.
func NewServiceSpecificCredentialStore(store storage.BasicStorage, accountID string) *ServiceSpecificCredentialStore {
	return &ServiceSpecificCredentialStore{
		bucket:    store.Bucket(serviceSpecificCredentialBucketName),
		accountID: accountID,
	}
}

// Get retrieves a service-specific credential by its ID.
func (s *ServiceSpecificCredentialStore) Get(credentialId string) (*ServiceSpecificCredential, error) {
	data, err := s.bucket.Get([]byte(credentialId))
	if err != nil {
		return nil, NewStoreError("get_service_specific_credential", err)
	}
	if data == nil {
		return nil, NewStoreError("get_service_specific_credential", ErrServiceSpecificCredentialNotFound)
	}
	var cred ServiceSpecificCredential
	if err := json.Unmarshal(data, &cred); err != nil {
		return nil, NewStoreError("get_service_specific_credential", err)
	}
	return &cred, nil
}

// Put stores a service-specific credential, keyed by its ID.
func (s *ServiceSpecificCredentialStore) Put(cred *ServiceSpecificCredential) error {
	data, err := json.Marshal(cred)
	if err != nil {
		return NewStoreError("put_service_specific_credential", err)
	}
	if err := s.bucket.Put([]byte(cred.ServiceSpecificCredentialId), data); err != nil {
		return NewStoreError("put_service_specific_credential", err)
	}
	return nil
}

// Delete removes a service-specific credential by its ID.
func (s *ServiceSpecificCredentialStore) Delete(credentialId string) error {
	if err := s.bucket.Delete([]byte(credentialId)); err != nil {
		return NewStoreError("delete_service_specific_credential", err)
	}
	return nil
}

// Exists reports whether a service-specific credential exists with the given ID.
func (s *ServiceSpecificCredentialStore) Exists(credentialId string) bool {
	return s.bucket.Has([]byte(credentialId))
}

// Create generates a new service-specific credential for the given user and service.
func (s *ServiceSpecificCredentialStore) Create(userName, serviceName string) (*ServiceSpecificCredential, error) {
	id, err := generateServiceCredentialID()
	if err != nil {
		return nil, err
	}
	password, err := generateServicePassword()
	if err != nil {
		return nil, err
	}
	credName := fmt.Sprintf("%s-at-%d", userName, time.Now().Unix())
	cred := &ServiceSpecificCredential{
		ServiceSpecificCredentialId:   id,
		ServiceSpecificCredentialName: credName,
		ServiceName:                   serviceName,
		UserName:                      userName,
		ServicePassword:               password,
		ServiceSpecificCredentialArn:  fmt.Sprintf("arn:aws:iam:%s:user/%s", s.accountID, userName),
		CreateDate:                    time.Now().UTC(),
		Status:                        "Active",
	}
	if err := s.Put(cred); err != nil {
		return nil, err
	}
	return cred, nil
}

// UpdateStatus changes the status of a service-specific credential (e.g. Active/Inactive).
func (s *ServiceSpecificCredentialStore) UpdateStatus(credentialId, status string) error {
	cred, err := s.Get(credentialId)
	if err != nil {
		return err
	}
	cred.Status = status
	return s.Put(cred)
}

// ListByUserName returns all service-specific credentials for the given user.
func (s *ServiceSpecificCredentialStore) ListByUserName(userName string) ([]*ServiceSpecificCredential, error) {
	var creds []*ServiceSpecificCredential
	err := s.bucket.ForEach(func(k, v []byte) error {
		var cred ServiceSpecificCredential
		if err := json.Unmarshal(v, &cred); err != nil {
			return err
		}
		if cred.UserName == userName {
			creds = append(creds, &cred)
		}
		return nil
	})
	if err != nil {
		return nil, NewStoreError("list_service_specific_credentials", err)
	}
	return creds, nil
}

// DeleteAllForUser removes all service-specific credentials belonging to the given user.
func (s *ServiceSpecificCredentialStore) DeleteAllForUser(userName string) error {
	var toDelete []string
	err := s.bucket.ForEach(func(k, v []byte) error {
		var cred ServiceSpecificCredential
		if err := json.Unmarshal(v, &cred); err != nil {
			return err
		}
		if cred.UserName == userName {
			toDelete = append(toDelete, cred.ServiceSpecificCredentialId)
		}
		return nil
	})
	if err != nil {
		return NewStoreError("delete_user_service_credentials", err)
	}
	for _, id := range toDelete {
		if err := s.Delete(id); err != nil {
			return err
		}
	}
	return nil
}

// Count returns the total number of service-specific credentials.
func (s *ServiceSpecificCredentialStore) Count() int {
	return s.bucket.Count()
}

func generateServiceCredentialID() (string, error) {
	return generateID("AGPA")
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

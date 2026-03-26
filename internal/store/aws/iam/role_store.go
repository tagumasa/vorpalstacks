package iam

// Package iam provides IAM (Identity and Access Management) data store implementations
// for vorpalstacks.

import (
	"encoding/json"
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
)

const roleBucketName = "iam_roles"

// RoleStore manages IAM roles.
type RoleStore struct {
	*common.BaseStore
	arnBuilder *ARNBuilder
}

// NewRoleStore creates a new RoleStore instance with the specified storage and account ID.
func NewRoleStore(store storage.BasicStorage, accountId string) *RoleStore {
	return &RoleStore{
		BaseStore:  common.NewBaseStore(store.Bucket(roleBucketName), "iam"),
		arnBuilder: NewARNBuilder(accountId),
	}
}

// Get retrieves an IAM role by its name from the store.
// Returns the role or an error if not found.
func (s *RoleStore) Get(roleName string) (*Role, error) {
	var role Role
	if err := s.BaseStore.Get(roleName, &role); err != nil {
		return nil, err
	}
	return &role, nil
}

// GetByID retrieves an IAM role by its ID from the store.
// Returns the role or an error if not found.
func (s *RoleStore) GetByID(roleID string) (*Role, error) {
	var found *Role
	err := s.ForEach(func(key string, value []byte) error {
		var role Role
		if err := json.Unmarshal(value, &role); err != nil {
			return err
		}
		if role.ID == roleID && found == nil {
			found = &role
		}
		return nil
	})
	if err != nil {
		return nil, NewStoreError("get_role_by_id", err)
	}
	if found == nil {
		return nil, NewStoreError("get_role_by_id", ErrRoleNotFound)
	}
	return found, nil
}

// List returns a list of IAM roles from the store with pagination support.
func (s *RoleStore) List(pathPrefix string, marker string, maxItems int) (*RoleListResult, error) {
	if maxItems <= 0 {
		maxItems = 100
	}

	var roles []*Role
	count := 0
	started := marker == ""
	hasMore := false

	err := s.ForEach(func(key string, value []byte) error {
		var role Role
		if err := json.Unmarshal(value, &role); err != nil {
			return err
		}

		if !started {
			if role.RoleName == marker {
				started = true
			}
			return nil
		}

		if pathPrefix != "" && !strings.HasPrefix(role.Path, pathPrefix) {
			return nil
		}

		if count < maxItems {
			roles = append(roles, &role)
			count++
		} else {
			hasMore = true
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_roles", err)
	}

	result := &RoleListResult{
		Roles:       roles,
		IsTruncated: hasMore,
	}
	if len(roles) > 0 {
		result.Marker = roles[len(roles)-1].RoleName
	}

	return result, nil
}

// Put creates or updates an IAM role in the store.
func (s *RoleStore) Put(role *Role) error {
	if role.CreateDate.IsZero() {
		role.CreateDate = time.Now().UTC()
	}
	if role.Path == "" {
		role.Path = "/"
	}
	if role.MaxSessionDuration == 0 {
		role.MaxSessionDuration = 3600
	}
	return s.BaseStore.Put(role.RoleName, role)
}

// Delete deletes an IAM role by its name from the store.
func (s *RoleStore) Delete(roleName string) error {
	return s.BaseStore.Delete(roleName)
}

// Exists checks whether an IAM role exists in the store by its name.
func (s *RoleStore) Exists(roleName string) bool {
	return s.BaseStore.Exists(roleName)
}

// Create creates a new IAM role in the store.
// Returns the created role or an error if the role already exists.
func (s *RoleStore) Create(roleName, path, accountId, assumeRolePolicyDocument, description string, maxSessionDuration int, tags []Tag) (*Role, error) {
	if s.Exists(roleName) {
		return nil, NewStoreError("create_role", ErrRoleAlreadyExists)
	}

	roleID, err := GenerateRoleID()
	if err != nil {
		return nil, NewStoreError("generate_role_id", err)
	}

	role := &Role{
		ID:                       roleID,
		Path:                     path,
		RoleName:                 roleName,
		AccountId:                accountId,
		CreateDate:               time.Now().UTC(),
		AssumeRolePolicyDocument: assumeRolePolicyDocument,
		Description:              description,
		MaxSessionDuration:       maxSessionDuration,
		Tags:                     tags,
	}
	role.Arn = s.arnBuilder.RoleARN(path, roleName)

	if err := s.Put(role); err != nil {
		return nil, err
	}
	return role, nil
}

// Update updates an existing IAM role in the store.
func (s *RoleStore) Update(role *Role) error {
	return s.Put(role)
}

// Count returns the number of IAM roles in the store.
func (s *RoleStore) Count() int {
	return s.BaseStore.Count()
}

// GetAssumeRolePolicyDocument returns the assume role policy document for a role.
// This implements the iam.RolePolicyProvider interface.
func (s *RoleStore) GetAssumeRolePolicyDocument(roleName string) (string, error) {
	role, err := s.Get(roleName)
	if err != nil {
		return "", err
	}
	return role.AssumeRolePolicyDocument, nil
}

// RoleExists checks if a role exists.
// This implements the iam.RolePolicyProvider interface.
func (s *RoleStore) RoleExists(roleName string) bool {
	return s.Exists(roleName)
}

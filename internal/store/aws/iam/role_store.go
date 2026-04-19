package iam

import (
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/utils/aws/types"
)

const roleBucketName = "iam_roles"

// RoleStore manages IAM roles.
type RoleStore struct {
	entityStore[Role]
	arnBuilder *ARNBuilder
}

// NewRoleStore creates a new RoleStore instance with the specified storage and account ID.
func NewRoleStore(store storage.BasicStorage, accountId string) *RoleStore {
	return &RoleStore{
		entityStore: newEntityStore[Role](store, roleBucketName),
		arnBuilder:  NewARNBuilder(accountId),
	}
}

// GetByID retrieves an IAM role by its ID from the store.
func (s *RoleStore) GetByID(roleID string) (*Role, error) {
	return getEntityByID(s.BaseStore, roleID, func(r *Role) string { return r.ID }, "get_role_by_id", ErrRoleNotFound)
}

// List returns a list of IAM roles from the store with pagination support.
func (s *RoleStore) List(pathPrefix string, marker string, maxItems int) (*RoleListResult, error) {
	items, truncated, nextMarker, err := listEntitiesWithPathPrefix(s.BaseStore, pathPrefix, marker, maxItems, func(r *Role) string { return r.Path })
	if err != nil {
		return nil, err
	}
	return &RoleListResult{
		Roles:       items,
		IsTruncated: truncated,
		Marker:      nextMarker,
	}, nil
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

// Create creates a new IAM role in the store.
func (s *RoleStore) Create(roleName, path, accountId, assumeRolePolicyDocument, description string, maxSessionDuration int, tags []types.Tag) (*Role, error) {
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

// GetAssumeRolePolicyDocument returns the assume role policy document for a role.
func (s *RoleStore) GetAssumeRolePolicyDocument(roleName string) (string, error) {
	role, err := s.Get(roleName)
	if err != nil {
		return "", err
	}
	return role.AssumeRolePolicyDocument, nil
}

// RoleExists checks if a role exists.
func (s *RoleStore) RoleExists(roleName string) bool {
	return s.Exists(roleName)
}

package iam

import (
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	"vorpalstacks/internal/utils/aws/types"
)

const instanceProfileBucketName = "iam_instance_profiles"

// InstanceProfileStore manages IAM instance profile data in persistent storage.
type InstanceProfileStore struct {
	entityStore[InstanceProfile]
	arnBuilder *ARNBuilder
}

// NewInstanceProfileStore creates a new store for IAM instance profiles.
func NewInstanceProfileStore(store storage.BasicStorage, accountId string) *InstanceProfileStore {
	return &InstanceProfileStore{
		entityStore: newEntityStore[InstanceProfile](store, instanceProfileBucketName),
		arnBuilder:  NewARNBuilder(accountId),
	}
}

// GetByID retrieves an instance profile by its ID.
func (s *InstanceProfileStore) GetByID(instanceProfileID string) (*InstanceProfile, error) {
	return getEntityByID(s.BaseStore, instanceProfileID, func(p *InstanceProfile) string { return p.ID }, "get_instance_profile_by_id", ErrInstanceProfileNotFound)
}

// List returns a paginated list of instance profiles.
func (s *InstanceProfileStore) List(pathPrefix string, marker string, maxItems int) (*InstanceProfileListResult, error) {
	items, truncated, nextMarker, err := listEntitiesWithPathPrefix(s.BaseStore, pathPrefix, marker, maxItems, func(p *InstanceProfile) string { return p.Path })
	if err != nil {
		return nil, err
	}
	return &InstanceProfileListResult{
		InstanceProfiles: items,
		IsTruncated:      truncated,
		Marker:           nextMarker,
	}, nil
}

// Put stores an instance profile.
func (s *InstanceProfileStore) Put(profile *InstanceProfile) error {
	if profile.CreateDate.IsZero() {
		profile.CreateDate = time.Now().UTC()
	}
	if profile.Path == "" {
		profile.Path = "/"
	}
	return s.BaseStore.Put(profile.InstanceProfileName, profile)
}

// Create creates a new IAM instance profile.
func (s *InstanceProfileStore) Create(instanceProfileName, path, accountId string, tags []types.Tag) (*InstanceProfile, error) {
	if s.Exists(instanceProfileName) {
		return nil, NewStoreError("create_instance_profile", ErrInstanceProfileAlreadyExists)
	}

	instanceProfileID, err := GenerateInstanceProfileID()
	if err != nil {
		return nil, NewStoreError("generate_instance_profile_id", err)
	}

	profile := &InstanceProfile{
		ID:                  instanceProfileID,
		Path:                path,
		InstanceProfileName: instanceProfileName,
		AccountId:           accountId,
		CreateDate:          time.Now().UTC(),
		Roles:               []string{},
		Tags:                tags,
	}
	profile.Arn = s.arnBuilder.InstanceProfileARN(path, instanceProfileName)

	if err := s.Put(profile); err != nil {
		return nil, err
	}
	return profile, nil
}

// AddRole adds a role to an instance profile.
func (s *InstanceProfileStore) AddRole(instanceProfileName, roleName string) error {
	return s.kl.WithLock(instanceProfileName, func() error {
		profile, err := s.Get(instanceProfileName)
		if err != nil {
			return err
		}

		for _, r := range profile.Roles {
			if r == roleName {
				return NewStoreError("add_role_to_instance_profile", ErrRoleAlreadyInInstanceProfile)
			}
		}

		profile.Roles = append(profile.Roles, roleName)
		return s.Put(profile)
	})
}

// RemoveRole removes a role from an instance profile.
func (s *InstanceProfileStore) RemoveRole(instanceProfileName, roleName string) error {
	return s.kl.WithLock(instanceProfileName, func() error {
		profile, err := s.Get(instanceProfileName)
		if err != nil {
			return err
		}

		found := false
		newRoles := []string{}
		for _, r := range profile.Roles {
			if r == roleName {
				found = true
			} else {
				newRoles = append(newRoles, r)
			}
		}

		if !found {
			return NewStoreError("remove_role_from_instance_profile", ErrRoleNotInInstanceProfile)
		}

		profile.Roles = newRoles
		return s.Put(profile)
	})
}

// HasRole checks whether a role is attached to an instance profile.
func (s *InstanceProfileStore) HasRole(instanceProfileName, roleName string) (bool, error) {
	profile, err := s.Get(instanceProfileName)
	if err != nil {
		return false, err
	}

	for _, r := range profile.Roles {
		if r == roleName {
			return true, nil
		}
	}
	return false, nil
}

// ListForRole returns all instance profiles that contain a specific role.
func (s *InstanceProfileStore) ListForRole(roleName string, marker string, maxItems int) (*InstanceProfileListResult, error) {
	filter := func(p *InstanceProfile) bool {
		for _, r := range p.Roles {
			if r == roleName {
				return true
			}
		}
		return false
	}
	result, err := common.List[InstanceProfile](s.BaseStore, common.ListOptions{Marker: marker, MaxItems: maxItems}, filter)
	if err != nil {
		return nil, err
	}
	return &InstanceProfileListResult{
		InstanceProfiles: result.Items,
		IsTruncated:      result.IsTruncated,
		Marker:           result.NextMarker,
	}, nil
}

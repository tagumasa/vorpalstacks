package iam

// Package iam provides IAM (Identity and Access Management) data store implementations
// for vorpalstacks.

import (
	"encoding/json"
	"strings"
	"time"

	"vorpalstacks/internal/core/storage"
)

const instanceProfileBucketName = "iam_instance_profiles"

// InstanceProfileStore manages IAM instance profile data in persistent storage.
type InstanceProfileStore struct {
	bucket     storage.Bucket
	arnBuilder *ARNBuilder
}

// NewInstanceProfileStore creates a new store for IAM instance profiles.
func NewInstanceProfileStore(store storage.BasicStorage, accountId string) *InstanceProfileStore {
	return &InstanceProfileStore{
		bucket:     store.Bucket(instanceProfileBucketName),
		arnBuilder: NewARNBuilder(accountId),
	}
}

// Get retrieves an instance profile by its name.
func (s *InstanceProfileStore) Get(instanceProfileName string) (*InstanceProfile, error) {
	data, err := s.bucket.Get([]byte(instanceProfileName))
	if err != nil {
		return nil, NewStoreError("get_instance_profile", err)
	}
	if data == nil {
		return nil, NewStoreError("get_instance_profile", ErrInstanceProfileNotFound)
	}
	var profile InstanceProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, NewStoreError("get_instance_profile", err)
	}
	return &profile, nil
}

// GetByID retrieves an instance profile by its ID.
func (s *InstanceProfileStore) GetByID(instanceProfileID string) (*InstanceProfile, error) {
	var found *InstanceProfile
	err := s.bucket.ForEach(func(k, v []byte) error {
		var profile InstanceProfile
		if err := json.Unmarshal(v, &profile); err != nil {
			return err
		}
		if profile.ID == instanceProfileID {
			found = &profile
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, NewStoreError("get_instance_profile_by_id", err)
	}
	if found == nil {
		return nil, NewStoreError("get_instance_profile_by_id", ErrInstanceProfileNotFound)
	}
	return found, nil
}

// List returns a paginated list of instance profiles.
func (s *InstanceProfileStore) List(pathPrefix string, marker string, maxItems int) (*InstanceProfileListResult, error) {
	if maxItems <= 0 {
		maxItems = 100
	}

	var profiles []*InstanceProfile
	count := 0
	started := marker == ""
	var hasMore bool

	err := s.bucket.ForEach(func(k, v []byte) error {
		var profile InstanceProfile
		if err := json.Unmarshal(v, &profile); err != nil {
			return err
		}

		if !started {
			if profile.InstanceProfileName == marker {
				started = true
			}
			return nil
		}

		if pathPrefix != "" && !strings.HasPrefix(profile.Path, pathPrefix) {
			return nil
		}

		if count < maxItems {
			profiles = append(profiles, &profile)
			count++
		} else {
			hasMore = true
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_instance_profiles", err)
	}

	result := &InstanceProfileListResult{
		InstanceProfiles: profiles,
		IsTruncated:      hasMore,
	}
	if len(profiles) > 0 {
		result.Marker = profiles[len(profiles)-1].InstanceProfileName
	}

	return result, nil
}

// Put stores an instance profile.
func (s *InstanceProfileStore) Put(profile *InstanceProfile) error {
	if profile.CreateDate.IsZero() {
		profile.CreateDate = time.Now().UTC()
	}
	if profile.Path == "" {
		profile.Path = "/"
	}

	data, err := json.Marshal(profile)
	if err != nil {
		return NewStoreError("put_instance_profile", err)
	}

	if err := s.bucket.Put([]byte(profile.InstanceProfileName), data); err != nil {
		return NewStoreError("put_instance_profile", err)
	}
	return nil
}

// Delete removes an instance profile by its name.
func (s *InstanceProfileStore) Delete(instanceProfileName string) error {
	if err := s.bucket.Delete([]byte(instanceProfileName)); err != nil {
		return NewStoreError("delete_instance_profile", err)
	}
	return nil
}

// Exists checks whether an instance profile exists.
func (s *InstanceProfileStore) Exists(instanceProfileName string) bool {
	return s.bucket.Has([]byte(instanceProfileName))
}

// Create creates a new IAM instance profile.
func (s *InstanceProfileStore) Create(instanceProfileName, path, accountId string, tags []Tag) (*InstanceProfile, error) {
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
}

// RemoveRole removes a role from an instance profile.
func (s *InstanceProfileStore) RemoveRole(instanceProfileName, roleName string) error {
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

// Count returns the total number of instance profiles.
func (s *InstanceProfileStore) Count() int {
	return s.bucket.Count()
}

// ListForRole returns all instance profiles that contain a specific role.
func (s *InstanceProfileStore) ListForRole(roleName string, marker string, maxItems int) (*InstanceProfileListResult, error) {
	if maxItems <= 0 {
		maxItems = 100
	}

	var profiles []*InstanceProfile
	count := 0
	started := marker == ""
	var hasMore bool

	err := s.bucket.ForEach(func(k, v []byte) error {
		var profile InstanceProfile
		if err := json.Unmarshal(v, &profile); err != nil {
			return err
		}

		hasRole := false
		for _, r := range profile.Roles {
			if r == roleName {
				hasRole = true
				break
			}
		}
		if !hasRole {
			return nil
		}

		if !started {
			if profile.InstanceProfileName == marker {
				started = true
			}
			return nil
		}

		if count < maxItems {
			profiles = append(profiles, &profile)
			count++
		} else {
			hasMore = true
		}
		return nil
	})

	if err != nil {
		return nil, NewStoreError("list_instance_profiles_for_role", err)
	}

	result := &InstanceProfileListResult{
		InstanceProfiles: profiles,
		IsTruncated:      hasMore,
	}
	if len(profiles) > 0 {
		result.Marker = profiles[len(profiles)-1].InstanceProfileName
	}

	return result, nil
}

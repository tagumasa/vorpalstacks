package iam

import (
	"errors"

	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// CreateInstanceProfileInput carries the parsed CreateInstanceProfile
// request.
type CreateInstanceProfileInput struct {
	InstanceProfileName string
	Path                string
	Tags                []tags.Tag
}

// InstanceProfileWithRoles pairs an instance profile with its resolved
// role records so the serialiser can populate the Roles member without
// store access.
type InstanceProfileWithRoles struct {
	Profile *iamstore.InstanceProfile
	Roles   []*iamstore.Role
}

// InstanceProfileListResult is the Core result for ListInstanceProfiles.
type InstanceProfileListResult struct {
	Profiles    []InstanceProfileWithRoles
	IsTruncated bool
	Marker      string
}

// createInstanceProfileCore validates input and creates an instance
// profile.
func (s *IAMService) createInstanceProfileCore(store *iamstore.IAMStore, input *CreateInstanceProfileInput) (*iamstore.InstanceProfile, error) {
	if input.InstanceProfileName == "" {
		return nil, NewInvalidInputError("InstanceProfileName", "cannot be empty")
	}
	if err := validateEntityName128(input.InstanceProfileName, "InstanceProfileName"); err != nil {
		return nil, err
	}

	path := input.Path
	if path == "" {
		path = "/"
	}
	if !validatePath(path) {
		return nil, NewInvalidInputError("Path", "must be a valid path starting and ending with /")
	}

	if err := validateNewTags(input.Tags); err != nil {
		return nil, err
	}

	profile, err := store.InstanceProfiles().Create(input.InstanceProfileName, path, s.accountID, input.Tags)
	if err != nil {
		if errors.Is(err, iamstore.ErrInstanceProfileAlreadyExists) {
			return nil, NewInstanceProfileAlreadyExistsError(input.InstanceProfileName)
		}
		return nil, err
	}
	return profile, nil
}

// resolveInstanceProfileRoles fetches the role records attached to the
// profile, preserving order and silently skipping roles that no longer
// exist.
func resolveInstanceProfileRoles(store *iamstore.IAMStore, profile *iamstore.InstanceProfile) []*iamstore.Role {
	roles := make([]*iamstore.Role, 0, len(profile.Roles))
	for _, roleName := range profile.Roles {
		if role, err := store.Roles().Get(roleName); err == nil {
			roles = append(roles, role)
		}
	}
	return roles
}

// getInstanceProfileCore validates input and retrieves the instance
// profile with its resolved roles.
func (s *IAMService) getInstanceProfileCore(store *iamstore.IAMStore, instanceProfileName string) (*InstanceProfileWithRoles, error) {
	if instanceProfileName == "" {
		return nil, NewValidationError("InstanceProfileName")
	}
	profile, err := store.InstanceProfiles().Get(instanceProfileName)
	if err != nil {
		return nil, NewNoSuchInstanceProfileError(instanceProfileName)
	}
	return &InstanceProfileWithRoles{
		Profile: profile,
		Roles:   resolveInstanceProfileRoles(store, profile),
	}, nil
}

// deleteInstanceProfileCore validates input and deletes an instance
// profile that has no roles attached.
func (s *IAMService) deleteInstanceProfileCore(store *iamstore.IAMStore, instanceProfileName string) error {
	if instanceProfileName == "" {
		return NewValidationError("InstanceProfileName")
	}
	profile, err := store.InstanceProfiles().Get(instanceProfileName)
	if err != nil {
		return NewNoSuchInstanceProfileError(instanceProfileName)
	}

	if len(profile.Roles) > 0 {
		return NewDeleteInstanceProfileConflictError(instanceProfileName)
	}

	return store.InstanceProfiles().Delete(instanceProfileName)
}

// listInstanceProfilesCore lists the instance profiles with their resolved
// roles.
func (s *IAMService) listInstanceProfilesCore(store *iamstore.IAMStore, pathPrefix, marker string, maxItems int) (*InstanceProfileListResult, error) {
	result, err := store.InstanceProfiles().List(pathPrefix, marker, maxItems)
	if err != nil {
		return nil, err
	}

	profiles := make([]InstanceProfileWithRoles, len(result.InstanceProfiles))
	for i, profile := range result.InstanceProfiles {
		profiles[i] = InstanceProfileWithRoles{
			Profile: profile,
			Roles:   resolveInstanceProfileRoles(store, profile),
		}
	}
	return &InstanceProfileListResult{
		Profiles:    profiles,
		IsTruncated: result.IsTruncated,
		Marker:      result.Marker,
	}, nil
}

// listInstanceProfilesForRoleCore lists the instance profiles associated
// with the given role, with their resolved roles.
func (s *IAMService) listInstanceProfilesForRoleCore(store *iamstore.IAMStore, roleName, marker string, maxItems int) (*InstanceProfileListResult, error) {
	if roleName == "" {
		return nil, NewValidationError("RoleName")
	}
	if !store.Roles().Exists(roleName) {
		return nil, NewNoSuchRoleError(roleName)
	}

	result, err := store.InstanceProfiles().ListForRole(roleName, marker, maxItems)
	if err != nil {
		return nil, err
	}

	profiles := make([]InstanceProfileWithRoles, len(result.InstanceProfiles))
	for i, profile := range result.InstanceProfiles {
		profiles[i] = InstanceProfileWithRoles{
			Profile: profile,
			Roles:   resolveInstanceProfileRoles(store, profile),
		}
	}
	return &InstanceProfileListResult{
		Profiles:    profiles,
		IsTruncated: result.IsTruncated,
		Marker:      result.Marker,
	}, nil
}

// addRoleToInstanceProfileCore validates input and attaches the role to
// the instance profile.  AWS enforces a maximum of one role per instance
// profile; the limit is enforced atomically inside the store layer.
func (s *IAMService) addRoleToInstanceProfileCore(store *iamstore.IAMStore, instanceProfileName, roleName string) error {
	if instanceProfileName == "" {
		return NewValidationError("InstanceProfileName")
	}
	if roleName == "" {
		return NewValidationError("RoleName")
	}

	if !store.InstanceProfiles().Exists(instanceProfileName) {
		return NewNoSuchInstanceProfileError(instanceProfileName)
	}

	if !store.Roles().Exists(roleName) {
		return NewNoSuchRoleError(roleName)
	}

	if err := store.InstanceProfiles().AddRole(instanceProfileName, roleName); err != nil {
		if errors.Is(err, iamstore.ErrRoleAlreadyInInstanceProfile) {
			return NewRoleAlreadyInInstanceProfileError(roleName, instanceProfileName)
		}
		if errors.Is(err, iamstore.ErrInstanceProfileRoleLimit) {
			return ErrInstanceProfileRoleLimit
		}
		return err
	}
	return nil
}

// removeRoleFromInstanceProfileCore validates input and detaches the role
// from the instance profile.
func (s *IAMService) removeRoleFromInstanceProfileCore(store *iamstore.IAMStore, instanceProfileName, roleName string) error {
	if instanceProfileName == "" {
		return NewValidationError("InstanceProfileName")
	}
	if roleName == "" {
		return NewValidationError("RoleName")
	}

	if !store.InstanceProfiles().Exists(instanceProfileName) {
		return NewNoSuchInstanceProfileError(instanceProfileName)
	}

	if !store.Roles().Exists(roleName) {
		return NewNoSuchRoleError(roleName)
	}

	if err := store.InstanceProfiles().RemoveRole(instanceProfileName, roleName); err != nil {
		return NewRoleNotInInstanceProfileError(roleName, instanceProfileName)
	}
	return nil
}

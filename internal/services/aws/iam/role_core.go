// Transport-agnostic Core functions for IAM roles: validation and store
// operations shared by the AWS-compatible HTTP API handlers and the admin
// gRPC-Web handler (the xxxCore pattern).
package iam

import (
	"errors"
	"fmt"

	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// CreateRoleInput holds the parameters for creating an IAM role.
type CreateRoleInput struct {
	RoleName                 string
	Path                     string
	AssumeRolePolicyDocument string
	Description              string
	MaxSessionDuration       int
	PermissionsBoundaryArn   string
	Tags                     []tags.Tag
}

// UpdateRoleInput holds the parameters for updating an IAM role's
// description and/or max session duration.
type UpdateRoleInput struct {
	RoleName           string
	Description        string
	MaxSessionDuration int
}

// DeleteRoleInput holds the parameters for deleting an IAM role.
// When Cascade is true, all dependent resources are removed before the role
// record (admin handler behaviour).  When false, a DeleteConflict error is
// returned if any dependent resource remains (AWS API behaviour).
type DeleteRoleInput struct {
	RoleName string
	Cascade  bool
}

// createRoleCore validates input and creates an IAM role in the store.
// Returns the created role or an IAM-formatted error.
func (s *IAMService) createRoleCore(store *iamstore.IAMStore, input *CreateRoleInput) (*iamstore.Role, error) {
	if input.RoleName == "" {
		return nil, NewInvalidInputError("RoleName", "cannot be empty")
	}
	if err := validateEntityName(input.RoleName, "RoleName"); err != nil {
		return nil, err
	}

	path := input.Path
	if path == "" {
		path = "/"
	}
	if !validatePath(path) {
		return nil, NewInvalidInputError("Path", "must be a valid path starting and ending with /")
	}

	if input.AssumeRolePolicyDocument == "" {
		return nil, ErrMalformedPolicyDocument
	}
	if !validateTrustPolicyDocument(input.AssumeRolePolicyDocument) {
		return nil, ErrMalformedPolicyDocument
	}

	if !validateRoleDescription(input.Description) {
		return nil, NewInvalidInputError("Description", "must be 0 to 1000 characters; allowed: tab, LF, CR, printable ASCII, Latin-1 supplement")
	}

	maxSessionDuration := input.MaxSessionDuration
	if maxSessionDuration == 0 {
		maxSessionDuration = defaultRoleSessionDuration
	}
	if !validateRoleMaxSessionDuration(maxSessionDuration) {
		return nil, NewInvalidInputError("MaxSessionDuration", fmt.Sprintf("must be between %d and %d seconds", minRoleSessionDuration, maxRoleSessionDuration))
	}

	if err := validateNewTags(input.Tags); err != nil {
		return nil, err
	}

	role, err := store.Roles().Create(
		input.RoleName, path, store.AccountID(),
		input.AssumeRolePolicyDocument, input.Description,
		maxSessionDuration, input.Tags,
	)
	if err != nil {
		if errors.Is(err, iamstore.ErrRoleAlreadyExists) {
			return nil, NewRoleAlreadyExistsError(input.RoleName)
		}
		return nil, err
	}

	// Apply permissions boundary if specified at creation time (Smithy
	// CreateRoleInput.PermissionsBoundary).
	if input.PermissionsBoundaryArn != "" {
		if err := putRolePermissionsBoundaryCore(store, role, input.PermissionsBoundaryArn); err != nil {
			return nil, err
		}
	}

	return role, nil
}

// putRolePermissionsBoundaryCore is the role equivalent of
// putUserPermissionsBoundaryCore.
func putRolePermissionsBoundaryCore(store *iamstore.IAMStore, role *iamstore.Role, pbArn string) error {
	if err := validateIAMPolicyArn(pbArn); err != nil {
		return err
	}
	if !store.Policies().Exists(pbArn) {
		return NewNoSuchPolicyError(pbArn)
	}
	if role.PermissionsBoundary != nil && role.PermissionsBoundary.PermissionsBoundaryArn == pbArn {
		return nil
	}
	if role.PermissionsBoundary != nil && role.PermissionsBoundary.PermissionsBoundaryArn != "" {
		_ = store.Policies().DecrementPermissionsBoundaryUsageCount(role.PermissionsBoundary.PermissionsBoundaryArn)
	}
	role.PermissionsBoundary = &iamstore.PermissionsBoundary{
		PermissionsBoundaryType: "Policy",
		PermissionsBoundaryArn:  pbArn,
	}
	if err := store.Roles().Put(role); err != nil {
		return err
	}
	_ = store.Policies().IncrementPermissionsBoundaryUsageCount(pbArn)
	return nil
}

// getRoleCore returns the IAM role with the given name.  Callers must
// validate that roleName is non-empty before calling.
func (s *IAMService) getRoleCore(store *iamstore.IAMStore, roleName string) (*iamstore.Role, error) {
	role, err := store.Roles().Get(roleName)
	if err != nil {
		return nil, NewNoSuchRoleError(roleName)
	}
	return role, nil
}

// listRolesCore returns a paginated list of IAM roles.
func (s *IAMService) listRolesCore(store *iamstore.IAMStore, pathPrefix, marker string, maxItems int) (*iamstore.RoleListResult, error) {
	return store.Roles().List(pathPrefix, marker, maxItems)
}

// updateAssumeRolePolicyCore validates the trust policy document and
// replaces the role's AssumeRolePolicyDocument field.  Used by the HTTP
// API UpdateAssumeRolePolicy operation; routing through core keeps the
// validation identical if a future admin handler mirrors the op.
func (s *IAMService) updateAssumeRolePolicyCore(store *iamstore.IAMStore, roleName, policyDocument string) error {
	role, err := s.getRoleCore(store, roleName)
	if err != nil {
		return err
	}
	if !validateTrustPolicyDocument(policyDocument) {
		return ErrMalformedPolicyDocument
	}
	role.AssumeRolePolicyDocument = policyDocument
	return store.Roles().Put(role)
}

// updateRoleCore validates input and updates an IAM role's description and/or
// max session duration.  Returns the updated role or an IAM-formatted error.
func (s *IAMService) updateRoleCore(store *iamstore.IAMStore, input *UpdateRoleInput) (*iamstore.Role, error) {
	if input.RoleName == "" {
		return nil, NewValidationError("RoleName")
	}

	maxSessionDuration := input.MaxSessionDuration
	if maxSessionDuration > 0 {
		if !validateRoleMaxSessionDuration(maxSessionDuration) {
			return nil, NewInvalidInputError("MaxSessionDuration", fmt.Sprintf("must be between %d and %d seconds", minRoleSessionDuration, maxRoleSessionDuration))
		}
	}

	if !validateRoleDescription(input.Description) {
		return nil, NewInvalidInputError("Description", "must be 0 to 1000 characters; allowed: tab, LF, CR, printable ASCII, Latin-1 supplement")
	}

	if err := store.UpdateRoleFields(input.RoleName, input.Description, maxSessionDuration); err != nil {
		return nil, err
	}

	role, err := store.Roles().Get(input.RoleName)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// deleteRoleCore validates input and deletes an IAM role.
// When input.Cascade is true, all dependent resources (inline/attached
// policies, instance profile associations) are removed first via
// cascadeDeleteRole.  When false, a DeleteConflict error is returned if any
// dependent resource remains.
// In both modes the permissions-boundary usage count is decremented
// best-effort before the role record is removed.
func (s *IAMService) deleteRoleCore(store *iamstore.IAMStore, input *DeleteRoleInput) error {
	if input.RoleName == "" {
		return ErrNoSuchRole
	}
	if !store.Roles().Exists(input.RoleName) {
		return NewNoSuchRoleError(input.RoleName)
	}

	if input.Cascade {
		return cascadeDeleteRole(store, input.RoleName)
	}

	// Conflict detection — AWS API path.
	instanceProfiles, err := store.InstanceProfiles().ListForRole(input.RoleName, "", 1)
	if err != nil {
		return err
	}
	if len(instanceProfiles.InstanceProfiles) > 0 {
		return NewDeleteRoleConflictError("Cannot delete entity, must remove role from instance profile first.")
	}

	inlinePolicies, err := store.InlinePolicies().List(PrincipalTypeRole, input.RoleName)
	if err != nil {
		return err
	}
	if len(inlinePolicies) > 0 {
		return NewDeleteRoleConflictError("Cannot delete entity, must delete policies first.")
	}

	attachedPolicies, err := store.AttachedPolicies().ListAttachedPolicies(PrincipalTypeRole, input.RoleName)
	if err != nil {
		return err
	}
	if len(attachedPolicies) > 0 {
		return NewDeleteRoleConflictError("Cannot delete entity, must detach policies first.")
	}

	// Decrement permissions boundary usage count before the role record is
	// removed. AWS allows deleting an entity that still has a permissions
	// boundary attached (it is not a deletion prerequisite), so the policy
	// counter must be adjusted here to avoid drift. Best-effort, matching the
	// PutRolePermissionsBoundary / DeleteRolePermissionsBoundary pattern.
	if role, gErr := store.Roles().Get(input.RoleName); gErr == nil {
		if role.PermissionsBoundary != nil && role.PermissionsBoundary.PermissionsBoundaryArn != "" {
			_ = store.Policies().DecrementPermissionsBoundaryUsageCount(role.PermissionsBoundary.PermissionsBoundaryArn)
		}
	}

	return store.Roles().Delete(input.RoleName)
}

// deleteRolePermissionsBoundaryCore removes the permissions boundary
// from a role and decrements the usage count on the previously-bound
// policy.  Mirrors deleteUserPermissionsBoundaryCore for roles.
func (s *IAMService) deleteRolePermissionsBoundaryCore(store *iamstore.IAMStore, roleName string) error {
	role, err := store.Roles().Get(roleName)
	if err != nil {
		return NewNoSuchRoleError(roleName)
	}

	if role.PermissionsBoundary != nil && role.PermissionsBoundary.PermissionsBoundaryArn != "" {
		_ = store.Policies().DecrementPermissionsBoundaryUsageCount(role.PermissionsBoundary.PermissionsBoundaryArn)
	}

	role.PermissionsBoundary = nil
	return store.Roles().Put(role)
}

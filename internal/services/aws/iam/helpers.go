package iam

import (
	"errors"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// resolveUserName returns userName if non-empty, otherwise defaults to the
// caller's IAM principal name.  Per AWS spec, several IAM operations allow
// UserName to be omitted, in which case it defaults to the authenticated
// caller.  If the caller is not an IAM user (e.g. anonymous), an error is
// returned.
func resolveUserName(reqCtx *request.RequestContext, userName string) (string, error) {
	if userName != "" {
		return userName, nil
	}
	if reqCtx.PrincipalType == request.PrincipalTypeUser && reqCtx.Principal != "" {
		return reqCtx.Principal, nil
	}
	// No user name is available to report; the caller's principal (possibly
	// empty for anonymous callers) is interpolated so no placeholder text
	// reaches the wire.
	return "", NewNoSuchUserError(reqCtx.Principal)
}

// buildAttachedManagedPolicies returns the attached managed policy list for
// a principal, collapsing the copy-paste across User/Group/Role in
// GetAccountAuthorizationDetails. Store read failures abort the response:
// the authorisation details must never silently omit a policy.
func buildAttachedManagedPolicies(store *iamstore.IAMStore, principalType, principalName string) ([]interface{}, error) {
	arns, err := store.AttachedPolicies().ListAttachedPolicies(principalType, principalName)
	if err != nil {
		return nil, storeListError(err)
	}
	policies := make([]interface{}, 0, len(arns))
	for _, arn := range arns {
		p, err := store.Policies().Get(arn)
		if err != nil {
			return nil, storeListError(err)
		}
		policies = append(policies, map[string]interface{}{
			"PolicyName": p.PolicyName,
			"PolicyArn":  p.Arn,
		})
	}
	return policies, nil
}

// buildInlinePolicyList returns the inline policy name list for a principal;
// a listing failure aborts the response rather than silently truncating it.
func buildInlinePolicyList(store *iamstore.IAMStore, principalType, principalName string) ([]interface{}, error) {
	names, err := store.InlinePolicies().List(principalType, principalName)
	if err != nil {
		return nil, storeListError(err)
	}
	list := make([]interface{}, 0, len(names))
	for _, pn := range names {
		list = append(list, map[string]interface{}{
			"PolicyName": pn,
		})
	}
	return list, nil
}

// attachPermissionsBoundaryCore runs the boundary-attach sequence shared by
// the user and role cores: validate the policy ARN, require the policy to
// exist, keep same-ARN attaches idempotent, set the field, persist, and only
// then move the usage tallies (decrement the previous boundary's policy,
// increment the new one).  boundary returns the address of the entity's
// PermissionsBoundary field; put persists the entity.
func attachPermissionsBoundaryCore[T any](store *iamstore.IAMStore, entity *T, pbArn string,
	boundary func(*T) **iamstore.PermissionsBoundary, put func(*T) error) error {
	if err := validateIAMPolicyArn(pbArn); err != nil {
		return err
	}
	// The policy record is resolved rather than probed: the resolved read
	// keeps an outage from masquerading as a missing policy.
	if _, err := store.Policies().Get(pbArn); err != nil {
		return storeReadError(err, iamstore.ErrPolicyNotFound, NewNoSuchPolicyError(pbArn))
	}
	field := boundary(entity)
	// Idempotent: same ARN already set — nothing to do.
	if *field != nil && (*field).PermissionsBoundaryArn == pbArn {
		return nil
	}
	previousArn := ""
	if *field != nil {
		previousArn = (*field).PermissionsBoundaryArn
	}
	*field = &iamstore.PermissionsBoundary{
		PermissionsBoundaryType: "Policy",
		PermissionsBoundaryArn:  pbArn,
	}
	if err := put(entity); err != nil {
		return err
	}
	// The tallies move only after the entity persisted: a failed put
	// leaves the entity and its usage counts exactly as they were.  A
	// tally failure after a successful put is logged and left best-effort
	// (an undercounting policy must never block the boundary change
	// itself); a failed increment drifts the count downwards, which
	// DeletePolicy relies on, so it is always surfaced in the log.
	if previousArn != "" {
		if err := store.Policies().DecrementPermissionsBoundaryUsageCount(previousArn); err != nil {
			logs.Warn("iam: previous permissions-boundary usage-count decrement failed; the policy's tally overcounts",
				logs.String("policyArn", previousArn), logs.Err(err))
		}
	}
	if err := store.Policies().IncrementPermissionsBoundaryUsageCount(pbArn); err != nil {
		logs.Warn("iam: permissions-boundary usage-count increment failed; the policy's tally undercounts",
			logs.String("policyArn", pbArn), logs.Err(err))
	}
	return nil
}

// deletePermissionsBoundaryCore clears a resolved entity's boundary field,
// decrementing the previously-bound policy's usage count, and persists the
// entity.  The caller resolves the entity and maps its not-found error; the
// wire-order member validation stays in the per-operation cores.  The tally
// moves only after the entity persisted, mirroring the attach sequence.
func deletePermissionsBoundaryCore[T any](store *iamstore.IAMStore, entity *T,
	boundary func(*T) **iamstore.PermissionsBoundary, put func(*T) error) error {
	field := boundary(entity)
	previous := *field
	*field = nil
	if err := put(entity); err != nil {
		return err
	}
	decrementBoundaryUsageCount(store, previous)
	return nil
}

// updateEntityCore runs the rename/repath sequence shared by UpdateGroup and
// UpdateUser: the current name is required, at least one of NewPath or
// New<Entity>Name must be given, the path and new-name predicates run in
// that order, the store rename's already-exists fault maps to the entity's
// own wire error, and the entity is re-read under its post-rename name.
// rename closes over the operation's input; the sentinel and wire factory
// arrive separately because the two entities model their collisions
// differently, and validateNewName carries each entity's name rule
// (128-character group names, 64-character user names).
func updateEntityCore[T any](entity, name, newPath, newName string,
	rename func() error,
	alreadyExistsSentinel error,
	alreadyExistsErr func(string) error,
	validateNewName func(string, string) error,
	get func(string) (*T, error)) (*T, error) {

	if name == "" {
		return nil, NewValidationError(entity + "Name")
	}

	if newPath == "" && newName == "" {
		return nil, NewInvalidInputError("Update"+entity, "at least one of NewPath or New"+entity+"Name must be specified")
	}
	if newPath != "" && !validatePath(newPath) {
		return nil, NewInvalidInputError("NewPath", "must be a valid path starting and ending with /")
	}
	if newName != "" {
		if err := validateNewName(newName, "New"+entity+"Name"); err != nil {
			return nil, err
		}
	}

	if err := rename(); err != nil {
		if errors.Is(err, alreadyExistsSentinel) {
			return nil, alreadyExistsErr(newName)
		}
		return nil, err
	}

	targetName := name
	if newName != "" {
		targetName = newName
	}
	return get(targetName)
}

// paginateAll exhausts a marker-paginated listing, following next markers
// until the final page.  page fetches one page given the current marker and
// reports whether a further page exists; the listAll* helpers below differ
// only in that closure.
func paginateAll[T any](page func(marker string) ([]T, bool, string, error)) ([]T, error) {
	var all []T
	marker := ""
	for {
		items, truncated, next, err := page(marker)
		if err != nil {
			return nil, err
		}
		all = append(all, items...)
		if !truncated {
			return all, nil
		}
		marker = next
	}
}

// listAllUsers paginates through all users matching pathPrefix.
func listAllUsers(store *iamstore.IAMStore, pathPrefix string) ([]*iamstore.User, error) {
	return paginateAll(func(marker string) ([]*iamstore.User, bool, string, error) {
		result, err := store.Users().List(pathPrefix, marker, pagination.AbsoluteMaxItems)
		if err != nil {
			return nil, false, "", err
		}
		return result.Users, result.IsTruncated, result.Marker, nil
	})
}

// listAllGroups paginates through all groups matching pathPrefix.
func listAllGroups(store *iamstore.IAMStore, pathPrefix string) ([]*iamstore.Group, error) {
	return paginateAll(func(marker string) ([]*iamstore.Group, bool, string, error) {
		result, err := store.Groups().List(pathPrefix, marker, pagination.AbsoluteMaxItems)
		if err != nil {
			return nil, false, "", err
		}
		return result.Groups, result.IsTruncated, result.Marker, nil
	})
}

// listAllRoles paginates through all roles matching pathPrefix.
func listAllRoles(store *iamstore.IAMStore, pathPrefix string) ([]*iamstore.Role, error) {
	return paginateAll(func(marker string) ([]*iamstore.Role, bool, string, error) {
		result, err := store.Roles().List(pathPrefix, marker, pagination.AbsoluteMaxItems)
		if err != nil {
			return nil, false, "", err
		}
		return result.Roles, result.IsTruncated, result.Marker, nil
	})
}

// listAllPolicies paginates through all policies matching the given filters.
func listAllPolicies(store *iamstore.IAMStore, scope, pathPrefix string, onlyAttached bool) ([]*iamstore.Policy, error) {
	return paginateAll(func(marker string) ([]*iamstore.Policy, bool, string, error) {
		result, err := store.Policies().List(scope, pathPrefix, onlyAttached, marker, pagination.AbsoluteMaxItems)
		if err != nil {
			return nil, false, "", err
		}
		return result.Policies, result.IsTruncated, result.Marker, nil
	})
}

// detachAllAttachedPolicies detaches every managed policy attached to a
// principal, decrementing each policy's attachment count.  The user, role
// and group cascade deletions share it.
func detachAllAttachedPolicies(store *iamstore.IAMStore, principalType, principalName string) error {
	attachedPolicies, err := store.AttachedPolicies().ListAttachedPolicies(principalType, principalName)
	if err != nil {
		return err
	}
	for _, policyArn := range attachedPolicies {
		if err := store.AttachedPolicies().Detach(principalType, principalName, policyArn); err != nil {
			return err
		}
		if err := store.Policies().DecrementAttachmentCount(policyArn); err != nil {
			return err
		}
	}
	return nil
}

// decrementBoundaryUsageCount adjusts the permissions-boundary usage count
// for a principal being deleted: AWS allows deleting an entity that still
// has a permissions boundary attached (it is not a deletion prerequisite),
// so the policy counter must be decremented here to avoid drift.
// Best-effort, matching the attach/delete boundary cores: a failure is
// logged (an undercounting policy can pass DeletePolicy while referenced)
// and never blocks the deletion.
func decrementBoundaryUsageCount(store *iamstore.IAMStore, boundary *iamstore.PermissionsBoundary) {
	if boundary != nil && boundary.PermissionsBoundaryArn != "" {
		if err := store.Policies().DecrementPermissionsBoundaryUsageCount(boundary.PermissionsBoundaryArn); err != nil {
			logs.Warn("iam: permissions-boundary usage-count decrement on entity deletion failed; the policy's tally undercounts",
				logs.String("policyArn", boundary.PermissionsBoundaryArn), logs.Err(err))
		}
	}
}

// cascadeDeleteUser removes all resources associated with a user before deleting
// the user record. Used by the admin handler to ensure consistent cleanup.
func cascadeDeleteUser(store *iamstore.IAMStore, userName string) error {
	// Decrement permissions boundary usage count before the cascade removes
	// the user record (see decrementBoundaryUsageCount).
	if user, gErr := store.Users().Get(userName); gErr == nil {
		decrementBoundaryUsageCount(store, user.PermissionsBoundary)
	}

	if store.LoginProfiles().Exists(userName) {
		if err := store.LoginProfiles().Delete(userName); err != nil {
			return err
		}
	}

	if err := store.InlinePolicies().DeleteAllForPrincipal(PrincipalTypeUser, userName); err != nil {
		return err
	}

	if err := detachAllAttachedPolicies(store, PrincipalTypeUser, userName); err != nil {
		return err
	}

	if err := store.UserGroups().RemoveAllGroupsForUser(userName); err != nil {
		return err
	}

	if err := store.SigningCertificates().DeleteAllForUser(userName); err != nil {
		return err
	}

	if err := store.SSHPublicKeys().DeleteAllForUser(userName); err != nil {
		return err
	}

	if err := store.ServiceSpecificCredentials().DeleteAllForUser(userName); err != nil {
		return err
	}

	if err := store.AccessKeys().DeleteByUserName(userName); err != nil {
		return err
	}

	mfaResult, err := store.MFADevices().ListForUser(userName, "", 1000)
	if err != nil {
		return err
	}
	for _, device := range mfaResult.MFADevices {
		if err := store.MFADevices().Deactivate(device.SerialNumber); err != nil {
			return err
		}
	}

	return store.Users().Delete(userName)
}

// cascadeDeleteRole removes all resources associated with a role before deleting
// the role record. Used by the admin handler to ensure consistent cleanup.
func cascadeDeleteRole(store *iamstore.IAMStore, roleName string) error {
	// Decrement permissions boundary usage count before the cascade removes
	// the role record (see decrementBoundaryUsageCount).
	if role, gErr := store.Roles().Get(roleName); gErr == nil {
		decrementBoundaryUsageCount(store, role.PermissionsBoundary)
	}

	if err := store.InlinePolicies().DeleteAllForPrincipal(PrincipalTypeRole, roleName); err != nil {
		return err
	}

	if err := detachAllAttachedPolicies(store, PrincipalTypeRole, roleName); err != nil {
		return err
	}

	instanceProfiles, err := store.InstanceProfiles().ListForRole(roleName, "", 1000)
	if err != nil {
		return err
	}
	for _, ip := range instanceProfiles.InstanceProfiles {
		if err := store.InstanceProfiles().RemoveRole(ip.InstanceProfileName, roleName); err != nil {
			return err
		}
	}

	return store.Roles().Delete(roleName)
}

// cascadeDeleteGroup removes all resources associated with a group before deleting
// the group record. Used by the admin handler to ensure consistent cleanup.
func cascadeDeleteGroup(store *iamstore.IAMStore, groupName string) error {
	if err := store.InlinePolicies().DeleteAllForPrincipal(PrincipalTypeGroup, groupName); err != nil {
		return err
	}

	if err := detachAllAttachedPolicies(store, PrincipalTypeGroup, groupName); err != nil {
		return err
	}

	if err := store.UserGroups().RemoveAllUsersFromGroup(groupName); err != nil {
		return err
	}

	return store.Groups().Delete(groupName)
}

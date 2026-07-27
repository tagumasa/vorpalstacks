// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"
	"encoding/json"
	"strconv"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/aws/types"
)

const (
	// MaxAccessKeysPerUser is the maximum number of access keys a user can have.
	MaxAccessKeysPerUser = 2
	// MaxPolicyVersions is the maximum number of policy versions allowed.
	MaxPolicyVersions = 5
	// MaxTagsPerResource is the maximum number of tags allowed on a single IAM resource.
	MaxTagsPerResource = 50
	// MaxTagKeyLength is the maximum length of a tag key.
	MaxTagKeyLength = 128
	// MaxTagValueLength is the maximum length of a tag value.
	MaxTagValueLength = 256
)

// policyValidationMode controls which fields are required for each statement.
type policyValidationMode int

const (
	// policyModeManaged is for identity-based policies (managed and inline).
	// Requires Effect + Action/NotAction + Resource/NotResource per statement.
	policyModeManaged policyValidationMode = iota
	// policyModeTrust is for AssumeRolePolicyDocument (resource-based trust
	// policies). Requires Effect + Action/NotAction + Principal/NotPrincipal.
	policyModeTrust
)

// validatePolicyDocument checks if a policy document is valid JSON and has
// the minimum required structure for an IAM identity-based policy: a
// top-level object with a "Statement" field containing at least one
// statement object with Effect, Action/NotAction, and Resource/NotResource.
func validatePolicyDocument(document string) bool {
	return validatePolicyDocumentMode(document, policyModeManaged)
}

// validateTrustPolicyDocument validates an AssumeRolePolicyDocument (trust
// policy).  Each statement must have Effect, Action/NotAction, and
// Principal/NotPrincipal.
func validateTrustPolicyDocument(document string) bool {
	return validatePolicyDocumentMode(document, policyModeTrust)
}

func validatePolicyDocumentMode(document string, mode policyValidationMode) bool {
	if document == "" {
		return false
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(document), &raw); err != nil {
		return false
	}
	statementsRaw, ok := raw["Statement"]
	if !ok {
		return false
	}

	// Statement can be a single object or an array of objects.
	var singleStmt map[string]interface{}
	if err := json.Unmarshal(statementsRaw, &singleStmt); err == nil {
		return validateStatement(singleStmt, mode)
	}

	var stmtArray []map[string]interface{}
	if err := json.Unmarshal(statementsRaw, &stmtArray); err != nil {
		return false
	}
	if len(stmtArray) == 0 {
		return false
	}
	for _, stmt := range stmtArray {
		if !validateStatement(stmt, mode) {
			return false
		}
	}
	return true
}

// validateStatement checks that a single policy statement has the required
// members for the given validation mode:
//   - Effect must be "Allow" or "Deny"
//   - Action or NotAction must be present
//   - For managed policies: Resource or NotResource must be present
//   - For trust policies: Principal or NotPrincipal must be present
func validateStatement(stmt map[string]interface{}, mode policyValidationMode) bool {
	effect, ok := stmt["Effect"].(string)
	if !ok {
		return false
	}
	if effect != "Allow" && effect != "Deny" {
		return false
	}

	// Action or NotAction is required in all policy statements.
	if !hasPolicyKey(stmt, "Action") && !hasPolicyKey(stmt, "NotAction") {
		return false
	}

	switch mode {
	case policyModeManaged:
		// Identity-based policies require Resource or NotResource.
		if !hasPolicyKey(stmt, "Resource") && !hasPolicyKey(stmt, "NotResource") {
			return false
		}
	case policyModeTrust:
		// Trust policies require Principal or NotPrincipal.
		if !hasPolicyKey(stmt, "Principal") && !hasPolicyKey(stmt, "NotPrincipal") {
			return false
		}
	}
	return true
}

// hasPolicyKey returns true if the statement map contains the given key
// with a non-nil value.  A key set to JSON null is treated as absent.
func hasPolicyKey(stmt map[string]interface{}, key string) bool {
	val, ok := stmt[key]
	if !ok {
		return false
	}
	return val != nil
}

// validateTagEntries validates the key and value length limits for each
// individual tag entry.  It does NOT check the total tag count — the caller
// is responsible for that because the acceptable count depends on context
// (new tags on Create vs. merged tags on TagResource).
func validateTagEntries(newTags []types.Tag) error {
	for _, t := range newTags {
		if len(t.Key) == 0 || len(t.Key) > MaxTagKeyLength {
			return NewInvalidInputError("TagKey", "must be 1 to "+strconv.Itoa(MaxTagKeyLength)+" characters")
		}
		if !tagKeyPattern.MatchString(t.Key) {
			return NewInvalidInputError("TagKey", "contains invalid characters")
		}
		if len(t.Value) > MaxTagValueLength {
			return NewInvalidInputError("TagValue", "must be 0 to "+strconv.Itoa(MaxTagValueLength)+" characters")
		}
	}
	return nil
}

// validateNewTags validates both per-tag entry limits and the total tag
// count for resources being created.  On Create operations there are no
// pre-existing tags, so the total count is simply len(newTags).
func validateNewTags(newTags []types.Tag) error {
	if err := validateTagEntries(newTags); err != nil {
		return err
	}
	if len(newTags) > MaxTagsPerResource {
		return NewInvalidInputError("Tags", "exceeds maximum of "+strconv.Itoa(MaxTagsPerResource)+" tags per resource")
	}
	return nil
}

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
	return "", ErrNoSuchUser
}

type tagOps[T any] struct {
	paramName  string
	emptyErr   error
	notFoundFn func(string) error
	getFn      func(*iamstore.IAMStore, string) (T, error)
	putFn      func(*iamstore.IAMStore, T) error
	tagsFn     func(T) *[]types.Tag
}

func tagResource[T any](ctx context.Context, s *IAMService, reqCtx *request.RequestContext, req *request.ParsedRequest, ops tagOps[T]) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, ops.paramName)
	if name == "" {
		return nil, ops.emptyErr
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	res, err := ops.getFn(store, name)
	if err != nil {
		return nil, ops.notFoundFn(name)
	}
	currentTags := ops.tagsFn(res)
	newTags := tags.ParseTagsWithQueryFallback(req.Parameters, "Tags")
	if err := validateTagEntries(newTags); err != nil {
		return nil, err
	}
	merged := tags.Apply(*currentTags, newTags)
	if len(merged) > MaxTagsPerResource {
		return nil, NewInvalidInputError("Tags", "exceeds maximum of "+strconv.Itoa(MaxTagsPerResource)+" tags per resource")
	}
	*currentTags = merged
	if err := ops.putFn(store, res); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func untagResource[T any](ctx context.Context, s *IAMService, reqCtx *request.RequestContext, req *request.ParsedRequest, ops tagOps[T]) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, ops.paramName)
	if name == "" {
		return nil, ops.emptyErr
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	res, err := ops.getFn(store, name)
	if err != nil {
		return nil, ops.notFoundFn(name)
	}
	currentTags := ops.tagsFn(res)
	*currentTags = tags.RemoveByTagKeys(*currentTags, tags.ParseTagKeysWithQueryFallback(req.Parameters, "TagKeys"))
	if err := ops.putFn(store, res); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func listResourceTags[T any](ctx context.Context, s *IAMService, reqCtx *request.RequestContext, req *request.ParsedRequest, ops tagOps[T]) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, ops.paramName)
	if name == "" {
		return nil, ops.emptyErr
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	res, err := ops.getFn(store, name)
	if err != nil {
		return nil, ops.notFoundFn(name)
	}
	return map[string]interface{}{
		"Tags":        tags.ToResponse(*ops.tagsFn(res)),
		"IsTruncated": false,
	}, nil
}

// buildAttachedManagedPolicies returns the attached managed policy list for
// a principal, collapsing the copy-paste across User/Group/Role in
// GetAccountAuthorizationDetails.
func buildAttachedManagedPolicies(store *iamstore.IAMStore, principalType, principalName string) []interface{} {
	arns, _ := store.AttachedPolicies().ListAttachedPolicies(principalType, principalName)
	policies := make([]interface{}, 0, len(arns))
	for _, arn := range arns {
		if p, err := store.Policies().Get(arn); err == nil {
			policies = append(policies, map[string]interface{}{
				"PolicyName": p.PolicyName,
				"PolicyArn":  p.Arn,
			})
		}
	}
	return policies
}

// buildInlinePolicyList returns the inline policy name list for a principal.
func buildInlinePolicyList(store *iamstore.IAMStore, principalType, principalName string) []interface{} {
	names, _ := store.InlinePolicies().List(principalType, principalName)
	list := make([]interface{}, 0, len(names))
	for _, pn := range names {
		list = append(list, map[string]interface{}{
			"PolicyName": pn,
		})
	}
	return list
}

// listAllUsers paginates through all users matching pathPrefix.
func listAllUsers(store *iamstore.IAMStore, pathPrefix string) ([]*iamstore.User, error) {
	var all []*iamstore.User
	marker := ""
	for {
		result, err := store.Users().List(pathPrefix, marker, 1000)
		if err != nil {
			return nil, err
		}
		all = append(all, result.Users...)
		if !result.IsTruncated {
			break
		}
		marker = result.Marker
	}
	return all, nil
}

// listAllGroups paginates through all groups matching pathPrefix.
func listAllGroups(store *iamstore.IAMStore, pathPrefix string) ([]*iamstore.Group, error) {
	var all []*iamstore.Group
	marker := ""
	for {
		result, err := store.Groups().List(pathPrefix, marker, 1000)
		if err != nil {
			return nil, err
		}
		all = append(all, result.Groups...)
		if !result.IsTruncated {
			break
		}
		marker = result.Marker
	}
	return all, nil
}

// listAllRoles paginates through all roles matching pathPrefix.
func listAllRoles(store *iamstore.IAMStore, pathPrefix string) ([]*iamstore.Role, error) {
	var all []*iamstore.Role
	marker := ""
	for {
		result, err := store.Roles().List(pathPrefix, marker, 1000)
		if err != nil {
			return nil, err
		}
		all = append(all, result.Roles...)
		if !result.IsTruncated {
			break
		}
		marker = result.Marker
	}
	return all, nil
}

// listAllPolicies paginates through all policies matching the given filters.
func listAllPolicies(store *iamstore.IAMStore, scope, pathPrefix string, onlyAttached bool) ([]*iamstore.Policy, error) {
	var all []*iamstore.Policy
	marker := ""
	for {
		result, err := store.Policies().List(scope, pathPrefix, onlyAttached, marker, 1000)
		if err != nil {
			return nil, err
		}
		all = append(all, result.Policies...)
		if !result.IsTruncated {
			break
		}
		marker = result.Marker
	}
	return all, nil
}

// cascadeDeleteUser removes all resources associated with a user before deleting
// the user record. Used by the admin handler to ensure consistent cleanup.
func cascadeDeleteUser(store *iamstore.IAMStore, userName string) error {
	if store.LoginProfiles().Exists(userName) {
		if err := store.LoginProfiles().Delete(userName); err != nil {
			return err
		}
	}

	if err := store.InlinePolicies().DeleteAllForPrincipal(PrincipalTypeUser, userName); err != nil {
		return err
	}

	attachedPolicies, err := store.AttachedPolicies().ListAttachedPolicies(PrincipalTypeUser, userName)
	if err != nil {
		return err
	}
	for _, policyArn := range attachedPolicies {
		if err := store.AttachedPolicies().Detach(PrincipalTypeUser, userName, policyArn); err != nil {
			return err
		}
		if err := store.Policies().DecrementAttachmentCount(policyArn); err != nil {
			return err
		}
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
	if err := store.InlinePolicies().DeleteAllForPrincipal(PrincipalTypeRole, roleName); err != nil {
		return err
	}

	attachedPolicies, err := store.AttachedPolicies().ListAttachedPolicies(PrincipalTypeRole, roleName)
	if err != nil {
		return err
	}
	for _, policyArn := range attachedPolicies {
		if err := store.AttachedPolicies().Detach(PrincipalTypeRole, roleName, policyArn); err != nil {
			return err
		}
		if err := store.Policies().DecrementAttachmentCount(policyArn); err != nil {
			return err
		}
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

	attachedPolicies, err := store.AttachedPolicies().ListAttachedPolicies(PrincipalTypeGroup, groupName)
	if err != nil {
		return err
	}
	for _, policyArn := range attachedPolicies {
		if err := store.AttachedPolicies().Detach(PrincipalTypeGroup, groupName, policyArn); err != nil {
			return err
		}
		if err := store.Policies().DecrementAttachmentCount(policyArn); err != nil {
			return err
		}
	}

	if err := store.UserGroups().RemoveAllUsersFromGroup(groupName); err != nil {
		return err
	}

	return store.Groups().Delete(groupName)
}

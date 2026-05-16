// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"
	"errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/aws/types"
	"vorpalstacks/internal/utils/timeutils"
)

// GetUser retrieves an IAM user by its name.
// UserName is required.
// Returns an error if the user does not exist.
func (s *IAMService) GetUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	if userName == "" {
		return nil, NewValidationError("UserName")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	user, err := store.Users().Get(userName)
	if err != nil {
		return nil, NewNoSuchUserError(userName)
	}

	return map[string]interface{}{
		"User": s.userToResponse(reqCtx, user),
	}, nil
}

// CreateUser creates a new IAM user.
// UserName is required and must not be empty.
// Path defaults to "/" if not specified.
// Tags are optional.
func (s *IAMService) CreateUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	if userName == "" {
		return nil, NewInvalidInputError("UserName", "cannot be empty")
	}

	path := request.GetStringParam(req.Parameters, "Path")
	if path == "" {
		path = "/"
	}

	newTags := tags.ParseTagsWithQueryFallback(req.Parameters, "Tags")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	user, err := store.Users().Create(userName, path, s.accountID, newTags)
	if err != nil {
		if errors.Is(err, iamstore.ErrUserAlreadyExists) {
			return nil, NewUserAlreadyExistsError(userName)
		}
		return nil, err
	}

	return map[string]interface{}{
		"User": s.userToResponse(reqCtx, user),
	}, nil
}

// DeleteUser deletes an IAM user by its name.
// UserName is required.
// Returns an error if the user has MFA devices, access keys, login profile, or attached policies.
// Also removes the user from all groups and deletes inline policies.
func (s *IAMService) DeleteUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	if userName == "" {
		return nil, NewValidationError("UserName")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Users().Exists(userName) {
		return nil, NewNoSuchUserError(userName)
	}

	mfaResult, err := store.MFADevices().ListForUser(userName, "", 1)
	if err != nil {
		return nil, err
	}
	if len(mfaResult.MFADevices) > 0 {
		return nil, NewDeleteConflictError("Cannot delete entity, must deactivate MFA devices first.")
	}

	keyCount, err := store.AccessKeys().CountByUserName(userName)
	if err != nil {
		return nil, err
	}
	if keyCount > 0 {
		return nil, NewDeleteConflictError("Cannot delete entity, must delete access keys first.")
	}

	if store.LoginProfiles().Exists(userName) {
		if err := store.LoginProfiles().Delete(userName); err != nil {
			return nil, err
		}
	}

	if err := store.InlinePolicies().DeleteAllForPrincipal(PrincipalTypeUser, userName); err != nil {
		return nil, err
	}

	attachedPolicies, err := store.AttachedPolicies().ListAttachedPolicies(PrincipalTypeUser, userName)
	if err != nil {
		return nil, err
	}
	for _, policyArn := range attachedPolicies {
		if err := store.AttachedPolicies().Detach(PrincipalTypeUser, userName, policyArn); err != nil {
			return nil, err
		}
		if err := store.Policies().DecrementAttachmentCount(policyArn); err != nil {
			return nil, err
		}
	}

	if err := store.UserGroups().RemoveAllGroupsForUser(userName); err != nil {
		return nil, err
	}

	if err := store.SigningCertificates().DeleteAllForUser(userName); err != nil {
		return nil, err
	}

	if err := store.SSHPublicKeys().DeleteAllForUser(userName); err != nil {
		return nil, err
	}

	if err := store.Users().Delete(userName); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UpdateUser updates the path and/or name of an IAM user.
// UserName is required.
// NewPath and NewUserName are optional parameters to update.
// If NewUserName is provided, migrates all user resources to the new name.
func (s *IAMService) UpdateUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	if userName == "" {
		return nil, NewValidationError("UserName")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	newPath := request.GetStringParam(req.Parameters, "NewPath")
	newUserName := request.GetStringParam(req.Parameters, "NewUserName")

	if err := store.RenameUser(userName, newUserName, newPath); err != nil {
		return nil, err
	}

	targetName := userName
	if newUserName != "" {
		targetName = newUserName
	}
	user, err := store.Users().Get(targetName)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"User": s.userToResponse(reqCtx, user),
	}, nil
}

// ListUsers lists IAM users.
// PathPrefix filters by path prefix.
// Supports pagination via Marker and MaxItems.
func (s *IAMService) ListUsers(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	pathPrefix := request.GetStringParam(req.Parameters, "PathPrefix")
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := store.Users().List(pathPrefix, marker, maxItems)
	if err != nil {
		return nil, err
	}

	users := make([]interface{}, len(result.Users))
	for i, user := range result.Users {
		users[i] = s.userToResponse(reqCtx, user)
	}

	response := map[string]interface{}{
		"Users":       users,
		"IsTruncated": result.IsTruncated,
	}

	if result.Marker != "" {
		response["Marker"] = result.Marker
	}

	return response, nil
}

var userTagOps = tagOps[*iamstore.User]{
	paramName:  "UserName",
	emptyErr:   NewValidationError("UserName"),
	notFoundFn: func(n string) error { return NewNoSuchUserError(n) },
	getFn:      func(s *iamstore.IAMStore, n string) (*iamstore.User, error) { return s.Users().Get(n) },
	putFn:      func(s *iamstore.IAMStore, r *iamstore.User) error { return s.Users().Put(r) },
	tagsFn:     func(r *iamstore.User) *[]types.Tag { return &r.Tags },
}

// TagUser adds tags to an IAM user.
// UserName is required.
// Tags are provided as a list of key-value pairs.
func (s *IAMService) TagUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagResource(ctx, s, reqCtx, req, userTagOps)
}

// UntagUser removes tags from an IAM user.
// UserName is required.
// TagKeys specifies which tags to remove.
func (s *IAMService) UntagUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return untagResource(ctx, s, reqCtx, req, userTagOps)
}

// ListUserTags lists the tags attached to an IAM user.
// UserName is required.
func (s *IAMService) ListUserTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return listResourceTags(ctx, s, reqCtx, req, userTagOps)
}

// PutUserPermissionsBoundary sets the permissions boundary for an IAM user.
// UserName is required.
// PermissionsBoundary is the ARN of a managed policy to use as the permissions boundary.
func (s *IAMService) PutUserPermissionsBoundary(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	if userName == "" {
		return nil, NewValidationError("UserName")
	}

	permissionsBoundary := request.GetStringParam(req.Parameters, "PermissionsBoundary")
	if permissionsBoundary == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	user, err := store.Users().Get(userName)
	if err != nil {
		return nil, NewNoSuchUserError(userName)
	}

	if !store.Policies().Exists(permissionsBoundary) {
		return nil, NewNoSuchPolicyError(permissionsBoundary)
	}

	user.PermissionsBoundary = &iamstore.PermissionsBoundary{
		PermissionsBoundaryType: "Policy",
		PermissionsBoundaryArn:  permissionsBoundary,
	}

	if err := store.Users().Put(user); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DeleteUserPermissionsBoundary removes the permissions boundary from an IAM user.
// UserName is required.
func (s *IAMService) DeleteUserPermissionsBoundary(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	userName := request.GetStringParam(req.Parameters, "UserName")
	if userName == "" {
		return nil, NewValidationError("UserName")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	user, err := store.Users().Get(userName)
	if err != nil {
		return nil, NewNoSuchUserError(userName)
	}

	user.PermissionsBoundary = nil

	if err := store.Users().Put(user); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// GetAccountAuthorizationDetails retrieves information about all IAM users, groups,
// roles, and policies in the account, including their relationships.
// Supports pagination via Marker and MaxItems.
func (s *IAMService) GetAccountAuthorizationDetails(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	filterParam := request.GetStringList(req.Parameters, "Filter")
	filters := make(map[string]bool)
	for _, f := range filterParam {
		filters[f] = true
	}
	if len(filters) == 0 {
		filters["User"] = true
		filters["Group"] = true
		filters["Role"] = true
		filters["LocalManagedPolicy"] = true
	}

	marker := pagination.GetMarker(req.Parameters)
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	type section struct {
		key    string
		filter string
		items  []interface{}
		marker func(i int) string
	}
	sections := make([]section, 0, 4)

	if filters["User"] {
		users, err := listAllUsers(store, "")
		if err != nil {
			return nil, err
		}
		items := make([]interface{}, 0, len(users))
		for _, user := range users {
			detail := map[string]interface{}{
				"UserId":     user.ID,
				"Path":       user.Path,
				"UserName":   user.UserName,
				"Arn":        user.Arn,
				"CreateDate": user.CreateDate.Format(timeutils.ISO8601SimpleFormat),
			}
			groupNames, _ := store.UserGroups().ListGroupsForUser(user.UserName)
			groupList := make([]interface{}, 0, len(groupNames))
			for _, gn := range groupNames {
				if g, err := store.Groups().Get(gn); err == nil {
					groupList = append(groupList, map[string]interface{}{
						"GroupName": g.GroupName,
						"GroupId":   g.ID,
						"Arn":       g.Arn,
					})
				}
			}
			detail["GroupList"] = groupList
			detail["AttachedManagedPolicies"] = buildAttachedManagedPolicies(store, PrincipalTypeUser, user.UserName)
			detail["UserPolicyList"] = buildInlinePolicyList(store, PrincipalTypeUser, user.UserName)
			items = append(items, detail)
		}
		sections = append(sections, section{
			key: "UserDetailList", filter: "User", items: items,
			marker: func(i int) string { return "user:" + users[i].UserName },
		})
	}

	if filters["Group"] {
		groups, err := listAllGroups(store, "")
		if err != nil {
			return nil, err
		}
		items := make([]interface{}, 0, len(groups))
		for _, group := range groups {
			detail := map[string]interface{}{
				"GroupId":    group.ID,
				"Path":       group.Path,
				"GroupName":  group.GroupName,
				"Arn":        group.Arn,
				"CreateDate": group.CreateDate.Format(timeutils.ISO8601SimpleFormat),
			}
			detail["GroupPolicyList"] = buildInlinePolicyList(store, PrincipalTypeGroup, group.GroupName)
			detail["AttachedManagedPolicies"] = buildAttachedManagedPolicies(store, PrincipalTypeGroup, group.GroupName)
			items = append(items, detail)
		}
		sections = append(sections, section{
			key: "GroupDetailList", filter: "Group", items: items,
			marker: func(i int) string { return "group:" + groups[i].GroupName },
		})
	}

	if filters["Role"] {
		roles, err := listAllRoles(store, "")
		if err != nil {
			return nil, err
		}
		items := make([]interface{}, 0, len(roles))
		for _, role := range roles {
			detail := map[string]interface{}{
				"RoleId":                   role.ID,
				"Path":                     role.Path,
				"RoleName":                 role.RoleName,
				"Arn":                      role.Arn,
				"CreateDate":               role.CreateDate.Format(timeutils.ISO8601SimpleFormat),
				"AssumeRolePolicyDocument": role.AssumeRolePolicyDocument,
			}
			detail["RolePolicyList"] = buildInlinePolicyList(store, PrincipalTypeRole, role.RoleName)
			detail["AttachedManagedPolicies"] = buildAttachedManagedPolicies(store, PrincipalTypeRole, role.RoleName)
			items = append(items, detail)
		}
		sections = append(sections, section{
			key: "RoleDetailList", filter: "Role", items: items,
			marker: func(i int) string { return "role:" + roles[i].RoleName },
		})
	}

	if filters["LocalManagedPolicy"] {
		policies, err := listAllPolicies(store, "Local", "", false)
		if err != nil {
			return nil, err
		}
		items := make([]interface{}, 0, len(policies))
		for _, policy := range policies {
			items = append(items, map[string]interface{}{
				"PolicyName":       policy.PolicyName,
				"PolicyId":         policy.ID,
				"Arn":              policy.Arn,
				"Path":             policy.Path,
				"DefaultVersionId": policy.DefaultVersionId,
			})
		}
		sections = append(sections, section{
			key: "Policies", filter: "LocalManagedPolicy", items: items,
			marker: func(i int) string { return "policy:" + policies[i].Arn },
		})
	}

	// Apply pagination across all sections combined.
	// Marker format: "<sectionType>:<itemName>" (e.g. "user:admin", "role:MyRole").
	resp := map[string]interface{}{}
	skipUntilMarker := marker != ""
	count := 0
	isTruncated := false
	nextMarker := ""

	for _, sec := range sections {
		// Always emit the key so AWS SDK deserialisation gets an empty list
		// instead of null for skipped/partial sections.
		secItems := make([]interface{}, 0)

		for i, item := range sec.items {
			itemMarker := sec.marker(i)

			if skipUntilMarker {
				if itemMarker == marker {
					skipUntilMarker = false
				}
				continue
			}

			if count >= maxItems {
				isTruncated = true
				nextMarker = itemMarker
				break
			}
			secItems = append(secItems, item)
			count++
		}

		resp[sec.key] = secItems

		if isTruncated {
			break
		}
	}

	if skipUntilMarker {
		for _, sec := range sections {
			resp[sec.key] = []interface{}{}
		}
	}

	resp["IsTruncated"] = isTruncated
	if isTruncated && nextMarker != "" {
		resp["Marker"] = nextMarker
	}

	return resp, nil
}

func (s *IAMService) userToResponse(reqCtx *request.RequestContext, user *iamstore.User) map[string]interface{} {
	resp := map[string]interface{}{
		"UserId":     user.ID,
		"Path":       user.Path,
		"UserName":   user.UserName,
		"Arn":        user.Arn,
		"CreateDate": user.CreateDate.Format(timeutils.ISO8601SimpleFormat),
	}

	if user.PasswordLastUsed != nil {
		resp["PasswordLastUsed"] = user.PasswordLastUsed.Format(timeutils.ISO8601SimpleFormat)
	}

	if user.PermissionsBoundary != nil {
		resp["PermissionsBoundary"] = map[string]interface{}{
			"PermissionsBoundaryType": user.PermissionsBoundary.PermissionsBoundaryType,
			"PermissionsBoundaryArn":  user.PermissionsBoundary.PermissionsBoundaryArn,
		}
	}

	if tags := tags.ToResponse(user.Tags); tags != nil {
		resp["Tags"] = tags
	}

	return resp
}

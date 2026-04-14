// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"
	"errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
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
	user, err := store.Users().Get(userName)
	if err != nil {
		return nil, NewNoSuchUserError(userName)
	}

	newPath := request.GetStringParam(req.Parameters, "NewPath")
	if newPath != "" {
		user.Path = newPath
		user.Arn = store.ARNBuilder().UserARN(newPath, user.UserName)
	}

	newUserName := request.GetStringParam(req.Parameters, "NewUserName")
	if newUserName != "" && newUserName != userName {
		if store.Users().Exists(newUserName) {
			return nil, NewUserAlreadyExistsError(newUserName)
		}

		user.UserName = newUserName
		user.Arn = store.ARNBuilder().UserARN(user.Path, newUserName)

		if err := store.Users().Put(user); err != nil {
			return nil, err
		}

		keys, err := store.AccessKeys().ListByUserNameWithSecret(userName)
		if err != nil {
			return nil, err
		}
		for _, key := range keys {
			key.UserName = newUserName
			if err := store.AccessKeys().Put(key); err != nil {
				return nil, err
			}
		}

		if store.LoginProfiles().Exists(userName) {
			profile, err := store.LoginProfiles().Get(userName)
			if err != nil {
				return nil, err
			}
			profile.UserName = newUserName
			if err := store.LoginProfiles().Put(profile); err != nil {
				return nil, err
			}
			if err := store.LoginProfiles().Delete(userName); err != nil {
				return nil, err
			}
		}

		if err := store.InlinePolicies().MigratePrincipal(userName, newUserName, PrincipalTypeUser); err != nil {
			return nil, err
		}

		if err := store.AttachedPolicies().MigratePrincipal(userName, newUserName, PrincipalTypeUser); err != nil {
			return nil, err
		}

		if err := store.UserGroups().MigrateUser(userName, newUserName); err != nil {
			return nil, err
		}

		if err := store.MFADevices().MigrateUser(userName, newUserName); err != nil {
			return nil, err
		}

		if err := store.Users().Delete(userName); err != nil {
			return nil, err
		}
	} else {
		if err := store.Users().Put(user); err != nil {
			return nil, err
		}
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
	maxItems := getMaxItems(req.Parameters)

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

// TagUser adds tags to an IAM user.
// UserName is required.
// Tags are provided as a list of key-value pairs.
func (s *IAMService) TagUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	user.Tags = tags.Apply(user.Tags, tags.ParseTagsWithQueryFallback(req.Parameters, "Tags"))

	if err := store.Users().Put(user); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UntagUser removes tags from an IAM user.
// UserName is required.
// TagKeys specifies which tags to remove.
func (s *IAMService) UntagUser(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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

	user.Tags = removeTags(user.Tags, parseTagKeys(req.Parameters))

	if err := store.Users().Put(user); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListUserTags lists the tags attached to an IAM user.
// UserName is required.
func (s *IAMService) ListUserTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
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
		"Tags":        tagsToResponse(user.Tags),
		"IsTruncated": false,
	}, nil
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

// GetAccountAuthorizationDetails retrieves information about all IAM users, groups, roles, and policies in the account, including their relationships.
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

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"IsTruncated": false,
	}

	if filters["User"] {
		userResult, err := store.Users().List("", "", 1000)
		if err != nil {
			return nil, err
		}
		userDetails := make([]interface{}, 0, len(userResult.Users))
		for _, user := range userResult.Users {
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

			attachedPolicyArns, _ := store.AttachedPolicies().ListAttachedPolicies(PrincipalTypeUser, user.UserName)
			attachedPolicies := make([]interface{}, 0, len(attachedPolicyArns))
			for _, arn := range attachedPolicyArns {
				if p, err := store.Policies().Get(arn); err == nil {
					attachedPolicies = append(attachedPolicies, map[string]interface{}{
						"PolicyName": p.PolicyName,
						"PolicyArn":  p.Arn,
					})
				}
			}
			detail["AttachedManagedPolicies"] = attachedPolicies

			inlinePolicyNames, _ := store.InlinePolicies().List(PrincipalTypeUser, user.UserName)
			userPolicyList := make([]interface{}, 0, len(inlinePolicyNames))
			for _, pn := range inlinePolicyNames {
				userPolicyList = append(userPolicyList, map[string]interface{}{
					"PolicyName": pn,
				})
			}
			detail["UserPolicyList"] = userPolicyList

			userDetails = append(userDetails, detail)
		}
		response["UserDetailList"] = userDetails
	}

	if filters["Group"] {
		groupResult, err := store.Groups().List("", "", 1000)
		if err != nil {
			return nil, err
		}
		groupDetails := make([]interface{}, 0, len(groupResult.Groups))
		for _, group := range groupResult.Groups {
			detail := map[string]interface{}{
				"GroupId":    group.ID,
				"Path":       group.Path,
				"GroupName":  group.GroupName,
				"Arn":        group.Arn,
				"CreateDate": group.CreateDate.Format(timeutils.ISO8601SimpleFormat),
			}

			inlinePolicyNames, _ := store.InlinePolicies().List(PrincipalTypeGroup, group.GroupName)
			groupPolicyList := make([]interface{}, 0, len(inlinePolicyNames))
			for _, pn := range inlinePolicyNames {
				groupPolicyList = append(groupPolicyList, map[string]interface{}{
					"PolicyName": pn,
				})
			}
			detail["GroupPolicyList"] = groupPolicyList

			attachedPolicyArns, _ := store.AttachedPolicies().ListAttachedPolicies(PrincipalTypeGroup, group.GroupName)
			attachedPolicies := make([]interface{}, 0, len(attachedPolicyArns))
			for _, arn := range attachedPolicyArns {
				if p, err := store.Policies().Get(arn); err == nil {
					attachedPolicies = append(attachedPolicies, map[string]interface{}{
						"PolicyName": p.PolicyName,
						"PolicyArn":  p.Arn,
					})
				}
			}
			detail["AttachedManagedPolicies"] = attachedPolicies

			groupDetails = append(groupDetails, detail)
		}
		response["GroupDetailList"] = groupDetails
	}

	if filters["Role"] {
		roleResult, err := store.Roles().List("", "", 1000)
		if err != nil {
			return nil, err
		}
		roleDetails := make([]interface{}, 0, len(roleResult.Roles))
		for _, role := range roleResult.Roles {
			detail := map[string]interface{}{
				"RoleId":                   role.ID,
				"Path":                     role.Path,
				"RoleName":                 role.RoleName,
				"Arn":                      role.Arn,
				"CreateDate":               role.CreateDate.Format(timeutils.ISO8601SimpleFormat),
				"AssumeRolePolicyDocument": role.AssumeRolePolicyDocument,
			}

			inlinePolicyNames, _ := store.InlinePolicies().List(PrincipalTypeRole, role.RoleName)
			rolePolicyList := make([]interface{}, 0, len(inlinePolicyNames))
			for _, pn := range inlinePolicyNames {
				rolePolicyList = append(rolePolicyList, map[string]interface{}{
					"PolicyName": pn,
				})
			}
			detail["RolePolicyList"] = rolePolicyList

			attachedPolicyArns, _ := store.AttachedPolicies().ListAttachedPolicies(PrincipalTypeRole, role.RoleName)
			attachedPolicies := make([]interface{}, 0, len(attachedPolicyArns))
			for _, arn := range attachedPolicyArns {
				if p, err := store.Policies().Get(arn); err == nil {
					attachedPolicies = append(attachedPolicies, map[string]interface{}{
						"PolicyName": p.PolicyName,
						"PolicyArn":  p.Arn,
					})
				}
			}
			detail["AttachedManagedPolicies"] = attachedPolicies

			roleDetails = append(roleDetails, detail)
		}
		response["RoleDetailList"] = roleDetails
	}

	if filters["LocalManagedPolicy"] {
		policyResult, err := store.Policies().List("Local", "", false, "", 1000)
		if err != nil {
			return nil, err
		}
		policies := make([]interface{}, 0, len(policyResult.Policies))
		for _, policy := range policyResult.Policies {
			policies = append(policies, map[string]interface{}{
				"PolicyName":       policy.PolicyName,
				"PolicyId":         policy.ID,
				"Arn":              policy.Arn,
				"Path":             policy.Path,
				"DefaultVersionId": policy.DefaultVersionId,
			})
		}
		response["Policies"] = policies
	}

	return response, nil
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

	if tags := tagsToResponse(user.Tags); tags != nil {
		resp["Tags"] = tags
	}

	return resp
}

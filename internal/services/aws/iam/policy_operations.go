// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"
	stderrors "errors"
	"fmt"
	"net/http"
	"strconv"

	"vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/timeutils"
)

// CreatePolicy creates a new managed policy.
// PolicyName is required and must not be empty.
// Path defaults to "/" if not specified.
// PolicyDocument must be a valid JSON policy document.
// Description and Tags are optional.
func (s *IAMService) CreatePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetStringParam(req.Parameters, "PolicyName")
	if policyName == "" {
		return nil, NewInvalidInputError("PolicyName", "cannot be empty")
	}

	path := request.GetStringParam(req.Parameters, "Path")
	if path == "" {
		path = "/"
	}

	document := request.GetStringParam(req.Parameters, "PolicyDocument")
	if !validatePolicyDocument(document) {
		return nil, ErrMalformedPolicyDocument
	}
	description := request.GetStringParam(req.Parameters, "Description")
	newTags := tags.ParseTagsWithQueryFallback(req.Parameters, "Tags")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	policy, err := store.Policies().Create(policyName, path, s.accountID, document, description, newTags)
	if err != nil {
		if stderrors.Is(err, iamstore.ErrPolicyAlreadyExists) {
			return nil, NewPolicyAlreadyExistsError(policyName)
		}
		return nil, err
	}

	return map[string]interface{}{
		"Policy": s.policyToResponse(reqCtx, policy),
	}, nil
}

// GetPolicy retrieves a managed policy by its ARN.
// Returns an error if the policy does not exist.
func (s *IAMService) GetPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyArn := request.GetStringParam(req.Parameters, "PolicyArn")
	if policyArn == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	policy, err := store.Policies().Get(policyArn)
	if err != nil {
		return nil, NewNoSuchPolicyError(policyArn)
	}

	return map[string]interface{}{
		"Policy": s.policyToResponse(reqCtx, policy),
	}, nil
}

// DeletePolicy deletes a managed policy by its ARN.
// Returns an error if the policy is attached to any users, groups, or roles.
// Returns an error if the policy does not exist.
func (s *IAMService) DeletePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyArn := request.GetStringParam(req.Parameters, "PolicyArn")
	if policyArn == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	policy, err := store.Policies().Get(policyArn)
	if err != nil {
		return nil, NewNoSuchPolicyError(policyArn)
	}

	if policy.AttachmentCount > 0 {
		return nil, NewDeletePolicyConflictError(policyArn)
	}

	if err := store.Policies().DeleteAllVersions(policyArn); err != nil {
		return nil, err
	}

	if err := store.Policies().Delete(policyArn); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListPolicies lists managed policies.
// Scope filters by policy scope (Local, AWS, All). Defaults to "Local".
// PathPrefix filters by path prefix.
// OnlyAttached filters to only attached policies.
// Supports pagination via Marker and MaxItems.
func (s *IAMService) ListPolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	scope := request.GetStringParam(req.Parameters, "Scope")
	if scope == "" {
		scope = "Local"
	}
	pathPrefix := request.GetStringParam(req.Parameters, "PathPrefix")
	onlyAttached := request.GetBoolParam(req.Parameters, "OnlyAttached")
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := store.Policies().List(scope, pathPrefix, onlyAttached, marker, maxItems)
	if err != nil {
		return nil, err
	}

	policies := make([]interface{}, len(result.Policies))
	for i, policy := range result.Policies {
		policies[i] = s.policyToResponse(reqCtx, policy)
	}

	response := map[string]interface{}{
		"Policies":    policies,
		"IsTruncated": result.IsTruncated,
	}

	if result.Marker != "" {
		response["Marker"] = result.Marker
	}

	return response, nil
}

// CreatePolicyVersion creates a new version of a managed policy.
// PolicyArn is required and must refer to an existing policy.
// PolicyDocument must be a valid JSON policy document.
// SetAsDefault specifies whether this version should be the default.
// Returns an error if the policy has reached the maximum number of versions.
func (s *IAMService) CreatePolicyVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyArn := request.GetStringParam(req.Parameters, "PolicyArn")
	if policyArn == "" {
		return nil, ErrNoSuchPolicy
	}

	document := request.GetStringParam(req.Parameters, "PolicyDocument")
	if !validatePolicyDocument(document) {
		return nil, ErrMalformedPolicyDocument
	}
	setAsDefault := request.GetBoolParam(req.Parameters, "SetAsDefault")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Policies().Exists(policyArn) {
		return nil, NewNoSuchPolicyError(policyArn)
	}

	versionCount, err := store.Policies().CountVersions(policyArn)
	if err != nil {
		return nil, err
	}
	if versionCount >= MaxPolicyVersions {
		return nil, ErrLimitExceededPolicyVersions
	}

	maxVersion, err := store.Policies().GetMaxVersion(policyArn)
	if err != nil {
		return nil, err
	}
	versionId := fmt.Sprintf("v%d", maxVersion+1)
	version := &iamstore.PolicyVersion{
		VersionId:        versionId,
		PolicyArn:        policyArn,
		IsDefaultVersion: setAsDefault,
		Document:         document,
	}

	if err := store.Policies().PutVersion(version); err != nil {
		return nil, err
	}

	if setAsDefault {
		if err := store.Policies().SetDefaultVersion(policyArn, versionId); err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"PolicyVersion": s.policyVersionToResponse(reqCtx, version),
	}, nil
}

// GetPolicyVersion retrieves a specific version of a managed policy.
// PolicyArn and VersionId are required.
// Returns an error if the policy or version does not exist.
func (s *IAMService) GetPolicyVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyArn := request.GetStringParam(req.Parameters, "PolicyArn")
	versionId := request.GetStringParam(req.Parameters, "VersionId")

	if policyArn == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	version, err := store.Policies().GetVersion(policyArn, versionId)
	if err != nil {
		return nil, NewNoSuchPolicyVersionError(versionId)
	}

	return map[string]interface{}{
		"PolicyVersion": s.policyVersionToResponse(reqCtx, version),
	}, nil
}

// DeletePolicyVersion deletes a specific version of a managed policy.
// PolicyArn and VersionId are required.
// Returns an error if attempting to delete the default version.
func (s *IAMService) DeletePolicyVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyArn := request.GetStringParam(req.Parameters, "PolicyArn")
	versionId := request.GetStringParam(req.Parameters, "VersionId")

	if policyArn == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	policy, err := store.Policies().Get(policyArn)
	if err != nil {
		return nil, NewNoSuchPolicyError(policyArn)
	}

	if policy.DefaultVersionId == versionId {
		return nil, errors.NewAWSError("DeleteConflict", "Cannot delete the default policy version.", http.StatusConflict)
	}

	if err := store.Policies().DeleteVersion(policyArn, versionId); err != nil {
		return nil, NewNoSuchPolicyVersionError(versionId)
	}

	return response.EmptyResponse(), nil
}

// ListPolicyVersions lists all versions of a managed policy.
// PolicyArn is required.
// Supports pagination via Marker and MaxItems.
func (s *IAMService) ListPolicyVersions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyArn := request.GetStringParam(req.Parameters, "PolicyArn")
	if policyArn == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Policies().Exists(policyArn) {
		return nil, NewNoSuchPolicyError(policyArn)
	}

	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	result, err := store.Policies().ListVersions(policyArn, marker, maxItems)
	if err != nil {
		return nil, err
	}

	versions := make([]interface{}, len(result.Versions))
	for i, version := range result.Versions {
		versions[i] = s.policyVersionToResponse(reqCtx, version)
	}

	response := map[string]interface{}{
		"Versions":    versions,
		"IsTruncated": result.IsTruncated,
	}

	if result.Marker != "" {
		response["Marker"] = result.Marker
	}

	return response, nil
}

// SetDefaultPolicyVersion sets a specific version of a managed policy as the default.
// PolicyArn and VersionId are required.
// Returns an error if the policy or version does not exist.
func (s *IAMService) SetDefaultPolicyVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyArn := request.GetStringParam(req.Parameters, "PolicyArn")
	versionId := request.GetStringParam(req.Parameters, "VersionId")

	if policyArn == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Policies().Exists(policyArn) {
		return nil, NewNoSuchPolicyError(policyArn)
	}

	if err := store.Policies().SetDefaultVersion(policyArn, versionId); err != nil {
		return nil, NewNoSuchPolicyVersionError(versionId)
	}

	return response.EmptyResponse(), nil
}

// TagPolicy adds tags to a managed policy.
// PolicyArn is required.
// Tags are provided as a list of key-value pairs.
func (s *IAMService) TagPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyArn := request.GetStringParam(req.Parameters, "PolicyArn")
	if policyArn == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	policy, err := store.Policies().Get(policyArn)
	if err != nil {
		return nil, NewNoSuchPolicyError(policyArn)
	}

	policy.Tags = tags.Apply(policy.Tags, tags.ParseTagsWithQueryFallback(req.Parameters, "Tags"))

	if err := store.Policies().Put(policy); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UntagPolicy removes tags from a managed policy.
// PolicyArn is required.
// TagKeys specifies which tags to remove.
func (s *IAMService) UntagPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyArn := request.GetStringParam(req.Parameters, "PolicyArn")
	if policyArn == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	policy, err := store.Policies().Get(policyArn)
	if err != nil {
		return nil, NewNoSuchPolicyError(policyArn)
	}

	policy.Tags = tags.RemoveByTagKeys(policy.Tags, tags.ParseTagKeysWithQueryFallback(req.Parameters, "TagKeys"))

	if err := store.Policies().Put(policy); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListPolicyTags lists the tags attached to a managed policy.
// PolicyArn is required.
func (s *IAMService) ListPolicyTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyArn := request.GetStringParam(req.Parameters, "PolicyArn")
	if policyArn == "" {
		return nil, ErrNoSuchPolicy
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	policy, err := store.Policies().Get(policyArn)
	if err != nil {
		return nil, NewNoSuchPolicyError(policyArn)
	}

	return map[string]interface{}{
		"Tags":        tags.ToResponse(policy.Tags),
		"IsTruncated": false,
	}, nil
}

// ListEntitiesForPolicy lists all IAM users, groups, and roles that the specified managed policy is attached to.
func (s *IAMService) ListEntitiesForPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyArn := request.GetStringParam(req.Parameters, "PolicyArn")
	if policyArn == "" {
		return nil, ErrNoSuchPolicy
	}

	entityFilter := request.GetStringParam(req.Parameters, "EntityFilter")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Policies().Exists(policyArn) {
		return nil, NewNoSuchPolicyError(policyArn)
	}

	refs, err := store.AttachedPolicies().ListPrincipalsForPolicy(policyArn)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"IsTruncated": false,
	}

	policyUsers := make([]interface{}, 0)
	policyGroups := make([]interface{}, 0)
	policyRoles := make([]interface{}, 0)

	for _, ref := range refs {
		switch ref.PrincipalType {
		case PrincipalTypeUser:
			if entityFilter != "" && entityFilter != "User" {
				continue
			}
			if user, err := store.Users().Get(ref.PrincipalName); err == nil {
				policyUsers = append(policyUsers, map[string]interface{}{
					"UserName": user.UserName,
					"UserId":   user.ID,
					"Arn":      user.Arn,
				})
			}
		case PrincipalTypeGroup:
			if entityFilter != "" && entityFilter != "Group" {
				continue
			}
			if group, err := store.Groups().Get(ref.PrincipalName); err == nil {
				policyGroups = append(policyGroups, map[string]interface{}{
					"GroupName": group.GroupName,
					"GroupId":   group.ID,
					"Arn":       group.Arn,
				})
			}
		case PrincipalTypeRole:
			if entityFilter != "" && entityFilter != "Role" {
				continue
			}
			if role, err := store.Roles().Get(ref.PrincipalName); err == nil {
				policyRoles = append(policyRoles, map[string]interface{}{
					"RoleName": role.RoleName,
					"RoleId":   role.ID,
					"Arn":      role.Arn,
				})
			}
		}
	}

	response["PolicyUsers"] = policyUsers
	response["PolicyGroups"] = policyGroups
	response["PolicyRoles"] = policyRoles

	return response, nil
}

func (s *IAMService) policyToResponse(reqCtx *request.RequestContext, policy *iamstore.Policy) map[string]interface{} {
	resp := map[string]interface{}{
		"PolicyId":         policy.ID,
		"Path":             policy.Path,
		"PolicyName":       policy.PolicyName,
		"Arn":              policy.Arn,
		"CreateDate":       policy.CreateDate.Format(timeutils.ISO8601SimpleFormat),
		"UpdateDate":       policy.UpdateDate.Format(timeutils.ISO8601SimpleFormat),
		"DefaultVersionId": policy.DefaultVersionId,
		"AttachmentCount":  policy.AttachmentCount,
		"IsAttachable":     policy.IsAttachable,
	}
	if policy.Description != "" {
		resp["Description"] = policy.Description
	}
	if tags := tags.ToResponse(policy.Tags); tags != nil {
		resp["Tags"] = tags
	}
	return resp
}

func (s *IAMService) policyVersionToResponse(reqCtx *request.RequestContext, version *iamstore.PolicyVersion) map[string]interface{} {
	return map[string]interface{}{
		"VersionId":        version.VersionId,
		"IsDefaultVersion": version.IsDefaultVersion,
		"CreateDate":       version.CreateDate.Format(timeutils.ISO8601SimpleFormat),
		"Document":         version.Document,
	}
}

// SimulatePrincipalPolicy simulates the effects of a set of IAM policies on a principal.
func (s *IAMService) SimulatePrincipalPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	var actionNames []string
	for i := 1; ; i++ {
		action := request.GetStringParam(req.Parameters, "ActionNames.member."+strconv.Itoa(i))
		if action == "" {
			break
		}
		actionNames = append(actionNames, action)
	}

	if len(actionNames) == 0 {
		if raw, ok := req.Parameters["ActionNames"].([]interface{}); ok {
			for _, a := range raw {
				if s, ok := a.(string); ok {
					actionNames = append(actionNames, s)
				}
			}
		}
	}

	resourceArns := []string{}
	for i := 1; ; i++ {
		arn := request.GetStringParam(req.Parameters, "ResourceArns.member."+strconv.Itoa(i))
		if arn == "" {
			break
		}
		resourceArns = append(resourceArns, arn)
	}

	evaluationResults := make([]interface{}, 0, len(actionNames))
	resources := resourceArns
	if len(resources) == 0 {
		resources = []string{"*"}
	}

	for _, action := range actionNames {
		for _, resource := range resources {
			evaluationResults = append(evaluationResults, map[string]interface{}{
				"EvalActionName":     action,
				"EvalResourceName":   resource,
				"EvalDecision":       "allowed",
				"MatchedStatements":  []interface{}{},
				"MissingContextKeys": []interface{}{},
				"OrganizationsDecisionDetail": map[string]interface{}{
					"AllowedByOrganizations": false,
				},
			})
		}
	}

	return map[string]interface{}{
		"EvaluationResults": evaluationResults,
		"IsTruncated":       false,
		"Marker":            "",
	}, nil
}

// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"
	stderrors "errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/iam/policy"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
	"vorpalstacks/internal/utils/aws/types"
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
	if !entityNamePattern128.MatchString(policyName) {
		return nil, NewInvalidInputError("PolicyName", "must be 1 to 128 alphanumeric characters or any of +=,.@-_")
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
	if err := validateNewTags(newTags); err != nil {
		return nil, err
	}

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

var policyTagOps = tagOps[*iamstore.Policy]{
	paramName:  "PolicyArn",
	emptyErr:   ErrNoSuchPolicy,
	notFoundFn: func(n string) error { return NewNoSuchPolicyError(n) },
	getFn:      func(s *iamstore.IAMStore, n string) (*iamstore.Policy, error) { return s.Policies().Get(n) },
	putFn:      func(s *iamstore.IAMStore, r *iamstore.Policy) error { return s.Policies().Put(r) },
	tagsFn:     func(r *iamstore.Policy) *[]types.Tag { return &r.Tags },
}

// TagPolicy adds tags to a managed policy.
// PolicyArn is required.
// Tags are provided as a list of key-value pairs.
func (s *IAMService) TagPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return tagResource(ctx, s, reqCtx, req, policyTagOps)
}

// UntagPolicy removes tags from a managed policy.
// PolicyArn is required.
// TagKeys specifies which tags to remove.
func (s *IAMService) UntagPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return untagResource(ctx, s, reqCtx, req, policyTagOps)
}

// ListPolicyTags lists the tags attached to a managed policy.
// PolicyArn is required.
func (s *IAMService) ListPolicyTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return listResourceTags(ctx, s, reqCtx, req, policyTagOps)
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

	type entityEntry struct {
		entityType string
		name       string
		data       map[string]interface{}
	}

	combined := make([]entityEntry, 0)

	for _, ref := range refs {
		switch ref.PrincipalType {
		case PrincipalTypeUser:
			if entityFilter != "" && entityFilter != "User" {
				continue
			}
			if user, err := store.Users().Get(ref.PrincipalName); err == nil {
				entry := map[string]interface{}{
					"UserName": user.UserName,
					"UserId":   user.ID,
					"Arn":      user.Arn,
				}
				combined = append(combined, entityEntry{"User", user.UserName, entry})
			}
		case PrincipalTypeGroup:
			if entityFilter != "" && entityFilter != "Group" {
				continue
			}
			if group, err := store.Groups().Get(ref.PrincipalName); err == nil {
				entry := map[string]interface{}{
					"GroupName": group.GroupName,
					"GroupId":   group.ID,
					"Arn":       group.Arn,
				}
				combined = append(combined, entityEntry{"Group", group.GroupName, entry})
			}
		case PrincipalTypeRole:
			if entityFilter != "" && entityFilter != "Role" {
				continue
			}
			if role, err := store.Roles().Get(ref.PrincipalName); err == nil {
				entry := map[string]interface{}{
					"RoleName": role.RoleName,
					"RoleId":   role.ID,
					"Arn":      role.Arn,
				}
				combined = append(combined, entityEntry{"Role", role.RoleName, entry})
			}
		}
	}

	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	paged := pagination.PaginateSlice(combined, marker, maxItems, func(item entityEntry) string {
		return item.entityType + ":" + item.name
	})

	for _, entry := range paged.Items {
		switch entry.entityType {
		case "User":
			policyUsers = append(policyUsers, entry.data)
		case "Group":
			policyGroups = append(policyGroups, entry.data)
		case "Role":
			policyRoles = append(policyRoles, entry.data)
		}
	}

	response["PolicyUsers"] = policyUsers
	response["PolicyGroups"] = policyGroups
	response["PolicyRoles"] = policyRoles
	response["IsTruncated"] = paged.IsTruncated
	if paged.NextMarker != "" {
		response["Marker"] = paged.NextMarker
	}

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

	if store, err := s.store(reqCtx); err == nil {
		count := 0
		userMarker := ""
		for {
			users, err := store.Users().List("", userMarker, 1000)
			if err != nil {
				break
			}
			for _, u := range users.Users {
				if u.PermissionsBoundary != nil && u.PermissionsBoundary.PermissionsBoundaryArn == policy.Arn {
					count++
				}
			}
			if !users.IsTruncated || users.Marker == "" {
				break
			}
			userMarker = users.Marker
		}
		roleMarker := ""
		for {
			roles, err := store.Roles().List("", roleMarker, 1000)
			if err != nil {
				break
			}
			for _, r := range roles.Roles {
				if r.PermissionsBoundary != nil && r.PermissionsBoundary.PermissionsBoundaryArn == policy.Arn {
					count++
				}
			}
			if !roles.IsTruncated || roles.Marker == "" {
				break
			}
			roleMarker = roles.Marker
		}
		resp["PermissionsBoundaryUsageCount"] = count
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

// SimulatePrincipalPolicy simulates the effects of IAM policies on a principal.
func (s *IAMService) SimulatePrincipalPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policySourceArn := request.GetStringParam(req.Parameters, "PolicySourceArn")
	if policySourceArn == "" {
		return nil, NewValidationError("PolicySourceArn")
	}

	actionNames := request.GetStringList(req.Parameters, "ActionNames")
	if len(actionNames) == 0 {
		return nil, NewValidationError("ActionNames")
	}

	resourceArns := request.GetStringList(req.Parameters, "ResourceArns")
	resources := resourceArns
	if len(resources) == 0 {
		resources = []string{"*"}
	}

	// Gather all applicable policies for the principal.
	policyDocs, err := s.gatherPrincipalPolicies(reqCtx, policySourceArn)
	if err != nil {
		return nil, err
	}

	// Add PolicyInputList (additional inline policy documents).
	policyInputList := request.GetStringList(req.Parameters, "PolicyInputList")
	for _, pDoc := range policyInputList {
		doc, pErr := policy.ParseDocument(pDoc)
		if pErr != nil {
			return nil, NewInvalidInputError("PolicyInputList", "contains a malformed policy document")
		}
		policyDocs = append(policyDocs, doc)
	}

	// Permissions boundary — if present, it limits the maximum permissions.
	boundaryInputList := request.GetStringList(req.Parameters, "PermissionsBoundaryPolicyInputList")
	var boundaryDocs []*policy.Document
	for _, bDoc := range boundaryInputList {
		doc, bErr := policy.ParseDocument(bDoc)
		if bErr != nil {
			return nil, NewInvalidInputError("PermissionsBoundaryPolicyInputList", "contains a malformed policy document")
		}
		boundaryDocs = append(boundaryDocs, doc)
	}

	// Build context entries from request.
	sessionCtx := buildSimulationContextEntries(req.Parameters)
	principalName := extractPrincipalNameFromARN(policySourceArn)
	principalAccount := extractAccountFromARN(policySourceArn)

	evaluator := policy.NewPolicyEvaluator()

	evaluationResults := make([]interface{}, 0, len(actionNames)*len(resources))
	for _, action := range actionNames {
		for _, resource := range resources {
			evalCtx := &policy.EvaluationContext{
				Principal:        policySourceArn,
				PrincipalAccount: principalAccount,
				Action:           action,
				Resource:         resource,
				RequestTime:      time.Now(),
				UserName:         principalName,
				SessionContext:   sessionCtx,
			}

			// Evaluate identity-based policies.
			decision := evaluator.Evaluate(evalCtx, policyDocs)

			effect := decision.Effect
			matchedStatements := []interface{}{}
			if decision.MatchedSid != "" {
				matchedStatements = append(matchedStatements, map[string]interface{}{
					"SourcePolicyId": policySourceArn,
					"StatementIds":   []interface{}{decision.MatchedSid},
				})
			}

			// If allowed by identity policies, check permissions boundary.
			allowedByBoundary := true
			if effect == policy.DecisionEffectAllow && len(boundaryDocs) > 0 {
				boundaryDecision := evaluator.Evaluate(evalCtx, boundaryDocs)
				if boundaryDecision.Effect != policy.DecisionEffectAllow {
					effect = policy.DecisionEffectDefaultDeny
					matchedStatements = []interface{}{}
					allowedByBoundary = false
				}
			}

			evalDecision := "implicitDeny"
			switch effect {
			case policy.DecisionEffectAllow:
				evalDecision = "allowed"
			case policy.DecisionEffectDeny:
				evalDecision = "explicitDeny"
			}

			resultEntry := map[string]interface{}{
				"EvalActionName":       action,
				"EvalResourceName":     resource,
				"EvalDecision":         evalDecision,
				"MatchedStatements":    matchedStatements,
				"MissingContextValues": []interface{}{},
				"OrganizationsDecisionDetail": map[string]interface{}{
					"AllowedByOrganizations": false,
				},
			}
			if len(boundaryDocs) > 0 {
				resultEntry["PermissionsBoundaryDecisionDetail"] = map[string]interface{}{
					"AllowedByPermissionsBoundary": allowedByBoundary,
				}
			}
			evaluationResults = append(evaluationResults, resultEntry)
		}
	}

	return map[string]interface{}{
		"EvaluationResults": evaluationResults,
		"IsTruncated":       false,
		"Marker":            "",
	}, nil
}

// gatherPrincipalPolicies collects all identity-based policies applicable to the given principal.
func (s *IAMService) gatherPrincipalPolicies(reqCtx *request.RequestContext, principalArn string) ([]*policy.Document, error) {
	entityType := resolveEntityType(principalArn)
	entityName := resolveEntityName(principalArn)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	var docs []*policy.Document

	switch entityType {
	case "User":
		if !store.Users().Exists(entityName) {
			return nil, NewNoSuchUserError(entityName)
		}
		// User inline policies.
		docs = append(docs, collectInlinePolicies(store, PrincipalTypeUser, entityName)...)
		// User attached managed policies.
		docs = append(docs, collectAttachedPolicies(store, PrincipalTypeUser, entityName)...)
		// Group inline + attached policies for each group the user belongs to.
		groupNames, err := store.UserGroups().ListGroupsForUser(entityName)
		if err != nil {
			return nil, err
		}
		for _, groupName := range groupNames {
			docs = append(docs, collectInlinePolicies(store, PrincipalTypeGroup, groupName)...)
			docs = append(docs, collectAttachedPolicies(store, PrincipalTypeGroup, groupName)...)
		}

	case "Role":
		if !store.Roles().Exists(entityName) {
			return nil, NewNoSuchRoleError(entityName)
		}
		docs = append(docs, collectInlinePolicies(store, PrincipalTypeRole, entityName)...)
		docs = append(docs, collectAttachedPolicies(store, PrincipalTypeRole, entityName)...)

	case "Group":
		if !store.Groups().Exists(entityName) {
			return nil, NewNoSuchGroupError(entityName)
		}
		docs = append(docs, collectInlinePolicies(store, PrincipalTypeGroup, entityName)...)
		docs = append(docs, collectAttachedPolicies(store, PrincipalTypeGroup, entityName)...)
	}

	return docs, nil
}

func collectInlinePolicies(store *iamstore.IAMStore, principalType, principalName string) []*policy.Document {
	policyNames, err := store.InlinePolicies().List(principalType, principalName)
	if err != nil {
		return nil
	}
	var docs []*policy.Document
	for _, pn := range policyNames {
		ip, err := store.InlinePolicies().Get(principalType, principalName, pn)
		if err != nil || ip == nil {
			continue
		}
		doc, err := policy.ParseDocument(ip.PolicyDocument)
		if err != nil {
			continue
		}
		docs = append(docs, doc)
	}
	return docs
}

func collectAttachedPolicies(store *iamstore.IAMStore, principalType, principalName string) []*policy.Document {
	arns, err := store.AttachedPolicies().ListAttachedPolicies(principalType, principalName)
	if err != nil {
		return nil
	}
	var docs []*policy.Document
	for _, arn := range arns {
		version, err := store.Policies().GetDefaultVersion(arn)
		if err != nil || version == nil {
			continue
		}
		doc, err := policy.ParseDocument(version.Document)
		if err != nil {
			continue
		}
		docs = append(docs, doc)
	}
	return docs
}

func extractPrincipalNameFromARN(arn string) string {
	_, name := parseIAMARNResource(arn)
	return name
}

func extractAccountFromARN(arn string) string {
	// arn:aws:iam::123456789012:user/Bob
	parts := strings.Split(arn, ":")
	if len(parts) >= 5 {
		return parts[4]
	}
	return ""
}

func buildSimulationContextEntries(params map[string]interface{}) map[string]string {
	result := make(map[string]string)
	i := 1
	for {
		nameKey := fmt.Sprintf("ContextEntries.member.%d.ContextKeyName", i)
		nameVal, ok := params[nameKey]
		if !ok {
			nameKey2 := fmt.Sprintf("ContextEntries.%d.ContextKeyName", i)
			nameVal, ok = params[nameKey2]
			if !ok {
				break
			}
		}
		ctxKey, _ := nameVal.(string)
		if ctxKey == "" {
			break
		}
		valKey := fmt.Sprintf("ContextEntries.member.%d.ContextKeyValues.member.1", i)
		valKeyAlt := fmt.Sprintf("ContextEntries.%d.ContextKeyValues.1", i)
		if v, ok := params[valKey]; ok {
			if s, ok := v.(string); ok {
				result[ctxKey] = s
			}
		} else if v, ok := params[valKeyAlt]; ok {
			if s, ok := v.(string); ok {
				result[ctxKey] = s
			}
		}
		i++
	}
	return result
}

// Package iam provides IAM service operations for vorpalstacks.
package iam

import (
	"context"
	"fmt"
	"time"

	"vorpalstacks/internal/common/iam/policy"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/common/tags"
	iamstore "vorpalstacks/internal/store/aws/iam"
	awsarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/timeutils"
)

// CreatePolicy creates a new managed policy.
// PolicyName is required and must not be empty.
// Path defaults to "/" if not specified.
// PolicyDocument must be a valid JSON policy document.
// Description and Tags are optional.
func (s *IAMService) CreatePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &CreatePolicyInput{
		PolicyName:     request.GetStringParam(req.Parameters, "PolicyName"),
		Path:           request.GetStringParam(req.Parameters, "Path"),
		PolicyDocument: request.GetStringParam(req.Parameters, "PolicyDocument"),
		Description:    request.GetStringParam(req.Parameters, "Description"),
		Tags:           tags.ParseTagsWithQueryFallback(req.Parameters, "Tags"),
	}
	policy, err := s.createPolicyCore(store, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Policy": s.policyToResponse(policy),
	}, nil
}

// GetPolicy retrieves a managed policy by its ARN.
// Returns an error if the policy does not exist.
func (s *IAMService) GetPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyArn := request.GetStringParam(req.Parameters, "PolicyArn")
	if policyArn == "" {
		return nil, NewValidationError("PolicyArn")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	policy, err := s.getPolicyCore(store, policyArn)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Policy": s.policyToResponse(policy),
	}, nil
}

// DeletePolicy deletes a managed policy by its ARN.
// Returns an error if the policy is attached to any users, groups, or roles.
// Returns an error if the policy does not exist.
func (s *IAMService) DeletePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &DeletePolicyInput{
		PolicyArn: request.GetStringParam(req.Parameters, "PolicyArn"),
	}
	if err := s.deletePolicyCore(store, input); err != nil {
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
	if !validatePolicyScope(scope) {
		return nil, NewInvalidInputError("Scope", "must be one of: All, AWS, Local")
	}
	pathPrefix := request.GetStringParam(req.Parameters, "PathPrefix")
	onlyAttached := request.GetBoolParam(req.Parameters, "OnlyAttached")
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.listPoliciesCore(store, scope, pathPrefix, marker, onlyAttached, maxItems)
	if err != nil {
		return nil, err
	}

	policies := make([]interface{}, len(result.Policies))
	for i, policy := range result.Policies {
		policies[i] = s.policyToResponse(policy)
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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &CreatePolicyVersionInput{
		PolicyArn:      request.GetStringParam(req.Parameters, "PolicyArn"),
		PolicyDocument: request.GetStringParam(req.Parameters, "PolicyDocument"),
		SetAsDefault:   request.GetBoolParam(req.Parameters, "SetAsDefault"),
	}
	version, err := s.createPolicyVersionCore(store, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"PolicyVersion": s.policyVersionToResponse(version),
	}, nil
}

// GetPolicyVersion retrieves a specific version of a managed policy.
// PolicyArn and VersionId are required.
// Returns an error if the policy or version does not exist.
func (s *IAMService) GetPolicyVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyArn := request.GetStringParam(req.Parameters, "PolicyArn")
	versionId := request.GetStringParam(req.Parameters, "VersionId")

	if policyArn == "" {
		return nil, NewValidationError("PolicyArn")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	version, err := s.getPolicyVersionCore(store, policyArn, versionId)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"PolicyVersion": s.policyVersionToResponse(version),
	}, nil
}

// DeletePolicyVersion deletes a specific version of a managed policy.
// PolicyArn and VersionId are required.
// Returns an error if attempting to delete the default version.
func (s *IAMService) DeletePolicyVersion(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyArn := request.GetStringParam(req.Parameters, "PolicyArn")
	versionId := request.GetStringParam(req.Parameters, "VersionId")

	if policyArn == "" {
		return nil, NewValidationError("PolicyArn")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.deletePolicyVersionCore(store, policyArn, versionId); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListPolicyVersions lists all versions of a managed policy.
// PolicyArn is required.
// Supports pagination via Marker and MaxItems.
func (s *IAMService) ListPolicyVersions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyArn := request.GetStringParam(req.Parameters, "PolicyArn")
	if policyArn == "" {
		return nil, NewValidationError("PolicyArn")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	result, err := s.listPolicyVersionsCore(store, policyArn, marker, maxItems)
	if err != nil {
		return nil, err
	}

	versions := make([]interface{}, len(result.Versions))
	for i, version := range result.Versions {
		versions[i] = s.policyVersionToResponse(version)
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
		return nil, NewValidationError("PolicyArn")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.setDefaultPolicyVersionCore(store, policyArn, versionId); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

var policyTagOps = tagOps[*iamstore.Policy]{
	paramName:  "PolicyArn",
	emptyErr:   ErrNoSuchPolicy,
	notFoundFn: func(n string) error { return NewNoSuchPolicyError(n) },
	getFn:      func(s *iamstore.IAMStore, n string) (*iamstore.Policy, error) { return s.Policies().Get(n) },
	putFn:      func(s *iamstore.IAMStore, r *iamstore.Policy) error { return s.Policies().Put(r) },
	tagsFn:     func(r *iamstore.Policy) *[]tags.Tag { return &r.Tags },
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
		return nil, NewValidationError("PolicyArn")
	}

	entityFilter := request.GetStringParam(req.Parameters, "EntityFilter")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	result, err := s.listEntitiesForPolicyCore(store, policyArn, entityFilter, marker, maxItems)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"PolicyUsers":  result.PolicyUsers,
		"PolicyGroups": result.PolicyGroups,
		"PolicyRoles":  result.PolicyRoles,
		"IsTruncated":  result.IsTruncated,
	}
	if result.Marker != "" {
		response["Marker"] = result.Marker
	}

	return response, nil
}

func (s *IAMService) policyToResponse(policy *iamstore.Policy) map[string]interface{} {
	resp := map[string]interface{}{
		"PolicyId":                      policy.ID,
		"Path":                          policy.Path,
		"PolicyName":                    policy.PolicyName,
		"Arn":                           policy.Arn,
		"CreateDate":                    policy.CreateDate.Format(timeutils.ISO8601SimpleFormat),
		"UpdateDate":                    policy.UpdateDate.Format(timeutils.ISO8601SimpleFormat),
		"DefaultVersionId":              policy.DefaultVersionId,
		"AttachmentCount":               policy.AttachmentCount,
		"PermissionsBoundaryUsageCount": policy.PermissionsBoundaryUsageCount,
		"IsAttachable":                  policy.IsAttachable,
	}

	if policy.Description != "" {
		resp["Description"] = policy.Description
	}
	if tags := tags.ToResponse(policy.Tags); tags != nil {
		resp["Tags"] = tags
	}
	return resp
}

func (s *IAMService) policyVersionToResponse(version *iamstore.PolicyVersion) map[string]interface{} {
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
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	policyDocs, err := s.gatherPrincipalPoliciesCore(store, policySourceArn)
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

func extractPrincipalNameFromARN(arn string) string {
	_, name := parseIAMARNResource(arn)
	return name
}

func extractAccountFromARN(arn string) string {
	_, _, _, accountID, _ := awsarn.SplitARN(arn)
	return accountID
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

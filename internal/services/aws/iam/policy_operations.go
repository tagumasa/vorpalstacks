package iam

import (
	"context"

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
// Scope filters by policy scope (All, AWS, Local); an omitted Scope
// defaults to All.
// PathPrefix filters by path prefix.
// OnlyAttached filters to only attached policies.
// Supports pagination via Marker and MaxItems.
func (s *IAMService) ListPolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	scope := request.GetStringParam(req.Parameters, "Scope")
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

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := s.setDefaultPolicyVersionCore(store, policyArn, versionId); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// TagPolicy adds tags to a managed policy.
// PolicyArn is required.
// Tags are provided as a list of key-value pairs.
func (s *IAMService) TagPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &TagResourceInput{
		ResourceName: request.GetStringParam(req.Parameters, "PolicyArn"),
		Tags:         tags.ParseTagsWithQueryFallback(req.Parameters, "Tags"),
	}
	if err := tagResourceCore(store, policyTagOps, input); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// UntagPolicy removes tags from a managed policy.
// PolicyArn is required.
// TagKeys specifies which tags to remove.
func (s *IAMService) UntagPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &UntagResourceInput{
		ResourceName: request.GetStringParam(req.Parameters, "PolicyArn"),
		TagKeys:      tags.ParseTagKeysWithQueryFallback(req.Parameters, "TagKeys"),
	}
	if err := untagResourceCore(store, policyTagOps, input); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

// ListPolicyTags lists the tags attached to a managed policy.
// PolicyArn is required.
func (s *IAMService) ListPolicyTags(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &ListResourceTagsInput{
		ResourceName: request.GetStringParam(req.Parameters, "PolicyArn"),
		Marker:       request.GetStringParam(req.Parameters, "Marker"),
		MaxItems:     pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems),
	}
	result, err := listResourceTagsCore(store, policyTagOps, input)
	if err != nil {
		return nil, err
	}
	resp := map[string]interface{}{
		"Tags":        tags.ToResponse(result.Tags),
		"IsTruncated": result.IsTruncated,
	}
	if result.Marker != "" {
		resp["Marker"] = result.Marker
	}
	return resp, nil
}

// ListEntitiesForPolicy lists all IAM users, groups, and roles that the specified managed policy is attached to.
func (s *IAMService) ListEntitiesForPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyArn := request.GetStringParam(req.Parameters, "PolicyArn")
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
	input := &SimulatePrincipalPolicyInput{
		PolicySourceArn:                    request.GetStringParam(req.Parameters, "PolicySourceArn"),
		ActionNames:                        request.GetStringList(req.Parameters, "ActionNames"),
		ResourceArns:                       request.GetStringList(req.Parameters, "ResourceArns"),
		PolicyInputList:                    request.GetStringList(req.Parameters, "PolicyInputList"),
		PermissionsBoundaryPolicyInputList: request.GetStringList(req.Parameters, "PermissionsBoundaryPolicyInputList"),
		ContextEntries:                     buildSimulationContextEntries(req.Parameters),
		PolicyExclusionList:                buildSimulationPolicyIdentifiers(req.Parameters),
		ResourcePolicy:                     request.GetStringParam(req.Parameters, "ResourcePolicy"),
		ResourceOwner:                      request.GetStringParam(req.Parameters, "ResourceOwner"),
		CallerArn:                          request.GetStringParam(req.Parameters, "CallerArn"),
		ResourceHandlingOption:             request.GetStringParam(req.Parameters, "ResourceHandlingOption"),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.simulatePrincipalPolicyCore(store, input)
	if err != nil {
		return nil, err
	}
	return simulationResponse(result.Evaluations, req.Parameters), nil
}

// SimulateCustomPolicy simulates the effects of caller-supplied policy
// documents, without gathering any principal's policies.
func (s *IAMService) SimulateCustomPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := &SimulateCustomPolicyInput{
		PolicyInputList:                    request.GetStringList(req.Parameters, "PolicyInputList"),
		PermissionsBoundaryPolicyInputList: request.GetStringList(req.Parameters, "PermissionsBoundaryPolicyInputList"),
		OrderedOrganizationPolicyInputList: buildOrderedOrganizationPolicies(req.Parameters),
		ActionNames:                        request.GetStringList(req.Parameters, "ActionNames"),
		ResourceArns:                       request.GetStringList(req.Parameters, "ResourceArns"),
		ResourcePolicy:                     request.GetStringParam(req.Parameters, "ResourcePolicy"),
		ResourceOwner:                      request.GetStringParam(req.Parameters, "ResourceOwner"),
		CallerArn:                          request.GetStringParam(req.Parameters, "CallerArn"),
		ContextEntries:                     buildSimulationContextEntries(req.Parameters),
		ResourceHandlingOption:             request.GetStringParam(req.Parameters, "ResourceHandlingOption"),
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := s.simulateCustomPolicyCore(store, input)
	if err != nil {
		return nil, err
	}
	return simulationResponse(result.Evaluations, req.Parameters), nil
}

// buildOrderedOrganizationPolicies reads the ordered organisation policy
// levels: each list member is one level of the Organizations hierarchy
// carrying its own list of SCP documents. Every document of every level
// bounds the evaluation, so the levels flatten to their documents in wire
// order.
func buildOrderedOrganizationPolicies(params map[string]interface{}) []string {
	var documents []string
	for _, level := range request.GetListParam(params, "OrderedOrganizationPolicyInputList") {
		documents = append(documents, request.GetStringList(level, "ServiceControlPolicyInputList")...)
	}
	return documents
}

// buildSimulationPolicyIdentifiers parses the flattened
// PolicyExclusionList into PolicyIdentifier union entries; each member
// names exactly one of PolicyType, PolicyArn, or the inline-policy
// identifier's fields.
func buildSimulationPolicyIdentifiers(params map[string]interface{}) []SimulationPolicyIdentifier {
	var identifiers []SimulationPolicyIdentifier
	for _, entry := range request.GetListParam(params, "PolicyExclusionList") {
		id := SimulationPolicyIdentifier{
			PolicyType:           request.GetStringParam(entry, "PolicyType"),
			PolicyArn:            request.GetStringParam(entry, "PolicyArn"),
			InlinePolicyName:     request.GetStringParam(entry, "InlinePolicyIdentifier.PolicyName"),
			InlineAttachmentType: request.GetStringParam(entry, "InlinePolicyIdentifier.AttachmentType"),
			InlineAttachmentName: request.GetStringParam(entry, "InlinePolicyIdentifier.AttachmentName"),
		}
		if id.PolicyType == "" && id.PolicyArn == "" && id.InlinePolicyName == "" {
			continue
		}
		identifiers = append(identifiers, id)
	}
	return identifiers
}

// simulationDecisionDetailsToResponse renders the per-policy-type
// decisions as the EvalDecisionDetails map.
func simulationDecisionDetailsToResponse(details []simulationDecisionDetail) map[string]string {
	if details == nil {
		return nil
	}
	out := make(map[string]string, len(details))
	for _, d := range details {
		out[d.PolicyType] = d.Decision
	}
	return out
}

// simulationStatementsToResponse renders matched statements in the wire
// Statement shape: the source policy's identifier and type. A policy
// without a type (an input-list document) omits the SourcePolicyType
// member.
func simulationStatementsToResponse(matched []simulationSourcePolicy) []interface{} {
	statements := make([]interface{}, 0, len(matched))
	for _, source := range matched {
		entry := map[string]interface{}{
			"SourcePolicyId": source.PolicyId,
		}
		if source.PolicyType != "" {
			entry["SourcePolicyType"] = source.PolicyType
		}
		statements = append(statements, entry)
	}
	return statements
}

// simulationResponse serialises the per-action evaluations, paginated over
// actions by Marker/MaxItems.
func simulationResponse(evaluations []SimulationEvaluation, params map[string]interface{}) map[string]interface{} {
	marker := request.GetStringParam(params, "Marker")
	maxItems := pagination.GetMaxItems(params, pagination.DefaultMaxItems)

	evaluationResults := make([]interface{}, 0, len(evaluations))
	for _, evaluation := range evaluations {
		resourceSpecific := make([]interface{}, 0, len(evaluation.ResourceSpecificResults))
		for _, resource := range evaluation.ResourceSpecificResults {
			entry := map[string]interface{}{
				"EvalResourceName":     resource.EvalResourceName,
				"EvalResourceDecision": resource.EvalResourceDecision,
				"MatchedStatements":    simulationStatementsToResponse(resource.MatchedStatements),
				"MissingContextValues": stringListToResponse(resource.MissingContextValues),
			}
			if evaluation.HasBoundary {
				entry["PermissionsBoundaryDecisionDetail"] = map[string]interface{}{
					"AllowedByPermissionsBoundary": resource.AllowedByBoundary,
				}
			}
			if details := simulationDecisionDetailsToResponse(resource.EvalDecisionDetails); details != nil {
				entry["EvalDecisionDetails"] = details
			}
			resourceSpecific = append(resourceSpecific, entry)
		}

		resultEntry := map[string]interface{}{
			"EvalActionName":          evaluation.EvalActionName,
			"EvalResourceName":        evaluation.EvalResourceName,
			"EvalDecision":            evaluation.EvalDecision,
			"MatchedStatements":       simulationStatementsToResponse(evaluation.MatchedStatements),
			"MissingContextValues":    stringListToResponse(evaluation.MissingContextValues),
			"ResourceSpecificResults": resourceSpecific,
		}
		if evaluation.HasBoundary {
			resultEntry["PermissionsBoundaryDecisionDetail"] = map[string]interface{}{
				"AllowedByPermissionsBoundary": evaluation.AllowedByBoundary,
			}
		}
		// The organisation decision detail exists only at the top level;
		// the per-resource shape carries no such member.
		if evaluation.HasOrgPolicies {
			resultEntry["OrganizationsDecisionDetail"] = map[string]interface{}{
				"AllowedByOrganizations": evaluation.AllowedByOrganizations,
			}
		}
		if details := simulationDecisionDetailsToResponse(evaluation.EvalDecisionDetails); details != nil {
			resultEntry["EvalDecisionDetails"] = details
		}
		evaluationResults = append(evaluationResults, resultEntry)
	}

	// Evaluation results carry no unique natural key — a request may repeat
	// an action name — so the pages walk positions with an opaque marker.
	paged := pagination.PaginateSliceByPosition(evaluationResults, marker, maxItems)

	resp := map[string]interface{}{
		"EvaluationResults": paged.Items,
		"IsTruncated":       paged.IsTruncated,
	}
	if paged.NextMarker != "" {
		resp["Marker"] = paged.NextMarker
	}
	return resp
}

// stringListToResponse renders a string slice as the wire list, with an
// empty slice (not nil) for an absent member value.
func stringListToResponse(values []string) []interface{} {
	out := make([]interface{}, 0, len(values))
	for _, v := range values {
		out = append(out, v)
	}
	return out
}

// buildSimulationContextEntries parses the flattened ContextEntries
// parameter into structured entries, collecting every value of each
// (possibly multi-valued) context key.
func buildSimulationContextEntries(params map[string]interface{}) []SimulationContextEntry {
	var entries []SimulationContextEntry
	for _, entry := range request.GetListParam(params, "ContextEntries") {
		name, _ := entry["ContextKeyName"].(string)
		if name == "" {
			continue
		}
		values := request.GetStringList(entry, "ContextKeyValues")
		if len(values) == 0 {
			continue
		}
		entries = append(entries, SimulationContextEntry{
			ContextKeyName:   name,
			ContextKeyValues: values,
		})
	}
	return entries
}

package cloudwatchlogs

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// PutResourcePolicy creates or updates a resource policy.
func (s *LogsService) PutResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetParamLowerFirst(req.Parameters, "PolicyName")
	policyDocument := request.GetParamLowerFirst(req.Parameters, "PolicyDocument")

	if policyName == "" {
		return nil, ErrMissingParameter
	}

	if err := validatePolicyDocument(policyDocument); err != nil {
		return nil, err
	}

	resourceArn := request.GetParamLowerFirst(req.Parameters, "ResourceArn")
	expectedRevisionId := request.GetParamLowerFirst(req.Parameters, "ExpectedRevisionId")
	policyScope := "ACCOUNT"
	if resourceArn != "" {
		policyScope = "RESOURCE"
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	existing, _ := store.GetResourcePolicy(policyName)
	if expectedRevisionId != "" && existing != nil && existing.RevisionId != expectedRevisionId {
		return nil, NewLogsError("InvalidParameterException",
			"Revision ID mismatch: expected "+expectedRevisionId+", got "+existing.RevisionId, 400)
	}
	rp := &logsstore.ResourcePolicy{
		PolicyName:     policyName,
		PolicyDocument: policyDocument,
		ResourceArn:    resourceArn,
		PolicyScope:    policyScope,
	}
	if existing != nil {
		rp.RevisionId = existing.RevisionId
	}

	if err := store.PutResourcePolicy(rp); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"resourcePolicy": formatResourcePolicy(rp),
		"revisionId":     rp.RevisionId,
	}, nil
}

// DeleteResourcePolicy deletes a resource policy.
func (s *LogsService) DeleteResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetParamLowerFirst(req.Parameters, "PolicyName")
	if policyName == "" {
		return nil, ErrMissingParameter
	}

	expectedRevisionId := request.GetParamLowerFirst(req.Parameters, "ExpectedRevisionId")
	if expectedRevisionId != "" {
		store, err := s.store(reqCtx)
		if err != nil {
			return nil, err
		}
		existing, _ := store.GetResourcePolicy(policyName)
		if existing != nil && existing.RevisionId != expectedRevisionId {
			return nil, NewLogsError("InvalidParameterException",
				"Revision ID mismatch: expected "+expectedRevisionId+", got "+existing.RevisionId, 400)
		}
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteResourcePolicy(policyName); err != nil {
		return nil, mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// DescribeResourcePolicies lists resource policies.
func (s *LogsService) DescribeResourcePolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetParamLowerFirst(req.Parameters, "ResourceArn")
	policyScopeFilter := request.GetParamLowerFirst(req.Parameters, "PolicyScope")
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")
	limit, err := validateListLimit(int32(request.GetIntParam(req.Parameters, "Limit")), 50, 50)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	allPolicies, err := store.ListResourcePolicies(resourceArn)
	if err != nil {
		return nil, mapStoreError(err)
	}

	if policyScopeFilter != "" {
		var filtered []*logsstore.ResourcePolicy
		for _, p := range allPolicies {
			if p.PolicyScope == policyScopeFilter {
				filtered = append(filtered, p)
			}
		}
		allPolicies = filtered
	}

	result := pagination.PaginateSlice(allPolicies, nextToken, int(limit), func(p *logsstore.ResourcePolicy) string {
		return p.PolicyName
	})

	policies := make([]map[string]interface{}, len(result.Items))
	for i, p := range result.Items {
		policies[i] = formatResourcePolicy(p)
	}

	resp := map[string]interface{}{
		"resourcePolicies": policies,
	}
	if result.NextMarker != "" {
		resp["nextToken"] = result.NextMarker
	}

	return resp, nil
}

func formatResourcePolicy(p *logsstore.ResourcePolicy) map[string]interface{} {
	result := map[string]interface{}{
		"policyName":      p.PolicyName,
		"policyDocument":  p.PolicyDocument,
		"lastUpdatedTime": p.LastUpdatedTime,
	}
	if p.ResourceArn != "" {
		result["resourceArn"] = p.ResourceArn
	}
	if p.PolicyScope != "" {
		result["policyScope"] = p.PolicyScope
	}
	if p.RevisionId != "" {
		result["revisionId"] = p.RevisionId
	}
	return result
}

package cloudwatchlogs

import (
	"context"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// --- Core methods ---

func (s *LogsService) putResourcePolicyCore(policyName, policyDocument, resourceArn, expectedRevisionId, region string) (*logsstore.ResourcePolicy, error) {
	if policyName == "" {
		return nil, ErrMissingParameter
	}
	if err := validatePolicyDocumentJSON(policyDocument); err != nil {
		return nil, err
	}

	policyScope := "ACCOUNT"
	if resourceArn != "" {
		policyScope = "RESOURCE"
	}

	store, err := s.getLogsStoreByRegion(region)
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
	return rp, nil
}

func (s *LogsService) deleteResourcePolicyCore(policyName, expectedRevisionId, region string) error {
	if policyName == "" {
		return ErrMissingParameter
	}

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return err
	}

	if expectedRevisionId != "" {
		existing, _ := store.GetResourcePolicy(policyName)
		if existing != nil && existing.RevisionId != expectedRevisionId {
			return NewLogsError("InvalidParameterException",
				"Revision ID mismatch: expected "+expectedRevisionId+", got "+existing.RevisionId, 400)
		}
	}

	if err := store.DeleteResourcePolicy(policyName); err != nil {
		return mapStoreError(err)
	}
	return nil
}

func (s *LogsService) describeResourcePoliciesCore(resourceArn, policyScopeFilter, nextToken, region string, limit int32) ([]*logsstore.ResourcePolicy, string, error) {
	l, err := validateListLimit(limit, 50, 50)
	if err != nil {
		return nil, "", err
	}

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return nil, "", err
	}

	allPolicies, err := store.ListResourcePolicies(resourceArn)
	if err != nil {
		return nil, "", mapStoreError(err)
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

	result := pagination.PaginateSlice(allPolicies, nextToken, int(l), func(p *logsstore.ResourcePolicy) string {
		return p.PolicyName
	})

	return result.Items, result.NextMarker, nil
}

// --- HTTP handlers ---

func (s *LogsService) PutResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetParamLowerFirst(req.Parameters, "PolicyName")
	policyDocument := request.GetParamLowerFirst(req.Parameters, "PolicyDocument")
	resourceArn := request.GetParamLowerFirst(req.Parameters, "ResourceArn")
	expectedRevisionId := request.GetParamLowerFirst(req.Parameters, "ExpectedRevisionId")

	rp, err := s.putResourcePolicyCore(policyName, policyDocument, resourceArn, expectedRevisionId, reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"resourcePolicy": formatResourcePolicy(rp),
		"revisionId":     rp.RevisionId,
	}, nil
}

func (s *LogsService) DeleteResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetParamLowerFirst(req.Parameters, "PolicyName")
	expectedRevisionId := request.GetParamLowerFirst(req.Parameters, "ExpectedRevisionId")

	if err := s.deleteResourcePolicyCore(policyName, expectedRevisionId, reqCtx.GetRegion()); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

func (s *LogsService) DescribeResourcePolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	resourceArn := request.GetParamLowerFirst(req.Parameters, "ResourceArn")
	policyScopeFilter := request.GetParamLowerFirst(req.Parameters, "PolicyScope")
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")
	limit := int32(request.GetIntParam(req.Parameters, "Limit"))

	policies, nextMarker, err := s.describeResourcePoliciesCore(resourceArn, policyScopeFilter, nextToken, reqCtx.GetRegion(), limit)
	if err != nil {
		return nil, err
	}

	formatted := make([]map[string]interface{}, len(policies))
	for i, p := range policies {
		formatted[i] = formatResourcePolicy(p)
	}

	resp := map[string]interface{}{
		"resourcePolicies": formatted,
	}
	if nextMarker != "" {
		resp["nextToken"] = nextMarker
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

package cloudwatchlogs

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// PutAccountPolicy creates or updates an account-level policy.
func (s *LogsService) PutAccountPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetParamLowerFirst(req.Parameters, "PolicyName")
	policyDocument := request.GetParamLowerFirst(req.Parameters, "PolicyDocument")
	policyType := request.GetParamLowerFirst(req.Parameters, "PolicyType")

	if policyName == "" || policyType == "" {
		return nil, ErrMissingParameter
	}

	if err := validatePolicyDocument(policyDocument); err != nil {
		return nil, err
	}

	if !validatePolicyType(policyType) {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid policyType: %s. Allowed values: DATA_PROTECTION_POLICY, SUBSCRIPTION_FILTER_POLICY, FIELD_INDEX_POLICY, TRANSFORMER_POLICY, METRIC_EXTRACTION_POLICY", policyType), 400)
	}

	scope := request.GetParamLowerFirst(req.Parameters, "Scope")
	if scope == "" {
		scope = "ALL"
	}
	if scope != "ALL" {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid scope: %s. Allowed values: ALL", scope), 400)
	}
	selectionCriteria := request.GetParamLowerFirst(req.Parameters, "SelectionCriteria")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	ap := &logsstore.AccountPolicy{
		PolicyName:        policyName,
		PolicyDocument:    policyDocument,
		PolicyType:        policyType,
		Scope:             scope,
		SelectionCriteria: selectionCriteria,
		AccountId:         s.accountID,
	}

	if err := store.PutAccountPolicy(ap); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"accountPolicy": formatAccountPolicy(ap),
	}, nil
}

// DeleteAccountPolicy deletes an account-level policy.
func (s *LogsService) DeleteAccountPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetParamLowerFirst(req.Parameters, "PolicyName")
	policyType := request.GetParamLowerFirst(req.Parameters, "PolicyType")

	if policyName == "" || policyType == "" {
		return nil, ErrMissingParameter
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteAccountPolicyEntry(policyType, policyName); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{}, nil
}

// DescribeAccountPolicies lists account-level policies.
func (s *LogsService) DescribeAccountPolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyType := request.GetParamLowerFirst(req.Parameters, "PolicyType")
	policyName := request.GetParamLowerFirst(req.Parameters, "PolicyName")
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")

	// accountIdentifiers is a multi-account feature and is scope-out for this
	// edge/on-premises platform. The parameter is intentionally not processed.

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	allPolicies, err := store.ListAccountPolicies(policyType, policyName)
	if err != nil {
		return nil, mapStoreError(err)
	}

	result := pagination.PaginateSlice(allPolicies, nextToken, 50, func(p *logsstore.AccountPolicy) string {
		return p.PolicyName
	})

	policies := make([]map[string]interface{}, len(result.Items))
	for i, p := range result.Items {
		policies[i] = formatAccountPolicy(p)
	}

	resp := map[string]interface{}{
		"accountPolicies": policies,
	}
	if result.NextMarker != "" {
		resp["nextToken"] = result.NextMarker
	}

	return resp, nil
}

func formatAccountPolicy(p *logsstore.AccountPolicy) map[string]interface{} {
	result := map[string]interface{}{
		"policyName":      p.PolicyName,
		"policyDocument":  p.PolicyDocument,
		"policyType":      p.PolicyType,
		"lastUpdatedTime": p.LastUpdatedTime,
		"scope":           p.Scope,
	}
	if p.AccountId != "" {
		result["accountId"] = p.AccountId
	}
	if p.SelectionCriteria != "" {
		result["selectionCriteria"] = p.SelectionCriteria
	}
	return result
}

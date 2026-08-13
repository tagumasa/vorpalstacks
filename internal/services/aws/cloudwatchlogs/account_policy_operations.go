package cloudwatchlogs

import (
	"context"
	"fmt"

	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// --- Core methods ---

func (s *LogsService) putAccountPolicyCore(policyName, policyDocument, policyType, scope, selectionCriteria, region string) (*logsstore.AccountPolicy, error) {
	if policyName == "" || policyType == "" {
		return nil, ErrMissingParameter
	}
	if err := validatePolicyNamePrefix(policyName); err != nil {
		return nil, err
	}
	if err := validatePolicyDocumentJSON(policyDocument); err != nil {
		return nil, err
	}
	if !validatePolicyType(policyType) {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid policyType: %s. Allowed values: DATA_PROTECTION_POLICY, SUBSCRIPTION_FILTER_POLICY, FIELD_INDEX_POLICY, TRANSFORMER_POLICY, METRIC_EXTRACTION_POLICY", policyType), 400)
	}
	if scope == "" {
		scope = "ALL"
	}
	if scope != "ALL" {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid scope: %s. Allowed values: ALL", scope), 400)
	}
	if err := validateSelectionCriteria(selectionCriteria); err != nil {
		return nil, err
	}

	store, err := s.getLogsStoreByRegion(region)
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
	return ap, nil
}

func (s *LogsService) deleteAccountPolicyCore(policyName, policyType, region string) error {
	if policyName == "" || policyType == "" {
		return ErrMissingParameter
	}

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return err
	}

	if err := store.DeleteAccountPolicyEntry(policyType, policyName); err != nil {
		return mapStoreError(err)
	}
	return nil
}

func (s *LogsService) describeAccountPoliciesCore(policyType, policyName, nextToken, region string) ([]*logsstore.AccountPolicy, string, error) {
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return nil, "", err
	}

	allPolicies, err := store.ListAccountPolicies(policyType, policyName)
	if err != nil {
		return nil, "", mapStoreError(err)
	}

	result := pagination.PaginateSlice(allPolicies, nextToken, 50, func(p *logsstore.AccountPolicy) string {
		return p.PolicyName
	})

	return result.Items, result.NextMarker, nil
}

// --- HTTP handlers ---

func (s *LogsService) PutAccountPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetParamLowerFirst(req.Parameters, "PolicyName")
	policyDocument := request.GetParamLowerFirst(req.Parameters, "PolicyDocument")
	policyType := request.GetParamLowerFirst(req.Parameters, "PolicyType")
	scope := request.GetParamLowerFirst(req.Parameters, "Scope")
	selectionCriteria := request.GetParamLowerFirst(req.Parameters, "SelectionCriteria")

	ap, err := s.putAccountPolicyCore(policyName, policyDocument, policyType, scope, selectionCriteria, reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"accountPolicy": formatAccountPolicy(ap),
	}, nil
}

func (s *LogsService) DeleteAccountPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetParamLowerFirst(req.Parameters, "PolicyName")
	policyType := request.GetParamLowerFirst(req.Parameters, "PolicyType")

	if err := s.deleteAccountPolicyCore(policyName, policyType, reqCtx.GetRegion()); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

func (s *LogsService) DescribeAccountPolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyType := request.GetParamLowerFirst(req.Parameters, "PolicyType")
	policyName := request.GetParamLowerFirst(req.Parameters, "PolicyName")
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")

	policies, nextMarker, err := s.describeAccountPoliciesCore(policyType, policyName, nextToken, reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	formatted := make([]map[string]interface{}, len(policies))
	for i, p := range policies {
		formatted[i] = formatAccountPolicy(p)
	}

	resp := map[string]interface{}{
		"accountPolicies": formatted,
	}
	if nextMarker != "" {
		resp["nextToken"] = nextMarker
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

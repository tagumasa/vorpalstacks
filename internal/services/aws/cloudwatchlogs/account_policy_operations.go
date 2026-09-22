package cloudwatchlogs

import (
	"context"
	"fmt"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// --- Core methods ---

func (s *LogsService) putAccountPolicyCore(policyName, policyDocument, policyType, scope, selectionCriteria, region string) (*logsstore.AccountPolicy, error) {
	if policyName == "" || policyType == "" {
		return nil, errRequiredMember("policyName and policyType")
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
	// The account-level data-protection policy carries the same document
	// contract as its group-level twin: "This policy must include two JSON
	// blocks" and "The JSON specified in policyDocument can be up to 30,720
	// characters long" — both from PutAccountPolicy's policyDocument
	// member documentation.
	if policyType == "DATA_PROTECTION_POLICY" {
		if err := validateDataProtectionPolicyDocument(policyDocument); err != nil {
			return nil, err
		}
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
	// "Specifying selectionCriteria is valid only when you specify
	// SUBSCRIPTION_FILTER_POLICY, FIELD_INDEX_POLICY or
	// TRANSFORMER_POLICY for policyType."
	if selectionCriteria != "" &&
		policyType != "SUBSCRIPTION_FILTER_POLICY" && policyType != "FIELD_INDEX_POLICY" && policyType != "TRANSFORMER_POLICY" {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("Specifying selectionCriteria is valid only when policyType is SUBSCRIPTION_FILTER_POLICY, FIELD_INDEX_POLICY or TRANSFORMER_POLICY, not %s", policyType), 400)
	}
	// "If policyType is SUBSCRIPTION_FILTER_POLICY, the only supported
	// selectionCriteria filter is LogGroupName NOT IN []".
	if policyType == "SUBSCRIPTION_FILTER_POLICY" {
		if err := validateSubscriptionSelectionCriteria(selectionCriteria); err != nil {
			return nil, err
		}
	}
	if policyType == "TRANSFORMER_POLICY" {
		if err := s.validateTransformerPolicySelection(region, policyName, selectionCriteria); err != nil {
			return nil, err
		}
		// "A transformer policy must include one JSON block with the
		// array of processors and their configurations" — the document
		// is the processor array itself, and it must be a valid
		// transformer (AWS rejects a malformed document at Put, so it
		// never reaches the ingestion seam's skip path).
		config := transformerPolicyConfig(policyDocument)
		if config == nil {
			return nil, NewLogsError("InvalidParameterException",
				"A transformer policy must include one JSON block with the array of processors and their configurations", 400)
		}
		if err := validateTransformerConfig(config); err != nil {
			return nil, err
		}
	}
	if policyType == "FIELD_INDEX_POLICY" {
		if err := s.validateFieldIndexPolicySelection(region, policyName, policyDocument, selectionCriteria); err != nil {
			return nil, err
		}
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
	// The invalidation follows the committed write, as the delete path
	// orders it: invalidating first opens a window where a concurrent
	// resolver re-caches the old policy with a fresh TTL after the
	// invalidation but before the new policy commits.
	if policyType == "TRANSFORMER_POLICY" {
		invalidateAllTransformerCaches()
	}
	return ap, nil
}

func (s *LogsService) deleteAccountPolicyCore(policyName, policyType, region string) error {
	if policyName == "" {
		return errRequiredMember("policyName")
	}
	if policyType == "" {
		return errRequiredMember("policyType")
	}
	if !validatePolicyType(policyType) {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid policyType: %s. Allowed values: DATA_PROTECTION_POLICY, SUBSCRIPTION_FILTER_POLICY, FIELD_INDEX_POLICY, TRANSFORMER_POLICY, METRIC_EXTRACTION_POLICY", policyType), 400)
	}

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return err
	}

	if err := store.DeleteAccountPolicyEntry(policyType, policyName); err != nil {
		return mapStoreError(err)
	}
	if policyType == "TRANSFORMER_POLICY" {
		invalidateAllTransformerCaches()
	}
	return nil
}

func (s *LogsService) describeAccountPoliciesCore(policyType, policyName string, accountIdentifiers []string, nextToken, region string) ([]*logsstore.AccountPolicy, string, error) {
	// The member is required ("Required: Yes", DescribeAccountPolicies
	// policyType — the API reference carries the trait the vendored model
	// revision predates) and rides the PolicyType enum; an invalid value
	// is a parameter error, not an empty listing.
	if policyType == "" {
		return nil, "", errRequiredMember("policyType")
	}
	if !validatePolicyType(policyType) {
		return nil, "", NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid policyType: %s. Allowed values: DATA_PROTECTION_POLICY, SUBSCRIPTION_FILTER_POLICY, FIELD_INDEX_POLICY, TRANSFORMER_POLICY, METRIC_EXTRACTION_POLICY", policyType), 400)
	}
	// "Currently, you can specify only one account ID in this parameter"
	// over the fixed-length-12 form; a source account other than the
	// local one owns no policies on this single-account platform.
	if err := validateAccountIdentifierList(accountIdentifiers); err != nil {
		return nil, "", err
	}
	if len(accountIdentifiers) > 1 {
		return nil, "", NewLogsError("InvalidParameterException",
			"Currently, you can specify only one account ID in the accountIdentifiers parameter", 400)
	}
	if len(accountIdentifiers) == 1 && accountIdentifiers[0] != s.accountID {
		return []*logsstore.AccountPolicy{}, "", nil
	}

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return nil, "", err
	}

	allPolicies, err := store.ListAccountPolicies(policyType, policyName)
	if err != nil {
		return nil, "", mapStoreError(err)
	}

	// The page marker rides the scoped listing vocabulary: the token
	// carries the request identity and its documented expiry rather than
	// a bare policy name. The operation carries no limit member, so the
	// page rides the default describe bound.
	scope := listingScope("accountpolicies", policyType, policyName, strings.Join(accountIdentifiers, ","))
	result, err := paginateScopedListing(scope, nextToken, allPolicies, logsstore.DefaultDescribeLimit, func(p *logsstore.AccountPolicy) string {
		return p.PolicyName
	})
	if err != nil {
		return nil, "", err
	}
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

	return response.EmptyResponse(), nil
}

func (s *LogsService) DescribeAccountPolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyType := request.GetParamLowerFirst(req.Parameters, "PolicyType")
	policyName := request.GetParamLowerFirst(req.Parameters, "PolicyName")
	accountIdentifiers := request.GetStringList(req.Parameters, "AccountIdentifiers")
	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")

	policies, nextMarker, err := s.describeAccountPoliciesCore(policyType, policyName, accountIdentifiers, nextToken, reqCtx.GetRegion())
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

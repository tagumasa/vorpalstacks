package cloudwatchlogs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// --- Core methods ---

func (s *LogsService) putResourcePolicyCore(policyName, policyDocument, resourceArn, expectedRevisionId, region string) (*logsstore.ResourcePolicy, error) {
	if policyName == "" {
		return nil, errRequiredMember("policyName")
	}
	if err := validatePolicyDocumentJSON(policyDocument); err != nil {
		return nil, err
	}
	if err := validateResourcePolicyArn(resourceArn); err != nil {
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

	// The admission — the quota census, the optimistic-revision checks,
	// the mint and the write — commits under the store's record mutex:
	// "Policy limits - An account can have a maximum of 10 policies
	// without resourceARN and one per LogGroup resourceARN"
	// (LimitExceededException is the operation's declared quota identity)
	// and the revision guard both decide on reads the bare write never
	// checks, so two concurrent Puts quoting the same revision admit
	// exactly one. Replacing a policy under its own name keeps its slot.
	var admitted *logsstore.ResourcePolicy
	if err := store.MutateResourcePolicy(policyName, func(existing *logsstore.ResourcePolicy) (*logsstore.ResourcePolicy, error) {
		if err := refuseResourcePolicyQuota(store, policyName, resourceArn, existing); err != nil {
			return nil, err
		}
		// "The expected revision ID of the resource policy. Required when
		// resourceArn is provided to prevent concurrent modifications. Use
		// null when creating a resource policy for the first time" — a put
		// against an existing resource-scoped policy must quote its revision;
		// the first-time null and the account-wide form carry none.
		if resourceArn != "" && existing != nil && expectedRevisionId == "" {
			return nil, NewLogsError("InvalidParameterException",
				"expectedRevisionId is required when updating a resource-scoped resource policy; quote the revision the last put returned", 400)
		}
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
		// "The revision ID of the created or updated resource policy. Only
		// returned for resource-scoped policies" — the account-wide form
		// mints none, so the revision identity is a resource-scoped fact
		// alone. Every successful resource-scoped put mints a fresh
		// revision — the revisionId the response returns and a later
		// ExpectedRevisionId must quote; carrying the previous value
		// forward would leave the optimistic-concurrency check comparing
		// stale against stale.
		if policyScope == "RESOURCE" {
			rp.RevisionId = newResourcePolicyRevisionId()
		}
		admitted = rp
		return rp, nil
	}); err != nil {
		return nil, mapStoreError(err)
	}
	return admitted, nil
}

// refuseResourcePolicyQuota enforces the documented policy limits: "An
// account can have a maximum of 10 policies without resourceARN and one
// per LogGroup resourceARN". A put replacing the named policy keeps its
// own slot under either counting.
func refuseResourcePolicyQuota(store *logsstore.Store, policyName, resourceArn string, existing *logsstore.ResourcePolicy) error {
	all, err := store.ListResourcePolicies("")
	if err != nil {
		return mapStoreError(err)
	}
	if resourceArn == "" {
		count := 0
		for _, p := range all {
			if p.ResourceArn != "" || p.PolicyName == policyName {
				continue
			}
			count++
		}
		if count+1 > logsstore.MaxAccountScopeResourcePolicies {
			return NewLogsError("LimitExceededException",
				fmt.Sprintf("An account can have a maximum of %d policies without resourceArn", logsstore.MaxAccountScopeResourcePolicies), 400)
		}
		return nil
	}
	for _, p := range all {
		if p.ResourceArn == resourceArn && p.PolicyName != policyName {
			return NewLogsError("LimitExceededException",
				"A LogGroup resourceArn can carry one resource policy; "+resourceArn+" already has one", 400)
		}
	}
	return nil
}

// newResourcePolicyRevisionId mints the revision identifier of one
// successful PutResourcePolicy: a random hex string (unpredictability is
// what makes the expected-revision check a concurrency guard rather than
// a guessable counter).
func newResourcePolicyRevisionId() string {
	return "rev-" + randomHexWithFallback(8, func() string { return fmt.Sprintf("%d", time.Now().UnixNano()) })
}

// validateResourcePolicyArn enforces the resourceArn form the member
// documents: "The ARN of the CloudWatch Logs resource to which the
// resource policy needs to be added or attached. Currently only supports
// LogGroup ARN." A logs-service ARN whose resource addresses a log group
// (tolerating the ":*" log-stream namespace suffix DescribeLogGroups
// appends to the object ARN); anything else is a parameter error, not a
// policy stored against an unusable target.
func validateResourcePolicyArn(resourceArn string) error {
	if resourceArn == "" {
		return nil
	}
	resource := resourceArn
	if strings.HasSuffix(resource, ":*") {
		resource = strings.TrimSuffix(resource, ":*")
	}
	_, service, _, _, res := svcarn.SplitARN(resource)
	if service != "logs" || !strings.HasPrefix(res, "log-group:") || strings.Contains(res, ":log-stream:") {
		return NewLogsError("InvalidParameterException",
			"resourceArn currently only supports LogGroup ARN: "+resourceArn, 400)
	}
	return nil
}

func (s *LogsService) deleteResourcePolicyCore(policyName, resourceArn, expectedRevisionId, region string) error {
	if policyName == "" {
		return errRequiredMember("policyName")
	}
	if err := validateResourcePolicyArn(resourceArn); err != nil {
		return err
	}

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return err
	}

	// The revision guard and the deletion commit under the store's
	// admission mutex, so a concurrent Put cannot slip a freshly minted
	// revision between the guard and the delete.
	if err := store.DeleteResourcePolicyAdmitted(policyName, func(existing *logsstore.ResourcePolicy) error {
		// "The expected revision ID of the resource policy. Required when
		// deleting a resource-scoped policy to prevent concurrent
		// modifications" — the account-wide form carries no revision
		// identity, so only the resource-scoped delete must quote one.
		if existing.PolicyScope == "RESOURCE" && expectedRevisionId == "" {
			return NewLogsError("InvalidParameterException",
				"expectedRevisionId is required when deleting a resource-scoped resource policy; quote the revision the last put returned", 400)
		}
		if expectedRevisionId != "" && existing.RevisionId != expectedRevisionId {
			return NewLogsError("InvalidParameterException",
				"Revision ID mismatch: expected "+expectedRevisionId+", got "+existing.RevisionId, 400)
		}
		return nil
	}); err != nil {
		return mapStoreError(err)
	}
	return nil
}

func (s *LogsService) describeResourcePoliciesCore(resourceArn, policyScopeFilter, nextToken, region string, limit int32) ([]*logsstore.ResourcePolicy, string, error) {
	l, err := validateListLimit(limit, logsstore.DefaultDescribeLimit, logsstore.MaxDescribeLimit)
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

	// "Valid values are ACCOUNT or RESOURCE. When not specified, defaults
	// to ACCOUNT." — the default is a scope selection, not the absence of
	// one, and a foreign value is a parameter error (the operation
	// declares InvalidParameterException), not an empty listing.
	if policyScopeFilter == "" {
		policyScopeFilter = "ACCOUNT"
	}
	if policyScopeFilter != "ACCOUNT" && policyScopeFilter != "RESOURCE" {
		return nil, "", NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid policyScope: %s. Allowed values: ACCOUNT, RESOURCE", policyScopeFilter), 400)
	}
	var filtered []*logsstore.ResourcePolicy
	for _, p := range allPolicies {
		if p.PolicyScope == policyScopeFilter {
			filtered = append(filtered, p)
		}
	}
	allPolicies = filtered

	// The listing pages through the scoped token vocabulary: the filter
	// members digest into the scope, so a token minted under one filter
	// repositions no other walk and a typed token rejects.
	scope := listingScope("resourcepolicies", resourceArn, policyScopeFilter)
	result, err := paginateScopedListing(scope, nextToken, allPolicies, int(l), func(p *logsstore.ResourcePolicy) string {
		return p.PolicyName
	})
	if err != nil {
		return nil, "", err
	}
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

	resp := map[string]interface{}{
		"resourcePolicy": formatResourcePolicy(rp),
	}
	// "The revision ID of the created or updated resource policy. Only
	// returned for resource-scoped policies."
	if rp.PolicyScope == "RESOURCE" {
		resp["revisionId"] = rp.RevisionId
	}
	return resp, nil
}

func (s *LogsService) DeleteResourcePolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policyName := request.GetParamLowerFirst(req.Parameters, "PolicyName")
	resourceArn := request.GetParamLowerFirst(req.Parameters, "ResourceArn")
	expectedRevisionId := request.GetParamLowerFirst(req.Parameters, "ExpectedRevisionId")

	if err := s.deleteResourcePolicyCore(policyName, resourceArn, expectedRevisionId, reqCtx.GetRegion()); err != nil {
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

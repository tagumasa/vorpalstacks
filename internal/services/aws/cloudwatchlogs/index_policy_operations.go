package cloudwatchlogs

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// The field index family: the log-group-level policy surface
// (PutIndexPolicy/DescribeIndexPolicies/DeleteIndexPolicy) and the
// derived field-index listing (DescribeFieldIndexes). The index's AWS
// contract is query acceleration with no API-observable output
// difference; the platform's faithful surface is the policy storage, the
// effective-policy resolution (a group-level policy overrides the
// account-level FIELD_INDEX_POLICY whose LogGroupNamePrefix — or
// account-wide scope — applies to the group) and the derivation of the
// listing over the same event scan the query plane reads.
// --- HTTP handlers ---

func (s *LogsService) PutIndexPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policy, err := s.putIndexPolicyCore(
		request.GetParamLowerFirst(req.Parameters, "LogGroupIdentifier"),
		request.GetParamLowerFirst(req.Parameters, "PolicyDocument"),
		reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"indexPolicy": formatIndexPolicy(policy.LogGroupIdentifier, policy.PolicyDocument, policy.LastUpdateTime, "LOG_GROUP", ""),
	}, nil
}

func (s *LogsService) DescribeIndexPolicies(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policies, err := s.describeIndexPoliciesCore(
		request.GetStringList(req.Parameters, "LogGroupIdentifiers"), reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}
	formatted := make([]map[string]interface{}, len(policies))
	for i, view := range policies {
		formatted[i] = formatIndexPolicy(view.groupARN, view.eff.policyDocument,
			view.eff.lastUpdateTime, view.eff.source, view.eff.policyName)
	}
	return map[string]interface{}{"indexPolicies": formatted}, nil
}

func (s *LogsService) DeleteIndexPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.deleteIndexPolicyCore(
		request.GetParamLowerFirst(req.Parameters, "LogGroupIdentifier"), reqCtx.GetRegion()); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func (s *LogsService) DescribeFieldIndexes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	entries, nextMarker, err := s.describeFieldIndexesCore(
		request.GetStringList(req.Parameters, "LogGroupIdentifiers"),
		request.GetStringList(req.Parameters, "IndexCategories"),
		reqCtx.GetRegion(),
		request.GetParamLowerFirst(req.Parameters, "NextToken"))
	if err != nil {
		return nil, err
	}
	formatted := make([]map[string]interface{}, len(entries))
	for i, entry := range entries {
		formatted[i] = formatFieldIndex(entry)
	}
	resp := map[string]interface{}{"fieldIndexes": formatted}
	if nextMarker != "" {
		resp["nextToken"] = nextMarker
	}
	return resp, nil
}

// formatIndexPolicy renders the model's IndexPolicy shape. The policy
// name is emitted only for account-level policies ("Responses about log
// group-level field index policies don't have this field, because those
// policies don't have names").
func formatIndexPolicy(logGroupIdentifier, policyDocument string, lastUpdateTime int64, source, policyName string) map[string]interface{} {
	out := map[string]interface{}{
		"logGroupIdentifier": logGroupIdentifier,
		"policyDocument":     policyDocument,
		"lastUpdateTime":     lastUpdateTime,
		"source":             source,
	}
	if policyName != "" {
		out["policyName"] = policyName
	}
	return out
}

// formatFieldIndex renders the model's FieldIndex shape; the time
// members emit only when the ingestion scan has produced them — the scan
// stamp for a group whose events were never scanned (the group ingested
// nothing) is absence, not a zero timestamp.
func formatFieldIndex(entry *fieldIndexEntry) map[string]interface{} {
	out := map[string]interface{}{
		"fieldIndexName": entry.fieldIndexName,
		"indexCategory":  entry.indexCategory,
		"type":           entry.indexType,
	}
	if entry.lastScanTime != 0 {
		out["lastScanTime"] = entry.lastScanTime
	}
	if entry.logGroupIdentifier != "" {
		out["logGroupIdentifier"] = entry.logGroupIdentifier
	}
	if entry.firstEventTime != 0 || entry.lastEventTime != 0 {
		out["firstEventTime"] = entry.firstEventTime
		out["lastEventTime"] = entry.lastEventTime
	}
	return out
}

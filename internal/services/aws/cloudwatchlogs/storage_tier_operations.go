package cloudwatchlogs

import (
	"context"

	"vorpalstacks/internal/common/request"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The account-level storage tier policy: the tier the account's log data
// is held in (STANDARD, or INTELLIGENT_TIERING where the service moves
// log data to the most cost-effective tier based on access frequency).
// The tiering choice has no API-observable output surface of its own —
// events read identically through every read operation — so the
// platform's faithful implementation is the policy record and its
// round-trip (recorded at the C5 adjudication that scheduled this
// family).

// putStorageTierPolicyCore validates and stores the account's storage
// tier policy. The storageTier member is the model's whole request: an
// absent member or a value outside the two-value enum rejects with the
// operation's InvalidParameterException.
func (s *LogsService) putStorageTierPolicyCore(storageTier, region string) (*logsstore.StorageTierPolicy, error) {
	if storageTier == "" {
		return nil, errRequiredMember("storageTier")
	}
	if storageTier != "STANDARD" && storageTier != "INTELLIGENT_TIERING" {
		return nil, NewLogsError("InvalidParameterException",
			"1 validation error detected: Value '"+storageTier+"' at 'storageTier' failed to satisfy constraint: Member must have value only [STANDARD, INTELLIGENT_TIERING]", 400)
	}
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return nil, err
	}
	return store.PutStorageTierPolicy(storageTier)
}

// getStorageTierPolicyCore reads the account's storage tier policy; an
// account that never set one has no policy record and answers the
// modelled ResourceNotFoundException (the family convention of the
// singleton policy getters: GetDataProtectionPolicy and
// GetDeliveryDestinationPolicy behave the same way).
func (s *LogsService) getStorageTierPolicyCore(region string) (*logsstore.StorageTierPolicy, error) {
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return nil, err
	}
	policy, err := store.GetStorageTierPolicy()
	if err != nil {
		return nil, mapStoreError(err)
	}
	return policy, nil
}

// --- HTTP handlers ---

func (s *LogsService) PutStorageTierPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policy, err := s.putStorageTierPolicyCore(
		request.GetParamLowerFirst(req.Parameters, "StorageTier"),
		reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}
	return formatStorageTierPolicy(policy), nil
}

func (s *LogsService) GetStorageTierPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	policy, err := s.getStorageTierPolicyCore(reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}
	return formatStorageTierPolicy(policy), nil
}

func formatStorageTierPolicy(policy *logsstore.StorageTierPolicy) map[string]interface{} {
	return map[string]interface{}{
		"storageTier":     policy.StorageTier,
		"lastUpdatedTime": policy.LastUpdatedTime,
	}
}

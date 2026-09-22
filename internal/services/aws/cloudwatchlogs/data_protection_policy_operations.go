package cloudwatchlogs

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// --- Core methods ---

func (s *LogsService) putDataProtectionPolicyCore(logGroupIdentifier, policyDocument, region string) (*logsstore.DataProtectionPolicy, error) {
	if logGroupIdentifier == "" {
		return nil, errRequiredMember("logGroupIdentifier")
	}
	// policyDocument is required, formatted as a JSON string and carrying
	// the two-block data protection structure the member documentation
	// mandates (the Audit block with its FindingsDestination and the
	// Deidentify block with its empty MaskConfig, over matching
	// DataIdentifer arrays).
	if err := validatePolicyDocumentJSON(policyDocument); err != nil {
		return nil, err
	}
	if err := validateDataProtectionPolicyDocument(policyDocument); err != nil {
		return nil, err
	}

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return nil, err
	}

	// The record keys on the resolved group name, so ARN and name input
	// address the same policy; the response echoes the request identifier.
	name := resolveLogGroupIdentifier(logGroupIdentifier)
	if _, err = store.GetLogGroup(name); err != nil {
		return nil, mapStoreError(err)
	}

	dpp := &logsstore.DataProtectionPolicy{
		LogGroupIdentifier: name,
		PolicyDocument:     policyDocument,
	}

	if err := store.PutDataProtectionPolicy(dpp); err != nil {
		return nil, mapStoreError(err)
	}
	// The group's display status follows the policy store: ACTIVATED while
	// a protection policy rides the group. The write runs after the policy
	// commit so a crash between the two leaves the status lagging one put
	// behind, never claiming protection a stored policy lacks.
	if err := store.MutateLogGroup(name, func(lg *logsstore.LogGroup) error {
		lg.DataProtectionStatus = "ACTIVATED"
		return nil
	}); err != nil {
		return nil, mapStoreError(err)
	}
	return dpp, nil
}

func (s *LogsService) getDataProtectionPolicyCore(logGroupIdentifier, region string) (*logsstore.DataProtectionPolicy, error) {
	if logGroupIdentifier == "" {
		return nil, errRequiredMember("logGroupIdentifier")
	}

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return nil, err
	}

	// The store sentinel must surface as the modelled
	// ResourceNotFoundException, like every other Core's store read.
	dpp, err := store.GetDataProtectionPolicy(resolveLogGroupIdentifier(logGroupIdentifier))
	if err != nil {
		return nil, mapStoreError(err)
	}
	return dpp, nil
}

func (s *LogsService) deleteDataProtectionPolicyCore(logGroupIdentifier, region string) error {
	if logGroupIdentifier == "" {
		return errRequiredMember("logGroupIdentifier")
	}

	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return err
	}

	if err := store.DeleteDataProtectionPolicy(resolveLogGroupIdentifier(logGroupIdentifier)); err != nil {
		return mapStoreError(err)
	}
	// "Displays whether this log group has a protection policy, or whether
	// it had one in the past": a deleted policy leaves the DELETED marker
	// behind — the group no longer has a policy but had one.
	if err := store.MutateLogGroup(resolveLogGroupIdentifier(logGroupIdentifier), func(lg *logsstore.LogGroup) error {
		lg.DataProtectionStatus = "DELETED"
		return nil
	}); err != nil {
		return mapStoreError(err)
	}
	return nil
}

// --- HTTP handlers ---

func (s *LogsService) PutDataProtectionPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupIdentifier := request.GetParamLowerFirst(req.Parameters, "LogGroupIdentifier")
	policyDocument := request.GetParamLowerFirst(req.Parameters, "PolicyDocument")

	dpp, err := s.putDataProtectionPolicyCore(logGroupIdentifier, policyDocument, reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	// The response echoes the identifier the request specified, which the
	// model documents as name-or-ARN.
	return map[string]interface{}{
		"logGroupIdentifier": logGroupIdentifier,
		"policyDocument":     dpp.PolicyDocument,
		"lastUpdatedTime":    dpp.LastUpdatedTime,
	}, nil
}

func (s *LogsService) GetDataProtectionPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupIdentifier := request.GetParamLowerFirst(req.Parameters, "LogGroupIdentifier")

	dpp, err := s.getDataProtectionPolicyCore(logGroupIdentifier, reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}

	// The response echoes the identifier the request specified, which the
	// model documents as name-or-ARN.
	return map[string]interface{}{
		"logGroupIdentifier": logGroupIdentifier,
		"policyDocument":     dpp.PolicyDocument,
		"lastUpdatedTime":    dpp.LastUpdatedTime,
	}, nil
}

func (s *LogsService) DeleteDataProtectionPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupIdentifier := request.GetParamLowerFirst(req.Parameters, "LogGroupIdentifier")

	if err := s.deleteDataProtectionPolicyCore(logGroupIdentifier, reqCtx.GetRegion()); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

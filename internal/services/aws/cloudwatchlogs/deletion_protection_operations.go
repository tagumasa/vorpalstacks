package cloudwatchlogs

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// PutLogGroupDeletionProtection updates the deletion protection setting
// for the specified log group. The modelled input is logGroupIdentifier
// (name or ARN) plus deletionProtectionEnabled, which is required — an
// explicit false is a legitimate value, so the wire presence flag travels
// alongside it.
func (s *LogsService) PutLogGroupDeletionProtection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_, deletionProtectionSet := req.Parameters["deletionProtectionEnabled"]
	if !deletionProtectionSet {
		_, deletionProtectionSet = req.Parameters["DeletionProtectionEnabled"]
	}
	input := PutLogGroupDeletionProtectionInput{
		LogGroupIdentifier:        request.GetParamLowerFirst(req.Parameters, "LogGroupIdentifier"),
		DeletionProtectionEnabled: request.GetBoolParam(req.Parameters, "DeletionProtectionEnabled"),
		DeletionProtectionSet:     deletionProtectionSet,
		Region:                    reqCtx.GetRegion(),
	}

	if err := s.putLogGroupDeletionProtectionCore(input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

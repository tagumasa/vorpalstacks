package cloudwatchlogs

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// PutLogGroupDeletionProtection updates the deletion protection setting
// for the specified log group.
func (s *LogsService) PutLogGroupDeletionProtection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := PutLogGroupDeletionProtectionInput{
		LogGroupIdentifier:        request.GetParamLowerFirst(req.Parameters, "LogGroupIdentifier"),
		LogGroupName:              request.GetParamLowerFirst(req.Parameters, "LogGroupName"),
		DeletionProtectionEnabled: request.GetBoolParam(req.Parameters, "DeletionProtectionEnabled"),
		Region:                    reqCtx.GetRegion(),
	}

	if err := s.putLogGroupDeletionProtectionCore(input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

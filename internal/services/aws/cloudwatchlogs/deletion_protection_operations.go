package cloudwatchlogs

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// PutLogGroupDeletionProtection updates the deletion protection setting
// for the specified log group.
func (s *LogsService) PutLogGroupDeletionProtection(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupIdentifier := request.GetParamLowerFirst(req.Parameters, "LogGroupIdentifier")
	if logGroupIdentifier == "" {
		logGroupIdentifier = request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	}

	if logGroupIdentifier == "" {
		return nil, ErrMissingParameter
	}

	enabled := request.GetBoolParam(req.Parameters, "DeletionProtectionEnabled")

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	lg, err := store.GetLogGroup(logGroupIdentifier)
	if err != nil {
		return nil, mapStoreError(err)
	}

	lg.DeletionProtectionEnabled = enabled
	if err := store.PutLogGroup(lg); err != nil {
		return nil, mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

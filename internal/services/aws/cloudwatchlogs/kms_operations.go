package cloudwatchlogs

import (
	"context"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/utils/aws/arn"
)

// AssociateKmsKey associates a KMS key with the specified log group.
func (s *LogsService) AssociateKmsKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	resourceIdentifier := request.GetParamLowerFirst(req.Parameters, "ResourceIdentifier")
	kmsKeyId := request.GetParamLowerFirst(req.Parameters, "KmsKeyId")

	if kmsKeyId == "" {
		return nil, ErrMissingParameter
	}
	if logGroupName == "" && resourceIdentifier == "" {
		return nil, ErrMissingParameter
	}

	target := logGroupName
	if target == "" {
		target = resourceIdentifier
	}

	parsed, err := arn.ParseARN(kmsKeyId)
	if err != nil || parsed.Service != "kms" {
		return nil, NewLogsError("InvalidParameterException",
			"kmsKeyId must be a valid KMS key ARN", 400)
	}
	if !strings.HasPrefix(parsed.Resource, "key/") && !strings.HasPrefix(parsed.Resource, "alias/") {
		return nil, NewLogsError("InvalidParameterException",
			"kmsKeyId resource must be a key UUID (key/...) or alias (alias/...)", 400)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	lg, err := store.GetLogGroup(target)
	if err != nil {
		return nil, mapStoreError(err)
	}

	lg.KmsKeyId = kmsKeyId
	if err := store.PutLogGroup(lg); err != nil {
		return nil, mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

// DisassociateKmsKey removes the KMS key association from the specified log group.
func (s *LogsService) DisassociateKmsKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")
	resourceIdentifier := request.GetParamLowerFirst(req.Parameters, "ResourceIdentifier")

	if logGroupName == "" && resourceIdentifier == "" {
		return nil, ErrMissingParameter
	}

	target := logGroupName
	if target == "" {
		target = resourceIdentifier
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	lg, err := store.GetLogGroup(target)
	if err != nil {
		return nil, mapStoreError(err)
	}

	lg.KmsKeyId = ""
	if err := store.PutLogGroup(lg); err != nil {
		return nil, mapStoreError(err)
	}

	return response.EmptyResponse(), nil
}

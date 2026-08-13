package cloudwatchlogs

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
)

// AssociateKmsKey associates a KMS key with the specified log group.
func (s *LogsService) AssociateKmsKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := AssociateKmsKeyInput{
		LogGroupName:       request.GetParamLowerFirst(req.Parameters, "LogGroupName"),
		ResourceIdentifier: request.GetParamLowerFirst(req.Parameters, "ResourceIdentifier"),
		KmsKeyId:           request.GetParamLowerFirst(req.Parameters, "KmsKeyId"),
		Region:             reqCtx.GetRegion(),
	}

	if err := s.associateKmsKeyCore(input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DisassociateKmsKey removes the KMS key association from the specified log group.
func (s *LogsService) DisassociateKmsKey(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := DisassociateKmsKeyInput{
		LogGroupName:       request.GetParamLowerFirst(req.Parameters, "LogGroupName"),
		ResourceIdentifier: request.GetParamLowerFirst(req.Parameters, "ResourceIdentifier"),
		Region:             reqCtx.GetRegion(),
	}

	if err := s.disassociateKmsKeyCore(input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

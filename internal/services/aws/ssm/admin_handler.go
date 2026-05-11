package ssm

import (
	"context"
	"fmt"
	"net/http"

	svcerrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/utils/timeutils"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/ssm"
	ssmconnect "vorpalstacks/internal/pb/aws/ssm/ssmconnect"
	ssmstore "vorpalstacks/internal/store/aws/ssm"
)

// AdminHandler implements the SSM admin console gRPC-Web handler.
type AdminHandler struct {
	ssmconnect.UnimplementedSSMServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ ssmconnect.SSMServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new SSM admin handler with the given storage manager and account ID.
func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
	}
}

func (h *AdminHandler) getStoreFromHeaders(headers http.Header) (*ssmstore.Store, error) {
	region := svccommon.GetRegionFromHeader(headers)
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return ssmstore.NewStore(regionStorage, h.accountId, region), nil
}

// DescribeParameters retrieves SSM parameters from the store, applying optional filters and pagination.
func (h *AdminHandler) DescribeParameters(ctx context.Context, req *connect.Request[pb.DescribeParametersRequest]) (*connect.Response[pb.DescribeParametersResult], error) {
	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	maxResults := req.Msg.Maxresults
	if maxResults <= 0 {
		maxResults = 50
	}

	filters := make(map[string]string)
	for _, f := range req.Msg.Filters {
		if len(f.Values) > 0 {
			filters[f.Key.String()] = f.Values[0]
		}
	}

	params, nextToken, err := store.DescribeParameters(filters, maxResults, req.Msg.Nexttoken)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	var metadataList []*pb.ParameterMetadata
	for _, p := range params {
		meta := &pb.ParameterMetadata{
			Name:             p.Name,
			Version:          p.Version,
			Lastmodifieddate: p.LastModifiedDate.Format(timeutils.ISO8601UTCFormat),
			Datatype:         p.DataType,
			Arn:              p.ARN,
		}
		if p.Description != "" {
			meta.Description = p.Description
		}
		if p.KeyID != "" {
			meta.Keyid = p.KeyID
		}
		if p.AllowedPattern != "" {
			meta.Allowedpattern = p.AllowedPattern
		}
		switch p.Type {
		case ssmstore.ParameterTypeString:
			meta.Type = pb.ParameterType_PARAMETER_TYPE_STRING
		case ssmstore.ParameterTypeStringList:
			meta.Type = pb.ParameterType_PARAMETER_TYPE_STRING_LIST
		case ssmstore.ParameterTypeSecureString:
			meta.Type = pb.ParameterType_PARAMETER_TYPE_SECURE_STRING
		}
		switch p.Tier {
		case ssmstore.ParameterTierStandard:
			meta.Tier = pb.ParameterTier_PARAMETER_TIER_STANDARD
		case ssmstore.ParameterTierAdvanced:
			meta.Tier = pb.ParameterTier_PARAMETER_TIER_ADVANCED
		case ssmstore.ParameterTierIntelligentTiering:
			meta.Tier = pb.ParameterTier_PARAMETER_TIER_INTELLIGENT_TIERING
		}
		metadataList = append(metadataList, meta)
	}

	return connect.NewResponse(&pb.DescribeParametersResult{
		Parameters: metadataList,
		Nexttoken:  nextToken,
	}), nil
}

// PutParameter creates or updates an SSM parameter via the admin console.
func (h *AdminHandler) PutParameter(ctx context.Context, req *connect.Request[pb.PutParameterRequest]) (*connect.Response[pb.PutParameterResult], error) {
	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("name is required"))
	}

	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	paramType := ssmstore.ParameterTypeString
	switch req.Msg.Type {
	case pb.ParameterType_PARAMETER_TYPE_STRING_LIST:
		paramType = ssmstore.ParameterTypeStringList
	case pb.ParameterType_PARAMETER_TYPE_SECURE_STRING:
		paramType = ssmstore.ParameterTypeSecureString
	}

	param := &ssmstore.Parameter{
		Name:        req.Msg.Name,
		Value:       req.Msg.Value,
		Type:        paramType,
		Description: req.Msg.Description,
		DataType:    req.Msg.Datatype,
		KeyID:       req.Msg.Keyid,
	}

	switch req.Msg.Tier {
	case pb.ParameterTier_PARAMETER_TIER_ADVANCED:
		param.Tier = ssmstore.ParameterTierAdvanced
	case pb.ParameterTier_PARAMETER_TIER_INTELLIGENT_TIERING:
		param.Tier = ssmstore.ParameterTierIntelligentTiering
	}

	if len(req.Msg.Tags) > 0 {
		tags := make(map[string]string, len(req.Msg.Tags))
		for _, t := range req.Msg.Tags {
			if t.Key != "" {
				tags[t.Key] = t.Value
			}
		}
		param.Tags = tags
	}

	overwrite := req.Msg.Overwrite
	version, err := store.PutParameter(param, overwrite)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.PutParameterResult{
		Version: version,
	}), nil
}

// DeleteParameter deletes an SSM parameter via the admin console.
func (h *AdminHandler) DeleteParameter(ctx context.Context, req *connect.Request[pb.DeleteParameterRequest]) (*connect.Response[pb.DeleteParameterResult], error) {
	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("name is required"))
	}

	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if err := store.DeleteParameter(req.Msg.Name); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DeleteParameterResult{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the SSM admin console.
func NewConnectHandler(sm *storage.RegionStorageManager, accountID string) (string, http.Handler) {
	return ssmconnect.NewSSMServiceHandler(NewAdminHandler(sm, accountID))
}

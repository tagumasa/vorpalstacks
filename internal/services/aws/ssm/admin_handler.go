package ssm

import (
	"context"
	"net/http"

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

func (h *AdminHandler) getStore(req *connect.Request[pb.DescribeParametersRequest]) (*ssmstore.Store, error) {
	region := svccommon.GetRegionFromHeader(req.Header())
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return ssmstore.NewStore(regionStorage, h.accountId, region), nil
}

// DescribeParameters retrieves SSM parameters from the store, applying optional filters and pagination.
func (h *AdminHandler) DescribeParameters(ctx context.Context, req *connect.Request[pb.DescribeParametersRequest]) (*connect.Response[pb.DescribeParametersResult], error) {
	store, err := h.getStore(req)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
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
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var metadataList []*pb.ParameterMetadata
	for _, p := range params {
		meta := &pb.ParameterMetadata{
			Name:             p.Name,
			Type:             pb.ParameterType_PARAMETER_TYPE_STRING,
			Tier:             pb.ParameterTier_PARAMETER_TIER_STANDARD,
			Version:          p.Version,
			Lastmodifieddate: p.LastModifiedDate.Format("2006-01-02T15:04:05.000Z"),
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

func NewConnectHandler(sm *storage.RegionStorageManager, accountID string) (string, http.Handler) {
	return ssmconnect.NewSSMServiceHandler(NewAdminHandler(sm, accountID))
}

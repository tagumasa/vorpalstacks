package ssm

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	svcerrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/utils/timeutils"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	pb "vorpalstacks/internal/pb/aws/ssm"
	ssmconnect "vorpalstacks/internal/pb/aws/ssm/ssmconnect"
	ssmstore "vorpalstacks/internal/store/aws/ssm"
)

// AdminHandler implements the SSM admin console gRPC-Web handler.
// It delegates to the shared SSMService store cache so that the same
// per-region store instances are used by both the HTTP API handlers and the
// admin console gRPC-Web handlers.
type AdminHandler struct {
	ssmconnect.UnimplementedSSMServiceHandler
	service *SSMService
}

var _ ssmconnect.SSMServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new SSM admin handler backed by the given
// service instance.
func NewAdminHandler(svc *SSMService) *AdminHandler {
	return &AdminHandler{
		service: svc,
	}
}

func (h *AdminHandler) getStoreFromHeaders(headers http.Header) (*ssmstore.Store, error) {
	region := svccommon.GetRegionFromHeader(headers)
	store, err := h.service.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}
	return store.(*ssmstore.Store), nil
}

// mapAdminError translates a store-level error to a connect error. The gRPC
// admin path loses the AWS error code when errors are passed through
// svcerrors.StoreErrorToGRPC, so this dispatcher recognises the specific
// ssmstore sentinels and maps them to their connect.Code counterparts.
func mapAdminError(err error) error {
	switch {
	case errors.Is(err, ssmstore.ErrParameterNotFound):
		return connect.NewError(connect.CodeNotFound, fmt.Errorf("parameter not found"))
	case errors.Is(err, ssmstore.ErrParameterAlreadyExists):
		return connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("parameter already exists"))
	case errors.Is(err, ssmstore.ErrInvalidParameterName),
		errors.Is(err, ssmstore.ErrInvalidParameterValue),
		errors.Is(err, ssmstore.ErrInvalidParameterType),
		errors.Is(err, ssmstore.ErrInvalidParameterVersion),
		errors.Is(err, ssmstore.ErrReservedParameterName),
		errors.Is(err, ssmstore.ErrInvalidAllowedPattern),
		errors.Is(err, ssmstore.ErrParameterPatternMismatch):
		return connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%v", err))
	case errors.Is(err, ssmstore.ErrParameterVersionNotFound),
		errors.Is(err, ssmstore.ErrParameterLabelNotFound):
		return connect.NewError(connect.CodeNotFound, fmt.Errorf("%v", err))
	}
	return svcerrors.StoreErrorToGRPC(err)
}

// DescribeParameters retrieves SSM parameters from the store, applying optional filters and pagination.
func (h *AdminHandler) DescribeParameters(ctx context.Context, req *connect.Request[pb.DescribeParametersRequest]) (*connect.Response[pb.DescribeParametersResult], error) {
	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, mapAdminError(err)
	}

	maxResults := req.Msg.Maxresults
	if maxResults <= 0 {
		maxResults = 50
	}

	var filters []ssmstore.ParameterFilter
	for _, f := range req.Msg.Filters {
		key := f.Key.String()
		if key == "" {
			continue
		}
		if !ssmstore.ValidateParameterFilterKey(key) {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid filter key: %s", key))
		}
		filters = append(filters, ssmstore.ParameterFilter{
			Key:    key,
			Option: "",
			Values: f.Values,
		})
	}

	params, nextToken, err := store.DescribeParameters(filters, maxResults, req.Msg.Nexttoken)
	if err != nil {
		return nil, mapAdminError(err)
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
	paramType := ""
	switch req.Msg.Type {
	case pb.ParameterType_PARAMETER_TYPE_STRING_LIST:
		paramType = "StringList"
	case pb.ParameterType_PARAMETER_TYPE_SECURE_STRING:
		paramType = "SecureString"
	}

	tier := ""
	switch req.Msg.Tier {
	case pb.ParameterTier_PARAMETER_TIER_ADVANCED:
		tier = "Advanced"
	case pb.ParameterTier_PARAMETER_TIER_INTELLIGENT_TIERING:
		tier = "Intelligent-Tiering"
	}

	param, err := normalisePutParameter(ParameterPutFields{
		Name:           req.Msg.Name,
		Value:          req.Msg.Value,
		Type:           paramType,
		Description:    req.Msg.Description,
		KeyID:          req.Msg.Keyid,
		AllowedPattern: req.Msg.Allowedpattern,
		DataType:       req.Msg.Datatype,
		Tier:           tier,
		Policies:       req.Msg.Policies,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("%v", err))
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

	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, mapAdminError(err)
	}

	overwrite := req.Msg.GetOverwrite()
	modifiedBy := req.Header().Get("X-Access-Key-Id")
	if modifiedBy == "" {
		modifiedBy = "vorpalstacks:admin"
	}
	version, err := h.service.putParameterWithEncryption(ctx, store, param, overwrite, modifiedBy)
	if err != nil {
		return nil, mapAdminError(err)
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
		return nil, mapAdminError(err)
	}

	return connect.NewResponse(&pb.DeleteParameterResult{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the SSM admin console.
func NewConnectHandler(svc *SSMService) (string, http.Handler) {
	return ssmconnect.NewSSMServiceHandler(NewAdminHandler(svc))
}

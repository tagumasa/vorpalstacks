package ssm

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"

	svcerrors "vorpalstacks/internal/common/errors"

	pb "vorpalstacks/internal/pb/aws/ssm"
	ssmconnect "vorpalstacks/internal/pb/aws/ssm/ssmconnect"
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
	return &AdminHandler{service: svc}
}

// DescribeParameters retrieves SSM parameters from the store, applying optional filters and pagination.
func (h *AdminHandler) DescribeParameters(ctx context.Context, req *connect.Request[pb.DescribeParametersRequest]) (*connect.Response[pb.DescribeParametersResult], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	filters, err := toStoreFilters(req.Msg.Filters)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	result, err := h.service.describeParametersCore(store, DescribeParametersInput{
		Filters:    filters,
		MaxResults: req.Msg.GetMaxresults(),
		NextToken:  req.Msg.Nexttoken,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(toSSMError(err))
	}

	var metadataList []*pb.ParameterMetadata
	for _, p := range result.Parameters {
		metadataList = append(metadataList, toPbParameterMetadataFromMeta(p))
	}

	return connect.NewResponse(&pb.DescribeParametersResult{
		Parameters: metadataList,
		Nexttoken:  result.NextToken,
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

	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	overwrite := req.Msg.GetOverwrite()
	modifiedBy := req.Header().Get("X-Access-Key-Id")
	if modifiedBy == "" {
		modifiedBy = "vorpalstacks:admin"
	}
	version, err := h.service.putParameterCore(ctx, store, param, overwrite, modifiedBy)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(toSSMError(err))
	}

	return connect.NewResponse(&pb.PutParameterResult{
		Version: proto.Int64(version),
	}), nil
}

// DeleteParameter deletes an SSM parameter via the admin console.
func (h *AdminHandler) DeleteParameter(ctx context.Context, req *connect.Request[pb.DeleteParameterRequest]) (*connect.Response[pb.DeleteParameterResult], error) {
	store, err := h.getStore(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	if err := h.service.deleteParameterCore(store, req.Msg.Name); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(toSSMError(err))
	}

	return connect.NewResponse(&pb.DeleteParameterResult{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the SSM admin console.
func NewConnectHandler(svc *SSMService) (string, http.Handler) {
	return ssmconnect.NewSSMServiceHandler(NewAdminHandler(svc))
}

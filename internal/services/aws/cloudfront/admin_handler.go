package cloudfront

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"

	svccommon "vorpalstacks/internal/common"
	svcerrors "vorpalstacks/internal/common/errors"

	pb "vorpalstacks/internal/pb/aws/cloudfront"
	cloudfrontconnect "vorpalstacks/internal/pb/aws/cloudfront/cloudfrontconnect"
	pbcommon "vorpalstacks/internal/pb/aws/common"
)

// AdminHandler implements the CloudFront admin console gRPC-Web handler.
// It delegates to the shared CloudFrontService store cache and core
// functions so that the same global store instance and validation logic
// are used by both the HTTP API handlers and the admin console gRPC-Web
// handlers.
type AdminHandler struct {
	cloudfrontconnect.UnimplementedCloudFrontServiceHandler
	service *CloudFrontService
}

var _ cloudfrontconnect.CloudFrontServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new CloudFront admin handler.
func NewAdminHandler(svc *CloudFrontService) *AdminHandler {
	return &AdminHandler{service: svc}
}

func (h *AdminHandler) getStoreFromHeaders(headers http.Header) (*cloudfrontStores, error) {
	return h.service.GetStoreForRegion("")
}

// CreateDistribution creates a new CloudFront distribution via the admin
// console. It delegates to createDistributionCore so that certificate
// checking, WAF association, and tag application follow the same code path
// as the HTTP API.
func (h *AdminHandler) CreateDistribution(ctx context.Context, req *connect.Request[pb.CreateDistributionRequest]) (*connect.Response[pb.CreateDistributionResult], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	cfg := req.Msg.GetDistributionconfig()
	if cfg == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("distributionconfig is required"))
	}

	originDomain := ""
	originID := "default-origin"
	if cfg.GetOrigins() != nil && len(cfg.GetOrigins().GetItems()) > 0 {
		firstOrigin := cfg.GetOrigins().GetItems()[0]
		originDomain = firstOrigin.GetDomainname()
		if firstOrigin.GetId() != "" {
			originID = firstOrigin.GetId()
		}
	}
	if originDomain == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("at least one origin with a domain name is required"))
	}

	callerRef := fmt.Sprintf("%d", time.Now().UnixNano())

	result, err := h.service.createDistributionFromAdmin(ctx, stores, AdminCreateDistributionInput{
		CallerReference: callerRef,
		Comment:         cfg.GetComment(),
		Enabled:         cfg.GetEnabled(),
		OriginID:        originID,
		OriginDomain:    originDomain,
		ACMRegion:       svccommon.GetRegionFromHeader(req.Header()),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateDistributionResult{
		Distribution: toPbDistribution(result.Distribution),
	}), nil
}

// ListDistributions returns all CloudFront distributions visible to the
// admin console. It delegates to listDistributionsCore.
func (h *AdminHandler) ListDistributions(ctx context.Context, req *connect.Request[pb.ListDistributionsRequest]) (*connect.Response[pb.ListDistributionsResult], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	maxItems := 100
	if req.Msg.Maxitems != nil {
		maxItems = int(*req.Msg.Maxitems)
	}

	result, err := h.service.listDistributionsCore(stores, ListDistributionsInput{
		Marker:   req.Msg.Marker,
		MaxItems: maxItems,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	items := make([]*pb.DistributionSummary, 0, len(result.Distributions))
	for _, d := range result.Distributions {
		items = append(items, toPbDistributionSummary(d))
	}

	return connect.NewResponse(&pb.ListDistributionsResult{
		Distributionlist: &pb.DistributionList{
			Quantity:    int32(len(items)),
			Items:       items,
			Istruncated: proto.Bool(result.IsTruncated),
			Nextmarker:  result.NextMarker,
			Marker:      req.Msg.Marker,
			Maxitems:    int32(maxItems),
		},
	}), nil
}

// DeleteDistribution deletes a CloudFront distribution via the admin
// console. It delegates to deleteDistributionCore so that ETag check,
// DistributionNotDisabled enforcement, WAF cleanup, tag cleanup, and
// invalidation cleanup all follow the same code path as the HTTP API.
func (h *AdminHandler) DeleteDistribution(ctx context.Context, req *connect.Request[pb.DeleteDistributionRequest]) (*connect.Response[pbcommon.Empty], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	if req.Msg.Id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("id is required"))
	}

	if err := h.service.deleteDistributionCore(ctx, stores, DeleteDistributionInput{
		Id:        req.Msg.Id,
		IfMatch:   "*",
		ACMRegion: svccommon.GetRegionFromHeader(req.Header()),
	}); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the CloudFront admin console.
func NewConnectHandler(svc *CloudFrontService) (string, http.Handler) {
	return cloudfrontconnect.NewCloudFrontServiceHandler(NewAdminHandler(svc))
}

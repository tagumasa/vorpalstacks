package wafv2

import (
	"context"
	"net/http"
	"vorpalstacks/internal/common/defaults"

	"connectrpc.com/connect"
	svcerrors "vorpalstacks/internal/common/errors"
	pb "vorpalstacks/internal/pb/aws/wafv2"
	wafv2connect "vorpalstacks/internal/pb/aws/wafv2/wafv2connect"
)

// AdminHandler implements the WAFv2 admin console gRPC-Web handler.
// It is a thin adapter: all business logic delegates to Core methods on
// WAFv2Service, and store types are confined to admin_handler_convert.go.
type AdminHandler struct {
	wafv2connect.UnimplementedWAFV2ServiceHandler
	service *WAFv2Service
}

var _ wafv2connect.WAFV2ServiceHandler = (*AdminHandler)(nil)

func NewAdminHandler(svc *WAFv2Service) *AdminHandler {
	return &AdminHandler{service: svc}
}

// ListWebACLs returns a paginated list of WebACL summaries via the admin
// console gRPC-Web interface.
func (h *AdminHandler) ListWebACLs(ctx context.Context, req *connect.Request[pb.ListWebACLsRequest]) (*connect.Response[pb.ListWebACLsResponse], error) {
	stores, err := h.service.GetStoresForRegion(defaults.GetRegionFromHeader(req.Header()))
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	scope := "REGIONAL"
	if req.Msg.Scope == pb.Scope_SCOPE_CLOUDFRONT {
		scope = "CLOUDFRONT"
	}

	result, err := h.service.listWebACLsCore(stores, ListWebACLsInput{
		Scope:      scope,
		Limit:      int(req.Msg.GetLimit()),
		NextMarker: req.Msg.Nextmarker,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	summaries := make([]*pb.WebACLSummary, 0, len(result.WebACLs))
	for _, wa := range result.WebACLs {
		summaries = append(summaries, toPbWebACLSummary(wa))
	}

	return connect.NewResponse(&pb.ListWebACLsResponse{
		Webacls:    summaries,
		Nextmarker: result.NextMarker,
	}), nil
}

// CreateWebACL creates a new WebACL via the admin console gRPC-Web
// interface.
func (h *AdminHandler) CreateWebACL(ctx context.Context, req *connect.Request[pb.CreateWebACLRequest]) (*connect.Response[pb.CreateWebACLResponse], error) {
	stores, err := h.service.GetStoresForRegion(defaults.GetRegionFromHeader(req.Header()))
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	scope := "REGIONAL"
	if req.Msg.Scope == pb.Scope_SCOPE_CLOUDFRONT {
		scope = "CLOUDFRONT"
	}

	// The console-side defaults for the required DefaultAction and
	// VisibilityConfig members are synthesised in the convert file,
	// which is the only admin-layer file allowed to reference store
	// types.
	webACL, err := h.service.createWebACLCore(stores, defaultCreateWebACLInput(req.Msg.GetName(), req.Msg.GetDescription(), scope))
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateWebACLResponse{
		Summary: toPbWebACLSummary(webACL),
	}), nil
}

// DeleteWebACL deletes an existing WebACL via the admin console gRPC-Web
// interface.
func (h *AdminHandler) DeleteWebACL(ctx context.Context, req *connect.Request[pb.DeleteWebACLRequest]) (*connect.Response[pb.DeleteWebACLResponse], error) {
	lockToken := req.Msg.GetLocktoken()
	stores, err := h.service.GetStoresForRegion(defaults.GetRegionFromHeader(req.Header()))
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	// DeleteWebACL must fail while the ACL is still associated with a
	// resource; the admin handler passes the request context so the
	// global-scope association store is included in the check.
	assocStores, err := h.service.allAssociationStoresForConnect(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	if _, err := h.service.deleteWebACLCore(stores, assocStores, req.Msg.GetId(), lockToken); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DeleteWebACLResponse{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the WAFv2
// admin console.
func NewConnectHandler(svc *WAFv2Service) (string, http.Handler) {
	return wafv2connect.NewWAFV2ServiceHandler(NewAdminHandler(svc))
}

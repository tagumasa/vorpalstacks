package wafv2

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	svcerrors "vorpalstacks/internal/common/errors"

	svccommon "vorpalstacks/internal/common"
	pb "vorpalstacks/internal/pb/aws/wafv2"
	wafv2connect "vorpalstacks/internal/pb/aws/wafv2/wafv2connect"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// AdminHandler implements the WAFv2 admin console gRPC-Web handler.
// It delegates to the shared WAFv2Service store cache so that the same
// per-region store instances are used by both the HTTP API handlers and the
// admin console gRPC-Web handlers.
type AdminHandler struct {
	wafv2connect.UnimplementedWAFV2ServiceHandler
	service *WAFv2Service
}

var _ wafv2connect.WAFV2ServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new WAFv2 admin handler backed by the given
// service instance.
func NewAdminHandler(svc *WAFv2Service) *AdminHandler {
	return &AdminHandler{
		service: svc,
	}
}

func (h *AdminHandler) getStoreFromHeaders(headers http.Header) (*wafstore.WebACLStore, error) {
	region := svccommon.GetRegionFromHeader(headers)
	return h.service.GetWebACLStoreForRegion(region)
}

// ListWebACLs returns a paginated list of WebACL summaries via the admin console gRPC-Web interface.
func (h *AdminHandler) ListWebACLs(ctx context.Context, req *connect.Request[pb.ListWebACLsRequest]) (*connect.Response[pb.ListWebACLsResponse], error) {
	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	maxItems := int(req.Msg.Limit)
	if maxItems <= 0 {
		maxItems = 100
	}

	result, err := store.List(req.Msg.Nextmarker, maxItems, "")
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	summaries := make([]*pb.WebACLSummary, 0, len(result.WebACLs))
	for _, wa := range result.WebACLs {
		summaries = append(summaries, &pb.WebACLSummary{
			Id:          wa.ID,
			Name:        wa.Name,
			Arn:         wa.ARN,
			Description: wa.Description,
			Locktoken:   wa.LockToken,
		})
	}

	return connect.NewResponse(&pb.ListWebACLsResponse{
		Webacls:    summaries,
		Nextmarker: result.NextMarker,
	}), nil
}

// CreateWebACL creates a new WebACL via the admin console gRPC-Web interface.
func (h *AdminHandler) CreateWebACL(ctx context.Context, req *connect.Request[pb.CreateWebACLRequest]) (*connect.Response[pb.CreateWebACLResponse], error) {
	if req.Msg.GetName() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("name is required"))
	}

	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	scope := "REGIONAL"
	if req.Msg.Scope == pb.Scope_SCOPE_CLOUDFRONT {
		scope = "CLOUDFRONT"
	}

	id, err := generateID()
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	webACL, err := store.Create(&wafstore.WebACL{
		ID:          id,
		Name:        req.Msg.GetName(),
		Description: req.Msg.GetDescription(),
		Scope:       scope,
		Capacity:    1500,
	})
	if err != nil {
		if wafstore.IsAlreadyExists(err) {
			return nil, connect.NewError(connect.CodeAlreadyExists, err)
		}
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateWebACLResponse{
		Summary: &pb.WebACLSummary{
			Id:          webACL.ID,
			Name:        webACL.Name,
			Arn:         webACL.ARN,
			Description: webACL.Description,
			Locktoken:   webACL.LockToken,
		},
	}), nil
}

// DeleteWebACL deletes an existing WebACL via the admin console gRPC-Web interface.
func (h *AdminHandler) DeleteWebACL(ctx context.Context, req *connect.Request[pb.DeleteWebACLRequest]) (*connect.Response[pb.DeleteWebACLResponse], error) {
	if req.Msg.GetId() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("id is required"))
	}

	lockToken := req.Msg.GetLocktoken()
	if lockToken == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("lock token is required"))
	}

	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	deleted, err := store.Delete(req.Msg.GetId(), lockToken)
	if err != nil {
		if wafstore.IsNotFound(err) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		if wafstore.IsLockTokenMismatch(err) {
			return nil, connect.NewError(connect.CodeFailedPrecondition, err)
		}
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if deleted.ARN != "" {
		if fullStores, fsErr := h.service.GetStoresForRegion(svccommon.GetRegionFromHeader(req.Header())); fsErr == nil {
			_ = fullStores.tags.Delete(deleted.ARN)
		}
	}

	return connect.NewResponse(&pb.DeleteWebACLResponse{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the WAFv2 admin console.
func NewConnectHandler(svc *WAFv2Service) (string, http.Handler) {
	return wafv2connect.NewWAFV2ServiceHandler(NewAdminHandler(svc))
}

package wafv2

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/wafv2"
	wafv2connect "vorpalstacks/internal/pb/aws/wafv2/wafv2connect"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

// AdminHandler implements the WAFv2 admin console gRPC-Web handler.
type AdminHandler struct {
	wafv2connect.UnimplementedWAFV2ServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
}

var _ wafv2connect.WAFV2ServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new WAFv2 admin handler with the given storage manager and account ID.
func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
	}
}

func (h *AdminHandler) getStoreFromHeaders(headers http.Header) (*wafstore.WebACLStore, error) {
	region := svccommon.GetRegionFromHeader(headers)
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	return wafstore.NewWebACLStore(regionStorage, h.accountId, region), nil
}

// ListWebACLs returns a paginated list of WebACL summaries via the admin console gRPC-Web interface.
func (h *AdminHandler) ListWebACLs(ctx context.Context, req *connect.Request[pb.ListWebACLsRequest]) (*connect.Response[pb.ListWebACLsResponse], error) {
	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	maxItems := int(req.Msg.Limit)
	if maxItems <= 0 {
		maxItems = 100
	}

	result, err := store.List(req.Msg.Nextmarker, maxItems)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
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
		return nil, connect.NewError(connect.CodeInvalidArgument, connect.NewError(connect.CodeInvalidArgument, nil))
	}

	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	scope := "REGIONAL"
	if req.Msg.Scope == pb.Scope_SCOPE_CLOUDFRONT {
		scope = "CLOUDFRONT"
	}

	id, err := generateID()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	webACL, err := store.Create(id, req.Msg.GetName(), req.Msg.GetDescription(), scope, 1500, nil, nil, nil)
	if err != nil {
		if wafstore.IsAlreadyExists(err) {
			return nil, connect.NewError(connect.CodeAlreadyExists, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
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
		return nil, connect.NewError(connect.CodeInvalidArgument, connect.NewError(connect.CodeInvalidArgument, nil))
	}

	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	lockToken := req.Msg.GetLocktoken()
	if lockToken == "" {
		webACL, getErr := store.Get(req.Msg.GetId())
		if getErr != nil {
			if wafstore.IsNotFound(getErr) {
				return nil, connect.NewError(connect.CodeNotFound, getErr)
			}
			return nil, connect.NewError(connect.CodeInternal, getErr)
		}
		lockToken = webACL.LockToken
	}

	if err := store.Delete(req.Msg.GetId(), lockToken); err != nil {
		if wafstore.IsNotFound(err) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		if wafstore.IsLockTokenMismatch(err) {
			return nil, connect.NewError(connect.CodeFailedPrecondition, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.DeleteWebACLResponse{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the WAFv2 admin console.
func NewConnectHandler(storageManager *storage.RegionStorageManager, accountID string) (string, http.Handler) {
	return wafv2connect.NewWAFV2ServiceHandler(NewAdminHandler(storageManager, accountID))
}

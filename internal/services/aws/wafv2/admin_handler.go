package wafv2

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/wafv2"
	wafv2connect "vorpalstacks/internal/pb/aws/wafv2/wafv2connect"
)

// AdminHandler implements the WAFv2 admin console gRPC-Web handler.
type AdminHandler struct {
	wafv2connect.UnimplementedWAFV2ServiceHandler
}

var _ wafv2connect.WAFV2ServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new WAFv2 admin handler.
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// ListWebACLs returns an empty list of WebACLs via the admin console gRPC-Web interface.
func (h *AdminHandler) ListWebACLs(ctx context.Context, req *connect.Request[pb.ListWebACLsRequest]) (*connect.Response[pb.ListWebACLsResponse], error) {
	return connect.NewResponse(&pb.ListWebACLsResponse{
		Webacls: []*pb.WebACLSummary{},
	}), nil
}

func NewConnectHandler() (string, http.Handler) {
	return wafv2connect.NewWAFV2ServiceHandler(NewAdminHandler())
}

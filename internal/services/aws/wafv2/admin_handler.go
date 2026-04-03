package wafv2

import (
	"context"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/wafv2"
	wafv2connect "vorpalstacks/internal/pb/aws/wafv2/wafv2connect"
)

type AdminHandler struct {
	wafv2connect.UnimplementedWAFV2ServiceHandler
}

var _ wafv2connect.WAFV2ServiceHandler = (*AdminHandler)(nil)

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) ListWebACLs(ctx context.Context, req *connect.Request[pb.ListWebACLsRequest]) (*connect.Response[pb.ListWebACLsResponse], error) {
	return connect.NewResponse(&pb.ListWebACLsResponse{
		Webacls: []*pb.WebACLSummary{},
	}), nil
}

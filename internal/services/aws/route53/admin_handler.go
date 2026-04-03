package route53

import (
	"context"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/route53"
	route53connect "vorpalstacks/internal/pb/aws/route53/route53connect"
)

type AdminHandler struct {
	route53connect.UnimplementedRoute53ServiceHandler
}

var _ route53connect.Route53ServiceHandler = (*AdminHandler)(nil)

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) ListHostedZones(ctx context.Context, req *connect.Request[pb.ListHostedZonesRequest]) (*connect.Response[pb.ListHostedZonesResponse], error) {
	return connect.NewResponse(&pb.ListHostedZonesResponse{
		Hostedzones: []*pb.HostedZone{},
	}), nil
}

package route53

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/route53"
	route53connect "vorpalstacks/internal/pb/aws/route53/route53connect"
)

// AdminHandler implements the Route 53 admin console gRPC-Web handler.
type AdminHandler struct {
	route53connect.UnimplementedRoute53ServiceHandler
}

var _ route53connect.Route53ServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Route 53 admin handler.
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// ListHostedZones returns all Route 53 hosted zones visible to the admin console.
func (h *AdminHandler) ListHostedZones(ctx context.Context, req *connect.Request[pb.ListHostedZonesRequest]) (*connect.Response[pb.ListHostedZonesResponse], error) {
	return connect.NewResponse(&pb.ListHostedZonesResponse{
		Hostedzones: []*pb.HostedZone{},
	}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Route53 admin console.
func NewConnectHandler() (string, http.Handler) {
	return route53connect.NewRoute53ServiceHandler(NewAdminHandler())
}

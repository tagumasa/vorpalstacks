package route53

import (
	"context"
	"net/http"
	"strconv"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/route53"
	route53connect "vorpalstacks/internal/pb/aws/route53/route53connect"
	route53store "vorpalstacks/internal/store/aws/route53"
)

// AdminHandler implements the Route 53 admin console gRPC-Web handler.
type AdminHandler struct {
	route53connect.UnimplementedRoute53ServiceHandler
	store     storage.BasicStorage
	accountId string
}

var _ route53connect.Route53ServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Route 53 admin handler with the given storage and account ID.
func NewAdminHandler(store storage.BasicStorage, accountId string) *AdminHandler {
	return &AdminHandler{
		store:     store,
		accountId: accountId,
	}
}

// ListHostedZones returns all Route 53 hosted zones visible to the admin console.
func (h *AdminHandler) ListHostedZones(ctx context.Context, req *connect.Request[pb.ListHostedZonesRequest]) (*connect.Response[pb.ListHostedZonesResponse], error) {
	zoneStore := route53store.NewHostedZoneStore(h.store, h.accountId)

	maxItems := 100
	if req.Msg.Maxitems != "" {
		mi, err := strconv.Atoi(req.Msg.Maxitems)
		if err == nil && mi > 0 {
			maxItems = mi
		}
	}

	result, err := zoneStore.List(req.Msg.Marker, maxItems)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var zones []*pb.HostedZone
	for _, z := range result.HostedZones {
		zones = append(zones, toPbHostedZone(z))
	}

	return connect.NewResponse(&pb.ListHostedZonesResponse{
		Hostedzones: zones,
		Istruncated: result.IsTruncated,
		Marker:      req.Msg.Marker,
		Nextmarker:  result.Marker,
	}), nil
}

// CreateHostedZone creates a new Route 53 hosted zone via the admin console.
func (h *AdminHandler) CreateHostedZone(ctx context.Context, req *connect.Request[pb.CreateHostedZoneRequest]) (*connect.Response[pb.CreateHostedZoneResponse], error) {
	zoneStore := route53store.NewHostedZoneStore(h.store, h.accountId)

	zone := &route53store.HostedZone{
		Name:            req.Msg.Name,
		CallerReference: req.Msg.Callerreference,
		AccountID:       h.accountId,
	}

	if req.Msg.Hostedzoneconfig != nil {
		zone.Config = &route53store.HostedZoneConfig{
			Comment:     req.Msg.Hostedzoneconfig.Comment,
			PrivateZone: req.Msg.Hostedzoneconfig.Privatezone,
		}
	}

	if req.Msg.Vpc != nil {
		zone.VPCs = []*route53store.VPC{{
			VPCID: req.Msg.Vpc.Vpcid,
		}}
		zone.Private = true
	}

	if err := zoneStore.Create(zone); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.CreateHostedZoneResponse{
		Hostedzone: toPbHostedZone(zone),
	}), nil
}

// DeleteHostedZone deletes a Route 53 hosted zone via the admin console.
func (h *AdminHandler) DeleteHostedZone(ctx context.Context, req *connect.Request[pb.DeleteHostedZoneRequest]) (*connect.Response[pb.DeleteHostedZoneResponse], error) {
	zoneStore := route53store.NewHostedZoneStore(h.store, h.accountId)

	if err := zoneStore.Delete(req.Msg.Id); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.DeleteHostedZoneResponse{}), nil
}

// toPbHostedZone converts a store HostedZone to a proto HostedZone.
func toPbHostedZone(z *route53store.HostedZone) *pb.HostedZone {
	pbZone := &pb.HostedZone{
		Id:                     z.ID,
		Name:                   z.Name,
		Callerreference:        z.CallerReference,
		Resourcerecordsetcount: int64(z.ResourceRecordSetCount),
	}
	if z.Config != nil {
		pbZone.Config = &pb.HostedZoneConfig{
			Comment:     z.Config.Comment,
			Privatezone: z.Config.PrivateZone,
		}
	}
	return pbZone
}

// NewConnectHandler creates a gRPC-Web connect handler for the Route53 admin console.
func NewConnectHandler(store storage.BasicStorage, accountID string) (string, http.Handler) {
	return route53connect.NewRoute53ServiceHandler(NewAdminHandler(store, accountID))
}

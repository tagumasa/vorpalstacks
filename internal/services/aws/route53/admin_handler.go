package route53

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"connectrpc.com/connect"
	svccommon "vorpalstacks/internal/common"
	svcerrors "vorpalstacks/internal/common/errors"

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
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid max_items: %s", req.Msg.Maxitems))
		}
		if mi > 0 {
			maxItems = mi
		}
	}

	result, err := zoneStore.List(req.Msg.Marker, maxItems)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
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
	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("name is required"))
	}

	zoneStore := route53store.NewHostedZoneStore(h.store, h.accountId)
	recordSetStore := route53store.NewRecordSetStore(h.store)

	zoneName := route53store.NormalizeZoneName(req.Msg.Name)
	nameServers := generateNameServers(4)

	privateZone := false
	if req.Msg.Vpc != nil {
		privateZone = true
	}

	comment := ""
	if req.Msg.Hostedzoneconfig != nil {
		comment = req.Msg.Hostedzoneconfig.Comment
		if req.Msg.Hostedzoneconfig.Privatezone {
			privateZone = true
		}
	}

	zone := &route53store.HostedZone{
		ID:              generateHostedZoneId(),
		Name:            zoneName,
		CallerReference: req.Msg.Callerreference,
		AccountID:       h.accountId,
		NameServers:     nameServers,
		Config:          &route53store.HostedZoneConfig{Comment: comment, PrivateZone: privateZone},
		Private:         privateZone,
		Region:          svccommon.GetRegionFromHeader(req.Header()),
	}

	if req.Msg.Vpc != nil {
		zone.VPCs = []*route53store.VPC{{
			VPCID:     req.Msg.Vpc.Vpcid,
			VPCRegion: protoVPCRegionToAWS(req.Msg.Vpc.Vpcregion),
		}}
	}

	if err := zoneStore.Create(zone); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	nsRecords := make([]*route53store.ResourceRecord, len(nameServers))
	for i, ns := range nameServers {
		nsRecords[i] = &route53store.ResourceRecord{Value: ns}
	}
	if err := recordSetStore.Create(zone.ID, &route53store.ResourceRecordSet{
		Name:            zoneName,
		Type:            "NS",
		TTL:             172800,
		ResourceRecords: nsRecords,
	}); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if err := recordSetStore.Create(zone.ID, &route53store.ResourceRecordSet{
		Name: zoneName,
		Type: "SOA",
		TTL:  900,
		ResourceRecords: []*route53store.ResourceRecord{
			{Value: fmt.Sprintf("%s %s 1 7200 900 1209600 86400", zoneName, nameServers[0])},
		},
	}); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	zone.ResourceRecordSetCount = 2
	if err := zoneStore.Update(zone); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateHostedZoneResponse{
		Hostedzone: toPbHostedZone(zone),
	}), nil
}

// DeleteHostedZone deletes a Route 53 hosted zone via the admin console.
func (h *AdminHandler) DeleteHostedZone(ctx context.Context, req *connect.Request[pb.DeleteHostedZoneRequest]) (*connect.Response[pb.DeleteHostedZoneResponse], error) {
	if req.Msg.Id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("id is required"))
	}

	zoneStore := route53store.NewHostedZoneStore(h.store, h.accountId)

	if err := zoneStore.Delete(req.Msg.Id); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
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

// protoVPCRegionToAWS converts a proto VPCRegion enum name (e.g. "V_P_C_REGION_US_EAST_1")
// to an AWS region string (e.g. "us-east-1").
func protoVPCRegionToAWS(region pb.VPCRegion) string {
	name := region.String()
	const prefix = "V_P_C_REGION_"
	if !strings.HasPrefix(name, prefix) {
		return strings.ToLower(name)
	}
	parts := strings.Split(strings.ToLower(strings.TrimPrefix(name, prefix)), "_")
	return strings.Join(parts, "-")
}

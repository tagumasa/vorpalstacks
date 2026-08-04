package route53

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
	svccommon "vorpalstacks/internal/common"
	svcerrors "vorpalstacks/internal/common/errors"

	pb "vorpalstacks/internal/pb/aws/route53"
	route53connect "vorpalstacks/internal/pb/aws/route53/route53connect"
	route53store "vorpalstacks/internal/store/aws/route53"
)

// AdminHandler implements the Route 53 admin console gRPC-Web handler.
type AdminHandler struct {
	route53connect.UnimplementedRoute53ServiceHandler
	service *Route53Service
}

var _ route53connect.Route53ServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Route 53 admin handler.
func NewAdminHandler(svc *Route53Service) *AdminHandler {
	return &AdminHandler{service: svc}
}

func (h *AdminHandler) getStoreFromHeaders(headers http.Header) (*route53store.Route53Stores, error) {
	region := svccommon.GetRegionFromHeader(headers)
	return h.service.GetStoreForRegion(region)
}

// ListHostedZones returns all Route 53 hosted zones visible to the admin console.
func (h *AdminHandler) ListHostedZones(ctx context.Context, req *connect.Request[pb.ListHostedZonesRequest]) (*connect.Response[pb.ListHostedZonesResponse], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

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

	result, err := stores.HostedZones().List(req.Msg.Marker, maxItems)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	var zones []*pb.HostedZone
	for _, z := range result.HostedZones {
		zones = append(zones, toPbHostedZone(z))
	}

	return connect.NewResponse(&pb.ListHostedZonesResponse{
		Hostedzones: zones,
		Istruncated: proto.Bool(result.IsTruncated),
		Marker:      req.Msg.Marker,
		Nextmarker:  result.Marker,
	}), nil
}

// CreateHostedZone creates a new Route 53 hosted zone via the admin console.
func (h *AdminHandler) CreateHostedZone(ctx context.Context, req *connect.Request[pb.CreateHostedZoneRequest]) (*connect.Response[pb.CreateHostedZoneResponse], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("name is required"))
	}

	zoneName := route53store.NormalizeZoneName(req.Msg.Name)
	nameServers := generateNameServers(4)

	privateZone := false
	if req.Msg.Vpc != nil {
		privateZone = true
	}

	comment := ""
	if req.Msg.Hostedzoneconfig != nil {
		comment = req.Msg.Hostedzoneconfig.Comment
		if req.Msg.Hostedzoneconfig.GetPrivatezone() {
			privateZone = true
		}
	}

	zone := &route53store.HostedZone{
		ID:              generateHostedZoneId(),
		Name:            zoneName,
		CallerReference: req.Msg.Callerreference,
		AccountID:       h.service.accountID,
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

	if err := stores.HostedZones().Create(zone); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	nsRecords := make([]*route53store.ResourceRecord, len(nameServers))
	for i, ns := range nameServers {
		nsRecords[i] = &route53store.ResourceRecord{Value: ns}
	}
	if err := stores.RecordSets().Create(zone.ID, &route53store.ResourceRecordSet{
		Name:            zoneName,
		Type:            "NS",
		TTL:             172800,
		ResourceRecords: nsRecords,
	}); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if err := stores.RecordSets().Create(zone.ID, &route53store.ResourceRecordSet{
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
	if err := stores.HostedZones().Update(zone); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateHostedZoneResponse{
		Hostedzone: toPbHostedZone(zone),
	}), nil
}

// DeleteHostedZone deletes a Route 53 hosted zone via the admin console.
func (h *AdminHandler) DeleteHostedZone(ctx context.Context, req *connect.Request[pb.DeleteHostedZoneRequest]) (*connect.Response[pb.DeleteHostedZoneResponse], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if req.Msg.Id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("id is required"))
	}

	// H_admin: Reject deletion if the zone contains user-created record
	// sets (anything other than the default NS and SOA records), matching
	// the HTTP API behaviour in hosted_zone_operations.go.
	recordSets, err := stores.RecordSets().List(req.Msg.Id)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}
	for _, rs := range recordSets {
		if rs.Type != "NS" && rs.Type != "SOA" {
			return nil, connect.NewError(connect.CodeFailedPrecondition,
				fmt.Errorf("hosted zone must be empty before it can be deleted"))
		}
	}

	stores.Tags().Raw().Delete("hostedzone/" + req.Msg.Id)

	if err := stores.HostedZones().Delete(req.Msg.Id); err != nil {
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
		Resourcerecordsetcount: proto.Int64(int64(z.ResourceRecordSetCount)),
	}
	if z.Config != nil {
		pbZone.Config = &pb.HostedZoneConfig{
			Comment:     z.Config.Comment,
			Privatezone: proto.Bool(z.Config.PrivateZone),
		}
	}
	// Output VPCs so the admin console can display associations.
	if len(z.VPCs) > 0 {
		pbZone.Vpcs = make([]*pb.VPC, len(z.VPCs))
		for i, vpc := range z.VPCs {
			pbZone.Vpcs[i] = &pb.VPC{
				Vpcid:     vpc.VPCID,
				Vpcregion: awsVPCRegionToProto(vpc.VPCRegion),
			}
		}
	}
	return pbZone
}

// NewConnectHandler creates a gRPC-Web connect handler for the Route53 admin console.
func NewConnectHandler(svc *Route53Service) (string, http.Handler) {
	return route53connect.NewRoute53ServiceHandler(NewAdminHandler(svc))
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

// awsVPCRegionToProto converts an AWS region string (e.g. "us-east-1")
// to a proto VPCRegion enum value.
func awsVPCRegionToProto(region string) pb.VPCRegion {
	name := "V_P_C_REGION_" + strings.ToUpper(strings.ReplaceAll(region, "-", "_"))
	if val, ok := pb.VPCRegion_value[name]; ok {
		return pb.VPCRegion(val)
	}
	return pb.VPCRegion(0)
}

package route53

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
	svccommon "vorpalstacks/internal/common"
	svcerrors "vorpalstacks/internal/common/errors"

	pb "vorpalstacks/internal/pb/aws/route53"
	route53connect "vorpalstacks/internal/pb/aws/route53/route53connect"
)

// AdminHandler implements the Route 53 admin console gRPC-Web handler.
// It delegates to shared Core methods so that validation and persistence
// follow a single code path shared with the HTTP API handlers.
type AdminHandler struct {
	route53connect.UnimplementedRoute53ServiceHandler
	service *Route53Service
}

var _ route53connect.Route53ServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Route 53 admin handler.
func NewAdminHandler(svc *Route53Service) *AdminHandler {
	return &AdminHandler{service: svc}
}

func (h *AdminHandler) getStoreFromHeaders(headers http.Header) (*route53Stores, error) {
	region := svccommon.GetRegionFromHeader(headers)
	return h.service.GetStoreForRegion(region)
}

// ListHostedZones returns all Route 53 hosted zones visible to the admin console.
func (h *AdminHandler) ListHostedZones(ctx context.Context, req *connect.Request[pb.ListHostedZonesRequest]) (*connect.Response[pb.ListHostedZonesResponse], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	maxItems := 100
	if req.Msg.Maxitems != "" {
		mi, err := strconv.Atoi(req.Msg.Maxitems)
		if err != nil {
			return nil, svcerrors.AWSErrorToGRPC(fmt.Errorf("invalid max_items: %s", req.Msg.Maxitems))
		}
		if mi > 0 {
			maxItems = mi
		}
	}

	result, err := h.service.listHostedZonesCore(stores, ListHostedZonesInput{
		Marker:   req.Msg.Marker,
		MaxItems: maxItems,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	var zones []*pb.HostedZone
	for _, z := range result.HostedZones {
		zones = append(zones, toPbHostedZone(z))
	}

	return connect.NewResponse(&pb.ListHostedZonesResponse{
		Hostedzones: zones,
		Istruncated: proto.Bool(result.IsTruncated),
		Marker:      result.Marker,
		Nextmarker:  result.NextMarker,
	}), nil
}

// CreateHostedZone creates a new Route 53 hosted zone via the admin console.
func (h *AdminHandler) CreateHostedZone(ctx context.Context, req *connect.Request[pb.CreateHostedZoneRequest]) (*connect.Response[pb.CreateHostedZoneResponse], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	input := CreateHostedZoneInput{
		Name:            req.Msg.Name,
		CallerReference: req.Msg.Callerreference,
	}

	if req.Msg.Hostedzoneconfig != nil {
		input.Comment = req.Msg.Hostedzoneconfig.Comment
		if req.Msg.Hostedzoneconfig.GetPrivatezone() {
			input.PrivateZone = true
		}
	}

	if req.Msg.Vpc != nil {
		input.VPCID = req.Msg.Vpc.Vpcid
		input.VPCRegion = protoVPCRegionToAWS(req.Msg.Vpc.Vpcregion)
	}

	result, err := h.service.createHostedZoneCore(stores, input)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateHostedZoneResponse{
		Hostedzone: toPbHostedZone(result.Zone),
	}), nil
}

// DeleteHostedZone deletes a Route 53 hosted zone via the admin console.
func (h *AdminHandler) DeleteHostedZone(ctx context.Context, req *connect.Request[pb.DeleteHostedZoneRequest]) (*connect.Response[pb.DeleteHostedZoneResponse], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	if req.Msg.Id == "" {
		return nil, svcerrors.AWSErrorToGRPC(fmt.Errorf("id is required"))
	}

	_, err = h.service.deleteHostedZoneCore(stores, DeleteHostedZoneInput{
		Id: req.Msg.Id,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DeleteHostedZoneResponse{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Route53 admin console.
func NewConnectHandler(svc *Route53Service) (string, http.Handler) {
	return route53connect.NewRoute53ServiceHandler(NewAdminHandler(svc))
}

package sesv2

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
	svcerrors "vorpalstacks/internal/common/errors"

	svccommon "vorpalstacks/internal/common"
	pb "vorpalstacks/internal/pb/aws/sesv2"
	sesv2connect "vorpalstacks/internal/pb/aws/sesv2/sesv2connect"
	storecommon "vorpalstacks/internal/store/aws/common"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

// AdminHandler implements the SESv2 gRPC-Web admin console handler.
type AdminHandler struct {
	sesv2connect.UnimplementedSESv2ServiceHandler
	service *SESv2Service
}

var _ sesv2connect.SESv2ServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new SESv2 admin console handler.
func NewAdminHandler(svc *SESv2Service) *AdminHandler {
	return &AdminHandler{service: svc}
}

func (h *AdminHandler) getStoreFromHeaders(headers http.Header) (sesv2store.SESv2StoreInterface, error) {
	region := svccommon.GetRegionFromHeader(headers)
	return h.service.GetStoreForRegion(region)
}

// ListEmailIdentities returns a paginated list of email identities in the
// requested region.
func (h *AdminHandler) ListEmailIdentities(ctx context.Context, req *connect.Request[pb.ListEmailIdentitiesRequest]) (*connect.Response[pb.ListEmailIdentitiesResponse], error) {
	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	limit := int(req.Msg.Pagesize)
	if limit <= 0 {
		limit = 100
	}

	opts := storecommon.ListOptions{
		MaxItems: limit,
		Marker:   req.Msg.Nexttoken,
	}

	result, err := store.ListEmailIdentities(opts)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	var identities []*pb.IdentityInfo
	for _, identity := range result.Items {
		info := &pb.IdentityInfo{
			Identityname:       identity.Identity,
			Identitytype:       pb.IdentityType_IDENTITY_TYPE_EMAIL_ADDRESS,
			Sendingenabled:     proto.Bool(identity.VerifiedForSending),
			Verificationstatus: verificationStatusToProto(identity.DkimAttributes),
		}
		if identity.IdentityType == "DOMAIN" {
			info.Identitytype = pb.IdentityType_IDENTITY_TYPE_DOMAIN
		}
		identities = append(identities, info)
	}

	return connect.NewResponse(&pb.ListEmailIdentitiesResponse{
		Emailidentities: identities,
		Nexttoken:       result.NextMarker,
	}), nil
}

// CreateEmailIdentity creates a new email identity via the admin console.
func (h *AdminHandler) CreateEmailIdentity(ctx context.Context, req *connect.Request[pb.CreateEmailIdentityRequest]) (*connect.Response[pb.CreateEmailIdentityResponse], error) {
	if req.Msg.Emailidentity == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("EmailIdentity is required"))
	}
	// Validate format to match the HTTP API behaviour (M2/L7 alignment).
	if !isValidIdentityFormat(req.Msg.Emailidentity) {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("EmailIdentity format is invalid"))
	}

	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	identity := sesv2store.NewEmailIdentity(req.Msg.Emailidentity)
	result, err := store.CreateEmailIdentity(identity)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	resp := &pb.CreateEmailIdentityResponse{
		Identitytype:             pb.IdentityType_IDENTITY_TYPE_EMAIL_ADDRESS,
		Verifiedforsendingstatus: proto.Bool(result.VerifiedForSending),
	}
	if result.IdentityType == "DOMAIN" {
		resp.Identitytype = pb.IdentityType_IDENTITY_TYPE_DOMAIN
	}
	if result.DkimAttributes != nil {
		resp.Dkimattributes = &pb.DkimAttributes{
			Signingenabled: proto.Bool(result.DkimAttributes.SigningEnabled),
			Tokens:         result.DkimAttributes.Tokens,
			Status:         pb.DkimStatus_DKIM_STATUS_SUCCESS,
		}
	}
	return connect.NewResponse(resp), nil
}

// DeleteEmailIdentity deletes an email identity via the admin console.
func (h *AdminHandler) DeleteEmailIdentity(ctx context.Context, req *connect.Request[pb.DeleteEmailIdentityRequest]) (*connect.Response[pb.DeleteEmailIdentityResponse], error) {
	if req.Msg.Emailidentity == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("EmailIdentity is required"))
	}

	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if err := store.DeleteEmailIdentity(req.Msg.Emailidentity); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DeleteEmailIdentityResponse{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Sesv2 admin console.
func NewConnectHandler(svc *SESv2Service) (string, http.Handler) {
	return sesv2connect.NewSESv2ServiceHandler(NewAdminHandler(svc))
}

func verificationStatusToProto(dkim *sesv2store.DkimAttributes) pb.VerificationStatus {
	if dkim == nil {
		return pb.VerificationStatus_VERIFICATION_STATUS_PENDING
	}
	switch dkim.Status {
	case "SUCCESS":
		return pb.VerificationStatus_VERIFICATION_STATUS_SUCCESS
	case "PENDING":
		return pb.VerificationStatus_VERIFICATION_STATUS_PENDING
	case "FAILED":
		return pb.VerificationStatus_VERIFICATION_STATUS_FAILED
	case "TEMPORARY_FAILURE":
		return pb.VerificationStatus_VERIFICATION_STATUS_TEMPORARY_FAILURE
	case "NOT_STARTED":
		return pb.VerificationStatus_VERIFICATION_STATUS_NOT_STARTED
	default:
		return pb.VerificationStatus_VERIFICATION_STATUS_PENDING
	}
}

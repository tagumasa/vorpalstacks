package sesv2

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"connectrpc.com/connect"
	svcerrors "vorpalstacks/internal/common/errors"

	svccommon "vorpalstacks/internal/common"
	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/sesv2"
	sesv2connect "vorpalstacks/internal/pb/aws/sesv2/sesv2connect"
	storecommon "vorpalstacks/internal/store/aws/common"
	sesv2store "vorpalstacks/internal/store/aws/sesv2"
)

// AdminHandler implements the SESv2 gRPC-Web admin console handler. It exposes
// list operations for email identities for the Flutter management UI.
type AdminHandler struct {
	sesv2connect.UnimplementedSESv2ServiceHandler
	storageManager *storage.RegionStorageManager
	accountId      string
	stores         sync.Map // region → *sesv2store.SESv2Store
}

var _ sesv2connect.SESv2ServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new SESv2 admin console handler.
func NewAdminHandler(storageManager *storage.RegionStorageManager, accountId string) *AdminHandler {
	return &AdminHandler{
		storageManager: storageManager,
		accountId:      accountId,
	}
}

func (h *AdminHandler) getSESv2StoreFromHeader(header http.Header) (*sesv2store.SESv2Store, error) {
	region := svccommon.GetRegionFromHeader(header)
	if cached, ok := h.stores.Load(region); ok {
		return cached.(*sesv2store.SESv2Store), nil
	}
	regionStorage, err := h.storageManager.GetStorage(region)
	if err != nil {
		return nil, err
	}
	s := sesv2store.NewSESv2Store(regionStorage, h.accountId, region)
	h.stores.Store(region, s)
	return s, nil
}

// ListEmailIdentities returns a paginated list of email identities in the
// requested region.
func (h *AdminHandler) ListEmailIdentities(ctx context.Context, req *connect.Request[pb.ListEmailIdentitiesRequest]) (*connect.Response[pb.ListEmailIdentitiesResponse], error) {
	store, err := h.getSESv2StoreFromHeader(req.Header())
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
			Sendingenabled:     identity.VerifiedForSending,
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

	store, err := h.getSESv2StoreFromHeader(req.Header())
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
		Verifiedforsendingstatus: result.VerifiedForSending,
	}
	if result.IdentityType == "DOMAIN" {
		resp.Identitytype = pb.IdentityType_IDENTITY_TYPE_DOMAIN
	}
	if result.DkimAttributes != nil {
		resp.Dkimattributes = &pb.DkimAttributes{
			Signingenabled: result.DkimAttributes.SigningEnabled,
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

	store, err := h.getSESv2StoreFromHeader(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if err := store.DeleteEmailIdentity(req.Msg.Emailidentity); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DeleteEmailIdentityResponse{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Sesv2 admin console.
func NewConnectHandler(sm *storage.RegionStorageManager, accountID string) (string, http.Handler) {
	return sesv2connect.NewSESv2ServiceHandler(NewAdminHandler(sm, accountID))
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

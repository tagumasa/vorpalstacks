package sesv2

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	svcerrors "vorpalstacks/internal/common/errors"
	pb "vorpalstacks/internal/pb/aws/sesv2"
	sesv2connect "vorpalstacks/internal/pb/aws/sesv2/sesv2connect"
	"vorpalstacks/internal/utils/aws/types"
)

// AdminHandler implements the SESv2 gRPC-Web admin console handler. It is
// a thin adapter: every operation delegates to the service-layer Core
// methods, ensuring identical validation to the HTTP API. No store
// packages are imported directly (AGENTS.md rule #29).
type AdminHandler struct {
	sesv2connect.UnimplementedSESv2ServiceHandler
	service *SESv2Service
}

var _ sesv2connect.SESv2ServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new SESv2 admin console handler backed by the
// given service instance.
func NewAdminHandler(svc *SESv2Service) *AdminHandler {
	return &AdminHandler{service: svc}
}

// ListEmailIdentities returns a paginated list of email identities in the
// requested region.
func (h *AdminHandler) ListEmailIdentities(ctx context.Context, req *connect.Request[pb.ListEmailIdentitiesRequest]) (*connect.Response[pb.ListEmailIdentitiesResponse], error) {
	store, err := h.getSESv2Store(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	result, err := h.service.listEmailIdentitiesCore(store, ListEmailIdentitiesInput{
		NextToken: req.Msg.GetNexttoken(),
		MaxItems:  int(req.Msg.GetPagesize()),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.ListEmailIdentitiesResponse{
		Emailidentities: toPbIdentityInfos(result.Identities),
		Nexttoken:       result.NextToken,
	}), nil
}

// CreateEmailIdentity creates a new email identity via the admin console.
// Tags, ConfigurationSetName, and DkimSigningAttributes are fully
// supported because the handler delegates to createEmailIdentityCore,
// which performs the same validation and processing as the HTTP API.
func (h *AdminHandler) CreateEmailIdentity(ctx context.Context, req *connect.Request[pb.CreateEmailIdentityRequest]) (*connect.Response[pb.CreateEmailIdentityResponse], error) {
	store, err := h.getSESv2Store(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	var parsedTags []types.Tag
	for _, t := range req.Msg.GetTags() {
		parsedTags = append(parsedTags, types.Tag{Key: t.GetKey(), Value: t.GetValue()})
	}

	input := CreateEmailIdentityInput{
		EmailIdentity:        req.Msg.GetEmailidentity(),
		ConfigurationSetName: req.Msg.GetConfigurationsetname(),
		Tags:                 parsedTags,
	}

	// Map proto DkimSigningAttributes into the Core input so
	// BYODKIM data is not silently discarded. Only the string-typed
	// fields (DomainSigningSelector, DomainSigningPrivateKey) are mapped
	// because proto3 does not provide field presence for enum scalars —
	// the proto getter returns the zero-value enum (AWS_SES_AP_NORTHEAST_3
	// for origin, RSA_1024_BIT for key length) even when the client did
	// not explicitly set them. Setting those defaults would override the
	// identity's correct origin/length with spurious values.
	//
	// applyDkimSigningAttributes handles missing origin by defaulting to
	// EXTERNAL, and missing key length by preserving the existing value.
	if dkimAttrs := req.Msg.GetDkimsigningattributes(); dkimAttrs != nil {
		input.DkimSigningAttrs = map[string]interface{}{
			"DomainSigningSelector":   dkimAttrs.GetDomainsigningselector(),
			"DomainSigningPrivateKey": dkimAttrs.GetDomainsigningprivatekey(),
		}
		input.DkimSigningProvided = true
	}

	result, err := h.service.createEmailIdentityCore(ctx, store, input)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(toPbCreateEmailIdentityResponse(result)), nil
}

// DeleteEmailIdentity deletes an email identity via the admin console.
func (h *AdminHandler) DeleteEmailIdentity(ctx context.Context, req *connect.Request[pb.DeleteEmailIdentityRequest]) (*connect.Response[pb.DeleteEmailIdentityResponse], error) {
	store, err := h.getSESv2Store(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if err := h.service.deleteEmailIdentityCore(store, req.Msg.GetEmailidentity()); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DeleteEmailIdentityResponse{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the SESv2 admin console.
func NewConnectHandler(svc *SESv2Service) (string, http.Handler) {
	return sesv2connect.NewSESv2ServiceHandler(NewAdminHandler(svc))
}

package cognitoidentity

import (
	"context"
	"errors"
	"net/http"

	"connectrpc.com/connect"
	svccommon "vorpalstacks/internal/common"
	svcerrors "vorpalstacks/internal/common/errors"

	pb "vorpalstacks/internal/pb/aws/cognitoidentity"
	"vorpalstacks/internal/pb/aws/cognitoidentity/cognitoidentityconnect"
	"vorpalstacks/internal/pb/aws/common"
	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// AdminHandler provides Cognito Identity service administration functionality.
// It delegates to the shared CognitoIdentityService store cache so that the same
// per-region store instances are used by both the HTTP API handlers and the
// admin console gRPC-Web handlers.
type AdminHandler struct {
	cognitoidentityconnect.UnimplementedCognitoIdentityServiceHandler
	service *CognitoIdentityService
}

var _ cognitoidentityconnect.CognitoIdentityServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Cognito Identity AdminHandler.
func NewAdminHandler(svc *CognitoIdentityService) *AdminHandler {
	return &AdminHandler{service: svc}
}

func (h *AdminHandler) getStoreFromHeaders(headers http.Header) (cognitoidentitystore.CognitoIdentityStoreInterface, error) {
	region := svccommon.GetRegionFromHeader(headers)
	return h.service.GetStoreForRegion(region)
}

// ListIdentityPools lists identity pools in Cognito Identity with pagination.
func (h *AdminHandler) ListIdentityPools(ctx context.Context, req *connect.Request[pb.ListIdentityPoolsInput]) (*connect.Response[pb.ListIdentityPoolsResponse], error) {
	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	maxResults := int(req.Msg.Maxresults)
	if maxResults <= 0 {
		maxResults = 60
	}

	opts := storecommon.ListOptions{
		MaxItems: maxResults,
		Marker:   req.Msg.Nexttoken,
	}

	result, err := store.ListIdentityPools(opts)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	descriptions := make([]*pb.IdentityPoolShortDescription, 0, len(result.Items))
	for _, pool := range result.Items {
		descriptions = append(descriptions, &pb.IdentityPoolShortDescription{
			Identitypoolid:   pool.ID,
			Identitypoolname: pool.Name,
		})
	}

	return connect.NewResponse(&pb.ListIdentityPoolsResponse{
		Identitypools: descriptions,
		Nexttoken:     result.NextMarker,
	}), nil
}

// CreateIdentityPool creates a new Cognito Identity Pool via the admin console.
func (h *AdminHandler) CreateIdentityPool(ctx context.Context, req *connect.Request[pb.CreateIdentityPoolInput]) (*connect.Response[pb.IdentityPool], error) {
	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if req.Msg.GetIdentitypoolname() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("IdentityPoolName is required"))
	}

	region := svccommon.GetRegionFromHeader(req.Header())

	pool := cognitoidentitystore.NewIdentityPool(
		req.Msg.GetIdentitypoolname(),
		req.Msg.GetAllowunauthenticatedidentities(),
		region,
	)

	if req.Msg.GetAllowclassicflow() {
		pool.AllowClassicFlow = req.Msg.GetAllowclassicflow()
	}
	if req.Msg.GetDeveloperprovidername() != "" {
		pool.DeveloperProviderName = req.Msg.GetDeveloperprovidername()
	}

	for _, p := range req.Msg.GetCognitoidentityproviders() {
		pool.CognitoIdentityProviders = append(pool.CognitoIdentityProviders, cognitoidentitystore.CognitoIdentityProvider{
			ProviderName:         p.GetProvidername(),
			ClientID:             p.GetClientid(),
			ServerSideTokenCheck: p.GetServersidetokencheck(),
		})
	}

	if m := req.Msg.GetSupportedloginproviders(); len(m) > 0 {
		pool.SupportedLoginProviders = m
	}
	if s := req.Msg.GetOpenidconnectproviderarns(); len(s) > 0 {
		pool.OpenIdConnectProviderARNs = s
	}
	if s := req.Msg.GetSamlproviderarns(); len(s) > 0 {
		pool.SamlProviderARNs = s
	}
	if t := req.Msg.GetIdentitypooltags(); len(t) > 0 {
		pool.Tags = t
	}

	created, err := store.CreateIdentityPool(pool)
	if err != nil {
		if errors.Is(err, cognitoidentitystore.ErrIdentityPoolAlreadyExists) {
			return nil, connect.NewError(connect.CodeAlreadyExists, err)
		}
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	resp := &pb.IdentityPool{
		Identitypoolid:                 created.ID,
		Identitypoolname:               created.Name,
		Allowunauthenticatedidentities: created.AllowUnauthenticatedIdentities,
		Allowclassicflow:               created.AllowClassicFlow,
		Developerprovidername:          created.DeveloperProviderName,
	}
	if len(created.CognitoIdentityProviders) > 0 {
		for _, p := range created.CognitoIdentityProviders {
			resp.Cognitoidentityproviders = append(resp.Cognitoidentityproviders, &pb.CognitoIdentityProvider{
				Providername:         p.ProviderName,
				Clientid:             p.ClientID,
				Serversidetokencheck: p.ServerSideTokenCheck,
			})
		}
	}
	if len(created.SupportedLoginProviders) > 0 {
		resp.Supportedloginproviders = created.SupportedLoginProviders
	}
	if len(created.OpenIdConnectProviderARNs) > 0 {
		resp.Openidconnectproviderarns = created.OpenIdConnectProviderARNs
	}
	if len(created.SamlProviderARNs) > 0 {
		resp.Samlproviderarns = created.SamlProviderARNs
	}
	if len(created.Tags) > 0 {
		resp.Identitypooltags = created.Tags
	}

	return connect.NewResponse(resp), nil
}

// DeleteIdentityPool deletes a Cognito Identity Pool via the admin console.
func (h *AdminHandler) DeleteIdentityPool(ctx context.Context, req *connect.Request[pb.DeleteIdentityPoolInput]) (*connect.Response[common.Empty], error) {
	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if req.Msg.GetIdentitypoolid() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("IdentityPoolId is required"))
	}

	if err := store.DeleteIdentityPool(req.Msg.GetIdentitypoolid()); err != nil {
		if errors.Is(err, cognitoidentitystore.ErrIdentityPoolNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&common.Empty{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Cognito Identity admin console.
func NewConnectHandler(svc *CognitoIdentityService) (string, http.Handler) {
	return cognitoidentityconnect.NewCognitoIdentityServiceHandler(NewAdminHandler(svc))
}

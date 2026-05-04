package cognitoidentity

import (
	"context"
	"errors"
	"net/http"

	"connectrpc.com/connect"
	svcerrors "vorpalstacks/internal/common/errors"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/cognitoidentity"
	"vorpalstacks/internal/pb/aws/cognitoidentity/cognitoidentityconnect"
	"vorpalstacks/internal/pb/aws/common"
	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// AdminHandler provides Cognito Identity service administration functionality.
type AdminHandler struct {
	cognitoidentityconnect.UnimplementedCognitoIdentityServiceHandler
	store *cognitoidentitystore.CognitoIdentityStore
}

var _ cognitoidentityconnect.CognitoIdentityServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Cognito Identity AdminHandler.
func NewAdminHandler(store *cognitoidentitystore.CognitoIdentityStore) *AdminHandler {
	return &AdminHandler{store: store}
}

// ListIdentityPools lists identity pools in Cognito Identity with pagination.
func (h *AdminHandler) ListIdentityPools(ctx context.Context, req *connect.Request[pb.ListIdentityPoolsInput]) (*connect.Response[pb.ListIdentityPoolsResponse], error) {
	maxResults := int(req.Msg.Maxresults)
	if maxResults <= 0 {
		maxResults = 60
	}

	opts := storecommon.ListOptions{
		MaxItems: maxResults,
		Marker:   req.Msg.Nexttoken,
	}

	result, err := h.store.ListIdentityPools(opts)
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
	if req.Msg.GetIdentitypoolname() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("IdentityPoolName is required"))
	}

	pool := cognitoidentitystore.NewIdentityPool(
		req.Msg.GetIdentitypoolname(),
		req.Msg.GetAllowunauthenticatedidentities(),
		h.store.Region(),
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

	created, err := h.store.CreateIdentityPool(pool)
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
	if req.Msg.GetIdentitypoolid() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("IdentityPoolId is required"))
	}

	if err := h.store.DeleteIdentityPool(req.Msg.GetIdentitypoolid()); err != nil {
		if errors.Is(err, cognitoidentitystore.ErrIdentityPoolNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&common.Empty{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Cognito Identity admin console.
func NewConnectHandler(st storage.BasicStorage, accountID, region string) (string, http.Handler) {
	return cognitoidentityconnect.NewCognitoIdentityServiceHandler(NewAdminHandler(cognitoidentitystore.NewCognitoIdentityStore(st, accountID, region)))
}

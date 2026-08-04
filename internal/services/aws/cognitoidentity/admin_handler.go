package cognitoidentity

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
	svccommon "vorpalstacks/internal/common"
	svcerrors "vorpalstacks/internal/common/errors"

	pb "vorpalstacks/internal/pb/aws/cognitoidentity"
	"vorpalstacks/internal/pb/aws/cognitoidentity/cognitoidentityconnect"
	"vorpalstacks/internal/pb/aws/common"
)

// AdminHandler provides Cognito Identity service administration functionality.
// It delegates to Core methods on CognitoIdentityService, ensuring that
// validation and persistence follow a single code path shared with the HTTP API.
type AdminHandler struct {
	cognitoidentityconnect.UnimplementedCognitoIdentityServiceHandler
	service *CognitoIdentityService
}

var _ cognitoidentityconnect.CognitoIdentityServiceHandler = (*AdminHandler)(nil)

func NewAdminHandler(svc *CognitoIdentityService) *AdminHandler {
	return &AdminHandler{service: svc}
}

// ListIdentityPools lists identity pools in Cognito Identity with pagination.
func (h *AdminHandler) ListIdentityPools(ctx context.Context, req *connect.Request[pb.ListIdentityPoolsInput]) (*connect.Response[pb.ListIdentityPoolsResponse], error) {
	region := svccommon.GetRegionFromHeader(req.Header())
	store, err := h.service.GetStoreForRegion(region)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	items, nextToken, err := h.service.listIdentityPoolsShortCore(store, ListIdentityPoolsInput{
		MaxResults: int(req.Msg.Maxresults),
		NextToken:  req.Msg.Nexttoken,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	descriptions := make([]*pb.IdentityPoolShortDescription, 0, len(items))
	for _, pool := range items {
		descriptions = append(descriptions, &pb.IdentityPoolShortDescription{
			Identitypoolid:   pool.ID,
			Identitypoolname: pool.Name,
		})
	}

	resp := &pb.ListIdentityPoolsResponse{
		Identitypools: descriptions,
	}
	if nextToken != "" {
		resp.Nexttoken = nextToken
	}
	return connect.NewResponse(resp), nil
}

// CreateIdentityPool creates a new Cognito Identity Pool via the admin console.
func (h *AdminHandler) CreateIdentityPool(ctx context.Context, req *connect.Request[pb.CreateIdentityPoolInput]) (*connect.Response[pb.IdentityPool], error) {
	region := svccommon.GetRegionFromHeader(req.Header())
	store, err := h.service.GetStoreForRegion(region)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	input := CreateIdentityPoolInput{
		IdentityPoolName:               req.Msg.GetIdentitypoolname(),
		AllowUnauthenticatedIdentities: req.Msg.GetAllowunauthenticatedidentities(),
		DeveloperProviderName:          req.Msg.GetDeveloperprovidername(),
		Region:                         region,
	}

	if req.Msg.GetAllowclassicflow() {
		input.AllowClassicFlow = req.Msg.GetAllowclassicflow()
		input.AllowClassicFlowProvided = true
	}

	for _, p := range req.Msg.GetCognitoidentityproviders() {
		input.CognitoIdentityProviders = append(input.CognitoIdentityProviders, ProviderOut{
			ProviderName:         p.GetProvidername(),
			ClientID:             p.GetClientid(),
			ServerSideTokenCheck: p.GetServersidetokencheck(),
		})
	}

	if m := req.Msg.GetSupportedloginproviders(); len(m) > 0 {
		input.SupportedLoginProviders = m
	}
	if s := req.Msg.GetOpenidconnectproviderarns(); len(s) > 0 {
		input.OpenIdConnectProviderARNs = s
	}
	if s := req.Msg.GetSamlproviderarns(); len(s) > 0 {
		input.SamlProviderARNs = s
	}
	if t := req.Msg.GetIdentitypooltags(); len(t) > 0 {
		input.Tags = t
		input.TagsProvided = true
	}

	result, err := h.service.createIdentityPoolCore(store, input)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(poolOutToProto(result)), nil
}

// DeleteIdentityPool deletes a Cognito Identity Pool via the admin console.
func (h *AdminHandler) DeleteIdentityPool(ctx context.Context, req *connect.Request[pb.DeleteIdentityPoolInput]) (*connect.Response[common.Empty], error) {
	region := svccommon.GetRegionFromHeader(req.Header())
	store, err := h.service.GetStoreForRegion(region)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if err := h.service.deleteIdentityPoolCore(store, req.Msg.GetIdentitypoolid()); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&common.Empty{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Cognito Identity admin console.
func NewConnectHandler(svc *CognitoIdentityService) (string, http.Handler) {
	return cognitoidentityconnect.NewCognitoIdentityServiceHandler(NewAdminHandler(svc))
}

func poolOutToProto(p *IdentityPoolOut) *pb.IdentityPool {
	resp := &pb.IdentityPool{
		Identitypoolid:                 p.ID,
		Identitypoolname:               p.Name,
		Allowunauthenticatedidentities: proto.Bool(p.AllowUnauthenticatedIdentities),
		Allowclassicflow:               proto.Bool(p.AllowClassicFlow),
		Developerprovidername:          p.DeveloperProviderName,
	}
	for _, cp := range p.CognitoIdentityProviders {
		resp.Cognitoidentityproviders = append(resp.Cognitoidentityproviders, &pb.CognitoIdentityProvider{
			Providername:         cp.ProviderName,
			Clientid:             cp.ClientID,
			Serversidetokencheck: proto.Bool(cp.ServerSideTokenCheck),
		})
	}
	if len(p.SupportedLoginProviders) > 0 {
		resp.Supportedloginproviders = p.SupportedLoginProviders
	}
	if len(p.OpenIdConnectProviderARNs) > 0 {
		resp.Openidconnectproviderarns = p.OpenIdConnectProviderARNs
	}
	if len(p.SamlProviderARNs) > 0 {
		resp.Samlproviderarns = p.SamlProviderARNs
	}
	if len(p.Tags) > 0 {
		resp.Identitypooltags = p.Tags
	}
	return resp
}

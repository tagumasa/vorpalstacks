package cognitoidentityprovider

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
	svccommon "vorpalstacks/internal/common"
	svcerrors "vorpalstacks/internal/common/errors"

	pb "vorpalstacks/internal/pb/aws/cognitoidentityprovider"
	"vorpalstacks/internal/pb/aws/cognitoidentityprovider/cognitoidentityproviderconnect"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	"vorpalstacks/internal/utils/timeutils"
)

// AdminHandler provides Cognito Identity Provider service administration functionality.
// It delegates to the shared CognitoService *Core methods so that the same
// validation and persistence code path is used by both the HTTP API handlers
// and the admin console gRPC-Web handlers.
type AdminHandler struct {
	cognitoidentityproviderconnect.UnimplementedCognitoIdentityProviderServiceHandler
	service *CognitoService
}

var _ cognitoidentityproviderconnect.CognitoIdentityProviderServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Cognito Identity Provider AdminHandler.
func NewAdminHandler(svc *CognitoService) *AdminHandler {
	return &AdminHandler{service: svc}
}

// getRegionFromHeaders extracts the AWS region from the gRPC-Web metadata.
func (h *AdminHandler) getRegionFromHeaders(headers http.Header) string {
	return svccommon.GetRegionFromHeader(headers)
}

// ListUserPools lists user pools in Cognito Identity Provider with pagination.
func (h *AdminHandler) ListUserPools(ctx context.Context, req *connect.Request[pb.ListUserPoolsRequest]) (*connect.Response[pb.ListUserPoolsResponse], error) {
	region := h.getRegionFromHeaders(req.Header())

	result, err := h.service.listUserPoolsCore(region, ListUserPoolsInput{
		MaxResults: int(req.Msg.GetMaxresults()),
		NextToken:  req.Msg.GetNexttoken(),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	descriptions := make([]*pb.UserPoolDescriptionType, 0, len(result.UserPools))
	for _, pool := range result.UserPools {
		desc := &pb.UserPoolDescriptionType{
			Id:   pool.ID,
			Name: pool.Name,
		}
		if !pool.CreationDate.IsZero() {
			desc.Creationdate = pool.CreationDate.Format(timeutils.ISO8601UTCFormat)
		}
		if !pool.LastModifiedDate.IsZero() {
			desc.Lastmodifieddate = pool.LastModifiedDate.Format(timeutils.ISO8601UTCFormat)
		}
		descriptions = append(descriptions, desc)
	}

	return connect.NewResponse(&pb.ListUserPoolsResponse{
		Userpools: descriptions,
		Nexttoken: result.NextToken,
	}), nil
}

// CreateUserPool creates a new Cognito user pool via the admin console.
func (h *AdminHandler) CreateUserPool(ctx context.Context, req *connect.Request[pb.CreateUserPoolRequest]) (*connect.Response[pb.CreateUserPoolResponse], error) {
	region := h.getRegionFromHeaders(req.Header())

	autoVerifiedAttrs := make([]string, 0, len(req.Msg.GetAutoverifiedattributes()))
	for _, attr := range req.Msg.GetAutoverifiedattributes() {
		switch attr {
		case pb.VerifiedAttributeType_VERIFIED_ATTRIBUTE_TYPE_EMAIL:
			autoVerifiedAttrs = append(autoVerifiedAttrs, "email")
		case pb.VerifiedAttributeType_VERIFIED_ATTRIBUTE_TYPE_PHONE_NUMBER:
			autoVerifiedAttrs = append(autoVerifiedAttrs, "phone_number")
		}
	}

	var pp *AdminPasswordPolicy
	if req.Msg.GetPolicies() != nil && req.Msg.GetPolicies().GetPasswordpolicy() != nil {
		protoPP := req.Msg.GetPolicies().GetPasswordpolicy()
		pp = &AdminPasswordPolicy{
			MinimumLength:                 int(protoPP.GetMinimumlength()),
			RequireUppercase:              protoPP.GetRequireuppercase(),
			RequireLowercase:              protoPP.GetRequirelowercase(),
			RequireNumbers:                protoPP.GetRequirenumbers(),
			RequireSymbols:                protoPP.GetRequiresymbols(),
			TemporaryPasswordValidityDays: int(protoPP.GetTemporarypasswordvaliditydays()),
			PasswordHistorySize:           int(protoPP.GetPasswordhistorysize()),
		}
	}

	created, err := h.service.createUserPoolFromAdmin(AdminCreateUserPoolInput{
		PoolName:          req.Msg.GetPoolname(),
		Region:            region,
		AutoVerifiedAttrs: autoVerifiedAttrs,
		PasswordPolicy:    pp,
		Tags:              req.Msg.GetUserpooltags(),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateUserPoolResponse{
		Userpool: userPoolToProto(created),
	}), nil
}

// DeleteUserPool deletes a Cognito user pool via the admin console.
func (h *AdminHandler) DeleteUserPool(ctx context.Context, req *connect.Request[pb.DeleteUserPoolRequest]) (*connect.Response[pbcommon.Empty], error) {
	region := h.getRegionFromHeaders(req.Header())

	if err := h.service.deleteUserPoolCore(region, req.Msg.GetUserpoolid()); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// DescribeUserPool retrieves the full configuration of a Cognito user pool.
func (h *AdminHandler) DescribeUserPool(ctx context.Context, req *connect.Request[pb.DescribeUserPoolRequest]) (*connect.Response[pb.DescribeUserPoolResponse], error) {
	region := h.getRegionFromHeaders(req.Header())

	pool, err := h.service.describeUserPoolCore(region, req.Msg.GetUserpoolid())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DescribeUserPoolResponse{
		Userpool: userPoolToProto(pool),
	}), nil
}

// ListUserPoolClients lists the user pool clients for a given pool.
func (h *AdminHandler) ListUserPoolClients(ctx context.Context, req *connect.Request[pb.ListUserPoolClientsRequest]) (*connect.Response[pb.ListUserPoolClientsResponse], error) {
	region := h.getRegionFromHeaders(req.Header())

	result, err := h.service.listUserPoolClientsCore(region, ListUserPoolClientsInput{
		UserPoolID: req.Msg.GetUserpoolid(),
		MaxResults: int(req.Msg.GetMaxresults()),
		NextToken:  req.Msg.GetNexttoken(),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	clients := make([]*pb.UserPoolClientDescription, 0, len(result.Clients))
	for _, c := range result.Clients {
		clients = append(clients, &pb.UserPoolClientDescription{
			Clientid:   c.ClientID,
			Userpoolid: c.UserPoolID,
			Clientname: c.ClientName,
		})
	}

	return connect.NewResponse(&pb.ListUserPoolClientsResponse{
		Userpoolclients: clients,
		Nexttoken:       result.NextToken,
	}), nil
}

// DescribeUserPoolClient retrieves the full configuration of a user pool client.
func (h *AdminHandler) DescribeUserPoolClient(ctx context.Context, req *connect.Request[pb.DescribeUserPoolClientRequest]) (*connect.Response[pb.DescribeUserPoolClientResponse], error) {
	region := h.getRegionFromHeaders(req.Header())

	client, err := h.service.describeUserPoolClientCore(region, req.Msg.GetUserpoolid(), req.Msg.GetClientid())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DescribeUserPoolClientResponse{
		Userpoolclient: userPoolClientToProto(client),
	}), nil
}

// DeleteUserPoolClient deletes a user pool client.
func (h *AdminHandler) DeleteUserPoolClient(ctx context.Context, req *connect.Request[pb.DeleteUserPoolClientRequest]) (*connect.Response[pbcommon.Empty], error) {
	region := h.getRegionFromHeaders(req.Header())

	if err := h.service.deleteUserPoolClientCore(region, req.Msg.GetUserpoolid(), req.Msg.GetClientid()); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// ListGroups lists the groups in a Cognito user pool.
func (h *AdminHandler) ListGroups(ctx context.Context, req *connect.Request[pb.ListGroupsRequest]) (*connect.Response[pb.ListGroupsResponse], error) {
	region := h.getRegionFromHeaders(req.Header())

	result, err := h.service.listGroupsCore(region, ListGroupsInput{
		UserPoolID: req.Msg.GetUserpoolid(),
		MaxResults: int(req.Msg.GetLimit()),
		NextToken:  req.Msg.GetNexttoken(),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	groups := make([]*pb.GroupType, 0, len(result.Groups))
	for _, g := range result.Groups {
		groups = append(groups, groupToProto(g))
	}

	return connect.NewResponse(&pb.ListGroupsResponse{
		Groups:    groups,
		Nexttoken: result.NextToken,
	}), nil
}

// CreateGroup creates a new Cognito group.
func (h *AdminHandler) CreateGroup(ctx context.Context, req *connect.Request[pb.CreateGroupRequest]) (*connect.Response[pb.CreateGroupResponse], error) {
	region := h.getRegionFromHeaders(req.Header())

	input := CreateGroupInput{
		UserPoolID:  req.Msg.GetUserpoolid(),
		GroupName:   req.Msg.GetGroupname(),
		Description: req.Msg.GetDescription(),
		RoleArn:     req.Msg.GetRolearn(),
	}
	if req.Msg.Precedence != nil {
		p := int(*req.Msg.Precedence)
		input.Precedence = &p
	}

	group, err := h.service.createGroupFromInputCore(region, input)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateGroupResponse{
		Group: groupToProto(group),
	}), nil
}

// DeleteGroup deletes a Cognito group.
func (h *AdminHandler) DeleteGroup(ctx context.Context, req *connect.Request[pb.DeleteGroupRequest]) (*connect.Response[pbcommon.Empty], error) {
	region := h.getRegionFromHeaders(req.Header())

	if err := h.service.deleteGroupCore(region, req.Msg.GetUserpoolid(), req.Msg.GetGroupname()); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// ListTagsForResource lists all tags assigned to a Cognito resource.
func (h *AdminHandler) ListTagsForResource(ctx context.Context, req *connect.Request[pb.ListTagsForResourceRequest]) (*connect.Response[pb.ListTagsForResourceResponse], error) {
	region := h.getRegionFromHeaders(req.Header())

	tags, err := h.service.listTagsForResourceCore(region, req.Msg.GetResourcearn())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.ListTagsForResourceResponse{
		Tags: tags,
	}), nil
}

// TagResource adds or overwrites tags on a Cognito resource.
func (h *AdminHandler) TagResource(ctx context.Context, req *connect.Request[pb.TagResourceRequest]) (*connect.Response[pb.TagResourceResponse], error) {
	region := h.getRegionFromHeaders(req.Header())

	if err := h.service.tagResourceCore(region, req.Msg.GetResourcearn(), req.Msg.GetTags()); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.TagResourceResponse{}), nil
}

// UntagResource removes tags from a Cognito resource.
func (h *AdminHandler) UntagResource(ctx context.Context, req *connect.Request[pb.UntagResourceRequest]) (*connect.Response[pb.UntagResourceResponse], error) {
	region := h.getRegionFromHeaders(req.Header())

	if err := h.service.untagResourceCore(region, req.Msg.GetResourcearn(), req.Msg.GetTagkeys()); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.UntagResourceResponse{}), nil
}

// ListUsers lists users in a Cognito user pool.
func (h *AdminHandler) ListUsers(ctx context.Context, req *connect.Request[pb.ListUsersRequest]) (*connect.Response[pb.ListUsersResponse], error) {
	region := h.getRegionFromHeaders(req.Header())

	result, err := h.service.listUsersCore(region, ListUsersInput{
		UserPoolID: req.Msg.GetUserpoolid(),
		MaxResults: int(req.Msg.GetLimit()),
		NextToken:  req.Msg.GetPaginationtoken(),
		Filter:     req.Msg.GetFilter(),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	users := make([]*pb.UserType, 0, len(result.Users))
	for _, u := range result.Users {
		users = append(users, userToProto(u))
	}

	return connect.NewResponse(&pb.ListUsersResponse{
		Users:           users,
		Paginationtoken: result.NextToken,
	}), nil
}

// AdminGetUser retrieves a user by username.
func (h *AdminHandler) AdminGetUser(ctx context.Context, req *connect.Request[pb.AdminGetUserRequest]) (*connect.Response[pb.AdminGetUserResponse], error) {
	region := h.getRegionFromHeaders(req.Header())

	user, err := h.service.adminGetUserCore(region, req.Msg.GetUserpoolid(), req.Msg.GetUsername())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.AdminGetUserResponse{
		Username:             user.Username,
		Userstatus:           userStatusToProto(user.UserStatus),
		Enabled:              proto.Bool(user.Enabled),
		Usercreatedate:       user.CreatedDate.Format(timeutils.ISO8601UTCFormat),
		Userlastmodifieddate: user.LastModifiedDate.Format(timeutils.ISO8601UTCFormat),
	}), nil
}

// AdminDeleteUser deletes a user from a Cognito user pool.
func (h *AdminHandler) AdminDeleteUser(ctx context.Context, req *connect.Request[pb.AdminDeleteUserRequest]) (*connect.Response[pbcommon.Empty], error) {
	region := h.getRegionFromHeaders(req.Header())

	if err := h.service.adminDeleteUserCore(region, req.Msg.GetUserpoolid(), req.Msg.GetUsername()); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// AdminEnableUser enables a user.
func (h *AdminHandler) AdminEnableUser(ctx context.Context, req *connect.Request[pb.AdminEnableUserRequest]) (*connect.Response[pb.AdminEnableUserResponse], error) {
	region := h.getRegionFromHeaders(req.Header())

	if err := h.service.adminEnableUserCore(region, req.Msg.GetUserpoolid(), req.Msg.GetUsername()); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.AdminEnableUserResponse{}), nil
}

// AdminDisableUser disables a user.
func (h *AdminHandler) AdminDisableUser(ctx context.Context, req *connect.Request[pb.AdminDisableUserRequest]) (*connect.Response[pb.AdminDisableUserResponse], error) {
	region := h.getRegionFromHeaders(req.Header())

	if err := h.service.adminDisableUserCore(region, req.Msg.GetUserpoolid(), req.Msg.GetUsername()); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.AdminDisableUserResponse{}), nil
}

// ListIdentityProviders lists identity providers in a user pool.
func (h *AdminHandler) ListIdentityProviders(ctx context.Context, req *connect.Request[pb.ListIdentityProvidersRequest]) (*connect.Response[pb.ListIdentityProvidersResponse], error) {
	region := h.getRegionFromHeaders(req.Header())

	result, err := h.service.listIdentityProvidersCore(region, ListIdentityProvidersInput{
		UserPoolID: req.Msg.GetUserpoolid(),
		MaxResults: int(req.Msg.GetMaxresults()),
		NextToken:  req.Msg.GetNexttoken(),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	providers := make([]*pb.ProviderDescription, 0, len(result.Providers))
	for _, ip := range result.Providers {
		providers = append(providers, providerDescriptionToProto(ip))
	}

	return connect.NewResponse(&pb.ListIdentityProvidersResponse{
		Providers: providers,
		Nexttoken: result.NextToken,
	}), nil
}

// DescribeIdentityProvider retrieves an identity provider by name.
func (h *AdminHandler) DescribeIdentityProvider(ctx context.Context, req *connect.Request[pb.DescribeIdentityProviderRequest]) (*connect.Response[pb.DescribeIdentityProviderResponse], error) {
	region := h.getRegionFromHeaders(req.Header())

	ip, err := h.service.describeIdentityProviderCore(region, req.Msg.GetUserpoolid(), req.Msg.GetProvidername())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.DescribeIdentityProviderResponse{
		Identityprovider: identityProviderToProto(ip),
	}), nil
}

// CreateIdentityProvider creates a new identity provider.
func (h *AdminHandler) CreateIdentityProvider(ctx context.Context, req *connect.Request[pb.CreateIdentityProviderRequest]) (*connect.Response[pb.CreateIdentityProviderResponse], error) {
	region := h.getRegionFromHeaders(req.Header())

	ip, err := h.service.createIdentityProviderFromInputCore(region, CreateIdentityProviderInput{
		UserPoolID:       req.Msg.GetUserpoolid(),
		ProviderName:     req.Msg.GetProvidername(),
		ProviderType:     identityProviderTypeFromProto(req.Msg.GetProvidertype()),
		ProviderDetails:  req.Msg.GetProviderdetails(),
		AttributeMapping: req.Msg.GetAttributemapping(),
		IdpIdentifiers:   req.Msg.GetIdpidentifiers(),
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateIdentityProviderResponse{
		Identityprovider: identityProviderToProto(ip),
	}), nil
}

// DeleteIdentityProvider deletes an identity provider.
func (h *AdminHandler) DeleteIdentityProvider(ctx context.Context, req *connect.Request[pb.DeleteIdentityProviderRequest]) (*connect.Response[pbcommon.Empty], error) {
	region := h.getRegionFromHeaders(req.Header())

	if err := h.service.deleteIdentityProviderCore(region, req.Msg.GetUserpoolid(), req.Msg.GetProvidername()); err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Cognito Identity Provider admin console.
func NewConnectHandler(svc *CognitoService) (string, http.Handler) {
	return cognitoidentityproviderconnect.NewCognitoIdentityProviderServiceHandler(NewAdminHandler(svc))
}

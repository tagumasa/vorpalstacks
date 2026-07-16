package cognitoidentityprovider

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	svccommon "vorpalstacks/internal/common"
	svcerrors "vorpalstacks/internal/common/errors"

	pb "vorpalstacks/internal/pb/aws/cognitoidentityprovider"
	"vorpalstacks/internal/pb/aws/cognitoidentityprovider/cognitoidentityproviderconnect"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
	storecommon "vorpalstacks/internal/store/aws/common"
	"vorpalstacks/internal/utils/timeutils"
)

// AdminHandler provides Cognito Identity Provider service administration functionality.
// It delegates to the shared CognitoService store cache so that the same
// per-region store instances are used by both the HTTP API handlers and the
// admin console gRPC-Web handlers.
type AdminHandler struct {
	cognitoidentityproviderconnect.UnimplementedCognitoIdentityProviderServiceHandler
	service *CognitoService
}

var _ cognitoidentityproviderconnect.CognitoIdentityProviderServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Cognito Identity Provider AdminHandler.
func NewAdminHandler(svc *CognitoService) *AdminHandler {
	return &AdminHandler{service: svc}
}

func (h *AdminHandler) getStoreFromHeaders(headers http.Header) (cognitostore.CognitoStoreInterface, error) {
	region := svccommon.GetRegionFromHeader(headers)
	return h.service.GetStoreForRegion(region)
}

// ListUserPools lists user pools in Cognito Identity Provider with pagination.
func (h *AdminHandler) ListUserPools(ctx context.Context, req *connect.Request[pb.ListUserPoolsRequest]) (*connect.Response[pb.ListUserPoolsResponse], error) {
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

	result, err := store.ListUserPoolsPaginated(opts)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	descriptions := make([]*pb.UserPoolDescriptionType, 0, len(result.Items))
	for _, pool := range result.Items {
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
		Nexttoken: result.NextMarker,
	}), nil
}

// CreateUserPool creates a new Cognito user pool via the admin console.
func (h *AdminHandler) CreateUserPool(ctx context.Context, req *connect.Request[pb.CreateUserPoolRequest]) (*connect.Response[pb.CreateUserPoolResponse], error) {
	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if req.Msg.GetPoolname() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("PoolName is required"))
	}

	region := svccommon.GetRegionFromHeader(req.Header())
	pool := cognitostore.NewUserPool(req.Msg.GetPoolname(), region)

	result, err := store.CreateUserPool(pool)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateUserPoolResponse{
		Userpool: &pb.UserPoolType{
			Id:   result.ID,
			Name: result.Name,
		},
	}), nil
}

// DeleteUserPool deletes a Cognito user pool via the admin console.
func (h *AdminHandler) DeleteUserPool(ctx context.Context, req *connect.Request[pb.DeleteUserPoolRequest]) (*connect.Response[pbcommon.Empty], error) {
	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if req.Msg.GetUserpoolid() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("UserPoolId is required"))
	}

	if err := store.DeleteUserPool(req.Msg.GetUserpoolid()); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Cognito Identity Provider admin console.
func NewConnectHandler(svc *CognitoService) (string, http.Handler) {
	return cognitoidentityproviderconnect.NewCognitoIdentityProviderServiceHandler(NewAdminHandler(svc))
}

package cognitoidentityprovider

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/cognitoidentityprovider"
	"vorpalstacks/internal/pb/aws/cognitoidentityprovider/cognitoidentityproviderconnect"
	pbcommon "vorpalstacks/internal/pb/aws/common"
	cognitostore "vorpalstacks/internal/store/aws/cognitoidentityprovider"
)

// AdminHandler provides Cognito Identity Provider service administration functionality.
type AdminHandler struct {
	cognitoidentityproviderconnect.UnimplementedCognitoIdentityProviderServiceHandler
	store *cognitostore.CognitoStore
}

var _ cognitoidentityproviderconnect.CognitoIdentityProviderServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Cognito Identity Provider AdminHandler.
func NewAdminHandler(store *cognitostore.CognitoStore) *AdminHandler {
	return &AdminHandler{store: store}
}

// ListUserPools lists user pools in Cognito Identity Provider.
func (h *AdminHandler) ListUserPools(ctx context.Context, req *connect.Request[pb.ListUserPoolsRequest]) (*connect.Response[pb.ListUserPoolsResponse], error) {
	pools, err := h.store.ListUserPools()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	maxResults := int(req.Msg.Maxresults)
	if maxResults <= 0 {
		maxResults = 60
	}

	descriptions := make([]*pb.UserPoolDescriptionType, 0, len(pools))
	for i, pool := range pools {
		if i >= maxResults {
			break
		}
		desc := &pb.UserPoolDescriptionType{
			Id:   pool.ID,
			Name: pool.Name,
		}
		if !pool.CreationDate.IsZero() {
			desc.Creationdate = pool.CreationDate.Format("2006-01-02T15:04:05.000Z")
		}
		if !pool.LastModifiedDate.IsZero() {
			desc.Lastmodifieddate = pool.LastModifiedDate.Format("2006-01-02T15:04:05.000Z")
		}
		descriptions = append(descriptions, desc)
	}

	return connect.NewResponse(&pb.ListUserPoolsResponse{
		Userpools: descriptions,
	}), nil
}

// CreateUserPool creates a new Cognito user pool via the admin console.
func (h *AdminHandler) CreateUserPool(ctx context.Context, req *connect.Request[pb.CreateUserPoolRequest]) (*connect.Response[pb.CreateUserPoolResponse], error) {
	if req.Msg.GetPoolname() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("PoolName is required"))
	}

	pool := &cognitostore.UserPool{
		Name: req.Msg.GetPoolname(),
	}

	result, err := h.store.CreateUserPool(pool)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
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
	if req.Msg.GetUserpoolid() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("UserPoolId is required"))
	}

	if err := h.store.DeleteUserPool(req.Msg.GetUserpoolid()); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pbcommon.Empty{}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Cognito Identity Provider admin console.
func NewConnectHandler(st storage.BasicStorage, accountID, region string) (string, http.Handler) {
	return cognitoidentityproviderconnect.NewCognitoIdentityProviderServiceHandler(NewAdminHandler(cognitostore.NewCognitoStore(st, accountID, region)))
}

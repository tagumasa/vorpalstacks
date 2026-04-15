package cognitoidentity

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/cognitoidentity"
	"vorpalstacks/internal/pb/aws/cognitoidentity/cognitoidentityconnect"
	cognitoidentitystore "vorpalstacks/internal/store/aws/cognitoidentity"
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

// ListIdentityPools lists identity pools in Cognito Identity.
func (h *AdminHandler) ListIdentityPools(ctx context.Context, req *connect.Request[pb.ListIdentityPoolsInput]) (*connect.Response[pb.ListIdentityPoolsResponse], error) {
	pools, err := h.store.ListIdentityPools()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	maxResults := int(req.Msg.Maxresults)
	if maxResults <= 0 {
		maxResults = 60
	}

	descriptions := make([]*pb.IdentityPoolShortDescription, 0, len(pools))
	for i, pool := range pools {
		if i >= maxResults {
			break
		}
		descriptions = append(descriptions, &pb.IdentityPoolShortDescription{
			Identitypoolid:   pool.ID,
			Identitypoolname: pool.Name,
		})
	}

	return connect.NewResponse(&pb.ListIdentityPoolsResponse{
		Identitypools: descriptions,
	}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Cognito Identity admin console.
func NewConnectHandler(st storage.BasicStorage, accountID, region string) (string, http.Handler) {
	return cognitoidentityconnect.NewCognitoIdentityServiceHandler(NewAdminHandler(cognitoidentitystore.NewCognitoIdentityStore(st, accountID, region)))
}

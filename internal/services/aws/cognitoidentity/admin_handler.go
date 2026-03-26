package cognitoidentity

import (
	"context"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/cognito_identity"
	cognitoidentityconnect "vorpalstacks/internal/pb/aws/cognito_identity/cognito_identityconnect"
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

package cognitoidentityprovider

import (
	"context"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/cognitoidentityprovider"
	"vorpalstacks/internal/pb/aws/cognitoidentityprovider/cognitoidentityproviderconnect"
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

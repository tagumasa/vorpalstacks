package apigateway

import (
	"context"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/apigateway"
	apigatewayconnect "vorpalstacks/internal/pb/aws/apigateway/apigatewayconnect"
)

type AdminHandler struct {
	apigatewayconnect.UnimplementedAPIGatewayServiceHandler
}

var _ apigatewayconnect.APIGatewayServiceHandler = (*AdminHandler)(nil)

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) GetRestApis(ctx context.Context, req *connect.Request[pb.GetRestApisRequest]) (*connect.Response[pb.RestApis], error) {
	return connect.NewResponse(&pb.RestApis{
		Items: []*pb.RestApi{},
	}), nil
}

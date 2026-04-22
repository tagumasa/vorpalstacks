// Package apigateway implements AWS API Gateway REST API operations including
// REST APIs, resources, methods, integrations, deployments, authorisers, and
// usage plans.
package apigateway

import (
	"context"
	"net/http"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/apigateway"
	apigatewayconnect "vorpalstacks/internal/pb/aws/apigateway/apigatewayconnect"
)

// AdminHandler implements the API Gateway admin console gRPC-Web handler.
type AdminHandler struct {
	apigatewayconnect.UnimplementedAPIGatewayServiceHandler
}

var _ apigatewayconnect.APIGatewayServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new API Gateway admin handler.
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// GetRestApis lists all REST APIs registered in the API Gateway service.
func (h *AdminHandler) GetRestApis(ctx context.Context, req *connect.Request[pb.GetRestApisRequest]) (*connect.Response[pb.RestApis], error) {
	return connect.NewResponse(&pb.RestApis{
		Items: []*pb.RestApi{},
	}), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the API Gateway admin console.
func NewConnectHandler() (string, http.Handler) {
	return apigatewayconnect.NewAPIGatewayServiceHandler(NewAdminHandler())
}

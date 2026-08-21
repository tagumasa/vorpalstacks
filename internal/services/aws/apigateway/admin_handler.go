package apigateway

import (
	"net/http"
	"vorpalstacks/internal/common/defaults"

	apigatewayconnect "vorpalstacks/internal/pb/aws/apigateway/apigatewayconnect"
)

// AdminHandler implements the gRPC admin console handlers for API Gateway.
type AdminHandler struct {
	apigatewayconnect.UnimplementedAPIGatewayServiceHandler
	service *APIGatewayService
}

var _ apigatewayconnect.APIGatewayServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new API Gateway admin handler.
func NewAdminHandler(svc *APIGatewayService) *AdminHandler {
	return &AdminHandler{service: svc}
}

func (h *AdminHandler) getStores(headers http.Header) (*apiGatewayStores, error) {
	region := defaults.GetRegionFromHeader(headers)
	return h.service.GetStoreForRegion(region)
}

// NewConnectHandler returns the connect RPC path and handler for API Gateway admin.
func NewConnectHandler(svc *APIGatewayService) (string, http.Handler) {
	return apigatewayconnect.NewAPIGatewayServiceHandler(NewAdminHandler(svc))
}

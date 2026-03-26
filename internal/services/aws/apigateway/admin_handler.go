package apigateway

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/apigateway"
	apigatewayconnect "vorpalstacks/internal/pb/aws/apigateway/apigatewayconnect"
	"vorpalstacks/internal/pb/aws/common"
)

// AdminHandler provides API Gateway service administration functionality.
// It implements the APIGatewayServiceHandler interface for gRPC-Web communication.
type AdminHandler struct{}

var _ apigatewayconnect.APIGatewayServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new API Gateway AdminHandler.
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// GetRestApis gets all REST APIs in API Gateway.
// REST APIs are the collection of HTTP resources and methods that you create in API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetRestApis(ctx context.Context, req *connect.Request[pb.GetRestApisRequest]) (*connect.Response[pb.RestApis], error) {
	return connect.NewResponse(&pb.RestApis{
		Items: []*pb.RestApi{},
	}), nil
}

// CreateApiKey creates an API key for API Gateway.
// API keys are used to control access to your APIs.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateApiKey(ctx context.Context, req *connect.Request[pb.CreateApiKeyRequest]) (*connect.Response[pb.ApiKey], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.CreateApiKey is not implemented"))
}

// CreateAuthorizer creates an authoriser for API Gateway.
// Authorisers control access to your APIs using various authentication methods.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateAuthorizer(ctx context.Context, req *connect.Request[pb.CreateAuthorizerRequest]) (*connect.Response[pb.Authorizer], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.CreateAuthorizer is not implemented"))
}

// CreateBasePathMapping creates a base path mapping for API Gateway.
// Base path mappings link a custom domain to an API and its stage.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateBasePathMapping(ctx context.Context, req *connect.Request[pb.CreateBasePathMappingRequest]) (*connect.Response[pb.BasePathMapping], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.CreateBasePathMapping is not implemented"))
}

// CreateDeployment creates a deployment for API Gateway.
// Deployments make your API accessible at a particular stage.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateDeployment(ctx context.Context, req *connect.Request[pb.CreateDeploymentRequest]) (*connect.Response[pb.Deployment], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.CreateDeployment is not implemented"))
}

// CreateDocumentationPart creates a documentation part for API Gateway.
// Documentation parts define the documentation for your API.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateDocumentationPart(ctx context.Context, req *connect.Request[pb.CreateDocumentationPartRequest]) (*connect.Response[pb.DocumentationPart], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.CreateDocumentationPart is not implemented"))
}

// CreateDocumentationVersion creates a documentation version for API Gateway.
// Documentation versions allow you to manage different versions of your API documentation.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateDocumentationVersion(ctx context.Context, req *connect.Request[pb.CreateDocumentationVersionRequest]) (*connect.Response[pb.DocumentationVersion], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.CreateDocumentationVersion is not implemented"))
}

// CreateDomainName creates a domain name for API Gateway.
// Domain names allow you to access your API using a custom domain.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateDomainName(ctx context.Context, req *connect.Request[pb.CreateDomainNameRequest]) (*connect.Response[pb.DomainName], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.CreateDomainName is not implemented"))
}

// CreateDomainNameAccessAssociation creates a domain name access association for API Gateway.
// This associates an access grant with a domain name.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateDomainNameAccessAssociation(ctx context.Context, req *connect.Request[pb.CreateDomainNameAccessAssociationRequest]) (*connect.Response[pb.DomainNameAccessAssociation], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.CreateDomainNameAccessAssociation is not implemented"))
}

// CreateModel creates a model for API Gateway.
// Models define the structure of request and response payloads.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateModel(ctx context.Context, req *connect.Request[pb.CreateModelRequest]) (*connect.Response[pb.Model], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.CreateModel is not implemented"))
}

// CreateRequestValidator creates a request validator for API Gateway.
// Request validators validate incoming requests before they reach your backend.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateRequestValidator(ctx context.Context, req *connect.Request[pb.CreateRequestValidatorRequest]) (*connect.Response[pb.RequestValidator], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.CreateRequestValidator is not implemented"))
}

// CreateResource creates a resource for API Gateway.
// Resources represent the URL paths of your API.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateResource(ctx context.Context, req *connect.Request[pb.CreateResourceRequest]) (*connect.Response[pb.Resource], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.CreateResource is not implemented"))
}

// CreateRestApi creates a REST API in API Gateway.
// REST APIs are the collection of HTTP resources and methods that you create in API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateRestApi(ctx context.Context, req *connect.Request[pb.CreateRestApiRequest]) (*connect.Response[pb.RestApi], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.CreateRestApi is not implemented"))
}

// CreateStage creates a stage for API Gateway.
// Stages represent deployed versions of your API (e.g., dev, test, prod).
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateStage(ctx context.Context, req *connect.Request[pb.CreateStageRequest]) (*connect.Response[pb.Stage], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.CreateStage is not implemented"))
}

// CreateUsagePlan creates a usage plan for API Gateway.
// Usage plans control access to your APIs by setting throttling and quota limits.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateUsagePlan(ctx context.Context, req *connect.Request[pb.CreateUsagePlanRequest]) (*connect.Response[pb.UsagePlan], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.CreateUsagePlan is not implemented"))
}

// CreateUsagePlanKey creates a usage plan key for API Gateway.
// Usage plan keys associate API keys with usage plans.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateUsagePlanKey(ctx context.Context, req *connect.Request[pb.CreateUsagePlanKeyRequest]) (*connect.Response[pb.UsagePlanKey], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.CreateUsagePlanKey is not implemented"))
}

// CreateVpcLink creates a VPC link for API Gateway.
// VPC links allow API Gateway to access resources in your VPC.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateVpcLink(ctx context.Context, req *connect.Request[pb.CreateVpcLinkRequest]) (*connect.Response[pb.VpcLink], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.CreateVpcLink is not implemented"))
}

// DeleteApiKey deletes an API key for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteApiKey(ctx context.Context, req *connect.Request[pb.DeleteApiKeyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.DeleteApiKey is not implemented"))
}

// DeleteAuthorizer deletes an authoriser for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteAuthorizer(ctx context.Context, req *connect.Request[pb.DeleteAuthorizerRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.DeleteAuthorizer is not implemented"))
}

// DeleteBasePathMapping deletes a base path mapping for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteBasePathMapping(ctx context.Context, req *connect.Request[pb.DeleteBasePathMappingRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.DeleteBasePathMapping is not implemented"))
}

// DeleteClientCertificate deletes a client certificate for API Gateway.
// Client certificates are used for mutual TLS authentication with your backend.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteClientCertificate(ctx context.Context, req *connect.Request[pb.DeleteClientCertificateRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.DeleteClientCertificate is not implemented"))
}

// DeleteDeployment deletes a deployment for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteDeployment(ctx context.Context, req *connect.Request[pb.DeleteDeploymentRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.DeleteDeployment is not implemented"))
}

// DeleteDocumentationPart deletes a documentation part for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteDocumentationPart(ctx context.Context, req *connect.Request[pb.DeleteDocumentationPartRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.DeleteDocumentationPart is not implemented"))
}

// DeleteDocumentationVersion deletes a documentation version for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteDocumentationVersion(ctx context.Context, req *connect.Request[pb.DeleteDocumentationVersionRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.DeleteDocumentationVersion is not implemented"))
}

// DeleteDomainName deletes a domain name for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteDomainName(ctx context.Context, req *connect.Request[pb.DeleteDomainNameRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.DeleteDomainName is not implemented"))
}

// DeleteDomainNameAccessAssociation deletes a domain name access association for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteDomainNameAccessAssociation(ctx context.Context, req *connect.Request[pb.DeleteDomainNameAccessAssociationRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.DeleteDomainNameAccessAssociation is not implemented"))
}

// DeleteGatewayResponse deletes a gateway response for API Gateway.
// Gateway responses define the responses returned when errors occur.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteGatewayResponse(ctx context.Context, req *connect.Request[pb.DeleteGatewayResponseRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.DeleteGatewayResponse is not implemented"))
}

// DeleteIntegration deletes an integration for API Gateway.
// Integrations connect API methods to backend endpoints.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteIntegration(ctx context.Context, req *connect.Request[pb.DeleteIntegrationRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.DeleteIntegration is not implemented"))
}

// DeleteIntegrationResponse deletes an integration response for API Gateway.
// Integration responses define the responses returned by the backend.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteIntegrationResponse(ctx context.Context, req *connect.Request[pb.DeleteIntegrationResponseRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.DeleteIntegrationResponse is not implemented"))
}

// DeleteMethod deletes a method for API Gateway.
// Methods define the HTTP operations (GET, POST, etc.) that your API supports.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteMethod(ctx context.Context, req *connect.Request[pb.DeleteMethodRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.DeleteMethod is not implemented"))
}

// DeleteMethodResponse deletes a method response for API Gateway.
// Method responses define the responses that your API returns to clients.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteMethodResponse(ctx context.Context, req *connect.Request[pb.DeleteMethodResponseRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.DeleteMethodResponse is not implemented"))
}

// DeleteModel deletes a model for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteModel(ctx context.Context, req *connect.Request[pb.DeleteModelRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.DeleteModel is not implemented"))
}

// DeleteRequestValidator deletes a request validator for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteRequestValidator(ctx context.Context, req *connect.Request[pb.DeleteRequestValidatorRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.DeleteRequestValidator is not implemented"))
}

// DeleteResource deletes a resource for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteResource(ctx context.Context, req *connect.Request[pb.DeleteResourceRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.DeleteResource is not implemented"))
}

// DeleteRestApi deletes a REST API for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteRestApi(ctx context.Context, req *connect.Request[pb.DeleteRestApiRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.DeleteRestApi is not implemented"))
}

// DeleteStage deletes a stage for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteStage(ctx context.Context, req *connect.Request[pb.DeleteStageRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.DeleteStage is not implemented"))
}

// DeleteUsagePlan deletes a usage plan for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteUsagePlan(ctx context.Context, req *connect.Request[pb.DeleteUsagePlanRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.DeleteUsagePlan is not implemented"))
}

// DeleteUsagePlanKey deletes a usage plan key for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteUsagePlanKey(ctx context.Context, req *connect.Request[pb.DeleteUsagePlanKeyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.DeleteUsagePlanKey is not implemented"))
}

// DeleteVpcLink deletes a VPC link for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteVpcLink(ctx context.Context, req *connect.Request[pb.DeleteVpcLinkRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.DeleteVpcLink is not implemented"))
}

// FlushStageAuthorizersCache flushes the authorisers cache for a stage in API Gateway.
// This forces API Gateway to re-fetch authoriser data.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) FlushStageAuthorizersCache(ctx context.Context, req *connect.Request[pb.FlushStageAuthorizersCacheRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.FlushStageAuthorizersCache is not implemented"))
}

// FlushStageCache flushes the cache for a stage in API Gateway.
// This forces API Gateway to re-fetch data from the backend.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) FlushStageCache(ctx context.Context, req *connect.Request[pb.FlushStageCacheRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.FlushStageCache is not implemented"))
}

// GenerateClientCertificate generates a client certificate for API Gateway.
// Client certificates are used for mutual TLS authentication with your backend.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GenerateClientCertificate(ctx context.Context, req *connect.Request[pb.GenerateClientCertificateRequest]) (*connect.Response[pb.ClientCertificate], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GenerateClientCertificate is not implemented"))
}

// GetAccount gets the account for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetAccount(ctx context.Context, req *connect.Request[pb.GetAccountRequest]) (*connect.Response[pb.Account], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetAccount is not implemented"))
}

// GetApiKey gets an API key for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetApiKey(ctx context.Context, req *connect.Request[pb.GetApiKeyRequest]) (*connect.Response[pb.ApiKey], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetApiKey is not implemented"))
}

// GetApiKeys gets all API keys for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetApiKeys(ctx context.Context, req *connect.Request[pb.GetApiKeysRequest]) (*connect.Response[pb.ApiKeys], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetApiKeys is not implemented"))
}

// GetAuthorizer gets an authoriser for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetAuthorizer(ctx context.Context, req *connect.Request[pb.GetAuthorizerRequest]) (*connect.Response[pb.Authorizer], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetAuthorizer is not implemented"))
}

// GetAuthorizers gets all authorisers for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetAuthorizers(ctx context.Context, req *connect.Request[pb.GetAuthorizersRequest]) (*connect.Response[pb.Authorizers], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetAuthorizers is not implemented"))
}

// GetBasePathMapping gets a base path mapping for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetBasePathMapping(ctx context.Context, req *connect.Request[pb.GetBasePathMappingRequest]) (*connect.Response[pb.BasePathMapping], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetBasePathMapping is not implemented"))
}

// GetBasePathMappings gets all base path mappings for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetBasePathMappings(ctx context.Context, req *connect.Request[pb.GetBasePathMappingsRequest]) (*connect.Response[pb.BasePathMappings], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetBasePathMappings is not implemented"))
}

// GetClientCertificate gets a client certificate for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetClientCertificate(ctx context.Context, req *connect.Request[pb.GetClientCertificateRequest]) (*connect.Response[pb.ClientCertificate], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetClientCertificate is not implemented"))
}

// GetClientCertificates gets all client certificates for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetClientCertificates(ctx context.Context, req *connect.Request[pb.GetClientCertificatesRequest]) (*connect.Response[pb.ClientCertificates], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetClientCertificates is not implemented"))
}

// GetDeployment gets a deployment for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetDeployment(ctx context.Context, req *connect.Request[pb.GetDeploymentRequest]) (*connect.Response[pb.Deployment], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetDeployment is not implemented"))
}

// GetDeployments gets all deployments for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetDeployments(ctx context.Context, req *connect.Request[pb.GetDeploymentsRequest]) (*connect.Response[pb.Deployments], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetDeployments is not implemented"))
}

// GetDocumentationPart gets a documentation part for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetDocumentationPart(ctx context.Context, req *connect.Request[pb.GetDocumentationPartRequest]) (*connect.Response[pb.DocumentationPart], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetDocumentationPart is not implemented"))
}

// GetDocumentationParts gets all documentation parts for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetDocumentationParts(ctx context.Context, req *connect.Request[pb.GetDocumentationPartsRequest]) (*connect.Response[pb.DocumentationParts], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetDocumentationParts is not implemented"))
}

// GetDocumentationVersion gets a documentation version for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetDocumentationVersion(ctx context.Context, req *connect.Request[pb.GetDocumentationVersionRequest]) (*connect.Response[pb.DocumentationVersion], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetDocumentationVersion is not implemented"))
}

// GetDocumentationVersions gets all documentation versions for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetDocumentationVersions(ctx context.Context, req *connect.Request[pb.GetDocumentationVersionsRequest]) (*connect.Response[pb.DocumentationVersions], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetDocumentationVersions is not implemented"))
}

// GetDomainName gets a domain name for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetDomainName(ctx context.Context, req *connect.Request[pb.GetDomainNameRequest]) (*connect.Response[pb.DomainName], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetDomainName is not implemented"))
}

// GetDomainNameAccessAssociations gets all domain name access associations for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetDomainNameAccessAssociations(ctx context.Context, req *connect.Request[pb.GetDomainNameAccessAssociationsRequest]) (*connect.Response[pb.DomainNameAccessAssociations], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetDomainNameAccessAssociations is not implemented"))
}

// GetDomainNames gets all domain names for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetDomainNames(ctx context.Context, req *connect.Request[pb.GetDomainNamesRequest]) (*connect.Response[pb.DomainNames], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetDomainNames is not implemented"))
}

// GetExport exports an API for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetExport(ctx context.Context, req *connect.Request[pb.GetExportRequest]) (*connect.Response[pb.ExportResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetExport is not implemented"))
}

// GetGatewayResponse gets a gateway response for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetGatewayResponse(ctx context.Context, req *connect.Request[pb.GetGatewayResponseRequest]) (*connect.Response[pb.GatewayResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetGatewayResponse is not implemented"))
}

// GetGatewayResponses gets all gateway responses for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetGatewayResponses(ctx context.Context, req *connect.Request[pb.GetGatewayResponsesRequest]) (*connect.Response[pb.GatewayResponses], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetGatewayResponses is not implemented"))
}

// GetIntegration gets an integration for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetIntegration(ctx context.Context, req *connect.Request[pb.GetIntegrationRequest]) (*connect.Response[pb.Integration], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetIntegration is not implemented"))
}

// GetIntegrationResponse gets an integration response for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetIntegrationResponse(ctx context.Context, req *connect.Request[pb.GetIntegrationResponseRequest]) (*connect.Response[pb.IntegrationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetIntegrationResponse is not implemented"))
}

// GetMethod gets a method for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetMethod(ctx context.Context, req *connect.Request[pb.GetMethodRequest]) (*connect.Response[pb.Method], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetMethod is not implemented"))
}

// GetMethodResponse gets a method response for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetMethodResponse(ctx context.Context, req *connect.Request[pb.GetMethodResponseRequest]) (*connect.Response[pb.MethodResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetMethodResponse is not implemented"))
}

// GetModel gets a model for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetModel(ctx context.Context, req *connect.Request[pb.GetModelRequest]) (*connect.Response[pb.Model], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetModel is not implemented"))
}

// GetModels gets all models for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetModels(ctx context.Context, req *connect.Request[pb.GetModelsRequest]) (*connect.Response[pb.Models], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetModels is not implemented"))
}

// GetModelTemplate gets a model template for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetModelTemplate(ctx context.Context, req *connect.Request[pb.GetModelTemplateRequest]) (*connect.Response[pb.Template], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetModelTemplate is not implemented"))
}

// GetRequestValidator gets a request validator for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetRequestValidator(ctx context.Context, req *connect.Request[pb.GetRequestValidatorRequest]) (*connect.Response[pb.RequestValidator], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetRequestValidator is not implemented"))
}

// GetRequestValidators gets all request validators for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetRequestValidators(ctx context.Context, req *connect.Request[pb.GetRequestValidatorsRequest]) (*connect.Response[pb.RequestValidators], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetRequestValidators is not implemented"))
}

// GetResource gets a resource for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetResource(ctx context.Context, req *connect.Request[pb.GetResourceRequest]) (*connect.Response[pb.Resource], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetResource is not implemented"))
}

// GetResources gets all resources for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetResources(ctx context.Context, req *connect.Request[pb.GetResourcesRequest]) (*connect.Response[pb.Resources], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetResources is not implemented"))
}

// GetRestApi gets a REST API for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetRestApi(ctx context.Context, req *connect.Request[pb.GetRestApiRequest]) (*connect.Response[pb.RestApi], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetRestApi is not implemented"))
}

// GetSdk gets an SDK for an API in API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetSdk(ctx context.Context, req *connect.Request[pb.GetSdkRequest]) (*connect.Response[pb.SdkResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetSdk is not implemented"))
}

// GetSdkType gets an SDK type for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetSdkType(ctx context.Context, req *connect.Request[pb.GetSdkTypeRequest]) (*connect.Response[pb.SdkType], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetSdkType is not implemented"))
}

// GetSdkTypes gets all SDK types for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetSdkTypes(ctx context.Context, req *connect.Request[pb.GetSdkTypesRequest]) (*connect.Response[pb.SdkTypes], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetSdkTypes is not implemented"))
}

// GetStage gets a stage for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetStage(ctx context.Context, req *connect.Request[pb.GetStageRequest]) (*connect.Response[pb.Stage], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetStage is not implemented"))
}

// GetStages gets all stages for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetStages(ctx context.Context, req *connect.Request[pb.GetStagesRequest]) (*connect.Response[pb.Stages], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetStages is not implemented"))
}

// GetTags gets tags for an API Gateway resource.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetTags(ctx context.Context, req *connect.Request[pb.GetTagsRequest]) (*connect.Response[pb.Tags], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetTags is not implemented"))
}

// GetUsage gets usage data for an API Gateway usage plan.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetUsage(ctx context.Context, req *connect.Request[pb.GetUsageRequest]) (*connect.Response[pb.Usage], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetUsage is not implemented"))
}

// GetUsagePlan gets a usage plan for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetUsagePlan(ctx context.Context, req *connect.Request[pb.GetUsagePlanRequest]) (*connect.Response[pb.UsagePlan], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetUsagePlan is not implemented"))
}

// GetUsagePlanKey gets a usage plan key for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetUsagePlanKey(ctx context.Context, req *connect.Request[pb.GetUsagePlanKeyRequest]) (*connect.Response[pb.UsagePlanKey], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetUsagePlanKey is not implemented"))
}

// GetUsagePlanKeys gets all usage plan keys for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetUsagePlanKeys(ctx context.Context, req *connect.Request[pb.GetUsagePlanKeysRequest]) (*connect.Response[pb.UsagePlanKeys], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetUsagePlanKeys is not implemented"))
}

// GetUsagePlans gets all usage plans for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetUsagePlans(ctx context.Context, req *connect.Request[pb.GetUsagePlansRequest]) (*connect.Response[pb.UsagePlans], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetUsagePlans is not implemented"))
}

// GetVpcLink gets a VPC link for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetVpcLink(ctx context.Context, req *connect.Request[pb.GetVpcLinkRequest]) (*connect.Response[pb.VpcLink], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetVpcLink is not implemented"))
}

// GetVpcLinks gets all VPC links for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetVpcLinks(ctx context.Context, req *connect.Request[pb.GetVpcLinksRequest]) (*connect.Response[pb.VpcLinks], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.GetVpcLinks is not implemented"))
}

// ImportApiKeys imports API keys for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ImportApiKeys(ctx context.Context, req *connect.Request[pb.ImportApiKeysRequest]) (*connect.Response[pb.ApiKeyIds], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.ImportApiKeys is not implemented"))
}

// ImportDocumentationParts imports documentation parts for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ImportDocumentationParts(ctx context.Context, req *connect.Request[pb.ImportDocumentationPartsRequest]) (*connect.Response[pb.DocumentationPartIds], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.ImportDocumentationParts is not implemented"))
}

// ImportRestApi imports a REST API for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ImportRestApi(ctx context.Context, req *connect.Request[pb.ImportRestApiRequest]) (*connect.Response[pb.RestApi], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.ImportRestApi is not implemented"))
}

// PutGatewayResponse creates or updates a gateway response for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutGatewayResponse(ctx context.Context, req *connect.Request[pb.PutGatewayResponseRequest]) (*connect.Response[pb.GatewayResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.PutGatewayResponse is not implemented"))
}

// PutIntegration creates or updates an integration for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutIntegration(ctx context.Context, req *connect.Request[pb.PutIntegrationRequest]) (*connect.Response[pb.Integration], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.PutIntegration is not implemented"))
}

// PutIntegrationResponse creates or updates an integration response for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutIntegrationResponse(ctx context.Context, req *connect.Request[pb.PutIntegrationResponseRequest]) (*connect.Response[pb.IntegrationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.PutIntegrationResponse is not implemented"))
}

// PutMethod creates or updates a method for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutMethod(ctx context.Context, req *connect.Request[pb.PutMethodRequest]) (*connect.Response[pb.Method], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.PutMethod is not implemented"))
}

// PutMethodResponse creates or updates a method response for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutMethodResponse(ctx context.Context, req *connect.Request[pb.PutMethodResponseRequest]) (*connect.Response[pb.MethodResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.PutMethodResponse is not implemented"))
}

// PutRestApi creates or updates a REST API for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutRestApi(ctx context.Context, req *connect.Request[pb.PutRestApiRequest]) (*connect.Response[pb.RestApi], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.PutRestApi is not implemented"))
}

// RejectDomainNameAccessAssociation rejects a domain name access association for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) RejectDomainNameAccessAssociation(ctx context.Context, req *connect.Request[pb.RejectDomainNameAccessAssociationRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.RejectDomainNameAccessAssociation is not implemented"))
}

// TagResource adds tags to an API Gateway resource.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TagResource(ctx context.Context, req *connect.Request[pb.TagResourceRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.TagResource is not implemented"))
}

// TestInvokeAuthorizer tests an authoriser for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TestInvokeAuthorizer(ctx context.Context, req *connect.Request[pb.TestInvokeAuthorizerRequest]) (*connect.Response[pb.TestInvokeAuthorizerResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.TestInvokeAuthorizer is not implemented"))
}

// TestInvokeMethod tests a method for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TestInvokeMethod(ctx context.Context, req *connect.Request[pb.TestInvokeMethodRequest]) (*connect.Response[pb.TestInvokeMethodResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.TestInvokeMethod is not implemented"))
}

// UntagResource removes tags from an API Gateway resource.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UntagResource(ctx context.Context, req *connect.Request[pb.UntagResourceRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.UntagResource is not implemented"))
}

// UpdateAccount updates the account for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateAccount(ctx context.Context, req *connect.Request[pb.UpdateAccountRequest]) (*connect.Response[pb.Account], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.UpdateAccount is not implemented"))
}

// UpdateApiKey updates an API key for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateApiKey(ctx context.Context, req *connect.Request[pb.UpdateApiKeyRequest]) (*connect.Response[pb.ApiKey], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.UpdateApiKey is not implemented"))
}

// UpdateAuthorizer updates an authoriser for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateAuthorizer(ctx context.Context, req *connect.Request[pb.UpdateAuthorizerRequest]) (*connect.Response[pb.Authorizer], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.UpdateAuthorizer is not implemented"))
}

// UpdateBasePathMapping updates a base path mapping for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateBasePathMapping(ctx context.Context, req *connect.Request[pb.UpdateBasePathMappingRequest]) (*connect.Response[pb.BasePathMapping], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.UpdateBasePathMapping is not implemented"))
}

// UpdateClientCertificate updates a client certificate for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateClientCertificate(ctx context.Context, req *connect.Request[pb.UpdateClientCertificateRequest]) (*connect.Response[pb.ClientCertificate], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.UpdateClientCertificate is not implemented"))
}

// UpdateDeployment updates a deployment for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateDeployment(ctx context.Context, req *connect.Request[pb.UpdateDeploymentRequest]) (*connect.Response[pb.Deployment], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.UpdateDeployment is not implemented"))
}

// UpdateDocumentationPart updates a documentation part for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateDocumentationPart(ctx context.Context, req *connect.Request[pb.UpdateDocumentationPartRequest]) (*connect.Response[pb.DocumentationPart], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.UpdateDocumentationPart is not implemented"))
}

// UpdateDocumentationVersion updates a documentation version for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateDocumentationVersion(ctx context.Context, req *connect.Request[pb.UpdateDocumentationVersionRequest]) (*connect.Response[pb.DocumentationVersion], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.UpdateDocumentationVersion is not implemented"))
}

// UpdateDomainName updates a domain name for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateDomainName(ctx context.Context, req *connect.Request[pb.UpdateDomainNameRequest]) (*connect.Response[pb.DomainName], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.UpdateDomainName is not implemented"))
}

// UpdateGatewayResponse updates a gateway response for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateGatewayResponse(ctx context.Context, req *connect.Request[pb.UpdateGatewayResponseRequest]) (*connect.Response[pb.GatewayResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.UpdateGatewayResponse is not implemented"))
}

// UpdateIntegration updates an integration for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateIntegration(ctx context.Context, req *connect.Request[pb.UpdateIntegrationRequest]) (*connect.Response[pb.Integration], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.UpdateIntegration is not implemented"))
}

// UpdateIntegrationResponse updates an integration response for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateIntegrationResponse(ctx context.Context, req *connect.Request[pb.UpdateIntegrationResponseRequest]) (*connect.Response[pb.IntegrationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.UpdateIntegrationResponse is not implemented"))
}

// UpdateMethod updates a method for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateMethod(ctx context.Context, req *connect.Request[pb.UpdateMethodRequest]) (*connect.Response[pb.Method], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.UpdateMethod is not implemented"))
}

// UpdateMethodResponse updates a method response for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateMethodResponse(ctx context.Context, req *connect.Request[pb.UpdateMethodResponseRequest]) (*connect.Response[pb.MethodResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.UpdateMethodResponse is not implemented"))
}

// UpdateModel updates a model for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateModel(ctx context.Context, req *connect.Request[pb.UpdateModelRequest]) (*connect.Response[pb.Model], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.UpdateModel is not implemented"))
}

// UpdateRequestValidator updates a request validator for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateRequestValidator(ctx context.Context, req *connect.Request[pb.UpdateRequestValidatorRequest]) (*connect.Response[pb.RequestValidator], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.UpdateRequestValidator is not implemented"))
}

// UpdateResource updates a resource for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateResource(ctx context.Context, req *connect.Request[pb.UpdateResourceRequest]) (*connect.Response[pb.Resource], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.UpdateResource is not implemented"))
}

// UpdateRestApi updates a REST API for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateRestApi(ctx context.Context, req *connect.Request[pb.UpdateRestApiRequest]) (*connect.Response[pb.RestApi], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.UpdateRestApi is not implemented"))
}

// UpdateStage updates a stage for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateStage(ctx context.Context, req *connect.Request[pb.UpdateStageRequest]) (*connect.Response[pb.Stage], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.UpdateStage is not implemented"))
}

// UpdateUsage updates usage data for an API Gateway usage plan.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateUsage(ctx context.Context, req *connect.Request[pb.UpdateUsageRequest]) (*connect.Response[pb.Usage], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.UpdateUsage is not implemented"))
}

// UpdateUsagePlan updates a usage plan for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateUsagePlan(ctx context.Context, req *connect.Request[pb.UpdateUsagePlanRequest]) (*connect.Response[pb.UsagePlan], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.UpdateUsagePlan is not implemented"))
}

// UpdateVpcLink updates a VPC link for API Gateway.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateVpcLink(ctx context.Context, req *connect.Request[pb.UpdateVpcLinkRequest]) (*connect.Response[pb.VpcLink], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("apigateway.APIGatewayService.UpdateVpcLink is not implemented"))
}

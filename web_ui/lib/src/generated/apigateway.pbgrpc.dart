// This is a generated file - do not edit.
//
// Generated from apigateway.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:grpc/service_api.dart' as $grpc;
import 'package:protobuf/protobuf.dart' as $pb;

import 'apigateway.pb.dart' as $0;
import 'common.pb.dart' as $1;

export 'apigateway.pb.dart';

/// APIGatewayService provides apigateway API operations.
@$pb.GrpcServiceName('apigateway.APIGatewayService')
class APIGatewayServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  APIGatewayServiceClient(super.channel, {super.options, super.interceptors});

  /// Create an ApiKey resource.
  /// HTTP: POST /apikeys
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ApiKey> createApiKey(
    $0.CreateApiKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createApiKey, request, options: options);
  }

  /// Adds a new Authorizer resource to an existing RestApi resource.
  /// HTTP: POST /restapis/{restApiId}/authorizers
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Authorizer> createAuthorizer(
    $0.CreateAuthorizerRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createAuthorizer, request, options: options);
  }

  /// Creates a new BasePathMapping resource.
  /// HTTP: POST /domainnames/{domainName}/basepathmappings
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.BasePathMapping> createBasePathMapping(
    $0.CreateBasePathMappingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createBasePathMapping, request, options: options);
  }

  /// Creates a Deployment resource, which makes a specified RestApi callable over the internet.
  /// HTTP: POST /restapis/{restApiId}/deployments
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Deployment> createDeployment(
    $0.CreateDeploymentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createDeployment, request, options: options);
  }

  /// Creates a documentation part.
  /// HTTP: POST /restapis/{restApiId}/documentation/parts
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DocumentationPart> createDocumentationPart(
    $0.CreateDocumentationPartRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createDocumentationPart, request,
        options: options);
  }

  /// Creates a documentation version
  /// HTTP: POST /restapis/{restApiId}/documentation/versions
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DocumentationVersion> createDocumentationVersion(
    $0.CreateDocumentationVersionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createDocumentationVersion, request,
        options: options);
  }

  /// Creates a new domain name.
  /// HTTP: POST /domainnames
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DomainName> createDomainName(
    $0.CreateDomainNameRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createDomainName, request, options: options);
  }

  /// Creates a domain name access association resource between an access association source and a private custom domain name.
  /// HTTP: POST /domainnameaccessassociations
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DomainNameAccessAssociation>
      createDomainNameAccessAssociation(
    $0.CreateDomainNameAccessAssociationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createDomainNameAccessAssociation, request,
        options: options);
  }

  /// Adds a new Model resource to an existing RestApi resource.
  /// HTTP: POST /restapis/{restApiId}/models
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Model> createModel(
    $0.CreateModelRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createModel, request, options: options);
  }

  /// Creates a RequestValidator of a given RestApi.
  /// HTTP: POST /restapis/{restApiId}/requestvalidators
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.RequestValidator> createRequestValidator(
    $0.CreateRequestValidatorRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createRequestValidator, request,
        options: options);
  }

  /// Creates a Resource resource.
  /// HTTP: POST /restapis/{restApiId}/resources/{parentId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Resource> createResource(
    $0.CreateResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createResource, request, options: options);
  }

  /// Creates a new RestApi resource.
  /// HTTP: POST /restapis
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.RestApi> createRestApi(
    $0.CreateRestApiRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createRestApi, request, options: options);
  }

  /// Creates a new Stage resource that references a pre-existing Deployment for the API.
  /// HTTP: POST /restapis/{restApiId}/stages
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Stage> createStage(
    $0.CreateStageRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createStage, request, options: options);
  }

  /// Creates a usage plan with the throttle and quota limits, as well as the associated API stages, specified in the payload.
  /// HTTP: POST /usageplans
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.UsagePlan> createUsagePlan(
    $0.CreateUsagePlanRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createUsagePlan, request, options: options);
  }

  /// Creates a usage plan key for adding an existing API key to a usage plan.
  /// HTTP: POST /usageplans/{usagePlanId}/keys
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.UsagePlanKey> createUsagePlanKey(
    $0.CreateUsagePlanKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createUsagePlanKey, request, options: options);
  }

  /// Creates a VPC link, under the caller's account in a selected region, in an asynchronous operation that typically takes 2-4 minutes to complete and become operational. The caller must have permissio...
  /// HTTP: POST /vpclinks
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.VpcLink> createVpcLink(
    $0.CreateVpcLinkRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createVpcLink, request, options: options);
  }

  /// Deletes the ApiKey resource.
  /// HTTP: DELETE /apikeys/{apiKey}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteApiKey(
    $0.DeleteApiKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteApiKey, request, options: options);
  }

  /// Deletes an existing Authorizer resource.
  /// HTTP: DELETE /restapis/{restApiId}/authorizers/{authorizerId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteAuthorizer(
    $0.DeleteAuthorizerRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteAuthorizer, request, options: options);
  }

  /// Deletes the BasePathMapping resource.
  /// HTTP: DELETE /domainnames/{domainName}/basepathmappings/{basePath}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteBasePathMapping(
    $0.DeleteBasePathMappingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteBasePathMapping, request, options: options);
  }

  /// Deletes the ClientCertificate resource.
  /// HTTP: DELETE /clientcertificates/{clientCertificateId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteClientCertificate(
    $0.DeleteClientCertificateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteClientCertificate, request,
        options: options);
  }

  /// Deletes a Deployment resource. Deleting a deployment will only succeed if there are no Stage resources associated with it.
  /// HTTP: DELETE /restapis/{restApiId}/deployments/{deploymentId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteDeployment(
    $0.DeleteDeploymentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteDeployment, request, options: options);
  }

  /// Deletes a documentation part
  /// HTTP: DELETE /restapis/{restApiId}/documentation/parts/{documentationPartId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteDocumentationPart(
    $0.DeleteDocumentationPartRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteDocumentationPart, request,
        options: options);
  }

  /// Deletes a documentation version.
  /// HTTP: DELETE /restapis/{restApiId}/documentation/versions/{documentationVersion}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteDocumentationVersion(
    $0.DeleteDocumentationVersionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteDocumentationVersion, request,
        options: options);
  }

  /// Deletes the DomainName resource.
  /// HTTP: DELETE /domainnames/{domainName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteDomainName(
    $0.DeleteDomainNameRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteDomainName, request, options: options);
  }

  /// Deletes the DomainNameAccessAssociation resource. Only the AWS account that created the DomainNameAccessAssociation resource can delete it. To stop an access association source in another AWS accou...
  /// HTTP: DELETE /domainnameaccessassociations/{domainNameAccessAssociationArn}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteDomainNameAccessAssociation(
    $0.DeleteDomainNameAccessAssociationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteDomainNameAccessAssociation, request,
        options: options);
  }

  /// Clears any customization of a GatewayResponse of a specified response type on the given RestApi and resets it with the default settings.
  /// HTTP: DELETE /restapis/{restApiId}/gatewayresponses/{responseType}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteGatewayResponse(
    $0.DeleteGatewayResponseRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteGatewayResponse, request, options: options);
  }

  /// Represents a delete integration.
  /// HTTP: DELETE /restapis/{restApiId}/resources/{resourceId}/methods/{httpMethod}/integration
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteIntegration(
    $0.DeleteIntegrationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteIntegration, request, options: options);
  }

  /// Represents a delete integration response.
  /// HTTP: DELETE /restapis/{restApiId}/resources/{resourceId}/methods/{httpMethod}/integration/responses/{statusCode}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteIntegrationResponse(
    $0.DeleteIntegrationResponseRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteIntegrationResponse, request,
        options: options);
  }

  /// Deletes an existing Method resource.
  /// HTTP: DELETE /restapis/{restApiId}/resources/{resourceId}/methods/{httpMethod}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteMethod(
    $0.DeleteMethodRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteMethod, request, options: options);
  }

  /// Deletes an existing MethodResponse resource.
  /// HTTP: DELETE /restapis/{restApiId}/resources/{resourceId}/methods/{httpMethod}/responses/{statusCode}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteMethodResponse(
    $0.DeleteMethodResponseRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteMethodResponse, request, options: options);
  }

  /// Deletes a model.
  /// HTTP: DELETE /restapis/{restApiId}/models/{modelName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteModel(
    $0.DeleteModelRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteModel, request, options: options);
  }

  /// Deletes a RequestValidator of a given RestApi.
  /// HTTP: DELETE /restapis/{restApiId}/requestvalidators/{requestValidatorId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteRequestValidator(
    $0.DeleteRequestValidatorRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteRequestValidator, request,
        options: options);
  }

  /// Deletes a Resource resource.
  /// HTTP: DELETE /restapis/{restApiId}/resources/{resourceId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteResource(
    $0.DeleteResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteResource, request, options: options);
  }

  /// Deletes the specified API.
  /// HTTP: DELETE /restapis/{restApiId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteRestApi(
    $0.DeleteRestApiRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteRestApi, request, options: options);
  }

  /// Deletes a Stage resource.
  /// HTTP: DELETE /restapis/{restApiId}/stages/{stageName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteStage(
    $0.DeleteStageRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteStage, request, options: options);
  }

  /// Deletes a usage plan of a given plan Id.
  /// HTTP: DELETE /usageplans/{usagePlanId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteUsagePlan(
    $0.DeleteUsagePlanRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteUsagePlan, request, options: options);
  }

  /// Deletes a usage plan key and remove the underlying API key from the associated usage plan.
  /// HTTP: DELETE /usageplans/{usagePlanId}/keys/{keyId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteUsagePlanKey(
    $0.DeleteUsagePlanKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteUsagePlanKey, request, options: options);
  }

  /// Deletes an existing VpcLink of a specified identifier.
  /// HTTP: DELETE /vpclinks/{vpcLinkId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteVpcLink(
    $0.DeleteVpcLinkRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteVpcLink, request, options: options);
  }

  /// Flushes all authorizer cache entries on a stage.
  /// HTTP: DELETE /restapis/{restApiId}/stages/{stageName}/cache/authorizers
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> flushStageAuthorizersCache(
    $0.FlushStageAuthorizersCacheRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$flushStageAuthorizersCache, request,
        options: options);
  }

  /// Flushes a stage's cache.
  /// HTTP: DELETE /restapis/{restApiId}/stages/{stageName}/cache/data
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> flushStageCache(
    $0.FlushStageCacheRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$flushStageCache, request, options: options);
  }

  /// Generates a ClientCertificate resource.
  /// HTTP: POST /clientcertificates
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ClientCertificate> generateClientCertificate(
    $0.GenerateClientCertificateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$generateClientCertificate, request,
        options: options);
  }

  /// Gets information about the current Account resource.
  /// HTTP: GET /account
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Account> getAccount(
    $0.GetAccountRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAccount, request, options: options);
  }

  /// Gets information about the current ApiKey resource.
  /// HTTP: GET /apikeys/{apiKey}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ApiKey> getApiKey(
    $0.GetApiKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getApiKey, request, options: options);
  }

  /// Gets information about the current ApiKeys resource.
  /// HTTP: GET /apikeys
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ApiKeys> getApiKeys(
    $0.GetApiKeysRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getApiKeys, request, options: options);
  }

  /// Describe an existing Authorizer resource.
  /// HTTP: GET /restapis/{restApiId}/authorizers/{authorizerId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Authorizer> getAuthorizer(
    $0.GetAuthorizerRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAuthorizer, request, options: options);
  }

  /// Describe an existing Authorizers resource.
  /// HTTP: GET /restapis/{restApiId}/authorizers
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Authorizers> getAuthorizers(
    $0.GetAuthorizersRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAuthorizers, request, options: options);
  }

  /// Describe a BasePathMapping resource.
  /// HTTP: GET /domainnames/{domainName}/basepathmappings/{basePath}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.BasePathMapping> getBasePathMapping(
    $0.GetBasePathMappingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBasePathMapping, request, options: options);
  }

  /// Represents a collection of BasePathMapping resources.
  /// HTTP: GET /domainnames/{domainName}/basepathmappings
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.BasePathMappings> getBasePathMappings(
    $0.GetBasePathMappingsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBasePathMappings, request, options: options);
  }

  /// Gets information about the current ClientCertificate resource.
  /// HTTP: GET /clientcertificates/{clientCertificateId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ClientCertificate> getClientCertificate(
    $0.GetClientCertificateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getClientCertificate, request, options: options);
  }

  /// Gets a collection of ClientCertificate resources.
  /// HTTP: GET /clientcertificates
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ClientCertificates> getClientCertificates(
    $0.GetClientCertificatesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getClientCertificates, request, options: options);
  }

  /// Gets information about a Deployment resource.
  /// HTTP: GET /restapis/{restApiId}/deployments/{deploymentId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Deployment> getDeployment(
    $0.GetDeploymentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDeployment, request, options: options);
  }

  /// Gets information about a Deployments collection.
  /// HTTP: GET /restapis/{restApiId}/deployments
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Deployments> getDeployments(
    $0.GetDeploymentsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDeployments, request, options: options);
  }

  /// Gets a documentation part.
  /// HTTP: GET /restapis/{restApiId}/documentation/parts/{documentationPartId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DocumentationPart> getDocumentationPart(
    $0.GetDocumentationPartRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDocumentationPart, request, options: options);
  }

  /// Gets documentation parts.
  /// HTTP: GET /restapis/{restApiId}/documentation/parts
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DocumentationParts> getDocumentationParts(
    $0.GetDocumentationPartsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDocumentationParts, request, options: options);
  }

  /// Gets a documentation version.
  /// HTTP: GET /restapis/{restApiId}/documentation/versions/{documentationVersion}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DocumentationVersion> getDocumentationVersion(
    $0.GetDocumentationVersionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDocumentationVersion, request,
        options: options);
  }

  /// Gets documentation versions.
  /// HTTP: GET /restapis/{restApiId}/documentation/versions
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DocumentationVersions> getDocumentationVersions(
    $0.GetDocumentationVersionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDocumentationVersions, request,
        options: options);
  }

  /// Represents a domain name that is contained in a simpler, more intuitive URL that can be called.
  /// HTTP: GET /domainnames/{domainName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DomainName> getDomainName(
    $0.GetDomainNameRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDomainName, request, options: options);
  }

  /// Represents a collection on DomainNameAccessAssociations resources.
  /// HTTP: GET /domainnameaccessassociations
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DomainNameAccessAssociations>
      getDomainNameAccessAssociations(
    $0.GetDomainNameAccessAssociationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDomainNameAccessAssociations, request,
        options: options);
  }

  /// Represents a collection of DomainName resources.
  /// HTTP: GET /domainnames
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DomainNames> getDomainNames(
    $0.GetDomainNamesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDomainNames, request, options: options);
  }

  /// Exports a deployed version of a RestApi in a specified format.
  /// HTTP: GET /restapis/{restApiId}/stages/{stageName}/exports/{exportType}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ExportResponse> getExport(
    $0.GetExportRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getExport, request, options: options);
  }

  /// Gets a GatewayResponse of a specified response type on the given RestApi.
  /// HTTP: GET /restapis/{restApiId}/gatewayresponses/{responseType}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GatewayResponse> getGatewayResponse(
    $0.GetGatewayResponseRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getGatewayResponse, request, options: options);
  }

  /// Gets the GatewayResponses collection on the given RestApi. If an API developer has not added any definitions for gateway responses, the result will be the API Gateway-generated default GatewayRespo...
  /// HTTP: GET /restapis/{restApiId}/gatewayresponses
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GatewayResponses> getGatewayResponses(
    $0.GetGatewayResponsesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getGatewayResponses, request, options: options);
  }

  /// Get the integration settings.
  /// HTTP: GET /restapis/{restApiId}/resources/{resourceId}/methods/{httpMethod}/integration
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Integration> getIntegration(
    $0.GetIntegrationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getIntegration, request, options: options);
  }

  /// Represents a get integration response.
  /// HTTP: GET /restapis/{restApiId}/resources/{resourceId}/methods/{httpMethod}/integration/responses/{statusCode}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.IntegrationResponse> getIntegrationResponse(
    $0.GetIntegrationResponseRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getIntegrationResponse, request,
        options: options);
  }

  /// Describe an existing Method resource.
  /// HTTP: GET /restapis/{restApiId}/resources/{resourceId}/methods/{httpMethod}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Method> getMethod(
    $0.GetMethodRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getMethod, request, options: options);
  }

  /// Describes a MethodResponse resource.
  /// HTTP: GET /restapis/{restApiId}/resources/{resourceId}/methods/{httpMethod}/responses/{statusCode}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.MethodResponse> getMethodResponse(
    $0.GetMethodResponseRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getMethodResponse, request, options: options);
  }

  /// Describes an existing model defined for a RestApi resource.
  /// HTTP: GET /restapis/{restApiId}/models/{modelName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Model> getModel(
    $0.GetModelRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getModel, request, options: options);
  }

  /// Describes existing Models defined for a RestApi resource.
  /// HTTP: GET /restapis/{restApiId}/models
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Models> getModels(
    $0.GetModelsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getModels, request, options: options);
  }

  /// Generates a sample mapping template that can be used to transform a payload into the structure of a model.
  /// HTTP: GET /restapis/{restApiId}/models/{modelName}/default_template
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Template> getModelTemplate(
    $0.GetModelTemplateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getModelTemplate, request, options: options);
  }

  /// Gets a RequestValidator of a given RestApi.
  /// HTTP: GET /restapis/{restApiId}/requestvalidators/{requestValidatorId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.RequestValidator> getRequestValidator(
    $0.GetRequestValidatorRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRequestValidator, request, options: options);
  }

  /// Gets the RequestValidators collection of a given RestApi.
  /// HTTP: GET /restapis/{restApiId}/requestvalidators
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.RequestValidators> getRequestValidators(
    $0.GetRequestValidatorsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRequestValidators, request, options: options);
  }

  /// Lists information about a resource.
  /// HTTP: GET /restapis/{restApiId}/resources/{resourceId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Resource> getResource(
    $0.GetResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getResource, request, options: options);
  }

  /// Lists information about a collection of Resource resources.
  /// HTTP: GET /restapis/{restApiId}/resources
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Resources> getResources(
    $0.GetResourcesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getResources, request, options: options);
  }

  /// Lists the RestApi resource in the collection.
  /// HTTP: GET /restapis/{restApiId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.RestApi> getRestApi(
    $0.GetRestApiRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRestApi, request, options: options);
  }

  /// Lists the RestApis resources for your collection.
  /// HTTP: GET /restapis
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.RestApis> getRestApis(
    $0.GetRestApisRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRestApis, request, options: options);
  }

  /// Generates a client SDK for a RestApi and Stage.
  /// HTTP: GET /restapis/{restApiId}/stages/{stageName}/sdks/{sdkType}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.SdkResponse> getSdk(
    $0.GetSdkRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSdk, request, options: options);
  }

  /// Gets an SDK type.
  /// HTTP: GET /sdktypes/{id}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.SdkType> getSdkType(
    $0.GetSdkTypeRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSdkType, request, options: options);
  }

  /// Gets SDK types
  /// HTTP: GET /sdktypes
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.SdkTypes> getSdkTypes(
    $0.GetSdkTypesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSdkTypes, request, options: options);
  }

  /// Gets information about a Stage resource.
  /// HTTP: GET /restapis/{restApiId}/stages/{stageName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Stage> getStage(
    $0.GetStageRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getStage, request, options: options);
  }

  /// Gets information about one or more Stage resources.
  /// HTTP: GET /restapis/{restApiId}/stages
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Stages> getStages(
    $0.GetStagesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getStages, request, options: options);
  }

  /// Gets the Tags collection for a given resource.
  /// HTTP: GET /tags/{resourceArn}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Tags> getTags(
    $0.GetTagsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getTags, request, options: options);
  }

  /// Gets the usage data of a usage plan in a specified time interval.
  /// HTTP: GET /usageplans/{usagePlanId}/usage
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Usage> getUsage(
    $0.GetUsageRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getUsage, request, options: options);
  }

  /// Gets a usage plan of a given plan identifier.
  /// HTTP: GET /usageplans/{usagePlanId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.UsagePlan> getUsagePlan(
    $0.GetUsagePlanRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getUsagePlan, request, options: options);
  }

  /// Gets a usage plan key of a given key identifier.
  /// HTTP: GET /usageplans/{usagePlanId}/keys/{keyId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.UsagePlanKey> getUsagePlanKey(
    $0.GetUsagePlanKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getUsagePlanKey, request, options: options);
  }

  /// Gets all the usage plan keys representing the API keys added to a specified usage plan.
  /// HTTP: GET /usageplans/{usagePlanId}/keys
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.UsagePlanKeys> getUsagePlanKeys(
    $0.GetUsagePlanKeysRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getUsagePlanKeys, request, options: options);
  }

  /// Gets all the usage plans of the caller's account.
  /// HTTP: GET /usageplans
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.UsagePlans> getUsagePlans(
    $0.GetUsagePlansRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getUsagePlans, request, options: options);
  }

  /// Gets a specified VPC link under the caller's account in a region.
  /// HTTP: GET /vpclinks/{vpcLinkId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.VpcLink> getVpcLink(
    $0.GetVpcLinkRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getVpcLink, request, options: options);
  }

  /// Gets the VpcLinks collection under the caller's account in a selected region.
  /// HTTP: GET /vpclinks
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.VpcLinks> getVpcLinks(
    $0.GetVpcLinksRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getVpcLinks, request, options: options);
  }

  /// Import API keys from an external source, such as a CSV-formatted file.
  /// HTTP: POST /apikeys?mode=import
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ApiKeyIds> importApiKeys(
    $0.ImportApiKeysRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$importApiKeys, request, options: options);
  }

  /// Imports documentation parts
  /// HTTP: PUT /restapis/{restApiId}/documentation/parts
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DocumentationPartIds> importDocumentationParts(
    $0.ImportDocumentationPartsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$importDocumentationParts, request,
        options: options);
  }

  /// A feature of the API Gateway control service for creating a new API from an external API definition file.
  /// HTTP: POST /restapis?mode=import
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.RestApi> importRestApi(
    $0.ImportRestApiRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$importRestApi, request, options: options);
  }

  /// Creates a customization of a GatewayResponse of a specified response type and status code on the given RestApi.
  /// HTTP: PUT /restapis/{restApiId}/gatewayresponses/{responseType}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GatewayResponse> putGatewayResponse(
    $0.PutGatewayResponseRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putGatewayResponse, request, options: options);
  }

  /// Sets up a method's integration.
  /// HTTP: PUT /restapis/{restApiId}/resources/{resourceId}/methods/{httpMethod}/integration
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Integration> putIntegration(
    $0.PutIntegrationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putIntegration, request, options: options);
  }

  /// Represents a put integration.
  /// HTTP: PUT /restapis/{restApiId}/resources/{resourceId}/methods/{httpMethod}/integration/responses/{statusCode}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.IntegrationResponse> putIntegrationResponse(
    $0.PutIntegrationResponseRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putIntegrationResponse, request,
        options: options);
  }

  /// Add a method to an existing Resource resource.
  /// HTTP: PUT /restapis/{restApiId}/resources/{resourceId}/methods/{httpMethod}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Method> putMethod(
    $0.PutMethodRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putMethod, request, options: options);
  }

  /// Adds a MethodResponse to an existing Method resource.
  /// HTTP: PUT /restapis/{restApiId}/resources/{resourceId}/methods/{httpMethod}/responses/{statusCode}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.MethodResponse> putMethodResponse(
    $0.PutMethodResponseRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putMethodResponse, request, options: options);
  }

  /// A feature of the API Gateway control service for updating an existing API with an input of external API definitions. The update can take the form of merging the supplied definition into the existin...
  /// HTTP: PUT /restapis/{restApiId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.RestApi> putRestApi(
    $0.PutRestApiRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putRestApi, request, options: options);
  }

  /// Rejects a domain name access association with a private custom domain name. To reject a domain name access association with an access association source in another AWS account, use this operation. ...
  /// HTTP: POST /rejectdomainnameaccessassociations
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> rejectDomainNameAccessAssociation(
    $0.RejectDomainNameAccessAssociationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$rejectDomainNameAccessAssociation, request,
        options: options);
  }

  /// Adds or updates a tag on a given resource.
  /// HTTP: PUT /tags/{resourceArn}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> tagResource(
    $0.TagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagResource, request, options: options);
  }

  /// Simulate the execution of an Authorizer in your RestApi with headers, parameters, and an incoming request body.
  /// HTTP: POST /restapis/{restApiId}/authorizers/{authorizerId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.TestInvokeAuthorizerResponse> testInvokeAuthorizer(
    $0.TestInvokeAuthorizerRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$testInvokeAuthorizer, request, options: options);
  }

  /// Simulate the invocation of a Method in your RestApi with headers, parameters, and an incoming request body.
  /// HTTP: POST /restapis/{restApiId}/resources/{resourceId}/methods/{httpMethod}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.TestInvokeMethodResponse> testInvokeMethod(
    $0.TestInvokeMethodRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$testInvokeMethod, request, options: options);
  }

  /// Removes a tag from a given resource.
  /// HTTP: DELETE /tags/{resourceArn}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> untagResource(
    $0.UntagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagResource, request, options: options);
  }

  /// Changes information about the current Account resource.
  /// HTTP: PATCH /account
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Account> updateAccount(
    $0.UpdateAccountRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateAccount, request, options: options);
  }

  /// Changes information about an ApiKey resource.
  /// HTTP: PATCH /apikeys/{apiKey}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ApiKey> updateApiKey(
    $0.UpdateApiKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateApiKey, request, options: options);
  }

  /// Updates an existing Authorizer resource.
  /// HTTP: PATCH /restapis/{restApiId}/authorizers/{authorizerId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Authorizer> updateAuthorizer(
    $0.UpdateAuthorizerRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateAuthorizer, request, options: options);
  }

  /// Changes information about the BasePathMapping resource.
  /// HTTP: PATCH /domainnames/{domainName}/basepathmappings/{basePath}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.BasePathMapping> updateBasePathMapping(
    $0.UpdateBasePathMappingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateBasePathMapping, request, options: options);
  }

  /// Changes information about an ClientCertificate resource.
  /// HTTP: PATCH /clientcertificates/{clientCertificateId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ClientCertificate> updateClientCertificate(
    $0.UpdateClientCertificateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateClientCertificate, request,
        options: options);
  }

  /// Changes information about a Deployment resource.
  /// HTTP: PATCH /restapis/{restApiId}/deployments/{deploymentId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Deployment> updateDeployment(
    $0.UpdateDeploymentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateDeployment, request, options: options);
  }

  /// Updates a documentation part.
  /// HTTP: PATCH /restapis/{restApiId}/documentation/parts/{documentationPartId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DocumentationPart> updateDocumentationPart(
    $0.UpdateDocumentationPartRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateDocumentationPart, request,
        options: options);
  }

  /// Updates a documentation version.
  /// HTTP: PATCH /restapis/{restApiId}/documentation/versions/{documentationVersion}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DocumentationVersion> updateDocumentationVersion(
    $0.UpdateDocumentationVersionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateDocumentationVersion, request,
        options: options);
  }

  /// Changes information about the DomainName resource.
  /// HTTP: PATCH /domainnames/{domainName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DomainName> updateDomainName(
    $0.UpdateDomainNameRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateDomainName, request, options: options);
  }

  /// Updates a GatewayResponse of a specified response type on the given RestApi.
  /// HTTP: PATCH /restapis/{restApiId}/gatewayresponses/{responseType}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GatewayResponse> updateGatewayResponse(
    $0.UpdateGatewayResponseRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateGatewayResponse, request, options: options);
  }

  /// Represents an update integration.
  /// HTTP: PATCH /restapis/{restApiId}/resources/{resourceId}/methods/{httpMethod}/integration
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Integration> updateIntegration(
    $0.UpdateIntegrationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateIntegration, request, options: options);
  }

  /// Represents an update integration response.
  /// HTTP: PATCH /restapis/{restApiId}/resources/{resourceId}/methods/{httpMethod}/integration/responses/{statusCode}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.IntegrationResponse> updateIntegrationResponse(
    $0.UpdateIntegrationResponseRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateIntegrationResponse, request,
        options: options);
  }

  /// Updates an existing Method resource.
  /// HTTP: PATCH /restapis/{restApiId}/resources/{resourceId}/methods/{httpMethod}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Method> updateMethod(
    $0.UpdateMethodRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateMethod, request, options: options);
  }

  /// Updates an existing MethodResponse resource.
  /// HTTP: PATCH /restapis/{restApiId}/resources/{resourceId}/methods/{httpMethod}/responses/{statusCode}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.MethodResponse> updateMethodResponse(
    $0.UpdateMethodResponseRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateMethodResponse, request, options: options);
  }

  /// Changes information about a model. The maximum size of the model is 400 KB.
  /// HTTP: PATCH /restapis/{restApiId}/models/{modelName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Model> updateModel(
    $0.UpdateModelRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateModel, request, options: options);
  }

  /// Updates a RequestValidator of a given RestApi.
  /// HTTP: PATCH /restapis/{restApiId}/requestvalidators/{requestValidatorId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.RequestValidator> updateRequestValidator(
    $0.UpdateRequestValidatorRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateRequestValidator, request,
        options: options);
  }

  /// Changes information about a Resource resource.
  /// HTTP: PATCH /restapis/{restApiId}/resources/{resourceId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Resource> updateResource(
    $0.UpdateResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateResource, request, options: options);
  }

  /// Changes information about the specified API.
  /// HTTP: PATCH /restapis/{restApiId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.RestApi> updateRestApi(
    $0.UpdateRestApiRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateRestApi, request, options: options);
  }

  /// Changes information about a Stage resource.
  /// HTTP: PATCH /restapis/{restApiId}/stages/{stageName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Stage> updateStage(
    $0.UpdateStageRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateStage, request, options: options);
  }

  /// Grants a temporary extension to the remaining quota of a usage plan associated with a specified API key.
  /// HTTP: PATCH /usageplans/{usagePlanId}/keys/{keyId}/usage
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Usage> updateUsage(
    $0.UpdateUsageRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateUsage, request, options: options);
  }

  /// Updates a usage plan of a given plan Id.
  /// HTTP: PATCH /usageplans/{usagePlanId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.UsagePlan> updateUsagePlan(
    $0.UpdateUsagePlanRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateUsagePlan, request, options: options);
  }

  /// Updates an existing VpcLink of a specified identifier.
  /// HTTP: PATCH /vpclinks/{vpcLinkId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.VpcLink> updateVpcLink(
    $0.UpdateVpcLinkRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateVpcLink, request, options: options);
  }

  // method descriptors

  static final _$createApiKey =
      $grpc.ClientMethod<$0.CreateApiKeyRequest, $0.ApiKey>(
          '/apigateway.APIGatewayService/CreateApiKey',
          ($0.CreateApiKeyRequest value) => value.writeToBuffer(),
          $0.ApiKey.fromBuffer);
  static final _$createAuthorizer =
      $grpc.ClientMethod<$0.CreateAuthorizerRequest, $0.Authorizer>(
          '/apigateway.APIGatewayService/CreateAuthorizer',
          ($0.CreateAuthorizerRequest value) => value.writeToBuffer(),
          $0.Authorizer.fromBuffer);
  static final _$createBasePathMapping =
      $grpc.ClientMethod<$0.CreateBasePathMappingRequest, $0.BasePathMapping>(
          '/apigateway.APIGatewayService/CreateBasePathMapping',
          ($0.CreateBasePathMappingRequest value) => value.writeToBuffer(),
          $0.BasePathMapping.fromBuffer);
  static final _$createDeployment =
      $grpc.ClientMethod<$0.CreateDeploymentRequest, $0.Deployment>(
          '/apigateway.APIGatewayService/CreateDeployment',
          ($0.CreateDeploymentRequest value) => value.writeToBuffer(),
          $0.Deployment.fromBuffer);
  static final _$createDocumentationPart = $grpc.ClientMethod<
          $0.CreateDocumentationPartRequest, $0.DocumentationPart>(
      '/apigateway.APIGatewayService/CreateDocumentationPart',
      ($0.CreateDocumentationPartRequest value) => value.writeToBuffer(),
      $0.DocumentationPart.fromBuffer);
  static final _$createDocumentationVersion = $grpc.ClientMethod<
          $0.CreateDocumentationVersionRequest, $0.DocumentationVersion>(
      '/apigateway.APIGatewayService/CreateDocumentationVersion',
      ($0.CreateDocumentationVersionRequest value) => value.writeToBuffer(),
      $0.DocumentationVersion.fromBuffer);
  static final _$createDomainName =
      $grpc.ClientMethod<$0.CreateDomainNameRequest, $0.DomainName>(
          '/apigateway.APIGatewayService/CreateDomainName',
          ($0.CreateDomainNameRequest value) => value.writeToBuffer(),
          $0.DomainName.fromBuffer);
  static final _$createDomainNameAccessAssociation = $grpc.ClientMethod<
          $0.CreateDomainNameAccessAssociationRequest,
          $0.DomainNameAccessAssociation>(
      '/apigateway.APIGatewayService/CreateDomainNameAccessAssociation',
      ($0.CreateDomainNameAccessAssociationRequest value) =>
          value.writeToBuffer(),
      $0.DomainNameAccessAssociation.fromBuffer);
  static final _$createModel =
      $grpc.ClientMethod<$0.CreateModelRequest, $0.Model>(
          '/apigateway.APIGatewayService/CreateModel',
          ($0.CreateModelRequest value) => value.writeToBuffer(),
          $0.Model.fromBuffer);
  static final _$createRequestValidator =
      $grpc.ClientMethod<$0.CreateRequestValidatorRequest, $0.RequestValidator>(
          '/apigateway.APIGatewayService/CreateRequestValidator',
          ($0.CreateRequestValidatorRequest value) => value.writeToBuffer(),
          $0.RequestValidator.fromBuffer);
  static final _$createResource =
      $grpc.ClientMethod<$0.CreateResourceRequest, $0.Resource>(
          '/apigateway.APIGatewayService/CreateResource',
          ($0.CreateResourceRequest value) => value.writeToBuffer(),
          $0.Resource.fromBuffer);
  static final _$createRestApi =
      $grpc.ClientMethod<$0.CreateRestApiRequest, $0.RestApi>(
          '/apigateway.APIGatewayService/CreateRestApi',
          ($0.CreateRestApiRequest value) => value.writeToBuffer(),
          $0.RestApi.fromBuffer);
  static final _$createStage =
      $grpc.ClientMethod<$0.CreateStageRequest, $0.Stage>(
          '/apigateway.APIGatewayService/CreateStage',
          ($0.CreateStageRequest value) => value.writeToBuffer(),
          $0.Stage.fromBuffer);
  static final _$createUsagePlan =
      $grpc.ClientMethod<$0.CreateUsagePlanRequest, $0.UsagePlan>(
          '/apigateway.APIGatewayService/CreateUsagePlan',
          ($0.CreateUsagePlanRequest value) => value.writeToBuffer(),
          $0.UsagePlan.fromBuffer);
  static final _$createUsagePlanKey =
      $grpc.ClientMethod<$0.CreateUsagePlanKeyRequest, $0.UsagePlanKey>(
          '/apigateway.APIGatewayService/CreateUsagePlanKey',
          ($0.CreateUsagePlanKeyRequest value) => value.writeToBuffer(),
          $0.UsagePlanKey.fromBuffer);
  static final _$createVpcLink =
      $grpc.ClientMethod<$0.CreateVpcLinkRequest, $0.VpcLink>(
          '/apigateway.APIGatewayService/CreateVpcLink',
          ($0.CreateVpcLinkRequest value) => value.writeToBuffer(),
          $0.VpcLink.fromBuffer);
  static final _$deleteApiKey =
      $grpc.ClientMethod<$0.DeleteApiKeyRequest, $1.Empty>(
          '/apigateway.APIGatewayService/DeleteApiKey',
          ($0.DeleteApiKeyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteAuthorizer =
      $grpc.ClientMethod<$0.DeleteAuthorizerRequest, $1.Empty>(
          '/apigateway.APIGatewayService/DeleteAuthorizer',
          ($0.DeleteAuthorizerRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteBasePathMapping =
      $grpc.ClientMethod<$0.DeleteBasePathMappingRequest, $1.Empty>(
          '/apigateway.APIGatewayService/DeleteBasePathMapping',
          ($0.DeleteBasePathMappingRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteClientCertificate =
      $grpc.ClientMethod<$0.DeleteClientCertificateRequest, $1.Empty>(
          '/apigateway.APIGatewayService/DeleteClientCertificate',
          ($0.DeleteClientCertificateRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteDeployment =
      $grpc.ClientMethod<$0.DeleteDeploymentRequest, $1.Empty>(
          '/apigateway.APIGatewayService/DeleteDeployment',
          ($0.DeleteDeploymentRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteDocumentationPart =
      $grpc.ClientMethod<$0.DeleteDocumentationPartRequest, $1.Empty>(
          '/apigateway.APIGatewayService/DeleteDocumentationPart',
          ($0.DeleteDocumentationPartRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteDocumentationVersion =
      $grpc.ClientMethod<$0.DeleteDocumentationVersionRequest, $1.Empty>(
          '/apigateway.APIGatewayService/DeleteDocumentationVersion',
          ($0.DeleteDocumentationVersionRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteDomainName =
      $grpc.ClientMethod<$0.DeleteDomainNameRequest, $1.Empty>(
          '/apigateway.APIGatewayService/DeleteDomainName',
          ($0.DeleteDomainNameRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteDomainNameAccessAssociation =
      $grpc.ClientMethod<$0.DeleteDomainNameAccessAssociationRequest, $1.Empty>(
          '/apigateway.APIGatewayService/DeleteDomainNameAccessAssociation',
          ($0.DeleteDomainNameAccessAssociationRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteGatewayResponse =
      $grpc.ClientMethod<$0.DeleteGatewayResponseRequest, $1.Empty>(
          '/apigateway.APIGatewayService/DeleteGatewayResponse',
          ($0.DeleteGatewayResponseRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteIntegration =
      $grpc.ClientMethod<$0.DeleteIntegrationRequest, $1.Empty>(
          '/apigateway.APIGatewayService/DeleteIntegration',
          ($0.DeleteIntegrationRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteIntegrationResponse =
      $grpc.ClientMethod<$0.DeleteIntegrationResponseRequest, $1.Empty>(
          '/apigateway.APIGatewayService/DeleteIntegrationResponse',
          ($0.DeleteIntegrationResponseRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteMethod =
      $grpc.ClientMethod<$0.DeleteMethodRequest, $1.Empty>(
          '/apigateway.APIGatewayService/DeleteMethod',
          ($0.DeleteMethodRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteMethodResponse =
      $grpc.ClientMethod<$0.DeleteMethodResponseRequest, $1.Empty>(
          '/apigateway.APIGatewayService/DeleteMethodResponse',
          ($0.DeleteMethodResponseRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteModel =
      $grpc.ClientMethod<$0.DeleteModelRequest, $1.Empty>(
          '/apigateway.APIGatewayService/DeleteModel',
          ($0.DeleteModelRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteRequestValidator =
      $grpc.ClientMethod<$0.DeleteRequestValidatorRequest, $1.Empty>(
          '/apigateway.APIGatewayService/DeleteRequestValidator',
          ($0.DeleteRequestValidatorRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteResource =
      $grpc.ClientMethod<$0.DeleteResourceRequest, $1.Empty>(
          '/apigateway.APIGatewayService/DeleteResource',
          ($0.DeleteResourceRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteRestApi =
      $grpc.ClientMethod<$0.DeleteRestApiRequest, $1.Empty>(
          '/apigateway.APIGatewayService/DeleteRestApi',
          ($0.DeleteRestApiRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteStage =
      $grpc.ClientMethod<$0.DeleteStageRequest, $1.Empty>(
          '/apigateway.APIGatewayService/DeleteStage',
          ($0.DeleteStageRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteUsagePlan =
      $grpc.ClientMethod<$0.DeleteUsagePlanRequest, $1.Empty>(
          '/apigateway.APIGatewayService/DeleteUsagePlan',
          ($0.DeleteUsagePlanRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteUsagePlanKey =
      $grpc.ClientMethod<$0.DeleteUsagePlanKeyRequest, $1.Empty>(
          '/apigateway.APIGatewayService/DeleteUsagePlanKey',
          ($0.DeleteUsagePlanKeyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteVpcLink =
      $grpc.ClientMethod<$0.DeleteVpcLinkRequest, $1.Empty>(
          '/apigateway.APIGatewayService/DeleteVpcLink',
          ($0.DeleteVpcLinkRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$flushStageAuthorizersCache =
      $grpc.ClientMethod<$0.FlushStageAuthorizersCacheRequest, $1.Empty>(
          '/apigateway.APIGatewayService/FlushStageAuthorizersCache',
          ($0.FlushStageAuthorizersCacheRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$flushStageCache =
      $grpc.ClientMethod<$0.FlushStageCacheRequest, $1.Empty>(
          '/apigateway.APIGatewayService/FlushStageCache',
          ($0.FlushStageCacheRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$generateClientCertificate = $grpc.ClientMethod<
          $0.GenerateClientCertificateRequest, $0.ClientCertificate>(
      '/apigateway.APIGatewayService/GenerateClientCertificate',
      ($0.GenerateClientCertificateRequest value) => value.writeToBuffer(),
      $0.ClientCertificate.fromBuffer);
  static final _$getAccount =
      $grpc.ClientMethod<$0.GetAccountRequest, $0.Account>(
          '/apigateway.APIGatewayService/GetAccount',
          ($0.GetAccountRequest value) => value.writeToBuffer(),
          $0.Account.fromBuffer);
  static final _$getApiKey = $grpc.ClientMethod<$0.GetApiKeyRequest, $0.ApiKey>(
      '/apigateway.APIGatewayService/GetApiKey',
      ($0.GetApiKeyRequest value) => value.writeToBuffer(),
      $0.ApiKey.fromBuffer);
  static final _$getApiKeys =
      $grpc.ClientMethod<$0.GetApiKeysRequest, $0.ApiKeys>(
          '/apigateway.APIGatewayService/GetApiKeys',
          ($0.GetApiKeysRequest value) => value.writeToBuffer(),
          $0.ApiKeys.fromBuffer);
  static final _$getAuthorizer =
      $grpc.ClientMethod<$0.GetAuthorizerRequest, $0.Authorizer>(
          '/apigateway.APIGatewayService/GetAuthorizer',
          ($0.GetAuthorizerRequest value) => value.writeToBuffer(),
          $0.Authorizer.fromBuffer);
  static final _$getAuthorizers =
      $grpc.ClientMethod<$0.GetAuthorizersRequest, $0.Authorizers>(
          '/apigateway.APIGatewayService/GetAuthorizers',
          ($0.GetAuthorizersRequest value) => value.writeToBuffer(),
          $0.Authorizers.fromBuffer);
  static final _$getBasePathMapping =
      $grpc.ClientMethod<$0.GetBasePathMappingRequest, $0.BasePathMapping>(
          '/apigateway.APIGatewayService/GetBasePathMapping',
          ($0.GetBasePathMappingRequest value) => value.writeToBuffer(),
          $0.BasePathMapping.fromBuffer);
  static final _$getBasePathMappings =
      $grpc.ClientMethod<$0.GetBasePathMappingsRequest, $0.BasePathMappings>(
          '/apigateway.APIGatewayService/GetBasePathMappings',
          ($0.GetBasePathMappingsRequest value) => value.writeToBuffer(),
          $0.BasePathMappings.fromBuffer);
  static final _$getClientCertificate =
      $grpc.ClientMethod<$0.GetClientCertificateRequest, $0.ClientCertificate>(
          '/apigateway.APIGatewayService/GetClientCertificate',
          ($0.GetClientCertificateRequest value) => value.writeToBuffer(),
          $0.ClientCertificate.fromBuffer);
  static final _$getClientCertificates = $grpc.ClientMethod<
          $0.GetClientCertificatesRequest, $0.ClientCertificates>(
      '/apigateway.APIGatewayService/GetClientCertificates',
      ($0.GetClientCertificatesRequest value) => value.writeToBuffer(),
      $0.ClientCertificates.fromBuffer);
  static final _$getDeployment =
      $grpc.ClientMethod<$0.GetDeploymentRequest, $0.Deployment>(
          '/apigateway.APIGatewayService/GetDeployment',
          ($0.GetDeploymentRequest value) => value.writeToBuffer(),
          $0.Deployment.fromBuffer);
  static final _$getDeployments =
      $grpc.ClientMethod<$0.GetDeploymentsRequest, $0.Deployments>(
          '/apigateway.APIGatewayService/GetDeployments',
          ($0.GetDeploymentsRequest value) => value.writeToBuffer(),
          $0.Deployments.fromBuffer);
  static final _$getDocumentationPart =
      $grpc.ClientMethod<$0.GetDocumentationPartRequest, $0.DocumentationPart>(
          '/apigateway.APIGatewayService/GetDocumentationPart',
          ($0.GetDocumentationPartRequest value) => value.writeToBuffer(),
          $0.DocumentationPart.fromBuffer);
  static final _$getDocumentationParts = $grpc.ClientMethod<
          $0.GetDocumentationPartsRequest, $0.DocumentationParts>(
      '/apigateway.APIGatewayService/GetDocumentationParts',
      ($0.GetDocumentationPartsRequest value) => value.writeToBuffer(),
      $0.DocumentationParts.fromBuffer);
  static final _$getDocumentationVersion = $grpc.ClientMethod<
          $0.GetDocumentationVersionRequest, $0.DocumentationVersion>(
      '/apigateway.APIGatewayService/GetDocumentationVersion',
      ($0.GetDocumentationVersionRequest value) => value.writeToBuffer(),
      $0.DocumentationVersion.fromBuffer);
  static final _$getDocumentationVersions = $grpc.ClientMethod<
          $0.GetDocumentationVersionsRequest, $0.DocumentationVersions>(
      '/apigateway.APIGatewayService/GetDocumentationVersions',
      ($0.GetDocumentationVersionsRequest value) => value.writeToBuffer(),
      $0.DocumentationVersions.fromBuffer);
  static final _$getDomainName =
      $grpc.ClientMethod<$0.GetDomainNameRequest, $0.DomainName>(
          '/apigateway.APIGatewayService/GetDomainName',
          ($0.GetDomainNameRequest value) => value.writeToBuffer(),
          $0.DomainName.fromBuffer);
  static final _$getDomainNameAccessAssociations = $grpc.ClientMethod<
          $0.GetDomainNameAccessAssociationsRequest,
          $0.DomainNameAccessAssociations>(
      '/apigateway.APIGatewayService/GetDomainNameAccessAssociations',
      ($0.GetDomainNameAccessAssociationsRequest value) =>
          value.writeToBuffer(),
      $0.DomainNameAccessAssociations.fromBuffer);
  static final _$getDomainNames =
      $grpc.ClientMethod<$0.GetDomainNamesRequest, $0.DomainNames>(
          '/apigateway.APIGatewayService/GetDomainNames',
          ($0.GetDomainNamesRequest value) => value.writeToBuffer(),
          $0.DomainNames.fromBuffer);
  static final _$getExport =
      $grpc.ClientMethod<$0.GetExportRequest, $0.ExportResponse>(
          '/apigateway.APIGatewayService/GetExport',
          ($0.GetExportRequest value) => value.writeToBuffer(),
          $0.ExportResponse.fromBuffer);
  static final _$getGatewayResponse =
      $grpc.ClientMethod<$0.GetGatewayResponseRequest, $0.GatewayResponse>(
          '/apigateway.APIGatewayService/GetGatewayResponse',
          ($0.GetGatewayResponseRequest value) => value.writeToBuffer(),
          $0.GatewayResponse.fromBuffer);
  static final _$getGatewayResponses =
      $grpc.ClientMethod<$0.GetGatewayResponsesRequest, $0.GatewayResponses>(
          '/apigateway.APIGatewayService/GetGatewayResponses',
          ($0.GetGatewayResponsesRequest value) => value.writeToBuffer(),
          $0.GatewayResponses.fromBuffer);
  static final _$getIntegration =
      $grpc.ClientMethod<$0.GetIntegrationRequest, $0.Integration>(
          '/apigateway.APIGatewayService/GetIntegration',
          ($0.GetIntegrationRequest value) => value.writeToBuffer(),
          $0.Integration.fromBuffer);
  static final _$getIntegrationResponse = $grpc.ClientMethod<
          $0.GetIntegrationResponseRequest, $0.IntegrationResponse>(
      '/apigateway.APIGatewayService/GetIntegrationResponse',
      ($0.GetIntegrationResponseRequest value) => value.writeToBuffer(),
      $0.IntegrationResponse.fromBuffer);
  static final _$getMethod = $grpc.ClientMethod<$0.GetMethodRequest, $0.Method>(
      '/apigateway.APIGatewayService/GetMethod',
      ($0.GetMethodRequest value) => value.writeToBuffer(),
      $0.Method.fromBuffer);
  static final _$getMethodResponse =
      $grpc.ClientMethod<$0.GetMethodResponseRequest, $0.MethodResponse>(
          '/apigateway.APIGatewayService/GetMethodResponse',
          ($0.GetMethodResponseRequest value) => value.writeToBuffer(),
          $0.MethodResponse.fromBuffer);
  static final _$getModel = $grpc.ClientMethod<$0.GetModelRequest, $0.Model>(
      '/apigateway.APIGatewayService/GetModel',
      ($0.GetModelRequest value) => value.writeToBuffer(),
      $0.Model.fromBuffer);
  static final _$getModels = $grpc.ClientMethod<$0.GetModelsRequest, $0.Models>(
      '/apigateway.APIGatewayService/GetModels',
      ($0.GetModelsRequest value) => value.writeToBuffer(),
      $0.Models.fromBuffer);
  static final _$getModelTemplate =
      $grpc.ClientMethod<$0.GetModelTemplateRequest, $0.Template>(
          '/apigateway.APIGatewayService/GetModelTemplate',
          ($0.GetModelTemplateRequest value) => value.writeToBuffer(),
          $0.Template.fromBuffer);
  static final _$getRequestValidator =
      $grpc.ClientMethod<$0.GetRequestValidatorRequest, $0.RequestValidator>(
          '/apigateway.APIGatewayService/GetRequestValidator',
          ($0.GetRequestValidatorRequest value) => value.writeToBuffer(),
          $0.RequestValidator.fromBuffer);
  static final _$getRequestValidators =
      $grpc.ClientMethod<$0.GetRequestValidatorsRequest, $0.RequestValidators>(
          '/apigateway.APIGatewayService/GetRequestValidators',
          ($0.GetRequestValidatorsRequest value) => value.writeToBuffer(),
          $0.RequestValidators.fromBuffer);
  static final _$getResource =
      $grpc.ClientMethod<$0.GetResourceRequest, $0.Resource>(
          '/apigateway.APIGatewayService/GetResource',
          ($0.GetResourceRequest value) => value.writeToBuffer(),
          $0.Resource.fromBuffer);
  static final _$getResources =
      $grpc.ClientMethod<$0.GetResourcesRequest, $0.Resources>(
          '/apigateway.APIGatewayService/GetResources',
          ($0.GetResourcesRequest value) => value.writeToBuffer(),
          $0.Resources.fromBuffer);
  static final _$getRestApi =
      $grpc.ClientMethod<$0.GetRestApiRequest, $0.RestApi>(
          '/apigateway.APIGatewayService/GetRestApi',
          ($0.GetRestApiRequest value) => value.writeToBuffer(),
          $0.RestApi.fromBuffer);
  static final _$getRestApis =
      $grpc.ClientMethod<$0.GetRestApisRequest, $0.RestApis>(
          '/apigateway.APIGatewayService/GetRestApis',
          ($0.GetRestApisRequest value) => value.writeToBuffer(),
          $0.RestApis.fromBuffer);
  static final _$getSdk = $grpc.ClientMethod<$0.GetSdkRequest, $0.SdkResponse>(
      '/apigateway.APIGatewayService/GetSdk',
      ($0.GetSdkRequest value) => value.writeToBuffer(),
      $0.SdkResponse.fromBuffer);
  static final _$getSdkType =
      $grpc.ClientMethod<$0.GetSdkTypeRequest, $0.SdkType>(
          '/apigateway.APIGatewayService/GetSdkType',
          ($0.GetSdkTypeRequest value) => value.writeToBuffer(),
          $0.SdkType.fromBuffer);
  static final _$getSdkTypes =
      $grpc.ClientMethod<$0.GetSdkTypesRequest, $0.SdkTypes>(
          '/apigateway.APIGatewayService/GetSdkTypes',
          ($0.GetSdkTypesRequest value) => value.writeToBuffer(),
          $0.SdkTypes.fromBuffer);
  static final _$getStage = $grpc.ClientMethod<$0.GetStageRequest, $0.Stage>(
      '/apigateway.APIGatewayService/GetStage',
      ($0.GetStageRequest value) => value.writeToBuffer(),
      $0.Stage.fromBuffer);
  static final _$getStages = $grpc.ClientMethod<$0.GetStagesRequest, $0.Stages>(
      '/apigateway.APIGatewayService/GetStages',
      ($0.GetStagesRequest value) => value.writeToBuffer(),
      $0.Stages.fromBuffer);
  static final _$getTags = $grpc.ClientMethod<$0.GetTagsRequest, $0.Tags>(
      '/apigateway.APIGatewayService/GetTags',
      ($0.GetTagsRequest value) => value.writeToBuffer(),
      $0.Tags.fromBuffer);
  static final _$getUsage = $grpc.ClientMethod<$0.GetUsageRequest, $0.Usage>(
      '/apigateway.APIGatewayService/GetUsage',
      ($0.GetUsageRequest value) => value.writeToBuffer(),
      $0.Usage.fromBuffer);
  static final _$getUsagePlan =
      $grpc.ClientMethod<$0.GetUsagePlanRequest, $0.UsagePlan>(
          '/apigateway.APIGatewayService/GetUsagePlan',
          ($0.GetUsagePlanRequest value) => value.writeToBuffer(),
          $0.UsagePlan.fromBuffer);
  static final _$getUsagePlanKey =
      $grpc.ClientMethod<$0.GetUsagePlanKeyRequest, $0.UsagePlanKey>(
          '/apigateway.APIGatewayService/GetUsagePlanKey',
          ($0.GetUsagePlanKeyRequest value) => value.writeToBuffer(),
          $0.UsagePlanKey.fromBuffer);
  static final _$getUsagePlanKeys =
      $grpc.ClientMethod<$0.GetUsagePlanKeysRequest, $0.UsagePlanKeys>(
          '/apigateway.APIGatewayService/GetUsagePlanKeys',
          ($0.GetUsagePlanKeysRequest value) => value.writeToBuffer(),
          $0.UsagePlanKeys.fromBuffer);
  static final _$getUsagePlans =
      $grpc.ClientMethod<$0.GetUsagePlansRequest, $0.UsagePlans>(
          '/apigateway.APIGatewayService/GetUsagePlans',
          ($0.GetUsagePlansRequest value) => value.writeToBuffer(),
          $0.UsagePlans.fromBuffer);
  static final _$getVpcLink =
      $grpc.ClientMethod<$0.GetVpcLinkRequest, $0.VpcLink>(
          '/apigateway.APIGatewayService/GetVpcLink',
          ($0.GetVpcLinkRequest value) => value.writeToBuffer(),
          $0.VpcLink.fromBuffer);
  static final _$getVpcLinks =
      $grpc.ClientMethod<$0.GetVpcLinksRequest, $0.VpcLinks>(
          '/apigateway.APIGatewayService/GetVpcLinks',
          ($0.GetVpcLinksRequest value) => value.writeToBuffer(),
          $0.VpcLinks.fromBuffer);
  static final _$importApiKeys =
      $grpc.ClientMethod<$0.ImportApiKeysRequest, $0.ApiKeyIds>(
          '/apigateway.APIGatewayService/ImportApiKeys',
          ($0.ImportApiKeysRequest value) => value.writeToBuffer(),
          $0.ApiKeyIds.fromBuffer);
  static final _$importDocumentationParts = $grpc.ClientMethod<
          $0.ImportDocumentationPartsRequest, $0.DocumentationPartIds>(
      '/apigateway.APIGatewayService/ImportDocumentationParts',
      ($0.ImportDocumentationPartsRequest value) => value.writeToBuffer(),
      $0.DocumentationPartIds.fromBuffer);
  static final _$importRestApi =
      $grpc.ClientMethod<$0.ImportRestApiRequest, $0.RestApi>(
          '/apigateway.APIGatewayService/ImportRestApi',
          ($0.ImportRestApiRequest value) => value.writeToBuffer(),
          $0.RestApi.fromBuffer);
  static final _$putGatewayResponse =
      $grpc.ClientMethod<$0.PutGatewayResponseRequest, $0.GatewayResponse>(
          '/apigateway.APIGatewayService/PutGatewayResponse',
          ($0.PutGatewayResponseRequest value) => value.writeToBuffer(),
          $0.GatewayResponse.fromBuffer);
  static final _$putIntegration =
      $grpc.ClientMethod<$0.PutIntegrationRequest, $0.Integration>(
          '/apigateway.APIGatewayService/PutIntegration',
          ($0.PutIntegrationRequest value) => value.writeToBuffer(),
          $0.Integration.fromBuffer);
  static final _$putIntegrationResponse = $grpc.ClientMethod<
          $0.PutIntegrationResponseRequest, $0.IntegrationResponse>(
      '/apigateway.APIGatewayService/PutIntegrationResponse',
      ($0.PutIntegrationResponseRequest value) => value.writeToBuffer(),
      $0.IntegrationResponse.fromBuffer);
  static final _$putMethod = $grpc.ClientMethod<$0.PutMethodRequest, $0.Method>(
      '/apigateway.APIGatewayService/PutMethod',
      ($0.PutMethodRequest value) => value.writeToBuffer(),
      $0.Method.fromBuffer);
  static final _$putMethodResponse =
      $grpc.ClientMethod<$0.PutMethodResponseRequest, $0.MethodResponse>(
          '/apigateway.APIGatewayService/PutMethodResponse',
          ($0.PutMethodResponseRequest value) => value.writeToBuffer(),
          $0.MethodResponse.fromBuffer);
  static final _$putRestApi =
      $grpc.ClientMethod<$0.PutRestApiRequest, $0.RestApi>(
          '/apigateway.APIGatewayService/PutRestApi',
          ($0.PutRestApiRequest value) => value.writeToBuffer(),
          $0.RestApi.fromBuffer);
  static final _$rejectDomainNameAccessAssociation =
      $grpc.ClientMethod<$0.RejectDomainNameAccessAssociationRequest, $1.Empty>(
          '/apigateway.APIGatewayService/RejectDomainNameAccessAssociation',
          ($0.RejectDomainNameAccessAssociationRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$tagResource =
      $grpc.ClientMethod<$0.TagResourceRequest, $1.Empty>(
          '/apigateway.APIGatewayService/TagResource',
          ($0.TagResourceRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$testInvokeAuthorizer = $grpc.ClientMethod<
          $0.TestInvokeAuthorizerRequest, $0.TestInvokeAuthorizerResponse>(
      '/apigateway.APIGatewayService/TestInvokeAuthorizer',
      ($0.TestInvokeAuthorizerRequest value) => value.writeToBuffer(),
      $0.TestInvokeAuthorizerResponse.fromBuffer);
  static final _$testInvokeMethod = $grpc.ClientMethod<
          $0.TestInvokeMethodRequest, $0.TestInvokeMethodResponse>(
      '/apigateway.APIGatewayService/TestInvokeMethod',
      ($0.TestInvokeMethodRequest value) => value.writeToBuffer(),
      $0.TestInvokeMethodResponse.fromBuffer);
  static final _$untagResource =
      $grpc.ClientMethod<$0.UntagResourceRequest, $1.Empty>(
          '/apigateway.APIGatewayService/UntagResource',
          ($0.UntagResourceRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updateAccount =
      $grpc.ClientMethod<$0.UpdateAccountRequest, $0.Account>(
          '/apigateway.APIGatewayService/UpdateAccount',
          ($0.UpdateAccountRequest value) => value.writeToBuffer(),
          $0.Account.fromBuffer);
  static final _$updateApiKey =
      $grpc.ClientMethod<$0.UpdateApiKeyRequest, $0.ApiKey>(
          '/apigateway.APIGatewayService/UpdateApiKey',
          ($0.UpdateApiKeyRequest value) => value.writeToBuffer(),
          $0.ApiKey.fromBuffer);
  static final _$updateAuthorizer =
      $grpc.ClientMethod<$0.UpdateAuthorizerRequest, $0.Authorizer>(
          '/apigateway.APIGatewayService/UpdateAuthorizer',
          ($0.UpdateAuthorizerRequest value) => value.writeToBuffer(),
          $0.Authorizer.fromBuffer);
  static final _$updateBasePathMapping =
      $grpc.ClientMethod<$0.UpdateBasePathMappingRequest, $0.BasePathMapping>(
          '/apigateway.APIGatewayService/UpdateBasePathMapping',
          ($0.UpdateBasePathMappingRequest value) => value.writeToBuffer(),
          $0.BasePathMapping.fromBuffer);
  static final _$updateClientCertificate = $grpc.ClientMethod<
          $0.UpdateClientCertificateRequest, $0.ClientCertificate>(
      '/apigateway.APIGatewayService/UpdateClientCertificate',
      ($0.UpdateClientCertificateRequest value) => value.writeToBuffer(),
      $0.ClientCertificate.fromBuffer);
  static final _$updateDeployment =
      $grpc.ClientMethod<$0.UpdateDeploymentRequest, $0.Deployment>(
          '/apigateway.APIGatewayService/UpdateDeployment',
          ($0.UpdateDeploymentRequest value) => value.writeToBuffer(),
          $0.Deployment.fromBuffer);
  static final _$updateDocumentationPart = $grpc.ClientMethod<
          $0.UpdateDocumentationPartRequest, $0.DocumentationPart>(
      '/apigateway.APIGatewayService/UpdateDocumentationPart',
      ($0.UpdateDocumentationPartRequest value) => value.writeToBuffer(),
      $0.DocumentationPart.fromBuffer);
  static final _$updateDocumentationVersion = $grpc.ClientMethod<
          $0.UpdateDocumentationVersionRequest, $0.DocumentationVersion>(
      '/apigateway.APIGatewayService/UpdateDocumentationVersion',
      ($0.UpdateDocumentationVersionRequest value) => value.writeToBuffer(),
      $0.DocumentationVersion.fromBuffer);
  static final _$updateDomainName =
      $grpc.ClientMethod<$0.UpdateDomainNameRequest, $0.DomainName>(
          '/apigateway.APIGatewayService/UpdateDomainName',
          ($0.UpdateDomainNameRequest value) => value.writeToBuffer(),
          $0.DomainName.fromBuffer);
  static final _$updateGatewayResponse =
      $grpc.ClientMethod<$0.UpdateGatewayResponseRequest, $0.GatewayResponse>(
          '/apigateway.APIGatewayService/UpdateGatewayResponse',
          ($0.UpdateGatewayResponseRequest value) => value.writeToBuffer(),
          $0.GatewayResponse.fromBuffer);
  static final _$updateIntegration =
      $grpc.ClientMethod<$0.UpdateIntegrationRequest, $0.Integration>(
          '/apigateway.APIGatewayService/UpdateIntegration',
          ($0.UpdateIntegrationRequest value) => value.writeToBuffer(),
          $0.Integration.fromBuffer);
  static final _$updateIntegrationResponse = $grpc.ClientMethod<
          $0.UpdateIntegrationResponseRequest, $0.IntegrationResponse>(
      '/apigateway.APIGatewayService/UpdateIntegrationResponse',
      ($0.UpdateIntegrationResponseRequest value) => value.writeToBuffer(),
      $0.IntegrationResponse.fromBuffer);
  static final _$updateMethod =
      $grpc.ClientMethod<$0.UpdateMethodRequest, $0.Method>(
          '/apigateway.APIGatewayService/UpdateMethod',
          ($0.UpdateMethodRequest value) => value.writeToBuffer(),
          $0.Method.fromBuffer);
  static final _$updateMethodResponse =
      $grpc.ClientMethod<$0.UpdateMethodResponseRequest, $0.MethodResponse>(
          '/apigateway.APIGatewayService/UpdateMethodResponse',
          ($0.UpdateMethodResponseRequest value) => value.writeToBuffer(),
          $0.MethodResponse.fromBuffer);
  static final _$updateModel =
      $grpc.ClientMethod<$0.UpdateModelRequest, $0.Model>(
          '/apigateway.APIGatewayService/UpdateModel',
          ($0.UpdateModelRequest value) => value.writeToBuffer(),
          $0.Model.fromBuffer);
  static final _$updateRequestValidator =
      $grpc.ClientMethod<$0.UpdateRequestValidatorRequest, $0.RequestValidator>(
          '/apigateway.APIGatewayService/UpdateRequestValidator',
          ($0.UpdateRequestValidatorRequest value) => value.writeToBuffer(),
          $0.RequestValidator.fromBuffer);
  static final _$updateResource =
      $grpc.ClientMethod<$0.UpdateResourceRequest, $0.Resource>(
          '/apigateway.APIGatewayService/UpdateResource',
          ($0.UpdateResourceRequest value) => value.writeToBuffer(),
          $0.Resource.fromBuffer);
  static final _$updateRestApi =
      $grpc.ClientMethod<$0.UpdateRestApiRequest, $0.RestApi>(
          '/apigateway.APIGatewayService/UpdateRestApi',
          ($0.UpdateRestApiRequest value) => value.writeToBuffer(),
          $0.RestApi.fromBuffer);
  static final _$updateStage =
      $grpc.ClientMethod<$0.UpdateStageRequest, $0.Stage>(
          '/apigateway.APIGatewayService/UpdateStage',
          ($0.UpdateStageRequest value) => value.writeToBuffer(),
          $0.Stage.fromBuffer);
  static final _$updateUsage =
      $grpc.ClientMethod<$0.UpdateUsageRequest, $0.Usage>(
          '/apigateway.APIGatewayService/UpdateUsage',
          ($0.UpdateUsageRequest value) => value.writeToBuffer(),
          $0.Usage.fromBuffer);
  static final _$updateUsagePlan =
      $grpc.ClientMethod<$0.UpdateUsagePlanRequest, $0.UsagePlan>(
          '/apigateway.APIGatewayService/UpdateUsagePlan',
          ($0.UpdateUsagePlanRequest value) => value.writeToBuffer(),
          $0.UsagePlan.fromBuffer);
  static final _$updateVpcLink =
      $grpc.ClientMethod<$0.UpdateVpcLinkRequest, $0.VpcLink>(
          '/apigateway.APIGatewayService/UpdateVpcLink',
          ($0.UpdateVpcLinkRequest value) => value.writeToBuffer(),
          $0.VpcLink.fromBuffer);
}

@$pb.GrpcServiceName('apigateway.APIGatewayService')
abstract class APIGatewayServiceBase extends $grpc.Service {
  $core.String get $name => 'apigateway.APIGatewayService';

  APIGatewayServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.CreateApiKeyRequest, $0.ApiKey>(
        'CreateApiKey',
        createApiKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateApiKeyRequest.fromBuffer(value),
        ($0.ApiKey value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateAuthorizerRequest, $0.Authorizer>(
        'CreateAuthorizer',
        createAuthorizer_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateAuthorizerRequest.fromBuffer(value),
        ($0.Authorizer value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateBasePathMappingRequest,
            $0.BasePathMapping>(
        'CreateBasePathMapping',
        createBasePathMapping_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateBasePathMappingRequest.fromBuffer(value),
        ($0.BasePathMapping value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateDeploymentRequest, $0.Deployment>(
        'CreateDeployment',
        createDeployment_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateDeploymentRequest.fromBuffer(value),
        ($0.Deployment value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateDocumentationPartRequest,
            $0.DocumentationPart>(
        'CreateDocumentationPart',
        createDocumentationPart_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateDocumentationPartRequest.fromBuffer(value),
        ($0.DocumentationPart value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateDocumentationVersionRequest,
            $0.DocumentationVersion>(
        'CreateDocumentationVersion',
        createDocumentationVersion_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateDocumentationVersionRequest.fromBuffer(value),
        ($0.DocumentationVersion value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateDomainNameRequest, $0.DomainName>(
        'CreateDomainName',
        createDomainName_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateDomainNameRequest.fromBuffer(value),
        ($0.DomainName value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateDomainNameAccessAssociationRequest,
            $0.DomainNameAccessAssociation>(
        'CreateDomainNameAccessAssociation',
        createDomainNameAccessAssociation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateDomainNameAccessAssociationRequest.fromBuffer(value),
        ($0.DomainNameAccessAssociation value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateModelRequest, $0.Model>(
        'CreateModel',
        createModel_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateModelRequest.fromBuffer(value),
        ($0.Model value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateRequestValidatorRequest,
            $0.RequestValidator>(
        'CreateRequestValidator',
        createRequestValidator_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateRequestValidatorRequest.fromBuffer(value),
        ($0.RequestValidator value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateResourceRequest, $0.Resource>(
        'CreateResource',
        createResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateResourceRequest.fromBuffer(value),
        ($0.Resource value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateRestApiRequest, $0.RestApi>(
        'CreateRestApi',
        createRestApi_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateRestApiRequest.fromBuffer(value),
        ($0.RestApi value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateStageRequest, $0.Stage>(
        'CreateStage',
        createStage_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateStageRequest.fromBuffer(value),
        ($0.Stage value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateUsagePlanRequest, $0.UsagePlan>(
        'CreateUsagePlan',
        createUsagePlan_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateUsagePlanRequest.fromBuffer(value),
        ($0.UsagePlan value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateUsagePlanKeyRequest, $0.UsagePlanKey>(
            'CreateUsagePlanKey',
            createUsagePlanKey_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateUsagePlanKeyRequest.fromBuffer(value),
            ($0.UsagePlanKey value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateVpcLinkRequest, $0.VpcLink>(
        'CreateVpcLink',
        createVpcLink_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateVpcLinkRequest.fromBuffer(value),
        ($0.VpcLink value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteApiKeyRequest, $1.Empty>(
        'DeleteApiKey',
        deleteApiKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteApiKeyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteAuthorizerRequest, $1.Empty>(
        'DeleteAuthorizer',
        deleteAuthorizer_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteAuthorizerRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteBasePathMappingRequest, $1.Empty>(
        'DeleteBasePathMapping',
        deleteBasePathMapping_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteBasePathMappingRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteClientCertificateRequest, $1.Empty>(
        'DeleteClientCertificate',
        deleteClientCertificate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteClientCertificateRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteDeploymentRequest, $1.Empty>(
        'DeleteDeployment',
        deleteDeployment_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteDeploymentRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteDocumentationPartRequest, $1.Empty>(
        'DeleteDocumentationPart',
        deleteDocumentationPart_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteDocumentationPartRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteDocumentationVersionRequest, $1.Empty>(
            'DeleteDocumentationVersion',
            deleteDocumentationVersion_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteDocumentationVersionRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteDomainNameRequest, $1.Empty>(
        'DeleteDomainName',
        deleteDomainName_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteDomainNameRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteDomainNameAccessAssociationRequest,
            $1.Empty>(
        'DeleteDomainNameAccessAssociation',
        deleteDomainNameAccessAssociation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteDomainNameAccessAssociationRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteGatewayResponseRequest, $1.Empty>(
        'DeleteGatewayResponse',
        deleteGatewayResponse_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteGatewayResponseRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteIntegrationRequest, $1.Empty>(
        'DeleteIntegration',
        deleteIntegration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteIntegrationRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteIntegrationResponseRequest, $1.Empty>(
            'DeleteIntegrationResponse',
            deleteIntegrationResponse_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteIntegrationResponseRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteMethodRequest, $1.Empty>(
        'DeleteMethod',
        deleteMethod_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteMethodRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteMethodResponseRequest, $1.Empty>(
        'DeleteMethodResponse',
        deleteMethodResponse_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteMethodResponseRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteModelRequest, $1.Empty>(
        'DeleteModel',
        deleteModel_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteModelRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteRequestValidatorRequest, $1.Empty>(
        'DeleteRequestValidator',
        deleteRequestValidator_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteRequestValidatorRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteResourceRequest, $1.Empty>(
        'DeleteResource',
        deleteResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteResourceRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteRestApiRequest, $1.Empty>(
        'DeleteRestApi',
        deleteRestApi_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteRestApiRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteStageRequest, $1.Empty>(
        'DeleteStage',
        deleteStage_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteStageRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteUsagePlanRequest, $1.Empty>(
        'DeleteUsagePlan',
        deleteUsagePlan_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteUsagePlanRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteUsagePlanKeyRequest, $1.Empty>(
        'DeleteUsagePlanKey',
        deleteUsagePlanKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteUsagePlanKeyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteVpcLinkRequest, $1.Empty>(
        'DeleteVpcLink',
        deleteVpcLink_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteVpcLinkRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.FlushStageAuthorizersCacheRequest, $1.Empty>(
            'FlushStageAuthorizersCache',
            flushStageAuthorizersCache_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.FlushStageAuthorizersCacheRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.FlushStageCacheRequest, $1.Empty>(
        'FlushStageCache',
        flushStageCache_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.FlushStageCacheRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GenerateClientCertificateRequest,
            $0.ClientCertificate>(
        'GenerateClientCertificate',
        generateClientCertificate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GenerateClientCertificateRequest.fromBuffer(value),
        ($0.ClientCertificate value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetAccountRequest, $0.Account>(
        'GetAccount',
        getAccount_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetAccountRequest.fromBuffer(value),
        ($0.Account value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetApiKeyRequest, $0.ApiKey>(
        'GetApiKey',
        getApiKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetApiKeyRequest.fromBuffer(value),
        ($0.ApiKey value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetApiKeysRequest, $0.ApiKeys>(
        'GetApiKeys',
        getApiKeys_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetApiKeysRequest.fromBuffer(value),
        ($0.ApiKeys value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetAuthorizerRequest, $0.Authorizer>(
        'GetAuthorizer',
        getAuthorizer_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetAuthorizerRequest.fromBuffer(value),
        ($0.Authorizer value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetAuthorizersRequest, $0.Authorizers>(
        'GetAuthorizers',
        getAuthorizers_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetAuthorizersRequest.fromBuffer(value),
        ($0.Authorizers value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetBasePathMappingRequest, $0.BasePathMapping>(
            'GetBasePathMapping',
            getBasePathMapping_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetBasePathMappingRequest.fromBuffer(value),
            ($0.BasePathMapping value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetBasePathMappingsRequest, $0.BasePathMappings>(
            'GetBasePathMappings',
            getBasePathMappings_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetBasePathMappingsRequest.fromBuffer(value),
            ($0.BasePathMappings value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetClientCertificateRequest,
            $0.ClientCertificate>(
        'GetClientCertificate',
        getClientCertificate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetClientCertificateRequest.fromBuffer(value),
        ($0.ClientCertificate value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetClientCertificatesRequest,
            $0.ClientCertificates>(
        'GetClientCertificates',
        getClientCertificates_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetClientCertificatesRequest.fromBuffer(value),
        ($0.ClientCertificates value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDeploymentRequest, $0.Deployment>(
        'GetDeployment',
        getDeployment_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDeploymentRequest.fromBuffer(value),
        ($0.Deployment value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDeploymentsRequest, $0.Deployments>(
        'GetDeployments',
        getDeployments_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDeploymentsRequest.fromBuffer(value),
        ($0.Deployments value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDocumentationPartRequest,
            $0.DocumentationPart>(
        'GetDocumentationPart',
        getDocumentationPart_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDocumentationPartRequest.fromBuffer(value),
        ($0.DocumentationPart value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDocumentationPartsRequest,
            $0.DocumentationParts>(
        'GetDocumentationParts',
        getDocumentationParts_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDocumentationPartsRequest.fromBuffer(value),
        ($0.DocumentationParts value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDocumentationVersionRequest,
            $0.DocumentationVersion>(
        'GetDocumentationVersion',
        getDocumentationVersion_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDocumentationVersionRequest.fromBuffer(value),
        ($0.DocumentationVersion value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDocumentationVersionsRequest,
            $0.DocumentationVersions>(
        'GetDocumentationVersions',
        getDocumentationVersions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDocumentationVersionsRequest.fromBuffer(value),
        ($0.DocumentationVersions value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDomainNameRequest, $0.DomainName>(
        'GetDomainName',
        getDomainName_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDomainNameRequest.fromBuffer(value),
        ($0.DomainName value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDomainNameAccessAssociationsRequest,
            $0.DomainNameAccessAssociations>(
        'GetDomainNameAccessAssociations',
        getDomainNameAccessAssociations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDomainNameAccessAssociationsRequest.fromBuffer(value),
        ($0.DomainNameAccessAssociations value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDomainNamesRequest, $0.DomainNames>(
        'GetDomainNames',
        getDomainNames_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDomainNamesRequest.fromBuffer(value),
        ($0.DomainNames value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetExportRequest, $0.ExportResponse>(
        'GetExport',
        getExport_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetExportRequest.fromBuffer(value),
        ($0.ExportResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetGatewayResponseRequest, $0.GatewayResponse>(
            'GetGatewayResponse',
            getGatewayResponse_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetGatewayResponseRequest.fromBuffer(value),
            ($0.GatewayResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetGatewayResponsesRequest, $0.GatewayResponses>(
            'GetGatewayResponses',
            getGatewayResponses_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetGatewayResponsesRequest.fromBuffer(value),
            ($0.GatewayResponses value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetIntegrationRequest, $0.Integration>(
        'GetIntegration',
        getIntegration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetIntegrationRequest.fromBuffer(value),
        ($0.Integration value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetIntegrationResponseRequest,
            $0.IntegrationResponse>(
        'GetIntegrationResponse',
        getIntegrationResponse_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetIntegrationResponseRequest.fromBuffer(value),
        ($0.IntegrationResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetMethodRequest, $0.Method>(
        'GetMethod',
        getMethod_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetMethodRequest.fromBuffer(value),
        ($0.Method value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetMethodResponseRequest, $0.MethodResponse>(
            'GetMethodResponse',
            getMethodResponse_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetMethodResponseRequest.fromBuffer(value),
            ($0.MethodResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetModelRequest, $0.Model>(
        'GetModel',
        getModel_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetModelRequest.fromBuffer(value),
        ($0.Model value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetModelsRequest, $0.Models>(
        'GetModels',
        getModels_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetModelsRequest.fromBuffer(value),
        ($0.Models value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetModelTemplateRequest, $0.Template>(
        'GetModelTemplate',
        getModelTemplate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetModelTemplateRequest.fromBuffer(value),
        ($0.Template value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetRequestValidatorRequest, $0.RequestValidator>(
            'GetRequestValidator',
            getRequestValidator_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetRequestValidatorRequest.fromBuffer(value),
            ($0.RequestValidator value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetRequestValidatorsRequest,
            $0.RequestValidators>(
        'GetRequestValidators',
        getRequestValidators_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetRequestValidatorsRequest.fromBuffer(value),
        ($0.RequestValidators value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetResourceRequest, $0.Resource>(
        'GetResource',
        getResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetResourceRequest.fromBuffer(value),
        ($0.Resource value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetResourcesRequest, $0.Resources>(
        'GetResources',
        getResources_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetResourcesRequest.fromBuffer(value),
        ($0.Resources value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetRestApiRequest, $0.RestApi>(
        'GetRestApi',
        getRestApi_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetRestApiRequest.fromBuffer(value),
        ($0.RestApi value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetRestApisRequest, $0.RestApis>(
        'GetRestApis',
        getRestApis_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetRestApisRequest.fromBuffer(value),
        ($0.RestApis value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetSdkRequest, $0.SdkResponse>(
        'GetSdk',
        getSdk_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetSdkRequest.fromBuffer(value),
        ($0.SdkResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetSdkTypeRequest, $0.SdkType>(
        'GetSdkType',
        getSdkType_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetSdkTypeRequest.fromBuffer(value),
        ($0.SdkType value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetSdkTypesRequest, $0.SdkTypes>(
        'GetSdkTypes',
        getSdkTypes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetSdkTypesRequest.fromBuffer(value),
        ($0.SdkTypes value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetStageRequest, $0.Stage>(
        'GetStage',
        getStage_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetStageRequest.fromBuffer(value),
        ($0.Stage value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetStagesRequest, $0.Stages>(
        'GetStages',
        getStages_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetStagesRequest.fromBuffer(value),
        ($0.Stages value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetTagsRequest, $0.Tags>(
        'GetTags',
        getTags_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetTagsRequest.fromBuffer(value),
        ($0.Tags value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetUsageRequest, $0.Usage>(
        'GetUsage',
        getUsage_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetUsageRequest.fromBuffer(value),
        ($0.Usage value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetUsagePlanRequest, $0.UsagePlan>(
        'GetUsagePlan',
        getUsagePlan_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetUsagePlanRequest.fromBuffer(value),
        ($0.UsagePlan value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetUsagePlanKeyRequest, $0.UsagePlanKey>(
        'GetUsagePlanKey',
        getUsagePlanKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetUsagePlanKeyRequest.fromBuffer(value),
        ($0.UsagePlanKey value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetUsagePlanKeysRequest, $0.UsagePlanKeys>(
            'GetUsagePlanKeys',
            getUsagePlanKeys_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetUsagePlanKeysRequest.fromBuffer(value),
            ($0.UsagePlanKeys value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetUsagePlansRequest, $0.UsagePlans>(
        'GetUsagePlans',
        getUsagePlans_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetUsagePlansRequest.fromBuffer(value),
        ($0.UsagePlans value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetVpcLinkRequest, $0.VpcLink>(
        'GetVpcLink',
        getVpcLink_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetVpcLinkRequest.fromBuffer(value),
        ($0.VpcLink value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetVpcLinksRequest, $0.VpcLinks>(
        'GetVpcLinks',
        getVpcLinks_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetVpcLinksRequest.fromBuffer(value),
        ($0.VpcLinks value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ImportApiKeysRequest, $0.ApiKeyIds>(
        'ImportApiKeys',
        importApiKeys_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ImportApiKeysRequest.fromBuffer(value),
        ($0.ApiKeyIds value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ImportDocumentationPartsRequest,
            $0.DocumentationPartIds>(
        'ImportDocumentationParts',
        importDocumentationParts_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ImportDocumentationPartsRequest.fromBuffer(value),
        ($0.DocumentationPartIds value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ImportRestApiRequest, $0.RestApi>(
        'ImportRestApi',
        importRestApi_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ImportRestApiRequest.fromBuffer(value),
        ($0.RestApi value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.PutGatewayResponseRequest, $0.GatewayResponse>(
            'PutGatewayResponse',
            putGatewayResponse_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.PutGatewayResponseRequest.fromBuffer(value),
            ($0.GatewayResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutIntegrationRequest, $0.Integration>(
        'PutIntegration',
        putIntegration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutIntegrationRequest.fromBuffer(value),
        ($0.Integration value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutIntegrationResponseRequest,
            $0.IntegrationResponse>(
        'PutIntegrationResponse',
        putIntegrationResponse_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutIntegrationResponseRequest.fromBuffer(value),
        ($0.IntegrationResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutMethodRequest, $0.Method>(
        'PutMethod',
        putMethod_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.PutMethodRequest.fromBuffer(value),
        ($0.Method value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.PutMethodResponseRequest, $0.MethodResponse>(
            'PutMethodResponse',
            putMethodResponse_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.PutMethodResponseRequest.fromBuffer(value),
            ($0.MethodResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutRestApiRequest, $0.RestApi>(
        'PutRestApi',
        putRestApi_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.PutRestApiRequest.fromBuffer(value),
        ($0.RestApi value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RejectDomainNameAccessAssociationRequest,
            $1.Empty>(
        'RejectDomainNameAccessAssociation',
        rejectDomainNameAccessAssociation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RejectDomainNameAccessAssociationRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TagResourceRequest, $1.Empty>(
        'TagResource',
        tagResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.TagResourceRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TestInvokeAuthorizerRequest,
            $0.TestInvokeAuthorizerResponse>(
        'TestInvokeAuthorizer',
        testInvokeAuthorizer_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.TestInvokeAuthorizerRequest.fromBuffer(value),
        ($0.TestInvokeAuthorizerResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TestInvokeMethodRequest,
            $0.TestInvokeMethodResponse>(
        'TestInvokeMethod',
        testInvokeMethod_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.TestInvokeMethodRequest.fromBuffer(value),
        ($0.TestInvokeMethodResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UntagResourceRequest, $1.Empty>(
        'UntagResource',
        untagResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UntagResourceRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateAccountRequest, $0.Account>(
        'UpdateAccount',
        updateAccount_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateAccountRequest.fromBuffer(value),
        ($0.Account value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateApiKeyRequest, $0.ApiKey>(
        'UpdateApiKey',
        updateApiKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateApiKeyRequest.fromBuffer(value),
        ($0.ApiKey value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateAuthorizerRequest, $0.Authorizer>(
        'UpdateAuthorizer',
        updateAuthorizer_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateAuthorizerRequest.fromBuffer(value),
        ($0.Authorizer value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateBasePathMappingRequest,
            $0.BasePathMapping>(
        'UpdateBasePathMapping',
        updateBasePathMapping_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateBasePathMappingRequest.fromBuffer(value),
        ($0.BasePathMapping value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateClientCertificateRequest,
            $0.ClientCertificate>(
        'UpdateClientCertificate',
        updateClientCertificate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateClientCertificateRequest.fromBuffer(value),
        ($0.ClientCertificate value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateDeploymentRequest, $0.Deployment>(
        'UpdateDeployment',
        updateDeployment_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateDeploymentRequest.fromBuffer(value),
        ($0.Deployment value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateDocumentationPartRequest,
            $0.DocumentationPart>(
        'UpdateDocumentationPart',
        updateDocumentationPart_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateDocumentationPartRequest.fromBuffer(value),
        ($0.DocumentationPart value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateDocumentationVersionRequest,
            $0.DocumentationVersion>(
        'UpdateDocumentationVersion',
        updateDocumentationVersion_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateDocumentationVersionRequest.fromBuffer(value),
        ($0.DocumentationVersion value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateDomainNameRequest, $0.DomainName>(
        'UpdateDomainName',
        updateDomainName_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateDomainNameRequest.fromBuffer(value),
        ($0.DomainName value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateGatewayResponseRequest,
            $0.GatewayResponse>(
        'UpdateGatewayResponse',
        updateGatewayResponse_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateGatewayResponseRequest.fromBuffer(value),
        ($0.GatewayResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateIntegrationRequest, $0.Integration>(
        'UpdateIntegration',
        updateIntegration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateIntegrationRequest.fromBuffer(value),
        ($0.Integration value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateIntegrationResponseRequest,
            $0.IntegrationResponse>(
        'UpdateIntegrationResponse',
        updateIntegrationResponse_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateIntegrationResponseRequest.fromBuffer(value),
        ($0.IntegrationResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateMethodRequest, $0.Method>(
        'UpdateMethod',
        updateMethod_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateMethodRequest.fromBuffer(value),
        ($0.Method value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateMethodResponseRequest, $0.MethodResponse>(
            'UpdateMethodResponse',
            updateMethodResponse_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateMethodResponseRequest.fromBuffer(value),
            ($0.MethodResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateModelRequest, $0.Model>(
        'UpdateModel',
        updateModel_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateModelRequest.fromBuffer(value),
        ($0.Model value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateRequestValidatorRequest,
            $0.RequestValidator>(
        'UpdateRequestValidator',
        updateRequestValidator_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateRequestValidatorRequest.fromBuffer(value),
        ($0.RequestValidator value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateResourceRequest, $0.Resource>(
        'UpdateResource',
        updateResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateResourceRequest.fromBuffer(value),
        ($0.Resource value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateRestApiRequest, $0.RestApi>(
        'UpdateRestApi',
        updateRestApi_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateRestApiRequest.fromBuffer(value),
        ($0.RestApi value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateStageRequest, $0.Stage>(
        'UpdateStage',
        updateStage_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateStageRequest.fromBuffer(value),
        ($0.Stage value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateUsageRequest, $0.Usage>(
        'UpdateUsage',
        updateUsage_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateUsageRequest.fromBuffer(value),
        ($0.Usage value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateUsagePlanRequest, $0.UsagePlan>(
        'UpdateUsagePlan',
        updateUsagePlan_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateUsagePlanRequest.fromBuffer(value),
        ($0.UsagePlan value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateVpcLinkRequest, $0.VpcLink>(
        'UpdateVpcLink',
        updateVpcLink_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateVpcLinkRequest.fromBuffer(value),
        ($0.VpcLink value) => value.writeToBuffer()));
  }

  $async.Future<$0.ApiKey> createApiKey_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateApiKeyRequest> $request) async {
    return createApiKey($call, await $request);
  }

  $async.Future<$0.ApiKey> createApiKey(
      $grpc.ServiceCall call, $0.CreateApiKeyRequest request);

  $async.Future<$0.Authorizer> createAuthorizer_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateAuthorizerRequest> $request) async {
    return createAuthorizer($call, await $request);
  }

  $async.Future<$0.Authorizer> createAuthorizer(
      $grpc.ServiceCall call, $0.CreateAuthorizerRequest request);

  $async.Future<$0.BasePathMapping> createBasePathMapping_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateBasePathMappingRequest> $request) async {
    return createBasePathMapping($call, await $request);
  }

  $async.Future<$0.BasePathMapping> createBasePathMapping(
      $grpc.ServiceCall call, $0.CreateBasePathMappingRequest request);

  $async.Future<$0.Deployment> createDeployment_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateDeploymentRequest> $request) async {
    return createDeployment($call, await $request);
  }

  $async.Future<$0.Deployment> createDeployment(
      $grpc.ServiceCall call, $0.CreateDeploymentRequest request);

  $async.Future<$0.DocumentationPart> createDocumentationPart_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateDocumentationPartRequest> $request) async {
    return createDocumentationPart($call, await $request);
  }

  $async.Future<$0.DocumentationPart> createDocumentationPart(
      $grpc.ServiceCall call, $0.CreateDocumentationPartRequest request);

  $async.Future<$0.DocumentationVersion> createDocumentationVersion_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateDocumentationVersionRequest> $request) async {
    return createDocumentationVersion($call, await $request);
  }

  $async.Future<$0.DocumentationVersion> createDocumentationVersion(
      $grpc.ServiceCall call, $0.CreateDocumentationVersionRequest request);

  $async.Future<$0.DomainName> createDomainName_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateDomainNameRequest> $request) async {
    return createDomainName($call, await $request);
  }

  $async.Future<$0.DomainName> createDomainName(
      $grpc.ServiceCall call, $0.CreateDomainNameRequest request);

  $async.Future<$0.DomainNameAccessAssociation>
      createDomainNameAccessAssociation_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.CreateDomainNameAccessAssociationRequest>
              $request) async {
    return createDomainNameAccessAssociation($call, await $request);
  }

  $async.Future<$0.DomainNameAccessAssociation>
      createDomainNameAccessAssociation($grpc.ServiceCall call,
          $0.CreateDomainNameAccessAssociationRequest request);

  $async.Future<$0.Model> createModel_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateModelRequest> $request) async {
    return createModel($call, await $request);
  }

  $async.Future<$0.Model> createModel(
      $grpc.ServiceCall call, $0.CreateModelRequest request);

  $async.Future<$0.RequestValidator> createRequestValidator_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateRequestValidatorRequest> $request) async {
    return createRequestValidator($call, await $request);
  }

  $async.Future<$0.RequestValidator> createRequestValidator(
      $grpc.ServiceCall call, $0.CreateRequestValidatorRequest request);

  $async.Future<$0.Resource> createResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateResourceRequest> $request) async {
    return createResource($call, await $request);
  }

  $async.Future<$0.Resource> createResource(
      $grpc.ServiceCall call, $0.CreateResourceRequest request);

  $async.Future<$0.RestApi> createRestApi_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateRestApiRequest> $request) async {
    return createRestApi($call, await $request);
  }

  $async.Future<$0.RestApi> createRestApi(
      $grpc.ServiceCall call, $0.CreateRestApiRequest request);

  $async.Future<$0.Stage> createStage_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateStageRequest> $request) async {
    return createStage($call, await $request);
  }

  $async.Future<$0.Stage> createStage(
      $grpc.ServiceCall call, $0.CreateStageRequest request);

  $async.Future<$0.UsagePlan> createUsagePlan_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateUsagePlanRequest> $request) async {
    return createUsagePlan($call, await $request);
  }

  $async.Future<$0.UsagePlan> createUsagePlan(
      $grpc.ServiceCall call, $0.CreateUsagePlanRequest request);

  $async.Future<$0.UsagePlanKey> createUsagePlanKey_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateUsagePlanKeyRequest> $request) async {
    return createUsagePlanKey($call, await $request);
  }

  $async.Future<$0.UsagePlanKey> createUsagePlanKey(
      $grpc.ServiceCall call, $0.CreateUsagePlanKeyRequest request);

  $async.Future<$0.VpcLink> createVpcLink_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateVpcLinkRequest> $request) async {
    return createVpcLink($call, await $request);
  }

  $async.Future<$0.VpcLink> createVpcLink(
      $grpc.ServiceCall call, $0.CreateVpcLinkRequest request);

  $async.Future<$1.Empty> deleteApiKey_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteApiKeyRequest> $request) async {
    return deleteApiKey($call, await $request);
  }

  $async.Future<$1.Empty> deleteApiKey(
      $grpc.ServiceCall call, $0.DeleteApiKeyRequest request);

  $async.Future<$1.Empty> deleteAuthorizer_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteAuthorizerRequest> $request) async {
    return deleteAuthorizer($call, await $request);
  }

  $async.Future<$1.Empty> deleteAuthorizer(
      $grpc.ServiceCall call, $0.DeleteAuthorizerRequest request);

  $async.Future<$1.Empty> deleteBasePathMapping_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteBasePathMappingRequest> $request) async {
    return deleteBasePathMapping($call, await $request);
  }

  $async.Future<$1.Empty> deleteBasePathMapping(
      $grpc.ServiceCall call, $0.DeleteBasePathMappingRequest request);

  $async.Future<$1.Empty> deleteClientCertificate_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteClientCertificateRequest> $request) async {
    return deleteClientCertificate($call, await $request);
  }

  $async.Future<$1.Empty> deleteClientCertificate(
      $grpc.ServiceCall call, $0.DeleteClientCertificateRequest request);

  $async.Future<$1.Empty> deleteDeployment_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteDeploymentRequest> $request) async {
    return deleteDeployment($call, await $request);
  }

  $async.Future<$1.Empty> deleteDeployment(
      $grpc.ServiceCall call, $0.DeleteDeploymentRequest request);

  $async.Future<$1.Empty> deleteDocumentationPart_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteDocumentationPartRequest> $request) async {
    return deleteDocumentationPart($call, await $request);
  }

  $async.Future<$1.Empty> deleteDocumentationPart(
      $grpc.ServiceCall call, $0.DeleteDocumentationPartRequest request);

  $async.Future<$1.Empty> deleteDocumentationVersion_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteDocumentationVersionRequest> $request) async {
    return deleteDocumentationVersion($call, await $request);
  }

  $async.Future<$1.Empty> deleteDocumentationVersion(
      $grpc.ServiceCall call, $0.DeleteDocumentationVersionRequest request);

  $async.Future<$1.Empty> deleteDomainName_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteDomainNameRequest> $request) async {
    return deleteDomainName($call, await $request);
  }

  $async.Future<$1.Empty> deleteDomainName(
      $grpc.ServiceCall call, $0.DeleteDomainNameRequest request);

  $async.Future<$1.Empty> deleteDomainNameAccessAssociation_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteDomainNameAccessAssociationRequest>
          $request) async {
    return deleteDomainNameAccessAssociation($call, await $request);
  }

  $async.Future<$1.Empty> deleteDomainNameAccessAssociation(
      $grpc.ServiceCall call,
      $0.DeleteDomainNameAccessAssociationRequest request);

  $async.Future<$1.Empty> deleteGatewayResponse_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteGatewayResponseRequest> $request) async {
    return deleteGatewayResponse($call, await $request);
  }

  $async.Future<$1.Empty> deleteGatewayResponse(
      $grpc.ServiceCall call, $0.DeleteGatewayResponseRequest request);

  $async.Future<$1.Empty> deleteIntegration_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteIntegrationRequest> $request) async {
    return deleteIntegration($call, await $request);
  }

  $async.Future<$1.Empty> deleteIntegration(
      $grpc.ServiceCall call, $0.DeleteIntegrationRequest request);

  $async.Future<$1.Empty> deleteIntegrationResponse_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteIntegrationResponseRequest> $request) async {
    return deleteIntegrationResponse($call, await $request);
  }

  $async.Future<$1.Empty> deleteIntegrationResponse(
      $grpc.ServiceCall call, $0.DeleteIntegrationResponseRequest request);

  $async.Future<$1.Empty> deleteMethod_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteMethodRequest> $request) async {
    return deleteMethod($call, await $request);
  }

  $async.Future<$1.Empty> deleteMethod(
      $grpc.ServiceCall call, $0.DeleteMethodRequest request);

  $async.Future<$1.Empty> deleteMethodResponse_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteMethodResponseRequest> $request) async {
    return deleteMethodResponse($call, await $request);
  }

  $async.Future<$1.Empty> deleteMethodResponse(
      $grpc.ServiceCall call, $0.DeleteMethodResponseRequest request);

  $async.Future<$1.Empty> deleteModel_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteModelRequest> $request) async {
    return deleteModel($call, await $request);
  }

  $async.Future<$1.Empty> deleteModel(
      $grpc.ServiceCall call, $0.DeleteModelRequest request);

  $async.Future<$1.Empty> deleteRequestValidator_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteRequestValidatorRequest> $request) async {
    return deleteRequestValidator($call, await $request);
  }

  $async.Future<$1.Empty> deleteRequestValidator(
      $grpc.ServiceCall call, $0.DeleteRequestValidatorRequest request);

  $async.Future<$1.Empty> deleteResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteResourceRequest> $request) async {
    return deleteResource($call, await $request);
  }

  $async.Future<$1.Empty> deleteResource(
      $grpc.ServiceCall call, $0.DeleteResourceRequest request);

  $async.Future<$1.Empty> deleteRestApi_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteRestApiRequest> $request) async {
    return deleteRestApi($call, await $request);
  }

  $async.Future<$1.Empty> deleteRestApi(
      $grpc.ServiceCall call, $0.DeleteRestApiRequest request);

  $async.Future<$1.Empty> deleteStage_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteStageRequest> $request) async {
    return deleteStage($call, await $request);
  }

  $async.Future<$1.Empty> deleteStage(
      $grpc.ServiceCall call, $0.DeleteStageRequest request);

  $async.Future<$1.Empty> deleteUsagePlan_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteUsagePlanRequest> $request) async {
    return deleteUsagePlan($call, await $request);
  }

  $async.Future<$1.Empty> deleteUsagePlan(
      $grpc.ServiceCall call, $0.DeleteUsagePlanRequest request);

  $async.Future<$1.Empty> deleteUsagePlanKey_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteUsagePlanKeyRequest> $request) async {
    return deleteUsagePlanKey($call, await $request);
  }

  $async.Future<$1.Empty> deleteUsagePlanKey(
      $grpc.ServiceCall call, $0.DeleteUsagePlanKeyRequest request);

  $async.Future<$1.Empty> deleteVpcLink_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteVpcLinkRequest> $request) async {
    return deleteVpcLink($call, await $request);
  }

  $async.Future<$1.Empty> deleteVpcLink(
      $grpc.ServiceCall call, $0.DeleteVpcLinkRequest request);

  $async.Future<$1.Empty> flushStageAuthorizersCache_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.FlushStageAuthorizersCacheRequest> $request) async {
    return flushStageAuthorizersCache($call, await $request);
  }

  $async.Future<$1.Empty> flushStageAuthorizersCache(
      $grpc.ServiceCall call, $0.FlushStageAuthorizersCacheRequest request);

  $async.Future<$1.Empty> flushStageCache_Pre($grpc.ServiceCall $call,
      $async.Future<$0.FlushStageCacheRequest> $request) async {
    return flushStageCache($call, await $request);
  }

  $async.Future<$1.Empty> flushStageCache(
      $grpc.ServiceCall call, $0.FlushStageCacheRequest request);

  $async.Future<$0.ClientCertificate> generateClientCertificate_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GenerateClientCertificateRequest> $request) async {
    return generateClientCertificate($call, await $request);
  }

  $async.Future<$0.ClientCertificate> generateClientCertificate(
      $grpc.ServiceCall call, $0.GenerateClientCertificateRequest request);

  $async.Future<$0.Account> getAccount_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetAccountRequest> $request) async {
    return getAccount($call, await $request);
  }

  $async.Future<$0.Account> getAccount(
      $grpc.ServiceCall call, $0.GetAccountRequest request);

  $async.Future<$0.ApiKey> getApiKey_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetApiKeyRequest> $request) async {
    return getApiKey($call, await $request);
  }

  $async.Future<$0.ApiKey> getApiKey(
      $grpc.ServiceCall call, $0.GetApiKeyRequest request);

  $async.Future<$0.ApiKeys> getApiKeys_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetApiKeysRequest> $request) async {
    return getApiKeys($call, await $request);
  }

  $async.Future<$0.ApiKeys> getApiKeys(
      $grpc.ServiceCall call, $0.GetApiKeysRequest request);

  $async.Future<$0.Authorizer> getAuthorizer_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetAuthorizerRequest> $request) async {
    return getAuthorizer($call, await $request);
  }

  $async.Future<$0.Authorizer> getAuthorizer(
      $grpc.ServiceCall call, $0.GetAuthorizerRequest request);

  $async.Future<$0.Authorizers> getAuthorizers_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetAuthorizersRequest> $request) async {
    return getAuthorizers($call, await $request);
  }

  $async.Future<$0.Authorizers> getAuthorizers(
      $grpc.ServiceCall call, $0.GetAuthorizersRequest request);

  $async.Future<$0.BasePathMapping> getBasePathMapping_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetBasePathMappingRequest> $request) async {
    return getBasePathMapping($call, await $request);
  }

  $async.Future<$0.BasePathMapping> getBasePathMapping(
      $grpc.ServiceCall call, $0.GetBasePathMappingRequest request);

  $async.Future<$0.BasePathMappings> getBasePathMappings_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetBasePathMappingsRequest> $request) async {
    return getBasePathMappings($call, await $request);
  }

  $async.Future<$0.BasePathMappings> getBasePathMappings(
      $grpc.ServiceCall call, $0.GetBasePathMappingsRequest request);

  $async.Future<$0.ClientCertificate> getClientCertificate_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetClientCertificateRequest> $request) async {
    return getClientCertificate($call, await $request);
  }

  $async.Future<$0.ClientCertificate> getClientCertificate(
      $grpc.ServiceCall call, $0.GetClientCertificateRequest request);

  $async.Future<$0.ClientCertificates> getClientCertificates_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetClientCertificatesRequest> $request) async {
    return getClientCertificates($call, await $request);
  }

  $async.Future<$0.ClientCertificates> getClientCertificates(
      $grpc.ServiceCall call, $0.GetClientCertificatesRequest request);

  $async.Future<$0.Deployment> getDeployment_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetDeploymentRequest> $request) async {
    return getDeployment($call, await $request);
  }

  $async.Future<$0.Deployment> getDeployment(
      $grpc.ServiceCall call, $0.GetDeploymentRequest request);

  $async.Future<$0.Deployments> getDeployments_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetDeploymentsRequest> $request) async {
    return getDeployments($call, await $request);
  }

  $async.Future<$0.Deployments> getDeployments(
      $grpc.ServiceCall call, $0.GetDeploymentsRequest request);

  $async.Future<$0.DocumentationPart> getDocumentationPart_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDocumentationPartRequest> $request) async {
    return getDocumentationPart($call, await $request);
  }

  $async.Future<$0.DocumentationPart> getDocumentationPart(
      $grpc.ServiceCall call, $0.GetDocumentationPartRequest request);

  $async.Future<$0.DocumentationParts> getDocumentationParts_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDocumentationPartsRequest> $request) async {
    return getDocumentationParts($call, await $request);
  }

  $async.Future<$0.DocumentationParts> getDocumentationParts(
      $grpc.ServiceCall call, $0.GetDocumentationPartsRequest request);

  $async.Future<$0.DocumentationVersion> getDocumentationVersion_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDocumentationVersionRequest> $request) async {
    return getDocumentationVersion($call, await $request);
  }

  $async.Future<$0.DocumentationVersion> getDocumentationVersion(
      $grpc.ServiceCall call, $0.GetDocumentationVersionRequest request);

  $async.Future<$0.DocumentationVersions> getDocumentationVersions_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDocumentationVersionsRequest> $request) async {
    return getDocumentationVersions($call, await $request);
  }

  $async.Future<$0.DocumentationVersions> getDocumentationVersions(
      $grpc.ServiceCall call, $0.GetDocumentationVersionsRequest request);

  $async.Future<$0.DomainName> getDomainName_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetDomainNameRequest> $request) async {
    return getDomainName($call, await $request);
  }

  $async.Future<$0.DomainName> getDomainName(
      $grpc.ServiceCall call, $0.GetDomainNameRequest request);

  $async.Future<$0.DomainNameAccessAssociations>
      getDomainNameAccessAssociations_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetDomainNameAccessAssociationsRequest>
              $request) async {
    return getDomainNameAccessAssociations($call, await $request);
  }

  $async.Future<$0.DomainNameAccessAssociations>
      getDomainNameAccessAssociations($grpc.ServiceCall call,
          $0.GetDomainNameAccessAssociationsRequest request);

  $async.Future<$0.DomainNames> getDomainNames_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetDomainNamesRequest> $request) async {
    return getDomainNames($call, await $request);
  }

  $async.Future<$0.DomainNames> getDomainNames(
      $grpc.ServiceCall call, $0.GetDomainNamesRequest request);

  $async.Future<$0.ExportResponse> getExport_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetExportRequest> $request) async {
    return getExport($call, await $request);
  }

  $async.Future<$0.ExportResponse> getExport(
      $grpc.ServiceCall call, $0.GetExportRequest request);

  $async.Future<$0.GatewayResponse> getGatewayResponse_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetGatewayResponseRequest> $request) async {
    return getGatewayResponse($call, await $request);
  }

  $async.Future<$0.GatewayResponse> getGatewayResponse(
      $grpc.ServiceCall call, $0.GetGatewayResponseRequest request);

  $async.Future<$0.GatewayResponses> getGatewayResponses_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetGatewayResponsesRequest> $request) async {
    return getGatewayResponses($call, await $request);
  }

  $async.Future<$0.GatewayResponses> getGatewayResponses(
      $grpc.ServiceCall call, $0.GetGatewayResponsesRequest request);

  $async.Future<$0.Integration> getIntegration_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetIntegrationRequest> $request) async {
    return getIntegration($call, await $request);
  }

  $async.Future<$0.Integration> getIntegration(
      $grpc.ServiceCall call, $0.GetIntegrationRequest request);

  $async.Future<$0.IntegrationResponse> getIntegrationResponse_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetIntegrationResponseRequest> $request) async {
    return getIntegrationResponse($call, await $request);
  }

  $async.Future<$0.IntegrationResponse> getIntegrationResponse(
      $grpc.ServiceCall call, $0.GetIntegrationResponseRequest request);

  $async.Future<$0.Method> getMethod_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetMethodRequest> $request) async {
    return getMethod($call, await $request);
  }

  $async.Future<$0.Method> getMethod(
      $grpc.ServiceCall call, $0.GetMethodRequest request);

  $async.Future<$0.MethodResponse> getMethodResponse_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetMethodResponseRequest> $request) async {
    return getMethodResponse($call, await $request);
  }

  $async.Future<$0.MethodResponse> getMethodResponse(
      $grpc.ServiceCall call, $0.GetMethodResponseRequest request);

  $async.Future<$0.Model> getModel_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetModelRequest> $request) async {
    return getModel($call, await $request);
  }

  $async.Future<$0.Model> getModel(
      $grpc.ServiceCall call, $0.GetModelRequest request);

  $async.Future<$0.Models> getModels_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetModelsRequest> $request) async {
    return getModels($call, await $request);
  }

  $async.Future<$0.Models> getModels(
      $grpc.ServiceCall call, $0.GetModelsRequest request);

  $async.Future<$0.Template> getModelTemplate_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetModelTemplateRequest> $request) async {
    return getModelTemplate($call, await $request);
  }

  $async.Future<$0.Template> getModelTemplate(
      $grpc.ServiceCall call, $0.GetModelTemplateRequest request);

  $async.Future<$0.RequestValidator> getRequestValidator_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetRequestValidatorRequest> $request) async {
    return getRequestValidator($call, await $request);
  }

  $async.Future<$0.RequestValidator> getRequestValidator(
      $grpc.ServiceCall call, $0.GetRequestValidatorRequest request);

  $async.Future<$0.RequestValidators> getRequestValidators_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetRequestValidatorsRequest> $request) async {
    return getRequestValidators($call, await $request);
  }

  $async.Future<$0.RequestValidators> getRequestValidators(
      $grpc.ServiceCall call, $0.GetRequestValidatorsRequest request);

  $async.Future<$0.Resource> getResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetResourceRequest> $request) async {
    return getResource($call, await $request);
  }

  $async.Future<$0.Resource> getResource(
      $grpc.ServiceCall call, $0.GetResourceRequest request);

  $async.Future<$0.Resources> getResources_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetResourcesRequest> $request) async {
    return getResources($call, await $request);
  }

  $async.Future<$0.Resources> getResources(
      $grpc.ServiceCall call, $0.GetResourcesRequest request);

  $async.Future<$0.RestApi> getRestApi_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetRestApiRequest> $request) async {
    return getRestApi($call, await $request);
  }

  $async.Future<$0.RestApi> getRestApi(
      $grpc.ServiceCall call, $0.GetRestApiRequest request);

  $async.Future<$0.RestApis> getRestApis_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetRestApisRequest> $request) async {
    return getRestApis($call, await $request);
  }

  $async.Future<$0.RestApis> getRestApis(
      $grpc.ServiceCall call, $0.GetRestApisRequest request);

  $async.Future<$0.SdkResponse> getSdk_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.GetSdkRequest> $request) async {
    return getSdk($call, await $request);
  }

  $async.Future<$0.SdkResponse> getSdk(
      $grpc.ServiceCall call, $0.GetSdkRequest request);

  $async.Future<$0.SdkType> getSdkType_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetSdkTypeRequest> $request) async {
    return getSdkType($call, await $request);
  }

  $async.Future<$0.SdkType> getSdkType(
      $grpc.ServiceCall call, $0.GetSdkTypeRequest request);

  $async.Future<$0.SdkTypes> getSdkTypes_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetSdkTypesRequest> $request) async {
    return getSdkTypes($call, await $request);
  }

  $async.Future<$0.SdkTypes> getSdkTypes(
      $grpc.ServiceCall call, $0.GetSdkTypesRequest request);

  $async.Future<$0.Stage> getStage_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetStageRequest> $request) async {
    return getStage($call, await $request);
  }

  $async.Future<$0.Stage> getStage(
      $grpc.ServiceCall call, $0.GetStageRequest request);

  $async.Future<$0.Stages> getStages_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetStagesRequest> $request) async {
    return getStages($call, await $request);
  }

  $async.Future<$0.Stages> getStages(
      $grpc.ServiceCall call, $0.GetStagesRequest request);

  $async.Future<$0.Tags> getTags_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetTagsRequest> $request) async {
    return getTags($call, await $request);
  }

  $async.Future<$0.Tags> getTags(
      $grpc.ServiceCall call, $0.GetTagsRequest request);

  $async.Future<$0.Usage> getUsage_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetUsageRequest> $request) async {
    return getUsage($call, await $request);
  }

  $async.Future<$0.Usage> getUsage(
      $grpc.ServiceCall call, $0.GetUsageRequest request);

  $async.Future<$0.UsagePlan> getUsagePlan_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetUsagePlanRequest> $request) async {
    return getUsagePlan($call, await $request);
  }

  $async.Future<$0.UsagePlan> getUsagePlan(
      $grpc.ServiceCall call, $0.GetUsagePlanRequest request);

  $async.Future<$0.UsagePlanKey> getUsagePlanKey_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetUsagePlanKeyRequest> $request) async {
    return getUsagePlanKey($call, await $request);
  }

  $async.Future<$0.UsagePlanKey> getUsagePlanKey(
      $grpc.ServiceCall call, $0.GetUsagePlanKeyRequest request);

  $async.Future<$0.UsagePlanKeys> getUsagePlanKeys_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetUsagePlanKeysRequest> $request) async {
    return getUsagePlanKeys($call, await $request);
  }

  $async.Future<$0.UsagePlanKeys> getUsagePlanKeys(
      $grpc.ServiceCall call, $0.GetUsagePlanKeysRequest request);

  $async.Future<$0.UsagePlans> getUsagePlans_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetUsagePlansRequest> $request) async {
    return getUsagePlans($call, await $request);
  }

  $async.Future<$0.UsagePlans> getUsagePlans(
      $grpc.ServiceCall call, $0.GetUsagePlansRequest request);

  $async.Future<$0.VpcLink> getVpcLink_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetVpcLinkRequest> $request) async {
    return getVpcLink($call, await $request);
  }

  $async.Future<$0.VpcLink> getVpcLink(
      $grpc.ServiceCall call, $0.GetVpcLinkRequest request);

  $async.Future<$0.VpcLinks> getVpcLinks_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetVpcLinksRequest> $request) async {
    return getVpcLinks($call, await $request);
  }

  $async.Future<$0.VpcLinks> getVpcLinks(
      $grpc.ServiceCall call, $0.GetVpcLinksRequest request);

  $async.Future<$0.ApiKeyIds> importApiKeys_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ImportApiKeysRequest> $request) async {
    return importApiKeys($call, await $request);
  }

  $async.Future<$0.ApiKeyIds> importApiKeys(
      $grpc.ServiceCall call, $0.ImportApiKeysRequest request);

  $async.Future<$0.DocumentationPartIds> importDocumentationParts_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ImportDocumentationPartsRequest> $request) async {
    return importDocumentationParts($call, await $request);
  }

  $async.Future<$0.DocumentationPartIds> importDocumentationParts(
      $grpc.ServiceCall call, $0.ImportDocumentationPartsRequest request);

  $async.Future<$0.RestApi> importRestApi_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ImportRestApiRequest> $request) async {
    return importRestApi($call, await $request);
  }

  $async.Future<$0.RestApi> importRestApi(
      $grpc.ServiceCall call, $0.ImportRestApiRequest request);

  $async.Future<$0.GatewayResponse> putGatewayResponse_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutGatewayResponseRequest> $request) async {
    return putGatewayResponse($call, await $request);
  }

  $async.Future<$0.GatewayResponse> putGatewayResponse(
      $grpc.ServiceCall call, $0.PutGatewayResponseRequest request);

  $async.Future<$0.Integration> putIntegration_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutIntegrationRequest> $request) async {
    return putIntegration($call, await $request);
  }

  $async.Future<$0.Integration> putIntegration(
      $grpc.ServiceCall call, $0.PutIntegrationRequest request);

  $async.Future<$0.IntegrationResponse> putIntegrationResponse_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutIntegrationResponseRequest> $request) async {
    return putIntegrationResponse($call, await $request);
  }

  $async.Future<$0.IntegrationResponse> putIntegrationResponse(
      $grpc.ServiceCall call, $0.PutIntegrationResponseRequest request);

  $async.Future<$0.Method> putMethod_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutMethodRequest> $request) async {
    return putMethod($call, await $request);
  }

  $async.Future<$0.Method> putMethod(
      $grpc.ServiceCall call, $0.PutMethodRequest request);

  $async.Future<$0.MethodResponse> putMethodResponse_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutMethodResponseRequest> $request) async {
    return putMethodResponse($call, await $request);
  }

  $async.Future<$0.MethodResponse> putMethodResponse(
      $grpc.ServiceCall call, $0.PutMethodResponseRequest request);

  $async.Future<$0.RestApi> putRestApi_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutRestApiRequest> $request) async {
    return putRestApi($call, await $request);
  }

  $async.Future<$0.RestApi> putRestApi(
      $grpc.ServiceCall call, $0.PutRestApiRequest request);

  $async.Future<$1.Empty> rejectDomainNameAccessAssociation_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RejectDomainNameAccessAssociationRequest>
          $request) async {
    return rejectDomainNameAccessAssociation($call, await $request);
  }

  $async.Future<$1.Empty> rejectDomainNameAccessAssociation(
      $grpc.ServiceCall call,
      $0.RejectDomainNameAccessAssociationRequest request);

  $async.Future<$1.Empty> tagResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagResourceRequest> $request) async {
    return tagResource($call, await $request);
  }

  $async.Future<$1.Empty> tagResource(
      $grpc.ServiceCall call, $0.TagResourceRequest request);

  $async.Future<$0.TestInvokeAuthorizerResponse> testInvokeAuthorizer_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.TestInvokeAuthorizerRequest> $request) async {
    return testInvokeAuthorizer($call, await $request);
  }

  $async.Future<$0.TestInvokeAuthorizerResponse> testInvokeAuthorizer(
      $grpc.ServiceCall call, $0.TestInvokeAuthorizerRequest request);

  $async.Future<$0.TestInvokeMethodResponse> testInvokeMethod_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.TestInvokeMethodRequest> $request) async {
    return testInvokeMethod($call, await $request);
  }

  $async.Future<$0.TestInvokeMethodResponse> testInvokeMethod(
      $grpc.ServiceCall call, $0.TestInvokeMethodRequest request);

  $async.Future<$1.Empty> untagResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UntagResourceRequest> $request) async {
    return untagResource($call, await $request);
  }

  $async.Future<$1.Empty> untagResource(
      $grpc.ServiceCall call, $0.UntagResourceRequest request);

  $async.Future<$0.Account> updateAccount_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateAccountRequest> $request) async {
    return updateAccount($call, await $request);
  }

  $async.Future<$0.Account> updateAccount(
      $grpc.ServiceCall call, $0.UpdateAccountRequest request);

  $async.Future<$0.ApiKey> updateApiKey_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateApiKeyRequest> $request) async {
    return updateApiKey($call, await $request);
  }

  $async.Future<$0.ApiKey> updateApiKey(
      $grpc.ServiceCall call, $0.UpdateApiKeyRequest request);

  $async.Future<$0.Authorizer> updateAuthorizer_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateAuthorizerRequest> $request) async {
    return updateAuthorizer($call, await $request);
  }

  $async.Future<$0.Authorizer> updateAuthorizer(
      $grpc.ServiceCall call, $0.UpdateAuthorizerRequest request);

  $async.Future<$0.BasePathMapping> updateBasePathMapping_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateBasePathMappingRequest> $request) async {
    return updateBasePathMapping($call, await $request);
  }

  $async.Future<$0.BasePathMapping> updateBasePathMapping(
      $grpc.ServiceCall call, $0.UpdateBasePathMappingRequest request);

  $async.Future<$0.ClientCertificate> updateClientCertificate_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateClientCertificateRequest> $request) async {
    return updateClientCertificate($call, await $request);
  }

  $async.Future<$0.ClientCertificate> updateClientCertificate(
      $grpc.ServiceCall call, $0.UpdateClientCertificateRequest request);

  $async.Future<$0.Deployment> updateDeployment_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateDeploymentRequest> $request) async {
    return updateDeployment($call, await $request);
  }

  $async.Future<$0.Deployment> updateDeployment(
      $grpc.ServiceCall call, $0.UpdateDeploymentRequest request);

  $async.Future<$0.DocumentationPart> updateDocumentationPart_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateDocumentationPartRequest> $request) async {
    return updateDocumentationPart($call, await $request);
  }

  $async.Future<$0.DocumentationPart> updateDocumentationPart(
      $grpc.ServiceCall call, $0.UpdateDocumentationPartRequest request);

  $async.Future<$0.DocumentationVersion> updateDocumentationVersion_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateDocumentationVersionRequest> $request) async {
    return updateDocumentationVersion($call, await $request);
  }

  $async.Future<$0.DocumentationVersion> updateDocumentationVersion(
      $grpc.ServiceCall call, $0.UpdateDocumentationVersionRequest request);

  $async.Future<$0.DomainName> updateDomainName_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateDomainNameRequest> $request) async {
    return updateDomainName($call, await $request);
  }

  $async.Future<$0.DomainName> updateDomainName(
      $grpc.ServiceCall call, $0.UpdateDomainNameRequest request);

  $async.Future<$0.GatewayResponse> updateGatewayResponse_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateGatewayResponseRequest> $request) async {
    return updateGatewayResponse($call, await $request);
  }

  $async.Future<$0.GatewayResponse> updateGatewayResponse(
      $grpc.ServiceCall call, $0.UpdateGatewayResponseRequest request);

  $async.Future<$0.Integration> updateIntegration_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateIntegrationRequest> $request) async {
    return updateIntegration($call, await $request);
  }

  $async.Future<$0.Integration> updateIntegration(
      $grpc.ServiceCall call, $0.UpdateIntegrationRequest request);

  $async.Future<$0.IntegrationResponse> updateIntegrationResponse_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateIntegrationResponseRequest> $request) async {
    return updateIntegrationResponse($call, await $request);
  }

  $async.Future<$0.IntegrationResponse> updateIntegrationResponse(
      $grpc.ServiceCall call, $0.UpdateIntegrationResponseRequest request);

  $async.Future<$0.Method> updateMethod_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateMethodRequest> $request) async {
    return updateMethod($call, await $request);
  }

  $async.Future<$0.Method> updateMethod(
      $grpc.ServiceCall call, $0.UpdateMethodRequest request);

  $async.Future<$0.MethodResponse> updateMethodResponse_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateMethodResponseRequest> $request) async {
    return updateMethodResponse($call, await $request);
  }

  $async.Future<$0.MethodResponse> updateMethodResponse(
      $grpc.ServiceCall call, $0.UpdateMethodResponseRequest request);

  $async.Future<$0.Model> updateModel_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateModelRequest> $request) async {
    return updateModel($call, await $request);
  }

  $async.Future<$0.Model> updateModel(
      $grpc.ServiceCall call, $0.UpdateModelRequest request);

  $async.Future<$0.RequestValidator> updateRequestValidator_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateRequestValidatorRequest> $request) async {
    return updateRequestValidator($call, await $request);
  }

  $async.Future<$0.RequestValidator> updateRequestValidator(
      $grpc.ServiceCall call, $0.UpdateRequestValidatorRequest request);

  $async.Future<$0.Resource> updateResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateResourceRequest> $request) async {
    return updateResource($call, await $request);
  }

  $async.Future<$0.Resource> updateResource(
      $grpc.ServiceCall call, $0.UpdateResourceRequest request);

  $async.Future<$0.RestApi> updateRestApi_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateRestApiRequest> $request) async {
    return updateRestApi($call, await $request);
  }

  $async.Future<$0.RestApi> updateRestApi(
      $grpc.ServiceCall call, $0.UpdateRestApiRequest request);

  $async.Future<$0.Stage> updateStage_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateStageRequest> $request) async {
    return updateStage($call, await $request);
  }

  $async.Future<$0.Stage> updateStage(
      $grpc.ServiceCall call, $0.UpdateStageRequest request);

  $async.Future<$0.Usage> updateUsage_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateUsageRequest> $request) async {
    return updateUsage($call, await $request);
  }

  $async.Future<$0.Usage> updateUsage(
      $grpc.ServiceCall call, $0.UpdateUsageRequest request);

  $async.Future<$0.UsagePlan> updateUsagePlan_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateUsagePlanRequest> $request) async {
    return updateUsagePlan($call, await $request);
  }

  $async.Future<$0.UsagePlan> updateUsagePlan(
      $grpc.ServiceCall call, $0.UpdateUsagePlanRequest request);

  $async.Future<$0.VpcLink> updateVpcLink_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateVpcLinkRequest> $request) async {
    return updateVpcLink($call, await $request);
  }

  $async.Future<$0.VpcLink> updateVpcLink(
      $grpc.ServiceCall call, $0.UpdateVpcLinkRequest request);
}

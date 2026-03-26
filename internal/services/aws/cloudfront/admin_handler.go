package cloudfront

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/cloudfront"
	cloudfrontconnect "vorpalstacks/internal/pb/aws/cloudfront/cloudfrontconnect"
	"vorpalstacks/internal/pb/aws/common"
)

// AdminHandler provides CloudFront service administration functionality.
// It implements the CloudFrontServiceHandler interface for gRPC-Web communication.
type AdminHandler struct{}

var _ cloudfrontconnect.CloudFrontServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new CloudFront AdminHandler.
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// ListDistributions lists all CloudFront distributions in the account.
// It returns a list of distribution summaries with their configurations and statuses.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListDistributions(ctx context.Context, req *connect.Request[pb.ListDistributionsRequest]) (*connect.Response[pb.ListDistributionsResult], error) {
	return connect.NewResponse(&pb.ListDistributionsResult{
		Distributionlist: &pb.DistributionList{
			Quantity:    0,
			Items:       []*pb.DistributionSummary{},
			Istruncated: false,
		},
	}), nil
}

// AssociateAlias associates an alias (alternate domain name) with a CloudFront distribution.
// This allows you to use your own domain name with your CloudFront distribution.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) AssociateAlias(ctx context.Context, req *connect.Request[pb.AssociateAliasRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.AssociateAlias is not implemented"))
}

// AssociateDistributionTenantWebACL associates a web ACL with a distribution tenant.
// This controls access to the distribution using AWS WAF rules.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) AssociateDistributionTenantWebACL(ctx context.Context, req *connect.Request[pb.AssociateDistributionTenantWebACLRequest]) (*connect.Response[pb.AssociateDistributionTenantWebACLResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.AssociateDistributionTenantWebACL is not implemented"))
}

// AssociateDistributionWebACL associates a web ACL with a CloudFront distribution.
// This controls access to the distribution using AWS WAF rules.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) AssociateDistributionWebACL(ctx context.Context, req *connect.Request[pb.AssociateDistributionWebACLRequest]) (*connect.Response[pb.AssociateDistributionWebACLResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.AssociateDistributionWebACL is not implemented"))
}

// CopyDistribution creates a copy of an existing CloudFront distribution.
// This allows you to create a new distribution with the same configuration.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CopyDistribution(ctx context.Context, req *connect.Request[pb.CopyDistributionRequest]) (*connect.Response[pb.CopyDistributionResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CopyDistribution is not implemented"))
}

// CreateAnycastIpList creates a list of IP addresses for Anycast CloudFront.
// Anycast allows multiple servers to serve content from the same IP address.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateAnycastIpList(ctx context.Context, req *connect.Request[pb.CreateAnycastIpListRequest]) (*connect.Response[pb.CreateAnycastIpListResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreateAnycastIpList is not implemented"))
}

// CreateCachePolicy creates a cache policy for CloudFront distributions.
// Cache policies control how CloudFront caches content at edge locations.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateCachePolicy(ctx context.Context, req *connect.Request[pb.CreateCachePolicyRequest]) (*connect.Response[pb.CreateCachePolicyResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreateCachePolicy is not implemented"))
}

// CreateCloudFrontOriginAccessIdentity creates an origin access identity for CloudFront.
// This is used to restrict access to your S3 bucket or other origins.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateCloudFrontOriginAccessIdentity(ctx context.Context, req *connect.Request[pb.CreateCloudFrontOriginAccessIdentityRequest]) (*connect.Response[pb.CreateCloudFrontOriginAccessIdentityResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreateCloudFrontOriginAccessIdentity is not implemented"))
}

// CreateConnectionFunction creates a connection function for CloudFront.
// Connection functions allow you to control connections to your origins.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateConnectionFunction(ctx context.Context, req *connect.Request[pb.CreateConnectionFunctionRequest]) (*connect.Response[pb.CreateConnectionFunctionResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreateConnectionFunction is not implemented"))
}

// CreateConnectionGroup creates a connection group for CloudFront.
// Connection groups allow you to manage multiple connections to origins.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateConnectionGroup(ctx context.Context, req *connect.Request[pb.CreateConnectionGroupRequest]) (*connect.Response[pb.CreateConnectionGroupResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreateConnectionGroup is not implemented"))
}

// CreateContinuousDeploymentPolicy creates a continuous deployment policy.
// This allows you to test changes to your distribution with a percentage of traffic.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateContinuousDeploymentPolicy(ctx context.Context, req *connect.Request[pb.CreateContinuousDeploymentPolicyRequest]) (*connect.Response[pb.CreateContinuousDeploymentPolicyResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreateContinuousDeploymentPolicy is not implemented"))
}

// CreateDistribution creates a new CloudFront distribution.
// Distributions serve content from one or more origins.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateDistribution(ctx context.Context, req *connect.Request[pb.CreateDistributionRequest]) (*connect.Response[pb.CreateDistributionResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreateDistribution is not implemented"))
}

// CreateDistributionTenant creates a distribution tenant for CloudFront.
// Distribution tenants allow you to create multi-tenant configurations.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateDistributionTenant(ctx context.Context, req *connect.Request[pb.CreateDistributionTenantRequest]) (*connect.Response[pb.CreateDistributionTenantResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreateDistributionTenant is not implemented"))
}

// CreateDistributionWithTags creates a CloudFront distribution with tags.
// Tags are key-value pairs used for billing and access control purposes.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateDistributionWithTags(ctx context.Context, req *connect.Request[pb.CreateDistributionWithTagsRequest]) (*connect.Response[pb.CreateDistributionWithTagsResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreateDistributionWithTags is not implemented"))
}

// CreateFieldLevelEncryptionConfig creates a field-level encryption configuration.
// This encrypts specific fields in the request body before forwarding to the origin.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateFieldLevelEncryptionConfig(ctx context.Context, req *connect.Request[pb.CreateFieldLevelEncryptionConfigRequest]) (*connect.Response[pb.CreateFieldLevelEncryptionConfigResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreateFieldLevelEncryptionConfig is not implemented"))
}

// CreateFieldLevelEncryptionProfile creates a field-level encryption profile.
// This defines the fields to encrypt and the encryption key to use.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateFieldLevelEncryptionProfile(ctx context.Context, req *connect.Request[pb.CreateFieldLevelEncryptionProfileRequest]) (*connect.Response[pb.CreateFieldLevelEncryptionProfileResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreateFieldLevelEncryptionProfile is not implemented"))
}

// CreateFunction creates a CloudFront function.
// Functions are lightweight scripts that can modify CloudFront behavior at edge locations.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateFunction(ctx context.Context, req *connect.Request[pb.CreateFunctionRequest]) (*connect.Response[pb.CreateFunctionResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreateFunction is not implemented"))
}

// CreateInvalidation creates an invalidation for a CloudFront distribution.
// Invalidations remove cached objects from CloudFront edge locations.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateInvalidation(ctx context.Context, req *connect.Request[pb.CreateInvalidationRequest]) (*connect.Response[pb.CreateInvalidationResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreateInvalidation is not implemented"))
}

// CreateInvalidationForDistributionTenant creates an invalidation for a distribution tenant.
// This removes cached objects for a specific tenant's distribution.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateInvalidationForDistributionTenant(ctx context.Context, req *connect.Request[pb.CreateInvalidationForDistributionTenantRequest]) (*connect.Response[pb.CreateInvalidationForDistributionTenantResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreateInvalidationForDistributionTenant is not implemented"))
}

// CreateKeyGroup creates a key group for CloudFront.
// Key groups are used to specify which public keys can be used for signed URLs or cookies.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateKeyGroup(ctx context.Context, req *connect.Request[pb.CreateKeyGroupRequest]) (*connect.Response[pb.CreateKeyGroupResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreateKeyGroup is not implemented"))
}

// CreateKeyValueStore creates a key value store for CloudFront.
// Key value stores are used to store data for CloudFront functions at edge locations.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateKeyValueStore(ctx context.Context, req *connect.Request[pb.CreateKeyValueStoreRequest]) (*connect.Response[pb.CreateKeyValueStoreResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreateKeyValueStore is not implemented"))
}

// CreateMonitoringSubscription creates a monitoring subscription for CloudFront.
// This enables additional monitoring and alerting for your distributions.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateMonitoringSubscription(ctx context.Context, req *connect.Request[pb.CreateMonitoringSubscriptionRequest]) (*connect.Response[pb.CreateMonitoringSubscriptionResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreateMonitoringSubscription is not implemented"))
}

// CreateOriginAccessControl creates an origin access control for CloudFront.
// This controls how CloudFront accesses your origin servers.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateOriginAccessControl(ctx context.Context, req *connect.Request[pb.CreateOriginAccessControlRequest]) (*connect.Response[pb.CreateOriginAccessControlResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreateOriginAccessControl is not implemented"))
}

// CreateOriginRequestPolicy creates an origin request policy for CloudFront.
// This controls what headers, cookies, and query strings are forwarded to the origin.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateOriginRequestPolicy(ctx context.Context, req *connect.Request[pb.CreateOriginRequestPolicyRequest]) (*connect.Response[pb.CreateOriginRequestPolicyResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreateOriginRequestPolicy is not implemented"))
}

// CreatePublicKey creates a public key for CloudFront.
// Public keys are used for signed URLs and signed cookies.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreatePublicKey(ctx context.Context, req *connect.Request[pb.CreatePublicKeyRequest]) (*connect.Response[pb.CreatePublicKeyResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreatePublicKey is not implemented"))
}

// CreateRealtimeLogConfig creates a real-time log configuration for CloudFront.
// Real-time logs provide detailed information about requests as they pass through CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateRealtimeLogConfig(ctx context.Context, req *connect.Request[pb.CreateRealtimeLogConfigRequest]) (*connect.Response[pb.CreateRealtimeLogConfigResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreateRealtimeLogConfig is not implemented"))
}

// CreateResponseHeadersPolicy creates a response headers policy for CloudFront.
// This defines which HTTP headers CloudFront includes in responses.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateResponseHeadersPolicy(ctx context.Context, req *connect.Request[pb.CreateResponseHeadersPolicyRequest]) (*connect.Response[pb.CreateResponseHeadersPolicyResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreateResponseHeadersPolicy is not implemented"))
}

// CreateStreamingDistribution creates a streaming distribution for CloudFront.
// Streaming distributions are used to serve live streaming content.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateStreamingDistribution(ctx context.Context, req *connect.Request[pb.CreateStreamingDistributionRequest]) (*connect.Response[pb.CreateStreamingDistributionResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreateStreamingDistribution is not implemented"))
}

// CreateStreamingDistributionWithTags creates a streaming distribution with tags.
// Tags are key-value pairs used for billing and access control purposes.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateStreamingDistributionWithTags(ctx context.Context, req *connect.Request[pb.CreateStreamingDistributionWithTagsRequest]) (*connect.Response[pb.CreateStreamingDistributionWithTagsResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreateStreamingDistributionWithTags is not implemented"))
}

// CreateTrustStore creates a trust store for CloudFront.
// Trust stores are used to verify client certificates for mutual TLS.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateTrustStore(ctx context.Context, req *connect.Request[pb.CreateTrustStoreRequest]) (*connect.Response[pb.CreateTrustStoreResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreateTrustStore is not implemented"))
}

// CreateVpcOrigin creates a VPC origin for CloudFront.
// VPC origins allow CloudFront to access your origins in a VPC.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateVpcOrigin(ctx context.Context, req *connect.Request[pb.CreateVpcOriginRequest]) (*connect.Response[pb.CreateVpcOriginResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.CreateVpcOrigin is not implemented"))
}

// DeleteAnycastIpList deletes an Anycast IP list for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteAnycastIpList(ctx context.Context, req *connect.Request[pb.DeleteAnycastIpListRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DeleteAnycastIpList is not implemented"))
}

// DeleteCachePolicy deletes a cache policy for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteCachePolicy(ctx context.Context, req *connect.Request[pb.DeleteCachePolicyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DeleteCachePolicy is not implemented"))
}

// DeleteCloudFrontOriginAccessIdentity deletes a CloudFront origin access identity.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteCloudFrontOriginAccessIdentity(ctx context.Context, req *connect.Request[pb.DeleteCloudFrontOriginAccessIdentityRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DeleteCloudFrontOriginAccessIdentity is not implemented"))
}

// DeleteConnectionFunction deletes a connection function for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteConnectionFunction(ctx context.Context, req *connect.Request[pb.DeleteConnectionFunctionRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DeleteConnectionFunction is not implemented"))
}

// DeleteConnectionGroup deletes a connection group for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteConnectionGroup(ctx context.Context, req *connect.Request[pb.DeleteConnectionGroupRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DeleteConnectionGroup is not implemented"))
}

// DeleteContinuousDeploymentPolicy deletes a continuous deployment policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteContinuousDeploymentPolicy(ctx context.Context, req *connect.Request[pb.DeleteContinuousDeploymentPolicyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DeleteContinuousDeploymentPolicy is not implemented"))
}

// DeleteDistribution deletes a CloudFront distribution.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteDistribution(ctx context.Context, req *connect.Request[pb.DeleteDistributionRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DeleteDistribution is not implemented"))
}

// DeleteDistributionTenant deletes a distribution tenant for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteDistributionTenant(ctx context.Context, req *connect.Request[pb.DeleteDistributionTenantRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DeleteDistributionTenant is not implemented"))
}

// DeleteFieldLevelEncryptionConfig deletes a field-level encryption configuration.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteFieldLevelEncryptionConfig(ctx context.Context, req *connect.Request[pb.DeleteFieldLevelEncryptionConfigRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DeleteFieldLevelEncryptionConfig is not implemented"))
}

// DeleteFieldLevelEncryptionProfile deletes a field-level encryption profile.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteFieldLevelEncryptionProfile(ctx context.Context, req *connect.Request[pb.DeleteFieldLevelEncryptionProfileRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DeleteFieldLevelEncryptionProfile is not implemented"))
}

// DeleteFunction deletes a CloudFront function.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteFunction(ctx context.Context, req *connect.Request[pb.DeleteFunctionRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DeleteFunction is not implemented"))
}

// DeleteKeyGroup deletes a key group for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteKeyGroup(ctx context.Context, req *connect.Request[pb.DeleteKeyGroupRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DeleteKeyGroup is not implemented"))
}

// DeleteKeyValueStore deletes a key value store for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteKeyValueStore(ctx context.Context, req *connect.Request[pb.DeleteKeyValueStoreRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DeleteKeyValueStore is not implemented"))
}

// DeleteMonitoringSubscription deletes a monitoring subscription for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteMonitoringSubscription(ctx context.Context, req *connect.Request[pb.DeleteMonitoringSubscriptionRequest]) (*connect.Response[pb.DeleteMonitoringSubscriptionResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DeleteMonitoringSubscription is not implemented"))
}

// DeleteOriginAccessControl deletes an origin access control for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteOriginAccessControl(ctx context.Context, req *connect.Request[pb.DeleteOriginAccessControlRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DeleteOriginAccessControl is not implemented"))
}

// DeleteOriginRequestPolicy deletes an origin request policy for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteOriginRequestPolicy(ctx context.Context, req *connect.Request[pb.DeleteOriginRequestPolicyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DeleteOriginRequestPolicy is not implemented"))
}

// DeletePublicKey deletes a public key for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeletePublicKey(ctx context.Context, req *connect.Request[pb.DeletePublicKeyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DeletePublicKey is not implemented"))
}

// DeleteRealtimeLogConfig deletes a real-time log configuration for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteRealtimeLogConfig(ctx context.Context, req *connect.Request[pb.DeleteRealtimeLogConfigRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DeleteRealtimeLogConfig is not implemented"))
}

// DeleteResourcePolicy deletes a resource policy for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteResourcePolicy(ctx context.Context, req *connect.Request[pb.DeleteResourcePolicyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DeleteResourcePolicy is not implemented"))
}

// DeleteResponseHeadersPolicy deletes a response headers policy for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteResponseHeadersPolicy(ctx context.Context, req *connect.Request[pb.DeleteResponseHeadersPolicyRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DeleteResponseHeadersPolicy is not implemented"))
}

// DeleteStreamingDistribution deletes a streaming distribution for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteStreamingDistribution(ctx context.Context, req *connect.Request[pb.DeleteStreamingDistributionRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DeleteStreamingDistribution is not implemented"))
}

// DeleteTrustStore deletes a trust store for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteTrustStore(ctx context.Context, req *connect.Request[pb.DeleteTrustStoreRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DeleteTrustStore is not implemented"))
}

// DeleteVpcOrigin deletes a VPC origin for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteVpcOrigin(ctx context.Context, req *connect.Request[pb.DeleteVpcOriginRequest]) (*connect.Response[pb.DeleteVpcOriginResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DeleteVpcOrigin is not implemented"))
}

// DescribeConnectionFunction describes a connection function for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeConnectionFunction(ctx context.Context, req *connect.Request[pb.DescribeConnectionFunctionRequest]) (*connect.Response[pb.DescribeConnectionFunctionResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DescribeConnectionFunction is not implemented"))
}

// DescribeFunction describes a CloudFront function.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeFunction(ctx context.Context, req *connect.Request[pb.DescribeFunctionRequest]) (*connect.Response[pb.DescribeFunctionResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DescribeFunction is not implemented"))
}

// DescribeKeyValueStore describes a key value store for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DescribeKeyValueStore(ctx context.Context, req *connect.Request[pb.DescribeKeyValueStoreRequest]) (*connect.Response[pb.DescribeKeyValueStoreResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DescribeKeyValueStore is not implemented"))
}

// DisassociateDistributionTenantWebACL disassociates a web ACL from a distribution tenant.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DisassociateDistributionTenantWebACL(ctx context.Context, req *connect.Request[pb.DisassociateDistributionTenantWebACLRequest]) (*connect.Response[pb.DisassociateDistributionTenantWebACLResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DisassociateDistributionTenantWebACL is not implemented"))
}

// DisassociateDistributionWebACL disassociates a web ACL from a CloudFront distribution.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DisassociateDistributionWebACL(ctx context.Context, req *connect.Request[pb.DisassociateDistributionWebACLRequest]) (*connect.Response[pb.DisassociateDistributionWebACLResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.DisassociateDistributionWebACL is not implemented"))
}

// GetAnycastIpList gets an Anycast IP list for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetAnycastIpList(ctx context.Context, req *connect.Request[pb.GetAnycastIpListRequest]) (*connect.Response[pb.GetAnycastIpListResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetAnycastIpList is not implemented"))
}

// GetCachePolicy gets a cache policy for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetCachePolicy(ctx context.Context, req *connect.Request[pb.GetCachePolicyRequest]) (*connect.Response[pb.GetCachePolicyResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetCachePolicy is not implemented"))
}

// GetCachePolicyConfig gets a cache policy configuration for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetCachePolicyConfig(ctx context.Context, req *connect.Request[pb.GetCachePolicyConfigRequest]) (*connect.Response[pb.GetCachePolicyConfigResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetCachePolicyConfig is not implemented"))
}

// GetCloudFrontOriginAccessIdentity gets a CloudFront origin access identity.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetCloudFrontOriginAccessIdentity(ctx context.Context, req *connect.Request[pb.GetCloudFrontOriginAccessIdentityRequest]) (*connect.Response[pb.GetCloudFrontOriginAccessIdentityResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetCloudFrontOriginAccessIdentity is not implemented"))
}

// GetCloudFrontOriginAccessIdentityConfig gets a CloudFront origin access identity configuration.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetCloudFrontOriginAccessIdentityConfig(ctx context.Context, req *connect.Request[pb.GetCloudFrontOriginAccessIdentityConfigRequest]) (*connect.Response[pb.GetCloudFrontOriginAccessIdentityConfigResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetCloudFrontOriginAccessIdentityConfig is not implemented"))
}

// GetConnectionFunction gets a connection function for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetConnectionFunction(ctx context.Context, req *connect.Request[pb.GetConnectionFunctionRequest]) (*connect.Response[pb.GetConnectionFunctionResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetConnectionFunction is not implemented"))
}

// GetConnectionGroup gets a connection group for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetConnectionGroup(ctx context.Context, req *connect.Request[pb.GetConnectionGroupRequest]) (*connect.Response[pb.GetConnectionGroupResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetConnectionGroup is not implemented"))
}

// GetConnectionGroupByRoutingEndpoint gets a connection group by routing endpoint for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetConnectionGroupByRoutingEndpoint(ctx context.Context, req *connect.Request[pb.GetConnectionGroupByRoutingEndpointRequest]) (*connect.Response[pb.GetConnectionGroupByRoutingEndpointResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetConnectionGroupByRoutingEndpoint is not implemented"))
}

// GetContinuousDeploymentPolicy gets a continuous deployment policy for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetContinuousDeploymentPolicy(ctx context.Context, req *connect.Request[pb.GetContinuousDeploymentPolicyRequest]) (*connect.Response[pb.GetContinuousDeploymentPolicyResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetContinuousDeploymentPolicy is not implemented"))
}

// GetContinuousDeploymentPolicyConfig gets a continuous deployment policy configuration for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetContinuousDeploymentPolicyConfig(ctx context.Context, req *connect.Request[pb.GetContinuousDeploymentPolicyConfigRequest]) (*connect.Response[pb.GetContinuousDeploymentPolicyConfigResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetContinuousDeploymentPolicyConfig is not implemented"))
}

// GetDistribution gets a CloudFront distribution.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetDistribution(ctx context.Context, req *connect.Request[pb.GetDistributionRequest]) (*connect.Response[pb.GetDistributionResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetDistribution is not implemented"))
}

// GetDistributionConfig gets a CloudFront distribution configuration.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetDistributionConfig(ctx context.Context, req *connect.Request[pb.GetDistributionConfigRequest]) (*connect.Response[pb.GetDistributionConfigResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetDistributionConfig is not implemented"))
}

// GetDistributionTenant gets a distribution tenant for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetDistributionTenant(ctx context.Context, req *connect.Request[pb.GetDistributionTenantRequest]) (*connect.Response[pb.GetDistributionTenantResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetDistributionTenant is not implemented"))
}

// GetDistributionTenantByDomain gets a distribution tenant by domain for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetDistributionTenantByDomain(ctx context.Context, req *connect.Request[pb.GetDistributionTenantByDomainRequest]) (*connect.Response[pb.GetDistributionTenantByDomainResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetDistributionTenantByDomain is not implemented"))
}

// GetFieldLevelEncryption gets field-level encryption for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetFieldLevelEncryption(ctx context.Context, req *connect.Request[pb.GetFieldLevelEncryptionRequest]) (*connect.Response[pb.GetFieldLevelEncryptionResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetFieldLevelEncryption is not implemented"))
}

// GetFieldLevelEncryptionConfig gets field-level encryption configuration for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetFieldLevelEncryptionConfig(ctx context.Context, req *connect.Request[pb.GetFieldLevelEncryptionConfigRequest]) (*connect.Response[pb.GetFieldLevelEncryptionConfigResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetFieldLevelEncryptionConfig is not implemented"))
}

// GetFieldLevelEncryptionProfile gets a field-level encryption profile for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetFieldLevelEncryptionProfile(ctx context.Context, req *connect.Request[pb.GetFieldLevelEncryptionProfileRequest]) (*connect.Response[pb.GetFieldLevelEncryptionProfileResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetFieldLevelEncryptionProfile is not implemented"))
}

// GetFieldLevelEncryptionProfileConfig gets a field-level encryption profile configuration for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetFieldLevelEncryptionProfileConfig(ctx context.Context, req *connect.Request[pb.GetFieldLevelEncryptionProfileConfigRequest]) (*connect.Response[pb.GetFieldLevelEncryptionProfileConfigResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetFieldLevelEncryptionProfileConfig is not implemented"))
}

// GetFunction gets a CloudFront function.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetFunction(ctx context.Context, req *connect.Request[pb.GetFunctionRequest]) (*connect.Response[pb.GetFunctionResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetFunction is not implemented"))
}

// GetInvalidation gets an invalidation for a CloudFront distribution.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetInvalidation(ctx context.Context, req *connect.Request[pb.GetInvalidationRequest]) (*connect.Response[pb.GetInvalidationResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetInvalidation is not implemented"))
}

// GetInvalidationForDistributionTenant gets an invalidation for a distribution tenant.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetInvalidationForDistributionTenant(ctx context.Context, req *connect.Request[pb.GetInvalidationForDistributionTenantRequest]) (*connect.Response[pb.GetInvalidationForDistributionTenantResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetInvalidationForDistributionTenant is not implemented"))
}

// GetKeyGroup gets a key group for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetKeyGroup(ctx context.Context, req *connect.Request[pb.GetKeyGroupRequest]) (*connect.Response[pb.GetKeyGroupResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetKeyGroup is not implemented"))
}

// GetKeyGroupConfig gets a key group configuration for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetKeyGroupConfig(ctx context.Context, req *connect.Request[pb.GetKeyGroupConfigRequest]) (*connect.Response[pb.GetKeyGroupConfigResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetKeyGroupConfig is not implemented"))
}

// GetManagedCertificateDetails gets managed certificate details for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetManagedCertificateDetails(ctx context.Context, req *connect.Request[pb.GetManagedCertificateDetailsRequest]) (*connect.Response[pb.GetManagedCertificateDetailsResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetManagedCertificateDetails is not implemented"))
}

// GetMonitoringSubscription gets a monitoring subscription for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetMonitoringSubscription(ctx context.Context, req *connect.Request[pb.GetMonitoringSubscriptionRequest]) (*connect.Response[pb.GetMonitoringSubscriptionResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetMonitoringSubscription is not implemented"))
}

// GetOriginAccessControl gets an origin access control for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetOriginAccessControl(ctx context.Context, req *connect.Request[pb.GetOriginAccessControlRequest]) (*connect.Response[pb.GetOriginAccessControlResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetOriginAccessControl is not implemented"))
}

// GetOriginAccessControlConfig gets an origin access control configuration for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetOriginAccessControlConfig(ctx context.Context, req *connect.Request[pb.GetOriginAccessControlConfigRequest]) (*connect.Response[pb.GetOriginAccessControlConfigResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetOriginAccessControlConfig is not implemented"))
}

// GetOriginRequestPolicy gets an origin request policy for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetOriginRequestPolicy(ctx context.Context, req *connect.Request[pb.GetOriginRequestPolicyRequest]) (*connect.Response[pb.GetOriginRequestPolicyResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetOriginRequestPolicy is not implemented"))
}

// GetOriginRequestPolicyConfig gets an origin request policy configuration for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetOriginRequestPolicyConfig(ctx context.Context, req *connect.Request[pb.GetOriginRequestPolicyConfigRequest]) (*connect.Response[pb.GetOriginRequestPolicyConfigResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetOriginRequestPolicyConfig is not implemented"))
}

// GetPublicKey gets a public key for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetPublicKey(ctx context.Context, req *connect.Request[pb.GetPublicKeyRequest]) (*connect.Response[pb.GetPublicKeyResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetPublicKey is not implemented"))
}

// GetPublicKeyConfig gets a public key configuration for CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetPublicKeyConfig(ctx context.Context, req *connect.Request[pb.GetPublicKeyConfigRequest]) (*connect.Response[pb.GetPublicKeyConfigResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetPublicKeyConfig is not implemented"))
}

// GetRealtimeLogConfig gets a real-time log configuration for CloudFront.
// Real-time logs provide detailed information about requests as they pass through CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetRealtimeLogConfig(ctx context.Context, req *connect.Request[pb.GetRealtimeLogConfigRequest]) (*connect.Response[pb.GetRealtimeLogConfigResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetRealtimeLogConfig is not implemented"))
}

// GetResourcePolicy gets a resource policy for CloudFront.
// Resource policies control access to CloudFront resources.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetResourcePolicy(ctx context.Context, req *connect.Request[pb.GetResourcePolicyRequest]) (*connect.Response[pb.GetResourcePolicyResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetResourcePolicy is not implemented"))
}

// GetResponseHeadersPolicy gets a response headers policy for CloudFront.
// Response headers policies define which HTTP headers are included in responses.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetResponseHeadersPolicy(ctx context.Context, req *connect.Request[pb.GetResponseHeadersPolicyRequest]) (*connect.Response[pb.GetResponseHeadersPolicyResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetResponseHeadersPolicy is not implemented"))
}

// GetResponseHeadersPolicyConfig gets a response headers policy configuration for CloudFront.
// This returns the detailed configuration of the response headers policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetResponseHeadersPolicyConfig(ctx context.Context, req *connect.Request[pb.GetResponseHeadersPolicyConfigRequest]) (*connect.Response[pb.GetResponseHeadersPolicyConfigResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetResponseHeadersPolicyConfig is not implemented"))
}

// GetStreamingDistribution gets a streaming distribution for CloudFront.
// Streaming distributions are used to serve live streaming content.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetStreamingDistribution(ctx context.Context, req *connect.Request[pb.GetStreamingDistributionRequest]) (*connect.Response[pb.GetStreamingDistributionResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetStreamingDistribution is not implemented"))
}

// GetStreamingDistributionConfig gets a streaming distribution configuration for CloudFront.
// This returns the detailed configuration of the streaming distribution.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetStreamingDistributionConfig(ctx context.Context, req *connect.Request[pb.GetStreamingDistributionConfigRequest]) (*connect.Response[pb.GetStreamingDistributionConfigResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetStreamingDistributionConfig is not implemented"))
}

// GetTrustStore gets a trust store for CloudFront.
// Trust stores are used to verify client certificates for mutual TLS.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetTrustStore(ctx context.Context, req *connect.Request[pb.GetTrustStoreRequest]) (*connect.Response[pb.GetTrustStoreResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetTrustStore is not implemented"))
}

// GetVpcOrigin gets a VPC origin for CloudFront.
// VPC origins allow CloudFront to access your origins in a VPC.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetVpcOrigin(ctx context.Context, req *connect.Request[pb.GetVpcOriginRequest]) (*connect.Response[pb.GetVpcOriginResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.GetVpcOrigin is not implemented"))
}

// ListAnycastIpLists lists all Anycast IP lists for CloudFront.
// Anycast allows multiple servers to serve content from the same IP address.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListAnycastIpLists(ctx context.Context, req *connect.Request[pb.ListAnycastIpListsRequest]) (*connect.Response[pb.ListAnycastIpListsResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListAnycastIpLists is not implemented"))
}

// ListCachePolicies lists all cache policies for CloudFront.
// Cache policies control how CloudFront caches content at edge locations.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListCachePolicies(ctx context.Context, req *connect.Request[pb.ListCachePoliciesRequest]) (*connect.Response[pb.ListCachePoliciesResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListCachePolicies is not implemented"))
}

// ListCloudFrontOriginAccessIdentities lists all CloudFront origin access identities.
// Origin access identities are used to restrict access to S3 buckets or other origins.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListCloudFrontOriginAccessIdentities(ctx context.Context, req *connect.Request[pb.ListCloudFrontOriginAccessIdentitiesRequest]) (*connect.Response[pb.ListCloudFrontOriginAccessIdentitiesResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListCloudFrontOriginAccessIdentities is not implemented"))
}

// ListConflictingAliases lists all conflicting aliases for CloudFront.
// Conflicting aliases occur when multiple distributions use the same alternate domain name.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListConflictingAliases(ctx context.Context, req *connect.Request[pb.ListConflictingAliasesRequest]) (*connect.Response[pb.ListConflictingAliasesResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListConflictingAliases is not implemented"))
}

// ListConnectionFunctions lists all connection functions for CloudFront.
// Connection functions allow you to control connections to your origins.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListConnectionFunctions(ctx context.Context, req *connect.Request[pb.ListConnectionFunctionsRequest]) (*connect.Response[pb.ListConnectionFunctionsResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListConnectionFunctions is not implemented"))
}

// ListConnectionGroups lists all connection groups for CloudFront.
// Connection groups allow you to manage multiple connections to origins.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListConnectionGroups(ctx context.Context, req *connect.Request[pb.ListConnectionGroupsRequest]) (*connect.Response[pb.ListConnectionGroupsResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListConnectionGroups is not implemented"))
}

// ListContinuousDeploymentPolicies lists all continuous deployment policies for CloudFront.
// Continuous deployment policies allow you to test changes with a percentage of traffic.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListContinuousDeploymentPolicies(ctx context.Context, req *connect.Request[pb.ListContinuousDeploymentPoliciesRequest]) (*connect.Response[pb.ListContinuousDeploymentPoliciesResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListContinuousDeploymentPolicies is not implemented"))
}

// ListDistributionsByAnycastIpListId lists all distributions associated with an Anycast IP list.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListDistributionsByAnycastIpListId(ctx context.Context, req *connect.Request[pb.ListDistributionsByAnycastIpListIdRequest]) (*connect.Response[pb.ListDistributionsByAnycastIpListIdResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListDistributionsByAnycastIpListId is not implemented"))
}

// ListDistributionsByCachePolicyId lists all distributions associated with a cache policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListDistributionsByCachePolicyId(ctx context.Context, req *connect.Request[pb.ListDistributionsByCachePolicyIdRequest]) (*connect.Response[pb.ListDistributionsByCachePolicyIdResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListDistributionsByCachePolicyId is not implemented"))
}

// ListDistributionsByConnectionFunction lists all distributions associated with a connection function.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListDistributionsByConnectionFunction(ctx context.Context, req *connect.Request[pb.ListDistributionsByConnectionFunctionRequest]) (*connect.Response[pb.ListDistributionsByConnectionFunctionResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListDistributionsByConnectionFunction is not implemented"))
}

// ListDistributionsByConnectionMode lists all distributions by connection mode.
// Connection modes define how CloudFront connects to your origins.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListDistributionsByConnectionMode(ctx context.Context, req *connect.Request[pb.ListDistributionsByConnectionModeRequest]) (*connect.Response[pb.ListDistributionsByConnectionModeResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListDistributionsByConnectionMode is not implemented"))
}

// ListDistributionsByKeyGroup lists all distributions associated with a key group.
// Key groups are used to specify which public keys can be used for signed URLs or cookies.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListDistributionsByKeyGroup(ctx context.Context, req *connect.Request[pb.ListDistributionsByKeyGroupRequest]) (*connect.Response[pb.ListDistributionsByKeyGroupResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListDistributionsByKeyGroup is not implemented"))
}

// ListDistributionsByOriginRequestPolicyId lists all distributions associated with an origin request policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListDistributionsByOriginRequestPolicyId(ctx context.Context, req *connect.Request[pb.ListDistributionsByOriginRequestPolicyIdRequest]) (*connect.Response[pb.ListDistributionsByOriginRequestPolicyIdResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListDistributionsByOriginRequestPolicyId is not implemented"))
}

// ListDistributionsByOwnedResource lists all distributions by owned resource.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListDistributionsByOwnedResource(ctx context.Context, req *connect.Request[pb.ListDistributionsByOwnedResourceRequest]) (*connect.Response[pb.ListDistributionsByOwnedResourceResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListDistributionsByOwnedResource is not implemented"))
}

// ListDistributionsByRealtimeLogConfig lists all distributions associated with a real-time log configuration.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListDistributionsByRealtimeLogConfig(ctx context.Context, req *connect.Request[pb.ListDistributionsByRealtimeLogConfigRequest]) (*connect.Response[pb.ListDistributionsByRealtimeLogConfigResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListDistributionsByRealtimeLogConfig is not implemented"))
}

// ListDistributionsByResponseHeadersPolicyId lists all distributions associated with a response headers policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListDistributionsByResponseHeadersPolicyId(ctx context.Context, req *connect.Request[pb.ListDistributionsByResponseHeadersPolicyIdRequest]) (*connect.Response[pb.ListDistributionsByResponseHeadersPolicyIdResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListDistributionsByResponseHeadersPolicyId is not implemented"))
}

// ListDistributionsByTrustStore lists all distributions associated with a trust store.
// Trust stores are used to verify client certificates for mutual TLS.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListDistributionsByTrustStore(ctx context.Context, req *connect.Request[pb.ListDistributionsByTrustStoreRequest]) (*connect.Response[pb.ListDistributionsByTrustStoreResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListDistributionsByTrustStore is not implemented"))
}

// ListDistributionsByVpcOriginId lists all distributions associated with a VPC origin.
// VPC origins allow CloudFront to access your origins in a VPC.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListDistributionsByVpcOriginId(ctx context.Context, req *connect.Request[pb.ListDistributionsByVpcOriginIdRequest]) (*connect.Response[pb.ListDistributionsByVpcOriginIdResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListDistributionsByVpcOriginId is not implemented"))
}

// ListDistributionsByWebACLId lists all distributions associated with a web ACL.
// Web ACLs control access to your distributions using AWS WAF rules.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListDistributionsByWebACLId(ctx context.Context, req *connect.Request[pb.ListDistributionsByWebACLIdRequest]) (*connect.Response[pb.ListDistributionsByWebACLIdResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListDistributionsByWebACLId is not implemented"))
}

// ListDistributionTenants lists all distribution tenants for CloudFront.
// Distribution tenants allow you to create multi-tenant configurations.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListDistributionTenants(ctx context.Context, req *connect.Request[pb.ListDistributionTenantsRequest]) (*connect.Response[pb.ListDistributionTenantsResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListDistributionTenants is not implemented"))
}

// ListDistributionTenantsByCustomization lists all distribution tenants by customisation.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListDistributionTenantsByCustomization(ctx context.Context, req *connect.Request[pb.ListDistributionTenantsByCustomizationRequest]) (*connect.Response[pb.ListDistributionTenantsByCustomizationResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListDistributionTenantsByCustomization is not implemented"))
}

// ListDomainConflicts lists all conflicting domain names for CloudFront.
// Conflicting domains occur when multiple distributions use the same alternate domain name.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListDomainConflicts(ctx context.Context, req *connect.Request[pb.ListDomainConflictsRequest]) (*connect.Response[pb.ListDomainConflictsResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListDomainConflicts is not implemented"))
}

// ListFieldLevelEncryptionConfigs lists all field-level encryption configurations for CloudFront.
// Field-level encryption encrypts specific fields in the request body before forwarding to the origin.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListFieldLevelEncryptionConfigs(ctx context.Context, req *connect.Request[pb.ListFieldLevelEncryptionConfigsRequest]) (*connect.Response[pb.ListFieldLevelEncryptionConfigsResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListFieldLevelEncryptionConfigs is not implemented"))
}

// ListFieldLevelEncryptionProfiles lists all field-level encryption profiles for CloudFront.
// Field-level encryption profiles define which fields to encrypt and the encryption key to use.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListFieldLevelEncryptionProfiles(ctx context.Context, req *connect.Request[pb.ListFieldLevelEncryptionProfilesRequest]) (*connect.Response[pb.ListFieldLevelEncryptionProfilesResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListFieldLevelEncryptionProfiles is not implemented"))
}

// ListFunctions lists all CloudFront functions.
// Functions are lightweight scripts that can modify CloudFront behaviour at edge locations.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListFunctions(ctx context.Context, req *connect.Request[pb.ListFunctionsRequest]) (*connect.Response[pb.ListFunctionsResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListFunctions is not implemented"))
}

// ListInvalidations lists all invalidations for a CloudFront distribution.
// Invalidations remove cached objects from CloudFront edge locations.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListInvalidations(ctx context.Context, req *connect.Request[pb.ListInvalidationsRequest]) (*connect.Response[pb.ListInvalidationsResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListInvalidations is not implemented"))
}

// ListInvalidationsForDistributionTenant lists all invalidations for a distribution tenant.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListInvalidationsForDistributionTenant(ctx context.Context, req *connect.Request[pb.ListInvalidationsForDistributionTenantRequest]) (*connect.Response[pb.ListInvalidationsForDistributionTenantResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListInvalidationsForDistributionTenant is not implemented"))
}

// ListKeyGroups lists all key groups for CloudFront.
// Key groups are used to specify which public keys can be used for signed URLs or cookies.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListKeyGroups(ctx context.Context, req *connect.Request[pb.ListKeyGroupsRequest]) (*connect.Response[pb.ListKeyGroupsResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListKeyGroups is not implemented"))
}

// ListKeyValueStores lists all key value stores for CloudFront.
// Key value stores are used to store data for CloudFront functions at edge locations.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListKeyValueStores(ctx context.Context, req *connect.Request[pb.ListKeyValueStoresRequest]) (*connect.Response[pb.ListKeyValueStoresResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListKeyValueStores is not implemented"))
}

// ListOriginAccessControls lists all origin access controls for CloudFront.
// Origin access controls control how CloudFront accesses your origin servers.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListOriginAccessControls(ctx context.Context, req *connect.Request[pb.ListOriginAccessControlsRequest]) (*connect.Response[pb.ListOriginAccessControlsResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListOriginAccessControls is not implemented"))
}

// ListOriginRequestPolicies lists all origin request policies for CloudFront.
// Origin request policies control what headers, cookies, and query strings are forwarded to the origin.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListOriginRequestPolicies(ctx context.Context, req *connect.Request[pb.ListOriginRequestPoliciesRequest]) (*connect.Response[pb.ListOriginRequestPoliciesResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListOriginRequestPolicies is not implemented"))
}

// ListPublicKeys lists all public keys for CloudFront.
// Public keys are used for signed URLs and signed cookies.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListPublicKeys(ctx context.Context, req *connect.Request[pb.ListPublicKeysRequest]) (*connect.Response[pb.ListPublicKeysResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListPublicKeys is not implemented"))
}

// ListRealtimeLogConfigs lists all real-time log configurations for CloudFront.
// Real-time logs provide detailed information about requests as they pass through CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListRealtimeLogConfigs(ctx context.Context, req *connect.Request[pb.ListRealtimeLogConfigsRequest]) (*connect.Response[pb.ListRealtimeLogConfigsResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListRealtimeLogConfigs is not implemented"))
}

// ListResponseHeadersPolicies lists all response headers policies for CloudFront.
// Response headers policies define which HTTP headers are included in responses.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListResponseHeadersPolicies(ctx context.Context, req *connect.Request[pb.ListResponseHeadersPoliciesRequest]) (*connect.Response[pb.ListResponseHeadersPoliciesResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListResponseHeadersPolicies is not implemented"))
}

// ListStreamingDistributions lists all streaming distributions for CloudFront.
// Streaming distributions are used to serve live streaming content.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListStreamingDistributions(ctx context.Context, req *connect.Request[pb.ListStreamingDistributionsRequest]) (*connect.Response[pb.ListStreamingDistributionsResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListStreamingDistributions is not implemented"))
}

// ListTagsForResource lists all tags for a CloudFront resource.
// Tags are key-value pairs used for billing and access control purposes.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListTagsForResource(ctx context.Context, req *connect.Request[pb.ListTagsForResourceRequest]) (*connect.Response[pb.ListTagsForResourceResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListTagsForResource is not implemented"))
}

// ListTrustStores lists all trust stores for CloudFront.
// Trust stores are used to verify client certificates for mutual TLS.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListTrustStores(ctx context.Context, req *connect.Request[pb.ListTrustStoresRequest]) (*connect.Response[pb.ListTrustStoresResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListTrustStores is not implemented"))
}

// ListVpcOrigins lists all VPC origins for CloudFront.
// VPC origins allow CloudFront to access your origins in a VPC.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListVpcOrigins(ctx context.Context, req *connect.Request[pb.ListVpcOriginsRequest]) (*connect.Response[pb.ListVpcOriginsResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.ListVpcOrigins is not implemented"))
}

// PublishConnectionFunction publishes a connection function for CloudFront.
// Publishing makes the function available for use in CloudFront distributions.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PublishConnectionFunction(ctx context.Context, req *connect.Request[pb.PublishConnectionFunctionRequest]) (*connect.Response[pb.PublishConnectionFunctionResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.PublishConnectionFunction is not implemented"))
}

// PublishFunction publishes a CloudFront function.
// Publishing makes the function available for use in CloudFront distributions.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PublishFunction(ctx context.Context, req *connect.Request[pb.PublishFunctionRequest]) (*connect.Response[pb.PublishFunctionResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.PublishFunction is not implemented"))
}

// PutResourcePolicy creates or updates a resource policy for CloudFront.
// Resource policies control access to CloudFront resources.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) PutResourcePolicy(ctx context.Context, req *connect.Request[pb.PutResourcePolicyRequest]) (*connect.Response[pb.PutResourcePolicyResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.PutResourcePolicy is not implemented"))
}

// TagResource adds tags to a CloudFront resource.
// Tags are key-value pairs used for billing and access control purposes.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TagResource(ctx context.Context, req *connect.Request[pb.TagResourceRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.TagResource is not implemented"))
}

// TestConnectionFunction tests a connection function for CloudFront.
// Testing validates the function's behaviour before publishing.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TestConnectionFunction(ctx context.Context, req *connect.Request[pb.TestConnectionFunctionRequest]) (*connect.Response[pb.TestConnectionFunctionResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.TestConnectionFunction is not implemented"))
}

// TestFunction tests a CloudFront function.
// Testing validates the function's behaviour before publishing.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TestFunction(ctx context.Context, req *connect.Request[pb.TestFunctionRequest]) (*connect.Response[pb.TestFunctionResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.TestFunction is not implemented"))
}

// UntagResource removes tags from a CloudFront resource.
// Tags are key-value pairs used for billing and access control purposes.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UntagResource(ctx context.Context, req *connect.Request[pb.UntagResourceRequest]) (*connect.Response[common.Empty], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.UntagResource is not implemented"))
}

// UpdateAnycastIpList updates an Anycast IP list for CloudFront.
// Anycast allows multiple servers to serve content from the same IP address.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateAnycastIpList(ctx context.Context, req *connect.Request[pb.UpdateAnycastIpListRequest]) (*connect.Response[pb.UpdateAnycastIpListResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.UpdateAnycastIpList is not implemented"))
}

// UpdateCachePolicy updates a cache policy for CloudFront.
// Cache policies control how CloudFront caches content at edge locations.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateCachePolicy(ctx context.Context, req *connect.Request[pb.UpdateCachePolicyRequest]) (*connect.Response[pb.UpdateCachePolicyResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.UpdateCachePolicy is not implemented"))
}

// UpdateCloudFrontOriginAccessIdentity updates a CloudFront origin access identity.
// Origin access identities are used to restrict access to S3 buckets or other origins.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateCloudFrontOriginAccessIdentity(ctx context.Context, req *connect.Request[pb.UpdateCloudFrontOriginAccessIdentityRequest]) (*connect.Response[pb.UpdateCloudFrontOriginAccessIdentityResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.UpdateCloudFrontOriginAccessIdentity is not implemented"))
}

// UpdateConnectionFunction updates a connection function for CloudFront.
// Connection functions allow you to control connections to your origins.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateConnectionFunction(ctx context.Context, req *connect.Request[pb.UpdateConnectionFunctionRequest]) (*connect.Response[pb.UpdateConnectionFunctionResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.UpdateConnectionFunction is not implemented"))
}

// UpdateConnectionGroup updates a connection group for CloudFront.
// Connection groups allow you to manage multiple connections to origins.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateConnectionGroup(ctx context.Context, req *connect.Request[pb.UpdateConnectionGroupRequest]) (*connect.Response[pb.UpdateConnectionGroupResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.UpdateConnectionGroup is not implemented"))
}

// UpdateContinuousDeploymentPolicy updates a continuous deployment policy.
// Continuous deployment policies allow you to test changes with a percentage of traffic.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateContinuousDeploymentPolicy(ctx context.Context, req *connect.Request[pb.UpdateContinuousDeploymentPolicyRequest]) (*connect.Response[pb.UpdateContinuousDeploymentPolicyResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.UpdateContinuousDeploymentPolicy is not implemented"))
}

// UpdateDistribution updates a CloudFront distribution.
// Distributions serve content from one or more origins.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateDistribution(ctx context.Context, req *connect.Request[pb.UpdateDistributionRequest]) (*connect.Response[pb.UpdateDistributionResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.UpdateDistribution is not implemented"))
}

// UpdateDistributionTenant updates a distribution tenant for CloudFront.
// Distribution tenants allow you to create multi-tenant configurations.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateDistributionTenant(ctx context.Context, req *connect.Request[pb.UpdateDistributionTenantRequest]) (*connect.Response[pb.UpdateDistributionTenantResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.UpdateDistributionTenant is not implemented"))
}

// UpdateDistributionWithStagingConfig updates a distribution with a staging configuration.
// This allows you to associate a staging distribution with a production distribution.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateDistributionWithStagingConfig(ctx context.Context, req *connect.Request[pb.UpdateDistributionWithStagingConfigRequest]) (*connect.Response[pb.UpdateDistributionWithStagingConfigResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.UpdateDistributionWithStagingConfig is not implemented"))
}

// UpdateDomainAssociation updates a domain association for CloudFront.
// Domain associations link alternate domain names to your distribution.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateDomainAssociation(ctx context.Context, req *connect.Request[pb.UpdateDomainAssociationRequest]) (*connect.Response[pb.UpdateDomainAssociationResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.UpdateDomainAssociation is not implemented"))
}

// UpdateFieldLevelEncryptionConfig updates a field-level encryption configuration.
// Field-level encryption encrypts specific fields in the request body before forwarding to the origin.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateFieldLevelEncryptionConfig(ctx context.Context, req *connect.Request[pb.UpdateFieldLevelEncryptionConfigRequest]) (*connect.Response[pb.UpdateFieldLevelEncryptionConfigResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.UpdateFieldLevelEncryptionConfig is not implemented"))
}

// UpdateFieldLevelEncryptionProfile updates a field-level encryption profile.
// Field-level encryption profiles define which fields to encrypt and the encryption key to use.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateFieldLevelEncryptionProfile(ctx context.Context, req *connect.Request[pb.UpdateFieldLevelEncryptionProfileRequest]) (*connect.Response[pb.UpdateFieldLevelEncryptionProfileResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.UpdateFieldLevelEncryptionProfile is not implemented"))
}

// UpdateFunction updates a CloudFront function.
// Functions are lightweight scripts that can modify CloudFront behaviour at edge locations.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateFunction(ctx context.Context, req *connect.Request[pb.UpdateFunctionRequest]) (*connect.Response[pb.UpdateFunctionResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.UpdateFunction is not implemented"))
}

// UpdateKeyGroup updates a key group for CloudFront.
// Key groups are used to specify which public keys can be used for signed URLs or cookies.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateKeyGroup(ctx context.Context, req *connect.Request[pb.UpdateKeyGroupRequest]) (*connect.Response[pb.UpdateKeyGroupResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.UpdateKeyGroup is not implemented"))
}

// UpdateKeyValueStore updates a key value store for CloudFront.
// Key value stores are used to store data for CloudFront functions at edge locations.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateKeyValueStore(ctx context.Context, req *connect.Request[pb.UpdateKeyValueStoreRequest]) (*connect.Response[pb.UpdateKeyValueStoreResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.UpdateKeyValueStore is not implemented"))
}

// UpdateOriginAccessControl updates an origin access control for CloudFront.
// Origin access controls control how CloudFront accesses your origin servers.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateOriginAccessControl(ctx context.Context, req *connect.Request[pb.UpdateOriginAccessControlRequest]) (*connect.Response[pb.UpdateOriginAccessControlResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.UpdateOriginAccessControl is not implemented"))
}

// UpdateOriginRequestPolicy updates an origin request policy for CloudFront.
// Origin request policies control what headers, cookies, and query strings are forwarded to the origin.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateOriginRequestPolicy(ctx context.Context, req *connect.Request[pb.UpdateOriginRequestPolicyRequest]) (*connect.Response[pb.UpdateOriginRequestPolicyResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.UpdateOriginRequestPolicy is not implemented"))
}

// UpdatePublicKey updates a public key for CloudFront.
// Public keys are used for signed URLs and signed cookies.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdatePublicKey(ctx context.Context, req *connect.Request[pb.UpdatePublicKeyRequest]) (*connect.Response[pb.UpdatePublicKeyResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.UpdatePublicKey is not implemented"))
}

// UpdateRealtimeLogConfig updates a real-time log configuration for CloudFront.
// Real-time logs provide detailed information about requests as they pass through CloudFront.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateRealtimeLogConfig(ctx context.Context, req *connect.Request[pb.UpdateRealtimeLogConfigRequest]) (*connect.Response[pb.UpdateRealtimeLogConfigResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.UpdateRealtimeLogConfig is not implemented"))
}

// UpdateResponseHeadersPolicy updates a response headers policy for CloudFront.
// Response headers policies define which HTTP headers are included in responses.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateResponseHeadersPolicy(ctx context.Context, req *connect.Request[pb.UpdateResponseHeadersPolicyRequest]) (*connect.Response[pb.UpdateResponseHeadersPolicyResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.UpdateResponseHeadersPolicy is not implemented"))
}

// UpdateStreamingDistribution updates a streaming distribution for CloudFront.
// Streaming distributions are used to serve live streaming content.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateStreamingDistribution(ctx context.Context, req *connect.Request[pb.UpdateStreamingDistributionRequest]) (*connect.Response[pb.UpdateStreamingDistributionResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.UpdateStreamingDistribution is not implemented"))
}

// UpdateTrustStore updates a trust store for CloudFront.
// Trust stores are used to verify client certificates for mutual TLS.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateTrustStore(ctx context.Context, req *connect.Request[pb.UpdateTrustStoreRequest]) (*connect.Response[pb.UpdateTrustStoreResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.UpdateTrustStore is not implemented"))
}

// UpdateVpcOrigin updates a VPC origin for CloudFront.
// VPC origins allow CloudFront to access your origins in a VPC.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateVpcOrigin(ctx context.Context, req *connect.Request[pb.UpdateVpcOriginRequest]) (*connect.Response[pb.UpdateVpcOriginResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.UpdateVpcOrigin is not implemented"))
}

// VerifyDnsConfiguration verifies a DNS configuration for CloudFront.
// This validates that your DNS records are correctly configured for your distribution.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) VerifyDnsConfiguration(ctx context.Context, req *connect.Request[pb.VerifyDnsConfigurationRequest]) (*connect.Response[pb.VerifyDnsConfigurationResult], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("cloudfront.CloudFrontService.VerifyDnsConfiguration is not implemented"))
}
